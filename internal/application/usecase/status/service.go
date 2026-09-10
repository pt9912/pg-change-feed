// Package status trägt den GetStatus Use Case (`ARC-002`, `ADR-0028`):
// die Status-Abfrage meldet den CDC-Zustand einer Tabelle
// (`LH-FA-CFG-003`).
package status

import (
	"context"
	"fmt"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

// GetStatusQuery und GetStatusResult sind die Transport-Typen des
// GetStatus Use Cases (`ADR-0042`): ihre Definition liegt am Inbound-Port
// — der Port trägt seinen Vertrag einschließlich der Transport-Typen,
// Ports referenzieren nur die Domain —; die Use-Case-Adressen sind Aliase
// desselben Typs.
type (
	GetStatusQuery  = inbound.GetStatusQuery
	GetStatusResult = inbound.GetStatusResult
)

// GetStatusService implementiert `inbound.GetStatusUseCase` und liest den
// Zustand über den `TableActivationPort`: die Existenz der physischen
// Tabelle ist die Vorbedingung (`LH-FA-CFG-003` Negative: expliziter
// Fehlerpfad), der Zustand „nicht aktiviert" liest sich aus der
// Abwesenheit der Bindungs-Zeile (`LH-FA-CFG-003` Boundary).
type GetStatusService struct {
	activation outbound.TableActivationPort
}

// NewGetStatusService verdrahtet den Use Case mit dem Aktivierungs-Port.
func NewGetStatusService(activation outbound.TableActivationPort) *GetStatusService {
	return &GetStatusService{activation: activation}
}

var _ inbound.GetStatusUseCase = (*GetStatusService)(nil)

// Status meldet den CDC-Zustand der Tabelle: „aktiviert" liest die
// Bindungs-Zeile, „nicht aktiviert" ihre Abwesenheit — der Bindungs-Zeilen-
// Bestand trägt den Aktivierungszustand (`SPEC-001`, `cdc.source_table`).
func (s *GetStatusService) Status(ctx context.Context, query GetStatusQuery) (GetStatusResult, error) {
	if query.Source == "" || query.Schema == "" || query.Table == "" {
		return GetStatusResult{}, domainerrors.ErrEmptyIdentifier
	}
	exists, err := s.activation.TableExists(ctx, query.Schema, query.Table)
	if err != nil {
		return GetStatusResult{}, err
	}
	if !exists {
		return GetStatusResult{}, fmt.Errorf("%w: %s.%s", inbound.ErrSourceTableMissing, query.Schema, query.Table)
	}
	_, registered, err := s.activation.Registered(ctx, query.Source, query.Schema, query.Table)
	if err != nil {
		return GetStatusResult{}, err
	}
	return GetStatusResult{Enabled: registered}, nil
}
