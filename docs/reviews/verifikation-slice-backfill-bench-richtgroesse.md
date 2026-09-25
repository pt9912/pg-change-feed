# Verifikations-Report: slice-backfill-bench-richtgroesse — 2026-09-25

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich + Entscheidungs-Konformität +
Plan-vs-Code-Diff + Gates. Review-Artefakt des Reviewers:
[`review-slice-backfill-bench-richtgroesse.md`](review-slice-backfill-bench-richtgroesse.md);
Formvorbild dieses Reports: [`verifikation-slice-backfill-e2e.md`](verifikation-slice-backfill-e2e.md).
Die vendored Baseline trägt kein eigenes Verifikations-Template
(`.harness/baseline/v6.9.0/templates/docs/reviews/` enthält nur `review-report.template.md`);
der Report folgt dem Formvorbild. Einen Verifier-Skill trägt das Repo nicht
(`.harness/skills/` führt Reviewer und Closure-Note-Reviewer).

**Gegenstand:** Slice-Plan `slice-backfill-bench-richtgroesse` (Welle `welle-backfill-bestand`),
Stand `HEAD` = `c97273b3`, Diff-Range `933ab054..HEAD`, 11 Commits. Slice-Inhalt: drei
Lifecycle-Commits (`82eba253`, `e736cf4a`, `6f9b391c`), Warn-Auswertung (`f69c1b3b`), Bench-Skript
(`51e43faf`), Handbuch/Plan-Nachzug (`673e3836`), DoD-Haken (`b98b5af7`), Fixrunde (`c97273b3`).
Nicht Slice-Inhalt, im Range enthalten: Architect-Verdikt samt
[`ADR-0120`](../plan/adr/0120-capture-slot-leerlauf-bestaetigung.md) (`a90555b6`), Plan des
Folge-Slices (`adb4b66b`), Review-Report (`50aa8da8`). Dieser Lauf ändert weder Code noch Plan,
Spec oder Doku; er schreibt nur diesen Report. Alle Mutationen liefen im Arbeitsbaum (Ausgabe in
eine Scratch-Datei, `cp`, kein `sed -i`) und sind je Lauf per `git checkout` zurückgenommen;
`tools/schema/plan.yaml`/`down.sql` (Nebeneffekt der Läufe mit Schema-Rollout) sind nach jedem
Lauf ebenfalls zurückgenommen (`git status --short` zeigt zuletzt nur diesen Report).

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — [`AGENTS.md`](../../AGENTS.md) §3.9)

Jeder Lauf schrieb in eine Log-Datei; der Exit-Code wurde im selben Aufruf gesondert gesichert, die
Logs danach gelesen. Stand aller Läufe ohne Mutation: `HEAD` = `c97273b3`, Arbeitsbaum sauber.

| Sensor | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `make gates` | **EXIT=0** | `baseline-verify: v6.9.0 OK — 54 Dateien` · `db-package-lists-check: OK — … dieselben 4 Pakete` · `coverage-gate: OK — Coverage 83.20% erfüllt Schwelle 80%` · `d-check: 1096 Datei(en) geprüft, 0 Befund(e)` · `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD"` · `generated-sync: OK` · a-check `gesamt: 0 Befund(e)` |
| `make test` (Race-Detector) | **EXIT=0** | 42 Zeilen `ok`, 0 Zeilen `FAIL`; `ok … internal/application/usecase/backfill 1.029s` |
| `make test-store` | **EXIT=0** | `ok … postgresstorage 7.269s`; `DB-Adapter-Coverage: 82.13% (gedeckt 850 von 1035 Statements; Profile gemergt: store,replication)`; `db-coverage: OK — … erfuellt Schwelle 70%` |
| `make commit-traceability RANGE=933ab054..HEAD` | **EXIT=0** | `OK — 11 Commit(s) in "933ab054..HEAD", Betreffs ohne Struktur-ID` |
| `make doc-commits RANGE=933ab054..HEAD` | **EXIT=0** | `d-check: 1097 Datei(en) geprüft, 0 Befund(e)` |
| `make doc-immutable RANGE=933ab054..HEAD` | **EXIT=0** | `d-check: 1097 Datei(en) geprüft, 0 Befund(e)` |
| `bash tools/bench-backfill.sh` (einmal, `free -m` vor dem Start 20.6 GB verfügbar) | **EXIT=0** | Lauf `20260925T015600Z`, Zeilen in §3 |
| `bash tools/bench-source-impact.sh` (einmal, einzeln — `make bench` als Ganzes nicht gefahren) | **EXIT=1** des Skripts | `Ergebnis (LH-QA-PER-001) — ohne CDC 16194 ms (Median von 5 Läufen), mit CDC 31089 ms (Median von 5 Läufen), Differenz 14895 ms (92.0%)`; `SCHWELLE ÜBERSCHRITTEN (LH-QA-PER-001, SPEC-025) — Overhead 92.0% liegt über der zulässigen 35%-Schwelle` |
| Store-Test der View, unmutiert, einzeln (Scratch-Kopie des Store-Runners, nur `postgresstorage`, `-v`) | Test grün | `--- PASS: TestBackfillStatusViewShowsTheLatestRunPerTable (0.05s)` — kein `SKIP` |
| Mutationen (§4) | drei Läufe von `make test`, ein Store-Lauf | siehe §4 |

Hygiene: dangling Volumes (`docker volume ls -q -f dangling=true | wc -l`) vor den Läufen **34**,
nach allen Läufen **34**; kein Container und kein Netz `pgc-bench-*` nach den Bench-Läufen
(Aufräumen über die `EXIT`-Falle der Skripte); kein `prune`. Ein Container `upbeat_wright`
(`bats/bats`) lief zeitweise auf dem Host; er stammt nicht aus meinen Läufen und ist unberührt.
Speicher vor den schweren Läufen (`free -m`, verfügbar): 20.2 bis 20.8 GB.

## 2. DoD — Verdikt je Zeile (§2 des Plans)

`[x]` sind acht Zeilen (Nr. 1–8), `[ ]` fünf (Nr. 9–13), zusammen 13 — gezählt am Plan
(`grep -c '^- \[x\]'` 8, `grep -c '^- \[ \]'` 5).

