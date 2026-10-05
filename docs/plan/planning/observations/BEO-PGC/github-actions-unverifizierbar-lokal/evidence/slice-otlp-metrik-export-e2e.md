**Vorgang:** slice-otlp-metrik-export-e2e (Planner-Closure).

**Fund:** Der Slice trug lokal `make test-integration` Exit 0 und alle Gates grün. Der erste
Post-Push-Lauf von e2e.yml (`37225751324`, Commit `19e83251`) war am Leg PostgreSQL 18 rot, Phase
„Wiederaufnahme nach docker start“. Ursache war ein Zählfehler des Runners, kein Runner-Umgebungs-
Unterschied im engeren Sinn (Collector legt seine Ausgabedatei beim Start leer neu an, siehe
`BEO-PGC/zaehlbasis-nach-neustart-ruecksetzende-ausgabe`); lokal blieb er verdeckt. Nach `a19c28cc`
und dem Abdeckungs-Neuschrieb `759d2f2f` beide Legs `success` (Lauf `37254000299`, PostgreSQL 17 etwa
27 min, PostgreSQL 18 etwa 29 min bei Limit 60 min).

**Form:** die Regel in `AGENTS.md` §3.10 hielt: das Risiko blieb offen, der Rotlauf wurde vor der
Closure gefunden und behoben. Neu gegenüber den früheren Belegen: der Rotlauf lag in der Test-Logik,
nicht im Workflow; ein lokaler Grünlauf eines Runners mit Zeitfenstern belegt kein Verhalten auf dem
gehosteten Runner.
