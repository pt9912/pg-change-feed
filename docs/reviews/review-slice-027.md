# Review-Report: slice-027 — 2026-09-12

**Review-Art:** Code — geprüft gegen Plan (`slice-027` §1–§3) und
Konventionen (`AGENTS.md` §3 Hard Rules, `.harness/skills/reviewer.md`).
Maintainability-Fokus (Modul 10 §Kontext-Zuschnitt) — keine DoD-/
Spec-Konformitätsprüfung (Verifier-Aufgabe, Modul 11); wo unten dennoch
DoD-Häkchen betrachtet werden, geschieht das ausschließlich zur Prüfung der
**Häkchen-Ehrlichkeit** (Beleg vorhanden ja/nein), nicht zur DoD-Abnahme
selbst.

**Gegenstand:** Commits `6daa232` (test: Black-Box-CLI-Rundlauf über
`docker exec` gegen den Feed-Container), `726e911` (docs/planning:
slice-027 DoD-Nachzug)

**Skill:** `.harness/skills/reviewer.md` @ `e9159ff` (HEAD zum Review-Zeitpunkt)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-12

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- `docs/plan/planning/in-progress/slice-027-black-box-e2e-test-cli.md`
  §1–§8 (inkl. Plan-Nachzug aus `726e911`)
- `docs/plan/planning/welle-8.md` §1 (Begründung: kein bestehender
  `*_endtoend_test.go` ist Black-Box, das Repo hat keinen echten E2E-Test),
  §3 (Closure-Trigger der Welle)
- `docs/plan/adr/0030-testpyramide.md` (Accepted — Integrationstest- vs.
  E2E-Tier-Unterscheidung)
