# Slice sdk-python-sse-client-flaeche: Öffentliche SSE-Stream-Client-Fläche (`SPEC-021`), Realserver-Beleg

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** [welle-sdk-python-vollabdeckung](../welle-sdk-python-vollabdeckung.md).

**Bezug:** [`LH-FA-SST-009`](../../../../spec/lastenheft.md),
[`LH-FA-SST-008`](../../../../spec/lastenheft.md) (Live-Change-Stream —
Boundary: keine Zustellgarantie, kein Stream-internes Replay),
[`ADR-0110`](../../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
§Entscheidung Festlegung 1/2 (Umfang, verschärfte Test-Pflicht).

**Berührte Spec-Stellen:** [`SPEC-021`](../../../../spec/pflichtenheft.md)
(Endpunkt, Event-Form, Nachrichtenschema).

**Verantwortlich:** Implementer-Agent, 2026-09-23.

**Autor:** Planner-Agent, direkt beauftragt (`ADR-0110` §Konsequenzen
Folgepflicht 1). **Datum:** 2026-09-21.

---

## 1. Ziel und Abgrenzung

**Ziel:** Eine öffentliche, stabile Python-API-Fläche im bestehenden
Package `pgchangefeed` für den SSE-Endpunkt von
[`SPEC-021`](../../../../spec/pflichtenheft.md) — `GET /changes/stream`
über `httpx` (bereits die einzige Fremdabhängigkeit dieses Packages,
Streaming-Response statt Einzel-Request), Bearer-Token-Auth, ein
Frame-Parser für `event:`/`data:`-Zeilen, der die zehn Nachrichtenfelder
liefert. **Zusätzlich** erweitert dieser Slice das in
`slice-sdk-python-grpc-client-flaeche` eingeführte
Realserver-Integrationstest-Werkzeug um einen SSE-Rundlauf (`ADR-0110`
§Entscheidung Festlegung 2/Folgepflicht 1 — die verschärfte Test-Pflicht
gilt für **jede** neue Fläche einzeln).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **HTTP-API- oder gRPC-Fläche** — HTTP bereits geliefert; gRPC ist der
  direkte Vorgänger-Slice dieser Welle.
- **NATS-Vollinhalts-Fläche** — eigener Folge-Slice
  (`slice-sdk-python-nats-stream-client-flaeche`), getrennter Draht-Vertrag
  (NATS statt HTTP), eigene Fremdabhängigkeit.
- **Ein `examples/python/`-Referenz-Client** — `ADR-0110` §Entscheidung
  Festlegung 2 mandatiert das ausdrücklich nicht.
- **Stream-internes Replay oder `Last-Event-ID`-Auswertung** — `SPEC-021`
  trägt beides nicht.
- **Version-Bump, `spec/pflichtenheft.md`-Träger-Nachzug,
  `docs/user/benutzerhandbuch.md`-Nachzug** — gebündelt im letzten
  Flächen-Slice dieser Welle (`slice-sdk-python-nats-stream-client-flaeche`).
- **Ein Neubau des Integrationstest-Werkzeugs** — dieser Slice erweitert
  das bereits bestehende Skript/Docker-Stufe um eine weitere Testdatei,
  baut keine zweite Infrastruktur auf.

## 2. Definition of Done

- [x] `LH-FA-SST-009` erfüllt: eine öffentliche, im Package sichtbare
      Client-Klasse öffnet `GET /changes/stream` (Bearer-Token in
      `Authorization`-Header über `httpx`), zerlegt SSE-Frames zu Events
      und liefert die zehn Nachrichtenfelder von
      [`SPEC-021`](../../../../spec/pflichtenheft.md) — Unit-Tests
      referenzieren `SPEC-021` (Frame-Parser netzlos, Authn-Boundary gegen
      einen `httpx.MockTransport`, dasselbe Fake-Muster wie die
      bestehende HTTP-Fläche dieses Packages). *(Sensor-Beleg: `make
      sdk-pack-python` EXIT=0 — 41 Unit-Tests grün, darunter 10 neue
      SSE-Tests `tests/test_sse_client.py` (`grep -c "^def test_"` = 10;
      20 + 3 + 8 + 10 = 41; Zahlen aus derselben Messung wie die
      Gesamtzahl gezogen).)*
- [x] **Realserver-Integrationstest erweitert** (`ADR-0110` §Entscheidung
      Festlegung 2/Folgepflicht 1): eine neue Testdatei läuft über
      `make test-sdk-python-integration` und belegt real, dass die SDK
      SSE-Fläche eine zuvor über `psql` eingefügte Änderung über den
      laufenden Feed-Container empfängt — belegt am Nachrichtenschema und
      über die `change_id` gegen `cdc.changes` gehalten. *(Sensor-Beleg:
      `make test-sdk-python-integration` EXIT=0 — der gRPC-Rundlauf
      (change_id=804-1) und der SSE-Rundlauf (change_id=807-1) liefen
      real grün, je mit SQL-Gegenprüfung gegen `cdc.changes`; der
      Öffnungsversuch ohne Token endete je mit gRPC-Status
      `Unauthenticated` bzw. HTTP-Status 401.)*
- [x] Kein Import aus `internal/**`/`cmd/**`/`gen/**` dieses Repos.
      *(strenge Import-Zeilen-Prüfung über `sdks/python/pgchangefeed/`:
      kein Treffer gegen internal/cmd/gen.)*
- [x] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update (`docs/user/benutzerhandbuch.md`,
      `spec/pflichtenheft.md`): bewusst **nicht** in diesem Slice —
      gebündelt im letzten Flächen-Slice (§1 Abgrenzung). *(Wie im
      Vorgänger-Slice: der gRPC-/SSE-Teil des Handbuchs wurde in der
      Fixrunde des Vorgänger-Slices bzw. im selben Muster je Fläche
      nachgezogen — der SSE-`**SDK:**`-Absatz gehört in denselben Zug wie
      die Fläche; der NATS-Handbuch-Teil und `spec/pflichtenheft.md`
      bleiben beim letzten Flächen-Slice.)*
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
| `sdks/python/pgchangefeed/src/pgchangefeed/sse_client.py` (Arbeitsname) | neu | öffentliche API-Fläche für den SSE-Stream (`SPEC-021`). |
| `sdks/python/pgchangefeed/tests/test_sse_client.py` (Arbeitsname) | neu | Frame-Parser-Grenzfälle, Authn-Boundary, Nachrichtenschema-Vollständigkeit — netzlos gegen `httpx.MockTransport`. |
| `sdks/python/pgchangefeed/tests/integration/test_sse_realserver.py` (Arbeitsname) | neu | Realserver-Rundlauf: SDK empfängt real eine committete Änderung über SSE. |
| `sdks/python/Dockerfile` | update | `integration`-Stufe kopiert/installiert zusätzlich `tests/integration/test_sse_realserver.py` (kein neuer Bau-Mechanismus, dieselbe Stufe). |

**Plan-Nachzug (im selben Lauf, vor dem Gate-Lauf):**

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `sdks/python/pgchangefeed/integration/test_sse_realserver.py` (statt `tests/integration/`) | Abweichung | der Vorgänger-Slice zog die Realserver-Tests in den Geschwister-Ordner `integration/` (Review-Begründung dort); die §3-Zeile trug noch den Arbeitsnamen. |
| `sdks/python/Dockerfile` (ENTRYPOINT flexibel statt kopiertem Testfile) | Abweichung | die `integration`-Stufe kopiert das Verzeichnis bereits **ganz** (`COPY pgchangefeed/integration pgchangefeed/integration`) — keine zusätzliche COPY-Zeile nötig. Real gezogen: der ENTRYPOINT von der gebundenen Fassung auf `python -u -m pytest` erweitert, der Runner reicht die Testdatei je Phase als docker run-Argument nach — jede Phase nennt ihre Testdatei explizit, kein stiller Ausschluss des Rests. |
| `sdks/python/pgchangefeed/src/pgchangefeed/models.py`: `StreamChange`-Datenklasse (10 Felder) | neu | der SSE-Wire-Vertrag trägt die zehn Domain-Felder, nicht die zwölf des HTTP-Lesezugriffs ([`SPEC-022`](../../../../spec/pflichtenheft.md)) — eigene getypte Klasse statt Wiederverwendung der HTTP-`Change`, dieselbe Mapping-Doktrin wie der Rest des Packages. |
| `tools/harness/run-sdk-python-integration-tests.sh` | update | der Runner bekommt eine `run_surface_phase`-Funktion: je Fläche ein Aufruf mit eigener Testdatei (explizit als docker run-Argument), eigenem Sentinel und ID-Wertebereich (gRPC 300ff., SSE 310ff.) und eigener Reject-Marker-Form (`Unauthenticated` bzw. `401`); die SQL-Gegenprüfung gegen `cdc.changes` läuft je Phase. |
| `sdks/python/pgchangefeed/src/pgchangefeed/__init__.py`-Docstring, `options.py`-Docstring, `sdks/python/README.md` §Status | update | Träger-Nachzug (`AGENTS.md` §3.13): die Satzform „SSE bleibt außerhalb" wird durch diesen Slice falsch. |

**Ansatz:** Referenzmaterial für die Frame-Zerlegung ist
`examples/csharp/sse-client/SseStream.cs`/`examples/kotlin/sse-client/…/SseStream.kt`
(Fremdsprachen-Referenzen gegen denselben Draht, `ADR-0110` §Kontext „was
sich geändert hat") — kein `examples/python/`-Vorbild. Draht-Vertrag
direkt aus [`SPEC-021`](../../../../spec/pflichtenheft.md); die
`internal/adapters/driving/http/sse_test.go`-Feldvollständigkeits-Prüfung
ist die serverseitige Gegenprobe (nur gelesen, kein Import).

## 4. Trigger

**Start** (`next` → `in-progress`): wenn
`slice-sdk-python-grpc-client-flaeche` in `done/` liegt (Welle-Plan §4
Reihenfolge — strikt sequentiell, das Integrationstest-Werkzeug muss
bereits existieren).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): entfällt aus
  heutiger Sicht — ein Endpunkt, ein Nachrichtenschema, bereits einmal in
  diesem Package (HTTP) und zweimal in Fremdsprachen (C#/Kotlin) erprobtes
  Muster; das Integrationstest-Werkzeug existiert bereits.
- `in-progress` → `open` (blockiert — Carveout?): `slice-sdk-python-grpc-client-flaeche`
  liegt noch nicht in `done/`, oder das Integrationstest-Werkzeug lässt
  sich aus einem noch unbekannten Grund nicht um eine zweite Testdatei
  erweitern (unwahrscheinlich).

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + ein realer, grüner
`make test-sdk-python-integration`-Lauf gegen die SSE-Fläche +
Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- `httpx`s Streaming-Response-API könnte SSE-Frames anders puffern als
  erwartet (Chunk-Grenzen fallen nicht mit Zeilen-Grenzen zusammen).
  **Ausgang:** weiter offen bis zum ersten realen Integrationstest-Lauf —
  der Frame-Parser wird bewusst zeilenweise über einen Iterator gebaut
  (Muster `examples/csharp/sse-client/SseStream.cs`s `ReadEvent`), nicht
  über eine Annahme fester Chunk-Grenzen.
- Der Realserver-SSE-Test könnte real langsamer terminieren als der
  gRPC-Test (Polling auf ein Event statt eines blockierenden Streams).
  **Ausgang:** entfallen — dasselbe Warte-/Poll-Muster wie
  `run-integration-tests.sh`s bestehender SSE-Rundlauf
  (`tools/harness/sseclient`) gilt unverändert.

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
`BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert),
`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`
(verkörpert, Doku-Update gebündelt im Folge-Slice),
`BEO-PGC/test-runner-stiller-ausschluss` (offen, 2×, nicht einschlägig).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Fortsetzung von
`slice-sdk-python-projektgeruest`).
