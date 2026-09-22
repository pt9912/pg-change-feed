# Review-Report: slice-sdk-kotlin-sse-client-flaeche — 2026-09-22

**Review-Art:** Code — geprüft gegen Plan
(`docs/plan/planning/in-progress/slice-sdk-kotlin-sse-client-flaeche.md`),
`ADR-0109` (Accepted) und `AGENTS.md` Hard Rules (Modul 10 §Drei
Review-Arten). Kein DoD-Abgleich — das ist Verifier-Aufgabe (Modul 11).

**Gegenstand:** Diff-Range `baf6e8bd..75fc1f3f`, ein Commit (`75fc1f3f`,
„feat(sdks/kotlin): SSE-Stream-Client-Fläche"), Slice
`slice-sdk-kotlin-sse-client-flaeche`, Welle
`welle-sdk-kotlin-vollabdeckung`.

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09"
(vier repo-spezifische HIGH-Regeln, seither um mehrere weitere HIGH-Klassen
ergänzt — u. a. `AGENTS.md` §3.7/§3.12/§3.13).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-22.

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-sdk-kotlin-sse-client-flaeche.md`
  (§1 Ziel/Abgrenzung, §2 DoD, §3 Plan, §6 Risiken, §7 Closure-Notiz, §8
  Sub-Area/Modus)
- `ADR-0109` (Accepted) — Festlegung 1 (letzter Absatz, SSE als
  Folge-Package antizipiert), Festlegung 3 (räumliche Trennung von
  `examples/kotlin/`, Import-Grenze)
- `spec/pflichtenheft.md` §2 `SPEC-021` (Endpunkt, Event-Form,
  Nachrichtenschema, zehn Felder)
- `AGENTS.md` §3.1 (Docker-only), §3.7 (Kommentar-Disziplin), §3.12
  (Herkunft von Aussagen), §3.13 (Träger-Nachzug)
- `harness/conventions.md` (MR-000 ID-Schema, MR-002 Slice-Kennungen)
- `harness/README.md` §Sensors/§Minimal-Agent-Workflow
- `examples/kotlin/sse-client/SseStream.kt` als Draht-Vorbild
  (Frame-Zerlegung)
- `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/.../http/HttpTransport.kt`,
  `.../http/PgChangeFeedHttpClient.kt`, `.../grpc/GrpcStreamTransport.kt`,
  `.../grpc/PgChangeFeedGrpcClient.kt` als bestehende, bereits gereviewte
  Flächen (Transport-Seam-Muster, Auth-Header-Form, Exception-Hierarchie,
  Flow-Nutzung)
- `docs/reviews/review-slice-sdk-csharp-sse-client-flaeche.md` —
  Kontrastfolie für F-2/F-3/F-4 dieses Reports
- `docs/reviews/review-slice-sdk-kotlin-http-client-flaeche.md` — Herkunft
  des dort gefundenen HIGH F-1 (Kotlin-`internal` fälschlich als
  JVM-Bytecode-Schutz behauptet), heute erneut geprüft

**Eigenständig durchgeführte Prüfungen (nicht nur Commit-Message/DoD-Text
übernommen):**

- `docker build -f sdks/kotlin/Dockerfile sdks/kotlin --target build`
  **ohne** `--build-context proto=proto` real ausgeführt: Bau bricht
  sichtbar exakt an der `COPY --from=proto`-Zeile ab („failed to resolve
  source metadata for docker.io/library/proto:latest …"), Dockerfile-Zeile
  68 im Fehler-Trace benannt — die Plan-Behauptung ist real bestätigt.
- `docker build … --build-context proto=proto --target build --no-cache`
  eigenständig (Cache umgangen, kein Übernehmen des Implementer-Caches)
  ausgeführt: `BUILD SUCCESSFUL`, `13 actionable tasks: 13 executed`.
- JUnit-XML-Reports aus dem selbst gebauten Image extrahiert
  (`docker cp` aus einem `docker create`-Container, kein Bind-Mount) und
  je Datei `tests=`/`failures=`/`errors=` gezählt: zehn Testklassen,
  Summe `2+1+9+5+5+6+3+7+2+7 = 47`, `0` Fehler/Errors überall — die drei
  neuen SSE-Klassen tragen exakt `7` (`SseFrameParserTest`) + `7`
  (`PgChangeFeedSseClientAuthBoundaryTest`) + `2`
  (`PgChangeFeedSseClientMessageSchemaTest`) `=16`; die Implementer-Zahlen
  „47 gesamt, 16 neu (7+7+2)" sind real bestätigt, nicht nur übernommen.
- `pgchangefeed-kotlin-0.1.0.jar` aus demselben Image extrahiert und
  `javap -p` (im gepinnten `eclipse-temurin:21-jdk`-Image, Digest identisch
  zur Dockerfile-`FROM`-Zeile) gegen `SseTransport`, `JdkSseTransport` und
  `PgChangeFeedSseClient` ausgeführt: `SseTransport` kompiliert zu einem
  gewöhnlichen `public interface`, `JdkSseTransport` zu einer `public final
  class`, und der `internal`-Konstruktor von `PgChangeFeedSseClient`
  (`SseTransport, PgChangeFeedClientOptions`) ist im Bytecode ein
  gewöhnlicher `public`-Konstruktor — genau das, was `SseTransport.kt`s
  KDoc behauptet („ordinary public symbols … a Java caller … can still
  implement `SseTransport` and invoke that constructor directly"). Anders
  als beim HIGH-Finding des vorigen Kotlin-Slice
  (`review-slice-sdk-kotlin-http-client-flaeche.md` F-1) trägt die
  KDoc-Aussage hier von Anfang an korrekt „compile-time Kotlin-Grenze,
  keine JVM-Bytecode-Schranke" — der Steering-Loop-Lerneintrag aus dem
  vorigen Slice ist wirksam angekommen, real durch eigenen `javap`-Lauf
  bestätigt, nicht nur aus der Closure-Notiz übernommen.
- Drei Mutationen real gesetzt, gebaut (`--no-cache` bzw. Cache-Invalidierung
  über den geänderten Quelltext) und wieder zurückgesetzt (Arbeitsbaum nach
  jedem Schritt über `git diff --stat`/`git status --short` auf Sauberkeit
  geprüft):
  1. `SseFrameParser.readFrame`: ein zusätzliches `return null` unmittelbar
     vor `return SseFrame(...)` eingefügt (identisch zur im Testkommentar
     beschriebenen Mutation) → `47 tests completed, 9 failed` — deckt sich
     exakt mit der erwarteten Frame-Parsing-Rotfärbung.
  2. `PgChangeFeedSseClient.buildException`s `401`-Zweig auf
     `PgChangeFeedForbiddenException` geändert → `47 tests completed, 1
     failed` (`missing or unknown token throws Unauthorized`) — exakt
     die erwartete Auth-Boundary-Rotfärbung.
  3. `streamChanges()`s `yield(parseChange(...))` um `.copy(schema = "")`
     ergänzt (identisch zur im Testkommentar beschriebenen Mutation) →
     `47 tests completed, 1 failed`
     (`streamChanges yields all SPEC-021 fields unchanged`) — exakt die
     erwartete Nachrichtenschema-Rotfärbung. Alle drei Mutationen sind
     damit real reproduziert, nicht nur aus dem Testkommentar geglaubt;
     die Commit-Message selbst nennt nur zwei der drei
     („Frame-Parsing-Logik und die 401-Auth-Boundary") — die dritte
     (Nachrichtenschema) steht nur im Testkommentar, siehe F-3.
- `grep -rn "internal/\|cmd/\|gen/" sdks/kotlin/ --include="*.kt"
  --include="*.md" --include="*.kts"` selbst gefahren: zwei Treffer, beide
  Doku-Kommentar-Zitate des Go-Testvorbilds
  (`PgChangeFeedSseClientMessageSchemaTest.kt:14` zitiert
  `internal/adapters/driving/http/sse_test.go`, sowie ein bereits
  bestehender, außerhalb dieses Diffs liegender Treffer im gRPC-Testpaket)
  — kein Import, kein Treffer in den neuen SSE-Produktionsdateien selbst.
- Alle zehn `SPEC-021`-Felder aus `spec/pflichtenheft.md` §`SPEC-021`
  einzeln gegen `sse/model/Change.kt`s `@SerializedName`-Annotationen
  gehalten: `change_id` (String), `transaction_id` (String),
  `source_table_id` (String), `sequence` (Long/int64), `operation`
  (String), `old_image`/`new_image` (`JsonElement?`, eingebettetes JSON
  oder `null`), `schema_version` (String), `schema` (String), `table`
  (String) — alle zehn vorhanden, keins fehlt, keins zusätzlich, Typen
  passen zur Spec-Tabelle.
- `SseFrameParser.readFrame` Zeile für Zeile gegen
  `examples/kotlin/sse-client/SseStream.kt`s `readEvent` gehalten:
  identisches Grenzverhalten (Leerzeile beendet ein Frame, führende
  Leerzeilen werden übersprungen, ein am Quellende unvollständiges Frame
  liefert `null`) — durch die sieben `SseFrameParserTest`-Fälle real
  abgesichert und durch Mutation 1 oben selbst nachvollzogen.
- Server-Gegenstück gelesen (`internal/adapters/driving/http/sse.go`):
  identisch zum C#-Review bereits festgestellt — der Server schreibt nie
  mehr als eine `data:`-Zeile pro Event; `SseFrameParser.kt`s
  `data = line.removePrefix(DATA_PREFIX)` überschreibt statt zu verketten,
  bei mehreren `data:`-Zeilen ginge Inhalt verloren, aber dieser Pfad wird
  vom aktuellen Server nie ausgelöst (siehe F-1).
- `buildException`/`extractErrorMessage` in
  `sse/PgChangeFeedSseClient.kt:127-142` Zeile für Zeile gegen
  `http/PgChangeFeedHttpClient.kt:180-197`s gleichnamige private Methoden
  gehalten: identischer `when`-Ausdruck (400/401/403/404/500 → je ein
  Exception-Typ, sonst `PgChangeFeedUnexpectedStatusException`),
  identisches `extractErrorMessage` — eine vollständige, wörtliche Kopie
  (siehe F-2).
- `git diff baf6e8bd..75fc1f3f -- sdks/kotlin/pgchangefeed-kotlin/README.md`
  gelesen: die §Status-Zeile „SSE and NATS-vollinhalt delivery remain out
  of scope for this package's planned first full release" ist ersetzt
  durch einen Satz, der die SSE-Fläche jetzt als Teil des Release-Standes
  nennt und ausschließlich NATS-Vollinhalt als offen belässt. Zusätzlich
  `git show baf6e8bd:sdks/kotlin/pgchangefeed-kotlin/README.md | grep -in
  sse` gegen den Elternstand gefahren: keine zweite, veraltete
  SSE-Uncovered-Aussage irgendwo im Dokument (nur die eine korrigierte
  Stelle nannte SSE als „out of scope").
- `git diff --stat baf6e8bd..75fc1f3f -- '*nats*' '*Nats*'
  'tools/harness/sdk-pack*' '.github/workflows/*sdk-kotlin*'
  'spec/architecture.md' '.a-check.yml' 'spec/pflichtenheft.md'
  'docs/user/benutzerhandbuch.md'` — leer: keiner der acht Bereiche
  berührt, wie im Plan §1 vorab benannt (Doku-/NATS-/Pack-Nachzug bewusst
  dem Folge-Slice vorbehalten). Voller `git diff --stat` zusätzlich
  gegengelesen: nur elf Dateien (Plan, README, sieben neue `.kt`-Dateien
  unter `sse/`, vier davon Test) berührt.
- `grep -n "^version" sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts` —
  `version = "0.1.0"`, real unverändert gegenüber dem Elternstand
  (`git diff baf6e8bd..75fc1f3f -- build.gradle.kts` liefert keinen Treffer
  für diese Datei).
- Neue Dateien nach `slice-`/`welle-`-Nennungen durchsucht (`AGENTS.md`
  §3.7): drei Treffer, alle drei in Testdateien als
  „Rot färbende Mutation (real geprüft, slice-sdk-kotlin-sse-client-flaeche): …"
  — Satzsubjekt ist in allen drei Fällen die Mutation/der Testfall, nicht
  ein Produktionscode-Pfad; das ist die vom Skill ausdrücklich zugelassene
  Testfall-Provenienz-Form, kein HIGH.
- Backtick-Parität real nachgezählt für beide in diesem Diff geänderten
  Markdown-Dateien: `slice-sdk-kotlin-sse-client-flaeche.md` → 252
  (gerade), `sdks/kotlin/pgchangefeed-kotlin/README.md` → 84 (gerade) —
  beide paarig.
- `docs/plan/planning/observations/` real per `git diff --stat` geprüft:
  keine Änderung — deckt sich mit der Closure-Notiz-Aussage „keine neue
  Beobachtung angefallen".
- Design-Entscheidung `Sequence<Change>` statt `Flow<Change>` gegen
  `http/PgChangeFeedHttpClient.kt` (rein blockierende, nicht-`suspend`
  Methoden, `grep -n "suspend" …PgChangeFeedHttpClient.kt` liefert keinen
  Treffer) und `grpc/PgChangeFeedGrpcClient.kt` (generierter
  Coroutine-Stub, echtes `suspend`/`Flow`) gehalten: Die KDoc-Begründung
  ist stichhaltig — `java.net.http.HttpClient.send()` ist ein gewöhnlicher
  blockierender Aufruf ohne nicht-blockierendes Pendant, exakt wie jede
  andere Methode der HTTP-Fläche dieses Packages; ein `flow { }`-Builder
  um diesen Aufruf bräuchte entweder einen expliziten
  `flowOn(Dispatchers.IO)`-Sprung oder würde Suspendierungspunkte
  suggerieren, die dieser Transport nicht bietet. `Sequence` ist damit dem
  synchronen Ausführungsmodell der bestehenden HTTP-Fläche **näher** als
  `Flow` es wäre — die gRPC-Fläche ist hier die Ausnahme (echter
  suspendierender Stub), nicht die Norm, an der SSE sich hätte ausrichten
  müssen. Der Plan öffnete diese Tür ausdrücklich („oder gleichwertig
  idiomatisch"); die Wahl ist begründet, keine unbegründete Abweichung
  (siehe F-4 zur benannten Kehrseite: drei Container-Typen über drei
  Surfaces).
- `make gates` ungefiltert laufen lassen, Exit-Code direkt (kein Pipe/
  Wrapper dazwischen) geprüft: `0`.

---

## Findings

### F-1 — `SseFrameParser` überschreibt statt verkettet mehrere `data:`-Zeilen

- `kategorie`: LOW
- `quelle`: Maintainability (SSE-Spezifikation, RFC-Konvention „mehrere
  aufeinanderfolgende `data:`-Zeilen werden mit `\n` verkettet")
- `pfad`: `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/io/github/pt9912/pgchangefeed/sse/SseFrameParser.kt:40`
- `befund`: `data = line.removePrefix(DATA_PREFIX)` überschreibt den
  bisher gesammelten Wert bei jeder weiteren `data:`-Zeile, statt die
  Zeilen zu verketten — bei mehreren `data:`-Zeilen in einem Frame ginge
  der vorherige Inhalt verloren. Der aktuelle Server
  (`internal/adapters/driving/http/sse.go`) schreibt nie mehr als eine
  `data:`-Zeile pro Event, und das gelesene Draht-Vorbild
  (`examples/kotlin/sse-client/SseStream.kt`) trägt exakt dasselbe
  Verhalten — kein neu eingeführter Fehler, sondern ein unverändert
  übernommenes, bislang folgenloses Verhalten; identisch zum bereits im
  C#-Sibling gefundenen F-1
  (`review-slice-sdk-csharp-sse-client-flaeche.md`).
- `verifizierbar`: ja — ein Test, der zwei `data:`-Zeilen in einem Frame
  sendet, würde den Verkettungs-Fall rot zeigen.
- `klasse`: „Vorbild-Verhalten unverändert übernommen, Spec-Lücke bleibt
  latent"

### F-2 — `buildException`/`extractErrorMessage` wörtlich aus `PgChangeFeedHttpClient.kt` kopiert

- `kategorie`: LOW
- `quelle`: Maintainability (DRY)
- `pfad`: `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/io/github/pt9912/pgchangefeed/sse/PgChangeFeedSseClient.kt:127-142`
  vs. `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/io/github/pt9912/pgchangefeed/http/PgChangeFeedHttpClient.kt:180-197`
- `befund`: Beide private Methoden sind funktional identisch (derselbe
  `when`-Ausdruck über 400/401/403/404/500, dieselbe
  `ErrorResponse`-Fallback-Logik) — eine künftige Änderung am
  `SPEC-018`-Fehler-Mapping müsste an beiden Stellen nachgezogen werden,
  ohne dass ein Compiler- oder Testfehler das an der zweiten Stelle
  erzwingt. Dieselbe Einordnung wie beim C#-Sibling (F-2,
  `review-slice-sdk-csharp-sse-client-flaeche.md`): Der Implementer
  benennt die Kopie in der Closure-Notiz selbst und begründet sie mit dem
  bereits dort etablierten Muster — die Bewertung „LOW, kein Blocker"
  trifft hier ebenfalls zu, weil kein Gate Code-Duplikation in diesem
  Repo prüft und die Wartungslast bislang nur hypothetisch ist (noch keine
  zweite Änderung an diesem Mapping seit dessen Entstehung im
  HTTP-Client-Slice). Anders als beim C#-Sibling ist hier zusätzlich zu
  bemerken, dass die Duplikation jetzt zum **zweiten** Mal auftritt
  (gRPC-Fläche bildet `Grpc.Core.StatusCode`/Kotlin-Äquivalent strukturell
  anders ab, kein drittes Vorkommen) — bei einem dritten Vorkommen in
  diesem Package (z. B. einer künftigen NATS-Vollinhalts-Fläche mit
  demselben HTTP-Fehler-Mapping) würde die Skill-Regel „Wiederholung eines
  Musters, das schon zweimal LOW war" (MEDIUM) greifen.
- `verifizierbar`: nein — kein Gate prüft Code-Duplikation in diesem Repo.
- `klasse`: „Wörtliche Kopie statt gemeinsamer Extraktion bei identischer
  Logik"

### F-3 — Commit-Message nennt nur zwei der drei real durchgeführten Mutationen

- `kategorie`: INFO
- `quelle`: Maintainability (Vollständigkeit der Selbstauskunft)
- `pfad`: Commit-Message `75fc1f3f` vs.
  `sdks/kotlin/pgchangefeed-kotlin/src/test/kotlin/io/github/pt9912/pgchangefeed/sse/PgChangeFeedSseClientMessageSchemaTest.kt:8-13`
- `befund`: Die Commit-Message nennt „Mutation-Testing real gebaut und rot
  gesehen für die Frame-Parsing-Logik und die 401-Auth-Boundary" — zwei
  von drei tatsächlich im Diff dokumentierten und (in diesem Review real
  reproduzierten) Mutationen. Die dritte, im Testkommentar von
  `PgChangeFeedSseClientMessageSchemaTest.kt` beschriebene Mutation
  (`.copy(schema = "")`) fehlt in der Commit-Message-Aufzählung. Kein
  Beleg-Verstoß im Sinne von „Beleg trägt seinen Satz nicht" — die
  Commit-Message behauptet keine Vollständigkeit („für die … und die …"
  ist keine geschlossene Aufzählung „alle drei") —, aber eine
  unvollständige Selbstauskunft an der Stelle, die ein Leser ohne
  Diff-Studium als vollständige Zusammenfassung läse.
- `verifizierbar`: nein — reine Text-Vollständigkeitsfrage, kein Gate.
- `klasse`: „Commit-Message zählt weniger Belege als der Diff trägt"

### F-4 — Drei Zustellwege, drei verschiedene Rückgabe-Container-Typen (`Flow`, `Sequence`, direkter Wert)

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/io/github/pt9912/pgchangefeed/sse/PgChangeFeedSseClient.kt:39-58`
  (KDoc-Begründung) vs. `grpc/PgChangeFeedGrpcClient.kt` (`Flow<Change>`)
  vs. `http/PgChangeFeedHttpClient.kt` (direkter Rückgabewert, keine
  Suspendierung)
- `befund`: Das Package liefert jetzt drei unterschiedliche
  Rückgabeformen für strukturell ähnliche „liefere mehrere/eine
  Change(s)"-Aufrufe: `Flow<Change>` (gRPC, echtes Suspend-Primitive),
  `Sequence<Change>` (SSE, kalt, aber synchron/blockierend) und ein
  direkter `List<Change>`-artiger Rückgabewert (HTTP `readChanges`, kein
  Streaming). Die KDoc-Begründung in `PgChangeFeedSseClient` ist
  stichhaltig (siehe Prüfungsblock oben) und referenziert `ADR-0109`
  Festlegung 1 ausdrücklich als „unifies only the shared connection
  denominator … not each surface's execution model" — die
  Inkonsistenz ist also bereits benannt und begründet abgegrenzt, kein
  unbemerkter Drift. Reiner Hinweis: Ein künftiger Consumer, der zwischen
  SSE und gRPC wechselt, wechselt auch den Container-Typ, den er
  konsumiert.
- `verifizierbar`: nein — keine Gegenprobe ohne eine hypothetische
  vierte, vereinheitlichte Form.
- `klasse`: „Begründete Inkonsistenz über parallele Oberflächen"

## Negativbefunde

- geprüft, ohne Befund: `SPEC-021`-Vollständigkeit — alle zehn
  Nachrichtenfelder Feld für Feld gegen `spec/pflichtenheft.md` §`SPEC-021`
  gehalten, `sse/model/Change.kt` deckt sie exakt (Namen und Typen).
- geprüft, ohne Befund: Kotlin-`internal`-Semantik — real per `javap -p`
  gegen ein selbst, unabhängig (`--no-cache`) gebautes Jar geprüft;
  `SseTransport.kt`s KDoc trägt die korrekte Aussage (compile-time
  Kotlin-Grenze, keine JVM-Bytecode-Schranke) von Anfang an — anders als
  beim HIGH-Finding des vorigen Kotlin-Slice, hier real bestätigt ohne
  Abweichung zwischen Aussage und Bytecode.
- geprüft, ohne Befund: Design-Entscheidung `Sequence<Change>` statt
  `Flow<Change>` — die KDoc-Begründung ist stichhaltig (blockierender
  `HttpClient.send()`, keine Suspend-Primitive in der zugrunde liegenden
  HTTP-Fläche) und konsistenter mit der übrigen (synchronen) HTTP-Fläche
  dieses Packages als ein `flow { }`-Wrapper es gewesen wäre; der Plan
  ließ diese Form ausdrücklich zu (siehe F-4 für die benannte, aber
  begründete Kehrseite über drei Surfaces).
- geprüft, ohne Befund: Exception-Mapping-Duplikation — dieselbe
  Einordnung wie beim C#-Sibling (LOW, kein Blocker) trifft zu (siehe
  F-2 für die Abweichung: zweites statt erstes Vorkommen dieses Musters
  in diesem Package).
