package bootstrap

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/mapper"
	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/receive"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// retryClock ist die deterministische Uhr der Retry-Tests: Now() liest den
// gesetzten Stand, die Aufzeichnung der Wartezüge zieht ihn weiter.
type retryClock struct{ now model.TimePoint }

func (c *retryClock) Now() model.TimePoint { return c.now }

// retrySleeper zeichnet die Warteschritte auf und zieht die Uhr des Tests
// um den gewarteten Schritt weiter.
type retrySleeper struct {
	clock *retryClock
	waits []time.Duration
}

func (s *retrySleeper) sleep(d time.Duration) {
	s.clock.now = model.NewTimePoint(s.clock.now.UnixNanos + int64(d))
	s.waits = append(s.waits, d)
}

// scriptCycle führt die Skript-Fehler der Reihe nach auf: ein leerer
// Skript-Eintrag bedeutet Erfolg.
type scriptCycle struct {
	errs   []error
	calls  int
	called func()
}

func (c *scriptCycle) run(context.Context) error {
	if c.called != nil {
		c.called()
	}
	if c.calls < len(c.errs) {
		e := c.errs[c.calls]
		c.calls++
		return e
	}
	c.calls++
	return nil
}

func retryTestClock() *retryClock {
	return &retryClock{now: model.NewTimePoint(1_000_000_000)}
}

func TestRunStreamWithRetryBackoffFolge(t *testing.T) {
	clk := &retryClock{now: model.NewTimePoint(0)}
	sleeper := &retrySleeper{clock: clk}
	cycle := &scriptCycle{errs: []error{
		receive.ErrReplication,
		outbound.ErrReplication,
		receive.ErrReplication,
		receive.ErrReplication,
	}}
	err := runStreamWithRetry(context.Background(), outbound.NoopLog, clk, cycle.run, sleeper.sleep)
	if err != nil {
		t.Fatalf("Erfolg nach Wiederholungen: %v", err)
	}
	want := []time.Duration{
		streamRetryInitialDelay,
		2 * streamRetryInitialDelay,
		4 * streamRetryInitialDelay,
		8 * streamRetryInitialDelay,
	}
	if len(sleeper.waits) != len(want) {
		t.Fatalf("Warteschritte = %d, erwartet %d", len(sleeper.waits), len(want))
	}
	for i, w := range want {
		if sleeper.waits[i] != w {
			t.Fatalf("Warteschritt %d = %v, erwartet %v", i, sleeper.waits[i], w)
		}
	}
}

func TestRunStreamWithRetryObergrenze(t *testing.T) {
	clk := &retryClock{now: model.NewTimePoint(0)}
	sleeper := &retrySleeper{clock: clk}
	cycle := &scriptCycle{errs: []error{
		receive.ErrReplication,
		receive.ErrReplication,
		receive.ErrReplication,
		receive.ErrReplication,
		receive.ErrReplication,
		receive.ErrReplication,
		receive.ErrReplication,
	}}
	err := runStreamWithRetry(context.Background(), outbound.NoopLog, clk, cycle.run, sleeper.sleep)
	if err != nil {
		t.Fatalf("Erfolg nach Wiederholungen: %v", err)
	}
	for i, w := range sleeper.waits {
		if w > streamRetryMaxDelay {
			t.Fatalf("Warteschritt %d = %v über der Obergrenze %v", i, w, streamRetryMaxDelay)
		}
	}
	if sleeper.waits[4] != streamRetryMaxDelay || sleeper.waits[5] != streamRetryMaxDelay {
		t.Fatalf("Obergrenze nicht getragen: %v", sleeper.waits)
	}
}

func TestRunStreamWithRetryErschoepfungEndetTransient(t *testing.T) {
	clk := &retryClock{now: model.NewTimePoint(0)}
	cycle := &scriptCycle{errs: []error{
		receive.ErrReplication,
		receive.ErrReplication,
		receive.ErrReplication,
		receive.ErrReplication,
	}}
	err := runStreamWithRetry(context.Background(), outbound.NoopLog, clk, cycle.run, func(d time.Duration) {
		// Jeder Wartezug zieht die Uhr über das Gesamtfenster hinaus —
		// die Erschöpfung trägt den Lauf-Fehler.
		clk.now = model.NewTimePoint(clk.now.UnixNanos + int64(6*time.Minute))
	})
	if err == nil {
		t.Fatal("Erschöpfung endet mit Fehler")
	}
	if !errors.Is(err, ErrTransientExhausted) {
		t.Fatalf("Erschöpfung trägt ErrTransientExhausted nicht: %v", err)
	}
	if got := classifyRunError(err); got != model.ErrorClassTransient {
		t.Fatalf("Klasse = %q, erwartet transient", got)
	}
}

