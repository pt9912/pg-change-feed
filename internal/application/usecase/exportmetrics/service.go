// Package exportmetrics trägt den Use Case des periodischen Metrik-Exports
// (`ADR-0149`): ein Zyklus liest die Kennzahlen der Sicht über den
// `MetricsReadPort`, ergänzt den zuletzt gemessenen WAL-Rückstand und
// überträgt sie über den `MetricExportPort`. Der Takt, die eigene Goroutine
// und die Verdrahtung gehören dem Bootstrap; der Use Case trägt die Frist je
// Versuch und die gedrosselte Warnung.
package exportmetrics

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/domain/messagecode"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// warnInterval ist der Mindestabstand zweier Warnungen während eines
// andauernden Fehlzustands.
const warnInterval = 5 * time.Minute

// walRetentionMetric ist der Name des Datenpunkts, den der Prozess selbst
// misst; die Sicht führt ihn nicht.
const walRetentionMetric = "cdc_wal_retention_bytes"

// Service implementiert `inbound.ExportMetricsUseCase`.
type Service struct {
	source   model.SourceID
	reader   outbound.MetricsReadPort
	wal      outbound.WALRetentionSource
	exporter outbound.MetricExportPort
	clock    outbound.ClockPort
	log      outbound.LogPort
	timeout  time.Duration

	mu       sync.Mutex
	failing  bool
	lastWarn model.TimePoint
}

var _ inbound.ExportMetricsUseCase = (*Service)(nil)

// Option konfiguriert den Service.
type Option func(*Service)

// WithLog setzt den `LogPort`; ohne ihn bleibt die Protokollierung ein No-Op.
func WithLog(log outbound.LogPort) Option { return func(s *Service) { s.log = log } }

// WithTimeout setzt die Frist je Versuch (Lesen und Übertragen zusammen); ohne
// sie (oder bei einem Wert bis null) gilt die Frist des übergebenen Kontexts.
func WithTimeout(d time.Duration) Option { return func(s *Service) { s.timeout = d } }

// New verdrahtet den Use Case.
func New(source model.SourceID, reader outbound.MetricsReadPort, wal outbound.WALRetentionSource, exporter outbound.MetricExportPort, clock outbound.ClockPort, opts ...Option) *Service {
	s := &Service{source: source, reader: reader, wal: wal, exporter: exporter, clock: clock, log: outbound.NoopLog}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Export führt einen Zyklus aus. Jede Zeile der Sicht trägt den Zeitpunkt, zu
// dem das Lesen begann; der WAL-Rückstand trägt den Zeitpunkt seiner Messung
// und entfällt, solange keiner gemessen ist. Ein Fehlschlag beim Lesen oder
// Übertragen geht zurück (`outbound.ErrMetricsRead`, `outbound.ErrMetricExport`)
// und ist zugleich protokolliert, außer wenn der übergebene Kontext selbst
// endete (Prozessende ist kein Fehlzustand des Empfängers); ein ausgefallener
// Zyklus wird nicht nachgeholt.
func (s *Service) Export(parent context.Context) error {
	ctx := parent
	if s.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(parent, s.timeout)
		defer cancel()
	}

	measuredAt := s.clock.Now()
	samples, err := s.reader.Read(ctx)
	if err != nil {
		if parent.Err() == nil {
			s.failed(ctx, messagecode.WarnMetricsReadFailed, "Lesen", err)
		}
		return err
	}
	for i := range samples {
		samples[i].At = measuredAt
	}
	if bytes, at, ok := s.wal.LastWALRetention(); ok {
		samples = append(samples, outbound.MetricSample{
			Name:  walRetentionMetric,
			Value: outbound.MetricValue{Float: float64(bytes), Int: bytes, IsInt: true},
			At:    at,
		})
	}
	if err := s.exporter.Export(ctx, outbound.MetricBatch{Source: s.source, Samples: samples}); err != nil {
		if parent.Err() == nil {
			s.failed(ctx, messagecode.WarnExportFailed, "Übertragung", err)
		}
		return err
	}
	s.succeeded(ctx, len(samples))
	return nil
}

// failed vermerkt einen Fehlschlag: der erste seit dem Start oder seit der
// letzten Wiederaufnahme warnt, weitere höchstens alle `warnInterval`, jeweils
// mit dem Code der Stufe des Fehlschlags. Die Zeile trägt Code, Stufe und —
// wenn der Fehler ihn führt — den Detail-Text des Adapters, nie den rohen
// Fehlertext einer Bibliothek.
func (s *Service) failed(ctx context.Context, code messagecode.Code, stage string, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.clock.Now()
	if s.failing && now.UnixNanos-s.lastWarn.UnixNanos < int64(warnInterval) {
		return
	}
	s.failing = true
	s.lastWarn = now
	attrs := []any{messagecode.LogKey, code, "stage", stage}
	var detail outbound.FailureDetail
	if errors.As(err, &detail) {
		attrs = append(attrs, "detail", detail.FailureDetail())
	}
	s.log.Warn(ctx, "metrics-export: Zyklus fehlgeschlagen", attrs...)
}

// succeeded vermerkt einen vollständig erfolgreichen Zyklus; auf einen
// Fehlzustand folgt die Info-Zeile der Wiederaufnahme.
func (s *Service) succeeded(ctx context.Context, points int) {
	s.mu.Lock()
	wasFailing := s.failing
	s.failing = false
	s.mu.Unlock()
	if wasFailing {
		s.log.Info(ctx, "metrics-export: Übertragung wieder aufgenommen", "points", points)
		return
	}
	s.log.Debug(ctx, "metrics-export: übertragen", "points", points)
}
