Zustand: **verkörpert** — Ausgang: **verkörpert** → `tools/schema/rollout.sh`
(Ablauf hinter `make schema-rollout`: Eingabe per `COPY`, Erzeugnisse per
`tar`-Stream nach `SCHEMA_ARTEFACT_DIR`, Default `.tmp/schema-rollout`, durch
`.tmp/` in `.gitignore` ausgenommen; der Rollout schreibt keine committete Datei).
Der Vertrag steht in `harness/targets/schema-rollout.md` §Erzeugnisse;
die Festlegung in
[`ADR-0142`](../../../../adr/0142-schema-rollout-erzeugnisse-ausserhalb-baum-eingabe-ohne-bind-mount.md)
· seit slice-schema-rollout-ohne-bind-mount.

Gegenstand: `make test-store`, `make test-replication`, `make test-integration`,
`make test-sdk-*-integration`, `make bench` und `make example-demo-up` rufen
`make schema-rollout` als Vorbedingung. `git status --short` nach diesen Läufen
ist leer (Verifikation `verifikation-slice-schema-rollout-ohne-bind-mount` §1,
gemessen im Lauf des Verifiers für `make test-store`, `make test-replication`,
`make test-integration`, `make example-demo-up`; die drei SDK-Runner und
`make bench` sind **hergeleitet** aus der gleichen Aufrufzeile, nicht gefahren).
Der Guard-Test `tools/harness/run-schema-rollout-guard-test.sh` prüft in Lauf 1
den Arbeitsbaum vor und nach dem Rollout. Ein Wrapper und ein Tabellentest für
eine Rücknahme existieren nicht mehr; ein weiterer Schreiber in eine committete
Datei ist ein neuer Beleg in dieser Ablage.

Zähler: 4× (Dateien unter `evidence/`).
