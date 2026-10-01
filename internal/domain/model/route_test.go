package model

import (
	stderrors "errors"
	"strings"
	"testing"

	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

// mustRoute legt eine Routing-Regel an oder bricht den Test ab.
func mustRoute(t testing.TB, name, target string, order int64, when *RouteCondition) RouteRule {
	t.Helper()
	rule, err := NewRouteRule(name, target, order, when)
	if err != nil {
		t.Fatalf("NewRouteRule(%q, %q, %d, %v): %v", name, target, order, when, err)
	}
	return rule
}

// text trägt einen Wert-Zeiger; nil bleibt die Abwesenheit.
func text(value string) *string { return &value }

// Das Alphabet des Zielnamens (`SPEC-032`), je Eingabe ein Fall. Rot färbende
// Mutationen je Fall: die Zeichenklasse der Folgezeichen um `.` erweitern
// (Punkt), die Prüfung des ersten Zeichens auf die Folgeklasse lockern
// (führendes `_`/`-`), die Längengrenze `> 63` auf `> 64` verschieben.
func TestNewRouteTargetAlphabet(t *testing.T) {
	cases := []struct {
		name  string
		input string
		valid bool
	}{
		{"Kleinbuchstaben", "eu", true},
		{"Ziffer am Anfang", "0eu", true},
		{"Unterstrich und Bindestrich in der Folge", "eu_west-1", true},
		{"ein Zeichen", "a", true},
		{"63 Zeichen", strings.Repeat("a", 63), true},
		{"64 Zeichen", strings.Repeat("a", 64), false},
		{"leer", "", false},
		{"Großbuchstabe", "Eu", false},
		{"führender Unterstrich", "_eu", false},
		{"führender Bindestrich", "-eu", false},
		{"Punkt", "eu.west", false},
		{"NATS-Platzhalter Stern", "eu*", false},
		{"NATS-Platzhalter Größer", "eu>", false},
		{"Leerzeichen", "eu west", false},
		{"Umlaut", "ärger", false},
		{"U+0000", "eu\x00", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			target, err := NewRouteTarget(c.input)
			if c.valid {
				if err != nil || string(target) != c.input {
					t.Fatalf("NewRouteTarget(%q) = %q, %v, wollen den Namen unverändert", c.input, target, err)
				}
			} else if !stderrors.Is(err, domainerrors.ErrInvalidRouteTarget) {
				t.Fatalf("NewRouteTarget(%q): Fehler = %v, wollen %v", c.input, err, domainerrors.ErrInvalidRouteTarget)
			}
			if IsValidRouteTarget(c.input) != c.valid {
				t.Fatalf("IsValidRouteTarget(%q) = %v, wollen %v", c.input, !c.valid, c.valid)
			}
		})
	}
}

// Der Konstruktor erzwingt die Invarianten einer einzelnen Regel
// (`SPEC-032`), je Eingabe ein Fall. Rot färbende Mutationen je Fall: die
// Prüfung der Eingabe im Konstruktor entfernen bzw. die Grenze verschieben
// (`order < 1` auf `order < 0`).
func TestNewRouteRuleInvariants(t *testing.T) {
	cases := []struct {
		name    string
		rule    func() (RouteRule, error)
		wantErr error
	}{
		{"ohne Bedingung", func() (RouteRule, error) { return NewRouteRule("rest", "sonstige", 100, nil) }, nil},
		{"mit Bedingung", func() (RouteRule, error) {
			return NewRouteRule("eu", "eu", 10, &RouteCondition{Column: "region", Equals: "eu"})
		}, nil},
		{"leerer Vergleichswert ist zulässig", func() (RouteRule, error) {
			return NewRouteRule("leer", "x", 1, &RouteCondition{Column: "region", Equals: ""})
		}, nil},
		{"order 1 ist die Untergrenze", func() (RouteRule, error) { return NewRouteRule("r", "x", 1, nil) }, nil},
		{"Regelname leer", func() (RouteRule, error) { return NewRouteRule("", "x", 1, nil) }, domainerrors.ErrEmptyIdentifier},
		{"Ziel außerhalb des Alphabets", func() (RouteRule, error) { return NewRouteRule("r", "X", 1, nil) }, domainerrors.ErrInvalidRouteTarget},
		{"order 0", func() (RouteRule, error) { return NewRouteRule("r", "x", 0, nil) }, domainerrors.ErrInvalidRoute},
		{"order negativ", func() (RouteRule, error) { return NewRouteRule("r", "x", -5, nil) }, domainerrors.ErrInvalidRoute},
		{"Bedingungsspalte leer", func() (RouteRule, error) {
			return NewRouteRule("r", "x", 1, &RouteCondition{Column: "", Equals: "a"})
		}, domainerrors.ErrInvalidRoute},
		{"Bedingungsspalte mit U+0000", func() (RouteRule, error) {
			return NewRouteRule("r", "x", 1, &RouteCondition{Column: "re\x00gion", Equals: "a"})
		}, domainerrors.ErrInvalidRoute},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			rule, err := c.rule()
			if c.wantErr != nil {
				if !stderrors.Is(err, c.wantErr) {
					t.Fatalf("Fehler = %v, wollen %v", err, c.wantErr)
				}
				if rule != (RouteRule{}) {
					t.Fatalf("Regel %+v bei Fehler, wollen den Nullwert", rule)
				}
				return
			}
			if err != nil {
				t.Fatalf("NewRouteRule: %v", err)
			}
			if rule.Name() == "" || rule.Target() == "" || rule.Order() < 1 {
				t.Fatalf("Regel %+v trägt ihre Eingaben nicht", rule)
			}
		})
	}
}

