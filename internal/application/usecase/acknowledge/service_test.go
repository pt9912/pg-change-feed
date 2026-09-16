package acknowledge_test

import (
	"context"
	stderrors "errors"
	"testing"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/acknowledge"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// fakeState trägt den ConsumerStatePort als Fake (`ADR-0030`); die
// Bestätigung meldet ihren Aufruf und die Eingabe — die Monotonie und die
// Quell-Bindung trägt der Domänen-Vergleich am Port, nicht der Use Case.
type fakeState struct {
	stored       model.ConsumerPosition
	lastPosition model.ConsumerPosition

	// ackErrFor bindet den Port-Fehler an die Consumer-Kennung der
	// Bestätigung (LP2): der Fake scheitert genau für diese Kennung und
	// trägt jede andere — so folgt der Ausgang der Kommando-Eingabe und
	// nicht dem Fake.
	ackErrFor model.ConsumerID
	ackErr    error

	registerCalls    int
	positionCalls    int
	acknowledgeCalls int
	removeCalls      int
}

func (f *fakeState) Register(ctx context.Context, consumer model.Consumer) (bool, error) {
	f.registerCalls++
	return true, nil
}

func (f *fakeState) Position(ctx context.Context, consumer model.ConsumerID) (model.ConsumerPosition, error) {
	f.positionCalls++
	return f.stored, nil
}

func (f *fakeState) Acknowledge(ctx context.Context, position model.ConsumerPosition) (model.ConsumerPosition, error) {
	f.acknowledgeCalls++
	f.lastPosition = position
	if f.ackErrFor != "" && position.ConsumerID == f.ackErrFor {
		return model.ConsumerPosition{}, f.ackErr
	}
	return position, nil
}

func (f *fakeState) Remove(ctx context.Context, consumer model.ConsumerID) (bool, error) {
	f.removeCalls++
	return true, nil
}

func (f *fakeState) Positions(ctx context.Context, source model.SourceID) ([]model.ConsumerPosition, error) {
	return nil, nil
}

func position(t *testing.T, offset uint64) model.SourcePosition {
	t.Helper()
	p, err := model.NewSourcePosition("src-1", offset)
	if err != nil {
		t.Fatalf("NewSourcePosition: %v", err)
	}
	return p
}

// TestAcknowledgeHappyPath trägt die Bestätigung (`LH-FA-CON-004` Happy
// Path): die bestätigte Position gilt als verarbeitete Position — die
// Rückkehr liest den fortgeführten Stand.
func TestAcknowledgeHappyPath(t *testing.T) {
	fake := &fakeState{}
	service := acknowledge.NewAcknowledgeConsumerService(fake)

	result, err := service.Acknowledge(context.Background(), inbound.AcknowledgeConsumerCommand{
		Consumer: "c-1",
		Position: position(t, 100),
	})
	if err != nil {
		t.Fatalf("Acknowledge: %v", err)
	}
	if result.Position.Position.Offset != 100 || result.Position.Position.SourceID != "src-1" {
		t.Fatalf("bestätigte Position: %+v", result.Position)
	}
	if fake.acknowledgeCalls != 1 || fake.lastPosition.ConsumerID != "c-1" {
		t.Fatalf("Bestätigungs-Züge: %d, Eingabe %+v", fake.acknowledgeCalls, fake.lastPosition)
	}
}

// TestAcknowledgeInvalidCommand trägt die Kennungs- und Positions-Grenzen:
// eine leere Kennung und ein Offset ohne Wert enden über die
// Domänen-Invarianten (`ADR-0029`), ohne den Zustands-Port zu berühren.
func TestAcknowledgeInvalidCommand(t *testing.T) {
	fake := &fakeState{}
	service := acknowledge.NewAcknowledgeConsumerService(fake)

	_, err := service.Acknowledge(context.Background(), inbound.AcknowledgeConsumerCommand{
		Position: position(t, 100),
	})
	if !stderrors.Is(err, domainerrors.ErrEmptyIdentifier) {
		t.Fatalf("leere Kennung: %v (Erwartung: ErrEmptyIdentifier)", err)
	}
	_, err = service.Acknowledge(context.Background(), inbound.AcknowledgeConsumerCommand{
		Consumer: "c-1",
		Position: model.SourcePosition{SourceID: "src-1"},
	})
	if !stderrors.Is(err, domainerrors.ErrInvalidPosition) {
		t.Fatalf("Position ohne Offset: %v (Erwartung: ErrInvalidPosition)", err)
	}
	if fake.acknowledgeCalls != 0 {
		t.Fatalf("ungültige Eingabe trägt Bestätigungs-Züge: %d", fake.acknowledgeCalls)
	}
}

// TestAcknowledgeStateErrorFollowsConsumer trägt den Fehlerpfad des
// Zustands-Schreibens: ein Fehler des `ConsumerStatePort` wird unverändert
// durchgereicht. Die Ablehnung ist an den Eingabewert gebunden — derselbe
// Fake trägt für eine andere Consumer-Kennung denselben Aufruf ohne Fehler,
// der Ausgang folgt also dem Kommando und nicht dem Fake.
func TestAcknowledgeStateErrorFollowsConsumer(t *testing.T) {
	wantErr := stderrors.New("Consumer-State nicht erreichbar")
	fake := &fakeState{ackErrFor: "c-boom", ackErr: wantErr}
	service := acknowledge.NewAcknowledgeConsumerService(fake)

	_, err := service.Acknowledge(context.Background(), inbound.AcknowledgeConsumerCommand{
		Consumer: "c-boom",
		Position: position(t, 100),
	})
	if !stderrors.Is(err, wantErr) {
		t.Fatalf("Port-Fehler: %v (Erwartung: %v)", err, wantErr)
	}

	result, err := service.Acknowledge(context.Background(), inbound.AcknowledgeConsumerCommand{
		Consumer: "c-1",
		Position: position(t, 100),
	})
	if err != nil {
		t.Fatalf("Bestätigung der Kennung c-1: %v", err)
	}
	if result.Position.Position.Offset != 100 || fake.lastPosition.ConsumerID != "c-1" {
		t.Fatalf("Rückkehr %+v, Port-Eingabe %+v", result, fake.lastPosition)
	}
}
