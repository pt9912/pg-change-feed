package postgresstorage_test

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage"
	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Die Antrags-Queue-Tests laufen gegen dieselbe reale PostgreSQL-Instanz wie
// die übrigen Store-Tests (`make test-store`, `ADR-0030`); ohne DSN
// überspringen sie. `cdc.administration_request` trägt der d-migrate-Rollout
// (`tools/schema/schema.yaml`, `ADR-0043`), `cdc.enable_table`/
// `cdc.disable_table` die Ausweichform (`tools/schema/nacharbeit-administration.sql`,
// `ADR-0050`) — beide Objektklassen entstehen im selben `make
// schema-rollout`-Lauf, vor diesem Testlauf.
//
// Kopplung: diese Datei muss vor `consumerstate_test.go`/`store_test.go`/
// `tableactivation_test.go` laufen — deren Tests räumen das Schema per
// `DROP SCHEMA cdc CASCADE` samt hand-DDL-Neuaufbau zurück (`newTestStore`/
// `newTestActivation`), der weder `cdc.administration_request` noch die
// beiden Funktionen mitträgt. `go test` fährt die Testdateien eines Pakets
// in Datei-Namensordnung; „administrationrequest_test.go" sortiert vor
// jeder der genannten Dateien (`a` < `c`/`st`/`t`) und läuft deshalb zuerst
// — dasselbe Muster wie bereits `schemastore_test.go`/`sqlviews_test.go` für
// dieselbe Ausgangslage.
const administrationRequestSource = "src-administration"

// newTestAdministrationRequestPool baut den Pool gegen die Test-Instanz und
// trägt den Quelle-Bestand idempotent fort (`ON CONFLICT DO NOTHING`, wie
// `newTestSchemaStore`).
func newTestAdministrationRequestPool(t *testing.T) (*pgxpool.Pool, string) {
	t.Helper()
	dsn := os.Getenv("CDC_STORE_TEST_DSN")
	if dsn == "" {
		t.Skip("CDC_STORE_TEST_DSN nicht gesetzt — reale PostgreSQL-Tests laufen über make test-store")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("Verbindungsaufbau: %v", err)
	}
	t.Cleanup(pool.Close)

	var functions int
	if err := pool.QueryRow(ctx,
		"SELECT count(*) FROM pg_proc p JOIN pg_namespace n ON n.oid = p.pronamespace WHERE n.nspname = 'cdc' AND p.proname IN ('enable_table', 'disable_table', 'exclude_column', 'include_column', 'backfill_table', 'set_transformation', 'remove_transformation')",
	).Scan(&functions); err != nil {
		t.Fatalf("Funktions-Prüfung: %v", err)
	}
	if functions != 7 {
		t.Fatalf("cdc.enable_table/cdc.disable_table/cdc.exclude_column/cdc.include_column/cdc.backfill_table/cdc.set_transformation/cdc.remove_transformation fehlen — der Schema-Rollout über make schema-rollout trägt sie (ADR-0050, LH-FA-CFG-005, LH-FA-CAP-009, LH-FA-CFG-007, tools/schema/nacharbeit-administration.sql); der test-store-Lauf rollt sie vor dem Testlauf aus")
	}

	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.source (source_id, name) VALUES ($1, 'Administrations-Quelle') ON CONFLICT (source_id) DO NOTHING", administrationRequestSource,
	); err != nil {
		t.Fatalf("Quelle-Zeile: %v", err)
	}
	return pool, dsn
}

// listenForAdministrationNotify baut eine dedizierte pgx-Verbindung für
// `LISTEN` auf — Pool-Verbindungen sind für `WaitForNotification` ungeeignet,
// weil jede Anfrage eine beliebige Pool-Verbindung ziehen kann. Der
// Kanal-Name `cdc_administration` ist Implementer-Entscheidung
// (tools/schema/nacharbeit-administration.sql).
func listenForAdministrationNotify(t *testing.T, ctx context.Context, dsn string) *pgx.Conn {
	t.Helper()
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		t.Fatalf("Listener-Verbindungsaufbau: %v", err)
	}
	t.Cleanup(func() { conn.Close(context.Background()) })
	if _, err := conn.Exec(ctx, "LISTEN cdc_administration"); err != nil {
		t.Fatalf("LISTEN cdc_administration: %v", err)
	}
	return conn
}

// TestAdministrationRequestEnableTableWritesPendingRequestAndNotifies trägt
// den realen Antrags-Weg (ADR-0050): der Funktionsaufruf legt eine Zeile mit
// Status `pending` an und sendet `pg_notify`, real über `LISTEN` beobachtet
// — kein direkter Zugriff auf `cdc.source_table`/`cdc.schema_version`/die
// Publication (die Verarbeitung trägt die Administrations-Goroutine).
func TestAdministrationRequestEnableTableWritesPendingRequestAndNotifies(t *testing.T) {
	pool, dsn := newTestAdministrationRequestPool(t)
	ctx := context.Background()

	listener := listenForAdministrationNotify(t, ctx, dsn)

	var requestID string
	if err := pool.QueryRow(ctx,
		"SELECT cdc.enable_table($1, $2, $3)", administrationRequestSource, "public", "orders_enable",
	).Scan(&requestID); err != nil {
		t.Fatalf("cdc.enable_table: %v", err)
	}
	if requestID == "" {
		t.Fatalf("cdc.enable_table lieferte eine leere Antrags-ID")
	}

	notifyCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	notification, err := listener.WaitForNotification(notifyCtx)
	if err != nil {
		t.Fatalf("WaitForNotification: %v", err)
	}
	if notification.Channel != "cdc_administration" {
		t.Fatalf("Notify-Kanal = %q, wollen cdc_administration", notification.Channel)
	}
	if notification.Payload != requestID {
		t.Fatalf("Notify-Payload = %q, wollen die Antrags-ID %q", notification.Payload, requestID)
	}

	var sourceID, schemaName, tableName, requestKind, status string
	if err := pool.QueryRow(ctx,
		"SELECT source_id, schema_name, table_name, request_kind, status FROM cdc.administration_request WHERE administration_request_id = $1",
		requestID,
	).Scan(&sourceID, &schemaName, &tableName, &requestKind, &status); err != nil {
		t.Fatalf("Antrags-Zeile lesen: %v", err)
	}
	if sourceID != administrationRequestSource || schemaName != "public" || tableName != "orders_enable" {
		t.Fatalf("Antrags-Zeile: source=%q schema=%q table=%q", sourceID, schemaName, tableName)
	}
	if requestKind != "enable" {
		t.Fatalf("request_kind = %q, wollen enable", requestKind)
	}
	if status != "pending" {
		t.Fatalf("status = %q, wollen pending", status)
	}
}

// TestAdministrationRequestEnableTableRequiresCdcAdminMembership belegt die
// Least-Privilege-Durchsetzung des REVOKE/GRANT-Paars in
// nacharbeit-administration.sql (ADR-0047): eine Rolle ohne
// cdc_admin-Mitgliedschaft scheitert am Aufruf — derselbe `SET
// ROLE`-Mechanismus wie roles_test.go, `permissionDenied` von dort
// wiederverwendet (SQLSTATE 42501, dieselbe PostgreSQL-Fehlerklasse, die
// „permission denied for function" auslöst). Der Spaltenausschluss trägt
// dieselbe Zusage; er prüft zusätzlich, dass die vierteilige Funktion im
// REVOKE/GRANT-Paar wirklich eingeschlossen ist.
func TestAdministrationRequestEnableTableRequiresCdcAdminMembership(t *testing.T) {
	pool, _ := newTestAdministrationRequestPool(t)
	ctx := context.Background()

	conn, err := pool.Acquire(ctx)
	if err != nil {
		t.Fatalf("Verbindung reservieren: %v", err)
	}
	defer func() {
		_, _ = conn.Exec(ctx, "RESET ROLE")
		conn.Release()
	}()

	if _, err := conn.Exec(ctx, "SET ROLE cdc_reader"); err != nil {
		t.Fatalf("SET ROLE cdc_reader: %v", err)
	}

	var requestID string
	err = conn.QueryRow(ctx,
		"SELECT cdc.enable_table($1, $2, $3)", administrationRequestSource, "public", "orders_reader_denied",
	).Scan(&requestID)
	if !permissionDenied(err) {
		t.Fatalf("cdc_reader SELECT cdc.enable_table(...): erwartet SQLSTATE 42501 (insufficient_privilege), erhalten %v", err)
	}

	err = conn.QueryRow(ctx,
		"SELECT cdc.exclude_column($1, $2, $3, $4)", administrationRequestSource, "public", "orders_reader_denied", "secret",
	).Scan(&requestID)
	if !permissionDenied(err) {
		t.Fatalf("cdc_reader SELECT cdc.exclude_column(...): erwartet SQLSTATE 42501 (insufficient_privilege), erhalten %v", err)
	}
}

