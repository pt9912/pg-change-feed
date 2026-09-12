package bootstrap_test

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pt9912/pg-change-feed/internal/bootstrap"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// TestAcknowledgeConsumerEndToEnd trägt den neuen Zugriffsweg
// (`LH-FA-CON-004.a`, `cmd/pg-change-feed/main.go` acknowledge-consumer)
// gegen eine reale PostgreSQL-Instanz (`make test-store`): der Aufruf
// bestätigt über `AcknowledgeConsumerUseCase`, nicht über ein
// Direktschreiben der Consumer-Positions-Tabelle. Happy Path, die
// Wiederholung derselben Position (Idempotenz, `LH-FA-CON-004` Boundary)
// und der echte Rückschritt (`ErrPositionRegression`) laufen im selben
// Test, weil jeder Schritt auf dem Stand des vorigen aufbaut — genau der
// Vorwärts-Invariante-Nachweis, den `LH-FA-CON-004.a` über diesen
// Zugriffsweg verlangt.
func TestAcknowledgeConsumerEndToEnd(t *testing.T) {
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

	const name = "cli-acknowledge-consumer"
	const source model.SourceID = "cli-acknowledge-source"
	if _, err := pool.Exec(ctx, "DELETE FROM cdc.consumer_position WHERE consumer_id = $1", name); err != nil {
		t.Fatalf("Datenstand-Rückbau (Position): %v", err)
	}
	if _, err := pool.Exec(ctx, "DELETE FROM cdc.consumer WHERE consumer_id = $1", name); err != nil {
		t.Fatalf("Datenstand-Rückbau (Consumer): %v", err)
	}
	// Die Positions-Zeile trägt den Fremdschlüssel auf `cdc.source`
	// (`SPEC-001`) — dieselbe Vorbedingung wie bei den
	// Adapter-Tests (`consumerstate_test.go`).
	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.source (source_id, name) VALUES ($1, 'Acknowledge-Quelle') ON CONFLICT (source_id) DO NOTHING",
		string(source),
	); err != nil {
		t.Fatalf("Quelle-Zeile: %v", err)
	}

	cfg := bootstrap.Config{DSN: dsn, Source: source}
	if code := bootstrap.RegisterConsumer(ctx, cfg, name); code != 0 {
		t.Fatalf("Consumer-Registrierung (Vorbedingung): Exit-Code = %d, wollen 0", code)
	}

	// Happy Path: die erste Bestätigung rückt die Position vor
	// (`LH-FA-CON-004.a`).
	var code int
	output := captureStdout(t, func() {
		code = bootstrap.AcknowledgeConsumer(ctx, cfg, name, 100)
	})
	if code != 0 {
		t.Fatalf("erste Bestätigung: Exit-Code = %d, wollen 0", code)
	}
	if !strings.Contains(output, name) {
		t.Fatalf("stdout = %q, wollen die Consumer-Kennung %q", output, name)
	}
	assertStoredOffset(t, pool, name, string(source), 100)

	// Boundary (`LH-FA-CON-004`): die Wiederholung derselben Position
	// bleibt idempotent — Exit-Code 0, gespeicherter Stand unverändert.
	output = captureStdout(t, func() {
		code = bootstrap.AcknowledgeConsumer(ctx, cfg, name, 100)
	})
	if code != 0 {
		t.Fatalf("wiederholte Bestätigung: Exit-Code = %d, wollen 0 (Idempotenz)", code)
	}
	assertStoredOffset(t, pool, name, string(source), 100)

	// Vorwärts-Invariante (`ErrPositionRegression`, `LH-FA-CON-004.a`):
	// ein echter Rückschritt über denselben Zugriffsweg wird abgelehnt,
	// der gespeicherte Stand bleibt bei 100 stehen.
	errOutput := captureStderr(t, func() {
		code = bootstrap.AcknowledgeConsumer(ctx, cfg, name, 50)
	})
	if code != 1 {
		t.Fatalf("rückläufige Bestätigung: Exit-Code = %d, wollen 1 (Vorwärts-Invariante)", code)
	}
	if !strings.Contains(errOutput, "acknowledge-consumer") {
		t.Fatalf("stderr = %q, wollen eine acknowledge-consumer-Diagnose-Zeile", errOutput)
	}
	if !strings.Contains(errOutput, "vor") {
		t.Fatalf("stderr = %q, wollen einen Hinweis auf die Vorwärts-Invariante (ErrPositionRegression)", errOutput)
	}
	assertStoredOffset(t, pool, name, string(source), 100)

	// Fortlaufende Bestätigung nach der abgewiesenen: der Zugriffsweg
	// bleibt für den regulären Vorwärtsfall offen.
	if code := bootstrap.AcknowledgeConsumer(ctx, cfg, name, 300); code != 0 {
		t.Fatalf("fortlaufende Bestätigung: Exit-Code = %d, wollen 0", code)
	}
	assertStoredOffset(t, pool, name, string(source), 300)
}

