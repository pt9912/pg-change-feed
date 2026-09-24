Zustand: offen — Ausgang: **weiter offen** (unter der 3×-Schärfungsschwelle).
Ursachenbehebung (die zwei vorbestehenden `pool.Query`-Stellen in
`internal/adapters/driven/postgresstorage/sqlviews_test.go` auf `defer rows.Close()`
umstellen und die übrigen Store-Tests auf offene `Rows` auditieren) berührt
bestehende Testdateien; ein Slice dafür existiert nicht. Die Herkunft ist
übernommen, nicht reproduziert (siehe `observation.md`). Gelesen wird der Eintrag
im Sichtungs-Schritt der Slice-Planung. Zähler (abgeleitet): 1×
(evidence/slice-backfill-change-origin.md).
