// Package integration_test trägt den MVP-Integrationstest gegen die
// Compose-Umgebung (Abschnitt 1, MVP-Schnitt): PostgreSQL starten → CDC
// aktivieren → INSERT → UPDATE → DELETE → Changes lesen → Reihenfolge und
// Inhalt prüfen (`LH-FA-CAP-001`…003, `LH-QA-POR-003`). Der Lauf verdrahtet
// die realen Adapter — Replication-Stream, ChangeStore und ACK — als
// Verdrahtungs-Test der Composition-Root-Ebene (`ADR-0026`) und läuft
// Ende-zu-Ende über den echten Capture Service mit Persist-before-ACK
// (`LH-QA-REL-001.a`).
//
// Die Instanz gehört der Compose-Umgebung (`make test-integration`); der
// d-migrate-Rollout (ADR-0043) trägt die CDC-Tabellen vor dem Lauf. Ohne
// DSN überspringt der Test.
package integration_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresack"
	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage"
	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/mapper"
	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/receive"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/capture"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// mvpSource trägt die Quell-Kennung der MVP-Läufe; die Bindungs-Zeilen der
// Aktivierung tragen sie in den CDC-Referenztabellen.
const mvpSource = "src-mvp"

// mvpEnv trägt die Compose-seitige Testumgebung eines MVP-Laufs: die
// Quell-Verbindung, die Verwaltungs-Namen und der verdrahtete
// Capture-Pfad.
type mvpEnv struct {
	dsn         string
	pool        *pgxpool.Pool
	store       *postgresstorage.PostgresChangeStoreAdapter
	publication string
	slot        string
	feed        string
	tableID     model.SourceTableID
	schemaV     model.SchemaVersionID
}

// newMVPEnv aktiviert CDC an einer Tabelle der Compose-Instanz
// (`LH-FA-CFG-001`): die Publication trägt die Tabelle, die
// Bindungs-Zeilen tragen die Port-Kennungen. Die CDC-Tabellen trägt der
// d-migrate-Rollout (ADR-0043); fehlen sie, endet der Test mit dem Verweis
// auf den Rollout — der Test legt sie nicht still selbst an. Slot,
// Publication und Feed-Tabelle räumt die Cleanup-Kette ab; der Stream-Lauf
// endet vorher über seinen Kontext, damit der Slot ohne Belegung droppt.
func newMVPEnv(t *testing.T, name string, replicaIdentityFull bool) *mvpEnv {
	t.Helper()
	dsn := os.Getenv("CDC_INTEGRATION_DSN")
	if dsn == "" {
		t.Skip("CDC_INTEGRATION_DSN nicht gesetzt — MVP-Integrationstest läuft über make test-integration gegen die Compose-Umgebung")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("Verbindungsaufbau: %v", err)
	}
	t.Cleanup(pool.Close)

	var tables int
	if err := pool.QueryRow(ctx,
		"SELECT count(*) FROM information_schema.tables WHERE table_schema = current_schema() AND table_name IN ('source', 'source_table', 'schema_version', 'transaction', 'change')",
	).Scan(&tables); err != nil {
		t.Fatalf("CDC-Schema-Prüfung: %v", err)
	}
	if tables != 5 {
		t.Fatalf("CDC-Tabellen unvollständig (%d von 5) — Schema-Rollout über make schema-rollout vor dem E2E-Lauf (ADR-0043)", tables)
	}

	env := &mvpEnv{
		dsn:         dsn,
		pool:        pool,
		publication: "pub_pgc_mvp_" + name,
		slot:        "slot_pgc_mvp_" + name,
		feed:        fmt.Sprintf("public.feed_mvp_%s", name),
		tableID:     model.SourceTableID("tbl-mvp-" + name),
		schemaV:     model.SchemaVersionID("sv-mvp-" + name),
	}
	statements := []string{
		fmt.Sprintf("CREATE TABLE %s (id int PRIMARY KEY, name text)", env.feed),
		fmt.Sprintf("CREATE PUBLICATION %s FOR TABLE %s", env.publication, env.feed),
		fmt.Sprintf("INSERT INTO cdc.source (source_id, name) VALUES ('%s', 'MVP-Quelle') ON CONFLICT (source_id) DO NOTHING", mvpSource),
		fmt.Sprintf("INSERT INTO cdc.source_table (source_table_id, source_id, schema_name, table_name) VALUES ('%s', '%s', 'public', 'feed_mvp_%s')", env.tableID, mvpSource, name),
		fmt.Sprintf("INSERT INTO cdc.schema_version (schema_version_id, source_table_id, version) VALUES ('%s', '%s', 1)", env.schemaV, env.tableID),
	}
	if replicaIdentityFull {
		statements = append(statements, "ALTER TABLE "+env.feed+" REPLICA IDENTITY FULL")
	}
	for _, statement := range statements {
		if _, err := pool.Exec(ctx, statement); err != nil {
			t.Fatalf("Aktivierung: %v", err)
		}
	}
	t.Cleanup(func() {
		dropCtx, dropCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer dropCancel()
		deadline := time.Now().Add(8 * time.Second)
		for time.Now().Before(deadline) {
			_, err := pool.Exec(dropCtx, fmt.Sprintf("SELECT pg_drop_replication_slot('%s')", env.slot))
			if err == nil {
				break
			}
			time.Sleep(200 * time.Millisecond)
		}
		for _, statement := range []string{
			"DROP PUBLICATION IF EXISTS " + env.publication,
			"DROP TABLE IF EXISTS " + env.feed,
		} {
			_, _ = pool.Exec(dropCtx, statement)
		}
	})

	store, err := postgresstorage.New(ctx, dsn)
	if err != nil {
		t.Fatalf("Store: %v", err)
	}
	t.Cleanup(store.Close)
	env.store = store
	return env
}

