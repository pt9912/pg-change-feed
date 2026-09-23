# Verifikationsbericht: slice-sdk-python-grpc-client-flaeche — 2026-09-23

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?" gegen den
DoD-Vertrag (`Slice-Plan` `slice-sdk-python-grpc-client-flaeche`
§2), den Plan-vs-Code-Diff (§3 + Plan-Nachzug inkl. Fixrunden-Zeilen) und
die ADR-Konformität ([`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md),
[`ADR-0107`](../plan/adr/0107-python-pypi-zweites-sdk-package.md)).
**Nicht** gegen den Diff als solchen (Reviewer-Aufgabe, abgeschlossen mit
[Review-Report](review-slice-sdk-python-grpc-client-flaeche.md)) und nicht
gegen realen Bedarf (Validator, hier nicht ausgelöst).

**Gegenstand:** `git diff a13c1d03..HEAD` — fünf Commits auf `main`:
`f1ce9be4` (Implementations-Commit), `b6f82269` (Review-Report),
`b7a993fe` (Fixrunde F-1…F-4), `f4112fe9`/`3b99fa8b` (Plan-Doku
Roadmap — nicht Slice-Substanz). Die DoD-Substanz liegt in `f1ce9be4` +
`b7a993fe`.

**Methode — Belege, nicht Behauptungen:** Jede DoD-Checkbox-Zeile wurde
gegen eine eigene Messung gehalten, nicht gegen den Berichtstext:
`grep -c "^def test_"` über alle vier Testdateien, Wheel-/SDist-Inhalt
entpackt und byte-gleich gegen den committeten Quelltext gehalten,
Import-Zeilen über `src`/`tests`/`integration` vollständig aufgelistet,
`find examples -iname "*python*"`, Target-Existenz in `harness/mk/sdk.mk`
und `harness/README.md` §Werkzeuge, die drei Lauf-Log-Dateien der
Fixrunde (Host-Temp-Verzeichnis, `sdk-pack-fixrunde.log`,
`sdk-python-integration-fixrunde.log`, `gates-fixrunde.log` — transient,
ihre belastbaren Zeilen sind hier zitiert) und der Gate-Nachweis-Stempel
`.harness/state/gates-passed.diffsha` gegen einen frisch berechneten
Baum-Hash (`tools/harness/working-tree-hash.sh`) gehalten — der Stempel
deckt den exakten Ist-Baum (`3b99fa8b`, Arbeitsbaum sauber). Kein
erneuter realer Integrationstest-Lauf (Auftragsgrenze); die
Runner-Assertionen wurden statisch gegen den Skripttext geprüft.

---

## 1. DoD je Zeile

| # | DoD-Zeile | Behauptet | Eigene Messung | Verdikt |
|---|---|---|---|---|
| 1 | `LH-FA-SST-009` gRPC-Fläche (`[x]`) | Sensor-Beleg im Plan: 29 Unit-Tests, 6 neue gRPC-Tests | erfüllt **mit Abweichung V-1** (Zahl veraltet, s. u.) | Substanz erfüllt |
| 2 | Realserver-Integrationstest-Werkzeug (`[x]`) | `EXIT=0`, change_id=804-1, Unauthenticated | erfüllt | erfüllt |
| 3 | Kein Import aus `internal/**`/`cmd/**`/`gen/**` (`[x]`) | strenge Import-Prüfung ohne Treffer | erfüllt (0 Content-Treffer) | erfüllt |
| 4 | `make gates` grün (`[x]`) | Stempel frisch | erfüllt | erfüllt |
| 5 | Review-Report (`[ ]`) | — | korrekt offen (Fixrunde existiert) | Trennung korrekt |
| 6 | Doku-Update (`[ ]`) | — | korrekt offen (gebündelt); README-Zeile im Zug erledigt | Trennung korrekt |
| 7 | Closure-Notiz (`[ ]`) | — | korrekt offen (§7 leer) | Trennung korrekt |
| 8 | Reconciliation (`[ ]`, „entfällt") | — | `docs/plan/planning/reconciliation.md` existiert nicht | korrekt |
| 9 | Beobachtungs-Register (`[ ]`) | — | korrekt offen (Antwort in §7 bei Closure) | Trennung korrekt |
| 10 | Risiko-Ausgänge (`[ ]`) | — | Ausgänge in §6 geschrieben (2× weiter offen, 1× entfallen) | Trennung korrekt |
| 11 | Drei Paarungen (`[ ]`) | — | [Welle](../plan/planning/done/welle-sdk-python-vollabdeckung.md) offen | Trennung korrekt |

### DoD 1 — gRPC-Client-Fläche (Substanz erfüllt, Beleg-Zahl veraltet)

- **Öffentliche Klasse:** `sdks/python/pgchangefeed/src/pgchangefeed/grpc_client.py`
  trägt `PgChangeFeedGrpcClient` (`grpc.Channel` injiziert,
  [`ClientOptions`](../../spec/pflichtenheft.md)-Konfiguration) und wird
  über `__init__.py` (`from pgchangefeed.grpc_client import
  PgChangeFeedGrpcClient`, Eintrag in `__all__`) Package-sichtbar exportiert.
- **RPC-Form:** `self._client.StreamChanges(_StreamChangesRequest(),
  metadata=metadata, timeout=timeout)` mit `metadata =
  (("authorization", "Bearer <token>"),)` — Bearer-Form am Ort,
  [`SPEC-020`](../../spec/pflichtenheft.md)-konform.
- **Zehn Felder:** `tests/test_grpc_client.py` hält die Feldliste
  (`SPEC_020_FIELD_NAMES`, zehn Namen) gegen den zur Bauzeit erzeugten
  Stub-Deskriptor (`Change.DESCRIPTOR.fields`) — Namen **und**
  Feldnummern 1–10 in Wire-Reihenfolge, gebunden an
  `proto/cdc/stream/v1/changestream.proto` (Proto-Drift bricht den Bau).
  Der Integrationstest hält dieselben zehn Felder zusätzlich typbewusst
  (`isinstance` je Feld, Nicht-`bool`-Guard) am realen Wire-Image.
- **Testzahlen — gemessen, nicht behauptet:** `grep -c "^def test_"`
  liefert `tests/test_grpc_client.py` = **8**,
  `tests/test_http_client.py` = 20, `tests/test_options.py` = 3 →
  **31 Unit-Tests**. Der Pack-Lauf der Fixrunde belegt im Log
  `============================== 31 passed in 0.12s
  ==============================` — 31, nicht die im DoD-Beleg
  behaupteten 29 mit 6 neuen gRPC-Tests. **Abweichung V-1** (Abschnitt 5).
- **Stub aus dem Bau, nicht committet:** `sdks/python/.gitignore`
  schließt `pgchangefeed/src/pgchangefeed/grpc_gen/*.py` aus und negiert
  `__init__.py`; `git ls-files sdks/python` zeigt in `grpc_gen/`
  ausschließlich den Package-Marker — kein Stub im committeten Baum.
  Der Bau liest die `.proto` über `COPY --from=proto
  cdc/stream/v1/changestream.proto changestream.proto`
  (`sdks/python/Dockerfile`) — der Aufruf trägt `--build-context
  proto=proto` (`tools/harness/sdk-pack-python.sh`, Log-Zeile
  `#2 [context proto]`).
- **Artefakte tragen die Stubs:** Wheel und SDist in
  `sdks/python/dist/` (Bauzeitstempel der Fixrunde) enthalten
  `pgchangefeed/grpc_gen/changestream_pb2.py`,
  `pgchangefeed/grpc_gen/changestream_pb2_grpc.py` und den
  Package-Marker; die Wheel-Module `grpc_client.py`, `__init__.py`,
  `options.py`, `grpc_gen/__init__.py` sind **byte-identisch** mit dem
  committeten Quelltext (entpackt und verglichen) — die
  `dist/`-Artefakte entsprechen dem Ist-Baum.
- **Authn-Boundary netzlos:** `test_unauthenticated_surfaces_from_the_enumeration`
  prüft gegen einen Fake-Channel (kein Socket), dass `UNAUTHENTICATED`
  aus der Enumeration hochkommt und nicht verschluckt wird —
  [`SPEC-020`](../../spec/pflichtenheft.md) Negative am Unit-Rand.

### DoD 2 — Realserver-Integrationstest-Werkzeug (erfüllt)

- **Alle vier Bausteine real existierend:** Skript
  `tools/harness/run-sdk-python-integration-tests.sh` (270 Zeilen,
  `set -euo pipefail`), Target `test-sdk-python-integration` in
  `harness/mk/sdk.mk` (`.PHONY`, Zeile 93/94), Docker-Stufe
  `integration` in `sdks/python/Dockerfile` (`FROM build AS
  integration`, ENTRYPOINT `python -u -m pytest integration -v
  --capture=no`), Testdatei `sdks/python/pgchangefeed/integration/
  test_grpc_realserver.py` (2 Tests).
- **Lauf-Beleg:** Die Log-Datei der Fixrunde trägt als letzte Zeile den
  Runner-Erfolgssatz — mit **interpolierter** `change_id=804-1`, also
  aus dem realen `RECEIVED`-Marker des Testcontainers extrahiert, nicht
  statisch behauptet. Dem Erfolgssatz gehen im Skript fünf explizite
  Fehl-Ausgänge voraus (kein READY, keine RECEIVED-Zeile mit
  `table=feed_e2e_full operation=INSERT … Sentinel`, kein
  `REJECTED code=Unauthenticated`, Test-Exit != 0, `change_id` nicht
  über `cdc.changes` lesbar — `SELECT count(*) FROM cdc.changes …
  change_id = '$grpc_change_id' AND table_name = '$TEST_TABLE' AND
  new_data->>'name' = '$TEST_SENTINEL'`), und das Skript läuft unter
  `set -euo pipefail`: Der Erfolgssatz ist nur nach allen fünf
  Bestehungen erreichbar. Das Log zeigt die Compose-Bring-up-Kette
  (PostgreSQL healthy, Feed-Container, Schema-Rollout über d-migrate,
  Vorbedingungen, Bild-Bau der `integration`-Stufe mit `[context
  proto]`), keine Fehlerzeile.
- **Grenze, sauber benannt:** Das Log trägt keinen expliziten
  `EXIT`-Marker und keine pytest-Ausgabe des Integrationstests — der
  Testcontainer läuft detached, seine Marker liest der Runner über
  `docker logs`; auf Erfolg druckt der Runner nur den Schlusssatz. Der
  Beleg ist die **Skriptstruktur** (Fehl-Ausgänge + `set -e`) plus der
  interpolierte `change_id`-Wert, nicht ein wörtliches `EXIT=0` im Log.
- **Kein stiller Ausschluss (Negativprobe):** `testpaths = ["tests"]`
  in `pyproject.toml` sammelt nur Unit-Tests (31 im Pack-Lauf — die 2
  Integrationstests fehlen dort korrekt), die `integration`-Stufe ruft
  ihren Pfad explizit — dieselbe Trennung, die der Plan-Nachzug
  deklariert.

### DoD 3 — Import-Grenze (erfüllt)

Alle 55 `from`/`import`-Zeilen unter `sdks/python/pgchangefeed/src`,
`tests` und `integration` aufgelistet: ausschließlich Standardbibliothek,
`grpc`, `httpx`, `pytest` und `pgchangefeed.*` — kein Import aus
`internal/**`, `cmd/**`, `gen/**` dieses Repos.

**Messnotiz zur Grep-Form:** Die in der Prüf-Anweisung genannte Form
(`grep -rnE "^\s*(from|import)" … | grep -E "internal|cmd/|gen/"`)
liefert **genau einen Treffer** — pfadgetrieben: der Dateiname
`grpc_gen/__init__.py` selbst matcht `gen/`, der Zeileninhalt
(`importable as ``pgchangefeed.grpc_gen``)`) nicht. Der Package-Name ist
das beabsichtigte Stub-Ziel, keine Verletzung. Die DoD-Formulierung im
Plan (Zweit-Grep `internal/cmd/gen` als Literalausdruck) liefert keinen
Treffer; die Content-only-Messung (Pfadanteil entfernt) liefert null
Treffer. Die DoD-Substanz trägt in allen Formen.

### DoD 4 — `make gates` grün (erfüllt)

Der Nachweis-Stempel `.harness/state/gates-passed.diffsha` wurde gegen
einen frisch berechneten Baum-Hash (`bash tools/harness/working-tree-hash.sh`)
gehalten — **byte-identisch** (`68be27fe…ef5bd7c`), der Stempel deckt
den exakten Ist-Baum (`3b99fa8b`, Arbeitsbaum sauber). Das Gates-Log
trägt alle sechs Einzelbelege:

| Gate | Log-Zeile |
|---|---|
| `baseline-verify` | `v6.9.0 OK — 54 Dateien (Integritaet + Vollstaendigkeit, netzlos)` |
| `docs-check` | `d-check: 908 Datei(en) geprüft, 0 Befund(e)` |
| `commit-traceability` | `OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` |
| `coverage-gate` | `coverage-gate: OK — Coverage 82.80% erfüllt Schwelle 80%` |
| `generated-sync` | `OK — das committete Erzeugnis ist byte-gleich der Ausgabe des gepinnten Generators (Stufe proto-export)` |
| `a-check` | `gesamt: 0 Befund(e)` |

`record-gates` läuft als letztes `gates`-Prerequisite nur bei grünen
Gates — der frische, zum Ist-Baum passende Stempel belegt den grünen
Gesamtlauf über den finalen Stand.

### DoD 5–11 — korrekt getrennt (offen, wo die Rollen-Sequenz noch nicht lief)

- **Review (`[ ]`):** Ein Review lief (Report commit), aber sein Verdikt
  verlangte eine Fixrunde (2 HIGH, 2 MEDIUM); der Report selbst führt
  die Checkbox als offen bis zum Nachzug nach der Fixrunde. Die
  Fixrunde (`b7a993fe`) lief, ein Re-Review liegt nicht vor — `[ ]` ist
  regulär korrekt, kein stiller Selbst-Nachzug.
- **Doku-Update (`[ ]`):** Die Bündelungs-Entscheidung (Handbuch-Teil
  für SSE/NATS + `spec/pflichtenheft.md`-Nachzug im letzten Flächen-Slice)
  steht im Plan; der in diesem Zug erledigte Teil
  (`harness/README.md` §Werkzeuge-Zeile, Nachzug auf der
  `make sdk-pack-python`-Zeile) ist real im Diff und im DoD-Text
  notiert. Der Fixrunden-Handbuch-Korrekturteil (F-3) ist als
  Plan-Nachzug-Zeile deklariert — trägt aber Abweichung V-2.
- **Closure-Notiz, Beobachtungs-Register, Risiko-Ausgänge, Drei
  Paarungen:** alle noch offen, korrekt als `[ ]` — §7 ist leer, die
  Risiko-Ausgänge stehen in §6 („weiter offen" ×2, „entfallen" ×1) und
  werden erst mit §7 geschlossen. Reconciliation: „entfällt" ist
  korrekt — die Datei existiert nicht.

---

## 2. Plan-vs-Code

**§3-Tabelle (8 Zeilen) — vollständig im Diff vertreten:**
`grpc_client.py` (neu), `pyproject.toml` (update:
`grpcio>=1.75`, `grpcio-tools` im `test`-Extra, `protobuf>=6`),
`sdks/python/Dockerfile` (Proto-Kontext-Zeile, Stub-Erzeugung inkl.
`sed`-Aufhebung des flachen Imports, Stufen `pack`/`pack-export`/
`integration`), `tests/test_grpc_client.py` (neu),
`integration/test_grpc_realserver.py` (neu, am Nachzug-Ort),
`tools/harness/run-sdk-python-integration-tests.sh` (neu),
`harness/mk/sdk.mk` (Target), `harness/README.md` §Werkzeuge (Zeile
`make test-sdk-python-integration` real vorhanden, Zeile 160, mit
Bindung-Spalte „kein Gate"). Keine still gestrichene Planzeile, keine
Nicht-Realisierung.

**Plan-Nachzug (10 Zeilen) — vertreten:** `integration/`-Ordner statt
`tests/integration/` (real am Nachzug-Ort, `testpaths`-Trennung trägt),
`protobuf>=6` (mit Begründung am Ort im `pyproject.toml`),
`grpc_gen/__init__.py` (Package-Marker, committet), `sdks/python/.gitignore`
(Stub-Ausschluss + Negation), `tools/harness/sdk-pack-python.sh`
(`--build-context proto=proto` zwingend, Log belegt den Kontext),
`stream_changes(timeout=None)` (real im Client, Fake zeichnet den
Timeout je Aufruf, zwei Tests binden Weiterleitung und Default),
`sdks/python/README.md` §Status + `__init__.py`-Docstring (real
nachgezogen), `docs/user/benutzerhandbuch.md` (F-3-Korrektur real,
Version 1.39, Historie-Zeile), Timeout-Tests (2 neu), typbewusste
Feldprüfung (real: `isinstance` je Feld + Inhalt: Operation, Alt-Bild
leer am INSERT, Sentinel im Neubild, Tabelle, Schema).

**Undeklarierte Änderung — Abweichung V-3:** `options.py` (Docstring +
Attribut-Doku, 21 Zeilen im Implementations-Commit) steht **nicht** als
eigene Plan-Nachzug-Zeile. Die Welle-Plan-Pflicht (Suchlauf prüft
`options.py`s Docstring in jedem der drei Slices) und die
Nachzug-Zeile-7-Begründung nennen `options.py` nur als
Klassen-Referenz. Die Änderung selbst ist korrekt und derselben
Träger-Nachzug-Klasse zugehörig (Satzform „no gRPC/SSE/NATS surface
exists yet" wäre durch diesen Slice falsch) — nur die Deklaration fehlt.

**Nicht-Slice-Substanz korrekt getrennt:** `f4112fe9`/`3b99fa8b`
(Roadmap: Ruhe-Marker entfernt, geplante Welle eingetragen) sind
Plan-Doku, tragen keine DoD-Substanz und berühren den Slice-Vertrag
nicht.

---

## 3. ADR-Konformität

- **[`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
  Festlegung 2 (verschärfte Test-Pflicht):** erfüllt — die gRPC-Fläche
  trägt `integration/test_grpc_realserver.py` (2 Tests) gegen eine
  reale, laufende Server-Instanz: der Runner fährt die
  `compose.yaml`-Umgebung hoch (Log: PostgreSQL/NATS/Feed-Container,
  d-migrate-Rollout, Aktivierungs-Vorbedingungen), der Prüfling ist das
  SDK selbst (`pgchangefeed.grpc_client`), der Empfang ist am Wire
  belegt (Tabelle, Operation, Sentinel) und über die `change_id`
  gegen `cdc.changes` gehalten, die Negative-Prüfung öffnet den Stream
  ohne Token und endet bei `UNAUTHENTICATED`. [`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
  §Fitness-Function-Zeile („mindestens einen Test gegen eine reale,
  laufende Server-Instanz") ist damit real getragen.
- **[`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
  Festlegung 4 / [`ADR-0107`](../plan/adr/0107-python-pypi-zweites-sdk-package.md)
  Festlegung 3 (Ort, Import-Grenze):** erfüllt — alles unter
  `sdks/python/`, Import-Messung ohne Treffer (Abschnitt 1, DoD 3);
  kein Eingriff in `internal/**`/`cmd/**`/`gen/**` (nicht im Diff),
  `.a-check.yml` unberührt.
- **Kein `examples/python/`:** `find examples -iname "*python*"` —
  kein Treffer, unverändert zum Parent (Festlegung 2 mandatiert keine
  Beispiel-Fläche).
- **Kein stiller ADR-Eingriff:** keine ADR-Datei im Diff; der
  Umsetzungs-Schnitt (eigenständiges Skript, geteilte Compose-Umgebung,
  kein Wegwerf-Client) liegt innerhalb der vom Slice-Plan §3
  deklarierten Architektur-Entscheidung und der
  [`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)-Delegation
  an den umsetzenden Zug (Folgepflicht 1).
- **[`AGENTS.md`](../../AGENTS.md) §3.13 Träger-Nachzug:** README §Status,
  `__init__.py`-Docstring, `options.py` real nachgezogen (V-3 betrifft
  nur die Deklaration); [`docs/user/benutzerhandbuch.md`](../../docs/user/benutzerhandbuch.md)
  nach F-3 im selben Zug korrigiert (trägt V-2).

---

## 4. SPEC-020-Konformität (im Quelltext geprüft)

| [`SPEC-020`](../../spec/pflichtenheft.md)-Festlegung | Beleg im Quelltext | Status |
|---|---|---|
| RPC `StreamChanges(StreamChangesRequest) returns (stream Change)` | `grpc_client.py` ruft `StreamChanges`; Unit-Test bindet den Pfad `/cdc.stream.v1.ChangeStream/StreamChanges` am Fake | erfüllt |
| Request filterlos | Unit-Test: `isinstance(request, StreamChangesRequest)` + `SerializeToString() == b""`; Proto: `message StreamChangesRequest {}` | erfüllt |
| Zehn Felder in Wire-Reihenfolge | Deskriptor-Assert (Namen + Nummern 1–10) gegen den zur Bauzeit erzeugten Stub; typbewusste Feldschleife am realen Wire | erfüllt |
| `authorization`-Metadata in Bearer-Form | `metadata = (("authorization", "Bearer <token>"),)`; Unit-Test hält die exakte Metadata-Paarung | erfüllt |
| `Unauthenticated` nicht verschluckt | Unit-Test (Iterator hebt den Fehler aus der Enumeration) + Integrationstest (`REJECTED code=Unauthenticated`, Runner-Assertion) | erfüllt |
| Kein Stream-internes Replay ([`LH-FA-SST-008`](../../spec/lastenheft.md) Boundary) | Client-Docstring trägt die Boundary; der Integrationstest verlangt Fire-and-Forget (Runner committet eine begrenzte Folge bis zum Empfang) | erfüllt |

---

## 5. Aufgedeckte Abweichungen

**V-1 — MEDIUM: DoD-Sensor-Beleg-Zahl driftet gegen die Messung
(29/6 statt 31/8).**
Fundstelle: `Slice-Plan` `slice-sdk-python-grpc-client-flaeche`
§2, Zeilen 78–82 („29 Unit-Tests grün, darunter 6 neue gRPC-Tests
(`grep -c "^def test_"` = 6; 20 + 3 + 6 = 29)"). Messung: die
Testdatei trägt **8** `def test_`-Funktionen, der Pack-Lauf der
Fixrunde zählt **31 passed** — 20 + 3 + 8. Die Fixrunde `b7a993fe`
setzte den korrigierten Beleg-Text (6/29) und die beiden neuen
F-4-Tests **in denselben Commit**, ohne die Zahl erneut zu ziehen; die
Plan-Nachzug-Zeile im selben Dokument („`def test_`-Zahl jetzt 8")
widerspricht der DoD-Zeile im selben Zug. Klasse: `AGENTS.md` §3.12
Instanz A — dieselbe Klasse wie Review F-1, jetzt beim zweiten
Nachziehen wieder aufgetreten (5 → 6 → 8, der Träger hinkt der Messung
jedes Mal hinterher). Die Substanz der DoD-Zeile trägt (mehr Tests
grün als behauptet, Wheel/SDist korrekt); die Zahl braucht eine
Korrektur vor Closure.

**V-2 — LOW: Änderungshistorie des Benutzerhandbuchs trägt die
Versionsnummer 1.39 doppelt.**
Fundstelle: [`docs/user/benutzerhandbuch.md`](../../docs/user/benutzerhandbuch.md)
Historie: neue Zeile „1.39 | 2026-09-23" (Python-gRPC, diese Fixrunde)
steht vor der bestehenden Zeile „1.39 | 2026-09-22" (Kotlin-SSE/NATS,
aus `c8c9e3ae`); der Kopf wurde 1.38 → 1.39 gehoben. Die neue Zeile
hätte 1.40 tragen müssen — `1.39` war durch die Kotlin-Zeile bereits
vergeben, und der Historie-Schwanz ist jetzt nicht mehr im
Versions-/Datums-Abwärtsgang (1.39/2026-09-23 steht vor
1.39/2026-09-22). Kein d-check-Befund (kein gebrochener Verweis), aber
eine öffentliche Versions-Aussage, die zwei verschiedene
Änderungseinträge unter derselben Nummer führt. Randnotiz für den
Planner: der Parent-Stand trug bereits die Unstimmigkeit „Kopfzeile
1.38, neueste Historie-Zeile 1.39" (aus dem Kotlin-NATS-Zug ohne
Kopfzeilen-Nachzug) — die Fixrunde hat den Kopf auf 1.39 gehoben und
damit die Kollision sichtbar gemacht statt aufgelöst.

**V-3 — LOW: `options.py`-Änderung ohne eigene Plan-Nachzug-Zeile.**
Fundstelle: `Slice-Plan` `slice-sdk-python-grpc-client-flaeche`
§3 Plan-Nachzug Zeile 7 (nennt nur `README.md` §Status und
`__init__.py`-Docstring; `options.py` steht dort nur als
Klassen-Referenz in der Begründung) gegen den realen Diff
(`options.py`, 21 Zeilen, Implementations-Commit). Die Änderung ist
sachlich korrekt und durch die Welle-Plan-Pflicht gedeckt — sie fehlt
nur als deklarierte Zeile in der Nachzug-Tabelle, die das
Übergabe-Artefakt der nächsten Rollen ist.

---

## 6. Verdikt

**DoD erfüllt mit Auflagen.**

Alle vier `[x]`-Checkboxen tragen real: Die gRPC-Fläche ist
[`SPEC-020`](../../spec/pflichtenheft.md)-konform gebunden (Pfad,
filterloser Request, zehn Felder, Bearer-Metadata,
Unauthenticated-Boundary — Unit und Realserver), das
Realserver-Integrationstest-Werkzeug existiert und lief real grün
(change_id=804-1, SQL-gegengeprüft), die Import-Grenze hält, und
`make gates` ist über den exakten Ist-Baum mit frischem Stempel grün
(alle sechs Gates einzeln im Log belegt). Die offenen `[ ]`-Zeilen
sind korrekt von den erledigten getrennt.

Auflagen vor Closure (alle drei reine Text-Korrekturen, kein Code):

1. **V-1:** DoD-Sensor-Beleg auf 31 Unit-Tests / 8 gRPC-Tests
   korrigieren (Zahl + Summenformel), damit die DoD-Zeile und die
   Plan-Nachzug-Zeile im selben Dokument nicht widersprechend bleiben.
2. **V-2:** die neue Historie-Zeile im Benutzerhandbuch auf 1.40
   setzen (Kopfzeile entsprechend), damit 1.39 eindeutig bleibt.
3. **V-3:** `options.py` als eigene Plan-Nachzug-Zeile deklarieren
   (Änderungs-Art update, Träger-Nachzug-Klasse).

Offen regulär: Review-Nachzug nach dem Re-Review der Fixrunde,
Closure-Notiz mit Steering-Loop-Eintrag, Beobachtungs-Register-Antwort
(keine Beobachtung angefallen), Risiko-Ausgänge bei Closure, Drei
Paarungen bei der Welle-Closure.

**Grenze dieses Laufs:** der reale `make test-sdk-python-integration`-
Lauf wurde nicht wiederholt (Auftragsgrenze) — sein Beleg ist die
Log-Datei der Fixrunde plus die Skript-Assertion-Kette, beide hier
ausgewertet; die drei Log-Dateien liegen im transienten Host-Temp und
sind durch die Zitate in diesem Bericht konserviert. Mutationsproben
(„Bewusstes Brechen", Modul 11) wurden für die Unit-Bindungen als
statische Bindungsprüfung geführt (je Assertion ist die rot machende
Client-Änderung benennbar: RPC-Pfad, Metadata-Paarung, filterloser
Request, Timeout-Weiterleitung, Unauthenticated-Weiterreichung); der
Integrationstest trägt seine Bindung über die Runner-Assertion-Kette
(`set -euo pipefail`, fünf explizite Fehl-Ausgänge), nicht über einen
erneuten Ausfall-Lauf.