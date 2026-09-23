# Slice sdk-python-http-reale2e: Python-SDK — HTTP-Fläche in den bestehenden Realserver-Runner nachziehen, Träger-Erweiterung + Matrix-Vollständigkeit

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** [welle-sdk-reale2e](../welle-sdk-reale2e.md).

**Bezug:** [`LH-FA-SST-009`](../../../../spec/lastenheft.md)
(Client-Bibliotheken — die Fläche existiert, der Beleg-Stand wird stärker),
[`LH-FA-SST-006`](../../../../spec/lastenheft.md) (HTTP-API),
[`ADR-0110`](../../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
§Entscheidung Festlegung 2 (die etablierte Mechanik-Klasse — hier auf die
eine bestehende Fläche angewandt, die die Ursprungs-Pflicht nicht trug),
[`ADR-0107`](../../adr/0107-python-pypi-zweites-sdk-package.md)
(Ort `sdks/python/`, Import-Grenze, unverändert gültig),
[`ADR-0057`](../../adr/0057-http-grpc-api.md) (HTTP-API-Server-Vertrag,
wird vom SDK benutzt, nicht erweitert).

**Berührte Spec-Stellen:** [`SPEC-018`](../../../../spec/pflichtenheft.md)
(HTTP-API, neun Fähigkeiten) — gelesen als Draht-Vertrag, nicht geändert.

**Verantwortlich:** Implementer-Agent, 2026-09-23.

**Autor:** Planner-Agent, Welle-Eröffnung
[welle-sdk-reale2e](../welle-sdk-reale2e.md). **Datum:** 2026-09-23.

---

## 1. Ziel und Abgrenzung

**Ziel:** Der Realserver-Beleg für die **Python-HTTP-Fläche**
(`PgChangeFeedHttpClient`, `SPEC-018`) — die einzige der zwölf
Zustellweg-Flächen ohne reale Server-Probe (die drei Python-Stream-Flächen
tragen ihre Belege bereits, siehe
`slice-sdk-python-grpc-client-flaeche`/`slice-sdk-python-sse-client-flaeche`/
`slice-sdk-python-nats-stream-client-flaeche` in `done/`). Umfang:
(a) eine neue Integrationstest-Datei
`sdks/python/pgchangefeed/integration/test_http_realserver.py`
(Arbeitsname, Geschwister der drei bestehenden Realserver-Testdateien);
(b) eine **HTTP-Phase** im bestehenden Runner
`tools/harness/run-sdk-python-integration-tests.sh` — ein HTTP-Rundlauf
im Muster des Server-E2E-HTTP-Rundlaufs: Registrierung eines
Wegwerf-Consumers und Listen-Aufruf über die SDK-Methoden, die
Registrierung über `cdc.consumer` gegen den SQL-Lesezugriffsweg gehalten,
ein Token-freier Aufruf endet mit HTTP-Status 401 (der Runner erweitert
seine Phasen-Liste und seinen Träger-Abschnitt-Generator, unverändert
sonst); (c) Träger-Nachzug: `docs/user/sdk-e2e-abdeckung.md` trägt den
Python-Abschnitt (idempotent vom Runner geschrieben); (d) die
**Matrix-Vollständigkeits-Prüfung** — nach diesem Slice deklariert der
Träger alle zwölf Flächen-Belege (3 Sprachen × 4 Wege); (e)
`harness/README.md`-Nachzug der bestehenden
`make test-sdk-python-integration`-Zeile (Phasen-Liste um HTTP erweitert).

**Zuschnitt der HTTP-Phase** (der Auftrag lässt die Stufe beim Slice
entscheiden — hier entschieden): Der Pflichtkern ist der Rundlauf
RegisterConsumer + ListTables mit SQL-Gegenprüfung + 401-Ablehnung. Die
neun Fähigkeiten einzeln zu fahren, ist **nicht** der Pflichtkern — die
Unit-Tests (`tests/test_http_client.py`) tragen die Methoden-Fläche
bereits netzlos gegen `httpx.MockTransport`; der Realserver-Rundlauf
belegt die Protokoll-Annahmen am Wire (Auth-Header-Form, JSON-Mapping,
Fehler-Verhalten). `ReadChanges` über HTTP-GET bleibt eine
**Erweiterungs-Option** des umsetzenden Zuges (Rückführung §4): sinnvoll,
weil sie den Lese-Weg am Wire belegt, aber nicht nötig, um die
Fläche real abzunehmen — die drei Stream-Phasen belegen das
Änderungs-Empfangen bereits.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die C#-/Kotlin-Flächen** — Vorgänger-Slices
  (`slice-sdk-csharp-reale2e`, `slice-sdk-kotlin-reale2e`), die ihre
  eigenen Runner-Infrastrukturen tragen.
- **Die drei Python-Stream-Phasen** — bereits real belegt
  (`done/`-Slices oben); dieser Slice ändert an ihnen nichts.
- **Version-Bump, `spec/pflichtenheft.md`-Nachzug, Handbuch-Nachzug** —
  der Realserver-Test ist kein Artefakt-Vertrag; `SPEC-027`s
  Deckungs-Aussage bleibt richtig (Welle-Plan §6); kein falsch werdender
  Träger real geprüft (Welle-Plan §6, Stand 2026-09-23).
- **Eine Aufnahme des Targets in `make gates`** — unverändert Werkzeug
  (Welle-Plan §6).
- **Ein Umbau des Runners über die HTTP-Phase hinaus** — der Runner ist
  bewährt real (drei grüne Läufe über die Welle
  `welle-sdk-python-vollabdeckung`); dieser Slice erweitert ihn additiv
  (vierte Phase), ohne seine Struktur zu ändern.

## 2. Definition of Done

- [x] `LH-FA-SST-006`/`LH-FA-SST-009`-Beleg-Stand stärker: die
      Python-HTTP-Fläche trägt einen realen Rundlauf gegen eine laufende
      Server-Instanz — Registrierung eines Wegwerf-Consumers und
      Listen-Aufruf über `PgChangeFeedHttpClient`, die Registrierung
      über `cdc.consumer` gegen den SQL-Lesezugriffsweg gehalten; ein
      Token-freier Aufruf endet mit HTTP-Status 401.
      *(Sensor-Beleg: `make test-sdk-python-integration` EXIT=0 mit vier
      Phasen — der Lauf trug gRPC change_id=804-1, SSE 807-1, NATS 813-1
      (je SQL-Gegenprüfung gegen `cdc.changes`) und die Consumer-Registrierung
      (consumer_id=python-sdk-e2e-4ad61b9f8a0d, unabhängig über
      `cdc.consumer` lesbar); die 401-Ablehnung real. Mutation real
      gefahren: gültiger reader-Token im 401-Test — der Lauf färbte rot
      („rejects_the_call_without_token FAILED“), Revert, Abschlusslauf
      grün.)*
- [x] **Runner-Erweiterung**: `tools/harness/run-sdk-python-integration-tests.sh`
      trägt die HTTP-Phase als vierte Phase (eigene Sentinel-/ID-
      Wertebereiche, `REJECTED status=401`-Marker, explizite
      Testdatei-Auswahl — kein stiller Ausschluss des Rests);
      `integration/test_http_realserver.py` (Arbeitsname) entsteht als
      Geschwister-Testdatei der drei bestehenden Realserver-Testdateien
      (kein Verweis auf `tests/`-Unit-Quelltext nötig).
- [x] **Träger-Nachzug im selben Zug**: `docs/user/sdk-e2e-abdeckung.md`
      trägt den Python-Abschnitt (aus derselben Messung, die ihn belegt,
      idempotent vom Runner geschrieben) — und die DoD-Verankerung der
      **Matrix-Vollständigkeit**: nach diesem Slice deklariert der
      Träger alle zwölf Flächen-Belege (3 Sprachen × 4 Wege); ein
      fehlender Abschnitt oder eine Phantom-Zeile (Beleg ohne realen
      Lauf) bricht diese DoD sichtbar.
- [x] **Runner-Kopf-Nachzug** (`BEO-PGC/nachzug-laesst-ueberholten-text-stehen`):
      die Phasen-Dokumentation im Skript-Kopf trägt die vierte Phase,
      ohne die überholte Dreier-Form still stehen zu lassen; die
      bestehende `harness/README.md`-Zeile `make test-sdk-python-integration`
      zieht im selben Zug nach (Phasen-Liste, Ablehnungs-Beleg-Form —
      der umsetzende Lauf ist ihr Beleg, [`AGENTS.md`](../../../../AGENTS.md) §4).
- [x] Kein Import aus `internal/**`/`cmd/**`/`gen/**` dieses Repos in
      `sdks/python/**` (Import-Zeilen-Prüfung, Muster der bestehenden
      SDK-Slices).
- [x] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6),
      kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register — entfällt: keine Reconciliation-Datei in
      diesem Repo.
