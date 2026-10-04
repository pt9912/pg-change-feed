// otlp.go trägt die Verdrahtung des OTLP-Metrik-Exports (`ADR-0149`): die
// Prüfung der drei Konfigurationswerte, die beide Zugriffswege der
// Konfiguration rufen, den Halter des zuletzt gemessenen WAL-Rückstands und
// den Takt, der den Export-Use-Case in eigener Goroutine ausführt.
package bootstrap

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/otlpexport"
	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage"
	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/exportmetrics"
	"github.com/pt9912/pg-change-feed/internal/domain/messagecode"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

const (
	// otlpDefaultInterval ist der Takt ohne Angabe.
	otlpDefaultInterval = 60 * time.Second
	// otlpMinIntervalSeconds und otlpMaxIntervalSeconds begrenzen den Takt.
	otlpMinIntervalSeconds = 5
	otlpMaxIntervalSeconds = 3600
	// otlpMaxAttemptTimeout begrenzt die Frist je Versuch nach oben.
	otlpMaxAttemptTimeout = 10 * time.Second
)

// ErrOTLPEndpointInvalid, ErrOTLPHeadersInvalid, ErrOTLPIntervalInvalid und
// ErrOTLPHeadersWithoutEndpoint tragen die Fehlerklasse `configuration` der
// Export-Konfiguration.
var (
	ErrOTLPEndpointInvalid        = messagecode.New(messagecode.OTLPEndpointInvalid, "OTLP-Endpunkt ungültig")
	ErrOTLPHeadersInvalid         = messagecode.New(messagecode.OTLPHeadersInvalid, "OTLP-Header ungültig")
	ErrOTLPIntervalInvalid        = messagecode.New(messagecode.OTLPIntervalInvalid, "OTLP-Takt ungültig")
	ErrOTLPHeadersWithoutEndpoint = messagecode.New(messagecode.OTLPHeadersNoEndpoint, "OTLP-Header ohne Endpunkt")
)

// otlpConfigError ist der Fehler der Export-Konfiguration: sein Text trägt den
// Code des Sentinels samt Einzelheit, und er ist zugleich ein
// `ErrConfiguration` (Start-Hindernis der Verdrahtung). Die Einzelheit nennt
// Variable und Position, nie einen Header-Wert oder den Benutzerteil der URL.
type otlpConfigError struct {
	sentinel error
	detail   string
}

func (e *otlpConfigError) Error() string { return e.sentinel.Error() + ": " + e.detail }

func (e *otlpConfigError) Unwrap() []error { return []error{e.sentinel, ErrConfiguration} }

// applyOTLP liest die Export-Konfiguration in `cfg` und prüft sie. Beide
// Zugriffswege (`ConfigFromEnv`, `mergeConfig`) rufen diese eine Funktion.
// Endpunkt und Header kommen allein aus der Umgebung (sie können
// Zugangsdaten tragen); den Takt liefert die Umgebung vor dem Datei-Feld
// `otlp_interval` (`fileInterval`, leer: nicht gesetzt) vor dem Default. Eine
// leere Umgebungsvariable gilt als nicht gesetzt.
func applyOTLP(cfg *Config, getenv func(string) string, fileInterval string) error {
	cfg.OTLPEndpoint = getenv(envOTLPEndpoint)
	if cfg.OTLPEndpoint != "" {
		if err := validateOTLPEndpoint(cfg.OTLPEndpoint); err != nil {
			return err
		}
	}
	headersRaw := getenv(envOTLPHeaders)
	headers, err := parseOTLPHeaders(headersRaw)
	if err != nil {
		return err
	}
	cfg.OTLPHeaders = headers

	name, intervalRaw := envOTLPInterval, getenv(envOTLPInterval)
	if intervalRaw == "" {
		name, intervalRaw = "otlp_interval", fileInterval
	}
	cfg.OTLPInterval = otlpDefaultInterval
	if intervalRaw != "" {
		interval, err := parseOTLPInterval(name, intervalRaw)
		if err != nil {
			return err
		}
		cfg.OTLPInterval = interval
	}

	if headersRaw != "" && cfg.OTLPEndpoint == "" {
		return &otlpConfigError{sentinel: ErrOTLPHeadersWithoutEndpoint,
			detail: fmt.Sprintf("%s gesetzt, aber %s fehlt", envOTLPHeaders, envOTLPEndpoint)}
	}
	return nil
}

