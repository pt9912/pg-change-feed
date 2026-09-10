// Package register trägt den RegisterConsumer Use Case (`ARC-002`,
// `ADR-0028`): die Registrierung eines benannten Consumers trägt die
// Consumer-Zeile (`LH-FA-CON-001`).
package register

import (
	"context"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// RegisterConsumerCommand und RegisterConsumerResult sind die
// Transport-Typen des RegisterConsumer Use Cases (`ADR-0042`): ihre
// Definition liegt am Inbound-Port — der Port trägt seinen Vertrag
// einschließlich der Transport-Typen, Ports referenzieren nur die Domain —;
// die Use-Case-Adressen sind Aliase desselben Typs.
type (
	RegisterConsumerCommand = inbound.RegisterConsumerCommand
	RegisterConsumerResult  = inbound.RegisterConsumerResult
)

// RegisterConsumerService implementiert `inbound.RegisterConsumerUseCase`
// und trägt die Registrierung über den `ConsumerStatePort`: die Kennung
// und der Name laufen über den Domänen-Konstruktor — eine leere Kennung
// endet über die Domänen-Invariante, ohne den Zustand zu berühren.
type RegisterConsumerService struct {
	state outbound.ConsumerStatePort
}

// NewRegisterConsumerService verdrahtet den Use Case mit dem
// Consumer-State-Port.
func NewRegisterConsumerService(state outbound.ConsumerStatePort) *RegisterConsumerService {
	return &RegisterConsumerService{state: state}
}

var _ inbound.RegisterConsumerUseCase = (*RegisterConsumerService)(nil)

// Register registriert den Consumer idempotent (`LH-FA-CON-001` Boundary):
// eine bereits registrierte Kennung bleibt unverändert; die Rückkehr
// meldet den Ausgang über `AlreadyRegistered`.
func (s *RegisterConsumerService) Register(ctx context.Context, command RegisterConsumerCommand) (RegisterConsumerResult, error) {
	consumer, err := model.NewConsumer(command.Consumer, command.Name)
	if err != nil {
		return RegisterConsumerResult{}, err
	}
	registered, err := s.state.Register(ctx, consumer)
	if err != nil {
		return RegisterConsumerResult{}, err
	}
	return RegisterConsumerResult{Consumer: consumer, AlreadyRegistered: !registered}, nil
}
