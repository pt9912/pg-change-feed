# Slice sdk-kotlin-projektgeruest: SDK-Projektgerüst `sdks/kotlin/pgchangefeed-kotlin/`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** [welle-sdk-kotlin-lh-fa-sst-009](../welle-sdk-kotlin-lh-fa-sst-009.md).

**Bezug:** [`LH-FA-SST-009`](../../../../spec/lastenheft.md) (Scope),
[`ADR-0109`](../../adr/0109-kotlin-github-packages-drittes-sdk-package.md)
Festlegung 1/2/3/4/5 (Umfang, Vertriebsweg, Ort, Versionierung,
Build-Mechanismus).

**Berührte Spec-Stellen:** [`LH-FA-SST-009.a`](../../../../spec/pflichtenheft.md)
(bleibt nach diesem Slice noch offen — der Träger-Nachzug folgt erst mit
`slice-sdk-kotlin-pack-werkzeug`, wenn das Package real paketierbar ist).

**Verantwortlich:** — (bis zur Priorisierung).

**Autor:** Planner-Agent, direkt beauftragt (`ADR-0109` §Konsequenzen
Folgepflicht 1). **Datum:** 2026-09-20.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein neues, eigenständiges Gradle-Projekt unter
`sdks/kotlin/pgchangefeed-kotlin/` anlegen — `build.gradle.kts` mit
Maven-Koordinate (`group = "io.github.pt9912"`,
`artifactId`/Projektname `pgchangefeed-kotlin`, `version = "0.1.0"`), ein
eigenes, digest-gepinntes Docker-Bau-Setup (analog `examples/kotlin/Dockerfile`,
Basis `eclipse-temurin:21-jdk` — `ADR-0109` Festlegung 5, **nicht**
`25-jdk`), ein englischsprachiges `README.md` und ein leeres/minimales
öffentliches API-Skelett — ohne jeden Import aus `internal/**`/`cmd/**`
dieses Repos.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Öffentliche HTTP-/gRPC-API-Fläche** — `slice-sdk-kotlin-http-client-flaeche`
  und `slice-sdk-kotlin-grpc-client-flaeche` übernehmen das; dieser Slice
  liefert nur den Ort, in den sie schreiben (das API-Skelett bleibt leer
  oder trägt höchstens eine gemeinsame Konfigurations-/Auth-Grundklasse,
  falls sich beim Schreiben ein echter, unstrittiger gemeinsamer Nenner
  zeigt — kein Vorgriff auf Endpunkt-Methoden, dieselbe Zurückhaltung wie
  bei `slice-sdk-csharp-projektgeruest`).
- **`make sdk-pack-kotlin` (Pack-/Publish-Werkzeug)** —
  `slice-sdk-kotlin-pack-werkzeug` übernimmt das; dieser Slice liefert nur
  ein Docker-Bau-Setup, das `./gradlew build`/`./gradlew test` trägt
  (analog `examples/kotlin/Dockerfile`s `build`-Stufe), kein
  `publish`-Aufruf und kein `make`-Ziel.
- **GitHub-Packages-Publish-Workflow** — `slice-sdk-kotlin-publish-workflow`
  übernimmt das; dieser Slice pflegt nur die `version`-Zeile in
  `build.gradle.kts`, kein Tag-Trigger, kein `GITHUB_TOKEN`-Bezug.
- **Umbau von `examples/kotlin/`** — `ADR-0109` Festlegung 3 verbietet das
  ausdrücklich; die Beispiele bleiben unverändertes Vorbild (`SPEC-023`).
- **Träger-Nachzug in `spec/pflichtenheft.md`** — `ADR-0109` §Konsequenzen
  Folgepflicht 2 bindet den Nachzug an den Zug, der das Package **real**
  existent macht (`./gradlew publishToMavenLocal`/ein Jar-Bau läuft real);
  das ist `slice-sdk-kotlin-pack-werkzeug`, nicht dieser Slice, der nur ein
  Projekt-Gerüst ohne Pack-Fähigkeit liefert.

## 2. Definition of Done

