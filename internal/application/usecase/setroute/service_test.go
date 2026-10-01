package setroute_test

import (
	"context"
	stderrors "errors"
	"strconv"
	"strings"
	"testing"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/setroute"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// fakeRouting trägt den `RoutingPort` als Fake (`ADR-0030`): der Regelstand ist
// an die Quelle gebunden, die der Aufruf nennt — ein Regelstand einer anderen
// Quelle erreicht den Use Case nicht.
type fakeRouting struct {
	source model.SourceID
	rules  map[string][]model.RouteRule // schema.table → Regelstand
	err    error
	calls  int
}

func (f *fakeRouting) RoutingRules(ctx context.Context, source model.SourceID) (map[string][]model.RouteRule, error) {
	f.calls++
	if f.err != nil {
		return nil, f.err
	}
	if source != f.source {
		return map[string][]model.RouteRule{}, nil
	}
	return f.rules, nil
}

var _ outbound.RoutingPort = (*fakeRouting)(nil)

// fakeColumns trägt den `ColumnExclusionPort` als Fake: Spaltenexistenz je
// Adresse `schema.table.column`, Ausschlussstand je Tabelle.
type fakeColumns struct {
	existing map[string]bool
	excluded map[string][]string

	existsErr   error
	excludedErr error

	existsCalls   int
	excludedCalls int
}

func (f *fakeColumns) ColumnExists(ctx context.Context, schema, table, column string) (bool, error) {
	f.existsCalls++
	if f.existsErr != nil {
		return false, f.existsErr
	}
	return f.existing[schema+"."+table+"."+column], nil
}

func (f *fakeColumns) ExcludedColumns(ctx context.Context, source model.SourceID) (map[string][]string, error) {
	f.excludedCalls++
	if f.excludedErr != nil {
		return nil, f.excludedErr
	}
	return f.excluded, nil
}

var _ outbound.ColumnExclusionPort = (*fakeColumns)(nil)

func route(t *testing.T, name, target string, order int64, column, equals string) model.RouteRule {
	t.Helper()
	var when *model.RouteCondition
	if column != "" {
		when = &model.RouteCondition{Column: column, Equals: equals}
	}
	built, err := model.NewRouteRule(name, target, order, when)
	if err != nil {
		t.Fatalf("NewRouteRule(%q): %v", name, err)
	}
	return built
}

func ruleSpec(target string, order int, column, equals string) string {
	spec := `{"target": "` + target + `", "order": ` + strconv.Itoa(order)
	if column != "" {
		spec += `, "when": {"column": "` + column + `", "equals": "` + equals + `"}`
	}
	return spec + "}"
}

// newPorts baut die Fakes. `public.orders` trägt `eu_orders` (order 10,
// region=eu), `us_orders` (order 20, region=us) und die Abschlussregel `rest`
// (order 100, ohne when); `public.nowhen` trägt zwei Regeln mit when und keine
// Abschlussregel; `public.other` trägt `asia_rule` (order 15) und `status_rule`
// (order 7). Spalten von `public.orders`: `id`, `name`, `region`, `status`,
// `secret` (ausgeschlossen).
func newPorts(t *testing.T) (*fakeRouting, *fakeColumns) {
	routing := &fakeRouting{
		source: "src-1",
		rules: map[string][]model.RouteRule{
			"public.orders": {
				route(t, "eu_orders", "eu", 10, "region", "eu"),
				route(t, "us_orders", "us", 20, "region", "us"),
				route(t, "rest", "sonstige", 100, "", ""),
			},
			"public.nowhen": {
				route(t, "first", "a", 10, "region", "a"),
				route(t, "second", "b", 20, "region", "b"),
			},
			"public.other": {
				route(t, "asia_rule", "asia", 15, "region", "asia"),
				route(t, "status_rule", "st", 7, "status", "x"),
			},
		},
	}
	columns := &fakeColumns{
		existing: map[string]bool{
			"public.orders.id": true, "public.orders.name": true, "public.orders.region": true,
			"public.orders.status": true, "public.orders.secret": true,
			"public.nowhen.region": true, "public.other.region": true, "public.other.status": true,
		},
		excluded: map[string][]string{"public.orders": {"secret"}},
	}
	return routing, columns
}

func set(routing *fakeRouting, columns *fakeColumns, source model.SourceID, table, ruleName, spec string) (model.RouteRule, error) {
	return setroute.NewSetRouteService(routing, columns).Set(context.Background(), setroute.SetRouteCommand{
		Source: source, Schema: "public", Table: table, RuleName: ruleName, RuleSpec: spec,
	})
}

// TestSetRouteAcceptsAConflictFreeRule trägt den Happy Path
// (`LH-FA-CFG-008`): eine Regel ohne Verstoß gegen R1 bis R5 wird als Regel mit
// Name, Ziel, `order` und Bedingung geliefert; der Use Case schreibt den
// Regelstand nicht. Die Gegenproben tragen die Grenzfälle der Konfliktprüfung:
// dasselbe Paar-Teil in einer anderen Spalte, die leere Zeichenkette als
// `equals`, eine gleichnamige Regel und eine gleiche `order` in einer anderen
// Tabelle, ein Regelstand einer anderen Quelle.
func TestSetRouteAcceptsAConflictFreeRule(t *testing.T) {
	for _, tc := range []struct {
		label    string
		source   model.SourceID
		table    string
		ruleName string
		spec     string
	}{
		{"Regel mit when", "src-1", "orders", "asia_orders", ruleSpec("asia", 15, "region", "asia")},
		{"gleiche equals in anderer Spalte (R5 bindet das Paar)", "src-1", "orders", "status_eu", ruleSpec("st", 16, "status", "eu")},
		{"leere Zeichenkette als equals", "src-1", "orders", "region_empty", ruleSpec("leer", 17, "region", "")},
		{"gleicher Name und gleiche order in anderer Tabelle (R1/R2 je Tabelle)", "src-1", "other", "eu_orders", ruleSpec("eu", 10, "region", "eu")},
		{"Abschlussregel mit höchster order (R4)", "src-1", "nowhen", "closing", ruleSpec("sonstige", 30, "", "")},
		{"Regelstand einer anderen Quelle", "src-2", "orders", "eu_orders", ruleSpec("eu", 10, "region", "eu")},
	} {
		routing, columns := newPorts(t)
		before := len(routing.rules["public."+tc.table])
		got, err := set(routing, columns, tc.source, tc.table, tc.ruleName, tc.spec)
		if err != nil {
			t.Fatalf("%s: Set = %v, wollen nil", tc.label, err)
		}
		if got.Name() != tc.ruleName {
			t.Fatalf("%s: Regelname = %q, wollen %q", tc.label, got.Name(), tc.ruleName)
		}
		if len(routing.rules["public."+tc.table]) != before {
			t.Fatalf("%s: der Use Case hat den Regelstand verändert", tc.label)
		}
	}
	routing, columns := newPorts(t)
	got, err := set(routing, columns, "src-1", "orders", "asia_orders", ruleSpec("asia", 15, "region", "asia"))
	if err != nil {
		t.Fatal(err)
	}
	condition, conditional := got.Condition()
	if got.Target() != "asia" || got.Order() != 15 || !conditional || condition.Column != "region" || condition.Equals != "asia" {
		t.Fatalf("Regel = %+v / %+v, wollen asia, 15, region=asia", got, condition)
	}
}

// TestSetRouteRejectsWithTheSpecTexts trägt je Zeile der Fehlertext-Tabelle
// (`SPEC-019`) den Auslöser, den exakten Fehlertext (Klartext, Doppelpunkt,
// Adresse) und den Grund; die Prüfreihenfolge ist die der Spec — bei zwei
// verletzten Zeilen bestimmt die erste den Text. Die Eingabe jedes Falls ist
// die Verletzung selbst. Rot färbende Mutation je Zeile: die zugehörige
// Prüfung im Use Case bzw. in der Domäne entfernen oder ihre Eingabe verändern
// (R1: `rule.name` gegen `r.name` vertauschen; R2: `rule.order` gegen
// `r.order`; R3 fehlt: `condition.Column` gegen `command.RuleName` in
// `ColumnExists`; R3 ausgeschlossen: `excluded[table]` gegen
// `excluded[command.Table]`; R4/R5: die Vergleichsfelder austauschen) — der Fall
// endet mit einem anderen Text oder angenommen.
func TestSetRouteRejectsWithTheSpecTexts(t *testing.T) {
	long := strings.Repeat("a", 64)
	for _, tc := range []struct {
		label    string
		table    string
		ruleName string
		spec     string
		text     string
		reason   error
	}{
		// Formzeilen, in der Reihenfolge der Tabelle.
		{"Regelname leer", "orders", "", ruleSpec("x", 5, "", ""), "Regelname ist ungültig: public.orders.", domainerrors.ErrInvalidRuleName},
		{"Regelname mit Großbuchstaben", "orders", "Eu", ruleSpec("x", 5, "", ""), "Regelname ist ungültig: public.orders.Eu", domainerrors.ErrInvalidRuleName},
		{"Regelname mit Punkt (Adresse zerlegbar)", "orders", "a.b", ruleSpec("x", 5, "", ""), "Regelname ist ungültig: public.orders.a.b", domainerrors.ErrInvalidRuleName},
		{"rule_spec SQL-NULL (leerer Text)", "orders", "neu", "", "rule_spec ist ungültig: public.orders.neu", domainerrors.ErrInvalidRuleSpec},
		{"rule_spec JSON-null", "orders", "neu", "null", "rule_spec ist ungültig: public.orders.neu", domainerrors.ErrInvalidRuleSpec},
		{"rule_spec ohne Objekt", "orders", "neu", `[1]`, "rule_spec ist ungültig: public.orders.neu", domainerrors.ErrInvalidRuleSpec},
		{"rule_spec eine Zeichenkette", "orders", "neu", `"x"`, "rule_spec ist ungültig: public.orders.neu", domainerrors.ErrInvalidRuleSpec},
		{"unbekannter Schlüssel", "orders", "neu", `{"target": "x", "order": 5, "zzz": 1}`, "unbekannter Schlüssel in rule_spec: zzz", domainerrors.ErrUnknownRuleSpecKey},
		{"unbekannter Schlüssel in when", "orders", "neu", `{"target": "x", "order": 5, "when": {"column": "region", "equals": "q", "extra": 1}}`, "unbekannter Schlüssel in rule_spec: extra", domainerrors.ErrUnknownRuleSpecKey},
		{"zwei unbekannte Schlüssel: der erste in aufsteigender Ordnung", "orders", "neu", `{"zzz": 1, "target": "x", "order": 5, "aaa": 1}`, "unbekannter Schlüssel in rule_spec: aaa", domainerrors.ErrUnknownRuleSpecKey},
		{"unbekannter Schlüssel geht der Form voraus", "orders", "neu", `{"zzz": 1}`, "unbekannter Schlüssel in rule_spec: zzz", domainerrors.ErrUnknownRuleSpecKey},
		{"target fehlt", "orders", "neu", `{"order": 5}`, "rule_spec ist ungültig: public.orders.neu", domainerrors.ErrInvalidRuleSpec},
		{"target keine Zeichenkette", "orders", "neu", `{"target": 7, "order": 5}`, "rule_spec ist ungültig: public.orders.neu", domainerrors.ErrInvalidRuleSpec},
		{"order fehlt", "orders", "neu", `{"target": "x"}`, "rule_spec ist ungültig: public.orders.neu", domainerrors.ErrInvalidRuleSpec},
		{"order eine Zeichenkette", "orders", "neu", `{"target": "x", "order": "5"}`, "rule_spec ist ungültig: public.orders.neu", domainerrors.ErrInvalidRuleSpec},
		{"order null", "orders", "neu", `{"target": "x", "order": null}`, "rule_spec ist ungültig: public.orders.neu", domainerrors.ErrInvalidRuleSpec},
		{"order 0", "orders", "neu", `{"target": "x", "order": 0}`, "rule_spec ist ungültig: public.orders.neu", domainerrors.ErrInvalidRuleSpec},
		{"order negativ", "orders", "neu", `{"target": "x", "order": -3}`, "rule_spec ist ungültig: public.orders.neu", domainerrors.ErrInvalidRuleSpec},
		{"order mit Bruch", "orders", "neu", `{"target": "x", "order": 5.5}`, "rule_spec ist ungültig: public.orders.neu", domainerrors.ErrInvalidRuleSpec},
		{"order ganzzahlig mit Nachkommastelle", "orders", "neu", `{"target": "x", "order": 5.0}`, "rule_spec ist ungültig: public.orders.neu", domainerrors.ErrInvalidRuleSpec},
		{"order mit Exponent", "orders", "neu", `{"target": "x", "order": 1e1}`, "rule_spec ist ungültig: public.orders.neu", domainerrors.ErrInvalidRuleSpec},
		{"order über der Obergrenze", "orders", "neu", `{"target": "x", "order": 2147483648}`, "rule_spec ist ungültig: public.orders.neu", domainerrors.ErrInvalidRuleSpec},
		{"order über int64", "orders", "neu", `{"target": "x", "order": 99999999999999999999}`, "rule_spec ist ungültig: public.orders.neu", domainerrors.ErrInvalidRuleSpec},
		{"when keine Objekt", "orders", "neu", `{"target": "x", "order": 5, "when": "region"}`, "rule_spec ist ungültig: public.orders.neu", domainerrors.ErrInvalidRuleSpec},
		{"when null", "orders", "neu", `{"target": "x", "order": 5, "when": null}`, "rule_spec ist ungültig: public.orders.neu", domainerrors.ErrInvalidRuleSpec},
		{"when leer", "orders", "neu", `{"target": "x", "order": 5, "when": {}}`, "rule_spec ist ungültig: public.orders.neu", domainerrors.ErrInvalidRuleSpec},
		{"when ohne equals", "orders", "neu", `{"target": "x", "order": 5, "when": {"column": "region"}}`, "rule_spec ist ungültig: public.orders.neu", domainerrors.ErrInvalidRuleSpec},
		{"when ohne column", "orders", "neu", `{"target": "x", "order": 5, "when": {"equals": "eu"}}`, "rule_spec ist ungültig: public.orders.neu", domainerrors.ErrInvalidRuleSpec},
		{"when.column leer", "orders", "neu", `{"target": "x", "order": 5, "when": {"column": "", "equals": "eu"}}`, "rule_spec ist ungültig: public.orders.neu", domainerrors.ErrInvalidRuleSpec},
		{"when.column mit U+0000", "orders", "neu", `{"target": "x", "order": 5, "when": {"column": "re\u0000gion", "equals": "eu"}}`, "rule_spec ist ungültig: public.orders.neu", domainerrors.ErrInvalidRuleSpec},
		{"when.column keine Zeichenkette", "orders", "neu", `{"target": "x", "order": 5, "when": {"column": 1, "equals": "eu"}}`, "rule_spec ist ungültig: public.orders.neu", domainerrors.ErrInvalidRuleSpec},
		{"when.equals keine Zeichenkette", "orders", "neu", `{"target": "x", "order": 5, "when": {"column": "region", "equals": 1}}`, "rule_spec ist ungültig: public.orders.neu", domainerrors.ErrInvalidRuleSpec},
		{"Zielname in Großbuchstaben", "orders", "neu", ruleSpec("EU", 5, "", ""), "Zielname ist ungültig: public.orders.EU", domainerrors.ErrInvalidRouteTarget},
		{"Zielname leer", "orders", "neu", ruleSpec("", 5, "", ""), "Zielname ist ungültig: public.orders.", domainerrors.ErrInvalidRouteTarget},
		{"Zielname mit Punkt (Adresse zerlegbar)", "orders", "neu", ruleSpec("a.b", 5, "", ""), "Zielname ist ungültig: public.orders.a.b", domainerrors.ErrInvalidRouteTarget},
		{"Zielname mit 64 Zeichen", "orders", "neu", ruleSpec(long, 5, "", ""), "Zielname ist ungültig: public.orders." + long, domainerrors.ErrInvalidRouteTarget},
		{"Zielname mit führendem Unterstrich", "orders", "neu", ruleSpec("_x", 5, "", ""), "Zielname ist ungültig: public.orders._x", domainerrors.ErrInvalidRouteTarget},
		{"die Formzeile geht dem Alphabet voraus", "orders", "neu", `{"target": "EU", "order": 0}`, "rule_spec ist ungültig: public.orders.neu", domainerrors.ErrInvalidRuleSpec},
		// Konfliktfreiheit.
		{"R1 Regelname vergeben", "orders", "eu_orders", ruleSpec("x", 5, "", ""), "Regelname bereits vergeben: public.orders.eu_orders", domainerrors.ErrRuleNameTaken},
		{"R2 order vergeben", "orders", "neu", ruleSpec("x", 10, "region", "q"), "order bereits vergeben: public.orders.10", domainerrors.ErrRouteOrderTaken},
		{"R1 geht R2 voraus", "orders", "eu_orders", ruleSpec("x", 10, "region", "q"), "Regelname bereits vergeben: public.orders.eu_orders", domainerrors.ErrRuleNameTaken},
		{"R3 Spalte fehlt an der Quelle", "orders", "neu", ruleSpec("x", 5, "nope", "q"), "Spalte existiert nicht an der Quelle: public.orders.nope", inbound.ErrSourceColumnMissing},
		{"R3 Spalte ist ausgeschlossen", "orders", "neu", ruleSpec("x", 5, "secret", "q"), "Spalte ist ausgeschlossen: public.orders.secret", domainerrors.ErrRoutingColumnExcluded},
		{"R3 Groß-/Kleinschreibung wird nicht gefaltet", "orders", "neu", ruleSpec("x", 5, "Region", "q"), "Spalte existiert nicht an der Quelle: public.orders.Region", inbound.ErrSourceColumnMissing},
		{"R2 geht R3 voraus", "orders", "neu", ruleSpec("x", 10, "nope", "q"), "order bereits vergeben: public.orders.10", domainerrors.ErrRouteOrderTaken},
		{"R3 geht R4 voraus", "orders", "neu", ruleSpec("x", 500, "nope", "q"), "Spalte existiert nicht an der Quelle: public.orders.nope", inbound.ErrSourceColumnMissing},
		{"R4 Abschlussregel bereits geführt", "orders", "neu", ruleSpec("x", 200, "", ""), "Regel ohne when bereits vorhanden: public.orders.rest", domainerrors.ErrRouteWithoutWhenTaken},
		{"R4 Abschlussregel bereits geführt, auch mit kleinerer order", "orders", "neu", ruleSpec("x", 5, "", ""), "Regel ohne when bereits vorhanden: public.orders.rest", domainerrors.ErrRouteWithoutWhenTaken},
		{"R4 Abschlussregel nicht mit höchster order", "nowhen", "neu", ruleSpec("x", 15, "", ""), "Regel ohne when trägt nicht die höchste order: public.nowhen.neu", domainerrors.ErrRouteWithoutWhenNotLast},
		{"R4 Regel mit when hinter der Abschlussregel", "orders", "neu", ruleSpec("x", 150, "region", "q"), "order liegt hinter der Regel ohne when: public.orders.rest", domainerrors.ErrRouteBehindWithoutWhen},
		{"R5 Paar bereits vergeben", "orders", "neu", ruleSpec("x", 15, "region", "eu"), "Bedingung bereits vergeben: public.orders.region", domainerrors.ErrRouteConditionTaken},
		{"R4 geht R5 voraus", "orders", "neu", ruleSpec("x", 150, "region", "eu"), "order liegt hinter der Regel ohne when: public.orders.rest", domainerrors.ErrRouteBehindWithoutWhen},
	} {
		routing, columns := newPorts(t)
		_, err := set(routing, columns, "src-1", tc.table, tc.ruleName, tc.spec)
		if err == nil {
			t.Fatalf("%s: Set = nil, wollen %q", tc.label, tc.text)
		}
		if err.Error() != tc.text {
			t.Fatalf("%s: Fehlertext = %q, wollen %q", tc.label, err.Error(), tc.text)
		}
		if !stderrors.Is(err, tc.reason) {
			t.Fatalf("%s: Fehler %v löst nicht auf den Grund %v auf", tc.label, err, tc.reason)
		}
	}
}

// TestSetRouteReadsNoPortForAFormViolation trägt: die Prüfungen der Regelform
// laufen vor dem ersten Lesezugriff; ein Antrag ohne Bedingung liest die
// Spaltenexistenz und den Ausschlussstand nicht. Rot färbende Mutation: den
// Lesezugriff vor `ParseRouteSpec` ziehen bzw. `ColumnExists` für eine Regel
// ohne `when` aufrufen — die Zähler färben rot.
func TestSetRouteReadsNoPortForAFormViolation(t *testing.T) {
	routing, columns := newPorts(t)
	if _, err := set(routing, columns, "src-1", "orders", "neu", `{"target": "x"}`); err == nil {
		t.Fatal("Set = nil, wollen eine Formverletzung")
	}
	if routing.calls != 0 || columns.existsCalls != 0 || columns.excludedCalls != 0 {
		t.Fatalf("Port-Aufrufe %d/%d/%d, wollen keinen vor der Formprüfung", routing.calls, columns.existsCalls, columns.excludedCalls)
	}
	if _, err := set(routing, columns, "src-1", "nowhen", "closing", ruleSpec("sonstige", 30, "", "")); err != nil {
		t.Fatal(err)
	}
	if columns.existsCalls != 0 || columns.excludedCalls != 0 {
		t.Fatalf("Spaltenzugriffe %d/%d für eine Regel ohne when, wollen keinen", columns.existsCalls, columns.excludedCalls)
	}
}

// TestSetRouteChecksTheColumnOfTheCommand trägt: R3 prüft genau die Spalte und
// Tabelle des Antrags gegen Katalog und Ausschlussstand dieser Tabelle — eine
// Spalte, die nur an einer Nachbartabelle ausgeschlossen ist, wird nicht
// abgelehnt. Rot färbende Mutation: den Ausschlussstand einer anderen Tabelle
// lesen (`excluded["public.other"]`) — die Gegenprobe färbt rot.
func TestSetRouteChecksTheColumnOfTheCommand(t *testing.T) {
	routing, columns := newPorts(t)
	columns.excluded["public.other"] = []string{"region"}
	if _, err := set(routing, columns, "src-1", "orders", "asia_orders", ruleSpec("asia", 16, "region", "asia")); err != nil {
		t.Fatalf("Set auf public.orders = %v, wollen nil (region ist nur an public.other ausgeschlossen)", err)
	}
	if _, err := set(routing, columns, "src-1", "other", "other_rule", ruleSpec("x", 40, "region", "q")); err == nil || err.Error() != "Spalte ist ausgeschlossen: public.other.region" {
		t.Fatalf("Set auf public.other = %v, wollen die Ablehnung der ausgeschlossenen Spalte", err)
	}
}

// TestSetRouteReturnsPortErrors trägt: ein Fehler eines Ports endet als
// derselbe Fehler, nicht als Spec-Fehlertext. Rot färbende Mutation: den
// Fehler eines Ports verwerfen — der Fall endet angenommen.
func TestSetRouteReturnsPortErrors(t *testing.T) {
	wantErr := stderrors.New("Katalog nicht erreichbar")
	routing, columns := newPorts(t)
	routing.err = wantErr
	if _, err := set(routing, columns, "src-1", "orders", "neu", ruleSpec("x", 5, "region", "q")); !stderrors.Is(err, wantErr) {
		t.Fatalf("Fehler des Routing-Ports = %v, wollen %v", err, wantErr)
	}
	routing, columns = newPorts(t)
	columns.existsErr = wantErr
	if _, err := set(routing, columns, "src-1", "orders", "neu", ruleSpec("x", 5, "region", "q")); !stderrors.Is(err, wantErr) {
		t.Fatalf("Fehler von ColumnExists = %v, wollen %v", err, wantErr)
	}
	routing, columns = newPorts(t)
	columns.excludedErr = wantErr
	if _, err := set(routing, columns, "src-1", "orders", "neu", ruleSpec("x", 5, "region", "q")); !stderrors.Is(err, wantErr) {
		t.Fatalf("Fehler von ExcludedColumns = %v, wollen %v", err, wantErr)
	}
}
