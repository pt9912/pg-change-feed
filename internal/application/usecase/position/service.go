// Package position trägt den GetConsumerPosition Use Case (`ARC-002`):
// der Lese der bestätigten Verarbeitungsposition eines Consumers
// (`LH-FA-CON-003`, `LH-FA-CON-005`).
package position

import (
	"context"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

// GetConsumerPositionQuery und GetConsumerPositionResult sind die
// Transport-Typen des GetConsumerPosition Use Cases (`ADR-0042`): ihre
// Definition liegt am Inbound-Port — der Port trägt seinen Vertrag
// einschließlich der Transport-Typen, Ports referenzieren nur die Domain —;
// die Use-Case-Adressen sind Aliase desselben Typs.
type (
	GetConsumerPositionQuery  = inbound.GetConsumerPositionQuery
	GetConsumerPositionResult = inbound.GetConsumerPositionResult
)

// GetConsumerPositionService implementiert
// `inbound.GetConsumerPositionUseCase` und liest die Position über den
// `ConsumerStatePort`: der Nullwert liest die definierte Anfangsposition
// eines Consumers ohne Bestätigung (`LH-FA-CON-005` Boundary).
type GetConsumerPositionService struct {
	state outbound.ConsumerStatePort
}

// NewGetConsumerPositionService verdrahtet den Use Case mit dem
// Consumer-State-Port.
func NewGetConsumerPositionService(state outbound.ConsumerStatePort) *GetConsumerPositionService {
	return &GetConsumerPositionService{state: state}
}

var _ inbound.GetConsumerPositionUseCase = (*GetConsumerPositionService)(nil)

// Position liest die bestätigte Position des Consumers (`LH-FA-CON-003`
// Happy Path); das Lesen verändert den gespeicherten Fortschritt nicht.
func (s *GetConsumerPositionService) Position(ctx context.Context, query GetConsumerPositionQuery) (GetConsumerPositionResult, error) {
	if query.Consumer == "" {
		return GetConsumerPositionResult{}, domainerrors.ErrEmptyIdentifier
	}
	carried, err := s.state.Position(ctx, query.Consumer)
	if err != nil {
		return GetConsumerPositionResult{}, err
	}
	return GetConsumerPositionResult{Position: carried}, nil
}
