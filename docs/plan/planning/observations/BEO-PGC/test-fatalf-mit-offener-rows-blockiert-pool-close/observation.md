# BEO-PGC/test-fatalf-mit-offener-rows-blockiert-pool-close

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die Testhygiene
der Store-Adapter-Tests gegen die reale PostgreSQL, keine eigene Sub-Area im Sinn
der Modus-Deklaration).

Die Beobachtung: Ein Store-Test iteriert über die `pgx.Rows` einer
`pool.Query`-Abfrage und bricht mit `t.Fatalf` innerhalb der Schleife ab, bevor
`rows.Close()` lief. Die Verbindung bleibt aus dem Pool ausgeliehen; das
`pool.Close()` im Cleanup wartet auf ihre Rückgabe, und der Testlauf blockiert,
statt rot zu enden. Ein roter Test wird so zum hängenden Lauf. Kein Gate kann das
unterscheiden von einem langsamen Lauf.

Herkunft der Aussage: Befund des Implementers von `slice-backfill-change-origin`
beim Schreiben der neuen Store-Tests, im Auftrag an den Planner übernommen und
**nicht** vom Planner reproduziert. Die neuen Tests tragen die Abhilfe
(`defer rows.Close()` samt Kommentar, `sqlviews_test.go` und `store_test.go`).
Gemessen bei der Closure (`git grep -n` gegen `HEAD` in
`internal/adapters/driven/postgresstorage/*_test.go`): vier Aufrufe
`pool.Query(` in Tests, zwei davon mit `defer` (die neuen Tests), die beiden
vorbestehenden in `sqlviews_test.go` (Bereichs- und Filter-Lesen, Zeilen 229 und
253) schließen die `Rows` erst nach der Schleife und tragen darin ein
`t.Fatalf` beim `Scan`-Fehler — dieselbe Form. Ob dort ein Hängen real eintritt,
ist nicht gemessen; ein Audit der übrigen Bestandstests auf offene `Rows` hat
nicht stattgefunden.

**Warum das zählt:** Die Klasse trifft den Fehlerpfad eines Tests, den nur ein
Fehlschlag ausführt — im grünen Regelbetrieb nie sichtbar. Ein Hängen ohne
Meldung verzögert die Diagnose des eigentlichen Fehlers.

Deklaration: `slice-backfill-change-origin` (Implementer-Befund, neue Tests).
