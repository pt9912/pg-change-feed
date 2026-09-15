# BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft Kommentare, die
ein Verhalten zusichern, das der Code nicht trägt, keine eigene Sub-Area im Sinn
der Modus-Deklaration).

Die Beobachtung: Ein Kommentar sagt einen **Fehlerpfad** zu — „der Startfehler
erreicht das Ergebnis von `Run`" —, während der Code den Fehler in eigener
Goroutine nur **loggt** und der Prozess weiterläuft. Die Aussage ist nicht bloß
ungenau, sie kehrt das Verhalten um, und sie steht dort, wo ein Leser die
Semantik sucht: am Deklarationsort der Variablen. Kein Gate prüft
Kommentar-Wahrheit — der Wächter ist das Review. Belegt **zweimal in derselben
Datei**: `review-slice-070` F-2 berichtigte einen Block (HIGH),
`review-slice-077` F-1 den nächsten (HIGH); beide standen gleichzeitig im Baum,
und der erste Durchgang fand nur einen.
