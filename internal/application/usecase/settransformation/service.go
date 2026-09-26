// Package settransformation trägt den SetTransformation Use Case (`ARC-002`,
// `ADR-0028`): das Anlegen einer Transformationsregel einer Tabelle
// (`LH-FA-CFG-007`).
package settransformation

import (
	"context"
	stderrors "errors"
	"fmt"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// SetTransformationCommand ist der Transport-Typ des SetTransformation Use
// Cases (`ADR-0042`): seine Definition liegt am Inbound-Port; die
// Use-Case-Adresse ist ein Alias desselben Typs.
type SetTransformationCommand = inbound.SetTransformationCommand

// SetTransformationService implementiert `inbound.SetTransformationUseCase`
// und prüft die Regel eines `set_transformation`-Antrags (`SPEC-019`): die
// Formzeilen in der Reihenfolge der Fehlertext-Tabelle, danach die
// Konfliktfreiheit K1 bis K4 gegen den Regelstand und die Spaltenliste der
// Quelltabelle, gelesen über den `TransformationPort`. Die Prüfungen selbst
// sind Domänen-Funktionen (`model.CheckRuleName`,
// `model.ParseTransformationSpec`, `TransformationSpec.CheckConflicts`); der
// Use Case bestimmt Reihenfolge und Adresse des Fehlertextes.
type SetTransformationService struct {
	rules outbound.TransformationPort
}

// NewSetTransformationService verdrahtet den Use Case mit dem
// Transformations-Port.
func NewSetTransformationService(rules outbound.TransformationPort) *SetTransformationService {
	return &SetTransformationService{rules: rules}
}

var _ inbound.SetTransformationUseCase = (*SetTransformationService)(nil)

// Set prüft den Antrag und liefert die geprüfte Regel. Eine verletzte
// Vorbedingung endet als Fehler, dessen Text der Fehlertext der Spec ist —
// Klartext, Doppelpunkt, Leerzeichen, Adresse (`SPEC-019`) — und der über
// `errors.Is` auf den Grund auflöst; der Regelstand bleibt unverändert, der
// Use Case schreibt ihn nicht. Die Prüfungen der Regelform laufen vor dem
// ersten Lesezugriff auf den Port.
//
// K3 prüft gegen die Spaltenliste zum Antragszeitpunkt: eine spätere
// Spalten-Erweiterung der Quelltabelle, die den Zielnamen kollidieren lässt,
// erreicht diese Prüfung nicht — sie endet als nicht anwendbare Regel im
// Erfassungspfad (`SPEC-030`, Anwendbarkeit).
func (s *SetTransformationService) Set(ctx context.Context, command SetTransformationCommand) (model.Transformation, error) {
	table := command.Schema + "." + command.Table
	if err := model.CheckRuleName(command.RuleName); err != nil {
		return model.Transformation{}, reject(err, table+"."+command.RuleName)
	}
	spec, err := model.ParseTransformationSpec(command.RuleSpec)
	if err != nil {
		var specErr *model.TransformationSpecError
		if stderrors.As(err, &specErr) {
			return model.Transformation{}, reject(specErr.Err, specErr.Detail)
		}
		return model.Transformation{}, reject(domainerrors.ErrInvalidRuleSpec, table+"."+command.RuleName)
	}

	state, err := s.rules.TransformationRules(ctx, command.Source)
	if err != nil {
		return model.Transformation{}, err
	}
	columns, err := s.rules.SourceColumns(ctx, command.Schema, command.Table)
	if err != nil {
		return model.Transformation{}, err
	}
	if err := spec.CheckConflicts(command.RuleName, state[table], columns); err != nil {
		return model.Transformation{}, reject(err, conflictAddress(err, table, command.RuleName, spec))
	}
	if !contains(columns, spec.Column()) {
		return model.Transformation{}, reject(inbound.ErrSourceColumnMissing, table+"."+spec.Column())
	}
	rule, err := spec.Build(command.RuleName)
	if err != nil {
		return model.Transformation{}, reject(domainerrors.ErrInvalidRuleSpec, table+"."+command.RuleName)
	}
	return rule, nil
}

// conflictAddress trägt die Adresse des Fehlertextes einer Konflikt-Zeile
// (`SPEC-019`): der Regelname bei K1, die Spalte bei K2, der Zielname bei K3.
func conflictAddress(reason error, table, ruleName string, spec model.TransformationSpec) string {
	switch {
	case stderrors.Is(reason, domainerrors.ErrRuleNameTaken):
		return table + "." + ruleName
	case stderrors.Is(reason, domainerrors.ErrColumnHasRule):
		return table + "." + spec.Column()
	default:
		return table + "." + spec.Target()
	}
}

// reject bildet den Fehlertext der Spec: der Klartext des Grundes, ein
// Doppelpunkt, ein Leerzeichen und die Adresse; der Grund bleibt über
// `errors.Is` erreichbar.
func reject(reason error, address string) error {
	return fmt.Errorf("%w: %s", reason, address)
}

// contains vergleicht zeichengenau (`SPEC-030`, Bezeichner).
func contains(names []string, name string) bool {
	for _, candidate := range names {
		if candidate == name {
			return true
		}
	}
	return false
}
