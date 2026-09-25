**Vorgang:** slice-backfill-speicher-untersuchung (Review F-1 und F-2, beide HIGH; Verifikation V-3)

**Fund:** Handbuch und Messbericht nannten für `GOGC=25` eine Senkung „um 18 bis 21 %“; aus den Zeilen des Zeilen-Dokuments gerechnet (Reihe K gegen Reihe A, drei Paare) sind es 11,3 bis 21,2 %. Der Höchstwert je Change (1,57 KiB) und seine Folgezahlen teilten die Spitze der beiden Runs 3 durch die Zahl der Changes **nach** dem Run statt der 2.000.000 davor; nachgerechnet 1,585 KiB (3.096,6 MiB × 1.024 / 2.000.000). Das erste Zählfeld des Suchlaufs (280 Zeilen, `docs/user` 71) galt für einen früheren Stand des Diffs (am Stand `989beef3`: 288 und 79), und „15 Reihen“ zählte die Reihen vor der Fixrunde (21 Überschriften am Stand `989beef3`). Gefunden haben es der Reviewer (F-1, F-2) und der Verifier (V-3) durch Nachrechnen gegen die gedruckten Zeilen.

Quelle: `docs/reviews/review-slice-backfill-speicher-untersuchung.md` (F-1, F-2) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-backfill-speicher-untersuchung.md` (§3 Zeilen F-1 und F-2, §10 V-3). <!-- d-check:status-provenance -->
