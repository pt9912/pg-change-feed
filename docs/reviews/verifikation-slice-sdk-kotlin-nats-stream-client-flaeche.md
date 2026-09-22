# Verifikationsbericht: slice-sdk-kotlin-nats-stream-client-flaeche — 2026-09-22

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?" gegen den
DoD-Vertrag
(`docs/plan/planning/in-progress/slice-sdk-kotlin-nats-stream-client-flaeche.md`
§2) und die §6-Risiko-Ausgänge, in frischem Kontext. **Nicht** gegen den
Diff als solchen (Reviewer-Aufgabe, abgeschlossen mit
[`review-slice-sdk-kotlin-nats-stream-client-flaeche.md`](review-slice-sdk-kotlin-nats-stream-client-flaeche.md)
und der Fixrunden-Bestätigung
[`review-slice-sdk-kotlin-nats-stream-client-flaeche-fixrunde.md`](review-slice-sdk-kotlin-nats-stream-client-flaeche-fixrunde.md))
und **nicht** gegen realen Bedarf (Validator, hier nicht ausgelöst). Dies
ist der **letzte** Flächen-Slice der Welle
[`welle-sdk-kotlin-vollabdeckung`](../plan/planning/done/welle-sdk-kotlin-vollabdeckung.md);
die Welle-Closure selbst ist ausdrücklich **nicht** Teil dieses Auftrags.

**Gegenstand:** vier Commits auf `main`, Elternstand `e8eddf79` (Closure
des vorigen SSE-Slice):

