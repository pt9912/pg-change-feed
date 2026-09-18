# Beleg: slice-099

Vorgang: `slice-099` — Kotlin-Sprachwurzel und HTTP-Client, zweiter realer Bau
der Beispiel-Client-Matrix (`ADR-0090`), parallel zu `slice-098`.

Fund: Der Slice erweitert den in `slice-098` neu angelegten Workflow
(`.github/workflows/examples.yml`) um einen **zweiten** Job
(`examples-kotlin`) — `AGENTS.md` §3.10 ist damit erneut dem Buchstaben nach
ausgelöst. Weder Implementer noch Reviewer noch Verifier konnten in dieser
Session einen realen Post-Push-Lauf auslösen oder prüfen (kein `git push`,
kein Runner-Zugriff) — dieselbe strukturelle Grenze wie bei den sechs
vorherigen Belegen. Alle drei Träger (Slice-Plan §5, Slice-Plan §6, der
Workflow-Datei-Kommentar selbst) führen den Punkt korrekt als offen, keiner
behauptet einen bereits erfolgten Lauf (Verifikationsbericht zu `slice-099`,
§5, dort ausdrücklich als „korrekt geführt“ bestätigt).

**Ausgang bei dieser Closure: weiter offen**, bis ein realer Post-Push-Lauf
sichtbar wird — kein neuer Träger, derselbe verkörperte Mechanismus
(`AGENTS.md` §3.10).

Quelle: `docs/plan/planning/in-progress/slice-099-kotlin-sprachwurzel-http-client.md`
§6 (dritte Risiko-Zeile) · Verifikationsbericht zu `slice-099`, §5.
