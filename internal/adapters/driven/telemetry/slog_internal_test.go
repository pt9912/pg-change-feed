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

// TestNewWithWriterFiltersBelowLevel trägt die Level-Filterung: eine Stufe
// unterhalb des konfigurierten Levels bleibt stumm — dieselbe Semantik, die
// `CDC_LOG_LEVEL` am Container-stdout steuert. Zwei Fälle, nicht nur einer
// (Verifier-Fund V-1, `docs/reviews/verify-slice-014.md`): `slog.
// HandlerOptions{}` defaultet ein unbesetztes `Level`-Feld selbst auf
// `LevelInfo` — ein Test, der ausschließlich mit `LevelInfo` konfiguriert,
// bleibt grün, selbst wenn `newWithWriter` das `Level:`-Feld gar nicht mehr
// an die `HandlerOptions` durchreicht (die Mutation, die V-1 beschreibt).
// Der zweite Fall (`LevelWarn`) liegt oberhalb dieses Zufalls-Defaults:
// ohne durchgereichten Level fiele die `Info`-Zeile *nicht* unter den
// (dann wirkungslosen) Default und der Test schlägt fehl — genau die
// Mutation, die der erste Fall allein nicht fängt.
func TestNewWithWriterFiltersBelowLevel(t *testing.T) {
	for _, testcase := range []struct {
		name       string
		configured slog.Level
	}{
		{name: "Info-Level filtert Debug", configured: slog.LevelInfo},
		{name: "Warn-Level filtert Info (Zufalls-Default-Probe V-1)", configured: slog.LevelWarn},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			var buf bytes.Buffer
			adapter := newWithWriter(&buf, testcase.configured)

			if testcase.configured == slog.LevelInfo {
				adapter.Debug(context.Background(), "heartbeat: Lebenszeichen geschrieben")
			} else {
				adapter.Info(context.Background(), "changestore: verbunden")
			}

			if buf.Len() != 0 {
				t.Fatalf("Zeile trotz konfiguriertem Level %s geschrieben: %s", testcase.configured, buf.String())
			}
		})
	}
}
