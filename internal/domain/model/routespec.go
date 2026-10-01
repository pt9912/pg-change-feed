package model

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

// MaxRouteOrder ist die größte `order`, die ein `set_route`-Antrag trägt:
// `SPEC-032` verlangt eine positive ganze Zahl und legt keine Obergrenze
// fest. Die Umsetzung wählt die Grenze der 32-Bit-Ganzzahl: sie liegt weit
// über jeder Zahl von Regeln einer Tabelle, und die `order` bleibt in jedem
// Ganzzahltyp eines Verbrauchers darstellbar.
const MaxRouteOrder int64 = math.MaxInt32

// RouteSpecError trägt zu einem Ablehnungsgrund den Wert, den die Adresse des
// Fehlertextes nennt (`SPEC-019`): den Schlüsselnamen bei
// `ErrUnknownRuleSpecKey`, den Zielnamen bei `ErrInvalidRouteTarget`, den
// Regelnamen, die `order` oder die Spalte bei den Konfliktzeilen R1, R2, R4
// und R5. `errors.Is` löst auf den Grund auf. Die Adresse mit
// `schema.table.`-Präfix bildet der Aufrufer.
type RouteSpecError struct {
	Err    error
	Detail string
}

// Error trägt den Grund und den Wert.
func (e *RouteSpecError) Error() string {
	return e.Err.Error() + ": " + e.Detail
}

// Unwrap gibt den Grund frei.
func (e *RouteSpecError) Unwrap() error { return e.Err }

// RouteSpec trägt die geprüfte Regelform eines `set_route`-Antrags
// (`SPEC-032`): Ziel, Auswertungsreihenfolge und die optionale Bedingung. Die
// Form ist gültig, die Konfliktfreiheit gegen den Regelstand einer Tabelle
// nicht geprüft. Die Felder sind nicht exportiert; ein Wert entsteht über
// `ParseRouteSpec`.
type RouteSpec struct {
	target     RouteTarget
	order      int64
	hasWhen    bool
	whenColumn string
	whenEquals string
}

// Order trägt die Auswertungsreihenfolge der Regel.
func (s RouteSpec) Order() int64 { return s.order }

// Target trägt das Ziel der Regel.
func (s RouteSpec) Target() RouteTarget { return s.target }

// Condition liefert die Bedingung `when`; `false`, wenn die Regel keine trägt.
func (s RouteSpec) Condition() (RouteCondition, bool) {
	return RouteCondition{Column: s.whenColumn, Equals: s.whenEquals}, s.hasWhen
}

var (
	routeSpecKeys = []string{"target", "order", "when"}
	routeWhenKeys = []string{"column", "equals"}
)

