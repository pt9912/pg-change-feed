**Vorgang:** architect-verdict-welle-transformationen-offene-fragen (Architect-Zug, Messung am laufenden Feed-Container)

**Fund:** `runStreamAfterAdministrationPass` in `internal/bootstrap/wiring.go` ist eine extrahierte Sequenz (Vorlauf, Goroutinen-Start, Stream) mit Whitebox-Tests und Fakes; die Aufrufstelle in `Run` ist über zwei Quelltext-Tests gebunden. Die Tests zeigen die Ordnung, nicht den Zustand der Verbindung zwischen Aufbau und Aufruf: `receive.NewStream` sendet `START_REPLICATION` (`receive.go`) lange vor der Sequenz, der Server beendet einen nicht gelesenen Strom nach `wal_sender_timeout`, und der Prozess endet nach einem Vorlauf über dieser Zeit mit der Klasse `replication`. Gefunden durch eine Messung am komponierten Prozess (drei Läufe, `ADR-0128` §Gemessen), von keinem Unit-Test und keinem Review.

**Form (Ausprägung):** die Klasse „Unit-Test mit Fakes verdeckt eine Lücke der realen Verdrahtung“ in einer neuen Ausprägung: nicht ein `nil`-Use-Case am Adapter, sondern eine **Zeit- und Zustandseigenschaft der komponierten Verbindungen**, die kein Fake abbildet. Schwere MEDIUM, gefunden nach dem Merge des Slice, der die Sequenz lieferte.

Quelle: `docs/reviews/architect-verdict-welle-transformationen-offene-fragen.md` (§2.1, §8) <!-- d-check:status-provenance -->
· `docs/plan/adr/0128-prozessstart-vorlauf-frist-und-beginn-des-replikationsstroms.md` (§Gemessen). <!-- d-check:status-provenance -->
