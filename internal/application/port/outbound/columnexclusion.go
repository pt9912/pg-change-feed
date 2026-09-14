package outbound

import "context"

// ColumnExclusionPort trägt die Fähigkeit der Spaltenexistenz-Prüfung an
// der Quelle (`ARC-004`, Fähigkeits-Port je `ADR-0034`): der
// Spaltenausschluss/-einschluss (`LH-FA-CFG-005`, `ADR-0059`) prüft
// darüber die Vorbedingung seines Negative-Pfads — dieselbe
// Katalog-Lesart wie `TableActivationPort.TableExists`, hier auf eine
// Spalte einer Tabelle gerichtet.
//
// Der Zuschnitt bleibt bei der einen Prüfung: der Ausschluss-Zustand
// selbst gehört keiner Zeile dieser Quelle, sondern der laufenden
// Erfassung (`ADR-0059` Teilfrage 3).
type ColumnExclusionPort interface {
	// ColumnExists prüft die physische Spalte an der Quelle; der
	// Negative-Pfad des Spaltenausschlusses (`LH-FA-CFG-005`) endet über
	// die Abwesenheit sichtbar statt still.
	ColumnExists(ctx context.Context, schema, table, column string) (bool, error)
}
