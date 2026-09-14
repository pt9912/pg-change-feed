# Slice slice-064: PostgreSQL-Versionsmatrix im E2E-Workflow

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** welle-17 — technisch unabhängig von `slice-062`/`slice-063`
implementierbar (parametrisiert nur `compose.yaml`/`Makefile`/`e2e.yml`);
der Welle-Closure-Trigger (grüner `e2e.yml`-Matrix-Lauf über beide Legs)
braucht `slice-062`/`slice-063` zusätzlich, damit jedes Matrix-Leg alle
fünf Testbelege trägt.

**Bezug:** [LH-QA-POR-001](../../../../spec/lastenheft.md),
[ADR-0058](../../adr/0058-testansatz-fuenf-luecken.md) (Entscheidung 4 —
Matrix-Job, Parametrisierung, Betroffene Dateien; vorab entschieden),
[ADR-0051](../../adr/0051-cicd-pipeline-github-actions.md)
(`ci.yml`/`e2e.yml`-Rollenteilung, nicht-blockierend).

**Berührte Spec-Stellen:**
[SPEC-012](../../../../spec/pflichtenheft.md) (`PG_MAJOR_VERSIONS = 17, 18`
— bereits festgelegt, dieser Slice liefert den fehlenden E2E-Beleg für
beide Versionen, ändert die Liste nicht).

**Verantwortlich:** pt9912.

**Autor:** pt9912 (Planner-Lauf). **Datum:** 2026-09-14.

---

## 1. Ziel und Abgrenzung

**Ziel:** `compose.yaml`s `postgres`-Service-Image-Zeile wird von einem
literalen Digest auf eine `${PG_TEST_IMAGE}`-Interpolation umgestellt
(Default bleibt der PostgreSQL-18-Digest); `.github/workflows/e2e.yml`
bekommt eine `strategy: matrix:` über die beiden in `SPEC-012` festgelegten
Digests (PostgreSQL 17, PostgreSQL 18), jedes Leg exportiert
`PG_TEST_IMAGE` vor `make image`/`make test-integration`. Der konkrete
PostgreSQL-17-Digest wird real über `docker manifest inspect
postgres:17-alpine` (amd64) ermittelt und mit Tag-Kommentar gepinnt
(`AGENTS.md` §3.8).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **`Makefile`s `PG_TEST_IMAGE`-Default für lokale/manuelle Läufe** —
  bleibt unverändert auf dem PostgreSQL-18-Digest (`ADR-0058` Entscheidung
  4); nur die CI-Matrix und die Compose-Interpolation sind Gegenstand
  dieses Slice.
- **Linux-Plattform-Assertion (`LH-QA-POR-002`)** — `slice-065`; eigener
  Workflow (`ci.yml` statt `e2e.yml`), eigene Messmethode, keine
  Abhängigkeit von der Versionsmatrix.
- **`LH-FA-SCH-003`/`LH-FA-DAT-006`/`LH-QA-OPS-005` (E2E-Testfunktionen und
  Upgrade-Phase)** — `slice-062`/`slice-063`; andere Dateien
  (`integration_test.go`, `run-integration-tests.sh`), andere
  Eigenschaftsklasse (Verhalten/Struktur/Betriebsmechanik statt
  Umgebungsmatrix). Der Matrix-Job dieses Slice führt sie im selben Lauf
  mit aus, ändert sie aber nicht.
- **Änderung der `SPEC-012`-Versionsliste selbst** (z. B. Aufnahme von
  PostgreSQL 19) — `ADR-0058` Re-Evaluierungs-Trigger 4: eine künftige
  Erweiterung der Liste ist eine eigene Pflichtenheft-Änderung mit
  eigenem Slice, nicht Teil dieser technischen Umsetzung der bereits
  festgelegten zwei Versionen.

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

- [x] `LH-QA-POR-001` erfüllt: `.github/workflows/e2e.yml` trägt eine
      `strategy: matrix:` über die PostgreSQL-17- und PostgreSQL-18-Digests;
      jedes Leg exportiert `PG_TEST_IMAGE`, ruft dieselben bestehenden
      Targets (`make image`, `make test-integration`) unverändert auf.
- [x] `compose.yaml`s `postgres`-Service-Image-Zeile nutzt
      `${PG_TEST_IMAGE}`-Interpolation mit dem PostgreSQL-18-Digest als
      Default; ein lokaler Lauf ohne gesetzte Variable bleibt unverändert
      lauffähig (Regressionstest).
