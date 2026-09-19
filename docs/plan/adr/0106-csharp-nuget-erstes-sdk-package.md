# ADR-0106: C#/NuGet als erstes SDK-Package für Client-Bibliotheken (`LH-FA-SST-009`)

**Status:** Accepted

**Datum:** 2026-09-19

**Autor:** pt9912 (Architect-Rolle, Nutzerentscheidung vom 2026-09-19 im Chat:
„Wir könnten mit C# NuGet package anfangen" — die Kern-Entscheidung; diese ADR
fasst sie vollwertig und entscheidet die daran hängenden, im Pflichtenheft
ausdrücklich offen markierten Sub-Fragen selbst, mit Alternativenvergleich,
Modul 8 §Rollen-Regeln: „Architect schreibt")

**Bezug:** [`LH-FA-SST-009`](../../../spec/lastenheft.md) (Client-Bibliotheken/
SDKs — die Lastenheft-Fähigkeit selbst), [`LH-FA-SST-006`](../../../spec/lastenheft.md)
(HTTP-/JSON-API), [`LH-FA-SST-008`](../../../spec/lastenheft.md) (Live-Change-Stream,
gRPC/SSE/NATS-Vollinhalt), [`ARC-005`](../../../spec/architecture.md)
(Driving Adapters — die Oberflächen, die das SDK als Consumer anspricht),
[`ADR-0057`](0057-http-grpc-api.md) (HTTP/JSON-API-Vertrag),
[`ADR-0060`](0060-grpc-streaming-mechanismus.md) (gRPC-Server-Streaming und
die `.proto`), [`ADR-0061`](0061-http-sse-zusaetzlich-zu-grpc.md) (SSE),
[`ADR-0100`](0100-nats-dritter-vollinhalts-zustellweg.md) (NATS-Vollinhalts-Stream),
[`ADR-0076`](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
(`examples/` — Ort, Import-Grenze, Rolle „Vorbild, kein Belegträger"),
[`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md) und
[`ADR-0090`](0090-beispiel-clients-volle-matrix.md) (Beispiel-Clients
C#/Kotlin — Präzedenzfall für Sprach-Wurzel-Form, Docker-only-Bauform,
Digest-/Paket-Pinnung; regeln aber ausdrücklich nur `examples/`, keine
veröffentlichten Pakete), [`ADR-0068`](0068-wegwerf-clients-begrenzte-import-berechtigung.md)
(Import-Berechtigungs-Disziplin für einen öffentlichen Pfad-Bereich),
[`ADR-0051`](0051-cicd-pipeline-github-actions.md) Entscheidung 3/8
(Release-Tag-Trigger, Docker-Hub-Secret-Muster — Vorbild für den
NuGet-Publish-Mechanismus), [`AGENTS.md`](../../../AGENTS.md) §3.1 (Docker-only)
§3.6 (Gates nur per ADR gelockert) §3.13 (Träger-Nachzug),
[`spec/pflichtenheft.md`](../../../spec/pflichtenheft.md) §2 `SPEC-018`/
`SPEC-020`/`SPEC-021`/`SPEC-024` (die vier Draht-Festlegungen, die das SDK
**benutzt**), `SPEC-023` (Beispiel-Client-Werkzeugketten — Vorbild, nicht
Vorgriff), `docs/user/version.md` (Server-SemVer — von dieser ADR
ausdrücklich **nicht** berührt)

**Schärft:** [`LH-FA-SST-009.a`](../../../spec/pflichtenheft.md) — die
Pflichtenheft-Stelle markiert Sprachmatrix und Vertriebsweg für
`LH-FA-SST-009` ausdrücklich als „offen, ADR-pflichtig"; diese ADR
beantwortet sie für die **erste** Sprache/den ersten Vertriebsweg
(C#/NuGet) und für die daran hängenden Umsetzungsfragen (Umfang, Ort,
Versionierung, Bau-/Publish-Mechanismus). Eine zweite Sprache oder ein
zweiter Vertriebsweg bleibt eine eigene, künftige ADR (§Re-Evaluierungs-Trigger).

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

**Die Anforderung.** [`LH-FA-SST-009`](../../../spec/lastenheft.md) verlangt
offizielle, versionierte Client-Bibliotheken (Packages) für die bestehenden
Zustellwege — HTTP-API (`SPEC-018`), gRPC-Stream (`SPEC-020`), SSE
(`SPEC-021`), NATS-Vollinhalts-Stream (`SPEC-024`) —, die ein Consumer über
den Paketmanager seiner Sprache einbindet, statt das Protokoll selbst zu
implementieren. Welche Sprache(n) zuerst bedient werden und über welchen
Vertriebsweg, ist ausdrücklich Architektur-/Spezifikationsfrage
(Out-of-Scope-Klausel des Lastenhefts). Ausdrücklich **nicht** ausreichend:
die bestehenden Wegwerf-/Beispielprogramme unter `examples/` — sie sind
unversioniert und nicht als eigenständiges, von Dritten konsumierbares
Package veröffentlicht ([`LH-FA-SST-009`](../../../spec/lastenheft.md)
Out-of-Scope; `SPEC-023`).

**Die Nutzerentscheidung.** Der Auftraggeber hat am 2026-09-19 im Chat
entschieden: „Wir könnten mit C# NuGet package anfangen." Das legt Sprache
(C#/.NET) und Vertriebsweg (NuGet) für das **erste** Package fest. Offen
bleiben die daran hängenden technischen Sub-Fragen, die diese ADR trifft:
Umfang des ersten Pakets (welche der vier Zustellwege deckt v1), Verhältnis
zu `examples/csharp/`, Versionierung, und Bau-/Publish-Mechanismus.

**Der Ist-Stand — gemessen, nicht erinnert.**

| Prozedur | Ergebnis |
|---|---|
| `find examples/csharp -maxdepth 1 -type d` | fünf Programme: `http-client`, `grpc-client`, `sse-client`, `nats-client` (Wecksignal, `SPEC-017`), `nats-stream-client` (Vollinhalt, `SPEC-024`) — **alle vier** für `LH-FA-SST-009` relevanten Zustellwege haben bereits ein C#-Vorbild |
| `cat examples/csharp/Directory.Packages.props` | zentral gepinnte Paketversionen: `Grpc.Net.Client` 2.83.0, `Grpc.Tools` 2.84.0, `Google.Protobuf` 3.36.1, `NATS.Net` 3.2.0 (übernommen, Messzeitpunkt 2026-09-17, [`ADR-0090`](0090-beispiel-clients-volle-matrix.md) Festlegung 4 — nicht heute neu gemessen) |
| `cat examples/csharp/http-client/http-client.csproj` | `TargetFramework net10.0`, keine Fremdabhängigkeit — reine BCL (`System.Net.Http`, `System.Text.Json`) |
| `cat examples/csharp/grpc-client/grpc-client.csproj` | liest die `.proto` über einen **zusätzlichen, benannten** Docker-Bau-Kontext (`--build-context proto=proto`), nicht committet im Sprach-Baum ([`ADR-0090`](0090-beispiel-clients-volle-matrix.md) Festlegung 2/3) |
| `grep -rn "sdks\|packages/csharp" .` | kein Treffer — kein bestehender SDK-Baum, keine Namenskollision |
| `cat docs/user/version.md` | `0.1.2` — Server-SemVer, eigenständig geführt, ohne erkennbare Kopplung an Paket-Versionen |
| `.a-check.yml` `languages:` | nur `go` — eine C#-Quelle wird von a-check strukturell nicht gelesen, unabhängig vom Ort |

**Was das ändert.** Die Prämisse „SSE und NATS als Folge-ADR/Folge-Slice
offenlassen, weil dort noch kein Vorbild existiert" trägt **nicht**: Ein
C#-Vorbild existiert bereits für **alle vier** Zustellwege, nicht nur für
HTTP und gRPC. Der Umfang von v1 ist deshalb keine Frage fehlenden
Draht-Verständnisses, sondern eine Frage des **Paket-Zuschnitts** —
Abhängigkeits-Footprint, Erst-Release-Größe, Zahl der zu pflegenden
öffentlichen APIs (§Entscheidung Festlegung 1).

**Fünf Bindungen, die den Lösungsraum prägen — alle aus dem Repo:**

1. **Docker-only** ([`AGENTS.md`](../../../AGENTS.md) §3.1): kein
   Host-`dotnet`. Bau und `pack` laufen im gepinnten Container, wie
   `examples-csharp` es bereits vormacht.
2. **`examples/` bleibt Vorbild, kein Belegträger** (`SPEC-023`): Ein SDK ist
   das Gegenteil — Paket-Konsumenten verlassen sich auf seine Kompatibilität
   und Versionsgrenzen. Beide Klassen dürfen nicht denselben Baum teilen.
3. **Kein neues Gate.** NuGet-Publish braucht Netz und ein Secret — dieselbe
   Klasse wie `make image`/Docker-Hub-Push
   ([`ADR-0051`](0051-cicd-pipeline-github-actions.md) Entscheidung 8):
   Werkzeug, kein Gate.
4. **a-check liest kein C#.** `sdks/csharp/**` braucht keine
   `.a-check.yml`-Gruppe — dieselbe Grenze, die
   [`ADR-0090`](0090-beispiel-clients-volle-matrix.md) Festlegung 6 für
   `examples/csharp` schon zieht.
5. **Eine Arbeit, die eine beschriebene Eigenschaft bewegt, zieht ihre Träger
   nach** ([`AGENTS.md`](../../../AGENTS.md) §3.13): `LH-FA-SST-009.a`s Satz
   „ist offen, ADR-pflichtig" wird mit dieser ADR für die erste
   Sprache/den ersten Vertriebsweg falsch — der Träger-Nachzug ist Sache des
   Implementierungs-Zuges (§Konsequenzen Folgepflicht), nicht dieser
   ausschließlich ADR-schreibenden Rolle (Aufgabenbindung dieses Zuges).

## Entscheidung

Wir wählen: **C#/.NET als erste SDK-Sprache, NuGet als erster
Vertriebsweg.** Fünf Festlegungen:

### 1 — Umfang des ersten Pakets: HTTP-API und gRPC-Stream, ein Package

v1 des Packages (Arbeitsname `PgChangeFeed.Client`) deckt **HTTP-API**
(`SPEC-018`) und **gRPC-Stream** (`SPEC-020`) — nicht SSE (`SPEC-021`) und
nicht NATS-Vollinhalt (`SPEC-024`), obwohl für alle vier bereits ein
C#-Vorbild existiert (§Kontext, Ist-Stand). Begründung, gegen die
Vorbild-Symmetrie abgewogen:

- **Der Vorbild-Einwand widerlegt sich nicht selbst — aber er trägt hier
  nicht.** Dass ein Vorbild existiert, senkt die Kosten für das
  *Wire-Verständnis*; es senkt **nicht** die Kosten für *Paket-Design*
  (öffentliche API-Form, Versionsgrenze, NuGet-Metadaten, Doku, Tests als
  Vertrag) — und diese Kosten fallen **je Zustellweg** an, unabhängig vom
  Vorbild.
- **Abhängigkeits-Footprint ist ein echter Grund, keine Bequemlichkeit.**
  HTTP braucht **keine** Fremdabhängigkeit (reine BCL). gRPC braucht
  `Grpc.Net.Client` + `Google.Protobuf` — für jeden Consumer, auch den, der
  nur HTTP will, sobald beide im selben Package liegen. SSE und NATS legen
  mit `NATS.Net` eine **weitere** Fremdabhängigkeit obendrauf; ein
  Erst-Release, das alle vier bündelt, zwingt jedem HTTP-only-Consumer drei
  Fremdmodule auf, die er nie lädt.
- **HTTP + gRPC ist die Kombination, die die meisten Consumer brauchen:**
  synchroner Abruf/Backfill (`GET /changes`, [`ADR-0081`](0081-changes-lesen-ueber-die-http-api.md))
  und Live-Zustellung (Server-Stream) — SSE und NATS sind **alternative**
  Transporte für dieselbe Live-Zustellungs-Fähigkeit
  ([`LH-FA-SST-008`](../../../spec/lastenheft.md) Out-of-Scope: „konkretes
  Übertragungsprotokoll … das sind Architektur-/Spezifikationsfragen"), kein
  zusätzlicher fachlicher Nutzen gegenüber gRPC für denselben Consumer.
- **Boundary-Akzeptanzkriterium bleibt erfüllt, ohne v1 zu bläuen:**
  [`LH-FA-SST-009`](../../../spec/lastenheft.md) Negative-AC verlangt, dass
  der direkte Zugriff auf die zugrunde liegende Schnittstelle nutzbar
  bleibt, wo keine Bibliothek existiert — SSE/NATS bleiben über
  `examples/csharp/sse-client`/`nats-stream-client` als **Vorbild** (nicht
  Ersatz) einsehbar, bis ein Folge-Paket sie deckt.
- **Ein Package, nicht vier.** HTTP und gRPC teilen Auth-Header-Form (Bearer
  Token) und Grundkonfiguration (Adresse, Token); ein Consumer, der beide
  Wege nutzt, referenziert ein Package statt zwei. Der Abhängigkeits-Preis
  von gRPC (`Grpc.Net.Client`, `Google.Protobuf`) trägt ein reiner
  HTTP-Consumer zwar mit — für **zwei** eng verwandte Zustellwege ist das
  vertretbar; bei vier wäre es das nicht mehr (§Verglichene Alternativen B).

SSE und NATS-Vollinhalt bleiben für ein **Folge-Package** offen — entweder
als v2 desselben Packages oder als eigenständiges Package, je nachdem, wie
sich der Abhängigkeits-Footprint dann bewertet (Entscheidung eines
Folge-Zuges, kein Vorgriff hier).

### 2 — Verhältnis zu `examples/csharp/`: eigenständiges Projekt, neuer Baum

Das Package entsteht in einem **neuen** Verzeichnis, `sdks/csharp/`
(z. B. `sdks/csharp/PgChangeFeed.Client/`), **nicht** durch Weiterentwicklung
eines bestehenden Beispiels. Begründung:

- **`SPEC-023` markiert `examples/` als „Vorbild, kein Belegträger, keine
  Zustandsmaschine"** — ein Paket-Konsument verlässt sich auf Kompatibilität
  über Versionsgrenzen hinweg; das ist die Definition eines Belegträgers.
  Beide Rollen im selben Baum zu führen hieße, `SPEC-023`s Grenze für genau
  die Dateien zu unterlaufen, die zum SDK würden.
- **Der Name `sdks/` statt `packages/`:** `packages/` kollidiert
  begrifflich mit NuGets eigenem Vokabular (ein NuGet-Erzeugnis heißt selbst
  „Package") und würde bei künftigen Sprachen zweideutig, ob ein
  Sprach-Unterordner unter `packages/<sprache>/` ein SDK oder ein
  Drittanbieter-Abhängigkeits-Cache meint. `sdks/` spiegelt exakt den
  Lastenheft-Begriff „Client-Bibliotheken (SDKs)".
- **Kein `ProjectReference` auf `examples/csharp/*`.** Die vorhandenen
  Beispiel-Klassen (`TablesClient.cs`, `TablesUrlBuilder.cs`,
  `Cli.cs`/`Config.cs`/`Format.cs`) bleiben **Referenzmaterial** für den
  Implementierungs-Zug — dieselbe Draht-Kenntnis, aber als eigenständiger,
  paketierbarer Code neu geschrieben, mit einer öffentlichen, stabilen API
  (nicht CLI-Argument-Parsing wie die Beispiele).
- **a-check bleibt unberührt** (§Kontext Bindung 4): `sdks/csharp/**` ist
  keine Go-Quelle und braucht keine `.a-check.yml`-Gruppe, wie
  `examples/csharp/**` es heute schon nicht braucht
  ([`ADR-0090`](0090-beispiel-clients-volle-matrix.md) Festlegung 6).
- **Die Import-Grenze der Beispiele gilt sinngemäß auch hier, verschärft:**
  Das SDK importiert ausschließlich die .NET-BCL und öffentliche
  NuGet-Pakete, **keinen** privaten Baum dieses Repositories — dieselbe
  Regel wie `SPEC-023`s Import-Grenze, hier ohne Ausnahme, weil ein SDK
  außerhalb dieses Repositories baubar sein muss, sobald es veröffentlicht
  ist (nicht nur kopierbar wie ein Beispiel).
- **Der gRPC-Teil erreicht die `.proto` genauso wie `examples/csharp/grpc-client`:**
  über einen zusätzlichen, benannten Docker-Bau-Kontext
  ([`ADR-0090`](0090-beispiel-clients-volle-matrix.md) Festlegung 2) — die
  `.proto` bleibt die einzige Quelle des Draht-Vertrags, kein committeter
  Stub im SDK-Baum.

### 3 — Versionierung: unabhängiges SemVer, losgelöst vom Server

Das Package führt sein **eigenes** SemVer 2.0, unabhängig von
`docs/user/version.md` (der Server-Version). Begründung:

- **Unterschiedliche Änderungstakte.** Ein SDK-Patch (z. B. eine
  Bugfix-Korrektur im Retry-Verhalten) braucht keinen Server-Release, und
  ein Server-Release (neue Fähigkeit ohne Draht-Änderung) braucht keinen
  SDK-Release. Kopplung erzwänge synchrone Releases für unabhängige
  Änderungen — die genannte Boundary-AC von
  [`LH-FA-SST-009`](../../../spec/lastenheft.md) verlangt ohnehin eine
  eigene, erkennbare Versionsgrenze *des Packages* bei einer inkompatiblen
  Protokolländerung, nicht eine geteilte mit dem Server.
- **Der Server selbst kennt bereits zwei unabhängige Versionsräume**
  (`docs/user/version.md` für den Server-Container, die Baseline-Version in
  `harness/conventions.md` für den Harness) — ein dritter, unabhängiger Raum
  für das SDK ist keine neue Klasse, sondern dieselbe Disziplin fortgesetzt.
- **Start bei `0.x.y`** (z. B. `0.1.0`), bis die öffentliche API als stabil
  gilt — Analogie zum Server, der ebenfalls bei `0.1.2` (Stand dieser ADR,
  `docs/user/version.md`) noch vorstabil ist. `1.0.0` markiert bewusst die
  erste Version mit versprochener Rückwärtskompatibilität, kein
  automatischer Übertrag vom Server-Stand.
- **Der SemVer-Major der AC-Boundary bindet an das Draht-Protokoll, nicht an
  den Server-Container:** Ändert sich `SPEC-018`/`SPEC-020` inkompatibel,
  hebt das SDK seinen eigenen Major — unabhängig davon, ob der
  Server-Container zeitgleich eine neue Version bekommt.

### 4 — Build-/Publish-Mechanismus: Docker-only bauen, Netz-Werkzeug veröffentlichen

- **Bauen und `pack` bleiben Docker-only** ([`AGENTS.md`](../../../AGENTS.md)
  §3.1), analog zu `make examples-csharp`: ein neues, netzlos **nicht**
  prüfbares (NuGet-Restore braucht Netz) Werkzeug-Ziel — Arbeitsname
  `make sdk-pack-csharp` — baut, testet und paketiert (`dotnet pack`) im
  gepinnten `mcr.microsoft.com/dotnet/sdk`-Image, erzeugt eine `.nupkg` als
  Artefakt. **Kein Gate** — dieselbe Begründung wie bei
  `make examples-csharp` ([`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md)
  Festlegung 4): Paket-Bezug braucht Netz, `make gates` bleibt netzlos.
- **Veröffentlichung nach NuGet.org braucht einen eigenen, Netz-bindenden
  Workflow und ein Secret** — Arbeitsname `NUGET_API_KEY`, dieselbe Klasse
  wie `DOCKERHUB_USERNAME`/`DOCKERHUB_TOKEN`
  ([`ADR-0051`](0051-cicd-pipeline-github-actions.md) Entscheidung 8): ein
  Repository-Secret, das der Workflow referenziert, nicht anlegt. Der
  Aufruf ist `dotnet nuget push … --api-key $NUGET_API_KEY --source
  https://api.nuget.org/v3/index.json`.
- **Der Trigger ist ein eigener Tag-Namensraum, nicht `v*`.**
  `.github/workflows/release.yml` reagiert bereits auf `push: tags: ['v*']`
  für den Server-Container ([`ADR-0051`](0051-cicd-pipeline-github-actions.md)
  Entscheidung 3) — dieselbe Präfix-Familie für das SDK zu verwenden würde
  beide Release-Räume kollidieren lassen (Festlegung 3, unabhängiges
  SemVer). Vorschlag für den umsetzenden Zug: `sdk-csharp-v*` als eigener
  Tag-Präfix, eigener Workflow (`sdk-csharp-release.yml`), der den Tag
  strikt gegen SemVer 2.0 validiert und gegen die im `.csproj` geführte
  `<Version>` abgleicht — dasselbe Muster wie `release.yml`s
  Tag-gegen-`docs/user/version.md`-Abgleich, nur gegen die
  Projektdatei statt gegen eine separate Markdown-Datei.
- **Beides ist Folgepflicht dieser ADR, nicht Teil ihres Umfangs** — diese
  ADR entscheidet die *Form*, sie implementiert **nichts**
  (§Konsequenzen Folgepflicht).

### 5 — Was diese ADR nicht ändert

- **Kein Produktionscode, kein Eingriff in `internal/**`/`cmd/**`.**
- **`examples/csharp/**` bleibt unverändert** — kein Umzug, keine
  Code-Extraktion in dieser ADR.
- **`spec/architecture.md` bleibt unberührt** — `sdks/**` ist ein
  Gate-Scope-Bereich außerhalb des Servers, keine `ARC-*`-Komponente.
- **`docs/user/version.md` bleibt unberührt** — die Server-Version bleibt
  eigenständig geführt (Festlegung 3).
- **Kein neues Gate, keine Schwellen-Senkung.**
- **`.a-check.yml` bleibt unberührt** — a-check liest kein C# (§Kontext
  Bindung 4).

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

### A — Erste Sprache und Vertriebsweg

| Option | Pro | Contra |
|---|---|---|
| A1 — Go zuerst (Go-Modul-Registry) | spiegelt die Server-Implementierungssprache; ein Go-SDK könnte interne Typen teilen | Go-Consumer haben ohnehin die **geringste** Einstiegshürde über die rohe API (`net/http`, `google.golang.org/grpc` direkt, keine Fremdsprache zu überbrücken) — der Nutzen eines SDKs ist hier am kleinsten, gerade wo er am billigsten wäre; keine bestehende Nutzerentscheidung dafür |
| A2 — TypeScript/npm zuerst | breites Ökosystem, viele potenzielle Consumer (Web/Node) | keine bestehende Beispiel-Client-Vorarbeit in dieser Sprache (kein `examples/typescript/`); höheres Risiko, Draht-Annahmen beim ersten Schreiben falsch zu treffen, ohne die Gegenprobe eines bereits laufenden Beispiels; keine bestehende Nutzerentscheidung dafür |
| **A3 — C#/NuGet zuerst (gewählt)** | bereits **vier** funktionierende C#-Beispiel-Clients als reale Draht-Vorarbeit (`examples/csharp/{http,grpc,sse,nats-stream}-client`); .NET trägt eine starke, gepflegte Konvention für generierte gRPC-Clients (`Grpc.Tools`, bereits im Bau erprobt); NuGet ist ein etablierter, gut automatisierbarer Vertriebsweg (`dotnet pack`/`dotnet nuget push`); explizite Nutzerentscheidung | NuGet.org-Publish braucht ein neues Secret und einen neuen Workflow (Folgepflicht); ein zweites Sprach-SDK bleibt offen und wartet auf eine künftige ADR |
| A4 — nichts tun / abwarten | kein Aufwand jetzt | `LH-FA-SST-009` bleibt ohne jedes SDK unbeantwortet; die AC „Consumer bindet ein Package ein, ohne das Protokoll selbst zu implementieren" bleibt strukturell unerfüllt; die explizite Nutzerentscheidung vom 2026-09-19 würde ignoriert |

### B — Umfang des ersten Pakets

| Option | Pro | Contra |
|---|---|---|
| B1 — nur HTTP-API | kleinster Erst-Release, keine Fremdabhängigkeit, geringstes Risiko | deckt Live-Zustellung nicht ab — die AC „Consumer empfängt Changes" wird nur über Polling erfüllt, nicht über einen der beiden Live-Wege, obwohl gRPC-Vorbild bereits existiert |
| **B2 — HTTP-API + gRPC-Stream, ein Package (gewählt)** | deckt Abruf **und** Live-Zustellung; deckt die zwei Wege, die die meisten Consumer kombinieren; ein Package statt vier hält die Einbindung für den Consumer einfach | zieht `Grpc.Net.Client`/`Google.Protobuf` auch HTTP-only-Consumern auf |
| B3 — alle vier Zustellwege in einem Package | „vollständige Matrix sofort", keine offene Zelle | größter Erst-Release-Umfang bei einem **ersten** Paket-Design (das noch niemand geprüft hat); zwingt jedem Consumer alle vier Fremdabhängigkeiten (`Grpc.*`, `NATS.Net`) auf, auch wenn er nur einen Weg nutzt — der Footprint-Einwand von B2 gilt hier verschärft |
| B4 — vier separate Packages, eines je Zustellweg | jeder Consumer lädt nur, was er braucht | vier NuGet-IDs, vier Versionsräume, vier Publish-Workflows für ein **erstes** SDK — Overhead, bevor auch nur ein Package real veröffentlicht wurde; verfrühte Modularisierung ohne Nutzungsdaten, die sie rechtfertigen |

### C — Verhältnis zu `examples/csharp/`

| Option | Pro | Contra |
|---|---|---|
| C1 — ein bestehendes Beispiel (z. B. `http-client`) wird zum SDK weiterentwickelt | kein neuer Baum, kein doppeltes Schreiben ähnlichen Codes | verletzt `SPEC-023`s Klassentrennung „Vorbild, kein Belegträger" an genau der Stelle, an der ein Dritter sich künftig auf Stabilität verlässt; CLI-Struktur der Beispiele passt nicht zu einer öffentlichen Bibliotheks-API |
| **C2 — eigenständiges Projekt unter `sdks/csharp/` (gewählt)** | klare räumliche und rollenbezogene Trennung (Vorbild vs. Belegträger); eigene Import-Grenze, eigene Versionierung, eigener Bauweg, ohne die Beispiel-Verträge zu berühren | ähnlicher Code entsteht zweimal (Beispiel und SDK) — vertretbar, weil beide unterschiedliche Verträge tragen (CLI-Demo vs. öffentliche API) |
| C3 — `examples/csharp/` wird komplett nach `sdks/csharp/` verschoben, Beispiele verschwinden | ein Baum weniger zu pflegen | zerstört `SPEC-023`s Beispiel-Anspruch (Vorbild, Docker-only-Startform, Handbuch-Bindung) ohne Not; Beispiele und SDK haben unterschiedliche Konsumenten (Repo-Leser vs. NuGet-Konsument) |

### D — Versionierung

| Option | Pro | Contra |
|---|---|---|
| D1 — gekoppelt an `docs/user/version.md` | ein Versionsraum, ein Blick genügt | erzwingt SDK-Releases bei jedem Server-Release ohne Draht-Änderung und umgekehrt; die AC-Boundary (SDK-eigene Versionsgrenze bei Draht-Änderung) passt nicht zu einer Server-Container-Version |
| **D2 — eigenständiges SemVer, ab `0.x.y` (gewählt)** | Release-Takt folgt der tatsächlichen Änderung (Draht vs. Server); Boundary-AC direkt abbildbar; Analogie zu bereits zwei bestehenden unabhängigen Versionsräumen im Repo | ein dritter Versionsraum, den ein Betrachter separat nachschlagen muss |
| D3 — Versionierung nach Datum (CalVer) | keine Debatte über Major/Minor/Patch | passt nicht zur AC-Boundary, die ausdrücklich „SemVer-Major" als Beispiel nennt; unüblich im NuGet-Ökosystem |

### E — Build-/Publish-Mechanismus

| Option | Pro | Contra |
|---|---|---|
| E1 — Host-`dotnet` für Pack/Publish | einfachster erster Schritt | verletzt [`AGENTS.md`](../../../AGENTS.md) §3.1 (Docker-only) direkt |
| **E2 — Docker-only Pack, separater Netz-Workflow mit eigenem Tag-Präfix und Secret (gewählt)** | konsistent mit `make examples-csharp`/`make image`; Trennung von Server- und SDK-Releases (Festlegung 3); Secret-Muster bereits bewährt (`DOCKERHUB_TOKEN`) | neuer Workflow und neues Secret sind Folgepflicht, nicht heute geliefert; ungeprüft bis zum ersten realen Lauf (`AGENTS.md` §3.10-Klasse) |
| E3 — SDK-Tag im selben `v*`-Namensraum wie der Server | ein Tag-Muster, ein Workflow-Trigger | kollidiert mit der unabhängigen SDK-Versionierung (Festlegung 3) — ein SDK-Patch-Tag `v0.1.1` wäre von einem Server-Tag `v0.1.1` nicht unterscheidbar |

**Fazit:** A3, B2, C2, D2, E2. A1/A2 scheitern am fehlenden Nutzen bzw. an
der fehlenden Vorarbeit, A4 an der unbeantworteten Anforderung; B1 lässt
Live-Zustellung aus, B3/B4 überziehen den Erst-Release bzw. modularisieren
verfrüht; C1/C3 verletzen `SPEC-023`s Klassentrennung; D1/D3 passen nicht
zur Boundary-AC bzw. zum Ökosystem; E1 verletzt Docker-only, E3 kollidiert
mit der unabhängigen Versionierung.

## Konsequenzen

- **Positiv:** `LH-FA-SST-009` bekommt einen konkreten, geschnittenen ersten
  Umsetzungspfad statt einer offenen Frage; die vier vorhandenen
  C#-Beispiel-Clients senken das Umsetzungsrisiko für Draht-Verständnis
  spürbar; die räumliche Trennung von `examples/` und `sdks/` hält beide
  Verträge (`SPEC-023` vs. ein SDK-Verhalten) sauber getrennt.
- **Negativ:** Ein HTTP-only-Consumer lädt `Grpc.Net.Client`/
  `Google.Protobuf` transitiv mit, obwohl er sie nie aufruft (Festlegung 1);
  SSE- und NATS-Consumer bleiben vorerst auf den direkten Zugriffsweg
  verwiesen (durch die Negative-AC von
  [`LH-FA-SST-009`](../../../spec/lastenheft.md) gedeckt, aber ein realer
  Komfortverlust); ein zusätzlicher Release-/Secret-Mechanismus (NuGet)
  entsteht neben dem bestehenden Docker-Hub-/GHCR-Mechanismus.
- **Folgepflicht:**
  1. Ein Planner-Zug schneidet die Umsetzung (voraussichtlich mehrere
     Slices: SDK-Projektgerüst, HTTP-Client-Fläche, gRPC-Client-Fläche,
     Pack-Werkzeug, Publish-Workflow) — diese ADR entscheidet die Form,
     nicht den Schnitt.
  2. `spec/pflichtenheft.md` §1 (`LH-FA-SST-009.a`) und §6 (Externe
     Verträge) brauchen einen Träger-Nachzug, sobald das Package real
     existiert (`AGENTS.md` §3.13) — Sache des umsetzenden Zuges, nicht
     dieser ausschließlich ADR-schreibenden Rolle.
  3. `harness/README.md` §Werkzeuge bekommt seine Zeilen für
     `make sdk-pack-csharp` und den Publish-Workflow erst, wenn sie
     real existieren (`AGENTS.md` §4: kein Träger nennt ein Target, das es
     nicht gibt).
  4. `docs/user/benutzerhandbuch.md` bekommt einen SDK-Hinweis je
     gedeckter Oberfläche, im selben Zug wie das jeweilige Client-Programm
     (Muster: `BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`).
  5. Der NuGet-Publish-Workflow gilt nach [`AGENTS.md`](../../../AGENTS.md)
     §3.10 erst nach einem realen, grünen Post-Push-Lauf als abgeschlossen.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `a-check` | keine — `sdks/csharp/**` ist keine Go-Quelle, a-check liest sie nicht (`.a-check.yml` `languages: go`) | — |
| Review (kein Sensor) | `sdks/csharp/**` importiert ausschließlich .NET-BCL und öffentliche NuGet-Pakete, keinen privaten Baum dieses Repos (Import-Grenze, analog `SPEC-023`) | — |
| künftiges Bau-Werkzeug (Folgepflicht) | `make sdk-pack-csharp` baut/testet/paketiert Docker-only, netzlos geprüfte Teile fahren im Bau, kein Gate | `make sdk-pack-csharp` (noch nicht existent) |

## Re-Evaluierungs-Trigger

1. **Eine zweite Sprache oder ein zweiter Vertriebsweg wird für
   `LH-FA-SST-009` verlangt** (z. B. TypeScript/npm, Go/Go-Modul-Registry,
   Kotlin/Maven) — eigene ADR, diese ADR bleibt für C#/NuGet unberührt.
2. **Das C#-Package wurde real auf NuGet.org veröffentlicht und
   Nutzungsdaten (Download-Zahlen, Issue-Nachfrage) legen eine andere
   Priorisierung nahe** — insbesondere ob SSE/NATS als v2 desselben Packages
   oder als eigenständiges Package sinnvoller sind (Festlegung 1 offen
   gelassen).
3. **Der Abhängigkeits-Footprint von gRPC im selben Package erweist sich als
   echtes Consumer-Problem** (z. B. wiederholte Nachfrage nach einem
   HTTP-only-Package ohne `Grpc.*`) — dann Aufspaltung in getrennte Packages
   nach dem in Alternative B4 verworfenen, aber dann neu zu bewertenden
   Muster.
4. **`docs/user/version.md`s Server-SemVer erreicht `1.0.0`** — kein
   automatischer Trigger für das SDK-SemVer (Festlegung 3 bleibt
   unabhängig), aber ein sinnvoller Anlass, die Reife-Erwartung beider
   Versionsräume gegeneinander zu prüfen.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-19 | Proposed | dieser Architect-Zug (kein vorausgehender Slice-Plan — ADR vor Implementierung, Modul 8) |
| 2026-09-19 | Accepted | Nutzerentscheidung im Chat vom 2026-09-19 („Wir könnten mit C# NuGet package anfangen") + diese vollwertige ADR-Fassung samt Alternativenvergleich |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0106` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
