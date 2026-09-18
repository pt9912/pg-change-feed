# Welle 17 — E2E-Testbelege für fünf testfreie Lastenheft-Kennungen — Closure-Notiz

**Welle:** welle-17
**Abschluss:** 2026-09-14
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — *was gelernt wurde*: geliefert · was
funktionierte · was anders lief. Mit ID-Bezug, wo es einen gibt.

- Fünf zuvor testfreie Lastenheft-Kennungen tragen jetzt einen realen
  E2E-Beleg: [`LH-FA-SCH-003`](../../../../spec/lastenheft.md) (entfernte
  Spalten), [`LH-FA-DAT-006`](../../../../spec/lastenheft.md)
  (Metadaten-Erweiterbarkeit), [`LH-QA-OPS-005`](../../../../spec/lastenheft.md)
  (Upgrade-Sicherheit), [`LH-QA-POR-001`](../../../../spec/lastenheft.md)
  (PostgreSQL-Versionsportabilität), [`LH-QA-POR-002`](../../../../spec/lastenheft.md)
  (Linux-Zielplattform).
- [`ADR-0058`](../../adr/0058-testansatz-fuenf-luecken.md) (Accepted) traf
  alle fünf Teilentscheidungen vorab; zwei ihrer Testformen erwiesen sich
  bei der Umsetzung als real nicht tragfähig und wurden über eigene
  Folge-ADRs korrigiert — [`ADR-0063`](../../adr/0063-lh-fa-sch-003-testform-korrektur.md)
  (Supersedes Entscheidung 1) und [`ADR-0064`](../../adr/0064-lh-qa-ops-005-testansatz-korrektur.md)
  (Supersedes Entscheidung 3); die übrigen drei Teilentscheidungen (2, 4, 5)
  wurden unverändert umgesetzt.
- `slice-062`: zwei neue E2E-Testfunktionen in
  `test/integration/integration_test.go` — `TestE2ESchemaChangeDropColumn`
  (`ADR-0063`-Testform: reale Spaltenentfernung löst denselben
  `ErrIncompatibleSchemaChange`-Pfad aus wie eine inkompatible
  Typänderung) und `TestE2EChangeTableMetadataExtensibility` (realer
  additiver `ALTER TABLE cdc.change ADD COLUMN`-Beleg mit `t.Cleanup`).
  Dabei real gefunden und behoben: Zwischen den beiden container-
  beendenden Schema-Testfunktionen reicht ein bloßer `docker start` nicht
  — Replication-Slot-Neuanlage und Schema-Version-Nachtrag waren nötig.
- `slice-063`: eine neue Phase in
  `tools/harness/run-integration-tests.sh` bildet einen simulierten
  Anwendungs-Upgrade-Zyklus nach — `ADR-0064`s korrigierter Mechanismus
  (`$COMPOSE up -d --force-recreate --no-deps pg-change-feed`, realer
  Container-Tausch) statt des in `ADR-0058` ursprünglich vorgesehenen
  zweiten `make schema-rollout`-Laufs, der real mit Exit 8 auf vier
  Fremdobjekten blockierte (das Blocker-Protokoll zu `slice-063`).
- `slice-064`: `compose.yaml`s `postgres`-Image ist auf eine
  `${PG_TEST_IMAGE}`-Interpolation umgestellt;
  `.github/workflows/e2e.yml` trägt eine `strategy: matrix:` über die
  beiden `SPEC-012`-Digests (PostgreSQL 17, PostgreSQL 18) mit
  `fail-fast: false`. Real belegt: Run `34822131377`, beide Legs
  `completed`/`success`.
- `slice-065`: ein neuer, benannter Schritt in
  `.github/workflows/ci.yml`, unmittelbar nach `Checkout` und vor
  `Gates`, prüft `uname -s`/`go env GOOS` gegen `Linux`/`linux` mit
  sichtbarem Fehlschlag bei Abweichung. Real belegt: Run `34823976027`
  (Commit `5d7672b`), `completed`/`success`.

## Was hat funktioniert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3.

- Beide Male, in denen `ADR-0058`s ursprüngliche Testannahme real nicht
  trug (`slice-062`/`LH-FA-SCH-003`, `slice-063`/`LH-QA-OPS-005`), stoppte
  der Implementer korrekt über die Rückführung, statt den Befund zu
  umgehen, und meldete ihn mit realer Reproduktion (Modul 8
  §Konflikt-Pfad). Beide Male fand der nachfolgende Architect-Zug
  (`ADR-0063`, `ADR-0064`) einen eng geschnittenen, minimal-invasiven
  Ersatzmechanismus, der beim zweiten Versuch sofort real funktionierte.
- `slice-064`s Mechanismus (Compose-Interpolation, Matrix-Job) lief im
  ersten Versuch real grün — beide Legs beim ersten Push.
