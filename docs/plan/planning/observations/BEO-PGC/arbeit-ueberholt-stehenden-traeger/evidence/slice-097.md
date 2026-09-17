**Vorgang:** slice-097

**Fund:** Der Umzug der erzeugten Vertragsfläche (`internal/adapters/driving/grpc/streamv1`
→ `gen/cdc/stream/v1`) machte einen Satz in
`docs/plan/planning/next/slice-095-beispiel-clients-drei.md` §4 falsch, der
ihn nicht im Diff hatte: „die öffentliche Vertragsfläche (`gen/cdc/stream/v1`)
existiert nicht (`no required module provides package …`, Exit 1)" — bei
Niederschrift (`slice-095`s Rückführung `in-progress` → `next`) wahr, seit
diesem Umzug falsch.

Der Implementer-Suchlauf für §3.13 fand die sieben im Architect-Verdikt
benannten Träger vollständig (`Dockerfile`, `.dockerignore`,
`harness/image-hash.txt`, `harness/sensors/{coverage-gate,generated-sync}.md`,
`harness/mk/coverage.mk`, `harness/README.md`, `AGENTS.md`), verfehlte aber
`next/slice-095` — ein Träger außerhalb dieser Liste, in einem anderen
Planning-Verzeichnis. Gefunden hat ihn erst ein zweiter, unabhängiger
Suchlauf des Reviewers (`docs/reviews/review-slice-097.md`, Abschnitt
„Beobachtung außerhalb des Diffs"). In dieser Closure direkt korrigiert
(Planner-Trägerpflege).

Quelle: `docs/reviews/review-slice-097.md` (Beobachtung außerhalb des Diffs) ·
`docs/reviews/verify-slice-097.md` (§3) · `AGENTS.md` §3.13.
