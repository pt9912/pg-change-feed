package integration_test

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"testing"
)

// runnerSkriptPfad ist der Pfad des Integrations-Runners relativ zur
// Repo-Wurzel.
const runnerSkriptPfad = "tools/harness/run-integration-tests.sh"

// runArgumentMuster findet ein `-run`-Argument samt Wert in Einfach- oder
// Doppelanführungszeichen (`-run '…'`, `-run "…"`, `-run='…'`); die dritte
// Gruppe fängt einen Wert ohne Anführungszeichen, den der Leser ablehnt.
var runArgumentMuster = regexp.MustCompile(`(?:^|\s)-run(?:=|\s+)(?:'([^']*)'|"([^"]*)"|(\S*))`)

// TestRunnerFuehrtJedeE2EFunktionAus ist die Vollständigkeits-Hälfte neben
// TestAbdeckungstabelleZeilen: jede `func TestE2E*` dieses Pakets steht in
// mindestens einem `-run`-Argumentwert von tools/harness/run-integration-tests.sh
// (`LH-QA-POR-003`). `go test -run` meldet keinen Fehler, solange ein anderer
// Name desselben Aufrufs trifft; eine Funktion außerhalb aller Muster liefe nie.
//
// Der Leser liest nur die `-run`-Werte (Anführungszeichen-Formen, siehe
// runArgumentMuster); ein Kommentar zählt nicht: ab einem `#` außerhalb von
// Einfach- und Doppelanführungszeichen, das am Zeilenanfang oder nach Leerraum
// steht, ist der Rest der Zeile entfernt (`$#`, `${#x}`, `a#b` sind keiner).
// Er prüft die Anwesenheit im Muster, nicht dass die Phase erreicht wird oder
// die Funktion ohne `--- SKIP` läuft. Benannte Grenze: der Leser parst keine
// Shell-Grammatik und liest ein `-run` in einem String (`echo "… -run '…'"`)
// oder Heredoc-Text wie ein Argument mit. Der Test braucht keinen Stack und
// keine Datenbank: er liest Quelltext und Skript aus dem Repo-Mount (`/src` in
// `make test`) und läuft deshalb auch dort.
func TestRunnerFuehrtJedeE2EFunktionAus(t *testing.T) {
	_, quelle, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Runner-Vollständigkeit: Quelldatei dieses Pakets nicht auflösbar")
	}
	verzeichnis := filepath.Dir(quelle)
	namen, err := e2eFunktionsnamen(verzeichnis)
	if err != nil {
		t.Fatalf("Runner-Vollständigkeit: %v", err)
	}
	if len(namen) == 0 {
		t.Fatal("Runner-Vollständigkeit: kein func TestE2E* im Paket gefunden")
	}
	skript, err := os.ReadFile(filepath.Join(verzeichnis, "..", "..", filepath.FromSlash(runnerSkriptPfad)))
	if err != nil {
		t.Fatalf("Runner-Vollständigkeit: %v", err)
	}
	fehlend, err := fehlendeImRunner(string(skript), namen)
	if err != nil {
		t.Fatalf("Runner-Vollständigkeit: %s: %v", runnerSkriptPfad, err)
	}
	t.Logf("Runner-Vollständigkeit: %d von %d func TestE2E* in einem -run-Wert von %s erfasst",
		len(namen)-len(fehlend), len(namen), runnerSkriptPfad)
	if len(fehlend) > 0 {
		t.Fatalf("von keinem -run-Wert in %s erfasst: %s — die Funktionen liefen nie",
			runnerSkriptPfad, strings.Join(fehlend, ", "))
	}
}

// e2eFunktionsnamen liefert die Namen aller `func TestE2E*` der Go-Dateien
// eines Verzeichnisses, sortiert.
func e2eFunktionsnamen(verzeichnis string) ([]string, error) {
	eintraege, err := os.ReadDir(verzeichnis)
	if err != nil {
		return nil, err
	}
	namen := make([]string, 0)
	for _, eintrag := range eintraege {
		if eintrag.IsDir() || !strings.HasSuffix(eintrag.Name(), ".go") {
			continue
		}
		geparst, err := parser.ParseFile(token.NewFileSet(), filepath.Join(verzeichnis, eintrag.Name()), nil, 0)
		if err != nil {
			return nil, err
		}
		for _, deklaration := range geparst.Decls {
			funktion, ok := deklaration.(*ast.FuncDecl)
			if ok && funktion.Recv == nil && strings.HasPrefix(funktion.Name.Name, "TestE2E") {
				namen = append(namen, funktion.Name.Name)
			}
		}
	}
	sort.Strings(namen)
	return namen, nil
}