- [ ] Beobachtungs-Register fortgeschrieben — kein Anfall ist ebenfalls
      eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen getragen — dieser Slice gehört zu
      [welle-sdk-reale2e](../welle-sdk-reale2e.md).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `sdks/python/pgchangefeed/integration/test_http_realserver.py` (Arbeitsname) | neu | der HTTP-Rundlauf als Realserver-Testdatei (Geschwister von `test_grpc_realserver.py`/`test_sse_realserver.py`/`test_nats_realserver.py`). |
| `tools/harness/run-sdk-python-integration-tests.sh` | update | vierte Phase (HTTP) in der Phasen-Liste, eigener Sentinel-/ID-Bereich, Träger-Abschnitt-Erweiterung, Kopf-Nachzug der Phasen-Dokumentation. |
| `harness/mk/sdk.mk` | update | der Kommentar-Block des Targets `test-sdk-python-integration` zieht um die HTTP-Phase nach (Zustandsfeld-Form: der Zustand + Beleg, keine Chronik). |
| `docs/user/sdk-e2e-abdeckung.md` | update | Python-Abschnitt als Erzeugnis des Runners (idempotent, Marker-gegrenzt). |
| `harness/README.md` §Werkzeuge | update | die bestehende `make test-sdk-python-integration`-Zeile trägt die vierte Phase (aus derselben Messung geschrieben). |
| `sdks/python/Dockerfile` | **keine Änderung erwartet** | die `integration`-Stufe liest die Testdatei je Phase über `PGCHANGEFEED_TEST_FILE` — eine neue Testdatei im committeten `integration/`-Ordner braucht keine Stufen-Änderung (der `COPY pgchangefeed/integration`-Layer trägt den Ordner bereits). Abweichung wäre ein Plan-Nachzug. |

