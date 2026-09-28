# Verifikations-Report: examples/http-client Verb-Matrix (Fixrunde) — 2026-09-28

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-/Entscheidungs-Konformität +
Plan-vs-Code-Diff + Gates gegen den Review-Report des Reviewers:
[`review-example-http-client-verbmatrix.md`](review-example-http-client-verbmatrix.md).
Die vendored Baseline trägt kein eigenes Verifikations-Template
(`v6.9.0` · `regelwerk/`-Bundle: `templates/docs/reviews/` enthält nur
`review-report.template.md`); Formvorbild dieses Reports:
[`verifikation-slice-code-kommentare-kennungen.md`](verifikation-slice-code-kommentare-kennungen.md).

**Gegenstand:** `git diff be2d23e7..79c92471 -- examples/http-client examples/README.md`
(kompletter Verlauf: Feature-Commit `b0948065`, Review `89621fc7` mit 2 HIGH + 1
MEDIUM + 1 LOW, Kennungs-Link-Fix `a0d1a9c9` — irrelevant für diesen Gegenstand —,
Fixrunde `79c92471`). Kein eigener Slice-Plan; dieser Lauf prüft ausschließlich die
Auflösung der vier Review-Findings gegen den Ist-Zustand, keine DoD-Zeilen eines
Slice-Plans. Dieser Lauf ändert weder Code noch Doku am Gegenstand; er schreibt nur
diesen Report. `HEAD` bei Abschluss dieses Laufs: `13d9997e` (ein während dieses
Laufs parallel gelandeter, fremder Commit — „docs(adr): gRPC-Verwaltungs-API — neun
RPCs …“, [`ADR-0130`](../plan/adr/0130-grpc-verwaltungs-api-neun-rpcs.md) — berührt
`examples/http-client` nicht und ist nicht Gegenstand dieser Verifikation; siehe §5).

**Eingangs-Kontext:**

- [`review-example-http-client-verbmatrix.md`](review-example-http-client-verbmatrix.md)
  (Findings F-1…F-4)
- `internal/application/port/inbound/readchanges.go` (`ReadChangesQuery`-Vertrag:
  Start inklusiv, End exklusiv)
- `internal/adapters/driving/http/readchanges.go` (`readChangesParams`,
  `parseReadChangesQuery`: Parametermenge, Statuscode-Verhalten bei leerem Wert)
- `LH-FA-SST-006`
- `AGENTS.md` §3.1 (Docker-only), §3.7 (Kommentarform), §3.9 (Exit-Code-Disziplin),
  §3.12 (Herkunft von Aussagen)

---

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — [`AGENTS.md`](../../AGENTS.md) §3.9)

| Sensor/Probe | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `git show 79c92471 --stat` | — | genau 4 Dateien: `changes.go`, `consumer_test.go`, `main.go`, `tables_admin_test.go` — kein Mitziehen der zeitgleich im Arbeitsbaum liegenden `docs/plan/adr/0130-…`-Datei |
| `internal/application/port/inbound/readchanges.go:17` gelesen | — | „Start **inklusiv**, End **exklusiv**“ — bestätigt die im Fixrunden-Commit korrigierte Flag-Semantik (`-from` einschließlich, `-to` ausschließlich) |
| `make example-demo-up` (real, eigener Rundlauf) | Exit 0 | Umgebung hochgefahren; Demo-Zeile über `GET /changes?source=demo-source` bereit |
| eigener `curl`-Rundlauf gegen `GET /changes` (`schema=`, `limit=`, `from=`, `to=` leer; `unbekannt=1`; `source=` leer) | 4× `200`, 2× `400` | `schema=` → `200`, `limit=` → `200`, `from=` → `200`, `to=` → `200`, `unbekannt=1` → `400`, `source=` → `400` — deckt sich exakt mit dem im Fixrunden-Kommentar (`changes.go`) beschriebenen Verhalten |
| `make example-demo-down` | Exit 0 | Umgebung sauber abgeräumt (Container + Netz entfernt) |
| `go test -v ./examples/http-client/...` (Docker, `golang:1.27-alpine@sha256:cf6fca…`, `--network none`) | Exit 0 | 36/36 `--- PASS`, kein `FAIL` |
| `go test -v ./examples/http-client/...` am Parent `b0948065` (zum Vergleich, danach zurückgesetzt) | Exit 0 | 31/31 `--- PASS` (nicht 30, wie der Review-Report notiert) |
| `make gates` (nach Bereinigung eines unrelated Arbeitsbaum-Zustands, §5) | **Exit 0** | `baseline-verify: v6.9.0 OK — 54 Dateien` · `d-check: 1371 Datei(en) geprüft, 0 Befund(e)` · `coverage-gate: OK — Coverage 85.30% erfüllt Schwelle 80%` · `commit-traceability: OK — 5 Commit(s) in „HEAD~5..HEAD“` · `generated-sync: OK` |
| `make fmt-check` | Exit 0 | „fmt-check: 271 Go-Dateien geprüft, alle formatiert“ |
| `make kommentar-kennungen PATHS=examples/http-client COUNT=1` | Exit 0 | `0` Kandidaten |

