Zustand: **verkörpert** — die real fehlenden Fähigkeiten sind gebaut:
`cdc.enable_table`/`cdc.disable_table` (SQL-Funktionen, `slice-036`), die
Administrations-Goroutine mit Live-Reload der laufenden `Assembler`-Bindung
(`slice-037`, `ADR-0050`), und der `diagnose`-CLI-Befehl, der
`LH-FA-ADM-002`…`005` real abdeckt (`slice-038`). Verkörpert in
`internal/bootstrap/wiring.go` (`runAdministration`,
`processAdministrationRequests`, `Diagnose`) und
`internal/adapters/driving/replication/mapper/mapper.go`
(`AddBinding`/`RemoveBinding`) · seit welle-12. Zähler (abgeleitet): 0× —
kein abgeschlossener Vorgang trug je einen Beleg (die Beobachtung wurde nie
über die 3×-Schwelle gezählt, sondern direkt über eine dedizierte
Feature-Welle aufgelöst, analog zu `BEO-PGC/schema-evolution-nicht-dynamisch`s
direkter Auflösung unter der Schwelle).
