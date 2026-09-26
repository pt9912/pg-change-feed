# Verifikations-Report: slice-harness-suchlauf-nachmessen — 2026-09-26

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich + Entscheidungs-Konformität +
Plan-vs-Code-Diff + Gates. Review-Artefakt des Reviewers:
[`review-slice-harness-suchlauf-nachmessen.md`](review-slice-harness-suchlauf-nachmessen.md);
Formvorbild dieses Reports: [`verifikation-slice-backfill-bench-richtgroesse.md`](verifikation-slice-backfill-bench-richtgroesse.md).
Die vendored Baseline trägt kein eigenes Verifikations-Template
(`.harness/baseline/v6.9.0/templates/docs/reviews/` enthält nur `review-report.template.md`); der
Report folgt dem Formvorbild.

**Gegenstand:** Slice-Plan `slice-harness-suchlauf-nachmessen` (wellenlos), Stand `HEAD` = `b1a6d76d`,
Diff-Range `71ff3e44..HEAD`, 9 Commits, 24 Dateien. Slice-Inhalt: drei Lifecycle-Commits (`42e5ebb8`,
`3c2bf16d`, `8717c4fb`), Bau (`b5d2e12d`), Plan-Nachzug und DoD-Haken (`40e22e09`, `fc0f5629`), Fixrunde
(`a124da45`, `b1a6d76d`). Nicht Slice-Inhalt, im Range enthalten: der Review-Report des Reviewers
(`89bb8a45`). Der Review-Report bezieht sich auf die Range bis `fc0f5629`; die Fixrunde ist damit von
keinem Reviewer gelesen — ihre Wirkung belegen die eigenen Mutationen in §4. Dieser Lauf ändert weder
Code noch Plan noch Doku; er schreibt nur diesen Report. Alle Mutationen liefen an Kopien der Skripte
in einem Scratch-Verzeichnis des Laufs (`sed` auf stdout, kein `sed -i`, Prüfling per `TOOL=<Kopie>`);
die Originale blieben unberührt, `git status --short` war nach jedem Lauf leer.

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — [`AGENTS.md`](../../AGENTS.md) §3.9)

Jeder Lauf schrieb in eine Log-Datei; der Exit-Code wurde im selben Aufruf gesondert gesichert, die
Logs danach gelesen. Stand aller Läufe ohne Mutation: `HEAD` = `b1a6d76d`, Arbeitsbaum sauber.

| Sensor | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `make gates` | **EXIT=0** | `baseline-verify: v6.9.0 OK — 54 Dateien` · `coverage-gate: OK — Coverage 83.40% erfüllt Schwelle 80%` · `d-check: 1199 Datei(en) geprüft, 0 Befund(e)` · `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` · `generated-sync: OK — das committete Erzeugnis ist byte-gleich …` · a-check `gesamt: 0 Befund(e)` |
| `make test-suchlauf-nachmessen` | **EXIT=0** | `run-suchlauf-nachmessen-tests: alle Fälle bestanden` |
| `make test-rollout-restore` | **EXIT=0** | `run-rollout-restore-tests: alle Fälle bestanden` |
| `make suchlauf-nachmessen PLAN=<Plan des Slice>` | **EXIT=0** | acht Zeilen `OK`, `suchlauf-nachmessen: 8 Zeilen stimmen` (soll = ist: 14 / 80 / 9 / 12 / 59 / 74 / 35 / 47) |
| `make test-store` | **EXIT=0** | `db-coverage: OK — DB-Adapter-Coverage 82.61% erfuellt Schwelle 80%`; danach `git status --short tools/schema` leer, `git status --short` leer |
| `make test-replication` | **EXIT=0** | `--- PASS: TestWALRetentionThresholdsFollowGrowthAtInactiveSlot`, `ok … internal/bootstrap`; danach `git status --short tools/schema` leer, `git status --short` leer |
| `make doc-commits RANGE=71ff3e44..HEAD` | **EXIT=0** | `d-check: 1199 Datei(en) geprüft, 0 Befund(e)` |
| `make doc-immutable RANGE=71ff3e44..HEAD` | **EXIT=0** | `d-check: 1199 Datei(en) geprüft, 0 Befund(e)` (ohne `RANGE` bricht das Ziel mit `flag needs an argument: --range` ab — Bedienung, kein Befund) |
| Mutationen (§4) | sieben Läufe des Zerleger-Tests rot, ein Lauf des Wrapper-Tests rot | siehe §4 |