// startCapture verdrahtet den Capture-Pfad am realen Treiber und startet
// den Stream (`ADR-0026`, `LH-QA-REL-001.a`): der ACK-Adapter trägt seine
// Bestätigung über die Stream-Verbindung, der Service orchestriert
// Persist-before-ACK. Der Rückgabe-Kanal trägt das Lauf-Ende; der Test
// beendet den Lauf über cancel.
func startCapture(t *testing.T, env *mvpEnv) (context.CancelFunc, <-chan error) {
	t.Helper()
	streamCtx, cancel := context.WithCancel(context.Background())
	stream, err := receive.NewStream(streamCtx, receive.Config{
		DSN:         env.dsn,
		Source:      mvpSource,
		Publication: env.publication,
		Slot:        env.slot,
		Tables: map[string]mapper.TableBinding{
			env.feed: {TableID: env.tableID, SchemaVersion: env.schemaV},
		},
	})
	if err != nil {
		cancel()
		t.Fatalf("NewStream: %v", err)
	}
	ack, err := postgresack.New(stream.Conn())
	if err != nil {
		cancel()
		t.Fatalf("ACK-Adapter: %v", err)
	}
	if err := stream.BindCapture(capture.NewCaptureService(env.store, ack)); err != nil {
		cancel()
		t.Fatalf("BindCapture: %v", err)
	}
	runDone := make(chan error, 1)
	go func() {
		runDone <- stream.Run(streamCtx)
	}()
	return cancel, runDone
}

// awaitPersistedChanges liest den Store, bis die erwarteten Changes der
// Tabelle dieser Umgebungen persistiert sind (Polling mit Test-Zeitgrenze);
// der Erfassungsweg läuft asynchron über den Stream. Der Tabellenfilter
// trennt die Läufe derselben Quelle.
func awaitPersistedChanges(t *testing.T, env *mvpEnv, limit int) []outbound.ChangeRecord {
	t.Helper()
	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		records, err := env.store.ReadChanges(context.Background(), outbound.ChangeQuery{
			Source: mvpSource, Table: &env.tableID,
		})
		if err != nil {
			t.Fatalf("ReadChanges: %v", err)
		}
		if len(records) >= limit {
			return records
		}
		time.Sleep(100 * time.Millisecond)
	}
	records, err := env.store.ReadChanges(context.Background(), outbound.ChangeQuery{
		Source: mvpSource, Table: &env.tableID,
	})
	if err != nil {
		t.Fatalf("ReadChanges: %v", err)
	}
	t.Fatalf("persistierte Changes innerhalb der Zeitspanne fehlgeschlagen (%d von %d)", len(records), limit)
	return nil
}

