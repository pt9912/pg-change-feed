// Package includecolumn trägt den IncludeColumn Use Case (`ARC-002`,
// `ADR-0028`): die Aufhebung eines Spaltenausschlusses (`LH-FA-CFG-005`).
package includecolumn

import (
	"context"
	"fmt"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
)

// IncludeColumnCommand ist der Transport-Typ des IncludeColumn Use Cases
// (`ADR-0042`): seine Definition liegt am Inbound-Port — der Port trägt
// seinen Vertrag einschließlich der Transport-Typen, Ports referenzieren
// nur die Domain —; die Use-Case-Adresse ist ein Alias desselben Typs.
type IncludeColumnCommand = inbound.IncludeColumnCommand

// IncludeColumnService implementiert `inbound.IncludeColumnUseCase` und
// trägt den Spalteneinschluss über den `ColumnExclusionPort`: dieselbe
// Vorbedingung wie der Ausschluss — die Spalte muss an der Quelle
// existieren (`LH-FA-CFG-005` Negative).
type IncludeColumnService struct {
	columns outbound.ColumnExclusionPort
}

// NewIncludeColumnService verdrahtet den Use Case mit dem
// Spaltenprüfungs-Port.
func NewIncludeColumnService(columns outbound.ColumnExclusionPort) *IncludeColumnService {
	return &IncludeColumnService{columns: columns}
}

var _ inbound.IncludeColumnUseCase = (*IncludeColumnService)(nil)

// Include prüft die Spalte an der Quelle und endet bei ihrer Abwesenheit
// über `ErrSourceColumnMissing` (`LH-FA-CFG-005` Negative) — der
// Antrags-Status trägt den Fehlschlag danach als `failed` samt Fehlertext.
func (s *IncludeColumnService) Include(ctx context.Context, command IncludeColumnCommand) error {
	exists, err := s.columns.ColumnExists(ctx, command.Schema, command.Table, command.Column)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("%w: %s.%s.%s", inbound.ErrSourceColumnMissing, command.Schema, command.Table, command.Column)
	}
	return nil
}