**Ansatz:** Der HTTP-Rundlauf folgt dem Muster des Server-E2E-HTTP-
Rundlaufs (`tools/harness/run-integration-tests.sh`, Phase
„HTTP-API-Rundlauf"): zwei Token-Klassen (admin für die Registrierung,
reader für den Listen-Aufruf — Container-Vertrag in `compose.yaml`), die
Registrierung gegen `cdc.consumer` gehalten. Die Phase wird in der
bestehenden `run_surface_phase`-Form geführt (eigene Sentinel-/ID-
Wertebereiche, die die Stream-Phasen nicht schneiden).

**Plan-Nachzug (im selben Lauf, vor dem Gate-Lauf):**

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/run-sdk-python-integration-tests.sh`: `run_surface_phase` um `received_grep`/`sql_kind` erweitert | Erweiterung | die HTTP-Phase braucht die `consumer`-SQL-Variante (Gegenprüfung gegen `cdc.consumer`) und eine phasenspezifische RECEIVED-Form — dieselbe Parametrisierung wie der C#-Runner (`run-sdk-csharp-integration-tests.sh`); die drei Stream-Phasen reichen ihre bisher feste Form als Parameter nach (kein Verhaltenstausch). |
| `tools/harness/run-sdk-python-integration-tests.sh`: Träger-Writer (Python-Abschnitt) | neu | der Python-Runner trägt seinen Abschnitt-Generator jetzt selbst (beidseitiger Rest-Erhalt, degenerater Pfad mit Kopf-Regeneration — die Form-Grenze der C#/Kotlin-Kette); die Matrix-Vollständigkeit (12 Zeilen) ist damit runner-generiert, nicht handgeschrieben. |
| `API_TOKEN_ADMIN`/`API_TOKEN_READER` im Runner-Kopf | neu | die HTTP-Phase braucht beide Token-Klassen (Risiko der C#-Kette, Ausgang dort belegt — dieselbe Ausgangslage hier). |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „die Python-HTTP-Fläche trägt einen realen Realserver-Beleg; der Träger deklariert alle zwölf Flächen-Belege"; beide Stände gemessen: Parent `d99768f2` und HEAD):**

| Träger | Befund | Behandlung |
|---|---|---|
| `docs/user/sdk-e2e-abdeckung.md` | Python-Abschnitt fehlte | in diesem Zug ergänzt (runner-generiert; C#/Kotlin-Abschnitte byte-identisch erhalten — beidseitiger Erhalt real gemessen) |
| `harness/README.md` §Werkzeuge (`make test-sdk-python-integration`-Zeile) | Phasen-Liste „grpc_client, sse_client, nats_stream_client“ ohne die HTTP-Fläche | gezogen: vier Flächen + HTTP-Ablehnungsform + cdc.consumer-Gegenprüfung |
| `tools/harness/run-sdk-python-integration-tests.sh` Kopf | Dreier-Phasen-Form (`nachzug-laesst-ueberholten-text-stehen`, im Plan §8 benannt) | gezogen: vierte Phase im Kopf |
| `sdks/python/README.md` §Status | geprüft — trägt keine Teststrategie-Aussage über Realserver-Läufe | kein Nachzug nötig |
| `docs/user/benutzerhandbuch.md` | geprüft — trägt die SDK-Hinweise, keine E2E-Beleg-Aussage | nichts zu ziehen |
| `spec/pflichtenheft.md` | geprüft — kein falsch werdender Träger (Welle-Plan §6) | nichts zu ziehen |

## 4. Trigger

**Start** (`next` → `in-progress`): wenn diese Welle eröffnet ist und kein
anderer Slice in `in-progress/` liegt (WIP-Limit 1); die C#-/Kotlin-
Vorgänger-Slices haben die Träger-Struktur real etabliert, die dieser
Slice vollendet.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): nicht erwartet —
  die Erweiterung einer bestehenden, bewährten Skript-Infrastruktur um
  eine Phase ist der kleinste Zuschnitt dieser Welle; falls der
  ReadChanges-Beleg (§1 Erweiterungs-Option) den Umfang sprengt, fällt
  er unter diese Rückführung.
- `in-progress` → `open` (blockiert — Carveout?): falls die
  `run_surface_phase`-Form für den zustandsbehafteten HTTP-Rundlauf
  (RegisterConsumer erzeugt Bestand, kein Stream-Empfang) wesentlich
  andere Assertions braucht — dann Rück zur Zerlegung mit einer eigenen
  HTTP-Runner-Phase-Form.

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + ein realer, grüner
`make test-sdk-python-integration`-Lauf mit allen vier Phasen +
Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **`run_surface_phase`-Form vs. HTTP-Rundlauf-Form:** die bestehende
  Funktion ist auf Stream-Empfang (Fire-and-Forget-Fenster,
  RECEIVED-Marker) zugeschnitten; der HTTP-Rundlauf erzeugt Bestand
  (ein registrierter Consumer) statt auf eine Change zu warten. *Erwartet,
  zu belegen durch:* der reale Lauf; falls die Form nicht passt, trägt
  der Fix eine eigene, schlankere HTTP-Phase-Funktion (Rückführung §4,
  zweiter Punkt).
- **Träger-Abschnitt-Kohärenz:** drei Runner schreiben je einen
  Abschnitt derselben Datei — eine Marker-Kollision (zwei Runner,
  dieselbe Grenz-Markierung) würde einen Abschnitt überschreiben.
  *Erwartet, zu belegen durch:* der reale Lauf dieses Slices liest und
  schreibt die Datei mit allen drei Abschnitten sichtbar; die
  Marker-Form ist je Sprache eindeutig (Sprach-Segment im Marker-Namen).
- **Sentinel-/ID-Kollisionen mit den Stream-Phasen:** die HTTP-Phase
  braucht eigene Wertebereiche (bestehende: 300/310/320 für gRPC/SSE/
  NATS). *Erwartet, zu belegen durch:* der reale Lauf.
- **Der ReadChanges-Beleg als Umfangs-Erweiterung** (§1 Erweiterungs-
  Option): falls der umsetzende Zug ihn fährt, trägt der DoD seine
  eigene SQL-Gegenprüfung gegen `cdc.changes` im Bereich `[from, to)`
  (Muster des bestehenden `GET /changes`-E2E-Belegs); falls nicht, wird
  die Option mit Ausgang „verworfen — Umfang" in §7 verzeichnet.

## 7. Closure-Notiz

- **Was hat funktioniert:** *(zu tragen bei Closure)*
- **Was ging anders als geplant:** *(zu tragen bei Closure)*
- **Steering-Loop-Eintrag (Lerneintrag):** *(zu tragen bei Closure —
  geschärfte Regel · neuer Sensor · benannte Spec-Lücke; ohne ihn kein
  `done/`-Übergang)*
- **Beobachtungs-Register (`../observations/`):** *(je Anfall Beleg oder
  „keine Beobachtung angefallen" als notierte Antwort)*
- **Risiken aus §6:** *(je ein Ausgang)*
- **Drei Paarungen:** dieser Slice gehört zu
  [welle-sdk-reale2e](../welle-sdk-reale2e.md) (offen) — die Prüfung
  läuft regelkonform bei deren Closure.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Sub-Area `sdks/python/` —
bereits mit `slice-sdk-python-projektgeruest` eröffnet (GF), keine
erneute Ausdifferenzierung nötig; `tools/harness/` ist Werkzeug-Area der
bestehenden Runner-Familie (GF, fortlaufend).

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/nachzug-laesst-ueberholten-text-stehen` (offen, 2×, einschlägig —
der Runner-Kopf trägt die überholte Dreier-Phasen-Form; der Nachzug
steht im DoD),
`BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert, Suchlauf
trägt den Runner-Kopf, die `harness/README.md`-Zeile und
`options.py`-Doku),
`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (verkörpert —
Träger-Abschnitt aus derselben Messung geschrieben),
`BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` (verkörpert),
`BEO-PGC/test-runner-stiller-ausschluss` (offen, 2×, nicht einschlägig —
explizite Phase-Auswahl),
`BEO-PGC/formvorbild-kopie-traegt-deutsches-wortfragment-weiter` (offen,
2×, nicht einschlägig — dieser Slice kopiert kein Form-Vorbild, er
erweitert die Origin).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Fortsetzung der
SDK-Bäume und der `tools/harness/`-Werkzeug-Familie; der Runner ist ein
GF-Artefakt mit bestehender Struktur, die Erweiterung ist additiv).