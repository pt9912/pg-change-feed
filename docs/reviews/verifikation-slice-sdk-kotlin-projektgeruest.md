# Verifikationsbericht: slice-sdk-kotlin-projektgeruest — 2026-09-20

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?" gegen den
DoD-Vertrag
(`docs/plan/planning/in-progress/slice-sdk-kotlin-projektgeruest.md` §2)
und die §6-Risiko-Ausgänge. **Nicht** gegen den Diff als solchen
(Reviewer-Aufgabe, mit
[`review-slice-sdk-kotlin-projektgeruest.md`](review-slice-sdk-kotlin-projektgeruest.md)
bereits abgeschlossen, 0 HIGH/MEDIUM, 3 LOW) und **nicht** gegen realen
Bedarf (Validator, hier nicht ausgelöst).

**Frischer Kontext:** Diese Sitzung hat `harness/README.md`, `AGENTS.md`,
den Slice-Plan (§1/§2/§3/§6/§7), den Review-Report und `ADR-0109`
(Festlegung 1–5, §Kontext, §Konsequenzen) gelesen. Behauptungen aus
Slice-Plan, Commit-Message und Review-Report wurden **nicht** übernommen,
sondern eigenständig nachgemessen: eigener realer `docker build
--no-cache`-Lauf, eigener `find`/`git ls-files`-Lauf gegen `gradlew.bat`,
eigene Live-Abfragen gegen Maven Central, `services.gradle.org` und die
Docker-Registry (drei unabhängige Versions-/Digest-Reverifikationen —
Implementer maß sie, Reviewer maß sie erneut, dies ist die dritte
Messung), eigener `grep` gegen `internal/`/`cmd/pg-change-feed`, eigener
`git diff`-Lauf gegen die von `ADR-0109` §6 als unberührt deklarierten
Träger, und ein eigener, ungepipter `make gates`-Lauf mit direkter
Exit-Code-Prüfung (`AGENTS.md` §3.9).

**Gegenstand:** `slice-sdk-kotlin-projektgeruest`, fünf Commits auf
`main` (`075a61ca`, `1e3ecbb6`, `5d7140bc`, `7500fefb`, `75734f84`). Der
Slice liegt weiterhin in `in-progress/` — erwartungsgemäß: Closure (§7
Notiz, Beobachtungs-Register, Risiko-Ausgänge im Plan selbst, `git mv`
nach `done/`) ist nicht Gegenstand dieser Verifikation.

---

## 1. DoD-Zeile für Zeile gegen den tatsächlichen Baum geprüft

