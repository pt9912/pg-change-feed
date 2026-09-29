# Verifikations-Report: slice-sdk-kotlin-grpc-administration-flaeche — 2026-09-29

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich +
Entscheidungs-Konformität + Plan-vs-Code-Diff + Gates. Review-Artefakte:
[`review-sdk-kotlin-grpc-administration-flaeche.md`](review-sdk-kotlin-grpc-administration-flaeche.md)
(1 HIGH F-1, 2 LOW F-2/F-3, 1 INFO F-4),
[`review-fixrunde-welle-sdk-grpc-administration-flaeche.md`](review-fixrunde-welle-sdk-grpc-administration-flaeche.md)
(0 HIGH/MEDIUM, Fixrunde bestätigt).

**Gegenstand:** Slice-Plan
[`slice-sdk-kotlin-grpc-administration-flaeche.md`](../plan/planning/done/slice-sdk-kotlin-grpc-administration-flaeche.md)
(Welle `welle-sdk-grpc-administration-flaeche`), Basis `76afcad4`, drei
Slice-Commits:

- `74b41657` — feat(sdk): Kotlin-Administration-Client und Stream-Filter (`LH-FA-SST-009`)
- `fb943092` — docs(sdk): Kotlin-SDK-README für Administration-Client nachgezogen (`LH-FA-SST-009`)
- `b239d849` — plan(slice): DoD, Suchlauf und Closure-Notiz (`LH-FA-SST-009`)

