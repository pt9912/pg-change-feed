-- Berichtete manuelle Nacharbeit am Schema-Rollout (ADR-0043
-- §Re-Evaluierungs-Trigger: zulässige Ausweichform für Rollout-Schritte,
-- die d-migrate nicht aus tools/schema/schema.yaml erzeugt): cdc.metrics trägt das
-- Metriken-Minimum nach LH-FA-SST-004/LH-QA-OPS-003 (SPEC-009) als vierte
-- Lese-View — reine Projektion/Aggregation über bereits persistierte
-- Zeilen, keine Domänenentscheidung (ARC-005 „SQL-Funktionen/Views",
-- ADR-0046 Rollenspaltung). Long-Format (metric_name, label, value): ein
-- Monitoring-System liest die Zeilen roh; die Klassifikation
-- (SPEC-007 HEALTH_STATES) bleibt Sache des lesenden Systems — diese View
-- trifft keine Schwellenwert-Entscheidung.
--
-- Abgedeckt (Minimum, nicht vollständig): cdc_transactions_total,
-- cdc_changes_processed, cdc_changes_pending{consumer},
-- cdc_oldest_change_age_seconds, cdc_consumer_position{consumer},
-- cdc_consumer_lag{consumer}, cdc_capture_lag, cdc_errors_total{class},
-- cdc_storage_bytes — alle sieben in `LH-QA-OPS-003` geforderten
-- Dimensionen (CDC-Lag, verarbeitete/ausstehende Changes,
-- Speicherverbrauch, Alter des ältesten Changes, Consumer-Positionen,
-- Fehler). Nicht abgedeckt: cdc_wal_retention_bytes (SPEC-009-Zeile,
-- keine Lastenheft-Dimension) — sie bräuchte Systemkatalog-Zugriffe
-- außerhalb des cdc-Schemas (pg_stat_replication), die die
-- Least-Privilege-Fläche von cdc_reader unnötig erweitern würden.
--
-- `cdc_changes_pending` (`SPEC-009`, `LH-QA-OPS-003`): je Consumer die
-- reale Anzahl noch nicht bestätigter Change-Zeilen seiner Quelle
-- (`commit_position` der Transaktion über der bestätigten Position) —
-- ein Zeilen-Zähler, anders als `cdc_consumer_lag`, das den
-- LSN-Byte-Abstand der Positionen trägt (`ADR-0005`): dieselbe
-- Quellzahl-Differenz kann bei wenigen großen oder vielen kleinen
-- Transaktionen sehr unterschiedliche Change-Zahlen bedeuten.
--
-- `cdc_errors_total` (`SPEC-009`, `LH-QA-OPS-003`, Klassen aus `ADR-0023`):
-- je Fehlerklasse die Anzahl Quellen, deren `cdc.process_heartbeat`
-- aktuell diese Klasse trägt — der current-state-Zähler, den das
-- vorhandene Heartbeat-Feld tatsächlich abbildet (kein historisches
-- Fehler-Log, das dieses Schema nicht führt); eine Quelle ohne aktuellen
-- Fehlerzustand trägt keine Zeile.
--
-- `cdc_capture_lag` (`SPEC-009`, `LH-FA-ADM-004`; Latenzschwellen
-- p95/Warn/Fehler `SPEC-013`): `committed_at` trägt den Quell-
-- Commit-Zeitpunkt aus dem WAL (`InsertTransaction`) — der Wert misst
-- den Abstand zwischen der letzten Quelländerung und der
-- CDC-Verfügbarkeit über `now() - max(committed_at)`.
--
-- `cdc_storage_bytes` (`SPEC-009`, `LH-FA-RET-006`): die physische
-- Speichergröße von `cdc.change` über `pg_relation_size` — die mit dem
-- Erfassungsvolumen wachsende Tabelle dieses Schemas (Row Images als
-- jsonb, `LH-FA-CAP-008`); die übrigen vier Tabellen tragen
-- Referenzdaten (Quelle, Tabellen-/Schema-Katalog, Transaktions-Kopf) und
-- bleiben dagegen klein. `pg_relation_size()` ist eine reguläre, für
-- `PUBLIC` ausführbare Systemfunktion (kein `SECURITY DEFINER`, kein
-- direkter Grant an `cdc_reader` nötig) und liefert reine Metadaten,
-- keinen Zeileninhalt; die View hardcodet den Tabellennamen in ihrer
-- `SELECT`-Klausel, `cdc_reader` kann ihn nicht selbst parametrisieren
-- (`ADR-0046` Kategorie C: Projektion ohne Domänenentscheidung).
--
-- Der Health-Endpoint (LH-FA-ADM-002, LH-QA-OPS-002) liegt bewusst nicht
-- in dieser Datei: eine reine Lese-View auf bereits persistierten Zustand
-- konnte den Lauf-Zustand des CDC-Prozesses selbst nicht bezeugen — ein
-- abgestürzter Prozess hinterlässt eine weiterhin erreichbare Datenbank,
-- die View würde „gesund" lesen. Aufgelöst über den Heartbeat-Mechanismus
-- (kein neuer Driving-Adapter-Zuschnitt, ADR-0020 bleibt unberührt): der
-- Capture-Prozess schreibt sein Lebenszeichen periodisch fort
-- (cdc.process_heartbeat), cdc.heartbeat (tools/schema/nacharbeit-heartbeat.sql)
-- projiziert dessen Alter.
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
SELECT 'cdc_capture_lag', NULL, COALESCE(extract(epoch FROM (now() - max(committed_at))), 0)::numeric
FROM cdc.transaction
UNION ALL
SELECT 'cdc_consumer_position', cp.consumer_id, cp.acknowledged_position::numeric
FROM cdc.consumer_position cp
UNION ALL
SELECT 'cdc_consumer_lag', cs.consumer_id, (cs.latest_commit_position - cs.acknowledged_position)::numeric
FROM cdc.consumer_status cs
WHERE cs.acknowledged_position IS NOT NULL
UNION ALL
SELECT 'cdc_changes_pending', cp.consumer_id,
  (SELECT count(*)
   FROM cdc.change c
   JOIN cdc.transaction t ON t.transaction_id = c.transaction_id
   WHERE t.source_id = cp.source_id AND t.commit_position > cp.acknowledged_position)::numeric
FROM cdc.consumer_position cp
UNION ALL
SELECT 'cdc_errors_total', ph.error_class, count(*)::numeric
FROM cdc.process_heartbeat ph
WHERE ph.error_class IS NOT NULL
GROUP BY ph.error_class
UNION ALL
SELECT 'cdc_storage_bytes', NULL, pg_relation_size('cdc.change')::numeric;

GRANT SELECT ON cdc.metrics TO cdc_reader;
