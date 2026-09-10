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
-- cdc_consumer_lag{consumer}. Seit slice-013 zusätzlich
-- cdc_capture_lag_approx — bewusst **nicht** unter dem SPEC-009-
-- kanonischen Namen `cdc_capture_lag`, siehe Grenze unten. Nicht
-- abgedeckt: der reale `cdc_capture_lag` selbst, cdc_changes_pending,
-- cdc_errors_total, cdc_wal_retention_bytes, cdc_storage_bytes
-- (SPEC-009-Zeilen) — sie brauchen entweder persistierten Zustand, den
-- dieses Schema noch nicht trägt (Fehler-Log, Quell-Commit-Zeitstempel)
-- oder Systemkatalog-Zugriffe außerhalb des cdc-Schemas
-- (pg_stat_replication, Relationsgrößen), die die Least-Privilege-Fläche
-- von cdc_reader unnötig erweitern würden; Folge-Slice
-- (`BEO-PGC/cdc-capture-lag-real`).
--
-- Grenze (cdc_capture_lag_approx, slice-013 §4 Rückführung): `committed_at`
-- trägt die Persistenz-Zeit der Instanz (DB-DEFAULT now(), s.
-- tools/schema/schema.yaml), nicht den tatsächlichen Quell-Commit-
-- Zeitpunkt aus dem WAL — pgoutput trägt ihn zwar bereits
-- (BEGIN/COMMIT-Nachrichten, github.com/jackc/pglogrepl), aber ihn bis in
-- diese Persistenz zu spiegeln berührt zusätzlich die
-- Replication-Adapter-Schicht (mehr als zwei Schichten, Folge-Slice). Der
-- Wert unten misst deshalb nur, wie lange die letzte persistierte
-- Transaktion zurückliegt (Pipeline-Frische) — kein LH-FA-ADM-004-Abstand
-- zur Quelländerung. Der `_approx`-Suffix im metric_name ist die
-- Interface-Markierung dieser Grenze (review-slice-013 F-2): ein
-- Monitoring-System, das die Zeilen roh liest (Datei-Kopfkommentar
-- oben), unterscheidet die Näherung sonst nicht vom SPEC-009-kanonischen
-- `cdc_capture_lag` — der Name selbst trägt die Zusage, nicht nur der
-- SQL-Kommentar. `SPEC-013`s Latenzschwellen (p95/Warn/Fehler) sind an
-- den kanonischen Namen `cdc_capture_lag` gebunden und bewusst NICHT an
-- `cdc_capture_lag_approx` anwendbar, bis der reale Abstand existiert.
--
-- Der Health-Endpoint (LH-FA-ADM-002, LH-QA-OPS-002) liegt bewusst nicht
-- in dieser Datei: eine reine Lese-View auf bereits persistierten Zustand
-- konnte den Lauf-Zustand des CDC-Prozesses selbst nicht bezeugen — ein
-- abgestürzter Prozess hinterlässt eine weiterhin erreichbare Datenbank,
-- die View würde „gesund" lesen. Aufgelöst über den Heartbeat-Mechanismus
-- (kein neuer Driving-Adapter-Zuschnitt, ADR-0020 bleibt unberührt,
-- Architect-Verdikt docs/plan/adr/architect-review-slice-011.md): der
-- Capture-Prozess schreibt sein Lebenszeichen periodisch fort
-- (cdc.process_heartbeat), cdc.heartbeat (tools/schema/nacharbeit-heartbeat.sql,
-- slice-012) projiziert dessen Alter.
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
SELECT 'cdc_capture_lag_approx', NULL, COALESCE(extract(epoch FROM (now() - max(committed_at))), 0)::numeric
FROM cdc.transaction
UNION ALL
SELECT 'cdc_consumer_position', cp.consumer_id, cp.acknowledged_position::numeric
FROM cdc.consumer_position cp
UNION ALL
SELECT 'cdc_consumer_lag', cs.consumer_id, (cs.latest_commit_position - cs.acknowledged_position)::numeric
FROM cdc.consumer_status cs
WHERE cs.acknowledged_position IS NOT NULL;

GRANT SELECT ON cdc.metrics TO cdc_reader;
