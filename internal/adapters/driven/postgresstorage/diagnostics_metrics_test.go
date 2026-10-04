package postgresstorage_test

import (
	"context"
	stderrors "errors"
	"math"
	"net/url"
	"os"
	"strconv"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
)

// Die Metrics-Adapter-Tests laufen gegen dieselbe reale PostgreSQL-Instanz wie
// die übrigen Store-Tests (`make test-store`); ohne DSN überspringen sie. Die
// Sicht `cdc.metrics` trägt die Nacharbeit
// (`tools/schema/nacharbeit-observability.sql`). Der Adapter liest unter einer
// Login-Identität mit der Mitgliedschaft `cdc_reader` und sonst keiner
// Berechtigung.
const (
	metricsTestSource   = "src-metrics-adapter"
	metricsTestConsumer = "metrics-adapter-consumer"
	metricsTestLogin    = "pgc_test_metrics_reader"
	metricsTestPassword = "test-login-password"

	// Positionen jenseits von 2^53: ein `float64` rundet sie, die ganze Zahl
	// des Adapters nicht.
	metricsTestAcknowledged = int64(9007199254740993)
	metricsTestCommit       = int64(9007199254740995)
)

func metricsFixtureCleanup(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()
	for _, stmt := range []string{
		"DELETE FROM cdc.consumer_position WHERE consumer_id = $1",
		"DELETE FROM cdc.consumer WHERE consumer_id = $1",
	} {
		if _, err := pool.Exec(ctx, stmt, metricsTestConsumer); err != nil {
			t.Fatalf("Datenstand-Rückbau (%s): %v", stmt, err)
		}
	}
	for _, stmt := range []string{
		"DELETE FROM cdc.process_heartbeat WHERE source_id = $1",
		"DELETE FROM cdc.transaction WHERE source_id = $1",
	} {
		if _, err := pool.Exec(ctx, stmt, metricsTestSource); err != nil {
			t.Fatalf("Datenstand-Rückbau (%s): %v", stmt, err)
		}
	}
}

// metricsLoginDSN legt die Login-Identität an (mit oder ohne
// `cdc_reader`-Mitgliedschaft) und liefert ihre DSN.
func metricsLoginDSN(t *testing.T, pool *pgxpool.Pool, baseDSN string, asReader bool) string {
	t.Helper()
	ctx := context.Background()
	if _, err := pool.Exec(ctx, "DROP ROLE IF EXISTS "+metricsTestLogin); err != nil {
		t.Fatalf("Vorab-Aufräumen der Login-Identität: %v", err)
	}
	if _, err := pool.Exec(ctx, "CREATE ROLE "+metricsTestLogin+" LOGIN PASSWORD '"+metricsTestPassword+"'"); err != nil {
		t.Fatalf("Login-Identität anlegen: %v", err)
	}
	t.Cleanup(func() { _, _ = pool.Exec(context.Background(), "DROP ROLE IF EXISTS "+metricsTestLogin) })
	if asReader {
		if _, err := pool.Exec(ctx, "GRANT cdc_reader TO "+metricsTestLogin); err != nil {
			t.Fatalf("cdc_reader-Mitgliedschaft: %v", err)
		}
	}
	parsed, err := url.Parse(baseDSN)
	if err != nil {
		t.Fatalf("Basis-DSN nicht parsebar: %v", err)
	}
	parsed.User = url.UserPassword(metricsTestLogin, metricsTestPassword)
	return parsed.String()
}

func newMetricsPool(t *testing.T) (*pgxpool.Pool, string) {
	t.Helper()
	dsn := os.Getenv("CDC_STORE_TEST_DSN")
	if dsn == "" {
		t.Skip("CDC_STORE_TEST_DSN nicht gesetzt — reale PostgreSQL-Tests laufen über make test-store")
	}
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		t.Fatalf("Verbindungsaufbau: %v", err)
	}
	t.Cleanup(pool.Close)
	return pool, dsn
}