- `AGENTS.md` §3 (Hard Rules, insbesondere §3.1, §3.7)
- `harness/conventions.md` (MR-000/MR-001)
- `harness/README.md` (Sensors-Tabelle, `make test-integration`-Zeile)
- Code als Review-Gegenstand: `tools/harness/run-integration-tests.sh`
  (Diff, neuer Abschnitt „Black-Box-CLI-Rundlauf"), `harness/README.md`
  (Diff)
- Code als Sachverhaltsprüfung (unverändert, gegengelesen):
  `cmd/pg-change-feed/main.go` (CLI-Unterbefehl-Signaturen
  `register-consumer <name>` / `acknowledge-consumer <consumer-id>
  <position>`), `compose.yaml` (Healthcheck-Vertrag, `restart: "no"`, kein
  `build:`-Block), `Dockerfile` (distroless-Runtime, keine Shell),
  `internal/application/usecase/enable/service.go` (`AlreadyEnabled` —
  Idempotenz der Aktivierung bei Neustart, unverändert)
- Eigene Gate-/Test-Läufe in dieser Sitzung: `make gates` (grün, alle vier
  inneren Gates), `make test-integration` (grün, Baseline-Lauf); eigene
  **Rot-Grün-Gegenprobe** an der zentralen Grenz-Assertion (`>=` statt `>`
  in der `resumed_ids`-Abfrage) — rot mit exakt der vom Implementer
  behaupteten Ausgabe, danach `git checkout` zurückgesetzt (`git diff`
  leer) und erneut grün bestätigt

---

## Findings

Keine HIGH- oder MEDIUM-Findings in diesem Lauf.

## Negativbefunde

- geprüft, ohne Befund: **Echtes Black-Box-Aufrufmuster** — `exec_feed()`
  (`tools/harness/run-integration-tests.sh`, neuer Abschnitt) ruft
  ausschließlich `docker exec "$FEED_CONTAINER" /pg-change-feed "$@"` auf;
  beide Aufrufstellen (`register-consumer "$CLI_CONSUMER"`,
  `acknowledge-consumer "$CLI_CONSUMER" "$first_position"` bzw.
  `"$second_position"`) laufen ausschließlich über diese Helfer-Funktion.
  Kein Go-Test-Prozess und kein internes Paket (`bootstrap.*`,
  `postgresstorage.*`) ist an der *Aktion* beteiligt; das Lesen (bewusst
  erlaubte Ausnahme laut Plan §1) läuft über `docker exec … psql`-Abfragen
  gegen `cdc.changes`/`cdc.consumer`/`cdc.consumer_position`. Die
  Signaturen stimmen mit `cmd/pg-change-feed/main.go:45-88` überein
  (`register-consumer <name>`, `acknowledge-consumer <consumer-id>
  <position>`).
- geprüft, ohne Befund: **„Kein Replay nach Neustart" ist real bewiesen,
  nicht nur behauptet** — die Positions-Prüfung nach dem Neustart liest
  `acknowledged_position` **frisch** aus `cdc.consumer_position`
  (`restored_position=$(docker exec … SELECT acknowledged_position …)`),
  statt die im Skript gehaltene `$first_position`-Variable wiederzuverwenden;
  erst danach wird `resumed_ids` mit `commit_position > $restored_position`
  gegen `cdc.changes` abgefragt und exakt gegen `"96"` verglichen. Eigene
  Rot-Grün-Gegenprobe: `commit_position >= $restored_position` (mutiert)
  liefert real `id-Folge '95,96', wollen '96'` und lässt `make
  test-integration` mit Ausgang 1 fehlschlagen — identisch mit der im
  DoD-Beleg (`726e911`) behaupteten Ausgabe. Nach Rücksetzen (`git
  checkout -- tools/harness/run-integration-tests.sh`, `git diff` leer)
  läuft der Lauf wieder grün.
- geprüft, ohne Befund: **Health-Gate nach dem simulierten Neustart** —
  nach `docker restart "$FEED_CONTAINER"` wartet das Skript in einer
  Schleife (60 × 1s) auf `docker inspect --format
  '{{.State.Health.Status}}'` == `healthy`, bevor es die zweite
  Quelländerung schreibt; dieselbe Struktur wie der bestehende
  Start-Healthcheck weiter oben im selben Skript. Der Compose-Healthcheck
  selbst läuft über den Binary-Exit-Code (`compose.yaml:91`, distroless,
  kein `sh -c`-Umweg), nicht über eine Shell-Prüfung — kein Zugriff auf
  einen noch nicht neu verdrahteten Container möglich, da der nachfolgende
  Schreib-/Lese-Zug ohnehin bis zu 30s (120 × 0,25s) auf das tatsächliche
  Erfassen der zweiten Änderung wartet.
- geprüft, ohne Befund: **Kein Ressourcen-/Zustands-Leck durch den
  Neustart** — der simulierte Neustart läuft ausschließlich über `docker
  restart "$FEED_CONTAINER"` (SIGTERM + Neustart desselben Containers,
  laut Kommentar); zwischen den beiden CLI-Bestätigungen findet **kein**
  `$COMPOSE down`/`up` statt (das passiert nur einmal am Skriptanfang und im
  `cleanup()`-Trap am Skriptende). Replication-Slot und Publication bleiben
  auf der PostgreSQL-Seite unberührt — ein `docker restart` betrifft nur den
  Feed-Container, nicht `cdc-test-postgres`.
- geprüft, ohne Befund: **Kein Chronik-Sprachgebrauch** (`AGENTS.md` §3.7)
  — alle neuen Kommentare im Skript (Skriptkopf-Ergänzung,
  Abschnittskommentar „Black-Box-CLI-Rundlauf", Neustart-Kommentar,
  Neustart-Beleg-Kommentar) beschreiben den geltenden Zustand im Indikativ
  oder eine Abgrenzungs-/Kopplungsbegründung gegenüber **bestehendem**,
  unverändertem Nachbar-Code im selben Skript („anders als der
  Go-Testlauf oben", „wie beim Lasttest-Beleg oben") — das ist ein
  Vergleich mit vorhandenem Code, keine Chronik über entfernten/verworfenen
  Text. Kein Konjunktiv über eine verworfene Alternative, kein abwesender
  Text, kein mitten im Satz abgebrochener Kommentar.
- geprüft, ohne Befund: **Docker-only / Suppression-Verbot** (`AGENTS.md`
  §3.1/§3.2) — kein lokales Toolchain-Setup im neuen Abschnitt (nur
  `docker exec`/`docker restart`/`docker inspect`), kein
  `//nolint`/`#noqa`-artiges Suppression-Muster im Diff.
- geprüft, ohne Befund: **Doku-Konsistenz `harness/README.md`** — die
  aktualisierte `make test-integration`-Zeile nennt exakt den neuen Umfang
  (`docker exec` gegen den Feed-Container, kein Go-Paket-Import, simulierter
  Neustart, Fortsetzen ab der bestätigten Position über `cdc.changes`) und
  überzeichnet ihn nicht — insbesondere wird der Lese-Zugriffsweg korrekt
  weiterhin als „bestehender SQL-Lesezugriffsweg" benannt, nicht als
  Black-Box-Lesen.
- geprüft, ohne Befund: **DoD-Checkbox-Deckung (`726e911`)** — die vier
  neu auf `[x]` gesetzten Punkte (`LH-QA-POR-003`-Umsetzung, voller
  Rundlauf-Beleg inkl. Rot-Grün-Gegenprobe, `make gates`/dreifacher
  `make test-integration`-Lauf, Doku-Update) sind durch das tatsächlich
  Gelieferte gedeckt (siehe obige Negativbefunde und eigene Gate-/Testläufe
  dieser Sitzung); die als entfällt markierte Reconciliation-Zeile ist
  korrekt begründet (Sub-Area `*`/`PGC` ist Greenfield, Datei existiert im
  Repo nicht). Review-Zeile bleibt korrekt offen (`[ ]`), ebenso
  Closure-Notiz, Beobachtungs-Register, Risiken-Ausgänge und die drei
  Paarungen — konsistent mit dem Lifecycle-Zustand `in-progress`.
- geprüft, ohne Befund: **Out-of-Scope-Disziplin §1** — kein neuer
  CLI-Lese-Unterbefehl (Bestand bleibt bewusst stehen, begründet), keine
  Rollen-DSN-Trennungs-Prüfung (Folge-Slice `slice-028` benannt und
  existiert als Plan-Datei), keine zusätzliche WAL-Rückstand-Prüfung gegen
  den Compose-Stack (bereits real gegen PostgreSQL bewiesen, `slice-026`),
  keine Umbenennung von `make test-integration`/`mvp_test.go` (an
  `slice-028` verwiesen).
- geprüft, ohne Befund: **Plan-Nachzug §3** — die drei dokumentierten
  Abweichungen vom ursprünglichen §3 (kein zusätzlicher Go-Testfall, keine
  `compose.yaml`-Änderung, `harness/README.md` ergänzt) sind im Diff exakt
  so umgesetzt, wie der Nachzug-Text sie beschreibt; keine stille
  Erweiterung über den nachgetragenen Plan hinaus.
- geprüft, ohne Befund: **Sub-Area-/Modus-Prüfung (§8 des Slice-Plans)** —
  einzige berührte Sub-Area ist die Default-`PGC` (Greenfield); beide
  zitierten Register-Treffer (`BEO-PGC/rollen-test-abdeckungsluecken`,
  `BEO-PGC/test-isolation-geteilter-zustand`) existieren im Register mit
  je einem Beleg (1×) und sind korrekt als nicht einschlägig für diesen
  Slice bewertet — keiner erreicht mit diesem Slice 3×.
- geprüft, ohne Befund: **Traceability** — Commit `6daa232` nennt
  `LH-QA-POR-003, ADR-0030`, Commit `726e911` nennt `LH-QA-POR-003`; kein
  `SPEC-*`/`ARC-*` im Betreff eines der beiden Commits. `ADR-0030` ist im
  ADR-Index (`docs/plan/adr/README.md:43`) mit Status `Accepted` geführt.
  `make commit-traceability` (Teil von `make gates`) meldet 0 Befunde über
  die letzten 5 Commits.
- geprüft, ohne Befund: **`make gates`** (in dieser Sitzung selbst
  ausgeführt) — `baseline-verify` (54 Dateien OK), `docs-check`
  (245 Dateien, 0 Befunde, inkl. Commit-Range-Lauf), `commit-traceability`
  (5 Commits, Betreffs ohne Struktur-ID), `a-check` (0 Befunde) — alle vier
  inneren Gates grün.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** keine (kein Finding oberhalb der
Negativbefund-Schwelle)

## Verdikt

**Merge-blockierend:** nein — kein HIGH, kein MEDIUM. Der neue
Black-Box-CLI-Rundlauf ruft `register-consumer`/`acknowledge-consumer`
nachweislich ausschließlich als externen Prozess auf, und die zentrale
„kein Replay nach Neustart"-Assertion hält der eigenen
Rot-Grün-Gegenprobe stand (mutierte Grenze `>=` reproduziert exakt die im
DoD-Beleg behauptete rote Ausgabe `'95,96'` statt `'96'`).

**Übergabe:** Da dieser Lauf keine Findings oberhalb der
Negativbefund-Schwelle trägt, geht nichts an den Implementer zurück und es
entsteht kein neuer Zähler-Beitrag für das Beobachtungs-Register. Dieser
Report ist ein **Lauf-Beleg** (Audit: dieser Diff, dieser Skill, dieses
Modell, dieses Verdikt) — er wird über Läufe hinweg nicht wieder gelesen,
und muss es nicht. Der Report ersetzt keine Verifikation — DoD-/
Spec-Konformität (inkl. der noch offenen DoD-Punkte: Review-Zeile selbst,
Closure-Notiz, Beobachtungs-Register, Risiken-Ausgänge, drei Paarungen)
prüft der Verifier separat.
