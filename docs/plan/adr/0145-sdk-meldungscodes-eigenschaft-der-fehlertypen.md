# ADR-0145: SDK-Fehlertypen tragen den Meldungscode als Eigenschaft `MessageCode` — Eigenlesung des gRPC-Statusdetails ohne neue Abhängigkeit (Schärft ADR-0144)

**Status:** Accepted — **kein** Supersedes.

**Datum:** 2026-10-03

**Autor:** pt9912 (Architect-Rolle, Modul 8; anderer Kontext als der Planner-Lauf
des Slice `sdk-meldungscodes-in-fehlertypen`, dessen Fragen A1 bis A3 diese ADR beantwortet)

**Bezug:** [`LH-FA-SST-009`](../../../spec/lastenheft.md) (SDK-Packages),
[ADR-0144](0144-meldungscodes-nutzerseitige-kennungen.md) (Meldungscodes; Festlegung 3
tabelliert die Wire-Form, der Satz „die SDKs ändern sich in dieser ADR nicht“ ist der
Anlass), [ADR-0106](0106-csharp-nuget-erstes-sdk-package.md),
[ADR-0107](0107-python-pypi-zweites-sdk-package.md),
[ADR-0109](0109-kotlin-github-packages-drittes-sdk-package.md) (Package-Form je
Sprache), [ADR-0134](0134-sdk-public-doc-check-gate-make-gates.md) (keine interne
Kennung unter `sdks/`), [ADR-0083](0083-herkunft-von-aussagen-in-traegern.md),
`AGENTS.md` §3.5, §3.8, §3.12.

**Schärft:** [ADR-0144](0144-meldungscodes-nutzerseitige-kennungen.md) Festlegung 3, Zeile
„Die SDKs unter `sdks/` ändern sich in dieser ADR nicht“: diese ADR legt fest, was die
SDKs mit dem Code tun; ADR-0144 bleibt unverändert (`AGENTS.md` §3.5). Spec-Stelle: die
SDK-Zeilen [`SPEC-026`](../../../spec/pflichtenheft.md) bis `SPEC-028` nennen Fehlertypen
nicht (gemessen im Plan, 0 Treffer); kein Spec-Nachzug nötig.

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

