# Verifikationsbericht: slice-sdk-kotlin-sse-client-flaeche — 2026-09-22

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?" gegen den
DoD-Vertrag
(`docs/plan/planning/in-progress/slice-sdk-kotlin-sse-client-flaeche.md`
§2) und die §6-Risiko-Ausgänge, in frischem Kontext. **Nicht** gegen den
Diff als solchen (Reviewer-Aufgabe, abgeschlossen mit
[`review-slice-sdk-kotlin-sse-client-flaeche.md`](review-slice-sdk-kotlin-sse-client-flaeche.md))
und **nicht** gegen realen Bedarf (Validator, hier nicht ausgelöst).

**Gegenstand:** zwei Commits auf `main`, Welle
[`welle-sdk-kotlin-vollabdeckung`](../plan/planning/done/welle-sdk-kotlin-vollabdeckung.md):

- `75fc1f3f` — Implementer-Commit („feat(sdks/kotlin): SSE-Stream-Client-Fläche"),
  neue Dateien unter `sdks/kotlin/pgchangefeed-kotlin/.../sse/**` (7
  Produktions-/Testdateien), Plan-Update, README-Korrektur.
- `88d4222d` (= `HEAD`) — Reviewer-Commit (0 HIGH, 0 MEDIUM, 2 LOW, 2
  INFO — keine Fixrunde nötig), inklusive Review-Report und
  DoD-Checkbox-Nachzug „Review durchgeführt" im selben Commit.

**Elternstand des vorigen Slice:** `baf6e8bd` (`next → in-progress`).

**Frischer Kontext:** Diese Sitzung hat `harness/README.md`,
`harness/conventions.md`, `AGENTS.md`, den vollständigen Slice-Plan
(§1–§8), `ADR-0109` (vollständig, Accepted), `spec/pflichtenheft.md`
§`SPEC-021` und den Review-Report (vollständig) selbst gelesen. Nichts aus
Slice-Plan, Commit-Message oder Review-Report ungeprüft übernommen: eigener
`grep`/`git diff`-Lauf gegen alle zehn `SPEC-021`-Felder in
`sse/model/Change.kt`; eigener `grep -rn "internal/\|cmd/\|gen/"`; eigener,
unabhängiger `docker build --no-cache` (kein Implementer-/Reviewer-Cache
übernommen) mit eigener JUnit-XML-Auszählung; eigener `javap -p`-Lauf gegen
ein selbst extrahiertes Jar; eigener `git diff --stat` über die volle
Slice-Range auf Fremdbereich-Berührung; eigene Backtick-Nachzählung;
eigener, ungepipter `make gates`-Lauf; eigene, isolierte
`make doc-commits`/`make doc-immutable`-Läufe über den exakten
Slice-Commit-Bereich; eigener Lese-Lauf des Beobachtungs-Registers
`BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut` samt
`git log`-Zeitstempel-Vergleich (siehe Abschnitt 6 — zentraler Befund
dieser Verifikation).

---

## 1. DoD-Vertrag (§2) — jede Checkbox einzeln geprüft

### 1.1 „`LH-FA-SST-009` erfüllt: … 16 neue Tests … Gesamtsuite 47 Tests, 0 Fehler"

Eigene Datei-Lektüre bestätigt die vier genannten Dateien real:
`sse/PgChangeFeedSseClient.kt`, `sse/SseFrameParser.kt`,
`sse/SseTransport.kt`, `sse/model/Change.kt` (plus vier Testdateien:
`FakeSseTransport.kt`, `PgChangeFeedSseClientAuthBoundaryTest.kt`,
`PgChangeFeedSseClientMessageSchemaTest.kt`, `SseFrameParserTest.kt`,
`TestClientFactory.kt`). Alle zehn `SPEC-021`-Felder aus
`spec/pflichtenheft.md` §`SPEC-021` einzeln gegen `sse/model/Change.kt`s
`@SerializedName`-Annotationen gehalten: `change_id`, `transaction_id`,
`source_table_id`, `sequence` (Long), `operation`, `old_image`/`new_image`
(`JsonElement?`), `schema_version`, `schema`, `table` — alle zehn
vorhanden, Typen passend zur Spec-Tabelle.

Eigener, **unabhängiger** Docker-Bau (kein Implementer-/Reviewer-Cache,
`--no-cache`, siehe Abschnitt 2) und eigene Auszählung der JUnit-XML-Reports
des selbst gebauten Images:

| Testklasse | `tests=` |
|---|---|
| `grpc.PgChangeFeedGrpcClientAuthBoundaryTest` | 2 |
| `grpc.PgChangeFeedGrpcClientMessageSchemaTest` | 1 |
| `http.PgChangeFeedHttpClientAuthBoundaryTest` | 9 |
| `http.PgChangeFeedHttpClientConsumerTest` | 5 |
| `http.PgChangeFeedHttpClientRetentionAndChangesTest` | 5 |
| `http.PgChangeFeedHttpClientTableTest` | 6 |
| `PgChangeFeedClientOptionsTest` | 3 |
| `sse.PgChangeFeedSseClientAuthBoundaryTest` | 7 |
| `sse.PgChangeFeedSseClientMessageSchemaTest` | 2 |
| `sse.SseFrameParserTest` | 7 |

Summe `2+1+9+5+5+6+3+7+2+7 = 47`, `failures="0" errors="0"` in **jeder**
Datei. Die drei neuen SSE-Klassen tragen exakt `7+7+2=16` — deckt sich
1:1 mit der DoD-Behauptung. **Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.2 „Kein Import aus `internal/**`/`cmd/**`/`gen/**`"

Eigener Lauf: `grep -rn "internal/\|cmd/\|gen/" sdks/kotlin/ --include="*.kt" --include="*.md" --include="*.kts"`
— zwei Treffer, beide Doku-Kommentar-Zitate des Go-Testvorbilds
(`PgChangeFeedSseClientMessageSchemaTest.kt:14` zitiert
`internal/adapters/driving/http/sse_test.go`; ein weiterer, außerhalb
dieses Diffs liegender Treffer im gRPC-Testpaket). Kein Treffer in einer
Produktionsdatei, kein Import. **Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.3 „`make gates` grün" — siehe eigenständiger Abschnitt 5 unten

### 1.4 „Review durchgeführt, Report unter `docs/reviews/` liegt vor …"

Eigene Lektüre des Reports: Summary-Tabelle 0 HIGH, 0 MEDIUM, 2 LOW, 2
INFO — DoD-Checkbox-Text (Plan §2) nennt exakt dieselbe Verteilung. Skill
`.harness/skills/reviewer.md`s Regel „DoD-Checkbox-Nachzug ohne Fixrunde"
korrekt angewendet: Da 0 HIGH/MEDIUM, keine Fixrunde nötig, Checkbox im
selben Commit (`88d4222d`) nachgezogen, der den Report anlegt — real per
`git diff 75fc1f3f..88d4222d -- docs/plan/planning/in-progress/slice-sdk-kotlin-sse-client-flaeche.md`
bestätigt (`[ ]` → `[x]`, Report-Pfad ergänzt). **Ergebnis: Checkbox
berechtigt auf `[x]`.**

### 1.5 „Doku-Update … bewusst nicht in diesem Slice … Träger-Nachzug-Suchlauf gegen README.md durchgeführt"

Eigener `git diff baf6e8bd..HEAD -- sdks/kotlin/pgchangefeed-kotlin/README.md`:
Die stehengebliebene Aussage „SSE and NATS-vollinhalt delivery remain out
of scope for this package's planned first full release" ist ersetzt durch
einen Satz, der SSE jetzt als gedeckt nennt und ausschließlich
NATS-Vollinhalt offen belässt. Eigener `git show baf6e8bd:…/README.md |
grep -in sse`: keine zweite, veraltete SSE-Aussage im Dokument. Eigener
`grep -n "^version" build.gradle.kts` → `version = "0.1.0"`, unverändert
(`git diff baf6e8bd..HEAD -- build.gradle.kts` liefert keinen Treffer für
diese Datei). Kein Bezug zu `spec/pflichtenheft.md`/
`docs/user/benutzerhandbuch.md` in diesem Diff (eigener
`git diff --stat`, siehe Abschnitt 4). **Ergebnis: Checkbox berechtigt auf
`[x]`.**

### 1.6 „Closure-Notiz mit Steering-Loop-Lerneintrag" / „Reconciliation-Register … entfällt" / „Beobachtungs-Register … keine Beobachtung angefallen … in §7 notiert" / „Jedes Risiko aus §6 trägt einen Ausgang"

Closure-Notiz (§7) ist real gefüllt (nicht `*(wird bei Bearbeitung
gefüllt.)*`), trägt einen Steering-Loop-Eintrag (Kotlin-`internal`-Lektion
von Anfang an korrekt angewendet). Eigener `find docs/plan/planning
-maxdepth 1 -iname "reconciliation*"` — leer: kein
`reconciliation.md` in diesem Repo, „entfällt" ist korrekt. Eigener
`git diff --stat baf6e8bd..HEAD -- 'docs/plan/planning/observations/**'`
— leer: keine Änderung an einem Beobachtungs-Eintrag, deckt sich mit der
Closure-Notiz-Aussage „keine neue Beobachtung angefallen". **Die
inhaltliche Aussage dieser Zeile ist jedoch teilweise faktisch veraltet —
siehe Abschnitt 6, zentraler Befund dieser Verifikation.** Beide §6-Risiken
tragen einen Ausgang (siehe Abschnitt 6). **Ergebnis der reinen
Checkbox-Prüfung: alle vier Zeilen berechtigt auf `[x]`; ein inhaltlicher
Befund unterhalb der Checkbox-Ebene wird in Abschnitt 6 gesondert
berichtet.**

### 1.7 „Die drei Paarungen … noch offen — Prüfung läuft regelkonform bei Wellen-Closure"

Unverändert `[ ]`. Die Welle `welle-sdk-kotlin-vollabdeckung` ist real noch
nicht geschlossen (eigener `find` — Datei liegt weiterhin unter
`docs/plan/planning/welle-sdk-kotlin-vollabdeckung.md`, kein `done/`-Pfad).
Korrekt offen für einen `in-progress`-Slice, dessen Welle noch läuft.

## 2. Eigener, unabhängiger Docker-Bau — dritte Bau-Bestätigung nach Implementer und Reviewer

```
$ rm -rf sdks/kotlin/dist && bash tools/harness/sdk-pack-kotlin.sh
… (Cache-Hit auf bereits vorhandene Layer) … EXIT=0
$ ls sdks/kotlin/dist/
pgchangefeed-kotlin-0.1.0.jar
```

Da dieser erste Lauf teils gecachte Layer traf, folgte ein zweiter,
**komplett cache-loser** Bau zur echten Unabhängigkeitsprüfung:

```
$ docker build -f sdks/kotlin/Dockerfile sdks/kotlin \
    --build-context proto=proto --target build --no-cache \
    -t verify-kotlin-sse-nocache
…
[build 10/11] RUN ./gradlew --no-daemon test   → 13 actionable tasks: 13 executed, BUILD SUCCESSFUL in 50s
[build 11/11] RUN ./gradlew --no-daemon build  → BUILD SUCCESSFUL
```

Alle 13 Gradle-Tasks real **executed** (nicht `UP-TO-DATE`/`CACHED`) in der
`test`-Stufe — ein vollständig frischer Testlauf, unabhängig von jedem
zuvor gebauten Layer. Die anschließende `build`-Stufe zeigt `UP-TO-DATE`,
weil sie auf der Dateisystem-Ausgabe der unmittelbar vorangegangenen
`RUN`-Anweisung **desselben** Bau-Vorgangs aufbaut (Docker-Layer-Verkettung
innerhalb eines `--no-cache`-Baus, kein Cache-Rückgriff auf einen
früheren Bauversuch — real durch die Executed-Zählung in Schritt `test`
belegt).

JUnit-XML-Reports aus dem selbst gebauten Image extrahiert (`docker
create`/`docker cp`, kein Bind-Mount) und ausgezählt — siehe Tabelle in
Abschnitt 1.1: **47 Tests gesamt, 16 neu (7+7+2), 0 Fehler/Errors** — real
bestätigt, dritte unabhängige Messung nach Implementer und Reviewer, alle
drei decken sich exakt.

## 3. Kotlin-`internal`-Semantik — dritte unabhängige Bestätigung per `javap -p`

Jar aus demselben, cache-losen Image extrahiert, mit `javap -p` im
gepinnten `eclipse-temurin:21-jdk`-Image (identischer Digest zur
Dockerfile-`FROM`-Zeile) geprüft:

```
public interface io.github.pt9912.pgchangefeed.sse.SseTransport { … }
public final class io.github.pt9912.pgchangefeed.sse.JdkSseTransport implements SseTransport { … }
public final class io.github.pt9912.pgchangefeed.sse.PgChangeFeedSseClient {
  public io.github.pt9912.pgchangefeed.sse.PgChangeFeedSseClient(SseTransport, PgChangeFeedClientOptions);  // war "internal" im Kotlin-Quelltext
  …
}
```

`SseTransport` kompiliert zu einem gewöhnlichen `public interface`,
`JdkSseTransport` zu einer `public final class`, und der im Kotlin-Quelltext
als `internal` deklarierte Konstruktor
`PgChangeFeedSseClient(SseTransport, PgChangeFeedClientOptions)` ist im
Bytecode ein gewöhnlicher `public`-Konstruktor — genau das, was
`SseTransport.kt`s KDoc behauptet (compile-time Kotlin-Grenze, keine
JVM-Bytecode-Schranke). Dritte unabhängige Bestätigung nach Implementer
(Closure-Notiz) und Reviewer (Review-Report) — die aus dem HIGH-Finding
des vorigen Kotlin-Slice (`review-slice-sdk-kotlin-http-client-flaeche.md`
F-1) gelernte Lektion ist real, dreifach unabhängig verifiziert, korrekt
umgesetzt. **Kein Rückfall.**

## 4. `git diff` gegen den Elternstand — unbeabsichtigte Fremdbereich-Berührung?

```
$ git diff --stat baf6e8bd..HEAD
 …/slice-sdk-kotlin-sse-client-flaeche.md         |  87 ++++-
 …/review-slice-sdk-kotlin-sse-client-flaeche.md  | 400 +++++++++++++++++++++
 sdks/kotlin/pgchangefeed-kotlin/README.md          |   4 +-
 …/sse/PgChangeFeedSseClient.kt                     | 173 +++++++++
 …/sse/SseFrameParser.kt                            |  45 +++
 …/sse/SseTransport.kt                              |  76 ++++
 …/sse/model/Change.kt                              |  39 ++
 …/sse/FakeSseTransport.kt                          |  33 ++
 …/sse/PgChangeFeedSseClientAuthBoundaryTest.kt      | 139 +++++++
 …/sse/PgChangeFeedSseClientMessageSchemaTest.kt     |  62 ++++
 …/sse/SseFrameParserTest.kt                         |  98 +++++
 …/sse/TestClientFactory.kt                          |  24 ++
 12 files changed, 1161 insertions(+), 19 deletions(-)
$ git diff --stat baf6e8bd..HEAD -- 'examples/kotlin/**' '.a-check.yml' \
    'spec/architecture.md' 'docs/user/version.md' 'sdks/csharp/**' \
    'sdks/python/**' '…/http/**' '…/grpc/**'
(leer)
```

Zwölf Dateien insgesamt, keine davon außerhalb von Plan/README/`sse/**`/
Review-Report. Alle neun ausdrücklich vom Auftrag genannten
Fremdbereiche (`examples/kotlin/**`, `.a-check.yml`,
`spec/architecture.md`, `docs/user/version.md`, `sdks/csharp/**`,
`sdks/python/**`, bestehende `http/`/`grpc/`-Flächen desselben Packages)
sind unberührt. **Kein unbeabsichtigter Fremdbereich-Zug.**

## 5. Backtick-Parität aller geänderten Markdown-Dateien — eigenständig nachgezählt

| Datei | Backtick-Anzahl (eigene Zählung) | Gerade? |
|---|---|---|
| `docs/plan/planning/in-progress/slice-sdk-kotlin-sse-client-flaeche.md` | 254 | ja |
| `sdks/kotlin/pgchangefeed-kotlin/README.md` | 84 | ja |
| `docs/reviews/review-slice-sdk-kotlin-sse-client-flaeche.md` | 522 | ja |

Der Slice-Plan zeigt 254 statt der vom Reviewer auf `75fc1f3f` gemessenen
252 — eigener `git diff 75fc1f3f..88d4222d`-Beleg klärt die Differenz:
Der Reviewer-Commit selbst fügt beim DoD-Checkbox-Nachzug ein zusätzliches
Backtick-Paar ein (der neu ergänzte Report-Pfad-Verweis) — beide
Messungen sind an ihrem jeweiligen Stand korrekt und paarig, keine Drift.
Alle drei Dateien paarig. **Kein Befund.**

## 6. §6-Risiken und das Beobachtungs-Register — zentraler Befund dieser Verifikation

Beide §6-Risiken tragen einen zulässigen Ausgang („weiter offen",
Register-Referenz bzw. transparente Benennung im Plan). Das erste Risiko
(gestubbter HTTP-Response-Stream vs. reales Chunked-Transfer-/Flush-
Verhalten) ist unverändert korrekt als „weiter offen" geführt — kein
Befund.

**Das zweite Risiko und die begleitende Beobachtungs-Referenz sind
faktisch veraltet.** Plan §3, §6, §7 und §8 behaupten durchgängig, der
Zähler von `BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut`
stehe bei **2×** und ein „dritter Kotlin-Treffer" würde ihn auf die
3×-Schwelle heben; die Closure-Notiz (§7) schließt daraus, der Zähler
„bleibe bei 2×", weil dieser Slice keinen dritten (isolierten
Bau-Versuch-)Treffer erzeugt habe.

