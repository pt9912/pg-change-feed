package register_test

import (
	"context"
	stderrors "errors"
	"testing"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/register"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// fakeState trägt den ConsumerStatePort als Fake (`ADR-0030`); die
// Registrierung meldet ihren Aufruf, die Stubs der übrigen Operationen
// tragen die Schnittstelle.
type fakeState struct {
	alreadyRegistered bool
	registerErr       error

	registerCalls    int
	lastConsumer     model.Consumer
	positionCalls    int
	acknowledgeCalls int
	removeCalls      int
}

func (f *fakeState) Register(ctx context.Context, consumer model.Consumer) (bool, error) {
	f.registerCalls++
	f.lastConsumer = consumer
	if f.registerErr != nil {
		return false, f.registerErr
	}
	return !f.alreadyRegistered, nil
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
	return true, nil
}

// TestRegisterHappyPath trägt die Registrierung (`LH-FA-CON-001` Happy
// Path): der Consumer ist fortan registriert — die Rückkehr liest die
// Kennung und den Namen.
func TestRegisterHappyPath(t *testing.T) {
	fake := &fakeState{}
	service := register.NewRegisterConsumerService(fake)

	result, err := service.Register(context.Background(), inbound.RegisterConsumerCommand{
		Consumer: "c-1", Name: "Lese-Consumer",
	})
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	if result.AlreadyRegistered {
		t.Fatalf("neue Registrierung meldet AlreadyRegistered")
	}
	if result.Consumer.ID != "c-1" || result.Consumer.Name != "Lese-Consumer" {
		t.Fatalf("registrierter Consumer: %+v", result.Consumer)
	}
	if fake.registerCalls != 1 || fake.lastConsumer.ID != "c-1" {
		t.Fatalf("Registrierungs-Züge: %d, Consumer %v", fake.registerCalls, fake.lastConsumer)
	}
}

// TestRegisterIdempotent trägt den Boundary-Pfad (`LH-FA-CON-001`): die
// bereits registrierte Kennung meldet das Ergebnis, ohne den Stand zu
// ändern.
func TestRegisterIdempotent(t *testing.T) {
	fake := &fakeState{alreadyRegistered: true}
	service := register.NewRegisterConsumerService(fake)

	result, err := service.Register(context.Background(), inbound.RegisterConsumerCommand{
		Consumer: "c-1", Name: "Lese-Consumer",
	})
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	if !result.AlreadyRegistered {
		t.Fatalf("erneute Registrierung meldet AlreadyRegistered=false")
	}
	if fake.registerCalls != 1 {
		t.Fatalf("Registrierungs-Züge: %d", fake.registerCalls)
	}
}

// TestRegisterInvalidCommand trägt die Kennungs-Grenzen: eine leere
// Kennung und ein leerer Name enden über die Domänen-Invariante
// (`ADR-0029`), ohne den Zustands-Port zu berühren.
func TestRegisterInvalidCommand(t *testing.T) {
	fake := &fakeState{}
	service := register.NewRegisterConsumerService(fake)

	_, err := service.Register(context.Background(), inbound.RegisterConsumerCommand{
		Name: "Lese-Consumer",
	})
	if !stderrors.Is(err, domainerrors.ErrEmptyIdentifier) {
		t.Fatalf("leere Kennung: %v (Erwartung: ErrEmptyIdentifier)", err)
	}
	_, err = service.Register(context.Background(), inbound.RegisterConsumerCommand{
		Consumer: "c-1",
	})
	if !stderrors.Is(err, domainerrors.ErrEmptyIdentifier) {
		t.Fatalf("leerer Name: %v (Erwartung: ErrEmptyIdentifier)", err)
	}
	if fake.registerCalls != 0 {
		t.Fatalf("ungültige Eingabe trägt Registrierungs-Züge: %d", fake.registerCalls)
	}
}