- geprüft, ohne Befund: Testzahlen — eigenständig aus den JUnit-XML-Reports
  eines selbst, unabhängig gebauten Images gezählt: 47 Tests gesamt, 16
  neu (7+7+2), 0 Fehler — deckt sich exakt mit der Implementer-Behauptung.
- geprüft, ohne Befund: Mutation-Testing-Behauptungen — alle drei im Diff
  dokumentierten Mutationen (Frame-Parsing, 401-Auth-Boundary,
  Nachrichtenschema) real gesetzt, gebaut und rot gesehen: `9/47`, `1/47`,
  `1/47` — exakt reproduziert, nicht nur geglaubt (siehe F-3 zur
  unvollständigen Commit-Message-Aufzählung).
- geprüft, ohne Befund: `sdks/kotlin/pgchangefeed-kotlin/README.md`-Korrektur
  — die stehengebliebene „SSE and NATS-vollinhalt delivery remain out of
  scope"-Aussage ist vollständig und korrekt behoben (SSE jetzt als
  gedeckt benannt, NATS bleibt offen benannt); kein Rest der veralteten
  Aussage an anderer Stelle des Dokuments.
- geprüft, ohne Befund: Version-Bump — `build.gradle.kts`s `version` real
  unverändert `"0.1.0"`.
