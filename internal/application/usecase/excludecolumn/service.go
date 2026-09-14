// Package excludecolumn trägt den ExcludeColumn Use Case (`ARC-002`,
// `ADR-0028`): der Ausschluss einer Spalte von der Erfassung
// (`LH-FA-CFG-005`).
package excludecolumn

import (
	"context"
	"fmt"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
)

// ExcludeColumnCommand ist der Transport-Typ des ExcludeColumn Use Cases
// (`ADR-0042`): seine Definition liegt am Inbound-Port — der Port trägt
// seinen Vertrag einschließlich der Transport-Typen, Ports referenzieren
// nur die Domain —; die Use-Case-Adresse ist ein Alias desselben Typs.
type ExcludeColumnCommand = inbound.ExcludeColumnCommand

// ExcludeColumnService implementiert `inbound.ExcludeColumnUseCase` und
// trägt den Spaltenausschluss über den `ColumnExclusionPort`: die Existenz
// der physischen Spalte an der Quelle ist die Vorbedingung
// (`LH-FA-CFG-005` Negative: expliziter Fehlerpfad).
type ExcludeColumnService struct {
	columns outbound.ColumnExclusionPort
}

// NewExcludeColumnService verdrahtet den Use Case mit dem
// Spaltenprüfungs-Port.
func NewExcludeColumnService(columns outbound.ColumnExclusionPort) *ExcludeColumnService {
	return &ExcludeColumnService{columns: columns}
}

var _ inbound.ExcludeColumnUseCase = (*ExcludeColumnService)(nil)

// Exclude prüft die Spalte an der Quelle und endet bei ihrer Abwesenheit
// über `ErrSourceColumnMissing` (`LH-FA-CFG-005` Negative) — der
// Antrags-Status trägt den Fehlschlag danach als `failed` samt Fehlertext.
func (s *ExcludeColumnService) Exclude(ctx context.Context, command ExcludeColumnCommand) error {
	exists, err := s.columns.ColumnExists(ctx, command.Schema, command.Table, command.Column)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("%w: %s.%s.%s", inbound.ErrSourceColumnMissing, command.Schema, command.Table, command.Column)
	}
	return nil
}
