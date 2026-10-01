package model

import (
	"fmt"
	"strings"

	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

// maxRouteTargetLength ist die Länge eines Zielnamens in Zeichen
// (`SPEC-032`, Zielname).
const maxRouteTargetLength = 63

// RouteTarget trägt den Namen eines Zustellziels (`SPEC-032`, Zielname): 1 bis
// 63 Zeichen, das erste aus `a`–`z` und `0`–`9`, die übrigen aus `a`–`z`,
// `0`–`9`, `_` und `-`. Die leere Zeichenkette ist „kein Ziel" (in der
// Datenbank `NULL`) und kein gültiger Name.
type RouteTarget string

// IsValidRouteTarget meldet, ob `raw` ein Zielname ist. Die eine Prüffunktion
// der Domäne für jeden Pfad, der einen Namen entgegennimmt.
func IsValidRouteTarget(raw string) bool {
	if raw == "" || len(raw) > maxRouteTargetLength {
		return false
	}
	for i := 0; i < len(raw); i++ {
		c := raw[i]
		alnum := (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')
		if !alnum && (i == 0 || (c != '_' && c != '-')) {
			return false
		}
	}
	return true
}

// NewRouteTarget legt einen Zielnamen an; ein Name außerhalb des Alphabets
// endet als `ErrInvalidRouteTarget`, er wird nicht gefaltet oder gekürzt.
func NewRouteTarget(raw string) (RouteTarget, error) {
	if !IsValidRouteTarget(raw) {
		return "", domainerrors.ErrInvalidRouteTarget
	}
	return RouteTarget(raw), nil
}

// RouteCondition ist die Bedingung `when` einer Routing-Regel: der Textwert
// der Spalte `Column` ist zeichengenau gleich `Equals`.
type RouteCondition struct {
	Column string
	Equals string
}

// RouteRule trägt eine Routing-Regel einer Tabelle (`SPEC-032`): Name, Ziel,
// Auswertungsreihenfolge und die optionale Bedingung. Die Felder sind nicht
// exportiert, ein Wert entsteht über `NewRouteRule` und ist danach
// unveränderlich; alle Felder sind vergleichbar, `RouteRule` bleibt über `==`
// vergleichbar.
type RouteRule struct {
	name       string
	target     RouteTarget
	order      int64
	hasWhen    bool
	whenColumn string
	whenEquals string
}

// NewRouteRule legt eine Routing-Regel an (`SPEC-032`):
//   - `name` ist nicht leer (`ErrEmptyIdentifier`); das Alphabet des Namens
//     prüft der Aufrufer, der ihn entgegennimmt;
//   - `target` liegt im Alphabet des Zielnamens (`ErrInvalidRouteTarget`);
//   - `order` ist mindestens 1 und `when.Column` ist nicht leer und trägt kein
//     U+0000 (beide Verletzungen: `ErrInvalidRoute`). `when.Equals` ist
//     unbeschränkt, die leere Zeichenkette zulässig.
//
// `when == nil` ist eine Regel ohne Bedingung: sie trifft jede Änderung der
// Tabelle.
func NewRouteRule(name, target string, order int64, when *RouteCondition) (RouteRule, error) {
	if name == "" {
		return RouteRule{}, fmt.Errorf("Regelname: %w", domainerrors.ErrEmptyIdentifier)
	}
	checked, err := NewRouteTarget(target)
	if err != nil {
		return RouteRule{}, err
	}
	if order < 1 {
		return RouteRule{}, fmt.Errorf("order: %w", domainerrors.ErrInvalidRoute)
	}
	rule := RouteRule{name: name, target: checked, order: order}
	if when != nil {
		if when.Column == "" || strings.IndexByte(when.Column, 0) >= 0 {
			return RouteRule{}, fmt.Errorf("when.column: %w", domainerrors.ErrInvalidRoute)
		}
		rule.hasWhen = true
		rule.whenColumn = when.Column
		rule.whenEquals = when.Equals
	}
	return rule, nil
}

// Name trägt den Regelnamen, mit dem die Regel je Tabelle adressiert wird.
func (r RouteRule) Name() string { return r.name }

// Target trägt das Ziel, das die Regel einer treffenden Änderung zuweist.
func (r RouteRule) Target() RouteTarget { return r.target }

// Order trägt die Auswertungsreihenfolge; eine kleinere Zahl wird zuerst
// geprüft.
func (r RouteRule) Order() int64 { return r.order }

// Condition liefert die Bedingung `when` der Regel; `false`, wenn die Regel
// keine trägt.
func (r RouteRule) Condition() (RouteCondition, bool) {
	return RouteCondition{Column: r.whenColumn, Equals: r.whenEquals}, r.hasWhen
}

// CheckApplicable prüft die Anwendbarkeit der Regel auf eine Änderung
// (`SPEC-032`, Anwendbarkeit): eine Regel ohne `when` ist auf jede Änderung
// anwendbar, eine Regel mit `when` ist es, wenn `when.column` in `columns`
// vorkommt. `columns` sind alle Spalten der Relation der Änderung,
// ausgeschlossene eingeschlossen. Die Prüfung hängt an Regel und Spaltenmenge,
// nie am Wert einer Zeile. Verletzung: `ErrRoutingColumnMissing`.
func (r RouteRule) CheckApplicable(columns []string) error {
	if r.hasWhen && !containsName(columns, r.whenColumn) {
		return fmt.Errorf("%w: %s", domainerrors.ErrRoutingColumnMissing, r.whenColumn)
	}
	return nil
}

// CheckNotExcluded prüft die Regel gegen die ausgeschlossenen Spalten der
// Tabelle (R3, `SPEC-019`): eine Bedingung auf eine ausgeschlossene Spalte
// verriete über das Ziel eine Eigenschaft des ausgeschlossenen Werts.
// Verletzung: `ErrRoutingColumnExcluded`.
func (r RouteRule) CheckNotExcluded(excluded []string) error {
	if r.hasWhen && containsName(excluded, r.whenColumn) {
		return fmt.Errorf("%w: %s", domainerrors.ErrRoutingColumnExcluded, r.whenColumn)
	}
	return nil
}

// matches meldet, ob die Regel die Zeile trifft, die `values` zu `columns`
// trägt. Ein Wert, den das Bild nicht trägt (nil: NULL oder unverändertes
// TOAST, eine Spalte ohne Wert, eine Spalte außerhalb von `columns`), trifft
// nicht — eine Abwesenheit ist keine Nichtanwendbarkeit.
func (r RouteRule) matches(columns []string, values []*string) bool {
	if !r.hasWhen {
		return true
	}
	for i, column := range columns {
		if column != r.whenColumn {
			continue
		}
		return i < len(values) && values[i] != nil && *values[i] == r.whenEquals
	}
	return false
}

// EvaluateRoute bestimmt das Zustellziel einer Änderung (`SPEC-032`,
// Auswertung): unter den treffenden Regeln gewinnt die mit der kleinsten
// `order`; bei gleicher `order` (R2 hält sie je Tabelle eindeutig) die in
// `rules` zuerst genannte. Trifft keine Regel, ist das Ziel leer. `values`
// ist die Bildbasis — das Neu-Bild bei INSERT und UPDATE, das Alt-Bild bei
// DELETE — als Quellwerte vor jeder Transformation; `values[i]` gehört zu
// `columns[i]`. Die Funktion ist rein und liest `rules` nur.
func EvaluateRoute(rules []RouteRule, columns []string, values []*string) RouteTarget {
	var best *RouteRule
	for i := range rules {
		rule := &rules[i]
		if !rule.matches(columns, values) {
			continue
		}
		if best == nil || rule.order < best.order {
			best = rule
		}
	}
	if best == nil {
		return ""
	}
	return best.target
}
