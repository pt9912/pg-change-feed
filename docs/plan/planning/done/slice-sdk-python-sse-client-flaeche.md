# Slice sdk-python-sse-client-flaeche: Öffentliche SSE-Stream-Client-Fläche (`SPEC-021`), Realserver-Beleg

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** [welle-sdk-python-vollabdeckung](welle-sdk-python-vollabdeckung.md).

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
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      *(Haupt-Review F-1…F-5; Fixrunde 1 löste F-1…F-5; Re-Review FR-1
      (MEDIUM) gelöst in Fixrunde 2; Re-Review-Endstand: kein offenes
      HIGH/MEDIUM.)*
- [ ] Doku-Update (`docs/user/benutzerhandbuch.md`,
      `spec/pflichtenheft.md`): bewusst **nicht** in diesem Slice —
      gebündelt im letzten Flächen-Slice (§1 Abgrenzung). *(Präzisiert nach
      Review F-1 (§3.13-Suchlauf-Feld oben trägt den vollen Stand): der
      Python-`**SDK:**`-Absatz im SSE-Handbuch-Abschnitt gehört in diesen
      Zug (Version 1.42, Fixrunde); gebündelt bleiben nur der
      NATS-Handbuch-Teil und `spec/pflichtenheft.md` — derselbe Schnitt
      wie beim Vorgänger-Slice.)*
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag. *(§7 unten —
      Lerneintrag: §3.13-Suchlauf je bewegter Eigenschaft, nicht je
      Slice; Beleg im Beobachtungs-Register.)*
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben,
      **falls dieser Slice einen Inventur-Fund auflöst** — entfällt: keine
      Reconciliation-Datei in diesem Repo.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder Beleg in `evidence/`; keine Beobachtung angefallen
      ist ebenfalls eine Antwort und wird in §7 notiert. *(Beleg:
      `BEO-PGC/arbeit-ueberholt-stehenden-traeger`,
      `evidence/slice-sdk-python-sse-client-flaeche.md`; übrige
      Kandidaten geprüft, siehe §7.)*
- [x] Jedes Risiko aus §6 trägt einen Ausgang. *(zwei Ausgänge in §6,
      final — Risiko 1 entfallen mit Beleg am Ort, Risiko 2 entfallen;
      siehe §7.)*
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      dieser Slice gehört zu
      [welle-sdk-python-vollabdeckung](welle-sdk-python-vollabdeckung.md)
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
| `sdks/python/Dockerfile`: CMD-Guard je Phase (Shell-Form mit `:?`) | update | Review F-3: der blanke docker run-Aufruf fuhr still die Unit-Tests grün (`testpaths = tests`), ohne jeden Integrationstest — verwandter Fall der Klasse `BEO-PGC/test-runner-stiller-ausschluss`. Real jetzt: `CMD` in Shell-Form mit `${PGCHANGEFEED_TEST_FILE:?…}`-Guard — blanker Aufruf scheitert laut, der Runner reicht die Testdatei als Umgebungsvariable nach. |
| `sdks/python/pgchangefeed/src/pgchangefeed/exceptions.py` (Docstring) | update | Review F-4: `PgChangeFeedMalformedResponseError` nennt jetzt auch die `SPEC-021`-Verletzungsform (der SSE-Client wirft ihn für Stream-Schema-Verletzungen). |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „SSE bleibt außerhalb des Packages"; beide Stände gemessen: Parent `b43251ee` und HEAD):**

| Träger | Befund | Behandlung |
|---|---|---|
| `sdks/python/pgchangefeed/src/pgchangefeed/__init__.py` (Docstring) | Satzform gefunden, gezogen in diesem Diff | erledigt |
| `sdks/python/pgchangefeed/src/pgchangefeed/options.py` (Docstring) | Satzform gefunden, gezogen | erledigt |
| `sdks/python/README.md` §Status | Satzform gefunden, gezogen | erledigt |
| `docs/user/benutzerhandbuch.md` §SSE (Python-`**SDK:**`-Absatz) | Absatz fehlte (Review F-1) | Fixrunde: Absatz ergänzt (Version 1.42 + Historie) |
| `harness/README.md` §Werkzeuge (`make test-sdk-python-integration`-Zeile) | Zeile trug die Ein-Flächen-Form (Review F-2) | Fixrunde: Zwei-Flächen-Form gezogen |
| `spec/pflichtenheft.md` `LH-FA-SST-009.a`/`SPEC-027` | Satzform „deckt HTTP-API" steht nicht im Diff — Bündelung im letzten Flächen-Slice | bleibt gebündelt (konsistent C#/Kotlin-Präzedenz, [`ADR-0110`](../../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md) Folgepflicht 2) |
| Wurzel-README | „deckt HTTP-API"-Form für das veröffentlichte 0.1.0-Package — wahr für den veröffentlichen Stand | bleibt (bündelige Praxis wie beim Vorgänger-Slice) |
| `docs/user/version.md` | Versionsstand des Python-Packages — Hebung im letzten Flächen-Slice | bleibt gebündelt |

