package backfill_test

import (
	"context"
	stderrors "errors"
	"fmt"
	"strings"
	"testing"

	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Die Tests dieser Datei tragen die Regelauswertung im Run (`LH-FA-CFG-007`,
// `ADR-0112` Folgepflicht 7, `ADR-0117`): der Rig hält die Tabelle
// `public.orders` mit den Spalten `id`, `name`, `secret`; der Regelstand
// kommt aus `fakeRules`. Die Lesungen des Regelstands eines Runs über die drei
// Blöcke des Rigs sind: Lesung 1 zu Beginn (nach dem Öffnen des Snapshots),
// Lesung 2 bis 4 je ein Block, Lesung 5 unmittelbar vor dem Commit.

const testQualified = "public.orders"

// ruleStates liefert je Lesung einen Stand; jenseits der genannten Lesungen
// gilt der letzte.
func ruleStates(states ...map[string][]model.Transformation) func(int) map[string][]model.Transformation {
	return func(call int) map[string][]model.Transformation {
		if call > len(states) {
			call = len(states)
		}
		return states[call-1]
	}
}

// ruled trägt eine Regel und ihre Wirkung auf die Spalte, auf der sie sitzt:
// den Schlüssel, unter dem die Spalte im Bild steht, und den Wert dort zu
// einem Quellwert.
type ruled struct {
	rule  model.Transformation
	key   string
	value func(source string) string
}

// ruleFor bildet einen Regeltyp der Domäne auf eine anwendbare Regel für
// `column` ab; `mapped` sind die Quellwerte, die `map_value` abbildet (jeder
// andere Quellwert bleibt). Ein Regeltyp ohne Fall bricht den Test ab: die
// Tests dieser Datei zählen die Regeltypen aus `model.TransformationKinds` auf.
func ruleFor(t *testing.T, kind model.TransformationKind, column string, mapped ...string) ruled {
	t.Helper()
	switch kind {
	case model.TransformationRenameColumn:
		return ruled{rename(t, "rule-"+column, column, column+"_renamed"), column + "_renamed", func(source string) string { return source }}
	case model.TransformationMapValue:
		values := make(map[string]string, len(mapped))
		for _, source := range mapped {
			values[source] = "ABGEBILDET-" + source
		}
		value := func(source string) string {
			if target, ok := values[source]; ok {
				return target
			}
			return source
		}
		return ruled{mapValue(t, "rule-"+column, column, values), column, value}
	}
	t.Fatalf("Regeltyp %q ohne Fall in diesem Test", kind)
	return ruled{}
}

func mapValue(t *testing.T, name, column string, values map[string]string) model.Transformation {
	t.Helper()
	rule, err := model.NewMapValue(name, column, values)
	if err != nil {
		t.Fatalf("NewMapValue(%q, %q, %v) = %v", name, column, values, err)
	}
	return rule
}

func completedImages(t *testing.T, r *rig) [][]string {
	t.Helper()
	var images [][]string
	for _, block := range r.writer.blocks {
		var blockImages []string
		for _, change := range block.changes {
			blockImages = append(blockImages, string(change.NewImage))
		}
		images = append(images, blockImages)
	}
	return images
}

// TestExecuteBuildsImagesWithTheRuleSet trägt die Regelauswertung im Run
// (`ADR-0112` Folgepflicht 7): jeder Change des Runs trägt das Bild, das
// `model.BuildRowImage` mit dem Regelstand der Tabelle liefert — für jeden
// Regeltyp der Domäne, an allen drei Blöcken, byte-genau; der Regelstand einer
// anderen Tabelle wirkt nicht; ein nicht abgebildeter Wert von `map_value` (`d`)
// bleibt. Rot färbende Mutationen (je eine): `nil` statt des Regelsatzes an
// `BuildRowImage` (Bild ohne Regel); den Regelstand unter dem Schlüssel einer
// anderen Tabelle lesen (Regel wirkt nicht); die Lesung mit einer anderen
// Quelle rufen (Quellen-Prüfung).
func TestExecuteBuildsImagesWithTheRuleSet(t *testing.T) {
	for _, kind := range model.TransformationKinds() {
		t.Run(string(kind), func(t *testing.T) {
			r := newRig()
			r.exclusion.stateFn = func(int) map[string][]string { return map[string][]string{testQualified: {"secret"}} }
			onName := ruleFor(t, kind, "name", "a", "c", "e")
			state := map[string][]model.Transformation{
				testQualified:  {onName.rule},
				"public.other": {ruleFor(t, kind, "id", "1", "2", "3", "4", "5").rule},
			}
			r.rules.stateFn = ruleStates(state)
			run := mustExecute(t, r)
			if run.Status != model.BackfillRunCompleted {
				t.Fatalf("Run = %+v, will completed", run)
			}
			image := func(id, name string) string {
				if name == "" {
					return `{"id":"` + id + `"}`
				}
				return `{"id":"` + id + `","` + onName.key + `":"` + onName.value(name) + `"}`
			}
			want := [][]string{
				{image("1", "a"), image("2", "")},
				{image("3", "c"), image("4", "d")},
				{image("5", "e")},
			}
			got := completedImages(t, r)
			if fmt.Sprint(got) != fmt.Sprint(want) {
				t.Fatalf("Bilder = %v, will %v", got, want)
			}
			if r.rules.gotSource != testSource {
				t.Fatalf("Regelstand der Quelle %q gelesen, will %q", r.rules.gotSource, testSource)
			}
			if r.rules.calls != 5 {
				t.Fatalf("Lesungen des Regelstands = %d, will 5 (Beginn, drei Blöcke, vor dem Commit)", r.rules.calls)
			}
		})
	}
}

// TestExecuteRulesNeverLeakExcludedColumns trägt `LH-QA-SEC-004` für den
// Backfill-Pfad (`ADR-0112` Teilfrage 5, Fitness Function): für jeden Regeltyp
// der Domäne, jede Spalte als Ziel der Regel und jede ausgeschlossene Spalte
// trägt kein Bild des Runs den Schlüssel der ausgeschlossenen Spalte, ihren
// Zielnamen, ihren Wert oder ihren abgebildeten Wert; die übrigen Spalten
// bleiben im Bild, die Spalte der Regel mit dem Wert, den die Regel ihr gibt.
// Rot färbende Mutation: der Ausschluss wird erst nach der Regel geprüft (der
// Zielname bzw. der abgebildete Wert der ausgeschlossenen Spalte erscheint im
// Bild).
func TestExecuteRulesNeverLeakExcludedColumns(t *testing.T) {
	columns := []string{"id", "name", "secret"}
	source := func(column string, block int) string { return fmt.Sprintf("%s-%d", strings.ToUpper(column), block) }
	for _, kind := range model.TransformationKinds() {
		for _, ruleColumn := range columns {
			for _, excluded := range columns {
				t.Run(fmt.Sprintf("%s/Regel an %s/ausgeschlossen %s", kind, ruleColumn, excluded), func(t *testing.T) {
					r := newRig()
					r.snapshot.blocks = [][][]*string{
						{{str(source("id", 1)), str(source("name", 1)), str(source("secret", 1))}},
						{{str(source("id", 2)), str(source("name", 2)), str(source("secret", 2))}},
					}
					r.exclusion.stateFn = func(int) map[string][]string { return map[string][]string{testQualified: {excluded}} }
					rule := ruleFor(t, kind, ruleColumn, source(ruleColumn, 1), source(ruleColumn, 2))
					r.rules.stateFn = ruleStates(rulesOf(testQualified, rule.rule))
					if run := mustExecute(t, r); run.Status != model.BackfillRunCompleted {
						t.Fatalf("Run = %+v, will completed", run)
					}
					for index, block := range r.writer.blocks {
						for _, change := range block.changes {
							image := string(change.NewImage)
							forbidden := []string{`"` + excluded + `"`, strings.ToUpper(excluded) + "-"}
							if ruleColumn == excluded {
								forbidden = append(forbidden, `"`+rule.key+`"`, rule.value(source(excluded, index+1)))
							}
							for _, item := range forbidden {
								if strings.Contains(image, item) {
									t.Fatalf("Bild %s trägt %s der ausgeschlossenen Spalte %s", image, item, excluded)
								}
							}
							for _, kept := range columns {
								if kept == excluded {
									continue
								}
								key, value := kept, source(kept, index+1)
								if kept == ruleColumn {
									key, value = rule.key, rule.value(value)
								}
								if !strings.Contains(image, `"`+key+`":"`+value+`"`) {
									t.Fatalf("Bild %s trägt die Spalte %s nicht als %q: %q", image, kept, key, value)
								}
							}
						}
					}
				})
			}
		}
	}
}

// TestExecuteInapplicableRuleEndsRunAsSchema trägt `ADR-0117` Festlegung 1
// bis 3: eine Regel, die auf die Spalten des Snapshots nicht anwendbar ist
// (Spalte fehlt, Zielname kollidiert), endet den Run `failed` mit der Klasse
// `schema` und dem Regelnamen samt Spalte im Fehlertext — nachdem der Snapshot
// seine Spalten lieferte und bevor eine Zeile gelesen und die
// Schreibtransaktion geöffnet wird, auch bei einer leeren Tabelle; es entsteht
// keine Change und kein Wecksignal, der Snapshot ist geschlossen. Die Regeln
// einer anderen Tabelle prüft der Run nicht, eine anwendbare Regel lässt ihn
// durchlaufen. Rot färbende Mutationen (je eine): die Prüfung
// `checkRulesApplicable` entfällt (Run `completed`, oder Fehler erst im Block);
// die Prüfung gegen `nil` statt der Snapshot-Spalten (die anwendbare Regel
// scheitert); die Prüfung über den Regelstand aller Tabellen (die fremde
// Regel scheitert); die Prüfung nach `Begin` (Begin-Zähler 1).
func TestExecuteInapplicableRuleEndsRunAsSchema(t *testing.T) {
	cases := []struct {
		name    string
		rules   func(t *testing.T) map[string][]model.Transformation
		empty   bool
		failed  bool
		wantErr error
		wantIn  []string
	}{
		{"Spalte fehlt", func(t *testing.T) map[string][]model.Transformation {
			return rulesOf(testQualified, rename(t, "r-fehlt", "gibt_es_nicht", "x"))
		}, false, true, domainerrors.ErrTransformationColumnMissing, []string{`"r-fehlt"`, `"gibt_es_nicht"`, testQualified}},
		{"Zielname kollidiert mit einer Spalte", func(t *testing.T) map[string][]model.Transformation {
			return rulesOf(testQualified, rename(t, "r-kollidiert", "name", "secret"))
		}, false, true, domainerrors.ErrTransformationTargetCollides, []string{`"r-kollidiert"`, `"name"`, "secret"}},
		{"Spalte fehlt bei leerer Tabelle", func(t *testing.T) map[string][]model.Transformation {
			return rulesOf(testQualified, rename(t, "r-fehlt", "gibt_es_nicht", "x"))
		}, true, true, domainerrors.ErrTransformationColumnMissing, []string{`"r-fehlt"`}},
		{"zweite Regel nicht anwendbar", func(t *testing.T) map[string][]model.Transformation {
			return rulesOf(testQualified, rename(t, "r-ok", "name", "label"), rename(t, "r-fehlt", "gibt_es_nicht", "x"))
		}, false, true, domainerrors.ErrTransformationColumnMissing, []string{`"r-fehlt"`}},
		{"Regel einer anderen Tabelle wird nicht geprüft", func(t *testing.T) map[string][]model.Transformation {
			return rulesOf("public.other", rename(t, "r-fremd", "gibt_es_nicht", "x"))
		}, false, false, nil, nil},
		{"anwendbare Regel", func(t *testing.T) map[string][]model.Transformation {
			return rulesOf(testQualified, rename(t, "r-ok", "name", "label"))
		}, false, false, nil, nil},
		{"map_value: Spalte fehlt", func(t *testing.T) map[string][]model.Transformation {
			return rulesOf(testQualified, mapValue(t, "r-fehlt", "gibt_es_nicht", map[string]string{"a": "b"}))
		}, false, true, domainerrors.ErrTransformationColumnMissing, []string{`"r-fehlt"`, `"gibt_es_nicht"`, testQualified}},
		{"map_value: abgebildeter Wert gleicht einer Spalte, anwendbar", func(t *testing.T) map[string][]model.Transformation {
			return rulesOf(testQualified, mapValue(t, "r-ok", "name", map[string]string{"a": "secret", "id": "name"}))
		}, false, false, nil, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := newRig()
			if tc.empty {
				r.snapshot.blocks = nil
			}
			r.rules.stateFn = ruleStates(tc.rules(t))
			run := mustExecute(t, r)
			if !tc.failed {
				if run.Status != model.BackfillRunCompleted {
					t.Fatalf("Run = %+v, will completed", run)
				}
				return
			}
			if run.Status != model.BackfillRunFailed || !strings.HasPrefix(run.ErrorMessage, "schema: ") {
				t.Fatalf("Run = %+v, will failed mit Klasse schema", run)
			}
			if !strings.Contains(run.ErrorMessage, tc.wantErr.Error()) {
				t.Fatalf("Fehlertext = %q, will die Ursache %q", run.ErrorMessage, tc.wantErr)
			}
			for _, want := range tc.wantIn {
				if !strings.Contains(run.ErrorMessage, want) {
					t.Fatalf("Fehlertext = %q, will %q enthalten", run.ErrorMessage, want)
				}
			}
			if r.trace.count("NextBlock") != 0 || r.trace.count("Begin") != 0 || r.trace.count("Append") != 0 || r.trace.count("Commit") != 0 || r.trace.count("Notify") != 0 {
				t.Fatalf("Zeile gelesen, Transaktion geöffnet oder Change geschrieben: %v", r.trace.events)
			}
			if r.snapshot.closed != 1 || r.writer.rollbacks != 0 || run.RowsCopied != 0 {
				t.Fatalf("Snapshot %d-mal geschlossen, Rollbacks %d, RowsCopied %d", r.snapshot.closed, r.writer.rollbacks, run.RowsCopied)
			}
			if r.rules.calls != 1 {
				t.Fatalf("Lesungen des Regelstands = %d, will 1 (Prüfung einmal je Run, vor dem ersten Block)", r.rules.calls)
			}
			if len(r.runs.finished) != 1 || r.runs.finished[0] != run {
				t.Fatalf("Finish = %+v", r.runs.finished)
			}
		})
	}
}

// TestExecuteRuleTargetCollisionInBlockEndsRunAsSchema trägt den Fall, der an
// der Zeile hängt: zwei Regeln mit demselben Zielnamen sind je für sich auf
// die Spalten anwendbar; tragen beide Quellspalten einer Zeile einen Wert,
// meldet der Bild-Bau `ErrTransformationTargetCollides`, und der Run endet
// atomar `failed` mit der Klasse `schema`: die geöffnete Transaktion ist
// zurückgerollt, nichts ist committet. Rot färbende Mutation: die Abbildung
// von `ErrTransformationTargetCollides` auf `schema` in `classifyError`
// entfernen (Klasse `internal`).
func TestExecuteRuleTargetCollisionInBlockEndsRunAsSchema(t *testing.T) {
	r := newRig()
	r.rules.stateFn = ruleStates(rulesOf(testQualified, rename(t, "r-a", "name", "gleich"), rename(t, "r-b", "secret", "gleich")))
	run := mustExecute(t, r)
	if run.Status != model.BackfillRunFailed || !strings.HasPrefix(run.ErrorMessage, "schema: ") || !strings.Contains(run.ErrorMessage, domainerrors.ErrTransformationTargetCollides.Error()) {
		t.Fatalf("Run = %+v, will failed mit Klasse schema und der Kollision", run)
	}
	if r.writer.begun != 1 || r.writer.rollbacks != 1 || len(r.writer.committed) != 0 || len(r.writer.blocks) != 0 {
		t.Fatalf("Begin=%d Rollback=%d Commits=%d Blöcke=%d, will 1/1/0/0", r.writer.begun, r.writer.rollbacks, len(r.writer.committed), len(r.writer.blocks))
	}
	if r.trace.count("Notify") != 0 {
		t.Fatalf("Wecksignal nach dem Fehler: %v", r.trace.events)
	}
}

// TestExecuteRuleStateChangeEndsRunAsConfiguration trägt `ADR-0117`
// Festlegung 5 und die Fail-closed-Prüfung des Regelstands (`ADR-0111`
// Teilfrage 4): jede Abweichung des Regelstands zwischen der Lesung zu Beginn
// und einer späteren Lesung (Blöcke, unmittelbar vor dem Commit) rollt zurück
// und endet den Run `failed` mit der Klasse `configuration` — auch eine
// Abweichung in einem Zwischenblock, die am Ende wieder gleich ist, und eine
// Regel, die nach dem Beginn gesetzt wird und nicht anwendbar wäre (der
// Zustand wechselt, `schema` bleibt der Nichtanwendbarkeit gegen die Spalten
// vorbehalten). Ein reiner Ordnungsunterschied und eine Doppelung sind keine
// Abweichung. Rot färbende Mutationen (je eine): der Vergleich je Block
// entfällt (die Fälle mit Zwischenabweichung enden `completed`); der Vergleich
// vor dem Commit entfällt (der Fall „erst vor dem Commit“); der Vergleich als
// Längenvergleich statt als Mengenvergleich (die Fälle „ersetzt“); `values`
// beim Bau der Regel in `newMapValue` nicht setzen (die Fälle „map_value
// ersetzt“); die Zuordnung in `encodeValueMap` ohne Sortierung kodieren (der
// Fall „je Lesung neu gebaut“).
func TestExecuteRuleStateChangeEndsRunAsConfiguration(t *testing.T) {
	var (
		none = map[string][]model.Transformation{}
	)
	label := func(t *testing.T) map[string][]model.Transformation {
		return rulesOf(testQualified, rename(t, "r-label", "name", "label"))
	}
	cases := []struct {
		name   string
		states func(t *testing.T) []map[string][]model.Transformation
		failed bool
		begun  int
	}{
		{"Regel ab Block 2 gesetzt", func(t *testing.T) []map[string][]model.Transformation {
			return []map[string][]model.Transformation{none, none, label(t)}
		}, true, 1},
		{"Regel ab Block 1 gesetzt (Lesung 2 weicht ab)", func(t *testing.T) []map[string][]model.Transformation {
			return []map[string][]model.Transformation{none, label(t)}
		}, true, 0},
		{"Zwischenblock weicht ab, am Ende wieder gleich", func(t *testing.T) []map[string][]model.Transformation {
			return []map[string][]model.Transformation{none, none, label(t), none, none}
		}, true, 1},
		{"Regel entfernt, dann wieder gesetzt", func(t *testing.T) []map[string][]model.Transformation {
			return []map[string][]model.Transformation{label(t), label(t), none, label(t), label(t)}
		}, true, 1},
		{"Abweichung erst vor dem Commit", func(t *testing.T) []map[string][]model.Transformation {
			return []map[string][]model.Transformation{none, none, none, none, label(t)}
		}, true, 1},
		{"Regel entfernt vor dem Commit", func(t *testing.T) []map[string][]model.Transformation {
			return []map[string][]model.Transformation{label(t), label(t), label(t), label(t), none}
		}, true, 1},
		{"Regel ersetzt (anderer Zielname)", func(t *testing.T) []map[string][]model.Transformation {
			return []map[string][]model.Transformation{label(t), label(t), rulesOf(testQualified, rename(t, "r-label", "name", "titel"))}
		}, true, 1},
		{"Regel ersetzt (andere Spalte, gleiche Zahl)", func(t *testing.T) []map[string][]model.Transformation {
			return []map[string][]model.Transformation{label(t), label(t), label(t), rulesOf(testQualified, rename(t, "r-label", "secret", "label"))}
		}, true, 1},
		{"nach dem Beginn gesetzte, nicht anwendbare Regel ist ein Zustandswechsel", func(t *testing.T) []map[string][]model.Transformation {
			return []map[string][]model.Transformation{none, none, rulesOf(testQualified, rename(t, "r-fehlt", "gibt_es_nicht", "x"))}
		}, true, 1},
		{"Regel einer anderen Tabelle ändert sich", func(t *testing.T) []map[string][]model.Transformation {
			return []map[string][]model.Transformation{none, rulesOf("public.other", rename(t, "r-fremd", "name", "x"))}
		}, false, 1},
		{"Ordnung und Doppelung sind keine Abweichung", func(t *testing.T) []map[string][]model.Transformation {
			a, b := rename(t, "r-a", "name", "a"), rename(t, "r-b", "secret", "b")
			return []map[string][]model.Transformation{
				rulesOf(testQualified, a, b), rulesOf(testQualified, b, a), rulesOf(testQualified, a, b, a), rulesOf(testQualified, b, a), rulesOf(testQualified, a, b),
			}
		}, false, 1},
		{"map_value ersetzt (andere Zuordnung, gleiche Zahl)", func(t *testing.T) []map[string][]model.Transformation {
			mapped := func(target string) map[string][]model.Transformation {
				return rulesOf(testQualified, mapValue(t, "r-map", "name", map[string]string{"a": target}))
			}
			return []map[string][]model.Transformation{mapped("x"), mapped("x"), mapped("y")}
		}, true, 1},
		{"map_value ersetzt (Zuordnung um ein Paar erweitert)", func(t *testing.T) []map[string][]model.Transformation {
			return []map[string][]model.Transformation{
				rulesOf(testQualified, mapValue(t, "r-map", "name", map[string]string{"a": "x"})),
				rulesOf(testQualified, mapValue(t, "r-map", "name", map[string]string{"a": "x"})),
				rulesOf(testQualified, mapValue(t, "r-map", "name", map[string]string{"a": "x", "b": "y"})),
			}
		}, true, 1},
		{"map_value gleichen Inhalts, je Lesung neu gebaut, ist keine Abweichung", func(t *testing.T) []map[string][]model.Transformation {
			built := func(reverse bool) map[string][]model.Transformation {
				values := map[string]string{}
				keys := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
				for i := range keys {
					key := keys[i]
					if reverse {
						key = keys[len(keys)-1-i]
					}
					values[key] = "v-" + key
				}
				return rulesOf(testQualified, mapValue(t, "r-map", "name", values))
			}
			return []map[string][]model.Transformation{built(false), built(true), built(false), built(true), built(false)}
		}, false, 1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := newRig()
			r.rules.stateFn = ruleStates(tc.states(t)...)
			run := mustExecute(t, r)
			if !tc.failed {
				if run.Status != model.BackfillRunCompleted || len(r.writer.committed) != 1 {
					t.Fatalf("Run = %+v, will completed", run)
				}
				return
			}
			if run.Status != model.BackfillRunFailed || !strings.HasPrefix(run.ErrorMessage, "configuration: ") ||
				!strings.Contains(run.ErrorMessage, domainerrors.ErrTransformationStateChanged.Error()) {
				t.Fatalf("Run = %+v, will failed mit Klasse configuration und dem Zustandswechsel", run)
			}
			if len(r.writer.committed) != 0 || r.trace.count("Commit") != 0 || r.trace.count("Notify") != 0 {
				t.Fatalf("Commit oder Wecksignal trotz Abweichung: %v", r.trace.events)
			}
			if r.writer.begun != tc.begun || r.writer.rollbacks != tc.begun {
				t.Fatalf("Begin=%d Rollback=%d, will %d/%d", r.writer.begun, r.writer.rollbacks, tc.begun, tc.begun)
			}
			if r.snapshot.closed != 1 {
				t.Fatalf("Snapshot %d-mal geschlossen, will 1", r.snapshot.closed)
			}
		})
	}
}

