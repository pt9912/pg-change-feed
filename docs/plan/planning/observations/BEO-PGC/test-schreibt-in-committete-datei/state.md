Zustand: **geplant** — Ausgang: **geplant** →
`slice-harness-suchlauf-nachmessen`,
dritter Liefer-Punkt: `tools/schema/apply-rollout.sh` sichert `tools/schema/plan.yaml`
und `tools/schema/down.sql` vor dem Rollout und stellt sie nach dem Lauf wieder her
(Vorbild: der Guard-Test `tools/harness/run-schema-rollout-guard-test.sh`) ·
seit welle-backfill-bestand
(Architect-Verdikt `architect-verdict-welle-backfill-bestand-lese-schritt`
§4.1, R3).

Gegenstand: `make test-store` und `make test-replication` rufen `apply-rollout.sh`,
das `make schema-rollout` mit Pflicht-Report und Rollback-Artefakt ausführt; beide
Dateien bleiben nach dem Lauf verändert.

Zähler: 4× (Dateien unter `evidence/`).
