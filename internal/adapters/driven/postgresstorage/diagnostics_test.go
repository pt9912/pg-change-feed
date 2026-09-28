package postgresstorage_test

import (
	"context"
	stderrors "errors"
	"net/url"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Die Diagnostics-Adapter-Tests laufen gegen dieselbe reale
// PostgreSQL-Instanz wie die übrigen Store-Tests (`make test-store`,
// `ADR-0030`); ohne DSN überspringen sie. `cdc.retention_blockers` und
// `cdc.backfill_status` trägt der d-migrate-Rollout
// (`tools/schema/schema.yaml`), `cdc.heartbeat`/`cdc.metrics` die Nacharbeit
// (`tools/schema/nacharbeit-heartbeat.sql`,
// `tools/schema/nacharbeit-observability.sql`) — beide trägt die
// handgeschriebene DDL des Store-Adapters nicht.
//
// Kopplung: diese Datei sortiert vor `heartbeat_test.go`/`store_test.go`
// (`d` < `h`/`s`) und läuft deshalb vor deren Neuaufbauten
// (`DROP SCHEMA cdc CASCADE`) — dieselbe Voraussetzung wie
// `backfilladmission_test.go`/`administrationrequest_test.go`: das per
// d-migrate ausgerollte Schema muss zu diesem Zeitpunkt bereits bestehen.
const diagnosticsTestSource = "src-diagnostics-adapter"
const diagnosticsTestConsumer = "diagnostics-adapter-consumer"
const diagnosticsTestTable = "diagnostics_adapter_table"

// newDiagnosticsFixture baut den Adapter gegen die Test-Instanz und liefert
// zusätzlich den Superuser-Pool für Fixture-Zeilen; der Datenstand der
// Test-Kennungen wird vorab geräumt.
func newDiagnosticsFixture(t *testing.T) (*postgresstorage.PostgresDiagnosticsAdapter, *pgxpool.Pool, string) {
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

	diagnosticsFixtureCleanup(t, pool)
	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.source (source_id, name) VALUES ($1, 'Diagnostics-Adapter-Quelle') ON CONFLICT (source_id) DO NOTHING",
		diagnosticsTestSource,
	); err != nil {
		t.Fatalf("Quelle-Zeile: %v", err)
	}

	adapter, err := postgresstorage.NewDiagnostics(ctx, dsn)
	if err != nil {
		t.Fatalf("NewDiagnostics: %v", err)
	}
	t.Cleanup(adapter.Close)
	return adapter, pool, dsn
}

// diagnosticsFixtureCleanup räumt den Datenstand der Test-Kennungen ab —
// vor und nach jedem Test, damit die vier Signalgruppen unabhängig
// voneinander steuerbar bleiben (Isolation, `BEO-PGC/test-isolation-geteilter-zustand`).
func diagnosticsFixtureCleanup(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()
	for _, stmt := range []string{
		"DELETE FROM cdc.consumer_position WHERE consumer_id = $1",
		"DELETE FROM cdc.consumer WHERE consumer_id = $1",
	} {
		if _, err := pool.Exec(ctx, stmt, diagnosticsTestConsumer); err != nil {
			t.Fatalf("Datenstand-Rückbau (%s): %v", stmt, err)
		}
	}
	for _, stmt := range []string{
		"DELETE FROM cdc.backfill_run WHERE source_id = $1",
		"DELETE FROM cdc.process_heartbeat WHERE source_id = $1",
		"DELETE FROM cdc.transaction WHERE source_id = $1",
	} {
		if _, err := pool.Exec(ctx, stmt, diagnosticsTestSource); err != nil {
			t.Fatalf("Datenstand-Rückbau (%s): %v", stmt, err)
		}
	}
}