// TestExecuteRuleReadFailure trägt den Lesefehler-Zweig des Regelstands je
// Aufrufstelle (`ADR-0111` Teilfrage 4): ein Stand, der nicht gelesen werden
// kann, ist kein bestätigter Stand. Der n-te Lesefehler (Lesung 1 zu Beginn,
// 2 bis 4 je Block, 5 unmittelbar vor dem Commit) endet den Run `failed` mit
// der Klasse der Ursache (`storage`), ohne Commit und ohne Wecksignal; die
// Lesung bricht den Lauf ab (keine weitere Lesung, kein weiterer Block), die
// Schreibtransaktion ist zurückgerollt, der Snapshot geschlossen. Weil nur der
// n-te Aufruf scheitert, färbt sich der Fall der Aufrufstelle rot, sobald sie
// ihren Fehler verwirft — dann läuft der Run bis zum Commit durch. Rot
// färbende Mutationen (je eine, je Aufrufstelle): den Fehler der Lesung zu
// Beginn, je Block oder vor dem Commit verwerfen (`_ =` mit leerem Stand).
func TestExecuteRuleReadFailure(t *testing.T) {
	cases := []struct {
		name      string
		call      int
		appended  int
		rollbacks int
		rows      int64
	}{
		{"zu Beginn, vor dem ersten Block", 1, 0, 0, 0},
		{"im ersten Block", 2, 0, 0, 0},
		{"im zweiten Block", 3, 1, 1, 2},
		{"im letzten Block", 4, 2, 1, 4},
		{"unmittelbar vor dem Commit", 5, 3, 1, 5},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := newRig()
			r.rules.err = fmt.Errorf("%w: Regelstand nicht lesbar", outbound.ErrStorage)
			r.rules.errCall = tc.call
			run := mustExecute(t, r)
			if run.Status != model.BackfillRunFailed || !strings.HasPrefix(run.ErrorMessage, "storage: ") || !strings.Contains(run.ErrorMessage, "Regelstand nicht lesbar") {
				t.Fatalf("Run = %+v, will failed mit storage-Klasse und Ursache", run)
			}
			if run.RowsCopied != tc.rows {
				t.Fatalf("RowsCopied = %d, will %d (zuletzt festgehaltener Zähler)", run.RowsCopied, tc.rows)
			}
			if r.rules.calls != tc.call {
				t.Fatalf("Regelstand-Lesungen = %d, will %d (Abbruch mit dem Lesefehler)", r.rules.calls, tc.call)
			}
			if len(r.writer.blocks) != tc.appended {
				t.Fatalf("angehängte Blöcke = %d, will %d", len(r.writer.blocks), tc.appended)
			}
			if r.trace.count("Commit") != 0 || len(r.writer.committed) != 0 || r.trace.count("Notify") != 0 {
				t.Fatalf("Commit/Wecksignal trotz Lesefehler: %v", r.trace.events)
			}
			if r.writer.rollbacks != tc.rollbacks {
				t.Fatalf("Rollbacks = %d, will %d", r.writer.rollbacks, tc.rollbacks)
			}
			if r.snapshot.closed != 1 {
				t.Fatalf("Snapshot %d-mal geschlossen, will 1", r.snapshot.closed)
			}
			if len(r.runs.finished) != 1 || r.runs.finished[0] != run {
				t.Fatalf("Finish = %+v", r.runs.finished)
			}
		})
	}
}

