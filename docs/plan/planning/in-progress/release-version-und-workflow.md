# Slice release-version-und-workflow: Versionierung + Release-Workflow (Kern-Mechanismus)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** [welle-release-pipeline-adr-0051](../welle-release-pipeline-adr-0051.md).

**Bezug:** [`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md)
Entscheidung 1/2/3 (Versionierung, Registry-Ziele, Release-Workflow).

**Berührte Spec-Stellen:** — (Prozess-ADR ohne Spec-Stratum, siehe
`ADR-0051` `Schärft:`-Feld).

**Verantwortlich:** Implementer-Agent (priorisiert 2026-09-19, direkter
Auftrag: "Ja, fang mit release-version-und-workflow an").

**Autor:** Planner-Agent, direkt beauftragt ("Release-Pipeline jetzt
wirklich bauen", "alles — leg eine Welle mit slices dafür an"). **Datum:**
2026-09-19.

---

## 1. Ziel und Abgrenzung

**Ziel:** Den in [`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md)
entschiedenen Kern-Release-Mechanismus real bauen: `docs/user/version.md`
als menschlich lesbarer Versions-Spiegel, `make image` um einen
Versions-/Tag-Parameter für den Multi-Registry-Push erweitert, und
`release.yml` (Trigger `push: tags: ['v*']`), das einen Tag strikt gegen
SemVer 2.0 und gegen `version.md` validiert, darüber baut und nach GHCR
**und** Docker Hub pusht (ein Build, Content-Mirror), abschließend ein
GitHub-Release mit Digest-Pin anlegt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- `image-scan.yml`, `upstream-drift.yml`, `hub-description.yml` — eigene
  Slices derselben Welle (`release-image-scan`, `release-upstream-drift`,
  `release-hub-description`); `hub-description.yml` braucht zudem den
  Job-Namen aus `release.yml`, den erst dieser Slice festlegt.
- `docs/user/releasing.md` — eigener Folge-Slice
  (`release-doku-releasing`), dokumentiert erst den vollständigen
  Endzustand aller vier Implementierungs-Slices.
- Ein tatsächlicher Release (`git tag` + Push) — bleibt außerhalb jeder
  Slice-DoD dieser Welle (§6 der Welle-Datei), eine irreversible, extern
  sichtbare Aktion nur nach gesonderter Rückfrage.
- Anlage der `DOCKERHUB_USERNAME`/`DOCKERHUB_TOKEN`-Repository-Secrets —
  externe, kontobezogene Handlung des Auftraggebers, kein Slice-Gegenstand;
  `release.yml` referenziert nur die Secret-Namen.
- P3–P9-Pin-Freshness und die zugehörigen neuen Make-Targets — eigener
  Slice (`release-upstream-drift`), thematisch getrennt von der
  Release-Auslösung selbst.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] `docs/user/version.md` existiert, trägt genau eine Zeile: die
      aktuelle Version als SemVer-2.0-String (`MAJOR.MINOR.PATCH[-PRERELEASE]`,
      ohne führendes `v`) — erste Version frei wählbar (z. B. `0.1.0` oder
      die bereits im Handbuch geführte `0.2.0-verdrahtung`-Zeile
      übernehmen, Entscheidung des Implementer-Laufs). Umgesetzt: `0.1.0`.
