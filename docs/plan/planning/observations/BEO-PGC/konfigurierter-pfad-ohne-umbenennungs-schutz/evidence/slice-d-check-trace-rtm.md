# Beleg: slice-d-check-trace-rtm

Vorgang: `slice-d-check-trace-rtm` — `d-check --trace`/Requirements
Traceability Matrix verdrahtet (Welle `welle-d-check`).

Fund (§6 Risiko 1 des Slice-Plans): `.d-check.yml`s neuer
`trace.coverage`-Block referenziert `docs/user/e2e-abdeckung.md` als
hartcodierten Pfad — das Erzeugnis von `make test-integration`
(`tools/harness/run-integration-tests.sh` schreibt es, `harness/README.md`
§Werkzeuge). Kein Sensor hält diese Referenz gegen eine künftige Umbenennung
oder Verschiebung der Datei; würde sie umbenannt, bliebe `trace.coverage`
syntaktisch gültig, aber wirkungslos (keine Coverage-Spalte mehr, Waisenzahl
stiege unbemerkt von 7 auf 9 zurück — `--trace` ist advisory, kein Gate
würde das melden).

Aktuell kein Umbenennungs-Anlass erkennbar (Planner-Einschätzung bei
Closure, 2026-09-17) — Risiko bleibt „weiter offen", kein Carveout, kein
Folge-Slice nötig, da kein akuter Bedarf.

Quelle: `docs/plan/planning/done/welle-d-check/slice-d-check-trace-rtm.md` §6 Risiko 1 ·
der Verifikationsbericht zu `slice-d-check-trace-rtm` §C Risiko 1.
