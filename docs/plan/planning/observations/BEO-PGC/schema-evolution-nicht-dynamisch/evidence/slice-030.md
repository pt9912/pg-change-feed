# Evidence: slice-030

`slice-030`s neue Black-Box-Testfälle (`TestMVPSchemaChangeAddColumn`,
`TestMVPSchemaChangeIncompatibleTypeChange`,
`test/integration/integration_test.go`) belegen erstmals real gegen den
laufenden Compose-Stack, dass die `ADR-0015`-Folgepflicht nicht eingelöst
ist: dieselbe `schema_version` vor und nach einem echten
`ALTER TABLE ADD COLUMN`, kein real herstellbarer Negative-Fall für eine
inkompatible Typänderung.

Reviewer (`docs/reviews/review-slice-030.md`, F-1, HIGH) hat den Fund
unabhängig am Code verifiziert. Architect-Verdikt
(`docs/reviews/architect-verdict-slice-030-adr-0015.md`) bestätigt: `ADR-0015`
gilt unverändert fort, die Folgepflicht wurde nie eingeplant — kein
früherer Slice hat sie fälschlich als geliefert behauptet.
