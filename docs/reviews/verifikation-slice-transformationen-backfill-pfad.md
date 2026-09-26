# Verifikations-Report: slice-transformationen-backfill-pfad — 2026-09-26

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich + Entscheidungs-Konformität +
Plan-vs-Code-Diff + Gates. Review-Artefakt des Reviewers:
[`review-slice-transformationen-backfill-pfad.md`](review-slice-transformationen-backfill-pfad.md);
Formvorbild dieses Reports:
[`verifikation-slice-transformationen-antragsweg-usecase.md`](verifikation-slice-transformationen-antragsweg-usecase.md).
Die vendored Baseline trägt kein eigenes Verifikations-Template
(`.harness/baseline/v6.9.0/templates/docs/reviews/` enthält nur `review-report.template.md`, Gerüst per `cp`
übernommen und durch die Form des Formvorbilds ersetzt).

**Gegenstand:** Slice-Plan `slice-transformationen-backfill-pfad` (Welle `welle-transformationen`), `HEAD` =
`17d0048a`, Diff-Range `3973390e..HEAD`, 9 Commits, 15 Dateien (+1330/−128 einschließlich zweier Lifecycle-Moves).
Slice-Inhalt: Lifecycle (`42ea6d7c`, `4661022c`, `24b1d7d9`), Implementer-Lauf (`210bf98d`, `fc0b8d38`, `56f40b93`,
`13aaf752`, `abebcf5b`), Review-Report (`17d0048a`, ohne Fixrunde). Der Stand ist nicht gepusht
(`origin/main` = `3973390e`). Dieser Lauf ändert weder Code noch Plan noch Doku; er schreibt nur diesen Report. Die
Go-Mutationen liefen an einer Kopie des Baums im Scratchpad (`git archive HEAD`, Ersetzung durch ein
Python-Skript auf der Kopie, Rücksetzen per `cp` vom unmutierten Text); die eine Mutation der E2E-Kette
(M-H) lief am Repo selbst und wurde per `git checkout` zurückgenommen (`git status --short` danach leer, bis auf
diesen Report).

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — [`AGENTS.md`](../../AGENTS.md) §3.9)

Jeder Lauf schrieb in eine Log-Datei; der Exit-Code wurde im selben Aufruf gesondert gesichert und danach gelesen.
Stand aller Läufe ohne Mutation: `HEAD` = `17d0048a`, Arbeitsbaum sauber. Je ein schwerer Docker-Lauf zugleich
(`free -m` vor dem Integrationslauf: 13,0 GB verfügbar; ein fremder Container mit Zufallsnamen lief während der Läufe,
nicht von mir).

| Sensor | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `make test` (Race-Detector) | **Exit 0** | 44 Zeilen `ok`, keine Zeile `FAIL`; darunter `ok …/usecase/backfill`, `ok …/internal/bootstrap`, `ok …/internal/domain/model` |
| `make test-store` | **Exit 0** | `ok …/internal/bootstrap 5.143s`, `DB-Adapter-Coverage: 82.56% (gedeckt 885 von 1072 Statements; Profile gemergt: store,replication)`, `db-coverage: OK — DB-Adapter-Coverage 82.56% erfuellt Schwelle 80%` |
| `make a-check` | **Exit 0** | `gesamt: 0 Befund(e)` |
| `make coverage-gate` | **Exit 0** | `coverage-gate: OK — Coverage 85.10% erfüllt Schwelle 80%` |
| `make gates` (einmal) | **Exit 0** | `baseline-verify: v6.9.0 OK — 54 Dateien` · `d-check: 1242 Datei(en) geprüft, 0 Befund(e)` · `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` · `generated-sync: OK — das committete Erzeugnis ist byte-gleich der Ausgabe des gepinnten Generators (Stufe proto-export)` · `coverage-gate: OK — Coverage 85.00% erfüllt Schwelle 80%` · `gesamt: 0 Befund(e)` (a-check) |
| `make suchlauf-nachmessen PLAN=<Plan des Slice>` | **Exit 0** | `suchlauf-nachmessen: 24 Zeilen stimmen` |
| `make commit-traceability RANGE=3973390e..HEAD` | **Exit 0** | `OK — 9 Commit(s) in "3973390e..HEAD", Betreffs ohne Struktur-ID` |
| `make doc-commits RANGE=3973390e..HEAD` | **Exit 0** | `d-check: 1242 Datei(en) geprüft, 0 Befund(e)` |
| `make doc-immutable RANGE=3973390e..HEAD` | **Exit 0** | `d-check: 1242 Datei(en) geprüft, 0 Befund(e)` (ohne `RANGE` bricht das Ziel mit `flag needs an argument: --range` ab — Bedienung, kein Befund) |
| gofmt (Docker, gepinntes Toolchain-Image, `gofmt -l` über die `*.go`-Dateien des Diffs) | **Exit 0** | keine Ausgabe |
| `make image`, dann `make test-integration` (ein Lauf) | **beide Exit 0** | Dauer 11:33:38 bis 11:38:52 = 5 min 14 s (gemessen, `make image` mit eingerechnet gegen den Layer-Cache); Zeilen unten |
| Mutationen (§4) | 7 Go-Mutationen (M-A bis M-G) und 1 E2E-Mutation (M-H) an der Eingabeseite | siehe §4 |

