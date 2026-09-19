# Verifikationsbericht: slice-sdk-python-http-client-flaeche — 2026-09-19

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?" gegen den
DoD-Vertrag
(`docs/plan/planning/in-progress/slice-sdk-python-http-client-flaeche.md`
§2/§3/§6/§7), in frischem Kontext. **Nicht** gegen den Diff als solchen
(Reviewer-Aufgabe, abgeschlossen mit
[`review-slice-sdk-python-http-client-flaeche.md`](review-slice-sdk-python-http-client-flaeche.md)
inklusive Fixrunden-Nachprüfung) und **nicht** gegen realen Bedarf
(Validator, hier nicht ausgelöst).

**Gegenstand:** vier Commits `0324b9e6..HEAD` (`8f76564e`, `225b5070`,
`1d971e6b` reine Moves/Slice-Metadaten; `314b51f6` Implementer-Inhalt;
`b9177f48` Review-Report; `d0da688b` Fixrunde; `3d884534`
Fixrunden-Nachprüfung), Slice `slice-sdk-python-http-client-flaeche`,
Welle `welle-sdk-python-lh-fa-sst-009`, `LH-FA-SST-009`, `ADR-0107`.

**Eingangs-Kontext (eigen gelesen, nicht aus Bericht übernommen):**
Slice-Plan §1/§2/§3/§6/§7/§8 vollständig, Review-Report vollständig
inkl. Fixrunden-Nachprüfung, `ADR-0107` Festlegung 1/3/§Konsequenzen,
`spec/pflichtenheft.md` §2 `SPEC-018`/`SPEC-022` im Volltext, alle sechs
Python-Quelldateien (`http_client.py`, `models.py`, `exceptions.py`,
`options.py`, `__init__.py`, `test_http_client.py`), `pyproject.toml`,
`sdks/python/README.md`, `docs/user/benutzerhandbuch.md`-Diff.

---

## 1. DoD-Vertrag (§2) — jede Zeile einzeln geprüft

### 1.1 „Öffentliche Client-Klasse mit einer Methode je der neun
Port-gedeckten Fähigkeiten … plus Changes-Lesen … Fehler-Antwortform …
typisierte Exception-Hierarchie … Bearer-Token bei Konstruktion"

Eigenes Nachzählen der öffentlichen Methoden in `http_client.py`:
`register_consumer`, `acknowledge_consumer`, `get_consumer_position`,
`remove_consumer`, `enable_table`, `disable_table`, `get_status`,
`list_tables`, `run_retention`, `read_changes` — **genau zehn**, keine
fehlt, keine zusätzliche (kein Diagnose-/Health-Endpunkt, `SPEC-018`
schließt diese ausdrücklich aus, §1 des Plans hält das ein).

`__init__(self, client: httpx.Client, options: ClientOptions)` — Bearer-
Token kommt ausschließlich aus `self._options.api_token`
(`_headers()` baut `Authorization: Bearer <token>` je Aufruf neu), kein
Modul-Level-/globaler Zustand.

**Ergebnis: berechtigt auf `[x]`.**

### 1.2 „Eigene Tests (`pytest`) decken je Fähigkeit mindestens den Happy
Path und die Auth-Boundary … ab — netzlos"

Eigenes Zählen in `test_http_client.py` (338 Zeilen komplett gelesen):
10 Happy-Path-Tests, einer je Fähigkeit (inkl. `read_changes` und ein
Zusatztest für den leeren Treffer), 2 Auth-Boundary-Tests (401 unbekanntes
Token, 403 `reader` gegen `admin`-Endpunkt über `remove_consumer` als
Repräsentant), 4 Statuscode-Mapping-Tests (400/404/500/undokumentiert
409), 1 Fallback-Test (Fehlerkörper ohne `error`-Feld), 2 Malformed-2xx-
Tests, 1 leerer-Treffer-Test = **20 Tests**, `test_options.py` = 3
vorbestehende Tests, Summe 23 — deckt sich mit dem realen Docker-Testlauf
(siehe §5).

