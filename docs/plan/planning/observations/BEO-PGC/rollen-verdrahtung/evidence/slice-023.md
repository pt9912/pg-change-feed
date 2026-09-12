# Beleg: slice-023 (Rollen-spezifische DSN-Verdrahtung)

Vorgang: slice-023 — der laut Architect-Verdikt (`architect-review-welle-6.md`
Zug 2, Ausgang `geplant`) zugewiesene Folge-Slice.

Fund: `internal/bootstrap/wiring.go` verdrahtet jetzt jeden Aufrufer über
die zur Aufgabe passende Rolle (`cdc_capture`/`cdc_admin`/`cdc_reader`,
[`ADR-0047`](../../../../../adr/0047-rollenspezifische-dsn-verdrahtung.md),
[`ADR-0048`](../../../../../adr/0048-heartbeat-grant-korrektur-select-ergaenzung.md)),
real getestet gegen PostgreSQL (`internal/bootstrap/roles_wiring_test.go`,
mehrfach durch Reviewer und Verifier unabhängig mit eigenen
Mutationstests gegengeprüft).

**Ausgang: eingetreten.** Die Beobachtung ist technisch vollständig
aufgelöst — die gemeinsame Instanz-DSN existiert nicht mehr, jeder
Aufrufer ist auf sein Least-Privilege beschränkt.

Quelle: `docs/reviews/review-slice-023.md`, `docs/reviews/verify-slice-023.md`.
