# BEO-PGC/blockgroesse-zaehlt-zeilen-nicht-bytes

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die Speichergrenze des
Snapshot-Lesers beim Backfill, keine eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Der Snapshot-Leser (`TableSnapshotPort`, Adapter `postgressnapshot`)
begrenzt einen Lese-Block auf `B` Zeilen, nicht auf Bytes. Der Speicherbedarf eines
`NextBlock`-Aufrufs ist `B` mal die Zeilenbreite, und der Block liegt bis zur Rückgabe
kurzzeitig doppelt vor (Treiberbytes und Zeichenketten der Werte). Bei Zeilen im
MB-Bereich (`jsonb`, `bytea`) trägt die Zusage „keine unbegrenzte RAM-Haltung"
(`LH-FA-CAP-006.a`) deshalb nur in der Zeilenzahl. Die Aussage ist aus dem Code und dem
Treiber-Quelltext abgeleitet, **nicht gemessen**; der Startwert `B` = 1.000 ist eine
Setzung ohne Messung.

**Warum das zählt:** Eine Speichergrenze in Zeilen ist für breite Zeilen keine
Speichergrenze; ohne Messung mit breiten Zeilen steht offen, ab welcher Zeilenbreite
der Block den Prozess belastet. Kein Gate liest die Zeilenbreite eines Bestands.

Deklaration: `slice-backfill-snapshot-reader`, Risiko §6 neunter Punkt („Blockgröße
zählt Zeilen, nicht Bytes"), Ausgang *weiter offen*.
