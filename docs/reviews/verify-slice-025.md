# Verifier-Report: slice-025 — 2026-09-12

**Review-Art:** Verifikation — *wogegen*: Slice-Plan §2 (Definition of
Done, 9 Punkte), §3 (Plan-vs-Code inkl. beider Plan-Nachzug-Blöcke), §6
(Risiko-Ausgang, bleibt Planner-Entscheidung), §8 (Sub-Area-Prüfung), sowie
Entscheidungs-Konformität gegen [`SPEC-009`](../../spec/pflichtenheft.md)
(Metrik-Zeile `cdc_wal_retention_bytes`, Zeile 266) und die
Reviewer-Findings F-1/F-2 aus
[`review-slice-025.md`](review-slice-025.md). Nicht geprüft: Diff gegen
Plan/Hard Rules im Detail über die DoD-Punkte hinaus (Reviewer-Aufgabe,
bereits erledigt: 0 HIGH, 2 MEDIUM, 0 LOW/INFO, kein Merge-Blocker), realer
Bedarf (Validator — hier nicht einschlägig, kein MVP-Grenz-Slice).

**Grundsatz:** Keine Behauptung wurde übernommen — jeder Beleg unten wurde
in diesem Lauf selbst gelesen oder ausgeführt: `walretention.go` und
`stream_test.go` im Volltext, `git show`/`git log`, `make gates` und
`make test-replication` selbst gestartet, und eine eigene
Rot-Grün-Gegenprobe am `reconnectAfterError`-Aufruf durchgeführt (Schritt 3
unten) — nicht nur die Fixrunden-Behauptung übernommen, dass eine
Rot-Grün-Verifikation stattgefunden habe.

**Gegenstand:** `5d310fb` (feat: `WALRetentionChecker` + Verdrahtung +
`TestWALRetentionMeasuresGrowingBytes`), `3c79f8d` (docs: Benutzerhandbuch),
`55fa4d7` (docs/planning: DoD-Häkchen + Plan-Nachzug), `de5591b` (Review-
Report, 2 MEDIUM: F-1 fehlender Negativtest, F-2 kein Reconnect), `68134d5`
(fix: `reconnectAfterError` + vier neue reale Fehlerpfad-Tests), `c7db69d`
(docs/planning: Plan-Nachzug Fixrunde). Dazwischen `842d0b6` (Roadmap-
Housekeeping, `LH-FA-SST-007`) — unbeteiligt an diesem Slice, nicht
Gegenstand dieser Verifikation.

**Modell:** Claude Sonnet 5 · **Datum:** 2026-09-12

**Eingangs-Kontext:**

- Slice-Plan §1–§8 am aktuellen Stand
  (`in-progress/slice-025-metrik-wal-rueckstand.md`), inkl. beider
  Plan-Nachzug-Blöcke in §3
