# Verifier-Report: slice-026 — 2026-09-12

**Review-Art:** Verifikation — *wogegen*: Slice-Plan §2 (Definition of
Done, 9 Punkte), §3 (Plan-vs-Code inkl. des Plan-Nachzug-Blocks der
Fixrunde), §6 (Risiko-Ausgang, bleibt Planner-Entscheidung), §8
(Sub-Area-Prüfung), sowie Entscheidungs-Konformität gegen
[`ADR-0049`](../plan/adr/0049-replication-fehlerklassen-schwellen.md) und
die Reviewer-Findings aus
[`review-slice-026.md`](review-slice-026.md) (F-1, MEDIUM). Nicht geprüft:
Diff gegen Plan/Hard Rules im Detail über die DoD-Punkte hinaus
(Reviewer-Aufgabe, bereits erledigt: 0 HIGH, 1 MEDIUM, 0 LOW/INFO, kein
Merge-Blocker), realer Bedarf (Validator — hier nicht einschlägig, kein
MVP-Grenz-Slice).

**Grundsatz:** Keine Behauptung wurde übernommen — jeder Beleg unten wurde
in diesem Lauf selbst gelesen oder ausgeführt: `wiring.go`,
`walretention.go`, `receive.go` (relevante Ausschnitte) und beide neuen
Testdateien im Volltext, `git show`/`git log`, `make gates` selbst
gestartet, `make test-replication` **dreimal** selbst gestartet (nicht nur
die Reviewer-Behauptung übernommen), `go test -race -count=5` auf den
Whitebox-Tests selbst gefahren, und **zwei eigene Rot-Grün-Gegenproben**
durchgeführt (Prioritäts-Garantie `mergeStreamAndWALFaultOutcome` UND
Zero-Value-Fallback `resolveWALRetentionThresholds`) — beide mit
anschließender Wiederherstellung, `git status` am Ende sauber.

**Gegenstand:** `e5e29af` (feat: WAL-Rückstand-Schwellen mit kontrollierter
Fortsetzung), `d4382b5` (docs/planning: DoD-Häkchen-Nachzug), `29a49ea`
(fix: Test für Zero-Value-Fallback, Review-Finding F-1 behoben),
`3c33b1e` (Review-Report, 0 HIGH/1 MEDIUM).

**Modell:** Claude Sonnet 5 · **Datum:** 2026-09-12

**Eingangs-Kontext:**

- Slice-Plan §1–§8 am aktuellen Stand
  (`in-progress/slice-026-schwellen-ueberwachung-kontrollierte-fortsetzung.md`),
  inkl. des Plan-Nachzug-Blocks in §3