- geprüft, ohne Befund: `--build-context proto=proto`-Zwang — real per
  eigenem `docker build` ohne das Flag reproduziert: Bau bricht exakt an
  der `COPY --from=proto`-Zeile ab; drittes reales Auftreten dieser
  Docker-Bau-Kopplung, im Plan selbst als 2×-Fund benannt und bewusst
  vermieden statt ein drittes Mal als neue Beobachtung ausgelöst.
- geprüft, ohne Befund: Import-Grenze (`internal/**`/`cmd/**`/`gen/**`) —
  eigener `grep`-Lauf über `sdks/kotlin/`, zwei Treffer außerhalb der
  neuen SSE-Dateien (Doku-Zitate des Go-Testvorbilds), keiner in den neuen
  Produktionsdateien.
- geprüft, ohne Befund: Out-of-Scope — kein Bezug zu NATS-Flächen, keine
  Änderung an einem Pack-Werkzeug- oder Publish-Workflow-Skript, kein
  Bezug zu `spec/architecture.md`, `.a-check.yml`, `spec/pflichtenheft.md`
  oder `docs/user/benutzerhandbuch.md` (`git diff --stat` gegen die volle
  Range bestätigt: nur Plan, README und `sse/**` berührt).
- geprüft, ohne Befund: `docs/plan/planning/observations/` — keine
  Änderung; deckt sich mit der Closure-Notiz-Aussage „keine neue
  Beobachtung angefallen", der 2×-Zähler
  (`BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut`) bleibt
  unverändert bei 2×.
