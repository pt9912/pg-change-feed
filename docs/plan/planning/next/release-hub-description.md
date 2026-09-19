# Slice release-hub-description: Docker-Hub-Beschreibungs-Sync nach erfolgreichem Release

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** [welle-release-pipeline-adr-0051](../welle-release-pipeline-adr-0051.md).

**Bezug:** [`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md)
Entscheidung 8 (`hub-description.yml`).

**Berührte Spec-Stellen:** — (Prozess-ADR ohne Spec-Stratum).

**Verantwortlich:** Implementer-Agent (priorisiert 2026-09-19, direkter
Auftrag: "dann mach weiter bis die welle geschlossen ist").

**Autor:** Planner-Agent, direkt beauftragt. **Datum:** 2026-09-19.

---

## 1. Ziel und Abgrenzung

**Ziel:** Nach einem erfolgreichen Release die Docker-Hub-Repository-
Beschreibung automatisch mit dem Repo-`README.md` synchronisieren:
`hub-description.yml` als eigener, `workflow_dispatch`-fähiger Workflow,
aus `release.yml` heraus als zusätzlicher Job mit `needs: <Release-Job>`
aufgerufen (`workflow_call`) — ein Fehlschlag hier macht das Release
selbst nicht rot, es ist Präsentation, kein Bestandteil der
Distributions-Zusage (`ADR-0051` Entscheidung 8).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Der eigentliche Release-Job (Build, Push, GitHub-Release) — liegt in
  `release-version-und-workflow`; dieser Slice fügt `release.yml`
  ausschließlich den zusätzlichen `needs`-Job hinzu.
- Inhaltliche Pflege der Docker-Hub-Beschreibung selbst — der Workflow
  spiegelt das bestehende `README.md`, verfasst keinen eigenen Text.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] `hub-description.yml` existiert: `on: {workflow_dispatch:,
      workflow_call:}`, synchronisiert die Docker-Hub-Beschreibung des
      Repos `pt9912/pg-change-feed` mit dem Inhalt von `README.md`
      (`DOCKERHUB_USERNAME`/`DOCKERHUB_TOKEN` referenziert, nicht angelegt
      — siehe §1 Abgrenzung der Welle-Datei).
- [ ] `release.yml` (aus `release-version-und-workflow`) bekommt einen
      zusätzlichen Job mit `needs: <Release-Job-Name>`, der
      `hub-description.yml` über `uses: ./.github/workflows/hub-description.yml`
      aufruft; ein Fehlschlag dieses Jobs lässt den Release-Job selbst
      unberührt (kein `needs`-Failure-Propagation-Block zurück).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update für `harness/README.md` §Sensors/§Werkzeuge — die
      `release.yml`-Zeile aus `release-version-und-workflow` nennt den
      zusätzlichen `hub-description`-Job knapp mit.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.github/workflows/hub-description.yml` | neu | `workflow_dispatch` + `workflow_call`, Docker-Hub-Beschreibungs-Sync. |
| `.github/workflows/release.yml` | update | zusätzlicher `needs: <Release-Job>`-Job, ruft `hub-description.yml` per `uses:` auf. |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `release-version-und-workflow` liegt in
`done/` (siehe Welle-Datei §5 Abhängigkeiten — braucht den realen
Release-Job-Namen).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): entfällt aus
  heutiger Sicht — Umfang ist ein Workflow plus ein Job.
- `in-progress` → `open` (blockiert — Carveout?): `release.yml` aus dem
  Vorgänger-Slice trägt eine Struktur, die einen `workflow_call` nicht
  sauber aufnimmt (z. B. Matrix-Job ohne klaren Einzel-Job-Namen) — dann
  Carveout mit Alternativ-Verdrahtung (`workflow_run`-Trigger statt
  `needs`-Job).

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + Closure-Notiz mit Lerneintrag
geschrieben.

## 6. Risiken und offene Punkte

- `AGENTS.md` §3.10 gilt für `hub-description.yml` und die Änderung an
  `release.yml` unverändert. **Ausgang:** weiter offen, strukturell
  (derselbe Fall wie `BEO-PGC/github-actions-unverifizierbar-lokal`,
  bereits verkörpert).
- Ein Docker-Hub-API-Aufruf zur Beschreibungs-Aktualisierung kann eine
  andere Authentifizierungsform verlangen als der reine Image-Push
  (`DOCKERHUB_TOKEN`s Scope reicht möglicherweise nicht). **Ausgang:**
  weiter offen bis zur Implementierung, dort real gegen die
  Docker-Hub-API zu prüfen (kein Repository-Secret, das dieser Slice
  technisch anlegen könnte — siehe Welle-Datei §6 Out-of-Scope).

## 7. Closure-Notiz

*(wird bei Bearbeitung gefüllt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Sub-Area „CI/CD-Pipeline"
(`.github/workflows/*.yml`) — bereits mehrfach berührt, Schwelle erfüllt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/github-actions-unverifizierbar-lokal` (7×, verkörpert, siehe §6)
betrifft diese Sub-Area; kein weiterer Treffer.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