// assertStoredOffset liest die gespeicherte Position direkt aus
// `cdc.consumer_position` — der Beleg, dass der Zugriffsweg tatsächlich
// über den Use Case geschrieben hat, nicht nur den Exit-Code meldet.
func assertStoredOffset(t *testing.T, pool *pgxpool.Pool, consumer, wantSource string, wantOffset int64) {
	t.Helper()
	var gotSource string
	var gotOffset int64
	if err := pool.QueryRow(context.Background(),
		"SELECT source_id, acknowledged_position FROM cdc.consumer_position WHERE consumer_id = $1",
		consumer,
	).Scan(&gotSource, &gotOffset); err != nil {
		t.Fatalf("gespeicherte Position lesen: %v", err)
	}
	if gotSource != wantSource || gotOffset != wantOffset {
		t.Fatalf("gespeicherte Position = (%q, %d), wollen (%q, %d)", gotSource, gotOffset, wantSource, wantOffset)
	}
}

// TestAcknowledgeConsumerReportsUnregistered trägt die
// Registrierungs-Grenze (`outbound.ErrConsumerUnregistered`) über
// denselben Zugriffsweg: ein ACK ohne registrierte Kennung endet
// sichtbar mit Exit-Code 1, ohne eine Positions-Zeile anzulegen — real
// gegen PostgreSQL (`make test-store`), weil der Verdrahtungsschritt
// davor (`postgresstorage.NewConsumerState`) eine echte Verbindung
// braucht.
func TestAcknowledgeConsumerReportsUnregistered(t *testing.T) {
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

	const name = "cli-acknowledge-unregistered"
	if _, err := pool.Exec(ctx, "DELETE FROM cdc.consumer_position WHERE consumer_id = $1", name); err != nil {
		t.Fatalf("Datenstand-Rückbau (Position): %v", err)
	}
	if _, err := pool.Exec(ctx, "DELETE FROM cdc.consumer WHERE consumer_id = $1", name); err != nil {
		t.Fatalf("Datenstand-Rückbau (Consumer): %v", err)
	}

	var code int
	output := captureStderr(t, func() {
		code = bootstrap.AcknowledgeConsumer(ctx, bootstrap.Config{DSN: dsn, Source: "cli-acknowledge-source"}, name, 100)
	})
	if code != 1 {
		t.Fatalf("Exit-Code = %d, wollen 1 (nicht registriert)", code)
	}
	if !strings.Contains(output, "acknowledge-consumer") {
		t.Fatalf("stderr = %q, wollen eine acknowledge-consumer-Diagnose-Zeile", output)
	}
	var count int
	if err := pool.QueryRow(ctx, "SELECT count(*) FROM cdc.consumer_position WHERE consumer_id = $1", name).Scan(&count); err != nil {
		t.Fatalf("Positions-Zeilen: %v", err)
	}
	if count != 0 {
		t.Fatalf("Positions-Zeilen nach abgewiesenem ACK = %d, wollen 0", count)
	}
}

