# Beleg: slice-012

Vorgang: slice-012 (Health-Endpoint per Heartbeat, Welle 4).

Fund: die Umsetzungs-Frage aus welle-3 (§6-Risiko slice-011, Architect-
Verdikt „kein Folge-ADR nötig") ist eingelöst — Heartbeat-Schreiber
(`HeartbeatPort` + PostgreSQL-Adapter), vierte SQL-Lese-View
(`cdc.heartbeat`) und ein `--healthcheck`-CLI-Modus (distroless-
Runtime-Zwang) sind real geliefert und real getestet (`make test-store`,
`make test-integration` mit echtem Docker-Health-Übergang
`starting`→`healthy`).

Quelle: docs/reviews/review-slice-012.md, docs/reviews/verify-slice-012.md.
