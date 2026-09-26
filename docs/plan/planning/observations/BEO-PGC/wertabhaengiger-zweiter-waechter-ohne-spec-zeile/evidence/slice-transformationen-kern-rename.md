**Vorgang:** slice-transformationen-kern-rename (Verifikations-Fund V-4, Review-Fund F-5)

**Fund:** Der zweite Wächter in `model.BuildRowImage` löst nur aus, wenn beide Quellspalten zweier Regeln mit gleichem Zielnamen in derselben Zeile einen Wert tragen (`TestConsumeTwoRulesWithSameTargetAreNotApplicableWhenBothMatch`, Gegenprobe mit einer wertlosen Quellspalte: die Änderung geht durch). `SPEC-030` (Anwendbarkeit) sagt: „nie am Wert einer Zeile“. Der Verifier hat den Punkt als INFO geführt und den Widerspruch zur Spec verneint (K3 am Antrag schließt die Konstellation aus, der Wächter ist der zweite); der Kommentar am Sentinel benennt die Grenze („an einer Zeile, die beide Werte trägt“). Vor der Fixrunde stand an dieser Stelle gar kein Wächter (Review F-5: zwei Regeln mit gleichem Ziel ergaben `{"id":"1","z":"x","z":"y"}` ohne Fehler).

Quelle: `docs/reviews/verifikation-slice-transformationen-kern-rename.md` (V-4) <!-- d-check:status-provenance -->
· `docs/reviews/review-slice-transformationen-kern-rename.md` (F-5) <!-- d-check:status-provenance -->
· `internal/adapters/driving/replication/mapper/transformation_test.go`.
