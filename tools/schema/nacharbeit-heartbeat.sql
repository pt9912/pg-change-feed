-- Berichtete manuelle Nacharbeit am Schema-Rollout (ADR-0043
-- §Re-Evaluierungs-Trigger: zulässige Ausweichform für Rollout-Schritte,
-- die d-migrate nicht aus tools/schema/schema.yaml erzeugt): cdc.heartbeat trägt den
-- Health-Endpoint (LH-FA-ADM-002, LH-QA-OPS-002; slice-012) als weitere
-- Lese-View — reine Projektion über die vom Capture-Prozess periodisch
-- fortgeschriebene Lebenszeichen-Zeile (cdc.process_heartbeat), keine
-- Domänenentscheidung (ARC-005 „SQL-Funktionen/Views", ADR-0046
-- Rollenspaltung, Kategorie C). Die Alters-Schwelle (SPEC-007
-- HEALTH_STATES: healthy/degraded/unhealthy) bleibt Sache des lesenden
-- Systems — diese View liefert nur den Zeitstempel und sein Alter in
-- Sekunden, keine Schwellenwert-Entscheidung (dasselbe Muster wie
-- cdc.metrics, nacharbeit-observability.sql). Diese Datei läuft im
-- schema-rollout-Lauf nach nacharbeit-roles.sql (Makefile) — der
-- Rollen-Schritt kann cdc.heartbeat deshalb nicht grants, die View
-- grantet sich selbst an cdc_reader, wie zuvor schon cdc.metrics.
--
-- error_class (slice-013, LH-FA-ADM-003, LH-QA-REL-003) projiziert den
-- zuletzt beobachteten Fehlerzustand derselben Zeile — NULL ist
-- Normalbetrieb und von einer der sieben ADR-0023-Kategorien
-- unterscheidbar, kein eigenes Fehler-Log.
CREATE OR REPLACE VIEW cdc.heartbeat AS
SELECT
    source_id,
    heartbeat_at,
    error_class,
    extract(epoch FROM (now() - heartbeat_at))::numeric AS age_seconds
FROM cdc.process_heartbeat;

GRANT SELECT ON cdc.heartbeat TO cdc_reader;
