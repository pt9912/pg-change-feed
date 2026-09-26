// Package removetransformation trägt den RemoveTransformation Use Case
// (`ARC-002`, `ADR-0028`): das Herausnehmen einer Transformationsregel einer
// Tabelle (`LH-FA-CFG-007`).
package removetransformation

import (
	"context"
	"fmt"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// RemoveTransformationCommand ist der Transport-Typ des RemoveTransformation
// Use Cases (`ADR-0042`): seine Definition liegt am Inbound-Port; die
// Use-Case-Adresse ist ein Alias desselben Typs.
type RemoveTransformationCommand = inbound.RemoveTransformationCommand

// RemoveTransformationService implementiert
// `inbound.RemoveTransformationUseCase`: ein `remove_transformation`-Antrag
// durchläuft die Zeile zum Regelnamen und K4 (`SPEC-019`) — der Name gehört
// zum Regelstand der Tabelle, gelesen über den `TransformationPort`.
type RemoveTransformationService struct {
	rules outbound.TransformationPort
}

// NewRemoveTransformationService verdrahtet den Use Case mit dem
// Transformations-Port.
func NewRemoveTransformationService(rules outbound.TransformationPort) *RemoveTransformationService {
	return &RemoveTransformationService{rules: rules}
}

var _ inbound.RemoveTransformationUseCase = (*RemoveTransformationService)(nil)

// Remove prüft den Regelnamen gegen das Alphabet (`ErrInvalidRuleName`) und
// gegen den Regelstand der Tabelle (`ErrRuleNotKept`); der Fehlertext ist der
// der Spec — Klartext, Doppelpunkt, Leerzeichen, Adresse
// `schema.table.rule_name` (`SPEC-019`) — und löst über `errors.Is` auf den
// Grund auf. Der Use Case schreibt den Regelstand nicht.
func (s *RemoveTransformationService) Remove(ctx context.Context, command RemoveTransformationCommand) error {
	address := command.Schema + "." + command.Table + "." + command.RuleName
	if err := model.CheckRuleName(command.RuleName); err != nil {
		return fmt.Errorf("%w: %s", err, address)
	}
	state, err := s.rules.TransformationRules(ctx, command.Source)
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
