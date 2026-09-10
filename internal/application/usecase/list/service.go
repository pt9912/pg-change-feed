// Package list trägt den ListTables Use Case (`ARC-002`, `ADR-0028`): die
// Liste trägt die aktivierten Tabellen einer Quelle (`LH-FA-CFG-004`).
package list

import (
	"context"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// ListTablesQuery und ListTablesResult sind die Transport-Typen des
// ListTables Use Cases (`ADR-0042`): ihre Definition liegt am
// Inbound-Port — der Port trägt seinen Vertrag einschließlich der
// Transport-Typen, Ports referenzieren nur die Domain —; die Use-Case-
// Adressen sind Aliase desselben Typs.
type (
	ListTablesQuery  = inbound.ListTablesQuery
	ListTablesResult = inbound.ListTablesResult
)

// ListTablesService implementiert `inbound.ListTablesUseCase` und liest
// die aktivierten Tabellen über den `TableActivationPort`; die
// Publication-Mitgliedschaft trennt die aktivierten Tabellen von den
// Bindungs-Zeilen, die nur noch Herkunft tragen.
type ListTablesService struct {
	activation outbound.TableActivationPort
}

// NewListTablesService verdrahtet den Use Case mit dem Aktivierungs-Port.
func NewListTablesService(activation outbound.TableActivationPort) *ListTablesService {
	return &ListTablesService{activation: activation}
}

var _ inbound.ListTablesUseCase = (*ListTablesService)(nil)

// ListTables liest die Tabellen der Quelle getrennt: `Tables` liest die
// aktivierten Tabellen (Bindungs-Zeile samt Publication-Mitgliedschaft),
// `Retained` liest Bindungs-Zeilen, die nur noch Herkunft persistierter
// Changes tragen (`LH-FA-CFG-002` Out-of-Scope); ohne Aktivierung tragen
// beide Rückgaben leere Listen (`LH-FA-CFG-004` Boundary).
func (s *ListTablesService) ListTables(ctx context.Context, query ListTablesQuery) (ListTablesResult, error) {
	if query.Source == "" || query.Publication == "" {
		return ListTablesResult{}, domainerrors.ErrEmptyIdentifier
	}
	bindings, err := s.activation.List(ctx, query.Source)
	if err != nil {
		return ListTablesResult{}, err
	}
	result := ListTablesResult{Tables: []model.SourceTable{}, Retained: []model.SourceTable{}}
	for _, binding := range bindings {
		published, err := s.activation.Published(ctx, query.Publication, binding.Schema, binding.Table)
		if err != nil {
			return ListTablesResult{}, err
		}
		if published {
			result.Tables = append(result.Tables, binding)
		} else {
			result.Retained = append(result.Retained, binding)
		}
	}
	return result, nil
}
