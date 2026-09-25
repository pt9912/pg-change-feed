Zustand: offen — Ausgang: **weiter offen**, adressiert. Die Grenze steht in der
Port-Doku (`internal/application/port/outbound/tablesnapshot.go`, `NextBlock`) und im
Kommentar an `DefaultBlockSize`; der Use Case des Runs nennt sie ebenfalls (Doc-Kommentar von
`BackfillTableService`, DoD-Punkt „Zeilenzahl im Speicher" in `slice-backfill-run-usecase`;
`TestExecuteStreamsBlocks` bindet die Streaming-Reihenfolge). Träger der Messung ist der Folge-Slice
`slice-backfill-bench-richtgroesse` (Plan in `done/`, §3 „Übergabe aus
`slice-backfill-snapshot-reader`": `B` und Zeilenbreite im Ursprung jeder Kopier-Zahl,
breite Zeilen ohne eigenen Lauf ungemessen); ein Bytelimit oder eine Zeilenbreiten-Wache
wäre Sache der Ausbaustufe für Durchsatz (`ADR-0111` Re-Evaluierung). Zähler
(abgeleitet): 1× (evidence/slice-backfill-snapshot-reader.md).

Messung des Folge-Slices `slice-backfill-bench-richtgroesse` (Vertrag
`harness/targets/bench-backfill.md`, Handbuch §Grenzwerte): der Bench nennt `B` und
die Zeilenbreite (etwa 74 Bytes) im Ursprung jeder Kopier-Zahl; breite Zeilen bleiben
ungemessen. Die Speicher-Spitze des Feed-Containers in der Stufe mit 200.000 Zeilen
liegt im Run bei 401,7 bis 467 MiB und 20 s nach dem letzten Run bei 368,6 bis
641,7 MiB (vier Läufe; zwei davon übernommen, Aufstellung im Handbuch und in der
Verifikation §3) gegen einen Bedarf von etwa 74 KB je Block; die Ursache ist nicht
untersucht. Adresse: der Lese-Schritt der Closure von `welle-backfill-bestand`
(Architect-Entscheidung, ob ein Untersuchungs-Slice entsteht); der Eintrag bleibt
offen.
