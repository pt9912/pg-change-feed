package bootstrap_test

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pt9912/pg-change-feed/internal/bootstrap"
)

// captureStdout (register_test.go) trägt bereits dasselbe
// Umleitungs-Muster für `os.Stdout`, das `Diagnose`s Bericht hier braucht —
// keine zweite Deklaration in diesem Paket.

// TestDiagnoseReportsConnectionFailure belegt dieselbe erste Fehlerklasse
// wie `TestHealthcheckReportsConnectionFailure`: eine nicht erreichbare
// Instanz trägt eine eigene, unterscheidbare stderr-Zeile — netzlos
// (`make test`), die Verbindung scheitert am geschlossenen lokalen Port.
func TestDiagnoseReportsConnectionFailure(t *testing.T) {
	var code int
	output := captureStderr(t, func() {
		code = bootstrap.Diagnose(context.Background(), "postgres://x:x@127.0.0.1:1/db?connect_timeout=1", "src-1")
	})
	if code != 1 {
		t.Fatalf("Diagnose-Exit-Code = %d, wollen 1 (Verbindungsfehler)", code)
	}
	if !strings.Contains(output, "nicht erreichbar") {
		t.Fatalf("stderr = %q, wollen eine Diagnose-Zeile mit dem Text 'nicht erreichbar'", output)
	}
}

const diagnoseTestSource = "src-diagnose-test"
const diagnoseTestConsumer = "diagnose-test-consumer"

// newDiagnoseTestFixture baut den Datenstand gegen die reale
// PostgreSQL-Testinstanz (`make test-store`, `ADR-0030`): eine Quelle mit
// einer Transaktion (trägt `cdc_capture_lag` real, wenn auch als
// tabellenweiter Wert — `cdc.metrics`, `nacharbeit-observability.sql`,
// filtert diese Zeile nicht nach Quelle) und ein Consumer mit einer
// bestätigten Position hinter der zuletzt erfassten (trägt
// `cdc_consumer_lag` > 0 für genau diesen Consumer). Ohne DSN überspringt
// der Test.
func newDiagnoseTestFixture(t *testing.T) *pgxpool.Pool {
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

	for _, stmt := range []string{
		"DELETE FROM cdc.consumer_position WHERE consumer_id = $1",
		"DELETE FROM cdc.consumer WHERE consumer_id = $1",
	} {
		if _, err := pool.Exec(ctx, stmt, diagnoseTestConsumer); err != nil {
			t.Fatalf("Datenstand-Rückbau (%s): %v", stmt, err)
		}
	}
	for _, stmt := range []string{
		"DELETE FROM cdc.process_heartbeat WHERE source_id = $1",
		"DELETE FROM cdc.transaction WHERE source_id = $1",
	} {
		if _, err := pool.Exec(ctx, stmt, diagnoseTestSource); err != nil {
			t.Fatalf("Datenstand-Rückbau (%s): %v", stmt, err)
		}
	}
	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.source (source_id, name) VALUES ($1, 'Diagnose-Testquelle') ON CONFLICT (source_id) DO NOTHING",
		diagnoseTestSource,
	); err != nil {
		t.Fatalf("Quelle-Zeile: %v", err)
	}
	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.transaction (transaction_id, source_id, commit_position, committed_at) VALUES ('tx-diagnose-1', $1, 1, current_timestamp)",
		diagnoseTestSource,
	); err != nil {
		t.Fatalf("Transaktion-Zeile 1: %v", err)
	}
	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.transaction (transaction_id, source_id, commit_position, committed_at) VALUES ('tx-diagnose-2', $1, 2, current_timestamp)",
		diagnoseTestSource,
	); err != nil {
		t.Fatalf("Transaktion-Zeile 2: %v", err)
	}
	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.consumer (consumer_id, name) VALUES ($1, $1)",
		diagnoseTestConsumer,
	); err != nil {
		t.Fatalf("Consumer-Zeile: %v", err)
	}
	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.consumer_position (consumer_id, source_id, acknowledged_position) VALUES ($1, $2, 1)",
		diagnoseTestConsumer, diagnoseTestSource,
	); err != nil {
		t.Fatalf("Consumer-Position-Zeile: %v", err)
	}
	return pool
}

