# Verifikations-Report: slice-backfill-run-usecase — 2026-09-24

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich +
Entscheidungs-Konformität + Plan-vs-Code-Diff + Gates. Review-Artefakt des
Reviewers: [`review-slice-backfill-run-usecase.md`](review-slice-backfill-run-usecase.md);
Formvorbild dieses Reports:
[`verifikation-slice-backfill-snapshot-reader.md`](verifikation-slice-backfill-snapshot-reader.md).

**Gegenstand:** Slice-Plan `slice-backfill-run-usecase` (Welle
`welle-backfill-bestand`), Diff-Range `455bcdef..HEAD` (`e01c9f6e`), 10 Commits:
Lifecycle und Verantwortlich (`d2e24986`, `3df8565d`, `643582b0`), Umsetzung
(`510b247e`, `9d4e9de8`, `a1ffb9be`, `f4ae3b2d`), Review-Report `882bdccd`,
Fixrunde (`b51a5494` Tests, Fakes, Port-Doc-Kommentare; `e01c9f6e` Plan und ein
minimaler Nachzug im Plan `slice-backfill-run-store`). Dieser Lauf ändert weder
Code noch Plan, Spec oder Doku; er schreibt nur diesen Report. Alle Mutationen
liefen im Arbeitsbaum und wurden per `git checkout -- internal` zurückgenommen;
`git status --short` und `git diff` waren nach jedem Lauf und am Ende leer.

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — [`AGENTS.md`](../../AGENTS.md) §3.9)

Jeder Lauf schrieb in eine Log-Datei; der Exit-Code wurde im selben Aufruf
gesondert gesichert (`make … > <log> 2>&1; echo …=$?`), die Logs danach gelesen.

