Zustand: offen — Ausgang: **weiter offen**, adressiert. Die Grenze steht in der
Port-Doku (`internal/application/port/outbound/tablesnapshot.go`, `NextBlock`) und im
Kommentar an `DefaultBlockSize`. Träger der Messung ist der Folge-Slice
`slice-backfill-bench-richtgroesse` (Plan in `open/`, §3 „Übergabe aus
`slice-backfill-snapshot-reader`": `B` und Zeilenbreite im Ursprung jeder Kopier-Zahl,
breite Zeilen ohne eigenen Lauf ungemessen); ein Bytelimit oder eine Zeilenbreiten-Wache
wäre Sache der Ausbaustufe für Durchsatz (`ADR-0111` Re-Evaluierung). Zähler
(abgeleitet): 1× (evidence/slice-backfill-snapshot-reader.md).
