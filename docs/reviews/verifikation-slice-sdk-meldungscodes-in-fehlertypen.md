# Verifikations-Report: slice-sdk-meldungscodes-in-fehlertypen — 2026-10-03

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich, Entscheidungs-Konformität
([`ADR-0145`](../plan/adr/0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md),
[`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md),
[`ADR-0134`](../plan/adr/0134-sdk-public-doc-check-gate-make-gates.md), Architect-Verdikt
[`architect-verdict-sdk-meldungscodes-in-fehlertypen`](architect-verdict-sdk-meldungscodes-in-fehlertypen.md))
und Plan-vs-Code-Diff. Review-Artefakt:
[`review-slice-sdk-meldungscodes-in-fehlertypen`](review-slice-sdk-meldungscodes-in-fehlertypen.md)
(0 HIGH/0 MEDIUM, F-1 bis F-3 LOW, F-4 bis F-7 INFO). Formvorbild:
[`verifikation-slice-wal-fehlerschwelle-ausgangsklasse`](verifikation-slice-wal-fehlerschwelle-ausgangsklasse.md).

**Gegenstand:** Plan
[`slice-sdk-meldungscodes-in-fehlertypen`](../plan/planning/in-progress/slice-sdk-meldungscodes-in-fehlertypen.md)
(wellenlos), Diff-Range `85f0506c..HEAD` (`6a589a19`): sechs Commits — `20ff9fe0` (C#), `5d8fe1af`
(Kotlin), `7efea1fd` (Python), `a1947e8a` (Docs/Runner/Handbuch), Review `3ab233e5`, Fixrunde
`6a589a19`. Die Fixrunde ändert Produktionslogik der drei Fehlerkörper-Parser und ist von mir voll
mitgeprüft (kein weiterer Review-Agent). Dieser Lauf ändert weder Code noch Plan noch Doku; er schreibt
nur diesen Report. Alle Mutationen liefen auf Kopien im Scratchpad (`git archive HEAD sdks proto`),
erzeugt per `sed … Original > Kopie` (Ausgabe nach stdout, kein `-i`), Bau der Test-Stufe per
`docker build --target build` aus der Kopie (Form von `make sdk-pack-*`), Rücknahme per `cp`.
`git status --short` im Echtrepo war am Ende leer; ich habe nicht gepusht.

**Verdikt: bestanden, mit Bedingungen (§7).** Kein HIGH, kein MEDIUM. Zwei LOW (V-1, V-2), drei INFO.
Alle Sensoren, die ich gefahren habe, sind grün; alle Mutationen färben sich rot; die Entscheidung zu F-2
ist in den drei Sprachen gleich umgesetzt.

## 1. Eigene Sensor-Belege (dieser Lauf, Exit-Code je Lauf direkt in eine Log-Datei, nie durch eine Pipe — [`AGENTS.md`](../../AGENTS.md) §3.9)

| Sensor | Exit | Beleg aus meinem Lauf |
|---|---|---|
| `make sdk-pack-csharp` / `-python` / `-kotlin` | 0 / 0 / 0 | alle drei aus dem Docker-Schicht-Cache (Quelle seit dem Implementer-Lauf unverändert), deshalb keine Testzeile — Ersatzbeleg unten |
| `docker build --no-cache --target build` der unmutierten Kopie (Bau-Kontext `proto` wie das Make-Ziel) | 0 / 0 / 0 | Python `195 passed, 6 warnings`; C# `Passed!  - Failed: 0, Passed: 201, Skipped: 0, Total: 201`; Kotlin `> Task :test` ausgeführt (nicht `UP-TO-DATE`), `BUILD SUCCESSFUL` (Gradle druckt keine Anzahl; die Wirkung belegen die Mutationen in §4) |
| `make examples-csharp` / `make examples-kotlin` | 0 / 0 | Bau bis `…:csharp-nats-stream done` bzw. `…:kotlin-nats-stream done`; Schichten aus dem Cache (`DONE 0.0s`) — die Beispiele berühren `sdks/` nicht (Plan §1, `git grep ErrorInfo -- examples` = 0), die Läufe belegen nur, dass sie weiter bauen |
| `make image` | 0 | `:dev`-Image gebaut (Server unverändert) |
| `make test-sdk-csharp-integration` | 0 | `run-sdk-csharp-integration-tests: Fehlercode-Belege (ADR-0145) grün — … HTTP (404) und gRPC (NotFound) mit dem Meldungscode PCF-E8025 als MessageCode des SDK-Fehlertyps, ein Reader-Token endete mit 403 bzw. PermissionDenied ohne Code` |
| `make test-sdk-python-integration` | 0 | dieselbe Zeile mit `(NOT_FOUND)`, `message_code`, `PCF-E8025` |
| `make test-sdk-kotlin-integration` | 0 | dieselbe Zeile mit `(NOT_FOUND)`, `messageCode`, `PCF-E8025` |
| `git status --short` nach den drei Realserver-Läufen | leer | `docs/user/sdk-e2e-abdeckung.md` (vom Runner geschrieben) ist byte-gleich dem committeten Stand |
| `make sdk-public-doc-check` / `make handbuch-public-doc-check` | 0 / 0 | `keine interne Kennung unter sdks` / `keine interne Kennung in 3 Nutzerdokumenten unter docs/user/` |
| `make docs-check` | 0 | `d-check: 1618 Datei(en) geprüft, 0 Befund(e)` (vor dem Report; Nachlauf mit dem Report in §8) |
| `make gates` | 0 | `baseline-verify: v6.13.0 OK — 54 Dateien`, `coverage-gate: OK — Coverage 82.40% erfüllt Schwelle 80%`, `generated-sync: OK`, `sdk-public-doc-check`, `handbuch-public-doc-check`, `ausgabe-kennungen-check`, `meldungscodes-check: 97 Codes in Tabelle und Katalog gleich`, `a-check gesamt: 0 Befund(e)` |
| `make doc-immutable RANGE=85f0506c..HEAD` | 0 | `d-check: 1618 Datei(en) geprüft, 0 Befund(e)` |
| `make doc-commits RANGE=85f0506c..HEAD` | 0 | `0 Befund(e)` |
| `make commit-traceability RANGE=85f0506c..HEAD` | 0 | `OK — 6 Commit(s) in "85f0506c..HEAD", Betreffs ohne Struktur-ID` (Standing-Gate über `HEAD~5..HEAD` in `make gates` ebenfalls 0) |
| `make fmt-check` | 0 | `334 Go-Dateien geprüft, alle formatiert` (Diff ohne Go-Dateien; Lauf der Vollständigkeit halber) |
| `make kommentar-kennungen DIFF=283d6175` | 0 | kein Kandidat |
| `make suchlauf-nachmessen PLAN=…slice-sdk-meldungscodes-in-fehlertypen.md` | 0 | `suchlauf-nachmessen: 21 Zeilen stimmen` |
| `make doc-trace` | 0 | `80 Anforderung(en), 0 Waise(n).` |

Nicht gelaufen: kein Lauf musste aus Zeit oder Netz ausfallen. Nicht gefahren: ein Gast-Assembly gegen die
neue C#-DLL (Binärkompatibilität bleibt *hergeleitet*, Plan §6) und ein SDK gegen einen Server vor 0.6.0.

## 2. DoD — Verdikt je Zeile (§2 des Plans)

Gezählt am Plan: 9 `[x]`-Zeilen (A, B, C, Nachzug zum Public-Doc-Gate, Diagnose-`error_code`, `make gates`, Review,
§3.13-Suchlauf, Doku-Update), 7 `[ ]`-Zeilen (Mutationsbeleg, Closure-Notiz, Reconciliation-entfällt,
Beobachtungs-Register, Risiko-Ausgänge, drei Paarungen — Slice liegt in `in-progress/`). Ich setze keinen
Haken (Planner-Vorrecht).

| # | DoD-Zeile | Verdikt | Beleg |
|---|---|---|---|
| 1 | (A) C#: `MessageCode`, Tabellen, Realserver, README, alter Konstruktor | **bestätigt** | Code gelesen (`ErrorBody.Parse`, `StatusDetail.ReadMessageCode`, je Typ zweiter Konstruktor, `HttpException_ConstructedWithoutMessageCode_…`, `GrpcException_ConstructedWithoutMessageCode_…`); 201 Tests grün im `--no-cache`-Bau; Realserver-Zeile mit `PCF-E8025`; Rot-Belege §4 |
| 2 | (B) Kotlin: `messageCode`, `@JvmOverloads` | **bestätigt** | `parseErrorBody`, `StatusDetail.readMessageCode`, Tests `… constructed without message code keeps working and has none`; `:test` ausgeführt; Realserver-Zeile; Rot-Belege §4 |
| 3 | (C) Python: `message_code` keyword-only, `code` bleibt `grpc.StatusCode` | **bestätigt** | `exceptions.py` gelesen (`*, message_code=None`, `self.code` unverändert); 195 Tests grün (`--no-cache`); Realserver-Zeile; Rot-Belege §4 |
| 4 | `make sdk-public-doc-check` Exit 0, READMEs, Handbuch 1.93 | **bestätigt, mit V-3** | Exit 0 (eigener Lauf); Handbuch am `HEAD` 1.94 (Fixrunde), Zeilen 1.93 und 1.94 in der Historie ohne Kennung, `make handbuch-public-doc-check` Exit 0; README-Diffs der drei Packages gelesen (Eigenschaft, Draht-Ursprung, „leer ohne Code“, Stream-Grenze) |
| 5 | Diagnose-`error_code` ohne Mapper-Code | **bestätigt** | je Sprache ein Test mit `PCF-E4003` und leer (`Diagnose_ErrorCodeArrivesUnchanged`, `diagnose error code arrives unchanged`, Python-Entsprechung); der Diff hat keine Mapper-Zeile dazu (F-4 des Reviews bleibt INFO: der Test prüft den Durchreichungsweg, den die Stubs tragen) |
| 6 | Mutationsbeleg je Sprache (`[ ]`) | **teilweise von mir nachgefahren** | siehe §4: M1 in allen drei Sprachen, M2 in Kotlin, die F-2-Mutation in allen drei, die Feldnummer in C# — alle rot. M3 und M4 sowie M2 in C#/Python habe ich **nicht** gefahren (Implementer-Angabe, *übernommen*); der Haken ist Planner-Vorrecht |
| 7 | `make gates` grün | **bestätigt** | eigener Lauf Exit 0 (§1) |
| 8 | Review durchgeführt, kein HIGH/MEDIUM | **bestätigt** | Report liegt vor; F-2 und F-3 sind in `6a589a19` behoben und von mir mitgeprüft (§3, §4) |
| 9 | §3.13-Suchlauf, `make suchlauf-nachmessen` Exit 0 | **bestätigt** | eigener Lauf: `21 Zeilen stimmen`, Exit 0 |
| 10 | Doku-Update, Package-Versionen bleiben | **bestätigt** | `git diff 85f0506c HEAD --stat` über `*.csproj`, `Directory.Packages.props`, `build.gradle.kts`, `pyproject.toml`, `docs/user/version.md` ist leer; csproj `<Version>0.5.0</Version>` (Z. 21), `pyproject.toml` `version = "0.5.0"` (Z. 7), `build.gradle.kts` `0.5.0` an Z. 86 und Z. 210, `docs/user/version.md` `0.6.0` (Server) |
| 11 bis 16 | Closure-Notiz, Reconciliation, Beobachtungs-Register, Risiko-Ausgänge, Paarungen (`[ ]`) | **korrekt offen** | Plan §7 trägt „—“; Vorschläge in §6 |

**Ein `[x]` trägt keinen Beleg mehr:** keiner.

## 3. Entscheidungs-Konformität und Plan-vs-Code-Diff

**Entscheidung zu Review F-2 (Plan §3, „Entscheidung des Implementers zu Zeile 13“).** Die Regel — jedes
der Felder `error`/`code` zählt nur als JSON-String, unabhängig vom anderen; ein nicht-String `error` liefert den
rohen Körper als Text — ist schlüssig und deckt sich mit Zeile 12 der Tabelle (kein JSON → roher Körper). Gelesen,
in allen drei Sprachen gleich umgesetzt:

| Sprache | Stelle | Lesung |
|---|---|---|
| C# | `Http/ErrorBody.cs` `Parse`, `Http/Models/ErrorResponse.cs` | `Error` und `Code` je `JsonElement`; `ValueKind == String` je Feld; `text ?? body`; `JsonException` und JSON-`null` → `(body, null)`; HTTP- und SSE-Client rufen dieselbe Funktion |
| Kotlin | `http/ErrorBody.kt` `parseErrorBody` | `JsonParser.parseString`, Objekt-Prüfung, `stringOf` nur für `isString`; `ifEmpty { null }`; `JsonParseException` → roher Körper; HTTP- und SSE-Client rufen dieselbe Funktion; `ErrorResponse.kt` entfällt |
| Python | `http_client.py` `_extract_error` | `isinstance(code, str) and code != ""`; `isinstance(text, str)` sonst `response.text`; `ValueError` und Nicht-`dict` → `(response.text, None)`; `sse_client.py` ruft `_build_error` |

Tabellen im Test, gezählt: **HTTP 13 Zeilen in allen drei Sprachen** (`HttpRows` C#, `httpRows` Kotlin, `_HTTP_ROWS` Python;
Zeile 13 `{"error":5,"code":"PCF-E8051"}` überall mit dem rohen Körper als erwartetem Text; die Tabelle läuft in C#, Kotlin
und Python je durch den HTTP- **und** den SSE-Test). **gRPC: 12 Zeilen in C# und Kotlin**, in Python 13 Einträge — die zwölf des
Plans plus ein Eintrag mit `None` statt leerem Trailer-Tupel (Zeilen 4 und 4′), Python-spezifisch und eine Obermenge, kein Widerspruch.
Die Statusdetails der Tests sind mit der Protobuf-Laufzeit der Sprache gebaut, nicht mit dem Parser unter Test.

**Weitere Konformitäts-Punkte.** (a) Keine neue Abhängigkeit (leere Diffs der vier Abhängigkeitsdateien, s. o.); die Statusdetail-Leser
arbeiten mit `Google.Protobuf`/`CodedInputStream`, `protobuf-java`, `google.protobuf.any_pb2` und eigenem Varint-Leser
([`ADR-0145`](../plan/adr/0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md) Festlegung 3). (b) Der Wert wird unverändert durchgereicht, kein Formatcheck (Tabellenzeile 11 `NO-PCF` belegt es in allen
drei Sprachen). (c) Überladung statt optionalem Parameter (Festlegung 4): C# zweiter Konstruktor je Typ, Kotlin `@JvmOverloads`,
Python keyword-only; je ein Test des alten Aufrufs. (d) Stream-Clients und NATS bleiben roh (Abgrenzung §1, `git diff` berührt sie nicht).
(e) Auth-Zweig: Tabellenzeilen 9/10 und gRPC 5/6 liefern leer.

| Plan-Zeile (§3) | Ist im Diff (`git diff 85f0506c HEAD --stat`, 40 Dateien) |
|---|---|
| C# `PgChangeFeedException.cs`, `ErrorResponse.cs`, `PgChangeFeedHttpClient.cs`, `PgChangeFeedSseClient.cs` | `M`; `Http/ErrorBody.cs` neu (Plan-Nachzug, benannt) |
| C# `PgChangeFeedGrpcException.cs`, `PgChangeFeedAdministrationClient.cs` + Hilfsklasse | `M`; `Grpc/StatusDetail.cs` neu (Plan-Nachzug) |
| Kotlin Fehlertypen, HTTP-, SSE-, Verwaltungs-Client + Hilfsfunktion | `M`; `http/ErrorBody.kt`, `grpc/StatusDetail.kt` neu, `http/model/ErrorResponse.kt` gelöscht (Fixrunde, benannt) |
| Python `exceptions.py`, `http_client.py`, `administration_client.py` + Hilfsfunktion | `M`; `_status_detail.py` neu (Plan-Nachzug) |
| Testquellen, Test-Hilfen mit Trailern | `MessageCodeTests.cs`, `MessageCodeTest.kt`, `test_message_code.py` neu; `FakeUnaryCallInvoker.cs`, `FakeAdministrationTransport.kt` `M` |
| Realserver-Fälle (Abweichung: eine Fehlercode-Phase je Runner) | `ErrorCodeRealserverTests.cs`, `ErrorCodeRealserverTest.kt`, `test_error_code_realserver.py` neu; die Abweichung vom Vorschlag „HTTP-Phase plus Verwaltungs-Phase“ ist in Plan §3 als Abweichung benannt und schlüssig: eine Phase trägt beide Wege und die Gegenprobe, der Runner-Mechanismus (`READY`/`RECEIVED`/`REJECTED`) bleibt unverändert |
| READMEs, Handbuch | drei READMEs, `benutzerhandbuch.md` (1.92 → 1.94) `M` |
| Runner, `sdk.mk`, `harness/README.md`, Abdeckungs-Träger | drei Runner (`sql_kind none`, `captured=1`), `sdk.mk` (dreizehn → vierzehn Phasen), `harness/README.md` (drei Integrationszeilen), `sdk-e2e-abdeckung.md` (+ je eine Zeile) `M` |

Kein unbenannter Nebeneffekt: `Directory.Packages.props`, `.csproj`, `build.gradle.kts`, `pyproject.toml`, `spec/**`, `examples/**` und der Go-Baum sind nicht im Diff.

**Plan-Abweichung „eine Fehlercode-Phase pro Runner“, `sql_kind none`, `captured=1`:** schlüssig. Die Phase hat keine Zeile in `cdc.changes`; ihre
Identität ist der Meldungscode am Wire beider Zugriffswege, und die Tests halten ihn seit der Fixrunde exakt auf `PCF-E8025` (F-3 behoben, gelesen
in `test_error_code_realserver.py` Z. 52–54, 108–110; der Runner-Regex `PCF-E[0-9]{4}` bleibt weit, aber der Test pinnt). Der Abdeckungs-Träger
`docs/user/sdk-e2e-abdeckung.md` trägt je Sprache genau eine neue Zeile, deren Testklassen-Name (`ErrorCodeRealserverTests`, `ErrorCodeRealserverTest`,
`test_error_code_realserver.py`) und Runner-Pfad zu den Dateien passen; der Runner-Lauf hat die Datei byte-gleich zurückgeschrieben.

## 4. Mutationen (eigene Läufe, Eingabeseite; Stelle · Instanz · gesehene Farbe)

Bau: `docker build --progress=plain --build-context proto=<Kopie>/proto --target build sdks/<Sprache>`; Rot = der Bau bricht an der Test-Zeile ab.

| # | Sprache | Mutation (Stelle) | Instanz | Gesehene Farbe |
|---|---|---|---|---|
| M1 | Python | `_build_error` ohne `message_code=` (`http_client.py`) | `test_http_client_error_carries_the_message_code`, `test_sse_client_…` | **rot, 12 failed, 183 passed** (HTTP und SSE, je Zeilen 1, 6, 7, 8, 11, 13) |
| M1 | C# | `PgChangeFeedHttpClient.BuildException` setzt `messageCode = null` | `HttpClient_ErrorCarriesTheMessageCode` | **rot, Failed: 6 von 201** (Zeilen 1, 6, 7, 8, 11, 13; die SSE-Fälle bleiben grün, weil nur die HTTP-Kopie mutiert ist) |
| M1 | Kotlin | `PgChangeFeedHttpClient.buildException` setzt `messageCode = null` | `http client error carries the message code` | **rot, 1 von 122 Tests** (`MessageCodeTest.kt:83`; die Schleife bricht an der ersten Zeile ab, welche Zeilen einzeln rot würden, ist nicht gemessen) |
| F-2 | Python | nicht-String `error` verwirft den Code (`_extract_error`) | HTTP- und SSE-Test, Zeile 13 | **rot, 2 failed, 193 passed** |
| F-2 | C# | `ErrorBody.Parse` liefert bei nicht-String `error` `(body, null)` | HTTP- und SSE-Test, Zeile 13 | **rot, Failed: 2 von 201** |
| F-2 | Kotlin | `parseErrorBody` liefert bei nicht-String `error` `ParsedError(body, null)` | `http …` und `sse client error carries the message code` | **rot, 2 von 122 Tests** (`MessageCodeTest.kt:83` und `:97`) |
| M2 | Kotlin | `domain`-Vergleich aus `StatusDetail.readFromAny` entfernt | `administration client error carries the message code` | **rot, 1 von 122 Tests** (`MessageCodeTest.kt:211`) |
| Feldnummer | C# | `StatusDetailsTag` Feld 3 → 4 (`StatusDetail.cs`) | `AdministrationClient_ErrorCarriesTheMessageCode` | **rot, Failed: 5 von 201** (gRPC-Zeilen 1, 2, 3, 11, 12) |

Nicht von mir gefahren und deshalb als **übernommen** (Implementer-Tabelle, Plan §3) oder **hergeleitet** gekennzeichnet: M2 in C# und Python, M3 (SSE-Fehlerbau
je Stelle, C# und Kotlin) und M4 (Mapper übergibt den Code nicht) in allen drei Sprachen, die Feldnummer in Kotlin und Python (Python ist vom Implementer gefahren,
Kotlin ist nirgends gefahren). Die Tabellen sind je Sprache gleich und decken diese Stellen **hergeleitet** gleichartig. Die Rot-Zahlen im Plan (§3 Mutationsbeleg:
Python 10, C# 5) sind Stände **vor** der Fixrunde (siehe V-2); meine Zahlen am `HEAD` sind 12 und 6.

## 5. Findings des Verifiers

| # | Kategorie | Befund |
|---|---|---|
| V-1 | LOW | **C# liest die Feldnamen des Fehlerkörpers ohne Groß-/Kleinschreibung, Kotlin und Python mit.** `JsonOptions` beider C#-Clients setzt `PropertyNameCaseInsensitive = true` (`PgChangeFeedHttpClient.cs:32`, `PgChangeFeedSseClient.cs:35`), `ErrorBody.Parse` deserialisiert `ErrorResponse` damit; ein Körper `{"Error":"x","Code":"PCF-E8051"}` liefert in C# den Code, in Kotlin (`fields.get("code")`) und Python (`body.get("code")`) nicht (hergeleitet aus dem Code, nicht gefahren). Die Tabelle hat keine Zeile dazu; der Server sendet immer klein (`json:"code"`), das Verhalten für `error` war schon vor diesem Slice so. Dieselbe Klasse wie F-2 des Reviews, zweiter Randfall der Drei-Sprachen-Kopie, ohne Wirkung am echten Server. |
| V-2 | LOW | **Stehende Zahlen im Plan sind Stände vor der Fixrunde, ohne Hinweis darauf.** Plan §3 „Belege des Implementers“ nennt `193 passed` (Python) und `Passed: 199` (C#); am `HEAD` sind es `195` und `201` (mein `--no-cache`-Bau). Die Mutationstabelle nennt „Python 10 Fälle rot“, „C# 5 Fälle rot“ für M1 (am `HEAD`: 12 und 6, wegen Zeile 13). Der Messzeitpunkt steht nicht an den Zahlen ([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz A). Dazu ein Tippfehler an Plan Z. 414/415 („Die\nDie Realserver-Tests …“). Kein Befund am Code. |
| V-3 | INFO | **Das Handbuch nennt eine Package-Version, die es noch nicht gibt.** `benutzerhandbuch.md` (1.94) sagt an zwei Stellen „ab Package-Version 0.6.0“; veröffentlicht ist 0.5.0 (`0.5.0` in `.csproj`, `pyproject.toml`, `build.gradle.kts`). Das ist mit dem geplanten Release konsistent (nächste Minor-Stufe), aber eine Aussage über ein künftiges Artefakt; trägt der Release eine andere Nummer, ist das Handbuch falsch. Bedingung B-1. |
| V-4 | INFO | **Fehlertext bei nicht-String `error` ändert sich gegenüber 0.5.0.** Python lieferte `str(body["error"])` (`"5"`, `"None"`), Kotlin koerzierte über Gson; jetzt steht in allen drei Sprachen der rohe Körper. Das Handbuch (1.94) sagt es, die READMEs nicht (`git grep -i "raw body"` über die drei READMEs: 0 Treffer). Der Server sendet solche Körper nicht; „additiv“ in `ADR-0145` Festlegung 4 gilt für die Signaturen, nicht für diesen Text. |
| V-5 | INFO | **Python-Tabelle hat einen Eintrag mehr als der Plan** (13 statt 12 gRPC-Einträge, `trailing_metadata() is None`), siehe §3 — Obermenge, kein Widerspruch. |

## 6. Vorschläge für die Closure (Planner entscheidet)

**Risiko-Ausgänge (Plan §6, je mit Beleg-Anker):**

| Risiko | Vorschlag | Beleg |
|---|---|---|
| Drei Sprachen, drei Parser, drei Lesarten | **eingetreten (F-2, behoben)**, Rest-Randfall V-1 weiter offen als benannte Grenze | Review F-2, Fixrunde `6a589a19`, Tabellenzeile 13 in drei Sprachen, Rot-Läufe §4 |
| Binärkompatibilität | **entfallen**, Erhalt *hergeleitet* (Gast-Assembly nicht gemessen) | Test des alten Konstruktors/Aufrufs je Sprache grün |
| Neue Abhängigkeit | **entfallen** | leere Diffs der vier Abhängigkeitsdateien (§2 Zeile 10) |
| Handgeschriebener Wire-Parser | **entfallen** | Realserver-Zeile je Sprache mit `PCF-E8025`, Tests pinnen den Wert; Statusdetails der Unit-Tests mit der Protobuf-Laufzeit gebaut |
| Server alt/neu | **entfallen**, *hergeleitet* (am alten Server nicht gelaufen) | Tabellenzeilen 2 und 4 (HTTP), 4 (gRPC) |
| Test läuft nicht, Bau grün | **entfallen** | `--no-cache`-Bau zeigt 195 und 201 Tests bzw. ausgeführtes `:test`; die Mutationen färben den Bau rot |
| Interne Kennung in öffentlichem Text | **entfallen** | `make sdk-public-doc-check` Exit 0, Suchlauf 21 Zeilen stimmen |
| Realserver deckt den Fehlerweg des Streams nicht | **weiter offen**, benannte Grenze (README sagt es) | Plan §1, Abgrenzung |

**Register-Fortschreibung:**
- F-1 des Reviews (Quellkompat-Mehrdeutigkeit des `null`-Literals, [`ADR-0145`](../plan/adr/0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md) Festlegung 4 „quellseitig erhalten“ breiter als ihr Beleg):
  die Klasse [`adr-aussage-breiter-als-ihre-messung`](../plan/planning/observations/BEO-PGC/adr-aussage-breiter-als-ihre-messung/state.md)
  trägt gemessen mit `ls evidence | wc -l` **11** Dateien (Zähler 11×, Ausgang **verkörpert**); der Eintrag nennt einen **Deckel bei 10×**: weitere Auftreten mit
  Schwere ≤ LOW, vor dem Merge vom Reviewer/Verifier/Architect gefunden, an einem bekannten Träger-Typ (ADR-Prosa-Aussage), bekommen **keine** `evidence/`-Datei, sondern
  stehen mit Finding-Kennung in der Closure-Notiz. F-1 ist LOW, vom Reviewer vor dem Merge gefunden, Träger ist eine ADR-Festlegung — nach dem Wortlaut des Deckels
  ein Eintrag in der Closure-Notiz (`F-1`, Review `review-slice-sdk-meldungscodes-in-fehlertypen`), keine neue `evidence/`-Datei. (Es wäre das 12. Vorkommen nach dem Zähler des Registers,
  nicht das 9.; die Zahl 9 ist am Register nicht nachvollziehbar.)
- F-2 des Reviews (und V-1): das ist **nicht** das erste Vorkommen. Die Klasse
  [`drei-sprachen-kopie-divergiert-am-randfall`](../plan/planning/observations/BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall/state.md)
  steht auf **gestrichen** (3×, akzeptiertes Negativ, Architect-Verdikt `welle-routing-lese-schritt`), mit dem Wiederaufnahme-Trigger „ein Divergenz-Fund **nach** dem Merge,
  oder ein Fund ab MEDIUM“. Beides trifft nicht zu (LOW, vor dem Merge, in der Fixrunde behoben); es ist aber der erste **Produktivcode**-Divergenz-Fund dieser Klasse (die
  bisherigen Belege waren Lesemodell, Doku-Satz und README), und die Gegenmaßnahme „derselbe Eingabesatz je Sprache im Plan“ lief an der **Tabelle** selbst ins Leere (`error` mit
  falschem Typ stand nicht in ihr). Vorschlag: Vermerk „kein neuer Anfall nach Wiederaufnahme-Trigger“ mit dieser Einordnung, keine neue Datei; V-1 als zweiter Randfall im Vermerk.
- Beobachtungs-Register neu: kein neuer Eintrag nötig; führt der Planner V-2 (Zahlen im Plan ohne Messzeitpunkt) als Klasse, ist das ein Fall der bestehenden Instanz-A-Regel
  ([`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md)), vom Verifier gefunden.

## 7. Bedingungen und offene Punkte für den Planner

**Bedingungen vor der Closure:**
- **B-1:** Der SDK-Release trägt die Nummer **0.6.0** je Sprache; sonst ist die Handbuch-Aussage „ab Package-Version 0.6.0“ (V-3) vor dem Release zu berichtigen.
- **B-2:** V-2 im Plan nachziehen (Messzeitpunkt an die Zahlen der Tabellen „Belege des Implementers“ und „Mutationsbeleg“, oder die Zahlen vom `HEAD` eintragen) und den
  Tippfehler an Plan Z. 414/415 beheben — Plan-Edit ist Planner-Sache.
- **B-3:** Der DoD-Punkt „Mutationsbeleg je Sprache“ trägt nach meiner Nachfahrt M1 (3 Sprachen), M2 (Kotlin), F-2 (3), Feldnummer (C#); M3, M4 und M2 in C#/Python stehen auf Implementer-Angabe. Der
  Planner entscheidet, ob das genügt oder ob er M3/M4 nachfahren lässt.

**Offene Punkte (kein Rückgabe-Pfeil):** V-1 (Groß-/Kleinschreibung der Feldnamen, nur C#), V-4 (Text-Änderung bei nicht-String `error` in den READMEs nicht genannt) — beides ohne Wirkung am echten Server.

## 8. Release-Bereitschaft SDK 0.6.0 (nur melden, nicht ausgeführt)

Vor dem Release je Sprache noch zu tun (Plan §5 Folgeschritt, Versionsänderung ist ausdrücklich nicht Teil des Slice):

| Datei | Stelle | Ist | Soll |
|---|---|---|---|
| `sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.csproj` | Z. 21 `<Version>` | 0.5.0 | 0.6.0 |
| `sdks/python/pgchangefeed/pyproject.toml` | Z. 7 `[project] version` | 0.5.0 | 0.6.0 |
| `sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts` | Z. 86 `version` und Z. 210 Publikations-`version` | 0.5.0 | 0.6.0 (beide) |
| `sdks/kotlin/pgchangefeed-kotlin/README.md` | Z. 36 Abhängigkeit `pgchangefeed-kotlin:0.5.0` | 0.5.0 | 0.6.0 |

Keine README-„Upgrading“-Zeile nötig (Festlegung 4, additive Änderung; die vorhandenen Einträge „0.5.0 — …“ bleiben Historie). Tags `sdk-csharp-v0.6.0`, `sdk-python-v0.6.0`,
`sdk-kotlin-v0.6.0`; der Tag-Check der Workflows gleicht jeden Tag gegen die Versionsdatei ab (Abbruch bei Abweichung). Nach dem Push je ein realer, grüner Post-Push-Lauf der drei Release-Workflows
([`AGENTS.md`](../../AGENTS.md) §3.10, `docs/maintainer/releasing.md` §4). Der Server-Release 0.6.0 (`docs/user/version.md`) ist davon unabhängig. Die Release-Notiz nennt die Eigenschaft,
ohne Neukompilierungs-Hinweis (Festlegung 4).
