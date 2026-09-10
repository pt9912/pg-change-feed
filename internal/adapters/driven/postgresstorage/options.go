package postgresstorage

import "github.com/pt9912/pg-change-feed/internal/application/port/outbound"

// Option konfiguriert einen `postgresstorage`-Adapter bei der Konstruktion
// (`New`, `NewTableActivation`, `NewHeartbeat`, `NewConsumerState`);
// aktuell trägt sie nur den optionalen `LogPort` (`LH-QA-OPS-004`,
// `ADR-0024`). Variadisch statt Pflichtparameter, damit bestehende
// Aufrufstellen (Tests) unverändert kompilieren — ungesetzt bleibt die
// Protokollierung beim No-Op (`outbound.NoopLog`).
type Option func(*options)

type options struct {
	log outbound.LogPort
}

// newOptions trägt den Default (`outbound.NoopLog`) und wendet die
// übergebenen `Option`-Werte in Aufrufreihenfolge an.
func newOptions(opts []Option) options {
	o := options{log: outbound.NoopLog}
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

// WithLog injiziert den `LogPort` (`ADR-0024`): die Composition Root
// übergibt den Telemetrie-Driven-Adapter (`internal/adapters/driven/telemetry`);
// ein Testdouble lässt sich an derselben Stelle einsetzen, ohne
// Paket-globalen Zustand zu berühren.
func WithLog(log outbound.LogPort) Option {
	return func(o *options) { o.log = log }
}
