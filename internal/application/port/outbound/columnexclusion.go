package outbound

import (
	"context"

	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// ColumnExclusionPort trägt die Fähigkeiten des Spaltenausschlusses an
// der Quelle (`ARC-004`, Fähigkeits-Port je `ADR-0034`): der
// Spaltenausschluss/-einschluss (`LH-FA-CFG-005`, `ADR-0059`) prüft
// darüber die Vorbedingung seines Negative-Pfads — dieselbe
// Katalog-Lesart wie `TableActivationPort.TableExists`, hier auf eine
// Spalte einer Tabelle gerichtet — und liest den dauerhaften
// Ausschlussstand, den jeder Pfad beim Anlegen einer Bindung mitführt
// (`ADR-0065`).
//
// Beide Fähigkeiten gehören derselben Objektklasse: der Spalte einer
// Tabelle an der Quelle. Der Ausschluss-Zustand selbst wird nicht hier
// geschrieben — Schreibpfad bleibt die Antrags-Queue (`ADR-0050`); diese
// Fähigkeit liest ihn nur zurück.
type ColumnExclusionPort interface {
	// ColumnExists prüft die physische Spalte an der Quelle; der
	// Negative-Pfad des Spaltenausschlusses (`LH-FA-CFG-005`) endet über
	// die Abwesenheit sichtbar statt still.
	ColumnExists(ctx context.Context, schema, table, column string) (bool, error)

	// ExcludedColumns liefert den dauerhaften Ausschlussstand je Tabelle
	// einer Quelle (`ADR-0065`): die `applied`-Zeilen der beiden
	// Spalten-Antragsarten in `cdc.administration_request`, ausgewertet in
	// der Reihenfolge ihres `requested_at` — bei gleichem Zeitstempel
	// deterministisch nach der Antrags-ID. Der Schlüssel der Rückgabe ist
	// der qualifizierte Tabellenname in der Form `schema.table`; eine
	// Tabelle ohne geführten Ausschluss trägt keinen Eintrag. Die
	// Ableitung ist die einzige Herkunft des Standes: sie deckt den
	// Prozessstart und einen Bindungs-Zyklus mit demselben Mechanismus.
	ExcludedColumns(ctx context.Context, source model.SourceID) (map[string][]string, error)
}