Hygiene: dangling Volumes (`docker volume ls -q -f dangling=true | wc -l`) am Ende **34**, kein
`prune`; ein schwerer Docker-Lauf zugleich; Speicher vor `make test-store` 16,7 GB verfügbar
(`free -m`), am Ende 17,6 GB. Nicht gefahren: `make test-integration`, `make test-sdk-*-integration`,
`make bench`, `make example-demo-up` (nicht beauftragt; ihre Aufrufer sind über die Aufrufer-Prüfung
und die Handlesung in §3 gebunden, nicht real gefahren).

## 2. DoD — Verdikt je Zeile (§2 des Plans)

`[x]` sind sieben Zeilen (Nr. 1–7), `[ ]` fünf (Nr. 8–12), zusammen zwölf — gezählt am Plan
(`grep -c '^- \[x\]'` 7, `grep -c '^- \[ \]'` 5). Der Doku-Update-Punkt (Nr. 7) ist `[x]` mit
„entfällt“ und wird als siebte `[x]`-Zeile geführt.

| # | DoD-Zeile | Verdikt | Realer Beleg |
|---|---|---|---|
| 1 | Werkzeug `suchlauf-nachmessen.sh` + `make suchlauf-nachmessen PLAN=`; Tabellentest mit elf Fällen, je an die Eingabe gebunden; drei benannte Mutationen rot | **erfüllt** | Skript gelesen: Blöcke mit Etikett `suchlauf`, `git grep -n <Optionen> [<Stand>] -- <Pathspec> <Ausschluss>`, Exit 0/1/2, `HEAD` und Namen abgelehnt (Hexform 7–40 oder `diff`), leerer Plan Exit 2; Makefile `$(if $(PLAN),,$(error …))` vor jedem Befehl. Test: elf `Fall`-Kommentare (1–11) gezählt, `make test-suchlauf-nachmessen` Exit 0 (§1). Mutationen aus dem DoD-Wortlaut selbst gefahren (§4): „Ausschluss der Plan-Datei entfernt“ **rot** (5 `FEHLER`), „Zahlenvergleich umgekehrt“ **rot** (6 `FEHLER`), „Optionsprüfung entfernt“ **rot** (Fall `-O` — „die Plan-Zeile hat ein Kommando ausgeführt“, Marker-Datei entstanden) |
| 2 | Träger nennen Werkzeug und Grenze: Sensor-Doku, `harness/README.md`-Zeile, `AGENTS.md` §3.13, `implement-slice` Schritt 18, Reviewer-Skill; Beleg Lesen + `make docs-check` | **erfüllt** | fünf Stellen im Diff gelesen (§3): jede nennt Aufruf und die Grenze „Zahlen und Stände, nicht Vollständigkeit von Suchraum und Muster“; keine behauptet ein Gate (Sensor-Doku „Kein Gate“, README-Zeile „kein Gate“, `AGENTS.md`: „ist kein Gate“); `make docs-check` innerhalb `make gates` Exit 0 (`d-check: 1199 Datei(en) geprüft, 0 Befund(e)`) |
| 3 | Rollout-Artefakte bleiben unverändert: Wrapper `rollout-restore.sh`, sieben Aufrufer, lokale Änderung bleibt; realer `make test-store` und `make test-replication`, danach `git status --short tools/schema` leer; Mutation „Wiederherstellung entfernt“ rot | **erfüllt, mit V-1 (INFO)** | Wrapper gelesen (`mktemp`-Sicherung, `trap … EXIT`, Exit-Code des Kommandos, lokale Änderung bleibt, fehlende Datei bleibt fehlend); die sieben Aufrufer von Hand über `git grep` gefunden, alle gehen durch `bash tools/schema/rollout-restore.sh make …` (`apply-rollout.sh`, `bench-lib.sh`, `run-integration-tests.sh`, drei `run-sdk-*-integration-tests.sh`, `examples/bootstrap.sh`); der Guard-Test trägt seine eigene Sicherung. **Realer Lauf beider Ziele** Exit 0 mit leerem `git status --short tools/schema` danach (§1). Mutation „Wiederherstellung entfernt“ (`trap restore EXIT` durch `trap : EXIT`) färbt `make test-rollout-restore` rot (6 `FEHLER`, §4). Das Gegenstück „der Parent-Stand lässt beide Dateien geändert“ ist **übernommen** aus dem Review-Report (dort: Wrapper aus `apply-rollout.sh` entfernt, `make test-store` Exit 0, danach `M tools/schema/plan.yaml`) und aus `BEO-PGC/test-schreibt-in-committete-datei` (4 Belege); von mir nicht gefahren |
| 4 | `make gates` grün, Exit gesondert | **erfüllt** | eigener Lauf EXIT=0 (§1) |
| 5 | Review durchgeführt, Report liegt vor | **erfüllt, mit V-2 (INFO)** | Report `review-slice-harness-suchlauf-nachmessen.md` liegt vor (1 HIGH, 2 MEDIUM, 1 LOW, 4 INFO; Verdikt merge-blockierend), Fixrunde `a124da45`; Finding für Finding in §5 nachgemessen — kein offenes HIGH/MEDIUM |
| 6 | §3.13-Suchlauf: Feld in §3, Gefundenes und Nichtgefundenes je Träger, beide Stände | **erfüllt** | acht Zeilen mit dem Werkzeug Exit 0 und alle acht Zeilen zusätzlich von Hand mit `git grep` nachgefahren (§6); die Befund-Zellen (Datei-Aufteilungen, „Nicht gefunden“) stimmen |
| 7 | Doku-Update: entfällt für das Benutzerhandbuch | **erfüllt (entfällt)** | keine Betreiber-Oberfläche berührt; `git diff --name-only 71ff3e44..HEAD -- docs/user` leer (das Handbuch trägt drei bestehende „nachmess“-Zeilen, im Diff nicht geändert) |
| 8 | Closure-Notiz mit Lerneintrag | **korrekt offen** | Plan §7 trägt „*(zu tragen bei Closure)*“ (Planner) |
| 9 | Reconciliation-Register — entfällt | **korrekt offen / entfällt** | Greenfield |
| 10 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | Closure-Pflicht (Planner) |
| 11 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | Plan §6: fünf Zeilen „Ausgang: *(bei Closure)*“ (Planner) |
| 12 | Die drei Paarungen | **korrekt offen** | hängen an der Closure der nächsten Welle |

