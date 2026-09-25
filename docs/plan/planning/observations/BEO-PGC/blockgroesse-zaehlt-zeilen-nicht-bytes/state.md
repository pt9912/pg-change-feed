Zustand: **geplant** — Ausgang: **geplant** →
`slice-retention-lauf-speicher-begrenzung` (Datei in `open/`): die Ursache der Speicher-Spitze
des Feed-Containers nach einem Backfill ist untersucht und belegt — der Retention-Lauf liest
alle Changes der Quelle samt Row Images, nicht der Block (Messbericht
`messbericht-slice-backfill-speicher-untersuchung` §1 und §4; Schalter „Bereinigung aus“:
flach). Herkunft: `slice-backfill-speicher-untersuchung`, seit welle-backfill-bestand
(Architect-Verdikt `architect-verdict-welle-backfill-bestand-lese-schritt` §4.3 und §5 (h)).

Gegenstand: der Block des Snapshot-Lesers zählt Zeilen, nicht Bytes (Port-Doku `NextBlock`
in `internal/application/port/outbound/tablesnapshot.go` und Kommentar an `DefaultBlockSize`).
Gemessen ist der Blockbedarf bei schmalen Zeilen (Spitze im Run 8,4 bis 8,6, 9,8 bis 10,1 und
24,3 bis 24,4 MiB bei der Blockgröße 100, 1.000 und 10.000, Reihen F, E, G) und der Einfluss
der Zeilenbreite bei der Blockgröße 1.000 (etwa 5,5 MiB mehr bei 1.273 statt 73 Bytes je
Zeile, Reihe I gegen E, abgeleitet). Zeilen im MB-Bereich (`jsonb`, `bytea`) bleiben
ungemessen; die Adresse dieser Messung ist der Re-Evaluierungs-Trigger „Ausbaustufe für
Durchsatz“ der Entscheidung zum Backfill-Bestand (`ADR-0111`).

Zähler: 2× (Dateien unter `evidence/`).
