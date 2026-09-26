package model

import (
	"fmt"
	"strings"

	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

// TransformationKind trägt den Regeltyp einer Transformationsregel
// (`LH-FA-CFG-007`, `SPEC-030`). Die geschlossene Menge der Regeltypen
// steht in `TransformationKinds`.
type TransformationKind string

// TransformationRenameColumn ist der Regeltyp `rename_column`: der Schlüssel
// `column` heißt im Row Image `to`, der Wert bleibt unverändert.
const TransformationRenameColumn TransformationKind = "rename_column"

// maxTransformationTargetBytes ist die Bezeichner-Länge der Quelle in
// Byte (UTF-8), die ein Zielname nicht überschreitet (`SPEC-030`,
// Bezeichner).
const maxTransformationTargetBytes = 63

// TransformationKinds liefert die Regeltypen als geschlossene Menge — die
// eine Quelle, aus der Aufrufer und Tests die Typen aufzählen. Jeder Aufruf
// liefert eine neue Liste.
func TransformationKinds() []TransformationKind {
	return []TransformationKind{TransformationRenameColumn}
}

// Transformation trägt eine Transformationsregel einer Tabelle: Name, Typ
// und die Parameter des Typs. Die Felder sind nicht exportiert, ein Wert
// entsteht über den Konstruktor seines Regeltyps und ist danach
// unveränderlich; der Nullwert trägt keine Spalte und wirkt auf kein Bild.
type Transformation struct {
	name   string
	kind   TransformationKind
	column string
	to     string
}

// NewRenameColumn legt eine Regel `rename_column` an (`SPEC-030`):
//   - `name` ist nicht leer (`ErrEmptyIdentifier`); das Alphabet des Namens
//     prüft der Aufrufer, der ihn entgegennimmt, nicht die Domäne;
//   - `column` ist nicht leer und trägt kein U+0000;
//   - `to` ist nicht leer, trägt kein U+0000 und hat höchstens 63 Byte in
//     UTF-8 (beide Verletzungen: `ErrInvalidTransformation`);
//   - `to` gleicht `column` nicht (`ErrTransformationTargetIsColumn`): der
//     Fall gehört zu K3 (`SPEC-019`), sein eigener Sentinel trennt ihn von
//     der Formverletzung.
func NewRenameColumn(name, column, to string) (Transformation, error) {
	if name == "" {
		return Transformation{}, fmt.Errorf("Regelname: %w", domainerrors.ErrEmptyIdentifier)
	}
	if column == "" || strings.IndexByte(column, 0) >= 0 {
		return Transformation{}, fmt.Errorf("column: %w", domainerrors.ErrInvalidTransformation)
	}
	if to == "" || strings.IndexByte(to, 0) >= 0 || len(to) > maxTransformationTargetBytes {
		return Transformation{}, fmt.Errorf("to: %w", domainerrors.ErrInvalidTransformation)
	}
	if to == column {
		return Transformation{}, domainerrors.ErrTransformationTargetIsColumn
	}
	return Transformation{name: name, kind: TransformationRenameColumn, column: column, to: to}, nil
}

// Name trägt den Regelnamen, mit dem die Regel je Tabelle adressiert wird.
func (t Transformation) Name() string { return t.name }

// Kind trägt den Regeltyp.
func (t Transformation) Kind() TransformationKind { return t.kind }

// Column trägt die Quellspalte der Regel.
func (t Transformation) Column() string { return t.column }

// To trägt den Zielnamen von `rename_column`.
func (t Transformation) To() string { return t.to }

// CheckApplicable prüft die Anwendbarkeit der Regel auf eine Änderung
// (`SPEC-030`, Anwendbarkeit): ihre Spalte kommt in `columns` vor, und der
// Zielname gleicht keiner Spalte aus `columns`, zeichengenau. `columns` sind
// alle Spalten der Relation der Änderung, ausgeschlossene eingeschlossen. Die
// Prüfung hängt an Regel und Spaltenmenge, nie am Wert einer Zeile.
// Verletzungen: `ErrTransformationColumnMissing`,
// `ErrTransformationTargetCollides`.
func (t Transformation) CheckApplicable(columns []string) error {
	if !containsName(columns, t.column) {
		return fmt.Errorf("%w: %s", domainerrors.ErrTransformationColumnMissing, t.column)
	}
	if containsName(columns, t.to) {
		return fmt.Errorf("%w: %s", domainerrors.ErrTransformationTargetCollides, t.to)
	}
	return nil
}

// applyTransformations wertet die Regeln für eine Spalte aus, deren Wert das
// Bild trägt: liefert den Schlüssel, unter dem die Spalte im Bild steht, und
// den Wert dort. Eine Regel trifft, wenn ihre Quellspalte `column` ist; die
// erste treffende Regel entscheidet (K2 hält je Quellspalte höchstens eine),
// keine treffende Regel lässt Schlüssel und Wert unverändert. Die Funktion
// ist rein.
func applyTransformations(rules []Transformation, column, value string) (string, string) {
	for _, rule := range rules {
		if rule.column != column {
			continue
		}
		switch rule.kind {
		case TransformationRenameColumn:
			return rule.to, value
		}
	}
	return column, value
}
