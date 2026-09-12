# BEO-PGC/test-isolation-geteilter-zustand

**Sub-Area:** Store-Adapter-Tests (Sub-Area-Kürzel `PGC` aus der
Modus-Deklaration)

Die Beobachtung: Mehrere `postgresstorage`-Paket-Tests mutieren den
geteilten `cdc`-Schema-/Tabellen-Zustand auf der Instanz, die sich alle
Pakete über `CDC_STORE_TEST_DSN` teilen — unskopiert (`DELETE` ohne
`WHERE` auf `cdc.consumer`/`cdc.consumer_position`) oder vollständig
(`DROP SCHEMA cdc CASCADE` samt partiellem Wiederaufbau, der die per
d-migrate ausgerollten Consumer-/View-Objekte nicht mitträgt). Bei
paralleler Paketausführung (Go-Standardverhalten) ist das ein reales
Race-Risiko gegen jedes andere Paket, das dieselbe Instanz nutzt. Statt
die Ursache (die unskopierten Löschungen/Drops selbst) zu beheben, trägt
das Repo bislang Workarounds, die die Symptomatik über die
Ausführungsreihenfolge umgehen.

Deklaration: `internal/adapters/driven/postgresstorage/sqlviews_test.go`
(Dateinamen-Sortierung erzwingt eine Ausführungsreihenfolge gegen
`store_test.go`), `tools/harness/run-store-tests.sh` (paketweise
Vorzieh-Isolierung für `internal/bootstrap`, `review-slice-022.md` F-2).