- `docs/reviews/review-slice-026.md` (0 HIGH, 1 MEDIUM F-1; Verdikt „nicht
  merge-blockierend")
- `internal/bootstrap/wiring.go` (Volltext, 761 Zeilen, aktueller Stand
  nach `29a49ea`)
- `internal/adapters/driving/replication/receive/receive.go`
  (`Stream.Run`/`process`, Zeilen 286–380) und
  `internal/adapters/driving/replication/receive/walretention.go`
  (Volltext, `slice-025`, unverändert)
- `internal/bootstrap/walretention_internal_test.go` (294 Zeilen, Volltext)
  und `walretention_endtoend_test.go` (245 Zeilen, Volltext)
- `docs/plan/adr/0049-replication-fehlerklassen-schwellen.md` (Accepted)
- `spec/pflichtenheft.md` Zeilen 224, 235–245 (Volltext gelesen)
- `docs/user/benutzerhandbuch.md` (Diff `e5e29af`)
- `docs/plan/planning/welle-7.md` (§1, §3, §4)
- `docs/plan/planning/observations/BEO-PGC/spec008-replication-luecke/state.md`
- `git log`/`git show`/`git diff` über den vollen Commit-Verlauf dieses
  Slice (`68134d5..HEAD`)

---

## Sensor-Belege (in diesem Lauf selbst gefahren)

| Sensor | Ausgabe-Kernaussage | Exit |
|---|---|---|
| `make gates` | `baseline-verify: v6.5.0 OK — 54 Dateien` · `d-check` Standardlauf: 239 Dateien, 0 Befunde · `d-check --range HEAD~5..HEAD` (Modul `commits`): 239 Dateien, 0 Befunde · `commit-traceability.sh "HEAD~5..HEAD"`: OK, 5 Commits, Betreffs ohne Struktur-ID · `a-check`: 0 Befunde | **0** |
| `make test-replication` (**dreimal** nacheinander) | alle Pakete `ok`, `.../bootstrap` konstant 18.68s–18.70s (Lauf 1: 18.696s, Lauf 2: 18.679s, Lauf 3: 18.695s), `.../receive` 4.3–4.6s | **0** (jeweils) |
| `go test -race -count=5 -run 'TestRunWALRetentionCheck\|TestMergeStreamAndWALFaultOutcome\|TestResolveWALRetentionThresholds\|TestClassifyWALRetentionThresholds' ./internal/bootstrap/...` (eigener Docker-Lauf mit CGO_ENABLED=1) | alle acht Testfunktionen 5× `PASS`, kein Data-Race-Befund | **0** |
| **Eigene Rot-Grün-Gegenprobe 1** (`mergeStreamAndWALFaultOutcome` Priorität umgedreht: `fault.get()` zuerst geprüft) | **FAIL**: `TestMergeStreamAndWALFaultOutcomePrioritizesStreamError` — „wollen die Stream-Ordnungs-Verletzung unverändert", liefert stattdessen den WAL-Fault-Fehler | **1** (rot, wie erwartet) |
| Datei wiederhergestellt (`git diff` danach leer) | `TestMergeStreamAndWALFaultOutcomePrioritizesStreamError` wieder `PASS` | **0** |
| **Eigene Rot-Grün-Gegenprobe 2** (`resolveWALRetentionThresholds`: beide `<= 0`-Bedingungen auf `>= 0` umgekehrt — exakt die im DoD-Beleg behauptete Mutation) | **FAIL**: beide neuen Tests — `TestResolveWALRetentionThresholdsDefaultsToSpec013/beide_Felder_negativ` liefert `-1` statt `104857600`, `TestResolveWALRetentionThresholdsKeepsPositiveOverride` liefert `104857600` statt des Overrides | **1** (rot, wie erwartet) |
| Datei wiederhergestellt (`git diff` danach leer) | beide Tests wieder `PASS` | **0** |
| `gofmt -l internal/bootstrap/` / `go vet ./internal/bootstrap/...` | keine Ausgabe / `VET_OK` | **0** |
| `git status` (am Ende dieses Laufs) | sauber, keine Restspur der beiden Gegenproben | — |

## DoD-Prüfung (Slice-Plan §2, Punkt für Punkt)

| # | DoD-Punkt | Verdikt | Beleg-Kernaussage |
|---|---|---|---|
| 1 | `SPEC-008`-Zeile `replication` (Transport-/Verbindungsstörungs-Anteil) erfüllt | **bestätigt** | `wiring.go:456-473` (`resolveWALRetentionThresholds`), `:480-489` (`classifyWALRetention`), `:555-583` (`runWALRetentionCheck`): Schwellen-Vergleich löst bei Überschreiten der Fehlerschwelle `fault.set` + `stopStream()` aus; `mergeStreamAndWALFaultOutcome` (`:528-533`) liefert diesen Fehler nur bei regulärem Stream-Ende zurück; `classifyRunError` (`:611-633`) kennt `outbound.ErrReplication` bereits als `replication`-Klasse — kein neuer Klassifikationspfad nötig, eigener Code-Lesung nachvollzogen |
| 2 | Ende-zu-Ende-Test, beide Seiten der Schwelle real, dreimal grün | **bestätigt, selbst dreimal ausgeführt** | `TestWALRetentionThresholdEndToEnd` erzeugt WAL-Wachstum über echte Inserts auf `wal_e2e_noise` (nicht publiziert), belegt Phase 1 (Warnschwelle, `Run()` läuft weiter — durch aktives Warten auf eine zweite verarbeitete Change über 8s bewiesen) und Phase 2 (Fehlerschwelle, `Run()` liefert `outbound.ErrReplication`, explizit **nicht** einer der drei Mapper-Sentinels) in einem Lauf; eigener dreifacher `make test-replication`-Lauf in dieser Sitzung, alle drei grün, `internal/bootstrap` konsistent ~18.7s |
| 3 | Stream-Ordnungs-Verletzungen bleiben hart abbrechend (Regressionstest) | **bestätigt, eigene Rot-Grün-Gegenprobe** | `TestMergeStreamAndWALFaultOutcomePrioritizesStreamError` real rot gesehen bei umgedrehter Priorität, wieder grün nach Wiederherstellung (Sensor-Tabelle oben); zusätzlich eigene statische Prüfung von `receive.go:302-323`: `process()` liefert einen Mapper-Sentinel-Fehler **unbedingt** über `return err` (Zeile 321-323) zurück — kein `ctx.Err()`-Check auf diesem Pfad; nur der `ReceiveMessage`-Fehlerpfad (Zeile 304-308) prüft `ctx.Err()`, trägt aber nie einen bereits entstandenen Mapper-Fehler. Beide Fehlerklassen können sich in einer sequenziellen Goroutine nicht überschneiden — die Priorität in `mergeStreamAndWALFaultOutcome` ist damit sowohl testbelegt als auch durch eigene Code-Lesung nachvollzogen |
| 4 | `make gates` grün | **bestätigt** | eigener Lauf in dieser Sitzung, alle vier inneren Gates 0 Befunde (Sensor-Tabelle) |
| 5 | Review durchgeführt, Report unter `docs/reviews/` liegt vor | **materiell erfüllt, Checkbox nicht nachgezogen — siehe V-1** | `docs/reviews/review-slice-026.md` existiert real (`3c33b1e`), 0 HIGH/1 MEDIUM/0 LOW/0 INFO, „nicht merge-blockierend" — Rollenwechsel fand statt (kein Self-Review); Plan-Checkbox in §2 steht weiterhin `[ ]`, auch nach der Fixrunde (`29a49ea`) |
| 6 | Doku-Update `docs/user/benutzerhandbuch.md` §4/§6/§9/Änderungshistorie 1.4 | **bestätigt** | `git show e5e29af -- docs/user/benutzerhandbuch.md`: §4 nennt exakt dieselben Schwellenwerte (100 MiB/1 GiB) und Log-Feldnamen (`metric`/`bytes`/`threshold_bytes`) wie `wiring.go:570-577`; §6 trennt `replication` korrekt in die beiden `ADR-0049`-Unterarten mit denselben Aktionen wie `spec/pflichtenheft.md:244`; §9 nennt dieselben Zahlenwerte; Änderungshistorie 1.4 vorhanden — Wortlaut gegen den tatsächlichen Log-Aufruf und gegen `SPEC-008` selbst gegengelesen, nicht nur den Reviewer-Befund übernommen |
| 7 | Closure-Notiz mit Steering-Loop-Lerneintrag | **korrekt offen** | §7 trägt weiterhin vollständig den Platzhalter-Vorlagentext — Planner-Arbeit, noch nicht fällig vor `git mv` |
| 8 | Reconciliation-Register, falls einschlägig | **entfällt — korrekt geprüft** | `docs/plan/planning/reconciliation.md` existiert nicht; Repo durchgehend Greenfield (`harness/conventions.md` §Modus-Deklaration); Plan-Text benennt das bereits korrekt |
| 9 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | `BEO-PGC/spec008-replication-luecke/state.md` trägt unverändert „weiter offen", 1× (`evidence/slice-020.md`); Plan-Nachzug benennt bereits korrekt, dass die Bedingung für „eingetreten" erfüllt ist — das Eintragen selbst bleibt Closure-Arbeit, `evidence/slice-026.md` existiert (erwartungsgemäß) noch nicht |
| 10 (Risiken + Paarungen) | Jedes Risiko aus §6 trägt einen Ausgang; drei Paarungen getragen | **korrekt offen** | Beide §6-Risiken tragen weiterhin `<bei Closure einzutragen>` (siehe unten); `slice-026` gehört zu `welle-7` (Kopf-Feld „Welle:") und ist laut `welle-7.md` §4/§3 der **letzte** der drei Slices dieser Welle — die drei Paarungen laufen daher regulär mit der bevorstehenden `welle-7`-Closure, nicht bei dieser Einzel-Slice-Closure (Modul 6 §Wann Arbeit eine Welle braucht) |

**Zwischenstand: 6/10 Kriterien materiell erfüllt und in diesem Lauf selbst
nachgeprüft (inkl. zwei eigener Rot-Grün-Gegenproben und einer eigenen
statischen Nebenläufigkeits-Analyse für Punkt 3), 1 Item korrekt entfallen
(Reconciliation-Register, Greenfield), 3 Items regulär offen als
Planner-/Welle-Closure-Arbeit (Closure-Notiz, Beobachtungs-Register,
Risiko-Ausgänge/Paarungen). Punkt 5 trägt einen eigenen Befund (V-1):
materiell erfüllt, Checkbox nicht nachgezogen — derselbe Musterbefund wie
in `verify-slice-024.md`/`verify-slice-025.md`.**

## Eigene Rot-Grün-Gegenprobe 1 — Prioritäts-Garantie (Detailprotokoll)

1. `internal/bootstrap/wiring.go`, `mergeStreamAndWALFaultOutcome`: Priorität
   umgedreht — `fault.get()` zuerst geprüft, `streamErr` nur als Fallback.
2. `go test -run TestMergeStreamAndWALFaultOutcomePrioritizesStreamError
   ./internal/bootstrap/... -v`: **FAIL** — Test verlangt für jeden der drei
   Mapper-Sentinels, dass er unverändert zurückkommt; mit umgedrehter
   Priorität liefert die Funktion stattdessen den gesetzten WAL-Fault.
3. Datei mit `git show HEAD:internal/bootstrap/wiring.go` wiederhergestellt
   (Kopie vor der Mutation gesichert) — `git diff` danach leer.
4. Derselbe Testlauf erneut: `PASS`.

**Ergebnis:** Der Test beweist die Priorität wirklich, nicht nur die
Existenz einer Fallunterscheidung. Zusätzlich eigene statische Analyse von
`receive.go` (siehe DoD-Punkt 3 oben): Der Mapper-Sentinel-Fehlerpfad in
`Stream.Run` prüft an keiner Stelle `ctx.Err()` — ein gleichzeitig
laufendes `stopStream()` kann eine bereits entstandene Stream-Ordnungs-
Verletzung strukturell nicht in `nil` verwandeln, unabhängig vom Timing.

## Eigene Rot-Grün-Gegenprobe 2 — Zero-Value-Fallback (Detailprotokoll, F-1-Fix)

1. `internal/bootstrap/wiring.go`, `resolveWALRetentionThresholds`: beide
   Bedingungen `<= 0` auf `>= 0` umgekehrt — exakt die im Fixrunden-Commit
   (`29a49ea`) behauptete Mutation.
2. `go test -run TestResolveWALRetentionThresholds
   ./internal/bootstrap/... -v`: **FAIL** — `beide_Felder_negativ` liefert
   `-1` statt `104857600`, `TestResolveWALRetentionThresholdsKeepsPositiveOverride`
   liefert `104857600` statt des 32768-Overrides (die `0`-Variante des
   ersten Falls bleibt zufällig grün, weil `0 >= 0` wahr ist wie `0 <= 0` —
   die negative Variante und der Override-Test sind die diskriminierenden
   Fälle, beide schlagen tatsächlich fehl).
3. Datei wiederhergestellt — `git diff` danach leer.
4. Derselbe Testlauf erneut: beide Tests wieder `PASS`.

**Ergebnis:** Die Fixrunden-Behauptung „Rot-Grün-Beleg: Bedingung testweise
auf `>= 0` umgekehrt, Test rot gesehen" ist durch eigene Wiederholung
bestätigt, nicht nur übernommen. Zusätzlich eigene Prüfung der
Produktionspfad-Anbindung: `ConfigFromEnv` (`wiring.go:144-179`) setzt
`WALRetentionWarnBytes`/`WALRetentionErrorBytes` an keiner Stelle — ein über
`main.go` gestarteter Produktionsprozess erreicht `Run` deshalb immer mit
`Config{}`-Zero-Values für beide Felder, und `resolveWALRetentionThresholds`
liefert dafür nachweislich die `SPEC-013`-Startwerte, nicht `0`.

## Vierte Frage — Leck der `WALRetentionChecker`-eigenen Verbindung bei `stopStream`?

Eigene Prüfung von `wiring.go:364-414` (Reihenfolge der Anweisungen und
`defer`-Registrierungen): `walRetentionCtx`/`stopWALRetention` (Zeile 401)
ist **unabhängig** von `streamCtx`/`stopStream` (Zeile 396) — beide leiten
sich getrennt von `ctx` ab. Löst die Prüf-Goroutine `stopStream()` aus
(Fehlerschwelle überschritten), kehrt sie **selbst sofort** zurück
(`case walRetentionError: … stopStream(); return`, Zeile 573-574) — kein
Warten auf eine externe Kündigung nötig. Unabhängig davon, ob die
Prüf-Goroutine sich selbst beendet oder über eine echte Stream-Fehlerursache
läuft: `Run` ruft nach `stream.Run()`s Rückkehr explizit `stopWALRetention()`
**und** `walRetentionDone.Wait()` auf (Zeile 412-413), **bevor** die Funktion
zurückkehrt und die deferred `walRetention.Close(closeCtx)` (Zeile 368-372)
ausgeführt wird. Die `sync.WaitGroup` garantiert, dass die Prüf-Goroutine
vollständig beendet ist, bevor die Verbindung geschlossen wird — kein
`Close()` gegen eine noch laufende `Measure()`-Anfrage, kein Leck.

Eigene Zusatzprüfung: `WALRetentionChecker.Measure` (`walretention.go:61-81`)
übergibt denselben `ctx` an `pglogrepl.IdentifySystem`/`querySingle` — bei
`stopWALRetention()` bricht ein gerade laufender Protokoll-/Katalog-Aufruf
zeitnah ab, statt bis zum nächsten Tick zu blockieren; `walRetentionDone.Wait()`
wartet in jedem Fall auf die tatsächliche Rückkehr der Goroutine, unabhängig
von der Blockierzeit. **Kein viertes, bisher unbemerktes Problem gefunden.**

## Plan-vs-Code-Diff (§3, gegen `e5e29af` und den Plan-Nachzug gegen `29a49ea`)

| Plan-Zeile (Ursprungsplan) | Behauptung | Abgleich gegen `git show e5e29af` |
|---|---|---|
| `internal/bootstrap/wiring.go` (`classifyRunError`, `Run`) | Schwellen-Vergleich einhängen, nur Transport-/Verbindungsstörungs-Zweig | **stimmt, mit Präzisierung** — `classifyRunError` selbst bleibt unverändert (`outbound.ErrReplication` war bereits als `replication`-Klasse geführt, `slice-025`); der neue Code liegt in `Run` (Verdrahtung der Prüf-Goroutine) und den neuen Funktionen `classifyWALRetention`/`runWALRetentionCheck`/`mergeStreamAndWALFaultOutcome`/`resolveWALRetentionThresholds` — keine Abweichung vom Ziel, nur genauer als die ursprüngliche Plan-Zeile |
| `receive/receive.go` oder Äquivalent | periodischer Schwellen-Check ruft Metrik ab | **abweichend, im Plan bereits vorgesehen als Möglichkeit** — der Check läuft in `internal/bootstrap/wiring.go` (`runWALRetentionCheck`), nicht in `receive.go`; `receive.go` selbst bleibt in diesem Commit unberührt (`git show --stat e5e29af` listet nur `wiring.go`, die beiden neuen Testdateien, `benutzerhandbuch.md`). Der Plansatz „Ort aus `slice-025`" war als Platzhalter formuliert („Implementer entscheidet"); die Composition Root ist der plausible und mit `ADR-0026` konsistente Ort (kennt bereits `walRetentionMeasurer`/`stopStream`) |
| `docs/user/benutzerhandbuch.md` §6 | Zeile `replication` nachgezogen | **stimmt** — siehe DoD-Punkt 6 oben |
| Ende-zu-Ende-Test | neu, beide Seiten der Schwelle | **stimmt** — `walretention_endtoend_test.go`, siehe DoD-Punkt 2 |
| Plan-Nachzug (Fixrunde) `wiring.go` | Fallback-Logik in `resolveWALRetentionThresholds` extrahiert | **stimmt** — `git show 29a49ea -- internal/bootstrap/wiring.go` zeigt genau diese Extraktion, `Run` ruft jetzt `resolveWALRetentionThresholds(cfg.WALRetentionWarnBytes, cfg.WALRetentionErrorBytes)` (Zeile 399) auf, keine Verhaltensänderung sonst |
| Plan-Nachzug (Fixrunde) `walretention_internal_test.go` | zwei neue Tests, Rot-Grün-Beleg | **stimmt** — siehe eigene Gegenprobe 2 oben |

`git show --stat e5e29af` zeigt genau vier Dateien (`benutzerhandbuch.md`,
zwei neue Testdateien, `wiring.go`) — keine unangekündigte fünfte. `git show
--stat 29a49ea` zeigt genau drei Dateien (Slice-Plan, ein Testfile, `wiring.go`)
— ebenfalls deckungsgleich mit dem Plan-Nachzug-Text.

**Out-of-Scope-Disziplin (§1) eigenständig geprüft:** `receive.go` ist in
keinem der drei Slice-Commits (`e5e29af`, `d4382b5`, `29a49ea`) enthalten —
der Stream-Ordnungs-Verletzungs-Zweig (`process()`) ist unverändert, die
Metrik-Erhebung selbst (`slice-025`) ist unberührt, keine
Alerting-Weiterleitung im Diff.

## §6 Risiken — Ausgang bei diesem Verifikationsstand

Beide im Slice-Plan §6 genannten Risiken tragen **weiterhin den
unausgefüllten Platzhalter** `<bei Closure einzutragen>`. Beide sind
inhaltlich bereits materiell beantwortbar (Planner-Urteil bleibt formal
zuständig):

1. „Ein Schwellen-Vergleich, der versehentlich auch den
   Stream-Ordnungs-Verletzungs-Zweig erreicht, würde echte Dateintegritäts-
   Korruption stillschweigend fortsetzen lassen." — Durch die eigene
   Rot-Grün-Gegenprobe 1 und die eigene statische Analyse von `receive.go`
   in diesem Lauf **nicht eingetreten**; die Architektur macht diesen Fehler
   strukturell unmöglich (kein `ctx.Err()`-Check im Mapper-Fehlerpfad).
   Plausibler Ausgang: **entfallen**, Begründung „Regressionstest und
   eigene Nebenläufigkeits-Analyse (Verifikation) zeigen: strukturell
   ausgeschlossen, nicht nur getestet".
2. „Der in `slice-024` vorgeschlagene Byte-Schwellenwert könnte sich beim
   realen Ende-zu-Ende-Test als unpraktikabel erweisen." — Ist **eingetreten**
   in abgeschwächter Form: Die Produktionswerte (100 MiB/1 GiB) selbst
   wurden nicht direkt im Test durchlaufen (das wäre für einen
   CI-tauglichen Testlauf unpraktikabel gewesen), sondern über
   `Config.WALRetentionWarnBytes`/`WALRetentionErrorBytes` durch kleinere
   Testfixture-Werte ersetzt — dieser Override-Mechanismus war jedoch
   bereits im *ursprünglichen* Slice-Plan (§2, vor jeder Implementierung)
   vorgesehen, nicht erst als Reaktion auf ein gescheitertes
   Testexperiment. Plausibler Ausgang: **entfallen**, Begründung
   „Override-Mechanismus von Anfang an eingeplant, kein nachträglicher
   Kompromiss — die Produktions-Startwerte selbst bleiben unverändert und
   ungetestet in Echtzeit, das ist eine bewusste Testzeit-Abwägung, kein
   eingetretenes Risiko". Alternative Lesart („weiter offen": die realen
   100 MiB/1 GiB-Werte sind nie in einem echten Lauf durchlaufen worden) ist
   ebenfalls vertretbar — Planner-Urteil.