// TestDiagnoseReportsNormalOperation belegt den Happy Path aller vier
// Signale (`LH-FA-SST-003`, deckt `LH-FA-ADM-002`…`005`): ein gesetztes
// Lebenszeichen ohne Fehlerzustand, ein numerischer `cdc_capture_lag`-Wert
// und der Rückstand des Test-Consumers (bestätigte Position 1 hinter der
// zuletzt erfassten Position 2, also Rückstand 1).
func TestDiagnoseReportsNormalOperation(t *testing.T) {
	pool := newDiagnoseTestFixture(t)
	ctx := context.Background()
	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.process_heartbeat (source_id, heartbeat_at, error_class) VALUES ($1, current_timestamp, NULL)",
		diagnoseTestSource,
	); err != nil {
		t.Fatalf("Lebenszeichen-Zeile: %v", err)
	}

	dsn := os.Getenv("CDC_STORE_TEST_DSN")
	var code int
	output := captureStdout(t, func() {
		code = bootstrap.Diagnose(ctx, dsn, diagnoseTestSource)
	})
	if code != 0 {
		t.Fatalf("Diagnose-Exit-Code = %d, wollen 0 (erfolgreicher Lesezugriff)", code)
	}
	if !strings.Contains(output, "Betriebsstatus (LH-FA-ADM-002): Lebenszeichen vor") {
		t.Fatalf("stdout = %q, wollen eine Betriebsstatus-Zeile", output)
	}
	if !strings.Contains(output, "Fehlerzustand (LH-FA-ADM-003): keiner (Normalbetrieb)") {
		t.Fatalf("stdout = %q, wollen 'keiner (Normalbetrieb)' im Normalfall", output)
	}
	if !strings.Contains(output, "CDC-Abstand cdc_capture_lag (LH-FA-ADM-004):") {
		t.Fatalf("stdout = %q, wollen eine cdc_capture_lag-Zeile", output)
	}
	if !strings.Contains(output, diagnoseTestConsumer+": 1") {
		t.Fatalf("stdout = %q, wollen den Rückstand des Test-Consumers (%s: 1)", output, diagnoseTestConsumer)
	}
}

// TestDiagnoseReportsErrorState belegt die Boundary von `LH-FA-ADM-003`:
// ein gesetzter Fehlerzustand ist über eine eigene, von Normalbetrieb
// unterscheidbare Zeile sichtbar — dieselbe Spalte
// (`cdc.process_heartbeat.error_class`), die `reportFault`
// (`internal/bootstrap/wiring.go`) im realen Fehlerfall schreibt, hier
// direkt gesetzt: ein realer, den laufenden Prozess beendender Fehler
// lässt sich an einem laufenden Feed-Container nicht mehr per CLI
// beobachten (siehe Plan-Nachzug slice-038 §3, Risiko 1 aus §6) — dieser
// Test belegt denselben Lesepfad (View → CLI-Ausgabe) unabhängig davon.
func TestDiagnoseReportsErrorState(t *testing.T) {
	pool := newDiagnoseTestFixture(t)
	ctx := context.Background()
	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.process_heartbeat (source_id, heartbeat_at, error_class) VALUES ($1, current_timestamp, 'schema')",
		diagnoseTestSource,
	); err != nil {
		t.Fatalf("Lebenszeichen-Zeile mit Fehlerzustand: %v", err)
	}

	dsn := os.Getenv("CDC_STORE_TEST_DSN")
	var code int
	output := captureStdout(t, func() {
		code = bootstrap.Diagnose(ctx, dsn, diagnoseTestSource)
	})
	if code != 0 {
		t.Fatalf("Diagnose-Exit-Code = %d, wollen 0 (erfolgreicher Lesezugriff, der Fehlerzustand ist Berichtsinhalt)", code)
	}
	if !strings.Contains(output, "Fehlerzustand (LH-FA-ADM-003): schema") {
		t.Fatalf("stdout = %q, wollen die Fehlerklasse 'schema' sichtbar und von Normalbetrieb unterscheidbar", output)
	}
	if strings.Contains(output, "keiner (Normalbetrieb)") {
		t.Fatalf("stdout = %q, trägt fälschlich die Normalbetrieb-Zeile trotz gesetztem Fehlerzustand", output)
	}
}

