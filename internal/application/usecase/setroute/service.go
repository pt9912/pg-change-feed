// Package setroute trägt den SetRoute Use Case (`ARC-002`): das Anlegen einer
// Routing-Regel einer Tabelle.
package setroute

import (
	"context"
	stderrors "errors"
	"fmt"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// SetRouteCommand ist der Transport-Typ des SetRoute Use Cases (`ADR-0042`):
// seine Definition liegt am Inbound-Port; die Use-Case-Adresse ist ein Alias
// desselben Typs.
type SetRouteCommand = inbound.SetRouteCommand

// SetRouteService implementiert `inbound.SetRouteUseCase` und prüft die Regel
// eines `set_route`-Antrags (`SPEC-019`): die Formzeilen in der Reihenfolge
// der Fehlertext-Tabelle, danach R1 und R2, R3 (die Spalte existiert und ist
// nicht ausgeschlossen) sowie R4 und R5 gegen den Routing-Regelstand
// (`RoutingPort`), die Spaltenexistenz und den Ausschlussstand
// (`ColumnExclusionPort`). Die Prüfungen selbst sind Domänen-Funktionen
// (`model.CheckRuleName`, `model.ParseRouteSpec`, `RouteRule.CheckIdentity`,
// `RouteRule.CheckNotExcluded`, `RouteRule.CheckArrangement`); der Use Case
// bestimmt Reihenfolge und Adresse des Fehlertextes.
type SetRouteService struct {
	routes  outbound.RoutingPort
	columns outbound.ColumnExclusionPort
}

// NewSetRouteService verdrahtet den Use Case mit dem Routing-Port und dem
// Spaltenprüfungs-Port.
func NewSetRouteService(routes outbound.RoutingPort, columns outbound.ColumnExclusionPort) *SetRouteService {
	return &SetRouteService{routes: routes, columns: columns}
}

var _ inbound.SetRouteUseCase = (*SetRouteService)(nil)

// Set prüft den Antrag und liefert die geprüfte Regel. Eine verletzte
// Vorbedingung endet als Fehler, dessen Text der Fehlertext der Spec ist —
// Klartext, Doppelpunkt, Leerzeichen, Adresse — und der über `errors.Is` auf
// den Grund auflöst; der Regelstand bleibt unverändert, der Use Case schreibt
// ihn nicht. Die Prüfungen der Regelform laufen vor dem ersten Lesezugriff auf
// einen Port; die Spaltenexistenz liest nur eine Regel mit `when`.
//
// R3 prüft gegen den Katalog und den Ausschlussstand zum Antragszeitpunkt:
// eine spätere Änderung an der Quelle erreicht diese Prüfung nicht — sie endet
// als nicht anwendbare Regel im Erfassungspfad (`SPEC-032`, Anwendbarkeit).
func (s *SetRouteService) Set(ctx context.Context, command SetRouteCommand) (model.RouteRule, error) {
	table := command.Schema + "." + command.Table
	if err := model.CheckRuleName(command.RuleName); err != nil {
		return model.RouteRule{}, reject(err, table+"."+command.RuleName)
	}
	spec, err := model.ParseRouteSpec(command.RuleSpec)
	if err != nil {
		var specErr *model.RouteSpecError
		switch {
		case !stderrors.As(err, &specErr):
			return model.RouteRule{}, reject(domainerrors.ErrInvalidRuleSpec, table+"."+command.RuleName)
		case stderrors.Is(specErr.Err, domainerrors.ErrUnknownRuleSpecKey):
			return model.RouteRule{}, reject(specErr.Err, specErr.Detail)
		default:
			return model.RouteRule{}, reject(specErr.Err, table+"."+specErr.Detail)
		}
	}
	rule, err := spec.Build(command.RuleName)
	if err != nil {
		return model.RouteRule{}, reject(domainerrors.ErrInvalidRuleSpec, table+"."+command.RuleName)
	}

	state, err := s.routes.RoutingRules(ctx, command.Source)
	if err != nil {
		return model.RouteRule{}, err
	}
	existing := state[table]
	if err := rule.CheckIdentity(existing); err != nil {
		return model.RouteRule{}, rejectConflict(err, table)
	}
	if condition, conditional := rule.Condition(); conditional {
		address := table + "." + condition.Column
		exists, err := s.columns.ColumnExists(ctx, command.Schema, command.Table, condition.Column)
		if err != nil {
			return model.RouteRule{}, err
		}
		if !exists {
			return model.RouteRule{}, reject(inbound.ErrSourceColumnMissing, address)
		}
		excluded, err := s.columns.ExcludedColumns(ctx, command.Source)
		if err != nil {
			return model.RouteRule{}, err
		}
		if err := rule.CheckNotExcluded(excluded[table]); err != nil {
			return model.RouteRule{}, reject(domainerrors.ErrRoutingColumnExcluded, address)
		}
	}
	if err := rule.CheckArrangement(existing); err != nil {
		return model.RouteRule{}, rejectConflict(err, table)
	}
	return rule, nil
}

// rejectConflict bildet den Fehlertext einer Konfliktzeile (R1, R2, R4, R5):
// der Klartext des Grundes und die Adresse `schema.table.<Wert>` mit dem Wert,
// den die Domänen-Funktion nennt.
func rejectConflict(err error, table string) error {
	var conflict *model.RouteSpecError
	if stderrors.As(err, &conflict) {
		return reject(conflict.Err, table+"."+conflict.Detail)
	}
	return err
}

// reject bildet den Fehlertext der Spec: der Klartext des Grundes, ein
// Doppelpunkt, ein Leerzeichen und die Adresse; der Grund bleibt über
// `errors.Is` erreichbar.
func reject(reason error, address string) error {
	return fmt.Errorf("%w: %s", reason, address)
}
