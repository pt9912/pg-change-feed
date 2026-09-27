# Verifikations-Report: slice-transformationen-e2e-abhilfe — 2026-09-27

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). Frischer Kontext, kein Vertrauen in
Implementer- oder Reviewer-Berichte. Formvorbild:
[`verifikation-slice-start-vorlauf-grenze.md`](verifikation-slice-start-vorlauf-grenze.md).

**Gegenstand:** `git diff e6c5d087..HEAD` — Lifecycle `c674d637`/`11195110`/`e6c5d087`
(vor dem Diff-Beginn liegend, nur `e6c5d087` als Basis), Implementierung `c883a5bf`
(neuer Go-Test `TestE2ETransformationRuleNotApplicableEndsCaptureWithSchemaClass`, neue
Runner-Phase „Transformationen-Nichtanwendbarkeit und Abhilfe", `docs/user/e2e-abdeckung.md`
als Erzeugnis, `harness/README.md`-Nachzug, zwei Kommentar-Korrekturen), Review `58e684fc`
(Report + DoD-Checkbox-Nachzug). 6 geänderte Dateien im Diff (5 vom Implementer + der
Review-Report selbst).

**Ablage:** Drei eigene Mutationen liefen **direkt am Arbeitsbaum** (kein `git worktree`
nötig — Docker-only, `make test` mit Race-Detector, `internal/domain/model/transformation.go`
und `internal/adapters/driving/replication/mapper/mapper.go`): je Mutation Sicherungskopie
zuerst (`cp` nach Scratchpad), Mutation per Write-Werkzeug bzw. `sed 'Ausdruck' Datei >
Scratch-Kopie` gefolgt von `cp Scratch-Kopie Datei` (kein `sed -i`, keine Umleitung direkt auf
eine Repo-Datei), Rücknahme per `cp`/`git checkout` aus der Sicherung, nach jeder Mutation
`md5sum`/`git status --porcelain` gegen den Originalzustand geprüft — leer nach jeder Rücknahme.
Kein `docker run … go test …` außerhalb von `make` (§3.1) — jede Mutation lief über den
sanktionierten `make test`-Aufruf.

**Repo-Zustand:** `HEAD` = `58e684fc`. Ich habe in diesem Lauf **nicht gepusht**.

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — [`AGENTS.md`](../../AGENTS.md) §3.9)

| Sensor | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `make gates` | **Exit 0** | `baseline-verify`: v6.9.0 OK, 54 Dateien; `docs-check`: 1360 Datei(en), 0 Befund(e); `commit-traceability` (`HEAD~5..HEAD`): OK, 5 Commits, ohne Struktur-ID; `coverage-gate`: OK — 85,30 % ≥ 80 %; `generated-sync`: OK, byte-gleich; `a-check`: gesamt 0 Befund(e) |
| `.harness/state/gates-passed.diffsha` vs. `tools/harness/working-tree-hash.sh` | **byte-gleich** | `3b26b59cb5cb650220069dacff61f462c13f7f2fce028484fb5eb85b899ee521` = `3b26b59cb5cb650220069dacff61f462c13f7f2fce028484fb5eb85b899ee521` |
| `make test` (Race, ganzer Baum, unmutiert) | **Exit 0** | alle Pakete `ok` (vor und nach den drei Mutationsläufen, je zur Bestätigung der Rücknahme) |
| `make fmt-check` | **Exit 0** | „260 Go-Dateien geprüft, alle formatiert" |
| `make kommentar-kennungen DIFF=e6c5d087` | **Exit 0** | keine Ausgabe, kein Kandidat |
| `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-transformationen-e2e-abhilfe.md` | **Exit 0** | „8 Zeilen stimmen" — alle vier Blöcke, je zwei Stände (`e6c5d087`, `diff`), `OK` |
| `make doc-immutable RANGE=e6c5d087..HEAD` | **Exit 0** | `d-check: 1360 Datei(en) geprüft, 0 Befund(e)` — Target **existiert real** (`d-check.mk` Zeile 38–39, per `harness/mk/doc-gate.mk` → `include d-check.mk` eingebunden); siehe §2 zum Reviewer-Fehler dazu |
| `make commit-traceability RANGE=c674d637~1..HEAD` | **Exit 0** | „5 Commit(s) …, Betreffs ohne Struktur-ID" (nicht 4, siehe §2 V-2) |
| **3 eigene Mutationen** an `Transformation.CheckApplicable`/`Assembler.RemoveTransformation` | alle wie erwartet rot, eine mit **relevantem Nebenbefund** | siehe §3 |
| Dangling-Docker-Volumes vor/nach allen eigenen Läufen | **39 / 39** | `docker volume ls -qf dangling=true \| wc -l`, kein `prune` |
| `git status --porcelain` nach jeder Mutation/Rücknahme und am Laufende | **leer** | Arbeitsbaum sauber |