Kein `[x]` ohne Beleg; kein `[ ]`, das über die Rollen-Sequenz hinaus belegt wäre.
Register, Risiko-Ausgänge, Paarungen und Closure-Notiz sind Planner-Arbeit und nicht Teil dieser Prüfung.

## 3. Plan-vs-Code-Diff

**§1 Ziel:** Werkzeug + Wrapper stehen; die Form der Zeile (`<Stand> <Soll> <Argumente>`, Etikett
`suchlauf`) ist im Skript, in der Sensor-Doku und in `AGENTS.md` §3.13 gleich beschrieben. Die vier
„Ausdrücklich NICHT“-Punkte sind eingehalten: kein Gate (`make gates` unverändert, `git diff
71ff3e44..HEAD -- Makefile` fügt nur drei `.PHONY`-Ziele hinzu), kein Parser über Tabellenzellen (das
Werkzeug meldet einen Plan ohne Block mit Exit 2), keine Prüfung der Vollständigkeit (Grenze an fünf
Stellen genannt), kein Mutations-Harness, kein weiterer Schreiber in committete Dateien angefasst.

**§3-Tabelle, Zeile für Zeile gegen den Diff (`git diff --stat 71ff3e44..HEAD`):**

| Plan-Zeile | Ist im Diff |
|---|---|
| `suchlauf-nachmessen.sh` neu | vorhanden, 245 Zeilen |
| `run-suchlauf-nachmessen-tests.sh` neu (elf Fälle nach Fixrunde) | vorhanden, 203 Zeilen, Fall 1–11 |
| `Makefile` (Ziele, `make help`) | drei Ziele, Hilfetext je Ziel (`make help` druckt alle drei, „elf Fälle“) |
| `harness/sensors/suchlauf-nachmessen.md` neu | vorhanden, Vertrag/Form/Ausgang/Grenze/Test |
| `harness/README.md` Werkzeug-Zeilen | drei Zeilen (`suchlauf-nachmessen`, `test-suchlauf-nachmessen`, `test-rollout-restore`); die Plan-Zeile nennt zwei, die dritte gehört zum Nachzug (Wrapper-Test) — im Rahmen |
| `AGENTS.md` §3.13, `implement-slice` Schritt 18, Reviewer-Skill | je Absatz im Diff |
| `apply-rollout.sh` + `rollout-restore.sh` (Wrapper) | Wrapper neu (51 Zeilen); `apply-rollout.sh` ruft ihn |
| sechs weitere Aufrufer | alle sechs im Diff (`run-integration-tests.sh`, drei SDK-Runner, `bench-lib.sh`, `examples/bootstrap.sh`), je eine Zeile |
| `run-rollout-restore-tests.sh`, `test-rollout-restore` | vorhanden, 86 Zeilen; Aufrufer-Prüfung liest `make … schema-rollout` unter `tools/*.sh`, `examples/*.sh` |
| `harness/targets/schema-rollout.md` (Nachzug + Fixrunde-Grenzen) | Abschnitt „Erzeugnisse in Test-, Bench- und Beispiel-Läufen“ samt „Grenzen der Rücknahme“ |
| Abweichung Ausschluss der Plan-Datei | im Skript: `:(exclude,glob)**/<Dateiname>`; im Test an Parent mit Plan unter `next/` gebunden; Mutation ohne den Ausschluss rot (§4) |
| Festlegung Form der Zeile (Allow-List) | im Skript `check_opts`/`check_paths`; Vertrag nennt die Liste; Test Fall 8 mit Marker-Datei |
| Fixrunde-Zeilen (Wrapper-Grenzen, Sensor-Doku, README, `make help`) | Skriptkopf des Wrappers nennt Nebenläufigkeit und `SIGKILL`; Sensor-Doku nennt Optionen, Magic, leeres Argument, Standardfehlerausgabe, `realpath`; „elf“ an vier Stellen gleich (Sensor-Doku, README, Makefile, Testkopf; `git grep` „fünf Fälle“ in `harness/README.md`: 0 Treffer) |

