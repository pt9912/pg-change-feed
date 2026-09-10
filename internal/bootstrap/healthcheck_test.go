package bootstrap_test

import (
	"bytes"
	"context"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/pt9912/pg-change-feed/internal/bootstrap"
)

// captureStderr leitet os.Stderr für die Dauer von fn auf einen Puffer um
// und liefert dessen Inhalt zurück — `bootstrap.Healthcheck` schreibt
// seine Diagnose direkt auf `os.Stderr` (Review-Finding F-2,
// `docs/reviews/review-slice-012.md`: der Exit-Code allein unterscheidet
// die drei Fehlerklassen nicht, die Diagnose-Zeile schon).
func captureStderr(t *testing.T, fn func()) string {
	t.Helper()
	original := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	os.Stderr = w
	defer func() { os.Stderr = original }()

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

// TestHealthcheckReportsConnectionFailure belegt die erste Fehlerklasse
// aus F-2: eine nicht erreichbare Instanz trägt eine eigene,
// unterscheidbare stderr-Zeile — netzlos (`make test`), die Verbindung
// scheitert am geschlossenen lokalen Port.
func TestHealthcheckReportsConnectionFailure(t *testing.T) {
	var code int
	output := captureStderr(t, func() {
		code = bootstrap.Healthcheck(context.Background(), "postgres://x:x@127.0.0.1:1/db?connect_timeout=1", "src-1")
	})
	if code != 1 {
		t.Fatalf("Healthcheck-Exit-Code = %d, wollen 1 (Verbindungsfehler)", code)
	}
	if !strings.Contains(output, "nicht erreichbar") {
		t.Fatalf("stderr = %q, wollen eine Diagnose-Zeile mit dem Text 'nicht erreichbar'", output)
	}
}

// TestHealthcheckReportsMissingHeartbeatRow belegt die zweite
// Fehlerklasse aus F-2: eine Quelle ohne Lebenszeichen-Zeile trägt eine
// andere Diagnose-Zeile als der Verbindungsfehler — gegen die reale
// Instanz (`make test-store`), `cdc.heartbeat` existiert dort bereits
// (Schema-Rollout des Runners); ohne DSN überspringt der Test.
func TestHealthcheckReportsMissingHeartbeatRow(t *testing.T) {
	dsn := os.Getenv("CDC_STORE_TEST_DSN")
	if dsn == "" {
		t.Skip("CDC_STORE_TEST_DSN nicht gesetzt — reale PostgreSQL-Tests laufen über make test-store")
	}
	var code int
	output := captureStderr(t, func() {
		code = bootstrap.Healthcheck(context.Background(), dsn, "src-healthcheck-never-beaten")
	})
	if code != 1 {
		t.Fatalf("Healthcheck-Exit-Code = %d, wollen 1 (kein Lebenszeichen)", code)
	}
	if !strings.Contains(output, "kein Lebenszeichen") {
		t.Fatalf("stderr = %q, wollen eine Diagnose-Zeile mit dem Text 'kein Lebenszeichen'", output)
	}
	if strings.Contains(output, "nicht erreichbar") {
		t.Fatalf("stderr = %q, trägt die Verbindungsfehler-Diagnose statt der Zeilen-Diagnose", output)
	}
}
