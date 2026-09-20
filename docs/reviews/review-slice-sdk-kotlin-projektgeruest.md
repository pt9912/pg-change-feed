# Review-Report: slice-sdk-kotlin-projektgeruest — 2026-09-20

**Review-Art:** Code — geprüft gegen Plan
(`docs/plan/planning/in-progress/slice-sdk-kotlin-projektgeruest.md`),
`ADR-0109` (Accepted) und `AGENTS.md` Hard Rules (Modul 10 §Drei
Review-Arten).

**Gegenstand:** Diff-Range `cae01ada..7500fefb` (Welle-Eröffnung bis
Implementer-Commit), Slice `slice-sdk-kotlin-projektgeruest`,
Welle `welle-sdk-kotlin-lh-fa-sst-009`.
Vier Commits: `075a61ca` (open→next, reiner Move), `1e3ecbb6`
(Verantwortlich gesetzt, nur Slice-Datei), `5d7140bc` (next→in-progress,
reiner Move), `7500fefb` (Inhalt: neuer Baum `sdks/kotlin/`).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09"
(vier repo-spezifische HIGH-Regeln, drei MEDIUM-Klassen).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-20.

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-sdk-kotlin-projektgeruest.md`
  (§1 Ziel/Abgrenzung, §2 DoD, §3 Plan inkl. Plan-Nachzug, §6 Risiken,
  §8 Sub-Area/Modus)
- `ADR-0109` (Accepted) — Festlegung 1/2/3/4/5, §Kontext Bindung 4
  (a-check liest kein Kotlin)
- `AGENTS.md` §3.1 (Docker-only), §3.3 (git mv + Inhalt = zwei Commits),
  §3.7 (Kommentar-Disziplin), §3.11 (kein host-lokaler Pfad), §3.12
  (Herkunft von Aussagen), §3.13 (Träger-Nachzug)
- `harness/conventions.md` (MR-000 ID-Schema, MR-002 Slice-Kennungen)
- `examples/kotlin/{Dockerfile,settings.gradle.kts,build.gradle.kts,http-client/build.gradle.kts}`
  als Vorbild-Referenz (`ADR-0109` §Kontext Ist-Stand)
- `sdks/csharp/**`, `sdks/python/**` als Formvorbild (Präzedenzfälle
  derselben SDK-Reihe)
- `docs/reviews/review-slice-sdk-csharp-projektgeruest.md` — Formvorbild
  für diesen Report

**Eigenständig durchgeführte Prüfungen (nicht nur Commit-Message/Plan
übernommen):**

- `docker build --no-cache -f sdks/kotlin/Dockerfile sdks/kotlin` real
  ausgeführt: Exit-Code direkt (ungepiped) geprüft, `0`. `./gradlew test`
  meldet `BUILD SUCCESSFUL`, 4 actionable tasks (inkl. `:test`, alle drei
  Testfälle in `PgChangeFeedClientOptionsTest.kt` liefen mit), `./gradlew
  build` danach ebenfalls `BUILD SUCCESSFUL`. Einzige Docker-Meldung:
  `InvalidBaseImagePlatform` (amd64-gepinntes Image auf arm64-Host) —
  identisches, bereits akzeptiertes Verhalten wie bei
  `examples/kotlin/Dockerfile`/`sdks/csharp/Dockerfile`, kein neues Muster.
  Test-Image danach mit `docker rmi` entfernt; `git status --short`
  anschließend leer.
- `make docs-check` real ausgeführt, Exit-Code direkt geprüft: `0`
  (`d-check: 848 Datei(en) geprüft, 0 Befund(e)`).
- Backtick-Paritäts-Check selbst nachgezählt (nicht dem Plan-Nachzug
  geglaubt): `sdks/kotlin/pgchangefeed-kotlin/README.md` 40 Backticks,
  `docs/plan/planning/in-progress/slice-sdk-kotlin-projektgeruest.md` 450
  Backticks — beide gerade.
- `find sdks/kotlin -iname "gradlew*"` und `git ls-files sdks/kotlin |
  grep -i gradlew` real ausgeführt: beide zeigen ausschließlich
  `sdks/kotlin/pgchangefeed-kotlin/gradlew`, kein `gradlew.bat` im
  Arbeitsbaum und nicht im Git-Index.
- Kotlin-Gradle-Plugin-Version selbst nachgemessen (nicht dem
  Implementer-Bericht geglaubt): `curl -s
  https://repo1.maven.org/maven2/org/jetbrains/kotlin/kotlin-gradle-plugin/maven-metadata.xml`
  → `<latest>2.4.20</latest>` — bestätigt, keine Drift.
- Gradle-8.14-Distribution-SHA256 selbst nachgemessen: `curl -sL
  https://services.gradle.org/distributions/gradle-8.14-bin.zip.sha256`
  → `61ad310d3c7d3e5da131b76bbf22b5a4c0786e9d892dae8c1658d4b484de3caa` —
  exakt identisch mit `distributionSha256Sum` in
  `gradle-wrapper.properties`.
- `junit-platform-launcher`-Version selbst nachgemessen (Maven Central
  `maven-metadata.xml`) → `<latest>6.1.3</latest>` — bestätigt, keine
  Drift ggü. `build.gradle.kts`.
- JDK-Basis-Digest selbst nachgemessen: `docker buildx imagetools inspect
  eclipse-temurin:21-jdk` → `linux/amd64`-Manifest-Digest
  `sha256:085eb93e049c7397f725bd8be31c4fd52ba4777a75168e851428508234aa224e`
  — identisch mit dem in `sdks/kotlin/Dockerfile` **und**
  `examples/kotlin/Dockerfile` gepinnten Wert.
- `git show 075a61ca --stat`/`git show 5d7140bc --stat` geprüft: beide
  zeigen ausschließlich den Rename `{open,next} => {next,in-progress}`,
  0 Insertions/Deletions — reine Moves. `git show 1e3ecbb6 --stat`
  geprüft: ändert ausschließlich die Slice-Datei selbst (1 Zeile) — das
  3-Commit-Muster (Move · Inhalt · Move) ist sauber getrennt, keine
  Move+Inhalt-Vermischung (`AGENTS.md` §3.3).
- `grep -rn "internal/\|cmd/pg-change-feed" sdks/kotlin/` — einziger
  Treffer liegt im Standard-Gradle-Wrapper-Skript selbst
  (`gradlew:60`, eine URL-Kommentarzeile
  `.../plugins/unixStartScript.txt`, Teil von Gradles eigenem generierten
  Wrapper, kein Bezug auf diesen Repo-Baum) — kein realer
  Import-Grenzverstoß.
- `git diff cae01ada..HEAD -- .a-check.yml examples/kotlin
  harness/README.md spec/ AGENTS.md harness/conventions.md
  docs/user/version.md` — leer; bestätigt `ADR-0109` §6 „Was diese ADR
  nicht ändert" wortgleich für diesen Slice.
- `git ls-files -s sdks/kotlin/pgchangefeed-kotlin/gradlew
  examples/kotlin/gradlew` — identischer Blob-Hash und identisches
  Executable-Bit (`100755`), also derselbe unveränderte Gradle-Wrapper.
- `grep -n "RUN \|COPY " examples/kotlin/Dockerfile` gegen
  `sdks/kotlin/Dockerfile` verglichen: dieselbe Bau-Reihenfolge
  (Wrapper/Manifeste zuerst, `gradlew help` als Cache-Aufwärmer, dann
  `src`, dann `test`, dann `build`) — kein Pfad-Zickzack, die
  „verschachtelte COPY-Struktur" (Präfix `pgchangefeed-kotlin/`,
  `WORKDIR`-Wechsel danach) ist sauber und folgt demselben Muster wie
  `sdks/python/Dockerfile`.
- `cat sdks/kotlin/pgchangefeed-kotlin/README.md` gegen die DoD-Liste
  geprüft: Englisch ✓, Installationsweg mit Maven-Koordinate und
  Registry-URL ✓, **expliziter PAT-Hinweis** vorhanden („GitHub Packages
  always requires authentication to read a package, even a public one …
  You need a GitHub account and a classic personal access token (PAT)
  with the `read:packages` scope.") ✓, Verweis auf Repo-Root-`README.md`
  ✓, kein Duplikat der Draht-Doku ✓.

---

## Findings

### F-1 — Deutsches Wortfragment in einem englischsprachigen öffentlichen KDoc-Kommentar

- `kategorie`: LOW
- `quelle`: Maintainability
- `pfad`: `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/io/github/pt9912/pgchangefeed/PgChangeFeedClientOptions.kt:9`
- `befund`: Der Klassen-Doc-Kommentar der öffentlichen Klasse
  `PgChangeFeedClientOptions` ist durchgängig Englisch, enthält aber das
  deutsche Wort „unstrittige" mitten im englischen Satz („this is the one
  unstrittige, shared configuration denominator identified while writing
  this project skeleton"). Wortgleiche Übernahme aus
  `sdks/csharp/PgChangeFeed.Client/PgChangeFeedClientOptions.cs:5-7`, dort
  bereits als F-1 (LOW, „Erstauftreten") in
  `docs/reviews/review-slice-sdk-csharp-projektgeruest.md` befundet — das
  ist damit die **zweite** formal befundete Instanz derselben Klasse
  (noch unter der 3×-Schwelle, die eine MEDIUM-Hochstufung nach §Pflege
  auslösen würde).
- `verifizierbar`: nein — reine Stilfrage, kein Gate prüft
  Kommentar-Sprache.
- `klasse`: Sprachbruch in öffentlichem Doc-Kommentar (2. Auftreten)

### F-2 — Deutsches Fachwort im öffentlichen README, dritte unbefundete Instanz

- `kategorie`: LOW
- `quelle`: Maintainability
- `pfad`: `sdks/kotlin/pgchangefeed-kotlin/README.md:9`
- `befund`: „SSE and NATS-vollinhalt delivery remain out of scope …" — das
  deutsche Wort „vollinhalt" steht unflektiert mitten im englischen Satz.
  Dieselbe Formulierung („NATS-vollinhalt(s) delivery") steht bereits in
  `sdks/python/README.md:9` und `sdks/csharp/README.md:9` — das ist real
  bereits die **dritte** Instanz derselben Sprachbruch-Klasse in dieser
  SDK-Reihe, aber keine der beiden Vorgänger-Reviews
  (`review-slice-sdk-csharp-projektgeruest.md`,
  vermutlich das Python-Pendant) hat sie bislang benannt — der Zähler für
  die §Pflege-Schärfung lief also bislang nicht mit. Dies ist die erste
  **formale** Befundung dieser Klasse.
- `verifizierbar`: nein — reine Stilfrage, kein Gate prüft
  Kommentar-/README-Sprache.
- `klasse`: Sprachbruch im öffentlichen README (Fachbegriff unflektiert,
  Erstbefundung trotz dreifachem realem Auftreten in `sdks/{csharp,python,kotlin}/README.md`)

### F-3 — Grammatik-Tippfehler in der DoD-Checkliste des Slice-Plans

- `kategorie`: LOW
- `quelle`: Maintainability
- `pfad`: `docs/plan/planning/in-progress/slice-sdk-kotlin-projektgeruest.md:86`
- `befund`: „`./gradlew build`/`test` laufen darin, kein Runtime-Stufe
  nötig" — Genus-Fehler, korrekt wäre „keine Runtime-Stufe nötig" (die
  Dockerfile-Kommentare selbst tragen die korrekte Form „Keine
  Runtime-Stufe hier"). Einmaliger Tippfehler ohne semantische
  Auswirkung.
- `verifizierbar`: nein — reine Stilfrage, kein Gate prüft Prosa-Grammatik.
- `klasse`: Tippfehler in Plan-Prosa (Einzelfall)

## Negativbefunde

- geprüft, ohne Befund: Docker-only-Disziplin (`AGENTS.md` §3.1) —
  `sdks/kotlin/Dockerfile` trägt `./gradlew`-Aufrufe ausschließlich im
  gepinnten Container; `.gitignore`-Kommentar benennt den Host-Bau
  ausdrücklich nur als Schutz gegen ein Versehen, nicht als vorgesehenen
  Pfad; kein `Makefile`-Target in diesem Slice.
- geprüft, ohne Befund: `git mv` + Inhaltsänderung als zwei/drei
  getrennte Commits (`AGENTS.md` §3.3) — siehe „Eigenständig
  durchgeführte Prüfungen" oben.
- geprüft, ohne Befund: `gradlew.bat`-Abwesenheit — weder im Arbeitsbaum
  noch im Git-Index, wie im Plan-Nachzug §3 zugesagt.
- geprüft, ohne Befund: die vier real gemessenen Versionswerte
  (Kotlin-Gradle-Plugin `2.4.20`, `eclipse-temurin:21-jdk`-Digest,
  Gradle-8.14-`distributionSha256Sum`, `junit-platform-launcher` `6.1.3`)
  — alle vier selbst gegen die genannten Live-Quellen nachgemessen, keine
  Drift ggü. dem Plan-Nachzug (`AGENTS.md` §3.12).
- geprüft, ohne Befund: `eclipse-temurin:21-jdk`-Basis statt `25-jdk`
  (`ADR-0109` Festlegung 5) — Dockerfile pinnt real `21-jdk`, keine
  Runtime-Stufe.
- geprüft, ohne Befund: eingebautes `maven-publish`-Plugin statt eines
  Dritt-Plugins (`ADR-0109` Festlegung 5, F1) — `build.gradle.kts` nutzt
  ausschließlich das Gradle-Kern-Plugin `` `maven-publish` `` (Backtick-Form
  des Plugin-Aufrufs), kein `com.vanniktech.maven-publish` oder
  vergleichbares.
- geprüft, ohne Befund: `group = "io.github.pt9912"`, `version =
  "0.1.0"`, Artefaktname `pgchangefeed-kotlin` — exakt wie in DoD und
  `ADR-0109` Festlegung 1/4 verlangt.
- geprüft, ohne Befund: Import-Grenze (`ADR-0109` §Entscheidung
  Festlegung 3, „verschärft") — kein realer Treffer für
  `internal/`/`cmd/pg-change-feed` unter `sdks/kotlin/` (der einzige
  Grep-Treffer ist Teil von Gradles eigenem generierten Wrapper-Skript,
  kein Bezug zu diesem Repo).
- geprüft, ohne Befund: `examples/kotlin/**` unangetastet (`ADR-0109`
  Festlegung 3, §6 „Was diese ADR nicht ändert") — leerer Diff.
- geprüft, ohne Befund: `.a-check.yml` unangetastet (a-check liest kein
  Kotlin, `ADR-0109` §Kontext Bindung 4).
- geprüft, ohne Befund: Kommentar-Disziplin (`AGENTS.md` §3.7) — die
  Dockerfile-/`build.gradle.kts`-/`settings.gradle.kts`-Kommentare
  beschreiben durchgängig den geltenden Zustand oder eine Abgrenzung zu
  Folge-Slices; keine Konjunktiv-Begründung über eine verworfene
  Alternative, kein abwesender Text, kein mitten im Satz abgebrochener
  Kommentar. Der bare Slice-Name im Dockerfile-Kopfkommentar
  (`slice-sdk-kotlin-projektgeruest`) folgt demselben, bereits
  akzeptierten Muster wie `examples/kotlin/Dockerfile`/
  `sdks/csharp/Dockerfile` — kein neues Muster, keine
  Vorher/Nachher-Sprache, deshalb kein HIGH nach der
  Chronik-in-Produktionscode-Regel.
- geprüft, ohne Befund: Backtick-Parität in den geänderten
  Markdown-Dateien (`README.md`, Slice-Plan) — beide gerade Zahlen.
- geprüft, ohne Befund: `README.md` dupliziert nicht die kanonische
  Draht-Doku — verweist auf `SPEC-018`/`SPEC-020` bzw. das
  Repo-Root-`README.md`, statt Endpunkte/Protokolldetails selbst zu
  beschreiben.
- geprüft, ohne Befund: host-lokale absolute Pfade (`AGENTS.md` §3.11) —
  kein Treffer unter `sdks/kotlin/` oder der Slice-Datei.
- geprüft, ohne Befund: Träger, die `ADR-0109` §6 ausdrücklich als
  „bleibt unberührt" benennt (`.a-check.yml`, `examples/kotlin/**`,
  `harness/README.md`, `spec/**`, `AGENTS.md`,
  `harness/conventions.md`, `docs/user/version.md`) — tatsächlich in
  diesem Diff unangetastet.
- geprüft, ohne Befund: Scope-Treue gegen §1 „Ausdrücklich NICHT in
  diesem Slice" — kein `make sdk-pack-kotlin`, kein Publish-Workflow,
  kein Umbau von `examples/kotlin/`, kein Pflichtenheft-Träger-Nachzug im
  Diff.
- geprüft, ohne Befund: Traceability — Commit-Betreff `7500fefb` nennt
  `LH-FA-SST-009` und `ADR-0109`, kein `SPEC-*`/`ARC-*` im Betreff;
  `harness/conventions.md` MR-002 (Slice-Kennungen sind Namen) —
  `slice-sdk-kotlin-projektgeruest` ist bereits ein Name.
- geprüft, ohne Befund: reale Docker-Build- und Testausführung —
  `BUILD SUCCESSFUL` für `test` und `build`, keine roten Tasks, die
  einzige Docker-Meldung ist der bereits bekannte plattformbedingte
  `InvalidBaseImagePlatform`-Hinweis (amd64-Pin auf arm64-Host).
- geprüft, ohne Befund: `make gates`-Nachweis — nicht separat erneut
  gefahren (dieser Review-Lauf beschränkt sich auf `make docs-check` als
  engsten einschlägigen Sensor für einen reinen neuen-Sprach-Baum-Diff,
  plus den realen Docker-Build als Gegenprobe zur
  Implementer-Behauptung); `make docs-check` real grün, siehe oben.
- geprüft, ohne Befund: DoD-Checkbox-Zustand — offene Checkboxen
  (Review, Closure-Notiz, Beobachtungs-/Reconciliation-Register, §6-
  Risiko-Ausgänge, drei Paarungen) sind tatsächlich noch offen im Plan,
  keine verfrühte Selbstabhakung.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 3 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Sprachbruch in öffentlichem
Doc-Kommentar (2. Auftreten) · Sprachbruch im öffentlichen README
(Fachbegriff unflektiert, 3. reales, aber 1. formal befundetes
Auftreten) · Tippfehler in Plan-Prosa (Einzelfall).

## Verdikt

**Merge-blockierend:** nein — 0 HIGH, 0 MEDIUM. Alle drei Findings sind
einmalige bzw. unter-der-Schwelle-liegende stilistische Beobachtungen
ohne semantische Auswirkung (LOW), keine Fixrunde erforderlich.

**Übergabe:** Keine Rückmeldung an den Implementer nötig — F-1/F-2/F-3
können bei Gelegenheit (z. B. im nächsten Slice, der dieselben Dateien
berührt) mitgezogen werden, sind aber keine eigene Fixrunde wert. F-2
verdient besondere Aufmerksamkeit beim nächsten Auftreten: Sie ist real
bereits die dritte Instanz derselben Klasse
(`sdks/csharp/README.md`, `sdks/python/README.md`, jetzt
`sdks/kotlin/README.md`) — ein viertes Auftreten (z. B. bei einem
künftigen vierten SDK) sollte die §Pflege-Schärfung auslösen
(Klassifikation prüfen / AGENTS.md-Eintrag / Fitness Function), auch
wenn die vorherigen zwei Instanzen nie formal als LOW gezählt wurden. Da
0 HIGH vorliegen (Skill-Regel „DoD-Checkbox-Nachzug ohne Fixrunde"),
ziehe ich die DoD-Checkbox „Review durchgeführt, Report unter
`docs/reviews/` liegt vor" im Slice-Plan selbst auf `[x]` nach, im
selben Commit wie diesen Report. Dieser Report ersetzt keine
Verifikation gegen die DoD — das bleibt Verifier-Aufgabe (Modul 11).
