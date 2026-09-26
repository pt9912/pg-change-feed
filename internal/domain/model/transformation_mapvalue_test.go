package model

import (
	stderrors "errors"
	"fmt"
	"sort"
	"strings"
	"testing"

	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

// mustMapValue legt eine `map_value`-Regel an oder bricht den Test ab.
func mustMapValue(t testing.TB, name, column string, values map[string]string) Transformation {
	t.Helper()
	rule, err := NewMapValue(name, column, values)
	if err != nil {
		t.Fatalf("NewMapValue(%q, %q, %v): %v", name, column, values, err)
	}
	return rule
}

var statusValues = map[string]string{"o": "open", "c": "closed"}

// Der Konstruktor erzwingt die Invarianten von `map_value` (`SPEC-030`,
// Randfälle), je Eingabe ein Fall. Rot färbende Mutationen je Fall: die
// Prüfung der Eingabe im Konstruktor entfernen (Regelname, Spalte leer,
// Spalte mit U+0000, leere Zuordnung `encoded == ""`).
func TestNewMapValueInvariants(t *testing.T) {
	cases := []struct {
		name    string
		rule    [2]string
		values  map[string]string
		wantErr error
	}{
		{"gültig", [2]string{"status_lesbar", "status"}, statusValues, nil},
		{"Abbildung auf sich selbst ist zulässig", [2]string{"r", "status"}, map[string]string{"a": "a"}, nil},
		{"mehrere Quellwerte, ein Zielwert", [2]string{"r", "status"}, map[string]string{"o": "x", "c": "x"}, nil},
		{"leere Zeichenkette als Schlüssel und als Wert", [2]string{"r", "status"}, map[string]string{"": ""}, nil},
		{"Regelname leer", [2]string{"", "status"}, statusValues, domainerrors.ErrEmptyIdentifier},
		{"Spalte leer", [2]string{"r", ""}, statusValues, domainerrors.ErrInvalidTransformation},
		{"Spalte mit U+0000", [2]string{"r", "sta\x00tus"}, statusValues, domainerrors.ErrInvalidTransformation},
		{"Zuordnung leer", [2]string{"r", "status"}, map[string]string{}, domainerrors.ErrInvalidTransformation},
		{"Zuordnung nil", [2]string{"r", "status"}, nil, domainerrors.ErrInvalidTransformation},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			rule, err := NewMapValue(c.rule[0], c.rule[1], c.values)
			if c.wantErr == nil {
				if err != nil {
					t.Fatalf("NewMapValue: %v", err)
				}
				if rule.Name() != c.rule[0] || rule.Column() != c.rule[1] || rule.Kind() != TransformationMapValue || rule.To() != "" {
					t.Fatalf("Regel %+v trägt nicht ihre Eingaben %q", rule, c.rule)
				}
				if fmt.Sprint(rule.Values()) != fmt.Sprint(c.values) {
					t.Fatalf("Zuordnung %v, wollen %v", rule.Values(), c.values)
				}
				return
			}
			if !stderrors.Is(err, c.wantErr) {
				t.Fatalf("Fehler = %v, wollen %v", err, c.wantErr)
			}
			if rule != (Transformation{}) {
				t.Fatalf("Regel %+v, wollen den Nullwert bei einem Fehler", rule)
			}
		})
	}
}

