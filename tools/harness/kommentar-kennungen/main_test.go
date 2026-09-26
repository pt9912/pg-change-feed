package main

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// TestScanText bindet die Zählregel an ihre Eingabe: verschiedene Kennungen
// zählen einzeln, dieselbe Kennung einmal, eine Kompaktform je Nummer, „ff.“
// nur hinter einer Kennung.
func TestScanText(t *testing.T) {
	cases := []struct {
		name string
		text string
		ids  []string
		ff   bool
	}{
		{"eine Kennung", "verteilt an alle (`ADR-0060`)", []string{"ADR-0060"}, false},
		{"zwei verschiedene Kennungen", "ADR-0060 und LH-FA-SST-008", []string{"ADR-0060", "LH-FA-SST-008"}, false},
		{"dieselbe Kennung zweimal", "ADR-0060, siehe ADR-0060", []string{"ADR-0060"}, false},
		{"alle vier Arten", "ADR-0001 LH-FA-CAP-001 LH-QA-SEC-001 SPEC-003 ARC-002",
			[]string{"ADR-0001", "LH-FA-CAP-001", "LH-QA-SEC-001", "SPEC-003", "ARC-002"}, false},
		{"Kompaktform mit Schrägstrich", "(`LH-FA-CFG-001`/`002`/`003`)",
			[]string{"LH-FA-CFG-001", "LH-FA-CFG-002", "LH-FA-CFG-003"}, false},
		{"Kompaktform ohne Backticks", "LH-QA-SEC-001/002", []string{"LH-QA-SEC-001", "LH-QA-SEC-002"}, false},
		{"Bereich mit Auslassungszeichen", "(`LH-FA-RET-002`…`004`)", []string{"LH-FA-RET-002", "LH-FA-RET-004"}, false},
		{"Kompaktform einer ADR", "ADR-0060/0066", []string{"ADR-0060", "ADR-0066"}, false},
		{"Nummer mit falscher Breite ist keine Kennung", "ADR-0060/005", []string{"ADR-0060"}, false},
		{"Nummer mit zu vielen Ziffern ist keine Kennung", "LH-FA-CFG-001/0021", []string{"LH-FA-CFG-001"}, false},
		{"ff. hinter einer Kennung", "(`LH-FA-REA-001` ff.)", []string{"LH-FA-REA-001"}, true},
		{"ff. ohne Backticks", "LH-FA-REA-001 ff. weiter", []string{"LH-FA-REA-001"}, true},
		{"ff. ohne Kennung davor", "siehe ff. und ADR-0001", []string{"ADR-0001"}, false},
		{"ff. nach dem Wort dazwischen", "LH-FA-REA-001 und ff.", []string{"LH-FA-REA-001"}, false},
		{"Unterkennung zählt als eine", "LH-FA-CAP-006.a", []string{"LH-FA-CAP-006"}, false},
		{"keine Kennung", "verteilt an alle Abonnenten", nil, false},
		{"unvollständige Kennungen", "ADR-60 LH-FA-CAP SPEC- ARC-1", nil, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ids, ff := scanText(c.text)
			if !reflect.DeepEqual(ids, c.ids) || ff != c.ff {
				t.Fatalf("scanText(%q) = %v, %v; erwartet %v, %v", c.text, ids, ff, c.ids, c.ff)
			}
		})
	}
}

