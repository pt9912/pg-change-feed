**Vorgang:** slice-backfill-e2e (Review-Fund F-8, Fixrunde, Verifikation §5)

**Fund:** Die Phase DDL-Fenster des E2E-Laufs druckte „failed (transient: Fehlerklasse transient: Quelle für den Tabellen-Snapshot vorübergehend nicht verfügbar: Tabelle public.feed_e2e_backfill_rewrite wurde nach dem Snapshot-Export umgeschrieben; ein neuer Antrag beginnt neu)“; die Zeile zu `storage` trug dasselbe Muster (F-8, INFO). Die Fixrunde entfernt die Klassen-Angabe des Fehlerwerts in `failureText` (`internal/application/usecase/backfill/service.go`) und bindet die einmalige Klasse im Test `TestExecuteFailureTextCarriesClassOnce`. Der Verifier mutierte die Entfernung (`strings.Replace(…, "", 0)`): `TestExecuteFailureTextCarriesClassOnce` rot (M4); die gedruckte Zeile beider CI-Legs des Laufs `36065957210` trägt die Klasse einmal („failed (transient: Quelle für den Tabellen-Snapshot vorübergehend nicht verfügbar: …)“).

Quelle: `docs/reviews/review-slice-backfill-e2e.md` (F-8) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-backfill-e2e.md` (§4 M4, §5 F-8). <!-- d-check:status-provenance -->