// TestExecuteUnreadableAppliedRowEndsRunAsInternal trägt die Klasse einer
// `applied`-Zeile des Regelstands, die die Faltung nicht mehr in eine Regel
// führt: der Fehler ist kein Fehler der Speicherung und keine Nichtanwendbarkeit
// gegen die Spalten, er endet als `internal` — je Aufrufstelle der Lesung (zu
// Beginn, je Block, unmittelbar vor dem Commit). Der Fehler ist der reale
// Fehler von `model.FoldTransformations` an einer nicht lesbaren Regelform.
// Rot färbende Mutation: die Klasse `internal` als Rückfall von
// `classifyError` durch eine andere Klasse ersetzen.
func TestExecuteUnreadableAppliedRowEndsRunAsInternal(t *testing.T) {
	_, foldErr := model.FoldTransformations([]model.TransformationRecord{
		{Kind: model.AdministrationRequestSetTransformation, Name: "r-kaputt", Spec: `{"type":`},
	})
	if foldErr == nil {
		t.Fatal("FoldTransformations meldet zu einer nicht lesbaren Regelform keinen Fehler")
	}
	if stderrors.Is(foldErr, outbound.ErrStorage) {
		t.Fatalf("der Faltungsfehler trägt %v", outbound.ErrStorage)
	}
	for _, call := range []int{1, 2, 4, 5} {
		t.Run(fmt.Sprintf("Lesung %d", call), func(t *testing.T) {
			r := newRig()
			r.rules.err, r.rules.errCall = foldErr, call
			run := mustExecute(t, r)
			if run.Status != model.BackfillRunFailed || !strings.HasPrefix(run.ErrorMessage, "internal: ") || !strings.Contains(run.ErrorMessage, "r-kaputt") {
				t.Fatalf("Run = %+v, will failed mit Klasse internal und dem Regelnamen", run)
			}
			if r.trace.count("Commit") != 0 || len(r.writer.committed) != 0 {
				t.Fatalf("Commit trotz Faltungsfehler: %v", r.trace.events)
			}
		})
	}
}

// TestExecuteRuleFailureIsReportedInTheRunOnly trägt „run-lokal“ (`ADR-0117`
// Festlegung 3): ein Run, der an einer Regel endet, meldet den Fehler im
// Run-Zustand und kehrt ohne Fehler zurück; der Use Case hält weder den
// Heartbeat- noch einen Port des Capture-Pfads
// (`TestPortsCarryNoCapturePathPort`), ein Run-Fehler kann den Erfassungspfad
// deshalb nicht anhalten.
func TestExecuteRuleFailureIsReportedInTheRunOnly(t *testing.T) {
	r := newRig()
	r.rules.stateFn = ruleStates(rulesOf(testQualified, rename(t, "r-fehlt", "gibt_es_nicht", "x")))
	result, err := r.execute(context.Background(), t)
	if err != nil {
		t.Fatalf("Execute = %v, will nil: der Fehler steht im Run-Zustand", err)
	}
	if result.Run.Status != model.BackfillRunFailed {
		t.Fatalf("Run = %+v", result.Run)
	}
	if len(r.runs.finished) != 1 {
		t.Fatalf("Finish-Aufrufe = %d, will 1", len(r.runs.finished))
	}
}
