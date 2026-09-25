Zustand: offen (**1×**) — unter der Schwelle, kein Ausgang zugewiesen. Ein Träger
ist **nicht** vorgeschlagen: ob ein Wert eine Schätzung oder eine Grenze ist,
steht in seiner **Formulierung**, nicht in seinem Wert — ein Sensor müsste
„≈45" von „45" unterscheiden und wüsste nicht, welches davon eine Decke
behauptet.

Zähler (abgeleitet): **1×** (evidence/slice-094.md). Das Erstauftreten fiel im
Review zu `slice-094` auf (F-1, HIGH): die **≈45** der `ADR-0082` wurde auf dem
Weg über den Slice-Plan in die Umsetzung zu einer Grenze und war um vier
Statements zu niedrig — gemessen `cmd` **49 von 49**.

**Nicht zu verwechseln** mit `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`
(dort trägt ein **genannter Beleg** seinen Satz nicht; hier trägt jeder Träger
seinen Satz und zitiert nur nicht als Zitat) und mit
`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (dort driftet ein Wert
**gegen eine Messung**; hier gibt es keine — das ist der Kern).

Der Ort, an dem der Fehler entstehen konnte (die Warn-Richtgröße von
`slice-backfill-bench-richtgroesse`), ist ohne zweites Auftreten geschlossen: jede
Nennung von Richtgröße und Toleranz trägt Lauf, Host und das Wort „Orientierung“
bzw. „Startwert, Setzung ohne Messung“, und keine Stelle lehnt einen Antrag wegen
der Größe ab (Verifikation §2 Nr. 3, §8). Der Zähler bleibt **1×**.

Angewandt ohne Anfall an `slice-retention-lauf-speicher-begrenzung`: der Bedarf einer Seite
(etwa 1,5 MiB, hergeleitet) steht als Orientierung neben der gemessenen Spitze (14,9 bis
17,6 MiB), keine Stelle lehnt oder begrenzt anhand der Herleitung; der Bedarf einer Seite
ist nicht getrennt gemessen. Der Zähler bleibt **1×**.
