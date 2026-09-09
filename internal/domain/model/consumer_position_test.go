package model

import (
	"testing"

	stderrors "errors"

	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

// acknowledgedConsumerPosition liefert eine Consumer-Position mit
// bestätigtem Offset der Quelle "src-1" für Tests; sie entsteht über den
// Konstruktor und `Advance`, nicht über ein Struktur-Literal.
func acknowledgedConsumerPosition(t *testing.T, offset uint64) ConsumerPosition {
	t.Helper()
	consumer, err := NewConsumerPosition("c-1")
	if err != nil {
		t.Fatalf("Consumer-Position: %v", err)
	}
	advanced, err := consumer.Advance(mustSourcePosition(t, "src-1", offset))
	if err != nil {
		t.Fatalf("Bestätigung: %v", err)
	}
	return advanced
}

// ADR-0029, Regel 2: der Consumer-ACK verläuft regulär nur vorwärts —
// eine frühere Position scheitert, eine spätere wird bestätigt.
func TestConsumerAcknowledgesForwardOnly(t *testing.T) {
	current := acknowledgedConsumerPosition(t, 200)
	if _, err := current.Advance(mustSourcePosition(t, "src-1", 100)); !stderrors.Is(err, domainerrors.ErrPositionRegression) {
		t.Fatalf("Fehler = %v, wollen ErrPositionRegression", err)
	}
	advanced, err := current.Advance(mustSourcePosition(t, "src-1", 300))
	if err != nil {
		t.Fatalf("vorwärts bestätigen: %v", err)
	}
	if !advanced.Position.After(current.Position) {
		t.Fatalf("Position rückt nicht vor: %d bleibt %d", advanced.Position.Offset, current.Position.Offset)
	}
}

// LH-FA-CON-004, Boundary: die Wiederholung derselben Position ist
// idempotent — der Stand bleibt, kein Fehler.
func TestLHFACON004RepeatAcknowledgementIsIdempotent(t *testing.T) {
	current := acknowledgedConsumerPosition(t, 200)
	advanced, err := current.Advance(current.Position)
	if err != nil {
		t.Fatalf("Wiederholung derselben Position: %v", err)
	}
	if advanced != current {
		t.Fatalf("Wiederholung ändert den Stand: %v vs %v", advanced, current)
	}
}

// ADR-0005: die bestätigte Quelle bleibt gebunden — die Position einer
// anderen Quelle ist keine Bestätigungs-Größe für einen Consumer, dessen
// Bestätigung an "src-1" gebunden ist.
func TestConsumerPositionRejectsOtherSourceAfterBinding(t *testing.T) {
	current := acknowledgedConsumerPosition(t, 200)
	if _, err := current.Advance(mustSourcePosition(t, "src-2", 300)); !stderrors.Is(err, domainerrors.ErrSourceMismatch) {
		t.Fatalf("Fehler = %v, wollen ErrSourceMismatch", err)
	}
}

// LH-FA-CON-003: der erste Consumer-ACK gilt ohne Vorgänger und bindet die
// Quelle seiner Position; danach meldet die Position „bestätigt“.
func TestConsumerPositionFirstAcknowledgementBindsSource(t *testing.T) {
	unacknowledged, err := NewConsumerPosition("c-1")
	if err != nil {
		t.Fatalf("Consumer-Position: %v", err)
	}
	if unacknowledged.Acknowledged() {
		t.Fatal("ohne Bestätigung meldet die Position „bestätigt“")
	}
	advanced, err := unacknowledged.Advance(mustSourcePosition(t, "src-1", 100))
	if err != nil {
		t.Fatalf("erste Bestätigung: %v", err)
	}
	if !advanced.Acknowledged() {
		t.Fatal("bestätigte Position meldet unbestätigt")
	}
	if advanced.Position.SourceID != "src-1" {
		t.Fatalf("erste Bestätigung bindet die Quelle nicht: %q", advanced.Position.SourceID)
	}
}

// Die Consumer-Position verlangt eine nichtleere Kennung.
func TestNewConsumerPositionRejectsEmptyIdentifier(t *testing.T) {
	if _, err := NewConsumerPosition(""); !stderrors.Is(err, domainerrors.ErrEmptyIdentifier) {
		t.Fatalf("Fehler = %v, wollen ErrEmptyIdentifier", err)
	}
}
