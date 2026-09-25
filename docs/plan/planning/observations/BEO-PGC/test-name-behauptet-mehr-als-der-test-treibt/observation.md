# Der Name eines Tests und seine Meldungen behaupten mehr, als der Test treibt

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die Aussagekraft von
Testnamen und Fehlermeldungen, keine eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Ein Test wird beim Umbau auf einen Ersatz gezogen, der einen schwächeren
Aufbau fährt (hier: kein Stream, kein `Run`), und behält Namen und Meldungen des stärkeren
(„…EndToEnd“, „der Stream-Lauf wurde … nicht beendet“). Wer den Namen liest, hält eine Kette
für gebunden, die nur die Herleitung trägt. `make gates` fängt es nicht: der Test ist grün,
die Behauptung steht im Namen. Gefunden hat es der Reviewer durch den Vergleich von
Vorgänger und Ersatz.

**Abgrenzung.** `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` beschreibt eine Zusage,
deren Test die Eingabeseite nicht bindet; hier bindet der Test seine Eingabe, und der Name
benennt einen anderen Gegenstand als den getriebenen.
`BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad` beschreibt einen Kommentar, der
einen Fehlerpfad zusagt; hier ist der Träger der Behauptung der Testname samt Meldungen.
