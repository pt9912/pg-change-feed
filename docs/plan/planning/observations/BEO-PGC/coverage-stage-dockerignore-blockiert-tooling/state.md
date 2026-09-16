Zustand: offen — Ausgang: **weiter offen** → beim Hinzufügen einer
weiteren Docker-Stage, die ein Skript unter `tools/` braucht (Kandidat:
`slice-050` Benchmark-Infrastruktur, `tools/bench-fixture.sh`-Muster),
`.dockerignore` UND die Alpine-Basis vorab prüfen, statt das
d-check-Muster blind zu kopieren. Kein Sensor denkbar, der einen
fehlenden `.dockerignore`-Ausnahme-Eintrag oder eine fehlende
`bash`-Installation vorab (ohne echten Docker-Build) erkennt — beides
zeigt sich erst am realen, roten Build. Zähler (abgeleitet): **2×** (evidence/slice-049.md,
evidence/slice-093.md) — unter der Schwelle. Der zweite Beleg ist der erste, der
kein Werkzeug-Fehler war, sondern eine **Lockerung**: die `.dockerignore`-Zeile
verletzte `ADR-0082`s `test-only`-Folgepflicht im Wortlaut, und der Architect hat
daraus `ADR-0085` gemacht (die Folgepflicht auf den **Zweck** geschärft, die
Ausnahme als **Klasse** statt als Liste). Der Eintrag war der Grund, warum die
Antwort eine Klasse wurde.