| # | DoD-Zeile | Verdikt | Realer Beleg |
|---|---|---|---|
| 1 | Bench: `tools/bench-backfill.sh` in `make bench` verdrahtet (zuletzt), druckt Kopierdauer je Tabellengröße und Abweichung der Schätzung; jede Zahl mit Größe und Lauf; Beleg auf den erreichbaren Beleg umformuliert | **erfüllt in der umformulierten Form, mit V-1 (LOW) und V-2 (INFO)** — das Original-Kriterium „ein realer `make bench`-Lauf mit Exit 0“ ist **nicht erfüllt**, und der Plan sagt das selbst | (a) Skript einzeln Exit 0: mein Lauf `20260925T015600Z` (§3), dazu der Review-Lauf `20260925T011036Z` (gedruckt im Review-Report); jede gedruckte Zeile trägt Lauf-Kennung im Präfix, Stufe und Zeilenzahl (`bench-backfill[20260925T015600Z]: Stufe 10000, Lauf 1/3 — Kopierdauer 1091 ms …`). (b) `make bench` endet am vorhandenen Skript `tools/bench-source-impact.sh`: mein Einzellauf des Skripts 92.0 % gegen 35 %, Exit 1 des Skripts (`make` meldet dafür Exit 2); die Plan-Zahl 91,7 % (16.094/30.848 ms) folgt rechnerisch aus ihren beiden Zeitangaben (30.848/16.094 − 1 = 0,917) und deckt sich mit meiner Messung; Ursache laut Architect-Verdikt: `fdatasync` 2.956 µs je Operation, Host-Last widerlegt, Kontrolle mit `synchronous_commit=off` −6,7 %; `git diff 933ab054..HEAD -- tools/bench-source-impact.sh` verschiebt nur `median_of` in die Bibliothek, `Dockerfile`, `harness/mk`, `spec/` und `.d-check.yml` sind nicht im Diff (keine Schwelle verändert, [`AGENTS.md`](../../AGENTS.md) §3.6). `make bench`-Rezept: `bench-source-impact`, `bench-scaling`, `bench-batch-vs-single`, `bench-backfill` (Makefile) |
| 2 | Toleranz, Richtgröße, Warnungen: je an genau einer Stelle; Auswertung im Use Case, nicht in `bootstrap.Diagnose`; keine Statusänderung, keine Ablehnung; Unit-Tests mit Fake-Uhr und Mutation; Store-Test der View; View-Signatur unverändert | **erfüllt, mit V-3 (LOW)** | `warn.go` trägt `copyDurationToleranceMinutes = 10` und `estimatedRowsGuideline = 4_000_000` je einmal (`git grep`, §6); `bootstrap/backfill.go` liest die Spalten der View und trägt weder Konstante noch Auswertung, `internal/bootstrap` ist nicht im Diff; die einzigen Aufrufer von `warnedEstimatedSize` (`service.go` `Request`, vor `Admit` als letztem Schritt) und `warnedCopyDuration` (`progress`, beide Abschlusspfade, `conclude`) ändern weder `Status` noch Ablauf; Mutationen MA/MB rot (§4), Store-Mutation rot mit gedruckter Meldung (§4); `git diff 933ab054..HEAD -- tools/schema` leer |
| 3 | Benennung: „geschätzt“, „Richtgröße“/„Orientierung“, „Startwert, Setzung ohne Messung“ | **erfüllt** | Suchlauf über die hinzugefügten Zeilen von Handbuch, Harness, Tools, Code: „Grenze“ nur in „Orientierung, keine Grenze“ und „Untergrenze“ (Speicher), kein „Limit“, kein „maximal“ am Begriff; jede Nennung von „Zeilenzahl“ trägt „geschätzt“/„Schätzung“ oder meint eine reale Zahl; `warn.go`: „der Startwert, eine Setzung ohne Messung“; Handbuch: „Startwert, Setzung ohne Messung“ |
| 4 | Träger: `make bench` auf vier Skripte (Makefile, `harness/README.md`, `bench-abdeckung.md`), Handbuch-Historie | **erfüllt** | Makefile-Kommentar „Vier eigenständige Skripte … drei mit Schwelle“, Hilfetext „vier Skripte“, vier Rezeptzeilen; `harness/README.md` Zeile `make bench` „vier eigenständige Bench-Skripte … (drei mit Schwelle, eines ohne)“ samt der WAL-Messung; `docs/user/bench-abdeckung.md` Kopftext „drei … mit Pass/Fail-Schwelle … `tools/bench-backfill.sh` misst ohne Schwelle und trägt keine Zeile“ (Datei trägt weiter genau die drei Zeilen `LH-QA-PER-001`…`003`); Handbuch `Version: 1.54`, Historienzeilen 1.53 und 1.54 |
| 5 | `make gates` grün, Exit gesondert | **erfüllt** | eigener Lauf EXIT=0 (§1) |
| 6 | Review durchgeführt, kein offenes HIGH/MEDIUM nach der Fixrunde | **erfüllt, mit V-4 (INFO)** | Report `review-slice-backfill-bench-richtgroesse.md` liegt vor (0 HIGH, 2 MEDIUM, 4 LOW, 6 INFO); F-1 bis F-6 und F-9 bis F-12 in §5 einzeln nachgemessen, F-7/F-8 ohne Aktion |
| 7 | §3.13-Suchlauf: Feld in §3, Gefundenes und Nichtgefundenes, beide Stände | **erfüllt, mit V-6 (INFO)** | elf Zeilen des Feldes nachgefahren, Build-Kontext und Abdeckungs-Träger über den Diff (§6): alle Fundstellen und Zählungen stimmen bis auf zwei Selbstverweis-Zählungen (V-6) |
| 8 | Doku-Update: Handbuch „Grenzwerte“ nennt Toleranz und Richtgröße mit Ursprung, Host und Lauf; Historienzeile | **erfüllt** | Handbuch §9 „Grenzwerte“: Toleranz „Startwert, Setzung ohne Messung“; Richtgröße „Ursprung (abgeleitet)“ mit Lauf `20260925T000439Z` (als übernommen gekennzeichnet), Rechnung 7.693 × 600 = 4.615.800 → 4.000.000 nachgerechnet; Host Linux 6.8.0-139-generic, Docker 29.8.1, 20 CPU (deckt sich mit der Host-Zeile meines Laufs) |
| 9 | Closure-Notiz mit Lerneintrag | **korrekt offen** | Plan §7 trägt „*(zu tragen bei Closure)*“ |
| 10 | Reconciliation-Register — entfällt | **korrekt offen / entfällt** | Greenfield |
| 11 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | Closure-Pflicht |
| 12 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | Plan §6: alle sieben Zeilen „Ausgang: *(bei Closure)*“ |
| 13 | Die drei Paarungen | **korrekt offen** | hängen an der Closure von `welle-backfill-bestand` |

Kein `[x]` ohne Beleg; kein `[ ]`, das über die Rollen-Sequenz hinaus belegt wäre.

