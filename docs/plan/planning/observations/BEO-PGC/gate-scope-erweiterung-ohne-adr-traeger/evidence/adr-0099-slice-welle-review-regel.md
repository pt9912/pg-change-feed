**Vorgang:** wellenloser Architect-Zug (Korrektur), `ADR-0099`
**Fund:** Commit `01b7b09` (Anlass: ein zu weiter `**`-Glob der
`slice`/`welle`-Matrixklassen in `.d-check.yml`) hat im selben Commit
zusätzlich zwei neue Matrix-Regeln ergänzt — `{from: slice, to: review,
allow: false}` und `{from: welle, to: review, allow: false}` — ohne ADR
und im direkten Widerspruch zu zwei bereits `Accepted` ADRs
(`ADR-0094`, `ADR-0097`), die genau diese Regel bewusst ausgelassen
hatten. Die Konfigurationsdatei selbst trug die Erweiterung, wie im
Erstauftreten (`evidence/slice-071.md`, dort `.a-check.yml`) — hier
`.d-check.yml`, dieselbe Fehlerklasse: Gate-Scope wächst durch die
Config-Datei, ohne den `AGENTS.md` §3.6-Träger. Aufgelöst über `ADR-0099`
— Regel zurückgenommen, `ADR-0094`s ursprüngliche Selbst-Zitat-Begründung
mit 11 real geprüften Fällen bestätigt (alle 11 Selbst-Zitate, keine
Fremd-Zitation). Zweites Auftreten.