// TestDiagnoseReportsNoHeartbeat belegt den zweiten Fehlerzweig
// (`errors.Is(err, pgx.ErrNoRows)`) im `cdc.heartbeat`-Scan von
// `Diagnose` (`internal/bootstrap/wiring.go`): eine Quelle, die noch nie
// ein Lebenszeichen geschrieben hat, meldet das textlich statt eines
// Lesefehlers — real gegen PostgreSQL (`make test-store`). Anders als
// `Healthcheck` kehrt `Diagnose` hier mit Exit 0 zurück, weil ein
// fehlendes Lebenszeichen Berichtsinhalt ist, kein Befehlsfehler
// (review-slice-038 F-1).
func TestDiagnoseReportsNoHeartbeat(t *testing.T) {
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

	const source = "src-diagnose-no-heartbeat"
	if _, err := pool.Exec(ctx, "DELETE FROM cdc.process_heartbeat WHERE source_id = $1", source); err != nil {
		t.Fatalf("Datenstand-Rückbau: %v", err)
	}

	var code int
	output := captureStdout(t, func() {
		code = bootstrap.Diagnose(ctx, dsn, source)
	})
	if code != 0 {
		t.Fatalf("Diagnose-Exit-Code = %d, wollen 0 (erfolgreicher Lesezugriff, fehlendes Lebenszeichen ist Berichtsinhalt)", code)
	}
	if !strings.Contains(output, "Betriebsstatus (LH-FA-ADM-002): kein Lebenszeichen — Instanz hat noch nie geschlagen") {
		t.Fatalf("stdout = %q, wollen die 'kein Lebenszeichen'-Zeile", output)
	}
	if !strings.Contains(output, "Fehlerzustand (LH-FA-ADM-003): unbekannt (kein Lebenszeichen)") {
		t.Fatalf("stdout = %q, wollen die 'unbekannt (kein Lebenszeichen)'-Zeile", output)
	}
}

const diagnoseNullLagTestSource = "src-diagnose-null-lag"
const diagnoseNullLagTestConsumer = "diagnose-null-lag-test-consumer"

// newDiagnoseNullLagFixture baut den Datenstand für den NULL-sicheren
// Scan-Zweig (`lag == nil` in `internal/bootstrap/wiring.go`): ein
// Consumer mit bestätigter Position, dessen gebundene Quelle noch nie
// eine Transaktion trug. `cdc.consumer_status`s
// `latest_commit_position`-Unterabfrage liefert dafür `NULL`
// (`tools/schema/schema.yaml`), `cdc.metrics.value` für die
// `cdc_consumer_lag`-Zeile dieses Consumers folglich ebenfalls
// (`NULL - x` bleibt `NULL`) — der im Plan-Nachzug slice-038 §3 genannte,
// real beim `test-store`-Lauf entdeckte Bugfix.
func newDiagnoseNullLagFixture(t *testing.T) *pgxpool.Pool {
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

	for _, stmt := range []string{
		"DELETE FROM cdc.consumer_position WHERE consumer_id = $1",
		"DELETE FROM cdc.consumer WHERE consumer_id = $1",
	} {
		if _, err := pool.Exec(ctx, stmt, diagnoseNullLagTestConsumer); err != nil {
			t.Fatalf("Datenstand-Rückbau (%s): %v", stmt, err)
		}
	}
	if _, err := pool.Exec(ctx, "DELETE FROM cdc.process_heartbeat WHERE source_id = $1", diagnoseNullLagTestSource); err != nil {
		t.Fatalf("Datenstand-Rückbau (Heartbeat): %v", err)
	}
	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.source (source_id, name) VALUES ($1, 'Diagnose-NULL-Lag-Testquelle') ON CONFLICT (source_id) DO NOTHING",
		diagnoseNullLagTestSource,
	); err != nil {
		t.Fatalf("Quelle-Zeile: %v", err)
	}
	// Bewusst KEINE cdc.transaction-Zeile für diese Quelle — genau das
	// Fehlen erzeugt den NULL-Fall, den dieser Fixture belegen soll.
	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.consumer (consumer_id, name) VALUES ($1, $1)",
		diagnoseNullLagTestConsumer,
	); err != nil {
		t.Fatalf("Consumer-Zeile: %v", err)
	}
	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.consumer_position (consumer_id, source_id, acknowledged_position) VALUES ($1, $2, 1)",
		diagnoseNullLagTestConsumer, diagnoseNullLagTestSource,
	); err != nil {
		t.Fatalf("Consumer-Position-Zeile: %v", err)
	}
	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.process_heartbeat (source_id, heartbeat_at, error_class) VALUES ($1, current_timestamp, NULL)",
		diagnoseNullLagTestSource,
	); err != nil {
		t.Fatalf("Lebenszeichen-Zeile: %v", err)
	}
	return pool
}

