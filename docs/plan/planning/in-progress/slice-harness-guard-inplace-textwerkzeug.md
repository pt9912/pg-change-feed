# Slice harness-guard-inplace-textwerkzeug: Der PreToolUse-Guard blockt in-place Textwerkzeuge und Host-Interpreter auf Repo-Pfaden — Härtung als `MR-003`, Tabellentest `make test-command-guard`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — Harness-Querschnitt: der Slice trägt keine
Closure-Bedingung, die von seiner DoD verschieden wäre. Er hat keine technische
Kante zu einem Slice der [welle-transformationen](../welle-transformationen.md);
die empfohlene Position steht in §4.

**Bezug:** [`AGENTS.md`](../../../../AGENTS.md) §3.1 (Docker-only, Verbot des
in-place Text-Umschreibens, Absatz „Durchsetzung“), §3.6 (ein Gate braucht eine ADR),
§3.7 (Kommentare) und §3.9;
[`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md) (Herkunft der
Aussagen im Plan; Grenze: kein Sensor über Dateiinhalte); Beobachtungs-Register
`BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel` (`state.md`: Adresse
dieses Slice) und `BEO-PGC/host-werkzeug-jenseits-docker-und-make-ohne-deklaration`;
Baseline-Regelwerk `modul-13-quality-gates.md` §Guard-Härtung (jede Härtung landet
als neuer `MR-<NNN>` mit Grenz-Zeile).

**Berührte Spec-Stellen:** — (Harness-Wächter; keine Spec-Stelle).

**Verantwortlich:** Implementer-Agent.

**Autor:** Planner-Agent, Auftrag des Auftraggebers (Architect-Zug zu
[`AGENTS.md`](../../../../AGENTS.md) §3.1, Commit `17cb4eb3`). **Datum:** 2026-09-26.

---

## 1. Ziel und Abgrenzung

**Ziel:** Der PreToolUse-Guard (`.claude/hooks/pretooluse-command-guard.sh`) blockt
die in-place Formen des Verbots aus [`AGENTS.md`](../../../../AGENTS.md) §3.1 —
`sed` mit `-i`/`--in-place`, `perl` mit `-i`, `awk`/`gawk` mit `-i inplace` — und
Host-`python`/`python3`/`perl` auf Repo-Pfaden; die Härtung landet als `MR-003` mit
ihrer Grenz-Zeile, und `make test-command-guard` bindet jede Zusage an ihre Eingabe
(Treffer, Nicht-Treffer, Bestandsregeln).

**Ausgangslage — Beleg-Anker im Suchlauf-Feld (§3) und in den Beleg-Dateien des Registers.**
Der Guard blockt am Kopf eines Kommando-Segments 15 Paketmanager-Namen (`BLOCKED=` im Skript,
gezählt) und liest ein `tools/harness/blocked/*`-Fragment, wenn es existiert; das Verzeichnis
existiert nicht (`ls tools/harness/blocked` meldet „nicht gefunden“, Stand `89d427e0`). Er ist
am Parent quote-blind (ein Trenner in einem Argument startet ein neues Segment; der Stand `diff` liest
Anführungszeichen, §3 „Fixrunde“) und sieht keine Umleitungen; sein Kopfkommentar nennt beides und „Bewusst NICHT geprüft: andere Interpreter“.
Das Verbot des in-place Text-Umschreibens steht seit dem Architect-Zug `17cb4eb3` in
[`AGENTS.md`](../../../../AGENTS.md) §3.1; die Durchsetzung dort nennt ausdrücklich, dass der
Guard es nicht liest. Der Fehlgriff ist in drei Vorgängen belegt (Register-Eintrag
`BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel`, drei Beleg-Dateien; übernommen aus
ihnen, nicht nachgemessen): `sed -i` in fünf Läufen (Implementer von
`slice-transformationen-antragsweg-usecase` mit Wirkung auf Repo-Dateien — die Spur war eine
`gofmt`-Abweichung —, Implementer, Reviewer und Verifier von
`slice-transformationen-backfill-pfad` und der Implementer von
`slice-backfill-speicher-untersuchung` ohne Wirkung), ein Host-Python-Heredoc auf einer
Repo-Datei in `slice-backfill-speicher-untersuchung`; `perl -pi` und `awk -i inplace` sind
nicht belegt. In `slice-transformationen-backfill-pfad` liefen die Mutationen danach über
Python-Ersetzungen auf einer Kopie im Scratchpad — der zulässige Weg nach §3.1, der weiter
durchgehen muss.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Ein Sensor über Dateiinhalte oder den Diff.** Ein Werkzeugaufruf hinterlässt in der Datei
  keine Signatur ([`AGENTS.md`](../../../../AGENTS.md) §3.1 Begründung; die Grenze
  von [`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md)). Was der Guard nicht
  sieht, bleibt Sache des Reviews (`.harness/skills/reviewer.md` §HIGH „Docker-only-Verstoß“).
- **Blockieren von Umleitungen und flaglosen Schreibwegen** (`> datei`, `>>`, `tee`, `dd of=`,
  `sed … > tmp && mv tmp datei`, `cp`/`mv` über eine Datei). Der Guard liest keine Umleitungen;
  ein Umleitungs-Leser wäre ein Shell-Parser, also der Sandbox-Anspruch,
  den sein Kopfkommentar ausdrücklich nicht erhebt. `> datei` ist zudem der Alltag von
  Logs und Scratchpad-Ausgaben. Die Lücke steht in der Grenz-Zeile von `MR-003`.
- **Die Kopf-Liste `tools/harness/blocked/go`** (`go gofmt python python3 node dotnet java
  gradle uv`). Sie blockte den Kopf eines Segments unbedingt: `docker run … go` ginge
  durch, weil `docker` der Kopf ist, und ein gelegentlicher Scratch-Aufruf von `python3`
  würde ebenso geblockt wie ein Schreibzugriff. Andere Klasse (Toolchain-Aufruf, nicht
  in-place Schreiben), und die drei Vorfälle nennen keinen Host-`go`-Aufruf. Die Mechanik
  liegt im Guard (Fragment-Ladung); die Entscheidung ist eine Nutzer-Frage, §6.
- **Weitere in-place-fähige Werkzeuge** (`ed`, `ex`, `vim -es`, `patch`, `ruby -i`,
  `git apply`, `truncate`). Kein Beleg unter den Vorfällen; nach der Baseline ist der Auslöser
  einer Guard-Härtung die dreimal beobachtete Umgehung, nicht das Bedrohungsmodell
  (Baseline-Regelwerk `modul-13-quality-gates.md` §Guard-Härtung). Jedes weitere Werkzeug ist ein
  neuer `MR` nach Beobachtung.
- **Ein Shell-Parser im Guard: Lexer-Stufe ja, Grammatik nein; Obergrenze und Falsch-Negativ
  in `MR-003`.** Derselbe Grund wie bei den Umleitungen. Der Guard liest
  Anführungszeichen und Backslash über einen Zustandsautomaten (`tools/harness/mask-quotes.awk`,
  Fixrunde, §3); Heredocs, Variablen, Umleitungen und Kommando-Substitution als Wert bleiben
  ungelesen, die Falsch-Positiv-Ränder benennt §6 und bindet der Tabellentest.
- **Die Aufnahme in `make gates`.** Ein Gate braucht eine ADR
  ([`AGENTS.md`](../../../../AGENTS.md) §3.6, §4), und ein Wächter gehört nicht in die
  Gate-Tabelle (Baseline-Regelwerk `modul-13-quality-gates.md` §Guard-Härtung: er verhindert
  eine Handlung, er prüft kein Ergebnis). `make test-command-guard` ist ein Werkzeug ohne
  Gate-Bindung, wie `make test-fmt-check`.
- **Eine Änderung der Regel in [`AGENTS.md`](../../../../AGENTS.md) §3.1 selbst.** Sie ist seit
  `17cb4eb3` beschlossen; dieser Slice zieht nur den Absatz „Durchsetzung“ auf den neuen Stand
  (Liefer-Punkt 3).

## 2. Definition of Done

- [x] **Liefer-Punkt 1 — die in-place Formen.** Der Guard blockt (Ausgabe
      `"decision": "block"`, Exit 0, wie im Bestand) ein Kommando-Segment, dessen Kopf —
      nach dem Überspringen der Zuweisungs- und Wrapper-Präfixe des Bestands und nach
      `-exec`/`-execdir`/`-ok` — ist: `sed` mit einem Flag `--in-place`,
      `--in-place=<Suffix>` (auch jede eindeutige Abkürzung ab `--i`) oder einem Bündel der
      Form `-[nEsrzub]*i` (`-i`, `-i.bak`, `-ni`, `-Ei`); `perl` mit einem Bündel der Form
      `-[0-7lanpsw]*i` (`-i`, `-pi`, `-i.bak`, `-0777pi`; `-MList::Util` und `-e` sind keines;
      die Optionen enden am Skriptnamen); `awk`/`gawk` mit `-i` und
      dem Wert `inplace` (auch `-iinplace`, `--include=inplace`, `--include inplace`).
      Die Segmentierung ist quote-bewusst (`tools/harness/mask-quotes.awk`): ein Trenner in
      Anführungszeichen startet kein Segment, ein Anführungszeichen-Argument ist ein Token,
      damit ein Muster mit `|` in Anführungszeichen (`'sed -i|perl -pi'`) nicht als
      Kommando `perl -pi` gilt. Der Block gilt unabhängig vom Ziel: ein `sed -i` auf einer
      Kopie im Scratchpad wird ebenso geblockt (der Mutationsweg ist `sed … Datei > Kopie`
      oder Edit/Write; offene Frage 2 in §6). Die Blockmeldung dieser Klasse nennt den
      Ersatzweg (Edit/Write; `sed` ohne `-i` nach stdout). Die Bestandsregeln (15
      Paketmanager-Namen, Sub-Shell-Rekursion, fail-closed bei Parse-Zweifel) gelten weiter.
      Die Kopf-Erkennung (Wrapper-Optionen, Schlüsselwörter) gilt für alle Klassen
      (§3 „Fixrunde“, F-3). *Zu belegen durch:* `make test-command-guard` (Tabellentest, netzlos,
      gegen ein Wegwerf-Repo im Temp-Verzeichnis) mit Treffern je Form und je Position
      (Kopf, nach `&&`, hinter `xargs`, `sudo`, `bash -c`, absoluter Pfad des Werkzeugs,
      `find … -exec`), **Nicht-Treffern** je Nachbarform (`sed -n`, `sed -E`, `sed s/a/b/ f`,
      `sed --version`, `awk '{print $1}'`, `awk -F, …`, `gawk -i /tmp/lib.awk …`,
      `perl -MList::Util -e …`, `perl -ne … /tmp/x`, `grep -i`, `cp -i`, `git grep -E
      'sed -i|perl -pi'`, `git commit -F <Datei>`, `make gates`) und den Bestandsregeln (`pip
      install`, `apt`, `bash -c "pip …"`, Tiefe über 3, defektes JSON); je Zusage an ihre
      Eingabe gebunden durch eine Mutation am Guard, an einer Kopie im Wegwerf-Repo gelaufen
      (`GUARD=<Kopie> make test-command-guard`, Exit ≠ 0, gesehenes Rot im Bericht): Regel
      entfernt · nur das nackte `-i` (Bündel, `-i.bak`, `--in-place` fallen durch) · jedes
      `sed` blockt (`sed -n` färbt rot) · Trenner in Anführungszeichen nicht
      maskiert (das Muster mit `|` färbt rot) · `-exec`-Kopf nicht gelesen · perl-Bündel mit
      beliebigem `i` (`-MList::Util` färbt rot) · awk-`-i` ohne `inplace` blockt (`gawk -i
      /tmp/lib.awk` färbt rot) · Block-Ausgabe entfernt; Menge der Erprobung: die Fälle des
      Tabellentests. Dazu ein **Live-Beleg in der Sitzung des Implementers**: ein Aufruf
      `sed -i s/a/b/ /nicht-vorhanden-scratch` (das Ziel existiert nicht, eine Wirkung ist
      ausgeschlossen) wird geblockt — der Wortlaut der Blockmeldung steht im Bericht —, ein
      `sed -n 1p Makefile` läuft.
- [x] **Liefer-Punkt 2 — Host-Interpreter auf Repo-Pfaden.** Ein Segment mit Kopf `python`,
      `python3`, `python3.<N>` oder `perl` blockt, wenn der **ganze Befehlsstring** ein
      Repo-Pfad-Muster trägt: den Absolutpfad der Repo-Wurzel `R/` (`R` aus dem Ort des
      Guards, zwei Ebenen über `.claude/hooks/`), oder den Namen eines Eintrags der obersten
      Repo-Ebene (`ls -A "$R"` ohne `.git`; Verzeichnis-Name mit folgendem `/`, Datei-Name
      als ganzes Wort), vor dem kein Pfad-Zeichen `[A-Za-z0-9_./~-]` steht (ein
      vorangestelltes `./` ist erlaubt). Treffer: `python3 tools/x.py`, `python3 ./tools/x.py`,
      `python3 - <<'EOF'` mit `open('docs/a.md', 'w')` im Text, `python3 -c "open('Makefile',
      'w')"`, `perl -e 'unlink q(Makefile)'`, `python3 <R>/docs/a.md`. Nicht-Treffer:
      `python3 --version`, `python3 -c 'print(1)'`, `python3 /tmp/scratch/mutate.py
      /tmp/scratch/copy.go`, `python3 /tmp/x/docs/a.py` (vor `docs/` steht `/`), `perl -ne
      'print' /tmp/x`. Die Blockmeldung nennt Edit/Write, ein Repo-Werkzeug hinter `make`
      und `sed s/a/b/ Datei > Scratchpad-Kopie`. *Zu belegen durch:* dieselben Tabellentest-Fälle, je an ihre Eingabe
      gebunden durch eine Mutation: Regel entfernt · nur das Segment statt des ganzen
      Befehls gelesen (der Heredoc-Fall färbt rot) · die Pfad-Zeichen-Bedingung entfernt
      (`/tmp/x/docs/a.py` färbt rot) · `./` nicht erlaubt (`python3 ./tools/x.py` färbt rot) ·
      Datei-Namen der obersten Ebene nicht gelesen (`Makefile`-Fall färbt rot). Beleg:
      Tabellentest und Mutationen an der Hook-Schnittstelle; der Ausgabeweg (`emit_block`) ist für
      alle drei Klassen derselbe, die Klasse `inplace` ist live belegt (Liefer-Punkt 1), die
      Klasse `interp` sah der Reviewer zweimal live geblockt (übernommen aus dem Review-Report,
      nicht nachgemessen). Ein Live-Aufruf mit Host-`python` entfällt (`AGENTS.md` §3.1).
- [x] **Liefer-Punkt 3 — die Träger.** (a) `harness/conventions/MR-003-…md` per `cp` aus
      `.harness/baseline/v6.9.0/templates/harness/conventions/MR-NNN-titel.template.md`, in
      place gefüllt (Auslöser: der Register-Eintrag mit seinen drei Beleg-Dateien; Adaption:
      was der Guard liest; **Grenz-Zeile**: was er nicht kann — die Punkte aus §6, Risiko 2;
      Ersetzt-Baseline-Regel: `grundlagen-durchsetzungsschicht.md` §Grenzen — ehrlich
      benannt; Auflösungs-Trigger), plus die Zeile in `harness/conventions.md`
      §Aktive Adaptionen (Anker `mr-003`); (b) der Kopfkommentar des Guards im Indikativ
      (Zusage · Abgrenzung · Grenze, [`AGENTS.md`](../../../../AGENTS.md) §3.7; der Satz
      „Bewusst NICHT geprüft: andere Interpreter“ nennt den neuen Stand), keine
      Slice-/Wellen-Nummer im Kommentar; (c) [`AGENTS.md`](../../../../AGENTS.md) §3.1 Absatz
      „Durchsetzung“ nennt, was der Guard jetzt blockt und was nicht (Herkunft `· seit
      slice-harness-guard-inplace-textwerkzeug`), und der Satz zur Mutationsprobe nennt den
      Weg auf der Kopie (`sed` ohne `-i` nach stdout, `-i` auch dort geblockt); (d) `make
      test-command-guard` in `Makefile` und eine Zeile unter Werkzeuge in
      [`harness/README.md`](../../../../harness/README.md) §Sensors („kein Gate“, Host-Werkzeuge
      `bash`, `awk`, `mktemp`). *Zu belegen durch:* Lesen der Stellen, `make docs-check`
      (Links, Anker) und der Suchlauf in §3.
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und gesondert
      ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff);
      `make suchlauf-nachmessen PLAN=<Plan-Datei>` läuft nach jeder Fixrunde
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [x] Doku-Update: [`harness/README.md`](../../../../harness/README.md) §Sensors
      (Liefer-Punkt 3d); das Benutzerhandbuch bleibt unberührt (keine
      Betreiber-Oberfläche).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel · neuer
      Sensor · benannte Spec-Lücke).
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — der Ausgang von
      `BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel` (`state.md`: die
      Durchsetzung heute, die Adresse aufgelöst, die zwei Nutzer-Fragen aus §6 mit
      Adresse); keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7
      notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von
      der Slice-Closure selbst (der Slice hat keine Welle; das Ereignis kann
      eintreten).

**Umfang:** S bis M — Schätzung, nicht gemessen: ein Guard-Skript um zwei Erkennungen,
ein Tabellentest mit rund vierzig Fällen, ein Makefile-Ziel, ein `MR`, drei kurze
Träger-Änderungen. Der zweite Liefer-Punkt ist der abtrennbare Teil (§4).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.claude/hooks/pretooluse-command-guard.sh` | update | Liefer-Punkt 1 und 2: Erkennung der in-place Formen (Flag-Tokens hinter der Anführungszeichen-Maskierung, Kopf nach `-exec`/`-execdir`/`-ok`) und des Repo-Pfad-Musters auf `python`/`python3`/`perl`; eine Blockmeldung je Klasse (gültiges JSON); Kopfkommentar im Indikativ ([`AGENTS.md`](../../../../AGENTS.md) §3.7). Die Bestandsregeln gelten weiter, der Guard ist fail-closed. |
| `tools/harness/run-command-guard-tests.sh` | neu | Tabellentest: baut ein Wegwerf-Repo im Temp-Verzeichnis (`.claude/hooks/` mit einer Kopie des Guards, `tools/harness/extract-command.awk`, oberste Ebene `docs/`, `internal/`, `Makefile`), füttert das Hook-JSON auf stdin und prüft Ausgabe und Exit; der Prüfling ist per `GUARD` übersteuerbar (Mutationsläufe an Kopien); Vorbild `tools/harness/run-fmt-check-tests.sh` und `tools/harness/run-suchlauf-nachmessen-tests.sh`. Schreibt nur ins Temp-Verzeichnis. |
| `tools/harness/mask-quotes.awk` | neu (Fixrunde) | Segmentierungs-Vorstufe des Guards: maskiert Trenner, Leerraum und Zeilenumbruch in Anführungszeichen und hinter Backslash (Finding F-2); 40 Zeilen, davon 12 Kommentar. |
| `Makefile` | update | Ziel `test-command-guard` mit `.PHONY` und Hilfe-Zeile (`GUARD`/`MASKER` übersteuerbar); kein Eintrag in den Gate-Zielen. |
| `harness/conventions/MR-003-…md` | neu, per `cp` | Härtung nach der Baseline-Regel (Auslöser, Adaption, Grenz-Zeile, Auflösungs-Trigger); der Vertrag des Guards steht dort und im Kopfkommentar, nicht in `harness/sensors/` — der Guard ist kein Target, und ein Wächter gehört nicht in die Gate-Tabelle. |
| `harness/conventions.md` | update | Zeile `MR-003` in §Aktive Adaptionen. |
| `harness/README.md` | update | Zeile unter Werkzeuge für `make test-command-guard` (kein Gate). |
| `AGENTS.md` §3.1 | update | Absatz „Durchsetzung“ auf den neuen Stand; Satz zur Mutationsprobe. |

**Der Aufbau des Tabellentests** (Ansatz, keine Vorschrift): je Fall ein Kommandostring,
der als `{"tool_input":{"command":"<escaped>"}}` auf die stdin des Guards geht;
erwartet ist entweder `"decision": "block"` in der Ausgabe oder keine Ausgabe (Pass-Fall,
[`AGENTS.md`](../../../../AGENTS.md) §3.1 nennt den Guard einen Stolperdraht, der im
Pass-Fall schweigt). Ein Nicht-Treffer-Fall steht neben jedem Treffer-Fall derselben Form,
damit eine Mutation „blockt zu viel“ rot wird. Falsch-Positiv-Ränder, die der Guard behält
(ein Trenner in einem Argument, eine Heredoc-Zeile mit `sed -i` am Anfang), stehen als
benannte Fälle mit dem erwarteten Block, damit die Grenz-Zeile von `MR-003` eine Messung
trägt und keine Behauptung.

**§3.13-Suchlauf (committetes Feld).** Bewegte Eigenschaft: „was der PreToolUse-Guard blockt
und was er nicht liest“ — Symbolname `pretooluse-command-guard`, seine Beschreibung
(„blockt Paketmanager“, „Stolperdraht“, „scannt den Command-String“) und das Zählwort der
in-place Formen (`sed -i`, `perl -pi`, `awk -i`). Suchraum: der ganze Baum; ausgenommen sind
`docs/reviews/**`, die Records unter `done/`, `.harness/baseline/**` und das
Beobachtungs-Register (es wird in der Closure geschrieben und hat unten eine eigene Zeile).
Stand ist der Parent `89d427e0`; der Implementer ergänzt die Zeilen mit Stand `diff`:

```suchlauf
89d427e0 3 -E 'pretooluse-command-guard' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline' ':!docs/plan/planning/observations'
89d427e0 6 -E 'blockt (Host-)?Paketmanager|Stolperdraht|scannt den Command-String' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline' ':!docs/plan/planning/observations'
89d427e0 9 -E 'sed -i|perl -pi|awk -i' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline' ':!docs/plan/planning/observations'
89d427e0 2 -E 'Durchsetzung heute|noch kein Plan angelegt' -- docs/plan/planning/observations
89d427e0 1 -E 'in keinem committeten Text' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline' ':!docs/plan/planning/observations'
89d427e0 3 -E 'inplace-textwerkzeug-am-repo-trotz-nutzerregel' -- docs/plan/planning/open docs/plan/planning/next docs/plan/planning/welle-transformationen.md
89d427e0 0 -E '\.claude/hooks|pretooluse|harness/conventions' -- docs/plan/planning/open docs/plan/planning/next docs/plan/planning/welle-transformationen.md
diff 12 -E 'pretooluse-command-guard' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline' ':!docs/plan/planning/observations'
diff 8 -E 'blockt (Host-)?Paketmanager|Stolperdraht|scannt den Command-String' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline' ':!docs/plan/planning/observations'
diff 148 -E 'sed -i|perl -pi|awk -i' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline' ':!docs/plan/planning/observations'
diff 0 -E 'Durchsetzung heute|noch kein Plan angelegt' -- docs/plan/planning/observations
diff 0 -E 'in keinem committeten Text' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline' ':!docs/plan/planning/observations'
89d427e0 3 -E 'in-place Textwerkzeuge liest er nicht|Bewusst NICHT gepr|blockt Host-Paketmanager und Host-Toolchains' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline' ':!docs/plan/planning/observations'
diff 0 -E 'in-place Textwerkzeuge liest er nicht|Bewusst NICHT gepr|blockt Host-Paketmanager und Host-Toolchains' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline' ':!docs/plan/planning/observations'
diff 9 -E 'test-command-guard' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline' ':!docs/plan/planning/observations'
```

Die ersten sieben Zeilen sind am Parent gemessen (Bezugs-Befehle im Block, Stand `89d427e0`; Ergebnis:
Symbol 3 Fundstellen — der Guard selbst, `.claude/settings.json`, [`AGENTS.md`](../../../../AGENTS.md)
§3.1; Beschreibung 6; Zählwort 9, darunter ein `sed -i` in `sdks/python/Dockerfile`, das im
Docker-Bau läuft und den Guard nicht berührt — er scannt nur die Bash-Aufrufe des Laufs, nicht
die Rezepte hinter `make`).

**Ergebnis am Stand `diff` (Arbeitsbaum bei der Closure, nach den Träger-Nachzügen; Gefundenes und Nichtgefundenes):**
Symbol 12 Fundstellen (Guard, `.claude/settings.json`, [`AGENTS.md`](../../../../AGENTS.md), `Makefile`,
`harness/README.md`, `harness/conventions.md`, `MR-003`, `mask-quotes.awk`, vier im Tabellentest) — die
neuen Träger tragen den Namen, keiner der 3 Bestandsträger sagt etwas Falsches. Beschreibung 8 (Guard 2,
[`AGENTS.md`](../../../../AGENTS.md) 3, `tools/harness/extract-command.awk` 1, `MR-003` 2; `implement-slice.md` trägt das
Stichwort nicht mehr): dessen Satz „der Guard scannt den Command-String“ ist
durch die Anführungszeichen-Lesung überholt und nachgezogen (Träger-Tabelle); der Kommentar in
`tools/harness/extract-command.awk` („Stolperdraht, keine Sandbox“) bleibt wahr und ist unberührt.
Zählwort 148: die Zahl ist gewachsen um die neuen Träger (`MR-003`, Tabellentest mit seinen Fällen,
Guard-Kommentar, [`AGENTS.md`](../../../../AGENTS.md)); die Bestandstreffer `reviewer.md` (Prüfgegenstand des Reviews),
`implement-slice.md` und `sdks/python/Dockerfile` sind unverändert wahr, die Treffer in
`slice-code-kommentare-bereinigung` §6/§8 und `welle-transformationen` §6 sind in der Closure nachgezogen
(Tabelle unten). Die Zeilen 4 und 5 stehen dadurch bei 0: „Durchsetzung heute“ im Register-`state.md` und
„in keinem committeten Text“ in `welle-transformationen` §6 sind getilgt; `test-command-guard` trägt
neun Fundstellen (die Nennung in `welle-transformationen` §6 ist die neunte).
Der überholte Wortlaut („liest er nicht“, „Bewusst NICHT geprüft: andere Interpreter“, „blockt … Host-Toolchains“:
3 am Parent, 0 am Stand `diff`) ist aus den zwei Trägern getilgt. **Nicht gefunden:** kein weiterer
Träger unter `docs/user/`, `spec/`, `.github/`, `harness/sensors/` oder `harness/targets/`, der den Guard
beschreibt. **Nicht durch Suchlauf fassbar:** ein Zeilen-Lokator auf den Guard — der Suchlauf trifft
Symbolnamen, keine Zahlen (Grenze in [`AGENTS.md`](../../../../AGENTS.md) §3.13); der Reviewer liest sie
gegen den Diff.

**Abweichungen vom Plan und Ergänzungen (Plan-Nachzug im Lauf; Anker: der Diff dieses Laufs).**

- **Beleg-Dateien.** Das Register trägt vier (Plan: drei): die vierte
  (`slice-antragsqueue-lesefehler-failed`) nennt zwei Host-`python3`-Fälle (Implementer ohne Wirkung,
  Verifier mit unklarem Ziel). `MR-003` zählt vier; der Host-Interpreter ist damit in zwei Vorgängen
  belegt (Risiko 5 nennt einen), `perl -pi` und `awk -i inplace` bleiben unbelegt.
- **Name des Eintrags:** `harness/conventions/MR-003-guard-inplace-textwerkzeug.md`.
- **Erkennung über den Plantext hinaus:** (a) die Optionen hinter einem Wrapper-Präfix werden übersprungen
  (`xargs -r sed -i`, `xargs -I{} sed -i`; der Plan nennt „hinter `xargs`“, die Alltagsform trägt
  Optionen); der Wert einer xargs-Option (`xargs -n 1`) wird übersprungen, ein Optionswert anderer Wrapper
  (`sudo -u x`) bleibt Grenze; (b) ein Trenner in Anführungszeichen startet kein Segment und ein
  Anführungszeichen-Argument ist ein Token, ein Flag mit Anführungszeichen (`sed -i''`, `sed '-i'`)
  bleibt ein Flag (Fixrunde, §3); (c) das find-Ende (`+`, `\`, `;`) begrenzt die Flag-Suche (`-iname` hinter
  `-exec … +` ist kein sed-Flag); (d) `-okdir` neben `-ok`; (e) der Tabellentest prüft je Block die
  Begründung der Klasse (`pkg`, `inplace`, `interp`) und die JSON-Form.
- **Tabellentest:** 317 Fälle statt der geschätzten rund vierzig (311 Tabellenzeilen und 6 Sonderfälle,
  gedruckt „alle 317 Fälle bestanden“ von `make test-command-guard`); die Gruppen sind Bestand, Kopf-Erkennung hinter Wrapper-Optionen und
  Schlüsselwörtern, in-place je Form, je Zeichenklassen-Mitglied und Position, Nicht-Treffer,
  Anführungszeichen-Lesung, Host-Interpreter, benannte Falsch-Positiv-Ränder, benannte Grenzen.
- **Liefer-Punkt 2 geliefert.** Die Falsch-Positiv-Fläche ist gebunden: Nicht-Treffer neben jedem Treffer
  (`/tmp/x/docs/a.py`, `python3 --version`, Scratchpad-Mutationsweg), benannte Ränder mit erwartetem
  Block (`cd /tmp/x && python3 tools/x.py`, Text mit Repo-Namen), Mutationen je Zusage rot gesehen (Bericht).
  Ein Live-Aufruf mit Host-`python3` entfällt: er verstieße gegen `AGENTS.md` §3.1; der Tabellentest trägt
  die Hook-Schnittstelle (Risiko 4).
- **Nutzer-Fragen aus §6 entschieden** (Auftrag des Auftraggebers an den Implementer, 2026-09-26):
  Frage 1 — `tools/harness/blocked/go` nicht in diesem Slice, zuerst die Wirkung von Liefer-Punkt 1
  abwarten; Frage 2 — keine Scratchpad-Ausnahme, der Guard blockt `sed -i` unbedingt (Mutationsproben
  laufen über Stdout-Umleitung auf einer Scratchpad-Kopie).

**Fixrunde (Anker: Review-Report `review-slice-harness-guard-inplace-textwerkzeug`, Finding je Zeile;
Stand `diff`).** Der Tabellentest zählt 317 Fälle (`make test-command-guard`: „alle 317 Fälle bestanden“,
311 Tabellenzeilen, 6 Sonderfälle; die Verifikation maß den Stand der Fixrunde mit 300 Fällen); gegen den
Guard am Parent `e98d419c` sind 168 rot, gegen den Guard am Review-Stand `a0472421` 81 (`GUARD=<Datei>`,
je ein Lauf bei dieser Closure am Stand mit 317 Fällen, gezählt mit `grep -c '^FEHLER:'` über das Log;
die Zahlen der Fixrunde waren 155 und 76, der Implementer-Nachzug nannte 80). Die Gruppe „Bestandsregeln“ (18 Tabellenzeilen) ist
am Parent grün.

| Finding | Entscheidung und Änderung | Beleg-Anker |
|---|---|---|
| F-2 | Im Guard gelöst, keine Rückführung nach `open/`: `tools/harness/mask-quotes.awk` ist ein Zustandsautomat über `'`, `"`, Backslash sowie `$(`/Backtick in `"…"` (kein Shell-Parser: Heredocs, Umleitungen, Variablen bleiben ungelesen). Ein Trenner, Leerraum oder Zeilenumbruch in Anführungszeichen wird vor der Segmentierung zu einem Steuerzeichen; ein Anführungszeichen-Argument ist ein Token. Ein unbalanciertes Anführungszeichen gibt den Rohstring zurück (Segmentierung ohne Anführungszeichen-Kenntnis: mehr Segmente, nie weniger; ein Apostroph im Heredoc-Text blockt also nur, wenn eine Kommando-Zeile mit `sed -i` folgt). `git commit -m "a & sed -i b"` und `grep -E 'a\|sed -i ' f` sind Nicht-Treffer, `\;` bleibt ein Trenner. Fällt der Maskierer aus, blockt der Guard (fail-closed, Fall „Maskierer nicht lesbar“). Latenz 15 ms je Aufruf gegen 11 ms (20 Aufrufe je Guard, gemessen). | 37 Fälle „Anführungszeichen“ (die Gruppe zwischen den Kommentarzeilen `# --- Anführungszeichen` und `# --- Liefer-Punkt 2` in `tools/harness/run-command-guard-tests.sh`, gezählt bei dieser Closure mit `grep -E '^(block\|pass) '` über ihre Zeilen; die Suchlauf-Form mit drei Gliedern, `grep -E "sed -i\|perl -pi" f`, `echo 'a\|sed -i x'`, Backslash, Tab, Trenner ohne Leerraum, `$(…)`/Backtick-Kontext, die Heredoc-Kehrseite); live: `git grep -c -E 'sed -i\|perl -pi\|awk -i' -- <Datei>` läuft in der Sitzung des Implementers ohne Block, `sed -i s/a/b/ /nicht-vorhanden-scratch` blockt mit der Meldung der Klasse `inplace`. Die Quote-Lesung ist als Lexer-Stufe ratifiziert, ihre Obergrenze steht in `MR-003` (Architect-Zug zu diesem Slice). |
| F-1 | Je Mitglied der Zeichenklassen ein Treffer-Fall und je Klassengrenze ein Nicht-Treffer: sed `[nEsrzub]` (Ergänzung `b`: `sed -bi`; Nicht-Treffer `-ei`, `-fi`, `-xi`), perl `[0-7lanpsw]` (Ziffernklasse von `0-9` auf `0-7` verengt: die Argumente von `-0` und `-l` sind oktal; Nicht-Treffer `-ei`, `-Minteger`, `-Ii`, `-Mfeature`), Pfadzeichen `[A-Za-z0-9_./~-]` (acht Mitglieder davor und dahinter, dazu die Zeichen `=` `:` `,` `"` `)` als Nicht-Mitglieder). Die Behauptung „jede Zusage gebunden“ in `MR-003` nennt die Menge (Tabellentest, Mutationen je Klassenmitglied). | 38 Fälle „Zeichenklassen“ (die perl-Ziffern 1 bis 7 je einzeln, 8 und 9 als Nicht-Treffer), 23 Fälle „Pfadzeichen davor/dahinter“ in der Gruppe Host-Interpreter (gezählt bei dieser Closure mit `grep -cE '"(kein )?Pfadzeichen (davor\|dahinter)'` über den Tabellentest); Mutationen unten. |
| F-3 | Das Überspringen der Wrapper-Optionen gilt für **alle** Klassen: die Lücke `env -i pip`, `xargs -n1 pip`, `time -p pip` ist für die Paketmanager-Klasse dieselbe wie bei `sed -i`, und eine Kopf-Erkennung für alle Klassen ist eine Wartungsstelle. Die Kopf-Erkennung gilt für alle Klassen und steht so in `MR-003` (Adaption) und im Tabellentest als eigene Gruppe „Kopf-Erkennung hinter Wrapper-Optionen“; `ls \| xargs -n1 pip` und `env -i pip install x` stehen dort (am Parent rot gemessen). `command -v`/`-V` und `-pv` zeigen an und führen nichts aus: sie blocken nicht; `type`, `which` sind ohnehin nicht Kopf einer Ausführung; `command -p pip` blockt. | Gruppe „Kopf-Erkennung“ (13 Fälle); `command -v`-Fälle am Guard `a0472421` rot (3), am neuen Guard grün. |
| F-4 | Schlüsselwörter `do then else elif if while until !`, führendes `{` und `)` (Funktionsdefinition) werden vor der Kopf-Erkennung übersprungen; ein `case` überspringt bis zum ersten Label, Labels gelten nur nach einem `case` im selben Befehl (`echo $(date) sed -i x` bleibt Nicht-Treffer). Die Form `xargs -n 1 sed -i` (Wert einer xargs-Option) ist mit gelöst. Grenz-Zeile in `MR-003`: Optionswerte anderer Wrapper (`sudo -u x`, `env -u X`, `nice -n 10`), Werkzeug aus einer Variablen, `eval "$cmd"`, Aliase. | Fälle „Schlüsselwort“ (`for`, `while`, `if`, `else`, `elif`, `until`, `!`, drei `case`-Formen, Funktion), am Parent rot. |
| F-5 | Die Meldung der Klasse `interp` nennt Edit/Write, ein Repo-Werkzeug hinter `make` und `sed s/a/b/ file > /path/to/scratch-copy` — die drei Wege aus `AGENTS.md` §3.1, keinen Host-Interpreter; `MR-003` (Begründung) sagt dasselbe und nennt den Interpreter auf einem Pfad ohne Repo-Namen eine Grenze, keinen zulässigen Weg. | Fall „Meldung interp: Weg nach AGENTS.md 3.1“ (prüft Edit/Write, den stdout-Weg und die Abwesenheit von `mutate`). |
| F-6 | Gelöst: `sed -i''`, `\sed`, `busybox`, `/usr/bin/env`, `/usr/bin/sudo`, `gsed`, die Abkürzung `--in-p`/`--i`, `find … -exec sh -c '…'`/`-exec python3 …` (der Rest hinter `-exec` ist ein eigenes Kommando), `eval "…"`; das Bündel `-input.txt` bei sed bleibt ein Treffer, weil GNU-sed es als `-i` mit Suffix `nput.txt` liest (real gemessen: `sed -n 1p -input.txt f` legt die Sicherungsdatei `fnput.txt` an), bei perl enden die Optionen am Skriptnamen (`perl x.pl -input a` blockt nicht). Rest als Grenze in `MR-003`: `-Wpi`, die awk-Abkürzung `--inc=inplace`, andere Interpreter (`node -e`, `ruby -e`, `uv run python`), Optionswerte anderer Wrapper. `MR-003` nennt die Grenz-Zeile als Kurzform in Kopfkommentar und `AGENTS.md`; „gleichlautend“ ist berichtigt (die `AGENTS.md`-Liste ist eine Teilmenge). | Fälle „Position“ (80) und „Grenzen“ (16). |
| F-7 | Der Konjunktiv-Kandidatenlauf in `implement-slice.md` Schritt 20 liest `'*.go' '*.sh' '*.awk'` mit `(//\|#)` und den transliterierten Formen (`waere`, `wuerde`, `haette`); die Urteilsregel nennt `else`-Zweige in Shell-Kommentaren und `${#var}` als zulässige Treffer. Der Lauf auf dem eigenen Diff (`git diff -U0 HEAD`, 80 hinzugefügte Kommentarzeilen) liefert 0 Treffer, auf dem Diff seit `e98d419c` 0 Treffer; die Bestandszeile „`waere` sonst“ im Guard ist entfernt, dazu das „wuerde“ im Kopfkommentar. | Befehl in `implement-slice.md`, Lauf oben. |
| F-8 | Kenntnis: `perl -i` und `awk -i inplace` haben keinen Beleg im Register (Risiko 5, `MR-003` Begründung); Liefer-Punkt 2 ist der Interpreter-Denylist näher als einer Zerlegung (`MR-003` Begründung: beide Regeln blocken nur auf einem Repo-Pfad); das Feld „Ersetzt-Baseline-Regel“ nennt den Satz, der für alles Ungelesene weiter gilt (im Feld ausgeführt). Keine Änderung der Baseline-Zuordnung. | `MR-003`. |
| F-9 | Keine Aktion (äquivalente Mutation durch Doppelschutz; Fixrunde: die Doppelmutation beider Zeilen ist rot). | Mutationen unten. |

**Mutationen der Fixrunde (Zusage · mutierte Eingabe · gesehenes Rot; je `GUARD=<Kopie> MASKER=<Kopie>` an einer
Kopie von Guard bzw. Maskierer).** Die Zeilen bis „Bestand und fail-closed“ und die Zeile „äquivalent“ (120
Mutationen, 118 rot, 2 äquivalent; Menge der Erprobung: die damals 300 Fälle) sind **übernommen** (Implementer-Lauf der
Fixrunde, nicht nachgemessen; kein committetes Skript, keine gedruckte Zeile); die drei Zeilen „Nachzug“ sind
**gemessen** (Implementer-Zug nach der Verifikation, Menge der Erprobung: die 317 Fälle, je ein Lauf von
`tools/harness/run-command-guard-tests.sh` mit `GUARD=<Kopie> MASKER=<Kopie>`, Exit 1, gedruckt „FEHLER“-Zeilen):

| Zusage | Mutierte Eingabe | Rot |
|---|---|---|
| sed-Klasse | je ein Mitglied aus `[nEsrzub]` (7 Läufe) | ja, 2 bis 3 Fälle je Mitglied |
| perl-Klasse | `0` und die Obergrenze `7` aus `0-7`, je `l a n p s w` (8), Präfix-Regel `[A-Za-z0-9]` | ja |
| perl `-e`/`-I`/Skriptname | Bündel `-e` nie, Bündel weit (`-Mfeature`), `-I` ohne Wert, Optionsende entfernt | ja (je 1 bis 2 Fälle) |
| Pfadzeichen-Klasse | je `A-Z`, `a-z`, `0-9`, `_`, `.`, `/`, `~`, `-` (8) | ja, 2 bis 5 Fälle |
| sed-Regeln | Regel entfernt / nur `-i` / jedes sed blockt / `--in-place`-Abkürzung ab 2 und nur voll / Anführungszeichen-Bereinigung entfernt / `gsed` (2) | ja |
| awk-Regeln | `-i` ohne `inplace` / `-iinplace` / `--include=inplace` entfernt | ja |
| `-exec` | Rekursion entfernt / `-execdir` / `-ok` / `-okdir` einzeln nicht gelesen / find-Ende nicht gelesen | ja |
| Wrapper | Optionen nicht übersprungen / xargs-Wert nie / Wert jedes Wrappers / je `-n -P -L -I -s -d -a -E` (8) / `sudo`, `env`, `busybox` kein Präfix / `eval` nicht rekursiv / kein Basename am Präfix und am Kopf / `\sed` | ja |
| `command -v` | Ausnahme entfernt / jede Option / `-V` fehlt / Bündel `-pv` fehlt | ja |
| Schlüsselwörter | je `do then else elif if while until !` (8), `{`, `)` | ja |
| `case` | Label ohne Zustand / nie übersprungen / Zustand nie gesetzt / `esac` setzt nicht zurück / kein Sprung zum ersten Label | ja |
| Maskierer | einfache/doppelte Anführungszeichen nicht maskiert oder öffnend, `$(`/Backtick in `"…"` maskiert, Backslash, `\;`, Zeilenfortsetzung, unbalanciert, Leerraum, Tab, Zeilenumbruch, `\|`, `&`, `;`, `(`, Backtick je einzeln, `$(`- und Backtick-Kontext endet nicht | ja (20) |
| Host-Interpreter | Regel entfernt / `perl` / `python3.N` / nur `python3` / Repo-Wurzel / `./` / Datei-Namen / Nachbedingung / Grund `pkg` / Meldung mit `mutate.py` | ja |
| Bestand und fail-closed | Tiefe 4 erlaubt / defektes JSON / Maskierer-Fehler nicht fail-closed / fehlendes awk (Präsenzprüfung und Exit-Code zusammen) | ja |
| Nachzug: Backslash-Escape des Maskierers | (a) `mk(d)` am Escape entfällt (1 rot: `echo a\|sed -i x`); (b) das Zeichen hinter `\` wird nicht verbraucht (7 rot); (c) Backslash außerhalb von `"…"` nicht gelesen (3 rot); (d) Backslash in `"…"` nicht gelesen (1 rot); (e) Backslash in `'…'` gelesen (1 rot); (f) das Paar `\\` nicht gepaart (2 rot: `echo a\\|sed -i x`, `echo "a\\"; …`) | ja, sechs von sechs |
| Nachzug: perl-Ziffernklasse `0-7` | `[0-9…]` (2 rot: `-8i`, `-9i`), je ein fehlendes Mitglied 0 bis 7 (acht Läufe, je 1 bis 4 rot: der Einzelfall `-3pi` bis `-6pi` und das Bündel `-l0127pi`), `[0-27…]` (4 rot: Ziffern 3 bis 6) | ja, zehn von zehn |
| Nachzug: Kehrseite der Quote-Lesung (Heredoc) | Maskierer maskiert nichts (31 rot, darunter der Grenzfall „zwei Apostrophe in zwei Heredoc-Zeilen“); Zeilenumbruch auch außerhalb von Anführungszeichen maskiert (4 rot, darunter „zwei Apostrophe in einer Heredoc-Zeile“); Rohstring-Rückfall entfernt (3 rot: die Fälle „unbalanciertes Anführungszeichen“) | ja, drei von drei |
| äquivalent | fehlendes awk (nur die Präsenzprüfung: der Exit-Code der `awk`-Extraktion blockt dieselbe Eingabe); Sprung hinter das erste `case`-Label (das Label-Überspringen leistet dasselbe) | grün, Doppelschutz |

| Träger | Befund | Behandlung |
|---|---|---|
| `.claude/hooks/pretooluse-command-guard.sh`, Kopfkommentar | „blockt Host-Paketmanager und Host-Toolchains“ und „Bewusst NICHT geprüft: andere Interpreter“ (Zeile 1 und 2 des Feldes) | Liefer-Punkt 3b: auf den neuen Stand. |
| [`AGENTS.md`](../../../../AGENTS.md) §3.1 Absatz „Durchsetzung“ | „Sprach-Toolchains und in-place Textwerkzeuge liest er nicht“ (Zeile 2 des Feldes) — der Satz wird falsch | Liefer-Punkt 3c. |
| [`AGENTS.md`](../../../../AGENTS.md) §3.1 Satz zur Mutationsprobe (Kopie im Scratchpad) | bleibt wahr; der Zusatz nennt, dass `-i` auch auf der Kopie geblockt wird | Liefer-Punkt 3c. |
| `.claude/commands/implement-slice.md` (Aufzählung der Regel; Commit via Message-Datei) | nennt die Formen und verweist auf §3.1; „der Guard scannt den Command-String, also nie eine Commit-Message inline, die ein geblocktes Tool-Token enthält“ (Zeile 2 des Feldes) | Fixrunde: der Satz „der Guard scannt den Command-String“ nennt den Stand `diff` nicht mehr (Anführungszeichen-Lesung: ein Tool-Token in den Anführungszeichen einer Inline-Message blockt nicht); der Punkt und die Schritt-20-Zeile zum Konjunktiv-Kandidatenlauf (F-7) sind angepasst. |
| `.harness/skills/reviewer.md` §HIGH „Docker-only-Verstoß“ | nennt die Formen als Prüfgegenstand des Reviews (Zeile 3 des Feldes) | bleibt wahr: der Guard fängt die Flag-Formen, das Review den Rest; unberührt. |
| Beobachtungs-Register `BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel`, `state.md` | „Durchsetzung heute: das Review“ und „noch kein Plan angelegt“ (Zeile 4 des Feldes am Parent) | Fremde Datei, deshalb Meldung des Implementers mit Frist „Closure“; die Closure zog den Text nach: der Ausgang nennt Guard, `MR-003` und Tabellentest, die Grenze des Guards und die Adresse der Nutzer-Frage zu `tools/harness/blocked/go` (Suchlauf-Zeile 4 steht bei 0). |
| `docs/plan/planning/welle-transformationen.md` §6 „Fragen für den nächsten Architect-Zug“ (a) | „in keinem committeten Text verboten“ (Zeile 5 des Feldes) — überholt seit dem Architect-Zug `17cb4eb3`, nicht erst durch diesen Slice | gemeldet, nicht mitgeändert (fremde Datei); die Closure zog den Absatz nach: Verbot in `AGENTS.md` §3.1, Guard, Zähler und Ausgang des Register-Eintrags (Suchlauf-Zeile 5 steht bei 0). |
| `docs/plan/planning/open/slice-code-kommentare-bereinigung.md` §8 („offen, Schwelle erreicht, 3×“) | die Sichtung nennt den Register-Eintrag „offen“, sein `state.md` sagt seit `17cb4eb3` „verkörpert“ (Zeile 6 des Feldes trifft drei Zeilen: diese, die Risiko-Zeile in §6, die nur die Zahl nennt und wahr bleibt, und die Notiz in `welle-transformationen`) | gemeldet wie oben; die Closure zog §6 und §8 nach („verkörpert, 4×“; der Guard blockt die Flag-Formen). |

## 4. Trigger

**Start** (`next` → `in-progress`): kein anderer Slice liegt in `in-progress/`
(WIP-Limit 1; beobachtbar: `ls docs/plan/planning/in-progress` nennt nur `roadmap.md`).
Die Änderung am Guard wirkt in jeder Sitzung, die das Repo öffnet, ab dem Commit — ein
zweiter Implementer-Lauf in derselben Zeit liefe unter einem Guard, der sich unter ihm
ändert; das WIP-Limit deckt das.

**Reihenfolge (Empfehlung an den Orchestrator):** nach `slice-antragsqueue-lesefehler-failed`
(dort §4: 1. `slice-code-kommentare-kennungen`, 2. `slice-harness-fmt-check`, 3. der
Queue-Fix, danach die Welle [welle-transformationen](../welle-transformationen.md)) und
**vor dem ersten Slice der Welle** (`slice-transformationen-map-value`); in jedem Fall vor
`slice-code-kommentare-bereinigung`, dessen §8 ihn den Ort nennt, an dem die Klasse eintreten
kann (Hunderte Kommentar-Änderungen). Grund für die Lage vor der Welle: zwei der drei
belegten Vorfälle liegen in Slices dieser Welle
(`slice-transformationen-antragsweg-usecase`, `slice-transformationen-backfill-pfad`), und
jeder Folge-Slice läuft danach unter dem Guard. **Eine technische Kante gibt es nicht**:
kein anderer offener Plan nennt `.claude/hooks/`, `pretooluse` oder `harness/conventions`
(`git grep -n -E '\.claude/hooks|pretooluse|harness/conventions'` über `docs/plan/planning/open`
und `welle-transformationen.md` trifft nur diesen Plan, gemessen am Parent `89d427e0`);
deshalb braucht [welle-transformationen](../welle-transformationen.md) §5 (Kanten zu
wellenlosen Slices) keinen Nachtrag, und die Reihenfolge-Sätze der anderen Pläne bleiben
wahr (sie zählen die Slices auf, die bei ihrer Anlage vorlagen).

**Nebenwirkung — für jeden Folge-Slice.** Ein Guard, der zu viel blockt, hält den Lauf eines
Implementers, Reviewers oder Verifiers an, ohne dass ein Test es merkt. Die Abwehr sitzt in
der Tabelle der Nicht-Treffer (§2 Liefer-Punkt 1) und in der Blockmeldung, die den Ersatzweg
nennt; der Rest ist ein bekannter Rand, §6 Risiko 1.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der zweite Liefer-Punkt (Host-Interpreter
  auf Repo-Pfaden) ist der abtrennbare Teil — er ist unter den drei Vorfällen nur einmal belegt und
  trägt die größte Falsch-Positiv-Fläche; trägt die Tabelle seiner Nicht-Treffer nicht, geht er als
  eigener Slice zurück und der erste Liefer-Punkt (`sed`, `perl`, `awk`) landet allein.
- `in-progress` → `open` (blockiert): (a) ein Alltags-Aufruf der Rollen, der nach den Regeln
  geblockt werden müsste, obwohl er zulässig ist (die Lesung der Flag-Tokens ohne
  Quote-Kenntnis trägt ihn nicht), und der Weg dahin verlangt Quote-Bewusstsein im Guard —
  dann Architect-Frage (Shell-Parser ja oder nein). Der Fall ist im Review eingetreten
  (Finding F-2: ein Muster mit `|` blockte); die Rückführung ist nicht eingetreten: der Guard
  löst ihn mit einem Zustandsautomaten über Anführungszeichen und Backslash (Fixrunde §3),
  der Architect hat ihn als Lexer-Stufe ratifiziert, `MR-003` trägt die Obergrenze; (b) der
  Hook liest die Block-Ausgabe in der Live-Sitzung nicht wie im Bestand (der Beleg in
  Liefer-Punkt 1 fehlt) — dann Architect-Frage zur Wirkung der Durchsetzungsschicht; nicht
  eingetreten (Risiko 4).

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + Review-Report ohne offenes HIGH oder MEDIUM +
Verifikation, dass die DoD trägt (`make test-command-guard` grün, die Mutationen rot gesehen,
der Live-Beleg der Blockmeldung im Bericht) + Closure-Notiz mit Lerneintrag geschrieben
(geschärfte Regel: [`AGENTS.md`](../../../../AGENTS.md) §3.1 Durchsetzung und `MR-003`;
neuer Sensor: `make test-command-guard`, ohne Gate).

## 6. Risiken und offene Punkte

- **1. Falsch-positive Blockaden legitimer Aufrufe.** Der Guard liest Anführungszeichen und
  Backslash (`tools/harness/mask-quotes.awk`): ein Trenner (`;`, `&`, `|`, `(`, Zeilenende) in
  Anführungszeichen startet kein Segment, ein Muster `'sed -i|perl -pi|awk -i'` in einem
  `git grep` blockt nicht. Ränder bleiben: die Zeilen eines Heredocs gelten als Kommando-Zeilen,
  ein unbalanciertes Anführungszeichen (Apostroph im Heredoc-Text) segmentiert ohne
  Anführungszeichen-Kenntnis, `\;` trennt, ein `sed -i` auf einer Scratchpad-Kopie blockt und ein
  `python3`-Aufruf, dessen Text einen Repo-Namen der obersten Ebene nennt, ohne die Datei zu
  schreiben. *Erwartet, zu belegen durch:* die Nicht-Treffer-Fälle des Tabellentests (die
  Alltags-Formen der Rollen: `sed -n`, `sed s/…/…/ f > g`, `awk` lesend, `git grep -E` mit
  Alternation, `git commit -F`) und die benannten Rand-Fälle mit erwartetem Block; die
  Blockmeldung nennt den Ersatzweg. **Ausgang:** *eingetreten und behandelt* — der Review
  fand ein `\|`-Muster in einem lesenden `git grep` geblockt (F-2, viermal live); die
  Rückführung nach `open/` trat nicht ein, der Guard löst es mit dem Anführungszeichen-Lexer
  (37 Fälle der Gruppe „Anführungszeichen“). Die Ränder stehen als 8 benannte Falsch-Positiv-Fälle
  und 16 benannte Grenzen im Tabellentest (gezählt bei dieser Closure, Gruppen `benannte
  Falsch-Positiv-Raender` und `benannte Grenzen`); die Kehrseite der Quote-Lesung (zwei Apostrophe
  in zwei Heredoc-Zeilen, Falsch-Negativ, V-2) steht in `MR-003` und ist gebunden.
- **2. Der Guard verspricht mehr, als er kann** (die Harness-Lüge der Baseline). Nicht
  gelesen werden: Umleitungen und flaglose Schreibwege (`> datei`, `tee`, `dd of=`,
  `sed … > tmp && mv`); ein `cd` im selben Kommando (der Pfad-Test nimmt die Repo-Wurzel als
  Arbeitsordner); `$PWD`, `$(git rev-parse --show-toplevel)`, `~`, Variablen und Glob-Muster als
  Pfad; ein Skript, das der Interpreter liest (`python3 /tmp/x.py`, dessen Text Repo-Dateien
  schreibt); `bash skript.sh`, dessen Inhalt `sed -i` trägt (nur `bash -c "…"` wird
  rekursiv gelesen); `-exec` nur für die drei Werkzeuge dieses Slice; jedes andere
  in-place-fähige Werkzeug (§1); die Rezepte hinter `make` und die Docker-Bauten. *Erwartet, zu
  belegen durch:* die Grenz-Zeile von `MR-003` nennt diese Punkte, der Kopfkommentar des Guards
  und [`AGENTS.md`](../../../../AGENTS.md) §3.1 „Durchsetzung“ sagen „Stolperdraht, keine
  Sandbox“ und verweisen darauf; der Reviewer liest die drei Stellen gegeneinander.
  **Ausgang:** *eingetreten und behandelt; Rest weiter offen als benannte Grenze* — der Review
  fand Lücken, die die Grenz-Zeile nicht nannte (F-4 Schlüsselwörter, F-6 weitere Formen); sie
  sind gelöst oder in `MR-003` genannt, Kopfkommentar und `AGENTS.md` §3.1 tragen die Kurzform
  (die Verifikation las die drei Stellen gegeneinander). Weiter offen als Kenntnis: Umleitungen
  und flaglose Schreibwege, der Optionswert hinter `sudo -u x`, Host-`python`/`perl` auf einem
  Pfad ohne Repo-Namen (kein Weg nach `AGENTS.md` §3.1), das Falsch-Negativ der Quote-Lesung
  (Heredoc), die Kopf-Liste `tools/harness/blocked/go` (Nutzer-Entscheidung, Adresse im
  `state.md` des Register-Eintrags).
- **3. Die Änderung schwächt den Bestandsschutz oder macht den Guard fail-open.** Der Guard
  hat heute keinen Test (Suchlauf-Zeile 1: drei Fundstellen, keine Testdatei darunter).
  *Erwartet, zu belegen durch:* die Bestandsregeln stehen als Fälle im
  Tabellentest (Paketmanager am Kopf, Präfixe, `bash -c`-Rekursion, Tiefe über 3, defektes
  JSON, fehlendes `awk`) und eine Mutation an der neuen Erkennung färbt sie nicht rot.
  **Ausgang:** *entfallen* — die Gruppe „Bestandsregeln“ (18 Tabellenzeilen) ist am Guard des
  Parents grün (Verifikation: 0 `FEHLER:`-Zeilen), die Fälle Tiefe 4, defektes JSON, abgeschnittenes
  JSON, `\u`-Escape und fehlendes `awk` binden fail-closed (Mutationen rot, Verifikation §4). Der
  Review fand eine Erweiterung der Bestandsklasse (F-3: die Kopf-Erkennung gilt für alle Klassen);
  sie verschärft, lockert kein Gate (`AGENTS.md` §3.6) und steht in `MR-003`.
- **4. Die Wirkung in der Live-Sitzung bleibt unbelegt, weil der Tabellentest nur die
  Hook-Schnittstelle (JSON in, JSON aus) prüft.** Ob Claude Code die Ausgabe wie im Bestand
  liest, ist eine Eigenschaft der Laufzeit, kein Ergebnis des Skripts
  ([`AGENTS.md`](../../../../AGENTS.md) §3.10 nennt dieselbe Klasse für externe Läufe).
  *Erwartet, zu belegen durch:* der Live-Beleg der Klasse `inplace` in Liefer-Punkt 1 (ein
  Aufruf mit nicht vorhandenem Ziel wird geblockt, der Wortlaut steht im Bericht); für die Klasse
  `interp` trägt der Tabellentest die Hook-Schnittstelle, ein Live-Aufruf mit Host-`python` entfällt
  (`AGENTS.md` §3.1). **Ausgang:** *entfallen* — der Live-Block der Klasse `inplace` ist in den
  Sitzungen von Implementer und Verifier gesehen, der Wortlaut im Verifikations-Report §1 gedruckt;
  der Reviewer sah beide Klassen live geblockt (übernommen aus dem Review-Report); der Ausgabeweg
  `emit_block` ist für alle drei Klassen derselbe. Die Closure-Sitzung des Planners sah die Klasse
  `interp` ebenfalls live: ein versehentlicher `python3 --version`-Aufruf mit einem Repo-Pfad im
  Befehlsstring kam als Hook-Fehler „Host python/perl on repo paths is blocked (AGENTS.md Hard Rule
  3.1)“ zurück und lief nicht (gemessen; ein Falsch-Positiv-Rand der Klasse, `MR-003`).
- **5. Die Evidenz trägt nicht alle Regeln** (Baseline: „Eine Wächter-Regel ohne
  Sensor-Evidenz ist Aufwand ohne Begründung und fällt beim ersten Fehlalarm“, Modul 13
  §Guard-Härtung). `sed -i` ist in fünf Läufen über drei Vorfälle belegt, `perl -pi` und
  `awk -i inplace` in keinem, der Host-Interpreter in zwei Vorgängen (ein Heredoc auf einer
  Repo-Datei, zwei Aufrufe im vierten Vorgang; Zahlen übernommen aus den Beleg-Dateien, §1,
  §3 Abweichungen). Die Regeln zu `perl` und `awk` sind dieselbe Klasse mit anderem Namen und
  kosten je eine Erkennung; die zum Host-Interpreter ist die schwächste (Substring-Heuristik,
  größte Falsch-Positiv-Fläche) und der abtrennbare Teil (§4). *Erwartet, zu belegen durch:* der
  Auslöser steht in `MR-003` mit den gezählten Beleg-Dateien; der Reviewer prüft, dass die Zahlen
  dort mit dem Register übereinstimmen; ein Auftreten **trotz** Guard ist eine neue Beleg-Datei
  im Register-Eintrag (Kenntnis, in diesem Slice nicht belegbar). **Ausgang:** *weiter offen als
  Kenntnis* — `MR-003` nennt den Auslöser mit den vier Beleg-Dateien (Review und Verifikation
  stellten fest, dass die Zahl mit dem Register übereinstimmt); `perl -i` und `awk -i inplace`
  bleiben ohne Beleg, die Host-Interpreter-Regel bleibt die schwächste (Neubewertungs-Trigger im
  Auflösungs-Trigger von `MR-003`); ein Auftreten trotz Guard ist eine weitere Beleg-Datei im
  Register-Eintrag.
- **6. Mutationsproben der Reviewer und Verifier werden behindert.** Sie arbeiten nach
  [`AGENTS.md`](../../../../AGENTS.md) §3.1 auf einer Kopie im Scratchpad; ein `sed -i` auf der
  Kopie wird jetzt geblockt. *Erwartet, zu belegen durch:* die Tabelle enthält den zulässigen
  Weg als Nicht-Treffer (`sed s/a/b/ Datei > /tmp/…/Kopie`, `python3 /tmp/…/mutate.py`), der
  Satz zur Mutationsprobe in §3.1 nennt ihn, und der Reviewer des Slice fährt eine Mutation
  auf diesem Weg. **Ausgang:** *entfallen* — der Weg (`sed s/a/b/ Datei > Kopie`) passiert den
  Guard (Tabellenfall; Review und Verifikation fuhren ihre Mutationen so, der Reviewer mit `sed`
  ohne `-i`, der Verifier mit `awk` nach stdout); ein Host-`python3 /tmp/…/mutate.py` passiert
  ihn ebenfalls, ist aber kein Weg nach `AGENTS.md` §3.1 (Satz im Absatz „Durchsetzung“).
- **7. Host-Werkzeuge außerhalb von Docker und `make`** (`bash`, `awk`, `mktemp` im Test,
  `ls` im Guard; `BEO-PGC/host-werkzeug-jenseits-docker-und-make-ohne-deklaration`, 2×).
  *Erwartet, zu belegen durch:* alle vier stehen in der Klasse „Host-Werkzeug ohne
  Installation“ ([`AGENTS.md`](../../../../AGENTS.md) §3.1: `bash`, `git`, `awk`, `mktemp` und
  die coreutils-Basis); die README-Zeile nennt sie. **Ausgang:** *entfallen* — der Kopf des
  Tabellentests und die README-Zeile nennen `bash`, `awk`, `mktemp`; `BEO-PGC/host-werkzeug-jenseits-docker-und-make-ohne-deklaration`
  bleibt bei 2×, der Trigger ist nicht eingetreten (`state.md` des Eintrags).

**Offene Fragen an den Nutzer (Entscheidung, nicht Umfang) — je mit Adresse:**

1. **Die Kopf-Liste `tools/harness/blocked/go`** (`go gofmt python python3 node dotnet java
   gradle uv`; die Fragment-Ladung liegt im Guard, das Verzeichnis existiert nicht). Sie
   blockte unbedingt am Kopf eines Segments: sie schlösse die Toolchain-Aufrufe auf dem Host,
   die kein Vorfall belegt, und machte Liefer-Punkt 2 für `python` überflüssig; `docker run … go`
   ginge durch, weil `docker` der Kopf ist. Empfehlung: nicht in diesem Slice; zuerst die
   Wirkung von Liefer-Punkt 1 abwarten. Adresse: der Nutzer; bei „ja“ ein eigener Plan, den der
   Planner anlegt. Die Frage steht mit dem Closure-Stand im `state.md` des Register-Eintrags.
   **Stand bei Closure:** in diesem Slice nicht umgesetzt (Entscheidung des Auftraggebers an den
   Implementer, 2026-09-26); die Neubewertung hat zwei Trigger in `MR-003` (Auflösungs-Trigger),
   das zweite Kriterium tritt mit der Closure von `slice-transformationen-map-value` ein
   (Adresse im `state.md` des Register-Eintrags).
2. **Eine Scratchpad-Ausnahme für `sed -i`** (erlaubt, wenn ein Token auf einen absoluten
   Pfad unter `/tmp/` zeigt und keines auf einen Repo-Pfad). Der Plan blockt unbedingt: die
   drei Vorfälle klassifizieren auch das `sed -i` auf `/dev/null` und auf einer
   Scratch-Datei als Fehlgriff, und der zulässige Mutationsweg braucht kein `-i`
   (Stdout-Umleitung). Eine Ausnahme kostet Pfad-Logik und ein Loch (`xargs sed -i < /tmp/liste`
   ginge durch). Adresse: der Nutzer; bei „ja“ ändert der Implementer den Fall im Tabellentest,
   und `MR-003` nennt die Ausnahme. **Stand bei Closure:** entschieden, keine Ausnahme (Auftrag des
   Auftraggebers an den Implementer, 2026-09-26); der Guard blockt `sed -i` unbedingt.

## 7. Closure-Notiz

- **Was hat funktioniert:** (1) Die Leser-Kette fand, was der Implementer-Lauf nicht fand:
  der Review nennt 1 HIGH, 4 MEDIUM, 2 LOW und 2 INFO, die Verifikation 8 Verifier-Findings
  (1 MEDIUM, 3 LOW, 4 INFO; übernommen aus den Reports). Der HIGH (F-1, Zeichenklassen ohne
  Mitglied-Fälle) und der Eintritt der Rückführung §4 (a) (F-2, ein Suchmuster mit `|` blockte
  viermal live) sind Funde durch Mutation und durch den Live-Lauf, nicht durch Lesen. (2) Die
  Mutationen der Eingabeseite: der Implementer fuhr 120 (118 rot, 2 äquivalent; übernommen aus
  dem Implementer-Bericht, ohne Lauf-Anker, V-6), der Reviewer 41 Kopien (übernommen aus dem
  Review-Report), der Verifier 60 (57 rot, 3 grün; übernommen aus der Verifikation §4); die drei
  grünen des Verifiers (V-3 Ziffernklasse `0-27` und `0-9`, V-4 Backslash-Escape) band der Nachzug
  `527d69dd` mit drei Mutationstabellen (6 + 10 + 3 Läufe, abgeleitet; gemessen vom
  Implementer-Zug). (3) Der Tabellentest trägt: 317 Fälle, 311 Tabellenzeilen und 6 Sonderfälle
  (gemessen bei dieser Closure: `make test-command-guard` Exit 0, „alle 317 Fälle bestanden“;
  `grep -cE '^(block|pass) '` 311; die 6 abgeleitet), gegen den Guard am Parent `e98d419c` 168
  rot und gegen den Guard am Review-Stand `a0472421` 81 rot (je ein Lauf bei dieser Closure,
  `grep -c '^FEHLER:'`); die Gruppe „Bestandsregeln“ (18 Tabellenzeilen) ist am Parent grün.
  (4) Die Wirkung ist live gesehen: die Klasse `inplace` in den Sitzungen von Implementer und
  Verifier (Wortlaut im Verifikations-Report §1), die Klasse `interp` vom Reviewer (übernommen)
  und in der Closure-Sitzung (ein versehentlicher `python3 --version` mit Repo-Pfad im Befehlsstring
  wurde geblockt, ohne Wirkung). (5) Der Architect-Zug `d55943b1` beantwortete die Fragen in einer
  Runde: Quote-Lesung als Lexer-Stufe ratifiziert mit Obergrenze in `MR-003`, die Erweiterung der
  Bestandsregel (F-3) akzeptiert, ein Satz zu Host-`python` ohne Repo-Pfad in `AGENTS.md` §3.1,
  der Tabellentest als Beleg für Liefer-Punkt 2. (6) Das Suchlauf-Feld trug: 15 Zeilen; der Verifier
  fuhr fünf von Hand nach; bei dieser Closure meldet `make suchlauf-nachmessen` mit den nachgezogenen
  Soll-Werten Exit 0, gedruckt „suchlauf-nachmessen: 15 Zeilen stimmen“ (gemessen).
- **Was ging anders als geplant:** (1) Der Plan schloss „Quote-Bewusstsein im Guard“ aus und führte
  die Rückführung §4 (a); die Fixrunde lieferte einen Zustandsautomaten (`mask-quotes.awk`, 40 Zeilen)
  statt der Rückführung, der Architect ratifizierte ihn nachträglich (V-1: die Umkehr eines
  Ausschlusses ohne Rückführung; die Bestätigung stand im Plan nur als Wunsch, jetzt steht sie
  als Ergebnis in §1 und §4). (2) Das Register trug vier Beleg-Dateien statt der geplanten drei;
  der Host-Interpreter ist in zwei Vorgängen belegt (Risiko 5). (3) Über den Plan hinaus entstanden
  `mask-quotes.awk`, die Gruppe „Kopf-Erkennung“ (Wrapper-Optionen und Schlüsselwörter für alle
  Klassen, F-3, F-4), der Konjunktiv-Kandidatenlauf auf Skripten (F-7) und 317 statt der geschätzten
  rund vierzig Fälle. (4) Zahlen im Plan wichen von der Nachmessung ab: „31 Fälle Anführungszeichen“
  gegen gemessen 37, „24 Fälle Pfadzeichen“ gegen 23, „80 rot am Review-Stand“ gegen 81; die Closure
  hat sie mit Anker nachgezogen (Deckel-Fall unten). (5) Die Fixrunde hat kein zweites Review (V-8);
  der Closure-Trigger „Review-Report ohne offenes HIGH oder MEDIUM“ stützt sich auf den Report am
  Stand `a0472421`, dessen Findings der Verifier je Zeile am Ist-Zustand nachmaß (Verifikation §5:
  F-1 bis F-7 aufgelöst, F-1 mit einem Rest, der im Nachzug `527d69dd` gebunden ist) und 60 Mutationen
  an der Fixrunde; ein Nach-Review ist nicht beauftragt, die Entscheidung liegt hier beim Planner.
  (6) Die Coverage-Zahl ist an diesem Stand nicht lauf-stabil: 85.10 % im Lauf dieser Closure und im
  Lauf der Verifikation (übernommen), 85.20 % im Auftrag zu dieser Closure genannt (übernommen); die
  Streuung ist 0,10 Prozentpunkte (abgeleitet), die Zahl ist Beleg des jeweiligen Laufs, keine
  Zustandsgröße; der Diff berührt keinen Go-Code. (7) Der Live-Beleg mit Host-`python3` (DoD-Zeile 2)
  entfällt: er verstieße gegen `AGENTS.md` §3.1; der Architect-Zug hat den Tabellentest als Beleg
  angenommen.
- **Steering-Loop-Eintrag (Lerneintrag):** *(a) Geschärfte Regel.* [`AGENTS.md`](../../../../AGENTS.md)
  §3.1 Absatz „Durchsetzung“ nennt, was der Guard blockt (`sed -i`/`--in-place`, `perl -i`,
  `awk -i inplace` unabhängig vom Ziel, Host-`python`/`perl` mit Repo-Pfad im Befehlsstring), was er nicht
  liest, und dass ein Host-Interpreter auch ohne Repo-Pfad kein zulässiger Weg ist; `MR-003` trägt Vertrag,
  Obergrenze der Quote-Lesung und Grenz-Zeile · seit slice-harness-guard-inplace-textwerkzeug. Herkunft:
  `BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel` (vier Beleg-Dateien),
  `BEO-PGC/host-werkzeug-jenseits-docker-und-make-ohne-deklaration`. Die Kopf-Erkennung (Wrapper-Optionen,
  Schlüsselwörter) gilt für alle Klassen, auch die Paketmanager (F-3). *(b) Neuer Sensor, kein Gate.*
  `make test-command-guard` (Tabellentest gegen den Guard, 317 Fälle; `harness/README.md` §Sensors,
  `Makefile`); ein Gate braucht eine ADR ([`AGENTS.md`](../../../../AGENTS.md) §3.6), und ein Wächter
  gehört nicht in die Gate-Tabelle (Baseline `modul-13-quality-gates.md` §Guard-Härtung). *(c) Benannte
  Lücken, je mit Adresse* (alle in der Grenz-Zeile von `MR-003`): erstens das Falsch-Negativ der
  Quote-Lesung (zwei Apostrophe in zwei Heredoc-Zeilen maskieren die Zeilen dazwischen, V-2; Adresse:
  die Obergrenze in `MR-003`, ein Fall, der mehr verlangt, ist eine Beleg-Datei im Register-Eintrag, keine
  Erweiterung des Maskierers); zweitens der Optionswert hinter `sudo -u x` und `env -u X` (Adresse:
  `MR-003`, Ereignis: ein Beleg); drittens Host-`python`/`perl` auf einem Pfad ohne Repo-Namen (Adresse:
  die Kopf-Liste `tools/harness/blocked/go`, Nutzer-Entscheidung, Trigger in `MR-003` und im `state.md`
  des Register-Eintrags); viertens Umleitungen und flaglose Schreibwege, Skripte, andere Interpreter
  (Adresse: das Review, Ausschluss eines Sensors über Dateiinhalte, §1). Die Nutzer-Frage zu
  `tools/harness/blocked/go` hat ihre Adresse: das zweite Trigger-Kriterium („die Closure des nächsten Slice,
  dessen Läufe unter diesem Guard liefen“) tritt mit der Closure von `slice-transformationen-map-value` ein;
  der Guard wirkt ab dem Commit in jeder Sitzung, die Läufe aller Folge-Slices laufen unter ihm (der Plan
  der Folge-Slices muss das nicht tragen). *(d) Finding-Klassen des Reviews und der Verifikation:*
  Zusage ohne Bindung an ihre Eingabeseite (F-1, V-3, V-4) · Falsch-Positiv-Rand ungebunden und nicht
  benannt (F-2) · Bestandsregel verändert, als unverändert ausgewiesen (F-3) · Grenz-Zeile nennt eine
  Lücke der Zusage nicht (F-4, F-6, V-2) · Meldung und Regel nennen verschiedene zulässige Wege (F-5) ·
  Selbstprüfung enger als der Prüfpunkt (F-7) · Wächter-Regel ohne Evidenz und Interpreter-Denylist (F-8) ·
  äquivalente Mutation durch Doppelschutz (F-9) · Ausschluss des Plans ohne Rückführung umgekehrt (V-1) ·
  Vorher-Nachher-Sprache in Doku-Prosa (V-7) · Mutationszahl ohne Lauf-Anker (V-6); ihre Zuordnung zu den
  Register-Zählern steht im nächsten Punkt.
- **Beobachtungs-Register (`../observations/`):** je Anfall eine Datei
  `evidence/slice-harness-guard-inplace-textwerkzeug.md`, Zähler = Zahl der Dateien. *Neue Belege:*
  `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` **17×** (F-1, HIGH; V-3 und V-4 im selben Vorgang;
  Ausprägung: die Zeichenklasse eines Musters, ihre Eingabeseite ist jedes Mitglied und jede Klassengrenze;
  Ausgang unverändert **verkörpert**, kein Kandidat der Schärfung: Reviewer und Verifier fanden ihn mit der
  bestehenden Regel); `BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung` **9×** (F-3 und F-5, je
  MEDIUM; „unverändert“ ohne Messung am Parent, „zulässiger Weg“ ohne Anker); `BEO-PGC/regel-weiter-als-ihr-sensor`
  **4×** (F-4, F-6, V-2; neue Domäne, Variante „benannte Grenze unvollständig“; Urteil: kein zusätzlicher Sensor,
  wie bei den drei früheren Belegen); `BEO-PGC/vorher-nachher-sprache-in-test-harness-kommentar` **6×** (V-7 und
  die Plan-Prosa dieser Closure, F-7; die Schärfung des Kandidatenlaufs in `implement-slice` Schritt 20 ist
  **verkörpert** und im `state.md` nachgezogen, der frühere Kandidat „Adresse Lese-Schritt“ war überholt).
  *`state.md` fortgeschrieben:* `BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel` (**verkörpert** durch
  Regel, Guard, `MR-003` und Tabellentest; die Nutzer-Fragen: Scratchpad-Ausnahme entschieden — keine,
  `blocked/go` mit Trigger und Adresse; die Beleg-Datei des vierten Vorgangs trägt als dritte Rolle den
  Planner der Closure-Sitzung von `slice-antragsqueue-lesefehler-failed`, Angabe des Auftraggebers,
  Zähler bleibt 4×; aus diesem Slice fällt kein Beleg an; ein versehentlicher `python3`-Aufruf der
  Closure-Sitzung wurde vom Guard geblockt),
  `BEO-PGC/host-werkzeug-jenseits-docker-und-make-ohne-deklaration` (Trigger geprüft, nicht eingetreten,
  Zähler bleibt 2×). *Deckel-Fälle ohne Datei, Kennung hier* (verkörpert ab 10×, vor dem Merge gefunden,
  Schwere ≤ LOW, bekannter Träger-Typ): V-6 (INFO, „120 Mutationen“ ohne Lauf-Anker) und die drei
  Zahl-Abweichungen im Plan (31/37, 24/23, 80/81; Leser: der Planner der Closure, INFO) zu
  `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (Deckel bei 23×). *Kein Anfall:*
  `BEO-PGC/werkzeug-liest-nutzerkonfiguration-ohne-pin` (der Tabellentest liest kein Fremdwerkzeug außer
  `ls -A` im Guard und keine Nutzerkonfiguration), `BEO-PGC/test-schreibt-in-committete-datei` (der Test
  schreibt nur ins Temp-Verzeichnis). *Lese-Schritt:* aus diesem Slice erreicht neu kein Eintrag die Schwelle
  ohne Ausgang (`BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung` liegt bei 9×, Ausgang
  **verkörpert**).
- **Folge-Slices:** keiner. Die Nutzer-Frage zu `tools/harness/blocked/go` ist eine Entscheidung mit
  Adresse (die Closure von `slice-transformationen-map-value`, der Planner legt sie dem Nutzer mit der
  gemessenen Wirkung vor), kein Slice. Träger-Meldungen aus §3 mit Frist „Closure dieses Slice“ sind
  erfüllt: `welle-transformationen` §6 (a), `slice-code-kommentare-bereinigung` §6 und §8 und der
  `state.md` des Register-Eintrags sind nachgezogen; der Link auf den `in-progress/`-Pfad dieses Slice
  steht in keinem anderen Dokument (`git grep` bei dieser Closure).
- **Risiken aus §6:** je ein Ausgang, mit Beleg in §6. *Eingetreten und behandelt:* Risiko 1 (F-2,
  Rückführung nicht eingetreten; der Architect ratifizierte den Lexer), Risiko 2 (Lücken der Grenz-Zeile,
  F-4, F-6, V-2, nachgezogen). *Entfallen:* Risiko 3 (Bestandsschutz gebunden), 4 (Live-Beleg der Klassen
  `inplace` und `interp`), 6 (Mutationsweg passiert den Guard), 7 (Host-Werkzeuge in der Klasse). *Weiter
  offen als Kenntnis:* Risiko 5 (`perl -i` und `awk -i inplace` ohne Beleg; Host-Interpreter-Regel die
  schwächste). Die Rückführung §4 (a) ist nicht eingetreten (ratifiziert, `MR-003` trägt die Obergrenze),
  die Rückführung §4 (b) ebenfalls nicht.
- **Drei Paarungen:** dieser Slice hat keine Welle; die Slice-Closure selbst trägt die drei
  Paarungen (Anker · Folge-Slice · Register), nach dem `git mv` nach `done/`. Die
  Ereignis-Adresse kann eintreten (die Closure dieses Slice), und keine der Paarungen hängt an
  einem Ereignis außerhalb: *Anker* — [`AGENTS.md`](../../../../AGENTS.md) §3.1 Durchsetzung trägt
  `· seit slice-harness-guard-inplace-textwerkzeug` (`git grep -n 'slice-harness-guard-inplace-textwerkzeug'
  -- AGENTS.md harness/README.md harness/conventions`), `MR-003` steht im Index von `harness/conventions.md`
  (Anker `mr-003`) und nennt den Slice als Herkunft der Obergrenze, das Ziel `test-command-guard` steht im
  `Makefile` und in `harness/README.md` (Zeile mit `· seit`), der Kopfkommentar des Guards verweist auf
  `MR-003`; *Folge-Slice* — keiner genannt; die Ereignis-Adresse der Nutzer-Frage (die Closure von
  `slice-transformationen-map-value`, Datei in `open/`) kann eintreten; *Register* —
  `BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel` und jede weitere genannte Kennung
  (`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`,
  `BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung`, `BEO-PGC/regel-weiter-als-ihr-sensor`,
  `BEO-PGC/vorher-nachher-sprache-in-test-harness-kommentar`,
  `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`,
  `BEO-PGC/host-werkzeug-jenseits-docker-und-make-ohne-deklaration`) existiert als Verzeichnis mit nicht
  leerem `evidence/` (geprüft mit `ls docs/plan/planning/observations/BEO-PGC/<slug>/evidence`).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area `*`/`PGC`
(Greenfield); `.claude/hooks`, `tools/harness`, `harness/` und `AGENTS.md` sind keine eigenen
Sub-Areas — kein Anlass zur Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (Zähler gemessen am
2026-09-26 mit `ls evidence | wc -l` je Eintrag):

- `BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel` (verkörpert, 3×): der Gegenstand
  dieses Slice; sein `state.md` nennt den Slice als Adresse des Guard-Ausbaus; Risiko 5 und 6
  in §6.
- `BEO-PGC/host-werkzeug-jenseits-docker-und-make-ohne-deklaration` (verkörpert, 2×): Risiko 7
  in §6.
- `BEO-PGC/regel-weiter-als-ihr-sensor` (verkörpert, teilweise, 3×): die Regel in
  [`AGENTS.md`](../../../../AGENTS.md) §3.1 reicht weiter als der Guard nach diesem Slice; die
  benannte Grenze (Risiko 2, `MR-003`) ist die Antwort dieses Slice auf die Klasse; ein
  Sensor über die Lücke ist ausgeschlossen (§1, erster Punkt).
- `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (verkörpert, 16×): die Mutationen der
  Eingabeseite je Zusage in Liefer-Punkt 1 und 2, mit Nicht-Treffer neben jedem Treffer.
- `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` (verkörpert, 14×),
  `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert, 32×),
  `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (verkörpert, 23×): das Suchlauf-Feld in
  §3 und die Zahlen mit Ursprung; die zwei überholten Träger sind gemeldet.
- `BEO-PGC/slice-chronik-in-code-kommentar` (verkörpert, 9×),
  `BEO-PGC/vorher-nachher-sprache-in-test-harness-kommentar` (verkörpert, 4×),
  `BEO-PGC/kommentar-herkunft-als-kette` (verkörpert, 2×): der Kopfkommentar des Guards und
  des Tabellentests sind Kommentare in Skripten (Liefer-Punkt 3b; Indikativ, höchstens eine
  Kennung, keine Slice-Nummer).
- `BEO-PGC/test-schreibt-in-committete-datei` (verkörpert, 4×): der Tabellentest schreibt nur
  ins Temp-Verzeichnis und liest den Guard; kein Auftreten möglich.
- Gesichtet, ohne Bezug zu diesem Slice: die übrigen Einträge des Registers.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
