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
- **Ein Shell-Parser im Guard.** Derselbe Grund wie bei den Umleitungen. Der Guard liest
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
      Paketmanager-Namen, Sub-Shell-Rekursion, fail-closed bei Parse-Zweifel) bleiben;
      die Kopf-Erkennung (Wrapper-Optionen, Schlüsselwörter) erweitert sie für alle Klassen
      (§3 „Fixrunde“). *Zu belegen durch:* `make test-command-guard` (Tabellentest, netzlos,
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
- [ ] **Liefer-Punkt 2 — Host-Interpreter auf Repo-Pfaden.** Ein Segment mit Kopf `python`,
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
      'print' /tmp/x`. Die Blockmeldung nennt Edit/Write und den Weg über eine absolute
      Scratchpad-Kopie. *Zu belegen durch:* dieselben Tabellentest-Fälle, je an ihre Eingabe
      gebunden durch eine Mutation: Regel entfernt · nur das Segment statt des ganzen
      Befehls gelesen (der Heredoc-Fall färbt rot) · die Pfad-Zeichen-Bedingung entfernt
      (`/tmp/x/docs/a.py` färbt rot) · `./` nicht erlaubt (`python3 ./tools/x.py` färbt rot) ·
      Datei-Namen der obersten Ebene nicht gelesen (`Makefile`-Fall färbt rot); Live-Beleg
      analog Liefer-Punkt 1 (`python3 tools/nicht-vorhanden.py` wird geblockt, `python3
      --version` läuft). **Offen für die nächste Rolle:** der Live-Beleg mit `python3`; der
      Auftrag an den Implementer verbietet Host-`python3` auch für einen Beleg, Tabellentest und
      Mutationen tragen die Hook-Schnittstelle (§3 „Abweichungen“), der Verifier fährt ihn.
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
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel · neuer
      Sensor · benannte Spec-Lücke).
- [ ] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — der Ausgang von
      `BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel` (`state.md`: die
      Durchsetzung heute, die Adresse aufgelöst, die zwei Nutzer-Fragen aus §6 mit
      Adresse); keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7
      notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von
      der Slice-Closure selbst (der Slice hat keine Welle; das Ereignis kann
      eintreten).

**Umfang:** S bis M — Schätzung, nicht gemessen: ein Guard-Skript um zwei Erkennungen,
ein Tabellentest mit rund vierzig Fällen, ein Makefile-Ziel, ein `MR`, drei kurze
Träger-Änderungen. Der zweite Liefer-Punkt ist der abtrennbare Teil (§4).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.claude/hooks/pretooluse-command-guard.sh` | update | Liefer-Punkt 1 und 2: Erkennung der in-place Formen (rohe Flag-Tokens, Kopf nach `-exec`/`-execdir`/`-ok`) und des Repo-Pfad-Musters auf `python`/`python3`/`perl`; eine Blockmeldung je Klasse (gültiges JSON); Kopfkommentar im Indikativ ([`AGENTS.md`](../../../../AGENTS.md) §3.7). Die Bestandsregeln bleiben, der Guard bleibt fail-closed. |
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
diff 6 -E 'blockt (Host-)?Paketmanager|Stolperdraht|scannt den Command-String' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline' ':!docs/plan/planning/observations'
diff 135 -E 'sed -i|perl -pi|awk -i' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline' ':!docs/plan/planning/observations'
diff 1 -E 'Durchsetzung heute|noch kein Plan angelegt' -- docs/plan/planning/observations
diff 1 -E 'in keinem committeten Text' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline' ':!docs/plan/planning/observations'
89d427e0 3 -E 'in-place Textwerkzeuge liest er nicht|Bewusst NICHT gepr|blockt Host-Paketmanager und Host-Toolchains' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline' ':!docs/plan/planning/observations'
diff 0 -E 'in-place Textwerkzeuge liest er nicht|Bewusst NICHT gepr|blockt Host-Paketmanager und Host-Toolchains' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline' ':!docs/plan/planning/observations'
diff 7 -E 'test-command-guard' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline' ':!docs/plan/planning/observations'
```

Die ersten sieben Zeilen sind am Parent gemessen (Bezugs-Befehle im Block, Stand `89d427e0`; Ergebnis:
Symbol 3 Fundstellen — der Guard selbst, `.claude/settings.json`, [`AGENTS.md`](../../../../AGENTS.md)
§3.1; Beschreibung 6; Zählwort 9, darunter ein `sed -i` in `sdks/python/Dockerfile`, das im
Docker-Bau läuft und den Guard nicht berührt — er scannt nur die Bash-Aufrufe des Laufs, nicht
die Rezepte hinter `make`).

**Ergebnis am Stand `diff` (Arbeitsbaum nach der Fixrunde; Gefundenes und Nichtgefundenes):**
Symbol 12 Fundstellen (Guard, `.claude/settings.json`, [`AGENTS.md`](../../../../AGENTS.md), `Makefile`,
`harness/README.md`, `harness/conventions.md`, `MR-003`, `mask-quotes.awk`, vier im Tabellentest) — die
neuen Träger tragen den Namen, keiner der 3 Bestandsträger sagt etwas Falsches. Beschreibung 6 (Bestand 6,
dazu `MR-003`, abzüglich `implement-slice.md`): dessen Satz „der Guard scannt den Command-String“ ist
durch die Anführungszeichen-Lesung überholt und nachgezogen (Träger-Tabelle); der Kommentar in
`tools/harness/extract-command.awk` („Stolperdraht, keine Sandbox“) bleibt wahr und ist unberührt.
Zählwort 135: die Zahl ist gewachsen um die neuen Träger (`MR-003`, Tabellentest mit seinen Fällen,
Guard-Kommentar, [`AGENTS.md`](../../../../AGENTS.md)); die Bestandstreffer `reviewer.md` (Prüfgegenstand des Reviews),
`implement-slice.md` und `sdks/python/Dockerfile` bleiben unverändert wahr, die zwei Treffer in
`slice-code-kommentare-bereinigung` §8 und `welle-transformationen` §6 sind gemeldet (Tabelle unten).
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
  Anführungszeichen-Argument ist ein Token (Fixrunde, §3; der erste Lauf las die Flag-Tokens roh und
  ließ ein Token mit Anführungszeichen als Nicht-Flag gelten); (c) das find-Ende (`+`, `\`, `;`) begrenzt die Flag-Suche (`-iname` hinter
  `-exec … +` ist kein sed-Flag); (d) `-okdir` neben `-ok`; (e) der Tabellentest prüft je Block die
  Begründung der Klasse (`pkg`, `inplace`, `interp`) und die JSON-Form.
- **Tabellentest:** mehr Fälle als der Plan schätzt (Zahl im Bericht des Laufs, gemessen mit
  `make test-command-guard`); die Gruppen sind Bestand, Kopf-Erkennung hinter Wrapper-Optionen und
  Schlüsselwörtern, in-place je Form, je Zeichenklassen-Mitglied und Position, Nicht-Treffer,
  Anführungszeichen-Lesung, Host-Interpreter, benannte Falsch-Positiv-Ränder, benannte Grenzen.
- **Liefer-Punkt 2 geliefert.** Die Falsch-Positiv-Fläche ist gebunden: Nicht-Treffer neben jedem Treffer
  (`/tmp/x/docs/a.py`, `python3 --version`, Scratchpad-Mutationsweg), benannte Ränder mit erwartetem
  Block (`cd /tmp/x && python3 tools/x.py`, Text mit Repo-Namen), Mutationen je Zusage rot gesehen (Bericht).
  Der Live-Beleg des zweiten Liefer-Punkts (ein realer `python3`-Aufruf in der Sitzung) entfällt im
  Implementer-Lauf: der Auftrag verbietet Host-`python3` auch für einen Beleg; der Tabellentest trägt die
  Hook-Schnittstelle, der Verifier wiederholt den Live-Beleg in seiner Sitzung (Risiko 4).
- **Nutzer-Fragen aus §6 entschieden** (Auftrag des Auftraggebers an den Implementer, 2026-09-26):
  Frage 1 — `tools/harness/blocked/go` nicht in diesem Slice, zuerst die Wirkung von Liefer-Punkt 1
  abwarten; Frage 2 — keine Scratchpad-Ausnahme, der Guard blockt `sed -i` unbedingt (Mutationsproben
  laufen über Stdout-Umleitung auf einer Scratchpad-Kopie).

**Fixrunde (Anker: Review-Report `review-slice-harness-guard-inplace-textwerkzeug`, Finding je Zeile;
Stand `diff`).** Der Tabellentest zählt 300 Fälle (`make test-command-guard`: „alle 300 Fälle bestanden“,
294 Tabellenzeilen, 6 Sonderfälle); gegen den Guard am Parent `e98d419c` sind 155 rot, gegen den Guard am
Review-Stand `a0472421` 76 (`GUARD=<Datei>`, gemessen). Die Gruppe „Bestandsregeln“ (18 Tabellenzeilen) ist
am Parent grün.

| Finding | Entscheidung und Änderung | Beleg-Anker |
|---|---|---|
| F-2 | Im Guard gelöst, keine Rückführung nach `open/`: `tools/harness/mask-quotes.awk` ist ein Zustandsautomat über `'`, `"`, Backslash sowie `$(`/Backtick in `"…"` (kein Shell-Parser: Heredocs, Umleitungen, Variablen bleiben ungelesen). Ein Trenner, Leerraum oder Zeilenumbruch in Anführungszeichen wird vor der Segmentierung zu einem Steuerzeichen; ein Anführungszeichen-Argument ist ein Token. Ein unbalanciertes Anführungszeichen gibt den Rohstring zurück (Segmentierung ohne Anführungszeichen-Kenntnis: mehr Segmente, nie weniger; ein Apostroph im Heredoc-Text blockt also nur, wenn eine Kommando-Zeile mit `sed -i` folgt). Die Bestandsregel „quote-blind“ ändert sich damit: `git commit -m "a & sed -i b"` und `grep -E 'a\|sed -i ' f` blocken nicht mehr (zwei benannte Ränder des ersten Laufs sind Nicht-Treffer), `\;` bleibt ein Trenner. Fällt der Maskierer aus, blockt der Guard (fail-closed, Fall „Maskierer nicht lesbar“). Latenz 15 ms je Aufruf gegen 11 ms (20 Aufrufe je Guard, gemessen). | 31 Fälle „Anführungszeichen“ (die Suchlauf-Form mit drei Gliedern, `grep -E "sed -i\|perl -pi" f`, `echo 'a\|sed -i x'`, Backslash, Tab, Trenner ohne Leerraum, `$(…)`/Backtick-Kontext); live: `git grep -c -E 'sed -i\|perl -pi\|awk -i' -- <Datei>` läuft in der Sitzung des Implementers ohne Block, `sed -i s/a/b/ /nicht-vorhanden-scratch` blockt mit der Meldung der Klasse `inplace`. Offen an den Architect: Bestätigung, dass ein Zustandsautomat dieser Größe unter „kein Shell-Parser“ fällt (Plan §4 Rückführung (a)). |
| F-1 | Je Mitglied der Zeichenklassen ein Treffer-Fall und je Klassengrenze ein Nicht-Treffer: sed `[nEsrzub]` (Ergänzung `b`: `sed -bi`; Nicht-Treffer `-ei`, `-fi`, `-xi`), perl `[0-7lanpsw]` (Ziffernklasse von `0-9` auf `0-7` verengt: die Argumente von `-0` und `-l` sind oktal; Nicht-Treffer `-ei`, `-Minteger`, `-Ii`, `-Mfeature`), Pfadzeichen `[A-Za-z0-9_./~-]` (acht Mitglieder davor und dahinter, dazu die Zeichen `=` `:` `,` `"` `)` als Nicht-Mitglieder). Die Behauptung „jede Zusage gebunden“ in `MR-003` nennt die Menge (Tabellentest, Mutationen je Klassenmitglied). | 29 Fälle „Zeichenklassen“, 24 Fälle „Pfadzeichen“ in der Gruppe Host-Interpreter; Mutationen unten. |
| F-3 | Das Überspringen der Wrapper-Optionen gilt für **alle** Klassen: die Lücke `env -i pip`, `xargs -n1 pip`, `time -p pip` ist dieselbe wie bei `sed -i` und bliebe für die Bestandsklasse bestehen, wenn das Überspringen nur für `inplace`/`interp` gälte; zwei Kopf-Erkennungen wären eine zweite Wartungsstelle. Das ist eine Erweiterung der Bestandsregel und steht so in `MR-003` (Adaption) und im Tabellentest als eigene Gruppe „Kopf-Erkennung hinter Wrapper-Optionen“; die zwei Fälle, die der erste Lauf „Bestand“ nannte, stehen dort (am Parent rot gemessen). `command -v`/`-V` und `-pv` zeigen an und führen nichts aus: sie blocken nicht; `type`, `which` sind ohnehin nicht Kopf einer Ausführung; `command -p pip` blockt. | Gruppe „Kopf-Erkennung“ (13 Fälle); `command -v`-Fälle am Guard `a0472421` rot (3), am neuen Guard grün. |
| F-4 | Schlüsselwörter `do then else elif if while until !`, führendes `{` und `)` (Funktionsdefinition) werden vor der Kopf-Erkennung übersprungen; ein `case` überspringt bis zum ersten Label, Labels gelten nur nach einem `case` im selben Befehl (`echo $(date) sed -i x` bleibt Nicht-Treffer). Die Form `xargs -n 1 sed -i` (Wert einer xargs-Option) ist mit gelöst. Grenz-Zeile in `MR-003`: Optionswerte anderer Wrapper (`sudo -u x`, `env -u X`, `nice -n 10`), Werkzeug aus einer Variablen, `eval "$cmd"`, Aliase. | Fälle „Schlüsselwort“ (`for`, `while`, `if`, `else`, `elif`, `until`, `!`, drei `case`-Formen, Funktion), am Parent rot. |
| F-5 | Die Meldung der Klasse `interp` nennt Edit/Write, ein Repo-Werkzeug hinter `make` und `sed s/a/b/ file > /path/to/scratch-copy` — die drei Wege aus `AGENTS.md` §3.1, keinen Host-Interpreter; `MR-003` (Begründung) sagt dasselbe und nennt den Interpreter auf einem Pfad ohne Repo-Namen eine Grenze, keinen zulässigen Weg. | Fall „Meldung interp: Weg nach AGENTS.md 3.1“ (prüft Edit/Write, den stdout-Weg und die Abwesenheit von `mutate`). |
| F-6 | Gelöst: `sed -i''`, `\sed`, `busybox`, `/usr/bin/env`, `/usr/bin/sudo`, `gsed`, die Abkürzung `--in-p`/`--i`, `find … -exec sh -c '…'`/`-exec python3 …` (der Rest hinter `-exec` ist ein eigenes Kommando), `eval "…"`; das Bündel `-input.txt` bei sed bleibt ein Treffer, weil GNU-sed es als `-i` mit Suffix `nput.txt` liest (real gemessen: `sed -n 1p -input.txt f` legt die Sicherungsdatei `fnput.txt` an), bei perl enden die Optionen am Skriptnamen (`perl x.pl -input a` blockt nicht). Rest als Grenze in `MR-003`: `-Wpi`, die awk-Abkürzung `--inc=inplace`, andere Interpreter (`node -e`, `ruby -e`, `uv run python`), Optionswerte anderer Wrapper. `MR-003` nennt die Grenz-Zeile als Kurzform in Kopfkommentar und `AGENTS.md`; „gleichlautend“ ist berichtigt (die `AGENTS.md`-Liste ist eine Teilmenge). | Fälle „Position“ (80) und „Grenzen“ (15). |
| F-7 | Der Konjunktiv-Kandidatenlauf in `implement-slice.md` Schritt 20 liest `'*.go' '*.sh' '*.awk'` mit `(//\|#)` und den transliterierten Formen (`waere`, `wuerde`, `haette`); die Urteilsregel nennt `else`-Zweige in Shell-Kommentaren und `${#var}` als zulässige Treffer. Der Lauf auf dem eigenen Diff (`git diff -U0 HEAD`, 80 hinzugefügte Kommentarzeilen) liefert 0 Treffer, auf dem Diff seit `e98d419c` 0 Treffer; die Bestandszeile „`waere` sonst“ im Guard ist entfernt, dazu das „wuerde“ im Kopfkommentar. | Befehl in `implement-slice.md`, Lauf oben. |
| F-8 | Kenntnis: `perl -i` und `awk -i inplace` haben keinen Beleg im Register (Risiko 5, `MR-003` Begründung); Liefer-Punkt 2 ist der Interpreter-Denylist näher als einer Zerlegung (`MR-003` Begründung: beide Regeln blocken nur auf einem Repo-Pfad); das Feld „Ersetzt-Baseline-Regel“ nennt den Satz, der für alles Ungelesene weiter gilt (im Feld ausgeführt). Keine Änderung der Baseline-Zuordnung. | `MR-003`. |
| F-9 | Keine Aktion (äquivalente Mutation durch Doppelschutz; Fixrunde: die Doppelmutation beider Zeilen ist rot). | Mutationen unten. |

**Mutationen der Fixrunde (Zusage · mutierte Eingabe · gesehenes Rot; je `GUARD=<Kopie> MASKER=<Kopie>` an einer
Kopie von Guard bzw. Maskierer, Menge der Erprobung: die 300 Fälle).** 120 Mutationen, 118 rot, 2 äquivalent:

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
| äquivalent | fehlendes awk (nur die Präsenzprüfung: der Exit-Code der `awk`-Extraktion blockt dieselbe Eingabe); Sprung hinter das erste `case`-Label (das Label-Überspringen leistet dasselbe) | grün, Doppelschutz |

| Träger | Befund | Behandlung |
|---|---|---|
| `.claude/hooks/pretooluse-command-guard.sh`, Kopfkommentar | „blockt Host-Paketmanager und Host-Toolchains“ und „Bewusst NICHT geprüft: andere Interpreter“ (Zeile 1 und 2 des Feldes) | Liefer-Punkt 3b: auf den neuen Stand. |
| [`AGENTS.md`](../../../../AGENTS.md) §3.1 Absatz „Durchsetzung“ | „Sprach-Toolchains und in-place Textwerkzeuge liest er nicht“ (Zeile 2 des Feldes) — der Satz wird falsch | Liefer-Punkt 3c. |
| [`AGENTS.md`](../../../../AGENTS.md) §3.1 Satz zur Mutationsprobe (Kopie im Scratchpad) | bleibt wahr; der Zusatz nennt, dass `-i` auch auf der Kopie geblockt wird | Liefer-Punkt 3c. |
| `.claude/commands/implement-slice.md` (Aufzählung der Regel; Commit via Message-Datei) | nennt die Formen und verweist auf §3.1; „der Guard scannt den Command-String, also nie eine Commit-Message inline, die ein geblocktes Tool-Token enthält“ (Zeile 2 des Feldes) | Fixrunde: der Satz „der Guard scannt den Command-String“ nennt den Stand `diff` nicht mehr (Anführungszeichen-Lesung: ein Tool-Token in den Anführungszeichen einer Inline-Message blockt nicht); der Punkt und die Schritt-20-Zeile zum Konjunktiv-Kandidatenlauf (F-7) sind angepasst. |
| `.harness/skills/reviewer.md` §HIGH „Docker-only-Verstoß“ | nennt die Formen als Prüfgegenstand des Reviews (Zeile 3 des Feldes) | bleibt wahr: der Guard fängt die Flag-Formen, das Review den Rest; unberührt. |
| Beobachtungs-Register `BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel`, `state.md` | „Durchsetzung heute: das Review“ und „noch kein Plan angelegt“ (Zeile 4 des Feldes am Parent) | „noch kein Plan angelegt“ ist mit dem Anlegen dieses Plans auf „der Plan liegt in `open/`“ gezogen; „Durchsetzung heute“ bleibt bis zur Closure wahr und ist dann Träger der Planner-Closure (fremde Datei, deshalb Meldung). Am Stand `diff` liegt der Plan in `in-progress/` und der Guard blockt die Formen: beide Sätze sind ab dem Commit dieses Laufs überholt; gemeldet, Frist: die Closure. |
| `docs/plan/planning/welle-transformationen.md` §6 „Fragen für den nächsten Architect-Zug“ (a) | „in keinem committeten Text verboten“ (Zeile 5 des Feldes) — überholt seit dem Architect-Zug `17cb4eb3`, nicht erst durch diesen Slice | gemeldet, nicht mitgeändert (fremde Datei; Frist: die Closure dieses Slice, der Planner zieht nach oder benennt den Träger mit Adresse). |
| `docs/plan/planning/open/slice-code-kommentare-bereinigung.md` §8 („offen, Schwelle erreicht, 3×“) | die Sichtung nennt den Register-Eintrag „offen“, sein `state.md` sagt seit `17cb4eb3` „verkörpert“ (Zeile 6 des Feldes trifft drei Zeilen: diese, die Risiko-Zeile in §6, die nur die Zahl nennt und wahr bleibt, und die Notiz in `welle-transformationen`) | gemeldet wie oben. |

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
  geblockt werden müsste, obwohl er zulässig ist (die rohe Lesung der Flag-Tokens trägt ihn
  nicht), und der Weg dahin verlangt Quote-Bewusstsein im Guard — dann Architect-Frage
  (Shell-Parser ja oder nein; eingetreten im Review, Finding F-2: ein Muster mit `|` blockte;
  gelöst im Guard durch einen Zustandsautomaten über Anführungszeichen und Backslash, kein
  Shell-Parser, Fixrunde §3 — die Frage bleibt an den Architect als Bestätigung offen); (b) der Hook liest die Block-Ausgabe in der Live-Sitzung nicht wie
  im Bestand (der Beleg in Liefer-Punkt 1 fehlt) — dann Architect-Frage zur Wirkung der
  Durchsetzungsschicht.

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
  Blockmeldung nennt den Ersatzweg. **Ausgang:** *(bei Closure)*
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
  **Ausgang:** *(bei Closure)*
- **3. Die Änderung schwächt den Bestandsschutz oder macht den Guard fail-open.** Der Guard
  hat heute keinen Test (Suchlauf-Zeile 1: drei Fundstellen, keine Testdatei darunter).
  *Erwartet, zu belegen durch:* die Bestandsregeln stehen als Fälle im
  Tabellentest (Paketmanager am Kopf, Präfixe, `bash -c`-Rekursion, Tiefe über 3, defektes
  JSON, fehlendes `awk`) und eine Mutation an der neuen Erkennung färbt sie nicht rot.
  **Ausgang:** *(bei Closure)*
- **4. Die Wirkung in der Live-Sitzung bleibt unbelegt, weil der Tabellentest nur die
  Hook-Schnittstelle (JSON in, JSON aus) prüft.** Ob Claude Code die Ausgabe wie im Bestand
  liest, ist eine Eigenschaft der Laufzeit, kein Ergebnis des Skripts
  ([`AGENTS.md`](../../../../AGENTS.md) §3.10 nennt dieselbe Klasse für externe Läufe).
  *Erwartet, zu belegen durch:* der Live-Beleg in Liefer-Punkt 1 und 2 (ein Aufruf mit
  nicht vorhandenem Ziel wird geblockt, der Wortlaut steht im Bericht); der Verifier wiederholt
  ihn in seiner Sitzung. **Ausgang:** *(bei Closure)*
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
  im Register-Eintrag (Kenntnis, in diesem Slice nicht belegbar). **Ausgang:** *(bei Closure)*
- **6. Mutationsproben der Reviewer und Verifier werden behindert.** Sie arbeiten nach
  [`AGENTS.md`](../../../../AGENTS.md) §3.1 auf einer Kopie im Scratchpad; ein `sed -i` auf der
  Kopie wird jetzt geblockt. *Erwartet, zu belegen durch:* die Tabelle enthält den zulässigen
  Weg als Nicht-Treffer (`sed s/a/b/ Datei > /tmp/…/Kopie`, `python3 /tmp/…/mutate.py`), der
  Satz zur Mutationsprobe in §3.1 nennt ihn, und der Reviewer des Slice fährt eine Mutation
  auf diesem Weg. **Ausgang:** *(bei Closure)*
- **7. Host-Werkzeuge außerhalb von Docker und `make`** (`bash`, `awk`, `mktemp` im Test,
  `ls` im Guard; `BEO-PGC/host-werkzeug-jenseits-docker-und-make-ohne-deklaration`, 2×).
  *Erwartet, zu belegen durch:* alle vier stehen in der Klasse „Host-Werkzeug ohne
  Installation“ ([`AGENTS.md`](../../../../AGENTS.md) §3.1: `bash`, `git`, `awk`, `mktemp` und
  die coreutils-Basis); die README-Zeile nennt sie. **Ausgang:** *(bei Closure)*

**Offene Fragen an den Nutzer (Entscheidung, nicht Umfang) — je mit Adresse:**

1. **Die Kopf-Liste `tools/harness/blocked/go`** (`go gofmt python python3 node dotnet java
   gradle uv`; die Fragment-Ladung liegt im Guard, das Verzeichnis existiert nicht). Sie
   blockte unbedingt am Kopf eines Segments: sie schlösse die Toolchain-Aufrufe auf dem Host,
   die kein Vorfall belegt, und machte Liefer-Punkt 2 für `python` überflüssig; `docker run … go`
   ginge durch, weil `docker` der Kopf ist. Empfehlung: nicht in diesem Slice; zuerst die
   Wirkung von Liefer-Punkt 1 abwarten. Adresse: der Nutzer; bei „ja“ ein eigener Plan, den der
   Planner anlegt. Die Frage steht mit dem Closure-Stand im `state.md` des Register-Eintrags.
2. **Eine Scratchpad-Ausnahme für `sed -i`** (erlaubt, wenn ein Token auf einen absoluten
   Pfad unter `/tmp/` zeigt und keines auf einen Repo-Pfad). Der Plan blockt unbedingt: die
   drei Vorfälle klassifizieren auch das `sed -i` auf `/dev/null` und auf einer
   Scratch-Datei als Fehlgriff, und der zulässige Mutationsweg braucht kein `-i`
   (Stdout-Umleitung). Eine Ausnahme kostet Pfad-Logik und ein Loch (`xargs sed -i < /tmp/liste`
   ginge durch). Adresse: der Nutzer; bei „ja“ ändert der Implementer den Fall im Tabellentest,
   und `MR-003` nennt die Ausnahme.

## 7. Closure-Notiz

- **Was hat funktioniert:** *(zu tragen bei Closure)*
- **Was ging anders als geplant:** *(zu tragen bei Closure)*
- **Steering-Loop-Eintrag:** *(zu tragen bei Closure — erwartet: geschärfte Regel
  [`AGENTS.md`](../../../../AGENTS.md) §3.1 Durchsetzung und `MR-003` (Guard blockt die in-place
  Formen; Grenz-Zeile), neuer Sensor `make test-command-guard` ohne Gate; Auslöser
  `BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel` (drei Vorfälle:
  `slice-backfill-speicher-untersuchung`, `slice-transformationen-antragsweg-usecase`,
  `slice-transformationen-backfill-pfad`); das Feld `liegt in` steht nur, wenn wirklich
  verkörpert; ohne Eintrag kein `done/`-Übergang)*
- **Beobachtungs-Register (`../observations/`):** *(zu tragen bei Closure — `state.md` von
  `BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel`: die Durchsetzung heute, die Adresse
  aufgelöst, die zwei Nutzer-Fragen aus §6 mit Adresse; ein Auftreten trotz Guard wäre eine
  weitere `evidence/`-Datei)*
- **Folge-Slices:** keiner erwartet; die zwei Nutzer-Fragen aus §6 sind Entscheidungen mit
  Adresse „Nutzer“, kein Slice. Träger-Meldungen aus §3 mit Frist „Closure dieses Slice“:
  `welle-transformationen` §6 (a) und `slice-code-kommentare-bereinigung` §8 — der Planner
  zieht nach oder benennt den Träger mit Adresse.
- **Risiken aus §6:** *(je ein Ausgang, zu tragen bei Closure)*
- **Drei Paarungen:** dieser Slice hat keine Welle; die Slice-Closure selbst trägt die drei
  Paarungen (Anker · Folge-Slice · Register), nach dem `git mv` nach `done/`. Die
  Ereignis-Adresse kann eintreten (die Closure dieses Slice), und keine der Paarungen hängt an
  einem Ereignis außerhalb: *Anker* — [`AGENTS.md`](../../../../AGENTS.md) §3.1 Durchsetzung und
  `MR-003` tragen `seit slice-harness-guard-inplace-textwerkzeug`, der Kopfkommentar des Guards
  und das Ziel `test-command-guard` im `Makefile` existieren; *Folge-Slice* — keiner genannt; *Register* —
  `BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel` (und jede weitere genannte Kennung)
  existiert als Verzeichnis mit nicht leerem `evidence/`.

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
