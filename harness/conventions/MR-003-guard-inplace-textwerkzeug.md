# MR-003 — Der PreToolUse-Guard blockt in-place Textwerkzeuge und Host-Interpreter auf Repo-Pfaden

Regeln dieser Datei: Pflichtfelder sind Datum, Geltungsbereich,
**Ersetzt-Baseline-Regel**, Adaption, Begründung und Auflösungs-Trigger.

- **Datum:** 2026-09-26
- **Geltungsbereich:** `.claude/hooks/pretooluse-command-guard.sh`,
  `tools/harness/run-command-guard-tests.sh` (`make test-command-guard`),
  [`AGENTS.md`](../../AGENTS.md) §3.1 Absatz „Durchsetzung“.
- **Ersetzt-Baseline-Regel:** [`grundlagen-durchsetzungsschicht.md`
  §Grenzen — ehrlich
  benannt](../../.harness/baseline/v6.9.0/regelwerk/grundlagen-durchsetzungsschicht.md#grenzen--ehrlich-benannt)
  — der Satz, ein Befehls-Guard prüfe nur Befehlspositionen, Interpreter-Umwege
  blieben möglich. Der Guard dieses Repos liest darüber hinaus die Flag-Tokens
  dreier Werkzeuge und ein Pfad-Muster im Befehlsstring; die Grenze rückt, sie
  fällt nicht (Grenz-Zeile unten).
- **Auslöser:** Der Register-Eintrag
  `BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel` mit vier
  Beleg-Dateien (Zähler aus seinem `state.md`, übernommen): `sed -i` in fünf
  Läufen dreier Vorgänge (Implementer, Reviewer und Verifier; mit Wirkung auf
  Repo-Dateien einmal, sonst ohne), ein Host-Python-Heredoc auf einer Repo-Datei
  und zwei Host-`python3`-Aufrufe im vierten Vorgang (ein Implementer ohne
  Wirkung, ein Verifier mit unklarem Ziel). Die Regel steht in
  [`AGENTS.md`](../../AGENTS.md) §3.1; ihre Durchsetzung war das Review. `perl -pi` und
  `awk -i inplace` sind in keinem Beleg genannt: sie sind dieselbe Klasse und
  kosten je eine Erkennung.
- **Adaption:** Der Guard blockt (Ausgabe `"decision": "block"`, Exit 0, wie im
  Bestand) ein Kommando-Segment, dessen Kopf — nach den Zuweisungs- und
  Wrapper-Präfixen des Bestands, deren Optionen (`xargs -r`) und nach
  `-exec`/`-execdir`/`-ok`/`-okdir` — ist:
  - `sed` mit `--in-place[=…]` oder einem Bündel `-[nEsrzu]*i` (`-i`, `-i.bak`,
    `-ni`, `-Ei`);
  - `perl` mit einem Bündel `-[0-9lanpsw]*i` (`-i`, `-pi`, `-0777pi`);
  - `awk`/`gawk` mit `-i inplace`, `-iinplace`, `--include=inplace` oder
    `--include inplace`;
  - `python`, `python3`, `python3.<N>` oder `perl`, wenn der **ganze
    Befehlsstring** den Absolutpfad der Repo-Wurzel oder den Namen eines
    Eintrags der obersten Repo-Ebene nennt (Verzeichnis mit folgendem `/`, Datei
    als ganzes Wort), vor dem kein Pfadzeichen `[A-Za-z0-9_./~-]` steht (ein
    vorangestelltes `./` ist erlaubt).

  Die Flag-Tokens werden **roh** gelesen; ein Token mit Anführungszeichen ist kein
  Flag, damit ein Muster `'sed -i|perl -pi'` in einem `grep`-Argument nicht
  blockt. Der Block gilt unabhängig vom Ziel, auch auf einer Scratchpad-Kopie;
  der Mutationsweg ist `sed … Datei > Kopie` (stdout) oder Edit/Write. Die
  Blockmeldung nennt den Ersatzweg. Bestandsregeln (Paketmanager-Liste,
  `bash -c`-Rekursion mit Tiefe 3, fail-closed bei Parse-Zweifel und ohne `awk`)
  bleiben. Der Tabellentest `make test-command-guard` bindet jede Zusage an ihre
  Eingabe (Treffer, Nicht-Treffer neben jedem Treffer, benannte
  Falsch-Positiv-Ränder, benannte Grenzen); er ist ein Werkzeug, kein Gate.
- **Grenz-Zeile — was der Guard nicht kann.** Ein Stolperdraht, keine Sandbox:
  - Umleitungen und flaglose Schreibwege (`> datei`, `>>`, `tee`, `dd of=`,
    `sed … > tmp && mv tmp datei`, `cp`/`mv` über eine Datei);
  - ein `cd` im selben Kommando (der Pfad-Test nimmt die Repo-Wurzel als
    Arbeitsordner; `cd /tmp/x && python3 tools/x.py` blockt als Falsch-Positiv,
    `cd <Repo> && python3 x.py` geht durch);
  - Variablen, Globs, `$PWD`, `~` und Kommando-Ersetzungen als Pfad;
  - ein Skript, das der Interpreter liest (`python3 /tmp/x.py`, dessen Text
    Repo-Dateien schreibt), und `bash skript.sh`, dessen Inhalt `sed -i` trägt
    (nur `bash -c "…"` wird rekursiv gelesen);
  - jedes andere in-place-fähige Werkzeug (`ed`, `ex`, `patch`, `ruby -i`,
    `git apply`, `truncate`) und ein Flag in Anführungszeichen (`sed '-i'`);
  - Optionen mit Wert hinter einem Wrapper-Präfix (`sudo -u x sed -i`);
  - die Rezepte hinter `make` und die Docker-Bauten (der Guard scannt die
    Bash-Aufrufe des Laufs);
  - Falsch-Positive durch die Quote-Blindheit: ein Trenner in einem Argument oder
    einer Heredoc-Zeile startet ein neues Segment, und steht dort `sed -i` am
    Kopf, blockt der Guard (`grep -E 'a|sed -i ' f`); ein `python3`-Aufruf, dessen
    Text einen Repo-Namen der obersten Ebene nennt, blockt, auch wenn er nichts
    schreibt. Die Ränder stehen als benannte Fälle im Tabellentest.

  Die Grenz-Zeile steht gleichlautend im Kopfkommentar des Guards und in
  [`AGENTS.md`](../../AGENTS.md) §3.1 „Durchsetzung“; was der Guard nicht liest,
  bleibt Sache des Reviews (`.harness/skills/reviewer.md`, HIGH
  „Docker-only-Verstoß“). Ein Sensor über Dateiinhalte ist ausgeschlossen: ein
  Werkzeugaufruf hinterlässt in der Datei keine Signatur.
- **Begründung:** Die Härtung folgt der Beobachtung, nicht dem Bedrohungsmodell
  (Baseline-Regelwerk `modul-13-quality-gates.md` §Guard-Härtung): vier
  Beleg-Dateien derselben Klasse in vier Vorgängen unter dem Review als einziger
  Durchsetzung. Gehärtet wird die Erkennung der Form, nicht die Denylist um den
  Interpreter: `python`/`perl` blocken nur auf einem Repo-Pfad, damit Scratchpad-
  Mutationsläufe (`python3 /tmp/…/mutate.py /tmp/…/kopie`) durchgehen.
- **Auflösungs-Trigger:** permanent, bis ein Folge-`MR` den Guard schärft oder
  eine Sandbox-Ausführung ihn ersetzt; ein Auftreten trotz Guard ist eine weitere
  Beleg-Datei im Register-Eintrag, keine Änderung dieses Eintrags.
