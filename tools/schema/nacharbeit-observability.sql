-- Berichtete manuelle Nacharbeit am Schema-Rollout (ADR-0043, s. o.
-- tools/schema/nacharbeit-operation-check.sql): cdc.metrics trägt das
-- Metriken-Minimum nach LH-FA-SST-004/LH-QA-OPS-003 (SPEC-009) als vierte
-- Lese-View — reine Projektion/Aggregation über bereits persistierte
-- Zeilen, keine Domänenentscheidung (ARC-005 „SQL-Funktionen/Views",
-- ADR-0046 Rollenspaltung). Long-Format (metric_name, label, value): ein
-- Monitoring-System liest die Zeilen roh; die Klassifikation
-- (SPEC-007 HEALTH_STATES) bleibt Sache des lesenden Systems — diese View
-- trifft keine Schwellenwert-Entscheidung.
--
-- Abgedeckt (Minimum, nicht vollständig — Slice-Plan slice-011 §1
-- Ausgrenzung): cdc_transactions_total, cdc_changes_processed,
-- cdc_oldest_change_age_seconds, cdc_consumer_position{consumer},
-- cdc_consumer_lag{consumer}. Nicht abgedeckt: cdc_changes_pending,
-- cdc_capture_lag, cdc_errors_total, cdc_wal_retention_bytes,
-- cdc_storage_bytes (SPEC-009-Zeilen) — sie brauchen entweder
-- persistierten Zustand, den dieses Schema noch nicht trägt (Fehler-Log,
-- Quell-Commit-Zeitstempel) oder Systemkatalog-Zugriffe außerhalb des
-- cdc-Schemas (pg_stat_replication, Relationsgrößen), die die
-- Least-Privilege-Fläche von cdc_reader unnötig erweitern würden; Folge-Slice.
--
-- Der Health-Endpoint (LH-FA-ADM-002, LH-QA-OPS-002) liegt bewusst nicht
-- in dieser Datei: eine SQL-View liest nur persistierten Zustand und kann
-- den Lauf-Zustand des CDC-Prozesses selbst nicht bezeugen — ein
-- abgestürzter Prozess hinterlässt eine weiterhin erreichbare Datenbank,
-- die View würde „gesund" lesen. Das verlangt einen Treiber, der den
-- laufenden Prozess selbst befragt (neuer Driving-Adapter-Zuschnitt,
-- ADR-0020 stellt HTTP/gRPC explizit zurück); offene Architekturfrage,
-- siehe Slice-Plan slice-011 §7.
CREATE OR REPLACE VIEW cdc.metrics AS
SELECT 'cdc_transactions_total'::text AS metric_name, NULL::text AS label, count(*)::numeric AS value
FROM cdc.transaction
UNION ALL
SELECT 'cdc_changes_processed', NULL, count(*)::numeric
FROM cdc.change
UNION ALL
SELECT 'cdc_oldest_change_age_seconds', NULL, COALESCE(extract(epoch FROM (now() - min(committed_at))), 0)::numeric
FROM cdc.transaction
UNION ALL
SELECT 'cdc_consumer_position', cp.consumer_id, cp.acknowledged_position::numeric
FROM cdc.consumer_position cp
UNION ALL
SELECT 'cdc_consumer_lag', cs.consumer_id, (cs.latest_commit_position - cs.acknowledged_position)::numeric
FROM cdc.consumer_status cs
WHERE cs.acknowledged_position IS NOT NULL;

GRANT SELECT ON cdc.metrics TO cdc_reader;
