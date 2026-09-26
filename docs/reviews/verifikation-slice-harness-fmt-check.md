# Verifikations-Report: slice-harness-fmt-check — 2026-09-26

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich + Entscheidungs-Konformität +
Plan-vs-Code-Diff + Gates. Review-Artefakt des Reviewers:
[`review-slice-harness-fmt-check.md`](review-slice-harness-fmt-check.md); Formvorbild dieses Reports:
[`verifikation-slice-code-kommentare-kennungen.md`](verifikation-slice-code-kommentare-kennungen.md).
Die vendored Baseline trägt kein eigenes Verifikations-Template (`v6.9.0` · `regelwerk/`-Bundle: das Verzeichnis
`templates/docs/reviews/` enthält nur `review-report.template.md`, Gerüst per `cp` übernommen und durch die Form des
Formvorbilds ersetzt).

**Gegenstand:** Slice-Plan `slice-harness-fmt-check` (wellenlos, Harness-Querschnitt), `HEAD` = `0f0cca7c`,
Diff-Range `17cb4eb3..HEAD`, 10 Commits, 15 Dateien (+680/−32 einschließlich zweier Lifecycle-Moves und des
Review-Reports). Slice-Inhalt: Lifecycle (`2d053bba`, `0c68671e`, `47926f65`), Implementer-Lauf (`d330c2be` Werkzeug,
`32637052` Format-Commit, `6047a7c1` Träger, `95ff5ab2` Plan-Nachzug), Review-Report (`cfaf4c5f`: 1 HIGH, 2 LOW,
3 INFO), Fixrunde (`dba47857`, `0f0cca7c`). Der Stand ist nicht gepusht (`origin/main` = `17cb4eb3`). Dieser Lauf ändert
weder Code noch Plan noch Doku; er schreibt nur diesen Report. Die Mutationen liefen an Kopien im Scratchpad der Sitzung
(Ersetzung per `sed` nach stdout in eine neue Datei, kein in-place-Werkzeug am Repo); `git status --short` danach leer.

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — [`AGENTS.md`](../../AGENTS.md) §3.9)

Jeder Lauf schrieb in eine Log-Datei; der Exit-Code wurde im selben Aufruf gesondert gesichert und danach gelesen.
Stand aller Läufe ohne Mutation: `HEAD` = `0f0cca7c`, Arbeitsbaum sauber. Je ein schwerer Docker-Lauf zugleich
(`free -m` vorher: 13,3 GB verfügbar).

| Sensor | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `make test-fmt-check` | **Exit 0** | `run-fmt-check-tests: alle Fälle bestanden` |
| `make fmt-check` | **Exit 0** | `fmt-check: 254 Go-Dateien geprüft, alle formatiert` |
| `make test` (Race-Detector) | **Exit 0** | letzte Zeile `ok …/tools/schema/rolloutguard 1.010s`, `grep -c FAIL` = 0 |
| `make a-check` | **Exit 0** | `gesamt: 0 Befund(e)` |
| `make coverage-gate` | **Exit 0** (zweimal) | `coverage-gate: OK — Coverage 85.10% erfüllt Schwelle 80%` (Lauf 1, allein), `… 85.00% …` (Lauf 2) |
| `make gates` (einmal) | **Exit 0** | `baseline-verify: v6.9.0 OK — 54 Dateien` · `d-check: 1267 Datei(en) geprüft, 0 Befund(e)` · `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD"` · `generated-sync: OK` · `gesamt: 0 Befund(e)`; das Log trägt **0** Vorkommen von `fmt-check` |
| `make suchlauf-nachmessen PLAN=<Plan des Slice>` | **Exit 0** | `suchlauf-nachmessen: 10 Zeilen stimmen` |
| `make kommentar-kennungen DIFF=17cb4eb3 COUNT=1` | **Exit 0** | `0` |
| `make commit-traceability RANGE=17cb4eb3..HEAD` | **Exit 0** | `OK — 10 Commit(s) in "17cb4eb3..HEAD", Betreffs ohne Struktur-ID` |
| `make doc-commits RANGE=17cb4eb3..HEAD` | **Exit 0** | `d-check: 1267 Datei(en) geprüft, 0 Befund(e)` |
| `make doc-immutable RANGE=17cb4eb3..HEAD` | **Exit 0** | `d-check: 1267 Datei(en) geprüft, 0 Befund(e)` (der Aufruf ohne `RANGE` endet mit Exit 2 „flag needs an argument: --range“: das Ziel verlangt die Range) |
| Vor-Format-Messung (§3) | Exit 1 | sechs Bestandsdateien genannt |
| Mutationen (§4) | 12 Läufe | alle rot |

