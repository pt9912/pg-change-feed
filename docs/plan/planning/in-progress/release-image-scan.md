# Slice release-image-scan: Trivy-CVE-Scan gegen das publizierte GHCR-Image

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** [welle-release-pipeline-adr-0051](../welle-release-pipeline-adr-0051.md).

**Bezug:** [`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md)
Entscheidung 6 (CVE-Scan).

**Berührte Spec-Stellen:** — (Prozess-ADR ohne Spec-Stratum).

**Verantwortlich:** Implementer-Agent (priorisiert 2026-09-19, direkter
Auftrag: "dann mach weiter bis die welle geschlossen ist").

**Autor:** Planner-Agent, direkt beauftragt. **Datum:** 2026-09-19.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die in `harness/README.md` bereits seit `slice-001` angekündigte,
aber unimplementierte Zusage einlösen: ein advisory Trivy-CVE-Scan gegen
das publizierte GHCR-`:latest`-Image, lokal über ein neues `make
image-cve`-Target (Docker-only, `AGENTS.md` §3.1) und automatisiert über
`image-scan.yml` (nächtlich + `workflow_dispatch`).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Ein Scan gegen das lokale `:dev`-Image — `ADR-0051` Entscheidung 4
  verwirft das explizit (Option D): das Dev-Image wird nie veröffentlicht,
  ein Ergebnis dagegen sagt nichts über den ausgelieferten Stand aus.
- CVE-Scan als blockierendes Gate — `ADR-0051` Entscheidung 4 verwirft das
  explizit (Option B): verletzt die Netzlos-Eigenschaft der bestehenden
  Gates und widerspricht `AGENTS.md` §3.6 im Ergebnis.
- Automatisches Beheben gefundener CVEs — bleibt bewusster Commit
  (Dockerfile-Basis-Update), das Scannen allein schließt keine Lücke
  (`ADR-0051` §Konsequenzen).
- Ein GHCR-`:latest`-Image existiert erst nach dem ersten echten Release
  (`release-version-und-workflow`, außerhalb dieser Welle) — bis dahin
  läuft `image-scan.yml` real, aber gegen ein noch nicht existierendes
  Tag; das ist eine strukturelle Konsequenz der Reihenfolge, kein
  Implementierungsfehler dieses Slice (siehe §6).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] `make image-cve` existiert (kein Gate, braucht Netz, analog
      `make image-stale`): Trivy CRITICAL/HIGH gegen das publizierte
      GHCR-`:latest`-Image von `ghcr.io/pt9912/pg-change-feed`, Ergebnis auf
      stdout, `exit 1` bei mindestens einem CRITICAL/HIGH-Fund,
      `exit 0` sonst. Digest-gepinntes `aquasec/trivy` (v0.74.0, real
      gegen `docker manifest inspect` verifiziert). Trägt optional
      `GHCR_USERNAME`/`GHCR_PASSWORD` für ein privates GHCR-Paket (Trivys
      eigene `--username`/`TRIVY_PASSWORD`-Mechanik statt eines
      allgemeinen Docker-Credential-Mounts — Letzteres färbte real auch
      Trivys eigenen, unauthentifizierten Vulnerability-DB-Bezug ein,
      sobald `~/.docker/config.json` einen nicht ausführbaren
      Credential-Helper referenziert). `--image-src remote` verhindert
      zusätzlich unnötige lokale Docker-/Containerd-/Podman-Socket-Proben.
      Real ausgeführt (mit und ohne Credentials): scheitert strukturell
      korrekt am fehlenden `ghcr.io/pt9912/pg-change-feed:latest` (kein
      echter Release bisher), Trivy selbst lief real und lud seine
      Vulnerability-DB.
- [x] `image-scan.yml` existiert: ruft `make image-cve` mit
      `GHCR_USERNAME`/`GHCR_PASSWORD` aus `github.actor`/
      `secrets.GITHUB_TOKEN` auf (`packages: read`-Permission), Trigger
      `schedule` (nächtlich) + `workflow_dispatch`; ein roter Lauf ist
      sichtbar, blockiert aber `make gates`/`ci.yml`/`release.yml` nicht
      (advisory, [`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md)
      Entscheidung 6). YAML-Struktur per Ruby-Stdlib-YAML-Parser geprüft.
- [x] `harness/README.md` §Werkzeuge trägt die reale `make image-cve`-Zeile
      und ersetzt den bisherigen „Nicht behauptet (geplant)"-Hinweis am
      Dateiende (`AGENTS.md` §4: kein behauptetes Gate/Target ohne Deckung).
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      2 MEDIUM (F-1: fehlende GHCR-Authentifizierung für private Pakete;
      F-2: instabiler Slice-Plan-Pfadverweis im Workflow-Kommentar) und
      2 LOW (F-3: fehlendes `--image-src remote`; F-4: Tippfehler) in
      derselben Fixrunde behoben, kein offenes HIGH.
