Zustand: verkörpert → geschärfte Test-Kadenz-Regel liegt in
`harness/README.md` §Sensors, `make schema-rollout`-Bindung (seit
slice-015). Zähler (abgeleitet): 8× (evidence/slice-006.md,
evidence/slice-010.md, evidence/slice-015.md, evidence/slice-016.md,
evidence/slice-036.md, evidence/slice-066.md,
evidence/slice-backfill-change-origin.md,
evidence/slice-transformationen-antragsweg-schema.md — Funktion mit
`jsonb`-Parameter: d-migrate meldet und rendert sie als `json`, Parameterart der
Nacharbeit-Funktion ist `json`, `ADR-0125`).
Nachrichtlich: Von den drei bislang betroffenen Objektklassen sind zwei
technisch aufgelöst — das CHECK-Constraint `chk_change_operation` seit
slice-015 (deklarativ in `tools/schema/schema.yaml`), die drei Views seit
slice-016 (ebenso deklarativ, `source_dialect`+`columns:`). Die dritte,
mit slice-036 neu belegte Klasse — SQL-Funktionen/Prozeduren — bleibt
**offen**: `schema migrate --execute` bricht für jede über den
`functions:`-Knoten deklarierte Funktion mit `POST_EXECUTE_DRIFT` (Exit 5)
ab, real isoliert reproduziert mit einer trivialen No-Arg-Funktion;
`cdc.enable_table`/`cdc.disable_table` laufen deshalb weiterhin über die
Ausweichform `tools/schema/nacharbeit-administration.sql`.
`tools/schema/nacharbeit-views.sql` und
`tools/schema/nacharbeit-operation-check.sql` sind beide gelöscht;
`tools/schema/nacharbeit-administration.sql` bleibt bestehen, bis
d-migrate die Funktionsklasse ebenso wie Views seit 1.3.1 aus dem
Post-Compare-Fingerabdruck ausblendet.
Mit slice-066 kommt eine vierte Objektklasse hinzu: eine **CHECK-Klausel
an einer bereits bestehenden Tabelle** (nicht die Erstanlage) ist
ebenfalls nicht deklarativ tragbar — Exit 5, die bestehende Klausel
entfällt beim Änderungsversuch. Sie liegt deshalb ebenfalls in
`nacharbeit-administration.sql`, neben der Funktionsklasse.
Mit `slice-backfill-change-origin` kommt eine fünfte Objektklasse hinzu: die
**Signaturänderung einer bereits angelegten View** (`ReplaceView` mit
`VIEW_SIGNATURE_INCOMPATIBLE`, Exit 8 im Precheck). Sie läuft nicht über eine
Nacharbeit-Datei, sondern über den Vorlauf `DROP VIEW` im Target
`schema-rollout` mit der Wache `tools/schema/rolloutguard`
(`harness/targets/schema-rollout.md`).
Die verkörperte Regel bleibt als Betriebsdisziplin bestehen (bewährt: der
d-migrate-Fix traf ein, der reale Test lief erst nach dem expliziten
Fix-Signal).
