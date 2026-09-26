// Command kommentar-kennungen listet die Kommentarblöcke der Go-Dateien, die
// ihre Herkunft nicht als ein auflösbares Feld tragen: Kandidat ist ein Block
// mit mindestens zwei verschiedenen Kennungen (ADR, LH-FA, LH-QA, SPEC, ARC)
// oder mit „ff.“ hinter einer Kennung. Es prüft die Form, nicht die Wahrheit
// eines Kommentars. Vertrag, Grenze und Exit-Codes:
// harness/sensors/kommentar-kennungen.md (AGENTS.md §3.7).
//
// Aufruf: kommentar-kennungen [-count] [-tests all|exclude|only] [-diff] [Pfad ...]
// Ohne Pfad liest es den Baum ab dem Arbeitsverzeichnis. Mit -diff liest es
// einen `git diff -U0`-Strom von stdin und meldet nur Blöcke, die eine
// hinzugefügte Zeile überlappen.
// Exit: 0 kein Kandidat (oder -count) · 1 mindestens ein Kandidat · 2 Eingabefehler
package main

import (
	"bufio"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Die vier Kennungsarten; die Ziffernbreite je Art steht im Muster.
var idPattern = regexp.MustCompile(`ADR-[0-9]{4}|LH-(?:FA|QA)-[A-Z]{3}-[0-9]{3}|SPEC-[0-9]{3}|ARC-[0-9]{3}`)

// Fortsetzung einer Kompaktform hinter einer Kennung: `/` oder `…`, optional
// in Backticks, dann die Nummer der nächsten Kennung derselben Familie.
var continuation = regexp.MustCompile("^`?(?:/|…)`?([0-9]{3,4})")

// „ff.“ hinter einer Kennung, optional mit schließenden Zeichen dazwischen.
var followingPages = regexp.MustCompile("^[`)\\]]*\\s+ff\\.")

// excludedRoots sind die Wurzelverzeichnisse, die das Werkzeug nicht liest:
// erzeugter Code, die SDK-Bäume (dort gilt die strengere Regel von
// `make sdk-public-doc-check`), das Repository-Innere und die vendored Baseline.
var excludedRoots = map[string]bool{"gen": true, "sdks": true, ".git": true, ".harness": true}

// block ist ein Kommentarblock mit Zeilenbereich und den gezählten Kennungen.
type block struct {
	file       string
	start, end int
	ids        []string
	ff         bool
}

// candidate meldet, ob der Block die Form verletzt: zwei verschiedene
// Kennungen oder „ff.“.
func (b block) candidate() bool { return len(b.ids) >= 2 || b.ff }

func (b block) String() string {
	parts := append([]string(nil), b.ids...)
	if b.ff {
		parts = append(parts, "ff.")
	}
	return fmt.Sprintf("%s:%d-%d  %s", b.file, b.start, b.end, strings.Join(parts, ", "))
}

// scanText liefert die verschiedenen Kennungen eines Kommentartexts in der
// Reihenfolge des ersten Auftretens und ob ein „ff.“ hinter einer Kennung
// steht. Eine Kompaktform zählt je Nummer als eine Kennung.
func scanText(text string) (ids []string, ff bool) {
	seen := map[string]bool{}
	add := func(id string) {
		if !seen[id] {
			seen[id] = true
			ids = append(ids, id)
		}
	}
	for _, loc := range idPattern.FindAllStringIndex(text, -1) {
		id := text[loc[0]:loc[1]]
		add(id)
		dash := strings.LastIndex(id, "-")
		family, width := id[:dash+1], len(id)-dash-1
		rest := text[loc[1]:]
		for {
			m := continuation.FindStringSubmatchIndex(rest)
			if m == nil || m[3]-m[2] != width {
				break
			}
			add(family + rest[m[2]:m[3]])
			rest = rest[m[1]:]
		}
		if followingPages.MatchString(rest) {
			ff = true
		}
	}
	return ids, ff
}

// commentText entfernt die Kommentar-Begrenzer.
func commentText(c *ast.Comment) string {
	switch {
	case strings.HasPrefix(c.Text, "//"):
		return c.Text[2:]
	case strings.HasPrefix(c.Text, "/*"):
		return strings.TrimSuffix(c.Text[2:], "*/")
	}
	return c.Text
}

// blocksOf liest die Kommentargruppen einer Quelle. Eine Gruppe ist eine
// Folge von Kommentaren ohne Leerzeile und ohne Code dazwischen; ein
// Endkommentar hinter Code bildet eine eigene Gruppe. Direktiven (`//go:…`)
// zählen nicht mit, eine Kennung in einem Zeichenketten-Literal ist kein
// Kommentar und wird nicht gelesen.
func blocksOf(name string, src []byte) ([]block, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, name, src, parser.ParseComments|parser.SkipObjectResolution)
	if err != nil {
		return nil, err
	}
	var out []block
	for _, g := range f.Comments {
		var texts []string
		start, end := 0, 0
		for _, c := range g.List {
			if strings.HasPrefix(c.Text, "//go:") {
				continue
			}
			texts = append(texts, commentText(c))
			line := fset.Position(c.Pos()).Line
			endLine := fset.Position(c.End()).Line
			if start == 0 {
				start = line
			}
			end = endLine
		}
		if len(texts) == 0 {
			continue
		}
		ids, ff := scanText(strings.Join(texts, " "))
		if len(ids) == 0 {
			continue
		}
		out = append(out, block{file: name, start: start, end: end, ids: ids, ff: ff})
	}
	return out, nil
}