- [ ] Doku-Update für `harness/README.md` §Werkzeuge — entfällt als
      eigener Punkt, da bereits §2 oben dieselbe Zeile explizit als
      DoD-Kriterium trägt.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `Makefile` | update | neues `image-cve`-Target, Trivy-Aufruf gegen `ghcr.io/pt9912/pg-change-feed:latest`, gepinntes Trivy-Image (Digest-Pin, Modul 14). |
| `.github/workflows/image-scan.yml` | neu | `schedule` + `workflow_dispatch`, ruft `make image-cve` auf, `permissions: {}` mit gezielter Lockerung. |
| `harness/README.md` §Werkzeuge | update | reale `make image-cve`-Zeile, „Nicht behauptet (geplant)"-Absatz entfernt. |
| `Makefile` (`image-cve`-Target) | update (Plan-Nachzug, Fixrunde) | Reviewer-Finding F-1: GHCR-Pakete aus `GITHUB_TOKEN`-Pushes entstehen unabhängig von der Repo-Sichtbarkeit privat — `GHCR_USERNAME`/`GHCR_PASSWORD` (optional) authentifizieren über Trivys eigene Flags; `--image-src remote` (F-3) verhindert unnötige lokale Runtime-Proben. |
| `.github/workflows/image-scan.yml` | update (Plan-Nachzug, Fixrunde) | `packages: read` + `GHCR_USERNAME`/`GHCR_PASSWORD` aus `github.actor`/`secrets.GITHUB_TOKEN` (F-1); stabiler Slice-Bezug statt Lifecycle-Pfad im Kommentar (F-2); Tippfehler behoben (F-4). |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): sofort — Welle eröffnet, unabhängig von
`release-version-und-workflow` (siehe Welle-Datei §5 Abhängigkeiten).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): entfällt aus
  heutiger Sicht — Umfang ist ein Make-Target plus ein Workflow.
- `in-progress` → `open` (blockiert — Carveout?): kein gepinntes,
  offizielles Trivy-Image auffindbar oder dessen Lizenz/Nutzung verbietet
  den vorgesehenen Einsatz — dann Carveout mit Alternativwerkzeug-Prüfung.

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + Closure-Notiz mit Lerneintrag
geschrieben.

## 6. Risiken und offene Punkte

- Bis zum ersten echten Release (`release-version-und-workflow`, siehe
  Welle-Datei §6 Out-of-Scope) existiert kein GHCR-`:latest`-Image —
  `image-scan.yml` läuft real, scheitert aber am fehlenden Ziel-Image, kein
  Implementierungsfehler. **Ausgang:** eingetreten und real bestätigt
  (`make image-cve` real ausgeführt: Trivy lud seine Vulnerability-DB und
  scheiterte danach korrekt mit einem GHCR-„DENIED"-Fehler auf das nicht
  existierende `:latest`-Tag).
- Review-Finding F-1: ein per `GITHUB_TOKEN` gepushtes GHCR-Paket entsteht
  unabhängig von der Repo-Sichtbarkeit privat — ohne Authentifizierung
  bliebe der `DENIED`-Fehler nach dem ersten Release bestehen, nur aus
  einem anderen Grund. **Ausgang:** eingetreten, in derselben Fixrunde
  behoben (`GHCR_USERNAME`/`GHCR_PASSWORD` über Trivys eigene Flags,
  `packages: read`-Permission in `image-scan.yml`) — der volle
  Authentifizierungspfad gegen ein tatsächlich privates GHCR-Paket bleibt
  bis zum ersten echten Release strukturell unverifiziert (ein lokaler
  Nachbau mit einer passwortgeschützten Test-Registry scheiterte an einer
  bekannten `htpasswd`/bcrypt-Inkompatibilität der `registry:2`-Referenz-
  implementierung, kein Befund gegen die hier gewählte Lösung selbst).
- `AGENTS.md` §3.10 gilt für `image-scan.yml` als neuen Workflow
  unverändert: `make gates` grün belegt nicht, dass der reale
  Post-Push-Lauf grün läuft. **Ausgang:** weiter offen, strukturell
  (derselbe Fall wie `BEO-PGC/github-actions-unverifizierbar-lokal`,
  bereits verkörpert als `AGENTS.md` §3.10).
- Trivys CVE-Datenbank ändert sich täglich unabhängig vom Code-Stand — ein
  heute grüner Scan kann morgen ohne Commit rot werden. **Ausgang:**
  entfallen als Risiko für dieses Slice: genau dieses Verhalten ist die
  bewusste Entscheidung von `ADR-0051` Entscheidung 4 (advisory, kein
  Gate, gerade weil ein externer Feed kein stabiler PR-Blocker sein darf).

## 7. Closure-Notiz

*(wird bei Bearbeitung gefüllt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Sub-Area „CI/CD-Pipeline"
(`.github/workflows/*.yml`, `Makefile`) — bereits mehrfach berührt
(`slice-039`, `ADR-0051` selbst), Schwelle erfüllt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/github-actions-unverifizierbar-lokal` (7×, bereits verkörpert,
siehe §6) betrifft diese Sub-Area; kein weiterer Treffer.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