// Die Regel ist unveränderlich und über `==` vergleichbar (der Backfill-Lauf
// vergleicht Regelstände als Menge): eine spätere Änderung der übergebenen
// Zuordnung und der zurückgegebenen Kopie wirkt nicht auf sie, gleiche
// Zuordnungen ergeben gleiche Regeln unabhängig von der Einfüge-Ordnung, eine
// andere Zuordnung eine ungleiche. Rot färbende Mutationen: die Schlüssel in
// `encodeValueMap` nicht sortieren (Fall „Einfüge-Ordnung“; die Zuordnung
// hat 32 Schlüssel, die Map-Ordnung von Go ist zufällig); `values` beim Bau
// der Regel nicht setzen (Fall „andere Zuordnung“).
func TestMapValueRuleIsImmutableAndComparable(t *testing.T) {
	input := map[string]string{"o": "open"}
	rule := mustMapValue(t, "r", "status", input)
	input["o"] = "verändert"
	input["neu"] = "neu"
	if got := rule.Values(); fmt.Sprint(got) != "map[o:open]" {
		t.Fatalf("Zuordnung nach Änderung der Eingabe = %v, wollen map[o:open]", got)
	}
	rule.Values()["o"] = "verändert"
	if got := rule.Values(); fmt.Sprint(got) != "map[o:open]" {
		t.Fatalf("Zuordnung nach Änderung der Kopie = %v, wollen map[o:open]", got)
	}

	keys := make([]string, 32)
	for i := range keys {
		keys[i] = fmt.Sprintf("k%02d", i)
	}
	build := func(order []string) Transformation {
		values := make(map[string]string, len(order))
		for _, key := range order {
			values[key] = "v-" + key
		}
		return mustMapValue(t, "r", "status", values)
	}
	forward := build(keys)
	reversed := append([]string(nil), keys...)
	sort.Sort(sort.Reverse(sort.StringSlice(reversed)))
	for i := 0; i < 20; i++ {
		if again := build(reversed); again != forward {
			t.Fatalf("Durchlauf %d: gleiche Zuordnung ergibt eine ungleiche Regel", i)
		}
	}
	if other := mustMapValue(t, "r", "status", map[string]string{"o": "open", "c": "zu"}); other == mustMapValue(t, "r", "status", statusValues) {
		t.Fatal("andere Zuordnung ergibt eine gleiche Regel")
	}
	if other := mustMapValue(t, "r", "status", map[string]string{"o": "open", "c": "closed", "x": "y"}); other == mustMapValue(t, "r", "status", statusValues) {
		t.Fatal("größere Zuordnung ergibt eine gleiche Regel")
	}
	if renamed := mustRename(t, "r", "status", "state"); renamed.Values() != nil {
		t.Fatalf("rename_column trägt eine Zuordnung: %v", renamed.Values())
	}
}

// Die Anwendbarkeit von `map_value` hängt an der Spalte (`SPEC-030`,
// Anwendbarkeit): die Regel trägt keinen Zielnamen und prüft keine Kollision,
// auch nicht mit einer Spalte, deren Name leer ist (der leere Zielname von
// `rename_column` wäre eine). Rot färbende Mutationen: `containsName(columns,
// t.column)` entfernen (Fall „Spalte fehlt“); die Prüfung `t.kind ==
// TransformationRenameColumn` in `CheckApplicable` entfernen (Fall „leerer
// Spaltenname“).
func TestCheckApplicableMapValue(t *testing.T) {
	rule := mustMapValue(t, "r", "status", statusValues)
	for _, c := range []struct {
		name    string
		columns []string
		wantErr error
	}{
		{"Spalte da", []string{"id", "status"}, nil},
		{"Spalte fehlt", []string{"id", "name"}, domainerrors.ErrTransformationColumnMissing},
		{"Spaltenmenge leer", nil, domainerrors.ErrTransformationColumnMissing},
		{"Groß-/Kleinschreibung wird nicht gefaltet", []string{"id", "Status"}, domainerrors.ErrTransformationColumnMissing},
		{"Wert der Zuordnung gleicht einer Spalte", []string{"status", "open", "closed"}, nil},
		{"leerer Spaltenname kollidiert nicht", []string{"status", ""}, nil},
	} {
		err := rule.CheckApplicable(c.columns)
		if !stderrors.Is(err, c.wantErr) && !(c.wantErr == nil && err == nil) {
			t.Fatalf("%s: Fehler = %v, wollen %v", c.name, err, c.wantErr)
		}
	}
}