// lineRange ist ein geschlossener Zeilenbereich.
type lineRange struct{ from, to int }

var hunkHeader = regexp.MustCompile(`^@@ -[0-9]+(?:,[0-9]+)? \+([0-9]+)(?:,([0-9]+))? @@`)

// parseDiff liest einen `git diff -U0`-Strom und liefert je Datei die
// hinzugefügten Zeilenbereiche. Ein Hunk mit null hinzugefügten Zeilen (reine
// Löschung) trägt keinen Bereich.
func parseDiff(r io.Reader) (map[string][]lineRange, error) {
	added := map[string][]lineRange{}
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)
	file := ""
	for sc.Scan() {
		line := sc.Text()
		switch {
		case strings.HasPrefix(line, "+++ "):
			path := strings.TrimPrefix(line, "+++ ")
			if path == "/dev/null" {
				file = ""
			} else {
				file = filepath.ToSlash(strings.TrimPrefix(path, "b/"))
			}
		case strings.HasPrefix(line, "@@ "):
			m := hunkHeader.FindStringSubmatch(line)
			if m == nil {
				return nil, fmt.Errorf("Hunk-Kopf nicht lesbar: %q", line)
			}
			if file == "" {
				continue
			}
			from, _ := strconv.Atoi(m[1])
			n := 1
			if m[2] != "" {
				n, _ = strconv.Atoi(m[2])
			}
			if n > 0 {
				added[file] = append(added[file], lineRange{from, from + n - 1})
			}
		}
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return added, nil
}

// overlaps meldet, ob der Block einen der Bereiche schneidet.
func (b block) overlaps(rs []lineRange) bool {
	for _, r := range rs {
		if b.start <= r.to && r.from <= b.end {
			return true
		}
	}
	return false
}

// goFiles sammelt die `.go`-Dateien unter den Pfaden, sortiert und ohne die
// ausgenommenen Wurzeln (relativ zum Arbeitsverzeichnis).
func goFiles(paths []string) ([]string, error) {
	seen := map[string]bool{}
	var files []string
	for _, p := range paths {
		info, err := os.Stat(p)
		if err != nil {
			return nil, err
		}
		if !info.IsDir() {
			if strings.HasSuffix(p, ".go") && !seen[p] && !excluded(p) {
				seen[p] = true
				files = append(files, p)
			}
			continue
		}
		err = filepath.WalkDir(p, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if excluded(path) {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			if !d.IsDir() && strings.HasSuffix(path, ".go") && !seen[path] {
				seen[path] = true
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	sort.Strings(files)
	return files, nil
}

// excluded prüft die erste Pfadkomponente relativ zum Arbeitsverzeichnis.
func excluded(path string) bool {
	rel := filepath.ToSlash(filepath.Clean(path))
	first := strings.SplitN(rel, "/", 2)[0]
	return excludedRoots[first]
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("kommentar-kennungen", flag.ContinueOnError)
	flags.SetOutput(stderr)
	count := flags.Bool("count", false, "nur die Zahl der Kandidaten drucken (Exit 0)")
	tests := flags.String("tests", "all", "all, exclude (ohne *_test.go) oder only (nur *_test.go)")
	diff := flags.Bool("diff", false, "git-diff-U0-Strom von stdin lesen, nur hinzugefügte Zeilen betrachten")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if *tests != "all" && *tests != "exclude" && *tests != "only" {
		fmt.Fprintf(stderr, "kommentar-kennungen: -tests '%s' ist weder all, exclude noch only\n", *tests)
		return 2
	}
	paths := flags.Args()
	if len(paths) == 0 {
		paths = []string{"."}
	}
	var added map[string][]lineRange
	if *diff {
		var err error
		if added, err = parseDiff(stdin); err != nil {
			fmt.Fprintf(stderr, "kommentar-kennungen: Diff-Strom: %v\n", err)
			return 2
		}
	}
	files, err := goFiles(paths)
	if err != nil {
		fmt.Fprintf(stderr, "kommentar-kennungen: %v\n", err)
		return 2
	}
	var found []block
	for _, name := range files {
		isTest := strings.HasSuffix(name, "_test.go")
		if (*tests == "exclude" && isTest) || (*tests == "only" && !isTest) {
			continue
		}
		src, err := os.ReadFile(name)
		if err != nil {
			fmt.Fprintf(stderr, "kommentar-kennungen: %v\n", err)
			return 2
		}
		blocks, err := blocksOf(filepath.ToSlash(filepath.Clean(name)), src)
		if err != nil {
			fmt.Fprintf(stderr, "kommentar-kennungen: %v\n", err)
			return 2
		}
		for _, b := range blocks {
			if !b.candidate() || (*diff && !b.overlaps(added[b.file])) {
				continue
			}
			found = append(found, b)
		}
	}
	if *count {
		fmt.Fprintln(stdout, len(found))
		return 0
	}
	for _, b := range found {
		fmt.Fprintln(stdout, b)
	}
	if len(found) > 0 {
		return 1
	}
	return 0
}