// TestMetricsReadEqualsTheView belegt `LH-FA-SST-010` (Quelle der Werte): die
// Ausgabe des Adapters unter einem `cdc_reader`-Login ist je Name und Label
// gleich einem `SELECT metric_name, label, value FROM cdc.metrics` derselben
// Datenbank — mit einem Consumer (Zeilen mit Label), einem Fehlerzustand
// (Zeile `cdc_errors_total` mit Klasse) und Positionen jenseits von 2^53. Die
// zwei zeitabhängigen Kennzahlen (`cdc_oldest_change_age_seconds`,
// `cdc_capture_lag`) hängen an `now()` der Abfrage und weichen um höchstens
// fünf Sekunden ab (benannte Toleranz); jede andere Zahl ist gleich.
// Rot färbende Mutation: in `Read` das Label verwerfen (`sample.Label` nicht
// setzen) — der Schlüssel der Consumer-Zeilen stimmt nicht mehr; `Int` aus
// `parseMetricValue` streichen — die Position nennt eine gerundete Zahl.
func TestMetricsReadEqualsTheView(t *testing.T) {
	pool, baseDSN := newMetricsPool(t)
	ctx := context.Background()
	metricsFixtureCleanup(t, pool)
	t.Cleanup(func() { metricsFixtureCleanup(t, pool) })

	for _, step := range []struct {
		sql  string
		args []any
	}{
		{"INSERT INTO cdc.source (source_id, name) VALUES ($1, 'Metrics-Adapter-Quelle') ON CONFLICT (source_id) DO NOTHING", []any{metricsTestSource}},
		{"INSERT INTO cdc.transaction (transaction_id, source_id, commit_position, committed_at) VALUES ('tx-metrics-adapter-1', $1, $2, current_timestamp)", []any{metricsTestSource, metricsTestCommit}},
		{"INSERT INTO cdc.consumer (consumer_id, name) VALUES ($1, $1)", []any{metricsTestConsumer}},
		{"INSERT INTO cdc.consumer_position (consumer_id, source_id, acknowledged_position) VALUES ($1, $2, $3)", []any{metricsTestConsumer, metricsTestSource, metricsTestAcknowledged}},
		{"INSERT INTO cdc.process_heartbeat (source_id, heartbeat_at, error_class, error_code) VALUES ($1, current_timestamp, 'schema', 'PCF-E4003')", []any{metricsTestSource}},
	} {
		if _, err := pool.Exec(ctx, step.sql, step.args...); err != nil {
			t.Fatalf("Fixture (%s): %v", step.sql, err)
		}
	}

	adapter, err := postgresstorage.NewMetrics(ctx, metricsLoginDSN(t, pool, baseDSN, true))
	if err != nil {
		t.Fatalf("NewMetrics: %v", err)
	}
	t.Cleanup(adapter.Close)

	samples, err := adapter.Read(ctx)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}

	rows, err := pool.Query(ctx, "SELECT metric_name, COALESCE(label, ''), value::text FROM cdc.metrics")
	if err != nil {
		t.Fatalf("Sicht lesen: %v", err)
	}
	defer rows.Close()
	want := map[string]string{}
	for rows.Next() {
		var name, label, value string
		if err := rows.Scan(&name, &label, &value); err != nil {
			t.Fatalf("Zeile der Sicht: %v", err)
		}
		want[name+"|"+label] = value
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("Zeilen der Sicht: %v", err)
	}

	got := map[string]outbound.MetricSample{}
	for _, s := range samples {
		got[s.Name+"|"+s.Label] = s
	}
	if len(got) != len(samples) {
		t.Fatalf("Adapter lieferte einen Schlüssel doppelt: %d Zeilen, %d Schlüssel", len(samples), len(got))
	}
	for key, text := range want {
		sample, ok := got[key]
		if !ok {
			t.Errorf("%s: die Sicht trägt die Zeile, der Adapter nicht", key)
			continue
		}
		wantFloat, err := strconv.ParseFloat(text, 64)
		if err != nil {
			t.Fatalf("%s: Wert der Sicht %q nicht lesbar: %v", key, text, err)
		}
		switch sample.Name {
		case "cdc_oldest_change_age_seconds", "cdc_capture_lag":
			if math.Abs(sample.Value.Float-wantFloat) > 5 {
				t.Errorf("%s: Adapter %v, Sicht %v, Abstand über der Toleranz von 5 s", key, sample.Value.Float, wantFloat)
			}
		default:
			if sample.Value.IsInt {
				if strconv.FormatInt(sample.Value.Int, 10) != text {
					t.Errorf("%s: Adapter %d, Sicht %s", key, sample.Value.Int, text)
				}
			} else if sample.Value.Float != wantFloat {
				t.Errorf("%s: Adapter %v, Sicht %s", key, sample.Value.Float, text)
			}
		}
	}
	for key := range got {
		if _, ok := want[key]; !ok {
			t.Errorf("%s: der Adapter trägt eine Zeile, die die Sicht nicht trägt", key)
		}
	}

	position, ok := got["cdc_consumer_position|"+metricsTestConsumer]
	if !ok || !position.Value.IsInt || position.Value.Int != metricsTestAcknowledged {
		t.Errorf("cdc_consumer_position = %+v, wollen die ganze Zahl %d", position, metricsTestAcknowledged)
	}
	lag, ok := got["cdc_consumer_lag|"+metricsTestConsumer]
	if !ok || !lag.Value.IsInt || lag.Value.Int != metricsTestCommit-metricsTestAcknowledged {
		t.Errorf("cdc_consumer_lag = %+v, wollen %d", lag, metricsTestCommit-metricsTestAcknowledged)
	}
	if _, ok := got["cdc_errors_total|schema"]; !ok {
		t.Error("cdc_errors_total mit Klasse schema fehlt")
	}
	if _, ok := got["cdc_transactions_total|"]; !ok {
		t.Error("cdc_transactions_total ohne Label fehlt")
	}
}

