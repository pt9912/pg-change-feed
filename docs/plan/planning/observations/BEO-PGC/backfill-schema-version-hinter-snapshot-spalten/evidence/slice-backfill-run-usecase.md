**Vorgang:** slice-backfill-run-usecase (Review-Fund F-8)

**Fund:** Der Reviewer las `service.go` (`currentVersion`) gegen `mapper.go` (`observeRelation`, Nachtrag der Spaltenform): eine Tabelle, die seit der Aktivierung oder seit einem `ALTER TABLE … ADD COLUMN` keine WAL-Änderung hatte, bekommt Backfill-Changes, deren Bild neuere Spalten trägt als ihre `schema_version`, oder die auf eine Version ohne `TableSchema` verweisen. Das akzeptierte Negativ von `ADR-0111` nennt den Fall „während des Runs", nicht die Fälle davor (aus dem Code abgeleitet, nicht als Lauf gemessen). Kein Fixrunden-Anlass, gemeldet an den Architect.

Quelle: `docs/reviews/review-slice-backfill-run-usecase.md` (F-8) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-backfill-run-usecase.md` (§10 Verdikt, Übergabe). <!-- d-check:status-provenance -->
