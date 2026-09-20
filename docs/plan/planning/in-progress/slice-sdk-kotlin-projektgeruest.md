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

**Verantwortlich:** Implementer-Agent, 2026-09-20.

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

- [x] `sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts` existiert:
      `group = "io.github.pt9912"`, Projektname/`artifactId`
      `pgchangefeed-kotlin`, `version = "0.1.0"` (`ADR-0109` Festlegung 4,
      Start bei `0.x.y`), Kotlin-Gradle-Plugin (dieselbe Version wie
      `examples/kotlin/build.gradle.kts`, real zum Bau-Zeitpunkt
      nachgemessen, nicht blind übernommen — `AGENTS.md` §3.12), das
      eingebaute `maven-publish`-Plugin (`ADR-0109` Festlegung 5, **kein**
      Dritt-Plugin wie `com.vanniktech.maven-publish`). Kein Import auf
      einen privaten Baum dieses Repos (`ADR-0109` §Kontext Bindung 5,
      Import-Grenze).
- [x] `sdks/kotlin/pgchangefeed-kotlin/settings.gradle.kts` existiert
      (analog `examples/kotlin/settings.gradle.kts`), eigenständiger
      Gradle-Wrapper (`gradlew`/`gradle/wrapper/` — **ohne**
      `gradlew.bat`, siehe Plan-Nachzug §3) mit derselben, bereits real
      erprobten Gradle-Version wie
      `examples/kotlin/gradle/wrapper/gradle-wrapper.properties`
      (`8.14`, real nachgemessen, `distributionSha256Sum` gepinnt).
- [x] `sdks/kotlin/Dockerfile` (Bau-Kontext `sdks/kotlin/`, eigenständig von
      `examples/kotlin/Dockerfile` und der Wurzel-`Dockerfile`) mit
      digest-gepinnter `eclipse-temurin:21-jdk`-Basis (real gemessener
      Digest zum Bau-Zeitpunkt, Kommentar-Pflicht analog
      `examples/kotlin/Dockerfile`); `./gradlew build`/`test` laufen darin,
      keine Runtime-Stufe nötig (ein SDK ist keine startbare Anwendung).
- [x] `sdks/kotlin/pgchangefeed-kotlin/README.md` (Englisch, Vorgabe aus
      dem Auftrag) beschreibt Zweck, Installationsweg (Gradle-/Maven-
      Koordinate `io.github.pt9912:pgchangefeed-kotlin`, Registry-URL
      `https://maven.pkg.github.com/pt9912/pg-change-feed`) **inklusive**
      des expliziten Hinweises auf die PAT-Pflicht zum Bezug
      (`ADR-0109` §Entscheidung Festlegung 2 — GitHub Packages verlangt
      immer eine Authentifizierung zum Lesen, auch für ein öffentliches
      Package) und verweist auf das Repo-Root-`README.md` für den vollen
      Kontext; kein Duplikat der Draht-Doku (`SPEC-018`/`SPEC-020` bleiben
      die kanonische Quelle).
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor —
      [`review-slice-sdk-kotlin-projektgeruest.md`](../../../reviews/review-slice-sdk-kotlin-projektgeruest.md)
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update für `harness/README.md` entfällt in diesem Slice — kein
      neues `make`-Target entsteht hier (Pack-Werkzeug folgt in
      `slice-sdk-kotlin-pack-werkzeug`).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben,
      **falls dieser Slice einen Inventur-Fund auflöst** — entfällt: keine
      Reconciliation-Datei in diesem Repo (kein Brownfield-Bootstrap).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis `BEO-PGC/<slug>/` oder eine weitere Datei in dessen
      `evidence/`; **kein Zähler wird gesetzt**, er folgt aus den Dateien.
      Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in
      §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen /
      weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
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

