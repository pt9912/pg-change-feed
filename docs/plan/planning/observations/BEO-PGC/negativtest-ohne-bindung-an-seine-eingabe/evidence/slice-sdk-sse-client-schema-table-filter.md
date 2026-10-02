**Vorgang:** slice-sdk-sse-client-schema-table-filter (Review F-1, MEDIUM)

**Fund:** Die Tests der SSE-Beispiel-Clients (Go, C#, Kotlin) banden den URL-Bau (`StreamURL`) und, in C# und Kotlin, den Parser je für sich, aber nicht ihre Verbindung: das Entfernen des `-table`-Flags im Go-Beispiel (Reviewer-Mutation M5) und das Ersetzen von `cfg.table` durch `""` im `StreamURL`-Aufruf (M6) ließen `go test` grün, während die Zusage des Plans „Flag → Anfrage“ lautete. C# und Kotlin tragen dieselbe Lücke (hergeleitet, nicht gemutet); das vorhandene `-target` trug sie schon am Parent. Die Fixrunde (`81ad7643`) zog die Verdrahtung in testbare Funktionen (Go `parseConfig`/`config.streamURL`, C# `SseStream.StreamUrl(Config)`, Kotlin `SseStream.streamUrl(Config)`), und die Tests fahren die echten Flag-Argumente; der Verifier sah sechs Verdrahtungs-Mutationen (G1, G2, C1, C2, K1, K2) rot.

**Form (Ausprägung):** Verdrahtung zwischen zwei einzeln gebundenen Gliedern (Flag-Parser und URL-Bau): die Eingabeseite der Zusage „Flag → Anfrage“ ist das Flag-Argument am Programmeingang, nicht der Parameter der Teilfunktion. Schwere MEDIUM, daher eine Datei trotz Deckel; vor dem Merge vom Reviewer gefunden, die Regel hat gewirkt, kein Kandidat der Schärfung.

Quelle: `docs/reviews/review-slice-sdk-sse-client-schema-table-filter.md` (F-1, M5, M6) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-sdk-sse-client-schema-table-filter.md` (§3 Zeile 5, Mutationen G1, G2, C1, C2, K1, K2). <!-- d-check:status-provenance -->