// TestAdministrationRequestDisableTableWritesPendingRequestAndNotifies
// spiegelt den vorigen Test für `cdc.disable_table` — dieselbe Zusage,
// anderer `request_kind`.
func TestAdministrationRequestDisableTableWritesPendingRequestAndNotifies(t *testing.T) {
	pool, dsn := newTestAdministrationRequestPool(t)
	ctx := context.Background()

	listener := listenForAdministrationNotify(t, ctx, dsn)

	var requestID string
	if err := pool.QueryRow(ctx,
		"SELECT cdc.disable_table($1, $2, $3)", administrationRequestSource, "public", "orders_disable",
	).Scan(&requestID); err != nil {
		t.Fatalf("cdc.disable_table: %v", err)
	}

	notifyCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	notification, err := listener.WaitForNotification(notifyCtx)
	if err != nil {
		t.Fatalf("WaitForNotification: %v", err)
	}
	if notification.Payload != requestID {
		t.Fatalf("Notify-Payload = %q, wollen die Antrags-ID %q", notification.Payload, requestID)
	}

	var requestKind, status string
	if err := pool.QueryRow(ctx,
		"SELECT request_kind, status FROM cdc.administration_request WHERE administration_request_id = $1",
		requestID,
	).Scan(&requestKind, &status); err != nil {
		t.Fatalf("Antrags-Zeile lesen: %v", err)
	}
	if requestKind != "disable" {
		t.Fatalf("request_kind = %q, wollen disable", requestKind)
	}
	if status != "pending" {
		t.Fatalf("status = %q, wollen pending", status)
	}
}

// TestAdministrationRequestColumnRequestsCarryColumnAndExitStatus trägt den
// realen Antrags-Weg des Spaltenausschlusses/-einschlusses
// (`LH-FA-CFG-005`, `ADR-0059`): die beiden vierteiligen Funktionen legen
// eine Zeile mit `column_name` und `request_kind` an, der Adapter liest
// beides zurück, und der Ergebnis-Vermerk endet für eine vorhandene Spalte
// als `applied`, für eine nicht existierende als `failed` samt Fehlertext
// (`LH-FA-CFG-005` Negative) — dieselbe asynchrone Sichtbarkeit wie bei den
// beiden Tabellen-Antragsarten.
func TestAdministrationRequestColumnRequestsCarryColumnAndExitStatus(t *testing.T) {
	pool, dsn := newTestAdministrationRequestPool(t)
	ctx := context.Background()

	adapter, err := postgresstorage.NewAdministrationRequest(ctx, dsn)
	if err != nil {
		t.Fatalf("NewAdministrationRequest: %v", err)
	}
	t.Cleanup(adapter.Close)

	var excludeID, includeID string
	if err := pool.QueryRow(ctx,
		"SELECT cdc.exclude_column($1, $2, $3, $4)", administrationRequestSource, "public", "orders_columns_exclude", "secret",
	).Scan(&excludeID); err != nil {
		t.Fatalf("cdc.exclude_column: %v", err)
	}
	if err := pool.QueryRow(ctx,
		"SELECT cdc.include_column($1, $2, $3, $4)", administrationRequestSource, "public", "orders_columns_include", "missing",
	).Scan(&includeID); err != nil {
		t.Fatalf("cdc.include_column: %v", err)
	}

	pending, err := adapter.ListPending(ctx)
	if err != nil {
		t.Fatalf("ListPending: %v", err)
	}
	byID := map[model.AdministrationRequestID]model.AdministrationRequest{}
	for _, request := range pending {
		byID[request.ID] = request
	}
	excludeRequest, found := byID[model.AdministrationRequestID(excludeID)]
	if !found {
		t.Fatalf("ListPending trägt nicht den exclude_column-Antrag %q: %+v", excludeID, pending)
	}
	if excludeRequest.Kind != model.AdministrationRequestExcludeColumn || excludeRequest.Column != "secret" ||
		excludeRequest.Table != "orders_columns_exclude" {
		t.Fatalf("exclude_column-Antrag: %+v", excludeRequest)
	}
	includeRequest, found := byID[model.AdministrationRequestID(includeID)]
	if !found {
		t.Fatalf("ListPending trägt nicht den include_column-Antrag %q: %+v", includeID, pending)
	}
	if includeRequest.Kind != model.AdministrationRequestIncludeColumn || includeRequest.Column != "missing" {
		t.Fatalf("include_column-Antrag: %+v", includeRequest)
	}

	if err := adapter.MarkApplied(ctx, excludeRequest.ID); err != nil {
		t.Fatalf("MarkApplied: %v", err)
	}
	failure := fmt.Sprintf("%v: public.orders_columns_include.missing", inbound.ErrSourceColumnMissing)
	if err := adapter.MarkFailed(ctx, includeRequest.ID, failure); err != nil {
		t.Fatalf("MarkFailed: %v", err)
	}

	var appliedStatus, failedStatus, failedMessage string
	if err := pool.QueryRow(ctx, "SELECT status FROM cdc.administration_request WHERE administration_request_id = $1", excludeID).Scan(&appliedStatus); err != nil {
		t.Fatalf("Status exclude_column-Antrag lesen: %v", err)
	}
	if appliedStatus != "applied" {
		t.Fatalf("Status exclude_column-Antrag = %q, wollen applied", appliedStatus)
	}
	if err := pool.QueryRow(ctx, "SELECT status, error_message FROM cdc.administration_request WHERE administration_request_id = $1", includeID).Scan(&failedStatus, &failedMessage); err != nil {
		t.Fatalf("Status include_column-Antrag lesen: %v", err)
	}
	if failedStatus != "failed" || failedMessage != failure {
		t.Fatalf("Status/Fehlertext include_column-Antrag = %q/%q", failedStatus, failedMessage)
	}
}

// TestAdministrationRequestAdapterListPendingMarkAppliedMarkFailed trägt
// die Gegenrichtung des Ports (`ADR-0050`): der Adapter liest die über SQL
// geschriebenen Anträge und vermerkt ihr Ergebnis — je ein Antrag pro
// Ausgang (erfolgreich/gescheitert), damit beide Update-Pfade real
// geprüft sind.
func TestAdministrationRequestAdapterListPendingMarkAppliedMarkFailed(t *testing.T) {
	pool, dsn := newTestAdministrationRequestPool(t)
	ctx := context.Background()

	adapter, err := postgresstorage.NewAdministrationRequest(ctx, dsn)
	if err != nil {
		t.Fatalf("NewAdministrationRequest: %v", err)
	}
	t.Cleanup(adapter.Close)

	var enableID, disableID string
	if err := pool.QueryRow(ctx,
		"SELECT cdc.enable_table($1, $2, $3)", administrationRequestSource, "public", "orders_adapter_enable",
	).Scan(&enableID); err != nil {
		t.Fatalf("cdc.enable_table: %v", err)
	}
	if err := pool.QueryRow(ctx,
		"SELECT cdc.disable_table($1, $2, $3)", administrationRequestSource, "public", "orders_adapter_disable",
	).Scan(&disableID); err != nil {
		t.Fatalf("cdc.disable_table: %v", err)
	}

	pending, err := adapter.ListPending(ctx)
	if err != nil {
		t.Fatalf("ListPending: %v", err)
	}
	byID := map[model.AdministrationRequestID]model.AdministrationRequest{}
	for _, request := range pending {
		byID[request.ID] = request
	}
	enableRequest, found := byID[model.AdministrationRequestID(enableID)]
	if !found {
		t.Fatalf("ListPending trägt nicht den enable-Antrag %q: %+v", enableID, pending)
	}
	if enableRequest.Source != administrationRequestSource || enableRequest.Schema != "public" ||
		enableRequest.Table != "orders_adapter_enable" || enableRequest.Kind != model.AdministrationRequestEnable {
		t.Fatalf("enable-Antrag: %+v", enableRequest)
	}
	disableRequest, found := byID[model.AdministrationRequestID(disableID)]
	if !found {
		t.Fatalf("ListPending trägt nicht den disable-Antrag %q: %+v", disableID, pending)
	}
	if disableRequest.Kind != model.AdministrationRequestDisable {
		t.Fatalf("disable-Antrag: %+v", disableRequest)
	}

	if err := adapter.MarkApplied(ctx, enableRequest.ID); err != nil {
		t.Fatalf("MarkApplied: %v", err)
	}
	if err := adapter.MarkFailed(ctx, disableRequest.ID, "Tabelle existiert nicht an der Quelle"); err != nil {
		t.Fatalf("MarkFailed: %v", err)
	}

	var appliedStatus string
	var failedStatus, failedMessage string
	if err := pool.QueryRow(ctx, "SELECT status FROM cdc.administration_request WHERE administration_request_id = $1", enableID).Scan(&appliedStatus); err != nil {
		t.Fatalf("Status enable-Antrag lesen: %v", err)
	}
	if appliedStatus != "applied" {
		t.Fatalf("Status enable-Antrag = %q, wollen applied", appliedStatus)
	}
	if err := pool.QueryRow(ctx, "SELECT status, error_message FROM cdc.administration_request WHERE administration_request_id = $1", disableID).Scan(&failedStatus, &failedMessage); err != nil {
		t.Fatalf("Status disable-Antrag lesen: %v", err)
	}
	if failedStatus != "failed" || failedMessage != "Tabelle existiert nicht an der Quelle" {
		t.Fatalf("Status/Fehlertext disable-Antrag = %q/%q", failedStatus, failedMessage)
	}

	remaining, err := adapter.ListPending(ctx)
	if err != nil {
		t.Fatalf("ListPending nach Vermerk: %v", err)
	}
	for _, request := range remaining {
		if request.ID == enableRequest.ID || request.ID == disableRequest.ID {
			t.Fatalf("ListPending nach Vermerk trägt weiterhin einen vermerkten Antrag: %+v", request)
		}
	}
}

