# Verifikations-Report: slice-sdk-kotlin-reale2e — 2026-09-23

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
(`AGENTS.md` §3.12 Instanz B). DoD-Abgleich + Plan-vs-Code-Diff + Draht-/
Grenz-/Gates-Prüfung. Review-Artefakt des Reviewers:
[`review-slice-sdk-kotlin-reale2e.md`](review-slice-sdk-kotlin-reale2e.md);
Hausform dieses Reports:
[`verifikation-slice-sdk-csharp-reale2e.md`](verifikation-slice-sdk-csharp-reale2e.md).

**Gegenstand:** [`../plan/planning/in-progress/slice-sdk-kotlin-reale2e.md`](../plan/planning/in-progress/slice-sdk-kotlin-reale2e.md),
Diff-Range `b58cb173..HEAD` — Substanz `060aa62f` (Implementation: Runner,
Docker-Stufe, Gradle-SourceSet, fünf Testdateien, Make-Target, Träger,
README-Zeile, Plan-Update; 12 Dateien, 825 Insertions / 10 Deletions —
Review-Angabe real nachgemessen) + `c6523009` (Fixrunde F-1, F-2, F-3, F-6
samt Review-Report, 3 Dateien). Committe wurde nichts; der Arbeitsbaum ist
sauber (auch nach allen vier Sensor-Läufen).

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — `AGENTS.md` §3.9)

| Sensor | Ausgang | Beleg aus meinem Lauf |
|---|---|---|
| `make test-sdk-kotlin-integration` | **EXIT=0** | alle vier Phasen grün: gRPC `change_id=805-1`, SSE `810-1`, NATS `817-1` (**laufgebunden** — die DoD-Werte `806-1`/`814-1`/`818-1` sind der Ursprungs-Lauf, meine IDs liegen wenige Positions-Schritte daneben; je ID vom Runner gegen `cdc.changes` mit `source_id`/`table_name`/`new_data->>'name'` = Sentinel gehalten); HTTP `consumer_id=kotlin-sdk-e2e-20260923105937` (Zeitstempel-Form, laufgebunden) gegen `cdc.consumer` gehalten; Runner meldete „Abdeckungs-Traeger unveraendert — docs/user/sdk-e2e-abdeckung.md entspricht dem Quelltext-Stand" — der committete Träger ist byte-identisch dem Runner-Erzeugnis am HEAD; Arbeitsbaum danach sauber |
| `make sdk-pack-kotlin` | **EXIT=0** | Risiko-§6-Ausgang 2 real belegt: die Layer-Kette meines Laufs zeigt `[build 8/11] COPY pgchangefeed-kotlin/src src` **CACHED** (dieser COPY trägt die neue `src/integrationTest`-Quellmenge — inhaltsgekeyte Cache-Kette) und darauf `[build 10/11] RUN ./gradlew --no-daemon test` + `[build 11/11] RUN ./gradlew --no-daemon build` ebenfalls **CACHED**, durch bis `pack-export` — der netzlose Test-/Bau-Pfad lief mit der Integrations-Quellmenge im Bau-Kontext grün und ohne sie; Erzeugnis real: `sdks/kotlin/dist/pgchangefeed-kotlin-0.2.0.jar` |
| `make doc-trace` | exit 0 (advisory) | **79 Anforderungen, 2 Waisen** — `LH-FA-CAP-009`, `LH-FA-CFG-007` (beide `WAISE`-Zeile real in der RTM-Ausgabe); `LH-FA-SST-009` trägt die Spalte `SDK-E2E`, Status `ok` |
| `make gates` | **EXIT=0** | baseline-verify v6.9.0 OK (54 Dateien) · d-check 931 Dateien, 0 Befunde · a-check 0 Befunde · commit-traceability OK (5 Commits, `HEAD~5..HEAD`) · coverage-gate OK — **82.80 %** gegen Schwelle 80 % · generated-sync OK (byte-gleich, Stufe `proto-export`) |
| Stempel | **GLEICH** | `bash tools/harness/working-tree-hash.sh` → `903afc337ff3dea8aac92fad1231bf3dbbcb43fb6f99ef26def22b87156d7ccb` = `.harness/state/gates-passed.diffsha` (Stempel 13:00:51 — von diesem, meinem Gate-Lauf geschrieben) |

