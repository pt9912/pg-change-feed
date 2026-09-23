# Review-Report: slice-sdk-python-nats-stream-client-flaeche — 2026-09-23

**Review-Art:** Code — geprüft gegen Plan
(`docs/plan/planning/in-progress/slice-sdk-python-nats-stream-client-flaeche.md`,
§2 DoD, §3 Plan + Plan-Nachzug, §6 Risiken), `ADR-0110` (Accepted) Festlegung
1/2/3, `ADR-0100` (NATS-Vollinhalts-Vertrag) und `AGENTS.md` Hard Rules
(Modul 10 §Drei Review-Arten). Kein DoD-Abgleich als solcher — das ist
Verifier-Aufgabe (Modul 11); wo eine DoD-Zeile eine **Zahl** behauptet, ist
das Nachmessen dieser Zahl trotzdem Reviewer-Aufgabe (`AGENTS.md` §3.12,
Reviewer-Skill §HIGH „Zahl im Träger … driftend").

**Gegenstand:** Commit `0f8cc4f2` („feat(sdk-python): NATS-Vollinhalts-
Client-Flaeche + Version-Hebung 0.2.0 + Traeger-Nachzug"), Diff-Range
`fa71cc96..HEAD`, Slice `slice-sdk-python-nats-stream-client-flaeche`,
**letzter** Flächen-Slice der Welle `welle-sdk-python-vollabdeckung`
(11 Dateien, 535 Insertions / 40 Deletions — bündelt die NATS-Fläche mit
Version-Hebung `0.1.0` → `0.2.0` und dem Träger-Nachzug für alle drei neuen
Flächen, gRPC + SSE + NATS).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09"
(seither um mehrere HIGH-Klassen ergänzt — u. a. „Zahl im Träger …
driftend", „Beleg trägt seinen Satz nicht" mit der eingefassten
Nachbar-Form „Zitat nennt die falsche Stelle", „Zusage ohne Bindung an ihre
Eingabeseite").
**Modell:** glm-5.3-flash (Claude-Agent-SDK-Subagent) · **Datum:** 2026-09-23.

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-sdk-python-nats-stream-client-flaeche.md`
  (§1 Ziel/Abgrenzung, §2 DoD, §3 Plan + Plan-Nachzug, §6 Risiken, §7
  Closure-Notiz)
- `ADR-0110` (Accepted) — Festlegung 1 (Vier-Wege-Umfang), Festlegung 2
  (verschärfte Test-Pflicht, reale Server-Instanz je neuer Fläche, kein
  `examples/python/`), Festlegung 3 (Struktur-Delegation an den umsetzenden
  Zug), §Konsequenzen Folgepflicht 2/3
- `ADR-0100` (NATS-Vollinhalts-Vertrag), `spec/pflichtenheft.md` §2 [`SPEC-024`]
  (Subjekt-Schema, zehn Felder, Verbindungsebene-Auth), `[SPEC-027]`-Zeile §6
- `AGENTS.md` §3.1 (Docker-only), §3.2 (Suppression-Verbot), §3.7
  (Kommentar-Disziplin), §3.12 (Herkunft von Aussagen), §3.13 (Träger-Nachzug)
- `harness/conventions.md` (MR-000/MR-002)
- Server-Gegenprobe (nur gelesen, nicht importiert):
  `internal/adapters/driven/natsstream/publisher.go` (zehn JSON-Keys,
  `subjectPrefix`)
- Vorbilder: `examples/csharp/nats-stream-client/`,
  `tools/harness/natsstreamsub`, die beiden Vorgänger-Slice-Pläne
  (`done/slice-sdk-python-grpc-client-flaeche.md`,
  `done/slice-sdk-python-sse-client-flaeche.md`) und deren Reviews
  (Lernklasse) sowie die beiden Schwester-NATS-Reviews
  (`docs/reviews/review-slice-sdk-csharp-nats-stream-client-flaeche.md` via
  `review-slice-sdk-kotlin-nats-stream-client-flaeche.md` als Kontrastfolie)

**Eigenständig durchgeführte Prüfungen (nicht aus dem Diff/Commit-Text
übernommen):**

- **Träger-Nachzug-Suchlauf (`AGENTS.md` §3.13) über beide Stände gemessen**
  (`fa71cc96` UND `HEAD`, jeweils `git grep`): erst der im Plan §2 deklarierte
  Mustersatz (`deckt HTTP-API` / `pgchangefeed-0.1.0` /
  `no gRPC/SSE/NATS surface exists yet`) über `spec/`, `docs/`,
  `sdks/python/`, dann bewusst breiter über die restlichen
  Formulierungsvarianten der bewegten Eigenschaft (`consume .* HTTP API`,
  `aktuell 0\.1\.0`, `real erzeugt.*pgchangefeed`) über den gesamten Baum
  inklusive `harness/`, Wurzel-READMEs und `.github/`. Ergebnis: **vier
  weitere, real übersehene Fundstellen** (F-1, F-4), alle in beiden Ständen
  identisch — die Arbeit bewegte die Eigenschaft, die Träger blieben stehen.
- **Test-Anzahl real nachgemessen** (`grep -c "^def test_"`
  über `sdks/python/pgchangefeed/tests/*.py`): 20 (http) + 3 (options) +
  8 (grpc) + 10 (sse) + 7 (nats) = **48** — deckungsgleich mit der
  Vorgänger-Formel der SSE-DoD („20 + 3 + 8 + 10 = 41") plus die 7 neuen
  NATS-Tests; keine Drift.
- **[`SPEC-024`](../../spec/pflichtenheft.md)-Draht Feld für Feld gegen den Server gehalten:** die zehn
  JSON-Keys des Client-Parsings (`change_id` … `table`) decken sich exakt
  mit `internal/adapters/driven/natsstream/publisher.go`s
  `json:`-Tags (Zeilen 140–149); Subjekt-Form `cdc.stream.<source_id>.>`
  deckt sich mit [`SPEC-024`]s Tabellenzeile („`cdc.stream.<source_id>.>`
  deckt alle Tabellen einer Quelle") und mit `subjectPrefix`/`subjectFor`
  des Publishers; Authentifizierung Verbindungsebene (`nats.connect(token=…)`,
  kein Per-Message-Header) — entspricht der [`SPEC-024`]-Zeile
  „Authentifizierung".
- **Versions-/Auth-Verkettung des Integrationstest-Werkzeugs real gelesen:**
  `compose.yaml` trägt `CDC_NATS_STREAM_TOKEN: e2e-nats-stream-token`
  (Zeile 114) — identisch mit der Runner-Konstante `NATS_STREAM_TOKEN`
  (Zeile 77); `PGCHANGEFEED_SOURCE_ID=src-e2e` deckt sich mit der
  `cdc.source`-Registrierung des Runners (Zeile 121). Kein Token-Bruch
  zwischen Feed-Container und Test-Container.
- **Smoke-Beleg der Version-Hebung real gemessen:** `sdks/python/dist/`
  trägt `pgchangefeed-0.2.0-py3-none-any.whl` (20481 Bytes) und
  `pgchangefeed-0.2.0.tar.gz` (24639 Bytes), Lauf-Stempel 2026-09-23 07:47 —
  das `0.2.0`-Artefakt-Paar existiert real, die „real geprüft"-Sätze in
  `spec/pflichtenheft.md` §1/§7-Historie sind durch die Messung gedeckt.
  `docs/user/releasing.md` bleibt korrekt (PyPI-Ist-Stand `0.1.0` ist
  historisch richtig — der Tag-Push ist Out-of-Scope dieses Slices).
- **Sync-Wrapper-Ausgang aus §6 am Ort benennt:** Modul-Docstring
  `nats_stream_client.py` („the connect and subscription run in a dedicated
  thread with its own event loop, the generator consumes a queue — one
  iterator form across all four surfaces, the asyncio boundary stays inside
  this module") — der Ausgang aus Plan §6 (Risiko 1) ist am Ort genannt,
  nicht verschwiegen.
- **Unit-Test-Bindungen an ihrer Eingabeseite gelesen:** Subjekt-Test
  (exakte String-Assertion), Schema-Tests (Payload-Mutationen: Nicht-JSON,
  Array, fehlende Felder), Verdrahtungs-Test (`connect_calls ==
  [([NATS_URL], TOKEN)]`), Timeout-Test (Zeitmessung) — alle an der
  Eingabe mutierbar und rot färbend; Ausnahme F-7 (vakuum-Nachklammer) und
  F-6 (Ursachen-Attribution des Realserver-Rejects).
- **`make gates` ungefiltert laufen lassen, Exit-Code direkt (kein Pipe/
  Wrapper dazwischen) geprüft: `0`** (u. a. `a-check: gesamt: 0 Befund(e)`,
  `generated-sync`, `coverage-gate`).
- Backtick-Zeichen je geänderter Markdown-Datei real ausgezählt:
  `docs/user/benutzerhandbuch.md` 2034, `spec/pflichtenheft.md` 1500,
  `sdks/python/README.md` 54, Slice-Plan 274 — alle vier gerade, alle paarig.
- `grep` über alle neuen Python-Dateien nach `internal/`/`cmd/`/`gen/`-Bezug:
  kein Treffer — Importe ausschließlich `nats`, Python-Standardbibliothek,
  `pgchangefeed.*` (Import-Grenze, `ADR-0110` Festlegung 3-Fußnote/
  `ADR-0107` Festlegung 3).

---

## Findings

### F-1 — `spec/pflichtenheft.md` §6: die [`SPEC-027`]-Vertragszeile blieb stehen — Versions-Zelle driftet gegen die eigene Messung, Planpunkt still gestrichen

- `kategorie`: HIGH
- `quelle`: `AGENTS.md` §3.12/§3.13 · `ADR-0110` §Konsequenzen
  Folgepflicht 2
- `pfad`: `spec/pflichtenheft.md:625`
- `befund`: Die `[SPEC-027]`-Zeile der §6-Vertragstabelle trägt weiter
  „PEP 440, `0.x.y` (aktuell `0.1.0`)" und „… deckt bereits dokumentierten
  Drahtvertrag ([`SPEC-018`])" — während derselbe Diff in derselben Datei
  (§1) auf „deckt HTTP-API, gRPC-Stream, SSE und NATS-Vollinhalt" +
  `0.2.0` zieht und `pyproject.toml` real `0.2.0` misst. Beide
  Schwester-Zeilen derselben Tabelle ([`SPEC-026`], [`SPEC-028`]) tragen durch
  die C#-/Kotlin-Vorgänger-Commits „aktuell `0.2.0`" **und** die
  Vier-Verträge-Liste — dieselbe Zeilen-Form, derselbe Nachzug. Der Plan
  §3 versprach die Zeile ausdrücklich („`spec/pflichtenheft.md` §1
  (`LH-FA-SST-009.a`), **§6 (`[SPEC-027]`-Zeile)** | update | Träger-Nachzug");
  der Plan-Nachzug deklariert das Streichen nicht — die einzige
  pflichtenheft-Zeile im Diff neben der §1-Prosa ist die Chronik-Zeile
  (§7-Historie). Damit ist zugleich `ADR-0110` Folgepflicht 2 („
  `spec/pflichtenheft.md` [`SPEC-027`] und `LH-FA-SST-009.a` brauchen einen
  Träger-Nachzug") zur Hälfte unerfüllt. Die INFO-Herabstufung der
  Zahl-drift-Klasse („Träger außerhalb des Diffs") greift nicht — die
  Datei liegt im Diff.
- `verifizierbar`: ja — `sed -n '625p' spec/pflichtenheft.md` gegen
  `grep '^version' sdks/python/pgchangefeed/pyproject.toml` (beide in
  diesem Review-Lauf ausgeführt).
- `klasse`: „Zahl im Träger … gegen die Messung driftend" +
  „Arbeit überholt stehenden Träger" (`BEO-PGC/arbeit-ueberholt-stehenden-traeger`)

### F-2 — `# noqa: BLE001` in der neuen Client-Klasse: Inline-Suppression ohne Linter und ohne ADR

- `kategorie`: HIGH
- `quelle`: `AGENTS.md` §3.2 (Suppression-Verbot) · Reviewer-Skill HIGH
  „Suppression eines Gates (`#noqa`, `//nolint`, …) ohne ADR"
- `pfad`: `sdks/python/pgchangefeed/src/pgchangefeed/nats_stream_client.py:89`
- `befund`: Der einzige `except`-Block des Connect-Pfads trägt
  `# noqa: BLE001` — ein Unterdrückungs-Token für einen Python-Linter
  (`BLE001` ist ein ruff-/flake8-Code), den dieses Repo strukturell nicht
  hat (kein `[tool.ruff]`/`[tool.flake8]` in `pyproject.toml`, kein
  lint-Schritt im `sdks/python/Dockerfile`; repo-weite `noqa`-Suche über
  `*.py`/`*.go`/`*.cs`/`*.kt` liefert genau diesen einen Treffer). Nach
  §3.2 ist das exakt die verbotene Form — tote Dekoration oder ein Vorgriff
  auf ein Profil, das niemand entschieden hat; es gibt auch keinen ADR, der
  die Ausnahme trüge. Der Kommentar **hinter** dem Token („connect-Fehlschlag
  wird sichtbar am Generator weitergereicht") trägt eine legitime
  Zusage-Klasse und bleibt — unterdrückt wird nur das Token selbst.
- `verifizierbar`: ja — `grep -rn "noqa" sdks/` zeigt die Stelle; die
  Linter-Abwesenheit ist am Fehlen jeder Linter-Konfiguration messbar.
- `klasse`: „Suppression ohne Linter/ADR" (§3.2)

### F-3 — `__init__.py`-Docstring zitiert `ADR-0110` Festlegung 3 für „v2 desselben Packages" — die Festlegung entscheidet ausdrücklich nichts

- `kategorie`: HIGH
- `quelle`: `AGENTS.md` §3.12 · Reviewer-Skill HIGH „Beleg trägt seinen
  Satz nicht" (Nachbar-Form „Zitat nennt die falsche Stelle", seit
  welle-d-check unter dieser Klasse gefasst)
- `pfad`: `sdks/python/pgchangefeed/src/pgchangefeed/__init__.py:16`
- `befund`: Die neue Klammer lautet wörtlich „([`SPEC-024`](../../spec/pflichtenheft.md), **[ADR-0110](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
  Festlegung 3**: v2 desselben Packages)". Im Original trägt Festlegung 3
  diese Aussage nicht — sie lautet: die ADR „entscheidet **nicht**, ob das
  Folge-Release eine `v2` desselben `pgchangefeed`-Packages ist oder ein
  eigenständiges Package wird — dieselbe Delegation an „einen Folge-Zug"",
  und `ADR-0110` Festlegung 4 bindet `examples/python/` separat. Die
  tatsächliche v2-Entscheidung trägt der Welle-Plan
  (`welle-sdk-python-vollabdeckung.md` §6: „diese Welle wählt v2 desselben
  Packages"). Wer den genannten Beleg fährt (die Festlegung im Original
  aufschlagen), findet das Gegenteil einer Entscheidung; die entscheidende
  Schicht ist von der Stelle aus unauffindbar.
- `verifizierbar`: ja — `ADR-0110` §Entscheidung Festlegung 3 im Original
  aufschlagen; der Welle-Plan §6 trägt den Satz wörtlich.
- `klasse`: „Beleg trägt seinen Satz nicht" / „Zitat nennt die falsche Stelle"
  (`BEO-PGC/zitat-nennt-die-falsche-stelle`)

### F-4 — §3.13-Träger-Nachzug unvollständig: vier weitere Fundstellen, alle von der Plan-Grep zu schmalen Suchform nicht getroffen

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §3.13 (Träger-Nachzug)
- `pfad`: `harness/README.md:158`, `harness/README.md:160`,
  `README.md:27`, `README.de.md:27`, `sdks/python/README.md:5`
- `befund`: Eigener Suchlauf über beide Stände (identische Treffer — von
  diesem Slice überholt, nicht gezogen): (a) `harness/README.md`
  `make sdk-pack-python`-Zeile nennt weiter „real erzeugt:
  `sdks/python/dist/pgchangefeed-0.1.0-py3-none-any.whl`,
  `…0.1.0.tar.gz`"; (b) `harness/README.md` `make
  test-sdk-python-integration`-Zeile nennt weiter „je Phase eine Fläche:
  `pgchangefeed.grpc_client`, `pgchangefeed.sse_client`" samt nur zwei
  Reject-Formen (gRPC `Unauthenticated`, HTTP `401`) — der dritte
  NATS-Schritt fehlt; (c) Wurzel-`README.md`/`README.de.md` nennen weiter
  „the HTTP API as an official Python client library" bzw. „die HTTP-API
  zusätzlich als offizielle Python-Client-Bibliothek"; (d)
  `sdks/python/README.md` Intro (Datei **im** Diff, nur §Status gezogen)
  trägt weiter „lets a Python application consume PG Change Feed's **HTTP
  API** without implementing the wire protocol itself". Beide
  Schwester-SDK-Slices (C# `d86d1965`, Kotlin `c8c9e3ae`) zogen genau diese
  Zeilen im selben Commit nach — der Kotlin-Slice-Commit benennt die
  Nachzugs-Pflicht für `harness/README.md`/`harness/mk/sdk.mk` sogar
  ausdrücklich in der Commit-Message; der SSE-Slice-Plan hatte den
  Wurzel-README-Nachzug ausdrücklich in diesen letzten Flächen-Slice
  geschoben („bleibt gebündelt"). Der Plan-deklarierte Suchmuster-Raum
  („deckt HTTP-API|pgchangefeed-0.1.0|no gRPC/SSE/NATS surface exists yet"
  über `spec/`, `docs/`, `sdks/python/`) trifft keine der fünf Stellen in
  deren tatsächlicher Schreibform und lässt `harness/` wie die Wurzel-READMEs
  außerhalb des Pfadraums — die §3.13-Form („grep über die Träger nach der
  bewegten Eigenschaft", nicht über die eigenen Muster) verlangt den
  breiteren Lauf.
- `verifizierbar`: ja — die fünf Stellen sind einzeln per `grep`
  aufzulösbar (alle in diesem Review-Lauf gegen beide Stände ausgeführt).
- `klasse`: „Arbeit überholt stehenden Träger"
  (`BEO-PGC/arbeit-ueberholt-stehenden-traeger`; Kalibrierung:
  Kotlin-Review F-2 und SSE-Review F-2, beide MEDIUM)

### F-5 — Suchlauf-Ergebnis fehlt im committeten Plan-Feld; der deklarierte Suchlauf selbst war zu schmal

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §3.13 · Plan §2 DoD (Suchlauf-Pflicht)
- `pfad`: `docs/plan/planning/in-progress/slice-sdk-python-nats-stream-client-flaeche.md:198-208`
- `befund`: Die DoD-Zeile verlangt das Suchlauf-Ergebnis „in §7 berichtet,
  einschließlich `options.py`s Docstring"; §7 steht noch vollständig auf
  Template (`<…>`), kein anderes committetes Feld des Plans trägt den
  Suchlauf-Befund — `AGENTS.md` §3.13 verlangt Gefundenes **und**
  Nichtgefundenes in einem committeten Feld, bevor die Closure-Notiz
  geschrieben wird. Bis zu diesem Stand ist der Befund also nirgends
  committet, und die F-1/F-4-Fundstellen zeigen, dass er drei bis fünf
  Fehltreffer zu berichten hätte. Der Plan-Nachzug-Zug (dieser Diff) war
  der Ort gewesen, an dem das Feld entstehen konnte.
- `verifizierbar`: ja — Plan §7 lesen (Template-Stand);
  `AGENTS.md` §3.11/§3.13 sind die Regeln.
- `klasse`: „Suchlauf-Ergebnis nicht committet" (§3.13-Träger-Seite)

### F-6 — Realserver-Reject-Test bindet die Ursache nicht: `pytest.raises(Exception)` färbt jeden Fehlschlag denselben Marker

- `kategorie`: MEDIUM
- `quelle`: Maintainability · Reviewer-Skill „Beleg trägt seinen Satz
  nicht" (schmale Form: Ursachen-Attribution)
- `pfad`: `sdks/python/pgchangefeed/integration/test_nats_realserver.py:184-190`
- `befund`: Der Rejection-Test hält `pytest.raises(Exception)` und
  schließt nur `TimeoutError` aus; der gedruckte Marker „REJECTED
  token-rejected: …" (der vom Runner als `reject_marker` gegenprüft wird)
  behauptet damit die Ursache „token-rejected", die Assertion hält sie
  nicht — ein nicht-authentifizierungsbedingter Verbindungsfehler erzeugte
  denselben Marker und denselben End-Echo-Satz des Runners („… wurde vom
  NATS-Server abgelehnt"). Beide Vorgänger-Tests binden eng
  (`pytest.raises(PgChangeFeedUnauthorizedError)`, `grpc.RpcError`);
  `nats-py` trägt `nats.errors.AuthorizationError` als engere, verfügbare
  Bindung. Die Abwesenheitsseite ist gebunden (lehnt der Server gar nicht
  ab, scheitert der Test real rot und der Runner über `test_exit != 0`) —
  die Lücke betrifft nur die Ursachen-Attribution des Belegs.
- `verifizierbar`: ja — Assertion lesen; Gegenprobe wäre eine
  Runner-Umgebung ohne Token-Konfiguration (Verifier-Lauf).
- `klasse`: „Beleg trägt seinen Satz nicht" (Ursachen-Attribution)

### F-7 — Vakuum-Assertion hinter der Sichtbarkeits-Behauptung des Rejection-Unit-Tests

- `kategorie`: LOW
- `quelle`: Maintainability
- `pfad`: `sdks/python/pgchangefeed/tests/test_nats_stream_client.py:161`
- `befund`: Nach der tragenden Bindung (`pytest.raises(RuntimeError, …)`
  am Generator) steht „Der Fehlschlag erreicht den Generator sichtbar
  (nicht still leer)" und darauf `assert rejection is not None` — die
  Behauptung kann nie rot werden (`rejection` wurde zwei Zeilen zuvor
  literal zugewiesen); die Kommentar-Aussage trägt ausschließlich die
  `pytest.raises`-Zeile davor. (Vergleich: gRPC-Review F-2 — dort war die
  vakuum-Form die Feldvollständigkeits-Assertion selbst, HIGH; hier ist es
  eine tote Nachklammer neben einer gebundenen Assertion.)
- `verifizierbar`: ja — Assertion lesen; Mutation „Entfernen der
  except-Behandlung" färbt über `pytest.raises` rot, über den Nachlauf
  nie.
- `klasse`: „Vakuum-Assertion" (tote Nachklammer)

### F-8 — `PgChangeFeedMalformedResponseError`-Docstring nennt [`SPEC-024`](../../spec/pflichtenheft.md) nicht als Form, die ihn trägt

- `kategorie`: LOW
- `quelle`: Maintainability (Wiederholung des SSE-Review-F-4-Musters)
- `pfad`: `sdks/python/pgchangefeed/src/pgchangefeed/exceptions.py:52-58`
- `befund`: Der Docstring zählt „[`SPEC-018`](../../spec/pflichtenheft.md)/[`SPEC-022`](../../spec/pflichtenheft.md) … or a [`SPEC-021`](../../spec/pflichtenheft.md) SSE
  event" als tragende Formen; der NATS-Weg ([`SPEC-024`](../../spec/pflichtenheft.md)) kommt in diesem Diff
  neu als dritter raisender Ort hinzu (drei neue Meldungen in
  `nats_stream_client.py`), bleibt in der Aufzählung aber unerwähnt —
  dieselbe Stelle, die der SSE-Slice in seiner Fixrunde für [`SPEC-021`](../../spec/pflichtenheft.md)
  nachziehen musste, verpasst jetzt [`SPEC-024`](../../spec/pflichtenheft.md). Zweites Auftreten desselben
  Musters in diesem Package.
- `verifizierbar`: ja — Docstring lesen gegen die drei neuen
  raising-Sites.
- `klasse`: „Docstring nennt seine tragenden Formen nicht vollständig"

### F-9 — Plan-Nachzug labelt die Chronik-Zeile als „§6", sie liegt in §7 Historie

- `kategorie`: LOW
- `quelle`: Maintainability (einmalige Fehletikettierung im Plan-Nachzug)
- `pfad`: `docs/plan/planning/in-progress/slice-sdk-python-nats-stream-client-flaeche.md:138`
- `befund`: Die Plan-Nachzug-Zeile „`spec/pflichtenheft.md` §6
  Chronik-Zeile | neu" meint die Historie-Tabelle, die im Pflichtenheft
  unter „## 7. Historie" (Zeile 630) liegt; die tatsächliche §6 trägt die
  Vertragszeilen. Die Etikett-Verschiebung legt die Doppeldeutigkeit frei,
  unter der der §3-Punkt „§6 (`[SPEC-027]`-Zeile)" beim Umsetzen als durch
  die Chronik-Zeile erfüllt gelesen werden konnte (F-1) — der
  Plan-Nachzug-Text selbst trägt die Verwirrung mit.
- `verifizierbar`: ja — Abschnitts-Köpfe des Pflichtenhefts (Zeilen 613/630)
  gegen die Plan-Zeile halten.
- `klasse`: „Fehletikettierter Plan-Nachzug"

### F-10 — Beide neuen Testdateien enden ohne Zeilenumbuch am Dateiende

- `kategorie`: LOW
- `quelle`: Maintainability (einmalige Form-Abweichung, package-weit
  etabliert)
- `pfad`: `sdks/python/pgchangefeed/tests/test_nats_stream_client.py`,
  `sdks/python/pgchangefeed/integration/test_nats_realserver.py`
- `befund`: Beide neuen Dateien enden ohne abschließenden Zeilenumbuch —
  dieselbe Form wie fünf bestehende Dateien des Packages
  (`grpc_client.py`, `test_grpc_client.py`, `test_sse_client.py`,
  `test_grpc_realserver.py`, `test_sse_realserver.py`); der SSE-Review
  F-5 führte dasselbe Muster einmal als LOW. Die neuen Dateien folgen der
  umgebenden Form, kein neuer Drift — die Stelle bleibt eine
  Konsistenz-Frage für einen künftigen Aufräum-Zug, kein Slices-Fehler.
- `verifizierbar`: ja — `tail -c 1` über die Python-Dateien des Packages
  (in diesem Review-Lauf ausgeführt, siehe Negativbefund-Liste).
- `klasse`: „Fehlende EOF-Newline" (etabliertes Package-Muster, 3. Auftreten)

### F-11 — Queue-Sentinel `None` hat keinen Produzenten; die `payload is None`-Verzweigung ist unerreichbar

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `sdks/python/pgchangefeed/src/pgchangefeed/nats_stream_client.py:81,96`
- `befund`: Der Callback legt ausschließlich `message.data` in die Queue;
  der consume-Thread erzeugt den `None`-Sentinel nie, die Verzweigung
  `if payload is None: return` ist an diesem Stand unerreichbar. Defensiver
  Platzhalter ohne Produzenten — Hinweis ohne erwartete Aktion (der
  Typ `bytes | None` trägt den Sentinel in der Signatur).
- `verifizierbar`: nein — Lesebefund.
- `klasse`: „Unerreichbare Verzweigung"

---

## Negativbefunde

- geprüft, ohne Befund: **Plan-vs-Code außerhalb von F-1** — alle sieben
  §3-Zeilen und alle vier Plan-Nachzug-Zeilen sind im Diff vertreten; der
  `tests/integration/`-Pfad-Punkt ist als Abweichung deklariert
  (Plan-Nachzug Zeile 1, Geschwister-Ordner `integration/` seit dem
  Vorgänger-Slice). Still gestrichen wurde nur die §6-`[SPEC-027]`-Zeile
  (F-1).
- geprüft, ohne Befund: **Test-Anzahl** — `grep -c "^def test_"` real
  20 + 3 + 8 + 10 + 7 = 48, deckungsgleich mit der Vorgänger-Formel (41
  bei SSE) plus 7; keine Drift (§3.12 Instanz A, Nachmessen statt Übernahme).
- geprüft, ohne Befund: **[`SPEC-024`](../../spec/pflichtenheft.md)-Draht** — Subjekt-Form
  (`cdc.stream.<source_id>.>`), alle zehn Nachrichtenfelder, Verbindungsebene-
  Authentifizierung, Fire-and-Forget/Zustellgarantie „keine" — einzeln
  gegen `spec/pflichtenheft.md` §2 [`SPEC-024`](../../spec/pflichtenheft.md) gehalten; serverseitige
  Gegenprobe (`publisher.go`) trägt dieselben zehn JSON-Keys und denselben
  Subjekt-Präfix. „Kein drittes Schema" gewahrt (`StreamChange` wiederverwendet).
- geprüft, ohne Befund: **Token-/Quelle-Verkettung des Realserver-Laufs** —
  `compose.yaml` `CDC_NATS_STREAM_TOKEN: e2e-nats-stream-token` = Runner-
  Konstante; `PGCHANGEFEED_SOURCE_ID=src-e2e` = `cdc.source`-Zeile des
  Runners; Reject-Marker `REJECTED token-rejected` deckt sich zwischen
  Test-Print und Runner-Grep.
- geprüft, ohne Befund: **Import-Grenze** — eigene `grep`-Läufe über alle
  neuen Python-Dateien: kein `internal/`/`cmd/`/`gen/`-Bezug; „kein
  `examples/python/`" eingehalten (`ADR-0110` Festlegung 2, Welle-Plan §6).
- geprüft, ohne Befund: **Runner-Mechanik** — `set -euo pipefail`; Phase-
  Rückgabewert über `$(...)` mit Exit-Propagation; Marker-Polls mit
  expliziten Fehlerpfaden (Vorbild-Muster der Vorgänger-Phasen
  unverändert); `run_surface_phase`-Refaktorisierung deckt alle drei
  Call-Sites; Env-Extra-Liste trägt ihre Leerzeichen-Kopplung als
  deklarierten Kommentar („Leerzeichen-getrennte Extra-Env-Liste"), keine
  aktuelle Wert-Form verletzt sie.
- geprüft, ohne Befund: **§3.7 Kommentar-Disziplin außerhalb von F-2** —
  keine Chronik, kein Vorher/Nachher, kein abwesender Text in den neuen
  Produktions-Dateien; Slice-Kennungen nur in Test-/Skript-Provenienz
  (zulässige Form, Satzsubjekt ist der Test/Skriptzustand).
- geprüft, ohne Befund: **Handbuch** — Version 1.42 → 1.43, neue
  Änderungshistorie-Zeile (1.43) vorhanden; der neue Absatz liegt real in
  §4 „Zugriff über den NATS-Vollinhalts-Stream" (Zeilen 940–1038, vor
  „## 5"), wie die Historien-Zeile behauptet; die referenzierte
  Zehn-Felder-Tabelle steht in derselben Sektion oberhalb. Die
  Handbuch-Versionshistorie-Regel (Skill-HIGH) ist erfüllt.
- geprüft, ohne Befund: **Version-Hebung** — `pyproject.toml` 0.1.0 →
  0.2.0 (Minor, PEP 440); reales 0.2.0-Artefakt-Paar in
  `sdks/python/dist/` gemessen (20481/24639 Bytes, 2026-09-23 07:47) —
  die „real geprüft"-Sätze des Diffs sind durch die Messung gedeckt;
  `docs/user/releasing.md` bleibt korrekt (PyPI-Ist-Stand `0.1.0` ist
  historisch wahr, Tag-Push Out-of-Scope).
- geprüft, ohne Befund: **Host-lokale Pfade (§3.11)** — keiner im Diff;
  Backtick-Parität der vier geänderten Markdown-Dateien paarig.
- geprüft, ohne Befund: **Traceability** — Commit-Betreff `0f8cc4f2` nennt
  `LH-FA-SST-009`/`ADR-0110`, kein `SPEC-*`/`ARC-*` im Betreff
  (Standing-Gate-Form).
- geprüft, ohne Befund: **Sync-Wrapper-Ausgang (Plan §6 Risiko 1) am Ort
  benannt** — Modul-Docstring benennt synchrones öffentliches Vokabular und
  asyncio-Grenze im Modul („one iterator form across all four surfaces");
  kein Blocker gemäß Plan-Ausgang.
- geprüft, ohne Befund: **`make gates`** — ungepiped, Exit-Code direkt
  geprüft: `0` (in diesem Review-Lauf real ausgeführt).
- geprüft, ohne Befund: **`spec/pflichtenheft.md` §1-Prosa und §7-Chronik**
  — beide korrekt nachgezogen („deckt HTTP-API, gRPC-Stream, SSE und
  NATS-Vollinhalt", `0.2.0`, Chronik-Zeile datiert 2026-09-23 im Muster der
  drei Vorgänger-Nachzüge); stale sind nur die in F-1/F-4 genannten Stellen.
- geprüft, ohne Befund: **Workflow/Make-Wiring** —
  `.github/workflows/sdk-python-release.yml` liest die Version dynamisch
  aus `pyproject.toml`, kein hartes Versionsliteral; keine neue
  Make-/Workflow-Verkabelung nötig (Runner-Target `test-sdk-python-integration`
  existiert seit dem gRPC-Slice, Dockerfile-Stufen `build`/`integration`
  tragen die neue Testdatei über `COPY pgchangefeed/tests`/`integration`
  und die Runtime-Dep `nats-py` über `pip install -e ".[test]"`).

---

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 3 |
| MEDIUM | 3 |
| LOW | 4 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** „Zahl im Träger … driftend" · „Arbeit
überholt stehenden Träger" · „Suppression ohne Linter/ADR" · „Beleg trägt
seinen Satz nicht" / „Zitat nennt die falsche Stelle" · „Suchlauf-Ergebnis
nicht committet" · „Beleg trägt seinen Satz nicht" (Ursachen-Attribution) ·
„Vakuum-Assertion" · „Docstring nennt tragende Formen unvollständig" ·
„Fehletikettierter Plan-Nachzug" · „Fehlende EOF-Newline" · „Unerreichbare
Verzweigung"

## Verdikt

**Merge-blockierend:** ja — drei HIGH-Findings (F-1 stale
[`SPEC-027`]-Vertragszeile mit Planpunkt still gestrichen; F-2 `#noqa`-Suppression
gegen §3.2; F-3 falscher Beleg-Anker im Package-Docstring). Die Fläche selbst
ist vollständig und konsistent designt: [`SPEC-024`](../../spec/pflichtenheft.md)-Draht Feld für Feld gegen
Server und Spec gehalten, Realserver-Phase im Runner mit Token-Verkettung
und SQL-Gegenprüfung, 48 Unit-Tests real nachgemessen, reales
`0.2.0`-Artefakt-Paar gemessen, `make gates` grün (eigener, ungepiped
geprüfter Lauf). Keine der drei HIGH-Stellen berührt Korrektheit oder
Sicherheit — sie sind Träger-/Form-Stellen mit einzeiliger bis kleiner
Fixfläche.

**Übergabe:** Rückgabe-Pfeil an den Implementer für eine gezielte Fixrunde —
kein Rollen-Widerspruch, keine Konflikt-Pfad-Sequenz über den Architect nötig
(Modul 8: Konflikt-Pfad greift bei Rollen-Widerspruch oder 3× gleichem
Konflikttyp; hier liegt keines von beiden vor). Erwarteter Umfang der
Fixrunde:

1. F-1: `spec/pflichtenheft.md` §6 `[SPEC-027]`-Zeile nachziehen
   („aktuell `0.2.0`" + Vier-Verträge-Liste [`SPEC-018`]/[`SPEC-020`]/
   [`SPEC-021`]/[`SPEC-024`] — Muster der [`SPEC-026`]-/[`SPEC-028`]-Zeilen);
   die still gestrichene Plan-§3-Zeile als ausgeführt oder als Abweichung
   deklarieren.
2. F-2: das `# noqa: BLE001`-Token entfernen; die Zusage-Kommentar-Zeile
   hinter dem `except` bleibt (die Begründung trägt der Kommentar, nicht
   das Token).
3. F-3: den Docstring-Anker auf die tatsächliche Entscheidungsschicht
   re-anchorn (Welle-Plan §6 — oder eine reine Zustandsform ohne
   Festlegungs-Zitat); Entscheidung selbst ist real und regelkonform
   getroffen, nur der Anker ist falsch.
4. F-4: die fünf benannten Träger-Stellen nachziehen (Muster der
   Schwester-Slices `d86d1965`/`c8c9e3ae`).
5. F-5: den Suchlauf-Befund (Gefundenes **und** Nichtgefundenes, inklusive
   der Korrekturen aus 1./4.) in §7 oder einem eigenen Plan-Feld
   committen, bevor die Closure-Notiz geschrieben wird.
6. F-6: den Rejection-Test an die engere Bindung heben (z. B.
   `nats.errors.AuthorizationError`) oder die Ursachen-Aussage des Markers
   auf das halten, was die Assertion misst.
7. F-7/F-8/F-9/F-10: kleinflächige Korrekturen (vakuum-Nachklammer
   entfernen oder aussagekräftig machen; Docstring-Enumeration um
   [`SPEC-024`] ergänzen; „§6" → „§7" im Plan-Nachzug; EOF-Newline nach
   Entscheidung der Fixrunde — Muster im Package etabliert).

Die Finding-Klassen gehen in die Slice-Closure §7 und von dort in den
Steering-Loop-Zähler. Dieser Report ist ein Lauf-Beleg; er ersetzt keine
Verifikation gegen die volle DoD — das bleibt Verifier-Aufgabe (Modul 11),
insbesondere der reale `make test-sdk-python-integration`-Lauf gegen die
NATS-Fläche (Welle-Closure-Trigger) bleibt dessen Nachweis.

**DoD-Checkbox-Nachzug:** entfällt — dieses Verdikt führt zu einer
Fixrunde; die DoD-Zeile „Review durchgeführt, Report unter `docs/reviews/`
liegt vor" wird regulär bei Schritt 21 des Implementer-Workflows
nachgezogen (Reviewer-Skill §DoD-Checkbox-Nachzug ohne Fixrunde, Grenzfall
„nur ohne Fixrunde" liegt nicht vor).