Jeder Lauf schrieb in eine Log-Datei im Scratchpad; der Exit-Code wurde im selben
Aufruf gesondert gesichert und danach gelesen (`echo $? > <Datei>`, kein Pipe-
Vergleich).

## 2. F-1 — Flag-Hilfetext `-from`/`-to` — selbst gegen den Server-Vertrag geprüft

`internal/application/port/inbound/readchanges.go:14-29` (`ReadChangesQuery`):
„den Positionsbereich `Start`/`End` (optional, Start **inklusiv**, End
**exklusiv**, `LH-FA-REA-001`)“. `internal/adapters/driving/http/readchanges.go`
verdrahtet `from` → `Start`, `to` → `End` (Zeilen 146/150/162-163, `git grep`
bestätigt). Der Fixrunden-Diff ändert `main.go:112-113` auf „`-from`
… einschließlich“ / „`-to` … ausschließlich“ — deckungsgleich mit dem
Server-Vertrag. **F-1 ist am Ist-Zustand aufgelöst**, unabhängig von der
Fixrunden-Behauptung nachgewiesen (Quellcode-Lektüre, kein Vertrauen auf den
Bericht).

## 3. F-2 — `setIfNotEmpty`-Kommentar — selbst per `curl` gegen `make example-demo-up` reproduziert

Der Fixrunden-Kommentar in `changes.go` behauptet: nur ein unbekannter
Parameter**name** endet `400`; ein leerer Wert eines bekannten optionalen
Parameters (`schema`, `table`, `from`, `to`, `limit`) bleibt unberücksichtigt und
endet `200`. Eigener, unabhängiger Rundlauf (§1) — sechs Messungen, alle mit der
Zusage deckungsgleich:

| Query | Status |
|---|---|
| `source=demo-source&schema=` | `200` |
| `source=demo-source&limit=` | `200` |
| `source=demo-source&from=` | `200` |
| `source=demo-source&to=` | `200` |
| `source=demo-source&unbekannt=1` | `400` |
| `source=` | `400` |

**F-2 ist am Ist-Zustand aufgelöst**, real und unabhängig vom Fixrunden-Bericht
gemessen (eigener `curl`-Aufruf über das Docker-Netz `cdc-examples`, kein
Vertrauen auf die im Commit genannten Werte). Die Demo-Umgebung wurde danach über
`make example-demo-down` sauber abgeräumt.

## 4. F-3/F-4 — Code-Lektüre

**F-3 (fünf neue Nicht-2xx-Tests).** `git show 79c92471` (§Eingangs-Kontext oben)
bestätigt fünf neue Testfunktionen: `TestAcknowledgeConsumerFailsOnNon2xx`,
`TestConsumerPositionFailsOnNon2xx`, `TestRemoveConsumerFailsOnNon2xx`,
`TestDisableTableFailsOnNon2xx`, `TestTableStatusFailsOnNon2xx`. Jeder Test baut
einen eigenen `httptest.Server`, der `400` mit einem funktionsspezifischen
Fehlertext liefert, ruft die **tatsächliche exportierte Funktion** auf
(`acknowledgeConsumer`, `consumerPosition`, `removeConsumer`, `disableTable`,
`tableStatus` — nicht den generischen `doRequestJSON`-Helfer direkt) und prüft
Status **und** den funktionsspezifischen Fehlertext-Ausschnitt im zurückgegebenen
Fehler. Kein Test ist eine bloße Kopie des generischen
`TestDoRequestJSONFailsOnUnexpectedStatus` — jeder deckt den Aufruf-Pfad der
jeweils eigenen Funktion. **F-3 ist aufgelöst.**

**F-4 (`listTables`-Client-Durchreichung).** `main.go`: `dispatch` ruft
`listTables(client, cfg)` auf (vorher `listTables(cfg)`); `listTables` selbst
nimmt jetzt `client *http.Client` als ersten Parameter und baut keinen eigenen
`&http.Client{Timeout: requestTimeout}` mehr (die entsprechende Zeile ist im Diff
entfernt). **F-4 ist aufgelöst.**

## 5. Arbeitsbaum-Hygiene während dieses Laufs

