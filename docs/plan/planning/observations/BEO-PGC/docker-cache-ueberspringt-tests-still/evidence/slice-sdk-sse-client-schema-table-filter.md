**Vorgang:** slice-sdk-sse-client-schema-table-filter

**Fund:** `make sdk-pack-python`, `make sdk-pack-csharp`, `make sdk-pack-kotlin`, `make examples-csharp`
und `make examples-kotlin` endeten beim Verifier mit Exit 0 aus dem Docker-Schicht-Cache, alle
Bau-Schritte `CACHED`, ohne Testzeile. Der Plan hatte die Falle vorab benannt (DoD (A), (B): der
Belegbefehl ist der Bau mit `--no-cache` der Test-Stufe oder eine rote Mutation). Der Verifier fuhr
`docker build --no-cache --target build` je Wurzel und las die gedruckten Zahlen (Python
`150 passed`, C# `Passed: 159`, C#-SSE-Beispiel `Passed: 25`; Kotlin ohne Zahl im grünen Lauf,
`BUILD SUCCESSFUL`) und fuhr zehn Einzelmutationen, die rot wurden.

Quelle: `docs/reviews/verifikation-slice-sdk-sse-client-schema-table-filter.md` (§1 Cache-Befund, §3 Zeile 1 und 2). <!-- d-check:status-provenance -->