**Zur Formulierung von Nr. 1 (Auftragsfrage „ehrlich oder mehr behauptet als belegt“).** Der Punkt
behauptet nur, was belegt ist: er trennt das Original-Kriterium ausdrücklich als „nicht erreichbar
und gilt nicht als erfüllt“ vom Ersatzbeleg (a) + (b) und nennt für (b) den roten Ausgang samt Zahl,
Schwelle und Ursache. Ich habe beide Teile nachgemessen (§1: Skript einzeln Exit 0; die
Schwellen-Überschreitung des Vorgängerskripts 92,0 %). Nicht behauptet und nicht belegt ist ein
grüner `make bench`; ich habe `make bench` als Ganzes auftragsgemäß nicht gefahren (der Beleg für
sein Ende ist das Verdikt plus meine zwei Einzelläufe, kein eigener Gesamtlauf). Zwei
Einschränkungen stehen in V-1 (wer die Umformulierung setzt) und V-2 (Herkunft einzelner Zahlen).

## 3. Eigener Lauf von `tools/bench-backfill.sh` (Lauf `20260925T015600Z`)

Host laut gedruckter Zeile: „Linux 6.8.0-139-generic, Docker 29.8.1, 20 CPU, 33362599936 Byte RAM,
PostgreSQL-Image postgres:18-alpine“; Default-Stufen 10.000/50.000/200.000 Zeilen, 3 Läufe je
Stufe, Blockgröße B = 1.000 Zeilen, Toleranz 10 min (aus dem Code gelesen). Gedruckt (gekürzt):

- „Stufe 10000 Ergebnis — Kopierdauer 1043–1127 (Median 1091, n=3) ms, Durchsatz 8873–9588 (Median
  9166, n=3) Zeilen/s“; „Stufe 50000 … Durchsatz 9369–9750 (Median 9634)“; „Stufe 200000 Ergebnis —
  Kopierdauer 22437–23921 (Median 23111, n=3) ms, Durchsatz 8361–8914 (Median 8654, n=3) Zeilen/s,
  Feed-Speicher-Spitze 441.0 MiB (Ruhe vor der Stufe 216.4 MiB, 20 s nach dem letzten Lauf 368.6
  MiB) … WAL-Rückstand-Spitze im Run Median 141 MiB (~739 B je Zeile, abgeleitet), vom Slot gehaltenes
  WAL Median 281 MiB“.
- Schätzung: frisch befüllt, mit und ohne Autovacuum `reltuples=-1, geschätzt unbekannt`; Autovacuum
  an, 35 s nach dem Befüllen `Abweichung +0.0%`; Autovacuum aus, nie analysiert, 35 s danach
  `geschätzt unbekannt`; nach `ANALYZE` `+0.0%`; nach `ANALYZE` und 20.000 weiteren Zeilen (+20 %)
  ohne erneutes `ANALYZE` `Abweichung -16.7%`.
- „cdc_capture_lag während des Runs (25 s, Stufe 200000, Live-Last 100/s): 0.089549–0.927493 (Median
  0.429046, n=20) s … Referenz ohne Run: 0.046642–1.007869 (Median 0.395197, n=25) s“.
- „Richtgröße (abgeleitet) — 8654 Zeilen/s (Median, Stufe 200000) × Toleranz 600 s = 5192400 Zeilen,
  abgerundet auf eine Stelle: 5000000 Zeilen“; „WAL-Rückstand-Schwellen (abgeleitet) — bei ~739 B je
  Zeile … Fehlerschwelle (1024 MiB) bei ~1452965 Zeilen“; jede Run-Zeile „Warnung Größe f, Warnung
  Dauer f“ (keine Stufe erreicht die Toleranz, jede geschätzte Zeilenzahl ist „unbekannt“).

Abgleich mit den Handbuch-Zahlen (§3.12): Kopierrate 8.654 (mein Lauf) gegen 7.693/8.933/8.559
(Handbuch: übernommen/Review/Fixrunde) — derselbe Bereich, die Richtgröße nach Rundung ist in drei
von vier Läufen 5.000.000 und einmal 4.000.000, wie das Handbuch sagt (die Konstante trägt den
kleinsten Wert); Schätz-Zeilen identisch zu den Handbuch-Werten; WAL-Rückstand-Spitze 141 MiB, 739
B je Zeile gegen Handbuch 140 MiB, 735 B (übernommen, gleicher Bereich); Live-Wirkung 0,43 gegen
0,40 s (im Run höher), damit weiter „in beiden Richtungen“ über die Läufe (0,48/0,50, 0,59/0,40,
0,456/0,413, 0,43/0,40) und die Aussage „kein Unterschied ableitbar“ ohne Widerspruch. Die Probe 20
s nach dem letzten Run (368,6 MiB) liegt in meinem Lauf **unter** der Spitze im Run (441,0 MiB); das
Handbuch bindet seinen Satz „über der Spitze“ an „die zwei gemessenen Läufe“ (437,0 und 641,7 MiB)
und bleibt damit wahr; die Speicher-Spanne über die vier Läufe reicht von 401,7 bis 467 MiB
(Spitze im Run) und bis 641,7 MiB (höchster gemessener Wert).

## 4. Mutationen der Eingabeseite (dieser Lauf)

Paket `internal/application/usecase/backfill`, Lauf `make test` je Mutation, Rot-Meldung gelesen,
Datei danach per `git checkout` zurückgenommen. Gewählt sind Mutationen, die nicht dieselben sind
wie die des Reviews (M1–M13 dort).

| # | Mutation | Lauf | Ergebnis (gedruckt) |
|---|---|---|---|
| MA | `warnsEstimatedSize`: `rows > estimatedRowsGuideline` → `rows > estimatedRowsGuideline+1` (Grenze nach oben verschoben) | `make test`, **EXIT=2** | `--- FAIL: TestWarnsEstimatedSizeBoundary` („warnsEstimatedSize = false, will true“), `--- FAIL: TestWarnedEstimatedSize`, `--- FAIL: TestRequestWarnsAtEstimatedSizeGuideline` („Warnung Größe: Admit false, Ergebnis false, will true“) — rot |
| MB | `warnsCopyDuration`: `now.Sub(started)` → `started.Sub(now)` (Vorzeichen der Dauer) | `make test`, **EXIT=2** | `--- FAIL: TestWarnsCopyDurationBoundary` („warnsCopyDuration(600000000001 ns) = false, will true“), `TestWarnedCopyDuration`, `TestExecuteWarnsCopyDurationAtTolerance` („Warnung Dauer je Fortschritts-Update = [false false false false], will [false false true true]“), `TestExecuteWarnsCopyDurationOnEmptyTable`, `TestExecuteWarnsCopyDurationOnFailureAndInterrupt` — rot |
| MC | `copyDurationToleranceNanos`: Faktor `* 60` entfernt (die Toleranz wäre 10 Sekunden statt 10 Minuten) | `make test`, **EXIT=0** | **überlebt** — alle Tests lesen die Konstante über `export_test.go` relativ zu sich selbst; siehe V-3 |
| MS | Store: View `cdc.backfill_status`, Spalte `r.warn_estimated_size` → `false AS warn_estimated_size` (`tools/schema/schema.yaml`); Lauf des Store-Runners nur für `postgresstorage` (Scratch-Kopie, `-v -run TestBackfillStatusView`) | **rot** | `backfillstatusview_test.go:67: vwst_e: Schätzung 7000000, Warnungen false/false — erwartet 7000000, true/false` — die Meldung nennt den **Wert** der Schätzung (kein Zeiger, Review F-6). Im echten `make test-store` färbt dieselbe Mutation zuerst `TestDiagnoseReportsTheLatestBackfillRunPerTable` (EXIT=2) rot |