Eigener Lese-Lauf des Registers zum Zeitpunkt dieser Verifikation:

```
$ cat docs/plan/planning/observations/BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut/state.md
Stand: offen (3×, Schwelle erreicht — Ausgang noch nicht zugewiesen).
Zähler (abgeleitet): 3× (evidence/slice-102.md, evidence/slice-103.md,
evidence/slice-sdk-csharp-sse-client-flaeche.md). …
```

Der Zähler steht real bei **3×**, nicht 2× — die dritte Instanz stammt aus
`slice-sdk-csharp-sse-client-flaeche` (Baum `sdks/csharp/**`, eine
strukturell andere, aber symptomgleiche Ursache, siehe
`observation.md`), **nicht** aus einem Kotlin-Treffer. Eigener
Zeitstempel-Abgleich klärt die Reihenfolge unzweideutig:

```
$ git show -s --format="%H %ci %s" d8cc10a1 a8923583 baf6e8bd
d8cc10a1 2026-09-22 00:10:21 +0200  docs(planning): Closure-Inhalt slice-sdk-csharp-sse-client-flaeche (LH-FA-SST-009, ADR-0106)
a8923583 2026-09-22 05:32:08 +0200  docs(plan): welle-sdk-csharp-vollabdeckung Closure-Notiz + Roadmap (LH-FA-SST-009, ADR-0106)
baf6e8bd 2026-09-22 06:06:49 +0200  docs(plan): slice-sdk-kotlin-sse-client-flaeche next -> in-progress (LH-FA-SST-009, ADR-0109)
```

