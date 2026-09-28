# Verifikations-Report: examples/csharp/http-client Verb-Matrix (Fixrunde) — 2026-09-28

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-/Entscheidungs-Konformität +
Plan-vs-Code-Diff + Gates gegen den Review-Report des Reviewers:
[`review-example-csharp-http-client-verbmatrix.md`](review-example-csharp-http-client-verbmatrix.md).
Kein eigener Slice-Plan für diesen Zug (analog zum Go-Vorbild); dieser Lauf
prüft ausschließlich die Auflösung des einen Review-Findings gegen den
Ist-Zustand sowie die vier vom Reviewer behaupteten „proaktiv vermiedenen“
Go-Findings, mit eigener, unabhängiger Gegenprobe. Formvorbild:
[`verifikation-example-http-client-verbmatrix.md`](verifikation-example-http-client-verbmatrix.md)
(Go-Geschwister-Lauf).

**Gegenstand:** `git diff 79c92471..d1b4c5f2 -- examples/csharp/http-client
examples/README.md` (kompletter Verlauf: Feature-Commit `9336350a`, Review
`9602405c` mit 0 HIGH + 1 MEDIUM, Fixrunde `d1b4c5f2`). Dieser Lauf ändert
weder Code noch Doku am Gegenstand; er schreibt nur diesen Report. `HEAD` bei
Beginn dieses Laufs: `d1b4c5f2` (Arbeitsbaum `git status` clean, drei
unpushte Commits — unverändert bei Abschluss).

**Eingangs-Kontext:**

- [`review-example-csharp-http-client-verbmatrix.md`](review-example-csharp-http-client-verbmatrix.md)
  (das eine Finding F-1, MEDIUM)
- `internal/application/port/inbound/readchanges.go:14-29` (`ReadChangesQuery`-
  Vertrag: Start inklusiv, End exklusiv)
- `examples/csharp/http-client/{Program.cs,Cli.cs,ChangesClient.cs,Dispatcher.cs}`
  (Volltext gelesen)
- `LH-FA-SST-006`
- `AGENTS.md` §3.1 (Docker-only), §3.7 (Kommentarform), §3.9
  (Exit-Code-Disziplin), §3.12 (Herkunft von Aussagen)

---

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — [`AGENTS.md`](../../AGENTS.md) §3.9)

| Sensor/Probe | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `git show d1b4c5f2 --stat` | — | genau **eine** Datei: `examples/csharp/http-client/Program.cs`, 2 Zeilen +/- |
| `Program.cs:3-15` gelesen | — | Kopfkommentar trägt jetzt genau **eine** Kennung (`LH-FA-SST-006`); der `ADR-0087`-Satz ist entfernt |
| `internal/application/port/inbound/readchanges.go:17` gelesen | — | „Start **inklusiv**, End **exklusiv**“ — Referenz-Vertrag für F-1-Vermeidung |
| `Cli.cs`/`Program.cs` nach Hilfetext für `--from`/`--to` durchsucht | kein Treffer | keine Flag-Hilfe-Ausgabe existiert in diesem manuell geschriebenen Parser — es gibt keine Stelle, die eine falsche Semantik behaupten könnte |
| `ChangesClient.cs:36-46` (`ChangesUrlBuilder`) gelesen | — | Kommentar beschreibt korrekt: leerer Wert eines bekannten optionalen Parameters bleibt unberücksichtigt, endet `200` — deckungsgleich mit Server-Verhalten (F-2-Vermeidung, Go-Analogie) |
| `grep -rn "FailsOnNon2xx" HttpClient.Tests/*.cs` | 9 Treffer | je einer für `ReadChangesAsync`, `RegisterConsumerAsync`, `AcknowledgeConsumerAsync`, `ConsumerPositionAsync`, `RemoveConsumerAsync`, `EnableTableAsync`, `DisableTableAsync`, `TableStatusAsync`, `RunRetentionAsync` — alle neun, keine fehlt |
| `grep -rn "System.Net.Http.HttpClient httpClient\|new System.Net.Http.HttpClient" examples/csharp/http-client/*.cs` (ohne Tests) | — | nur `Program.cs:40` konstruiert einen `HttpClient`; jeder andere Client (inkl. `TablesClient`) nimmt ihn als Parameter — F-4-Vermeidung bestätigt |
| `docker build --build-context proto=proto --target build --no-cache -t verify-csharp-build examples/csharp` | Exit 0 | `HttpClient.Tests`: „Passed! - Failed: 0, Passed: 49, Skipped: 0, Total: 49“ — frischer, nicht gecachter Lauf gegen den `d1b4c5f2`-Stand |
| `make gates` | Exit 0 | `baseline-verify: v6.9.0 OK — 54 Dateien` · `d-check: 1373 Datei(en) geprüft, 0 Befund(e)` · `a-check: gesamt: 0 Befund(e)` · `commit-traceability: OK — 5 Commit(s) in „HEAD~5..HEAD“` · `coverage-gate: OK — Coverage 85.30% erfüllt Schwelle 80%` · `generated-sync: OK` |
| `make kommentar-kennungen PATHS=examples/csharp/http-client COUNT=1` | Exit 0 | `0` (bestätigt nur die bekannte Werkzeug-Grenze — `.cs` liegt außerhalb des Suchraums, siehe Review-Report F-1) |
| `make example-demo-up` (real, eigener Rundlauf) | Exit 0 | Umgebung hochgefahren; Demo-Zeile über `GET /changes?source=demo-source` bereit |
| `make example-run-csharp SURFACE=http ARGS="--verb=register-consumer --consumer-id verify-vf --name verify-vf-name --source demo-source"` | Exit 0 | `{"consumer_id":"verify-vf","name":"verify-vf-name","already_registered":false}` |
| `make example-run-csharp SURFACE=http ARGS="--verb=consumer-position --consumer-id verify-vf"` | Exit 0 | `{"consumer_id":"verify-vf","source_id":"","offset":0,"acknowledged":false}` |
| eigene DB-Gegenprobe (`psql` gegen `cdc.consumer`) | — | `verify-vf \| verify-vf-name` — real in der Datenbank angelegt, unabhängig vom Client-Output |
| `make example-demo-down` | Exit 0 | Umgebung sauber abgeräumt (Container + Netz entfernt) |
| `git status` nach allen Läufen | clean | Arbeitsbaum unverändert; dieser Lauf hat keinen Code/Doku-Diff erzeugt |

