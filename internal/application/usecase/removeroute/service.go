// Package removeroute trägt den RemoveRoute Use Case (`ARC-002`): das
// Herausnehmen einer Routing-Regel einer Tabelle.
package removeroute

import (
	"context"
	"fmt"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// RemoveRouteCommand ist der Transport-Typ des RemoveRoute Use Cases
// (`ADR-0042`): seine Definition liegt am Inbound-Port; die Use-Case-Adresse
// ist ein Alias desselben Typs.
type RemoveRouteCommand = inbound.RemoveRouteCommand

// RemoveRouteService implementiert `inbound.RemoveRouteUseCase`: ein
// `remove_route`-Antrag durchläuft die Zeile zum Regelnamen und R6
// (`SPEC-019`) — der Name gehört zum Routing-Regelstand der Tabelle, gelesen
// über den `RoutingPort`.
type RemoveRouteService struct {
	routes outbound.RoutingPort
}

// NewRemoveRouteService verdrahtet den Use Case mit dem Routing-Port.
func NewRemoveRouteService(routes outbound.RoutingPort) *RemoveRouteService {
	return &RemoveRouteService{routes: routes}
}

var _ inbound.RemoveRouteUseCase = (*RemoveRouteService)(nil)

// Remove prüft den Regelnamen gegen das Alphabet (`ErrInvalidRuleName`) und
// gegen den Routing-Regelstand der Tabelle (`ErrRuleNotKept`); der Fehlertext
// ist der der Spec — Klartext, Doppelpunkt, Leerzeichen, Adresse
// `schema.table.rule_name` (`SPEC-019`) — und löst über `errors.Is` auf den
// Grund auf. Der Use Case schreibt den Regelstand nicht.
func (s *RemoveRouteService) Remove(ctx context.Context, command RemoveRouteCommand) error {
	address := command.Schema + "." + command.Table + "." + command.RuleName
	if err := model.CheckRuleName(command.RuleName); err != nil {
		return fmt.Errorf("%w: %s", err, address)
	}
	state, err := s.routes.RoutingRules(ctx, command.Source)
	if err != nil {
		return err
	}
	for _, rule := range state[command.Schema+"."+command.Table] {
		if rule.Name() == command.RuleName {
			return nil
		}
	}
	return fmt.Errorf("%w: %s", domainerrors.ErrRuleNotKept, address)
}