// TestAdministrationRequestAdapterMarkAppliedIsIdempotent trägt die
// Idempotenz des Vermerks (`ADR-0050`s Konsequenz: ein Antrag über einen
// Prozess-Neustart hinweg wird beim nächsten Poll erneut abgeholt) — ein
// zweiter `MarkApplied`-Aufruf auf einen bereits vermerkten Antrag bleibt
// ohne Fehler und ohne Wirkung.
func TestAdministrationRequestAdapterMarkAppliedIsIdempotent(t *testing.T) {
	pool, dsn := newTestAdministrationRequestPool(t)
	ctx := context.Background()

	adapter, err := postgresstorage.NewAdministrationRequest(ctx, dsn)
	if err != nil {
		t.Fatalf("NewAdministrationRequest: %v", err)
	}
	t.Cleanup(adapter.Close)

	var requestID string
	if err := pool.QueryRow(ctx,
		"SELECT cdc.enable_table($1, $2, $3)", administrationRequestSource, "public", "orders_adapter_idempotent",
	).Scan(&requestID); err != nil {
		t.Fatalf("cdc.enable_table: %v", err)
	}

	if err := adapter.MarkApplied(ctx, model.AdministrationRequestID(requestID)); err != nil {
		t.Fatalf("erstes MarkApplied: %v", err)
	}
	if err := adapter.MarkFailed(ctx, model.AdministrationRequestID(requestID), "zu spät"); err != nil {
		t.Fatalf("MarkFailed nach MarkApplied: %v", err)
	}

	var status string
	if err := pool.QueryRow(ctx, "SELECT status FROM cdc.administration_request WHERE administration_request_id = $1", requestID).Scan(&status); err != nil {
		t.Fatalf("Status lesen: %v", err)
	}
	if status != "applied" {
		t.Fatalf("status = %q, wollen applied (MarkFailed nach MarkApplied bleibt ohne Wirkung)", status)
	}
}

// infoRecorder sammelt die Info-Meldungen eines Adapters.
type infoRecorder struct {
	mu       sync.Mutex
	messages []string
}

func (r *infoRecorder) Debug(context.Context, string, ...any) {}
func (r *infoRecorder) Warn(context.Context, string, ...any)  {}
func (r *infoRecorder) Error(context.Context, string, ...any) {}
func (r *infoRecorder) Info(_ context.Context, msg string, _ ...any) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.messages = append(r.messages, msg)
}

func (r *infoRecorder) count(msg string) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	n := 0
	for _, m := range r.messages {
		if m == msg {
			n++
		}
	}
	return n
}

// TestAdministrationRequestAdapterMarkAppliedLogsOnlyARequestItMarked trägt
// die Bedeutung der Erfolgsmeldung: „Antrag erledigt“ steht nur, wenn
// `MarkApplied` eine `pending`-Zeile vermerkt hat; ein bereits vermerkter
// Antrag (bei `backfill` die Annahme) bleibt ohne Meldung. Rot färbende
// Mutation: die `RowsAffected`-Bedingung in `MarkApplied` entfernen — der
// zweite Aufruf meldet den Antrag erneut.
func TestAdministrationRequestAdapterMarkAppliedLogsOnlyARequestItMarked(t *testing.T) {
	pool, dsn := newTestAdministrationRequestPool(t)
	ctx := context.Background()
	log := &infoRecorder{}

	adapter, err := postgresstorage.NewAdministrationRequest(ctx, dsn, postgresstorage.WithLog(log))
	if err != nil {
		t.Fatalf("NewAdministrationRequest: %v", err)
	}
	t.Cleanup(adapter.Close)

	var requestID string
	if err := pool.QueryRow(ctx,
		"SELECT cdc.enable_table($1, $2, $3)", administrationRequestSource, "public", "orders_adapter_logged",
	).Scan(&requestID); err != nil {
		t.Fatalf("cdc.enable_table: %v", err)
	}
	const done = "administrationrequest: Antrag erledigt"
	if err := adapter.MarkApplied(ctx, model.AdministrationRequestID(requestID)); err != nil {
		t.Fatalf("erstes MarkApplied: %v", err)
	}
	if got := log.count(done); got != 1 {
		t.Fatalf("Meldungen %q nach dem ersten MarkApplied = %d, erwartet 1", done, got)
	}
	if err := adapter.MarkApplied(ctx, model.AdministrationRequestID(requestID)); err != nil {
		t.Fatalf("zweites MarkApplied: %v", err)
	}
	if got := log.count(done); got != 1 {
		t.Fatalf("Meldungen %q nach dem zweiten MarkApplied = %d, erwartet 1 (kein Antrag vermerkt)", done, got)
	}
}