Der Commit, der den Zähler auf 3× hebt (`d8cc10a1`), landete **vor**
`baf6e8bd` — dem Commit, mit dem dieser Kotlin-SSE-Slice überhaupt erst
`in-progress` ging. Der Zähler stand also bereits bei 3×, **bevor** dieser
Slice begann, nicht erst danach. Die im Plan wiederholt behauptete
Ausgangslage „2×, unter der Schwelle" war zum Zeitpunkt ihrer
Niederschrift bereits falsch — eine stale, nicht am aktuellen
Register-Stand re-verifizierte Tatsachenbehauptung (`AGENTS.md` §3.12
Instanz B: Aussagen über den Zustand eines Trägers, die als geprüft
formuliert sind, brauchen einen Beleg-Anker am tatsächlichen Stand zum
Lesezeitpunkt). §8s „Vorgelagert — offene Beobachtungen sichten"-Schritt
hat den Eintrag offenbar aus einer älteren Erinnerung oder einem
Vorgänger-Slice übernommen, statt das Register zum eigenen Lesezeitpunkt
neu zu öffnen.

**Einordnung — kein DoD-Blocker, aber ein realer Befund für den
Planner:**

- Die im Diff dieses Slice geprüfbare Tatsache bleibt korrekt: `git diff
  --stat baf6e8bd..HEAD -- 'docs/plan/planning/observations/**'` ist real
  leer — dieser Slice selbst hat **keinen** vierten Beleg erzeugt (jeder
  Docker-Bau trug `--build-context proto=proto` von Anfang an, wie
  geplant). Die Handlung war richtig.
