package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/mapper"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
)

// Whitebox-Test (`package bootstrap`, nicht `bootstrap_test`): der
// Schwellen-Vergleich, der Abbruch-Zug und die Rückgabewert-Priorität sind
// unexportierte Verdrahtungsdetails (`slice-026`, `ADR-0049`) — der reale
// Ende-zu-Ende-Beleg über beide Seiten der Schwelle liegt in
// `walretention_endtoend_test.go` (`make test-replication`); dieser Test
// belegt die Vergleichs-/Prioritäts-Logik selbst, ohne reale PostgreSQL-
// Instanz.

// TestResolveWALRetentionThresholdsDefaultsToSpec013 belegt den
// Zero-Value-Fallback von `Config.WALRetentionWarnBytes`/
// `WALRetentionErrorBytes` (Review-Finding F-1, `review-slice-026.md`): ein
// unbesetzter (0 oder negativer) Override übernimmt die
// SPEC-013-Startwerte (100 MiB/1 GiB), nicht eine Schwelle von 0 — ein
// über `ConfigFromEnv` gestarteter Produktionsprozess erreicht `Run` immer
// mit `Config{}`-Zero-Values für beide Felder (kein Umgebungsname mappt
// darauf, siehe `wiring.go` Kommentar an `Config.WALRetentionWarnBytes`).
// Rot färbende Mutation: die Bedingung in `resolveWALRetentionThresholds`
// von `<= 0` auf `>= 0` umkehren — dann liefert dieser Test 0/0 statt der
// SPEC-013-Startwerte.
func TestResolveWALRetentionThresholdsDefaultsToSpec013(t *testing.T) {
	cases := []struct {
		name          string
		warnOverride  int64
		errorOverride int64
	}{
		{"beide Felder 0 (Zero-Value-Config)", 0, 0},
		{"beide Felder negativ", -1, -1},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			warnBytes, errorBytes := resolveWALRetentionThresholds(c.warnOverride, c.errorOverride)
			if warnBytes != walRetentionWarnBytes {
				t.Fatalf("resolveWALRetentionThresholds(%d, %d) warnBytes = %d, wollen SPEC-013-Startwert %d", c.warnOverride, c.errorOverride, warnBytes, walRetentionWarnBytes)
			}
			if errorBytes != walRetentionErrorBytes {
				t.Fatalf("resolveWALRetentionThresholds(%d, %d) errorBytes = %d, wollen SPEC-013-Startwert %d", c.warnOverride, c.errorOverride, errorBytes, walRetentionErrorBytes)
			}
		})
	}
}

// TestResolveWALRetentionThresholdsKeepsPositiveOverride belegt die
// Gegenseite: ein gesetzter positiver Override (Testfixture-Schwellen, wie
// `TestWALRetentionThresholdEndToEnd` sie nutzt) bleibt unverändert und wird
// nicht durch die SPEC-013-Startwerte ersetzt.
func TestResolveWALRetentionThresholdsKeepsPositiveOverride(t *testing.T) {
	const warnOverride, errorOverride int64 = 32 * 1024, 512 * 1024
	warnBytes, errorBytes := resolveWALRetentionThresholds(warnOverride, errorOverride)
	if warnBytes != warnOverride {
		t.Fatalf("resolveWALRetentionThresholds(%d, %d) warnBytes = %d, wollen den Override unverändert", warnOverride, errorOverride, warnBytes)
	}
	if errorBytes != errorOverride {
		t.Fatalf("resolveWALRetentionThresholds(%d, %d) errorBytes = %d, wollen den Override unverändert", warnOverride, errorOverride, errorBytes)
	}
}

// TestClassifyWALRetentionThresholds belegt die `>`-Semantik aus `SPEC-013`
// („Warn > 100 MiB · Fehler > 1 GiB"): an der Schwelle selbst gilt noch die
// niedrigere Stufe — dieselbe Semantik wie `healthcheckVerdict`
// (`TestHealthcheckVerdictThreshold`). Rot färbende Mutation: `>` durch
// `>=` ersetzen — dann kippt die Klassifikation exakt an der Schwelle
// selbst um eine Stufe.
func TestClassifyWALRetentionThresholds(t *testing.T) {
	const warn, errorThreshold int64 = 1024, 8192
	cases := []struct {
		name  string
		bytes int64
		want  walRetentionLevel
	}{
		{"weit unterhalb", 0, walRetentionOK},
		{"an der Warnschwelle", warn, walRetentionOK},
		{"knapp über der Warnschwelle", warn + 1, walRetentionWarn},
		{"an der Fehlerschwelle", errorThreshold, walRetentionWarn},
		{"knapp über der Fehlerschwelle", errorThreshold + 1, walRetentionError},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := classifyWALRetention(c.bytes, warn, errorThreshold); got != c.want {
				t.Fatalf("classifyWALRetention(%d, %d, %d) = %d, wollen %d", c.bytes, warn, errorThreshold, got, c.want)
			}
		})
	}
}

