package outbound

import (
	"context"

	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// TransformationPort trägt die Lese-Fähigkeiten der Transformationsregeln
// an der Quelle (`ARC-004`, `ADR-0112` Teilfrage 6, Fähigkeits-Port je
// `ADR-0034`): den dauerhaften Regelstand, den jeder Pfad beim Anlegen einer
// Bindung mitführt und die Konfliktprüfung K1 bis K3 liest, und die
// Spaltennamen der Quelltabelle, gegen die K3 und K4 prüfen. Beide
// Fähigkeiten gehören derselben Frage — welche Regeln die Spalten einer
// Tabelle tragen dürfen — und derselben Adapter-Instanz wie die Aktivierung
// (im MVP eine Instanz). Der Regelstand selbst wird nicht hier geschrieben:
// der Schreibpfad bleibt die Antrags-Queue (`ADR-0050`).
type TransformationPort interface {
	// TransformationRules liefert den dauerhaften Regelstand je Tabelle
	// einer Quelle (`SPEC-019`): die `applied`-Zeilen der beiden
	// Transformations-Antragsarten in `cdc.administration_request`,
	// ausgewertet in der Reihenfolge ihres `requested_at` — bei gleichem
	// Zeitstempel deterministisch nach der Antrags-ID. Der Schlüssel der
	// Rückgabe ist der qualifizierte Tabellenname in der Form
	// `schema.table`; eine Tabelle ohne geführte Regel trägt keinen Eintrag.
	// Jede Lesung liefert eine frisch aufgebaute Liste je Tabelle, die der
	// Aufrufer behalten darf. Die Ableitung ist die einzige Herkunft des
	// Standes: sie deckt den Prozessstart und einen Bindungs-Zyklus mit
	// demselben Mechanismus.
	TransformationRules(ctx context.Context, source model.SourceID) (map[string][]model.Transformation, error)

	// SourceColumns liefert die Spaltennamen der Quelltabelle in der
	// Reihenfolge der Tabelle, ausgeschlossene Spalten eingeschlossen
	// (Katalog-Lesart wie `ColumnExclusionPort.ColumnExists`); eine nicht
	// vorhandene Tabelle liefert die leere Liste.
	SourceColumns(ctx context.Context, schema, table string) ([]string, error)
}
