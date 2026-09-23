# Slice sdk-python-grpc-client-flaeche: Öffentliche gRPC-Stream-Client-Fläche (`SPEC-020`), Realserver-Integrationstest-Werkzeug

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** [welle-sdk-python-vollabdeckung](../welle-sdk-python-vollabdeckung.md).

**Bezug:** [`LH-FA-SST-009`](../../../../spec/lastenheft.md),
[`LH-FA-SST-008`](../../../../spec/lastenheft.md) (Live-Change-Stream —
Boundary: kein Replay im Stream selbst),
[`ADR-0110`](../../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
§Entscheidung Festlegung 1/2 (Umfang, verschärfte Test-Pflicht — **dieser
Slice führt den Mechanismus dafür ein**), [`ADR-0107`](../../adr/0107-python-pypi-zweites-sdk-package.md)
§Entscheidung Festlegung 3 (Ort, Import-Grenze, unverändert gültig),
[`ADR-0060`](../../adr/0060-grpc-streaming-mechanismus.md)
(gRPC-Server-Streaming-Vertrag, wird vom SDK benutzt, nicht erweitert).

**Berührte Spec-Stellen:** [`SPEC-020`](../../../../spec/pflichtenheft.md)
(Nachrichtenschema, RPC-Name, Stream-Semantik).

**Verantwortlich:** — bis zur Priorisierung.

**Autor:** Planner-Agent, direkt beauftragt (`ADR-0110` §Konsequenzen
Folgepflicht 1). **Datum:** 2026-09-21.

---

## 1. Ziel und Abgrenzung

**Ziel:** Zwei Dinge in einem Slice, weil sie zusammengehören (das erste
ist ohne das zweite laut `ADR-0110` nicht abnahmefähig): (a) eine
öffentliche, stabile Python-API-Fläche für den `StreamChanges`-RPC von
[`SPEC-020`](../../../../spec/pflichtenheft.md) — ein Server-Streaming-
Aufruf über `grpcio`, der die zehn Nachrichtenfelder als Python-Iterator/
-Generator an den Consumer liefert, Authentifizierung über den
`authorization`-Metadata-Eintrag (`Bearer <token>`); (b) das
**Realserver-Integrationstest-Werkzeug**, das `ADR-0110` §Entscheidung
Festlegung 2/Folgepflicht 1 für **jede** neue Zustellweg-Fläche verlangt —
neu, weil kein Python-Vorbild und keine bestehende Infrastruktur dafür
existiert (`ADR-0110` §Kontext: kein `examples/python/`, andere Sprache
als der bestehende Go-Toolchain-Container in
`tools/harness/run-integration-tests.sh`).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **HTTP-API-Fläche** — bereits geliefert
  (`slice-sdk-python-http-client-flaeche`).
- **SSE- oder NATS-Vollinhalts-Fläche** — eigene Folge-Slices
  (`slice-sdk-python-sse-client-flaeche`,
  `slice-sdk-python-nats-stream-client-flaeche`), die das hier eingeführte
  Werkzeug **erweitern**, nicht neu bauen.
- **Ein `examples/python/`-Referenz-Client** — `ADR-0110` §Entscheidung
  Festlegung 2 mandatiert das ausdrücklich nicht (Welle-Plan §6).
- **Stream-internes Replay** — `SPEC-020`/[`LH-FA-SST-008`](../../../../spec/lastenheft.md)
  Boundary schließt das aus.
- **Version-Bump, `spec/pflichtenheft.md`-Träger-Nachzug,
  `docs/user/benutzerhandbuch.md`-Nachzug** — gebündelt im letzten
  Flächen-Slice dieser Welle (`slice-sdk-python-nats-stream-client-flaeche`).
- **Ein Ausbau des Integrationstest-Werkzeugs auf die bestehende
  HTTP-Fläche** — `ADR-0110` verlangt die verschärfte Pflicht ausdrücklich
  nur für neue Flächen (Welle-Plan §6).

## 2. Definition of Done

- [ ] `LH-FA-SST-009` erfüllt: eine öffentliche, im Package sichtbare
      Client-Klasse öffnet den `StreamChanges`-RPC (`grpcio`,
      `authorization`-Metadata mit Bearer-Token) und liefert die zehn
      Nachrichtenfelder von
      [`SPEC-020`](../../../../spec/pflichtenheft.md) — Unit-Tests
      referenzieren `SPEC-020` (Nachrichtenschema-Vollständigkeit,
      Authn-Boundary gegen einen Fake-Channel/-Stub, netzlos, Muster
      `examples/csharp/PgChangeFeed.Client/.../FakeCallInvoker.cs`-Analogie).
      Der Stub entsteht im Bau aus der `.proto` über einen zusätzlichen,
      benannten Docker-Bau-Kontext (`--build-context proto=proto`,
      `grpcio-tools`/`protoc`), kein committeter Stub im SDK-Baum —
      dasselbe Muster wie bei C#/Kotlin, hier zum ersten Mal für Python
      eingeführt.
- [ ] **Realserver-Integrationstest-Werkzeug eingeführt** (`ADR-0110`
      §Entscheidung Festlegung 2/Folgepflicht 1): ein neues Skript
      `tools/harness/run-sdk-python-integration-tests.sh` und ein neues
      Make-Target `make test-sdk-python-integration` (Arbeitsnamen,
      `harness/mk/sdk.mk`) — bringt die bestehende `compose.yaml`-Umgebung
      hoch (Muster `tools/harness/run-integration-tests.sh`s
      Bring-up-Phase: PostgreSQL/NATS/Feed-Container, Schema-Rollout,
      Tabellen-Aktivierung), baut eine neue, additive Docker-Stufe
      `integration` in `sdks/python/Dockerfile` (baut auf `build` auf,
      installiert das SDK samt `grpcio` und `pytest`-Integrationstests
      unter einem eigenen Verzeichnis, z. B.
      `pgchangefeed/tests/integration/`), startet diesen Container **im
      selben Docker-Netz** wie der Feed-Container (`docker run --network
      cdc-feed-test`, dasselbe benannte Netz wie
      `run-integration-tests.sh`s Wegwerf-Go-Clients) und lässt ihn real
      per gRPC eine zuvor über `psql` eingefügte Änderung empfangen —
      belegt am Nachrichtenschema und über die `change_id` gegen
      `cdc.changes` gehalten (Muster: `run-integration-tests.sh`s
      gRPC-Rundlauf, hier mit dem SDK selbst statt einem Wegwerf-Client als
      Prüfling). Kein Gate — braucht DB-Zugang/Docker/Netz, dieselbe
      Klasse wie `make test-integration`.
- [ ] Kein Import aus `internal/**`/`cmd/**`/`gen/**` dieses Repos — real
      geprüft (`grep -rn "internal/\|cmd/\|gen/" sdks/python/`).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update (`docs/user/benutzerhandbuch.md`,
      `spec/pflichtenheft.md`): bewusst **nicht** in diesem Slice —
      gebündelt im letzten Flächen-Slice (§1 Abgrenzung); `harness/README.md`
      §Sensors/§Werkzeuge bekommt die neue `make test-sdk-python-integration`-
      Zeile **in diesem** Slice, weil das Werkzeug hier real entsteht
      (`AGENTS.md` §4: kein Träger nennt ein Target, das es nicht gibt —
      umgekehrt gilt auch: ein real existierendes Target bekommt seine
      Zeile im selben Zug, nicht erst später).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben,
      **falls dieser Slice einen Inventur-Fund auflöst** — entfällt: keine
      Reconciliation-Datei in diesem Repo.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder Beleg in `evidence/`; keine Beobachtung angefallen
      ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      dieser Slice gehört zu
      [welle-sdk-python-vollabdeckung](../welle-sdk-python-vollabdeckung.md)
      (noch offen); die Prüfung läuft regelkonform bei deren Closure.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `sdks/python/pgchangefeed/src/pgchangefeed/grpc_client.py` (Arbeitsname) | neu | öffentliche API-Fläche für `StreamChanges` (`SPEC-020`). |
| `sdks/python/pgchangefeed/pyproject.toml` | update | `grpcio`-Laufzeitabhängigkeit ergänzen (analog `Grpc.Net.Client` bei C#), `grpcio-tools` als Test-/Build-Extra für die Stub-Erzeugung. |
| `sdks/python/Dockerfile` | update | zusätzlicher `proto`-Bau-Kontext (`COPY --from=proto …`), `protoc`-Aufruf über `grpcio-tools`, neue Stufe `integration`. |
| `sdks/python/pgchangefeed/tests/test_grpc_client.py` (Arbeitsname) | neu | Nachrichtenschema-Vollständigkeit, Authn-Boundary — netzlos, gegen einen Fake-Stub. |
| `sdks/python/pgchangefeed/tests/integration/test_grpc_realserver.py` (Arbeitsname) | neu | Realserver-Rundlauf: SDK empfängt real eine committete Änderung über gRPC. |
| `tools/harness/run-sdk-python-integration-tests.sh` | neu | Bring-up der `compose.yaml`-Umgebung, Bau/Start der `integration`-Docker-Stufe, Abbau. |
| `harness/mk/sdk.mk` | update | neues Target `test-sdk-python-integration`. |
| `harness/README.md` §Werkzeuge | update | neue Zeile für `make test-sdk-python-integration`. |

**Ansatz:** Draht-Kenntnis-Quelle für das Nachrichtenschema ist
[`SPEC-020`](../../../../spec/pflichtenheft.md) direkt und
`internal/adapters/driving/grpc/server_test.go` (nur **gelesen**, nicht
importiert — kein Python-Import eines privaten Baums dieses Repos,
`ADR-0110` §Kontext nennt genau diese Klasse Referenz als „schwächere,
aber reale Mitigation"). Für die Draht-**Semantik** (Stream-Lebenszyklus,
Metadata-Header-Form) zusätzlich `examples/csharp/grpc-client`/
`examples/kotlin/grpc-client` als funktionierende Fremdsprachen-Referenzen
gegen denselben Server (`ADR-0110` §Kontext „was sich geändert hat").
Für das Integrationstest-Werkzeug ist `tools/harness/run-integration-tests.sh`
das Struktur-Vorbild für Bring-up/Abbau der Compose-Umgebung — dieses
Skript bleibt aber **unverändert**; das neue Skript ist eigenständig,
kein Patch daran (Begründung: unterschiedliche Zuständigkeit — Server-E2E
vs. SDK-Protokoll-Beleg, unterschiedlicher Toolchain-Container Go vs.
Python).

**Architektur-Entscheidung für das Integrationstest-Werkzeug (zu klären,
hier entschieden, kein Vorgriff auf eine neue ADR — `ADR-0110` delegiert
den Mechanismus ausdrücklich an den umsetzenden Zug, §Konsequenzen
Folgepflicht 1):**

- **Eigenständiges Skript, geteilte Compose-Umgebung.** Kein neuer,
  separater Compose-Stack — `compose.yaml` bringt bereits
  PostgreSQL/NATS/Feed-Container mit allen benötigten Adressen
  (`CDC_HTTP_ADDR`, `CDC_GRPC_ADDR`, `CDC_NATS_URL`) hoch; ein zweiter,
  paralleler Compose-Stack würde dieselbe ~150-zeilige Bring-up-Logik
  duplizieren, ohne einen Vorteil zu bieten.
- **Kein Wegwerf-Client wie `tools/harness/grpcclient`.** Anders als die
  bestehenden Go-Wegwerf-Clients (die kein Produktionscode sind) **ist**
  hier das SDK selbst der Prüfling — der Integrationstest importiert
  `pgchangefeed.grpc_client` direkt, kein Duplikat-Client daneben.
- **Eigenes Skript statt Erweiterung von `run-integration-tests.sh`.**
  Das bestehende Skript ist Server-E2E-Eigentum (`LH-QA-POR-003`), fährt
  einen Go-Toolchain-Container und ist bereits > 2900 Zeilen; ein
  Python-Toolchain-Aufruf mittendrin würde beide Zuständigkeiten
  vermischen. Das neue Skript bleibt schlanker: nur Bring-up, ein
  Docker-Bau/-Lauf der `integration`-Stufe, Abbau — keine der
  Rollen-DSN-/Retention-/Diagnose-Prüfungen, die `run-integration-tests.sh`
  sonst trägt.

## 4. Trigger

**Start** (`next` → `in-progress`): wenn diese Welle eröffnet ist und kein
anderer Slice in `in-progress/` liegt (WIP-Limit 1).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls das
  Integrationstest-Werkzeug selbst (Bring-up/Abbau/Docker-Stufe) mehr
  Aufwand verlangt als die gRPC-Fläche — dann Abspaltung eines eigenen
  „Integrationstest-Werkzeug"-Slice, den die gRPC-Fläche als Vorgänger
  bräuchte.
- `in-progress` → `open` (blockiert — Carveout?): `--build-context
  proto=proto` lässt sich aus einem noch unbekannten Grund nicht auf den
  Python-Baum übertragen (unwahrscheinlich — C#/Kotlin belegen den
  Mechanismus bereits real, `grpcio-tools`/`protoc` ist eine etablierte
  Python-Toolchain).

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + ein realer, grüner
`make test-sdk-python-integration`-Lauf gegen die gRPC-Fläche +
Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- Das neue Integrationstest-Werkzeug könnte real langsamer/flakier sein
  als erwartet (Compose-Bring-up + Docker-Bau + realer Netzwerk-Rundlauf
  in einem einzigen Skript). **Ausgang:** weiter offen — reale Laufzeit
  und Stabilität werden erst beim tatsächlichen Lauf sichtbar; ein
  wiederholtes Auftreten von Flakiness wäre ein Kandidat für
  `BEO-PGC/test-integration-retention-timing-flake`-artige Beobachtung
  (eigener, neuer Registereintrag, falls es real auftritt).
- `grpcio-tools`s `protoc`-Python-Plugin könnte einen anderen
  Stub-Codestil erzeugen als erwartet (z. B. keine Typ-Stubs `.pyi` ohne
  Zusatz-Flag). **Ausgang:** weiter offen bis zum ersten realen Bau —
  kein Blocker, falls die erzeugten Typen funktional korrekt sind, auch
  ohne `.pyi`.
- Das neue Skript teilt sich das benannte Docker-Netz `cdc-feed-test` mit
  einem eventuell noch laufenden `make test-integration`-Lauf, falls
  beide gleichzeitig ausgeführt werden. **Ausgang:** entfallen — dieselbe
  Grenze gilt bereits für `run-integration-tests.sh` selbst (kein
  paralleler zweiter Lauf vorgesehen); das neue Skript folgt derselben
  Erwartung, kein neues Risiko.

## 7. Closure-Notiz

- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Steering-Loop-Eintrag:** <…>
- **Beobachtungs-Register (`../observations/`):** <…>
- **Folge-Slices:** <…>
- **Risiken aus §6:** <…>
- **Drei Paarungen:** dieser Slice gehört zu
  [welle-sdk-python-vollabdeckung](../welle-sdk-python-vollabdeckung.md)
  (noch offen) — die Prüfung läuft regelkonform bei deren Closure.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Sub-Area `sdks/python/` — bereits
mit `slice-sdk-python-projektgeruest` eröffnet (GF), keine erneute
Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert, Suchlauf-Pflicht
trägt insbesondere `options.py`s Docstring),
`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`
(verkörpert, DoD-Punkt „Doku-Update" liegt gebündelt im Folge-Slice),
`BEO-PGC/test-runner-stiller-ausschluss` (offen, 2×, nicht einschlägig —
anderes Skript), `BEO-PGC/deutsches-fachwort-im-englischen-sdk-readme`
(offen, 3×, Architect-Entscheidung, kein Handlungsbedarf hier).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Fortsetzung von
`slice-sdk-python-projektgeruest`).
