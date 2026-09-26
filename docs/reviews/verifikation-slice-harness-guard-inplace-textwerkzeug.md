# Verifikations-Report: slice-harness-guard-inplace-textwerkzeug — 2026-09-26

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich + Entscheidungs-Konformität +
Plan-vs-Code-Diff + Gates. Review-Artefakt des Reviewers:
[`review-slice-harness-guard-inplace-textwerkzeug.md`](review-slice-harness-guard-inplace-textwerkzeug.md); Formvorbild
dieses Reports: [`verifikation-slice-harness-fmt-check.md`](verifikation-slice-harness-fmt-check.md). Die vendored
Baseline trägt kein eigenes Verifikations-Template (`v6.9.0` · `templates/docs/reviews/` enthält nur
`review-report.template.md`, Gerüst per `cp` übernommen und durch die Form des Formvorbilds ersetzt).

**Gegenstand:** Slice-Plan `slice-harness-guard-inplace-textwerkzeug` (wellenlos, Harness-Querschnitt), `HEAD` =
`baff1593`, Diff-Range `e98d419c..HEAD`, 8 Commits, 11 Dateien (+1205/−80 einschließlich zweier Lifecycle-Moves und des
Review-Reports). Inhalt: Lifecycle (`05273900`, `67ba0dd6`, `87491cca`), Implementer-Lauf (`43408474` Guard, Tabellentest,
`MR-003`, Träger; `a0472421` Plan-Nachzug), Review-Report (`209df712`: 1 HIGH, 4 MEDIUM, 2 LOW, 2 INFO), Fixrunde
(`baff1593`); davor `c625b572` (`implement-slice` Schritt 20). Der Stand ist nicht gepusht (`origin/main` = `e98d419c`,
Baum sauber). Dieser Lauf ändert weder Code noch Plan noch Doku; er schreibt nur diesen Report. Alle Mutationen liefen an
Kopien im Scratchpad der Sitzung (literale Ersetzung per `awk` nach stdout in eine neue Datei, kein in-place-Werkzeug am
Repo, der Guard war in der Sitzung aktiv); `git status --short` danach: nur dieser Report.

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — [`AGENTS.md`](../../AGENTS.md) §3.9)

Jeder Lauf schrieb in eine Log-Datei im Scratchpad; der Exit-Code wurde im selben Aufruf gesondert gesichert. Stand aller
Läufe ohne Mutation: `HEAD` = `baff1593`, Arbeitsbaum sauber. Ein schwerer Docker-Lauf zugleich; dangling Volumes
(`docker volume ls -qf dangling=true | wc -l`) vor dem ersten Lauf **34**, nach allen Läufen **34**; kein `prune`.

| Sensor | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `make test-command-guard` | **Exit 0** | `run-command-guard-tests: alle 300 Fälle bestanden`; `grep -cE '^(block\|pass) '` der Tabelle **294**, dazu **6** Sonderfälle (`run_raw`/PATH-Umbau/`judge_block` ohne Tabellenzeile) — die Zahl 300 stimmt |
| dasselbe mit `GUARD=<Guard am Parent e98d419c>` | Exit 1 | **155** Zeilen `FEHLER:` von 300 (`git show e98d419c:.claude/hooks/pretooluse-command-guard.sh` in eine Scratchpad-Datei); davon **0** in der Gruppe „Bestand“ (18 Tabellenzeilen) — die Gruppe ist am Parent grün, wie die Fixrunde sagt |
| dasselbe mit `GUARD=<Guard am Review-Stand a0472421>` | Exit 1 | **76** Zeilen `FEHLER:` — die Zahl der Fixrunde stimmt |
| `make docs-check` | **Exit 0** | `d-check: 1279 Datei(en) geprüft, 0 Befunde` |
| `make gates` (einmal) | **Exit 0** | `baseline-verify: v6.9.0 OK — 54 Dateien` · `d-check: 1279 Datei(en) geprüft, 0 Befund(e)` · `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD"` · `coverage-gate: OK — Coverage 85.10% erfüllt Schwelle 80%` · `generated-sync: OK` · `a-check … gesamt: 0 Befund(e)`; das Log trägt keine Zeile zum Guard (er ist kein Gate) |
| `make suchlauf-nachmessen PLAN=<Plan des Slice>` | **Exit 0** | `suchlauf-nachmessen: 15 Zeilen stimmen` |
| `make doc-immutable RANGE=e98d419c..HEAD` | **Exit 0** | `d-check: 1279 Datei(en) geprüft, 0 Befund(e)` |
| `make commit-traceability RANGE=e98d419c..HEAD` | **Exit 0** | `OK — 8 Commit(s) in "e98d419c..HEAD", Betreffs ohne Struktur-ID` |
| `make doc-commits RANGE=e98d419c..HEAD` | **Exit 0** | `d-check: 1279 Datei(en) geprüft, 0 Befund(e)` |
| `make kommentar-kennungen DIFF=e98d419c COUNT=1` | **Exit 0** | `0` (das Werkzeug liest nur Go; der Diff enthält keine Go-Datei — die Shell-Kommentare von Hand, §6) |
| Konjunktiv-Kandidatenlauf aus `implement-slice` Schritt 20 (der Befehl des Schritts, Diff seit `e98d419c`, `'*.go' '*.sh' '*.awk'`) | 0 Treffer | `git diff -U0 e98d419c -- '*.go' '*.sh' '*.awk' \| grep -nE '^\+.*(//\|#).*(wäre\|waere\|würde\|wuerde\|hielte\|hätte\|haette\|sonst\|statt)'` gibt keine Zeile (F-7 trägt) |
| Mutationen (§4) | 60 Läufe | 57 rot, **3 grün** (V-3, V-4) |