// fakeWALRetentionMeasurer liefert eine vorab festgelegte Folge von
// Byte-Werten je `Measure`-Aufruf — Stellvertreter für
// `*receive.WALRetentionChecker` über das `walRetentionMeasurer`-Interface.
// Nach dem letzten Wert blockiert `Measure`, bis der Kontext endet: die
// Prüf-Goroutine hat dann bereits reagiert (Log, `stopStream`, `fault.set`)
// und braucht keinen weiteren Tick.
type fakeWALRetentionMeasurer struct {
	mu     sync.Mutex
	values []int64
	calls  int
	block  chan struct{}
}

func (f *fakeWALRetentionMeasurer) Measure(ctx context.Context) (int64, error) {
	f.mu.Lock()
	idx := f.calls
	f.calls++
	f.mu.Unlock()
	if idx < len(f.values) {
		return f.values[idx], nil
	}
	select {
	case <-f.block:
	case <-ctx.Done():
	}
	return 0, ctx.Err()
}

// recordingLog trägt jeden Log-Aufruf zur Prüfung, welche Stufe geloggt
// wurde — `Debug` bleibt ungenutzt (kein Aufrufer dieses Tests nutzt ihn).
type recordingLog struct {
	mu    sync.Mutex
	warns int
	errs  int
	infos int
}

func (r *recordingLog) Debug(ctx context.Context, msg string, args ...any) {}
func (r *recordingLog) Info(ctx context.Context, msg string, args ...any) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.infos++
}
func (r *recordingLog) Warn(ctx context.Context, msg string, args ...any) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.warns++
}
func (r *recordingLog) Error(ctx context.Context, msg string, args ...any) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.errs++
}

var _ outbound.LogPort = (*recordingLog)(nil)

// TestRunWALRetentionCheckContinuesBelowAndAtWarnLevel belegt die
// „kontrollierte Fortsetzung" zwischen Warn- und Fehlerschwelle
// (`SPEC-008`): mehrere Ticks über der Warnschwelle protokollieren jeweils
// eine Warnung, brechen die Schleife aber nicht ab — `stopStream` wird
// nicht aufgerufen, `fault` bleibt leer. Rot färbende Mutation: den
// `walRetentionWarn`-Fall ebenfalls `stopStream`/`fault.set` aufrufen
// lassen — dann bricht dieser Test, weil die Schleife nach dem ersten
// Warn-Tick endet statt weiterzulaufen.
func TestRunWALRetentionCheckContinuesBelowAndAtWarnLevel(t *testing.T) {
	measurer := &fakeWALRetentionMeasurer{values: []int64{0, 2000, 2000}, block: make(chan struct{})}
	log := &recordingLog{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	streamCtx, stopStream := context.WithCancel(context.Background())
	defer stopStream()
	var fault walRetentionFault

	done := make(chan struct{})
	go func() {
		defer close(done)
		runWALRetentionCheck(ctx, measurer, log, time.Millisecond, 1024, 8192, stopStream, &fault)
	}()

	deadline := time.After(time.Second)
	for {
		log.mu.Lock()
		warns := log.warns
		log.mu.Unlock()
		if warns >= 2 {
			break
		}
		select {
		case <-deadline:
			t.Fatal("erwartete mindestens zwei Warnungen (zwei Ticks über der Warnschwelle) nicht innerhalb 1s")
		case <-time.After(time.Millisecond):
		}
	}
	if streamCtx.Err() != nil {
		t.Fatal("stopStream wurde aufgerufen, obwohl der WAL-Rückstand nur die Warnschwelle überschreitet")
	}
	if err := fault.get(); err != nil {
		t.Fatalf("fault = %v, wollen nil (Warnschwelle bricht nicht ab)", err)
	}
	cancel()
	<-done
}

// TestRunWALRetentionCheckStopsStreamAboveErrorThreshold belegt den
// kontrollierten Abbruch oberhalb der Fehlerschwelle (`SPEC-008`,
// `ADR-0049`): der Zug protokolliert einen Fehler, ruft `stopStream` genau
// einmal auf und setzt `fault` auf einen `outbound.ErrReplication`-Fehler.
// Rot färbende Mutation: den `walRetentionError`-Fall wie `walRetentionWarn`
// behandeln (nur loggen, kein `stopStream`/`fault.set`) — dann bleibt
// `streamCtx` hier unbeendet und dieser Test läuft in den 1s-Timeout.
func TestRunWALRetentionCheckStopsStreamAboveErrorThreshold(t *testing.T) {
	measurer := &fakeWALRetentionMeasurer{values: []int64{9000}, block: make(chan struct{})}
	log := &recordingLog{}
	streamCtx, stopStream := context.WithCancel(context.Background())
	defer stopStream()
	var fault walRetentionFault

	done := make(chan struct{})
	go func() {
		defer close(done)
		runWALRetentionCheck(context.Background(), measurer, log, time.Millisecond, 1024, 8192, stopStream, &fault)
	}()

	select {
	case <-streamCtx.Done():
	case <-time.After(time.Second):
		t.Fatal("stopStream wurde nach Überschreiten der Fehlerschwelle nicht innerhalb 1s aufgerufen")
	}
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("runWALRetentionCheck ist nach dem Schwellen-Fehler nicht zurückgekehrt")
	}
	got := fault.get()
	if got == nil {
		t.Fatal("fault = nil, wollen einen gesetzten WAL-Schwellen-Fehler")
	}
	if !errors.Is(got, outbound.ErrReplication) {
		t.Fatalf("fault = %v, wollen einen Fehler der Klasse replication (outbound.ErrReplication)", got)
	}
	log.mu.Lock()
	errs := log.errs
	log.mu.Unlock()
	if errs != 1 {
		t.Fatalf("Error-Log-Aufrufe = %d, wollen 1", errs)
	}
}