## Negativbefunde

- geprüft, ohne Befund: **`make gates`** — eigener Lauf, 0 Befunde in allen
  vier Gates.
- geprüft, ohne Befund: **`make test-replication`, dreimal** — eigener
  Lauf, alle Pakete grün, konsistente Laufzeiten (`internal/bootstrap`
  18.68s–18.70s über alle drei Läufe).
- geprüft, ohne Befund: **`go test -race -count=5`** auf allen acht
  `internal/bootstrap`-Testfunktionen — kein Data-Race.
- geprüft, ohne Befund: **Eigene Rot-Grün-Gegenprobe Prioritäts-Garantie**
  — Test scheitert sichtbar bei umgedrehter Priorität, besteht mit ihr;
  kein Restschaden am Arbeitsverzeichnis nach Wiederherstellung.
- geprüft, ohne Befund: **Eigene Rot-Grün-Gegenprobe Zero-Value-Fallback**
  — dieselbe Mutation wie im Fixrunden-Commit behauptet, real reproduziert;
  kein Restschaden nach Wiederherstellung.
- geprüft, ohne Befund: **Kein Verbindungsleck bei `stopStream`** — eigene
  Lese-Analyse der `defer`-/`WaitGroup`-Kette in `Run`, siehe „Vierte Frage"
  oben.
