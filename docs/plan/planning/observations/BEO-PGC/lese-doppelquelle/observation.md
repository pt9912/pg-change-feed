# BEO-PGC/lese-doppelquelle

**Sub-Area:** Store-Adapter/Lesepfad (SQL-Views vs. Go-Use-Cases;
Sub-Area-Kürzel `PGC` aus der Modus-Deklaration)

Die Beobachtung: Die SQL-Views im `cdc`-Schema (`active_tables`,
`consumer_status`, `changes`) und die Go-Use-Case-Lesepfade
(`ReadChangesUseCase`/`ChangeStorePort`) implementieren dieselbe
Lese-Semantik zweimal — einmal als SQL-Projektion, einmal als
Go-Code. Beide sind laut [`ADR-0046`](../../../../../plan/adr/README.md)
zulässig (reine Lese-Views ohne Entscheidungslogik sind kein
Businesslogik-Duplikat), aber die **Feld-/Filter-Semantik** kann
zwischen den zwei Implementierungen driften, wenn eine Seite geändert
wird und die andere nicht mitgezogen wird — es gibt keinen Sensor, der
beide gegeneinander hält.

Deklaration: `tools/schema/nacharbeit-views.sql`,
`internal/adapters/driven/postgresstorage/queries/queries.go`
(`SelectChanges`).
