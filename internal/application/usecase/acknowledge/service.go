// Package acknowledge trägt den AcknowledgeConsumer Use Case (`ARC-002`,
// `ADR-0028`): die Bestätigung der verarbeiteten Position eines Consumers
// (`LH-FA-CON-004`) — der reguläre ACK verläuft nur vorwärts
// (`ADR-0029`, Regel 2).
package acknowledge

import (
	"context"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// AcknowledgeConsumerCommand und AcknowledgeConsumerResult sind die
// Transport-Typen des AcknowledgeConsumer Use Cases (`ADR-0042`): ihre
// Definition liegt am Inbound-Port — der Port trägt seinen Vertrag
// einschließlich der Transport-Typen, Ports referenzieren nur die Domain —;
// die Use-Case-Adressen sind Aliase desselben Typs.
type (
	AcknowledgeConsumerCommand = inbound.AcknowledgeConsumerCommand
	AcknowledgeConsumerResult  = inbound.AcknowledgeConsumerResult
)

// AcknowledgeConsumerService implementiert
// `inbound.AcknowledgeConsumerUseCase` und trägt die Bestätigung über den
// `ConsumerStatePort`: die Position läuft über den Domänen-Konstruktor —
// ein Offset ohne Wert endet über die Domänen-Invariante, ohne den
// gespeicherten Fortschritt zu berühren; die Monotonie trägt der
// fortgeführte Vergleich am Port.
type AcknowledgeConsumerService struct {
	state outbound.ConsumerStatePort
}

// NewAcknowledgeConsumerService verdrahtet den Use Case mit dem
// Consumer-State-Port.
func NewAcknowledgeConsumerService(state outbound.ConsumerStatePort) *AcknowledgeConsumerService {
	return &AcknowledgeConsumerService{state: state}
}

var _ inbound.AcknowledgeConsumerUseCase = (*AcknowledgeConsumerService)(nil)

// Acknowledge trägt die bestätigte Position fort (`LH-FA-CON-004`): die
// Wiederholung derselben Position ist idempotent (Boundary), eine frühere
// Position und eine Position einer anderen Quelle enden über die
// Invarianten-Sentinels (`ADR-0029`, Regel 2) — kein administrativer Reset
// läuft als ACK durch (`ADR-0013`).
func (s *AcknowledgeConsumerService) Acknowledge(ctx context.Context, command AcknowledgeConsumerCommand) (AcknowledgeConsumerResult, error) {
	if command.Consumer == "" {
		return AcknowledgeConsumerResult{}, domainerrors.ErrEmptyIdentifier
	}
	position, err := model.NewSourcePosition(command.Position.SourceID, command.Position.Offset)
	if err != nil {
		return AcknowledgeConsumerResult{}, err
	}
	carried, err := s.state.Acknowledge(ctx, model.ConsumerPosition{
		ConsumerID: command.Consumer,
		Position:   position,
	})
	if err != nil {
		return AcknowledgeConsumerResult{}, err
	}
	return AcknowledgeConsumerResult{Position: carried}, nil
}