// `map_value` ersetzt den Wert, wenn er als Zeichenkette ein Schlüssel der
// Zuordnung ist, und lässt Schlüssel und Position der Quellspalte, jeden
// anderen Wert und jede Abwesenheit unverändert (`SPEC-030`). Rot färbende
// Mutationen je Fall: in `applyTransformations` den Zweig `TransformationMapValue`
// entfernen (Fälle mit Zuordnung), den Wert immer ersetzen bzw. den Schlüssel
// aus `rule.to` setzen (Fälle „nicht zugeordnet“, „Schlüssel bleibt“), die
// Suche in `lookupMappedValue` nach der Länge statt nach dem Inhalt (Fall
// „gleiche Länge“) bzw. ohne die Länge des Schlüssels (Fall „Präfix“),
// den Wert vor der Suche normalisieren (`strings.ToLower`, `TrimSpace`; Fälle
// „Groß-/Kleinschreibung“, „Leerraum“) bzw. das Ergebnis erneut nachschlagen
// (Fall „keine Kette“), `BuildRowImage` ohne `applyTransformations` für
// abwesende Werte auswerten (Fall „NULL“).
func TestBuildRowImageMapValue(t *testing.T) {
	cases := []struct {
		name     string
		values   []*string
		excluded []string
		rules    []Transformation
		want     string
	}{
		{
			name:   "zugeordneter Wert",
			values: []*string{textValue("7"), nil, textValue("Ada"), textValue("o")},
			rules:  []Transformation{mustMapValue(t, "r", "status", statusValues)},
			want:   `{"id":"7","name":"Ada","status":"open"}`,
		},
		{
			name:   "zweiter zugeordneter Wert",
			values: []*string{textValue("7"), nil, textValue("Ada"), textValue("c")},
			rules:  []Transformation{mustMapValue(t, "r", "status", statusValues)},
			want:   `{"id":"7","name":"Ada","status":"closed"}`,
		},
		{
			name:   "nicht zugeordneter Wert bleibt",
			values: []*string{textValue("7"), nil, textValue("Ada"), textValue("x")},
			rules:  []Transformation{mustMapValue(t, "r", "status", statusValues)},
			want:   `{"id":"7","name":"Ada","status":"x"}`,
		},
		{
			name:   "Schlüssel und Position der Quellspalte bleiben (mittlere Spalte)",
			values: []*string{textValue("7"), nil, textValue("Ada"), textValue("o")},
			rules:  []Transformation{mustMapValue(t, "r", "name", map[string]string{"Ada": "Ada Lovelace"})},
			want:   `{"id":"7","name":"Ada Lovelace","status":"o"}`,
		},
		{
			name:   "erste Spalte",
			values: []*string{textValue("7"), nil, textValue("Ada"), textValue("o")},
			rules:  []Transformation{mustMapValue(t, "r", "id", map[string]string{"7": "sieben"})},
			want:   `{"id":"sieben","name":"Ada","status":"o"}`,
		},
		{
			name:   "NULL bleibt Abwesenheit, kein Schlüssel, kein Platzhalter",
			values: []*string{textValue("7"), nil, textValue("Ada"), nil},
			rules:  []Transformation{mustMapValue(t, "r", "status", statusValues)},
			want:   `{"id":"7","name":"Ada"}`,
		},
		{
			name:   "Spalte ohne Wert bleibt Abwesenheit",
			values: []*string{textValue("7")},
			rules:  []Transformation{mustMapValue(t, "r", "status", statusValues)},
			want:   `{"id":"7"}`,
		},
		{
			name:     "ausgeschlossene Spalte trägt weder Schlüssel noch Wert noch abgebildeten Wert",
			values:   []*string{textValue("7"), textValue("s"), textValue("Ada"), textValue("o")},
			excluded: []string{"status"},
			rules:    []Transformation{mustMapValue(t, "r", "status", statusValues)},
			want:     `{"id":"7","secret":"s","name":"Ada"}`,
		},
		{
			name:   "Groß-/Kleinschreibung zählt",
			values: []*string{textValue("7"), nil, textValue("Ada"), textValue("O")},
			rules:  []Transformation{mustMapValue(t, "r", "status", statusValues)},
			want:   `{"id":"7","name":"Ada","status":"O"}`,
		},
		{
			name:   "Leerraum wird nicht abgeschnitten",
			values: []*string{textValue("7"), nil, textValue("Ada"), textValue("o ")},
			rules:  []Transformation{mustMapValue(t, "r", "status", statusValues)},
			want:   `{"id":"7","name":"Ada","status":"o "}`,
		},
		{
			name:   "Präfix eines Schlüssels ist kein Schlüssel",
			values: []*string{textValue("7"), nil, textValue("Ada"), textValue("op")},
			rules:  []Transformation{mustMapValue(t, "r", "status", statusValues)},
			want:   `{"id":"7","name":"Ada","status":"op"}`,
		},
		{
			name:   "gleiche Länge, anderer Inhalt",
			values: []*string{textValue("7"), nil, textValue("Ada"), textValue("x")},
			rules:  []Transformation{mustMapValue(t, "r", "status", map[string]string{"o": "open"})},
			want:   `{"id":"7","name":"Ada","status":"x"}`,
		},
		{
			name:   "keine Kette: der abgebildete Wert wird nicht erneut nachgeschlagen",
			values: []*string{textValue("7"), nil, textValue("Ada"), textValue("a")},
			rules:  []Transformation{mustMapValue(t, "r", "status", map[string]string{"a": "b", "b": "c"})},
			want:   `{"id":"7","name":"Ada","status":"b"}`,
		},
		{
			name:   "Abbildung auf sich selbst ändert den Wert nicht",
			values: []*string{textValue("7"), nil, textValue("Ada"), textValue("o")},
			rules:  []Transformation{mustMapValue(t, "r", "status", map[string]string{"o": "o"})},
			want:   `{"id":"7","name":"Ada","status":"o"}`,
		},
		{
			name:   "mehrere Quellwerte auf denselben Zielwert",
			values: []*string{textValue("7"), nil, textValue("Ada"), textValue("c")},
			rules:  []Transformation{mustMapValue(t, "r", "status", map[string]string{"o": "offen", "c": "offen"})},
			want:   `{"id":"7","name":"Ada","status":"offen"}`,
		},
		{
			name:   "leere Zeichenkette ist ein Wert, kein NULL",
			values: []*string{textValue("7"), nil, textValue("Ada"), textValue("")},
			rules:  []Transformation{mustMapValue(t, "r", "status", map[string]string{"": "leer"})},
			want:   `{"id":"7","name":"Ada","status":"leer"}`,
		},
		{
			name:   "abgebildeter Wert trägt die Maskierung von encoding/json",
			values: []*string{textValue("7"), nil, textValue("Ada"), textValue("o")},
			rules:  []Transformation{mustMapValue(t, "r", "status", map[string]string{"o": "<a>&\"b\"\n"})},
			want:   "{\"id\":\"7\",\"name\":\"Ada\",\"status\":\"\x5cu003ca\x5cu003e\x5cu0026\\\"b\\\"\x5cn\"}",
		},
		{
			name:   "Regel auf eine Spalte außerhalb der Spaltenliste wirkt nicht",
			values: []*string{textValue("7"), nil, textValue("Ada"), textValue("o")},
			rules:  []Transformation{mustMapValue(t, "r", "unbekannt", statusValues)},
			want:   `{"id":"7","name":"Ada","status":"o"}`,
		},
		{
			name:   "Umbenennung und Wertabbildung auf verschiedenen Spalten",
			values: []*string{textValue("7"), nil, textValue("Ada"), textValue("o")},
			rules: []Transformation{
				mustRename(t, "a", "name", "customer_name"),
				mustMapValue(t, "b", "status", statusValues),
			},
			want: `{"id":"7","customer_name":"Ada","status":"open"}`,
		},
		{
			name:   "die Ordnung der Regeln ändert das Bild nicht",
			values: []*string{textValue("7"), nil, textValue("Ada"), textValue("o")},
			rules: []Transformation{
				mustMapValue(t, "b", "status", statusValues),
				mustRename(t, "a", "name", "customer_name"),
			},
			want: `{"id":"7","customer_name":"Ada","status":"open"}`,
		},
		{
			name:   "abgebildeter Wert gleicht dem Namen einer Spalte: kein Schlüsselkonflikt",
			values: []*string{textValue("7"), nil, textValue("Ada"), textValue("o")},
			rules:  []Transformation{mustMapValue(t, "r", "status", map[string]string{"o": "name"})},
			want:   `{"id":"7","name":"Ada","status":"name"}`,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := BuildRowImage(transformationColumns, c.values, c.excluded, c.rules)
			if err != nil {
				t.Fatalf("BuildRowImage: %v", err)
			}
			if string(got) != c.want {
				t.Fatalf("Bild %q, wollen %q", got, c.want)
			}
		})
	}
}

