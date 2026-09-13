# BEO-PGC/schema-rollout-braucht-compose-init

**Sub-Area:** Schemamigration (d-migrate; Sub-Area-Kürzel `PGC` aus der
Modus-Deklaration)

Die Beobachtung: `make schema-rollout` gegen eine PostgreSQL-Instanz, die
nicht über `compose.yaml` (mit dem `tools/schema/compose-init`-Mount)
gestartet wurde, scheitert real mit `POST_EXECUTE_DRIFT` (Exit 5,
`executionError: relation "cdc.source_table" does not exist`) beim Anlegen
der ersten View. Ursache: Die generierte DDL aus `tools/schema/schema.yaml`
trägt unqualifizierte Tabellennamen (`CREATE TABLE "source_table" (...)`)
und verlässt sich auf den `search_path` der Rollout-Verbindung, um im
Schema `cdc` statt `public` zu landen — dieser `search_path` und das leere
Schema `cdc` selbst entstehen nicht durch den Rollout, sondern durch
`tools/schema/compose-init/01-cdc-schema.sql` (`CREATE SCHEMA IF NOT EXISTS
cdc; ALTER ROLE postgres IN DATABASE cdc SET search_path = cdc;`), das
`compose.yaml` als PostgreSQL-Init-Skript mountet. Die generierten Views
referenzieren `cdc.<table>` explizit — ohne das Init-Skript landen die
Tabellen in `public`, und die Views scheitern sofort. Real aufgetreten
beim Bau einer eigenständigen, von `compose.yaml` unabhängigen
Bench-Umgebung (`tools/bench-lib.sh`, `docker network`/`docker run` statt
Compose), die den Mount zunächst nicht nachbildete.

## Benannt, nicht gezählt

Keine.