// TestDiagnosticsReadNormalOperation belegt den Happy Path aller sechs
// Signalgruppen in einem Aufruf (`ADR-0132`): ein gesetztes Lebenszeichen mit
// Fehlerzustand, ein numerischer CDC-Abstand, der Rückstand des
// Test-Consumers, derselbe Consumer als Retention-Blocker, ein numerischer
// Speicherverbrauch und ein Backfill-Run.
func TestDiagnosticsReadNormalOperation(t *testing.T) {
	adapter, pool, _ := newDiagnosticsFixture(t)
	ctx := context.Background()
	t.Cleanup(func() { diagnosticsFixtureCleanup(t, pool) })

	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.transaction (transaction_id, source_id, commit_position, committed_at) VALUES ('tx-diagnostics-adapter-1', $1, 1, current_timestamp)",
		diagnosticsTestSource,
	); err != nil {
		t.Fatalf("Transaktion-Zeile 1: %v", err)
	}
	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.transaction (transaction_id, source_id, commit_position, committed_at) VALUES ('tx-diagnostics-adapter-2', $1, 2, current_timestamp)",
		diagnosticsTestSource,
	); err != nil {
		t.Fatalf("Transaktion-Zeile 2: %v", err)
	}
	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.consumer (consumer_id, name) VALUES ($1, $1)",
		diagnosticsTestConsumer,
	); err != nil {
		t.Fatalf("Consumer-Zeile: %v", err)
	}
	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.consumer_position (consumer_id, source_id, acknowledged_position) VALUES ($1, $2, 1)",
		diagnosticsTestConsumer, diagnosticsTestSource,
	); err != nil {
		t.Fatalf("Consumer-Position-Zeile: %v", err)
	}
	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.process_heartbeat (source_id, heartbeat_at, error_class) VALUES ($1, current_timestamp, 'schema')",
		diagnosticsTestSource,
	); err != nil {
		t.Fatalf("Lebenszeichen-Zeile: %v", err)
	}
	if _, err := pool.Exec(ctx,
		`INSERT INTO cdc.backfill_run
		    (run_id, source_id, schema_name, table_name, status, requested_at, rows_copied, estimated_rows, warn_estimated_size, warn_duration, error_message)
		 VALUES ('run-diagnostics-adapter-1', $1, 'public', $2, 'completed', current_timestamp, 12, 10, false, false, NULL)`,
		diagnosticsTestSource, diagnosticsTestTable,
	); err != nil {
		t.Fatalf("Backfill-Run-Zeile: %v", err)
	}

	snapshot, err := adapter.Read(ctx, model.SourceID(diagnosticsTestSource))
	if err != nil {
		t.Fatalf("Read: %v", err)
	}

	if snapshot.HeartbeatAgeSeconds == nil {
		t.Fatalf("HeartbeatAgeSeconds = nil, wollen einen gesetzten Wert")
	}
	if *snapshot.HeartbeatAgeSeconds < 0 || *snapshot.HeartbeatAgeSeconds > 5 {
		t.Fatalf("HeartbeatAgeSeconds = %v, wollen einen kleinen Wert direkt nach dem Schreiben", *snapshot.HeartbeatAgeSeconds)
	}
	if snapshot.ErrorClass == nil || *snapshot.ErrorClass != "schema" {
		t.Fatalf("ErrorClass = %v, wollen \"schema\"", snapshot.ErrorClass)
	}
	if snapshot.CaptureLag < 0 {
		t.Fatalf("CaptureLag = %v, wollen einen nichtnegativen Wert", snapshot.CaptureLag)
	}
	if len(snapshot.ConsumerLags) != 1 {
		t.Fatalf("ConsumerLags = %v, wollen genau einen Eintrag", snapshot.ConsumerLags)
	}
	if lag := snapshot.ConsumerLags[0]; lag.ConsumerID != diagnosticsTestConsumer || lag.Lag == nil || *lag.Lag != 1 {
		t.Fatalf("ConsumerLags[0] = %+v, wollen ConsumerID %q und Lag 1", lag, diagnosticsTestConsumer)
	}
	if snapshot.RetentionBlocker == nil {
		t.Fatalf("RetentionBlocker = nil, wollen den Test-Consumer als Blocker")
	}
	if snapshot.RetentionBlocker.ConsumerID != diagnosticsTestConsumer || snapshot.RetentionBlocker.Name != diagnosticsTestConsumer {
		t.Fatalf("RetentionBlocker = %+v, wollen ConsumerID/Name %q", snapshot.RetentionBlocker, diagnosticsTestConsumer)
	}
	if snapshot.RetentionBlocker.AcknowledgedPosition != 1 {
		t.Fatalf("RetentionBlocker.AcknowledgedPosition = %d, wollen 1", snapshot.RetentionBlocker.AcknowledgedPosition)
	}
	if snapshot.RetentionBlocker.Backlog == nil || *snapshot.RetentionBlocker.Backlog != 1 {
		t.Fatalf("RetentionBlocker.Backlog = %v, wollen 1", snapshot.RetentionBlocker.Backlog)
	}
	if snapshot.StorageBytes < 0 {
		t.Fatalf("StorageBytes = %v, wollen einen nichtnegativen Wert", snapshot.StorageBytes)
	}
	if len(snapshot.Backfill) != 1 {
		t.Fatalf("Backfill = %v, wollen genau einen Eintrag", snapshot.Backfill)
	}
	table := snapshot.Backfill[0]
	if table.Schema != "public" || table.Table != diagnosticsTestTable || table.Status != "completed" {
		t.Fatalf("Backfill[0] = %+v, wollen Schema public, Table %q, Status completed", table, diagnosticsTestTable)
	}
	if table.RowsCopied != 12 {
		t.Fatalf("Backfill[0].RowsCopied = %d, wollen 12", table.RowsCopied)
	}
	if table.EstimatedRows == nil || *table.EstimatedRows != 10 {
		t.Fatalf("Backfill[0].EstimatedRows = %v, wollen 10", table.EstimatedRows)
	}
	if table.WarnEstimatedSize || table.WarnDuration {
		t.Fatalf("Backfill[0] Warnungen = %v/%v, wollen false/false", table.WarnEstimatedSize, table.WarnDuration)
	}
	if table.ErrorMessage != "" {
		t.Fatalf("Backfill[0].ErrorMessage = %q, wollen leer (COALESCE)", table.ErrorMessage)
	}
}