- [x] Der PostgreSQL-17-Digest ist real über `docker manifest inspect
      postgres:17-alpine` (amd64) ermittelt, mit Tag-Kommentar gepinnt
      (`AGENTS.md` §3.8) und in der Matrix-Zeile referenziert.
- [x] `Makefile`s `PG_TEST_IMAGE`-Default bleibt unverändert (siehe §1
      Out-of-Scope) — Regressionstest: `make test-integration` ohne
      gesetzte Variable nutzt weiterhin den PostgreSQL-18-Digest.
- [x] `make gates` grün.
- [x] `make test-integration` lokal mit explizit exportiertem
      PostgreSQL-17-`PG_TEST_IMAGE` mindestens einmal grün belegt (lokaler
      Vorab-Beleg vor dem ersten echten CI-Matrix-Lauf; kein Gate,
      [ADR-0030](../../adr/0030-testpyramide.md)).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      Siehe `docs/reviews/review-slice-064.md` (0 HIGH, 0 MEDIUM, 1 LOW,
      keine Fixrunde nötig — DoD-Checkbox-Nachzug ohne Fixrunde,
      `.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne Fixrunde).
- [x] Doku-Update: `harness/README.md` §Werkzeuge, `.github/workflows/e2e.yml`-Zeile
      um die Matrix-Beschreibung ergänzt (kein neues Gate, `e2e.yml` bleibt
      nicht-blockierend).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben,
      **falls dieser Slice einen Inventur-Fund auflöst** — entfällt: Repo
      ist GF (`harness/conventions.md` Modus-Deklaration `PGC`), keine
      `reconciliation.md` vorhanden.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen
      `evidence/`; keine Beobachtung angefallen ist ebenfalls eine Antwort
      und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen /
      weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      **verschoben auf `welle-17`-Closure** (dieser Slice trägt
      `Welle: welle-17`).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `compose.yaml` | update | Image-Zeile parametrisiert (`${PG_TEST_IMAGE}`-Interpolation, Default PostgreSQL 18) |
| `Makefile` | keine funktionale Änderung | `PG_TEST_IMAGE`-Default bleibt (siehe §1 Out-of-Scope); ggf. Kommentar-Klarstellung, dass die Variable auch von der CI-Matrix genutzt wird |
| `.github/workflows/e2e.yml` | update | `strategy: matrix:` über beide `SPEC-012`-Digests, `PG_TEST_IMAGE`-Export je Leg |
| `harness/README.md` | update | `e2e.yml`-Zeile um die Matrix-Beschreibung ergänzt |

## 4. Trigger

**Start** (`next` → `in-progress`): `welle-17` eröffnet, `Verantwortlich:`
gesetzt, WIP-Limit (1 je Implementer) frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Zeigt sich, dass
  die Compose-Parametrisierung eine tiefere Änderung an
  `run-integration-tests.sh` braucht (z. B. weil ein Skript-Teil den
  Digest literal voraussetzt statt über die Variable zu gehen), gehört das
  zurück zur Zerlegung.
- `in-progress` → `open` (blockiert — Carveout?): Falls
  `docker manifest inspect postgres:17-alpine` netzseitig nicht erreichbar
  ist oder kein amd64-Manifest liefert, ist die Digest-Ermittlung
  blockiert — Carveout mit Folge-Slice, sobald Netzzugriff verfügbar ist.

## 5. Closure-Trigger

DoD vollständig **und** `make gates` grün **und** Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- Ein `docker run`/`docker manifest inspect`-Aufruf gegen die Registry
  braucht Netzzugriff, den ein Agentenlauf ggf. nicht hat — die
  Digest-Ermittlung ist damit potenziell blockierend (siehe §4
  Rückführung). — **Ausgang: entfallen** — real gelang der Zugriff
  mehrfach unabhängig (Implementer, Reviewer, Verifier), identischer
  Digest über alle drei Läufe.
- `BEO-PGC/github-actions-unverifizierbar-lokal` (3× mit diesem Slice —
  siehe §8): Der reale Matrix-Lauf (Run `34822131377`) ist bereits grün,
  beide Legs `completed`/`success`, vom Planner-Koordinator UND
  unabhängig vom Verifier per `gh run view` bestätigt. — **Ausgang:
  entfallen** — real belegt. Zusätzlich `evidence/slice-064.md` im
  Register ergänzt (3. Beleg der strukturellen Verifikationsgrenze-
  Beobachtung, Ausgang folgt bei `welle-17`-Closure).
- PostgreSQL 17 könnte real ein vom bestehenden PostgreSQL-18-Verhalten
  abweichendes `pgoutput`-/Replication-Detail zeigen. — **Ausgang:
  entfallen** — realer `make test-integration`-Lauf gegen PostgreSQL 17
  (Implementer + isolierter Verifier-Lauf, plus der reale CI-Matrix-Leg)
  lief vollständig grün, keine Abweichung beobachtet.

## 7. Closure-Notiz

- **Was hat funktioniert:** Der Mechanismus (`${PG_TEST_IMAGE}`-
  Interpolation, Matrix-Job) funktionierte im ersten Versuch real —
  beide Legs liefen beim ersten Push grün. Der Verifier fand und schloss
  eigenständig eine Lücke (DoD-Punkt „lokaler PG-17-Vorab-Lauf" war
  abgehakt, aber ohne Artefakt-Beleg im Repo) — reale, unabhängige Prüfung
  trug hier über eine reine Selbstauskunfts-Übernahme hinaus.
- **Was ging anders als geplant:** Nichts Wesentliches — Implementierung
  entsprach `ADR-0058` Entscheidung 4 nahezu wörtlich; einzige Ergänzung
  war `fail-fast: false` (Implementer-Entscheidung, von Reviewer und
  Verifier unabhängig als legitim bestätigt).
- **Steering-Loop-Eintrag:** keiner neu verkörpert — `BEO-PGC/github-actions-unverifizierbar-lokal`
  erreicht mit diesem Slice zwar 3×, aber der Lese-Schritt (Verkörperung
  ja/nein) läuft bei der `welle-17`-Closure, nicht hier.
- **Beobachtungs-Register (`../observations/`):** `evidence/slice-064.md`
  in `BEO-PGC/github-actions-unverifizierbar-lokal/` ergänzt — Zähler
  steht bei 3×, Ausgang folgt bei `welle-17`-Closure.
- **Folge-Slices:** keine.
- **Risiken aus §6:** alle drei entfallen — siehe §6.
- **Drei Paarungen:** verschoben auf `welle-17`-Closure (dieser Slice trägt
  `Welle: welle-17`, siehe DoD-Item).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Einzige berührte Sub-Area ist die
Repo-weite Default-Sub-Area `*`/`PGC` (`harness/conventions.md` führt keine
feinere Sub-Area für CI-Workflows oder Compose-Verdrahtung).

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen
(`docs/plan/planning/observations/BEO-PGC/`), insbesondere gegen
GitHub-Actions-Matrix-Verhalten:

- `BEO-PGC/github-actions-unverifizierbar-lokal` — **offen**, 2×
  (`evidence/slice-039.md`, `evidence/slice-056.md`; unter der
  3×-Schwelle) — direkter Treffer: jeder GitHub-Actions-Workflow-Aspekt
  (hier: ob beide Matrix-Legs real parallel laufen und grün werden) löst
  sich strukturell erst mit dem ersten echten Lauf auf GitHub nach einem
  Push/PR/Merge auf; kein Slice kann das repo-intern vorwegnehmen. Würde
  dieser Slice als dritte Instanz gezählt (3×), wäre die Beobachtung eine
  Lücke statt einer Notiz — als Risiko in §6 aufgenommen, damit die
  Closure-Notiz den tatsächlichen Zähler-Stand nach diesem Slice
  dokumentiert.
- `BEO-PGC/dod-checkbox-nachzug`, `dod-checkbox-nachzug-architect-pfad`,
  `dod-checkbox-nachzug-review-ohne-fixrunde` — geprüft, kein Treffer:
  betreffen Planning-Harness-Prozessdisziplin (DoD-Häkchen-Nachzug), nicht
  CI-Workflow-Inhalte.
- `BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit` — geprüft: prüft
  auf Alarmmüdigkeit durch wiederholt rote, nicht-blockierende
  `e2e.yml`-Läufe. Thematisch benachbart (dieser Slice verdoppelt die
  Runner-Minuten von `e2e.yml` und damit potenziell die Zahl roter Legs),
  aber kein direkter Treffer: Diese Beobachtung betrifft die
  Wahrnehmungsdisziplin bei bereits rotem `e2e.yml`, nicht die
  Matrix-Einführung selbst.
- Weitere durchgesehen (`schema-rollout-braucht-compose-init`,
  `schema-rollout-fremdobjekte`): keine Treffer — beide betreffen
  d-migrate-Rolloutverhalten, nicht die Compose-Image-Parametrisierung
  oder die CI-Matrix-Struktur.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF (nur
`*`/`PGC`).

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
