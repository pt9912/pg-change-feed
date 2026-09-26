# MR-003 — Der PreToolUse-Guard blockt in-place Textwerkzeuge und Host-Interpreter auf Repo-Pfaden

Regeln dieser Datei: Pflichtfelder sind Datum, Geltungsbereich,
**Ersetzt-Baseline-Regel**, Adaption, Begründung und Auflösungs-Trigger.

- **Datum:** 2026-09-26
- **Geltungsbereich:** `.claude/hooks/pretooluse-command-guard.sh`,
  `tools/harness/mask-quotes.awk`,
  `tools/harness/run-command-guard-tests.sh` (`make test-command-guard`),
  [`AGENTS.md`](../../AGENTS.md) §3.1 Absatz „Durchsetzung“.
- **Ersetzt-Baseline-Regel:** [`grundlagen-durchsetzungsschicht.md`
  §Grenzen — ehrlich
  benannt](../../.harness/baseline/v6.9.0/regelwerk/grundlagen-durchsetzungsschicht.md#grenzen--ehrlich-benannt)
  — der Satz, ein Befehls-Guard prüfe nur Befehlspositionen, Interpreter-Umwege
  blieben möglich. Der Guard dieses Repos liest darüber hinaus die Flag-Tokens
  dreier Werkzeuge und ein Pfad-Muster im Befehlsstring; die Grenze rückt, sie
  fällt nicht: der Satz gilt weiter für alles, was der Guard nicht liest
  (Grenz-Zeile unten).
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
- **Adaption:** Der Guard segmentiert den Befehlsstring **quote-bewusst**
  (`tools/harness/mask-quotes.awk`: ein Trenner, Leerraum oder Zeilenumbruch in
  einfachen oder doppelten Anführungszeichen und hinter einem Backslash trennt
  nicht; `$(` und Backtick in doppelten Anführungszeichen führen aus und trennen
  weiter; `\;` bleibt ein Trenner, weil es das Ende eines `find -exec` ist).
  Ein Anführungszeichen-Argument ist ein Token. Ein unbalanciertes
  Anführungszeichen ist Parse-Zweifel: der Guard segmentiert dann ohne
  Anführungszeichen-Kenntnis (mehr Segmente, nie weniger); kann der Maskierer
  nicht laufen, blockt der Guard (fail-closed).

  Er blockt (Ausgabe `"decision": "block"`, Exit 0, wie im Bestand) ein
  Kommando-Segment, dessen Kopf ist:
  - `sed`/`gsed` mit `--in-place[=…]` (auch jede eindeutige Abkürzung ab `--i`)
    oder einem Bündel `-[nEsrzub]*i` (`-i`, `-i.bak`, `-ni`, `-Ei`, `-bi`); das
    Bündel gilt auch hinter einem Dateinamen, weil GNU-sed Optionen vertauscht
    liest: `sed -n 1p -input.txt` ist `-i` mit dem Suffix `nput.txt`;
  - `perl` mit einem Bündel `-[0-7lanpsw]*i` (`-i`, `-pi`, `-i.bak`, `-0777pi`;
    `-MList::Util`, `-e`, `-Ii` sind keines); die Optionen enden am ersten
    Nicht-Options-Token (dem Skriptnamen), das Argument von `-e`/`-E`/`-I` wird
    übersprungen — `perl x.pl -input a` ist ein Skript-Argument;
  - `awk`/`gawk` mit `-i inplace`, `-iinplace`, `--include=inplace` oder
    `--include inplace`;
  - `python`, `python3`, `python3.<N>` oder `perl`, wenn der **ganze
    Befehlsstring** den Absolutpfad der Repo-Wurzel oder den Namen eines
    Eintrags der obersten Repo-Ebene nennt (Verzeichnis mit folgendem `/`, Datei
    als ganzes Wort), vor dem kein Pfadzeichen `[A-Za-z0-9_./~-]` steht (ein
    vorangestelltes `./` ist erlaubt).

  Der **Kopf** ist das erste Wort des Segments nach diesen Übersprüngen:
  Zuweisungen (`VAR=…`); Wrapper-Präfixe (`sudo env command exec nice time xargs
  eval busybox`, auch mit absolutem Pfad) samt ihren Optionen (`xargs -r`,
  `env -i`, `time -p`) und dem Wert einer xargs-Option (`xargs -n 1`); ein
  führendes `\` (`\sed`); die Schlüsselwörter `do then else elif if while until
  !`, ein führendes `{` oder `)` (Funktionsdefinition) und, nach einem `case`
  im selben Befehl, die Case-Labels (`x) sed -i …`). `command -v` und
  `command -V` zeigen an und führen nichts aus; der Rest hinter `eval`,
  `bash -c` (auch `-lc`, `-ec`) und hinter `-exec`/`-execdir`/`-ok`/`-okdir`
  ist ein eigenes Kommando (Rekursionstiefe 3, darüber fail-closed).
  Anführungszeichen gehören nicht zum Flag (`sed -i''` und `sed '-i'` blocken); ein
  Anführungszeichen-Argument mit Leerraum ist ein Token und kein Flag
  (`grep -E 'sed -i|perl -pi'` blockt nicht). Diese Übersprünge gelten für alle
  Klassen, auch für die Paketmanager: `env -i pip`, `xargs -n1 pip`, `time -p pip`
  und `for … do pip …` blocken. Das erweitert die Bestandsregel (der Bestand las
  nach einem Wrapper-Präfix keine Optionen und kannte keine Schlüsselwörter); der
  Tabellentest führt die Formen als eigene Gruppe.

  Der Block gilt unabhängig vom Ziel, auch auf einer Scratchpad-Kopie. Der Weg
  nach [`AGENTS.md`](../../AGENTS.md) §3.1 ist Edit/Write, `sed … Datei >
  Kopie` (stdout) oder ein Repo-Werkzeug hinter `make`; beide Blockmeldungen
  nennen diesen Ersatzweg. Bestandsregeln (Paketmanager-Liste,
  `bash -c`-Rekursion mit Tiefe 3, fail-closed bei Parse-Zweifel und ohne `awk`)
  bleiben. Der Tabellentest `make test-command-guard` bindet die Zusagen an ihre
  Eingabe (Treffer, Nicht-Treffer neben jedem Treffer, ein Fall je Mitglied der
  Zeichenklassen, benannte Falsch-Positiv-Ränder, benannte Grenzen); er ist ein
  Werkzeug, kein Gate.
- **Grenz-Zeile — was der Guard nicht kann.** Ein Stolperdraht, keine Sandbox:
  - Umleitungen und flaglose Schreibwege (`> datei`, `>>`, `tee`, `dd of=`,
    `sed … > tmp && mv tmp datei`, `cp`/`mv` über eine Datei);
  - ein `cd` im selben Kommando (der Pfad-Test nimmt die Repo-Wurzel als
    Arbeitsordner; `cd /tmp/x && python3 tools/x.py` blockt als Falsch-Positiv,
    `cd <Repo> && python3 x.py` geht durch);
  - Variablen, Globs, `$PWD`, `~` und Kommando-Ersetzungen als Pfad, ein Werkzeug
    aus einer Variablen (`$x -i`), `eval "$cmd"` und Aliase;
  - ein Skript, das der Interpreter liest (`python3 /tmp/x.py`, dessen Text
    Repo-Dateien schreibt; der Aufruf auf einem Pfad ohne Repo-Namen geht durch,
    ist aber kein Weg nach `AGENTS.md` §3.1), und `bash skript.sh`, dessen
    Inhalt `sed -i` trägt (nur `bash -c "…"` wird rekursiv gelesen);
  - jedes andere in-place-fähige Werkzeug (`ed`, `ex`, `patch`, `ruby -i`,
    `git apply`, `truncate`) und jeder andere Interpreter (`node -e`, `ruby -e`,
    `uv run python`);
  - ein perl-Bündel mit Buchstaben außerhalb der Klasse (`-Wpi`) und die
    Abkürzung einer awk-Langoption (`--inc=inplace`);
  - Optionen mit Wert hinter einem Wrapper außer `xargs` (`sudo -u x sed -i`,
    `env -u X`, `nice -n 10`);
  - die Rezepte hinter `make` und die Docker-Bauten (der Guard scannt die
    Bash-Aufrufe des Laufs);
  - Falsch-Positive: die Zeilen eines Heredocs gelten als Kommando-Zeilen (`sed -i`
    am Zeilenanfang blockt); ein unbalanciertes Anführungszeichen, auch ein
    Apostroph im Heredoc-Text, segmentiert ohne Anführungszeichen-Kenntnis und
    blockt ein `sed -i` hinter einem Trenner im Text; `echo a\;sed -i x` blockt
    (`\;` trennt); ein `python3`-Aufruf, dessen Text einen Repo-Namen der obersten
    Ebene nennt, blockt, auch wenn er nichts schreibt; ein Case-Label gilt nur
    nach einem `case` im selben Befehl. Die Ränder stehen als benannte Fälle im
    Tabellentest.

  Die Grenz-Zeile steht in Kurzform im Kopfkommentar des Guards und in
  [`AGENTS.md`](../../AGENTS.md) §3.1 „Durchsetzung“ (deren „er liest nicht“-Liste
  ist eine Teilmenge dieser); was der Guard nicht liest, bleibt Sache des Reviews
  (`.harness/skills/reviewer.md`, HIGH „Docker-only-Verstoß“). Ein Sensor über
  Dateiinhalte ist ausgeschlossen: ein Werkzeugaufruf hinterlässt in der Datei
  keine Signatur.
- **Begründung:** Die Härtung folgt der Beobachtung, nicht dem Bedrohungsmodell
  (Baseline-Regelwerk `modul-13-quality-gates.md` §Guard-Härtung): vier
  Beleg-Dateien derselben Klasse in vier Vorgängen unter dem Review als einziger
  Durchsetzung. Gehärtet wird die Erkennung der Form, nicht die Denylist um den
  Interpreter: `python`/`perl` blocken nur auf einem Repo-Pfad (zwei
  Beleg-Dateien; die Heuristik trägt die größte Falsch-Positiv-Fläche und ist der
  abtrennbare Teil). Ein Muster mit `|` in einem `git grep`/`grep` ist der
  Alltag der Rollen (das Suchmuster des Suchlaufs nach `AGENTS.md` §3.13): die
  Quote-Bewusstsein-Stufe (ein Zustandsautomat über Anführungszeichen und
  Backslash, kein Shell-Parser) hält diesen Aufruf frei. Ohne Beleg im Register
  sind `perl -i` und `awk -i inplace` (dieselbe Klasse, je eine Erkennung, am
  Tabellentest gebunden).
- **Auflösungs-Trigger:** permanent, bis ein Folge-`MR` den Guard schärft oder
  eine Sandbox-Ausführung ihn ersetzt; ein Auftreten trotz Guard ist eine weitere
  Beleg-Datei im Register-Eintrag, keine Änderung dieses Eintrags.
