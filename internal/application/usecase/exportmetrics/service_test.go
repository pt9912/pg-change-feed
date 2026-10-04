package exportmetrics_test

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/exportmetrics"
	"github.com/pt9912/pg-change-feed/internal/domain/messagecode"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

type fakeClock struct {
	mu  sync.Mutex
	now int64
}

func (c *fakeClock) Now() model.TimePoint {
	c.mu.Lock()
	defer c.mu.Unlock()
	return model.NewTimePoint(c.now)
}

func (c *fakeClock) advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now += int64(d)
}

type fakeReader struct {
	samples []outbound.MetricSample
	err     error
	calls   int
}

func (r *fakeReader) Read(context.Context) ([]outbound.MetricSample, error) {
	r.calls++
	if r.err != nil {
		return nil, r.err
	}
	return append([]outbound.MetricSample(nil), r.samples...), nil
}

type fakeWAL struct {
	bytes int64
	at    model.TimePoint
	ok    bool
}

func (w fakeWAL) LastWALRetention() (int64, model.TimePoint, bool) { return w.bytes, w.at, w.ok }

type fakeExporter struct {
	batches []outbound.MetricBatch
	err     error
	// blockUntilDone lässt Export auf das Ende des Kontexts warten.
	blockUntilDone bool
}

func (e *fakeExporter) Export(ctx context.Context, batch outbound.MetricBatch) error {
	if e.blockUntilDone {
		<-ctx.Done()
		return fmt.Errorf("%w: %w", outbound.ErrMetricExport, ctx.Err())
	}
	e.batches = append(e.batches, batch)
	return e.err
}

type logEntry struct {
	level string
	msg   string
	attrs map[string]any
}

type recordingLog struct {
	mu      sync.Mutex
	entries []logEntry
}

