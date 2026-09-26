package model

import (
	"fmt"
	"sort"
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

// TransformationMapValue ist der Regeltyp `map_value`: ist der Wert von
// `column` als Zeichenkette ein Schlüssel der Zuordnung, steht der
// zugeordnete Wert im Bild; jeder andere Wert bleibt unverändert.
const TransformationMapValue TransformationKind = "map_value"

// maxTransformationTargetBytes ist die Bezeichner-Länge der Quelle in
// Byte (UTF-8), die ein Zielname nicht überschreitet (`SPEC-030`,
// Bezeichner).
const maxTransformationTargetBytes = 63

// TransformationKinds liefert die Regeltypen als geschlossene Menge — die
// eine Quelle, aus der Aufrufer und Tests die Typen aufzählen. Jeder Aufruf
// liefert eine neue Liste.
func TransformationKinds() []TransformationKind {
	return []TransformationKind{TransformationRenameColumn, TransformationMapValue}
}

// Transformation trägt eine Transformationsregel einer Tabelle: Name, Typ
// und die Parameter des Typs. Die Felder sind nicht exportiert, ein Wert
// entsteht über den Konstruktor seines Regeltyps und ist danach
// unveränderlich; der Nullwert trägt keine Spalte und wirkt auf kein Bild.
//
// `values` trägt die Zuordnung von `map_value` in kanonischer Kodierung
// (`encodeValueMap`) und ist bei jedem anderen Regeltyp leer. Alle Felder
// sind Zeichenketten, damit `Transformation` über `==` vergleichbar bleibt:
// `sameSet` in `internal/application/usecase/backfill/service.go` vergleicht
// Regelstände als Menge über diesen Typ.
type Transformation struct {
	name   string
	kind   TransformationKind
	column string
	to     string
	values string
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

// NewMapValue legt eine Regel `map_value` an (`SPEC-030`):
//   - `name` ist nicht leer (`ErrEmptyIdentifier`);
//   - `column` ist nicht leer und trägt kein U+0000;
//   - `values` ist nicht leer (beide Verletzungen:
//     `ErrInvalidTransformation`).
//
// Die Abbildung eines Werts auf sich selbst und mehrere Quellwerte auf
// denselben Zielwert sind zulässig. Die Regel kopiert die Zuordnung: eine
// spätere Änderung von `values` wirkt nicht auf sie.
func NewMapValue(name, column string, values map[string]string) (Transformation, error) {
	return newMapValue(name, column, encodeValueMap(values))
}

// newMapValue legt die Regel aus der kanonischen Kodierung der Zuordnung an;
// die leere Kodierung ist die leere Zuordnung.
func newMapValue(name, column, encoded string) (Transformation, error) {
	if name == "" {
		return Transformation{}, fmt.Errorf("Regelname: %w", domainerrors.ErrEmptyIdentifier)
	}
	if column == "" || strings.IndexByte(column, 0) >= 0 {
		return Transformation{}, fmt.Errorf("column: %w", domainerrors.ErrInvalidTransformation)
	}
	if encoded == "" {
		return Transformation{}, fmt.Errorf("values: %w", domainerrors.ErrInvalidTransformation)
	}
	return Transformation{name: name, kind: TransformationMapValue, column: column, values: encoded}, nil
}

// Name trägt den Regelnamen, mit dem die Regel je Tabelle adressiert wird.
func (t Transformation) Name() string { return t.name }

// Kind trägt den Regeltyp.
func (t Transformation) Kind() TransformationKind { return t.kind }

// Column trägt die Quellspalte der Regel.
func (t Transformation) Column() string { return t.column }

// To trägt den Zielnamen von `rename_column`; ein Regeltyp ohne Zielname
// trägt den leeren Wert.
func (t Transformation) To() string { return t.to }

// Values liefert eine neue Kopie der Zuordnung von `map_value`; ein anderer
// Regeltyp liefert `nil`.
func (t Transformation) Values() map[string]string {
	if t.values == "" {
		return nil
	}
	values := make(map[string]string)
	for rest := t.values; rest != ""; {
		var from, to string
		from, rest = readValueField(rest)
		to, rest = readValueField(rest)
		values[from] = to
	}
	return values
}

// CheckApplicable prüft die Anwendbarkeit der Regel auf eine Änderung
// (`SPEC-030`, Anwendbarkeit): ihre Spalte kommt in `columns` vor, und der
// Zielname von `rename_column` gleicht keiner Spalte aus `columns`,
// zeichengenau; `map_value` trägt keinen Zielnamen und prüft nur die Spalte.
// `columns` sind alle Spalten der Relation der Änderung, ausgeschlossene
// eingeschlossen. Die Prüfung hängt an Regel und Spaltenmenge, nie am Wert
// einer Zeile. Verletzungen: `ErrTransformationColumnMissing`,
// `ErrTransformationTargetCollides`.
func (t Transformation) CheckApplicable(columns []string) error {
	if !containsName(columns, t.column) {
		return fmt.Errorf("%w: %s", domainerrors.ErrTransformationColumnMissing, t.column)
	}
	if t.kind == TransformationRenameColumn && containsName(columns, t.to) {
		return fmt.Errorf("%w: %s", domainerrors.ErrTransformationTargetCollides, t.to)
	}
	return nil
}

// applyTransformations wertet die Regeln für eine Spalte aus, deren Wert das
// Bild trägt: liefert den Schlüssel, unter dem die Spalte im Bild steht, und
// den Wert dort. Eine Regel trifft, wenn ihre Quellspalte `column` ist; die
// erste treffende Regel entscheidet (K2 hält je Quellspalte höchstens eine),
// keine treffende Regel lässt Schlüssel und Wert unverändert. `map_value`
// lässt den Schlüssel und ersetzt den Wert, wenn die Zuordnung ihn als
// Schlüssel führt (zeichengenau, ohne Normalisierung). Die Funktion ist rein.
func applyTransformations(rules []Transformation, column, value string) (string, string) {
	for _, rule := range rules {
		if rule.column != column {
			continue
		}
		switch rule.kind {
		case TransformationRenameColumn:
			return rule.to, value
		case TransformationMapValue:
			if mapped, ok := lookupMappedValue(rule.values, value); ok {
				return column, mapped
			}
			return column, value
		}
	}
	return column, value
}

// encodeValueMap kodiert eine Zuordnung kanonisch: die Paare in aufsteigender
// Ordnung des Schlüssels, je Paar Schlüssel und Wert mit einer Länge von vier
// Byte (Big-Endian) davor. Gleiche Zuordnungen ergeben dieselbe Zeichenkette,
// die leere Zuordnung die leere Zeichenkette.
func encodeValueMap(values map[string]string) string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var encoded strings.Builder
	for _, key := range keys {
		writeValueField(&encoded, key)
		writeValueField(&encoded, values[key])
	}
	return encoded.String()
}

func writeValueField(out *strings.Builder, field string) {
	length := uint32(len(field))
	out.WriteByte(byte(length >> 24))
	out.WriteByte(byte(length >> 16))
	out.WriteByte(byte(length >> 8))
	out.WriteByte(byte(length))
	out.WriteString(field)
}

// readValueField liest ein Feld der kanonischen Kodierung und liefert es mit
// dem Rest; die Kodierung stammt ausschließlich aus `encodeValueMap`.
func readValueField(encoded string) (field, rest string) {
	length := int(uint32(encoded[0])<<24 | uint32(encoded[1])<<16 | uint32(encoded[2])<<8 | uint32(encoded[3]))
	return encoded[4 : 4+length], encoded[4+length:]
}

// lookupMappedValue sucht `value` als Schlüssel der kodierten Zuordnung, ohne
// sie zu dekodieren. Die Suche ist linear in der Zahl der Zuordnungen.
func lookupMappedValue(encoded, value string) (string, bool) {
	for rest := encoded; rest != ""; {
		var from, to string
		from, rest = readValueField(rest)
		to, rest = readValueField(rest)
		if from == value {
			return to, true
		}
	}
	return "", false
}
