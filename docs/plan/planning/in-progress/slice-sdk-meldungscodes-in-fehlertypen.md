# Slice sdk-meldungscodes-in-fehlertypen: Meldungscodes des Servers als Eigenschaft der Fehlertypen der drei SDK-Packages

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — es gibt keine Closure-Bedingung, die von der DoD dieses
Slice verschieden ist (Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht).

**Bezug:** [`LH-FA-SST-009`](../../../../spec/lastenheft.md) (SDK-Packages),
[`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md) (Meldungscodes;
Teilfrage 3: HTTP-Feld `code`, gRPC-Statusdetail `ErrorInfo`; dort der Satz „die SDKs
ändern sich in dieser ADR nicht“),
[`ADR-0145`](../../adr/0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md) (entscheidet
A1 bis A3 dieses Plans; schärft `ADR-0144` Festlegung 3),
[Architect-Verdikt](../../../reviews/architect-verdict-sdk-meldungscodes-in-fehlertypen.md),
[`ADR-0134`](../../adr/0134-sdk-public-doc-check-gate-make-gates.md) (keine interne
Kennung unter `sdks/`),
[`ADR-0110`](../../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md) §Entscheidung
Festlegung 2 (Mechanik der SDK-Realserver-Tiers).

**Berührte Spec-Stellen:** — (die Wire-Form steht in
[`SPEC-018`](../../../../spec/pflichtenheft.md) und `SPEC-031`; die drei SDK-Zeilen
`SPEC-026` bis `SPEC-028` nennen Fehlertypen nicht, gemessen:
`git grep -n -i -E 'PgChangeFeed[A-Za-z]*(Exception|Error)|Fehlertyp' 81d0fa96 -- spec`
druckt 0 Zeilen; die SDKs setzen die Wire-Form um und ändern sie nicht).

**Verantwortlich:** — (gesetzt beim Übergang `open` → `next`).

**Autor:** Planner-Agent, Freigabe des Auftraggebers vom 2026-10-03 (auf die Frage, ob
geplant und umgesetzt werden soll: „Ja das können wir noch machen“), Folge der
Meldungscodes-Slices des Servers (Release 0.6.0). **Datum:** 2026-10-03.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein Anwender der drei SDK-Packages (C# `PgChangeFeed.Client`, Kotlin
`pgchangefeed-kotlin`, Python `pgchangefeed`; Stand 0.5.0) liest den Meldungscode
einer Fehlerantwort des Servers (`PCF-<E|W|I><4 Ziffern>`) als Eigenschaft des
gefangenen Fehlertyps, statt den Klartext zu parsen: beim HTTP-/SSE-Client aus dem
Feld `code` des Fehlerkörpers, beim gRPC-Administrations-Client aus `reason` des
Statusdetails `google.rpc.ErrorInfo`; ohne Code auf dem Draht ist die Eigenschaft
„leer“ (`null`/`None`). Das Feld `error_code` der Diagnose kommt über die aus der
`.proto` erzeugten Stubs ohne Mapper-Code an und wird belegt, nicht gebaut. Drei
Liefer-Punkte:

- (A) **C#** `PgChangeFeed.Client`: Eigenschaft an `PgChangeFeedException` und
  `PgChangeFeedGrpcException`, Parser für Fehlerkörper und Statusdetail, Tests mit
  Fake-Transport, ein Realserver-Fall je Zugriffsweg (HTTP, gRPC-Administration), README;
- (B) **Kotlin** `pgchangefeed-kotlin`: dasselbe an `PgChangeFeedException` und
  `PgChangeFeedGrpcException`;
- (C) **Python** `pgchangefeed`: dasselbe an `PgChangeFeedError` und
  `PgChangeFeedGrpcError`.

Je Sprache dieselbe Eingabetabelle (§3), dieselben Randfälle, derselbe Realserver-Fall
(`BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall`).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Versionsänderung, Tag und Release der Packages (0.6.0)** — jede Veröffentlichung ist
  eine Freigabe mit externen Konten (NuGet, PyPI, GitHub Packages, Cloudsmith); der
  Hauptlauf führt den Release als Folgeschritt **nach** der Verifikation aus
  (`docs/maintainer/releasing.md` §4, SDK-Release je Sprache). Der Slice ändert weder
  `PgChangeFeed.Client.csproj` (`<Version>`), `build.gradle.kts` (zwei Stellen) noch
  `pyproject.toml`; er nennt den Release als Folgeschritt und die Release-Hinweise (§6).
- **Server-Änderungen** — der Server trägt die Codes seit `ADR-0144`; ein gefundener
  Server-Fehler wird gemeldet, nicht hier repariert.
- **Codes für Auth-Fehler** (HTTP `401`/`403`, gRPC `Unauthenticated`/`PermissionDenied`)
  — der Server vergibt dort bewusst keinen Code (Entscheidung `ADR-0144`, Frage E9 dort
  ist offen); die SDKs liefern für diese Fälle leer. Ein späterer Code ist eine
  Server-Entscheidung, die SDKs lesen ihn dann ohne Änderung (Durchreichen, §3).
- **Typisierte Fehler der gRPC-Stream-Clients** — der Stream-Client jeder Sprache lässt
  den gRPC-Fehler roh durch (`RpcException`, `StatusException`, `grpc.RpcError`; in den
  README-Zeilen „… as `grpc.RpcError` while iterating“ so zugesagt); nur der
  Administrations-Client bildet auf eigene Typen ab. Ein Stream-Fehler mit `ErrorInfo`
  (der Server sendet ihn bei `ChangeStream` ohne Broadcaster, `internal/adapters/driving/grpc/server.go`)
  bleibt über die native Statusdetail-API des gRPC-Clients lesbar; das README sagt das.
  Eine Typisierung der Streams wäre eine neue Fehlerhierarchie, kein Zusatzfeld.
- **NATS-Vollinhalts-Client** — der Weg trägt keinen Fehlerkörper des Servers; sein
  `MalformedMessage`-Fehler bekommt keine Eigenschaft.
- **Beispiele unter `examples/`** — gemessen, dass sie die Fehlerform nicht auswerten:
  `git grep -n -E 'ErrorInfo|status-details' 81d0fa96 -- examples` druckt 0 Zeilen; sie
  bleiben unverändert.
- **Neue Laufzeit-Abhängigkeit** — entschieden: keine (§4, A1; `ADR-0145`
  Festlegung 3); das Statusdetail liest je Sprache eine interne Hilfsfunktion mit der
  vorhandenen Protobuf-Laufzeit.

## 2. Definition of Done

- [ ] [`LH-FA-SST-009`](../../../../spec/lastenheft.md) (A) C#: `PgChangeFeedException`
      (HTTP und SSE) und `PgChangeFeedGrpcException` tragen die Eigenschaft
      `MessageCode` (`string?`); Code aus dem Feld `code` des Fehlerkörpers bzw. aus
      `ErrorInfo.reason`; Tests für jede Zeile der Eingabetabelle (§3) mit Fake-Transport
      (HTTP-Handler, SSE-Handler, `FakeUnaryCallInvoker` mit Trailern); ein Realserver-Fall
      je Zugriffsweg in `PgChangeFeed.Client.Integration`; README (Fehlerbehandlung,
      kein „Upgrading“-Eintrag, `ADR-0145` Festlegung 4); je Typ inklusive der Basis der
      bisherige Konstruktor neben dem neuen mit `string? messageCode` hinten (Überladung);
      ein Test, dass der alte Konstruktor weiter funktioniert und `MessageCode` leer ist.
      *Zu belegen durch:* `make sdk-pack-csharp` (Tests laufen
      im Bau; ein Lauf aus dem Docker-Schicht-Cache druckt keine Testzeile — Belegbefehl ist
      der Bau mit `--no-cache` der Test-Stufe oder eine Mutation, die den Test rot färbt),
      `make test-sdk-csharp-integration`, `make examples-csharp`.
- [ ] [`LH-FA-SST-009`](../../../../spec/lastenheft.md) (B) Kotlin:
      `PgChangeFeedException` und `PgChangeFeedGrpcException` tragen `messageCode`
      (`String?`) über `@JvmOverloads` mit `messageCode: String? = null` am
      Primär-Konstruktor; Tests, Realserver-Fälle und README (ohne „Upgrading“-Eintrag)
      wie (A) mit derselben Eingabetabelle; ein Test, dass der Aufruf ohne `messageCode`
      weiter funktioniert und `messageCode` leer ist. *Zu belegen durch:* `make sdk-pack-kotlin` (gleiche Cache-Bedingung),
      `make test-sdk-kotlin-integration`, `make examples-kotlin`.
- [ ] [`LH-FA-SST-009`](../../../../spec/lastenheft.md) (C) Python: `PgChangeFeedError`
      und `PgChangeFeedGrpcError` tragen `message_code` (`str | None`, keyword-only
      `*, message_code=None`; der vorhandene Name `PgChangeFeedGrpcError.code` bleibt der
      `grpc.StatusCode`); Tests, Realserver-Fälle und README wie (A) mit derselben
      Eingabetabelle; ein Test, dass der Aufruf ohne `message_code` weiter funktioniert
      und die Eigenschaft `None` ist; `test_readme_examples` bleibt grün. *Zu belegen durch:*
      `make sdk-pack-python` (gleiche Cache-Bedingung), `make test-sdk-python-integration`.
- [ ] [`ADR-0134`](../../adr/0134-sdk-public-doc-check-gate-make-gates.md) und Nachzug:
      `make sdk-public-doc-check` Exit 0 (die Kennung `PCF-…` ist eine Anwender-Kennung
      und kein Treffer des Wächters; gemessen: das Muster
      `\b(SPEC|ADR|ARC)-[0-9]+|\bLH-(FA|QA)-[A-Z]{3}-[0-9]+|\b(slice|welle)-[a-z0-9]` in
      `tools/harness/sdk-public-doc-check.sh` Zeile 30 trifft `PCF-E8025` nicht); die
      drei READMEs (Englisch, ohne interne Kennung) beschreiben die Eigenschaft, den
      Draht-Ursprung, „leer ohne Code“ und die Stream-Grenze; das Benutzerhandbuch nennt
      die SDK-Eigenschaft (der Satz „die Beispiele und SDKs dieses Projekts werten es
      nicht aus“ ist berichtigt), Version 1.92 → 1.93 und eine Zeile der
      Änderungshistorie in Betreibersicht ohne Kennung
      (`.harness/skills/nutzerdoku-schreiben.md`). *Zu belegen durch:*
      `make sdk-public-doc-check`, `make docs-check`, Suchlauf in §3.
- [ ] Diagnose-`error_code` (Teil von (A)–(C)): je Sprache ein Test, dass eine
      `HeartbeatStatus` mit gesetztem und mit leerem `error_code` über den
      Administrations-Client unverändert ankommt (`ErrorCode`/`errorCode`/`error_code`),
      und ein README-Satz je Sprache. Kein Mapper-Code. *Zu belegen durch:* die
      Testzeilen im `--no-cache`-Bau je Sprache.
- [ ] Mutationsbeleg je Sprache (Eingabe wie Verdrahtung, `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`):
      (M1) die Zuweisung von `code` im HTTP-Fehlerbau entfernt → Test rot; (M2) die
      `domain`-Prüfung des Statusdetail-Parsers entfernt → Test rot; (M3) die Zuweisung im
      SSE-Fehlerbau entfernt → Test rot; (M4) der gRPC-Mapper übergibt den Code nicht →
      Test rot. Je Mutation auf einer Kopie im Scratchpad (`AGENTS.md` §3.1). Der Bericht
      des Implementers nennt je Mutation die **Stelle**, die **Instanz** (Go-/C#-/Kotlin-/
      Python-Test, Parser oder Fehlerbau) und die **gesehene Farbe**; jede Verallgemeinerung
      darüber hinaus (von einer Stelle auf alle, von einer Sprache auf drei) steht als
      *hergeleitet*, „der Implementer fährt sie“ ist eine Erwartung, keine Erprobung
      (`AGENTS.md` §3.12, Verfasser-Regel). *Zu belegen durch:* Berichte von Implementer,
      Reviewer und Verifier.
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und gesondert
      ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor, kein offenes
      HIGH/MEDIUM (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff);
      `make suchlauf-nachmessen PLAN=docs/plan/planning/<Verzeichnis>/slice-sdk-meldungscodes-in-fehlertypen.md`
      endet mit Exit 0 ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Doku-Update: siehe (A)–(C) und der Nachzug-Punkt; kein SPEC-/ARC-Eintrag. Die
      Package-Versionen bleiben (Release ist Folgeschritt, §1).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls eine
      Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo ohne
      Welle für diesen Slice hier geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `sdks/csharp/PgChangeFeed.Client/Http/PgChangeFeedException.cs`, `Http/Models/ErrorResponse.cs`, `Http/PgChangeFeedHttpClient.cs` (`BuildException`, `ExtractErrorMessage`), `Sse/PgChangeFeedSseClient.cs` (eigenes `BuildException`) | update | `MessageCode` an der Basis, je Typ (Basis eingeschlossen) ein zweiter öffentlicher Konstruktor mit `string? messageCode` hinten, der bisherige bleibt; Parser liest `code` nur als JSON-String. HTTP und SSE haben je eine eigene Kopie des Fehlerbaus (gemessen, `git grep BuildException`), beide ziehen nach. |
| `sdks/csharp/PgChangeFeed.Client/Grpc/PgChangeFeedGrpcException.cs`, `Grpc/PgChangeFeedAdministrationClient.cs` (`MapException`) + neue interne Hilfsklasse für das Statusdetail | update / neu | `MessageCode` an der Basis und den sechs Typen; Parser liest `grpc-status-details-bin` aus `RpcException.Trailers`. |
| `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/io/github/pt9912/pgchangefeed/http/PgChangeFeedException.kt`, `http/model/ErrorResponse.kt`, `http/PgChangeFeedHttpClient.kt`, `sse/PgChangeFeedSseClient.kt`, `grpc/PgChangeFeedGrpcException.kt`, `grpc/PgChangeFeedAdministrationClient.kt` (`mapException`) + interne Hilfsfunktion für das Statusdetail | update / neu | wie C#; Gson koerziert eine Zahl zu einem String, deshalb liest der Parser das Feld über `JsonElement` und nimmt nur einen JSON-String (Eingabetabelle Zeile 5). |
| `sdks/python/pgchangefeed/src/pgchangefeed/exceptions.py`, `http_client.py` (`_build_error`, `_extract_error_message`; `sse_client.py` ruft dieselbe Funktion), `administration_client.py` (`_map_error`) + interne Hilfsfunktion für das Statusdetail | update / neu | Keyword `message_code=None` an beiden Basen und allen Unterklassen (sie erben `__init__`); `str(body["code"])` würde eine Zahl koerzieren — der Parser prüft `isinstance(..., str)`. |
| Testquellen der drei Packages (HTTP-, SSE-, Administrations-Fehlerabbildung; `FakeUnaryCallInvoker.cs` um Trailer erweitert, Kotlin-Transport-Fake und Python-Fake-`RpcError` um Metadaten) | update / neu | je Sprache dieselben Tabellen (unten), dazu Diagnose-`error_code`. |
| `sdks/csharp/PgChangeFeed.Client.Integration/` (HTTP- und gRPC-Phase), `sdks/kotlin/pgchangefeed-kotlin/src/integrationTest/`, `sdks/python/pgchangefeed/integration/` (je Fläche eine Testdatei) | update | Realserver-Fall je Sprache und Zugriffsweg (unten); ein Fall am echten Wire, weil ein Test gegen Fake-Transport das Zusammenspiel von Server-Fehlerkörper und SDK-Parser nicht belegt. Der gRPC-Administrations-Client hat bislang keinen Realserver-Test (gemessen: `git grep -n AdministrationClient 81d0fa96 -- sdks/csharp/PgChangeFeed.Client.Integration sdks/python/pgchangefeed/integration sdks/kotlin/pgchangefeed-kotlin/src/integrationTest` druckt 0 Zeilen). |
| `sdks/csharp/README.md`, `sdks/kotlin/pgchangefeed-kotlin/README.md`, `sdks/python/README.md` | update | Fehlerbehandlung je Sprache, kein „Upgrading“-Eintrag (die Änderung ist additiv, `ADR-0145` Festlegung 4); Englisch, nur `PCF-…`-Beispiele. |
| `docs/user/benutzerhandbuch.md` (Fehlerantworten, Fehlerform der gRPC-API, Version, Änderungshistorie) | update | Nachzug; Satz bei „Fehlerform“ berichtigt. |

**Ansatz.**

*Quellen des Codes (gemessen am Server).* HTTP: `errorResponse{Error, Code}` mit
`json:"code,omitempty"` (`internal/adapters/driving/http/middleware.go`); `401`/`403`
rufen `writeError(..., "")` und tragen kein Feld (`middleware.go`, `withToken`).
gRPC: `statusError` hängt `ErrorInfo{Reason, Domain: "pg-change-feed"}` an
`InvalidArgument`/`NotFound`/`Internal` (`internal/adapters/driving/grpc/errors.go`,
`administration.go`); die drei Auth-Statuswerte kommen als `status.Error` ohne Detail
(`interceptor.go`).

*Form der Eigenschaft (Frage A2, bestätigt in `ADR-0145` Festlegung 1).* Name `MessageCode` (C#), `messageCode` (Kotlin),
`message_code` (Python) statt `Code`/`code`: `PgChangeFeedGrpcError.code` ist in Python
bereits die `grpc.StatusCode` (`exceptions.py`, Klasse `PgChangeFeedGrpcError`), `StatusCode` in
C# und `statusCode` in Kotlin liegen daneben — ein Name `Code` wäre dort missverständlich und in Python
eine Kollision; „Meldungscode“ ist der Begriff des Handbuchs. Typ `string?`/`String?`/
`str | None`. Der Wert wird **unverändert durchgereicht**: kein Formatcheck gegen
`PCF-…` (ein künftiges Präfix bräuchte sonst ein SDK-Release), leerer String und
JSON-`null` lesen als leer. Statusdetail: der erste `Any`-Eintrag, dessen `type_url` auf `google.rpc.ErrorInfo` endet
und dessen `domain == "pg-change-feed"` ist, liefert `reason`; fremde Domäne, anderer Detailtyp, leerer
`reason`, nicht lesbare Bytes → leer, und der typisierte Fehler wird trotzdem geworfen
(der Parser wirft nie).

*API-additiv (Überladung statt optionalem Parameter, `ADR-0145` Festlegung 4).* Die
Eigenschaft ist neu, jede 0.5.x-Signatur bleibt: C# — der bisherige Konstruktor bleibt,
je Typ (die Basen eingeschlossen) kommt ein **zweiter öffentlicher Konstruktor** mit
`string? messageCode` **hinten** dazu; Kotlin — `@JvmOverloads` am Primär-Konstruktor mit
`messageCode: String? = null` hinten (erzeugt die alte JVM-Signatur weiter); Python —
keyword-only `*, message_code=None`. Bruchfläche null (quell- und binärseitig); ein
„Upgrading“-Eintrag in den READMEs entfällt. *C#-Mehrdeutigkeit:* an der Basis stehen
`(int, string, Exception)` und `(int, string, string?)` nebeneinander, ein Aufruf mit dem
`null`-Literal an der dritten Stelle ist mehrdeutig; der Code ruft die Überladung mit
benannten Argumenten (`messageCode: null`) oder gar nicht mit `null` auf. Dass der alte
Aufruf funktioniert und die Eigenschaft leer ist, belegt je Sprache ein Test (DoD), statt
es zu behaupten.

*Statusdetail ohne neue Abhängigkeit (Frage A1, entschieden in `ADR-0145` Festlegung 3).* Die Stubs kennen `google.rpc.*`
nicht (gemessen: `git grep -n -E 'ErrorInfo|grpc-status-details|google\.rpc|StatusProto|grpcio-status|CommonProtos' 81d0fa96 -- sdks`
druckt 0 Zeilen; `proto/` importiert kein `google/rpc`). Vorhandene Laufzeit-Abhängigkeiten:
C# `Google.Protobuf` 3.36.2 und `Grpc.Net.Client` 2.83.0
(`sdks/csharp/Directory.Packages.props`); Kotlin `io.grpc:grpc-protobuf` (BOM 1.84.0) und
`protobuf-java` 4.36.2 (`build.gradle.kts`); Python `protobuf>=6` und `grpcio>=1.75`
(`pyproject.toml`), **ohne** `grpcio-status`. Der Plan liest `grpc-status-details-bin`
(Binär-Trailer, Inhalt `google.rpc.Status`: Feld 3 wiederholt `Any`; `Any`: Feld 1
`type_url`, Feld 2 `value`; `ErrorInfo`: Feld 1 `reason`, Feld 2 `domain`) mit den
vorhandenen Protobuf-Laufzeiten: C# `CodedInputStream` und `Any` aus `Google.Protobuf`,
Kotlin `CodedInputStream`/`Any` aus `protobuf-java`, Python mit `google.protobuf.any_pb2.Any`
und einem eigenen Leser für Varint und längenbegrenzte Felder (`protobuf` hat für
`Status`/`ErrorInfo` ohne `googleapis-common-protos` keine Klassen). Die Feldnummern sind
die der Googleapis-Definition, *gelesen* vom Architect am 2026-10-03 in einer lokalen
Kopie der googleapis-Protos eines anderen Projekts (Version der Kopie nicht festgestellt;
Ursprung: Architect-Verdikt §1, `ADR-0145` Kontext 4); der Server-Wire wird im
Realserver-Fall je Sprache gemessen. Ein `Any` wertet der Parser nur aus, wenn sein
`type_url` auf `google.rpc.ErrorInfo` endet (deckt Tabellenzeile 8, `RetryInfo`).
Abgelehnt (entschieden, `ADR-0145` Alternative D): C# `Google.Api.CommonProtos`,
Python `grpcio-status` samt `googleapis-common-protos` (je eine neue Laufzeit-Abhängigkeit
im veröffentlichten Package für drei Felder). Kotlin: ob `com.google.rpc.Status` über
`grpc-protobuf` transitiv bereits im Klassenpfad liegt (`io.grpc.protobuf.StatusProto`),
ist **nicht gemessen** — der Implementer prüft es mit `./gradlew dependencies` in der
`build`-Stufe; selbst bei einem Treffer bleibt der Plan bei der Eigenlesung in allen
drei Sprachen, damit die drei Parser dieselbe Tabelle prüfen.

*Eingabetabelle HTTP-Fehlerkörper (HTTP-Client, SSE-Client beim Öffnen; je Sprache
identisch).* Spalten: Status · Körper · erwarteter Typ · `MessageCode` · Fehlertext.

| # | Status | Körper | Typ | Code | Text |
|---|---|---|---|---|---|
| 1 | 400 | `{"error":"x","code":"PCF-E8051"}` | BadRequest | `PCF-E8051` | `x` |
| 2 | 400 | `{"error":"x"}` | BadRequest | leer | `x` |
| 3 | 400 | `{"error":"x","code":""}` | BadRequest | leer | `x` |
| 4 | 400 | `{"error":"x","code":null}` | BadRequest | leer | `x` |
| 5 | 400 | `{"error":"x","code":5}` | BadRequest | leer | `x` |
| 6 | 404 | `{"error":"x","code":"PCF-E8025"}` | NotFound | `PCF-E8025` | `x` |
| 7 | 500 | `{"error":"x","code":"PCF-E7000"}` | ServerError | `PCF-E7000` | `x` |
| 8 | 503 | `{"error":"x","code":"PCF-E2001"}` | UnexpectedStatus | `PCF-E2001` | `x` |
| 9 | 401 | `{"error":"x"}` | Unauthorized | leer | `x` |
| 10 | 403 | `{"error":"x"}` | Forbidden | leer | `x` |
| 11 | 400 | `{"error":"x","code":"NO-PCF","extra":1}` | BadRequest | `NO-PCF` (durchgereicht) | `x` |
| 12 | 502 | `Bad Gateway` (kein JSON) | UnexpectedStatus | leer | `Bad Gateway` |

Zeile 5 trägt die Randfall-Divergenz (Gson und `str()` koerzieren die Zahl, `System.Text.Json`
wirft): Soll ist in allen drei Sprachen „leer, Text bleibt `x`“. Zeilen 9 und 10 sind der
Auth-Zweig; ein Server, der dort künftig einen Code sendet, wird durchgereicht (Zeile 11
belegt das Durchreichen ohne Formatprüfung).

*Eingabetabelle gRPC-Statusdetail (Administrations-Client; je Sprache identisch).*
Spalten: Status · Detail · erwarteter Typ · `MessageCode`.

| # | Status | `grpc-status-details-bin` | Typ | Code |
|---|---|---|---|---|
| 1 | InvalidArgument | `Status{ Any(ErrorInfo{reason "PCF-E8051", domain "pg-change-feed"}) }` | InvalidArgument | `PCF-E8051` |
| 2 | NotFound | wie 1 mit `PCF-E8025` | NotFound | `PCF-E8025` |
| 3 | Internal | wie 1 mit `PCF-E7000` | Internal | `PCF-E7000` |
| 4 | InvalidArgument | keine Trailer | InvalidArgument | leer |
| 5 | Unauthenticated | keine Trailer | Unauthenticated | leer |
| 6 | PermissionDenied | keine Trailer | PermissionDenied | leer |
| 7 | InvalidArgument | `ErrorInfo` mit `domain "example.com"` | InvalidArgument | leer |
| 8 | InvalidArgument | `Any` mit `type_url` `…/google.rpc.RetryInfo` | InvalidArgument | leer |
| 9 | InvalidArgument | `ErrorInfo` mit leerem `reason` | InvalidArgument | leer |
| 10 | InvalidArgument | nicht lesbare Bytes (`0xFF 0xFF`) | InvalidArgument | leer, Fehler wird geworfen |
| 11 | InvalidArgument | zwei `Any`: erst fremde Domäne, dann `pg-change-feed` | InvalidArgument | der zweite |
| 12 | Unavailable | wie 1 | UnexpectedStatus | `PCF-E8051` (statusunabhängig) |

Dazu je Sprache: `Diagnose` mit gesetztem (`PCF-E4003`) und leerem `error_code` kommt
unverändert an; der `Diagnose`-Aufruf trägt keinen Fehlerzweig mehr als die Tabelle.

*Realserver-Fälle (Vorschlag: ja, je Sprache zwei).* In der HTTP-Phase und in einer neuen
Verwaltungs-Phase der gRPC-Fläche (je Sprache die vorhandene Mechanik: C# `dotnet test --filter`
je Testklasse, Kotlin `integrationTest --tests`, Python `PGCHANGEFEED_TEST_FILE`):
`EnableTable` mit dem `admin`-Token auf `public.<nicht vorhandene Tabelle>` endet am
laufenden Feed-Container mit `NotFound` und dem Code `PCF-E8025` — HTTP (`404`) und gRPC
(`NotFound`, `ErrorInfo`); Gegenprobe: derselbe Aufruf mit dem `reader`-Token endet mit
`Forbidden` bzw. `PermissionDenied` und leerem Code. Der Fall belegt den Wire des Servers
(Feldnummern des Statusdetails, `domain`) am echten Server — der Unit-Test gegen
Fake-Transport belegt ihn nicht (Lehre aus dem Realserver-Slice des SSE-Filters). Der
Runner schreibt die Zeilen in `docs/user/sdk-e2e-abdeckung.md` wie bisher (marker-gegrenzt);
fällt `PCF-E8025` am Server anders aus, ist das ein Server-Fund (melden, §4 Rückführung),
und der Fall nimmt den am Server gemessenen Code. *Erwartung, nicht gemessen:* `EnableTable`
auf eine fehlende Tabelle liefert `404`/`NotFound` mit `PCF-E8025`
(Handbuch §„Fehlerantworten“ nennt `PCF-E8025` für „eine an der Quelle fehlende Tabelle“).

*Lerneintrag aus den Meldungscodes-Slices.* Je Zweig ein Test: Tabelle und Mutationen
oben decken Code gesetzt, fehlend, leer, falscher Typ, jeden Fehlertyp, SSE-Zweig (eigene
Kopie in C# und Kotlin) und Auth-Zweig.

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „die Fehlertypen der SDK-Packages
tragen keinen Meldungscode; gRPC-Statusdetail und `error_code` werden von den SDKs nicht
gelesen oder beschrieben“; Stand der Messung ist `de24629d` (Plan-Nachzug nach dem
Architect-Verdikt; **gemessen** am 2026-10-03 mit `make suchlauf-nachmessen`; die
Zahlen sind am Stand `81d0fa96` und am Stand `de24629d` gleich, beide Stände gemessen,
zuvor vom Architect nicht nachgemessen; der Implementer ergänzt die `diff`-Zeilen und
trägt Gefundenes und Nichtgefundenes ein):**

```suchlauf
de24629d 47 -n -E 'class PgChangeFeed[A-Za-z]*(Exception|Error)\b' -- sdks ':!*Test*' ':!*test*' ':!*/obj/*' ':!dist'
de24629d 28 -n -E 'BuildException|_build_error|mapException|_map_error|MapException|ExtractErrorMessage|_extract_error_message|ErrorResponse' -- sdks ':!*Test*' ':!*test*' ':!*/obj/*' ':!dist'
de24629d 13 -n -i -E 'error handling|^## Upgrading|raw .?grpc\.RpcError|StatusException' -- sdks/csharp/README.md sdks/python/README.md sdks/kotlin/pgchangefeed-kotlin/README.md
de24629d 0 -n -E 'PCF-' -- sdks
de24629d 0 -n -E 'error_code|ErrorCode' -- sdks ':!*/obj/*' ':!dist'
de24629d 0 -n -E 'ErrorInfo|grpc-status-details|google\.rpc|StatusProto|grpcio-status|CommonProtos' -- sdks ':!*/obj/*' ':!dist'
de24629d 0 -n -E 'ErrorInfo|status-details' -- examples
de24629d 1 -n -E 'SDKs dieses Projekts' -- docs/user
de24629d 0 -n -i -E 'PgChangeFeed[A-Za-z]*(Exception|Error)|Fehlertyp' -- harness/README.md spec
```

| Träger | Messung am Stand (`de24629d`, gemessen am 2026-10-03; Zeile 9 ergänzt) | Behandlung (Befund am Diff trägt der Implementer ein) |
|---|---|---|
| Fehlertyp-Klassen der Packages | Zeile 1: 47 Zeilen (C#: 7 gRPC + 8 HTTP + 1 NATS; Kotlin: 7 + 8 + 1; Python: 15) | die Basen und die HTTP- und gRPC-Typen tragen die Eigenschaft; die NATS-`MalformedMessage`-Typen nicht (Abgrenzung §1). |
| Orte, die einen Fehler bauen oder abbilden | Zeile 2: 28 Zeilen (C#: HTTP- und SSE-Bau, `MapException`; Kotlin: HTTP, SSE, `mapException`; Python: `_build_error` (SSE ruft sie), `_map_error`) | jeder Ort setzt den Code oder die Auslassung ist begründet (Stream-Clients roh, §1). |
| READMEs: Fehlerbehandlung | Zeile 3: 13 Zeilen | je Sprache der Abschnitt „Error handling“ und die gRPC-Fehlertabelle nennen die Eigenschaft; die Zeilen „… as `grpc.RpcError` while iterating“ und der StatusException-Satz bleiben als Stream-Grenze stehen. |
| `PCF-` unter `sdks/` | Zeile 4: 0 | wächst um die README-Beispiele und Tests; `make sdk-public-doc-check` Exit 0 ist der Wächter. |
| `error_code` in den SDKs | Zeile 5: 0 | nur Tests und README-Sätze (die Stubs tragen das Feld); kein Mapper-Code. |
| Statusdetail-Leser | Zeile 6: 0 | je Sprache genau eine interne Hilfsfunktion. |
| Beispiele | Zeile 7: 0 | bleibt 0 (Abgrenzung §1). |
| Handbuch: „SDKs werten es nicht aus“ | Zeile 8: 1 (Zeile 1938, Absatz „Fehlerform“) | der Satz ist berichtigt; die Historienzeilen 1.33/1.35 („Fehlerklasse“) bleiben Historie. Nicht gesucht und nicht gefunden: weitere Träger in `spec/` und `harness/README.md` (Zeile 9: 0 Zeilen, gemessen; die Muster `Exception|Error` allein treffen dort nur fremde Zeilen). Das Muster von `tools/harness/sdk-public-doc-check.sh` (Zeile 30) trifft `PCF-E8025` nicht (gemessen mit `grep -c -E` auf diese Zeichenkette: 0). |

## 4. Trigger

**Start** (`next` → `in-progress`): kein anderer Slice liegt in `in-progress/` (WIP-Limit 1);
das Server-Release 0.6.0 trägt die Codes (Commit `81d0fa96`); Docker mit Netzzugang für
die drei `make sdk-pack-*`-Ziele und `make examples-*`, ein geladenes `:dev`-Image
(`make image`) für die drei `make test-sdk-*-integration`-Ziele.

**Architect-Fragen — entschieden** in
[`ADR-0145`](../../adr/0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md) und im
[Architect-Verdikt](../../../reviews/architect-verdict-sdk-meldungscodes-in-fehlertypen.md)
(Startbedingung „Architect-Entscheid“ erfüllt):

- **A1 — neue Laufzeit-Abhängigkeit für das Statusdetail?** Entschieden: nein, Eigenlesung
  mit der vorhandenen Protobuf-Laufzeit in allen drei Sprachen (§3 Ansatz; `ADR-0145`
  Festlegung 3, Verdikt §1).
- **A2 — öffentliche API-Form.** Entschieden: `MessageCode`/`messageCode`/`message_code`,
  leer als `null`/`None`, Wert unverändert durchgereicht (kein Formatcheck),
  **Überladung statt optionalem Parameter** (C# zweiter Konstruktor, Kotlin `@JvmOverloads`,
  Python keyword-only) — weicht vom ursprünglichen Plan-Vorschlag ab (`ADR-0145`
  Festlegung 1, 2 und 4, Verdikt §2).
- **A3 — Umfang der gRPC-Fehler.** Entschieden: nur der Administrations-Client (§1);
  die Stream-Clients behalten rohe Fehler (`ADR-0145` Festlegung 5, Verdikt §3).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): braucht eine Sprache mehr als
  drei Fixrunden, teilt sich der Slice an der Sprach-Wurzel (`sdks/<sprache>/`).
- `in-progress` → `open` (blockiert): der Server sendet am Wire einen anderen Körper oder
  ein anderes Statusdetail als in `ADR-0144` und im Handbuch beschrieben (gemessen im
  Realserver-Fall) — der Fund geht an den Architect, kein Umgehen im SDK.

## 5. Closure-Trigger

DoD vollständig, Review ohne offenes HIGH/MEDIUM, `make gates` grün (Exit-Code
ungefiltert), die drei `make sdk-pack-*`-Ziele und `make examples-csharp`/
`make examples-kotlin` mit ausgeführten Tests belegt (Bau ohne Cache der Test-Stufe oder
rote Mutation), die drei `make test-sdk-*-integration`-Ziele grün (die gedruckten Zeilen
der neuen Fälle im Bericht), `make sdk-public-doc-check` Exit 0, Suchlauf-Block
nachgemessen, Closure-Notiz mit Lerneintrag geschrieben. **Folgeschritt nach der
Verifikation (nicht Teil des Slice, Hauptlauf):** Release-Notiz je Sprache nennt die neue
Eigenschaft, ohne Neukompilierungs-Hinweis (`ADR-0145` Festlegung 4); Versionsänderung auf 0.6.0 je Sprache
(`PgChangeFeed.Client.csproj`, `build.gradle.kts` an zwei Stellen, `pyproject.toml`),
Tags `sdk-csharp-v0.6.0`/`sdk-python-v0.6.0`/`sdk-kotlin-v0.6.0`, Beobachtung der drei
Release-Workflows (`AGENTS.md` §3.10, `docs/maintainer/releasing.md` §4).

## 6. Risiken und offene Punkte

- **Drei Sprachen, drei Parser, drei Lesarten** (Zahl statt String im Feld `code`; leer
  gegen fehlend; Statusdetail mit fremder Domäne). — **Ausgang:** weiter offen bis zur
  Closure: Dieselbe Tabelle je Sprache, Zeile 5 und die gRPC-Zeilen 7 bis 11 binden die
  Randfälle (`BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall`).
- **Binärkompatibilität.** — **Ausgang:** entfallen: `ADR-0145` Festlegung 4 (Überladung
  statt optionalem Parameter) erhält jede 0.5.x-Signatur; die Verdikt-Aussage ist
  *hergeleitet*, an einem Gast-Assembly nicht gemessen. Belegt wird der Erhalt durch je einen
  Test des alten Konstruktors bzw. Aufrufs (DoD).
- **Neue Abhängigkeit.** — **Ausgang:** entfallen: `ADR-0145` Festlegung 3 (Eigenlesung,
  keine neue Abhängigkeit); der Diff ändert `Directory.Packages.props`/`.csproj` (Pakete),
  `build.gradle.kts` (Abhängigkeiten) und `pyproject.toml` (`dependencies`) nicht
  (Review-Prüfung).
- **Handgeschriebener Wire-Parser.** Die Feldnummern von `Status`/`Any`/`ErrorInfo` sind
  *gelesen* (Architect, lokale Kopie der googleapis-Protos eines anderen Projekts, Version
  der Kopie nicht festgestellt), nicht aus einer Datei dieses Repos. —
  **Ausgang:** entfallen mit dem Realserver-Fall je Sprache, der den Server-Wire liest;
  die Unit-Tests bauen das Detail mit der Protobuf-Laufzeit der Sprache und nicht mit dem
  Parser selbst nach, damit ein Feldnummernfehler auf beiden Seiten nicht still grün ist.
- **Server alt/neu.** SDK 0.6.0 gegen einen Server vor 0.6.0 (kein `code`, kein
  `ErrorInfo`). — **Ausgang:** entfallen mit den Tabellenzeilen 2 und 4 (HTTP) und 4 (gRPC):
  ohne Code auf dem Draht liest die Eigenschaft leer, Text und Typ bleiben; am echten
  alten Server **nicht gelaufen** (hergeleitet aus der Tabelle).
- **Test läuft nicht, Bau grün.** Ein Lauf aus dem Docker-Schicht-Cache druckt keine
  Testzeile (`BEO-PGC/docker-cache-ueberspringt-tests-still`). — **Ausgang:** weiter offen
  bis zur Closure: Belegbefehl ist `docker build --no-cache --target build` je Sprache oder
  eine rote Mutation (DoD, Mutationsbeleg).
- **Interne Kennung in öffentlichem Text.** README, Docstrings, KDoc, XML-Doku und
  Fehlertexte erreichen die Anwender über die Packages; `PCF-…` ist erlaubt, `ADR-0144`
  nicht. — **Ausgang:** entfallen mit `make sdk-public-doc-check` Exit 0 und
  Suchlauf-Zeile 4 am Diff.
- **Realserver-Fall deckt den Fehlerweg des Streams nicht.** Die Stream-Clients bleiben roh
  (§1). — **Ausgang:** weiter offen, benannte Grenze: README sagt es; eine Typisierung wäre
  ein eigener Slice mit Architect-Entscheid (A3).

## 7. Closure-Notiz

Wird bei der Closure gefüllt (vor dem `git mv` nach `done/`); Ursprung der Angaben
(gemessen / übernommen) je Zeile, `AGENTS.md` §3.12.

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag:** —
- **Beobachtungs-Register (`../observations/`):** —
- **Folge-Slices:** — (SDK-Release 0.6.0 je Sprache ist ein Schritt des Hauptlaufs, kein Slice)
- **Risiken aus §6:** —
- **Drei Paarungen:** —

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area `*`
(`harness/conventions.md`, Modus-Deklaration: Greenfield, Kürzel `PGC`) mit den
Pfaden `sdks/` und `docs/user/` — eine Sub-Area.

**Vorgelagert — offene Beobachtungen sichten:** das Register wurde durchgegangen
(Verzeichnisliste `docs/plan/planning/observations/BEO-PGC/`, 2026-10-03); Treffer:
`BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall` (gestrichen, 3×; die Eingabetabelle
dieses Plans ist die Gegenmaßnahme),
`BEO-PGC/docker-cache-ueberspringt-tests-still` (offen),
`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (verkörpert),
`BEO-PGC/deutsches-fachwort-im-englischen-sdk-readme` (verkörpert: README bleibt Englisch),
`BEO-PGC/intern-kennungen-in-ausgelieferten-texten` (offen),
`BEO-PGC/nachzug-laesst-ueberholten-text-stehen` (verkörpert: Suchlauf §3),
`BEO-PGC/sdk-decoder-verhalten-am-neuen-feld-ungemessen` (gestrichen: berührt die
Toleranz der Decoder gegen ein unbekanntes Feld; Tabellenzeile 11 der HTTP-Tabelle mit
Zusatzfeld `extra` misst sie je Sprache am neuen Feld `code`). Zähler nicht neu
ausgezählt (übernommen aus den `state.md`-Dateien).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