Nicht im Plan, im Diff: die vier Register-`state.md`-Zeilen (Link auf `open/…` durch Kennung ersetzt,
weil der Lifecycle-Move den Link bricht; im Plan §3 Befund-Zelle als Meldung geführt) und der
Review-Report. Im Plan, nicht im Diff: nichts. `git diff --name-only 71ff3e44..HEAD -- docs/plan/adr
spec .d-check.yml .a-check.yml Dockerfile harness/mk` ist leer: keine ADR/Spec/Schwelle/Gate-Konfiguration
verändert ([`AGENTS.md`](../../AGENTS.md) §3.5, §3.6).

**Abweichungen Plan ↔ Code:** keine unbenannte. Die zwei Abweichungen vom Plan-Wortlaut (Ausschluss
über Dateinamen-Glob statt Pfad; Rücknahme im Wrapper statt in `apply-rollout.sh`) stehen als Zeilen in
§3 und im dritten DoD-Punkt/§1 gleichlautend (Review F-3 behoben, §5).

**§3.13-Wortlaut und Vertrag gegen den realen Makefile-Stand (`make help`):** `make help` druckt genau die
Zeilen `suchlauf-nachmessen`, `test-rollout-restore`, `test-suchlauf-nachmessen`; `AGENTS.md` §3.13
(„endet bei jeder Abweichung, bei `HEAD` als Stand und bei einem Plan ohne Block mit Exit ≠ 0“) deckt sich
mit dem Skript (Exit 1 Abweichung; Exit 2 `HEAD`/leer). Die Sensor-Doku sagt „Über `make` kommt jeder
Exit ≠ 0 als der Make-eigene Exit `2` an“ — im eigenen Lauf mit einem Wegwerf-Plan (`diff 1 -e … -- AGENTS.md`,
Ist 0) gemessen: gedruckt `ABWEICHUNG  soll=1 ist=0`, `make: *** … Fehler 1`, Make-Exit `2`; ohne `PLAN`
gedruckt `PLAN fehlt, z.B. make suchlauf-nachmessen PLAN=…`, Make-Exit `2`, bevor ein Befehl läuft. Die
`harness/README.md`-Zeilen benennen `bash` + `git`; die Sensor-Doku nennt zusätzlich `realpath`
(Review F-6, dort als Vertragsgenauigkeit gefordert): gedeckt.

