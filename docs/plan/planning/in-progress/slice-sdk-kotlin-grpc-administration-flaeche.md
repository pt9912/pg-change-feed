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

- [x] `LH-FA-SST-009` erfüllt: `PgChangeFeedAdministrationClient` deckt alle
      elf RPCs mit typisierten Requests/Responses und einer versiegelten
      (`sealed class`) Fehlerklasse für die gRPC-Status-Codes ab (Formvorbild
      die bestehende HTTP-Fehlerklasse); Unit-Tests je RPC. (Requests/
      Responses sind die generierten Protobuf-Nachrichten direkt, kein
      eigener Modell-Layer — siehe „Abweichung vom Plan" unten.)
- [x] `ADR-0133` erfüllt: `streamChanges()` trägt optionale `schema`/`table`-
      Parameter, leer = ungefiltert (Regressionstest für den parameterlosen
      Aufruf).
- [x] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`), kein Self-Review.
- [x] `docs/user/benutzerhandbuch.md`: beide gRPC-Abschnitte nennen Kotlin
      jetzt mit der vollen Fläche — mit dieser Zeile ist die Drei-Sprachen-
      SDK-Matrix für die gRPC-Verwaltungs-API vollständig; Versionshistorie
      nachgezogen. `sdks/kotlin/pgchangefeed-kotlin/README.md` nachgezogen.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — kein neuer
      Eintrag angefallen, siehe §7.
- [x] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen sind getragen — von der Welle-Closure.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/io/github/pt9912/pgchangefeed/grpc/PgChangeFeedAdministrationClient.kt` | neu | elf RPC-Methoden, Formvorbild `examples/kotlin/grpc-client/src/main/kotlin/cdcexamples/grpc/{ConsumerClient,TablesAdminClient,RetentionClient,ChangesClient,DiagnoseClient}.kt` |
| `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/io/github/pt9912/pgchangefeed/grpc/PgChangeFeedGrpcException.kt` | neu | versiegelte (`sealed class`) Fehlerklasse für die gRPC-Status-Codes |
| `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/io/github/pt9912/pgchangefeed/grpc/PgChangeFeedGrpcClient.kt` | update | `streamChanges()` um optionale `schema`/`table`-Parameter erweitern |
| `sdks/kotlin/Dockerfile` | update (nicht geplant, siehe „Abweichung vom Plan") | zweite `COPY --from=proto`-Zeile für `administration.proto` — ohne sie bricht `make sdk-pack-kotlin` am Kotlin-Compiler ab |
| `sdks/kotlin/pgchangefeed-kotlin/src/test/kotlin/io/github/pt9912/pgchangefeed/grpc/{FakeAdministrationTransport,AdministrationTestClientFactory,PgChangeFeedAdministrationClient{Consumer,Table,RetentionAndChanges,Diagnose,ErrorMapping}Test,PgChangeFeedGrpcClientFilterTest}.kt` | neu | je RPC ein Happy-/Boundary-/Negative-Fall, Regressionstest für `streamChanges()` ohne Filter |
| `docs/user/benutzerhandbuch.md` | update | beide gRPC-Abschnitte, Versionshistorie |
| `sdks/kotlin/pgchangefeed-kotlin/README.md` | update | Administration-Client dokumentieren |

**Abweichung vom Plan:** `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/io/github/pt9912/pgchangefeed/grpc/model/*.kt`
entsteht **nicht** — die elf RPC-Methoden von `PgChangeFeedAdministrationClient`
nehmen die generierten Protobuf-Nachrichten aus `AdministrationOuterClass`
direkt als Parameter/Rückgabetyp, ohne eigene Zwischenschicht (dieselbe
Entscheidung wie beim C#-Geschwister-SDK und wie beim bestehenden
`PgChangeFeedGrpcClient` für den Stream). Zusätzlich, nicht im ursprünglichen
Plan: `sdks/kotlin/Dockerfile` braucht eine zweite `COPY --from=proto`-Zeile
für `administration.proto` (siehe §7 „Was ging anders als geplant").

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
  bestätigt: `a40b4809 fix(bootstrap): EnableTable/DisableTable über
  HTTP/gRPC aktualisieren den laufenden Assembler`, Datei
  `internal/bootstrap/assemblersync.go` vorhanden) — dieser Slice fährt
  ohnehin keinen Realserver-Test (Unit-Tests gegen
  `FakeAdministrationTransport`/`FakeGrpcStreamTransport`, kein
  Docker-Compose-Rundlauf), der Realserver-Beleg bleibt Gegenstand von
  `make test-sdk-kotlin-integration`.

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „`pgchangefeed-kotlin`
deckt die gRPC-Verwaltungs-API/den Stream-Filter noch nicht ab" in
`docs/user/benutzerhandbuch.md`/`sdks/kotlin/pgchangefeed-kotlin/README.md`;
beide Stände gemessen, Parent `fd39b68b`):**

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| „nehmen ihn noch nicht als eigenen Aufrufparameter" (Stream-Filter-Absatz) | `git grep -n "nehmen ihn noch nicht als eigenen Aufrufparameter"` über `docs/user/benutzerhandbuch.md`, Parent `fd39b68b` und Diff (Block unten, Zeilen 1–2) | Parent: 1 Trefferzeile (Zeile 1312, direkt nach dem C#-/Python-Filterhinweis) — nannte `pgchangefeed`/`pgchangefeed-kotlin` gemeinsam als noch nicht filterfähig; zum Zeitpunkt dieser Messung hatte der parallele Python-Slice den Python-Teil bereits uncommittet im Arbeitsbaum korrigiert (`pgchangefeed` filtert bereits), der Kotlin-Teil war noch offen. Diff: 0 — Satz umgeschrieben, `pgchangefeed-kotlin` jetzt namentlich mit `streamChanges(schema, table)` genannt. Nichtgefunden: keine weitere Stelle, die den Stream-Filter-Stand von `pgchangefeed-kotlin` nennt. | Absatz auf alle drei Packages als filterfähig umgeschrieben (Zeile 1301 ff.); der `**SDK:**`-Absatz im gRPC-Change-Stream-Abschnitt (Zeile ~1372 ff.) ergänzt zusätzlich die zwei Parameter im Fließtext. |
| „Methodensatz dafür" (Verwaltungs-API-Absatz, gemeinsame Offen-Aussage für `pgchangefeed`/`pgchangefeed-kotlin`) | `git grep -n "Methodensatz dafür"` über `docs/user/benutzerhandbuch.md`, Parent `fd39b68b` und Diff (Block unten, Zeilen 3–4) | Parent: 1 Trefferzeile (Zeile 1510) — „keines von ihnen trägt bislang einen Methodensatz dafür" für beide verbliebenen Packages. Diff: 0 — Absatz durch einen eigenen Kotlin-`**SDK:**`-Absatz ersetzt (der Python-Absatz stand zum Messzeitpunkt bereits uncommittet im Arbeitsbaum). Nichtgefunden: keine weitere Stelle. | Neuer Absatz nach dem Python-Absatz eingefügt, nennt alle elf RPC-Methoden und die versiegelte Fehlerklasse; schließt mit „Drei-Sprachen-SDK-Matrix … vollständig". |
| „cannot be filtered by table" (Kotlin-README, veraltete Pauschalaussage) | `git grep -n "cannot be filtered by table"` über `sdks/kotlin/pgchangefeed-kotlin/README.md`, Parent `fd39b68b` und Diff (Block unten, Zeile 5–6) | Parent: 1 Trefferzeile („The gRPC and SSE streams cannot be filtered by table…"), jetzt falsch für gRPC. Diff: 0. Nichtgefunden: keine weitere Pauschalaussage dieser Art in der README. | Satz präzisiert: „The gRPC stream can be filtered by schema/table (`streamChanges(schema, table)`); the SSE stream cannot yet; …" — SSE-Teilaussage bleibt korrekt bestehen. |