// Die Regel gibt ihre Bedingung zurück und kopiert die Eingabe: eine spätere
// Änderung des übergebenen Werts wirkt nicht auf sie.
func TestRouteRuleKeepsItsConditionIndependentOfTheInput(t *testing.T) {
	when := &RouteCondition{Column: "region", Equals: "eu"}
	rule := mustRoute(t, "eu", "eu", 10, when)
	when.Column = "geändert"
	got, ok := rule.Condition()
	if !ok || got.Column != "region" || got.Equals != "eu" {
		t.Fatalf("Condition() = %+v, %v, wollen region/eu", got, ok)
	}
	if _, ok := mustRoute(t, "rest", "x", 99, nil).Condition(); ok {
		t.Fatalf("eine Regel ohne Bedingung meldet eine")
	}
}

// Die Auswertung (`SPEC-032`, Auswertung, Abwesenheit, Bildbasis): je Fall
// eine Zeile über die Spalten `id`, `name`, `region`. Rot färbende Mutationen
// je Fall stehen am Fall.
func TestEvaluateRoute(t *testing.T) {
	columns := []string{"id", "name", "region"}
	eu := mustRoute(t, "eu", "eu", 10, &RouteCondition{Column: "region", Equals: "eu"})
	us := mustRoute(t, "us", "us", 20, &RouteCondition{Column: "region", Equals: "us"})
	rest := mustRoute(t, "rest", "sonstige", 100, nil)
	euLate := mustRoute(t, "eu-spät", "eu-spaet", 30, &RouteCondition{Column: "region", Equals: "eu"})
	emptyEquals := mustRoute(t, "leer", "leer", 5, &RouteCondition{Column: "region", Equals: ""})

	cases := []struct {
		name   string
		rules  []RouteRule
		values []*string
		want   RouteTarget
	}{
		// Happy
		{"Bedingung trifft", []RouteRule{eu, us}, []*string{text("7"), text("Ada"), text("eu")}, "eu"},
		{"zweite Regel trifft", []RouteRule{eu, us}, []*string{text("7"), text("Ada"), text("us")}, "us"},
		{"Regel ohne Bedingung trifft jede Zeile", []RouteRule{rest}, []*string{text("7"), text("Ada"), text("zz")}, "sonstige"},
		{"Abschlussregel fängt, was keine Bedingung trifft", []RouteRule{eu, us, rest}, []*string{text("7"), text("Ada"), text("zz")}, "sonstige"},
		// Boundary: Reihenfolge
		{"die kleinere order gewinnt, obwohl später gelistet", []RouteRule{euLate, eu}, []*string{text("7"), text("Ada"), text("eu")}, "eu"},
		{"die kleinere order gewinnt vor der Regel ohne Bedingung", []RouteRule{rest, eu}, []*string{text("7"), text("Ada"), text("eu")}, "eu"},
		{"gleiche order: die zuerst gelistete gewinnt", []RouteRule{
			mustRoute(t, "a", "erste", 7, nil), mustRoute(t, "b", "zweite", 7, nil),
		}, []*string{text("1"), nil, nil}, "erste"},
		// Boundary: Abwesenheit ist Nicht-Treffer, die nächste Regel wird geprüft
		{"NULL trifft nicht, die Abschlussregel bestimmt", []RouteRule{eu, rest}, []*string{text("7"), text("Ada"), nil}, "sonstige"},
		{"NULL trifft nicht, ohne Abschlussregel bleibt das Ziel leer", []RouteRule{eu, us}, []*string{text("7"), text("Ada"), nil}, ""},
		{"Bild kürzer als die Spaltenliste: Wert abwesend", []RouteRule{eu, rest}, []*string{text("7")}, "sonstige"},
		{"kein Bild: nur die Regel ohne Bedingung trifft", []RouteRule{eu, rest}, nil, "sonstige"},
		{"kein Bild und keine Regel ohne Bedingung", []RouteRule{eu}, nil, ""},
		{"leerer Vergleichswert trifft den leeren Wert, nicht die Abwesenheit", []RouteRule{emptyEquals, rest}, []*string{text("7"), text("Ada"), text("")}, "leer"},
		{"leerer Vergleichswert trifft die Abwesenheit nicht", []RouteRule{emptyEquals, rest}, []*string{text("7"), text("Ada"), nil}, "sonstige"},
		// Boundary: zeichengenauer Vergleich
		{"Groß-/Kleinschreibung wird nicht gefaltet", []RouteRule{eu}, []*string{text("7"), text("Ada"), text("EU")}, ""},
		{"Leerzeichen am Rand werden nicht gekürzt", []RouteRule{eu}, []*string{text("7"), text("Ada"), text("eu ")}, ""},
		// Negative
		{"keine Regel", nil, []*string{text("7"), text("Ada"), text("eu")}, ""},
		{"keine Regel trifft", []RouteRule{eu, us}, []*string{text("7"), text("Ada"), text("zz")}, ""},
		{"Bedingungsspalte nicht in der Spaltenliste trifft nicht", []RouteRule{
			mustRoute(t, "x", "x", 1, &RouteCondition{Column: "fehlt", Equals: "eu"}),
		}, []*string{text("7"), text("Ada"), text("eu")}, ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			before := append([]RouteRule(nil), c.rules...)
			if got := EvaluateRoute(c.rules, columns, c.values); got != c.want {
				t.Fatalf("EvaluateRoute = %q, wollen %q", got, c.want)
			}
			for i := range before {
				if c.rules[i] != before[i] {
					t.Fatalf("die Auswertung hat die Regelliste verändert")
				}
			}
		})
	}
}

