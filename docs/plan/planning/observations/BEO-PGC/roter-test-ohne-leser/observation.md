# BEO-PGC/roter-test-ohne-leser

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft das Verhältnis
zwischen einem vorhandenen Beleg und der Frage, ob ihn ein Gate abholt, keine
eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Ein Test ist **rot**, sein Target läuft aber **weder im Gate
noch in der CI** — es läuft nur, wenn jemand es von Hand aufruft. Der Beleg
existiert also, nur liest ihn niemand: der Fehler überlebt beliebig viele
Slices, und beide Seiten sind sich sicher, geprüft zu sein (die Testseite, weil
sie einen Test hat; die Gate-Seite, weil sie grün ist). Belegt an `slice-080`:
`make test-replication`s tier-weites `go test ./...` scheitert seit `ADR-0050` an
einem Fixture, das die seither gelesenen Tabellen nicht mitbringt — sichtbar
geworden erst, als ein neuer Träger (die DB-Adapter-Coverage) den Lauf aufrief.
