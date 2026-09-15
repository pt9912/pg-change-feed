# Beleg: slice-036

Vorgang: `slice-036` — der Slice, der die SQL-Administration geliefert hat
(`cdc.enable_table`/`cdc.disable_table` als SQL-Funktionen und die
Administrations-Goroutine mit Live-Reload der laufenden Assembler-Bindung).

Fund: Der Eintrag führte seinen Ausgang `verkörpert` bereits, ohne dass eine
Belegdatei angelegt war — der abgeleitete Zähler stand damit auf 0, obwohl die
Beobachtung aufgetreten und aufgelöst ist. Nachgetragen bei der Closure von
`slice-073` (Register-Paarung der Closure), rekonstruiert aus dem `state.md`
des Eintrags, der den Vorgang selbst nennt.

Quelle: `docs/plan/planning/observations/BEO-PGC/verwaltung-keine-sql-administration/state.md`.
