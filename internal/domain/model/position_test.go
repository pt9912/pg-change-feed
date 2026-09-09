package model

import (
	"math/rand"
	"sort"
	"testing"

	stderrors "errors"

	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

// LH-FA-DAT-004, Happy Path: die Positionen zweier committed Changes
// bestimmen die Ordnung; sie entspricht der logischen Reihenfolge
// (LH-FA-CAP-004).
func TestLHFADAT004PositionsAreSortable(t *testing.T) {
	first, err := NewSourcePosition("src-1", 100)
	if err != nil {
		t.Fatalf("erste Position: %v", err)
	}
	second, err := NewSourcePosition("src-1", 200)
	if err != nil {
		t.Fatalf("zweite Position: %v", err)
	}
	if !first.Before(second) || !second.After(first) {
		t.Fatalf("Ordnung verletzt: Offset %d liegt vor %d, Compare %d", first.Offset, second.Offset, first.Compare(second))
	}
}

// Die Position-Invarianten (`SPEC-003`) scheitern am Konstruktor: keine
// Quelle, kein Offset.
func TestNewSourcePositionRejectsInvariantViolations(t *testing.T) {
	if _, err := NewSourcePosition("", 100); !stderrors.Is(err, domainerrors.ErrEmptyIdentifier) {
		t.Fatalf("Fehler = %v, wollen ErrEmptyIdentifier", err)
	}
	if _, err := NewSourcePosition("src-1", 0); !stderrors.Is(err, domainerrors.ErrInvalidPosition) {
		t.Fatalf("Fehler = %v, wollen ErrInvalidPosition", err)
	}
}

// LH-FA-REA-004: die Sortierung ist bei identischer Eingabe
// deterministisch — derselbe Shuffle zweimal ergibt dieselbe Reihenfolge.
func TestSourcePositionOrderingIsDeterministic(t *testing.T) {
	offsets := []uint64{7, 3, 9, 1, 5, 8, 2, 6, 4, 10}
	positions := make([]SourcePosition, 0, len(offsets))
	for _, offset := range offsets {
		position, err := NewSourcePosition("src-1", offset)
		if err != nil {
			t.Fatalf("Position %d: %v", offset, err)
		}
		positions = append(positions, position)
	}
	shuffle := func() []SourcePosition {
		shuffled := append([]SourcePosition(nil), positions...)
		rng := rand.New(rand.NewSource(42))
		rng.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })
		sort.Slice(shuffled, func(i, j int) bool { return shuffled[i].Before(shuffled[j]) })
		return shuffled
	}
	firstRun, secondRun := shuffle(), shuffle()
	for i := range firstRun {
		if firstRun[i] != secondRun[i] {
			t.Fatalf("Sortierung nicht deterministisch: Stelle %d trägt %v vs %v", i, firstRun[i].Offset, secondRun[i].Offset)
		}
	}
}

// Der Quell-Vergleich in Compare ist ein stabiler Sortier-Schlüssel, keine
// fachliche Ordnung über Quellen hinweg (`ADR-0005`): gleiche Quelle
// ordnet nach Offset.
func TestSourcePositionCompareAcrossSourcesIsStableTiebreak(t *testing.T) {
	first, err := NewSourcePosition("src-a", 500)
	if err != nil {
		t.Fatalf("Position a: %v", err)
	}
	second, err := NewSourcePosition("src-b", 100)
	if err != nil {
		t.Fatalf("Position b: %v", err)
	}
	if c := first.Compare(second); c >= 0 {
		t.Fatalf("Tiebreak nicht stabil: Quelle \"src-a\" sortiert nicht vor \"src-b\" (Compare %d)", c)
	}
	if first.After(second) {
		t.Fatal("Compare und After widersprechen sich")
	}
}
