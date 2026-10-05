# Slice sdk-tls-optionen: die drei SDKs verbinden über TLS mit eigenem Vertrauensanker

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — es gibt keine Closure-Bedingung jenseits der DoD
dieses Slice, siehe Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht (Modul 6).

**Abhängigkeit:** harte Voraussetzung: Lastenheft 0.16.0 mit `LH-FA-SST-013` und
`SPEC-037` im Pflichtenheft sind committet (die Folge-Anforderung, auf die
`slice-examples-grpc-tls` wartete). Server-seitig liegt TLS vor
([`slice-tls-http-grpc-server`](../done/slice-tls-http-grpc-server.md)); die
Beispiele unter `examples/` hängen nicht an den SDKs und sind nicht Teil dieses
Slice.

**Bezug:** [`LH-FA-SST-013`](../../../../spec/lastenheft.md) (Scope),
[`LH-FA-SST-009`](../../../../spec/lastenheft.md) (die Packages),
[`LH-FA-SST-011`](../../../../spec/lastenheft.md) (Serverseite),
[`ADR-0150`](../../adr/0150-tls-und-mehrfach-token.md),
[`ADR-0110`](../../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md) (Mechanik der
Realserver-Phasen).

**Berührte Spec-Stellen:** `SPEC-037` · `SPEC-026` · `SPEC-027` · `SPEC-028`
(Package-Versionen, Minor-Sprung).

**Verantwortlich:** — bis zur Priorisierung.

**Autor:** pt9912 (Architect-Zuschnitt). **Datum:** 2026-10-05.

---

## 1. Ziel und Abgrenzung