- geprüft, ohne Befund: **`gofmt -l` / `go vet`** auf `internal/bootstrap` —
  keine Befunde.
- geprüft, ohne Befund: **Realer Ende-zu-Ende-Test, kein Vortäuschen** —
  eigene Lektüre von `walretention_endtoend_test.go`: WAL-Wachstum kommt
  über echte Inserts auf einer nicht publizierten Tabelle zustande, beide
  Schwellenseiten werden real und nacheinander in einem Lauf durchlaufen,
  die Fehlerklassifikation wird explizit gegen alle drei Mapper-Sentinels
  ausgeschlossen (Zeile 190-192).
- geprüft, ohne Befund: **Doku-Konsistenz** — Log-Feldnamen und -Werte im
  Benutzerhandbuch stimmen wörtlich mit `wiring.go` überein, Zahlenwerte
  stimmen mit `SPEC-013` überein.
- geprüft, ohne Befund: **Traceability** — beide inhaltstragenden
  Commit-Betreffs (`e5e29af`, `29a49ea`) nennen `ADR-0049`, kein
  `SPEC-*`/`ARC-*` im Betreff; `make commit-traceability` (Teil von `make
  gates`) lief grün.
- geprüft, ohne Befund: **Beobachtungs-Register unverändert wie
  angekündigt** — `state.md` trägt weiterhin „weiter offen", 1×; keine
  unangekündigte Änderung durch diesen Slice.
