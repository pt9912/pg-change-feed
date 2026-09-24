**Vorgang:** slice-backfill-sql-administration (Review F-4, Behebung durch Verifikation nachgemessen)

**Fund:** Der Handbuch-Abschnitt „Bestand als Backfill überführen" nannte als Voraussetzung eine Login-Identität mit `cdc_admin`-Mitgliedschaft und führte danach die Status-Abfrage auf `cdc.backfill_status`. Der Reviewer führte das SQL gegen den ausgerollten Arbeitsbaum unter einem Login `IN ROLE cdc_admin` aus: „permission denied for view backfill_status" (`cdc_admin` trägt kein `SELECT` auf die View, nur `cdc_reader`). Behoben in der Fixrunde: der Abschnitt nennt `cdc_reader` (`CDC_READER_DSN`) als Lese-Identität und sagt, dass `cdc_admin` kein `SELECT` auf die View trägt.

Quelle: `docs/reviews/review-slice-backfill-sql-administration.md` (F-4) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-backfill-sql-administration.md` (§4 Zeile F-4). <!-- d-check:status-provenance -->