**Ziel:** Jedes der drei SDK-Packages (C#, Python, Kotlin) nimmt eine Option für
einen eigenen Vertrauensanker (PEM-Datei) an und bedient damit HTTP, SSE,
gRPC-Stream und gRPC-Verwaltung gegen einen Server mit TLS; ein Realserver-Lauf
je SDK belegt es gegen ein Zertifikat, das nicht im Systemspeicher steht.

**Entscheidungen des Zuschnitts** (Quelle: Pflichtenheft `SPEC-037`):

- Eine Option, ein Inhalt: Pfad einer PEM-Datei; mit Anker vertraut die
  Verbindung genau diesen Zertifikaten, ohne Anker gilt der Systemspeicher.
- Die Prüfung ist nie abschaltbar; kein Überschreiben des Servernamens.
- NATS-Fläche bleibt unberührt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Client-Zertifikate (mTLS) — nicht gefordert (`LH-FA-SST-013` Out-of-Scope);
  der Server verlangt keines (`SPEC-034`).
- Eine „unsichere“ Option (Prüfung aus) und ein Servernamen-Override — bewusst
  nicht angeboten; der Override wäre in der JDK-HTTP-Schicht nur mit einem
  eigenen Trust-Manager machbar und öffnet den Weg zur abgeschalteten Prüfung.
  Wer den Namen nicht treffen kann, stellt das Server-Zertifikat auf den Namen
  der Adresse aus. Ein Folge-Wunsch ist eine eigene Anforderung.
- TLS der NATS-Verbindung — Sache der NATS-Konfiguration (wie `LH-FA-SST-011`).
- Die Beispiele und die Server-seitige Konfiguration — anderer Gegenstand
  (`slice-examples-grpc-tls`, `slice-tls-http-grpc-server`).
- Veröffentlichung der neuen Versionen (Tags, Registries) — braucht eine neue
  Freigabe des Auftraggebers; der Slice endet mit gebauten Artefakten
  (`make sdk-pack-*`).

## 2. Definition of Done

- [x] **C#:** `PgChangeFeedClientOptions` trägt die Anker-Option; der HTTP-, SSE-
      und gRPC-Weg (Stream und Verwaltung) nutzt sie; Unit-Tests: Happy
      (Verbindung über TLS mit gültigem Anker), Boundary (kein Anker → Systemspeicher,
      `http://` unverändert), Negative (fremdes Zertifikat, falscher Name,
      nicht lesbarer Pfad, Datei ohne PEM → sichtbarer Fehler, kein Rückfall).
      `LH-FA-SST-013` erfüllt, Test referenziert.
      Belege: §3 „Belege des Implementers“ (Tests `Tls/TlsClientTests.cs`,
      `Tls/TlsOptionsTests.cs`, Mutationen mit gesehener Farbe).
- [x] **Kotlin** und **Python:** wie C#, je Sprache dieselben Testfälle und
      derselbe Inhalt der Option (`SPEC-037`). Belege: §3 „Belege des Implementers“
      (`tls/TlsClientTest.kt`, `tests/test_tls.py`).
- [x] **Realserver-Beleg:** je SDK eine TLS-Phase im Runner
      (`make test-sdk-csharp-integration`, `…-kotlin-…`, `…-python-…`): der Feed-Container
      läuft mit Zertifikat aus `tools/harness/certgen` (Muster der TLS-Phase von
      `make test-integration`, Override-Datei im Temp-Verzeichnis, `compose.yaml`
      bleibt unverändert); je Fläche empfängt/ruft der Client mit Anker; ohne
      Anker scheitert dieselbe Verbindung sichtbar. Der Abdeckungs-Träger
      `docs/user/sdk-e2e-abdeckung.md` trägt die neuen Phasen (Runner schreibt ihn).
      Belege: §3 „Belege des Implementers“ (gedruckte Zeilen je Runner, Mutationsläufe).
- [x] `make gates` grün (inkl. `make sdk-public-doc-check`: keine interne Kennung in
      Quellen, Tests, README der SDKs). Beleg: §3 „Belege des Implementers“,
      Punkt „Gates“.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor, kein Self-Review.
- [x] Doku-Update: README der drei Packages (Englisch, Abschnitt zu TLS; die
      Ist-Aussagen „plaintext gRPC“ in Docstrings und README werden nachgezogen),
      Benutzerhandbuch nur, falls es die SDK-Nutzung beschreibt (es tut es an einer
      Stelle: Abschnitt „Schnittstellen mit TLS verschlüsseln“, Version 1.100, mit
      Historienzeile). Suchlauf: §3.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register fortgeschrieben oder „keine Beobachtung“ in §7.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

(Liefer-Punkte: drei — je SDK Option, Tests, Realserver-Phase.)

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `sdks/csharp/PgChangeFeed.Client/PgChangeFeedClientOptions.cs` und `Http/`, `Sse/`, `Grpc/` | update | Option und Verdrahtung in `HttpClientHandler`/`SocketsHttpHandler` (benutzerdefinierte Vertrauenskette) und `GrpcChannel` (`HttpHandler`) — [`LH-FA-SST-013`](../../../../spec/lastenheft.md) |
| `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/io/github/pt9912/pgchangefeed/PgChangeFeedClientOptions.kt` und `http/`, `sse/`, `grpc/` | update | Option; `SSLContext`/Trust-Manager für `HttpTransport`/`SseTransport`, `TlsChannelCredentials` für gRPC statt `usePlaintext()` bei TLS-Adresse |
| `sdks/python/pgchangefeed/src/pgchangefeed/` (`http_client.py`, SSE-Client, `grpc_client.py`, `administration_client.py`) | update | Option; `ssl`-Kontext für HTTP/SSE, `grpc.ssl_channel_credentials(root_certificates=…)` für gRPC |
| Tests je SDK (Unit, neue Testdateien) | neu | Happy/Boundary/Negative nach `LH-FA-SST-013`; Zertifikate zur Laufzeit erzeugt oder als kleine Test-Fixtures außerhalb der Veröffentlichung |
| `tools/harness/run-sdk-{csharp,kotlin,python}-integration-tests.sh` und `tools/harness/lib-*` | update | TLS-Phase je Runner, Wiederverwendung von `tools/harness/certgen`; **kein Teil des aktuellen parallelen Zuges an `tools/harness`** — erst nach dessen Abschluss anfassen |
| SDK-READMEs, Docstrings | update | Ist-Zustand TLS, keine interne Kennung (`make sdk-public-doc-check`) |
| Package-Versionen (`.csproj`, `build.gradle.kts`, `pyproject.toml`) | **nicht geändert** | Auftrag des Auftraggebers an den Implementer: die Versionen bleiben (`0.6.1`), der Minor-Sprung gehört in den Release-Zug mit neuer Freigabe; `make sdk-pack-*` baut deshalb Artefakte der Version `0.6.1` |

**Plan-Nachzug (Dateien und Entscheidungen des Laufs, über die Zeilen oben hinaus):**

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `sdks/csharp/PgChangeFeed.Client/TlsTransport.cs` | neu | Handler mit Prüfung gegen genau die Anker (`CustomRootTrust`, `NoCheck`), Servername und `RemoteCertificateNotAvailable` aus der Plattform-Prüfung übernommen; ein präsentiertes Zertifikat, das byte-gleich einem Anker ist, gilt nach der Gültigkeitsprüfung (Server-Zertifikat eines Ausstellers als alleiniger Anker); `https` Pflicht bei Anker |
| `PgChangeFeedClientOptions.cs`, `Http/PgChangeFeedHttpClient.cs`, `Sse/PgChangeFeedSseClient.cs`, `Grpc/*Client.cs` (C#) | update | `TrustAnchorFile` (Konstruktor mit dritter Angabe), PEM-Prüfung beim Erzeugen; HTTP- und SSE-Client haben keinen Konstruktor, der die Verbindung selbst baut — neu `PgChangeFeedHttpClient(options)` und `PgChangeFeedSseClient(options)` (`IDisposable`, besitzen ihren `HttpClient`); die gRPC-Clients bauen ihren Kanal über denselben Handler |
| `Grpc/PgChangeFeedAdministrationClient.cs` `MapException` | update | ein TLS-Fehler kommt als `RpcException` mit Status `Internal` und dem `HttpRequestException` als `DebugException` an; er wäre als `PgChangeFeedGrpcInternalException` („Fehler im Server“) gemeldet worden — jetzt `PgChangeFeedGrpcUnexpectedStatusException` ohne Meldungscode (`SPEC-037`: Verbindungsfehler, Server nicht erreicht) |
| `sdks/csharp/.../PgChangeFeed.Client.Tests/Tls/*` und `PgChangeFeed.Client.Tests.csproj` | neu / update | Kestrel-Server auf Loopback (`FrameworkReference Microsoft.AspNetCore.App`, kein neues Paket), Zertifikate zur Laufzeit (`CertificateRequest`), vier Flächen je Fall |
| `sdks/csharp/PgChangeFeed.Client.Integration/*Tls*.cs`, `TlsScenarios.cs`, `PhaseEnvironment.cs` | neu / update | vier TLS-Phasen (HTTP, SSE, gRPC-Stream, gRPC-Verwaltung) mit den drei Verweigerungen |
| `sdks/kotlin/.../TlsSupport.kt`, `PgChangeFeedClientOptions.kt`, `http/`, `sse/`, `grpc/` | neu / update | `trustAnchorFile: Path?` (Zweit-Konstruktor erhält die alte Signatur), `SSLContext`/`TrustManager`; `https`-Adresse → TLS-Kanal (`TlsChannelCredentials`), jedes andere Schema bleibt Klartext; neu `PgChangeFeedHttpClient(options)`, `PgChangeFeedSseClient(options)`; der Trust-Manager der JDK prüft den Gültigkeitszeitraum eines **als Anker eingehängten Server-Zertifikats nicht** (am Test `an expired certificate fails the connection` gesehen) — `ValidityCheckingTrustManager` prüft ihn zusätzlich |
| `sdks/kotlin/.../src/test/.../tls/*`, `src/integrationTest/.../*Tls*.kt`, `TlsScenarios.kt` | neu | Server aus dem JDK (`HttpsServer`) und `grpc-netty-shaded`, Zertifikate zur Laufzeit mit `keytool` im Temp-Verzeichnis |
| `sdks/python/.../options.py`, `tls.py`, `__init__.py` | update / neu | `ClientOptions.trust_anchor_file`; die Python-Clients bekommen ihre `httpx.Client`/`grpc.Channel` vom Aufrufer, deshalb neu `create_http_client(options)` und `create_grpc_channel(options)` (`https://host:port` → TLS, `host:port` mit Anker → TLS) statt eines Parameters an vier Konstruktoren |
| `sdks/python/pgchangefeed/tests/test_tls.py`, `integration/*tls*.py`, `integration/tls_scenarios.py`, `pyproject.toml` (Test-Extra `cryptography`) | neu / update | Server mit `ThreadingHTTPServer`/`grpc.server`, Zertifikate zur Laufzeit mit `cryptography` (nur Test-Extra, keine Laufzeit-Abhängigkeit) |
| `tools/harness/lib-sdk-tls-fixture.sh` | neu | gemeinsame TLS-Vorbereitung der drei Runner (certgen `good`/`other` mit den Namen `pg-change-feed`, `localhost`; der Container-Name `cdc-test-feed` steht bewusst nicht im Zertifikat und ist die Namensabweichung gegen denselben Server) |
| `tools/harness/run-sdk-{csharp,kotlin,python}-integration-tests.sh`, `harness/mk/sdk.mk` | update | vier TLS-Phasen vor dem Abdeckungs-Träger, `RUN_PHASE_DOCKER_ARGS` für das Zertifikatsverzeichnis, `TOOLCHAIN_IMAGE`/`GO_MODCACHE_VOLUME` vom Makefile (Muster `fmt-check`) |
| `docs/user/sdk-e2e-abdeckung.md` | Runner-Erzeugnis | je Sprache eine Zeile (die Runner haben sie geschrieben) |
| `docs/user/benutzerhandbuch.md`, `harness/README.md` | update | Handbuch 1.100 (Abschnitt „Schnittstellen mit TLS verschlüsseln“, der Satz „bieten keine TLS-Einstellung“ entfällt); Sensors-Zeilen der drei Integrationsziele |

```suchlauf
71c6a886 9 -i -E 'plaintext|no TLS' -- sdks
diff 38 -i -E 'plaintext|no TLS' -- sdks
71c6a886 6 -i -E 'does not serve TLS|speaks plaintext|plaintext gRPC' -- sdks
diff 0 -i -E 'does not serve TLS|speaks plaintext|plaintext gRPC' -- sdks
71c6a886 2 -i -E 'keine TLS-Einstellung' -- docs/user
diff 1 -i -E 'keine TLS-Einstellung' -- docs/user
```

Stand der Zeilen: **gemessen** mit `git grep -n` am Parent `71c6a886` und am Arbeitsbaum
(`diff`); der Implementer wiederholt sie mit `make suchlauf-nachmessen PLAN=<Plan-Datei>`.
Die 38 `diff`-Trefferzeilen zu „plaintext“ sind gelesen (die 38. ist der Upgrading-Eintrag
des Kotlin-READMEs der Fixrunde): 15 liegen in Dateien, die
schon am Parent standen (READMEs, Docstrings der Konstruktoren), 23 in den neuen
Dateien (`TlsTransport.cs`, `TlsSupport.kt`, `tls.py`, Tests und Test-Server);
jede nennt Klartext als Zweig neben TLS (`https`/`http`-Gegenüberstellung,
„kein Rückfall auf Klartext“, Test des Klartext-Zweigs); keine behauptet mehr,
der Server spreche nur Klartext. Die eine verbleibende Zeile zu „keine TLS-Einstellung“
steht in der Historienzeile 1.96 des Handbuchs (Historie, unverändert). Der
Suchlauf prüft Zahlen und Stände, nicht die Vollständigkeit des Musters.

**Belege des Implementers (gemessen am Arbeitsstand auf dem Parent `267a58b6`).**
Host: Linux 6.8.0-139-generic, Docker 29.8.2; Basis-Images wie in den Dockerfiles
gepinnt (`dotnet/sdk:10.0`, `eclipse-temurin:21-jdk`, `python:3.14-slim`,
`postgres:18-alpine` aus `compose.yaml`, `golang:1.27-alpine` für `certgen`); das
`:dev`-Image kommt aus `make image` dieses Laufs (Exit 0, Digest
`sha256:f0116af4e15985b208d1ee8515963c4167d701ee042df9148ddbdaccacbff3b2` in
`harness/image-hash.txt`). Der Server ist unverändert.

- **Unit-Tests** (Docker-Bau der Stufe `build` jedes SDK ohne Cache der Stufe,
  netzlose Loopback-Server, Zertifikate zur Laufzeit, keines im Repo):
  C# `Passed!  - Failed:     0, Passed:   246, Skipped:     0, Total:   246`
  (`Tls/TlsClientTests.cs`: zehn Fälle je vier Flächen — Happy mit Server-Zertifikat
  als Anker, Happy mit **Aussteller** als Anker, Aussteller eines anderen Servers,
  Anker-Datei mit mehreren Zertifikaten, fremder Anker, **Namensfehler**
  (Zertifikat nennt `localhost`, Adresse `127.0.0.1`), abgelaufenes Zertifikat,
  kein Anker, `http://` ohne Anker unverändert, Anker mit `http://` abgelehnt;
  `Tls/TlsOptionsTests.cs`: fünf Fälle zu Pfad und PEM); Kotlin
  `TlsClientTest` 11 Tests, 0 Fehler (`./gradlew test` BUILD SUCCESSFUL); Python
  `242 passed` (`tests/test_tls.py`, dazu die README-Wächter-Tests
  `tests/test_readme_examples.py`, die die neuen README-Codeblöcke gegen die
  Signaturen binden).
- **Mutationsproben** (Zusage · mutierte **Eingabe**/Stelle · gesehenes Rot; an
  Kopien im Scratchpad, Bau der Kopie in Docker, kein Textwerkzeug am Repo; die
  Nenner sind der Testbestand zur Zeit der Probe — C# 238, Kotlin 131, Python 230,
  also vor den zwei Ausstellerfällen):
  - C#, Servername · das Flag `RemoteCertificateNameMismatch` aus der Prüfung
    des Callbacks gestrichen · rot: `ServerNameNotInCertificate_FailsTheConnection`
    × 4 Flächen (4 von 238).
  - C#, Kette/Gültigkeit · das Ergebnis von `chain.Build` verworfen (`return true`)
    · rot: `ForeignAnchor_…` × 4 und `ExpiredCertificate_…` × 4 (8 von 238).
  - C#, Anker wird gelesen · Callback nicht gesetzt; Anker-mit-`http`-Prüfung und
    Leer-Prüfung der PEM-Datei gestrichen (eine Kopie, drei Stellen) · rot:
    `WithTrustAnchor_…` × 4, `AnchorFileWithSeveralCertificates_…` × 4,
    `TrustAnchorWithPlaintextAddress_…` × 4, `EmptyFile_…`, `FileWithoutPemCertificate_…`.
  - Kotlin, Gültigkeit · `ValidityCheckingTrustManager` nicht eingehängt · rot:
    `an expired certificate fails the connection`; die Mutation wurde nicht
    ausgedacht, sondern am grün gebliebenen Test **gesehen**: der Trust-Manager der
    JDK akzeptierte ein abgelaufenes, als Anker eingehängtes Server-Zertifikat.
  - Kotlin, Kette und Name · die Prüfung des Delegaten im Wrapper gestrichen
    (nur Gültigkeit) · rot: `a foreign anchor fails the connection`,
    `a server name that is not in the certificate fails the connection`.
  - Kotlin, Anker wird gelesen · die Anker in `httpClient` bzw. in
    `channelCredentials` verworfen (zwei Kopien) · je rot: `with a trust anchor
    every surface connects over TLS`, `an anchor file with several certificates
    trusts the one that matches` (2 von 131 je Kopie).
  - Kotlin, PEM-Prüfung · die Leer-Prüfung gestrichen (dieselbe Kopie wie die
    Gültigkeits-Mutation oben) · rot: `the anchor option is checked when the
    options are created` (zusammen 2 von 131).
  - Kotlin, **kein Rot**: `endpointIdentificationAlgorithm = ""` am `HttpClient`
    blieb grün (das JDK setzt `HTTPS` selbst); die Zeile ist entfernt, die
    Namensprüfung liegt am Trust-Manager (Mutation „Delegat gestrichen“ oben).
  - Python, Servername · `check_hostname = False` am Kontext des HTTP-Clients und
    `grpc.ssl_target_name_override` am Kanal · rot: `test_server_name_not_in_certificate_…`
    × 4 (dazu `test_trust_anchor_connects_over_tls[grpc_*]`, `…several_certificates[grpc_*]`
    und `test_grpc_address_without_scheme_…`, weil der Override den Namen der
    Verbindung verfälscht; 9 von 230).
  - Python, Anker wird gelesen · `httpx.Client` ohne Kontext und `ssl_channel_credentials()`
    ohne Wurzeln · rot: `test_trust_anchor_connects_over_tls` × 4,
    `…several_certificates` × 4, `test_grpc_address_without_scheme_…` (9 von 230).
  - Python, PEM-Prüfung · `load_verify_locations` aus `_read_trust_anchor` gestrichen
    · rot: `test_anchor_option_is_checked_when_the_options_are_created` (1 von 230);
    die vorher mitgeführte Zählung der Zertifikate war redundant (grün unter
    der Mutation) und ist entfernt.
  - Python, `http`-Prüfungen · beide `raise ValueError` der Fabriken gestrichen ·
    rot: `test_trust_anchor_with_plaintext_address_is_refused` × 4 (4 von 230).
  - **Nicht durch eine Mutation gebunden:** (a) „die Vertrauensanker des
    Betriebssystems gelten bei gesetztem Anker nicht zusätzlich“ — ein
    Test bräuchte ein vom System anerkanntes Zertifikat; der Code setzt in allen drei
    SDKs nur den Anker-Speicher (C# `CustomRootTrust`, Kotlin Schlüsselspeicher
    nur mit den Ankern, Python `cadata` ohne Standard-Speicher), (b) das
    abgelaufene Zertifikat unter Python (der Test ist grün; eine Mutation, die
    allein den Zeitraum abschaltet, wurde nicht gefahren — *hergeleitet*: OpenSSL
    prüft den Zeitraum der Kette), (c) die Ausstellerfälle (`IssuerAsTrustAnchor`) in allen
    drei Sprachen sind nicht mutiert (Stand der Fixrunde: (a) Python HTTP/SSE
    gebunden, (c) für C# teilweise durch die Mutationen der Fixrunde belegt).
- **Fixrunde nach dem Review (F-2, F-3, F-4, F-6), gemessen am Arbeitsstand dieser Runde.**
  - **F-2 (Server-Zertifikat eines Ausstellers als alleiniger Anker).** Je ein Fall
    `Blatt-Zertifikat eines Ausstellers als alleiniger Anker` (Erfolg) und `Blatt eines
    anderen Servers als Anker` (Fehler) auf allen vier Flächen, in allen drei Sprachen.
    **Gemessen vor dem Fix:** Kotlin und Python bestehen den Erfolgsfall unverändert
    (`make sdk-pack-kotlin` Exit 0; `make sdk-pack-python` Exit 0); C# scheitert an
    allen vier Flächen (`Failed:     4, Passed:   250, Skipped:     0, Total:   254`,
    Ursache `PartialChain`, Aussteller fehlt in der Kette). **Fix nur in C#**
    (`TlsTransport.cs`): ist das präsentierte Zertifikat byte-gleich (`RawData`)
    einem Anker, gelten nach der Namens- und Verfügbarkeitsprüfung der Plattform der
    Gültigkeitszeitraum als Prüfung, sonst die Kettenprüfung mit `CustomRootTrust`;
    die Prüfung bleibt an, kein `AllowUnknownCertificateAuthority`. **Danach:** C#
    `Passed!  - Failed:     0, Passed:   254, Skipped:     0, Total:   254`, Kotlin
    `BUILD SUCCESSFUL` (13 Fälle in `TlsClientTest`), Python `252 passed`.
    Mutationen (Kopie im Scratchpad, Bau in Docker, Eingabeseite = das
    präsentierte Zertifikat): Gleichheitsprüfung unwirksam (`false &&`) · rot
    `IssuedServerCertificateAsOnlyTrustAnchor_ConnectsOverTls` × 4 (4 von 254);
    Gültigkeitsprüfung im Gleichheitszweig durch `return true` ersetzt · rot
    `ExpiredCertificate_FailsTheConnection` × 4 (4 von 254); Gleichheitsprüfung
    durch „jeder Anker gleich“ ersetzt (zusammen mit `return true` im Zweig) · rot
    `ForeignAnchor_…`, `IssuedServerCertificateOfAnotherServerAsTrustAnchor_…`,
    `IssuerOfAnotherServerAsTrustAnchor_…` und `ExpiredCertificate_…`, je × 4.
    Die Fälle in Kotlin und Python sind nicht mutiert (die Implementierung
    unverändert, der Fall bestand vor der Änderung; *hergeleitet*: die Anker-Speicher
    von JDK und OpenSSL nehmen ein Blatt als Anker auf).
  - **F-3 (Systemanker gelten bei gesetztem Anker nicht zusätzlich).** Python HTTP und
    SSE gebunden: `test_with_anchor_the_system_trust_does_not_apply[http|sse]` setzt
    `SSL_CERT_FILE` auf das Server-Zertifikat (Kontrolle: ohne Anker verbindet der
    Client), mit fremdem Anker scheitert dieselbe Verbindung; Mutation
    `context.load_default_certs()` nach dem Anker-Kontext · rot (2 von 252 Fällen).
    **Nicht gebunden bleiben:** Python gRPC (der Wurzelspeicher von gRPC liest
    `GRPC_DEFAULT_SSL_ROOTS_FILE_PATH` einmal je Prozess, ein Test im selben Prozess
    wäre reihenfolgeabhängig), C# (`SSL_CERT_FILE` wirkt auf den Wurzelspeicher der
    Plattform, der prozessweit beim ersten Gebrauch gelesen wird) und Kotlin (der
    Standard-`SSLContext` der JDK ist prozessweit initialisiert); ein Test dort wäre
    ein Scheintest oder bräuchte einen Kindprozess je Fall. Der Code setzt in diesen
    Fällen nur den Anker-Speicher (*hergeleitet*, wie oben).
  - **F-4 (Verhaltensänderungen, Kompat-Messung).** `make test-sdk-kompat` (Standard,
    `SDK_KOMPAT_NEU=dist`) endet rot (Exit 2) an der Auflösung des Pakets: das Ziel
    erwartet `0.6.0` in `sdks/csharp/dist`, dort liegt `0.6.1`
    (`NU1603 … PgChangeFeed.Client 0.6.0 was not found. … 0.6.1 was resolved instead`);
    die Ursache steht seit dem Versionsstand `0.6.1` und liegt nicht an dieser
    Änderung. `SDK_KOMPAT_NEU=registry` endet mit Exit 0 (`Kompatibilitätsmessung
    (registry) grün für: csharp kotlin python`) und misst die veröffentlichten
    Pakete, **nicht** den Arbeitsstand; das neue C#-Mapping (`Internal` mit
    `HttpRequestException` → `UnexpectedStatus`) und `https` → TLS in Kotlin sind
    damit von dem Ziel **nicht gemessen** (das Ziel liest die Fehlertypen, nicht
    das Mapping), ihre Belege sind die Unit-Fälle `ExpectConnectionFailureAsync`
    bzw. `expectConnectionFailure`. Beide Änderungen stehen als Eintrag
    „Upgrading“ in den READMEs von C# und Kotlin (Überschrift `0.7.0`, die nächste
    Minor-Version laut `SPEC-037`; die Package-Versionen bleiben bis zum Release-Zug).
  - **F-6 (Kotlin-Realserver, Gegenlauf).** `SDK_TLS_CA_FILE=/tls/other.pem make
    test-sdk-kotlin-integration` · Exit 2; rot in der HTTP-TLS-Phase:
    `registersAConsumerAndListsTablesOverTlsWithTheTrustAnchor() FAILED
    javax.net.ssl.SSLHandshakeException` (die Ablehnungs-Fälle der Phase blieben
    `REJECTED tls no_anchor=ok foreign_anchor=ok name_mismatch=ok`). Die
    Ausstellerfälle am Realserver sind nicht gefahren: `certgen` erzeugt
    selbstsignierte Zertifikate, ein Ausstellerfall bräuchte eine Erweiterung der
    Fixture; sie sind in den Unit-Fällen aller drei Sprachen gebunden.
  - **Realserver, Wiederholung nach dem Fix** (`make image` Exit 0 zuvor): die drei
    `make test-sdk-*-integration` je Exit 0, die TLS-Zeilen wie oben
    (C# `csharp-sdk-tls-20261005064620`, Kotlin `kotlin-sdk-tls-20261005065038`,
    Python `python-sdk-tls-4acfb622a1fd`).