// Die kanonische Kodierung trägt Felder jeder Länge: Schlüssel und Werte mit
// 0, 1, 255, 256, 65535, 65536 und 70000 Byte (jede Stufe, an der ein weiteres
// Längen-Byte gebraucht wird) kommen aus `Values()` unverändert zurück, und
// die Suche findet jeden Schlüssel und lässt einen Wert anderer Länge stehen.
// Rot färbende Mutationen: das Längen-Präfix in `writeValueField` bzw.
// `readValueField` auf die niederen Bytes kürzen (Fälle ab 256 Byte).
func TestMapValueEncodingCarriesFieldsOfAnyLength(t *testing.T) {
	lengths := []int{0, 1, 255, 256, 65535, 65536, 70000}
	values := make(map[string]string, len(lengths))
	for i, length := range lengths {
		values[strings.Repeat("k", length)+fmt.Sprint(i)] = strings.Repeat("v", lengths[len(lengths)-1-i]) + fmt.Sprint(i)
	}
	rule := mustMapValue(t, "r", "status", values)
	if got := rule.Values(); !equalMaps(got, values) {
		t.Fatalf("Zuordnung kommt nicht unverändert zurück: %d Paare, wollen %d", len(got), len(values))
	}
	for key, want := range values {
		image, err := BuildRowImage([]string{"status"}, []*string{textValue(key)}, nil, []Transformation{rule})
		if err != nil || string(image) != `{"status":"`+want+`"}` {
			t.Fatalf("Schlüssel der Länge %d: Bild der Länge %d, Fehler %v, wollen den zugeordneten Wert der Länge %d", len(key), len(image), err, len(want))
		}
	}
	stays := strings.Repeat("k", 256) + "x"
	image, err := BuildRowImage([]string{"status"}, []*string{textValue(stays)}, nil, []Transformation{rule})
	if err != nil || string(image) != `{"status":"`+stays+`"}` {
		t.Fatalf("Wert ohne Zuordnung: Bild der Länge %d, Fehler %v, wollen den Wert unverändert", len(image), err)
	}
}

