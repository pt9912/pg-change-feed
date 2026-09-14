# Slice slice-065: Linux-Plattform-Assertion im Gate-Workflow

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** welle-17 — vollständig unabhängig von `slice-062`/`slice-063`/
`slice-064` (eigener Workflow, eigene Messmethode); trägt keinen eigenen
Beitrag zum Welle-Closure-Trigger (der bezieht sich ausschließlich auf den
`e2e.yml`-Matrix-Lauf), ist aber Teil der DoD-vollständigen Menge, die
`welle-17` vor Closure einsammelt.

**Bezug:** [LH-QA-POR-002](../../../../spec/lastenheft.md),
[ADR-0058](../../adr/0058-testansatz-fuenf-luecken.md) (Entscheidung 5 —
Workflow-Wahl, Messmethode, Betroffene Dateien; vorab entschieden),
[ADR-0051](../../adr/0051-cicd-pipeline-github-actions.md)
(`ci.yml`/`e2e.yml`-Rollenteilung, `ci.yml` blockierend).

**Berührte Spec-Stellen:**
[`LH-QA-POR-002`](../../../../spec/lastenheft.md) §Primäre Zielplattform
Linux — bereits im Lastenheft festgelegt, dieser Slice liefert den
fehlenden sichtbaren Log-Beleg, ändert die Zusage nicht.

**Verantwortlich:** pt9912.

**Autor:** pt9912 (Planner-Lauf). **Datum:** 2026-09-14.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein zusätzlicher, benannter Schritt in `.github/workflows/ci.yml`,
unmittelbar nach dem Checkout und vor `make gates`, gibt `uname -s` und
`go env GOOS` aus und prüft beide explizit gegen `Linux`/`linux` mit
sichtbarem Fehlschlag bei Abweichung.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **`.github/workflows/e2e.yml`** — `ADR-0058` Entscheidung 5 bindet die
  Messmethode bewusst an `ci.yml` (blockierender Gate-Workflow, bei jedem
  Pull Request), nicht an den langsameren, nicht-blockierenden
  Compose-Stack von `e2e.yml`; eine Doppelung der Prüfung in `e2e.yml`
  liefert keinen zusätzlichen Beleg.
- **PostgreSQL-Versionsmatrix (`LH-QA-POR-001`) und die drei
  E2E-Testfälle/-Phasen** — `slice-062`/`slice-063`/`slice-064`; völlig
  unabhängige Kennungen, Dateien und Workflows.
- **Ein eigener, neuer Workflow ausschließlich für die
  Plattform-Assertion** — `ADR-0058` §Verglichene Alternativen (Option B)
  verwirft das ausdrücklich: unverhältnismäßiger Infrastruktur-Aufwand für
  eine Ein-Zeilen-Prüfung, widerspricht „kleinste sinnvolle Änderung"
  (`AGENTS.md` §6).
- **Eine Aussage über künftige Runner-Wechsel** — `ADR-0058` benennt das
  ausdrücklich als Grenze der gewählten Option: Der Schritt bleibt eine
  Momentaufnahme je Lauf, wie jeder andere CI-Schritt auch; eine
  weitergehende Zusicherung wäre ein anderer Vorgang.

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

- [x] `LH-QA-POR-002` erfüllt: neuer, benannter Schritt in
      `.github/workflows/ci.yml`, unmittelbar nach `Checkout` und vor dem
      `Gates`-Schritt, gibt `uname -s` und `go env GOOS` aus und schlägt
      sichtbar fehl, wenn eines von beiden nicht `Linux`/`linux` ist.
- [x] Der Schritt trägt eine sprechende `name:`-Zeile (Log-Auffindbarkeit).
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8);
      Report: `docs/reviews/review-slice-065.md`, Verdikt 0 HIGH/MEDIUM/LOW,
      1 INFO, keine Fixrunde.
- [x] Doku-Update: `harness/README.md` §Sensors, `ci.yml`-Beschreibung
      (falls dort vorhanden) um den neuen Schritt ergänzt — kein neues
      Gate, `ci.yml` bleibt derselbe blockierende Workflow.
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
| `.github/workflows/ci.yml` | update | ein neuer Schritt, zwei Befehlszeilen (`uname -s`, `go env GOOS`), platziert zwischen `Checkout` und `Gates` |

## 4. Trigger

**Start** (`next` → `in-progress`): `welle-17` eröffnet, `Verantwortlich:`
gesetzt, WIP-Limit (1 je Implementer) frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Nicht plausibel
  für einen Zwei-Zeilen-Schritt in einer bestehenden Workflow-Datei — kein
  bekannter Zerlegungsfall.