// validateOTLPEndpoint verlangt das Schema `http` oder `https` und einen Host.
func validateOTLPEndpoint(raw string) error {
	u, err := url.Parse(raw)
	switch {
	case err != nil:
		return &otlpConfigError{sentinel: ErrOTLPEndpointInvalid, detail: envOTLPEndpoint + " ist nicht als URL lesbar"}
	case u.Scheme != "http" && u.Scheme != "https":
		return &otlpConfigError{sentinel: ErrOTLPEndpointInvalid, detail: envOTLPEndpoint + ": das Schema muss http oder https sein"}
	case u.Hostname() == "":
		return &otlpConfigError{sentinel: ErrOTLPEndpointInvalid, detail: envOTLPEndpoint + ": der Host fehlt"}
	}
	return nil
}

// parseOTLPInterval liest den Takt als ganze Zahl von 5 bis 3600 Sekunden.
func parseOTLPInterval(name, raw string) (time.Duration, error) {
	seconds, err := strconv.ParseUint(raw, 10, 32)
	if err != nil || seconds < otlpMinIntervalSeconds || seconds > otlpMaxIntervalSeconds {
		return 0, &otlpConfigError{sentinel: ErrOTLPIntervalInvalid,
			detail: fmt.Sprintf("%s: %q ist keine ganze Zahl von %d bis %d", name, raw, otlpMinIntervalSeconds, otlpMaxIntervalSeconds)}
	}
	return time.Duration(seconds) * time.Second, nil
}

// parseOTLPHeaders zerlegt `k=v,k2=v2`: das Komma trennt die Elemente, das
// erste `=` Schlüssel und Wert (der Wert darf weitere `=` tragen). Ein Element
// ohne `=`, mit leerem Schlüssel, mit Leerraum oder einem unzulässigen Zeichen
// im Schlüssel oder mit einem Steuerzeichen im Wert ist ein Fehler; Leerraum im
// Wert bleibt erhalten. Eine leere Zeichenkette liefert keine Header. Der Text
// des Fehlers nennt die Position des Elements, nie dessen Inhalt.
func parseOTLPHeaders(raw string) ([]otlpexport.Header, error) {
	if raw == "" {
		return nil, nil
	}
	elements := strings.Split(raw, ",")
	headers := make([]otlpexport.Header, 0, len(elements))
	for i, element := range elements {
		key, value, found := strings.Cut(element, "=")
		fail := func(reason string) error {
			return &otlpConfigError{sentinel: ErrOTLPHeadersInvalid, detail: fmt.Sprintf("%s: Element %d %s", envOTLPHeaders, i+1, reason)}
		}
		switch {
		case !found:
			return nil, fail("trägt kein =")
		case key == "":
			return nil, fail("hat einen leeren Schlüssel")
		case strings.IndexFunc(key, unicode.IsSpace) >= 0:
			return nil, fail("hat Leerraum im Schlüssel")
		case strings.IndexFunc(key, func(r rune) bool { return !isHeaderTokenRune(r) }) >= 0:
			return nil, fail("hat ein unzulässiges Zeichen im Schlüssel")
		case strings.IndexFunc(value, isHeaderControlRune) >= 0:
			return nil, fail("hat ein Steuerzeichen im Wert")
		}
		headers = append(headers, otlpexport.Header{Key: key, Value: value})
	}
	return headers, nil
}

// isHeaderTokenRune meldet, ob das Zeichen in einem Header-Namen stehen darf
// (RFC 9110 `token`).
func isHeaderTokenRune(r rune) bool {
	switch {
	case r >= '0' && r <= '9', r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z':
		return true
	}
	return strings.ContainsRune("!#$%&'*+-.^_`|~", r)
}

// isHeaderControlRune meldet ein Steuerzeichen, das ein Header-Wert nicht
// tragen kann (der Tabulator ist zulässig).
func isHeaderControlRune(r rune) bool { return (r < 0x20 && r != '\t') || r == 0x7f }

// otlpAttemptTimeout ist die Frist je Versuch: der Takt, höchstens zehn
// Sekunden.
func otlpAttemptTimeout(interval time.Duration) time.Duration {
	if interval < otlpMaxAttemptTimeout {
		return interval
	}
	return otlpMaxAttemptTimeout
}

// walRetentionGauge hält den zuletzt gemessenen WAL-Rückstand und seinen
// Messzeitpunkt. Der Prüfzug des Prozesses schreibt ihn (über
// `recordingWALMeasurer`), der Export liest ihn (`outbound.WALRetentionSource`).
type walRetentionGauge struct {
	mu    sync.Mutex
	bytes int64
	at    model.TimePoint
	ok    bool
}

