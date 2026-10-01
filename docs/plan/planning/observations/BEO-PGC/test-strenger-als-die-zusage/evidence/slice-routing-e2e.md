**Vorgang:** slice-routing-e2e (Review F-2, MEDIUM)

**Fund:** `TestRunStreamWithRetrySlotStillActive` prüft `len(delivered[2]) != 1` („erwartet genau eine") für den zweiten Versuch des Retry-Zyklus. Beim ersten Lauf von `make test-replication` an PostgreSQL 17 lieferte der zweite Versuch zwei Positionen (`[94993648 95321768]`), die Wiederholung war grün. Die bestätigte Position einer Transaktion ist ihre `CommitLSN`; liegt der Slot-Stand beim Aufbau des zweiten Versuchs genau darauf, stellt der Server die bereits persistierte Halter-Transaktion erneut zu — unter At-least-once zulässig, die Persistierung ist idempotent. Ob der Stand schon darüber hinaus gerückt ist, entscheidet das Zeitverhalten der Leerlauf-Bestätigung des Halters (hergeleitet, nicht reproduziert). Der Test liegt in einem anderen Tier als der Slice (`make test-replication`, im Diff unverändert) und beeinflusst dessen Belege nicht.

**Form (Ausprägung):** Erwartung genauer als die Zusage (Anzahl „genau eine" statt „die Retry-Change einmal persistiert"). Schwere MEDIUM; vor dem Merge vom Reviewer gefunden, kein Fix im Slice (anderer Tier), Folge-Slice angelegt.

Quelle: `docs/reviews/review-slice-routing-e2e.md` (F-2) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-routing-e2e.md` (§5 Zeile F-2). <!-- d-check:status-provenance -->
