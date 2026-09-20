# Slice sdk-kotlin-grpc-client-flaeche: Öffentliche gRPC-Stream-Client-Fläche (`SPEC-020`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** [welle-sdk-kotlin-lh-fa-sst-009](../welle-sdk-kotlin-lh-fa-sst-009.md).

**Bezug:** [`LH-FA-SST-009`](../../../../spec/lastenheft.md),
[`LH-FA-SST-008`](../../../../spec/lastenheft.md) (Live-Change-Stream —
Boundary: kein Replay im Stream selbst),
[`ADR-0109`](../../adr/0109-kotlin-github-packages-drittes-sdk-package.md)
Festlegung 1/3 (Umfang, `.proto`-Bezug über Zusatzkontext),
[`ADR-0060`](../../adr/0060-grpc-streaming-mechanismus.md)
(gRPC-Server-Streaming-Vertrag, wird vom SDK benutzt, nicht erweitert).

**Berührte Spec-Stellen:** [`SPEC-020`](../../../../spec/pflichtenheft.md)
(Nachrichtenschema, RPC-Name, Stream-Semantik — das SDK benutzt diese
Festlegungen, verändert sie nicht).

**Verantwortlich:** Implementer-Agent, 2026-09-20.

**Autor:** Planner-Agent, direkt beauftragt (`ADR-0109` §Konsequenzen
Folgepflicht 1). **Datum:** 2026-09-20.

---

## 1. Ziel und Abgrenzung

**Ziel:** Eine öffentliche, stabile Kotlin-API-Fläche für den
`StreamChanges`-RPC von [`SPEC-020`](../../../../spec/pflichtenheft.md) — ein
Server-Streaming-Aufruf, der ein `kotlinx.coroutines.flow.Flow<Change>`
(oder gleichwertig idiomatisch, analog dem Coroutine-Stub-Muster in
`examples/kotlin/grpc-client`) an den Consumer liefert, Authentifizierung
über den `authorization`-Metadata-Eintrag (`Bearer <token>`), dieselbe
`.proto`-Bezugsform wie `examples/kotlin/grpc-client` (zusätzlicher,
benannter Docker-Bau-Kontext `--build-context proto=proto`, kein
committeter Stub — real bestätigt: `examples/kotlin/grpc-client` nutzt
denselben Mechanismus wie `examples/csharp/grpc-client`, siehe
`harness/mk/examples.mk`).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **HTTP-API** — `slice-sdk-kotlin-http-client-flaeche` übernimmt das;
  getrennter Draht-Vertrag, keine Überschneidung.
