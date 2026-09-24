**Vorgang:** slice-backfill-sql-administration (Review F-5, Behebung durch Verifikation nachgemessen)

**Fund:** Das Handbuch nannte an zwei Stellen (Einleitung des Backfill-Abschnitts, Glossar „Backfill") „die Zeilen, die die Tabelle zum Zeitpunkt des Antrags bereits enthält". Das Pflichtenheft (`LH-FA-CAP-009.a`: „zum Startzeitpunkt eines Runs") und der Mechanismus (`snapshot.go`, Paketkommentar: der Snapshot entsteht mit dem Slot des Runs beim Beginn der Ausführung) tragen den Start des Runs; der Handbuch-Abschnitt „Überlappung" und die Position `X` trugen ihn ebenfalls, die beiden Stellen widersprachen ihm. Behoben in der Fixrunde (Handbuch 1.49).

Quelle: `docs/reviews/review-slice-backfill-sql-administration.md` (F-5) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-backfill-sql-administration.md` (§4 Zeile F-5). <!-- d-check:status-provenance -->
