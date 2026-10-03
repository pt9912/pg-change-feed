-- Berichtete manuelle Nacharbeit am Schema-Rollout (ADR-0043
-- §Re-Evaluierungs-Trigger: zulässige Ausweichform für Rollout-Schritte,
-- die d-migrate nicht aus tools/schema/schema.yaml erzeugt): cdc.heartbeat trägt den
-- Health-Endpoint als weitere
-- Lese-View — reine Projektion über die vom Capture-Prozess periodisch
-- fortgeschriebene Lebenszeichen-Zeile (cdc.process_heartbeat), keine
-- Domänenentscheidung ("SQL-Funktionen/Views",
-- Rollenspaltung, Kategorie C). Die Alters-Schwelle
-- (HEALTH_STATES: healthy/degraded/unhealthy) bleibt Sache des lesenden
-- Systems — diese View liefert nur den Zeitstempel und sein Alter in
-- Sekunden, keine Schwellenwert-Entscheidung (dasselbe Muster wie
-- cdc.metrics, nacharbeit-observability.sql). Diese Datei läuft im
-- schema-rollout-Lauf nach nacharbeit-roles.sql (Makefile) — der
-- Rollen-Schritt kann cdc.heartbeat deshalb nicht grants, die View
-- grantet sich selbst an cdc_reader, wie zuvor schon cdc.metrics.
--
-- error_class (LH-FA-ADM-003) projiziert den
-- zuletzt beobachteten Fehlerzustand derselben Zeile — NULL ist
-- Normalbetrieb und von einer der sieben Kategorien
-- unterscheidbar, kein eigenes Fehler-Log.
--
-- error_code projiziert den Meldungscode desselben Fehlerzustands. Die View ist
-- ein bekanntes Fremdobjekt der Wache (DropView in knownForeignObjects,
-- tools/schema/rolloutguard/guard.go): der Rollout löscht sie vor dieser Datei
-- und legt sie hier neu an, ein Vorlauf entfällt. Die Spalte steht hinter
-- age_seconds, die bestehenden Spalten behalten Namen und Reihenfolge. Der Code
-- gilt nur neben einer Klasse: ein Beat eines Servers ohne error_code löscht nur error_class,
-- die Sicht zeigt dann keinen Code ohne Klasse.
CREATE OR REPLACE VIEW cdc.heartbeat AS
SELECT
    source_id,
    heartbeat_at,
    error_class,
    extract(epoch FROM (now() - heartbeat_at))::numeric AS age_seconds,
    CASE WHEN error_class IS NULL THEN NULL ELSE error_code END AS error_code
FROM cdc.process_heartbeat;

GRANT SELECT ON cdc.heartbeat TO cdc_reader;
