# ADR-0107: Python/PyPI als zweites SDK-Package für `LH-FA-SST-009`

**Status:** Accepted

**Datum:** 2026-09-19

**Autor:** pt9912 (Architect-Rolle, Nutzerentscheidung vom 2026-09-19 im Chat:
„Können wir noch ein python-SDK erstellen" — die Kern-Entscheidung „Python als
zweite Sprache"; diese ADR fasst sie vollwertig, wählt den Vertriebsweg
begründet (nicht angenommen) und entscheidet die daran hängenden Sub-Fragen
selbst, mit Alternativenvergleich, Modul 8 §Rollen-Regeln: „Architect
schreibt")

**Bezug:** [`LH-FA-SST-009`](../../../spec/lastenheft.md) (Client-Bibliotheken/
SDKs — die Lastenheft-Fähigkeit selbst), [`LH-FA-SST-006`](../../../spec/lastenheft.md)
(HTTP-/JSON-API), [`LH-FA-SST-008`](../../../spec/lastenheft.md) (Live-Change-Stream,
gRPC/SSE/NATS-Vollinhalt), [`ARC-005`](../../../spec/architecture.md)
(Driving Adapters — die Oberflächen, die das SDK als Consumer anspricht),
[`ADR-0106`](0106-csharp-nuget-erstes-sdk-package.md) (erstes SDK-Package,
C#/NuGet — Formvorbild dieser ADR, von ihr **unberührt**: `§Re-Evaluierungs-Trigger`
Punkt 1 sieht eine zweite Sprache/einen zweiten Vertriebsweg als eigene ADR
ausdrücklich vor), [`ADR-0057`](0057-http-grpc-api.md) (HTTP/JSON-API-Vertrag),
[`ADR-0060`](0060-grpc-streaming-mechanismus.md) (gRPC-Server-Streaming und
die `.proto`), [`ADR-0061`](0061-http-sse-zusaetzlich-zu-grpc.md) (SSE),
[`ADR-0100`](0100-nats-dritter-vollinhalts-zustellweg.md) (NATS-Vollinhalts-Stream),
[`ADR-0076`](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
(`examples/` — Ort, Import-Grenze, Rolle „Vorbild, kein Belegträger" —
diese ADR stellt fest, dass für Python **kein** solches Vorbild existiert),
[`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md) und
[`ADR-0090`](0090-beispiel-clients-volle-matrix.md) (Beispiel-Clients
C#/Kotlin — Präzedenzfall für Sprach-Wurzel-Form und Docker-only-Bauform;
regeln ausdrücklich nur `examples/`, keine veröffentlichten Pakete),
[`ADR-0068`](0068-wegwerf-clients-begrenzte-import-berechtigung.md)
(Import-Berechtigungs-Disziplin für einen öffentlichen Pfad-Bereich),
[`ADR-0051`](0051-cicd-pipeline-github-actions.md) Entscheidung 3/8
(Release-Tag-Trigger, Registry-Secret-Muster — Vorbild für den
PyPI-Publish-Mechanismus, wie `ADR-0106` es für NuGet bereits nutzte),
[`AGENTS.md`](../../../AGENTS.md) §3.1 (Docker-only) §3.6 (Gates nur per ADR
gelockert) §3.13 (Träger-Nachzug), [`spec/pflichtenheft.md`](../../../spec/pflichtenheft.md)
§2 `SPEC-018`/`SPEC-020`/`SPEC-021`/`SPEC-024` (die vier Draht-Festlegungen,
die das SDK **benutzt**), `SPEC-026` (das bereits existierende C#-Package —
Formvorbild für einen künftigen Python-Sprach-Eintrag), `docs/user/version.md`
(Server-SemVer — von dieser ADR ausdrücklich **nicht** berührt)

**Schärft:** [`LH-FA-SST-009.a`](../../../spec/pflichtenheft.md) — dieselbe
Pflichtenheft-Stelle, die [`ADR-0106`](0106-csharp-nuget-erstes-sdk-package.md)
für C#/NuGet beantwortet hat, bleibt nach ihrem eigenen Wortlaut für „eine
zweite Sprache oder einen zweiten Vertriebsweg" offen; diese ADR beantwortet
sie jetzt für Python/PyPI als **zweite** Sprache/zweiten Vertriebsweg, und für
die daran hängenden Umsetzungsfragen (Umfang, Ort, Versionierung,
Bau-/Publish-Mechanismus). Eine dritte Sprache bleibt eine eigene, künftige
ADR (§Re-Evaluierungs-Trigger).

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

**Die Anforderung.** [`LH-FA-SST-009`](../../../spec/lastenheft.md) verlangt
offizielle, versionierte Client-Bibliotheken (Packages) für die bestehenden
Zustellwege — HTTP-API (`SPEC-018`), gRPC-Stream (`SPEC-020`), SSE
(`SPEC-021`), NATS-Vollinhalts-Stream (`SPEC-024`). Welche Sprache(n) bedient
werden und über welchen Vertriebsweg, ist ausdrücklich Architektur-/
Spezifikationsfrage (Out-of-Scope-Klausel des Lastenhefts,
[`LH-FA-SST-009.a`](../../../spec/pflichtenheft.md)). Für die **erste**
Sprache hat [`ADR-0106`](0106-csharp-nuget-erstes-sdk-package.md) diese Frage
bereits beantwortet: C#/NuGet, real veröffentlicht als
`PgChangeFeed.Client` (`SPEC-026`, `sdk-csharp-v0.1.0`).

**Die Nutzerentscheidung.** Der Auftraggeber hat am 2026-09-19 im Chat
entschieden: „Können wir noch ein python-SDK erstellen." Das legt die
**Sprache** (Python) für das **zweite** Package fest. Der Vertriebsweg ist
damit **nicht** vorentschieden — er wird in dieser ADR gegen Alternativen
geprüft (§Verglichene Alternativen B), nicht angenommen. Offen bleiben zudem
dieselben Sub-Fragen, die `ADR-0106` für C# getroffen hat: Umfang des ersten
Python-Pakets, Ort im Baum, Versionierung, Bau-/Publish-Mechanismus.

**Der Ist-Stand — gemessen, nicht erinnert.**

| Prozedur | Ergebnis |
|---|---|
| `find examples -maxdepth 1` | `csharp/`, `kotlin/`, die flache Go-Wurzel (`http-client/`, `grpc-client/`, `sse-client/`, `nats-client/`, `nats-stream-client/`) — **kein** `examples/python/`. Anders als bei C# (vier funktionierende Beispiel-Clients, `ADR-0106` §Kontext Ist-Stand Zeile 1) existiert für Python **keine** Draht-Vorarbeit |
| `find sdks -maxdepth 2` | `sdks/csharp/PgChangeFeed.Client/**` — kein `sdks/python/`, keine Namenskollision |
| `grep -rn "sdks/python\|pypi\|PgChangeFeed" spec/ docs/` (sinngemäß über `spec/pflichtenheft.md`) | `SPEC-026` beschreibt ausschließlich das C#-Package; kein bestehender Python-Eintrag |
| `spec/pflichtenheft.md` `SPEC-018` | die vollständige HTTP-API (neun Fähigkeiten, JSON-Schema, Token-Header-Form, Fehler-Antwortform) steht **textuell vollständig** in Klartext-Tabellenform — kein OpenAPI/Swagger-Dokument im Repo (`grep -rin "openapi\|swagger"` ohne Treffer), aber auch keine Lücke: die Tabelle ist selbst der vollständige Vertrag |
| `.a-check.yml` `languages:` | nur `go` — eine Python-Quelle wird von a-check strukturell nicht gelesen, unabhängig vom Ort (dieselbe Grenze wie bei C#, `ADR-0106` §Kontext Bindung 4) |
| `cat docs/user/version.md` | `0.1.2` — Server-SemVer, unverändert eigenständig geführt |

**Was das ändert — der zentrale Unterschied zu `ADR-0106`.** Bei C# senkte
das bereits vorhandene Vier-Wege-Beispiel-Set das Draht-Verständnis-Risiko
für **alle** vier Zustellwege gleichermaßen; der dortige Scope-Schnitt
(HTTP + gRPC in v1) war deshalb ausschließlich eine Frage des
**Paket-Zuschnitts** (Abhängigkeits-Footprint, Erst-Release-Größe), nicht des
Wire-Risikos (`ADR-0106` §Kontext „Was das ändert"). Für Python gilt das
**nicht**: Ohne ein laufendes Referenzprogramm muss das Python-SDK-Team den
Draht-Vertrag für jeden gedeckten Zustellweg **direkt** aus
`spec/pflichtenheft.md` (und im Zweifel dem Go-Server-Code) erschließen, ohne
die Gegenprobe eines bereits funktionierenden Beispiels. Dieses Risiko ist
für die vier Zustellwege **nicht gleich groß**: Die HTTP-API ist ein
zustandsloses REST/JSON-Protokoll, dessen vollständiger Vertrag bereits als
Klartext-Tabelle vorliegt (`SPEC-018`, siehe Ist-Stand oben) — ein Python-Team
kann ihn mit einer stdlib-nahen Bibliothek (`httpx`/`requests`) Zeile für
Zeile nachbauen und mit Unit-Tests gegen die dokumentierten Statuscodes/
Fehlerformen prüfen, ohne eine laufende Gegenprobe zu brauchen. gRPC dagegen
trägt zusätzliche, im Fließtext schwerer zu verifizierende Protokoll-Semantik
(Server-Streaming-Lebenszyklus, Channel-Credentials, Metadata-Header-Form,
generierter Stub aus der `.proto`) — hier senkt das fehlende Vorbild die
Fehler-Erkennungswahrscheinlichkeit beim ersten Schreiben spürbar. Diese
Asymmetrie ist der tragende Grund für einen **kleineren** Erst-Scope als bei
C# (§Entscheidung Festlegung 1) — keine reine Kopie von `ADR-0106`s
Festlegung 1.

**Fünf Bindungen, die den Lösungsraum prägen — analog `ADR-0106`, für Python
neu geprüft:**

1. **Docker-only** ([`AGENTS.md`](../../../AGENTS.md) §3.1): kein
   Host-`python`/`pip`. Bau und Paketierung laufen im gepinnten Container,
   analog `make sdk-pack-csharp`.
2. **`examples/` bleibt Vorbild, kein Belegträger** (`SPEC-023`) — hier ohne
   praktische Wirkung für Python, weil kein `examples/python/` existiert
   (Kontext-Unterschied oben); die Grenze gilt trotzdem für den Fall, dass
   ein künftiger Zug doch Beispiel-Clients unter `examples/python/` anlegt.
3. **Kein neues Gate.** PyPI-Publish braucht Netz und ein Secret — dieselbe
   Klasse wie NuGet-Publish ([`ADR-0106`](0106-csharp-nuget-erstes-sdk-package.md)
   §Entscheidung Festlegung 4) und Docker-Hub-Push
   ([`ADR-0051`](0051-cicd-pipeline-github-actions.md) Entscheidung 8):
   Werkzeug, kein Gate.
4. **a-check liest kein Python.** `sdks/python/**` braucht keine
   `.a-check.yml`-Gruppe — dieselbe Grenze wie `sdks/csharp/**`.
5. **Eine Arbeit, die eine beschriebene Eigenschaft bewegt, zieht ihre Träger
   nach** ([`AGENTS.md`](../../../AGENTS.md) §3.13): `LH-FA-SST-009.a`s Satz
   „Eine zweite Sprache … bleibt offen" wird mit dieser ADR falsch — der
   Träger-Nachzug ist Sache des Implementierungs-Zuges (§Konsequenzen
   Folgepflicht), nicht dieser ausschließlich ADR-schreibenden Rolle.

## Entscheidung

Wir wählen: **Python als zweite SDK-Sprache, PyPI als Vertriebsweg.** Fünf
Festlegungen:

### 1 — Umfang des ersten Pakets: nur HTTP-API, gRPC als Folge-Paket/-Release

Anders als bei C# (`ADR-0106` Festlegung 1: HTTP + gRPC in einem
Erst-Release) deckt v1 des Python-Packages (Arbeitsname `pgchangefeed`)
**ausschließlich HTTP-API** (`SPEC-018`) — **nicht** gRPC (`SPEC-020`), SSE
(`SPEC-021`) oder NATS-Vollinhalt (`SPEC-024`). Begründung:

- **Das fehlende Referenzprogramm ist hier ein echter Grund, kein
  Textbaustein.** Bei C# senkte das Vorbild das Wire-Risiko für alle vier
  Wege gleichermaßen, der Scope-Schnitt dort war reiner Abhängigkeits-/
  Umfangs-Entscheid (`ADR-0106` §Kontext „Was das ändert"). Für Python gibt
  es keine Gegenprobe — das Risiko, eine Protokoll-Annahme beim ersten
  Schreiben falsch zu treffen, ist für gRPC (Streaming-Lebenszyklus,
  generierter Stub, Channel-Credentials) real höher als für ein
  zustandsloses REST/JSON-Protokoll, dessen kompletter Vertrag bereits als
  Klartext-Tabelle vorliegt (`SPEC-018`, §Kontext Ist-Stand).
- **Abhängigkeits-Footprint bleibt zusätzlich ein Grund, wie bei C#.** HTTP
  braucht in Python praktisch keine schwere Fremdabhängigkeit
  (`httpx`, ein etabliertes, reines HTTP-Client-Paket ohne Codegenerierungs-
  Kette). gRPC braucht `grpcio` **und** einen aus der `.proto` generierten
  Stub (`grpcio-tools`/`protoc`) — eine zusätzliche Bau-Stufe, die für ein
  Erst-Release, dessen Draht-Annahmen noch ungeprüft sind, das Fehlerrisiko
  weiter erhöht (ein falscher Stub-Aufruf scheitert oft erst zur Laufzeit
  gegen einen echten Server, nicht beim Bauen).
- **Ein möglichst kleiner, gut geprüfter erster Schnitt senkt das
  Gesamtrisiko eines Erst-Release ohne Vorbild stärker als ein größerer
  Schnitt mit mehr Testfläche gegen ungeprüfte Annahmen.** `SPEC-018` deckt
  bereits neun Fähigkeiten (Consumer-Verwaltung, Tabellen-Verwaltung,
  Retention) plus das Changes-Lesen (`SPEC-022`) — genug Fläche für ein
  eigenständiges, sinnvolles Erst-Release, ohne die zusätzliche
  Stub-Generierungs-Kette.
- **Boundary-Akzeptanzkriterium bleibt erfüllt:**
  [`LH-FA-SST-009`](../../../spec/lastenheft.md) Negative-AC verlangt, dass
  der direkte Zugriff nutzbar bleibt, wo keine Bibliothek existiert — gRPC,
  SSE und NATS-Vollinhalt bleiben für Python-Consumer über die
  entsprechenden Go-Server-Endpunkte direkt ansprechbar (kein
  Referenz-Client nötig, um den Draht selbst zu nutzen — nur um ihn
  komfortabel zu verpacken).
- **Ein Folge-Release trägt gRPC, sobald ein erstes Python-Package real
  existiert und der HTTP-Teil real geprüft ist** — dieselbe inkrementelle
  Logik wie bei C#, nur mit dem HTTP/gRPC-Schnitt eine Stufe früher gezogen,
  weil das Wire-Risiko hier ungleich verteilt ist (§Kontext „Was das
  ändert").

### 2 — Vertriebsweg: PyPI

PyPI ist der Vertriebsweg für das Python-Package, unter dem Arbeitsnamen
`pgchangefeed` (`pip install pgchangefeed`). Begründung, gegen echte
Alternativen abgewogen (§Verglichene Alternativen B):

- **PyPI ist der de-facto-Standardweg für Python-Pakete** — `pip`, das
  Standard-Installationswerkzeug jeder CPython-Distribution, löst
  Abhängigkeiten standardmäßig gegen PyPI auf; ein Consumer braucht keinen
  zusätzlichen Registry-Eintrag oder Kanal, um das Package zu finden.
- **Zielgruppen-Passung.** [`spec/lastenheft.md`](../../../spec/lastenheft.md)
  §Systemkontext nennt „ETL-/ELT-Prozesse, Data-Warehouse-Loader" explizit
  als Consumer-Klasse — genau das Segment, in dem Python (`pandas`,
  Airflow, dbt-Python-Modelle, Data-Pipelines) das dominante Ökosystem ist;
  PyPI ist dort der einzige praktisch erwartete Bezugsweg.
- **Conda scheitert an derselben Kosten-Nutzen-Abwägung wie eine zusätzliche
  Registry:** ein zweiter Vertriebsweg für dieselbe Sprache, bevor auch nur
  ein PyPI-Release real existiert, ist verfrühte Parallelität ohne
  Nutzungsdaten, die sie rechtfertigen (analog `ADR-0106`s B4-Verwerfung).
- **Eine private/interne Registry** hätte keinen erkennbaren Vorteil hier —
  das Package soll von Dritten öffentlich konsumierbar sein
  ([`LH-FA-SST-009`](../../../spec/lastenheft.md) AC: „ein Consumer bindet
  es über den Paketmanager seiner Sprache ein"), eine private Registry
  verlangt zusätzlich eine Zugangsverwaltung, die diese Anforderung nicht
  fordert.
- **„Nur Source, kein Vertriebsweg"** (`git clone` + lokale Installation)
  erfüllt die AC nicht: „über den Paketmanager … einbinden, ohne das
  Protokoll selbst zu implementieren" verlangt einen echten Paketmanager-
  Bezugsweg, kein manuelles Vendoring.

### 3 — Paket-/Repo-Struktur: `sdks/python/`, eigenständiges Projekt

Das Package entsteht unter `sdks/python/` (z. B.
`sdks/python/pgchangefeed/`), analog zu `sdks/csharp/`
([`ADR-0106`](0106-csharp-nuget-erstes-sdk-package.md) §Entscheidung
Festlegung 2). Begründung:

- **`sdks/` ist bereits die etablierte Sprach-Wurzel-Konvention** für
  veröffentlichte SDK-Packages (`grep -rn "sdks/" .` zeigt ausschließlich
  `sdks/csharp/**`) — ein neues Sprach-Verzeichnis unter derselben Wurzel
  ist keine neue Konvention, sondern deren zweite Anwendung.
- **Kein Konflikt mit `examples/`,** weil `examples/python/` ohnehin nicht
  existiert (§Kontext Ist-Stand) — die räumliche Trennungslogik von
  `ADR-0106` §Entscheidung Festlegung 2 (Vorbild vs. Belegträger) gilt hier
  vorsorglich, nicht als Konfliktauflösung eines bestehenden Bestands.
- **a-check bleibt unberührt** (§Kontext Bindung 4): `sdks/python/**` ist
  keine Go-Quelle.
- **Import-Grenze wie bei C#, verschärft:** Das SDK importiert ausschließlich
  die Python-Standardbibliothek und öffentliche PyPI-Pakete, keinen
  privaten Baum dieses Repositories — ein SDK muss außerhalb dieses
  Repositories baubar sein, sobald es veröffentlicht ist.

### 4 — Versionierung: unabhängiges Versionsschema, PEP 440, kompatibel zu SemVer

Das Package führt sein **eigenes** Versionsschema nach PEP 440 — Pythons
offiziellem Versions-Standard —, unabhängig von `docs/user/version.md` (der
Server-Version) **und** unabhängig vom C#-Package-Versionsraum. Begründung:

- **PEP 440 statt reinem SemVer, weil PyPI-Werkzeuge (`pip`, `build`,
  `setuptools`) PEP 440 als kanonisches Schema erwarten** — für den hier
  relevanten einfachen Fall (`MAJOR.MINOR.PATCH`, z. B. `0.1.0`) ist PEP 440
  zu SemVer 2.0 kompatibel; die Boundary-AC von
  [`LH-FA-SST-009`](../../../spec/lastenheft.md) (erkennbare Versionsgrenze,
  „z. B. SemVer-Major") bleibt dadurch unverändert erfüllbar — ein
  PEP-440-Major-Sprung ist derselbe Vorgang wie ein SemVer-Major-Sprung.
- **Unterschiedliche Änderungstakte, dieselbe Begründung wie bei C#**
  (`ADR-0106` §Entscheidung Festlegung 3): ein Python-SDK-Patch braucht
  keinen Server-Release und umgekehrt; ein Python-SDK-Release braucht
  ebenso wenig einen C#-SDK-Release, weil beide Pakete unabhängige
  Draht-Implementierungen mit eigenem Entwicklungstakt sind.
- **Start bei `0.x.y`** (z. B. `0.1.0`), aus denselben Gründen wie beim
  C#-Package — vorstabile öffentliche API, `1.0.0` markiert bewusst
  versprochene Rückwärtskompatibilität.
- **Ein vierter, unabhängiger Versionsraum** neben Server, Harness-Baseline
  und C#-SDK ist dieselbe fortgesetzte Disziplin, keine neue Klasse
  (`ADR-0106` §Entscheidung Festlegung 3, zweiter Punkt).

### 5 — Build-/Publish-Mechanismus: Docker-only bauen, `build`+`twine`, eigener Tag-Namensraum

- **Bauen und Paketieren bleiben Docker-only**
  ([`AGENTS.md`](../../../AGENTS.md) §3.1), analog `make sdk-pack-csharp`:
  ein neues, netzlos **nicht** prüfbares (PyPI-Paketbezug für Test-
  Abhängigkeiten braucht Netz) Werkzeug-Ziel — Arbeitsname
  `make sdk-pack-python` — baut, testet (`pytest`) und paketiert
  (`python -m build`) im gepinnten `python`-Basis-Image, erzeugt ein
  Source-Distribution (`.tar.gz`) und ein Wheel (`.whl`) als Artefakt.
  **Kein Gate** — dieselbe Begründung wie bei `make sdk-pack-csharp`.
- **Paket-Werkzeug: `build` (PyPA-Referenzwerkzeug) + `setuptools`-Backend,
  nicht Poetry oder Hatch.** `build`/`setuptools`/`twine` ist die von der
  offiziellen Python-Packaging-Anleitung (packaging.python.org) empfohlene
  Minimalkette für ein einzelnes, unabhängig veröffentlichtes Package ohne
  Mehrfach-Umgebungs-Verwaltung; Poetry und Hatch sind vollwertige
  Projekt-/Umgebungs-Manager mit eigenem Lockfile-/Environment-Format —
  zusätzliche Werkzeug-Fläche, die dieses schlanke, HTTP-only-Erst-Release
  nicht braucht (§Verglichene Alternativen E).
- **Veröffentlichung nach PyPI braucht einen eigenen, Netz-bindenden
  Workflow und ein Secret** — Arbeitsname `PYPI_API_TOKEN`, dieselbe Klasse
  wie `NUGET_API_KEY` ([`ADR-0106`](0106-csharp-nuget-erstes-sdk-package.md)
  §Entscheidung Festlegung 4) und `DOCKERHUB_USERNAME`/`DOCKERHUB_TOKEN`
  ([`ADR-0051`](0051-cicd-pipeline-github-actions.md) Entscheidung 8): ein
  Repository-Secret, das der Workflow referenziert, nicht anlegt. Der
  Aufruf ist `twine upload --repository pypi dist/* -u __token__ -p
  $PYPI_API_TOKEN`. PyPIs „Trusted Publishing" (OIDC, kein langlebiges
  Secret) ist real verfügbar und der PyPI-eigenen Doku nach der langfristig
  empfohlene Weg — sie bleibt hier bewusst **nicht** die Erstwahl, weil sie
  eine vorherige Projekt-Registrierung auf PyPI voraussetzt (Henne-Ei für
  ein noch nie veröffentlichtes Package) und weil das API-Token-Muster
  bereits zweimal etabliert und geprüft ist (NuGet, Docker Hub) — ein
  konsistentes drittes Mal senkt die Fehlerfläche stärker als der
  Sicherheitsgewinn eines vierten, neuen Mechanismus für ein Erst-Release.
  Eine spätere Umstellung auf Trusted Publishing bleibt möglich, sobald das
  Package einmal real existiert (§Re-Evaluierungs-Trigger).
- **Der Trigger ist ein eigener Tag-Namensraum, weder `v*` noch
  `sdk-csharp-v*`.** `.github/workflows/release.yml` reagiert bereits auf
  `push: tags: ['v*']` für den Server-Container
  ([`ADR-0051`](0051-cicd-pipeline-github-actions.md) Entscheidung 3), ein
  künftiger `sdk-csharp-release.yml`
  ([`ADR-0106`](0106-csharp-nuget-erstes-sdk-package.md) §Entscheidung
  Festlegung 4) auf `sdk-csharp-v*` — dieselbe Präfix-Familie für Python zu
  verwenden würde alle drei Release-Räume kollidieren lassen (Festlegung 4,
  unabhängige Versionierung). Vorschlag für den umsetzenden Zug:
  `sdk-python-v*` als eigener Tag-Präfix, eigener Workflow
  (`sdk-python-release.yml`), der den Tag strikt gegen PEP 440 validiert
  und gegen die im Paket-Metadatenfeld (`pyproject.toml` `[project]
  version`) geführte Version abgleicht — dasselbe Muster wie
  `release.yml`s Tag-gegen-`docs/user/version.md`-Abgleich bzw.
  `sdk-csharp-release.yml`s Tag-gegen-`.csproj`-Abgleich, nur gegen die
  `pyproject.toml`.
- **Beides ist Folgepflicht dieser ADR, nicht Teil ihres Umfangs** — diese
  ADR entscheidet die *Form*, sie implementiert **nichts**
  (§Konsequenzen Folgepflicht).

### 6 — Was diese ADR nicht ändert

- **Kein Produktionscode, kein Eingriff in `internal/**`/`cmd/**`.**
- **`ADR-0106` und `sdks/csharp/**` bleiben unverändert** — kein Umzug,
  keine Zusammenlegung der beiden Sprach-SDKs in dieser ADR.
- **`spec/architecture.md` bleibt unberührt** — `sdks/**` ist ein
  Gate-Scope-Bereich außerhalb des Servers, keine `ARC-*`-Komponente.
- **`docs/user/version.md` bleibt unberührt** — die Server-Version bleibt
  eigenständig geführt (Festlegung 4).
- **Kein neues Gate, keine Schwellen-Senkung.**
- **`.a-check.yml` bleibt unberührt** — a-check liest kein Python.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

### A — Zweite Sprache

| Option | Pro | Contra |
|---|---|---|
| A1 — Go zuerst (Go-Modul-Registry) | spiegelt die Server-Implementierungssprache | bereits in [`ADR-0106`](0106-csharp-nuget-erstes-sdk-package.md) §Verglichene Alternativen A1 verworfen — Go-Consumer haben die geringste Einstiegshürde über die rohe API, der SDK-Nutzen ist hier am kleinsten, gerade wo er am billigsten wäre; keine bestehende Nutzerentscheidung dafür (weder damals noch heute) |
| A2 — TypeScript/npm | breites Web-/Node-Ökosystem, viele potenzielle Consumer | keine bestehende Beispiel-Client-Vorarbeit (kein `examples/typescript/`) — dasselbe Risiko wie bei Python (§Kontext), aber ohne die Zielgruppen-Passung: Web-/Node-Anwendungen sind seltener der CDC-Konsument als Data-Pipelines; keine bestehende Nutzerentscheidung dafür |
| A3 — Kotlin/Maven | bereits vier funktionierende Kotlin-Beispiel-Clients unter `examples/kotlin/` (`ADR-0087`/`ADR-0090`) — dasselbe niedrige Wire-Risiko wie C# bei Erstwahl | Kotlin/JVM-Consumer haben mit Java/Maven bereits ein etabliertes, sprachnahes Ökosystem für HTTP-/gRPC-Clients (z. B. generierte gRPC-Java-Stubs) — der SDK-Mehrwert ist geringer als für ein Ökosystem ohne vergleichbare Bordmittel; keine bestehende Nutzerentscheidung dafür |
| **A4 — Python/PyPI (gewählt)** | explizite Nutzerentscheidung; stärkste Zielgruppen-Passung — [`spec/lastenheft.md`](../../../spec/lastenheft.md) §Systemkontext nennt „ETL-/ELT-Prozesse, Data-Warehouse-Loader" ausdrücklich, genau das Python-dominierte Segment; PyPI ist ein etablierter, gut automatisierbarer Vertriebsweg | keine bestehende Beispiel-Client-Vorarbeit (kein `examples/python/`) — höheres Wire-Risiko als bei C#/Kotlin, adressiert durch einen kleineren Erst-Scope (§Entscheidung Festlegung 1) |
| A5 — nichts tun / abwarten | kein Aufwand jetzt | die explizite Nutzerentscheidung vom 2026-09-19 würde ignoriert; `LH-FA-SST-009.a`s „zweite Sprache offen" bliebe unbeantwortet, obwohl die Zielgruppen-Passung (ETL/ELT) bereits im Lastenheft dokumentiert ist |

### B — Vertriebsweg für Python

| Option | Pro | Contra |
|---|---|---|
| **B1 — PyPI (gewählt)** | de-facto-Standardweg, `pip`-Default-Index, keine zusätzliche Registry-Konfiguration für den Consumer; passt zur Zielgruppe (ETL-/Data-Pipeline-Tooling installiert praktisch ausschließlich über PyPI) | Publish braucht Netz + eigenes Secret (Folgepflicht) |
| B2 — Conda(-Forge) | in Data-Science-/Data-Engineering-Umgebungen verbreitet, löst native Abhängigkeiten mit | zusätzlicher, paralleler Vertriebsweg für dieselbe Sprache ohne Nutzungsdaten, die ihn rechtfertigen (analog `ADR-0106`s B4-Verwerfung „verfrühte Modularisierung"); höherer Pflegeaufwand (eigenes Feedstock-Repo, eigener Review-Prozess bei conda-forge) für ein reines HTTP-Package ohne native Erweiterungen |
| B3 — private/interne Package-Registry | volle Kontrolle über Zugriff | widerspricht der AC „Consumer bindet über den Paketmanager seiner Sprache ein" ohne fachlichen Grund — das Package soll öffentlich, nicht zugriffsbeschränkt konsumierbar sein |
| B4 — nur Source (kein Paketmanager-Vertriebsweg) | kein Publish-Mechanismus, kein Secret nötig | erfüllt die AC nicht: „über den Paketmanager … einbinden, ohne das Protokoll selbst zu implementieren" verlangt einen echten Bezugsweg, nicht manuelles Vendoring |

### C — Umfang des ersten Pakets

| Option | Pro | Contra |
|---|---|---|
| **C1 — nur HTTP-API (gewählt)** | kleinster, am besten prüfbarer Erst-Scope ohne Referenz-Client; `SPEC-018` liegt bereits vollständig als Klartext-Vertrag vor, keine Codegenerierungs-Kette nötig; deckt bereits neun Fähigkeiten plus Changes-Lesen | Live-Zustellung (gRPC/SSE/NATS) wird von v1 nicht abgedeckt |
| C2 — HTTP-API + gRPC-Stream, ein Package (wie bei C#) | Konsistenz zum C#-Zuschnitt, deckt Abruf **und** Live-Zustellung sofort | ignoriert die reale Risiko-Asymmetrie: ohne Referenz-Client ist die gRPC-Streaming-Semantik beim ersten Schreiben ungleich schwerer zu verifizieren als bei C# (§Kontext „Was das ändert"); zusätzliche Stub-Generierungs-Kette (`grpcio-tools`) im allerersten Release |
| C3 — alle vier Zustellwege in einem Package | „vollständige Matrix sofort" | größtmöglicher Erst-Release-Umfang bei größtem Draht-Risiko (kein Vorbild für keinen der vier Wege); überträgt den bei C# bereits für B3 verworfenen Footprint-Einwand verschärft auf eine Sprache ohne jede Referenzimplementierung |
| C4 — separates Package je Zustellweg von Anfang an | jeder Consumer lädt nur, was er braucht | vier Package-Namen, vier Versionsräume, vier Publish-Workflows für eine Sprache, die noch nicht einmal ein erstes Release hat — verfrühte Modularisierung ohne Nutzungsdaten (analog `ADR-0106`s B4) |

### D — Versionierung

| Option | Pro | Contra |
|---|---|---|
| D1 — gekoppelt an `docs/user/version.md` oder an das C#-SDK | ein Blick genügt | erzwingt SDK-Releases bei jedem Server- bzw. C#-SDK-Release ohne eigenen Draht-Änderungsgrund; passt nicht zur Boundary-AC (SDK-eigene Versionsgrenze) |
| **D2 — eigenständiges PEP-440-Schema, ab `0.x.y` (gewählt)** | Release-Takt folgt der tatsächlichen Python-SDK-Änderung; PEP 440 ist das von PyPI-Werkzeugen erwartete Schema und für einfache Fälle SemVer-kompatibel; Analogie zu bereits drei bestehenden unabhängigen Versionsräumen (Server, Harness-Baseline, C#-SDK) | ein vierter Versionsraum, den ein Betrachter separat nachschlagen muss |
| D3 — Versionierung nach Datum (CalVer) | keine Debatte über Major/Minor/Patch | passt nicht zur AC-Boundary („z. B. SemVer-Major"); unüblich für PyPI-Bibliotheken, die von anderen Paketen als Abhängigkeit eingebunden werden |

### E — Paket-Werkzeug (Python)

| Option | Pro | Contra |
|---|---|---|
| **E1 — `build` (PyPA) + `setuptools`-Backend + `twine` (gewählt)** | von packaging.python.org als Minimalweg für ein einzelnes Package empfohlen; keine zusätzliche Projekt-/Umgebungs-Verwaltungsschicht; direkt analog zu `dotnet pack`/`dotnet nuget push`s Einfachheit bei C# | mehr manuelle `pyproject.toml`-Pflege als bei einem All-in-one-Werkzeug |
| E2 — Poetry | beliebtes All-in-one-Werkzeug (Dependency-Resolution, Build, Publish in einem) | eigenes Lockfile-/Environment-Format, zusätzliche Werkzeug-Fläche, die ein schlankes HTTP-only-Erst-Release ohne komplexe Abhängigkeitsauflösung nicht braucht |
| E3 — Hatch | offiziell von der PyPA als Alternative geführt, moderne Plugin-Architektur | ähnliches Mehrgewicht wie Poetry für den hier benötigten einfachen Fall; ein weiteres, noch unbewährtes Werkzeug in diesem Repo statt einer bereits von PyPA als Minimalweg dokumentierten Kette |

**Fazit:** A4, B1, C1, D2, E1. A1–A3 scheitern an fehlendem Nutzen, fehlender
Zielgruppen-Passung bzw. bereits vorhandenen sprachnahen Bordmitteln, A5 an
der unbeantworteten Nutzerentscheidung; B2–B4 passen nicht zur AC bzw.
modularisieren verfrüht ohne Nutzungsdaten; C2–C4 ignorieren die
Risiko-Asymmetrie eines fehlenden Referenz-Clients bzw. überziehen den
Erst-Release; D1/D3 passen nicht zur Boundary-AC bzw. zum Ökosystem; E2/E3
bringen Werkzeug-Gewicht, das dieser schlanke Erst-Scope nicht braucht.

## Konsequenzen

- **Positiv:** `LH-FA-SST-009.a`s zweite offene Frage bekommt eine konkrete,
  begründete Antwort inklusive Vertriebsweg-Abwägung statt einer
  Übernahme der Nutzeräußerung ohne Prüfung; die Zielgruppen-Passung
  (ETL-/Data-Warehouse-Consumer, `spec/lastenheft.md` §Systemkontext) senkt
  das Risiko, ein SDK für ein Segment zu bauen, das es nicht wirklich
  nachfragt; der bewusst kleinere Erst-Scope (nur HTTP) hält das
  Draht-Risiko ohne Referenzprogramm handhabbar.
- **Negativ:** gRPC-, SSE- und NATS-Consumer in Python bleiben vorerst auf
  den direkten Zugriffsweg verwiesen (durch die Negative-AC von
  [`LH-FA-SST-009`](../../../spec/lastenheft.md) gedeckt, aber ein realer
  Komfortverlust, größer als bei C#, wo immerhin gRPC bereits in v1
  enthalten ist); ein dritter Release-/Secret-Mechanismus (PyPI) entsteht
  neben den bestehenden (Docker Hub/GHCR, NuGet); PEP-440-`SemVer`-
  Kompatibilität gilt nur für den hier genutzten einfachen Fall
  (`MAJOR.MINOR.PATCH`), nicht für PEP 440s vollen Funktionsumfang
  (Epochs, Pre-/Post-Releases) — bewusst ungenutzt, kein Verlust hier.
- **Folgepflicht:**
  1. Ein Planner-Zug schneidet die Umsetzung (voraussichtlich mehrere
     Slices: SDK-Projektgerüst, HTTP-Client-Fläche, Pack-Werkzeug,
     Publish-Workflow) — diese ADR entscheidet die Form, nicht den Schnitt.
  2. `spec/pflichtenheft.md` §1 (`LH-FA-SST-009.a`) und §6 (Externe
     Verträge, analog `SPEC-026`) brauchen einen Träger-Nachzug, sobald das
     Package real existiert (`AGENTS.md` §3.13) — Sache des umsetzenden
     Zuges, nicht dieser ausschließlich ADR-schreibenden Rolle.
  3. `harness/README.md` §Werkzeuge bekommt seine Zeilen für
     `make sdk-pack-python` und den Publish-Workflow erst, wenn sie real
     existieren (`AGENTS.md` §4: kein Träger nennt ein Target, das es nicht
     gibt).
  4. `docs/user/benutzerhandbuch.md` bekommt einen SDK-Hinweis, im selben
     Zug wie das jeweilige Client-Programm.
  5. Der PyPI-Publish-Workflow gilt nach [`AGENTS.md`](../../../AGENTS.md)
     §3.10 erst nach einem realen, grünen Post-Push-Lauf als abgeschlossen.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `a-check` | keine — `sdks/python/**` ist keine Go-Quelle, a-check liest sie nicht (`.a-check.yml` `languages: go`) | — |
| Review (kein Sensor) | `sdks/python/**` importiert ausschließlich die Python-Standardbibliothek und öffentliche PyPI-Pakete, keinen privaten Baum dieses Repos (Import-Grenze, analog `ADR-0106` §Entscheidung Festlegung 2) | — |
| künftiges Bau-Werkzeug (Folgepflicht) | `make sdk-pack-python` baut/testet/paketiert Docker-only, netzlos geprüfte Teile fahren im Bau, kein Gate | `make sdk-pack-python` (noch nicht existent) |

## Re-Evaluierungs-Trigger

1. **Eine dritte Sprache oder ein dritter Vertriebsweg wird für
   `LH-FA-SST-009` verlangt** (z. B. TypeScript/npm, Kotlin/Maven) — eigene
   ADR, diese ADR bleibt für Python/PyPI unberührt (dieselbe Regel wie
   `ADR-0106` §Re-Evaluierungs-Trigger 1 sie für C# formuliert).
2. **Das Python-Package wurde real auf PyPI veröffentlicht und
   Nutzungsdaten (Download-Zahlen, Issue-Nachfrage) legen eine andere
   Priorisierung nahe** — insbesondere ob gRPC als nächstes Folge-Release
   sinnvoller ist als SSE/NATS oder umgekehrt (Festlegung 1 offen gelassen).
3. **Ein Python-Referenz-Client entsteht nachträglich** (z. B. weil ein
   künftiger Zug doch `examples/python/` anlegt) — dann verliert die in
   §Kontext „Was das ändert" begründete Risiko-Asymmetrie ihre Grundlage,
   und ein größerer Erst-Scope für ein *Folge*-Release ließe sich ohne die
   hier geltende Vorsicht rechtfertigen.
4. **PyPIs Trusted Publishing wird für ein bereits real existierendes
   Package geprüft** — Festlegung 5 verwarf es nur für das *Erst-Release*
   wegen der Henne-Ei-Registrierungsvoraussetzung; nach dem ersten
   erfolgreichen `PYPI_API_TOKEN`-Publish ist eine Umstellung ein sinnvoller
   Anlass zur Neubewertung, kein struktureller Ausschluss.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-19 | Proposed | dieser Architect-Zug (kein vorausgehender Slice-Plan — ADR vor Implementierung, Modul 8) |
| 2026-09-19 | Accepted | Nutzerentscheidung im Chat vom 2026-09-19 („Können wir noch ein python-SDK erstellen") + diese vollwertige ADR-Fassung samt Alternativenvergleich (Sprache **und** Vertriebsweg) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0107` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
