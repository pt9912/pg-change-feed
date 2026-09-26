# Verifikations-Report: slice-transformationen-kern-rename — 2026-09-26

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich + Entscheidungs-Konformität +
Plan-vs-Code-Diff + Gates. Review-Artefakt des Reviewers:
[`review-slice-transformationen-kern-rename.md`](review-slice-transformationen-kern-rename.md);
Formvorbild dieses Reports: [`verifikation-slice-harness-suchlauf-nachmessen.md`](verifikation-slice-harness-suchlauf-nachmessen.md).
Die vendored Baseline trägt kein eigenes Verifikations-Template
(`.harness/baseline/v6.9.0/templates/docs/reviews/` enthält nur `review-report.template.md`); der
Report folgt dem Formvorbild.

**Gegenstand:** Slice-Plan `slice-transformationen-kern-rename` (Welle `welle-transformationen`), Stand
`HEAD` = `aec21d2c`, Diff-Range `7310dbd1..HEAD`, 10 Commits, 16 Dateien (13 Go-Dateien). Slice-Inhalt:
drei Lifecycle-/Plan-Commits (`64b8eb9e`, `44a8682f`, `1852b2f9`), Bau (`06719d4f`), Plan-Nachzüge und
DoD-Haken (`0d0190be`, `cfca864b`, `770fc754`), Fixrunde (`275f7f5c` Code, `aec21d2c` Plan). Nicht
Slice-Inhalt, im Range enthalten: der Review-Report des Reviewers (`d626312f`), gelesen am Stand
`770fc754`. Die Fixrunde ist damit von keinem Reviewer gelesen — ihre Wirkung belegen die eigenen
Mutationen in §4. Dieser Lauf ändert weder Code noch Plan noch Doku; er schreibt nur diesen Report. Alle
Mutationen liefen an einer Kopie des Baums (`git archive HEAD` in ein Scratch-Verzeichnis, Bind-Mount der
Kopie, Zeichenketten-Ersetzung mit `python3` auf der Kopie, kein `sed -i`); der Arbeitsbaum blieb
unberührt, `git status --short` war nach jedem Lauf leer.

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — [`AGENTS.md`](../../AGENTS.md) §3.9)

Jeder Lauf schrieb in eine Log-Datei; der Exit-Code wurde im selben Aufruf gesondert gesichert, die Logs
danach gelesen. Stand aller Läufe ohne Mutation: `HEAD` = `aec21d2c`, Arbeitsbaum sauber. Es lief ein
schwerer Docker-Lauf zugleich.

| Sensor | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `make test` (Race-Detector) | **EXIT=0** | 42 Zeilen `ok`, keine Zeile `FAIL`/`panic`/`DATA RACE`; darunter `ok … replication/mapper`, `ok … internal/bootstrap`, `ok … internal/domain/model` |
| `make a-check` | **EXIT=0** | `gesamt: 0 Befund(e)` |
| `make coverage-gate` | **EXIT=0**, zweimal gefahren | Lauf 1 `coverage-gate: OK — Coverage 83.90% erfüllt Schwelle 80%`, Lauf 2 `coverage-gate: OK — Coverage 83.80% erfüllt Schwelle 80%` (V-2) |
| `make gates` (einmal) | **EXIT=0** | `coverage-gate: OK — Coverage 83.80% erfüllt Schwelle 80%` · `d-check: 1209 Datei(en) geprüft, 0 Befund(e)` · `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` · `generated-sync: OK — das committete Erzeugnis ist byte-gleich …` · a-check `gesamt: 0 Befund(e)`; `git status --short` danach leer |
| `make suchlauf-nachmessen PLAN=<Plan des Slice>` | **EXIT=0** | 14 Zeilen `OK`, `suchlauf-nachmessen: 14 Zeilen stimmen` (soll = ist: 32 / 55 / 10 / 10 / 43 / 52 / 5 / 5 / 3 / 3 / 79 / 81 / 29 / 32) |
| `make commit-traceability RANGE=7310dbd1..HEAD` | **EXIT=0** | `OK — 10 Commit(s) in "7310dbd1..HEAD", Betreffs ohne Struktur-ID` |
| `make doc-commits RANGE=7310dbd1..HEAD` | **EXIT=0** | `d-check: 1209 Datei(en) geprüft, 0 Befund(e)` |
| `make doc-immutable RANGE=7310dbd1..HEAD` | **EXIT=0** | `d-check: 1209 Datei(en) geprüft, 0 Befund(e)` (ohne `RANGE` bricht das Ziel mit `flag needs an argument: --range` ab — Bedienung, kein Befund) |
| Benchmark `BenchmarkAssemblerChange…` | Parent und Diff je `-count 6 -benchtime 2s` | siehe §3, Zeile Kosten der Prüfung |
| Mutationen (§4) | sieben Läufe, alle rot | siehe §4 |

Hygiene: dangling Volumes (`docker volume ls -q -f dangling=true | wc -l`) am Ende **34**, kein `prune`, kein
`system prune`; freier Speicher am Ende 17,9 GB (`free -m`). Nicht gefahren: `make test-store`,
`make test-replication`, `make test-integration`, `make test-sdk-*-integration`, `make bench` — der Diff führt
keinen Weg zur Datenbank, keinen SQL-Zugriff, keine Betreiber-Oberfläche und keinen Zustellweg ein (Plan §1:
„kein SQL-Weg“).