- `docs/reviews/review-slice-025.md` (0 HIGH, 2 MEDIUM: F-1, F-2; Verdikt
  „nicht merge-blockierend")
- `internal/adapters/driving/replication/receive/walretention.go`
  (Volltext, aktueller Stand nach `68134d5`)
- `internal/bootstrap/wiring.go` (Diff `68134d5`, Umfeld `runWALRetentionCheck`
  Zeilen 330–430)
- `internal/adapters/driving/replication/receive/stream_test.go` (Diff
  `68134d5`, alle vier neuen Tests im Volltext)
- `spec/pflichtenheft.md` Zeilen 224, 244, 266 (Volltext gelesen)
- `docs/user/benutzerhandbuch.md` (Diff `3c79f8d`, Volltext des neuen
  Abschnitts)
- `docs/plan/planning/observations/BEO-PGC/spec008-replication-luecke/state.md`,
  `.../dod-checkbox-nachzug/state.md`
- `git log`/`git show`/`git diff` über den vollen Commit-Verlauf
  (`b918e27..HEAD`)

---

## Sensor-Belege (in diesem Lauf selbst gefahren)

| Sensor | Ausgabe-Kernaussage | Exit |
|---|---|---|
| `make gates` | `baseline-verify: v6.5.0 OK — 54 Dateien` · `d-check` Standardlauf: 237 Dateien, 0 Befunde · `d-check --range HEAD~5..HEAD` (Modul `commits`): 237 Dateien, 0 Befunde · `commit-traceability.sh "HEAD~5..HEAD"`: OK, 5 Commits, Betreffs ohne Struktur-ID · `a-check`: 0 Befunde | **0** |
| `make test-replication` | alle Pakete `ok`, u. a. `.../receive` 4.613s, `.../bootstrap` 8.434s | **0** |
| `go test ./internal/.../receive/... -run TestWALRetentionMeasure -v` (eigener Docker-Lauf gegen frischen Testcontainer, DSN über eigenes Docker-Netz) | alle fünf `TestWALRetentionMeasure*`-Tests `PASS`, keiner geskippt | **0** |
| **Eigene Rot-Grün-Gegenprobe** (`reconnectAfterError(ctx)`-Aufruf in `Measure`s `IDENTIFY_SYSTEM`-Fehlerpfad auskommentiert, `TestWALRetentionMeasureReconnectsAfterConnectionLoss` erneut gegen denselben Testcontainer) | **FAIL**: `Measure nach Reconnect: … IDENTIFY_SYSTEM: conn closed` — Test scheitert sichtbar am dritten `Measure`-Aufruf, exakt am erwarteten Punkt | **1** (rot, wie erwartet) |
| Dieselbe Datei wiederhergestellt (`git diff` danach leer), Test erneut ausgeführt | alle fünf Tests wieder `PASS` | **0** |
| `git show --stat 68134d5` / `git show --stat c7db69d` | genau drei geänderte Dateien (`walretention.go`, `wiring.go`, `stream_test.go`) bzw. genau die Plan-Datei | — |

## DoD-Prüfung (Slice-Plan §2, Punkt für Punkt)

| # | DoD-Punkt | Verdikt | Beleg-Kernaussage |
|---|---|---|---|
| 1 | [`SPEC-009`](../../spec/pflichtenheft.md) erfüllt: `WALRetentionChecker` misst periodisch, real getestet, `make test-replication` grün, rot färbende Mutation einmal gesehen | **bestätigt** (Kernanspruch), Mutationsnachweis nicht in diesem Lauf wiederholt | Code gelesen: `Measure` bildet `int64(current.XLogPos - confirmed)` (`walretention.go:80`) — Richtung korrekt (wächst bei inaktivem Slot). `TestWALRetentionMeasuresGrowingBytes` lief in diesem Lauf real gegen PostgreSQL, `PASS`. Die im DoD-Text behauptete Rot-Grün-Mutation (LSN-Subtraktion vertauscht) stammt aus der *ursprünglichen* Implementierungsrunde (vor dem Review) — dieser Lauf hat sie nicht erneut reproduziert (außerhalb des beauftragten Prüfumfangs, der die *Fixrunden*-Rot-Grün-Probe verlangte, siehe Sensor-Tabelle); die Kernaussage „die Metrik ist real erfüllt" ist über eigenen Code-/Testlauf trotzdem eigenständig bestätigt |
| 2 | Messintervall an `heartbeatInterval` gebunden, keine neue Konfigurationsachse | **bestätigt** | `wiring.go:373`: `runWALRetentionCheck(walRetentionCtx, walRetention, log, heartbeatInterval)` — dieselbe Variable wie beim Heartbeat-Zug (`wiring.go:361`), kein neues Config-Feld im Diff |
| 3 | `make gates` grün | **bestätigt** | eigener Lauf in dieser Sitzung, alle vier inneren Gates 0 Befunde (siehe Sensor-Tabelle) |
| 4 | Review durchgeführt, Report unter `docs/reviews/` liegt vor | **materiell erfüllt, Checkbox nicht nachgezogen — siehe V-1** | `docs/reviews/review-slice-025.md` existiert real (`de5591b`), 0 HIGH/2 MEDIUM/0 LOW/0 INFO, „nicht merge-blockierend" — die Bedingung ist materiell erfüllt (Rollenwechsel fand statt, kein Self-Review), die Checkbox in §2 steht weiterhin auf `[ ]`, auch nach der Fixrunde (`68134d5`/`c7db69d` ändern sie nicht) |
| 5 | Doku-Update `docs/user/benutzerhandbuch.md` §4/§9/Änderungshistorie 1.3 | **bestätigt** | `git show 3c79f8d` zeigt genau den neuen Abschnitt „WAL-Rückstand prüfen" (nach „Metriken lesen"), den `max_wal_senders`-Hinweis unter „Grenzwerte" und die Änderungshistorie-Zeile 1.3 — Log-Beispiel im Dokument (`metric": "cdc_wal_retention_bytes", "bytes": 12345`) stimmt mit dem tatsächlichen Log-Aufruf überein (durch eigene Lektüre von `wiring.go` bestätigt, nicht nur Reviewer-Aussage übernommen) |
| 6 | Closure-Notiz mit Steering-Loop-Lerneintrag | **korrekt offen** | §7 trägt weiterhin vollständig den Platzhalter-Vorlagentext — Planner-Arbeit, noch nicht fällig vor `git mv` |
| 7 | Reconciliation-Register fortgeschrieben, falls einschlägig | **entfällt — korrekt geprüft** | `docs/plan/planning/reconciliation.md` existiert nicht; Repo durchgehend Greenfield (`harness/conventions.md` §Modus-Deklaration) |
| 8 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | `BEO-PGC/spec008-replication-luecke/state.md` trägt unverändert „weiter offen", 1× (`evidence/slice-020.md`) — Plan §8 benennt den Treffer bereits korrekt als „dieser Slice ist Teil der Antwort", ohne die Registerdatei selbst zu ändern (Planner-Arbeit bei Closure) |
| 9 (inkl. Risiken + Paarungen) | Jedes Risiko aus §6 trägt einen Ausgang; drei Paarungen getragen | **korrekt offen** | Beide §6-Risiken tragen weiterhin `<bei Closure einzutragen>` (siehe unten); Paarungen laufen laut DoD-Zeile bei dieser Closure selbst (Slice gehört zu `welle-7`, siehe Kopf — Paarungen verschieben sich auf die Welle-Closure, analog zu `slice-024`; die Checkbox-Formulierung im Plan lässt beide Fälle offen, siehe V-2) |

**Zwischenstand: 5/9 Kriterien materiell erfüllt und in diesem Lauf selbst
nachgeprüft (inkl. eigener Rot-Grün-Gegenprobe für F-2), 1 Item korrekt
entfallen (Reconciliation-Register, Greenfield), 3 Items regulär offen als
Planner-Closure-Arbeit (Closure-Notiz, Beobachtungs-Register,
Risiko-Ausgänge/Paarungen). Punkt 4 trägt einen eigenen Befund (V-1):
materiell erfüllt, Checkbox nicht nachgezogen.**

## Eigene Rot-Grün-Gegenprobe (F-2, Detailprotokoll)

1. `internal/adapters/driving/replication/receive/walretention.go:64`
   auskommentiert: `// VERIFY-DISABLED: c.reconnectAfterError(ctx)` im
   `IDENTIFY_SYSTEM`-Fehlerpfad von `Measure`.
