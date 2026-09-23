# Verifikationsbericht: slice-sdk-python-nats-stream-client-flaeche — 2026-09-23

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?" gegen den
DoD-Vertrag ([Slice-Plan](../plan/planning/in-progress/slice-sdk-python-nats-stream-client-flaeche.md)
§2, alle 15 Zeilen), den Plan-vs-Code-Diff (§3 + Plan-Nachzug + das
committete §3.13-Suchlauf-Feld), die
[`SPEC-024`](../../spec/pflichtenheft.md)-Konformität und
[`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
§Entscheidung Festlegung 2/3. **Nicht** gegen den Diff als solchen
(Reviewer-Aufgabe, abgeschlossen mit
[Haupt-Review](review-slice-sdk-python-nats-stream-client-flaeche.md)
(F-1…F-11) und
[Fixrunden-Report](review-slice-sdk-python-nats-stream-client-flaeche-fixrunde.md)
(R-1…R-4)) und nicht gegen realen Bedarf (Validator, hier nicht
ausgelöst).

**Gegenstand:** `git diff fa71cc96..HEAD` — vier Commits auf `main`:
`0f8cc4f2` (Implementations-Commit: NATS-Fläche + Version-Hebung +
Träger-Nachzug), `b4d1d352` (Fixrunde 1, F-1…F-11), `264498e5`
(ID-Verlinkung im Haupt-Review-Report, gemessen: rein mechanische
Link-Formen), `a3e9d64e` (Fixrunde 2, R-1…R-3, trägt zugleich den
Fixrunden-Report).

**Methode — Belege, nicht Behauptungen** ([`AGENTS.md`](../../AGENTS.md)
§3.12 Instanz B): Jede DoD-Zeile wurde gegen eine eigene Messung gehalten,
nicht gegen den Berichtstext. Testzahlen per `grep -c "^def test_"` über
alle fünf Testdateien; die **Unit-Suite selbst gefahren**, netzlos im
vorhandenen `sdk-python-integration`-Image mit der **aktuellen**
`nats_stream_client.py` als Einzeldatei-Mount (48 passed, Exit 0 — der
Lauf prüft damit den HEAD-Stand, nicht den Image-Stand); `make gates`
**selbst gefahren**, ungepiped, `make`-eigener Exit-Code separat
gesichert ([`AGENTS.md`](../../AGENTS.md) §3.9); der Nachweis-Stempel
**vor** dem eigenen Gates-Lauf gegen den Baum-Hash gehalten
(`bash tools/harness/working-tree-hash.sh | diff -
.harness/state/gates-passed.diffsha` — byte-gleich, Exit 0) und der
Baum-Hash inhaltsbasiert über getrackte **und** untracked Dateien
(Skript gelesen); `make commit-traceability
RANGE=fa71cc96..a3e9d64e` **selbst gefahren** (OK, 4 Commits); der
0.2.0-Wheel **extrahiert** und Feld für Feld gelesen (METADATA,
Abhängigkeiten, Quellstand je Modul, byte-gleich gehalten gegen
`git show 0f8cc4f2:…`); die reale
`make test-sdk-python-integration`-Evidenz gegen beide
Implementer-Logs (`/tmp/integration-s3e.log` 08:07 **rot**,
`/tmp/integration-s3f.log` 08:10 **grün**) an den Traceback-Zeilennummern
und Testzahlen gegen den HEAD-Quellstand gehalten; beide
§3.13-Suchlauf-Stände per `grep` an Parent und HEAD gemessen;
`spec/pflichtenheft.md` §2 [`SPEC-024`](../../spec/pflichtenheft.md),
§1 [`LH-FA-SST-009.a`](../../spec/pflichtenheft.md), §6/§7 und
[`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
im Original gelesen; die serverseitige Gegenprobe
`internal/adapters/driven/natsstream/publisher.go` **nur gelesen**.

---

## 1. DoD je Zeile

| # | DoD-Zeile (§2) | Verdikt |
|---|---|---|
| 1 | [`LH-FA-SST-009`](../../spec/lastenheft.md) erfüllt: öffentliche Client-Klasse, Token-Auth, Vollinhalts-Namensraum, zehn Felder; Unit-Tests referenzieren `SPEC-024` (Schema-Vollständigkeit, Subjekt-Format netzlos, Verbindungs-Fehlerpfad gegen Fake) | **erfüllt** (§2.1) |
| 2 | `pyproject.toml`: NATS-Python-Bibliothek als gepinnte Laufzeitabhängigkeit, real zum Bau gemessen | **erfüllt** (§2.2, eine Anmerkung zur Pin-Form) |
| 3 | Realserver-Integrationstest erweitert, läuft über `make test-sdk-python-integration`, Empfang real belegt, `change_id` gegen `cdc.changes` | **erfüllt** (§2.3) |
| 4 | Version-Hebung `0.1.0` → `0.2.0` | **erfüllt** (§2.4) |
| 5 | Träger-Nachzug `spec/pflichtenheft.md` `LH-FA-SST-009.a`/`SPEC-027` + Suchlauf-Pflicht (Ergebnis committet) | **erfüllt** (§2.5; Residuen der Re-Reviews am HEAD gemessen — gezogen) |
| 6 | Kein Import aus `internal/**`/`cmd/**`/`gen/**` | **erfüllt** (§2.6) |
| 7 | `make gates` grün | **erfüllt** (eigener Lauf, Exit 0, §2.7) |
| 8 | Reales Artefakt-Paar (`make sdk-pack-python`) mit `version=0.2.0` als Smoke-Beleg, alle vier Flächen im selben Artefakt | **erfüllt mit Abweichung V-1** (Artefakt-Frische; §2.8) |
| 9 | Review durchgeführt, Report unter `docs/reviews/` liegt vor | **erfüllt mit Abweichung V-2** (keine Re-Review-Stufe der Fixrunde 2; §2.9) |
| 10 | Doku-Update: [`docs/user/benutzerhandbuch.md`](../../docs/user/benutzerhandbuch.md) gRPC/SSE/NATS gebündelt | **erfüllt** (§2.10) |
| 11 | Closure-Notiz mit Steering-Loop-Lerneintrag | **korrekt offen** (§7 trägt Template; Closure-Rest nach diesem Lauf) |
| 12 | Reconciliation-Register — „entfällt" | **korrekt** (`docs/plan/planning/reconciliation.md` existiert nicht, gemessen) |
| 13 | Beobachtungs-Register fortgeschrieben / „keine Beobachtung" in §7 | **korrekt offen** (Antwort in §7 bei Closure; Register trägt keinen Eintrag dieses Slices) |
| 14 | Jedes Risiko aus §6 trägt einen Ausgang | **erfüllt** (alle vier Ausgänge im Plan geschrieben: 1× am Ort benannt — Modul-Docstring; 2× entfallen mit Begründung; 1× weiter offen → Review) |
| 15 | Drei Paarungen (Anker · Folge-Slice · Register) | **erfüllt als Deklaration** (Welle [welle-sdk-python-vollabdeckung](../plan/planning/welle-sdk-python-vollabdeckung.md) noch offen; Endprüfung regelkonform bei deren Closure) |

Kein DoD-Punkt wurde still gestrichen; die Zeilen 11/13 sind die
regelmäßige Trennung Implementer-Bericht → Verifikation → Closure
(dieselbe Form wie bei den drei Vorgänger-Verifikationen dieser Welle).

### 2.1 DoD 1 — Client-Fläche und Unit-Tests (erfüllt)

- **Öffentliche Klasse:**
  `sdks/python/pgchangefeed/src/pgchangefeed/nats_stream_client.py`
  trägt `PgChangeFeedNatsStreamClient`; `__init__.py` importiert sie und
  führt sie in `__all__` (gemessen, Zeilen 27–36).
- **Verbindungsebene-Auth:** `nats.connect(servers=[address],
  token=api_token, connect_timeout=15, allow_reconnect=False)` (Zeilen
  83–88) — der Token geht in den Connect, nicht per Message
  ([`SPEC-024`](../../spec/pflichtenheft.md) Zeile „Authentifizierung").
  Ein Connect-Fehlschlag wird über `connection_error` gesammelt und am
  Generator **erhoben** wieder ausgelöst (`raise connection_error[0]`,
  Zeilen 111–112) — kein stiller leerer Stream; die Zusage steht im
  Methoden-Docstring („not a swallowed empty stream", Zeilen 70–73) und
  im Unit-Test `test_connect_rejection_surfaces_and_is_not_swallowed`.
- **Subjekt-Form:** `_subject_namespace` baut
  `cdc.stream.<source_id>.>` (Zeilen 127–131) —
  [`SPEC-024`](../../spec/pflichtenheft.md): „`cdc.stream.<source_id>.>`
  deckt alle Tabellen einer Quelle". Exakte String-Assertion im
  Subjekt-Test (zwei Quell-IDs).
- **Zehn Felder, kein drittes Schema:** `StreamChange.from_json`
  liest genau die zehn [`SPEC-024`](../../spec/pflichtenheft.md)-Felder
  (`models.py` Zeilen 268–293); der Test
  `test_parse_stream_change_maps_the_ten_fields` hält alle zehn einzeln
  am JSON-Image. Server-Gegenprobe (nur gelesen):
  `internal/adapters/driven/natsstream/publisher.go`, `streamMessage`
  (zehn `json:`-Tags, identische Reihenfolge) und `subjectPrefix =
  "cdc.stream."` — deckungsgleich.
- **Fire-and-Forget, kein Replay:** Docstring („Delivery guarantee:
  none (Core NATS, fire-and-forget, no replay)"), Granularität „eine
  Nachricht je Zeilen-Change" über einen Yield je Nachricht; kein
  Replay-Pfad im Modul (kein Lesen aus `cdc.changes`).
- **Testzahlen — gemessen, nicht behauptet:** `grep -c "^def test_"`
  über `sdks/python/pgchangefeed/tests/*.py`: `test_http_client.py` =
  **20**, `test_options.py` = **3**, `test_grpc_client.py` = **8**,
  `test_sse_client.py` = **10**, `test_nats_stream_client.py` = **7** —
  Summe **48**, deckungsgleich mit der Vorgänger-Formel (41 bei SSE)
  plus 7. Die pytest-Zahl im Pack-Bau
  (`/tmp/sdk-pack-s3-final.log`, Zeile 227): „**48 passed**, 2 warnings".
  **Eigener, netzloser Suite-Lauf gegen den HEAD-Stand** (Image +
  Einzeldatei-Mount der aktuellen `nats_stream_client.py`): 48 passed,
  Exit 0.
- **Unit-Test-Bindungen (Mutations-Fläche) am Quelltext geprüft:**
  Subjekt-Format (zwei exakte String-Equality-Assertions),
  Verdrahtung (`connect_calls == [([NATS_URL], TOKEN)]` und
  `subscribe_subjects == ["cdc.stream.src-e2e.>"]` — Token und Subjekt
  an der Eingabeseite gebunden), Schema-Violations (drei
  Payload-Mutationen: kein JSON, Array, fehlendes Feld → typisierter
  Fehler), Verbindungs-Fehlerpfad gegen den Fake-Seam
  (`monkeypatch` auf `nats.connect`). Der reale Ablehnungs-Nachweis der
  Ablehnungs-Klasse ist laufgebunden (§2.3).

### 2.2 DoD 2 — Abhängigkeit (erfüllt)

- `pyproject.toml`: `"nats-py>=2"` in `dependencies` — in derselben
  Untergranz-Form wie die drei Geschwister-Abhängigkeiten (`httpx>=0.27`,
  `grpcio>=1.75`, `protobuf>=6`); kein Upper-Pin, kein Lock-Mechanismus
  im Repo — konsistente Package-Form.
- **Real zum Bau gemessen:** das Pack-Bau-Log
  (`/tmp/sdk-pack-s3-final.log`, Zeilen 84–86, 121) löst `nats-py>=2`
  real auf (`Downloading nats_py-2.16.0-py3-none-any.whl`), der
  Integrations-Bau (`/tmp/integration-s3f.log`, Zeilen 137–195)
  installiert `nats-py-2.16.0` — die DoD-Forderung „nicht blind aus den
  Schwester-SDKs übernommen" ist durch die Eigenmessung der Bibliothek
  in diesem Bau erfüllt.

### 2.3 DoD 3 — Realserver-Integrationstest (erfüllt)

- **Testdatei:** `sdks/python/pgchangefeed/integration/test_nats_realserver.py`
  (Pfad-Abweichung ggü. §3 im Plan-Nachzug deklariert — Geschwister-Ordner
  `integration/` seit dem gRPC-Slice).
- **Runner:** dritte Phase im `run_surface_phase`-Muster, Sentinel
  `PythonNatsSdkE2ESentinel`, ID-Bereich 320ff., Reject-Marker
  `REJECTED token-rejected`, Env-Extra-Liste
  (`PGCHANGEFEED_NATS_URL=nats://nats:4222`,
  `PGCHANGEFEED_NATS_STREAM_TOKEN=e2e-nats-stream-token`,
  `PGCHANGEFEED_SOURCE_ID=src-e2e`); Token-Verkettung real geprüft —
  `compose.yaml` Zeile 114 `CDC_NATS_STREAM_TOKEN: e2e-nats-stream-token`
  = Runner-Konstante, `CDC_SOURCE_ID: src-e2e` (Zeile 92) =
  `PGCHANGEFEED_SOURCE_ID`; SQL-Gegenprüfung gegen `cdc.changes`
  (Runner-Zeilen 280–283, count ≥ 1 je `change_id`).
- **Reale Läufe (Implementer-Logs, selbst gelesen):**
  `/tmp/integration-s3e.log` (08:07) ist **rot** — der Reject-Test band
  die falsche Klasse (`pytest.raises(AuthorizationError)`); der
  NATS-Server lehnte real ab, aber am Connect-Init-Pfad fliegt
  `nats.errors.Error` mit dem Server-Rohwortlaut — der Lauf dokumentiert
  den echten Befund der Klasse.
  `/tmp/integration-s3f.log` (08:10) ist **grün**: Endzeile nennt
  **alle drei** Flächen mit je `change_id` — gRPC `807-1`, SSE `810-1`,
  **NATS `812-1`** — „change_id je unabhaengig ueber cdc.changes
  lesbar", „ein NATS-Verbindungsversuch mit falschem Token wurde vom
  NATS-Server abgelehnt". Die s3e→s3f-Sequenz ist zugleich ein lebender
  Mutations-Beweis für die Ablehnungs-Assertion: die Bindungsänderung
  kippte das Ergebnis real rot→grün gegen denselben Server.
- **Laufdeckung zum HEAD-Stand:** der s3f-Lauf lief zwischen Fixrunde 1
  (Arbeitsbaum, vor Commit 08:11) und deren Commit; die
  Traceback-Zeilennummern im Log (`nats_stream_client.py` Zeilen 101 =
  `await nc.drain()`, 103 = `asyncio.run`) decken sich mit dem
  HEAD-Stand (der Prä-F-11-Stand läge bei 102/104), die Unit-Teilmenge
  trägt 7 NATS-Tests. Fixrunde 2 (08:36) änderte danach nur einen
  Typtoken (Zeile 75, keine Zeilenverschiebung, kein Verhalten) — der
  Lauf deckt den HEAD-Stand der NATS-Fläche ab.
- **Feldvollständigkeit am realen Wire:** der Receive-Test hält alle
  zehn Felder je in seinem Wire-Typ (inkl. `old_image is None` für
  INSERT, Non-empty-Identitäten, `operation == "INSERT"`), Sentinel
  real im `new_image`.

### 2.4 DoD 4 — Version-Hebung (erfüllt)

- `pyproject.toml` `version = "0.2.0"` (Zeile 7, gemessen); Wheel-METADATA
  `Version: 0.2.0` (extrahiert aus
  `pgchangefeed-0.2.0-py3-none-any.whl`); PEP-440-Minor-Inkrement
  (Plan-§6-Risiko-4-Ausgang konsistent).

### 2.5 DoD 5 — Träger-Nachzug und Suchlauf (erfüllt)

- `spec/pflichtenheft.md` §1 (Python-Absatz, Zeilen 188–192): „deckt
  HTTP-API, gRPC-Stream, SSE und NATS-Vollinhalt, real Docker-only
  paketierbar … und real geprüft (`pgchangefeed-0.2.0-py3-none-any.whl`,
  `pgchangefeed-0.2.0.tar.gz`)" — gemessen; §6 Vertragszeile
  [`SPEC-027`](../../spec/pflichtenheft.md) (Zeile 625): „aktuell
  `0.2.0`" + Vier-Verträge-Liste in der Zeilenform der Schwester-Zeilen
  (624/626); §7 Historie (Zeile 672): Chronik-Zeile datiert 2026-09-23
  im Muster der C#-/Kotlin-Vorgänger.
- **Suchlauf-Feld committet** (Plan-Zeilen 140–155, vor der
  Closure-Notiz) — der F-5/R-1-Streit ist am HEAD aufgelöst; **eigene
  Messung der genannten Fundstellen am HEAD:**
  `README.de.md:27` trägt „dieselben vier Zustellwege zusätzlich als
  offizielle Python-Client-Bibliothek" (R-1 gezogen);
  `harness/README.md` sdk-pack-Zeile trägt 0.2.0 mit Endstand-Anker,
  `test-sdk-python-integration`-Zeile trägt drei Flächen **und** drei
  Reject-Formen („… bzw. der laut ablehnenden
  NATS-Verbindungsablehnung („Authorization Violation")"), die
  widersprüchliche Endklause ist entzogen (R-2 gezogen);
  Plan-Zeilen 148/149 sind wahr (Befund- und Behandlungs-Zellen
  decken sich mit dem Baum).
- **Eigener Suchlauf-Grep am HEAD** (`grep -rn "deckt HTTP-API\|
  pgchangefeed-0.1.0\|no gRPC/SSE/NATS surface exists yet" spec/ docs/
  sdks/python/`): die Resttreffer sind sämtlich historisch korrekt —
  die korrigierte Vier-Wege-Form selbst (§1/§6), §7-Historie-Zeile
  664 (0.1.0-Namen, Chronik), die Welle-Plan-Trigger-Zeile 56
  (Trigger-Ist-Stand 2026-09-21, im Plan-Feld als „belassen"
  deklariert und verifiziert), `done/`-Records und
  `docs/reviews/**`. Kein lebender Träger trägt eine stale Form.

### 2.6 DoD 6 — Import-Grenze (erfüllt)

- Eigener Grep über alle drei neuen Python-Dateien: die einzige
  `internal/`-Nennung ist die Docstring-Prosa (Server-Gegenprobe „gelesen,
  nicht importiert"); Imports ausschließlich Python-Standardbibliothek,
  `nats`/`nats.errors` und `pgchangefeed.*`. Kein `examples/python/`
  (Verzeichnis gemessen — existiert nicht); keine Package-Spaltung
  (ein Package, Version `0.2.0`).

### 2.7 DoD 7 — `make gates` (erfüllt, eigener Lauf)

- **Eigener Lauf:** `make gates > /tmp/v-make-gates.log 2>&1`,
  `make`-eigener Exit-Code direkt gesichert (kein Pipe-/Wrapper-Glied
  dazwischen, [`AGENTS.md`](../../AGENTS.md) §3.9): **Exit 0**. Alle
  sechs Gates sichtbar im Log: `baseline-verify: v6.9.0 OK — 54
  Dateien`; `d-check: 917 Datei(en) geprüft, 0 Befund(e)` (Modul-Bündel
  inkl. links/anchors/ids/matrix/versions/structure/hostpaths/tracked);
  `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs
  ohne Struktur-ID` (der Range schließt alle vier Slice-Commits ein);
  `coverage-gate: OK — Coverage 82.80% erfüllt Schwelle 80%`;
  `generated-sync: OK — das committete Erzeugnis ist byte-gleich …`;
  `a-check … gesamt: 0 Befund(e)`.
- **Stempel:** vor dem eigenen Lauf byte-gleich gehalten
  (`7c729d5a55877a6b3e97c48c286be470979352d94c49057fa308ef8369ddeade`
  = `.harness/state/gates-passed.diffsha`, Skript-Exit 0); der Stempel
  trägt den letzten offiziellen Lauf (08:36, nach Fixrunde 2).
- **Slice-Range-Traceability:** `make commit-traceability
  RANGE=fa71cc96..a3e9d64e` (eigener Lauf): OK — 4 Commits, Betreffs
  ohne Struktur-ID, je Betreff
  [`LH-FA-SST-009`](../../spec/lastenheft.md)/[`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md).
  (Die rollen-generischen Sensor-Namen `doc-commits`/`doc-immutable`
  existieren in diesem Repo nicht; hier trägt
  [`ADR-0045`](../plan/adr/0045-commit-traceability-standing-gate.md)s
  Standing-Gate mit `RANGE`-Überschreibung dieselbe Prüfung.)

### 2.8 DoD 8 — Artefakt-Paar (erfüllt mit Abweichung V-1)

- **Paar real vorhanden:** `sdks/python/dist/` trägt
  `pgchangefeed-0.2.0-py3-none-any.whl` (20481 Bytes) und
  `pgchangefeed-0.2.0.tar.gz` (24639 Bytes), Stempel 2026-09-23 07:47,
  neben dem älteren 0.1.0-Paar (06:41); Pack-Log endet mit „Successfully
  built /out/pgchangefeed-0.2.0.tar.gz / …whl". **Alle vier Flächen im
  selben Artefakt** (Wheel-Extraktion: `http_client.py`, `grpc_client.py`
  + `grpc_gen`, `sse_client.py`, `nats_stream_client.py`) — der
  Buchstabe der DoD-Zeile ist erfüllt.
- **Abweichung V-1 — Artefakt-Frische:** das 0.2.0-Paar ist byte-gleich
  dem `0f8cc4f2`-Quellstand (Wheel-`nats_stream_client.py` == `git show
  0f8cc4f2:…`, byte-gleich gemessen) und **nicht** dem HEAD nach den
  zwei Fixrunden: das Wheel trägt noch `# noqa: BLE001` (Zeile 89 —
  die §3.2-Form, die der Quelltext seit Fixrunde 1 nicht mehr trägt),
  `queue.Queue[bytes | None]` (Zeile 75, seit Fixrunde 2 entzogen) und
  ein `exceptions.py` ohne die [`SPEC-024`](../../spec/pflichtenheft.md)-Nennung.
  Kein Neubau nach 08:11 (dist-Stempel 07:47). Die funktionale Aussage
  des Smoke-Belegs (Version + vier Flächen in einem Artefakt) hält;
  die Byte-Frische nicht. **Maßnahme vor der Closure:** entweder ein
  erneuter `make sdk-pack-python` (die Testfläche des Endstands ist
  bereits grün belegt — 48 passed im 08:10er-Integrations-Bau gegen den
  fixierten Quellstand und in meinem eigenen netzlosen Lauf), oder die
  Staleness wird in §7 explizit deklariert. Kein Tag-Push ist erfolgt
  (`sdks/python/dist/` ist `.gitignore`t, PyPI-Ist-Stand bleibt
  `0.1.0`), die Abweichung ist also Evidenz-Frische, keine
  Veröffentlichungs-Lüge.

### 2.9 DoD 9 — Review (erfüllt mit Abweichung V-2)

- Zwei Reports liegen unter `docs/reviews/`:
  [Haupt-Review](review-slice-sdk-python-nats-stream-client-flaeche.md)
  (F-1…F-11, Merge-blockierend ja) und
  [Fixrunden-Report](review-slice-sdk-python-nats-stream-client-flaeche-fixrunde.md)
  (R-1…R-4, Merge-blockierend ja).
- **Abweichung V-2 — keine Re-Review-Stufe der Fixrunde 2:** Commit
  `a3e9d64e` trägt die Fixes zu R-1/R-2/R-3 und committet zugleich den
  Re-Review-Report der **ersten** Fixrunde; eine dritte Reviewer-Stufe
  über die Fixrunde-2-Fixes selbst existiert nicht. Diese Lücke schließt
  diese Verifikation **inhaltlich** — alle drei Residuen am HEAD selbst
  gemessen und gezogen: `README.de.md:27` (vier Zustellwege),
  `harness/README.md` test-sdk-Zeile (drei Reject-Formen, Endklause
  entzogen), `nats_stream_client.py:75` (`queue.Queue[bytes]`). R-4
  (Suchlauf-Raum-Kopf) ist am HEAD erweitert. Für den
  Planner bleibt als Prozess-Notiz: die zweite Fixrunde lief ohne
  Reviewer-Gegenprüfung; das DoD-Kriterium in seiner Wortform („Report
  unter `docs/reviews/` liegt vor") ist durch die zwei Reports erfüllt.

### 2.10 DoD 10 — Doku-Update (erfüllt)

- [`docs/user/benutzerhandbuch.md`](../../docs/user/benutzerhandbuch.md):
  `Version: 1.43` (Kopf), neuer §4-Absatz im NATS-Abschnitt
  (`PgChangeFeedNatsStreamClient.stream_changes()`, Subjekt-Namensraum,
  zehn Felder, verbindungsseitige Auth, `pip install pgchangefeed`),
  Änderungshistorie-Zeile 1.43 datiert 2026-09-23 mit Slice-Kennung; die
  gRPC-(1.40/1.41)- und SSE-(1.42)-Anteile tragen die Vorgänger-Slices —
  die Historie-Zeile 1.43 behauptet die vollständige Matrix, und die
  drei Absätze stehen real in den drei Abschnitten.

---

## 3. Plan-vs-Code-Diff

Alle sieben §3-Zeilen und alle vier Plan-Nachzug-Zeilen sind im Diff
(`fa71cc96..HEAD`) vertreten; keine Zeile still gestrichen:

| §3-Zeile | Im Diff | Gemessen |
|---|---|---|
| `nats_stream_client.py` (neu) | ✓ | 159 Zeilen, `0f8cc4f2` |
| `pyproject.toml` (update) | ✓ | `nats-py>=2` + `version = "0.2.0"` + Beschreibung Vier-Wege |
| `tests/test_nats_stream_client.py` (neu) | ✓ | 7 Tests, netzlos, Fake-Seam |
| Integrationstest (neu) | ✓ | als Abweichung deklariert (Plan-Nachzug Zeile 1, `integration/`) |
| `options.py` (update) | ✓ | Docstring trägt Vier-Wege-Matrix, `SPEC-024`-Form |
| `spec/pflichtenheft.md` §1/§6 (update) | ✓ | nach Fixrunde 1 vollständig (F-1) |
| `docs/user/benutzerhandbuch.md` (update) | ✓ | §2.10 |

Plan-Nachzug: `integration/test_nats_realserver.py` (Abweichung
deklariert) ✓; `tools/harness/run-sdk-python-integration-tests.sh`
(dritte Phase, `run_surface_phase`-Refaktorisierung mit Env-Extra-Liste,
Sentinel 320ff., Reject-Marker, SQL-Gegenprüfung — die SQL-Gegenprüfung
war Bestand der Vorgänger-Phasen und bleibt je Phase wirksam) ✓;
`__init__.py` + `options.py` + `sdks/python/README.md` + pyproject-
Beschreibung ✓; `spec/pflichtenheft.md` §7 Chronik-Zeile ✓.

Die Plan-eigenen Änderungen im Diff (Plan-Nachzug-Tabelle, §3.13-Teil
als committetes Feld, die Fixrunden-Korrekturen der Zeilen 148/149 und
des Feld-Kopfs) sind selbst Teil des Diffs und oben je gemessen. Die
still gestrichene §3-Zeile des Erstzugs (F-1, §6-`SPEC-027`) ist seit
Fixrunde 1 als ausgeführt nachgezogen (Commit-Message + Pflichtenheft-
Zeile 625 gemessen). Keine weitere still gestrichene Zeile.

---

## 4. `SPEC-024`-Konformität (Feld für Feld)

| Merkmal (§2 [`SPEC-024`](../../spec/pflichtenheft.md)) | Client-Seite | Server-Gegenprobe (nur gelesen) | Beleg |
|---|---|---|---|
| Subjekt-Schema `cdc.stream.<source_id>.<schema>.<table>`; `cdc.stream.<source_id>.>` deckt alle Tabellen einer Quelle | Abonnement exakt `cdc.stream.<source_id>.>` (`_subject_namespace`), `<schema>.<table>` kommen je Nachricht | `subjectPrefix = "cdc.stream."`, Subjekt je Tabelle | Unit-Test exakt, Integrationstest real |
| Nachrichteninhalt: zehn Felder wie `SPEC-021` | `StreamChange.from_json` liest genau die zehn Keys | `streamMessage`: dieselben zehn `json:`-Tags | Unit-Test Feld für Feld; Realserver-Test je Wire-Typ |
| Granularität: eine Nachricht je Zeilen-Change | ein `StreamChange` je Message, kein Dedup im Client | `Run`-Loop veröffentlicht je Kanal-Change | s3f-Log (INSERT je Sentinel-Zeile) |
| Zustellgarantie keine (Fire-and-Forget, kein Replay) | Docstring + kein Replay-Pfad; Catch-up bleibt beim HTTP-Lesezugriffsweg | Publisher ohne Zustellgarantie | Modul-Docstring, Runner-Fire-and-Forget-Fenster |
| Auth Verbindungsebene; Ablehnung laut | Token im `nats.connect` (nicht per Message); Connect-Fehler wird erhoben ausgelöst, nicht geschluckt | Server verlangt Token (`--auth` in `compose.yaml`), Feed trägt `CDC_NATS_STREAM_TOKEN` | Unit-Test (Fake), Realserver-Reject-Test mit Ursachen-Bindung `nats_errors.Error, match="Authorization Violation"` + Runner-Marker; s3e→s3f lebend rot→grün |
| Aktivierung (beide Bedingungen, Server-Seite) | Client konfiguriert nur `address`/`api_token` | `compose.yaml` trägt beide Werte (`CDC_NATS_URL` + `CDC_NATS_STREAM_TOKEN`) | compose.yaml Zeilen 105/114 (nur gelesen) |

Kein drittes Schema: der Client instanziiert denselben
`StreamChange`-Typ wie die SSE-Fläche (`type(change) is StreamChange`
in zwei Tests).

---

## 5. `ADR-0110` Festlegung 2/3 (erfüllt)

- **Festlegung 2/Folgepflicht 1 — reale Server-Instanz je Fläche:** die
  dritte Runner-Phase mit drei Beleg-Klassen (Empfang am Wire,
  SQL-Gegenprüfung, Ablehnung) — die drei Phasen des Runners decken
  gRPC/SSE/NATS je mit eigenem Sentinel/ID-Bereich (300/310/320) und
  je eigener Reject-Form; der grüne Lauf ist im
  `/tmp/integration-s3f.log` belegt (Erfolgszeile mit `change_id=812-1`
  für die NATS-Fläche, geprüft).
- **Festlegung 2 — kein `examples/python/`:** Verzeichnis gemessen
  (existiert nicht), Diff trägt keinen solchen Pfad.
- **Festlegung 3 — keine Package-Spaltung:** ein Package (`pgchangefeed`
  `0.2.0`), kein zweites; die Struktur-Entscheidung ist im
  `__init__.py`-Docstring an die richtige Schicht verankert (Welle-Plan
  §6 Zeile 137 „diese Welle wählt v2 desselben Packages" — Anker im
  Original aufgeschlagen; [`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
  Festlegung 3 delegiert, entscheidet nicht).
- **Import-Grenze** (Fitness-Function-Zeile der ADR): §2.6.

---

## 6. Aufgedeckte Abweichungen (Verifier-Funde)

**Zwei** Abweichungen, beide Evidenz-/Prozess-Klasse, keine berührt
Code-Korrektheit, [`SPEC-024`](../../spec/pflichtenheft.md)-Draht oder
Sicherheit:

- **V-1 (Artefakt-Frische, Fundstelle `sdks/python/dist/` vs.
  HEAD-Quellstand):** das 0.2.0-Artefakt-Paar (07:47) ist byte-gleich
  dem `0f8cc4f2`-Stand, nicht dem Endstand nach zwei Fixrunden — im
  Wheel: `# noqa: BLE001` (`nats_stream_client.py:89`),
  `queue.Queue[bytes | None]` (:75), `exceptions.py` ohne
  [`SPEC-024`](../../spec/pflichtenheft.md)-Nennung. DoD-Wortlaut im
  Buchstaben erfüllt; Maßnahme vor Closure (§2.8).
- **V-2 (keine Re-Review-Stufe der Fixrunde 2, Fundstelle
  `a3e9d64e`/`docs/reviews/`):** die zweite Fixrunde lief ohne
  Reviewer-Gegenprüfung; inhaltlich durch diese Verifikation geschlossen
  (alle drei Residuen am HEAD gezogen, gemessen); Prozess-Notiz an den
  Planner (§2.9).

---

## 7. Offene Restpunkte für die Closure (erwartete Trennung, keine Abweichung)

- DoD-Checkboxen (alle `[ ]`) und §7 Closure-Notiz inkl.
  Steering-Loop-Lerneintrag — bei Closure nachzuziehen (Muster der
  Vorgänger-Slices); die Finding-Klassen der beiden Review-Reports
  (dort enumeriert) gehen in den Steering-Loop-Zähler.
- Beobachtungs-Register-Antwort in §7 („keine Beobachtung angefallen"
  ist ebenfalls eine Antwort); die Fixrunden-Funde der Klasse
  `arbeit-ueberholt-stehenden-traeger` (F-4, R-1/R-2) sind Kandidaten
  für den bestehenden Register-Eintrag.
- Risiken §6: Ausgänge stehen; der Risiko-1-Ausgang ist am Ort benannt
  (Modul-Docstring, „one iterator form across all four surfaces").
- Nach Closure-Commits: frischer Gate-Lauf + neuer Stempel (der
  aktuelle Stempel deckt den Stand vor der Closure-Notiz).

---

## 8. Verdikt

**DoD je Zeile:** 11 von 15 Zeilen erfüllt (davon 2 mit Abweichung
V-1/V-2), 4 Zeilen korrekt auf die Closure getrennt (Closure-Notiz,
Beobachtungs-Register, Paarungen-Endprüfung bei Welle-Closure,
Reconciliation-korrekt-entfällt).

**Aufgedeckte Abweichungen:** 2 (V-1 Artefakt-Frische — Fundstelle
Wheel-Inhalt vs. HEAD-Quellstand, Maßnahme vor Closure; V-2 — fehlende
Re-Review-Stufe der Fixrunde 2, inhaltlich geschlossen, Prozess-Notiz).

**Gesamt-Urteil:** Die DoD-Substanz des Slices ist real und belegt:
die [`SPEC-024`](../../spec/pflichtenheft.md)-Client-Fläche trägt den
Draht Feld für Feld gegen Spec und Server-Gegenprobe, die Realserver-
Evidenz ist echt (roter Fehlschlag der falschen Bindung → grüner Lauf
mit `change_id=812-1` gegen `cdc.changes`), 48 Unit-Tests grün gegen
den HEAD-Stand in meinem eigenen netzlosen Lauf, Version `0.2.0` in
pyproject und realem Wheel, `make gates` grün in meinem eigenen,
ungepiped geprüften Lauf mit byte-gleichem Stempel, Träger-Nachzug
vollständig und am HEAD verifiziert. Die zwei Abweichungen sind vor der
Closure zu tragen (V-1 durch Neubau oder Deklaration, V-2 als
Prozess-Notiz an den Planner); sie blockieren die Substanz nicht, aber
V-1 sollte vor der Closure-Notiz real aufgelöst werden, weil
`harness/README.md` und `spec/pflichtenheft.md` das 0.2.0-Paar als
Endstand-Anker nennen.