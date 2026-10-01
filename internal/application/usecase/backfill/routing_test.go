package backfill_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Die Tests dieser Datei tragen das Zustellziel der Backfill-Changes und den
// Routing-Regelstand des Runs (`LH-FA-CFG-008`): der Rig hält die Tabelle
// `public.orders` mit den Spalten `id`, `name`, `secret` und die Zeilen
// 1 bis 5 (Namen `a`, NULL, `c`, `d`, `e`) in drei Blöcken. Die Lesungen des
// Routing-Regelstands eines Runs über die drei Blöcke sind: Lesung 1 zu
// Beginn (nach dem Öffnen des Snapshots), Lesung 2 bis 4 je ein Block,
// Lesung 5 unmittelbar vor dem Commit.

func routeRule(t *testing.T, name, target string, order int64, when *model.RouteCondition) model.RouteRule {
	t.Helper()
	rule, err := model.NewRouteRule(name, target, order, when)
	if err != nil {
		t.Fatalf("NewRouteRule(%q, %q, %d) = %v", name, target, order, err)
	}
	return rule
}

func whenName(value string) *model.RouteCondition {
	return &model.RouteCondition{Column: "name", Equals: value}
}

func routesOf(table string, rules ...model.RouteRule) map[string][]model.RouteRule {
	return map[string][]model.RouteRule{table: rules}
}

// routeStates liefert je Lesung einen Stand; jenseits der genannten Lesungen
// gilt der letzte.
func routeStates(states ...map[string][]model.RouteRule) func(int) map[string][]model.RouteRule {
	return func(call int) map[string][]model.RouteRule {
		if call > len(states) {
			call = len(states)
		}
		return states[call-1]
	}
}

func completedTargets(r *rig) []model.RouteTarget {
	var targets []model.RouteTarget
	for _, block := range r.writer.blocks {
		for _, change := range block.changes {
			targets = append(targets, change.RouteTarget)
		}
	}
	return targets
}

// TestExecuteBuildsRouteTargetsFromTheRuleSet trägt die Auswertung im Run:
// jede Change des Runs trägt das Ziel der treffenden Regel mit der kleinsten
// `order` über den Quellwert der Zeile (Auffangregel `rest` mit der größten
// `order` zuletzt; ein abwesender Wert trifft nur die Auffangregel), unabhängig
// vom Block der Zeile; gleich dem Ziel, das `model.EvaluateRoute` für dieselbe
// Zeile und Regelliste bestimmt (die Funktion, die der WAL-Pfad ruft). Der
// Quellwert ist der Wert vor der Transformation: `rename_column` an der
// Bedingungsspalte ändert kein Ziel. Die Regeln einer anderen Tabelle wirken
// nicht; ohne Regel ist das Ziel leer. Rot färbende Mutationen (je eine, an der
// Eingabeseite): `nil` statt der Regelliste an `EvaluateRoute` im Run (alle
// Ziele leer); `nil` statt der Zeile (nur die Auffangregel trifft); der
// Regelstand aller Tabellen statt der Tabelle des Runs (der Fall „Regeln einer
// anderen Tabelle“ trägt Ziele).
func TestExecuteBuildsRouteTargetsFromTheRuleSet(t *testing.T) {
	rules := []model.RouteRule{
		routeRule(t, "rest", "rest", 99, nil),
		routeRule(t, "d", "ziel_d", 20, whenName("d")),
		routeRule(t, "a", "ziel_a", 10, whenName("a")),
	}
	columns := []string{"id", "name", "secret"}
	cases := []struct {
		name           string
		state          map[string][]model.RouteRule
		transformation map[string][]model.Transformation
		want           []model.RouteTarget
	}{
		{"Regelliste", routesOf(testQualified, rules...), nil,
			[]model.RouteTarget{"ziel_a", "rest", "rest", "ziel_d", "rest"}},
		{"Regelliste mit rename_column an der Bedingungsspalte", routesOf(testQualified, rules...),
			rulesOf(testQualified, rename(t, "r-name", "name", "label")),
			[]model.RouteTarget{"ziel_a", "rest", "rest", "ziel_d", "rest"}},
		{"Regeln einer anderen Tabelle", routesOf("public.other", rules...), nil,
			[]model.RouteTarget{"", "", "", "", ""}},
		{"keine Regel", nil, nil,
			[]model.RouteTarget{"", "", "", "", ""}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := newRig()
			r.routes.stateFn = routeStates(tc.state)
			r.rules.stateFn = ruleStates(tc.transformation)
			run := mustExecute(t, r)
			if run.Status != model.BackfillRunCompleted {
				t.Fatalf("Run = %+v, will completed", run)
			}
			got := completedTargets(r)
			if fmt.Sprint(got) != fmt.Sprint(tc.want) {
				t.Fatalf("Ziele = %q, will %q", got, tc.want)
			}
			// Vertrag mit der gemeinsamen Auswertung: dasselbe Ziel, das die
			// Domänen-Funktion für die Quellwerte der Zeile bestimmt.
			index := 0
			for _, block := range r.snapshot.blocks {
				for _, row := range block {
					if want := model.EvaluateRoute(tc.state[testQualified], columns, row); got[index] != want {
						t.Fatalf("Zeile %d: Ziel %q, EvaluateRoute %q", index+1, got[index], want)
					}
					index++
				}
			}
			if r.routes.gotSource != testSource || r.routes.calls != 5 {
				t.Fatalf("Routing-Regelstand der Quelle %q in %d Lesungen, will %q in 5 (Beginn, drei Blöcke, vor dem Commit)", r.routes.gotSource, r.routes.calls, testSource)
			}
		})
	}
}