## 4. Mutationen der Eingabeseite (dieser Lauf)

Je Mutation eine Kopie des Prüflings, Testlauf mit `TOOL=<Kopie>`, Ausgang und `FEHLER`-Meldungen
gelesen; die Originale unberührt. Die ersten drei stammen aus dem DoD-Wortlaut, die übrigen sind
Eingabeseiten der Fixrunde-Zusagen.

| # | Mutation | Lauf | Ergebnis (gedruckt) |
|---|---|---|---|
| MA | Optionsprüfung entfernt: `check_opts` liefert immer 0 | `run-suchlauf-nachmessen-tests.sh`, **EXIT=1** | `FEHLER: -O abgelehnt — Exit 0, erwartet 2`; `… — es lief trotzdem eine Messung`; `… — die Plan-Zeile hat ein Kommando ausgeführt`; ebenso `--open-files-in-pager`, `--no-index`, `-f`, `-e ohne Muster` — rot |
| MB | Stand `HEAD` als gültige Kennung zugelassen (Hex-Regex um `\|HEAD` erweitert) | **EXIT=1** | `FEHLER: HEAD abgelehnt — Exit 0, erwartet 2`; `… — es lief trotzdem eine Messung` — rot |
| MC | Ausschluss der Plan-Datei entfernt (`excludes=()`) | **EXIT=1**, 5 `FEHLER` | Fälle „stimmt“ und „Selbstverweis“ zählen den Plan selbst — rot |
| MD | Zahlenvergleich umgekehrt (`-eq` → `-ne`) | **EXIT=1**, 6 `FEHLER` | rot |
| ME | `git grep`-Exit über 1 ignoriert (Schwelle `-gt 1` → `-gt 200`) | **EXIT=1** | `FEHLER: git grep-Fehler — Exit 0, erwartet 2` — rot (Review F-2, dort grün, jetzt gebunden) |
| MF | nicht geschlossener Block ignoriert (Prüfung durch `if false` ersetzt) | **EXIT=1**, 2 `FEHLER` | rot (Review F-2) |
| MG | Commit-Stand ohne Pathspec ausgeführt (`"${paths[@]}"` aus dem Aufruf am Commit-Stand entfernt) | **EXIT=1** | `FEHLER: Commit-Stand mit Pathspec — Exit 1, erwartet 0` — rot (Review F-2, dort grün) |
| MW | Wrapper: `trap restore EXIT` durch `trap : EXIT` ersetzt (Wiederherstellung entfernt) | `run-rollout-restore-tests.sh`, **EXIT=1**, 6 `FEHLER` | rot |

Zwei erste Versuche waren äquivalent bzw. zu schwach und zählen nicht: `-gt 99` statt `-gt 1` (Exit 128
löst weiter aus) und die Entfernung nur des `--` am Commit-Stand-Aufruf (`git grep` liest den Pathspec auch
ohne `--`) blieben grün; ich habe sie durch ME und MG ersetzt, die die Zusage tatsächlich aufheben.
Nicht gebunden bleibt, wie im Vertrag benannt, die Trennung der Standardfehlerausgabe von der
Trefferzählung (der Test erzeugt keinen `git grep`-Lauf mit Hinweis und Exit 0 oder 1).

## 5. Findings des Reviews nachgemessen (nicht dem Fixrunden-Bericht geglaubt)