Nicht gefahren: `make test`, `make test-store`, `make test-integration`, `make bench` — der Diff berührt weder Go-Code noch
Schema, Compose oder Runner.

**Live-Beleg der Blockmeldung (Plan §2 Liefer-Punkt 1, in dieser Sitzung).** Der Aufruf
`sed -i s/a/b/ <Scratchpad>/nicht-vorhanden` (Ziel nicht vorhanden) kam als Hook-Fehler zurück, der Aufruf lief nicht; der
Wortlaut: „In-place text tools (sed -i, perl -i, awk -i inplace) rewrite repo files without a trace and are blocked, also on
a scratch copy (AGENTS.md Hard Rule 3.1). Change a file with the Edit/Write tools; to try a change on a copy, write to
stdout: sed s/a/b/ file > /path/to/scratch-copy.“ Er stimmt mit `git grep -n 'rewrite repo files without a trace' -- .claude
tools harness` überein (ein Treffer, `pretooluse-command-guard.sh:68`). Derselbe Lauf führte `sed -n 1p Makefile` aus (kein Block)
und mehrere `git grep`/`grep`-Aufrufe mit `|` im Muster in Anführungszeichen (kein Block) — die quote-bewusste Segmentierung
wirkt live. Der Live-Beleg mit `python3` ist **nicht** erbracht (Auftrag: kein Host-`python3`; V-5).

## 2. DoD — Verdikt je Zeile (§2 des Plans)

Gezählt am Plan: sechs `[x]`-Zeilen (`grep -c '^- \[x\]'` 6), sechs `[ ]`-Zeilen (`grep -c '^- \[ \]'` 6), zusammen zwölf.

| # | DoD-Zeile | Verdikt | Realer Beleg |
|---|---|---|---|
| 1 | Liefer-Punkt 1 — die in-place Formen | **erfüllt** | 300 Fälle grün, 155/76 rot an Parent/Review-Stand (§1); Treffer je Form und Position, Nicht-Treffer der DoD-Liste und Bestandsregeln als Fälle vorhanden; alle acht im Plan genannten Mutationen rot gesehen (§4: Regel entfernt 117 rot · nur `-i` 12 · jedes `sed` 22 · Trenner nicht maskiert 18 · `-exec`-Kopf nicht gelesen 8 · perl-Bündel beliebig 5 · awk ohne `inplace` 2 · Block-Ausgabe entfernt 170); Live-Beleg der Blockmeldung und `sed -n 1p Makefile` (§1); die Sätze „unabhängig vom Ziel“, „Ersatzweg in der Meldung“ und „15 Paketmanager“ gegen den Guard gelesen (`BLOCKED` zählt 15) — Vorbehalt V-1 (der Text der Zeile ist in der Fixrunde umgeschrieben) |
| 2 | Liefer-Punkt 2 — Host-Interpreter auf Repo-Pfaden | **Code und Tabellentest erfüllt, Live-Beleg offen** (Box zu Recht `[ ]`) | Treffer und Nicht-Treffer der DoD-Liste an der Hook-Schnittstelle gefahren (§3); die fünf Mutationen des Plans rot (Regel entfernt · Segment statt Befehl 1 · Pfadzeichen-Bedingung 21 · `./` 1 · Datei-Namen 7); der Live-Beleg mit `python3` steht aus (V-5) |
| 3 | Liefer-Punkt 3 — die Träger (a–d) | **erfüllt** | (a) `MR-003` per `cp` aus der Vorlage, Pflichtfelder vorhanden, Grenz-Zeile führt die Punkte aus Plan §6 Risiko 2, Index-Zeile mit Anker `mr-003` in `harness/conventions.md` (`make docs-check` und `make doc-immutable` Exit 0); (b) Kopfkommentar im Indikativ, keine Slice-/Wellen-Nummer, Zusage · Abgrenzung · Grenze tragen (§6); (c) `AGENTS.md` §3.1 „Durchsetzung“ und Mutationsprobe gelesen, `· seit slice-harness-guard-inplace-textwerkzeug`; (d) `make test-command-guard` in `Makefile` mit `.PHONY` und Hilfe-Zeile, README-Zeile „kein Gate“ mit Host-Werkzeugen; Suchlauf (§3) |
| 4 | `make gates` grün, Exit gesondert | **erfüllt** | eigener Lauf Exit 0 (§1) |
| 5 | Review durchgeführt, Report liegt vor | **erfüllt, mit Vorbehalt** | Report liegt vor (1 HIGH, 4 MEDIUM, 2 LOW, 2 INFO, merge-blockierend); Finding für Finding in §5 gegen den Ist-Zustand nachgemessen — F-1 bis F-7 aufgelöst, mit einem Rest an F-1 (V-3); die Fixrunde selbst (Guard +Maskierer, 300 Fälle) hat kein zweites Review (V-8) |
| 6 | §3.13-Suchlauf, beide Stände, `make suchlauf-nachmessen` | **erfüllt** | 15 Zeilen Exit 0; fünf Zeilen von Hand nachgefahren (§3); Gefundenes und Nichtgefundenes je Träger stehen im Feld |
| 7 | Doku-Update `harness/README.md` §Sensors, Handbuch unberührt | **erfüllt** | eine README-Zeile im Diff; `git diff --stat e98d419c..HEAD -- docs/user` leer |
| 8 | Closure-Notiz mit Steering-Loop-Lerneintrag | **korrekt offen** | Plan §7 trägt „*(zu tragen bei Closure)*“ (Planner) |
| 9 | Reconciliation-Register — entfällt | **korrekt offen** | Greenfield, keine Datei (Box `[ ]` im Plan, Text „entfällt“) |
| 10 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | Closure-Pflicht (Planner); der Diff berührt keine Register-Datei |
| 11 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | sieben Zeilen tragen „Ausgang: *(bei Closure)*“ |
| 12 | Die drei Paarungen | **korrekt offen** | hängen an der Closure (Planner) |