- `slice-065`, der kleinste Slice der Welle (Zwei-Zeilen-Schritt), lief
  ebenfalls im ersten Versuch grün, keine Fixrunde nötig.
- Über alle vier Slices hinweg fanden Reviewer und Verifier zusammen 0
  HIGH, 0 MEDIUM, 0 LOW, 4 INFO — keine Fixrunde in der gesamten Welle
  nötig.
- Der reale Welle-Closure-Trigger (grüner `e2e.yml`-Matrix-Lauf) lag
  bereits vor Beginn dieser Closure dreifach unabhängig bestätigt vor
  (Planner-Koordinator, Reviewer, Verifier bei `slice-064`) — kein neuer
  Lauf für diese Closure nötig.

## Was ging anders als geplant?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — jede Zeile moeglichst mit der Konsequenz,
die daraus schon gezogen wurde (Folge-Slice, Spec-Version).

- `ADR-0058`s Happy-Path-Annahme für `LH-FA-SCH-003` (stille
  Spalten-Auslassung) widersprach dem bereits von `ADR-0059` akzeptierten
  Verhalten (`ErrIncompatibleSchemaChange`) — korrigiert über `ADR-0063`;
  Folge: `run-integration-tests.sh` musste strukturell umgebaut werden
  (Drop-Column-Test wanderte hinter die Container-Ende-Grenze).
- `ADR-0058`s ursprünglicher Upgrade-Mechanismus (zweiter
  `make schema-rollout`-Lauf) war strukturell nicht lauffähig — real mit
  und ohne `--execute` reproduziert (Exit 8, vier statt der angenommenen
  zwei Fremdobjekte) — korrigiert über `ADR-0064`. Der zugrunde liegende
  strukturelle Befund (`BEO-PGC/schema-rollout-fremdobjekte`) bleibt
  offen (siehe §Trigger-Audit unten).
