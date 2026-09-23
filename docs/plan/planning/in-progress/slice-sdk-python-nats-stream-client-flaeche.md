# Slice sdk-python-nats-stream-client-flaeche: Öffentliche NATS-Vollinhalts-Client-Fläche (`SPEC-024`), Realserver-Beleg, Version-Hebung, Träger-Nachzug

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** [welle-sdk-python-vollabdeckung](../welle-sdk-python-vollabdeckung.md).

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

- [ ] `LH-FA-SST-009` erfüllt: eine öffentliche, im Package sichtbare
      Client-Klasse verbindet sich über eine NATS-Python-Client-Bibliothek
      mit Token-Auth, abonniert den Vollinhalts-Namensraum und liefert die
      zehn Nachrichtenfelder von
      [`SPEC-024`](../../../../spec/pflichtenheft.md) — Unit-Tests
      referenzieren `SPEC-024` (Nachrichtenschema-Vollständigkeit,
      Subjekt-Formatierung netzlos, Verbindungs-Fehlerpfad gegen einen
      Fake/gestubbten Client).
- [ ] `sdks/python/pgchangefeed/pyproject.toml` bekommt die
      NATS-Python-Client-Bibliothek als gepinnte Laufzeitabhängigkeit,
      real zum Bau-Zeitpunkt gemessen (`AGENTS.md` §3.12 — nicht blind aus
      den Schwester-SDKs übernommen, andere Sprache, andere Bibliothek).
- [ ] **Realserver-Integrationstest erweitert** (`ADR-0110` §Entscheidung
      Festlegung 2/Folgepflicht 1): eine neue Testdatei unter
      `sdks/python/pgchangefeed/tests/integration/`
      (`test_nats_stream_realserver.py`, Arbeitsname) läuft über
      `make test-sdk-python-integration` und belegt real den Empfang
      einer committeten Änderung über NATS — belegt am Nachrichtenschema
      und über die `change_id` gegen `cdc.changes` gehalten.
- [ ] `pyproject.toml`s `version` wird gehoben (von `0.1.0` auf `0.2.0` —
      additive, rückwärtskompatible Erweiterung).
- [ ] `spec/pflichtenheft.md` `LH-FA-SST-009.a`/`SPEC-027` nachgezogen: der
      Satz „deckt HTTP-API" wird zu „deckt HTTP-API, gRPC-Stream, SSE und
      NATS-Vollinhalt" (`AGENTS.md` §3.13) — Suchlauf-Pflicht: `grep -rn
      "deckt HTTP-API\|pgchangefeed-0.1.0\|no gRPC/SSE/NATS surface exists yet"`
      über `spec/`, `docs/`, `sdks/python/` (Ergebnis in §7 berichtet,
      **einschließlich** `options.py`s Docstring, dem bereits bekannten
      Fundort dieser Klasse, `BEO-PGC/arbeit-ueberholt-stehenden-traeger`).
- [ ] Kein Import aus `internal/**`/`cmd/**`/`gen/**` dieses Repos.
- [ ] `make gates` grün.
- [ ] Ein real neu gebautes Artefakt-Paar (`make sdk-pack-python`) mit
      `version=0.2.0` als Smoke-Beleg — alle vier Client-Flächen im
      selben Artefakt.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: `docs/user/benutzerhandbuch.md` bekommt einen
      SDK-Hinweis für gRPC, SSE **und** NATS-Vollinhalt gebündelt.
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

**§3.13-Suchlauf (committetes Feld — Suchlauf-Pflicht des DoD, Raum `spec/`, `docs/`, `sdks/python/`; das Ergebnis steht hier, vor der Closure-Notiz):**

| Fundstelle (Schreibform am Ist-Baum) | Befund | Behandlung |
|---|---|---|
| `spec/pflichtenheft.md` §1 (Python-Absatz) | „deckt HTTP-API" gefunden | gezogen: „deckt HTTP-API, gRPC-Stream, SSE und NATS-Vollinhalt" + Endstand-Artefaktnamen 0.2.0 |
| `spec/pflichtenheft.md` §6 Vertragszeile `SPEC-027` | „aktuell `0.1.0`" + Vertrag-Liste `SPEC-018` gefunden (Review F-1) | gezogen: `0.2.0` + volle Vertragsliste |
| `spec/pflichtenheft.md` §7 Historie | Chronik-Zeile fehlte | ergänzt (Plan-Nachzug) |
| `harness/README.md` §Werkzeuge (sdk-pack-python-Zeile) | 0.1.0-Artefaktnamen gefunden | gezogen auf 0.2.0 (Endstand) |
| `harness/README.md` §Werkzeuge (test-sdk-python-integration-Zeile) | Zwei-Flächen-/Zwei-Reject-Form gefunden | gezogen auf drei Flächen |
| Wurzel-`README.md` + `README.de.md` (Distribution-Zeile) | „HTTP-API als offizielle Python-Bibliothek" gefunden | gezogen auf dieselben vier Zustellwege |
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
  **Ausgang:** weiter offen bis zur Implementierung — falls sich ein
  synchroner Wrapper (`asyncio.run`) als unpraktikabel erweist, wird die
  Diskrepanz explizit im README dokumentiert statt verschwiegen (kein
  Blocker für diesen Slice, aber ein Kandidat für einen künftigen
  Registereintrag, falls Nutzer sich daran stoßen).
- Ein gefakter NATS-Verbindungsfehler-Test könnte das reale
  Token-Ablehnungsverhalten des NATS-Servers nicht exakt nachbilden.
  **Ausgang:** entfallen — der reale Beleg liegt im
  Integrationstest-Werkzeug dieses Slice selbst (anders als bei
  C#/Kotlin, wo dieser Beleg bewusst weiter offen bleibt).
- Der Träger-Nachzug (drei Flächen gebündelt, `options.py`-Docstring
  zusätzlich) könnte eine Stelle übersehen. **Ausgang:** weiter offen —
  Review prüft den Diff unabhängig gegen denselben Suchraum.
- Version-Hebung `0.1.0` → `0.2.0` ohne explizites Minor-Inkrement-Schema
  in einer Datei dieses Repos. **Ausgang:** entfallen — additive,
  rückwärtskompatible Erweiterung ist in PEP 440 (kompatibel zu SemVer,
  `ADR-0107` Festlegung 4) unmissverständlich ein Minor-Bump.

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
trägt den Träger-Nachzug **und** `options.py`s Docstring dieses Slice),
`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`
(verkörpert, DoD-Punkt „Doku-Update" trägt sie), `BEO-PGC/release-mechanismus-nicht-in-releasing-doku-nachgezogen`
(offen, 1×, nicht einschlägig).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Fortsetzung von
`slice-sdk-python-projektgeruest`).