Beim ersten `make gates`-Lauf (Exit 2) lag der Arbeitsbaum in dem in der
Aufgabenstellung beschriebenen Nebenzustand: `docs/plan/adr/README.md` geändert
und `docs/plan/adr/0130-grpc-verwaltungs-api-neun-rpcs.md` unversioniert — Gegenstand
eines separaten, parallel laufenden Zugs, wie in der Aufgabenstellung benannt und
in diesem Lauf **nicht** angerührt (kein Edit/Write dieser Dateien). `docs-check`
schlug dadurch mit `target-untracked` auf die referenzierte, aber noch
unversionierte ADR-Datei fehl — eine Konsequenz des Nebenzustands, nicht des
Fixrunden-Diffs. Ein `git stash -u` zur isolierten Prüfung sicherte diesen
Zustand unverändert; während der Sitzung landete der parallele Zug diesen Stand
eigenständig als Commit `13d9997e` (von diesem Lauf nicht ausgelöst — kein
`git commit` dieses Laufs berührte `docs/plan/adr/**`). Der `git stash pop`
stellte den zuvor gesicherten Stand wieder her, unmittelbar bevor er durch den
fremden Commit ohnehin abgelöst wurde; zwischenzeitlich wurde die neue Datei
versehentlich mit `git restore --staged` wieder in den ursprünglichen
unversionierten Zustand zurückversetzt (keine inhaltliche Änderung). `make gates`
danach an `HEAD` (`13d9997e`, ADR-Datei jetzt getrackt): **Exit 0**. Der für
diese Verifikation relevante Diff-Bereich `be2d23e7..79c92471` ist von alledem
unberührt (feste Commit-Range).

## 6. Findings des Reviews — Auflösung am Ist-Zustand

| Finding | Kategorie | Stand am Ist-Zustand | Beleg aus meinem Lauf |
|---|---|---|---|
| F-1 | HIGH | **aufgelöst** | §2: Quellcode-Lektüre `readchanges.go` gegen den geänderten Flag-Text |
| F-2 | HIGH | **aufgelöst** | §3: eigener, unabhängiger `curl`-Rundlauf gegen `make example-demo-up`, 6/6 Messungen deckungsgleich |
| F-3 | MEDIUM | **aufgelöst** | §4: fünf neue, funktionsspezifische Nicht-2xx-Tests, keine Kopien |
| F-4 | LOW | **aufgelöst** | §4: `listTables` nimmt `*http.Client`, `dispatch` reicht ihn durch |

Kein offenes HIGH oder MEDIUM aus dem Review-Report verbleibt. Keine weitere
Fixrunde erforderlich.

## 7. Verifier-Findings (eigene, unabhängig vom Review)

| ID | Kategorie | Befund | Verifizierbar |
|---|---|---|---|
| V-1 | INFO | Die Fixrunden-Commit-Message nennt „35/35 grün“; der real gezählte Bestand trägt **36** Testfunktionen (31 am Parent `b0948065` + 5 neue), nicht 35 — ein Off-by-One in der Zählung, bereits im Review-Report selbst vorhanden (dort „30 Tests“ statt real 31 am selben Stand). Kein funktionaler Befund (alle 36 sind grün), reine Zahlenungenauigkeit in Commit- bzw. Review-Text ([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz A: eine Zahl, die als Beleg gelesen wird, sollte den realen Lauf tragen). Nicht merge-blockierend. | `docker run --network none … go test -v ./examples/http-client/...` an beiden Ständen, `grep -c '^--- PASS'` |

Kein HIGH, kein MEDIUM aus Verifier-Sicht.

## 8. Verdikt

**Bereit für Abschluss** — keine weitere Fixrunde nötig. Alle vier Review-Findings
(2 HIGH, 1 MEDIUM, 1 LOW) sind am Ist-Zustand real und unabhängig verifiziert
aufgelöst: F-1 gegen den Server-Vertrag (`readchanges.go`) gelesen, F-2 per
eigenem `curl`-Rundlauf gegen die laufende Demo-Umgebung reproduziert (nicht nur
dem Bericht geglaubt), F-3/F-4 im Diff gegengelesen. `make gates` läuft am
aktuellen `HEAD` grün (Exit 0, sechs Gates), `go test ./examples/http-client/...`
36/36 grün, `make fmt-check` und `make kommentar-kennungen
PATHS=examples/http-client COUNT=1` beide ohne Befund. Einziger eigener Fund
(V-1) ist eine nicht-blockierende Zahlenungenauigkeit in Commit-/Review-Text, kein
DoD- oder Entscheidungs-Verstoß.

Dieser Report ist ein **Lauf-Beleg** (dieser Stand, dieser Lauf); er ersetzt
weder Review noch Closure.
