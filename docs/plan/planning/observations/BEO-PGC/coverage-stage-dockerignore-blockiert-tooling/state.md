Zustand: offen — Ausgang: **weiter offen** → beim Hinzufügen einer
weiteren Docker-Stage, die ein Skript unter `tools/` braucht (Kandidat:
`slice-050` Benchmark-Infrastruktur, `tools/bench-fixture.sh`-Muster),
`.dockerignore` UND die Alpine-Basis vorab prüfen, statt das
d-check-Muster blind zu kopieren. Kein Sensor denkbar, der einen
fehlenden `.dockerignore`-Ausnahme-Eintrag oder eine fehlende
`bash`-Installation vorab (ohne echten Docker-Build) erkennt — beides
zeigt sich erst am realen, roten Build. Zähler (abgeleitet): 1×
(evidence/slice-049.md) — unter der Schwelle.
