**Vorgang:** slice-routing-kern-label (Review F-1, MEDIUM; siebzehnte Datei des Eintrags)

**Fund:** Der Slice fügte `cdc.changes` die Spalte `route_target` als **letzte** Spalte hinzu. Der Godoc von `SelectChanges` (`internal/adapters/driven/postgresstorage/queries/queries.go`) sagte weiter „`origin` steht als letzte Spalte“ und im nächsten Satz „`route_target` steht danach als letzte Spalte“ — zwei „letzte“ im selben Block. Der Test-Kommentar `internal/adapters/driven/postgresstorage/sqlexec/translate_test.go` („Die letzte Spalte der Projektion trägt die Herkunft“) in einer vom Diff geänderten Datei war seit dem Diff falsch; `sqlviews_test.go` hatte die Berichtigung („vorletzte“) bekommen, diese beiden Stellen nicht. Die Fixrunde (`2629d544`) zog beide nach; der Verifier schloss F-1 am Text (`git grep -n 'letzte Spalte' -- internal tools`).

**Form (Ausprägung):** dieselbe Klasse am Träger **Go-Kommentar** (Godoc und Testkommentar): ein hinzugefügter Satz macht den Nachbarsatz im selben Block und eine Stelle in einer vom Diff berührten Nachbardatei falsch. Gefunden mit der Lese-Probe des Reviewers (Kontext um die hinzugefügten Zeilen); vor dem Merge gefunden, in der Fixrunde behoben, Ausgang unverändert **verkörpert**.

Quelle: `docs/reviews/review-slice-routing-kern-label.md` (F-1) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-routing-kern-label.md` (§5 Zeile F-1). <!-- d-check:status-provenance -->
