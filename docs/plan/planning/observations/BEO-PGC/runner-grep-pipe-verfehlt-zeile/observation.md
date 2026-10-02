# BEO-PGC/runner-grep-pipe-verfehlt-zeile

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die Hilfsskripte der
Tier-Runner unter `tools/harness/`, keine eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Die Warte- und Prüfschleifen der Runner lesen den Ausgang eines Containers mit
`docker logs … | grep -qF "<Marke>"` unter `set -euo pipefail`. `grep -q` beendet sich beim
ersten Treffer; schreibt `docker logs` danach weitere Zeilen in die geschlossene Pipe, erhält es
SIGPIPE, und unter `pipefail` ist die Pipeline falsch, obwohl die Zeile vorhanden ist. Der
Runner liest dann „kein Treffer“ und meldet eine Phase als fehlgeschlagen, deren Test die Zeile
gedruckt hat.

Erstes Auftreten: im ersten C#-Tier-Lauf des Reviewers von
`slice-sdk-sse-filter-phase-verbindung-haertung` endete die NATS-Phase mit „der Ablehnungs-Beleg
blieb aus (REJECTED token-rejected fehlt)“, obwohl die Zeile im ausgegebenen Container-Log stand;
der zweite Lauf war grün (gemessen vom Reviewer, Review F-1). Der Mechanismus ist **hergeleitet**,
teilweise gemessen: ein beendeter Container mit einer Treffer-Zeile und 20 Folgezeilen,
`docker logs | grep -qF` unter `pipefail`, 300 Aufrufe, 3 Fehlschläge (vom Reviewer gemessen).

Stelle der Prüfung (Anker, am Stand `a420e223` gemessen mit `git grep`): 35 Zeilen der Form
`docker logs … | grep -q…` in sechs Dateien unter `tools/harness/`
(`run-integration-tests.sh` 21, die drei `run-sdk-*-integration-tests.sh` je 3,
`lib-sdk-filter-fixture.sh` 3, `lib-sdk-route-fixture.sh` 2).

Nachbar, nicht dasselbe: `BEO-PGC/pipe-maskiert-make-exit-code` (dort maskiert eine Pipe einen
roten Exit-Code; hier färbt sie einen vorhandenen Treffer falsch). Ob die Form den einmaligen
roten Lauf von `BEO-PGC/e2e-routing-abhilfe-phase-einmal-rot` erklärt, ist **nicht untersucht**;
an der betroffenen Phase trägt `run-integration-tests.sh` die Form nach Lesen nicht.