2. `go test ./internal/.../receive/... -run TestWALRetentionMeasureReconnectsAfterConnectionLoss -v`
   gegen einen frisch gestarteten PostgreSQL-Testcontainer (eigenes
   Docker-Netz `cdc-repl-verify`, Image-Digest identisch zum Makefile-Pin):
   **FAIL** — `Measure nach Reconnect: … IDENTIFY_SYSTEM: conn closed`. Der
   Test scheitert exakt an der dritten `Measure`-Erwartung (Wiederherstellung
   nach Verbindungsabbruch), nicht an der zweiten (sichtbares Scheitern nach
   Abbruch bleibt auch ohne Reconnect korrekt, wie erwartet).
3. Datei mit der ursprünglichen Fassung überschrieben (`git diff` danach
   leer — keine Restspur im Arbeitsverzeichnis).
4. Derselbe Testlauf erneut: alle fünf `TestWALRetentionMeasure*`-Tests
   `PASS`.

**Ergebnis: Der Test beweist echte Recovery, nicht nur „kein Absturz".**
Ohne den Reconnect-Pfad schlägt exakt der Aufruf fehl, der die
Wiederherstellung behauptet — die Gegenprobe bestätigt, dass der Test die
Behauptung tatsächlich trägt.

**Backend-Identifikation geprüft:** Der Test markiert die
Checker-Verbindung mit einem eindeutigen `application_name`
(`dsnWithApplicationName`, `stream_test.go:576-589`) und zielt
`pg_terminate_backend` ausschließlich auf Backends mit diesem Namen
(`terminateBackend`, `stream_test.go:604-618`). Eigene Prüfung:
`env.pool.Config().ConnString()` liefert die ursprüngliche `postgres://…`-DSN
unverändert zurück (pgx `ConnConfig.ConnString()` gibt die beim Parsen
übergebene Zeichenkette zurück, kein synthetisiertes Format) — die
URL-Erkennung in `dsnWithApplicationName` trifft daher zu, `&application_name=…`
wird korrekt anghängt. Der Stream (`env.pool`, INSERT/Abfragen) und der
Checker verbinden sich mit unterschiedlichen `application_name`-Werten
(Stream: kein expliziter Wert; Checker: `walretention_reconnect_test`) —
`pg_terminate_backend` trifft nachweislich **nicht** die Stream-Verbindung
(die zum Testzeitpunkt ohnehin bereits über `cancelRun()`/`stream.Run`-Ende
beendet ist, siehe `stream_test.go:648-660`, bevor der Checker überhaupt
aufgebaut wird).

