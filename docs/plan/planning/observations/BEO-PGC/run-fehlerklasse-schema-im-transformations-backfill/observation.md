# BEO-PGC/run-fehlerklasse-schema-im-transformations-backfill

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die Fehlerklassen eines
Backfill-Runs, keine eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: `ADR-0111` Teilfrage 5 nennt für den Run fünf Fehlerklassen (`permission`,
`configuration`, `storage`, `transient`, `replication`). Der Use Case vergibt genau diese;
einen nicht erkannten Fehler fängt `internal` (in der geschlossenen Menge von `SPEC-008` und
`ADR-0023`), die Klasse `schema` vergibt der Run nicht (`classifyError` in
`internal/application/usecase/backfill/service.go`). Ob die Nichtanwendbarkeit einer
Transformationsregel im Run (Spalte fehlt in der Spaltenliste des Snapshots, Zielname
kollidiert) die Klasse `schema` braucht, entscheidet weder `ADR-0111` noch `ADR-0112`:
`ADR-0112` beschreibt die Nichtanwendbarkeit nur für den Erfassungspfad.

Deklaration: `slice-backfill-run-usecase`, Risiko §6 („An Architect gemeldet: Klasse `schema`
im Backfill-Pfad der Transformationen"), Ausgang *weiter offen*; Review F-9 (INFO).
