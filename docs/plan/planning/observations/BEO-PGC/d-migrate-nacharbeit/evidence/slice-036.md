# Beleg: slice-036 (Antragsqueue-SQL-Funktionen, ADR-0050)

Der `functions:`-Knoten des neutralen Schemamodells wurde real gegen den
gepinnten d-migrate 1.3.1 getestet (Digest
`sha256:862dfb04c34dd17278b1bab46961363c12eeb8d464cf1776565d6285603d2c89`):

- `schema generate --source … --target postgresql` erzeugt die
  `CREATE OR REPLACE FUNCTION …`-DDL korrekt — inklusive `parameters`,
  `returns`, `security: definer` (`SECURITY DEFINER`) und `search_path`
  (`SET search_path = …`).
- `schema migrate --source … --target db:… --execute` bricht dagegen für
  **jede** im `functions:`-Knoten deklarierte Funktion mit
  `POST_EXECUTE_DRIFT` (Exit 5) ab — real isoliert: eine triviale
  No-Arg-Funktion ohne Parameter/Security/`search_path` reproduziert den
  Fehler ebenso wie `cdc.enable_table`/`cdc.disable_table` selbst, gegen
  eine frische PostgreSQL-18-Instanz. Eine Tabelle im selben Lauf (ohne
  Funktion) rollt dagegen Exit 0.

Ausweichform: `tools/schema/nacharbeit-administration.sql` trägt die beiden
Funktionen (dieselbe Klasse wie `nacharbeit-roles.sql`); die Antrags-Tabelle
`cdc.administration_request` läuft weiterhin über `tools/schema/schema.yaml`
— nur die Objektklasse „Funktion" ist betroffen, nicht Tabellen. Real gegen
PostgreSQL getestet: Antrag mit Status `pending` entsteht,
`pg_notify('cdc_administration', …)` wird über `LISTEN` empfangen, ein
Login ohne `cdc_admin`-Mitgliedschaft scheitert mit „permission denied for
function" (Least-Privilege, `REVOKE … FROM PUBLIC` vor dem gezielten
`GRANT … TO cdc_admin`).
