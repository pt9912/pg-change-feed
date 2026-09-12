# Evidence: slice-029

`slice-029` liefert den in der Beobachtung benannten fehlenden Sensor: der
neue Testfall `TestMVPChangesViewMatchesReadChanges`
(`test/integration/integration_test.go`) liest denselben INSERT/UPDATE/
DELETE-Rundlauf sowohl über die externe SQL-Sicht `cdc.changes` (rohe
`pgx`-Query, kein `postgresstorage`-Import) als auch über den internen
Go-Adapter `ReadChanges` und vergleicht beide auf Reihenfolge
(`commit_position`, `sequence`) und Feldinhalt (`operation`,
`old_data`/`new_data`, `schema_version`).

Der Sensor ist real scharf, nicht nur behauptet: Verifikation
(`docs/reviews/verify-slice-029.md`) hat die `changes`-View-Spaltenprojektion
testweise vertauscht und bestätigt, dass **ausschließlich** dieser Testfall
rot wird, während die übrigen vier Testfälle grün bleiben.

Der Testfall läuft mit jedem `make test-integration`-Lauf mit — er ist
damit ein dauerhafter Sensor gegen genau die Drift, die die Beobachtung
benennt, nicht ein einmaliger Beleg.
