// Package excludecolumn trägt den ExcludeColumn Use Case (`ARC-002`):
// der Ausschluss einer Spalte von der Erfassung.
package excludecolumn

import (
	"context"
	"fmt"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// ExcludeColumnCommand ist der Transport-Typ des ExcludeColumn Use Cases
// (`ADR-0042`): seine Definition liegt am Inbound-Port — der Port trägt
// seinen Vertrag einschließlich der Transport-Typen, Ports referenzieren
// nur die Domain —; die Use-Case-Adresse ist ein Alias desselben Typs.
type ExcludeColumnCommand = inbound.ExcludeColumnCommand

// ExcludeColumnService implementiert `inbound.ExcludeColumnUseCase` und
// trägt den Spaltenausschluss über den `ColumnExclusionPort`: die Existenz
// der physischen Spalte an der Quelle ist die Vorbedingung
// (`LH-FA-CFG-005` Negative: expliziter Fehlerpfad). Die zweite
// Vorbedingung ist die Gegenrichtung von R3: eine Spalte, die eine
// Routing-Bedingung des Regelstands trägt, wird nicht ausgeschlossen; den
// Regelstand liest der `RoutingPort`.
type ExcludeColumnService struct {
	columns outbound.ColumnExclusionPort
	routes  outbound.RoutingPort
}

// NewExcludeColumnService verdrahtet den Use Case mit dem
// Spaltenprüfungs-Port und dem Routing-Port.
func NewExcludeColumnService(columns outbound.ColumnExclusionPort, routes outbound.RoutingPort) *ExcludeColumnService {
	return &ExcludeColumnService{columns: columns, routes: routes}
}

var _ inbound.ExcludeColumnUseCase = (*ExcludeColumnService)(nil)

// Exclude prüft die Spalte an der Quelle und endet bei ihrer Abwesenheit
// über `ErrSourceColumnMissing` (`LH-FA-CFG-005` Negative) — der
// Antrags-Status trägt den Fehlschlag danach als `failed` samt Fehlertext.
// Eine vorhandene Spalte mit Routing-Bedingung endet über
// `ErrColumnHasRouteCondition` (R3); der Regelstand wird erst gelesen, wenn
// die Spalte existiert.
func (s *ExcludeColumnService) Exclude(ctx context.Context, command ExcludeColumnCommand) error {
	exists, err := s.columns.ColumnExists(ctx, command.Schema, command.Table, command.Column)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("%w: %s.%s.%s", inbound.ErrSourceColumnMissing, command.Schema, command.Table, command.Column)
	}
	state, err := s.routes.RoutingRules(ctx, command.Source)
	if err != nil {
		return err
	}
	if model.RoutesConditionOn(state[command.Schema+"."+command.Table], command.Column) {
		return fmt.Errorf("%w: %s.%s.%s", domainerrors.ErrColumnHasRouteCondition, command.Schema, command.Table, command.Column)
	}
	return nil
}
