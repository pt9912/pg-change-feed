# BEO-PGC/walsender-wirksamkeit

**Sub-Area:** Replication-Stream (Walsender-Verwaltung; Sub-Area-Kürzel
`PGC` aus der Modus-Deklaration)

Die Beobachtung: Der Publication-Entzug (Disable) wirkt am laufenden
Walsender erst nach dessen Neuaufbau (PostgreSQL-Verhalten); der
Katalog-Beleg (`pg_publication_tables`) trägt die Bindung, nicht das
Capture-Zeitverhalten eines laufenden Streams. Der Happy Path von
Disable ist für neue/neugestartete Feeds belegt, nicht für einen live
laufenden Stream.

Deklaration: `internal/application/port/outbound/tableactivation.go`
(grenze-Kommentar am Publication-Entzug);
[`ADR-0028`](../../../../../plan/adr/README.md) trägt die
Verwaltungs-Verträge.
