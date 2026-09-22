# Slice sdk-csharp-sse-client-flaeche: Öffentliche SSE-Stream-Client-Fläche (`SPEC-021`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** [welle-sdk-csharp-vollabdeckung](welle-sdk-csharp-vollabdeckung.md).

**Bezug:** [`LH-FA-SST-009`](../../../../spec/lastenheft.md),
[`LH-FA-SST-008`](../../../../spec/lastenheft.md) (Live-Change-Stream —
Boundary: keine Zustellgarantie, kein Stream-internes Replay),
[`ADR-0106`](../../adr/0106-csharp-nuget-erstes-sdk-package.md)
§Entscheidung Festlegung 1, letzter Absatz (SSE als Folge-Package
ausdrücklich antizipiert — kein neuer ADR-Bedarf), Festlegung 2
(räumliche Trennung von `examples/csharp/`, hier unverändert gültig).

**Berührte Spec-Stellen:** [`SPEC-021`](../../../../spec/pflichtenheft.md)
(Endpunkt, Event-Form, Nachrichtenschema — das SDK benutzt diese
Festlegungen, verändert sie nicht).

**Verantwortlich:** Implementer-Agent, 2026-09-21.

**Autor:** Planner-Agent, direkt beauftragt (Nutzerauftrag „alle drei SDKs
auf volle Vier-Wege-Parität"). **Datum:** 2026-09-21.

---

## 1. Ziel und Abgrenzung

**Ziel:** Eine öffentliche, stabile .NET-API-Fläche im bestehenden Package
`PgChangeFeed.Client` für den SSE-Endpunkt von
[`SPEC-021`](../../../../spec/pflichtenheft.md) — `GET /changes/stream`,
Bearer-Token-Auth (`reader` oder `admin`), ein Frame-Parser, der
`event:`/`data:`-Zeilen zu strukturierten Events zusammensetzt und die
zehn Nachrichtenfelder von `SPEC-021` an den Consumer liefert (z. B. als
`IAsyncEnumerable<Change>`, dieselbe idiomatische Form wie die bestehende
gRPC-Fläche).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **HTTP-API- oder gRPC-Fläche** — beide bereits geliefert
  (`slice-sdk-csharp-http-client-flaeche`,
  `slice-sdk-csharp-grpc-client-flaeche`); dieser Slice fügt nur SSE
  hinzu.
- **NATS-Vollinhalts-Fläche** — eigener Slice
  (`slice-sdk-csharp-nats-stream-client-flaeche`), getrennter Draht-Vertrag
  (NATS statt HTTP), unabhängige Fremdabhängigkeit (`NATS.Net`).
- **Stream-internes Replay oder `Last-Event-ID`-Auswertung** — `SPEC-021`
  trägt beides nicht („der `Last-Event-ID`-Header wird weder gesendet noch
  ausgewertet"); das SDK kann keine Fähigkeit anbieten, die der Draht nicht
  hat.
- **Version-Bump, `spec/pflichtenheft.md`-Träger-Nachzug,
  `docs/user/benutzerhandbuch.md`-Nachzug** — gebündelt im
  Folge-Slice `slice-sdk-csharp-nats-stream-client-flaeche` (Welle-Plan §4
  Reihenfolge), damit beide neuen Flächen zusammen, nicht zweimal
  nacheinander, in die Träger-Dokumente einfließen.
- **CLI-Argument-Parsing wie `examples/csharp/sse-client`** — dasselbe
  Argument wie bei HTTP/gRPC: ein SDK-Consumer ruft eine Methode auf,
  parst keine `argv`.
- **Ein Realserver-Integrationstest** — anders als bei der
  Python-Vollabdeckungs-Welle (`ADR-0110` Folgepflicht 1) verlangt weder
  `ADR-0106` noch diese Welle das für C# (Welle-Plan §6, Begründung: die
  bestehende Vorarbeits-Parität senkt das Wire-Risiko bereits, dieselbe
  Teststrategie wie bei HTTP/gRPC gilt fort).

## 2. Definition of Done

- [x] `LH-FA-SST-009` erfüllt: eine öffentliche, im Package sichtbare
      Client-Klasse öffnet `GET /changes/stream` (Bearer-Token in
      `Authorization`-Header), zerlegt SSE-Frames zu Events und liefert die
      zehn Nachrichtenfelder von
      [`SPEC-021`](../../../../spec/pflichtenheft.md) — Tests referenzieren
      `SPEC-021` (Frame-Parser netzlos, Authn-Boundary gegen eine gestubbte
      HTTP-Response ohne echten Server, Muster
      `PgChangeFeedHttpClientAuthBoundaryTests.cs`).
- [x] Kein Import aus `internal/**`/`cmd/**`/`gen/**` dieses Repos — real
      geprüft (`grep -rn "internal/\|cmd/\|gen/" sdks/csharp/`), Ausnahme
      nur Doku-Zitate von Test-Vorbildern.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      Bericht: `docs/reviews/review-slice-sdk-csharp-sse-client-flaeche.md`
      (0 HIGH, 0 MEDIUM, 3 LOW, 1 INFO — keine Fixrunde nötig,
      DoD-Checkbox-Nachzug ohne Fixrunde nach Skill-Regel §DoD-Checkbox-
      Nachzug ohne Fixrunde).
- [x] Doku-Update (`docs/user/benutzerhandbuch.md`,
      `spec/pflichtenheft.md`): bewusst **nicht** in diesem Slice —
      gebündelt im Folge-Slice `slice-sdk-csharp-nats-stream-client-flaeche`
      (§1 Abgrenzung), damit beide neuen Flächen zusammen in die
      Träger-Dokumente einfließen. Beide Dateien real unberührt
      (`git diff --stat` gegen diesen Diff), `PgChangeFeed.Client.csproj`s
      `<Version>` unverändert `0.1.0`.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben,
      **falls dieser Slice einen Inventur-Fund auflöst** — entfällt: keine
      Reconciliation-Datei in diesem Repo.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neuer
      Beleg in `BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut/evidence/`
      (Zähler jetzt 3×, Schwelle erreicht, Ausgang an die Closure von
      `welle-sdk-csharp-vollabdeckung` verwiesen — siehe §7).
- [x] Jedes Risiko aus §6 trägt einen Ausgang.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      dieser Slice gehört zu
      [welle-sdk-csharp-vollabdeckung](welle-sdk-csharp-vollabdeckung.md)
      (noch offen); die Prüfung läuft regelkonform bei deren Closure.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `sdks/csharp/PgChangeFeed.Client/Sse/PgChangeFeedSseClient.cs` (Arbeitsname) | neu | öffentliche API-Fläche für den SSE-Stream (`SPEC-021`). |
| `sdks/csharp/PgChangeFeed.Client/Sse/SseFrameParser.cs` (Arbeitsname) | neu | Frame-Parser (`event:`/`data:` → `Event`) — netzlos testbar, reines Zeilen-Verarbeiten, kein Fremdmodul (Muster `examples/csharp/sse-client/SseStream.cs`: „**Keine neue Abhängigkeit**"). |
| `sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.Tests/Sse/*.cs` (Arbeitsname, mehrere Dateien nach Testart) | neu | Frame-Parser-Grenzfälle (unvollständiges Frame, Leerzeile-Trennung), Authn-Boundary, Nachrichtenschema-Vollständigkeit gegen `SPEC-021`. |

**Ansatz:** Referenzmaterial ist `examples/csharp/sse-client/SseStream.cs`
(Frame-Zerlegung, bereits real erprobt, keine neue Fremdabhängigkeit) und
die bestehende HTTP-Fläche dieses Packages
(`PgChangeFeed.Client/Http/PgChangeFeedHttpClient.cs`,
`PgChangeFeedException`-Hierarchie — derselbe Auth-Header, dieselbe
Fehlerklasse für einen konsistenten Consumer-seitigen Fehlerpfad über SSE
und HTTP).

**Bekannte Docker-Bau-Kopplung (`BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut`,
2×, unter der Schwelle):** `sdks/csharp/Dockerfile` trägt eine einzige
`build`-Stufe für alle Flächen — ein `docker build` dieses Slice (auch ohne
gRPC-Bezug) braucht deshalb weiterhin `--build-context proto=proto`, sonst
bricht der Bau an der bereits bestehenden `COPY --from=proto`-Zeile ab.
`make sdk-pack-csharp` trägt den Flag bereits unconditional
(`tools/harness/sdk-pack-csharp.sh`) — kein Änderungsbedarf an Skript oder
Makefile-Fragment für diesen Slice, nur eine bewusst benannte, keine
stille Kopplung.

**Plan-Nachzug (Implementer-Zug):**

- Zusätzliche neue Datei gegenüber der §3-Tabelle:
  `sdks/csharp/PgChangeFeed.Client/Sse/Models/Change.cs` — das
  `SPEC-021`-Nachrichten-DTO (`JsonPropertyName`-Attribute) bekam eine eigene
  Datei statt im Client selbst zu liegen, dasselbe Trennungsmuster wie
  `Http/Models/Changes.cs` gegenüber `Http/PgChangeFeedHttpClient.cs`.
- Statt einer einzelnen `PgChangeFeed.Client.Tests/Sse/*.cs`-Sammlung liegen
  die Tests in fünf Dateien (`SseFrameParserTests.cs`, `TestClientFactory.cs`,
  `PgChangeFeedSseClientTests.cs`, `PgChangeFeedSseClientAuthBoundaryTests.cs`,
  `ChangeMessageSchemaTests.cs`) — dasselbe Gruppierungs-Muster wie bei
  `slice-sdk-csharp-grpc-client-flaeche`s Plan-Nachzug (`Grpc/`-Ordner),
  reine Lesbarkeits-Entscheidung, kein fachlicher Umfangsunterschied.
  `TestClientFactory.cs` referenziert den bestehenden
  `PgChangeFeed.Client.Tests.Http.FakeHttpMessageHandler` direkt
  (`internal` ist assembly-, nicht namespace-scoped) statt einen zweiten
  Fake-Handler anzulegen.
- Fehlerabbildung (`BuildException`/`ExtractErrorMessage`) liegt als eigene,
  kleine private Methodenpaar-Kopie in `PgChangeFeedSseClient.cs`, statt sie
  aus `PgChangeFeedHttpClient.cs` zu extrahieren — bewusste Design-
  Entscheidung: dieselbe Unabhängigkeit, die bereits zwischen HTTP- und
  gRPC-Fläche besteht (Letztere bildet gRPC-`StatusCode` nativ ab, kein
  gemeinsamer Mapper), statt eines produktiven Refactors an einer
  bestehenden, stabilen Datei außerhalb dieses Slice-Umfangs.
- `AGENTS.md` §3.13-Suchlauf (Träger-Nachzug) real ausgeführt:
  `grep -rln "SSE" sdks/csharp/` traf `sdks/csharp/README.md` §Status
  („SSE and NATS-vollinhalts delivery remain uncovered by this package") —
  korrigiert, damit die Aussage nach diesem Slice wieder zutrifft (SSE jetzt
  gedeckt, NATS-Vollinhalt bleibt offen). `spec/pflichtenheft.md`/
  `docs/user/benutzerhandbuch.md` bleiben davon unberührt (siehe DoD-Punkt
  „Doku-Update" — bewusst gebündelt im Folge-Slice, kein Träger-Nachzug-Fund
  dieser Klasse).

## 4. Trigger

**Start** (`next` → `in-progress`): wenn diese Welle eröffnet ist und kein
anderer Slice in `in-progress/` liegt (WIP-Limit 1).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): entfällt aus
  heutiger Sicht — ein Endpunkt, ein Nachrichtenschema, bereits zweimal in
  diesem Package erprobtes Muster (HTTP/gRPC).
- `in-progress` → `open` (blockiert — Carveout?): der Frame-Parser lässt
  sich aus einem noch unbekannten Grund nicht ohne echten Server
  netzlos testen (unwahrscheinlich — `examples/csharp/sse-client` belegt
  den Mechanismus bereits real, und die bestehende HTTP-Fläche belegt das
  Fake-Handler-Muster für diesen Baum).

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + Closure-Notiz mit Lerneintrag
geschrieben.

## 6. Risiken und offene Punkte

- Ein gestubbter/gefakter HTTP-Response-Stream könnte das reale Chunked-
  Transfer-/Flush-Verhalten des Server-`http.Flusher` nicht exakt
  nachbilden. **Ausgang:** weiter offen — ein realer Rundlauf-Beleg bleibt
  `make test-integration`s bestehendem `tools/harness/sseclient`
  vorbehalten (Wegwerf-Client, kein SDK-Import); dieses SDK bekommt (wie
  bei HTTP/gRPC) keinen eigenen Integrationsbeleg in dieser Welle
  (Welle-Plan §6, Begründung Vorarbeits-Parität).
- Die Docker-Bau-Kopplung (`--build-context proto=proto` auch ohne
  gRPC-Bezug) könnte bei einem künftigen, isolierten Bau-Versuch dieses
  Slice übersehen werden, wenn jemand die DoD-Begründung „SSE braucht kein
  proto" wörtlich nimmt. **Ausgang:** eingetreten — der Verifier prüfte
  real (per `git log`) den Ursprung der Kopplung im SDK-Baum
  (`sdks/csharp/Dockerfile`, seit `slice-sdk-csharp-grpc-client-flaeche`)
  und der Planner hat sie als dritten Beleg der Klasse
  `BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut` registriert —
  die vorab benannte 3×-Schwelle ist damit real erreicht
  (`evidence/slice-sdk-csharp-sse-client-flaeche.md`). Die daran hängende
  Regelschärfungs-Frage (Architect-Entscheidung, Modul 4/8) bleibt der
  Closure von `welle-sdk-csharp-vollabdeckung` vorbehalten — die Welle
  liegt noch offen (ein weiterer Slice folgt), kein Teil dieser
  Einzel-Slice-Closure.

## 7. Closure-Notiz

- **Was hat funktioniert:** Ein glatter Durchlauf ohne Fixrunde — der
  Reviewer fand 0 HIGH/0 MEDIUM (3 LOW, 1 INFO,
  `docs/reviews/review-slice-sdk-csharp-sse-client-flaeche.md`) und der
  Verifier bestätigte „DoD erfüllt"
  (`docs/reviews/verifikation-slice-sdk-csharp-sse-client-flaeche.md`), je
  mit eigenständig nachgefahrenen, nicht übernommenen Docker-Build-Läufen
  (mit und ohne `--build-context proto=proto`, insgesamt vier unabhängige
  Bau-Bestätigungen über Implementer/Reviewer/Verifier). Der Frame-Parser
  und alle zehn `SPEC-021`-Felder wurden von beiden Rollen eigenständig
  gegen den Quellcode gehalten, nicht nur aus Bericht oder Commit-Message
  übernommen. Der DoD-Checkbox-Nachzug ohne Fixrunde griff hier korrekt:
  der Reviewer zog die Checkbox „Review durchgeführt" im selben Commit
  selbst nach.
- **Was ging anders als geplant:** Keine inhaltliche Abweichung vom Plan —
  die drei im Implementer-Zug benannten Plan-Nachzüge (§3: eigene
  `Sse/Models/Change.cs`-Datei, fünf statt eine Testdatei, kein
  Mapper-Extrakt aus `Http/`) sind alle begründete Ausgestaltung, kein
  Plan-Konflikt (von Reviewer und Verifier unabhängig geprüft). Der
  Verifier ergänzte einen über die reine DoD-Prüfung hinausgehenden
  Hinweis: die Docker-Bau-Kopplung in `sdks/csharp/Dockerfile` (§6) ist
  strukturell derselben Beobachtungsklasse zuordenbar wie die bereits
  registrierten `examples/csharp`/`examples/kotlin`-Fälle, aber bislang nie
  für den SDK-Baum selbst belegt — der Planner hat das bei der Closure
  real per `git log` nachgezogen (siehe unten).
- **Steering-Loop-Eintrag:** kein neuer Sensor, keine neue Hard Rule durch
  diese Closure selbst. Der Verifier-Hinweis zur Docker-Bau-Kopplung
  bestätigt eine bereits vorab (Welle-Eröffnungs-Sichtung, Slice-Plan §3/§6)
  erwartete Schwellen-Erreichung — die vorab benannte Konsequenz
  („dritter Treffer würde die 3×-Schwelle erreichen und eine
  Reviewer-Skill-Schärfung auslösen") ist eingetreten; die Schärfung selbst
  bleibt eine Architect-Entscheidung bei der Welle-Closure (siehe unten),
  kein Ergebnis dieser Slice-Closure. Die vier Reviewer-Findings (F-1…F-4,
  3 LOW/1 INFO) lösen **keine** neuen Beobachtungs-Register-Einträge aus:
  jedes ist ein Erstauftreten ohne belegtes Wiederholungsmuster in diesem
  Repo (eigene Prüfung: kein bestehender `BEO-PGC/*`-Eintrag trägt „Kopie
  statt Extraktion", „Vorbild-Verhalten übernommen, Spec-Lücke latent",
  „Discriminator ungeprüft" oder „Randfall ohne Testbeleg"), und der Plan
  begründet den zugrunde liegenden Design-Verzicht (keine Mapper-Extraktion
  aus `Http/`) bereits explizit als bewusste Entscheidung (§3). F-1
  (fehlende `data:`-Verkettung) und F-2 (wörtliche Kopie des
  Fehler-Mappings) sind die beiden mit dem größten Wiederholungsrisiko —
  falls der Folge-Slice `slice-sdk-csharp-nats-stream-client-flaeche`
  denselben Frame-/Mapping-Zuschnitt wiederholt (eigener Parser, eigene
  Fehler-Mapping-Kopie), wäre das ein zweites Auftreten und ein Kandidat für
  einen neuen Registereintrag — an dieser Stelle benannt, damit der
  Folge-Slice-Reviewer es nicht neu herleiten muss.
- **Beobachtungs-Register (`../observations/`):** Beleg ergänzt zu
  `BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut/` — neue
  Evidenz-Datei `evidence/slice-sdk-csharp-sse-client-flaeche.md`. Der
  Planner hat den Verifier-Hinweis geprüft: `git log --diff-filter=A
  --oneline -- sdks/csharp/Dockerfile` und `git log --oneline --
  sdks/csharp/Dockerfile` zeigen real, dass die unbedingte
  `COPY --from=proto …`-Zeile bereits mit
  `06e78545 feat(sdk): C#-gRPC-Client-Fläche für PgChangeFeed.Client`
  (Slice `slice-sdk-csharp-grpc-client-flaeche`) in die einzige
  `build`-Stufe kam — kein neuer Fund dieses Slice, sondern eine
  vorbestehende, hier erstmals belegte Instanz. Zähler jetzt real
  ausgezählt **3×** (`evidence/slice-102.md`, `evidence/slice-103.md`,
  `evidence/slice-sdk-csharp-sse-client-flaeche.md`) — die 3×-Schwelle ist
  erreicht, aber die dritte Instanz hat eine strukturell andere Ursache
  als die ersten beiden (ein bewusst gebündeltes `.csproj`/Package,
  `ADR-0106` Festlegung 1, statt vermeidbar geteilter Docker-Stufen über
  eigenständige Programme — Details in der Evidenz-Datei §Einordnung). Der
  `Ausgang` dieser Beobachtung ist **nicht** Teil dieser Slice-Closure:
  `welle-sdk-csharp-vollabdeckung` liegt noch offen (ein weiterer Slice,
  `slice-sdk-csharp-nats-stream-client-flaeche`, folgt) und hat den
  dritten Treffer bereits bei ihrer eigenen Eröffnungs-Sichtung (§6)
  vorausgesehen — der Lese-Schritt für die Regelschärfungs-Frage
  (Architect-Entscheidung, Modul 4/8) bleibt der Closure dieser Welle
  vorbehalten (`state.md` trägt den Verweis). Zusätzlich geprüft, kein
  Handlungsbedarf: `BEO-PGC/deutsches-fachwort-im-englischen-sdk-readme`
  (bereits 3×, Architect-Sichtung ausstehend, unverändert durch diesen
  Slice), `BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`
  (verkörpert, DoD-Punkt „Doku-Update" bewusst im Folge-Slice gebündelt,
  kein Verstoß), `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert,
  Suchlauf real durchgeführt — siehe Plan-Nachzug §3, Fund in
  `sdks/csharp/README.md` behoben).
- **Folge-Slices:** keine aus diesem Slice selbst erwartet — der nächste
  Slice `slice-sdk-csharp-nats-stream-client-flaeche` ist bereits Teil
  derselben Welle und bereits geplant (Welle-Plan §4), kein neuer
  Folge-Slice-Bedarf durch diesen Zug.
- **Risiken aus §6:**
  - „Gestubbter HTTP-Response-Stream könnte reales Chunked-Transfer-/
    Flush-Verhalten nicht exakt nachbilden" — **Ausgang: weiter offen**,
    unverändert zur Plan-Begründung; ein realer Rundlauf-Beleg bleibt
    `make test-integration`s `tools/harness/sseclient` vorbehalten, dieses
    SDK bekommt (wie bei HTTP/gRPC) keinen eigenen Integrationsbeleg in
    dieser Welle.
  - „Docker-Bau-Kopplung könnte bei einem isolierten Bau-Versuch übersehen
    werden" — **Ausgang: eingetreten** — der vorab benannte dritte Treffer
    der Klasse `BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut`
    ist real eingetreten und belegt (siehe Beobachtungs-Register oben); die
    daran hängende Regelschärfungs-Entscheidung ist an die Closure von
    `welle-sdk-csharp-vollabdeckung` verwiesen, kein offener Punkt dieses
    Einzel-Slice mehr.
- **Drei Paarungen:** dieser Slice gehört zu
  [welle-sdk-csharp-vollabdeckung](welle-sdk-csharp-vollabdeckung.md)
  (noch offen) — die Prüfung läuft regelkonform bei deren Closure.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Sub-Area `sdks/csharp/` — bereits
mit `slice-sdk-csharp-projektgeruest` eröffnet (GF), keine erneute
Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut` (2×, §3 dieses
Plans trägt sie bereits), `BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`
(verkörpert, DoD-Punkt „Doku-Update" liegt gebündelt im Folge-Slice),
`BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert, Suchlauf-Pflicht
`AGENTS.md` §3.13 gilt unverändert).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Fortsetzung von
`slice-sdk-csharp-projektgeruest`).
