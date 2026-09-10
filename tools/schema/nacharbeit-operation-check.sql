-- Berichtete manuelle Nacharbeit am Schema-Rollout (ADR-0043,
-- Re-Evaluierungs-Trigger): der CHECK über der Operation ist die benannte
-- Grenze des neutralen Schemamodells (tools/schema/schema.yaml) — d-migrate
-- 1.2.0 (gepinnt) kann den Ausdruck nicht überführen; der Rollout-Lauf
-- konvergiert am Post-execute-Vergleich nicht. Diese Datei trägt den
-- Constraint als Schritt des `make schema-rollout`-Laufs, und die
-- handgeschriebene DDL
-- (internal/adapters/driven/postgresstorage/schema.sql) trägt dieselbe Form.
-- Die Form gilt mit dem Ablauf der Nacharbeit auch ohne d-migrate-Zug: der
-- DROP-Guard macht den Schritt wiederholbar.
ALTER TABLE cdc.change DROP CONSTRAINT IF EXISTS chk_change_operation;
ALTER TABLE cdc.change
    ADD CONSTRAINT chk_change_operation
    CHECK (operation IN ('INSERT', 'UPDATE', 'DELETE'));
