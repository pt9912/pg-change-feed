-- Berichtete manuelle Nacharbeit am Schema-Rollout (ADR-0043
-- §Re-Evaluierungs-Trigger: zulässige Ausweichform für Rollout-Schritte,
-- die d-migrate nicht aus tools/schema/schema.yaml erzeugt): cdc.enable_table/
-- cdc.disable_table sind die schreibenden SQL-Funktionen der Antrags-Queue
-- (ADR-0050, LH-FA-ADM-001). d-migrate 1.3.1 generiert ihre DDL korrekt über
-- den `functions:`-Knoten (`schema generate`), aber `schema migrate --execute`
-- bricht für jede dort deklarierte Funktion mit POST_EXECUTE_DRIFT (Exit 5)
-- ab — real reproduziert auch mit einer trivialen No-Arg-Funktion gegen eine
-- frische PostgreSQL-18-Instanz, nicht auf diese beiden Funktionen beschränkt.
-- Herkunft: `docs/plan/planning/observations/BEO-PGC/d-migrate-nacharbeit`.
-- Diese Datei läuft im schema-rollout-Lauf nach nacharbeit-roles.sql
-- (Makefile) — sie grantet an cdc_admin, die Rolle entsteht erst dort.
--
-- ADR-0050s zentrale Entscheidung: eine schreibende SQL-Funktion schreibt
-- ausschließlich einen Antrags-Datensatz (Status `pending`) und sendet
-- `pg_notify` — sie berührt cdc.source_table, cdc.schema_version und die
-- Publication nicht direkt (das bleibt exklusiv TableActivationAdapter über
-- EnableTableUseCase/DisableTableUseCase vorbehalten, ausgeführt von der
-- Administrations-Goroutine eines Folge-Slice). Kanal-Name `cdc_administration`
-- ist Implementer-Entscheidung, real mit `LISTEN`/pgx WaitForNotification
-- getestet (administrationrequest_test.go).
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
        (administration_request_id, source_id, schema_name, table_name, request_kind, status)
    VALUES (v_id, p_source_id, p_schema_name, p_table_name, 'enable', 'pending');
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
        (administration_request_id, source_id, schema_name, table_name, request_kind, status)
    VALUES (v_id, p_source_id, p_schema_name, p_table_name, 'disable', 'pending');
    PERFORM pg_notify('cdc_administration', v_id);
    RETURN v_id;
END;
$$;

REVOKE EXECUTE ON FUNCTION cdc.enable_table(text, text, text), cdc.disable_table(text, text, text) FROM PUBLIC;
GRANT EXECUTE ON FUNCTION cdc.enable_table(text, text, text), cdc.disable_table(text, text, text) TO cdc_admin;
