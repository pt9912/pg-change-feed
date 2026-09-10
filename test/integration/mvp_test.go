// Package integration_test trägt den MVP-Integrationstest gegen die
// Compose-Umgebung (Abschnitt 1, MVP-Schnitt): PostgreSQL starten → CDC
// aktivieren → INSERT → UPDATE → DELETE → Changes lesen → Reihenfolge und
// Inhalt prüfen (`LH-FA-CAP-001`…003, `LH-QA-POR-003`). Der Lauf fährt das
// verdrahtete System: der Feed-Container trägt die CDC-Runtime des
// Binarys — die Verdrahtung (Store, Aktivierung, Stream, Service, ACK)
// liegt am Composition Root (`ADR-0026`) und läuft im Container, nicht
// hier. Der Test verdrahtet keinen Stream und betreibt keinen; seinen
// Lese-Pfad trägt der Store-Adapter
// (`PostgresChangeStoreAdapter.ReadChanges`) gegen dieselbe Instanz, in
// die das Binary persistiert — persistierte Changes am Ende-zu-Ende-Pfad
// haben keinen anderen Schreiber als den Feed-Container
// (`LH-QA-REL-001.a`).
//
// Die Instanz gehört der Compose-Umgebung (`make test-integration`); der
// d-migrate-Rollout (ADR-0043) trägt die CDC-Tabellen und der Runner die
// Aktivierungs-Vorbedingung (Quell-Tabellen, Quelle-Zeile) vor dem
// Container-Start — die Aktivierung selbst läuft über die Verdrahtung des
// Containers als EnableTable-Aufrufe (`ADR-0028`), und der Test liest den
// Stand über die Status- und Listen-Use-Cases (`LH-FA-CFG-003`/`004`).
// Ohne DSN überspringt der Test.
package integration_test

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage"
	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/list"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/status"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// mvpSource und mvpSlot tragen dieselben Werte wie der Container-Vertrag
// in compose.yaml (CDC_SOURCE_ID, CDC_SLOT).
const (
	mvpSource = "src-mvp"
	mvpSlot   = "slot_pgc_mvp"
)

// mvpEnv trägt die Compose-seitige Testumgebung eines MVP-Laufs: die
// Quell-Verbindung, den Store-Lese-Pfad und die Port-Kennung der
// Feed-Tabelle, die der Test aus den CDC-Referenztabellen liest — der
// Test führt die Bindungs-Kennungen nicht selbst.
type mvpEnv struct {
	dsn     string
	pool    *pgxpool.Pool
	store   *postgresstorage.PostgresChangeStoreAdapter
	feed    string
	tableID model.SourceTableID
}

// newMVPEnv verbindet gegen die Compose-Instanz und liest die
// Port-Kennung der Feed-Tabelle. Der verdrahtete Feed-Container ist
// Vorbedingung: der Slot trägt seinen Lauf — fehlt er, endet der Test mit
// dem Verweis auf den Container-Start (Klasse `configuration`), nicht mit
// einem stillen Warten auf Changes, die niemand streamt. Slot, Publication
// und Feed-Tabellen räumt der Runner mit der Compose-Umgebung ab
// (`compose down -v`).
func newMVPEnv(t *testing.T, feedTable string) *mvpEnv {
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

	var slotCount int
	if err := pool.QueryRow(ctx,
		"SELECT count(*) FROM pg_replication_slots WHERE slot_name = $1", mvpSlot,
	).Scan(&slotCount); err != nil {
		t.Fatalf("Slot-Prüfung: %v", err)
	}
	if slotCount != 1 {
		t.Fatalf("Replication-Slot %q fehlt — der Feed-Container trägt die CDC-Verdrahtung nicht (Container-Start im Runner, compose.yaml)", mvpSlot)
	}

	var tableID string
	if err := pool.QueryRow(ctx,
		"SELECT source_table_id FROM cdc.source_table WHERE schema_name = 'public' AND table_name = $1", feedTable,
	).Scan(&tableID); err != nil {
		t.Fatalf("Bindungs-Zeile für %q: %v — die Aktivierung trägt die Verdrahtung des Feed-Containers als EnableTable-Aufruf (ADR-0028)", feedTable, err)
	}

	store, err := postgresstorage.New(ctx, dsn)
	if err != nil {
		t.Fatalf("Store: %v", err)
	}
	t.Cleanup(store.Close)
	return &mvpEnv{
		dsn:     dsn,
		pool:    pool,
		store:   store,
		feed:    "public." + feedTable,
		tableID: model.SourceTableID(tableID),
	}
}

// awaitPersistedChanges liest den Store, bis die erwarteten Changes der
// Tabelle dieser Umgebungen persistiert sind (Polling mit Test-Zeitgrenze);
// der Erfassungsweg läuft im Feed-Container. Der Tabellenfilter trennt die
// Läufe der beiden Feed-Tabellen.
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

