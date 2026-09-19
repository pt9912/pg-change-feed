# Slice release-hub-description: Docker-Hub-Beschreibungs-Sync nach erfolgreichem Release

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** [welle-release-pipeline-adr-0051](welle-release-pipeline-adr-0051.md).

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

- [x] `hub-description.yml` existiert: `on: {workflow_dispatch:,
      workflow_call: {secrets: {DOCKERHUB_USERNAME, DOCKERHUB_TOKEN}}}`,
      synchronisiert die Docker-Hub-Beschreibung des Repos
      `pt9912/pg-change-feed` mit dem Inhalt von `README.md` über
      `POST /v2/auth/token` ({identifier, secret} → access_token) +
      `PATCH /v2/repositories/pt9912/pg-change-feed` ({full_description}),
      reiner curl/jq-Aufruf auf dem Runner (kein gepinnter Action-Fork).
      Real gegen die öffentliche Docker-Hub-API geprüft: 400 bei zu
      kurzem Fake-Secret, 401 bei korrekt geformtem, falschem Secret
      (bestätigt Endpoint + Feldnamen + dass das Repo existiert); der
      eigene Fehlerpfad (kein `access_token` in der Antwort → `exit 1`
      mit `::error::`) real mit Fake-Zugangsdaten ausgelöst und bestätigt.
      `DOCKERHUB_USERNAME`/`DOCKERHUB_TOKEN` referenziert, nicht angelegt
      — siehe §6 Out-of-Scope der Welle-Datei. Token-Extraktion über das
      neue, netzlos testbare `tools/harness/dockerhub-token.sh`
      (`make test-dockerhub-token`, Review-Finding F-3) statt Inline-jq.
- [x] `release.yml` (aus `release-version-und-workflow`) bekommt einen
      zusätzlichen Job `hub-description` mit `needs: release`, der
      `hub-description.yml` über `uses: ./.github/workflows/hub-description.yml`
      mit `secrets: inherit` aufruft; ein Fehlschlag dieses Jobs lässt den
      bereits abgeschlossenen `release`-Job unberührt — GitHub Actions
      ändert dessen Status nicht rückwirkend (kein
      `needs`-Failure-Propagation-Block zurück).
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      2 HIGH (F-1/F-2: „Zitat nennt die falsche Stelle" — ADR-Abschnitt
      bzw. Welle-Datei-Abschnitt falsch benannt, Aussagen selbst korrekt)
      und 1 MEDIUM (F-3: fehlende netzlose Negativtest-Abdeckung für die
      Token-Extraktion) sowie 1 LOW (F-4: `curl -f`-Asymmetrie nur in der
      Commit-Message erklärt) und 2 INFO (F-5/F-6: Erfolgspfad und
      Permissions-Vererbung erst durch realen Lauf klärbar) in derselben
      Fixrunde behoben bzw. dokumentiert, kein offenes HIGH.