**Übernommen, nicht erneut gefahren:** ein realer, voller `make test-integration`-Lauf
(≈8 Min. Compose-Last). `git diff --name-only 58e684fc..HEAD` ist leer — seit dem
Review-Commit wurde nichts mehr geändert; Implementer **und** Reviewer haben je einen
eigenen, grünen `make test-integration`-Lauf mit zitierten Ausgabezeilen (Sentinel-Text,
`error_class=schema`, `pending`→`applied`, `docs/user/e2e-abdeckung.md` unverändert)
gefahren. Ich habe keinen konkreten Zweifel an einer bestimmten Aussage dieses Laufs — der
dritte Lauf hätte keine neue Erkenntnis gegenüber den bereits vorliegenden, wörtlich
zitierten Belegen gebracht. Anker: Review-Report §„Eigene Messungen" (Zeilen zum realen
`make test-integration`-Lauf).

## 2. Reviewer-Faktenfehler zu `make doc-immutable` — eigener Verifier-Fund (V-1, LOW)

Der Review-Report behauptet unter INFO: „**`make doc-immutable` existiert in diesem Repo
nicht**" — Beleg: `grep -rn immutable Makefile harness/mk/*.mk harness/README.md` liefere 0
Treffer.

**Das ist falsch, real widerlegt:**

```
$ grep -n "doc-immutable" d-check.mk
38:.PHONY: doc-immutable
39:doc-immutable: ## Doc-/ADR-Immutabilität via git-Diff (Modul vcs); RANGE=base..head oder STAGED=1 (DC-FA-VCS-001)
$ grep -n "d-check.mk" harness/mk/doc-gate.mk
4:include d-check.mk
$ make doc-immutable RANGE=e6c5d087..HEAD
docker run --rm --network none … --enable vcs --disable links --disable anchors … --range e6c5d087..HEAD
d-check: 1360 Datei(en) geprüft, 0 Befund(e)
$ echo $?
0
```

`d-check.mk` liegt am Repo-Wurzelverzeichnis (nicht unter `harness/mk/`) und wird über
`harness/mk/doc-gate.mk`s `include d-check.mk` eingebunden. Der Reviewer-Grep
(`Makefile harness/mk/*.mk harness/README.md`) deckt genau diesen Pfad **nicht ab** — sein
Befehl lieferte ein wahres Ergebnis (0 Treffer in den drei genannten Dateien/Mustern), aber
dieses Ergebnis **trägt die Schlussfolgerung „das Target existiert nicht" nicht**: Andere
Verifikations-Reports desselben Repos (`docs/reviews/verifikation-slice-transformationen-kern-rename.md`
Zeile 38, `docs/reviews/verifikation-slice-transformationen-antragsweg-schema.md` Zeile 41)
haben `make doc-immutable` bereits vor diesem Slice real und erfolgreich gefahren (1209 bzw.
1219 Dateien, 0 Befunde) — das Target ist kein neuer Zuwachs, sondern länger etablierter
Bestand, den ein vollständigerer Suchraum (`d-check.mk` inklusive) sofort gefunden hätte.

