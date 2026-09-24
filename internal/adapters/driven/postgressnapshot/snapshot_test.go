package postgressnapshot

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jackc/pglogrepl"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgproto3"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgressnapshot/snapshotlogic"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Jeder Test dieses Pakets setzt einen PostgreSQL mit `wal_level=logical`
// voraus und überspringt ohne `CDC_REPLICATION_TEST_DSN` (`make
// test-replication`): das Paket steht aus dem Nenner des Coverage-Gates
// heraus (`ADR-0071` Punkt 1), solange kein Test netzlos läuft. Die
// Tests tragen die Messungen M1–M5 aus `ADR-0111`.

const (
	dsnEnv          = "CDC_REPLICATION_TEST_DSN"
	exclusiveDSNEnv = "CDC_SNAPSHOT_TEST_EXCLUSIVE_DSN"
)

var nameCounter atomic.Int64

// quoteIdent ist die Quoting-Funktion des Adapters, hier für die
// Test-Statements.
var quoteIdent = snapshotlogic.QuoteIdent

// suffix liefert einen je Prozess und Aufruf eindeutigen Namensteil im
// Bezeichner-Alphabet der Quelle.
func suffix() string {
	return fmt.Sprintf("%d_%d", os.Getpid(), nameCounter.Add(1))
}

// testDSN liefert den DSN des Testcontainers; ohne ihn überspringt der
// Test.
func testDSN(t *testing.T) string {
	t.Helper()
	dsn := os.Getenv(dsnEnv)
	if dsn == "" {
		t.Skip(dsnEnv + " nicht gesetzt — die Snapshot-Adapter-Tests laufen über make test-replication")
	}
	return dsn
}

func testCtx(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	t.Cleanup(cancel)
	return ctx
}

