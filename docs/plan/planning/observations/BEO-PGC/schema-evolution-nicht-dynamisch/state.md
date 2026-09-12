Zustand: **verkörpert** — die drei Bausteine (`SchemaStorePort`/
`TableSchema`, dynamische Re-Versionierung, Typ-Auswertung/Fehlerklasse
`schema`) sind real geliefert und dreifach gegen den Compose-Stack
verifiziert: `LH-FA-SCH-005`s Boundary und `LH-FA-SCH-004`s Negative-Fall
sind erfüllt. Der Code verhält sich wie das in
[`ADR-0015`](../../../../adr/0015-schema-evolution.md) beschlossene
Option C — verkörpert in
`internal/adapters/driving/replication/mapper/mapper.go`
(`Assembler.observeRelation`, `classifyRelationColumns`) `seit slice-033`.
Zähler (abgeleitet): 2× (evidence/slice-030.md, evidence/slice-033.md) —
unter der 3×-Schwelle, aber direkt aufgelöst: der Fund
(`slice-030`) und seine vollständige Behebung (`slice-033`, nach der
Zwischenstufe `slice-032`) tragen dieselbe Beobachtung, analog zu
`BEO-PGC/spec008-replication-luecke`s direkter Auflösung bei 2×
(Modul 5 „Offene Risiken werden bei Closure aufgelöst" gilt sinngemäß).
`slice-032` selbst legte keinen eigenen Beleg an (Zwischenschritt ohne
eigenständigen Fund).
