# Verifikations-Report: slice-code-kommentare-kennungen — 2026-09-26

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich + Entscheidungs-Konformität +
Plan-vs-Code-Diff + Gates. Review-Artefakt des Reviewers:
[`review-slice-code-kommentare-kennungen.md`](review-slice-code-kommentare-kennungen.md);
Formvorbild dieses Reports:
[`verifikation-slice-transformationen-backfill-pfad.md`](verifikation-slice-transformationen-backfill-pfad.md).
Die vendored Baseline trägt kein eigenes Verifikations-Template
(`v6.9.0` · `regelwerk/`-Bundle: das Verzeichnis `templates/docs/reviews/` enthält nur
`review-report.template.md`, Gerüst per `cp` übernommen und durch die Form des Formvorbilds ersetzt).

**Gegenstand:** Slice-Plan `slice-code-kommentare-kennungen` (wellenlos, Harness-Querschnitt), `HEAD` = `f049bf20`,
Diff-Range `1021f6fe..HEAD`, 9 Commits, 13 Dateien (+1574/−39 einschließlich zweier Lifecycle-Moves und des
Review-Reports). Slice-Inhalt: Lifecycle (`ec6cb07e`, `0d333120`), Plan-Feld (`0484a748`), Implementer-Lauf
(`c3a172c6`, `e52b6780`, `e5efeee3`, `71cf9537`), Review-Report (`c75abb87`), Fixrunde (`f049bf20`). Der Stand ist nicht
gepusht (`origin/main` = `1021f6fe`). Dieser Lauf ändert weder Code noch Plan noch Doku; er schreibt nur diesen Report.
Die Mutationen liefen an Kopien im Scratchpad der Sitzung (Ersetzung mit Bash-Parameterexpansion auf der Kopie, kein
Textwerkzeug am Repo); `git status --short` danach leer.

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — [`AGENTS.md`](../../AGENTS.md) §3.9)

Jeder Lauf schrieb in eine Log-Datei; der Exit-Code wurde im selben Aufruf gesondert gesichert und danach gelesen. Stand
aller Läufe ohne Mutation: `HEAD` = `f049bf20`, Arbeitsbaum sauber. Je ein schwerer Docker-Lauf zugleich.

| Sensor | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `make test-kommentar-kennungen` | **Exit 0** | `ok …/tools/harness/kommentar-kennungen 0.006s`, `run-kommentar-kennungen-tests: alle Fälle bestanden` |
| `make test` (Race-Detector) | **Exit 0** | `ok …/tools/harness/kommentar-kennungen 1.029s`, keine Zeile `FAIL` (`grep -c FAIL` = 0) |
| `make a-check` | **Exit 0** | `gesamt: 0 Befund(e)` |
| `make coverage-gate` | **Exit 0** | `coverage-gate: OK — Coverage 85.00% erfüllt Schwelle 80%` |
| `make gates` (einmal) | **Exit 0** | `baseline-verify: v6.9.0 OK — 54 Dateien` · `d-check: 1256 Datei(en) geprüft, 0 Befund(e)` · `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD"` · `generated-sync: OK` · `gesamt: 0 Befund(e)` |
| `make suchlauf-nachmessen PLAN=<Plan des Slice>` | **Exit 0** | `suchlauf-nachmessen: 16 Zeilen stimmen` |
| `make commit-traceability RANGE=1021f6fe..HEAD` | **Exit 0** | `OK — 9 Commit(s) in "1021f6fe..HEAD", Betreffs ohne Struktur-ID` |
| `make doc-commits RANGE=1021f6fe..HEAD` | **Exit 0** | `d-check: 1256 Datei(en) geprüft, 0 Befund(e)` |
| `make doc-immutable RANGE=1021f6fe..HEAD` | **Exit 0** | `d-check: 1256 Datei(en) geprüft, 0 Befund(e)` |
| gofmt (Docker, gepinntes Toolchain-Image, `gofmt -l` über die drei `*.go`-Dateien des Diffs, Befehl aus Schritt 18 des Implementer-Ablaufs) | **Exit 0** | keine Ausgabe |
| Mutationen (§4) | 7 Programm-Mutationen (M1–M7), 4 Aufrufer-Mutationen (S1–S4), 1 Nenner-Messung | alle rot, siehe §4 |

