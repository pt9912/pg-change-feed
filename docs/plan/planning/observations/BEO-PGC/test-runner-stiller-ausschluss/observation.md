# BEO-PGC/test-runner-stiller-ausschluss

**Sub-Area:** Test-Infrastruktur (`tools/harness/run-integration-tests.sh`;
Sub-Area-Kürzel `PGC` aus der Modus-Deklaration)

Die Beobachtung: `tools/harness/run-integration-tests.sh` filtert
`test/integration/integration_test.go`s Testfunktionen über mehrere
separate `go test -run '^(...)$'`-Aufrufe statt eines einzigen,
ungefilterten `go test ./test/integration/...`. Eine künftig zu dieser
Datei hinzugefügte Testfunktion, die in **keinem** der bestehenden
`-run`-Muster auftaucht, wird von `make test-integration` dauerhaft und
**stillschweigend** nie ausgeführt — `go test -run` meldet keinen Fehler,
solange mindestens eine andere Funktion im selben Aufruf matcht.