**Einordnung nach `AGENTS.md` §3.13/Beobachtungs-Register:** Dies ist eine Instanz von
[`BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`](../plan/planning/observations/BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht/observation.md)
(Form **Befehl**, aktuell 16× verkörpert, Deckel bei 14×/Schwelle erreicht seit
`welle-backfill-bestand`): ein genannter Beleg-Befehl liefert ein wahres Teilergebnis, das
die behauptete Aussage nicht trägt, weil Suchraum/Bereich unvollständig war — exakt die Form,
die dieses Register bereits für `git diff --name-only` ohne Pathspec (`slice-084`) und
`go list -f '{{len .TestGoFiles}}'` ohne `XTestGoFiles` (`slice-085`) führt. Der Deckel des
Eintrags ist erreicht (≥14×, verkörpert seit `welle-20`/`welle-backfill-bestand`, HIGH-Punkt
im Reviewer-Skill bereits eingezogen); dieses 17. Auftreten braucht deshalb **keine eigene
`evidence/`-Datei**, sondern wird — regelkonform mit `state.md` dieses Eintrags — mit
Finding-Kennung in der Closure-Notiz des Slice genannt (§7 „Was hat funktioniert"/Lerneintrag).

**Schwere:** LOW, nicht HIGH/MEDIUM — die falsche Behauptung stand als **INFO**, hat keinen
Merge-blockierenden Effekt gehabt (der Reviewer hat die Immutability trotzdem korrekt über
den trivialen Ersatzweg „`ADR-0112` nicht im Diff" bestätigt) und berührt keine DoD-Zeile
dieses Slice. Sie ist trotzdem ein realer Sorgfaltsmangel: Der Skill verlangt, einen
genannten Beleg zu **fahren**, nicht seine Abwesenheit aus einem zu engen Suchraum zu
schließen — dieselbe Klasse, die das Register bereits mehrfach für andere Suchbefehle
belegt.

**V-2 (INFO, nebenbei gefunden):** Der Review-Report nennt für
`make commit-traceability RANGE=c674d637~1..HEAD` das Ergebnis „4 Commits" — mein eigener
Lauf mit exakt demselben `RANGE`-Wert liefert **5 Commit(s)** (`git log --oneline
c674d637~1..HEAD` zeigt `c674d637`/`11195110`/`e6c5d087`/`c883a5bf`/`58e684fc`). Mit
`RANGE=c674d637..HEAD` (ohne `~1`) ergäben sich tatsächlich 4. Vermutlich ein
Kopier-/Tippfehler zwischen dem im Report gezeigten Befehl und der tatsächlich gefahrenen
Variante — inhaltlich ohne Wirkung (alle fünf bzw. vier Commits tragen `(ADR-0112)`, keiner
`SPEC-*`/`ARC-*`), aber ein weiteres kleines Beleg-Zahl-Wackeln im selben Report. Kein
eigener Registereintrag nötig (zu geringfügig, keine Fehlschlussfolge), aber zusammen mit
V-1 ein Muster: dieser Reviewer-Lauf hat zweimal eine Zahl/Existenzaussage ungenau geführt.

## 3. Drei eigene Mutationen gegen den Domänenmechanismus (Schwerpunkt 2)

**Zuerst wörtlich gelesen:** Der Review-Report behauptet, `TestConsumeRuleTargetCollisionAfterCompatibleExtension`
und `TestConsumeRuleRemedyRestoresCapture` (beide in `internal/adapters/driving/replication/mapper/transformation_test.go`)
„tragen im Kommentar eine benannte, rot färbende Mutation" und zitiert dazu einen
Kommentartext. **Das Zitat gehört tatsächlich zu einer dritten, anderen Testfunktion**
(`TestConsumeRuleColumnMissingInRelationIsNotApplicable`, Zeilen 174–178) — weder
`TestConsumeRuleTargetCollisionAfterCompatibleExtension` (Zeilen 251–255) noch
`TestConsumeRuleRemedyRestoresCapture` (Zeilen 299–302) tragen selbst einen
Mutations-Kommentar; ich habe alle drei Kommentarblöcke wörtlich gegengelesen. Das ist
dieselbe Fehlerklasse wie V-1 (Beleg trägt die Aussage nicht) — hier: das zitierte Textstück
belegt nicht, was über die zwei genannten Tests behauptet wird.

Um zu prüfen, ob die zugrunde liegende Behauptung („der Domänenmechanismus ist bereits
mutationsgebunden") trotz des Fehlzitats **inhaltlich** stimmt, habe ich drei eigene
Mutationen gegen `internal/domain/model/transformation.go` (`Transformation.CheckApplicable`)
und `internal/adapters/driving/replication/mapper/mapper.go` (`Assembler.RemoveTransformation`)
gefahren — jede über den vollen `make test`-Lauf (Race-Detector):

| # | Mutation | Erwartet (laut Review-Zitat) | Real gesehen |
|---|---|---|---|
| M-1 | Zielnamen-Kollisionsprüfung in `CheckApplicable` entfernt (zweiter `if`) | `TestConsumeRuleTargetCollidesWithRelationColumnIsNotApplicable` und `TestConsumeRuleTargetCollisionAfterCompatibleExtension` rot | **3 andere Tests rot** (`TestConsumeEveryRuleIsCheckedForApplicability`, `TestExecuteInapplicableRuleEndsRunAsSchema`, `TestCheckApplicable`) — **weder** `TestConsumeRuleTargetCollidesWithRelationColumnIsNotApplicable` **noch** `TestConsumeRuleTargetCollisionAfterCompatibleExtension` färben sich. Grund (real gelesen, `internal/domain/model/rowimage.go:54-64`): `BuildRowImage` trägt eine **zweite, unabhängige** Kollisionsprüfung (`containsName(columns, key)` beim Aufbau des Row Image), die denselben `ErrTransformationTargetCollides`-Sentinel liefert und über `mapper.imageError` (`mapper.go:600-605`) auf denselben `ErrTransformationNotApplicable` abgebildet wird — Verteidigung in der Tiefe, aber sie bedeutet: dieser konkrete Mutant an `CheckApplicable` wird von den beiden genannten Tests **nicht** gefangen, weil der redundante Pfad ihn maskiert. |
| M-2 | Spalten-Existenz-Prüfung in `CheckApplicable` entfernt (erster `if`) | — (Review nennt diese Mutation nicht) | **4 Tests rot** (`TestConsumeMapValueApplicabilityHangsOnTheColumnNotTheValue`, `TestConsumeRuleColumnMissingInRelationIsNotApplicable`, `TestConsumeRuleRemedyRestoresCapture`, `TestConsumeEveryRuleIsCheckedForApplicability` in `mapper`; `TestExecuteInapplicableRuleEndsRunAsSchema`, `TestExecuteRuleFailureIsReportedInTheRunOnly` in `backfill`; `TestCheckApplicableMapValue`, `TestCheckApplicable` in `model`) — **hier** fällt `TestConsumeRuleRemedyRestoresCapture` real (dieser Test nutzt eine Relation ganz ohne die Quellspalte, trifft also den ersten `if`, nicht den zweiten). |
| M-3 | `Assembler.RemoveTransformation` zum No-op gemacht (`withoutTransformation`-Aufruf entfernt) | `TestConsumeRuleRemedyRestoresCapture` rot (Abhilfe-Mechanismus selbst) | **Rot wie erwartet**: `TestConsumeRuleRemedyRestoresCapture`, plus `TestTransformationListsAreReplacedNotMutated`, `TestSetAndRemoveMapValueOnLiveBinding`, `TestSetAndRemoveTransformationOnLiveBinding` (mapper), `TestProcessAdministrationRequestsSetAndRemoveTransformationTakeEffectLive`, `TestRunStreamAfterAdministrationPassAppliesTheRuleRemovalBeforeTheStreamAssemblesTheFirstTransaction` (bootstrap) — **dieser** Mechanismus ist real und eindeutig gebunden. |

**Bewertung dieses Befunds:** Die Behauptung des Reviewers, dass die zwei genannten
Domänen-Unit-Tests „dieselbe Story wie der neue E2E-Test" mit einer **rot-färbenden
Mutation** tragen, ist in zwei Hinsichten ungenau:

1. **Das zitierte Mutations-Zitat gehört zu einer anderen Testfunktion** — reines
   Beleg-Zitier-Problem (siehe oben, gleiche Fehlerklasse wie V-1).
2. **`TestConsumeRuleTargetCollisionAfterCompatibleExtension` selbst ist gegen die exakte
   Mutation, die die Kollisionsprüfung in `CheckApplicable` entfernt, nicht empfindlich** —
   der Test bleibt grün, weil `BuildRowImage`s eigene, unabhängige Kollisionserkennung densel­ben
   äußerlich sichtbaren Fehler produziert. Das ist **kein Fehlverhalten des Systems** (die
   Verteidigung in der Tiefe ist real vorhanden und für die Anwendungsebene sogar wünschenswert
   — der Nutzer sieht in beiden Fällen `error_class=schema`), aber es bedeutet: Der Test bindet
   nicht spezifisch die Codezeile `t.kind == TransformationRenameColumn && containsName(columns,
   t.to)` in `CheckApplicable`, sondern das **Gesamtverhalten** über zwei redundante Pfade. Wer
   diesen Test als „mutationsdokumentierten Beleg für `CheckApplicable`" liest (wie der
   Review-Report es tut), überschätzt seine Trennschärfe.

**Ist das ein Blocker für diesen Slice?** **Nein.** Der neue E2E-Test
(`TestE2ETransformationRuleNotApplicableEndsCaptureWithSchemaClass`) prüft das
**Gesamtverhalten am System** (Fehlerklasse `schema`, kein Change, Boundary) — und dieses
Gesamtverhalten ist durch die Kombination aus `CheckApplicable` und `BuildRowImage` real und
mehrfach abgesichert; M-1 bis M-3 bestätigen, dass **irgendein** Pfad des zugrunde liegenden
Mechanismus bei jeder der drei Mutationen bricht (nur eben nicht immer der von den zwei
zitierten Unit-Tests). Die vom Plan/DoD verlangte Aussage — „die Nichtanwendbarkeit einer
Regel endet den Erfassungspfad sichtbar mit Klasse `schema`" — bleibt durch M-1/M-2 (Domäne)
und den realen E2E-Lauf (System) gebunden. Der Verzicht des Reviewers auf eine eigene
Text-Mutation am E2E-Test selbst ist dadurch **im Ergebnis vertretbar**, aber seine
**Begründung** (Zitat aus der falschen Testfunktion, Behauptung einer Trennschärfe, die
`TestConsumeRuleTargetCollisionAfterCompatibleExtension` speziell für `CheckApplicable` nicht
hat) ist ungenau. Ich führe dies als **V-3 (LOW)** — ein Beleg-Ungenauigkeits-Fund derselben
Klasse wie V-1, diesmal an einer Mutations-Kommentar-Zuschreibung statt an einem
Existenz-Grep.

Alle drei Mutationen wurden nach dem jeweiligen Lauf per `cp`/`git checkout` zurückgenommen
und gegen die Sicherungskopie per `md5sum` identisch bestätigt; `make test` lief nach jeder
Rücknahme erneut grün (alle Pakete `ok`).

## 4. Runner-Phase — Kriterium (a)/(b) wörtlich am Skript geprüft (Schwerpunkt 3)

`tools/harness/run-integration-tests.sh` (neue Phase „Transformationen-Nichtanwendbarkeit und
Abhilfe", nach `TestE2ESchemaChangeIncompatibleTypeChange`), Reihenfolge real gelesen:

1. `go test -run '^TestE2ETransformationRuleNotApplicableEndsCaptureWithSchemaClass$'` löst die
   Kollision aus, der Container endet real (`env.pool` schreibt direkt, kein `docker exec`).
2. **Kriterium (a):** `ta_running=$(docker inspect --format '{{.State.Running}}' …)`,
   `[ "$ta_running" = "false" ] || bf_fail …` — Container-Stopp wird **vor** jedem weiteren
   Schritt gemessen. Danach der Log-Sentinel-Grep.
3. **Kriterium (b):** *erst danach* `ta_remove_request=$(bf_sql "SELECT
   cdc.remove_transformation(...)")` — ein reiner SQL-Aufruf, der den bereits (Schritt 2)
   bestätigt gestoppten Container voraussetzt, gefolgt von `bf_expect … pending …`.

Die Reihenfolge im Skript-Text ist eindeutig: die `docker inspect`-Messung (Zeile mit
`ta_running=`) steht **vor** dem `bf_sql "SELECT cdc.remove_transformation…"`-Aufruf — keine
Möglichkeit, dass der SQL-Aufruf vor der Running-Prüfung liefe (sequenzielles Bash-Skript,
kein Hintergrund-Prozess dazwischen). Kriterium (a)/(b) ist damit **real als Reihenfolge im
Skript verankert**, nicht nur behauptet. Kriterium (c) misst — wie Plan §1/§2 und Review
korrekt benennen — ausschließlich die **Folge** (`bf_await_applied`, kein zweiter
`schema`-Fehler, `cdc.changes`-Zeile), keine Zeitordnung „applied vor der ersten
Assemblierung"; ich bestätige diese Grenzziehung als zutreffend benannt, kein stiller Gap.

## 5. Escape-Mechanismus (Schwerpunkt 4)

`set -euo pipefail` (Zeile 62) gilt für das ganze Skript; `bf_fail()` (Zeile 2737–2740) ruft
`exit 1`. `trap cleanup EXIT` (Zeile 232) fängt **jeden** Skript-Exit — regulär am Ende wie
bei einem `bf_fail`/`set -e`-Abbruch mitten in der neuen Phase. `cleanup()` (Zeile 227–231)
räumt unbedingt ab (`docker unpause … || true`, `$COMPOSE down -v --remove-orphans … || true`,
`rm -rf "${WAL_TMP:-}"`) — kein Zustand, der einen Fehlschlag der Phase überlebt, weil sie der
letzte Rundlauf des Skripts ist. Reine Lektüre bestätigt die Behauptung eindeutig, kein
eigener Fehlerpfad-Test nötig.

## 6. Suchlauf (§3.13) und Diff-Umfang

`make suchlauf-nachmessen` bestätigt alle acht Zeilen exakt (§1). `git diff --name-only
e6c5d087..HEAD` liefert 6 Dateien: die 5 vom Implementer geänderten plus der Review-Report
selbst (entstanden erst mit `58e684fc`) — konsistent mit der Review-eigenen Zählung „5
geänderte Dateien" (die den eigenen Report nicht mitzählt, wie im Formvorbild bereits als
Konvention etabliert), kein unbenannter Zuwachs.

## 7. Traceability, Lifecycle, Baum

- 5 Commits (`c674d637~1..HEAD`), jeder trägt `(ADR-0112)`, keiner `SPEC-*`/`ARC-*` —
  eigen per `git log --oneline` gelesen, `make commit-traceability` bestätigt „5 Commit(s) …,
  Betreffs ohne Struktur-ID" (siehe V-2 zur abweichenden Reviewer-Zahl).
- Lifecycle-Commits `c674d637`/`11195110`/`e6c5d087` liegen **vor** dem Diff-Bereich dieses
  Verifikationsauftrags (`e6c5d087..HEAD`) — sie sind bereits Teil der Basis, nicht Gegenstand
  dieses Reviews; ihre Form (Move/Ein-Zeilen-Feld) hat der Reviewer bereits geprüft.
- `git status --porcelain`: leer während des gesamten Laufs (auch zwischen den drei
  Mutationen). Kein `docker prune`, dangling Volumes 39/39 unverändert.
- `ADR-0112` selbst ist nicht Teil des Diffs — Immutability real bestätigt (`make
  doc-immutable`, §1/§2).

## 8. DoD — Verdikt je Zeile (§2 des Plans, aktueller Stand)

| # | DoD-Zeile | Verdikt | Bemerkung |
|---|---|---|---|
| 1 | Nichtanwendbarkeit sichtbar, Kriterium (a) (`[x]`) | **bestätigt** | Go-Test wörtlich gelesen (§3.13-Text im Test selbst: Gegenprobe aktiv geprüft, Boundary bestätigt); Runner-Kriterium (a) real per `docker inspect` vor jedem Folgeschritt (§4) |
| 2 | Abhilfe-Akzeptanzkriterium (a)–(d) (`[x]`) | **bestätigt, mit benannter Grenze bei (c)** | (a)/(b) real als Skript-Reihenfolge verankert (§4); (c) korrekt als Folge-Messung gekennzeichnet (Plan §1/§2, `AGENTS.md` §3.12 Instanz B); (d) real über `cdc.changes`-Zeile geprüft (Review-Zitat der Runner-Ausgabe) |
| 3 | Abdeckung/Klassen getragen (`[x]`) | **bestätigt** | `make docs-check` Teil von `make gates` (0 Befunde); `docs/user/e2e-abdeckung.md` als Runner-Erzeugnis, keine achte Fehlerklasse (`ErrIncompatibleSchemaChange`/`ErrTransformationNotApplicable` bereits am Parent nebeneinander, Suchlauf §6 des Plans bestätigt) |
| 4 | `make gates` grün (`[x]`) | **bestätigt** | eigener Lauf, Exit 0, `.harness/state/gates-passed.diffsha` deckungsgleich mit dem aktuellen Arbeitsbaum-Hash |
| 5 | Review durchgeführt (`[x]`) | **bestätigt, mit Einschränkung** | Report vorhanden, 0 HIGH/MEDIUM/LOW; V-1/V-2/V-3 dieses Reports sind eigene Verifier-Funde zur **Sorgfalt** des Reviews (Beleg-Ungenauigkeiten), keine inhaltlichen Lücken am Slice-Gegenstand selbst |
| 6 | §3.13-Suchlauf (`[x]`) | **bestätigt** | §6 dieses Reports, `make suchlauf-nachmessen` Exit 0, „8 Zeilen stimmen" |
| 7 | Doku-Update `harness/README.md` (`[x]`) | **bestätigt** | Sensor-Zeile `make test-integration` nennt den neuen Beleg (eigen per `git diff` gelesen) |
| 8–12 | Closure-Notiz, Beobachtungs-Register, Risiken-Ausgänge, drei Paarungen, Reconciliation-entfällt (`[ ]`/`[x]` gemischt) | **korrekt offen bzw. korrekt entfällt** | Planner-Vorrecht, siehe §9 dieses Reports für die Übergaben |

**Alle sieben `[x]`-Zeilen sind real bestätigt.** Die fünf verbleibenden `[ ]`-Zeilen (bzw.
die eine korrekt als „entfällt" markierte `[x]`) bleiben regelkonform der Planner-Closure
vorbehalten.

## 9. Übergaben an den Planner

1. **Bereit für Closure**, sofern der Planner die verbleibenden `[ ]`-Zeilen abarbeitet:
   Closure-Notiz mit Steering-Loop-Lerneintrag, Beobachtungs-Register-Fortschreibung, §6-
   Risiken-Ausgänge, drei Paarungen (bleiben regelkonform an der Closure von
   [welle-transformationen](../plan/planning/done/welle-transformationen.md) hängen).
2. **§6-Risiko „Die Ordnung … wird als gemessen ausgegeben, obwohl der Lauf nur ihre Folge
   sieht"**: **entfallen/nicht eingetreten** — Plan §1 und §2 kennzeichnen die Grenze bereits
   korrekt beim Schreiben (kein nachträglicher Fund nötig); Ausgang für die Closure-Notiz:
   „vermieden, weil im Plan von Anfang an als Grenze benannt".
3. **§6-Risiko „Kein Mutationslauf an der Bash-Runner-Phase selbst"**: **entfallen/vertretbar**
   — dieser Verifikationslauf hat die Reihenfolge (a)/(b) durch reine, eindeutige Lektüre des
   sequenziellen Bash-Skripts bestätigt (§4); ein eigener Fehlerpfad-Mutationslauf an der
   Bash-Phase selbst wäre eine schwere Compose-Last ohne neue Erkenntnis gegenüber der bereits
   eindeutigen Textlektüre.
4. **Neuer Registereintrag-Beleg (kein neuer Eintrag, 17. Vorkommen eines bestehenden):**
   [`BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`](../plan/planning/observations/BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht/observation.md)
   — der Reviewer hat einen `grep`-Befehl über einen zu engen Suchraum gefahren
   (`Makefile harness/mk/*.mk harness/README.md`, ohne das per `include` eingebundene
   `d-check.mk` am Repo-Wurzelverzeichnis) und daraus fälschlich „das Target existiert nicht"
   geschlossen (V-1, §2). Deckel bereits erreicht (≥14×) — nach `state.md` dieses Eintrags
   braucht dieses Vorkommen **keine eigene `evidence/`-Datei**, sondern eine Zeile mit
   Finding-Kennung in der Closure-Notiz dieses Slice (§7 „Was hat funktioniert" oder ein
   eigener Punkt zu den Verifier-Funden).
5. **Mutations-Zuschreibungs-Ungenauigkeit im Review (V-3, §3):** Kein neuer Registereintrag
   nötig (zu geringfügig für einen eigenen Kandidaten, aber verwandt mit derselben Klasse wie
   V-1) — dem Planner zur Kenntnis: die zwei vom Reviewer zitierten Domänen-Unit-Tests
   (`TestConsumeRuleTargetCollisionAfterCompatibleExtension`,
   `TestConsumeRuleRemedyRestoresCapture`) sind **teilweise**, nicht **vollständig**
   mutationsgebunden gegenüber `Transformation.CheckApplicable` — die Kollisionsprüfung dort
   wird durch `BuildRowImage`s redundante Prüfung maskiert (§3, M-1). Kein Handlungsbedarf für
   diesen Slice; falls ein künftiger Slice `CheckApplicable` isoliert ändert, ist diese
   Redundanz-Maskierung zu beachten.
6. **Kein neuer schwerer Docker-Lauf nötig:** `make test-integration` wurde nicht ein drittes
   Mal gefahren (§1) — `git diff --name-only 58e684fc..HEAD` ist leer, die Begründung ist von
   mir eigenständig verifiziert, nicht nur übernommen.

## 10. Verdikt

**DoD bestätigt: inhaltlich vollständig, keine Fixrunde am Slice-Gegenstand nötig.** Der neue
E2E-Test bindet die Nichtanwendbarkeit real an die Namenskollision (Gegenprobe aktiv geprüft),
die neue Runner-Phase bindet Kriterium (a)/(b) real als Skript-Reihenfolge (nicht nur
behauptet), Kriterium (c) ist korrekt als Folge-Messung gekennzeichnet, der
Escape-Mechanismus räumt unbedingt ab. `make gates`, `make test`, `make fmt-check`, `make
kommentar-kennungen`, `make suchlauf-nachmessen` und `make doc-immutable` liefen alle
eigenständig grün gegen den Arbeitsbaum.

**Reviewer-Sorgfalt (nicht Slice-Inhalt):** Drei eigene Funde (V-1 LOW, V-2 INFO, V-3 LOW) —
alle vom Typ „ein genannter Beleg trägt die Aussage nicht vollständig", keiner davon ändert
das Gesamtverdikt zum Slice, aber V-1 ist ein klarer Faktenfehler (Behauptung der
Nicht-Existenz eines real existierenden, bereits mehrfach genutzten Gate-Ziels, aus einem zu
eng gewählten Suchraum), den ein realer Lauf sofort widerlegt hätte. Empfehlung an den
Planner: das 17. Vorkommen von `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` in der
Closure-Notiz vermerken.

**Mutationen:** 3 eigene Mutationsläufe (M-1 bis M-3, §3), alle real gegen `make test`
(Race-Detector) gefahren, alle Zustände nach Rücknahme identisch zum Original bestätigt
(`md5sum`, `git status --porcelain`, abschließender grüner `make test`-Lauf). M-1 zeigt einen
Nebenbefund (Verteidigung in der Tiefe maskiert eine spezifische Mutation gegenüber zwei
konkret zitierten Tests) — kein Blocker, aber Grund für V-3.

**Gates:** `make gates`, `make test`, `make fmt-check`, `make kommentar-kennungen
DIFF=e6c5d087`, `make suchlauf-nachmessen`, `make doc-immutable RANGE=e6c5d087..HEAD`, `make
commit-traceability RANGE=c674d637~1..HEAD` — alle Exit 0 im eigenen, ungepipten Lauf.
Arbeitsbaum sauber, dangling Docker-Volumes unverändert (39/39), kein `prune`, kein Push.

**Commit-Kennung dieses Reports:** wird mit dem Commit gesetzt, der diese Datei anlegt (siehe
Commit-Message).

Dieser Report ist ein **Lauf-Beleg** (dieser Stand, dieser Lauf) und ersetzt weder Review noch
Closure.
