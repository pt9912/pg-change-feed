Zustand: verkörpert → geschärfte Test-Kadenz-Regel liegt in
`harness/README.md` §Sensors, `make schema-rollout`-Bindung (seit
slice-015). Zähler (abgeleitet): 4× (evidence/slice-006.md,
evidence/slice-010.md, evidence/slice-015.md, evidence/slice-016.md).
Nachrichtlich: Beide ursprünglich betroffenen Fälle sind jetzt technisch
aufgelöst — `chk_change_operation` seit slice-015 (deklarativ in
`tools/schema/schema.yaml`), die drei Views seit slice-016 (ebenso
deklarativ, `source_dialect`+`columns:`). `tools/schema/nacharbeit-views.sql`
und `tools/schema/nacharbeit-operation-check.sql` sind beide gelöscht.
Die verkörperte Regel bleibt als Betriebsdisziplin bestehen (bewährt: der
d-migrate-Fix traf ein, der reale Test lief erst nach dem expliziten
Fix-Signal).