// ohneKommentar entfernt einen Shell-Kommentar aus einer Zeile: ab dem ersten
// `#` außerhalb von Anführungszeichen, das am Zeilenanfang oder nach Leerraum
// steht. Ein Backslash außerhalb von Einfachanführungszeichen schützt das
// folgende Zeichen.
func ohneKommentar(zeile string) string {
	einfach, doppelt := false, false
	for i := 0; i < len(zeile); i++ {
		c := zeile[i]
		switch {
		case einfach:
			if c == '\'' {
				einfach = false
			}
		case c == '\\':
			i++
		case doppelt:
			if c == '"' {
				doppelt = false
			}
		case c == '\'':
			einfach = true
		case c == '"':
			doppelt = true
		case c == '#' && (i == 0 || zeile[i-1] == ' ' || zeile[i-1] == '\t'):
			return zeile[:i]
		}
	}
	return zeile
}

// runMuster liest die `-run`-Werte eines Skripttexts als Muster, ohne
// Kommentare (ohneKommentar). Ein `-run` ohne Anführungszeichen oder mit einem
// ungültigen Muster ist ein Fehler.
func runMuster(skript string) ([]*regexp.Regexp, error) {
	muster := make([]*regexp.Regexp, 0)
	for nummer, zeile := range strings.Split(skript, "\n") {
		zeile = ohneKommentar(zeile)
		for _, treffer := range runArgumentMuster.FindAllStringSubmatch(zeile, -1) {
			if treffer[3] != "" || (treffer[1] == "" && treffer[2] == "") {
				return nil, fmt.Errorf("Zeile %d: -run ohne Anführungszeichen oder ohne Wert (%q) — Form nicht lesbar", nummer+1, strings.TrimSpace(treffer[0]))
			}
			wert := treffer[1] + treffer[2]
			kompiliert, err := regexp.Compile(wert)
			if err != nil {
				return nil, fmt.Errorf("Zeile %d: -run-Wert %q ist kein Muster: %v", nummer+1, wert, err)
			}
			muster = append(muster, kompiliert)
		}
	}
	return muster, nil
}

// fehlendeImRunner liefert die Namen, die kein `-run`-Muster des Skripttexts
// trifft (Trefferregel wie `go test -run`: Teilübereinstimmung, ein
// unverankerter Wert wie `TestE2E` erfasst alle Namen mit diesem Teilstring;
// Reihenfolge der Läufe liest der Wächter nicht), sortiert. Ein Skript ohne
// jedes `-run` ist ein Fehler.
func fehlendeImRunner(skript string, namen []string) ([]string, error) {
	muster, err := runMuster(skript)
	if err != nil {
		return nil, err
	}
	if len(muster) == 0 {
		return nil, fmt.Errorf("kein -run-Argument gefunden")
	}
	fehlend := make([]string, 0)
	for _, name := range namen {
		erfasst := false
		for _, m := range muster {
			if m.MatchString(name) {
				erfasst = true
				break
			}
		}
		if !erfasst {
			fehlend = append(fehlend, name)
		}
	}
	sort.Strings(fehlend)
	return fehlend, nil
}