- [ ] `sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts` existiert:
      `group = "io.github.pt9912"`, Projektname/`artifactId`
      `pgchangefeed-kotlin`, `version = "0.1.0"` (`ADR-0109` Festlegung 4,
      Start bei `0.x.y`), Kotlin-Gradle-Plugin (dieselbe Version wie
      `examples/kotlin/build.gradle.kts`, real zum Bau-Zeitpunkt
      nachgemessen, nicht blind übernommen — `AGENTS.md` §3.12), das
      eingebaute `maven-publish`-Plugin (`ADR-0109` Festlegung 5, **kein**
      Dritt-Plugin wie `com.vanniktech.maven-publish`). Kein Import auf
      einen privaten Baum dieses Repos (`ADR-0109` §Kontext Bindung 5,
      Import-Grenze).
- [ ] `sdks/kotlin/pgchangefeed-kotlin/settings.gradle.kts` existiert
      (analog `examples/kotlin/settings.gradle.kts`), eigenständiger
      Gradle-Wrapper (`gradlew`/`gradlew.bat`/`gradle/wrapper/`) mit
      derselben, bereits real erprobten Gradle-Version wie
      `examples/kotlin/gradle/wrapper/gradle-wrapper.properties`
      (`8.14`, real nachgemessen, `distributionSha256Sum` gepinnt).
- [ ] `sdks/kotlin/Dockerfile` (Bau-Kontext `sdks/kotlin/`, eigenständig von
      `examples/kotlin/Dockerfile` und der Wurzel-`Dockerfile`) mit
      digest-gepinnter `eclipse-temurin:21-jdk`-Basis (real gemessener
      Digest zum Bau-Zeitpunkt, Kommentar-Pflicht analog
      `examples/kotlin/Dockerfile`); `./gradlew build`/`test` laufen darin,
      kein Runtime-Stufe nötig (ein SDK ist keine startbare Anwendung).
