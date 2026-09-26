// Die beiden Transformations-Use-Cases trägt diese Datei an einer Stelle
// (`ARC-003`): der Antragsweg der SQL-Administration (`ARC-005`,
// `LH-FA-CFG-007`) ruft sie über den Port auf (`ADR-0028`); die
// Orchestrierung liegt in den Application Services (`ARC-002`).

package inbound

import (
	"context"

	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// SetTransformationCommand trägt die Eingabe von `set_transformation`
// (`LH-FA-CFG-007`): die Tabelle, der Regelname und die Regelform als
// JSON-Text, wie der Antrags-Datensatz sie hält (`SPEC-019`,
// `SPEC-030`). Beide Textfelder können leer sein; der Use Case prüft sie und
// meldet die Verletzung mit dem Fehlertext der Spec.
type SetTransformationCommand struct {
	Source   model.SourceID
	Schema   string
	Table    string
	RuleName string
	RuleSpec string
}

// RemoveTransformationCommand trägt die Eingabe von `remove_transformation`
// (`LH-FA-CFG-007`): die Tabelle und der Regelname, den sie führt.
type RemoveTransformationCommand struct {
	Source   model.SourceID
	Schema   string
	Table    string
	RuleName string
}

// SetTransformationUseCase legt eine Transformationsregel einer Tabelle an
// (`LH-FA-CFG-007`, `ADR-0028`): die Form der Regel und die
// Konfliktfreiheit K1 bis K4 (`SPEC-019`) sind die Vorbedingungen. Der
// Aufruf liefert die geprüfte Regel, die der Aufrufer in die laufende
// Bindung einträgt; jede verletzte Vorbedingung endet als Fehler mit dem
// Fehlertext der Spec (Klartext, Doppelpunkt, Adresse), und der Regelstand
// bleibt unverändert. Der Use Case schreibt den Regelstand nicht: er ist die
// Ableitung aus den `applied`-Zeilen der Antrags-Queue (`ADR-0112`
// Teilfrage 6).
type SetTransformationUseCase interface {
	Set(ctx context.Context, command SetTransformationCommand) (model.Transformation, error)
}

// RemoveTransformationUseCase nimmt eine Transformationsregel einer Tabelle
// heraus (`LH-FA-CFG-007`, `ADR-0028`): der Regelname gehört zum Regelstand
// der Tabelle (K4, `SPEC-019`); sonst endet der Aufruf als Fehler mit dem
// Fehlertext der Spec.
type RemoveTransformationUseCase interface {
	Remove(ctx context.Context, command RemoveTransformationCommand) error
}