// TestTableActivationExcludedColumnsDerivesAppliedColumnRequests trägt die
// Ableitung des dauerhaften Ausschlussstandes gegen die reale PostgreSQL
// (`LH-FA-CFG-005`, `ADR-0065`): die `applied`-Zeilen der beiden
// Spalten-Antragsarten tragen den Stand, `exclude_column` trägt einen Namen
// ein, `include_column` nimmt ihn wieder heraus. Der Test setzt die
// Antrags-Zeilen direkt (`requested_at` und Antrags-ID sind hier
// Prüfgegenstand); die schreibenden Funktionen und ihren `pending`-Ausgang
// decken die Tests dieses Pakets darüber ab.
//
// Die beiden ersten Zeilen tragen denselben `requested_at`: die eingefügte
// Reihenfolge wäre `b-exclude` vor `a-include` (nicht ausgeschlossen), die
// `administration_request_id` als Zweitschlüssel dreht das um (ausgeschlossen)
// — der Test belegt damit die deterministische Ordnung, nicht nur ihr
// Ergebnis.
func TestTableActivationExcludedColumnsDerivesAppliedColumnRequests(t *testing.T) {
	pool, dsn := newTestAdministrationRequestPool(t)
	ctx := context.Background()

	adapter, err := postgresstorage.NewTableActivation(ctx, dsn)
	if err != nil {
		t.Fatalf("NewTableActivation: %v", err)
	}
	t.Cleanup(adapter.Close)

	const (
		table        = "orders_exclusion_derivation"
		otherTable   = "orders_exclusion_other"
		insertColumn = `INSERT INTO cdc.administration_request
    (administration_request_id, source_id, schema_name, table_name, column_name, request_kind, requested_at, status)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	)
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(),
			"DELETE FROM cdc.administration_request WHERE administration_request_id LIKE 'exclusion-derivation-%'")
	})

	tiedAt := time.Date(2026, 9, 15, 12, 0, 0, 0, time.UTC)
	insert := func(id, targetTable, kind, column, status string, requestedAt time.Time) {
		t.Helper()
		if _, err := pool.Exec(ctx, insertColumn,
			id, administrationRequestSource, "public", targetTable, column, kind, requestedAt, status,
		); err != nil {
			t.Fatalf("Antrags-Zeile %q schreiben: %v", id, err)
		}
	}

	// Reihenfolge-Eins: dieselbe Spalte, derselbe Zeitstempel.
	insert("exclusion-derivation-b-exclude", table, "exclude_column", "secret", "applied", tiedAt)
	insert("exclusion-derivation-a-include", table, "include_column", "secret", "applied", tiedAt)
	// Zeilen ohne Anteil am Stand: fremde Tabelle, andere Antragsart,
	// anderer Ausgang.
	insert("exclusion-derivation-e-other-table", otherTable, "exclude_column", "tenant", "applied", tiedAt)
	insert("exclusion-derivation-f-enable", table, "enable", "", "applied", tiedAt)
	insert("exclusion-derivation-g-pending", table, "exclude_column", "pending_column", "pending", tiedAt)
	insert("exclusion-derivation-h-failed", table, "exclude_column", "failed_column", "failed", tiedAt)

	excluded, err := adapter.ExcludedColumns(ctx, administrationRequestSource)
	if err != nil {
		t.Fatalf("ExcludedColumns: %v", err)
	}
	if got := excluded["public."+table]; len(got) != 1 || got[0] != "secret" {
		t.Fatalf("Ausschlussstand public.%s = %v, wollen [secret] (include_column hebt den früheren exclude_column desselben Zeitstempels auf)", table, got)
	}
	if got := excluded["public."+otherTable]; len(got) != 1 || got[0] != "tenant" {
		t.Fatalf("Ausschlussstand public.%s = %v, wollen [tenant]", otherTable, got)
	}

	// Einschluss nach dem Ausschluss: der Stand fällt weg — kein Eintrag,
	// keine leere Liste.
	insert("exclusion-derivation-c-include", table, "include_column", "secret", "applied", tiedAt.Add(time.Second))
	excluded, err = adapter.ExcludedColumns(ctx, administrationRequestSource)
	if err != nil {
		t.Fatalf("ExcludedColumns nach dem Einschluss: %v", err)
	}
	if got, present := excluded["public."+table]; present {
		t.Fatalf("Ausschlussstand public.%s = %v, wollen keinen Eintrag (include_column hat den letzten Namen genommen)", table, got)
	}

	// Eine Quelle ohne Antrag trägt keinen Stand.
	excluded, err = adapter.ExcludedColumns(ctx, model.SourceID("src-administration-ohne-antraege"))
	if err != nil {
		t.Fatalf("ExcludedColumns ohne Anträge: %v", err)
	}
	if len(excluded) != 0 {
		t.Fatalf("Ausschlussstand einer Quelle ohne Anträge = %v, wollen leer", excluded)
	}
}

// TestAdministrationListenerWaitForNotification trägt die
// `LISTEN`-Wecksignal-Fähigkeit real: ein `NOTIFY` löst die Rückkehr aus,
// eine ausbleibende Benachrichtigung endet über den Timeout des
// übergebenen `ctx` (`context.DeadlineExceeded`) — der Fallback-Poll-Takt
// der Administrations-Goroutine.
func TestAdministrationListenerWaitForNotification(t *testing.T) {
	pool, dsn := newTestAdministrationRequestPool(t)
	ctx := context.Background()

	listener, err := postgresstorage.NewAdministrationListener(ctx, dsn)
	if err != nil {
		t.Fatalf("NewAdministrationListener: %v", err)
	}
	t.Cleanup(func() { _ = listener.Close(context.Background()) })

	timeoutCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()
	if err := listener.WaitForNotification(timeoutCtx); !stderrors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("WaitForNotification ohne NOTIFY: %v, wollen context.DeadlineExceeded", err)
	}

	var requestID string
	if err := pool.QueryRow(ctx,
		"SELECT cdc.enable_table($1, $2, $3)", administrationRequestSource, "public", "orders_listener",
	).Scan(&requestID); err != nil {
		t.Fatalf("cdc.enable_table: %v", err)
	}

	notifyCtx, cancelNotify := context.WithTimeout(ctx, 5*time.Second)
	defer cancelNotify()
	if err := listener.WaitForNotification(notifyCtx); err != nil {
		t.Fatalf("WaitForNotification nach NOTIFY: %v", err)
	}
}

// TestAdministrationRequestBackfillTableWritesOnlyTheRequestAndNotifies trägt
// den Antrags-Weg der Bestands-Antragsart (`LH-FA-CAP-009`, `ADR-0111`
// Teilfrage 5, `ADR-0050`): `cdc.backfill_table` legt genau eine Zeile der Art
// `backfill` ohne Spalte mit Status `pending` an und sendet `pg_notify` mit
// der Antrags-Kennung; sie legt weder eine Run-Zeile an noch berührt sie
// Bindung oder Publication. Rot färbende Mutation: den `INSERT` der Funktion
// um `INSERT INTO cdc.backfill_run …` erweitern — die Prüfung „keine
// Run-Zeile“ meldet sie; die Art `'backfill'` im `INSERT` durch `'enable'`
// ersetzen — die Prüfung der Art meldet sie.
func TestAdministrationRequestBackfillTableWritesOnlyTheRequestAndNotifies(t *testing.T) {
	pool, dsn := newTestAdministrationRequestPool(t)
	ctx := context.Background()

	listener := listenForAdministrationNotify(t, ctx, dsn)

	var requestID string
	if err := pool.QueryRow(ctx,
		"SELECT cdc.backfill_table($1, $2, $3)", administrationRequestSource, "public", "orders_backfill",
	).Scan(&requestID); err != nil {
		t.Fatalf("cdc.backfill_table: %v", err)
	}
	if requestID == "" {
		t.Fatalf("cdc.backfill_table lieferte eine leere Antrags-ID")
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM cdc.administration_request WHERE administration_request_id = $1", requestID)
	})

	notifyCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	notification, err := listener.WaitForNotification(notifyCtx)
	if err != nil {
		t.Fatalf("WaitForNotification: %v", err)
	}
	if notification.Channel != "cdc_administration" || notification.Payload != requestID {
		t.Fatalf("Notify = %q/%q, wollen cdc_administration/%q", notification.Channel, notification.Payload, requestID)
	}

	var sourceID, schemaName, tableName, requestKind, status string
	var column *string
	if err := pool.QueryRow(ctx,
		"SELECT source_id, schema_name, table_name, column_name, request_kind, status FROM cdc.administration_request WHERE administration_request_id = $1",
		requestID,
	).Scan(&sourceID, &schemaName, &tableName, &column, &requestKind, &status); err != nil {
		t.Fatalf("Antrags-Zeile lesen: %v", err)
	}
	if sourceID != administrationRequestSource || schemaName != "public" || tableName != "orders_backfill" {
		t.Fatalf("Antrags-Zeile: source=%q schema=%q table=%q", sourceID, schemaName, tableName)
	}
	if requestKind != "backfill" || status != "pending" || column != nil {
		t.Fatalf("Antrags-Zeile: Art %q, Status %q, Spalte %v — erwartet backfill, pending, NULL", requestKind, status, column)
	}
	var runs int
	if err := pool.QueryRow(ctx, "SELECT count(*) FROM cdc.backfill_run WHERE run_id = $1", requestID).Scan(&runs); err != nil {
		t.Fatalf("Run-Zeilen zählen: %v", err)
	}
	if runs != 0 {
		t.Fatalf("cdc.backfill_table legte %d Run-Zeile(n) an — die Funktion schreibt nur den Antrag (ADR-0050)", runs)
	}
}

// TestAdministrationRequestBackfillTableRequiresCdcAdminMembership belegt das
// `REVOKE … FROM PUBLIC`/`GRANT … TO cdc_admin`-Paar für `cdc.backfill_table`
// (`ADR-0047`, `LH-QA-SEC-001`…`003`): eine Rolle ohne `cdc_admin`-Mitgliedschaft
// scheitert mit SQLSTATE 42501 („permission denied for function“). Rot färbende
// Mutation: `cdc.backfill_table(text, text, text)` aus der `REVOKE`-Zeile der
// Nacharbeit-Datei streichen — `PUBLIC` behält `EXECUTE`, der Aufruf gelingt.
func TestAdministrationRequestBackfillTableRequiresCdcAdminMembership(t *testing.T) {
	pool, _ := newTestAdministrationRequestPool(t)
	ctx := context.Background()

	for _, role := range []string{"cdc_reader", "cdc_capture"} {
		conn, err := pool.Acquire(ctx)
		if err != nil {
			t.Fatalf("Verbindung reservieren: %v", err)
		}
		func() {
			defer func() {
				_, _ = conn.Exec(ctx, "RESET ROLE")
				conn.Release()
			}()
			if _, err := conn.Exec(ctx, "SET ROLE "+role); err != nil {
				t.Fatalf("SET ROLE %s: %v", role, err)
			}
			var requestID string
			err := conn.QueryRow(ctx,
				"SELECT cdc.backfill_table($1, $2, $3)", administrationRequestSource, "public", "orders_backfill_denied",
			).Scan(&requestID)
			if !permissionDenied(err) {
				t.Fatalf("%s SELECT cdc.backfill_table(...): erwartet SQLSTATE 42501 (insufficient_privilege), erhalten %v", role, err)
			}
		}()
	}
}

// TestAdministrationRequestKindCheckCarriesExactlyTheSevenKinds belegt die
// geschlossene `request_kind`-Menge (`chk_administration_request_kind`,
// `tools/schema/nacharbeit-administration.sql`): jede der sieben Arten wird
// angenommen, eine achte endet mit SQLSTATE 23514. Rot färbende Mutationen
// (je eine): `'set_transformation'` aus der CHECK-Klausel streichen — die Art
// `set_transformation` endet mit 23514; `'truncate'` in die Klausel
// aufnehmen — die achte Art wird angenommen.
func TestAdministrationRequestKindCheckCarriesExactlyTheSevenKinds(t *testing.T) {
	pool, _ := newTestAdministrationRequestPool(t)
	ctx := context.Background()
	const requestID = "kind-check-request"
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM cdc.administration_request WHERE administration_request_id = $1", requestID)
	})

	insert := func(kind string) error {
		_, err := pool.Exec(ctx,
			`INSERT INTO cdc.administration_request (administration_request_id, source_id, schema_name, table_name, request_kind, status)
			 VALUES ($1, $2, 'public', 'kind_check', $3, 'pending')`, requestID, administrationRequestSource, kind)
		return err
	}
	for _, kind := range []string{"enable", "disable", "exclude_column", "include_column", "backfill", "set_transformation", "remove_transformation"} {
		if err := insert(kind); err != nil {
			t.Fatalf("Art %q: erwartet angenommen, %v", kind, err)
		}
		if _, err := pool.Exec(ctx, "DELETE FROM cdc.administration_request WHERE administration_request_id = $1", requestID); err != nil {
			t.Fatalf("Aufräumen: %v", err)
		}
	}
	err := insert("truncate")
	var pgErr *pgconn.PgError
	if !stderrors.As(err, &pgErr) || pgErr.Code != "23514" {
		t.Fatalf("Art %q: erwartet SQLSTATE 23514 (check_violation), erhalten %v", "truncate", err)
	}
}

// TestAdministrationRequestTransformationRequestsCarryRuleAndNotify trägt den
// realen Antrags-Weg der Transformations-Antragsarten (`LH-FA-CFG-007`,
// `ADR-0112` Teilfrage 1): `cdc.set_transformation` schreibt eine `pending`-Zeile
// mit Regelname und Regelform (Spalte `jsonb`) und ohne Spalte, `cdc.remove_transformation`
// eine mit Regelname und ohne Regelform; beide senden `pg_notify` mit der
// Antrags-ID, und der Adapter liest beide Felder über `ListPending` zurück.
// Rot färbende Mutationen (Eingabeseite): die Art `'set_transformation'` in
// `cdc.set_transformation` durch `'remove_transformation'` ersetzen — die Zeile
// trägt die falsche Art; den Ausdruck `COALESCE(rule_spec::text, …)` in
// `SelectPendingAdministrationRequests` durch die leere Zeichenkette ersetzen —
// `ListPending` endet mit der Konstruktor-Invariante `ErrEmptyIdentifier`.
func TestAdministrationRequestTransformationRequestsCarryRuleAndNotify(t *testing.T) {
	pool, dsn := newTestAdministrationRequestPool(t)
	ctx := context.Background()

	var requestIDs []string
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM cdc.administration_request WHERE administration_request_id = ANY($1)", requestIDs)
	})
	adapter, err := postgresstorage.NewAdministrationRequest(ctx, dsn)
	if err != nil {
		t.Fatalf("NewAdministrationRequest: %v", err)
	}
	t.Cleanup(adapter.Close)
	listener := listenForAdministrationNotify(t, ctx, dsn)

	const ruleSpec = `{"kind": "rename_column", "column": "name", "to": "title"}`
	var setID, removeID, nullSpecID string
	if err := pool.QueryRow(ctx,
		"SELECT cdc.set_transformation($1, $2, $3, $4, $5::json)", administrationRequestSource, "public", "orders_rules", "umbenennung", ruleSpec,
	).Scan(&setID); err != nil {
		t.Fatalf("cdc.set_transformation: %v", err)
	}
	requestIDs = append(requestIDs, setID)
	notifyCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	notification, err := listener.WaitForNotification(notifyCtx)
	if err != nil {
		t.Fatalf("WaitForNotification: %v", err)
	}
	if notification.Channel != "cdc_administration" || notification.Payload != setID {
		t.Fatalf("Notify = %q/%q, wollen cdc_administration/%q", notification.Channel, notification.Payload, setID)
	}
	if err := pool.QueryRow(ctx,
		"SELECT cdc.remove_transformation($1, $2, $3, $4)", administrationRequestSource, "public", "orders_rules", "umbenennung",
	).Scan(&removeID); err != nil {
		t.Fatalf("cdc.remove_transformation: %v", err)
	}
	requestIDs = append(requestIDs, removeID)
	notification, err = listener.WaitForNotification(notifyCtx)
	if err != nil {
		t.Fatalf("WaitForNotification: %v", err)
	}
	if notification.Payload != removeID {
		t.Fatalf("Notify-Payload = %q, wollen die Antrags-ID %q", notification.Payload, removeID)
	}
	// Ein JSON-`null` ist eine Regelform, kein SQL-NULL: die Spalte trägt den
	// Text `null`.
	if err := pool.QueryRow(ctx,
		"SELECT cdc.set_transformation($1, $2, $3, $4, 'null'::json)", administrationRequestSource, "public", "orders_rules", "json_null",
	).Scan(&nullSpecID); err != nil {
		t.Fatalf("cdc.set_transformation(JSON null): %v", err)
	}
	requestIDs = append(requestIDs, nullSpecID)

	for _, tc := range []struct {
		id, kind string
		wantSpec bool
	}{{setID, "set_transformation", true}, {removeID, "remove_transformation", false}} {
		var kind, status, ruleName string
		var column, spec, message *string
		if err := pool.QueryRow(ctx,
			"SELECT request_kind, status, rule_name, column_name, rule_spec::text, error_message FROM cdc.administration_request WHERE administration_request_id = $1", tc.id,
		).Scan(&kind, &status, &ruleName, &column, &spec, &message); err != nil {
			t.Fatalf("Antrags-Zeile %s lesen: %v", tc.kind, err)
		}
		if kind != tc.kind || status != "pending" || ruleName != "umbenennung" || column != nil || message != nil || (spec != nil) != tc.wantSpec {
			t.Fatalf("Antrags-Zeile %s: Art %q, Status %q, Regelname %q, Spalte gesetzt %t, Regelform gesetzt %t (erwartet %t), Fehlertext gesetzt %t",
				tc.kind, kind, status, ruleName, column != nil, spec != nil, tc.wantSpec, message != nil)
		}
	}

	pending, err := adapter.ListPending(ctx)
	if err != nil {
		t.Fatalf("ListPending: %v", err)
	}
	byID := map[model.AdministrationRequestID]model.AdministrationRequest{}
	for _, request := range pending {
		byID[request.ID] = request
	}
	setRequest, found := byID[model.AdministrationRequestID(setID)]
	if !found {
		t.Fatalf("ListPending trägt nicht den set_transformation-Antrag %q: %+v", setID, pending)
	}
	if setRequest.Kind != model.AdministrationRequestSetTransformation || setRequest.RuleName != "umbenennung" || setRequest.Column != "" || setRequest.Table != "orders_rules" {
		t.Fatalf("set_transformation-Antrag: %+v", setRequest)
	}
	var decoded map[string]string
	if err := json.Unmarshal([]byte(setRequest.RuleSpec), &decoded); err != nil {
		t.Fatalf("Regelform %q ist kein JSON-Objekt: %v", setRequest.RuleSpec, err)
	}
	if decoded["kind"] != "rename_column" || decoded["column"] != "name" || decoded["to"] != "title" || len(decoded) != 3 {
		t.Fatalf("Regelform = %v, wollen kind/column/to wie beantragt", decoded)
	}
	removeRequest, found := byID[model.AdministrationRequestID(removeID)]
	if !found {
		t.Fatalf("ListPending trägt nicht den remove_transformation-Antrag %q: %+v", removeID, pending)
	}
	if removeRequest.Kind != model.AdministrationRequestRemoveTransformation || removeRequest.RuleName != "umbenennung" || removeRequest.RuleSpec != "" {
		t.Fatalf("remove_transformation-Antrag: %+v", removeRequest)
	}
	if nullRequest := byID[model.AdministrationRequestID(nullSpecID)]; nullRequest.RuleSpec != "null" {
		t.Fatalf("Regelform des JSON-null-Antrags = %q, wollen den Text null", nullRequest.RuleSpec)
	}
}

// TestAdministrationRequestTransformationFunctionsRequireCdcAdminMembership
// belegt das `REVOKE … FROM PUBLIC`/`GRANT … TO cdc_admin`-Paar für
// `cdc.set_transformation` und `cdc.remove_transformation` (`ADR-0047`,
// `LH-QA-SEC-001`…`003`): eine Rolle ohne `cdc_admin`-Mitgliedschaft scheitert
// mit SQLSTATE 42501 („permission denied for function“), und der gescheiterte
// Aufruf hinterlässt keine Zeile. Rot färbende Mutation: die jeweilige
// Signatur aus der `REVOKE`-Zeile der Nacharbeit-Datei streichen — `PUBLIC`
// behält `EXECUTE`, der Aufruf gelingt.
func TestAdministrationRequestTransformationFunctionsRequireCdcAdminMembership(t *testing.T) {
	pool, _ := newTestAdministrationRequestPool(t)
	ctx := context.Background()
	const table = "orders_rules_denied"
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM cdc.administration_request WHERE table_name = $1", table)
	})

	for _, role := range []string{"cdc_reader", "cdc_capture"} {
		conn, err := pool.Acquire(ctx)
		if err != nil {
			t.Fatalf("Verbindung reservieren: %v", err)
		}
		func() {
			defer func() {
				_, _ = conn.Exec(ctx, "RESET ROLE")
				conn.Release()
			}()
			if _, err := conn.Exec(ctx, "SET ROLE "+role); err != nil {
				t.Fatalf("SET ROLE %s: %v", role, err)
			}
			var requestID string
			err := conn.QueryRow(ctx,
				"SELECT cdc.set_transformation($1, $2, $3, $4, $5::json)", administrationRequestSource, "public", table, "verboten", `{"kind": "rename_column"}`,
			).Scan(&requestID)
			if !permissionDenied(err) {
				t.Fatalf("%s SELECT cdc.set_transformation(...): erwartet SQLSTATE 42501 (insufficient_privilege), erhalten %v", role, err)
			}
			err = conn.QueryRow(ctx,
				"SELECT cdc.remove_transformation($1, $2, $3, $4)", administrationRequestSource, "public", table, "verboten",
			).Scan(&requestID)
			if !permissionDenied(err) {
				t.Fatalf("%s SELECT cdc.remove_transformation(...): erwartet SQLSTATE 42501 (insufficient_privilege), erhalten %v", role, err)
			}
		}()
	}
	var rows int
	if err := pool.QueryRow(ctx, "SELECT count(*) FROM cdc.administration_request WHERE table_name = $1", table).Scan(&rows); err != nil {
		t.Fatalf("Zeilen zählen: %v", err)
	}
	if rows != 0 {
		t.Fatalf("die abgelehnten Aufrufe hinterließen %d Antragszeile(n)", rows)
	}
}

// TestAdministrationRequestRuleColumnsAreNullable trägt die Spaltenform der
// zwei Regel-Spalten (`SPEC-019`): `rule_name` ist `text`, `rule_spec` ist
// `jsonb`, beide nullable — die fünf Antragsarten ohne Regel lassen sie NULL.
// Rot färbende Mutation: `rule_spec` in `tools/schema/schema.yaml` mit
// `type: text` statt `type: json` deklarieren — der Test meldet
// `rule_spec:text:YES`.
func TestAdministrationRequestRuleColumnsAreNullable(t *testing.T) {
	pool, _ := newTestAdministrationRequestPool(t)
	ctx := context.Background()

	var shape string
	if err := pool.QueryRow(ctx,
		`SELECT string_agg(column_name || ':' || data_type || ':' || is_nullable, ',' ORDER BY column_name)
		 FROM information_schema.columns
		 WHERE table_schema = 'cdc' AND table_name = 'administration_request' AND column_name IN ('rule_name', 'rule_spec')`,
	).Scan(&shape); err != nil {
		t.Fatalf("Spaltenform lesen: %v", err)
	}
	if shape != "rule_name:text:YES,rule_spec:jsonb:YES" {
		t.Fatalf("Spaltenform = %q, wollen rule_name:text:YES,rule_spec:jsonb:YES", shape)
	}
}

// TestAdministrationRequestSetTransformationAcceptanceSet trägt die
// Annahmemenge des Parameters `rule_spec` von `cdc.set_transformation`
// (`ADR-0126` Festlegung 1, `SPEC-019`): jede Form läuft als Text, der über
// `::text::json` zum Parameter wird — der Weg eines Clients mit Literal. Eine
// angenommene Form schreibt genau eine `pending`-Zeile und liest ihren Wert
// als `rule_spec::text` zurück (`jsonb` normalisiert: der letzte Wert eines
// doppelten Schlüssels gilt); eine abgelehnte Form endet mit einem Fehler des
// Aufrufs und hinterlässt keine Zeile. Die Ablehnung von `\u0000` und der
// Zahl außerhalb des Bereichs entsteht beim Cast `p_rule_spec::jsonb` in der
// Funktion. Rot färbende Mutationen (Eingabeseite): den Ausdruck
// `p_rule_spec::jsonb` in `tools/schema/nacharbeit-administration.sql` durch
// `NULL::jsonb` ersetzen — jede abgelehnte Form wird angenommen, jede
// angenommene Form mit Wert liest NULL.
func TestAdministrationRequestSetTransformationAcceptanceSet(t *testing.T) {
	pool, _ := newTestAdministrationRequestPool(t)
	ctx := context.Background()
	const table = "rule_spec_acceptance"
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM cdc.administration_request WHERE table_name = $1", table)
	})

	call := func(ruleName string, spec *string) error {
		var requestID string
		return pool.QueryRow(ctx,
			"SELECT cdc.set_transformation($1, $2, $3, $4, $5::text::json)", administrationRequestSource, "public", table, ruleName, spec,
		).Scan(&requestID)
	}
	stored := func(ruleName string) (rows int, spec *string) {
		if err := pool.QueryRow(ctx,
			"SELECT count(*), min(rule_spec::text) FROM cdc.administration_request WHERE table_name = $1 AND rule_name = $2", table, ruleName,
		).Scan(&rows, &spec); err != nil {
			t.Fatalf("Zeilen von %q lesen: %v", ruleName, err)
		}
		return rows, spec
	}
	text := func(s string) *string { return &s }

	accepted := []struct {
		name string
		spec *string
		want *string
	}{
		{"sql_null", nil, nil},
		{"json_null", text(`null`), text(`null`)},
		{"leeres_objekt", text(`{}`), text(`{}`)},
		{"leeres_array", text(`[]`), text(`[]`)},
		{"zahl", text(`1`), text(`1`)},
		{"zeichenkette", text(`"x"`), text(`"x"`)},
		{"doppelter_schluessel", text(`{"kind":"a","kind":"b"}`), text(`{"kind": "b"}`)},
	}
	for _, tc := range accepted {
		if err := call(tc.name, tc.spec); err != nil {
			t.Errorf("Form %s: erwartet angenommen, erhalten %v", tc.name, err)
			continue
		}
		rows, spec := stored(tc.name)
		if rows != 1 || (spec == nil) != (tc.want == nil) || (spec != nil && *spec != *tc.want) {
			t.Errorf("Form %s: %d Zeile(n), rule_spec gelesen %s, erwartet 1 Zeile mit %s", tc.name, rows, describeSpec(spec), describeSpec(tc.want))
		}
	}

	rejected := []struct {
		name string
		spec string
	}{
		{"nul_im_wert", `{"a":"\u0000"}`},
		{"nul_im_schluessel", `{"\u0000":1}`},
		{"nul_im_array", `["\u0000"]`},
		{"nul_als_skalar", `"\u0000"`},
		{"zahl_ausserhalb_des_bereichs", `{"a":1e200000}`},
		{"syntaxfehler", `{oops`},
		{"leerer_text", ``},
		{"einzelnes_surrogat", `"\ud83d"`},
	}
	for _, tc := range rejected {
		if err := call(tc.name, &tc.spec); err == nil {
			t.Errorf("Form %s: erwartet abgelehnt, der Aufruf gelang", tc.name)
		}
		if rows, _ := stored(tc.name); rows != 0 {
			t.Errorf("Form %s: der abgelehnte Aufruf hinterließ %d Zeile(n)", tc.name, rows)
		}
	}
}

// TestAdministrationFunctionsPinSecurityDefinerAndSearchPath bindet die Zusage
// „`SECURITY DEFINER` mit gepinntem `search_path`“ an den Katalog der realen
// Instanz: jede der sieben schreibenden Funktionen trägt `prosecdef` und
// `proconfig = {search_path=cdc, pg_temp}` (`ADR-0050`, `LH-QA-SEC-002`).
// Rot färbende Mutation (Eingabeseite): die Zeile `SET search_path = cdc,
// pg_temp` einer Funktion in `tools/schema/nacharbeit-administration.sql`
// streichen — die Funktion erscheint in der Meldung.
func TestAdministrationFunctionsPinSecurityDefinerAndSearchPath(t *testing.T) {
	pool, _ := newTestAdministrationRequestPool(t)

	var unpinned []string
	if err := pool.QueryRow(context.Background(),
		`SELECT coalesce(array_agg(p.proname::text ORDER BY p.proname), '{}')
		 FROM pg_proc p JOIN pg_namespace n ON n.oid = p.pronamespace
		 WHERE n.nspname = 'cdc'
		   AND p.proname IN ('enable_table', 'disable_table', 'exclude_column', 'include_column', 'backfill_table', 'set_transformation', 'remove_transformation')
		   AND (NOT p.prosecdef OR p.proconfig IS DISTINCT FROM ARRAY['search_path=cdc, pg_temp'])`,
	).Scan(&unpinned); err != nil {
		t.Fatalf("Katalog lesen: %v", err)
	}
	if len(unpinned) != 0 {
		t.Fatalf("ohne SECURITY DEFINER und gepinnten search_path (cdc, pg_temp): %v", unpinned)
	}
}

// describeSpec zeigt eine nullable Regelform lesbar: SQL-NULL als `NULL`, alles
// andere als der gelesene Text.
func describeSpec(spec *string) string {
	if spec == nil {
		return "NULL"
	}
	return *spec
}

// TestTableActivationTransformationRulesDeriveAppliedRuleRequests trägt die
// Ableitung des dauerhaften Regelstandes gegen die reale PostgreSQL
// (`LH-FA-CFG-007`, `SPEC-019`): die `applied`-Zeilen der beiden
// Transformations-Antragsarten tragen den Stand, `set_transformation` trägt
// eine Regel ein, `remove_transformation` nimmt sie unter dem Namen heraus.
// Der Test setzt die Antrags-Zeilen direkt (`requested_at` und Antrags-ID sind
// Prüfgegenstand); Bereinigung und Zeilen sind auf die Kennungs-Vorsilbe
// `rules-derivation-` begrenzt.
//
// Die beiden ersten Zeilen tragen denselben `requested_at`: die eingefügte
// Reihenfolge wäre `b-set` vor `a-remove` (Regel herausgenommen), die
// `administration_request_id` als Zweitschlüssel dreht das um (Regel geführt).
// Rot färbende Mutation: `administration_request_id` aus dem `ORDER BY` von
// `SelectAppliedTransformationRequests` streichen — die Reihenfolge der
// gleichzeitigen Zeilen folgt der Einfügung, die Regel `tie` fehlt.
func TestTableActivationTransformationRulesDeriveAppliedRuleRequests(t *testing.T) {
	pool, dsn := newTestAdministrationRequestPool(t)
	ctx := context.Background()

	adapter, err := postgresstorage.NewTableActivation(ctx, dsn)
	if err != nil {
		t.Fatalf("NewTableActivation: %v", err)
	}
	t.Cleanup(adapter.Close)
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(),
			"DELETE FROM cdc.administration_request WHERE administration_request_id LIKE 'rules-derivation-%'")
	})

	const (
		table       = "orders_rules_derivation"
		otherTable  = "orders_rules_other"
		insertRule  = `INSERT INTO cdc.administration_request
    (administration_request_id, source_id, schema_name, table_name, column_name, rule_name, rule_spec, request_kind, requested_at, status)
VALUES ($1, $2, $3, $4, $5, $6, $7::text::jsonb, $8, $9, $10)`
		renameSpec = `{"kind": "rename_column", "column": "%s", "to": "%s"}`
	)
	tiedAt := time.Date(2026, 9, 26, 12, 0, 0, 0, time.UTC)
	insert := func(id, targetTable, kind, ruleName string, spec *string, status string, requestedAt time.Time) {
		t.Helper()
		var column, name *string
		if kind == "exclude_column" {
			secret := "secret"
			column = &secret
		}
		if ruleName != "" {
			name = &ruleName
		}
		if _, err := pool.Exec(ctx, insertRule,
			id, administrationRequestSource, "public", targetTable, column, name, spec, kind, requestedAt, status,
		); err != nil {
			t.Fatalf("Antrags-Zeile %q schreiben: %v", id, err)
		}
	}
	spec := func(column, to string) *string {
		text := fmt.Sprintf(renameSpec, column, to)
		return &text
	}

	// Die Zeilen stehen in aufsteigender Ordnung von `requested_at` in der
	// Tabelle: ohne Zweitschlüssel läse die Abfrage sie in Einfüge-Ordnung,
	// und der Test färbte sich an der Ordnung der Gleichzeitigen rot.
	// Gleicher Zeitstempel, zwei Paare: Set eingefügt vor Remove, das Remove
	// ordnet nach der Antrags-ID zuerst (`a-…` vor `b-…`), die Regel bleibt.
	insert("rules-derivation-b-set", table, "set_transformation", "tie", spec("name", "customer_name"), "applied", tiedAt)
	insert("rules-derivation-a-remove", table, "remove_transformation", "tie", nil, "applied", tiedAt)
	insert("rules-derivation-d-set", table, "set_transformation", "tie_two", spec("note", "notiz"), "applied", tiedAt)
	insert("rules-derivation-c-remove", table, "remove_transformation", "tie_two", nil, "applied", tiedAt)
	// Eine Regel, die wieder herausgenommen ist, lässt keinen Eintrag zurück.
	insert("rules-derivation-f-set", otherTable, "set_transformation", "gone", spec("name", "x"), "applied", tiedAt)
	// Zyklus über drei Zeitpunkte: Set, Remove, Set mit anderem Ziel.
	insert("rules-derivation-g-remove", otherTable, "remove_transformation", "gone", nil, "applied", tiedAt.Add(time.Second))
	insert("rules-derivation-h-set", table, "set_transformation", "cycle", spec("status", "state_a"), "applied", tiedAt.Add(time.Second))
	insert("rules-derivation-i-remove", table, "remove_transformation", "cycle", nil, "applied", tiedAt.Add(2*time.Second))
	insert("rules-derivation-j-set", table, "set_transformation", "cycle", spec("status", "state_b"), "applied", tiedAt.Add(3*time.Second))
	// Zeilen ohne Anteil am Stand: anderer Ausgang, andere Antragsart.
	insert("rules-derivation-h-pending", table, "set_transformation", "pending_rule", spec("note", "n1"), "pending", tiedAt)
	insert("rules-derivation-i-failed", table, "set_transformation", "failed_rule", spec("note", "n2"), "failed", tiedAt)
	insert("rules-derivation-j-exclude", table, "exclude_column", "", nil, "applied", tiedAt)

	state, err := adapter.TransformationRules(ctx, administrationRequestSource)
	if err != nil {
		t.Fatalf("TransformationRules: %v", err)
	}
	describe := func(rules []model.Transformation) string {
		out := ""
		for _, rule := range rules {
			out += fmt.Sprintf("%s:%s>%s;", rule.Name(), rule.Column(), rule.To())
		}
		return out
	}
	if got, want := describe(state["public."+table]), "tie:name>customer_name;tie_two:note>notiz;cycle:status>state_b;"; got != want {
		t.Fatalf("Regelstand public.%s = %q, wollen %q (Zweitschlüssel, Zyklus, ohne pending/failed/fremde Art)", table, got, want)
	}
	if got, present := state["public."+otherTable]; present {
		t.Fatalf("Regelstand public.%s = %v, wollen keinen Eintrag (remove_transformation hat die letzte Regel genommen)", otherTable, got)
	}

	// Ein Remove nach dem letzten Set lässt die Tabelle ohne Eintrag.
	insert("rules-derivation-k-remove", table, "remove_transformation", "tie", nil, "applied", tiedAt.Add(4*time.Second))
	insert("rules-derivation-l-remove", table, "remove_transformation", "cycle", nil, "applied", tiedAt.Add(5*time.Second))
	insert("rules-derivation-m-remove", table, "remove_transformation", "tie_two", nil, "applied", tiedAt.Add(6*time.Second))
	state, err = adapter.TransformationRules(ctx, administrationRequestSource)
	if err != nil {
		t.Fatalf("TransformationRules nach dem Herausnehmen: %v", err)
	}
	if got, present := state["public."+table]; present {
		t.Fatalf("Regelstand public.%s = %v, wollen keinen Eintrag", table, got)
	}

	// Eine Quelle ohne Antrag trägt keinen Stand.
	state, err = adapter.TransformationRules(ctx, model.SourceID("src-administration-ohne-antraege"))
	if err != nil || len(state) != 0 {
		t.Fatalf("Regelstand einer Quelle ohne Anträge = %v (%v), wollen leer", state, err)
	}
}