// TestExecuteRoutingStateChangeEndsRunAsConfiguration trägt die
// Fail-closed-Prüfung des Routing-Regelstands (`ADR-0139`): der Stand der
// Lesung zu Beginn ist der Stand des Runs; jede Abweichung in einer späteren
// Lesung — ein `set_route` oder `remove_route` zwischen zwei Blöcken, ein
// Wechsel unmittelbar vor dem Commit, eine ersetzte Regel gleichen Namens —
// rollt zurück und endet den Run `failed` mit der Klasse `configuration`; ohne
// Wechsel läuft der Run durch, ein reiner Ordnungsunterschied und eine
// Doppelung sind keine Abweichung. Nur der Routing-Stand wechselt, der
// Transformations- und der Ausschlussstand bleiben. Rot färbende Mutationen
// (je eine): der Vergleich je Block entfällt (die Fälle zwischen Blöcken enden
// `completed`); der Vergleich vor dem Commit entfällt (die Fälle „vor dem
// Commit“); der Vergleich als Längenvergleich (der Fall „Regel ersetzt“).
func TestExecuteRoutingStateChangeEndsRunAsConfiguration(t *testing.T) {
	rule := func(target string) model.RouteRule { return routeRule(t, "regel", target, 10, whenName("a")) }
	other := routeRule(t, "andere", "ziel_b", 20, whenName("c"))
	none := routesOf(testQualified)
	one := routesOf(testQualified, rule("ziel_a"))
	cases := []struct {
		name       string
		states     func(int) map[string][]model.RouteRule
		wantFailed bool
		wantBlocks int
	}{
		{"ohne Wechsel", routeStates(one), false, 3},
		{"Ordnung und Doppelung sind keine Abweichung", routeStates(
			routesOf(testQualified, rule("ziel_a"), other), routesOf(testQualified, other, rule("ziel_a")),
			routesOf(testQualified, other, rule("ziel_a"), other), routesOf(testQualified, rule("ziel_a"), other)), false, 3},
		{"set_route zwischen Block 1 und 2", routeStates(none, none, one), true, 1},
		{"remove_route zwischen Block 1 und 2", routeStates(one, one, none), true, 1},
		{"set_route zwischen Block 2 und 3", routeStates(none, none, none, one), true, 2},
		{"set_route vor dem Commit", routeStates(none, none, none, none, one), true, 3},
		{"remove_route vor dem Commit", routeStates(one, one, one, one, none), true, 3},
		{"Regel ersetzt (anderes Ziel) zwischen Block 1 und 2", routeStates(one, one, routesOf(testQualified, rule("ziel_neu"))), true, 1},
		{"zweite Regel gesetzt zwischen Block 1 und 2", routeStates(one, one, routesOf(testQualified, rule("ziel_a"), other)), true, 1},
		{"Zwischenblock weicht ab, am Ende wieder gleich", routeStates(one, one, none, one), true, 1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := newRig()
			r.routes.stateFn = tc.states
			run := mustExecute(t, r)
			if !tc.wantFailed {
				if run.Status != model.BackfillRunCompleted || len(r.writer.committed) != 1 {
					t.Fatalf("Run = %+v, will completed", run)
				}
				return
			}
			if run.Status != model.BackfillRunFailed || !strings.HasPrefix(run.ErrorMessage, "configuration: ") ||
				!strings.Contains(run.ErrorMessage, domainerrors.ErrRoutingStateChanged.Error()) {
				t.Fatalf("Run = %+v, will failed mit Klasse configuration und dem Wechsel des Routing-Regelstands", run)
			}
			if len(r.writer.committed) != 0 || r.trace.count("Commit") != 0 || r.trace.count("Notify") != 0 {
				t.Fatalf("Commit oder Wecksignal trotz Abweichung: %v", r.trace.events)
			}
			if r.writer.rollbacks != 1 || r.snapshot.closed != 1 || len(r.writer.blocks) != tc.wantBlocks {
				t.Fatalf("Rollbacks=%d Snapshot geschlossen=%d Blöcke=%d, will 1/1/%d", r.writer.rollbacks, r.snapshot.closed, len(r.writer.blocks), tc.wantBlocks)
			}
		})
	}
}

