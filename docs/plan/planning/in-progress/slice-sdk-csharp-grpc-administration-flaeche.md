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

**Verantwortlich:** Implementer-Agent.

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

- [x] `LH-FA-SST-009` erfüllt: `PgChangeFeedAdministrationClient` deckt alle
      elf RPCs mit typisierten Requests/Responses und einer typisierten
      Fehlerklasse für die gRPC-Status-Codes ab (analog
      `PgChangeFeedHttpClient`s `PgChangeFeedException`-Muster); Unit-Tests
      je RPC. (Requests/Responses sind die generierten Protobuf-Nachrichten
      direkt, kein eigener Modell-Layer — siehe „Abweichung vom Plan" unten.)
- [x] `ADR-0133` erfüllt: `StreamChangesAsync` trägt optionale
      `schema`/`table`-Parameter, leer = ungefiltert (Regressionstest für
      den parameterlosen Aufruf).
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`), kein Self-Review.
- [x] `docs/user/benutzerhandbuch.md`: beide gRPC-Abschnitte („Zugriff über
      den gRPC-Change-Stream", „Zugriff über die gRPC-Verwaltungs-API")
      nennen C# jetzt mit der vollen Fläche statt „SDK folgt"; Versions­historie
      nachgezogen. `sdks/csharp/README.md` nachgezogen.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder eine weitere `evidence/`-Datei; keine Beobachtung
      angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang.
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
  gegenzuprüfen. — **Ausgang:** weiter offen, kein neuer Fund: die
  proto3-Zero-Value-Semantik (`ByteString.Empty` für ein leeres
  `old_image`/`new_image`, `Known = false` statt `null` für
  `ConsumerLag`/`HeartbeatStatus`) ist in C# real getestet
  (`PgChangeFeedAdministrationClientRetentionAndChangesTests`,
  `PgChangeFeedAdministrationClientDiagnoseTests`) und deckungsgleich mit
  dem Go-/Kotlin-Formvorbild — kein Randfall gefunden, Zähler bleibt bei 2×
  → Register unverändert.
- **Server-Fehler analog zum Kotlin-Beispiel-Client-Fund** (`EnableTable`
  über gRPC aktualisiert den laufenden Capture-Prozess nicht ohne
  Neustart, separat gefixt) — ein Realserver-Test dieses Slice könnte
  denselben (jetzt gefixten) oder einen verwandten Effekt erneut zeigen,
  falls der Fix zum Zeitpunkt der Implementierung noch nicht gepusht ist.
  — **Ausgang:** entfallen (nicht ausgelöst) — dieser Slice fährt keinen
  Realserver-Test (Unit-Tests gegen `FakeUnaryCallInvoker`/`FakeCallInvoker`,
  kein Docker-Compose-Rundlauf); ein Realserver-Beleg für
  `PgChangeFeedAdministrationClient` bleibt Gegenstand eines späteren
  Integrationstest-Slices (`make test-sdk-csharp-integration`), zu dessen
  Zeitpunkt der Assembler-Sync-Fix (`internal/bootstrap/assemblersync.go`,
  siehe Aufgabenstellung) bereits gepusht ist.
- **`protoc`-Namenskonflikt** (Kotlin-Beispiel-Client-Fund:
  `AdministrationOuterClass` statt `Administration` als generierter
  Klassenname) — für C# (`Grpc.Tools`) tritt vermutlich ein anderes,
  eigenes Namensschema auf; vor dem Schreiben von Business-Code real am
  generierten Code verifizieren, nicht annehmen. — **Ausgang:** entfallen
  (kein Konflikt in C#s Codegen) — real verifiziert: `Administration.AdministrationClient`
  (Grpc.Tools generiert keinen `OuterClass`-Namen für den C#-Namespace
  `Cdc.Administration.V1`), `make sdk-pack-csharp` baut und testet grün
  (109 Tests, `Passed: 109`) mit genau diesem Bezeichner, identisch zum
  bereits gepushten `examples/csharp/grpc-client`.

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „die drei
SDK-Packages decken die gRPC-Verwaltungs-API/den Stream-Filter noch nicht
ab" in `docs/user/benutzerhandbuch.md`/`sdks/csharp/README.md`; beide
Stände gemessen, Parent `046782ae`):**

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| „drei SDK-Packages" als offener Folge-Schritt (gRPC-Fläche) | `git grep -n "drei SDK-Packages"` über `docs/user/benutzerhandbuch.md sdks/csharp/README.md`, Parent `046782ae` und Diff (Block unten, Zeilen 1–2) | Parent: 8 Trefferzeilen — Zeile 1201 (HTTP `GET /diagnose`-SDK-Methode, außerhalb des Scopes dieses Slice), Zeile 1309 (gRPC-Stream-Filterung), Zeile 1490 (gRPC-Verwaltungs-API), Zeile 1519 (SSE-Filterung, außerhalb des Scopes — nur der gRPC-Stream und die Verwaltungs-API sind Gegenstand), fünf Versionshistorie-Zeilen 1.74–1.77 (Chronik, unberührt). Diff: 6 (Zeilen 1309 und 1490 aufgelöst — C# jetzt namentlich genannt; 1201, 1519 und die fünf Chronik-Zeilen bleiben unverändert korrekt). Nichtgefunden: keine weitere Stelle außerhalb dieser beiden Dateien, die den gRPC-Fähigkeitsstand der drei SDK-Packages nennt (Suchraum bewusst auf die beiden Nutzer-Dokumente beschränkt — `docs/plan/adr/*` sind `Accepted`/unberührbar, `docs/plan/planning/welle-sdk-grpc-administration-flaeche.md` ist Eigentum der Welle-Closure, nicht dieses Slice). | Zeile 1309 (Filterungs-Absatz) und Zeile 1490 (Verwaltungs-API-Absatz) auf `PgChangeFeed.Client` als abdeckend umgeschrieben; Zeile 1201 (HTTP-Diagnose) und Zeile 1519 (SSE-Filter) bleiben unverändert korrekt offen — beide außerhalb des Scopes dieses Slice; die fünf Chronik-Zeilen bleiben unverändert (Versionshistorie ist append-only). |
| „cannot be filtered by table" (README, veraltete Pauschalaussage nach dem Stream-Filter-Zusatz) | `git grep -n "cannot be filtered"` über `sdks/csharp/README.md`, Parent `046782ae` und Diff (Block unten, Zeile 3) | Parent: 1 Trefferzeile („The gRPC and SSE streams cannot be filtered by table…"), jetzt falsch für gRPC. Diff: 0. Nichtgefunden: keine weitere Pauschalaussage dieser Art in der README. | Satz präzisiert: „The gRPC stream can be filtered by schema/table (…); the SSE stream cannot yet; …" — SSE-Teilaussage bleibt korrekt bestehen. |

```suchlauf
046782ae 8 -n "drei SDK-Packages" -- docs/user/benutzerhandbuch.md sdks/csharp/README.md
diff 6 -n "drei SDK-Packages" -- docs/user/benutzerhandbuch.md sdks/csharp/README.md
046782ae 1 -n "cannot be filtered" -- sdks/csharp/README.md
diff 0 -n "cannot be filtered" -- sdks/csharp/README.md
```

**Abweichung vom Plan (§3):** `sdks/csharp/PgChangeFeed.Client/Grpc/Models/*.cs`
entsteht **nicht** — die elf RPC-Methoden von `PgChangeFeedAdministrationClient`
nehmen die generierten Protobuf-Nachrichten aus `Cdc.Administration.V1`
direkt als Parameter/Rückgabetyp, ohne eigene Zwischenschicht. Begründung:
Ein Protobuf-Message ist bereits eine typisierte C#-Klasse; eine
handgeschriebene Spiegelung von elf Nachrichtenformen hätte keinen
Deserialisierungs-Bedarf zu rechtfertigen (anders als bei
`PgChangeFeedHttpClient`, dessen `Http/Models/*.cs` json-Zieltypen sind).
Dieselbe Entscheidung trägt bereits `PgChangeFeedGrpcClient` für den
Stream (eigener Doc-Kommentar: „There is no separate DTO layer … the
generated message already is the typed form"); diese Klasse folgt
demselben, bereits im Repo etablierten Muster statt es für die
Verwaltungs-API zu brechen. Real geprüft: alle elf Methoden kompilieren
und elf RPC-Testdateien (`PgChangeFeed.Client.Tests/Grpc/PgChangeFeedAdministrationClient*Tests.cs`)
laufen grün gegen die generierten Typen.

## 7. Closure-Notiz

- **Was hat funktioniert:** Das direkte fachliche Vorbild (`examples/csharp/grpc-client`,
  bereits vollständig implementiert/gereviewt/verifiziert) übertrug sich
  ohne Reibung auf die SDK-Fläche — Rechtsklassen, Nachrichtenschema und
  Fehlerform stimmten beim ersten Entwurf bereits exakt. Die Rollen-Sequenz
  (Implementer → Reviewer → Fixrunde → Verifier) fing eine reale
  Beleg-Lücke (F-1) vor dem Merge.
- **Was ging anders als geplant:** kein eigener `Models/*.cs`-DTO-Layer
  (siehe §3 „Abweichung vom Plan") — dieselbe, bereits im Repo etablierte
  Entscheidung wie beim Stream-Client, real geprüft und vom Reviewer als
  vertretbar bestätigt (F-4, kein Befund).
- **Steering-Loop-Eintrag:** keiner — F-1 ist eine neue Instanz der
  bereits verkörperten Regel „Beleg trägt seinen Satz nicht"
  (`.harness/skills/reviewer.md`), kein neuer Mechanismus nötig.
- **Beobachtungs-Register (`../observations/`):** `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`
  um `evidence/slice-sdk-csharp-grpc-administration-flaeche.md` ergänzt
  (HIGH, neue Form „Assertion" trotz Deckel bei 14×, da Schwere > LOW) —
  Zähler steht bei 18×.
- **Folge-Slices:** keiner — der offene Realserver-Beleg für
  `PgChangeFeedAdministrationClient` ist Gegenstand eines späteren,
  eigenständigen `make test-sdk-csharp-integration`-Zugs, hier nicht neu
  benannt (bereits bekannter Folge-Bedarf, siehe §6).
- **Risiken aus §6:** alle drei mit Ausgang — siehe §6 (1× weiter offen
  unter der Register-Schwelle, 2× entfallen — einmal mangels Auslösung,
  einmal real widerlegt).
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
