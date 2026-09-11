-- Berichtete manuelle Nacharbeit am Schema-Rollout (ADR-0043
-- §Re-Evaluierungs-Trigger: zulässige Ausweichform für Rollout-Schritte,
-- die d-migrate nicht aus tools/schema/schema.yaml erzeugt): die drei
-- Least-Privilege-Rollen nach LH-QA-SEC-001…003 sind kein Tabellen- oder
-- View-Objekt und liegen deshalb außerhalb von tools/schema/schema.yaml
-- (dessen tables:-Knoten trägt nur d-migrate-überführbare Objekte,
-- ADR-0043). CREATE ROLE kennt kein IF NOT EXISTS — der DO-Block macht
-- den Schritt wiederholbar wie die übrigen Nacharbeit-Dateien.
--
-- Rollenschnitt (Slice-Plan slice-011, LH-QA-SEC-001 Least-Privilege,
-- LH-QA-SEC-002 getrennte Berechtigbarkeit): cdc_capture trägt den
-- Erfassungspfad — REPLICATION-Attribut für den Replication-Stream
-- (ADR-0008), INSERT auf transaction/change (queries.go InsertTransaction/
-- InsertChange), SELECT auf source_table/schema_version für die
-- Bindungs-Auflösung. cdc_admin trägt die Registrierungs-/
-- Verwaltungspfade — DML auf source_table/schema_version (LH-FA-CFG-001.a)
-- und consumer/consumer_position (LH-FA-CON-001…006), CREATE auf der
-- Datenbank für `CREATE PUBLICATION` selbst
-- (tableactivation.go CREATE/ALTER PUBLICATION) — für das Hinzufügen von
-- Tabellen zur Publication reicht das allein nicht, siehe die Grenze
-- weiter unten. cdc_reader trägt
-- ausschließlich die Lese-Views (LH-FA-SST-002) — kein Grant auf eine
-- Basistabelle: PostgreSQL-Views laufen mit den Rechten des
-- View-Eigentümers (Definer-Semantik ohne `security_invoker`), das
-- SELECT-Grant auf die View allein trägt den Lesezugriff (LH-QA-SEC-003).
-- cdc.metrics (tools/schema/nacharbeit-observability.sql) grantet sich
-- selbst an cdc_reader, weil die View dort erst entsteht — diese Datei
-- läuft im schema-rollout-Lauf davor (Makefile).
--
-- Grenze: alle drei Rollen bleiben NOLOGIN — sie bündeln Privilegien
-- (Gruppenrollen), keine Anmelde-Identität mit Passwort. Die Verdrahtung
-- (internal/bootstrap/wiring.go) trägt in diesem Slice-Stand weiterhin
-- eine einzige Instanz-DSN für Store-, Aktivierungs- und
-- Stream-Verbindung; das Binden einzelner Anmelde-Rollen je Adapter ist
-- offenes Risiko dieses Slice (Slice-Plan §6) und Folge-Arbeit.
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'cdc_capture') THEN
        CREATE ROLE cdc_capture NOLOGIN REPLICATION;
    END IF;
    IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'cdc_admin') THEN
        CREATE ROLE cdc_admin NOLOGIN;
    END IF;
    IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'cdc_reader') THEN
        CREATE ROLE cdc_reader NOLOGIN;
    END IF;
END
$$;

GRANT USAGE ON SCHEMA cdc TO cdc_capture, cdc_admin, cdc_reader;

-- cdc_capture: Erfassungspfad (ADR-0008, SPEC-002).
GRANT SELECT ON cdc.source_table, cdc.schema_version TO cdc_capture;
GRANT INSERT ON cdc.transaction, cdc.change TO cdc_capture;

-- cdc_admin: Registrierungs- und Verwaltungspfad (LH-FA-CFG-001.a,
-- LH-FA-CON-001…006). Kein Grant auf transaction/change — der
-- Erfassungspfad bleibt cdc_capture vorbehalten (LH-QA-SEC-002).
GRANT SELECT, INSERT, UPDATE, DELETE ON cdc.source_table, cdc.schema_version, cdc.consumer, cdc.consumer_position TO cdc_admin;
GRANT SELECT ON cdc.transaction, cdc.change TO cdc_admin;

-- CREATE PUBLICATION/ALTER PUBLICATION (tableactivation.go) verlangt das
-- CREATE-Privileg auf der Zieldatenbank; der Datenbankname wechselt
-- zwischen den Testumgebungen (cdc, cdc_test) — current_database() trägt
-- ihn dynamisch, GRANT nimmt keinen Funktionsaufruf als Objektnamen.
DO $$
BEGIN
    EXECUTE format('GRANT CREATE ON DATABASE %I TO cdc_admin', current_database());
END
$$;

-- Grenze (real geprüft, roles_test.go
-- TestCdcAdminPublicationRequiresTableOwnership/
-- TestCdcAdminAlterPublicationAddTableRequiresTableOwnership,
-- Review-Finding F-3 review-slice-011.md): CREATE ON DATABASE trägt nur
-- `CREATE PUBLICATION` selbst — für `… FOR TABLE`/`ALTER PUBLICATION …
-- ADD TABLE` verlangt PostgreSQL zusätzlich Eigentümerrechte an JEDER
-- hinzugefügten Tabelle (SQLSTATE 42501 „must be owner of table …" ohne
-- sie, empirisch bestätigt). Diese Datei kann das nicht generisch
-- grantet: die zu aktivierenden Quelltabellen liegen außerhalb des
-- cdc-Schemas, gehören dem Aufrufer und existieren beim Schema-Rollout
-- noch nicht. Eine differenziertere Grant-Strategie ist deshalb kein
-- fester GRANT hier, sondern eine Betriebs-Vorbedingung je aktivierter
-- Tabelle — der Aufrufer überträgt die Eigentümerschaft
-- (`ALTER TABLE <quelle> OWNER TO cdc_admin`) oder nimmt cdc_admin in die
-- Eigentümer-Rolle auf (`GRANT <eigentümer-rolle> TO cdc_admin`), bevor
-- die Tabelle aktiviert wird — dieselbe Klasse offener Vorbedingung wie
-- REPLICA IDENTITY (tools/harness/run-integration-tests.sh). Diese
-- Vorbedingung wird erst wirksam, sobald wiring.go die Rolle tatsächlich
-- einsetzt (offenes Risiko, Slice-Plan §6); bis dahin ist sie dokumentiert
-- und getestet, nicht angenommen.

-- cdc_reader: ausschließlich die drei bestehenden Lese-Views.
GRANT SELECT ON cdc.active_tables, cdc.consumer_status, cdc.changes TO cdc_reader;
