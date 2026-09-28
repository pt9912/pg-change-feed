# Slice slice-sdk-kotlin-grpc-administration-flaeche: Kotlin-SDK — Administration-Client + Stream-Filter

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`.

**Welle:** welle-sdk-grpc-administration-flaeche.

**Bezug:** [`LH-FA-SST-009`](../../../../spec/lastenheft.md) (Scope),
[`ADR-0130`](../../adr/0130-grpc-verwaltungs-api-neun-rpcs.md),
[`ADR-0131`](../../adr/0131-grpc-readchanges-zehnter-rpc.md),
[`ADR-0132`](../../adr/0132-diagnose-ueber-inbound-port-http-grpc.md),
[`ADR-0133`](../../adr/0133-tabellen-granulare-filterung-grpc-sse.md).

**Berührte Spec-Stellen:** — (kein SPEC-/ARC-Eintrag über `LH-FA-SST-009`
hinaus).

**Verantwortlich:** Implementer-Agent.

**Autor:** Planner-Agent. **Datum:** 2026-09-28.

---

## 1. Ziel und Abgrenzung

**Ziel:** Das Kotlin-SDK-Package `pgchangefeed-kotlin` bekommt einen
`PgChangeFeedAdministrationClient` mit allen elf `Administration`-RPCs, und
der bestehende `PgChangeFeedGrpcClient.streamChanges()` bekommt die
optionalen `schema`/`table`-Filter-Parameter (`ADR-0133`) — analog zu
`examples/kotlin/grpc-client` (bereits vollständig umgesetzt, direktes
fachliches Vorbild).

**Wichtiger, real verifizierter Vorbild-Fund (aus dem Beispiel-Client-Review
`docs/reviews/review-example-kotlin-grpc-client-verbmatrix.md`):** Wegen des
Namenskonflikts zwischen der Datei `administration.proto` und dem darin
definierten Dienst `Administration` legt `protoc`s Kotlin-/Java-Codegen alle
Administration-Nachrichtentypen unter `AdministrationOuterClass.*` ab (nicht
`Administration.*`), während der Coroutine-Stub weiterhin
`AdministrationGrpcKt.AdministrationCoroutineStub` heißt — real am
generierten Bytecode verifiziert. Dieser Slice übernimmt dieselbe
Import-Form, statt sie erneut zu entdecken.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Server-seitige Änderungen** — anderer Vorgang, siehe C#-Geschwister-Slice
  §1.
- **Der real gefundene `EnableTable`/`DisableTable`-Server-Fehler** (siehe
  §6) — separat gefixt, nicht Gegenstand dieses Slice.
- **C#-/Python-SDK** — eigene Slices derselben Welle.
- **Neue Beispiel-Clients** — `examples/kotlin/grpc-client` ist bereits
  fertig, dieser Slice liest ihn nur als Vorbild.
- **`diagnose()`-SDK-Methode (HTTP)** — dieselbe Fähigkeit steht bereits
  über `diagnose()` (gRPC) im Scope; HTTP-Diagnose ist Gegenstand einer
  eigenen, hier nicht geplanten Folge-Welle.
- **Breaking Change an `streamChanges()`** — die neuen Parameter sind
  optional/additiv, die bestehende Signatur ohne Filter bleibt aufrufbar.

## 2. Definition of Done

- [ ] `LH-FA-SST-009` erfüllt: `PgChangeFeedAdministrationClient` deckt alle
      elf RPCs mit typisierten Requests/Responses und einer versiegelten
      (`sealed class`) Fehlerklasse für die gRPC-Status-Codes ab (Formvorbild
      die bestehende HTTP-Fehlerklasse); Unit-Tests je RPC.
- [ ] `ADR-0133` erfüllt: `streamChanges()` trägt optionale `schema`/`table`-
      Parameter, leer = ungefiltert (Regressionstest für den parameterlosen
      Aufruf).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`), kein Self-Review.
