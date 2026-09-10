// Package remove trägt den RemoveConsumer Use Case (`ARC-002`): die
// administrative Entfernung eines Consumers (`LH-FA-CON-006`) — der
// Consumer trägt danach keine Rolle in der Retention
// (`LH-FA-RET-004`).
package remove

import (
	"context"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

// RemoveConsumerCommand und RemoveConsumerResult sind die Transport-Typen
// des RemoveConsumer Use Cases (`ADR-0042`): ihre Definition liegt am
// Inbound-Port — der Port trägt seinen Vertrag einschließlich der
// Transport-Typen, Ports referenzieren nur die Domain —; die Use-Case-
// Adressen sind Aliase desselben Typs.
type (
	RemoveConsumerCommand = inbound.RemoveConsumerCommand
	RemoveConsumerResult  = inbound.RemoveConsumerResult
)

// RemoveConsumerService implementiert `inbound.RemoveConsumerUseCase` und
// trägt die Entfernung über den `ConsumerStatePort`: die bestätigte
// Position geht mit der Zeile (dokumentiertes Verhalten,
// `LH-FA-CON-006` Boundary), die Entfernung löscht keine Changes.
type RemoveConsumerService struct {
	state outbound.ConsumerStatePort
}

// NewRemoveConsumerService verdrahtet den Use Case mit dem
// Consumer-State-Port.
func NewRemoveConsumerService(state outbound.ConsumerStatePort) *RemoveConsumerService {
	return &RemoveConsumerService{state: state}
}

var _ inbound.RemoveConsumerUseCase = (*RemoveConsumerService)(nil)

// Remove entfernt den Consumer idempotent (`LH-FA-CON-006`): die Rückkehr
// meldet, ob der Consumer entfernt war; ein erneuter Aufruf bleibt ohne
// Wirkung.
func (s *RemoveConsumerService) Remove(ctx context.Context, command RemoveConsumerCommand) (RemoveConsumerResult, error) {
	if command.Consumer == "" {
		return RemoveConsumerResult{}, domainerrors.ErrEmptyIdentifier
	}
	removed, err := s.state.Remove(ctx, command.Consumer)
	if err != nil {
		return RemoveConsumerResult{}, err
	}
	return RemoveConsumerResult{Removed: removed}, nil
}
