package outbound

import "context"

// LogPort trägt strukturiertes Betriebs-Logging als Fähigkeit (`ARC-004`,
// `ADR-0024`: „Logging-/Metrics-Frameworks bleiben Infrastruktur. …
// werden durch Driven Adapters implementiert.", geschärft durch `ARC-011`
// „Metriken und strukturierte Logs … über den Outbound Port
// substituierbar"). Die vier Stufen spiegeln `log/slog`
// (Go-Standardbibliothek, `LH-QA-OPS-004`), damit der Driven Adapter sie
// ohne Übersetzung durchreicht — der Port selbst bindet keine konkrete
// Logging-Bibliothek; austauschbar bliebe er auch für ein anderes Backend.
type LogPort interface {
	Debug(ctx context.Context, msg string, attrs ...any)
	Info(ctx context.Context, msg string, attrs ...any)
	Warn(ctx context.Context, msg string, attrs ...any)
	Error(ctx context.Context, msg string, attrs ...any)
}

// NoopLog trägt den Default für Aufrufer, die keinen Telemetrie-Adapter
// injizieren (Tests, `postgresstorage.Option`/`postgresack.Option`
// unbenutzt) — vermeidet einen Nil-Check an jeder Logging-Aufrufstelle der
// Driven-/Driving-Adapter.
var NoopLog LogPort = noopLog{}

type noopLog struct{}

func (noopLog) Debug(context.Context, string, ...any) {}
func (noopLog) Info(context.Context, string, ...any)  {}
func (noopLog) Warn(context.Context, string, ...any)  {}
func (noopLog) Error(context.Context, string, ...any) {}
