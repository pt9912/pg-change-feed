Zustand: offen — Ausgang: **weiter offen**, adressiert. Träger der Messung ist der Folge-Slice
`slice-backfill-bench-richtgroesse` (Plan in `done/`, §3 „Übergabe aus
`slice-backfill-run-store`": die Kopierdauer je Tabellengröße schließt die Einfügeform ein, eine
Blockdauer nahe oder über der Frist steht als Befund im Bericht). Die Fristen tragen die
Kennzeichnung „Startwerte ohne Messung" im Kommentar an den Konstanten; die Einfügeform trägt sie
am Code nicht (Review F-7). Verwandt: `BEO-PGC/blockgroesse-zaehlt-zeilen-nicht-bytes`
(dieselbe Messung, andere Größe: Speicher je Block). Zähler (abgeleitet): 1×
(evidence/slice-backfill-run-store.md).

Messung des Folge-Slices `slice-backfill-bench-richtgroesse`: die mittlere Blockdauer
liegt bei etwa 0,12 bis 0,13 s (abgeleitet: Dauer durch Blockzahl, Handbuch
§Grenzwerte) und damit weit unter der Frist von 5 min je Block; die längste Dauer
eines einzelnen Blocks und die Dauer des Commits sind nicht gemessen, die Einfügeform
(zeilenweise, eine Schreibtransaktion) ist im Ursprung jeder Kopier-Zahl genannt.
Die Fristen und die Einfügeform bleiben Startwerte ohne Messung der Spitze; Adresse: der
Re-Evaluierungs-Trigger „Kopierdauer über der Betriebs-Toleranz“ der Entscheidung zum
Backfill-Bestand.