- **Realserver** (`make image` zuvor; je Runner ein Lauf, Exit 0; die Zeile ist die
  gedruckte Abschlusszeile des Runners, `change_id`/`consumer_id` je unabhängig
  über `cdc.changes` bzw. `cdc.consumer` gegengelesen):
  - C#, `make test-sdk-csharp-integration`: `run-sdk-csharp-integration-tests: TLS-Belege
    (LH-FA-SST-013, ADR-0150) grün — … HTTP (consumer_id=csharp-sdk-tls-20261005055841),
    SSE (change_id=954-1), gRPC-Stream (change_id=957-1) und gRPC-Verwaltung
    (consumer_id=csharp-sdk-tls-admin-20261005055847) arbeiteten über TLS mit
    TrustAnchorFile …; ohne Anker, mit fremdem Anker und mit einem Servernamen
    außerhalb des Zertifikats scheiterte dieselbe Verbindung an der TLS-Prüfung`.
  - Kotlin, `make test-sdk-kotlin-integration`: `… TLS-Belege … HTTP
    (consumer_id=kotlin-sdk-tls-20261005055447), SSE (change_id=994-1), gRPC-Stream
    (change_id=998-1) und gRPC-Verwaltung (consumer_id=kotlin-sdk-tls-admin-20261005055518)
    … mit trustAnchorFile …`.
  - Python, `make test-sdk-python-integration`: `… TLS-Belege … HTTP
    (consumer_id=python-sdk-tls-cba8b6dfcd75), SSE (change_id=945-1), gRPC-Stream
    (change_id=948-1) und gRPC-Verwaltung (consumer_id=python-sdk-tls-admin-244a1bf1ff2d)
    … mit trust_anchor_file …`.
  - Mutationsläufe am Realserver (Eingabeseite: der Pfad des Ankers im
    Test-Container; `SDK_TLS_FOREIGN_CA_FILE=/tls/good.pem` macht den „fremden“
    Anker zum richtigen, `SDK_TLS_CA_FILE=/tls/other.pem` den richtigen zum
    fremden): C# beide Läufe, Python beide Läufe, Kotlin der Lauf mit dem fremden
    Anker; jeder Lauf endete mit Exit 2 von `make` (Exit 1 des Runners) in der
    HTTP-TLS-Phase — mit dem richtigen Anker als fremdem: „der Ablehnungs-Beleg
    blieb aus (REJECTED tls … fehlt)“, mit dem fremden als richtigem: „der Test
    empfing keine der committeten Änderungen“ bzw. „lieferte den Happy-Path-Marker
    nicht“. Der Kotlin-Lauf mit `SDK_TLS_CA_FILE=/tls/other.pem` steht in der Fixrunde (F-6).
- **Pakete:** `make sdk-pack-csharp`, `make sdk-pack-kotlin`, `make sdk-pack-python`
  je Exit 0; Artefakte `PgChangeFeed.Client.0.6.1.nupkg`,
  `pgchangefeed-kotlin-0.6.1.jar` und `…-0.6.1-sources.jar`,
  `pgchangefeed-0.6.1-py3-none-any.whl` und `pgchangefeed-0.6.1.tar.gz` — die
  Versionen sind absichtlich nicht angehoben (siehe Tabelle oben).
- **Gates:** `make gates` am Endstand dieses Arbeitszuges Exit 0 (Exit direkt
  gelesen, nicht durch eine Pipe; der erste Lauf endete an einem Befund von
  `docs-check` — eine nackte Kennung in diesem Plan —, behoben und wiederholt);
  `make sdk-public-doc-check` und `make handbuch-public-doc-check` Exit 0 für sich
  gefahren; `make suchlauf-nachmessen PLAN=…` stimmt in allen sechs Zeilen.
  `make fmt-check` entfällt: der Zug trägt keine Go-Datei (`tools/harness/certgen`
  ist unverändert).
- **Kommentar-Läufe (Schritt 20):** der Chronik-Kandidatenlauf über die geänderten
  und neuen `.cs`/`.kt`/`.py`/`.sh`/`.mk`-Dateien und der Konjunktiv-Lauf über die
  hinzugefügten Zeilen: 0 Treffer; `make kommentar-kennungen DIFF=HEAD` Exit 0
  ohne Kandidat (Form, nicht Wahrheit).

## 4. Trigger

**Start** (`next` → `in-progress`): Lastenheft 0.16.0 und `SPEC-037` committet;
der parallele Zug an `tools/harness` ist geschlossen oder berührt die Runner
nicht; ein geladenes `:dev`-Image (`make image`).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): die drei Liefer-Punkte tragen nicht in einem
  Review — dann je Sprache ein Slice (`slice-sdk-tls-optionen-csharp`, `…-kotlin`,
  `…-python`), dieser Plan wird zum ersten, die beiden anderen entstehen in `open/`.
- `in-progress` → `open` (blockiert): eine Sprache kann den Anker ohne
  Abschalten der Prüfung nicht an HTTP **und** gRPC binden — dann Befund an
  `SPEC-037`, Folge-ADR.

## 5. Closure-Trigger

DoD vollständig, `make gates` grün, Review- und Verifier-Bericht liegen vor,
Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- Die JDK-/Java-HTTP-Schicht und das Python-SSE brauchen jeweils einen eigenen
  Weg zum Trust-Store — **Ausgang:** weiter offen bis zur Umsetzung; nicht lösbar
  heißt Rückführung (§4).
- Die Realserver-Runner liegen unter `tools/harness`, das ein anderer Zug gerade
  ändert — **Ausgang:** weiter offen; entfällt, wenn der Slice nach dessen
  Abschluss startet.
- Der Suchlauf trifft Symbolnamen, nicht verschobene Zahlen (`AGENTS.md` §3.13
  Grenze) — **Ausgang:** Lese-Handlung des Reviewers.

## 7. Closure-Notiz

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag:** —
- **Beobachtungs-Register:** —
- **Folge-Slices:** —
- **Risiken aus §6:** —
- **Drei Paarungen:** —

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist `*` (Default; `sdks/`, `tools/harness`,
`docs/user/`); eine Ausdifferenzierung ist nicht nötig.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen; kein Treffer
für TLS in den SDKs.

Alle berührten Sub-Areas GF.
