package model

import (
	stderrors "errors"
	"strconv"
	"testing"

	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

// ParseRouteSpec liest die Regelform (`SPEC-032`) und gibt sie zeichengenau
// zurück; die Grenzen von `order` (1 und `MaxRouteOrder`) sind zulässig, die
// leere Zeichenkette als `equals` ebenso. Rot färbende Mutation je Eingabe:
// `order < 1` zu `order < 2` (der Fall `order 1` endet abgelehnt),
// `order > MaxRouteOrder` zu `order >= MaxRouteOrder` (der Fall `MaxRouteOrder`
// endet abgelehnt), die Prüfung `column == ""` für `equals` wiederholen (der
// Fall `equals ""` endet abgelehnt).
func TestParseRouteSpecAcceptsTheRuleForm(t *testing.T) {
	for _, tc := range []struct {
		label   string
		text    string
		target  RouteTarget
		order   int64
		column  string
		equals  string
		hasWhen bool
	}{
		{"ohne when", `{"target": "sonstige", "order": 100}`, "sonstige", 100, "", "", false},
		{"mit when", `{"target": "eu", "order": 10, "when": {"column": "region", "equals": "eu"}}`, "eu", 10, "region", "eu", true},
		{"order 1", `{"target": "a", "order": 1}`, "a", 1, "", "", false},
		{"order an der Obergrenze", `{"target": "a", "order": ` + strconv.FormatInt(MaxRouteOrder, 10) + `}`, "a", MaxRouteOrder, "", "", false},
		{"equals leer", `{"target": "a", "order": 5, "when": {"column": "region", "equals": ""}}`, "a", 5, "region", "", true},
		{"Spaltenname mit Großbuchstaben und Leerzeichen", `{"target": "a", "order": 5, "when": {"column": "Re Gion", "equals": "x"}}`, "a", 5, "Re Gion", "x", true},
		{"Schlüssel in anderer Reihenfolge", `{"order": 5, "when": {"equals": "x", "column": "c"}, "target": "a"}`, "a", 5, "c", "x", true},
	} {
		spec, err := ParseRouteSpec(tc.text)
		if err != nil {
			t.Fatalf("%s: ParseRouteSpec = %v, wollen nil", tc.label, err)
		}
		condition, conditional := spec.Condition()
		if spec.Target() != tc.target || spec.Order() != tc.order || conditional != tc.hasWhen || condition.Column != tc.column || condition.Equals != tc.equals {
			t.Fatalf("%s: Regelform = %+v / %+v, wollen %q %d %q=%q", tc.label, spec, condition, tc.target, tc.order, tc.column, tc.equals)
		}
		rule, err := spec.Build("r")
		if err != nil {
			t.Fatalf("%s: Build = %v, wollen nil", tc.label, err)
		}
		if rule.Name() != "r" || rule.Target() != tc.target || rule.Order() != tc.order {
			t.Fatalf("%s: Regel = %+v, wollen r %q %d", tc.label, rule, tc.target, tc.order)
		}
	}
}

// ParseRouteSpec lehnt jede Abweichung von der Regelform mit dem Grund ab, den
// die Fehlertext-Tabelle von `SPEC-019` nennt; die Eingabe jedes Falls ist die
// Verletzung. Die Ablehnung mit Wert (Schlüssel, Zielname) trägt ihn als
// `RouteSpecError.Detail`. Rot färbende Mutation je Zeile: die Prüfung der
// Form streichen bzw. die Prüfreihenfolge umkehren (Alphabet vor Form,
// Schlüssel nach Form).
func TestParseRouteSpecRejectsWithTheReason(t *testing.T) {
	for _, tc := range []struct {
		label  string
		text   string
		reason error
		detail string
	}{
		{"leerer Text", ``, domainerrors.ErrInvalidRuleSpec, ""},
		{"JSON-null", `null`, domainerrors.ErrInvalidRuleSpec, ""},
		{"Array", `[]`, domainerrors.ErrInvalidRuleSpec, ""},
		{"Nachlauf", `{"target": "a", "order": 1} x`, domainerrors.ErrInvalidRuleSpec, ""},
		{"ungültiges UTF-8", "{\"target\": \"a\xff\", \"order\": 1}", domainerrors.ErrInvalidRuleSpec, ""},
		{"unbekannter Schlüssel", `{"target": "a", "order": 1, "kind": "x"}`, domainerrors.ErrUnknownRuleSpecKey, "kind"},
		{"unbekannter Schlüssel in when", `{"target": "a", "order": 1, "when": {"column": "c", "equals": "e", "op": "eq"}}`, domainerrors.ErrUnknownRuleSpecKey, "op"},
		{"Schlüssel von rule_spec vor Schlüssel von when", `{"zzz": 1, "target": "a", "order": 1, "when": {"aaa": 1}}`, domainerrors.ErrUnknownRuleSpecKey, "zzz"},
		{"target fehlt", `{"order": 1}`, domainerrors.ErrInvalidRuleSpec, ""},
		{"order ohne Ziffernform", `{"target": "a", "order": 1.0}`, domainerrors.ErrInvalidRuleSpec, ""},
		{"order über der Obergrenze", `{"target": "a", "order": ` + strconv.FormatInt(MaxRouteOrder+1, 10) + `}`, domainerrors.ErrInvalidRuleSpec, ""},
		{"order 0", `{"target": "a", "order": 0}`, domainerrors.ErrInvalidRuleSpec, ""},
		{"when ohne Objekt", `{"target": "a", "order": 1, "when": []}`, domainerrors.ErrInvalidRuleSpec, ""},
		{"when.equals fehlt", `{"target": "a", "order": 1, "when": {"column": "c"}}`, domainerrors.ErrInvalidRuleSpec, ""},
		{"Zielname außerhalb des Alphabets", `{"target": "Ost", "order": 1}`, domainerrors.ErrInvalidRouteTarget, "Ost"},
		{"Form vor Alphabet", `{"target": "Ost", "order": 0}`, domainerrors.ErrInvalidRuleSpec, ""},
	} {
		_, err := ParseRouteSpec(tc.text)
		if !stderrors.Is(err, tc.reason) {
			t.Fatalf("%s: Fehler = %v, wollen Grund %v", tc.label, err, tc.reason)
		}
		var specErr *RouteSpecError
		if tc.detail != "" {
			if !stderrors.As(err, &specErr) || specErr.Detail != tc.detail {
				t.Fatalf("%s: Fehler = %v, wollen RouteSpecError mit Wert %q", tc.label, err, tc.detail)
			}
		} else if stderrors.As(err, &specErr) {
			t.Fatalf("%s: Fehler = %v, wollen keinen Wert in der Adresse", tc.label, err)
		}
	}
}

func routeRule(t *testing.T, name string, order int64, column string) RouteRule {
	t.Helper()
	var when *RouteCondition
	if column != "" {
		when = &RouteCondition{Column: column, Equals: "e"}
	}
	rule, err := NewRouteRule(name, "t", order, when)
	if err != nil {
		t.Fatal(err)
	}
	return rule
}

// FoldRoutes leitet den Regelstand aus den `applied`-Zeilen in der Ordnung der
// Eingabe ab: `set_route` trägt ein, eine gleichnamige Regel wird an ihrer
// Stelle ersetzt, `remove_route` nimmt heraus, eine fremde Art trägt nichts;
// ohne Regel ist die Rückgabe `nil`. Die Ordnung der Eingabe entscheidet:
// dasselbe Paar in umgekehrter Reihenfolge liefert den umgekehrten Stand. Rot
// färbende Mutation: `remove_route` als Eintrag behandeln (der Stand trägt die
// entfernte Regel), das Ersetzen nach Namen durch Anhängen ersetzen (zwei
// Regeln gleichen Namens), die Eingabe vor dem Falten sortieren.
func TestFoldRoutesFollowsTheOrderOfTheRecords(t *testing.T) {
	set := func(name string, order int) RouteRecord {
		return RouteRecord{Kind: AdministrationRequestSetRoute, Name: name, Spec: `{"target": "t` + strconv.Itoa(order) + `", "order": ` + strconv.Itoa(order) + `}`}
	}
	remove := func(name string) RouteRecord { return RouteRecord{Kind: AdministrationRequestRemoveRoute, Name: name} }

	rules, err := FoldRoutes([]RouteRecord{set("a", 10), set("b", 20), remove("a"), set("a", 30)})
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) != 2 || rules[0].Name() != "b" || rules[1].Name() != "a" || rules[1].Order() != 30 {
		t.Fatalf("Stand = %+v, wollen b(20), a(30)", rules)
	}
	rules, err = FoldRoutes([]RouteRecord{set("a", 10), set("a", 11)})
	if err != nil || len(rules) != 1 || rules[0].Order() != 11 {
		t.Fatalf("Stand = %+v / %v, wollen eine Regel a(11)", rules, err)
	}
	rules, err = FoldRoutes([]RouteRecord{set("a", 10), remove("a")})
	if err != nil || rules != nil {
		t.Fatalf("Stand = %+v / %v, wollen nil", rules, err)
	}
	rules, err = FoldRoutes([]RouteRecord{{Kind: AdministrationRequestSetTransformation, Name: "x", Spec: "{}"}, remove("nie")})
	if err != nil || rules != nil {
		t.Fatalf("Stand = %+v / %v, wollen nil für fremde Arten und einen nicht geführten Namen", rules, err)
	}
	if _, err := FoldRoutes([]RouteRecord{{Kind: AdministrationRequestSetRoute, Name: "kaputt", Spec: `{"target": "a"}`}}); !stderrors.Is(err, domainerrors.ErrInvalidRuleSpec) {
		t.Fatalf("Fehler = %v, wollen ErrInvalidRuleSpec: der Stand wird nie um eine Zeile verkürzt", err)
	}
}

