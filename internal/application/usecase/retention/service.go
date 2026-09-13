// Package retention trägt den RunRetention Use Case (`ARC-002`,
// `ADR-0014`): die Freigabe je Change trägt
// `model.RetentionPolicy.AllowsDeletion`, die physische Löschung der
// freigegebenen Menge trägt der `ChangeStorePort`.
package retention

import (
	"context"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// RunRetentionCommand und RunRetentionResult sind die Transport-Typen des
// RunRetention Use Cases (`ADR-0042`): ihre Definition liegt am
// Inbound-Port — der Port trägt seinen Vertrag einschließlich der
// Transport-Typen, Ports referenzieren nur die Domain —; die Use-Case-
// Adressen sind Aliase desselben Typs.
type (
	RunRetentionCommand = inbound.RunRetentionCommand
	RunRetentionResult  = inbound.RunRetentionResult
)

// RunRetentionService implementiert `inbound.RunRetentionUseCase` und
// orchestriert die Bereinigung über drei Ports: der `ConsumerStatePort`
// trägt die bestätigten Consumer-Positionen der Quelle, der `ClockPort`
// die Wanduhr für das Alter jedes Changes (`ADR-0040`), der
// `ChangeStorePort` liest die Kandidaten und führt die physische Löschung
// der freigegebenen Menge aus.
type RunRetentionService struct {
	store outbound.ChangeStorePort
	state outbound.ConsumerStatePort
	clock outbound.ClockPort
}

// NewRunRetentionService verdrahtet den Use Case mit seinen drei Ports.
func NewRunRetentionService(store outbound.ChangeStorePort, state outbound.ConsumerStatePort, clock outbound.ClockPort) *RunRetentionService {
	return &RunRetentionService{store: store, state: state, clock: clock}
}

var _ inbound.RunRetentionUseCase = (*RunRetentionService)(nil)

// Run liest die Kandidaten-Changes der Quelle und die bestätigten
// Consumer-Positionen, befragt `RetentionPolicy.AllowsDeletion` je
// betrachtetem Change real (Alter aus `CommittedAt` gegen die Wanduhr,
// Change-Position, alle bestätigten Consumer-Positionen der Quelle,
// `LH-FA-RET-002`…`004`) und übergibt ausschließlich die freigegebene
// Menge an `ChangeStorePort.DeleteChanges`. Eine leere Quellen-Kennung ist
// eine ungültige Konfiguration und endet über einen expliziten Fehlerpfad,
// keine stille Übernahme (`LH-FA-RET-002` Negative).
func (s *RunRetentionService) Run(ctx context.Context, command RunRetentionCommand) (RunRetentionResult, error) {
	if command.Source == "" {
		return RunRetentionResult{}, domainerrors.ErrEmptyIdentifier
	}

	consumerPositions, err := s.state.Positions(ctx, command.Source)
	if err != nil {
		return RunRetentionResult{}, err
	}

	records, err := s.store.ReadChanges(ctx, outbound.ChangeQuery{Source: command.Source})
	if err != nil {
		return RunRetentionResult{}, err
	}

	now := s.clock.Now()
	eligible := make([]model.ChangeID, 0, len(records))
	for _, record := range records {
		age := now.Sub(record.CommittedAt)
		if command.Policy.AllowsDeletion(age, record.Position, consumerPositions) {
			eligible = append(eligible, record.Change.ID)
		}
	}

	if err := s.store.DeleteChanges(ctx, eligible); err != nil {
		return RunRetentionResult{}, err
	}
	return RunRetentionResult{Deleted: len(eligible)}, nil
}
