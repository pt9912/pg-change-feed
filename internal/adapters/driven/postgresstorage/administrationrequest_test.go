package postgresstorage_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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