- [ ] Doku-Update für `harness/README.md` §Sensors/§Werkzeuge — entfällt
      als eigener Punkt, da bereits §2 oben dieselbe Zeile explizit als
      DoD-Kriterium trägt.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — entfällt: keine Reconciliation-Datei in diesem Repo (kein Brownfield-Bootstrap).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neuer Beleg `evidence/slice-release-hub-description.md` unter der bestehenden `BEO-PGC/zitat-nennt-die-falsche-stelle/` (bereits verkörpert, jetzt 5./6. Beleg über zwei Fundstellen F-1/F-2).
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit). Dieser Slice gehört zu `welle-release-pipeline-adr-0051` (noch offen) — Prüfung folgt regelkonform bei deren Closure.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.github/workflows/hub-description.yml` | neu | `workflow_dispatch` + `workflow_call`, Docker-Hub-Beschreibungs-Sync. |
| `.github/workflows/release.yml` | update | zusätzlicher `hub-description`-Job (`needs: release`), ruft `hub-description.yml` per `uses:`/`secrets: inherit` auf. |
| `harness/README.md` §Werkzeuge | update (Plan-Nachzug) | `release.yml`-Zeile um den `hub-description`-Job ergänzt, neue eigene `hub-description.yml`-Zeile. |
| `tools/harness/dockerhub-token.sh` | neu (Plan-Nachzug, Fixrunde) | Reviewer-Finding F-3 (MEDIUM): netzlos testbare Token-Extraktion (grep/sed, kein jq) statt Inline-jq im Workflow, analog `tools/harness/release-tag-info.sh`. |
| `tools/harness/run-dockerhub-token-tests.sh` + `make test-dockerhub-token` | neu (Plan-Nachzug, Fixrunde) | Tabellentest mit canned JSON-Antworten (real gegen die echte API beobachtete Formen) gegen `dockerhub-token.sh`. |

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
  eingetreten — im Schwester-Repo d-check real dokumentiert
  (`packaging/dockerhub/README.md` §Transport): ein Token mit
  `read/write`-Scope scheiterte dort am `PATCH`-Aufruf mit `403
  Forbidden`, obwohl derselbe Token den Image-Push erfolgreich
  authentifizierte — erst `read/write/delete`-Scope behob es, ohne den
  Token-Wert selbst zu ändern. Der Kopfkommentar von
  `hub-description.yml` trägt diese Anforderung jetzt explizit. Der
  volle Erfolgspfad mit einem real gesetzten, korrekt skopierten Token
  bleibt trotzdem bis zum ersten echten Lauf unbewiesen (siehe
  Risiko-Einträge unten).
- Reviewer-Finding F-5 (INFO): die Erfolgspfad-Feldnamen (`access_token`
  bei `POST /v2/auth/token`, Erfolgs-Statuscode bei `PATCH
  .../repositories/...`) bleiben strukturell unbewiesen, bis der
  Workflow einmal mit echten `DOCKERHUB_USERNAME`/`DOCKERHUB_TOKEN`
  läuft — alle realen Prüfungen (Implementer, Review) trafen
  ausschließlich Fehlerpfade (`400`/`401` ohne gültige Zugangsdaten).
  **Ausgang:** weiter offen, strukturell (dieselbe Kategorie wie
  `AGENTS.md` §3.10, hier auf die konkreten Feldnamen zugespitzt).
- Reviewer-Finding F-6 (INFO): ob GitHub Actions bei `uses: ./…yml` ohne
  expliziten `permissions:`-Block im Aufrufer-Job (`release.yml`s
  `hub-description`-Job) die in `hub-description.yml` selbst
  deklarierten Job-Permissions (`contents: read`) gewährt oder den
  Aufrufer-Default (`permissions: {}`) durchreicht, lässt sich nicht
  netzlos/lokal klären. **Ausgang:** eingetreten und real behoben — der
  erste echte Tag-Push (`v0.1.0`, 2026-09-19) zeigte real: GitHub Actions
  reicht den Aufrufer-Default (`permissions: {}`) durch und lässt einen
  reusable Workflow **nicht** eigenständig auf `contents: read`
  eskalieren (dokumentiertes Verhalten: „Reusable workflows can only
  downgrade permissions from the caller, never elevate them") — der
  gesamte Lauf scheiterte mit `startup_failure` vor jedem Job-Start,
  ohne dass ein einziger Job (auch nicht `release` selbst) lief. Behoben
  durch expliziten `permissions: {contents: read}`-Block am
  `hub-description`-Job in `release.yml`; der Fix selbst noch nicht
  erneut real gegen einen Tag-Push verifiziert im Moment dieses
  Nachtrags (folgt unmittelbar).

## 7. Closure-Notiz

- **Was hat funktioniert:** Die reale, hands-on-Verifikation der
  Docker-Hub-API (vier `curl`-Aufrufe mit absichtlich ungültigen
  Zugangsdaten, sowohl vom Implementer als auch unabhängig vom Reviewer)
  bestätigte Endpoint-URLs und Feldnamen (`identifier`/`secret`/
  `access_token`/`full_description`), ohne echte Zugangsdaten zu
  brauchen — inklusive der Nebenerkenntnis, dass `pt9912/pg-change-feed`
  auf Docker Hub bereits existiert (401 statt 404 auf den PATCH-Versuch).
  Der Hinweis des Nutzers auf das Schwester-Repo d-check
  (`packaging/dockerhub/README.md`) lieferte einen bereits real
  dokumentierten Produktionsfehler (Token-Scope `read/write` reicht für
  den Image-Push, aber nicht für den Beschreibungs-`PATCH` —
  `403 Forbidden`) und ließ sich direkt als Kopfkommentar-Hinweis in
  `hub-description.yml` übernehmen, bevor der Nutzer das eigene Token
  überhaupt angelegt hat.
- **Was ging anders als geplant:** Der Reviewer fand 2 HIGH (F-1/F-2:
  „Zitat nennt die falsche Stelle" — ein ADR-Abschnitt und ein
  Welle-Datei-Abschnitt waren falsch benannt, die dahinterstehenden
  Aussagen selbst aber beide sachlich korrekt) und 1 MEDIUM (F-3:
  fehlende committete, netzlose Negativtest-Abdeckung für die einzige
  nicht-triviale Logik dieses Slices — die Token-Extraktion) sowie 1 LOW
  (F-4: `curl -f`-Asymmetrie nur in der Commit-Message erklärt) und
  2 INFO (F-5/F-6: Erfolgspfad-Feldnamen und Reusable-Workflow-
  Permissions-Vererbung bleiben strukturell erst durch einen realen Lauf
  klärbar). Alle in einer Fixrunde behoben bzw. dokumentiert
  (`5cc7fa49`) — F-3 führte zu einem neuen, wiederverwendbaren
  Testartefakt (`tools/harness/dockerhub-token.sh` +
  `run-dockerhub-token-tests.sh` + `make test-dockerhub-token`), F-5/F-6
  blieben bewusst als offene §6-Risiken statt code-seitig „gelöst" zu
  werden, da beide strukturell erst durch einen realen Post-Push-Lauf
  klärbar sind.
- **Steering-Loop-Eintrag:** kein neuer Sensor, keine geschärfte Regel —
  F-1/F-2 sind der 5./6. Beleg der bereits verkörperten Beobachtung
  `BEO-PGC/zitat-nennt-die-falsche-stelle`. F-3 bestätigt und erweitert
  aber ein wiederverwendbares Struktur-Muster dieser Welle: Inline-
  Shell-Logik in einem GitHub-Actions-Workflow bekommt ein eigenständiges,
  netzlos testbares Skript statt unbelegter Inline-Logik — nach
  `tools/harness/release-tag-info.sh` (`release-version-und-workflow`)
  jetzt zum zweiten Mal angewandt (`tools/harness/dockerhub-token.sh`).
- **Beobachtungs-Register (`../observations/`):** neuer Beleg
  `evidence/slice-release-hub-description.md` unter der bestehenden
  `BEO-PGC/zitat-nennt-die-falsche-stelle/` (bereits verkörpert, jetzt
  6× über zwei Fundstellen in diesem Slice; Zähler in `state.md`
  nachgezogen).
- **Folge-Slices:** `release-doku-releasing` — bereits als Datei in
  `open/` vorhanden (Welle `welle-release-pipeline-adr-0051`, letzter
  verbleibender Slice). Kein neuer Folge-Slice aus diesem Slice selbst.
- **Risiken aus §6:** vier Risiken, vier Ausgänge — (1) `AGENTS.md` §3.10
  realer Post-Push-Lauf → **weiter offen**, bereits verkörpert; (2)
  Docker-Hub-Token-Scope für den Beschreibungs-`PATCH` →
  **eingetreten**, real im Schwester-Repo d-check dokumentiert
  (`read/write/delete` statt `read/write` nötig), Kopfkommentar trägt den
  Hinweis jetzt; (3) F-5 Erfolgspfad-Feldnamen unbewiesen →
  **weiter offen**, strukturell; (4) F-6 Permissions-Vererbung bei
  Reusable-Workflow-Aufruf ungeklärt → **weiter offen**, strukturell.
- **Drei Paarungen:** dieser Slice gehört zu
  [welle-release-pipeline-adr-0051](welle-release-pipeline-adr-0051.md)
  (noch offen) — die Prüfung läuft regelkonform bei deren Closure, nicht
  hier (§2 DoD-Zeile „im Repo mit Wellen von der nächsten
  Welle-Closure"). Vorab-Hinweis für diese spätere Prüfung: kein
  `liegt in`-Feld in diesem Slice; der einzige Folge-Slice
  (`release-doku-releasing`) existiert bereits als Datei unter
  `docs/plan/planning/open/`; der neue Beobachtungs-Beleg liegt unter
  `BEO-PGC/zitat-nennt-die-falsche-stelle/evidence/`.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Sub-Area „CI/CD-Pipeline"
(`.github/workflows/*.yml`) — bereits mehrfach berührt, Schwelle erfüllt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/github-actions-unverifizierbar-lokal` (7×, verkörpert, siehe §6)
betrifft diese Sub-Area; kein weiterer Treffer.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