| Finding | Was ich gemessen habe | Verdikt |
|---|---|---|
| F-1 (HIGH) Argument-Injektion `-O`/`--open-files-in-pager` | `check_opts` gelesen (Allow-List; jede andere Option Exit 2, gebündelte Kurzoptionen nur aus `[iwEFGPvIahHocLln]`, höchstens ein Suchwort, keines neben `-e`); `check_paths` (Pathspec-Magic auf `:!`, `:^`, `:(exclude\|glob\|literal\|icase\|top)`); Test Fall 8 mit Marker-Datei für `-O`, `--open-files-in-pager`, `-O` ohne Kommando, gebündelt `-inO`, `--no-index`, `-f`; MA rot mit „die Plan-Zeile hat ein Kommando ausgeführt“ | **behoben**, mutiert bestätigt |
| F-2 (MEDIUM) Vertrags-Exit-Zweige ohne Testbindung | Fälle 6, 9, 10, 11 neu (Commit-Stand mit Pathspec, `git grep`-Fehler, nicht geschlossener Block, Zeilenform mit Soll/Stand/Muster/leerem Argument); ME, MF, MG rot; Soll ohne Zahl, unbekannte Kennung und Zeile ohne Muster über Fall 11 gebunden (Ausgaben `keine Zahl`, `kein Commit`, `ohne Suchmuster`) | **behoben**, mutiert bestätigt |
| F-3 (MEDIUM) Plan §1/DoD nennen `apply-rollout.sh` als Ort der Rücknahme, §3 den Wrapper | Plan §1 sagt „der Wrapper `tools/schema/rollout-restore.sh` stellt … wieder her“, DoD-Punkt 3 nennt Wrapper und sieben Aufrufer, §3-Zeile `apply-rollout.sh` sagt „nicht im Skript selbst, sondern über `rollout-restore.sh`“; die Begründung „der Plan lasse … offen“ ist entfernt | **behoben** |
| F-4 (LOW) leeres Argument verschwindet; `2>&1` zählt Stderr | Skript: Wörter mit `-z` geprüft, Meldung „leeres Argument“ Exit 2 (Test Fall 11 `leeres Argument`); Standardfehlerausgabe in `$errf` getrennt, nur als „Hinweis von git grep (nicht gezählt)“ auf stderr; die Trennung ist nicht test-gebunden und im Vertrag so genannt | **behoben**; Trennung als Grenze benannt |
| F-5 (INFO) Wrapper nicht nebenläufig, SIGKILL | Skriptkopf „Grenzen: zwei gleichzeitig laufende Aufrufe …; ein SIGKILL lässt die Erzeugnisse verändert“, `schema-rollout.md` „Grenzen der Rücknahme“ | **benannt** |
| F-6 (INFO) Host-Werkzeuge (`realpath`), §3.1 | Sensor-Doku nennt `bash`, `git`, `realpath` und die Klasse „wie `git rev-parse` in `apply-rollout.sh`, es wird nichts installiert“; die Vereinbarkeit mit [`AGENTS.md`](../../AGENTS.md) §3.1 ist eine Architect-/Planner-Aussage (Plan §6, Risiko 1, Ausgang offen) | **Vertrag genau**; Ausgang beim Planner |
| F-7 (INFO) Aufrufer-Prüfung nur Einzelzeilen | `schema-rollout.md` „Grenzen der Rücknahme“ nennt Fortsetzungszeile, `$MAKE`, Makefile-Rezept, Compose-Kommando; von Hand: `git grep -n -E 'schema-rollout' -- 'tools/*.sh' 'examples/*.sh'` liefert sieben Aufruf-Zeilen (alle durch den Wrapper) und sonst Kommentare/Echo/Fehlertexte | **benannt** |
| F-8 (INFO) vier Register-`state.md` führen den Slice als „geplant“ | Planner-Closure-Arbeit, im Plan §3 als Meldung mit Frist „Closure dieses Slice“ geführt; die vier `state.md` sind im Diff (Link → Kennung), der Zustandstext unverändert | **bewusst offen** (Meldung, Planner) |

Kein offenes HIGH, kein offenes MEDIUM. Ein offener Punkt über INFO hinaus: keiner.

## 6. Suchlauf-Feld (§3 des Plans), an beiden Ständen nachgefahren

Parent `8717c4fb`. Von Hand mit `git grep -n … | wc -l` und dem Ausschluss
`':(exclude,glob)**/slice-harness-suchlauf-nachmessen.md'` (ohne das Werkzeug), alle acht Zeilen:

| Zeile | Stand | Soll | Von Hand |
|---|---|---|---|
| 1 `-i -E 'suchlauf\|nachmess\|Suchform'` über die Träger-Pfade | `8717c4fb` | 14 | **14** (`AGENTS.md` 9, Handbuch 3, Reviewer-Skill 2) |
| 2 dasselbe | `diff` | 80 | **80** (je Datei 5/4/13/7/3/3/10/12/23; Werkzeug-Skripte 12 + 23 = 35 wie in der Befund-Zelle) |
| 3 Formmuster, ganzer Baum ohne Berichte/`done`/Baseline | `8717c4fb` | 9 | **9** |
| 4 dasselbe | `diff` | 12 | **12** |
| 5 `plan\.yaml\|down\.sql`, ganzer Baum ohne Berichte/`done`/Baseline | `8717c4fb` | 59 | **59** |
| 6 dasselbe | `diff` | 74 | **74** |
| 7 `make .*schema-rollout\|apply-rollout` über die Aufrufer-Bäume | `8717c4fb` | 35 | **35** (`git grep -c`: 15 Dateien) |
| 8 dasselbe | `diff` | 47 | **47** (`git grep -c`: 18 Dateien) |

Die Datei-Aufteilungen der Befund-Zelle (`AGENTS.md` 9 → 13, Reviewer-Skill 2 → 4, Handbuch 3 unverändert,
„15 → 18 Dateien“) stimmen. „Nicht gefunden: eine weitere Beschreibung der Suchform in
`.claude/commands/*` am Parent“: am Parent trägt nur `AGENTS.md`, Reviewer-Skill und Handbuch das Wort
(`git grep -c` am Parent, drei Dateien) — bestätigt. Der Stand `diff` ist der Arbeitsbaum dieses Laufs
(`HEAD` = `b1a6d76d`, sauber); die Zahlen streuen mit jedem Commit, der die Muster berührt.

## 7. Entscheidungs-Konformität

- **[`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md) §Entscheidung 4 („es gibt keinen
  Sensor“ auf Prosa):** konform. Das Werkzeug prüft nicht, ob eine Zahl ihren Ursprung trägt; es wiederholt
  eine **deklarierte** Messung (Befehl, Stand, Zahl) und sagt das in der Grenze der Sensor-Doku, in
  `AGENTS.md` §3.13, in der README-Zeile, im Schritt 18 und im Reviewer-Skill. Es ist kein Gate (nicht in
  `make gates`, Begründung „`diff` bewegt sich“), und die Vollständigkeit von Suchraum und Muster bleibt
  Lese-Handlung des Reviewers. Der Reviewer-Skill-Absatz ändert die Rolle des Reviewers nicht.
- **[`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md) Entscheidung 3** (Pflicht-Report und
  Rollback-Artefakt „werden je Rollout aufbewahrt“): konform mit Anmerkung V-1. Das Target
  `schema-rollout` selbst ist unverändert; im Betrieb bleiben `plan.yaml` und `down.sql` sein Erzeugnis.
  Der Wrapper wirkt nur in Läufen, die den Rollout gegen Wegwerf-Datenbanken als Vorbedingung fahren.
- **[`AGENTS.md`](../../AGENTS.md):** §3.3 (die Lifecycle-Moves `42e5ebb8`, `8717c4fb` sind reine Renames,
  der Inhalt steht in eigenen Commits), §3.5/§3.6 (kein ADR/Spec/Schwelle im Diff), §3.7 (Kommentare der
  neuen Skripte tragen Zusage/Kopplung/Grenze im Indikativ; der Wrapper-Kopf nennt Grenzen, keine
  verworfene Alternative), §3.9 (Wrapper endet mit `"$@"` und trägt den Exit-Code des Kommandos; meine
  Läufe ungepiped), §3.11 (kein host-lokaler Pfad — `docs-check` Exit 0), §3.13 (Träger-Meldung fremder
  Dateien statt stiller Mitänderung: die vier `state.md` nur in der Verweis-Form angefasst, Zustandstext
  als Meldung geführt). §3.1: Host-Werkzeuge in V-3 benannt.