func equalMaps(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for key, value := range a {
		if other, ok := b[key]; !ok || other != value {
			return false
		}
	}
	return true
}

// Gleiche Regelmenge und gleiche Relation ergeben dasselbe Bild, auch aus
// gleichzeitigen Aufrufen (`-race`); die Regelliste bleibt unverändert und
// das Ergebnis teilt keinen Speicher mit ihr.
func TestBuildRowImageMapValueIsDeterministicAndPure(t *testing.T) {
	rules := []Transformation{
		mustRename(t, "a", "name", "customer_name"),
		mustMapValue(t, "b", "status", statusValues),
	}
	values := []*string{textValue("7"), textValue("s"), textValue("Ada"), textValue("o")}
	first, err := BuildRowImage(transformationColumns, values, []string{"secret"}, rules)
	if err != nil {
		t.Fatalf("BuildRowImage: %v", err)
	}
	want := string(first)
	first[0] = '!'
	results := make(chan string, 8)
	for i := 0; i < cap(results); i++ {
		go func() {
			got, _ := BuildRowImage(transformationColumns, values, []string{"secret"}, rules)
			results <- string(got)
		}()
	}
	for i := 0; i < cap(results); i++ {
		if got := <-results; got != want {
			t.Fatalf("Aufruf %d: Bild %q, wollen %q", i, got, want)
		}
	}
	if got := rules[1].Values(); fmt.Sprint(got) != "map[c:closed o:open]" {
		t.Fatalf("Zuordnung verändert: %v", got)
	}
}

