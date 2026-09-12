package bootstrap_test

import (
	"bytes"
	"context"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pt9912/pg-change-feed/internal/bootstrap"
)

// captureStdout leitet os.Stdout für die Dauer von fn auf einen Puffer um
// und liefert dessen Inhalt zurück — `bootstrap.RegisterConsumer` schreibt
// sein Ergebnis direkt auf `os.Stdout` (dasselbe Umleitungs-Muster wie
// `captureStderr` in `healthcheck_test.go`, dort für die
// Diagnose-Ausgabe auf `os.Stderr`).
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	original := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	os.Stdout = w
	defer func() { os.Stdout = original }()

	fn()

	if err := w.Close(); err != nil {
		t.Fatalf("Pipe schließen: %v", err)
	}
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("Pipe lesen: %v", err)
	}
	return buf.String()
}

// TestRegisterConsumerEndToEnd trägt den neuen Zugriffsweg
// (`LH-FA-CON-001.a`, `cmd/pg-change-feed/main.go` register-consumer)
// gegen eine reale PostgreSQL-Instanz (`make test-store`): der Aufruf
// registriert über `RegisterConsumerUseCase`, nicht über ein
// Direktschreiben der CDC-Speichertabelle. Happy Path und Boundary
// (`LH-FA-CON-001`) laufen im selben Test, weil der zweite Aufruf auf dem
// Stand des ersten aufbaut.
func TestRegisterConsumerEndToEnd(t *testing.T) {
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

	const name = "cli-register-consumer"
	if _, err := pool.Exec(ctx, "DELETE FROM cdc.consumer_position WHERE consumer_id = $1", name); err != nil {
		t.Fatalf("Datenstand-Rückbau (Position): %v", err)
	}
	if _, err := pool.Exec(ctx, "DELETE FROM cdc.consumer WHERE consumer_id = $1", name); err != nil {
		t.Fatalf("Datenstand-Rückbau (Consumer): %v", err)
	}

	var code int
	output := captureStdout(t, func() {
		code = bootstrap.RegisterConsumer(ctx, bootstrap.Config{AdminDSN: dsn}, name)
	})
	if code != 0 {
		t.Fatalf("erste Registrierung: Exit-Code = %d, wollen 0", code)
	}
	if !strings.Contains(output, name) {
		t.Fatalf("stdout = %q, wollen die Consumer-Kennung %q", output, name)
	}

	var count int
	if err := pool.QueryRow(ctx, "SELECT count(*) FROM cdc.consumer WHERE consumer_id = $1", name).Scan(&count); err != nil {
		t.Fatalf("Consumer-Zeilen: %v", err)
	}
	if count != 1 {
		t.Fatalf("Consumer-Zeilen nach Registrierung = %d, wollen 1 (der Aufruf lief über den Use Case, nicht per Direktschreiben)", count)
	}

	// Boundary (`LH-FA-CON-001`): der erneute Aufruf über denselben
	// Zugriffsweg bleibt idempotent und meldet weiterhin Exit-Code 0, mit
	// einer eigenen Ausgabe-Zeile für den bereits registrierten Stand.
	output = captureStdout(t, func() {
		code = bootstrap.RegisterConsumer(ctx, bootstrap.Config{AdminDSN: dsn}, name)
	})
	if code != 0 {
		t.Fatalf("erneute Registrierung: Exit-Code = %d, wollen 0", code)
	}
	if !strings.Contains(output, "bereits registriert") {
		t.Fatalf("stdout = %q, wollen einen Hinweis auf die bereits bestehende Registrierung", output)
	}
	if err := pool.QueryRow(ctx, "SELECT count(*) FROM cdc.consumer WHERE consumer_id = $1", name).Scan(&count); err != nil {
		t.Fatalf("Consumer-Zeilen nach erneuter Registrierung: %v", err)
	}
	if count != 1 {
		t.Fatalf("Consumer-Zeilen nach erneuter Registrierung = %d, wollen weiterhin 1", count)
	}
}

// TestRegisterConsumerReportsStorageFailure trägt die
// Verdrahtungs-Fehlerklasse (Exit-Code 1) — netzlos (`make test`), die
// Verbindung scheitert am geschlossenen lokalen Port, dieselbe
// Diagnose-Form wie `Healthcheck` (`healthcheck_test.go`).
func TestRegisterConsumerReportsStorageFailure(t *testing.T) {
	var code int
	output := captureStderr(t, func() {
		code = bootstrap.RegisterConsumer(context.Background(), bootstrap.Config{AdminDSN: "postgres://x:x@127.0.0.1:1/db?connect_timeout=1"}, "irrelevant")
	})
	if code != 1 {
		t.Fatalf("Exit-Code = %d, wollen 1 (Verdrahtungsfehler)", code)
	}
	if !strings.Contains(output, "register-consumer") {
		t.Fatalf("stderr = %q, wollen eine register-consumer-Diagnose-Zeile", output)
	}
}

// TestRegisterConsumerReportsDomainFailure trägt den zweiten,
// unabhängigen Fehler-Zweig in `RegisterConsumer` — den nach
// `register.Register(...)` (`review-slice-021.md` F-1): eine leere
// Consumer-Kennung scheitert an `model.NewConsumer`
// (`domainerrors.ErrEmptyIdentifier`), bevor der Consumer-State-Port
// berührt wird — real gegen PostgreSQL (`make test-store`), weil der
// Verdrahtungsschritt davor (`postgresstorage.NewConsumerState`) eine
// echte Verbindung braucht, um diesen Zweig überhaupt zu erreichen.
func TestRegisterConsumerReportsDomainFailure(t *testing.T) {
	dsn := os.Getenv("CDC_STORE_TEST_DSN")
	if dsn == "" {
		t.Skip("CDC_STORE_TEST_DSN nicht gesetzt — reale PostgreSQL-Tests laufen über make test-store")
	}

	var code int
	output := captureStderr(t, func() {
		code = bootstrap.RegisterConsumer(context.Background(), bootstrap.Config{AdminDSN: dsn}, "")
	})
	if code != 1 {
		t.Fatalf("Exit-Code = %d, wollen 1 (Domänenfehler)", code)
	}
	if !strings.Contains(output, "register-consumer") {
		t.Fatalf("stderr = %q, wollen eine register-consumer-Diagnose-Zeile", output)
	}
}
