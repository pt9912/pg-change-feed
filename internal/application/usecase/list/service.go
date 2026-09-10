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
// die aktivierten Tabellen über den `TableActivationPort`.
type ListTablesService struct {
	activation outbound.TableActivationPort
}

// NewListTablesService verdrahtet den Use Case mit dem Aktivierungs-Port.
func NewListTablesService(activation outbound.TableActivationPort) *ListTablesService {
	return &ListTablesService{activation: activation}
}

var _ inbound.ListTablesUseCase = (*ListTablesService)(nil)

// ListTables liest die aktivierten Tabellen der Quelle; ohne Aktivierung
// trägt die Rückkehr eine leere Liste (`LH-FA-CFG-004` Boundary).
func (s *ListTablesService) ListTables(ctx context.Context, query ListTablesQuery) (ListTablesResult, error) {
	if query.Source == "" {
		return ListTablesResult{}, domainerrors.ErrEmptyIdentifier
	}
	tables, err := s.activation.List(ctx, query.Source)
	if err != nil {
		return ListTablesResult{}, err
	}
	if tables == nil {
		tables = []model.SourceTable{}
	}
	return ListTablesResult{Tables: tables}, nil
}
