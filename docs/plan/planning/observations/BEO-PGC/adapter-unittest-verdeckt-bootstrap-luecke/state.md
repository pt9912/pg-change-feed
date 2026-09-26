Zustand: **geplant** — Ausgang: **geplant** → `slice-start-vorlauf-grenze` (Phase in `make
test-integration`, die den Prozessstart mit einem wartenden Antrag fährt; Architect-Verdikt
`architect-verdict-welle-transformationen-offene-fragen` §8, `ADR-0128`). Zähler (abgeleitet):
**4×** (evidence/slice-061.md, evidence/slice-072.md,
evidence/slice-transformationen-start-reihenfolge.md,
evidence/architect-verdict-welle-transformationen-offene-fragen.md).

Der vierte Beleg ist die Messung des Architect-Zugs: die Start-Sequenz aus `Run` trägt Unit-Tests
mit Fakes und zwei Quelltext-Tests; die Eigenschaft, die sie nicht zeigen (der Replikationsstrom
läuft schon, während der Vorlauf wartet, und der Server beendet ihn), erschien erst am
komponierten Prozess. Der bisher vorgeschlagene Träger dieses Eintrags, ein Bootstrap-Smoke-Test
(`ConfigFromEnv` + `Run()`-Teilaufruf, kein `nil`-Use-Case), hätte diesen Fall nicht gefunden: er
prüft die Verdrahtung, nicht den Verlauf einer Verbindung, und braucht die Datenbank wie `Run`
selbst (hergeleitet, nicht erprobt). Der Träger der Klasse ist die Phase in `make test-integration`
in dem Slice, der die Eigenschaft einführt (Testpyramide, `ADR-0030`); kein eigener Smoke-Test.
Trigger der Neubewertung: ein Auftreten, das eine Phase von `make test-integration` gefunden hätte
und in keinem Slice steht.
