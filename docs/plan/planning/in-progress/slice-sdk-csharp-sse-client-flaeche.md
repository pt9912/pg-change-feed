# Slice sdk-csharp-sse-client-flaeche: Öffentliche SSE-Stream-Client-Fläche (`SPEC-021`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** [welle-sdk-csharp-vollabdeckung](../welle-sdk-csharp-vollabdeckung.md).

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
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update (`docs/user/benutzerhandbuch.md`,
      `spec/pflichtenheft.md`): bewusst **nicht** in diesem Slice —
      gebündelt im Folge-Slice `slice-sdk-csharp-nats-stream-client-flaeche`
      (§1 Abgrenzung), damit beide neuen Flächen zusammen in die
      Träger-Dokumente einfließen. Beide Dateien real unberührt
      (`git diff --stat` gegen diesen Diff), `PgChangeFeed.Client.csproj`s
      `<Version>` unverändert `0.1.0`.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben,
      **falls dieser Slice einen Inventur-Fund auflöst** — entfällt: keine
      Reconciliation-Datei in diesem Repo.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder Beleg in `evidence/`; keine Beobachtung angefallen
      ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      dieser Slice gehört zu
      [welle-sdk-csharp-vollabdeckung](../welle-sdk-csharp-vollabdeckung.md)
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
  proto" wörtlich nimmt. **Ausgang:** weiter offen, aber transparent
  benannt (§3) — dritter Treffer der Klasse
  `BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut` würde die
  3×-Schwelle erreichen und eine Reviewer-Skill-Schärfung auslösen.

## 7. Closure-Notiz

- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Steering-Loop-Eintrag:** <…>
- **Beobachtungs-Register (`../observations/`):** <…>
- **Folge-Slices:** <…>
- **Risiken aus §6:** <…>
- **Drei Paarungen:** dieser Slice gehört zu
  [welle-sdk-csharp-vollabdeckung](../welle-sdk-csharp-vollabdeckung.md)
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