und die Fixrunde `bd10c391` (Kotlin: zwei Fixture-Strings übersetzt —
`git show bd10c391` zeigt beide Zeilenersetzungen, `grep -n "Rechtsklasse\|interner
Fehler" …ErrorMappingTest.kt` liefert Exit 1) samt `48e04899`/`fdb52c2b`
(C#-Altbestand, Fixrunden-Review).

Nebenläufigkeit: der Python-Slice derselben Welle lief parallel im selben
Arbeitsbaum — die dadurch entstandenen Überlagerungen (Handbuch-Zeilen 1.80/1.81
in einem Commit, Ein-Zeilen-Link-Fix in `verifikation-sdk-csharp-…`) sind vom
Ausgangs-Review (F-4, INFO) und vom Fixrunden-Review transparent berichtet und
hier nicht Gegenstand der Prüfung.

## 1. Eigene Sensor-Belege (dieser Lauf, Exit-Code je ungepiped direkt gesichert — [`AGENTS.md`](../../AGENTS.md) §3.9)

| Sensor | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `make gates` | **Exit 0** | `baseline-verify: v6.9.0 OK — 54 Dateien` · `d-check: 1405 Datei(en) geprüft, 0 Befund(e)` · `commit-traceability: OK — 5 Commit(s) …, Betreffs ohne Struktur-ID` · `coverage-gate: OK — Coverage 80.40% erfüllt Schwelle 80%` · `generated-sync: OK` (beide `.proto`-Quellen, vier Generate-Dateien) · `gesamt: 0 Befund(e)` (a-check) |
| `make sdk-pack-kotlin` | **Exit 0** | Docker-Bau lief gecacht durch (`RUN ./gradlew --no-cache test` → `CACHED`); deshalb kein Test-Beleg, sondern nur Bau-Konsistenz — Test-Beleg im No-Cache-Lauf (nächste Zeile) |
| `docker build --no-cache --build-context proto=proto --target build` (eigenes Image `pg-change-feed:verifier-kotlin-nocache`) | **Exit 0** | `> Task :compileKotlin` … `> Task :compileTestKotlin` … `> Task :test` · `BUILD SUCCESSFUL in 31s` · `13 actionable tasks: 13 executed`; kein „Unresolved reference 'administration'" — die zweite `COPY --from=proto`-Zeile (`sdks/kotlin/Dockerfile` Zeile 60, `administration.proto`) trägt real |
| Testresultate aus dem No-Cache-Image (`grep "<testsuite"` über `build/test-results/test/*.xml`) | **19 Suites, 89 Testfälle, 0 failures, 0 errors** | gRPC-Paket: ConsumerTest 7 · DiagnoseTest 3 · ErrorMappingTest 6 · RetentionAndChangesTest 4 · TableTest **6** (dazu V-1) · GrpcClientFilterTest 3 (inkl. `streamChanges without arguments sends an empty request`) · GrpcClientAuthBoundaryTest 2 · GrpcClientMessageSchemaTest 1 |
| `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-sdk-kotlin-grpc-administration-flaeche.md` | **Exit 0** | `suchlauf-nachmessen: 6 Zeilen stimmen` (Parent `fd39b68b` und `diff` je nachgemessen, alle sechs Zeilen `OK`) |
| `grep -riE "SPEC-\|ADR-\|ARC-\|LH-FA-\|LH-QA-" sdks/kotlin/` | **0 Treffer** | Exit 1 = kein Match (`sdk-public-doc-check`-Gegenstand, selbst gefahren) |

Nicht gefahren: `make test-sdk-kotlin-integration` (Realserver-E2E ist nicht
Scope dieses Slice — Plan §7 nennt es ausdrücklich als Gegenstand eines eigenen
Folge-Slices); `make test` (Go-Fläche nicht berührt; der Diff enthält keine
Go-Quelldatei, die Generaten prüft `generated-sync`).

## 2. DoD — Verdikt je Zeile (§2 des Plans)

Gezählt: 8 `[x]`-Zeilen, 1 `[ ]`-Zeile.

| # | DoD-Zeile | Verdikt | Realer Beleg |
|---|---|---|---|
| 1 | `LH-FA-SST-009` erfüllt: `PgChangeFeedAdministrationClient` deckt alle elf RPCs, versiegelte Fehlerklasse, Unit-Tests je RPC | **bestätigt** | `PgChangeFeedAdministrationClient.kt` trägt elf `suspend fun` (registerConsumer, acknowledgeConsumer, getConsumerPosition, removeConsumer, enableTable, disableTable, getTableStatus, listTables, runRetention, readChanges, diagnose) 1:1 gegen `proto/cdc/administration/v1/administration.proto` (elf `rpc`-Zeilen 222–232, selbst gezählt); `PgChangeFeedGrpcException`-Hierarchie: `sealed class` mit fünf typisierten Unterklassen + `PgChangeFeedGrpcUnexpectedStatusException`-Fallback, Mapping über `mapException`; 32 laufende Testfälle im gRPC-Paket (No-Cache-Lauf, §1) — je RPC mindestens ein laufender Test; eine der beanspruchten negativen Tabellen-Abdeckungen läuft nicht (V-1) |
| 2 | `ADR-0133` erfüllt: `streamChanges()` mit optionalen `schema`/`table`-Parametern, Regressionstest | **bestätigt** | `PgChangeFeedGrpcClient.streamChanges(schema: String? = null, table: String? = null)`; Regressionstest `streamChanges without arguments sends an empty request` läuft (FilterTest: 3 Testfälle, 0 failures); Felder werden nur bei nicht-`null` gesetzt — leere `StreamChangesRequest{}` = ungefiltert |
| 3 | `make gates` grün | **bestätigt** | eigener Lauf, Exit 0 (§1) |
| 4 | Review durchgeführt, Report unter `docs/reviews/`, kein Self-Review | **bestätigt** | beide Report-Pfade lösen auf (Ausgangs-Review mit 1 HIGH; Fixrunden-Review 0/0); F-1 (HIGH) durch `bd10c391` geschlossen — beide übersetzten Strings im Baum, grep Exit 1 (§1) |
| 5 | `docs/user/benutzerhandbuch.md`: beide gRPC-Abschnitte nennen Kotlin mit voller Fläche; `sdks/kotlin/pgchangefeed-kotlin/README.md` nachgezogen | **bestätigt** | Handbuch: Stream-Absatz nennt `pgchangefeed-kotlin` (`streamChanges(schema, table)`, beide Parameter optional), Verwaltungs-API-Absatz nennt alle elf `suspend fun`-Methoden namentlich und schließt „Drei-Sprachen-SDK-Matrix … vollständig"; Versionshistorie 1.81 vorhanden; Version: 1.81; README: Abschnitt „Manage tables and consumers over gRPC" samt Codebeispiel, API-Tabelle und pro-Stream-differenzierter Filteraussage |
| 6 | Closure-Notiz mit Steering-Loop-Lerneintrag | **bestätigt** | Plan §7 ausgefüllt — Steering-Loop-Eintrag: keiner, mit Begründung (Bau-Konfigurationslücke, kein wiederkehrendes Agenten-Verhalten) |
| 7 | Beobachtungs-Register fortgeschrieben — kein neuer Eintrag | **bestätigt** | `BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall` existiert unter `docs/plan/planning/observations/BEO-PGC/`, Zähler unverändert (2×); kein neues Register-Ereignis, konsistent mit den Findings beider Reviews (F-3 ist eine Testlücke, keine Divergenz) |
| 8 | Jedes Risiko aus §6 trägt einen Ausgang | **bestätigt** | beide Einträge tragen einen ausformulierten Ausgang; der Server-Fehler-Ausgang nachgeprüft: Fix-Commit `a40b4809` existiert (`git log --oneline -1`), `internal/bootstrap/assemblersync.go` existiert |
| 9 | Die drei Paarungen — von der Welle-Closure | **korrekt offen** | bewusst offene Welle-Closure-Aufgabe, nicht als Finding gezählt |

Kein `[x]` ohne Beleg. Die offene Zeile 9 ist konsistent Welle-Closure-Arbeit.

## 3. F-1 (HIGH) — Fixrunde nachgemessen, nicht dem Bericht geglaubt

Der Befund lautete: zwei wortgleich aus dem C#-Vorbild übernommene deutsche
Fixture-Strings (`"Rechtsklasse unzureichend für diese RPC"`,
`"interner Fehler"`) brechen die Sprache der Testdatei. Nachgemessen:
`git show bd10c391` zeigt beide Ersetzungen auf Englisch
(`"insufficient role for this rpc"`, `"internal error"`), die aktuelle Datei
lautet durchgängig englisch (grep Exit 1). Die Assertion-Aussage des
Fixrunden-Reviews („der Text wird in keinem Test assertiert") stimmt am
Quelltext: beide `permissionDenied`/`internal`-Tests prüfen ausschließlich
`ex.statusCode`. F-1 ist geschlossen.

## 4. F-2/F-3/F-4 — Bewertung

- **F-2 (LOW, konsolidierte Negativ-Abdeckung):** durch V-1 (§8) in einem Punkt
  unterlaufen — die zwei notFound-Tests der Tabellen-RPCs, die F-2 als
  verifizierbar angab, sind real nur einer (`enableTable`); der
  `disableTable`-notFound-Test läuft nicht. Die LOW-Einordnung der
  Plan-Abweichung selbst bleibt richtig, der Nachweis-Satz trägt seinen Satz
  nur halb.
- **F-3 (LOW, `ConsumerLag.Known = false` / `ListTablesResponse.Retained` mit
  Einträgen ohne eigenen Boundary-Test):** am Quelltext bestätigt
  (DiagnoseTest setzt nur `Known(true)`; TableTest prüft nur
  `retainedCount, 0`) — schmal, nicht blockierend, Stand wie im Review.
- **F-4 (INFO, unbenannte Fremdübernahme in `b239d849`):** vom Fixrunden-Review
  als Klasse bestätigt; die Welle-Closure trägt sie als Prozess-Beobachtung.

## 5. Eigene Stichprobe: Rechtsklassen und Nachrichtenschema

Gegen `proto/cdc/administration/v1/administration.proto`, `administrationRPCRoles`
(`internal/adapters/driving/grpc/interceptor.go`) und die ADR-Tabellen gehalten,
ohne die Review-Negativbefunde zu übernehmen:

| RPC | Rechtsklasse laut ADR | Client-KDoc | Nachrichtenschema | Verdikt |
|---|---|---|---|---|
| `RegisterConsumer` | `roleAdmin` ([ADR-0130](../plan/adr/0130-grpc-verwaltungs-api-neun-rpcs.md) §Teilfrage 4) | „(admin token)“ | `RegisterConsumerRequest{consumer_id, name}` → `…Response{consumer_id, name, already_registered}` | konform |
| `ReadChanges` | `roleReader` ([ADR-0131](../plan/adr/0131-grpc-readchanges-zehnter-rpc.md) §Teilfrage 4, Tabelleneintrag real geprüft) | „(reader or admin token)“ | `ReadChangesRequest{source, schema, table, from, to, limit}` → `ReadChangesResponse{repeated ChangeRecord changes}`, `ChangeRecord` mit dreizehn Feldern in derselben Reihenfolge wie die ADR-Tabelle — am `.proto` gelesen | konform |
| `Diagnose` | `roleReader` ([ADR-0132](../plan/adr/0132-diagnose-ueber-inbound-port-http-grpc.md), Tabelleneintrag real geprüft) | „(reader or admin token)“ | `DiagnoseRequest{source}` → `DiagnoseResponse{heartbeat, capture_lag, consumer_lags, retention_blocker, storage_bytes, backfill}` — die sechs Diagnose-Signale | konform |
| `streamChanges`-Filter | optional, unabhängig, leer = ungefiltert ([ADR-0133](../plan/adr/0133-tabellen-granulare-filterung-grpc-sse.md) §Teilfrage 1/2) | KDoc „both left `null` (the default) delivers every change … unchanged from the original, filter-less contract" | `StreamChangesRequest{schema, table}` — Felder nur gesetzt, wenn nicht-`null` | konform |

Da das SDK die generierten Protobuf-Nachrichten direkt führt (kein DTO-Layer),
ist die Feld-Kongruenz strukturell gegeben; die Stichprobe bestätigt die
Verwendung. Die volle elfköpfige Rollen-Zuordnung (sechs admin, fünf
reader/admin) stimmt mit `administrationRPCRoles` überein — von mir als Block
abgeglichen.

## 6. Plan-vs-Code-Diff

Plan §3 nennt sieben Zeilen; `…/grpc/model/*.kt` entfällt als dokumentierte
Abweichung (§3 „Abweichung vom Plan“ — generierte Protobuf-Nachrichten direkt,
dieselbe Entscheidung wie C#-Geschwister und `PgChangeFeedGrpcClient`). Im Diff
(`git diff --stat 76afcad4..b239d849`, 17 Dateien) sind enthalten:

- `PgChangeFeedAdministrationClient.kt` (neu, 300 Zeilen), `PgChangeFeedGrpcException.kt` (neu, 68), `PgChangeFeedGrpcClient.kt` (Filter-Update, 29 Zeilen geändert) — Plan
- `sdks/kotlin/Dockerfile` (+3) — als Abweichung vom Plan im Plan nachgetragen, sachlich notwendig, analog `sdks/csharp/Dockerfile` (Zeile 51/54)
- fünf Administration-Testklassen + `AdministrationTestClientFactory.kt` + `FakeAdministrationTransport.kt` + `PgChangeFeedGrpcClientFilterTest.kt` — Plan
- `docs/user/benutzerhandbuch.md`, `sdks/kotlin/pgchangefeed-kotlin/README.md` — Plan
- `FakeGrpcStreamTransport.kt` (+9) — **nicht** in der Plan-Tabelle: mechanische Folgenänderung der `GrpcStreamTransport`-Signatur (Request-Parameter für die Filter), V-2
- `docs/user/benutzerhandbuch.md`-Zug (Python-Absatz 1.80) und ein Link-Fix im C#-Verifikations-Report — die im Ausgangs-Review (F-4) transparent berichtete Nebenläufigkeits-Race, kein Kotlin-Inhalt

Keine `Accepted`-ADR berührt, kein Gate gelockert, kein Suppression-Mechanismus
im Diff. Keine unbemerkte Erweiterung: der Client führt exakt elf RPC-Methoden,
zwei Konstruktor-Arten, `close()` und die Transport-Schnittstelle — nichts
darüber hinaus.

## 7. Entscheidungs-Konformität

- **[ADR-0130](../plan/adr/0130-grpc-verwaltungs-api-neun-rpcs.md)/[ADR-0131](../plan/adr/0131-grpc-readchanges-zehnter-rpc.md)/[ADR-0132](../plan/adr/0132-diagnose-ueber-inbound-port-http-grpc.md):** Nachrichtenschema, Rechtsklassen und Fehlerform
  (§5) konform; die Fehlerform-Tabelle (400→`InvalidArgument`,
  401→`Unauthenticated`, 403→`PermissionDenied`, 404→`NotFound`,
  500→`Internal`) ist in `mapException` eins zu eins übernommen, mit
  Fallback-Klasse für übrige Codes — durch `PgChangeFeedAdministrationClientErrorMappingTest`
  (6 laufende Testfälle, 0 failures) belegt.
- **[ADR-0133](../plan/adr/0133-tabellen-granulare-filterung-grpc-sse.md):** additive, optionale, unabhängige Filter-Parameter; kein
  Breaking Change — die parameterlose Signatur bleibt aufrufbar und ist
  regressionsgetestet.
- **AGENTS.md §3.1 (Docker-only):** Build/Test liefen im gepinnten
  `eclipse-temurin:21-jdk`-Image (Digest-Zeile im Bau-Log), kein
  Host-`gradlew`/`kotlinc`. §3.7 (Kommentar-Herkunft): die neuen KDoc-Blöcke
  tragen keine Kennungs-Ketten, kein „ff." — grep über die drei gRPC-Quelltexte
  0 Treffer. §3.9: alle eigenen Gate-Aufrufe einzeln mit direkt gesichertem
  Exit-Code, Log-Filterung nur gegen die geschriebene Log-Datei.
- **Traceability:** alle drei Slice-Commit-Betreffs nennen `LH-FA-SST-009`,
  keiner trägt `SPEC-*`/`ARC-*` im Betreff (`commit-traceability: OK`).

## 8. Findings dieser Verifikation

| # | Kategorie | Befund | Verifizierbar |
|---|---|---|---|
| V-1 | MEDIUM | `PgChangeFeedAdministrationClientTableTest.disableTable table missing at source throws notFound` (Zeilen 76–89) **läuft nie** — still. Der letzte Ausdruck des `runBlocking`-Blocks ist `assertFailsWith<PgChangeFeedGrpcNotFoundException> { … }`, damit trägt die `@Test`-Methode den Rückgabetyp `PgChangeFeedGrpcNotFoundException` statt `void` (`javap` am kompilierten Klassenfile: `public final …PgChangeFeedGrpcNotFoundException disableTable table missing at source throws notFound();`), und die JUnit-Plattform entdeckt Testmethoden mit Nicht-void-Rückgabe nicht — die Testresultate zeigen für die Klasse `tests="6"`, `skipped="0"` (sieben `@Test`-Methoden); ein `--tests`-Filter auf genau diese Methode bricht `:test` mit „no tests found“ ab, dieselbe Probe auf die void-Methode `enableTable table missing at source throws notFound` läuft grün. Wirkung: der `NOT_FOUND`-Negative-Fall von `disableTable` ist ohne laufenden Test; der Review-Satz (F-2, „zwei notFound-Tests … enableTable, disableTable“) trägt diese Hälfte nicht. Kein DoD-Bruch: alle elf RPCs behalten je mindestens einen laufenden Test (Happy-Pfad), und die `NOT_FOUND`→`PgChangeFeedGrpcNotFoundException`-Zuordnung ist generisch über `ErrorMappingTest` und über den `enableTable`-notFound-Test belegt. Behebung: Ergebnis verwerfen (`val ex = assertFailsWith<…> { … }` samt Statuscode-Assertion, Form des `enableTable`-Gegenstücks) — eine Zeile, keine Strukturänderung | ja — `javap`, Testresultate im `--no-cache`-Image, `--tests`-Probe beidseitig |
| V-2 | INFO | `FakeGrpcStreamTransport.kt` (+9 Zeilen: `lastRequest`, Request-Parameter, KDoc-Satz) steht nicht in der Plan-§3-Tabelle — mechanische Folgenänderung der `GrpcStreamTransport`-Signatur, die `streamChanges()`-Filter brauchen den Request im Fake; sachlich notwendig, keine unbenannte Erweiterung im strengen Sinn | ja — `git diff 76afcad4..b239d849 -- …FakeGrpcStreamTransport.kt` |
| V-3 | INFO | Systemprüfung aller Testklassen des Packages auf dieselbe Form (nicht-void `@Test`-Methode): V-1 ist der einzige Fall — alle übrigen nicht-void-Methoden der Testklassen sind Fake-/Factory-Helfer ohne `@Test` | ja — `javap`-Loop über `build/classes/kotlin/test`, Ausgabe ohne zweiten Treffer |

Kein HIGH, keine DoD-Verletzung.

## 9. Verdikt

**DoD erfüllt: ja.** Alle acht `[x]`-Zeilen sind am Ist-Zustand belegt; die
eine `[ ]`-Zeile (drei Paarungen) ist korrekt Welle-Closure-Arbeit. Das
HIGH-Finding des Reviews ist durch die Fixrunde real geschlossen. Die
Rechtsklassen- und Nachrichtenschema-Stichprobe gegen die vier ADRs findet
keinen Widerspruch. `make gates` ist im eigenen Lauf grün (Exit 0), der
Kotlin-Testbestand ist im eigenen `--no-cache`-Lauf real ausgeführt und grün
(19 Suites, 89 Testfälle, 0 failures, 0 errors).

V-1 (MEDIUM) vermindert die beanspruchte negative Abdeckung um einen Fall und
widerruft eine Halbachse des Review-Nachweises F-2 — es bricht keine DoD-Zeile,
weil jede RPC-Fläche laufende Unit-Tests behält und die `NOT_FOUND`-Zuordnung
zweifach anders belegt ist. Empfehlung an den Planner: die Ein-Zeilen-Korrektur
aus V-1 in der Welle-Closure oder dem nächsten geplanten SDK-Zug miterledigen;
bis dahin V-1 als offenen Punkt im Welle-Record führen.

**Übergabe:** an den Planner — die DoD des Slices ist verifiziert, der Slice
kann in die Welle-Closure gehen (drei Paarungen bleiben dort gemäß §2). V-1
wird als MEDIUM in den Steering-Loop-Zähler der Welle übernommen; V-2/V-3 sind
Prozess-Beobachtungen ohne Handlungsbedarf. Dieser Report ist ein **Lauf-Beleg**
(dieser Stand, dieser Lauf) und ersetzt weder Review noch Closure.