// TestTableActivationTransformationRulesFailVisiblyOnUnparsableRow trägt: eine
// vermerkte Zeile, deren Regelform nicht mehr zu einer Regel führt, endet
// sichtbar, statt den Stand um sie zu verkürzen.
func TestTableActivationTransformationRulesFailVisiblyOnUnparsableRow(t *testing.T) {
	pool, dsn := newTestAdministrationRequestPool(t)
	ctx := context.Background()
	adapter, err := postgresstorage.NewTableActivation(ctx, dsn)
	if err != nil {
		t.Fatalf("NewTableActivation: %v", err)
	}
	t.Cleanup(adapter.Close)
	const source = "src-administration-unparsable"
	if _, err := pool.Exec(ctx, "INSERT INTO cdc.source (source_id, name) VALUES ($1, 'unparsable') ON CONFLICT (source_id) DO NOTHING", source); err != nil {
		t.Fatalf("Quelle-Zeile: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM cdc.administration_request WHERE source_id = $1", source)
		_, _ = pool.Exec(context.Background(), "DELETE FROM cdc.source WHERE source_id = $1", source)
	})
	if _, err := pool.Exec(ctx, `INSERT INTO cdc.administration_request
    (administration_request_id, source_id, schema_name, table_name, rule_name, rule_spec, request_kind, status)
VALUES ('rules-unparsable-1', $1, 'public', 'orders_unparsable', 'regel', '{"kind": "unbekannt"}'::jsonb, 'set_transformation', 'applied')`, source); err != nil {
		t.Fatalf("Antrags-Zeile: %v", err)
	}
	if _, err := adapter.TransformationRules(ctx, source); !stderrors.Is(err, domainerrors.ErrUnknownTransformationKind) {
		t.Fatalf("TransformationRules = %v, wollen ErrUnknownTransformationKind", err)
	}
}