**Coverage-Zahl als Lauf-Beleg ([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz A).** Zwei Läufe desselben Stands drucken 85.10 %
und 85.00 %; die Zahl ist an diesem Stand **nicht lauf-stabil** (Plan und Review nennen 85.00 %, jeweils Beleg ihres Laufs).
`tools/` gehört nicht zur gemessenen Fläche und der Slice ändert an `./internal/...` nur Kommentare und Whitespace; die
Schwelle 80 % ist in beiden Läufen erfüllt.

**Hygiene:** dangling Volumes (`docker volume ls -qf dangling=true | wc -l`) vor dem ersten Lauf **34**, nach allen
Läufen **34**; kein `prune`, kein `system prune`. Nicht gefahren: `make test-integration`, `make test-store`,
`make test-replication`, `make bench` (der Diff berührt in `internal/` und `test/` nur Kommentare und Whitespace, weder
Schema noch Compose noch Runner).

## 2. DoD — Verdikt je Zeile (§2 des Plans)

Gezählt am Plan: acht `[x]`-Zeilen, vier `[ ]`-Zeilen, zusammen zwölf (`grep -c '^- \[x\]'` 8,
`grep -c '^- \[ \]'` 4).

| # | DoD-Zeile | Verdikt | Realer Beleg |
|---|---|---|---|
| 1 | Liefer-Punkt 1 — das Werkzeug (Skript, Ziel, Tabellentest, Vertrag, README-Zeile, nicht in `make gates`) | **erfüllt** | `make test-fmt-check` Exit 0; Exit-Codes 0/1/2 des Skripts gefahren (Fälle des Tabellentests, §3: formatiert · unformatiert genannt Exit 1 · Unterverzeichnis · Syntaxfehler Exit 2 · Verzeichnis ohne Go-Datei Exit 2); die drei im Plan genannten Mutationen (Ausgabe → nur Exit-Code; nur oberste Ebene; Zählung entfernt) rot (§4 M1/M2/M3); Vertrag und README-Zeile gelesen (§5–§6); `make gates` ohne `fmt-check` |
| 2 | Liefer-Punkt 2 — die sechs Bestandsdateien, `make fmt-check` Exit 0 | **erfüllt** | `make fmt-check` Exit 0, 254 Dateien; Vor-Format-Messung an `17cb4eb3` meldet genau die sechs Dateien; der Format-Commit `32637052` ist mit `git show -w` auf zwei Kommentare und einen Umbruch reduziert (§3); `make test` Exit 0 |
| 3 | Liefer-Punkt 3 — Schritt 18 und Reviewer-LOW-Zeile | **erfüllt** | beide Stellen im Diff gelesen (§5); Suchlauf-Feld nachgemessen (§3) |
| 4 | `make gates` grün, Exit gesondert | **erfüllt** | eigener Lauf Exit 0 (§1) |
| 5 | Review durchgeführt, Report liegt vor | **erfüllt** | `review-slice-harness-fmt-check` liegt vor (1 HIGH, 0 MEDIUM, 2 LOW, 3 INFO); Finding für Finding in §7 gegen den Ist-Zustand nachgemessen — HIGH aufgelöst, **ein LOW (F-2) nicht aufgelöst** (V-1, §8); kein offenes HIGH oder MEDIUM |
| 6 | §3.13-Suchlauf: Feld in §3, Gefundenes und Nichtgefundenes je Träger, beide Stände, `make suchlauf-nachmessen` nach jeder Fixrunde | **erfüllt** | 10 Zeilen mit dem Werkzeug Exit 0; alle zehn von Hand mit `git grep` nachgefahren (§3); die Träger-Tabelle nennt Gefundenes je Träger; die Prosa unter dem Feld nummeriert die Zeilen falsch (V-2, INFO) |
| 7 | Doku-Update: `harness/README.md` §Sensors, Handbuch unberührt | **erfüllt** | zwei README-Zeilen im Diff; `git diff --stat 17cb4eb3..HEAD -- docs/user` leer |
| 8 | Reconciliation-Register — entfällt | **erfüllt (entfällt)** | Greenfield, keine Datei |
| 9 | Closure-Notiz mit Steering-Loop-Lerneintrag | **korrekt offen** | Plan §7 trägt „*(zu tragen bei Closure)*“ (Planner) |
| 10 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | Closure-Pflicht (Planner); der Diff berührt keine Register-Datei |
| 11 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | Plan §6: die Zeilen tragen „Ausgang: *(bei Closure)*“ (Planner) |
| 12 | Die drei Paarungen | **korrekt offen** | hängen an der Closure des Slice (Planner) |

(Die Tabelle führt die Zeilen in der Reihenfolge des Plans.) Kein `[x]` ohne Beleg, kein `[ ]`, das über die
Rollen-Sequenz hinaus belegt wäre. Register, Risiko-Ausgänge, Paarungen und Closure-Notiz sind Planner-Arbeit und nicht
Teil dieser Prüfung.

## 3. Plan-vs-Code-Diff und Zahlen mit Ursprung ([`AGENTS.md`](../../AGENTS.md) §3.12)

**§1 Ziel und „Ausdrücklich NICHT“.** Eingehalten: kein Eintrag in `GATE_CHECKS` (`git diff 17cb4eb3..HEAD -- Makefile |
grep GATE` leer; das `make gates`-Log trägt kein `fmt-check`); kein Linter und kein `//nolint`; kein schreibender Pfad im
Werkzeug (Mount `:ro`, kein `gofmt -w`, kein `--user`); keine Formatierung anderer Sprachen. Der Trigger der Gate-Aufnahme
steht im Vertrag als Kenntnis („ein weiteres Auftreten, das der Reviewer trotz gelaufenem Schritt 18 findet“).

**§3-Tabelle, Zeile für Zeile gegen `git diff --stat`:**

| Plan-Zeile | Ist im Diff |
|---|---|
| `tools/harness/fmt-check.sh` neu | 70 Zeilen; Exit-Codes 0/1/2, Auswertung der Ausgabe, Zählung mit `! -type d` (Zeile 44), `--network none`, Mount `:ro` |
| `tools/harness/run-fmt-check-tests.sh` neu | 188 Zeilen; zwölf nummerierte Fälle (1–11 samt 8b) |
| `Makefile` update | `fmt-check` und `test-fmt-check` mit `.PHONY` und Hilfe-Zeile; `make help` druckt beide; kein Gate-Eintrag |
| `harness/sensors/fmt-check.md` neu | 124 Zeilen; Vertrag, Aufruf, Exit-Codes, Grenzen 1–6, Mutationstabelle mit 16 Zeilen |
| `harness/README.md` update | zwei Zeilen unter Werkzeuge; „kein Gate“ in beiden |
| `implement-slice.md` Schritt 18, `reviewer.md` LOW | 15 bzw. 7 Zeilen |
| sechs Go-Dateien | `queries.go` 4, `transformation_test.go` 2, `seam_test.go` 2, `log_test.go` 16, `service_test.go` 2, `integration_test.go` 2 Zeilen |

Über den Plan hinaus, im Plan als **Nachzug** und **Fixrunde** benannt und begründet: `gofmt -d`-Aufruf im Vertrag,
Mutationszeilen je Zusage, `ADR-0014`-Ersatz in `queries.go`, Fixrunde F-1/F-3/F-4/F-6. Abweichungen: 254 statt 252 Dateien
(zwei Go-Dateien von `slice-code-kommentare-kennungen`; `git ls-files '*.go' | wc -l` und der Lauf des Werkzeugs drucken
254); der `ADR-0014`-Ersatz in `queries.go` geht über `gofmt` hinaus (der Block trägt danach eine Kennung, `ADR-0124`,
`make kommentar-kennungen DIFF=17cb4eb3 COUNT=1` **0**; `RetentionPolicy` steht in `internal/domain/model/retention.go`,
[`AGENTS.md`](../../AGENTS.md) §3.7 Kopplung). Keine `Accepted` ADR geändert (`make doc-immutable` Exit 0), keine Schwelle
und keine Gate-Konfiguration verändert ([`AGENTS.md`](../../AGENTS.md) §3.5, §3.6). Die beiden Lifecycle-Commits (`2d053bba`
open → next, `47926f65` next → in-progress) sind reine Renames ([`AGENTS.md`](../../AGENTS.md) §3.3).

**Zahlen mit Ursprung.**

| Zahl | Stand und Ursprung im Plan | Mein Lauf |
|---|---|---|
| 254 Go-Dateien | Plan: gemessen am Stand `17cb4eb3` (`git ls-files '*.go'`); 252 am Stand `7b70b34a` | `make fmt-check` druckt 254 (Stand `HEAD` und Kopie von `17cb4eb3`, jeweils 254) |
| sechs Dateien, sechs Hunks, 86 Zeilen | Plan §1 (Stand `7b70b34a`), am Start neu gemessen (`17cb4eb3`) | `git archive 17cb4eb3` in eine Wegwerf-Kopie, `tools/harness/fmt-check.sh <Kopie>` Exit 1 nennt genau `queries.go`, `transformation_test.go`, `seam_test.go`, `log_test.go`, `service_test.go`, `integration_test.go`; `gofmt -d` über diese sechs im Toolchain-Image: `grep -c '^@@'` **6**, `wc -l` **86**, Hunk `@@ -1236,7 +1236,7 @@` in `integration_test.go` |
| Lokatoren `docs/user/e2e-abdeckung.md` | Plan: verschieben sich nicht | `git diff 17cb4eb3..HEAD -- docs/user/e2e-abdeckung.md` leer (0 Zeilen) |
| 85.00 % (Fixrunde-Zeile im Plan) | Beleg des Laufs der Fixrunde | mein Lauf 85.10 % und 85.00 % (§1: nicht lauf-stabil; beide über der Schwelle) |

**Format-Commit `32637052` (nur Kommentare und Whitespace).** Der Token-Vergleich ohne Go-Werkzeug: `git show 32637052 -w
-U0 -- '*.go'` lässt von den sechs Dateien nur drei mit Zeilen stehen — `queries.go` (zwei Kommentarzeilen: `''` → „leerer
Text“, `ADR-0014` → `RetentionPolicy`-Stelle), `integration_test.go` (eine Kommentarzeile: Backtick-Paar → `/004`),
`log_test.go` (vier einzeilige Methoden auf je drei Zeilen; identische Anweisungen, die `{` bleibt am Zeilenende, keine
Semikolon-Einfügung). `transformation_test.go`, `seam_test.go`, `service_test.go` verschwinden unter `-w`: reine
Ausrichtung. Damit ist der Nicht-Kommentar-Nicht-Whitespace-Strom der Dateien vor und nach dem Commit gleich, abgesehen
vom Umbruch in `log_test.go`. `make test` Exit 0. Die drei geänderten Kommentare sind wahr gegen den Code (Review §Eigenständige
Prüfungen; `retention.go` trägt `RetentionPolicy`, `git grep` bestätigt).

**Selbstanwendung.** `make fmt-check` misst den ganzen Baum (254 Dateien) und damit alle Go-Dateien des Diffs; Exit 0.

**Suchlauf-Feld von Hand nachgefahren** (`git grep -n -E …` ohne das Werkzeug, Plan-Datei ausgeschlossen — am Stand
`17cb4eb3` liegt sie unter `open/`/`next/`, deshalb `:(exclude,glob)**/slice-harness-fmt-check.md`; Ist gleich Soll in
allen zehn Zeilen):

| Zeile | Stand | Soll | Ist von Hand |
|---|---|---|---|
| 1 `gofmt` in `.claude/commands .harness/skills` | `7b70b34a` | 4 | 4 |
| 2 „Bestandsdateien außerhalb des Diffs“ | `7b70b34a` | 1 | 1 |
| 3 `fmt-check` in Ablauf-Trägern und `tools` | `7b70b34a` | 0 | 0 |
| 4 „sechs Bestandsdateien\|Bestandsdateien, die“ in `open`/`next` | `17cb4eb3` | 2 | 2 |
| 5 `gofmt` | `diff` | 3 | 3 (`implement-slice.md` Zeilen 183, 185, 186) |
| 6 „Bestandsdateien außerhalb des Diffs“ | `diff` | 0 | 0 |
| 7 `fmt-check` | `diff` | 37 | 37 |
| 8 wie 4 | `diff` | 2 | 2 (`slice-antragsqueue-lesefehler-failed` Zeile 243, `slice-transformationen-e2e-wirkung` Zeile 166) |
| 9 `Bestandsdatei` inkl. `welle-transformationen` | `17cb4eb3` | 3 | 3 |
| 10 wie 9 | `diff` | 3 | 3 (dritter Treffer `welle-transformationen.md` Zeile 305) |

## 4. Mutationen der Eingabeseite (dieser Lauf)

Je eine Kopie von `tools/harness/fmt-check.sh` im Scratchpad (`sed` nach stdout, die Ersetzung wurde auf tatsächliche
Änderung geprüft), `TOOL=<Kopie> make test-fmt-check`-Form (`bash tools/harness/run-fmt-check-tests.sh` mit `TOOL` und
`TOOLCHAIN_IMAGE`), Exit ungepiped in Log-Dateien: jeder Lauf endete mit Exit 1; das Original blieb unverändert
(`git status --short` leer).

| # | Zusage | Mutation | Rote Fälle (gedruckt) |
|---|---|---|---|
| M1 | Abweichung heißt Exit 1, die Ausgabe wertet | Auswertung der Ausgabe → `if false` (nur der Exit-Code des Formatierers) | unformatierte Datei · gemischt · Unterverzeichnis · Eingabe unverändert · Pfad mit Leerzeichen · Symlink |
| M2 | der ganze Baum wird gelesen | `gofmt -l ./*.go` statt `gofmt -l .` | dieselben sechs (der Pfad trägt `./`, das Unterverzeichnis wird nicht gelesen) |
| M3 | leer ist nicht bestanden | Zählschwelle `-eq 0` → `-eq -1` (Zählung wirkungslos) | leeres Verzeichnis · Verzeichnis ohne Go-Datei |
| M4 | Syntaxfehler heißt Exit 2 | Exit-Code des Formatierers ignoriert | Syntaxfehler |
| M5 | Docker-Fehler heißt Exit 2 | Exit des Docker-Aufrufs 125 durchgereicht | Docker-Fehler |
| M6 | lesender Mount | `:ro` aus dem Mount entfernt | Docker-Argument lesender Mount |
| M6b | Arbeitsverzeichnis `/src` | `-w /src` entfernt | Eingabe unverändert · formatierte Datei · gemischt · Pfad mit Leerzeichen · Symlink · unformatierte Datei · Unterverzeichnis |
| M7 | der Container läuft ohne Netz | `--network none` entfernt | Docker-Argument `--network` (Wert `none` fehlt) |
| M7b | dasselbe | `--network host` | Docker-Argument `--network` trägt nicht den Wert `none` |
| M8 | ein Verzeichnis-Argument, ein Mount | `"$dir"` ohne Anführungszeichen | Pfad mit Leerzeichen · Docker-Argument lesender Mount |
| M9 | ein `.go`-Symlink zählt | Zählung `-type f` statt `! -type d` | Symlink mit Endung `.go` |
| M10 | das Image ist `TOOLCHAIN_IMAGE` | festes Image statt der Variable | Docker-Argument Image · Docker-Fehler |

Die sieben geforderten Klassen der Eingabeseite (Ausgabe-Auswertung; `./*.go` statt Baum; Zählung; Docker-Exit;
`:ro`/`-w`; `--network none`; Anführungszeichen um `$dir`) sowie Formatierer-Exit, Symlink-Zählung und Image sind je rot
gesehen; die Zeilen der Mutationstabelle im Vertrag (§Test) stimmen mit den roten Fällen dieser Läufe überein
(Ausnahme: für M6b nennt der Vertrag keine Zeile, die Mutation färbt sieben Fälle rot). Die Grenze des Stub-Falls steht im
Vertrag: er belegt, was der Aufrufer übergibt, nicht, dass der Daemon es einhält.

## 5. Entscheidungs-Konformität und Träger

**Architect-Entscheidung „Werkzeug ohne Gate“** (Register `BEO-PGC/formatierungs-drift-ohne-gate`, `state.md`): eingehalten
— `make gates` unverändert, `make fmt-check` außerhalb; der Vertrag und beide README-Zeilen sagen „kein Gate“; der Trigger
der Gate-Aufnahme steht als **Kenntnis** (Vertrag §Kein Gate), nicht als Umfang. Die Aufnahme als Gate bliebe eine ADR
([`AGENTS.md`](../../AGENTS.md) §3.6, §4).

**[`AGENTS.md`](../../AGENTS.md) §3.1 (Host-Werkzeuge).** Das Skript nutzt kein Host-`gofmt` (`gofmt` läuft im gepinnten
`TOOLCHAIN_IMAGE`, `docker run --rm --network none`, Mount `:ro`); Host-Werkzeuge `bash`, `git` (Repo-Wurzel), `realpath`
und `docker` stehen im Vertrag mit Verweis auf die Klasse „Host-Werkzeug ohne Installation“ (`AGENTS.md` Zeile 80;
`realpath` ist auch bei `suchlauf-nachmessen` so geführt). Im Diff kein `sed -i`, `perl -pi` oder `awk -i`
(`git diff 17cb4eb3..HEAD | grep -nE '^\+.*(sed -i|perl -pi|awk -i)'` trifft nur Prosa im Review-Report, keine
Kommandozeile).

**Schritt 18 in `.claude/commands/implement-slice.md`.** Absatz „Format“ ruft `make fmt-check`, nennt `gofmt -d` als
Korrekturvorlage mit Verweis auf den Vertrag (Rang-Zeiger statt Wiederholung des Docker-Befehls), „Bestandsdateien
außerhalb des Diffs“ entfällt zu Recht (der Bestand ist formatiert, der Lauf misst den Baum), die Grenze „erste, nicht
tragende Linie“ bleibt; die Schritte 19–21 tragen unverändert ihre Nummern (`git diff` gelesen: Hunk endet vor „19.“).
**Reviewer-Skill LOW-Zeile:** die Probe ist `make fmt-check`, Verweis auf den Vertrag.

**`harness/README.md`-Zeilen gegen das reale Makefile.** `make help` druckt `fmt-check` und `test-fmt-check`; die Zeilen
nennen Exit 0/1/2 und „über `make` kommt jeder Exit ≠ 0 als Exit 2 an“ (gelesen; der Make-Exit-Weg ist der des Vorgängerziels
`make kommentar-kennungen`, hier nicht eigens gefahren), „kein Gate“, Host-Werkzeuge und die Fälle des Tabellentests. Beide Ziele stehen mit `.PHONY` im
Makefile, keines in einem Gate-Bündel.

**Handbuch-Kandidatenlauf.** Keine Betreiber-Oberfläche berührt: der Diff nennt weder `docs/user/`, `compose.yaml`,
`CDC_*`-Konfiguration, Schema noch eine API; `make fmt-check` ist ein Harness-Werkzeug. Das Benutzerhandbuch bleibt
unberührt, wie der Plan es sagt.

## 6. Hard Rules

Docker-only (Werkzeug und Test im Container, Host nur `bash`/`git`/`realpath`/`docker`); kein `//nolint`; kein Gate ohne
ADR (keine Gate-Änderung); Exit-Codes ungepiped gesichert (§1); ein Betreff je Commit mit `ADR-0083` und ohne `SPEC-`/`ARC-`
(`make commit-traceability` Exit 0, 10 Commits); der Format-Commit ist ein eigener Commit ohne Werkzeug-Änderung; der Plan
bleibt in `in-progress/` (die Closure ist Planner-Arbeit).

## 7. Findings des Reviews — Auflösung am Ist-Zustand

| Finding | Kategorie | Stand am Ist-Zustand | Beleg aus meinem Lauf |
|---|---|---|---|
| F-1 | HIGH | **aufgelöst** | `git grep -nE 'wäre\|würde\|sonst\|meldete' -- tools/harness/fmt-check.sh tools/harness/run-fmt-check-tests.sh harness/sensors/fmt-check.md` trifft zwei Zeilen im Vertrag (Zeile 10 „Eine gemeldete Datei“ — Partizip, kein Konjunktiv; Zeile 45 „bei Abweichung auf stderr, sonst auf stdout“ — normaler Zweig); der Kommentar des Skripts nennt die Zusage („ein Verzeichnis ohne Go-Datei endet mit Exit 2, nicht mit Exit 0, und ein Exit 0 sagt zu, dass mindestens eine Go-Datei geprüft wurde“); im Vertrag steht „Leer ist nicht bestanden“ im Indikativ |
| F-2 | LOW | **nicht aufgelöst** (V-1, §8) | Plan Zeile 162 ist eine Leerzeile zwischen der Zeile „Schritt-18-Absatz“ und den sechs Go-Datei-Zeilen; die Zeilen 163–172 stehen ohne Kopf; Die Fixrunden-Zeile 169–172 hängt an derselben Kette |
| F-3 | LOW | **aufgelöst** | Fall 11 (Stub-`docker`) trägt `--network none`, `<absoluter Pfad>:/src:ro`, Image; M6, M7, M7b, M10 rot (§4) |
| F-4 | INFO | **aufgelöst** | Zählung `! -type d` (Skript Zeile 44), M9 (`-type f`) färbt Fall „Symlink“ rot; Vertrag Grenzen 5 (`:` im Pfad) und 6 (Symlink) stehen |
| F-5 | INFO | **als Träger-Meldung geführt** | Plan §3 Fixrunde-Absatz nennt `welle-transformationen.md` Zeilen 304–305 mit Frist „Closure dieses Slice“; `git grep -n 'Bestandsdatei'` trifft dort Zeile 305; Zeile 306 („eine der sechs;“) trägt dasselbe Zählwort ohne das Muster (V-3) |
| F-6 | INFO | **aufgelöst** | Kopf des Tabellentests trennt „Container läuft mit `--network none`“ von dem Image-Zugriff des Daemons im Fall Docker-Fehler; Vertrag §Test und Hilfe-Zeile von `test-fmt-check` („Container ohne Netz“) ziehen nach |

## 8. Verifier-Findings

| ID | Kategorie | Befund | Verifizierbar |
|---|---|---|---|
| V-1 | LOW | F-2 des Reviews ist nicht behoben, obwohl die Fixrunde „Tabelle wiederhergestellt“ meldet: die Datei-Tabelle in §3 des Plans endet nach der Zeile „Schritt-18-Absatz“ mit einer Leerzeile (Zeile 162), die sechs Go-Datei-Zeilen und die vier Fixrunden-Zeilen (163–172) stehen danach ohne Kopf und ohne Trennzeile und rendern als Fließtext. Kein Sensor meldet es (`make docs-check` Exit 0). Kein Einfluss auf eine DoD-Zeile. | `sed -n '158,166p' docs/plan/planning/in-progress/slice-harness-fmt-check.md` (Zeile 162 leer) |
| V-2 | INFO | Die Prosa unter dem Suchlauf-Block nummeriert die Zeilen nach der Umordnung falsch: „Muster der Zeilen vier und fünf“ trifft die Zeilen vier und acht, „Zeilen sechs und sieben“ (Nebenwort `Bestandsdatei`) die Zeilen neun und zehn; das Feld selbst stimmt (10 Zeilen, alle von Hand nachgemessen). | Plan §3, Absätze unter dem Suchlauf-Block, gegen den Block gelesen |
| V-3 | INFO | Der Träger `welle-transformationen.md` trägt das Zählwort an zwei Stellen (304–305 „sechs Bestandsdateien“ und 306 „eine der sechs“); Feld und Meldung nennen nur 304–305. Beide Sätze bleiben nach `done/` wahr, die Meldung ist Kenntnis. Die Lokator-Angabe „304–305“ ist eine Zeilenzahl in einem Plan, keine Adresse mit Bestand. | `sed -n '302,308p' docs/plan/planning/welle-transformationen.md` |
| V-4 | INFO | Die Coverage-Zahl dieses Stands ist über zwei Läufe 85.10 % und 85.00 %: der Lauf-Beleg im Plan (85.00 %) gilt für seinen Lauf, nicht als Zustandsgröße (§3.12). | §1 |

Kein HIGH, kein MEDIUM aus Verifier-Sicht. V-1 ist ein Nachzug des Plans (Implementer-Fixrunde oder Planner bei der
Closure); V-2 bis V-4 sind ohne erwartete Aktion.

## 9. Verdikt

**DoD trägt:** ja, mit einem offenen LOW (V-1). Acht `[x]`-Zeilen sind am Ist-Zustand belegt (die Läufe sind meine, nicht
die des Implementers oder des Reviewers); vier `[ ]`-Zeilen sind korrekt offen (Closure, Planner). Aus Verifier-Sicht ist
der Slice closure-fähig, sobald V-1 (Tabelle in Plan §3) nachgezogen ist.

**Zusammenfassung der Gates:** `make test-fmt-check`, `make fmt-check` (254 Dateien), `make test`, `make a-check`,
`make coverage-gate` (85.10 % / 85.00 %), `make gates`, `make commit-traceability`, `make doc-commits`,
`make doc-immutable RANGE=…`, `make suchlauf-nachmessen` (10 Zeilen), `make kommentar-kennungen DIFF=17cb4eb3 COUNT=1` (0)
— alle Exit 0; Vor-Format-Messung an einer Kopie von `17cb4eb3` nennt genau die sechs Bestandsdateien (`gofmt -d`: 6 Hunks,
86 Zeilen); zwölf Eingabeseiten-Mutationen rot; Format-Commit auf Kommentare und Whitespace reduziert; `docs/user/e2e-abdeckung.md`
unverändert.

**Übergabe an den Planner (Closure):**
1. V-1 (Leerzeile im Plan §3) nachziehen bzw. an den Implementer geben.
2. Träger-Meldungen mit Frist „Closure dieses Slice“: `slice-transformationen-e2e-wirkung` §4 und
   `slice-antragsqueue-lesefehler-failed` §4 (Sätze über die „sechs Bestandsdateien“ überholt), `welle-transformationen.md`
   Zeilen 304–306 (Kenntnis, der Satz bleibt wahr; V-3).
3. Risiko-Ausgänge §6, Closure-Notiz, Register `BEO-PGC/formatierungs-drift-ohne-gate` (`state.md`: Werkzeug geliefert;
   Trigger der Gate-Aufnahme bleibt), drei Paarungen.
4. Finding-Klassen des Reviews in die Closure §7.

Dieser Report ist ein **Lauf-Beleg** (dieser Stand, dieser Lauf); er ersetzt weder Review noch Closure.