// ParseTransformationSpec liest `map_value` strikt: Schlüssel, Pflichtschlüssel
// und die Form von `values` (`SPEC-030`, `SPEC-019` Fehlertext-Tabelle). Rot
// färbende Mutationen je Zeile: den Zweig `TransformationMapValue` in
// `ParseTransformationSpec` entfernen (Fälle ohne Fehler enden als unbekannter
// Regeltyp); `"values"` aus `allowedRuleKeys` streichen (der gültige Fall
// endet als unbekannter Schlüssel); `to` in die erlaubten Schlüssel
// aufnehmen (Fall „to“); in `jsonStringMap` den Typ des Werts nicht prüfen
// (Fälle „Zahl“, „null“, „Objekt“, „Array“ als Wert); die leere Zuordnung
// nicht ablehnen (Fälle „values leer“, „values fehlt“).
func TestParseTransformationSpecMapValue(t *testing.T) {
	invalid := domainerrors.ErrInvalidRuleSpec
	for _, tc := range []struct {
		name   string
		text   string
		reason error
		detail string
	}{
		{"values fehlt", `{"kind": "map_value", "column": "status"}`, invalid, ""},
		{"values leer", `{"kind": "map_value", "column": "status", "values": {}}`, invalid, ""},
		{"values null", `{"kind": "map_value", "column": "status", "values": null}`, invalid, ""},
		{"values ist ein Array", `{"kind": "map_value", "column": "status", "values": ["o"]}`, invalid, ""},
		{"values ist eine Zeichenkette", `{"kind": "map_value", "column": "status", "values": "o"}`, invalid, ""},
		{"Wert ist eine Zahl", `{"kind": "map_value", "column": "status", "values": {"o": 1}}`, invalid, ""},
		{"Wert ist null", `{"kind": "map_value", "column": "status", "values": {"o": null}}`, invalid, ""},
		{"Wert ist ein Objekt", `{"kind": "map_value", "column": "status", "values": {"o": {}}}`, invalid, ""},
		{"Wert ist ein Array", `{"kind": "map_value", "column": "status", "values": {"o": ["x"]}}`, invalid, ""},
		{"ein Wert ist keine Zeichenkette, ein anderer schon", `{"kind": "map_value", "column": "status", "values": {"o": "open", "c": true}}`, invalid, ""},
		{"column fehlt", `{"kind": "map_value", "values": {"o": "open"}}`, invalid, ""},
		{"column leer", `{"kind": "map_value", "column": "", "values": {"o": "open"}}`, invalid, ""},
		{"column ist eine Zahl", `{"kind": "map_value", "column": 1, "values": {"o": "open"}}`, invalid, ""},
		{"column mit NUL", "{\"kind\": \"map_value\", \"column\": \"a\\u0000b\", \"values\": {\"o\": \"open\"}}", invalid, ""},
		{"to gehört nicht zu map_value", `{"kind": "map_value", "column": "status", "to": "x", "values": {"o": "open"}}`, domainerrors.ErrUnknownRuleSpecKey, "to"},
		{"unbekannter Schlüssel", `{"kind": "map_value", "column": "status", "values": {"o": "open"}, "extra": 1}`, domainerrors.ErrUnknownRuleSpecKey, "extra"},
		{"unbekannter Schlüssel vor ungültigem values", `{"kind": "map_value", "column": "status", "values": {}, "extra": 1}`, domainerrors.ErrUnknownRuleSpecKey, "extra"},
		{"Schlüssel mit anderer Schreibweise", `{"kind": "map_value", "column": "status", "Values": {"o": "open"}}`, domainerrors.ErrUnknownRuleSpecKey, "Values"},
	} {
		_, err := ParseTransformationSpec(tc.text)
		if !stderrors.Is(err, tc.reason) {
			t.Fatalf("%s: Fehler = %v, wollen %v", tc.name, err, tc.reason)
		}
		if tc.detail != "" {
			var specErr *TransformationSpecError
			if !stderrors.As(err, &specErr) || specErr.Detail != tc.detail {
				t.Fatalf("%s: Fehler = %v, wollen den Wert %q", tc.name, err, tc.detail)
			}
		}
	}

	spec, err := ParseTransformationSpec(` {"kind": "map_value", "column": "status", "values": {"o": "open", "c": "closed"}} `)
	if err != nil {
		t.Fatalf("gültige Form: %v", err)
	}
	if spec.Kind() != TransformationMapValue || spec.Column() != "status" || spec.Target() != "" {
		t.Fatalf("Regelform = %+v, wollen map_value auf status ohne Zielname", spec)
	}
	for _, text := range []string{
		`{"kind": "map_value", "column": "status", "values": {"o": "o"}}`,
		`{"kind": "map_value", "column": "status", "values": {"o": "x", "c": "x"}}`,
		`{"kind": "map_value", "column": "status", "values": {"": ""}}`,
		`{"kind": "map_value", "column": "spalte ä.\"x\"", "values": {"ü": "ä😀"}}`,
	} {
		if _, err := ParseTransformationSpec(text); err != nil {
			t.Fatalf("%s: %v", text, err)
		}
	}

	// Ein doppelter Schlüssel kommt aus einer `jsonb`-Spalte nie an; der Text
	// trägt den letzten Wert.
	last, err := ParseTransformationSpec(`{"kind": "map_value", "column": "status", "values": {"o": "erster", "o": "letzter"}}`)
	if err != nil {
		t.Fatalf("doppelter Schlüssel: %v", err)
	}
	rule, err := last.Build("r")
	if err != nil || fmt.Sprint(rule.Values()) != "map[o:letzter]" {
		t.Fatalf("Regel = %v, Fehler %v, wollen map[o:letzter]", rule.Values(), err)
	}
}