// TestTableActivationSourceColumnsReadsTheCatalog trägt die Katalog-Lesart der
// Spaltenliste gegen die reale PostgreSQL (`LH-FA-CFG-007`, K3/K4): die
// Spaltennamen der Tabelle in ihrer Reihenfolge, zeichengenau (auch ein Name
// mit Großbuchstaben), einschließlich einer später hinzugefügten Spalte; eine
// nicht vorhandene Tabelle liefert die leere Liste, ein Bezeichner außerhalb
// des Alphabets endet in der Fehlerklasse `configuration`. Rot färbende
// Mutation: `ORDER BY ordinal_position` streichen ist unsichtbar (die Liste
// bleibt eine Menge); `table_name = $2` gegen `table_name = $1` tauschen — die
// Liste ist leer.
func TestTableActivationSourceColumnsReadsTheCatalog(t *testing.T) {
	pool, dsn := newTestAdministrationRequestPool(t)
	ctx := context.Background()
	adapter, err := postgresstorage.NewTableActivation(ctx, dsn)
	if err != nil {
		t.Fatalf("NewTableActivation: %v", err)
	}
	t.Cleanup(adapter.Close)

	const table = "orders_rule_columns"
	if _, err := pool.Exec(ctx, `CREATE TABLE public.`+table+` (id integer PRIMARY KEY, "Name" text, status text)`); err != nil {
		t.Fatalf("Tabelle anlegen: %v", err)
	}
	t.Cleanup(func() { _, _ = pool.Exec(context.Background(), "DROP TABLE IF EXISTS public."+table) })

	columns, err := adapter.SourceColumns(ctx, "public", table)
	if err != nil {
		t.Fatalf("SourceColumns: %v", err)
	}
	if fmt.Sprint(columns) != "[id Name status]" {
		t.Fatalf("Spalten = %v, wollen [id Name status] in Tabellen-Reihenfolge, zeichengenau", columns)
	}
	if _, err := pool.Exec(ctx, "ALTER TABLE public."+table+" ADD COLUMN extra text"); err != nil {
		t.Fatalf("Spalte hinzufügen: %v", err)
	}
	columns, err = adapter.SourceColumns(ctx, "public", table)
	if err != nil || fmt.Sprint(columns) != "[id Name status extra]" {
		t.Fatalf("Spalten nach ADD COLUMN = %v (%v), wollen [id Name status extra]", columns, err)
	}

	missing, err := adapter.SourceColumns(ctx, "public", "orders_rule_columns_fehlt")
	if err != nil || missing == nil || len(missing) != 0 {
		t.Fatalf("Spalten einer fehlenden Tabelle = %v (%v), wollen die leere Liste", missing, err)
	}
	if _, err := adapter.SourceColumns(ctx, "public; DROP", table); !stderrors.Is(err, postgresstorage.ErrActivationConfiguration) {
		t.Fatalf("Schema außerhalb des Alphabets: Fehler = %v, wollen ErrActivationConfiguration", err)
	}
	if _, err := adapter.SourceColumns(ctx, "public", "Orders"); !stderrors.Is(err, postgresstorage.ErrActivationConfiguration) {
		t.Fatalf("Tabelle außerhalb des Alphabets: Fehler = %v, wollen ErrActivationConfiguration", err)
	}
}

