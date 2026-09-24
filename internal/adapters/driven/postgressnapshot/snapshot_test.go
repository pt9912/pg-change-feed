package postgressnapshot

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jackc/pglogrepl"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgproto3"

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
	if raw != "-1" {
		t.Fatalf("reltuples der nie analysierten Tabelle = %s auf %s, erwartet -1", raw, scalar(t, admin, "SHOW server_version"))
	}
	rows, known, err := adapter.EstimatedRows(ctx, "public", table)
	if err != nil || known {
		t.Fatalf("EstimatedRows vor ANALYZE = %d, %v, %v; erwartet unbekannt", rows, known, err)
	}

	mustExec(t, admin, "ANALYZE "+qualified)
	rows, known, err = adapter.EstimatedRows(ctx, "public", table)
	if err != nil || !known || rows != 40 {
		t.Fatalf("EstimatedRows nach ANALYZE = %d, %v, %v; erwartet 40, bekannt", rows, known, err)
	}
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

// parityColumns trägt die Typen der Bild-Parität; `gen` ist generiert und
// fehlt in beiden Bildern.
const parityColumns = `id int PRIMARY KEY, ts timestamptz, f float8, b bytea, iv interval, n numeric,
	j jsonb, arr text[], d date, nul text, gen int GENERATED ALWAYS AS (id * 2) STORED`

const parityInsert = `INSERT INTO %s (id, ts, f, b, iv, n, j, arr, d) VALUES
	(1, '2026-09-24 08:11:12.5+00', 0.1::float8 + 0.2::float8, '\xdeadbeef', '1 day 02:03:04', 12345.6700,
	 '{"a": [1, 2], "b": null}', '{x,"y z",NULL}', '2026-09-24'),
	(2, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL)`

// TestImageParityWalAndBackfill trägt die Bild-Parität (M5): dieselben
// Zeilen ergeben über den WAL-Pfad (Publication, Slot, Walsender) und den
// Backfill-Pfad byte-gleiche Row Images — unter einer Rolle, deren
// Sitzungs-GUC vom Standard abweichen.
func TestImageParityWalAndBackfill(t *testing.T) {
	dsn := testDSN(t)
	ctx := testCtx(t)
	admin := connect(t, dsn)
	role, password := newRole(t, admin, "REPLICATION")
	for _, guc := range []string{"timezone = 'Asia/Tokyo'", "datestyle = 'German, DMY'", "intervalstyle = 'sql_standard'"} {
		mustExec(t, admin, "ALTER ROLE "+role+" SET "+guc)
	}
	table := newTable(t, admin, parityColumns)
	qualified := "public." + quoteIdent(table)
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
	mustExec(t, admin, fmt.Sprintf(parityInsert, qualified))
	if err := pglogrepl.StartReplication(ctx, wal, slot, start, pglogrepl.StartReplicationOptions{
		Mode:       pglogrepl.LogicalReplication,
		PluginArgs: []string{"proto_version '1'", "publication_names '" + publication + "'"},
	}); err != nil {
		t.Fatalf("START_REPLICATION: %v", err)
	}
	walImages := collectWalImages(t, ctx, wal, 2)

	// Backfill-Pfad: derselbe DSN, dieselbe Rolle.
	snap, err := newAdapter(t, memberDSN).OpenSnapshot(ctx, "parity"+suffix(), "public", table)
	if err != nil {
		t.Fatalf("OpenSnapshot: %v", err)
	}
	t.Cleanup(func() { _ = snap.Close(ctx) })
	_, rows := drain(t, ctx, snap)
	backfillImages := map[string][]byte{}
	for _, row := range rows {
		image, err := model.BuildRowImage(snap.Columns(), row, nil)
		if err != nil {
			t.Fatalf("BuildRowImage: %v", err)
		}
		backfillImages[*row[0]] = image
	}

	for id, walImage := range walImages {
		backfillImage, ok := backfillImages[id]
		if !ok {
			t.Errorf("Zeile %s fehlt im Backfill", id)
			continue
		}
		if !bytes.Equal(walImage, backfillImage) {
			t.Errorf("Zeile %s: WAL-Bild %s ≠ Backfill-Bild %s", id, walImage, backfillImage)
		}
	}
	if len(backfillImages) != len(walImages) {
		t.Errorf("Backfill trägt %d Zeilen, WAL %d", len(backfillImages), len(walImages))
	}

	// Die Bindung an die Aussage: die Rollen-GUC prägen das WAL-Bild
	// tatsächlich, die generierte Spalte fehlt, NULL entfällt.
	image := string(walImages["1"])
	if !strings.Contains(image, "JST") || !regexp.MustCompile(`"d":"\d{2}\.\d{2}\.\d{4}"`).MatchString(image) || !strings.Contains(image, `"iv":"1 2:03:04"`) {
		t.Errorf("das WAL-Bild trägt die Rollen-GUC nicht: %s", image)
	}
	if strings.Contains(image, `"gen"`) || strings.Contains(image, `"nul"`) {
		t.Errorf("das Bild trägt die generierte oder die NULL-Spalte: %s", image)
	}
	if got := string(walImages["2"]); got != `{"id":"2"}` {
		t.Errorf("Bild der Zeile 2 = %s, erwartet {\"id\":\"2\"}", got)
	}
	t.Logf("Bild Zeile 1 auf %s: %s", scalar(t, admin, "SHOW server_version"), image)
}

// collectWalImages liest Insert-Nachrichten der Publication vom
// Walsender und baut je Zeile das Bild mit derselben Funktion wie der
// Backfill; der Schlüssel ist der Wert der ersten Spalte.
func collectWalImages(t *testing.T, ctx context.Context, wal *pgconn.PgConn, want int) map[string][]byte {
	t.Helper()
	relations := map[uint32]*pglogrepl.RelationMessage{}
	images := map[string][]byte{}
	for len(images) < want {
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
			columns := make([]string, len(m.Tuple.Columns))
			values := make([]*string, len(m.Tuple.Columns))
			for i, column := range m.Tuple.Columns {
				columns[i] = relation.Columns[i].Name
				if column.DataType == pglogrepl.TupleDataTypeText {
					text := string(column.Data)
					values[i] = &text
				}
			}
			image, err := model.BuildRowImage(columns, values, nil)
			if err != nil {
				t.Fatalf("BuildRowImage: %v", err)
			}
			images[*values[0]] = image
		}
	}
	return images
}

// TestSlotReserveExhaustedIsConfiguration trägt die Betriebs-Vorbedingung
// der Slot-Reserve: sind alle Slots belegt, endet die Anlage als
// `configuration`. Der Test braucht einen eigenen PostgreSQL mit
// `max_replication_slots=1` — auf dem gemeinsamen Testcontainer würde er
// die Slots paralleler Pakete blockieren — und überspringt ohne
// `CDC_SNAPSHOT_TEST_EXCLUSIVE_DSN`.
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
	_, err = newAdapter(t, dsn).OpenSnapshot(ctx, "reserve"+suffix(), "public", "t")
	if !errors.Is(err, outbound.ErrSnapshotConfiguration) {
		t.Fatalf("OpenSnapshot ohne freien Slot: %v, erwartet Klasse configuration", err)
	}
}