// Das Ziel hängt nicht an der Listenposition, solange `order` eindeutig ist
// (R2): jede Permutation derselben Regelmenge liefert für jede Zeile dasselbe
// Ziel. Rot färbende Mutation: in `EvaluateRoute` den Vergleich
// `rule.order < best.order` entfernen — dann gewinnt die zuerst gelistete
// Regel und die Permutationen weichen ab.
func TestEvaluateRouteIsIndependentOfListOrder(t *testing.T) {
	columns := []string{"id", "region", "tier"}
	rules := []RouteRule{
		mustRoute(t, "a", "platin", 10, &RouteCondition{Column: "tier", Equals: "platin"}),
		mustRoute(t, "b", "eu", 20, &RouteCondition{Column: "region", Equals: "eu"}),
		mustRoute(t, "c", "gold", 30, &RouteCondition{Column: "tier", Equals: "gold"}),
		mustRoute(t, "d", "rest", 100, nil),
	}
	rows := [][]*string{
		{text("1"), text("eu"), text("platin")},
		{text("2"), text("eu"), text("gold")},
		{text("3"), text("us"), text("gold")},
		{text("4"), text("us"), text("silver")},
		{text("5"), nil, nil},
	}
	want := make([]RouteTarget, len(rows))
	for i, row := range rows {
		want[i] = EvaluateRoute(rules, columns, row)
	}
	if want[0] != "platin" || want[1] != "eu" || want[2] != "gold" || want[3] != "rest" || want[4] != "rest" {
		t.Fatalf("Ausgangsziele = %v", want)
	}
	permute(rules, func(permutation []RouteRule) {
		for i, row := range rows {
			if got := EvaluateRoute(permutation, columns, row); got != want[i] {
				t.Fatalf("Permutation %v: Zeile %d = %q, wollen %q", names(permutation), i, got, want[i])
			}
		}
	})
}

func names(rules []RouteRule) []string {
	out := make([]string, len(rules))
	for i, rule := range rules {
		out[i] = rule.Name()
	}
	return out
}

// permute ruft `visit` für jede Permutation von `rules` (Heap-Algorithmus auf
// einer Kopie).
func permute(rules []RouteRule, visit func([]RouteRule)) {
	work := append([]RouteRule(nil), rules...)
	var heap func(k int)
	heap = func(k int) {
		if k == 1 {
			visit(work)
			return
		}
		for i := 0; i < k; i++ {
			heap(k - 1)
			if k%2 == 0 {
				work[i], work[k-1] = work[k-1], work[i]
			} else {
				work[0], work[k-1] = work[k-1], work[0]
			}
		}
	}
	heap(len(work))
}

