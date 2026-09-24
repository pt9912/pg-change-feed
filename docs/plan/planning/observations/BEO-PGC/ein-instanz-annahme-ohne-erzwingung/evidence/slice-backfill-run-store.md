**Vorgang:** slice-backfill-run-store (Review F-5, Verifikation §5 Zeile `ADR-0113` Festlegung 1 Punkt 2)

**Fund:** Der Review fand im Annahme-Adapter weder eine Sperre noch eine Unique-Kante auf aktive Runs (`cdc.backfill_run`: Primärschlüssel auf `run_id`, kein Index auf Quelle, Schema und Tabelle für aktive Runs) und keinen Test mit zwei gleichzeitig annehmenden Verbindungen; der Godoc, `ADR-0113` und der Slice-Plan nennen die Annahme (eine Goroutine, eine Instanz je Quelle) ausdrücklich, die Vereinbarkeit mit der ADR ist gegeben (F-5, INFO). Der Verifier bestätigte den Befund am Code (kein Lock, keine Unique-Kante) und hielt fest, dass er zwei gleichzeitig annehmende Verbindungen **nicht** gemessen hat. Der Slice-Plan führte den Punkt als Risiko mit Ausgang zur Closure; der Ausgang ist *weiter offen*.

Quelle: `docs/reviews/review-slice-backfill-run-store.md` (F-5) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-backfill-run-store.md` (§5 Zeile `ADR-0113` Festlegung 1 Punkt 2, §8 Punkt 2). <!-- d-check:status-provenance -->
