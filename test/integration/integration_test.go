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
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
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

// e2eSource, e2ePublication und e2eSlot tragen dieselben Werte wie der
// Container-Vertrag in compose.yaml (CDC_SOURCE_ID, CDC_PUBLICATION,
// CDC_SLOT).
const (
	e2eSource      = "src-e2e"
	e2ePublication = "pub_pgc_e2e"
	e2eSlot        = "slot_pgc_e2e"
)

// e2eEnv trägt die Compose-seitige Testumgebung eines E2E-Laufs: die
// Quell-Verbindung, den Store-Lese-Pfad, den Klartext-Tabellennamen der
// Feed-Tabelle und ihre Port-Kennung, die der Test aus den
// CDC-Referenztabellen liest — der Test führt die Bindungs-Kennungen nicht
// selbst. Der Tabellenfilter des Leseports läuft über den Klartext-Namen
// (`ADR-0081` Teilfrage 3).
type e2eEnv struct {
	dsn     string
	pool    *pgxpool.Pool
	store   *postgresstorage.PostgresChangeStoreAdapter
	feed    string
	table   string
	tableID model.SourceTableID
}

// newE2EEnv verbindet gegen die Compose-Instanz und liest die
// Port-Kennung der Feed-Tabelle. Der verdrahtete Feed-Container ist
// Vorbedingung: der Slot trägt seinen Lauf — fehlt er, endet der Test mit
// dem Verweis auf den Container-Start (Klasse `configuration`), nicht mit
// einem stillen Warten auf Changes, die niemand streamt. Slot, Publication
// und Feed-Tabellen räumt der Runner mit der Compose-Umgebung ab
// (`compose down -v`).
func newE2EEnv(t *testing.T, feedTable string) *e2eEnv {
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
		"SELECT count(*) FROM pg_replication_slots WHERE slot_name = $1", e2eSlot,
	).Scan(&slotCount); err != nil {
		t.Fatalf("Slot-Prüfung: %v", err)
	}
	if slotCount != 1 {
		t.Fatalf("Replication-Slot %q fehlt — der Feed-Container trägt die CDC-Verdrahtung nicht (Container-Start im Runner, compose.yaml)", e2eSlot)
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
	return &e2eEnv{
		dsn:     dsn,
		pool:    pool,
		store:   store,
		feed:    "public." + feedTable,
		table:   feedTable,
		tableID: model.SourceTableID(tableID),
	}
}

