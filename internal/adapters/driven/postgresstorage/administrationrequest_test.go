package postgresstorage_test

import (
	"context"
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
		"SELECT count(*) FROM pg_proc p JOIN pg_namespace n ON n.oid = p.pronamespace WHERE n.nspname = 'cdc' AND p.proname IN ('enable_table', 'disable_table', 'exclude_column', 'include_column', 'backfill_table')",
	).Scan(&functions); err != nil {
		t.Fatalf("Funktions-Prüfung: %v", err)
	}
	if functions != 5 {
		t.Fatalf("cdc.enable_table/cdc.disable_table/cdc.exclude_column/cdc.include_column/cdc.backfill_table fehlen — der Schema-Rollout über make schema-rollout trägt sie (ADR-0050, LH-FA-CFG-005, LH-FA-CAP-009, tools/schema/nacharbeit-administration.sql); der test-store-Lauf rollt sie vor dem Testlauf aus")
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
// Publication (Slice-Abgrenzung, Verarbeitung folgt in einem Folge-Slice).
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

// TestAdministrationRequestKindCheckCarriesExactlyTheFiveKinds belegt die
// geschlossene `request_kind`-Menge (`chk_administration_request_kind`,
// `tools/schema/nacharbeit-administration.sql`): jede der fünf Arten wird
// angenommen, eine sechste endet mit SQLSTATE 23514. Rot färbende Mutationen
// (je eine): `'backfill'` aus der CHECK-Klausel streichen — die Art `backfill`
// endet mit 23514; `'truncate'` in die Klausel aufnehmen — die sechste Art
// wird angenommen.
func TestAdministrationRequestKindCheckCarriesExactlyTheFiveKinds(t *testing.T) {
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
	for _, kind := range []string{"enable", "disable", "exclude_column", "include_column", "backfill"} {
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