**Coverage-Zahl als Lauf-Beleg ([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz A).** 85.00 % ist der Ausdruck meines Laufs
(einmal allein, einmal im Schritt `coverage-gate` von `make gates`, beide 85.00 %); `tools/` gehört nicht zur gemessenen
Fläche, der Slice bewegt die Zahl nicht (Plan §3 „Coverage-Gate: geprüft, unverändert“, gegen `harness/mk/coverage.mk`
und den Diff gelesen: kein Bestandteil des Diffs).

**Hygiene:** dangling Volumes (`docker volume ls -qf dangling=true | wc -l`) vor dem ersten Lauf **34**, nach allen
Läufen **34**; kein `prune`, kein `system prune`. Nicht gefahren: `make test-integration`, `make test-store`,
`make test-replication`, `make bench` (der Diff berührt weder `internal/` außer Kommentarzeilen einer Datei noch
Schema, Compose noch Runner).

## 2. DoD — Verdikt je Zeile (§2 des Plans)

Gezählt am Plan: sieben `[x]`-Zeilen, fünf `[ ]`-Zeilen, zusammen zwölf (`grep -c '^- \[x\]'` 7,
`grep -c '^- \[ \]'` 5).

| # | DoD-Zeile | Verdikt | Realer Beleg |
|---|---|---|---|
| 1 | Liefer-Punkt 1 — Regel und Träger (`AGENTS.md` §3.7, Schritt 20, Reviewer-Unterpunkt) | **erfüllt** | die drei Stellen im Diff gelesen (§5); `make docs-check` im `make gates`-Lauf Exit 0; `make suchlauf-nachmessen` 16 Zeilen stimmen; Einstufung MEDIUM für die Form, HIGH-Eskalation für eine falsche Wiedergabe steht im Unterpunkt mit Begründung |
| 2 | Liefer-Punkt 2 — das Werkzeug (Ziel, Programm, Aufrufer, Vertrag, README-Zeile, `a-check` grün) | **erfüllt** | `make kommentar-kennungen` Exit 0/1/2 gefahren (Kandidaten: Exit 1 des Skripts, über `make` Exit 2 mit `Fehler 1`); netzlos (`--network none`, Toolchain-Image gepinnt); Tabellentest grün; Mutationen der Eingabeseite in §4 rot; Selbstanwendung 0 Kandidaten (§3); Basis-Messung nachgemessen (§3); `make a-check` Exit 0; Zeilen in `harness/README.md` gegen `make help` (§6) |
| 3 | Liefer-Punkt 3 — Erstbeleg `changestream.go` | **erfüllt** | 0 Kandidaten (Ausgabe leer, Exit 0, `COUNT=1` druckt 0); der Diff der Datei ohne Nicht-`//`-Zeile: **0**; `Publish` trägt alle Zusagen des Plans und genau einen Anker; gofmt leer; `make test` Exit 0 (§6) |
| 4 | `make gates` grün, Exit gesondert | **erfüllt** | eigener Lauf Exit 0 (§1) |
| 5 | Review durchgeführt, Report liegt vor | **erfüllt** | Report liegt vor (0 HIGH, 1 MEDIUM, 4 LOW, 4 INFO); Finding für Finding in §7 gegen den Ist-Zustand nachgemessen — kein offenes HIGH oder MEDIUM |
| 6 | §3.13-Suchlauf: Feld in §3, Gefundenes und Nichtgefundenes je Träger, beide Stände | **erfüllt** | 16 Zeilen mit dem Werkzeug Exit 0 (§1); acht Zeilen von Hand mit `git grep` nachgefahren (§3); die Träger-Tabelle nennt Gefundenes je Träger und das Nichtgefundene (kein Träger in `README.md`/`docs/user/`, Handbuch unberührt) |
| 7 | Doku-Update: `AGENTS.md` §3.7, `harness/README.md` §Sensors, Handbuch unberührt | **erfüllt** | beide Träger im Diff; der Diff berührt `docs/user/` nicht (`git diff --stat 1021f6fe..HEAD -- docs/user` leer) |
| 8 | Closure-Notiz mit Steering-Loop-Lerneintrag | **korrekt offen** | Plan §7 trägt „*(zu tragen bei Closure)*“ (Planner) |
| 9 | Reconciliation-Register — entfällt (Greenfield) | **korrekt offen (entfällt)** | keine Datei, Zeile steht als „entfällt“ im Plan (Häkchen-Form: der Planner setzt es bei der Closure) |
| 10 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | Closure-Pflicht (Planner); der Diff berührt keine Register-Datei |
| 11 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | Plan §6: alle sieben Zeilen tragen „Ausgang: *(bei Closure)*“ (Planner) |
| 12 | Die drei Paarungen | **korrekt offen** | hängen an der Closure des Slice (Planner) |

Kein `[x]` ohne Beleg; kein `[ ]`, das über die Rollen-Sequenz hinaus belegt wäre. Register, Risiko-Ausgänge, Paarungen
und Closure-Notiz sind Planner-Arbeit und nicht Teil dieser Prüfung.

## 3. Plan-vs-Code-Diff und Zahlen mit Ursprung ([`AGENTS.md`](../../AGENTS.md) §3.12)

**§1 Ziel und Zielform.** Regel (§3.7), Aufruf in Schritt 20, Reviewer-Probe und Erstbeleg stehen. Die
„Ausdrücklich NICHT“-Punkte sind eingehalten: kein Eintrag in einem Gate-Bündel (`make gates` läuft ohne das Werkzeug;
`Makefile`-Diff nennt nur `.PHONY`-Ziele außerhalb der Gate-Liste); kein Satz-Subjekt und keine Slice-/Wellen-Nummer im
Programm (`main.go` gelesen: Zählregel über Kennungen und „ff.“, Blockgrenzen aus `go/parser`); Bestand außer der
Datei des Erstbelegs unberührt (`git diff --stat` nennt nur `changestream.go` unter `internal/`); `sdks/` und `gen/`
ausgenommen (`excludedRoots`); keine Ausnahme-Marker im Quelltext.

**§3-Tabelle, Zeile für Zeile gegen `git diff --stat`:**

| Plan-Zeile | Ist im Diff |
|---|---|
| `AGENTS.md` §3.7 update | 31 Zeilen: Regelabsatz, Falsch/Richtig-Paar, Werkzeug-Verweis samt Grenze |
| `implement-slice.md` Schritt 20 update | 16 Zeilen; Nummerierung 20/21 unverändert (`grep -nE '^(19\|20\|21)\. '`: Zeilen 192, 218, 273) |
| `reviewer.md` update | 20 Zeilen, MEDIUM-Abschnitt |
| `tools/harness/kommentar-kennungen/` (`main.go`, `main_test.go`) neu | 351 bzw. 366 Zeilen |
| `tools/harness/kommentar-kennungen.sh` neu | 75 Zeilen |
| `Makefile` update | Ziel `kommentar-kennungen` mit `.PHONY` und Hilfe-Zeile; kein Eintrag in Gate-Zielen |
| `harness/sensors/kommentar-kennungen.md` neu | 168 Zeilen |
| `harness/README.md` update | zwei Zeilen unter Werkzeuge |
| `changestream.go` update | drei Blöcke, nur Kommentarzeilen |
| `.a-check.yml` prüfen | unverändert (Diff nennt sie nicht), `make a-check` Exit 0 |

Über den Plan hinaus, als **Nachzug** im Plan benannt und begründet: `test-kommentar-kennungen` (Makefile,
`harness/README.md`), `run-kommentar-kennungen-tests.sh` (Fixrunde, F-2 des Reviews), Kompaktform mit `…`, Temp-Datei
statt Pipe, `go build` + `exec` statt `go run`. Ein Punkt der Nachzug-Tabelle weicht vom Wortlaut der DoD ab
(Liefer-Punkt 2: „Pipe unter `bash` mit `set -o pipefail`“ gegen Temp-Datei) — die Abweichung steht mit Begründung im Plan
und im Vertrag, die Begründung trägt ([`AGENTS.md`](../../AGENTS.md) §3.9: in der Pipe verschwände der Exit des linken
Glieds). Keine `Accepted` ADR wurde geändert (`make doc-immutable` Exit 0); keine Schwelle und keine Gate-Konfiguration
verändert ([`AGENTS.md`](../../AGENTS.md) §3.5, §3.6). Die beiden Lifecycle-Commits sind reine Renames (`git show -M
--stat`: je ein Rename, 0 Zeilen; [`AGENTS.md`](../../AGENTS.md) §3.3).

**Zahlen mit Ursprung.**

| Zahl | Stand und Ursprung im Plan | Mein Lauf |
|---|---|---|
| 600 / 403 / 197 (gesamt / Nicht-Test / Test) | Stand `0d333120`, Werkzeug dort nur im Arbeitsbaum — der Plan weist sie als **mit keinem Repo-Befehl wiederholbar** aus | nicht nachmessbar (bestätigt: Stand `0d333120` trägt das Werkzeug nicht im Commit); konsistent mit 597/400/197 nach dem Erstbeleg (600 − 3 Blöcke = 597, 403 − 3 = 400, Test unverändert 197) |
| 597 / 400 / 197 | Diff-Stand, der Plan nennt Befehl `make kommentar-kennungen COUNT=1` (gemessen) | `COUNT=1` → **597**, `TESTS=exclude` → **400**, `TESTS=only` → **197** (Exit 0, HEAD `f049bf20`) |
| 1474 (Nenner) | Plan: derselbe Zähler bei auf „mindestens eine Kennung“ gesenkter Schwelle, Änderung zurückgenommen | Kopie des Programms mit `>= 1`, gegen den Baum: **1474** (782 Nicht-Test, 692 Test) |
| 40,5 % | abgeleitet aus 597/1474 | 597 / 1474 = 0,4050 (**abgeleitet**), unter der Hälfte; Rückführung (a) tritt nicht ein |
| 396 / 192 (Prototyp) | Plan §1: „erwartet, nicht belegt“, Wegwerf-Prototyp | nicht nachmessbar (Prototyp nicht committet, Plan benennt das) |

Beide Stände sind im Plan ausgewiesen (Stand `0d333120` als nicht wiederholbar, Diff-Stand als wiederholbar).
Die Träger-Meldung zu den Prototyp-Zahlen (`slice-code-kommentare-bereinigung` §1) ist im Plan geführt.

**Selbstanwendung und Erstbeleg (gemessen):** `make kommentar-kennungen PATHS=tools/harness/kommentar-kennungen` Exit 0,
Ausgabe leer; `PATHS=internal/application/port/outbound/changestream.go` Exit 0, Ausgabe leer, mit `COUNT=1` **0**;
`DIFF=1021f6fe COUNT=1` (der ganze Slice-Diff) **0**. `DIFF=HEAD~40 COUNT=1` meldet **39** Blöcke, mit
`GIT_CONFIG_COUNT=1 GIT_CONFIG_KEY_0=diff.mnemonicPrefix GIT_CONFIG_VALUE_0=true` ebenfalls **39** (Befund F-1 des
Reviews behoben, §7).

**Suchlauf-Feld von Hand nachgefahren** (`git grep` ohne das Werkzeug, Ist gleich Soll), acht der 16 Zeilen: Parent `7b70b34a`
Kennungs-Zeilen **2233** (Zeile 1), `-l` **243** (2), `-l` ohne Testdateien **125** (3), „`ff.`“ **8** (5), „auflösbares
Feld“ **2** (6), „§3.7“ **6** (7); Stand `diff` Kennungs-Zeilen **2236**, „`ff.`“ **8**, „auflösbares Feld“ **4**, „§3.7“
**13**; `Zeile 266` am Stand `1021f6fe` **1** (`slice-harness-fmt-check`, Zeile 193). Zusammensetzung der Diff-Zahlen (F-6):
2236 = 2233 − 8 + 11 Fixture-Zeilen in `main_test.go` (nicht von Hand zerlegt, aber die Größenordnung deckt sich mit der
Diff-Zahl: `changestream.go` elf Zeilen mit Kennung vorher, drei nachher). Die Zeilenverschiebung der LOW-Zeile
`gofmt -l` im Reviewer-Skill ist real: **286** jetzt (`grep -n` gedruckt), Zeile 266 am Stand `1021f6fe` (Meldung des
Implementers an die fremde Plan-Datei trifft zu).

## 4. Mutationen der Eingabeseite (dieser Lauf)

Programm: je eine Kopie von `go.mod`, `main_test.go` und einer mutierten `main.go` im Scratchpad, Einzellauf `go test`
im gepinnten Toolchain-Image ohne Netz. Die Ersetzung prüft, dass sich der Text tatsächlich ändert; eine erste
Mutation der Kompaktform-Schleife (Schleife mit `if true { break }`) brach am Compiler („declared and not used“) und
zählt nicht — sie wurde durch eine kompilierende Mutation ersetzt.

| # | Zusage | Mutation | Ergebnis (gedruckt) |
|---|---|---|---|
| M1 | Kandidat ab **zwei** verschiedenen Kennungen | Schwelle `>= 2` auf `>= 3` | **rot**: `TestCandidate`, `TestRunModes`, `TestRunDiffMode` |
| M2 | Kompaktform zählt je Nummer | Fortsetzungsschleife nach der ersten Nummer abgebrochen | **rot**: `TestScanText` |
| M3 | Zieldatei-Zeile ohne Präfix `b/` endet mit Exit 2 | Präfixprüfung in `parseDiff` (`case strings.HasPrefix(path, "b/")`) zu `case true` | **rot**: `TestParseDiffPrefix`, `TestRunInputErrors` |
| M4 | Diff-Modus meldet nur Blöcke, die eine hinzugefügte Zeile überlappen | Überlappungsprüfung invertiert | **rot**: `TestOverlaps`, `TestRunDiffMode` |
| M5 | Kennung im Zeichenketten-Literal ist kein Kandidat | Literale zusätzlich gescannt (`ast.Inspect` über `BasicLit`) | **rot**: `TestBlocksOf` |
| M6 | Wurzel `.git` wird nicht gelesen | Eintrag aus `excludedRoots` entfernt | **rot**: `TestRunModes` (die Bindung aus F-8 des Reviews trägt) |
| M7 | Direktiven `//go:…` zählen nicht mit | Ausschluss entfernt | **rot**: `TestBlocksOf` |

Aufrufer (`run-kommentar-kennungen-tests.sh` gegen eine mutierte Kopie des Skripts über `TOOL`; Lauf der unmutierten
Fassung: Exit 0):

| # | Zusage | Mutation | Ergebnis |
|---|---|---|---|
| S1 | Eine Diff-Basis, die kein Commit ist, endet mit Exit 2 (keine Optionsinjektion) | Commit-Prüfung durch `if false` ersetzt | **rot**, Exit 1: „unbekannte Diff-Basis — Ausgabe trägt ‚kein Commit‘ nicht“, `fatal: bad revision` |
| S2 | Keine Temp-Datei bleibt zurück | `trap` entfernt | **rot**, Exit 1 |
| S3 | Die Form des Diff-Stroms ist gegen fremde Konfiguration gepinnt | `--src-prefix=a/ --dst-prefix=b/` entfernt | **rot**, Exit 1 |
| S4 | Der Exit des Programms wird weitergegeben | `exit 0` nach dem Docker-Aufruf | **rot**, Exit 1 |

Acht Mutationsklassen der Eingabeseite (Schwelle, Kompaktform, Präfixprüfung, Überlappung, Literal-Ausschluss,
`.git`-Ausschluss, Commit-Prüfung der Basis, `trap`) sind je einmal rot gesehen. `git status --short` danach leer.

## 5. Regel-Träger und Entscheidungs-Konformität

**`AGENTS.md` §3.7 (Wortlaut).** Konkret und prüfbar: höchstens eine Kennung je Kommentar, keine Kette, keine
Kompaktform (`…-003/005`, `…-002`…`004`), kein „ff.“, keine Spec-Wiedergabe in eigenen Worten, Kopplung nennt die Stelle;
das Werkzeug wird mit seiner Grenze genannt. Baseline (`v6.9.0` · `regelwerk/grundlagen-harness-dateien.md` §Was ein
Kommentar trägt, Zeile 107: „nennt sie als **ein** auflösbares Feld“) und §3.7 widersprechen sich nicht: „höchstens eine
Kennung“ konkretisiert „ein Feld“. Die Zählregel ist **strenger** als der Baseline-Wortlaut (auch zwei Anker mit je einer
eigenen Aussage werden gemeldet); das ist im Vertrag begründet (`harness/sensors/kommentar-kennungen.md` §Belege der
Definition „Strenger als der Baseline-Wortlaut“, Grenze 3) und im Plan §1 als bewusste Zielform gesetzt. Das
Falsch/Richtig-Paar ist selbst konform: das „Falsch“-Beispiel ist die gekürzte Zeile des Vorgängertexts von `Publish`
(Zwei Entscheidungen, eine Anforderungsreihe mit „ff.“, gegen `git show 1021f6fe:…/changestream.go` gelesen) und trägt
seine Kennungen absichtlich, das „Richtig“-Beispiel trägt genau einen Anker. Es liegt in einem Textblock der Markdown-Datei,
nicht im Go-Kommentar; das Werkzeug liest es nicht, die Regel deckt es auch nicht. Die Reihenfolge der Absätze in §3.7
(Falsch/Richtig der Kommentar-Regel → „Zustandsfelder ebenso“ → Herkunfts-Absatz → „Begründung“) ist nach der Fixrunde
wieder zusammenhängend (F-4, `git diff` gelesen).

**`.claude/commands/implement-slice.md` Schritt 20.** Aufruf `make kommentar-kennungen DIFF=<Basis>` mit dem Hinweis auf
`git add` neuer Dateien, der Chronik-Kandidatenlauf bleibt daneben, die Grenze (Form, nicht Wahrheit; Verweis auf „Grenze
dieser Selbstprüfung“, die es an Zeile 243 gibt) ist genannt; die Nummern 20 und 21 sind unverändert; der Herkunfts-Block
steht hinter der Selbstprüfungs-Grenze (F-4).

**`.harness/skills/reviewer.md`.** Neuer MEDIUM-Unterpunkt: die Form (Reihung, Kompaktform, „ff.“) ist MEDIUM, eine falsch
wiedergegebene Spec-Aussage oder ein nicht getragenes Verhalten geht unter den bestehenden HIGH-Punkt „Kommentar trägt keine
der Kommentar-Klassen“; der Lauf ist ausdrücklich Probe, kein Beleg; Bestandskandidaten außerhalb des Diffs sind nicht
Gegenstand. Die Einstufung ist begründet und doppelt den HIGH-Punkt nicht. Die Herkunft nennt
`slice-chronik-in-code-kommentar`, nicht das erwartete Register-Verzeichnis `kommentar-herkunft-als-kette`: der Plan
führt das als Meldung (F-9, Closure des Planner legt das Verzeichnis an) — kein Befund des Verifiers.

**[`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md)-Konformität.** Das Werkzeug zählt die **Form** (Zahl
verschiedener Kennungen), nicht die Wahrheit von Prosa; der Vertrag nennt die Grenze im ersten Absatz und in fünf
Grenzpunkten, sagt „kein Gate“ und „keine Ausnahmeliste“ ([`AGENTS.md`](../../AGENTS.md) §3.2, §3.6). Ein „grün“ wird an
keiner Stelle als Konformitätsaussage ausgegeben (Vertrag, `harness/README.md`-Zeile, `AGENTS.md` §3.7, Schritt 20 und
Reviewer-Skill tragen denselben Satz). Der Architect-Verdikt
[`architect-verdict-slice-chronik-in-code-kommentar`](architect-verdict-slice-chronik-in-code-kommentar.md) („kein
Textmuster-Sensor über Satz-Subjekte“) ist eingehalten: `main.go` wertet weder Satz-Subjekt noch Slice-/Wellen-Nummer aus
(`git grep` im Programm nach `slice-`/`welle-`: keine Auswertung; der Vertrag grenzt den Chronik-Lauf ab).
Erstbeleg-Anker: [`ADR-0060`](../plan/adr/0060-grpc-streaming-mechanismus.md) (einziger Anker im Godoc von `Publish`).

## 6. Erstbeleg, Handbuch, README-Zeilen

**Erstbeleg `changestream.go`.** Der Diff der Datei, gefiltert auf Zeilen, die weder `+++`/`---` noch mit `//` beginnen:
**0** (Befehl: `git diff 1021f6fe..HEAD -- <Datei> | grep -E '^[+-]' | grep -vE '^(\+\+\+|---)' | grep -vE '^[+-]\s*//'`,
Zeilenzahl 0). Die Kennungen sind auf einen Anker je Block reduziert (elf Zeilen mit Kennung vorher, drei nachher, alle
`ADR-0060`). `Publish` trägt die Zusagen des Plans vollständig: verteilt an die zum Aufrufzeitpunkt registrierten
Abonnenten, blockiert nie, begrenzte Warteschlange je Abonnent, Überlauf für ihn verworfen (Drop-Newest), keine
Zustellgarantie, Fehler nur aus ungültigem Aufruf oder beendetem `ctx`, ohne Einfluss auf Persistierung und
Bestätigung, Nachvollziehbarkeit über den Lesezugriffsweg (`ChangeStorePort`, View `cdc.changes` — beide existieren,
`git grep` bestätigt Fundstellen in `changestore.go` bzw. `readchanges.go`). `ErrChangeStream` und `ErrNotify`: die
Aussage „ordnet sich keiner Fehlerklasse zu“ ist wahr gegen `changenotification.go` (dort „Fehlerklasse `transient`“) und
`broadcaster.go` (Fehler nur bei `nil`-Change, `%w: Publish ohne Change`); der Rang-Zeiger nennt die Datei (F-3).

**Handbuch-Kandidatenlauf.** Keine Betreiber-Oberfläche berührt: der Diff nennt weder `docs/user/`, `compose.yaml`,
Konfigurationsvariablen (`CDC_*`), Schema noch eine API; `make kommentar-kennungen` ist ein Harness-Werkzeug ohne
Betreiber-Bezug. Das Benutzerhandbuch bleibt unberührt, wie der Plan es sagt.

**`harness/README.md`-Zeilen gegen das reale Makefile.** `make help` druckt beide Ziele (`kommentar-kennungen`,
`test-kommentar-kennungen`); die Zeile nennt den Aufruf `[PATHS] [COUNT] [TESTS] [DIFF]` wie das Makefile, „Exit 1 bei
mindestens einem Kandidaten (über `make` als Exit 2)“ — gefahren (`make kommentar-kennungen PATHS=internal/domain/errors`:
Kandidatenzeilen, `make: *** … Fehler 1`, Make-Exit 2). Beide Ziele stehen in `.PHONY`, keines in einem Gate-Bündel
(`make gates` ohne das Werkzeug grün).

**Keine Prozess-/Host-Werkzeug-Verstöße im Diff.** `git diff 1021f6fe..HEAD | grep -nE '^\+.*(sed -i|perl -pi)'` trifft
genau eine Zeile: Prosa im Review-Report („Keine Spur von `sed -i`/`perl -pi` in hinzugefügten Zeilen“), kein Kommando.
Host-Werkzeuge des Werkzeugs (`bash`, `git`, `docker`) sind im Vertrag und in der `harness/README.md`-Zeile genannt; im
Programm und in den Skripten keine Chronik-Sprache (`grep` nach „früher“, „vorher“, „bisher“, „slice-“, „welle-“: ein
Treffer, ein sachlicher Satz über einen Hunk).

## 7. Findings des Reviews — Auflösung am Ist-Zustand

| Finding | Kategorie | Stand am Ist-Zustand | Beleg aus meinem Lauf |
|---|---|---|---|
| F-1 | MEDIUM | **aufgelöst** | Aufrufer pinnt `--src-prefix=a/ --dst-prefix=b/ --no-ext-diff --no-textconv --no-color` (Skript gelesen), `parseDiff` endet bei `+++` ohne `b/` mit Exit 2 (M3 rot); Reproduktion `DIFF=HEAD~40 COUNT=1` mit `diff.mnemonicPrefix=true`: **39**, gleich der Zahl ohne Konfiguration |
| F-2 | LOW | **aufgelöst** | `run-kommentar-kennungen-tests.sh` (198 Zeilen): Stub-`docker`, Zweige, Injektion; `make test-kommentar-kennungen` Exit 0; S1–S4 rot |
| F-3 | LOW | **aufgelöst** | Godoc von `ErrChangeStream` widerspruchsfrei und wahr gegen `broadcaster.go`/`changenotification.go` (§6) |
| F-4 | LOW | **aufgelöst** | Bezugspaare wieder zusammen (§5) |
| F-5 | LOW | **als Träger-Meldung an den Planner geführt** | `observations/BEO-PGC/zitat-nennt-die-falsche-stelle/evidence/welle-d-check-verkoerperung.md:28` zitiert weiter „`AGENTS.md:451-455`“; `grep -n '^### 3.13' AGENTS.md` **480** (Parent `1021f6fe`: **449**); die Meldung steht im Plan §3 mit Frist „Closure dieses Slice“ |
| F-6 | INFO | **aufgelöst** | Zusammensetzung der Diff-Zahlen im Plan (Absatz „Zusammensetzung der `diff`-Zahlen“); Größenordnung von Hand bestätigt (§3) |
| F-7 | INFO | **aufgelöst** | Stichprobe und „strenger als der Baseline-Wortlaut“ im Vertrag §Belege der Definition (§5) |
| F-8 | INFO | **aufgelöst** | M6 (`.git`-Ausschluss) ist rot |
| F-9 | INFO | **als Meldung geführt** | beide Stände der Basis-Messung mit Ursprung im Plan (§3); Register-Verzeichnis `kommentar-herkunft-als-kette` legt die Closure an |

## 8. Verifier-Findings

| ID | Kategorie | Befund | Verifizierbar |
|---|---|---|---|
| V-1 | INFO | Die Zeilenzahl der Suchlauf-Zeile „`Zeile 266`“ ist an beiden Ständen 1, weil sie eine Datei in `open/` (`slice-harness-fmt-check`) trifft; die Zeile belegt die **Meldung**, nicht die Nachziehung: bis der Planner den Träger nachzieht, bleibt „Zeile 266“ dort ein falscher Lokator (der Skill trägt die Zeile jetzt auf 286). Frist und Adresse stehen im Plan (Closure). | `git grep -n 'Zeile 266' -- docs/plan/planning` |
| V-2 | INFO | Die Kandidatenzahl vor dem Erstbeleg (600/403/197) ist mit keinem Repo-Befehl nachmessbar, weil der Parent des Erstbelegs das Programm nicht trägt; der Plan weist sie als solche aus, kein Widerspruch zu [`AGENTS.md`](../../AGENTS.md) §3.12. Die Nachmessung am Diff-Stand (597/400/197) ist wiederholbar. | Plan §3 Basis-Messung |

Kein HIGH, kein MEDIUM, kein LOW aus Verifier-Sicht. Beide INFO sind Planner-Arbeit bei der Closure bzw. bereits im Plan
benannt.

## 9. Verdikt

**DoD trägt:** ja. Sieben `[x]`-Zeilen sind am Ist-Zustand belegt (die Läufe sind meine, nicht die des Implementers oder
des Reviewers); fünf `[ ]`-Zeilen sind korrekt offen (Closure, Planner). Der Slice ist aus Verifier-Sicht
closure-fähig.

**Zusammenfassung der Gates:** `make test-kommentar-kennungen`, `make test`, `make a-check`, `make coverage-gate`
(85.00 %), `make gates`, `make commit-traceability`, `make doc-commits`, `make doc-immutable`,
`make suchlauf-nachmessen` (16 Zeilen) — alle Exit 0; gofmt ohne Ausgabe; acht Klassen von Eingabeseiten-Mutationen
rot; Basis-Messung 597/400/197 (Nenner 1474) nachgemessen; Selbstanwendung und Erstbeleg 0 Kandidaten.

**Übergabe an den Planner (Closure):**
1. Träger-Meldungen mit Frist „Closure dieses Slice“: `slice-harness-fmt-check` („Zeile 266“ → aktueller Lokator),
   `slice-code-kommentare-bereinigung` §1 (Prototyp-Zahlen 396/192 gegen 400/197 gemessen), Register-Beleg
   `zitat-nennt-die-falsche-stelle/evidence/welle-d-check-verkoerperung.md:28` (`AGENTS.md:451-455`, F-5),
   `BEO-PGC/slice-chronik-in-code-kommentar/state.md` (Verweis auf das Werkzeug samt Abgrenzung).
2. Neues Register-Verzeichnis `kommentar-herkunft-als-kette` (Beleg dieser Slice) und die Herkunft im Reviewer-Skill danach
   nachziehen (F-9).
3. Risiko-Ausgänge §6 (sieben), Closure-Notiz, drei Paarungen; die Stichproben-Zahlen (25 Verstoß / 5 Grenzfall von 30,
   Review: vier/vier von acht) sind das Material für den Ausgang des ersten Risikos.
4. Finding-Klassen des Reviews in die Closure §7.

Dieser Report ist ein **Lauf-Beleg** (dieser Stand, dieser Lauf); er ersetzt weder Review noch Closure.
