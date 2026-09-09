package model

import (
	"testing"
)

// ADR-0040: die Zeit läuft als Wert durch die Domain — Before/After und
// Arithmetik tragen den Zustand, ohne das Paket `time` zu importieren.
func TestTimePointArithmetic(t *testing.T) {
	base := NewTimePoint(1000)
	later := NewTimePoint(2000)
	if base.Unset() || later.Unset() {
		t.Fatal("belegte Zeitpunkte melden „nicht gesetzt“")
	}
	if !later.After(base) || !base.Before(later) {
		t.Fatalf("Ordnung verletzt: %d vs %d", base.UnixNanos, later.UnixNanos)
	}
	if got := base.Add(Duration{Nanos: 500}); got.UnixNanos != 1500 {
		t.Fatalf("Add = %d, wollen 1500", got.UnixNanos)
	}
	if got := later.Sub(base); got.Nanos != 1000 {
		t.Fatalf("Sub = %d, wollen 1000", got.Nanos)
	}
	if !NewTimePoint(0).Unset() {
		t.Fatal("Nullwert meldet nicht „nicht gesetzt“")
	}
}