**Auth-Boundary „je Fähigkeit" — eigene Einschätzung, kein Fund:** Die
zwei Auth-Boundary-Tests laufen nur gegen `remove_consumer`
(repräsentativ), nicht gegen jede der zehn Methoden einzeln. Geprüft, ob
das ein DoD-Verstoß ist: `_handle` ist der **einzige** Ort, an dem
Statuscode-Fehler entstehen, und **jede** `_get`/`_post`-Stelle aller
zehn Methoden ruft ausschließlich `self._handle` auf (eigene Lektüre der
kompletten Datei, siehe 1.6) — es gibt keinen methodenspezifischen
Zweitpfad, den ein Auth-Test pro Methode zusätzlich prüfen müsste. Exakt
dasselbe Muster trägt bereits das C#-Geschwister-Package
(`PgChangeFeedHttpClientAuthBoundaryTests.cs`: zwei generische Auth-Tests,
kein Test je Methode) — die Interpretation ist repo-konsistent, keine
Python-spezifische Lücke. **Kein neuer Fund**, DoD-Checkbox berechtigt.

**Ergebnis: berechtigt auf `[x]`.**

### 1.3 „Kein Import aus `internal/**`/`cmd/**`"

Eigener Lauf: `grep -rn "internal/\|cmd/" sdks/python/` → kein Treffer,
Exit `1`. **Ergebnis: berechtigt auf `[x]`.**

### 1.4 „`make gates` grün"

Siehe §5 unten. **Ergebnis: berechtigt auf `[x]`.**

### 1.5 „Review durchgeführt … Fixrunde durchlaufen … Fixrunden-Nachprüfung
bestätigt beide Findings real behoben"

