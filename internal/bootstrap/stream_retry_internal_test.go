package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/mapper"
	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/receive"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// retryClock ist die deterministische Uhr der Retry-Tests: Now() liest den
// gesetzten Stand, die Aufzeichnung der Wartezüge zieht ihn weiter.
type retryClock struct{ now model.TimePoint }

func (c *retryClock) Now() model.TimePoint { return c.now }

func (c *retryClock) advance(d time.Duration) {
	c.now = model.NewTimePoint(c.now.UnixNanos + int64(d))
}

// retrySleeper zeichnet die Warteschritte auf und zieht die Uhr des Tests
// um den gewarteten Schritt weiter.
type retrySleeper struct {
	clock *retryClock
	waits []time.Duration
}

func (s *retrySleeper) sleep(d time.Duration) {
	s.clock.advance(d)
	s.waits = append(s.waits, d)
}

// scriptCycle führt die Skript-Fehler der Reihe nach auf: ein leerer
// Skript-Eintrag bedeutet Erfolg. setups[i] ist die Dauer des Aufbaus des
// i-ten Zyklus vor dem Streaming-Signal, runtimes[i] seine Streaming-Dauer
// bis zum Fehler, beide auf der Test-Uhr. Ein Zyklus in noStream erreicht
// den Streaming-Zustand nie und meldet kein Signal; seine Laufzeit steht
// dann in setups[i].
type scriptCycle struct {
	clock    *retryClock
	errs     []error
	setups   []time.Duration
	runtimes []time.Duration
	noStream map[int]bool
	calls    int
}

func (c *scriptCycle) run(_ context.Context, streaming func()) error {
	i := c.calls
	c.calls++
	if i < len(c.setups) {
		c.clock.advance(c.setups[i])
	}
	if !c.noStream[i] {
		streaming()
	}
	if i < len(c.runtimes) {
		c.clock.advance(c.runtimes[i])
	}
	if i < len(c.errs) {
		return c.errs[i]
	}
	return nil
}

// attrLog zeichnet Info- und Warn-Einträge samt Attributen auf.
type attrLog struct {
	outbound.LogPort
	infos []logEntry
	warns []logEntry
}

type logEntry struct {
	msg   string
	attrs []any
}

func (l *attrLog) Info(_ context.Context, msg string, attrs ...any) {
	l.infos = append(l.infos, logEntry{msg, attrs})
}

func (l *attrLog) Warn(_ context.Context, msg string, attrs ...any) {
	l.warns = append(l.warns, logEntry{msg, attrs})
}

func (e logEntry) attr(key string) any {
	for i := 0; i+1 < len(e.attrs); i += 2 {
		if e.attrs[i] == key {
			return e.attrs[i+1]
		}
	}
	return nil
}

