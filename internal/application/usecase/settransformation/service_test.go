package settransformation_test

import (
	"context"
	stderrors "errors"
	"testing"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/settransformation"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// fakeTransformationPort trägt den `TransformationPort` als Fake (`ADR-0030`):
// Regelstand und Spaltenliste sind an Quelle und Tabelle gebunden, die der
// Aufruf nennt — eine Regel oder Spalte einer anderen Quelle bzw. Tabelle
// erreicht den Use Case nicht.
type fakeTransformationPort struct {
	source  model.SourceID
	rules   map[string][]model.Transformation // schema.table → Regelstand
	columns map[string][]string               // schema.table → Spalten

	rulesErr   error
	columnsErr error

	ruleCalls   int
	columnCalls int
}

func (f *fakeTransformationPort) TransformationRules(ctx context.Context, source model.SourceID) (map[string][]model.Transformation, error) {
	f.ruleCalls++
	if f.rulesErr != nil {
		return nil, f.rulesErr
	}
	if source != f.source {
		return map[string][]model.Transformation{}, nil
	}
	return f.rules, nil
}

func (f *fakeTransformationPort) SourceColumns(ctx context.Context, schema, table string) ([]string, error) {
	f.columnCalls++
	if f.columnsErr != nil {
		return nil, f.columnsErr
	}
	return f.columns[schema+"."+table], nil
}

var _ outbound.TransformationPort = (*fakeTransformationPort)(nil)

func rule(t *testing.T, name, column, to string) model.Transformation {
	t.Helper()
	built, err := model.NewRenameColumn(name, column, to)
	if err != nil {
		t.Fatalf("NewRenameColumn(%q, %q, %q): %v", name, column, to, err)
	}
	return built
}

func spec(column, to string) string {
	return `{"kind": "rename_column", "column": "` + column + `", "to": "` + to + `"}`
}

// newPort baut den Fake mit der Tabelle `public.orders` (Spalten `id`, `name`,
// `status`, `secret`) und einer Regel `kundenname` (`name` → `customer_name`);
// die Nachbartabelle `public.other` trägt eine gleichnamige Regel und andere
// Spalten, die Nachbarquelle nichts.
func newPort(t *testing.T) *fakeTransformationPort {
	return &fakeTransformationPort{
		source: "src-1",
		rules: map[string][]model.Transformation{
			"public.orders": {rule(t, "kundenname", "name", "customer_name")},
			"public.other":  {rule(t, "statusname", "status", "state")},
		},
		columns: map[string][]string{
			"public.orders": {"id", "name", "status", "secret"},
			"public.other":  {"id", "name", "status", "state"},
		},
	}
}

func set(port *fakeTransformationPort, source model.SourceID, ruleName, ruleSpec string) (model.Transformation, error) {
	return settransformation.NewSetTransformationService(port).Set(context.Background(), settransformation.SetTransformationCommand{
		Source: source, Schema: "public", Table: "orders", RuleName: ruleName, RuleSpec: ruleSpec,
	})
}

// TestSetTransformationAcceptsAConflictFreeRule trägt den Happy Path
// (`LH-FA-CFG-007`): eine Regel ohne Verstoß gegen K1 bis K4 wird als Regel
// mit Name, Typ, Spalte und Zielname geliefert. Der Use Case trägt keinen
// Schreibpfad; der Regelstand der Tabelle bleibt, wie er war.
func TestSetTransformationAcceptsAConflictFreeRule(t *testing.T) {
	port := newPort(t)
	got, err := set(port, "src-1", "statusname", spec("status", "state"))
	if err != nil {
		t.Fatalf("Set = %v, wollen nil", err)
	}
	if got.Name() != "statusname" || got.Kind() != model.TransformationRenameColumn || got.Column() != "status" || got.To() != "state" {
		t.Fatalf("Regel = %+v, wollen statusname rename_column status -> state", got)
	}
	if len(port.rules["public.orders"]) != 1 {
		t.Fatalf("Regelstand der Tabelle = %v, wollen unverändert eine Regel", port.rules["public.orders"])
	}
}

// TestSetTransformationRejectsWithTheSpecTexts trägt je Zeile der
// Fehlertext-Tabelle (`SPEC-019`) den Auslöser, den exakten Fehlertext
// (Klartext, Doppelpunkt, Adresse) und den Grund; die Prüfreihenfolge ist die
// der Spec — bei zwei verletzten Zeilen bestimmt die erste den Text. Rot
// färbende Mutation je Zeile: die zugehörige Prüfung im Use Case bzw. in der
// Domäne entfernen oder ihre Eingabe verändern (K1: `rule.Name()` gegen
// `command.RuleName` vertauschen; K2: `spec.Column()` gegen `command.RuleName`;
// K3: `spec.Target()` gegen `spec.Column()`; K4: `spec.Column()` gegen
// `spec.Target()` in der Katalogprüfung) — der Fall endet mit einem anderen
// Text oder angenommen.
func TestSetTransformationRejectsWithTheSpecTexts(t *testing.T) {
	for _, tc := range []struct {
		label    string
		ruleName string
		ruleSpec string
		text     string
		reason   error
	}{
		// Formzeilen, in der Reihenfolge der Tabelle.
		{"Regelname leer", "", spec("status", "state"), "Regelname ist ungültig: public.orders.", domainerrors.ErrInvalidRuleName},
		{"Regelname mit Großbuchstaben", "Status", spec("status", "state"), "Regelname ist ungültig: public.orders.Status", domainerrors.ErrInvalidRuleName},
		{"Regelname mit Punkt (Adresse zerlegbar)", "a.b", spec("status", "state"), "Regelname ist ungültig: public.orders.a.b", domainerrors.ErrInvalidRuleName},
		{"rule_spec SQL-NULL (leerer Text)", "statusname", "", "rule_spec ist ungültig: public.orders.statusname", domainerrors.ErrInvalidRuleSpec},
		{"rule_spec JSON-null", "statusname", "null", "rule_spec ist ungültig: public.orders.statusname", domainerrors.ErrInvalidRuleSpec},
		{"rule_spec ohne Objekt", "statusname", `[1]`, "rule_spec ist ungültig: public.orders.statusname", domainerrors.ErrInvalidRuleSpec},
		{"kind fehlt", "statusname", `{"column": "status", "to": "state"}`, "rule_spec ist ungültig: public.orders.statusname", domainerrors.ErrInvalidRuleSpec},
		{"kind keine Zeichenkette", "statusname", `{"kind": 1}`, "rule_spec ist ungültig: public.orders.statusname", domainerrors.ErrInvalidRuleSpec},
		{"unbekannter kind", "statusname", `{"kind": "explode"}`, "unbekannter Regeltyp: explode", domainerrors.ErrUnknownTransformationKind},
		{"unbekannter Schlüssel", "statusname", `{"kind": "rename_column", "column": "status", "to": "state", "extra": 1}`, "unbekannter Schlüssel in rule_spec: extra", domainerrors.ErrUnknownRuleSpecKey},
		{"Pflichtschlüssel fehlt", "statusname", `{"kind": "rename_column", "column": "status"}`, "rule_spec ist ungültig: public.orders.statusname", domainerrors.ErrInvalidRuleSpec},
		{"Pflichtschlüssel mit falschem Typ", "statusname", `{"kind": "rename_column", "column": "status", "to": 7}`, "rule_spec ist ungültig: public.orders.statusname", domainerrors.ErrInvalidRuleSpec},
		{"Zielname leer (Form von to)", "statusname", spec("status", ""), "rule_spec ist ungültig: public.orders.statusname", domainerrors.ErrInvalidRuleSpec},
		// Konfliktfreiheit.
		{"K1 Regelname vergeben", "kundenname", spec("status", "state"), "Regelname bereits vergeben: public.orders.kundenname", domainerrors.ErrRuleNameTaken},
		{"K2 Spalte trägt eine Regel", "zweite", spec("name", "state"), "Spalte trägt bereits eine Regel: public.orders.name", domainerrors.ErrColumnHasRule},
		{"K3 Zielname gleicht dem Ziel einer anderen Regel", "zweite", spec("status", "customer_name"), "Zielname kollidiert mit einer anderen Regel: public.orders.customer_name", domainerrors.ErrTargetCollidesWithRule},
		{"K3 Zielname gleicht einer Spalte", "zweite", spec("status", "id"), "Zielname kollidiert mit einer Spalte der Tabelle: public.orders.id", domainerrors.ErrTargetCollidesWithColumn},
		{"K3 Zielname gleicht einer ausgeschlossenen Spalte", "zweite", spec("status", "secret"), "Zielname kollidiert mit einer Spalte der Tabelle: public.orders.secret", domainerrors.ErrTargetCollidesWithColumn},
		{"K3 Zielname gleicht der Quellspalte (Sentinel der Domäne)", "zweite", spec("status", "status"), "Zielname kollidiert mit einer Spalte der Tabelle: public.orders.status", domainerrors.ErrTargetCollidesWithColumn},
		{"K4 Spalte fehlt an der Quelle", "zweite", spec("gibt_es_nicht", "state"), "Spalte existiert nicht an der Quelle: public.orders.gibt_es_nicht", inbound.ErrSourceColumnMissing},
		// Prüfreihenfolge: die erste verletzte Zeile bestimmt den Text.
		{"Regelname vor rule_spec", "", "", "Regelname ist ungültig: public.orders.", domainerrors.ErrInvalidRuleName},
		{"Form von rule_spec vor kind", "statusname", `{"kind": "explode", "z": 1}` + "x", "rule_spec ist ungültig: public.orders.statusname", domainerrors.ErrInvalidRuleSpec},
		{"kind vor Schlüssel", "statusname", `{"kind": "explode", "extra": 1}`, "unbekannter Regeltyp: explode", domainerrors.ErrUnknownTransformationKind},
		{"Schlüssel vor Pflichtschlüssel", "statusname", `{"kind": "rename_column", "column": "status", "extra": 1}`, "unbekannter Schlüssel in rule_spec: extra", domainerrors.ErrUnknownRuleSpecKey},
		{"Form vor K1", "kundenname", spec("name", ""), "rule_spec ist ungültig: public.orders.kundenname", domainerrors.ErrInvalidRuleSpec},
		{"K1 vor K2", "kundenname", spec("name", "state"), "Regelname bereits vergeben: public.orders.kundenname", domainerrors.ErrRuleNameTaken},
		{"K2 vor K3", "zweite", spec("name", "id"), "Spalte trägt bereits eine Regel: public.orders.name", domainerrors.ErrColumnHasRule},
		{"K3 vor K4", "zweite", spec("gibt_es_nicht", "id"), "Zielname kollidiert mit einer Spalte der Tabelle: public.orders.id", domainerrors.ErrTargetCollidesWithColumn},
		{"K3 (Ziel gleich Quelle) vor K4 bei fehlender Spalte", "zweite", spec("gibt_es_nicht", "gibt_es_nicht"), "Zielname kollidiert mit einer Spalte der Tabelle: public.orders.gibt_es_nicht", domainerrors.ErrTargetCollidesWithColumn},
	} {
		port := newPort(t)
		_, err := set(port, "src-1", tc.ruleName, tc.ruleSpec)
		if err == nil {
			t.Errorf("%s: Set = nil, wollen %q", tc.label, tc.text)
			continue
		}
		if err.Error() != tc.text {
			t.Errorf("%s: Fehlertext = %q, wollen %q", tc.label, err.Error(), tc.text)
		}
		if !stderrors.Is(err, tc.reason) {
			t.Errorf("%s: Fehler = %v, wollen Grund %v", tc.label, err, tc.reason)
		}
		if len(port.rules["public.orders"]) != 1 {
			t.Errorf("%s: Regelstand = %v, wollen unverändert", tc.label, port.rules["public.orders"])
		}
	}
}

// TestSetTransformationChecksTheFormBeforeReadingTheStore trägt: eine
// Formverletzung endet, ohne dass der Use Case Regelstand oder Spaltenliste
// liest. Rot färbende Mutation: `CheckRuleName` bzw. `ParseTransformationSpec`
// hinter die beiden Port-Aufrufe stellen — die Zähler stehen dann bei 1.
func TestSetTransformationChecksTheFormBeforeReadingTheStore(t *testing.T) {
	for _, tc := range []struct{ name, spec string }{
		{"", spec("status", "state")},
		{"statusname", ""},
		{"statusname", `{"kind": "explode"}`},
	} {
		port := newPort(t)
		if _, err := set(port, "src-1", tc.name, tc.spec); err == nil {
			t.Fatalf("Set(%q, %q) = nil, wollen einen Fehler", tc.name, tc.spec)
		}
		if port.ruleCalls != 0 || port.columnCalls != 0 {
			t.Fatalf("Set(%q, %q): Port-Aufrufe = %d/%d, wollen 0/0", tc.name, tc.spec, port.ruleCalls, port.columnCalls)
		}
	}
}

// TestSetTransformationBindsTheTableAndSourceOfTheCommand trägt: die
// Prüfung liest Regelstand und Spalten der Tabelle und Quelle des Kommandos —
// der gleiche Regelname und ein gleicher Zielname an der Nachbartabelle bzw.
// in der Nachbarquelle sind kein Verstoß, eine Spalte nur der Nachbartabelle
// ist an dieser Tabelle nicht vorhanden. Rot färbende Mutation: `state[table]`
// durch `state["public.other"]` ersetzen (K1 färbt rot), in `Set` die
// Tabelle beim Lesen der Spalten vertauschen (K4 und K3 färben rot).
func TestSetTransformationBindsTheTableAndSourceOfTheCommand(t *testing.T) {
	port := newPort(t)
	// `statusname` heißt an `public.other` bereits eine Regel: an `orders` frei.
	if _, err := set(port, "src-1", "statusname", spec("status", "state")); err != nil {
		t.Fatalf("Regelname der Nachbartabelle: Set = %v, wollen nil", err)
	}
	// `state` ist eine Spalte nur von `public.other`: als Zielname an
	// `orders` frei (oben angenommen), als Quellspalte dort nicht vorhanden.
	if _, err := set(port, "src-1", "zweite", spec("state", "neu")); !stderrors.Is(err, inbound.ErrSourceColumnMissing) {
		t.Fatalf("Spalte der Nachbartabelle: Fehler = %v, wollen ErrSourceColumnMissing", err)
	}
	// Die Nachbarquelle trägt keinen Regelstand: derselbe Regelname ist dort frei.
	if _, err := set(port, "src-2", "kundenname", spec("status", "state")); err != nil {
		t.Fatalf("Regelname der Nachbarquelle: Set = %v, wollen nil", err)
	}
}

// TestSetTransformationComparesNamesCharacterExact trägt `SPEC-030`
// (Bezeichner): Großschreibung faltet nicht — `Name` ist eine andere Spalte
// als `name`, `ID` ein anderer Zielname als `id`. Rot färbende Mutation:
// die Vergleiche mit `strings.EqualFold` ersetzen (K3 lehnt `ID` ab, K4
// nimmt `Name` an).
func TestSetTransformationComparesNamesCharacterExact(t *testing.T) {
	port := newPort(t)
	if _, err := set(port, "src-1", "zweite", spec("status", "ID")); err != nil {
		t.Fatalf("Zielname ID neben Spalte id: Set = %v, wollen nil", err)
	}
	if _, err := set(port, "src-1", "dritte", spec("Name", "x")); !stderrors.Is(err, inbound.ErrSourceColumnMissing) {
		t.Fatalf("Spalte Name neben Spalte name: Fehler = %v, wollen ErrSourceColumnMissing", err)
	}
}

// TestSetTransformationReturnsPortErrors trägt: ein Lesefehler des Ports
// endet unverändert, ohne Regel.
func TestSetTransformationReturnsPortErrors(t *testing.T) {
	wantErr := stderrors.New("Antrags-Historie nicht lesbar")
	port := newPort(t)
	port.rulesErr = wantErr
	if _, err := set(port, "src-1", "statusname", spec("status", "state")); !stderrors.Is(err, wantErr) {
		t.Fatalf("Regelstand nicht lesbar: Fehler = %v, wollen %v", err, wantErr)
	}
	port = newPort(t)
	port.columnsErr = wantErr
	if _, err := set(port, "src-1", "statusname", spec("status", "state")); !stderrors.Is(err, wantErr) {
		t.Fatalf("Spaltenliste nicht lesbar: Fehler = %v, wollen %v", err, wantErr)
	}
}

// TestSetTransformationAcceptsAndChecksMapValueThroughTheDomain trägt den
// zweiten Regeltyp am unveränderten Use Case (`LH-FA-CFG-007`): ein
// `map_value`-Antrag durchläuft dieselben Prüfungen wie `rename_column` — die
// Formzeilen aus der Domäne, K1, K2 und K4 mit dem Text und der Adresse der
// Spec —, und K3 trifft ihn nicht: er trägt keinen Zielname, auch ein
// abgebildeter Wert, der einer Spalte gleicht, ist kein Konflikt. Rot färbende
// Mutation je Fall: der Use Case nennt einen Regeltyp beim Namen (etwa ein
// `if spec.Kind() == model.TransformationRenameColumn` um die Konflikt-Prüfung
// oder die Katalog-Prüfung) — der `map_value`-Fall endet angenommen bzw. mit
// einem anderen Text.
func TestSetTransformationAcceptsAndChecksMapValueThroughTheDomain(t *testing.T) {
	mapSpec := func(column, values string) string {
		return `{"kind": "map_value", "column": "` + column + `", "values": ` + values + `}`
	}
	port := newPort(t)
	got, err := set(port, "src-1", "statuswert", mapSpec("status", `{"o": "open", "c": "closed"}`))
	if err != nil {
		t.Fatalf("Set = %v, wollen nil", err)
	}
	if got.Name() != "statuswert" || got.Kind() != model.TransformationMapValue || got.Column() != "status" || got.To() != "" ||
		len(got.Values()) != 2 || got.Values()["o"] != "open" || got.Values()["c"] != "closed" {
		t.Fatalf("Regel = %+v, wollen statuswert map_value status mit zwei Zuordnungen", got)
	}
	if len(port.rules["public.orders"]) != 1 {
		t.Fatalf("Regelstand der Tabelle = %v, wollen unverändert eine Regel", port.rules["public.orders"])
	}
	if _, err := set(newPort(t), "src-1", "statuswert", mapSpec("status", `{"o": "id", "c": "name", "x": "status"}`)); err != nil {
		t.Fatalf("abgebildete Werte gleich Spaltennamen: Set = %v, wollen nil (K3 trifft map_value nicht)", err)
	}

	for _, tc := range []struct {
		label    string
		ruleName string
		ruleSpec string
		text     string
		reason   error
	}{
		{"K1 Regelname vergeben", "kundenname", mapSpec("status", `{"o": "open"}`), "Regelname bereits vergeben: public.orders.kundenname", domainerrors.ErrRuleNameTaken},
		{"K2 Spalte trägt eine Umbenennung", "statuswert", mapSpec("name", `{"a": "b"}`), "Spalte trägt bereits eine Regel: public.orders.name", domainerrors.ErrColumnHasRule},
		{"K4 Spalte fehlt an der Quelle", "statuswert", mapSpec("gibt_es_nicht", `{"a": "b"}`), "Spalte existiert nicht an der Quelle: public.orders.gibt_es_nicht", inbound.ErrSourceColumnMissing},
		{"Form: values leer", "statuswert", mapSpec("status", `{}`), "rule_spec ist ungültig: public.orders.statuswert", domainerrors.ErrInvalidRuleSpec},
		{"Form: Wert keine Zeichenkette", "statuswert", mapSpec("status", `{"o": 1}`), "rule_spec ist ungültig: public.orders.statuswert", domainerrors.ErrInvalidRuleSpec},
		{"Schlüssel: to gehört nicht zu map_value", "statuswert", `{"kind": "map_value", "column": "status", "to": "x", "values": {"o": "open"}}`, "unbekannter Schlüssel in rule_spec: to", domainerrors.ErrUnknownRuleSpecKey},
	} {
		port := newPort(t)
		_, err := set(port, "src-1", tc.ruleName, tc.ruleSpec)
		if !stderrors.Is(err, tc.reason) || err.Error() != tc.text {
			t.Fatalf("%s: Fehler = %v, wollen %q mit dem Grund %v", tc.label, err, tc.text, tc.reason)
		}
	}
}
