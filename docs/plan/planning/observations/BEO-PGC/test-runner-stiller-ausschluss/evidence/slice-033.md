# Evidence: slice-033

`slice-033` musste `tools/harness/run-integration-tests.sh`s bisherigen
einzelnen `go test ./test/integration/...`-Aufruf in zwei
`-run`-gefilterte Aufrufe aufteilen (sechs benannte Testfunktionen vorn,
`TestMVPSchemaChangeIncompatibleTypeChange` separat und zuletzt) — Grund:
der von diesem Slice gemeldete `schema`-Fehler beendet den geteilten
Feed-Container dauerhaft (`restart: "no"`), ein gemeinsamer Testlauf
hätte nachfolgende Tests mitgerissen.

Der Reviewer (`docs/reviews/review-slice-033.md`, F-2 MEDIUM) fand dabei,
dass dieser Split real das Risiko eines stillen Ausschlusses künftiger
Testfunktionen trägt: `go test -run` meldet keinen Fehler, wenn eine neue
Testfunktion in keinem der beiden Muster auftaucht, solange andere
Funktionen im selben Aufruf matchen.
