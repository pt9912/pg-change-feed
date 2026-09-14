# Slice slice-063: Upgrade-Sicherheit — simulierter Container-Stopp/Rollout/Start-Zyklus

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** welle-17 — unabhängig von `slice-062`/`slice-064`/`slice-065`
implementierbar; der Welle-Closure-Trigger (grüner `e2e.yml`-Matrix-Lauf)
braucht diesen Slice zusammen mit `slice-062` in jedem Matrix-Leg.

**Bezug:** [LH-QA-OPS-005](../../../../spec/lastenheft.md),
[ADR-0064](../../adr/0064-lh-qa-ops-005-testansatz-korrektur.md)
(Supersedes `ADR-0058` Entscheidung 3 — korrigierter Mechanismus:
Container-Tausch über `$COMPOSE up -d --force-recreate --no-deps
pg-change-feed` statt eines zweiten `make schema-rollout`-Laufs;
maßgeblich für diesen Slice), [ADR-0058](../../adr/0058-testansatz-fuenf-luecken.md)
(Entscheidung 3 — ursprüngliche, real blockierte Fassung, nur noch als
Kontext), [ADR-0030](../../adr/0030-testpyramide.md) (E2E-Tier-Definition).

**Berührte Spec-Stellen:** [`LH-QA-OPS-005`](../../../../spec/lastenheft.md)
§Upgrade-Sicherheit — bereits im Lastenheft festgelegt, dieser Slice liefert
den fehlenden Testbeleg, ändert die Zusage nicht.

**Verantwortlich:** pt9912.

**Autor:** pt9912 (Planner-Lauf). **Datum:** 2026-09-14.

---

## 1. Ziel und Abgrenzung

**Ziel:** Eine neue Phase in `tools/harness/run-integration-tests.sh`, nach
dem bestehenden Black-Box-CLI-Rundlauf und vor
`TestE2ESchemaChangeIncompatibleTypeChange`, bildet den Mechanismus eines
Anwendungs-Upgrades nach — korrigierte Form nach `ADR-0064` (Supersedes
`ADR-0058` Entscheidung 3, ein zweiter `make schema-rollout`-Lauf blockiert
real mit Exit 8 auf Fremdobjekten und entfällt ersatzlos): Zeile einfügen
und Position über `cdc.changes` festhalten → `$COMPOSE up -d
--force-recreate --no-deps pg-change-feed` (realer Container-Tausch —
neue Instanz desselben `:dev`-Images, `container_name` bleibt stabil,
`postgres`/`nats` bleiben unberührt) → Health-Poll → Datenstand vor dem
Tausch identisch lesbar, eine danach eingefügte Zeile wird weiterhin
erfasst.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Echter Zwei-Image-Vergleich (altes vs. neues Image)** — `ADR-0058`
  Entscheidung 3/Re-Evaluierungs-Trigger 3 verschiebt das ausdrücklich auf
  den Zeitpunkt, sobald eine echte Release-Historie existiert
  ([ADR-0051](../../adr/0051-cicd-pipeline-github-actions.md) Folgepflicht:
  `slice-039` liegt bereits in `done/`,
  [`docs/plan/planning/done/slice-039-ci-workflow-dependabot.md`](../done/slice-039-ci-workflow-dependabot.md);
  `slice-040`, Release-Pipeline/Tags, ist noch nicht angelegt) — ohne
  Git-Tags/Releases wäre jeder frühere Commit als „Vorgängerversion"
  willkürlich gewählt.
- **`LH-FA-SCH-003`/`LH-FA-DAT-006` (E2E-Testfunktionen)** — `slice-062`;
  andere Eigenschaftsklasse (Verhalten/Struktur statt Betriebsmechanik) und
  andere Datei (`integration_test.go` statt `run-integration-tests.sh`-Phase).
- **PostgreSQL-Versionsmatrix (`LH-QA-POR-001`) und
  Linux-Plattform-Assertion (`LH-QA-POR-002`)** — `slice-064`/`slice-065`;
  beide berühren ausschließlich CI-Workflow-/Compose-Dateien.
- **Ein neues Make-Target für den Upgrade-Zyklus** — der Zyklus nutzt
  ausschließlich bestehende Bausteine (`docker stop`/`start`, `make
  schema-rollout`, `cdc.changes`-Lesung) innerhalb der bestehenden
  Orchestrierung; ein eigenes Target wäre ein zweiter Mechanismus für
  denselben Zweck.

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

- [x] `LH-QA-OPS-005` erfüllt: neue Phase in
      `tools/harness/run-integration-tests.sh` (nach dem
      Black-Box-CLI-Rundlauf, vor `TestE2ESchemaChangeIncompatibleTypeChange`)
      führt real `$COMPOSE up -d --force-recreate --no-deps
      pg-change-feed` aus (`ADR-0064`) und belegt den vor dem Tausch
      erfassten Datenstand über `cdc.changes` identisch lesbar; eine danach
      eingefügte Zeile wird weiterhin erfasst.