## Die drei anderen neuen Fehlerpfad-Tests (F-1)

| Test | Prüft real? | Befund |
|---|---|---|
| `TestWALRetentionMeasureInvalidSlotName` | ja | Slot-Name außerhalb `identifierShape` (`"Ungueltiger Slot!"`) — `NewWALRetentionChecker` weist ihn vor jeder DB-Interaktion mit `ErrConfiguration` ab; kein Testcontainer nötig, echter Regex-Ablehnungspfad, kein Zirkelschluss |
| `TestWALRetentionMeasureConnectionRefusedFails` | ja | Verbindung gegen `127.0.0.1:1` (kein Listener) mit `connect_timeout=2` — `ErrReplication` real durch TCP-Verbindungsfehler ausgelöst, nicht simuliert |
| `TestWALRetentionMeasureMissingSlot` | ja | Checker gegen real existierende PostgreSQL, aber nie angelegten Slot-Namen (`slot_pgc_test_missing_retention`) — `IDENTIFY_SYSTEM` gelingt (Slot-unabhängig), die Katalogabfrage liefert `exists=false`, `Measure` meldet `ErrReplication` über den `!exists`-Zweig — trägt exakt den in F-1 benannten Fehlerpfad, kein triviales Assert |

Alle drei sind reale, nicht-zirkuläre Tests gegen tatsächliche Fehlerzustände
(Konfiguration, Transport, fehlender Katalogeintrag) — keine Mock-Stubs, kein
Test, der nur die eigene Testfixture zurückspiegelt.

## `reconnectAfterError` — Verdrahtung und Dauerhänger-Frage

`Measure` ruft `c.reconnectAfterError(ctx)` in **beiden** Fehlerpfaden auf,
die eine defekte Verbindung anzeigen können (`IDENTIFY_SYSTEM`-Fehler,
Katalogabfrage-Fehler; `walretention.go:63-72`) — **vor** der Rückgabe des
Fehlers, sodass der *nächste* `Measure`-Aufruf (nächster Tick in
`runWALRetentionCheck`) bereits die neue Verbindung nutzt. Der
`!exists`-Zweig und der `ParseLSN`-Fehlerzweig lösen **keinen** Reconnect
aus — korrekt so: Beide sind Katalog-/Daten-Zustände (fehlender Slot,
kaputter LSN-Text), keine Verbindungsstörung; ein Reconnect würde daran
nichts ändern.

