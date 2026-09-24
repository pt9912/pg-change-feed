**Vorgang:** slice-backfill-e2e (Review F-6, Fixrunde)

**Fund:** Der Review maß mit einem Scratch-Test, dass `importSnapshot` mit einem Kontext von
2 s gegen eine fremde Transaktion mit `ACCESS EXCLUSIVE` nach 2,001 s mit der Klasse
`transient` („Tabellensperre: timeout: context deadline exceeded“) zurückkehrt und keine
Sitzung des Runs hinterlässt; ohne Kontext-Limit wartet der Run unbegrenzt. Die Fixrunde
bindet den Kontext-Abbruch in einem committeten Store-Tier-Test und nennt die Gegenrichtung
im Handbuch; eine Zeitgrenze des Runs ist nicht eingeführt.

Quelle: `docs/reviews/review-slice-backfill-e2e.md` (F-6) <!-- d-check:status-provenance -->
· Slice-Plan `slice-backfill-e2e` §6. <!-- d-check:status-provenance -->