Jeder Lauf wurde mit direkt geprüftem Exit-Code ausgeführt, nie durch eine
Pipe/einen Wrapper hindurch.

## 2. Das eine Review-Finding (F-1, MEDIUM) — Auflösung am Ist-Zustand

Der Review-Report fand einen Kopfkommentar in `Program.cs` mit zwei Kennungen
(`LH-FA-SST-006`, `ADR-0087`) im selben Block, verstößt gegen `AGENTS.md` §3.7
(höchstens eine Kennung je Kommentarblock). Eigene Lektüre von `Program.cs:3-15`
am `d1b4c5f2`-Stand: der Block trägt jetzt genau eine Kennung (`LH-FA-SST-006`);
der Satz „Startform ist ein Container-Aufruf, kein Host-Aufruf (`ADR-0087`
Festlegung 3)“ ist auf „Startform ist ein Container-Aufruf, kein Host-Aufruf“
gekürzt — die Aussage bleibt (reiner Ist-Zustand ohne Kennung), nur der
zweite Anker ist entfernt. `git show d1b4c5f2` bestätigt: die Änderung ist
minimal (2 Zeilen ersetzt) und betrifft ausschließlich diesen Kommentarblock,
keine funktionale Änderung. **F-1 ist am Ist-Zustand aufgelöst**, unabhängig
von der Fixrunden-Commit-Message nachgewiesen (eigene Quelltext-Lektüre).

## 3. Die vier „proaktiv vermiedenen“ Go-Findings — eigene Gegenprobe (nicht nur dem Review-Bericht geglaubt)

**F-1-Vermeidung (from/to-Semantik).** Das Go-Vorbild hatte eine Flag-Hilfe
mit genau umgekehrter Positionsgrenzen-Semantik. Die C#-Implementierung
(`Cli.cs`) hat **keine** Flag-Hilfe-Ausgabe für `--from`/`--to` (kein
`--help`/Usage-Text im gesamten Parser, eigener `grep` ohne Treffer) — es
gibt strukturell keine Stelle, die eine falsche Aussage treffen könnte. Der
Server-Vertrag (`readchanges.go:17`: Start inklusiv, End exklusiv) bleibt
unverändert die Referenz; die C#-Implementierung reicht `--from`/`--to`
unkommentiert an den Server durch (`ChangesUrlBuilder.Build`). **Bestätigt,
kein Fund.**

**F-2-Vermeidung (leerer-Parameter-Kommentar).** Das Go-Vorbild behauptete
fälschlich einen `400` für jeden leeren bekannten Parameter. Eigene Lektüre
von `ChangesUrlBuilder`s Kommentar (`ChangesClient.cs:36-46`): „ein leerer
Wert eines bekannten optionalen Parameters … bleibt für den Server
unberücksichtigt und endet mit `200`“ — korrekt gemäß Server-Verhalten
(`parseReadChangesQuery`/`readChangesPosition`/`readChangesLimit`); der
Client selbst sendet nie einen leeren Wert (`AddIfNotEmpty`). **Bestätigt,
kein Fund.**

