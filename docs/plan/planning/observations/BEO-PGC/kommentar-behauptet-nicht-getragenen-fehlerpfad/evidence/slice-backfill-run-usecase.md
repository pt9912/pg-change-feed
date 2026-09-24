**Vorgang:** slice-backfill-run-usecase (Review-Fund F-12)

**Fund:** Der Kommentar in `conclude` (`internal/application/usecase/backfill/service.go`) sagte zu, ein `queued`-Run „wird nach einem Neustart aufgenommen" und ein nicht festgehaltener Run bleibe „bis zum Abgleich beim Prozessstart `running`". Beide Aussagen betreffen den Fehlerpfad des Runs und werden vom Code dieses Diffs nicht getragen: Aufnahme und Abgleich liegen bei `slice-backfill-sql-administration`, hier steht nur der Port-Vertrag `InterruptRunning`. Die Aussagen stimmen mit der Architektur-Sicht und `ADR-0113` Festlegung 2 überein, der Kommentar trug keinen Rang-Zeiger darauf. Der Plan führte die Klasse in §8 als einschlägig; gefunden hat es der Reviewer. Behoben in der Fixrunde: der Kommentar nennt die Zusage des Use Cases und zeigt auf `ADR-0113` Festlegung 2.

Quelle: `docs/reviews/review-slice-backfill-run-usecase.md` (F-12) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-backfill-run-usecase.md` (V-3, §8 Zeile §3.7). <!-- d-check:status-provenance -->
