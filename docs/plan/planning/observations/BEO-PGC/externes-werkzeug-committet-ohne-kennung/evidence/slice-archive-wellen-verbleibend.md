# Beleg: slice-archive-wellen-verbleibend

Vorgang: `slice-archive-wellen-verbleibend` — 17 reale `archive-welle`-Läufe
(`welle-2`…`welle-11`, `welle-13`…`welle-17`, `welle-19`, `welle-20`).

Fund: Dieselbe zugrunde liegende Ursache wie beim ersten Beleg (fest
einprogrammierte, kennungslose Commit-Message des externen Werkzeugs) trat
hier in **17-facher** Wiederholung auf — 34 Commits statt vormals 4, alle
ohne `LH-*`/`ADR-*`-Kennung. Ein erster Versuch, dies als
Hintergrund-Zug mit Gate-Prüfung nach **jedem einzelnen** Lauf zu
behandeln, brach beim zweiten Lauf korrekt ab, statt die erwartungsgemäß
rote `commit-traceability` durch improvisierte Commits zu erzwingen — die
tatsächlich tragfähige Strategie erwies sich als umgekehrt: alle 17
mechanischen Läufe zuerst (mit `make docs-check` je Lauf statt des vollen,
teuren `make gates`), danach **eine** gebündelte Fenster-Räumung am Ende
über die ohnehin anfallenden, Kennung-tragenden Closure-Commits dieses
Slices. Die Räumungskosten sind konstant (die Fenstergröße des Standing-Gates,
5 Commits) und nicht proportional zur Anzahl der Läufe — dieselbe
Erkenntnis, die der erste Beleg bereits als offene Frage aufwarf, ist hiermit
real bestätigt.

Quelle: `docs/plan/planning/done/slice-archive-wellen-verbleibend.md` §2/§6/§7 ·
Commits `82647c4`…`96b313b` (34 Commits, siehe Slice-Plan §2 LP2 für die
vollständige Paar-Liste) · Closure-Commits dieses Slices.
