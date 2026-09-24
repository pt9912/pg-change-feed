# BEO-PGC/test-schreibt-in-committete-datei

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft Test- und
Sensor-Läufe, die ein committetes Erzeugnis als Nebenwirkung überschreiben, keine
eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Ein Test- oder Sensor-Lauf ruft ein Target auf, dessen Ausgabe in
einer committeten Datei landet, und lässt sie im Arbeitsbaum verändert zurück.
Belegt am Pflicht-Report `tools/schema/plan.yaml` und am Rollback-Artefakt
`tools/schema/down.sql`, die `make schema-rollout` je Lauf schreibt (Vorgabe des
Targets): der Guard-Test
`tools/harness/run-schema-rollout-guard-test.sh` (seine Läufe im Arbeitsbaum
rufen das Target ohne `-C`) hinterließ beide Dateien modifiziert (Review F-6, INFO), und `make test-store`
sowie `make test-replication` überschreiben über `tools/schema/apply-rollout.sh`
die `target`-Zeile von `plan.yaml` (Verifikation V-6, INFO). Der Arbeitsbaum ist
nach dem Lauf dirty; ein Commit des Erzeugnisses trägt dann die Bezeichnung des
zuletzt gelaufenen Tiers statt die des Compose-Rollouts.

Behoben ist nur der Guard-Test: Fixrunde 2 sichert beide Dateien beim Start in ein
`mktemp`-Verzeichnis und stellt sie in `cleanup` wieder her (jeder Ausgang). Die
Läufe der zwei DB-Tiers bleiben, wie sie sind; der Verifier ordnet das als
Eigenschaft ihres Rollout-Pfads ein, nicht als Befund dieses Slice.

**Warum das zählt:** Eine committete Datei, die je Testlauf mitgeschrieben wird,
kommt mit dem Ziel des zuletzt gelaufenen Tiers in den nächsten Commit, wenn
niemand sie zurücknimmt — kein Gate liest den Inhalt der Ziel-Zeile
(`5f126971` musste sie nachträglich auf die Compose-Bezeichnung zurückstellen).

Deklaration: `slice-backfill-change-origin`.
