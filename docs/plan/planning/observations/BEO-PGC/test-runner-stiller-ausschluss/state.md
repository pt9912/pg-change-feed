Zustand: **verkörpert** — Ausgang: **verkörpert** → `test/integration/runner_vollstaendigkeit_test.go`, `TestRunnerFuehrtJedeE2EFunktionAus` (liegt in dieser Datei; Feld und Zielort auf einer Zeile) · seit slice-harness-integration-runner-vollstaendigkeit.
Zähler (abgeleitet): **3×** (evidence/slice-033.md, evidence/slice-074.md,
evidence/slice-routing-nats-subjekt.md) — Schwelle erreicht.

Gegenstand: ob jede Testfunktion in `test/integration` von mindestens einem `-run`-Muster in
`tools/harness/run-integration-tests.sh` erfasst wird. Der Wächter liest die `-run`-Argumentwerte
des Skripts (Kommentare zählen nicht) und färbt rot, wenn eine `func TestE2E*` in keinem steht;
sein Leser ist an Eingaben je Zweig gebunden (`TestRunnerLeserDreiZustaende`). Die
Deklarations-Hälfte (neue Runner-Phase ohne Deklaration) trägt `TestAbdeckungstabelleZeilen`
(`slice-074`); für `tools/harness/run-notify-tests.sh` ist die Lücke geschlossen (ganzes Paket,
Exit 1 bei `--- SKIP`, `slice-routing-nats-subjekt`).

Grenze: der Wächter prüft die Anwesenheit im Muster, nicht dass die Funktion ohne `--- SKIP`
läuft und nicht die Phasen ohne Go-Testfunktion; die Shell-Grammatik liest er nicht
(Godoc der Funktion). Wirksamkeit **erprobt**: die Mutation „Name aus dem Sammelmuster entfernt“ färbt
rot (Verifikation `verifikation-slice-harness-integration-runner-vollstaendigkeit` §2 c, vom
Verifier gemessen).