- Falsch ist die **Begründung** dafür in Prosa: „der Zähler bleibt bei
  2×" beschreibt einen Zustand, der zum Zeitpunkt dieses Slice-Starts
  bereits nicht mehr zutraf. Korrekt wäre gewesen: „der Zähler steht
  bereits bei 3× (Schwelle erreicht, ausgelöst durch
  `slice-sdk-csharp-sse-client-flaeche` vor Beginn dieses Slice) und
  bleibt durch diesen Slice unverändert bei 3×."
- Diese Klasse von Befund ist genau die, für die die Verifier-Rolle
  existiert: unsichtbar für den Reviewer, der die geprüfte Aussage „keine
  Änderung an `observations/`" korrekt, aber ohne den absoluten
  Register-Stand zum eigenen Lesezeitpunkt gegenzuprüfen, bestätigt hat
  (`review-slice-sdk-kotlin-sse-client-flaeche.md`, Negativbefund
  „docs/plan/planning/observations/ … deckt sich mit der
  Closure-Notiz-Aussage … der 2×-Zähler … bleibt unverändert bei 2×") —
  eine Behauptung ohne erneuerte Prüfung des tatsächlichen
  Registerstands, an genau der Stelle, an der dieser Bericht sie fängt.
- Kein DoD-Verstoß im engeren Sinn (keine Checkbox wird dadurch
  unberechtigt), aber eine Textkorrektur-Pflicht in §3/§6/§7/§8 dieses
  Slice-Plans, bevor er als kanonische Closure-Aussage in `done/`
  wandert, und ein Hinweis für die anstehende Wellen-Closure von
  `welle-sdk-kotlin-vollabdeckung` (die Wellen-Closure von
  `welle-sdk-csharp-vollabdeckung` hat den dritten Treffer bereits selbst
  vorausgesehen und im Register verankert — das ist die Quelle, die der
  Kotlin-Slice hätte lesen müssen).

## 7. `make gates` real, ungepiped ausgeführt

```
$ git status --short
(leer)
$ make gates > gates1.log 2>&1; echo $?
0
```

Einzelbelege aus demselben Lauf (Log vollständig eingesehen):

| Gate | Ergebnis |
|---|---|
| `baseline-verify` | `v6.9.0 OK — 54 Dateien (Integritaet + Vollstaendigkeit, netzlos)` |
| `docs-check` | `d-check: 899 Datei(en) geprüft, 0 Befund(e)` (volle Modul-Liste: links, anchors, ids, matrix, versions, structure, hostpaths, tracked) |
| `commit-traceability` | `d-check`-Modul `commits`: `899 Datei(en) geprüft, 0 Befund(e)`; `commit-traceability.sh`: `OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` |
| `coverage-gate` | `coverage-gate: OK — Coverage 82.80% erfüllt Schwelle 80%` |
| `generated-sync` | `OK — das committete Erzeugnis ist byte-gleich der Ausgabe des gepinnten Generators (Stufe proto-export)` |
| `a-check` | `gesamt: 0 Befund(e)` |

Zusätzlich isoliert über den exakten Slice-Commit-Bereich:

```
$ make doc-commits RANGE=baf6e8bd..HEAD
d-check: 899 Datei(en) geprüft, 0 Befund(e)
$ make doc-immutable RANGE=baf6e8bd..HEAD
d-check: 899 Datei(en) geprüft, 0 Befund(e)
```

Kein ADR-Datei wurde in diesem Slice berührt (`git diff --stat
baf6e8bd..HEAD -- 'docs/plan/adr/**'` ist leer) — die
`doc-immutable`-Prüfung ist damit für diesen Slice strukturell moot,
läuft aber gleichwohl grün. **Ergebnis: Die DoD-Checkbox „`make gates`
grün." ist berechtigt auf `[x]` gesetzt** — real, ungepiped, `EXIT=0`,
zusätzlich durch zwei gezielte Modul-Läufe über den exakten Slice-Bereich
bestätigt.

## 8. Traceability der beiden Slice-Commits

```
$ git log -1 --format=%s 75fc1f3f
feat(sdks/kotlin): SSE-Stream-Client-Fläche (LH-FA-SST-009, ADR-0109)
$ git log -1 --format=%s 88d4222d
docs(reviews): Review slice-sdk-kotlin-sse-client-flaeche (LH-FA-SST-009, ADR-0109)
```

Beide Commits tragen `LH-FA-SST-009`/`ADR-0109` im Betreff, kein
`SPEC-*`/`ARC-*`-Betreff. **Kein Befund.**

---

## Verdikt

**DoD erfüllt — mit einem nicht blockierenden, aber zu korrigierenden
Befund.** Alle sechs beauftragten Prüfblöcke wurden real und unabhängig
nachgemessen, nicht aus Bericht oder Commit-Message übernommen:

1. Alle DoD-Checkboxen sind gegen den realen Code-/Doku-Stand berechtigt
   gesetzt — SSE-Fläche real vorhanden, alle zehn `SPEC-021`-Felder
   gedeckt, Import-Grenze eingehalten, README-Korrektur vollständig ohne
   Rest, Version unverändert `0.1.0`, kein Bezug zu
   `spec/pflichtenheft.md`/`docs/user/benutzerhandbuch.md`.
2. Ein eigener, **cache-loser** Docker-Bau (dritte unabhängige
   Bau-Bestätigung nach Implementer und Reviewer) bestätigt: `BUILD
   SUCCESSFUL`, 13 Tasks `executed`, 47 Tests gesamt/16 neu/0 Fehler —
   exakt deckungsgleich mit beiden Vorbehauptungen.
3. `javap -p` gegen ein selbst extrahiertes Jar bestätigt dritte
   unabhängige Mal: die Kotlin-`internal`-KDoc-Aussage ist korrekt
   (compile-time Kotlin-Grenze, keine JVM-Bytecode-Schranke) — kein
   Rückfall auf das HIGH-Finding des vorigen Kotlin-Slice.
4. `git diff --stat` über die volle Range zeigt ausschließlich die zwölf
   erwarteten Dateien; keiner der neun genannten Fremdbereiche berührt.
5. Backtick-Parität aller drei geänderten/neuen Markdown-Dateien
   eigenständig nachgezählt, alle paarig; die kleine Differenz zum
   Reviewer-Wert (252 vs. 254) ist durch den Reviewer-eigenen
   Checkbox-Nachzug-Commit erklärt, keine Drift.
6. **Zentraler Befund:** Der im Plan (§3/§6/§7/§8) durchgängig behauptete
   Zählerstand „2×" für
   `BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut` ist zum
   Zeitpunkt des Slice-Starts bereits veraltet — das Register stand real
   schon bei **3× (Schwelle erreicht)**, ausgelöst durch
   `slice-sdk-csharp-sse-client-flaeche` (Commit `d8cc10a1`,
   2026-09-22 00:10) — **vor** dem `next → in-progress`-Übergang dieses
   Kotlin-Slice (`baf6e8bd`, 2026-09-22 06:06). Die **Handlung** dieses
   Slice war richtig (kein vierter Beleg erzeugt, `--build-context
   proto=proto` von Anfang an gesetzt); die **Begründung** in Prosa
   („Zähler bleibt bei 2×") ist eine stale Tatsachenbehauptung
   (`AGENTS.md` §3.12 Instanz B) und braucht eine Textkorrektur, bevor
   dieser Plan als kanonische Closure-Aussage nach `done/` wandert. Der
   Reviewer hat dieselbe Stelle geprüft, aber nur auf „keine Diff-Änderung
   an `observations/`" — nicht auf den absoluten Register-Stand zum
   eigenen Lesezeitpunkt; genau diese Lücke ist die Verifier-only-Klasse,
   für die diese Rolle existiert.
7. `make gates` lief eigenständig, ungepiped, mit `EXIT=0`; zusätzlich
   `make doc-commits`/`make doc-immutable` über den exakten Slice-Bereich,
   beide `0 Befund(e)`.

**Freigabe an den Planner:** Der Slice ist DoD-konform für seinen
aktuellen `in-progress`-Stand. Vor Closure (`git mv` nach `done/`) und vor
der anstehenden `welle-sdk-kotlin-vollabdeckung`-Closure ist die
Zählerstand-Aussage in §3/§6/§7/§8 dieses Plans auf den tatsächlichen
Registerstand (3×, unverändert durch diesen Slice) zu korrigieren — reine
Textkorrektur, kein Code-/Test-Bezug, kein Gate betroffen. Die
Registerstand-Frage selbst (Ausgang des 3×-Fundes,
Regelschärfungs-Entscheidung) bleibt — wie im Register selbst bereits
vermerkt — eine künftige Architect-Entscheidung, keine Pflicht dieses
Slice.