**Plan-Nachzug (Implementer-Zug, `AGENTS.md`-Konvention „im selben Lauf
nachtragen"):**

- **Kein `gradlew.bat`.** Kurskorrektur während der Umsetzung: Dieses Repo
  ist Docker-only (`AGENTS.md` §3.1) — der einzige Ort, an dem `./gradlew`
  je läuft, ist das Linux-basierte `eclipse-temurin:21-jdk`-Image in
  `sdks/kotlin/Dockerfile`. `gradlew.bat` (Windows-Batch-Pendant) wäre
  totes Gewicht ohne Verwendungszweck in diesem Repo. `examples/kotlin/`
  führt es zwar mit (Standard-`gradle wrapper`-Ergebnis), bleibt aber nach
  `ADR-0109` Festlegung 3 unangetastet — für dieses neue Projektgerüst
  gilt bewusst nur `gradlew` + `gradle/wrapper/*`.
- **Docker-Bau-Struktur mit verschachteltem Projektverzeichnis.** Anders
  als `examples/kotlin/Dockerfile` (Bau-Kontext = Projektwurzel selbst)
  liegt das Gradle-Projekt hier eine Ebene tiefer
  (`sdks/kotlin/pgchangefeed-kotlin/`, Bau-Kontext bleibt `sdks/kotlin/`)
  — dasselbe Strukturmuster wie `sdks/python/Dockerfile` (`pyproject.toml`
  unter `pgchangefeed/`, `WORKDIR`-Wechsel vor dem Bau). Die `COPY`-Zeilen
  tragen deshalb das Präfix `pgchangefeed-kotlin/`, und `WORKDIR` wechselt
  nach dem Wrapper-Kopieren auf `/src/pgchangefeed-kotlin` — ein
  Implementierungsdetail, keine Abweichung von `ADR-0109` Festlegung 3/5.
- **Reale Nachmessung ohne Drift** (`AGENTS.md` §3.12, Bau-Zeitpunkt
  2026-09-20): Kotlin-Gradle-Plugin weiterhin `2.4.20`
  (`repo1.maven.org/maven2/org/jetbrains/kotlin/kotlin-gradle-plugin/maven-metadata.xml`,
  `lastUpdated` 2026-09-07 — keine neue Version seit der `examples/kotlin/`-Messung
  vom 2026-09-17); `eclipse-temurin:21-jdk`-Digest
  (`docker buildx imagetools inspect`, linux/amd64) weiterhin
  `sha256:085eb93e049c7397f725bd8be31c4fd52ba4777a75168e851428508234aa224e`;
  Gradle-8.14-Distribution-`sha256`
  (`curl -sL https://services.gradle.org/distributions/gradle-8.14-bin.zip.sha256`)
  weiterhin `61ad310d3c7d3e5da131b76bbf22b5a4c0786e9d892dae8c1658d4b484de3caa`;
  Test-Framework-Versionen (`kotlin-test-junit5` `2.4.20`,
  `junit-platform-launcher` `6.1.3`) ebenfalls unverändert gegen Maven
  Central re-verifiziert. Alle vier Werte sind reale Wiederholungsmessungen,
  keine Übernahme aus `ADR-0109`s Kontext-Messung.
- **Test-Framework-Wahl: `kotlin-test-junit5` + JUnit Platform, kein neuer
  Kandidat.** Vor dem ersten Bau-/Testlauf geprüft (`BEO-PGC/workaround-uebersieht-etabliertes-muster-im-bestand`,
  siehe unten): `examples/kotlin/http-client/build.gradle.kts` setzt
  bereits real `org.jetbrains.kotlin:kotlin-test-junit5` +
  `org.junit.platform:junit-platform-launcher` mit `tasks.test {
  useJUnitPlatform() }` ein — dasselbe Muster übernimmt
  `sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts` unverändert, keine neue
  Framework-Wahl.
- **Der gemeinsame Nenner zeigte sich real**, analog C#/Python:
  `PgChangeFeedClientOptions` (`address: java.net.URI`,
  `apiToken: String`, `require(apiToken.isNotBlank())`) unter
  `src/main/kotlin/io/github/pt9912/pgchangefeed/`, mit drei
  Konstruktions-/Validierungstests unter `src/test/kotlin/…` (kein
  Draht-Verhalten).

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

- **Was hat funktioniert:** Vier unabhängige Docker-Builds über den
  gesamten Zyklus (Implementer, Reviewer, Verifier, und diese
  Planner-Closure-Lektüre — ohne einen fünften eigenen Bau, da Review und
  Verifikation bereits je einen eigenständigen `docker build --no-cache`
  fuhren und übereinstimmend `BUILD SUCCESSFUL` meldeten) bestätigten
  dasselbe Ergebnis, keine Divergenz. Der real nachgemessene
  Versions-/Digest-Satz (Kotlin-Gradle-Plugin `2.4.20`,
  `eclipse-temurin:21-jdk`-Digest, Gradle-8.14-`distributionSha256Sum`,
  `junit-platform-launcher` `6.1.3`) blieb über drei unabhängige
  Nachmessungen (Implementer, Reviewer, Verifier) driftfrei
  (`AGENTS.md` §3.12). Das etablierte 3-Commit-Move-Muster (`git mv`
  open→next · Inhalt · `git mv` next→in-progress) hielt `AGENTS.md` §3.3
  sauber, von Reviewer und Verifier je eigenständig bestätigt.
- **Was ging anders als geplant:** Nichts Wesentliches — der im Plan
  selbst dokumentierte Nachzug (kein `gradlew.bat`, verschachtelte
  Docker-Bau-Struktur, reale Nachmessung ohne Drift, unveränderte
  Test-Framework-Wahl) lief bereits im Implementer-Zug korrekt. Der
  Reviewer fand 0 HIGH/0 MEDIUM/3 LOW (F-1 deutsches Wortfragment
  „unstrittige" im KDoc-Kommentar von `PgChangeFeedClientOptions.kt`,
  wortgleiche Übernahme aus der C#-Fassung; F-2 deutsches Fachwort
  „vollinhalt" im README-Satz zu SSE-/NATS-Umfang, real bereits die
  dritte Instanz derselben Formulierung in `sdks/{csharp,python,kotlin}/README.md`;
  F-3 Genus-Tippfehler in der DoD-Zeile dieses Slice-Plans selbst). Als
  Planner-Entscheidung zu den drei LOW-Findings: F-3 (reine Plan-Prosa,
  keine semantische Wirkung, keine Downstream-Berührung) habe ich in
  dieser Closure direkt korrigiert (§2, „kein Runtime-Stufe" →
  „keine Runtime-Stufe") — das ist eine Bearbeitung des Plan-Dokuments
  selbst, kein neuer Implementer-/Review-/Verifikations-Zyklus. F-1/F-2
  bleiben **unverändert im ausgelieferten Code/README** — beide betreffen
  öffentlich sichtbaren, bereits Review- und Verifier-geprüften Text in
  `sdks/kotlin/pgchangefeed-kotlin/`, und der Reviewer selbst hat
  ausdrücklich **keine** Fixrunde verlangt („können bei Gelegenheit …
  mitgezogen werden, sind aber keine eigene Fixrunde wert" —
  `review-slice-sdk-kotlin-projektgeruest.md` §Verdikt). Ein isolierter
  Wort-Fix nach bereits abgeschlossener Verifikation würde einen erneuten
  Docker-Build/Review-Zyklus für eine rein stilistische Änderung ohne
  Downstream-Nutzen erzwingen; die Korrektur bleibt deshalb bewusst dem
  nächsten Slice überlassen, der dieselben Dateien ohnehin berührt
  (`slice-sdk-kotlin-http-client-flaeche`/`slice-sdk-kotlin-grpc-client-flaeche`
  schreiben in dasselbe `PgChangeFeedClientOptions.kt`-Umfeld bzw. dessen
  README-Abschnitt).
- **Steering-Loop-Eintrag:** Kein neuer Sensor — Sprachbruch in Prosa ist
  keine Gate-fähige Eigenschaft (kein Werkzeug prüft Kommentar-/
  README-Sprache, wie beide Findings selbst vermerken). Zwei neue
  Beobachtungs-Einträge (siehe unten) fassen die zugrunde liegende Klasse:
  Eine Formvorbild-Kopie zwischen SDK-Sprachpaketen trägt nicht nur
  korrekte Muster, sondern auch unbemerkte Fehler weiter — ein Review, das
  nur den Diff des neuen Pakets liest, statt das übernommene Vorbild
  erneut auf Sprachreinheit zu prüfen, findet die Übernahme nicht von
  selbst.
- **Beobachtungs-Register (`../observations/`):** Zwei neue Verzeichnisse
  angelegt:
  - `BEO-PGC/formvorbild-kopie-traegt-deutsches-wortfragment-weiter`
    (KDoc-/Doc-Kommentar-Klasse, F-1): Zähler real 2×
    (`evidence/slice-sdk-csharp-projektgeruest.md`,
    `evidence/slice-sdk-kotlin-projektgeruest.md`) — weiter unter der
    3×-Schwelle.
  - `BEO-PGC/deutsches-fachwort-im-englischen-sdk-readme` (README-Klasse,
    F-2): Zähler real 3× (`evidence/slice-sdk-csharp-projektgeruest.md`,
    `evidence/slice-sdk-python-projektgeruest.md`,
    `evidence/slice-sdk-kotlin-projektgeruest.md`) — real nachgemessen
    (`grep -n "vollinhalt" sdks/csharp/README.md sdks/python/README.md
    sdks/kotlin/pgchangefeed-kotlin/README.md`, alle drei tragen den
    Fehler) und trägt bei dieser Closure die 3×-Schwelle direkt bei ihrer
    Anlage — die csharp- und python-Instanzen waren real vorhanden, aber
    von keinem der beiden vorigen Reviews formal benannt worden. Dieser
    Eintrag ist deshalb bewusst als **erreichte** Schwelle angelegt
    (`state.md`), nicht als offen unter der Schwelle — die
    Pflege-Schärfung selbst (Klassifikation/Reviewer-Skill-Ergänzung/
    Fitness Function, `.harness/skills/reviewer.md` §Pflege) ist
    Architect-Aufgabe (Modul 4/8) und **kein** Teil dieser
    Slice-/Wellen-Closure — dieser Slice legt den Lese-Beleg nur bereit,
    er löst ihn nicht auf.
- **Folge-Slices:** `slice-sdk-kotlin-http-client-flaeche`,
  `slice-sdk-kotlin-grpc-client-flaeche` — beide bereits als Dateien in
  `open/` vorhanden.
- **Risiken aus §6:**
  - „Leeres API-Skelett könnte keinen echten gemeinsamen Nenner zeigen" —
    **Ausgang: weiter offen**, entschieden beim Schreiben der beiden
    Folge-Slices; dieser Slice liefert bereits einen realen gemeinsamen
    Nenner (`PgChangeFeedClientOptions`), ob er für beide Flächen trägt,
    zeigt sich erst dort.
  - „JVM-Ziel-Level (`jvmToolchain`) könnte Konsumentenkreis unnötig
    einschränken" — **Ausgang: weiter offen**, Entscheidung bleibt bei
    diesem Slice bestehen (`jvmToolchain(21)`, Analogie zu
    `examples/kotlin/`, keine Nutzungsdaten, die eine andere Wahl
    rechtfertigen); ein Multi-Target-Wechsel wäre eine spätere,
    eigenständige, Gradle-Versions-getriebene Entscheidung (`ADR-0109`
    §Re-Evaluierungs-Trigger 3).
  - „Neues `build.gradle.kts` könnte ältere/inkompatible Kotlin-Plugin-
    Version wählen" — **Ausgang: entfallen** — real durch drei
    unabhängige Nachmessungen (Implementer, Reviewer, Verifier) gegen
    Maven Central widerlegt: Kotlin-Gradle-Plugin `2.4.20` ist die
    tatsächlich aktuelle Version, keine Drift ggü. `ADR-0109`s
    Kontext-Messung.
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
