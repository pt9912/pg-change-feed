# BEO-PGC/schema-evolution-nicht-dynamisch

**Sub-Area:** Replication-Decoder/Mapper (Schema-Evolution entlang
[`ADR-0015`](../../../../adr/0015-schema-evolution.md); Sub-Area-Kürzel
`PGC` aus der Modus-Deklaration)

Die Beobachtung: [`ADR-0015`](../../../../adr/0015-schema-evolution.md)
(Accepted, `permanent`) hat Option C beschlossen — `TableSchema`-/
`SchemaVersion`-Modelle je Change, `SchemaStorePort` als Outbound-Port-
Folgepflicht, Fehlerklasse `schema` für nicht sicher interpretierbare
Änderungen — und Option A („Relation Metadata 1:1 durchreichen")
ausdrücklich verworfen. Der ausgelieferte Code verhält sich aber wie das
verworfene Option A: `mapper.TableBinding.SchemaVersion` ist bei
Aktivierung statisch gebunden und ändert sich nie zur Laufzeit;
`mapper.Assembler.Consume` verwirft `*decode.Relation`-Ereignisse
ungeprüft; `decode.tupleValues` liest jeden Spaltenwert nur als Text ohne
Typ-/OID-Auswertung. Zwei Lastenheft-Akzeptanzkriterien sind damit
strukturell unerfüllbar: `LH-FA-SCH-005`s Boundary (zwei Changes vor/nach
einer Schemaänderung müssen sich anhand ihrer Schema-Version
unterscheiden lassen) und `LH-FA-SCH-004`s Negative-Fall (eine
inkompatible Typänderung muss erkennbar gemeldet werden). Die
`SchemaStorePort`-Folgepflicht aus `ADR-0015` existiert im gesamten Repo
nicht (`grep -rln SchemaStorePort` → 0 Treffer).

Deklaration: `internal/adapters/driving/replication/mapper/mapper.go`
(`TableBinding.SchemaVersion`, `Assembler.Consume`),
`internal/adapters/driving/replication/decode/decode.go` (`tupleValues`),
`internal/bootstrap/wiring.go:328-332`.