Die Zahlen der [`harness/README.md`](../../harness/README.md)-Zeile
`make doc-trace` (79 Anforderungen, 2 Waisen `LH-FA-CAP-009`/`LH-FA-CFG-007`,
`LH-FA-SST-009` über SDK-E2E nicht Waise) stimmen mit meiner Messung
überein (`AGENTS.md` §3.12 Instanz A nachgemessen).

## 2. DoD — Verdikt je Zeile (§2 des Plans)

| # | DoD-Zeile | Verdikt | Realer Beleg |
|---|---|---|---|
| 1 | Vier Kotlin-Flächen tragen je einen realen Rundlauf, je SQL-Gegenprüfung (`cdc.changes` gRPC/SSE/NATS; `cdc.consumer` HTTP), je Ablehnungs-Beleg (`Unauthenticated` / `401` / NATS-Verbindungsablehnung) | **erfüllt** | Eigener Lauf EXIT=0 (siehe §1): drei Stream-Phasen je über `RECEIVED change_id=…`-Marker mit `received_grep`-Form (`table=`/`operation=INSERT`/`new_image=`+Sentinel) und SQL-Variante `changes` (`count(*) >= 1`) gehalten; HTTP-Phase über SQL-Variante `consumer`; REJECTED-Marker je Phase vom Runner erzwungen (`grep -qF "$reject_marker"`: gRPC `REJECTED code=Unauthenticated`, SSE/HTTP `REJECTED status=401`, NATS `REJECTED token-rejected` — fehlt der Marker, bricht der Lauf rot, Runner-Z. 245–248) |
| 2 | Mechanik: additive Stufe `integration` auf `FROM build`, dieselben Pins, kein neuer Pin; Runner mit `set -euo pipefail`, Cleanup-Trap, explizite Phasen-Auswahl; Make-Target in `harness/mk/sdk.mk` | **erfüllt** | Dockerfile-Diff: einzige `FROM`-Neuheit ist `FROM build AS integration` (Z. 148) — alle `FROM`-Zeilen Parent↔HEAD gegengeprüft, der Pin `eclipse-temurin:21-jdk@sha256:085e…` (Z. 43) unverändert, kein neuer Pin; `:?`-Guard im Stufen-CMD real (`${PGCHANGEFEED_TEST_NAME:?…}`); Runner: Trap + Cleanup + vier explizite `run_phase`-Aufrufe (Phase-Auswahl je Testklasse — `BEO-PGC/test-runner-stiller-ausschluss`-Disziplin); Target + `.PHONY` + `##`-Hilfe vorhanden |
| 3 | Träger-Erweiterung im selben Zug: Kotlin-Abschnitt, idempotent aus derselben Messung | **erfüllt** | Mein Lauf: „Abdeckungs-Traeger unveraendert" — der committete Abschnitt (4 Zeilen + 2 Marker) ist byte-identisch dem Runner-Erzeugnis am HEAD; C#-Abschnitt byte-identisch erhalten (Diff: nur Anhängung hinter `csharp-end`); Writer-Erhalt beidseitig im Quelltext (`awk` vor `kotlin-begin` / ab `kotlin-end`) |
| 4 | Kein Import aus `internal/**`/`cmd/**`/`gen/**` in `sdks/kotlin/**` | **erfüllt** | Grep über `sdks/kotlin/`: die `cdc.stream.v1.*`-Imports liegen ausschließlich in bestehendem `src/main`/`src/test` (SDK-eigene, im Docker-Bau aus der `.proto` erzeugte Stubs — kein Bezug auf `gen/**` dieses Repos); die fünf neuen integrationTest-Dateien importieren nur `io.github.pt9912.pgchangefeed.*`, `io.grpc.Status`/`StatusException`, `java.*`, `kotlin.test.*`, `kotlinx.coroutines.*`; kein `examples/`-Bezug; keine `@Suppress`-Form (§3.2) |
| 5 | `make gates` grün | **erfüllt** | eigener Lauf EXIT=0 + Stempel-Gleichheit (siehe §1) |
| 6 | Review durchgeführt, Report unter `docs/reviews/` | **korrekt offen** | Review-Report liegt vor (in `c6523009` committet); die Fixrunde ist derselbe Commit — der Nachzug setzt das Häkchen bei Schritt 21 (`BEO-PGC/dod-checkbox-nachzug`-Form wie im C#-Vorgänger), bewusst nicht vorschnell |
| 7 | Doku-Update: `harness/README.md`-Zeile nach dem realen Lauf | **erfüllt** | Zeile existiert (Parent: 0 Treffer, HEAD: 1); ihre Verhaltens-Aussagen trägt mein Lauf real (Bring-up im Python-Runner-Muster, vier Phasen, Reject-Formen, SQL-Gegenprüfungen, Träger-Schreibung, `:?`-Guard, `:dev`-Voraussetzung, „Gradle-Task hängt bewusst nicht an `check` — `make sdk-pack-kotlin` bleibt netzlos grün, real gemessen" — mein Pack-Lauf) |
| 8–11 | Closure-Notiz, Reconciliation-Register (selbst „entfällt"), Beobachtungs-Register, §6-Risiko-Ausgänge, drei Paarungen | **offen — Closure-zeitlich** | Slice steht in `in-progress/`; diese Zeilen sind die Closure-Pflicht, kein DoD-Verstoß |

Die DoD-Werte sind nicht driftend: die quotierten `change_id`s sind
laufgebunden — die DoD-Zeile trägt ihren Ursprungs-Lauf namentlich, mein
Lauf trägt seine eigenen Werte samt wirksamer SQL-Gegenprüfung
(§3.12-konform); nur der Form nach abweichend, nicht in der Aussage.

## 3. Plan-vs-Code-Diff (§3 + Plan-Nachzug + §3.13-Feld)

Alle 6 §3-Zeilen und alle 3 Plan-Nachzug-Zeilen sind im Diff vertreten, jede
Zeile einer realen Änderung zugeordnet:

- Runner-Skript (378 Zeilen am HEAD), fünf Kotlin-Testdateien
  (`PhaseEnvironment.kt` + vier Testklassen à Happy-Path + Reject),
  Dockerfile-Stufe `integration`, `harness/mk/sdk.mk`-Target, Träger
  (6 Zeilen), README-Zeile — alles real.
- Plan-Nachzug: SourceSet `integrationTest` + Task **nicht** an `check`
  (Grep: keine `check.dependsOn`-Form), zwei eigene Konfigurationen mit
  Kommentar-Begründungszug (`integrationTestImplementation` erbt
  `testImplementation`; `integrationTestRuntimeClasspath` resolvable, erbt
  `testRuntimeOnly` **und** `runtimeOnly` — `grpc-netty-shaded`,
  Kommentar trägt die ProviderNotFoundException-Begründung), Task-Classpath
  explizit aus diesen Konfigurationen, `showStandardStreams` live
  (Risiko-§6-Ausgang 3); gson-Form `oldImage == null || oldImage!!.isJsonNull`
  in SSE/NATS-Test real (SSE `SseRealserverTest.kt:42`, NATS
  `NatsRealserverTest.kt:42`). Keine still gestrichene Plan-Zeile.
- Commit-Split sauber: `060aa62f` = Code/Träger/Verdrahtung/Plan-Update,
  `c6523009` = Fixrunde F-1/F-2/F-3/F-6 + Review-Report.

**§3.13-Suchlauf-Feld — gegen beide Stände (`b58cb173` und HEAD) nachgemessen,
alle 6 Zeilen bestätigt:**

1. `harness/README.md` neue Target-Zeile: Parent 0 Treffer → HEAD 1 ✓
2. `docs/user/sdk-e2e-abdeckung.md` Kotlin-Abschnitt: Parent ohne Marker →
   HEAD 6 Zeilen (2 Marker + 4 Zeilen) ✓; C#-Abschnitt byte-identisch
   erhalten ✓
3. `make test-sdk-csharp-integration`-Zeile: Parent und HEAD byte-identisch
   (nur Zeilenverschiebung 160→161 durch die Einfügung) — „gemeldet, nicht
   gezogen" ist die wahre Behandlung, Fixrunde F-6 korrekt ✓
4. `sdks/kotlin/pgchangefeed-kotlin/README.md` §Status: existiert real
   (`sdks/kotlin/` trägt kein eigenes `README.md` — F-1-Korrektur
   verifiziert), trägt keine Teststrategie-Aussage über Realserver-Läufe
   (0 Treffer auf Realserver/E2E) ✓
5. `docs/user/benutzerhandbuch.md`: Kotlin-Nennungen nur als
   Beispiele-Verweise, 0 Realserver-/Integrationstest-Treffer, im Range
   unverändert ✓
6. `spec/pflichtenheft.md`: im Range unverändert (`SPEC-027`s Deckungs-Aussage
   bleibt richtig) ✓

Breiterer eigener Suchlauf über beide Stände (`git grep` Kotlin × Realserver
über `docs/ harness/ spec/`): dieselbe Trägermenge an beiden Ständen; die
`done/`-Records (`slice-sdk-kotlin-nats-stream-client-flaeche.md:279`,
`slice-sdk-kotlin-sse-client-flaeche.md:53`, `welle-sdk-kotlin-vollabdeckung.md:129`)
tragen ihre „kein eigener Realserver-Test"-Aussagen sliced-/welle-scoped
(„diese Fläche", „diese Welle") und bleiben als Records wahr; Roadmap/Welle
tragen Planungs-Aussagen, keine falsch werdenden Träger. Keine weitere
Fundstelle, die durch diesen Zug falsch geworden wäre.

## 4. Draht-Konformität (Task-Punkt 3)

- **gRPC:** Client setzt `authorization`-Metadata mit `Bearer <token>`
  (`PgChangeFeedGrpcClient.kt:115,127-130`); Reject-Test fordert
  `StatusException` mit `Status.Code.UNAUTHENTICATED`
  (`GrpcRealserverTest.kt:71-74`), nicht eine leere Stream-Annahme ✓
- **SSE:** `Authorization: Bearer` via `options.apiToken`
  (`PgChangeFeedSseClient.kt:98`); Reject als getippte
  `PgChangeFeedUnauthorizedException` mit `statusCode == 401`
  (`SseRealserverTest.kt:67-71`) ✓
- **NATS:** Auth auf Verbindungsebene (Token aus `PgChangeFeedClientOptions`,
  kein per-call-Auth — Klassen-Doku `PgChangeFeedNatsStreamClient.kt:22-29`);
  Subjekt über `buildSourceSubject(sourceId)`; Reject gebunden an den
  Server-Rohwortlaut `Authorization` über die Ausnahme-Kette
  (`exception.toString().contains("Authorization")`, mit Kommentar zur
  io.nats-Wrapping-Form) ✓
- **HTTP:** `Authorization: Bearer` (`PgChangeFeedHttpClient.kt:143`);
  admin-/reader-Token-Paar je Phase-Env gebunden; Reject `401` ✓
- **SQL-Gegenprüfungen im Runner:** `changes`-Variante prüft `cdc.changes`
  gegen `source_id`/`change_id`/`table_name`/`new_data->>'name'` = Sentinel;
  `consumer`-Variante prüft `cdc.consumer` gegen `consumer_id` — je Phase
  über die explizite Variante verdrahtet (Runner-Z. 259–265) ✓
- **Ablehnungs-Marker je Phase** im Runner erzwungen (siehe DoD-Zeile 1);
  Phase-Adressen gegen den `compose.yaml`-Container-Vertrag gehalten
  (`pg-change-feed:9090`/`:8090`, `nats://nats:4222`, Netz `cdc-feed-test`,
  Token `e2e-*`, `pub_pgc_e2e`, `slot_pgc_e2e`, `src-e2e`) ✓

## 5. Gradle-Form (Task-Punkt 4) — eigener Beleg

- `integrationTest`-Task hängt **nicht** an `check`: keine
  `check.dependsOn`-Form in `build.gradle.kts` (Grep); die Pack-Closure
  (`--target pack-export`: `build`→`pack`→`pack-export`) enthält die Stufe
  `integration` (`FROM build AS integration`) strukturell nicht — und die
  Layer-Kette meines `make sdk-pack-kotlin`-Laufs (§1) belegt das
  Verhalten real: der `RUN ./gradlew test`-Layer ist über dem COPY der
  aktuellen Quellmenge (inkl. `src/integrationTest`) CACHED und grün.
- Zwei Konfigurationen mit Kommentar-Begründungszug (Kopplung/Grenze-Form,
  `grpc-netty-shaded`-Begründung) ✓; Pins unverändert ✓.

## 6. Abweichungen und Beobachtungen (keine DoD-Verletzung)

1. **Lauf-Belege ohne committetes Artefakt** (INFO, `AGENTS.md` §3.12
   Klasse): die DoD-Mutation („gültiger Token im gRPC-Reject-Test — der
   Lauf färbte rot („der Ablehnungs-Beleg blieb aus“), Revert,
   Abschlusslauf grün") und die gson-Fassung („real im vierten Lauf als
   rote Assertion gesehen") tragen keinen auflösbaren Beleg-Anker (kein
   Log, keine Evidence-Datei unter `docs/plan/planning/evidence/`). Die
   Falsifikation, die ich erreichen konnte, ist die Quelltext-Ebene: der
   Rot-Pfad existiert real und trägt wortgleich die quotierte Zeile
   (Runner-Z. 246: „der Ablehnungs-Beleg blieb aus ($reject_marker
   fehlt)"); die gson-Assertion ist rot-fähig bei real vorhandenem
   Alt-Bild. Die Behauptung bleibt damit konsistent, aber
   nachfahrbar-nicht-archiviert — dieselbe Klasse, die der C#-Review F-8
   benannt hat.
2. **ID-Werte Parent↔mein Lauf** (erwartet, keine Abweichung): die C#-Nachmessung
   reproduzierte ihre `change_id`s byte-identisch; hier driften sie um wenige
   Positions-Schritte (`806-1`→`805-1`, `814-1`→`810-1`, `818-1`→`817-1`).
   Beide Werte sind laufgebunden und je gegen `cdc.changes` validiert — die
   DoD-Aussage („je SQL-Gegenprüfung") ist in meinem Lauf wirksam belegt;
   nur die Byte-Identität der IDs ist keine stabile Eigenschaft dieser
   Bring-up-Kette.
3. **d-check-Dateizahl im Review-Report** (INFO): der Review-Report trägt
   „930 Dateien, 0 Befunde" — meine Messung am HEAD zeigt 931; die Differenz
   ist der Review-Report selbst, der seit der Messung committet ist
   (selbst-referentiell, kein Drift).

Sonst **0 DoD-Abweichungen**: keine behauptete Zahl driftet gegen meine
Messung, keine still gestrichene Plan-Zeile, keine Grenzverletzung, keine
Pin-Änderung, kein Gate rot, Arbeitsbaum nach allen Läufen sauber.

## Gesamt-Urteil

**DoD-Verdikt: erfüllt** für alle fünf realprüfbaren `[x]`-Zeilen gegen
eigene Sensoren-Läufe (Pflichtbeleg EXIT=0 mit allen vier Phasen,
`sdk-pack-kotlin` EXIT=0 mit Layer-Ketten-Beleg, doc-trace 79/2, gates
EXIT=0 + Stempel gleich); die offenen `[ ]`-Zeilen sind Closure-zeitlich
(Review-Nachzug nach Fixrunde, Closure-Notiz samt §6-Ausgängen und
Lerneintrag, Register, Paarungen) und korrekt offen, solange der Slice in
`in-progress/` liegt. Der Fixrunde-Commit zieht F-1/F-2/F-3/F-6 real nach
(F-1-Träger und F-6-Behandlung am HEAD verifiziert); F-4/F-5 bleiben INFO
ohne erwartete Aktion (Review-Verdikt). Die RTM-Zahlen der README (79/2,
`LH-FA-SST-009` via SDK-E2E) sind durch meine eigene Messung gedeckt.
Verbleibende Reste: eine nicht-blockierende Beobachtung zur Belegbarkeit der
Mutations-/gson-Rot-Läufe (§6.1). Der Slice ist bereit für den
Review-Nachzug und die Closure-Notiz; `done/`-Übergang erst nach deren
Tragen (§7 des Plans).