## 2. DoD — Verdikt je Zeile (§2 des Plans)

Gezählt am Plan: neun `[x]`-Zeilen, vier `[ ]`-Zeilen, zusammen dreizehn (`grep -c '^- \[x\]'` 9,
`grep -c '^- \[ \]'` 4).

| # | DoD-Zeile | Verdikt | Realer Beleg |
|---|---|---|---|
| 1 | `rename_column` wirkt auf beide Images (Position der Quellspalte, Wert unverändert, Abwesenheit bleibt, Metadaten unverändert, leere Regelmenge = Bytes wie vor dem Slice) | **erfüllt** | `make test` Exit 0 (§1). Code gelesen: `BuildRowImage` prüft je Spalte `i >= len(values)`/`nil`/`excluded` und ruft erst danach `applyTransformations`; der Schlüssel wird in der Schleife an der Stelle der Quellspalte geschrieben. **Byte-Gleichheit ohne Regel:** `git diff -w 1852b2f9 HEAD -- internal/domain/model/rowimage_test.go` gelesen — geändert sind nur das Verschieben der Tabelle in die Paketvariable `rowImageByteCases` und das vierte Argument `nil`; keine Erwartungs-Zeichenkette geändert; `TestBuildRowImageEmptyRuleSetKeepsBytes` fährt dieselbe Tabelle mit leerer Regelliste. `mapper_test.go` nicht im Diff, `snapshot_test.go` zwei Zeilen (`git diff --stat`). Mutation M7 (Ausschluss erst gegen den Zielschlüssel) rot (§4) |
| 2 | `Assembler` trägt Regelstand und Prüfung (Ersetzen unter `tablesMu`, Erhalt bei `AddBinding`-Merge und `setSchemaVersion`, Prüfung vor Serialisierung, `classifyRunError` → `schema`, Negativtests an ihre Eingabe gebunden, Fall „Kollision nach kompatibler Erweiterung“, spalten-entfernende Fälle enden weiter vorher an `relationOther`) | **erfüllt** | Code gelesen (§3). Tests `TestAddBindingKeepsRuleState`, `TestConsumeRuleSurvivesSchemaBump`, `TestConsumeRuleTargetCollisionAfterCompatibleExtension`, `TestConsumeEveryRuleIsCheckedForApplicability`, `TestClassifyRunErrorMapsKnownSentinelsToADR0023Classes` vorhanden und in `make test` grün. Eingabeseiten-Mutationen selbst gefahren, alle rot: M3 (Prüfung nur der ersten Regel), M5 (Abbildung auf `schema` entfernt), M6 (`AddBinding` erhält den Regelstand nicht), M4 (§4). Der Pfad `receive.Stream.process` gibt den Fehler von `Consume` unverändert weiter und erreicht `Capture` nicht (`receive.go:543-546`, gelesen) |
| 3 | Fitness Function [`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md): Eigenschaftstest über die Regeltypen der Domänen-Menge × ausgeschlossene Spalte; Determinismus; Nebenläufigkeit; `make a-check`; `make coverage-gate` | **erfüllt** | `TestExcludedColumnIsUnreachableForEveryRuleKind` gelesen: `for _, kind := range model.TransformationKinds()` × jede der vier Spalten als ausgeschlossene, prüft weder Quellschlüssel noch Zielname noch Quellwert in Alt- und Neu-Image, ein unbekannter Typ bricht mit „nicht abgedeckt“ ab (`TransformationKinds` wird außerhalb von `transformation.go` nur im Domänen-Test der Menge und in diesem Test aufgerufen — die eine Quelle). `TestRenameColumnImagesAreDeterministic` (getrennte Assembler-Instanzen, wiederholt) und `TestAssemblerTransformationsAreRaceFree` (zwei Goroutinen, jedes Image ist einer der beiden vollständigen Stände) gelesen, beide in `make test` unter `-race` grün. `make a-check` Exit 0, `make coverage-gate` Exit 0 (§1). `coverage`-Stufe des Dockerfile schließt nur `postgresstorage`, `postgresack`, `postgressnapshot`, `replication/receive` aus (im Log `grep -vE` gelesen) — kein neues Paket |
| 4 | Kommentar-Träger: `TableBinding`/`AddBinding` nennen den Regelstand; `ErrTransformationNotApplicable` sagt nur zu, was der Code trägt | **erfüllt, mit V-1 (LOW) am Nachbarkommentar** | Diff gelesen: `TableBinding`, `AddBinding`, `setSchemaVersion`, `Assembler` nennen `Transformations`. Der Kommentar am Sentinel („`receive.Stream` beendet damit den Lauf, die Transaktion erreicht `Capture` nicht und wird nicht bestätigt“) ist wahr: `process` gibt den Fehler weiter, `Capture` wird nicht gerufen. Er nimmt den Kollisionsfall zweier Regeln ausdrücklich aus („erst die Row-Image-Konstruktion … meldet“) und trägt keine Behauptung, die Abhilfe wirke im gescheiterten Prozess. V-1: der Doc-Kommentar von `change` sagt „ohne dass … die Sequenz vorrückt“, was für den Kollisionsfall nicht gilt |
| 5 | `make gates` grün, Exit gesondert | **erfüllt** | eigener Lauf EXIT=0 (§1) |
| 6 | Review durchgeführt, Report liegt vor | **erfüllt, mit V-3 (INFO)** | Report liegt vor (1 HIGH, 1 MEDIUM, 2 LOW, 3 INFO; Verdikt merge-blockierend), Fixrunde `275f7f5c`/`aec21d2c`; Finding für Finding in §5 nachgemessen — kein offenes HIGH/MEDIUM. V-3: die Fixrunde hat keinen zweiten Reviewer-Lauf |
| 7 | §3.13-Suchlauf: committetes Feld in §3, Gefundenes und Nichtgefundenes je Träger, beide Stände | **erfüllt** | 14 Zeilen mit dem Werkzeug Exit 0 (§1); acht Zeilen und die Nichtgefundenes-Behauptungen von Hand nachgefahren (§6); die Befund-Zellen stimmen |
| 8 | Doku-Update entfällt (kein öffentlicher Vertrag berührt) | **erfüllt (entfällt)** | `git diff --name-only 7310dbd1 HEAD -- docs/user harness tools spec compose.yaml Makefile` leer; Kandidatenlauf auf `CDC_*`, `os.Getenv`, `cdc.<name>(`, `HandleFunc`, Pfad-Literale in den zugefügten Go-Zeilen: 0 Treffer — keine neue Betreiber-Oberfläche (Handbuch-Kandidatenlauf) |
| 9 | Reconciliation-Register — entfällt | **erfüllt (entfällt)** | Greenfield |
| 10 | Closure-Notiz mit Lerneintrag | **korrekt offen** | Plan §7 trägt „*(zu tragen bei Closure)*“ (Planner) |
| 11 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | Closure-Pflicht (Planner) |
| 12 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | Plan §6: acht Zeilen mit „Ausgang: *(bei Closure …)*“ (Planner) |
| 13 | Die drei Paarungen | **korrekt offen** | hängen an der Closure der Welle |

Kein `[x]` ohne Beleg; kein `[ ]`, das über die Rollen-Sequenz hinaus belegt wäre. Register,
Risiko-Ausgänge, Paarungen und Closure-Notiz sind Planner-Arbeit und nicht Teil dieser Prüfung.

## 3. Plan-vs-Code-Diff

**§1 Ziel:** Domäne (`transformation.go`), erweiterte gemeinsame Funktion (`BuildRowImage`), Regelstand und
Prüfung im `Assembler`, `classifyRunError` stehen. Die „Ausdrücklich NICHT“-Punkte sind eingehalten: kein
SQL-/Antragsweg (Diff nennt weder `tools/schema` noch einen Use Case oder Port), keine K1–K4-Prüfung des
Regelsatzes im `Assembler` (nur Anwendbarkeit je Change), kein `map_value` (`TransformationKinds` liefert
genau `rename_column`), Backfill-Aufrufer übergibt `nil` (`service.go:489`), keine Startreihenfolge.

**§3-Tabelle, Zeile für Zeile gegen den Diff (`git diff --stat 7310dbd1 HEAD`):**

| Plan-Zeile | Ist im Diff |
|---|---|
| `transformation.go` neu | vorhanden, 113 Zeilen: `TransformationKind`, `NewRenameColumn`, `CheckApplicable`, `applyTransformations`, `TransformationKinds` |
| `rowimage.go` update (vierter Parameter `rules`, Kollisionsprüfung) | Signatur `BuildRowImage(columns, values, excluded, rules)`; Kollision gegen `columns` und gegen `renamed` liefert `ErrTransformationTargetCollides` und kein Bild |
| `mapper.go` update | `TableBinding.Transformations`, `SetTransformation`/`RemoveTransformation`, Erhalt in `AddBinding`, Prüfung in `change`, `ErrTransformationNotApplicable`, `imageError` — alle im Diff; `setSchemaVersion` setzt nur die Version und lässt den Regelstand stehen (`TestConsumeRuleSurvivesSchemaBump`) |
| `wiring.go` update | eine `case`-Zeile in `classifyRunError` |
| `service.go` update | `nil` als Regelmenge, Kommentar nennt die Adresse `slice-transformationen-backfill-pfad` |
| `errors.go` update, vier Sentinels | `ErrInvalidTransformation`, `ErrTransformationTargetIsColumn`, `ErrTransformationColumnMissing`, `ErrTransformationTargetCollides` |
| `rowimage_test.go`, `snapshot_test.go` update | Aufrufer nachgezogen; Tabelle als `rowImageByteCases`, Erwartungen unverändert |
| `transformation_test.go` (Domäne) neu | vorhanden, 339 Zeilen, `TestBuildRowImageNeverWritesTwoEqualKeys` eingeschlossen |
| Mapper-Tests in eigenen Dateien, Benchmark | `transformation_test.go` (722 Zeilen), `transformation_internal_test.go`, `mapper_bench_test.go` vorhanden; `mapper_test.go` nicht im Diff |
| `heartbeat_internal_test.go` update | zwei Zeilen in der Sentinel-Tabelle |

Nicht im Plan, im Diff: der Review-Report des Reviewers. Im Plan, nicht im Diff: nichts.
`git diff --name-only 7310dbd1 HEAD -- docs/plan/adr spec .d-check.yml .a-check.yml Dockerfile harness/mk`
ist leer: keine ADR, keine Spec, keine Schwelle, keine Gate-Konfiguration verändert
([`AGENTS.md`](../../AGENTS.md) §3.5, §3.6).

**Festlegungen des Plans gegen den Code:**

| Festlegung | Ist |
|---|---|
| Regelname Teil der Domänen-Regel, Invariante „nicht leer“ | `NewRenameColumn(name, column, to)`; `name == ""` → `ErrEmptyIdentifier` |
| Anwendbarkeit als Domänen-Methode | `Transformation.CheckApplicable(columns)`; Mapper wrappt Regelname, Tabelle und Grund mit `%w` |
| `BuildRowImage` prüft Anwendbarkeit nicht, schreibt aber nie zwei gleichnamige Schlüssel | `containsName(columns, key) \|\| containsName(renamed, key)` an der Schreibstelle; Mapper ordnet über `imageError` als `ErrTransformationNotApplicable` ein; Mutationen M1, M2, M4 rot (§4) |
| `SetTransformation` ersetzt unter demselben Namen an der Stelle | `withTransformation` — Test `TestSetAndRemoveTransformationOnLiveBinding` und der Schnappschuss-Fall „Ersetzen“ |
| erste treffende Regel entscheidet | `applyTransformations` gibt bei der ersten Regel mit `rule.column == column` zurück |

**Konstruktor gegen [`SPEC-030`](../../spec/pflichtenheft.md) (Bezeichner):** `column` nicht leer und ohne
U+0000 → `ErrInvalidTransformation`; `to` nicht leer, ohne U+0000, `len(to) > 63` (Bytes, nicht Runes) →
`ErrInvalidTransformation`; `to == column` → eigener Sentinel `ErrTransformationTargetIsColumn`
(K3-Fall, [`SPEC-019`](../../spec/pflichtenheft.md)/`SPEC-030` Randfälle: „verletzt K3 … endet den Antrag `failed`“);
Vergleich zeichengenau (kein Falten). Alle drei Bezeichner-Invarianten der Prüfaufträge sind im Code
(`transformation.go:55-63`) und im Test `TestNewRenameColumnInvariants` je Eingabe gebunden (Reviewer-Mutationen M5/M20 zusätzlich).

**Kosten der Prüfung (Plan §6, Benchmark), eigene Messung** (`go test -run '^$' -bench BenchmarkAssemblerChange
-count 6 -benchtime 2s`, im Toolchain-Container von `make test` ohne `-race`, Parent `1852b2f9` mit einer
Benchmark-Datei ohne Regel-Variante, Diff = `HEAD`), gedruckte Zeilen:

| Stand | ns/op (sechs Läufe) | B/op | allocs/op |
|---|---|---|---|
| Parent, ohne Regel | 5519, 7305, 5801, 4945, 5137, 5092 | 2444 | 86 |
| Diff, ohne Regel | 5213, 5221, 5321, 5197, 5312, 5301 | 2444 | 86 |
| Diff, mit zwei Regeln (`…WithRules`) | 5450, 5361, 5550, 5673, 5520, 5424 | 2444 | 86 |

Die Allokationen (2444 B/op, 86 allocs/op) stimmen mit der Plan-Aussage „unverändert“ überein, auch nach der
Kollisionsprüfung (`renamed` bleibt ohne Regel auf dem Stapel). Die Laufzeit-Spannen des Plans (Parent
4777–4958, Diff ohne Regel 5005–5213 bzw. 4871–5275, mit zwei Regeln 5233–5439 bzw. 5102–5356 ns/op) sind
Läufe der Implementer-Messung, die meine nicht Zeile für Zeile wiederholt (Parent-Spanne meines Laufs
4945–7305 ns/op, der Ausreißer 7305 zeigt das Rauschen); die Aussage „0,2–0,6 µs, nicht gegen Rauschen
abgesichert“ trägt auch meine Messung. Abgeleitet (Mediane meiner Läufe): Parent 5,3 µs, Diff ohne Regel
5,3 µs, mit zwei Regeln 5,5 µs — der Aufschlag liegt innerhalb der Parent-Streuung.

## 4. Mutationen der Eingabeseite (dieser Lauf)

Je Mutation eine Zeichenketten-Ersetzung an der Kopie, dann `go test -race` über
`internal/domain/model`, `internal/adapters/driving/replication/mapper`, `internal/bootstrap` im selben
Toolchain-Image wie `make test`; Ausgang und `--- FAIL`-Zeilen gelesen; die Kopie danach aus der
Sicherungsdatei zurückgesetzt (`diff -r` Kopie gegen Arbeitsbaum: identisch). Ausgangslauf der unveränderten
Kopie: Exit 0, drei Pakete `ok`.

| # | Mutation | Ergebnis (gedruckt) |
|---|---|---|
| M1 | Kollisionsprüfung gegen `columns` in `BuildRowImage` entfernt (`if containsName(renamed, key)`) | **EXIT=1**: `--- FAIL: TestBuildRowImageNeverWritesTwoEqualKeys`. Der Mapper-Test bleibt grün, weil `checkTransformations` den Fall vorher fängt — beide Wächter tragen ihren Teil |
| M2 | Prüfung gegen die umbenannten Schlüssel entfernt (`if containsName(columns, key)`) | **EXIT=1**: `TestBuildRowImageNeverWritesTwoEqualKeys`, `TestConsumeTwoRulesWithSameTargetAreNotApplicableWhenBothMatch` |
| M3 | `checkTransformations` über `rules[:min(1, len(rules))]` (Prüfung nur der ersten Regel) | **EXIT=1**: `--- FAIL: TestConsumeEveryRuleIsCheckedForApplicability` (Review-Finding F-2, dort grün, jetzt gebunden) |
| M4 | `imageError` beim Neu-Image-Aufruf ausgelassen (`return nil, err`) | **EXIT=1**: `TestConsumeTwoRulesWithSameTargetAreNotApplicableWhenBothMatch` (der Fehler trägt `ErrTransformationNotApplicable` nicht, `classifyRunError` bildete ihn nicht auf `schema` ab) |
| M5 | `ErrTransformationNotApplicable` aus `classifyRunError` entfernt | **EXIT=1**: `--- FAIL: TestClassifyRunErrorMapsKnownSentinelsToADR0023Classes` |
| M6 | `AddBinding` erhält `existing.Transformations` nicht | **EXIT=1**: `--- FAIL: TestAddBindingKeepsRuleState` |
| M7 | Ausschluss gegen den Zielschlüssel statt gegen die Quellspalte (Regel vor Ausschluss) | **EXIT=1**: `TestBuildRowImageRenameColumn`, `TestBuildRowImageNeverWritesTwoEqualKeys`, `TestExcludedColumnIsUnreachableForEveryRuleKind` |

Die Eingabeseiten der beiden gemeinsam geforderten Fälle (Kollisionsprüfung und Prüfung jeder Regel) sind
damit an ihre Tests gebunden; die übrigen Mutationen des Reviewers (M1–M20) habe ich nicht wiederholt
(übernommen aus dem Review-Report, dort je einzeln gefahren).

## 5. Findings des Reviews nachgemessen (nicht dem Fixrunden-Bericht geglaubt)

| Finding | Was ich gemessen habe | Verdikt |
|---|---|---|
| F-1 (HIGH) Doc-Kommentar von `BuildRowImage` sagt „eine nicht anwendbare Regel wirkt hier nicht“, der Code trägt den Kollisionsfall nicht | Code (`rowimage.go:55-64`) und Doc-Kommentar (`:34-39`) gelesen: der Kommentar sagt jetzt „Eine Regel, deren Spalte nicht in `columns` steht, trifft keine Spalte und wirkt nicht“ und „schreibt nie zwei gleichnamige Schlüssel … liefert `ErrTransformationTargetCollides` und kein Bild“; beides trägt der Code. Test `TestBuildRowImageNeverWritesTwoEqualKeys` (vier Spaltenfälle, ein Regelfall, Gegenprobe mit einer wertlosen Quellspalte) an die Eingabe gebunden (M1, M2 rot). **`imageError`-Einordnung:** `errors.Is(err, ErrTransformationTargetCollides)` → `%w`-Wrapping mit `ErrTransformationNotApplicable` und dem Grund; beide Sentinels sind über `errors.Is` erreichbar (Test prüft beide), `classifyRunError` bildet auf `schema` ab (M5 rot); ohne die Einordnung (M4) rot. **Kommentar am Sentinel** gelesen: wahr (siehe DoD 4). Rest: V-1, V-4 | **behoben**, mutiert bestätigt |
| F-2 (MEDIUM) Prüfung nur für die erste Regel gebunden | `TestConsumeEveryRuleIsCheckedForApplicability` (vier Fälle: zweite, erste, letzte von drei, mittlere von drei); M3 rot | **behoben**, mutiert bestätigt |
| F-3 (LOW) Plan-Tabelle trägt überholte Zeilen | Plan §3 gelesen: keine Zeile „falls vorhanden“, keine `mapper_test.go`-Zeile „update / neu“, keine Zeile „Datei und Name am Start gemessen“ mehr; je Gegenstand eine wahre Zeile | **behoben** |
| F-4 (LOW) Coverage-Zahl ohne Lauf-Stand | Plan-Belege nennen jetzt „83.90%“ mit dem Zusatz „Alle Messungen dieses Punkts stehen am Endstand des Codes“, die alte Zahl „83.70%“ ist entfernt. Eigene Läufe am Endstand: 83.90 %, 83.80 %, 83.80 % (V-2) | **behoben** (Form); Zahl ist ein Lauf-Beleg, keine Konstante |
| F-5 (INFO) Konfliktfreiheit zwischen Regeln ohne Laufzeit-Wächter | jetzt ein zweiter Wächter in `BuildRowImage`/`imageError`; K3 bleibt Prüfung des Antragswegs (Plan §1, `SetTransformation`-Kommentar) | **behoben/benannt** |
| F-6 (INFO) drei Randfälle | Plan §6 führt einen Punkt „Randfälle ohne eigene Bindung oder mit geteiltem Speicher (benannt)“ mit (a) Regelliste geteilt, (b) Prüfreihenfolge in `CheckApplicable`, (c) UTF-8-Gültigkeit; Ausgang „(bei Closure)“ | **benannt** |
| F-7 (INFO) Kommentar an `blockBuilder.build`, gemeldete fremde Träger | Träger selbst gefunden: `docs/plan/adr/0115-backfill-spaltenwerte-text-ergebnisformat.md` Zeilen 129 und 196, `docs/plan/planning/open/slice-transformationen-backfill-pfad.md` Zeile 172 tragen die Signatur mit drei Argumenten (`git grep -n 'BuildRowImage(columns' -- docs spec harness` mit den Ausschlüssen: 3 Zeilen); Plan §3 meldet sie mit Frist „Closure dieses Slice“. Kenntnis: die ADR bleibt unberührt (`Accepted`), der Planner zieht die Plan-Zeile nach | **gemeldet wie behauptet** |

Kein offenes HIGH, kein offenes MEDIUM.

## 6. Suchlauf-Feld (§3 des Plans), an beiden Ständen nachgefahren

Parent `1852b2f9`. Von Hand mit `git grep -n … | wc -l` (ohne das Werkzeug), acht der 14 Zeilen und die
„Nichtgefunden“-Behauptungen; Stand `diff` ist der Arbeitsbaum dieses Laufs (`HEAD` = `aec21d2c`, sauber):

| Zeile | Stand | Soll | Von Hand |
|---|---|---|---|
| 1 `BuildRowImage` über `internal/*.go` | `1852b2f9` | 32 | **32** |
| 2 dasselbe | `diff` | 55 | **55** |
| 4 dasselbe ohne `_test.go` | `diff` | 10 | **10** (Definition, drei Aufrufstellen: `mapper.go:270`, `:274`, `service.go:489`; sechs Doc-Kommentare) |
| 6 `TableBinding{` | `diff` | 52 | **52** |
| 8 dasselbe ohne Tests | `diff` | 5 | **5** (`config_file.go:250`, `wiring.go:370`, `:381`, `:425`, `:1404`) |
| 10 `BuildRowImage(columns` über `docs spec harness` mit Ausschlüssen, Plan-Datei zusätzlich ausgeschlossen | `diff` | 3 | **3** ([`ADR-0115`](../plan/adr/0115-backfill-spaltenwerte-text-ergebnisformat.md) Zeilen 129, 196; `backfill-pfad` Zeile 172) |
| 12 `ExcludedColumns` | `diff` | 81 | **81**; je Datei gegenüber dem Parent (`git grep -c`, sortiert, `diff`): nur `mapper.go` 9 → 10 und `transformation_test.go` 0 → 1 mehr — wie in der Befund-Zelle |
| 14 `ErrIncompatibleSchemaChange` | `diff` | 32 | **32** |

Die übrigen sechs Zeilen (Parent-Stände 10, 43, 5, 3, 79 und 29) sind nur mit dem Werkzeug nachgemessen; bei `ExcludedColumns` bestätigt der
Vergleich der Trefferzahlen je Datei zwischen Parent und Diff den Abstand von zwei Zeilen (23 Dateien mit Plan-Datei, 22 ohne). „Nichtgefunden“:
`git grep -c Transformations 1852b2f9 -- internal` druckt keine Datei (bestätigt); `json.Marshal` an Row
Images: nur `rowimage.go:65,69` (bestätigt, kein zweiter Bild-Erzeuger, [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 2). Die Zahlen streuen mit jedem Commit, der die Muster berührt.

## 7. Entscheidungs-Konformität

- **[`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) Teilfrage 2/3:**
  geschlossener Satz (ein Typ `rename_column`, als Menge geführt), Position der Quellspalte behalten (Schleife
  in Relation-Spaltenreihenfolge, Tests `mittlere Spalte`/`erste Spalte`), Ausschluss zuerst, einmalige
  Auswertung gegen die Original-Spaltennamen. K1–K4 bleiben Sache des Antragswegs, der Assembler prüft nur
  Anwendbarkeit — wie Plan §1 festlegt.
- **Teilfrage 4:** `Assembler.change` prüft die Regeln der Bindung gegen die Spalten der Relation **vor** der
  Serialisierung, meldet `ErrTransformationNotApplicable`, `classifyRunError` bildet auf `ErrorClassSchema`
  ab; die Transaktion erreicht `Capture` nicht. „nach Abhilfe setzt die Erfassung fort“: `TestConsumeRuleRemedyRestoresCapture`.
  Die spalten-entfernenden Fälle enden weiter an `relationOther` (Test mit `errors.Is` in beide Richtungen, übernommen aus dem Review-Report, von mir nicht gesondert gefahren).
- **Teilfrage 5:** Ausschluss und Regel in einer Schleife, Ausschluss vor Regel (M7 rot). Der Schlüssel der
  ausgeschlossenen Spalte erscheint weder unter Quell- noch unter Zielname ([`LH-QA-SEC-004`](../../spec/lastenheft.md),
  Eigenschaftstest).
- **Teilfrage 6:** Regelstand als unveränderliche Liste in der Bindung, Ersetzen unter `tablesMu`, Erhalt bei
  `AddBinding`-Merge und `setSchemaVersion` (M6 rot; `setSchemaVersion` ändert nur die Version und lässt den
  Rest, Test `TestConsumeRuleSurvivesSchemaBump`).
- **Folgepflicht 2:** alle Bausteine der Aufzählung vorhanden (`model`-Regeltyp mit Konstruktor-Invarianten
  und reiner Auswertung, `TableBinding.Transformations`, `SetTransformation`/`RemoveTransformation`,
  Anwendbarkeits-Prüfung, `ErrTransformationNotApplicable`, `classifyRunError`); die Nebenläufigkeitstests
  gibt es. **Ort der Auswertung:** die ADR sagt „in `rowImage`“ (Teilfrage 5, Folgepflicht 2), der Code wertet in
  der gemeinsamen Domänen-Funktion `BuildRowImage` aus; das private `rowImage` des Mappers ist mit
  `slice-backfill-row-image-gemeinsam` entfallen. Die Abweichung vom Wortlaut ist im Plan §6 als Risiko
  benannt, sie steht als Auslegung und nicht als Entscheidung; die Wirkung (eine Schleife, Ausschluss zuerst) ist
  erhalten. Ich lese sie als konform mit dem Zweck der Entscheidung, nicht als Verstoß.
- **Fitness Function (drei Zeilen der Tabelle in der ADR):** Eigenschaftstest im `mapper`-Paket über die
  Regeltypen der Domänen-Menge — vorhanden; Determinismus und Nebenläufigkeit — vorhanden (beide unter `-race`);
  `make a-check` — Exit 0; die vierte Zeile (`make test-integration`) ist Folgepflicht 5 und nicht Teil dieses
  Slice.
- **[`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 2:** genau **eine**
  Row-Image-Konstruktion — die drei Aufrufstellen im Produktivcode sind `Assembler.change` (zwei) und
  `blockBuilder.build` (eine), beide rufen `model.BuildRowImage`; `json.Marshal` an Bildern nur in `rowimage.go`.
  Der Backfill-Pfad übergibt `nil` (bis `backfill-pfad`, [`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) Folgepflicht 7).
- **[`ADR-0059`](../plan/adr/0059-spaltenauswahl-mechanismus.md):** Muster `ExcludeColumn` gespiegelt
  (Listen neu aufgebaut, Set/Remove auf nicht getragener Bindung ohne Wirkung und ohne Beleben,
  Schnappschuss ohne eigene Sperre).
- **[`ADR-0023`](../plan/adr/0023-fehlerklassifikation.md):** keine achte Klasse, Abbildung auf `schema`.
- **[`SPEC-030`](../../spec/pflichtenheft.md):** Bezeichner-Invarianten im Konstruktor (siehe §3); Anwendbarkeit
  hängt für die Regel-gegen-Spalten-Prüfung an Regel und Spaltenmenge, nie am Wert (V-4 zum zweiten Wächter).
- **[`AGENTS.md`](../../AGENTS.md):** §3.3 (die Moves `64b8eb9e` und `1852b2f9` sind reine Renames, 0 Zeilen
  geändert; Inhalt steht in eigenen Commits), §3.5/§3.6 (keine ADR/Spec/Schwelle im Diff), §3.7 (§3.7-Probe
  über die zugefügten Go-Zeilen auf `slice-`/`welle-`-Kennungen: ein Treffer, die namentliche Adresse
  `slice-transformationen-backfill-pfad` in `service.go` als Rang-Zeiger, benannte Grenze), §3.9 (meine
  Läufe ungepiped), §3.12 (Zahlen im Plan tragen Lauf und Ursprung; V-2), §3.13 (fremde Träger gemeldet statt
  still mitgeändert). §3.2: kein `//nolint`.
- **Commits:** `make doc-commits` und `make commit-traceability` über die Range Exit 0; jeder der zehn
  Betreffs nennt `LH-FA-CFG-007` und `ADR-0112`, keiner trägt `SPEC-`/`ARC-`.
- **Handbuch-Kandidatenlauf:** `docs/user/benutzerhandbuch.md` nicht im Diff; keine neue Umgebungsvariable,
  keine SQL-Funktion, kein Endpunkt, keine Flagge — keine Betreiber-Oberfläche, der Nachzug bleibt
  `slice-transformationen-betriebsdoku` (Adresse im Plan).

## 8. Findings dieser Verifikation

| # | Kategorie | Befund | Quelle | Verifizierbar |
|---|---|---|---|---|
| V-1 | LOW | Der Doc-Kommentar von `Assembler.change` sagt, eine nicht anwendbare Regel ende als `ErrTransformationNotApplicable`, „ohne dass ein Bild entsteht oder die Sequenz vorrückt“. Für die Kollision zweier Regeln (Row-Image-Konstruktion, `imageError`) rückt die Sequenz vor (`a.open.sequence++` steht vor dem `BuildRowImage`-Aufruf), und beim Update kann das Neu-Image entstehen, bevor das Alt-Image die Kollision meldet. Folgenlos für den Lauf (er endet, nichts wird persistiert oder bestätigt); der Sentinel-Kommentar nimmt den Fall aus, der Nachbarkommentar nicht. | [`AGENTS.md`](../../AGENTS.md) §3.7 (Kommentar-Zusage); `internal/adapters/driving/replication/mapper/mapper.go` (Doc von `change`, `a.open.sequence++` Zeile 255) | ja — Lesen der beiden Stellen; Kollisions-Test zeigt den Pfad |
| V-2 | INFO | Die Coverage-Zahl am unveränderten Endstand schwankt zwischen Läufen: eigene Läufe 83.90 %, 83.80 %, 83.80 % (der Reviewer am Stand `770fc754` zweimal 83.80 %, der Plan-Beleg 83.90 %). Der Plan nennt 83.90 % für `make coverage-gate` und für `make gates`; in meinem `make gates`-Lauf stand 83.80 %. Die Zahl ist ein Beleg eines Laufs, keine Konstante des Codes; die Schwelle 80 % ist in allen Läufen erfüllt. | [`AGENTS.md`](../../AGENTS.md) §3.12 Instanz A; Plan §3 „Belege des Laufs“ | ja — `make coverage-gate` wiederholt |
| V-3 | INFO | Der Review-Report liest die Range bis `770fc754`; die Fixrunde (`275f7f5c`: Kollisionsprüfung in `BuildRowImage`, `imageError`, zwei neue Tests) hat keinen zweiten Reviewer-Lauf. Ihre Wirkung ist in §4 (M1–M4) und §5 mutiert bestätigt; die Fixrunde liegt im heißen Pfad (`BuildRowImage`), belegt unverändert 86 allocs/op (§3). Ein Reviewer-Blick bleibt eine Rollen-Entscheidung des Planners. | `internal/domain/model/rowimage.go`, `mapper.go` `imageError` | ja — Lesen; §4 |
| V-4 | INFO | [`SPEC-030`](../../spec/pflichtenheft.md) sagt: „Anwendbarkeit hängt an Regelmenge und Spaltenmenge, nie am Wert einer Zeile“. Der zweite Wächter in `BuildRowImage` (zwei Regeln mit gleichem Zielnamen) löst nur aus, wenn beide Quellspalten in derselben Zeile einen Wert tragen (Test-Gegenprobe: mit einer wertlosen Quellspalte geht die Änderung durch). Das ist kein Widerspruch zur Spec — die Konstellation ist nach K3 am Antrag ausgeschlossen, der Wächter ist der zweite —, aber der einzige Anwendbarkeits-Pfad, der am Wert hängt. Im Kommentar am Sentinel benannt („an einer Zeile, die beide Werte trägt“). | [`SPEC-030`](../../spec/pflichtenheft.md) Anwendbarkeit, [`SPEC-019`](../../spec/pflichtenheft.md) K3; `mapper.go` Kommentar an `ErrTransformationNotApplicable` | ja — Lesen; `TestConsumeTwoRulesWithSameTargetAreNotApplicableWhenBothMatch` Gegenprobe |
| V-5 | INFO | Zwei Erwartungen des Plans sind Zusagen der Folge-Slices und hier nicht belegbar: der Kommentar an `blockBuilder.build` nennt einen Slice-Namen (zeitgebunden, Rang-Zeiger; muss mit `slice-transformationen-backfill-pfad` umgeschrieben werden), und die drei fremden Träger (`backfill-pfad` Zeile 172, [`ADR-0115`](../plan/adr/0115-backfill-spaltenwerte-text-ergebnisformat.md) Zeilen 129 und 196) tragen die alte Signatur mit drei Parametern. Beides ist im Plan gemeldet, mit Frist (Closure dieses Slice, Planner). | Plan §3 Suchlauf-Feld Zeile 2, §6 „Zwischenzustand im Backfill-Pfad“ | ja — `git grep -n 'BuildRowImage(columns'` |

Kein HIGH, kein MEDIUM. Keine DoD-Verletzung.

## 9. Verdikt

**DoD bestätigt:** ja — jede der neun `[x]`-Zeilen (Nr. 1–9) ist am Ist-Zustand belegt; die vier `[ ]`-Zeilen
(Nr. 10–13) sind korrekt offen (Planner-Closure). **Plan-vs-Code:** keine unbenannte Abweichung; die
Abweichung vom ADR-Wortlaut „in `rowImage`“ steht als benannte Auslegung im Plan §6.
**Entscheidungs-Konformität:** [`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
(Teilfrage 2/3/4/5/6, Folgepflicht 2, Fitness Function),
[`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 2,
[`ADR-0059`](../plan/adr/0059-spaltenauswahl-mechanismus.md) und
[`ADR-0023`](../plan/adr/0023-fehlerklassifikation.md) konform; [`SPEC-030`](../../spec/pflichtenheft.md)-Bezeichner-Invarianten
im Konstruktor gedeckt. **Gates:** `make test`, `make a-check`, `make coverage-gate` (83.80–83.90 %),
`make gates` und `make suchlauf-nachmessen` Exit 0 im eigenen Lauf; sieben Eingabeseiten-Mutationen rot.

**Übergabe:** an den Planner — Closure-Notiz mit Lerneintrag, Beobachtungs-Register, Ausgänge der §6-Risiken
(darunter der Ausgang der Zeile „Der Ort der Auswertung weicht vom ADR-Wortlaut ab“ und der Randfall-Zeile),
Paarungen bei der Closure der Welle; die gemeldeten fremden Träger (`slice-transformationen-backfill-pfad`
Zeile 172, [`ADR-0115`](../plan/adr/0115-backfill-spaltenwerte-text-ergebnisformat.md) Zeilen 129 und 196) nach; V-1 als
Kommentar-Nachzug (Implementer) oder Anmerkung, V-2/V-3 als Anmerkungen zur Entscheidung des Planners.
Dieser Report ist ein **Lauf-Beleg** (dieser Stand, dieser Lauf) und ersetzt weder Review noch Closure.