// TestAdministrationRequestListPendingCarriesRuleRowsWithMissingFields trägt
// den Lese-Pfad der Transformations-Antragsarten gegen die reale PostgreSQL
// (`LH-FA-CFG-007`, `SPEC-019`): eine Zeile mit NULL-Regelnamen, mit
// SQL-NULL-Regelform oder mit JSON-`null` als Regelform entsteht als
// `pending`-Antrag (die Funktionen prüfen nichts) und wird von `ListPending`
// als Antrag geliefert, nicht als Fehler — der Text `null` bzw. der leere
// Text tragen SQL- und JSON-`null` unterscheidbar. Rot färbende Mutation: die
// Prüfung `ruleName == ""` in `model.NewAdministrationRequest` zurücklegen —
// `ListPending` endet mit `ErrEmptyIdentifier`, und jeder Antrag der Queue
// bleibt ungelesen.
func TestAdministrationRequestListPendingCarriesRuleRowsWithMissingFields(t *testing.T) {
	pool, dsn := newTestAdministrationRequestPool(t)
	ctx := context.Background()
	adapter, err := postgresstorage.NewAdministrationRequest(ctx, dsn)
	if err != nil {
		t.Fatalf("NewAdministrationRequest: %v", err)
	}
	t.Cleanup(adapter.Close)

	const table = "orders_rule_pending_rows"
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM cdc.administration_request WHERE table_name = $1", table)
	})
	create := func(call string, args ...any) model.AdministrationRequestID {
		t.Helper()
		var id string
		if err := pool.QueryRow(ctx, call, args...).Scan(&id); err != nil {
			t.Fatalf("%s: %v", call, err)
		}
		return model.AdministrationRequestID(id)
	}
	nullName := create("SELECT cdc.set_transformation($1, 'public', $2, NULL, '{\"kind\": \"rename_column\"}'::json)", administrationRequestSource, table)
	nullSpec := create("SELECT cdc.set_transformation($1, 'public', $2, 'regel_a', NULL::json)", administrationRequestSource, table)
	jsonNull := create("SELECT cdc.set_transformation($1, 'public', $2, 'regel_b', 'null'::json)", administrationRequestSource, table)
	removeNull := create("SELECT cdc.remove_transformation($1, 'public', $2, NULL)", administrationRequestSource, table)
	valid := create("SELECT cdc.set_transformation($1, 'public', $2, 'regel_c', '{\"kind\": \"rename_column\", \"column\": \"a\", \"to\": \"b\"}'::json)", administrationRequestSource, table)

	pending, err := adapter.ListPending(ctx)
	if err != nil {
		t.Fatalf("ListPending = %v, wollen nil (die Zeilen sind Anträge, keine Lesefehler)", err)
	}
	byID := map[model.AdministrationRequestID]model.AdministrationRequest{}
	for _, request := range pending {
		byID[request.ID] = request
	}
	for id, want := range map[model.AdministrationRequestID][2]string{
		nullName:   {"", `{"kind": "rename_column"}`},
		nullSpec:   {"regel_a", ""},
		jsonNull:   {"regel_b", "null"},
		removeNull: {"", ""},
	} {
		request, found := byID[id]
		if !found {
			t.Fatalf("ListPending trägt den Antrag %q nicht: %+v", id, pending)
		}
		if request.RuleName != want[0] || request.RuleSpec != want[1] {
			t.Fatalf("Antrag %q = Regelname %q, Regelform %q, wollen %q und %q", id, request.RuleName, request.RuleSpec, want[0], want[1])
		}
	}
	// Die gültige Zeile hinter den unvollständigen bleibt lesbar; `jsonb` ordnet
	// die Schlüssel um, der Inhalt bleibt.
	request, found := byID[valid]
	if !found || request.RuleName != "regel_c" {
		t.Fatalf("ListPending trägt die gültige Zeile nicht: %+v", request)
	}
	spec, err := model.ParseTransformationSpec(request.RuleSpec)
	if err != nil || spec.Column() != "a" || spec.Target() != "b" {
		t.Fatalf("Regelform der gültigen Zeile = %q (%v), wollen column a, to b", request.RuleSpec, err)
	}
}