// awaitPersistedChanges liest den Store, bis die erwarteten Changes der
// Tabelle dieser Umgebungen persistiert sind (Polling mit Test-Zeitgrenze);
// der Erfassungsweg läuft im Feed-Container. Der Tabellenfilter trennt die
// Läufe der beiden Feed-Tabellen.
func awaitPersistedChanges(t *testing.T, env *e2eEnv, limit int) []outbound.ChangeRecord {
	t.Helper()
	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		records, err := env.store.ReadChanges(context.Background(), outbound.ChangeQuery{
			Source: e2eSource, Table: env.table,
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
		Source: e2eSource, Table: env.table,
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

// TestE2ECaptureFlow trägt den E2E-Ablauf am verdrahteten Feed-Container:
// INSERT, UPDATE und DELETE an der aktivierten Tabelle werden Ende-zu-Ende
// durch das Binary erfasst, in Commit-Reihenfolge persistiert und mit
// Inhalt gelesen (`LH-FA-CAP-001`…003, `LH-FA-CAP-004`, `LH-QA-POR-003`).
// Trägt zugleich `LH-FA-SST-001`: die definierte Schnittstelle zur
// PostgreSQL-Quelle (Logical Replication) läuft hier real Ende-zu-Ende,
// nicht nur strukturell.
// Row-Images mit der Default-Replica-Identity: das Neu-Bild trägt die
// Zeile, das Alt-Bild der DELETE-Änderung trägt die Schlüsselspalte; der
// Alt-Stand eines UPDATE bleibt bei unverändertem Schlüssel abwesend —
// Quellverhalten der Default-Identity (`LH-FA-CAP-008` Boundary).
// Trägt zugleich `LH-FA-DAT-002` (Quelltabelle über `SourceTableID`
// identifizierbar), `LH-FA-DAT-003` (INSERT/UPDATE/DELETE je Change
// unterscheidbar) und `LH-FA-DAT-005` (Row-Image-Werte je Operationstyp
// real gelesen). Das Wiederlesen desselben Bereichs unten trägt zugleich
// `LH-QA-REL-004`: zweimaliges Lesen liefert dieselben Changes.
func TestE2ECaptureFlow(t *testing.T) {
	env := newE2EEnv(t, "feed_e2e_flow")
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
		Source: e2eSource, Table: env.table,
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
	after, err := model.NewSourcePosition(e2eSource, records[len(records)-1].Position.Offset+0x10000)
	if err != nil {
		t.Fatalf("Bereichs-Position: %v", err)
	}
	tail, err := env.store.ReadChanges(ctx, outbound.ChangeQuery{
		Source: e2eSource, Table: env.table, Start: &after,
	})
	if err != nil {
		t.Fatalf("Bereichslesen: %v", err)
	}
	if len(tail) != 0 {
		t.Fatalf("Bereich hinter der letzten Position: %d Changes (Erwartung: 0)", len(tail))
	}
}

// TestE2EUpdateOldImageWithFullReplicaIdentity trägt die
// REPLICA-IDENTITY-FULL-Seite der Row Images am verdrahteten
// Feed-Container (`LH-FA-CAP-008`): mit voller Identity trägt die Quelle
// den kompletten Alt-Stand — `old_data` trägt beide Spalten (`ADR-0016`),
// nicht nur den Schlüssel.
func TestE2EUpdateOldImageWithFullReplicaIdentity(t *testing.T) {
	env := newE2EEnv(t, "feed_e2e_full")
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
func queryChangesView(ctx context.Context, env *e2eEnv, rowID string) ([]changesViewRow, error) {
	rows, err := env.pool.Query(ctx, `
SELECT commit_position, sequence, operation, old_data, new_data, schema_version
FROM cdc.changes
WHERE source_id = $1 AND source_table_id = $2
  AND (new_data->>'id' = $3 OR old_data->>'id' = $3)
ORDER BY commit_position, sequence`,
		e2eSource, string(env.tableID), rowID)
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
func awaitChangesViewRows(t *testing.T, env *e2eEnv, rowID string, limit int) []changesViewRow {
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

// TestE2EChangesViewMatchesReadChanges belegt `BEO-PGC/lese-doppelquelle`:
// derselbe INSERT/UPDATE/DELETE-Rundlauf, gelesen über die externe
// SQL-Sicht `cdc.changes` (`LH-FA-SST-002`) statt über den internen
// Store-Adapter (`ReadChanges`), trägt dieselbe Reihenfolge
// (`commit_position`, `sequence`) und denselben Feldinhalt (`operation`,
// `old_data`/`new_data`, `schema_version`) wie die `ReadChanges`-Lesung
// desselben Datensatzes (`LH-FA-REA-002`…`006`). Die Zeilen-ID ist
// isoliert von den übrigen Testfällen dieser Datei auf derselben Tabelle
// (id=1 in TestE2EUpdateOldImageWithFullReplicaIdentity, id=90/91 im
// nachgelagerten Lasttest-Beleg und id=95/96 im CLI-E2E-Abschnitt, beide in
// run-integration-tests.sh) — die Lesung filtert auf den Feldwert, nicht auf
// die Testreihenfolge.
func TestE2EChangesViewMatchesReadChanges(t *testing.T) {
	env := newE2EEnv(t, "feed_e2e_full")
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

	all, err := env.store.ReadChanges(ctx, outbound.ChangeQuery{Source: e2eSource, Table: env.table})
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

// TestE2ERetentionBlockersViewShowsFurthestBehindConsumer belegt
// `LH-FA-RET-005`: `cdc.retention_blockers` zeigt je Quelle real den
// Consumer mit der am weitesten zurückliegenden bestätigten Position — der
// Consumer, dessen Position `RetentionPolicy.AllowsDeletion`
// (`RunRetentionUseCase`) aktuell als Löschgrenze der Quelle behandelt.
// Registrierung und Bestätigung laufen direkt über den
// `ConsumerStatePort`-Adapter, ohne CLI-Rundlauf — isoliert von den über
// `register-consumer`/`acknowledge-consumer` geführten Consumern des
// externen Black-Box-Rundlaufs (`tools/harness/run-integration-tests.sh`),
// die erst nach diesem Go-Testlauf entstehen. Real gegen zwei Consumer
// getestet (`LH-FA-RET-005` Happy Path/Boundary; trägt zugleich einen
// Teilbeleg für `LH-FA-CON-002`: `behindConsumer`s gespeicherte Position
// bleibt exakt sein eigener bestätigter Wert, unverändert durch
// `aheadConsumer`s spätere, unabhängige Bestätigung — die Happy-Path-/
// Boundary-Formulierung aus `spec/lastenheft.md` selbst, denselben
// Bereich gleichzeitig lesend, prüft dieser Testfall nicht):
// einer blockiert (weiter
// zurückliegende Position), einer nicht (bereits weiter bestätigt). Die
// Sicht berechnet nichts neu, was die Domain-Policy nicht bereits real
// entscheidet — sie macht nur sichtbar, welche Position aktuell die
// Löschgrenze trägt.
func TestE2ERetentionBlockersViewShowsFurthestBehindConsumer(t *testing.T) {
	env := newE2EEnv(t, "feed_e2e_full")
	ctx := context.Background()

	for _, statement := range []string{
		"INSERT INTO " + env.feed + " (id, name) VALUES (70, 'RetentionViewBehind')",
		"INSERT INTO " + env.feed + " (id, name) VALUES (71, 'RetentionViewAhead')",
	} {
		if _, err := env.pool.Exec(ctx, statement); err != nil {
			t.Fatalf("Quelländerung: %v (%s)", err, statement)
		}
	}
	behindRows := awaitChangesViewRows(t, env, "70", 1)
	aheadRows := awaitChangesViewRows(t, env, "71", 1)
	behindPosition, err := model.NewSourcePosition(e2eSource, uint64(behindRows[0].commitPosition))
	if err != nil {
		t.Fatalf("Positions-Konstruktion (behind): %v", err)
	}
	aheadPosition, err := model.NewSourcePosition(e2eSource, uint64(aheadRows[0].commitPosition))
	if err != nil {
		t.Fatalf("Positions-Konstruktion (ahead): %v", err)
	}

	state, err := postgresstorage.NewConsumerState(ctx, env.dsn)
	if err != nil {
		t.Fatalf("ConsumerState-Adapter: %v", err)
	}
	t.Cleanup(state.Close)

	const behindConsumer = model.ConsumerID("retention-view-behind")
	const aheadConsumer = model.ConsumerID("retention-view-ahead")
	for _, consumer := range []model.Consumer{
		{ID: behindConsumer, Name: "Retention View Behind"},
		{ID: aheadConsumer, Name: "Retention View Ahead"},
	} {
		if _, err := state.Register(ctx, consumer); err != nil {
			t.Fatalf("Register(%s): %v", consumer.ID, err)
		}
	}
	// Beide Consumer entfernen, sobald dieser Testfall endet: eine
	// bestätigte, nie wieder fortgeschriebene Position dieser Consumer
	// würde sonst über das Testende hinaus als reale, dauerhaft
	// zurückliegende Position in `Positions(e2eSource)` weiterleben und
	// den späteren Retention-Beleg des Compose-Laufs
	// (`tools/harness/run-integration-tests.sh`, `RunRetentionUseCase`)
	// dauerhaft blockieren — die Zeilen-Abwesenheit ist hier die
	// Rücknahme, dieselbe Lesart wie bei jeder administrativen Entfernung
	// (`ConsumerStatePort.Remove`-Doku).
	t.Cleanup(func() {
		if _, err := state.Remove(context.Background(), behindConsumer); err != nil {
			t.Errorf("Remove(%s) nach Testende: %v", behindConsumer, err)
		}
		if _, err := state.Remove(context.Background(), aheadConsumer); err != nil {
			t.Errorf("Remove(%s) nach Testende: %v", aheadConsumer, err)
		}
	})
	if _, err := state.Acknowledge(ctx, model.ConsumerPosition{ConsumerID: behindConsumer, Position: behindPosition}); err != nil {
		t.Fatalf("Acknowledge(%s): %v", behindConsumer, err)
	}
	if _, err := state.Acknowledge(ctx, model.ConsumerPosition{ConsumerID: aheadConsumer, Position: aheadPosition}); err != nil {
		t.Fatalf("Acknowledge(%s): %v", aheadConsumer, err)
	}

	rows, err := env.pool.Query(ctx,
		"SELECT consumer_id, acknowledged_position, backlog FROM cdc.retention_blockers WHERE source_id = $1 AND consumer_id = ANY($2)",
		e2eSource, []string{string(behindConsumer), string(aheadConsumer)})
	if err != nil {
		t.Fatalf("cdc.retention_blockers-Lesung: %v", err)
	}
	defer rows.Close()
	type blockerRow struct {
		consumerID string
		position   int64
		backlog    int64
	}
	var blockers []blockerRow
	for rows.Next() {
		var row blockerRow
		if err := rows.Scan(&row.consumerID, &row.position, &row.backlog); err != nil {
			t.Fatalf("cdc.retention_blockers-Scan: %v", err)
		}
		blockers = append(blockers, row)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("cdc.retention_blockers-Lesung: %v", err)
	}

	// Happy Path (LH-FA-RET-005): der weiter zurückliegende Consumer
	// erscheint als aktueller Blocker der Quelle, mit realem Rückstand.
	if len(blockers) != 1 {
		t.Fatalf("cdc.retention_blockers für %q/%q: %d Zeilen (Erwartung: 1 — genau der aktuell blockierende Consumer je Quelle)", behindConsumer, aheadConsumer, len(blockers))
	}
	if blockers[0].consumerID != string(behindConsumer) {
		t.Fatalf("cdc.retention_blockers zeigt %q als Blocker (Erwartung: %q, der weiter zurückliegende Consumer)", blockers[0].consumerID, behindConsumer)
	}
	if blockers[0].position != behindRows[0].commitPosition {
		t.Fatalf("cdc.retention_blockers.acknowledged_position: %d (Erwartung: %d)", blockers[0].position, behindRows[0].commitPosition)
	}
	if blockers[0].backlog <= 0 {
		t.Fatalf("cdc.retention_blockers.backlog: %d (Erwartung: > 0 — der zurückliegende Consumer hat einen realen Rückstand)", blockers[0].backlog)
	}

	// Boundary (LH-FA-RET-005): der bereits weiter bestätigende Consumer
	// blockiert nicht — nur der am weitesten zurückliegende Consumer wird
	// je Quelle gezeigt, nicht auch der bereits vorbeigezogene.
	if blockers[0].consumerID == string(aheadConsumer) {
		t.Fatalf("cdc.retention_blockers zeigt den bereits weiter bestätigenden Consumer %q als Blocker", aheadConsumer)
	}
}

// TestE2EMetricsCarriesStorageBytes belegt `LH-FA-RET-006` (`SPEC-009`
// `cdc_storage_bytes`) am verdrahteten Feed-Container: `cdc.metrics` trägt
// nach einer realen CDC-Erfassung einen positiven Wert für die physische
// Speichergröße von `cdc.change` — derselbe externe SQL-Lesezugriffsweg wie
// der bestehende `cdc_capture_lag`-Beleg
// (`tools/harness/run-integration-tests.sh`), hier gegen dieselbe View.
func TestE2EMetricsCarriesStorageBytes(t *testing.T) {
	env := newE2EEnv(t, "feed_e2e_full")
	ctx := context.Background()

	if _, err := env.pool.Exec(ctx,
		"INSERT INTO "+env.feed+" (id, name) VALUES (80, 'StorageBytesMetric')",
	); err != nil {
		t.Fatalf("Quelländerung: %v", err)
	}
	awaitChangesViewRows(t, env, "80", 1)

	var storageBytes float64
	if err := env.pool.QueryRow(ctx,
		"SELECT value FROM cdc.metrics WHERE metric_name = 'cdc_storage_bytes'",
	).Scan(&storageBytes); err != nil {
		t.Fatalf("cdc_storage_bytes-Lesen: %v", err)
	}
	if storageBytes <= 0 {
		t.Fatalf("cdc_storage_bytes = %v, wollen > 0 nach realer CDC-Erfassung", storageBytes)
	}
}

// TestE2EActivationState liest den Aktivierungsstand am verdrahteten
// Feed-Container über die Status- und Listen-Use-Cases (`LH-FA-CFG-003`,
// `LH-FA-CFG-004`, `ADR-0028`): die aktivierte Feed-Tabelle meldet
// „aktiviert", die nie aktivierte Tabelle meldet „nicht aktiviert"
// (Boundary), die fehlende Tabelle endet sichtbar (Negative), und die
// Liste trägt die aktivierten Tabellen der Quelle.
func TestE2EActivationState(t *testing.T) {
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
		Source: e2eSource, Schema: "public", Table: "feed_e2e_flow", Publication: e2ePublication,
	})
	if err != nil {
		t.Fatalf("Status der aktivierten Tabelle: %v", err)
	}
	if !enabled.Enabled {
		t.Fatalf("Status von feed_e2e_flow: nicht aktiviert — die Verdrahtung des Feed-Containers trägt die Aktivierung als EnableTable-Aufruf (ADR-0028)")
	}

	// Boundary: die nie aktivierte Tabelle meldet „nicht aktiviert"
	// (`LH-FA-CFG-003`).
	idle, err := statusCase.Status(ctx, inbound.GetStatusQuery{
		Source: e2eSource, Schema: "public", Table: "feed_e2e_idle", Publication: e2ePublication,
	})
	if err != nil {
		t.Fatalf("Status der nie aktivierten Tabelle: %v", err)
	}
	if idle.Enabled {
		t.Fatalf("Status von feed_e2e_idle: aktiviert (nie aktivierte Tabelle)")
	}

	// Negative: die fehlende Tabelle endet sichtbar (`LH-FA-CFG-003`).
	if _, err := statusCase.Status(ctx, inbound.GetStatusQuery{
		Source: e2eSource, Schema: "public", Table: "feed_e2e_missing", Publication: e2ePublication,
	}); !errors.Is(err, inbound.ErrSourceTableMissing) {
		t.Fatalf("Status der fehlenden Tabelle: %v (Erwartung: Fehlerklasse der fehlenden Quelle-Tabelle)", err)
	}

	// Liste: die Quelle trägt beide aktivierten Feed-Tabellen
	// (`LH-FA-CFG-004` Happy Path); die nie aktivierte Tabelle bleibt
	// außerhalb.
	tables, err := listCase.ListTables(ctx, inbound.ListTablesQuery{Source: e2eSource, Publication: e2ePublication})
	if err != nil {
		t.Fatalf("Tabellen-Liste: %v", err)
	}
	activated := map[string]bool{}
	for _, table := range tables.Tables {
		activated[table.QualifiedName()] = true
	}
	for _, expected := range []string{"public.feed_e2e_flow", "public.feed_e2e_full"} {
		if !activated[expected] {
			t.Fatalf("Tabellen-Liste ohne %q: %d Tabellen", expected, len(tables.Tables))
		}
	}
	if activated["public.feed_e2e_idle"] || len(tables.Retained) != 0 {
		t.Fatalf("Tabellen-Listen ohne Deaktivierung: aktiviert %v, Herkunft %v", tables.Tables, tables.Retained)
	}
}

// TestE2EActiveTablesViewMatchesActivationState liest den
// Aktivierungsstand am verdrahteten Feed-Container über die externe
// SQL-Sicht `cdc.active_tables` (`LH-FA-SST-002`) statt über den
// internen Status-Use-Case: die aktivierte Feed-Tabelle erscheint in der
// Sicht (`LH-FA-CFG-003`/`004` Happy Path), die nie aktivierte Tabelle
// nicht (`LH-FA-CFG-003`/`004` Boundary) — beide Lesewege tragen dieselbe
// Aussage über denselben Bindungszustand. `feed_e2e_full` bleibt über den
// gesamten Compose-Lauf aktiviert (Lasttest-Beleg in
// run-integration-tests.sh) und ist damit unabhängig von der
// Deaktivierung in TestE2EDisableRetainedState, die nur feed_e2e_flow
// betrifft.
func TestE2EActiveTablesViewMatchesActivationState(t *testing.T) {
	env := newE2EEnv(t, "feed_e2e_full")
	ctx := context.Background()

	activation, err := postgresstorage.NewTableActivation(ctx, env.dsn)
	if err != nil {
		t.Fatalf("Aktivierungs-Adapter: %v", err)
	}
	t.Cleanup(activation.Close)
	statusCase := status.NewGetStatusService(activation)

	// Happy Path: die aktivierte Tabelle erscheint in `cdc.active_tables`,
	// gefiltert auf ihre eigene Bindungs-Kennung — isoliert von den
	// übrigen im selben Compose-Lauf aktivierten Tabellen
	// (`feed_e2e_flow`, `feed_e2e_schema`).
	var activeCount int
	if err := env.pool.QueryRow(ctx,
		"SELECT count(*) FROM cdc.active_tables WHERE source_table_id = $1",
		string(env.tableID),
	).Scan(&activeCount); err != nil {
		t.Fatalf("cdc.active_tables-Lesung (feed_e2e_full): %v", err)
	}
	if activeCount != 1 {
		t.Fatalf("cdc.active_tables ohne feed_e2e_full: %d Zeilen (Erwartung: 1)", activeCount)
	}
	enabled, err := statusCase.Status(ctx, inbound.GetStatusQuery{
		Source: e2eSource, Schema: "public", Table: "feed_e2e_full", Publication: e2ePublication,
	})
	if err != nil {
		t.Fatalf("Status der aktivierten Tabelle: %v", err)
	}
	if !enabled.Enabled {
		t.Fatalf("Status von feed_e2e_full: nicht aktiviert, obwohl cdc.active_tables eine Zeile trägt")
	}

	// Boundary: die nie aktivierte Tabelle (kein `source_table`-Eintrag)
	// erscheint nicht in `cdc.active_tables`.
	var idleCount int
	if err := env.pool.QueryRow(ctx,
		"SELECT count(*) FROM cdc.active_tables WHERE source_id = $1 AND schema_name = 'public' AND table_name = $2",
		e2eSource, "feed_e2e_idle",
	).Scan(&idleCount); err != nil {
		t.Fatalf("cdc.active_tables-Lesung (feed_e2e_idle): %v", err)
	}
	if idleCount != 0 {
		t.Fatalf("cdc.active_tables mit feed_e2e_idle: %d Zeilen (Erwartung: 0)", idleCount)
	}
	idle, err := statusCase.Status(ctx, inbound.GetStatusQuery{
		Source: e2eSource, Schema: "public", Table: "feed_e2e_idle", Publication: e2ePublication,
	})
	if err != nil {
		t.Fatalf("Status der nie aktivierten Tabelle: %v", err)
	}
	if idle.Enabled {
		t.Fatalf("Status von feed_e2e_idle: aktiviert, obwohl cdc.active_tables keine Zeile trägt")
	}
}

// TestE2EDisableRetainedState trägt die Deaktivierung mit Change-Bestand
// am verdrahteten Feed-Container (`LH-FA-CFG-002` Out-of-Scope,
// `ADR-0028`): die Erfassung stoppt über den Publication-Entzug, die
// Bindungs-Zeile bleibt als Herkunft der persistierten Changes — Status
// und Liste trennen den Herkunfts-Bestand vom Erfassungs-Zustand. Der
// Test läuft nach den Capture-Läufen (Quell-Reihenfolge) und deaktiviert
// die Tabelle als letztes.
func TestE2EDisableRetainedState(t *testing.T) {
	env := newE2EEnv(t, "feed_e2e_flow")
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
		Source: e2eSource, Table: env.table,
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
		Source: e2eSource, Schema: "public", Table: "feed_e2e_flow", Publication: e2ePublication,
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
		Source: e2eSource, Schema: "public", Table: "feed_e2e_flow", Publication: e2ePublication,
	})
	if err != nil {
		t.Fatalf("Status nach der Deaktivierung: %v", err)
	}
	if state.Enabled || !state.Retained {
		t.Fatalf("Status nach der Deaktivierung: %+v (Erwartung: Retained)", state)
	}
	tables, err := listCase.ListTables(ctx, inbound.ListTablesQuery{Source: e2eSource, Publication: e2ePublication})
	if err != nil {
		t.Fatalf("Tabellen-Liste nach der Deaktivierung: %v", err)
	}
	for _, table := range tables.Tables {
		if table.QualifiedName() == "public.feed_e2e_flow" {
			t.Fatalf("aktivierte Liste trägt die deaktivierte Tabelle")
		}
	}
	retained := map[string]bool{}
	for _, table := range tables.Retained {
		retained[table.QualifiedName()] = true
	}
	if !retained["public.feed_e2e_flow"] {
		t.Fatalf("Herkunfts-Liste ohne die deaktivierte Tabelle: %v", tables.Retained)
	}
}

// TestE2ESchemaChangeAddColumn trägt eine reale `ALTER TABLE … ADD
// COLUMN` auf der aktivierten Tabelle `feed_e2e_schema` am verdrahteten
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
func TestE2ESchemaChangeAddColumn(t *testing.T) {
	env := newE2EEnv(t, "feed_e2e_schema")
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

// TestE2ESchemaChangeDropColumn trägt eine reale `ALTER TABLE … DROP
// COLUMN` auf der aktivierten Tabelle `feed_e2e_schema` am verdrahteten
// Feed-Container (`LH-FA-SCH-003` Happy Path und Boundary, `ADR-0063`
// Supersedes `ADR-0058` Entscheidung 1): eine real entfernte Spalte fehlt
// in der eingehenden Relation-Nachricht genauso wie jede andere gelöschte
// Spalte und löst denselben `relationOther`-Pfad mit
// `mapper.ErrIncompatibleSchemaChange` aus wie eine inkompatible
// Typänderung (`LH-FA-SCH-004`, bereits von `ADR-0059` Teilfrage 4
// akzeptierter Präzedenzfall) — sichtbarer `schema`-Fehler, kein stilles
// Auslassen. Der Erfassungspfad des Feed-Containers endet darüber
// dauerhaft (`restart: "no"` in `compose.yaml`), genau wie bei
// `TestE2ESchemaChangeIncompatibleTypeChange` — deshalb läuft diese
// Funktion nach der bisherigen Container-Ende-Grenze, mit einem expliziten
// Container-Neustart davor (`tools/harness/run-integration-tests.sh`). Die
// Boundary-Klausel bleibt unverändert: die vor der Entfernung erfasste
// Change bleibt über `cdc.changes` inklusive historischem Wert lesbar. Die
// Spalte `removable` ist eine eigene, wegwerfbare Spalte, getrennt von
// `amount`/`extra`, die bereits Zustand für
// `TestE2ESchemaChangeIncompatibleTypeChange` tragen.
func TestE2ESchemaChangeDropColumn(t *testing.T) {
	env := newE2EEnv(t, "feed_e2e_schema")
	ctx := context.Background()

	if _, err := env.pool.Exec(ctx, "ALTER TABLE "+env.feed+" ADD COLUMN removable text"); err != nil {
		t.Fatalf("ALTER TABLE ADD COLUMN (removable): %v", err)
	}
	if _, err := env.pool.Exec(ctx,
		"INSERT INTO "+env.feed+" (id, name, removable) VALUES (3, 'BeforeDrop', 'ToBeRemoved')"); err != nil {
		t.Fatalf("INSERT vor der Entfernung: %v", err)
	}
	beforeRows := awaitChangesViewRows(t, env, "3", 1)

	if _, err := env.pool.Exec(ctx, "ALTER TABLE "+env.feed+" DROP COLUMN removable"); err != nil {
		t.Fatalf("ALTER TABLE DROP COLUMN: %v", err)
	}
	if _, err := env.pool.Exec(ctx,
		"INSERT INTO "+env.feed+" (id, name) VALUES (4, 'AfterDrop')"); err != nil {
		t.Fatalf("INSERT nach der Entfernung: %v", err)
	}

	// LH-FA-SCH-003 Happy Path (ADR-0063, Supersedes ADR-0058 Entscheidung 1):
	// die reale Spaltenentfernung löst denselben ErrIncompatibleSchemaChange-
	// Pfad aus wie jede andere inkompatible Relation-Änderung (Konvergenz
	// mit LH-FA-SCH-004, bereits von ADR-0059 Teilfrage 4 akzeptiert) —
	// sichtbarer schema-Fehler, kein stilles Auslassen.
	if got := awaitHeartbeatErrorClass(t, env, "schema"); got != "schema" {
		t.Fatalf("cdc.heartbeat.error_class nach DROP COLUMN: %q, wollen \"schema\"", got)
	}
	rows, err := queryChangesView(ctx, env, "4")
	if err != nil {
		t.Fatalf("cdc.changes-Lesung für id=4: %v", err)
	}
	if len(rows) != 0 {
		t.Fatalf("cdc.changes trägt id=4 nach dem gemeldeten schema-Fehler: %+v (Erwartung: der Erfassungspfad endete vor dem Commit dieser Transaktion)", rows)
	}

	// LH-FA-SCH-003 Boundary (unverändert ggü. ADR-0058): die vor der
	// Entfernung erfasste Change bleibt über cdc.changes unverändert lesbar,
	// inklusive des historischen Werts der entfernten Spalte.
	rereadBefore := awaitChangesViewRows(t, env, "3", 1)
	if string(rereadBefore[0].newData) != string(beforeRows[0].newData) {
		t.Fatalf("ältere Change nach DROP COLUMN: %s (vor der Entfernung: %s)", rereadBefore[0].newData, beforeRows[0].newData)
	}
	beforeImage := imageJSON(t, rereadBefore[0].newData)
	if beforeImage["removable"] != "ToBeRemoved" {
		t.Fatalf("ältere Change trägt den historischen Wert der entfernten Spalte nicht: %s", rereadBefore[0].newData)
	}
}

// TestE2EChangeTableMetadataExtensibility trägt `LH-FA-DAT-006`: eine
// reale, additive Erweiterung des internen Change-Schemas (`cdc.change`,
// `ADR-0017`) lässt bereits gespeicherte Changes unverändert über
// `cdc.changes` lesbar — Schreib- und Lesepfad (`InsertChange`/
// `SelectChanges`,
// `internal/adapters/driven/postgresstorage/queries/queries.go`) tragen
// explizite Spaltenlisten statt `SELECT *`/positioneller Vollständigkeit,
// und die deklarative `changes`-View (`tools/schema/schema.yaml`) trägt
// dieselbe Disziplin mit einer sichtbaren `columns:`-Signatur. Die neue
// Spalte bleibt additiv, nullable, und wird per `t.Cleanup` vor den
// nachfolgenden Testphasen wieder entfernt (`ADR-0058` Entscheidung 2).
func TestE2EChangeTableMetadataExtensibility(t *testing.T) {
	env := newE2EEnv(t, "feed_e2e_full")
	ctx := context.Background()

	if _, err := env.pool.Exec(ctx,
		"INSERT INTO "+env.feed+" (id, name) VALUES (500, 'BeforeMetadataExtension')"); err != nil {
		t.Fatalf("INSERT vor der Metadaten-Erweiterung: %v", err)
	}
	beforeRows := awaitChangesViewRows(t, env, "500", 1)

	const metadataColumn = "e2e_metadata_probe"
	if _, err := env.pool.Exec(ctx,
		"ALTER TABLE cdc.change ADD COLUMN "+metadataColumn+" jsonb DEFAULT NULL"); err != nil {
		t.Fatalf("ALTER TABLE cdc.change ADD COLUMN: %v", err)
	}
	t.Cleanup(func() {
		if _, err := env.pool.Exec(context.Background(),
			"ALTER TABLE cdc.change DROP COLUMN "+metadataColumn); err != nil {
			t.Errorf("ALTER TABLE cdc.change DROP COLUMN nach Testende: %v", err)
		}
	})

	// Happy Path (LH-FA-DAT-006): die vor der Erweiterung erfasste Change
	// bleibt über cdc.changes unverändert lesbar.
	rereadBefore := awaitChangesViewRows(t, env, "500", 1)
	if string(rereadBefore[0].newData) != string(beforeRows[0].newData) {
		t.Fatalf("vor der Erweiterung erfasste Change nach ALTER TABLE cdc.change ADD COLUMN: %s (davor: %s)", rereadBefore[0].newData, beforeRows[0].newData)
	}

	if _, err := env.pool.Exec(ctx,
		"INSERT INTO "+env.feed+" (id, name) VALUES (501, 'AfterMetadataExtension')"); err != nil {
		t.Fatalf("INSERT nach der Metadaten-Erweiterung: %v", err)
	}
	afterRows := awaitChangesViewRows(t, env, "501", 1)

	// Boundary (LH-FA-DAT-006): die danach eingefügte Zeile ist ebenfalls
	// unverändert lesbar.
	afterImage := imageJSON(t, afterRows[0].newData)
	if afterImage["id"] != "501" || afterImage["name"] != "AfterMetadataExtension" {
		t.Fatalf("Row Image nach der Metadaten-Erweiterung: %s", afterRows[0].newData)
	}

	// Boundary (LH-FA-DAT-006): die neue interne Spalte ist über die
	// bestehende cdc.changes-View nicht sichtbar — die View trägt eine
	// explizite columns:-Signatur (tools/schema/schema.yaml), keine
	// SELECT *-Kopplung an die vollständige Spaltenmenge von cdc.change.
	var exposed int
	if err := env.pool.QueryRow(ctx,
		"SELECT count(*) FROM information_schema.columns WHERE table_schema = 'cdc' AND table_name = 'changes' AND column_name = $1",
		metadataColumn,
	).Scan(&exposed); err != nil {
		t.Fatalf("cdc.changes-Spalten-Prüfung: %v", err)
	}
	if exposed != 0 {
		t.Fatalf("cdc.changes exponiert die neue interne Spalte %q (Erwartung: nicht sichtbar)", metadataColumn)
	}
}

// TestE2EHeartbeatHealthy trägt `LH-FA-ADM-002` (Betriebsstatus, Happy
// Path) am verdrahteten Feed-Container: `cdc.heartbeat` trägt für die
// laufende Quelle eine frische, fehlerfreie Lebenszeichen-Zeile —
// derselbe externe SQL-Lesezugriffsweg wie `awaitHeartbeatErrorClass`
// unten, gegen dieselbe Projektion (`tools/schema/nacharbeit-heartbeat.sql`,
// `LH-QA-OPS-002`). Die Alters-Schwelle spiegelt den
// Produktions-Healthcheck (`internal/bootstrap/wiring.go`,
// `heartbeatStaleAfter` = 3 × `heartbeatInterval` = 15s) — derselbe Wert,
// den der Compose-Healthcheck bereits gegen dieselbe Zeile prüft
// (`docker inspect --format '{{.State.Health.Status}}'` im Runner-Skript
// oben). Dieser Testfall läuft im ersten `go test`-Aufruf des
// Runner-Skripts, vor `TestE2ESchemaChangeIncompatibleTypeChange` — jener
// setzt `error_class` dauerhaft auf `schema` und würde den
// Happy-Path-Beleg sonst verdecken.
func TestE2EHeartbeatHealthy(t *testing.T) {
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

	const freshnessThresholdSeconds = 15.0
	deadline := time.Now().Add(30 * time.Second)
	var errorClass string
	var ageSeconds float64
	for time.Now().Before(deadline) {
		if err := pool.QueryRow(ctx,
			"SELECT coalesce(error_class, ''), age_seconds FROM cdc.heartbeat WHERE source_id = $1", e2eSource,
		).Scan(&errorClass, &ageSeconds); err != nil {
			t.Fatalf("cdc.heartbeat-Lesung: %v", err)
		}
		if errorClass == "" && ageSeconds < freshnessThresholdSeconds {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatalf("cdc.heartbeat innerhalb der Zeitspanne nicht frisch/fehlerfrei: error_class=%q age_seconds=%.1f (Schwelle %.0fs)", errorClass, ageSeconds, freshnessThresholdSeconds)
}

// awaitHeartbeatErrorClass liest `cdc.heartbeat.error_class` der Quelle
// über `env.pool`, bis der erwartete Fehlerzustand ansteht (Polling mit
// Test-Zeitgrenze) — derselbe externe SQL-Lesezugriffsweg wie
// `cdc.changes` (`awaitChangesViewRows` oben), gegen die Projektion aus
// `tools/schema/nacharbeit-heartbeat.sql` (`LH-FA-ADM-003`,
// `LH-QA-REL-003`). Eine leere Zeichenkette trägt `NULL` (Normalbetrieb,
// `error_class` noch nicht gesetzt).
func awaitHeartbeatErrorClass(t *testing.T, env *e2eEnv, want string) string {
	t.Helper()
	ctx := context.Background()
	deadline := time.Now().Add(30 * time.Second)
	got := ""
	for time.Now().Before(deadline) {
		if err := env.pool.QueryRow(ctx,
			"SELECT coalesce(error_class, '') FROM cdc.heartbeat WHERE source_id = $1", e2eSource,
		).Scan(&got); err != nil {
			t.Fatalf("cdc.heartbeat-Lesung: %v", err)
		}
		if got == want {
			return got
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatalf("cdc.heartbeat.error_class innerhalb der Zeitspanne %q, wollen %q", got, want)
	return ""
}

// TestE2ESchemaChangeIncompatibleTypeChange trägt `LH-FA-SCH-004` an der
// aktivierten Tabelle `feed_e2e_schema` (nach
// TestE2ESchemaChangeAddColumn, mit der dort hinzugefügten Spalte
// `extra`) in beiden im Slice-Plan (`slice-030`, §6) benannten Fällen:
//
//  1. PostgreSQL lehnt eine tatsächlich inkompatible Typänderung bereits
//     selbst über die DDL ab, bevor CDC sie überhaupt sieht — der reale
//     Boundary-Fall aus `LH-FA-SCH-004`.
//  2. Eine Typänderung, die PostgreSQL über eine `USING`-Klausel bei
//     durchgehend konvertierbaren Bestandsdaten zulässt, sendet
//     `pgoutput` dennoch eine neue Relation-Nachricht für die geänderte
//     Spalte (andere PostgreSQL-Typ-OID): nicht sicher als Obermenge
//     erkennbar (`relationOther`,
//     `internal/adapters/driving/replication/mapper/mapper.go`
//     `classifyRelationColumns`) — der Negative-Fall aus `LH-FA-SCH-004`.
//     `Assembler.observeRelation` meldet ihn sichtbar als Fehler der
//     Klasse `schema` (`mapper.ErrIncompatibleSchemaChange`), statt die
//     Änderung still zu übernehmen.
//
// Der sichtbare Fehler beendet den Erfassungspfad des Feed-Containers
// (`Assembler.Consume` -> `receive.Stream.Run` -> `bootstrap.Run` ->
// `os.Exit(1)`, `restart: "no"` in compose.yaml trägt keinen
// Neustart-Vertrag) — der Container bleibt danach beendet stehen. Dieser
// Test läuft deshalb als eigener, letzter `go test`-Aufruf des
// Compose-Laufs (`tools/harness/run-integration-tests.sh`), nach jedem
// Schritt, der den laufenden Feed-Container noch braucht (Lasttest-Beleg,
// Black-Box-CLI-Rundlauf). Beobachtet wird der Fehlerzustand über den
// bestehenden externen Lesezugriffsweg `cdc.heartbeat` (`error_class`,
// `LH-FA-ADM-003`) — dieselbe SQL-Lesedisziplin wie `cdc.changes`. Die
// Zeile id=11 bleibt dabei unerfasst: ihre Transaktion trägt die
// auslösende Relation-Nachricht vor ihrem eigenen Commit, und der
// Erfassungspfad endet, bevor sie committed wird.
func TestE2ESchemaChangeIncompatibleTypeChange(t *testing.T) {
	env := newE2EEnv(t, "feed_e2e_schema")
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

	// Fall 2 (Negative-Fall, s. Funktionskommentar): dieselbe Typänderung
	// gelingt an PostgreSQL mit durchgehend konvertierbaren Bestandsdaten
	// (id=1 "100", id=2 "200" aus TestE2ESchemaChangeAddColumn) — CDC
	// meldet sie trotzdem sichtbar als Fehlerklasse `schema`, weil die
	// Spalte `amount` ihre PostgreSQL-Typ-OID wechselt.
	if _, err := env.pool.Exec(ctx,
		"ALTER TABLE "+env.feed+" ALTER COLUMN amount TYPE integer USING amount::integer"); err != nil {
		t.Fatalf("ALTER COLUMN TYPE integer mit konvertierbaren Bestandsdaten: %v", err)
	}
	if _, err := env.pool.Exec(ctx,
		"INSERT INTO "+env.feed+" (id, name, amount, extra) VALUES (11, 'TypeChanged', 300, 'Y')"); err != nil {
		t.Fatalf("INSERT nach der zugelassenen Typänderung: %v", err)
	}

	if got := awaitHeartbeatErrorClass(t, env, "schema"); got != "schema" {
		t.Fatalf("cdc.heartbeat.error_class nach der Typänderung: %q, wollen \"schema\"", got)
	}

	rows, err := queryChangesView(ctx, env, "11")
	if err != nil {
		t.Fatalf("cdc.changes-Lesung für id=11: %v", err)
	}
	if len(rows) != 0 {
		t.Fatalf("cdc.changes trägt id=11 nach dem gemeldeten schema-Fehler: %+v (Erwartung: der Erfassungspfad endete vor dem Commit dieser Transaktion)", rows)
	}
}

// abdeckungZeilenPraefix kennzeichnet die Zeilen, die der Runner aus dem
// Testausgang liest (`tools/harness/run-integration-tests.sh`): dahinter
// steht ein über `|` getrennter Satz aus Quelldatei, Quellzeile, Nachweis,
// Kennungen und Kurzbeschreibung. Der Präfix macht sie von der übrigen
// Testausgabe unterscheidbar.
const abdeckungZeilenPraefix = "ABDECKUNG|"

// abdeckungKennungMuster trifft die Kennungen eines E2E-Doc-Kommentars:
// die Vertrags- und Technik-Kennungen des Lastenhefts/Pflichtenhefts (samt
// der Verfeinerungs-Form `…-001.a`) adressieren eine Zeile der
// Abdeckungstabelle; die ADR- und Sicht-Kennungen stehen im Kommentar,
// benennen aber eine Entscheidung statt einer Anforderung, und die
// Lebenszyklus-Kennungen (`slice-NNN`, `welle-NN`) nennen einen Vorgang
// statt eines Nachweises — beide fallen deshalb aus der
// Beschreibungsspalte heraus: die Tabelle führt Nachweise, keine
// Entstehungsgeschichte.
var abdeckungKennungMuster = regexp.MustCompile(`LH-(?:FA|QA)-[A-Z]{3}-\d{3}(?:\.[a-z])?|SPEC-\d{3}|ARC-\d{3}|ADR-\d{4}|slice-\d{3}|welle-\d{1,2}`)

// abdeckungKurzformEinleitungen sind die Zeichen, mit denen ein
// E2E-Kommentar eine Bereichs- oder Nachbar-Angabe hinter einer Kennung
// einleitet (`…003`, `/`004``, `...005`). Hinter einer Kennung wird daraus
// die fehlende Kennungs-Familie abgeleitet; eine Einleitung ohne laufende
// Nummer ist eine Auslassung im Fließtext, keine Kurzform.
var abdeckungKurzformEinleitungen = []string{"…", "...", "/"}

// abdeckungFamilieMuster trennt eine Kennung in ihre Präfix-Familie und
// ihre laufende Nummer (`LH-FA-CAP-001` -> `LH-FA-CAP`, 1).
var abdeckungFamilieMuster = regexp.MustCompile(`^(.+)-(\d{3,4})$`)

// abdeckungKurzformSpanne ist die größte Spanne, die eine Kurzform
// aufspannen darf. Sie begrenzt einen Tippfehler in der Endzahl, damit aus
// `…999` keine 998 erfundenen Zeilen werden.
const abdeckungKurzformSpanne = 40

// abdeckungAufraeumRegeln normalisieren die Kurzbeschreibung nach dem
// Entfernen der Kennungen: Aufzählungs-Trenner, leere Klammern und die
// Leerzeichen, die eine Kennung hinterlässt, fallen weg.
var abdeckungAufraeumRegeln = []struct {
	muster *regexp.Regexp
	ersatz string
}{
	// Ein Aufzählungs-Trenner direkt hinter der öffnenden Klammer bleibt
	// stehen, wenn die erste Kennung der Aufzählung wegfällt.
	{regexp.MustCompile(`([(\[])\s*[,;]\s*`), "$1"},
	// Leere Klammer und Klammer mit zurückgebliebenen Trennern.
	{regexp.MustCompile(`\(\s*(?:[,;]\s*)*\)`), ""},
	{regexp.MustCompile(`[,;]\s*\)`), ")"},
	{regexp.MustCompile(`\s+([,.;:)\]])`), "$1"},
	{regexp.MustCompile(`([(\[])\s+`), "$1"},
	{regexp.MustCompile(`\s{2,}`), " "},
	{regexp.MustCompile(`^[,;:]\s*`), ""},
}

// TestAbdeckungstabelleZeilen leitet den Go-Anteil der
// E2E-Abdeckungstabelle (`docs/user/e2e-abdeckung.md`) aus dem Quelltext
// dieses Pakets ab — nicht aus dem, was gelaufen ist: eine
// `func TestE2E*`, die kein `-run`-Muster des Runners trifft, erscheint
// trotzdem in der Tabelle (`BEO-PGC/test-runner-stiller-ausschluss`).
// Funktionsname und Quellzeile kommen aus dem AST, die Spec-Kennungen und
// die Kurzbeschreibung aus dem Doc-Kommentar. Eine `func TestE2E*` ohne
// Spec-Kennung, ohne Doc-Kommentar oder mit einer nicht auflösbaren
// Kurzform bricht den Erzeuger sichtbar ab, statt die Zeile
// stillschweigend wegzulassen.
func TestAbdeckungstabelleZeilen(t *testing.T) {
	_, quelle, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("E2E-Abdeckungstabelle: Quelldatei dieses Pakets nicht auflösbar")
	}
	zeilen, err := abdeckungsZeilen(quelle)
	if err != nil {
		t.Fatalf("E2E-Abdeckungstabelle: %v", err)
	}
	if len(zeilen) == 0 {
		t.Fatal("E2E-Abdeckungstabelle: kein func TestE2E* im Paket gefunden — der Zeilensatz wäre leer")
	}
	for _, zeile := range zeilen {
		t.Log(abdeckungZeilenPraefix + zeile)
	}
}

// abdeckungsZeilen liefert je `func TestE2E*` der Go-Dateien neben
// quelleDatei eine Zeile `Quelldatei|Quellzeile|Nachweis|Kennungen|
// Kurzbeschreibung`, in Datei- und Quellzeilen-Reihenfolge.
func abdeckungsZeilen(quelleDatei string) ([]string, error) {
	verzeichnis := filepath.Dir(quelleDatei)
	eintraege, err := os.ReadDir(verzeichnis)
	if err != nil {
		return nil, err
	}
	zeilen := make([]string, 0)
	for _, eintrag := range eintraege {
		if eintrag.IsDir() || !strings.HasSuffix(eintrag.Name(), ".go") {
			continue
		}
		datei := filepath.Join(verzeichnis, eintrag.Name())
		quelldatei, err := abdeckungsModulpfad(datei)
		if err != nil {
			return nil, err
		}
		zeilen, err = abdeckungsZeilenEinerDatei(datei, quelldatei, zeilen)
		if err != nil {
			return nil, err
		}
	}
	return zeilen, nil
}

// abdeckungsZeilenEinerDatei hängt die Zeilen der `func TestE2E*` einer
// Datei an zeilen an.
func abdeckungsZeilenEinerDatei(datei, quelldatei string, zeilen []string) ([]string, error) {
	baumsatz := token.NewFileSet()
	geparst, err := parser.ParseFile(baumsatz, datei, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	for _, deklaration := range geparst.Decls {
		funktion, ok := deklaration.(*ast.FuncDecl)
		if !ok || funktion.Recv != nil || !strings.HasPrefix(funktion.Name.Name, "TestE2E") {
			continue
		}
		zeile, err := abdeckungsZeile(baumsatz, quelldatei, funktion)
		if err != nil {
			return nil, err
		}
		zeilen = append(zeilen, zeile)
	}
	return zeilen, nil
}

// abdeckungsZeile baut die Zeile einer einzelnen Testfunktion; die
// Spec-Kennung ist Pflicht, sonst wäre die Zeile nicht adressierbar.
func abdeckungsZeile(baumsatz *token.FileSet, quelldatei string, funktion *ast.FuncDecl) (string, error) {
	nachweis := funktion.Name.Name
	if funktion.Doc == nil {
		return "", fmt.Errorf("%s: kein Doc-Kommentar — die Zeile trägt sonst keine Spec-Kennung", nachweis)
	}
	kennungen, ohneKennungen, err := abdeckungsKennungen(nachweis, funktion.Doc.Text())
	if err != nil {
		return "", err
	}
	if len(kennungen) == 0 {
		return "", fmt.Errorf("%s: keine Spec-Kennung im Doc-Kommentar — der Nachweis wäre keiner Anforderung zugeordnet", nachweis)
	}
	beschreibung, err := abdeckungsKurzbeschreibung(nachweis, ohneKennungen)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{
		quelldatei,
		strconv.Itoa(baumsatz.Position(funktion.Pos()).Line),
		nachweis,
		strings.Join(kennungen, ","),
		beschreibung,
	}, "|"), nil
}

// abdeckungsKennungen trennt den Doc-Kommentar in die adressierten
// Spec-Kennungen (Lastenheft/Pflichtenheft, Kurzformen aufgelöst) und den
// Text ohne jede Kennung — der Text ist die Quelle der Kurzbeschreibung.
func abdeckungsKennungen(nachweis, kommentar string) ([]string, string, error) {
	kennungen := make([]string, 0)
	var rest strings.Builder
	text := kommentar
	for {
		stelle := abdeckungKennungMuster.FindStringIndex(text)
		if stelle == nil {
			rest.WriteString(text)
			break
		}
		anfang, ende := stelle[0], stelle[1]
		kennung := text[anfang:ende]
		// Ein Inline-Code-Span, der die Kennung trägt, fällt mit ihr weg:
		// bliebe er stehen, entstünde eine leere Spanne, und ein
		// halbierter Span (`` `LH-…`…003 ``) ließe ein unpaariges
		// Backtick zurück.
		if anfang > 0 && text[anfang-1] == '`' && ende < len(text) && text[ende] == '`' {
			anfang--
			ende++
		}
		weitere, kurzformLaenge, err := abdeckungsKurzform(nachweis, kennung, text[ende:])
		if err != nil {
			return nil, "", err
		}
		ende += kurzformLaenge
		rest.WriteString(text[:anfang])
		if abdeckungsAdressiert(kennung) {
			kennungen = append(kennungen, kennung)
			kennungen = append(kennungen, weitere...)
		}
		text = text[ende:]
	}
	if strings.Count(rest.String(), "`")%2 != 0 {
		return nil, "", fmt.Errorf("%s: Doc-Kommentar hinterlässt nach dem Entfernen der Kennungen ein unpaariges Inline-Code-Zeichen", nachweis)
	}
	gesehen := make(map[string]bool, len(kennungen))
	eindeutig := kennungen[:0]
	for _, kennung := range kennungen {
		if !gesehen[kennung] {
			gesehen[kennung] = true
			eindeutig = append(eindeutig, kennung)
		}
	}
	return eindeutig, rest.String(), nil
}

// abdeckungsAdressiert sagt, ob eine Kennung in die Kennungsspalte gehört:
// die Anforderungs-Kennungen der beiden Spec-Straten Vertrag und Technik.
func abdeckungsAdressiert(kennung string) bool {
	return strings.HasPrefix(kennung, "LH-") || strings.HasPrefix(kennung, "SPEC-")
}

// abdeckungsKurzform löst die Kurzform hinter einer Kennung auf und
// liefert die dadurch zusätzlich adressierten Kennungen samt der Länge, die
// die Kurzform im Kommentar einnimmt. Trägt die Kennung keine laufende
// Nummer, liegt die Angabe nicht aufsteigend, sprengt sie die
// Kurzform-Spanne oder bleibt der Inline-Code-Span der Nummer offen, bricht
// der Erzeuger sichtbar ab, statt zu raten.
func abdeckungsKurzform(nachweis, kennung, rest string) ([]string, int, error) {
	laenge := 0
	for _, einleitung := range abdeckungKurzformEinleitungen {
		if strings.HasPrefix(rest, einleitung) {
			laenge = len(einleitung)
			break
		}
	}
	if laenge == 0 {
		return nil, 0, nil
	}
	zahlText := rest[laenge:]
	umspannt := strings.HasPrefix(zahlText, "`")
	if umspannt {
		zahlText = zahlText[1:]
	}
	ziffern := 0
	for ziffern < len(zahlText) && zahlText[ziffern] >= '0' && zahlText[ziffern] <= '9' {
		ziffern++
	}
	if ziffern == 0 {
		return nil, 0, nil
	}
	if ziffern != 3 {
		return nil, 0, fmt.Errorf("%s: Kurzform %q hinter %s trägt keine dreistellige laufende Nummer", nachweis, rest[:laenge+ziffern], kennung)
	}
	gelesen := laenge + ziffern
	if umspannt {
		gelesen++
		if !strings.HasPrefix(zahlText[ziffern:], "`") {
			return nil, 0, fmt.Errorf("%s: Kurzform %q hinter %s öffnet einen Inline-Code-Span ohne Abschluss", nachweis, rest[:gelesen], kennung)
		}
		gelesen++
	}
	familie, anfangsZahl, err := abdeckungsFamilie(kennung)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: Kurzform %q hinter %s ist nicht auflösbar (%v)", nachweis, rest[:gelesen], kennung, err)
	}
	endZahl, err := strconv.Atoi(zahlText[:ziffern])
	if err != nil {
		return nil, 0, fmt.Errorf("%s: Kurzform %q hinter %s trägt keine Zahl: %v", nachweis, rest[:gelesen], kennung, err)
	}
	if endZahl <= anfangsZahl || endZahl-anfangsZahl > abdeckungKurzformSpanne {
		return nil, 0, fmt.Errorf("%s: Kurzform %q hinter %s nennt keine aufsteigende Nachbar-Familie", nachweis, rest[:gelesen], kennung)
	}
	weitere := make([]string, 0, endZahl-anfangsZahl)
	for zahl := anfangsZahl + 1; zahl <= endZahl; zahl++ {
		weitere = append(weitere, fmt.Sprintf("%s-%03d", familie, zahl))
	}
	return weitere, gelesen, nil
}

// abdeckungsFamilie zerlegt eine Kennung in Präfix-Familie und laufende
// Nummer; eine Kennung ohne laufende Nummer (etwa die Verfeinerungs-Form
// `…-001.a`) trägt keine Kurzform und endet hier sichtbar.
func abdeckungsFamilie(kennung string) (string, int, error) {
	teile := abdeckungFamilieMuster.FindStringSubmatch(kennung)
	if teile == nil {
		return "", 0, fmt.Errorf("die Kennung trägt keine laufende Nummer")
	}
	zahl, err := strconv.Atoi(teile[2])
	if err != nil {
		return "", 0, err
	}
	return teile[1], zahl, nil
}

// abdeckungsKommentarAbsatz liefert den ersten Absatz eines
// Doc-Kommentars. Die Kurzbeschreibung bleibt damit beim
// Zusammenfassungs-Absatz: eine anschließende nummerierte Aufzählung trägt
// ihre Kennungen, aber nicht mehr ihren Fließtext in die Tabellenzelle.
func abdeckungsKommentarAbsatz(kommentar string) string {
	if trenner := strings.Index(kommentar, "\n\n"); trenner >= 0 {
		return kommentar[:trenner]
	}
	return kommentar
}

// abdeckungsKurzbeschreibung normalisiert den kennungsfreien Kommentartext
// zum ersten Satz seines ersten Absatzes.
func abdeckungsKurzbeschreibung(nachweis, text string) (string, error) {
	beschreibung := strings.Join(strings.Fields(abdeckungsKommentarAbsatz(text)), " ")
	beschreibung = strings.TrimSpace(strings.TrimPrefix(beschreibung, nachweis))
	if ende := abdeckungsSatzende(beschreibung); ende >= 0 {
		beschreibung = beschreibung[:ende]
	}
	for _, regel := range abdeckungAufraeumRegeln {
		beschreibung = regel.muster.ReplaceAllString(beschreibung, regel.ersatz)
	}
	beschreibung = strings.TrimSpace(beschreibung)
	if beschreibung == "" {
		return "", fmt.Errorf("%s: Doc-Kommentar trägt nach dem Entfernen der Kennungen keine Kurzbeschreibung", nachweis)
	}
	if strings.Count(beschreibung, "`")%2 != 0 {
		return "", fmt.Errorf("%s: Kurzbeschreibung hinterlässt ein unpaariges Inline-Code-Zeichen", nachweis)
	}
	// Der Zeilensatz ist über `|` getrennt; ein Trennerzeichen im Text
	// würde ihn zerlegen.
	return strings.ReplaceAll(beschreibung, "|", "/"), nil
}

// abdeckungsSatzende liefert den Index hinter dem ersten Satzende — dem
// ersten Punkt außerhalb eines Inline-Code-Spans, dem Leerraum oder das
// Textende folgt und dem keine Ziffer vorausgeht (Aufzählungs- und
// Zahlmarken wie „1." bleiben damit im Satz) — oder -1, wenn der Text kein
// solches Ende trägt.
func abdeckungsSatzende(text string) int {
	imSpan := false
	for i := 0; i < len(text); i++ {
		switch {
		case text[i] == '`':
			imSpan = !imSpan
		case text[i] == '.' && !imSpan:
			if i > 0 && text[i-1] >= '0' && text[i-1] <= '9' {
				continue
			}
			if i+1 == len(text) || text[i+1] == ' ' {
				return i + 1
			}
		}
	}
	return -1
}

// abdeckungsModulpfad übersetzt einen absoluten Dateipfad in den
// modul-relativen (Repository-) Pfad: der Aufstieg endet an der Datei
// `go.mod`.
func abdeckungsModulpfad(datei string) (string, error) {
	verzeichnis := filepath.Dir(datei)
	for {
		if _, err := os.Stat(filepath.Join(verzeichnis, "go.mod")); err == nil {
			return filepath.Rel(verzeichnis, datei)
		}
		uebergeordnet := filepath.Dir(verzeichnis)
		if uebergeordnet == verzeichnis {
			return "", fmt.Errorf("%s: kein go.mod im Aufstieg gefunden", datei)
		}
		verzeichnis = uebergeordnet
	}
}