**(1) Anlass.** Der Server trägt seit dem Release 0.6.0 den Meldungscode im HTTP-Fehlerkörper
(Feld `code`) und im gRPC-Statusdetail (`google.rpc.ErrorInfo`, `reason = <Code>`,
`domain = pg-change-feed`). Die drei SDK-Packages (C#, Kotlin, Python; Stand 0.5.0) reichen
heute nur den Text durch; ein Anwender müsste ihn parsen. Das SDK-Package ist eine öffentliche
API (NuGet, PyPI, GitHub Packages, Cloudsmith): Name, Form und Semantik der Eigenschaft sind
ein Vertrag, den ein späterer Wechsel bricht — deshalb eine ADR statt Plan-Vorschlag.

**(2) Ist-Stand der Fehlertypen** (gelesen, 2026-10-03): die HTTP-Fehlertypen und die
gRPC-Fehlertypen haben je Sprache eine abstrakte Basis; die Blatt-Typen sind in C# `public sealed`
mit **öffentlichem** Konstruktor, in Kotlin Unterklassen einer `sealed class`. Python:
`PgChangeFeedGrpcError.code` ist bereits die `grpc.StatusCode` (`exceptions.py`). Der
README-Hinweis „Upgrading 0.5.0“ in C# zeigt, dass dieses Package eine binäre Signaturänderung
bisher als Hinweis behandelt hat, nicht als Hindernis.

**(3) Das Statusdetail ist in keinem Package gelesen.** Die Stubs kennen `google.rpc.*` nicht
(Plan-Messung: `git grep` auf `ErrorInfo|grpc-status-details|google\.rpc|StatusProto|grpcio-status|CommonProtos`
unter `sdks/` druckt am Parent 0 Zeilen). Vorhandene Laufzeit-Abhängigkeiten laut Plan
(*übernommen*, nicht von mir nachgemessen): C# `Google.Protobuf` und `Grpc.Net.Client`,
Kotlin `grpc-protobuf` und `protobuf-java`, Python `protobuf` und `grpcio`, **ohne**
`grpcio-status`.

**(4) Wire-Definition** (*gelesen* am 2026-10-03 in einer lokalen Kopie des googleapis-Include-Baums
eines anderen Projekts; Version der Kopie nicht festgestellt, die Nummern sind in Protobuf
nicht änderbar, ohne das Wire-Format zu brechen): `google.rpc.Status` Feld 1 `code`, Feld 2
`message`, Feld 3 `repeated google.protobuf.Any details`; `google.protobuf.Any` Feld 1
`type_url`, Feld 2 `value`; `google.rpc.ErrorInfo` Feld 1 `reason`, Feld 2 `domain`, Feld 3
`metadata` (Map). Der Trailer heißt `grpc-status-details-bin` (gRPC-Spezifikation,
*übernommen*; der Realserver-Fall belegt ihn am Server-Wire).

## Entscheidung

Wir wählen **eine Eigenschaft `MessageCode` je Fehlertyp-Hierarchie, ohne neue
Laufzeit-Abhängigkeit, mit Eigenlesung des Statusdetails, und mit Überladungen statt
Parameter-Änderung**. Die vier Fragen des Plans in Festlegungen:

1. **Name und Typ (A2).** `MessageCode` (C#, `string?`), `messageCode` (Kotlin, `String?`),
   `message_code` (Python, `str | None`). Nicht `Code`/`code`: in Python ist `code` am
   gRPC-Fehler die `grpc.StatusCode`; C# `StatusCode` und Kotlin `statusCode` liegen daneben.
   Die Eigenschaft hängt an **beiden** Basen je Sprache (HTTP und gRPC), alle Blatt-Typen
   erben sie. `MalformedResponse` und die NATS-Fehler bekommen keinen Wert (kein Fehlerkörper
   des Servers; Eigenschaft ist dort leer).
2. **Semantik.** Leer ist `null`/`None`. Der Wert wird **unverändert durchgereicht**, ohne
   Formatprüfung gegen `PCF-…` (ein künftiges Präfix braucht kein SDK-Release). HTTP: nur ein
   JSON-**String** zählt; fehlend, `null`, leerer String und jeder andere JSON-Typ (z. B. eine
   Zahl) lesen leer, der Fehlertext bleibt erhalten und der typisierte Fehler wird geworfen.
   gRPC: der erste `ErrorInfo` mit `domain == "pg-change-feed"` liefert `reason`; fremde Domäne,
   anderer Detailtyp, leerer `reason`, fehlende oder nicht lesbare Bytes → leer. Der Parser
   wirft nie. Der Code ist unabhängig vom Status: ein Code auf einem unerwarteten Status
   (`Unavailable`, HTTP `503`) wird geliefert; Auth-Fehler tragen beim Server keinen Code und
   lesen deshalb leer.
3. **Statusdetail ohne Abhängigkeit (A1).** Je Sprache genau **eine** interne Hilfsfunktion
   liest den Binär-Trailer `grpc-status-details-bin` mit der vorhandenen Protobuf-Laufzeit
   (C# `CodedInputStream`, Kotlin `CodedInputStream`, Python `google.protobuf` mit einem
   Leser für Varint und längenbegrenzte Felder) und nur die Felder aus Kontext 4. Keine neue
   Laufzeit-Abhängigkeit in `.csproj`/`Directory.Packages.props`, `build.gradle.kts`,
   `pyproject.toml`. Die drei Parser prüfen dieselbe Tabelle (Plan §3), auch wenn eine
   Sprache transitiv eine fertige Klasse hätte.
4. **Konstruktoren (A2).** **Überladung statt Parameter-Änderung:** C#: der bisherige
   Konstruktor bleibt, je Typ kommt ein zweiter mit `string? messageCode` hinten dazu (an
   der Basis ebenso); Kotlin: `@JvmOverloads` am Primär-Konstruktor mit `messageCode: String? = null`
   hinten (erzeugt die alte JVM-Signatur weiter); Python: Keyword-only `message_code=None`
   (`*, message_code=None`). Damit bleibt jede 0.5.x-Signatur binär und quellseitig erhalten;
   die Änderung ist rein additiv. Ein Eintrag „Upgrading 0.6.0“ in den READMEs entfällt, die
   Release-Notiz nennt die neue Eigenschaft ohne Hinweis auf Neukompilierung.
5. **Umfang (A3).** Typisiert wird, was heute typisiert ist: HTTP-Client, SSE-Client beim
   Öffnen, gRPC-Administrations-Client. Die gRPC-Stream-Clients bleiben roh
   (`RpcException`, `StatusException`, `grpc.RpcError`); das README sagt, dass der Code dort
   über die native Statusdetail-API des gRPC-Clients lesbar bleibt. NATS-Fehler bleiben
   außerhalb. Eine Typisierung der Streams ist eine eigene Entscheidung.
6. **Version.** Der Schritt ist additiv (Minor); 0.6.0 je Sprache als Folgeschritt nach der
   Verifikation, nicht im Slice. Die `Diagnose`-Antwort trägt `error_code` über die erzeugten
   Stubs ohne Mapper-Code.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun, Text parsen | kein Aufwand | der Text ist laut ADR-0144 Festlegung 4 nicht stabil; der Anwender bräuchte dafür einen Parser |
| B — Code als Präfix im Text der Ausnahme | kein neuer Typ | verschiebt die Parser-Last zum Anwender; Text ist frei |
| **C — Eigenschaft, Eigenlesung des Statusdetails (gewählt)** | keine neue Abhängigkeit im veröffentlichten Package; drei Parser, eine Tabelle; additiv | handgeschriebener Wire-Leser (drei Felder; Fehler am Wire belegt der Realserver-Fall) |
| D — fertige Klassen: C# `Google.Api.CommonProtos`, Python `grpcio-status` + `googleapis-common-protos`, Kotlin `StatusProto` | weniger eigener Code | je eine neue Laufzeit-Abhängigkeit im Package, mit Version, Lizenz, Pin (AGENTS §3.8-Geist) und Folgepflege; drei Sprachen mit verschiedenen API-Formen, die Tabelle wäre dreimal anders geprüft; für drei Felder |
| E — optionaler Parameter statt Überladung | kürzer | ändert die binäre Signatur in C# und Kotlin; bricht gegen 0.5.x kompilierte Aufrufer, die Fehler selbst bauen (Tests fremden Codes) |
| F — Name `Code`/`code` | kürzer | Kollision mit `PgChangeFeedGrpcError.code` in Python, Verwechslung mit `StatusCode` |

## Konsequenzen

- Positiv: ein Anwender liest den Code als Eigenschaft; die 0.5.x-API bleibt binär und quellseitig
  unverändert; keine neue Lieferketten-Fläche.
- Negativ: drei handgeschriebene Wire-Leser. Milderung: der Plan verlangt, dass die Unit-Tests das
  Detail mit der Protobuf-Laufzeit nachbauen und nicht mit dem Parser, und ein Realserver-Fall je
  Sprache liest den Server-Wire.
- Negativ: zwei Konstruktoren je Typ in C#, `@JvmOverloads` in Kotlin; ein Name, der nach
  Veröffentlichung nicht mehr zu wechseln ist.
- Akzeptiertes Negativ: der Code gilt nur dort, wo der Server einen setzt; ein Server vor 0.6.0
  liefert leer. Der Wire-Leser wertet ein `Any` nur aus, wenn sein `type_url` auf
  `google.rpc.ErrorInfo` endet (Tabellenzeile 8 des Plans: `RetryInfo` → leer).
- Folgepflicht: Plan-Nachzug (Verdikt, Liste); die Release-Notiz je Sprache; `docs/user/benutzerhandbuch.md`
  berichtigt den Satz, die SDKs werteten den Code nicht aus. Keine Spec-Änderung.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Unit-Tabellen je Sprache (hergeleitet; **nicht erprobt**) | die Zeilen der HTTP-Tabelle (12) und der Statusdetail-Tabelle (12) des Plans liefern je Sprache dieselben Codes, einschließlich Zahl im Feld `code` (leer) und fremde Domäne (leer); Erwartung, der Implementer schreibt sie | `make sdk-pack-csharp`, `make sdk-pack-kotlin`, `make sdk-pack-python` |
| Realserver-Fall (hergeleitet; **nicht erprobt**) | `EnableTable` auf eine fehlende Tabelle endet je Sprache an HTTP und gRPC mit `NotFound` und einem Code; belegt Feldnummern und `domain` am Server-Wire (Erwartung `PCF-E8025`, am Server nicht gemessen) | `make test-sdk-csharp-integration`, `make test-sdk-kotlin-integration`, `make test-sdk-python-integration` |
| Mutation (hergeleitet; **nicht erprobt**) | das Entfernen der `domain`-Prüfung oder der Zuweisung färbt je einen Test rot; weder Stelle noch Instanz noch Farbe sind gefahren — „der Implementer fährt sie“ ist eine Erwartung | — |
| Abhängigkeits-Aussage (hergeleitet) | Diff ohne Änderung an `Directory.Packages.props`/`PgChangeFeed.Client.csproj` (Pakete), `build.gradle.kts` (Abhängigkeiten), `pyproject.toml` (`dependencies`) | Review |
| Öffentlicher Text | `PCF-…` ist keine interne Kennung; Messung im Plan (Muster in `tools/harness/sdk-public-doc-check.sh` trifft `PCF-E8025` nicht, *übernommen* aus dem Plan) | `make sdk-public-doc-check` |

## Re-Evaluierungs-Trigger

Verlangt ein Anwender den Code an den Stream-Clients, oder bekommt der Server ein weiteres
Statusdetail, das ein SDK auswerten soll, öffnet das die Frage nach einer fertigen Bibliothek
(Alternative D) neu. Ein Server-Code für Auth-Fehler wird ohne SDK-Änderung durchgereicht.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-03 | Accepted — Architect-Entscheidung vor dem Start des Slice | `sdk-meldungscodes-in-fehlertypen` (Name des Plans) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0145` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