// TestMergeStreamAndWALFaultOutcomePrioritizesStreamError belegt die
// Sentinel-Trennung aus `ADR-0049`(a) auf Ebene der Rückgabewert-Priorität
// (`slice-026` §6, Regressionsrisiko): eine Stream-Ordnungs-Verletzung
// erreicht `Run`s Rückgabewert unverändert, selbst wenn zufällig auch ein
// WAL-Schwellen-Fehler aufgelaufen ist — der neue Fortsetzungspfad
// überschreibt oder verdeckt sie nie. Rot färbende Mutation: die Priorität
// in `mergeStreamAndWALFaultOutcome` umdrehen (`fault.get()` zuerst prüfen)
// — dann liefert dieser Test den WAL-Schwellen-Fehler statt der
// Stream-Ordnungs-Verletzung.
func TestMergeStreamAndWALFaultOutcomePrioritizesStreamError(t *testing.T) {
	var fault walRetentionFault
	fault.set(fmt.Errorf("%w: WAL-Rückstand über Fehlerschwelle", outbound.ErrReplication))

	streamOrderSentinels := []error{
		mapper.ErrChangeWithoutBegin,
		mapper.ErrCommitWithoutBegin,
		mapper.ErrBeginWithoutCommit,
	}
	for _, sentinel := range streamOrderSentinels {
		got := mergeStreamAndWALFaultOutcome(sentinel, &fault)
		if !errors.Is(got, sentinel) {
			t.Fatalf("mergeStreamAndWALFaultOutcome(%v, gesetzter WAL-Fault) = %v, wollen die Stream-Ordnungs-Verletzung unverändert", sentinel, got)
		}
		if errors.Is(got, outbound.ErrReplication) && !errors.Is(sentinel, outbound.ErrReplication) {
			t.Fatalf("mergeStreamAndWALFaultOutcome(%v, …) trägt zusätzlich outbound.ErrReplication — der WAL-Fault hat die Stream-Ordnungs-Verletzung verdeckt", sentinel)
		}
	}
}

// TestMergeStreamAndWALFaultOutcomeFallsBackToFaultOnRegularStreamEnd
// belegt die Gegenseite: erst ein regulärer Stream-Abschluss (`nil`, über
// `stopStream` ausgelöst) lässt den aufgelaufenen WAL-Schwellen-Fehler
// durch — genau der Fall, den `Run` nach einem Fehlerschwellen-Abbruch
// erreicht.
func TestMergeStreamAndWALFaultOutcomeFallsBackToFaultOnRegularStreamEnd(t *testing.T) {
	var fault walRetentionFault
	walErr := fmt.Errorf("%w: WAL-Rückstand über Fehlerschwelle", outbound.ErrReplication)
	fault.set(walErr)

	got := mergeStreamAndWALFaultOutcome(nil, &fault)
	if !errors.Is(got, outbound.ErrReplication) {
		t.Fatalf("mergeStreamAndWALFaultOutcome(nil, gesetzter WAL-Fault) = %v, wollen den WAL-Schwellen-Fehler", got)
	}

	if got := mergeStreamAndWALFaultOutcome(nil, &walRetentionFault{}); got != nil {
		t.Fatalf("mergeStreamAndWALFaultOutcome(nil, leerer WAL-Fault) = %v, wollen nil", got)
	}
}
