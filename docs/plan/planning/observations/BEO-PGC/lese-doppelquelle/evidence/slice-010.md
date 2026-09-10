# Beleg: slice-010

Vorgang: slice-010 (Lesen-Vollabdeckung und SQL-Schnittstelle, Welle 3).

Fund: §6-Risiko (a) — SQL-Views (`active_tables`, `consumer_status`,
`changes`) und die Go-Use-Case-Lesepfade implementieren dieselbe
Lese-Semantik zweimal, ohne gemeinsamen Vertrags-Test. `ADR-0046`
erlaubt die Views (keine Entscheidungslogik), löst aber nicht die
Drift-Gefahr zwischen den beiden Implementierungen.

Quelle: slice-010-Plan §6 (Closure-Notiz), `ADR-0046`.
