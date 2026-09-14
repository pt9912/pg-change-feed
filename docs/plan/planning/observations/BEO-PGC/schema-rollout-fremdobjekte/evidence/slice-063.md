**Vorgang:** slice-063
**Fund:** Ein zweiter `make schema-rollout`-Lauf gegen eine bereits
migrierte Compose-DB blockierte real mit Exit 8
(`DESTRUCTIVE_OPERATION_REQUIRES_CONFIRMATION`) — diesmal auf vier statt
zwei Objekten: zusätzlich zu `cdc.heartbeat`/`cdc.metrics` (Views) auch
`cdc.disable_table`/`cdc.enable_table` (Funktionen, `ADR-0050`). Real
reproduziert mit UND ohne `--execute` (`--dry-run` liefert denselben
blockierten Report). Aufgelöst über `ADR-0064` (Supersedes `ADR-0058`
Entscheidung 3): der Migrationsschritt entfällt ersatzlos, ein
Container-Tausch ersetzt ihn — die strukturelle Ursache selbst
(Fremdobjekte außerhalb des neutralen Modells) bleibt unberührt offen.
Zweites, unabhängiges Auftreten.
