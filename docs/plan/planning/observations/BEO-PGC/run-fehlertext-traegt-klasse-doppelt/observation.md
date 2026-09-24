# BEO-PGC/run-fehlertext-traegt-klasse-doppelt

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die Zusammensetzung des
Fehlertexts eines Backfill-Runs, keine eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Der Fehlertext eines Runs (`error_message`) entsteht aus zwei Schichten,
die je die Fehlerklasse voranstellen — `Fail` setzt „<Klasse>: “ vor den Text, und der
Sentinel-Text des Ports beginnt selbst mit „Fehlerklasse <Klasse>: …“. Der Betreiber liest
„transient: Fehlerklasse transient: Quelle für den Tabellen-Snapshot vorübergehend nicht
verfügbar: …“. Kein Gate liest die Zusammensetzung; sichtbar wurde sie in der gedruckten
Zeile der Phase DDL-Fenster des E2E-Laufs. Der Reviewer fand sie (F-8, INFO): die Form des
Sentinel-Textes bestand bereits, der Diff machte sie sichtbar.

**Warum das zählt:** Die Klasse ist die Handlungsanweisung des Betreibers (`SPEC-008`); ein
Text, der sie zweimal nennt, liest sich wie zwei Befunde, und jede weitere Text-Schicht
kann sie ein drittes Mal voranstellen.

Deklaration: `slice-backfill-e2e` (Review F-8).
