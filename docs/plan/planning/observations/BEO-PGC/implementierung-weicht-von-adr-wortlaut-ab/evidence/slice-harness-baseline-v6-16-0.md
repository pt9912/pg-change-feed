**Vorgang:** slice-harness-baseline-v6-16-0 (Review F-1, MEDIUM; Re-Review zur Fixrunde R-1, LOW, gleich Verifikation V-1)

**Fund:** Folgepflicht 6 von `ADR-0161` verlangte, in `AGENTS.md` §3.5 den Verweis „für einen
MR-Eintrag `ADR-0159` Entscheidung 2“ zu **ersetzen**; der Nachzug strich ihn nur, und
`harness/targets/zitat-vergleich.md` beschränkte die Beleg-Rolle auf die `Accepted` ADR — zwei
Träger sagten Verschiedenes über den Beleg des Form-Commits (Review F-1). Die Fixrunde nahm den
Fixrunden-Auftrag „mit dem `formnorm`-`cmp` als Beleg“ als **Zusatz** zur `zitat-vergleich`-Zeile
statt als Ersatz: §3.5 verlangte danach für die Form-Korrektur an einem MR-Eintrag beide
Messungen, während `ADR-0161` Entscheidung 4 nur `teilrange` und den `formnorm`-`cmp` nennt und
§Konsequenzen die Referent-Messung als Bedingung einer Zitat-Korrektur an MR-Einträgen streicht
(Re-Review R-1, Verifikation V-1). Bei der Closure an den Wortlaut der ADR angeglichen: §3.5 und
`zitat-vergleich.md` nennen `teilrange` und `formnorm`-`cmp` **an Stelle** der
`zitat-vergleich`-Zeile.

**Form (Ausprägung):** weder verdeckt noch benannt — eine **unbeabsichtigte** Abweichung: der
Träger überschritt den Wortlaut nicht als Entscheidung, sondern durch die Lesart eines
Auftrags, der den Wortlaut nur wiedergeben sollte (Ursache: die Vorgabe des Fixrunden-Auftrags
lag nicht als Artefakt im Repo, `BEO-PGC/rollen-uebergabe-ohne-committetes-artefakt`).
Träger-Typ `AGENTS.md` und Target-Vertrag. Vor dem Merge gefunden. Zwei Funde, ein Vorgang,
eine Datei.

Quelle: `docs/reviews/review-slice-harness-baseline-v6-16-0.md` (F-1), <!-- d-check:status-provenance -->
`docs/reviews/review-slice-harness-baseline-v6-16-0-fixrunde.md` (R-1), <!-- d-check:status-provenance -->
`docs/reviews/verify-slice-harness-baseline-v6-16-0.md` (V-1). <!-- d-check:status-provenance -->
