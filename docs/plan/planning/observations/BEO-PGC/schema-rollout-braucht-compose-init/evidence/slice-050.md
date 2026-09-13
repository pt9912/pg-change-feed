**Vorgang:** slice-050
**Fund:** `tools/bench-lib.sh` baute die erste eigenständige, von
`compose.yaml` unabhängige PostgreSQL-Umgebung für `make schema-rollout`
und scheiterte zunächst real mit Exit 5 (`POST_EXECUTE_DRIFT`, `relation
"cdc.source_table" does not exist`), weil der
`tools/schema/compose-init`-Mount fehlte. Behoben durch denselben Mount
wie in `compose.yaml` (`docker-entrypoint-initdb.d`).