| Sensor | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `make gates` (Stand `e01c9f6e`, vor dem Report-Commit) | **EXIT=0** | baseline-verify `v6.9.0 OK — 54 Dateien` · `db-package-lists-check: OK — … dieselben 4 Pakete` · `coverage-gate: OK — Coverage 84.50% erfüllt Schwelle 80%` · `d-check: 1016 Datei(en) geprüft, 0 Befund(e)` · `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD"` · `generated-sync: OK` · a-check `gesamt: 0 Befund(e)` |
| `make test` (Race-Detector) | **EXIT=0** | `ok … internal/application/usecase/backfill`, kein `FAIL` |
| `make a-check` | **EXIT=0** | `gesamt: 0 Befund(e)` |
| `make coverage-gate` | **EXIT=0** | `total: (statements) 84.4%`, `coverage-gate: OK — Coverage 84.40% erfüllt Schwelle 80%` (der `make gates`-Lauf druckte 84.50 %: dieselbe Codebasis, zwei Läufe, Lauf-Belege, V-4) |
| `make commit-traceability RANGE=455bcdef..HEAD` | **EXIT=0** | `OK — 10 Commit(s) in "455bcdef..HEAD", Betreffs ohne Struktur-ID`; `git log --format=%s 455bcdef..HEAD` auf `SPEC-`/`ARC-`: 0 Treffer; jede Message trägt `LH-*`/`ADR-*` (Schleife über alle zehn, keine „untraceable") |
| `make doc-commits RANGE=455bcdef..HEAD` | **EXIT=0** | 1016 Dateien, 0 Befunde |
| `make doc-immutable RANGE=455bcdef..HEAD` | **EXIT=0** | 1016 Dateien, 0 Befunde |
| `make doc-trace` (advisory, nur zur Kenntnis) | EXIT=0 | 80 Anforderung(en), 3 Waise(n): `LH-FA-CAP-009`, `LH-FA-CFG-007`, `LH-FA-CFG-008` — die Backfill-Waise bleibt bis zum Liefer-Slice mit Datenbank erwartet |

**Kein DB-Sensor nötig, bestätigt am Diff:** `git diff --name-only 455bcdef..HEAD`
nennt außer Plan-Dateien und Review-Report nur `internal/domain/…`,
`internal/application/…` und eine Zeile in
`internal/adapters/driving/replication/mapper/mapper.go`; kein `postgres*`-Paket,
kein Skript, kein Dockerfile, kein Workflow. `make test-store`/`make test-replication`
wurden deshalb nicht gefahren; die Mapper-Tests laufen in `make test`.

## 2. DoD — Verdikt je Zeile (§2 des Plans)

| # | DoD-Zeile | Verdikt | Realer Beleg |
|---|---|---|---|
| 1 | Happy Path gegen Fakes | **erfüllt** | `TestExecuteWritesAllBlocksInOneTransaction` (eine Transaktion, ein Commit, drei Blöcke `0bf-run-1-0000000{1,2,3}`, Position `X` je Block, `INSERT`/`backfill`, Bild ohne ausgeschlossene Spalte, ein Wecksignal), `TestExecuteEmptyTable` (`completed`, 0 Zeilen, keine Transaktion); `make test` EXIT=0 mit Race-Detector. Mutationen M17-Klasse, M19, M20 rot (§4) |
| 2 | Negative gegen Fakes, je an ihrer Eingabe gebunden | **erfüllt** | `TestExecuteClassifiesFailures` (17 Fälle, fünf Snapshot-Sentinels je ein Fall), `TestExecuteFailClosed` (Zwischenabweichung mit gleichem Ende, Einschluss, Ordnung/Doppelung keine Abweichung, Bindung fehlt/andere Kennung/nicht lesbar), `TestExecuteContextEnded…`, `TestExecuteRechecksPreconditions`, `TestRequestPreconditionFailuresLeaveNoRun` (`openCalls == 0`, keine Schätzung, kein `Admit`). Die Mutation je Prüfung färbt den Test rot: 20 Mutationen, alle rot (§4) |
| 3 | Annahme gegen Fakes | **erfüllt** | `TestRequestAdmitsLastAfterPreconditionsAndEstimate` (Reihenfolge `Registered`, `Published`, `EstimatedRows`, `Admit`; kein Snapshot), `TestRequestCarriesEstimate` (vier Fälle, „unbekannt" ≠ 0, bekannte 0), `TestRequestActiveRunEndsWithoutSecondRow`, `TestPortSchnitt` (`BackfillAdmissionPort` = `Admit`, `BackfillRunPort` ohne Anlage); `make a-check` EXIT=0; Port-Definitionen im Diff gelesen (§3) |
| 4 | Zeilenzahl im Speicher durch `B` je Block begrenzt | **erfüllt** | `TestExecuteStreamsBlocks` prüft die Fake-Aufruf-Reihenfolge (Schreiber erhält Block `n`, bevor der Leser `n+1` liefert); der Code hält `rows` nur je Schleifendurchlauf (`copyBlocks`); `make a-check`, `make coverage-gate` grün. Der Plan macht keine Aussage über Bytes, nur über Zeilen |
| 5 | `make gates` grün, Exit gesondert | **erfüllt** | eigener Lauf EXIT=0 (§1) |
| 6 | Review durchgeführt, Report liegt vor | **erfüllt, mit V-3 (INFO)** | Report liegt vor (1 HIGH, 3 MEDIUM, 3 LOW, 6 INFO). Die Fixrunde ist ohne Nachprüfungs-Review abgeschlossen; F-1…F-7 und F-12 habe ich selbst am Code und an den Tests nachgemessen (§4, §6). Der Haken steht nach der Fixrunde (Skill-Regel: bei Fixrunde regulär beim Implementer-Nachzug), nicht im Report-Commit — regelkonform |
| 7 | §3.13-Suchlauf im committeten Feld | **erfüllt** | Zahlen aller Zeilen an den genannten Ständen nachgemessen (§5); Gefundenes und Nichtgefundenes je Träger; Übergaben als Meldung statt Mitänderung |
| 8 | Doku-Update entfällt (Handbuch unberührt) | **erfüllt** | `git diff --name-only 455bcdef -- docs/user/benutzerhandbuch.md` leer; keine `CDC_*`-Variable, keine SQL-Funktion, kein Endpunkt, kein `make`-Target im Diff |
| 9 | Closure-Notiz mit Lerneintrag | **korrekt offen** | Plan §7 trägt „*(zu tragen bei Closure)*" |
| 10 | Reconciliation-Register — entfällt | **erfüllt** | Greenfield, `[x]` (anders als im Vorgänger-Report V-5 des Slice `backfill-snapshot-reader`) |
| 11 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | Closure-Pflicht; der Plan nennt die gesichteten Einträge in §8 |
| 12 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | zehn Risiko-Punkte in §6 (gezählt), alle tragen `Ausgang: (bei Closure)` (neun auf einer Zeile, der dritte mit Zeilenumbruch) |
| 13 | Drei Paarungen | **korrekt offen** | hängen an der Closure von `welle-backfill-bestand` |

Kein `[x]` ohne Beleg; kein `[ ]`, das über die Rollen-Sequenz hinaus belegt wäre.

## 3. Plan-vs-Code-Diff und Fixrunden-Diff

**Fixrunde `882bdccd..b51a5494`, `.go`-Dateien** (`git diff 882bdccd..b51a5494 --stat`):
fünf Dateien. Die **Nicht-Test-Änderungen** sind vollständig benannt:

| Datei | Änderung |
|---|---|
| `internal/application/usecase/backfill/service.go` | nur Kommentare: Doc von `copyBlocks` (Grenze der Prüfung, Lesefehler wie Abweichung) und Doc von `conclude` (Zusage des Use Cases statt Verhalten eines anderen Slice, Rang-Zeiger auf Festlegung 2 von [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md)); **keine Code-Zeile**, das Diff-Hunk enthält ausschließlich `//`-Zeilen |
| `internal/application/port/outbound/backfillrun.go` | Doc: Zeitbegrenzung der Operationen, `Finish` auf einen beendeten Run = wirkungsloser Erfolg, `InterruptRunning` zeitbegrenzt |
| `internal/application/port/outbound/backfillwriter.go` | Doc zu `Rollback` (abgelöster Kontext, Adapter begrenzt selbst) |
| `internal/application/port/outbound/tablesnapshot.go` | Doc zu `Close` (dasselbe) |

`service_test.go` (+257/−16) ist der Test-Anteil: Fakes mit Fehler ab dem n-ten
Aufruf, schrittweise laufende Uhr, vier neue Tests. **Kein geändertes
Produktionsverhalten** — bestätigt; die Mutationen (§4) treffen den unverändert
gebliebenen Code der Erstlieferung.

**Plan §3 gegen den Diff `455bcdef..HEAD`** (15 Dateien): jede Zeile der
Plan-Tabelle ist im Diff vertreten; ungeplante Dateien gibt es nicht (der
Review-Report und der fremde Plan `slice-backfill-run-store` sind im Plan benannt).

- **Abweichungen vom Plan, im Plan benannt und begründet:** (1) `Ports` trägt
  zusätzlich `SchemaStorePort` als achten Pflicht-Port (`model.NewChange` verlangt
  eine Schema-Version-Referenz); (2) `model.ChangeIDFor`, vom WAL-Mapper gerufen
  (+1 Zeile, `git diff 455bcdef..HEAD -- internal/adapters/driving/replication/mapper/`
  zeigt genau `model.ChangeID(fmt.Sprintf("%s-%d", a.open.tx.ID, sequence))` →
  `model.ChangeIDFor(a.open.tx.ID, sequence)`); (3) weitere Sentinels in
  `internal/domain/errors/` (Statuswechsel, Fortschritts-Rückschritt, negative
  Zeilenzahl, Blocknummer-Überlauf); (4) Klasse `schema` wird nicht vergeben,
  `ErrSchemaVersionUnknown` → `configuration`, nicht erkannter Fehler → `internal`.
  Alle vier sind ADR-konform (§7).
- **Byte-Gleichheit der WAL-Change-IDs:** `ChangeIDFor` ist
  `fmt.Sprintf("%s-%d", tx, sequence)` — dieselbe Formatzeichenfolge wie der Parent.
  Mutation M19 (`_` statt `-`) färbt `TestConsumeFullTransaction` und
  `TestConsumeTwoTransactions` im Mapper-Paket rot (§4), die Mapper-Tests sind
  unverändert grün.
- **Port-Schnitt gelesen:** `BackfillAdmissionPort` hat genau `Admit`;
  `BackfillRunPort` hat `Queued`, `MarkRunning`, `RecordProgress`, `Finish`,
  `InterruptRunning` und **keine** Anlage; `BackfillWriterPort.Begin` mit
  `BackfillTransaction` (`AppendBlock`, `Commit`, `Rollback`); der Inbound Port
  trägt `Request` und `Execute` mit Transport-Typen am Port
  ([`ADR-0042`](../plan/adr/0042-transport-typen-am-port.md)). `usecase/backfill`
  importiert nur `port/inbound`, `port/outbound` und Domäne
  ([`ADR-0002`](../plan/adr/0002-abhaengigkeitsrichtung.md),
  [`ADR-0003`](../plan/adr/0003-physische-modulgrenzen.md); `make a-check` EXIT=0).
- **Port-Verträge und ihr Träger:** `Finish` auf einen beendeten Run ist ein
  wirkungsloser Erfolg (Doc-Kommentar, `backfillrun.go`); die Zeitbegrenzungs-Pflicht
  steht in den Doc-Kommentaren von `BackfillRunPort`, `InterruptRunning`,
  `BackfillTransaction.Rollback` und `TableSnapshot.Close`. Der Träger im Plan
  `slice-backfill-run-store` (`git diff 455bcdef..HEAD` auf diese Datei: **+16
  Zeilen, 0 gelöscht**): ein DoD-Punkt „Adapter-Pflichten aus den Port-Verträgen"
  mit *Zu belegen durch* (`make test-store`, netzloser Test an der Ausführungs-Naht)
  und eine §3-Zeile. Bestehender Text, Absicht und übrige Zusagen des Plans sind
  unverändert (nur Ergänzung). Bestätigt: der Träger ist vorhanden. Ebenfalls
  vorhanden und unberührt (F-13 des Reviews): der DoD-Punkt „Ordnung und Zustand"
  (`slice-backfill-run-store`, „mit einem WAL-Commit **auf derselben Position** `X`
  liest `cdc.changes` die Backfill-Blöcke **vor** dem WAL-Commit … Kollation der
  Datenbank (`datcollate`)") und ein Risiko-Punkt „Die Textsortierung trägt die
  Ordnung `0bf-…`" — die reale `ORDER BY`-Sortierung trägt also `slice-backfill-run-store`,
  nicht dieser Slice.
- **Ordnung im Slice selbst:** `BackfillTransactionID` bildet `0bf-<run>-%08d`
  (Blocknummer ab 1, Überlauf ab 100000000 ein Fehler); `TestBackfillTransactionID`
  bindet den Byte-Vergleich von Go; Mutation M20 (`bf-`) färbt ihn rot. Dass
  WAL-Kennungen (XID als Dezimalzahl, `strconv.FormatUint`) nie mit `0` beginnen,
  ist im Reviewer-Report hergeleitet; die Datenbank-Sortierung gehört `run-store`.

## 4. Mutationen (Eingabeseite, selbst gesehen)

Je Mutation: exakter Austausch einer Zeile (Skript mit Zählprüfung „genau ein
Treffer"), `make test`, Exit-Code und `FAIL`-Zeilen gelesen,
`git checkout -- internal`. Ausgangslauf ohne Mutation: `make test` EXIT=0.
**20 Mutationen, alle rot (EXIT=2).** Eine erste Fassung von M10 scheiterte am
Übersetzen (`declared and not used`) und ist in korrigierter Form (M10b) gezählt.

| # | Zusage | Mutation an `service.go` (`Z:` = Zeile im Stand `e01c9f6e`) | rot: |
|---|---|---|---|
| M1 | `finished_at` von `interrupted` | `Interrupt(model.TimePoint{})` statt der Uhr (`conclude`, Z:295) | `TestExecuteFinishedAt/interrupted` |
| M2 | `finished_at` von `failed` | `Fail(model.TimePoint{}, …)` (Z:297) | `TestExecuteFinishedAt/failed_beim_Kopieren`, `…/failed_bei_der_erneuten_Prüfung` |
| M3 | `finished_at` der leeren Tabelle | `Complete(model.TimePoint{}, 0)` (Z:247) | `TestExecuteFinishedAt/completed,_leere_Tabelle` |
| M4 | **Lesefehler unmittelbar vor dem Commit** | `excluded, _ := s.excludedColumns(…)` (Z:260) | `TestExecuteExclusionReadFailure/unmittelbar_vor_dem_Commit` |
| M5 | **Lesefehler je Block** | `excluded, _ := s.excludedColumns(…)` (Z:220) | `…ReadFailure/vor_dem_ersten_Block`, `…/vor_dem_zweiten_Block`, `…/vor_dem_letzten_Block` |
| M6 | Fortschritts-Fehler je Block | `run, _ = s.progress(…, copied)` (Z:241) | `TestExecuteProgressFailure/nach_dem_ersten`, `…zweiten`, `…letzten_Block` |
| M7 | Fortschritts-Fehler mit der Anfangsposition | `run, _ = s.progress(…, 0)` (Z:196) | `TestExecuteProgressFailure/mit_der_Anfangsposition` |
| M8 | `Finish`-Fehler bei leerer Tabelle | Fehler von `Finish(completed)` verworfen (Z:251) | `TestExecuteEmptyTableFinishFailure/erstes_Finish_scheitert`, `…/jedes_Finish_scheitert` |
| M9 | Vergleich je Block | `else if false` (Z:229) | `TestExecuteFailClosed/Zwischenblock_weicht_ab,_am_Ende_wieder_gleich` |
| M10b | Vergleich vor dem Commit | Vergleich entfernt, `_ = excluded` (Z:264) | `TestExecuteFailClosed/Abweichung_erst_vor_dem_Commit`, `…/Einschluss_statt_Ausschluss` |
| M11 | `interrupted` gegen `failed` am eigenen Kontext | Zweig `ctx.Err() != nil` → `false` (Z:294) | `TestExecuteFinishedAt/interrupted`, `TestExecuteContextEndedInterrupts` (3 Fälle) |
| M12 | `queued` bleibt `queued` bei beendetem Kontext | Zweig → `false` (Z:292) | `TestExecuteContextEndedBeforeStartLeavesQueued` |
| M13 | `Finish` auf abgelöstem Kontext | `Finish(ctx, final)` statt `context.WithoutCancel(ctx)` (Z:302) | `TestExecuteContextEndedInterrupts` (3 Fälle) |
| M14 | `Rollback` auf abgelöstem Kontext | `writer.Rollback(ctx)` (Z:432) | `TestExecuteContextEndedInterrupts` (2 Fälle) |
| M15 | `Close` auf abgelöstem Kontext | `snapshot.Close(ctx)` (Z:421) | `TestExecuteContextEndedInterrupts` (2 Fälle) |
| M16 | Bindung vor dem Commit (`stillBound`) | Aufruf entfernt (Z:257) | `TestExecuteFailClosed` (3 Bindungs-Fälle) |
| M17 | Rollback auf jedem Pfad | `defer`-Rollback entfernt (Z:208) | `TestExecuteClassifiesFailures` (4 Fälle), `TestExecuteMidCopyFailureKeepsProgressAndWritesNothing`, `TestExecuteFailClosed/Ausschluss_ab_Block_2` |
| M18 | Commit-Fehler endet den Run | Fehler von `Commit` verworfen (Z:271) | `TestExecuteClassifiesFailures/Schreiber_Commit`, `…/ChangeStore-Fehler_am_Schreiber` |
| M19 | Byte-Gleichheit der Change-Kennung | `ChangeIDFor` mit `_` (`change.go`) | `TestChangeIDFor`, `TestExecuteWritesAllBlocksInOneTransaction`, `TestConsumeFullTransaction`, `TestConsumeTwoTransactions` (Mapper) |
| M20 | Ordnungs-Präfix | `bf-` statt `0bf-` (`backfillrun.go`) | `TestBackfillTransactionID`, `TestExecuteWritesAllBlocksInOneTransaction` |

**Zu den acht Review-Mutationen (Nr. 24–31 im Review):** `finished_at` in drei
Wegen (M1–M3), Lesefehler je Block und vor dem Commit (M4, M5), Fortschritts-Fehler
in zwei Wegen (M6, M7) und `Finish`-Fehler bei leerer Tabelle (M8) sind alle acht
rot — nicht nur fünf. Die beiden Lesefehler-Mutationen (Kern von
[`LH-QA-SEC-004`](../../spec/lastenheft.md)) färben nun je den eigenen Fall: M4 allein
den Fall der Lesung vor dem Commit, M5 die drei Blockfälle — die Bindung ist
nicht mehr über den jeweils anderen Aufruf maskiert.

**Fake-Realismus („grün ohne Aussage"?).** Geprüft, ob die Assertion das Erwartete
prüft und nicht nur „kein Panic":

- `fakeExclusion.errCall`, `fakeRuns.progressErrCall`/`finishErrCall` lassen nur
  den n-ten Aufruf scheitern (Zähler `calls`, `len(progresses)`, `len(finished)`);
  `errCall == 0` bewahrt das „jeder Aufruf scheitert" der alten Tests.
- `TestExecuteExclusionReadFailure` bindet je Fall Klasse (`storage: `-Präfix und
  Ursachentext), zuletzt festgehaltenen Zähler (`RowsCopied` 0/2/4/5), die **Anzahl
  der Lesungen** (`exclusion.calls == call`, Abbruch mit dem Lesefehler), Zahl der
  angehängten Blöcke, kein Commit, kein Wecksignal, Zahl der Rollbacks, Snapshot
  genau einmal geschlossen und `Finish` genau einmal mit dem Ergebnis-Run.
- `TestExecuteProgressFailure` bindet zusätzlich „kein Fortschritt nach dem
  Fehler" (`len(progresses) == call`); `TestExecuteEmptyTableFinishFailure` bindet
  die Reihenfolge der Endzustände (`completed` dann `failed`) und den Ergebnis-Run
  `running`, wenn auch das zweite `Finish` scheitert.
- `TestExecuteFinishedAt` liest die Uhr mit `step = 10` je Aufruf, vergleicht
  `finished_at` mit der **letzten Lesung** (`clock.last`), mit dem Wert, den der Port
  festhält (Commit oder `Finish`), und `finished_at` **nach** `started_at`.
  Grenze: die Bindung „letzte Lesung der Uhr" würde eine zusätzliche Lesung nach dem
  Endzustand als Abweichung melden — sie ist damit strenger als die Zusage, nicht
  lascher.

## 5. Zahlen der Suchlauf-Felder (§3.12), an den genannten Ständen selbst gemessen

| Plan-Zahl | Meine Messung | Ergebnis |
|---|---|---|
| Zeile 1: `git grep -n -i backfill <Stand> -- 'internal/**/*.go'` — `643582b0`: 80 Zeilen/18 Dateien; `882bdccd`: 362/26; `b51a5494`: 386/26 | 80/18, 362/26, 386/26 | bestätigt |
| Zeile 1: drei Fundstellen der alten Bedeutung (`mapper.go:334`, `schemastore.go:106`, `queries.go:286`) und `TestConsumeRelationBackfillsMissingTableSchema` (`mapper_test.go:349`), an allen drei Ständen gleich | drei Zeilen an allen drei Ständen; Test bei `mapper_test.go:349` | bestätigt |
| Zeile 2: `ports/outbound`/`port/outbound` in `docs spec harness` — `643582b0`: 44 Zeilen/26 Dateien; `882bdccd`/`b51a5494`: 52/27; ohne historische Träger nur dieser Plan (3 Zeilen) | 44/26, 52/27, 52/27; nach dem Filter je genau drei Zeilen aus dem Plan (172, 196, 223) | bestätigt |
| Zeile 3: `classifyRunError` in `wiring.go:1378`; `git diff --stat 643582b0 b51a5494 -- internal/bootstrap` leer; Verweis-Befehl 0 Treffer | `wiring.go:1378`; Diff leer; 0 Treffer | bestätigt |
| Zeile 4 (a): 11 Zeilen an `643582b0` und `882bdccd`; `run-store` 46/54/59, `sql-administration` 52/54, `bench-richtgroesse` 63/155, `e2e` 84/89/134 | 11 Zeilen an beiden Ständen; Zeilennummern wie genannt | bestätigt — mit V-2 (Escape `\|` im Befehl) |
| Zeile 5 (b): `slice-transformationen-*` je 6 Treffer; `antragsweg-usecase` 156/161/167, `backfill-pfad` 73/135/150 | 6 an beiden Ständen, Zeilennummern wie genannt | bestätigt — mit V-2 |
| Zeile 6 (Adapter-Pflichten): außerhalb dieses Plans die Rollback-Treffer `run-store` 52/63/105/112/219/223, `spec/architecture.md:293`, `benutzerhandbuch.md:595`, `harness/README.md:151`, `harness/targets/schema-rollout.md:61`, `releasing.md:315`/`:317`, `lastenheft.md:112`, `pflichtenheft.md:57`; nur `postgressnapshot` implementiert einen der Ports (`snapshot.go:70`, `:262`), `closeTimeout` `snapshot.go:46`, `closeConn` `:319` | exakt diese 14 Fundorte; Adapter-Treffer `snapshot.go:70`/`:262`; `closeTimeout = 5 * time.Second` in Z. 46, `func closeConn` in Z. 319 | bestätigt |
| Nicht gefunden: ein Träger von `BackfillRunPort`/`BackfillWriterPort` außerhalb der Pläne | keiner in `internal/adapters` | bestätigt |

## 6. Kernaussagen, selbst am Code gelesen

- **Fail-closed und Atomarität:** `copyBlocks` liest je Block den Ausschlussstand
  neu; der erste Block setzt die Baseline (`Begin` danach), jeder weitere Block und
  der Zustand unmittelbar vor dem Commit müssen mengengleich sein
  (`sameNames`, Reihenfolge und Doppelung unerheblich); ein Lesefehler ist jeweils
  ein Fehler des Laufs. Vor dem Commit prüft `stillBound` die Bindung unter derselben
  Tabellen-Kennung. Ein Fehler an irgendeiner Stelle vor `Commit` geht über
  `return run, err` durch den `defer`-Rollback (`committed == false`) und das
  Schließen des Snapshots (`defer closeSnapshot`, registriert zuerst und damit
  **nach** dem Rollback ausgeführt). Kein Commit bei einem Fehler vor dem Commit:
  M4, M5, M9, M10b, M16, M18 färben je den passenden Test rot.
- **`context.WithoutCancel`:** `Finish` (Z:302), `Rollback` (Z:432) und `Close`
  (Z:421) laufen auf einem vom Abbruch gelösten Kontext; M13–M15 sind rot; die Fakes
  halten `ctx.Err()` je Aufruf (`finishCtxErr`, `rollCtx`, `closeCtx`).
- **`interrupted` gegen `failed`:** entschieden am eigenen Kontext in `conclude`
  (`ctx.Err() != nil` → `Interrupt`, sonst `Fail` mit `classifyError`); ein noch
  `queued` Run bleibt bei beendetem Kontext `queued` und der Aufruf meldet den
  Kontext-Fehler; M11/M12 rot.
- **`Request`:** Vorbedingungen (`Registered`, `Published`), Schätzung
  (`EstimatedRows`: `known`=falsch → `UnknownRowEstimate`, bekannte `0` bleibt
  bekannt), dann `Admit` als letzter Schritt; kein Snapshot.
- **`Execute`:** zuerst erneute Prüfung von Bindung und Publication
  (`publishedTable`, Abweichung → `failed`/`configuration`, ohne Snapshot),
  dann die Schema-Version (`CurrentVersion`, unbekannt → `configuration`), dann
  `Start` und `MarkRunning`, dann `copyBlocks`; ein Wecksignal nur nach einem
  Commit mit mindestens einer Zeile.
- **Fehler des Runs sind run-lokal:** der Use Case hält keinen Heartbeat-, Ack-,
  Store- oder Stream-Port (`Ports`, acht Felder); `TestPortsCarryNoCapturePathPort`
  bindet das über Typnamen (Reichweite wie F-13 des Reviews, unverändert — strukturell
  trägt die Aussage, weil der Use Case keinen Zugang besitzt).

### Urteil zu „Ergebnis kann der Zeile widersprechen"

Der Fall: `Commit` wirkt serverseitig, der Client erhält einen Verbindungsfehler.
`copyBlocks` liefert den Fehler, der `defer` ruft `Rollback` (wirkungslos nach
einem serverseitig gewirkten Commit), `conclude` ruft `Finish(failed)` — vertraglich
ein wirkungsloser Erfolg auf einen `completed`-Run —, und das Ergebnis trägt `failed`
mit Fehlertext, während die Zeile `completed` trägt; der Use Case liest die Zeile
nicht zurück. **Bewertung: LOW (V-1).**

- **Warum nicht höher:** die Daten sind vollständig und atomar committet (kein
  Teilbestand, keine stille Lücke); die dauerhafte Wahrheit ist die Zeile, das
  Ergebnis der Zustand des Aufrufs; Plan §6 führt den Punkt als eigenes Risiko
  („Commit mit unbekanntem Ausgang") samt offener Frage, ob das Ergebnis die Zeile
  zurücklesen soll; die Gegenseite (`Finish` als wirkungsloser Erfolg) steht als
  Port-Vertrag und als DoD-Punkt „Adapter-Pflichten" in
  `slice-backfill-run-store` mit einem Store-Test als Beleg.
- **Warum nicht INFO:** (a) die Richtigkeit hängt an einer Adapter-Pflicht, die
  heute nur ein Plan trägt — verletzt ein Adapter den Vertrag (`Finish(failed)`
  überschreibt `completed`), meldet die Zeile `failed` über committete Daten, und ein
  Betreiber würde neu beantragen und den Bestand doppelt kopieren; (b) das Wecksignal
  entfällt bei diesem Ergebnis (`Execute` sendet nur bei fehlerfreiem Ende), obwohl
  Daten sichtbar sind — folgenlos, solange Leser pollen, aber ein Unterschied zu
  „completed"; (c) der Aufrufer (`slice-backfill-sql-administration`) muss
  `result.Run.Status` als Zustand des Aufrufs lesen und nicht als Zeilenzustand —
  diese Übergabe steht **nicht** unter den drei gemeldeten Übergaben im Suchlauf-Feld
  des Plans. Kein Blocker; Ausgang bei Closure, Übergabe an den Planner (siehe
  Übergabe).

## 7. Entscheidungs-Konformität

| Entscheidung | Festlegung | Beleg | Verdikt |
|---|---|---|---|
| [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 3 | Position `X` an jedem Block, `committed_at` | `blockBuilder.position` aus `snapshot.Offset()`; `TestExecuteWritesAllBlocksInOneTransaction` (Position je Block) | konform |
| [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 4 | eine Schreibtransaktion, ein Commit am Ende mit der Run-Zeile `completed`; jeder Block liest den Ausschlussstand neu; vor dem Commit Bindung und Ausschlussstand; jede Abweichung rollt zurück (`failed`, Grund) | §6; M4/M5/M9/M10b/M16 rot; Grenze der Prüfung im Doc von `copyBlocks` und in Plan §6 benannt | konform |
| [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 5 | Vorbedingungen (`Registered`/`Published`), Klassen `permission`/`configuration`/`storage`/`transient`/`replication`, run-lokal, ein Wecksignal je Tabelle über den `ChangeNotificationPort` (best effort, [`ADR-0055`](../plan/adr/0055-nats-change-notification-wecksignal.md)) | `classifyError` (fünf Snapshot-Sentinels je ein Fall), `TestExecuteNotifyIsBestEffort`, `TestExecuteEmptyTable` (kein Wecksignal); `internal` als Rückfall liegt in der geschlossenen Menge von [`SPEC-008`](../../spec/pflichtenheft.md)/[`ADR-0023`](../plan/adr/0023-fehlerklassifikation.md) (sieben Klassen); `schema` wird nicht vergeben, im Plan benannt und an den Architect gemeldet | konform |
| [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 6 | Kennung `0bf-<run>-<8 Stellen>`, Sequenz `1…B`, `change_id` = `<Transaktions-ID>-<Sequenz>` mit derselben Bildungsregel wie der WAL-Pfad | `BackfillTransactionID`, `ChangeIDFor`; M19, M20 rot; Mapper-Tests unverändert grün | konform |
| [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 7 | `committed_at` = Snapshot-Zeitpunkt über `ClockPort` ([`ADR-0040`](../plan/adr/0040-clockport.md)) | Uhr nach `OpenSnapshot` gelesen (`blockBuilder.at`); `TestExecuteCommittedAtIsSnapshotTime` | konform |
| [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 1 | drei Ports; Annahme-Port `Admit` als eigener Fähigkeits-Port; Run-Zustands-Port ohne Anlage; `Request` endet mit `Admit` | `TestPortSchnitt`, `TestRequestAdmitsLastAfterPreconditionsAndEstimate`; Port-Dateien gelesen | konform |
| [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 2 Punkt 4 | erneute Prüfung von Bindung und Publication vor dem Öffnen des Slots; Abweichung → `failed`/`configuration`, kein Slot, keine Kopie; `queued` überlebt den Neustart | `TestExecuteRechecksPreconditions` (`openCalls == 0`), `TestExecuteContextEndedBeforeStartLeavesQueued`; M12 rot | konform (Aufnahme beim Start und Wecksignal liegen bei `sql-administration`, dort nicht Gegenstand) |
| [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 3 | Schätzung „unbekannt" ≠ 0; keine Schwelle im Code | `RowEstimate` (Nullwert unbekannt), `TestRequestCarriesEstimate`; Warn-Felder bleiben `false` | konform |
| [`ADR-0115`](../plan/adr/0115-backfill-spaltenwerte-text-ergebnisformat.md) | Bild über die eine gemeinsame Funktion, Werte als Text (`*string`, `nil` = NULL) | `blockBuilder.build` ruft `model.BuildRowImage` (keine zweite Bild-Erzeugung); `git diff 455bcdef..HEAD -- internal/domain/model/rowimage*` leer | konform |
| [`ADR-0034`](../plan/adr/0034-ports-nach-faehigkeiten.md), [`ADR-0028`](../plan/adr/0028-inbound-use-cases.md), [`ADR-0027`](../plan/adr/0027-capture-application-service.md) | Ports nach Fähigkeit, Inbound Port, Application Service | `BackfillTableUseCase`, `BackfillTableService`; `Ports.Schemas` (`SchemaStorePort`) ist der bestehende Fähigkeits-Port, kein neuer Schnitt | konform |
| [`ADR-0002`](../plan/adr/0002-abhaengigkeitsrichtung.md)/[`ADR-0003`](../plan/adr/0003-physische-modulgrenzen.md) | Hexagon-Kanten | `make a-check` EXIT=0, 0 Befunde | konform |
| [`SPEC-029`](../../spec/pflichtenheft.md) | Feldform des Run-Zustands | `BackfillRun`: `started_at` leer bis `running` (auch bei `queued` → `failed`), `finished_at` je Endzustand (M1–M3), `error_message` = `<Klasse>: <Ursache>` nur bei `failed`, `interrupted` ohne Text, `estimated_rows` unbekannt ≠ 0, beide Warn-Felder `false` | konform |
| [`LH-FA-CAP-006.a`](../../spec/pflichtenheft.md) | keine unbegrenzte RAM-Haltung | Streaming-Reihenfolge, ein Block je Durchlauf | konform (Zeilen, nicht Bytes — im Plan so benannt) |

## 8. Harte Regeln

- **[`AGENTS.md`](../../AGENTS.md) §3.3 — Moves rein:** `git show --stat -M d2e24986 3df8565d`:
  `{open => next}/… | 0` und `{next => in-progress}/… | 0`, je 0 Einfügungen und
  0 Löschungen. **Verantwortlich-Vermerk (`643582b0`):** ein eigener Commit nach dem
  Move, ein Ein-Zeilen-Diff im Kopf (`— (noch nicht priorisiert)` → `Implementer-Agent,
  2026-09-24`). Urteil mit dem Regelwerk der Baseline (`v6.9.0` ·
  `regelwerk/modul-05-planning-harness.md` §Lifecycle als State Machine, §Trigger je
  Lifecycle-Übergang): **§3.3 ist eingehalten** (Move rein, Inhaltsänderung im eigenen
  Commit). **Streng gelesen** verlangt der Trigger `open→next` das gesetzte
  `Verantwortlich:` — hier stand es beim Move noch auf „—" und wurde erst in
  `in-progress/`, vor dem ersten Code-Commit (`510b247e`), gesetzt. Das Regelwerk
  nennt das Feld ausdrücklich „Deklaration", ein Sensor prüft es nicht; die Folge ist
  gering (V-5, INFO), eine Nachbesserung gibt es nicht.
- **§3.5 (ADR-Immutabilität):** `make doc-immutable` EXIT=0; im Diff keine ADR-Datei.
- **§3.6:** kein Gate gelockert, keine Schwelle berührt; `THRESHOLD` und
  `DB_COVERAGE_THRESHOLD` nicht im Diff.
- **§3.2:** `git diff 455bcdef..HEAD -U0 -- '*.go' | grep -c -i nolint` = 0.
- **§3.7 (diff-skopierter Kandidatenlauf):** über alle `+`-Zeilen der `*.go`-Dateien
  und der geänderten Plan-Dateien (Chronik-/Konjunktiv-Vokabular: `früher`, `zuvor`,
  `jetzt`, `nicht mehr`, `Fixrunde`, `Review`, `Befund`, `F-<n>`, `slice-`, `welle-`,
  `wäre`, `würde`, `hätte`, `ersetzt`, `wurde`, `statt`, `vorher`, `bisher`,
  `künftig`, `später`, `nur noch`, `sonst`, `seit `): in `.go` fünf Treffer, alle
  unschädlich (die Fehlermeldung „die Bindung … besteht nicht mehr", der Fall-Name
  „Einschluss statt Ausschluss", zwei Substring-Treffer in „gebaut wurden"/„sonst ist
  er leer" als Indikativ-Beschreibung des Zustands, eine Testmeldung „der Empfänger
  wurde verändert"); in den Plan-Dateien 0 Treffer. Kommentare zu Verhalten
  anderer Slices tragen einen Rang-Zeiger (`ADR-0113 Festlegung 2` in `conclude`,
  Anker im Kopf-Kommentar der Inbound-Datei); der Fixrunden-Nachzug des Kommentars
  in `conclude` ist damit geschlossen (F-12).
- **§3.11:** `docs-check` 0 Befunde (`hostpaths`); ein Suchlauf über alle `+`-Zeilen
  des Diffs auf die Wurzel-Segmente der Präfixliste `hostpaths.prefixes` und auf
  Laufwerks-Muster: 0 Treffer;
  dieser Report trägt keinen host-lokalen Pfad.
- **§3.12:** Zahlen im Suchlauf-Feld nachgemessen (§5): alle genannten Zahlen und
  Zeilennummern stimmen; Coverage-Zahlen stehen im Plan nicht.
- **§3.13:** das Feld trägt Gefundenes und Nichtgefundenes je Träger an beiden
  Ständen (Parent `643582b0`, Stand `b51a5494`, dazu `882bdccd`); fremde Träger
  sind gemeldet statt mitgeändert, mit einer benannten Ausnahme: der minimale
  Nachzug im Plan `slice-backfill-run-store` (Zusatz eines DoD-Punkts und einer
  §3-Zeile) ist im Plan als „fremder Plan, minimal" benannt und ändert die Absicht
  jenes Plans nicht (Diff: +16, −0).
- **Handbuch unberührt:** `git diff --name-only 455bcdef -- docs/user/benutzerhandbuch.md`
  ist leer.
- **Commit-Traceability:** zehn Commits, jede Message nennt `LH-*` oder `ADR-*`,
  keine `SPEC-`/`ARC-`-Kennung im Betreff (§1).
- **§3.1:** der Diff trägt keine Host-Toolchain; alle Läufe liefen über `make` im
  gepinnten Toolchain-Image.

## 9. Befunde

Kein V-Befund blockiert. V-1 ist LOW, V-2…V-5 sind INFO.

- **V-1 (LOW) — Ergebnis von `Execute` kann der Zeile widersprechen.** Siehe §6
  „Urteil". `internal/application/usecase/backfill/service.go:302-305` (`conclude`,
  `Finish(failed)` nach einem Commit mit unbekanntem Ausgang) gegen
  `internal/application/port/outbound/backfillrun.go` (`Finish` als wirkungsloser
  Erfolg auf einen beendeten Run). Plan §6 führt das Risiko („Commit mit unbekanntem
  Ausgang", Ausgang bei Closure) und den Träger der Adapter-Pflicht
  (`slice-backfill-run-store`, DoD „Adapter-Pflichten"). **Nicht geführt:** die
  Übergabe an `slice-backfill-sql-administration` — der Aufrufer liest
  `result.Run.Status` als Zustand des Aufrufs, nicht als Zeilenzustand, und die
  Ausgabe in `cdc.backfill_status` ist die Zeile. Vorschlag an den Planner: einen
  Satz in den Plan von `slice-backfill-sql-administration` (Übergabe, keine
  Mitänderung durch den Implementer dieses Slice) und den Ausgang bei Closure führen
  („entfallen" erst mit dem Store-Test des `run-store`).
- **V-2 (INFO) — Der Escape `\|` in den Befehlszellen des Suchlauf-Felds.** Die
  Befehle stehen in Tabellenzellen mit `\|` (Tabellen-Escape); als roher Quelltext
  kopiert liefert `git grep -n -E 'Admit\|Annahme-Port\|Run-Zustands\|Schreiber' 643582b0 -- 'docs/plan/planning/open/slice-backfill-*'`
  **0** Zeilen (`\|` ist in `-E` ein wörtliches Pipe-Zeichen); mit `|` (so, wie eine
  Markdown-Darstellung die Zelle zeigt) **11**. Ebenso Zeile 5 (`usecase/backfill\|Ausschlussstand`,
  0 gegen 6). Die Fixrunde hat den zweiten Teil von F-6 des Reviews behoben
  (Abgrenzung der Zeilen, `e2e`-Treffer werden genannt, der Transformations-Plan
  hat eine eigene Zeile); den ersten Teil (Kopierbarkeit des Rohtexts) nicht.
  Dieselbe Schreibweise trägt der Vorgänger-Slice `backfill-snapshot-reader`, dessen
  Verifikation sie ebenso auflöste — Repo-Konvention, kein Neuverstoß; die Zahlen
  stimmen. Kein Nachzug im Slice nötig; ein Satz „`\|` steht für `|`" im Feld würde
  die Kopierbarkeit herstellen.
- **V-3 (INFO) — Fixrunde ohne Nachprüfungs-Review.** Die DoD-Zeile „Review
  durchgeführt" ist durch den Report erfüllt; F-1…F-7 und F-12 habe ich selbst
  nachgemessen: F-1 (Lesefehler) M4/M5; F-2 (`finished_at`) M1–M3; F-3 (Fortschritt,
  Leer-Endzustand) M6–M8; F-4 (Zeitbegrenzung) Port-Doku in drei Dateien plus
  DoD-Punkt in `run-store`; F-5 (`Finish` auf beendeten Run) Doc-Kommentar; F-6 zur
  Hälfte (V-2); F-7 (Grenze der Prüfung) Doc von `copyBlocks` und Plan §6; F-12
  (Kommentar in `conclude`) Rang-Zeiger. F-8/F-9 (Architect) stehen als „An Architect
  gemeldet" in Plan §6, F-10 ist nicht als Änderung beauftragt (Ausschluss-Lesung je
  Block, Kosten nicht gemessen), F-11 (`Verantwortlich` nach dem Move) ist nicht
  mehr behebbar (V-5), F-13 (Ordnungs-Beleg) trägt `run-store` (§3).
- **V-4 (INFO) — Coverage-Zahl lauf-gebunden.** `make gates` druckte 84.50 %,
  `make coverage-gate` unmittelbar danach 84.40 %, der Reviewer 84.40 % — dieselbe
  Codebasis; die Zahl schwankt im Zehntel-Bereich zwischen Läufen. Der Plan nennt
  keine Coverage-Zahl; die Schwelle (80 %) ist mit Abstand erfüllt.
- **V-5 (INFO) — `Verantwortlich:` nach dem Move gesetzt.** Siehe §8 (§3.3): kein
  Verstoß gegen die Hard Rule, eine Abweichung vom Trigger-Wortlaut von
  `open→next` im Regelwerk, geheilt vor dem ersten Code-Commit; der Vorgänger-Slice
  setzte das Feld in `next/`. Kein Nachzug.

## 10. Ergebnis

| Prüfpunkt | Ergebnis |
|---|---|
| DoD §2 — `[x]`-Zeilen (9) | **9 von 9 erfüllt**, je mit eigenem Beleg (Zeile 6 mit V-3) |
| DoD §2 — übrige `[ ]`-Zeilen (4) | **4 von 4 korrekt offen** (Closure-Notiz, Beobachtungs-Register, Risiko-Ausgänge, Paarungen) |
| Plan-vs-Code-Diff | **deckungsgleich**: jede Plan-Zeile im Diff vertreten, vier benannte Abweichungen mit Begründung, nichts Ungeplantes |
| Fixrunde `882bdccd..b51a5494` | **kein Produktionsverhalten geändert**: `service.go` nur Kommentare, Ports nur Doc-Kommentare, Rest Tests/Fakes |
| Acht Review-Mutationen | **8 von 8 rot**, darunter beide Lesefehler-Mutationen und `finished_at` bei `interrupted` |
| Weitere Mutationen | 12 von 12 rot (Vergleich je Block/vor dem Commit, `WithoutCancel` an drei Stellen, `interrupted`, `queued`, Bindung, Rollback, Commit-Fehler, `ChangeIDFor`, Präfix) |
| ADR-/Spec-Konformität | **konform** (siehe §7: die ADRs zu Position, Atomarität, Klassen, Ordnung, Retention, Rollenschnitt, Ports, Hexagon-Kanten sowie die Spec-Stellen zu Run-Zustand und Fehlerklassen) |
| Harte Regeln (§3.1, §3.2, §3.3, §3.5, §3.6, §3.7, §3.9, §3.11, §3.12, §3.13, Handbuch) | **erfüllt** (V-2, V-5 INFO) |
| Commit-Traceability, MR-Immutabilität | **grün** (`RANGE=455bcdef..HEAD`) |
| Gates | **`make gates` EXIT=0**, `make test` EXIT=0, `make a-check` EXIT=0, `make coverage-gate` EXIT=0 (`84.40%`; `make gates`: `84.50%`) |
| DB-Sensoren | nicht nötig (kein Store-Paket im Diff, §1) |

## Verdikt

**Bestätigt.** Die DoD-Behauptung des Implementers ist durch eigene Messung belegt:
der Use Case hält die Fail-closed-Zusage an ihrer Eingabe (jeder Lesefehler, jede
Abweichung, jede fehlende Bindung führt zu Rollback, `failed` und keinem Commit — je
Fall eine rot färbende Mutation), die Endzustände tragen `finished_at` aus der Uhr,
Snapshot-Schließen, Rollback und `Finish` laufen auf einem abgelösten Kontext,
`interrupted` gegen `failed` und `queued` bleibt `queued` entscheidet der eigene
Kontext, die Kennungs- und Ordnungsregeln tragen (Mapper unverändert grün). Die
Fixrunde hat kein Produktionsverhalten geändert. Kein Befund blockiert; V-1 (LOW)
und V-2…V-5 (INFO) gehen ohne Änderung im Slice an den Planner.

**Übergabe:** Verifier → Planner. Offen bleibt die Closure: Closure-Notiz mit
Lerneintrag (die Finding-Klassen des Reviews, darunter „Zusage ohne Bindung an ihre
Eingabeseite" — Fakes, die nur dauerhaft scheitern, sind die Ursache — und „Beleg
trägt seinen Satz nicht"), Ausgänge der zehn §6-Risiken (V-1: die Übergabe an
`slice-backfill-sql-administration` ergänzen, „Ergebnis gegen Zeile" erst mit dem
Store-Test des `run-store` schließen), Beobachtungs-Register, Welle-Paarungen. Die
beiden an den Architect gemeldeten Punkte (`schema_version` der Backfill-Changes vor
dem Run, Klasse `schema` im Backfill-Pfad der Transformationen) bleiben dort. Danach
der reine `git mv` nach `done/` ([`AGENTS.md`](../../AGENTS.md) §3.3, Fall 2).
