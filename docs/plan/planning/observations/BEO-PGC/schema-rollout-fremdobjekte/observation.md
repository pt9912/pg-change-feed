# BEO-PGC/schema-rollout-fremdobjekte

**Sub-Area:** Schemamigration (d-migrate; Sub-Area-Kürzel `PGC` aus der
Modus-Deklaration)

Die Beobachtung: Ein zweiter `make schema-rollout`-Lauf gegen eine
bereits migrierte DB blockiert mit Exit 8
(`DESTRUCTIVE_OPERATION_REQUIRES_CONFIRMATION`, `DropView`) auf
`cdc.heartbeat`/`cdc.metrics` — beide Views entstehen außerhalb des
neutralen Modells (`tools/schema/schema.yaml`), über separate
psql-Nacharbeit-Schritte (`tools/schema/nacharbeit-heartbeat.sql`,
`tools/schema/nacharbeit-observability.sql`). d-migrate sieht sie als
nicht-deklarierte Fremdobjekte und plant ihren Abbau. Pin-unabhängig:
real gegen d-migrate 1.3.0 und 1.3.1 reproduziert
([`ADR-0043`](../../../../../plan/adr/README.md)).