**§3.13-Suchlauf, zweite Eigenschaft (Fixrunde FR-1 — bewegte Eigenschaft: „die Testdatei wird je Phase als docker run-Argument übergeben"; beide Stände gemessen: Vorher-Stand `3c941b0b` — dort kam die Phrase in den `harness/README`-Träger — und HEAD):**

| Träger | Befund | Behandlung |
|---|---|---|
| `tools/harness/run-sdk-python-integration-tests.sh` Kopf („Stufe-ENTRYPOINT") | Phrase gefunden, durch die CMD-Umstellung falsch | Fixrunde 2: „Stufe-CMD" gezogen |
| `tools/harness/run-sdk-python-integration-tests.sh` Kopf („docker run-Argument") | Phrase gefunden, gezogen | Fixrunde 2: Umgebungsvariable `PGCHANGEFEED_TEST_FILE`-Form |
| `harness/README.md` §Werkzeuge („docker run-Argument") | Phrase gefunden, gezogen | Fixrunde 2: ENV-Form samt Guard-Beschreibung |

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
  **Ausgang:** entfallen — der zeilenweise Parser bindet Chunk-Grenzen
  nicht: die Unit-Probe fährt vier Chunk-Formen (CRLF, CR/NL-Split über
  Chunk-Grenzen, Mid-Line-Split, LF) netzlos grün im Integration-Image
  (Haupt-Review F-6-Negativprobe, Unit-Test
  `test_chunk_boundaries_need_not_coincide_with_lines`), und die reale
  `make test-sdk-python-integration`-Fläche läuft grün gegen den
  laufenden Feed-Container (Implementer-Lauf EXIT=0; Verifikation §1
  DoD 2: eigener Verifier-Lauf EXIT=0). Ein Pufferungs-Abweichfall tritt
  in keinem der Läufe auf.
- Der Realserver-SSE-Test könnte real langsamer terminieren als der
  gRPC-Test (Polling auf ein Event statt eines blockierenden Streams).
  **Ausgang:** entfallen — dasselbe Warte-/Poll-Muster wie
  `run-integration-tests.sh`s bestehender SSE-Rundlauf
  (`tools/harness/sseclient`) gilt unverändert.

## 7. Closure-Notiz

- **Was hat funktioniert:** der Plan-Schnitt trägt sich in einem Zug —
  die SSE-Fläche und die Runner-Erweiterung je Fläche gehören in einen
  Slice, weil die Fläche ohne den Werkzeug-Lauf nicht abnahmefähig ist
  ([`ADR-0110`](../../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
  §Entscheidung Festlegung 2/Folgepflicht 1). Beide Rückführungen aus §4
  stehen nicht an: der Vorgänger-Slice liegt in `done/`, und das Werkzeug
  nahm die zweite Testdatei ohne zweite Infrastruktur auf —
  `make test-sdk-python-integration` lief real grün (Implementer-Lauf
  EXIT=0: gRPC change_id=804-1, SSE change_id=807-1, je SQL-Gegenprüfung
  gegen `cdc.changes`; Verifier-Lauf EXIT=0: gRPC 804-1 deckend, SSE
  808-1 — laufgebundene Messgröße, Verifikation N-1). Die Rollen-Kette
  läuft unabhängig: Haupt-Review F-1…F-9, Fixrunde 1 löste F-1…F-5
  (`3c941b0b`), Re-Review FR-1 gelöst in Fixrunde 2 (`beeddc3c`), Endstand
  kein offenes HIGH/MEDIUM (Addendum im Fixrunden-Report); die
  Verifikation fährt eigene Läufe (Unit-Suite netzlos, beide Mutationen
  rot aus dem richtigen Grund, `make gates` Exit 0) und ihr Urteil
  „erfüllt mit Auflagen“ ist mit den beiden Auflagen V-1/V-2 gelöst
  (`25fad0e1`).
- **Was ging anders als geplant:** die Lücken-Struktur wiederholt sich
  in der Korrektur-Kette desselben Slice, nicht in der Umsetzung — der
  §3.13-Suchlauf war je Slice und je **einer** bewegten Eigenschaft
  gebunden („SSE bleibt außerhalb des Packages“), während der Fix-Zug
  eine **zweite** Eigenschaft bewegte (Testdatei-Übergabe:
  docker run-Argument → Umgebungsvariable) ohne eigenen Suchlauf.
  Re-Review FR-1 fand die drei Phrase-Stellen, die ein zweiter Suchlauf
  gemeldet hätte; Haupt-Review F-2 traf denselben Fall auf der ersten
  Eigenschaft (die `harness/README.md`-Werkzeug-Zeile stand nicht im
  Plan-Nachzug, gefunden vom Reviewer). Verifikation V-1 trägt die
  dritte Stelle: die Stand-Deklaration der zweiten Suchlauf-Tabelle
  nannte `6bbe99d9` als Vorher-Stand des `harness/README.md`-Trägers —
  der wahre Vorher-Stand ist `3c941b0b` (die Phrase kam erst mit
  Fixrunde 1 in den Träger); gelöst im Verifikations-Commit `25fad0e1`.
- **Steering-Loop-Eintrag:** geschärfte Regel (Anwendungs-Schärfung der
  bereits verkörperten Klasse
  `BEO-PGC/arbeit-ueberholt-stehenden-traeger`, Anker
  [`AGENTS.md`](../../../../AGENTS.md) §3.13 · seit welle-20): der
  §3.13-Suchlauf wird je **bewegter Eigenschaft** geführt, nicht je
  Slice — jede mechanische Umstellung im Fix-Diff (hier:
  docker run-Argument → Umgebungsvariable `PGCHANGEFEED_TEST_FILE`) ist
  eine eigene Eigenschaft mit eigenem Suchlauf über beide Stände. Kein
  neuer Sensor: die verfügbare Falsifikation bleibt die Messung an
  beiden Ständen (§3.13s Grenze für Zahlen/Prosa-Umformulierungen
  bleibt; die Mechanik-Phrase ist ein neuer Fall derselben Grenze — ein
  Träger, den ein Suchlauf mit einem Wortmuster je Eigenschaft träfe,
  bräuchte eine Semantik-Entscheidung, welche Umstellung zur welchen
  Eigenschaft gehört). Keine benannte Spec-Lücke. Ob die Schärfung einen
  eigenen Satz in `AGENTS.md` §3.13 trägt, prüft der Lese-Schritt der
  Welle-Closure — die Klasse steht mit diesem Beleg bei 22×
  (Datei-Anzahl unter `evidence/`, real ausgezählt), bereits verkörpert.
- **Beobachtungs-Register (`../observations/`):**
  - **`BEO-PGC/arbeit-ueberholt-stehenden-traeger`** — neuer,
    zweiundzwanzigster Beleg:
    `evidence/slice-sdk-python-sse-client-flaeche.md`; die Kette
    (Haupt-Review F-2, Re-Review FR-1, Verifikation V-1) trägt je ihre
    Fundstellen; `state.md` trägt den abgeleiteten Zähler.
  - **`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`** — geprüft:
    kein eigenständiger Fund in dieser Kette (die V-1-Abweichung ist
    eine Stand-Deklaration des Suchlaufs, kein Zahlenwert, der gegen
    seine Messung driftet). Die Zahlen des Closure-Zugs sind je aus
    ihrer Messung gezogen: Zähler 22× (Datei-Anzahl unter `evidence/`,
    real ausgezählt — 21× vor diesem Beleg), Testzahlen 20/3/8/10 = 41
    (nachgemessen, Verifikation §1 DoD 1), `change_id`-Werte je
    laufgebunden (N-1).
  - **`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`**
    — geprüft: die Klasse traf den Verbund (Haupt-Review F-1, MEDIUM)
    und ist in der Kette gelöst (Fixrunde 1: Python-`**SDK:**`-Absatz im
    SSE-Handbuch-Abschnitt, Version 1.42, Historie-Zeile; Verifikation
    §1 DoD 5 misst alle drei Träger); der gebündelte NATS-Handbuch-Teil
    und `spec/pflichtenheft.md` bleiben beim Folgelauf.
  - **`BEO-PGC/test-runner-stiller-ausschluss`** — geprüft: Haupt-Review
    F-3 ist verwandt (stiller Ausschluss im Default-Aufruf), aber kein
    Vorkommen der Klasse — anderes Skript und anderer Mechanismus
    (pytest-Sammlung über `testpaths` beim blanken Image-Aufruf statt
    `go test -run`-Filterung in `tools/harness/run-integration-tests.sh`);
    der Fall ist im Zug gelöst (CMD-Guard, Negativprobe Ausgang 2,
    Verifikation §1 DoD 1) und in den Reports konserviert — kein
    Registereintrag, der Zähler der Klasse bleibt 2×.
  - Re-Review FR-2 (Link-Nachzug `228dfc9e` an committeten
    Review-Reports, INFO) — Prozess-Merkung ohne Registereintrag; der
    Fund ist in der Kette gelöst und im Report konserviert.
- **Folge-Slices:** `slice-sdk-python-nats-stream-client-flaeche`
  (Reihenfolge Welle-Plan §4) — erweitert das hier erweiterte Werkzeug
  um die dritte Fläche und bündelt den NATS-Handbuch-Teil, den
  Version-Bump und den `spec/pflichtenheft.md`-Nachzug. Kein
  zusätzlicher Slice aus dieser Closure.
- **Risiken aus §6:** Risiko 1 (httpx-Pufferung) — **entfallen** (der
  Parser bindet Chunk-Grenzen nicht: vier Chunk-Formen netzlos grün,
  reale Integrationstest-Läufe grün; Begründung und Beleg am Ort in
  §6). Risiko 2 (SSE-Test-Laufzeit) — **entfallen** (unverändert; §6).
  Kein Ausgang „weiter offen“ ins Register — ein Registereintrag für ein
  nie aufgetretenes Muster trüge kein `evidence/` (Paarung (c)).
- **Drei Paarungen:** dieser Slice gehört zu
  [welle-sdk-python-vollabdeckung](welle-sdk-python-vollabdeckung.md)
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
