package model

import (
	stderrors "errors"
	"strings"
	"sync"
	"testing"

	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

// mustRename legt eine `rename_column`-Regel an oder bricht den Test ab.
func mustRename(t testing.TB, name, column, to string) Transformation {
	t.Helper()
	rule, err := NewRenameColumn(name, column, to)
	if err != nil {
		t.Fatalf("NewRenameColumn(%q, %q, %q): %v", name, column, to, err)
	}
	return rule
}

// Der Konstruktor erzwingt die Invarianten einer einzelnen Regel
// (`SPEC-030`, Bezeichner), je Eingabe ein Fall. Rot färbende Mutationen je
// Fall: die Prüfung der Eingabe im Konstruktor entfernen bzw. die Grenze
// verschieben (`> 63` auf `> 64`).
func TestNewRenameColumnInvariants(t *testing.T) {
	byte63 := strings.Repeat("a", 63)
	cases := []struct {
		name    string
		rule    [3]string
		wantErr error
	}{
		{"gültig", [3]string{"kundenname", "name", "customer_name"}, nil},
		{"Zielname mit 63 Byte", [3]string{"r", "name", byte63}, nil},
		{"Zielname mit 63 Byte über Mehrbyte-Zeichen", [3]string{"r", "name", strings.Repeat("ä", 31) + "a"}, nil},
		{"Zielname mit Sonderzeichen ist zulässig", [3]string{"r", "name", "Kunden Name.\"x\""}, nil},
		{"Regelname leer", [3]string{"", "name", "x"}, domainerrors.ErrEmptyIdentifier},
		{"Spalte leer", [3]string{"r", "", "x"}, domainerrors.ErrInvalidTransformation},
		{"Spalte mit U+0000", [3]string{"r", "na\x00me", "x"}, domainerrors.ErrInvalidTransformation},
		{"Zielname leer", [3]string{"r", "name", ""}, domainerrors.ErrInvalidTransformation},
		{"Zielname mit U+0000", [3]string{"r", "name", "x\x00y"}, domainerrors.ErrInvalidTransformation},
		{"Zielname mit 64 Byte", [3]string{"r", "name", byte63 + "x"}, domainerrors.ErrInvalidTransformation},
		{"Zielname mit 64 Byte über Mehrbyte-Zeichen", [3]string{"r", "name", strings.Repeat("ä", 32)}, domainerrors.ErrInvalidTransformation},
		{"Zielname gleicht der Quellspalte", [3]string{"r", "name", "name"}, domainerrors.ErrTransformationTargetIsColumn},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			rule, err := NewRenameColumn(c.rule[0], c.rule[1], c.rule[2])
			if c.wantErr == nil {
				if err != nil {
					t.Fatalf("NewRenameColumn: %v", err)
				}
				if rule.Name() != c.rule[0] || rule.Column() != c.rule[1] || rule.To() != c.rule[2] || rule.Kind() != TransformationRenameColumn {
					t.Fatalf("Regel %+v trägt nicht ihre Eingaben %q", rule, c.rule)
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

// `to` gleich `column` trägt einen eigenen Sentinel und ist keine
// Formverletzung: der Antragsweg bildet ihn auf den Text von K3 ab, nicht auf
// `rule_spec ist ungültig` (`SPEC-019`).
func TestNewRenameColumnTargetIsColumnIsNotAFormViolation(t *testing.T) {
	_, err := NewRenameColumn("r", "name", "name")
	if stderrors.Is(err, domainerrors.ErrInvalidTransformation) {
		t.Fatalf("Fehler %v ist als Formverletzung erkennbar, wollen ausschließlich ErrTransformationTargetIsColumn", err)
	}
	if !stderrors.Is(err, domainerrors.ErrTransformationTargetIsColumn) {
		t.Fatalf("Fehler = %v, wollen ErrTransformationTargetIsColumn", err)
	}
}

// Die Regeltypen sind eine geschlossene Menge; ein Aufruf liefert eine neue
// Liste, ein Aufrufer kann die Menge nicht verändern.
func TestTransformationKindsIsAClosedSet(t *testing.T) {
	kinds := TransformationKinds()
	if len(kinds) != 1 || kinds[0] != TransformationRenameColumn || string(kinds[0]) != "rename_column" {
		t.Fatalf("Regeltypen = %v, wollen genau rename_column", kinds)
	}
	kinds[0] = "verändert"
	if again := TransformationKinds(); again[0] != TransformationRenameColumn {
		t.Fatalf("Regeltypen nach Veränderung der ersten Liste = %v, wollen rename_column", again)
	}
}

// Die Anwendbarkeit hängt an Regel und Spaltenmenge (`SPEC-030`), zeichengenau
// verglichen. Rot färbende Mutationen je Fall: `containsName(columns,
// t.column)` bzw. `containsName(columns, t.to)` in `CheckApplicable`
// entfernen.
func TestCheckApplicable(t *testing.T) {
	rule := mustRename(t, "r", "name", "customer_name")
	cases := []struct {
		name    string
		columns []string
		wantErr error
	}{
		{"Spalte da, Zielname frei", []string{"id", "name", "status"}, nil},
		{"Spalte fehlt", []string{"id", "status"}, domainerrors.ErrTransformationColumnMissing},
		{"Spaltenmenge leer", nil, domainerrors.ErrTransformationColumnMissing},
		{"Zielname kollidiert mit einer Spalte", []string{"id", "name", "customer_name"}, domainerrors.ErrTransformationTargetCollides},
		{"Groß-/Kleinschreibung wird nicht gefaltet: Spalte", []string{"id", "Name"}, domainerrors.ErrTransformationColumnMissing},
		{"Groß-/Kleinschreibung wird nicht gefaltet: Zielname", []string{"id", "name", "Customer_Name"}, nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := rule.CheckApplicable(c.columns)
			if c.wantErr == nil {
				if err != nil {
					t.Fatalf("CheckApplicable: %v", err)
				}
				return
			}
			if !stderrors.Is(err, c.wantErr) {
				t.Fatalf("Fehler = %v, wollen %v", err, c.wantErr)
			}
		})
	}
	// Der Nullwert trägt keine Spalte und ist auf keine Änderung anwendbar.
	if err := (Transformation{}).CheckApplicable([]string{"id"}); !stderrors.Is(err, domainerrors.ErrTransformationColumnMissing) {
		t.Fatalf("Nullwert: Fehler = %v, wollen ErrTransformationColumnMissing", err)
	}
}

var transformationColumns = []string{"id", "secret", "name", "status"}

// `rename_column` setzt den Schlüssel unter den Zielnamen, an die Position
// seiner Quellspalte, mit unverändertem Wert (`SPEC-030`, `ADR-0112`
// Teilfrage 3). Rot färbende Mutationen je Fall: in `applyTransformations`
// `rule.to` durch `column` ersetzen (Schlüssel bleibt), den Wert verändern,
// bzw. das Bild in `BuildRowImage` nach Regel-Schlüsseln hinten anfügen
// (Position).
func TestBuildRowImageRenameColumn(t *testing.T) {
	cases := []struct {
		name     string
		values   []*string
		excluded []string
		rules    []Transformation
		want     string
	}{
		{
			name:   "mittlere Spalte behält ihre Position",
			values: []*string{textValue("7"), nil, textValue("Ada"), textValue("o")},
			rules:  []Transformation{mustRename(t, "r", "name", "customer_name")},
			want:   `{"id":"7","customer_name":"Ada","status":"o"}`,
		},
		{
			name:   "erste Spalte",
			values: []*string{textValue("7"), nil, textValue("Ada"), textValue("o")},
			rules:  []Transformation{mustRename(t, "r", "id", "pk")},
			want:   `{"pk":"7","name":"Ada","status":"o"}`,
		},
		{
			name:   "letzte Spalte",
			values: []*string{textValue("7"), nil, textValue("Ada"), textValue("o")},
			rules:  []Transformation{mustRename(t, "r", "status", "state")},
			want:   `{"id":"7","name":"Ada","state":"o"}`,
		},
		{
			name:   "zwei Regeln auf verschiedenen Spalten",
			values: []*string{textValue("7"), nil, textValue("Ada"), textValue("o")},
			rules: []Transformation{
				mustRename(t, "a", "name", "customer_name"),
				mustRename(t, "b", "status", "state"),
			},
			want: `{"id":"7","customer_name":"Ada","state":"o"}`,
		},
		{
			name:   "NULL bleibt Abwesenheit, kein Zielschlüssel",
			values: []*string{textValue("7"), nil, nil, textValue("o")},
			rules:  []Transformation{mustRename(t, "r", "name", "customer_name")},
			want:   `{"id":"7","status":"o"}`,
		},
		{
			name:   "Spalte ohne Wert bleibt Abwesenheit",
			values: []*string{textValue("7")},
			rules:  []Transformation{mustRename(t, "r", "name", "customer_name")},
			want:   `{"id":"7"}`,
		},
		{
			name:     "ausgeschlossene Spalte trägt weder Quell- noch Zielschlüssel noch Wert",
			values:   []*string{textValue("7"), textValue("geheim"), textValue("Ada"), textValue("o")},
			excluded: []string{"secret"},
			rules:    []Transformation{mustRename(t, "r", "secret", "hidden")},
			want:     `{"id":"7","name":"Ada","status":"o"}`,
		},
		{
			name:   "Wert bleibt unverändert, auch mit Sonderzeichen",
			values: []*string{textValue("7"), nil, textValue("<a>&\"b\" \xc3\xbc\n"), nil},
			rules:  []Transformation{mustRename(t, "r", "name", "customer_name")},
			want:   "{\"id\":\"7\",\"customer_name\":\"\x5cu003ca\x5cu003e\x5cu0026\\\"b\\\" \xc3\xbc\x5cn\"}",
		},
		{
			name:   "Zielname trägt die Maskierung von encoding/json",
			values: []*string{textValue("7"), nil, textValue("Ada"), nil},
			rules:  []Transformation{mustRename(t, "r", "name", "a\"b<")},
			want:   "{\"id\":\"7\",\"a\\\"b\x5cu003c\":\"Ada\"}",
		},
		{
			name:   "Regel auf eine Spalte außerhalb der Spaltenliste wirkt nicht",
			values: []*string{textValue("7"), nil, textValue("Ada"), textValue("o")},
			rules:  []Transformation{mustRename(t, "r", "unbekannt", "x")},
			want:   `{"id":"7","name":"Ada","status":"o"}`,
		},
		{
			name:   "die erste treffende Regel entscheidet",
			values: []*string{textValue("7"), nil, textValue("Ada"), nil},
			rules: []Transformation{
				mustRename(t, "a", "name", "erster"),
				mustRename(t, "b", "name", "zweiter"),
			},
			want: `{"id":"7","erster":"Ada"}`,
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

// Ein fehlendes Bild bleibt fehlend, mit und ohne Regel (`LH-FA-CAP-008`).
func TestBuildRowImageRuleKeepsAbsentImageAbsent(t *testing.T) {
	got, err := BuildRowImage(transformationColumns, nil, nil, []Transformation{mustRename(t, "r", "name", "customer_name")})
	if err != nil {
		t.Fatalf("BuildRowImage: %v", err)
	}
	if got != nil {
		t.Fatalf("Bild %q, wollen nil", got)
	}
}

// Eine leere Regelmenge liefert dieselben Bytes wie ein Aufruf ohne Regeln:
// dieselbe Referenz-Tabelle wie `TestBuildRowImageBytes`, mit einer leeren,
// nicht-nil Liste.
func TestBuildRowImageEmptyRuleSetKeepsBytes(t *testing.T) {
	for _, c := range rowImageByteCases {
		t.Run(c.name, func(t *testing.T) {
			got, err := BuildRowImage(rowImageColumns, c.values, c.excluded, []Transformation{})
			if err != nil {
				t.Fatalf("BuildRowImage: %v", err)
			}
			if string(got) != c.want {
				t.Fatalf("Bild %q, wollen %q", got, c.want)
			}
		})
	}
}

// Die Auswertung ist rein: die Regelliste bleibt unverändert, das Ergebnis
// teilt keinen Speicher mit ihr, und gleichzeitige Aufrufe (`-race`)
// liefern dieselben Bytes.
func TestBuildRowImageWithRulesPureAndConcurrent(t *testing.T) {
	rules := []Transformation{mustRename(t, "r", "name", "customer_name")}
	values := []*string{textValue("7"), textValue("s"), textValue("Ada"), textValue("o")}
	first, err := BuildRowImage(transformationColumns, values, []string{"secret"}, rules)
	if err != nil {
		t.Fatalf("BuildRowImage: %v", err)
	}
	want := string(first)
	first[0] = '!'
	var wg sync.WaitGroup
	results := make([][]byte, 8)
	for i := range results {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results[i], _ = BuildRowImage(transformationColumns, values, []string{"secret"}, rules)
		}()
	}
	wg.Wait()
	for i, result := range results {
		if string(result) != want {
			t.Fatalf("Aufruf %d: Bild %q, wollen %q", i, result, want)
		}
	}
	if rules[0].Name() != "r" || rules[0].Column() != "name" || rules[0].To() != "customer_name" {
		t.Fatalf("Regel verändert: %+v", rules[0])
	}
}
