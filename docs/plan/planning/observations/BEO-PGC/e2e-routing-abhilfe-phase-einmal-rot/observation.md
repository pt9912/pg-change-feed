# BEO-PGC/e2e-routing-abhilfe-phase-einmal-rot

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft den E2E-Runner
`tools/harness/run-integration-tests.sh`, Phase „Routing-Nichtanwendbarkeit und Abhilfe“,
und den GitHub-Workflow `e2e`).

Die Beobachtung: Der Job `image + test-integration (PostgreSQL 18)` des Workflows `e2e`
endete im Lauf 36983752407 (Push-Commit 4a43f6ac, 2026-10-02) mit `failure` in der Phase
„Routing-Nichtanwendbarkeit und Abhilfe (b) — ein zweiter schema-Fehler nach der Abhilfe
(error_class=schema)“; `gh run rerun --failed` desselben Laufs endete beide Legs mit
`success`. Ursprung: gemessen vom Hauptlauf mit `gh run view`/`gh run list`.

**Stelle der Prüfung (Anker, gemessen am Stand 4a43f6ac durch Lesen):**
`tools/harness/run-integration-tests.sh`, Funktion `rn_abhilfe` (Zeilen 5331–5343), Aufruf für
Fall (b) in Zeile 5358; die Meldung stammt aus Zeile 5340. Die dort und in den aufgerufenen
Hilfen verwendeten Fristen (gelesen, nicht gemessen):

- `bf_await_healthy` (Zeilen 3497–3507): bis zu 60 Versuche im Abstand 1 s;
- `bf_await_applied` (Zeilen 3444–3458): bis zu 60 Versuche im Abstand 0,5 s;
- `bf_await_sql` für die zuvor nicht bestätigte Transaktion (id=1): 30 s, Schritt 0,25 s;
- `rn_await_end` (Zeilen 5305–5314) vor `rn_abhilfe`: bis zu 90 Versuche im Abstand 1 s.

**Was nicht untersucht ist.** Die Ursache ist weder gemessen noch benannt. Gemessen ist nur:
zwischen dem letzten grünen `e2e` (Commit 4e664547) und dem roten Lauf änderte sich kein
Server-Code, kein `tools/harness/run-integration-tests.sh` und nichts unter `test/`
(`git diff --stat 4e664547 HEAD -- internal cmd tools/harness/run-integration-tests.sh test`
leer). Ein Zeitfenster-Flake ist *hergeleitet*, nicht gemessen; ebenso offen ist, ob die Prüfung in
Zeile 5339 (Lesen von `error_class` aus `cdc.heartbeat` nach `applied`) einen Zustand aus dem
Lauf des Vorgängerprozesses oder einen neuen Fehler des neu gestarteten Prozesses liest. Das Log
des roten Jobs ist in dieser Eintragung nicht gesichert.

**Offene Frage für das nächste Auftreten:** das Log des roten Jobs vor einem Re-Run sichern
(`gh run view <Lauf> --log-failed`), die letzten Zeilen des Feed-Container-Logs der Phase lesen
und die Fristen der Hilfen gegen die gemessenen Zeitpunkte halten. Eine Entscheidung über Fristen
oder Prüfform ist damit nicht getroffen.