// TestBlocksOf bindet die Blockgrenzen an die Quelle: Zeilenbereich, Trennung
// durch Leerzeile, Endkommentar hinter Code, Direktive, Zeichenketten-Literal.
func TestBlocksOf(t *testing.T) {
	cases := []struct {
		name string
		src  string
		want []block
	}{
		{"zwei Kennungen in einem Block",
			"package p\n\n// Erste ADR-0001\n// zweite LH-FA-CAP-001\nfunc F() {}\n",
			[]block{{file: "p.go", start: 3, end: 4, ids: []string{"ADR-0001", "LH-FA-CAP-001"}}}},
		{"Leerzeile trennt die Blöcke",
			"package p\n\n// Erste ADR-0001\n\n// zweite LH-FA-CAP-001\nfunc F() {}\n",
			[]block{
				{file: "p.go", start: 3, end: 3, ids: []string{"ADR-0001"}},
				{file: "p.go", start: 5, end: 5, ids: []string{"LH-FA-CAP-001"}},
			}},
		{"Endkommentar hinter Code bildet einen eigenen Block",
			"package p\n\nvar x = 1 // ADR-0001 LH-FA-CAP-001\n// SPEC-003\nvar y = 2\n",
			[]block{
				{file: "p.go", start: 3, end: 3, ids: []string{"ADR-0001", "LH-FA-CAP-001"}},
				{file: "p.go", start: 4, end: 4, ids: []string{"SPEC-003"}},
			}},
		{"Kennung in einem Zeichenketten-Literal ist kein Kommentar",
			"package p\n\nvar s = \"ADR-0001 LH-FA-CAP-001\"\nvar r = `// ADR-0001 LH-FA-CAP-001`\n",
			nil},
		{"Direktive zählt nicht mit",
			"package p\n\n// ADR-0002\n//go:generate tool ADR-0001 LH-FA-CAP-001\nfunc F() {}\n",
			[]block{{file: "p.go", start: 3, end: 3, ids: []string{"ADR-0002"}}}},
		{"Direktive allein ist kein Block",
			"//go:build linux && ADR-0001 && LH-FA-CAP-001\n\npackage p\n",
			nil},
		{"Blockkommentar über mehrere Zeilen",
			"package p\n\n/* ADR-0001\n   LH-FA-CAP-001 */\nfunc F() {}\n",
			[]block{{file: "p.go", start: 3, end: 4, ids: []string{"ADR-0001", "LH-FA-CAP-001"}}}},
		{"ff. über einen Zeilenumbruch",
			"package p\n\n// (`LH-FA-REA-001`\n// ff.)\nfunc F() {}\n",
			[]block{{file: "p.go", start: 3, end: 4, ids: []string{"LH-FA-REA-001"}, ff: true}}},
		{"Kommentar ohne Kennung",
			"package p\n\n// nichts\nfunc F() {}\n",
			nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := blocksOf("p.go", []byte(c.src))
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("blocksOf = %+v; erwartet %+v", got, c.want)
			}
		})
	}
	if _, err := blocksOf("p.go", []byte("package p\nfunc (")); err == nil {
		t.Fatal("ein unlesbarer Quelltext liefert keinen Fehler")
	}
}

func TestCandidate(t *testing.T) {
	cases := []struct {
		name string
		b    block
		want bool
	}{
		{"eine Kennung", block{ids: []string{"ADR-0001"}}, false},
		{"zwei Kennungen", block{ids: []string{"ADR-0001", "SPEC-003"}}, true},
		{"eine Kennung und ff.", block{ids: []string{"LH-FA-REA-001"}, ff: true}, true},
	}
	for _, c := range cases {
		if got := c.b.candidate(); got != c.want {
			t.Errorf("%s: candidate() = %v, erwartet %v", c.name, got, c.want)
		}
	}
}

