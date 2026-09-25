# Slice sdk-readme-nutzerdoku: Die drei SDK-READMEs und Paket-Metadaten als Anwender-Dokumentation — Aufbau, belegte Beispiele, API-Übersicht; Version 0.2.1

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — es gibt keine Closure-Bedingung, die von der DoD dieses
Slice verschieden ist (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine
Welle braucht).

**Bezug:** [`LH-FA-SST-009`](../../../../spec/lastenheft.md) (offizielle Client-Bibliotheken — die drei
Packages, deren Paketbeschreibung dieser Slice ist),
[`ADR-0106`](../../adr/0106-csharp-nuget-erstes-sdk-package.md), [`ADR-0107`](../../adr/0107-python-pypi-zweites-sdk-package.md), [`ADR-0109`](../../adr/0109-kotlin-github-packages-drittes-sdk-package.md)
(die drei Package-Entscheidungen: Vertriebsweg, Metadaten-Quellen, Version unabhängig
vom Server), [`ADR-0110`](../../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md) (Sprachmatrix, Flächen der Packages).

**Berührte Spec-Stellen:** [`SPEC-026`](../../../../spec/pflichtenheft.md), [`SPEC-027`](../../../../spec/pflichtenheft.md), [`SPEC-028`](../../../../spec/pflichtenheft.md)
(Package-Zeilen mit der aktuellen Version) und §1 `LH-FA-SST-009.a` (real erzeugte
Artefaktnamen) — geändert; [`SPEC-018`](../../../../spec/pflichtenheft.md) bis [`SPEC-024`](../../../../spec/pflichtenheft.md) (Drahtverträge, aus denen die
Anwender-Beschreibung ihre Aussagen liest) — gelesen.

**Verantwortlich:** Implementer-Agent, 2026-09-25.

**Autor:** Implementer-Agent (Nutzer-Anweisung „ohne Nachfragen weiter“, Rückmeldung zur
PyPI-Paketseite). **Datum:** 2026-09-25.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die README jedes der drei SDK-Packages beschreibt in Anwender-Sprache, was
das Package tut, wie man es installiert, einsetzt und welche Klassen und Methoden es
bietet — belegt durch Beispiele aus dem Quelltext und den Tests der SDKs, ohne
interne Kennungen und ohne Projekt-Chronik; die Metadaten-Felder, die Anwender auf den
Paket-Seiten sehen, sind ebenso frei von Kennungen; die Quell-Version aller drei
Packages steht auf `0.2.1`, damit die Beschreibung mit dem nächsten Release erscheint.

**Ausgangslage (Nutzer-Rückmeldung, wörtlich):** (1) Der PyPI-Status-Text sei „nicht
aktuell/sinnvoll“ (er nennt, es gebe noch keinen Referenz-Client, und zählt Kennungen
auf); (2) „Warum steht dort nicht die API Beschreibung“ — die README enthält keine
API-Beschreibung, nur Verweise auf die Spezifikation; (3) mit den `SPEC-`-Kennungen
könne kein Benutzer etwas anfangen. Die README ist die Paketbeschreibung auf PyPI.org
(`pyproject.toml` `readme`) und auf NuGet.org (`PackageReadmeFile`); die Kotlin-README
gelangt nicht in das Artefakt (§3).

