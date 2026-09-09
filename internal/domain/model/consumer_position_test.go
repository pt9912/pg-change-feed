package model

import (
	"testing"

	stderrors "errors"

	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

// acknowledgedConsumerPosition liefert eine Consumer-Position der Quelle
// "src-1" mit bestätigtem Offset für Tests.
func acknowledgedConsumerPosition(t *testing.T, offset uint64) ConsumerPosition {
	t.Helper()
	position, err := NewSourcePosition("src-1", offset)
	if err != nil {
		t.Fatalf("Position: %v", err)
	}
	return ConsumerPosition{ConsumerID: "c-1", SourceID: "src-1", Position: position}
}

// ADR-0029, Regel 2: der Consumer-ACK verläuft regulär nur vorwärts —
// eine frühere Position scheitert, eine spätere wird bestätigt.
func TestConsumerAcknowledgesForwardOnly(t *testing.T) {
	current := acknowledgedConsumerPosition(t, 200)
	backwards, err := NewSourcePosition("src-1", 100)
	if err != nil {
		t.Fatalf("frühere Position: %v", err)
	}
	if _, err := current.Advance(backwards); !stderrors.Is(err, domainerrors.ErrPositionRegression) {
		t.Fatalf("Fehler = %v, wollen ErrPositionRegression", err)
	}
	forwards, err := NewSourcePosition("src-1", 300)
	if err != nil {
		t.Fatalf("spätere Position: %v", err)
	}
	advanced, err := current.Advance(forwards)
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

// ADR-0005: die Position gehört zur Quelle des Consumers — eine Position
// einer anderen Quelle ist keine Bestätigungs-Größe.
func TestConsumerPositionRejectsOtherSource(t *testing.T) {
	current := acknowledgedConsumerPosition(t, 200)
	otherSource, err := NewSourcePosition("src-2", 300)
	if err != nil {
		t.Fatalf("fremde Position: %v", err)
	}
	if _, err := current.Advance(otherSource); !stderrors.Is(err, domainerrors.ErrSourceMismatch) {
		t.Fatalf("Fehler = %v, wollen ErrSourceMismatch", err)
	}
}

// Der erste Consumer-ACK gilt ohne Vorgänger und bindet die Quelle;
// danach meldet die Position „bestätigt“ (LH-FA-CON-003).
func TestConsumerPositionFirstAcknowledgementNeedsNoPredecessor(t *testing.T) {
	unacknowledged, err := NewConsumerPosition("c-1", "src-1")
	if err != nil {
		t.Fatalf("Consumer-Position: %v", err)
	}
	if unacknowledged.Acknowledged() {
		t.Fatal("ohne Bestätigung meldet die Position „bestätigt“")
	}
	position, err := NewSourcePosition("src-1", 100)
	if err != nil {
		t.Fatalf("Position: %v", err)
	}
	advanced, err := unacknowledged.Advance(position)
	if err != nil {
		t.Fatalf("erste Bestätigung: %v", err)
	}
	if !advanced.Acknowledged() {
		t.Fatal("bestätigte Position meldet unbestätigt")
	}
	if _, err := unacknowledged.Advance(otherSourcePosition(t)); !stderrors.Is(err, domainerrors.ErrSourceMismatch) {
		t.Fatalf("unbestätigte Position akzeptiert eine fremde Quelle: %v", err)
	}
}

func otherSourcePosition(t *testing.T) SourcePosition {
	t.Helper()
	position, err := NewSourcePosition("src-2", 100)
	if err != nil {
		t.Fatalf("fremde Position: %v", err)
	}
	return position
}

// Die Consumer-Position verlangt nichtleere Kennung und Quelle.
func TestNewConsumerPositionRejectsEmptyIdentifiers(t *testing.T) {
	if _, err := NewConsumerPosition("", "src-1"); !stderrors.Is(err, domainerrors.ErrEmptyIdentifier) {
		t.Fatalf("Fehler = %v, wollen ErrEmptyIdentifier", err)
	}
	if _, err := NewConsumerPosition("c-1", ""); !stderrors.Is(err, domainerrors.ErrEmptyIdentifier) {
		t.Fatalf("Fehler = %v, wollen ErrEmptyIdentifier", err)
	}
}
