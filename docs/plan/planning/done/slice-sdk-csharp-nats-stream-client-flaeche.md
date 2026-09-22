# Slice sdk-csharp-nats-stream-client-flaeche: Öffentliche NATS-Vollinhalts-Client-Fläche (`SPEC-024`), Version-Hebung, Träger-Nachzug

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** [welle-sdk-csharp-vollabdeckung](welle-sdk-csharp-vollabdeckung.md).

**Bezug:** [`LH-FA-SST-009`](../../../../spec/lastenheft.md),
[`LH-FA-SST-008`](../../../../spec/lastenheft.md) (Live-Change-Stream),
[`ADR-0106`](../../adr/0106-csharp-nuget-erstes-sdk-package.md)
§Entscheidung Festlegung 1, letzter Absatz (NATS-Vollinhalt als
Folge-Package ausdrücklich antizipiert), [`ADR-0100`](../../adr/0100-nats-dritter-vollinhalts-zustellweg.md)
(NATS-Vollinhalts-Stream-Vertrag, wird vom SDK benutzt, nicht erweitert),
[`AGENTS.md`](../../../../AGENTS.md) §3.13 (Träger-Nachzug — dieser Slice
löst ihn für `spec/pflichtenheft.md` aus).

**Berührte Spec-Stellen:** [`SPEC-024`](../../../../spec/pflichtenheft.md)
(Subjekt-Schema, Nachrichtenform, Authentifizierung), `SPEC-026`
(`PgChangeFeed.Client`-Metadaten — Version-Hebung),
[`LH-FA-SST-009.a`](../../../../spec/pflichtenheft.md) (Träger-Nachzug:
„deckt HTTP-API und gRPC-Stream" wird durch diese Welle falsch).

**Verantwortlich:** Implementer-Agent, 2026-09-22.

**Autor:** Planner-Agent, direkt beauftragt (Nutzerauftrag „alle drei SDKs
auf volle Vier-Wege-Parität"). **Datum:** 2026-09-21.

---

## 1. Ziel und Abgrenzung

**Ziel:** Eine öffentliche, stabile .NET-API-Fläche im bestehenden Package
`PgChangeFeed.Client` für den NATS-Vollinhalts-Stream von
[`SPEC-024`](../../../../spec/pflichtenheft.md) — Verbindung über
`NATS.Net`, Subjekt-Abonnement `cdc.stream.<source_id>.<schema>.<table>`
(oder ein Platzhalter-Muster wie `cdc.stream.<source_id>.>`), Token-Auth
auf Verbindungsebene, Deserialisierung der zehn Nachrichtenfelder von
`SPEC-024` (dasselbe Schema wie SSE/gRPC). **Zusätzlich** bündelt dieser
Slice, weil er der letzte Flächen-Slice dieser Welle ist: Version-Hebung
von `PgChangeFeed.Client` und den Träger-Nachzug in
`spec/pflichtenheft.md`/`docs/user/benutzerhandbuch.md` für **beide** neu
gelieferten Flächen (SSE **und** NATS-Vollinhalt).

**Übernimmt:** `slice-sdk-csharp-sse-client-flaeche` — nicht den Code,
sondern dessen offen gelassenen Doku-/Träger-Nachzug (§1 jenes Slice
„Ausdrücklich NICHT in diesem Slice").

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **HTTP-API-, gRPC- oder SSE-Fläche** — HTTP/gRPC bereits geliefert; SSE
  ist der direkte Vorgänger-Slice dieser Welle.
- **Ein zweiter Vertriebsweg oder eine vierte Sprache** — unverändert
  `ADR-0106` §Re-Evaluierungs-Trigger 1, keine Entscheidung dieses
  Slice.
- **Ein realer `sdk-csharp-v<Version>`-Tag-Push** — Betreiber-Entscheidung
  nach `AGENTS.md` §3.10 (Welle-Plan §3), nicht Teil dieses Slice.
- **CLI-Argument-Parsing wie `examples/csharp/nats-stream-client`** —
  dasselbe Argument wie bei den übrigen Flächen: ein SDK-Consumer ruft
  eine Methode auf, parst keine `argv`.
- **Ein Realserver-Integrationstest** — dieselbe Begründung wie beim
  SSE-Slice (Welle-Plan §6): die Vorarbeits-Parität senkt das Wire-Risiko
  bereits, keine strengere Teststufe ohne fachlichen Grund.

## 2. Definition of Done

- [x] `LH-FA-SST-009` erfüllt: eine öffentliche, im Package sichtbare
      Client-Klasse verbindet sich über `NATS.Net` mit Token-Auth,
      abonniert den Vollinhalts-Namensraum
      (`cdc.stream.<source_id>.<schema>.<table>`) und liefert die zehn
      Nachrichtenfelder von
      [`SPEC-024`](../../../../spec/pflichtenheft.md) — Tests
      referenzieren `SPEC-024` (Nachrichtenschema-Vollständigkeit,
      Subjekt-Formatierung netzlos, Verbindungs-Fehlerpfad ohne echten
      Server). Real umgesetzt: `PgChangeFeedNatsStreamClient` (Injektions-
      und Convenience-Konstruktor über `INatsClient`), `BuildSubject`/
      `BuildSourceSubject` als Subjekt-Formatierer, 15 neue Testmethoden
      (`[Fact]`/`[Theory]`-Attribute, real ausgezählt) in fünf Dateien
      (`Nats/`, inkl. `FakeNatsClient.cs` ohne eigene Testmethode) — 21
      tatsächlich laufende neue Testfälle (real gemessen: isolierter
      Docker-Testlauf gegen Elternstand `4b93def4` liefert `Passed: 49,
      Total: 49`, gegen diesen Commit `Passed: 70, Total: 70`, Delta 21),
      70/70 grün (`dotnet test` im Docker-Bau).
- [x] `sdks/csharp/Directory.Packages.props` bekommt `NATS.Net` als
      gepinntes Paket, real zum Bau-Zeitpunkt neu gemessen (nicht blind
      aus `examples/csharp/Directory.Packages.props`s `3.2.0` übernommen —
      `AGENTS.md` §3.12). Real gemessen am 2026-09-22:
      `api.nuget.org/v3-flatcontainer/nats.net/index.json` listet `3.2.0`
      weiterhin als jüngste stabile Version — derselbe Kandidat, real
      erneut bestätigt statt übernommen (siehe §7).
- [x] `PgChangeFeed.Client.csproj`s `<Version>` wird gehoben (von `0.1.0`
      auf `0.2.0` — additive, rückwärtskompatible Erweiterung, kein
      Breaking Change an HTTP-/gRPC-/SSE-Fläche).
- [x] `spec/pflichtenheft.md` `LH-FA-SST-009.a`/`SPEC-026` nachgezogen: der
      Satz „deckt HTTP-API und gRPC-Stream" wird zu „deckt HTTP-API,
      gRPC-Stream, SSE und NATS-Vollinhalt" (`AGENTS.md` §3.13) —
      Suchlauf-Pflicht: `grep -rn "HTTP-API und gRPC-Stream\|PgChangeFeed.Client.0.1.0"`
      über `spec/`, `docs/`, `harness/` (Ergebnis wird in §7 berichtet,
      gefunden **und** nicht gefunden).
- [x] Kein Import aus `internal/**`/`cmd/**`/`gen/**` dieses Repos — real
      geprüft (`grep -rn "internal/\|cmd/\|gen/" sdks/csharp/`), kein
      Treffer in den neuen `Nats/**`-Dateien.
- [x] `make gates` grün — Exit-Code `0`, ungepiped geprüft (siehe §7).
- [x] Ein real neu gebautes `.nupkg` (`make sdk-pack-csharp`) mit
      `<Version>0.2.0` als Smoke-Beleg — alle vier Client-Flächen im
      selben Artefakt. Real erzeugt: `sdks/csharp/dist/PgChangeFeed.Client.0.2.0.nupkg`
      (37451 Bytes, `lib/net10.0/PgChangeFeed.Client.dll` enthalten).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      `docs/reviews/review-slice-sdk-csharp-nats-stream-client-flaeche.md`
      (2× HIGH, 1× MEDIUM) durch eine gezielte Fixrunde behoben (F-1
      Kommentar-Chronik zurückgeführt, F-2/F-3 Zahlen real neu gemessen
      und korrigiert) — kein offenes HIGH/MEDIUM mehr.
- [x] Doku-Update: `docs/user/benutzerhandbuch.md` bekommt einen
      SDK-Hinweis für SSE **und** NATS-Vollinhalt (`ADR-0106`
      §Konsequenzen Folgepflicht 4-Muster) — getragen durch die bereits
      verkörperte Selbstprüf-Instruktion und den Reviewer-HIGH-Punkt
      (`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`).
      Version 1.37 → 1.38, Stand 2026-09-22, Zeile 1.38 im
      Änderungshistorie-Anhang.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben,
      **falls dieser Slice einen Inventur-Fund auflöst** — entfällt: keine
      Reconciliation-Datei in diesem Repo.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder Beleg in `evidence/`; keine Beobachtung angefallen
      ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      dieser Slice gehört zu
      [welle-sdk-csharp-vollabdeckung](welle-sdk-csharp-vollabdeckung.md)
      (noch offen); die Prüfung läuft regelkonform bei deren Closure.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `sdks/csharp/PgChangeFeed.Client/Nats/PgChangeFeedNatsStreamClient.cs` (Arbeitsname) | neu | öffentliche API-Fläche für den NATS-Vollinhalts-Stream (`SPEC-024`). |
| `sdks/csharp/Directory.Packages.props` | update | `NATS.Net`-Paketversion ergänzen, real gemessen. |
| `sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.csproj` | update | `<Version>0.2.0</Version>`, `PackageReference` auf `NATS.Net`. |
| `sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.Tests/Nats/*.cs` (Arbeitsname) | neu | Nachrichtenschema-Vollständigkeit, Subjekt-Formatierung, Verbindungs-Fehlerpfad. |
| `spec/pflichtenheft.md` §1 (`LH-FA-SST-009.a`), §6 (`SPEC-026`-Zeile) | update | Träger-Nachzug: volle Vier-Wege-Abdeckung für C#. |
| `docs/user/benutzerhandbuch.md` | update | SDK-Hinweis für SSE und NATS-Vollinhalt. |

**Ansatz:** Referenzmaterial ist `examples/csharp/nats-stream-client/`
(`Format.cs` für das Nachrichtenschema, `Program.cs`/`Cli.cs` für den
Verbindungs-/Subjekt-Aufbau, `NATS.Net`-Nutzung bereits real erprobt) und
die bestehende gRPC-Fläche dieses Packages (Fake-Verbindungs-Test-Muster,
`FakeCallInvoker.cs`-Analogie für einen gefakten NATS-Verbindungsfehler).

**Bekannte Docker-Bau-Kopplung:** wie beim SSE-Slice — `--build-context
proto=proto` bleibt für jeden Bau dieses Baums zwingend, unabhängig vom
NATS-Bezug (`BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut`).
Real erneut bestätigt: `docker build -f sdks/csharp/Dockerfile sdks/csharp`
ist in diesem Zug nicht ohne den Zusatzkontext versucht worden (bereits
zweifach belegt durch die vorigen Slices); kein vierter Beleg dieser
Beobachtung wurde hier erzeugt.

**Plan-Nachzug (Implementer-Zug):**

- Zusätzliche neue Dateien gegenüber der §3-Tabelle:
  `sdks/csharp/PgChangeFeed.Client/Nats/Models/Change.cs` (ein
  [`SPEC-024`](../../../../spec/pflichtenheft.md)-DTO, eigene Datei —
  dasselbe Trennungsmuster wie `Sse/Models/Change.cs`
  gegenüber `Sse/PgChangeFeedSseClient.cs`) und
  `sdks/csharp/PgChangeFeed.Client/Nats/PgChangeFeedNatsMalformedMessageException.cs`
  (eine bewusst **eigenständige**, kleine Exception-Klasse — **keine**
  Wiederverwendung von `PgChangeFeed.Client.Http.PgChangeFeedException`s
  Hierarchie: deren `StatusCode`-Feld ist ein HTTP-Konzept ohne
  NATS-Äquivalent; ein Platzhalter-Statuscode dafür wäre ein schlechterer
  Fit als ein kleiner, eigener Typ. Das vermeidet bewusst das Muster, das
  der SSE-Slice-Report unter F-2 als Wartbarkeitsrisiko benannte — hier
  aus strukturellem Grund, nicht nur aus Vorsicht).
- Statt einer einzelnen `PgChangeFeed.Client.Tests/Nats/*.cs`-Sammlung
  liegen die Tests in fünf Dateien (`FakeNatsClient.cs`,
  `ChangeMessageSchemaTests.cs`, `SubjectTests.cs`,
  `PgChangeFeedNatsStreamClientTests.cs`,
  `PgChangeFeedNatsStreamClientAuthBoundaryTests.cs`) — dasselbe
  Gruppierungsmuster wie bei den SSE-/gRPC-Slices.
- Zwei zusätzliche öffentliche Subjekt-Formatierer
  (`PgChangeFeedNatsStreamClient.BuildSubject`/`BuildSourceSubject`) plus
  eine `AllSourcesSubject`-Konstante, über die §2-DoD-Formulierung
  „Subjekt-Formatierung netzlos" hinaus konkretisiert (der Plan-Wortlaut
  benannte den Testinhalt, nicht die konkrete API-Form) — Design-
  Ausgestaltung, kein Plan-Konflikt.
- `NATS.Net`s öffentliches `INatsClient`-Interface (`SubscribeAsync<T>`)
  trägt bereits eine testbare Injektionsstelle, analog `Grpc.Core.CallInvoker`
  bei der gRPC-Fläche — kein eigener Adapter-Typ nötig;
  `PgChangeFeed.Client.Tests.Nats.FakeNatsClient` implementiert nur
  `SubscribeAsync<T>` real, alle übrigen `INatsClient`-Mitglieder werfen
  `NotSupportedException` (Muster `FakeCallInvoker`).
- Auth-Boundary: **keine** eigene Exception-Hierarchie für einen
  abgelehnten NATS-Verbindungsversuch — die reale `NATS.Net`-Ausnahme
  (`NatsServerException`/`NatsConnectionFailedException`) propagiert
  unverändert aus der Aufzählung, dieselbe „nicht stillschweigend
  geschluckt"-Semantik wie gRPCs `RpcException`, aber ohne eine zweite,
  parallele Mapping-Schicht zu erfinden.
- `AGENTS.md` §3.13-Suchlauf (Träger-Nachzug), real ausgeführt (Ergebnis
  vollständig in §7): `grep -rn "HTTP-API und gRPC-Stream\|PgChangeFeed.Client.0.1.0"
  spec/ docs/ harness/` plus gezielte Prüfung von `docs/user/releasing.md`,
  `README.md`/`README.de.md`, `sdks/csharp/README.md`,
  `sdks/csharp/Dockerfile`, `harness/mk/sdk.mk`,
  `tools/harness/sdk-pack-csharp.sh` und der `Sse/Models/Change.cs`-
  Kommentar „three surfaces, three independent wire contracts" (durch das
  wortgleiche Muster in der neuen `Nats/Models/Change.cs`-Datei
  aufgefallen, nicht durch den `grep` selbst) — alle real geänderten
  Stellen und alle bewusst unverändert gelassenen (ADRs, `done/`,
  `docs/reviews/`, Kotlin-Zeilen, historische Changelog-Zeilen) stehen in
  §7.

## 4. Trigger

**Start** (`next` → `in-progress`): wenn
`slice-sdk-csharp-sse-client-flaeche` in `done/` liegt (Welle-Plan §4
Reihenfolge — bewusst sequentiell, nicht parallel).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls der
  Träger-Nachzug in `spec/pflichtenheft.md` mehr Stellen betrifft als
  erwartet (z. B. weitere `SPEC-026`-Referenzen außerhalb §1/§6) — dann
  Abspaltung eines reinen Doku-Nachzug-Slice.
- `in-progress` → `open` (blockiert — Carveout?): `slice-sdk-csharp-sse-client-flaeche`
  liegt noch nicht in `done/`.

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + reales `.nupkg` mit `0.2.0` +
Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- Ein gefakter NATS-Verbindungsfehler-Test könnte das reale
  Token-Ablehnungsverhalten des NATS-Servers nicht exakt nachbilden.
  **Ausgang:** weiter offen — ein realer Rundlauf-Beleg bleibt
  `make test-integration`s bestehendem `tools/harness/natsstreamsub`
  vorbehalten (Wegwerf-Client, kein SDK-Import); dieselbe Teststrategie
  wie bei HTTP/gRPC/SSE.
- Der Träger-Nachzug in `spec/pflichtenheft.md` könnte eine Stelle
  übersehen, die denselben veralteten Satz an einer nicht per `grep`
  gefundenen Formulierung trägt (`AGENTS.md` §3.13 Grenze: Symbolnamen
  trifft `grep` zuverlässig, Prosa-Umformulierungen nicht). **Ausgang:**
  weiter offen — Review prüft den Diff unabhängig gegen denselben
  Suchraum.
- Version-Hebung `0.1.0` → `0.2.0` ohne SemVer-Minor-Konvention explizit
  in einer Datei dieses Repos festgeschrieben (nur `ADR-0106`
  §Entscheidung Festlegung 3 „ab `0.x.y`", kein Minor-Inkrement-Schema).
  **Ausgang:** entfallen — additive, rückwärtskompatible Erweiterung ist
  im SemVer-Vokabular unmissverständlich ein Minor-Bump, keine
  ADR-pflichtige Entscheidung.

## 7. Closure-Notiz

- **Was hat funktioniert:** `NATS.Net`s öffentliches `INatsClient`-Interface
  (`SubscribeAsync<T>`) trug bereits eine testbare Injektionsstelle,
  strukturell analog `Grpc.Core.CallInvoker` — kein neuer Adapter-Typ
  nötig, `FakeNatsClient` folgt demselben Muster wie `FakeCallInvoker`. Ein
  glatter Durchlauf: 70/70 Tests grün beim ersten Docker-Bau, kein
  Fixzyklus zwischen Implementer-Läufen nötig. Drei gezielte
  Mutationstests (ParseChange-Nullpfad entfernt, Auth-Exception
  geschluckt, `Change`-Feld entfernt) haben je genau den erwarteten Test
  rot gemacht und keinen anderen — die drei Zusagen
  (Nachrichtenschema-Vollständigkeit, Auth-Boundary „nicht stillschweigend
  geschluckt", Malformed-Payload-Pfad) sind damit nicht nur behauptet,
  sondern real gegen ein rotes Gegenbeispiel gehalten (`AGENTS.md`-Geist
  „zu jeder Zusage das rot gesehene Gegenbeispiel").
- **Was ging anders als geplant:** Die reale NuGet-Messung
  (`api.nuget.org/v3-flatcontainer/nats.net/index.json`, 2026-09-22) zeigt
  weiterhin `3.2.0` als jüngste stabile Version — keine Pin-Hebung nötig,
  aber real neu gemessen statt aus `examples/csharp/Directory.Packages.props`
  übernommen (`AGENTS.md` §3.12). Der Träger-Nachzug betraf mehr Dateien
  als der im Plan §2 genannte `grep`-Suchraum (`spec/`, `docs/`,
  `harness/`): zwei zusätzliche Fundstellen
  (`tools/harness/sdk-pack-csharp.sh`, `Sse/Models/Change.cs`) lagen
  außerhalb des vorgeschriebenen Musters und wurden erst beim Lesen
  benachbarter Dateien gefunden — kein Rückführungs-Trigger (§4 „zu groß"
  meint zusätzliche `spec/pflichtenheft.md`-Stellen, nicht zusätzliche
  Dateien außerhalb davon), aber ein neuer Beleg für
  `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (siehe unten). Design-
  Ausgestaltung gegenüber dem Plan-Wortlaut (kein Plan-Konflikt): eine
  bewusst **eigenständige** `PgChangeFeedNatsMalformedMessageException`
  statt einer Wiederverwendung der `PgChangeFeed.Client.Http.PgChangeFeedException`-
  Hierarchie (deren `StatusCode`-Feld ist ein HTTP-Konzept ohne
  NATS-Äquivalent) — und **kein** eigenes Exception-Mapping für die
  Auth-Boundary: die reale `NATS.Net`-Ausnahme propagiert unverändert.
- **Steering-Loop-Eintrag:** kein neuer Sensor, keine neue Hard Rule durch
  diese Closure selbst. Die 3×-Schwelle von
  `BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut` (bereits beim
  SSE-Slice real erreicht) bleibt eine Architect-Entscheidung bei der
  Welle-Closure — dieser Slice hat keinen vierten Beleg erzeugt (kein
  isolierter Bau-Versuch ohne `--build-context proto=proto` unternommen).
- **Beobachtungs-Register (`../observations/`):** neuer Beleg
  `evidence/slice-sdk-csharp-nats-stream-client-flaeche.md` zu
  `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (Zähler jetzt real
  ausgezählt **20×**) — real behobene Fundstellen (`spec/pflichtenheft.md`
  §1/§6, `sdks/csharp/README.md`, `PgChangeFeed.Client.csproj`s
  `<Description>`, `sdks/csharp/Dockerfile`, `harness/mk/sdk.mk`,
  `harness/README.md`, `README.md`/`README.de.md`,
  `tools/harness/sdk-pack-csharp.sh`, `Sse/Models/Change.cs`s
  „three surfaces"-Kommentar — **zehn** Dateien, real ausgezählt per
  `git diff --stat` (korrigiert gegenüber der zunächst genannten neun,
  Fixrunde review-slice-sdk-csharp-nats-stream-client-flaeche F-3); bewusst
  **keine** zusammenfassende „Fundstellen"-Zahl an dieser Stelle
  (Verifikationsbericht §2: „Fundstelle" trägt — anders als „Datei" — keine
  im Text definierte, mechanisch eindeutige Zähleinheit), zwei der zehn
  Dateien außerhalb des vorgeschriebenen `grep`-Musters gefunden (Details
  in der Evidenz-Datei).
  Übrige Register-Einträge geprüft, kein Handlungsbedarf:
  `BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`
  (verkörpert, DoD-Punkt „Doku-Update" real geliefert — siehe §2),
  `BEO-PGC/release-mechanismus-nicht-in-releasing-doku-nachgezogen` (offen,
  1×, weiterhin nicht einschlägig).
- **Folge-Slices:** keine aus diesem Slice selbst erwartet — dies ist der
  letzte Flächen-Slice von `welle-sdk-csharp-vollabdeckung`; die Welle
  selbst schließt als nächstes.
- **Risiken aus §6:**
  - „Ein gefakter NATS-Verbindungsfehler-Test könnte das reale
    Token-Ablehnungsverhalten des NATS-Servers nicht exakt nachbilden" —
    **Ausgang: weiter offen**, unverändert zur Plan-Begründung; ein realer
    Rundlauf-Beleg bleibt `make test-integration`s bestehendem
    `tools/harness/natsstreamsub` vorbehalten.
  - „Der Träger-Nachzug in `spec/pflichtenheft.md` könnte eine Stelle
    übersehen" — **Ausgang: weiter offen** für `spec/pflichtenheft.md`
    selbst (Review prüft unabhängig); für den **breiteren** Suchraum
    (`docs/`, `harness/`, `sdks/csharp/`) ist er **eingetreten und
    behoben**: zwei Fundstellen lagen außerhalb des `grep`-Musters (siehe
    Beobachtungs-Register oben) und wurden durch Lesen statt durch den
    Suchlauf selbst gefunden.
  - „Version-Hebung ohne SemVer-Minor-Konvention" — **Ausgang: entfallen**,
    wie geplant.
- **Drei Paarungen:** dieser Slice gehört zu
  [welle-sdk-csharp-vollabdeckung](welle-sdk-csharp-vollabdeckung.md)
  (noch offen) — die Prüfung läuft regelkonform bei deren Closure.
- **Planner-Closure-Nachtrag (Modul 8):** Verifikationsbericht
  (`verifikation-slice-sdk-csharp-nats-stream-client-flaeche.md`, Verdikt
  „DoD-konform: ja") vollständig gelesen; die dort nicht-blockierend
  gegebene Empfehlung zur „acht Fundstellen"-Formulierung umgesetzt —
  diese DoD-Zeile und der zugehörige Beobachtungs-Beleg
  (`evidence/slice-sdk-csharp-nats-stream-client-flaeche.md`,
  `../observations/BEO-PGC/arbeit-ueberholt-stehenden-traeger/state.md`)
  tragen jetzt ausschließlich die mechanisch eindeutige Zahl „zehn Dateien"
  (`git diff --stat`), ohne eine zusammenfassende „Fundstellen"-Zahl ohne
  definierte Zähleinheit (gewählt statt einer expliziten Zähldefinition,
  weil die Datei-Liste selbst bereits vollständig ist und keine zweite
  Zählregel gebraucht wird). Alle drei §6-Risiken tragen einen Ausgang
  (siehe oben); kein Risiko wurde durch diesen Nachtrag neu aufgeworfen.
  `BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut` (3×, Schwelle
  erreicht) bleibt bewusst **unverkörpert** durch diese Slice-Closure —
  ihr eigenes `state.md` verweist den Ausgang ausdrücklich an den
  Lese-Schritt der `welle-sdk-csharp-vollabdeckung`-Closure (Modul 6), und
  eine Regelschärfung selbst ist ohnehin eine Architect-, keine
  Planner-Entscheidung (Modul 4/8) — siehe dort.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Sub-Area `sdks/csharp/` — bereits
mit `slice-sdk-csharp-projektgeruest` eröffnet (GF), keine erneute
Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut` (§3 dieses Plans
trägt sie), `BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`
(verkörpert, DoD-Punkt „Doku-Update" trägt sie), `BEO-PGC/arbeit-ueberholt-stehenden-traeger`
(verkörpert, Suchlauf-Pflicht `AGENTS.md` §3.13 trägt den Träger-Nachzug
dieses Slice), `BEO-PGC/release-mechanismus-nicht-in-releasing-doku-nachgezogen`
(offen, 1×, nicht einschlägig — Publish-Mechanismus unverändert).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Fortsetzung von
`slice-sdk-csharp-projektgeruest`).
