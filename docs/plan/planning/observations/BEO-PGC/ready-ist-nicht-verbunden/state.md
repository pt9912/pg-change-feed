Zustand: offen (**2×**) — unter der Schwelle (3×), kein Ausgang zugewiesen. Beide bekannten
Träger sind im Vorgang gehärtet: die Filter-Phase (Träger `tools/harness/lib-sdk-filter-fixture.sh`,
Szenario-Dateien der drei Tiers) und die Routing-Phase (Träger `tools/harness/lib-sdk-route-fixture.sh`,
`RouteScenario` der drei Tiers) committen nach `SEEN` eine zweite Gruppe und lassen das
Ruhefenster erst nach `SEEN_SECOND` beginnen. Ein Ausgang *verkörpert* ist bei 2× nicht
zugewiesen; ein drittes Auftreten (eine weitere Phase mit Negativ über Abwesenheit und
Stream-Client) löst den Lese-Schritt der Welle-Closure aus.

Nicht betroffen (hergeleitet beziehungsweise gelesen): die Regel-Phasen der Tiers (kein Negativ
über Abwesenheit), die HTTP-Pull-Flächen (Ruhefenster 0, keine Verbindung), der Server-Rundlauf
`tools/harness/run-integration-tests.sh` (wartet vor der festen Menge auf alle Clients, Review
F-4 des Routing-Slice, gelesen, nicht gefahren).

Zähler (abgeleitet): **2×** (evidence/slice-sdk-sse-filter-phase-verbindung-haertung.md,
evidence/slice-sdk-routing-phase-verbindung-haertung.md).