| DoD-Zeile (§2) | Eigene Prüfung | Ergebnis |
|---|---|---|
| `build.gradle.kts`: `group = "io.github.pt9912"`, Projektname/`artifactId` `pgchangefeed-kotlin`, `version = "0.1.0"`, Kotlin-Gradle-Plugin, eingebautes `maven-publish` | `Read` der Datei: `group = "io.github.pt9912"` (Z. 25), `version = "0.1.0"` (Z. 26, zusätzlich in der `publishing`-Block-Koordinate Z. 46), `plugins { kotlin("jvm") version "2.4.20"; `maven-publish` }` (Z. 20–23) — kein Dritt-Plugin, `MavenPublication`-Block trägt `artifactId = "pgchangefeed-kotlin"` (Z. 45) | ✓ erfüllt |
| Kein Import auf einen privaten Baum | eigener `grep -rn "internal/\|cmd/pg-change-feed" sdks/kotlin/` | einziger Treffer: `gradlew:60`, eine URL-Kommentarzeile aus Gradles eigenem generiertem Wrapper (`.../gradle/api/internal/plugins/unixStartScript.txt`) — kein Bezug zu diesem Repo. ✓ erfüllt |
| `settings.gradle.kts` existiert, eigenständiger Wrapper (`gradlew`/`gradle/wrapper/`), **ohne** `gradlew.bat`, Gradle `8.14`, `distributionSha256Sum` gepinnt | `Read` bestätigt Wurzelprojekt-Name `pgchangefeed-kotlin`; `gradle-wrapper.properties` trägt `distributionUrl=…gradle-8.14-bin.zip`, `distributionSha256Sum=61ad3…caa` | ✓ erfüllt (Wrapper-Abwesenheit siehe §2 unten) |
| `sdks/kotlin/Dockerfile`, Bau-Kontext `sdks/kotlin/`, digest-gepinnte `eclipse-temurin:21-jdk`-Basis, `./gradlew build`/`test`, keine Runtime-Stufe | `Read`: `FROM eclipse-temurin:21-jdk@sha256:085eb93e…224e AS build` (einzige Stufe), `RUN ./gradlew --no-daemon test` (Z. 44), `RUN ./gradlew --no-daemon build` (Z. 45) | ✓ erfüllt |
| `README.md` Englisch, Installationsweg (Maven-Koordinate, Registry-URL), **inklusive** PAT-Hinweis, Verweis auf Repo-Root-`README.md`, kein Duplikat der Draht-Doku | `Read`: durchgängig Englisch; „This package is distributed via GitHub Packages …"; „GitHub Packages always requires authentication to read a package, even a public one … a classic personal access token (PAT) with the `read:packages` scope"; Verweis auf `github.com/pt9912/pg-change-feed/blob/main/README.md`; §Documentation verweist explizit auf `spec/pflichtenheft.md`/`SPEC-018`/`SPEC-020` statt zu duplizieren | ✓ erfüllt |
| `make gates` grün | eigener Lauf, siehe §5 unten | ✓ erfüllt |
| Review durchgeführt, Report liegt vor | `docs/reviews/review-slice-sdk-kotlin-projektgeruest.md` existiert, Verdikt „Merge-blockierend: nein — 0 HIGH, 0 MEDIUM" | ✓ erfüllt |
| Doku-Update `harness/README.md` entfällt | kein neues `make`-Target in diesem Diff (kein `sdk-pack-kotlin`, kein Publish-Workflow) | ✓ zutreffend begründet |
| Closure-Notiz / Reconciliation / Beobachtungs-Register / §6-Risiko-Ausgänge / drei Paarungen | §7 des Plans trägt weiterhin ausschließlich Platzhaltertext (`<wird beim Abschluss ergänzt>`); §6-Risiken tragen im Fließtext bereits „Ausgang: weiter offen", aber keine DoD-Checkbox ist auf `[x]` gesetzt | erwartungsgemäß offen — **kein Mangel**, siehe §6 unten |

## 2. `gradlew.bat`-Abwesenheit — dritte unabhängige Bestätigung

Eigene Läufe:

```
find sdks/kotlin -iname "gradlew*"
git ls-files sdks/kotlin | grep -i gradlew
```

Beide Befehle liefern ausschließlich `sdks/kotlin/pgchangefeed-kotlin/gradlew`
— kein `gradlew.bat` im Arbeitsbaum und keiner im Git-Index. Deckungsgleich
mit Implementer-Zusage (Plan-Nachzug §3) und Reviewer-Befund.

## 3. Reale, eigenständige Docker-Bau-Ausführung — dritte unabhängige Bau-Bestätigung

Eigener Lauf (nach Implementer und Reviewer die dritte unabhängige
Bau-Bestätigung):

```
docker build --no-cache -f sdks/kotlin/Dockerfile sdks/kotlin -t verifier-kotlin-sdk-check:tmp
```

Exit-Code direkt (ungepiped) geprüft: **`0`**. `./gradlew --no-daemon test`
meldet `BUILD SUCCESSFUL in 26s`, 4 actionable tasks (alle drei Testfälle aus
`PgChangeFeedClientOptionsTest.kt` liefen mit den erwarteten Testnamen mit);
`./gradlew --no-daemon build` danach ebenfalls `BUILD SUCCESSFUL in 9s`, 4
actionable tasks up-to-date. Einzige Docker-Meldung:
`InvalidBaseImagePlatform` (amd64-gepinntes Image auf arm64-Host) — dasselbe,
bereits akzeptierte Verhalten wie bei
`examples/kotlin/Dockerfile`/`sdks/csharp/Dockerfile`, kein neues Muster.
Test-Image danach per `docker rmi` entfernt; `git status --short` anschließend
leer — kein Nebeneffekt im Arbeitsbaum.

## 4. Drei Versions-/Digest-Reverifikationen — jeweils dritte unabhängige Messung

