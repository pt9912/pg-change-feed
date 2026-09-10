package telemetry_test

import (
	"log/slog"
	"testing"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/telemetry"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
)

// TestSlogAdapterImplementsLogPort trägt die Port-Kante (`ARC-004`): `New`
// liefert einen Wert, den die Composition Root ohne weitere Anpassung als
// `outbound.LogPort` an die Driven-/Driving-Adapter injiziert (`WithLog`).
func TestSlogAdapterImplementsLogPort(t *testing.T) {
	var _ outbound.LogPort = telemetry.New(slog.LevelInfo)
}
