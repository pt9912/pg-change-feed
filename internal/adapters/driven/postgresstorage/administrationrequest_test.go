package postgresstorage_test

import (
	"context"
	stderrors "errors"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage"
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
		"SELECT count(*) FROM pg_proc p JOIN pg_namespace n ON n.oid = p.pronamespace WHERE n.nspname = 'cdc' AND p.proname IN ('enable_table', 'disable_table')",
	).Scan(&functions); err != nil {
		t.Fatalf("Funktions-Prüfung: %v", err)
	}
	if functions != 2 {
		t.Fatalf("cdc.enable_table/cdc.disable_table fehlen — der Schema-Rollout über make schema-rollout trägt sie (ADR-0050, tools/schema/nacharbeit-administration.sql); der test-store-Lauf rollt sie vor dem Testlauf aus")
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
// „permission denied for function" auslöst).
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
