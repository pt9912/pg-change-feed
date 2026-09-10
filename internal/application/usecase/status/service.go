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
// Abwesenheit der Bindungs-Zeile (`LH-FA-CFG-003` Boundary); die
// Publication-Mitgliedschaft trennt den Erfassungs-Zustand von der
// Bindungs-Zeile als Herkunft.
type GetStatusService struct {
	activation outbound.TableActivationPort
}

// NewGetStatusService verdrahtet den Use Case mit dem Aktivierungs-Port.
func NewGetStatusService(activation outbound.TableActivationPort) *GetStatusService {
	return &GetStatusService{activation: activation}
}

var _ inbound.GetStatusUseCase = (*GetStatusService)(nil)

// Status meldet den CDC-Zustand der Tabelle: „aktiviert" liest die
// Bindungs-Zeile samt Publication-Mitgliedschaft; eine Bindungs-Zeile
// ohne Mitgliedschaft liest sich als Herkunft persistierter Changes
// (Retained — `LH-FA-CFG-002` Out-of-Scope: ihr Verhalten folgt der
// Retention), ihre Abwesenheit als „nicht aktiviert" (`LH-FA-CFG-003`
// Boundary).
func (s *GetStatusService) Status(ctx context.Context, query GetStatusQuery) (GetStatusResult, error) {
	if query.Source == "" || query.Schema == "" || query.Table == "" || query.Publication == "" {
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
	if !registered {
		return GetStatusResult{}, nil
	}
	published, err := s.activation.Published(ctx, query.Publication, query.Schema, query.Table)
	if err != nil {
		return GetStatusResult{}, err
	}
	if published {
		return GetStatusResult{Enabled: true}, nil
	}
	return GetStatusResult{Retained: true}, nil
}