// TestRunnerLeserDreiZustaende fährt den Leser am Skripttext: die drei
// Kernzustände (vollständig, ein Name aus dem Muster entfernt, derselbe Name nur
// noch im Kommentar) und weitere Fälle zu Quoting, Kommentar-Erkennung, den
// benannten Grenzen und den Fehlerfällen. Ein Fall vergleicht die fehlende
// Namensliste mit DeepEqual, ein Fehlerfall den Fehlertext. Fail-open bleiben
// Formen außerhalb des Zustandsautomaten (`;#`, ANSI-C-Quotes `$'…'`, ein
// Anführungszeichen über Zeilengrenzen); sie sind nicht gebunden.
func TestRunnerLeserDreiZustaende(t *testing.T) {
	namen := []string{"TestE2EAlpha", "TestE2EBeta", "TestE2EGamma"}
	vollstaendig := "go test -v \\\n  -run '^(TestE2EAlpha|TestE2EBeta)$' \\\n  ./x\n" +
		"go test -v -run \"^TestE2EGamma$\" ./x\n"
	faelle := []struct {
		name       string
		skript     string
		fehlend    []string
		fehlerText string
	}{
		{"vollständig", vollstaendig, []string{}, ""},
		{
			"Name aus dem Muster entfernt",
			strings.Replace(vollstaendig, "|TestE2EBeta", "", 1),
			[]string{"TestE2EBeta"}, "",
		},
		{
			"Name nur im Zeilenkommentar",
			strings.Replace(vollstaendig, "|TestE2EBeta", "", 1) +
				"# TestE2EBeta läuft separat\n  # go test -run '^TestE2EBeta$'\n",
			[]string{"TestE2EBeta"}, "",
		},
		{
			"Name nur im Zeilenende-Kommentar",
			strings.Replace(vollstaendig, "|TestE2EBeta", "", 1) +
				"true # -run '^TestE2EBeta$'\n",
			[]string{"TestE2EBeta"}, "",
		},
		{
			"# in Anführungszeichen im Muster bleibt gelesen",
			"go test -run '^TestE2EAlpha$| #|TestE2EBeta|TestE2EGamma' ./x\n",
			[]string{}, "",
		},
		{
			"# in Doppelanführungszeichen vor dem -run beginnt keinen Kommentar",
			strings.Replace(vollstaendig, "|TestE2EBeta", "", 1) +
				"echo \"a #b\" -run '^TestE2EBeta$'\n",
			[]string{}, "",
		},
		{
			"maskiertes \\\" in Doppelanführungszeichen schließt sie nicht",
			strings.Replace(vollstaendig, "|TestE2EBeta", "", 1) +
				"echo \"a \\\" #b\" -run '^TestE2EBeta$'\n",
			[]string{}, "",
		},
		{
			"maskiertes \\\" außerhalb öffnet keine Anführungszeichen, das folgende # ist Kommentar",
			strings.Replace(vollstaendig, "|TestE2EBeta", "", 1) +
				"echo \\\" #x -run '^TestE2EBeta$'\n",
			[]string{"TestE2EBeta"}, "",
		},
		{
			"$# vor dem -run beginnt keinen Kommentar",
			"echo $# && go test -run '^(TestE2EAlpha|TestE2EBeta|TestE2EGamma)$' ./x\n",
			[]string{}, "",
		},
		{
			"-run nur in echo-String zählt als erfasst (benannte Grenze)",
			strings.Replace(vollstaendig, "|TestE2EBeta", "", 1) +
				"echo \"go test -run '^TestE2EBeta$' ./x\"\n",
			[]string{}, "",
		},
		{
			"unverankerter Wert erfasst alle Namen (Semantik von go test -run)",
			"go test -run 'TestE2E' ./x\n",
			[]string{}, "",
		},
		{
			"Name außerhalb eines -run-Arguments",
			strings.Replace(vollstaendig, "|TestE2EBeta", "", 1) + "echo \"TestE2EBeta\"\n",
			[]string{"TestE2EBeta"}, "",
		},
		{
			"alle Namen fehlen",
			"go test -run '^TestE2EAndere$' ./x\n",
			namen, "",
		},
		{
			"-run ohne Anführungszeichen",
			"go test -run TestE2EAlpha ./x\n",
			nil, "ohne Anführungszeichen",
		},
		{
			"kein -run im Skript",
			"echo nichts\n",
			nil, "kein -run-Argument",
		},
	}
	for _, f := range faelle {
		t.Run(f.name, func(t *testing.T) {
			fehlend, err := fehlendeImRunner(f.skript, namen)
			if f.fehlerText != "" {
				if err == nil || !strings.Contains(err.Error(), f.fehlerText) {
					t.Fatalf("Fehler %v, erwartet Text %q", err, f.fehlerText)
				}
				return
			}
			if err != nil {
				t.Fatalf("unerwarteter Fehler: %v", err)
			}
			if !reflect.DeepEqual(fehlend, f.fehlend) {
				t.Fatalf("fehlend = %v, erwartet %v", fehlend, f.fehlend)
			}
		})
	}
}
