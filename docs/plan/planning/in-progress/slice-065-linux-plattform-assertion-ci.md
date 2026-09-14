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

- [ ] `LH-QA-POR-002` erfüllt: neuer, benannter Schritt in
      `.github/workflows/ci.yml`, unmittelbar nach `Checkout` und vor dem
      `Gates`-Schritt, gibt `uname -s` und `go env GOOS` aus und schlägt
      sichtbar fehl, wenn eines von beiden nicht `Linux`/`linux` ist.
- [ ] Der Schritt trägt eine sprechende `name:`-Zeile (Log-Auffindbarkeit).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: `harness/README.md` §Sensors, `ci.yml`-Beschreibung
      (falls dort vorhanden) um den neuen Schritt ergänzt — kein neues
      Gate, `ci.yml` bleibt derselbe blockierende Workflow.
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

- `BEO-PGC/github-actions-unverifizierbar-lokal` (offen, 2× — siehe §8):
  Ob der neue Schritt real auf `ubuntu-latest` grün ausgibt, löst sich erst
  mit dem ersten echten `ci.yml`-Lauf auf GitHub nach einem Push auf —
  lokal simulierbar nur über `act` oder manuelles Nachvollziehen der
  Befehle, kein repo-interner Beweis. — **Ausgang:** <bei Closure zu
  füllen>
- Der neue Schritt könnte versehentlich vor `Checkout` oder nach `Gates`
  platziert werden, wenn die Workflow-Datei beim Einfügen nicht sorgfältig
  gelesen wird — die Messmethode verlangt „unmittelbar nach Checkout, vor
  `make gates`" (`ADR-0058` Entscheidung 5). — **Ausgang:** <bei Closure zu
  füllen>

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