**Dauerhänger-Frage beantwortet:** Schlägt `reconnectAfterError` selbst
fehl (`connectReplication` liefert einen Fehler), bleibt `c.conn` auf der
bereits geschlossenen alten Verbindung stehen (`conn`-Variable wird nur bei
`err == nil` übernommen, `walretention.go:96-98`). Der **übernächste**
`Measure`-Aufruf versucht erneut `IdentifySystem` auf dieser toten
Verbindung, scheitert wieder, ruft erneut `reconnectAfterError` auf — **kein
dauerhaftes Hängenbleiben**, sondern ein Retry pro Tick, begrenzt durch den
Heartbeat-Takt. Es gibt keinen Pfad, auf dem der Checker nach einem
einmaligen Fehler dauerhaft aufgibt; jeder Tick ist ein neuer Versuch. Kein
Ressourcen-Leck dabei: Die alte Verbindung wird vor jedem Neuaufbau-Versuch
geschlossen (`_ = c.conn.Close(ctx)`, Zeile 95), ein fehlgeschlagener
`connectReplication`-Versuch hält kein offenes Handle zurück.

## Race-Condition- und Ressourcen-Prüfung (eigene Suche nach einem dritten Problem)

- geprüft, ohne Befund: **Keine parallele Nutzung von `checker`** —
  `runWALRetentionCheck` läuft in genau einer Goroutine
  (`wiring.go:371-374`), `checker.Measure`/`checker.Close` werden nirgends
  sonst aufgerufen; `Close` läuft über `defer` erst nach
  `walRetentionDone.Wait()` (`wiring.go:344-348,380-381`) — kein `Close`
  gegen eine noch laufende `Measure`-Anfrage.
- geprüft, ohne Befund: **Kein Verbindungsleck im Fehlerfall** — siehe
  oben, jeder `reconnectAfterError`-Aufruf schließt die Alt-Verbindung vor
  dem Neuaufbau-Versuch, ein gescheiterter Versuch hinterlässt kein
  offenes Handle.
- geprüft, ohne Befund: **Kein zusätzlicher `IDENTIFY_SYSTEM`-Aufruf pro
  Tick durch den Reconnect** — nach einem Reconnect beginnt der *nächste*
  Tick wieder regulär bei `IdentifySystem`; es gibt keinen doppelten
  Aufruf innerhalb eines einzelnen `Measure`-Aufrufs.
- geprüft, ohne Befund: **Kontext-Weitergabe beim Reconnect** —
  `reconnectAfterError(ctx)` nutzt denselben `ctx`, den `Measure` erhalten
  hat; bei Prozess-Shutdown (ctx bereits storniert) scheitert der
  Neuverbindungsversuch schnell und wird verworfen — kein hängender
  Reconnect-Versuch nach Shutdown-Signal.
- Kein drittes, bisher unbemerktes Problem gefunden, das über die bereits
  von Reviewer (F-1, F-2) und dieser Verifikation (V-1) benannten Befunde
  hinausgeht.

## Plan-vs-Code-Diff (Plan-Nachzug Fixrunde, §3, gegen `68134d5`/`c7db69d`)

| Plan-Nachzug-Zeile | Behauptung | Abgleich gegen `git show` |
|---|---|---|
| `walretention.go` | `dsn`-Feld + `reconnectAfterError`-Ersatzpfad bei Protokoll-/Katalogfehler | **stimmt** — `git show 68134d5 -- walretention.go` zeigt genau `dsn`-Feld, zwei neue `reconnectAfterError`-Aufrufe, neue Methode |
| `wiring.go` | Nur Kommentar-Nachzug an `runWALRetentionCheck`, keine Verhaltensänderung der Schleife | **stimmt** — `git show 68134d5 -- wiring.go` zeigt ausschließlich einen Kommentar-Zusatz (4 Zeilen), keine Code-Zeile geändert |
| `stream_test.go` | drei neue Fehlerpfad-Tests + `TestWALRetentionMeasureReconnectsAfterConnectionLoss` | **stimmt** — `git show 68134d5 -- stream_test.go` zeigt genau die vier benannten Testfunktionen plus die Hilfsfunktionen `dsnWithApplicationName`/`countBackends`/`terminateBackend` |

