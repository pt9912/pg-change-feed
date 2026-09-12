// Package integration_test trägt den Kern-CDC-Erfassungspfad-Ausschnitt
// des Compose-Integrationstests (`make test-integration`): PostgreSQL
// starten → CDC aktivieren → INSERT → UPDATE → DELETE → Changes lesen →
// Reihenfolge und Inhalt prüfen (`LH-FA-CAP-001`…003, `LH-QA-POR-003`).
// Der über den ursprünglichen MVP-Schnitt (Abschnitt 1 Lastenheft)
// hinausgewachsene Scope desselben Compose-Integrationstests — Rollen-DSN-
// Verifikation (`ADR-0047`, `LH-QA-SEC-001`…`003`), `cdc_capture_lag`-
// Lasttest-Beleg und der Black-Box-CLI-Rundlauf — läuft im Runner-Skript
// (`tools/harness/run-integration-tests.sh`), nicht in diesem Paket. Der
// Lauf hier fährt das verdrahtete System: der Feed-Container trägt die
// CDC-Runtime des Binarys — die Verdrahtung (Store, Aktivierung, Stream,
// Service, ACK) liegt am Composition Root (`ADR-0026`) und läuft im
// Container, nicht hier. Der Test verdrahtet keinen Stream und betreibt
// keinen; seinen Lese-Pfad trägt der Store-Adapter
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
	"github.com/pt9912/pg-change-feed/internal/application/usecase/disable"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/list"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/status"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// mvpSource, mvpPublication und mvpSlot tragen dieselben Werte wie der