// ParseRouteSpec liest die Regelform aus dem JSON-Text von `rule_spec`
// (`SPEC-032`) strikt und prüft die Formzeilen in der Reihenfolge der
// Fehlertext-Tabelle der Antrags-Spec (die erste Formzeile, den Regelnamen,
// prüft `CheckRuleName`):
//
//	a) der Text ist gültiges UTF-8 und ein JSON-Objekt (der leere Text, JSON-
//	   `null` und jeder Wert ohne Objekt gehören hierher): `ErrInvalidRuleSpec`;
//	b) jeder Schlüssel von `rule_spec` und von `when` gehört zur Regelform:
//	   sonst `ErrUnknownRuleSpecKey` mit dem Schlüssel, bei mehreren dem
//	   ersten in aufsteigender Ordnung, `rule_spec` vor `when`;
//	c) `target` steht als Zeichenkette, `order` als ganze Zahl von 1 bis
//	   `MaxRouteOrder` in der Schreibweise einer Ganzzahl (ein Bruch oder
//	   ein Exponent, auch `10.0`, endet hier), `when` — wenn es steht — als
//	   Objekt mit den Zeichenketten `column` (nicht leer, ohne U+0000) und
//	   `equals`: sonst `ErrInvalidRuleSpec`;
//	d) `target` liegt im Alphabet des Zielnamens: sonst
//	   `ErrInvalidRouteTarget` mit dem Zielnamen, wie beantragt.
//
// Ein doppelter Schlüssel ist keine Eingabe dieser Funktion: der Text kommt
// aus einer `jsonb`-Spalte, die den letzten Wert je Schlüssel hält.
func ParseRouteSpec(text string) (RouteSpec, error) {
	if !utf8.ValidString(text) {
		return RouteSpec{}, fmt.Errorf("%w: kein gültiges UTF-8", domainerrors.ErrInvalidRuleSpec)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal([]byte(text), &fields); err != nil {
		return RouteSpec{}, fmt.Errorf("%w: %w", domainerrors.ErrInvalidRuleSpec, err)
	}
	if fields == nil {
		return RouteSpec{}, fmt.Errorf("%w: kein JSON-Objekt", domainerrors.ErrInvalidRuleSpec)
	}
	if key, unknown := firstUnknownKey(fields, routeSpecKeys); unknown {
		return RouteSpec{}, &RouteSpecError{Err: domainerrors.ErrUnknownRuleSpecKey, Detail: key}
	}
	var when map[string]json.RawMessage
	if raw, present := fields["when"]; present && len(raw) > 0 && raw[0] == '{' {
		if err := json.Unmarshal(raw, &when); err != nil {
			return RouteSpec{}, fmt.Errorf("%w: when: %w", domainerrors.ErrInvalidRuleSpec, err)
		}
		if key, unknown := firstUnknownKey(when, routeWhenKeys); unknown {
			return RouteSpec{}, &RouteSpecError{Err: domainerrors.ErrUnknownRuleSpecKey, Detail: key}
		}
	}

	target, ok := jsonString(fields["target"])
	if !ok {
		return RouteSpec{}, fmt.Errorf("%w: target fehlt oder ist keine Zeichenkette", domainerrors.ErrInvalidRuleSpec)
	}
	order, ok := jsonRouteOrder(fields["order"])
	if !ok {
		return RouteSpec{}, fmt.Errorf("%w: order fehlt oder ist keine ganze Zahl von 1 bis %d", domainerrors.ErrInvalidRuleSpec, MaxRouteOrder)
	}
	spec := RouteSpec{order: order}
	if _, present := fields["when"]; present {
		if when == nil {
			return RouteSpec{}, fmt.Errorf("%w: when ist kein Objekt", domainerrors.ErrInvalidRuleSpec)
		}
		column, columnOK := jsonString(when["column"])
		equals, equalsOK := jsonString(when["equals"])
		if !columnOK || !equalsOK || column == "" || strings.IndexByte(column, 0) >= 0 {
			return RouteSpec{}, fmt.Errorf("%w: when trägt column und equals nicht in der Form der Regel", domainerrors.ErrInvalidRuleSpec)
		}
		spec.hasWhen, spec.whenColumn, spec.whenEquals = true, column, equals
	}
	checked, err := NewRouteTarget(target)
	if err != nil {
		return RouteSpec{}, &RouteSpecError{Err: err, Detail: target}
	}
	spec.target = checked
	return spec, nil
}

// firstUnknownKey liefert den ersten Schlüssel von `fields` in aufsteigender
// Ordnung, der nicht in `allowed` steht.
func firstUnknownKey(fields map[string]json.RawMessage, allowed []string) (string, bool) {
	keys := make([]string, 0, len(fields))
	for key := range fields {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		if !containsName(allowed, key) {
			return key, true
		}
	}
	return "", false
}

// jsonRouteOrder liest einen JSON-Wert, der eine ganze Zahl von 1 bis
// `MaxRouteOrder` in der Schreibweise einer Ganzzahl sein muss (nur Ziffern);
// abwesend, eine Zeichenkette, ein Vorzeichen, ein Bruch und ein Exponent
// liefern `false`.
func jsonRouteOrder(raw json.RawMessage) (int64, bool) {
	if len(raw) == 0 || raw[0] < '0' || raw[0] > '9' {
		return 0, false
	}
	order, err := strconv.ParseInt(string(raw), 10, 64)
	if err != nil || order < 1 || order > MaxRouteOrder {
		return 0, false
	}
	return order, true
}

// Build legt die Regel unter dem Namen an. Der Name gehört zuvor durch
// `CheckRuleName`.
func (s RouteSpec) Build(name string) (RouteRule, error) {
	var when *RouteCondition
	if s.hasWhen {
		when = &RouteCondition{Column: s.whenColumn, Equals: s.whenEquals}
	}
	return NewRouteRule(name, string(s.target), s.order, when)
}

// CheckIdentity prüft R1 und R2 (`SPEC-019`) der Regel gegen den Regelstand
// `existing` der Tabelle, in der Reihenfolge R1, R2:
//   - R1 `ErrRuleNameTaken`: eine Regel trägt den Namen (Adresse: der Name);
//   - R2 `ErrRouteOrderTaken`: eine Regel trägt dieselbe `order` (Adresse: die
//     `order` in Dezimalschreibweise).
//
// Die Funktion ist rein.
func (r RouteRule) CheckIdentity(existing []RouteRule) error {
	for _, rule := range existing {
		if rule.name == r.name {
			return &RouteSpecError{Err: domainerrors.ErrRuleNameTaken, Detail: r.name}
		}
	}
	for _, rule := range existing {
		if rule.order == r.order {
			return &RouteSpecError{Err: domainerrors.ErrRouteOrderTaken, Detail: strconv.FormatInt(r.order, 10)}
		}
	}
	return nil
}

// CheckArrangement prüft R4 und R5 (`SPEC-019`) der Regel gegen den
// Regelstand `existing` der Tabelle, in der Reihenfolge der Fehlertext-Tabelle:
//   - eine Regel ohne `when`: `ErrRouteWithoutWhenTaken`, wenn eine Regel ohne
//     `when` geführt wird (Adresse: deren Name); `ErrRouteWithoutWhenNotLast`,
//     wenn eine geführte Regel eine höhere `order` trägt (Adresse: der Name
//     dieser Regel);
//   - eine Regel mit `when`: `ErrRouteBehindWithoutWhen`, wenn die geführte
//     Regel ohne `when` eine kleinere `order` trägt (Adresse: deren Name);
//     `ErrRouteConditionTaken`, wenn eine geführte Regel dasselbe Paar
//     (`column`, `equals`) trägt (Adresse: die Spalte).
//
// Die Funktion ist rein.
func (r RouteRule) CheckArrangement(existing []RouteRule) error {
	if !r.hasWhen {
		for _, rule := range existing {
			if !rule.hasWhen {
				return &RouteSpecError{Err: domainerrors.ErrRouteWithoutWhenTaken, Detail: rule.name}
			}
		}
		for _, rule := range existing {
			if rule.order > r.order {
				return &RouteSpecError{Err: domainerrors.ErrRouteWithoutWhenNotLast, Detail: r.name}
			}
		}
		return nil
	}
	for _, rule := range existing {
		if !rule.hasWhen && rule.order < r.order {
			return &RouteSpecError{Err: domainerrors.ErrRouteBehindWithoutWhen, Detail: rule.name}
		}
	}
	for _, rule := range existing {
		if rule.hasWhen && rule.whenColumn == r.whenColumn && rule.whenEquals == r.whenEquals {
			return &RouteSpecError{Err: domainerrors.ErrRouteConditionTaken, Detail: r.whenColumn}
		}
	}
	return nil
}

// RoutesConditionOn meldet, ob eine Regel des Regelstands ihre Bedingung auf
// die Spalte trägt (R3, Gegenrichtung, `SPEC-019`): der Vergleich ist
// zeichengenau, eine Regel ohne `when` trägt keine Spalte.
func RoutesConditionOn(rules []RouteRule, column string) bool {
	for _, rule := range rules {
		if rule.hasWhen && rule.whenColumn == column {
			return true
		}
	}
	return false
}

// RouteRecord trägt eine `applied`-Zeile der beiden Routing-Antragsarten
// einer Tabelle, wie der Regelstand sie auswertet: Art, Regelname und
// Regelform als JSON-Text (bei `remove_route` leer).
type RouteRecord struct {
	Kind AdministrationRequestKind
	Name string
	Spec string
}

// FoldRoutes leitet den Routing-Regelstand einer Tabelle aus ihren
// `applied`-Zeilen ab, die in der Ordnung `requested_at`, bei gleichem
// Zeitstempel nach der Antrags-Kennung stehen (`SPEC-019`): `set_route`
// trägt die Regel ein — eine Regel gleichen Namens wird an ihrer Stelle
// ersetzt, sonst hinten angehängt —, `remove_route` nimmt sie unter dem
// Namen heraus; jede andere Art trägt keinen Stand. Eine Zeile, deren
// Regelform nicht mehr zu einer Regel führt, endet als Fehler: der Stand
// wird nie um eine Zeile verkürzt. Die Rückgabe ist eine frische Liste
// (`nil` ohne Regel), die der Aufrufer nicht mit dem Speicher der Eingabe
// teilt. Die Funktion ist rein.
func FoldRoutes(records []RouteRecord) ([]RouteRule, error) {
	var rules []RouteRule
	for _, record := range records {
		switch record.Kind {
		case AdministrationRequestSetRoute:
			spec, err := ParseRouteSpec(record.Spec)
			if err != nil {
				return nil, fmt.Errorf("Routing-Regelstand: Regel %q: %w", record.Name, err)
			}
			rule, err := spec.Build(record.Name)
			if err != nil {
				return nil, fmt.Errorf("Routing-Regelstand: Regel %q: %w", record.Name, err)
			}
			rules = routeReplaceOrAppend(rules, rule)
		case AdministrationRequestRemoveRoute:
			rules = routeDropByName(rules, record.Name)
		}
	}
	return rules, nil
}

// routeReplaceOrAppend liefert eine neue Liste mit der Regel: eine Regel
// gleichen Namens wird an ihrer Stelle ersetzt, sonst hinten angehängt.
func routeReplaceOrAppend(rules []RouteRule, rule RouteRule) []RouteRule {
	next := make([]RouteRule, 0, len(rules)+1)
	replaced := false
	for _, existing := range rules {
		if existing.name == rule.name {
			next = append(next, rule)
			replaced = true
			continue
		}
		next = append(next, existing)
	}
	if !replaced {
		next = append(next, rule)
	}
	return next
}

// routeDropByName liefert eine neue Liste ohne die Regel unter dem Namen; sie
// ist `nil`, wenn keine Regel bleibt.
func routeDropByName(rules []RouteRule, name string) []RouteRule {
	var next []RouteRule
	for _, rule := range rules {
		if rule.name != name {
			next = append(next, rule)
		}
	}
	return next
}
