**Vorgang:** slice-backfill-speicher-untersuchung (Messbericht §3.4, §3.5 und §4, Verifikation §2 Zeile 2)

**Fund:** Die Speicher-Spitze des Feed-Containers nach einem Backfill hängt nicht am Block: die Spitze des Go-Prozesses im Run (`anon`) liegt bei der Blockgröße 100, 1.000 und 10.000 bei 8,4 bis 8,6, 9,8 bis 10,1 und 24,3 bis 24,4 MiB (Reihen F, E, G, Stufe 200.000 Zeilen, `n` = 3 je Größe, gemessen); breite Zeilen (1.273 statt 73 Bytes) heben sie bei der Blockgröße 1.000 um etwa 5,5 MiB (15,2 bis 15,7 gegen 9,8 bis 10,1 MiB, Reihe I gegen Reihe E; die Differenz ist abgeleitet). Die Ursache der Spitze ist der Retention-Lauf, der alle Changes der Quelle samt Row Images liest (Schalter „Bereinigung aus“: flach). Die Grenze „Block zählt Zeilen, nicht Bytes“ bleibt eine Aussage der Port-Doku; Zeilen im MB-Bereich sind ungemessen.

Quelle: `docs/reviews/messbericht-slice-backfill-speicher-untersuchung.md` (§3.4, §3.5, §4) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-backfill-speicher-untersuchung.md` (§2 Zeile 2). <!-- d-check:status-provenance -->