| Wert | Live-Quelle (eigener Abruf) | Ergebnis | Deckung mit Dockerfile/`build.gradle.kts`/`gradle-wrapper.properties` |
|---|---|---|---|
| Kotlin-Gradle-Plugin | `curl -s https://repo1.maven.org/maven2/org/jetbrains/kotlin/kotlin-gradle-plugin/maven-metadata.xml` | `<latest>2.4.20</latest>`, `<lastUpdated>20260907091159</lastUpdated>` | ✓ identisch mit `kotlin("jvm") version "2.4.20"` |
| Gradle-8.14-Distribution-SHA256 | `curl -sL https://services.gradle.org/distributions/gradle-8.14-bin.zip.sha256` | `61ad310d3c7d3e5da131b76bbf22b5a4c0786e9d892dae8c1658d4b484de3caa` | ✓ byte-identisch mit `distributionSha256Sum` |
| `eclipse-temurin:21-jdk` linux/amd64-Digest | `docker buildx imagetools inspect eclipse-temurin:21-jdk` (raw Manifest-Liste) | `sha256:085eb93e049c7397f725bd8be31c4fd52ba4777a75168e851428508234aa224e` (Platform `linux/amd64`) | ✓ identisch mit dem in `sdks/kotlin/Dockerfile` gepinnten `FROM`-Digest |

Alle drei Werte sind heute (2026-09-20) real und unabhängig — zum dritten
Mal in dieser Kette (Implementer → Reviewer → Verifier) — nachgemessen,
keine Drift gegenüber dem committeten Stand.

## 5. `make gates` real ausgeführt

Eigener, ungepipter Lauf, Exit-Code direkt geprüft (`AGENTS.md` §3.9):

```
make gates > /tmp/verifier-make-gates.log 2>&1; ec=$?
```

Ergebnis: **`EXIT=0`**. Einzelbelege aus demselben Lauf:

| Gate | Ergebnis |
|---|---|
| `baseline-verify` | `v6.9.0 OK — 54 Dateien (Integritaet + Vollstaendigkeit, netzlos)` |
| `docs-check` | `d-check: 849 Datei(en) geprüft, 0 Befund(e)` |
| `commit-traceability` | `d-check: 849 Datei(en) geprüft, 0 Befund(e)` (Modul `commits`) + `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` |
| `coverage-gate` | `coverage-gate: OK — Coverage 82.80% erfüllt Schwelle 80%` |
| `generated-sync` | `generated-sync: OK — das committete Erzeugnis ist byte-gleich der Ausgabe des gepinnten Generators (Stufe proto)` |
| `a-check` | `gesamt: 0 Befund(e)` |

Die geprüfte Commit-Range (`HEAD~5..HEAD`) deckt exakt die fünf Commits
dieses Slice — `commit-traceability` prüft damit denselben Umfang, den
dieser Bericht als Gegenstand führt.

## 6. §6-Risiken — korrekt NICHT als erledigt markiert?

| # | Risiko (Kurzform) | Zustand im Plan | Zulässig für diesen Bearbeitungsstand? |
|---|---|---|---|
| 1 | Leeres API-Skelett könnte keinen echten gemeinsamen Nenner zeigen | Fließtext trägt bereits „Ausgang: weiter offen" — inhaltlich korrekt, denn `PgChangeFeedClientOptions` zeigte sich real als gemeinsamer Nenner (§3 Plan-Nachzug) | ✓ zulässig — kein DoD-Häkchen vorgezogen |
| 2 | JVM-Ziel-Level (`jvmToolchain`) könnte Konsumentenkreis unnötig einschränken | „Ausgang: weiter offen" — `build.gradle.kts` trägt real `jvmToolchain(21)`, Entscheidung bleibt bei diesem Slice, keine andere Wahl nachträglich behauptet | ✓ zulässig |
| 3 | Neues `build.gradle.kts` könnte ältere/inkompatible Kotlin-Plugin-Version wählen | „Ausgang: weiter offen" — real durch §4 dieses Berichts widerlegt (keine Drift), aber die Plan-Zeile behauptet das nicht vorschnell als „entfallen", sondern lässt den Ausgang textuell offen | ✓ zulässig |

Die zugehörige DoD-Checkbox „Jedes Risiko aus §6 trägt einen Ausgang
(eingetreten / entfallen / weiter offen)" bleibt im Plan `[ ]` — korrekt,
denn ein formaler Nachzug in die Checkliste selbst ist noch nicht erfolgt,
obwohl der Fließtext bereits einen Ausgang trägt. Keine der drei
Risiko-Beschreibungen wurde als „erledigt"/„entfallen" fehlklassifiziert,
um die Closure vorzuziehen — dieser Bearbeitungsstand ist ehrlich
dargestellt.

