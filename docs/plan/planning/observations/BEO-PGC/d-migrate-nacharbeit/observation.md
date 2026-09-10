# BEO-PGC/d-migrate-nacharbeit

**Sub-Area:** Schemamigration (d-migrate; Sub-Area-Kürzel `PGC` aus der
Modus-Deklaration)

Die Beobachtung: d-migrate 1.2.0 konvergiert am CHECK-Ausdruck mit
String-Literalen nicht (`raw-sql-text-drift`, upstream offener Issue) —
der Constraint `chk_change_operation` läuft als berichtete manuelle
Nacharbeit (`nacharbeit-operation-check.sql`) im Rollout-Lauf statt im
neutralen Modell (`tools/schema/schema.yaml`).

Deklaration: `tools/schema/schema.yaml` (Benannte Grenze),
`Makefile` (psql-Nacharbeit-Schritt im `schema-rollout`-Target);
[`ADR-0043`](../../../../docs/plan/adr/README.md) trägt die
Re-Evaluierungs-Regel (Ausweichform/Retirement).
