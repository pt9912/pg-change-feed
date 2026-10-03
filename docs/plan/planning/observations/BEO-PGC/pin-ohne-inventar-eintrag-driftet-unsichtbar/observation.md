# BEO-PGC/pin-ohne-inventar-eintrag-driftet-unsichtbar

**Sub-Area:** `*`/`PGC` (Build-, Test- und Gate-Konfiguration).

Die Beobachtung: Ein Upstream-Pin, der im Pin-Inventar von
[`ADR-0051`](../../../../adr/0051-cicd-pipeline-github-actions.md)
Entscheidung 7 keine Zeile führt, driftet, ohne dass ein Sensor es meldet.
Das Inventar ist die Liste dessen, was `make image-stale`, die
`make pin-stale-*`-Ziele und der nächtliche Workflow `upstream-drift.yml`
lesen; ein Pin außerhalb dieser Liste hat keinen Leser, seine Veralterung
fällt nur auf, wenn jemand ihn von Hand misst. Zwei Formen sind belegt:
der Pin `postgres:17-alpine` als YAML-Matrix-Wert in
`.github/workflows/e2e.yml` (das Messwerkzeug `tools/harness/pin-stale.sh`
liest Makefile-Variablen, keinen Matrix-Wert), und die fünf gepinnten
Basis-Images der SDK- und Beispiel-Dockerfiles (`dotnet/sdk`,
`dotnet/runtime`, `eclipse-temurin` JDK und JRE, `python:3.14-slim`).

Deklaration: Verifikationsbericht zu `slice-pin-digests-aktualisieren-2026-10`
§8.1 und §8.6.
