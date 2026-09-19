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
      unverändert — real gegen zwei lokale Registry-Container geprüft
      (`registry:2`, zwei Ports): derselbe Manifest-Digest auf allen vier
      Tag/Registry-Kombinationen (`docker buildx imagetools inspect`,
      Beleg im Fixrunden-Commit).
- [x] `release.yml` existiert: Trigger `push: tags: ['v*']`; die
      Tag-/Version-/Stabilitäts-Ermittlung läuft über das eigenständige,
      netzlos getestete `tools/harness/release-tag-info.sh`
      (`make test-release-tag-info`, real gegen zwölf Fälle inklusive der
      zwei vom Review gefundenen SemVer-2.0-Randfälle — führende Null in
      einem numerischen Prerelease-Identifier, Bindestrich in der
      Build-Metadata) statt einer unbelegten Inline-Prüfung; ein
      zusätzlicher Schritt gleicht die Version gegen `docs/user/version.md`
      am getaggten Commit ab (Abbruch bei jeder Abweichung, vor jedem
      Login/Build/Push); baut über `make image VERSION=<validierte
      Version>`; legt danach ein GitHub-Release an, dessen
      Beschreibungstext den Image-Digest trägt. YAML-Struktur und jeder
      `run:`-Schritt syntaktisch geprüft (Ruby-Stdlib-YAML-Parser bzw.
      `bash -n`) — der reale Post-Push-Lauf bleibt unverifiziert
      (`AGENTS.md` §3.10, siehe §6).
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      1 HIGH (F-1, unbelegter §7-Verweis) und 2 MEDIUM (F-2/F-3, reale
      SemVer-2.0-Abweichungen) sowie F-4 (fehlender Negativtest) in
      derselben Fixrunde behoben; F-5 (`ADR-0103`-Trigger-Frage) als
      Risiko in §6 übernommen statt hier aufgelöst; F-6 (Checkbox-Nachzug)
      behoben; kein offenes HIGH.
- [x] Doku-Update für `harness/README.md` (`make image`-Zeile mit dem
      neuen `VERSION`/`LATEST`-Verhalten, neue `release.yml`-Zeile).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — entfällt: keine Reconciliation-Datei in diesem Repo (kein Brownfield-Bootstrap).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neuer Beleg in `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht/evidence/`.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit). Dieser Slice gehört zu `welle-release-pipeline-adr-0051` (noch offen) — Prüfung folgt regelkonform bei deren Closure.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/user/version.md` | neu | Versions-Spiegel ([`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md) Entscheidung 1). |
| `Makefile` (`image`-Target) | update | `VERSION`-Parameter, bedingter `--push`-Pfad mit Multi-Registry-/`:latest`-Tags statt des bestehenden `--load`-Pfads. |
| `.github/workflows/release.yml` | neu | Tag-Validierung (SemVer 2.0 + `version.md`-Abgleich), Build über `make image VERSION=...`, GitHub-Release mit Digest-Pin. |
| `AGENTS.md` §3.8 | keine Änderung | bereits vorhanden (Action-Pinning) — jede neue `uses:`-Zeile in `release.yml` folgt der bestehenden Regel, kein neuer Regeltext. |
| `harness/README.md` §Sensors/§Werkzeuge | update | `release.yml`-Zeile analog zu `e2e.yml` (kein Gate, `ADR-0051`). |
| `tools/harness/release-tag-info.sh` | neu (Plan-Nachzug, Fixrunde) | Reviewer-Finding F-1/F-2/F-3: die zunächst inline in `release.yml` geführte SemVer-2.0-Validierung/Stabilitäts-Ermittlung trug zwei reale Fehler (führende Null im Prerelease-Identifier akzeptiert, Bindestrich in der Build-Metadata fälschlich als Prerelease gewertet) und keinen automatisierten Beleg — jetzt eigenständiges, netzlos testbares Skript. |
| `tools/harness/run-release-tag-info-tests.sh` | neu (Plan-Nachzug, Fixrunde) | Reviewer-Finding F-4: automatisierter Tabellentest gegen zwölf Fälle, deckt beide real gefundenen Randfälle ab. |
| `Makefile` (`test-release-tag-info`-Target) | neu (Plan-Nachzug, Fixrunde) | macht den neuen Testlauf über `make gates`-analoge Disziplin aufrufbar (kein Gate, netzlos). |

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
  lokal geladenes Single-Platform-Image. **Ausgang:** eingetreten,
  real geprüft und unauffällig: `--metadata-file`s `containerimage.digest`
  liefert in beiden Modi denselben Digest-Wert (real gegen zwei lokale
  Registry-Container verifiziert, `docker buildx imagetools inspect`
  bestätigt Übereinstimmung auf allen vier Tag/Registry-Kombinationen) —
  kein Widerspruch zu `ADR-0044`/`ADR-0103`, beide ADRs bleiben für den
  lokalen `:dev`-Pfad unverändert zutreffend.
