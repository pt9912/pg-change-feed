# MR-004 — Der PreToolUse-Guard sperrt Host-`python` und `python3` am Kopf eines Kommando-Segments unbedingt

Regeln dieser Datei: Pflichtfelder sind Datum, Geltungsbereich,
**Ersetzt-Baseline-Regel**, Adaption, Begründung und Auflösungs-Trigger;
`Löst auf` und `Ausgelöst durch Baseline-Stand` nur, wenn dieser Eintrag einen
früheren ablöst. `Ersetzt-Baseline-Regel` nennt **genau eine** Regel der
Baseline, an deren Stelle dieser Eintrag tritt — als Link mit
Abschnitts-Anker in die vendored Fassung; ein Datei-Link benennt keine Regel.

- **Datum:** 2026-09-27
- **Geltungsbereich:** das Fragment `tools/harness/blocked/python`, der
  PreToolUse-Guard (`.claude/hooks/pretooluse-command-guard.sh`), der
  Tabellentest (`tools/harness/run-command-guard-tests.sh`,
  `make test-command-guard`), [`AGENTS.md`](../../AGENTS.md) §3.1 Absatz
  „Durchsetzung“.
- **Ersetzt-Baseline-Regel:** [`modul-13-quality-gates.md`
  §Guard-Härtung: Wächter reifen in Wellen — „Gehärtet wird die Zerlegung,
  nicht die
  Denylist“](../../.harness/baseline/v6.14.0/regelwerk/modul-13-quality-gates.md#guard-haertung).
- **Adaption:** Der Guard sperrt zwei Kopf-Wörter unbedingt — `python` und
  `python3`, geladen aus `tools/harness/blocked/python` über die bestehende
  Fragment-Ladung (`BLOCKED`, Kopfkommentar des Guards). Ein Kommando-Segment,
  dessen Kopf nach den bestehenden Übersprüngen (Zuweisungen, Wrapper-Präfixe
  samt Optionen, `\`, Schlüsselwörter, Basename) `python` oder `python3` ist,
  blockt mit der Klasse `pkg`, unabhängig vom Rest des Befehlsstrings: auch
  hinter `cd <Repo> &&`, auch auf einem Pfad ohne Repo-Namen, auch bei
  `python3 --version`. Die Repo-Pfad-Regel von `MR-003` (Klasse `interp`)
  bleibt für `python3.<N>` und `perl` unverändert in Kraft — diese Namen
  erreichen den neuen Boden nicht, weil das Fragment nur die zwei exakten
  Wörter trägt. Die übrigen Namen der ursprünglich erwogenen Liste (`go`,
  `gofmt`, `node`, `dotnet`, `java`, `gradle`, `uv`) bleiben ungesperrt.
- **Begründung:** Die Baseline nennt als Schaden einer Denylist um den
  Interpreter, dass sie legitime Shell-Arbeit inklusive `make` blockiert
  (`modul-13-quality-gates.md` §Guard-Härtung, dritter Punkt). Die Kopf-Liste
  dieses Guards liest nur den **Kopf** eines Segments: `make` und `docker`
  bleiben frei, ebenso `command -v python3` und jedes `python3` als Argument,
  in Anführungszeichen oder als Teil eines Suchmusters (*hergeleitet* aus dem
  Guard-Code, *belegt* durch die Pass-Fälle des Tabellentests). Der Schaden,
  den die Baseline benennt, tritt hier nicht ein.

  Die zwei Ränder der Grenz-Zeile von `MR-003` — `cd <Repo> && python3 x.py`
  geht durch, ein Aufruf auf einem Pfad ohne Repo-Namen ebenso — schließt eine
  Zerlegung nicht: ein Pfad ohne Repo-Namen ist keine Zerlegungs-Frage (es gibt
  keine Sub-Shell zu entpacken), und ein `cd` im selben Kommando ist
  Shell-Zustand hinter der Obergrenze der Quote-Lesung (`MR-003`, Adaption
  „Obergrenze der Quote-Lesung“: der Maskierer kennt keine Heredocs, keine
  Umleitungen, keine Kommando-Werte — und auch kein `cd`). Die unbedingte
  Kopf-Liste schließt beide Ränder für `python`/`python3`, weil sie den
  Befehlsstring nach dem Kopf gar nicht mehr liest.

  Auslöser: das zweite Auflösungs-Kriterium von `MR-003` (die Closure des
  nächsten Slice, dessen Läufe unter diesem Guard liefen — eingetreten mit der
  Closure von `slice-transformationen-map-value`) und die Nutzer-Entscheidung
  „Weg 3“ in derselben Sitzung (2026-09-27): nicht die volle Liste, sondern nur
  die zwei Namen, für die eine Beleg-Datei vorliegt (Register-Eintrag
  `BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel`).
- **Grenz-Zeile — was der Guard nicht kann.** Ein Stolperdraht, keine Sandbox:
  - jeder nicht gelistete Name (`go`, `gofmt`, `node`, `dotnet`, `java`,
    `gradle`, `uv`, `python3.<N>`, `python2`, `pypy`, `ipython`,
    `uv run python`) — `python3.<N>` und `python2` bleiben unter der
    Repo-Pfad-Regel von `MR-003`, die übrigen Namen ganz ungelesen;
  - ein `python3`-Aufruf aus einer Variablen oder einem Alias (der Kopf ist
    dann nicht das Literal `python3`);
  - Umleitungen und flaglose Schreibwege, ein Skript, das ein anderer
    Interpreter liest (`ruby x.py`, `node x.py`) — dieselbe Grenze wie in
    `MR-003`;
  - die Rezepte hinter `make` und die Docker-Bauten (der Guard scannt die
    Bash-Aufrufe des Laufs, nicht was `python3` innerhalb eines Containers
    tut — `docker run … python3 …` bleibt frei, weil der Kopf `docker` ist);
  - eine Heredoc-Zeile, deren erstes Wort `python3` ist, blockt trotzdem (die
    Zeilen eines Heredocs gelten als Kommando-Zeilen, unverändert aus
    `MR-003`) — kein neuer Falsch-Negativ-Rand, aber ein neuer
    Falsch-Positiv-Rand: `python3 --version` blockt jetzt, ausdrücklich
    gewollt (Nutzer-Entscheidung).

  Schärft `MR-003` in genau diesen Sätzen: „`cd <Repo> && python3 x.py` geht
  durch“, „ein Skript, das der Interpreter liest (`python3 /tmp/x.py`) geht
  durch“ und der Falsch-Positiv-Rand „ein `python3`-Aufruf, dessen Text einen
  Repo-Namen nennt, blockt, auch wenn er nichts schreibt“ gelten für
  `python`/`python3` nicht mehr — sie bleiben wahr für `python3.<N>` und
  `perl`, die weiter über die Repo-Pfad-Regel laufen. `MR-003` selbst bleibt
  unverändert (`Accepted`, immutable).
- **Auflösungs-Trigger:** permanent; die Neubewertung der nicht gelisteten
  Namen (`go`, `gofmt`, `node`, `dotnet`, `java`, `gradle`, `uv`) folgt einer
  weiteren Beleg-Datei des Register-Eintrags
  `BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel`, die einen
  Host-Aufruf eines dieser Namen nennt; die Ergänzung der Liste bleibt eine
  Nutzer-Entscheidung.
