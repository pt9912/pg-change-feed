# Beleg: slice-039 (CI-Workflow und Dependabot)

Vorgang: slice-039 (wellenlos, erster Slice der CI/CD-Pipeline-Arbeit).

Fund: `.github/workflows/ci.yml` und `.github/dependabot.yml` wurden
ausschließlich statisch geprüft (`python3`/PyYAML-Parse, `yamllint`,
`actionlint`, `jsonschema` gegen die SchemaStore-Schemata) — vierfach
unabhängig, aber keine dieser Prüfungen ersetzt einen echten Lauf auf
einem GitHub-Actions-Runner. Zusätzlich ist unklar, ob `make test`
(Race-Detector, Debian-Toolchain-Image) im Runner ein anderes
Ressourcen-/Zeitverhalten zeigt als lokal — auch das ist ohne realen
Lauf nicht ausschließbar.

**Ausgang: weiter offen.** Löst sich erst mit dem ersten echten
PR/Push-getriggerten Lauf nach dem Merge (Sache des Nutzers, `slice-039`
§1 nennt das ausdrücklich als Out-of-Scope). Dieselbe Grenze gilt für
jeden künftigen GitHub-Actions-Workflow dieses Repos (`slice-040`:
`release.yml`, `hub-description.yml`, `image-scan.yml`,
`upstream-drift.yml`) — dort potenziell verschärft, weil diese zusätzlich
Registry-Secrets brauchen, die vor dem ersten echten Lauf ebenfalls
ungeprüft bleiben.

Quelle: `docs/plan/planning/in-progress/slice-039-ci-workflow-dependabot.md`
§6, `docs/reviews/verify-slice-039.md`.