- `in-progress` → `open` (blockiert — Carveout?): Kein bekannter Blocker —
  `ADR-0058` liegt bereits Accepted vor und trifft alle Fragen (Workflow,
  Platzierung, Messmethode).

## 5. Closure-Trigger

DoD vollständig **und** `make gates` grün **und** Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- `BEO-PGC/github-actions-unverifizierbar-lokal` (Hinweis: der Plan-Text
  nannte bei Eröffnung 2×; `slice-064`s Closure hatte den Zähler bereits
  auf 3× gehoben, bevor dieser Slice `in-progress` erreichte — vom
  Reviewer als INFO F-1 vermerkt, vom Verifier bestätigt): Ob der neue
  Schritt real auf `ubuntu-latest` grün ausgibt, löst sich erst mit dem
  ersten echten `ci.yml`-Lauf auf GitHub auf. — **Ausgang: entfallen** —
  real belegt: Run `34823976027` (Implementierungs-Commit `5d7672b`)
  `completed`/`success`, von Planner UND Verifier unabhängig per `gh run
  view` bestätigt.
- Der neue Schritt könnte versehentlich vor `Checkout` oder nach `Gates`
  platziert werden. — **Ausgang: entfallen** — Reviewer UND Verifier
  bestätigten unabhängig die korrekte Position (unmittelbar nach
  Checkout, vor dem Gates-Schritt).

## 7. Closure-Notiz

- **Was hat funktioniert:** Der kleinste Slice der Welle — ein Zwei-
  Zeilen-Schritt — lief im ersten Versuch real grün, keine Fixrunde nötig.
- **Was ging anders als geplant:** Der Plan zitierte bei Eröffnung
  `BEO-PGC/github-actions-unverifizierbar-lokal` mit 2×, obwohl
  `slice-064`s Closure den Zähler bereits auf 3× gehoben hatte — der
  Reviewer fand die veraltete Momentaufnahme (INFO F-1), der Verifier
  bestätigte den aktuellen Stand. Kein Fehler in der Implementierung
  selbst, aber ein Hinweis: Slice-Pläne zitieren den Registerstand zum
  Planungszeitpunkt, nicht den zur Closure-Zeit — bei mehreren Slices
  derselben Welle kann der Stand zwischen Eröffnung und Closure
  weiterlaufen.
- **Steering-Loop-Eintrag:** keiner — die Beobachtung ist bereits durch
  `slice-064` auf 3× gehoben, der Lese-Schritt (Ausgang-Zuweisung) folgt
  bei der jetzt anstehenden `welle-17`-Closure, nicht hier.
- **Beobachtungs-Register (`../observations/`):** keine Beobachtung
  dieses Slices angefallen — der bereits erreichte 3×-Stand von
  `BEO-PGC/github-actions-unverifizierbar-lokal` stammt aus `slice-064`.
- **Folge-Slices:** keine.
- **Risiken aus §6:** beide entfallen — siehe §6.
- **Drei Paarungen:** verschoben auf `welle-17`-Closure (dieser Slice trägt
  `Welle: welle-17`, siehe DoD-Item).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Einzige berührte Sub-Area ist die
Repo-weite Default-Sub-Area `*`/`PGC` (`harness/conventions.md` führt keine
feinere Sub-Area für CI-Workflows).

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen
(`docs/plan/planning/observations/BEO-PGC/`), insbesondere gegen
GitHub-Actions-Verifikationsgrenzen:

- `BEO-PGC/github-actions-unverifizierbar-lokal` — **offen**, 2×
  (`evidence/slice-039.md`, `evidence/slice-056.md`; unter der
  3×-Schwelle) — direkter Treffer: dieselbe strukturelle Grenze wie bei
  `slice-064` (§8 dort) — ob der neue Schritt real grün läuft, löst sich
  erst mit dem ersten echten `ci.yml`-Lauf auf GitHub nach einem Push auf.
  Würde dieser Slice als dritte Instanz gezählt (3×), wäre die Beobachtung
  eine Lücke statt einer Notiz — als Risiko in §6 aufgenommen.
- `BEO-PGC/architect-verdikt-ablageort-uneinheitlich`,
  `architect-verdikt-rollen-scope-luecke` — geprüft, kein Treffer: beide
  betreffen Planning-Harness-Ablageort-/Scope-Fragen bei
  Architect-Verdikten, nicht CI-Workflow-Inhalte.
- Weitere durchgesehen (`commit-traceability-kein-vorab-hook`,
  `handbuch-versionshistorie-uebersprungen`): keine Treffer — betreffen
  Commit-Hook-Zeitpunkt bzw. Handbuch-Versionshistorie, nicht die
  Plattform-Assertion in `ci.yml`.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF (nur
`*`/`PGC`).

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
