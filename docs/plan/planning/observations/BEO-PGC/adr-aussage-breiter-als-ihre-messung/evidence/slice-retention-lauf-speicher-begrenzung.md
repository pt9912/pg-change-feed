**Vorgang:** slice-retention-lauf-speicher-begrenzung (Verifikation V-3, INFO)

**Fund:** Die dritte Zeile der Fitness-Function-Tabelle der `Accepted`-ADR `ADR-0124` erwartet, dass die Spitze nach dem Run bei 1.000.000 bis 3.000.000 Changes „im Bereich der Spitze im Run (8,9 bis 10,5 MiB) plus dem Bedarf einer Seite“ liegt; die Erwartung ist hergeleitet, die Zeile ist als „Messung ohne Pass/Fail“ geführt. Gemessen sind 14,9 bis 17,6 MiB (`memory.peak`, zwei Läufe) und 15,5 bis 17,3 MiB im Verifier-Lauf, der `anon`-Wert der Spitze im Run liegt bei 10,0 bis 13,2 MiB. Wörtlich ist die Erwartung nicht erfüllt, der Form nach (flach über der Zahl der Changes, Takte laufen weiter) ist sie es. Der Slice setzte die ADR korrekt um; die Aussage war breiter als die Menge der Messgrößen, an der sie geprüft war. Kein Supersede: keine der sechs Festlegungen hängt an der Zahl, die Lesart steht in der Closure-Notiz des Slice.

Quelle: `docs/reviews/verifikation-slice-retention-lauf-speicher-begrenzung.md` (§4, §9 V-3) <!-- d-check:status-provenance -->
· `docs/reviews/messbericht-slice-retention-lauf-speicher-begrenzung.md` (Abschnitt 7.1). <!-- d-check:status-provenance -->
