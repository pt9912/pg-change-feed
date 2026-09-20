# Review-Report: slice-sdk-kotlin-http-client-flaeche — 2026-09-20

**Review-Art:** Code — geprüft gegen Plan
(`docs/plan/planning/in-progress/slice-sdk-kotlin-http-client-flaeche.md`,
inkl. Plan-Nachzug in §3), `ADR-0109` (Accepted) und `AGENTS.md` Hard Rules
(Modul 10 §Drei Review-Arten). Kein DoD-Abgleich — das ist Verifier-Aufgabe
(Modul 11).

**Gegenstand:** Diff `df6acd80..7d01359b` (Abschluss von
`slice-sdk-kotlin-projektgeruest` bis Implementer-Commit dieses Slice),
ein einziger Commit `7d01359b` (feat(sdk): Kotlin-HTTP-Client-Fläche für
pgchangefeed-kotlin), Slice `slice-sdk-kotlin-http-client-flaeche`, Welle
`welle-sdk-kotlin-lh-fa-sst-009`.

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09"
(vier repo-spezifische HIGH-Regeln, drei MEDIUM-Klassen, seither um mehrere
weitere HIGH-Klassen ergänzt — u. a. `AGENTS.md` §3.12/§3.13, „Beleg trägt
seinen Satz nicht").
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-20.

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-sdk-kotlin-http-client-flaeche.md`
  (§1 Ziel/Abgrenzung, §2 DoD, §3 Plan inkl. Plan-Nachzug, §6 Risiken, §8
  Sub-Area/Modus)
- `ADR-0109` (Accepted) — Festlegung 1/2/3/5, §Verglichene Alternativen A
  (Java-Binärkompatibilität als genannter Vorteil), §Konsequenzen
  Folgepflicht 2/4
- `spec/pflichtenheft.md` §2 `SPEC-018` (Endpunkte, Token-Header-Form,
  Fehler-Antwortform, alle neun Fähigkeiten), `SPEC-022` (`GET /changes`)
- `AGENTS.md` §3.1 (Docker-only), §3.7 (Kommentar-Disziplin), §3.12
  (Herkunft von Aussagen), §3.13 (Träger-Nachzug)
- `harness/conventions.md` (MR-000 ID-Schema, MR-002 Slice-Kennungen)
- `sdks/csharp/PgChangeFeed.Client/Http/PgChangeFeedException.cs` als
  Vergleichsfassung für die versiegelte Fehlerhierarchie
- `docs/reviews/review-slice-sdk-csharp-http-client-flaeche.md` — Formvorbild
  für diesen Report, sowie Quelle des dort bereits behobenen F-2-Musters
  (malformter `2xx`-Erfolgs-Body)

**Eigenständig durchgeführte Prüfungen (nicht nur Commit-Message/Plan-Text
übernommen):**

- Alle neun `SPEC-018`-Fähigkeiten (`RegisterConsumer`,
  `AcknowledgeConsumer`, `GetConsumerPosition`, `RemoveConsumer`,
  `EnableTable`, `DisableTable`, `GetStatus`, `ListTables`, `RunRetention`)
  plus `readChanges` (`SPEC-022`) Feld für Feld gegen die
  `spec/pflichtenheft.md`-Tabellen gehalten — Endpunkt, Methode, Request-/
  Response-Feldnamen, Pflicht/optional, Statuscodes: keine Abweichung.
- `PgChangeFeedException.kt` (sieben Unterklassen) Zeile für Zeile gegen
  `sdks/csharp/PgChangeFeed.Client/Http/PgChangeFeedException.cs` gehalten:
  identische Klassennamen, identische Statuscode-Zuordnung
  (`400`/`401`/`403`/`404`/`500` plus `UnexpectedStatus`/
  `MalformedResponse`), identisches Verhalten.
- `grep -rn "internal/\|cmd/" sdks/kotlin/` real ausgeführt: einziger
  Treffer ist `gradlew:60`, ein Kommentar im Gradle-Wrapper-Skript, der auf
  eine GitHub-URL innerhalb der Gradle-Quelle verweist
  (`.../api/internal/plugins/unixStartScript.txt`) — kein Import dieses
  Repos, kein Treffer im SDK-Code selbst.
- Kotlins `internal`-Sichtbarkeit real gegen die JVM-Bytecode-Realität
  geprüft (nicht nur die Plan-/KDoc-Aussage übernommen) — siehe F-1.
- `docker build -f sdks/kotlin/Dockerfile sdks/kotlin` real ausgeführt
  (Cache-Hit auf dem bereits vom Implementer gebauten Layer): Exit `0`,
  `./gradlew test`/`./gradlew build` beide grün.
- Aus demselben Bau-Image den Jar extrahiert und mit `javap -p` gegen
  `PgChangeFeedHttpClient.class`/`HttpTransport.class` geprüft (Beleg für
  F-1, Kommandos siehe dort).
- `com.google.code.gson:gson`-Version real gegen
  `https://repo1.maven.org/maven2/com/google/code/gson/gson/maven-metadata.xml`
  abgerufen: `<latest>`/`<release>` beide `2.14.0` — deckt sich mit der im
  Plan-Nachzug genannten Messung.
- Backtick-Parität real nachgezählt (Backtick-Zeichen je Datei gezählt,
  gerade Anzahl erwartet) für alle drei im Diff geänderten Markdown-Dateien
  (`docs/plan/planning/in-progress/slice-sdk-kotlin-http-client-flaeche.md`:
  370, `docs/user/benutzerhandbuch.md`: 1742,
  `sdks/kotlin/pgchangefeed-kotlin/README.md`: 58) — alle drei gerade.
- `grep -rniE "folgt erst|follow in subsequent|coming soon|not yet implemented|TODO|FIXME"
  sdks/kotlin/` real ausgeführt: kein Treffer — bestätigt den im Plan-Nachzug
  behaupteten Träger-Nachzug-Suchlauf des Implementers (README-Korrektur
  bereits im Diff).
- `git show df6acd80..7d01359b -- .../README.md docs/user/benutzerhandbuch.md
  .../build.gradle.kts` gelesen: README-Korrektur, Handbuch-Version
  1.35→1.36 samt neuer Änderungshistorie-Zeile, `gson`-Abhängigkeit mit
  begründendem Kommentar — alle drei wie im Plan-Nachzug beschrieben.
- Out-of-Scope-Grep: kein `fun main(` (kein CLI-Parsing) und keine
  gRPC-/SSE-/NATS-Codepfade unter `sdks/kotlin/pgchangefeed-kotlin/src/`
  (die wenigen Treffer für „sse"/„nats" sind Fehlalarme —
  `assertEquals`/`assertFailsWith`/„internal wrapped" — einzeln
  gegengelesen).
- Mutation-Behauptungen (Schritt 19 des Plan-Nachzugs) an der Testlogik
  nachvollzogen, nicht nur übernommen: `listTables percent-encodes reserved
  characters in query parameters`
  (`PgChangeFeedHttpClientTableTest.kt:105`) prüft die exakt gebaute URL
  inkl. `%20`/`%26`/`%2F` — eine Identitäts-`encode()` würde die Assertion
  sichtbar brechen; `reader token against an admin endpoint throws
  Forbidden` (`PgChangeFeedHttpClientAuthBoundaryTest.kt:34`) erwartet
  `PgChangeFeedForbiddenException` via `assertFailsWith` — eine vertauschte
  `403`→`BadRequest`-Zuordnung in `buildException` würfe den falschen Typ
  und ließe `assertFailsWith` fehlschlagen. Beide Mutationen wurden anhand
  der Logik nachvollzogen, kein eigener Mutations-Lauf zusätzlich zum
  bereits vom Implementer dokumentierten gefahren (Docker-Cache-Hit oben
  bestätigt den aktuellen, unmutierten Grünzustand unabhängig).
- `make gates` real ausgeführt, Exit-Code direkt (ungepiped) geprüft: `0`
  (`d-check: 859 Datei(en) geprüft, 0 Befund(e)`, `generated-sync: OK`,
  `a-check: gesamt: 0 Befund(e)`).

---

## Findings

### F-1 — Der `internal`-Transport-Seam ist entgegen der KDoc-/Plan-Aussage kein JVM-Sichtbarkeitsschutz — für den von `ADR-0109` selbst benannten Java-Konsumenten technisch aufrufbar

- `kategorie`: **HIGH**
- `quelle`: `AGENTS.md` §3.7 (Kommentar beschreibt, was da ist) i. V. m.
  der Reviewer-Skill-Klasse „Beleg trägt seinen Satz nicht" — die Aussage
  trägt einen Beleg (Kotlin-Modul-Sichtbarkeit), der geprüfte Beleg
  (JVM-Bytecode) stützt die volle Aussage nicht.
- `pfad`: `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/io/github/pt9912/pgchangefeed/http/HttpTransport.kt:47-53`
  (KDoc: „is invisible outside this module, so it adds no public API
  surface"); dieselbe Aussage im Plan-Nachzug,
  `docs/plan/planning/in-progress/slice-sdk-kotlin-http-client-flaeche.md:169-172`
  („ist nur innerhalb des Gradle-Moduls sichtbar")
- `befund`: Kotlins `internal`-Modifikator ist eine **compile-time**-Prüfung
  des Kotlin-Compiler-Frontends gegen die Modul-Metadaten eines anderen
  **Kotlin**-Moduls — er ist **keine** JVM-Bytecode-Zugriffsbeschränkung.
  Real mit `javap -p` gegen das im Diff selbst gebaute Jar geprüft:
  ```
  public io.github.pt9912.pgchangefeed.http.PgChangeFeedHttpClient(io.github.pt9912.pgchangefeed.http.HttpTransport, io.github.pt9912.pgchangefeed.PgChangeFeedClientOptions);
  ```
  (unmangled `public` Konstruktor — Konstruktoren werden von Kotlins
  Namens-Mangling für `internal`-Member nicht erfasst, weil ein
  JVM-Konstruktor immer `<init>` heißt) und
  ```
  public interface io.github.pt9912.pgchangefeed.http.HttpTransport {
    public abstract io.github.pt9912.pgchangefeed.http.TransportResponse send(io.github.pt9912.pgchangefeed.http.TransportRequest);
  }
  ```
  (das `internal fun interface` selbst kompiliert zu einem ganz gewöhnlichen
  `public interface`). Ein Java-Konsument dieses Jars — den `ADR-0109`
  §Verglichene Alternativen A ausdrücklich als Zielgruppe nennt
  („Kotlin-Artefakte sind Java-binärkompatibel, erreichen also zusätzlich
  reine Java-Konsumenten ohne eigenes Kotlin-Package") — kann
  `HttpTransport` implementieren und den zweiten Konstruktor direkt
  aufrufen; Java prüft Kotlins Modul-Metadaten nicht. Die Aussage stimmt
  nur für einen **Kotlin**-Konsumenten, der über den Kotlin-Compiler und
  Gradles Modul-Mechanismus kompiliert — nicht für „außerhalb dieses
  Moduls" im Allgemeinen, und nicht „ohne zusätzliche öffentliche
  API-Fläche" für den JVM-Bytecode als Ganzes. Das funktionale Testbarkeits-
  Ziel des Slice (DoD-Punkt „netzlos prüfbar") wird davon nicht berührt —
  die Tests laufen tatsächlich ohne Socket; betroffen ist ausschließlich
  die Zusatzbehauptung „adds no public API surface", die für die im
  selben ADR benannte Zielgruppe schlicht nicht zutrifft.
- `verifizierbar`: ja — `docker build -f sdks/kotlin/Dockerfile sdks/kotlin
  --target build`, danach `jar xf …jar … .class` und `javap -p` gegen die
  beiden genannten Klassen (Kommandos oben, real ausgeführt und
  protokolliert).
- `klasse`: Kotlin-`internal` fälschlich als JVM-Zugriffsschutz behauptet

## Negativbefunde

- geprüft, ohne Befund: alle neun `SPEC-018`-Fähigkeiten plus `readChanges`
  (`SPEC-022`) — Signatur, Request-/Response-Form und Fehler-Antwortform
  spiegeln die Spec exakt, keine fehlt, keine zusätzliche (kein Diagnose-/
  Health-Endpunkt — Plan §1 „Ausdrücklich NICHT in diesem Slice" korrekt
  eingehalten).
- geprüft, ohne Befund: `PgChangeFeedException`-Sieben-Klassen-Hierarchie
  konsistent mit der C#-Fassung — identische Namen, identische
  Statuscode-Zuordnung, identisches Fallback-/Malformed-Verhalten.
  Insbesondere ist der in `review-slice-sdk-csharp-http-client-flaeche.md`
  F-2 gefundene Fehler (rohe `JsonException` bei malformtem `2xx`-Body)
  in der Kotlin-Fassung von Anfang an vermieden —
  `PgChangeFeedHttpClient.kt:parseSuccessBody` fängt `JsonSyntaxException`
  bereits ab und wirft `PgChangeFeedMalformedResponseException`.
- geprüft, ohne Befund: Import-Grenze — `grep -rn "internal/\|cmd/"
  sdks/kotlin/` liefert nur einen Treffer im Gradle-Wrapper-Kommentar
  (`gradlew:60`, eine fremde GitHub-URL), keinen Import dieses Repos. Das
  Kotlin-Schlüsselwort `internal` selbst ist korrekt von dieser Prüfung
  unterschieden (F-1 behandelt es separat unter einem anderen Aspekt).
- geprüft, ohne Befund: Netzlosigkeit der Tests — `FakeHttpTransport`/
  `TestClientFactory` verdrahten den `internal`-Konstruktor mit einer
  reinen In-Memory-Antwortfunktion; kein Test öffnet einen Socket oder
  einen Loopback-Server (die verwendete Basis-Adresse
  `http://example.invalid:8080` wird nie real aufgelöst, weil
  `FakeHttpTransport.send()` nie an `java.net.http.HttpClient` delegiert).
- geprüft, ohne Befund: `offset` als `Long` statt `ULong` — die Grenze
  (obere `uint64`-Hälfte nicht darstellbar) ist im KDoc von
  `model/Consumers.kt` benannt und mit dem konkreten Gson-Grund belegt,
  konsistent mit den übrigen 64-bit-Feldern des SDK; für ein
  Pre-1.0-Release vertretbar.
- geprüft, ohne Befund: Backtick-Parität in allen drei geänderten
  Markdown-Dateien (gerade Anzahl je Datei, real nachgezählt).
- geprüft, ohne Befund: `AGENTS.md` §3.13 Träger-Nachzug — die im
  Plan-Nachzug behauptete README-Korrektur ist real im Diff vorhanden und
  korrekt; ein eigener, unabhängiger `grep` nach stehenden
  „folgt erst"/„follow in subsequent"/„coming soon"/„not yet
  implemented"-Aussagen über `sdks/kotlin/` liefert keinen weiteren
  Treffer.
- geprüft, ohne Befund: Out-of-Scope-Einhaltung (Plan §1) — kein
  `fun main(`, keine gRPC-/SSE-/NATS-Codepfade unter
  `sdks/kotlin/pgchangefeed-kotlin/src/`.
- geprüft, ohne Befund: `docs/user/benutzerhandbuch.md` — Version
  1.35→1.36 korrekt hochgezogen, neue Änderungshistorie-Zeile vorhanden,
  inklusive des in `ADR-0109` Festlegung 2 verlangten PAT-Hinweises
  (`read:packages`).
- geprüft, ohne Befund: `com.google.code.gson:gson:2.14.0` — real gegen
  Maven Central verifiziert, `<latest>`/`<release>` stimmen überein, keine
  Drift gegenüber der im Kommentar genannten Messung.
- geprüft, ohne Befund: Mutation-Testing-Behauptungen (Schritt 19) — beide
  genannten Testfälle existieren und ihre Assertions würden bei den
  beschriebenen Mutationen (Identitäts-`encode()`, vertauschte
  `403`-Zuordnung) nachweisbar rot werden; ein eigener zusätzlicher
  Mutations-Docker-Lauf wurde als nicht nötig bewertet, da die
  Testlogik die Aussage eindeutig trägt und der reguläre (unmutierte)
  Docker-Build real grün lief.
- geprüft, ohne Befund: `make gates` real gefahren, Exit-Code direkt
  (ungepiped) geprüft, `0`.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Kotlin-`internal` fälschlich als
JVM-Zugriffsschutz behauptet (neue Klasse, kein bereits gezähltes
Vorkommen aus dem Beobachtungs-Register).

## Verdikt

**Merge-blockierend:** ja — F-1 ist eine über zwei Träger (Produktionscode-
KDoc **und** Slice-Plan §3) wiederholte, mechanisch widerlegte Aussage
über eine von `ADR-0109` selbst benannte Zielgruppe (Java-Konsumenten).
Sie ist kein Sicherheitsloch und bricht keinen bestehenden Test, aber eine
Fixrunde ist nötig, weil ein künftiger Leser (Consumer, Implementer eines
Folge-Slice) sich sonst auf eine falsche API-Surface-Zusicherung verlässt:

1. `HttpTransport.kt`s KDoc auf die tatsächliche Grenze umschreiben —
   Kotlin-`internal` ist eine **compile-time**-Sichtbarkeit für andere
   Kotlin-Module über den Gradle-/Kotlin-Compiler-Mechanismus, keine
   JVM-Bytecode-Zugriffsbeschränkung; ein Java- oder reflection-basierter
   Aufrufer kann den Konstruktor und das Interface technisch weiterhin
   erreichen. Alternativ: die Grenze bewusst als benannte, akzeptierte
   Grenze formulieren (kein Sicherheitsproblem, nur kein vollständiger
   API-Surface-Schutz) statt als absolute Unsichtbarkeits-Zusage.
2. Denselben Wortlaut im Plan-Nachzug (§3) entweder korrigieren oder — da
   Slice-Pläne in `in-progress/` typischerweise nicht rückwirkend über
   bereits geschriebene Prosa hinaus verändert werden — den Punkt in der
   Closure-Notiz (§7) als Lerneintrag festhalten, damit er nicht in einen
   künftigen Plan unkorrigiert weiterwandert.

**Übergabe:** Rückmeldung an den Implementer-Agenten mit diesem Report.
Da eine Fixrunde nötig ist, bleibt die DoD-Checkbox „Review durchgeführt,
Report unter `docs/reviews/` liegt vor" im Slice-Plan **offen** — sie wird
regulär bei Schritt 21 des Implementer-Workflows nach der Fixrunde
nachgezogen (Skill-Regel „DoD-Checkbox-Nachzug ohne Fixrunde" greift hier
nicht). Dieser Report ersetzt keine Verifikation gegen die DoD — das bleibt
Verifier-Aufgabe (Modul 11).
