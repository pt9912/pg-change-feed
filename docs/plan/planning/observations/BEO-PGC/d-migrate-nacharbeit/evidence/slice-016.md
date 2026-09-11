# Beleg: slice-016 (d-migrate-1.3.1-Views-Retirement)

d-migrate-Pin auf v1.3.1 gehoben (Digest
`sha256:862dfb04c34dd17278b1bab46961363c12eeb8d464cf1776565d6285603d2c89`).
Die drei Views (`active_tables`, `consumer_status`, `changes`) real gegen
den neuen Pin getestet, mit `source_dialect: postgresql` und
`columns:`-Signatur (Anhang F.11 des d-migrate-Anwenderhandbuchs)
deklarativ in `tools/schema/schema.yaml` überführt. Dreifach unabhängig
real reproduziert (Implementer, Reviewer, Verifier — je eigene
Testcontainer-Instanz):

- `schema migrate --execute` konvergiert gegen eine leere DB (Erstanlage,
  `CreateView`, Exit 0) UND gegen eine bereits migrierte DB (Folgelauf,
  `ReplaceView`, Exit 0 — kein `VIEW_SIGNATURE_UNKNOWN`).
- Gegenprobe: dieselbe Situation ohne `columns:`-Signatur reproduziert
  `VIEW_SIGNATURE_UNKNOWN` auf allen drei Views — belegt die Kausalität.

`tools/schema/nacharbeit-views.sql` ist gelöscht. Die Views-Hälfte der
Beobachtung ist damit **technisch aufgelöst**. Der mit slice-015
verkörperte Ausgang (Test-Kadenz-Regel in `harness/README.md`) bleibt als
Regel bestehen — sie hat sich bewährt: das genannte Fix-Signal traf ein,
der reale Test lief erst danach, nicht bei jedem Pin-Bump dazwischen.

**Nebenfund (eigener Register-Eintrag):** Ein zweiter `make
schema-rollout`-Lauf gegen eine bereits migrierte DB blockiert mit Exit 8
(`DropView` auf `cdc.heartbeat`/`cdc.metrics`, die außerhalb des
neutralen Modells liegen) — pin-unabhängig, real gegen 1.3.0 UND 1.3.1
reproduziert, keine Regression dieses Slices. Siehe
`BEO-PGC/schema-rollout-fremdobjekte`.
