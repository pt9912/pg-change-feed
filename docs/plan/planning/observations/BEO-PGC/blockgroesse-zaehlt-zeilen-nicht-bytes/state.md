Zustand: **verkörpert** — Ausgang: **verkörpert** →
`internal/application/usecase/retention/service.go` (Kandidaten seitenweise ohne Row Images,
`ReadRetentionCandidates` am `ChangeStorePort`) · seit slice-retention-lauf-speicher-begrenzung
(Entscheidung `ADR-0124`): die Ursache der Speicher-Spitze des Feed-Containers nach einem
Backfill ist der Retention-Lauf, nicht der Block (Messbericht
`messbericht-slice-backfill-speicher-untersuchung` §1 und §4; Schalter „Bereinigung aus“:
flach). Beleg der Behebung: Spitze 14,9 bis 17,6 MiB bei 1.000.000 bis 3.000.000 Changes
(`messbericht-slice-retention-lauf-speicher-begrenzung` Abschnitt 3; Verifier-Lauf
15,5 bis 17,3 MiB, `verifikation-slice-retention-lauf-speicher-begrenzung` Abschnitt 4).
Herkunft: `slice-backfill-speicher-untersuchung`, seit welle-backfill-bestand
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
