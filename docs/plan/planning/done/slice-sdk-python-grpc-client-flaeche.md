# Slice sdk-python-grpc-client-flaeche: Öffentliche gRPC-Stream-Client-Fläche (`SPEC-020`), Realserver-Integrationstest-Werkzeug

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** [welle-sdk-python-vollabdeckung](welle-sdk-python-vollabdeckung.md).

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

**Verantwortlich:** Implementer-Agent, 2026-09-23.

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

- [x] `LH-FA-SST-009` erfüllt: eine öffentliche, im Package sichtbare
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
      eingeführt. *(Sensor-Beleg: `make sdk-pack-python` EXIT=0 —
      31 Unit-Tests grün, darunter 8 neue gRPC-Tests
      `tests/test_grpc_client.py` (`grep -c "^def test_"` = 8; 20 + 3 + 8
      = 31; Endstand nach der Fixrunde, die F-4 mit zwei
      Timeout-Bindungs-Tests aufgelöst hat); Wheel+SDist tragen die
      `grpc_gen`-Stub-Module.)*
- [x] **Realserver-Integrationstest-Werkzeug eingeführt** (`ADR-0110`
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
      Klasse wie `make test-integration`. *(Sensor-Beleg: `make
      test-sdk-python-integration` EXIT=0 — der SDK-Client empfing real
      eine committete Änderung über `pg-change-feed:9090`
      (change_id=804-1, unabhängig über `cdc.changes` lesbar); der
      Öffnungsversuch ohne Token endete mit gRPC-Status
      `Unauthenticated`.)*
- [x] Kein Import aus `internal/**`/`cmd/**`/`gen/**` dieses Repos — real
      geprüft (`grep -rn "internal/\|cmd/\|gen/" sdks/python/`). *(Der
      Wörter-Grep trägt jetzt Prosadiskurs über die Stub-Erzeugung
      (`grpc_gen`-Pfad, „`gen/**` bleibt die Go-Bindung") in Kommentaren;
      die strenge Import-Zeilen-Prüfung
      (`grep -E "^\s*(from|import)" … \| grep internal/cmd/gen`) liefert
      keinen Treffer.)*
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      *(Haupt-Review: review-slice-sdk-python-grpc-client-flaeche.md
      F-1…F-6; Fixrunde 1 löste F-1 bis F-4/F-6; Verifikation V-1…V-3
      gelöst; Fixrunden-Re-Review: review-…-fixrunde.md — 0 HIGH, MEDIUM
      F-1 in Fixrunde 2 gelöst; Fixrunde 2: kein offenes HIGH/MEDIUM.)*
- [x] Doku-Update (`docs/user/benutzerhandbuch.md`,
      `spec/pflichtenheft.md`): bewusst **nicht** in diesem Slice —
      gebündelt im letzten Flächen-Slice (§1 Abgrenzung); `harness/README.md`
      §Sensors/§Werkzeuge bekommt die neue `make test-sdk-python-integration`-
      Zeile **in diesem** Slice, weil das Werkzeug hier real entsteht
      (`AGENTS.md` §4: kein Träger nennt ein Target, das es nicht gibt —
      umgekehrt gilt auch: ein real existierendes Target bekommt seine
      Zeile im selben Zug, nicht erst später). *(In diesem Lauf erledigt:
      Zeile in `harness/README.md` §Werkzeuge, dazu Nachzug auf der
      bestehenden `make sdk-pack-python`-Zeile — der pack-Bau trägt jetzt
      denselben Bau-Kontext.)* *(Fixrunden-Nachzug, Re-Review F-2/F-3: der
      gRPC-Teil des Handbuchs ist hier trotzdem korrigiert — der
      PyPI-Absatz (1.40) und der `**SDK:**`-Absatz im gRPC-Abschnitt (1.41)
      tragen die neue Fläche; der Plan-Satz „Handbuch nicht in diesem
      Slice" bleibt nur für den SSE-/NATS-Handbuch-Teil und
      `spec/pflichtenheft.md` geltend.)* *(Closure: die in diesem Slice
      geltenden Teile tragen reale Belege — `harness/README.md`-Zeile im
      Implementations-Commit `f1ce9be4`, gRPC-Handbuch-Teil 1.40/1.41 in
      `b7a993fe`/`2eb18947` (Kopf 1.41, `Stand:` 2026-09-23); die
      Bündelungs-Formel bleibt für den SSE-/NATS-Handbuch-Teil und
      `spec/pflichtenheft.md` geltend.)*
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag. *(§7 unten —
      Lerneintrag: Zahl-/Versions-Nachzug aus derselben Messung, die sie
      belegt; Beleg im Beobachtungs-Register.)*
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben,
      **falls dieser Slice einen Inventur-Fund auflöst** — entfällt: keine
      Reconciliation-Datei in diesem Repo.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder Beleg in `evidence/`; keine Beobachtung angefallen
      ist ebenfalls eine Antwort und wird in §7 notiert. *(Beleg:
      `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`,
      `evidence/slice-sdk-python-grpc-client-flaeche.md`; übrige
      Kandidaten geprüft, siehe §7.)*
- [x] Jedes Risiko aus §6 trägt einen Ausgang. *(drei Ausgänge in §6,
      final — 3× entfallen mit Begründung am Ort; siehe §7.)*
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      dieser Slice gehört zu
      [welle-sdk-python-vollabdeckung](welle-sdk-python-vollabdeckung.md)
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

**Plan-Nachzug (im selben Lauf, vor dem Gate-Lauf — über den Plan
hinausgehende oder abweichende Änderungen):**

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `sdks/python/pgchangefeed/integration/test_grpc_realserver.py` (statt `tests/integration/`) | Abweichung | Der blanke `RUN pytest`-Lauf der `build`-Stufe sammelt mit `testpaths = ["tests"]` auch Unterordner von `tests/`; ein `--ignore=tests/integration` in `addopts` würde auch den expliziten Integrations-Aufruf ausschließen. Ein eigener Top-Level-Ordner `integration/` (Geschwister von `tests/`) trägt beide Läufe ohne Ausschluss-Tricks: `testpaths` begrenzt den Unit-Lauf auf `tests/`, die `integration`-Stufe ruft ihren Pfad explizit auf. |
| `sdks/python/pgchangefeed/pyproject.toml`: zusätzlich `protobuf>=6` als Laufzeitabhängigkeit | Erweiterung | Die zur Bauzeit erzeugten Stub-Module importieren die `protobuf`-Runtime; das Wheel muss sie als Abhängigkeit deklarieren, sonst scheitert der Import beim Consumer. `grpcio-tools` trägt sie nur für den Bau, nicht als Paket-Vertrag. |
| `sdks/python/pgchangefeed/src/pgchangefeed/grpc_gen/__init__.py` | neu | Package-Marker: macht aus dem zur Bauzeit erzeugten Stub-Ordner ein reguläres Unterpackage (Wheel-Einbindung über `setuptools packages.find`); die Stub-Module selbst bleiben uncommittet (`.gitignore`-Nachzug, dieselbe Tabelle). |
| `sdks/python/.gitignore` | update | die zwei Stub-Module (`changestream_pb2.py`/`changestream_pb2_grpc.py`, protoc leitet die Modulnamen aus dem `.proto`-**Dateinamen** ab — nicht aus dem Proto-Package) bleiben uncommittet, der Package-Marker bleibt committet. |
| `tools/harness/sdk-pack-python.sh` | update | der `pack-export`-Bau läuft über dieselbe `COPY --from=proto`-Zeile des Dockerfiles — ohne `--build-context proto=proto` bricht `make sdk-pack-python` ab; der Aufruf trägt den Kontext zwingend (Muster `tools/harness/sdk-pack-csharp.sh`). |
| `stream_changes(timeout=None)` | Erweiterung | der Integrationstest braucht eine fristbare Empfangsschleife — ein blockierendes `next()` ohne Call-Deadline hängt endlos, wenn der Server nichts sendet; der `timeout`-Parameter ist der grpcio-native Weg und bleibt optional (`None` = unverändert unbegrenzt). |
| `sdks/python/README.md` §Status, `src/pgchangefeed/__init__.py`-Docstring | update | Träger-Nachzug (`AGENTS.md` §3.13, dieselbe Klasse wie der im Welle-Plan genannte `options.py`-Docstring): die Satzform „gRPC bleibt außerhalb" wird durch diesen Slice falsch. |
| `docs/user/benutzerhandbuch.md` (SDK-Absatz + Version 1.40 + Historie) | update | Review F-3: der Handbuch-Absatz zum PyPI-Package trägt dieselbe falsch werdende Satzform mit dem superseded `ADR-0107`-Zitat; statt auf den gebündelten Nachzug im letzten Flächen-Slice zu warten, korrigiert die Fixrunde den gRPC-Teil im selben Slice (Präzedenz: die C#/Kotlin-Vollabdeckungs-Slices korrigierten den Handbuch-Hinweis je Fläche im selben Zug, Änderungshistorie 1.34/1.37). Der SSE-/NATS-Handbuch-Teil bleibt beim letzten Flächen-Slice. (V-2-Nachzug: die Historien-Zeile trug 1.40, nicht 1.39 — 1.39 war bereits von der Kotlin-Welle belegt.) |
| `tests/test_grpc_client.py`: `timeout`-Weiterleitungs-Tests (2 neu, `def test_`-Zahl jetzt 8) | update | Review F-4: der neue Parameter `stream_changes(timeout=…)` ist ohne Aufzeichnung im Fake ungebunden — eine Regression, die ihn still fallen lässt, blieb grün. Der Fake zeichnet den Timeout je Aufruf mit; zwei Tests binden Weiterleitung (30.0) und Default (None). |
| `sdks/python/pgchangefeed/src/pgchangefeed/options.py` | update | Verifikation V-3 (Deklarations-Nachzug, sachlich unverändert): Docstring + Attribut-Doku ziehen mit der neuen Fläche nach — [`ADR-0107`](../../adr/0107-python-pypi-zweites-sdk-package.md)s Erst-Scope-Satz wird durch diesen Slice falsch, `address`/`api_token` tragen jetzt die gRPC-Form mit (dieselbe Träger-Klasse wie README/`__init__`, im Welle-Plan als Suchlauf-Ziel benannt). |
| `integration/test_grpc_realserver.py`: Nichtleer-Asserts zurück | update | Re-Review F-1 (MEDIUM): die typbewusste Neuschreibung (F-2 des Haupt-Reviews) hatte für `transaction_id`, `source_table_id`, `schema_version` die korrekte Nichtleer-Bindung ersetzt — ein proto3-Wire-Image ohne diese Felder (`""`) blieb grün. Drei `!= ""`-Asserts zurück (neben `change_id`). |
| `docs/user/benutzerhandbuch.md`: Python-`**SDK:**`-Absatz im gRPC-Abschnitt (1.41 + `Stand:`) | update | Re-Review F-2/F-4 (LOW): der gRPC-Handbuch-Abschnitt trug `**SDK:**`-Absätze nur für C#/Kotlin, die 1.40-Historie nannte eine Klasse, die der Körper nirgends nennt; der dritte Sprach-Absatz zieht nach (Präzedenz Änderungshistorie 1.37), Version 1.41, `Stand:`-Datum mitgezogen. |
| `integration/test_grpc_realserver.py`: Feldvollständigkeits-Prüfung typbewusst | update | Review F-2: das frühere `!= ""` war für `old_image`/`new_image` (bytes) und `sequence` (int) vakuum (nie ungleich-rot). Real jetzt: je Feld ein Typ-Assert, dazu Inhalt (Operation, Alt-Bild-Leere am INSERT, Sentinel, Tabelle, Schema). |

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
  in einem einzigen Skript). **Ausgang:** entfallen — der erste reale
  Lauf des Werkzeugs ist grün mit vollem Behauptungspfad (Verifikation
  §1 DoD 2: `change_id=804-1` real empfangen und über `cdc.changes`
  gegengeprüft, fünf explizite Fehl-Ausgänge vor dem Erfolgssatz,
  `set -euo pipefail`); Flakiness tritt in keinem der Läufe auf. Ein
  Registereintrag für ein nie aufgetretenes Muster trägt kein `evidence/`
  (Paarung (c) verlangt nicht-leeren Beleg); ein reales späteres
  Auftreten von Flakiness bekommt seinen Registereintrag beim dann
  geschlossenen Vorgang (Kandidat:
  `BEO-PGC/test-integration-retention-timing-flake`-artige Beobachtung).
- `grpcio-tools`s `protoc`-Python-Plugin könnte einen anderen
  Stub-Codestil erzeugen als erwartet (z. B. keine Typ-Stubs `.pyi` ohne
  Zusatz-Flag). **Ausgang:** entfallen — der erste reale Bau erzeugt
  funktionale Stubs ohne Blocker: Modul-Namen folgen dem
  `.proto`-Dateinamen, kein `.pyi`; die Deskriptor-Bindung der
  Unit-Tests und die Wheel-/SDist-Einbindung tragen die Fläche
  (Verifikation §1 DoD 1).
- Das neue Skript teilt sich das benannte Docker-Netz `cdc-feed-test` mit
  einem eventuell noch laufenden `make test-integration`-Lauf, falls
  beide gleichzeitig ausgeführt werden. **Ausgang:** entfallen — dieselbe
  Grenze gilt bereits für `run-integration-tests.sh` selbst (kein
  paralleler zweiter Lauf vorgesehen); das neue Skript folgt derselben
  Erwartung, kein neues Risiko.

## 7. Closure-Notiz

- **Was hat funktioniert:** der Plan-Schnitt trägt sich in einem Zug — die
  gRPC-Fläche und das Realserver-Integrationstest-Werkzeug gehören in
  einen Slice, weil die Fläche ohne den Werkzeug-Lauf nicht abnahmefähig
  ist ([`ADR-0110`](../../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
  §Entscheidung Festlegung 2). Beide Rückführungen aus §4 stehen nicht an:
  der Bau-Kontext `proto=proto` trägt den Python-Baum (C#/Kotlin belegen
  den Mechanismus bereits), und das Werkzeug ist mit der Fläche in
  derselben Größe geliefert — `make test-sdk-python-integration` lief
  real grün (`change_id=804-1` über die Runner-Assertion-Kette empfangen,
  über `cdc.changes` gegengeprüft; Lauf-Beleg zitiert im
  Verifikationsbericht §1 DoD 2). Reviewer (F-1…F-6) und Verifier
  (V-1…V-3) laufen unabhängig voneinander; die beiden Fixrunden lösen
  ihre Befunde je in einem Commit (`b7a993fe`, `2eb18947`), der
  Re-Review-Endstand trägt 0 HIGH.
- **Was ging anders als geplant:** drei Zahl-/Stand-Driften entstehen in
  der Korrektur-Kette selbst, nicht in der Umsetzung — jede
  Nachzug-Zeile, die eine Zahl oder Version trägt, zog den Wert aus dem
  Vorgängerstand des Trägers statt aus der Messung, die sie belegt:
  5/29 → 6/29 → 1.39/1.40 (Haupt-Review F-1, Verifikation V-1, Re-Review
  F-5; je vom nächsten Leser durch Nachzählen statt Übernahme gefunden).
  Die Plan-Nachzug-Tabelle trägt dadurch die V-3-Deklaration
  (`options.py`) und die zwei Fixrunden-2-Zeilen — sie bleibt das
  Übergabe-Artefakt der nächsten Rollen, der Code-Diff ist es nicht.
- **Steering-Loop-Eintrag:** geschärfte Regel (Anwendungs-Schärfung der
  bereits verkörperten Klasse `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`,
  Anker [`AGENTS.md`](../../../../AGENTS.md) §3.12 Instanz A ·
  seit slice-089): eine Zahl-/Versions-Nachzug-Zeile wird aus
  **derselben Messung** gezogen, die sie belegt — Messung zuerst, dann
  Text —, nicht aus dem Vorgängerstand des Trägers. Kein neuer Sensor:
  die verfügbare Falsifikation bleibt die Messung selbst, so wie der
  Eintrag in seiner `state.md` sie führt. Keine benannte Spec-Lücke. Ob
  die Schärfung einen eigenen Satz in `AGENTS.md` §3.12 Instanz A trägt,
  prüft der Lese-Schritt der Welle-Closure — die Klasse steht mit diesem
  Beleg bei 13× (Datei-Anzahl unter `evidence/`, real ausgezählt),
  bereits verkörpert.
- **Beobachtungs-Register (`../observations/`):**
  - **`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`** — neuer,
    dreizehnter Beleg:
    `evidence/slice-sdk-python-grpc-client-flaeche.md`; die drei
    Auftreten der Kette (Haupt-Review F-1, Verifikation V-1, Re-Review
    F-5) tragen je ihre Fundstellen; `state.md` trägt den abgeleiteten
    Zähler.
  - **`BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`** — geprüft: die
    Klasse traf den Verbund (Haupt-Review F-2, Re-Review F-1/F-2); beide
    Funde sind in der Kette gelöst (Fixrunde 1 / Fixrunde 2) und in den
    Review-Reports konserviert — kein zusätzlicher Beleg in dieser
    Closure; der Lernwert der Kette läuft über die Zahl-Drift-Beobachtung
    (deren drei Funde dieselbe Wurzel tragen).
  - **`BEO-PGC/arbeit-ueberholt-stehenden-traeger`** — geprüft: F-7 des
    Re-Reviews (die überholte `options.py`-Planzeile des offenen
    NATS-Folge-Slices) ist ein Übergabe-Artefakt an den Folgelauf, kein
    Registereintrag; die fixgenerierten Funde F-3/F-4 sind in der Kette
    gelöst und in den Reports konserviert.
  - Re-Review F-6 (Verifikations-Report und Auflagen-Korrekturen in
    einem Commit) — Prozessbeobachtung ohne Registereintrag; die
    Übergabe-Artefakte (Plan-Nachzug-Notizen) tragen die Zuordnung
    (Re-Review: „erwartet keine Aktion").
- **Folge-Slices:** `slice-sdk-python-sse-client-flaeche` und
  `slice-sdk-python-nats-stream-client-flaeche` (Reihenfolge Welle-Plan
  §4) — sie erweitern das hier eingeführte Werkzeug; der NATS-Folgelauf
  reduziert die überholte `options.py`-Planzeile (F-7, Übergabe oben) und
  bündelt den SSE-/NATS-Handbuch-Teil, den Version-Bump und den
  `spec/pflichtenheft.md`-Nachzug. Kein zusätzlicher Slice aus dieser
  Closure.
- **Risiken aus §6:** Risiko 1 (Werkzeug-Laufzeit/Stabilität) —
  **entfallen** (der erste reale Lauf ist grün; Flakiness tritt in keinem
  der Läufe auf; Begründung und Beleg am Ort in §6). Risiko 2
  (Stub-Codestil) — **entfallen** (funktionale Stubs aus dem realen Bau,
  kein Blocker; §6). Risiko 3 — **entfallen** (unverändert; §6). Kein
  Ausgang „weiter offen" ins Register — ein Registereintrag für ein nie
  aufgetretenes Muster trüge kein `evidence/` (Paarung (c)).
- **Drei Paarungen:** dieser Slice gehört zu
  [welle-sdk-python-vollabdeckung](welle-sdk-python-vollabdeckung.md)
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
