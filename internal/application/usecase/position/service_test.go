package position_test

import (
	"context"
	stderrors "errors"
	"testing"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/position"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// fakeState trägt den ConsumerStatePort als Fake (`ADR-0030`); der Lese
// meldet seinen Stand, die Stubs der übrigen Operationen tragen die
// Schnittstelle.
type fakeState struct {
	stored model.ConsumerPosition

	registerCalls    int
	positionCalls    int
	acknowledgeCalls int
	removeCalls      int
	lastConsumer     model.ConsumerID
}

func (f *fakeState) Register(ctx context.Context, consumer model.Consumer) (bool, error) {
	f.registerCalls++
	return true, nil
}

func (f *fakeState) Position(ctx context.Context, consumer model.ConsumerID) (model.ConsumerPosition, error) {
	f.positionCalls++
	f.lastConsumer = consumer
	return f.stored, nil
}

func (f *fakeState) Acknowledge(ctx context.Context, position model.ConsumerPosition) (model.ConsumerPosition, error) {
	f.acknowledgeCalls++
	return position, nil
}

func (f *fakeState) Remove(ctx context.Context, consumer model.ConsumerID) (bool, error) {
	f.removeCalls++
	return true, nil
}

// TestPositionHappyPath trägt den Positions-Lese (`LH-FA-CON-003` Happy
// Path): der Consumer erhält seine bestätigte Position.
func TestPositionHappyPath(t *testing.T) {
	stored, err := model.NewConsumerPosition("c-1")
	if err != nil {
		t.Fatalf("NewConsumerPosition: %v", err)
	}
	pos, err := model.NewSourcePosition("src-1", 100)
	if err != nil {
		t.Fatalf("NewSourcePosition: %v", err)
	}
	stored, err = stored.Advance(pos)
	if err != nil {
		t.Fatalf("Advance: %v", err)
	}
	fake := &fakeState{stored: stored}
	service := position.NewGetConsumerPositionService(fake)

	result, err := service.Position(context.Background(), inbound.GetConsumerPositionQuery{Consumer: "c-1"})
	if err != nil {
		t.Fatalf("Position: %v", err)
	}
	if !result.Position.Acknowledged() || result.Position.Position.Offset != 100 {
		t.Fatalf("bestätigte Position: %+v", result.Position)
	}
	if fake.positionCalls != 1 || fake.lastConsumer != "c-1" {
		t.Fatalf("Lese-Züge: %d, Consumer %q", fake.positionCalls, fake.lastConsumer)
	}
}

// TestPositionWithoutAcknowledgement trägt den Boundary-Pfad
// (`LH-FA-CON-005`): die definierte Anfangsposition liest sich als
// Nullwert — der Lese trägt ihn, ohne einen Fehler zu melden.
func TestPositionWithoutAcknowledgement(t *testing.T) {
	fake := &fakeState{}
	service := position.NewGetConsumerPositionService(fake)

	result, err := service.Position(context.Background(), inbound.GetConsumerPositionQuery{Consumer: "c-1"})
	if err != nil {
		t.Fatalf("Position: %v", err)
	}
	if result.Position.Acknowledged() {
		t.Fatalf("Position ohne Bestätigung liest sich bestätigt: %+v", result.Position)
	}
	if fake.positionCalls != 1 {
		t.Fatalf("Lese-Züge: %d", fake.positionCalls)
	}
}

// TestPositionInvalidCommand trägt die Kennungs-Grenze: eine leere
// Kennung endet über die Domänen-Invariante (`ADR-0029`), ohne den
// Zustands-Port zu berühren.
func TestPositionInvalidCommand(t *testing.T) {
	fake := &fakeState{}
	service := position.NewGetConsumerPositionService(fake)

	_, err := service.Position(context.Background(), inbound.GetConsumerPositionQuery{})
	if !stderrors.Is(err, domainerrors.ErrEmptyIdentifier) {
		t.Fatalf("leere Kennung: %v (Erwartung: ErrEmptyIdentifier)", err)
	}
	if fake.positionCalls != 0 {
		t.Fatalf("ungültige Eingabe trägt Lese-Züge: %d", fake.positionCalls)
	}
}
