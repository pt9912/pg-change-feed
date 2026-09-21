# Slice sdk-kotlin-publish-workflow: GitHub-Packages-Publish-Workflow (`sdk-kotlin-release.yml`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** [welle-sdk-kotlin-lh-fa-sst-009](welle-sdk-kotlin-lh-fa-sst-009.md).

**Bezug:** [`LH-FA-SST-009`](../../../../spec/lastenheft.md),
[`ADR-0109`](../../adr/0109-kotlin-github-packages-drittes-sdk-package.md)
Festlegung 2/5 (Vertriebsweg, Trigger, `GITHUB_TOKEN`, Tag-Präfix),
[`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md)
Entscheidung 3/8 (Release-Tag-Trigger, Secret-Muster — Kontrastfolie: dort
brauchte jeder Vertriebsweg ein externes Registry-Secret, hier ausdrücklich
**nicht**).

**Berührte Spec-Stellen:** — (Prozess-/CI-Artefakt ohne eigene
`SPEC-*`-Kennung, analog `sdk-csharp-release.yml`/`sdk-python-release.yml`).

**Verantwortlich:** Implementer-Agent (dietmar.burkard@nerdware.dev), 2026-09-20.

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

- [x] `.github/workflows/sdk-kotlin-release.yml` existiert: Trigger
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
- [x] `sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts` trägt den
      `publishing { repositories { maven { … } } }`-Block mit dem
      eingebauten `maven-publish`-Plugin (bereits in
      `slice-sdk-kotlin-projektgeruest` angelegt) — falls dieser Block
      dort noch nicht vollständig war, ergänzt dieser Slice ihn.
- [x] Jede `uses:`-Zeile ist auf einen vollständigen Commit-SHA gepinnt,
      mit Tag-Kommentar (`AGENTS.md` §3.8) — `actions/checkout`
      wiederverwendet denselben, bereits im ganzen Repo gepinnten SHA
      (`3d3c42e5aac5ba805825da76410c181273ba90b1 # v7.0.1`); **kein**
      `actions/setup-java` nötig — JDK 21 ist auf `ubuntu-latest` bereits
      vorinstalliert (`JAVA_HOME_21_X64`, real gemessen, siehe §3 „Beim
      Schreiben getroffene Entscheidungen") — keine zusätzliche
      `uses:`-Zeile über `actions/checkout` hinaus.
- [x] YAML-Struktur geprüft (Muster `sdk-csharp-release.yml`/
      `sdk-python-release.yml`: Ruby-Stdlib-YAML-Parser, netzlos —
      `ruby -ryaml -e "YAML.load_file(...)"`, zusätzlich jeder `run:`-Block
      einzeln mit `bash -n` auf Syntaxfehler geprüft).
- [x] `harness/README.md` §Werkzeuge bekommt die reale Zeile für
      `.github/workflows/sdk-kotlin-release.yml` (kein Gate, Bindung auf
      [`ADR-0109`](../../adr/0109-kotlin-github-packages-drittes-sdk-package.md)) —
      erst jetzt zulässig, weil der Workflow jetzt real existiert
      (`AGENTS.md` §4).
- [x] `docs/user/releasing.md` zieht den neuen SDK-Release-Weg **vorab als
      eigener DoD-Punkt** nach (analog dem bestehenden C#-/Python-Eintrag,
      Abschnitt „SDK-Release" — hier **ohne** eine neue Secret-Tabellenzeile,
      stattdessen ein expliziter Hinweis: „kein externes Secret,
      `GITHUB_TOKEN` mit `packages: write`") — nicht erst nach einem
      Reviewer-Finding wie bei der C#-Welle
      (`BEO-PGC/release-mechanismus-nicht-in-releasing-doku-nachgezogen`,
      1×, Welle-Plan §6).
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      `docs/reviews/review-slice-sdk-kotlin-publish-workflow.md` fand ein
      HIGH (F-1: `./gradlew publish` lief direkt auf dem Runner statt im
      gepinnten Docker-Image, entgegen `ADR-0109` Festlegung 5 wörtlich) —
      aufgelöst durch eine Fixrunde: eine neue Docker-Stufe `publish`
      (`sdks/kotlin/Dockerfile`, baut auf `build` auf) trägt den
      Publish-Schritt jetzt Docker-only, `.github/workflows/sdk-kotlin-release.yml`
      baut/ruft diese Stufe per `docker build --target publish` /
      `docker run --rm -e GITHUB_ACTOR=… -e GITHUB_TOKEN=… … ./gradlew
      --no-daemon publish` auf, der manuelle Host-seitige `.proto`-Kopier-
      Schritt entfällt (die `.proto` fließt über denselben benannten
      Bau-Kontext `proto` ein wie bei `make sdk-pack-kotlin`). Kein offenes
      HIGH mehr — bestätigt durch einen frischen Zweit-Review-Lauf
      (`docs/reviews/review-slice-sdk-kotlin-publish-workflow-fixrunde.md`,
      eigener `docker build --target publish` + `docker inspect`-Beleg, 0
      HIGH/MEDIUM, 1 INFO).
- [x] Doku-Update: `harness/README.md` §Werkzeuge und
      `docs/user/releasing.md` (siehe oben) — entfällt als eigener Punkt,
      da bereits oben als DoD-Kriterium geführt.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben,
      **falls dieser Slice einen Inventur-Fund auflöst** — entfällt: keine
      Reconciliation-Datei in diesem Repo.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neuer
      Eintrag `BEO-PGC/implementierung-weicht-von-adr-wortlaut-ab`
      (Erstauftreten, 1×, siehe §7).
- [x] Jedes Risiko aus §6 trägt einen Ausgang (das Post-Push-Risiko bleibt
      nach `AGENTS.md` §3.10 strukturell **weiter offen**, auch bei
      grüner DoD im Übrigen; die übrigen drei Risiken sind aufgelöst bzw.
      entfallen strukturell — siehe §7).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind auf
      Slice-Ebene getragen — dieser Slice gehört zu
      [welle-sdk-kotlin-lh-fa-sst-009](welle-sdk-kotlin-lh-fa-sst-009.md)
      (Slice-Bezug oben gesetzt, kein Folge-Slice — letzter Slice der
      Welle —, Register-Eintrag siehe §7); die **volle**
      Drei-Paarungen-Prüfung der Welle selbst (Roadmap-Rückbindung,
      Wellen-Closure-Notiz) läuft regelkonform bei der separaten,
      unmittelbar nachfolgenden Welle-Closure von
      [welle-sdk-kotlin-lh-fa-sst-009](welle-sdk-kotlin-lh-fa-sst-009.md).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.github/workflows/sdk-kotlin-release.yml` | neu | Tag-Trigger, SemVer-/`build.gradle.kts`-Abgleich, `./gradlew publish` mit `GITHUB_TOKEN`. |
| `sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts` | update | `publishing.repositories.maven`-Block mit GitHub-Packages-Registry-Adresse ergänzt (der `publications`-Block existierte bereits aus `slice-sdk-kotlin-projektgeruest`). |
| `harness/README.md` §Werkzeuge | update | reale Workflow-Zeile. |
| `docs/user/releasing.md` | update | SDK-Release-Weg-Abschnitt für Kotlin, vorab, kein externes Secret; §1-Intro und Änderungshistorie (Version 1.4 → 1.5) mitgezogen. |
| `tools/harness/sdk-kotlin-release-tag-info.sh` (Plan-Nachzug, real umgesetzt) | neu | validiert `sdk-kotlin-v<SemVer-2.0>`-Tags, gibt `version=` aus (kein `latest=`) — Struktur-Vorbild `tools/harness/sdk-csharp-release-tag-info.sh`, teilt sich `tools/harness/semver-regex.sh` (Entscheidung siehe unten). |
| `tools/harness/run-sdk-kotlin-release-tag-info-tests.sh` + `Makefile`-Target `test-sdk-kotlin-release-tag-info` (Plan-Nachzug, real umgesetzt) | neu | netzloser Tabellentest, Muster `run-sdk-csharp-release-tag-info-tests.sh`/`test-sdk-csharp-release-tag-info` (kein Gate); real rot gesehen über eine Mutation (SemVer-Check auf `if false` gesetzt, 4 Fälle scheitern), danach revertiert. |
| `tools/harness/semver-regex.sh` (Plan-Nachzug, real umgesetzt) | update | Kopf-Kommentar auf den dritten Konsumenten (`sdk-kotlin-release-tag-info.sh`) erweitert — kein Verhaltens-, nur ein Dokumentationsnachzug (der dritte Konsument teilt sich die bereits bestehende `SEMVER_RE` unverändert). |

**Beim Schreiben getroffene Entscheidungen (Auflösung der unten zuvor als
„Möglicher Plan-Nachzug" benannten offenen Punkte — keine Abweichung vom
Plan, der Plan hat sie ausdrücklich auf diesen Zeitpunkt verschoben):**

- **Geteilte SemVer-2.0-Regex, wie bei C#.** `tools/harness/sdk-python-release-tag-info.sh`
  führt eine eigenständige PEP-440-Regex (reale Prüfung: dessen
  Closure-Notiz bestätigt „eigenständige Regex, nicht geteilt" — kein
  drittes Sourcing-Ziel dort). Kotlin/Gradle-Maven-Versionierung ist reines
  SemVer 2.0 (`ADR-0109` Festlegung 4) — dieselbe Grammatik wie beim
  C#-SDK. `tools/harness/sdk-kotlin-release-tag-info.sh` teilt sich deshalb
  `tools/harness/semver-regex.sh` mit `tools/harness/release-tag-info.sh`
  und `tools/harness/sdk-csharp-release-tag-info.sh` (jetzt drei
  Konsumenten) — der Kopf-Kommentar von `semver-regex.sh` ist entsprechend
  nachgezogen.
- **Versionsprüfung direkt auf dem Runner, wie geplant** —
  `grep -oE '^version = "[^"]+"' … | sed -E 's/^version = "([^"]+)"/\1/'`
  gegen die **Top-Level**-`version`-Zeile in `build.gradle.kts` (Anker
  `^`, kein Leerraum davor) trifft ausschließlich Zeile 77
  (`version = "0.1.0"`), nicht die gleichlautende, aber eingerückte Zeile
  innerhalb des `publications.create<MavenPublication>`-Blocks — real
  gegen die Datei geprüft (`grep -oE` liefert genau einen Treffer).
- **Publish-Schritt läuft auf dem Runner, nicht im Docker-Bau — mit einem
  zusätzlichen, im ursprünglichen Plan nicht ausbuchstabierten Schritt:**
  Da die `.proto`-Quelle nicht im committeten SDK-Baum liegt (`ADR-0109`
  Festlegung 3) und der Docker-Bau sie nur über den zusätzlichen
  Bau-Kontext `--build-context proto=proto` bezieht (`sdks/kotlin/Dockerfile`),
  braucht der Host-seitige `./gradlew publish`-Lauf denselben Kopier-Schritt
  manuell: `cp proto/cdc/stream/v1/changestream.proto
  sdks/kotlin/pgchangefeed-kotlin/src/main/proto/changestream.proto` VOR
  dem Publish-Schritt — ohne ihn bricht die Protobuf-Codegenerierung ab
  (`publish` hängt von `jar` ab, `jar` von `compileKotlin`, das wiederum
  den generierten gRPC-Stub braucht). Dieser Schritt stand nicht explizit
  im ursprünglichen §3-Plan, war aber durch `ADR-0109` Festlegung 3
  bereits impliziert (kein committeter Stub) — kein struktureller
  Plan-Bruch, eine notwendige Konkretisierung beim Schreiben.
- **JDK 21 auf dem Runner, kein `actions/setup-java`.** Real gemessen
  (`actions/runner-images`-Ubuntu-24.04-Readme, live abgerufen):
  `ubuntu-latest` trägt JDK 21 bereits vorinstalliert unter
  `JAVA_HOME_21_X64` (Default-`JAVA_HOME` zeigt auf JDK 17) — der
  Publish-Schritt setzt `JAVA_HOME: ${{ env.JAVA_HOME_21_X64 }}` explizit,
  kein zusätzlicher Action-Pin nötig (anders als bei
  `sdk-python-release.yml`s `astral-sh/setup-uv`, wo `uv` auf dem Runner
  nicht vorinstalliert ist).

**Hinweise aus dem Beobachtungs-Register (proaktiv):**

- `BEO-PGC/release-mechanismus-nicht-in-releasing-doku-nachgezogen` (1×):
  bereits als eigener, vorab geplanter DoD-Punkt oben aufgenommen — nicht
  erst nach einem Reviewer-Finding.
- **Backtick-Paritäts-Check** vor jedem Commit dieses Slice.
- `BEO-PGC/github-actions-unverifizierbar-lokal` (verkörpert als
  `AGENTS.md` §3.10): Das Post-Push-Risiko dieses Slice bleibt strukturell
  offen, bis ein realer Tag-Push erfolgt — kein `make gates`-Ersatz.
- Aus `slice-sdk-kotlin-pack-werkzeug` (Review F-2, MEDIUM, nicht
  merge-blockierend, vorgemerkt): Dockerfile-Kommentar/`harness/README.md`/
  Commit-Message dieses Vorgänger-Slice schreiben die
  Sources-/Javadoc-Jar-Freistellung der zitierten GitHub-Doku-Seite zu,
  obwohl diese Seite dazu **schweigt** (real per `curl` bestätigt, dreifach
  unabhängig). Der praktische Schluss bleibt haltbar (getragen vom
  tatsächlichen Gradle-Kern-Faktum `from(components["java"])` ohne
  `withSourcesJar()`/`withJavadocJar()`), nur die Attribution ist unpräzise.
  Berührt dieser Slice die Formulierung ohnehin (z. B. beim
  `docs/user/releasing.md`-Nachzug), die Zuschreibung auf das tragende
  Gradle-Faktum umstellen statt sie zu wiederholen — kein eigener DoD-Punkt,
  da nicht zwingend berührt.

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
  Kopie würde die Stellen unabhängig driften lassen. **Ausgang:** aufgelöst
  — `tools/harness/sdk-kotlin-release-tag-info.sh` teilt sich
  `tools/harness/semver-regex.sh` (dritter Konsument, Kopf-Kommentar dort
  entsprechend nachgezogen), siehe §3 „Beim Schreiben getroffene
  Entscheidungen"; netzlos getestet über
  `make test-sdk-kotlin-release-tag-info` (16 Fälle, alle bestanden,
  zusätzlich real rot gesehen über eine Mutation der Regex-Prüfung).
- Das reale `publishing.repositories.maven`-Minimalrezept aus `ADR-0109`s
  Recherche könnte beim tatsächlichen Schreiben eine zusätzliche
  Gradle-Eigenheit zeigen (z. B. Group-/Artifact-ID-Ableitung aus
  `settings.gradle.kts` statt `build.gradle.kts`), die die Doku-Recherche
  nicht erfasst hat. **Ausgang:** teilweise geprüft, im Kern weiter offen —
  `make sdk-pack-kotlin` (Docker-Bau, `./gradlew build`/`test`) bestätigt,
  dass Gradle den ergänzten `publishing.repositories.maven`-Block
  fehlerfrei einliest (kein Konfigurationsfehler beim `build.gradle.kts`-
  Parsen); ob `./gradlew publish` gegen die reale GitHub-Packages-Registry
  tatsächlich erfolgreich schreibt, bleibt nach `AGENTS.md` §3.10 bis zum
  ersten realen Tag-Push unbewiesen — falls der reale erste
  Post-Push-Lauf daran scheitert, ist das nach `ADR-0109`
  §Re-Evaluierungs-Trigger 4 eine Nachbesserung, keine automatische
  Supersession (reiner Workflow-Bugfix ohne Entscheidungsänderung).

## 7. Closure-Notiz

**Träger-Nachzug-Suchlauf (`AGENTS.md` §3.13):** Bereits vom Implementer
und beiden Reviewer-Läufen durchgeführt (`grep -rln "sdk-csharp-v\|
sdk-python-v\|NUGET_API_KEY\|PYPI_API_TOKEN\|SDK-Release"` sowie gezielt
über `spec/pflichtenheft.md`) — kein weiterer, fremder Träger außerhalb
des Diffs gefunden, der eine Eigenschaft dieses Slice beschreibt und nicht
bereits nachgezogen wäre (siehe Negativbefunde beider Review-Reports).
Diese Planner-Closure hat den Suchlauf mit `grep -rln "sdk-kotlin-release\|
gradlew publish"` eigenständig wiederholt — kein zusätzlicher Treffer über
den bereits im Diff bearbeiteten Bestand hinaus.

**Was hat funktioniert:** Die real funktionierende Fixrunden-Sequenz mit
zweitem unabhängigem Review ist der zentrale Beleg dieses Slice: Der erste
Reviewer-Lauf fand ein reales, mechanisch nachvollziehbares HIGH (F-1 —
`./gradlew publish` lief auf dem Runner statt Docker-only, entgegen
`ADR-0109` §Entscheidung Festlegung 5 wörtlich) und empfahl explizit
**keine** einseitige Implementer-Rückgabe, sondern eine Fixrunde mit
frischem Zweit-Review — genau das lief so: Der Fix (`21872c3d`, neue
Docker-Stufe `publish`) löste das Finding real auf, und der
Fixrunden-Reviewer übernahm nichts ungeprüft (eigener `docker build
--target publish` + `docker inspect`-Lauf, eigenständige Wiederholung
aller elf ursprünglich unauffälligen Prüfpunkte, eigener Träger-Nachzug-
Suchlauf) statt der Implementer-Einschätzung zu vertrauen. Der Verifier
wiederholte das Muster ein drittes Mal unabhängig (vierter Docker-Build
insgesamt) und bestätigte DoD-Konformität ohne Übernahme. Das 3-Commit-
Lifecycle-Muster (Implementer → Review → Fix → Fixrunden-Review →
Verifikation → Planner-Closure) hielt über den gesamten Zyklus, ohne dass
irgendeine Rolle die Einschätzung der vorherigen unbesehen übernahm.

**Was ging anders als geplant:** Der Slice-Plan sah beim Schreiben (§3
„Beim Schreiben getroffene Entscheidungen") den Publish-Schritt bewusst
auf dem Runner vor, mit einer eigenen Lesart von `AGENTS.md` §3.1 als
Begründung — das erwies sich beim Review als nicht tragfähig gegen den
wörtlichen `ADR-0109`-Text (Docker-only bis einschließlich `publish`).
Anders als bei den beiden Fixrunden der Vorgänger-SDK-Wellen (C#: Slice-
Chronik im Docstring; Python: Träger-Nachzug-Lücke) ist dies die erste
Fixrunde dieser dritten SDK-Welle, deren Ursache eine **Implementierungs-
Abweichung von einer bereits `Accepted`-ADR** war, nicht ein Textdefekt
oder eine Nachzug-Lücke — eine neue Fehlerklasse für das
Beobachtungs-Register dieser Welle (siehe unten). Die Auflösung blieb
trotzdem im Rahmen dieses einen Slice: ein reiner Code-Fix, keine
Supersede-ADR — `ADR-0109` §Re-Evaluierungs-Trigger 4 hatte diesen
Unterschied (Workflow-Bugfix ohne Entscheidungsänderung vs. inhaltliche
Korrektur) bereits selbst vorgesehen, bevor der Fall real eintrat.

**Steering-Loop-Eintrag:** Kein neuer Sensor, keine geschärfte
`AGENTS.md`-Regel — `AGENTS.md` §3.5 (Accepted-ADR-Immutabilität) und §3.1
(Docker-only) tragen die Regel bereits vollständig; was fehlte, war die
Disziplin, eine Abweichung als **offene Frage** statt als eigene
Entscheidung zu formulieren. Festgehalten als neuer Beobachtungs-Register-
Eintrag statt als Regelschärfung, weil ein Erstauftreten (1×) unter der
3×-Schwelle liegt, siehe unten.

**Beobachtungs-Register (`../observations/`):** Neuer Eintrag
`BEO-PGC/implementierung-weicht-von-adr-wortlaut-ab` (Erstauftreten, 1×,
`evidence/slice-sdk-kotlin-publish-workflow.md`) — geprüft, keine
bestehende Klasse passt: `BEO-PGC/fitness-function-gegen-eigene-
entscheidung` betrifft eine ADR, die sich selbst widerspricht (Fitness-
Function vs. Entscheidungstext), nicht eine Implementierung, die vom
Entscheidungstext abweicht; `BEO-PGC/gate-scope-erweiterung-ohne-adr-
traeger` betrifft einen wachsenden Gate-Scope ohne ADR-Träger, keine
Bau-/Ausführungsort-Abweichung. Die neue Klasse: ein Implementer weicht
von der wörtlichen Festlegung einer `Accepted`-ADR ab und stellt das als
bewusste Entscheidung dar, statt es als offene Frage zu kennzeichnen —
aufgelöst hier über eine Fixrunde (Code-Fix), nicht über eine
Supersede-ADR, weil der Fund den Ausführungsort betraf, keine inhaltliche
ADR-Korrektur.

**Folge-Slices:** keine — dies ist der letzte Slice der Welle
`welle-sdk-kotlin-lh-fa-sst-009`. Die Welle-Closure selbst
(Drei-Paarungen-Prüfung auf Wellen-Ebene, Roadmap-Rückbindung,
Wellen-Closure-Notiz) ist ein separater, unmittelbar nachfolgender Zug.

**Risiken aus §6 — Ausgänge:**

- `AGENTS.md` §3.10 (realer Post-Push-Lauf unverifiziert) — **weiter
  offen**, strukturell, unverändert; kein Tag-Push in diesem Zug
  (`git tag -l | grep kotlin` leer).
- Kein Secret-Anlage-Risiko — **entfallen strukturell**, unverändert:
  `GITHUB_TOKEN` ist immer verfügbar, kein externes Secret nötig.
- SemVer-Tag-Parser teilt sich Code mit `semver-regex.sh` — **aufgelöst**:
  `tools/harness/sdk-kotlin-release-tag-info.sh` ist dritter Konsument,
  `make test-sdk-kotlin-release-tag-info` grün (16 Fälle).
- Reales `publishing.repositories.maven`-Minimalrezept könnte
  Gradle-Eigenheit zeigen — **teilweise geprüft, im Kern weiter offen**:
  `make sdk-pack-kotlin` bestätigt fehlerfreies Parsen; der reale
  Schreibzugriff bleibt nach `AGENTS.md` §3.10 bis zum ersten Tag-Push
  offen.

**Drei Paarungen:** dieser Slice gehört zu
[welle-sdk-kotlin-lh-fa-sst-009](welle-sdk-kotlin-lh-fa-sst-009.md) —
dies ist der letzte Slice dieser Welle; die volle Drei-Paarungen-Prüfung
auf Wellen-Ebene läuft regelkonform bei deren eigener, separater, hier
unmittelbar nachfolgender Closure.

Review (`docs/reviews/review-slice-sdk-kotlin-publish-workflow.md`, 1
HIGH, Fixrunde `21872c3d` real behoben) +
`docs/reviews/review-slice-sdk-kotlin-publish-workflow-fixrunde.md` (0
HIGH/MEDIUM, 1 INFO) + `docs/reviews/verifikation-slice-sdk-kotlin-publish-workflow.md`
(DoD konform) liegen vor. Diese Planner-Closure schließt den Slice nach
`done/` ab (kein Self-Review — anderer Rollenwechsel, `AGENTS.md` §6).

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
