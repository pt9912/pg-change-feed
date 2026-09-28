# Slice slice-sdk-csharp-grpc-administration-flaeche: C#-SDK — Administration-Client + Stream-Filter

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`.

**Welle:** welle-sdk-grpc-administration-flaeche.

**Bezug:** [`LH-FA-SST-009`](../../../../spec/lastenheft.md) (Scope),
[`ADR-0130`](../../adr/0130-grpc-verwaltungs-api-neun-rpcs.md),
[`ADR-0131`](../../adr/0131-grpc-readchanges-zehnter-rpc.md),
[`ADR-0132`](../../adr/0132-diagnose-ueber-inbound-port-http-grpc.md),
[`ADR-0133`](../../adr/0133-tabellen-granulare-filterung-grpc-sse.md).

**Berührte Spec-Stellen:** — (kein SPEC-/ARC-Eintrag über `LH-FA-SST-009`
hinaus; das Nachrichtenschema selbst steht bereits in `SPEC-031`/`SPEC-020`,
unverändert durch diesen Slice).

**Verantwortlich:** —.

**Autor:** Planner-Agent. **Datum:** 2026-09-28.

---

## 1. Ziel und Abgrenzung

**Ziel:** Das C#-SDK-Package `PgChangeFeed.Client` bekommt einen
`PgChangeFeedAdministrationClient` mit allen elf `Administration`-RPCs
(`RegisterConsumer`, `AcknowledgeConsumer`, `GetConsumerPosition`,
`RemoveConsumer`, `EnableTable`, `DisableTable`, `GetTableStatus`,
`ListTables`, `RunRetention`, `ReadChanges`, `Diagnose`), und der bestehende
`PgChangeFeedGrpcClient.StreamChangesAsync` bekommt die optionalen
`schema`/`table`-Filter-Parameter (`ADR-0133`) — analog zu
`examples/csharp/grpc-client` (bereits vollständig umgesetzt, direktes
fachliches Vorbild für Nachrichtenschema, Rechtsklassen, Fehlerform).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Server-seitige Änderungen** — alle vier ADRs sind serverseitig
  abgeschlossen; ein während dieses Slice gefundener Server-Fehler (wie
  bereits einmal bei der Beispiel-Client-Arbeit geschehen, siehe
  `docs/reviews/review-example-kotlin-grpc-client-verbmatrix.md` INFO-1)
  wird gemeldet, nicht hier repariert — anderer Vorgang.
- **Python-/Kotlin-SDK** — eigene Slices derselben Welle
  (`slice-sdk-python-grpc-administration-flaeche`,
  `slice-sdk-kotlin-grpc-administration-flaeche`); unabhängig baubar/testbar.
- **Neue Beispiel-Clients** — `examples/csharp/grpc-client` ist bereits
  fertig; dieser Slice liest ihn nur als Vorbild, ändert ihn nicht.
- **`GET /diagnose`-SDK-Methode (HTTP)** — dieselbe Fähigkeit steht bereits
  über `Diagnose` (gRPC) im Scope; eine HTTP-Diagnose-SDK-Methode ist
  Gegenstand einer eigenen, hier nicht geplanten Folge-Welle (siehe
  Welle-Plan §6).
- **Breaking Change an `StreamChangesAsync`** — die neuen Parameter sind
  optional/additiv, die bestehende Signatur ohne Filter bleibt
  aufrufbar — kein zweiter Vorgang, aber ausdrücklich als Erhaltungs-Ziel
  benannt, nicht nur implizit angenommen.

## 2. Definition of Done

- [ ] `LH-FA-SST-009` erfüllt: `PgChangeFeedAdministrationClient` deckt alle
      elf RPCs mit typisierten Requests/Responses und einer typisierten
      Fehlerklasse für die gRPC-Status-Codes ab (analog
      `PgChangeFeedHttpClient`s `PgChangeFeedException`-Muster); Unit-Tests
      je RPC.
- [ ] `ADR-0133` erfüllt: `StreamChangesAsync` trägt optionale
      `schema`/`table`-Parameter, leer = ungefiltert (Regressionstest für
      den parameterlosen Aufruf).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`), kein Self-Review.
- [ ] `docs/user/benutzerhandbuch.md`: beide gRPC-Abschnitte („Zugriff über
      den gRPC-Change-Stream", „Zugriff über die gRPC-Verwaltungs-API")
      nennen C# jetzt mit der vollen Fläche statt „SDK folgt"; Versions­historie
      nachgezogen. `sdks/csharp/README.md` nachgezogen.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder eine weitere `evidence/`-Datei; keine Beobachtung
      angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      von der nächsten Welle-Closure (`welle-sdk-grpc-administration-flaeche`).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `sdks/csharp/PgChangeFeed.Client/Grpc/PgChangeFeedAdministrationClient.cs` | neu | elf RPC-Methoden, Formvorbild `examples/csharp/grpc-client/{ConsumerClient,TablesAdminClient,RetentionClient,ChangesClient,DiagnoseClient}.cs` |