// imageJSON liest ein Row Image als JSON-Objekt; ein fehlendes Bild ist
// Abwesenheit, kein Fehler (`LH-FA-CAP-008` Boundary).
func imageJSON(t *testing.T, raw []byte) map[string]any {
	t.Helper()
	image := map[string]any{}
	if raw == nil {
		return image
	}
	if err := json.Unmarshal(raw, &image); err != nil {
		t.Fatalf("Row-Image lesen: %v (%s)", err, raw)
	}
	return image
}

// awaitRunEnd wartet auf das reguläre Stream-Ende nach dem Abbruch; die
// Cleanup-Kette droppt den Slot erst nach dem Verbindungs-Ende.
func awaitRunEnd(t *testing.T, cancel context.CancelFunc, runDone <-chan error) {
	t.Helper()
	cancel()
	select {
	case err := <-runDone:
		if err != nil {
			t.Fatalf("Stream-Lauf: %v", err)
		}
	case <-time.After(10 * time.Second):
		t.Fatalf("Stream-Lauf endet nach dem Abbruch nicht")
	}
}

// TestMVPCaptureFlow trägt den MVP-Ablauf: INSERT, UPDATE und DELETE an
// der aktivierten Tabelle werden Ende-zu-Ende erfasst, in Commit-Reihenfolge
// persistiert und mit Inhalt gelesen (`LH-FA-CAP-001`…003,
// `LH-FA-CAP-004`, `LH-QA-POR-003`). Row-Images mit der Default-Replica-
// Identity: das Neu-Bild trägt die Zeile, das Alt-Bild der DELETE-Änderung
// trägt die Schlüsselspalte; der Alt-Stand eines UPDATE bleibt bei
// unverändertem Schlüssel abwesend — Quellverhalten der Default-Identity
// (`LH-FA-CAP-008` Boundary, Bewertung `slice-006` §6 Risiko (b)).
func TestMVPCaptureFlow(t *testing.T) {
	env := newMVPEnv(t, "flow", false)
	cancel, runDone := startCapture(t, env)
	defer cancel()
	ctx := context.Background()

	for _, statement := range []string{
		"INSERT INTO " + env.feed + " (id, name) VALUES (1, 'Alpha')",
		"UPDATE " + env.feed + " SET name = 'Bravo' WHERE id = 1",
		"DELETE FROM " + env.feed + " WHERE id = 1",
	} {
		if _, err := env.pool.Exec(ctx, statement); err != nil {
			t.Fatalf("Quelländerung: %v (%s)", err, statement)
		}
	}

	records := awaitPersistedChanges(t, env, 3)
	if len(records) != 3 {
		t.Fatalf("persistierte Changes: %d (Erwartung: 3)", len(records))
	}

	operations := []model.Operation{
		model.OperationInsert, model.OperationUpdate, model.OperationDelete,
	}
	for i, record := range records {
		if record.Change.Operation != operations[i] {
			t.Fatalf("Change %d: Operation %s, Erwartung %s", i, record.Change.Operation, operations[i])
		}
		if record.Change.SourceTableID != env.tableID {
			t.Fatalf("Change %d: Quelltabelle %q", i, record.Change.SourceTableID)
		}
		if record.Change.Sequence != 1 {
			t.Fatalf("Change %d: Sequenz %d (eine Änderung je Quelltransaktion)", i, record.Change.Sequence)
		}
		if record.Change.TransactionID == "" {
			t.Fatalf("Change %d ohne Transaktionskennung", i)
		}
	}

	// Die Commit-Positionen laufen in Commit-Reihenfolge
	// (`LH-FA-CAP-004`, `LH-FA-REA-004.a`).
	for i := 1; i < len(records); i++ {
		if !records[i-1].Position.Before(records[i].Position) {
			t.Fatalf("Positionsordnung: %v läuft nicht vor %v", records[i-1].Position, records[i].Position)
		}
	}

	insertImage := imageJSON(t, records[0].Change.NewImage)
	if insertImage["id"] != "1" || insertImage["name"] != "Alpha" {
		t.Fatalf("INSERT-Neu-Image: %s", records[0].Change.NewImage)
	}
	if records[0].Change.OldImage != nil {
		t.Fatalf("INSERT-Alt-Image: %s (Abwesenheit erwartet)", records[0].Change.OldImage)
	}
	updateImage := imageJSON(t, records[1].Change.NewImage)
	if updateImage["id"] != "1" || updateImage["name"] != "Bravo" {
		t.Fatalf("UPDATE-Neu-Image: %s", records[1].Change.NewImage)
	}
	if records[1].Change.OldImage != nil {
		t.Fatalf("UPDATE-Alt-Image: %s (Default-Identity ohne Schlüsselwechsel sendet keinen Alt-Stand)", records[1].Change.OldImage)
	}
	deleteImage := imageJSON(t, records[2].Change.OldImage)
	if deleteImage["id"] != "1" {
		t.Fatalf("DELETE-Alt-Image: %s", records[2].Change.OldImage)
	}
	if records[2].Change.NewImage != nil {
		t.Fatalf("DELETE-Neu-Image: %s (Abwesenheit erwartet)", records[2].Change.NewImage)
	}

	// Das Wiederlesen desselben Bereichs trägt dieselben Changes in
	// derselben Reihenfolge (`LH-FA-REA-004.a`, `LH-FA-REA-005`).
	reread, err := env.store.ReadChanges(ctx, outbound.ChangeQuery{
		Source: mvpSource, Table: &env.tableID,
	})
	if err != nil {
		t.Fatalf("Wiederlesen: %v", err)
	}
	if len(reread) != len(records) {
		t.Fatalf("Wiederlesen: %d Changes (Erstsatz: %d)", len(reread), len(records))
	}
	for i := range records {
		if reread[i].Change.ID != records[i].Change.ID ||
			reread[i].Position.Offset != records[i].Position.Offset ||
			string(reread[i].Change.NewImage) != string(records[i].Change.NewImage) {
			t.Fatalf("Wiederlesen: Change %d weicht ab", i)
		}
	}

	// Der Bereich hinter der letzten Position trägt keinen Change
	// (`LH-FA-REA-002` Boundary).
	after, err := model.NewSourcePosition(mvpSource, records[len(records)-1].Position.Offset+0x10000)
	if err != nil {
		t.Fatalf("Bereichs-Position: %v", err)
	}
	tail, err := env.store.ReadChanges(ctx, outbound.ChangeQuery{
		Source: mvpSource, Table: &env.tableID, Start: &after,
	})
	if err != nil {
		t.Fatalf("Bereichslesen: %v", err)
	}
	if len(tail) != 0 {
		t.Fatalf("Bereich hinter der letzten Position: %d Changes (Erwartung: 0)", len(tail))
	}

	awaitRunEnd(t, cancel, runDone)
}

