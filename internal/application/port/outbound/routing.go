package outbound

import (
	"context"

	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// RoutingPort trägt die Lese-Fähigkeit der Routing-Regeln an der Quelle
// (`ARC-004`, Fähigkeits-Port): den dauerhaften Regelstand, den jeder Pfad
// beim Anlegen einer Bindung mitführt und den die Konfliktprüfung R1 bis R6
// sowie die Gegenrichtung von R3 beim Spaltenausschluss liest. Dieselbe
// Adapter-Instanz wie die Aktivierung (im MVP eine Instanz). Der Regelstand
// selbst wird nicht hier geschrieben: der Schreibpfad bleibt die
// Antrags-Queue.
type RoutingPort interface {
	// RoutingRules liefert den dauerhaften Routing-Regelstand je Tabelle
	// einer Quelle (`SPEC-019`): die `applied`-Zeilen der beiden
	// Routing-Antragsarten in `cdc.administration_request`, ausgewertet in
	// der Reihenfolge ihres `requested_at` — bei gleichem Zeitstempel
	// deterministisch nach der Antrags-ID. Der Schlüssel der Rückgabe ist der
	// qualifizierte Tabellenname in der Form `schema.table`; eine Tabelle
	// ohne geführte Regel trägt keinen Eintrag. Jede Lesung liefert eine
	// frisch aufgebaute Liste je Tabelle, die der Aufrufer behalten darf. Die
	// Ableitung ist die einzige Herkunft des Standes: sie deckt den
	// Prozessstart und einen Bindungs-Zyklus mit demselben Mechanismus.
	RoutingRules(ctx context.Context, source model.SourceID) (map[string][]model.RouteRule, error)
}