- geprüft, ohne Befund: Kommentar-Disziplin (`AGENTS.md` §3.7) — kein
  Konjunktiv über eine verworfene Alternative, kein abwesender Text, keine
  unzulässige Slice-/Wellen-Chronik in Produktionscode; die drei
  `slice-sdk-kotlin-sse-client-flaeche`-Nennungen stehen ausschließlich in
  Testdateien mit dem Testfall/der Mutation als Satzsubjekt (zulässige
  Provenienz-Form).
- geprüft, ohne Befund: Backtick-Parität — beide in diesem Diff geänderten
  Markdown-Dateien real nachgezählt, beide paarig (252, 84).
- geprüft, ohne Befund: Traceability — Commit-Betreff nennt
  `LH-FA-SST-009`/`ADR-0109`, kein `SPEC-*`/`ARC-*` im Betreff.
- geprüft, ohne Befund: `make gates` real gefahren, Exit-Code direkt
  (ungepiped) geprüft, `0`.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 2 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** „Vorbild-Verhalten unverändert
übernommen, Spec-Lücke bleibt latent" · „Wörtliche Kopie statt
gemeinsamer Extraktion bei identischer Logik" · „Commit-Message zählt
weniger Belege als der Diff trägt" · „Begründete Inkonsistenz über
parallele Oberflächen"