func TestRunStreamWithRetryEpisodeZurueckgesetzt(t *testing.T) {
	clk := &retryClock{now: model.NewTimePoint(0)}
	sleeper := &retrySleeper{clock: clk}
	// Erste Episode: zwei Fehlversuche, dann Erfolg — der erfolgreiche
	// Zyklus setzt die Episode zurück.
	cycle := &scriptCycle{errs: []error{
		receive.ErrReplication,
		receive.ErrReplication,
		nil, // erfolgreicher Zyklus setzt die Episode zurück
	}}
	err := runStreamWithRetry(context.Background(), outbound.NoopLog, clk, cycle.run, sleeper.sleep)
	if err != nil {
		t.Fatalf("Erfolg nach Episoden: %v", err)
	}
	// Zweite Episode: der erste Warteschritt trägt wieder die
	// Anfangsverzögerung, nicht die Verdopplung der ersten Episode.
	cycle2 := &scriptCycle{errs: []error{
		receive.ErrReplication,
		nil,
	}}
	err = runStreamWithRetry(context.Background(), outbound.NoopLog, clk, cycle2.run, sleeper.sleep)
	if err != nil {
		t.Fatalf("Erfolg in der zweiten Episode: %v", err)
	}
	want := []time.Duration{
		streamRetryInitialDelay,
		2 * streamRetryInitialDelay,
		// Erfolgreicher Zyklus — die Episode beginnt neu:
		streamRetryInitialDelay,
	}
	if len(sleeper.waits) != len(want) {
		t.Fatalf("Warteschritte = %d, erwartet %d", len(sleeper.waits), len(want))
	}
	for i, w := range want {
		if sleeper.waits[i] != w {
			t.Fatalf("Warteschritt %d = %v, erwartet %v", i, sleeper.waits[i], w)
		}
	}
}

func TestRunStreamWithRetryKlassenEndenOhneWiederholung(t *testing.T) {
	cases := []struct {
		name string
		err  error
	}{
		{"configuration", ErrConfiguration},
		{"schema", mapper.ErrIncompatibleSchemaChange},
		{"Stream-Ordnungs-Verletzung", mapper.ErrChangeWithoutBegin},
		{"storage", outbound.ErrStorage},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			clk := &retryClock{now: model.NewTimePoint(0)}
			sleeper := &retrySleeper{clock: clk}
			cycle := &scriptCycle{errs: []error{c.err}}
			err := runStreamWithRetry(context.Background(), outbound.NoopLog, clk, cycle.run, sleeper.sleep)
			if !errors.Is(err, c.err) {
				t.Fatalf("Fehler %v trägt die Klasse nicht: %v", c.err, err)
			}
			if len(sleeper.waits) != 0 {
				t.Fatalf("Nicht wiederholbare Klasse wartet %d Mal", len(sleeper.waits))
			}
		})
	}
}

func TestRunStreamWithRetryKontextEndeWährendWartensEndetRegulär(t *testing.T) {
	clk := &retryClock{now: model.NewTimePoint(0)}
	sleeper := &retrySleeper{clock: clk}
	ctx, cancel := context.WithCancel(context.Background())
	cycle := &scriptCycle{errs: []error{receive.ErrReplication}}
	err := runStreamWithRetry(ctx, outbound.NoopLog, clk, cycle.run, func(time.Duration) { cancel() })
	if err != nil {
		t.Fatalf("Kontext-Ende während des Wartens endet regulär: %v", err)
	}
	if len(sleeper.waits) != 0 {
		t.Fatalf("Warteschritte nach Kontext-Ende = %d, erwartet 0", len(sleeper.waits))
	}
}