// TestDiagnosticsReadNoHeartbeat belegt den `sqlexec.IsAbsent`-Zweig von
// `SelectDiagnosticsHeartbeat` (`ADR-0132`): eine Quelle ohne Lebenszeichen
// trägt `HeartbeatAgeSeconds == nil` und `ErrorClass == nil`, der
// Lesezugriff scheitert trotzdem nicht — die übrigen fünf Signalgruppen
// bleiben unabhängig lesbar (hier: leere, gesetzte Listen, kein Blocker).
func TestDiagnosticsReadNoHeartbeat(t *testing.T) {
	adapter, pool, _ := newDiagnosticsFixture(t)
	t.Cleanup(func() { diagnosticsFixtureCleanup(t, pool) })

	snapshot, err := adapter.Read(context.Background(), model.SourceID(diagnosticsTestSource))
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if snapshot.HeartbeatAgeSeconds != nil {
		t.Fatalf("HeartbeatAgeSeconds = %v, wollen nil (kein Lebenszeichen)", *snapshot.HeartbeatAgeSeconds)
	}
	if snapshot.ErrorClass != nil {
		t.Fatalf("ErrorClass = %v, wollen nil (kein Lebenszeichen)", *snapshot.ErrorClass)
	}
}

// TestDiagnosticsReadNoRetentionBlockerAndEmptyListsAreNeverNil belegt zwei
// Dinge in einem Aufruf: den `sqlexec.IsAbsent`-Zweig von
// `SelectDiagnosticsRetentionBlocker` (kein Consumer hat je gegen diese
// Quelle bestätigt) und die Zusage von `scanConsumerLags`/
// `scanBackfillStatus`, ohne Zeilen eine leere, gesetzte Liste zu liefern —
// nie `nil`.
func TestDiagnosticsReadNoRetentionBlockerAndEmptyListsAreNeverNil(t *testing.T) {
	adapter, pool, _ := newDiagnosticsFixture(t)
	t.Cleanup(func() { diagnosticsFixtureCleanup(t, pool) })

	snapshot, err := adapter.Read(context.Background(), model.SourceID(diagnosticsTestSource))
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if snapshot.RetentionBlocker != nil {
		t.Fatalf("RetentionBlocker = %+v, wollen nil (kein Consumer hat je bestätigt)", snapshot.RetentionBlocker)
	}
	if snapshot.ConsumerLags == nil || len(snapshot.ConsumerLags) != 0 {
		t.Fatalf("ConsumerLags = %#v, wollen eine leere, gesetzte Liste", snapshot.ConsumerLags)
	}
	if snapshot.Backfill == nil || len(snapshot.Backfill) != 0 {
		t.Fatalf("Backfill = %#v, wollen eine leere, gesetzte Liste", snapshot.Backfill)
	}
}