**Umfangserweiterung der Fixrunde (Nutzer-Anweisung „Ja, öffentliche API-Kommentare
auch bereinigen -- Auch keine Chronik/Forensik in den Kommentaren betreiben. Kommentare
beschreiben die Funktionen/Module.“):** Der Slice umfasst zusätzlich die
Anwender-Dokumentation **im Code**: die öffentlichen API-Kommentare der drei SDKs
(Python-Docstrings, C#-XML-Kommentare, Kotlin-KDoc, jeweils auch Modul-/Typkommentare),
interne Kommentare und Testkommentare der SDK-Dateien, die Kommentare der Build-Dateien
und die Laufzeit-Fehlertexte. Regel: ein Kommentar beschreibt, was die Funktion, das
Modul oder der Typ tut (Zweck, Parameter, Rückgabe, Fehler, Besonderheiten), in Worten
eines Anwenders — ohne Kennungen (`SPEC-`/`ADR-`/`LH-`/`ARC-`), ohne Slice-/Welle-/
Paragraphen-Bezüge und ohne Chronik oder Abwägung verworfener Alternativen.
*Begründung der Aufnahme:* Docstrings und Fehlertexte sind Teil des ausgelieferten
Python-Pakets (Wheel und sdist), die XML-Kommentare gehören nach der Aufnahme von
`GenerateDocumentationFile` zum NuGet-Paket, die KDoc zum Sources-Jar; PyPI und NuGet
erlauben keinen erneuten Upload einer Version — was in `0.2.1` steht, bleibt dort, jede
spätere Bereinigung bräuchte `0.2.2`. Deshalb gehört die Bereinigung in `0.2.1`.
Die Signaturen und das Verhalten der SDKs bleiben unberührt; geändert ist der Wortlaut
der Fehlertexte (20 Texte, Python 8, C# 6, Kotlin 6).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Ein Release / ein Tag** — ein Tag-Push (`sdk-*-v*`) ist Betreiber-Handlung
  ([`AGENTS.md`](../../../../AGENTS.md) §3.10); die drei `sdk-*-release.yml`-Workflows bleiben unverändert. Der
  Slice hebt nur die Quell-Version, damit ein späterer Tag eine noch nicht
  veröffentlichte Version trägt.
- **Änderungen an Signaturen oder Verhalten der SDKs** — die Beschreibung folgt dem
  Code, nicht umgekehrt. Auffälligkeiten im Code, die eine README-Aussage
  erschweren (z. B. transitive Abhängigkeiten der Kotlin-Flächen), werden gemeldet,
  nicht mitgeändert. Ausnahme der Fixrunde: der Wortlaut der Laufzeit-Fehlertexte
  (siehe Umfangserweiterung).
- **Das Benutzerhandbuch** (`docs/user/benutzerhandbuch.md`) — ein anderes Dokument
  mit eigener Konvention für Betreiber und Integratoren; es verweist auf die READMEs,
  ohne Aussagen zu ihrem Inhalt oder Paketstand zu machen (Suchlauf §3). Ändert sich
  daran nichts, entfällt auch Version und Historie.
- **Dauerhafter Wächter für die C#- und Kotlin-Beispiele** — die Beispiele werden
  einmalig durch Übersetzen belegt (§3). Für Python läuft der Wächter als Test im
  Package-Bau, weil die README dort bereits im Bau-Kontext liegt; für C# und Kotlin
  läge die README außerhalb der Test-Bau-Kontexte, ein Wächter dort wäre eine
  Änderung der Bau-Kontexte (ein anderer Vorgang).
- **Ein neues Gate in `make gates`** für den Kennungs-Wächter der SDK-Dateien — jedes
  bestehende Gate trägt eine ADR
  ([`ADR-0041`](../../adr/0041-a-check-maschinenform-architekturpruefung.md),
  [`ADR-0045`](../../adr/0045-commit-traceability-standing-gate.md),
  [`ADR-0054`](../../adr/0054-coverage-gate-und-benchmark-infrastruktur.md),
  [`ADR-0084`](../../adr/0084-sync-gate-fuer-generierte-artefakte.md)), und die ADR ist
  Architect-Arbeit. Die Fixrunde liefert den Wächter als Werkzeug mit Tabellentest, das die
  drei `make sdk-pack-*`-Ziele als Vorstufe tragen und das die drei Release-Workflows über
  diese Ziele vor jedem Publish fahren (§3); ein Rückfall wird damit nicht ausgeliefert,
  fällt aber nicht bei jedem Push auf. Die Aufnahme in `make gates` ist übergeben an
  `slice-sdk-public-doc-check-gate` (netzlos, gemessen `real 0m0,013s` für
  `make sdk-public-doc-check`, 2026-09-25); der Slice schwächt keine bestehende Schwelle.

## 2. Definition of Done

- [x] Die drei READMEs tragen den Anwender-Aufbau (Was ist das · Installation · Schnellstart ·
      Live-Streams · API-Übersicht · Change-Objekt mit `origin` · Fehlerbehandlung ·
      Zustellwege im Vergleich · Links · Lizenz), jedes Code-Beispiel trägt einen
      Beleg-Anker in §3 und ist übersetzt (C#, Kotlin) bzw. durch den Wächter im Bau
      geprüft (Python); keine `SPEC-`/`ADR-`/`LH-`/`ARC-`-Kennung, keine Chronik.
- [x] Die Paket-Metadaten-Felder (Beschreibung, Schlagwörter, URLs) tragen keine Kennung; die
      Quell-Version der drei Packages ist `0.2.1`, alle Träger, die sie nennen, stehen
      auf demselben Stand (Suchlauf §3, beide Stände gemessen).
- [x] `make sdk-pack-csharp`, `make sdk-pack-python`, `make sdk-pack-kotlin` grün, die
      Artefakte tragen die neue README bzw. Metadaten (Listing im Bericht),
      `make docs-check` und `make gates` grün.
- [x] Öffentliche API-Kommentare, interne Kommentare, Testkommentare, Build-Datei-Kommentare
      und Laufzeit-Fehlertexte der drei SDKs tragen keine Kennung und beschreiben den
      Ist-Zustand; ein Wächter hält es (Python: `test_public_text.py` im Bau,
      alle drei: `make sdk-public-doc-check` als Vorstufe von `make sdk-pack-*`) —
      Beleg: Suchlauf §3, Mutationen §3.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8); kein
      offenes HIGH/MEDIUM nach der Fixrunde (F-1 bis F-3 behoben, Belege §3);
      Verifikation: [`verifikation-slice-sdk-readme-nutzerdoku`](../../../reviews/verifikation-slice-sdk-readme-nutzerdoku.md)
      (Verdikt Bestätigt).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung
      angefallen“ in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **mit**
      Wellen von der nächsten Welle-Closure geprüft (auch für Slices ohne
      Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Ist-Zustand gemessen am Parent (`6839f738`): Die drei READMEs tragen einen
„Status“-Absatz mit Kennungen (`SPEC-`/`ADR-`/`LH-`) und den Satz „no `examples/python/`
reference client exists yet“ (Python), verweisen für die API auf `spec/pflichtenheft.md`
und haben keine API-Übersicht, keine Beispiele (Ausnahme: der Kotlin-Installationsweg).
Wohin die README gelangt: Python `pyproject.toml` `readme = "README.md"` (der Docker-Bau
kopiert `sdks/python/README.md` an `pgchangefeed/README.md`); C# `PackageReadmeFile` (das
`.csproj` packt `../README.md` an die Paket-Wurzel); Kotlin: das `jar` trägt keine README,
der `publishing`-Block der `build.gradle.kts` trägt keine `pom`-Beschreibung — Anwender
sehen dort nur Koordinate und Version.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `sdks/python/README.md` | update | Anwender-Aufbau; Beispiele aus `sdks/python/pgchangefeed/integration/test_*_realserver.py` und `src/pgchangefeed/*.py` (Beleg-Anker unten je Beispiel). |
| `sdks/csharp/README.md` | update | dito, Beispiele aus `sdks/csharp/PgChangeFeed.Client.Integration/*RealserverTests.cs`. |
| `sdks/kotlin/pgchangefeed-kotlin/README.md` | update | dito, Beispiele aus `src/integrationTest/kotlin/…/*RealserverTest.kt`; der Installationsweg (GitHub Packages, Zugangsdaten) bleibt in Anwender-Sprache. |
| `sdks/python/pgchangefeed/pyproject.toml` | update | `keywords`, `[project.urls]` (Dokumentation, Fehlerberichte), `version = "0.2.1"`. |
| `sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.csproj` | update | `<PackageTags>`, `<Version>0.2.1</Version>`, Versions-Kommentar. |
| `sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts` | update | `pom { name, description, url, licenses }` im `publishing`-Block (das Feld, das GitHub Packages anzeigt), `version` `0.2.1` an beiden Stellen (Kopf und `publishing`-Block), Versions-Kommentar. |
| `sdks/python/pgchangefeed/tests/test_readme_examples.py` | neu, in der Fixrunde erweitert | Wächter: jeder ` ```python `-Block der README übersetzt (`compile`), jeder `from pgchangefeed… import`-Name löst auf; **jeder Aufruf** eines Beispiels (Client-Methode, Konstruktor, Funktion eines importierten Moduls) bindet gegen die echte Signatur (`inspect.signature(...).bind`, unbekanntes Schlüsselwort oder fehlendes Pflichtargument färbt rot); jede öffentliche Client-Methode und jeder Stream-Client hat eine Zeile der API-Tabelle mit den Parametern des Codes (Konstruktoren eingeschlossen; Parameter per `ast` gelesen, nicht per Kommagrenze); die Feldliste der Stream-Nachricht gleicht `StreamChange` **und** dem generierten gRPC-`Change`; die Fehler-Tabelle nennt jede Exception des Pakets mit dem Statuscode, auf den der Client sie abbildet. Feste Zähler (`13`, `7`, `>= 5`) sind durch aus der API abgeleitete Mengen ersetzt. |
| `sdks/python/pgchangefeed/tests/test_public_text.py` | neu (Fixrunde) | Wächter gegen Rückfall: keine Datei des installierten Pakets (`pgchangefeed/**/*.py`, einschließlich der zur Bauzeit erzeugten Stubs) trägt eine Kennung `SPEC-`/`ADR-`/`ARC-`/`LH-FA-`/`LH-QA-`. Läuft im Bau vor `pack`. |
| `sdks/python/pgchangefeed/src/pgchangefeed/*.py` (8 Dateien) | update (Fixrunde) | Docstrings, Modul- und `#`-Kommentare und acht Laufzeit-Fehlertexte beschreiben, was Klassen und Funktionen tun; öffentliche Methoden tragen Docstrings mit Verhalten und Fehlern (`register_consumer`, `acknowledge_consumer`, …). |
| `sdks/python/pgchangefeed/tests/*.py`, `sdks/python/pgchangefeed/integration/*.py` (8 Dateien), `pgchangefeed/pyproject.toml` | update (Fixrunde) | Kennungen und Slice-/ADR-Bezüge aus Testkommentaren und Build-Kommentaren entfernt; Testnamen `…spec_020…`/`…broadcaster…` sachlich umbenannt; die sdist trägt Tests und `pyproject.toml`. |
| `proto/cdc/stream/v1/changestream.proto`, `gen/cdc/stream/v1/changestream.pb.go`, `gen/cdc/stream/v1/changestream_grpc.pb.go` | update (Fixrunde) | Die Kommentare der `.proto` landen in den erzeugten Stubs aller Sprachen (im Python-Wheel: `changestream_pb2_grpc.py`, gemessen: 4 Zeilen mit Kennung). Neue englische Anwender-Kommentare ohne Kennung; Draht und Feldnummern unverändert; der Go-Code ist mit `make proto-generate` neu erzeugt (nur Kommentare, +20/−30 Zeilen). |
| `sdks/csharp/PgChangeFeed.Client/**/*.cs` (ohne Tests), `PgChangeFeed.Client.csproj`, `Directory.Packages.props` (zusammen 17 Dateien mit Kennung vor der Fixrunde) | update (Fixrunde) | XML-Kommentare, Typ-/Modulkommentare und sechs Laufzeit-Fehlertexte; der h2c-Kommentar von `PgChangeFeedGrpcClient` stimmt jetzt mit README und Realserver-Test überein (kein `AppContext`-Schalter nötig, `http://` genügt). Konstruktoren und `ReadFrameAsync` tragen `<summary>`, ein `cref` ist aufgelöst (`global::Grpc.Core.RpcException`), der Bau ist warnungsfrei. `<GenerateDocumentationFile>true</GenerateDocumentationFile>` liefert die Kommentare als `lib/net10.0/PgChangeFeed.Client.xml` im nupkg aus (F-10). |
| `sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.Tests/**`, `sdks/csharp/PgChangeFeed.Client.Integration/*.cs` | update (Fixrunde) | Kennungen aus Testkommentaren entfernt; Testnamen `AllTenSpec0NNFieldsRoundTrip` zu `AllTenFieldsRoundTrip`, `NoBroadcasterWired_…` zu `StreamNotAvailable_…`. |
| `sdks/kotlin/pgchangefeed-kotlin/src/main/**/*.kt` (17 Dateien), `build.gradle.kts`, `settings.gradle.kts` | update (Fixrunde) | KDoc, Modulkommentare und sechs Laufzeit-Fehlertexte; die Kommentare in `build.gradle.kts` beschreiben den Ist-Zustand der POM-Felder ohne Aussage über die Anzeige auf GitHub Packages (F-6). `java { withSourcesJar() }` liefert die Quellen samt KDoc als `pgchangefeed-kotlin-0.2.1-sources.jar` neben dem Haupt-Jar aus (F-10); die Publikation verweist darauf (gemessen mit `generateMetadataFileForMavenPublication`). |
| `sdks/kotlin/pgchangefeed-kotlin/src/test/**`, `src/integrationTest/**` | update (Fixrunde) | Kennungen und Slice-Bezüge aus Testkommentaren entfernt; Testnamen `…all SPEC-0NN fields…` zu `…all ten fields…`, `no Broadcaster wired…` zu `stream not available…`; die „Ausgang bewusst offen (Slice-Plan §6)“-Sätze der Fakes beschreiben, was der Fake nachbildet. |
| `sdks/csharp/Dockerfile`, `sdks/kotlin/Dockerfile`, `sdks/python/Dockerfile`, `sdks/csharp/.gitignore`, `sdks/kotlin/.gitignore` | update (Fixrunde) | Kommentare als Beschreibung des Ist-Zustands ohne Kennungen (nur Kommentarzeilen; die Kotlin-`pack`-Stufe nennt das Sources-Jar). Der Wächter-Umfang ist damit einfach: keine Kennung unter `sdks/`. |
| `tools/harness/sdk-public-doc-check.sh`, `tools/harness/run-sdk-public-doc-check-tests.sh`, `harness/mk/sdk.mk` | neu / update (Fixrunde) | `make sdk-public-doc-check` (grep, netzlos) prüft alle Textdateien unter `sdks/` ohne Bau-Ausgaben und erzeugten Code; die drei `sdk-pack-*`-Ziele hängen davon ab. `make test-sdk-public-doc-check` fährt den Tabellentest. Begründung, warum kein Gate: §1. |
| `sdks/*/README.md` | update (Fixrunde) | F-2: Kotlin-`oldImage`/`newImage` sind bei INSERT/DELETE `JsonNull`, nicht `null`; F-5: Transportfehler werden nicht umgesetzt (je Sprache benannt, Kotlin: zehn Sekunden Timeout je HTTP-Aufruf); F-7: `enable_table`-Felder erklärt, „Kotlin or Java“ zu „Kotlin“; F-9/F-12: `kotlin-stdlib` ist `compile`, TLS-Aussage (der Server bedient kein TLS, Belegbefehl §3), gRPC-Bilder „empty when there is none“ und der Weg zu einem eigenen Kanal. |
| `harness/README.md` (§Sensors), `docs/user/releasing.md` | update (Fixrunde) | Zeilen für die neuen Ziele, Sources-Jar und XML-Dokumentationsdatei; die drei SDK-Workflow-Zeilen und `releasing.md` (Version 1.7) auf den Ist-Stand der veröffentlichten Versionen (F-4). |
| `spec/pflichtenheft.md` §1 (`LH-FA-SST-009.a`, drei Artefakt-Sätze), §6 `SPEC-026`/`-027`/`-028`, §7 | update | die Zeilen tragen die aktuelle Version bzw. den real erzeugten Artefaktnamen; eine Historie-Zeile. |
| `harness/README.md` §Sensors (`make sdk-pack-*`), `harness/mk/sdk.mk`, `sdks/csharp/Dockerfile`, `sdks/kotlin/Dockerfile` | update | Träger, die den Artefaktnamen mit der Version nennen (Suchlauf unten). |
| `docs/user/benutzerhandbuch.md` | entfällt (kein Diff) | Die zwölf Verweise auf die READMEs (12 Zeilen, alle mit `sdks/`-Pfad) (`git grep -n -i 'README' docs/user/benutzerhandbuch.md`) nennen nur den Pfad; die Aussage zum Kotlin-Installationsweg (Zeile 1054, „vollständigen Installationsweg samt Gradle-Zugangsdaten-Konfiguration“) gilt weiter, die README-Installationssektion trägt beides. Keine Aussage zu README-Inhalt oder Paketstand → weder `Version:` noch Historie berührt. |
| `sdks/csharp/Dockerfile`, `sdks/kotlin/Dockerfile`, `harness/mk/sdk.mk` (Kommentare mit Artefaktnamen) | update | reine Namensnachführung `0.2.0` → `0.2.1` in Kommentaren. |
| `sdks/*/` `ReadmeSnippetsTemp` (C#-Testprojekt, Kotlin-Testquellmenge) | einmalig, nicht geliefert | Die README-Beispiele wurden je einmal in einer temporären Testdatei des jeweiligen Test-Bau-Kontexts übersetzt (`make sdk-pack-csharp`/`-kotlin` grün, danach mit einer Falsch-Mutation rot gesehen, Datei entfernt); die Dateien stehen nicht im Diff. |

**Beleg-Anker der Code-Beispiele** (je Beispiel die Quelle im Repo, aus der Aufrufform,
Typen und Namen gelesen sind; das Beispiel ist daraus zusammengesetzt, Werte sind
Platzhalter). Übersetzungs-Beleg: Python durch den Wächter `test_readme_examples.py` im
Bau (`make sdk-pack-python`, gedruckte Zeile `64 passed, 2 warnings in 0.51s`, 2026-09-25); C# und Kotlin durch
die einmalige Übersetzung in der temporären Testdatei (`make sdk-pack-csharp`: `Passed!  -
Failed:     0, Passed:    76`; `make sdk-pack-kotlin`: `BUILD SUCCESSFUL in 28s`; je mit
einer Falsch-Mutation `StreamChangeAsync`/`streamChange` rot gesehen: `error CS1061 …
'StreamChangeAsync'`, `e: … Unresolved reference 'streamChange'`).

| Sprache · Beispiel | Beleg-Anker (Quelle im Repo) |
|---|---|
| Python · HTTP-Schnellstart (`register_consumer`, `get_consumer_position`, `read_changes(from_=)`, `acknowledge_consumer`) | `sdks/python/pgchangefeed/integration/test_http_realserver.py` Z. 40–47 (`PgChangeFeedHttpClient(...)`, `register_consumer(RegisterConsumerRequest(...))`); `sdks/python/pgchangefeed/tests/test_http_client.py` Z. 79 (`AcknowledgeConsumerRequest(consumer_id=, source_id=, offset=)`), Z. 97 (`get_consumer_position`), Z. 232 (`read_changes(source=, from_=)`) |
| Python · gRPC | `…/integration/test_grpc_realserver.py` Z. 42–44 (`grpc.insecure_channel`, `PgChangeFeedGrpcClient(channel, ClientOptions(...))`, `stream_changes`); `new_image.decode()` Z. 55 |
| Python · SSE | `…/integration/test_sse_realserver.py` Z. 43–45; `httpx.Timeout(10.0, read=None)` ist die `httpx`-API für einen Stream ohne Lese-Frist (das Testbeispiel nutzt `timeout=95.0` für seine 90-s-Frist) |
| Python · NATS | `…/integration/test_nats_realserver.py` Z. 47–50 |
| Python · Fehlerbehandlung | `…/tests/test_http_client.py` Z. 314 (`pytest.raises(PgChangeFeedForbiddenError)`), Klassen aus `src/pgchangefeed/exceptions.py` |
| C# · HTTP-Schnellstart | `sdks/csharp/PgChangeFeed.Client.Integration/HttpRealserverTests.cs` Z. 27–30 (`PgChangeFeedClientOptions`, `PgChangeFeedHttpClient`, `RegisterConsumerAsync`); `…Tests/Http/PgChangeFeedHttpClientConsumerTests.cs` Z. 47 (`AcknowledgeConsumerRequest("c-1", "s-1", 42)`), Z. 68 (`GetConsumerPositionAsync`); `…Tests/Http/PgChangeFeedHttpClientRetentionAndChangesTests.cs` Z. 47 (`ReadChangesAsync(… from:, to:, limit:)`) |
| C# · gRPC | `…Integration/GrpcRealserverTests.cs` Z. 24, 27, 36 (`new PgChangeFeedGrpcClient(options)`, `StreamChangesAsync`, `NewImage?.ToStringUtf8()`) |
| C# · SSE | `…Integration/SseRealserverTests.cs` Z. 23–28; `Timeout.InfiniteTimeSpan` ist die .NET-API für einen `HttpClient` ohne Anfrage-Frist |
| C# · NATS | `…Integration/NatsRealserverTests.cs` Z. 23–28 (`BuildSourceSubject`, `StreamChangesAsync(subject, …)`); `BuildSubject` aus `Nats/PgChangeFeedNatsStreamClient.cs` |
| C# · Fehlerbehandlung | `…Integration/HttpRealserverTests.cs` Z. 51 (`ThrowsAsync<PgChangeFeedUnauthorizedException>`); Klassen aus `Http/PgChangeFeedException.cs` |
| Kotlin · HTTP-Schnellstart | `sdks/kotlin/pgchangefeed-kotlin/src/integrationTest/kotlin/io/github/pt9912/pgchangefeed/integration/HttpRealserverTest.kt` Z. 31–34 (`PgChangeFeedHttpClient(httpClient, options)`, `registerConsumer`); `src/test/kotlin/…/http/PgChangeFeedHttpClientConsumerTest.kt` Z. 47, 66 (`acknowledgeConsumer(AcknowledgeConsumerRequest("c-1", "s-1", 42L))`, `getConsumerPosition`); `…RetentionAndChangesTest.kt` Z. 52 (`readChanges(source =, … from = 1L …)`) |
| Kotlin · gRPC | `…integration/GrpcRealserverTest.kt` Z. 28, 32, 36 (`PgChangeFeedGrpcClient(options).use`, `streamChanges()`, `newImage?.toStringUtf8()`) |
| Kotlin · SSE | `…integration/SseRealserverTest.kt` Z. 23, 27 |
| Kotlin · NATS | `…integration/NatsRealserverTest.kt` Z. 22–27 |
| Kotlin · Fehlerbehandlung | `src/test/kotlin/…/http/PgChangeFeedHttpClientAuthBoundaryTest.kt` Z. 24, 39 (`assertFailsWith<PgChangeFeedUnauthorizedException>`/`…ForbiddenException>`) |
| Kotlin · Installation (Registry, Koordinate) | die bisherige README-Fassung des Pakets (Registry-Block, Gradle-Zugangsdaten), `pom`-Ausgabe: `groupId io.github.pt9912`, `artifactId pgchangefeed-kotlin`; die Abhängigkeits-Hinweise stützt die erzeugte `pom` (acht der neun Abhängigkeiten stehen mit `<scope>runtime</scope>` — `gson`, `grpc-kotlin-stub`, `grpc-protobuf`, `grpc-stub`, `protobuf-java`, `kotlinx-coroutines-core-jvm`, `jnats`, `grpc-netty-shaded` —, `kotlin-stdlib` mit `compile`; gemessen mit `./gradlew generatePomFileForMavenPublication` im Bau-Image, Fixrunde) |

**Belege für die Artefakt-Inhalte (Läufe 2026-09-25):** `make sdk-pack-python` (Exit 0):
das Wheel-`METADATA` trägt `Version: 0.2.1`, die vier `Project-URL`-Zeilen, `Keywords:`,
`Description-Content-Type: text/markdown` und den README-Text als Beschreibung
(0 Kennungs-Treffer, 219 Zeilen); das `tar.gz` trägt `README.md` byte-gleich
(`cmp` ohne Ausgabe); `twine check` im gepinnten `python:3.14-slim`-Image:
`PASSED` für beide Artefakte. `make sdk-pack-csharp` (Exit 0): `unzip -l` des
`PgChangeFeed.Client.0.2.1.nupkg` nennt `README.md` (12049 Bytes), byte-gleich der
Quelle (`cmp`), das `nuspec` trägt `<version>0.2.1</version>`, `<readme>README.md</readme>` und
`<tags>postgresql change-data-capture cdc logical-replication change-feed</tags>`.
`make sdk-pack-kotlin` (Exit 0): `pgchangefeed-kotlin-0.2.1.jar`; die `pom` trägt `name`,
`description`, `url`, `licenses`/`license` (MIT) und `scm`; das `jar` trägt keine README
(`unzip -l … | grep -ci readme` → 0).

**Wächter-Mutationen der ersten Fassung (Zusage · mutierte Eingabe · gesehenes Rot):**
Python-Wächter, je einzeln in die README eingeführt und mit `make sdk-pack-python` gefahren
(Exit 2), danach zurückgenommen: (a) API-Tabelle `read_changes(…, from_, …)` → `from`:
`test_the_api_overview_names_real_methods_with_their_parameters` rot; (b) Beispiel
`client.read_changes(` → `client.read_change(`: rot (`not a client method: ['read_change']`);
(c) Feld-Tabelle `committed_at` → `committed_on`:
`test_the_change_table_lists_exactly_the_fields_of_change` rot; (d) Import im Beispiel
`AcknowledgeConsumerRequest` → `…Requests`: `test_a_readme_example_compiles_and_imports_existing_names[0]` rot;
(e) Fehler-Tabelle `PgChangeFeedNotFoundError` → `…Errors`: rot (jeweils `1 failed, 63 passed`).
Diese fünf banden die Aufruf-Argumente der Beispiele, die Feldliste der Stream-Nachricht und die
Statuscodes nicht (Review-Finding F-1: `api_token=` → `token=`, `offset=` → `position=`,
`from_=` → `start=`, `schema_version` aus der Feldliste, `403` → `404` blieben grün).

**Wächter-Mutationen der Fixrunde** (je einzeln in `sdks/python/README.md` bzw. eine Quelldatei
eingeführt, `make sdk-pack-python` in `/tmp/claude-1000/…/scratchpad/mut-<x>.log`, danach `git checkout`;
gedruckte Zeile jeweils `1 failed, 82 passed`, bei der Konstruktor-Zeile der Stream-Tabelle `3 failed, 80 passed`, Exit 2):

| Zusage · mutierte Eingabe | gesehenes Rot |
|---|---|
| Argument-Namen der Beispiele: `ClientOptions(…, api_token=` → `token=` | `test_every_call_of_a_readme_example_fits_the_real_signature[0]` (`call does not fit the signature of ClientOptions`) |
| Argument-Namen im Request-Konstruktor: `AcknowledgeConsumerRequest(… offset=` → `position=` | dieselbe Prüfung `[0]` (`… AcknowledgeConsumerRequest`) |
| Argument-Namen einer Client-Methode: `client.read_changes("my-source", from_=` → `start=` | dieselbe Prüfung `[0]` (`client.read_changes called with arguments its signature does not accept`) |
| Feldliste der Stream-Nachricht: `schema_version` aus dem Absatz gestrichen | `test_the_stream_paragraph_lists_exactly_the_fields_of_the_stream_message` |
| Statuscode der Fehler-Tabelle: `PgChangeFeedForbiddenError` `403` → `404` | `test_the_error_table_names_every_exception_with_its_status_code` (`404 == 403`) |
| Abdeckung der API: Zeile `remove_consumer` aus der API-Tabelle gelöscht | `test_the_api_overview_covers_every_public_method` (`extra … remove_consumer`) |
| Parameter der API-Tabelle: `get_status(source, schema, table, publication)` ohne `publication` | `test_the_api_overview_names_real_methods_with_their_parameters` |
| Konstruktor der Stream-Tabelle: `PgChangeFeedSseClient(client, options)` → `(http, options)` | `test_the_api_overview_constructors_match_the_real_ones` (`['http', 'options'] == ['client', 'options']`) |
| Beispiel je Stream-Client: `PgChangeFeedNatsStreamClient(` im Beispiel → `PgChangeFeedNatsStreamClient2(` | `test_every_client_has_a_readme_example` |
| Kennungs-Wächter des Pakets: `(SPEC-018)` in `options.py` eingefügt | `test_package_sources_carry_no_internal_identifier[options.py]` |
| Kennungs-Wächter der SDK-Dateien: `(ADR-0109)` in `sdks/kotlin/pgchangefeed-kotlin/settings.gradle.kts` eingefügt, `make sdk-pack-kotlin` | Exit 2 vor dem Docker-Bau: `sdk-public-doc-check: interne Kennung in den SDK-Dateien` (`sdk.mk:22`) |
| Tabellentest des Kennungs-Wächters: je Kennungsart eine Datei (Python/C#/Kotlin/README/Build/Dockerfile) und die Ausnahmen (`obj`, `dist`, `grpc_gen`, Wort ohne Kennung) | `make test-sdk-public-doc-check`: „alle Fälle bestanden“ (elf Fälle; die Treffer-Fälle erwarten Exit 1) |

Nicht gefahren, Grund: eine Mutation auf der Code-Seite (Methode oder Feld im SDK
umbenannt) — die Eingabeseite der Wächter ist die README bzw. der Quelltext-Text; die
Code-Seite färbt denselben Vergleich (Tabelle/Beispiel gegen `inspect`/`dataclasses`) und
zusätzlich die bestehenden SDK-Tests. Für C# und Kotlin wird der Text der Beispiele weiterhin
nicht dauerhaft gewacht (§1); ihre Kennungs-Freiheit halten `make sdk-public-doc-check` und
der Tabellentest.

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaften: „die Version der drei
Packages (0.2.0 → 0.2.1)“, „Kennungen in Paket-Beschreibung und Metadaten-Feldern“,
„der Text der README als Paketbeschreibung“; beide Stände gemessen, Zahlen wie
gedruckt; Parent `6839f738`, Diff-Stand = Arbeitsbaum am 2026-09-25):**

| Träger / Eigenschaft | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| Kennungen in den drei READMEs | `git grep -nE 'SPEC-\|ADR-\|LH-FA\|LH-QA\|ARC-' 6839f738 -- sdks/python/README.md sdks/csharp/README.md sdks/kotlin/pgchangefeed-kotlin/README.md \| wc -l` (Parent); `grep -rnE '<dasselbe Muster>' <dieselben drei Dateien> \| wc -l` (Diff-Stand) | Parent **11** Zeilen (Python 4, C# 3, Kotlin 4); Diff-Stand **0** | gelöst (Umschreibung der drei READMEs) |
| Kennungen in den Nutzer-Feldern der Metadaten (`<Description>`, `<PackageTags>`, `description`, `keywords`, `[project.urls]`-Zeilen, `description.set`/`name.set`/`url.set`) | dasselbe Muster, gefiltert auf diese Feld-Zeilen, je Datei über `git show 6839f738:<Datei>` bzw. Arbeitsbaum | Parent **0**, Diff-Stand **0** in `.csproj`, `pyproject.toml`, `build.gradle.kts` | kein Handlungsbedarf; die Felder waren frei, der Slice hält es bei den neuen Feldern ein |
| Kennungen in den drei Build-Dateien insgesamt (Kommentare) — Stand der ersten Fassung | `git show 6839f738:<Datei> \| grep -cE '<Muster>'` / `grep -cE '<Muster>' <Datei>` | `.csproj` **7 / 7**, `pyproject.toml` **5 / 5**, `build.gradle.kts` **14 / 13** | in der Fixrunde gelöst: siehe die Zeile „Kennungen in den SDK-Dateien“ unten (Parent `97363562` 5/7/… → HEAD 0) |
| Kennungen und Slice-/Welle-Namen in den SDK-Dateien insgesamt (Quellen, Tests, Build-Dateien, Dockerfiles, `.gitignore`; Fixrunde) | `git grep -nE 'SPEC-[0-9]\|ADR-[0-9]\|LH-(FA\|QA)\|ARC-[0-9]\|\bslice-[a-z0-9]\|\bwelle-[0-9a-z]' <Stand> -- <Pfad-Gruppe>` (ohne `grpc_gen`), Zeilen/Dateien je Gruppe, Parent `97363562` → HEAD (`make sdk-public-doc-check` am Arbeitsbaum: Exit 0, gedruckt „keine interne Kennung unter sdks“) | Python `src` 93 Zeilen/8 Dateien → **0**; Python Tests+Integration 35/8 → **0**; `pyproject.toml` 5/1 → **0**; C# Produktion+`.csproj`+`props` 105/17 → **0**; C# Tests+Integration 44/21 → **0**; Kotlin `main`+`build.gradle.kts`+`settings.gradle.kts` 132/19 → **0**; Kotlin Tests+`integrationTest` 42/19 → **0**; Dockerfiles+`.gitignore` 43/5 → **0**; READMEs 0 → 0; `sdks` gesamt **499 Zeilen in 98 Dateien → 0** | gelöst (Umschreibung je Datei, Wächter §3) |
| Laufzeit-Fehlertexte mit Kennung | `git grep -nE 'protocol violation outside' <Stand> -- sdks` | Parent Python **8**, C# **6**, Kotlin **6**; HEAD **0** | gelöst; ein Test oder Integrationstest, der diese Texte prüft, existiert nicht (`grep` der SDK-Test- und Integrationsquellen nach den Textfragmenten: 0 Treffer — die `assertEquals(…, ex.message)`-Zeilen der Tests prüfen Fehlertexte, die die Fakes selbst liefern) |
| Chronik-/Jargon-Wörter in den SDK-Dateien | `git grep -nE 'Draht-Kenntnis\|unstrittig\|Welle-Plan\|Boundary \(\|zuvor\|bisher' <Stand> -- sdks` | Parent **28** Zeilen (Python 7, Kotlin `main` 11 und `test` 1, C# 9); HEAD **0** | gelöst |
| Sources-Jar / XML-Dokumentationsdatei als Träger der Kommentare | `unzip -l sdks/csharp/dist/PgChangeFeed.Client.0.2.1.nupkg`, `unzip -l sdks/kotlin/dist/pgchangefeed-kotlin-0.2.1-sources.jar`, `unzip -p … \| strings [-e l] \| grep -cE 'SPEC-\|ADR-\|LH-(FA\|QA)'` über `.xml`, `.dll` (auch UTF-16-Literale), Jars | nupkg trägt `lib/net10.0/PgChangeFeed.Client.xml` (59422 Bytes), Sources-Jar 37 Einträge; 0 Treffer in `.xml`, `.dll`, den Jar-Dateien und im Wheel (`.py`-Dateien, einschließlich `changestream_pb2_grpc.py`) | Belege der Auslieferung |
| Versionsangaben `0.2.0` | `git grep -c '0\.2\.0' <Stand> -- spec harness sdks docs/user tools .github Makefile` | Parent **29** Zeilen in 11 Dateien (`benutzerhandbuch.md` 4, `harness/README.md` 3, `harness/mk/sdk.mk` 2, `sdks/csharp/Dockerfile` 1, `.csproj` 2, `sdks/kotlin/Dockerfile` 2, Kotlin-README 1, `build.gradle.kts` 3, `pyproject.toml` 1, `lastenheft.md` 1, `pflichtenheft.md` 9); Diff-Stand (Stand der ersten Fassung, `92e2ca6c`) **8** Zeilen in 3 Dateien (`benutzerhandbuch.md` 4, `lastenheft.md` 1, `pflichtenheft.md` 3); Fixrunde-Stand (HEAD nach der Fixrunde) **21** Zeilen in 5 Dateien (die acht wie zuvor, dazu `docs/user/releasing.md` 10 und `harness/README.md` 3) | die 21 verschobenen Zeilen stehen auf `0.2.1` (`git grep -n '0\.2\.1'`: Parent `97363562` 22 Zeilen, HEAD 22 Zeilen — die Fixrunde hebt keine Version); die 8 ursprünglichen sind Records: vier Handbuch-Historie-Zeilen (1.32, 1.38, 1.39, 1.43), eine Lastenheft-Historie-Zeile (Fassung 0.2.0 des Dokuments), drei Pflichtenheft-Historie-Zeilen (2026-09-22/23); die 13 der Fixrunde sind Belege veröffentlichter Tags und Läufe (`sdk-…-v0.2.0`, F-4) in `releasing.md` und den drei Workflow-Zeilen der `harness/README.md` — sie berichten ein vergangenes Ereignis und sind keine Versionsträger der Quellversion |
| Veröffentlichte Tags | `git tag -l 'sdk-*'` | `sdk-csharp-v0.1.0`, `sdk-csharp-v0.2.0`, `sdk-kotlin-v0.2.0`, `sdk-python-v0.1.0`, `sdk-python-v0.2.0` (Stand der Planung); `git ls-remote --tags origin 'sdk-*0.2.1*'` (Closure, 2026-09-25): `sdk-csharp-v0.2.1`, `sdk-kotlin-v0.2.1`, `sdk-python-v0.2.1`, alle drei auf `7df775c8` | `0.2.0` ist an allen drei Stellen veröffentlicht: die Beschreibung erscheint nur mit einer neuen Version — Grund der Hebung auf `0.2.1` (Versionsentscheidung) |
| Publish-Workflows (Tag-Abgleich gegen die Metadaten-Quelle) | `bash tools/harness/sdk-<sprache>-release-tag-info.sh sdk-<sprache>-v0.2.1` (je Sprache) und `make test-sdk-<sprache>-release-tag-info` | je `version=0.2.1`, Exit 0; die drei Testläufe `alle Fälle bestanden` | die Extraktion der Version aus `<Version>`/`^version = "…"` trifft die Kommentar-Umformulierungen nicht |
| Träger, die den README-Inhalt beschreiben | `git grep -n -i 'README' -- spec harness AGENTS.md docs/user docs/plan/adr sdks tools .github Makefile` gefiltert auf SDK/Package/NuGet/PyPI | nur Pfad-Verweise (Handbuch) und die Bau-Kommentare der Dockerfiles/`.csproj`/`pyproject.toml` über den Weg der Datei in das Paket; keine Beschreibung des Inhalts | kein Nachzug nötig |
| README-Aussagen zu `null`/leer der Row Images (F-2), je Sprache gegen den Quelltext | `git grep -nE 'isJsonNull\|Assert.Null\(change.OldImage\|old_image is None' HEAD -- sdks` | Kotlin: KDoc `Changes.kt:24-28` und Test `…RetentionAndChangesTest.kt:72` belegen `JsonElement.JsonNull` (`isJsonNull == true`) statt `null`; C#: `Assert.Null(change.OldImage)` (`JsonElement?`); Python: `assert change.old_image is None`; Kotlin-SSE/-NATS-Tests prüfen `oldImage?.isJsonNull` | gelöst: Kotlin-README nennt `JsonNull`, C#/Python bleiben `null`/`None`; gRPC-Bilder „empty when there is none“ (`ByteString`/`bytes`) in allen drei READMEs |
| README-Aussage „der Server bedient kein TLS“ (F-12) | `git grep -nEi 'tls\.Config\|ListenAndServeTLS\|ServeTLS\|credentials.NewTLS\|tls\.Listen' HEAD -- internal cmd \| wc -l` und `git grep -nEi '\btls\b' HEAD -- docs/user/benutzerhandbuch.md \| wc -l` | **0** und **0** — weder Server-Code noch Handbuch führen TLS | Aussage in den drei READMEs; „Python ≥ 3.14“: `requires-python = ">=3.14"` (`pyproject.toml`) nennt keinen Grund, die README nennt keinen (§6) |
| Träger der Aussage „Kotlin-Zustellweg unbewiesen“ (F-4) | `git grep -nE 'unbewiesen' HEAD -- docs/user harness/README.md \| wc -l`; `gh run list --workflow sdk-kotlin-release.yml`; `gh api /users/pt9912/packages/maven/io.github.pt9912.pgchangefeed-kotlin/versions --jq '.[].name'` | HEAD **2** Treffer: `releasing.md` Historie-Zeile 1.2 (ein Record) und die `hub-description.yml`-Zeile der `harness/README.md` (anderer Gegenstand); Lauf 36117929552 `success`, Paketversion `0.2.0` (2026-09-25) | `releasing.md` (Version 1.7) und die drei SDK-Workflow-Zeilen der `harness/README.md` nachgezogen |
| Deutsche Wortfragmente und Chronik in den englischen READMEs | `grep -niE 'vollinhalt\|welle\|harness\|rohform\|boundary\|derzeit\|bisher\|vorher\|urspr\|pre-1\|noch nicht\|\byet\b\|previous\|formerly\|no longer\|already\|slice\|\bADR\b\|\bSPEC\b\|fire-and-forget\|currently\|until now\|used to\|earlier\|intentionally' <drei READMEs>` | Stand der ersten Fassung (`92e2ca6c`): 9 Treffer; Arbeitsbaum über `d9e048fa` (gemessen 2026-09-25, Closure; Verifikation §7 nennt dieselbe Zahl): **12** Treffer, je vier in den drei READMEs — alle „already“/„no longer“ in Sachaussagen zum Server-Zustand („`retained` … already stored“, „no longer captured“, Felder `already_registered` und `already_enabled`) — keine Chronik, kein deutsches Fragment | belassen |

## 4. Trigger

**Start** (`next` → `in-progress`): kein anderer Slice liegt in `in-progress/` (WIP-Limit 1);
Nutzer-Anweisung liegt vor.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): wenn die Beispiel-Belege
  Änderungen an den SDK-Signaturen verlangten (dann Code-Slice je Package).
- `in-progress` → `open` (blockiert — Carveout?): wenn ein Pack-Lauf am Netz oder Speicher
  des Hosts scheitert und nicht wiederholbar ist.

## 5. Closure-Trigger

DoD vollständig, `make gates` grün, Review-Report liegt vor und ist aufgelöst, Verifikation
bestätigt die DoD, Closure-Notiz mit Steering-Loop-Eintrag geschrieben.

## 6. Risiken und offene Punkte

- Eine README-Aussage stimmt nicht mit dem Verhalten des Servers überein, weil sie aus
  dem Quelltext des SDK statt aus dem Server gelesen wurde — **Ausgang: eingetreten**, im
  Slice aufgelöst: zwei README-Aussagen stimmten nicht mit dem Code überein (Kotlin
  `JsonNull` statt `null`, Review F-2; Transportfehler außerhalb der Fehlerklassen, Review
  F-5), beide behoben (Verifikation §3); die Server-Seite (Positionen, Tokens, `limit`,
  `from`/`to`, `origin`, Subjekt-Namensraum, `already_enabled`) las der Verifier gegen
  Server-Code und Pflichtenheft (`SPEC-018`, `SPEC-022`) ohne Befund (Verifikation §11,
  „README-Beispiele (Anwendersicht)“), der Reviewer gegen `SPEC-022` und das
  Benutzerhandbuch (Review, Negativbefunde); gelesen, nicht gegen einen
  Server-Container gefahren. Register: `BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall`
  (2×).
- Die Beschreibung erscheint erst nach einem Release (Tag-Push, Betreiber-Handlung) auf den
  Paket-Seiten; ob sie dort wie erwartet gerendert wird, ist lokal nicht prüfbar
  ([`AGENTS.md`](../../../../AGENTS.md) §3.10 sinngemäß: externe Oberfläche) — **Ausgang:
  entfallen** für den Inhalt an den Registries: die Publish-Läufe 36129672841 (C#/NuGet),
  36129672859 (Python/PyPI) und 36129672852 (Kotlin/GitHub Packages) endeten je `success`
  (`gh run view <id> --json conclusion,headBranch`, 2026-09-25; Tags `sdk-*-v0.2.1` auf
  `7df775c8`, `git ls-remote --tags origin`). PyPI: `https://pypi.org/pypi/pgchangefeed/0.2.1/json`
  trägt die neue Beschreibung (README-Text als `description`, 0 Kennungs-Treffer; im
  Closure-Lauf gelesen, auch vom Auftraggeber des Closure-Zugs). NuGet: das `.nupkg` aus
  `api.nuget.org/v3-flatcontainer` trägt `README.md` (13150 Bytes), `nuspec` (Version
  `0.2.1`, Beschreibung, Tags) und `lib/net10.0/PgChangeFeed.Client.xml` (59422 Bytes),
  0 Kennungs-Treffer in README und XML (gemessen 11:36 UTC, sechs Minuten nach dem Push;
  der erste Abruf um 11:34 UTC lieferte noch `404`, die Indexierung folgt dem Push).
  **Ungesehen** bleibt die grafische Darstellung der Paket-Seiten (PyPI, NuGet.org,
  GitHub Packages): sie ist eine externe Oberfläche, an der kein Träger im Repo hängt.
- Die Kotlin-Beschreibung liegt im `pom`; ob GitHub Packages sie anzeigt, ist lokal nicht
  prüfbar (die Kommentare in `build.gradle.kts` machen dazu keine Aussage mehr) — **Ausgang:
  entfallen**: kein Träger im Repo behauptet die Anzeige (`build.gradle.kts` sagt, die
  POM-Felder seien die mit den Artefakten veröffentlichten Metadaten), und die
  veröffentlichte `pom` trägt `name`, `description`, `url` und Lizenz MIT (abgerufen von
  `maven.pkg.github.com/pt9912/pg-change-feed/io/github/pt9912/pgchangefeed-kotlin/0.2.1/`
  mit dem `gh`-Token, 2026-09-25); die Anzeige auf der GitHub-Packages-Seite ist ungesehen.
- Das Sources-Jar (`java { withSourcesJar() }`) gehört über `from(components["java"])` zur
  Veröffentlichung nach GitHub Packages; die Publikations-Metadaten verweisen darauf
  (`generateMetadataFileForMavenPublication`), der Upload gegen `maven.pkg.github.com` ist
  lokal nicht prüfbar — **Ausgang: entfallen**: `pgchangefeed-kotlin-0.2.1-sources.jar`,
  `pgchangefeed-kotlin-0.2.1.jar` und `pgchangefeed-kotlin-0.2.1.pom` liefern am selben
  Pfad je HTTP `200` (`curl -I -L` mit dem `gh`-Token, 2026-09-25); der Publish-Schritt
  „Nach GitHub Packages veroeffentlichen“ des Laufs 36129672852 ist `success`
  ([`AGENTS.md`](../../../../AGENTS.md) §3.10).
- Die geänderten Fehlertexte können Integrationstests betreffen. Die SDK-Test- und
  Integrationsquellen tragen keine Prüfung dieser Texte (`grep` nach den Fragmenten: 0 Treffer,
  §3; Verifikation §6 bestätigt es, dazu AST-Gleichheit bis auf die Texte) —
  **Ausgang: entfallen**: kein Test kann an den Texten rot werden. **Benannte Grenze:**
  `make test-sdk-python-integration`, `make test-sdk-csharp-integration` und
  `make test-sdk-kotlin-integration` sind nach der Bereinigung nicht gefahren (Nutzer-Auftrag:
  keine Realserver-Läufe) und laufen in keinem Workflow (`git grep -n 'test-sdk-.*-integration'
  -- .github` druckt 0 Zeilen, gemessen 2026-09-25); ein Slice, der SDK-Code ändert, fährt sie
  nach seiner DoD, dieser ändert kein Verhalten.
- Der Kennungs-Wächter für C# und Kotlin ist ein Repo-Werkzeug (`make sdk-public-doc-check`),
  kein Gate und kein Teil der Test-Bau-Kontexte; ein Rückfall fällt beim `make sdk-pack-*` auf,
  nicht bei `make gates` — **Ausgang: weiter offen** → `slice-sdk-public-doc-check-gate`
  (Datei in `open/`; Entscheidung der Closure: Aufnahme als Gate, netzlos, `real 0m0,013s`;
  die ADR als Gate-Träger ist Architect-Arbeit, siehe §1).
- Die README nennt den Grund von „Python 3.14 or newer“ nicht (`requires-python = ">=3.14"`
  nennt ihn nicht; `grep -rln 'requires-python\|Python 3\.14\|python:3\.14' docs/plan/adr spec` trifft nur die Begründung des Basis-Images in einer ADR, nicht die Untergrenze) — **Ausgang: weiter offen**
  → Register `BEO-PGC/sdk-python-untergrenze-ohne-anwender-begruendung` (1×), keine Anwender-Aktion.

## 7. Closure-Notiz

- **Was hat funktioniert:** Die Rollen-Kette lief in getrennten Kontexten, und der
  Sensor war die Mutation der Eingabeseite. Der Reviewer fuhr Mutationen (M1 bis M8) gegen
  den Python-Wächter und fand den HIGH (F-1: Argumente, Feldliste und Statuscodes der
  README ungebunden) durch Nachfahren, nicht durch Lesen. Der Verifier fuhr sechs eigene
  README-Mutationen (A bis F, je genau ein Test rot), zwei Wächter-Pfade (G: Stub über die
  `.proto`, H: C#-Datei), übersetzte die C#- und Kotlin-Beispiele in Temp-Dateien der
  Test-Bau-Kontexte (Exit 0, je ein Gegenversuch rot) und scannte die fünf ausgelieferten
  Artefakte auf Kennungen (0 Treffer). Die Kennungs-Bereinigung ist gemessen: `git grep`
  über `sdks` ohne `grpc_gen`, Parent `97363562` **499 Zeilen in 98 Dateien**, `HEAD` 0
  (Verifikation §3, **übernommen**; Muster und Gruppen im Suchlauf §3). Die Gate-Zeilen
  des Verifiers stehen in seinem Report §1 (`make gates` Exit 0, `coverage-gate: OK —
  Coverage 83.10%`, `d-check: 1144 Datei(en) geprüft, 0 Befund(e)`, **übernommen**). Der
  Slice ändert an Go nur Kommentare in `gen/cdc/stream/v1/*.pb.go` (Verifikation §6:
  `git diff -U0` ohne Kommentarzeilen ist für `proto/` und `gen/` leer, `make
  generated-sync` Exit 0); Zahlen in `harness/sensors/*.md` sind davon nicht betroffen
  (`git grep -n -iE 'sdks|0\.2\.[01]' -- harness/sensors` druckt 0 Zeilen, gemessen 2026-09-25;
  der Nenner der Coverage-Stufe in `harness/sensors/coverage-gate.md` hängt an den
  Statements, und die Kommentarzeilen ändern keines, **abgeleitet**), kein Nachzug.
- **Was ging anders als geplant:** (1) *Umfangserweiterung auf Nutzeranweisung mitten im
  Slice* (Feststellung, kein Fehler): der Plan nahm README und Metadaten-Felder in den
  Umfang; auf den Review F-3 hin und die Anweisung „öffentliche API-Kommentare auch
  bereinigen“ kamen die Docstrings, XML-Kommentare, KDoc, Testkommentare, Build-Dateien,
  20 Laufzeit-Fehlertexte und die `.proto`-Kommentare hinzu (106 Dateien in `sdks proto
  gen`, Verifikation §6). Begründung im Plan §1: ein veröffentlichter Paketstand ist
  unveränderlich, die Bereinigung gehört in `0.2.1`. (2) *Die Fixrunde lief ohne zweiten
  Reviewer-Durchgang* — **benannte Grenze** (V-2): der Verifier prüfte ihre Semantik durch
  AST-Vergleich (Python), Nicht-Kommentar-Diff (C#, Kotlin), Artefakt-Scan und Mutationen
  (Verifikation §6, §4, §5); eine Reviewer-Lesung der neuen Kommentar-Texte fehlt. Dieselbe
  Grenze steht in den Closure-Notizen von `slice-backfill-run-store`, `slice-backfill-e2e`,
  `slice-backfill-bench-richtgroesse`, `slice-backfill-slot-leerlauf-bestaetigung` und
  `slice-backfill-sdk-origin` (`git grep -n -i -E 'Fixrunde lief ohne' -- docs/plan/planning/done`,
  gemessen 2026-09-25: fünf Slice-Pläne); die Vorgänger führen keinen Register-Eintrag, dieser
  Slice hält es gleich (siehe Register unten). (3) *Der Wächter blieb kein Gate*: die
  Begründung im Plan §1 („ein weiteres Gate ändert die Gate-Liste“) ist ein Aufwandsargument
  (Verifikation V-4); die Closure-Entscheidung steht unten. (4) *Neue Artefakt-Inhalte ohne
  ADR*: die XML-Dokumentationsdatei im `.nupkg` und das Sources-Jar (Review F-10,
  Verifikation V-6), Entscheidung unten.
- **Entscheidung V-4 (Wächter als Gate):** Aufnahme in `make gates`, ausgeführt von einem
  Folge-Slice, nicht in dieser Closure. Grund für die Aufnahme: `make sdk-public-doc-check`
  ist netzlos und ein `grep` (`time make sdk-public-doc-check` → `real 0m0,013s`,
  `time make test-sdk-public-doc-check` → `real 0m0,195s`, gemessen 2026-09-25); die
  Alternative „Vorstufe der Pack-Ziele bleibt“ trägt nur für die Auslieferung (die drei
  Release-Workflows fahren `make sdk-pack-*` vor jedem Publish, gemessen in der
  Verifikation §6/§8), nicht für die Früherkennung je Push. Grund gegen die Ausführung hier:
  jedes bestehende Gate trägt eine ADR (`ADR-0041`, `-0045`, `-0054`, `-0084`), und die
  Gate-Liste (`harness/README.md` §Sensors), `harness/mk/sdk.mk` und die Sensor-Doku gehören
  dem Implementer-Zug nach dieser ADR. Adresse: `slice-sdk-public-doc-check-gate` (Datei in
  `open/`). Bis dahin tragen `harness/README.md` (Zeile `make sdk-public-doc-check`) und
  `harness/mk/sdk.mk` die Aufwands-Begründung „ein weiteres Gate ändert die Gate-Liste“;
  der Folge-Slice ersetzt sie (Suchlauf im Slice-Plan, DoD 2). Keine Schwelle gesenkt
  ([`AGENTS.md`](../../../../AGENTS.md) §3.6).
- **Entscheidung V-6 (ADR für die neuen Artefakt-Inhalte):** keine neue ADR. Begründung:
  `ADR-0106` sagt, der Bau „erzeugt eine `.nupkg` als Artefakt“ — das Paket bleibt eines, die
  XML-Datei liegt darin; `ADR-0109` nennt den Artefaktumfang nicht (`grep -n -iE
  'sources|javadoc|\.jar' docs/plan/adr/0109-kotlin-github-packages-drittes-sdk-package.md` druckt 0
  Zeilen, `grep -n -iE 'jar|Artefakt'` zwei Vergleichszeilen ohne Umfangsaussage, gemessen
  2026-09-25); `ADR-0107` (Wheel und `sdist`) bleibt unberührt. Der Zusatz ist additiv:
  kein Anwender-Code bricht, die Signaturen sind unverändert. Belegt am veröffentlichten
  Stand: das `.nupkg` von NuGet.org trägt `lib/net10.0/PgChangeFeed.Client.xml`, GitHub
  Packages liefert `pgchangefeed-kotlin-0.2.1-sources.jar` (Messungen in §6). Die
  Eigenschaft steht in den `make sdk-pack-*`-Zeilen der `harness/README.md` (§Sensors) und
  im Plan §3. Trigger für eine ADR (dann als neue ADR, die `ADR-0106`/`ADR-0109` schärft,
  nicht als Änderung der `Accepted`): eine Festlegung zum Artefaktumfang, die über „Kommentare
  reisen mit dem Paket“ hinausgeht (Javadoc-Jar, Symbol-Package, Umfang der `sdist`).
- **Post-Push-Belege** ([`AGENTS.md`](../../../../AGENTS.md) §3.10 sinngemäß; kein Workflow im
  Diff, die drei Publish-Workflows liefen erstmals mit `0.2.1`): Tags `sdk-csharp-v0.2.1`,
  `sdk-python-v0.2.1`, `sdk-kotlin-v0.2.1` auf `7df775c8` (`git ls-remote --tags origin`);
  Läufe 36129672841 (C#/NuGet), 36129672859 (Python/PyPI), 36129672852 (Kotlin/GitHub Packages)
  je `completed`/`success` (`gh run view <id> --json conclusion,headBranch,headSha`, gemessen
  2026-09-25); CI von `7df775c8` (`gh run list --commit 7df775c888a2ec987a6c8b7ae675f94c2a21cdc6`):
  `ci` 36127605194, `examples` 36127605177, `e2e` 36127605205, je `success`. Die Inhalte an den
  Registries stehen in §6 (Risiko zwei bis vier).
- **Benannte Grenzen (Verifikation V-3, V-5, V-7, V-8), je mit Adresse oder Begründung:**
  V-3: die 20 Fehlertexte sind Laufzeitverhalten ohne Test; die SDK-Realserver-Tiers sind
  nach der Bereinigung nicht gefahren (§6, fünfter Punkt); Adresse: jeder Slice, der SDK-Code
  ändert, fährt sie nach seiner DoD; dieser Slice ändert keine Anweisung des Codes.
  V-5: 16 von 60 öffentlichen Python-Definitionen tragen keinen Docstring (13 `from_json`,
  drei verschachtelte Funktionen des NATS-Clients; Klassen und Client-Methoden tragen alle
  einen), die Kotlin-KDoc-Abdeckung ist nicht gemessen — keine Aktion: der Einstieg des
  Anwenders (Klassen, Client-Methoden, README-API-Tabelle) ist vollständig, die
  `from_json`-Methoden sind Dekodier-Fabriken der Modelle; die Kotlin-Messung ist ein
  Befehl, kein Träger, den der Zug ausführt, der Kotlin-Modelle ändert. V-7: der Satz „The
  fields match the change of the HTTP and SSE surfaces (without the commit position, commit
  time and origin)“ in `proto/cdc/stream/v1/changestream.proto` ist nicht falsch (auch die
  SSE-Nachricht trägt die zehn Felder ohne diese drei), aber die Klammer liest sich wie eine
  HTTP-Aussage; er steht in den Stubs von `0.2.1` aller Sprachen. Adresse: die nächste
  Änderung an dieser `.proto` (Kommentar im selben Zug, ein Stub-Neubau via `make
  proto-generate`, der Wächter `make generated-sync` hält Go). V-8: keine Sicht auf die
  gerenderten Paket-Seiten — §6 zweiter und dritter Punkt (Inhalt gemessen, Darstellung
  ungesehen).
- **Steering-Loop-Eintrag (Lerneintrag):** *Geschärfte Regel (Kandidat, nicht entschieden,
  1×):* die Anwender-Sichtbarkeit von Kennungen. Die Träger-Regel „Kennungen in Prosa
  verlinken“ (d-check `ids`) ist für Leser im Repo gebaut; sie galt auch in Texten, die das
  Repo verlassen (README als PyPI-/NuGet-Beschreibung, Docstrings im Wheel, XML-Doku im
  `.nupkg`, KDoc im Sources-Jar, Fehlertexte, `.proto`-Kommentare in den Stubs) und erzeugte
  dort Kennungen, die kein Anwender auflösen kann — der Nutzer nannte sie Rauschen. Der
  Regel-Kandidat lautet: *Texte, die das Repo verlassen, tragen keine internen Kennungen*.
  Planner-Entscheidung: **keine** `AGENTS.md`-Regel aus einem Auftreten (Register-Regel: erst
  ab 3×); der Eintrag `BEO-PGC/intern-kennungen-in-ausgelieferten-texten` (1×) trägt den
  Kandidaten, der Slice trägt Wächter je Träger. *Neuer Sensor:* zwei, beide kein Gate:
  `sdks/python/pgchangefeed/tests/test_public_text.py` (im Bau, für die zur Bauzeit erzeugten
  Stubs) und `make sdk-public-doc-check` mit Tabellentest `make test-sdk-public-doc-check`
  (`grep` über `sdks/`, Vorstufe der drei `make sdk-pack-*`-Ziele); die Aufnahme des zweiten
  in `make gates` übernimmt `slice-sdk-public-doc-check-gate`. *Benannte Spec-Lücke:* kein
  Dokument des Repos benennt die Anwender-Sicht auf ausgelieferte Texte als eigene
  Leser-Klasse (`ids` verlangt Links, verbietet keine Kennung; §3.7 deckt Kommentar-Klassen,
  keine ausgelieferte Aussage); die Lücke ist der Register-Eintrag, kein ADR-Auftrag.
- **Beobachtungs-Register (`../observations/`):** je Anfall eine Datei
  `evidence/slice-sdk-readme-nutzerdoku.md`, Zähler = Zahl der Dateien.
  *Bestehende Klassen mit weiterer Datei:* `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`
  (F-1, HIGH) **12×**, verkörpert (kein Deckel, weil Schwere ≥ MEDIUM);
  `BEO-PGC/nachzug-laesst-ueberholten-text-stehen` (F-4) **9×**, verkörpert;
  `BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall` (F-2) **2×**, offen;
  `BEO-PGC/zwei-quellen-drift-handbuch-gegen-pflichtenheft` (F-11, Ausprägung Quelltext-Kommentar
  gegen README) **3×**, offen, ohne Ausgang: **Lese-Schritt der nächsten Welle-Closure
  (`welle-transformationen`)**, der `state.md` trägt den Vermerk. *Neue Klassen:*
  `BEO-PGC/intern-kennungen-in-ausgelieferten-texten` (F-3) **1×**, offen;
  `BEO-PGC/sdk-python-untergrenze-ohne-anwender-begruendung` (F-12, Plan §6) **1×**, offen.
  *Gedeckelt, ohne neue Datei (Deckel für verkörperte Einträge ab 10×, `../observations/README.md`):*
  `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (21×, verkörpert) — V-1 (LOW, der
  Verifier fand vor dem Merge die veraltete Trefferzahl 9 statt 12 im Suchlauf-Feld
  „Deutsche Wortfragmente und Chronik“, in dieser Closure auf 12 gesetzt, Ursprung im Feld) und
  F-9 (LOW, Zählwort „alle“ Abhängigkeiten `runtime`, Reviewer vor dem Merge; die `pom`
  trägt acht `runtime` und `kotlin-stdlib` `compile`), Träger-Typ bekannt (Zahl im
  Slice-Plan). *Ohne Eintrag:* F-5 (Aussage reicht weiter als der Code — Transportfehler
  außerhalb der Fehlerklassen, behoben, je Sprache benannt), F-7 (Aussage ohne Anleitung:
  `enable_table`-Felder), F-8 (Wächter-Zählkonstanten, gehört zur Klasse von F-1 und ist mit ihr
  behoben), F-10 (Anwender-Doku-Fläche, umgesetzt als XML-Datei und Sources-Jar) sind
  einmalige, im Slice behobene Ausprägungen ohne Klasse mit Wiederholung im Register; F-6
  (Kommentar behauptet Nichtgeprüftes, LOW) fällt unter den verkörperten Reviewer-Punkt
  „Kommentar trägt keine der Kommentar-Klassen“ (Klausel Zusage) und weitet
  `BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad` (dort: ein Fehlerpfad, hier die
  Anzeige einer externen Oberfläche) nicht zu einem vierten Beleg. V-2 trägt keinen Eintrag
  (siehe „Was ging anders“ 2, Vorgänger-Praxis). *Kein Anfall:*
  `BEO-PGC/deutsches-fachwort-im-englischen-sdk-readme` und
  `BEO-PGC/formvorbild-kopie-traegt-deutsches-wortfragment-weiter` (Suchlauf §3: 12 Treffer,
  alle „already“/„no longer“ in Sachaussagen, kein deutsches Fragment).
  *Einträge bei 3× oder darüber ohne Ausgang mit dieser Datei:* `zwei-quellen-drift-handbuch-gegen-pflichtenheft`
  (3×) — Adresse: Lese-Schritt der Closure von
  [welle-transformationen](../welle-transformationen.md).
- **Folge-Slices:** `slice-sdk-public-doc-check-gate` (Datei in `open/`, Entscheidung V-4).
  **Übergaben an offene Pläne** ([`AGENTS.md`](../../../../AGENTS.md) §3.13): Suchlauf
  `git grep -n -E '0\.2\.[01]|sdk-pack|SPEC-02[678]|Package-Version|Paketbeschreibung' --
  docs/plan/planning/open docs/plan/planning/welle-transformationen.md` (gemessen
  2026-09-25 am Stand `d9e048fa`, vor den Änderungen dieser Closure): keine Zeile trägt
  `0.2.0` oder `0.2.1`; Treffer in `slice-transformationen-betriebsdoku` (`SPEC-026`/`-027`/`-028`
  als berührte Stellen, Z. 21 bis 23) und `welle-transformationen.md` (Z. 273 Kennungsvergabe,
  ohne Bezug; Z. 395 „Package-Versionen der SDKs bleiben unberührt“ — gilt weiter). Ein
  zweiter Suchlauf `git grep -n -E 'sdks/|README' -- docs/plan/planning/open` trifft
  `slice-transformationen-betriebsdoku` (u. a. DoD Z. 72 bis 75, §3 Z. 136 und 137) und die
  `harness/README.md`-Verweise anderer Pläne; `slice-transformationen-e2e-wirkung` trägt den Satz „die SDK-Tiers
  laufen gegen einen Server ohne Regeln“ (`git grep -n 'test-sdk-'`), eine Aussage über
  Testläufe, die der Slice nicht bewegt. **Nachgezogen:** `slice-transformationen-betriebsdoku`
  §3 Zeile „SDK-READMEs“ (Z. 137). Fund am Stand `d9e048fa`: die Zeile nannte zwei von drei
  READMEs und sagte, eine Änderung dort hebe die Package-Version nicht (Metadaten-Text); die
  README ist die Paketbeschreibung auf PyPI und NuGet und erscheint dort mit einer neuen
  Version, und `make sdk-public-doc-check` färbt jede Kennung in einer Datei unter `sdks/`.
  Die Zeile nennt alle drei READMEs, den Abschnitt zum Change-Objekt, die Eigenschaft
  „Paketbeschreibung“ und den Wächter. Nicht betroffen:
  `slice-transformationen-e2e-wirkung` (SDK-Tiers fahren gegen einen Server ohne Regeln,
  keine bewegte Eigenschaft), `slice-capture-*`, `slice-harness-suchlauf-nachmessen`,
  `slice-backfill-speicher-untersuchung`.
- **Risiken aus §6:** je ein Ausgang, mit Beleg in §6. *Eingetreten:* README-Aussage gegen
  den Code (F-2, F-5; im Slice behoben). *Entfallen:* Darstellung und Inhalt auf den
  Paket-Seiten (Registry-Inhalte gemessen, Darstellung ungesehen) · Anzeige der Kotlin-`pom` ·
  Upload des Sources-Jars · Fehlertexte und Integrationstests (Grenze: Tiers nicht gefahren).
  *Weiter offen:* Kennungs-Wächter ohne Gate → `slice-sdk-public-doc-check-gate` · Python-Untergrenze
  → Register `BEO-PGC/sdk-python-untergrenze-ohne-anwender-begruendung`.
- **Drei Paarungen:** dieser Slice ist wellenlos; die Paarungen (Anker · Folge-Slice · Register)
  prüft die Closure von [welle-transformationen](../welle-transformationen.md) (offen; die
  Ereignis-Adresse steht dort §5 „Ereignis-Adresse“, die Roadmap führt sie unter *Offene
  Wellen*). Anker: die Wächter (`test_public_text.py`, `make sdk-public-doc-check`); Folge-Slice:
  `slice-sdk-public-doc-check-gate`; Register: die sechs Einträge oben.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist `sdks/*` (die SDK-Sprachpakete) und
`spec/`/`harness/` als Träger von Versionsangaben; die Modus-Deklaration führt nur die
Default-Sub-Area `*` (`PGC`, Greenfield) — keine Ausdifferenzierung nötig, keine Schwelle
verfehlt.

**Vorgelagert — offene Beobachtungen sichten:** Register `BEO-PGC` durchgegangen; Treffer je
berührter Sub-Area: `deutsches-fachwort-im-englischen-sdk-readme` (Zähler 3×, verkörpert:
Sprachreinheit der englischen READMEs — die drei neuen READMEs sind je Form-Teil gesichtet,
Suchlauf auf deutsche Wortfragmente im Bericht), `formvorbild-kopie-traegt-deutsches-wortfragment-weiter`
(3×, verkörpert), `zahl-in-traeger-driftet-gegen-die-messung` (21×, verkörpert: die Trefferzahlen
im Suchlauf stehen wie gedruckt), `nachzug-laesst-ueberholten-text-stehen` (8×, verkörpert: die
Versionsträger werden am Ende gegen die Beschreibung gesucht), `arbeit-ueberholt-stehenden-traeger`
(31×, verkörpert: Suchlauf §3 als committetes Feld), `beleg-befehl-traegt-seinen-satz-nicht`
(13×, verkörpert: jeder Beleg nennt den Lauf und die gedruckte Zeile),
`drei-sprachen-kopie-divergiert-am-randfall` (1×, offen: die Aussagen zu `origin` und Fehlerklassen
je Sprache werden gegen den Quelltext der jeweiligen Sprache gelesen, nicht aus einer Schwester-README
kopiert).

**Modus:** alle berührten Sub-Areas GF.