// TestMVPCaptureFlow trägt den MVP-Ablauf am verdrahteten Feed-Container:
// INSERT, UPDATE und DELETE an der aktivierten Tabelle werden Ende-zu-Ende
// durch das Binary erfasst, in Commit-Reihenfolge persistiert und mit
// Inhalt gelesen (`LH-FA-CAP-001`…003, `LH-FA-CAP-004`, `LH-QA-POR-003`).
// Row-Images mit der Default-Replica-Identity: das Neu-Bild trägt die
// Zeile, das Alt-Bild der DELETE-Änderung trägt die Schlüsselspalte; der
// Alt-Stand eines UPDATE bleibt bei unverändertem Schlüssel abwesend —
// Quellverhalten der Default-Identity (`LH-FA-CAP-008` Boundary).
func TestMVPCaptureFlow(t *testing.T) {
	env := newMVPEnv(t, "feed_mvp_flow")
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
}

// TestMVPUpdateOldImageWithFullReplicaIdentity trägt die
// REPLICA-IDENTITY-FULL-Seite der Row Images am verdrahteten
// Feed-Container (`LH-FA-CAP-008`): mit voller Identity trägt die Quelle
// den kompletten Alt-Stand — `old_data` trägt beide Spalten (`ADR-0016`),
// nicht nur den Schlüssel.
func TestMVPUpdateOldImageWithFullReplicaIdentity(t *testing.T) {
	env := newMVPEnv(t, "feed_mvp_full")
	ctx := context.Background()

	if _, err := env.pool.Exec(ctx,
		"INSERT INTO "+env.feed+" (id, name) VALUES (1, 'Alpha')"); err != nil {
		t.Fatalf("INSERT: %v", err)
	}
	if _, err := env.pool.Exec(ctx,
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
}

// TestMVPActivationState liest den Aktivierungsstand am verdrahteten
// Feed-Container über die Status- und Listen-Use-Cases (`LH-FA-CFG-003`,
// `LH-FA-CFG-004`, `ADR-0028`): die aktivierte Feed-Tabelle meldet
// „aktiviert", die nie aktivierte Tabelle meldet „nicht aktiviert"
// (Boundary), die fehlende Tabelle endet sichtbar (Negative), und die
// Liste trägt die aktivierten Tabellen der Quelle.
func TestMVPActivationState(t *testing.T) {
	dsn := os.Getenv("CDC_INTEGRATION_DSN")
	if dsn == "" {
		t.Skip("CDC_INTEGRATION_DSN nicht gesetzt — MVP-Integrationstest läuft über make test-integration gegen die Compose-Umgebung")
	}
	ctx := context.Background()
	activation, err := postgresstorage.NewTableActivation(ctx, dsn)
	if err != nil {
		t.Fatalf("Aktivierungs-Adapter: %v", err)
	}
	t.Cleanup(activation.Close)
	statusCase := status.NewGetStatusService(activation)
	listCase := list.NewListTablesService(activation)

	// Happy Path: die aktivierte Tabelle meldet „aktiviert" (`LH-FA-CFG-003`).
	enabled, err := statusCase.Status(ctx, inbound.GetStatusQuery{
		Source: mvpSource, Schema: "public", Table: "feed_mvp_flow",
	})
	if err != nil {
		t.Fatalf("Status der aktivierten Tabelle: %v", err)
	}
	if !enabled.Enabled {
		t.Fatalf("Status von feed_mvp_flow: nicht aktiviert — die Verdrahtung des Feed-Containers trägt die Aktivierung als EnableTable-Aufruf (ADR-0028)")
	}

	// Boundary: die nie aktivierte Tabelle meldet „nicht aktiviert"
	// (`LH-FA-CFG-003`).
	idle, err := statusCase.Status(ctx, inbound.GetStatusQuery{
		Source: mvpSource, Schema: "public", Table: "feed_mvp_idle",
	})
	if err != nil {
		t.Fatalf("Status der nie aktivierten Tabelle: %v", err)
	}
	if idle.Enabled {
		t.Fatalf("Status von feed_mvp_idle: aktiviert (nie aktivierte Tabelle)")
	}

	// Negative: die fehlende Tabelle endet sichtbar (`LH-FA-CFG-003`).
	if _, err := statusCase.Status(ctx, inbound.GetStatusQuery{
		Source: mvpSource, Schema: "public", Table: "feed_mvp_missing",
	}); !errors.Is(err, inbound.ErrSourceTableMissing) {
		t.Fatalf("Status der fehlenden Tabelle: %v (Erwartung: Fehlerklasse der fehlenden Quelle-Tabelle)", err)
	}

	// Liste: die Quelle trägt beide aktivierten Feed-Tabellen
	// (`LH-FA-CFG-004` Happy Path); die nie aktivierte Tabelle bleibt
	// außerhalb.
	tables, err := listCase.ListTables(ctx, inbound.ListTablesQuery{Source: mvpSource})
	if err != nil {
		t.Fatalf("Tabellen-Liste: %v", err)
	}
	activated := map[string]bool{}
	for _, table := range tables.Tables {
		activated[table.QualifiedName()] = true
	}
	for _, expected := range []string{"public.feed_mvp_flow", "public.feed_mvp_full"} {
		if !activated[expected] {
			t.Fatalf("Tabellen-Liste ohne %q: %d Tabellen", expected, len(tables.Tables))
		}
	}
	if activated["public.feed_mvp_idle"] {
		t.Fatalf("Tabellen-Liste trägt die nie aktivierte Tabelle")
	}
}
