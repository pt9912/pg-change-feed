**Vorgang:** slice-beispiele-compose-bootstrap
**Fund:** Ein zweiter `make schema-rollout`-Lauf gegen eine bereits
migrierte Demo-Datenbank (`examples/compose.yaml`, eigenständig von der
Wurzel-`compose.yaml`) blockierte real mit Exit 8
(`DESTRUCTIVE_OPERATION_REQUIRES_CONFIRMATION`) — diesmal auf sechs
Objekten: zusätzlich zu `cdc.heartbeat`/`cdc.metrics` (Views) und
`cdc.disable_table`/`cdc.enable_table` (Funktionen, `ADR-0050`) auch
`cdc.exclude_column`/`cdc.include_column` (Funktionen,
`tools/schema/nacharbeit-administration.sql`). Real reproduziert mit
`--execute` (Report `tools/schema/plan.yaml`, `"blockers":
[{"reason":"DESTRUCTIVE_OPERATION_REQUIRES_CONFIRMATION", …}]`).
Aufgelöst, ohne die strukturelle Ursache zu beheben: `examples/bootstrap.sh`
prüft vor jedem Rollout-Aufruf `to_regclass('cdc.source_table')` und
überspringt `make schema-rollout` vollständig, wenn das Ziel bereits
migriert ist — ein Existenz-Check als lokale Wache dieses einen Aufrufers,
kein Fix des zugrunde liegenden Fremdobjekte-Problems. Drittes,
unabhängiges Auftreten (nach slice-016, slice-063) — der Objektumfang
wächst mit jedem neuen `nacharbeit-*.sql`-Skript weiter, die Ursache bleibt
unverändert offen.
