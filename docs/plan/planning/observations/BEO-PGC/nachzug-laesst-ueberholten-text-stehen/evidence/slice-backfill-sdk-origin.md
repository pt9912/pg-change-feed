**Vorgang:** slice-backfill-sdk-origin (Review F-6, Verifikation §3 Zeile F-6)

**Fund:** Der Handbuch-Absatz zum Python-Package sagte, SSE und der NATS-Vollinhalts-Stream „folgen im selben Folge-Release“, obwohl die Python-README und die Package-Zeile im Pflichtenheft beide Zustellwege als geliefert führen (`0.2.0`). Der Satz stammt aus einer früheren Lieferung des Packages, deren Nachzug den Absatz nicht las; er lag außerhalb des Diffs dieses Slice und ist keine bewegte Eigenschaft von `origin`, sondern ein vorbestehender Träger, den der Reviewer beim Nachlesen der SDK-Absätze fand (INFO, Zuständigkeit Planner). Der Nachzug lief in der Fixrunde; der Verifier maß danach die Zukunftsaussagen im Handbuch (`git grep -i 'folgen\|im selben Folge-Release\|noch nicht'`: je 12 Zeilen an `e80b4f64` und `HEAD`, die vier Historie-Zeilen und die acht Betriebs-Zustände ohne Package-Bezug) und in den drei READMEs (0).

Quelle: `docs/reviews/review-slice-backfill-sdk-origin.md` (F-6) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-backfill-sdk-origin.md` (§3 Zeile F-6). <!-- d-check:status-provenance -->
