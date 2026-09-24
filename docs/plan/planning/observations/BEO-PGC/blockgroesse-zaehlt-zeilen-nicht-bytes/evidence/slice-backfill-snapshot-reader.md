**Vorgang:** slice-backfill-snapshot-reader (Nachprüfungs-Review F-8, Planner-Closure)

**Fund:** Der Nachprüfungs-Review hielt fest, dass die Grenze `B` Zeilen zählt und der Block bis zur Rückgabe doppelt vorliegt (F-8, INFO, aus dem Treiber-Quelltext abgeleitet, nicht gemessen); die zweite Fixrunde zog die Port-Doku und den Kommentar an `DefaultBlockSize` nach. Der Plan führt das Risiko in §6 mit dem Ausgang *weiter offen*: eine Messung mit breiten Zeilen gehört in die Ausbaustufe für Durchsatz; der Bench-Slice trägt die Übergabe.

Quelle: `docs/reviews/review-slice-backfill-snapshot-reader-fixrunde.md` (F-8) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-backfill-snapshot-reader.md` (§3.7, Zeile F-8). <!-- d-check:status-provenance -->