| `sdks/csharp/PgChangeFeed.Client/Grpc/Models/*.cs` | neu | typisierte Request-/Response-Modelle, 1:1 aus `proto/cdc/administration/v1/administration.proto` |
| `sdks/csharp/PgChangeFeed.Client/Grpc/PgChangeFeedGrpcException.cs` (oder analog benannt) | neu | typisierte Fehlerklasse für `InvalidArgument`/`Unauthenticated`/`PermissionDenied`/`NotFound`/`Internal`, Formvorbild `PgChangeFeedException` (HTTP) |
| `sdks/csharp/PgChangeFeed.Client/Grpc/PgChangeFeedGrpcClient.cs` | update | `StreamChangesAsync` um optionale `schema`/`table`-Parameter erweitern (additiv) |
| `sdks/csharp/PgChangeFeed.Client.Tests/Grpc/*Tests.cs` | neu | je RPC ein Happy-/Boundary-/Negative-Fall (`LH-FA-SST-009`), Regressionstest für `StreamChangesAsync` ohne Filter |
| `docs/user/benutzerhandbuch.md` | update | beide gRPC-Abschnitte, Versionshistorie |
| `sdks/csharp/README.md` | update | Administration-Client dokumentieren |

## 4. Trigger

**Start** (`next` → `in-progress`): Welle `welle-sdk-grpc-administration-flaeche`
eröffnet (erfüllt — diese Datei entsteht im selben Zug).

**Rückführungen:**

- `in-progress` → `next` (zu groß): falls sich die elf RPCs nicht in einem
  Slice liefern lassen — Anzeichen: mehr als drei Fixrunden nötig.
- `in-progress` → `open` (blockiert): ein Server-Fehler analog zum
  Kotlin-Beispiel-Client-Fund blockiert einen bestimmten RPC — Carveout auf
  den betroffenen RPC, die übrigen zehn liefern.

## 5. Closure-Trigger

DoD vollständig, Review ohne offenes HIGH/MEDIUM, `make gates` grün,
Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- **`BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall`** (2×, offen,
  unter der 3×-Schwelle) — die drei SDK-Sprachen könnten denselben
  Randfall (z. B. ein leeres `old_image`/`new_image`-Byte-Array bei
  `ReadChanges`, ein `null`-Consumer-Lag) unterschiedlich lesen, wenn der
  Implementer nur das Formvorbild kopiert, ohne die eigene Sprach-Laufzeit
  gegenzuprüfen. — **Ausgang:** weiter offen: → Register (unter der
  Schwelle; falls dieser Slice einen neuen Fund beisteuert, `evidence/`
  ergänzen).
- **Server-Fehler analog zum Kotlin-Beispiel-Client-Fund** (`EnableTable`
  über gRPC aktualisiert den laufenden Capture-Prozess nicht ohne
  Neustart, separat gefixt) — ein Realserver-Test dieses Slice könnte
  denselben (jetzt gefixten) oder einen verwandten Effekt erneut zeigen,
  falls der Fix zum Zeitpunkt der Implementierung noch nicht gepusht ist.
  — **Ausgang:** eingetreten: Blockade melden, nicht selbst reparieren |
  entfallen: Fix bereits gepusht, kein erneutes Auftreten | weiter offen.
- **`protoc`-Namenskonflikt** (Kotlin-Beispiel-Client-Fund:
  `AdministrationOuterClass` statt `Administration` als generierter
  Klassenname) — für C# (`Grpc.Tools`) tritt vermutlich ein anderes,
  eigenes Namensschema auf; vor dem Schreiben von Business-Code real am
  generierten Code verifizieren, nicht annehmen. — **Ausgang:** entfallen
  (kein Konflikt in C#s Codegen) | weiter offen: → eigener Beobachtungs-Eintrag,
  falls doch.

## 7. Closure-Notiz

- **Was hat funktioniert:** <wird bei Closure gefüllt>
- **Was ging anders als geplant:** <wird bei Closure gefüllt>
- **Steering-Loop-Eintrag:** <wird bei Closure gefüllt, falls einer entsteht>
- **Beobachtungs-Register (`../observations/`):** <wird bei Closure gefüllt>
- **Folge-Slices:** <wird bei Closure gefüllt, falls einer entsteht>
- **Risiken aus §6:** <wird bei Closure gefüllt — siehe §6>
- **Drei Paarungen:** von der Welle-Closure (`welle-sdk-grpc-administration-flaeche`).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** eine Sub-Area, `sdks/csharp/`
(bestehendes, aktiv gepflegtes Package mit eigenem Release-Workflow,
`harness/conventions.md` §Modus-Deklaration — Repo-Default `PGC`/Greenfield
gilt unverändert für neue Dateien in einer bestehenden Sub-Area).

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen
(vollständige Liste unter `docs/plan/planning/observations/BEO-PGC/`).
Treffer mit Zähler-Stand unter der 3×-Schwelle, als Risiko in §6
aufgenommen: `drei-sprachen-kopie-divergiert-am-randfall` (2×). Kein
Treffer ≥ 3× offen und einschlägig für diese Sub-Area.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (neue Dateien in
einer bestehenden, bereits Greenfield geführten Sub-Area — kein
Brownfield-Bestand, keine Inventur nötig).