// Die Ordnung der Schlüssel in `values` ändert die Regel nicht: die Form ist
// die eine Quelle der Regel (Build), Regeln aus zwei Schreibweisen desselben
// Inhalts sind gleich. Rot färbende Mutation: die Sortierung in
// `encodeValueMap` entfernen.
func TestTransformationSpecMapValueBuild(t *testing.T) {
	build := func(text string) Transformation {
		spec, err := ParseTransformationSpec(text)
		if err != nil {
			t.Fatalf("ParseTransformationSpec: %v", err)
		}
		rule, err := spec.Build("status_lesbar")
		if err != nil {
			t.Fatalf("Build: %v", err)
		}
		return rule
	}
	first := build(`{"kind": "map_value", "column": "status", "values": {"o": "open", "c": "closed", "p": "pending", "d": "done"}}`)
	second := build(`{"values": {"d": "done", "p": "pending", "c": "closed", "o": "open"}, "column": "status", "kind": "map_value"}`)
	if first != second {
		t.Fatalf("Regeln aus zwei Schreibweisen desselben Inhalts sind ungleich: %+v, %+v", first, second)
	}
	if first.Name() != "status_lesbar" || first.Kind() != TransformationMapValue || first.Column() != "status" || first.To() != "" || len(first.Values()) != 4 {
		t.Fatalf("Regel = %+v", first)
	}
	if _, err := (TransformationSpec{kind: TransformationMapValue, column: "status"}).Build("r"); !stderrors.Is(err, domainerrors.ErrInvalidTransformation) {
		t.Fatalf("Build ohne Zuordnung: Fehler = %v, wollen ErrInvalidTransformation", err)
	}
}

// `map_value` trägt keinen Zielnamen: K3 wirkt für ihn nicht, K1 und K2 wirken
// wie für jede Regel (`SPEC-019`; K2 ist die Aussage „eine Quellspalte trägt
// höchstens eine Spaltenregel“). Rot färbende Mutation je Fall: K2 (die
// Spaltenschleife) streichen (die K2-Fälle); `s.to != ""` in `CheckConflicts`
// entfernen (Fall „kein K3 gegen eine zweite Wertabbildung“: der leere
// Zielname zweier Wertabbildungen gälte als Kollision).
func TestCheckConflictsMapValue(t *testing.T) {
	specOf := func(text string) TransformationSpec {
		spec, err := ParseTransformationSpec(text)
		if err != nil {
			t.Fatalf("ParseTransformationSpec(%s): %v", text, err)
		}
		return spec
	}
	mapSpec := specOf(`{"kind": "map_value", "column": "status", "values": {"o": "open"}}`)
	columns := []string{"id", "name", "status", "secret"}
	renamed := mustRename(t, "kundenname", "name", "customer_name")
	mapped := mustMapValue(t, "status_lesbar", "status", statusValues)

	for _, tc := range []struct {
		label string
		name  string
		spec  TransformationSpec
		rules []Transformation
		want  error
	}{
		{"frei", "zweite", mapSpec, []Transformation{renamed}, nil},
		{"leerer Regelstand", "zweite", mapSpec, nil, nil},
		{"K1 Name vergeben", "kundenname", mapSpec, []Transformation{renamed}, domainerrors.ErrRuleNameTaken},
		{"K2 Spalte trägt eine Umbenennung", "zweite", specOf(`{"kind": "map_value", "column": "name", "values": {"a": "b"}}`), []Transformation{renamed}, domainerrors.ErrColumnHasRule},
		{"K2 Spalte trägt eine Wertabbildung", "zweite", mapSpec, []Transformation{mapped}, domainerrors.ErrColumnHasRule},
		{"K2 Umbenennung gegen eine Wertabbildung derselben Spalte", "zweite", specOf(`{"kind": "rename_column", "column": "status", "to": "state"}`), []Transformation{mapped}, domainerrors.ErrColumnHasRule},
		{"kein K3 gegen eine Umbenennung", "zweite", mapSpec, []Transformation{renamed, mustRename(t, "x", "id", "pk")}, nil},
		{"kein K3 gegen eine zweite Wertabbildung auf anderer Spalte", "zweite", mapSpec, []Transformation{mustMapValue(t, "x", "name", statusValues)}, nil},
	} {
		got := tc.spec.CheckConflicts(tc.name, tc.rules, columns)
		if !stderrors.Is(got, tc.want) && !(tc.want == nil && got == nil) {
			t.Fatalf("%s: Fehler = %v, wollen %v", tc.label, got, tc.want)
		}
	}
}

