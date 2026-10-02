**Vorgang:** slice-sdk-sse-filter-phase-verbindung-haertung (Review F-1, MEDIUM, außerhalb des Diffs)

**Fund:** Im ersten C#-Tier-Lauf des Reviewers endete die NATS-Phase mit „der Ablehnungs-Beleg blieb aus
(REJECTED token-rejected fehlt)“, obwohl die Zeile im ausgegebenen Container-Log stand; der zweite Lauf war grün
(gemessen vom Reviewer, Fehlertext des Laufs). Die Ursache ist hergeleitet: `docker logs … | grep -qF` in der
Ablehnungs-Schleife (`tools/harness/run-sdk-csharp-integration-tests.sh` Zeile 246, ebenso Kotlin Zeile 244 und Python
Zeile 280) verfehlt unter `pipefail` eine vorhandene Zeile (SIGPIPE). Gemessen vom Reviewer an einem Nachbau (ein
beendeter Container, eine Treffer-Zeile, 20 Folgezeilen): 300 Aufrufe, 3 Fehlschläge.

Weitere Vorkommen: der Implementer desselben Vorgangs meldete dasselbe Symptom (gehört zum selben Vorgang, keine
eigene Datei, im Bericht des Implementers **übernommen**, ohne auflösbaren Anker); in den fünf C#-Gesamtläufen des
Verifiers trat der Ausfall nicht auf (gemessen, kein Gegenbeleg gegen die Rate von 1 %). Nicht gezählt: die
Implementer-Meldung eines einmaligen Ausfalls der NATS-Phase im Vorgänger-Slice (Review F-4 dort); ihre Ursache ist dort
als Zeitfenster 20 s gegen 15 s hergeleitet, abweichend von der hier genannten, und ein Lauf-Beleg fehlt.

Quelle: `docs/reviews/review-slice-sdk-sse-filter-phase-verbindung-haertung.md` (F-1) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-sdk-sse-filter-phase-verbindung-haertung.md` (§5, §6 Zeile NATS-Phasen-Flake). <!-- d-check:status-provenance -->