// TestMVPUpdateOldImageWithFullReplicaIdentity trägt die
// REPLICA-IDENTITY-FULL-Seite der Row Images (`LH-FA-CAP-008`): mit
// voller Identity trägt die Quelle den kompletten Alt-Stand — `old_data`
// trägt beide Spalten (`ADR-0016`), nicht nur den Schlüssel.
func TestMVPUpdateOldImageWithFullReplicaIdentity(t *testing.T) {
	env := newMVPEnv(t, "full", true)
	cancel, runDone := startCapture(t, env)
	defer cancel()

	if _, err := env.pool.Exec(context.Background(),
		"INSERT INTO "+env.feed+" (id, name) VALUES (1, 'Alpha')"); err != nil {
		t.Fatalf("INSERT: %v", err)
	}
	if _, err := env.pool.Exec(context.Background(),
		"UPDATE "+env.feed+" SET name = 'Delta' WHERE id = 1"); err != nil {
		t.Fatalf("UPDATE: %v", err)
	}

	records := awaitPersistedChanges(t, env, 2)
	if len(records) != 2 {
		t.Fatalf("persistierte Changes: %d (Erwartung: 2)", len(records))
	}
	if records[0].Change.Operation != model.OperationInsert ||
		records[1].Change.Operation != model.OperationUpdate {
		t.Fatalf("Operationsfolge: %s, %s", records[0].Change.Operation, records[1].Change.Operation)
	}
	oldImage := imageJSON(t, records[1].Change.OldImage)
	if oldImage["id"] != "1" || oldImage["name"] != "Alpha" {
		t.Fatalf("UPDATE-Alt-Image mit voller Identity: %s", records[1].Change.OldImage)
	}
	newImage := imageJSON(t, records[1].Change.NewImage)
	if newImage["id"] != "1" || newImage["name"] != "Delta" {
		t.Fatalf("UPDATE-Neu-Image: %s", records[1].Change.NewImage)
	}

	awaitRunEnd(t, cancel, runDone)
}