Eigene Lektüre des vollständigen Review-Reports inkl.
§Fixrunden-Nachprüfung: 1 HIGH (F-1, `pyproject.toml`-Kommentar behauptet
„folgt erst" für eine real bereits gelieferte Fläche, `AGENTS.md` §3.13),
1 MEDIUM (F-2, Commit-Message-Testzahl-Kategorisierung). Eigene,
unabhängige Prüfung des F-1-Fixes: `pyproject.toml:22-24` gelesen — trägt
jetzt „wird von der HTTP-API-Client-Flaeche
(src/pgchangefeed/http_client.py) fuer die HTTP-Requests genutzt", keine
Slice-Namen, keine „folgt erst"-Formulierung. Eigener erweiterter
Suchlauf (unabhängig vom Reviewer-Suchlauf formuliert):

```
grep -rniE "follow-up|added by|folgt erst|folgt mit|noch nicht|steht noch aus|is added|will be added|slice-sdk-python" sdks/python/
```

Ergebnis: drei Treffer, alle bereits vom Reviewer identisch geprüft und
zutreffend (`README.md:9`, `__init__.py:7` — beide ausschließlich
gRPC/SSE/NATS betreffend, `ADR-0107` Festlegung 1 konform; `Dockerfile:8`
— betrifft eine noch nicht existierende `pack`-Stufe eines anderen
Folge-Slice). Keine vierte, vom Reviewer übersehene Stelle gefunden — F-1
ist real geschlossen.

F-2: eigenes Nachzählen (siehe 1.2) bestätigt die im Plan-Nachzug
dokumentierte Aufschlüsselung (20 Kategorisierungs-Tests + 3
vorbestehende); die Commit-Message-Zahl „23" ist als Gesamtsumme korrekt,
die Kategorisierungs-Zuordnung war ungenau — sachgerecht dokumentiert
statt per Amend korrigiert, kein neuer Handlungsbedarf.

**Ergebnis: berechtigt auf `[x]`.**

### 1.6 Eigener, vollständiger Blick auf `_handle`

Komplette Datei `http_client.py` (258 Zeilen) gelesen: `_get`/`_post`
rufen ausschließlich `self._handle(response, mapper)` auf; `_handle`
deckt (a) jeden Nicht-2xx-Statuscode über `_build_error`
(400/401/403/404/500 → typisierte Klasse, jeder andere Code →
`PgChangeFeedUnexpectedStatusError`), (b) den malformten-2xx-Body-Fall
über zwei separate `try/except`-Blöcke — ungültiges JSON
(`except ValueError`) und valides JSON ohne erwartetes Feld
(`except (KeyError, TypeError)`), beide → typisierte
`PgChangeFeedMalformedResponseError`. Da jede der zehn Methoden
ausschließlich über `_get`/`_post` geht, deckt `_handle` **strukturell
alle zehn Methoden einheitlich** ab — kein methodenspezifischer Zweitpfad
gefunden, der ihn umginge. Bestätigt eigenständig, nicht nur aus dem
Review-Report übernommen.

### 1.7 „Doku-Update: `docs/user/benutzerhandbuch.md`"

Eigener `git diff 0324b9e6..HEAD -- docs/user/benutzerhandbuch.md`:
Kopfzeile `Version: 1.34` → `1.35`; neue Zeile in `### Änderungshistorie`:
„1.35 | 2026-09-19 | Python-SDK-Hinweis für die HTTP-Oberfläche ergänzt
…" — inhaltlich zutreffend: nennt „dieselben zehn Fähigkeiten", die
reale Methodenzahl stimmt (siehe 1.1). Neuer Absatz in §4 nennt
`PgChangeFeedHttpClient`, `pip install pgchangefeed`, die Ausklammerung
von gRPC/SSE/NATS — deckt sich mit `ADR-0107` Festlegung 1 und dem
tatsächlichen Lieferumfang. `sdks/python/README.md` §Status wurde
proaktiv vor dem ersten Review nachgezogen (eigene Lektüre bestätigt:
„provides … a full HTTP API client surface", kein „folgt erst" mehr für
die HTTP-Fläche selbst).

**Ergebnis: berechtigt auf `[x]`.**

### 1.8 Reconciliation, Beobachtungs-Register, Closure-Notiz, §6-Risiko-
Ausgänge, drei Paarungen

Der Slice liegt weiterhin unter `in-progress/` — §7 „Closure-Notiz" trägt
weiterhin `<…>` (ungefüllt), die drei zugehörigen `[ ]`-Checkboxen sind
konsistent mit diesem Zustand (Closure-Pflichten, keine Liefer-Punkte,
regulär bei Slice-Closure fällig, nicht Gegenstand dieser
Zwischen-Verifikation vor `done/`). Reconciliation „entfällt" ist
zutreffend begründet (keine Reconciliation-Datei in diesem Repo, real
über `find docs/plan/planning -iname "reconciliation*"` bestätigt: kein
Treffer außerhalb der Referenz im Plantext selbst). Beobachtungs-Register
„keine Beobachtung angefallen" — eigene Prüfung: Der Plan-Nachzug (§3)
selbst dokumentiert zwei reale Beobachtungen (Träger-Nachzug-Suchlauf vor
Review erfolgreich, Import-Grenzen-Grep-Kollision vermieden) sowie den
F-1-Fund in der Fixrunde — diese sind bereits in `AGENTS.md` §3.13s
`BEO-PGC/arbeit-ueberholt-stehenden-traeger`-Familie referenziert (der
Review-Report nennt sie explizit als „dieselbe Fehlerklasse wie F-1 des
C#-Geschwister-Reviews"); kein neuer, noch nicht verkörperter
Beobachtungs-Typ gefunden, DoD-Text ist zutreffend.

Die §6-Risiken (drei Einträge — Draht-Treue ohne Referenz-Client,
Exception- vs. Result-Design, `httpx`-Mock-Verdeckungsrisiko) tragen alle
drei explizit „weiter offen" mit nachvollziehbarer, strukturell
begründeter Auflösung (Folge-Slice/Design-Endgültigkeit erst ab `1.0.0`);
das ist eine zulässige Ausgangsklasse (Baseline-Regelwerk
`modul-05-planning-harness.md` §Offene Risiken werden bei Closure
aufgelöst) — formaler Häkchen-Nachzug bleibt regulär Closure-Sache.

**Kein DoD-Blocker in dieser Gruppe.**

---

## 2. Eigenständige Feld-für-Feld-Gegenprobe `models.py` gegen
`SPEC-018`/`SPEC-022` (unabhängig vom Review-Report, zweite Instanz)

`spec/pflichtenheft.md` Zeilen 333–486 (`SPEC-018`/`SPEC-022`, Volltext)
und `models.py` (261 Zeilen, Volltext) nebeneinander gelesen, jede
Dataclass einzeln gegen die Spec-Tabellenzeile gehalten:

| Fähigkeit | Spec-Feldform (Request/Response) | `models.py` | Abweichung? |
|---|---|---|---|
| `RegisterConsumer` | Req `{consumer_id, name}`; Resp `{consumer_id, name, already_registered}` | `RegisterConsumerRequest{consumer_id, name}`; `RegisterConsumerResponse{consumer_id, name, already_registered}` | keine |
| `AcknowledgeConsumer` | Req `{consumer_id, source_id, offset}`; Resp `{consumer_id, source_id, offset}` | identisch | keine |
| `GetConsumerPosition` | Resp `{consumer_id, source_id, offset, acknowledged}` | `ConsumerPositionResponse{consumer_id, source_id, offset, acknowledged}` | keine |
| `RemoveConsumer` | Req `{consumer_id}`; Resp `{consumer_id, removed}` | `RemoveConsumerResponse{consumer_id, removed}` | keine |
| `EnableTable` | Req `{source, schema, table, table_id, schema_version_id, version, publication}`; Resp `{table_id, source, schema, table, already_enabled}` | identisch (Feld für Feld, inkl. `version: int` für int64 ≥ 1) | keine |
| `DisableTable` | Req `{source, schema, table, publication}`; Resp `{removed, retained}` | identisch | keine |
| `GetStatus` | Resp `{enabled, retained}` | `TableStatusResponse{enabled, retained}` | keine |
| `ListTables` | Resp `{tables: [{table_id, source, schema, table}], retained: [...]}` | `ListTablesResponse{tables: list[TableInfo], retained: list[TableInfo]}`, `TableInfo{table_id, source, schema, table}` — `retained` trägt korrekt dieselbe Form wie `tables` | keine |
| `RunRetention` | Req `{source, min_age_nanos}`; Resp `{deleted}` | identisch | keine |
| `ReadChanges` | Resp-Element: `commit_position, change_id, transaction_id, source_table_id, schema, table, sequence, operation, old_image, new_image, schema_version, committed_at` (zwölf Felder, `old_image`/`new_image` statt `old_data`/`new_data`) | `Change{commit_position, change_id, transaction_id, source_table_id, schema, table, sequence, operation, old_image, new_image, schema_version, committed_at}` — exakt zwölf Felder, exakt dieselben Namen inkl. `old_image`/`new_image` | keine |
| `ReadChangesResponse` | `{"changes": [...]}`, leere Liste statt `null`/`404` | `ReadChangesResponse{changes: list[Change] = field(default_factory=list)}` | keine |

**Ergebnis der zweiten unabhängigen Instanz: keine Abweichung gefunden.**
Diese Prüfung deckt sich mit dem Review-Ergebnis, wurde aber komplett neu
und ohne Rückgriff auf den Review-Text durchgeführt (eigene
Spec-Zeilen-Lektüre plus eigene Datei-Lektüre in dieser Sitzung). Die
Bedeutung dieser zweiten Instanz ist real, nicht dekorativ: `ADR-0107`
§Kontext benennt ausdrücklich das Fehlen eines Python-Referenz-Clients
als strukturelles Risiko für genau diese Klasse von Fehler (§6 Risiko 1
des Plans) — zwei unabhängige, komplette Feld-für-Feld-Durchläufe (Review
+ diese Verifikation) senken dieses Risiko, ersetzen aber laut Plan
weiterhin nicht den später fälligen realen Rundlauf-Beleg.

---

## 3. `make gates` real, ungepiped, in dieser Sitzung ausgeführt

```
$ git log -1 --format=%H
3d884534...
$ make gates > /tmp/verifier-make-gates.log 2>&1; echo $?
0
```

Einzelbelege aus demselben Lauf (Log vollständig eingesehen):

| Gate | Ergebnis |
|---|---|
| `baseline-verify` | `v6.9.0 OK` |
| `docs-check` | `d-check: 828 Datei(en) geprüft, 0 Befund(e)` |
| `a-check` | `gesamt: 0 Befund(e)` |
| `commit-traceability` | `OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` |
| `coverage-gate` | `OK — Coverage 82.80% erfüllt Schwelle 80%` |
| `generated-sync` | `OK — das committete Erzeugnis ist byte-gleich der Ausgabe des gepinnten Generators` |

**Ergebnis: DoD-Checkbox „`make gates` grün" berechtigt auf `[x]`.**

## 4. Vierte unabhängige Docker-Build-Bestätigung

```
$ docker build --no-cache -f sdks/python/Dockerfile sdks/python > /tmp/verifier-python-build.log 2>&1; ec=$?
DOCKER_BUILD_EXIT=0
...
tests/test_http_client.py ....................  [ 86%]
tests/test_options.py ...                        [100%]
23 passed in 0.04s
```

Exit-Code direkt (ungepiped) geprüft, `0`. 23/23 Tests grün — 20 in
`test_http_client.py`, 3 in `test_options.py`, deckt sich mit den
Zählungen in §1.2. Dies ist die **vierte** unabhängige Bau-Bestätigung
nach Implementer (Plan-Nachzug), Reviewer-Erstlauf und
Reviewer-Fixrunden-Nachprüfung — jeweils mit demselben Ergebnis, ohne
dass diese Sitzung eine der drei vorherigen Behauptungen unbesehen
übernommen hätte. Test-Image danach per `docker image prune` entfernt.

## 5. ADR-0107 §5/§Konsequenzen — unberührte Pfade real geprüft

```
$ git diff 0324b9e6..HEAD --stat -- .a-check.yml spec/architecture.md docs/user/version.md sdks/csharp/
(leer)
```

Kein Treffer — alle vier genannten Pfade sind über den vollständigen
Diff-Bereich dieses Slices unberührt, wie `ADR-0107` §6 „Was diese ADR
nicht ändert" es für a-check, `spec/architecture.md` und
`docs/user/version.md` festlegt, sowie §Konsequenzen Folgepflicht 2 (kein
Umzug/keine Zusammenlegung von `sdks/csharp/**`).

Vollständiger Diff-Stat (`0324b9e6..HEAD`, zehn Dateien) zusätzlich
gelesen: ausschließlich `sdks/python/**`,
`docs/plan/planning/in-progress/slice-sdk-python-http-client-flaeche.md`,
`docs/reviews/review-slice-sdk-python-http-client-flaeche.md` und
`docs/user/benutzerhandbuch.md` — kein Treffer außerhalb des erwarteten
Slice-Umfangs.

## 6. `pyproject.toml` — F-1-Fix eigenständig gegengelesen

Eigene Lektüre von `sdks/python/pgchangefeed/pyproject.toml` Zeile 1–37
komplett: Zeile 22–24 trägt „Einzige Fremdabhaengigkeit (`ADR-0107`
Festlegung 1: httpx, keine weitere, keine schwere Fremdabhaengigkeit) --
wird von der HTTP-API-Client-Flaeche
(src/pgchangefeed/http_client.py) fuer die HTTP-Requests genutzt." — kein
Slice-Name, keine „folgt erst"/„folgt mit"-Formulierung mehr, eine
indikative Ist-Zustands-Beschreibung (`AGENTS.md` §3.7-konform).
**Bestätigt: F-1 real und dauerhaft behoben.**

---

## Verdikt

**DoD erfüllt** für den aktuellen `in-progress`-Stand des Slice — die
Closure-Pflichten in §2 (Closure-Notiz, formaler Risiko-Ausgangs-Häkchen-
Nachzug, drei Paarungen bei Welle-Closure) sind bewusst noch offen und
kein Gegenstand dieser Verifikation. Alle beauftragten Prüfpunkte wurden
real und unabhängig nachgemessen, nicht aus Bericht oder Commit-Message
übernommen:

1. Alle zehn öffentlichen Methoden real nachgezählt, exakt `SPEC-018`
   (neun) + `SPEC-022` (`ReadChanges`) — keine fehlt, keine zusätzliche.
2. Testabdeckung real nachgezählt: 10 Happy-Path (je Fähigkeit) + 2
   Auth-Boundary (repräsentativ über `_handle`, repo-konsistent mit dem
   C#-Geschwister-Package) + 4 Statuscode-Mapping + 1 Fallback + 2
   Malformed-2xx + 1 Leerer-Treffer = 20, plus 3 vorbestehende
   `test_options.py`-Tests = 23, deckt sich mit dem eigenen Docker-Lauf.
3. `grep -rn "internal/\|cmd/" sdks/python/` — kein Treffer.
4. **Zweite, unabhängige Feld-für-Feld-Gegenprobe** von `models.py` gegen
   `SPEC-018`/`SPEC-022`: **keine Abweichung gefunden** — für alle zehn
   Fähigkeiten plus die zwölf `Change`-Felder, inklusive
   `old_image`/`new_image`-Namensgebung und `ListTablesResponse.retained`.
5. `_handle` komplett gelesen: deckt strukturell alle zehn Methoden
   einheitlich ab, inklusive des malformten-2xx-Body-Falls über zwei
   getrennte `except`-Zweige.
6. `git diff 0324b9e6..HEAD --stat` gegen `.a-check.yml`,
   `spec/architecture.md`, `docs/user/version.md`, `sdks/csharp/**` —
   leer, alle vier unberührt.
7. `pyproject.toml:22-24` eigenständig gelesen: F-1 real und dauerhaft
   behoben, kein Slice-Name, keine veraltete „folgt erst"-Formulierung.
8. `make gates` eigenständig, ungepiped, `EXIT=0` — alle sechs Gates
   grün (baseline-verify, docs-check 828/0, a-check 0, commit-
   traceability OK, coverage-gate 82.80% ≥ 80%, generated-sync OK).
9. `docker build --no-cache -f sdks/python/Dockerfile sdks/python` —
   vierte unabhängige Bau-Bestätigung, Exit `0`, 23/23 Tests grün.
10. `docs/user/benutzerhandbuch.md`: Version 1.34→1.35 korrekt,
    Changelog-Zeile inhaltlich zutreffend (Methodenzahl „zehn" stimmt).

**Freigabe an den Planner:** Der Slice ist DoD-konform für seinen
aktuellen `in-progress`-Stand. Vor Closure (`git mv` nach `done/`)
bleiben die im Plan selbst bereits als offen geführten Closure-Pflichten
zu erfüllen (§7 Closure-Notiz mit Steering-Loop-Lerneintrag, formaler
Häkchen-Nachzug für „Jedes Risiko aus §6 trägt einen Ausgang", die drei
Paarungen bei der Closure von `welle-sdk-python-lh-fa-sst-009`).
