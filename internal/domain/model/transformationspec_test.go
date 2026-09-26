package model

import (
	stderrors "errors"
	"strings"
	"testing"

	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

// CheckRuleName prüft das Alphabet des Regelnamens (`SPEC-030`, Bezeichner):
// 1 bis 63 Zeichen aus `a`–`z`, `0`–`9`, `_`. Rot färbende Mutation: die
// Grenze `{1,63}` in `{1,64}` ändern — der Name mit 64 Zeichen wird
// angenommen; `A-Z` in das Alphabet aufnehmen — der Name mit Großbuchstaben
// wird angenommen.
func TestCheckRuleNameBindsTheAlphabet(t *testing.T) {
	for _, name := range []string{"a", "kundenname", "regel_1", "_", "0", strings.Repeat("a", 63)} {
		if err := CheckRuleName(name); err != nil {
			t.Fatalf("CheckRuleName(%q) = %v, wollen nil", name, err)
		}
	}
	for _, name := range []string{"", "Kundenname", "regel-1", "regel.1", "regel 1", "regel\n", "rëgel", strings.Repeat("a", 64)} {
		if err := CheckRuleName(name); !stderrors.Is(err, domainerrors.ErrInvalidRuleName) {
			t.Fatalf("CheckRuleName(%q) = %v, wollen ErrInvalidRuleName", name, err)
		}
	}
}

// ParseTransformationSpec liest die Regelform strikt und prüft in der
// Reihenfolge der Spec (`SPEC-019` Fehlertext-Tabelle, Formzeilen zwei bis
// fünf): jeder Fall bindet Auslöser, Grund und — wo die Adresse ihn nennt —
// den Wert. Rot färbende Mutation je Zeile: die zugehörige Prüfung entfernen
// (`ValidString`, die `kind`-Prüfung, `allowedRuleKeys`, die
// Schlüssel-Schleife, die Typprüfung von `column`/`to`, die Form über
// `NewRenameColumn`) — der Fall endet dann mit einem anderen Grund oder
// angenommen.
func TestParseTransformationSpecRejectsInSpecOrder(t *testing.T) {
	invalid := domainerrors.ErrInvalidRuleSpec
	for _, tc := range []struct {
		name   string
		text   string
		reason error
		detail string
	}{
		{"leerer Text (SQL-NULL)", "", invalid, ""},
		{"nur Leerraum", "  \n", invalid, ""},
		{"JSON-null", "null", invalid, ""},
		{"Array", `[{"kind": "rename_column"}]`, invalid, ""},
		{"Zahl", "1", invalid, ""},
		{"Zeichenkette", `"rename_column"`, invalid, ""},
		{"kein JSON", `{oops`, invalid, ""},
		{"Objekt mit Nachlauf", `{"kind": "rename_column"}x`, invalid, ""},
		{"ungültiges UTF-8", "{\"kind\": \"rename_column\", \"column\": \"a\xff\", \"to\": \"b\"}", invalid, ""},
		{"kind fehlt", `{"column": "name", "to": "kundenname"}`, invalid, ""},
		{"kind ist null", `{"kind": null, "column": "name", "to": "x"}`, invalid, ""},
		{"kind ist eine Zahl", `{"kind": 1, "column": "name", "to": "x"}`, invalid, ""},
		{"kind ist ein Objekt", `{"kind": {}}`, invalid, ""},
		{"unbekannter kind", `{"kind": "map_column", "column": "name"}`, domainerrors.ErrUnknownTransformationKind, "map_column"},
		{"leerer kind", `{"kind": ""}`, domainerrors.ErrUnknownTransformationKind, ""},
		{"unbekannter kind vor unbekanntem Schlüssel", `{"kind": "nope", "zzz": 1}`, domainerrors.ErrUnknownTransformationKind, "nope"},
		{"unbekannter Schlüssel", `{"kind": "rename_column", "column": "name", "to": "x", "extra": 1}`, domainerrors.ErrUnknownRuleSpecKey, "extra"},
		{"unbekannter Schlüssel vor fehlendem Pflichtschlüssel", `{"kind": "rename_column", "column": "name", "values": {}}`, domainerrors.ErrUnknownRuleSpecKey, "values"},
		{"erster unbekannter Schlüssel in aufsteigender Ordnung", `{"kind": "rename_column", "column": "name", "to": "x", "z": 1, "b": 2}`, domainerrors.ErrUnknownRuleSpecKey, "b"},
		{"Schlüssel mit anderer Schreibweise", `{"kind": "rename_column", "Column": "name", "to": "x"}`, domainerrors.ErrUnknownRuleSpecKey, "Column"},
		{"column fehlt", `{"kind": "rename_column", "to": "x"}`, invalid, ""},
		{"to fehlt", `{"kind": "rename_column", "column": "name"}`, invalid, ""},
		{"column ist eine Zahl", `{"kind": "rename_column", "column": 1, "to": "x"}`, invalid, ""},
		{"to ist null", `{"kind": "rename_column", "column": "name", "to": null}`, invalid, ""},
		{"column leer", `{"kind": "rename_column", "column": "", "to": "x"}`, invalid, ""},
		{"to leer", `{"kind": "rename_column", "column": "name", "to": ""}`, invalid, ""},
		{"to mit NUL", "{\"kind\": \"rename_column\", \"column\": \"name\", \"to\": \"a\\u0000b\"}", invalid, ""},
		{"column mit NUL", "{\"kind\": \"rename_column\", \"column\": \"a\\u0000b\", \"to\": \"x\"}", invalid, ""},
		{"to länger als 63 Byte", `{"kind": "rename_column", "column": "name", "to": "` + strings.Repeat("a", 64) + `"}`, invalid, ""},
	} {
		_, err := ParseTransformationSpec(tc.text)
		if !stderrors.Is(err, tc.reason) {
			t.Fatalf("%s: Fehler = %v, wollen %v", tc.name, err, tc.reason)
		}
		if tc.detail != "" || errorsAs(err) {
			var specErr *TransformationSpecError
			if !stderrors.As(err, &specErr) || specErr.Detail != tc.detail {
				t.Fatalf("%s: Fehler = %v, wollen den Wert %q", tc.name, err, tc.detail)
			}
		}
	}
}

// errorsAs meldet, ob der Fehler ein `TransformationSpecError` trägt: dessen
// Fälle nennen den Wert der Adresse, die übrigen nicht.
func errorsAs(err error) bool {
	var specErr *TransformationSpecError
	return stderrors.As(err, &specErr)
}

// Eine gültige Regelform liefert Typ, Spalte und Zielname; ein Zielname
// gleich der Quellspalte ist keine Formverletzung (K3, `SPEC-030` Randfälle),
// 63 Byte Zielname sind zulässig, ein Zielname mit Mehrbyte-Zeichen zählt in
// Byte. Rot färbende Mutation: den Ausschluss von
// `ErrTransformationTargetIsColumn` in `ParseTransformationSpec` streichen —
// der Zielname gleich der Quellspalte endet als `ErrInvalidRuleSpec`.
func TestParseTransformationSpecAcceptsValidForms(t *testing.T) {
	spec, err := ParseTransformationSpec(` {"kind": "rename_column", "column": "name", "to": "kundenname"} `)
	if err != nil {
		t.Fatalf("gültige Form: %v", err)
	}
	if spec.Kind() != TransformationRenameColumn || spec.Column() != "name" || spec.Target() != "kundenname" {
		t.Fatalf("Regelform = %+v, wollen rename_column name -> kundenname", spec)
	}
	same, err := ParseTransformationSpec(`{"kind": "rename_column", "column": "id", "to": "id"}`)
	if err != nil {
		t.Fatalf("Zielname gleich Quellspalte: Fehler = %v, wollen nil (K3 lehnt ab, nicht die Form)", err)
	}
	if same.Column() != "id" || same.Target() != "id" {
		t.Fatalf("Regelform = %+v, wollen id -> id", same)
	}
	if _, err := ParseTransformationSpec(`{"kind": "rename_column", "column": "name", "to": "` + strings.Repeat("a", 63) + `"}`); err != nil {
		t.Fatalf("Zielname mit 63 Byte: %v", err)
	}
	if _, err := ParseTransformationSpec(`{"kind": "rename_column", "column": "spalte ä.\"x\"", "to": "ü-ziel"}`); err != nil {
		t.Fatalf("Namen mit Punkt, Anführungszeichen und Mehrbyte-Zeichen: %v", err)
	}
	if _, err := ParseTransformationSpec(`{"kind": "rename_column", "column": "name", "to": "` + strings.Repeat("ü", 32) + `"}`); !stderrors.Is(err, domainerrors.ErrInvalidRuleSpec) {
		t.Fatalf("Zielname mit 64 Byte in 32 Zeichen: Fehler = %v, wollen ErrInvalidRuleSpec", err)
	}
}

// CheckConflicts prüft K1 bis K3 an ihrer Eingabe und in der Reihenfolge der
// Spec. Rot färbende Mutation je Zeile: die Prüfung streichen (K1: die
// Namensschleife, K2: die Spaltenschleife, K3a: die Zielschleife, K3b:
// `s.to == s.column` bzw. `containsName(columns, s.to)`) — der Fall endet
// mit `nil` bzw. dem Grund der Folgeprüfung.
func TestCheckConflictsBindsK1ToK3AndTheirOrder(t *testing.T) {
	mustRule := func(name, column, to string) Transformation {
		rule, err := NewRenameColumn(name, column, to)
		if err != nil {
			t.Fatalf("NewRenameColumn(%q, %q, %q): %v", name, column, to, err)
		}
		return rule
	}
	spec := func(column, to string) TransformationSpec {
		return TransformationSpec{kind: TransformationRenameColumn, column: column, to: to}
	}
	columns := []string{"id", "name", "status", "secret"}
	rules := []Transformation{mustRule("kundenname", "name", "customer_name")}

	for _, tc := range []struct {
		label string
		name  string
		spec  TransformationSpec
		rules []Transformation
		want  error
	}{
		{"frei", "statusname", spec("status", "state"), rules, nil},
		{"leerer Regelstand", "kundenname", spec("name", "customer_name"), nil, nil},
		{"K1 Name vergeben", "kundenname", spec("status", "state"), rules, domainerrors.ErrRuleNameTaken},
		{"K1 zeichengenau, Großschreibung ist ein anderer Name", "Kundenname", spec("status", "state"), rules, nil},
		{"K2 Spalte trägt eine Regel", "zweite", spec("name", "state"), rules, domainerrors.ErrColumnHasRule},
		{"K3 Ziel gleicht Ziel einer anderen Regel", "zweite", spec("status", "customer_name"), rules, domainerrors.ErrTargetCollidesWithRule},
		{"K3 Ziel gleicht Spalte der Tabelle", "zweite", spec("status", "id"), rules, domainerrors.ErrTargetCollidesWithColumn},
		{"K3 Ziel gleicht ausgeschlossener Spalte (Katalog)", "zweite", spec("status", "secret"), rules, domainerrors.ErrTargetCollidesWithColumn},
		{"K3 Ziel gleicht der Quellspalte, auch ohne Katalogeintrag", "zweite", spec("fehlt", "fehlt"), rules, domainerrors.ErrTargetCollidesWithColumn},
		{"K3 zeichengenau, Großschreibung ist ein anderer Name", "zweite", spec("status", "ID"), rules, nil},
		{"K1 vor K2", "kundenname", spec("name", "state"), rules, domainerrors.ErrRuleNameTaken},
		{"K2 vor K3", "zweite", spec("name", "id"), rules, domainerrors.ErrColumnHasRule},
		{"K3 Regel vor K3 Spalte (Ziel gleicht beidem)", "zweite", spec("status", "id"), append(append([]Transformation{}, rules...), mustRule("dritte", "secret", "id")), domainerrors.ErrTargetCollidesWithRule},
	} {
		got := tc.spec.CheckConflicts(tc.name, tc.rules, columns)
		if !stderrors.Is(got, tc.want) && !(tc.want == nil && got == nil) {
			t.Fatalf("%s: Fehler = %v, wollen %v", tc.label, got, tc.want)
		}
	}
}

// Build legt die Regel unter dem Namen an; die Quellspalte als Ziel bleibt
// ein Fehler des Konstruktors, den K3 vorher abfängt.
func TestTransformationSpecBuild(t *testing.T) {
	spec, err := ParseTransformationSpec(`{"kind": "rename_column", "column": "name", "to": "kundenname"}`)
	if err != nil {
		t.Fatalf("ParseTransformationSpec: %v", err)
	}
	rule, err := spec.Build("umbenennung")
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if rule.Name() != "umbenennung" || rule.Kind() != TransformationRenameColumn || rule.Column() != "name" || rule.To() != "kundenname" {
		t.Fatalf("Regel = %+v", rule)
	}
	same, err := ParseTransformationSpec(`{"kind": "rename_column", "column": "id", "to": "id"}`)
	if err != nil {
		t.Fatalf("ParseTransformationSpec: %v", err)
	}
	if _, err := same.Build("regel"); !stderrors.Is(err, domainerrors.ErrTransformationTargetIsColumn) {
		t.Fatalf("Build mit Ziel gleich Quelle: Fehler = %v, wollen ErrTransformationTargetIsColumn", err)
	}
	if _, err := (TransformationSpec{kind: "map_value"}).Build("regel"); !stderrors.Is(err, domainerrors.ErrUnknownTransformationKind) {
		t.Fatalf("Build mit unbekanntem Typ: Fehler = %v, wollen ErrUnknownTransformationKind", err)
	}
}

// FoldTransformations leitet den Regelstand aus den Zeilen in der übergebenen
// Ordnung ab: Set trägt ein, Remove nimmt heraus, ein Set-Remove-Set-Zyklus
// endet mit der letzten Regel, ein Set unter vorhandenem Namen ersetzt an der
// Stelle, ein Remove eines unbekannten Namens ändert nichts, eine andere Art
// trägt keinen Stand. Rot färbende Mutation: die Ordnung der Eingabe
// vertauschen (Remove vor Set) — der Zyklus-Fall liefert die Regel nicht
// bzw. eine gelöschte; `dropByName` gegen `nil`-Rückgabe ohne Filter ersetzen
// — der Remove-Fall behält die Regel.
func TestFoldTransformationsFollowsTheOrderOfTheRows(t *testing.T) {
	set := func(name, column, to string) TransformationRecord {
		return TransformationRecord{
			Kind: AdministrationRequestSetTransformation, Name: name,
			Spec: `{"kind": "rename_column", "column": "` + column + `", "to": "` + to + `"}`,
		}
	}
	remove := func(name string) TransformationRecord {
		return TransformationRecord{Kind: AdministrationRequestRemoveTransformation, Name: name}
	}
	names := func(rules []Transformation) string {
		out := make([]string, 0, len(rules))
		for _, rule := range rules {
			out = append(out, rule.Name()+":"+rule.Column()+">"+rule.To())
		}
		return strings.Join(out, ",")
	}
	for _, tc := range []struct {
		label   string
		records []TransformationRecord
		want    string
	}{
		{"keine Zeile", nil, ""},
		{"Set", []TransformationRecord{set("a", "name", "x")}, "a:name>x"},
		{"zwei Sets in Reihenfolge", []TransformationRecord{set("a", "name", "x"), set("b", "status", "y")}, "a:name>x,b:status>y"},
		{"Set, Remove", []TransformationRecord{set("a", "name", "x"), remove("a")}, ""},
		{"Set, Remove, Set (Zyklus)", []TransformationRecord{set("a", "name", "x"), remove("a"), set("a", "name", "z")}, "a:name>z"},
		{"Remove vor Set", []TransformationRecord{remove("a"), set("a", "name", "x")}, "a:name>x"},
		{"Set unter vorhandenem Namen ersetzt an der Stelle", []TransformationRecord{set("a", "name", "x"), set("b", "status", "y"), set("a", "id", "z")}, "a:id>z,b:status>y"},
		{"Remove eines unbekannten Namens", []TransformationRecord{set("a", "name", "x"), remove("b")}, "a:name>x"},
		{"andere Art trägt keinen Stand", []TransformationRecord{{Kind: AdministrationRequestExcludeColumn, Name: "a", Spec: "x"}}, ""},
	} {
		got, err := FoldTransformations(tc.records)
		if err != nil {
			t.Fatalf("%s: %v", tc.label, err)
		}
		if names(got) != tc.want {
			t.Fatalf("%s: Regelstand = %q, wollen %q", tc.label, names(got), tc.want)
		}
	}
}

// Eine Zeile, deren Regelform nicht zu einer Regel führt, endet als Fehler
// statt den Stand um sie zu verkürzen. Rot färbende Mutation: den Fehler in
// `FoldTransformations` überspringen (`continue`) — der Stand ist leer, der
// Test färbt rot.
func TestFoldTransformationsFailsVisiblyOnUnparsableRow(t *testing.T) {
	_, err := FoldTransformations([]TransformationRecord{
		{Kind: AdministrationRequestSetTransformation, Name: "a", Spec: `{"kind": "unbekannt"}`},
	})
	if !stderrors.Is(err, domainerrors.ErrUnknownTransformationKind) {
		t.Fatalf("Fehler = %v, wollen ErrUnknownTransformationKind", err)
	}
	if !strings.Contains(err.Error(), `"a"`) {
		t.Fatalf("Fehler = %v, wollen den Regelnamen genannt", err)
	}
	if _, err := FoldTransformations([]TransformationRecord{
		{Kind: AdministrationRequestSetTransformation, Name: "a", Spec: `{"kind": "rename_column", "column": "id", "to": "id"}`},
	}); !stderrors.Is(err, domainerrors.ErrTransformationTargetIsColumn) {
		t.Fatalf("Zeile mit Ziel gleich Quelle: Fehler = %v, wollen ErrTransformationTargetIsColumn", err)
	}
}