// TestMetricsReadWithoutReaderMembershipFailsAsReadError belegt den Lesefehler
// am Adapter: eine Login-Identität ohne `cdc_reader`-Mitgliedschaft hat kein
// `USAGE` auf dem Schema; `Read` meldet den Berechtigungsfehler als
// `ErrMetricsRead`. Die Gegenprobe mit der Mitgliedschaft ist
// `TestMetricsReadEqualsTheView`.
func TestMetricsReadWithoutReaderMembershipFailsAsReadError(t *testing.T) {
	pool, baseDSN := newMetricsPool(t)
	adapter, err := postgresstorage.NewMetrics(context.Background(), metricsLoginDSN(t, pool, baseDSN, false))
	if err != nil {
		t.Fatalf("NewMetrics: %v", err)
	}
	t.Cleanup(adapter.Close)

	_, err = adapter.Read(context.Background())
	if !stderrors.Is(err, outbound.ErrMetricsRead) {
		t.Fatalf("Read-Fehler = %v, wollen ErrMetricsRead", err)
	}
}

// TestMetricsReadOfUnreachableInstanceFailsAsReadError belegt, dass der Pool
// erst beim Lesen verbindet: die Konstruktion einer nicht erreichbaren
// Instanz gelingt, das Lesen endet als `ErrMetricsRead`. Der Test braucht
// keine Datenbank (Port 1 verwirft lokal, netzlos).
func TestMetricsReadOfUnreachableInstanceFailsAsReadError(t *testing.T) {
	adapter, err := postgresstorage.NewMetrics(context.Background(), "postgres://cdc:cdc@127.0.0.1:1/cdc_test?sslmode=disable")
	if err != nil {
		t.Fatalf("NewMetrics (verbindet erst beim Lesen): %v", err)
	}
	t.Cleanup(adapter.Close)

	_, err = adapter.Read(context.Background())
	if !stderrors.Is(err, outbound.ErrMetricsRead) {
		t.Fatalf("Read-Fehler = %v, wollen ErrMetricsRead", err)
	}
}

// TestNewMetricsRejectsUnparsableDSN belegt den Konstruktionsfehler: eine
// DSN, die der Treiber nicht liest, endet als `ErrMetricsRead`.
func TestNewMetricsRejectsUnparsableDSN(t *testing.T) {
	_, err := postgresstorage.NewMetrics(context.Background(), "postgres://[::1")
	if !stderrors.Is(err, outbound.ErrMetricsRead) {
		t.Fatalf("Fehler = %v, wollen ErrMetricsRead", err)
	}
}
