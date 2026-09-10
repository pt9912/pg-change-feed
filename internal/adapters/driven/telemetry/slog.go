// Package telemetry trägt den `SlogAdapter` als Driven-Implementierung des
// `LogPort` (`ARC-006`, `ADR-0024`): das strukturierte
// Logging-Framework (`log/slog`, Go-Standardbibliothek) bleibt
// Infrastruktur und lebt ausschließlich hier — kein anderer Adapter und
// nicht die Composition Root importieren `log/slog`, um selbst zu
// protokollieren; sie bedienen den Port (`outbound.LogPort`).
package telemetry

import (
	"context"
	"io"
	"log/slog"
	"os"

	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
)

// SlogAdapter implementiert `outbound.LogPort` über einen JSON-Handler
// (`LH-QA-OPS-004`): die vier Port-Stufen reichen unverändert an die
// gleichnamigen `*Context`-Methoden von `*slog.Logger` durch.
type SlogAdapter struct {
	logger *slog.Logger
}

var _ outbound.LogPort = (*SlogAdapter)(nil)

// New baut den JSON-Handler auf `stdout` — der Betriebsbeleg des
// Containers (`compose.yaml`) — mit dem übergebenen Level
// (`bootstrap.parseLogLevel`, `CDC_LOG_LEVEL`).
func New(level slog.Level) *SlogAdapter {
	return newWithWriter(os.Stdout, level)
}

// newWithWriter trägt den Test-Zugang: ein `io.Writer` statt `stdout`
// macht die JSON-Struktur ohne Container-Lauf prüfbar
// (`slog_internal_test.go`).
func newWithWriter(w io.Writer, level slog.Level) *SlogAdapter {
	return &SlogAdapter{logger: slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{Level: level}))}
}

func (a *SlogAdapter) Debug(ctx context.Context, msg string, attrs ...any) {
	a.logger.DebugContext(ctx, msg, attrs...)
}

func (a *SlogAdapter) Info(ctx context.Context, msg string, attrs ...any) {
	a.logger.InfoContext(ctx, msg, attrs...)
}

func (a *SlogAdapter) Warn(ctx context.Context, msg string, attrs ...any) {
	a.logger.WarnContext(ctx, msg, attrs...)
}

func (a *SlogAdapter) Error(ctx context.Context, msg string, attrs ...any) {
	a.logger.ErrorContext(ctx, msg, attrs...)
}
