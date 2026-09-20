# Verifikations-Report: slice-sdk-kotlin-http-client-flaeche — 2026-09-20

**Rolle:** Verifier (Modul 11) — Frage: „Bauen wir es richtig?" (DoD- und
Entscheidungs-Konformität), nicht Review (Diff gegen Plan/Hard Rules) und
nicht Validierung (Bauen wir das Richtige?).

**Gegenstand:** `docs/plan/planning/in-progress/slice-sdk-kotlin-http-client-flaeche.md`
(`LH-FA-SST-009`, `ADR-0109`), Diff `df6acd80..9181ae21` — Implementer-Commit
`7d01359b` (feat(sdk): Kotlin-HTTP-Client-Fläche), Review-Commit `ab3f9485`
(1 HIGH F-1), Fix-Commit `9181ae21` (KDoc + Plan-Nachzug korrigiert,
DoD-Checkbox „Review durchgeführt" gesetzt).

**Eingangs-Kontext (gelesen, nicht übernommen):** `harness/README.md`,
`AGENTS.md` (§3.1, §3.7, §3.9, §3.12, §3.13), `harness/conventions.md`,
der vollständige Slice-Plan (§1 Ziel/Abgrenzung, §2 DoD, §3 Plan inkl.
Plan-Nachzug, §6 Risiken), `ADR-0109` (Accepted, vollständig gelesen inkl.
§Entscheidung Festlegung 1–6, §Verglichene Alternativen, §Konsequenzen),
`spec/pflichtenheft.md` §`SPEC-018`/`SPEC-022`,
`docs/reviews/review-slice-sdk-kotlin-http-client-flaeche.md`.

**Modell:** claude-sonnet-5 · **Datum:** 2026-09-20.

---

## 1. DoD-Zeilen einzeln gegen Code-/Doku-Stand geprüft (§2)

Jede Zeile unten wurde **selbst nachgemessen**, nicht aus Implementer-
oder Reviewer-Bericht übernommen.

### 1.1 Öffentliche Client-Klasse, alle neun `SPEC-018`-Fähigkeiten + `readChanges` (`SPEC-022`)

Gelesen: `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/io/github/pt9912/pgchangefeed/http/PgChangeFeedHttpClient.kt`
Zeile für Zeile gegen `spec/pflichtenheft.md` `SPEC-018` (Zeilen 340–380) und
`SPEC-022` (Zeilen 464–492) gehalten:

| Fähigkeit | Endpunkt/Methode Spec | Endpunkt/Methode Kotlin | Response-Feldnamen | Ergebnis |
|---|---|---|---|---|
| `RegisterConsumer` | `POST /consumers` | `registerConsumer()` → `post("/consumers", …)` | `consumer_id`/`name`/`already_registered` | deckungsgleich |
| `AcknowledgeConsumer` | `POST /consumers/acknowledge` | `acknowledgeConsumer()` | `consumer_id`/`source_id`/`offset` | deckungsgleich |
| `GetConsumerPosition` | `GET /consumers/position?consumer_id=` | `getConsumerPosition()` | `+acknowledged` | deckungsgleich |
| `RemoveConsumer` | `POST /consumers/remove` | `removeConsumer()` | `consumer_id`/`removed` | deckungsgleich |
| `EnableTable` | `POST /tables/enable`, 7 Pflichtfelder | `enableTable()`, `EnableTableRequest` 7 Felder | `table_id`/`source`/`schema`/`table`/`already_enabled` | deckungsgleich |
| `DisableTable` | `POST /tables/disable`, 4 Pflichtfelder | `disableTable()`, `DisableTableRequest` 4 Felder | `removed`/`retained` | deckungsgleich |
| `GetStatus` | `GET /tables/status` 4 Query-Parameter | `getStatus()` | `enabled`/`retained` | deckungsgleich |
| `ListTables` | `GET /tables` 2 Query-Parameter | `listTables()` | `tables`/`retained` je Liste | deckungsgleich |
| `RunRetention` | `POST /retention/run` | `runRetention()` | `deleted` | deckungsgleich |
| `readChanges` (`SPEC-022`) | `GET /changes`, `source` Pflicht, `schema`/`table`/`from`/`to`/`limit` optional, `from` inklusiv/`to` exklusiv | `readChanges()` — identische Parameterliste, gleiche Semantik im KDoc benannt | 10-Felder-`Change` (`commit_position`…`committed_at`) | deckungsgleich |

Fehler-Antwortform `400`/`401`/`403`/`404`/`500` → `PgChangeFeedException`-
Sieben-Klassen-Hierarchie (`PgChangeFeedException.kt`), Bearer-Token bei
Konstruktion übergeben (`PgChangeFeedClientOptions`, kein globaler State) —
bestätigt. **Ergebnis: DoD-Zeile erfüllt.**

### 1.2 Tests: Happy Path + Auth-Boundary je Fähigkeit, netzlos

Gelesen: alle fünf Testdateien
(`PgChangeFeedHttpClientConsumerTest.kt`, `…TableTest.kt`,
`…RetentionAndChangesTest.kt`, `…AuthBoundaryTest.kt`,
`PgChangeFeedClientOptionsTest.kt`) plus die zwei Test-Helfer
(`FakeHttpTransport.kt`, `TestClientFactory.kt`).

- Happy Path je Fähigkeit real vorhanden (Consumer-Vier, Tabellen-Vier,
  Retention + readChanges) — jeweils Methode/URL/Body/Response-Parsing
  geprüft.
- Auth-Boundary: `missing or unknown token throws Unauthorized` (401),
  `reader token against an admin endpoint throws Forbidden` (403), plus
  400/404/500/`UnexpectedStatus`/`MalformedResponse`-Fälle in
  `PgChangeFeedHttpClientAuthBoundaryTest.kt`.
- Netzlosigkeit real geprüft: `TestClientFactory` verdrahtet
  `PgChangeFeedHttpClient` über den `internal`-Sekundärkonstruktor direkt
  mit `FakeHttpTransport` (In-Memory-Antwortfunktion); die Basis-Adresse
  `http://example.invalid:8080` wird nie aufgelöst, weil `FakeHttpTransport.send()`
  nie an `java.net.http.HttpClient` delegiert — kein Socket in keinem der
  fünf Testdateien. **Ergebnis: DoD-Zeile erfüllt.**

### 1.3 Kein Import aus `internal/**`/`cmd/**`

Eigener Lauf:

```
grep -rn "internal/\|cmd/" sdks/kotlin/
→ sdks/kotlin/pgchangefeed-kotlin/gradlew:60 (Kommentar auf eine
  fremde GitHub-URL im Gradle-Wrapper-Skript, kein Import dieses Repos)
```

Kein Treffer im SDK-Code selbst. **Ergebnis: DoD-Zeile erfüllt.**

### 1.4 `make gates` grün

Eigener, ungepipter Lauf (`AGENTS.md` §3.9), Exit-Code direkt geprüft:

```
make gates > /tmp/verifier-gates.log 2>&1; ec=$?; echo "EXIT_CODE=$ec"
EXIT_CODE=0
```

Belege aus dem Log: `baseline-verify: v6.9.0 OK`, `coverage-gate: OK —
Coverage 82.70% erfüllt Schwelle 80%`, `d-check: 860 Datei(en) geprüft, 0
Befund(e)` (zweimal — docs-check-Modulliste), `commit-traceability: OK —
5 Commit(s) in "HEAD~5..HEAD"`, `generated-sync: OK`, `a-check … gesamt: 0
Befund(e)`. **Ergebnis: DoD-Zeile erfüllt.**

### 1.5 Review durchgeführt, F-1 real aufgelöst

`docs/reviews/review-slice-sdk-kotlin-http-client-flaeche.md` fand 1 HIGH
(F-1): Kotlin-`internal` fälschlich als absolute JVM-Unsichtbarkeit
behauptet in `HttpTransport.kt`s KDoc und im Plan-Nachzug §3, real mit
`javap -p` widerlegt (Konstruktor und Interface kompilieren zu
gewöhnlichen `public`-Symbolen).

Fix-Commit `9181ae21` gelesen (voller Diff): beide Träger korrigiert.
Sachliche Prüfung der neuen Formulierung (nicht nur, dass sie geändert
wurde):

- `HttpTransport.kt` neue KDoc: „`internal` here is a compile-time
  visibility boundary that the Kotlin compiler enforces against other
  Kotlin modules' metadata; … it is not a JVM bytecode access restriction.
  … both this constructor and [HttpTransport] itself are ordinary `public`
  symbols … a Java caller, or reflection from any language, can still
  implement [HttpTransport] and invoke this constructor directly." —
  **sachlich korrekt**: Kotlins `internal`-Modifikator wird vom
  Kotlin-Compiler-Frontend über `@Metadata`/Modul-Namen-Mangling
  durchgesetzt (`internal`-Member erhalten einen modul-spezifischen
  Namens-Suffix im Bytecode-Namen, keine JVM-`private`/`package`-Grenze);
  Konstruktoren heißen im Bytecode immer `<init>` und Interfaces werden
  nicht gemangelt — beides bleibt gewöhnlicher `public`-Bytecode, exakt
  wie im Fix beschrieben. Kein Widerspruch zur JVM-Spezifikation, keine
  Übertreibung in die Gegenrichtung (der Fix behauptet nicht „unwirksam",
  sondern benennt korrekt Zweck und Grenze: Schutz gegen versehentliche
  Kotlin/Gradle-Modul-Nutzung, nicht gegen jeden JVM-Aufrufer).
- Plan-Nachzug §3: derselbe Wortlaut, konsistent mit der KDoc-Fassung.
- Eigener `grep`-Lauf nach dem alten, falschen Wortlaut:
  ```
  grep -rn "invisible outside\|no public API surface\|nur innerhalb des Gradle-Moduls" sdks/kotlin/ docs/plan/planning/in-progress/slice-sdk-kotlin-http-client-flaeche.md
  → kein Treffer
  ```
  Kein Restvorkommen der widerlegten Aussage. Ein zusätzlicher eigener
  `javap -p`-Lauf wurde als nicht nötig bewertet — die Textlogik des Fixes
  ist in sich korrekt und deckt sich mit der bereits im Review protokollierten
  `javap`-Ausgabe (Konstruktor-Signatur, Interface-Deklaration), die dieser
  Verifikationslauf im Report gegengelesen hat.
- DoD-Checkbox „Review durchgeführt … kein offenes HIGH" ist im Fix-Commit
  auf `[x]` gesetzt, mit einer neuen Erläuterungszeile, die F-1 und seine
  Auflösung benennt — Checkbox-Zustand und Text stimmen überein.

**Ergebnis: DoD-Zeile erfüllt — F-1 ist inhaltlich korrekt aufgelöst, kein
offenes HIGH.**

### 1.6 Doku-Update `docs/user/benutzerhandbuch.md`

Gelesen: §4 „Zugriff über die HTTP-/JSON-API", `**SDK:**`-Absatz. Der
Kotlin-Absatz ist vorhanden, benennt Koordinate
(`io.github.pt9912:pgchangefeed-kotlin`), die zehn gedeckten Fähigkeiten,
die sealed-class-Fehlerhierarchie, den Vertriebsweg **GitHub Packages**
inklusive explizitem PAT-Hinweis (`read:packages`-Scope, „GitHub Packages
verlangt immer eine Authentifizierung zum Lesen, auch für ein öffentliches
Package") — deckt `ADR-0109` Festlegung 2 vollständig.
Versionshistorie: `Version: 1.36` im Kopf, Zeile „1.36 | 2026-09-20 |
Kotlin-SDK-Hinweis für die HTTP-Oberfläche ergänzt …" vorhanden.
**Ergebnis: DoD-Zeile erfüllt.**

---

## 2. Eigener Docker-Build (dritte unabhängige Bau-Bestätigung)

```
docker build --no-cache -f sdks/kotlin/Dockerfile sdks/kotlin
```

Ergebnis: beide Gradle-Aufrufe im Dockerfile grün —
`./gradlew --no-daemon test` → `BUILD SUCCESSFUL in 34s, 4 actionable
tasks: 4 executed`; `./gradlew --no-daemon build` → `BUILD SUCCESSFUL in
11s`. Exportiertes Image erfolgreich gebaut, Shell-Exit `0`. Dies ist die
dritte unabhängige Bau-Bestätigung nach Implementer (Plan-Nachzug §3
Mutation-Beleg) und Reviewer (Cache-Hit-Build) — hier mit `--no-cache`,
also ein vollständig frischer Bau ohne Layer-Wiederverwendung.

**Ergebnis: bestätigt.**

## 3. Diff-Scope-Prüfung gegen unbeabsichtigte Berührung

```
git diff --stat df6acd80..HEAD -- examples/kotlin/ .a-check.yml \
  spec/architecture.md docs/user/version.md sdks/csharp/ sdks/python/
→ (leer)
```

Keine dieser sechs Pfad-Gruppen wurde berührt — deckt sich mit `ADR-0109`
§6 „Was diese ADR nicht ändert" (Punkte 2/3/5/7). Vollständiger Diff
(`git diff --stat df6acd80..HEAD`) umfasst genau 19 Dateien, alle unter
`sdks/kotlin/**`, `docs/plan/planning/in-progress/slice-sdk-kotlin-http-client-flaeche.md`,
`docs/reviews/review-slice-sdk-kotlin-http-client-flaeche.md` und
`docs/user/benutzerhandbuch.md` — kein Fremdtreffer.

**Ergebnis: bestätigt, keine unbeabsichtigte Berührung.**

## 4. Backtick-Parität — eigenständig nachgezählt

Alle vier im Diff (`df6acd80..HEAD`) geänderten Markdown-Dateien, Backtick-
Zeichen je Datei gezählt (`tr -cd '`' | wc -c`):

| Datei | Anzahl Backticks |
|---|---|
| `docs/plan/planning/in-progress/slice-sdk-kotlin-http-client-flaeche.md` | 386 (gerade) |
| `docs/reviews/review-slice-sdk-kotlin-http-client-flaeche.md` | 326 (gerade) |
| `docs/user/benutzerhandbuch.md` | 1742 (gerade) |
| `sdks/kotlin/pgchangefeed-kotlin/README.md` | 58 (gerade) |

Alle vier Dateien tragen eine gerade Anzahl Backtick-Zeichen — keine
unclosed Inline-Code-Spanne. **Ergebnis: Parität bestätigt.**

## 5. §6-Risiken — korrekt nicht als erledigt markiert

§6 trägt drei Risiken, jedes mit Prosa-„**Ausgang:**"-Feld, keine
Checkbox auf Risiko-Ebene. Die zugehörige DoD-Sammelzeile in §2,
„Jedes Risiko aus §6 trägt einen Ausgang.", ist real weiterhin `[ ]`
(unchecked) — korrekt: zwei der drei Risiken sind im Text explizit als
„weiter offen" geführt (Test-Fake-Realitätsnähe, `HttpClient`-
Ersetzbarkeit — Letzteres tatsächlich beim Schreiben gelöst, aber die
Formulierung selbst hält den Prosa-Stand fest, keine Checkbox-Aktualisierung
vorgenommen), das dritte ist im Fließtext als entschieden vermerkt. Diese
DoD-Zeile bündelt formal alle drei — ihr Ausgang bleibt Planner-
Closure-Arbeit (§7 „Risiken aus §6" ist noch als Platzhalter
`<wird beim Abschluss ergänzt>` offen), nicht Implementer- oder
Verifier-Aufgabe. Ebenso korrekt offen: „Closure-Notiz" (§7 vollständig
Platzhalter), „Beobachtungs-Register fortgeschrieben", „Die drei
Paarungen sind getragen".

**Ergebnis: korrekt unerledigt — kein Verstoß, kein vorzeitiger Closure-
Versuch.**

---

## Verdikt

**DoD-konform: ja.**

Alle Implementer-/Reviewer-skopierten DoD-Zeilen in §2 sind real erfüllt
und wurden in diesem Lauf **eigenständig** nachgemessen (keine Übernahme
der Implementer-/Reviewer-Behauptung): Fähigkeitsdeckung Feld für Feld
gegen `SPEC-018`/`SPEC-022` gehalten, Tests real gelesen (Happy Path +
Auth-Boundary, netzlos verifiziert über den `internal`-Transport-Seam),
Import-Grenze per eigenem `grep` bestätigt, `make gates` eigenständig und
ungepiped mit Exit `0` gefahren, F-1s Fix inhaltlich (nicht nur formal)
als korrekt beurteilt, Doku-Update mit PAT-Hinweis vorhanden. Ein dritter,
unabhängiger `docker build --no-cache`-Lauf bestätigt den Bau zusätzlich.
Diff-Scope-Prüfung zeigt keine unbeabsichtigte Berührung von
`examples/kotlin/**`, `.a-check.yml`, `spec/architecture.md`,
`docs/user/version.md`, `sdks/csharp/**`, `sdks/python/**`. Backtick-
Parität aller vier geänderten Markdown-Dateien bestätigt. Die noch offenen
DoD-Zeilen (Closure-Notiz, Beobachtungs-Register, §6-Risiko-Ausgänge, drei
Paarungen) sind korrekt als Planner-Closure-Arbeit unerledigt — kein
DoD-Verstoß, sondern der erwartete Zwischenstand vor Closure.

**Übergabe:** Bestätigung an den Planner — dieser Slice ist bereit für die
Closure-Schritte (§7, Beobachtungs-Register, §6-Risiko-Ausgänge, drei
Paarungen); dieser Report selbst nimmt keine dieser Handlungen vor.

## Ausgeführte Sensor-Läufe dieses Verifikationslaufs (Zusammenfassung)

- `make gates` — eigenständig, ungepiped, Exit `0`.
- `docker build --no-cache -f sdks/kotlin/Dockerfile sdks/kotlin` —
  eigenständig, Exit `0`, `BUILD SUCCESSFUL` für `test` und `build`.
- `grep -rn "internal/\|cmd/" sdks/kotlin/` — eigenständig, kein
  Import-Treffer.
- `grep -rn "invisible outside\|no public API surface\|nur innerhalb des
  Gradle-Moduls" sdks/kotlin/ docs/plan/planning/in-progress/slice-sdk-kotlin-http-client-flaeche.md`
  — eigenständig, kein Treffer (F-1-Restvorkommen ausgeschlossen).
- `git diff --stat df6acd80..HEAD -- examples/kotlin/ .a-check.yml
  spec/architecture.md docs/user/version.md sdks/csharp/ sdks/python/` —
  eigenständig, leer (keine unbeabsichtigte Berührung).
- Backtick-Zählung (`tr -cd '`' | wc -c`) über alle vier geänderten
  Markdown-Dateien — eigenständig, alle vier gerade.