// connect öffnet eine reguläre Verbindung; die Verbindung endet mit dem
// Test.
func connect(t *testing.T, dsn string) *pgconn.PgConn {
	t.Helper()
	conn, err := pgconn.Connect(context.Background(), dsn)
	if err != nil {
		t.Fatalf("Verbindungsaufbau: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close(context.Background()) })
	// Ein Aufräum-Statement wartet höchstens 10 s auf die Sperre einer
	// ungeschlossenen Lese-Transaktion.
	mustExec(t, conn, "SET lock_timeout = '10s'")
	return conn
}

func mustExec(t *testing.T, conn *pgconn.PgConn, sql string) {
	t.Helper()
	if _, err := conn.Exec(context.Background(), sql).ReadAll(); err != nil {
		t.Fatalf("%s: %v", sql, err)
	}
}

// bestEffort führt ein Aufräum-Statement ohne Fehlerprüfung aus.
func bestEffort(conn *pgconn.PgConn, sql string) {
	_, _ = conn.Exec(context.Background(), sql).ReadAll()
}

// scalar liest die erste Spalte der ersten Zeile als Text.
func scalar(t *testing.T, conn *pgconn.PgConn, sql string) string {
	t.Helper()
	results, err := conn.Exec(context.Background(), sql).ReadAll()
	if err != nil {
		t.Fatalf("%s: %v", sql, err)
	}
	if len(results) != 1 || len(results[0].Rows) != 1 {
		t.Fatalf("%s: erwartet genau eine Zeile", sql)
	}
	return string(results[0].Rows[0][0])
}

// newTable legt eine eigene Tabelle an und räumt sie ab; sie liefert
// den unqualifizierten Namen im Schema `public`.
func newTable(t *testing.T, admin *pgconn.PgConn, columns string) string {
	t.Helper()
	name := "snap_" + suffix()
	mustExec(t, admin, "CREATE TABLE public."+quoteIdent(name)+" ("+columns+")")
	t.Cleanup(func() { bestEffort(admin, "DROP TABLE IF EXISTS public."+quoteIdent(name)) })
	return name
}

// newRole legt eine Login-Rolle an und räumt sie samt ihren Rechten ab.
func newRole(t *testing.T, admin *pgconn.PgConn, attributes string) (role, password string) {
	t.Helper()
	sfx := suffix()
	role, password = "snap_role_"+sfx, "pw"+sfx
	mustExec(t, admin, "CREATE ROLE "+role+" LOGIN "+attributes+" PASSWORD '"+password+"'")
	t.Cleanup(func() {
		bestEffort(admin, "DROP OWNED BY "+role)
		bestEffort(admin, "DROP ROLE IF EXISTS "+role)
	})
	return role, password
}

// roleDSN ersetzt Nutzer und Passwort des Test-DSN.
func roleDSN(t *testing.T, dsn, role, password string) string {
	t.Helper()
	u, err := url.Parse(dsn)
	if err != nil {
		t.Fatalf("DSN: %v", err)
	}
	u.User = url.UserPassword(role, password)
	return u.String()
}

func newAdapter(t *testing.T, dsn string, opts ...Option) *PostgresTableSnapshotAdapter {
	t.Helper()
	a, err := New(dsn, opts...)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return a
}

func slotCount(t *testing.T, admin *pgconn.PgConn, slot string) string {
	t.Helper()
	return scalar(t, admin, "SELECT count(*) FROM pg_replication_slots WHERE slot_name = '"+slot+"'")
}

// waitSlotGone wartet, bis der Slot nicht mehr im Katalog steht: der
// Walsender endet asynchron zum Verbindungsende des Clients.
func waitSlotGone(t *testing.T, admin *pgconn.PgConn, slot string) {
	t.Helper()
	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		if slotCount(t, admin, slot) == "0" {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatalf("Slot %s steht 15 s nach dem Verbindungsende noch im Katalog", slot)
}

// drain liest den Snapshot bis zum Ende; es liefert die Blockgrößen samt
// dem abschließenden leeren Block und alle Zeilen.
func drain(t *testing.T, ctx context.Context, snap outbound.TableSnapshot) (sizes []int, rows [][]*string) {
	t.Helper()
	for {
		block, err := snap.NextBlock(ctx)
		if err != nil {
			t.Fatalf("NextBlock: %v", err)
		}
		sizes = append(sizes, len(block))
		if len(block) == 0 {
			return sizes, rows
		}
		rows = append(rows, block...)
	}
}

// render macht eine Zeile vergleichbar; NULL trägt ein eigenes Zeichen.
func render(row []*string) string {
	parts := make([]string, len(row))
	for i, value := range row {
		if value == nil {
			parts[i] = "∅"
		} else {
			parts[i] = *value
		}
	}
	return strings.Join(parts, "|")
}

func renderAll(rows [][]*string) []string {
	out := make([]string, len(rows))
	for i, row := range rows {
		out[i] = render(row)
	}
	sort.Strings(out)
	return out
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// TestSnapshotPairedWithPoint trägt M1: der Slot-Snapshot sieht die vor
// `X` committeten Zeilen und keine danach committete — auch keine, die
// zwischen der Slot-Anlage und dem Import committet wird. `X` steht
// zwischen den beiden Commit-Positionen.
func TestSnapshotPairedWithPoint(t *testing.T) {
	dsn := testDSN(t)
	ctx := testCtx(t)
	admin := connect(t, dsn)
	table := newTable(t, admin, "id int PRIMARY KEY, v text")
	qualified := "public." + quoteIdent(table)
	mustExec(t, admin, "INSERT INTO "+qualified+" VALUES (1, 'a'), (2, 'b')")
	beforePoint := scalar(t, admin, "SELECT pg_current_wal_lsn()")

	adapter := newAdapter(t, dsn)
	export, err := adapter.exportSnapshot(ctx, "cdc_bf_paired_"+suffix())
	if err != nil {
		t.Fatalf("exportSnapshot: %v", err)
	}
	mustExec(t, admin, "INSERT INTO "+qualified+" VALUES (3, 'c')")
	afterPoint := scalar(t, admin, "SELECT pg_current_wal_lsn()")

	snap, err := adapter.importSnapshot(ctx, export, "public", table)
	if err != nil {
		t.Fatalf("importSnapshot: %v", err)
	}
	t.Cleanup(func() { _ = snap.Close(ctx) })

	point := pglogrepl.LSN(snap.Offset()).String()
	if scalar(t, admin, "SELECT '"+beforePoint+"'::pg_lsn <= '"+point+"'::pg_lsn") != "t" {
		t.Errorf("X=%s liegt vor dem Commit der Zeilen 1 und 2 (%s)", point, beforePoint)
	}
	if scalar(t, admin, "SELECT '"+point+"'::pg_lsn < '"+afterPoint+"'::pg_lsn") != "t" {
		t.Errorf("X=%s liegt nicht vor dem Commit der Zeile 3 (%s)", point, afterPoint)
	}
	_, rows := drain(t, ctx, snap)
	if got, want := renderAll(rows), []string{"1|a", "2|b"}; !equalStrings(got, want) {
		t.Fatalf("Snapshot-Zeilen = %v, erwartet %v (Zeile 3 committet nach X)", got, want)
	}
}

// TestTemporarySlotEndsWithConnection trägt M4: der Slot besteht, solange
// die Replication-Verbindung den Snapshot hält, ist nach Import und
// Verbindungsende weg, und der Cursor liest weiter den Snapshot-Stand —
// weder die danach geänderte noch die danach eingefügte Zeile.
func TestTemporarySlotEndsWithConnection(t *testing.T) {
	dsn := testDSN(t)
	ctx := testCtx(t)
	admin := connect(t, dsn)
	table := newTable(t, admin, "id int PRIMARY KEY, v text")
	qualified := "public." + quoteIdent(table)
	mustExec(t, admin, "INSERT INTO "+qualified+" VALUES (1, 'a'), (2, 'b')")

	slot := "cdc_bf_ends_" + suffix()
	adapter := newAdapter(t, dsn)
	export, err := adapter.exportSnapshot(ctx, slot)
	if err != nil {
		t.Fatalf("exportSnapshot: %v", err)
	}
	if got := slotCount(t, admin, slot); got != "1" {
		t.Fatalf("Slot %s zwischen Anlage und Import: %s Zeilen im Katalog, erwartet 1", slot, got)
	}
	snap, err := adapter.importSnapshot(ctx, export, "public", table)
	if err != nil {
		t.Fatalf("importSnapshot: %v", err)
	}
	t.Cleanup(func() { _ = snap.Close(ctx) })
	waitSlotGone(t, admin, slot)

	mustExec(t, admin, "UPDATE "+qualified+" SET v = 'geändert' WHERE id = 1")
	mustExec(t, admin, "INSERT INTO "+qualified+" VALUES (3, 'c')")
	_, rows := drain(t, ctx, snap)
	if got, want := renderAll(rows), []string{"1|a", "2|b"}; !equalStrings(got, want) {
		t.Fatalf("Snapshot-Zeilen nach Verbindungsende = %v, erwartet %v", got, want)
	}
}

// TestOpenSnapshotLeavesNoSlot trägt den Regelweg über `OpenSnapshot`:
// Run-Kennung ohne Bindestriche im Slot-Namen, kein Slot nach der
// Rückkehr.
func TestOpenSnapshotLeavesNoSlot(t *testing.T) {
	dsn := testDSN(t)
	ctx := testCtx(t)
	admin := connect(t, dsn)
	table := newTable(t, admin, "id int PRIMARY KEY")
	mustExec(t, admin, "INSERT INTO public."+quoteIdent(table)+" VALUES (1)")

	sfx := suffix()
	adapter := newAdapter(t, dsn)
	snap, err := adapter.OpenSnapshot(ctx, "0a1b-2c3d-"+sfx, "public", table)
	if err != nil {
		t.Fatalf("OpenSnapshot: %v", err)
	}
	t.Cleanup(func() { _ = snap.Close(ctx) })
	waitSlotGone(t, admin, "cdc_bf_0a1b2c3d"+sfx)
	if snap.Offset() == 0 {
		t.Fatal("Offset ist 0")
	}
}

// TestPermissionClassWithoutSelect trägt M2: eine Rolle nur mit `LOGIN
// REPLICATION` legt den temporären Slot an, scheitert aber ohne `SELECT`
// an der Tabelle mit der Klasse `permission` und hinterlässt keinen Slot;
// mit dem Grant gelingt derselbe Aufruf.
func TestPermissionClassWithoutSelect(t *testing.T) {
	dsn := testDSN(t)
	ctx := testCtx(t)
	admin := connect(t, dsn)
	role, password := newRole(t, admin, "REPLICATION")
	table := newTable(t, admin, "id int PRIMARY KEY")
	mustExec(t, admin, "INSERT INTO public."+quoteIdent(table)+" VALUES (1)")

	adapter := newAdapter(t, roleDSN(t, dsn, role, password))
	runID := "perm" + suffix()
	_, err := adapter.OpenSnapshot(ctx, runID, "public", table)
	if !errors.Is(err, outbound.ErrSnapshotPermission) {
		t.Fatalf("OpenSnapshot ohne SELECT: %v, erwartet Klasse permission", err)
	}
	waitSlotGone(t, admin, "cdc_bf_"+runID)

	mustExec(t, admin, "GRANT SELECT ON public."+quoteIdent(table)+" TO "+role)
	snap, err := adapter.OpenSnapshot(ctx, runID, "public", table)
	if err != nil {
		t.Fatalf("OpenSnapshot mit SELECT: %v", err)
	}
	t.Cleanup(func() { _ = snap.Close(ctx) })
	if _, rows := drain(t, ctx, snap); len(rows) != 1 {
		t.Fatalf("Zeilen mit SELECT = %d, erwartet 1", len(rows))
	}
}

// TestWrongPasswordIsPermission trägt die Anmelde-Klasse: ein falsches
// Passwort endet als `permission`, nicht als vorübergehender Fehler.
func TestWrongPasswordIsPermission(t *testing.T) {
	dsn := testDSN(t)
	ctx := testCtx(t)
	admin := connect(t, dsn)
	role, _ := newRole(t, admin, "REPLICATION")
	table := newTable(t, admin, "id int PRIMARY KEY")

	adapter := newAdapter(t, roleDSN(t, dsn, role, "falsch"))
	_, err := adapter.OpenSnapshot(ctx, "pw"+suffix(), "public", table)
	if !errors.Is(err, outbound.ErrSnapshotPermission) {
		t.Fatalf("OpenSnapshot mit falschem Passwort: %v, erwartet Klasse permission", err)
	}
}

// TestSlotCreationTimeout trägt M3: die Slot-Anlage wartet auf die beim
// Aufruf laufende Schreibtransaktion; das Zeitlimit beendet den Versuch
// als `transient`, der abgebrochene Versuch hinterlässt keinen Slot, und
// nach dem Ende der Schreibtransaktion gelingt derselbe Aufruf.
func TestSlotCreationTimeout(t *testing.T) {
	dsn := testDSN(t)
	ctx := testCtx(t)
	admin := connect(t, dsn)
	writer := connect(t, dsn)
	table := newTable(t, admin, "id int PRIMARY KEY")
	mustExec(t, writer, "BEGIN")
	mustExec(t, writer, "INSERT INTO public."+quoteIdent(table)+" VALUES (1)")
	committed := false
	t.Cleanup(func() {
		if !committed {
			bestEffort(writer, "ROLLBACK")
		}
	})

	runID := "timeout" + suffix()
	adapter := newAdapter(t, dsn, WithSlotTimeout(time.Second))
	type outcome struct {
		err     error
		elapsed time.Duration
	}
	done := make(chan outcome, 1)
	start := time.Now()
	go func() {
		_, err := adapter.OpenSnapshot(ctx, runID, "public", table)
		done <- outcome{err: err, elapsed: time.Since(start)}
	}()
	var got outcome
	select {
	case got = <-done:
	case <-time.After(20 * time.Second):
		t.Fatal("OpenSnapshot kehrt trotz Zeitlimit von 1 s nicht binnen 20 s zurück")
	}
	if !errors.Is(got.err, outbound.ErrSnapshotTransient) {
		t.Fatalf("OpenSnapshot gegen offene Schreibtransaktion: %v, erwartet Klasse transient", got.err)
	}
	if got.elapsed < 900*time.Millisecond {
		t.Errorf("OpenSnapshot endete nach %s, vor dem Zeitlimit von 1 s", got.elapsed)
	}

	mustExec(t, writer, "COMMIT")
	committed = true
	waitSlotGone(t, admin, "cdc_bf_"+runID)
	snap, err := newAdapter(t, dsn).OpenSnapshot(ctx, runID, "public", table)
	if err != nil {
		t.Fatalf("OpenSnapshot nach dem Ende der Schreibtransaktion: %v", err)
	}
	t.Cleanup(func() { _ = snap.Close(ctx) })
	if _, rows := drain(t, ctx, snap); len(rows) != 1 {
		t.Fatalf("Zeilen nach dem Ende der Schreibtransaktion = %d, erwartet 1", len(rows))
	}
}

// TestBlocks trägt die Blockbildung: Blöcke von höchstens `B` Zeilen,
// exakt an der Blockgrenze kein Rest-Block, eine leere Tabelle liefert
// keinen Block, und das Ende bleibt nach dem ersten leeren Block leer.
func TestBlocks(t *testing.T) {
	dsn := testDSN(t)
	ctx := testCtx(t)
	admin := connect(t, dsn)
	cases := []struct {
		name  string
		rows  int
		sizes []int
	}{
		{"mehrere Blöcke mit Rest", 7, []int{3, 3, 1, 0}},
		{"exakt an der Blockgrenze", 6, []int{3, 3, 0}},
		{"weniger als ein Block", 2, []int{2, 0}},
		{"leere Tabelle", 0, []int{0}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			table := newTable(t, admin, "id int PRIMARY KEY, v text")
			if tc.rows > 0 {
				mustExec(t, admin, fmt.Sprintf("INSERT INTO public.%s SELECT g, 'v'||g FROM generate_series(1, %d) g", quoteIdent(table), tc.rows))
			}
			snap, err := newAdapter(t, dsn, WithBlockSize(3)).OpenSnapshot(ctx, "blocks"+suffix(), "public", table)
			if err != nil {
				t.Fatalf("OpenSnapshot: %v", err)
			}
			t.Cleanup(func() { _ = snap.Close(ctx) })
			sizes, rows := drain(t, ctx, snap)
			if fmt.Sprint(sizes) != fmt.Sprint(tc.sizes) {
				t.Fatalf("Blockgrößen = %v, erwartet %v", sizes, tc.sizes)
			}
			if len(rows) != tc.rows {
				t.Fatalf("Zeilen = %d, erwartet %d", len(rows), tc.rows)
			}
			again, err := snap.NextBlock(ctx)
			if err != nil || len(again) != 0 {
				t.Fatalf("NextBlock nach dem Ende = %d Zeilen, %v; erwartet leer", len(again), err)
			}
		})
	}
}

// TestColumnsSkipDroppedAndGenerated trägt die Spaltenliste: gelöschte und
// generierte Spalten fehlen, ein Name mit Sonderzeichen wird gelesen, NULL
// bleibt NULL, der Text-Stand der Quelle geht unverändert durch.
func TestColumnsSkipDroppedAndGenerated(t *testing.T) {
	dsn := testDSN(t)
	ctx := testCtx(t)
	admin := connect(t, dsn)
	table := newTable(t, admin, `"Id" int PRIMARY KEY, "we ""ird" text, dropme int, gen int GENERATED ALWAYS AS ("Id" * 2) STORED, n text`)
	qualified := "public." + quoteIdent(table)
	mustExec(t, admin, "ALTER TABLE "+qualified+" DROP COLUMN dropme")
	mustExec(t, admin, "INSERT INTO "+qualified+` ("Id", "we ""ird") VALUES (1, 'x')`)

	snap, err := newAdapter(t, dsn).OpenSnapshot(ctx, "cols"+suffix(), "public", table)
	if err != nil {
		t.Fatalf("OpenSnapshot: %v", err)
	}
	t.Cleanup(func() { _ = snap.Close(ctx) })
	if got, want := snap.Columns(), []string{"Id", `we "ird`, "n"}; !equalStrings(got, want) {
		t.Fatalf("Spalten = %q, erwartet %q", got, want)
	}
	_, rows := drain(t, ctx, snap)
	if got, want := renderAll(rows), []string{"1|x|∅"}; !equalStrings(got, want) {
		t.Fatalf("Zeilen = %v, erwartet %v", got, want)
	}
}

// TestEstimatedRows trägt die Katalog-Schätzung: eine frisch befüllte,
// nie analysierte Tabelle trägt `reltuples = -1` (unbekannt, gemessen an
// der Version dieses Servers), nach `ANALYZE` liegt die Schätzung bei der
// tatsächlichen Zeilenzahl.
func TestEstimatedRows(t *testing.T) {
	dsn := testDSN(t)
	ctx := testCtx(t)
	admin := connect(t, dsn)
	table := newTable(t, admin, "id int PRIMARY KEY, v text")
	qualified := "public." + quoteIdent(table)
	mustExec(t, admin, "INSERT INTO "+qualified+" SELECT g, 'v' FROM generate_series(1, 40) g")
	adapter := newAdapter(t, dsn)

	raw := scalar(t, admin, "SELECT reltuples::text FROM pg_class WHERE oid = '"+qualified+"'::regclass")
	version := scalar(t, admin, "SHOW server_version")
	if raw != "-1" {
		t.Fatalf("reltuples der nie analysierten Tabelle = %s auf %s, erwartet -1", raw, version)
	}
	t.Logf("reltuples der nie analysierten Tabelle = %s auf PostgreSQL %s", raw, version)
	rows, known, err := adapter.EstimatedRows(ctx, "public", table)
	if err != nil || known {
		t.Fatalf("EstimatedRows vor ANALYZE = %d, %v, %v; erwartet unbekannt", rows, known, err)
	}

	mustExec(t, admin, "ANALYZE "+qualified)
	rows, known, err = adapter.EstimatedRows(ctx, "public", table)
	if err != nil || !known || rows != 40 {
		t.Fatalf("EstimatedRows nach ANALYZE = %d, %v, %v; erwartet 40, bekannt", rows, known, err)
	}
	t.Logf("EstimatedRows nach ANALYZE = %d, bekannt %v", rows, known)
}

// TestConfigurationClass trägt die Klasse `configuration`: eine nicht
// vorhandene Tabelle, eine ungültige Kennung und eine ungültige
// Konstruktion enden vor jedem Lesen — und ohne Slot.
func TestConfigurationClass(t *testing.T) {
	dsn := testDSN(t)
	ctx := testCtx(t)
	admin := connect(t, dsn)
	adapter := newAdapter(t, dsn)

	runID := "conf" + suffix()
	if _, err := adapter.OpenSnapshot(ctx, runID, "public", "gibt_es_nicht_"+suffix()); !errors.Is(err, outbound.ErrSnapshotConfiguration) {
		t.Errorf("OpenSnapshot einer fehlenden Tabelle: %v, erwartet Klasse configuration", err)
	}
	waitSlotGone(t, admin, "cdc_bf_"+runID)
	if _, _, err := adapter.EstimatedRows(ctx, "public", "gibt_es_nicht_"+suffix()); !errors.Is(err, outbound.ErrSnapshotConfiguration) {
		t.Errorf("EstimatedRows einer fehlenden Tabelle: %v, erwartet Klasse configuration", err)
	}
	for _, bad := range []string{"", "GROSS", "mit leerzeichen", strings.Repeat("a", 60), "-"} {
		if _, err := adapter.OpenSnapshot(ctx, bad, "public", "t"); !errors.Is(err, outbound.ErrSnapshotConfiguration) {
			t.Errorf("OpenSnapshot mit Run-Kennung %q: %v, erwartet Klasse configuration", bad, err)
		}
	}
	if _, err := adapter.OpenSnapshot(ctx, runID, "", "t"); !errors.Is(err, outbound.ErrSnapshotConfiguration) {
		t.Errorf("OpenSnapshot ohne Schema: %v, erwartet Klasse configuration", err)
	}
	if _, _, err := adapter.EstimatedRows(ctx, "public", ""); !errors.Is(err, outbound.ErrSnapshotConfiguration) {
		t.Errorf("EstimatedRows ohne Tabelle: %v, erwartet Klasse configuration", err)
	}
	for name, build := range map[string]func() error{
		"leerer DSN":       func() error { _, err := New(""); return err },
		"DSN nicht lesbar": func() error { _, err := New("postgres://%zz"); return err },
		"Blockgröße 0":     func() error { _, err := New(dsn, WithBlockSize(0)); return err },
		"Zeitlimit 0":      func() error { _, err := New(dsn, WithSlotTimeout(0)); return err },
	} {
		if err := build(); !errors.Is(err, outbound.ErrSnapshotConfiguration) {
			t.Errorf("New mit %s: %v, erwartet Klasse configuration", name, err)
		}
	}
}

// TestNextBlockAfterClose trägt das Ende der Lese-Transaktion: ein
// geschlossener Snapshot liest nicht mehr, ein zweites Close bleibt ohne
// Wirkung.
func TestNextBlockAfterClose(t *testing.T) {
	dsn := testDSN(t)
	ctx := testCtx(t)
	admin := connect(t, dsn)
	table := newTable(t, admin, "id int PRIMARY KEY")
	snap, err := newAdapter(t, dsn).OpenSnapshot(ctx, "close"+suffix(), "public", table)
	if err != nil {
		t.Fatalf("OpenSnapshot: %v", err)
	}
	if err := snap.Close(ctx); err != nil {
		t.Fatalf("Close: %v", err)
	}
	if err := snap.Close(ctx); err != nil {
		t.Fatalf("zweites Close: %v", err)
	}
	if _, err := snap.NextBlock(ctx); !errors.Is(err, outbound.ErrSnapshotStorage) {
		t.Fatalf("NextBlock nach Close: %v, erwartet Klasse storage", err)
	}
}

// gucLage ist eine Sitzungs-GUC-Lage der Rolle, unter der WAL-Pfad und
// Backfill-Pfad lesen; `inImage` sind Teilzeichenfolgen, die das WAL-Bild
// nur trägt, wenn die Lage tatsächlich wirkt.
type gucLage struct {
	name     string
	settings []string
	inImage  []string
}

var gucLagen = []gucLage{
	{name: "Standard"},
	{
		name: "Tokyo",
		settings: []string{
			"timezone = 'Asia/Tokyo'", "datestyle = 'German, DMY'", "intervalstyle = 'sql_standard'",
			"bytea_output = 'escape'", "extra_float_digits = 0",
		},
		inImage: []string{"JST", `+1-2 +3 +4:05:06.7`, `\\336\\255\\276\\357`},
	},
	{
		name: "Verbose",
		settings: []string{
			"intervalstyle = 'postgres_verbose'", "datestyle = 'SQL, MDY'", "extra_float_digits = -3",
			"timezone = 'America/New_York'",
		},
		inImage: []string{"EDT", "@ 1 year 2 mons 3 days 4 hours 5 mins 6.7 secs"},
	},
}

// TestImageParityWalAndBackfill trägt die Typ-Parität (`ADR-0115`
// Festlegung 4, M5 aus `ADR-0111`): über die Typ-Tabelle ist der Roh-Text
// jeder Spalte im Backfill-Pfad byte-gleich dem des WAL-Pfads (Publication,
// Slot, Walsender), und das über `BuildRowImage` gebaute Bild ebenso — in
// drei Sitzungs-GUC-Lagen der Rolle.
func TestImageParityWalAndBackfill(t *testing.T) {
	dsn := testDSN(t)
	admin := connect(t, dsn)
	set := typeTable(suffix())
	for _, statement := range set.create {
		mustExec(t, admin, statement)
	}
	t.Cleanup(func() {
		for _, statement := range set.drop {
			bestEffort(admin, statement)
		}
	})
	versionNum, err := strconv.Atoi(scalar(t, admin, "SHOW server_version_num"))
	if err != nil {
		t.Fatalf("server_version_num: %v", err)
	}
	t.Logf("PostgreSQL %s, %d Typ-Spalten", scalar(t, admin, "SHOW server_version"), len(set.columns))
	for _, lage := range gucLagen {
		t.Run(lage.name, func(t *testing.T) { checkParity(t, dsn, admin, set, versionNum >= 180000, lage) })
	}
}

func checkParity(t *testing.T, dsn string, admin *pgconn.PgConn, set typeSet, virtualGenerated bool, lage gucLage) {
	ctx := testCtx(t)
	role, password := newRole(t, admin, "REPLICATION")
	for _, setting := range lage.settings {
		mustExec(t, admin, "ALTER ROLE "+role+" SET "+setting)
	}
	table := newTable(t, admin, set.definition(virtualGenerated))
	qualified := "public." + quoteIdent(table)
	mustExec(t, admin, "ALTER TABLE "+qualified+" DROP COLUMN dropme")
	mustExec(t, admin, "GRANT SELECT ON "+qualified+" TO "+role)
	publication := "snap_pub_" + suffix()
	mustExec(t, admin, "CREATE PUBLICATION "+publication+" FOR TABLE "+qualified)
	t.Cleanup(func() { bestEffort(admin, "DROP PUBLICATION IF EXISTS "+publication) })
	memberDSN := roleDSN(t, dsn, role, password)

	// WAL-Pfad: Slot vor den Inserts, Walsender der Rolle.
	config, err := pgconn.ParseConfig(memberDSN)
	if err != nil {
		t.Fatalf("DSN: %v", err)
	}
	config.RuntimeParams["replication"] = "database"
	wal, err := pgconn.ConnectConfig(ctx, config)
	if err != nil {
		t.Fatalf("Replication-Verbindung: %v", err)
	}
	t.Cleanup(func() { _ = wal.Close(context.Background()) })
	slot := "snap_wal_" + suffix()
	created, err := pglogrepl.CreateReplicationSlot(ctx, wal, slot, outputPlugin, pglogrepl.CreateReplicationSlotOptions{
		Temporary: true, SnapshotAction: "NOEXPORT_SNAPSHOT", Mode: pglogrepl.LogicalReplication,
	})
	if err != nil {
		t.Fatalf("Slot-Anlage: %v", err)
	}
	start, err := pglogrepl.ParseLSN(created.ConsistentPoint)
	if err != nil {
		t.Fatalf("consistent_point: %v", err)
	}
	mustExec(t, admin, set.insert(qualified))
	if err := pglogrepl.StartReplication(ctx, wal, slot, start, pglogrepl.StartReplicationOptions{
		Mode:       pglogrepl.LogicalReplication,
		PluginArgs: []string{"proto_version '1'", "publication_names '" + publication + "'"},
	}); err != nil {
		t.Fatalf("START_REPLICATION: %v", err)
	}
	walRows := collectWalRows(t, ctx, wal, 3)

	// Backfill-Pfad: derselbe DSN, dieselbe Rolle.
	snap, err := newAdapter(t, memberDSN).OpenSnapshot(ctx, "parity"+suffix(), "public", table)
	if err != nil {
		t.Fatalf("OpenSnapshot: %v", err)
	}
	t.Cleanup(func() { _ = snap.Close(ctx) })
	_, rows := drain(t, ctx, snap)
	if len(rows) != len(walRows) {
		t.Fatalf("Backfill trägt %d Zeilen, WAL %d", len(rows), len(walRows))
	}

	columns := snap.Columns()
	if want := len(set.columns) + 1; len(columns) != want {
		t.Fatalf("Backfill-Spalten = %d (%q), erwartet %d: die gelöschte und die generierten Spalten fehlen", len(columns), columns, want)
	}
	images := map[string][]byte{}
	for _, row := range rows {
		id := *row[0]
		walRow, ok := walRows[id]
		if !ok {
			t.Errorf("Zeile %s fehlt im WAL", id)
			continue
		}
		if !equalStrings(walRow.columns, columns) {
			t.Fatalf("Spalten des WAL %q ≠ Spalten des Backfill %q", walRow.columns, columns)
		}
		for i, name := range columns {
			if !sameText(walRow.values[i], row[i]) {
				t.Errorf("Zeile %s, Spalte %s (Typ %s): WAL %s ≠ Backfill %s", id, name, set.ddlOf(name), showText(walRow.values[i]), showText(row[i]))
			}
		}
		image, err := model.BuildRowImage(columns, row, nil)
		if err != nil {
			t.Fatalf("BuildRowImage: %v", err)
		}
		images[id] = image
		walImage, err := model.BuildRowImage(walRow.columns, walRow.values, nil)
		if err != nil {
			t.Fatalf("BuildRowImage (WAL): %v", err)
		}
		if !bytes.Equal(walImage, image) {
			t.Errorf("Zeile %s: WAL-Bild %s ≠ Backfill-Bild %s", id, walImage, image)
		}
	}

	// Die Bindung an die Aussage: die Rollen-GUC prägen das Bild tatsächlich,
	// NULL entfällt, die gelöschte und die generierten Spalten fehlen.
	image := string(images["1"])
	for _, want := range lage.inImage {
		if !strings.Contains(image, want) {
			t.Errorf("das Bild der Lage %s trägt %q nicht: %s", lage.name, want, image)
		}
	}
	for _, absent := range []string{`"gen"`, `"genv"`, `"dropme"`} {
		if strings.Contains(image, absent) {
			t.Errorf("das Bild trägt %s: %s", absent, image)
		}
	}
	if got := string(images["3"]); got != `{"id":"3"}` {
		t.Errorf("Bild der NULL-Zeile = %s, erwartet {\"id\":\"3\"}", got)
	}
	t.Logf("Bild Zeile 1 in der Lage %s: %s", lage.name, image)
}

func sameText(a, b *string) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

func showText(value *string) string {
	if value == nil {
		return "NULL"
	}
	return fmt.Sprintf("%q", *value)
}

// walRow trägt die Spalten und Textwerte einer Insert-Nachricht.
type walRow struct {
	columns []string
	values  []*string
}

// collectWalRows liest Insert-Nachrichten der Publication vom Walsender; der
// Schlüssel ist der Wert der ersten Spalte. Ein Wert ist der Text des
// Feldes, `nil` ist NULL.
func collectWalRows(t *testing.T, ctx context.Context, wal *pgconn.PgConn, want int) map[string]walRow {
	t.Helper()
	relations := map[uint32]*pglogrepl.RelationMessage{}
	rows := map[string]walRow{}
	for len(rows) < want {
		receiveCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
		message, err := wal.ReceiveMessage(receiveCtx)
		cancel()
		if err != nil {
			t.Fatalf("ReceiveMessage: %v", err)
		}
		data, ok := message.(*pgproto3.CopyData)
		if !ok || len(data.Data) == 0 || data.Data[0] != pglogrepl.XLogDataByteID {
			continue
		}
		xlog, err := pglogrepl.ParseXLogData(data.Data[1:])
		if err != nil {
			t.Fatalf("ParseXLogData: %v", err)
		}
		logical, err := pglogrepl.Parse(xlog.WALData)
		if err != nil {
			t.Fatalf("Parse: %v", err)
		}
		switch m := logical.(type) {
		case *pglogrepl.RelationMessage:
			relations[m.RelationID] = m
		case *pglogrepl.InsertMessage:
			relation := relations[m.RelationID]
			if relation == nil {
				t.Fatalf("Insert ohne Relation %d", m.RelationID)
			}
			row := walRow{columns: make([]string, len(m.Tuple.Columns)), values: make([]*string, len(m.Tuple.Columns))}
			for i, column := range m.Tuple.Columns {
				row.columns[i] = relation.Columns[i].Name
				if column.DataType == pglogrepl.TupleDataTypeText {
					text := string(column.Data)
					row.values[i] = &text
				}
			}
			rows[*row.values[0]] = row
		}
	}
	return rows
}

// TestSlotReserveExhaustedIsConfiguration trägt die Betriebs-Vorbedingung
// der Slot-Reserve: sind alle Slots belegt, endet die Anlage als
// `configuration`. Der Test braucht einen eigenen PostgreSQL mit
// `max_replication_slots=1`, weil er auf dem gemeinsamen Testcontainer die
// Slots paralleler Pakete blockiert; den Container und die Variable
// `CDC_SNAPSHOT_TEST_EXCLUSIVE_DSN` stellt die Phase `tier` von
// `tools/harness/run-replication-tests.sh`, ohne sie überspringt der Test.
func TestSlotReserveExhaustedIsConfiguration(t *testing.T) {
	dsn := os.Getenv(exclusiveDSNEnv)
	if dsn == "" {
		t.Skip(exclusiveDSNEnv + " nicht gesetzt — die Slot-Reserve braucht einen eigenen PostgreSQL mit max_replication_slots=1")
	}
	ctx := testCtx(t)
	config, err := pgconn.ParseConfig(dsn)
	if err != nil {
		t.Fatalf("DSN: %v", err)
	}
	config.RuntimeParams["replication"] = "database"
	holder, err := pgconn.ConnectConfig(ctx, config)
	if err != nil {
		t.Fatalf("Replication-Verbindung: %v", err)
	}
	t.Cleanup(func() { _ = holder.Close(context.Background()) })
	if _, err := pglogrepl.CreateReplicationSlot(ctx, holder, "snap_holder", outputPlugin, pglogrepl.CreateReplicationSlotOptions{
		Temporary: true, SnapshotAction: "NOEXPORT_SNAPSHOT", Mode: pglogrepl.LogicalReplication,
	}); err != nil {
		t.Fatalf("Slot des Halters: %v", err)
	}
	admin := connect(t, dsn)
	table := newTable(t, admin, "id int PRIMARY KEY")
	snap, err := newAdapter(t, dsn).OpenSnapshot(ctx, "reserve"+suffix(), "public", table)
	if err == nil {
		_ = snap.Close(ctx)
	}
	if !errors.Is(err, outbound.ErrSnapshotConfiguration) || !strings.Contains(err.Error(), "SQLSTATE 53400") {
		t.Fatalf("OpenSnapshot ohne freien Slot: %v, erwartet Klasse configuration aus SQLSTATE 53400", err)
	}
}

// TestQuotedSchemaAndTableNames trägt das Quoting von Schema und Tabelle:
// Namen mit Anführungszeichen, Leerzeichen und Großbuchstaben werden
// gelesen, ihre Spalten- und Schätzungs-Abfragen finden sie ebenfalls.
func TestQuotedSchemaAndTableNames(t *testing.T) {
	dsn := testDSN(t)
	ctx := testCtx(t)
	admin := connect(t, dsn)
	schema := `Sch "ema ` + suffix()
	table := `Mixed Case "x ` + suffix()
	mustExec(t, admin, "CREATE SCHEMA "+quoteIdent(schema))
	t.Cleanup(func() { bestEffort(admin, "DROP SCHEMA IF EXISTS "+quoteIdent(schema)+" CASCADE") })
	mustExec(t, admin, "CREATE TABLE "+quoteIdent(schema)+"."+quoteIdent(table)+" (id int PRIMARY KEY, v text)")
	mustExec(t, admin, "INSERT INTO "+quoteIdent(schema)+"."+quoteIdent(table)+" VALUES (1, 'a'), (2, 'b')")
	mustExec(t, admin, "ANALYZE "+quoteIdent(schema)+"."+quoteIdent(table))

	adapter := newAdapter(t, dsn)
	snap, err := adapter.OpenSnapshot(ctx, "quoted"+suffix(), schema, table)
	if err != nil {
		t.Fatalf("OpenSnapshot: %v", err)
	}
	t.Cleanup(func() { _ = snap.Close(ctx) })
	_, rows := drain(t, ctx, snap)
	if got, want := renderAll(rows), []string{"1|a", "2|b"}; !equalStrings(got, want) {
		t.Fatalf("Zeilen = %v, erwartet %v", got, want)
	}
	if estimate, known, err := adapter.EstimatedRows(ctx, schema, table); err != nil || !known || estimate != 2 {
		t.Fatalf("EstimatedRows = %d, %v, %v; erwartet 2, bekannt", estimate, known, err)
	}
}

// TestEstimatedRowsAnalyzedEmptyTableIsKnownZero trägt die Grenze der
// Schätzung: eine analysierte leere Tabelle trägt `reltuples = 0` und ist
// eine bekannte Schätzung von null Zeilen, nicht „unbekannt".
func TestEstimatedRowsAnalyzedEmptyTableIsKnownZero(t *testing.T) {
	dsn := testDSN(t)
	ctx := testCtx(t)
	admin := connect(t, dsn)
	table := newTable(t, admin, "id int PRIMARY KEY")
	mustExec(t, admin, "ANALYZE public."+quoteIdent(table))
	if raw := scalar(t, admin, "SELECT reltuples::text FROM pg_class WHERE oid = 'public."+quoteIdent(table)+"'::regclass"); raw != "0" {
		t.Fatalf("reltuples der analysierten leeren Tabelle = %s, erwartet 0", raw)
	}
	rows, known, err := newAdapter(t, dsn).EstimatedRows(ctx, "public", table)
	if err != nil || !known || rows != 0 {
		t.Fatalf("EstimatedRows = %d, %v, %v; erwartet 0, bekannt", rows, known, err)
	}
}

// TestSlotNameCollisionIsReplication trägt die Klasse `replication`: ein
// bereits belegter Slot-Name ist eine Störung der Slot-Anlage ohne
// spezifischere Klasse.
func TestSlotNameCollisionIsReplication(t *testing.T) {
	dsn := testDSN(t)
	ctx := testCtx(t)
	admin := connect(t, dsn)
	table := newTable(t, admin, "id int PRIMARY KEY")
	runID := "collide" + suffix()
	holder := newAdapter(t, dsn)
	export, err := holder.exportSnapshot(ctx, "cdc_bf_"+runID)
	if err != nil {
		t.Fatalf("exportSnapshot des Halters: %v", err)
	}
	t.Cleanup(func() { closeConn(export.conn) })

	_, err = newAdapter(t, dsn).OpenSnapshot(ctx, runID, "public", table)
	if !errors.Is(err, outbound.ErrSnapshotReplication) {
		t.Fatalf("OpenSnapshot mit belegtem Slot-Namen: %v, erwartet Klasse replication", err)
	}
}

// appDSN setzt den `application_name` der Sitzungen, damit ein Test die
// Sitzungen des Adapters im Katalog findet.
func appDSN(t *testing.T, dsn, application string) string {
	t.Helper()
	u, err := url.Parse(dsn)
	if err != nil {
		t.Fatalf("DSN: %v", err)
	}
	q := u.Query()
	q.Set("application_name", application)
	u.RawQuery = q.Encode()
	return u.String()
}

// waitSessionsGone wartet, bis keine Sitzung mit dem `application_name`
// mehr im Katalog steht: das Backend endet asynchron zum Verbindungsende des
// Clients.
func waitSessionsGone(t *testing.T, admin *pgconn.PgConn, application string) {
	t.Helper()
	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		if scalar(t, admin, "SELECT count(*) FROM pg_stat_activity WHERE application_name = '"+application+"'") == "0" {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatalf("Sitzungen von %s stehen 15 s nach dem Fehler noch im Katalog", application)
}

// TestImportOfUnknownSnapshotIsStorage trägt den Abbruch des Imports: ein
// Snapshot-Name, den der Server nicht kennt, endet als Klasse `storage`, und
// beide Verbindungen des Adapters sind danach geschlossen (der Slot ist weg,
// keine Sitzung bleibt).
func TestImportOfUnknownSnapshotIsStorage(t *testing.T) {
	dsn := testDSN(t)
	ctx := testCtx(t)
	admin := connect(t, dsn)
	table := newTable(t, admin, "id int PRIMARY KEY")
	application := "snap_import_" + suffix()
	adapter := newAdapter(t, appDSN(t, dsn, application))
	slot := "cdc_bf_unknown_" + suffix()
	export, err := adapter.exportSnapshot(ctx, slot)
	if err != nil {
		t.Fatalf("exportSnapshot: %v", err)
	}
	export.snapshotName = "00000003-00000002-1"
	if _, err := adapter.importSnapshot(ctx, export, "public", table); !errors.Is(err, outbound.ErrSnapshotStorage) {
		t.Fatalf("importSnapshot mit unbekanntem Snapshot: %v, erwartet Klasse storage", err)
	}
	waitSlotGone(t, admin, slot)
	waitSessionsGone(t, admin, application)
}

// unreachableDSN ersetzt Host und Port des Test-DSN durch eine Adresse ohne
// Listener.
func unreachableDSN(t *testing.T, dsn string) string {
	t.Helper()
	u, err := url.Parse(dsn)
	if err != nil {
		t.Fatalf("DSN: %v", err)
	}
	u.Host = "127.0.0.1:1"
	return u.String()
}

// TestUnreachableSourceIsTransient trägt die Klasse `transient` an den
// Verbindungsaufbauten: eine Quelle ohne Listener endet an der
// Replication-Verbindung, an der Lese-Verbindung des Imports und an der
// Schätzung als `transient`, und die Replication-Verbindung des Imports ist
// danach geschlossen.
func TestUnreachableSourceIsTransient(t *testing.T) {
	dsn := testDSN(t)
	ctx := testCtx(t)
	admin := connect(t, dsn)
	table := newTable(t, admin, "id int PRIMARY KEY")
	broken := newAdapter(t, unreachableDSN(t, dsn))

	if _, err := broken.OpenSnapshot(ctx, "unreach"+suffix(), "public", table); !errors.Is(err, outbound.ErrSnapshotTransient) {
		t.Errorf("OpenSnapshot ohne Listener: %v, erwartet Klasse transient", err)
	}
	if _, _, err := broken.EstimatedRows(ctx, "public", table); !errors.Is(err, outbound.ErrSnapshotTransient) {
		t.Errorf("EstimatedRows ohne Listener: %v, erwartet Klasse transient", err)
	}

	slot := "cdc_bf_unreach_" + suffix()
	export, err := newAdapter(t, dsn).exportSnapshot(ctx, slot)
	if err != nil {
		t.Fatalf("exportSnapshot: %v", err)
	}
	if _, err := broken.importSnapshot(ctx, export, "public", table); !errors.Is(err, outbound.ErrSnapshotTransient) {
		t.Errorf("importSnapshot ohne Listener für die Lese-Verbindung: %v, erwartet Klasse transient", err)
	}
	waitSlotGone(t, admin, slot)
}

// TestCatalogQueryFailuresKeepTheirClass trägt die Fehler der Katalog-
// Abfragen: in einer eigenen Datenbank ohne `SELECT` auf `pg_attribute`
// und `pg_class` enden die Spaltenliste und die Schätzung mit der
// Fehlerklasse ihres SQLSTATE (`permission`), nicht als Lesefehler.
func TestCatalogQueryFailuresKeepTheirClass(t *testing.T) {
	dsn := testDSN(t)
	ctx := testCtx(t)
	admin := connect(t, dsn)
	database := "snap_cat_" + suffix()
	mustExec(t, admin, "CREATE DATABASE "+database+" TEMPLATE template0")
	t.Cleanup(func() { bestEffort(admin, "DROP DATABASE IF EXISTS "+database+" WITH (FORCE)") })
	u, err := url.Parse(dsn)
	if err != nil {
		t.Fatalf("DSN: %v", err)
	}
	u.Path = "/" + database
	isolated := connect(t, u.String())
	mustExec(t, isolated, "REVOKE SELECT ON pg_catalog.pg_attribute, pg_catalog.pg_class FROM PUBLIC")
	mustExec(t, isolated, "CREATE TABLE public.t (id int PRIMARY KEY)")
	role, password := newRole(t, admin, "REPLICATION")
	mustExec(t, isolated, "GRANT SELECT ON public.t TO "+role)
	adapter := newAdapter(t, roleDSN(t, u.String(), role, password))

	runID := "catalog" + suffix()
	_, err = adapter.OpenSnapshot(ctx, runID, "public", "t")
	if !errors.Is(err, outbound.ErrSnapshotPermission) || !strings.Contains(err.Error(), "Spaltenliste") {
		t.Errorf("OpenSnapshot ohne SELECT auf pg_attribute: %v, erwartet Klasse permission in der Phase Spaltenliste", err)
	}
	waitSlotGone(t, isolated, "cdc_bf_"+runID)
	_, _, err = adapter.EstimatedRows(ctx, "public", "t")
	if !errors.Is(err, outbound.ErrSnapshotPermission) || !strings.Contains(err.Error(), "Zeilenschätzung") {
		t.Errorf("EstimatedRows ohne SELECT auf pg_class: %v, erwartet Klasse permission in der Phase Zeilenschätzung", err)
	}
}

// TestReadAfterBackendTermination trägt den Verbindungsabbruch mitten im
// Lesen: nach dem Beenden der Lese-Sitzung endet der nächste Block als
// sichtbarer Fehler mit einer Klasse, und `Close` bleibt ohne Wirkung.
func TestReadAfterBackendTermination(t *testing.T) {
	dsn := testDSN(t)
	ctx := testCtx(t)
	admin := connect(t, dsn)
	table := newTable(t, admin, "id int PRIMARY KEY")
	mustExec(t, admin, "INSERT INTO public."+quoteIdent(table)+" SELECT g FROM generate_series(1, 5) g")
	snap, err := newAdapter(t, dsn, WithBlockSize(2)).OpenSnapshot(ctx, "term"+suffix(), "public", table)
	if err != nil {
		t.Fatalf("OpenSnapshot: %v", err)
	}
	t.Cleanup(func() { _ = snap.Close(ctx) })
	if block, err := snap.NextBlock(ctx); err != nil || len(block) != 2 {
		t.Fatalf("erster Block = %d Zeilen, %v; erwartet 2", len(block), err)
	}
	pid := snap.(*snapshot).conn.PID()
	mustExec(t, admin, fmt.Sprintf("SELECT pg_terminate_backend(%d)", pid))

	_, err = snap.NextBlock(ctx)
	t.Logf("NextBlock nach pg_terminate_backend: %v", err)
	if !errors.Is(err, outbound.ErrSnapshotTransient) {
		t.Fatalf("NextBlock nach dem Ende der Sitzung: %v, erwartet Klasse transient", err)
	}
	if _, err := snap.NextBlock(ctx); !errors.Is(err, outbound.ErrSnapshotTransient) {
		t.Fatalf("zweiter NextBlock auf der beendeten Verbindung: %v, erwartet Klasse transient", err)
	}
	if err := snap.Close(ctx); err != nil {
		t.Fatalf("Close nach dem Verbindungsabbruch: %v", err)
	}
}

// TestReadWithCancelledContextIsTransient trägt das Kontext-Ende beim
// Lesen: ein beendeter Kontext endet als Klasse `transient`.
func TestReadWithCancelledContextIsTransient(t *testing.T) {
	dsn := testDSN(t)
	ctx := testCtx(t)
	admin := connect(t, dsn)
	table := newTable(t, admin, "id int PRIMARY KEY")
	mustExec(t, admin, "INSERT INTO public."+quoteIdent(table)+" VALUES (1)")
	snap, err := newAdapter(t, dsn).OpenSnapshot(ctx, "cancel"+suffix(), "public", table)
	if err != nil {
		t.Fatalf("OpenSnapshot: %v", err)
	}
	t.Cleanup(func() { _ = snap.Close(ctx) })
	cancelled, cancel := context.WithCancel(ctx)
	cancel()
	if _, err := snap.NextBlock(cancelled); !errors.Is(err, outbound.ErrSnapshotTransient) {
		t.Fatalf("NextBlock mit beendetem Kontext: %v, erwartet Klasse transient", err)
	}
}
