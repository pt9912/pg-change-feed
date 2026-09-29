// Die beiden Transformations-Use-Cases trägt diese Datei an einer Stelle
// (`ARC-003`): der Antragsweg der SQL-Administration ruft sie über den Port
// auf; die
// Orchestrierung liegt in den Application Services.

package inbound

import (
	"context"

	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// SetTransformationCommand trägt die Eingabe von `set_transformation`
// (`LH-FA-CFG-007`): die Tabelle, der Regelname und die Regelform als
// JSON-Text, wie der Antrags-Datensatz sie hält. Beide Textfelder können
// leer sein; der Use Case prüft sie und
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
// (`LH-FA-CFG-007`): die Form der Regel und die
// Konfliktfreiheit K1 bis K4 sind die Vorbedingungen. Der
// Aufruf liefert die geprüfte Regel, die der Aufrufer in die laufende
// Bindung einträgt; jede verletzte Vorbedingung endet als Fehler mit dem
// Fehlertext der Spec (Klartext, Doppelpunkt, Adresse), und der Regelstand
// bleibt unverändert. Der Use Case schreibt den Regelstand nicht: er ist die
// Ableitung aus den `applied`-Zeilen der Antrags-Queue.
type SetTransformationUseCase interface {
	Set(ctx context.Context, command SetTransformationCommand) (model.Transformation, error)
}

// RemoveTransformationUseCase nimmt eine Transformationsregel einer Tabelle
// heraus (`LH-FA-CFG-007`): der Regelname gehört zum Regelstand
// der Tabelle (K4); sonst endet der Aufruf als Fehler mit dem
// Fehlertext der Spec.
type RemoveTransformationUseCase interface {
	Remove(ctx context.Context, command RemoveTransformationCommand) error
}