func TestParseDiff(t *testing.T) {
	diff := strings.Join([]string{
		"diff --git a/a.go b/a.go",
		"--- a/a.go",
		"+++ b/a.go",
		"@@ -3 +3,2 @@ func F() {",
		"@@ -10,2 +11 @@",
		"@@ -20,3 +21,0 @@",
		"diff --git a/gone.go b/gone.go",
		"--- a/gone.go",
		"+++ /dev/null",
		"@@ -1,3 +0,0 @@",
		"diff --git a/d/b.go b/d/b.go",
		"--- /dev/null",
		"+++ b/d/b.go",
		"@@ -0,0 +1,4 @@",
	}, "\n") + "\n"
	got, err := parseDiff(strings.NewReader(diff))
	if err != nil {
		t.Fatal(err)
	}
	want := map[string][]lineRange{
		"a.go":   {{3, 4}, {11, 11}},
		"d/b.go": {{1, 4}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseDiff = %v; erwartet %v", got, want)
	}
	if _, err := parseDiff(strings.NewReader("+++ b/a.go\n@@ kaputt @@\n")); err == nil {
		t.Fatal("ein unlesbarer Hunk-Kopf liefert keinen Fehler")
	}
}

func TestOverlaps(t *testing.T) {
	b := block{start: 5, end: 8}
	cases := []struct {
		name string
		rs   []lineRange
		want bool
	}{
		{"Bereich innerhalb", []lineRange{{6, 6}}, true},
		{"Bereich schneidet den Anfang", []lineRange{{3, 5}}, true},
		{"Bereich schneidet das Ende", []lineRange{{8, 10}}, true},
		{"Bereich davor", []lineRange{{1, 4}}, false},
		{"Bereich dahinter", []lineRange{{9, 12}}, false},
		{"kein Bereich", nil, false},
	}
	for _, c := range cases {
		if got := b.overlaps(c.rs); got != c.want {
			t.Errorf("%s: overlaps = %v, erwartet %v", c.name, got, c.want)
		}
	}
}

// writeTree legt die Dateien unter dem Arbeitsverzeichnis des Tests an.
func writeTree(t *testing.T, files map[string]string) {
	t.Helper()
	t.Chdir(t.TempDir())
	for name, src := range files {
		if err := os.MkdirAll(filepath.Dir(name), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(name, []byte(src), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

const (
	candidateSrc = "package p\n\n// ADR-0001 und LH-FA-CAP-001\nfunc F() {}\n"
	okSrc        = "package p\n\n// ADR-0001\nfunc F() {}\n"
)

func runTool(t *testing.T, stdin string, args ...string) (code int, out, errOut string) {
	t.Helper()
	var so, se bytes.Buffer
	code = run(args, strings.NewReader(stdin), &so, &se)
	return code, so.String(), se.String()
}

func TestRunModes(t *testing.T) {
	writeTree(t, map[string]string{
		"a/x.go":        candidateSrc,
		"a/y.go":        candidateSrc,
		"a/x_test.go":   candidateSrc,
		"a/ok.go":       okSrc,
		"gen/g.go":      candidateSrc,
		"sdks/s.go":     candidateSrc,
		".harness/h.go": candidateSrc,
		"a/notgo.txt":   "// ADR-0001 LH-FA-CAP-001\n",
	})
	cases := []struct {
		name     string
		args     []string
		wantCode int
		wantOut  string
	}{
		{"Baum, alle Dateien", nil, 1,
			"a/x.go:3-3  ADR-0001, LH-FA-CAP-001\na/x_test.go:3-3  ADR-0001, LH-FA-CAP-001\na/y.go:3-3  ADR-0001, LH-FA-CAP-001\n"},
		{"ohne Testdateien", []string{"-tests", "exclude"}, 1,
			"a/x.go:3-3  ADR-0001, LH-FA-CAP-001\na/y.go:3-3  ADR-0001, LH-FA-CAP-001\n"},
		{"nur Testdateien", []string{"-tests", "only"}, 1, "a/x_test.go:3-3  ADR-0001, LH-FA-CAP-001\n"},
		{"Zahl gesamt", []string{"-count"}, 0, "3\n"},
		{"Zahl ohne Testdateien", []string{"-count", "-tests", "exclude"}, 0, "2\n"},
		{"Zahl nur Testdateien", []string{"-count", "-tests", "only"}, 0, "1\n"},
		{"Pfad ohne Kandidat", []string{"a/ok.go"}, 0, ""},
		{"Pfad eines Verzeichnisses", []string{"a"}, 1,
			"a/x.go:3-3  ADR-0001, LH-FA-CAP-001\na/x_test.go:3-3  ADR-0001, LH-FA-CAP-001\na/y.go:3-3  ADR-0001, LH-FA-CAP-001\n"},
		{"ausgenommene Wurzel gen", []string{"gen"}, 0, ""},
		{"ausgenommene Wurzel sdks", []string{"sdks/s.go"}, 0, ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			code, out, errOut := runTool(t, "", c.args...)
			if code != c.wantCode || out != c.wantOut {
				t.Fatalf("Exit %d, Ausgabe %q (stderr %q); erwartet Exit %d, Ausgabe %q",
					code, out, errOut, c.wantCode, c.wantOut)
			}
		})
	}
}

func TestRunDiffMode(t *testing.T) {
	writeTree(t, map[string]string{"a/x.go": candidateSrc + "\n// SPEC-003 und ARC-002\nvar v = 1\n"})
	overlapping := "+++ b/a/x.go\n@@ -3,0 +3,1 @@\n"
	touchingOther := "+++ b/a/x.go\n@@ -9,0 +9,1 @@\n"
	otherFile := "+++ b/a/y.go\n@@ -3,0 +3,1 @@\n"
	cases := []struct {
		name     string
		diff     string
		args     []string
		wantCode int
		wantOut  string
	}{
		{"hinzugefügte Zeile im ersten Block", overlapping, []string{"-diff"}, 1,
			"a/x.go:3-3  ADR-0001, LH-FA-CAP-001\n"},
		{"hinzugefügte Zeile im zweiten Block", "+++ b/a/x.go\n@@ -1 +6 @@\n", []string{"-diff"}, 1,
			"a/x.go:6-6  SPEC-003, ARC-002\n"},
		{"hinzugefügte Zeile außerhalb der Blöcke", touchingOther, []string{"-diff"}, 0, ""},
		{"reine Löschung im Block", "+++ b/a/x.go\n@@ -3,1 +2,0 @@\n", []string{"-diff"}, 0, ""},
		{"hinzugefügte Zeile in einer anderen Datei", otherFile, []string{"-diff"}, 0, ""},
		{"leerer Diff", "", []string{"-diff"}, 0, ""},
		{"Zahl im Diff-Modus", overlapping, []string{"-diff", "-count"}, 0, "1\n"},
		{"ohne -diff bleibt der Diff unbeachtet", touchingOther, nil, 1,
			"a/x.go:3-3  ADR-0001, LH-FA-CAP-001\na/x.go:6-6  SPEC-003, ARC-002\n"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			code, out, errOut := runTool(t, c.diff, c.args...)
			if code != c.wantCode || out != c.wantOut {
				t.Fatalf("Exit %d, Ausgabe %q (stderr %q); erwartet Exit %d, Ausgabe %q",
					code, out, errOut, c.wantCode, c.wantOut)
			}
		})
	}
}

func TestRunInputErrors(t *testing.T) {
	writeTree(t, map[string]string{"a/x.go": candidateSrc, "a/broken.go": "package p\nfunc ("})
	cases := []struct {
		name  string
		stdin string
		args  []string
	}{
		{"unbekannter Wert für -tests", "", []string{"-tests", "alle"}},
		{"unbekannte Option", "", []string{"-nope"}},
		{"Pfad fehlt", "", []string{"gibt-es-nicht"}},
		{"Quelltext nicht lesbar", "", []string{"a/broken.go"}},
		{"Diff-Strom nicht lesbar", "+++ b/a/x.go\n@@ kaputt @@\n", []string{"-diff", "a/x.go"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			code, out, errOut := runTool(t, c.stdin, c.args...)
			if code != 2 || out != "" || errOut == "" {
				t.Fatalf("Exit %d, Ausgabe %q, stderr %q; erwartet Exit 2, keine Ausgabe, eine Meldung", code, out, errOut)
			}
		})
	}
}
