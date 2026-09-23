# Slice sdk-python-nats-stream-client-flaeche: Öffentliche NATS-Vollinhalts-Client-Fläche (`SPEC-024`), Realserver-Beleg, Version-Hebung, Träger-Nachzug

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** [welle-sdk-python-vollabdeckung](welle-sdk-python-vollabdeckung.md).

**Bezug:** [`LH-FA-SST-009`](../../../../spec/lastenheft.md),
[`LH-FA-SST-008`](../../../../spec/lastenheft.md) (Live-Change-Stream),
[`ADR-0110`](../../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
§Entscheidung Festlegung 1/2/3 (Umfang, verschärfte Test-Pflicht,
Struktur bleibt v2 desselben Packages), [`ADR-0100`](../../adr/0100-nats-dritter-vollinhalts-zustellweg.md)
(NATS-Vollinhalts-Stream-Vertrag), [`AGENTS.md`](../../../../AGENTS.md)
§3.13 (Träger-Nachzug — dieser Slice löst ihn für
`spec/pflichtenheft.md` **und** `options.py`s Docstring aus).

**Berührte Spec-Stellen:** [`SPEC-024`](../../../../spec/pflichtenheft.md)
(Subjekt-Schema, Nachrichtenform, Authentifizierung), `SPEC-027`
(`pgchangefeed`-Metadaten — Version-Hebung),
[`LH-FA-SST-009.a`](../../../../spec/pflichtenheft.md) (Träger-Nachzug:
„deckt HTTP-API" wird durch diese Welle falsch).

**Verantwortlich:** Implementer-Agent, 2026-09-23.

**Autor:** Planner-Agent, direkt beauftragt (`ADR-0110` §Konsequenzen
Folgepflicht 1/2). **Datum:** 2026-09-21.

---

## 1. Ziel und Abgrenzung

**Ziel:** Eine öffentliche, stabile Python-API-Fläche im bestehenden
Package `pgchangefeed` für den NATS-Vollinhalts-Stream von
[`SPEC-024`](../../../../spec/pflichtenheft.md) — Verbindung über eine
NATS-Python-Client-Bibliothek (`nats-py`, Arbeitsname, asyncio-basiert —
real zum Bau-Zeitpunkt zu prüfen/messen, analog `NATS.Net`/`io.nats:jnats`
bei den Schwester-SDKs), Subjekt-Abonnement
`cdc.stream.<source_id>.<schema>.<table>`, Token-Auth auf
Verbindungsebene, Deserialisierung der zehn Nachrichtenfelder von
`SPEC-024`. **Zusätzlich** erweitert dieser Slice das
Realserver-Integrationstest-Werkzeug um einen NATS-Rundlauf und bündelt,
als letzter Flächen-Slice dieser Welle: Version-Hebung von
`pgchangefeed` und den Träger-Nachzug in
`spec/pflichtenheft.md`/`docs/user/benutzerhandbuch.md` für **alle drei**
neu gelieferten Flächen (gRPC, SSE **und** NATS-Vollinhalt).

**Übernimmt:** `slice-sdk-python-grpc-client-flaeche` und
`slice-sdk-python-sse-client-flaeche` — nicht deren Code, sondern deren
offen gelassenen Doku-/Träger-Nachzug (§1 jeweils „Ausdrücklich NICHT in
diesem Slice").

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **HTTP-API-, gRPC- oder SSE-Fläche** — HTTP bereits geliefert; gRPC und
  SSE sind die direkten Vorgänger-Slices dieser Welle.
- **Ein `examples/python/`-Referenz-Client** — `ADR-0110` §Entscheidung
  Festlegung 2 mandatiert das ausdrücklich nicht.
- **Aufspaltung in ein eigenständiges Folge-Package** — `ADR-0110`
  §Entscheidung Festlegung 3 überlässt das einem künftigen Zug.
- **Ein realer `sdk-python-v<Version>`-Tag-Push** —
  Betreiber-Entscheidung nach `AGENTS.md` §3.10 (Welle-Plan §3).
- **Trusted Publishing (PyPI OIDC)** — unverändert `ADR-0107`
  §Re-Evaluierungs-Trigger 4.

## 2. Definition of Done

- [x] `LH-FA-SST-009` erfüllt: eine öffentliche, im Package sichtbare
      Client-Klasse verbindet sich über eine NATS-Python-Client-Bibliothek
      mit Token-Auth, abonniert den Vollinhalts-Namensraum und liefert die
      zehn Nachrichtenfelder von
      [`SPEC-024`](../../../../spec/pflichtenheft.md) — Unit-Tests
      referenzieren `SPEC-024` (Nachrichtenschema-Vollständigkeit,
      Subjekt-Formatierung netzlos, Verbindungs-Fehlerpfad gegen einen
      Fake/gestubbten Client). *(Sensor-Beleg: Unit-Suite 48 Tests grün —
      `grep -c "^def test_"` = 20 http + 3 options + 8 grpc + 10 sse + 7
      nats = 48; Pack-Lauf „48 passed"; eigener Verifier-Lauf netzlos
      48 passed gegen den HEAD-Stand; Subjekt-Form, Token-im-Connect,
      zehn Felder und Fire-and-Forget einzeln gegen `SPEC-024` und die
      Server-Gegenprobe `publisher.go` gehalten — Verifikation §2.1/§4.)*
- [x] `sdks/python/pgchangefeed/pyproject.toml` bekommt die
      NATS-Python-Client-Bibliothek als gepinnte Laufzeitabhängigkeit,
      real zum Bau-Zeitpunkt gemessen (`AGENTS.md` §3.12 — nicht blind aus
      den Schwester-SDKs übernommen, andere Sprache, andere Bibliothek).
      *(Sensor-Beleg: `"nats-py>=2"` in `dependencies`, Untergranz-Form der
      drei Geschwister-Abhängigkeiten; real aufgelöst im Pack-Bau —
      „Downloading nats_py-2.16.0-py3-none-any.whl"; Verifikation §2.2.)*
- [x] **Realserver-Integrationstest erweitert** (`ADR-0110` §Entscheidung
      Festlegung 2/Folgepflicht 1): eine neue Testdatei unter
      `sdks/python/pgchangefeed/tests/integration/`
      (`test_nats_stream_realserver.py`, Arbeitsname) läuft über
      `make test-sdk-python-integration` und belegt real den Empfang
      einer committeten Änderung über NATS — belegt am Nachrichtenschema
      und über die `change_id` gegen `cdc.changes` gehalten. *(Sensor-Beleg:
      Pfad-Abweichung im Plan-Nachzug deklariert (`integration/`-Ordner
      seit dem gRPC-Slice); s3f-Lauf 2026-09-23 08:10 grün — gRPC
      change_id=807-1, SSE 810-1, NATS 812-1, je SQL-Gegenprüfung
      (Verifikation §2.3); eigener frischer Welle-Closure-Lauf Exit 0 —
      gRPC 804-1, SSE 808-1, NATS 810-1, je SQL-Gegenprüfung; Reject-
      Ablehnung real am Wire (s3e→s3f rot→grün).)*
- [x] `pyproject.toml`s `version` wird gehoben (von `0.1.0` auf `0.2.0` —
      additive, rückwärtskompatible Erweiterung). *(Sensor-Beleg:
      `version = "0.2.0"` (Zeile 7); Wheel-METADATA `Version: 0.2.0`
      (extrahiert); Verifikation §2.4.)*
- [x] `spec/pflichtenheft.md` `LH-FA-SST-009.a`/`SPEC-027` nachgezogen: der
      Satz „deckt HTTP-API" wird zu „deckt HTTP-API, gRPC-Stream, SSE und
      NATS-Vollinhalt" (`AGENTS.md` §3.13) — Suchlauf-Pflicht: `grep -rn
      "deckt HTTP-API\|pgchangefeed-0.1.0\|no gRPC/SSE/NATS surface exists yet"`
      über `spec/`, `docs/`, `sdks/python/` (Ergebnis in §7 berichtet,
      **einschließlich** `options.py`s Docstring, dem bereits bekannten
      Fundort dieser Klasse, `BEO-PGC/arbeit-ueberholt-stehenden-traeger`).
      *(Sensor-Beleg: Suchlauf-Feld committet (§3, vor der Closure-Notiz);
      F-1 (§6 `SPEC-027`-Zeile) und F-4 (fünf Träger-Stellen) in Fixrunde 1
      gelöst, R-1/R-2 (drei Reststellen) in Fixrunde 2; Verifikation §2.5
      hält beide Suchlauf-Stände am HEAD nach — kein lebender Träger
      trägt eine stale Form.)*
- [x] Kein Import aus `internal/**`/`cmd/**`/`gen/**` dieses Repos.
      *(Sensor-Beleg: Grep über alle drei neuen Python-Dateien — kein
      `internal/`/`cmd/`/`gen/`-Bezug; `examples/python/` existiert nicht
      (Verzeichnis gemessen); Verifikation §2.6.)*
- [x] `make gates` grün. *(Sensor-Beleg: Verifier-Lauf Exit 0, ungepiped,
      byte-gleiches Stempel-Precheck; Closure-Läufe je Commit unten in der
      Commit-Liste.)*
- [x] Ein real neu gebautes Artefakt-Paar (`make sdk-pack-python`) mit
      `version=0.2.0` als Smoke-Beleg — alle vier Client-Flächen im
      selben Artefakt. *(Sensor-Beleg: frisches Paar nach Fixrunde 2 —
      20503/24655 Bytes, Stempel 2026-09-23 08:50; Wheel-Inhalt byte-gleich
      gegen den HEAD-Quellstand gehalten (`nats_stream_client.py`,
      `exceptions.py`, `options.py` je diff-leer), METADATA
      `Version: 0.2.0`, kein `noqa`-Treffer — die Abweichung V-1 der
      Verifikation (Artefakt-Frische) ist damit real gelöst; alle vier
      Flächen im selben Artefakt.)*
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      *(Haupt-Review F-1…F-11 (3 HIGH/3 MEDIUM/4 LOW/1 INFO) → Fixrunde 1
      (`b4d1d352`) → Re-Review R-1…R-4 im Fixrunden-Report (1 HIGH/1
      MEDIUM/1 LOW/1 INFO) → Fixrunde 2 (`a3e9d64e`) → Verifikation
      (V-1/V-2). V-2: Fixrunde 2 lief ohne Re-Review-Stufe — die
      Verifikation schloss sie durch Nachmessen der drei Residuen am HEAD;
      Prozess-Notiz in §7.)*
- [x] Doku-Update: `docs/user/benutzerhandbuch.md` bekommt einen
      SDK-Hinweis für gRPC, SSE **und** NATS-Vollinhalt gebündelt.
      *(Sensor-Beleg: Handbuch `Version: 1.43`, Änderungshistorie-Zeile
      1.43 datiert 2026-09-23; gRPC-Anteile 1.40/1.41 und SSE-Anteil 1.42
      tragen die Vorgänger-Slices, der NATS-Absatz liegt in §4; Verifikation
      §2.10.)*
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag. *(§7 unten.)*
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben,
      **falls dieser Slice einen Inventur-Fund auflöst** — entfällt: keine
      Reconciliation-Datei in diesem Repo. *(gemessen: Datei existiert
      nicht.)*
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder Beleg in `evidence/`; keine Beobachtung angefallen
      ist ebenfalls eine Antwort und wird in §7 notiert. *(drei Belege:
      `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (23.),
      `BEO-PGC/zitat-nennt-die-falsche-stelle` (F-3),
      `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` (F-6/R-1); übrige
      Finding-Klassen geprüft, siehe §7.)*
- [x] Jedes Risiko aus §6 trägt einen Ausgang. *(alle vier Ausgänge final
      in §6 — siehe §7.)*
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      dieser Slice gehört zu
      [welle-sdk-python-vollabdeckung](welle-sdk-python-vollabdeckung.md)
      (noch offen); die Prüfung läuft regelkonform bei deren Closure.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `sdks/python/pgchangefeed/src/pgchangefeed/nats_stream_client.py` (Arbeitsname) | neu | öffentliche API-Fläche für den NATS-Vollinhalts-Stream (`SPEC-024`). |
| `sdks/python/pgchangefeed/pyproject.toml` | update | NATS-Python-Client-Bibliothek ergänzen, `version` auf `0.2.0` heben. |
| `sdks/python/pgchangefeed/tests/test_nats_stream_client.py` (Arbeitsname) | neu | Nachrichtenschema-Vollständigkeit, Subjekt-Formatierung, Verbindungs-Fehlerpfad — netzlos. |
| `sdks/python/pgchangefeed/tests/integration/test_nats_stream_realserver.py` (Arbeitsname) | neu | Realserver-Rundlauf: SDK empfängt real eine committete Änderung über NATS. |
| `sdks/python/pgchangefeed/src/pgchangefeed/options.py` | update | Docstring-Satz „no gRPC/SSE/NATS surface exists yet" entfernen/korrigieren (`AGENTS.md` §3.13). |
| `spec/pflichtenheft.md` §1 (`LH-FA-SST-009.a`), §6 (`SPEC-027`-Zeile) | update | Träger-Nachzug: volle Vier-Wege-Abdeckung für Python. |
| `docs/user/benutzerhandbuch.md` | update | SDK-Hinweis für gRPC, SSE und NATS-Vollinhalt. |

**Plan-Nachzug (im selben Lauf, vor dem Gate-Lauf):**

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `sdks/python/pgchangefeed/integration/test_nats_realserver.py` (Arbeitsname, statt `tests/integration/`-Pfad) | Abweichung | die Realserver-Tests liegen seit dem Vorgänger-Slice im Geschwister-Ordner `integration/`; die §3-Zeile trug noch den Arbeitsnamen `test_nats_stream_realserver.py` unter `tests/`. |
| `tools/harness/run-sdk-python-integration-tests.sh` | update | dritte Phase im `run_surface_phase`-Muster (Env-Liste je Fläche), NATS-Sentinel/-ID-Bereich 320ff., Reject-Marker `REJECTED token-rejected` (Ablehnung am Verbindungsversuch, `allow_reconnect=False`), SQL-Gegenprüfung gegen `cdc.changes`. |
| `sdks/python/pgchangefeed/src/pgchangefeed/__init__.py` (Docstring + Export `PgChangeFeedNatsStreamClient`), `options.py`-Docstring, `sdks/python/README.md` §Status, pyproject-Beschreibung | update | Träger-Nachzug (`AGENTS.md` §3.13): die Satzform „NATS bleibt außerhalb" wird durch diesen Slice falsch; das Package trägt jetzt die volle Vier-Wege-Matrix. |
| `spec/pflichtenheft.md` §7 Historie-Chronik-Zeile | neu | der Nachzug trägt seine Chronik-Zeile in §7 Historie (Muster der Vorgänger-Nachzüge 2026-09-22). |

**§3.13-Suchlauf (committetes Feld — Suchlauf-Pflicht des DoD; Raum `spec/`, `docs/`, `sdks/python/`, darüber hinaus `harness/` und die Wurzel-READMEs, wie die Fixrunden-Funde belegen; das Ergebnis steht hier, vor der Closure-Notiz):**

| Fundstelle (Schreibform am Ist-Baum) | Befund | Behandlung |
|---|---|---|
| `spec/pflichtenheft.md` §1 (Python-Absatz) | „deckt HTTP-API" gefunden | gezogen: „deckt HTTP-API, gRPC-Stream, SSE und NATS-Vollinhalt" + Endstand-Artefaktnamen 0.2.0 |
| `spec/pflichtenheft.md` §6 Vertragszeile `SPEC-027` | „aktuell `0.1.0`" + Vertrag-Liste `SPEC-018` gefunden (Review F-1) | gezogen: `0.2.0` + volle Vertragsliste |
| `spec/pflichtenheft.md` §7 Historie | Chronik-Zeile fehlte | ergänzt (Plan-Nachzug) |
| `harness/README.md` §Werkzeuge (sdk-pack-python-Zeile) | 0.1.0-Artefaktnamen gefunden | gezogen auf 0.2.0 (Endstand) |
| `harness/README.md` §Werkzeuge (test-sdk-python-integration-Zeile) | Zwei-Flächen-Form gefunden; die Reject-Form-Zeile trägt drei Formen (Fixrunde 2, Re-Review R-2) | gezogen |
| Wurzel-`README.md` (Distribution-Zeile) | „HTTP-API als offizielle Python-Bibliothek" gefunden | gezogen auf dieselben vier Zustellwege |
| `README.de.md` (Distribution-Zeile) | dieselbe Zwilling-Zeile gefunden, im ersten Suchlauf übersehen (Re-Review R-1/R-2) | Fixrunde 2: gezogen auf dieselben vier Zustellwege |
| `sdks/python/README.md` §Intro | „consume the HTTP API" gefunden | gezogen auf die vier Zustellwege |
| `sdks/python/pgchangefeed/src/pgchangefeed/options.py` | Docstring-Satz „NATS still follows" gefunden | gezogen auf `PgChangeFeedNatsStreamClient` |
| `docs/plan/planning/welle-sdk-python-vollabdeckung.md` (Trigger-Zeile 56) | 0.1.0-Artefaktnamen gefunden — **historisch korrekt** (Trigger-Ist-Stand von 2026-09-21) | belassen (Records-Herkunfts-Form) |
| historische Records (`done/`, `docs/reviews/**` bereits geschlossener Züge) | dieselben Sprachformen — Records tragen keine §Geschichte | belassen |

**Ansatz:** Referenzmaterial für Nachrichtenschema/Subjekt-Form ist
`examples/csharp/nats-stream-client/Format.cs`/
`examples/kotlin/nats-stream-client/…/Format.kt` (Fremdsprachen-Referenzen
gegen denselben Draht) und `internal/adapters/driven/natsstream/publisher_test.go`
(nur **gelesen**, serverseitige Feldvollständigkeits-Gegenprobe, kein
Import — `ADR-0110` nennt genau diese Datei als Draht-Kenntnis-Quelle).
Die konkrete NATS-Python-Bibliothek (voraussichtlich `nats-py`, das
offizielle asyncio-Client-Paket des NATS-Projekts) wird zum
Implementierungszeitpunkt real gegen PyPI gemessen, nicht hier
vorentschieden.

## 4. Trigger

**Start** (`next` → `in-progress`): wenn
`slice-sdk-python-sse-client-flaeche` in `done/` liegt (Welle-Plan §4
Reihenfolge — strikt sequentiell).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls der
  Träger-Nachzug (drei Flächen gebündelt) mehr Stellen betrifft als
  erwartet, oder falls `nats-py`s asyncio-Modell eine grundlegend andere
  Integrationsform braucht als die synchronen HTTP-/gRPC-Flächen —
  dann Abspaltung eines eigenen Doku-Nachzug- bzw.
  Async-Integrations-Slice.
- `in-progress` → `open` (blockiert — Carveout?): `slice-sdk-python-sse-client-flaeche`
  liegt noch nicht in `done/`.

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + reales Artefakt-Paar mit `0.2.0` +
ein realer, grüner `make test-sdk-python-integration`-Lauf gegen die
NATS-Fläche + Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- Die gewählte NATS-Python-Bibliothek (`nats-py`, Arbeitsname) ist
  asyncio-basiert, während die bestehende HTTP-/gRPC-/SSE-Fläche dieses
  Packages synchron ist (`httpx.Client`, synchroner gRPC-Stub) —
  ein Stilbruch im öffentlichen API-Vokabular des Packages.
  **Ausgang (Closure):** entfallen — der Sync-Wrapper-Weg wurde real
  umgesetzt (Connect und Subscription laufen in einem dedizierten Thread
  mit eigener Event-Loop, der Generator konsumiert eine Queue; ein
  Iterator-Vokabular über alle vier Flächen); die dokumentierte Diskrepanz
  trägt der Modul-Docstring am Ort („one iterator form across all four
  surfaces" — Fixrunden-Report, Negativbefund „Sync-Wrapper-Exit am Ort
  benannt"). Kein Registereintrag angefallen — ein Nutzer-Stoß an der
  Grenze ist in keinem Lauf oder Review belegt.
- Ein gefakter NATS-Verbindungsfehler-Test könnte das reale
  Token-Ablehnungsverhalten des NATS-Servers nicht exakt nachbilden.
  **Ausgang:** entfallen — der reale Beleg liegt im
  Integrationstest-Werkzeug dieses Slice selbst (anders als bei
  C#/Kotlin, wo dieser Beleg bewusst weiter offen bleibt).
- Der Träger-Nachzug (drei Flächen gebündelt, `options.py`-Docstring
  zusätzlich) könnte eine Stelle übersehen. **Ausgang (Closure):**
  eingetreten, im selben Slice gelöst — der Implementer-Suchlauf
  überließ real Fundstellen dem unabhängigen Reviewer (F-4: fünf
  Träger-Stellen in der tatsächlichen Schreibform, Fixrunde 1) und die
  Fixrunde hinterließ Reste samt zweier Behauptungen im committeten
  Suchlauf-Feld, die der Baum widerlegte (R-1/R-2: drei Reststellen,
  Fixrunde 2); beide Fixrunden zogen die Stellen real. Die Lücken-Struktur
  geht als 23. Beleg der bereits verkörperten Klasse
  `BEO-PGC/arbeit-ueberholt-stehenden-traeger` ins Register
  (`evidence/slice-sdk-python-nats-stream-client-flaeche.md`) —
  kein Folge-Slice, kein Carveout.
- Version-Hebung `0.1.0` → `0.2.0` ohne explizites Minor-Inkrement-Schema
  in einer Datei dieses Repos. **Ausgang:** entfallen — additive,
  rückwärtskompatible Erweiterung ist in PEP 440 (kompatibel zu SemVer,
  `ADR-0107` Festlegung 4) unmissverständlich ein Minor-Bump.

## 7. Closure-Notiz

- **Was hat funktioniert:** der Plan-Schnitt trägt sich in einem Zug —
  die NATS-Vollinhalts-Fläche, die dritte Runner-Phase, die
  Version-Hebung und der Träger-Nachzug für alle drei neuen Flächen
  gehören in einen Slice, weil der Nachzug erst mit dem letzten
  Flächen-Zug seine volle Wahrheit erreicht ([`ADR-0110`](../../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
  §Entscheidung Festlegung 2, Folgepflicht 2). Beide Rückführungen aus §4
  stehen nicht an: der Vorgänger-Slice liegt in `done/`, und die
  asyncio-Bibliothek fügt sich in ein synchrones Iterator-Vokabular —
  `make test-sdk-python-integration` lief real grün über alle drei
  Flächen (s3f-Lauf 2026-09-23 08:10: gRPC change_id=807-1, SSE 810-1,
  NATS 812-1, je SQL-Gegenprüfung gegen `cdc.changes`; eigener frischer
  Welle-Closure-Lauf Exit 0: gRPC 804-1, SSE 808-1, NATS 810-1). Die
  Rollen-Kette läuft unabhängig: Haupt-Review F-1…F-11, Fixrunde 1
  (`b4d1d352`) löste F-1…F-11, Re-Review R-1…R-4 im Fixrunden-Report,
  Fixrunde 2 (`a3e9d64e`) löste R-1…R-3 und zog R-4 (Raum-Kopf) nach;
  die Verifikation fährt eigene Läufe (Unit-Suite netzlos 48 passed gegen
  den HEAD-Stand, `make gates` Exit 0, Wheel-Extraktion Feld für Feld,
  Suchlauf-Stände an Parent und HEAD gemessen) und trägt die zwei
  Abweichungen V-1/V-2. Der s3e→s3f-Wechsel ist zugleich ein lebender
  Mutations-Beweis der Reject-Bindung: die falsche Klasse
  (`AuthorizationError`) kippte den Lauf real rot→grün gegen denselben
  Server, die korrekte Bindung (`nats.errors.Error`,
  `match="Authorization Violation"`) hält.
- **Was ging anders als geplant:** drei Lücken-Klassen in derselben
  Korrektur-Kette, keine in der Umsetzung: (1) der §3.13-Suchlauf
  arbeitete mit der schmalen Plan-Mustersatz-Grep-Form über `spec/`,
  `docs/`, `sdks/python/` und verfehlte alle fünf Träger-Stellen in deren
  tatsächlicher Schreibform — `harness/README.md` und die Wurzel-READMEs
  lagen außerhalb des Pfadraums (F-4, gefunden vom Reviewer; Fixrunde 1).
  (2) Das committete Suchlauf-Feld behauptete zwei vollzogene Nachzüge,
  die der Baum widerlegte (R-1/R-2; Fixrunde 2) — das Feld ist selbst
  ein Träger, dessen Behandlungs-Angaben gemessen werden müssen, bevor
  die Closure-Notiz sie liest. (3) Die Verifikation fand das 0.2.0-Artefakt-Paar
  byte-stale gegen den nach zwei Fixrunden weitergeänderten Quellstand
  (V-1) — gelöst durch einen frischen Neubau nach Fixrunde 2
  (20503/24655 Bytes, Wheel-Inhalt byte-gleich gegen den HEAD-Stand
  gemessen). Prozess-Notiz: die zweite Fixrunde lief ohne Re-Review-Stufe
  (V-2); die Verifikation schloss die drei Residuen durch eigenes
  Nachmessen am HEAD — für künftige Fixrunden bleibt der Regelfall, dass
  eine Fixrunde eine Reviewer-Gegenprüfung bekommt, oder der Verifier
  die Lücke explizit schließt und als Prozess-Notiz trägt.
- **Steering-Loop-Eintrag:** Anwendungs-Schärfung der bereits
  verkörperten Klasse `BEO-PGC/arbeit-ueberholt-stehenden-traeger`
  (Anker [`AGENTS.md`](../../../../AGENTS.md) §3.13 · seit welle-20), zwei
  Lehren aus dieser Kette: (a) der Suchlauf folgt §3.13s Wortlaut —
  `grep` über die Träger nach der bewegten Eigenschaft — und nicht dem
  Plan-Mustersatz; die schmale Form verfehlte alle fünf F-4-Stellen in
  ihrer Schreibform und ließ `harness/` und die Wurzel-READMEs außerhalb
  des Pfadraums. (b) Das committete Suchlauf-Feld ist selbst ein Träger:
  seine Behandlungs-Angaben werden am Baum nachgemessen (R-1 widerlegte
  zwei von ihnen), nicht deklariert. Kein neuer Sensor — die verfügbare
  Falsifikation bleibt die Messung an beiden Ständen (§3.13s Grenze für
  Zahlen/Prosa-Umformulierungen trägt unverändert); keine benannte
  Spec-Lücke. Ob die Schärfung einen eigenen Satz in `AGENTS.md` §3.13
  trägt, prüft der Lese-Schritt der Welle-Closure — die Klasse steht mit
  diesem Beleg bei 23× (Datei-Anzahl unter `evidence/`, real ausgezählt),
  bereits verkörpert.
- **Beobachtungs-Register (`../observations/`):**
  - **`BEO-PGC/arbeit-ueberholt-stehenden-traeger`** — neuer, 23. Beleg:
    `evidence/slice-sdk-python-nats-stream-client-flaeche.md`; die Kette
    (F-1 dual, F-4, R-1 dual, R-2) trägt je ihre Fundstellen; `state.md`
    trägt den real ausgezählten Zähler.
  - **`BEO-PGC/zitat-nennt-die-falsche-stelle`** — neuer Beleg (F-3,
    HIGH): der `__init__.py`-Docstring zitierte `ADR-0110` Festlegung 3
    als entscheidende Schicht für „v2 desselben Packages"; die Festlegung
    entscheidet ausdrücklich nichts, die entscheidende Schicht ist der
    Welle-Plan §6. Fixrunde 1 re-anchorte und das Fixrunden-Re-Review
    maß beide Anker nach (`evidence/slice-sdk-python-nats-stream-client-flaeche.md`
    dort).
  - **`BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`** — neuer Beleg
    (F-6, R-1): die Reject-Assertion band die Ursache nicht, die der
    Marker behauptete (gelöst: `nats.errors.Error` +
    `match="Authorization Violation"`, gegen die nats-py-Quellen
    v2.11.0/2.16.0 nachgemessen, s3e→s3f rot→grün am selben Server); R-1
    zählt derselben Klasse über die Behandlungs-Angaben des Suchlauf-Felds
    (`evidence/slice-sdk-python-nats-stream-client-flaeche.md` dort).
  - **`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`** — geprüft:
    F-1 ist ein echter Zahl-Drift („aktuell `0.1.0`" in der §6-Zelle gegen
    real gemessene `0.2.0`), aber die Zahl war bei ihrer letzten
    Niederschrift wahr und wurde durch diese Arbeit falsch — nach der
    Klassen-Grenze des Eintrags („dort driftet ein Träger schon; hier
    driftet er durch diese Arbeit") zählt der Vorgang bei
    `arbeit-ueberholt-stehenden-traeger` (F-1 dual, dort geführt); kein
    eigenständiger Fund der Zahl-Drift-Klasse in dieser Kette.
  - **`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`**
    — geprüft: kein Vorkommen; der Handbuch-Hinweis war als eigener
    DoD-Punkt gebündelt und real geliefert (1.43, Verifikation §2.10).
  - Übrige Finding-Klassen ohne Registereintrag (je in ihrer
    Fixrunde gelöst, Reports konservieren die Fundstellen): F-2
    (Suppression, §3.2-Träger), F-5 (Suchlauf-Ergebnis nicht committet —
    durch das §3-Feld gelöst), F-7/F-11/R-3 (Vakuum-Assertion,
    Unerreichbare Verzweigung — keine Klasse im Register), F-8 (Docstring-
    Enumeration unvollständig, 2. Auftreten, unter der Schwelle), F-9
    (Fehletikettierter Plan-Nachzug, einmalig), F-10 (fehlende
    EOF-Newline, 3. Auftreten eines etablierten Package-Musters, an die
    künftige Aufrum-Entscheidung gebunden), R-4 (Deklarations-Raum enger
    als dokumentierter Lauf, INFO — der Feld-Kopf trägt am HEAD den
    breiteren Raum).
- **Folge-Slices:** keiner — letzter Flächen-Slice der
  [Welle](welle-sdk-python-vollabdeckung.md); die Welle-Closure
  (Results-Datei, Roadmap-Nachzug, Lese-Schritt) folgt als eigener
  Planner-Zug. Ein realer `sdk-python-v0.2.0`-Tag-Push bleibt
  Betreiber-Entscheidung ([`AGENTS.md`](../../../../AGENTS.md) §3.10,
  Welle-Plan §3) — kein Folge-Slice.
- **Risiken aus §6:** Risiko 1 (asyncio-Stilbruch) — **entfallen** (der
  Sync-Wrapper-Weg real umgesetzt, dokumentiert am Ort; Beleg in §6).
  Risiko 2 (gefakter Reject-Test) — **entfallen** (der reale Beleg liegt
  im Integrationstest-Werkzeug selbst; s3e→s3f rot→grün am Wire).
  Risiko 3 (Träger-Nachzug übersehen) — **eingetreten, im selben Slice
  gelöst** (F-4 in Fixrunde 1, R-1/R-2 in Fixrunde 2; Beleg im Register,
  siehe oben). Risiko 4 (Minor-Bump-Schema) — **entfallen** (PEP 440,
  `ADR-0107` Festlegung 4; unverändert §6). Kein Ausgang „weiter offen"
  ins Register — der einzige materialisierte Fall trägt seinen Beleg als
  Beleg einer bestehenden Klasse (Paarung (c)).
- **Drei Paarungen:** dieser Slice gehört zu
  [welle-sdk-python-vollabdeckung](welle-sdk-python-vollabdeckung.md)
  (noch offen) — die Prüfung läuft regelkonform bei deren Closure:
  Anker (Welle-Link oben löst am Ruheort auf), Folge-Slice (keiner,
  letzter Flächen-Slice der Welle), Register (drei Evidenzdateien dieser
  Closure, jede mit Beleg).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Sub-Area `sdks/python/` — bereits
mit `slice-sdk-python-projektgeruest` eröffnet (GF), keine erneute
Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert, Suchlauf-Pflicht
trägt den Träger-Nachzug **und** `options.py`s Docstring dieses Slice),
`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`
(verkörpert, DoD-Punkt „Doku-Update" trägt sie), `BEO-PGC/release-mechanismus-nicht-in-releasing-doku-nachgezogen`
(offen, 1×, nicht einschlägig).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Fortsetzung von
`slice-sdk-python-projektgeruest`).