func (l *recordingLog) add(level, msg string, attrs []any) {
	m := map[string]any{}
	for i := 0; i+1 < len(attrs); i += 2 {
		m[fmt.Sprint(attrs[i])] = attrs[i+1]
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries = append(l.entries, logEntry{level: level, msg: msg, attrs: m})
}

func (l *recordingLog) Debug(_ context.Context, msg string, attrs ...any) { l.add("debug", msg, attrs) }
func (l *recordingLog) Info(_ context.Context, msg string, attrs ...any)  { l.add("info", msg, attrs) }
func (l *recordingLog) Warn(_ context.Context, msg string, attrs ...any)  { l.add("warn", msg, attrs) }
func (l *recordingLog) Error(_ context.Context, msg string, attrs ...any) { l.add("error", msg, attrs) }

func (l *recordingLog) level(level string) []logEntry {
	l.mu.Lock()
	defer l.mu.Unlock()
	var out []logEntry
	for _, e := range l.entries {
		if e.level == level {
			out = append(out, e)
		}
	}
	return out
}

type fixture struct {
	clock    *fakeClock
	reader   *fakeReader
	wal      *fakeWAL
	exporter *fakeExporter
	log      *recordingLog
	service  *exportmetrics.Service
}

const startNanos = int64(1_700_000_000_000_000_000)

func newFixture(opts ...exportmetrics.Option) *fixture {
	f := &fixture{
		clock:    &fakeClock{now: startNanos},
		reader:   &fakeReader{samples: []outbound.MetricSample{{Name: "cdc_transactions_total", Value: outbound.MetricValue{Float: 3, Int: 3, IsInt: true}}}},
		wal:      &fakeWAL{},
		exporter: &fakeExporter{},
		log:      &recordingLog{},
	}
	opts = append([]exportmetrics.Option{exportmetrics.WithLog(f.log)}, opts...)
	f.service = exportmetrics.New("src-1", f.reader, f.wal, f.exporter, f.clock, opts...)
	return f
}

// TestExportAddsTheLastWALMeasurementWithItsOwnTime belegt den Zyklus: die
// Zeilen der Sicht tragen den Zeitpunkt, zu dem das Lesen begann, der
// WAL-Rückstand den Zeitpunkt seiner Messung; die Quellkennung geht mit.
// Rot färbende Mutation: im Service die Zuweisung `samples[i].At = measuredAt`
// streichen — die Zeilen der Sicht tragen keinen Zeitpunkt; den Zeitpunkt des
// WAL-Datenpunkts durch `measuredAt` ersetzen — die Messzeit der WAL-Messung
// geht verloren.
func TestExportAddsTheLastWALMeasurementWithItsOwnTime(t *testing.T) {
	f := newFixture()
	walAt := model.NewTimePoint(startNanos - int64(30*time.Second))
	*f.wal = fakeWAL{bytes: 4096, at: walAt, ok: true}

	if err := f.service.Export(context.Background()); err != nil {
		t.Fatalf("Export: %v", err)
	}
	if len(f.exporter.batches) != 1 {
		t.Fatalf("Übertragungen = %d, wollen 1", len(f.exporter.batches))
	}
	batch := f.exporter.batches[0]
	if batch.Source != "src-1" || len(batch.Samples) != 2 {
		t.Fatalf("Batch = %+v, wollen Quelle src-1 und 2 Zeilen", batch)
	}
	if got := batch.Samples[0]; got.Name != "cdc_transactions_total" || got.At.UnixNanos != startNanos {
		t.Errorf("Zeile der Sicht = %+v, wollen Zeitpunkt %d", got, startNanos)
	}
	wal := batch.Samples[1]
	if wal.Name != "cdc_wal_retention_bytes" || !wal.Value.IsInt || wal.Value.Int != 4096 || wal.At != walAt {
		t.Errorf("WAL-Zeile = %+v, wollen 4096 Byte mit dem Zeitpunkt der Messung", wal)
	}
}

// TestExportOmitsWALDatapointUntilMeasured belegt `SPEC-033` (Quelle der
// Werte): ist noch kein WAL-Rückstand gemessen, entfällt genau dieser
// Datenpunkt; die Zeilen der Sicht gehen hinaus.
// Rot färbende Mutation: die Bedingung `ok` am Halter streichen — der Zyklus
// überträgt einen Datenpunkt mit Wert 0.
func TestExportOmitsWALDatapointUntilMeasured(t *testing.T) {
	f := newFixture()
	if err := f.service.Export(context.Background()); err != nil {
		t.Fatalf("Export: %v", err)
	}
	samples := f.exporter.batches[0].Samples
	if len(samples) != 1 || samples[0].Name != "cdc_transactions_total" {
		t.Fatalf("Zeilen = %+v, wollen nur die der Sicht", samples)
	}
}

// TestFailureCodesFollowTheStage belegt die zwei Codes: ein Fehlschlag beim
// Lesen trägt `PCF-W6002` und überträgt nichts, ein Fehlschlag der Übertragung
// `PCF-W6001`; beide tragen Stufe und — soweit der Fehler ihn führt — den
// Detail-Text, nie den rohen Fehlertext (`geheim`).
// Rot färbende Mutation: die beiden Codes im Service vertauschen — die Prüfung
// des Codes je Stufe wird rot; `"error", err` in die Warn-Zeile aufnehmen — die
// Prüfung auf den rohen Text wird rot.
func TestFailureCodesFollowTheStage(t *testing.T) {
	readErr := fmt.Errorf("%w: geheim", outbound.ErrMetricsRead)
	f := newFixture()
	f.reader.err = readErr
	if err := f.service.Export(context.Background()); !errors.Is(err, outbound.ErrMetricsRead) {
		t.Fatalf("Fehler = %v, wollen ErrMetricsRead", err)
	}
	if len(f.exporter.batches) != 0 {
		t.Fatalf("nach einem Lesefehler wurde übertragen")
	}
	warns := f.log.level("warn")
	if len(warns) != 1 || warns[0].attrs[messagecode.LogKey] != messagecode.WarnMetricsReadFailed || warns[0].attrs["stage"] != "Lesen" {
		t.Fatalf("Warnungen = %+v, wollen eine mit PCF-W6002 (Lesen)", warns)
	}

	g := newFixture()
	g.exporter.err = &detailError{detail: "Status 500", raw: "geheim"}
	if err := g.service.Export(context.Background()); !errors.Is(err, outbound.ErrMetricExport) {
		t.Fatalf("Fehler = %v, wollen ErrMetricExport", err)
	}
	warns = g.log.level("warn")
	if len(warns) != 1 || warns[0].attrs[messagecode.LogKey] != messagecode.WarnExportFailed || warns[0].attrs["stage"] != "Übertragung" {
		t.Fatalf("Warnungen = %+v, wollen eine mit PCF-W6001 (Übertragung)", warns)
	}
	if warns[0].attrs["detail"] != "Status 500" {
		t.Errorf("detail = %v, wollen Status 500", warns[0].attrs["detail"])
	}
	for _, w := range append(f.log.level("warn"), warns...) {
		for k, v := range w.attrs {
			if fmt.Sprint(v) == "geheim" || k == "error" {
				t.Errorf("Warn-Zeile trägt rohen Fehlertext: %s=%v", k, v)
			}
		}
	}
}

// detailError ist ein Fehler im Sinn von `ErrMetricExport` mit einem
// Detail-Text; sein roher Text (`raw`) darf in keiner Log-Zeile stehen.
type detailError struct{ detail, raw string }

func (e *detailError) Error() string         { return outbound.ErrMetricExport.Error() + ": " + e.raw }
func (e *detailError) Unwrap() error         { return outbound.ErrMetricExport }
func (e *detailError) FailureDetail() string { return e.detail }

// TestWarningIsThrottledToFiveMinutes belegt die Drosselung (Fake-Uhr, keine
// Wartezeit): der erste Fehlschlag warnt, ein Fehlschlag eine Sekunde vor
// Ablauf der fünf Minuten nicht, einer genau nach fünf Minuten wieder; die
// Wiederaufnahme erzeugt genau eine Info-Zeile, und der erste Fehlschlag danach
// warnt sofort.
// Rot färbende Mutation: die Abstandsprüfung in `failed` entfernen — jeder
// Fehlschlag warnt, die Zählung `[1 0 1 …]` wird rot; die Grenze `<` auf `<=`
// ändern — der Fehlschlag genau nach fünf Minuten bleibt stumm.
func TestWarningIsThrottledToFiveMinutes(t *testing.T) {
	f := newFixture()
	f.exporter.err = outbound.ErrMetricExport

	steps := []struct {
		advance  time.Duration
		wantWarn int
	}{
		{0, 1},
		{1 * time.Second, 1},
		{4*time.Minute + 58*time.Second, 1},
		{1 * time.Second, 2},
		{4*time.Minute + 59*time.Second, 2},
		{1 * time.Second, 3},
	}
	for i, step := range steps {
		f.clock.advance(step.advance)
		_ = f.service.Export(context.Background())
		if got := len(f.log.level("warn")); got != step.wantWarn {
			t.Fatalf("Schritt %d: %d Warnungen, wollen %d", i, got, step.wantWarn)
		}
	}
	if got := len(f.log.level("info")); got != 0 {
		t.Fatalf("Info-Zeilen vor der Wiederaufnahme = %d, wollen 0", got)
	}

	f.exporter.err = nil
	if err := f.service.Export(context.Background()); err != nil {
		t.Fatalf("Export nach Wiederaufnahme: %v", err)
	}
	if got := len(f.log.level("info")); got != 1 {
		t.Fatalf("Info-Zeilen der Wiederaufnahme = %d, wollen 1", got)
	}
	if err := f.service.Export(context.Background()); err != nil {
		t.Fatalf("zweiter Export: %v", err)
	}
	if got := len(f.log.level("info")); got != 1 {
		t.Fatalf("ein zweiter Erfolg erzeugte eine Info-Zeile (%d)", got)
	}

	f.exporter.err = outbound.ErrMetricExport
	f.clock.advance(time.Second)
	_ = f.service.Export(context.Background())
	if got := len(f.log.level("warn")); got != 4 {
		t.Fatalf("der erste Fehlschlag nach der Wiederaufnahme warnte nicht sofort: %d Warnungen", got)
	}
}

// TestReadAndTransmitShareOneFailureState belegt den gemeinsamen
// Fehlerzustand: ein Lesefehler und danach ein Übertragungsfehler innerhalb von
// fünf Minuten warnen zusammen einmal; nach fünf Minuten warnt der nächste
// Fehlschlag mit dem Code seiner eigenen Stufe.
func TestReadAndTransmitShareOneFailureState(t *testing.T) {
	f := newFixture()
	f.reader.err = outbound.ErrMetricsRead
	_ = f.service.Export(context.Background())
	f.reader.err = nil
	f.exporter.err = outbound.ErrMetricExport
	f.clock.advance(time.Minute)
	_ = f.service.Export(context.Background())
	if got := len(f.log.level("warn")); got != 1 {
		t.Fatalf("Warnungen = %d, wollen 1 (gemeinsamer Zustand)", got)
	}
	f.clock.advance(4 * time.Minute)
	_ = f.service.Export(context.Background())
	warns := f.log.level("warn")
	if len(warns) != 2 || warns[1].attrs[messagecode.LogKey] != messagecode.WarnExportFailed {
		t.Fatalf("Warnungen = %+v, wollen eine zweite mit PCF-W6001", warns)
	}
}

// TestExportEndsAtTheDeadline belegt die Frist je Versuch: ein Empfänger, der
// nicht antwortet, hält den Zyklus nur bis `WithTimeout`; der Zyklus endet mit
// dem Fehlschlag und der Warnung. Der Test hängt nicht, wenn die Frist nicht
// greift: er bricht nach zwei Sekunden mit einem Befund ab.
// Rot färbende Mutation: im Service die Ableitung mit `context.WithTimeout`
// streichen — der Zyklus kehrt nicht zurück, der Test meldet „Frist greift nicht“.
func TestExportEndsAtTheDeadline(t *testing.T) {
	f := newFixture(exportmetrics.WithTimeout(50 * time.Millisecond))
	f.exporter.blockUntilDone = true

	done := make(chan error, 1)
	go func() { done <- f.service.Export(context.Background()) }()
	select {
	case err := <-done:
		if !errors.Is(err, outbound.ErrMetricExport) {
			t.Fatalf("Fehler = %v, wollen ErrMetricExport", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Frist greift nicht: der Zyklus kehrt nach zwei Sekunden nicht zurück")
	}
	if got := len(f.log.level("warn")); got != 1 {
		t.Fatalf("Warnungen = %d, wollen 1", got)
	}
}

// TestShutdownIsNoFailure belegt, dass das Ende des übergebenen Kontexts
// (Prozessende) keine Warnung erzeugt.
// Rot färbende Mutation: die Bedingung `parent.Err() == nil` vor dem Aufruf von
// `failed` streichen — der Zyklus warnt beim Prozessende.
func TestShutdownIsNoFailure(t *testing.T) {
	f := newFixture()
	f.exporter.blockUntilDone = true
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := f.service.Export(ctx); err == nil {
		t.Fatal("Export mit beendetem Kontext meldet keinen Fehler")
	}
	if got := len(f.log.level("warn")); got != 0 {
		t.Fatalf("Warnungen beim Prozessende = %d, wollen 0", got)
	}
}