- [x] `make image` akzeptiert einen `VERSION`-Parameter (`make image
      VERSION=1.2.3`): baut über **einen** `docker buildx build`-Lauf und
      pusht (`--push`, kein `--load`) nach GHCR **und** Docker Hub mit dem
      Versions-Tag auf beiden Registries (Content-Mirror,
      [`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md)
      Entscheidung 2 Option D); `:latest` wird nur bei einem stabilen
      (nicht-Prerelease) `VERSION`-Wert zusätzlich gesetzt, auf beiden
      Registries. Ohne `VERSION` bleibt das bestehende Verhalten
      (`--load`, nur `:dev`, [`ADR-0044`](../../adr/0044-image-beleg-semantik.md))
      unverändert — real gegen zwei lokale Registry-Container geprüft:
      derselbe Digest auf allen vier Tag/Registry-Kombinationen
      (`docker buildx imagetools inspect`, siehe §7).
- [x] `release.yml` existiert: Trigger `push: tags: ['v*']`; ein erster
      Schritt validiert den Tag fail-fast strikt gegen SemVer 2.0 **und**
      gegen den Wert in `docs/user/version.md` am getaggten Commit (Abbruch
      bei jeder Abweichung, vor jedem Login/Build/Push); baut über
      `make image VERSION=<validierte Version>`; legt danach ein
      GitHub-Release an, dessen Beschreibungstext den Image-Digest trägt.
      YAML-Struktur und jeder `run:`-Schritt syntaktisch geprüft (Ruby-
      Stdlib-YAML-Parser bzw. `bash -n`, siehe §7) — der reale
      Post-Push-Lauf bleibt unverifiziert (`AGENTS.md` §3.10, §6).
- [x] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update für `harness/README.md` (`make image`-Zeile mit dem
      neuen `VERSION`/`LATEST`-Verhalten, neue `release.yml`-Zeile).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/user/version.md` | neu | Versions-Spiegel ([`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md) Entscheidung 1). |
| `Makefile` (`image`-Target) | update | `VERSION`-Parameter, bedingter `--push`-Pfad mit Multi-Registry-/`:latest`-Tags statt des bestehenden `--load`-Pfads. |
| `.github/workflows/release.yml` | neu | Tag-Validierung (SemVer 2.0 + `version.md`-Abgleich), Build über `make image VERSION=...`, GitHub-Release mit Digest-Pin. |
| `AGENTS.md` §3.8 | keine Änderung | bereits vorhanden (Action-Pinning) — jede neue `uses:`-Zeile in `release.yml` folgt der bestehenden Regel, kein neuer Regeltext. |
| `harness/README.md` §Sensors/§Werkzeuge | update | `release.yml`-Zeile analog zu `e2e.yml` (kein Gate, `ADR-0051`). |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): sofort — Welle eröffnet, kein externer
Trigger jenseits der Welle-Eröffnung selbst.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): der
  `make image`-Umbau erweist sich als invasiver als erwartet (z. B. weil
  der bestehende `--load`-Pfad und der neue `--push`-Pfad sich nicht
  sauber in einem Target vereinen lassen) — dann Aufteilung in einen
  reinen `make image`-Umbau-Slice und einen `release.yml`-Slice.
- `in-progress` → `open` (blockiert — Carveout?): GitHub Actions verlangt
  für den GHCR-/Docker-Hub-Push eine Secret-Konfiguration, die real nur
  der Auftraggeber anlegen kann — blockiert nicht die Implementierung
  (die Datei referenziert nur Secret-Namen), wohl aber jeden Versuch, den
  Workflow real laufen zu lassen; kein Blocker für die DoD dieses Slice.

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + Closure-Notiz mit Lerneintrag
geschrieben.

## 6. Risiken und offene Punkte

- `release.yml` kann erst nach einem echten Tag-Push real (grün) laufen —
  `AGENTS.md` §3.10 gilt unverändert: der Slice gilt bei `make gates` grün
  als DoD-erfüllt, der reale Post-Push-Lauf bleibt bis zum ersten
  tatsächlichen Release unbewiesen. **Ausgang:** weiter offen,
  strukturell (kein Sensor kann das schließen), dokumentiert analog
  `BEO-PGC/github-actions-unverifizierbar-lokal` (bereits verkörpert als
  `AGENTS.md` §3.10 — kein neuer Eintrag nötig, nur derselbe Fall erneut).
- Ein Build unter `--push` statt `--load` ändert den Digest-Ermittlungspfad
  von `make image` (`harness/image-hash.raw`/`harness/image-hash.txt`,
  [`ADR-0103`](../../adr/0103-image-hash-lokal-statt-committet.md)) — ein
  gepushtes Multi-Platform-Manifest hat einen anderen Digest-Typ als ein
  lokal geladenes Single-Platform-Image. **Ausgang:** weiter offen bis zur
  Implementierung, dort real gegen `ADR-0044`/`ADR-0103` zu prüfen (kein
  Widerspruch erwartet, da beide ADRs den lokalen `:dev`-Pfad unverändert
  lassen — zu verifizieren, nicht anzunehmen).
- Docker Hub verlangt eine andere Namensform als GHCR
  (`docker.io/<user>/<repo>` vs. `ghcr.io/pt9912/pg-change-feed`) — der
  genaue Docker-Hub-Repository-Name ist noch nicht festgelegt. **Ausgang:**
  weiter offen, Implementer-Entscheidung (mutmaßlich `pt9912/pg-change-feed`,
  gegen ein real existierendes Docker-Hub-Konto zu prüfen).

## 7. Closure-Notiz

*(wird bei Bearbeitung gefüllt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Sub-Area „CI/CD-Pipeline"
(`.github/workflows/*.yml`, `Makefile`) — bereits mehrfach berührt
(`slice-039`, `ADR-0051` selbst), Schwelle erfüllt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/github-actions-unverifizierbar-lokal` (7×, bereits verkörpert als
`AGENTS.md` §3.10, siehe §6) und
`BEO-PGC/kein-echter-versionswechsel-upgrade-test` (1×, dessen
Re-Evaluierungs-Trigger diese Welle nur zur Hälfte erfüllt, siehe
Welle-Datei §6 Out-of-Scope) betreffen diese Sub-Area; kein weiterer
Treffer.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Modus-Deklaration
`harness/conventions.md`, Default `PGC`/Greenfield).