- Ein Planner-Koordinator-Fehler während des `slice-063`-Blocker-Reports
  (roter, korrekt ungepipter Exit-Code sichtbar, `git push` aber im
  selben Arbeitsschritt-Batch beauftragt statt in einem eigenen,
  auswertenden Schritt) war das vierte Auftreten von
  `BEO-PGC/report-nackte-id-ohne-link` — verkörpert **vor** dieser
  Welle-Closure, als eigener, wellen-interner Lese-Schritt bei
  `slice-063` (`AGENTS.md` §3.9, neuer Absatz „Prüfung und Folgehandlung
  sind zwei Schritte, nicht einer"; Herkunfts-Anker `seit slice-063`).
- `BEO-PGC/github-actions-unverifizierbar-lokal` erreichte mit
  `slice-064` real 3× — der Lese-Schritt dafür lief bereits als eigener,
  separater Architect-Zug (nicht als Teil dieser Planner-Closure, siehe
  §Steering-Loop-Einträge).
- `slice-065`s Plan zitierte den Registerstand von
  `BEO-PGC/github-actions-unverifizierbar-lokal` bei Eröffnung noch mit
  2×, obwohl `slice-064`s Closure ihn bereits auf 3× gehoben hatte — vom
  Reviewer als INFO gefunden, vom Verifier bestätigt. Kein
  Implementierungsfehler, aber ein Hinweis: Slice-Pläne zitieren den
  Registerstand zum Planungszeitpunkt, der zwischen Eröffnung und
  Closure mehrerer Slices derselben Welle weiterlaufen kann.

## Steering-Loop-Einträge

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 (hier stehen **nur** Beobachtungen, die im
Register 3× erreicht haben; jeder Eintrag nennt seine `BEO-<NNN>`) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (Feld
und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der Backticks; die
**Spec-Lücke** trägt statt `liegt in` ihre `LH-*`-ID — das ist kein Versehen).

- **Neue Hard Rule, dieser Lese-Schritt als echter, separater
  Architect-Zug — nicht Teil dieser Planner-Closure:**
  `BEO-PGC/github-actions-unverifizierbar-lokal` erreichte mit
  `slice-064` real 3×. Verdikt: **verkörpert** → `AGENTS.md` §3.10 „Ein
  neuer oder strukturell geänderter GitHub-Actions-Workflow gilt erst
  nach einem realen, grünen Post-Push-Lauf als abgeschlossen" — liegt in
  `AGENTS.md §3.10`. Kein neuer Sensor (strukturell nicht möglich:
  externer, gehosteter Dienst, Docker-only/netzloser Geltungsbereich
  schließt den Gegenstand aus, `AGENTS.md` §3.1) — der Architect-Verdikt
  dazu, dass sich GitHub-Actions-Workflows nicht lokal verifizieren lassen
  (3×).
  Auslöser: `BEO-PGC/github-actions-unverifizierbar-lokal`
  (`slice-039`, `slice-056`, `slice-064` — 3×).
- **Bereits vor dieser Welle-Closure verkörpert, nur bestätigend
  erwähnt** (Feststellung, kein neuer Eintrag dieser Closure):
  `BEO-PGC/report-nackte-id-ohne-link` erreichte mit dem
  `slice-063`-Blocker-Report real 4× und wurde noch während der Welle,
  als eigener Lese-Schritt bei `slice-063` selbst, verkörpert (`AGENTS.md`
  §3.9, neuer Absatz „Prüfung und Folgehandlung sind zwei Schritte, nicht
  einer") — liegt in `AGENTS.md §3.9`. Auslöser: `BEO-PGC/report-nackte-id-ohne-link`
  (`slice-054`, `slice-055`, `slice-056`, `slice-063-blocker` — 4×,
  Herkunfts-Anker `seit slice-063`).

## Beobachtungs-Register (Zeiger)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register — der Zähler wird **nicht** hier gepflegt; diese
Sektion ist ein Zeiger und trägt keine Daten.

Der Bestand liegt in [`../observations/`](../observations/)`BEO-PGC/`. In
dieser Welle **neu angelegt**: `kein-admin-weg-schema-fehler-recovery`
(1×, `slice-062`, Ausgang *weiter offen*, unter der Schwelle),
`kein-echter-versionswechsel-upgrade-test` (1×, `slice-063`, Ausgang
*weiter offen*, unter der Schwelle). In dieser Welle **fortgeschrieben,
weiter unter der Schwelle**: `schema-rollout-fremdobjekte` (2×,
`slice-016`, `slice-063` — kein Ausgang fällig; `ADR-0064`
Re-Evaluierungs-Trigger 2 bleibt unverändert unerreicht, siehe
§Trigger-Audit unten). In dieser Welle **auf 3× gehoben und mit Ausgang
versehen** (als separater Architect-Zug, siehe Steering-Loop-Einträge):
`github-actions-unverifizierbar-lokal` (3×, Ausgang *verkörpert* →
`AGENTS.md` §3.10). **Bereits innerhalb der Welle, vor dieser
Closure, auf 4× gehoben und mit Ausgang versehen**:
`report-nackte-id-ohne-link` (4×, Ausgang *verkörpert* → `AGENTS.md`
§3.9, Herkunfts-Anker `seit slice-063`). Durchgesehen ohne Zähler-Änderung
in mehreren Slices dieser Welle (§8-Sichtungen, keine Treffer oder unter
der Schwelle unverändert): `test-integration-retention-timing-flake`
(1×), `test-runner-stiller-ausschluss` (1×),
`test-isolation-geteilter-zustand`, `schema-rollout-braucht-compose-init`,
`nicht-blockierender-workflow-alarmmuedigkeit`,
`dod-checkbox-nachzug`/`-architect-pfad`/`-review-ohne-fixrunde`,
`architect-verdikt-ablageort-uneinheitlich`,
`architect-verdikt-rollen-scope-luecke`,
`commit-traceability-kein-vorab-hook`,
`handbuch-versionshistorie-uebersprungen`, `d-migrate-nacharbeit`,
`schema-evolution-nicht-dynamisch`, `spec008-replication-luecke`,
`slice-chronik-in-code-kommentar` — kein Treffer in dieser Welle.
Unverändert, nicht von dieser Welle berührt: alle übrigen Registereinträge
aus `welle-16` und früher.

## Folge-Slices

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — **derivativ**: Diese Liste zeigt nur,
das Original ist die Slice-Datei. Jeder genannte Folge-Slice muss als Datei im
Planning-Lifecycle existieren; genannt ohne angelegt ist dieselbe Klasse wie
ein halluziniertes Gate.

Keine, die diese Welle selbst fortsetzen — `welle-17` schließt vollständig
mit ihren vier Slices (`slice-062`, `slice-063`, `slice-064`,
`slice-065`); alle fünf ursprünglich testfreien Lastenheft-Kennungen sind
belegt. `slice-073` (lokaler `commit-msg`-Git-Hook) ist ein bereits vor
dieser Welle geplanter, wellenloser Folge-Slice aus der `welle-16`-Closure,
kein Folge-Slice dieser Welle.

## Verifikation

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 1 — keine Behauptung ohne nachprüfbaren
Anker (Hash, Lauf, Zahl).

- Alle vier Slices (`slice-062`, `slice-063`, `slice-064`, `slice-065`)
  liegen in `done/`.
- `make gates` grün — Planner-Lauf zur Closure, Exit-Code explizit
  geprüft, nicht durch eine Pipe maskiert (`AGENTS.md` §3.9): Coverage-Gate
  44,40 % ≥ 35 % Schwelle, `d-check` 0 Befunde (516 Dateien geprüft),
  `commit-traceability` OK, `a-check` 0 Befunde.
- **Welle-17-Closure-Trigger (c)** — der reale, grüne
  `e2e.yml`-Matrix-Lauf über beide PostgreSQL-Legs: Dieser Beleg liegt
  bereits **dreifach unabhängig bestätigt** vor, aus `slice-064`s Lauf,
  und wird hier zusammengefasst statt erneut erzeugt.
  - **Implementer-Behauptung:** Run `34822131377` — beide Matrix-Legs
    (PostgreSQL 17, PostgreSQL 18) `completed`/`success`.
  - **Planner-Koordinator:** per `gh run view` unabhängig bestätigt
    (siehe `slice-064` §6 Risiken).
  - **Verifier** (der Verifikationsbericht zu `slice-064`): eigener,
    unabhängiger `gh run view`/`gh run list`-Abruf desselben Laufs,
    identisches Ergebnis.
  - Kein neuer Testlauf in dieser Closure nötig — der vorliegende,
    dreifach unabhängig reproduzierte Beleg trägt den Closure-Trigger.
    Ergänzend real belegt: `slice-065`s eigener `ci.yml`-Lauf
    (Run `34823976027`, Commit `5d7672b`, `completed`/`success`,
    Planner **und** Verifier unabhängig bestätigt) — kein Bestandteil
    des Welle-Closure-Triggers selbst (der bezieht sich ausschließlich
    auf `e2e.yml`), aber Teil der DoD-vollständigen Menge dieser Welle.
- **Trigger-Audit der Welle** (Modul 6 §Wellen-Closure-Prozedur, Schritt 2
  — drei Artefaktklassen):
  - **Carveouts:** `docs/plan/carveouts/` enthält ausschließlich
    `.gitkeep` — kein offener Carveout im Repo.
  - **Bootstrap-aware Gates:** keines Teil dieser Welle.
  - **ADRs mit Re-Evaluierungs-Trigger:** `ADR-0058` (fünf
    Teilentscheidungen), `ADR-0063`, `ADR-0064` geprüft. `ADR-0063`s
    Trigger ist *permanent* (folgt aus bereits bestehendem,
    unveränderlichem Code-Verhalten). `ADR-0058`s Teilentscheidungen 1
    und 5 sind *permanent*; Teilentscheidung 2 (`LH-FA-DAT-006`) wartet
    auf eine noch nicht implementierte reale Metadatenerweiterung von
    `cdc.change` — nicht eingetreten; Teilentscheidung 3 ist bereits
    durch `ADR-0064` korrigiert (kein eigener Trigger mehr relevant);
    Teilentscheidung 4 (`SPEC-012`-Versionsliste) — unverändert, nicht
    eingetreten. `ADR-0064`s Trigger 1 (echte Release-Historie) — nicht
    eingetreten (`slice-040` weiterhin nicht angelegt). **`ADR-0064`s
    Trigger 2 (`BEO-PGC/schema-rollout-fremdobjekte` aufgelöst) ist
    explizit geprüft und NICHT eingetreten** — die Beobachtung steht
    weiterhin bei 2× (`evidence/slice-016.md`, `evidence/slice-063.md`),
    unter der 3×-Registerschwelle; kein Trigger fällig. **Bewusste
    Feststellung: kein ADR-Trigger dieser Welle ist fällig, keine
    Nacharbeit nötig.**
- Drei Paarungen (Anker · Folge-Slice · Register) — geprüft **nach** dem
  `git mv` dieser Welle-Plan-Datei nach `done/`:
  - **Anker** — der eine Steering-Loop-Eintrag mit `liegt in`-artigem
    Verweis (`AGENTS.md §3.10`) existiert real; der zweite
    (`AGENTS.md §3.9`, bereits vor dieser Closure verkörpert) existiert
    ebenfalls real.
  - **Folge-Slice** — keiner dieser Welle genannt (siehe §Folge-Slices);
    nichts zu prüfen.
  - **Register** — alle in dieser Welle berührten Verzeichnisse
    (`kein-admin-weg-schema-fehler-recovery`,
    `kein-echter-versionswechsel-upgrade-test`,
    `schema-rollout-fremdobjekte`, `github-actions-unverifizierbar-lokal`,
    `report-nackte-id-ohne-link`) existieren mit nicht leerem
    `evidence/` — grün.

## Archivierung

Feststellung: das Repo führt weiterhin **kein Archivierungs-Werkzeug**
(`archiv.zip`-Target existiert nicht) — die Archivierungs-Bedingung ist in
diesem Zug **nicht eingetreten**; alle vier Slice-Dateien, ihre Review-/
Verifier-Reports sowie dieser Welle-Plan bleiben vollständig in `done/`.
