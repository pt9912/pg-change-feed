# Slice sdk-csharp-grpc-client-flaeche: Öffentliche gRPC-Stream-Client-Fläche (`SPEC-020`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** [welle-sdk-csharp-lh-fa-sst-009](../welle-sdk-csharp-lh-fa-sst-009.md).

**Bezug:** [`LH-FA-SST-009`](../../../../spec/lastenheft.md),
[`LH-FA-SST-008`](../../../../spec/lastenheft.md) (Live-Change-Stream —
Boundary: kein Replay im Stream selbst),
[`ADR-0106`](../../adr/0106-csharp-nuget-erstes-sdk-package.md)
Festlegung 1/2 (Umfang, `.proto`-Bezug über Zusatzkontext),
[`ADR-0060`](../../adr/0060-grpc-streaming-mechanismus.md)
(gRPC-Server-Streaming-Vertrag, wird vom SDK benutzt, nicht erweitert).

**Berührte Spec-Stellen:** [`SPEC-020`](../../../../spec/pflichtenheft.md)
(Nachrichtenschema, RPC-Name, Stream-Semantik — das SDK benutzt diese
Festlegungen, verändert sie nicht).

**Verantwortlich:** — bis zur Priorisierung.

**Autor:** Planner-Agent, direkt beauftragt (`ADR-0106` §Konsequenzen
Folgepflicht 1). **Datum:** 2026-09-19.

---

## 1. Ziel und Abgrenzung

**Ziel:** Eine öffentliche, stabile .NET-API-Fläche für den
`StreamChanges`-RPC von [`SPEC-020`](../../../../spec/pflichtenheft.md) — ein
Server-Streaming-Aufruf, der `IAsyncEnumerable<Change>` (oder gleichwertig
idiomatisch) an den Consumer liefert, Authentifizierung über den
`authorization`-Metadata-Eintrag (`Bearer <token>`), dieselbe `.proto`-Bezugsform
wie `examples/csharp/grpc-client` (zusätzlicher, benannter Docker-Bau-Kontext,
kein committeter Stub).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **HTTP-API** — `slice-sdk-csharp-http-client-flaeche` übernimmt das;
  getrennter Draht-Vertrag, keine Überschneidung.
