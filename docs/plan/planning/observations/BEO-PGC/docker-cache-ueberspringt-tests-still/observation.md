# BEO-PGC/docker-cache-ueberspringt-tests-still

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die Bau-Ziele, deren
Test-Stufe Teil eines `docker build` ist: `make sdk-pack-csharp`, `make sdk-pack-python`,
`make sdk-pack-kotlin`, `make examples-csharp`, `make examples-kotlin`).

Die Beobachtung: Die Tests dieser Ziele laufen als Stufe eines `docker build`. Ändert sich
der Eingabe-Kontext der Stufe nicht, liefert Docker die Schicht aus dem Cache: das Ziel
endet mit Exit 0, **ohne dass ein Test ausgeführt wird**, und es druckt keine Testzeile.
Ein grüner Lauf dieser Ziele ist deshalb nur dann ein Beleg „die Tests sind gelaufen“,
wenn der Lauf eine Testzeile druckt oder wenn der Cache umgangen wird (`--no-cache` an der
Test-Stufe) oder eine Mutation die Stufe rot färbt. Hinzu kommt: der grüne Gradle-Lauf
druckt keine Testzahl, der rote nennt sie (`N tests completed, M failed`) — die Kotlin-Zahl
ist ohne roten Lauf nicht lesbar.

**Abgrenzung.** `BEO-PGC/test-runner-stiller-ausschluss` zählt den Test, den das eigene
Runner-Skript über ein `-run`-Muster auslässt; `BEO-PGC/test-methode-lauft-still-nicht`
den Test, dessen Signatur-Form die Plattform nicht entdeckt; `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`
den genannten Befehl, der die behauptete Aussage nicht stützt. Hier ist der Test
vorhanden, entdeckbar und im Runner erfasst — der Docker-Schicht-Cache führt ihn nicht aus,
und der Befehl (`make sdk-pack-*`) meldet trotzdem Erfolg.

**Warum das zählt:** Ein Bericht „sdk-pack grün“ trägt die Aussage „die Tests sind grün“
nur, wenn die Stufe lief. Kein Gate liest die Differenz, und der Fall tritt bei jedem Lauf
ohne Änderung der Eingabe auf, also gerade bei Reviewer und Verifier, die denselben Stand
noch einmal bauen.
