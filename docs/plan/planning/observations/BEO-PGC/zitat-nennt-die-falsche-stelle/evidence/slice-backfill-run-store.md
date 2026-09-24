**Vorgang:** slice-backfill-run-store (Review F-3)

**Fund:** Der Kopplungs-Kommentar in `backfillhelpers_test.go` (und die Plan-Zeile `schema.sql` im Slice-Plan) nannte `consumerstate_test.go` neben `store_test.go` und `tableactivation_test.go` als Datei, deren Tests das Schema per `DROP SCHEMA cdc CASCADE` neu aufbauen. Gemessen (`grep -n "DROP SCHEMA" *_test.go`, `grep -ln "newTestStore(\|newTestActivation("`): die Ausführenden sind `store_test.go` und `tableactivation_test.go`; `consumerstate_test.go` enthält weder `DROP SCHEMA` noch einen Aufruf der beiden Helfer. Die Aussage über die Reihenfolge trug (`make test-store` grün), die zitierte Stelle trug sie nicht — ein Verweis auf eine Datei statt auf einen Abschnitt, dieselbe Form wie beim Slice-Kennungs-Zitat. Behoben in der Fixrunde (Kommentar und Plan-Zeile nennen die zwei Dateien).

Quelle: `docs/reviews/review-slice-backfill-run-store.md` (F-3). <!-- d-check:status-provenance -->