**Coverage-Zahl als Lauf-Beleg mit Streuung ([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz A).** Meine zwei Läufe im
selben Baum drucken 85.10 % (`make coverage-gate` allein) und 85.00 % (Schritt `coverage-gate` in `make gates`);
der Implementer nannte 85,10 %, der Reviewer 85.00 % in zwei Läufen. Die Zahl ist der Beleg eines Laufs, nicht ein
Ist-Stand; die beobachtete Spanne über vier bis fünf Läufe ist 85,00–85,10 %, der Abstand zur Schwelle (80 %) ist
fünf Punkte. Der Plan trägt die Zahl nicht (Diff gelesen).

**Integrationslauf (gedruckt).** „Backfill-Regelstand ([`LH-FA-CFG-007`](../../spec/lastenheft.md)) belegt — Run … übernahm 3 Zeilen von
feed_e2e_backfill_rule mit umbenanntem Schlüssel customer_name (Werte RegelAlpha,RegelBeta, Schlüsselmenge
customer_name,id,note gleich der WAL-Change), lesbar über cdc.changes und GET /changes; Run … nach exclude_column auf
name trägt in 4 Bildern weder name noch customer_name noch einen Wert von name; Run … auf feed_e2e_backfill_rulebad
endete failed (schema: Regel "kundenname" (Spalte "name") an public.feed_e2e_backfill_rulebad auf die Spalten des
Snapshots nicht anwendbar: Zielname kollidiert mit einer Spalte der Änderung: customer_name) ohne Change, der
Erfassungspfad lief weiter, der neue Antrag … nach dem Entfernen der Regel übernahm 2 Zeilen in Rohform“; danach
„E2E-Abdeckungstabelle unverändert — docs/user/e2e-abdeckung.md entspricht dem Quelltext-Stand“ und „Lauf
abgeschlossen — E2E-Abdeckungstabelle aus 14 Go-Zeilen und 36 Bash-Zeilen“. `git status --short` danach leer.

Hygiene: dangling Volumes (`docker volume ls -qf dangling=true | wc -l`) vor dem ersten Lauf **34**, nach dem
Integrationslauf **34**, nach der E2E-Mutation und dem Rücksetzen des Images **34**; kein `prune`, kein
`system prune`. Nicht gefahren: `make test-replication` (der Diff berührt den Replication-Tier nicht; der
Paritätstest liegt in `make test`, Plan §3 Zeile „Paritätstest“ am Start gelesen), `make bench`, eine Messung der
Kosten der Lesung gegen eine Queue mit vielen Zeilen (V-1).

## 2. DoD — Verdikt je Zeile (§2 des Plans)

Gezählt am Plan: acht `[x]`-Zeilen, vier `[ ]`-Zeilen, zusammen zwölf (`grep -c '^- \[x\]'` 8,
`grep -c '^- \[ \]'` 4).

| # | DoD-Zeile | Verdikt | Realer Beleg |
|---|---|---|---|
| 1 | Regelauswertung im Run: je Block Regelstand neu, Bild über die gemeinsame Funktion mit Regelsatz; byte-gleiches Bild wie WAL-Change; Paritätstest über die Regeltypen der Domäne; Aufrufer mit leerer Regelmenge ersetzt | **erfüllt** | `make test` Exit 0. `blockBuilder.build` ruft `model.BuildRowImage(b.columns, row, excluded, rules)` (Diff gelesen); Paritätstest `TestBackfillAndWALImagesAreByteEqualWithRules` in `internal/bootstrap` zählt `model.TransformationKinds()` auf (`git grep -n TransformationKinds`: Zeilen 112 der Paritätsdatei, 70 und 111 der Regel-Testdatei des Runs); Suchlauf `BuildRowImage\(.*nil\)` ohne Tests: Parent 1, Arbeitsbaum 0. M-A (`nil` statt `rules`) rot, dazu der Paritätstest rot. `make test-replication` nicht gefahren: der Plan bindet es nur an einen Paritätstest im Replication-Tier; der Ort ist im Plan §3 am Start gelesen und liegt in `make test` (`TestImageParityWalAndBackfill` im Replication-Tier bleibt regelunabhängig, ein Aufruf mit `nil` als Regelsatz, Suchlauf: zwei Treffer in `snapshot_test.go`) |
| 2 | Fail-closed und Nichtanwendbarkeit: Regelstand-Vergleich als Menge, Nichtanwendbarkeit endet `failed`/`schema` einmal je Run vor der Schreibtransaktion mit der Prüffunktion der Domäne; `classifyError` bildet `schema` und `configuration` ab; run-lokal; Eigenschaftstest (Regeltyp × ausgeschlossene Spalte); Lesefehler je Aufrufstelle | **erfüllt** | `make test` Exit 0. Code gelesen: Lesung nach `OpenSnapshot`, `progress` und vor der ersten `NextBlock`-Lesung und `Begin` (`checkRulesApplicable` gegen `builder.columns`), Vergleich je Block, Vergleich unmittelbar vor `Commit`, `classifyError`. Mutationen M-B (Vergleich vor dem Commit), M-C (Abbildung `schema`), M-D (Prüfung entfällt), M-E (Lesefehler zu Beginn), M-G (Mengenvergleich) je rot (§4). `TestExecuteRuleFailureIsReportedInTheRunOnly` und `TestPortsCarryNoCapturePathPort` tragen „run-lokal“ (der Use Case hält keinen Heartbeat-Port) |
| 3 | E2E-Beleg in `make test-integration`: `rename_column` im Backfill (über `cdc.changes` und `GET /changes`, `origin = 'backfill'`), `exclude_column` auf der umbenannten Spalte (`LH-QA-SEC-004`), nicht anwendbare Regel → `failed`/`schema` ohne Change, Erfassungspfad läuft, nach Abhilfe neuer Run `completed`; Zeile im Runner-Erzeugnis | **erfüllt** | eigener Lauf Exit 0 (§1); die gedruckte Zeile trägt alle vier Teile; die Erwartung von [`ADR-0117`](../plan/adr/0117-backfill-run-fehlerklasse-schema.md) Festlegung 4 („ohne Prozessneustart“, dort *erwartet*) ist damit erprobt: die Phase startet den Feed-Container nicht neu (Diff der Phase gelesen: kein `docker restart`). Die Zeile im Erzeugnis steht in `docs/user/e2e-abdeckung.md` (Zeile 62, trägt `LH-FA-CFG-007`, `LH-FA-CAP-009`, `LH-QA-SEC-004`) |
| 4 | `make gates` grün, Exit gesondert | **erfüllt** | eigener Lauf Exit 0 (§1) |
| 5 | Review durchgeführt, Report liegt vor | **erfüllt** | Report liegt vor (0 HIGH, 0 MEDIUM, 0 LOW, 8 INFO; nicht merge-blockierend); Finding für Finding in §5 nachgemessen — kein offenes HIGH/MEDIUM |
| 6 | §3.13-Suchlauf: committetes Feld in §3, Gefundenes und Nichtgefundenes je Träger, beide Stände | **erfüllt** | 24 Zeilen mit dem Werkzeug Exit 0 (§1); neun Paare von Hand nachgefahren (§6); die Befund-Zellen stimmen |
| 7 | Doku-Update: entfällt, Aufschub mit Adresse `slice-transformationen-betriebsdoku` §2 | **erfüllt (Aufschub mit Adresse; die Adresse trägt den Gegenstand nur teilweise, V-4)** | Handbuch-Kandidatenlauf §7 |
| 8 | Reconciliation-Register — entfällt (Greenfield) | **erfüllt (entfällt)** | keine Datei |
| 9 | Closure-Notiz mit Steering-Loop-Lerneintrag | **korrekt offen** | Plan §7 trägt „*(zu tragen bei Closure)*“ (Planner) |
| 10 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | Closure-Pflicht (Planner); der Diff berührt genau eine Register-Datei (V-6) |
| 11 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | Plan §6: alle sieben Zeilen tragen „Ausgang: *(bei Closure …)*“ (Planner) |
| 12 | Die drei Paarungen | **korrekt offen** | hängen an der Closure der Welle |

Kein `[x]` ohne Beleg; kein `[ ]`, das über die Rollen-Sequenz hinaus belegt wäre. Register, Risiko-Ausgänge,
Paarungen und Closure-Notiz sind Planner-Arbeit und nicht Teil dieser Prüfung.

## 3. Plan-vs-Code-Diff

**§1 Ziel:** der Run liest den Regelstand über den neunten Port, baut das Bild mit dem Regelsatz, prüft den Stand
fail-closed, die Nichtanwendbarkeit endet den Run `schema`, die E2E-Phase steht. Die „Ausdrücklich NICHT“-Punkte
sind eingehalten: `git diff --stat 3973390e..HEAD` nennt weder `internal/adapters/**` noch
`internal/application/usecase/` außerhalb des Pakets `backfill` noch `spec/` noch `docs/plan/adr/` noch
`test/integration/`; ein zweiter Bild-Bau-Weg im Run besteht nicht (Suchlauf Zeile 1–2: `BuildRowImage(` 4/4).

**§3-Tabelle, Zeile für Zeile gegen `git diff --stat`:**

| Plan-Zeile | Ist im Diff |
|---|---|
| `service.go` update | Port `Transformations`, Lesung zu Beginn, `checkRulesApplicable`, Vergleich je Block und vor dem Commit, `sameSet`, `classifyError`, Doc-Kommentare (Diff gelesen) |
| `errors.go` update (Plan-Drift benannt) | neuer Sentinel `ErrTransformationStateChanged`; der Plan nennt ihn als Plan-Drift |
| `rowimage.go` update (Kommentar) | ein Satz |
| `wiring.go` update | eine Zeile `Transformations: activation` |
| zwei Test-Verdrahtungen | `administration_roles_internal_test.go`, `backfill_endtoend_test.go` |
| `service_test.go` update, `transformation_test.go` neu | 76 bzw. 464 Zeilen; acht Testfunktionen in der neuen Datei |
| `backfill_image_parity_test.go` neu | 206 Zeilen, tabellengetrieben über `model.TransformationKinds()` |
| `run-integration-tests.sh` update | Phase „Backfill-Regelstand“ nach „Backfill-Boundary“, Kopfkommentar nachgezogen |
| `harness/README.md` update | Zeile `make test-integration`: „acht Backfill-Rundläufe“ |
| `docs/user/e2e-abdeckung.md` Erzeugnis | 71 Zeilen (Zeilen-Lokatoren) mit einer hinzugefügten Zeile der Phase |

Nicht im Plan-Feld, im Diff: die Register-Datei `state.md` von `BEO-PGC/run-fehlerklasse-schema-im-transformations-backfill`
(eine Zeile, Commit `13aaf752` mit eigener Begründung im Betreff; Reviewer F-7). Im Plan, nicht im Diff: nichts. Keine
`Accepted` ADR wurde inhaltlich geändert (`make doc-immutable RANGE=3973390e..HEAD` Exit 0). Keine Schwelle, keine
Gate-Konfiguration verändert ([`AGENTS.md`](../../AGENTS.md) §3.5, §3.6). Die beiden Lifecycle-Commits sind reine
Renames (`git show -M --stat` von `42ea6d7c` und `24b1d7d9`: je ein Rename, 0 Zeilen; [`AGENTS.md`](../../AGENTS.md) §3.3).

**Festlegungen des Implementers gegen den Code:**

| Festlegung | Ist |
|---|---|
| Stand des Runs ist der Stand nach dem Öffnen des Snapshots | `transformationRules` nach `OpenSnapshot` und `progress`, vor der Schleife; Vergleichsstand `baselineRules` |
| Kosten der Lesung: Regelstand `Blockzahl + 2`, Ausschlussstand `Blockzahl + 1` | am Code nachgezählt: eine Lesung zu Beginn, je nicht leerem Block eine, eine vor dem Commit; die Formel gilt ab einem Block — bei einer leeren Tabelle endet der Run ohne `Begin` und liest den Regelstand einmal, den Ausschlussstand nie (V-1) |
| Faltungsfehler im Run: Klasse `internal` | Rückfall von `classifyError`; M-E und die Reviewer-Mutation M18 rot |
| leere Tabelle: Prüfung läuft auch dann | die Prüfung liegt vor der Schleife, hängt nicht an einer Zeile; M-D rot in `TestExecuteInapplicableRuleEndsRunAsSchema` (Fall leere Tabelle) |

## 4. Mutationen der Eingabeseite (dieser Lauf)

M-A bis M-G an einer Kopie des Baums, Einzellauf `go test -race` im gepinnten Toolchain-Image ohne Netz über
`./internal/application/usecase/backfill/` und `./internal/bootstrap/` (Ausgangslauf der unmutierten Kopie: beide
`ok`). Die Ersetzung setzt genau eine Fundstelle (das Skript bricht bei einer anderen Trefferzahl als eins ab). M-H läuft
als `make image` und `make test-integration` am mutierten Repo.

| # | Zusage | Mutation | Ergebnis (gedruckt) |
|---|---|---|---|
| M-A | Der Run baut das Bild mit dem Regelsatz des Blocks | `nil` statt `rules` an `BuildRowImage` (`build`) | **rot**: `TestExecuteBuildsImagesWithTheRuleSet`, `TestExecuteRulesNeverLeakExcludedColumns`, `TestExecuteRuleTargetCollisionInBlockEndsRunAsSchema`; in `internal/bootstrap`: `TestBackfillAndWALImagesAreByteEqualWithRules` |
| M-B | Der Regelstand wird unmittelbar vor dem Commit verglichen | der Vergleich vor `run.Complete` entfällt (`_ = rules`) | **rot**: `TestExecuteRuleStateChangeEndsRunAsConfiguration` |
| M-C | Nichtanwendbarkeit (Spalte fehlt) endet `schema` | Abbildung von `ErrTransformationColumnMissing` aus `classifyError` entfernt | **rot**: `TestExecuteInapplicableRuleEndsRunAsSchema` |
| M-D | Nichtanwendbarkeit wird vor der ersten Zeile geprüft | `checkRulesApplicable(nil, …)` statt der Regeln | **rot**: `TestExecuteInapplicableRuleEndsRunAsSchema`, `TestExecuteRuleFailureIsReportedInTheRunOnly` |
| M-E | Lesefehler des Regelstands endet den Run (Stelle: zu Beginn) | Fehler der ersten Lesung verworfen (`baselineRules, _ :=`) | **rot**: `TestExecuteRuleReadFailure`, `TestExecuteUnreadableAppliedRowEndsRunAsInternal` |
| M-G | Der Regelstand wird als Menge verglichen | `sameSet` prüft nur die Mengengröße | **rot**: `TestExecuteRuleStateChangeEndsRunAsConfiguration` |
| M-F | Der Run erhält den Regelstand-Port aus der Verdrahtung | die Zeile `Transformations: activation` in `wiring.go` entfernt (`make test`-Pakete) | **grün**: `ok …/usecase/backfill (cached)`, `ok …/internal/bootstrap 1.128s` |
| M-H | dieselbe Zeile, E2E-Kette | dieselbe Mutation im Repo, `make image` (Exit 0), `make test-integration` | **Exit 2**: „Backfill-Happy-Path — Run … erreichte completed nicht innerhalb von 60s (status=running)“, `make: *** [Makefile:191: test-integration] Fehler 1`; die Mutation ist per `git checkout` zurückgenommen (`git status --short` leer), `make image` danach erneut Exit 0 |

Die Mutationen des Reviews (M1–M21, 20 rot, M21 = die Verdrahtungszeile grün) sind **übernommen**, nicht wiederholt;
meine sieben Go-Mutationen sind Stichproben derselben Zusagen an `HEAD`. M-F und M-H bestätigen die Reviewer-Beobachtung
F-4 (M21) und ergänzen die E2E-Seite (V-3).

## 5. Findings des Reviews nachgemessen (nicht dem Bericht geglaubt)

| Finding | Was ich gemessen habe | Verdikt |
|---|---|---|
| F-1 (INFO) Kosten der Lesung ohne Messung und Adresse | `git grep -n -E 'abgeleitet\|gemessen\|Größenordnung' -- internal/application/usecase/backfill/service.go`: 0 Treffer — der Doc-Kommentar von `copyBlocks` trägt die Kennzeichnung „abgeleitet, nicht gemessen“ **nicht** (die Behauptung des Implementers trifft wörtlich nicht zu); die Kennzeichnung steht im Plan §3 (Festlegung „Kosten der Lesung“). Der Auslöser im Plan („falls die Queue einer Quelle in die Größenordnung der Blockzahl wächst“) trägt keine Adresse (kein Slice, kein Register-Eintrag; §6 führt kein Kostenrisiko) | **bestätigt**, offen als V-1 |
| F-2 (INFO) Bezugspunkt des Regelstands im Run | ADR-Wortlaut gelesen ([`ADR-0117`](../plan/adr/0117-backfill-run-fehlerklasse-schema.md) Festlegung 2 „nachdem der Snapshot seine Spalten liefert“, Festlegung 5 „zwischen zwei Lesungen des Runs“, [`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) Folgepflicht 7 ohne Zeitpunkt); [`SPEC-002`](../../spec/pflichtenheft.md) sagt „Regelstand zum Erfassungszeitpunkt“ | **konform**, kein Widerspruch; Frage des Implementers bleibt beim Architect |
| F-3 (INFO) `sameSet[T comparable]` und `map_value` | `git grep -n 'sameSet\|comparable' -- docs/plan/planning/open/slice-transformationen-map-value.md`: 0 Treffer; der Plan verlangt in DoD Punkt 2 einen Diff ohne `internal/application/usecase/` und führt die Änderung dort als Rückführungs-Fall (`in-progress` → `open`); `model.Transformation` trägt vier Zeichenketten-Felder (`transformation.go`), `map_value` trägt ein Objekt `values` | **bestätigt**, offen als V-2 |
| F-4 (INFO) Verdrahtungszeile nicht in `make test` gebunden | M-F grün in `make test`; M-H rot in `make test-integration` (Backfill-Happy-Path, Run bleibt `running`); `git grep -n 'bootstrap.Run('` in Tests: nur `run_test.go` (endet am ersten Konstruktor ohne erreichbare Quelle) | **bestätigt**; die Bindung trägt allein der E2E, V-3 |
| F-5 (INFO) Prüfkern einmal, Iteration zweimal | `git grep -n -E 'CheckApplicable\(' -- internal ':!*_test.go'`: drei Treffer (Definition, `mapper.go`, `service.go`); `checkTransformations` (Mapper) und `checkRulesApplicable` (Run) tragen je eine eigene Schleife mit eigenem Fehlertext, beide rufen dieselbe Prüffunktion der Domäne | **konform** mit [`ADR-0117`](../plan/adr/0117-backfill-run-fehlerklasse-schema.md) Festlegung 2 und der Fitness-Function-Zeile „Review-Prüfpflicht“ (Prüffunktion an einer Stelle der Domäne, beide Pfade rufen sie); die Iteration ist eine benannte Grenze |
| F-6 (INFO) Coverage 85,10 % gegen 85.00 % | eigene Läufe 85.10 % und 85.00 % (§1) | **bestätigt** als Streuung, Lauf-Beleg |
| F-7 (INFO) Register-Reparatur | `git show 13aaf752`: eine Zeile in `state.md`, Link auf `open/…` durch die Kennung in Inline-Code ersetzt; Ausgang und Zähler des Eintrags unverändert | **konform**; die Closure schreibt das Register (Planner) |
| F-8 (INFO) `sed -i` auf `/dev/null` | Kenntnis; der Diff trägt keine Spur eines Textwerkzeugs (Diff gelesen) | **Kenntnis** |

Kein offenes HIGH, kein offenes MEDIUM.

## 6. Suchlauf-Feld (§3 des Plans), an beiden Ständen nachgefahren

Parent `3973390e`; Stand `diff` ist der Arbeitsbaum (`HEAD` = `17d0048a`, sauber). Neun der zwölf Messpaare von Hand
mit `git grep -n … | wc -l` (ohne das Werkzeug), die Plan-Datei ausgeschlossen (am Parent liegt sie unter `open/`, am
Arbeitsbaum unter `in-progress/`; mit dem Pathspec `':!docs/plan/planning/*/slice-transformationen-backfill-pfad.md'`
gemessen):

| Zeile | Stand | Soll | Von Hand |
|---|---|---|---|
| `BuildRowImage(` ohne Tests | `3973390e` / `diff` | 4 / 4 | **4** / **4** |
| `BuildRowImage\(.*nil\)` ohne Tests | `3973390e` / `diff` | 1 / 0 | **1** / **0** |
| `Ausschlussstand` (ohne `docs/reviews`, `done/`, Baseline) | `3973390e` / `diff` | 92 / 93 | **92** / **93** |
| `-i Backfill` in `docs/user harness spec` | `3973390e` / `diff` | 236 / 237 | **236** / **237** |
| `sieben Backfill` | `3973390e` / `diff` | 1 / 0 | **1** / **0** |
| `acht Backfill` | `3973390e` / `diff` | 0 / 1 | **0** / **1** |
| `slice-transformationen-backfill-pfad` im Code | `3973390e` / `diff` | 1 / 0 | **1** / **0** |
| `TransformationRules` ohne Tests | `3973390e` / `diff` | 13 / 15 | **13** / **15** |
| `CheckApplicable\(` ohne Tests | `3973390e` / `diff` | 2 / 3 | **2** / **3** |

Alle Zahlen stimmen. Die `Ausschlussstand`-Zeile misst am Parent 96 Treffer, solange der Pathspec nur den
`in-progress/`-Pfad ausschließt (vier Treffer in der Plan-Datei unter `open/`); mit dem Pfad-Muster über alle
Lifecycle-Verzeichnisse sind es 92 — die Zahl des Plans gilt für den Suchraum „ohne die Plan-Datei“. Die
Nichtgefunden-Aussagen habe ich an zwei Stellen selbst gesucht: kein zweiter Bild-Bau-Weg im Run (`BuildRowImage(`
4/4, der eine Aufrufer im Run trägt `rules`), keine Kopie der Prüffunktion im Run (`CheckApplicable(` drei Treffer,
davon ein Aufrufer je Pfad).

## 7. Entscheidungs-Konformität

- **[`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) Teilfrage 5 (Ausschluss gilt
  zuerst) und Folgepflicht 7 (ein Backfill-Pfad trägt dieselbe Regelauswertung):** der Run ruft dieselbe
  `BuildRowImage` mit Ausschluss- und Regelsatz; der Eigenschaftstest im Run (`TestExecuteRulesNeverLeakExcludedColumns`,
  Regeltyp × Regel-Spalte × ausgeschlossene Spalte) und die E2E-Phase (`exclude_column` auf der umbenannten Spalte:
  vier Bilder ohne Quellname, Zielname und Wert) tragen [`LH-QA-SEC-004`](../../spec/lastenheft.md) für den Backfill-Pfad.
  Konform.
- **[`ADR-0117`](../plan/adr/0117-backfill-run-fehlerklasse-schema.md):** Festlegung 1 (Klasse `schema` für beide
  Sentinels der Prüfung, `internal` als Rückfall) im Code und mit M-C rot; Festlegung 2 (dieselbe Prüffunktion der
  Domäne, einmal je Run, nach dem Öffnen des Snapshots und vor `Begin` und der ersten Zeile) im Code gelesen, M-D rot;
  Festlegung 3 (run-lokal, Sichtbarkeit) mit der E2E-Zeile („Lebenszeichen der Quelle ohne Fehlerzustand“, Feed-Container
  läuft weiter, Sonden-Zeile wird erfasst) erprobt; Festlegung 4 („ohne Prozessneustart“, dort *erwartet*) mit dem
  E2E-Lauf erprobt; Festlegung 5 (Klasse `configuration` für den Regelstand-Wechsel): der Text nennt die
  „`ErrExclusionStateChanged`-Abbildung“, der Diff führt einen eigenen Sentinel mit derselben Klasse — Plan-Drift
  im Plan benannt, kein Widerspruch zum Wortlaut der Klasse. Folgepflicht 2 (Prüfung, Sentinel, Abbildung, Kommentar
  „vergibt `schema` nicht“ entfällt: `git grep -n -F 'vergibt der Run nicht'` 2/0, Negativtests je Eingabe) und
  Folgepflicht 3 (E2E-Beleg) erfüllt; Verdikt `architect-verdict-backfill-schema-klasse-rollen` als Start-Vorbedingung
  im Übergangs-Commit `24b1d7d9` genannt. Folgepflicht 4 (Handbuch, Betriebshinweis zur Abhilfe) ist nicht Teil dieses
  Slice (V-4).
- **[`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md):** Teilfrage 2 (eine Funktion für WAL- und
  Backfill-Pfad): der Paritätstest bindet beide Pfade an denselben Regelsatz (M-A, Reviewer M2 und M20 rot); Teilfrage 4
  (Fail-closed vor dem Commit) um den Regelstand erweitert, Bindung, Ausschluss- und Regelstand vor dem Commit in
  `copyBlocks` gelesen; Teilfrage 5 (run-lokal) siehe oben. Konform.
- **[`ADR-0115`](../plan/adr/0115-backfill-spaltenwerte-text-ergebnisformat.md) Festlegung 4 (Byte-Gleichheit an der
  Bild-Konstruktion):** der Paritätstest vergleicht die Bytes am Ausgang der gemeinsamen Funktion; `LH-FA-CAP-009.a`
  sagt am `jsonb`-Lesepfad nur Schlüsselmenge und Werte zu — die E2E-Phase misst diese Hälfte (Schlüsselmenge der
  Backfill-Change gleich der WAL-Change, gedruckt: `customer_name,id,note`). Konform.
- **[`SPEC-030`](../../spec/pflichtenheft.md) (Anwendbarkeit), [`SPEC-019`](../../spec/pflichtenheft.md),
  [`SPEC-008`](../../spec/pflichtenheft.md) (Zeile `schema`: „im Erfassungspfad und im Run“),
  [`LH-FA-CAP-009.a`](../../spec/pflichtenheft.md) („bevor die erste Zeile gelesen wird“, ohne Change, ohne den
  Erfassungspfad zu berühren):** am Code und am E2E-Lauf gegengelesen; kein Widerspruch.
- **Die drei Fragen des Implementers** (Bewertung des Reviewers nachgeprüft): (i) Lesekosten je Block — am Code
  nachgezählt (`Blockzahl + 2` und `Blockzahl + 1` ab einem Block), nicht gemessen; benannte Grenze mit dem Nachtrag V-1.
  (ii) Prüfung auch bei leerer Tabelle — konform mit Festlegung 2 („einmal je Run“ meint die Häufigkeit, nicht die
  Bedingung einer Zeile); M-D rot im Fall der leeren Tabelle. (iii) Regelstand ab dem Öffnen des Snapshots —
  konform (Wortlaut schweigt zum Antragszeitpunkt), Grenze: ein zwischen zwei Lesungen gesetzter und zurückgenommener
  Stand bleibt unsichtbar und steht im Doc-Kommentar von `copyBlocks`; die Grenze ändert kein gebautes Bild.
- **[`AGENTS.md`](../../AGENTS.md):** §3.3 (Moves rein), §3.5/§3.6 (keine `Accepted` ADR überschrieben, keine Schwelle
  gesenkt), §3.7 (Kommentar-Zusagen der neuen Zweige an den Code gelesen: „auch bei einer leeren Tabelle“, „bevor eine
  Zeile gelesen … wird“ mit M-D; die Formel `Blockzahl plus zwei` gilt ab einem Block, V-1), §3.9 (meine Läufe ungepiped),
  §3.12 (Zahlen des Plans tragen Herkunft; Coverage als Lauf-Beleg), §3.13 (Suchlauf beide Stände; die Meldung
  „Regeln gelten auch für einen Backfill“ an `slice-transformationen-betriebsdoku` steht dort in §2 als „Backfill-Bezug“).
  §3.2: kein `//nolint` im Diff (`git grep -n nolint` über die geänderten Go-Dateien: 0 Treffer). §3.10: kein
  Workflow im Diff.
- **Commits:** `make doc-commits` und `make commit-traceability` über die Range Exit 0; jeder der neun Betreffs nennt
  `LH-FA-CFG-007`, keiner trägt `SPEC-`/`ARC-` im Betreff.
- **Handbuch-Kandidatenlauf** (`git diff --name-only 24b1d7d9 -- internal/bootstrap/ tools/schema/
  internal/adapters/driving/`): vier Dateien, alle in `internal/bootstrap/` (drei Tests, `wiring.go`); der Diff von
  `wiring.go` trägt keine `CDC_*`-Variable, keine SQL-Funktion und keinen Endpunkt — **keine neue
  Betreiber-Oberfläche**; `docs/user/benutzerhandbuch.md` liegt nicht im Diff. Die Aussage „Regeln gelten auch für den
  Backfill“ nennt `slice-transformationen-betriebsdoku` §2 mit dem Wort „Backfill-Bezug“ (V-4).
- **Kenntnis (Prozess):** Implementer und Reviewer meldeten je einen wirkungslosen `sed -i` auf `/dev/null` (Review F-8);
  in meinem eigenen Lauf trat derselbe Fehlgriff einmal auf einer Scratch-Skriptdatei auf (nicht das Repo, ohne Wirkung
  auf den Baum). Die Mutationen liefen über Python-Ersetzungen auf der Kopie, am Repo selbst nur M-H per
  `git checkout` zurückgenommen.

## 8. Findings dieser Verifikation

| # | Kategorie | Befund | Quelle | Verifizierbar |
|---|---|---|---|---|
| V-1 | LOW | **Kosten der Lesung: Kommentar ohne Kennzeichnung, Auslöser ohne Adresse, Formel nur ab einem Block.** Der Doc-Kommentar von `copyBlocks` nennt die Zahl der Lesungen („Blockzahl plus eine“, „Blockzahl plus zwei“) und die Ausdehnung auf alle Tabellen der Quelle, aber weder „abgeleitet“ noch „nicht gemessen“ (0 Treffer für `abgeleitet\|gemessen`); die Behauptung des Implementers, der Kommentar trage die Kennzeichnung, trifft wörtlich nicht zu — die Kennzeichnung steht im Plan §3. Der Auslöser des Plans („ein tabellenbezogener Lesezugriff … gehört in einen eigenen Plan, falls die Queue einer Quelle in die Größenordnung der Blockzahl wächst“) trägt keine Adresse (kein Slice, kein Register-Eintrag); ohne Adresse liest ihn niemand, der die Queue wachsen sieht. Die Formel `Blockzahl + 2` gilt ab einem Block: bei einer leeren Tabelle liest der Run den Regelstand einmal und den Ausschlussstand nie. Keine der Zahlen ist falsch als Strukturaussage; sie sind an Code und Test (fünf Lesungen bei drei Blöcken) gebunden. **Schwere:** LOW — die Zusage ist keine Größenaussage, der Ausgang ist eine Grenze ohne Träger | `internal/application/usecase/backfill/service.go` (Doc-Kommentar `copyBlocks`); Plan §3 („Kosten der Lesung“) | ja — `git grep -n -E 'abgeleitet\|gemessen' -- internal/application/usecase/backfill/service.go` |
| V-2 | LOW | **Träger-Nachzug ohne Adresse: `sameSet[T comparable]` und der Regeltyp `map_value`.** `sameSet` verlangt einen vergleichbaren Typ und wird mit `model.Transformation` (vier Zeichenketten-Felder) instanziiert; der Kommentar sagt, der Typ sei über alle Felder vergleichbar. Ein Regeltyp mit einem Objekt `values` macht den Typ je nach Ablage unvergleichbar (Map, Slice); der Bau bricht dann am Übersetzer, sichtbar, nicht still. Der Plan `slice-transformationen-map-value` nennt `sameSet` nicht (0 Treffer) und verlangt in DoD Punkt 2 einen Diff ohne `internal/application/usecase/`, in §4 die Rückführung `in-progress` → `open`, falls der Wirkort ändern muss. Die Meldung des Reviewers (F-3) hat damit keinen Träger im Plan, an den sie gerichtet ist. **Meldung an den Planner:** `slice-transformationen-map-value` §3 um die Fundstelle ergänzen (`sameSet`, Doc-Kommentar „vergleichbar“) und den Zwei-Wege-Ausweg festhalten (Vergleichbarkeit der Regel erhalten oder den Vergleich in der Domäne führen); die Ausnahme in DoD Punkt 2 benennen. Frist: der Start dieses Slice | `internal/application/usecase/backfill/service.go` (`sameSet`); Plan `slice-transformationen-map-value` §2 (DoD Punkt 2), §3, §4 | ja — Übersetzer und `git grep -n sameSet -- docs/plan/planning` |
| V-3 | INFO | **Verdrahtungszeile `Transformations: activation` ist allein durch den E2E gebunden.** M-F: `make test` grün; M-H: `make test-integration` Exit 2 an der ersten Backfill-Phase (Happy Path, Run bleibt `running` bis zum Zeitlimit). Kein Test in `make test`/`make test-store` erreicht die Backfill-Verdrahtung von `Run` (`run_test.go` endet am ersten Konstruktor ohne erreichbare Quelle; die drei Tests mit realem Adapter bauen ihre `backfill.Ports` selbst). Die Zeile `Exclusion: activation` daneben hat dieselbe Form und dieselbe Bindung. Der Ausfall ist laut (jeder Backfill-Run hängt), nicht still. **Bewertung:** benannte Grenze, kein Test-Nachtrag verlangt — ein Test der Verdrahtung bräuchte ein `Run` gegen reale Datenbank, und die Tier-Aussage („belegt am laufenden System“) ist der Gegenstand des E2E-Beleges. Der Plan benennt die Grenze nicht; der Planner nimmt sie in die Closure-Notiz auf | `internal/bootstrap/wiring.go` (`backfill.Ports`); `internal/bootstrap/run_test.go` | ja — M-F, M-H |
| V-4 | LOW | **Handbuch-Träger `slice-transformationen-betriebsdoku` trägt die Run-Abhilfe nur dem Wort nach.** Die Meldung „Regeln gelten auch für einen Backfill“ steht als „Backfill-Bezug“ in §2 (Adresse vorhanden). [`ADR-0117`](../plan/adr/0117-backfill-run-fehlerklasse-schema.md) Folgepflicht 4 verlangt einen Betriebshinweis zur Abhilfe im Run: Regel entfernen, danach ein **neuer** `cdc.backfill_table`-Antrag, **kein** Prozessneustart — der E2E belegt genau diesen Ablauf. Die Abhilfe-Zeile des Betriebsdoku-Plans nennt für die Nichtanwendbarkeit „Prozess starten — so, wie `e2e-abhilfe` sie belegt hat“ (das ist der Erfassungspfad); die Kennung der ADR und der neue Antrag stehen dort nicht (`git grep -n '0117\|neuer Antrag' -- docs/plan/planning/open/slice-transformationen-betriebsdoku.md`: 0 Treffer). **Meldung an den Planner** (Träger in fremder Datei, Frist: Closure dieses Slice): den Run-Ablauf der Abhilfe in §2 von `slice-transformationen-betriebsdoku` benennen | `docs/plan/planning/open/slice-transformationen-betriebsdoku.md` §2; [`ADR-0117`](../plan/adr/0117-backfill-run-fehlerklasse-schema.md) Folgepflicht 4 | ja — `git grep` wie genannt |
| V-5 | INFO | **Coverage-Zahl ist ein Lauf-Beleg mit Streuung.** 85.10 % (`make coverage-gate` allein) und 85.00 % (im selben Baum, in `make gates`); Implementer 85,10 %, Reviewer 85.00 %. Spanne 85,00–85,10 %, Schwelle 80 %. Der Plan trägt die Zahl nicht | §1 | ja — `make coverage-gate` |
| V-6 | INFO | **Register-Datei im Diff.** `BEO-PGC/run-fehlerklasse-schema-im-transformations-backfill/state.md` trägt seit `13aaf752` die Kennung des Slice in Inline-Code statt eines Links auf `open/…` (eine Zeile); der Ausgang und der Zähler des Eintrags sind unverändert, die Fortschreibung schreibt die Closure. Der Plan §3 führt die Datei nicht; die Commit-Message nennt Grund und Umfang | `git show 13aaf752` | ja — `make docs-check` |
| V-7 | INFO | **Übernommen, nicht nachgemessen:** die 21 Mutationen des Reviews (M1–M21), soweit nicht in §4 wiederholt; die Zeitangabe 5 min 23 s des Reviewer-Integrationslaufs (meine eigene Messung: 5 min 14 s); die Nullprobe der Kosten der Lesung gegen eine Queue mit vielen Zeilen (nicht gefahren, V-1); `make test-replication` (nicht Teil des Diffs) | §1, §4 | nein |

Kein HIGH, kein MEDIUM. Keine DoD-Verletzung.

## 9. Verdikt

**DoD bestätigt:** ja — jede der acht `[x]`-Zeilen (Nr. 1–8) ist am Ist-Zustand belegt; die vier `[ ]`-Zeilen
(Nr. 9–12) sind korrekt offen (Planner-Closure). **Plan-vs-Code:** keine unbenannte Abweichung im Code; die
einzige nicht im Plan-Feld stehende Änderung ist die Register-Zeile (V-6, mit Grund im Commit). **Entscheidungs-Konformität:**
[`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) (Teilfrage 5, Folgepflicht 7),
[`ADR-0117`](../plan/adr/0117-backfill-run-fehlerklasse-schema.md) (Festlegung 1 bis 5, Folgepflicht 2 und 3),
[`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) (Teilfrage 2, 4, 5),
[`ADR-0115`](../plan/adr/0115-backfill-spaltenwerte-text-ergebnisformat.md) (Festlegung 4) und
[`SPEC-030`](../../spec/pflichtenheft.md)/[`SPEC-019`](../../spec/pflichtenheft.md)/[`SPEC-008`](../../spec/pflichtenheft.md)
konform; Folgepflicht 4 von [`ADR-0117`](../plan/adr/0117-backfill-run-fehlerklasse-schema.md) (Handbuch) ist an
`slice-transformationen-betriebsdoku` übergeben und dort im Detail nicht getragen (V-4). **Review-Findings:** F-1 bis F-8
gemessen — F-5 konform, F-2 konform, F-6 bestätigt als Streuung, F-7 konform, F-8 Kenntnis, F-1/F-3/F-4 bestätigt und als
V-1, V-2, V-3 mit Adresse geführt; kein offenes HIGH/MEDIUM. **Gates:** `make test`, `make test-store`, `make a-check`,
`make coverage-gate` (85.10 %, im Gate-Lauf 85.00 %), `make gates`, `make suchlauf-nachmessen` (24 Zeilen),
`make doc-commits`, `make doc-immutable`, `make commit-traceability`, `gofmt -l` (keine Ausgabe) und der reale
`make image` + `make test-integration`-Lauf (Phase „Backfill-Regelstand“ belegt) im eigenen Lauf grün; sechs von
sieben Go-Eingabeseiten-Mutationen rot, die siebte (Verdrahtungszeile) in `make test` grün und im E2E rot.

**Übergabe:** an den Planner — Closure-Notiz mit Lerneintrag (Klassen aus Review und diesem Report: „Kosten der Lesung
ohne Messung und Adresse“, „Träger-Nachzug: Vergleichbarkeit von `Transformation`“, „Verdrahtungszeile ohne Test in
`make test`“), Beobachtungs-Register, Ausgänge der §6-Risiken (das Fenster zwischen Regel-Setzbarkeit und
Backfill-Bindung entfällt mit der Closure dieses Slice: der E2E-Beleg trägt den Bestand mit umbenanntem Schlüssel;
die Laufzeit von `make test-integration` trägt den Lauf 5 min 14 s, gemessen 2026-09-26), Paarungen bei der Closure
der Welle; Meldungen V-2 (`slice-transformationen-map-value`, Frist: dessen Start) und V-4
(`slice-transformationen-betriebsdoku` §2, Frist: Closure dieses Slice); V-1 mit Adresse für den Auslöser
(Slice oder Register-Eintrag) und der Kennzeichnung im Kommentar (Zusatz an den Implementer, keine Fixrunde
nötig, wenn der Planner die Grenze mit Adresse führt); V-3, V-5 bis V-7 als Anmerkungen. Dieser Report ist ein
**Lauf-Beleg** (dieser Stand, dieser Lauf) und ersetzt weder Review noch Closure.
