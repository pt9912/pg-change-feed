# Verifikationsbericht: slice-sdk-python-sse-client-flaeche — 2026-09-23

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?" gegen den
DoD-Vertrag (`Slice-Plan` `slice-sdk-python-sse-client-flaeche`
§2), den Plan-vs-Code-Diff (§3 + Plan-Nachzug + beide
§3.13-Suchlauf-Felder) und die ADR-Konformität
([`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
§Entscheidung Festlegung 2, §Konsequenzen Folgepflicht 1).
**Nicht** gegen den Diff als solchen (Reviewer-Aufgabe, abgeschlossen mit
[Haupt-Review](review-slice-sdk-python-sse-client-flaeche.md) und
[Fixrunden-Report](review-slice-sdk-python-sse-client-flaeche-fixrunde.md))
und nicht gegen realen Bedarf (Validator, hier nicht ausgelöst).

**Gegenstand:** `git diff b43251ee..HEAD` — fünf Commits auf `main`:
`6bbe99d9` (Implementations-Commit), `3c941b0b` (Fixrunde 1, F-1…F-5),
`228dfc9e` (Link-Nachzug im Haupt-Review-Report, gemessen: nur
Link-Formen), `beeddc3c` (Fixrunde 2, FR-1), `b8345825` (Linkziele im
Fixrunden-Report). Die DoD-Substanz liegt in `6bbe99d9` + `3c941b0b` +
`beeddc3c`.

**Methode — Belege, nicht Behauptungen:** Jede DoD-Checkbox-Zeile wurde
gegen eine eigene Messung gehalten, nicht gegen den Berichtstext:
Testzahlen per `grep -c "^def test_"` über alle vier Testdateien; die
Unit-Suite **selbst gefahren** netzlos im lokalen
Integration-Image (`docker run --network none` — 41 passed; SSE-Teilmenge
10 passed); **beide Mutations-Nachweise** (Modul 11 „Bewusstes Brechen")
selbst gefahren — Bearer-Header-Form mutiert und Statuscode-Prüfung
umgangen, je über Einzeldatei-Mount in das Image (Arbeitsbaum blieb
sauber, nichts committet); `make gates` **selbst gefahren** mit
separat gesichertem `make`-eigenem Exit-Code (`AGENTS.md` §3.9);
`make test-sdk-python-integration` **selbst gefahren** (Exit-Code
separat gesichert); Import-Zeilen über `src`/`tests`/`integration`
vollständig aufgelistet; beide §3.13-Suchlauf-Stände per `git show`
an Parent und HEAD gemessen; Wheel **und** SDist in
`sdks/python/dist/` auf die SSE-Quelle geprüft. Kein Stempel-Glaube:
der Nachweis-Stempel wurde vor dem eigenen Gates-Lauf gegen den
Baum-Hash gehalten (byte-gleich) und nach dem eigenen Lauf erneut
gelesen (unverändert).

---

## 1. DoD je Zeile

| # | DoD-Zeile | Behauptet | Eigene Messung | Verdikt |
|---|---|---|---|---|
| 1 | [`LH-FA-SST-009`](../../spec/lastenheft.md) SSE-Fläche (`[x]`) | 41 Unit-Tests grün, 10 neue SSE-Tests, `make sdk-pack-python` EXIT=0 | erfüllt (Zahlen nachgemessen, Suite selbst gefahren, Mutationen rot aus dem richtigen Grund) | erfüllt |
| 2 | Realserver-Integrationstest erweitert (`[x]`) | `make test-sdk-python-integration` EXIT=0, change_id=804-1/807-1, Unauthenticated/401 | erfüllt (Lauf selbst gefahren, EXIT=0; Werte laufgebunden — Notiz N-1) | erfüllt |
| 3 | Kein Import aus `internal/**`/`cmd/**`/`gen/**` (`[x]`) | strenge Import-Prüfung ohne Treffer | erfüllt (0 Content-Treffer) | erfüllt |
| 4 | `make gates` grün (`[x]`) | Stempel frisch | erfüllt (eigener Lauf, Exit 0, alle sechs Gates) | erfüllt |
| 5 | Review-Report (`[x]`) | F-1…F-5 + FR-1 gelöst, Endstand kein offenes HIGH/MEDIUM | Substanz erfüllt — **mit Abweichung V-2** (Re-Review-Beleg der Fixrunde 2 nicht committet) | erfüllt mit Abweichung |
| 6 | Doku-Update (`[ ]`) | — | korrekt offen (NATS-Handbuch-Teil + `spec/pflichtenheft.md` gebündelt); der SSE-Anteil ist im Zug real erledigt (Absatz, Version 1.42, Historie gemessen) | Trennung korrekt |
| 7 | Closure-Notiz (`[ ]`) | — | korrekt offen (§7 trägt Platzhalter) | Trennung korrekt |
| 8 | Reconciliation (`[ ]`, „entfällt") | — | `docs/plan/planning/reconciliation.md` existiert nicht | korrekt |
| 9 | Beobachtungs-Register (`[ ]`) | — | korrekt offen (Antwort in §7 bei Closure) | Trennung korrekt |
| 10 | Risiko-Ausgänge (`[ ]`) | — | Ausgänge in §6 geschrieben (1× weiter offen bis zum ersten realen Lauf — durch den eigenen Lauf dieses Berichts jetzt belegt; 1× entfallen mit Begründung) | Trennung korrekt |
| 11 | Drei Paarungen (`[ ]`) | — | [Welle](../plan/planning/welle-sdk-python-vollabdeckung.md) offen | Trennung korrekt |

### DoD 1 — SSE-Client-Fläche (erfüllt)

- **Öffentliche Klasse:** `sdks/python/pgchangefeed/src/pgchangefeed/sse_client.py`
  trägt `PgChangeFeedSseClient` (`httpx.Client` injiziert, nicht
  besessen, [`ClientOptions`](../../spec/pflichtenheft.md)-Konfiguration)
  und wird über `__init__.py` (Import + Eintrag in `__all__`,
  Zeile 35) Package-sichtbar exportiert.
- **Endpunkt- und Header-Form:** `GET /changes/stream`
  (`_STREAM_PATH`, Zeile 55; URL-Bildung Zeile 75), Bearer-Token im
  `Authorization`-Header (Zeile 76) — [`SPEC-021`](../../spec/pflichtenheft.md)-konform.
- **Frame-Parser:** zeilenweise über `iter_lines()` — `event: `/`data: `
  sammeln die Felder, die Leerzeile schließt das Frame ab, ein Frame
  vor seiner Leerzeile wird am Stream-Ende verworfen, Nicht-`change`-Eventnamen
  werden übersprungen (Zeilen 84–106). Chunk-Grenzen sind explizit
  ungebunden (Test `test_chunk_boundaries_need_not_coincide_with_lines`
  mit eigenem `_ChunkedStream`).
- **Zehn Felder:** `StreamChange` in `models.py` (10 Felder, eigene
  getypte Klasse, Plan-Nachzug-Zeile) mit `from_json`; fehlendes Bild
  `None`; ein fehlendes Feld endet typisiert
  (`PgChangeFeedMalformedResponseError`, drei Befund-Formen:
  kein JSON / kein Objekt / fehlendes Feld).
- **Testzahlen — gemessen, nicht behauptet:** `grep -c "^def test_"`
  liefert `test_http_client.py` = **20**, `test_options.py` = **3**,
  `test_grpc_client.py` = **8**, `test_sse_client.py` = **10** →
  **41**; Summe stimmt. Die Suite lief in meinem eigenen netzlosen
  Image-Lauf: `41 passed in 0.10s`, die SSE-Teilmenge separat
  `10 passed in 0.08s`. Der Pack-Lauf-Log der Fixrunde (transient,
  Host-Temp-Verzeichnis, `sdk-pack-s2-fix.log`) trägt `41 passed in 0.13s` — die
  `sdks/python/**`-Dateien sind seit diesem Bau unverändert
  (`git diff 3c941b0b..HEAD` berührt sie nicht).
- **Artefakte tragen die Quelle, keine Stubs:** Wheel und SDist in
  `sdks/python/dist/` enthalten `pgchangefeed/sse_client.py` samt der
  `StreamChange`-Quelle (entpackt und geprüft) — die SSE-Fläche ist
  hand-gemappter Quellcode; ein `grpc_gen`-Nachzug ist gegenstandlos,
  die Stub-Fläche bleibt die gRPC-Fläche.
- **Mutations-Nachweise (Modul 11), beide rot aus dem richtigen Grund:**
  1. **Bearer-Header-Form** (Zeile 76) mutiert auf `Token <token>` →
     genau ein Test rot:
     `test_stream_changes_yields_typed_events_with_the_ten_fields`,
     `AssertionError` am Handler-Assert des Fakes
     (`tests/test_sse_client.py:63`, „Token e2e-reader-token" ≠
     „Bearer …") — 9 passed. Die Header-Zusage ist an der Eingabeseite
     gebunden, nicht an einer späteren Stelle.
  2. **Statuscode-Prüfung** (Zeile 78) umgangen (`if False and …`) →
     beide Status-Tests rot: `test_unknown_token_surfaces_as_typed_unauthorized_error`
     und `test_503_without_broadcaster_surfaces_as_typed_error`, je
     `DID NOT RAISE` — 8 passed. Die 401-/503-Bindung trägt die
     [`SPEC-021`](../../spec/pflichtenheft.md)-Negativen am Unit-Rand;
     `401` landet über `_STATUS_TO_ERROR` (`http_client.py:70`) auf
     `PgChangeFeedUnauthorizedError`.
  Beide Mutationen waren flüchtig (`/tmp`-Mount), der Arbeitsbaum blieb
  sauber (`git status` leer, gemessen).
- **Kein stiller Ausschluss (F-3-Negativprobe, real gemessen):** der
  blanke Image-Aufruf ohne `PGCHANGEFEED_TEST_FILE` endet laut —
  `PGCHANGEFEED_TEST_FILE: Testdatei je Flaeche erforderlich`, Ausgang 2
  (eigener Lauf) — nicht still mit „41 passed". Der Guard feuert bei
  unset **und** leer (`:?`); der Runner setzt je Phase explizit.

### DoD 2 — Realserver-Integrationstest erweitert (erfüllt, eigener Lauf)

- **Lauf-Beleg aus eigener Ausführung:** `make test-sdk-python-integration`
  **selbst gefahren**, `make`-eigener Exit-Code separat gesichert:
  **0**. Der Runnerschlusssatz (nur nach allen Bestehungen beider
  Phasen erreichbar — fünf Fehl-Ausgänge je Phase vor ihm, `set -euo
  pipefail`) trägt: gRPC-Fläche (`pgchangefeed.grpc_client`,
  change_id=**804-1**) und SSE-Fläche (`pgchangefeed.sse_client`,
  change_id=**808-1**) empfingen je eine danach committete Änderung
  (Tabelle, Operation, Sentinel am Wire), ein Stream-Öffnungsversuch
  ohne Token endete je mit gRPC-Status `Unauthenticated` bzw. HTTP-Status `401`.
- **Verdrahtung der SSE-Phase:** `run_surface_phase` mit
  `PGCHANGEFEED_HTTP_ADDR=http://pg-change-feed:8090`, Sentinel
  `PythonSseSdkE2ESentinel`, ID-Basis 310, Testdatei
  `integration/test_sse_realserver.py` als Umgebungsvariable
  `PGCHANGEFEED_TEST_FILE`, Reject-Marker `REJECTED status=401` —
  der Marker stammt ausschließlich aus der SSE-Testdatei (Zeile 103),
  der Phasen-Erfolg ist daher an diese Datei gebunden.
- **SQL-Gegenprüfung:** die `RECEIVED`-`change_id` wird je Phase gegen
  `cdc.changes` gehalten (`SELECT count(*) … change_id = … AND
  table_name = … AND new_data->>'name' = '<Sentinel>'`, Runner-Zeilen
  269–274) — dieselbe Disziplin wie der Server-E2E-Runner.
- **Nachrichtenschema am realen Wire:** der SSE-Integrationstest hält
  alle zehn Felder typbewusst (`isinstance` je Feld,
  Nicht-`bool`-Guard), Identitäten non-empty, `INSERT` ohne Alt-Bild,
  Sentinel im Row Image, `schema = public` (Zeilen 63–86).
- **Grenze, sauber benannt:** der Testcontainer läuft detached, seine
  Marker liest der Runner über `docker logs`; der Beleg ist die
  Skriptstruktur (Fehl-Ausgänge + `set -e`) plus der interpolierte
  `change_id`-Wert, kein wörtliches `EXIT=0` im Log — dieselbe Grenze,
  die der Vorgänger-Verifikationsbericht bereits benennt.
- **Notiz N-1 (keine Verletzung):** der DoD-Beleg nennt für den
  Implementer-Lauf change_id=804-1 (gRPC) und 807-1 (SSE); mein Lauf
  maß 804-1 (gRPC, exakt deckend) und **808-1** (SSE). `change_id` ist
  `<tx.ID>-<sequence>` (`internal/adapters/driving/replication/mapper/mapper.go:243`)
  und damit eine **laufgebunde** Messgröße; der DoD bindet die Werte an
  ihren Lauf (§3.12 Instanz A erfüllt). Die Differenz zwischen den
  Läufen (zwei Commit-Vorgänge weiter) ist per-Run-Zustand, kein
  Widerspruch — der Nachweis trägt über die Struktur (Phasen, Marker,
  SQL-Gegenprüfung, Reject), nicht über die Zahl.

### DoD 3 — Import-Grenze (erfüllt)

Alle `from`/`import`-Zeilen unter `src`, `tests` und `integration`
aufgelistet: ausschließlich Standardbibliothek, `httpx`, `pytest`,
`grpc` und `pgchangefeed.*`. Der Prüf-Grep liefert einen einzigen,
pfadgetriebenen Treffer — der Docstring-Text ``gen/**`` im
Package-Marker `grpc_gen/__init__.py` matcht das Muster `gen/`; kein
Import. Kein Treffer gegen `internal/**` oder `cmd/**`.

### DoD 4 — `make gates` grün (erfüllt, eigener Lauf)

`make gates` **selbst gefahren** (Exit-Code ungepiped separat
gesichert: **0**). Vor dem Lauf war der Stempel
`.harness/state/gates-passed.diffsha` gegen den frischen Baum-Hash
(`tools/harness/working-tree-hash.sh`) **byte-gleich**; nach dem Lauf
unverändert. Die sechs Einzelbelege aus dem eigenen Lauf:

| Gate | Log-Zeile |
|---|---|
| `baseline-verify` | `v6.9.0 OK — 54 Dateien (Integritaet + Vollstaendigkeit, netzlos)` |
| `docs-check` | `d-check: 913 Datei(en) geprüft, 0 Befund(e)` |
| `commit-traceability` | `OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` |
| `coverage-gate` | `coverage-gate: OK — Coverage 82.80% erfüllt Schwelle 80%` |
| `generated-sync` | `OK — das committete Erzeugnis ist byte-gleich der Ausgabe des gepinnten Generators (Stufe proto-export)` |
| `a-check` | `gesamt: 0 Befund(e)` |

### DoD 5 — Review-Report (erfüllt, mit Abweichung V-2)

Beide Reports liegen committe vor (Haupt-Review F-1…F-9, Fixrunden-Report
F-1…F-5/FR-1). Die Fixrunden-Substanz ist **real und vom Verifier
gemessen**: der Handbuch-`**SDK:**`-Absatz steht (Zeile 874, `pip
install pgchangefeed`, Endpunkt, Iterator über getypte `StreamChange`
mit zehn Feldern, Bearer-Header, `PgChangeFeedUnauthorizedError` — je
Assertion gegen den Code gehalten), Version-Kopf `1.42` und Historie-Zeile
`1.42` im selben Träger; die `harness/README.md`-Zeile trägt die
Zwei-Flächen-ENV-Form; der F-3-Guard scheitert laut (Ausgang 2, gemessen);
der F-4-Docstring nennt die [`SPEC-021`](../../spec/pflichtenheft.md)-Form
(`exceptions.py:53`); `sse_client.py` endet mit Zeilenumbuch (gemessen).
**Abweichung V-2** (Abschnitt 5): die behauptete Endstand-Form
„Re-Review-Endstand: kein offenes HIGH/MEDIUM" trägt keinen committeten
Re-Review-Beleg der Fixrunde 2 — der Fixrunden-Report endet mit seinem
Verdikt („FR-1 braucht eine kurze zweite Fixrunde"), ein Addendum zur
gelösten FR-1 fehlt. Die Substanz selbst ist hier nachgemessen (alle drei
Träger tragen die ENV-Form; kein Treffer der alten Phrasen an HEAD).

### DoD 6–11 — korrekt getrennt (offen, wo die Rollen-Sequenz noch nicht lief)

- **Doku-Update (`[ ]`):** der im Zug erledigte SSE-Handbuch-Anteil ist
  real gemessen (Absatz + `1.42` + Historie); gebündelt bleiben
  NATS-Handbuch-Teil und `spec/pflichtenheft.md` — dieselbe
  Bündelungs-Praxis wie beim Vorgänger-Slice. Die Box bleibt korrekt
  offen, weil der gebündelte Teil fehlt.
- **§6-Risiko 1 (Streaming-Pufferung):** sein Ausgang „weiter offen bis
  zum ersten realen Integrationstest-Lauf" ist durch den eigenen
  Integrationstest-Lauf dieses Berichts nun belegt (der zeilenweise
  Parser lief real gegen den laufenden Feed-Container) — die Schließung
  des Ausgangs gehört in die Closure-Notiz, nicht in die Checkbox.
- **Closure-Notiz, Beobachtungs-Register, Risiko-Ausgänge, Drei
  Paarungen:** korrekt als `[ ]` — §7 trägt Platzhalter; die Prüfung
  läuft bei der Wellen-Closure. Reconciliation „entfällt" ist korrekt.

---

## 2. Plan-vs-Code

**§3-Tabelle (4 Zeilen) — vollständig im Diff vertreten:**
`sse_client.py` (neu, 134 Zeilen), `tests/test_sse_client.py` (neu,
236 Zeilen, 10 Tests), `sdks/python/Dockerfile` (update —
`integration`-Stufe, CMD-Guard), `integration/test_sse_realserver.py`
(neu, am Nachzug-Ort). Keine still gestrichene Planzeile.

**Plan-Nachzug (7 Zeilen) — vertreten:**
`integration/`-Ordner statt `tests/integration/` (real am Nachzug-Ort),
Dockerfile-`CMD` mit `:?`-Guard je Phase (real, Negativprobe Ausgang 2),
`models.py`-`StreamChange` (10 Felder, real), Runner-`run_surface_phase`
(real, gRPC 300ff./SSE 310ff., eigene Reject-Marker-Formen, SQL-Gegenprüfung
je Phase), `__init__.py`-/`options.py`-Docstrings + `sdks/python/README.md`
§Status (real nachgezogen, beide Stände gemessen — Parent
`b43251ee` trug die „SSE bleibt außerhalb"-Form an allen drei, HEAD in
keinem), `exceptions.py`-Docstring (real), Träger-Nachzug nach
[`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
Folgepflicht 3 (Handbuch-Absatz, real). Die zweite Fixrunde (`beeddc3c`)
trug die drei Phrase-Stellen auf die ENV-Form nach — beide Stände
gemessen: an `6bbe99d9` trug der Runner-Kopf „Stufe-ENTRYPOINT"/„docker
run-Argument", an HEAD in keinem der Träger (Grep ohne Treffer).
Abweichung V-1 betrifft die Stand-Deklaration dieser Tabelle (Abschnitt 5).

---

## 3. [`SPEC-021`](../../spec/pflichtenheft.md)-Konformität

| Merkmal | Festlegung | Im Code | Verdikt |
|---|---|---|---|
| Endpunkt / Methode | `GET /changes/stream` | `_STREAM_PATH` + `self._client.stream("GET", …)`; Unit-Test hält Methode und Pfad am Handler-Assert | erfüllt |
| Rechtsklasse | `reader` oder `admin` | Token aus `ClientOptions` (ein Bearer-Token); Docstring nennt die beiden Klassen; der Integrationstest trägt die Reader-Klasse (`compose.yaml`-`CDC_API_TOKEN_READER`) | erfüllt |
| Event-Typ | `event: change` | Parser sammelt `event: `-Zeile, yieldt nur `event_name == "change"`; Nicht-`change`-Namen übersprungen (Unit-Test `keepalive`) | erfüllt |
| Event-Daten | `data:` mit JSON-Objekt, denselben zehn Feldern, eingebettete Row Images, fehlendes Bild `null` | `_parse_stream_change` → `StreamChange.from_json`; Unit-Tests halten die zehn Felder und `null`-Bild; Integrationstest hält sie typbewusst am realen Wire | erfüllt |
| Frame-Form | Event-Zeile, `data:`-Zeile, Leerzeile schließt | zeilenweiser Parser, Leerzeile schließt; unvollständiges Frame verworfen (Unit-Test) | erfüllt |
| Zustellgarantie / Replay | Fire-and-Forget, **kein** Stream-internes Replay, `Last-Event-ID` weder gesendet noch ausgewertet | kein `Last-Event-ID`-Code im Package (grep über `src`/`tests`/`integration`); Docstring trägt die Boundary; serverseitige Gegenprobe `internal/adapters/driving/http/sse.go` (nur gelesen) trägt dieselbe Grenze | erfüllt |
| 401 | fehlender/unbekannter Token | Unit-Test (Mock, `status_code == 401`), Integrationstest (`pytest.raises(PgChangeFeedUnauthorizedError)`, Marker `REJECTED status=401`), Mapping über `_STATUS_TO_ERROR` | erfüllt |
| `503` ohne Broadcaster | — | Unit-Test hält den Status, Integration nicht berührt (kein DoD-Anspruch) | erfüllt |

**Abgrenzung zu [`SPEC-022`](../../spec/pflichtenheft.md):** `StreamChange`
trägt zehn Felder, nicht die zwölf des HTTP-Lesezugriffs (`committed_at`
fehlt korrekt) — der Kommentar in `models.py` trägt dieselbe Abgrenzung.

---

## 4. [`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md) Festlegung 2 / Folgepflicht 1

- **Realer Server-Rundlauf für die SSE-Fläche:** erfüllt — der Runner
  verdrahtet die SSE-Phase (Adresse, Sentinel, ID-Bereich, Reject-Marker,
  SQL-Gegenprüfung), und mein eigener Lauf belegt sie real (Abschnitt 1,
  DoD 2). Nicht ausschließlich Unit-Tests gegen einen Mock — die
  Fitness-Function-Zeile der ADR (Review-Aufgabe) trägt.
- **Import-Grenze (Festlegung 3 aus [`ADR-0107`](../plan/adr/0107-python-pypi-zweites-sdk-package.md),
  unverändert durch [`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)):** erfüllt (Abschnitt 1, DoD 3).
- **Kein `examples/python/`:** erfüllt — `ls examples/` trägt
  `csharp`/`kotlin` und die flache Go-Wurzel, kein `python`-Verzeichnis;
  die Referenzmaterial-Basis des Plans (C#/Kotlin-SSE-Clients, nur
  gelesen) hält.

---

## 5. Abweichungen (aufgedeckt, beide LOW)

- **V-1 — §3.13-Suchlauf, zweite Eigenschaft: Stand-Deklaration für den
  `harness/README.md`-Träger unpräzise.** Die Tabelle deklariert
  „beide Stände gemessen: Parent `6bbe99d9` und HEAD" und führt den
  `harness/README.md`-Träger als „Phrase gefunden". Gemessen: an
  `6bbe99d9` trug die Werkzeug-Zeile die Ein-Flächen-Form **ohne** die
  Phrase „docker run-Argument" — die Phrase kam erst mit der Fixrunde 1
  (`3c941b0b`) in den Träger und wurde in Fixrunde 2 (`beeddc3c`)
  gezogen. Der wahre Vorher-Stand dieses Trägers ist `3c941b0b`, nicht
  `6bbe99d9`. Die Substanz ist unberührt (an HEAD trägt keiner der drei
  Träger die alte Phrase, gemessen; ENV-Form an allen drei). Korrektur
  gehört in die Closure-Notiz-Vorbereitung, kein Code-Pfeil.
- **V-2 — DoD-Zeile 5: der „Re-Review-Endstand" trägt keinen committeten
  Beleg der Fixrunde 2.** Der Fixrunden-Report endet mit dem Verdikt,
  FR-1 brauche eine zweite Fixrunde; ein Addendum, das die gelöste FR-1
  bestätigt, fehlt. Der Implementer-Werkflow (Schritt 21) zieht die
  Checkbox regulär nach der Fixrunde — das trägt die Checkbox selbst,
  aber die Endstand-Behauptung im Plan-Text ist bis zu dieser
  Verifikation eine Behauptung ohne committeten Beleg (§3.12 Instanz B).
  Die Substanz ist hier nachgemessen und bestätigt (alle drei Träger
  tragen die ENV-Form, die alte Phrase ist an HEAD nullfachig) — die
  Lücke ist die Form, nicht der Inhalt. Korrektur: einzeiliges
  Re-Review-Addendum im Fixrunden-Report vor der Closure-Notiz.
- **N-1 — Notiz, keine Verletzung:** die im DoD-Beleg genannten
  `change_id`-Werte sind laufgebunden (`<tx.ID>-<sequence>`,
  `internal/adapters/driving/replication/mapper/mapper.go:243`); mein
  eigener Lauf maß 808-1 statt der behaupteten 807-1 für die SSE-Phase
  (gRPC deckt sich exakt mit 804-1). Der DoD bindet die Werte an seinen
  Lauf; der Nachweis trägt über die Struktur, nicht über die Zahl.

## Negativbefunde (geprüft, ohne Befund)

- Wheel **und** SDist in `sdks/python/dist/` tragen `sse_client.py` mit
  `StreamChange`-Quelle (entpackt und geprüft) — die SSE-Fläche ist
  Quellcode im Artefakt, keine Stub-Form nötig.
- Kein stiller Ausschluss der Testdatei: blanker Image-Aufruf scheitert
  laut (Ausgang 2, gemessen); je Phase nennt der Runner seine Testdatei
  als Umgebungsvariable; der `pytest`-Lauf der `build`-Stufe sammelt
  nur `testpaths = tests` (41 Unit-Tests — die beiden
  Realserver-Dateien fehlen dort korrekt).
- Commit-Message-Traceability aller fünf Commits des Range: je
  [`LH-FA-SST-009`](../../spec/lastenheft.md) +
  [`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md),
  kein `SPEC-*`/`ARC-*` im Betreff; das Standing-Gate bestätigte
  `HEAD~5..HEAD` im eigenen Gates-Lauf.
- `AGENTS.md` §3.11: der Diff nennt keinen host-lokalen Absolutpfad;
  §3.5: [`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
  (`Accepted`) unberührt; §3.2: kein `//nolint`-Form im neuen Python-Code.
- Fake-Muster-Parität: `httpx.MockTransport` in `test_sse_client.py`
  wie in der bestehenden `test_http_client.py` — dasselbe Seam.

---

## Summary

| Kategorie | Anzahl |
|---|---|
| DoD-Zeilen erfüllt (inkl. Substanz) | 11 von 11; 5 offen-korrekt getrennt |
| Abweichungen (LOW) | 2 (V-1, V-2) |
| Notizen ohne Verletzung | 1 (N-1) |

**Eigene Sensor-Läufe dieses Berichts:** `make gates` (Exit 0, alle sechs
Gates), `make test-sdk-python-integration` (Exit 0, beide Phasen),
Unit-Suite im Integration-Image netzlos (41 passed; SSE-Teilmenge
10 passed), zwei Mutations-Läufe (je rot aus dem richtigen Grund),
Stempel-Abgleich vor und nach dem Gates-Lauf, Wheel-/SDist-Inspektion.

## Verdikt

**Gesamt-Urteil: erfüllt mit Auflagen.** Der DoD-Vertrag (§2) ist in
seiner Substanz vollständig und durch eigene Messung bestätigt — die
beiden LOW-Abweichungen sind Form- und Buchhaltungs-Auflagen, keine
Code-Pfeile: (1) die Stand-Deklaration des zweiten §3.13-Suchlaufs für
den `harness/README.md`-Träger korrigieren (wahrer Vorher-Stand
`3c941b0b`), (2) den Re-Review-Endstand der Fixrunde 2 durch ein
kurzes Addendum im Fixrunden-Report committet belegen — beides vor der
Closure-Notiz. Keine DoD-Verletzung, kein offener HIGH/MEDIUM-Gegenstand
aus dieser Verifikation.