- [ ] `sdks/kotlin/pgchangefeed-kotlin/README.md` (Englisch, Vorgabe aus
      dem Auftrag) beschreibt Zweck, Installationsweg (Gradle-/Maven-
      Koordinate `io.github.pt9912:pgchangefeed-kotlin`, Registry-URL
      `https://maven.pkg.github.com/pt9912/pg-change-feed`) **inklusive**
      des expliziten Hinweises auf die PAT-Pflicht zum Bezug
      (`ADR-0109` §Entscheidung Festlegung 2 — GitHub Packages verlangt
      immer eine Authentifizierung zum Lesen, auch für ein öffentliches
      Package) und verweist auf das Repo-Root-`README.md` für den vollen
      Kontext; kein Duplikat der Draht-Doku (`SPEC-018`/`SPEC-020` bleiben
      die kanonische Quelle).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update für `harness/README.md` entfällt in diesem Slice — kein
      neues `make`-Target entsteht hier (Pack-Werkzeug folgt in
      `slice-sdk-kotlin-pack-werkzeug`).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben,
      **falls dieser Slice einen Inventur-Fund auflöst** — entfällt: keine
      Reconciliation-Datei in diesem Repo (kein Brownfield-Bootstrap).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis `BEO-PGC/<slug>/` oder eine weitere Datei in dessen
      `evidence/`; **kein Zähler wird gesetzt**, er folgt aus den Dateien.
      Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in
      §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen /
      weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      dieser Slice gehört zu
      [welle-sdk-kotlin-lh-fa-sst-009](../welle-sdk-kotlin-lh-fa-sst-009.md)
      (noch offen); die Prüfung läuft regelkonform bei deren Closure.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts` | neu | Maven-Koordinate, Kotlin-Plugin, `maven-publish`, `version = "0.1.0"` (`ADR-0109` Festlegung 3/4). |
| `sdks/kotlin/pgchangefeed-kotlin/settings.gradle.kts` | neu | Projektname, analog `examples/kotlin/settings.gradle.kts`. |
| `sdks/kotlin/pgchangefeed-kotlin/gradlew`/`gradlew.bat`/`gradle/wrapper/*` | neu | Gradle-Wrapper, dieselbe Version wie `examples/kotlin/` (real nachgemessen). |
| `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/…` (leeres/minimales Skelett, z. B. `PgChangeFeedClientOptions.kt`) | neu | gemeinsamer Konfigurations-Nenner (Adresse, Token) für die beiden Folge-Slices, falls sich beim Schreiben ein echter zeigt — kein Vorgriff auf Endpunkt-Methoden. |
| `sdks/kotlin/Dockerfile` | neu | digest-gepinnter Bau, analog `examples/kotlin/Dockerfile`s `build`-Stufe, ohne Runtime-Stufe. |
| `sdks/kotlin/pgchangefeed-kotlin/README.md` | neu | Englisch, Installationsweg **inklusive** PAT-Hinweis, Verweis auf Repo-Root-`README.md`. |
| `sdks/kotlin/.gitignore` | neu (optional) | `build/`/`.gradle/`, analog `examples/kotlin/.gitignore`. |
| Testdatei (Platzhalter, falls das Skelett bereits eine echte Klasse trägt) | neu | Konstruktions-/Validierungstest der Optionsklasse (kein Draht-Verhalten, das kommt mit den Folge-Slices). |

**Hinweise aus dem Beobachtungs-Register (vor dem ersten Bau-Lauf zu
prüfen, nicht erst nachträglich):**

- `BEO-PGC/workaround-uebersieht-etabliertes-muster-im-bestand` (Zähler
  real 2×, weiter unter der 3×-Schwelle, siehe Welle-Plan §6): Vor dem
  ersten Test-/Bau-Lauf einen Suchlauf über `sdks/csharp/**`,
  `sdks/python/**` **und** `examples/kotlin/**` fahren — hier zusätzlich
  zu den beiden Vorbild-Sprachen eine reale Kotlin/Gradle-Referenz im
  selben Repo (z. B. `grep -rl "using Xunit"`-artige Suchläufe für ein
  Kotlin/JUnit-Analogon: prüfen, welches Test-Framework/welche
  Test-Annotation `examples/kotlin/*/build.gradle.kts` bereits real
  einsetzt, bevor eine neue Wahl getroffen wird).
- **Backtick-Paritäts-Check** (`BEO-PGC/unclosed-backtick-taeuscht-nackte-id-vor`,
  Zähler real 1×): vor jedem Commit dieses Slice die Gesamtzahl der
  Backtick-Zeichen jeder geänderten Markdown-Datei zählen und gegen eine
  gerade Zahl prüfen.
- **Träger-Nachzug-Suchlauf** (`BEO-PGC/arbeit-ueberholt-stehenden-traeger`,
  Zähler real 19×, `AGENTS.md` §3.13): sobald ein Folge-Slice dieser Welle
  eine hier geschriebene Aussage überholt (z. B. „HTTP-Fläche folgt erst"
  in diesem `README.md`), deckt der Suchlauf des Folge-Slice **beide**
  Achsen ab — welche Datei (`README.md`, `build.gradle.kts`,
  `gradle.properties`) und welche Formulierung (deutsch **und** englisch).

## 4. Trigger

**Start** (`next` → `in-progress`): sofort — `ADR-0109` ist `Accepted`,
keine weitere Vorbedingung.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): entfällt aus
  heutiger Sicht — Umfang ist ein `build.gradle.kts`, ein Dockerfile, ein
  README.
- `in-progress` → `open` (blockiert — Carveout?): kein gepinntes
  `eclipse-temurin:21-jdk`-Image mit funktionierendem Kotlin-Gradle-Bau
  auffindbar — unwahrscheinlich, `examples/kotlin/Dockerfile` nutzt
  dieselbe Basis bereits real.

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + Closure-Notiz mit Lerneintrag
geschrieben.

## 6. Risiken und offene Punkte

- Das leere API-Skelett könnte beim Schreiben der beiden Folge-Slices
  (`slice-sdk-kotlin-http-client-flaeche`,
  `slice-sdk-kotlin-grpc-client-flaeche`) doch keinen echten gemeinsamen
  Nenner zeigen (unterschiedliche Auth-Form trotz beide Bearer-Token) —
  dann bleibt das Skelett minimal (nur `build.gradle.kts`/Bau/README), kein
  erzwungener gemeinsamer Code. **Ausgang:** weiter offen, entschieden
  beim Schreiben der Folge-Slices.
- Ein Kotlin-SDK-Package hat potenziell breitere Zielgruppen als ein
  Docker-only-Beispielprogramm — die Wahl des JVM-Ziel-Levels
  (`jvmToolchain`) in `build.gradle.kts` könnte den Konsumentenkreis
  unnötig auf sehr aktuelle JVM-Runtimes einschränken, falls sie unreflektiert
  von `examples/kotlin/` übernommen wird. **Ausgang:** weiter offen —
  Entscheidung bleibt bei diesem Slice (Analogie zu `examples/kotlin/`,
  keine Nutzungsdaten, die eine andere Wahl rechtfertigen); ein
  Multi-Target-Wechsel wäre eine spätere, eigenständige Entscheidung
  (`ADR-0109` §Re-Evaluierungs-Trigger 3, Gradle-Versions-getrieben).
- `examples/kotlin/*/build.gradle.kts` pinnt Fremdabhängigkeiten mit
  ausführlichen, real recherchierten Versions-Kommentaren (z. B. BOM-
  Angleichungen) — ein neues, eigenständiges `build.gradle.kts` für das
  SDK könnte dieselbe Sorgfalt unterlassen und eine ältere/inkompatible
  Kotlin-Plugin-Version wählen. **Ausgang:** weiter offen, entschieden
  beim Schreiben — die Kotlin-Plugin-Version wird zum Bau-Zeitpunkt dieses
  Slice real gegen Maven Central nachgemessen, nicht aus `ADR-0109`s
  Kontext-Messung (2026-09-17) unbesehen übernommen (`AGENTS.md` §3.12).

## 7. Closure-Notiz

- **Was hat funktioniert:** <wird beim Abschluss ergänzt>
- **Was ging anders als geplant:** <wird beim Abschluss ergänzt>
- **Steering-Loop-Eintrag:** <wird beim Abschluss ergänzt>
- **Beobachtungs-Register (`../observations/`):** <wird beim Abschluss
  ergänzt>
- **Folge-Slices:** `slice-sdk-kotlin-http-client-flaeche`,
  `slice-sdk-kotlin-grpc-client-flaeche` — beide bereits als Dateien in
  `open/` vorhanden.
- **Risiken aus §6:** <wird beim Abschluss ergänzt>
- **Drei Paarungen:** dieser Slice gehört zu
  [welle-sdk-kotlin-lh-fa-sst-009](../welle-sdk-kotlin-lh-fa-sst-009.md)
  (noch offen) — die Prüfung läuft regelkonform bei deren Closure.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Neue Sub-Area `sdks/kotlin/` —
noch nie berührt, entsteht mit diesem Slice erstmalig. Analogie zur
bereits etablierten Sub-Area `examples/kotlin/` (`ADR-0087`/`ADR-0090`,
Docker-only, digest-/paket-gepinnt) — dieselbe Konventionen-Dichte
übertragen, keine eigene Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/workaround-uebersieht-etabliertes-muster-im-bestand` (2×, siehe
§3/§6) und `BEO-PGC/unclosed-backtick-taeuscht-nackte-id-vor` (1×) betreffen
diesen Slice unmittelbar (siehe §3). Der einzige weitere inhaltlich
einschlägige Eintrag
(`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`) ist
bereits verkörpert und betrifft die Folge-Slices (HTTP-/gRPC-Fläche), nicht
dieses Gerüst-Slice ohne Betreiber-Oberfläche.

**Modus-Begründungsblock:**

### Sub-Area: `sdks/kotlin/`

- **Modus:** Greenfield — neuer Baum, Doku (`ADR-0109`) führt vor Code.
- **Konventionen-Dichte:** übernommen von `examples/kotlin/`
  (`harness/conventions.md` §Modus-Deklaration, Default-Sub-Area `*`/`PGC`
  Greenfield) — Docker-only, digest-/paket-gepinnt, kein Host-`gradle`/
  `kotlinc`/`java`.
- **Phase-Reife:** Phase 0 (Erstanlage) — Doc (`ADR-0109`) vollständig vor
  dem ersten Code dieser Sub-Area.
- **Evidenz-/Diskrepanz-Risiko:** niedrig (GF, keine Inventur nötig).
- **Reconciliation-Aufwand:** keiner — kein Brownfield-Bestand in diesem
  Baum.