// Container-Vertrag in compose.yaml (CDC_SOURCE_ID, CDC_PUBLICATION,
// CDC_SLOT).
const (
	mvpSource      = "src-mvp"
	mvpPublication = "pub_pgc_mvp"
	mvpSlot        = "slot_pgc_mvp"
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
		t.Skip("CDC_INTEGRATION_DSN nicht gesetzt — Compose-Integrationstest läuft über make test-integration gegen die Compose-Umgebung")
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

// changesViewRow trägt eine über `cdc.changes` gelesene Zeile — den
// externen SQL-Lesezugriffsweg (`LH-FA-SST-002`,
// `tools/schema/schema.yaml` `views.changes`), unabhängig vom internen
// `postgresstorage`-Adapter (`SelectChanges` in
// `internal/adapters/driven/postgresstorage/queries/queries.go`).
type changesViewRow struct {
	commitPosition int64
	sequence       int64
	operation      string
	oldData        []byte
	newData        []byte
	schemaVersion  string
}

// queryChangesView liest `cdc.changes`, gefiltert auf Quelle, Tabelle und
// den `id`-Feldwert im Row Image (dieselbe `->>`-Filterform wie im
// Black-Box-CLI-Rundlauf, `tools/harness/run-integration-tests.sh`).
func queryChangesView(ctx context.Context, env *mvpEnv, rowID string) ([]changesViewRow, error) {
	rows, err := env.pool.Query(ctx, `
SELECT commit_position, sequence, operation, old_data, new_data, schema_version
FROM cdc.changes
WHERE source_id = $1 AND source_table_id = $2
  AND (new_data->>'id' = $3 OR old_data->>'id' = $3)
ORDER BY commit_position, sequence`,
		mvpSource, string(env.tableID), rowID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	view := make([]changesViewRow, 0)
	for rows.Next() {
		var row changesViewRow
		if err := rows.Scan(&row.commitPosition, &row.sequence, &row.operation,
			&row.oldData, &row.newData, &row.schemaVersion); err != nil {
			return nil, err
		}
		view = append(view, row)
	}
	return view, rows.Err()
}

// awaitChangesViewRows liest `cdc.changes`, bis die erwartete Anzahl
// Zeilen für `rowID` vorliegt (Polling mit Test-Zeitgrenze) — dieselbe
// Warte-Disziplin wie `awaitPersistedChanges`, gegen den externen
// SQL-Lesezugriffsweg statt gegen den Store-Adapter.
func awaitChangesViewRows(t *testing.T, env *mvpEnv, rowID string, limit int) []changesViewRow {
	t.Helper()
	ctx := context.Background()
	deadline := time.Now().Add(30 * time.Second)
	var view []changesViewRow
	for time.Now().Before(deadline) {
		var err error
		view, err = queryChangesView(ctx, env, rowID)
		if err != nil {
			t.Fatalf("cdc.changes-Lesung: %v", err)
		}
		if len(view) >= limit {
			return view
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatalf("cdc.changes-Lesung innerhalb der Zeitspanne fehlgeschlagen (%d von %d)", len(view), limit)
	return nil
}

// changeRowID liest den `id`-Feldwert eines Changes aus seinem Neu- oder
// Alt-Image; beide Testrichtungen (INSERT/UPDATE tragen `id` im
// Neu-Image, UPDATE/DELETE mit voller Replica-Identity auch im
// Alt-Image) treffen denselben Wert.
func changeRowID(t *testing.T, change model.Change) string {
	t.Helper()
	if id, ok := imageJSON(t, change.NewImage)["id"].(string); ok {
		return id
	}
	if id, ok := imageJSON(t, change.OldImage)["id"].(string); ok {
		return id
	}
	return ""
}

// TestMVPChangesViewMatchesReadChanges belegt `BEO-PGC/lese-doppelquelle`:
// derselbe INSERT/UPDATE/DELETE-Rundlauf, gelesen über die externe
// SQL-Sicht `cdc.changes` (`LH-FA-SST-002`) statt über den internen
// Store-Adapter (`ReadChanges`), trägt dieselbe Reihenfolge
// (`commit_position`, `sequence`) und denselben Feldinhalt (`operation`,
// `old_data`/`new_data`, `schema_version`) wie die `ReadChanges`-Lesung
// desselben Datensatzes (`LH-FA-REA-002`…`006`). Die Zeilen-ID ist
// isoliert von den übrigen Testfällen dieser Datei auf derselben Tabelle
// (id=1 in TestMVPUpdateOldImageWithFullReplicaIdentity, id=90/91 im
// nachgelagerten Lasttest-Beleg und id=95/96 im CLI-E2E-Abschnitt, beide in
// run-integration-tests.sh) — die Lesung filtert auf den Feldwert, nicht auf
// die Testreihenfolge.
func TestMVPChangesViewMatchesReadChanges(t *testing.T) {
	env := newMVPEnv(t, "feed_mvp_full")
	ctx := context.Background()

	const rowID = "60"
	for _, statement := range []string{
		"INSERT INTO " + env.feed + " (id, name) VALUES (60, 'ViewAlpha')",
		"UPDATE " + env.feed + " SET name = 'ViewBravo' WHERE id = 60",
		"DELETE FROM " + env.feed + " WHERE id = 60",
	} {
		if _, err := env.pool.Exec(ctx, statement); err != nil {
			t.Fatalf("Quelländerung: %v (%s)", err, statement)
		}
	}

	viewRows := awaitChangesViewRows(t, env, rowID, 3)
	if len(viewRows) != 3 {
		t.Fatalf("cdc.changes-Lesung: %d Zeilen (Erwartung: 3)", len(viewRows))
	}

	all, err := env.store.ReadChanges(ctx, outbound.ChangeQuery{Source: mvpSource, Table: &env.tableID})
	if err != nil {
		t.Fatalf("ReadChanges: %v", err)
	}
	records := make([]outbound.ChangeRecord, 0, 3)
	for _, record := range all {
		if changeRowID(t, record.Change) == rowID {
			records = append(records, record)
		}
	}
	if len(records) != 3 {
		t.Fatalf("ReadChanges-Lesung: %d Zeilen mit id=%s (Erwartung: 3)", len(records), rowID)
	}

	operations := []model.Operation{model.OperationInsert, model.OperationUpdate, model.OperationDelete}
	for i := range records {
		if records[i].Change.Operation != operations[i] {
			t.Fatalf("ReadChanges Change %d: Operation %s, Erwartung %s", i, records[i].Change.Operation, operations[i])
		}
		if viewRows[i].operation != string(records[i].Change.Operation) {
			t.Fatalf("Change %d: operation Sicht=%s ReadChanges=%s", i, viewRows[i].operation, records[i].Change.Operation)
		}
		if uint64(viewRows[i].commitPosition) != records[i].Position.Offset {
			t.Fatalf("Change %d: commit_position Sicht=%d ReadChanges=%d", i, viewRows[i].commitPosition, records[i].Position.Offset)
		}
		if viewRows[i].sequence != records[i].Change.Sequence {
			t.Fatalf("Change %d: sequence Sicht=%d ReadChanges=%d", i, viewRows[i].sequence, records[i].Change.Sequence)
		}
		if string(viewRows[i].oldData) != string(records[i].Change.OldImage) {
			t.Fatalf("Change %d: old_data Sicht=%s ReadChanges=%s", i, viewRows[i].oldData, records[i].Change.OldImage)
		}
		if string(viewRows[i].newData) != string(records[i].Change.NewImage) {
			t.Fatalf("Change %d: new_data Sicht=%s ReadChanges=%s", i, viewRows[i].newData, records[i].Change.NewImage)
		}
		if viewRows[i].schemaVersion != string(records[i].Change.SchemaVersion) {
			t.Fatalf("Change %d: schema_version Sicht=%s ReadChanges=%s", i, viewRows[i].schemaVersion, records[i].Change.SchemaVersion)
		}
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
		t.Skip("CDC_INTEGRATION_DSN nicht gesetzt — Compose-Integrationstest läuft über make test-integration gegen die Compose-Umgebung")
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
		Source: mvpSource, Schema: "public", Table: "feed_mvp_flow", Publication: mvpPublication,
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
		Source: mvpSource, Schema: "public", Table: "feed_mvp_idle", Publication: mvpPublication,
	})
	if err != nil {
		t.Fatalf("Status der nie aktivierten Tabelle: %v", err)
	}
	if idle.Enabled {
		t.Fatalf("Status von feed_mvp_idle: aktiviert (nie aktivierte Tabelle)")
	}

	// Negative: die fehlende Tabelle endet sichtbar (`LH-FA-CFG-003`).
	if _, err := statusCase.Status(ctx, inbound.GetStatusQuery{
		Source: mvpSource, Schema: "public", Table: "feed_mvp_missing", Publication: mvpPublication,
	}); !errors.Is(err, inbound.ErrSourceTableMissing) {
		t.Fatalf("Status der fehlenden Tabelle: %v (Erwartung: Fehlerklasse der fehlenden Quelle-Tabelle)", err)
	}

	// Liste: die Quelle trägt beide aktivierten Feed-Tabellen
	// (`LH-FA-CFG-004` Happy Path); die nie aktivierte Tabelle bleibt
	// außerhalb.
	tables, err := listCase.ListTables(ctx, inbound.ListTablesQuery{Source: mvpSource, Publication: mvpPublication})
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
	if activated["public.feed_mvp_idle"] || len(tables.Retained) != 0 {
		t.Fatalf("Tabellen-Listen ohne Deaktivierung: aktiviert %v, Herkunft %v", tables.Tables, tables.Retained)
	}
}

// TestMVPDisableRetainedState trägt die Deaktivierung mit Change-Bestand
// am verdrahteten Feed-Container (`LH-FA-CFG-002` Out-of-Scope,
// `ADR-0028`): die Erfassung stoppt über den Publication-Entzug, die
// Bindungs-Zeile bleibt als Herkunft der persistierten Changes — Status
// und Liste trennen den Herkunfts-Bestand vom Erfassungs-Zustand. Der
// Test läuft nach den Capture-Läufen (Quell-Reihenfolge) und deaktiviert
// die Tabelle als letztes.
func TestMVPDisableRetainedState(t *testing.T) {
	env := newMVPEnv(t, "feed_mvp_flow")
	ctx := context.Background()

	activation, err := postgresstorage.NewTableActivation(ctx, env.dsn)
	if err != nil {
		t.Fatalf("Aktivierungs-Adapter: %v", err)
	}
	t.Cleanup(activation.Close)
	disableCase := disable.NewDisableTableService(activation)
	statusCase := status.NewGetStatusService(activation)
	listCase := list.NewListTablesService(activation)

	// Change-Bestand vor der Deaktivierung sichern: eine Zeile läuft über
	// den Feed-Container in den Store, unabhängig von der
	// Test-Reihenfolge der Capture-Läufe. Der Bestand liest vor der
	// Quelländerung, die Erwartung zählt ihn hoch.
	before, err := env.store.ReadChanges(ctx, outbound.ChangeQuery{
		Source: mvpSource, Table: &env.tableID,
	})
	if err != nil {
		t.Fatalf("Change-Bestand lesen: %v", err)
	}
	if _, err := env.pool.Exec(ctx,
		"INSERT INTO "+env.feed+" (id, name) VALUES (7, 'Kilo')"); err != nil {
		t.Fatalf("Quelländerung: %v", err)
	}
	awaitPersistedChanges(t, env, len(before)+1)

	result, err := disableCase.Disable(ctx, inbound.DisableTableCommand{
		Source: mvpSource, Schema: "public", Table: "feed_mvp_flow", Publication: mvpPublication,
	})
	if err != nil {
		t.Fatalf("Disable: %v", err)
	}
	if !result.Retained || result.Removed {
		t.Fatalf("Deaktivierung mit Change-Bestand: %+v (Erwartung: Retained)", result)
	}

	// Der Zustand nach der Deaktivierung: die Bindungs-Zeile liest sich
	// als Herkunft, die Publication trägt die Tabelle nicht mehr.
	state, err := statusCase.Status(ctx, inbound.GetStatusQuery{
		Source: mvpSource, Schema: "public", Table: "feed_mvp_flow", Publication: mvpPublication,
	})
	if err != nil {
		t.Fatalf("Status nach der Deaktivierung: %v", err)
	}
	if state.Enabled || !state.Retained {
		t.Fatalf("Status nach der Deaktivierung: %+v (Erwartung: Retained)", state)
	}
	tables, err := listCase.ListTables(ctx, inbound.ListTablesQuery{Source: mvpSource, Publication: mvpPublication})
	if err != nil {
		t.Fatalf("Tabellen-Liste nach der Deaktivierung: %v", err)
	}
	for _, table := range tables.Tables {
		if table.QualifiedName() == "public.feed_mvp_flow" {
			t.Fatalf("aktivierte Liste trägt die deaktivierte Tabelle")
		}
	}
	retained := map[string]bool{}
	for _, table := range tables.Retained {
		retained[table.QualifiedName()] = true
	}
	if !retained["public.feed_mvp_flow"] {
		t.Fatalf("Herkunfts-Liste ohne die deaktivierte Tabelle: %v", tables.Retained)
	}
}

// TestMVPSchemaChangeAddColumn trägt eine reale `ALTER TABLE … ADD
// COLUMN` auf der aktivierten Tabelle `feed_mvp_schema` am verdrahteten
// Feed-Container (`LH-FA-SCH-001` Happy Path, `LH-FA-SCH-002` Boundary,
// `LH-FA-SCH-005` Boundary): die danach erfassten Changes tragen die neue
// Spalte im Row Image, die zuvor erfasste Change bleibt über
// `cdc.changes` unverändert lesbar und ohne die neue Spalte — und beide
// Changes tragen unterscheidbare Schema-Versionen. Die erste real
// eintreffende Relation-Nachricht trägt die Spaltenform der statisch
// gebundenen Erstversion nach (`mapper.Assembler.Consume`, `ADR-0015`
// Folgepflicht); `ADD COLUMN` ist eine kompatible Erweiterung und
// registriert eine neue `SchemaVersionID` über den `SchemaStorePort`, auf
// die die `TableBinding` gehoben wird — künftige Changes referenzieren
// sie, die bereits erfasste Change bleibt bei ihrer ursprünglichen
// Version.
func TestMVPSchemaChangeAddColumn(t *testing.T) {
	env := newMVPEnv(t, "feed_mvp_schema")
	ctx := context.Background()

	if _, err := env.pool.Exec(ctx,
		"INSERT INTO "+env.feed+" (id, name, amount) VALUES (1, 'Before', '100')"); err != nil {
		t.Fatalf("INSERT vor der Schemaänderung: %v", err)
	}
	beforeRows := awaitChangesViewRows(t, env, "1", 1)

	if _, err := env.pool.Exec(ctx, "ALTER TABLE "+env.feed+" ADD COLUMN extra text"); err != nil {
		t.Fatalf("ALTER TABLE ADD COLUMN: %v", err)
	}
	if _, err := env.pool.Exec(ctx,
		"INSERT INTO "+env.feed+" (id, name, amount, extra) VALUES (2, 'After', '200', 'NewCol')"); err != nil {
		t.Fatalf("INSERT nach der Schemaänderung: %v", err)
	}
	afterRows := awaitChangesViewRows(t, env, "2", 1)

	// LH-FA-SCH-001 Happy Path: die neue Spalte ist im Row Image der
	// danach erfassten Änderung erkennbar.
	afterImage := imageJSON(t, afterRows[0].newData)
	if afterImage["extra"] != "NewCol" {
		t.Fatalf("Row Image nach ADD COLUMN: %s (Erwartung: extra=NewCol)", afterRows[0].newData)
	}

	// LH-FA-SCH-002 Boundary: die ältere Change bleibt unverändert lesbar
	// — sie trägt weiterhin nur die zum Erfassungszeitpunkt bekannten
	// Spalten, ohne die erst danach hinzugefügte.
	rereadBefore := awaitChangesViewRows(t, env, "1", 1)
	if string(rereadBefore[0].newData) != string(beforeRows[0].newData) {
		t.Fatalf("ältere Change nach der Schemaänderung: %s (vor der Änderung: %s)", rereadBefore[0].newData, beforeRows[0].newData)
	}
	beforeImage := imageJSON(t, rereadBefore[0].newData)
	if _, present := beforeImage["extra"]; present {
		t.Fatalf("ältere Change trägt die erst danach hinzugefügte Spalte: %s", rereadBefore[0].newData)
	}

	// LH-FA-SCH-005 Boundary: die Schema-Versionen vor und nach der
	// Erweiterung sind unterscheidbar — `ADD COLUMN` registriert real eine
	// neue Version (`ADR-0015` Folgepflicht).
	if beforeRows[0].schemaVersion == afterRows[0].schemaVersion {
		t.Fatalf("Schema-Version unterscheidet sich nicht: davor %s, danach %s (Erwartung: ADD COLUMN registriert eine neue Version)", beforeRows[0].schemaVersion, afterRows[0].schemaVersion)
	}
}

// TestMVPSchemaChangeIncompatibleTypeChange trägt `LH-FA-SCH-004` an der
// aktivierten Tabelle `feed_mvp_schema` (nach
// TestMVPSchemaChangeAddColumn, mit der dort hinzugefügten Spalte
// `extra`) in beiden im Slice-Plan (`slice-030`, §6) benannten Fällen:
//
//  1. PostgreSQL lehnt eine tatsächlich inkompatible Typänderung bereits
//     selbst über die DDL ab, bevor CDC sie überhaupt sieht — der reale
//     Boundary-Fall aus `LH-FA-SCH-004`.
//  2. Eine Typänderung, die PostgreSQL über eine `USING`-Klausel bei
//     durchgehend konvertierbaren Bestandsdaten zulässt, wird von CDC
//     unverändert übernommen.
//
// Der Negative-Fall aus `LH-FA-SCH-004` — PostgreSQL lässt eine
// Typänderung zu, CDC kann sie aber nicht verlustfrei decodieren, und
// meldet das sichtbar (Fehlerklasse `schema`) — ist mit dem aktuellen
// System nicht real herstellbar: Der Decoder-/Mapper-Pfad
// (`internal/adapters/driving/replication/decode/decode.go` `tupleValues`,
// `internal/adapters/driving/replication/mapper/mapper.go` `rowImage`)
// interpretiert jeden Spaltenwert ausschließlich als Text ohne eigene
// Typprüfung — eine „nicht sicher interpretierbare Schemaänderung“ im
// Sinn von `LH-FA-SCH-004.a` entsteht dabei nicht, weil der Pfad keinen
// Typ interpretiert, den er verlieren könnte. Fall 2 unten belegt das:
// die Typänderung gelingt an PostgreSQL, und CDC übernimmt den
// konvertierten Wert ohne jede sichtbare Fehlermeldung.
func TestMVPSchemaChangeIncompatibleTypeChange(t *testing.T) {
	env := newMVPEnv(t, "feed_mvp_schema")
	ctx := context.Background()

	// Fall 1 (Boundary aus slice-030 §6): PostgreSQL lehnt die DDL ab,
	// weil "NichtNumerisch" nicht in integer konvertiert — CDC sieht die
	// Änderung nie.
	if _, err := env.pool.Exec(ctx,
		"INSERT INTO "+env.feed+" (id, name, amount, extra) VALUES (10, 'Reject', 'NichtNumerisch', 'Z')"); err != nil {
		t.Fatalf("INSERT der nicht konvertierbaren Zeile: %v", err)
	}
	awaitChangesViewRows(t, env, "10", 1)

	if _, err := env.pool.Exec(ctx,
		"ALTER TABLE "+env.feed+" ALTER COLUMN amount TYPE integer USING amount::integer"); err == nil {
		t.Fatalf("ALTER COLUMN TYPE integer mit nicht konvertierbaren Bestandsdaten: kein Fehler (Erwartung: PostgreSQL lehnt ab)")
	} else {
		t.Logf("PostgreSQL lehnt die inkompatible Typänderung selbst ab (slice-030 §6): %v", err)
	}

	if _, err := env.pool.Exec(ctx, "DELETE FROM "+env.feed+" WHERE id = 10"); err != nil {
		t.Fatalf("Aufräumen der nicht konvertierbaren Zeile: %v", err)
	}

	// Fall 2 (Fund, s. Funktionskommentar): dieselbe Typänderung gelingt
	// mit durchgehend konvertierbaren Bestandsdaten (id=1 "100", id=2
	// "200" aus TestMVPSchemaChangeAddColumn) — CDC übernimmt den neuen
	// Wert unverändert als Text, ohne eine Fehlerklasse `schema` zu
	// melden.
	if _, err := env.pool.Exec(ctx,
		"ALTER TABLE "+env.feed+" ALTER COLUMN amount TYPE integer USING amount::integer"); err != nil {
		t.Fatalf("ALTER COLUMN TYPE integer mit konvertierbaren Bestandsdaten: %v", err)
	}
	if _, err := env.pool.Exec(ctx,
		"INSERT INTO "+env.feed+" (id, name, amount, extra) VALUES (11, 'TypeChanged', 300, 'Y')"); err != nil {
		t.Fatalf("INSERT nach der zugelassenen Typänderung: %v", err)
	}
	rows := awaitChangesViewRows(t, env, "11", 1)
	image := imageJSON(t, rows[0].newData)
	if image["amount"] != "300" {
		t.Fatalf("Row Image nach der zugelassenen Typänderung: %s (Erwartung: amount=300 als Text, unverändert übernommen)", rows[0].newData)
	}
}
