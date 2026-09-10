package telemetry

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
)

// TestNewWithWriterWritesStructuredJSON trägt den Beleg für
// `LH-QA-OPS-004`/`ADR-0024`, den slice-014 bisher nur am
// Integrationstest-stdout-Diff hatte (Review-Finding F-2): die Ausgabe von
// `New` (über den paketinternen Test-Zugang `newWithWriter`, derselbe
// Handler-Aufbau wie am `stdout`-Pfad) ist gültiges JSON mit `msg`,
// `level` und den übergebenen Attributen.
func TestNewWithWriterWritesStructuredJSON(t *testing.T) {
	var buf bytes.Buffer
	adapter := newWithWriter(&buf, slog.LevelDebug)

	adapter.Info(context.Background(), "changestore: verbunden", "source", "src-1")

	line := strings.TrimSpace(buf.String())
	var decoded map[string]any
	if err := json.Unmarshal([]byte(line), &decoded); err != nil {
		t.Fatalf("Log-Zeile ist kein gültiges JSON: %v (%s)", err, line)
	}
	if decoded["msg"] != "changestore: verbunden" {
		t.Fatalf("msg-Feld: %v", decoded["msg"])
	}
	if decoded["level"] != "INFO" {
		t.Fatalf("level-Feld: %v", decoded["level"])
	}
	if decoded["source"] != "src-1" {
		t.Fatalf("source-Feld: %v", decoded["source"])
	}
}

// TestNewWithWriterFiltersBelowLevel trägt die Level-Filterung: `Debug`
// bleibt unterhalb des konfigurierten `Info`-Levels stumm — dieselbe
// Semantik, die `CDC_LOG_LEVEL` am Container-stdout steuert.
func TestNewWithWriterFiltersBelowLevel(t *testing.T) {
	var buf bytes.Buffer
	adapter := newWithWriter(&buf, slog.LevelInfo)

	adapter.Debug(context.Background(), "heartbeat: Lebenszeichen geschrieben")

	if buf.Len() != 0 {
		t.Fatalf("Debug-Zeile trotz Info-Level geschrieben: %s", buf.String())
	}
}