**F-3-Vermeidung (fehlende Nicht-2xx-Tests).** Eigenes Zählen (nicht aus dem
Review-Bericht übernommen): `grep -rn "FailsOnNon2xx"
HttpClient.Tests/*.cs` liefert genau neun Treffer, einen je der neun neuen
Aufruf-Funktionen (`ReadChangesAsync`, `RegisterConsumerAsync`,
`AcknowledgeConsumerAsync`, `ConsumerPositionAsync`, `RemoveConsumerAsync`,
`EnableTableAsync`, `DisableTableAsync`, `TableStatusAsync`,
`RunRetentionAsync`). Keine der neun fehlt — anders als beim Go-Vorbild
(dort fehlten fünf von neun). **Bestätigt, kein Fund.**

**F-4-Vermeidung (eigener HttpClient).** Eigener `grep` über alle
Nicht-Test-Dateien nach `System.Net.Http.HttpClient httpClient`/
`new System.Net.Http.HttpClient`: nur `Program.cs:40` konstruiert einen
`HttpClient`; jede Aufruf-Funktion (inklusive `TablesClient.ListTablesAsync`)
nimmt ihn als Parameter entgegen. Kein zweiter, unabhängiger Client wie beim
Go-Vorbild (`listTables(cfg)` baute dort einen eigenen). **Bestätigt, kein
Fund.**

## 4. Reale E2E-Gegenprobe (eigener, vom Bericht unabhängiger Verb)

Reviewer und Implementer belegten bereits `table-status`/`enable-table`/
`changes`/`disable-table` sowie die `from`/`to`-Grenzsemantik. Für diesen
Verifikationslauf wurde ein anderer Verb-Pfad gewählt:
`register-consumer` gefolgt von `consumer-position`, gegen die frisch
hochgefahrene Demo-Umgebung (`make example-demo-up`). Beide Aufrufe liefen
real gegen den laufenden Feed-Container über das Docker-Netz
`cdc-examples`, kein `docker exec`, kein Mock. Die Registrierung wurde
zusätzlich unabhängig gegen `cdc.consumer` per `psql` bestätigt (`verify-vf`
/ `verify-vf-name`, real in der Datenbank sichtbar) — nicht nur über den
Client-Output geglaubt. Die Umgebung wurde danach über
`make example-demo-down` sauber abgeräumt.

## 5. Arbeitsbaum-Prüfung

`git show d1b4c5f2 --stat`: genau eine Datei,
`examples/csharp/http-client/Program.cs` (2 Zeilen geändert). Kein
Mitziehen anderer Dateien. `git status` vor und nach diesem Lauf: clean.
Der Arbeitsbaum ist drei Commits vor `origin/main` (unverändert seit
Aufgabenbeginn) — dieser Lauf hat weder committet noch gepusht, wie es
seine Rolle vorsieht (Verifier berichtet, ändert nicht).

## 6. Verifier-Findings (eigene, unabhängig vom Review)

Kein eigener Fund. Alle geprüften Flächen (Kopfkommentar-Form, F-1…F-4-
Vermeidung, Testabdeckung, Docker-only, Arbeitsbaum-Hygiene, Gates, reale
E2E-Probe) decken sich mit den Berichten von Implementer und Reviewer, real
und unabhängig nachgemessen.

## 7. Verdikt

**Bereit für Abschluss** — keine weitere Fixrunde nötig. Das eine
Review-Finding (F-1, MEDIUM, Kommentar-Kette in `Program.cs`) ist am
Ist-Zustand real und unabhängig verifiziert aufgelöst: der Kopfkommentar
trägt jetzt genau eine Kennung, `git show d1b4c5f2` bestätigt eine minimale,
auf diesen Block beschränkte Änderung. Die vier vom Reviewer als „proaktiv
vermieden“ deklarierten Go-Findings (F-1…F-4 des Vorbild-Reviews) wurden in
diesem Lauf **eigenständig gegengeprüft**, nicht nur dem Review-Bericht
geglaubt — für F-1/F-2 durch Quelltext-/Kommentarlektüre gegen den
Server-Vertrag, für F-3 durch eigenes Nachzählen (9/9 Nicht-2xx-Tests
vorhanden), für F-4 durch eigenen `grep` (ein einziger `HttpClient`-
Konstruktionspunkt). `make gates` läuft am aktuellen Stand grün (Exit 0,
sechs Gates), ein frischer `--no-cache`-Docker-Build bestätigt 49/49
`HttpClient.Tests` grün, und eine eigene, vom Bericht unabhängige
E2E-Probe (`register-consumer`/`consumer-position`) lief real gegen die
Demo-Umgebung, mit Gegenprobe direkt in der Datenbank. Kein offenes HIGH
oder MEDIUM verbleibt.

Dieser Report ist ein **Lauf-Beleg** (dieser Stand, dieser Lauf); er ersetzt
weder Review noch Closure.
