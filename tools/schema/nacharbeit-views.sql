-- Berichtete manuelle Nacharbeit am Schema-Rollout (ADR-0043,
-- Re-Evaluierungs-Trigger): die drei SQL-Views nach LH-FA-SST-002 sind die
-- benannte Grenze des neutralen Schemamodells (tools/schema/schema.yaml) —
-- d-migrate 1.2.0 (gepinnt) kann eine `CREATE VIEW` nicht überführen, deren
-- Katalogform (`pg_get_viewdef`) von der Autorenform abweicht; der
-- Rollout-Lauf konvergiert am Post-execute-Vergleich nicht (derselbe
-- `raw-sql-text-drift`-Befund wie bei `chk_change_operation`,
-- `nacharbeit-operation-check.sql`). Diese Datei trägt die Views als
-- Schritt des `make schema-rollout`-Laufs; `CREATE OR REPLACE VIEW` macht
-- den Schritt wiederholbar.
CREATE OR REPLACE VIEW cdc.active_tables AS
SELECT
    st.source_table_id,
    st.source_id,
    st.schema_name,
    st.table_name,
    (
        SELECT max(sv.version)
        FROM cdc.schema_version sv
        WHERE sv.source_table_id = st.source_table_id
    ) AS current_schema_version
FROM cdc.source_table st;

CREATE OR REPLACE VIEW cdc.consumer_status AS
SELECT
    c.consumer_id,
    c.name,
    cp.source_id,
    cp.acknowledged_position,
    (
        SELECT max(t.commit_position)
        FROM cdc.transaction t
        WHERE t.source_id = cp.source_id
    ) AS latest_commit_position
FROM cdc.consumer c
LEFT JOIN cdc.consumer_position cp ON cp.consumer_id = c.consumer_id;

CREATE OR REPLACE VIEW cdc.changes AS
SELECT
    t.source_id,
    t.commit_position,
    c.change_id,
    c.transaction_id,
    c.source_table_id,
    st.schema_name,
    st.table_name,
    c.sequence,
    c.operation,
    c.old_data,
    c.new_data,
    c.schema_version,
    t.committed_at
FROM cdc.change c
JOIN cdc.transaction t ON t.transaction_id = c.transaction_id
JOIN cdc.source_table st ON st.source_table_id = c.source_table_id
ORDER BY t.commit_position, c.transaction_id, c.sequence;
