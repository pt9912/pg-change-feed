**Vorgang:** slice-backfill-run-store (Review F-7, Plan §6 „Gemeldet — Startwerte ohne Messung")

**Fund:** Der Implementer meldete die Fristen der Adapter (30 s je Zustands-Operation, 5 min je Block und Commit) und die zeilenweise Einfügung eines Blocks als Setzungen ohne Messung; die Fristen tragen die Kennzeichnung im Kommentar an den Konstanten (`backfill.go`), die Einfügeform steht nur im Plan (F-7, INFO: ein Leser von `AppendBlock` sieht die Grenze nicht). Die Messung gehört zu `slice-backfill-bench-richtgroesse`, dessen Plan sie bis zur Closure nicht trug; die Closure zog sie als „Übergabe aus `slice-backfill-run-store`" in dessen §3.

Quelle: `docs/reviews/review-slice-backfill-run-store.md` (F-7). <!-- d-check:status-provenance -->
