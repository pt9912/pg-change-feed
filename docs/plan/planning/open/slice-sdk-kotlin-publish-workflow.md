# Slice sdk-kotlin-publish-workflow: GitHub-Packages-Publish-Workflow (`sdk-kotlin-release.yml`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** [welle-sdk-kotlin-lh-fa-sst-009](../welle-sdk-kotlin-lh-fa-sst-009.md).

**Bezug:** [`LH-FA-SST-009`](../../../../spec/lastenheft.md),
[`ADR-0109`](../../adr/0109-kotlin-github-packages-drittes-sdk-package.md)
Festlegung 2/5 (Vertriebsweg, Trigger, `GITHUB_TOKEN`, Tag-Präfix),
[`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md)
Entscheidung 3/8 (Release-Tag-Trigger, Secret-Muster — Kontrastfolie: dort
brauchte jeder Vertriebsweg ein externes Registry-Secret, hier ausdrücklich
**nicht**).

**Berührte Spec-Stellen:** — (Prozess-/CI-Artefakt ohne eigene
`SPEC-*`-Kennung, analog `sdk-csharp-release.yml`/`sdk-python-release.yml`).

**Verantwortlich:** — (bis zur Priorisierung).

**Autor:** Planner-Agent, direkt beauftragt (`ADR-0109` §Konsequenzen
Folgepflicht 1/5). **Datum:** 2026-09-20.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein neuer, Netz-bindender GitHub-Actions-Workflow
`.github/workflows/sdk-kotlin-release.yml`, Trigger ausschließlich
`push: tags: ['sdk-kotlin-v*']`, der den Tag strikt gegen SemVer 2.0
validiert, ihn gegen die in `build.gradle.kts` geführte `version`
abgleicht (Abbruch bei Abweichung, vor jedem Build/Publish — Muster
`sdk-csharp-release.yml`s Tag-gegen-Projektdatei-Abgleich, hier gegen
`build.gradle.kts` statt `.csproj`/`pyproject.toml`), das Artefakt aus
`make sdk-pack-kotlin` erzeugt und per `./gradlew publish` gegen
`https://maven.pkg.github.com/pt9912/pg-change-feed` veröffentlicht —
authentifiziert **ausschließlich** über das eingebaute `GITHUB_TOKEN` mit
`permissions: contents: read` / `packages: write` (`ADR-0109` §Entscheidung
Festlegung 2/5, real dokumentiertes Minimalrezept) — **kein** externes
Repository-Secret, der zentrale Unterschied zu
`sdk-csharp-release.yml`/`sdk-python-release.yml` (`NUGET_API_KEY`/
`PYPI_API_TOKEN`).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Anlage eines Repository-Secrets** — entfällt strukturell: GitHub
  Packages braucht keins (`ADR-0109` §Entscheidung Festlegung 2/5). Anders
  als bei den beiden Vorgänger-SDKs gibt es hier **keinen** externen
  Wert, den ein Repository-Betreiber anlegen müsste — `GITHUB_TOKEN` wird
  von GitHub Actions automatisch je Lauf bereitgestellt.
- **Denselben Tag-Namensraum wie `release.yml`/`sdk-csharp-release.yml`/
  `sdk-python-release.yml`** — `ADR-0109` Festlegung 5 begründet
  ausdrücklich den eigenen Präfix `sdk-kotlin-v*`, um die unabhängige
  SDK-Versionierung (Festlegung 4) nicht mit den anderen Versionsräumen
  kollidieren zu lassen.
- **`:latest`-Äquivalent für das SDK** — GitHub Packages kennt kein
  Analogon zu einem Docker-`:latest`-Tag; ein veröffentlichtes Package ist
  über seine exakte Version referenzierbar, kein zusätzlicher Mechanismus
  nötig (dieselbe Begründung wie bei NuGet/PyPI).
- **GitHub-Release-Eintrag für das SDK** (analog `release.yml`s
  `gh release create`) — nicht Teil von `ADR-0109` Festlegung 5; kann ein
  späterer, eigenständiger Nachzug sein, wenn Bedarf entsteht, aber kein
  Bestandteil dieser Folgepflicht.
- **Der reale, grüne Post-Push-Lauf gegen GitHub Packages** — bleibt nach
  [`AGENTS.md`](../../../../AGENTS.md) §3.10 ein offenes Risiko dieses
  Slice bis zum ersten tatsächlichen Tag-Push (siehe §6); `make gates`
  grün und ein plausibler YAML-Aufbau sind keine Ersatz-Bestätigung.

## 2. Definition of Done

- [ ] `.github/workflows/sdk-kotlin-release.yml` existiert: Trigger
      ausschließlich `push: tags: ['sdk-kotlin-v*']` (`ci.yml`/`e2e.yml`/
      `release.yml`/`sdk-csharp-release.yml`/`sdk-python-release.yml`
      schließen diesen Tag-Namensraum durch ihre jeweils eigenen
      `tags`/`tags-ignore`-Filter aus, kein Doppellauf — Beleg im Bericht
      dieses Slice, real per `grep -n "tags"` nachgemessen, nicht
      übernommen); validiert den Tag-Suffix (nach `sdk-kotlin-v`) strikt
      gegen SemVer 2.0; liest `version` aus
      `sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts` und bricht bei
      Abweichung ab, **vor** jedem Login/Build/Publish (Muster
      `release.yml`/`sdk-csharp-release.yml`); setzt
      `permissions: contents: read` / `packages: write` auf Job-Ebene
      (`ADR-0109` §Kontext Recherche, real dokumentiertes Minimalrezept);
      ruft `make sdk-pack-kotlin` auf; veröffentlicht mit
      `./gradlew publish` gegen `https://maven.pkg.github.com/pt9912/pg-change-feed`,
      `username = System.getenv("GITHUB_ACTOR")` /
      `password = System.getenv("GITHUB_TOKEN")` im
      `publishing.repositories.maven`-Block von `build.gradle.kts` (real
      recherchiertes Minimalrezept, `ADR-0109` §Kontext) — **kein**
      `secrets.<NAME>`-Verweis auf ein externes Repository-Secret.
- [ ] `sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts` trägt den
      `publishing { repositories { maven { … } } }`-Block mit dem
      eingebauten `maven-publish`-Plugin (bereits in
      `slice-sdk-kotlin-projektgeruest` angelegt) — falls dieser Block
      dort noch nicht vollständig war, ergänzt dieser Slice ihn.
- [ ] Jede `uses:`-Zeile ist auf einen vollständigen Commit-SHA gepinnt,
      mit Tag-Kommentar (`AGENTS.md` §3.8) — `actions/checkout` und
      `actions/setup-java` (oder gleichwertig für Gradle) wiederverwenden
      denselben, bereits im Repo etablierten SHA, falls bereits gepinnt,
      sonst neu recherchiert und gepinnt.
- [ ] YAML-Struktur geprüft (Muster `sdk-csharp-release.yml`/
      `sdk-python-release.yml`: Ruby-Stdlib-YAML-Parser, netzlos —
      `ruby -ryaml -e "YAML.load_file(...)"`, zusätzlich jeder `run:`-Block
      einzeln mit `bash -n` auf Syntaxfehler geprüft).
- [ ] `harness/README.md` §Werkzeuge bekommt die reale Zeile für
      `.github/workflows/sdk-kotlin-release.yml` (kein Gate, Bindung auf
      [`ADR-0109`](../../adr/0109-kotlin-github-packages-drittes-sdk-package.md)) —
      erst jetzt zulässig, weil der Workflow jetzt real existiert
      (`AGENTS.md` §4).
- [ ] `docs/user/releasing.md` zieht den neuen SDK-Release-Weg **vorab als
      eigener DoD-Punkt** nach (analog dem bestehenden C#-/Python-Eintrag,
      Abschnitt „SDK-Release" — hier **ohne** eine neue Secret-Tabellenzeile,
      stattdessen ein expliziter Hinweis: „kein externes Secret,
      `GITHUB_TOKEN` mit `packages: write`") — nicht erst nach einem
      Reviewer-Finding wie bei der C#-Welle
      (`BEO-PGC/release-mechanismus-nicht-in-releasing-doku-nachgezogen`,
      1×, Welle-Plan §6).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: `harness/README.md` §Werkzeuge und
      `docs/user/releasing.md` (siehe oben) — entfällt als eigener Punkt,
      da bereits oben als DoD-Kriterium geführt.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben,
      **falls dieser Slice einen Inventur-Fund auflöst** — entfällt: keine
      Reconciliation-Datei in diesem Repo.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder Beleg in `evidence/`; keine Beobachtung angefallen
      ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (das Post-Push-Risiko bleibt
      nach `AGENTS.md` §3.10 strukturell **weiter offen**, auch bei
      grüner DoD im Übrigen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      dieser Slice gehört zu
      [welle-sdk-kotlin-lh-fa-sst-009](../welle-sdk-kotlin-lh-fa-sst-009.md)
      (noch offen); die Prüfung läuft regelkonform bei deren Closure.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.github/workflows/sdk-kotlin-release.yml` | neu | Tag-Trigger, SemVer-/`build.gradle.kts`-Abgleich, `./gradlew publish` mit `GITHUB_TOKEN`. |
| `sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts` | update (falls nötig) | `publishing`-Block mit GitHub-Packages-Registry-Adresse. |
| `harness/README.md` §Werkzeuge | update | reale Workflow-Zeile. |
| `docs/user/releasing.md` | update | SDK-Release-Weg-Abschnitt für Kotlin, vorab, kein externes Secret. |

**Möglicher Plan-Nachzug (analog den beiden Vorgänger-Wellen, hier vorab
benannt statt erst im Nachhinein entdeckt):**

- Ein gemeinsames Tag-Validierungs-Skript
  (`tools/harness/semver-regex.sh`, bereits aus der C#-Welle vorhanden)
  wird für ein neues `tools/harness/sdk-kotlin-release-tag-info.sh`
  wiederverwendet, analog `tools/harness/sdk-csharp-release-tag-info.sh`/
  einer möglichen Python-Entsprechung — real zu prüfen, ob eine
  Python-Entsprechung bereits existiert oder ob `sdk-python-release.yml`
  eine eigenständige Regex führt (siehe dessen Closure-Notiz: „eigenständige
  PEP-440-Einfachfall-Regex, nicht geteilt"). Für Kotlin/SemVer 2.0 ist
  eine Teilung mit `tools/harness/semver-regex.sh` naheliegend (dieselbe
  Grammatik wie bei C#), aber erst beim Schreiben zu entscheiden.
- Die `build.gradle.kts`-Versionsprüfung läuft direkt auf dem Runner
  (`grep`/`sed` gegen die feste `version = "…"`-Zeile) statt im gepinnten
  Toolchain-Container — dieselbe Begründung wie bei
  `sdk-csharp-release.yml` (kein `xmllint`/Gradle-Parser im
  `golang:1.27-alpine`-Toolchain-Image nötig für eine Ein-Zeilen-Extraktion).

**Hinweise aus dem Beobachtungs-Register (proaktiv):**

- `BEO-PGC/release-mechanismus-nicht-in-releasing-doku-nachgezogen` (1×):
  bereits als eigener, vorab geplanter DoD-Punkt oben aufgenommen — nicht
  erst nach einem Reviewer-Finding.
- **Backtick-Paritäts-Check** vor jedem Commit dieses Slice.
- `BEO-PGC/github-actions-unverifizierbar-lokal` (verkörpert als
  `AGENTS.md` §3.10): Das Post-Push-Risiko dieses Slice bleibt strukturell
  offen, bis ein realer Tag-Push erfolgt — kein `make gates`-Ersatz.

## 4. Trigger

**Start** (`next` → `in-progress`): wenn `slice-sdk-kotlin-pack-werkzeug`
in `done/` liegt (siehe Welle-Plan §4 Reihenfolge — der Workflow
veröffentlicht, was das Pack-Werkzeug real erzeugt).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): entfällt aus
  heutiger Sicht — ein Workflow, analog `sdk-csharp-release.yml`/
  `sdk-python-release.yml`.
- `in-progress` → `open` (blockiert — Carveout?): `slice-sdk-kotlin-pack-werkzeug`
  liegt noch nicht in `done/`.

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + Closure-Notiz mit Lerneintrag
geschrieben — **das Post-Push-Risiko (§6) bleibt unabhängig davon nach
`AGENTS.md` §3.10 offen**, bis ein realer Tag-Push bestätigt oder ein roter
Befund mit Folgemaßnahme dokumentiert ist; dieser Slice kann trotzdem nach
`done/` gehen, wenn das Risiko explizit als „weiter offen" geführt wird
(dieselbe Praxis wie bei den beiden Vorgänger-SDK-Wellen).

## 6. Risiken und offene Punkte

- **`AGENTS.md` §3.10 gilt unverändert:** `make gates` grün und ein
  plausibler YAML-Aufbau belegen nicht, dass der reale Post-Push-Lauf auf
  GitHub grün läuft (`./gradlew publish`-Verhalten auf dem Runner,
  GitHub-Packages-Registry-Antwortverhalten). **Ausgang:** weiter offen,
  strukturell — bestätigt oder widerlegt erst durch einen realen Tag-Push
  `sdk-kotlin-v0.1.0` (oder gleichwertig), dieselbe Klasse wie
  `BEO-PGC/github-actions-unverifizierbar-lokal` (bereits verkörpert als
  `AGENTS.md` §3.10).
- **Kein Secret-Anlage-Risiko** (anders als bei den beiden Vorgänger-SDKs:
  `NUGET_API_KEY` existierte bei der C#-Welle bereits real,
  `PYPI_API_TOKEN` fehlte bei der Python-Welle real) — GitHub Packages
  braucht kein externes Secret, `GITHUB_TOKEN` ist immer verfügbar.
  **Ausgang:** entfällt strukturell, kein Risiko dieser Klasse für diesen
  Slice.
- Ein SemVer-Tag-Parser für `sdk-kotlin-v*` teilt sich möglicherweise Code
  mit `tools/harness/semver-regex.sh` (bereits geteilt zwischen
  `release-tag-info.sh` und `sdk-csharp-release-tag-info.sh`) — eine naive
  Kopie würde die Stellen unabhängig driften lassen. **Ausgang:** weiter
  offen, zu entscheiden beim Schreiben (siehe §3 Plan-Nachzug-Hinweis).
- Das reale `publishing.repositories.maven`-Minimalrezept aus `ADR-0109`s
  Recherche könnte beim tatsächlichen Schreiben eine zusätzliche
  Gradle-Eigenheit zeigen (z. B. Group-/Artifact-ID-Ableitung aus
  `settings.gradle.kts` statt `build.gradle.kts`), die die Doku-Recherche
  nicht erfasst hat. **Ausgang:** weiter offen — falls der reale erste
  Post-Push-Lauf daran scheitert, ist das nach `ADR-0109`
  §Re-Evaluierungs-Trigger 4 eine Nachbesserung, keine automatische
  Supersession (reiner Workflow-Bugfix ohne Entscheidungsänderung).

## 7. Closure-Notiz

- **Was hat funktioniert:** <wird beim Abschluss ergänzt>
- **Was ging anders als geplant:** <wird beim Abschluss ergänzt>
- **Steering-Loop-Eintrag:** <wird beim Abschluss ergänzt>
- **Beobachtungs-Register (`../observations/`):** <wird beim Abschluss
  ergänzt>
- **Folge-Slices:** keine aus diesem Slice selbst erwartet — letzter Slice
  der Welle.
- **Risiken aus §6:** <wird beim Abschluss ergänzt>
- **Drei Paarungen:** dieser Slice gehört zu
  [welle-sdk-kotlin-lh-fa-sst-009](../welle-sdk-kotlin-lh-fa-sst-009.md)
  (noch offen) — dies ist der letzte Slice dieser Welle; die Prüfung der
  drei Paarungen läuft regelkonform bei deren eigener, separater Closure
  (nicht Teil dieses Slice-Zugs).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Sub-Area `.github/workflows/`
(CI/CD-Pipeline) — bereits mehrfach berührt (`ci.yml`, `e2e.yml`,
`release.yml`, `sdk-csharp-release.yml`, `sdk-python-release.yml`,
`image-scan.yml`, `upstream-drift.yml`, `hub-description.yml`), Schwelle
erfüllt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/github-actions-unverifizierbar-lokal` (bereits verkörpert als
`AGENTS.md` §3.10) betrifft diese Sub-Area unmittelbar (siehe §6);
`BEO-PGC/release-mechanismus-nicht-in-releasing-doku-nachgezogen` (1×)
betrifft diesen Slice über den vorab eingeplanten DoD-Punkt;
`BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit` geprüft — betrifft
mehrfach rote, aber nie behobene advisory Workflows, nicht diesen Slice
(der Workflow ist neu, keine bestehende Rot-Historie); kein weiterer
Treffer.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
