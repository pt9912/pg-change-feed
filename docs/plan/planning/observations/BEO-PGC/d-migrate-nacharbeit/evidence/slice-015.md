# Beleg: slice-015 (d-migrate-1.3.0-Retirement)

d-migrate-Pin auf v1.3.0 gehoben; real gegen einen frischen Rollout
getestet (Implementer + unabhängig Verifier, zwei getrennte
Testcontainer-Instanzen, gleiches Ergebnis):

- **`chk_change_operation` (CHECK-Ausdruck): konvergiert.** Exit 0, keine
  Post-execute-Drift — die Operation ist jetzt ausdrückbar, die
  Ausweichform (`nacharbeit-operation-check.sql`) ist zurückgebaut, der
  Constraint lebt deklarativ in `tools/schema/schema.yaml`.
- **Die drei Views (`active_tables`, `consumer_status`, `changes`):
  konvergieren weiterhin nicht.** Exit 5 ("Post-execute compare detected
  drift"), einzeln und kombiniert reproduziert, sowohl im Implementer- als
  auch im unabhängigen Verifier-Lauf. Die Ausweichform
  (`nacharbeit-views.sql`) bleibt bestehen.

**Neue Erkenntnis gegenüber slice-006/slice-010:** Eine eigene, isolierte
Reproduktion (außerhalb des Arbeitsbaums, Repro-Paket an d-migrate
übergeben) zeigt, dass die Drift bei frisch angelegten `CREATE VIEW`-
Statements auftritt, die **Teil desselben Migrations-Plans** sind — laut
`spec/cli-spec.md:940-951` sollte der Post-Compare nur Objekte prüfen, die
der Plan *nicht* angefasst hat. Das CHECK-Constraint im selben Lauf löst
korrekt keine Drift aus; nur `CREATE VIEW` tut es. Der Rollout-Report
selbst (`plan.yaml`) zeigt für den betroffenen Lauf `status: ok,
exitCode: 0`, obwohl der Prozess mit Exit 5 endet — die Drift-Information
steht nirgendwo im Report. Dieser Befund wurde an das d-migrate-Team
gemeldet; die Rückmeldung bestätigt einen eigenen internen Check und
kündigt eine Fehlerbehebung an (Stand 2026-09-11: Arbeit begonnen, kein
Release-Termin bekannt).

**Ausgang für dieses Vorkommen:** die Views-Hälfte der Beobachtung bleibt
bestätigt offen — jetzt mit einem von d-migrate anerkannten, in Arbeit
befindlichen Bug als Ursache statt einer vermuteten Werkzeug-Grenze.