`git show --stat 68134d5` zeigt genau diese drei Dateien, keine weitere —
Abweichungsfreiheit auch auf Datei-Ebene bestätigt. `c7db69d` ändert
ausschließlich die Slice-Plan-Datei (Plan-Nachzug-Tabelle selbst) — konsistent
mit dem Nachzugs-Charakter.

## §6 Risiken — Ausgang bei diesem Verifikationsstand

Beide im Slice-Plan §6 genannten Risiken tragen **weiterhin den
unausgefüllten Platzhalter** `<bei Closure einzutragen>`:

1. Künstlich erzeugter WAL-Rückstand im Testcontainer könnte flaky sein.
2. Der Expositionsweg war zum Planungszeitpunkt noch nicht entschieden.

Beide sind inhaltlich bereits **materiell** beantwortet (nicht verwechseln
mit dem formalen §6-Ausgang, der Planner-Pflicht bleibt): Risiko 1 hat sich
über mehrere reale Testläufe in dieser Sitzung (inkl. der eigenen
Rot-Grün-Gegenprobe) nicht als flaky gezeigt — plausibler Ausgang
„entfallen" oder „weiter offen mit Beobachtung", Planner-Urteil. Risiko 2
ist durch den Plan-Nachzug (§3, „Entscheidung nachgetragen: … Log, nicht …
Tabellenspalte") bereits entschieden — plausibler Ausgang „entfallen,
Begründung: Entscheidung getroffen (strukturiertes Log)".

## Eigene Befunde

### V-1 — DoD-Checkbox „Review durchgeführt" nicht nachgezogen, obwohl die Bedingung materiell erfüllt ist

- `kategorie`: LOW
- `pfad`: `docs/plan/planning/in-progress/slice-025-metrik-wal-rueckstand.md:108-110`
  (Checkbox weiterhin `[ ]`) vs. `docs/reviews/review-slice-025.md`
  (Report existiert, 0 HIGH/2 MEDIUM, „nicht merge-blockierend")
- `befund`: Derselbe Musterbefund wie `verify-slice-024.md` V-1: Die
  DoD-Zeile „Review durchgeführt, Report unter `docs/reviews/` liegt vor …
  kein Self-Review" ist materiell erfüllt (Rollenwechsel fand statt, Report
  liegt vor), die Checkbox wurde weder im Review-Commit (`de5591b`) noch in
  der Fixrunde (`68134d5`/`c7db69d`) auf `[x]` gesetzt. Die zugrunde
  liegende Regel-Klasse ist in `BEO-PGC/dod-checkbox-nachzug` bereits als
  **verkörpert** geführt (`.claude/commands/implement-slice.md`, Schritt 18/21,
  seit `welle-5`) — hier betrifft es allerdings nicht den engeren
  Schritt-21-Fall (die Fixrunde löst keinen *bestehenden* DoD-Punkt auf,
  sondern behebt zwei MEDIUM-Findings, die keine eigene DoD-Zeile tragen),
  sondern denselben Fall wie bei `slice-024`: eine materiell erfüllte
  Review-Bedingung, deren Checkbox nicht mitgezogen wurde.
- `verifizierbar`: ja — `grep -n '^\- \[' slice-025-….md` zeigt die
  Checkbox weiterhin `[ ]`; kein Commit nach `de5591b` ändert diese Zeile.
- **Für die Closure:** kein Blocker — der Planner sollte die Checkbox vor
  `git mv` nach `done/` auf `[x]` setzen (materiell gedeckt).

### V-2 — Drei-Paarungen-Checkbox lässt bei einem wellen-zugehörigen Slice beide Lesarten offen

- `kategorie`: LOW
- `pfad`: `docs/plan/planning/in-progress/slice-025-metrik-wal-rueckstand.md:119`
- `befund`: Wie bereits in `verify-slice-024.md` (§Closure-Bereitschaft
  Punkt 4) angemerkt: `slice-025` gehört zu `welle-7` (Kopf-Feld
  „Welle:"), die drei Paarungen laufen bei einem wellen-zugehörigen Slice
  regulär erst mit der Welle-Closure, nicht bei der einzelnen
  Slice-Closure (Modul 6 §Wann Arbeit eine Welle braucht). Die DoD-Zeile
  benennt zwar beide Fälle nebeneinander, sagt aber nicht ausdrücklich, wer
  für *diesen* Slice zuständig ist — dieselbe Unschärfe wie bei
  `slice-024` unbehoben in die Fixrunde durchgereicht.
- **Für die Closure:** kein Blocker, aber der Planner sollte beim
  Checkbox-Nachzug vermerken, dass diese Pflicht auf die `welle-7`-Closure
  verschoben ist, statt die Zeile unkommentiert `[ ]` zu lassen.

## Negativbefunde

- geprüft, ohne Befund: **`make gates`** — eigener Lauf, 0 Befunde in
  allen vier Gates.
- geprüft, ohne Befund: **`make test-replication`** — eigener Lauf, alle
  Pakete grün, inkl. `.../receive` und `.../bootstrap`.
- geprüft, ohne Befund: **Eigene Rot-Grün-Gegenprobe F-2** — Test scheitert
  sichtbar ohne Reconnect-Pfad, besteht mit ihm; kein Restschaden am
  Arbeitsverzeichnis nach Wiederherstellung (`git status` sauber).
- geprüft, ohne Befund: **Backend-Identifikation im Reconnect-Test** —
  eindeutiger `application_name`, `pg_terminate_backend` trifft nachweislich
  nicht die (zu diesem Zeitpunkt bereits beendete) Stream-Verbindung.
- geprüft, ohne Befund: **Drei weitere Fehlerpfad-Tests (F-1)** — real,
  nicht zirkulär, decken Konfigurations-, Transport- und
  Katalog-Fehlerpfade tatsächlich ab.
- geprüft, ohne Befund: **Kein dauerhaftes Hängenbleiben nach
  gescheitertem Reconnect** — jeder Tick ist ein neuer Versuch, kein
  Ressourcen-Leck.
- geprüft, ohne Befund: **Keine Race Condition** — `checker` hat genau
  einen Aufrufer (die `runWALRetentionCheck`-Goroutine), `Close` läuft erst
  nach `Wait()`.
- geprüft, ohne Befund: **Plan-Nachzug Fixrunde gegen realen Diff** — alle
  drei genannten Dateien stimmen exakt mit `git show --stat 68134d5`
  überein, keine unangekündigte vierte Datei.
- geprüft, ohne Befund: **Doku-Update-Inhalt** — Log-Beispiel im
  Benutzerhandbuch stimmt wörtlich mit dem tatsächlichen `log.Info`-Aufruf
  überein (eigene Lektüre beider Stellen).
- geprüft, ohne Befund: **Reconciliation-Register** — Datei existiert nicht,
  Repo durchgehend Greenfield, DoD-Punkt korrekt als entfallend behandelt.
- geprüft, ohne Befund: **Beobachtungs-Register unverändert wie
  angekündigt** — `BEO-PGC/spec008-replication-luecke/state.md` trägt
  weiterhin „weiter offen", 1×; keine unangekündigte Änderung durch diesen
  Slice.
- geprüft, ohne Befund: **Traceability** — alle sechs inhaltstragenden
  Commit-Betreffs nennen `ADR-0049`, kein `SPEC-*`/`ARC-*` im Betreff;
  `make commit-traceability` (Teil von `make gates`) lief grün.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 2 (V-1, V-2) |
| INFO | 0 |

**Zusammenfassung DoD:** 5/9 Kriterien materiell erfüllt und in diesem Lauf
selbst geprüft (inkl. eigener Rot-Grün-Gegenprobe für F-2 und Einzelprüfung
der drei F-1-Tests), 1 Item korrekt entfallen (Reconciliation-Register,
Greenfield), 3 Items regulär offen als Planner-Closure-Arbeit
(Closure-Notiz, Beobachtungs-Register, §6-Risiko-Ausgänge/Paarungen). Kein
DoD-Defekt im Sinn eines unbelegten „bestätigt"-Punkts — beide eigenen
Befunde (V-1, V-2, je LOW) sind Checkbox-Diskrepanzen, die die materielle
Substanz nicht widerlegen.

## Verdikt

**DoD-/Entscheidungs-Konformität: bestätigt, mit zwei offenen
Klein-Befunden (2× LOW, V-1/V-2; non-blocking für diesen Slice, empfohlene
Korrektur vor `git mv` nach `done/`).** Die Fixrunde hat beide
Reviewer-Findings tatsächlich behoben, nicht nur behauptet:
`reconnectAfterError` ist korrekt verdrahtet (beide Verbindungsfehlerpfade,
kein Dauerhänger, kein Ressourcen-Leck), und
`TestWALRetentionMeasureReconnectsAfterConnectionLoss` beweist echte
Recovery — durch eigene Rot-Grün-Gegenprobe in diesem Lauf verifiziert,
nicht aus der Fixrunden-Behauptung übernommen. Die drei übrigen neuen
Tests (F-1) sind reale, nicht-zirkuläre Fehlerpfad-Tests. Der Plan-Nachzug
der Fixrunde stimmt exakt mit dem tatsächlichen Diff überein. `make gates`
und `make test-replication` liefen in diesem Lauf selbst grün.

**Plan-vs-Code-Diff:** vollständige Deckung — keine unangekündigte Datei,
Out-of-Scope-Punkte aus §1 gewahrt (keine `cdc.metrics`-Erweiterung, keine
Änderung der bestehenden `START_REPLICATION`-Lookup-Nutzung, kein
Schwellen-Vergleich implementiert).

**Closure-Bereitschaft: noch nicht — drei reguläre Planner-Schritte stehen
aus**, keiner davon ein Defekt an bereits gelieferter Substanz:

1. Closure-Notiz §7 schreiben (Was hat funktioniert / anders als geplant /
   Steering-Loop-Eintrag / Beobachtungs-Register / Folge-Slices / §6-Risiken
   / Drei Paarungen).
2. Beide §6-Risiken auf einen der drei Ausgänge setzen — plausible Ausgänge
   oben skizziert, formales Setzen bleibt Planner-Urteil.
3. Beobachtungs-Register-Entscheidung: löst dieser Slice
   `BEO-PGC/spec008-replication-luecke` weiter nicht auf (Umsetzung folgt
   erst mit `slice-026`) — voraussichtlich unverändert „weiter offen".

Zusätzlich empfohlen: V-1 (Review-Checkbox) und V-2 (Drei-Paarungen-Checkbox-
Lesart) vor `git mv` nach `done/` beheben bzw. vermerken.

**Übergabe:** Bericht an den Planner. Keine Reparaturen — Plan-Datei und
Code wurden von diesem Lauf nicht dauerhaft verändert (die Rot-Grün-Gegenprobe
wurde vollständig rückgängig gemacht, `git status` am Ende sauber).

---

**Gate-Beleg:** `make gates` einmal in diesem Lauf ausgeführt, Exit 0, 0
Befunde (`baseline-verify`: 54 Dateien OK; `d-check` Standardlauf: 237
Dateien, 0 Befunde; `d-check --range HEAD~5..HEAD` Modul `commits`: 237
Dateien, 0 Befunde; `commit-traceability.sh`: OK, 5 Commits; `a-check`: 0
Befunde). `make test-replication` einmal vollständig ausgeführt, alle
Pakete grün. Zusätzlich ein gezielter Docker-Lauf gegen einen eigens
gestarteten Testcontainer für die Rot-Grün-Gegenprobe (Container und
Netzwerk danach abgeräumt). `git status` am Ende dieses Laufs sauber (keine
Arbeitsverzeichnis-Änderung durch die Verifikation selbst).