- **Tabellen-granulare Filterung des Streams** — `SPEC-020` selbst trägt
  sie nicht („eine tabellen-granulare Filterung ist nicht Teil dieser
  Version"); das SDK kann keine Fähigkeit anbieten, die der Draht nicht
  hat.
- **Stream-internes Replay** — `SPEC-020`/[`LH-FA-SST-008`](../../../../spec/lastenheft.md)
  Boundary schließt das ausdrücklich aus; ein Consumer, der Replay
  braucht, nutzt den bestehenden Lesezugriffsweg (`GET /changes`,
  `slice-sdk-kotlin-http-client-flaeche`), nicht dieses Stream-API.
- **SSE- oder NATS-Vollinhalts-Stream** — `ADR-0109` Festlegung 1 grenzt
  v1 ausdrücklich auf HTTP-API und gRPC-Stream ein (Welle-Plan §6
  Out-of-Scope), obwohl für Kotlin bereits funktionierende Beispiele für
  beide existieren.
- **CLI-Argument-Parsing wie `examples/kotlin/grpc-client`** — dasselbe
  Argument wie bei der HTTP-Fläche: ein SDK-Consumer ruft eine Methode
  auf, parst keine `argv`.
- **Committeter Protobuf-/Coroutine-Stub im SDK-Baum** — `ADR-0109`
  Festlegung 3 verlangt denselben Bezugsweg wie `examples/kotlin/grpc-client`:
  die `.proto` bleibt die einzige Quelle, der Stub entsteht im Bau über
  das Gradle-Plugin `com.google.protobuf` (`protoc-gen-grpc-java` +
  `protoc-gen-grpc-kotlin`).

## 2. Definition of Done

- [x] `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/…/grpc/` (oder
      gleichwertiger Namensraum) trägt eine öffentliche Client-Klasse mit
      einer Methode, die den `StreamChanges`-RPC öffnet und die Nachrichten
      von [`SPEC-020`](../../../../spec/pflichtenheft.md) (`change_id`,
      `transaction_id`, `source_table_id`, `sequence`, `operation`,
      `old_image`, `new_image`, `schema_version`, `schema`, `table`) an den
      Consumer weiterreicht — Bearer-Token wird bei Konstruktion oder
      Aufruf übergeben, landet im `authorization`-Metadata-Eintrag.
      Umgesetzt als `io.github.pt9912.pgchangefeed.grpc.PgChangeFeedGrpcClient`.
- [x] `sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts` bekommt dieselben
      gRPC-/Coroutine-Koordinaten wie `examples/kotlin/grpc-client/build.gradle.kts`
      (`io.grpc:grpc-kotlin-stub`, `io.grpc:grpc-netty-shaded`,
      `io.grpc:grpc-bom`, `com.google.protobuf:protobuf-java`,
      `org.jetbrains.kotlinx:kotlinx-coroutines-core`, das
      `com.google.protobuf`-Gradle-Plugin) — Versionen zum Bau-Zeitpunkt
      dieses Slice real neu gemessen, nicht aus `examples/kotlin/` oder
      `ADR-0109`s Kontext-Messung (2026-09-17) unbesehen übernommen
      (`AGENTS.md` §3.12). Reale Drift gefunden und übernommen:
      `com.google.protobuf:protoc`/`protobuf-java` `4.36.1` → `4.36.2`
      (Maven Central, 2026-09-20); alle übrigen Koordinaten unverändert
      (siehe `build.gradle.kts`-Kommentar für den vollständigen Messbeleg
      inkl. der als Nicht-Release erkannten Commit-Hash-Metadaten-Anomalie
      bei `grpc-kotlin-stub`/`protoc-gen-grpc-kotlin`).
- [x] `sdks/kotlin/Dockerfile` bekommt den zusätzlichen, benannten
      Bau-Kontext `proto` (`--build-context proto=proto`, `COPY --from=proto
      cdc/stream/v1/changestream.proto …`) — ohne ihn bricht der Bau an der
      `COPY`-Zeile ab, kein stiller Fallback (Muster
      `examples/kotlin/Dockerfile`/`harness/mk/examples.mk`, hier auf den
      SDK-Baum übertragen). Real geprüft: `docker build --build-context
      proto=proto -f sdks/kotlin/Dockerfile sdks/kotlin` (ohne den
      Zusatzkontext läuft `generateProto` als `NO-SOURCE`, mit ihm real
      erzeugt — beide Fälle real ausgeführt).
- [x] Eigene Tests decken mindestens: Nachrichtenschema-Vollständigkeit
      (Feld-für-Feld gegen `SPEC-020`, analog
      `internal/adapters/driving/grpc/server_test.go`s Feldvollständigkeits-
      Test) und den Authn-Boundary-Pfad (fehlendes/ungültiges Token →
      `Unauthenticated`, ohne echten Server — ein Fake/Stub des
      generierten Coroutine-Stubs oder gleichwertig, netzlos).
      `PgChangeFeedGrpcClientMessageSchemaTest`/`PgChangeFeedGrpcClientAuthBoundaryTest`
      gegen `FakeGrpcStreamTransport`; beide real mutationsgetestet (rot
      färbende Mutation je Zusage, im Implementer-Bericht dokumentiert).
- [x] Kein Import aus `internal/**`/`cmd/**`/`gen/**` dieses Repos —
      der generierte Stub entsteht im SDK-eigenen Bau aus der `.proto`,
      nicht als Kopie von `gen/**` (`ADR-0109` §Kontext Bindung „Import-
      Grenze, hier ohne Ausnahme") — real geprüft:
      `grep -rn "internal/\|cmd/\|gen/" sdks/kotlin/` liefert höchstens
      einen Treffer als Doku-Kommentar-Zitat des Test-Vorbilds
      (`internal/adapters/driving/grpc/server_test.go`), keinen Import.
      Real gemessen: genau ein solcher Treffer (in
      `PgChangeFeedGrpcClientMessageSchemaTest.kt`s KDoc) plus ein
      unveränderter, vor diesem Slice bereits vorhandener Treffer im
      Gradle-Wrapper-Skript `gradlew` (Upstream-Kommentarzeile
      `org/gradle/api/internal/plugins/…`, kein Repo-Pfad, nicht Teil
      dieses Diffs).
- [x] `make gates` grün. Exit-Code `0`, direkt geprüft (`AGENTS.md` §3.9) —
      Beleg im Implementer-Bericht.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update: `docs/user/benutzerhandbuch.md` bekommt einen
      SDK-Hinweis für die Kotlin-gRPC-Oberfläche (`ADR-0109` §Konsequenzen
      Folgepflicht 4, **inklusive** des PAT-Hinweises für den Bezug über
      GitHub Packages) — getragen durch die bereits verkörperte
      Selbstprüf-Instruktion und den Reviewer-HIGH-Punkt
      (`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`).
      Zusätzlich die veraltete Aussage „gRPC-Change-Stream folgt in einem
      Folge-Release" im bestehenden HTTP-`**SDK:**`-Absatz korrigiert
      (Träger-Nachzug, `AGENTS.md` §3.13) — Version `1.36` → `1.37`.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben,
      **falls dieser Slice einen Inventur-Fund auflöst** — entfällt: keine
      Reconciliation-Datei in diesem Repo.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder Beleg in `evidence/`; keine Beobachtung angefallen
      ist ebenfalls eine Antwort und wird in §7 notiert. Keine neue
      Beobachtung — siehe §7.
- [x] Jedes Risiko aus §6 trägt einen Ausgang.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      dieser Slice gehört zu
      [welle-sdk-kotlin-lh-fa-sst-009](../welle-sdk-kotlin-lh-fa-sst-009.md)
      (noch offen); die Prüfung läuft regelkonform bei deren Closure.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/…/grpc/PgChangeFeedGrpcClient.kt` (Arbeitsname) | neu | öffentliche API-Fläche für `StreamChanges` (`SPEC-020`). |
| `sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts` | update | gRPC-/Coroutine-/Protobuf-Gradle-Plugin-Koordinaten, real neu gemessen. |
| `sdks/kotlin/Dockerfile` | update | zusätzlicher `proto`-Bau-Kontext, `COPY --from=proto …`-Zeile. |
| `sdks/kotlin/pgchangefeed-kotlin/src/test/kotlin/…/grpc/*Test.kt` | neu | Nachrichtenschema-Vollständigkeit, Authn-Boundary. |
| `docs/user/benutzerhandbuch.md` | update | SDK-Hinweis für die Kotlin-gRPC-Oberfläche inkl. PAT-Hinweis, im selben Zug (`ADR-0109` Folgepflicht 4). |

**Ansatz:** Referenzmaterial ist `examples/kotlin/grpc-client/src/main/kotlin/cdcexamples/grpc/Main.kt`
(Draht-Kenntnis für Kanal-Aufbau und Metadata-Header, kein
`project(":…")`-Abhängigkeitspfad — `ADR-0109` Festlegung 3: „dieselbe
Draht-Kenntnis, aber als eigenständiger, paketierbarer Code neu
geschrieben") und `internal/adapters/driving/grpc/server_test.go`
(Feldvollständigkeits-Testmuster, serverseitig, als Vorbild für die
Consumer-seitige Prüfung).

**Hinweise aus dem Beobachtungs-Register (proaktiv):**

- `examples/kotlin/grpc-client/build.gradle.kts` trägt ausführliche,
  real recherchierte Versions-Kommentare (BOM-Angleichungen zwischen
  `io.grpc:grpc-bom`, `protoc`, `protobuf-java`, `kotlinx-coroutines-core`)
  — dieselbe Sorgfalt gilt für dieses Slice: keine Version blind aus dem
  Beispiel kopieren, ohne die dort dokumentierten Versions-Kopplungen
  (BOM-Zeile, `protoc`-Major-Gleichlauf mit `protobuf-java`) real
  nachzuvollziehen (`BEO-PGC/workaround-uebersieht-etabliertes-muster-im-bestand`,
  Zähler 2×, siehe Welle-Plan §6 — hier in der positiven Form: das
  etablierte Muster **wird** übernommen, nicht übersehen).
- `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (19×): Träger-Nachzug-
  Suchlauf vor Abschluss über `README.md`, `build.gradle.kts` — beide
  Sprachachsen (deutsch/englisch).
- **Backtick-Paritäts-Check** vor jedem Commit dieses Slice.

## 4. Trigger

**Start** (`next` → `in-progress`): wenn `slice-sdk-kotlin-projektgeruest`
in `done/` liegt (siehe Welle-Plan §4 Reihenfolge) — parallelisierbar zu
`slice-sdk-kotlin-http-client-flaeche`, keine gegenseitige Abhängigkeit.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): entfällt aus
  heutiger Sicht — ein RPC, ein Nachrichtentyp.
- `in-progress` → `open` (blockiert — Carveout?): `slice-sdk-kotlin-projektgeruest`
  liegt noch nicht in `done/`, oder der Docker-Bau-Kontext-Mechanismus
  (`--build-context proto=proto`) lässt sich aus einem noch unbekannten
  Grund nicht auf den SDK-Baum übertragen (unwahrscheinlich —
  `examples/kotlin/grpc-client` belegt den Mechanismus bereits real, siehe
  `harness/mk/examples.mk`).

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + Closure-Notiz mit Lerneintrag
geschrieben.

## 6. Risiken und offene Punkte

- Der Bau dieses Slice ist ohne ein aufrufendes `make`-Ziel (das erst
  `slice-sdk-kotlin-pack-werkzeug` liefert) nur über einen direkten
  `docker build --build-context proto=proto …`-Aufruf prüfbar — kein
  Komfort-Ziel für diesen Zwischenstand. **Ausgang:** eingetreten,
  akzeptiert: der direkte Aufruf reicht als Beleg für die DoD dieses
  Slice; das komfortable `make`-Ziel folgt bewusst erst mit dem
  Pack-Werkzeug-Slice (Welle-Plan §4 Reihenfolge), dieselbe Struktur wie
  bei der C#-Welle.
- Die Kotlin-/Gradle-gRPC-Werkzeugkette braucht real mehr einzeln gepinnte
  Koordinaten als die C#-Variante (`Grpc.Tools` bündelt bei .NET intern,
  was bei Kotlin/Gradle als eigenständige Artefakte — `protoc`,
  `protoc-gen-grpc-java`, `protoc-gen-grpc-kotlin`, das
  `com.google.protobuf`-Gradle-Plugin — einzeln gepinnt werden muss,
  `ADR-0109` §Kontext). Eine unvollständige oder inkonsistente
  BOM-Angleichung (siehe `examples/kotlin/grpc-client/build.gradle.kts`s
  Kommentare zu `protobuf-java`-Versionsdrift) könnte den Bau mit
  `cannot find symbol` scheitern lassen. **Ausgang:** eingetreten und
  aufgelöst — real gebaut (`docker build --build-context proto=proto -f
  sdks/kotlin/Dockerfile sdks/kotlin`), `compileKotlin`/`compileJava` beide
  ohne `cannot find symbol`; die real neu gemessene BOM-Angleichung
  (`protoc`/`protobuf-java` explizit auf `4.36.2` angehoben, `grpc-protobuf`
  zieht weiterhin transitiv `protobuf-java:3.25.9`) trägt dieselbe
  funktionierende Kombination wie `examples/kotlin/grpc-client`. Ein
  unabhängiger, realer Kopierfehler (Docker-`COPY`-Zielpfad relativ zum zu
  diesem Zeitpunkt bereits gewechselten `WORKDIR /src/pgchangefeed-kotlin`,
  nicht zur Bau-Kontext-Wurzel `/src` — ein doppelter Pfadanteil
  `pgchangefeed-kotlin/pgchangefeed-kotlin/…`) wurde vom ersten realen
  Docker-Bau-Lauf sichtbar gemacht (`generateProto NO-SOURCE`) und noch vor
  dem ersten Commit dieses Slice korrigiert — kein Nacharbeits-Slice nötig.
- Ein Fake/Stub für die Authn-Boundary könnte den realen
  gRPC-`Unauthenticated`-Status-Pfad nicht exakt nachbilden. **Ausgang:**
  weiter offen — ein realer Rundlauf-Beleg bleibt
  `make test-integration`s bestehendem `tools/harness/grpcclient`
  vorbehalten (Wegwerf-Client, kein SDK-Import); dieses SDK bekommt
  frühestens mit einem Folge-Slice einen eigenen Integrationsbeleg.

## 7. Closure-Notiz

- **Was hat funktioniert:** Die reale Neu-Messung (`AGENTS.md` §3.12) fand
  genau eine echte Drift ggü. `examples/kotlin/grpc-client`s
  2026-09-17-Messung (`protoc`/`protobuf-java` `4.36.1` → `4.36.2`,
  drei Tage später) und erkannte korrekt eine zweite, irreführende
  Metadaten-Anomalie (Commit-Hash-„Releases" für `grpc-kotlin-stub`/
  `protoc-gen-grpc-kotlin`) als Nicht-Release statt sie blind zu
  übernehmen — beide Ergebnisse real gegen `repo1.maven.org`
  nachvollzogen, nicht behauptet. Der Fake-Transport-Seam
  (`GrpcStreamTransport`/`FakeGrpcStreamTransport`) hielt das exakte
  Analogie-Muster von `HttpTransport`/`FakeHttpTransport` ein — inklusive
  der aus dem vorigen Slice gelernten, korrekten `internal`-Kommentierung
  (compile-time Kotlin-Grenze, keine JVM-Bytecode-Schranke) von Anfang an,
  ohne die HIGH-Fixrunde des vorigen Slice zu wiederholen. Beide
  Kern-Zusagen (Authn-Boundary, Nachrichtenschema-Vollständigkeit) wurden
  real mutationsgetestet: `.catch { }` (verschluckt den
  `StatusException`) und `.map { it.toBuilder().clearSchema().build() }`
  (entfernt ein Feld) machten je einen Testlauf real rot, vor der Rückkehr
  zur sauberen Fassung.
- **Was ging anders als geplant:** Der erste reale
  `docker build --build-context proto=proto …`-Lauf schlug fehl
  (`generateProto NO-SOURCE`, dann `Unresolved reference 'cdc'`) — die
  `COPY --from=proto …`-Zielpfadangabe im Dockerfile duplizierte das
  Segment `pgchangefeed-kotlin/`, weil `WORKDIR` an dieser Stelle bereits
  auf `/src/pgchangefeed-kotlin` gewechselt war (Docker-`COPY`-Ziele sind
  relativ zum aktuellen `WORKDIR`, nicht zur Bau-Kontext-Wurzel) — ein
  Unterschied zu `sdks/csharp/Dockerfile`, das `WORKDIR /src` nie
  wechselt. Der Fehler wurde durch den realen Bau-Lauf selbst sichtbar
  (kein stiller Fallback) und vor dem ersten Commit korrigiert (§6, Risiko
  2 „Ausgang").
- **Steering-Loop-Eintrag:** Ein Docker-`Dockerfile`, das `WORKDIR`
  zwischen dem Kopieren des Wrapper-Baums und dem eigentlichen
  Quell-Kopierschritt wechselt (Muster: `sdks/kotlin/Dockerfile`,
  `examples/kotlin/Dockerfile` bleibt bei `WORKDIR /src`), braucht bei
  jedem `COPY --from=<zusatzkontext>`-Ziel eine bewusste Prüfung, ob der
  Zielpfad relativ zum *aktuellen* `WORKDIR` oder — fälschlich in Analogie
  zum Bau-Kontext-Pfad des `--from`-Quellsegments — zur Bau-Kontext-Wurzel
  gemeint war; der reale `docker build`-Lauf ist der einzige Sensor, der
  diesen Unterschied zuverlässig zeigt (ein `generateProto NO-SOURCE` bei
  vorhandener `.proto`-Quelle ist der Leitbefund).
- **Beobachtungs-Register (`../observations/`):** Keine neue Beobachtung.
  Der Copy-Pfad-Fehler oben ist kein wiederkehrendes, bereits im Register
  geführtes Muster (geprüft: kein Treffer für „doppelter Pfadanteil"/
  „COPY --from=proto" außerhalb des bereits bekannten, andersartigen
  `zusatzkontext-kopplung-breiter-als-dod-wortlaut`-Eintrags) und trat nur
  im eigenen Entwurf auf, nie committet — kein Repo-weites Muster, das
  einen eigenen `BEO-PGC`-Eintrag rechtfertigt; der Lerneintrag oben hält
  ihn stattdessen fest. `BEO-PGC/workaround-uebersieht-etabliertes-muster-im-bestand`
  (2×) und `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (19×) wurden
  korrekt vermieden (siehe „Was hat funktioniert" und Träger-Nachzug
  unten) — kein drittes bzw. zwanzigstes Vorkommen.
- **Träger-Nachzug (`AGENTS.md` §3.13, real durchgeführter Suchlauf):**
  Gesucht (`grep -rn "Kotlin"`/`"pgchangefeed-kotlin"`) über
  `docs/user/benutzerhandbuch.md`, `sdks/kotlin/pgchangefeed-kotlin/README.md`
  und das Root-`README.md`. Gefunden und nachgezogen: das
  `pgchangefeed-kotlin/README.md`s „Status"-Absatz (behauptete
  „gRPC … follows in a subsequent release") und
  `docs/user/benutzerhandbuch.md`s HTTP-`**SDK:**`-Absatz (dieselbe
  veraltete Aussage auf Deutsch) — beide korrigiert, plus neuer
  `**SDK:**`-Absatz im gRPC-Abschnitt des Handbuchs (Version `1.36` →
  `1.37`). Nicht gefunden/nicht nachgezogen: das Root-`README.md` nennt
  `sdks/csharp/`/`sdks/python/` mit ihren Registry-Links, aber **keine**
  `sdks/kotlin/`-Zeile — dieser Zustand bestand bereits vor diesem Slice
  (seit `slice-sdk-kotlin-projektgeruest`/`-http-client-flaeche`, keine
  Eigenschaft, die *dieser* Slice bewegt) und liegt außerhalb des in §3
  geplanten Datei-Umfangs; an den Reviewer/Koordinator gemeldet statt
  still mitgeändert.
- **Folge-Slices:** keine aus diesem Slice selbst erwartet — Umfang bleibt
  innerhalb der Welle (Pack-Werkzeug, Publish-Workflow bleiben eigene,
  bereits geplante Slices).
- **Risiken aus §6:** alle drei tragen einen Ausgang (siehe §6) — zwei
  eingetreten/aufgelöst, einer (Fake/Stub-Grenze der Authn-Boundary)
  bewusst weiter offen bis zu einem künftigen Integrationsbeleg.
- **Drei Paarungen:** dieser Slice gehört zu
  [welle-sdk-kotlin-lh-fa-sst-009](../welle-sdk-kotlin-lh-fa-sst-009.md)
  (noch offen) — die Prüfung läuft regelkonform bei deren Closure.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Sub-Area `sdks/kotlin/` — mit
`slice-sdk-kotlin-projektgeruest` bereits eröffnet (GF, siehe dessen §8),
keine erneute Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`
(bereits verkörpert) betrifft diesen Slice über den DoD-Punkt „Doku-Update"
oben; `BEO-PGC/workaround-uebersieht-etabliertes-muster-im-bestand` (2×) und
`BEO-PGC/arbeit-ueberholt-stehenden-traeger` (19×) über §3/§6. Kein weiterer
Treffer für diese Sub-Area.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Fortsetzung von
`slice-sdk-kotlin-projektgeruest`).