- **Tabellen-granulare Filterung des Streams** — `SPEC-020` selbst trägt
  sie nicht („eine tabellen-granulare Filterung ist nicht Teil dieser
  Version"); das SDK kann keine Fähigkeit anbieten, die der Draht nicht
  hat.
- **Stream-internes Replay** — `SPEC-020`/[`LH-FA-SST-008`](../../../../spec/lastenheft.md)
  Boundary schließt das ausdrücklich aus; ein Consumer, der Replay
  braucht, nutzt den bestehenden Lesezugriffsweg (`GET /changes`,
  `slice-sdk-csharp-http-client-flaeche`), nicht dieses Stream-API.
- **SSE- oder NATS-Vollinhalts-Stream** — `ADR-0106` Festlegung 1 grenzt
  v1 ausdrücklich auf HTTP-API und gRPC-Stream ein (Welle-Plan §6
  Out-of-Scope).
- **CLI-Argument-Parsing wie `examples/csharp/grpc-client`** — dasselbe
  Argument wie bei der HTTP-Fläche: ein SDK-Consumer ruft eine Methode
  auf, parst keine `argv`.
- **Committeter Protobuf-Stub im SDK-Baum** — `ADR-0106` Festlegung 2
  verlangt denselben Bezugsweg wie `examples/csharp/grpc-client`: die
  `.proto` bleibt die einzige Quelle, der Stub entsteht im Bau.

## 2. Definition of Done

- [ ] `sdks/csharp/PgChangeFeed.Client/Grpc/` (oder gleichwertiger
      Namensraum) trägt eine öffentliche Client-Klasse mit einer Methode,
      die den `StreamChanges`-RPC öffnet und die Nachrichten von
      [`SPEC-020`](../../../../spec/pflichtenheft.md) (`change_id`,
      `transaction_id`, `source_table_id`, `sequence`, `operation`,
      `old_image`, `new_image`, `schema_version`, `schema`, `table`) an den
      Consumer weiterreicht — Bearer-Token wird bei Konstruktion oder
      Aufruf übergeben, landet im `authorization`-Metadata-Eintrag.
- [ ] `sdks/csharp/PgChangeFeed.Client.csproj` bekommt die drei
      gRPC-Pakete (`Grpc.Net.Client`, `Grpc.Tools` `PrivateAssets="All"`,
      `Google.Protobuf`) über eine `Directory.Packages.props`
      (zentral gepinnt, exakte Versionen, keine Bereiche — Muster
      `examples/csharp/Directory.Packages.props`, Versionen zum
      Bau-Zeitpunkt neu gemessen, nicht blind übernommen — `AGENTS.md`
      §3.12).
- [ ] `sdks/csharp/Dockerfile` bekommt den zusätzlichen, benannten
      Bau-Kontext `proto` (`--build-context proto=proto`, `COPY --from=proto
      cdc/stream/v1/changestream.proto …`) — ohne ihn bricht der Bau an der
      `COPY`-Zeile ab, kein stiller Fallback (Muster
      `examples/csharp/Dockerfile`/`ADR-0090` Festlegung 2, hier auf den
      SDK-Baum übertragen).
- [ ] Eigene Tests (xUnit) decken mindestens: Nachrichtenschema-Vollständigkeit
      (Feld-für-Feld gegen `SPEC-020`, analog
      `internal/adapters/driving/grpc/server_test.go`s Feldvollständigkeits-
      Test) und den Authn-Boundary-Pfad (fehlendes/ungültiges Token →
      `Unauthenticated`, ohne echten Server — Fake-`CallInvoker` oder
      gleichwertig).
- [ ] Kein Import aus `internal/**`/`cmd/**`/`gen/**` dieses Repos —
      der generierte Stub entsteht im SDK-eigenen Bau aus der `.proto`,
      nicht als Kopie von `gen/**` (`ADR-0106` §Kontext Bindung „Import-
      Grenze, hier ohne Ausnahme").
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: `docs/user/benutzerhandbuch.md` bekommt einen
      SDK-Hinweis für die gRPC-Oberfläche (`ADR-0106` §Konsequenzen
      Folgepflicht 4) — getragen durch die bereits verkörperte
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
      [welle-sdk-csharp-lh-fa-sst-009](../welle-sdk-csharp-lh-fa-sst-009.md)
      (noch offen); die Prüfung läuft regelkonform bei deren Closure.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `sdks/csharp/PgChangeFeed.Client/Grpc/PgChangeFeedGrpcClient.cs` (Arbeitsname) | neu | öffentliche API-Fläche für `StreamChanges` (`SPEC-020`). |
| `sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.csproj` | update | `PackageReference` auf `Grpc.Net.Client`/`Grpc.Tools`/`Google.Protobuf`, `Protobuf Include`-Zeile für den im Bau kopierten Stub. |
| `sdks/csharp/Directory.Packages.props` | neu | zentral gepinnte Versionen der drei gRPC-Pakete, real gemessen zum Bau-Zeitpunkt (Muster `examples/csharp/Directory.Packages.props`). |
| `sdks/csharp/Dockerfile` | update | zusätzlicher `proto`-Bau-Kontext, `COPY --from=proto …`-Zeile. |
| `harness/mk/examples.mk` bzw. neues `harness/mk/sdk.mk` | ggf. update | falls der Pack-Werkzeug-Slice (Folge-Slice) den Aufruf bereits vorwegnehmen muss, um den `--build-context proto=proto`-Zwang für den Bau dieses Slice testbar zu machen — sonst entfällt diese Zeile hier und wandert vollständig in `slice-sdk-csharp-pack-werkzeug`. |
| `sdks/csharp/PgChangeFeed.Client.Tests/GrpcClientTests.cs` (Arbeitsname) | neu | Nachrichtenschema-Vollständigkeit, Authn-Boundary. |
| `docs/user/benutzerhandbuch.md` | update | SDK-Hinweis für die gRPC-Oberfläche, im selben Zug (`ADR-0106` Folgepflicht 4). |

**Ansatz:** Referenzmaterial ist `examples/csharp/grpc-client/Cli.cs`/
`Program.cs` (Draht-Kenntnis für Kanal-Aufbau und Metadata-Header, kein
`ProjectReference`) und
`internal/adapters/driving/grpc/server_test.go` (Feldvollständigkeits-
Testmuster, serverseitig, als Vorbild für die Consumer-seitige Prüfung).

## 4. Trigger

**Start** (`next` → `in-progress`): wenn `slice-sdk-csharp-projektgeruest`
in `done/` liegt (siehe Welle-Plan §4 Reihenfolge) — parallelisierbar zu
`slice-sdk-csharp-http-client-flaeche`, keine gegenseitige Abhängigkeit.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): entfällt aus
  heutiger Sicht — ein RPC, ein Nachrichtentyp.
- `in-progress` → `open` (blockiert — Carveout?): `slice-sdk-csharp-projektgeruest`
  liegt noch nicht in `done/`, oder der Docker-Bau-Kontext-Mechanismus
  (`--build-context proto=proto`) lässt sich aus einem noch unbekannten
  Grund nicht auf den SDK-Baum übertragen (unwahrscheinlich —
  `examples/csharp/grpc-client` belegt den Mechanismus bereits real).

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + Closure-Notiz mit Lerneintrag
geschrieben.

## 6. Risiken und offene Punkte

- Der Bau dieses Slice ist ohne ein aufrufendes `make`-Ziel (das erst
  `slice-sdk-csharp-pack-werkzeug` liefert) nur über einen direkten
  `docker build --build-context proto=proto …`-Aufruf prüfbar — kein
  Komfort-Ziel für diesen Zwischenstand. **Ausgang:** eingetreten,
  akzeptiert: der direkte Aufruf reicht als Beleg für die DoD dieses
  Slice; das komfortable `make`-Ziel folgt bewusst erst mit dem
  Pack-Werkzeug-Slice (Welle-Plan §4 Reihenfolge).
- Ein Fake-`CallInvoker`-Test für die Authn-Boundary könnte den realen
  gRPC-`Unauthenticated`-Status-Pfad nicht exakt nachbilden. **Ausgang:**
  weiter offen — ein realer Rundlauf-Beleg bleibt
  `make test-integration`s bestehendem `tools/harness/grpcclient`
  vorbehalten (Wegwerf-Client, kein SDK-Import); dieses SDK bekommt
  frühestens mit einem Folge-Slice einen eigenen Integrationsbeleg.

## 7. Closure-Notiz

- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Steering-Loop-Eintrag:** <Guide oder Sensor geschärft/ergänzt, oder
  „kein neuer Sensor" — je nach Lauf>.
- **Beobachtungs-Register (`../observations/`):** <neu angelegt | Beleg
  ergänzt | keine Beobachtung angefallen>.
- **Folge-Slices:** keine aus diesem Slice selbst erwartet — Umfang bleibt
  innerhalb der Welle.
- **Risiken aus §6:** <jedes mit genau einem Ausgang — siehe §6>.
- **Drei Paarungen:** dieser Slice gehört zu
  [welle-sdk-csharp-lh-fa-sst-009](../welle-sdk-csharp-lh-fa-sst-009.md)
  (noch offen) — die Prüfung läuft regelkonform bei deren Closure.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Sub-Area `sdks/csharp/` — mit
`slice-sdk-csharp-projektgeruest` bereits eröffnet (GF, siehe dessen §8),
keine erneute Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`
(bereits verkörpert) betrifft diesen Slice über den DoD-Punkt „Doku-Update"
oben; kein weiterer Treffer für diese Sub-Area.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Fortsetzung von
`slice-sdk-csharp-projektgeruest`).