## 7. `git diff` gegen die von `ADR-0109` §6 als unberührt deklarierten Träger

Eigener Lauf:

```
git diff cae01ada..HEAD --stat -- examples/kotlin .a-check.yml spec/architecture.md docs/user/version.md sdks/csharp sdks/python
```

Ergebnis: **leerer Diff** — kein einziges Byte in einem der sechs
genannten Pfade geändert. Deckungsgleich mit `ADR-0109` §6 „Was diese ADR
nicht ändert" und mit dem Review-Befund.

## 8. Reconciliation

Entfällt — kein Brownfield-Bootstrap in diesem Repo, keine
`reconciliation.md`-Datei existiert (`harness/conventions.md` §Modus-
Deklaration: gesamtes Repo Greenfield). DoD-Zeile korrekt als „entfällt"
geführt.

---

## Offene Closure-Arbeit (erwartet, kein Befund)

- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (§7 des Plans ist noch
  Platzhaltertext).
- [ ] Beobachtungs-Register fortgeschrieben.
- [ ] §6-Risiko-Ausgänge formal in die DoD-Checkliste nachgezogen (Inhalt
  bereits im Fließtext vorhanden, siehe §6 oben).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) — Prüfung läuft
  regelkonform bei Closure der Welle
  [welle-sdk-kotlin-lh-fa-sst-009](../plan/planning/done/welle-sdk-kotlin-lh-fa-sst-009.md).
- `git mv` nach `done/` — nicht erfolgt, erwartungsgemäß.

## Verdikt

**DoD-Konformität bestätigt.** Jede DoD-Zeile aus §2 wurde einzeln gegen
den tatsächlichen Baum geprüft, nicht aus Plan, Commit-Message oder
Review-Bericht übernommen:

1. `build.gradle.kts`/`settings.gradle.kts`/Gradle-Wrapper tragen exakt
   die verlangten Werte (`group`, `artifactId`, `version`, Kotlin-Plugin,
   eingebautes `maven-publish`).
2. `sdks/kotlin/Dockerfile` ist digest-gepinnt auf `eclipse-temurin:21-jdk`
   (real re-verifiziert, nicht `25-jdk`), ohne Runtime-Stufe.
3. `README.md` ist Englisch und trägt den PAT-Hinweis für GitHub Packages.
4. `make gates` lief eigenständig, ungepiped, mit `EXIT=0` über alle
   sechs Gates.
5. Der eigene, dritte Docker-Build (`--no-cache`) liefert Exit `0` und
   zweimal `BUILD SUCCESSFUL` (Test, Build) — dieselbe Bestätigung wie
   Implementer und Reviewer, unabhängig wiederholt.
6. `gradlew.bat` ist real weder im Arbeitsbaum noch im Git-Index — dritte
   unabhängige Bestätigung.
7. Alle drei Versions-/Digest-Werte (Kotlin-Gradle-Plugin, Gradle-8.14-
   `distributionSha256Sum`, `eclipse-temurin:21-jdk`-Digest) sind real
   gegen ihre Live-Quellen nachgemessen — dritte unabhängige Messung ohne
   Drift.
8. Kein Import aus `internal/**`/`cmd/**` unter `sdks/kotlin/` (einziger
   Grep-Treffer ist Teil von Gradles eigenem Wrapper-Skript).
9. `examples/kotlin/**`, `.a-check.yml`, `spec/architecture.md`,
   `docs/user/version.md`, `sdks/csharp/**`, `sdks/python/**` sind über den
   gesamten Slice-Range unverändert (leerer Diff).
10. Review durchgeführt, 0 HIGH/MEDIUM, 3 LOW ohne Fixrunden-Pflicht —
    Checkbox berechtigt gesetzt.
11. Reconciliation entfällt korrekt (kein Brownfield-Bestand).
12. §6-Risiken tragen inhaltlich einen zulässigen Ausgang, ohne die
    zugehörige DoD-Checkbox vorzeitig abzuhaken.

**Keine offenen Punkte, die der DoD-Konformität dieses Diffs
entgegenstehen.**

**Freigabe an den Planner:** Der Slice ist bereit für Closure nach
`done/`, sobald Closure-Notiz, Beobachtungs-Register-Eintrag, der formale
Risiko-Ausgangs-Nachzug in die DoD-Checkliste und die drei Paarungen
(letztere bei Wellen-Closure) ergänzt sind.