MA und MB färben je mehrere Tests an genau der mutierten Zusage rot (Grenze der Richtgröße, Richtung
der Dauer), an ihrer Eingabeseite (Schätzung, Uhr), nicht an einer Nebenwirkung. MC ist keine
Regression, sondern die belegte Bindungslücke der Zusage „Startwert 10 Minuten“ (V-3).

## 5. Findings des Reviews nachgemessen (nicht dem Fixrunden-Bericht geglaubt)

| Finding | Was ich gemessen habe | Verdikt |
|---|---|---|
| F-1 (MEDIUM) Träger tragen den Stand vor dem Verdikt | Handbuch §4, Punkt „WAL-Rückstand des Capture-Slots“: „Der Rückstand ist **nicht an den Backfill gebunden** …“, Folge-Slice `slice-backfill-slot-leerlauf-bestaetigung`, Abhilfen (Commit auf aktivierter Tabelle, Datei-Feld `wal_retention_error_bytes`); das Feld steht in [`SPEC-013`](../../spec/pflichtenheft.md)/`internal/bootstrap/config_file.go` und im Pflichtenheft als „— (kein Env-Gegenstück)“, deckt die Aussage; Plan §3 „Befunde der Messung“ Punkt 1 (beantwortet durch Verdikt und [`ADR-0120`](../plan/adr/0120-capture-slot-leerlauf-bestaetigung.md), Richtgröße nicht an die WAL-Schwelle gekoppelt) und Punkt 3 (Ursache Flush-Latenz, Verdikt); DoD-Punkt „Bench“ auf den erreichbaren Beleg gesetzt; die Werte des Verdikts (2.956 µs, 33,5 MiB, 3 s) stehen im Verdikt | **behoben** |
| F-2 (MEDIUM) Speicher-„Spitze“ kleiner als ein Wert derselben Messung | Handbuch §9: Spitze im Run **und** Probe 20 s danach (437,0 und 641,7 MiB) mit der Aussage „höchster gemessener Wert der Stufe ist 641,7 MiB, die Spitze im Run ist eine Untergrenze“, „Bemessen Sie den Speicher nicht knapp an der Spitze im Run“; die 641,7 MiB und 416,5 MiB decken sich mit den gedruckten Zeilen des Review-Reports; meine Messung (441,0/368,6) widerspricht dem Satz nicht (§3) | **behoben** |
| F-3 (LOW) Vertrag/README ohne WAL-Messung | `harness/targets/bench-backfill.md` §Ablauf Schritt 2/3/4 und §Grenzen nennen WAL-Rückstand (`confirmed_flush_lsn`), gehaltenes WAL (`restart_lsn`), Live-Commit-Freigabe (≤ 120 s), Zeile „WAL-Rückstand-Schwellen (abgeleitet)“, Speicher 20 s nach dem letzten Run; `harness/README.md` `make bench`-Zeile nennt sie; jede beschriebene Größe erscheint in meiner gedruckten Ausgabe (§3) | **behoben** |
| F-4 (LOW) Einheit `_S` ohne Bindung | `tools/bench-backfill.sh`: `deadline=$((SECONDS + RUN_TIMEOUT_S))`, Abbruch bei `"$SECONDS" -gt "$deadline"` mit Meldung „… nach ${RUN_TIMEOUT_S} s“; Vertrag: „Sekunden (Uhr der Shell)“; meine Läufe endeten regulär, die Abbruchgrenze selbst nicht ausgelöst (Lesung) | **behoben** (Lesung; Abbruchzweig nicht gefahren) |
| F-5 (LOW) zwei Median-Formeln | `range_of` ruft `bench::median_of`; `git grep -n 'NR + 1' HEAD -- tools` ohne Treffer in `bench-backfill.sh`; `bench::median_of` steht einmal in `tools/bench-lib.sh`; meine Läufe drucken plausible Bereiche mit Median (`Kopierdauer 22437–23921 (Median 23111, n=3)`, Median = mittlerer der drei) | **behoben** |
| F-6 (LOW) Zeiger statt Wert | Diff gelesen (`text()`-Hilfsfunktion, `strconv.FormatInt`); MS zeigt die gedruckte Meldung mit dem Wert 7000000 | **behoben**, mutiert bestätigt |
| F-7 (INFO) Konstanten nicht an Träger gebunden | ohne Aktion laut Reviewer; meine Mutation MC verschärft den Befund (V-3) | **bestätigt** |
| F-8 (INFO) äquivalente Mutation `known &&` | ohne Aktion | unverändert |
| F-9 (INFO) Richtgröße reproduziert sich nicht | Handbuch nennt drei Läufe (7.693 → 4.000.000; 8.933 und 8.559 → je 5.000.000), Konstante = kleinster Wert; mein vierter Lauf (8.654 → 5.000.000) bestätigt die Streuung | **im Träger gekennzeichnet** |
| F-10 (INFO) Live-Wirkung Einzellauf | Handbuch: „jeder Vergleich ist ein Einzellauf ohne Wiederholung“, drei Läufe genannt, „kein Unterschied ableitbar“; mein vierter Lauf liegt im selben Band | **im Träger gekennzeichnet** |
| F-11 (INFO) übernommene Läufe nicht auflösbar | Handbuch nennt „übernommen aus dem Lauf-Bericht, im Repository nicht auflösbar“ an jeder Stelle mit den Läufen `20260924T233628Z`/`20260925T000439Z`; Lauf `20260925T012459Z` ist als „(gemessen)“ geführt, seine Zeilen liegen ebenfalls nicht im Repository (V-2) | **gekennzeichnet, mit V-2** |
| F-12 (INFO) Reihenfolge der Skripte | Makefile ruft `bench-backfill.sh` zuletzt, Kommentar nennt den Grund („damit ihr Abbruch keine Schwellen-Prüfung und nicht die Erzeugung von docs/user/bench-abdeckung.md verdeckt“); Plan §3 Zeile `Makefile` und Suchlauf-Zeile tragen ihn | **behoben** |

## 6. Suchlauf-Feld (§3 des Plans), an beiden Ständen nachgefahren

