# Beleg: slice-038

Vorgang: slice-038 (CLI-Diagnose-Befehl, `in-progress/` zum Zeitpunkt
dieses Belegs — der Beleg wird vor dem `git mv` nach `done/` geschrieben,
siehe Baseline-Regelwerk `modul-06-roadmap.md` §Das Beobachtungs-Register).

Fund: Der Plan-Nachzug (§3, `docs/plan/planning/in-progress/
slice-038-cli-diagnose.md`) bestätigt real dieselbe strukturelle
Eigenschaft erneut, die diese Beobachtung seit `slice-007` führt:
`reportFault` (`internal/bootstrap/wiring.go`) schreibt `error_class` nur
unmittelbar vor `os.Exit` — jeder von `Run()` klassifizierte Fehler
beendet den Prozess, `restart: "no"` hält den Container danach beendet
stehen. Damit existiert aktuell kein am **laufenden** Prozess
beobachtbarer, nicht-terminaler Fehlerzustand — für `Diagnose`s
Fehlerzustands-Zweig (`LH-FA-ADM-003`, §6 Risiko 1 des Slice-Plans) heißt
das: Weder der Unit-Test (`TestDiagnoseReportsErrorState`,
`internal/bootstrap/diagnose_test.go`) noch der E2E-Beleg
(`tools/harness/run-integration-tests.sh`) lösen den Fehlerzustand über
den realen Erfassungspfad aus. Beide schreiben stattdessen denselben
Spaltenwert (`cdc.process_heartbeat.error_class`), den `reportFault` im
echten Fehlerfall schriebe, direkt per SQL — derselbe Lesepfad (View →
CLI-Ausgabe), weil der reale Auslösepfad strukturell keinen laufenden
Prozess zum Beobachten übrig lässt. Der Ersatzbeleg ist damit nicht
Testschwäche dieses Slice, sondern eine erneute reale Bestätigung der seit
`slice-007` offenen Beobachtung.

Quelle: `docs/reviews/review-slice-038.md` (F-3), Plan-Nachzug §3
(`docs/plan/planning/in-progress/slice-038-cli-diagnose.md`).