// Die Anwendbarkeit hängt an Regel und Spaltenmenge, nie am Wert
// (`SPEC-032`, Anwendbarkeit). Rot färbende Mutation: in `CheckApplicable`
// die Bedingung `r.hasWhen &&` entfernen — dann ist eine Regel ohne `when`
// nicht anwendbar; oder `!containsName` zu `containsName` drehen.
func TestRouteRuleCheckApplicable(t *testing.T) {
	columns := []string{"id", "region"}
	withWhen := mustRoute(t, "eu", "eu", 10, &RouteCondition{Column: "region", Equals: "eu"})
	missing := mustRoute(t, "x", "x", 20, &RouteCondition{Column: "fehlt", Equals: "eu"})
	always := mustRoute(t, "rest", "rest", 100, nil)

	if err := withWhen.CheckApplicable(columns); err != nil {
		t.Fatalf("Bedingung auf vorhandene Spalte: %v", err)
	}
	if err := always.CheckApplicable(nil); err != nil {
		t.Fatalf("Regel ohne Bedingung ist auf jede Änderung anwendbar: %v", err)
	}
	err := missing.CheckApplicable(columns)
	if !stderrors.Is(err, domainerrors.ErrRoutingColumnMissing) || !strings.Contains(err.Error(), "fehlt") {
		t.Fatalf("fehlende Spalte: Fehler = %v, wollen %v mit dem Spaltennamen", err, domainerrors.ErrRoutingColumnMissing)
	}
	if err := withWhen.CheckApplicable([]string{"Region"}); !stderrors.Is(err, domainerrors.ErrRoutingColumnMissing) {
		t.Fatalf("Groß-/Kleinschreibung wird nicht gefaltet: Fehler = %v", err)
	}
}

// Eigenschaftstest zu R3, Richtung Regel → Ausschluss (`SPEC-019`): über alle
// Teilmengen der Spalten einer Tabelle als Ausschlussmenge nennt eine Regel
// genau dann eine ausgeschlossene Spalte, wenn `CheckNotExcluded` sie
// abweist; eine Regel ohne Bedingung nennt nie eine. Rot färbende Mutation:
// in `CheckNotExcluded` die Bedingung `r.hasWhen &&` entfernen oder die
// Suche auf `r.name` statt `r.whenColumn` richten.
func TestRouteRuleCheckNotExcludedHoldsForEveryExclusionSet(t *testing.T) {
	columns := []string{"id", "name", "region", "secret"}
	rules := []RouteRule{mustRoute(t, "rest", "rest", 100, nil)}
	for _, column := range columns {
		rules = append(rules, mustRoute(t, "on-"+column, "t", int64(len(rules)+1), &RouteCondition{Column: column, Equals: "v"}))
	}
	for mask := 0; mask < 1<<len(columns); mask++ {
		var excluded []string
		for i, column := range columns {
			if mask&(1<<i) != 0 {
				excluded = append(excluded, column)
			}
		}
		for _, rule := range rules {
			condition, hasWhen := rule.Condition()
			names := hasWhen && containsName(excluded, condition.Column)
			err := rule.CheckNotExcluded(excluded)
			if names != (err != nil) {
				t.Fatalf("Regel %q, ausgeschlossen %v: nennt ausgeschlossene Spalte = %v, Fehler = %v", rule.Name(), excluded, names, err)
			}
			if err != nil && !stderrors.Is(err, domainerrors.ErrRoutingColumnExcluded) {
				t.Fatalf("Regel %q: Fehler = %v, wollen %v", rule.Name(), err, domainerrors.ErrRoutingColumnExcluded)
			}
		}
	}
}

// `WithRouteTarget` setzt das Ziel, die leere Zeichenkette löscht es, ein Wert
// außerhalb des Alphabets wird abgelehnt und lässt den Change unverändert.
// `NewChange` lässt das Ziel leer. Rot färbende Mutation: in
// `WithRouteTarget` die Prüfung `NewRouteTarget` entfernen.
func TestChangeWithRouteTarget(t *testing.T) {
	change, err := NewChange("tx-1", "tx", "tbl", 1, OperationInsert, nil, nil, "sv")
	if err != nil {
		t.Fatalf("NewChange: %v", err)
	}
	if change.RouteTarget != "" {
		t.Fatalf("NewChange: RouteTarget = %q, wollen leer", change.RouteTarget)
	}
	routed, err := change.WithRouteTarget("eu")
	if err != nil || routed.RouteTarget != "eu" {
		t.Fatalf("WithRouteTarget(eu) = %q, %v", routed.RouteTarget, err)
	}
	cleared, err := routed.WithRouteTarget("")
	if err != nil || cleared.RouteTarget != "" {
		t.Fatalf("WithRouteTarget(\"\") = %q, %v", cleared.RouteTarget, err)
	}
	unchanged, err := routed.WithRouteTarget("EU")
	if !stderrors.Is(err, domainerrors.ErrInvalidRouteTarget) || unchanged.RouteTarget != "eu" {
		t.Fatalf("WithRouteTarget(EU) = %q, %v, wollen Ablehnung und unverändertes Ziel", unchanged.RouteTarget, err)
	}
}