- geprüft, ohne Befund: **Plan-vs-Code-Diff** — vollständige Deckung, keine
  unangekündigte Datei in `e5e29af` oder `29a49ea`.

## Eigene Befunde

### V-1 — DoD-Checkbox „Review durchgeführt" nicht nachgezogen, obwohl die Bedingung materiell erfüllt ist

- `kategorie`: LOW
- `pfad`: `docs/plan/planning/in-progress/slice-026-schwellen-ueberwachung-kontrollierte-fortsetzung.md:107-109`
  (Checkbox weiterhin `[ ]`) vs. `docs/reviews/review-slice-026.md`
  (Report existiert, 0 HIGH/1 MEDIUM, „nicht merge-blockierend")
- `befund`: Derselbe Musterbefund wie `verify-slice-024.md`/
  `verify-slice-025.md` V-1: Die DoD-Zeile „Review durchgeführt, Report
  unter `docs/reviews/` liegt vor … kein Self-Review" ist materiell
  erfüllt (Rollenwechsel fand statt, Report liegt vor, 0 HIGH), die
  Checkbox wurde weder im Review-Commit (`3c33b1e`) noch in der Fixrunde
  (`29a49ea`) auf `[x]` gesetzt. Diese Finding-Klasse ist bereits in
  `BEO-PGC/dod-checkbox-nachzug` als **verkörpert** geführt
  (`.claude/commands/implement-slice.md`, seit `welle-5`); hier betrifft es
  denselben wiederkehrenden Fall wie bei `slice-024`/`slice-025`: eine
  materiell erfüllte Review-Bedingung, deren Checkbox nicht mitgezogen
  wurde.
- `verifizierbar`: ja — die Checkbox bleibt `[ ]`, kein Commit nach
  `3c33b1e` ändert diese Zeile.
- **Für die Closure:** kein Blocker — der Planner sollte die Checkbox vor
  `git mv` nach `done/` auf `[x]` setzen (materiell gedeckt).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 (V-1) |
| INFO | 0 |

**Zusammenfassung DoD:** 6/10 Kriterien materiell erfüllt und in diesem
Lauf selbst geprüft (inkl. zweier eigener Rot-Grün-Gegenproben und einer
eigenen statischen Nebenläufigkeits-Analyse), 1 Item korrekt entfallen
(Reconciliation-Register, Greenfield), 3 Items regulär offen als
Planner-/Welle-Closure-Arbeit (Closure-Notiz, Beobachtungs-Register,
§6-Risiko-Ausgänge, drei Paarungen — letztere laufen turnusgemäß mit der
`welle-7`-Closure). Kein DoD-Defekt im Sinn eines unbelegten
„bestätigt"-Punkts — der einzige eigene Befund (V-1, LOW) ist eine
Checkbox-Diskrepanz, die die materielle Substanz nicht widerlegt.

## Verdikt

**DoD-/Entscheidungs-Konformität: bestätigt, mit einem offenen
Klein-Befund (LOW, V-1; non-blocking für diesen Slice, empfohlene Korrektur
vor `git mv` nach `done/`).** Die Prioritäts-Garantie zwischen
Stream-Fehlern und WAL-Schwellen-Fehler — der sicherheitskritische Kern
dieses Slices — hält sowohl unter eigener Rot-Grün-Gegenprobe als auch
unter eigener statischer Nebenläufigkeits-Analyse von `receive.go`: ein
WAL-Schwellen-Fehler kann eine Stream-Ordnungs-Verletzung strukturell nicht
maskieren. Die Fixrunde hat das Reviewer-Finding F-1 tatsächlich behoben,
nicht nur behauptet: `resolveWALRetentionThresholds` sitzt nachweislich im
Produktionspfad (`Run` ruft sie mit den `ConfigFromEnv`-Werten auf, die für
beide Felder immer Zero-Value tragen), und die eigene Rot-Grün-Gegenprobe
reproduziert exakt die im Commit behauptete Mutation. Der Ende-zu-Ende-Test
ist real (kein simuliertes WAL-Wachstum), durchläuft beide Seiten der
Schwelle in einem Lauf und wurde in dieser Sitzung dreimal grün
reproduziert. Kein viertes, bisher unbemerktes Problem gefunden — die
`WALRetentionChecker`-eigene Verbindung wird über eine `sync.WaitGroup`
korrekt synchronisiert geschlossen, kein Leck.

**Plan-vs-Code-Diff:** vollständige Deckung, eine benannte, unschädliche
Abweichung (Ort des periodischen Checks: Composition Root statt
`receive.go` — plausibel und mit `ADR-0026` konsistent, im Ursprungsplan
als offene Entscheidung markiert). Out-of-Scope-Punkte aus §1 gewahrt
(kein Eingriff in `process()`, keine Änderung an `slice-025`s
Metrik-Erhebung, keine Alerting-Weiterleitung).

**Closure-Bereitschaft: noch nicht — reguläre Planner-/Welle-Closure-Schritte
stehen aus**, keiner davon ein Defekt an bereits gelieferter Substanz:

1. Closure-Notiz §7 schreiben (Was hat funktioniert / anders als geplant /
   Steering-Loop-Eintrag / Beobachtungs-Register / Folge-Slices / §6-Risiken).
2. Beide §6-Risiken auf einen der drei Ausgänge setzen — plausible Ausgänge
   oben skizziert (beide „entfallen" mit Begründung, oder Risiko 2 „weiter
   offen"), formales Setzen bleibt Planner-Urteil.
3. Beobachtungs-Register: `BEO-PGC/spec008-replication-luecke` auf
   „eingetreten" setzen (die Bedingung — Ende-zu-Ende-Beleg über beide
   Seiten der Schwelle — ist durch diesen Lauf eigenständig dreifach
   bestätigt) und `evidence/slice-026.md` anlegen.
4. Drei Paarungen: laufen turnusgemäß mit der bevorstehenden
   `welle-7`-Closure (letzter Slice dieser Welle, `slice-024`/`slice-025`
   liegen bereits in `done/`) — nicht bei dieser Einzel-Slice-Closure.

Zusätzlich empfohlen: V-1 (Review-Checkbox) vor `git mv` nach `done/`
beheben.

**Übergabe:** Bericht an den Planner. Keine Reparaturen — Plan-Datei und
Code wurden von diesem Lauf nicht dauerhaft verändert (beide
Rot-Grün-Gegenproben wurden vollständig rückgängig gemacht, `git status`
am Ende sauber).

---

**Gate-Beleg:** `make gates` einmal in diesem Lauf ausgeführt, Exit 0, 0
Befunde (`baseline-verify`: 54 Dateien OK; `d-check` Standardlauf: 239
Dateien, 0 Befunde; `d-check --range HEAD~5..HEAD` Modul `commits`: 239
Dateien, 0 Befunde; `commit-traceability.sh`: OK, 5 Commits;
`a-check`: 0 Befunde). `make test-replication` **dreimal** vollständig
ausgeführt, alle Pakete grün, konsistente Laufzeiten. Zusätzlich
`go test -race -count=5` auf den acht `internal/bootstrap`-Testfunktionen
(kein Data-Race) und zwei gezielte Rot-Grün-Gegenproben direkt gegen
`wiring.go` (Priorität, Fallback), beide mit vollständiger Wiederherstellung.
`git status` am Ende dieses Laufs sauber (keine Arbeitsverzeichnis-Änderung
durch die Verifikation selbst).
