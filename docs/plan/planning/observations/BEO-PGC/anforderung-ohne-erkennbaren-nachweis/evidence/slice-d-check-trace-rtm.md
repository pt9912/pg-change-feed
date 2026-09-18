# Beleg: slice-d-check-trace-rtm

Vorgang: `slice-d-check-trace-rtm` — `d-check --trace`/Requirements
Traceability Matrix verdrahtet (Welle `welle-d-check`).

Fund: Nach Aktivierung von `trace.requirements.id-pattern` **und**
`trace.coverage` (auf `docs/user/e2e-abdeckung.md`) meldet `make doc-trace`
real **7 von 76** Anforderungen als Waise: `LH-FA-CFG-006`, `LH-FA-CON-002`,
`LH-FA-DAT-002`, `LH-FA-DAT-003`, `LH-FA-SST-001`, `LH-FA-SST-005`,
`LH-QA-REL-004`. Unabhängig vom Implementer (Planungsmessung), Reviewer und
Verifier reproduziert — alle drei Läufe liefern identisch 76/7 (Verifier
zusätzlich: temporäre Entfernung von `trace.coverage` liefert 76/9, exakte
Differenzmenge `LH-FA-CAP-002`/`LH-FA-CAP-003`).

Der Slice-Plan (§1 Abgrenzung, „Sanierung der sieben realen Waisen") schließt
die inhaltliche Prüfung bewusst aus: ob echte Lücke oder fehlende Zitierung,
ist pro Anforderung zu klären, kein Konfigurationsvorgang. Damit bleibt der
Fund ohne diesen Register-Eintrag ohne jeden weiteren Leser — `--trace` ist
advisory, kein Gate liest die sieben Kennungen ein zweites Mal.

Quelle: `docs/plan/planning/done/welle-d-check/slice-d-check-trace-rtm.md` §1/§2 ·
der Verifikationsbericht zu `slice-d-check-trace-rtm` §A/§G.