Parent `933ab054` per `git grep … 933ab054`, Fixrunden-Parent `50aa8da8` für die Fixrunden-Zeilen,
Diff-Stand per `git grep … HEAD` bzw. `grep` im Arbeitsbaum. Gedruckte Zahlen:

| Plan-Zeile | Meine Messung | Ergebnis |
|---|---|---|
| Zeile `make bench` / Zählwörter (`drei\|vier\|zwei\|beiden\|sechs\|sieben\|fünf` nahe „bench“) | Parent: `docs/user/bench-abdeckung.md:3`, `harness/README.md:150`, `tools/bench-lib.sh:9` und `:171`; dazu per `drei\|beiden\|Drei`: `Makefile:181` „Drei eigenständige Skripte“, `Makefile:188` „drei Skripte“, `tools/bench-lib.sh:145` „der beiden anderen“, `:160` „letzten der drei“. HEAD: `bench-abdeckung.md:3` „drei … mit Pass/Fail-Schwelle“, `README.md:150` „vier … (drei mit Schwelle, eines ohne)“, `bench-lib.sh:9` „vier“, `:154` „der anderen“, `:169` „drei Schwellen-Skripte“, `:181`, `Makefile:181/182/192` „Vier“/„drei mit Schwelle“/„vier Skripte“. Nicht gefunden (Suche über `README.md`, `AGENTS.md`, `docs/user`, `spec`, `harness`, `.github`, `tools`, `Makefile`, `examples`, `sdks`): kein weiterer Träger, der die Zahl der Bench-Skripte führt (`tools/bench-scaling.sh:100` „drei [`SPEC-014`](../../spec/pflichtenheft.md)-Lastenstufen“ meint einen anderen Gegenstand) | bestätigt |
| Hilfetext/Kommentar des Targets `bench` | Parent-Kommentar `Makefile:181–186`, Hilfetext `:188`, drei Rezeptzeilen; HEAD vier Rezeptzeilen (`git diff` gelesen) | bestätigt |
| Bench-Abdeckungs-Träger | `git diff 933ab054..HEAD -- docs/user/bench-abdeckung.md tools/bench-lib.sh`: nur der erzeugte Kopftext, keine Tabellenzeile, Generator `bench::render_abdeckung` einmal in `tools/bench-lib.sh` | bestätigt |
| Handbuch-Grenzwerte und Zählwörter im Backfill-Absatz | Aufzählungspunkte in „Grenzwerte“: Parent **5**, HEAD **12** (5 + 7 Backfill-Punkte: Toleranz, Richtgröße, gemessene Werte, Speicher, WAL-Rückstand, Live-Wirkung, Schätzung); „gelten drei Betriebs-Vorbedingungen“ (Parent) → „vier“ (HEAD), vier Punkte gezählt | bestätigt |
| Toleranz an genau einer Stelle; Wortwahl | `git grep -n 'copyDurationToleranceMinutes\|estimatedRowsGuideline' HEAD`: Definition je einmal in `warn.go` (Zeilen 13 und 23), sonst Tests über `export_test.go`, das Skript (`bench-backfill.sh:38` liest `copyDurationToleranceMinutes` aus `warn.go`) und der Vertrag; `grep -rn -e '4_000_000' -e '4\.000\.000'` trifft `warn.go:23`, das Handbuch (drei Zeilen einer Angabe) und den Plan | bestätigt |
| Build-Kontext | `git diff --stat 933ab054..HEAD -- Dockerfile .dockerignore` leer; `warn.go` liegt unter `internal/` | bestätigt |
| zweite Eigenschaft „Warn-Kennzeichnungen ohne Auswertung `false`“ | Parent: `benutzerhandbuch.md:524` „In dieser Version setzt keine Auswertung sie: beide bleiben `false`“, `backfillrun.go:64` „bleiben `false`, solange keine“, `tools/schema/schema.yaml:273` und `:408`; HEAD: die zwei ersten Fundstellen entfallen, `schema.yaml:273/408` unverändert (Datei nicht im Diff) | bestätigt |
| Fixrunde: WAL-Rückstand | `git grep -c -i -E 'WAL-Rückstand\|Fehlerschwelle\|wal_retention_(error\|warn)' … ':!docs/plan' ':!docs/reviews'`: `50aa8da8` 11 Dateien (Handbuch 23 Treffer), HEAD 13 Dateien (Handbuch 27, `harness/README.md` 1, `harness/targets/bench-backfill.md` 3) | bestätigt |
| Fixrunde: `RUN_TIMEOUT` | `50aa8da8`: `tools/bench-backfill.sh` 3, `harness/targets/bench-backfill.md` 1; HEAD dieselben zwei Dateien (dazu der Plan, 2 Selbstverweise) | bestätigt |
| Fixrunde: Reihenfolge `bench-batch-vs-single` | `50aa8da8` und HEAD gleiche Verteilung (`Makefile` 1, `harness/README.md` 1, `tools/bench-lib.sh` 1, `docs/user/bench-abdeckung.md` 1, `tools/bench-batch-vs-single.sh` 10, Records unter `done/`/`observations/`), Plan 2 → 4 | bestätigt |
| Fixrunde: Werte im Handbuch | `467 MiB`: `50aa8da8` Plan 1 und Handbuch 1, HEAD Plan **2** und Handbuch 1 (Plan nennt „je 1“, V-6); `7\.693`: `50aa8da8` Handbuch 2, HEAD Handbuch 3 und Plan 1 (wie im Plan); `kein Unterschied`: `50aa8da8` Handbuch 1, HEAD Handbuch 2 und Plan **2** (Plan nennt 1, V-6) | bestätigt bis auf V-6 |

## 7. Träger-Nachzüge

| Träger | Beleg | Verdikt |
|---|---|---|
| Handbuch `Version: 1.54`, Historienzeilen 1.53 und 1.54 | `docs/user/benutzerhandbuch.md` Zeile 3, Zeilen 1772 und 1773; §4 trägt die vierte Betriebs-Vorbedingung und die Deutung beider Kennzeichnungen, „Diagnose ausführen“ deutet sie, §9 die sieben Backfill-Punkte | erfüllt |
| `harness/targets/bench-backfill.md` | neu, Vertrag, Ablauf, Parameter, Grenzen; jede beschriebene Größe erscheint in meiner Ausgabe (§3) | erfüllt |
| `harness/README.md`, Zeile `make bench` | vier Skripte, WAL-Messung genannt; Träger verweist auf den Vertrag; Zeile behält die Angabe „real gemessen: 25,0 %/28,2 %“ für `LH-QA-PER-001` (V-5) | erfüllt, mit V-5 |
| `docs/user/bench-abdeckung.md` | Kopftext nennt drei Schwellen-Skripte und dass das vierte ohne Schwelle keine Zeile trägt; Entscheidung „keine Zeile für [`LH-FA-CAP-009`](../../spec/lastenheft.md)“ im Plan §3 begründet | erfüllt |
| Makefile-Kommentar und Hilfetext | vier Skripte, Reihenfolge mit Begründung | erfüllt |
| `spec/pflichtenheft.md` §3 | nicht im Diff: weder Toleranz noch Richtgröße ist ein Spec-Wert (Plan §1 „Ausdrücklich NICHT“) | erfüllt |