// TestDiagnoseReportsUnknownLagForSourceWithoutTransactions belegt den
// NULL-sicheren Scan über `*float64` für `cdc_consumer_lag`
// (review-slice-038 F-1): ein Consumer mit bestätigter Position, dessen
// Quelle noch nie eine Transaktion trug, meldet „unbekannt“ statt eines
// irreführenden numerischen Rückstands oder eines Scan-Fehlers.
func TestDiagnoseReportsUnknownLagForSourceWithoutTransactions(t *testing.T) {
	newDiagnoseNullLagFixture(t)
	ctx := context.Background()
	dsn := os.Getenv("CDC_STORE_TEST_DSN")

	var code int
	output := captureStdout(t, func() {
		code = bootstrap.Diagnose(ctx, dsn, diagnoseNullLagTestSource)
	})
	if code != 0 {
		t.Fatalf("Diagnose-Exit-Code = %d, wollen 0 (erfolgreicher Lesezugriff)", code)
	}
	if !strings.Contains(output, diagnoseNullLagTestConsumer+": unbekannt (Quelle trug noch nie eine Transaktion)") {
		t.Fatalf("stdout = %q, wollen die NULL-sichere Zeile für %s", output, diagnoseNullLagTestConsumer)
	}
}

// TestDiagnoseReportsNoConfirmedConsumer belegt den `!found`-Zweig von
// `Diagnose` (review-slice-038 F-1): die `cdc_consumer_lag`-Abfrage trägt
// keinen Quellenfilter (`cdc.metrics`, `nacharbeit-observability.sql` —
// reine Aggregation über `cdc.consumer_status`), der Zweig ist deshalb
// nur beobachtbar, wenn system-weit kein Consumer eine bestätigte
// Position trägt. Der Test sichert den bestehenden
// `cdc.consumer_position`-Bestand, leert ihn temporär und stellt ihn über
// `t.Cleanup` wieder her (netzlos gegen den `make test-store`-Testcontainer,
// dessen DB-Daten containerlokal sind) — der einzige Weg, diesen Zweig in
// einem geteilten Paketlauf isoliert herzustellen, ohne einen
// zweiten Container zu brauchen.
func TestDiagnoseReportsNoConfirmedConsumer(t *testing.T) {
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

	const source = "src-diagnose-no-confirmed-consumer"
	if _, err := pool.Exec(ctx, "DELETE FROM cdc.process_heartbeat WHERE source_id = $1", source); err != nil {
		t.Fatalf("Datenstand-Rückbau (Heartbeat): %v", err)
	}
	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.source (source_id, name) VALUES ($1, 'Diagnose-Kein-Consumer-Testquelle') ON CONFLICT (source_id) DO NOTHING",
		source,
	); err != nil {
		t.Fatalf("Quelle-Zeile: %v", err)
	}
	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.process_heartbeat (source_id, heartbeat_at, error_class) VALUES ($1, current_timestamp, NULL)",
		source,
	); err != nil {
		t.Fatalf("Lebenszeichen-Zeile: %v", err)
	}

	type savedPosition struct {
		consumerID string
		sourceID   string
		position   int64
	}
	rows, err := pool.Query(ctx, "SELECT consumer_id, source_id, acknowledged_position FROM cdc.consumer_position")
	if err != nil {
		t.Fatalf("Bestand sichern (Query): %v", err)
	}
	var saved []savedPosition
	for rows.Next() {
		var p savedPosition
		if err := rows.Scan(&p.consumerID, &p.sourceID, &p.position); err != nil {
			rows.Close()
			t.Fatalf("Bestand sichern (Scan): %v", err)
		}
		saved = append(saved, p)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("Bestand sichern (Iteration): %v", err)
	}
	rows.Close()

	if _, err := pool.Exec(ctx, "DELETE FROM cdc.consumer_position"); err != nil {
		t.Fatalf("Bestand temporär leeren: %v", err)
	}
	// Wiederherstellung läuft vor dem Pool-Schluss (t.Cleanup ist LIFO,
	// pool.Close wurde weiter oben registriert) — sonst bleibt der
	// system-weite Bestand für nachfolgende Tests im selben Paketlauf
	// verändert.
	t.Cleanup(func() {
		for _, p := range saved {
			if _, err := pool.Exec(context.Background(),
				"INSERT INTO cdc.consumer_position (consumer_id, source_id, acknowledged_position) VALUES ($1, $2, $3) ON CONFLICT (consumer_id) DO NOTHING",
				p.consumerID, p.sourceID, p.position,
			); err != nil {
				t.Errorf("Bestand wiederherstellen (%s): %v", p.consumerID, err)
			}
		}
	})

	var code int
	output := captureStdout(t, func() {
		code = bootstrap.Diagnose(ctx, dsn, source)
	})
	if code != 0 {
		t.Fatalf("Diagnose-Exit-Code = %d, wollen 0 (erfolgreicher Lesezugriff)", code)
	}
	if !strings.Contains(output, "(keiner — kein Consumer mit bestätigter Position)") {
		t.Fatalf("stdout = %q, wollen die 'keiner'-Zeile bei leerem cdc.consumer_position-Bestand", output)
	}
}
