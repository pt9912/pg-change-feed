# Beleg: slice-082

Vorgang: `slice-082` — der Fixture-Nachzug, der `make test-replication` wieder
grün macht.

Fund: **Zwei** Stellen im Plan des eigenen Slice, beide vom Reviewer bzw.
Verifier am Gegenstand widerlegt — der Code-Diff war korrekt, die Aussage über
ihn nicht.

1. Das DoD-Kriterium behauptete, der Tier-Lauf beweise die reale Anwesenheit
   **beider** Tabellen. Er beweist eine: `cdc.administration_request` ist für
   `bootstrap.Run` fatal (`42P01` beendet den Lauf), `cdc.process_heartbeat`
   wird **best-effort** geschrieben (`internal/bootstrap/wiring.go:844`,
   `_ = port.Beat(...)`) — nimmt man nur sie weg, bleibt der Lauf grün.
2. Dasselbe Kriterium nannte „15 Tabellen, 6 Sichten, 4 Funktionen". Die Zahl
   war aus dem **Review-Report** übernommen und nie selbst gemessen; der
   Rollout-Report (`tools/schema/plan.yaml`) führt **10** Tabellen und **4**
   Sichten.

Quelle: `docs/reviews/review-slice-082.md` F-1 · `docs/reviews/verify-slice-082.md`
V-1 · `tools/schema/plan.yaml` (Report des Rollouts) · `internal/bootstrap/wiring.go:844`.
