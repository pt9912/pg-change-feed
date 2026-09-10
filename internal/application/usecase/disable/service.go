// Package disable trägt den DisableTable Use Case (`ARC-002`,
// `ADR-0028`): die Deaktivierung einer Tabelle entzieht sie der
// Publication und trägt den Ausgang der Bindungs-Zeile
// (`LH-FA-CFG-002`).
package disable

import (
	"context"
	"fmt"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

// DisableTableCommand und DisableTableResult sind die Transport-Typen des
// DisableTable Use Cases (`ADR-0042`): ihre Definition liegt am
// Inbound-Port — der Port trägt seinen Vertrag einschließlich der
// Transport-Typen, Ports referenzieren nur die Domain —; die Use-Case-
// Adressen sind Aliase desselben Typs.
type (
	DisableTableCommand = inbound.DisableTableCommand
	DisableTableResult  = inbound.DisableTableResult
)

// DisableTableService implementiert `inbound.DisableTableUseCase` und
// trägt die Deaktivierung über den `TableActivationPort`: der
// Publication-Entzug trägt den Stopp der Erfassung (`LH-FA-CFG-002`
// Happy Path); die Bindungs-Zeile folgt dem Change-Bestand — bei
// persistierten Changes bleibt sie als Herkunft bestehen (`LH-FA-CFG-002`
// Out-of-Scope: ihr Verhalten folgt der Retention), ohne Bestand wird sie
// entfernt. Die `REPLICA IDENTITY` bleibt von der Deaktivierung
// unberührt.
type DisableTableService struct {
	activation outbound.TableActivationPort
}

// NewDisableTableService verdrahtet den Use Case mit dem Aktivierungs-Port.
func NewDisableTableService(activation outbound.TableActivationPort) *DisableTableService {
	return &DisableTableService{activation: activation}
}

var _ inbound.DisableTableUseCase = (*DisableTableService)(nil)

// Disable deaktiviert die Tabelle idempotent (`LH-FA-CFG-002` Boundary):
// der Publication-Entzug läuft gegen jede Deaktivierung, der
// Bindungs-Zeilen-Entzug nur gegen einen vorhandenen Stand — die Rückkehr
// meldet den Zeilen-Ausgang über `Removed` und `Retained`.
func (s *DisableTableService) Disable(ctx context.Context, command DisableTableCommand) (DisableTableResult, error) {
	if command.Source == "" || command.Schema == "" || command.Table == "" || command.Publication == "" {
		return DisableTableResult{}, domainerrors.ErrEmptyIdentifier
	}
	exists, err := s.activation.TableExists(ctx, command.Schema, command.Table)
	if err != nil {
		return DisableTableResult{}, err
	}
	if !exists {
		return DisableTableResult{}, fmt.Errorf("%w: %s.%s", inbound.ErrSourceTableMissing, command.Schema, command.Table)
	}
	if err := s.activation.Unpublish(ctx, command.Publication, command.Schema, command.Table); err != nil {
		return DisableTableResult{}, err
	}
	table, registered, err := s.activation.Registered(ctx, command.Source, command.Schema, command.Table)
	if err != nil {
		return DisableTableResult{}, err
	}
	if !registered {
		return DisableTableResult{}, nil
	}
	removal, err := s.activation.Unregister(ctx, table)
	if err != nil {
		return DisableTableResult{}, err
	}
	switch removal {
	case outbound.ActivationRemoved:
		return DisableTableResult{Removed: true}, nil
	case outbound.ActivationRetained:
		return DisableTableResult{Retained: true}, nil
	default:
		return DisableTableResult{}, nil
	}
}