## Verdikt

**Merge-blockierend:** nein — 0 HIGH, 0 MEDIUM, zwei LOW und zwei INFO,
keines davon mit Sicherheits-, Korrektheits- oder ADR-Verstoß-Charakter.
Insbesondere: die aus dem vorigen Kotlin-Slice gelernte HIGH-Lektion
(`internal`-Sichtbarkeit fälschlich als JVM-Bytecode-Schutz behauptet) ist
in diesem Diff von Anfang an korrekt umgesetzt — real per `javap -p`
gegen ein unabhängig gebautes Jar bestätigt, kein Rückfall. Keine
Fixrunde am Implementer nötig; die vier Findings sind Hinweise für eine
künftige Gelegenheit (F-2s Duplikations-Zähler ist bei 2× — ein drittes
Vorkommen würde die Skill-Regel „zweimal LOW → MEDIUM" auslösen), keine
Blocker dieses Slice.

**DoD-Checkbox-Nachzug ohne Fixrunde** (Skill-Regel,
`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne Fixrunde): Da
dieses Verdikt zu keiner Fixrunde führt, wird die DoD-Zeile „Review
durchgeführt, Report unter `docs/reviews/` liegt vor" im Slice-Plan
(`docs/plan/planning/in-progress/slice-sdk-kotlin-sse-client-flaeche.md`)
im selben Commit, der diesen Report anlegt, auf `[x]` nachgezogen, mit
Verweis auf diesen Report-Pfad.

**Übergabe:** kein Rückgabe-Pfeil an den Implementer nötig. Die vier
Findings (zwei LOW, zwei INFO) gehen in die Slice-Closure §7 und von dort
in den Steering-Loop-Zähler — insbesondere F-2s Duplikations-Zähler
(2. Vorkommen in diesem Package). Dieser Report ist ein Lauf-Beleg; er
ersetzt keine Verifikation gegen die volle DoD — das bleibt
Verifier-Aufgabe (Modul 11).
