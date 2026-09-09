package outbound_test

import (
	"testing"

	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// fakeClock ist die Fake-Uhr der Tests (`ADR-0040`): sie trägt einen
// gesetzten Zeitpunkt und liefert ihn deterministisch, bis der Test sie
// vorstellt.
type fakeClock struct {
	now model.TimePoint
}

func (f *fakeClock) Now() model.TimePoint {
	return f.now
}

// ClockPort erfüllt die Port-Kante (`ARC-004`): die Rückgabe ist ein
// Domänenwert, kein `time`-Typ.
func TestClockPortReturnsDomainTimePoint(t *testing.T) {
	clock := &fakeClock{now: model.NewTimePoint(1000)}
	var port outbound.ClockPort = clock
	if got := port.Now(); got != model.NewTimePoint(1000) {
		t.Fatalf("Now = %d, wollen 1000", got.UnixNanos)
	}
}

// Die Fake-Uhr ist deterministisch: zwei Aufrufe ohne Vorstellen tragen
// denselben Zeitpunkt; nach dem Vorstellen tragen sie den neuen.
func TestFakeClockIsDeterministicUntilAdvanced(t *testing.T) {
	clock := &fakeClock{now: model.NewTimePoint(1000)}
	if clock.Now() != clock.Now() {
		t.Fatal("zwei Aufrufe tragen verschiedene Zeitpunkte")
	}
	advanced := model.NewTimePoint(2000)
	clock.now = advanced
	if clock.Now() != advanced {
		t.Fatalf("Now = %d, wollen %d", clock.Now().UnixNanos, advanced.UnixNanos)
	}
	if !clock.Now().After(model.NewTimePoint(1000)) {
		t.Fatal("vorgestellte Uhr meldet ihren Zeitpunkt nicht als späteren")
	}
}
