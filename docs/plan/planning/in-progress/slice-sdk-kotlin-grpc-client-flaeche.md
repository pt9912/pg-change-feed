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

- [ ] `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/…/grpc/` (oder
      gleichwertiger Namensraum) trägt eine öffentliche Client-Klasse mit
      einer Methode, die den `StreamChanges`-RPC öffnet und die Nachrichten
      von [`SPEC-020`](../../../../spec/pflichtenheft.md) (`change_id`,
      `transaction_id`, `source_table_id`, `sequence`, `operation`,
      `old_image`, `new_image`, `schema_version`, `schema`, `table`) an den
      Consumer weiterreicht — Bearer-Token wird bei Konstruktion oder
      Aufruf übergeben, landet im `authorization`-Metadata-Eintrag.
- [ ] `sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts` bekommt dieselben
      gRPC-/Coroutine-Koordinaten wie `examples/kotlin/grpc-client/build.gradle.kts`
      (`io.grpc:grpc-kotlin-stub`, `io.grpc:grpc-netty-shaded`,
      `io.grpc:grpc-bom`, `com.google.protobuf:protobuf-java`,
      `org.jetbrains.kotlinx:kotlinx-coroutines-core`, das
      `com.google.protobuf`-Gradle-Plugin) — Versionen zum Bau-Zeitpunkt
      dieses Slice real neu gemessen, nicht aus `examples/kotlin/` oder
      `ADR-0109`s Kontext-Messung (2026-09-17) unbesehen übernommen
      (`AGENTS.md` §3.12).
- [ ] `sdks/kotlin/Dockerfile` bekommt den zusätzlichen, benannten
      Bau-Kontext `proto` (`--build-context proto=proto`, `COPY --from=proto
      cdc/stream/v1/changestream.proto …`) — ohne ihn bricht der Bau an der
      `COPY`-Zeile ab, kein stiller Fallback (Muster
      `examples/kotlin/Dockerfile`/`harness/mk/examples.mk`, hier auf den
      SDK-Baum übertragen).
- [ ] Eigene Tests decken mindestens: Nachrichtenschema-Vollständigkeit
      (Feld-für-Feld gegen `SPEC-020`, analog
      `internal/adapters/driving/grpc/server_test.go`s Feldvollständigkeits-
      Test) und den Authn-Boundary-Pfad (fehlendes/ungültiges Token →
      `Unauthenticated`, ohne echten Server — ein Fake/Stub des
      generierten Coroutine-Stubs oder gleichwertig, netzlos).
- [ ] Kein Import aus `internal/**`/`cmd/**`/`gen/**` dieses Repos —
      der generierte Stub entsteht im SDK-eigenen Bau aus der `.proto`,
      nicht als Kopie von `gen/**` (`ADR-0109` §Kontext Bindung „Import-
      Grenze, hier ohne Ausnahme") — real geprüft:
      `grep -rn "internal/\|cmd/\|gen/" sdks/kotlin/` liefert höchstens
      einen Treffer als Doku-Kommentar-Zitat des Test-Vorbilds
      (`internal/adapters/driving/grpc/server_test.go`), keinen Import.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: `docs/user/benutzerhandbuch.md` bekommt einen
      SDK-Hinweis für die Kotlin-gRPC-Oberfläche (`ADR-0109` §Konsequenzen
      Folgepflicht 4, **inklusive** des PAT-Hinweises für den Bezug über
      GitHub Packages) — getragen durch die bereits verkörperte
      Selbstprüf-Instruktion und den Reviewer-HIGH-Punkt
      (`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben,
      **falls dieser Slice einen Inventur-Fund auflöst** — entfällt: keine
      Reconciliation-Datei in diesem Repo.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder Beleg in `evidence/`; keine Beobachtung angefallen
      ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
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
  `cannot find symbol` scheitern lassen. **Ausgang:** weiter offen,
  entschieden beim Schreiben — `examples/kotlin/grpc-client` belegt bereits
  eine funktionierende Kombination als Referenz.
- Ein Fake/Stub für die Authn-Boundary könnte den realen
  gRPC-`Unauthenticated`-Status-Pfad nicht exakt nachbilden. **Ausgang:**
  weiter offen — ein realer Rundlauf-Beleg bleibt
  `make test-integration`s bestehendem `tools/harness/grpcclient`
  vorbehalten (Wegwerf-Client, kein SDK-Import); dieses SDK bekommt
  frühestens mit einem Folge-Slice einen eigenen Integrationsbeleg.

## 7. Closure-Notiz

- **Was hat funktioniert:** <wird beim Abschluss ergänzt>
- **Was ging anders als geplant:** <wird beim Abschluss ergänzt>
- **Steering-Loop-Eintrag:** <wird beim Abschluss ergänzt>
- **Beobachtungs-Register (`../observations/`):** <wird beim Abschluss
  ergänzt>
- **Folge-Slices:** keine aus diesem Slice selbst erwartet — Umfang bleibt
  innerhalb der Welle (Pack-Werkzeug, Publish-Workflow bleiben eigene,
  bereits geplante Slices).
- **Risiken aus §6:** <wird beim Abschluss ergänzt>
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
