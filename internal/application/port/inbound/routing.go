// Die beiden Routing-Use-Cases trägt diese Datei an einer Stelle (`ARC-003`):
// der Antragsweg der SQL-Administration ruft sie über den Port auf; die
// Orchestrierung liegt in den Application Services.

package inbound

import (
	"context"

	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// SetRouteCommand trägt die Eingabe von `set_route` (`LH-FA-CFG-008`): die
// Tabelle, der Regelname und die Regelform als JSON-Text, wie der
// Antrags-Datensatz sie hält. Beide Textfelder können leer sein; der Use Case
// prüft sie und meldet die Verletzung mit dem Fehlertext der Spec.
type SetRouteCommand struct {
	Source   model.SourceID
	Schema   string
	Table    string
	RuleName string
	RuleSpec string
}

// RemoveRouteCommand trägt die Eingabe von `remove_route`
// (`LH-FA-CFG-008`): die Tabelle und der Regelname, den sie führt.
type RemoveRouteCommand struct {
	Source   model.SourceID
	Schema   string
	Table    string
	RuleName string
}

// SetRouteUseCase legt eine Routing-Regel einer Tabelle an
// (`LH-FA-CFG-008`): die Form der Regel und die Konfliktfreiheit R1 bis R5
// sind die Vorbedingungen. Der Aufruf liefert die geprüfte Regel, die der
// Aufrufer in die laufende Bindung einträgt; jede verletzte Vorbedingung
// endet als Fehler mit dem Fehlertext der Spec (Klartext, Doppelpunkt,
// Adresse), und der Regelstand bleibt unverändert. Der Use Case schreibt den
// Regelstand nicht: er ist die Ableitung aus den `applied`-Zeilen der
// Antrags-Queue.
type SetRouteUseCase interface {
	Set(ctx context.Context, command SetRouteCommand) (model.RouteRule, error)
}

// RemoveRouteUseCase nimmt eine Routing-Regel einer Tabelle heraus
// (`LH-FA-CFG-008`): der Regelname gehört zum Routing-Regelstand der Tabelle
// (R6); sonst endet der Aufruf als Fehler mit dem Fehlertext der Spec.
type RemoveRouteUseCase interface {
	Remove(ctx context.Context, command RemoveRouteCommand) error
}
