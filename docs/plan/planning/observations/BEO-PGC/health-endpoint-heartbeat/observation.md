# BEO-PGC/health-endpoint-heartbeat

**Sub-Area:** Observability (Health-Endpoint; Sub-Area-Kürzel `PGC` aus
der Modus-Deklaration)

Die Beobachtung: [`LH-FA-ADM-002`](../../../../../../spec/lastenheft.md)/
[`LH-QA-OPS-002`](../../../../../../spec/lastenheft.md) (Health-Endpoint)
ist nicht geliefert. Der Architect-Verdikt zu `slice-011`:
Ein Heartbeat-Pattern (Prozess schreibt periodisch in eine
`cdc.process_heartbeat`-Tabelle, eine vierte SQL-View liest sie) deckt
den Bedarf ohne neue ADR — kein Architektur-Loch, sondern eine noch
nicht geschnittene Folge-Slice-Arbeit (Application-/Bootstrap-Schicht,
Timer im Capture-Prozess).

Deklaration: der Architect-Verdikt zu `slice-011`,
slice-011-Plan §1/§5/§6.