```suchlauf
fd39b68b 1 -n "nehmen ihn noch nicht als eigenen Aufrufparameter" -- docs/user/benutzerhandbuch.md
diff 0 -n "nehmen ihn noch nicht als eigenen Aufrufparameter" -- docs/user/benutzerhandbuch.md
fd39b68b 1 -n "Methodensatz dafür" -- docs/user/benutzerhandbuch.md
diff 0 -n "Methodensatz dafür" -- docs/user/benutzerhandbuch.md
fd39b68b 1 -n "cannot be filtered by table" -- sdks/kotlin/pgchangefeed-kotlin/README.md
diff 0 -n "cannot be filtered by table" -- sdks/kotlin/pgchangefeed-kotlin/README.md
```

## 7. Closure-Notiz

- **Was hat funktioniert:** Das direkte fachliche Vorbild
  (`examples/kotlin/grpc-client`) und das C#-Geschwister-SDK
  (`PgChangeFeedAdministrationClient`/`PgChangeFeedGrpcException`) trugen die
  Nachrichtenschema-, Rechtsklassen- und Fehlerform-Entscheidungen bereits
  vollständig vor — Übertragung auf Kotlin ohne neue Design-Fragen. Der
  `AdministrationTransport`-Schnitt (ein Interface mit elf suspend-Methoden,
  gespiegelt von `GrpcStreamTransport`) hielt die elf RPC-Tests netzlos ohne
  einen echten `Channel`.
