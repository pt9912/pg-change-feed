# Beleg: Review zu `slice-080`

Vorgang: das Review von `slice-080` (F-1).

Fund: `harness/sensors/db-adapter-coverage.md` und `tools/harness/db-coverage.sh`
erklärten, `-coverpkg` instrumentiere „in jeder Testbinary den ganzen Gegenstand".
Falsch: im Profil erscheinen nur die **verlinkten** Pakete — ein Lauf über
`postgresstorage` allein erzeugt 610 Statements und **null** Zeilen für
`postgresack`/`receive`. Die Deduplizierung ist dennoch nötig, aber aus anderem
Grund: im Replication-Profil erscheint jede Position 2× (einmal mit
`count > 0`, einmal mit 0), und ohne die Regel „mindestens ein Vorkommen" fiele
`receive` von 112/155 auf 0/155.

Quelle: Review zu `slice-080`, F-1 · `docs/plan/planning/in-progress/slice-080-db-adapter-coverage.md`.