Kein `[x]` ohne Beleg, kein `[ ]`, das über die Rollen-Sequenz hinaus belegt wäre. Register, Risiko-Ausgänge, Paarungen
und Closure-Notiz sind Planner-Arbeit und nicht Teil dieser Prüfung.

## 3. Plan-vs-Code-Diff und Zahlen mit Ursprung ([`AGENTS.md`](../../AGENTS.md) §3.12)

**§3-Tabelle gegen `git diff --stat e98d419c..HEAD`:**

| Plan-Zeile | Ist im Diff |
|---|---|
| `.claude/hooks/pretooluse-command-guard.sh` update | +191/−36, drei Klassen `pkg`/`inplace`/`interp`, eine Meldung je Klasse (gültiges JSON, Test prüft die Form), Kopfkommentar im Indikativ |
| `tools/harness/run-command-guard-tests.sh` neu | 461 Zeilen, 294 Tabellenzeilen + 6 Sonderfälle; Wegwerf-Repo im Temp-Verzeichnis, `GUARD`/`MASKER` übersteuerbar |
| `tools/harness/mask-quotes.awk` neu (Fixrunde) | 40 Zeilen (Plan nennt „davon 12 Kommentar“: `grep -c '^#'` der Datei 12, stimmt) |
| `Makefile` update | Ziel mit `.PHONY` und Hilfe-Zeile, kein Eintrag in `GATE_CHECKS` |
| `MR-003` neu per `cp`, `harness/conventions.md` update | 134 Zeilen bzw. eine Index-Zeile |
| `harness/README.md`, `AGENTS.md` §3.1 update | eine Zeile bzw. der Absatz „Durchsetzung“ und die Mutationsprobe |
| über den Plan hinaus | `.claude/commands/implement-slice.md` (`c625b572` Schritt 20 und Fixrunde) — im Plan-Nachzug als Träger-Zeile benannt |

**Entscheidungen des Plans gegen den Code.** Kein Gate-Eintrag (`git diff e98d419c..HEAD -- Makefile | grep GATE` leer);
kein Sensor über Dateiinhalte; `tools/harness/blocked/go` nicht angelegt (`ls tools/harness/blocked` nicht vorhanden); keine
Änderung der Regel in `AGENTS.md` §3.1 selbst, nur der Absatz „Durchsetzung“ und ein Satz zur Mutationsprobe.
**Abweichung mit Gewicht — die Fixrunde kehrt einen Ausschluss des Plans um (V-1):** Plan §1 „Ausdrücklich NICHT“ führte
ursprünglich „Quote-Bewusstsein im Guard (Shell-Parser)“, Plan §4 Rückführung (a) benannte den Fall als Architect-Frage; die
Fixrunde liefert einen 40-Zeilen-Zustandsautomaten und schreibt §1, §4 und die DoD-Zeile 1 auf den neuen Stand um
(`git diff a0472421..baff1593 -- <Plan>`).

**Suchlauf-Feld — von Hand nachgefahren** (`git grep -n -E …` ohne das Werkzeug, Plan-Datei ausgeschlossen; Suchraum wie im
Feld; Ist gleich Soll in allen geprüften Zeilen):

| Zeile | Stand | Soll | Ist von Hand |
|---|---|---|---|
| 1 `pretooluse-command-guard` | `89d427e0` | 3 | 3 |
| 1 dasselbe | `diff` | 12 | 12 |
| 3 Zählwort `sed -i\|perl -pi\|awk -i` | `89d427e0` | 9 | 9 |
| 3 dasselbe | `diff` | 135 | 135 |
| 5 „in keinem committeten Text“ | `diff` | 1 | 1 (`welle-transformationen.md`; mit Plan-Datei 4) |
| 7 `test-command-guard` | `diff` | 7 | 7 (Guard, `AGENTS.md`, `Makefile` 2, README, `MR-003` 2) |
| 6 überholter Wortlaut („liest er nicht“, „Bewusst NICHT gepr“, „blockt Host-Paketmanager und Host-Toolchains“) | `diff` | 0 | 0 (mit Plan-Datei 7 — der Plan zitiert den Wortlaut als Suchmuster; ohne sie 0) |

Der Stand `89d427e0` des Feldes ist der Parent bei Anlage des Plans; zwischen `89d427e0` und dem Range-Parent `e98d419c`
liegt an den Trägern des Feldes (`.claude`, `tools/harness`, `AGENTS.md`, `harness`) nur `harness/sensors/fmt-check.md`
(`git diff --stat 89d427e0 e98d419c`), das keine Zeile des Feldes bewegt. **Zahlen mit Ursprung:** „300 Fälle: 294 + 6“, „155
rot am Parent“, „76 rot am Review-Stand“, „Bestand 18 Zeilen, am Parent grün“, „`BLOCKED` = 15“: gemessen und bestätigt.
„Latenz 15 ms gegen 11 ms“ (20 Aufrufe je Guard) und „120 Mutationen, 118 rot, 2 äquivalent“ stammen aus dem Lauf des
Implementers; nicht nachgemessen (Latenz laut Auftrag), die Mutationszahl ohne Lauf-Anker (V-6).

## 4. Mutationen der Eingabeseite (dieser Lauf)