- [x] Real bestätigt: `--force-recreate` erzeugt eine neue Container-
      Instanz desselben `:dev`-Images (`container_name` bleibt
      `cdc-test-feed`), `postgres`/`nats` bleiben durch `--no-deps`
      unberührt und healthy.
- [x] Health-Poll nach dem Tausch analog zum bestehenden simulierten
      `docker restart`-Rundlauf (`LH-QA-REL-001`).
- [x] `make gates` grün.
- [x] `make test-integration` grün mit der neuen Phase sichtbar im Log
      (kein Gate, [ADR-0030](../../adr/0030-testpyramide.md)).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      Report: [`docs/reviews/review-slice-063.md`](../../../reviews/review-slice-063.md)
      — 0 HIGH/MEDIUM/LOW, keine Fixrunde.
- [x] Doku-Update: `harness/README.md` §Sensors/§Werkzeuge, `make
      test-integration`-Zeile um die neue Upgrade-Sicherheits-Phase ergänzt
      (kein neues Gate, kein neues Target).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben,
      **falls dieser Slice einen Inventur-Fund auflöst** — entfällt: Repo
      ist GF (`harness/conventions.md` Modus-Deklaration `PGC`), keine
      `reconciliation.md` vorhanden.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen
      `evidence/`; keine Beobachtung angefallen ist ebenfalls eine Antwort
      und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen /
      weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      **verschoben auf `welle-17`-Closure** (dieser Slice trägt
      `Welle: welle-17`).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/run-integration-tests.sh` | update | neue Phase (Container-Tausch über `$COMPOSE up -d --force-recreate --no-deps`, `ADR-0064`), nutzt ausschließlich bestehende Bausteine |
| `harness/README.md` | update | `make test-integration`-Sensor-Zeile um die neue Phase ergänzt |

## 4. Trigger

**Start** (`next` → `in-progress`): `welle-17` eröffnet, `Verantwortlich:`
gesetzt, WIP-Limit (1 je Implementer) frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Zeigt sich, dass
  `docker stop`/`start` in Kombination mit `make schema-rollout` einen
  bislang unbekannten Seiteneffekt auf laufende Replication-Slots oder
  Consumer-Positionen hat, der eine eigene Fehlerbehandlung braucht, gehört
  das zurück zur Zerlegung.
- `in-progress` → `open` (blockiert — Carveout?): Wenn der zweite `make
  schema-rollout`-Lauf real mit Exit 8 auf `cdc.heartbeat`/`cdc.metrics`
  blockiert (`BEO-PGC/schema-rollout-fremdobjekte`, §6) und keine
  unmittelbare, nicht-destruktive Umgehung im Slice-Umfang liegt, ist das
  ein Blocker außerhalb dieses Slice — Carveout oder Architect-Rückfrage,
  ob der zweite Rollout-Lauf in dieser Form (`ADR-0058` Entscheidung 3)
  tragfähig bleibt oder eine Folge-ADR braucht.

  **Nachtrag — was tatsächlich eintrat (2026-09-14):** Der Implementer
  reproduzierte den Fall real (Exit 8, sogar auf vier statt zwei
  Objekten — zusätzlich `cdc.disable_table`/`cdc.enable_table`,
  `ADR-0050`) und prüfte drei Umgehungen (`--allow-destructive`,
  Objekt-Ausschluss-Flag, Überführung ins deklarative Schema) — alle
  verworfen, keine im Slice-Umfang liegende, nicht-destruktive Lösung.
  Befund: `docs/reviews/blocker-slice-063.md`. Rückführung nach `open`
  ausgeführt; ein Architect-Zug entscheidet über Carveout vs. Folge-ADR
  zu `ADR-0058` Entscheidung 3.

## 5. Closure-Trigger

DoD vollständig **und** `make gates` grün **und** Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- **`BEO-PGC/schema-rollout-fremdobjekte` — bereits eingetreten
  (2×, real reproduziert mit vier statt zwei Objekten,
  `docs/reviews/blocker-slice-063.md`):** Der ursprünglich geplante zweite
  `make schema-rollout`-Lauf blockierte real mit Exit 8 — dieser Slice
  vermeidet das jetzt durch `ADR-0064`s korrigierten Mechanismus
  (Container-Tausch statt Migrationsschritt) vollständig, statt den
  Blocker zu umgehen. Kein weiteres Risiko für die aktuelle Umsetzung,
  aber `BEO-PGC/schema-rollout-fremdobjekte` bleibt als eigenständige,
  ungelöste strukturelle Lücke bestehen (`ADR-0064` Re-Evaluierungs-
  Trigger 2). — **Ausgang:** eingetreten → `ADR-0064` (Folge-Entscheidung
  mit korrigiertem Mechanismus, kein neuer Slice nötig, da derselbe
  `slice-063` die korrigierte Form umsetzt).
- Der Container-Stopp könnte den Replication-Slot oder eine offene
  Transaktion in einem Zustand hinterlassen, der den nachfolgenden Start
  real verzögert oder einen Health-Check-Timeout auslöst (anders als beim
  bestehenden `docker restart`-Rundlauf, der denselben Prozess ohne
  Zwischenschritt neu startet). — **Ausgang:** <bei Closure zu füllen>
- `BEO-PGC/test-integration-retention-timing-flake` (offen, 1× — siehe
  §8): eine zusätzliche Phase mit Stopp/Rollout/Start-Timing im selben
  Compose-Lauf könnte bestehende Alters-/Lag-Schwellen-Wartephasen
  verschieben. — **Ausgang:** <bei Closure zu füllen>
- Kein echter Versionswechsel (`ADR-0058` benennt das offen): „alt" und
  „neu" sind dasselbe Image — die Lücke bleibt bestehen, bis eine echte
  Release-Historie existiert (Re-Evaluierungs-Trigger 3 der ADR). —
  **Ausgang:** <bei Closure zu füllen>

## 7. Closure-Notiz

- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Steering-Loop-Eintrag:** <…>
- **Beobachtungs-Register (`../observations/`):** <…>
- **Folge-Slices:** keine.
- **Risiken aus §6:** <jedes mit genau einem Ausgang — siehe §6>
- **Drei Paarungen:** verschoben auf `welle-17`-Closure (dieser Slice trägt
  `Welle: welle-17`, siehe DoD-Item).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Einzige berührte Sub-Area ist die
Repo-weite Default-Sub-Area `*`/`PGC` (`harness/conventions.md` führt keine
feinere Sub-Area für Orchestrierungs-Skripte).

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen
(`docs/plan/planning/observations/BEO-PGC/`), insbesondere gegen neue
Testphasen im selben Compose-Lauf:

- `BEO-PGC/test-integration-retention-timing-flake` — **offen**, 1×
  (`evidence/slice-057.md`; unter der 3×-Schwelle). Thematisch relevant:
  dieser Slice fügt eine weitere zeitkritische Phase (Stopp/Rollout/Start
  mit Health-Poll) in denselben Compose-Lauf ein, in dem die bekannte
  Timing-Flake-Historie liegt — als Risiko in §6 aufgenommen, nicht als
  eigener Treffer gezählt (kein zweites Auftreten *dieser* Beobachtung,
  solange kein realer Flake in diesem Slice auftritt).
- `BEO-PGC/schema-rollout-braucht-compose-init` — geprüft: betrifft die
  initiale Compose-Hochfahr-Reihenfolge (Rollout vor erstem Feed-Start),
  nicht einen wiederholten Rollout gegen eine bereits migrierte,
  laufende DB — kein direkter Treffer, aber die dort dokumentierte
  Abhängigkeit (Rollout braucht eine erreichbare Postgres-Instanz) gilt
  unverändert auch für den wiederholten Lauf dieses Slice.
- **`BEO-PGC/schema-rollout-fremdobjekte` — direkter Treffer, tatsächlich
  eingetreten (2×, real reproduziert mit vier statt zwei Objekten,
  `docs/reviews/blocker-slice-063.md`; Beleg wird bei der Closure dieses
  Slice nachgetragen):** Der ursprünglich geplante zweite
  `make schema-rollout`-Lauf blockierte real mit Exit 8, sogar mit
  `--dry-run` (Architect-Zug, `ADR-0064` §Kontext). Dieser Slice umgeht
  das nicht, sondern nutzt jetzt `ADR-0064`s korrigierten Mechanismus
  (Container-Tausch statt Migrationsschritt), der den blockierten Lauf
  gar nicht mehr braucht.
- Weitere durchgesehen (`d-migrate-nacharbeit`,
  `test-runner-stiller-ausschluss`): kein direkter Treffer —
  `d-migrate-nacharbeit` betrifft eine inzwischen zurückgebaute
  Ausweichform (`nacharbeit-views.sql`, seit `slice-016` entfallen),
  `test-runner-stiller-ausschluss` betrifft `go test -run`-Musterzeilen für
  Go-Testfunktionen; dieser Slice fügt keine neue Go-Testfunktion hinzu,
  sondern eine reine Shell-Phase außerhalb des `-run`-Filters.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF (nur
`*`/`PGC`).

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
