package telemetry

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
)

// TestNewWithWriterWarnAndErrorCarryTheirLevel trägt den Vertrag des
// `LogPort` für die beiden Stufen, die der Adapter über den
// JSON-Handler-Standard hinausreichen muss (`LH-QA-OPS-004`,
// `ARC-006`, `ADR-0024`): `Warn`/`Error` rufen genau die gleichnamige
// `*Context`-Methode ihres Levels auf, der JSON-Handler setzt `level` auf
// `WARN`/`ERROR` und trägt die übergebenen Attribute unverändert mit.
//
// Die zwei Stufen sind kein Detail: `runWALRetentionCheck` protokolliert
// den Warn-Bereich („kontrollierte Fortsetzung") über `Warn`, den
// Abbruch über `Error` (`SPEC-008`, `ADR-0049`) — eine Vertauschung der
// beiden Methodenrümpfe ließe den Betreiber die beiden Lagen nicht mehr
// unterscheiden, ohne den Testlauf zu stören.
//
// Rot färbende Mutation: den Rumpf von `Warn` auf
// `a.logger.InfoContext(ctx, msg, attrs...)` umstellen — dann trägt die
// Warn-Zeile `level` `INFO` statt `WARN` und dieser Test bricht. Dieselbe
// Probe greift für `Error` gegen `Warn`.
func TestNewWithWriterWarnAndErrorCarryTheirLevel(t *testing.T) {
	cases := []struct {
		name      string
		log       func(adapter *SlogAdapter, ctx context.Context, msg string, attrs ...any)
		wantLevel string
	}{
		{
			name:      "Warn",
			log:       func(a *SlogAdapter, ctx context.Context, msg string, attrs ...any) { a.Warn(ctx, msg, attrs...) },
			wantLevel: "WARN",
		},
		{
			name:      "Error",
			log:       func(a *SlogAdapter, ctx context.Context, msg string, attrs ...any) { a.Error(ctx, msg, attrs...) },
			wantLevel: "ERROR",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var buf bytes.Buffer
			// LevelWarn als Handler-Level: beide geprüften Stufen liegen auf
			// oder über der Konfiguration und werden deshalb real geschrieben.
			adapter := newWithWriter(&buf, slog.LevelWarn)

			c.log(adapter, context.Background(), "replication: Lage", "metric", "cdc_wal_retention_bytes", "bytes", 4096)

			line := strings.TrimSpace(buf.String())
			var decoded map[string]any
			if err := json.Unmarshal([]byte(line), &decoded); err != nil {
				t.Fatalf("Log-Zeile ist kein gültiges JSON: %v (%s)", err, line)
			}
			if decoded["level"] != c.wantLevel {
				t.Fatalf("level-Feld: %v, wollen %s", decoded["level"], c.wantLevel)
			}
			if decoded["msg"] != "replication: Lage" {
				t.Fatalf("msg-Feld: %v", decoded["msg"])
			}
			if decoded["metric"] != "cdc_wal_retention_bytes" {
				t.Fatalf("Attribut metric: %v", decoded["metric"])
			}
			if decoded["bytes"] != float64(4096) {
				t.Fatalf("Attribut bytes: %v", decoded["bytes"])
			}
		})
	}
}