- **Commits:** `make doc-commits RANGE=71ff3e44..HEAD` Exit 0, `make commit-traceability` (Standing-Gate)
  Exit 0; jede der neun Betreff-Zeilen nennt `ADR-0083` (Bau- und Fixrunden-Commit zusätzlich
  `ADR-0043`), keine trägt `SPEC-`/`ARC-`.

## 8. Findings dieser Verifikation

| # | Kategorie | Befund | Quelle | Verifizierbar |
|---|---|---|---|---|
| V-1 | INFO | [`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md) Entscheidung 3 sagt „je Rollout aufbewahrt“ ohne Unterscheidung von Betrieb und Test-Vorbedingung; der Slice verwirft die Erzeugnisse der Test-, Bench- und Beispiel-Läufe (Wegwerf-Datenbanken) und begründet das im Vertrag `schema-rollout.md` (im Betrieb bleibt der Beleg). Die Lesart ist plausibel, aber eine Interpretation der ADR-Zeile, keine ihrer Aussagen. | `harness/targets/schema-rollout.md` §Erzeugnisse in Test-, Bench- und Beispiel-Läufen; [`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md) Entscheidung 3 | ja — Lesen der beiden Stellen |
| V-2 | INFO | Der Review-Report liest die Range bis `fc0f5629`; die Fixrunde (`a124da45`, `b1a6d76d`: Allow-List, neue Testfälle, Wrapper-Grenzen, Plan-Nachzug) hat keinen zweiten Reviewer-Lauf. Ihre Wirkung ist in §4 und §5 mutiert bestätigt; ein Reviewer-Blick auf `check_opts`/`check_paths` (neue Sicherheitsfläche des Werkzeugs) bleibt eine Rollen-Entscheidung des Planners. | `tools/harness/suchlauf-nachmessen.sh` `check_opts`, `check_paths` | ja — Lesen; §4 MA |
| V-3 | INFO | Das §6-Risiko „Docker-only für `git`“ ist mit der Fixrunde nicht ausgesprochen: das Werkzeug ruft `git`, `bash` und `realpath` auf dem Host (Präzedenz `apply-rollout.sh`, `run-release-tag-info-tests.sh`); der Vertrag nennt die Klasse, der Ausgang steht als *(bei Closure)* im Plan. Ausgang und Architect-Aussage sind Planner-Arbeit. | Plan §6 Risiko 1; Sensor-Doku §Vertrag | nein |
| V-4 | INFO | Das Gegenstück zur Rücknahme („der Parent-Stand lässt `plan.yaml` geändert“) habe ich nicht real gefahren; es ist aus dem Review-Report und dem Beobachtungs-Eintrag übernommen. Die Wirkung der Rücknahme selbst ist real gemessen (beide DB-Tier-Läufe, leerer Status, und die Wrapper-Mutation rot). | §2 Nr. 3 | ja — `make test-store` am Parent mit entferntem Wrapper |

Kein HIGH, kein MEDIUM, kein LOW. Keine DoD-Verletzung.

## 9. Verdikt

**DoD bestätigt:** ja — jede der sieben `[x]`-Zeilen (Nr. 1–7) ist am Ist-Zustand belegt; die fünf
`[ ]`-Zeilen (Nr. 8–12) sind korrekt offen (Planner-Closure). **Plan-vs-Code:** keine unbenannte
Abweichung. **Entscheidungs-Konformität:** [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md) und [`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md) konform, V-1 als Auslegungs-Anmerkung.
**Gates:** `make gates` Exit 0 im eigenen Lauf; `make test-store` und `make test-replication` real Exit 0
mit unverändertem `tools/schema`.

**Übergabe:** an den Planner — Closure-Notiz mit Lerneintrag, Beobachtungs-Register (darunter F-8: die
vier `state.md` führen den Slice weiter als „geplant“), Ausgänge der fünf §6-Risiken (V-3 für Risiko 1),
Paarungen bei der Closure der nächsten Welle; V-1/V-2 als Anmerkungen zur Entscheidung des Planners.
Dieser Report ist ein **Lauf-Beleg** (dieser Stand, dieser Lauf) und ersetzt weder Review noch Closure.
