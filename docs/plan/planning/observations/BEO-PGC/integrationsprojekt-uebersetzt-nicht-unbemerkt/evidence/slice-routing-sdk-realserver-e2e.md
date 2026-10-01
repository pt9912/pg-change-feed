**Vorgang:** slice-routing-sdk-realserver-e2e (Review F-2, MEDIUM)

**Fund:** Das C#-Integrationsprojekt übersetzte am Parent des Slice nicht: drei Aufrufe
`client.StreamChangesAsync(<CancellationToken>)` übergaben das Token positional als erstes Argument,
seit der Stream-Filter-Signatur steht dort `string? schema` (Compiler-Fehler `CS1503` in der Stufe
`integration`). Der Implementer fand es am Start des Slice beim Bau; kein Gate, kein Workflow und
keiner der Pack-Läufe hatte es zuvor gezeigt (der Reviewer las die Dockerfile-Stufen und suchte in
`.github`: kein Workflow ruft `test-sdk-*-integration` auf, `GATE_CHECKS` enthält keines der Ziele).
Der Fix im Slice ist reiner Test-Code (benanntes `cancellationToken:`); der Kotlin-Tier hat dieselbe
Form (`integrationTest` hängt bewusst nicht an `check`).

**Form (Ausprägung):** Test-Tier ohne regelmäßigen Übersetzungslauf. Schwere MEDIUM; vor dem Merge
vom Reviewer als Lücke benannt, außerhalb des Diffs, daher Register-Eintrag statt Fixrunde.

Quelle: `docs/reviews/review-slice-routing-sdk-realserver-e2e.md` (F-2) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-routing-sdk-realserver-e2e.md` (§6 Zeile F-2, §8 V-2). <!-- d-check:status-provenance -->