- [ ] `docs/user/benutzerhandbuch.md`: beide gRPC-Abschnitte nennen Kotlin
      jetzt mit der vollen Fläche — mit dieser Zeile ist die Drei-Sprachen-
      SDK-Matrix für die gRPC-Verwaltungs-API vollständig; Versionshistorie
      nachgezogen. `sdks/kotlin/pgchangefeed-kotlin/README.md` nachgezogen.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen sind getragen — von der Welle-Closure.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/io/github/pt9912/pgchangefeed/grpc/PgChangeFeedAdministrationClient.kt` | neu | elf RPC-Methoden, Formvorbild `examples/kotlin/grpc-client/src/main/kotlin/cdcexamples/grpc/{ConsumerClient,TablesAdminClient,RetentionClient,ChangesClient,DiagnoseClient}.kt` |
| `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/io/github/pt9912/pgchangefeed/grpc/model/*.kt` | neu | typisierte Request-/Response-Modelle, Import über `AdministrationOuterClass.*` (siehe §1-Fund) |
| `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/io/github/pt9912/pgchangefeed/grpc/PgChangeFeedGrpcException.kt` (oder analog benannt) | neu | `sealed class` für die gRPC-Status-Codes |
| `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/io/github/pt9912/pgchangefeed/grpc/PgChangeFeedGrpcClient.kt` | update | `streamChanges()` um optionale `schema`/`table`-Parameter erweitern |
| `sdks/kotlin/pgchangefeed-kotlin/src/test/kotlin/.../grpc/*Test.kt` | neu | je RPC ein Happy-/Boundary-/Negative-Fall, Regressionstest für `streamChanges()` ohne Filter |
| `docs/user/benutzerhandbuch.md` | update | beide gRPC-Abschnitte, Versionshistorie |
| `sdks/kotlin/pgchangefeed-kotlin/README.md` | update | Administration-Client dokumentieren |

## 4. Trigger

**Start** (`next` → `in-progress`): Welle `welle-sdk-grpc-administration-flaeche`
eröffnet (erfüllt).

**Rückführungen:**

- `in-progress` → `next` (zu groß): mehr als drei Fixrunden nötig.
- `in-progress` → `open` (blockiert): ein Server-Fehler blockiert einen
  bestimmten RPC — Carveout auf den betroffenen RPC.

## 5. Closure-Trigger

DoD vollständig, Review ohne offenes HIGH/MEDIUM, `make gates` grün,
Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- **`BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall`** (2×, offen) —
  dieselbe Randfall-Divergenz-Gefahr wie bei den Geschwister-Slices, hier
  speziell: der `protoc`-Namenskonflikt (§1) könnte den Implementer dazu
  verleiten, den generierten Namen anders zu behandeln als das
  Beispiel-Client-Vorbild — real gegen den generierten Code prüfen. —
  **Ausgang:** weiter offen: → Register (unter der Schwelle).
- **Real gefundener, bereits gefixter Server-Fehler:** `EnableTable`/
  `DisableTable` über gRPC/HTTP aktualisierten den laufenden Capture-Prozess
  nicht (dauerhafter Datenverlust bis zum nächsten Neustart) —
  gefunden beim Testen von `examples/kotlin/grpc-client`, in einem
  separaten, hochprioritären Zug bereits behoben. Prüfen, ob der Fix zum
  Zeitpunkt der Implementierung dieses Slice gepusht ist (`internal/bootstrap/assemblersync.go`),
  bevor ein Realserver-Test dieses Slice `enable-table` gegen eine frische
  Tabelle nutzt. — **Ausgang:** entfallen (Fix bereits gepusht, real
  bestätigt) | weiter offen (Fix noch nicht gepusht — dann keinen eigenen
  Realserver-Test auf frischer Tabelle bauen, sondern auf eine bereits
  aktivierte Tabelle ausweichen).

## 7. Closure-Notiz

- **Was hat funktioniert:** <wird bei Closure gefüllt>
- **Was ging anders als geplant:** <wird bei Closure gefüllt>
- **Steering-Loop-Eintrag:** <wird bei Closure gefüllt, falls einer entsteht>
- **Beobachtungs-Register (`../observations/`):** <wird bei Closure gefüllt>
- **Folge-Slices:** <wird bei Closure gefüllt, falls einer entsteht>
- **Risiken aus §6:** <wird bei Closure gefüllt — siehe §6>
- **Drei Paarungen:** von der Welle-Closure (`welle-sdk-grpc-administration-flaeche`).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** eine Sub-Area, `sdks/kotlin/`
(bestehendes, aktiv gepflegtes Package mit eigenem Release-Workflow).

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen.
Treffer unter der 3×-Schwelle, als Risiko in §6 aufgenommen:
`drei-sprachen-kopie-divergiert-am-randfall` (2×). Kein Treffer ≥ 3× offen
und einschlägig.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (neue Dateien in
einer bestehenden, bereits Greenfield geführten Sub-Area).
