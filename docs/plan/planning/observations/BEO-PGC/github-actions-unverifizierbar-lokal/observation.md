# BEO-PGC/github-actions-unverifizierbar-lokal

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die
CI/CD-Infrastruktur, keine eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: GitHub-Actions-Workflows
([`ADR-0051`](../../../../adr/0051-cicd-pipeline-github-actions.md)) lassen
sich nicht in Docker simulieren — jede Prüfung vor einem echten Push/PR
bleibt auf statische Verfahren beschränkt (YAML-Parse, `yamllint`,
`actionlint`, JSON-Schema-Validierung gegen die offiziellen
SchemaStore-Schemata). Ob ein Workflow auf dem echten GitHub-Actions-Runner
tatsächlich grün läuft — inklusive Ressourcen-/Zeitverhalten,
Registry-Zugangsdaten-Pfade, Runner-spezifischer Umgebungsunterschiede —
bleibt bis zum ersten realen Lauf unbewiesen.

Deklaration: `slice-039-ci-workflow-dependabot.md` §6.