## 8. Entscheidungs-Konformität

| Entscheidung | Festlegung | Beleg | Verdikt |
|---|---|---|---|
| [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 3, Punkt 1 (Toleranz) | Startwert 10 Minuten, benannte Konstante an genau einer Stelle im Use-Case-Paket, Uhr über `ClockPort`, kein Konfigurationsschlüssel, kein §3-Eintrag | `warn.go` `copyDurationToleranceMinutes = 10`; Uhr über `s.ports.Clock` in `service.go`; keine `CDC_*`-Variable und kein Datei-Feld im Diff; `spec/` nicht im Diff; Kopierdauer beginnt bei `StartedAt` (`warnedCopyDuration` lässt `queued` frei) | konform, mit V-3 (Wert ohne Test-Bindung) |
| Festlegung 3, Punkt 2 (Richtgröße aus der Messung) | steht vorher nirgends; Zeilenzahl in der Toleranz, abgerundet, abgeleitet gekennzeichnet, Host und Lauf | Konstante 4.000.000 mit Herleitung im Handbuch (Rate × 600 s, eine Stelle abgerundet, „abgeleitet“, Lauf, Host); die Bench-Zeile „Richtgröße (abgeleitet)“ druckt dieselbe Regel (mein Lauf 5.000.000) | konform |
| Festlegung 3, Punkt 3 (zwei Warnungen; keine Ablehnung, kein Abbruch, keine Statusänderung) | (1) beim Antrag über der Richtgröße, (2) zur Laufzeit je Fortschritts-Update und beim Abschluss | `Request` (Warnung 1, `Admit` als letzter Schritt), `progress`, beide Abschlusspfade und `conclude` (Warnung 2); MA/MB rot; kein Aufrufer verzweigt auf das Ergebnis | konform |
| Festlegung 3, Punkt 4 (Ablage, Sichtbarkeit) | Auswertung an einer Stelle im Use Case, Ergebnis in den Warn-Spalten, View reicht durch, `diagnose` liest die View | `warn.go`; Adapter schreiben `WarnEstimatedSize` (`backfilladmission.go:109`) und `WarnDuration` (`backfillrun.go`, `backfillwriter.go`); `internal/bootstrap` nicht im Diff | konform |
| Festlegung 3, Punkt 5 (unbekannte Schätzung) | Warnung (1) entfällt, „unbekannt“, nie `0` | `warnsEstimatedSize` (`known &&`); mein Lauf: frisch befüllt `reltuples=-1, geschätzt unbekannt`, alle Run-Zeilen „geschätzt unbekannt“, Warnung Größe f; Handbuch §4 nennt den Zusammenhang | konform |
| Festlegung 3, Punkt 6 (Benennung) | „geschätzt“, „Richtgröße“/„Orientierung“, „Startwert, Setzung ohne Messung“ | Plan-DoD Nr. 3 (§2) | konform |
| [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) Festlegung 3 (Warnung, keine Ablehnung), Teilfrage 4 (Ein-Transaktions-Form) | keine Ablehnung wegen der Größe; Kopierdauer der Ein-Transaktions-Form gemessen | keine Stelle im Code verzweigt auf die Größe (`estimatedRowsGuideline` nur in `warn.go`); Bench nennt die Einfügeform (Vertrag) | konform |
| [`ADR-0054`](../plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md) §(b), [`ADR-0104`](../plan/adr/0104-benchmark-schwellen-per-001-002-003.md) | Bench-Infrastruktur, Schwellen der drei `LH-QA-PER`-Kennungen; [`AGENTS.md`](../../AGENTS.md) §3.6 keine Schwelle ohne ADR gesenkt | `SPEC-025` und `tools/bench-source-impact.sh` (bis auf die Verschiebung von `median_of`) unverändert; das vierte Skript trägt ausdrücklich keine Schwelle („Messung ohne Pass/Fail“, `bench::record_row` nicht gerufen) | konform |
| [`ADR-0114`](../plan/adr/0114-schema-rollout-vorlauf-view-signatur.md) | keine View-Signatur-Änderung ohne Lesefenster | `git diff 933ab054..HEAD -- tools/schema` leer | konform |
| [`ADR-0120`](../plan/adr/0120-capture-slot-leerlauf-bestaetigung.md) (nur Kontext) | Ursache des WAL-Rückstands, Folge-Slice | im Handbuch als Kontext und Folge-Slice genannt, keine Umsetzung in diesem Slice (Diff ohne `internal/adapters/driving/replication`) | konform |

## 9. Plan-vs-Code-Diff

Jede Zeile der Plan-Tabelle §3 ist im Diff (`git diff --stat 933ab054..HEAD`, 23 Dateien)
vertreten oder als „prüfen“ ohne Änderung geführt: `tools/bench-backfill.sh` (neu),
`Makefile`, `internal/application/usecase/backfill/{warn.go,export_test.go,warn_test.go,service.go,service_test.go}`,
`docs/user/benutzerhandbuch.md`, `harness/README.md`, `harness/targets/bench-backfill.md` (neu),
`docs/user/bench-abdeckung.md`, `tools/bench-lib.sh`, `tools/bench-source-impact.sh`,
`internal/domain/model/backfillrun.go` (Kommentar), `internal/adapters/driven/postgresstorage/backfillstatusview_test.go`.
„Prüfen, nicht ändern“ bestätigt: Run-Zustands-Port, Annahme-Port und Adapter, `internal/bootstrap`
(`backfill.go`, `diagnose_test.go`), `tools/schema/schema.yaml` sind nicht im Diff. Im Diff, aber nicht
Slice-Inhalt: `ADR-0120` samt Index-Zeile, Architect-Verdikt, Plan des Folge-Slices, Welle und Roadmap
(Planner-Nachzug), Review-Report. Ungeplant im Slice-Inhalt: nichts. Lifecycle-Commits `82eba253` und
`6f9b391c` sind reine `git mv` (`0 insertions, 0 deletions`, [`AGENTS.md`](../../AGENTS.md) §3.3);
`Verantwortlich` steht in `e736cf4a` als eigener Commit zwischen den Moves.

## 10. Harte Regeln

- **§3.1** — Läufe über `make` und gepinnte Images; die Bench-Skripte sind repo-eigene Host-Skripte
  (`bash`, `docker`) wie ihre Schwester; keine Host-Toolchain, kein `sed -i`, kein Host-Python.
- **§3.2** — `git diff 933ab054..HEAD` über `*.go`, `*.sh`, `*.sql`, `*.yaml`, `+`-Zeilen: 0 × `nolint`.
- **§3.3** — siehe §9.
- **§3.5** — `make doc-immutable RANGE=933ab054..HEAD` EXIT=0; im ADR-Verzeichnis nur die neue ADR und
  die Index-Zeile im Diff.
- **§3.6** — `Dockerfile`, `harness/mk`, `.a-check.yml`, `.d-check.yml`, `spec/` nicht im Diff;
  keine Schwelle verändert.
- **§3.9** — Exit-Codes nie durch eine Pipe gemessen (§1); Gate-Läufe und Folgehandlungen getrennt.
- **§3.10** — der Diff berührt keinen Workflow (`git diff --stat -- .github` leer).
- **§3.11** — `make docs-check` (in `make gates`) 0 Befunde; `git diff` `+`-Zeilen ohne host-lokalen
  Pfad (0 Treffer); dieser Report trägt keinen.
- **§3.12** — jede Zahl dieses Reports ist **gemessen** (Lauf oder Befehl genannt), **abgeleitet**
  (gekennzeichnet) oder **übernommen** (gekennzeichnet: Handbuch- und Plan-Zahlen aus nicht
  auflösbaren Läufen).
- **§3.13** — Suchlauf-Feld an beiden Ständen nachgefahren (§6).
- **Handbuch-Pflicht** — `Version: 1.52`/`1.53` → `1.54` mit Historienzeilen im Diff.
- **Commit-Traceability** — 11 Commits, jede Message nennt eine Kennung, keine Struktur-Kennung im
  Betreff (§1).

## 11. Befunde

Kein V-Befund blockiert. V-1 und V-3 sind LOW, V-2 und V-4 bis V-7 INFO.

- **V-1 (LOW) — Die Umformulierung des DoD-Punkts „Bench“ ist Planner-Sache; das Kästchen gilt dem
  Ersatzbeleg.** Das Architect-Verdikt weist die Umformulierung dem Planner zu („Feststellung des
  Planners bei der Closure“). Im Repository steht sie in der Fixrunde des Slice (`c97273b3`), ohne
  dass die Rolle erkennbar ist. Inhaltlich trägt sie (§2 Nr. 1). Zusätzlich bindet Plan §5
  (Closure-Trigger) den Abschluss an „`make bench` mit Exit 2 an `LH-QA-PER-001`“: das ist eine
  Eigenschaft dieses Hosts (Verdikt: Flush-Latenz) und keine des Slice. Aktion: der Planner ratifiziert
  bei der Closure ausdrücklich, dass das Original-Kriterium entfällt und (a) + (b) an seine Stelle
  treten, und trägt die Host-Abhängigkeit in die Closure-Notiz.
- **V-2 (INFO) — Als „gemessen“ geführte Fixrunden-Zahlen tragen Lauf-Kennung und keine
  auflösbare gedruckte Zeile.** Lauf `20260925T012459Z` (Handbuch, Plan §3: 8.559 Zeilen/s,
  401,7/437,0 MiB, 0,456/0,413 s) und die 91,7 % im DoD-Punkt „Bench“ (16.094/30.848 ms) stehen ohne
  die gedruckte Zeile im Repository; der Review-Lauf `20260925T011036Z` ist über den Review-Report
  auflösbar. Meine Nachmessung (§3, 92,0 % in §1) liegt im selben Band und stützt die Zahlen, ersetzt
  aber ihren Anker nicht. Aktion: Closure-Notiz nennt die Zeilen oder kennzeichnet die Zahlen als
  übernommen.
- **V-3 (LOW) — Die Zusage „Startwert 10 Minuten“ hat keine maschinelle Bindung, auch nicht an die
  Größenordnung.** MC (Faktor `* 60` entfernt: 10 Sekunden statt 10 Minuten) lässt `make test` grün,
  weil jeder Test die Konstante relativ zu sich selbst prüft; das Bench-Skript liest
  `copyDurationToleranceMinutes` und nicht die Nanosekunden-Umrechnung. Der Review-Befund F-7 (10 → 11)
  nennt nur die Ebene des Zahlenwerts. Die einzige Bindung ist die Prosa des Handbuchs und der
  Suchlauf des Plans; ein Test auf den Wert stünde im Spannungsfeld zur Regel „an genau einer
  Stelle“ ([`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md)). Aktion:
  Planner entscheidet, ob die Umrechnung `Minuten → Nanosekunden` (nicht der Wert) einen eigenen Test
  bekommt; bis dahin bleibt es eine benannte Lücke.
- **V-4 (INFO) — Die Fixrunde hat keinen eigenen Review-Report.** Die DoD-Zeile „Review
  durchgeführt“ ist durch den Report des Reviewers erfüllt; die Fixrunde habe ich in §5 selbst
  nachgemessen (F-6 durch Mutation, F-2/F-3/F-4/F-5/F-12 an Trägern und Ausgabe). Eine
  Reviewer-Durchsicht der Fixrunde wäre ein eigener Schritt.
- **V-5 (INFO) — Die `make bench`-Zeile in `harness/README.md` nennt „real gemessen:
  25,0 %/28,2 %“ ohne Host und Datum, während dieser Host bei 87,5 bis 95,8 % liegt.** Die Zeile
  stammt aus einem früheren Slice und ist durch diesen nicht falsch geworden; ein Leser schließt aus
  ihr auf ein grünes `LH-QA-PER-001`. Das Verdikt benennt den Re-Evaluierungs-Trigger für
  [`ADR-0104`](../plan/adr/0104-benchmark-schwellen-per-001-002-003.md) (zweite Umgebung), ändert die
  Zeile aber nicht. Zur Kenntnis des Planners; kein Umfang dieses Slice.
- **V-6 (INFO) — Zwei Zählungen des Suchlauf-Feldes nennen „je 1“, der Plan trägt selbst 2
  Treffer; die 95,8 % sind dem Parent zugeschrieben, das Verdikt führt alle drei Werte als Läufe des
  Implementers.** `467 MiB` im Plan: 1 vor der Fixrunde, 2 danach (die Suchbefehl-Zeile nennt das
  Muster selbst); `kein Unterschied` im Plan: 2 gegen die genannte 1. Die Einordnung „95,8 % am Parent
  `933ab054`“ im DoD-Punkt „Bench“ stimmt mit der Zeile des Verdikts („am Parent-Stand … 87,5 % /
  93,9 % / 95,8 %“) nur, wenn alle drei Werte am Parent liegen. Kein Einfluss auf ein Verdikt; Zahlen
  ohne Folgen für einen Träger.
- **V-7 (INFO) — Eine gesetzte Warnung ist in keinem realen Lauf gesehen.** Alle Run-Zeilen der
  Bench-Läufe (meiner und der des Reviews) tragen „Warnung Größe f, Warnung Dauer f“, weil keine
  Stufe die Toleranz erreicht und jede frisch befüllte Tabelle „unbekannt“ trägt. Die Setz-Pfade
  belegen Unit-Tests (MA/MB rot) und der View-Test (MS rot); ein Ende-zu-Ende-Beleg einer gesetzten
  Warnung ist nach Plan nicht verlangt. Restrisiko benannt, kein Befund gegen die DoD.

## 12. Grenzen dieses Laufs

- `make bench` als Ganzes nicht gefahren (Auftrag; bekannt rot am Messhost, Verdikt-Beleg gilt); die
  Aussage „`make bench` endet mit Exit 2“ stützt sich auf das Verdikt und meine zwei Einzelläufe
  (`bench-source-impact.sh` Exit 1 mit 92,0 %, `bench-backfill.sh` Exit 0). Die Reihenfolge der vier
  Rezeptzeilen ist gelesen, nicht durch einen Gesamtlauf gezeigt.
- `bench-backfill.sh --full` (Stufe 1.000.000) nicht gefahren (der dritte Run überschreitet laut
  Handbuch und Verdikt die WAL-Fehlerschwelle); die Werte der 1.000.000-Zeilen-Läufe im Handbuch sind
  übernommen und im Repository nicht auflösbar.
- Der Abbruchzweig von `RUN_TIMEOUT_S` (F-4) ist gelesen, nicht ausgelöst.
- Store-Tier nur gegen PostgreSQL 18 (Repo-Default); `make test-store` läuft mit dem gepinnten Image.
- Übernommen, nicht selbst gefahren: die Review-Mutationen M1–M13, die Läufe `20260924T233628Z`,
  `20260925T000439Z`, `20260925T012459Z`, die Messungen des Verdikts (`pg_test_fsync`,
  `synchronous_commit=off`, `wal-verdikt[20260925T004731Z]`).

## 13. Ergebnis

| Prüfpunkt | Ergebnis |
|---|---|
| DoD §2 — `[x]`-Zeilen (8) | **8 von 8 erfüllt**, je mit eigenem Beleg (Nr. 1 in der umformulierten Form mit V-1/V-2, Nr. 2 mit V-3, Nr. 6 mit V-4, Nr. 7 mit V-6) |
| DoD §2 — übrige `[ ]`-Zeilen (5) | **5 von 5 korrekt offen** (Closure-/Rollen-Sequenz; Reconciliation „entfällt“) |
| Original-Kriterium „realer `make bench`-Lauf mit Exit 0“ | **nicht erfüllt und im Plan als solches benannt**; Ersatzbeleg (a) + (b) nachgemessen |
| Review-Findings F-1 bis F-12 | **F-1 bis F-6 behoben** (einzeln nachgemessen), F-7/F-8 unverändert ohne Aktion, F-9 bis F-12 im Träger gekennzeichnet bzw. behoben |
| Plan-vs-Code-Diff | **deckungsgleich**; „prüfen“-Zeilen unverändert; ungeplant nichts im Slice-Inhalt |
| Entscheidungen | [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md), [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md), [`ADR-0054`](../plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md), [`ADR-0104`](../plan/adr/0104-benchmark-schwellen-per-001-002-003.md), [`ADR-0114`](../plan/adr/0114-schema-rollout-vorlauf-view-signatur.md) **konform**; [`ADR-0120`](../plan/adr/0120-capture-slot-leerlauf-bestaetigung.md) nur Kontext, nicht umgesetzt |
| Mutationen der Eingabeseite | **drei rot** (MA, MB, MS), **eine überlebt** (MC, V-3) |
| Suchlauf-Feld | **elf Zeilen bestätigt** (neun an beiden Ständen, zwei über den Diff), zwei Selbstverweis-Zählungen (V-6) |
| Harte Regeln | **erfüllt** (§10) |
| Commit-Traceability, MR-Immutabilität | **grün** (`RANGE=933ab054..HEAD`) |
| Gates | **`make gates` EXIT=0**, `make test` EXIT=0, `make test-store` EXIT=0 |

## Verdikt

**Bestätigt.** Die DoD-Behauptung des Implementers ist durch eigene Messung belegt:
`make gates`, `make test` und `make test-store` laufen am Stand `c97273b3` mit Exit 0; das neue Skript
läuft einzeln mit Exit 0 (Lauf `20260925T015600Z`), und das Ende von `make bench` am vorhandenen
Skript `tools/bench-source-impact.sh` (92,0 % gegen 35 %, Ursache laut Verdikt die Flush-Latenz des
Hosts) ist nachgemessen. Der DoD-Punkt „Bench“ ist ehrlich formuliert: er nennt das Original-Kriterium
als nicht erfüllt und behauptet nur den Ersatzbeleg. Die Auswertung der beiden Warnungen färbt bei
Mutationen ihrer Eingabeseite rot (Grenze der Richtgröße, Vorzeichen der Dauer, View-Spalte), die
Findings des Reviews sind behoben, keine ADR-Festlegung und keine Schwelle sind berührt. Kein Befund
blockiert; V-1 und V-3 (LOW) sind Entscheidungen für den Planner, V-2 und V-4 bis V-7 INFO.

**Übergabe:** Verifier → Planner. Offen bleibt die Closure: Closure-Notiz mit Lerneintrag (die
Finding-Klassen des Reviews gehen in den Zähler), Ratifizierung der DoD-Umformulierung samt
Host-Abhängigkeit des Closure-Triggers (V-1), Anker oder Kennzeichnung der nicht auflösbaren Läufe
(V-2), Entscheidung zur Bindung der Toleranz-Umrechnung (V-3), Ausgänge der sieben §6-Risiken
(Toleranz: Startwert ohne Messung geführt; Richtgröße-wird-zur-Grenze: jede Nennung trägt Lauf und
„Orientierung“, keine Stelle lehnt ab; Schätzung so alt wie das letzte `ANALYZE`: gemessen, frisch
befüllte Tabelle „unbekannt“, Warnung (1) schweigt dort; Host-Abhängigkeit: Handbuch nennt Host und
Lauf; Laufzeit: Lauf von `bench-backfill.sh` etwa fünf Minuten (Ende 04:01 gegen Start 03:56, gemessen an den Zeitstempeln der Log-Datei), Gesamtlaufzeit von `make bench` durch das
rote Vorgängerskript nicht erreicht), Beobachtungs-Register, Reconciliation-Zeile, Welle-Paarungen.
Danach der reine `git mv` nach `done/` ([`AGENTS.md`](../../AGENTS.md) §3.3, Fall 2).