var _ outbound.WALRetentionSource = (*walRetentionGauge)(nil)

func (g *walRetentionGauge) set(bytes int64, at model.TimePoint) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.bytes, g.at, g.ok = bytes, at, true
}

// LastWALRetention liefert den letzten Wert; `ok` ist falsch, solange keiner
// gemessen ist.
func (g *walRetentionGauge) LastWALRetention() (int64, model.TimePoint, bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.bytes, g.at, g.ok
}

// recordingWALMeasurer reicht jede Messung unverändert durch und hält jede
// erfolgreiche im Halter fest; Schwellen-Vergleich und Abbruchpfade von
// `runWALRetentionCheck` bleiben unberührt.
type recordingWALMeasurer struct {
	inner walRetentionMeasurer
	gauge *walRetentionGauge
	clock outbound.ClockPort
}

func (m recordingWALMeasurer) Measure(ctx context.Context) (int64, error) {
	bytes, err := m.inner.Measure(ctx)
	if err == nil {
		m.gauge.set(bytes, m.clock.Now())
	}
	return bytes, err
}

// runMetricExport führt den Export-Use-Case bis zum Ende des Kontexts im
// Takt `interval` aus; der erste Versuch folgt nach einem vollen Takt. Ein
// Zyklus läuft höchstens bis zu seiner Frist und blockiert nur diese
// Goroutine; ein Fehlschlag ist im Use Case protokolliert und wird nicht
// wiederholt, der nächste Takt ist die Wiederholung.
func runMetricExport(ctx context.Context, useCase inbound.ExportMetricsUseCase, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	runMetricExportTicks(ctx, useCase, ticker.C)
}

func runMetricExportTicks(ctx context.Context, useCase inbound.ExportMetricsUseCase, ticks <-chan time.Time) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticks:
			_ = useCase.Export(ctx)
		}
	}
}

// startMetricExport startet den Export, wenn `cfg.OTLPEndpoint` gesetzt ist:
// Lese-Adapter unter `cfg.ReaderDSN`, OTLP-Adapter, Use Case und die eigene
// Goroutine des Takts. Ist kein Endpunkt gesetzt, entsteht nichts (kein
// Adapter, keine Verbindung, keine Goroutine) und `started` ist falsch. Die
// zurückgegebene Funktion beendet die Goroutine, wartet auf sie und schließt
// den Lese-Pool; sie ist in jedem Fall aufrufbar.
func startMetricExport(ctx context.Context, cfg Config, log outbound.LogPort, clock outbound.ClockPort, wal outbound.WALRetentionSource) (stop func(), started bool, err error) {
	if cfg.OTLPEndpoint == "" {
		return func() {}, false, nil
	}
	reader, err := postgresstorage.NewMetrics(ctx, cfg.ReaderDSN, postgresstorage.WithLog(log))
	if err != nil {
		return nil, false, err
	}
	stopLoop, err := startMetricExportLoop(ctx, cfg, reader, log, clock, wal)
	if err != nil {
		reader.Close()
		return nil, false, err
	}
	return func() {
		stopLoop()
		reader.Close()
	}, true, nil
}

// startMetricExportLoop baut OTLP-Adapter und Use Case über dem übergebenen
// Lese-Port und startet die Goroutine des Takts; die zurückgegebene Funktion
// beendet sie und wartet auf sie.
func startMetricExportLoop(ctx context.Context, cfg Config, reader outbound.MetricsReadPort, log outbound.LogPort, clock outbound.ClockPort, wal outbound.WALRetentionSource) (stop func(), err error) {
	interval := cfg.OTLPInterval
	if interval <= 0 {
		interval = otlpDefaultInterval
	}
	exporter, err := otlpexport.New(cfg.OTLPEndpoint, cfg.OTLPHeaders)
	if err != nil {
		return nil, err
	}
	service := exportmetrics.New(cfg.Source, reader, wal, exporter, clock,
		exportmetrics.WithLog(log), exportmetrics.WithTimeout(otlpAttemptTimeout(interval)))

	exportCtx, cancel := context.WithCancel(ctx)
	var done sync.WaitGroup
	done.Add(1)
	go func() {
		defer done.Done()
		runMetricExport(exportCtx, service, interval)
	}()
	log.Info(ctx, "metrics-export: gestartet", "interval_ms", interval.Milliseconds())
	return func() {
		cancel()
		done.Wait()
	}, nil
}