func TestRunStreamWithRetryBackoffFolge(t *testing.T) {
	clk := &retryClock{now: model.NewTimePoint(0)}
	sleeper := &retrySleeper{clock: clk}
	cycle := &scriptCycle{clock: clk, errs: []error{
		receive.ErrReplication,
		outbound.ErrReplication,
		receive.ErrReplication,
		receive.ErrReplication,
	}}
	err := runStreamWithRetry(context.Background(), outbound.NoopLog, clk, cycle.run, sleeper.sleep)
	if err != nil {
		t.Fatalf("Erfolg nach Wiederholungen: %v", err)
	}
	// Werte der ADR-0135 Festlegung 2: Anfangsverzögerung 2 s, Faktor 2.
	want := []time.Duration{2 * time.Second, 4 * time.Second, 8 * time.Second, 16 * time.Second}
	if len(sleeper.waits) != len(want) {
		t.Fatalf("Warteschritte = %v, erwartet %v", sleeper.waits, want)
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
	cycle := &scriptCycle{clock: clk, errs: []error{
		receive.ErrReplication, receive.ErrReplication, receive.ErrReplication,
		receive.ErrReplication, receive.ErrReplication, receive.ErrReplication,
	}}
	err := runStreamWithRetry(context.Background(), outbound.NoopLog, clk, cycle.run, sleeper.sleep)
	if err != nil {
		t.Fatalf("Erfolg nach Wiederholungen: %v", err)
	}
	// Obergrenze 30 s je Warteschritt (ADR-0135 Festlegung 2).
	want := []time.Duration{
		2 * time.Second, 4 * time.Second, 8 * time.Second, 16 * time.Second,
		30 * time.Second, 30 * time.Second,
	}
	if len(sleeper.waits) != len(want) {
		t.Fatalf("Warteschritte = %v, erwartet %v", sleeper.waits, want)
	}
	for i, w := range want {
		if sleeper.waits[i] != w {
			t.Fatalf("Warteschritt %d = %v, erwartet %v", i, sleeper.waits[i], w)
		}
	}
}

// TestRunStreamWithRetryGesamtfensterFuenfMinuten bindet das Gesamtfenster an
// 5 Minuten: der Fehler bei genau 5 Minuten nach dem ersten wiederholten
// Fehler wird noch wiederholt, der bei 6 Minuten 40 Sekunden erschöpft.
func TestRunStreamWithRetryGesamtfensterFuenfMinuten(t *testing.T) {
	clk := &retryClock{now: model.NewTimePoint(0)}
	cycle := &scriptCycle{clock: clk, errs: []error{
		receive.ErrReplication, receive.ErrReplication, receive.ErrReplication,
		receive.ErrReplication, receive.ErrReplication, receive.ErrReplication,
		receive.ErrReplication, receive.ErrReplication,
	}}
	// Jeder Warteschritt zieht die Uhr um 100 s: Fehler bei t = 0, 100, 200,
	// 300 (genau 5 min, innerhalb), 400 (über dem Fenster).
	err := runStreamWithRetry(context.Background(), outbound.NoopLog, clk, cycle.run, func(time.Duration) {
		clk.advance(100 * time.Second)
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
	if cycle.calls != 5 {
		t.Fatalf("Zyklen bis zur Erschöpfung = %d, erwartet 5 (Fehler bei 0/100/200/300/400 s)", cycle.calls)
	}
}

// TestRunStreamWithRetryEpisodeNachLangemZyklusZurueckgesetzt belegt die
// Rücksetzung in einem Aufruf: Zyklus 1 scheitert sofort, Zyklus 2 läuft eine
// Stunde und scheitert, Zyklus 3 scheitert sofort, Zyklus 4 gelingt. Die
// Wartezüge zeigen die Anfangsverzögerung nach dem langen Zyklus, das
// Fenster ist zurückgesetzt (kein ErrTransientExhausted nach einer Stunde).
func TestRunStreamWithRetryEpisodeNachLangemZyklusZurueckgesetzt(t *testing.T) {
	clk := &retryClock{now: model.NewTimePoint(0)}
	sleeper := &retrySleeper{clock: clk}
	log := &attrLog{LogPort: outbound.NoopLog}
	cycle := &scriptCycle{
		clock: clk,
		errs: []error{
			receive.ErrReplication, receive.ErrReplication, receive.ErrReplication, nil,
		},
		runtimes: []time.Duration{0, time.Hour, 0, 0},
	}
	err := runStreamWithRetry(context.Background(), log, clk, cycle.run, sleeper.sleep)
	if err != nil {
		t.Fatalf("Wiederholung statt Erschöpfung erwartet: %v", err)
	}
	want := []time.Duration{2 * time.Second, 2 * time.Second, 4 * time.Second}
	if len(sleeper.waits) != len(want) {
		t.Fatalf("Warteschritte = %v, erwartet %v", sleeper.waits, want)
	}
	for i, w := range want {
		if sleeper.waits[i] != w {
			t.Fatalf("Warteschritt %d = %v, erwartet %v", i, sleeper.waits[i], w)
		}
	}
}

// TestRunStreamWithRetryStabilitaetsschwelle bindet die Schwelle des
// erfolgreichen Zyklus an 30 s: 29 s Laufzeit setzen die Episode nicht
// zurück, 30 s setzen sie zurück.
func TestRunStreamWithRetryStabilitaetsschwelle(t *testing.T) {
	cases := []struct {
		name    string
		runtime time.Duration
		want    []time.Duration
	}{
		{"29 s setzt nicht zurück", 29 * time.Second, []time.Duration{2 * time.Second, 4 * time.Second}},
		{"30 s setzt zurück", 30 * time.Second, []time.Duration{2 * time.Second, 2 * time.Second}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			clk := &retryClock{now: model.NewTimePoint(0)}
			sleeper := &retrySleeper{clock: clk}
			cycle := &scriptCycle{
				clock:    clk,
				errs:     []error{receive.ErrReplication, receive.ErrReplication, nil},
				runtimes: []time.Duration{0, c.runtime},
			}
			if err := runStreamWithRetry(context.Background(), outbound.NoopLog, clk, cycle.run, sleeper.sleep); err != nil {
				t.Fatalf("Erfolg erwartet: %v", err)
			}
			if len(sleeper.waits) != len(c.want) {
				t.Fatalf("Warteschritte = %v, erwartet %v", sleeper.waits, c.want)
			}
			for i, w := range c.want {
				if sleeper.waits[i] != w {
					t.Fatalf("Warteschritt %d = %v, erwartet %v", i, sleeper.waits[i], w)
				}
			}
		})
	}
}

// TestRunStreamWithRetrySichtbarkeit belegt ADR-0135 Festlegung 4: ein WARN
// je Wiederholung mit hochzählendem Versuchszähler und ein INFO bei der
// Fortsetzung nach einer Wiederholung.
func TestRunStreamWithRetrySichtbarkeit(t *testing.T) {
	clk := &retryClock{now: model.NewTimePoint(0)}
	sleeper := &retrySleeper{clock: clk}
	log := &attrLog{LogPort: outbound.NoopLog}
	cycle := &scriptCycle{clock: clk, errs: []error{receive.ErrReplication, receive.ErrReplication, nil}}
	if err := runStreamWithRetry(context.Background(), log, clk, cycle.run, sleeper.sleep); err != nil {
		t.Fatalf("Erfolg erwartet: %v", err)
	}
	if len(log.warns) != 2 {
		t.Fatalf("WARN-Einträge = %d, erwartet 2", len(log.warns))
	}
	for i, w := range log.warns {
		if got := w.attr("versuch"); got != i+1 {
			t.Errorf("WARN %d: versuch = %v, erwartet %d", i, got, i+1)
		}
		if got := w.attr("warteschritt"); got != sleeper.waits[i] {
			t.Errorf("WARN %d: warteschritt = %v, erwartet %v", i, got, sleeper.waits[i])
		}
		if w.attr("error") == nil {
			t.Errorf("WARN %d: Fehlertext fehlt", i)
		}
	}
	// Zwei Zyklen erreichen das Streaming nach einer Wiederholung
	// (versuche 1 und 2), der dritte endet regulär (versuche 2).
	if len(log.infos) != 3 {
		t.Fatalf("INFO-Einträge = %d, erwartet 3: %v", len(log.infos), log.infos)
	}
	for i, want := range []any{1, 2, 2} {
		if got := log.infos[i].attr("versuche"); got != want {
			t.Errorf("INFO %d: versuche = %v, erwartet %v", i, got, want)
		}
	}
	if log.infos[2].msg != "pg-change-feed: Stream-Zyklus nach Wiederholung regulär beendet" {
		t.Errorf("INFO 2 = %q, erwartet reguläres Ende", log.infos[2].msg)
	}
}

// TestRunStreamWithRetryAufbauGrenzeGesamtfenster belegt ADR-0136 Festlegung
// 1: Zyklen, die 31 s im Aufbau hängen und den Streaming-Zustand nie
// erreichen, setzen die Episode nie zurück — das Gesamtfenster erschöpft.
func TestRunStreamWithRetryAufbauGrenzeGesamtfenster(t *testing.T) {
	clk := &retryClock{now: model.NewTimePoint(0)}
	sleeper := &retrySleeper{clock: clk}
	const n = 60
	cycle := &scriptCycle{clock: clk, noStream: map[int]bool{}}
	for i := 0; i < n; i++ {
		cycle.errs = append(cycle.errs, receive.ErrReplication)
		cycle.setups = append(cycle.setups, 31*time.Second)
		cycle.noStream[i] = true
	}
	err := runStreamWithRetry(context.Background(), outbound.NoopLog, clk, cycle.run, sleeper.sleep)
	if !errors.Is(err, ErrTransientExhausted) {
		t.Fatalf("Fehler = %v, erwartet ErrTransientExhausted", err)
	}
	if got := classifyRunError(err); got != model.ErrorClassTransient {
		t.Fatalf("Klasse = %q, erwartet transient", got)
	}
	if cycle.calls >= n {
		t.Fatalf("Zyklen = %d: das Gesamtfenster erschöpft nicht", cycle.calls)
	}
}

// TestRunStreamWithRetryAufbauZaehltNichtZurStabilitaet belegt, dass die
// Dauer des Aufbaus vor dem Signal nicht zur Stabilität zählt: 31 s Aufbau
// und 1 s Streaming setzen die Episode nicht zurück, 30 s Streaming nach
// demselben Aufbau setzen sie zurück.
func TestRunStreamWithRetryAufbauZaehltNichtZurStabilitaet(t *testing.T) {
	cases := []struct {
		name   string
		stream time.Duration
		want   []time.Duration
	}{
		{"1 s Streaming", time.Second, []time.Duration{2 * time.Second, 4 * time.Second}},
		{"30 s Streaming", 30 * time.Second, []time.Duration{2 * time.Second, 2 * time.Second}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			clk := &retryClock{now: model.NewTimePoint(0)}
			sleeper := &retrySleeper{clock: clk}
			cycle := &scriptCycle{
				clock:    clk,
				errs:     []error{receive.ErrReplication, receive.ErrReplication, nil},
				setups:   []time.Duration{0, 31 * time.Second},
				runtimes: []time.Duration{0, c.stream},
			}
			if err := runStreamWithRetry(context.Background(), outbound.NoopLog, clk, cycle.run, sleeper.sleep); err != nil {
				t.Fatalf("Erfolg erwartet: %v", err)
			}
			if len(sleeper.waits) != len(c.want) {
				t.Fatalf("Warteschritte = %v, erwartet %v", sleeper.waits, c.want)
			}
			for i, w := range c.want {
				if sleeper.waits[i] != w {
					t.Fatalf("Warteschritt %d = %v, erwartet %v", i, sleeper.waits[i], w)
				}
			}
		})
	}
}

// TestRunStreamWithRetryOhneSignalKeinFortgesetzt belegt ADR-0136 Festlegung
// 4: das INFO der Fortsetzung folgt dem Streaming-Signal; Zyklen ohne Signal
// schreiben es nicht.
func TestRunStreamWithRetryOhneSignalKeinFortgesetzt(t *testing.T) {
	clk := &retryClock{now: model.NewTimePoint(0)}
	log := &attrLog{LogPort: outbound.NoopLog}
	cycle := &scriptCycle{
		clock:    clk,
		errs:     []error{receive.ErrReplication, receive.ErrReplication, nil},
		noStream: map[int]bool{1: true, 2: true},
	}
	if err := runStreamWithRetry(context.Background(), log, clk, cycle.run, func(time.Duration) {}); err != nil {
		t.Fatalf("Erfolg erwartet: %v", err)
	}
	for _, e := range log.infos {
		if e.msg == "pg-change-feed: Stream-Zyklus nach Wiederholung fortgesetzt" {
			t.Fatalf("INFO der Fortsetzung ohne Streaming-Signal: %v", log.infos)
		}
	}
}

// TestRunStreamWithRetryFortsetzungImLog belegt das INFO der Fortsetzung: ein
// Zyklus nach einer Wiederholung, der bis zu seinem Fehler gestreamt hat,
// schreibt "fortgesetzt" mit der Zahl der Versuche der beendeten Episode.
func TestRunStreamWithRetryFortsetzungImLog(t *testing.T) {
	clk := &retryClock{now: model.NewTimePoint(0)}
	log := &attrLog{LogPort: outbound.NoopLog}
	cycle := &scriptCycle{
		clock:    clk,
		errs:     []error{receive.ErrReplication, receive.ErrReplication, nil},
		runtimes: []time.Duration{0, time.Hour},
	}
	if err := runStreamWithRetry(context.Background(), log, clk, cycle.run, func(time.Duration) {}); err != nil {
		t.Fatalf("Erfolg erwartet: %v", err)
	}
	if len(log.infos) == 0 || log.infos[0].msg != "pg-change-feed: Stream-Zyklus nach Wiederholung fortgesetzt" {
		t.Fatalf("INFO der Fortsetzung fehlt: %v", log.infos)
	}
	if got := log.infos[0].attr("versuche"); got != 1 {
		t.Errorf("INFO: versuche = %v, erwartet 1", got)
	}
}

// TestRunStreamWithRetryOhneFehlerKeinLog belegt, dass ein Lauf ohne
// Wiederholung weder WARN noch INFO schreibt.
func TestRunStreamWithRetryOhneFehlerKeinLog(t *testing.T) {
	clk := &retryClock{now: model.NewTimePoint(0)}
	log := &attrLog{LogPort: outbound.NoopLog}
	cycle := &scriptCycle{clock: clk}
	if err := runStreamWithRetry(context.Background(), log, clk, cycle.run, func(time.Duration) {}); err != nil {
		t.Fatalf("Erfolg erwartet: %v", err)
	}
	if len(log.warns)+len(log.infos) != 0 {
		t.Fatalf("Log-Einträge ohne Wiederholung: %v %v", log.warns, log.infos)
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
		{"permission Sentinel", receive.ErrPermission},
		{"permission 42501 in der Kette", fmt.Errorf("%w: %v", receive.ErrPermission, &pgconn.PgError{Code: "42501"})},
		{"Server-Abweisung", fmt.Errorf("%w: %w", receive.ErrReplication, receive.ErrRejected)},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			clk := &retryClock{now: model.NewTimePoint(0)}
			sleeper := &retrySleeper{clock: clk}
			cycle := &scriptCycle{clock: clk, errs: []error{c.err}}
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

func TestClassifyRunErrorPermission(t *testing.T) {
	err := fmt.Errorf("%w: Verbindungsaufbau: %v", receive.ErrPermission, &pgconn.PgError{Code: "28P01"})
	if got := classifyRunError(err); got != model.ErrorClassPermission {
		t.Fatalf("Klasse = %q, erwartet permission", got)
	}
}

func TestRunStreamWithRetryKontextEndeWährendWartensEndetRegulär(t *testing.T) {
	clk := &retryClock{now: model.NewTimePoint(0)}
	sleeper := &retrySleeper{clock: clk}
	ctx, cancel := context.WithCancel(context.Background())
	cycle := &scriptCycle{clock: clk, errs: []error{receive.ErrReplication}}
	err := runStreamWithRetry(ctx, outbound.NoopLog, clk, cycle.run, func(time.Duration) { cancel() })
	if err != nil {
		t.Fatalf("Kontext-Ende während des Wartens endet regulär: %v", err)
	}
	if len(sleeper.waits) != 0 {
		t.Fatalf("Warteschritte nach Kontext-Ende = %d, erwartet 0", len(sleeper.waits))
	}
}
