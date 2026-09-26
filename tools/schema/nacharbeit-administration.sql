-- Berichtete manuelle Nacharbeit am Schema-Rollout (ADR-0043
-- §Re-Evaluierungs-Trigger: zulässige Ausweichform für Rollout-Schritte,
-- die d-migrate nicht aus tools/schema/schema.yaml erzeugt). Zwei
-- Objektklassen der Antrags-Queue liegen hier:
--
-- (1) cdc.enable_table/cdc.disable_table/cdc.exclude_column/
-- cdc.include_column/cdc.backfill_table/cdc.set_transformation/
-- cdc.remove_transformation sind ihre schreibenden SQL-Funktionen
-- (ADR-0050, LH-FA-ADM-001, LH-FA-CFG-005, LH-FA-CAP-009, LH-FA-CFG-007,
-- ADR-0112). d-migrate 1.3.1 generiert ihre
-- DDL korrekt über
-- den `functions:`-Knoten (`schema generate`), aber `schema migrate --execute`
-- bricht für jede dort deklarierte Funktion mit POST_EXECUTE_DRIFT (Exit 5)
-- ab — real reproduziert auch mit einer trivialen No-Arg-Funktion gegen eine
-- frische PostgreSQL-18-Instanz, nicht auf diese Funktionen beschränkt.
--
-- (2) Die geschlossene request_kind-Menge (`enable`, `disable`,
-- `exclude_column`, `include_column`, `backfill`, `set_transformation`,
-- `remove_transformation`) trägt der CHECK
-- chk_administration_request_kind. Er steht nicht als deklarativer Check in
-- tools/schema/schema.yaml, weil d-migrate 1.3.1 eine CHECK-Änderung an
-- einer bestehenden Tabelle nicht konvergiert: real gemessen — den
-- Schema-Stand ohne die beiden Spalten-Antragsarten ausgerollt (Klausel
-- `enable`/`disable`), dann denselben Rollout mit der erweiterten Klausel
-- — endet mit POST_EXECUTE_DRIFT (Exit 5), und die bestehende Klausel
-- entfällt dabei; dasselbe Ergebnis mit umbenanntem Constraint, also
-- unabhängig vom Namen, und ohne dass die neue Klausel entsteht. Die
-- Anlage der Spalte `column_name` an derselben bestehenden Tabelle
-- konvergiert dagegen (real gemessen, Exit 0) und bleibt deshalb im
-- deklarativen Modell; ebenso die Spalten `rule_name` und `rule_spec`.
-- Die beiden Zeilen unten setzen die Menge
-- idempotent: DROP IF EXISTS vor ADD, damit ein Erstanlage-Lauf und jeder
-- Folgelauf denselben Endzustand tragen.
--
-- Herkunft beider Klassen:
-- `docs/plan/planning/observations/BEO-PGC/d-migrate-nacharbeit`.
-- Diese Datei läuft im schema-rollout-Lauf nach nacharbeit-roles.sql
-- (Makefile) — sie grantet an cdc_admin, die Rolle entsteht erst dort.
--
-- ADR-0050s zentrale Entscheidung: eine schreibende SQL-Funktion schreibt
-- ausschließlich einen Antrags-Datensatz (Status `pending`) und sendet
-- `pg_notify` — sie berührt cdc.source_table, cdc.schema_version und die
-- Publication nicht direkt (das bleibt exklusiv TableActivationAdapter über
-- EnableTableUseCase/DisableTableUseCase vorbehalten, ausgeführt von der
-- Administrations-Goroutine). Dieselbe Disziplin trägt
-- der Spaltenausschluss (`LH-FA-CFG-005`, ADR-0059): cdc.exclude_column/
-- cdc.include_column schreiben ausschließlich den Antrags-Datensatz samt
-- Spaltennamen und prüfen nichts an der Quelle — die Spaltenexistenz trägt
-- der Use Case über ColumnExclusionPort. cdc.backfill_table (`LH-FA-CAP-009`,
-- ADR-0111) schreibt ebenso ausschließlich den Antrags-Datensatz der Art
-- `backfill` (ohne Spalte); Vorbedingungen und Annahme trägt der Use Case,
-- und `applied` heißt dort „angenommen“ (SPEC-019).
-- cdc.set_transformation/cdc.remove_transformation (`LH-FA-CFG-007`, ADR-0112)
-- schreiben ebenso ausschließlich den Antrags-Datensatz samt Regelname, die
-- erste zusätzlich die Regelform, und prüfen nichts — weder die Regelform
-- noch die Spalte; beides trägt der Capture-Prozess (ADR-0046). Die Regelform
-- kommt als `json`-Parameter und geht als `jsonb` in die Spalte `rule_spec`
-- (ADR-0125); der Parametertyp ist `json`, weil d-migrate 1.3.1 jede Funktion mit `json`-
-- oder `jsonb`-Parameter als `in:json` meldet und ihren Abbau im zweiten
-- Rollout als `DROP FUNCTION … (…, json)` rendert — für eine Funktion mit
-- `jsonb`-Parameter existiert diese Signatur nicht, der Abbau scheitert dort
-- mit Exit 5 (real gemessen). Ein Literal und ein `::json`-Wert werden
-- angenommen, ein `::jsonb`-Wert nicht (real gemessen: „function … does not
-- exist“).
-- Jede Funktion schreibt `requested_at` je Aufruf mit `clock_timestamp()`
-- (der Spalten-Default `current_timestamp`, der Transaktionsbeginn, greift nur
-- für einen INSERT ohne Wert): Aufrufe derselben Transaktion tragen
-- verschiedene Zeitstempel in der Reihenfolge des Aufrufs, und die Ordnung
-- der Antrags-Queue ist die Ordnung des Aufrufs (ADR-0127).
-- Der Kanal-Name `cdc_administration` ist Implementer-Entscheidung, real mit
-- `LISTEN`/pgx WaitForNotification getestet (administrationrequest_test.go).
--
-- SECURITY DEFINER mit gepinntem search_path (dieselbe Absicherung, die
-- eine SECURITY-DEFINER-Routine braucht: ohne ihn entscheidet der Suchpfad
-- des Aufrufers, welche Tabelle der Rumpf trifft) — Definer-Semantik, wie sie
-- die drei Lese-Views bereits ohne `security_invoker` tragen (nacharbeit-roles.sql).
-- Least-Privilege (LH-QA-SEC-001…003, ADR-0047): PostgreSQL grantet EXECUTE
-- auf eine neue Funktion standardmäßig an PUBLIC — die REVOKE-Zeile schließt
-- das, bevor cdc_admin exklusiv das Recht bekommt (real geprüft: ein Login
-- ohne cdc_admin-Mitgliedschaft scheitert mit „permission denied for
-- function", ein Login mit cdc_admin-Mitgliedschaft gelingt).
ALTER TABLE cdc.administration_request DROP CONSTRAINT IF EXISTS chk_administration_request_kind;
ALTER TABLE cdc.administration_request ADD CONSTRAINT chk_administration_request_kind CHECK (request_kind IN ('enable', 'disable', 'exclude_column', 'include_column', 'backfill', 'set_transformation', 'remove_transformation'));

CREATE OR REPLACE FUNCTION cdc.enable_table(p_source_id text, p_schema_name text, p_table_name text)
RETURNS text
LANGUAGE plpgsql
VOLATILE
SECURITY DEFINER
SET search_path = cdc, pg_temp
AS $$
DECLARE
    v_id text;
BEGIN
    v_id := gen_random_uuid()::text;
    INSERT INTO cdc.administration_request
        (administration_request_id, requested_at, source_id, schema_name, table_name, request_kind, status)
    VALUES (v_id, clock_timestamp(), p_source_id, p_schema_name, p_table_name, 'enable', 'pending');
    PERFORM pg_notify('cdc_administration', v_id);
    RETURN v_id;
END;
$$;

CREATE OR REPLACE FUNCTION cdc.disable_table(p_source_id text, p_schema_name text, p_table_name text)
RETURNS text
LANGUAGE plpgsql
VOLATILE
SECURITY DEFINER
SET search_path = cdc, pg_temp
AS $$
DECLARE
    v_id text;
BEGIN
    v_id := gen_random_uuid()::text;
    INSERT INTO cdc.administration_request
        (administration_request_id, requested_at, source_id, schema_name, table_name, request_kind, status)
    VALUES (v_id, clock_timestamp(), p_source_id, p_schema_name, p_table_name, 'disable', 'pending');
    PERFORM pg_notify('cdc_administration', v_id);
    RETURN v_id;
END;
$$;

CREATE OR REPLACE FUNCTION cdc.exclude_column(p_source_id text, p_schema_name text, p_table_name text, p_column_name text)
RETURNS text
LANGUAGE plpgsql
VOLATILE
SECURITY DEFINER
SET search_path = cdc, pg_temp
AS $$
DECLARE
    v_id text;
BEGIN
    v_id := gen_random_uuid()::text;
    INSERT INTO cdc.administration_request
        (administration_request_id, requested_at, source_id, schema_name, table_name, column_name, request_kind, status)
    VALUES (v_id, clock_timestamp(), p_source_id, p_schema_name, p_table_name, p_column_name, 'exclude_column', 'pending');
    PERFORM pg_notify('cdc_administration', v_id);
    RETURN v_id;
END;
$$;

CREATE OR REPLACE FUNCTION cdc.include_column(p_source_id text, p_schema_name text, p_table_name text, p_column_name text)
RETURNS text
LANGUAGE plpgsql
VOLATILE
SECURITY DEFINER
SET search_path = cdc, pg_temp
AS $$
DECLARE
    v_id text;
BEGIN
    v_id := gen_random_uuid()::text;
    INSERT INTO cdc.administration_request
        (administration_request_id, requested_at, source_id, schema_name, table_name, column_name, request_kind, status)
    VALUES (v_id, clock_timestamp(), p_source_id, p_schema_name, p_table_name, p_column_name, 'include_column', 'pending');
    PERFORM pg_notify('cdc_administration', v_id);
    RETURN v_id;
END;
$$;

CREATE OR REPLACE FUNCTION cdc.backfill_table(p_source_id text, p_schema_name text, p_table_name text)
RETURNS text
LANGUAGE plpgsql
VOLATILE
SECURITY DEFINER
SET search_path = cdc, pg_temp
AS $$
DECLARE
    v_id text;
BEGIN
    v_id := gen_random_uuid()::text;
    INSERT INTO cdc.administration_request
        (administration_request_id, requested_at, source_id, schema_name, table_name, request_kind, status)
    VALUES (v_id, clock_timestamp(), p_source_id, p_schema_name, p_table_name, 'backfill', 'pending');
    PERFORM pg_notify('cdc_administration', v_id);
    RETURN v_id;
END;
$$;

CREATE OR REPLACE FUNCTION cdc.set_transformation(p_source_id text, p_schema_name text, p_table_name text, p_rule_name text, p_rule_spec json)
RETURNS text
LANGUAGE plpgsql
VOLATILE
SECURITY DEFINER
SET search_path = cdc, pg_temp
AS $$
DECLARE
    v_id text;
BEGIN
    v_id := gen_random_uuid()::text;
    INSERT INTO cdc.administration_request
        (administration_request_id, requested_at, source_id, schema_name, table_name, rule_name, rule_spec, request_kind, status)
    VALUES (v_id, clock_timestamp(), p_source_id, p_schema_name, p_table_name, p_rule_name, p_rule_spec::jsonb, 'set_transformation', 'pending');
    PERFORM pg_notify('cdc_administration', v_id);
    RETURN v_id;
END;
$$;

CREATE OR REPLACE FUNCTION cdc.remove_transformation(p_source_id text, p_schema_name text, p_table_name text, p_rule_name text)
RETURNS text
LANGUAGE plpgsql
VOLATILE
SECURITY DEFINER
SET search_path = cdc, pg_temp
AS $$
DECLARE
    v_id text;
BEGIN
    v_id := gen_random_uuid()::text;
    INSERT INTO cdc.administration_request
        (administration_request_id, requested_at, source_id, schema_name, table_name, rule_name, request_kind, status)
    VALUES (v_id, clock_timestamp(), p_source_id, p_schema_name, p_table_name, p_rule_name, 'remove_transformation', 'pending');
    PERFORM pg_notify('cdc_administration', v_id);
    RETURN v_id;
END;
$$;

REVOKE EXECUTE ON FUNCTION cdc.enable_table(text, text, text), cdc.disable_table(text, text, text), cdc.exclude_column(text, text, text, text), cdc.include_column(text, text, text, text), cdc.backfill_table(text, text, text), cdc.set_transformation(text, text, text, text, json), cdc.remove_transformation(text, text, text, text) FROM PUBLIC;
GRANT EXECUTE ON FUNCTION cdc.enable_table(text, text, text), cdc.disable_table(text, text, text), cdc.exclude_column(text, text, text, text), cdc.include_column(text, text, text, text), cdc.backfill_table(text, text, text), cdc.set_transformation(text, text, text, text, json), cdc.remove_transformation(text, text, text, text) TO cdc_admin;