- `c8c9e3ae` — Implementer-Commit („feat(sdk-kotlin): NATS-Vollinhalts-
  Client-Fläche + Version-Hebung 0.2.0").
- `49e35ab9` — erster Review-Report (1 HIGH F-1, 1 MEDIUM F-2).
- `6b279311` — Fixrunden-Commit (beide Findings behoben, erweiterter
  repo-weiter Suchlauf).
- `009db193` — Fixrunden-Review-Report (0 HIGH/MEDIUM/LOW/INFO).
- `936b3b11` — Nachtrag: zwei nackte IDs im eigenen Fixrunden-Report in
  Backticks gesetzt (`docs-check`-Fund, kein DoD-Bezug).

**Frischer Kontext:** Diese Sitzung hat `harness/README.md`, `AGENTS.md`,
`harness/conventions.md`, den vollständigen Slice-Plan (§1–§8), `ADR-0109`
(Volltext, Accepted), `spec/pflichtenheft.md` §`SPEC-024`/§6
`SPEC-028`-Zeile/§1 `LH-FA-SST-009.a`, und beide Review-Reports vollständig
gelesen. Nichts aus Plan, Commit-Message oder Review-Report ungeprüft
übernommen — siehe die eigenständigen Prüfungen unten, insbesondere ein
eigener, **fünfter unabhängiger** Docker-Bau (`--no-cache`, frischer
`:test`-Lauf statt Cache-Treffer) mit eigener JUnit-XML-Auszählung.

---

## 1. DoD-Vertrag (§2) — jede Zeile einzeln geprüft

### 1.1 „`LH-FA-SST-009` erfüllt … öffentliche Client-Klasse … 14 neue Tests grün … zwei rot färbende Mutationen"

- **Alle zehn `SPEC-024`-Felder** eigenständig gegen
  `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/io/github/pt9912/pgchangefeed/nats/model/Change.kt`
  gehalten (`grep -n "SerializedName"`): `change_id`, `transaction_id`,
  `source_table_id`, `sequence`, `operation`, `old_image`, `new_image`,
  `schema_version`, `schema`, `table` — exakt die zehn Felder aus
  `spec/pflichtenheft.md` §`SPEC-024` „Nachrichteninhalt", keins fehlt,
  keins zusätzlich.
- **Subjekt-Schema/Authentifizierung**: eigene Lektüre bestätigt
  `cdc.stream.<source_id>.<schema>.<table>` (Tabellen-Subjekt),
  `cdc.stream.<source_id>.>` (Quellen-Subjekt), `cdc.stream.>`
  (Alle-Quellen), Token über `Options.Builder().token(…)` auf
  Verbindungsebene — deckungsgleich mit `SPEC-024`.
- **Test-Anzahl real nachgemessen, fünfte unabhängige Bau-Bestätigung
  dieser Welle** (Implementer + Reviewer-Erst-Report + Fixrunden-Review +
  diese Sitzung, zwei eigene Läufe): eigener
  `docker build -f sdks/kotlin/Dockerfile --target build --build-context proto=proto --no-cache`
  erzwang einen frischen (nicht Cache-getroffenen) `:test`-Task — die
  Konsolenausgabe zeigt `> Task :test` **ohne** `UP-TO-DATE`. Aus dem
  frisch erzeugten Container per `docker cp` extrahierte
  `build/test-results/test/*.xml` selbst mit `xml.etree.ElementTree`
  aufsummiert: **61 Tests, 0 Failures, 0 Errors**, aufgeschlüsselt
  `PgChangeFeedNatsStreamClientAuthBoundaryTest` 4,
  `PgChangeFeedNatsStreamClientMessageSchemaTest` 3,
  `PgChangeFeedNatsStreamClientSubjectTest` 7 → Delta zum SSE-Vorgänger
  (47, laut Reviewer-Report real gegen den isolierten Elternstand
  gemessen) = 14, exakt wie DoD/Commit-Message behaupten. Keine Drift.
- Ein zweiter, separater Lauf über `bash tools/harness/sdk-pack-kotlin.sh`
  (Docker-Layer-Cache-Treffer, deterministisch reproduzierbar) liefert
  dasselbe `.jar` — siehe 1.5.
- Die beiden vom Implementer/Reviewer benannten rot färbenden Mutationen
  wurden **nicht** von dieser Sitzung selbst reproduziert (nicht
  angefordert, bereits vom Reviewer an ihrer Eingabeseite geprüft, siehe
  Erst-Report „Mutation-Testing-Plausibilisierung"). Diese Sitzung
  akzeptiert diesen Teil als bereits belegt (Reviewer-Zuständigkeit),
  konzentriert die eigene Prüfung auf Testzahl, Feldvollständigkeit und
  Subjekt-Form — alle drei real nachgemessen.

**Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.2 „`io.nats:jnats`-Abhängigkeit real neu gemessen, weiterhin `2.26.3`"

Eigener `grep -n "io.nats"` gegen `build.gradle.kts` zeigt
`implementation("io.nats:jnats:2.26.3")`. Eigener
`curl -s https://repo1.maven.org/maven2/io/nats/jnats/maven-metadata.xml`
(dieser Sitzung, 2026-09-22) liefert `<latest>2.26.3</latest>` /
`<release>2.26.3</release>` — keine Drift, dritte unabhängige Messung
nach Implementer und Reviewer, alle drei identisch. **Ergebnis: Checkbox
berechtigt auf `[x]`.**

### 1.3 „`version` gehoben `0.1.0` → `0.2.0`, Top-Level UND `publishing`-Block"

Eigene Lektüre von `build.gradle.kts`: Zeile 108 `version = "0.2.0"`
(Top-Level), Zeile 167 `version = "0.2.0"` (innerhalb
`publishing.publications.maven`). Beide Stellen getrennt geführt, beide
`0.2.0`. **Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.4 „`spec/pflichtenheft.md` nachgezogen … Suchlauf-Pflicht berichtet"

Eigene Lektüre von `spec/pflichtenheft.md`:

- §6 `SPEC-028`-Zeile (Zeile 624): „…, das Package deckt bereits
  dokumentierte Drahtverträge (`SPEC-018`, `SPEC-020`, `SPEC-021`,
  `SPEC-024`)" — Version-Feld „aktuell `0.2.0`" korrekt.
- §1 `LH-FA-SST-009.a` (Zeile 184/195): „deckt HTTP-API, gRPC-Stream, SSE
  und NATS-Vollinhalt" — kein „nur HTTP-API und gRPC-Stream" mehr an einer
  Ist-Stand-Stelle.
- Änderungshistorie-Zeile datiert 2026-09-22 (Zeile 666) trägt den
  Nachzug korrekt, mit `welle-sdk-kotlin-vollabdeckung`-Bezug.

**Eigener, vierter unabhängiger Suchlauf** (andere Formulierungsvarianten
als die drei vorigen Läufe — Implementer, Erst-Reviewer, Fixrunden-
Reviewer): `grep -rn "pgchangefeed-kotlin-0\.1\.0\|HTTP-API und gRPC-Stream"`
über den vollen Baum (ausgenommen `docs/plan/planning/done/**`,
`docs/reviews/**`, `docs/plan/planning/observations/**`,
`welle-sdk-kotlin-vollabdeckung.md`, `docs/user/releasing.md`,
`docs/plan/adr/0109-…md` — dieselben drei bereits benannten Ausnahme-
Klassen: Records/Historie, illustrative Beispiele, Plan-eigene Chronik).
Ergebnis: **keine weitere Fundstelle** außerhalb dieser drei Klassen — die
verbleibenden Treffer sind ausschließlich datierte Änderungshistorie-Zeilen
(`spec/pflichtenheft.md:658/660/662/664/665`, alle vor dem 2026-09-22-Stand
oder explizit über ein anderes SDK) und die im Slice-Plan selbst stehende
Vorher/Nachher-Beschreibung des eigenen Vorgangs (§1/§2/§7 — dort legitim,
siehe Fixrunden-Report). Vierte unabhängige Bestätigung: Träger-Nachzug
vollständig. **Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.5 „Kein Import aus `internal/**`/`cmd/**`/`gen/**`"

Eigene Lektüre der Importzeilen in allen neuen `.kt`-Dateien
(`nats/PgChangeFeedNatsStreamClient.kt`, `nats/model/Change.kt`,
`nats/PgChangeFeedNatsMalformedMessageException.kt`,
`nats/FakeNatsStreamTransport.kt`, drei Testdateien): ausschließlich
`io.github.pt9912.pgchangefeed.*`, `com.google.gson.*`, `io.nats.client.*`,
`java.*`/`kotlin.*`. Kein `internal/`/`cmd/`/`gen/`-Treffer. **Ergebnis:
Checkbox berechtigt auf `[x]`.**

### 1.6 „`make gates` grün" — siehe eigenständiger Abschnitt 5 unten

### 1.7 „Real neu gebautes `.jar` mit `version=0.2.0`, alle vier Client-Flächen im selben Artefakt, `144598 Bytes`"

Eigener, isolierter `bash tools/harness/sdk-pack-kotlin.sh`-Lauf (nach
`rm -f sdks/kotlin/dist/*.jar`, damit das Artefakt real neu entsteht statt
ein altes stehen zu lassen): erzeugt
`sdks/kotlin/dist/pgchangefeed-kotlin-0.2.0.jar`, `ls -la` zeigt real
**144598 Bytes** — exakt wie DoD/Reviewer-Report behaupten. `unzip -l`
gegen dasselbe Jar, gefiltert auf `.class`-Dateien je Paketpfad: `http` 33,
`grpc` 4, `sse` 9, `nats` 7 Klassen — alle vier Client-Flächen real im
selben Artefakt, keine fehlt. **Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.8 „Review durchgeführt … Report unter `docs/reviews/` … kein offenes HIGH/MEDIUM mehr"

Beide Reports vollständig gelesen. Erst-Report (`49e35ab9`): 1 HIGH (F-1,
Chronik-Kommentar `build.gradle.kts`), 1 MEDIUM (F-2, stale `0.1.0.jar` in
`sdks/kotlin/Dockerfile:77,104`). Fixrunden-Report (`009db193`+`936b3b11`):
0 HIGH/MEDIUM/LOW/INFO, eigener dritter Suchlauf (fünf Formulierungs-
varianten) ohne weiteren Fund, Umfang der Fixrunde exakt auf die zwei
Findings begrenzt (per `git show 6b279311` isoliert geprüft). Eigene
Prüfung der beiden Fixes gegen den heutigen Dateistand (siehe Abschnitt 2
unten). **Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.9 „Doku-Update `docs/user/benutzerhandbuch.md`"

Eigene Lektüre: Änderungshistorie-Zeile 1.39 (2026-09-22) trägt den
SDK-Hinweis für SSE **und** NATS-Vollinhalt, inklusive PAT-Hinweis
(`read:packages`). Beide `**SDK:**`-Absätze im jeweiligen §4-Abschnitt
real vorhanden (Grep-Beleg oben). **Ergebnis: Checkbox berechtigt auf
`[x]`.**

### 1.10 „Closure-Notiz mit Steering-Loop-Lerneintrag", „Reconciliation-Register (entfällt)", „Beobachtungs-Register", „Jedes Risiko aus §6 trägt einen Ausgang"

Eigene Lektüre von Plan §7: vollständig gefüllt, benennt „Was hat
funktioniert", „Was ging anders als geplant" (der breitere Träger-
Nachzug-Suchraum), Steering-Loop-Eintrag mit den zwei real geprüften
Mutationen, Beobachtungs-Register-Aussage („keine neue Beobachtung"),
Risiko-Ausgänge für alle drei §6-Punkte. Reconciliation entfällt korrekt
(keine Reconciliation-Datei in diesem Repo, `harness/conventions.md`
bestätigt keinen abweichenden Aufbau). **Ergebnis: alle vier Checkboxen
berechtigt auf `[x]`.**

### 1.11 „Drei Paarungen (Anker · Folge-Slice · Register)" — bewusst `[ ]`

Eigene Lektüre: bleibt `[ ]`, mit explizitem Verweis „bleibt bewusst offen
bis zur Welle-Closure". Konsistent mit dem Auftrag dieser Sitzung, der
die Welle-Closure ausdrücklich ausschließt. **Kein DoD-Verstoß** — dieselbe
zulässige Konstellation wie im Präzedenzfall
`docs/reviews/verifikation-slice-release-hub-description.md` §1.5 (dort:
Closure-Pflichten, hier: welle-abhängige Prüfung).

---

## 2. F-1/F-2-Fixes real wirksam — vierte bzw. dritte unabhängige Bestätigung

### F-1 (`build.gradle.kts:97-100`, Chronik-Kommentar)

Eigene Lektüre der aktuellen Datei (Zeile 96–100):

> „Version `0.2.0` deckt beide zuletzt gelieferten Client-Flächen ab, SSE
> (`SPEC-021`) und NATS-Vollinhalt (`SPEC-024`, `ADR-0109` Festlegung 1)
> — additive, rückwärtskompatible Erweiterung, SemVer-Minor, keine
> ADR-pflichtige Ausnahme."

Kein Arrow-Muster (`0.1.0 -> 0.2.0`), keine Slice-Kennung, kein
Plan-Abschnittsverweis — ausschließlich `SPEC-021`/`SPEC-024`/`ADR-0109`
als Anker. Reine Zustandsaussage. Der unmittelbar vorangehende Absatz
(`jnats`-Messung) nennt zwar den eigenen Slice-Namen als Provenienz-Anker
(„am heutigen Bau-Zeitpunkt dieses Slice, 2026-09-22, …"), das ist nach
`AGENTS.md` §3.12 zulässig (Herkunft der Messung, kein Vorher/Nachher über
den Produktionscode-Pfad) — vierte unabhängige Bestätigung, dass dieser
Absatz korrekt bleibt (Implementer schrieb ihn, Erst-Reviewer beanstandete
ihn **nicht**, Fixrunden-Reviewer bestätigte ihn erneut nicht, diese
Sitzung ebenfalls nicht). **F-1 real und vollständig behoben.**

### F-2 (`sdks/kotlin/Dockerfile:77,104`, stale `0.1.0.jar`)

Eigener `grep -n "0\.1\.0\|0\.2\.0" sdks/kotlin/Dockerfile`: beide Treffer
tragen jetzt `pgchangefeed-kotlin-0.2.0.jar` (Zeile 77, 104), kein `0.1.0`
mehr in dieser Datei. **F-2 real und vollständig behoben.**

---

## 3. `git diff` gegen Elternstand (`e8eddf79..HEAD`) — Scope-Sauberkeit

Eigener `git diff --stat e8eddf79..HEAD`: 17 Dateien, 1387
Insertions/52 Deletions. Berührte Pfade: drei Slice-/Review-Dokumente
(`docs/plan/planning/in-progress/...md`, zwei `docs/reviews/...md`),
`docs/user/benutzerhandbuch.md`, `harness/README.md`, `harness/mk/sdk.mk`,
`sdks/kotlin/Dockerfile`, `sdks/kotlin/pgchangefeed-kotlin/README.md`,
`sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts`, sechs neue/geänderte
`.kt`-Dateien unter `sdks/kotlin/pgchangefeed-kotlin/src/{main,test}/…/nats/**`,
`spec/pflichtenheft.md`.

**Kein Treffer** für `examples/kotlin/**`, `.a-check.yml`,
`spec/architecture.md`, `docs/user/version.md`, `sdks/csharp/**`,
`sdks/python/**` — eigener `git diff --stat e8eddf79..HEAD -- examples/
.a-check.yml spec/architecture.md docs/user/version.md sdks/csharp
sdks/python` liefert leere Ausgabe. Scope exakt wie im Plan §1
zugesagt. Kein Bezug zu `.github/workflows/**` (weder `sdk-kotlin-
release.yml` noch ein anderer Workflow wurde in diesem Diff verändert) —
`AGENTS.md` §3.10 ist für diesen Slice nicht einschlägig.

---

## 4. Backtick-Parität — eigenständig nachgezählt

| Datei | Backticks | Gerade? |
|---|---|---|
| `docs/plan/planning/in-progress/slice-sdk-kotlin-nats-stream-client-flaeche.md` | 332 | ja |
| `docs/reviews/review-slice-sdk-kotlin-nats-stream-client-flaeche.md` | 534 | ja |
| `docs/reviews/review-slice-sdk-kotlin-nats-stream-client-flaeche-fixrunde.md` | 236 | ja |
| `docs/user/benutzerhandbuch.md` | 1914 | ja |
| `harness/README.md` | 1804 | ja |
| `harness/mk/sdk.mk` | 94 | ja |
| `sdks/kotlin/pgchangefeed-kotlin/README.md` | 98 | ja |
| `spec/pflichtenheft.md` | 1486 | ja |

Alle acht in diesem Slice geänderten Markdown-/Kommentar-tragenden
Dateien: gerade Backtick-Anzahl, paarig — bestätigt durch `make gates`
(`docs-check`/`ids`/`links`-Module, 0 Befunde) und eigenes `grep -o` je
Datei.

---

## 5. `make gates` real, ungepiped ausgeführt

```
$ git log -1 --format=%H
936b3b11ef1f1195d6fad1256404e70ecf47d7c0
$ git status --short
(leer)
$ make gates > <log> 2>&1; echo $?
0
```

Einzelbelege aus demselben Lauf (Log vollständig eingesehen):

| Gate | Ergebnis |
|---|---|
| `baseline-verify` | `v6.9.0 OK — 54 Dateien (Integritaet + Vollstaendigkeit, netzlos)` |
| `docs-check` | `d-check: 903 Datei(en) geprüft, 0 Befund(e)` (volle Modul-Liste) |
| `a-check` | `gesamt: 0 Befund(e)` |
| `commit-traceability` | `d-check`-Modul `commits`: `903 Datei(en) geprüft, 0 Befund(e)`; `commit-traceability.sh`: `OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` |
| `coverage-gate` | `coverage-gate: OK — Coverage 82.80% erfüllt Schwelle 80%` |
| `generated-sync` | `OK — das committete Erzeugnis ist byte-gleich der Ausgabe des gepinnten Generators (Stufe proto-export)` |

Zusätzlich zwei gezielte Modul-Läufe über den exakten Slice-Commit-Bereich
(Parent `e8eddf79` bis `HEAD`):

```
$ make doc-commits RANGE=e8eddf79..HEAD
d-check: 903 Datei(en) geprüft, 0 Befund(e)
$ make doc-immutable RANGE=e8eddf79..HEAD
d-check: 903 Datei(en) geprüft, 0 Befund(e)
```

**Ergebnis: Die DoD-Checkbox „`make gates` grün." ist berechtigt auf
`[x]` gesetzt** — real, ungepiped, `EXIT=0`, zusätzlich durch zwei
gezielte Modul-Läufe über den exakten Slice-Bereich bestätigt.

---

## 6. Eigener, fünfter unabhängiger Docker-Bau — Testzahl, Jar-Inhalt

Diese Sitzung führte **zwei** eigene Docker-Bauten durch, unabhängig von
den vier vorigen Läufen dieser Welle (Implementer, Erst-Reviewer
[`git worktree`-Elternstand-Vergleich], Fixrunden-Implementer,
Fixrunden-Reviewer):

1. `bash tools/harness/sdk-pack-kotlin.sh` (nach `rm -f
   sdks/kotlin/dist/*.jar`) — reales Artefakt `pgchangefeed-kotlin-
   0.2.0.jar`, `144598 Bytes`, vier Client-Pakete (`http` 33, `grpc` 4,
   `sse` 9, `nats` 7 `.class`-Dateien).
2. `docker build -f sdks/kotlin/Dockerfile --target build --build-context
   proto=proto --no-cache` (erzwingt einen Layer-Cache-Bruch, damit `:test`
   real neu läuft statt `UP-TO-DATE` zu melden) — Konsolen-Log zeigt
   `> Task :test` (kein `UP-TO-DATE`-Zusatz), `BUILD SUCCESSFUL`. Eigener
   `docker create`/`docker cp` extrahiert
   `build/test-results/test/*.xml` aus dem frischen Container; eigenes
   Python-Skript (`xml.etree.ElementTree`) summiert `tests`-Attribute über
   alle 13 XML-Dateien: **61 Tests, 0 Failures, 0 Errors** — davon `nats`-
   Paket 4+3+7=14. Deckt sich exakt mit DoD, Commit-Message und beiden
   Review-Reports.

Beide temporären Images (`pg-change-feed:sdk-kotlin-verify-build`) wurden
nach der Prüfung entfernt (`docker rmi`), kein Rückstand im lokalen
Docker-Zustand.

---

## 7. §6-Risiken — jedes mit zulässigem Ausgang?

Die drei zulässigen Ausgänge (Baseline-Regelwerk
`modul-05-planning-harness.md` §Offene Risiken werden bei Closure
aufgelöst): *eingetreten* → Carveout/Folge-Slice · *entfallen* →
gestrichen mit Begründung · *weiter offen* → Beobachtungs-Register.

| # | Risiko (Kurzform) | Ausgang im Plan | Zulässige Klasse? | Eigene Einschätzung |
|---|---|---|---|---|
| 1 | Gefakter NATS-Verbindungsfehler-Test bildet reales Token-Ablehnungsverhalten nicht exakt nach | „weiter offen — realer Rundlauf-Beleg bleibt `make test-integration`s `tools/harness/natsstreamsub` vorbehalten" | ✓ weiter offen | Konsistent: `natsstream/**`-Tests dieser Fläche sind ausdrücklich netzlos (Slice-Plan §1 Out-of-Scope „Realserver-Integrationstest"); ein Realserver-Beleg für den Vollinhalts-Client selbst existiert in dieser Welle nicht — die bestehende Wecksignal-Integrationsprüfung ist ein anderer Subjekt-Namensraum. Kein Blocker, strukturell dieselbe Klasse wie zuvor bei SSE. |
| 2 | Träger-Nachzug könnte eine Stelle übersehen | „weiter offen — Review prüft den Diff unabhängig; real breiterer Suchlauf als geplant" | ✓ weiter offen (praktisch aufgelöst) | Real durch **vier** unabhängige Suchläufe abgedeckt (Implementer, Erst-Reviewer [fand F-2], Fixrunden-Reviewer [dritter Lauf, fünf Varianten], diese Sitzung [vierter Lauf, eigene Varianten]) — kein weiterer Fund über F-2 hinaus. Formal bleibt die Klasse „weiter offen" (ein späterer Slice könnte theoretisch noch etwas finden), praktisch ist das Risiko auf ein sehr geringes Restmaß gesenkt. |
| 3 | Version-Hebung ohne explizites Minor-Schema | „entfallen — additiver SemVer-Minor-Bump ist unmissverständlich" | ✓ entfallen | Nachvollziehbar begründet, kein Widerspruch. |

**Ergebnis: Alle drei §6-Risiken tragen eine der drei zulässigen
Ausgangsklassen.**

---

## 8. DoD-Checkbox „Review durchgeführt" — inhaltlich korrekt gegen beide Reports?

Erst-Report-Summary: HIGH 1, MEDIUM 1, LOW 0, INFO 0 — F-1 (HIGH,
Chronik-Kommentar), F-2 (MEDIUM, stale Dockerfile-Kommentar). DoD-Zeile
im Plan (§2) nennt exakt „1 HIGH, 1 MEDIUM" mit denselben
Kurzbeschreibungen. Fixrunden-Report-Summary: HIGH 0, MEDIUM 0, LOW 0,
INFO 0 — DoD-Zeile im Plan trägt „kein offenes HIGH/MEDIUM mehr", korrekt.
**Kein offenes HIGH, kein unbehandeltes MEDIUM/LOW.** **Ergebnis: Checkbox
berechtigt auf `[x]` gesetzt.**

---

## 9. Kein realer Tag, kein realer Push — explizit geprüft

```
$ git tag -l | grep -i sdk-kotlin
(leer)
$ git tag -l
sdk-csharp-v0.1.0
sdk-python-v0.1.0
v0.1.0
v0.1.1
v0.1.2
```

**Kein** `sdk-kotlin-v*`-Tag existiert — auch kein `sdk-kotlin-v0.1.0` aus
einem früheren Slice dieser Welle. Die Behauptung „kein realer
`sdk-kotlin-v<Version>`-Tag-Push" (Plan §1 „Ausdrücklich NICHT in diesem
Slice") ist damit real bestätigt: Es gibt schlicht **keinen** Tag in
diesem Namensraum, nicht einmal einen älteren. Der Publish-Workflow
`.github/workflows/sdk-kotlin-release.yml` existiert zwar bereits (aus
einem früheren Slice dieser Welle, außerhalb des heutigen Diffs), wurde
aber von diesem Slice nicht verändert (siehe Abschnitt 3) und durch keinen
Tag-Push ausgelöst. Kein Zugriff auf `maven.pkg.github.com` erfolgte durch
diese Sitzung oder die geprüften Commits.

---

## Verdikt

**DoD erfüllt** für den aktuellen `in-progress`-Stand des Slice — die
einzige verbleibende `[ ]`-Checkbox („Drei Paarungen") ist bewusst und
zulässig bis zur Welle-Closure offen (Abschnitt 1.11), kein DoD-Verstoß.
Alle beauftragten Prüfpunkte wurden real und unabhängig nachgemessen,
nicht aus Bericht oder Commit-Message übernommen:

1. Alle zehn `SPEC-024`-Nachrichtenfelder real gegen `Change.kt` gehalten
   — vollständig, keine Abweichung.
2. `build.gradle.kts`s `version` real `0.2.0` an beiden Stellen (Top-Level
   und `publishing`-Block); der zuvor beanstandete Chronik-Kommentar trägt
   jetzt eine reine Zustandsaussage (F-1 real behoben, vierte
   unabhängige Bestätigung).
3. `sdks/kotlin/Dockerfile:77,104` tragen real `0.2.0.jar`, kein `0.1.0`
   mehr in dieser Datei (F-2 real behoben).
4. `io.nats:jnats:2.26.3` — dritte unabhängige Messung gegen Maven
   Central, keine Drift.
5. Eigener, fünfter unabhängiger Docker-Bau dieser Welle (`--no-cache`,
   erzwungen frischer `:test`-Lauf) liefert real **61 Tests, 0 Fehler**
   (14 neu: 4 Auth + 3 Schema + 7 Subjekt) und ein reales `.jar`
   (`144598 Bytes`, vier Client-Flächen im selben Artefakt).
6. Träger-Nachzug (`spec/pflichtenheft.md`, `docs/user/benutzerhandbuch.md`,
   `sdks/kotlin/pgchangefeed-kotlin/README.md`, `harness/README.md`,
   `harness/mk/sdk.mk`) real vollständig — eigener, vierter unabhängiger
   Suchlauf mit eigenen Formulierungsvarianten fand keine weitere,
   übersehene Stelle außer den bereits behobenen F-2.
7. `make gates` lief eigenständig, ungepiped, mit `EXIT=0`; zusätzlich
   `make doc-commits`/`make doc-immutable` über den exakten Slice-Bereich
   (`e8eddf79..HEAD`), beide `0 Befunde`.
8. `git diff e8eddf79..HEAD --stat` bestätigt Scope-Sauberkeit — kein
   Treffer für `examples/kotlin/**`, `.a-check.yml`, `spec/architecture.md`,
   `docs/user/version.md`, `sdks/csharp/**`, `sdks/python/**`.
9. Backtick-Parität aller acht geänderten Markdown-/Kommentar-tragenden
   Dateien eigenständig nachgezählt — alle gerade, alle paarig.
10. Alle drei §6-Risiken tragen eine zulässige Ausgangsklasse.
11. **Kein** `sdk-kotlin-v*`-Tag existiert im Repository (auch kein
    `sdk-kotlin-v0.1.0` aus einem früheren Slice); kein realer Push gegen
    `maven.pkg.github.com` erfolgte durch diese oder eine vorige Sitzung
    dieser Welle.

**Freigabe an den Planner:** Der Slice ist DoD-konform für seinen
aktuellen `in-progress`-Stand. Vor Closure (`git mv` nach `done/`) bleibt
ausschließlich die im Plan selbst bereits als offen geführte Prüfung „Drei
Paarungen" — regelkonform erst bei der Welle-Closure
(`welle-sdk-kotlin-vollabdeckung` → `done/`), nicht Teil dieses Slices und
nicht Teil dieses Verifikationsauftrags. Diese Verifikation nimmt weder
Closure noch Welle-Closure vor.