Je eine Kopie von Guard bzw. Maskierer im Scratchpad (literale Ersetzung, die Ersetzung wurde auf tatsächlichen Treffer
geprüft; zwei Ersetzungen ohne Treffer wurden mit einem kürzeren Literal wiederholt), `GUARD=<Kopie>` bzw.
`MASKER=<Kopie>` gegen `bash tools/harness/run-command-guard-tests.sh`, Exit ungepiped: **60 angewandte Mutationen, 57 rot,
3 grün**. Das Original blieb unverändert (`git status --short`).

| Zusage | Mutation | Rote Fälle |
|---|---|---|
| sed-Klasse `[nEsrzub]` | je `s`, `b`, `n` entfernt (3) | 2, 2, 3 |
| sed-Klassengrenze | `e` in die Klasse aufgenommen | 1 (`sed -ei`) |
| perl-Klasse `[0-7lanpsw]` | `w`, `l`, `0` entfernt; Obergrenze `7` → `6`; `M` aufgenommen | 2, 3, 4, 2, 1 |
| perl-Ziffernklasse | **`0-9` (Obergrenze geweitet)** | **0 — grün (V-3)** |
| perl-Ziffernklasse | **`0-27` (Ziffern 3–6 fehlen)** | **0 — grün (V-3)** |
| Pfadzeichen-Klasse | je `~`, `-`, `_`, `A-Z`, `0-9` entfernt (5); nur `A-M` (Rest fehlt); `=`, `:` aufgenommen (2) | 2, 2, 2, 2, 3; 1; 2, 1 |
| in-place Regeln | Regel entfernt · nur `-i` · jedes `sed` blockt · Abkürzung ab 9 Zeichen · `gsed` nicht gelesen · Anführungszeichen-Bereinigung entfernt · perl beliebiges `i` · awk `-i` ohne Wert · `-iinplace`/`--include=inplace` nicht gelesen | 117 · 12 · 22 · 2 · 1 · 1 · 5 · 2 · 1 |
| `-exec`-Familie | Rekursion entfernt · `-execdir`/`-ok`/`-okdir` nicht gelesen · find-Ende nicht gelesen | 8 · 3 · 1 |
| Wrapper und Kopf | `sudo` kein Präfix · Optionen nicht übersprungen · xargs-Wert nie übersprungen · `eval` nicht rekursiv · `-c`-Bündel nur `-c` · Tiefe 4 erlaubt | 4 · 21 · 9 · 2 · 2 · 1 |
| `command -v` | Ausnahme entfernt · `-V` fehlt | 3 · 1 |
| Schlüsselwörter, `case` | `until`, `!`, `do` entfernt (3); Case-Label ohne Zustand | 1, 2, 3; 3 |
| Host-Interpreter | Segment statt Befehl · Pfadzeichen-Bedingung entfernt · `./` nicht erlaubt (zwei Formen) · Datei-Namen nicht gelesen · Absolutpfad der Wurzel · Punkt-Verzeichnisse | 1 · 21 · 1, 1 · 7 · 1 · 1 |
| Maskierer | einfache Anführungszeichen nicht maskiert · `$(` in `"…"` maskiert · `\;` nicht Trenner · Zeilenumbruch nicht maskiert · **Rohstring-Fallback entfernt** · Fallback → `out` | 18 · 1 · 1 · 1 · 3 · 3 |
| **Maskierer, Backslash** | **das Zeichen hinter `\` nicht maskiert (`\|`, `\&`, `\ `)** | **0 — grün (V-4)** |
| fail-closed | Maskierer-Ausfall nicht fail-closed · Exit-Code der `awk`-Extraktion ignoriert · Block-Ausgabe entfernt | 1 · 3 · 170 |

Die Mutationen der Zusagen aus Plan §2 sind an den genannten Fällen rot gesehen. F-1 ist für sed-Klasse und Pfadzeichen je
Mitglied gebunden (5 von 5 Pfadzeichen-Mutationen, 3 von 3 sed-Mitgliedern rot); die drei grünen Mutationen stehen in §7.

**Die Alltags- und Umgehungsformen der Review-Findings, an der Hook-Schnittstelle gefahren** (Treiber im Scratchpad, Ausgabe je
Zeile Klasse oder `pass`; Ist wie zugesagt):

- F-2 (nicht geblockt): `git grep -E 'sed -i|perl -pi|awk -i' -- x`, `grep -E "sed -i|perl -pi" f`, `echo 'a|sed -i x'`,
  `echo a\|sed -i x`; (geblockt): `cat f | sed -i …`, `ls; sed -i …`, `x && sed -i …`, `x || sed -i …`, `(sed -i …)`,
  `echo $(sed -i …)`, `echo "$(sed -i …)"`.
- F-3: `env -i pip install x`, `xargs -n1 pip install`, `time -p pip install`, `for f in a; do pip install; done`,
  `command -p pip install x` blocken; `command -v pip`, `command -V pip`, `command -pv pip`, `type pip`, `which pip` passieren.
- F-4: `for … do sed -i`, `while … do`, `if … then`, `! sed -i`, `case x in x) sed -i …;; esac` blocken; `echo $(date) sed -i x`
  passiert (Label nur nach `case`).
- F-6 (gelöst): `sed -i''`, `perl -i'' -pe`, `/usr/bin/env python3 tools/x.py`, `/usr/bin/sudo sed -i`, `busybox sed -i`, `\sed -i`, `gsed -i`,
  `sed --in-p`, `find … -exec sh -c 'sed -i …' {} \;`, `sed -n 1p -input.txt f` blocken; (Grenzen, wie in `MR-003`) `perl x.pl -input a`,
  `sudo -u x sed -i`, `uv run python tools/x.py`, `node -e`, `perl -Wpi`, `x=sed; $x -i …` passieren.
- Mutationsweg nach `AGENTS.md` §3.1: `sed s/a/b/ Makefile > /tmp/copy` passiert.
- Interp: `python3 tools/x.py`, `python3 ./tools/x.py`, `python3 -c "open('Makefile', 'w')"`, `perl -e 'unlink q(Makefile)'`, `python3 <Repo-Wurzel>/docs/a.md`
  und ein `python3`-Heredoc mit `docs/a.md` im Text blocken (Klasse `interp`); `python3 --version`, `python3 -c 'print(1)'`,
  `python3 /tmp/scratch/mutate.py /tmp/scratch/copy.go`, `python3 /tmp/x/docs/a.py`, `perl -ne 'print' /tmp/x` passieren.

## 5. Findings des Reviews — Auflösung am Ist-Zustand

| Finding | Kategorie | Stand am Ist-Zustand | Beleg aus meinem Lauf |
|---|---|---|---|
| F-1 | HIGH | **aufgelöst, mit Rest (V-3)** | sed-Klasse (3 Mitglieder), Pfadzeichen (5) und perl-Buchstaben (4) je rot (§4); der Rest: die perl-Ziffernklasse `0-7` ist an den Ziffern 0, 1, 2, 7 gebunden, nicht an 3–6 und nicht an der Obergrenze (`0-9` und `0-27` bleiben grün) |
| F-2 | MEDIUM | **aufgelöst** | Formen oben; 31 Fälle „Anführungszeichen“ im Tabellentest, live gesehen; **neue Kehrseite V-2** (balancierte Anführungszeichen über Heredoc-Zeilen) |
| F-3 | MEDIUM | **aufgelöst** | die Gruppe „Bestand“ (18 Zeilen) ist am Parent grün (0 `FEHLER:`); die Kopf-Erkennung steht als eigene Gruppe; `command -v/-V/-pv` passieren, `command -p pip` blockt; `MR-003` (Adaption) benennt die Erweiterung der Bestandsregel |
| F-4 | MEDIUM | **aufgelöst** | Formen oben; Grenzen (`sudo -u x`, `env -u X`, `nice -n 10`) in `MR-003` und Kopfkommentar |
| F-5 | MEDIUM | **aufgelöst** | Meldung `interp` und `MR-003` (Begründung) nennen Edit/Write, ein Repo-Werkzeug hinter `make` und `sed … > Kopie` (`git grep`, Guard Zeile 69); kein Host-Interpreter mehr als Weg; ein Tabellenfall prüft die Abwesenheit von `mutate` (Mutation „Meldung mit `mutate.py`“ des Implementers; von mir nicht wiederholt) |
| F-6 | LOW | **aufgelöst** | Formen oben (Treffer wie zugesagt, Grenzen wie benannt) |
| F-7 | LOW | **aufgelöst** | Befehl in `implement-slice.md` Schritt 20 liest `'*.go' '*.sh' '*.awk'` mit `(//\|#)` und transliterierten Formen; mein Lauf auf dem Diff seit `e98d419c`: 0 Treffer (§1) |
| F-8 | INFO | Kenntnis | `MR-003` Begründung nennt die Denylist-Nähe von Liefer-Punkt 2 und die unbelegten Regeln `perl`/`awk` |
| F-9 | INFO | Kenntnis | doppelte Sicherung (Präsenzprüfung und Exit-Code): Einzelmutation äquivalent, beide zusammen rot (Mutation „Exit-Code der `awk`-Extraktion ignoriert“: 3 rot) |

## 6. Entscheidungs-Konformität und Träger

**[`AGENTS.md`](../../AGENTS.md) §3.1.** Absatz „Durchsetzung“: die genannten Blocks (`sed -i`/`--in-place`, `perl -i`,
`awk -i inplace`, unabhängig vom Ziel, hinter `find -exec`, Host-`python`/`perl` mit Repo-Pfad) und Nicht-Lesungen
(Umleitungen, `tee`, Skript im Interpreter, andere Werkzeuge, `cd`, Sprach-Toolchains) sind an der Hook-Schnittstelle wahr
(§4). Der Satz zur Mutationsprobe nennt den Weg über Edit/Write oder `sed … Datei > Kopie`; der Weg passiert den Guard.
`Herkunft` ist ein Feld (`· seit slice-harness-guard-inplace-textwerkzeug`).

**`MR-003`** (Baseline `v6.9.0` · `regelwerk/modul-13-quality-gates.md` §Guard-Härtung: jede Härtung ein `MR` mit
Grenz-Zeile): per `cp` aus `templates/harness/conventions/MR-NNN-titel.template.md`, das Feld „Ersetzt-Baseline-Regel“ nennt
den Satz aus `regelwerk/grundlagen-durchsetzungsschicht.md` §Grenzen — ehrlich benannt mit der Begründung, warum er für alles
Ungelesene weiter gilt; Auslöser mit vier Beleg-Dateien (stimmt mit dem Register überein: Review-Report, Nachmessung dort);
Grenz-Zeile führt die Punkte aus Plan §6 Risiko 2 und die Falsch-Positive; Auflösungs-Trigger permanent. Der Index in
`harness/conventions.md` trägt die Zeile mit Anker `mr-003`; das `structure`-Modul und `make doc-immutable` sind grün.
Die Grenz-Zeile nennt **nicht** die balancierte Kehrseite der Quote-Lesung (V-2).

**Kommentar-Regeln (§3.7), von Hand gelesen** (die Skripte und der Maskierer sind Shell/awk, `make kommentar-kennungen` liest
sie nicht): der Kopfkommentar des Guards (Zeilen 2–58), der Kopf des Maskierers (Zeilen 1–12) und der Kopf des Tabellentests
sind im Indikativ, tragen Zusage · Abgrenzung · Grenze · Rang-Zeiger (`MR-003`, `make test-command-guard`, `AGENTS.md`), eine
Herkunftsangabe je Kommentar, keine Slice- oder Wellen-Nummer (`git grep -n -E 'slice-|welle'` über die drei Dateien: 0 Treffer;
ein Verweis auf `MR-003` steht als Rang-Zeiger). Der Konjunktiv-Lauf ist leer (§1). Eine Vorher-Nachher-Formulierung steht in der **Doku-Prosa** von
`MR-003` (Adaption: „der Bestand las nach einem Wrapper-Präfix keine Optionen und kannte keine Schlüsselwörter“) — INFO V-7.

**Hard Rules.** Docker-only: das Werkzeug ist reines `bash`/`awk`/`mktemp`/`cp`/`grep`/`ln` und schreibt nur ins Temp-Verzeichnis
(`trap`-Aufräumen); kein Host-`go`/`python`. Kein `//nolint`. Kein Gate gelockert oder neu (§3.6: `GATE_CHECKS` unverändert).
`git diff e98d419c..HEAD | grep -cE '^\+.*(sed -i|perl -pi|awk -i)'` trifft 154 Zeilen, alle in Prosa, Kommentaren,
Meldungstexten, Tabellenfällen (Kommandostrings als Testeingabe) und dem Review-Report; kein Aufruf im Diff. Lifecycle-Commits
`05273900` (open → next) und `87491cca` (next → in-progress) sind reine Renames (`git show -M --stat`: 0 Zeilen). Ein Betreff je Commit
mit [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md) und ohne Struktur-ID (`make commit-traceability` Exit 0, 8 Commits).
Handbuch-Kandidatenlauf: der Diff nennt weder `docs/user/`, `compose.yaml`, `CDC_*`-Konfiguration, Schema noch eine API — keine
Betreiber-Oberfläche, das Benutzerhandbuch bleibt unberührt.

**Die Fragen des Implementers an den Architect — Bewertung und Empfehlung (kein Eingriff).**

1. **Quote-bewusst statt quote-blind, mit Rohstring-Fallback** (Plan §4 Rückführung (a)).
   *Semantik-Änderung:* der Guard segmentierte am Parent an jedem Trenner; jetzt maskiert `mask-quotes.awk` Trenner, Leerraum
   und Zeilenumbruch in `'…'`, `"…"` und hinter `\`, lässt `$(`/Backtick in `"…"` lebendig und gibt bei unbalancierten
   Anführungszeichen den Rohstring zurück; fällt der Maskierer aus, blockt der Guard. *Nutzen (gemessen):* der Alltag der Rollen
   — das Suchmuster des Suchlaufs nach [`AGENTS.md`](../../AGENTS.md) §3.13 — blockt nicht mehr (live gesehen; 31 Fälle);
   der Review hatte den Block viermal live erlebt. *Risiko:* (a) neue Falsch-Negative: balancierte Anführungszeichen über
   Heredoc-Zeilen maskieren dazwischenliegende Kommando-Zeilen (V-2, gemessen); (b) ein Lexer als eigene Wartungsstelle
   (40 Zeilen, 20 einzeln rot gesehene Mutationen); (c) was in Anführungszeichen ausgeführt wird, ohne `bash -c`/`eval`/`-exec`
   (`awk 'BEGIN{system(…)}'`, `git -c alias.x='!…'`), war schon am Parent ungelesen (Grenze, nicht neu).
   *Empfehlung:* Antwort auf „fällt ein Zustandsautomat dieser Größe unter ‚kein Shell-Parser‘?“: **ja, als Lexer-Stufe**
   (Zeichenklasse und Verschachtelung von `$(`/Backtick), **nein als Grammatik** (Heredoc, Umleitung, Variablen, Kommando als
   Wert bleiben ungelesen) — der Architect sollte das ausdrücklich ratifizieren, die balancierte Kehrseite (V-2) in die
   Grenz-Zeile von `MR-003` aufnehmen lassen und eine Obergrenze setzen: jede weitere Quote-/Heredoc-Semantik ist der Schritt
   zur Sandbox, nicht zu diesem Guard. Die Umkehr des Ausschlusses aus Plan §1 ohne Rückführung ist V-1.
2. **Erweiterung der Bestandsregel (F-3: Wrapper-Optionen und Schlüsselwörter für alle Klassen).**
   *Nutzen:* schließt am Parent gemessene Umgehungen der Paketmanager-Klasse (`env -i pip`, `xargs -n1 pip`, `time -p pip`,
   `for … do pip`; die Gruppe „Kopf-Erkennung“ ist am Parent rot); eine Kopf-Erkennung statt zwei Wartungsstellen. *Risiko:*
   ein Existenzcheck (`command -v pip`) hätte plötzlich geblockt — durch die Ausnahme für `command -v/-V/-pv` gelöst
   (gemessen: passiert; `command -p pip` blockt); die Optionswerte anderer Wrapper (`sudo -u x pip`, `nice -n 10 pip`) bleiben
   Grenze (gemessen: passiert). *Empfehlung:* **akzeptieren**; die Änderung verschärft (kein Gate gelockert, §3.6), steht in
   `MR-003` (Adaption) und ist als eigene Tabellengruppe gebunden. Formulierungs-Vorbehalt: DoD-Zeile 1 sagt „Bestandsregeln …
   bleiben“ und nennt die Erweiterung in einem Nebensatz (der Text ist in der Fixrunde angepasst, V-1).
3. **Host-`python`/`perl` auf einem Pfad ohne Repo-Namen geht durch den Guard, ist nach §3.1 aber kein genannter Weg.**
   *Befund:* `python3 /tmp/x/docs/a.py` passiert (gemessen); [`AGENTS.md`](../../AGENTS.md) §3.1 führt `python` schon unter
   „läuft im Container“ und verbietet einen Host-Interpreter „auf der Datei“; `MR-003` (Grenz-Zeile) sagt es ausdrücklich
   („kein Weg nach §3.1“). *Empfehlung:* **kein neuer Regelsatz**, aber ein Satz im Absatz „Durchsetzung“ (`AGENTS.md`
   §3.1): ein Host-Interpreter ist auch ohne Repo-Pfad kein zulässiger Weg, der Guard blockt ihn nur auf einem Repo-Pfad
   (Stolperdraht). Die Spannung „Guard lässt durch, was §3.1 verbietet“ ist zugleich das Argument der offenen Nutzer-Frage 1
   (`tools/harness/blocked/go`, die `python3` unbedingt blockte); sie gehört mit dieser Empfehlung in die Antwort auf die Frage.

## 7. Verifier-Findings

| ID | Kategorie | Befund | Verifizierbar |
|---|---|---|---|
| V-1 | MEDIUM | Die Fixrunde kehrt einen ausdrücklichen Ausschluss des Plans um: §1 „Ausdrücklich NICHT“ führte „Quote-Bewusstsein im Guard (Shell-Parser)“, §4 Rückführung (a) benannte den Fall als Architect-Frage mit Rückführung nach `open/`; die Fixrunde lieferte `mask-quotes.awk` und schrieb §1, §4 und die DoD-Zeile 1 auf den neuen Stand um. Der Plan hält die Bestätigung als „offen an den Architect“, ist aber nicht mehr am Text der ursprünglichen Zusage prüfbar. Kein Fehler im Code — eine Steuerungs-Lücke: die Closure hängt an einer Bestätigung, die der Plan nur als Wunsch führt. | `git diff a0472421..baff1593 -- docs/plan/planning/in-progress/slice-harness-guard-inplace-textwerkzeug.md` (Zeilen §1, §2, §4) |
| V-2 | LOW | Balancierte Anführungszeichen über Heredoc-Zeilen erzeugen ein neues Falsch-Negativ: `cat <<EOF`, Zeile mit einem Apostroph, `EOF`, `sed -i s/a/b/ f`, `cat <<EOF`, Zeile mit einem Apostroph, `EOF` — die zwei Apostrophe umschließen die `sed -i`-Zeile, der Guard lässt sie passieren (gemessen; ein Apostroph allein oder zwei in derselben Zeile blocken). Der Kommentar „mehr Segmente, nie weniger“ gilt nur für unbalancierte Anführungszeichen; `MR-003` und Kopfkommentar nennen diese Kehrseite nicht, die Grenz-Zeile nennt nur die Falsch-Positive des Heredocs. Kein Tabellenfall. | `bash <Treiber> "$(printf 'cat <<EOF\ndon'"'"'t\nEOF\nsed -i s/a/b/ f\ncat <<EOF\nwon'"'"'t\nEOF')"` → `pass` |
| V-3 | LOW | Rest von F-1 (Klasse „Zusage ohne Bindung an ihre Eingabeseite“): die perl-Ziffernklasse `0-7` ist an den Ziffern 0, 1, 2 und 7 gebunden (`-0pi`, `-l0127pi`); die Mutationen `[0-27lanpsw]` (Ziffern 3–6 fehlen) und `[0-9lanpsw]` (Obergrenze geweitet) färben keinen Fall rot. Die Fixrunde nennt die Verengung von `0-9` auf `0-7` als Begründung (oktal), der Tabellentest hat kein `-8i`/`-9i` als Nicht-Treffer. Kein Verhaltensrisiko (`perl -8i` ist keine gültige Option), aber die Behauptung „je Klassengrenze ein Nicht-Treffer“ (Fixrunde F-1) ist für diese Klasse nicht ganz getragen. | `MASKER`-freie Kopie des Guards mit `[0-9lanpsw]` bzw. `[0-27lanpsw]`: `GUARD=<Kopie> make test-command-guard` Exit 0 |
| V-4 | LOW | Die Maskierung des Zeichens hinter einem Backslash (`\|`, `\&`, `\ `) ist ungebunden: die Mutation „`mk(d)` entfällt“ färbt keinen Fall rot, obwohl `echo a\|sed -i x` dann blockt (gemessen mit einer Guard-Kopie samt mutiertem Maskierer). Der Fall „Backslash vor \|“ (`grep a\|sed\ -i\|b f`) passiert zufällig — das Kopf-Token trägt den Backslash `sed\`. | `echo a\|sed -i x` am mutierten Maskierer: Klasse `inplace`; am Original: kein Block |
| V-5 | INFO | Liefer-Punkt 2 bleibt ohne Live-Beleg mit `python3`: Plan §2 und Risiko 4 weisen ihn dem Verifier zu, der Auftrag dieses Laufs verbietet Host-`python3` auch für einen Beleg. Getragen ist die Hook-Schnittstelle (Tabellentest, Mutationen, §4) und die Klasse `interp` als Ausgabe: der Block-Pfad `emit_block` ist für alle drei Klassen derselbe (nur die Meldung unterscheidet sich), die Klasse `inplace` ist live belegt (§1); der Review-Report nennt außerdem zwei unbeabsichtigte Live-Blocks der Klasse `interp` in der Sitzung des Reviewers (übernommen, nicht nachgemessen). Die DoD-Box bleibt `[ ]`, bis der Planner den Beleg annimmt oder eine Sitzung ihn erbringt. | Auftrag dieses Laufs; Plan §2 Zeile 2 |
| V-6 | INFO | „120 Mutationen, 118 rot, 2 äquivalent“ (Plan §3 Fixrunde) hat keinen Lauf-Anker: es gibt kein committetes Skript und keine gedruckte Zeile; §3.12 verlangt für eine gemessene Zahl den Lauf. Meine eigenen 60 Mutationen (§4) stützen den Befund bis auf die drei grünen; die Tabelle im Plan trägt „ja“ je Gruppe. | §4; Plan §3 |
| V-7 | INFO | Vorher-Nachher-Sprache in Doku-Prosa: `MR-003` (Adaption) sagt „der Bestand las nach einem Wrapper-Präfix keine Optionen und kannte keine Schlüsselwörter“; der Plan-Abschnitt „Abweichungen“ sagt „der erste Lauf las die Flag-Tokens roh“. Beim `MR` ist die Abgrenzung zur Bestandsregel der Gegenstand (Herkunfts-Anker), beim Plan die Chronik der Arbeit; ohne erwartete Aktion, Kenntnis für den Reviewer. | `git grep -n 'der Bestand las' -- harness/conventions` |
| V-8 | INFO | Die Fixrunde (`git diff --stat a0472421..baff1593`: Guard 154 Zeilen geändert, Maskierer 40 neu, Tabellentest 241 geändert) hat kein zweites Review; der Closure-Trigger („Review-Report ohne offenes HIGH oder MEDIUM“) stützt sich auf den Report am Stand `a0472421` und diese Verifikation. Die Auflösung ist hier finding-weise gemessen (§5); ob ein Nach-Review verlangt wird, entscheidet der Planner. | §5 |

Kein HIGH aus Verifier-Sicht. V-1 ist eine Steuerungs-Frage für Planner/Architect, V-2 bis V-4 sind kleine Nachzüge im
Guard-Kommentar, in `MR-003` und im Tabellentest (je ein Fall, ein Satz), V-5 bis V-8 ohne erwartete Aktion im Slice.

## 8. Verdikt

**DoD trägt:** ja für Liefer-Punkt 1 und 3 sowie die Zeilen `make gates`, Review, Suchlauf, Doku-Update; **Liefer-Punkt 2:
Code und Tabellentest tragen, der Live-Beleg mit `python3` steht aus** (V-5, Entscheidung beim Planner). Die Sensoren sind
meine, nicht die des Implementers: 300 Fälle grün, 155/76 rot an Parent/Review-Stand, `make gates` Exit 0, 57 von 60
Mutationen rot. Aus Verifier-Sicht ist der Slice **closure-fähig unter zwei Bedingungen**: (1) der Architect ratifiziert die
Quote-Lesung (Frage 1) und den Wortlaut der DoD-Zeile 1 (V-1), (2) der Planner entscheidet über V-5 (Live-Beleg annehmen oder
in einer Sitzung ohne Host-Python-Verbot erbringen). V-2 bis V-4 sind vor oder mit der Closure als kleiner Nachzug
empfohlen (ein Tabellenfall für `echo a\|sed -i x`, zwei Nicht-Treffer `perl -8i`/`-9i` und ein Fall mit den Ziffern 3–6,
ein Satz zur balancierten Kehrseite in `MR-003` und im Kopfkommentar).

**Zusammenfassung der Gates:** `make test-command-guard` (300), `make docs-check`, `make gates` (`coverage-gate` 85.10 %,
`generated-sync`, `a-check` 0 Befunde), `make suchlauf-nachmessen` (15 Zeilen), `make doc-immutable RANGE=…`,
`make commit-traceability RANGE=…`, `make doc-commits RANGE=…`, `make kommentar-kennungen DIFF=e98d419c COUNT=1` (0) — alle Exit 0;
der Live-Beleg der Blockmeldung `inplace` in dieser Sitzung; 60 Eingabeseiten-Mutationen (57 rot, 3 grün).

**Übergabe an den Planner:**
1. V-1 und die drei Architect-Fragen (§6) an den Architect: Quote-Lesung ratifizieren (Lexer ja, Grammatik nein, Obergrenze),
   F-3-Erweiterung akzeptieren, einen Satz zum Host-Interpreter in `AGENTS.md` §3.1 „Durchsetzung“.
2. V-5: Live-Beleg `python3` (Klasse `interp`) annehmen oder erbringen lassen; danach die DoD-Zeile 2 abhaken.
3. V-2 bis V-4 als Nachzug an den Implementer oder bei der Closure mit Adresse führen.
4. Träger-Meldungen aus Plan §3 mit Frist „Closure dieses Slice“: `state.md` des Register-Eintrags
   `BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel` („Durchsetzung heute: das Review“ überholt), `welle-transformationen.md`
   §6 (a) („in keinem committeten Text“ überholt), `slice-code-kommentare-bereinigung` §8.
5. Risiko-Ausgänge §6, Closure-Notiz (Finding-Klassen des Reviews und dieses Reports: „Zusage ohne Bindung an ihre
   Eingabeseite“ (Rest), „Kehrseite einer Regel nicht in der Grenz-Zeile“, „Ausschluss des Plans ohne Rückführung umgekehrt“),
   Register, drei Paarungen.

Dieser Report ist ein **Lauf-Beleg** (dieser Stand, dieser Lauf); er ersetzt weder Review noch Closure.