// Der Regelstand faltet `map_value`-Zeilen wie jede Regel: Set trägt ein,
// Remove nimmt heraus, ein Set unter vorhandenem Namen ersetzt an der Stelle
// (die Zuordnung ist dann die neue); eine nicht mehr lesbare Zeile endet als
// Fehler. Rot färbende Mutation: den Zweig `TransformationMapValue` in
// `TransformationSpec.Build` entfernen (die Fälle enden als unbekannter
// Regeltyp).
func TestFoldTransformationsMapValue(t *testing.T) {
	set := func(name, column, values string) TransformationRecord {
		return TransformationRecord{
			Kind: AdministrationRequestSetTransformation, Name: name,
			Spec: `{"kind": "map_value", "column": "` + column + `", "values": ` + values + `}`,
		}
	}
	remove := func(name string) TransformationRecord {
		return TransformationRecord{Kind: AdministrationRequestRemoveTransformation, Name: name}
	}
	describe := func(rules []Transformation) string {
		out := make([]string, 0, len(rules))
		for _, rule := range rules {
			out = append(out, fmt.Sprintf("%s:%s:%s:%v", rule.Name(), rule.Kind(), rule.Column(), rule.Values()))
		}
		return strings.Join(out, ",")
	}
	for _, tc := range []struct {
		label   string
		records []TransformationRecord
		want    string
	}{
		{"Set", []TransformationRecord{set("a", "status", `{"o": "open"}`)}, "a:map_value:status:map[o:open]"},
		{"Set, Remove", []TransformationRecord{set("a", "status", `{"o": "open"}`), remove("a")}, ""},
		{"Set, Remove, Set (Zyklus)", []TransformationRecord{set("a", "status", `{"o": "open"}`), remove("a"), set("a", "status", `{"o": "offen"}`)}, "a:map_value:status:map[o:offen]"},
		{"Set unter vorhandenem Namen ersetzt die Zuordnung", []TransformationRecord{set("a", "status", `{"o": "open"}`), set("a", "status", `{"c": "closed"}`)}, "a:map_value:status:map[c:closed]"},
		{"neben einer Umbenennung", []TransformationRecord{
			{Kind: AdministrationRequestSetTransformation, Name: "u", Spec: `{"kind": "rename_column", "column": "name", "to": "x"}`},
			set("a", "status", `{"o": "open"}`),
		}, "u:rename_column:name:map[],a:map_value:status:map[o:open]"},
	} {
		got, err := FoldTransformations(tc.records)
		if err != nil {
			t.Fatalf("%s: %v", tc.label, err)
		}
		if describe(got) != tc.want {
			t.Fatalf("%s: Regelstand = %q, wollen %q", tc.label, describe(got), tc.want)
		}
	}
	if _, err := FoldTransformations([]TransformationRecord{set("a", "status", `{}`)}); !stderrors.Is(err, domainerrors.ErrInvalidRuleSpec) || !strings.Contains(err.Error(), `"a"`) {
		t.Fatalf("Zeile mit leerer Zuordnung: Fehler = %v, wollen ErrInvalidRuleSpec mit dem Regelnamen", err)
	}
}
