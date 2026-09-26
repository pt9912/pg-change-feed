Zustand: **verkörpert** — Ausgang: **verkörpert** → `tools/schema/rollout-restore.sh`
(Wrapper: sichert `tools/schema/plan.yaml` und `tools/schema/down.sql` in ein
`mktemp`-Verzeichnis, führt das Kommando aus und stellt beide danach wieder her,
auch bei einem Fehlschlag) und `make test-rollout-restore` (Tabellentest des
Wrappers und Aufrufer-Prüfung: jedes Skript unter `tools/` und `examples/`, das
`make schema-rollout` aufruft, geht durch den Wrapper); der Vertrag steht in
`harness/targets/schema-rollout.md` §Erzeugnisse in Test-, Bench- und Beispiel-Läufen
· seit slice-harness-suchlauf-nachmessen
(Architect-Verdikt `architect-verdict-welle-backfill-bestand-lese-schritt`
§4.1, R3).

Gegenstand: `make test-store`, `make test-replication`, `make test-integration`,
`make test-sdk-*-integration`, `make bench` und `make example-demo-up` rufen
`make schema-rollout` mit Pflicht-Report und Rollback-Artefakt als Vorbedingung;
beide Dateien bleiben ohne Wrapper nach dem Lauf verändert. Der Guard-Test
`tools/harness/run-schema-rollout-guard-test.sh` trägt seine eigene Sicherung; im
Betrieb bleiben beide Dateien das Erzeugnis des Rollouts.

Belegt (Verifikation `verifikation-slice-harness-suchlauf-nachmessen` §1, gemessen im
Lauf des Verifiers): `make test-store` und `make test-replication` enden mit Exit 0,
danach ist `git status --short tools/schema` leer; die Mutation „Wiederherstellung
entfernt“ färbt `make test-rollout-restore` rot. Benannte Grenzen des Wrappers
(Review F-5, F-7): zwei gleichzeitige Aufrufe teilen dieselben zwei Dateien, ein
`SIGKILL` lässt die Erzeugnisse verändert, die Aufrufer-Prüfung liest nur Einzelzeilen.
Ein weiterer Schreiber in eine committete Datei ist ein neuer Beleg in dieser
Ablage; der Slice hat keinen weiteren gefunden.

Zähler: 4× (Dateien unter `evidence/`).
