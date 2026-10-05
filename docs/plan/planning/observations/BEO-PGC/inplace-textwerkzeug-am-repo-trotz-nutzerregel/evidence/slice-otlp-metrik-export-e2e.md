**Vorgang:** slice-otlp-metrik-export-e2e (Review F-1, Planner-Nachtrag bei der Closure).

**Fund:** Die Umleitungs-Hälfte von `AGENTS.md` §3.1 in der Rolle Implementer: Text wurde per
`cat >>` an `tools/harness/otlpcheck/check_test.go` angehängt, statt per Edit/Write. Der
Implementer meldete es selbst; der Reviewer (Review-Report
zu `slice-otlp-metrik-export-e2e`, F-1, HIGH) hält fest, dass der Guard
Umleitungen nicht liest und der Diff die Schreibweise nicht zeigt. Der Inhalt ist geprüft und
unauffällig (Test grün, Mutation der Soll-Einheit färbt `TestCheckRejectsEachDeviation` rot).
Verifizierbar am Repo: nein; Ursprung der Angabe **übernommen**, nicht am Guard gemessen.

**Form (Ausprägung):** Umleitungs-Anhang trotz geltender Regel; kein neuer Mechanismus, die
Regel in `AGENTS.md` §3.1 galt und der Guard-Vertrag (`MR-003`) nennt die Grenze (Umleitungen werden
nicht gelesen). Elfter Beleg dieses Eintrags.
