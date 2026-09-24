**Vorgang:** slice-backfill-e2e (Review-Fund F-4, Fixrunde, Verifikation §5)

**Fund:** In `bf_ddl_window` (Runner `tools/harness/run-integration-tests.sh`) las die Assertion zum Erfassungspfad `coalesce(error_class, '')` aus `cdc.heartbeat` und verglich mit dem leeren String; dieser Wert entsteht bei einer Zeile mit NULL **und** bei keiner Zeile. Eine fehlende Heartbeat-Zeile (Quelle unbekannt, Filter falsch) hielte die Zusage „der Run-Fehler ist run-lokal, der Erfassungspfad läuft weiter“ ebenso grün (F-4, LOW; gelesen, nicht gemutet — ein E2E-Lauf dauert rund 5 Minuten). Die Fixrunde liest `SELECT count(*) FROM cdc.heartbeat WHERE source_id = 'src-e2e' AND error_class IS NULL` gleich 1: bei fehlender Zeile 0, die Assertion färbt rot. Der Verifier prüfte die Form am Diff und beide CI-Legs grün; die Mutation der Eingabeseite (Quellfilter auf eine nicht vorhandene Quelle) fuhr er nicht (Verifikations-Report §12).

Quelle: `docs/reviews/review-slice-backfill-e2e.md` (F-4) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-backfill-e2e.md` (§5 F-4, §12). <!-- d-check:status-provenance -->