// TestAcknowledgeConsumerReportsInvalidPosition trägt einen der zwei
// externen Domänenfehler-Pfade aus `review-slice-022.md` F-1: ein Offset
// von 0 scheitert an `model.NewSourcePosition`
// (`domainerrors.ErrInvalidPosition`), bevor der Consumer-State-Port
// berührt wird — real gegen PostgreSQL (`make test-store`), weil der
// Verdrahtungsschritt davor (`postgresstorage.NewConsumerState`) eine
// echte Verbindung braucht, um diesen Zweig überhaupt zu erreichen
// (dasselbe Muster wie `TestRegisterConsumerReportsDomainFailure` in
// `register_test.go`).
func TestAcknowledgeConsumerReportsInvalidPosition(t *testing.T) {
	dsn := os.Getenv("CDC_STORE_TEST_DSN")
	if dsn == "" {
		t.Skip("CDC_STORE_TEST_DSN nicht gesetzt — reale PostgreSQL-Tests laufen über make test-store")
	}

	var code int
	output := captureStderr(t, func() {
		code = bootstrap.AcknowledgeConsumer(context.Background(),
			bootstrap.Config{DSN: dsn, Source: "cli-acknowledge-source"},
			"cli-acknowledge-invalid-position", 0)
	})
	if code != 1 {
		t.Fatalf("Exit-Code = %d, wollen 1 (Domänenfehler: Offset ohne Wert)", code)
	}
	if !strings.Contains(output, "acknowledge-consumer") {
		t.Fatalf("stderr = %q, wollen eine acknowledge-consumer-Diagnose-Zeile", output)
	}
	if !strings.Contains(output, "Offset") {
		t.Fatalf("stderr = %q, wollen einen Hinweis auf die fehlende Offset-Invariante (ErrInvalidPosition)", output)
	}
}

// TestAcknowledgeConsumerReportsEmptyIdentifier trägt den zweiten der
// zwei externen Domänenfehler-Pfade aus `review-slice-022.md` F-1: eine
// leere Consumer-Kennung scheitert an `AcknowledgeConsumerService.Acknowledge`
// (`domainerrors.ErrEmptyIdentifier`), noch vor der Positions-Prüfung und
// bevor der Consumer-State-Port berührt wird — real gegen PostgreSQL
// (`make test-store`), aus demselben Verdrahtungsgrund wie oben.
func TestAcknowledgeConsumerReportsEmptyIdentifier(t *testing.T) {
	dsn := os.Getenv("CDC_STORE_TEST_DSN")
	if dsn == "" {
		t.Skip("CDC_STORE_TEST_DSN nicht gesetzt — reale PostgreSQL-Tests laufen über make test-store")
	}

	var code int
	output := captureStderr(t, func() {
		code = bootstrap.AcknowledgeConsumer(context.Background(),
			bootstrap.Config{DSN: dsn, Source: "cli-acknowledge-source"},
			"", 100)
	})
	if code != 1 {
		t.Fatalf("Exit-Code = %d, wollen 1 (Domänenfehler: leere Kennung)", code)
	}
	if !strings.Contains(output, "acknowledge-consumer") {
		t.Fatalf("stderr = %q, wollen eine acknowledge-consumer-Diagnose-Zeile", output)
	}
	if !strings.Contains(output, "Kennung") {
		t.Fatalf("stderr = %q, wollen einen Hinweis auf die leere Kennung (ErrEmptyIdentifier)", output)
	}
}

// TestAcknowledgeConsumerReportsStorageFailure trägt die
// Verdrahtungs-Fehlerklasse (Exit-Code 1) — netzlos (`make test`), die
// Verbindung scheitert am geschlossenen lokalen Port, dieselbe
// Diagnose-Form wie `RegisterConsumer` (`register_test.go`).
func TestAcknowledgeConsumerReportsStorageFailure(t *testing.T) {
	var code int
	output := captureStderr(t, func() {
		code = bootstrap.AcknowledgeConsumer(context.Background(), bootstrap.Config{DSN: "postgres://x:x@127.0.0.1:1/db?connect_timeout=1"}, "irrelevant", 1)
	})
	if code != 1 {
		t.Fatalf("Exit-Code = %d, wollen 1 (Verdrahtungsfehler)", code)
	}
	if !strings.Contains(output, "acknowledge-consumer") {
		t.Fatalf("stderr = %q, wollen eine acknowledge-consumer-Diagnose-Zeile", output)
	}
}