- Docker Hub verlangt eine andere Namensform als GHCR
  (`docker.io/<user>/<repo>` vs. `ghcr.io/pt9912/pg-change-feed`). **Ausgang:**
  eingetreten, real geprüft: `docker.io/<user>/<repo>` ist syntaktisch
  korrekt und lösst identisch zur präfixlosen Docker-Hub-Form auf (real
  gegen `docker.io/library/alpine:3.20` vs. `alpine:3.20` verifiziert,
  identischer Digest) — `pt9912/pg-change-feed` als Repository-Name
  gewählt; ob dieses Repository auf Docker Hub real existiert, bleibt bis
  zum ersten echten Push unbewiesen (externe Kontoabhängigkeit, siehe
  Welle-Datei §6).
- Review-Finding F-5: `ADR-0103`s eigener Re-Evaluierungs-Trigger nennt
  „ein `docker push`-Workflow" namentlich als Auslöser für eine
  Folge-ADR („Archiv-Bedingung wird real erfüllt"). Dieser Slice fügt
  genau das hinzu. **Ausgang:** weiter offen, bewusst nicht in diesem
  Slice aufgelöst — der Trigger liest sich als „ein Lauf-Zweig *braucht*"
  einen historischen Beleg, was erst mit einem tatsächlich durchgeführten
  Push real eintritt, nicht bereits mit der bloßen Existenz des
  Mechanismus (dieselbe Abgrenzung wie beim ersten Risiko oben und wie
  `BEO-PGC/kein-echter-versionswechsel-upgrade-test`, siehe Welle-Datei
  §6: „ein tatsächlicher Release" ist explizit außerhalb dieser Welle).
  Braucht eine Architect-Entscheidung spätestens vor dem ersten echten
  Release: entweder eine Folge-ADR (`Supersedes ADR-0103`, committeter
  Digest für den Release-Pfad) oder eine explizite Begründung, warum der
  GitHub-Release-Beschreibungstext als „historisch nachschlagbarer
  Beleg" im Sinne des Triggers ausreicht, ohne `git`-Commit-Historie.

## 7. Closure-Notiz

- **Was hat funktioniert:** Die reale Zwei-lokale-Registries-Verifikation
  (`registry:2` auf zwei Ports statt echter GHCR-/Docker-Hub-Zugangsdaten)
  bewies die Content-Mirror-Eigenschaft (identischer Manifest-Digest über
  vier Tag/Registry-Kombinationen) ohne echte Zugangsdaten zu brauchen —
  ein realer, wiederholbarer Beleg statt einer bloßen Behauptung „sollte
  funktionieren". Ebenso die `docker.io/library/alpine`-Gegenprobe für die
  Docker-Hub-Namensform.
- **Was ging anders als geplant:** Der Reviewer fand 1 HIGH (F-1: zwei
  DoD-Checkboxen verwiesen auf den zu diesem Zeitpunkt noch leeren §7 —
  sechster Beleg der bereits verkörperten Reviewer-Regel „Beleg trägt
  seinen Satz nicht") und 4 MEDIUM: F-2/F-3 waren echte SemVer-2.0-Bugs
  in der ursprünglich inline in `release.yml` geführten Validierung
  (führende Null im Prerelease-Identifier akzeptiert; Bindestrich in der
  Build-Metadata fälschlich als Prerelease gewertet), F-4 der fehlende
  automatisierte Negativtest dafür, F-5 eine echte, unentschiedene
  `ADR-0103`-Governance-Frage (ihr eigener Trigger nennt „ein
  `docker push`-Workflow" als Auslöser für eine Folge-ADR). Alle außer
  F-5 in einer Fixrunde behoben (`31cf6e8a`) — F-5 bewusst als offenes
  Risiko übernommen statt vom Implementer einseitig entschieden, der
  Verifier bestätigte diese Einordnung als angemessen.
- **Steering-Loop-Eintrag:** kein neuer Sensor — F-1 ist ein weiterer
  (sechster) Beleg der bereits verkörperten Reviewer-Regel „Beleg trägt
  seinen Satz nicht" (`.harness/skills/reviewer.md`); F-2/F-3/F-4 führten
  aber zu einem echten Struktur-Ergebnis: die Tag-Validierungslogik lebt
  jetzt als eigenständiges, netzlos testbares Skript
  (`tools/harness/release-tag-info.sh` + `make test-release-tag-info`)
  statt als unbelegte Inline-Logik in einem GitHub-Actions-Workflow — ein
  wiederverwendbares Muster für künftige Workflow-interne
  Validierungslogik in diesem Repo.
- **Beobachtungs-Register (`../observations/`):** neuer Beleg
  `evidence/slice-release-version-und-workflow.md` unter der bestehenden
  `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht/` (bereits verkörpert,
  jetzt 6. Beleg).
- **Folge-Slices:** `release-image-scan`, `release-upstream-drift`,
  `release-hub-description`, `release-doku-releasing` — alle vier bereits
  als Dateien in `open/` vorhanden (Welle
  `welle-release-pipeline-adr-0051`). Zusätzlich benannt (kein
  Datei-Anlegen in diesem Zug): eine Architect-Entscheidung zu F-5 vor
  dem ersten echten Release (Folge-ADR `Supersedes ADR-0103` oder
  explizite Begründung, siehe §6).
- **Risiken aus §6:** vier Risiken, vier Ausgänge — (1) `release.yml`
  real erst nach echtem Tag-Push verifizierbar → **weiter offen**,
  strukturell (`AGENTS.md` §3.10); (2) Digest-Ermittlungspfad `--push`
  vs. `--load` → **eingetreten**, real geprüft und unauffällig (identischer
  Digest über beide Modi); (3) Docker-Hub-Namensform `docker.io/<user>/<repo>`
  → **eingetreten**, real geprüft und korrekt, reale Docker-Hub-Existenz
  bleibt bis zum ersten Push unbewiesen; (4) `ADR-0103`-Trigger-Frage
  (F-5) → **weiter offen**, Architect-Entscheidung vor dem ersten echten
  Release fällig.
- **Drei Paarungen:** dieser Slice gehört zu
  [welle-release-pipeline-adr-0051](../welle-release-pipeline-adr-0051.md)
  (noch offen) — die Prüfung läuft regelkonform bei deren Closure, nicht
  hier (§2 DoD-Zeile „im Repo mit Wellen von der nächsten
  Welle-Closure"). Vorab-Hinweis für diese spätere Prüfung: kein
  `liegt in`-Feld in diesem Slice; alle vier Folge-Slices existieren
  bereits als Dateien unter `docs/plan/planning/open/`; der neue
  Beobachtungs-Beleg liegt unter
  `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht/evidence/`.

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
