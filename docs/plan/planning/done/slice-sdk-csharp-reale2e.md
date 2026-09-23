# Slice sdk-csharp-reale2e: C#-SDK — vier Zustellweg-Flächen mit Realserver-Beleg, Abdeckungs-Träger + `trace.coverage`-Eintrag

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** [welle-sdk-reale2e](../welle-sdk-reale2e.md).

**Bezug:** [`LH-FA-SST-009`](../../../../spec/lastenheft.md)
(Client-Bibliotheken — die Fläche existiert, der Beleg-Stand wird stärker),
[`LH-FA-SST-008`](../../../../spec/lastenheft.md) (Live-Streaming —
gRPC/SSE/NATS-Flächen, Boundary: kein Replay im Stream selbst),
[`LH-FA-SST-006`](../../../../spec/lastenheft.md) (HTTP-API),
[`ADR-0110`](../../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
§Entscheidung Festlegung 2/Folgepflicht 1 (die etablierte Mechanik-Klasse,
hier auf den C#-Baum gespiegelt), [`ADR-0106`](../../adr/0106-csharp-nuget-erstes-sdk-package.md)
(Ort `sdks/csharp/`, Import-Grenze, unverändert gültig), [`ADR-0060`](../../adr/0060-grpc-streaming-mechanismus.md)/[`ADR-0061`](../../adr/0061-http-sse-zusaetzlich-zu-grpc.md)/[`ADR-0100`](../../adr/0100-nats-dritter-vollinhalts-zustellweg.md)
(Server-Verträge der drei Stream-Wege, werden vom SDK benutzt, nicht erweitert).

**Berührte Spec-Stellen:** [`SPEC-018`](../../../../spec/pflichtenheft.md),
[`SPEC-020`](../../../../spec/pflichtenheft.md),
[`SPEC-021`](../../../../spec/pflichtenheft.md),
[`SPEC-024`](../../../../spec/pflichtenheft.md) — gelesen als Draht-Vertrag,
nicht geändert.

**Verantwortlich:** Implementer-Agent, 2026-09-23.

**Autor:** Planner-Agent, Welle-Eröffnung
[welle-sdk-reale2e](../welle-sdk-reale2e.md). **Datum:** 2026-09-23.

---

## 1. Ziel und Abgrenzung

**Ziel:** Der Realserver-Beleg für die **vier C#-Zustellweg-Flächen** des
Packages `PgChangeFeed.Client` (HTTP `SPEC-018`, gRPC `SPEC-020`, SSE
`SPEC-021`, NATS-Vollinhalt `SPEC-024`) — Mechanik im Muster des
Python-Vorbilds ([`ADR-0110`](../../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
§Entscheidung Festlegung 2, real gebaut in
`slice-sdk-python-grpc-client-flaeche`, in `done/`): (a) eine additive
Docker-Stufe `integration` in `sdks/csharp/Dockerfile` (baut auf `build`
auf — dasselbe installierte SDK, kein zweiter Build-Pfad); (b) ein
Runner-Skript `tools/harness/run-sdk-csharp-integration-tests.sh`
(Arbeitsname) — Bring-up der `compose.yaml`-Umgebung im Muster
`tools/harness/run-sdk-python-integration-tests.sh` (PostgreSQL/NATS/
Feed-Container, Schema-Rollout über d-migrate, Vorbedingungen der
Aktivierung), Bau/Stau der `integration`-Stufe im selben Docker-Netz wie
der Feed-Container, eine Phase je Fläche mit READY/RECEIVED/REJECTED-
Markern; der C#-Prüfling ist die **kompilierte Client-Assembly** — der
Integrationstest importiert `PgChangeFeed.Client` direkt, kein
Wegwerf-Duplikat-Client daneben; (c) ein Make-Target
`test-sdk-csharp-integration` in `harness/mk/sdk.mk`; (d) der Abdeckungs-
Träger `docs/user/sdk-e2e-abdeckung.md` samt `.d-check.yml`
`trace.coverage`-Eintrag (Label `SDK-E2E`) — beide in diesem Slice, weil
der C#-Runner der erste Abschnitts-Erzeuger der Datei ist
([Welle-Plan §4](../welle-sdk-reale2e.md), Reihenfolge-Grund 1).

**Je Fläche der Beleg (die vier Phasen des Runners):**

1. **gRPC (`SPEC-020`):** `PgChangeFeedGrpcClient` öffnet real den
   Server-Stream gegen den laufenden Feed-Container
   (`pg-change-feed:9090`, `authorization`-Metadata mit Bearer-Token) und
   empfängt eine danach über `psql` committete Änderung — belegt am
   Stream-Image (Tabelle, Operation, Row-Image) und über die `change_id`
   gegen `cdc.changes` gehalten; ein Öffnungsversuch ohne gültiges Token
   endet mit gRPC-Status `Unauthenticated`.
2. **SSE (`SPEC-021`):** `PgChangeFeedSseClient` öffnet real
   `GET /changes/stream` (`http://pg-change-feed:8090`, reader-Token) und
   empfängt eine danach committete Änderung; `change_id` gegen
   `cdc.changes` gehalten; ein Aufruf ohne gültiges Token endet mit
   HTTP-Status 401.
3. **NATS-Vollinhalt (`SPEC-024`):** `PgChangeFeedNatsStreamClient`
   verbindet sich real (`nats://nats:4222`, `CDC_NATS_STREAM_TOKEN`)
   und empfängt eine danach committete Änderung als vollständiges
   JSON-Event über `cdc.stream.<source_id>.<schema>.<table>`;
   `change_id` gegen `cdc.changes` gehalten; ein Verbindungsversuch mit
   falschem Token wird vom NATS-Server abgelehnt.
4. **HTTP (`SPEC-018`):** `PgChangeFeedHttpClient` läuft real gegen
   `http://pg-change-feed:8090` — ein Rundlauf im Muster des
   Server-E2E-HTTP-Rundlaufs (`run-integration-tests.sh`-Phase
   „HTTP-API-Rundlauf"): Registrierung eines Wegwerf-Consumers und
   Listen-Aufruf über die SDK-Methoden, die Registrierung über
   `cdc.consumer` gegen den SQL-Lesezugriffsweg gehalten; ein Aufruf
   ohne gültiges Token endet mit HTTP-Status 401. (Die neun Fähigkeiten
   einzeln zu fahren, ist hier nicht der Pflichtkern — die Unit-Tests
   tragen die Methoden-Fläche bereits netzlos gegen Fakes; der
   Realserver-Rundlauf belegt die Protokoll-Annahmen am Wire:
   Auth-Header-Form, JSON-Mapping, Abbruch-Verhalten.)

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Kotlin-Flächen** — eigener Folge-Slice
  (`slice-sdk-kotlin-reale2e`), der dieselbe, hier durchlaufene Form
  spiegelt (Welle-Plan §4, Grund 2).
- **Die Python-HTTP-Fläche** — eigener Folge-Slice
  (`slice-sdk-python-http-reale2e`); er erweitert das bestehende
  Python-Runner-Skript, nicht dieses.
- **Version-Bump oder Publish-Workflow-Änderung** — der Realserver-Test
  ist kein Artefakt-Vertrag; `<Version>` in der `.csproj` und
  `.github/workflows/sdk-csharp-release.yml` bleiben unverändert
  (Welle-Plan §6).
- **Ein `examples/csharp/`-Bezug** — die Beispiele bleiben Doku mit
  Bau-Bindung (`SPEC-023`), unberührt.
- **`spec/pflichtenheft.md`-Nachzug** — kein falsch werdender Träger
  (Welle-Plan §6, real geprüft 2026-09-23); `SPEC-027`s Deckungs-Aussage
  bleibt richtig.
- **Eine Aufnahme des Targets in `make gates`** — braucht DB-Zugang/
  Docker/Netz, dieselbe Klasse wie `make test-integration` (Welle-Plan §6).

## 2. Definition of Done

- [x] `LH-FA-SST-009`-Beleg-Stand stärker: alle vier C#-Flächen tragen
      je einen realen Rundlauf gegen eine laufende Server-Instanz
      (§1, Phasen 1–4), jede Phase mit SQL-Gegenprüfung der empfangenen
      `change_id` gegen `cdc.changes` (gRPC/SSE/NATS) bzw. der
      Registrierung gegen `cdc.consumer` (HTTP) und je einem
      Ablehnungs-Beleg ohne gültiges Token (gRPC `Unauthenticated`,
      SSE/HTTP Status 401, NATS-Token-Abweisung durch den Server).
      *(Sensor-Beleg: `make test-sdk-csharp-integration` EXIT=0 — der
      Lauf trug gRPC change_id=804-1, SSE 807-1, NATS 810-1 (je
      SQL-Gegenprüfung gegen `cdc.changes`) und die Consumer-Registrierung
      (consumer_id=csharp-sdk-e2e-20260923075217, unabhängig über
      `cdc.consumer` lesbar); die Ablehnungs-Belege trugen je Phase real
      gRPC `Unauthenticated`, HTTP `401` und die NATS-Verbindungsablehnung.
      Der Endstand nach den Lauf-Korrekturen (Linktiefe des Trägers)
      wurde in einem erneuten Lauf bestätigt.)*
- [x] **Mechanik** ([`ADR-0110`](../../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
      §Entscheidung Festlegung 2, gespiegelt): additive Docker-Stufe
      `integration` in `sdks/csharp/Dockerfile` (baut auf `build` auf,
      **dieselben Pins wiederverwendet** — `mcr.microsoft.com/dotnet/sdk:10.0@sha256:60a2…`
      aus der bestehenden `build`-Stufe, **keine neuen Pins**),
      Runner-Skript `tools/harness/run-sdk-csharp-integration-tests.sh`
      (Bring-up im Python-Runner-Muster, `set -euo pipefail`, Cleanup
      je Ausgang, Testdatei-/Phase-Auswahl explizit — kein stiller
      Ausschluss, `BEO-PGC/test-runner-stiller-ausschluss`-Disziplin),
      Make-Target `test-sdk-csharp-integration` in `harness/mk/sdk.mk`.
      Kein Gate. *(Sensor-Beleg: realer Lauf des Targets EXIT=0 — siehe
      oben; der Runner fuhr alle vier Phasen gegen den laufenden
      Feed-Container.)*
- [x] **Abdeckungs-Träger + `trace.coverage`-Eintrag im selben Zug**:
      `docs/user/sdk-e2e-abdeckung.md` entsteht real (generiert vom
      Runner, Marker-gegrenzter C#-Abschnitt — die Datei trägt nur Zeilen
      real existierender Runner-Phasen, Muster
      [`e2e-abdeckung.md`](../../../../docs/user/e2e-abdeckung.md)) und
      `.d-check.yml` trägt den `trace.coverage`-Eintrag (Label
      `SDK-E2E`) — die RTM sieht die C#-Belege ab diesem Slice; die
      Kotlin-/Python-Abschnitte kommen mit deren Slices (Welle-Plan §4).
- [x] Kein Import aus `internal/**`/`cmd/**`/`gen/**` dieses Repos in
      `sdks/csharp/**` — der neue Integrationstest-Quelltext inklusive
      (Import-Zeilen-Prüfung, Muster der bestehenden SDK-Slices).
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6),
      kein Self-Review (Modul 8). *(Report `review-slice-sdk-csharp-reale2e`:
      0 HIGH · 1 MEDIUM · 4 LOW · 2 INFO; Fixrunde `86892bdb` löste
      F-1 bis F-6 (F-7 INFO als Runner-Kommentar); die Verifikation
      urteilte „erfüllt" mit zwei nicht-blockierenden Beobachtungen,
      beide gelöst in `7bede9f2` — der Nachzug bei Schritt 21 des
      Minimal Agent Workflow geschieht in diesem Closure-Commit
      (`BEO-PGC/dod-checkbox-nachzug`).)*
- [x] Doku-Update im selben Zug: `harness/README.md` bekommt
      die neue `make test-sdk-csharp-integration`-Zeile, weil das Target
      hier real entsteht ([`AGENTS.md`](../../../../AGENTS.md) §4: kein
      Träger nennt ein Target, das es nicht gibt — umgekehrt: eine Zeile
      über ein Target trägt ihren realen Lauf, der vor der Zeile gefahren
      ist; `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag. *(§7 unten —
      Lerneintrag: der Suchraum des §3.13-Suchlaufs folgt der
      Schreib-Verbreiterung der Bewegung; Beleg im Beobachtungs-Register.)*
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in
      diesem Repo.
- [x] Beobachtungs-Register fortgeschrieben — kein Anfall ist ebenfalls
      eine Antwort und wird in §7 notiert. *(Beleg:
      `BEO-PGC/arbeit-ueberholt-stehenden-traeger`,
      `evidence/slice-sdk-csharp-reale2e.md`; übrige Kandidaten geprüft,
      siehe §7.)*
- [x] Jedes Risiko aus §6 trägt einen Ausgang. *(vier Ausgänge in §6, je
      entfallen mit Beleg am Ort; siehe §7.)*
- [x] Die drei Paarungen getragen — dieser Slice gehört zu
      [welle-sdk-reale2e](../welle-sdk-reale2e.md). *(Anker: die Regel im
      Lerneintrag liegt in `AGENTS.md` §3.13 (`seit welle-20`, am Ort
      existent); Folge-Slice: beide genannten Pläne existieren in
      `open/`; Register: der neue Beleg liegt in `evidence/` — §7.)*

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/run-sdk-csharp-integration-tests.sh` | neu | Runner: Bring-up der Compose-Umgebung, Bau/Start der `integration`-Stufe im Netz `cdc-feed-test`, vier Phasen mit READY/RECEIVED/REJECTED-Markern, SQL-Gegenprüfung je Phase, Cleanup. Form-Vorbild: `tools/harness/run-sdk-python-integration-tests.sh` (unverändert als Quelle). |
| `sdks/csharp/PgChangeFeed.Client.Integration/` (Arbeitsname; eigenes Testprojekt) | neu | der Integrationstest importiert die kompilierte `PgChangeFeed.Client`-Assembly direkt; eigenes Projekt statt zweitem Test-Ordner im bestehenden Tests-Projekt, damit der Unit-Lauf (`dotnet test PgChangeFeed.Client.Tests`) unberührt bleibt und die `integration`-Stufe das Integrationsprojekt gezielt aufruft (kein stiller Ausschluss des Rests — das C#-Äquivalent der Python-`integration/`-Ordner-Form). |
| `sdks/csharp/Dockerfile` | update | additive Stufe `integration` (baut auf `build` auf); Pins unverändert (`mcr.microsoft.com/dotnet/sdk:10.0@sha256:60a2…`, Wiederverwendung aus `build`, kein neuer Pin). |
| `harness/mk/sdk.mk` | update | neues Target `test-sdk-csharp-integration` (Werkzeug-Block im Kommentar-Stil der bestehenden Targets). |
| `docs/user/sdk-e2e-abdeckung.md` | neu | der Abdeckungs-Träger, C#-Abschnitt als Erzeugnis des Runners (idempotent, Marker-gegrenzt); Form nach `docs/user/e2e-abdeckung.md`. |
| `.d-check.yml` | update | `trace.coverage`-Eintrag `- files: [docs/user/sdk-e2e-abdeckung.md] / label: SDK-E2E`. |
| `harness/README.md` §Werkzeuge | update | neue Zeile für `make test-sdk-csharp-integration` (nach dem realen Lauf geschrieben). |

**Plan-Nachzug (im selben Lauf, vor dem Gate-Lauf):**

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| Testklassen-Form (4 Klassen à 2 Tests statt 4 Einzelfiles) | Abweichung | der Runner selektiert je Phase per `dotnet test --filter FullyQualifiedName~<Klasse>` — die Testklasse trägt Happy-Path **und** Reject-Beleg je Phase (dieselbe Zusammenfassung wie die Python-Testdateien); `PGCHANGEFEED_TEST_NAME` (mit `:?`-Guard, Muster Python-CMD-Guard) ersetzt den Arbeitsnamen „Testdatei". |
| `sdks/csharp/PgChangeFeed.Client.Integration/PhaseEnvironment.cs` | neu | gemeinsame Env-Auslese + Marker-Druck (`Console.Out.Flush()` — Risiko-§6 Pufferung) der vier Phasen; jede Variable an ihrer Eingabeseite gebunden (`Required`). |
| HTTP-Phase: `PGCHANGEFEED_SOURCE_ID`/`PGCHANGEFEED_HTTP_PUBLICATION` als Env statt festen Strings im Test | Erweiterung | die ListTables-Eingabe (Quelle, Publikation) ist eine Eingabeseite der Zusage — gebunden über die Umgebung, dieselbe Disziplin wie die Token-Form (Review-Klasse des Vorgänger-Slices F-6). |
| HTTP-Phase: eigener Sentinel-/ID-Bereich (430ff., eigener Sentinel-Name) | Erweiterung | die Phase-Funktion committet Fire-and-Forget-Inserts in derselben Tabelle; ein wiederverwendeter ID-Bereich kollidierte real mit der gRPC-Phase (PK-Konflikt sichtbar im ersten Lauf) — eigene Bereiche je Phase (400/410/420/430). |
| Träger-Schreiber erhält fremde Abschnitte | Erweiterung | der Runner ersetzt nur seinen marker-gegrenzten C#-Abschnitt und erhält den Inhalt hinter dem end-Marker — die Kotlin-/Python-Abschnitte der Folge-Slices bleiben bei jedem Re-Run bestehen. |
| C#-Namespaces: `Change`-DTO je Fläche eigen (Sse/Nats/Models, Cdc.Stream.V1) | Abweichung | die drei Stream-Flächen tragen drei eigenständige DTOs (vier unabhängige Wire-Verträge, Kommentar im SSE-Model); der Test importiert je Fläche ihren eigenen Typ — kein geteilter Alias. |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „die C#-Flächen tragen reale Realserver-Belege"; beide Stände gemessen: Parent `fce7af10` und HEAD):**

| Träger | Befund | Behandlung |
|---|---|---|
| `harness/README.md` §Werkzeuge | Zeile fehlte (Target real, Zeile nicht) | in diesem Zug ergänzt (nach dem realen Lauf, `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`) |
| `harness/README.md` §Werkzeuge (`make doc-trace`-Zeile) | Coverage-Aufzählung ohne den vierten Eintrag gefunden (Review F-1a); Zahlen „76/55/0" gegen die Messung 79/2 überholt | gezogen: vier Dimensionen + Messung 2026-09-23 (79 Anforderungen, 2 Waisen) |
| `harness/sensors/docs-check.md` §Grenze | Aufzählung mit einer Coverage-Datei gefunden (Review F-1b) | gezogen: vier Dateien |
| `docs/plan/planning/welle-sdk-reale2e.md` §6 (Träger-Form-Behauptung) | „`Datei:Zeile`-Orte (Muster e2e-abdeckung.md)" gefunden — der reale Träger trägt Runner-Verweis-Orte ohne Zeilenanker (Review F-1c) | gezogen: echte Form am Ort benannt |
| `sdks/csharp/README.md` §Status | geprüft — trägt keine Teststrategie-Aussage über Realserver-Läufe | kein Nachzug nötig |
| `docs/user/benutzerhandbuch.md` | geprüft — trägt die SDK-Hinweise, keine E2E-Beleg-Aussage | nichts zu ziehen |
| `spec/pflichtenheft.md` | geprüft — `SPEC-027`s Deckungs-Aussage bleibt richtig (Welle-Plan §6) | kein falsch werdender Träger |

**Ansatz:** Struktur-Vorbild ist der Python-Runner — eigenständiges
Skript, geteilte Compose-Umgebung, SDK als Prüfling, keine Server-E2E-
Rollen-DSN-/Retention-Prüfungen (entschieden in
`slice-sdk-python-grpc-client-flaeche` §3, `done/`; hier gespiegelt).
Die vier Phasen folgen derselbe `run_surface_phase`-Form: begrenztes
Fire-and-Forget-Fenster, eindeutige Sentinel- und ID-Wertebereiche je
Phase, Reject-Marker je Protokoll (`REJECTED code=Unauthenticated`,
`REJECTED status=401`, `REJECTED token-rejected`, HTTP-Phase:
`REJECTED status=401`). Der Runner schreibt den Träger-Abschnitt
idempotent (Marker-gegrenzt) aus derselben Messung, die ihn belegt
(`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`-Disziplin).

## 4. Trigger

**Start** (`next` → `in-progress`): wenn diese Welle eröffnet ist und kein
anderer Slice in `in-progress/` liegt (WIP-Limit 1).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls die vier
  Phasen (HTTP + drei Streams) im Runner mehr Aufwand verlangen als der
  Mechanismus-Kern — dann Abspaltung eines Mechanik-Slices (HTTP-Phase
  als eigener Folgeslice).
- `in-progress` → `open` (blockiert — Carveout?): die `integration`-Stufe
  scheitert an einem C#-spezifischen Bau-/Laufzeit-Verhalten, das der
  Python-Vorbild-Mechanik widerspricht (unwahrscheinlich — die
  `build`-Stufe baut die Assembly bereits real, `sdk-pack-csharp` läuft
  grün; der Unterschied ist die Container-Start-Form, kein neuer Pin).

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + ein realer, grüner
`make test-sdk-csharp-integration`-Lauf mit allen vier Phasen +
Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Form-Vorbild-Kopie trägt deutsches Wortfragment weiter**
  (`BEO-PGC/formvorbild-kopie-traegt-deutsches-wortfragment-weiter`,
  offen, 2×): Runner-Skript, Dockerfile-Kommentar und
  Make-Target-Kommentar entstehen als Kopie des Python-Vorbilds. Der
  Eintrag benennt genau den Prüfpunkt für ein drittes Auftreten: das
  übernommene Vorbild selbst noch einmal auf Sprachreinheit prüfen,
  statt die Übereinstimmung mit dem Vorbild als ausreichenden Beleg zu
  werten. *Zusage:* die übernommenen Form-Teile werden je separat
  gesichtet; der Suchlauf im Bericht trägt das Ergebnis.
  **Ausgang:** entfallen — die je-teilige Sichtung lief real (Review
  „Kommentar-Sichtung je übernommenem Form-Teil": Runner-Kopf,
  Phasen-Prosa, Dockerfile-Kommentar, Make-Target-Kommentar je separat
  gegen das Python-Vorbild); kein unübersetztes deutsches Wortfragment,
  der Zähler der Klasse bleibt 2×. Gefunden wurden stattdessen zwei
  andere Klassen (F-3 Grammatik-Slip, F-4 ASCII-Mischform), beide gelöst
  (Fixrunde `86892bdb`; Umlaut-Residuum des F-4-Fixes `7bede9f2`).
- **.NET-stdout-Pufferung vs. `docker logs`-Marker-Polling:** der Runner
  liest READY/RECEIVED/REJECTED über `docker logs`, während der Test
  läuft — .NET-Pufferung könnte die Marker verzögern (Python-Lösung:
  `python -u`). *Erwartet, zu belegen durch:* der erste reale Lauf
  zeigt die Marker fristnah; falls nicht, trägt der Fix den
  ungepufferten Ausgabe-Form (z. B. `Console.Out.Flush()` je Marker)
  und den Beleg, dass der Runner sie liest. **Ausgang:** entfallen —
  die Marker trafen fristnah in allen drei Läufen (Implementer-Lauf,
  Review-Lauf, Verifier-Lauf, je EXIT=0); die Phase-Form druckt je
  Marker mit `Console.Out.Flush()` (Plan-Nachzug, `PhaseEnvironment.cs`).
- **zwei Token-Klassen in der HTTP-Phase:** der Server-E2E-HTTP-Rundlauf
  nutzt das `admin`-Token für RegisterConsumer und das `reader`-Token
  für Listen/Lesen (Container-Vertrag `CDC_API_TOKEN_ADMIN`/
  `CDC_API_TOKEN_READER`); der Python-Runner trägt nur den reader-Token.
  Die HTTP-Phase des C#-Runners braucht beide — der Ablehnungs-Beleg
  (401) bleibt beim Token-freien Aufruf. *Erwartet, zu belegen durch:*
  der reale Lauf; die Env-Form folgt dem bestehenden
  `run-integration-tests.sh`-Muster. **Ausgang:** entfallen — die
  HTTP-Phase fährt beide Klassen real (admin-/reader-Paar je Env-gebunden
  `PGCHANGEFEED_API_TOKEN_ADMIN`/`PGCHANGEFEED_API_TOKEN_READER`, Runner);
  der Ablehnungs-Beleg (401) lief am token-freien Aufruf, der
  Listen-Aufruf am reader-Token.
- **Sentinel-/ID-Kollisionen mit anderen Läufen:** die vier Phasen
  tragen je eigene Sentinel- und ID-Wertebereiche (Muster
  Python-Runner: 300/310/320); der C#-Runner wählt eigene Bereiche,
  damit sich parallele oder nacheinander gelaufene SDK-Runner nicht in
  die Quere kommen. *Erwartet, zu belegen durch:* der reale Lauf.
  **Ausgang:** entfallen — eigene Bereiche je Phase (400/410/420/430,
  Plan-Nachzug); drei unabhängige Läufe (Implementer, Review, Verifier)
  ohne Kollision, `consumer_id` je laufgebunden statt kollidierend.

## 7. Closure-Notiz

- **Was hat funktioniert:** der Mechanik-Spiegel trug den ersten Zug ohne
  zweite Infrastruktur — die additive `integration`-Stufe baut auf der
  bestehenden `build`-Stufe auf (kein neuer Pin), und der Runner fuhr
  alle vier Phasen gegen dieselbe Compose-Umgebung:
  `make test-sdk-csharp-integration` lief real grün (Implementer-Lauf
  EXIT=0: gRPC `change_id=804-1`, SSE `807-1`, NATS `810-1`, je
  SQL-Gegenprüfung gegen `cdc.changes`, die Registrierung gegen
  `cdc.consumer`); Review- und Verifier-Läufe reproduzierten die Werte
  byte-identisch (deterministische Bring-up-Kette — nur `consumer_id`
  ist laufgebunden, Zeitstempel-Form, Verifikation §1). Der Abdeckungs-
  Träger ist byte-stabil idempotent („Abdeckungs-Traeger unveraendert"
  in allen drei Läufen). Die Rollen-Kette lief unabhängig: Haupt-Review
  F-1…F-8, Fixrunde `86892bdb` (F-1 bis F-6, F-7 als Runner-Kommentar),
  Verifikation „erfüllt" (eigener Lauf EXIT=0, `make doc-trace` —
  gemessen 2026-09-23: 79 Anforderungen, 2 Waisen —, `make gates` EXIT=0
  mit Stempel-Gleichheit); die zwei nicht-blockierenden Beobachtungen
  der Verifikation (Umlaut-Residuum im Träger, überholter Report-Satz)
  sind in `7bede9f2` gelöst.
- **Was ging anders als geplant:** der §3.13-Suchlauf trug zunächst vier
  Träger-Zeilen — alle vier einzeln bestätigt (Verifikation §3) — und
  verfehlte drei weitere (Review F-1, MEDIUM): die
  `make doc-trace`-Coverage-Aufzählung in `harness/README.md`, die
  Coverage-Dimensionen-Aufzählung in `harness/sensors/docs-check.md`
  §Grenze, die Träger-Form-Behauptung in
  [welle-sdk-reale2e](../welle-sdk-reale2e.md) §6. Alle drei hängen an
  derselben Wurzel: der `.d-check.yml`-`trace.coverage`-Eintrag dieses
  Diff verbreitert die Menge der kuratierten Coverage-Dimensionen, und
  genau diese Verbreiterung überholt jede Prosa-Zeile, die die Aufzählung
  trägt — der deklarierte Mustersatz des Plans (geplante Träger,
  `Datei:Zeile`-Form) beschrieb den Suchraum, nicht die Bewegung.
  Denselben Satz-Zug traf F-2 (die neuen Träger-Zeilen citen die
  Kopf-Anforderung nicht — RTM-Waise), gelöst in der Fixrunde; die RTM
  sieht [`LH-FA-SST-009`](../../../../spec/lastenheft.md) über das Label
  `SDK-E2E` (real gemessen, `make doc-trace` 2026-09-23, Status `ok`).
- **Steering-Loop-Eintrag (Lerneintrag):** geschärfte Regel
  (Anwendungs-Schärfung der bereits verkörperten Klasse
  `BEO-PGC/arbeit-ueberholt-stehenden-traeger`, Anker
  [`AGENTS.md`](../../../../AGENTS.md) §3.13 · seit welle-20): der
  Suchraum eines §3.13-Suchlaufs folgt der **Schreib-Verbreiterung der
  Bewegung**, nicht dem deklarierten Mustersatz — jede
  Verdrahtungs-Änderung, die eine aufgezählte Menge verbreitert (hier:
  jeder `.d-check.yml`-`trace.coverage`-Eintrag), überholt **alle**
  Prosa-Zeilen, die dieselbe Aufzählung tragen — README-Werkzeug-Zeile,
  Sensor-Doku, Welle-Plan —, auch wenn sie keinem geplanten Träger
  zugeordnet waren. Kein neuer Sensor: die verfügbare Falsifikation
  bleibt die Messung an beiden Ständen (§3.13s Grenze für
  Prosa-Aufzählungen bleibt; ein Wortmuster, das jede
  Verbreiterungs-Überholung träfe, bräuchte eine Semantik-Entscheidung,
  welche Aufzählung zu welcher Bewegung gehört). Keine benannte
  Spec-Lücke. Ob die Schärfung einen eigenen Satz in `AGENTS.md` §3.13
  trägt, prüft der Lese-Schritt der Welle-Closure — die Klasse steht mit
  diesem Beleg bei 24× (Datei-Anzahl unter `evidence/`, real
  ausgezählt), bereits verkörpert.
- **Beobachtungs-Register (`../observations/`):**
  - **`BEO-PGC/arbeit-ueberholt-stehenden-traeger`** — neuer,
    vierundzwanzigster Beleg:
    `evidence/slice-sdk-csharp-reale2e.md`; Review F-1 trägt die drei
    Fundstellen, die Fixrunde zog sie; `state.md` trägt den abgeleiteten
    Zähler.
  - **`BEO-PGC/formvorbild-kopie-traegt-deutsches-wortfragment-weiter`** —
    geprüft: kein drittes Auftreten (die je-teilige Sichtung lief real,
    Review §Eigenständig durchgeführte Prüfungen); Zähler bleibt 2×.
  - **F-3/F-4/F-5** — geprüft: keine Register-Klasse trägt sie (F-3
    „einmalig" — Grammatik-Slip; F-4 „Orthografie-Split in der
    Träger-Familie"; F-5 schmale Erzeuger-Glob-Form der
    `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`-Klasse, deren Zähler
    dadurch bei 8× bleibt) — alle drei in der Kette gelöst und im Report
    konserviert; notiert als Antwort statt Registereintrag.
  - **F-2/F-8** — geprüft: F-8 ist durch die Nachmessung entkräftet
    (kein Vorkommen); F-2 ist eine neue Kurzform („Abdeckungszuordnung
    nennt ihre Anforderung nicht"), in der Kette gelöst und als Vorlauf
    auf die Welle-Closure-Prüfung konserviert — kein Registereintrag.
- **Folge-Slices:** `slice-sdk-kotlin-reale2e` und
  `slice-sdk-python-http-reale2e` (beide in `open/`, Welle-Plan §4) —
  Review F-7 (INFO) gibt ihnen die Writer-Form-Randbedingung mit: ein
  Folge-Runner ersetzt nur seinen eigenen marker-gegrenzten Abschnitt
  und erhält den Inhalt **hinter** seinem end-Marker, nicht den
  oberhalb (der C#-Writer regeneriert Kopf + Tabellenkopf neu, weil sein
  begin-Marker der erste der Datei ist); die Randbedingung trägt der
  Runner-Kommentar seit `86892bdb`. Kein zusätzlicher Slice aus dieser
  Closure.
- **Risiken aus §6:** Risiko 1 (Form-Vorbild-Kopie) — **entfallen** mit
  Beleg am Ort; Risiko 2 (Marker-Pufferung) — **entfallen**; Risiko 3
  (zwei Token-Klassen) — **entfallen**; Risiko 4 (ID-Kollisionen) —
  **entfallen** (alle Begründungen und Belege am Ort in §6). Kein Ausgang
  „weiter offen" ins Register — ein Registereintrag für ein nie
  aufgetretenes Muster trüge kein `evidence/` (Paarung (c)).
- **Drei Paarungen:** dieser Slice gehört zu
  [welle-sdk-reale2e](../welle-sdk-reale2e.md) (noch offen) — die Prüfung
  läuft regelkonform bei deren Closure. (a) Anker: die Regel des
  Lerneintrags liegt in `AGENTS.md` §3.13 (`seit welle-20`) — am Ort
  existent; (b) Folge-Slice: beide genannten Pläne existieren als
  Dateien in `open/`; (c) Register: jede genannte Kennung existiert als
  Verzeichnis, und das `evidence/` des genannten Eintrags trägt 24
  Dateien (real ausgezählt).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Sub-Area `sdks/csharp/` —
bereits mit `slice-sdk-csharp-projektgeruest` eröffnet (GF), keine
erneute Ausdifferenzierung nötig; `tools/harness/` ist Werkzeug-Area der
bestehenden Runner-Familie (GF, fortlaufend).

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/formvorbild-kopie-traegt-deutsches-wortfragment-weiter` (offen,
2×, einschlägig — Risiko §6), `BEO-PGC/test-runner-stiller-ausschluss`
(offen, 2×, nicht einschlägig — explizite Phase-Auswahl statt
`-run`-Muster; Deklarations-Hälfte über den Träger gedeckt),
`BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert, Suchlauf
trägt README-Status-Abschnitt und `harness/README.md`-Zeilen),
`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (verkörpert —
Träger-Abschnitt aus derselben Messung geschrieben), `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`
(verkörpert — README-Zeile nach dem realen Lauf).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Fortsetzung der
SDK-Bäume und der `tools/harness/`-Werkzeug-Familie; keine BF-Fläche
berührt — die Dockerfile-Stufe ist eine additive Erweiterung eines
GF-Artefakts).