// TestExecuteOnlyRoutingStateChangesEndsRunAsConfiguration trägt die
// Unabhängigkeit der beiden Regelstände: bleibt der Transformationsstand über
// alle Lesungen gleich, genügt der Wechsel des Routing-Stands allein. Rot
// färbende Mutation: der Vergleich des Routing-Standes entfällt in jedem Block
// und vor dem Commit, während der Transformationsvergleich bleibt (der Run
// endet `completed`).
func TestExecuteOnlyRoutingStateChangesEndsRunAsConfiguration(t *testing.T) {
	r := newRig()
	r.rules.stateFn = ruleStates(rulesOf(testQualified, rename(t, "r-name", "name", "label")))
	r.routes.stateFn = routeStates(routesOf(testQualified), routesOf(testQualified, routeRule(t, "regel", "ziel", 10, nil)))
	run := mustExecute(t, r)
	if run.Status != model.BackfillRunFailed || !strings.HasPrefix(run.ErrorMessage, "configuration: ") {
		t.Fatalf("Run = %+v, will failed mit Klasse configuration", run)
	}
	if r.rules.calls < 2 {
		t.Fatalf("Lesungen des Transformationsstands = %d, will mindestens 2", r.rules.calls)
	}
}

// TestExecuteRoutingReadFailureEndsRun trägt den Lesefehler-Zweig: ein
// Routing-Regelstand, der nicht gelesen werden kann, ist kein bestätigter
// Stand. Der n-te Lesefehler (Lesung 1 zu Beginn, 2 bis 4 je Block, 5 vor dem
// Commit) endet den Run `failed` mit der Klasse `configuration` und der
// Ursache im Fehlertext, ohne Commit und ohne Wecksignal. Die Eingabe ist der
// Lesefehler an der jeweiligen Lesung; weil nur der n-te Aufruf scheitert,
// färbt sich der Test rot, sobald ein Lesefehler verworfen wird (der Run läuft
// bis zum Commit durch) oder wieder mit der Klasse der Ursache (`storage`)
// endet.
func TestExecuteRoutingReadFailureEndsRun(t *testing.T) {
	for call := 1; call <= 5; call++ {
		t.Run(fmt.Sprintf("Lesung %d", call), func(t *testing.T) {
			r := newRig()
			r.routes.err = outbound.ErrStorage
			r.routes.errCall = call
			run := mustExecute(t, r)
			if run.Status != model.BackfillRunFailed || !strings.HasPrefix(run.ErrorMessage, "configuration: ") ||
				!strings.Contains(run.ErrorMessage, outbound.ErrStorage.Error()) {
				t.Fatalf("Run = %+v, will failed mit Klasse configuration und der Ursache im Text", run)
			}
			if r.routes.calls != call {
				t.Fatalf("Lesungen = %d, will %d (der Lauf bricht an der Lesung ab)", r.routes.calls, call)
			}
			if len(r.writer.committed) != 0 || r.trace.count("Commit") != 0 || r.trace.count("Notify") != 0 {
				t.Fatalf("Commit oder Wecksignal trotz Lesefehler: %v", r.trace.events)
			}
			if r.snapshot.closed != 1 {
				t.Fatalf("Snapshot %d-mal geschlossen, will 1", r.snapshot.closed)
			}
		})
	}
}

