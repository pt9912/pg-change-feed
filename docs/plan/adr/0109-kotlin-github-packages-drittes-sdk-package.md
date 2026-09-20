# ADR-0109: Kotlin/GitHub Packages als drittes SDK-Package für `LH-FA-SST-009`

**Status:** Accepted

**Datum:** 2026-09-20

**Autor:** pt9912 (Architect-Rolle, Nutzerentscheidung vom 2026-09-20 im Chat:
„Können wir noch ein kotlin-SDK erstellen" — die Kern-Entscheidung „Kotlin als
dritte Sprache"; eine zweite Nutzerentscheidung im selben Zug, nach
Vorstellung von Maven Central als Erst-Kandidat: „Anstatt maven können wir
github verwenden" — die Vertriebsweg-Entscheidung. Diese ADR fasst beide
vollwertig, wählt den Vertriebsweg begründet gegen echte Alternativen
(einschließlich des zunächst vorgeschlagenen, dann vom Nutzer verworfenen
Maven Central) und entscheidet die daran hängenden Sub-Fragen selbst, mit
Alternativenvergleich, Modul 8 §Rollen-Regeln: „Architect schreibt")

**Bezug:** [`LH-FA-SST-009`](../../../spec/lastenheft.md) (Client-Bibliotheken/
SDKs — die Lastenheft-Fähigkeit selbst), [`LH-FA-SST-006`](../../../spec/lastenheft.md)
(HTTP-/JSON-API), [`LH-FA-SST-008`](../../../spec/lastenheft.md) (Live-Change-Stream,
gRPC/SSE/NATS-Vollinhalt), [`ARC-005`](../../../spec/architecture.md)
(Driving Adapters — die Oberflächen, die das SDK als Consumer anspricht),
[`ADR-0106`](0106-csharp-nuget-erstes-sdk-package.md) (erstes SDK-Package,
C#/NuGet — Formvorbild für den Umfangs-/Struktur-Schnitt dieser ADR),
[`ADR-0107`](0107-python-pypi-zweites-sdk-package.md) (zweites SDK-Package,
Python/PyPI — Formvorbild für die begründete Vertriebsweg-Wahl statt
Übernahme; beide von dieser ADR **unberührt**: `§Re-Evaluierungs-Trigger`
Punkt 1 beider ADRs sieht eine weitere Sprache/einen weiteren Vertriebsweg
ausdrücklich als eigene ADR vor), [`ADR-0108`](0108-python-sdk-uv-statt-build-twine.md)
(Frontend-Wechsel-Präzedenzfall — Beleg-Pflicht bei Werkzeug-Reifegrad,
`AGENTS.md` §3.12), [`ADR-0057`](0057-http-grpc-api.md)
(HTTP/JSON-API-Vertrag), [`ADR-0060`](0060-grpc-streaming-mechanismus.md)
(gRPC-Server-Streaming und die `.proto`), [`ADR-0061`](0061-http-sse-zusaetzlich-zu-grpc.md)
(SSE), [`ADR-0100`](0100-nats-dritter-vollinhalts-zustellweg.md)
(NATS-Vollinhalts-Stream), [`ADR-0076`](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
(`examples/` — Ort, Import-Grenze, Rolle „Vorbild, kein Belegträger"),
[`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md) und
[`ADR-0090`](0090-beispiel-clients-volle-matrix.md) (Beispiel-Clients
C#/Kotlin — Präzedenzfall für Sprach-Wurzel-Form, Docker-only-Bauform,
digest-/paket-gepinnte Gradle-Werkzeugkette; regeln ausdrücklich nur
`examples/`, keine veröffentlichten Pakete), [`ADR-0098`](0098-beispiel-clients-start-ueber-make-dockerfile.md)
(Startform der Beispiele — unberührt von dieser ADR),
[`ADR-0068`](0068-wegwerf-clients-begrenzte-import-berechtigung.md)
(Import-Berechtigungs-Disziplin für einen öffentlichen Pfad-Bereich),
[`ADR-0051`](0051-cicd-pipeline-github-actions.md) Entscheidung 3/8
(Release-Tag-Trigger, Registry-Secret-Muster — Kontrastfolie: dort brauchte
jeder Vertriebsweg ein externes Registry-Secret, hier ausdrücklich **nicht**),
[`AGENTS.md`](../../../AGENTS.md) §3.1 (Docker-only) §3.5
(Accepted-ADR-Immutabilität) §3.6 (Gates nur per ADR gelockert) §3.12
(Herkunft von Aussagen — die Vorarbeits-Parität, die Gradle/JDK-Recherche
und die GitHub-Packages-Auth-Recherche unten tragen ihren Beleg) §3.13
(Träger-Nachzug), [`spec/pflichtenheft.md`](../../../spec/pflichtenheft.md)
§2 `SPEC-018`/`SPEC-020`/`SPEC-021`/`SPEC-024` (die vier Draht-Festlegungen,
die das SDK **benutzt**), `SPEC-026`/`SPEC-027` (die bereits existierenden
C#-/Python-Packages — Formvorbild für einen künftigen Kotlin-Sprach-Eintrag),
`docs/user/version.md` (Server-SemVer — von dieser ADR ausdrücklich **nicht**
berührt)

**Schärft:** [`LH-FA-SST-009.a`](../../../spec/pflichtenheft.md) — dieselbe
Pflichtenheft-Stelle, die [`ADR-0106`](0106-csharp-nuget-erstes-sdk-package.md)
für C#/NuGet und [`ADR-0107`](0107-python-pypi-zweites-sdk-package.md) für
Python/PyPI beantwortet hat, bleibt nach ihrem eigenen Wortlaut für „eine
dritte Sprache oder einen dritten Vertriebsweg" offen; diese ADR beantwortet
sie jetzt für Kotlin/GitHub Packages als **dritte** Sprache/dritten
Vertriebsweg, und für die daran hängenden Umsetzungsfragen (Umfang, Ort,
Versionierung, Bau-/Publish-Mechanismus). Eine vierte Sprache bleibt eine
eigene, künftige ADR (§Re-Evaluierungs-Trigger).

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

**Die Anforderung.** [`LH-FA-SST-009`](../../../spec/lastenheft.md) verlangt
offizielle, versionierte Client-Bibliotheken (Packages) für die bestehenden
Zustellwege — HTTP-API (`SPEC-018`), gRPC-Stream (`SPEC-020`), SSE
(`SPEC-021`), NATS-Vollinhalts-Stream (`SPEC-024`). Für die **erste** und
**zweite** Sprache haben [`ADR-0106`](0106-csharp-nuget-erstes-sdk-package.md)
(C#/NuGet, `SPEC-026`) und [`ADR-0107`](0107-python-pypi-zweites-sdk-package.md)
(Python/PyPI, `SPEC-027`) diese Frage bereits beantwortet.

**Die Nutzerentscheidung — in zwei Schritten.** Der Auftraggeber hat am
2026-09-20 im Chat zunächst entschieden: „Können wir noch ein kotlin-SDK
erstellen." Das legt die **Sprache** (Kotlin) für das **dritte** Package
fest. Ein erster Architect-Entwurf dieser ADR schlug Maven Central (Central
Portal) als Vertriebsweg vor, begründet gegen Alternativen abgewogen. Der
Auftraggeber hat diesen Vorschlag darauf explizit **abgelehnt**: „Anstatt
maven können wir github verwenden." Das ist die zweite, für den
Vertriebsweg maßgebliche Nutzerentscheidung. Diese ADR trägt beide
Entscheidungen und wägt GitHub Packages ehrlich gegen Maven Central und die
übrigen Alternativen ab (§Verglichene Alternativen B) — sie übernimmt die
Nutzeräußerung nicht unbesehen, sondern prüft sie, wie
[`ADR-0107`](0107-python-pypi-zweites-sdk-package.md) es für PyPI bereits
vorgemacht hat.

**Der Ist-Stand — gemessen, nicht erinnert (heutiger Zug).**

| Prozedur | Ergebnis |
|---|---|
| `find examples/kotlin -maxdepth 1 -type d` | fünf Programme: `http-client`, `sse-client`, `nats-client` (Wecksignal, `SPEC-017`), `grpc-client`, `nats-stream-client` (Vollinhalt, `ADR-0100`) |
| `find examples/csharp -maxdepth 1 -type d` (zum Vergleich, heute erneut gemessen — **nicht** aus `ADR-0106` übernommen) | dieselben fünf Programme: `http-client`, `grpc-client`, `sse-client`, `nats-client`, `nats-stream-client` |
| `find sdks -maxdepth 2 -type d` | `sdks/csharp/PgChangeFeed.Client`, `sdks/python/pgchangefeed` — kein `sdks/kotlin`, keine Namenskollision |
| `cat examples/kotlin/gradle/wrapper/gradle-wrapper.properties` | Gradle `8.14`, `distributionSha256Sum` gepinnt |
| `cat examples/kotlin/build.gradle.kts` | Kotlin-Gradle-Plugin `2.4.20` (gemessen 2026-09-17, Maven Central), `com.google.protobuf`-Gradle-Plugin `0.10.0` (Gradle Plugin Portal) |
| `cat examples/kotlin/Dockerfile` | Basis `eclipse-temurin:21-jdk`/`21-jre`, digest-gepinnt (gemessen 2026-09-17) |
| `.a-check.yml` `languages:` | nur `go` — eine Kotlin-Quelle wird von a-check strukturell nicht gelesen, unabhängig vom Ort (dieselbe Grenze wie bei C#/Python) |
| `cat docs/user/version.md` | `0.1.2` — Server-SemVer, unverändert eigenständig geführt |
| GitHub-Repository dieses Projekts | `pt9912/pg-change-feed` — trägt die Adresse des GitHub-Packages-Endpunkts (§Entscheidung Festlegung 2) |

**Was das ändert — der zentrale Unterschied zu `ADR-0107`, die genaue
Einordnung ggü. `ADR-0106`.** Anders als bei Python (kein
`examples/python/`, `ADR-0107` §Kontext) existiert für Kotlin bereits ein
**vollständiges** Fünf-Programm-Beispiel-Set — dieselbe Anzahl und
dieselbe Abdeckung aller vier `LH-FA-SST-009`-relevanten Zustellwege
(HTTP, gRPC, SSE, NATS-Vollinhalt) plus das NATS-Wecksignal, die
`ADR-0106` für C# beim Schreiben dieser ADR bereits vorfand. Das ist eine
**Präzision, keine Übernahme** der Aufgabenstellung: Die heutige Messung
zeigt **Parität** mit C# — nicht mehr Vorarbeit als C#, sondern **exakt so
viel**. Für den Umfangs-Schnitt (§Entscheidung Festlegung 1) heißt das: Das
Wire-Verständnis-Risiko ist für Kotlin so niedrig wie es für C# bei
`ADR-0106` war (nicht niedriger) — der dortige Scope-Schnitt (HTTP + gRPC in
einem Erst-Release, SSE/NATS als Folge) war dort **ausschließlich** eine
Frage des Paket-Zuschnitts (Abhängigkeits-Footprint, Erst-Release-Größe),
nicht des Wire-Risikos (`ADR-0106` §Kontext „Was das ändert"). Dieselbe
Trennung gilt hier: Die Footprint-Frage ist von der Vorarbeits-Frage
**unabhängig**, und die Parität mit C# ist ein Grund, `ADR-0106`s
Scope-Logik direkt zu übertragen — kein Grund, sie zu übertreffen (kein
größerer Erst-Scope allein aus „es gibt ja schon alle fünf Beispiele");
kein Grund, sie zu unterschreiten wie bei Python (kein fehlendes Vorbild,
das eine zusätzliche Vorsicht rechtfertigen würde).

**Fünf Bindungen, die den Lösungsraum prägen — analog `ADR-0106`/`ADR-0107`,
für Kotlin neu geprüft:**

1. **Docker-only** ([`AGENTS.md`](../../../AGENTS.md) §3.1): kein
   Host-`gradle`/`kotlinc`/`java`. Bau, Test und Paketierung laufen im
   gepinnten JDK-Container, analog `make examples-kotlin` und
   `make sdk-pack-csharp`/`make sdk-pack-python`.
2. **`examples/` bleibt Vorbild, kein Belegträger** (`SPEC-023`) — hier mit
   praktischer Wirkung, weil `examples/kotlin/` real existiert: dieselbe
   räumliche Trennungslogik wie bei C# (`ADR-0106` §Entscheidung
   Festlegung 2), kein Konflikt, aber ein realer Abgrenzungsbedarf
   (§Entscheidung Festlegung 3, §Verglichene Alternativen D).
3. **Kein neues Gate.** GitHub-Packages-Publish braucht Netz (§Entscheidung
   Festlegung 5) — dieselbe Klasse wie NuGet-/PyPI-Publish: Werkzeug, kein
   Gate.
4. **a-check liest kein Kotlin.** `sdks/kotlin/**` braucht keine
   `.a-check.yml`-Gruppe — dieselbe Grenze wie `sdks/csharp/**`/
   `sdks/python/**`.
5. **Eine Arbeit, die eine beschriebene Eigenschaft bewegt, zieht ihre
   Träger nach** ([`AGENTS.md`](../../../AGENTS.md) §3.13):
   `LH-FA-SST-009.a`s Satz „eine dritte Sprache … bleibt offen" wird mit
   dieser ADR falsch — der Träger-Nachzug ist Sache des
   Implementierungs-Zuges (§Konsequenzen Folgepflicht), nicht dieser
   ausschließlich ADR-schreibenden Rolle.

**Recherche zum GitHub-Packages-Publish-/Auth-Mechanismus — real
durchgeführt, mit Beleg (`AGENTS.md` §3.12, keine Behauptung ohne Anker).**

- **Live abgerufen** (`curl`, heute,
  `docs.github.com/…/publishing-java-packages-with-gradle`): Das
  dokumentierte Beispiel-Workflow verwendet **ausschließlich** das
  eingebaute `GITHUB_TOKEN` — kein separates Repository-Secret — mit den
  Workflow-Berechtigungen `permissions: contents: read` / `packages:
  write`. Wörtlich: „`GITHUB_TOKEN` to publish packages associated with the
  workflow repository." Für **Publish** ist damit **kein einziges**
  zusätzliches Secret nötig — ein echter, dokumentiert bestätigter
  Unterschied zu NuGet (`NUGET_API_KEY`), PyPI (`PYPI_API_TOKEN`) und dem
  ursprünglich erwogenen Maven-Central-Weg (fünf Werte, siehe unten).
- **Live abgerufen**, dieselbe Seite: das reale Gradle-`build.gradle`-
  Fragment für GitHub Packages lautet
  `maven { name = "GitHubPackages"; url =
  "https://maven.pkg.github.com/<OWNER>/<REPOSITORY>"; credentials {
  username = System.getenv("GITHUB_ACTOR"); password =
  System.getenv("GITHUB_TOKEN") } }` — mit dem **eingebauten**,
  Gradle-Kern-Plugin `id 'maven-publish'`, **kein** zusätzliches
  Dritt-Plugin.
- **Live abgerufen** (`curl`, heute,
  `docs.github.com/…/working-with-the-apache-maven-registry`): Die
  Distributions-URL folgt dem Muster `https://maven.pkg.github.com/<OWNER
  >/<REPOSITORY>` — für dieses Repository also
  `https://maven.pkg.github.com/pt9912/pg-change-feed`. Ein Repository
  kann mehrere Packages tragen, solange die `distributionManagement`-URL
  (Maven) bzw. die `publishing.repositories.maven.url` (Gradle) explizit
  gesetzt ist — keine 1:1-Bindung „ein Package pro Repository" (real
  bestätigt, entkräftet einen möglichen Einwand gegen Mehrfach-Nutzung
  desselben Repos für C#-/Python-/Kotlin-Metadaten, auch wenn nur das
  Kotlin-Package GitHub Packages als Vertriebsweg nutzt).
- **Live abgerufen** (`curl`, heute,
  `docs.github.com/…/working-with-the-gradle-registry`): **Der real
  bestätigte, benannte Nachteil.** Wörtlich: „You need an access token to
  publish, install, and delete **private, internal, and public**
  packages." — GitHub Packages verlangt **immer** eine Authentifizierung
  zum Lesen, unabhängig von der Sichtbarkeit des Pakets, und ausschließlich
  über ein „personal access token (classic)" (keine anonyme, token-lose
  Installation wie bei Maven Central/NuGet/PyPI). Für einen Kotlin-Consumer
  außerhalb dieses Repositories bedeutet das: ein eigenes GitHub-Konto,
  ein klassischer PAT mit `read:packages`-Scope, und eine
  Repository-/`~/.m2/settings.xml`-Konfiguration mit Zugangsdaten, **auch**
  wenn das SDK öffentlich und quelloffen ist. Dieser Punkt wird in
  §Verglichene Alternativen B und §Konsequenzen **nicht verschwiegen** —
  er ist der reale Preis der heutigen Nutzerentscheidung.
- **Live abgerufen**, dieselbe Recherche: der Central-Portal-Weg (Maven
  Central), den ein früherer Entwurf dieser ADR vorschlug, hätte
  demgegenüber **fünf** Secrets gebraucht (`mavenCentralUsername`/
  `mavenCentralPassword` für den Central-Portal-User-Token,
  `signingInMemoryKey`/`signingInMemoryKeyId`/`signingInMemoryKeyPassword`
  für die von Maven Central verlangte GPG-Signatur) plus eine vorherige
  Namespace-Registrierung (`io.github.pt9912`, real ohne DNS verifizierbar,
  aber ein separater, externer Registrierungsschritt) — GitHub Packages
  braucht **keins** davon: kein GPG-Schlüsselpaar, kein
  Sonatype-Central-Portal-Konto, kein Namespace-Antrag. Das ist der reale
  Kern der Vereinfachung, die der Nutzer mit „Anstatt maven können wir
  github verwenden" gewählt hat.

**Recherche zur JDK-/Gradle-Version — real durchgeführt (analog der
Python-3.14-Reverifikation in `ADR-0107`/`ADR-0108`), unverändert von der
Vertriebsweg-Frage.**

- **Live abgerufen** (`curl`, heute, Adoptium-API
  `api.adoptium.net/v3/info/available_releases`): aktuellste LTS-Version
  ist JDK **25** (`most_recent_lts: 25`), nicht mehr JDK 21 — `examples/kotlin`
  ist mit JDK 21 auf der **vorletzten** LTS-Zeile gepinnt (gemessen
  2026-09-17).
- **Live abgerufen** (`curl`, heute, `docs.gradle.org/current/userguide/compatibility.html`,
  Tabelle „Java Compatibility"): Gradle **8.14** — dieselbe Version, die
  `examples/kotlin/gradle/wrapper/gradle-wrapper.properties` bereits
  einsetzt — unterstützt JDK-Toolchains bis einschließlich **24**; JDK **25**
  als Toolchain-Ziel verlangt laut derselben Tabelle **Gradle 9.1.0**.
  Ein Sprung auf JDK 25 für das neue SDK würde also **zusätzlich** einen
  Gradle-Versionssprung erzwingen, der in diesem Repo bisher **ungeprüft**
  ist — eine neue, ungeprüfte Kombination für ein Erst-Release, das gerade
  **von** der Parität mit einer bereits erprobten Gradle/Kotlin-Werkzeugkette
  profitieren soll (§Entscheidung Festlegung 5).

## Entscheidung

Wir wählen: **Kotlin als dritte SDK-Sprache, GitHub Packages (Gradle-/Maven-
Registry dieses Repositories) als Vertriebsweg.** Sechs Festlegungen:

### 1 — Umfang des ersten Pakets: HTTP-API und gRPC-Stream, ein Package —
wie bei C#, nicht größer

v1 des Packages (Arbeitsname `pgchangefeed-kotlin`, Koordinate
`io.github.pt9912:pgchangefeed-kotlin`) deckt **HTTP-API** (`SPEC-018`) und
**gRPC-Stream** (`SPEC-020`) — nicht SSE (`SPEC-021`) und nicht
NATS-Vollinhalt (`SPEC-024`), obwohl für **alle fünf** Kotlin-Beispielprogramme
bereits real funktionierende Referenzen existieren (§Kontext, Ist-Stand).
Begründung, gegen den größeren Zuschnitt abgewogen (§Verglichene
Alternativen C):

- **Die Vorarbeits-Parität mit C# ist der tragende Grund, `ADR-0106`s
  Scope-Logik zu übertragen — nicht sie zu übertreffen.** `ADR-0106` §Kontext
  „Was das ändert" trennt bereits explizit Wire-Verständnis-Risiko
  (durch ein Vorbild gesenkt) von Paket-Design-Kosten (öffentliche
  API-Form, Versionsgrenze, Dokumentation, Tests als Vertrag — fallen **je
  Zustellweg** an, unabhängig vom Vorbild). Kotlin hat heute exakt die
  Vorbild-Parität, die C# bei `ADR-0106` hatte — kein Grund, hier eine
  andere Kosten-Rechnung anzusetzen.
- **Abhängigkeits-Footprint ist unverändert ein echter Grund, unabhängig von
  der Vorarbeit und unabhängig vom Vertriebsweg.** HTTP braucht in Kotlin
  keine schwere Fremdabhängigkeit (`java.net.http`, siehe
  `examples/kotlin/http-client/build.gradle.kts`). gRPC zieht bereits real
  gemessen mehrere Koordinaten (`io.grpc:grpc-kotlin-stub`,
  `io.grpc:grpc-netty-shaded`, `io.grpc:grpc-bom`,
  `com.google.protobuf:protobuf-java`,
  `org.jetbrains.kotlinx:kotlinx-coroutines-core` —
  `examples/kotlin/grpc-client/build.gradle.kts`). SSE und
  NATS-Vollinhalt (`io.nats:jnats`, siehe
  `examples/kotlin/nats-stream-client`) legen **weitere, eigenständige**
  Fremdmodule obendrauf — ein Erst-Release, das alle vier bündelt, zwingt
  jedem HTTP-only-Consumer drei Fremdmodul-Gruppen auf, die er nie lädt.
- **HTTP + gRPC bleibt die Kombination, die die meisten Consumer
  brauchen** (`ADR-0106` §Entscheidung Festlegung 1, unverändert gültig):
  synchroner Abruf/Backfill und Live-Zustellung über den Transport mit dem
  stärksten Typsicherheits-/Tooling-Vorteil für ein JVM-Ökosystem
  (generierter Kotlin-Coroutine-Stub) — SSE und NATS sind **alternative**
  Transporte für dieselbe Live-Zustellungs-Fähigkeit
  ([`LH-FA-SST-008`](../../../spec/lastenheft.md) Out-of-Scope), kein
  zusätzlicher fachlicher Nutzen gegenüber gRPC für denselben Consumer.
- **Boundary-Akzeptanzkriterium bleibt erfüllt, ohne v1 zu bläuen:**
  SSE/NATS bleiben über `examples/kotlin/sse-client`/`nats-stream-client`
  als **Vorbild** (nicht Ersatz) einsehbar, bis ein Folge-Paket sie deckt —
  dieselbe Boundary-Erfüllung wie bei C#.
- **Ein Package, nicht vier/fünf.** Dieselbe Konsolidierungs-Begründung wie
  bei C# (`ADR-0106` §Entscheidung Festlegung 1, letzter Punkt): ein
  Consumer, der beide Wege nutzt, referenziert eine Maven-Koordinate statt
  mehrerer.

SSE und NATS-Vollinhalt bleiben für ein **Folge-Package** offen — dieselbe
Formulierung wie in `ADR-0106`, jetzt zum dritten Mal wiederholt, weil sie
sich als tragfähiges Muster bewährt (kein Vorgriff auf den konkreten
Umsetzungs-Zug hier).

### 2 — Vertriebsweg: GitHub Packages (Gradle-/Maven-Registry dieses
Repositories)

GitHub Packages ist der Vertriebsweg, unter der Koordinate
`io.github.pt9912:pgchangefeed-kotlin`, veröffentlicht in die
Registry-Adresse `https://maven.pkg.github.com/pt9912/pg-change-feed`
(real recherchiertes URL-Muster, §Kontext). Begründung — Nutzerentscheidung,
gegen echte Alternativen geprüft, nicht unbesehen übernommen
(§Verglichene Alternativen B):

- **Radikal einfacherer Publish-Mechanismus als jede geprüfte Alternative.**
  Kein GPG-Schlüsselpaar, kein Sonatype-Central-Portal-Konto, kein
  Namespace-Antrag, **kein einziges** zusätzliches Repository-Secret — das
  eingebaute `GITHUB_TOKEN` mit `packages: write`-Berechtigung genügt
  (real dokumentiert, §Kontext Recherche). Das ist einfacher als NuGet
  (ein externes Secret) und PyPI (ein externes Secret), und **erheblich**
  einfacher als der ursprünglich erwogene Maven-Central-Weg (fünf externe
  Werte plus Namespace-Registrierung).
- **Liegt bereits im selben Hosting wie Quellcode und CI** — keine neue
  externe Account-Beziehung (Sonatype, npm, PyPI) für dieses Package; die
  Registry-Adresse folgt direkt aus dem bereits bestehenden Repository
  `pt9912/pg-change-feed`.
- **Der reale Preis wird nicht verschwiegen: GitHub Packages verlangt eine
  Authentifizierung zum Lesen, auch für ein öffentliches Package** (real
  dokumentiert, §Kontext Recherche — anders als Maven Central, NuGet und
  PyPI, die anonymen Download erlauben). Ein Kotlin-Consumer braucht ein
  GitHub-Konto und einen klassischen PAT mit `read:packages`-Scope, den er
  in seiner eigenen Build-Konfiguration hinterlegt. Das ist ein **echter**
  Reibungspunkt gegenüber Maven Central — die explizite
  Nutzerentscheidung nimmt ihn bewusst in Kauf, zugunsten eines
  drastisch einfacheren Erst-Release-Aufwands (§Re-Evaluierungs-Trigger 5
  hält den Rückweg offen, falls sich das als echtes Consumer-Problem
  erweist).
- **`LH-FA-SST-009`s AC verlangt „über den Paketmanager seiner Sprache
  einbinden, ohne das Protokoll selbst zu implementieren"** — eine
  Authentifizierungspflicht beim Bezug verletzt diese AC nicht: Der
  Consumer implementiert weiterhin nicht das Draht-Protokoll selbst, er
  konfiguriert nur zusätzlich Zugangsdaten für die Registry, ein in der
  Gradle-/Maven-Welt (z. B. bei jeder privaten Firmen-Registry) übliches
  Muster.

### 3 — Verhältnis zu `examples/kotlin/`: eigenständiges Projekt, neuer Baum

Das Package entsteht unter `sdks/kotlin/pgchangefeed-kotlin/`, analog zu
`sdks/csharp/PgChangeFeed.Client/` und `sdks/python/pgchangefeed/`. Begründung
— hier mit realem Konfliktpotential, weil `examples/kotlin/` real existiert
(anders als Python):

- **`SPEC-023` markiert `examples/` als „Vorbild, kein Belegträger, keine
  Zustandsmaschine"** — ein Paket-Konsument verlässt sich auf Kompatibilität
  über Versionsgrenzen hinweg; das ist die Definition eines Belegträgers.
  `examples/kotlin/{http,grpc,sse,nats,nats-stream}-client` bleiben genau
  das, was sie heute sind: CLI-Demonstrationsprogramme mit eigenem
  `application`-Plugin-Entrypoint, keine öffentliche Bibliotheks-API.
- **`sdks/` ist bereits die etablierte Sprach-Wurzel-Konvention** für
  veröffentlichte SDK-Packages (`grep -rn "sdks/" .` zeigt
  `sdks/csharp/**` und `sdks/python/**`) — ein neues Sprach-Verzeichnis
  unter derselben Wurzel ist die dritte Anwendung derselben Konvention,
  keine neue.
- **Kein `project(":…")`-Abhängigkeitspfad auf `examples/kotlin/*`.** Die
  vorhandenen Beispiel-Klassen bleiben **Referenzmaterial** für den
  Implementierungs-Zug — dieselbe Draht-Kenntnis, aber als eigenständiger,
  paketierbarer Code neu geschrieben, mit einer öffentlichen, stabilen API
  (nicht CLI-Argument-Parsing wie die Beispiele).
- **a-check bleibt unberührt** (§Kontext Bindung 4): `sdks/kotlin/**` ist
  keine Go-Quelle.
- **Der gRPC-Teil erreicht die `.proto` genauso wie
  `examples/kotlin/grpc-client`:** über einen zusätzlichen, benannten
  Docker-Bau-Kontext (`--build-context proto=proto`) — die `.proto` bleibt
  die einzige Quelle des Draht-Vertrags, kein committeter Stub im
  SDK-Baum.
- **Import-Grenze wie bei C#/Python, verschärft:** Das SDK importiert
  ausschließlich die Kotlin-/JVM-Standardbibliothek und öffentliche
  Maven-Koordinaten, keinen privaten Baum dieses Repositories.

### 4 — Versionierung: unabhängiges SemVer, ab `0.x.y`

Das Package führt sein **eigenes** SemVer 2.0 — das im JVM-/Gradle-/Maven-
Ökosystem übliche Schema, nicht Pythons PEP 440 (`ADR-0107` Festlegung 4
galt nur für Python, aus einem PyPI-spezifischen Grund) —, unabhängig von
`docs/user/version.md` und unabhängig von den C#- und Python-Versionsräumen.
Begründung:

- **Maven-Koordinaten erwarten SemVer-artige Versionen** (`groupId:artifactId:version`);
  ein `-SNAPSHOT`-Suffix ist die einzige ökosystem-spezifische Erweiterung,
  die für Vorab-Releases genutzt werden kann, aber für dieses Erst-Release
  nicht gebraucht wird.
- **Unterschiedliche Änderungstakte**, dieselbe Begründung wie bei C#/Python
  (`ADR-0106` Festlegung 3, `ADR-0107` Festlegung 4): ein Kotlin-SDK-Patch
  braucht keinen Server- oder Schwester-SDK-Release und umgekehrt.
- **Start bei `0.x.y`** (z. B. `0.1.0`), aus denselben Gründen wie bei den
  beiden bestehenden Packages — vorstabile öffentliche API, `1.0.0`
  markiert bewusst versprochene Rückwärtskompatibilität.
- **Ein fünfter, unabhängiger Versionsraum** neben Server, Harness-Baseline,
  C#-SDK und Python-SDK ist dieselbe fortgesetzte Disziplin, keine neue
  Klasse.
- **Von der Vertriebsweg-Wahl unberührt** — GitHub Packages verlangt kein
  eigenes Versionsschema, dieselbe SemVer-Wahl gälte unverändert auch bei
  Maven Central (§Verglichene Alternativen B).

### 5 — Build-/Publish-Mechanismus: Gradle, eingebautes `maven-publish`,
`GITHUB_TOKEN`, Docker-only, JDK 21 (nicht 25)

- **Build-System: Gradle**, wie `examples/kotlin/` es bereits vormacht —
  keine Parallel-Werkzeugkette (Maven/`pom.xml`) für dieselbe Sprache
  innerhalb dieses Repos (§Verglichene Alternativen F).
- **Gradle-Plugin: das eingebaute `maven-publish`, kein Dritt-Plugin.**
  Anders als beim ursprünglich erwogenen Maven-Central-Weg (der
  `com.vanniktech.maven-publish` gebraucht hätte, u. a. wegen der dort
  verlangten GPG-Signierung und Central-Portal-Validierung) braucht GitHub
  Packages **keine** Signierung und **keine** Portal-spezifische Logik —
  das von GitHub selbst dokumentierte Minimalrezept (§Kontext Recherche)
  ist ein einfacher `publishing { repositories { maven { … } } }`-Block mit
  dem Gradle-Kern-Plugin. Ein zusätzliches externes Plugin wäre hier
  unnötige Komplexität ohne Gegenwert (§Verglichene Alternativen F).
- **JDK-Basis-Image: `eclipse-temurin:21-jdk`, nicht `25-jdk`.** Real
  gemessen (§Kontext Recherche): JDK 25 ist die aktuelle LTS-Version, aber
  Gradle 8.14 — die von `examples/kotlin/` bereits erprobte Version —
  unterstützt JDK 25 nur ab Gradle 9.1.0 als Toolchain-Ziel. Ein
  gleichzeitiger Sprung auf eine neue JDK-**und**-Gradle-Kombination für
  ein Erst-Release ohne Referenz in diesem Repo würde genau die Vorsicht
  aufgeben, die die Vorarbeits-Parität mit C# eigentlich rechtfertigt (§1).
  Das SDK bleibt deshalb bei JDK 21/Gradle 8.14 — derselben, bereits real
  gebauten Kombination wie `examples/kotlin/` —, und der umsetzende Zug
  übernimmt den in `examples/kotlin/Dockerfile` bereits real gepinnten
  `eclipse-temurin:21-jdk`-Digest (`ADR-0087` §Entscheidung Festlegung 3-
  Muster: Kandidat + Digest in der ADR, Pin-Commit beim umsetzenden Zug —
  hier bereits vollständig vorhanden, kein neuer Digest-Bezug nötig).
- **Kein Gate.** GitHub-Packages-Paketbezug (Kotlin-Gradle-Plugin, gRPC-/
  Coroutine-Bibliotheken von Maven Central weiterhin als
  Abhängigkeits-Quelle, unabhängig vom Publish-Ziel) und der Publish-Schritt
  selbst brauchen Netz — ein neues, netzlos **nicht** prüfbares
  Werkzeug-Ziel, Arbeitsname `make sdk-pack-kotlin`, baut/testet/paketiert
  (`./gradlew build`, `./gradlew publish`) im gepinnten
  `eclipse-temurin:21-jdk`-Image.
- **Kein zusätzliches Repository-Secret.** `GITHUB_TOKEN` ist der von
  GitHub Actions **automatisch** bereitgestellte Workflow-Token; der
  Publish-Workflow setzt lediglich `permissions: contents: read` /
  `packages: write` (real dokumentiertes Minimalrezept, §Kontext
  Recherche) — kein Arbeitsname für ein neu anzulegendes Secret nötig, im
  Gegensatz zu `NUGET_API_KEY`/`PYPI_API_TOKEN`.
- **Der Trigger ist ein eigener Tag-Namensraum, weder `v*` noch
  `sdk-csharp-v*`/`sdk-python-v*`.** Vorschlag für den umsetzenden Zug:
  `sdk-kotlin-v*` als eigener Tag-Präfix, eigener Workflow
  (`sdk-kotlin-release.yml`), der den Tag strikt gegen SemVer 2.0 validiert
  und gegen die im Gradle-Build (`version = "…"` in
  `sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts`) geführte Version
  abgleicht — dasselbe Muster wie bei C#/Python, nur gegen die
  Gradle-Versionszeile statt `.csproj`/`pyproject.toml`.
- **Beides ist Folgepflicht dieser ADR, nicht Teil ihres Umfangs** — diese
  ADR entscheidet die *Form*, sie implementiert **nichts**
  (§Konsequenzen Folgepflicht).

### 6 — Was diese ADR nicht ändert

- **Kein Produktionscode, kein Eingriff in `internal/**`/`cmd/**`.**
- **`ADR-0106`/`ADR-0107`/`ADR-0108` und `sdks/csharp/**`/`sdks/python/**`
  bleiben unverändert** — kein Umzug, keine Zusammenlegung der drei
  Sprach-SDKs in dieser ADR.
- **`examples/kotlin/**` bleibt unverändert** — kein Umzug, keine
  Code-Extraktion in dieser ADR.
- **`spec/architecture.md` bleibt unberührt** — `sdks/**` ist ein
  Gate-Scope-Bereich außerhalb des Servers, keine `ARC-*`-Komponente.
- **`docs/user/version.md` bleibt unberührt** — die Server-Version bleibt
  eigenständig geführt.
- **Kein neues Gate, keine Schwellen-Senkung.**
- **`.a-check.yml` bleibt unberührt** — a-check liest kein Kotlin.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

### A — Dritte Sprache

| Option | Pro | Contra |
|---|---|---|
| A1 — Go (Go-Modul-Registry) | spiegelt die Server-Implementierungssprache | bereits in `ADR-0106`/`ADR-0107` §Verglichene Alternativen A1 zweimal verworfen — Go-Consumer haben die geringste Einstiegshürde über die rohe API, der SDK-Nutzen ist hier am kleinsten, gerade wo er am billigsten wäre; keine bestehende Nutzerentscheidung dafür |
| A2 — TypeScript/npm | breites Web-/Node-Ökosystem | keine bestehende Beispiel-Client-Vorarbeit (kein `examples/typescript/`) — dasselbe Risiko wie bei Python bei `ADR-0107`, aber ohne dessen kompensierende Zielgruppen-Passung; keine bestehende Nutzerentscheidung dafür |
| **A3 — Kotlin/GitHub Packages (gewählt)** | explizite Nutzerentscheidung; reale Vorarbeits-**Parität** mit C# (fünf funktionierende Beispiel-Clients, alle vier relevanten Zustellwege bereits gedeckt, heute erneut gemessen); erweitert die erreichte Zielgruppe über .NET/Python hinaus — Enterprise-Java-/Android-/JVM-Server-Backend-Entwickler, die C#- oder Python-Pakete gar nicht in Betracht ziehen; Kotlin-Artefakte sind Java-binärkompatibel, erreichen also **zusätzlich** reine Java-Konsumenten ohne eigenes Kotlin-Package | `ADR-0107` §Verglichene Alternativen A3 hatte Kotlin bereits einmal als Kandidaten für die **zweite** Sprache geprüft und mit dem Argument verworfen, dass Java/Maven-Consumer bereits etablierte Bordmittel für HTTP-/gRPC-Clients haben (generierte gRPC-Java-Stubs) — dieser Einwand bleibt inhaltlich richtig, trägt aber gegen die **heutige**, explizite Nutzerentscheidung nicht mehr: „etablierte Bordmittel" senken den SDK-Mehrwert graduell, sie eliminieren ihn nicht (derselbe generierte-Stub-Umweg gilt für jede Sprache ohne SDK, auch für Go bei A1) |
| A4 — nichts tun / abwarten | kein Aufwand jetzt | die explizite Nutzerentscheidung vom 2026-09-20 würde ignoriert; `LH-FA-SST-009.a`s „dritte Sprache offen" bliebe unbeantwortet, obwohl die Vorarbeits-Parität mit C# bereits vorliegt |

### B — Vertriebsweg für Kotlin

| Option | Pro | Contra |
|---|---|---|
| **B1 — GitHub Packages (gewählt)** | explizite, zweite Nutzerentscheidung; **kein** zusätzliches Repository-Secret (`GITHUB_TOKEN` genügt, real dokumentiert); kein GPG-Schlüsselpaar, kein Sonatype-Konto, kein Namespace-Antrag — der mit Abstand einfachste Publish-Weg aller vier bisher genutzten/erwogenen Vertriebswege dieses Repos; liegt bereits im selben Hosting wie Quellcode und CI | verlangt **immer** eine Authentifizierung zum **Lesen**, auch bei einem öffentlichen Package (real dokumentiert, §Kontext) — ein Consumer außerhalb dieses Repos braucht ein GitHub-Konto und einen klassischen PAT mit `read:packages`-Scope; kein anonymer `pip install`-/`dotnet add package`-artiger Bezug |
| B2 — Maven Central über den Central Portal | de-facto-Standardweg für JVM-/Kotlin-Bibliotheken; anonymer Lesezugriff ohne Konto/Token für jeden Consumer; `io.github.pt9912`-Namespace real ohne DNS-Eintrag verifizierbar (Recherche eines früheren Entwurfs dieser ADR) | **vom Nutzer explizit abgelehnt** („Anstatt maven können wir github verwenden"); braucht ein Central-Portal-Konto, eine Namespace-Registrierung und **fünf** Secrets (Portal-Token-Paar + GPG-Schlüssel-Tripel) statt keines — der deutlich aufwendigere Erst-Release-Weg |
| B3 — JCenter | historisch verbreitet in der Android-/JVM-Welt | **real abgeschaltet seit 2021** — keine lebende Option, muss aber genannt werden, um sie korrekt auszuschließen |
| B4 — private/interne Package-Registry | volle Zugriffskontrolle | widerspricht der AC „öffentlich über den Paketmanager konsumierbar" ohne fachlichen Grund — strukturell nicht besser als B1s Auth-Pflicht, aber ohne den Vorteil, bereits im bestehenden Repository-Hosting zu liegen |
| B5 — nur Source (kein Paketmanager-Vertriebsweg, z. B. JitPack aus einem Git-Tag) | kein Registrierungsprozess, kein Konto nötig | erfüllt die AC nicht vollwertig: JitPack baut aus einem öffentlichen Git-Tag zur Abrufzeit beim Consumer statt aus einem vorab geprüften Artefakt — geringere Lieferzusicherung als ein regulärer Registry-Release; kein etablierter Standardweg für ein Erst-Release, das Vertrauen aufbauen soll |

### C — Umfang des ersten Pakets

| Option | Pro | Contra |
|---|---|---|
| C1 — nur HTTP-API (wie Pythons v1) | kleinster Erst-Release-Umfang | ignoriert die reale Vorarbeits-Parität mit C# — es gibt kein Wire-Risiko, das eine derart kleine Fläche rechtfertigt (§Kontext „Was das ändert") |
| **C2 — HTTP-API + gRPC-Stream, ein Package (wie C#, gewählt)** | deckt Abruf **und** Live-Zustellung; nutzt die Vorarbeits-Parität voll aus, ohne den strukturellen Footprint-Einwand zu ignorieren; Konsistenz mit dem bereits zweimal bewährten Zuschnitt (`ADR-0106`) | zieht `io.grpc:*`/`kotlinx-coroutines-core` auch HTTP-only-Consumern auf |
| C3 — HTTP-API + gRPC-Stream + SSE | ein zusätzlicher Live-Transport sofort verfügbar, alle drei „synchron-relevanten" Wege in einem Release | SSE ist fachlich redundant zu gRPC für denselben Consumer (beide sind alternative Transporte derselben Live-Zustellungs-Fähigkeit, `LH-FA-SST-009` Out-of-Scope) — vergrößert die Erst-Release-Fläche ohne neuen fachlichen Nutzen, nur weil ein Vorbild existiert |
| C4 — alle vier Zustellwege in einem Package | „vollständige Matrix sofort", volle Ausnutzung der fünffachen Vorarbeit | derselbe Footprint-Einwand wie in `ADR-0106` §Verglichene Alternativen B3 — verschärft, weil hier zusätzlich `io.nats:jnats` und ein SSE-Client-Stack hinzukommen; größter Erst-Release-Umfang bei einem **ersten** Paket-Design, das noch niemand geprüft hat |
| C5 — vier/fünf separate Packages, eines je Zustellweg | jeder Consumer lädt nur, was er braucht | fünf Maven-Koordinaten, fünf Versionsräume, fünf Publish-Konfigurationen für ein **erstes** SDK — Overhead, bevor auch nur ein Package real veröffentlicht wurde; verfrühte Modularisierung ohne Nutzungsdaten (analog `ADR-0106`s B4-Verwerfung) |

### D — Verhältnis zu `examples/kotlin/`

| Option | Pro | Contra |
|---|---|---|
| D1 — ein bestehendes Beispiel (z. B. `http-client`) wird zum SDK weiterentwickelt | kein neuer Baum, kein doppeltes Schreiben ähnlichen Codes | verletzt `SPEC-023`s Klassentrennung „Vorbild, kein Belegträger" an genau der Stelle, an der ein Dritter sich künftig auf Stabilität verlässt; das `application`-Plugin-CLI-Gerüst der Beispiele passt nicht zu einer öffentlichen Bibliotheks-API |
| **D2 — eigenständiges Projekt unter `sdks/kotlin/` (gewählt)** | klare räumliche und rollenbezogene Trennung (Vorbild vs. Belegträger); eigene Import-Grenze, eigene Versionierung, eigener Bauweg, ohne die Beispiel-Verträge zu berühren; identisches Muster wie bei C#/Python | ähnlicher Code entsteht zweimal (Beispiel und SDK) — vertretbar, weil beide unterschiedliche Verträge tragen (CLI-Demo vs. öffentliche API) |
| D3 — `examples/kotlin/` wird komplett nach `sdks/kotlin/` verschoben, Beispiele verschwinden | ein Baum weniger zu pflegen | zerstört `SPEC-023`s Beispiel-Anspruch (Vorbild, Docker-only-Startform, Handbuch-Bindung, `ADR-0098`) ohne Not; Beispiele und SDK haben unterschiedliche Konsumenten (Repo-Leser vs. Package-Konsument) |

### E — Versionierung

| Option | Pro | Contra |
|---|---|---|
| E1 — gekoppelt an `docs/user/version.md` oder an ein anderes SDK | ein Blick genügt | erzwingt SDK-Releases bei jedem Server- bzw. Schwester-SDK-Release ohne eigenen Draht-Änderungsgrund; passt nicht zur Boundary-AC (SDK-eigene Versionsgrenze) |
| **E2 — eigenständiges SemVer 2.0, ab `0.x.y` (gewählt)** | Release-Takt folgt der tatsächlichen Kotlin-SDK-Änderung; SemVer ist das im Maven-/Gradle-Ökosystem native Schema, kein Kompatibilitäts-Umweg wie Pythons PEP 440; Analogie zu bereits vier bestehenden unabhängigen Versionsräumen (Server, Harness-Baseline, C#-SDK, Python-SDK) | ein fünfter Versionsraum, den ein Betrachter separat nachschlagen muss |
| E3 — Versionierung nach Datum (CalVer) | keine Debatte über Major/Minor/Patch | passt nicht zur AC-Boundary („z. B. SemVer-Major"); unüblich für Maven-/Gradle-Bibliotheken |

### F — Build-/Publish-Werkzeug

| Option | Pro | Contra |
|---|---|---|
| **F1 — Gradles eingebautes `maven-publish`-Plugin (gewählt)** | kein zusätzliches Dritt-Plugin; exakt das von GitHub selbst dokumentierte Minimalrezept für GitHub Packages (real geprüft, README/Anleitung); GitHub Packages braucht keine Signierung — der Hauptvorteil eines schwereren Publish-Plugins (native GPG-/Central-Portal-Integration) entfällt hier ersatzlos | mehr manuelle `build.gradle.kts`-Pflege als ein All-in-one-Werkzeug, falls später doch signiert werden soll |
| F2 — `com.vanniktech.maven-publish` | ein Plugin für Signieren **und** Publizieren, falls später zusätzlich Maven Central bedient werden soll; aktiver Entwicklungstakt (`0.37.0`, real gemessen) | sein Haupt-Mehrwert (native GPG-Signierung, Central-Portal-Validierung) ist für GitHub Packages **ungenutzt** — ein zusätzliches externes Plugin ohne Gegenwert für den heute gewählten Vertriebsweg |
| F3 — Maven (`pom.xml`) statt Gradle als Build-System | Maven ist ebenfalls ein etablierter JVM-Standard | führt eine **zweite** Build-Werkzeugkette für dieselbe Sprache in diesem Repo ein (`examples/kotlin/` nutzt bereits Gradle) — kein erkennbarer Vorteil, der den doppelten Pflegeaufwand rechtfertigt |

**Fazit:** A3, B1, C2, D2, E2, F1. A1/A2 scheitern an fehlendem Nutzen bzw.
fehlender Vorarbeit ohne kompensierende Zielgruppen-Passung, A4 an der
unbeantworteten Nutzerentscheidung; B2 wurde vom Nutzer explizit
zurückgewiesen und ist ohnehin der aufwendigere Weg, B3 ist real
abgeschaltet, B4/B5 passen nicht zur AC bzw. senken die Lieferzusicherung;
C1 ignoriert die Vorarbeits-Parität, C3/C4/C5 überziehen den Erst-Release
bzw. modularisieren verfrüht; D1/D3 verletzen `SPEC-023`s Klassentrennung
ohne Not; E1/E3 passen nicht zur Boundary-AC bzw. zum Ökosystem; F2 bringt
Werkzeug-Gewicht ohne Gegenwert für GitHub Packages, F3 führt eine
Parallel-Werkzeugkette ohne Grund ein.

## Konsequenzen

- **Positiv:** `LH-FA-SST-009.a`s dritte offene Frage bekommt eine
  konkrete, begründete Antwort inklusive ehrlicher Vertriebsweg-Abwägung
  (Maven Central geprüft und vom Nutzer bewusst abgelehnt, nicht
  übergangen); die gemessene Vorarbeits-**Parität** mit C# (nicht:
  Überlegenheit) erlaubt es, einen bereits einmal bewährten, konservativen
  Scope-Schnitt direkt zu übertragen; der Publish-Mechanismus ist der
  bislang einfachste aller drei SDK-Sprachen — kein externes Secret, kein
  externes Konto; die Zielgruppen-Erweiterung (Enterprise-Java/Android/
  JVM-Server-Backend, zusätzlich Java-binärkompatibel nutzbar) ergänzt die
  bereits bedienten .NET-/Python-Segmente.
- **Negativ:** SSE- und NATS-Consumer in Kotlin bleiben vorerst auf den
  direkten Zugriffsweg verwiesen (durch die Negative-AC gedeckt, wie bei
  C#); ein Kotlin-Consumer braucht — anders als bei C#/NuGet und
  Python/PyPI — zwingend ein GitHub-Konto und einen klassischen PAT mit
  `read:packages`-Scope, um das Package überhaupt zu **installieren** (real
  dokumentierter, bewusst in Kauf genommener Reibungspunkt, §Entscheidung
  Festlegung 2); das Package ist an dieses eine GitHub-Repository
  gebunden (`maven.pkg.github.com/pt9912/pg-change-feed`), nicht an eine
  vom Hosting unabhängige Registry-Identität.
- **Folgepflicht:**
  1. Ein Planner-Zug schneidet die Umsetzung (voraussichtlich mehrere
     Slices: SDK-Projektgerüst, HTTP-Client-Fläche, gRPC-Client-Fläche,
     Pack-/Publish-Werkzeug, Publish-Workflow) — diese ADR entscheidet die
     Form, nicht den Schnitt.
  2. `spec/pflichtenheft.md` §1 (`LH-FA-SST-009.a`) und §6 (Externe
     Verträge, analog `SPEC-026`/`SPEC-027`) brauchen einen Träger-Nachzug,
     sobald das Package real existiert (`AGENTS.md` §3.13) — Sache des
     umsetzenden Zuges.
  3. `harness/README.md` §Werkzeuge bekommt seine Zeilen für
     `make sdk-pack-kotlin` und den Publish-Workflow erst, wenn sie real
     existieren (`AGENTS.md` §4).
  4. `docs/user/benutzerhandbuch.md` bekommt einen SDK-Hinweis, im selben
     Zug wie das jeweilige Client-Programm — inklusive eines expliziten
     Hinweises auf die PAT-Pflicht zum Bezug (Festlegung 2), damit ein
     Consumer nicht überrascht wird.
  5. Der GitHub-Packages-Publish-Workflow gilt nach
     [`AGENTS.md`](../../../AGENTS.md) §3.10 erst nach einem realen,
     grünen Post-Push-Lauf als abgeschlossen.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `a-check` | keine — `sdks/kotlin/**` ist keine Go-Quelle, a-check liest sie nicht (`.a-check.yml` `languages: go`) | — |
| Review (kein Sensor) | `sdks/kotlin/**` importiert ausschließlich die Kotlin-/JVM-Standardbibliothek und öffentliche Maven-Koordinaten, keinen privaten Baum dieses Repos (Import-Grenze, analog `ADR-0106`/`ADR-0107`) | — |
| Review (kein Sensor) | `sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts` verwendet `eclipse-temurin:21-jdk` (nicht `25-jdk`) als Docker-Bau-Basis, solange Gradle in diesem Repo nicht auf ≥ 9.1.0 angehoben ist (§Entscheidung Festlegung 5) | — |
| Review (kein Sensor) | der Publish-Workflow verwendet ausschließlich `GITHUB_TOKEN` (kein neues Repository-Secret) für den Registry-Zugriff (§Entscheidung Festlegung 5) | — |
| künftiges Bau-Werkzeug (Folgepflicht) | `make sdk-pack-kotlin` baut/testet/paketiert Docker-only, netzlos geprüfte Teile fahren im Bau, kein Gate | `make sdk-pack-kotlin` (noch nicht existent) |

## Re-Evaluierungs-Trigger

1. **Eine vierte Sprache oder ein vierter Vertriebsweg wird für
   `LH-FA-SST-009` verlangt** (z. B. TypeScript/npm) — eigene ADR, diese
   ADR bleibt für Kotlin/GitHub Packages unberührt (dieselbe Regel wie
   `ADR-0106` §Re-Evaluierungs-Trigger 1 und `ADR-0107`
   §Re-Evaluierungs-Trigger 1 sie bereits zweimal formuliert haben).
2. **Das Kotlin-Package wurde real auf GitHub Packages veröffentlicht und
   Nutzungsdaten (Download-Zahlen, Issue-Nachfrage) legen eine andere
   Priorisierung nahe** — insbesondere ob SSE oder NATS-Vollinhalt als
   nächstes Folge-Release sinnvoller ist (Festlegung 1 offen gelassen).
3. **Gradle wird in diesem Repo real auf ≥ 9.1.0 angehoben** (z. B. weil
   `examples/kotlin/` selbst auf JDK 25 wechselt) — dann verliert die in
   §Entscheidung Festlegung 5 begründete Zurückhaltung bei der
   JDK-Basis-Image-Wahl ihre Grundlage, und ein Wechsel des SDK-Bau-Images
   auf `eclipse-temurin:25-jdk` ließe sich ohne die hier geltende Vorsicht
   rechtfertigen.
4. **Der reale GitHub-Packages-Publish-Workflow (`sdk-kotlin-release.yml`)
   scheitert im ersten realen Lauf** (`AGENTS.md` §3.10) an einer
   Berechtigungs- oder Registry-Eigenheit, die die Doku-Recherche nicht
   erfasst hat — Nachbesserung als eigene Folge-ADR mit
   `Supersedes ADR-0109`, sofern der Publish-Mechanismus selbst betroffen
   ist (ein reiner Workflow-Bugfix ohne Entscheidungsänderung braucht
   keine neue ADR).
5. **Die Authentifizierungspflicht beim Lesen erweist sich als echtes
   Consumer-Problem** — z. B. weil wiederholt Nachfrage nach einem
   anonym installierbaren Package aufkommt, oder weil ein Consumer ohne
   GitHub-Konto das Package strukturell nicht beziehen kann — dann Wechsel
   zu Maven Central (§Verglichene Alternativen B2, bereits einmal
   vollständig recherchiert: Namespace `io.github.pt9912`,
   `com.vanniktech.maven-publish` `0.37.0`, fünf benannte Secrets) als
   eigene Folge-ADR mit `Supersedes ADR-0109`.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-20 | Proposed | dieser Architect-Zug (kein vorausgehender Slice-Plan — ADR vor Implementierung, Modul 8); erster Entwurf schlug Maven Central vor |
| 2026-09-20 | Accepted | Nutzerentscheidung im Chat vom 2026-09-20 („Können wir noch ein kotlin-SDK erstellen", dann „Anstatt maven können wir github verwenden") + diese vollwertige ADR-Fassung samt Alternativenvergleich (Sprache **und** Vertriebsweg) und real recherchiertem Beleg (GitHub-Packages-Auth-Mechanismus, Maven-Central-Kontrastfolie, JDK-/Gradle-Kompatibilität) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0109` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