// TestNewDiagnosticsConnectsAndPingsSuccessfully belegt den Erfolgsfall von
// `NewDiagnostics` (`ADR-0132`) explizit — dieselbe Konstruktion, die
// `newDiagnosticsFixture` bereits implizit in jedem anderen Test dieser Datei
// durchläuft.
func TestNewDiagnosticsConnectsAndPingsSuccessfully(t *testing.T) {
	dsn := os.Getenv("CDC_STORE_TEST_DSN")
	if dsn == "" {
		t.Skip("CDC_STORE_TEST_DSN nicht gesetzt — reale PostgreSQL-Tests laufen über make test-store")
	}
	adapter, err := postgresstorage.NewDiagnostics(context.Background(), dsn)
	if err != nil {
		t.Fatalf("NewDiagnostics: %v", err)
	}
	defer adapter.Close()
}

// TestNewDiagnosticsCarriesStorageClass belegt denselben Fehlerpfad wie
// `TestNewConsumerStateCarriesStorageClass` (`consumerstate_test.go`): eine
// nicht erreichbare Instanz meldet den Aufbau als Fehler der Klasse
// `storage` über den eigenen Sentinel des Ports; der Test braucht keine
// Datenbank (Port 1 verwirft lokal, netzlos).
func TestNewDiagnosticsCarriesStorageClass(t *testing.T) {
	_, err := postgresstorage.NewDiagnostics(context.Background(), "postgres://cdc:cdc@127.0.0.1:1/cdc_test?sslmode=disable")
	if !stderrors.Is(err, outbound.ErrDiagnosticsStorage) {
		t.Fatalf("Fehler = %v, wollen Klasse storage (%v)", err, outbound.ErrDiagnosticsStorage)
	}
}

// TestDiagnosticsReadReportsUnreadableViewAsStorageFailure belegt den
// Storage-Fehlerpfad von `Read` selbst: eine Login-Identität ohne
// `cdc_reader`-Mitgliedschaft (keine `USAGE`-Berechtigung auf dem Schema
// `cdc`) scheitert am ersten Diagnose-Lesezugriff (`cdc.heartbeat`) mit
// einem Berechtigungsfehler, den `Read` als `ErrDiagnosticsStorage` meldet —
// dasselbe Rollen-Muster wie
// `internal/bootstrap/diagnose_test.go`s
// `TestDiagnoseReportsUnreadableViewAsStorageFailure`, hier direkt am
// Adapter statt über den CLI-Sondermodus.
func TestDiagnosticsReadReportsUnreadableViewAsStorageFailure(t *testing.T) {
	baseDSN := os.Getenv("CDC_STORE_TEST_DSN")
	if baseDSN == "" {
		t.Skip("CDC_STORE_TEST_DSN nicht gesetzt — reale PostgreSQL-Tests laufen über make test-store")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, baseDSN)
	if err != nil {
		t.Fatalf("Verbindungsaufbau: %v", err)
	}
	t.Cleanup(pool.Close)

	const login, password = "pgc_test_diagnostics_no_grant", "test-login-password"
	if _, err := pool.Exec(ctx, "DROP ROLE IF EXISTS "+login); err != nil {
		t.Fatalf("Vorab-Aufräumen der Login-Identität: %v", err)
	}
	if _, err := pool.Exec(ctx, "CREATE ROLE "+login+" LOGIN PASSWORD '"+password+"'"); err != nil {
		t.Fatalf("Login-Identität anlegen (ohne cdc_reader-Mitgliedschaft): %v", err)
	}
	t.Cleanup(func() { _, _ = pool.Exec(context.Background(), "DROP ROLE IF EXISTS "+login) })

	parsed, err := url.Parse(baseDSN)
	if err != nil {
		t.Fatalf("Basis-DSN nicht parsebar: %v", err)
	}
	parsed.User = url.UserPassword(login, password)

	adapter, err := postgresstorage.NewDiagnostics(ctx, parsed.String())
	if err != nil {
		t.Fatalf("NewDiagnostics (Verbindung/Ping gelingt auch ohne cdc_reader-Mitgliedschaft): %v", err)
	}
	t.Cleanup(adapter.Close)

	_, err = adapter.Read(ctx, model.SourceID(diagnosticsTestSource))
	if !stderrors.Is(err, outbound.ErrDiagnosticsStorage) {
		t.Fatalf("Read-Fehler = %v, wollen Klasse storage (%v)", err, outbound.ErrDiagnosticsStorage)
	}
}