// RoutesConditionOn vergleicht die Spalte zeichengenau und kennt nur Regeln
// mit `when`. Rot färbende Mutation: den Vergleich auf `strings.EqualFold`
// stellen (der Fall `Region` endet gefunden), die Prüfung `hasWhen` streichen
// (die Abschlussregel trägt die leere Spalte; der Fall `""` endet gefunden).
func TestRoutesConditionOn(t *testing.T) {
	rules := []RouteRule{routeRule(t, "a", 1, "region"), routeRule(t, "b", 2, "")}
	for column, want := range map[string]bool{"region": true, "Region": false, "status": false, "": false} {
		if got := RoutesConditionOn(rules, column); got != want {
			t.Fatalf("RoutesConditionOn(%q) = %t, wollen %t", column, got, want)
		}
	}
	if RoutesConditionOn(nil, "region") {
		t.Fatal("RoutesConditionOn(nil) = true")
	}
}

// CheckIdentity und CheckArrangement tragen R1/R2 und R4/R5 mit dem Wert der
// Adresse (`SPEC-019`); die Reihenfolge innerhalb der Funktion ist die der
// Tabelle. Rot färbende Mutation: R1 und R2 vertauschen (der Fall „R1 geht R2
// voraus" endet mit dem Text von R2).
func TestRouteRuleChecksCarryTheAddressValue(t *testing.T) {
	existing := []RouteRule{routeRule(t, "a", 10, "region"), routeRule(t, "z", 100, "")}
	for _, tc := range []struct {
		label  string
		check  func(RouteRule, []RouteRule) error
		rule   RouteRule
		reason error
		detail string
	}{
		{"R1", RouteRule.CheckIdentity, routeRule(t, "a", 11, "q"), domainerrors.ErrRuleNameTaken, "a"},
		{"R2", RouteRule.CheckIdentity, routeRule(t, "n", 10, "q"), domainerrors.ErrRouteOrderTaken, "10"},
		{"R1 geht R2 voraus", RouteRule.CheckIdentity, routeRule(t, "a", 10, "q"), domainerrors.ErrRuleNameTaken, "a"},
		{"R4 bereits geführt", RouteRule.CheckArrangement, routeRule(t, "n", 200, ""), domainerrors.ErrRouteWithoutWhenTaken, "z"},
		{"R4 hinter der Abschlussregel", RouteRule.CheckArrangement, routeRule(t, "n", 150, "q"), domainerrors.ErrRouteBehindWithoutWhen, "z"},
		{"R5", RouteRule.CheckArrangement, routeRule(t, "n", 50, "region"), domainerrors.ErrRouteConditionTaken, "region"},
	} {
		err := tc.check(tc.rule, existing)
		var specErr *RouteSpecError
		if !stderrors.Is(err, tc.reason) || !stderrors.As(err, &specErr) || specErr.Detail != tc.detail {
			t.Fatalf("%s: Fehler = %v, wollen %v mit Wert %q", tc.label, err, tc.reason, tc.detail)
		}
	}
	if err := routeRule(t, "n", 50, "q").CheckArrangement(existing); err != nil {
		t.Fatalf("Gegenprobe: Fehler = %v, wollen nil", err)
	}
	if err := routeRule(t, "n", 50, "q").CheckIdentity(existing); err != nil {
		t.Fatalf("Gegenprobe: Fehler = %v, wollen nil", err)
	}
}