// TestExecuteInapplicableRoutingRuleEndsRunAsSchema trägt die Negative der
// Routing-Regel im Run (`ADR-0138`): eine Regel, deren `when.column` in den
// Spalten des Snapshots fehlt, endet den Run `failed` mit der Klasse `schema`
// und Regelname und Spalte im Fehlertext — nachdem der Snapshot seine Spalten
// lieferte und bevor eine Zeile gelesen und die Schreibtransaktion geöffnet
// wird, auch bei einer leeren Tabelle; es entsteht keine Change, der Snapshot
// ist geschlossen, der Regelstand wird einmal gelesen. Eine Regel ohne
// Bedingung, eine Regel auf eine vorhandene Spalte und die Regeln einer
// anderen Tabelle lassen den Run durchlaufen. Rot färbende Mutationen (je
// eine, an der Eingabeseite): die Prüfung `checkRoutesApplicable` entfällt
// (Run `completed`); die Prüfung gegen `nil` statt der Snapshot-Spalten (die
// anwendbare Regel scheitert); die Prüfung über den Regelstand aller Tabellen
// (die fremde Regel scheitert).
func TestExecuteInapplicableRoutingRuleEndsRunAsSchema(t *testing.T) {
	missing := func(t *testing.T) model.RouteRule {
		return routeRule(t, "r-fehlt", "ziel", 10, &model.RouteCondition{Column: "gibt_es_nicht", Equals: "x"})
	}
	cases := []struct {
		name   string
		state  func(t *testing.T) map[string][]model.RouteRule
		empty  bool
		failed bool
		wantIn []string
	}{
		{"Spalte fehlt", func(t *testing.T) map[string][]model.RouteRule { return routesOf(testQualified, missing(t)) },
			false, true, []string{`"r-fehlt"`, `"gibt_es_nicht"`, testQualified}},
		{"Spalte fehlt bei leerer Tabelle", func(t *testing.T) map[string][]model.RouteRule { return routesOf(testQualified, missing(t)) },
			true, true, []string{`"r-fehlt"`}},
		{"zweite Regel nicht anwendbar", func(t *testing.T) map[string][]model.RouteRule {
			return routesOf(testQualified, routeRule(t, "r-ok", "ziel", 5, whenName("a")), missing(t))
		}, false, true, []string{`"r-fehlt"`}},
		{"Regel einer anderen Tabelle wird nicht geprüft", func(t *testing.T) map[string][]model.RouteRule { return routesOf("public.other", missing(t)) },
			false, false, nil},
		{"Regel ohne Bedingung ist anwendbar", func(t *testing.T) map[string][]model.RouteRule {
			return routesOf(testQualified, routeRule(t, "r-alle", "ziel", 5, nil))
		}, false, false, nil},
		{"Regel auf eine vorhandene Spalte", func(t *testing.T) map[string][]model.RouteRule {
			return routesOf(testQualified, routeRule(t, "r-ok", "ziel", 5, whenName("a")))
		}, false, false, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := newRig()
			if tc.empty {
				r.snapshot.blocks = nil
			}
			r.routes.stateFn = routeStates(tc.state(t))
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
			if !strings.Contains(run.ErrorMessage, domainerrors.ErrRoutingColumnMissing.Error()) {
				t.Fatalf("Fehlertext = %q, will die Ursache %q", run.ErrorMessage, domainerrors.ErrRoutingColumnMissing)
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
			if r.routes.calls != 1 {
				t.Fatalf("Lesungen des Routing-Regelstands = %d, will 1 (Prüfung einmal je Run, vor dem ersten Block)", r.routes.calls)
			}
			if len(r.runs.finished) != 1 || r.runs.finished[0] != run {
				t.Fatalf("Finish = %+v", r.runs.finished)
			}
		})
	}
}

// TestExecuteInapplicabilityPrecedesStateChange trägt die Reihenfolge: eine
// nicht anwendbare Regel zu Beginn endet den Run `schema`, bevor eine spätere
// Lesung einen Wechsel sähe — der Run liest den Routing-Stand nur einmal und
// kein Block läuft. Eine nicht anwendbare Transformationsregel wird vor der
// nicht anwendbaren Routing-Regel gemeldet. Rot färbende Mutation: die
// Anwendbarkeitsprüfung des Routings entfällt und der Wechsel in Lesung 2
// endet den Run `configuration` statt `schema`.
func TestExecuteInapplicabilityPrecedesStateChange(t *testing.T) {
	missing := routeRule(t, "r-fehlt", "ziel", 10, &model.RouteCondition{Column: "gibt_es_nicht", Equals: "x"})
	t.Run("nicht anwendbar vor Wechsel", func(t *testing.T) {
		r := newRig()
		r.routes.stateFn = routeStates(routesOf(testQualified, missing), routesOf(testQualified))
		run := mustExecute(t, r)
		if run.Status != model.BackfillRunFailed || !strings.HasPrefix(run.ErrorMessage, "schema: ") {
			t.Fatalf("Run = %+v, will failed mit Klasse schema", run)
		}
		if r.routes.calls != 1 || r.trace.count("NextBlock") != 0 {
			t.Fatalf("Lesungen = %d, NextBlock = %d, will 1 und 0", r.routes.calls, r.trace.count("NextBlock"))
		}
	})
	t.Run("Transformationsregel zuerst", func(t *testing.T) {
		r := newRig()
		r.rules.stateFn = ruleStates(rulesOf(testQualified, rename(t, "t-fehlt", "gibt_es_nicht", "x")))
		r.routes.stateFn = routeStates(routesOf(testQualified, missing))
		run := mustExecute(t, r)
		if run.Status != model.BackfillRunFailed || !strings.HasPrefix(run.ErrorMessage, "schema: ") || !strings.Contains(run.ErrorMessage, `"t-fehlt"`) {
			t.Fatalf("Run = %+v, will schema mit der Transformationsregel t-fehlt", run)
		}
	})
}
