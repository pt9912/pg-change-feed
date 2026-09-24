# BEO-PGC/backfill-adapter-startwerte-ohne-messung

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die Postgres-Adapter des
Backfill-Runs, keine eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Die drei Postgres-Adapter des Backfill-Runs (`postgresstorage`) tragen
Setzungen ohne Messung. Zwei Fristen begrenzen jede Operation, weil der Use Case
`Finish`, `InterruptRunning` und `Rollback` auf einem vom Abbruch gelösten Kontext ohne Frist
ruft: `backfillStateTimeout` (30 s je Zustands-Operation) und `backfillBlockTimeout` (5 min je
Block und Commit) in `internal/adapters/driven/postgresstorage/backfill.go`, im Kommentar als
„Startwerte ohne Messung" gekennzeichnet. Dazu die Einfügeform: `AppendBlock` schreibt je
Change eine `Exec`-Anweisung in der einen Schreibtransaktion; diese Grenze steht im Slice-Plan
und im Review, am Code des Adapters steht sie nicht. Ob 5 min je Block bei der gemessenen
Blockgröße und Zeilenbreite trägt und ob die zeilenweise Einfügung die Kopierdauer bestimmt,
ist ungemessen.

**Warum das zählt:** Eine Frist unter der Blockdauer bricht einen gesunden Run mit einem
Zeitfehler ab (Klasse `storage`, Run `failed`), eine Einfügeform ohne Messung verschiebt die
Kopierdauer, an der die Warn-Richtgröße hängt. Kein Gate liest die Dauer eines Blocks.

Deklaration: `slice-backfill-run-store`, Risiko §6 zehnter Punkt („Startwerte ohne
Messung"), Review F-7 (INFO), Ausgang *weiter offen*.
