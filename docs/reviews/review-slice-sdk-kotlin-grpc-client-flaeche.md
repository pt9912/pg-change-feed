# Review-Report: slice-sdk-kotlin-grpc-client-flaeche — 2026-09-20

**Review-Art:** Code — geprüft gegen Plan
(`docs/plan/planning/in-progress/slice-sdk-kotlin-grpc-client-flaeche.md`,
inkl. DoD-Häkchen und §7 Closure-Notiz-Entwurf des Implementers), `ADR-0109`
(Accepted), `ADR-0060` (Accepted), `SPEC-020` und `AGENTS.md` Hard Rules
(Modul 10 §Drei Review-Arten). Kein DoD-Abgleich als eigenständige Prüfung —
das ist Verifier-Aufgabe (Modul 11); wo DoD-Text zitiert wird, dient das nur
der Einordnung des geprüften Diffs.

**Gegenstand:** Commit `473f3ee8` (feat(sdk): Kotlin-gRPC-Client-Fläche für
pgchangefeed-kotlin) gegen Elternstand `69c9e44b`
(`git diff 69c9e44b..473f3ee8`), Slice `slice-sdk-kotlin-grpc-client-flaeche`,
Welle `welle-sdk-kotlin-lh-fa-sst-009`.

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09"
(vier repo-spezifische HIGH-Regeln, seither erweitert um mehrere weitere
HIGH-Klassen, u. a. `AGENTS.md` §3.12/§3.13, „Beleg trägt seinen Satz
nicht", die Kotlin-`internal`-Lektion aus
`docs/reviews/review-slice-sdk-kotlin-http-client-flaeche.md` F-1).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-20.

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-sdk-kotlin-grpc-client-flaeche.md`
  (§1 Ziel/Abgrenzung, §2 DoD, §3 Plan, §6 Risiken samt Ausgängen, §7
  Closure-Notiz-Entwurf, §8 Sub-Area/Modus)
- `ADR-0109` (Accepted) — Festlegung 1/3/5, §Kontext (Ist-Stand-Messung,
  GitHub-Packages-Recherche, JDK-/Gradle-Kompatibilität), §Konsequenzen
  Folgepflicht 4
- `ADR-0060` (Accepted) — gRPC-Server-Streaming-Vertrag, vom SDK benutzt
- `spec/pflichtenheft.md` `SPEC-020` (Nachrichtenschema, zehn Felder,
  RPC-Name, Authn-Boundary)
- `AGENTS.md` §3.1 (Docker-only), §3.7 (Kommentar-Disziplin), §3.12
  (Herkunft von Aussagen), §3.13 (Träger-Nachzug)
- `harness/conventions.md` (MR-000 ID-Schema, MR-002 Slice-Kennungen)
- `sdks/csharp/PgChangeFeed.Client/Grpc/PgChangeFeedGrpcClient.cs` als
  Vergleichsfassung für den Transport-Seam-Ansatz
- `docs/reviews/review-slice-sdk-kotlin-http-client-flaeche.md` — Quelle
  des dort gefundenen HIGH F-1 (Kotlin-`internal` fälschlich als
  JVM-Zugriffsschutz behauptet) — tragende Referenz für Prüfpunkt 2 dieses
  Laufs

**Eigenständig durchgeführte Prüfungen (nicht nur Commit-Message/Plan-Text
übernommen):**

- Alle zehn `SPEC-020`-Nachrichtenfelder (`change_id`, `transaction_id`,
  `source_table_id`, `sequence`, `operation`, `old_image`, `new_image`,
  `schema_version`, `schema`, `table`) Feld für Feld gegen
  `PgChangeFeedGrpcClientMessageSchemaTest.kt` und
  `PgChangeFeedGrpcClient.kt` gehalten — Proto-Konstruktion im Test, Fluss
  durch `streamChanges()`, Assertions im Test: keine Abweichung, keine
  fehlende, keine zusätzliche Zuordnung.
- **Kotlin-`internal`-Semantik real erneut geprüft, nicht nur die
  KDoc-Aussage übernommen** (gelernte Lektion aus F-1 des vorigen Slice):
  `docker build --no-cache --build-context proto=proto -f
  sdks/kotlin/Dockerfile -t kotlin-grpc-review-test sdks/kotlin` real
  ausgeführt (Exit 0, `generateProto` real ausgeführt statt `NO-SOURCE`,
  `compileKotlin`/`compileJava`/`test`/`build` alle grün), Jar extrahiert,
  `javap -p` gegen `PgChangeFeedGrpcClient.class`, `GrpcStreamTransport.class`,
  `GeneratedStubTransport.class` im gepinnten `eclipse-temurin:21-jdk`-Image
  ausgeführt. Ergebnis deckt sich diesmal **exakt** mit der KDoc-Aussage
  (siehe Negativbefunde) — anders als beim vorigen Slice ist die
  `internal`-Kommentierung von Anfang an korrekt: sie behauptet nur die
  compile-time-Kotlin-Grenze, nicht mehr.
- **`COPY --from=proto`-Zielpfad-Fix real nachvollzogen:** derselbe
  `docker build --no-cache …`-Lauf (kein Cache-Hit) bestätigt, dass
  `generateProto` die `.proto` real findet und kompiliert — der im
  Implementer-Bericht (§6, §7) beschriebene, vor dem ersten Commit bereits
  korrigierte Doppelpfad-Fehler (`pgchangefeed-kotlin/pgchangefeed-kotlin/…`)
  tritt im committeten Dockerfile nicht auf.
- **`gradlew`-Grep-Treffer real gegen `git blame` gehalten:**
  `git blame -L 58,62 sdks/kotlin/pgchangefeed-kotlin/gradlew` zeigt
  Commit `7500fefb8` (2026-09-20 05:01, Projektgerüst-Commit) als
  Ursprung der Zeile — vor diesem Slice (Implementer-Commit
  06:57) entstanden, unverändert im Diff dieses Slice; kein neuer Import.
- **README-Lücke (Punkt 2 des Implementer-Berichts) real gegen das
  C#-/Python-Präzedens gehalten:** `git log --oneline -- README.md
  README.de.md` zeigt `1b3b4009` (C#) und `3be590f9` (Python) als die
  jeweiligen Commits, die die Root-README um die SDK-Zeile ergänzten —
  beide **nach** dem realen Publish (`1b3b4009`s Commit-Message: „das
  jetzt real auf NuGet.org veroeffentlichte … Package"; `3be590f9` ist die
  Python-Release-Realisierung selbst), nicht bereits beim jeweiligen
  HTTP-Client-Flächen-Slice. Für Kotlin existiert weder ein Pack-Werkzeug
  noch ein Publish-Workflow — das reale Präzedens deckt die
  Implementer-Einschätzung „vorbestehend, kein Zug dieses Slice",
  nicht das Gegenteil.
- Gradle-/gRPC-/Protobuf-Koordinaten real gegen Maven Central und das
  Gradle Plugin Portal nachverifiziert (`curl` gegen
  `maven-metadata.xml`/das Plugin-Portal-Metadatum, 2026-09-20):
  `io.grpc:grpc-bom` `1.84.0` ✓, `io.grpc:protoc-gen-grpc-java` `1.84.0`
  ✓, `io.grpc:grpc-kotlin-stub`/`io.grpc:protoc-gen-grpc-kotlin`: `latest`/
  `release` zeigen real einen Commit-Hash-Metadateneintrag
  (`6f774052d1d6923f8af2e0023886d69949b695ee`) — bestätigt die im
  `build.gradle.kts`-Kommentar behauptete Anomalie exakt, `1.5.0` bleibt
  die zuletzt echte Release-Version; `com.google.protobuf:protoc`/
  `protobuf-java`: `latest`/`release` zeigen real `21.0-rc-1` (ein
  Release-Candidate einer neuen Versionszählung, keine stabile Version) —
  `4.36.2` bestätigt als zuletzt reguläre stabile Version;
  `org.jetbrains.kotlinx:kotlinx-coroutines-core` `1.11.0` ✓;
  `com.google.protobuf`-Gradle-Plugin (Plugin Portal) `0.10.0` ✓. Zusätzlich
  `io.grpc:grpc-protobuf:1.84.0`s POM real abgerufen: transitive
  `protobuf-java`-Version `3.25.9`, bestätigt die im Kommentar behauptete
  BOM-Angleichungs-Begründung für die explizite Anhebung.
- Mutation-Testing-Behauptungen an der Testlogik nachvollzogen: eine
  `.catch { }`-Mutation in `streamChanges()` würde den `StatusException`
  aus dem Fake-Fehlerpfad verschlucken —
  `PgChangeFeedGrpcClientAuthBoundaryTest.kt`s `assertFailsWith<StatusException>`
  würde dann sichtbar fehlschlagen; eine
  `.map { it.toBuilder().clearSchema().build() }`-Mutation würde das
  `schema`-Feld vor der Ausgabe entfernen — die `assertEquals("public",
  received.schema)`-Zeile in
  `PgChangeFeedGrpcClientMessageSchemaTest.kt` würde dann sichtbar
  fehlschlagen. Beide Mutationen wurden anhand der Logik nachvollzogen,
  kein eigener Mutations-Docker-Lauf zusätzlich zum bereits real
  ausgeführten, unmutierten Bau gefahren (derselbe Prüfumfang wie im
  Vorgänger-Review).
- Out-of-Scope-Grep: kein `fun main(` unter
  `sdks/kotlin/pgchangefeed-kotlin/src/**/grpc/`, kein committetes `.proto`
  im SDK-Baum (`find sdks/kotlin -iname "*.proto"` liefert nichts), kein
  SSE-/NATS-Code (die wenigen „sse"/„nats"-Grep-Treffer sind Fehlalarme —
  `assertEquals`, KDoc-Wortbestandteile — einzeln gegengelesen), die
  gesendete `StreamChangesRequest` trägt
  `StreamChangesRequest.getDefaultInstance()` (kein Filter-Feld gesetzt,
  `SPEC-020`: „eine tabellen-granulare Filterung ist nicht Teil dieser
  Version").
- Import-Grenze: `grep -rn "internal/\|cmd/\|gen/" sdks/kotlin/` liefert
  genau zwei Treffer — der bereits oben per `git blame` bestätigte
  vorbestehende `gradlew`-Kommentar und ein KDoc-Zitat des
  Test-Vorbilds (`internal/adapters/driving/grpc/server_test.go`) in
  `PgChangeFeedGrpcClientMessageSchemaTest.kt` — beide kein Import.
- Backtick-Parität real nachgezählt (gerade Anzahl je Datei erwartet) für
  alle neun im Diff geänderten/neuen Dateien: Slice-Plan (510), Handbuch
  (1786), SDK-`README.md` (70), `Dockerfile` (36), `build.gradle.kts` (98),
  `PgChangeFeedGrpcClient.kt` (98), `FakeGrpcStreamTransport.kt` (16),
  `PgChangeFeedGrpcClientAuthBoundaryTest.kt` (20),
  `PgChangeFeedGrpcClientMessageSchemaTest.kt` (18) — alle neun gerade.
- `git show 473f3ee8 -- docs/user/benutzerhandbuch.md
  sdks/kotlin/pgchangefeed-kotlin/README.md` gelesen: Handbuch-Version
  1.36→1.37 mit neuer Änderungshistorie-Zeile, neuer `**SDK:**`-Absatz im
  gRPC-Abschnitt inkl. PAT-Hinweis, und die stehende „folgt in einem
  Folge-Release"-Aussage im bestehenden HTTP-Absatz korrigiert; dieselbe
  Korrektur im SDK-`README.md`s Status-Absatz (Englisch) — beide wie im
  Plan-Nachzug (§7 „Träger-Nachzug") beschrieben.
- `make gates` real ausgeführt, Exit-Code direkt (ungepiped) geprüft: `0`
  (`d-check: 861 Datei(en) geprüft, 0 Befund(e)`, `generated-sync: OK`,
  `a-check: gesamt: 0 Befund(e)`).

---

## Findings

Keine HIGH-, MEDIUM- oder LOW-Findings in diesem Lauf.

## Negativbefunde

- geprüft, ohne Befund: `SPEC-020`-Nachrichtenschema-Vollständigkeit — alle
  zehn Felder erreichen den Aufrufer über `streamChanges()` unverändert,
  Feld-für-Feld-Assertion in `PgChangeFeedGrpcClientMessageSchemaTest.kt`
  deckt jedes Feld einzeln, keine Lücke.
- geprüft, ohne Befund: **Kotlin-`internal`-Semantik korrekt beschrieben**
  — real per `docker build --no-cache` + `javap -p` gegen den frischen
  Bau geprüft: `internal constructor(transport: GrpcStreamTransport, …)`
  kompiliert zu einem gewöhnlichen `public`-Konstruktor
  (`public io.github.pt9912.pgchangefeed.grpc.PgChangeFeedGrpcClient(io.github.pt9912.pgchangefeed.grpc.GrpcStreamTransport, io.github.pt9912.pgchangefeed.PgChangeFeedClientOptions);`),
  `internal fun interface GrpcStreamTransport` zu einem gewöhnlichen
  `public interface`, `internal class GeneratedStubTransport` zu einer
  gewöhnlichen `public final class` — exakt das, was die KDoc in
  `PgChangeFeedGrpcClient.kt:98-108` behauptet (compile-time
  Kotlin-Compiler-Grenze, keine JVM-Bytecode-Zugriffsbeschränkung, ein
  Java-Aufrufer oder Reflection kann den Konstruktor weiterhin erreichen).
  Anders als beim vorigen Slice (F-1) gibt es hier **keine** Überzeichnung
  der Aussage — die Lektion wurde übernommen, nicht nur behauptet,
  übernommen zu haben.
- geprüft, ohne Befund: `COPY --from=proto`-Zielpfad im Dockerfile —
  realer `docker build --no-cache`-Lauf (kein Cache-Treffer) zeigt
  `generateProto` real ausgeführt (nicht `NO-SOURCE`),
  `compileKotlin`/`compileJava`/`test`/`build` alle grün; der im
  Implementer-Bericht beschriebene, vor dem ersten Commit korrigierte
  Doppelpfad-Fehler tritt im committeten Zustand nicht mehr auf.
- geprüft, ohne Befund: `gradlew`-Grep-Treffer (`internal/…`) ist real
  vorbestehend — `git blame` zeigt Commit `7500fefb8`
  (Projektgerüst-Commit, 2026-09-20 05:01) als Ursprung, vor dem
  Implementer-Commit dieses Slice (06:57) entstanden, unverändert im Diff.
- geprüft, ohne Befund: fehlende `sdks/kotlin/`-Zeile in Root-`README.md`/
  `README.de.md` ist real vorbestehend und **kein** Versäumnis dieses
  Slice — das reale C#-/Python-Präzedens (`1b3b4009`, `3be590f9`) zeigt,
  dass die Root-README erst beim realen Package-Publish nachgezogen wird,
  nicht bereits bei der Client-Flächen-Ebene; für Kotlin existiert weder
  Pack-Werkzeug noch Publish-Workflow. Die Implementer-Meldung an
  Reviewer/Koordinator statt stiller Mitänderung ist die korrekte Form.
- geprüft, ohne Befund: Gradle-/gRPC-/Protobuf-Koordinaten — alle sechs
  im DoD-Punkt genannten real gegen Maven Central/Gradle Plugin Portal
  nachgemessen (2026-09-20), keine Abweichung von der im
  `build.gradle.kts`-Kommentar dokumentierten Messung; die behauptete
  Commit-Hash-Anomalie bei `grpc-kotlin-stub`/`protoc-gen-grpc-kotlin` und
  die `21.0-rc-1`-Release-Candidate-Anomalie bei `protoc`/`protobuf-java`
  sind real reproduzierbar, keine Behauptung ohne Beleg (`AGENTS.md`
  §3.12).
- geprüft, ohne Befund: Mutation-Testing-Behauptungen (Authn-Boundary
  `.catch{}`, Schema-Vollständigkeit `.clearSchema()`) — beide Mutationen
  anhand der Testlogik nachvollzogen, beide würden die jeweils genannte
  Assertion sichtbar rot färben.
- geprüft, ohne Befund: Out-of-Scope-Einhaltung (Plan §1) — kein
  `fun main(`, kein committetes `.proto` im SDK-Baum, keine SSE-/
  NATS-Codepfade, keine tabellen-granulare Filterung
  (`StreamChangesRequest.getDefaultInstance()`, kein Replay-Mechanismus
  im Client.
- geprüft, ohne Befund: Import-Grenze — `grep -rn "internal/\|cmd/\|gen/"
  sdks/kotlin/` liefert nur die zwei bereits benannten Nicht-Import-Treffer
  (vorbestehender `gradlew`-Kommentar, KDoc-Zitat des Server-Test-Vorbilds).
- geprüft, ohne Befund: Kommentar-Disziplin (`AGENTS.md` §3.7) — keine
  Slice-/Wellen-Chronik im Produktionscode-KDoc von
  `PgChangeFeedGrpcClient.kt`; die einzigen slice-benannten Kommentare
  stehen in Testdateien und tragen die zulässige Testfall-Provenienz-Form
  („Rot färbende Mutation (real geprüft, …): … dieser Test schlägt dann
  fehl" — Satzsubjekt ist der Test, nicht der Produktionscode-Pfad).
- geprüft, ohne Befund: Backtick-Parität in allen neun geänderten/neuen
  Dateien (gerade Anzahl je Datei, real nachgezählt).
- geprüft, ohne Befund: `AGENTS.md` §3.13 Träger-Nachzug — Handbuch-Version
  1.36→1.37 samt Änderungshistorie-Zeile und PAT-Hinweis, SDK-`README.md`s
  Status-Absatz korrigiert; die stehende „folgt in einem Folge-Release"-
  Aussage ist in beiden Trägern beseitigt.
- geprüft, ohne Befund: `make gates` real gefahren, Exit-Code direkt
  (ungepiped) geprüft, `0`.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** keine — 0 Findings.

## Verdikt

**Merge-blockierend:** nein — 0 HIGH, 0 MEDIUM, 0 LOW. Alle vom
Implementer selbst benannten offenen Punkte (Dockerfile-Copy-Pfad-Fix,
`gradlew`-Grep-Treffer, README-Lücke) wurden unabhängig real nachgeprüft
und bestätigt, keiner davon verlangt eine Fixrunde. Insbesondere die aus
dem vorigen Slice gelernte Kotlin-`internal`-Lektion (F-1 in
`docs/reviews/review-slice-sdk-kotlin-http-client-flaeche.md`) wurde
diesmal von Anfang an korrekt angewendet — real per `javap -p` gegen einen
frischen, unge-cachten Docker-Build verifiziert, keine Wiederholung des
vorigen Findings.

**Übergabe:** Da keine Fixrunde nötig ist, zieht dieser Report die
DoD-Zeile „Review durchgeführt, Report unter `docs/reviews/` liegt vor" im
Slice-Plan selbst auf `[x]` nach (Reviewer-Skill §DoD-Checkbox-Nachzug
ohne Fixrunde), im selben Commit, der diesen Report anlegt. Dieser Report
ersetzt keine Verifikation gegen die DoD — das bleibt Verifier-Aufgabe
(Modul 11).
