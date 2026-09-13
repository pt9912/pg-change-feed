Zustand: **verkörpert** — die real fehlende Löschausführung ist gebaut:
`RunRetentionUseCase`/`RunRetentionService`
(`internal/application/usecase/retention/service.go`, `slice-043`), der
Hintergrundzug `runRetentionCleanup`
(`internal/bootstrap/wiring.go`, `slice-044`), die Sichtbarkeit
blockierender Consumer über `cdc.retention_blockers`
(`tools/schema/schema.yaml`, `slice-045`) und die Metrik
`cdc_storage_bytes` (`tools/schema/nacharbeit-observability.sql`,
`slice-046`). Verkörpert in den vier genannten Zielorten · seit welle-13.
Zähler (abgeleitet): 0× — kein abgeschlossener Vorgang trug je einen Beleg
(die Beobachtung wurde nie über die 3×-Schwelle gezählt, sondern direkt
über eine dedizierte Feature-Welle aufgelöst, analog zu
`BEO-PGC/verwaltung-keine-sql-administration`s und
`BEO-PGC/schema-evolution-nicht-dynamisch`s direkter Auflösung unter der
Schwelle).