- **Was ging anders als geplant:** `sdks/kotlin/Dockerfile` kopierte im
  bestehenden Bau-Kontext nur `changestream.proto`, nicht
  `administration.proto` — der erste `make sdk-pack-kotlin`-Lauf scheiterte
  am Kotlin-Compiler mit „Unresolved reference 'administration'" über die
  gesamte neue Datei. Ursache: `ADR-0133`s Vorbild-Slice (C#) hatte diese
  zweite `COPY --from=proto`-Zeile bereits in `sdks/csharp/Dockerfile`
  ergänzt (slice-sdk-csharp-grpc-administration-flaeche), die Kotlin-Fläche
  war zu diesem Zeitpunkt noch nicht nachgezogen — kein Plan-Fehler dieses
  Slice, sondern eine fehlende Parallel-Übertragung, hier direkt behoben
  (zweite `COPY`-Zeile analog zu `sdks/csharp/Dockerfile`). Plan §3 wird
  dafür um `sdks/kotlin/Dockerfile` ergänzt (siehe „Abweichung vom Plan"
  unten).
- **Steering-Loop-Eintrag:** keiner — die Dockerfile-Lücke ist eine
  Bau-Konfigurationslücke, kein wiederkehrendes Agenten-Verhalten.
- **Beobachtungs-Register (`../observations/`):** kein neuer Eintrag —
  `BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall` bleibt bei 2×
  (kein neuer Randfall-Fund, siehe §6); die Dockerfile-Lücke ist keine
  Agenten-Verhaltens-Beobachtung im Sinne des Registers, sondern eine
  zwischen den drei SDK-Slices der Welle nicht mitgezogene Bau-Datei — die
  Welle-Closure prüft, ob die übrigen Sprach-Dockerfiles (Python) dieselbe
  Lücke trugen.
- **Folge-Slices:** keiner aus diesem Slice heraus — `make
  test-sdk-kotlin-integration` (Realserver-Beleg der neuen Fläche) bleibt
  Gegenstand eines eigenen, hier nicht geplanten Folge-Slices, analog zum
  C#-Geschwister.
- **Risiken aus §6:** beide aufgelöst — `drei-sprachen-kopie-divergiert-am-
  randfall` bleibt unter der Schwelle (weiter offen, Register unverändert);
  der Server-Fehler-Risiko ist entfallen (Fix bereits gepusht und real
  bestätigt).
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
