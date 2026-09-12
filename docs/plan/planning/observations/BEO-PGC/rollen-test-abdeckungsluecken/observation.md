# BEO-PGC/rollen-test-abdeckungsluecken

**Sub-Area:** Bootstrap/Rollen-Verdrahtung (Sub-Area-Kürzel `PGC` aus der
Modus-Deklaration)

Die Beobachtung: Die realen Rollen-Tests (`internal/bootstrap/roles_wiring_test.go`)
belegen die Kern-Fitness-Function aus `ADR-0047`/`ADR-0048`, haben aber
zwei Deckungslücken, die derselbe Vorgang (`slice-023`-Verifikation)
aufdeckte: (1) der Heartbeat-Grant-Test entzieht/erteilt die Rechte
selbst per `REVOKE`/`GRANT` im Test und liest nie den tatsächlichen
`nacharbeit-roles.sql`-Inhalt — eine reale Regression in der
Rollout-Datei selbst bliebe unsichtbar und würde vom Test-`Cleanup`
sogar still „repariert". (2) Replication-Stream- und ACK-Adapter (beide
`cdc_capture`-gebunden) sind gegen Rollen-Vertauschung nicht
testgesichert, weil `tools/harness/run-replication-tests.sh` keine
rollenbeschränkten Login-Test-Identitäten bereitstellt.

Deklaration: `internal/bootstrap/roles_wiring_test.go` (Code-Kommentar),
`docs/reviews/verify-slice-023.md` (Nachtrag, V-1/V-2-Zusatzbefunde).
