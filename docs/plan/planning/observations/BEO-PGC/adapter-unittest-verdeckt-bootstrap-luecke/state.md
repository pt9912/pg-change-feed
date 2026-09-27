Zustand: **verkörpert** — Ausgang: **verkörpert** → Phase
„Prozessstart-Vorlauf-Frist“ in `tools/harness/run-integration-tests.sh` ·
seit slice-start-vorlauf-grenze. Der reale Rundlauf (Sperre, `pending`-Antrag,
Neustart, Healthcheck-Poll, `cdc.changes`-Poll) liefert genau den Beleg, den
die Whitebox-Tests von `internal/bootstrap` mit Fakes allein nicht zeigen
konnten — der Prozess bleibt über den Healthcheck gesund, während der
Replikationsstrom bereits läuft und der Vorlauf noch wartet; diese
Eigenschaft erscheint erst am komponierten, realen Prozess (dieselbe Klasse
wie beim vierten Beleg unten). Zähler (abgeleitet):
**5×** (evidence/slice-061.md, evidence/slice-072.md,
evidence/slice-transformationen-start-reihenfolge.md,
evidence/architect-verdict-welle-transformationen-offene-fragen.md,
evidence/architect-verdict-wal-fehlerschwelle-ausgangsklasse.md).

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

Der fünfte Beleg trägt den Träger bestätigend: die Runner-Phase „Fehlerschwelle beendet den
Container“ fand die Abweichung der Ausgangs-Klasse, die die Unit-Tests der Kette mit Fakes nicht
zeigten, und ein Slice trägt die Korrektur (`slice-wal-fehlerschwelle-ausgangsklasse`). Der
Trigger der Neubewertung ist damit nicht eingetreten. Die Phase ist am realen Prozess
falsifiziert: die Mutation „Aufruf von `stopStream` in der Schwellen-Prüfung entfernt“ färbt
sie rot (der Container lief 90 s nach einer Last über der Fehlerschwelle weiter; Kopie des Repos,
Image aus der Kopie, voller `make test-integration`, vom Verifier gefahren,
`verifikation-slice-capture-leerlauf-quellbelege` §4 P2). Kein weiterer Beleg der Klasse aus
diesem Vorgang: Zähler unverändert 5×.
