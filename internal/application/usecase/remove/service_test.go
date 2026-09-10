package remove_test

import (
	"context"
	stderrors "errors"
	"testing"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/remove"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// fakeState trägt den ConsumerStatePort als Fake (`ADR-0030`); die
// Entfernung meldet ihren Ausgang, die Stubs der übrigen Operationen
// tragen die Schnittstelle.
type fakeState struct {
	exists    bool
	removeErr error

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
	return model.ConsumerPosition{}, nil
}

func (f *fakeState) Acknowledge(ctx context.Context, position model.ConsumerPosition) (model.ConsumerPosition, error) {
	f.acknowledgeCalls++
	return position, nil
}

func (f *fakeState) Remove(ctx context.Context, consumer model.ConsumerID) (bool, error) {
	f.removeCalls++
	f.lastConsumer = consumer
	if f.removeErr != nil {
		return false, f.removeErr
	}
	return f.exists, nil
}

// TestRemoveHappyPath trägt die administrative Entfernung
// (`LH-FA-CON-006` Happy Path): der Consumer führt danach keine Rolle
// mehr in der Retention — die Rückkehr meldet die Entfernung.
func TestRemoveHappyPath(t *testing.T) {
	fake := &fakeState{exists: true}
	service := remove.NewRemoveConsumerService(fake)

	result, err := service.Remove(context.Background(), inbound.RemoveConsumerCommand{
		Consumer: "c-1",
	})
	if err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if !result.Removed {
		t.Fatalf("Entfernung meldet Removed=false")
	}
	if fake.removeCalls != 1 || fake.lastConsumer != "c-1" {
		t.Fatalf("Entfernungs-Züge: %d, Consumer %q", fake.removeCalls, fake.lastConsumer)
	}
}

// TestRemoveIdempotent trägt den Boundary-Pfad (`LH-FA-CON-006`): ein
// Consumer ohne Zeile bleibt ohne Wirkung; die Rückkehr meldet den
// Ausgang, ohne einen Fehler zu tragen.
func TestRemoveIdempotent(t *testing.T) {
	fake := &fakeState{}
	service := remove.NewRemoveConsumerService(fake)

	result, err := service.Remove(context.Background(), inbound.RemoveConsumerCommand{
		Consumer: "c-1",
	})
	if err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if result.Removed {
		t.Fatalf("Entfernung ohne Bestand meldet Removed=true")
	}
	if fake.removeCalls != 1 {
		t.Fatalf("Entfernungs-Züge: %d", fake.removeCalls)
	}
}

// TestRemoveInvalidCommand trägt die Kennungs-Grenze: eine leere Kennung
// endet über die Domänen-Invariante (`ADR-0029`), ohne den
// Zustands-Port zu berühren.
func TestRemoveInvalidCommand(t *testing.T) {
	fake := &fakeState{}
	service := remove.NewRemoveConsumerService(fake)

	_, err := service.Remove(context.Background(), inbound.RemoveConsumerCommand{})
	if !stderrors.Is(err, domainerrors.ErrEmptyIdentifier) {
		t.Fatalf("leere Kennung: %v (Erwartung: ErrEmptyIdentifier)", err)
	}
	if fake.removeCalls != 0 {
		t.Fatalf("ungültige Eingabe trägt Entfernungs-Züge: %d", fake.removeCalls)
	}
}
