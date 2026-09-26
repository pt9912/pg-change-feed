Zustand: **verkörpert** — Ausgang: **verkörpert** → `AGENTS.md` §3.1, Absatz „Text-Umschreiben im
Repo ist Sache der Datei-Werkzeuge des Laufs“ (`sed -i`, `perl -pi`, `awk -i inplace`, Host-Interpreter
auf einer Repo-Datei sind verboten, ebenso das Schreiben und Anhängen von Text an eine Repo-Datei per
Umleitung — `>`, `>>`, `tee`, Heredoc; Mutationsproben laufen auf einer Kopie im Scratchpad), dem
PreToolUse-Guard (`.claude/hooks/pretooluse-command-guard.sh`: blockt `sed -i`/`--in-place`, `perl -i`,
`awk -i inplace` unabhängig vom Ziel und Host-`python`/`perl` mit einem Repo-Pfad im Befehlsstring;
Vertrag und Grenz-Zeile: `MR-003`, Tabellentest `make test-command-guard`, kein Gate),
`.harness/skills/reviewer.md` (Unterpunkt „Docker-only-Verstoß“) und `.claude/commands/implement-slice.md`
(Verweis) · seit slice-harness-guard-inplace-textwerkzeug, Umleitung geschärft seit welle-transformationen
(Architect-Verdikt `architect-verdict-welle-transformationen-offene-fragen` §5; Beleg-Anker:
`git grep -n 'Umleitung' -- AGENTS.md .harness/skills/reviewer.md .claude/commands/implement-slice.md`
und `git grep -n 'rewrite repo files without a trace' -- .claude`).
Was der Guard nicht liest, bleibt Sache des Reviews: Umleitungen und flaglose Schreibwege, ein Skript,
das ein Interpreter liest, Host-`python`/`perl` auf einem Pfad ohne Repo-Namen, andere
in-place-fähige Werkzeuge (Grenz-Zeile in `MR-003`). Ein Sensor über Dateiinhalte ist ausgeschlossen
(ein Werkzeugaufruf hinterlässt keine Signatur in der Datei); ein Guard-Ausbau gegen `> Repo-Pfad` ist
verworfen (er blockte `make gates > /tmp/log` und ließe `tee` und Variablen offen, Verdikt §5).
Offene Nutzer-Entscheidung, Adresse: die Kopf-Liste `tools/harness/blocked/go` (`go gofmt python
python3 node dotnet java gradle uv`; das Verzeichnis existiert nicht, die Fragment-Ladung liegt im
Guard). Trigger der Neubewertung, wörtlich wie in `MR-003` (Auflösungs-Trigger): eine weitere Beleg-Datei
dieses Eintrags mit einem Host-Interpreter-Aufruf ohne Repo-Pfad im Befehlsstring und Wirkung auf eine
Repo-Datei, oder die Closure des nächsten Slice, dessen Läufe unter diesem Guard liefen. Das zweite
Kriterium ist mit der Closure von `slice-transformationen-map-value` eingetreten, das erste nicht
(Beleg-Datei `evidence/slice-transformationen-map-value.md`: keine der drei Stellen ist ein
Host-Interpreter-Aufruf ohne Repo-Pfad mit Wirkung auf eine Repo-Datei). Der Planner dieser Closure hat
die Frage mit der gemessenen Wirkung des Guards schriftlich vorgelegt: Closure-Notiz
`done/slice-transformationen-map-value.md` §7, Punkt „Frage an den Nutzer“. Die Entscheidung des Nutzers
steht aus. Die Scratchpad-Ausnahme für
`sed -i` ist entschieden: keine, der Guard blockt unbedingt.
Zähler (abgeleitet): 5× (evidence/slice-backfill-speicher-untersuchung.md,
evidence/slice-transformationen-antragsweg-usecase.md,
evidence/slice-transformationen-backfill-pfad.md,
evidence/slice-antragsqueue-lesefehler-failed.md,
evidence/slice-transformationen-map-value.md). Der fünfte Beleg ist eine **Regelgrenze**, kein neuer
Fehlgriff: ein Anhängen per `cat >>` (der Guard liest Umleitungen nicht; die Regelfrage ist mit dem
Architect-Verdikt entschieden: `AGENTS.md` §3.1 verbietet den Weg),
ein vom Guard geblockter `sed -i` des Verifiers und ein `cd <Repo> && python3 --version` des Planners, das
den Guard an der in `MR-003` benannten Grenze passierte. Der vierte Beleg trifft drei Rollen: ein
Host-`python3`-Aufruf des Implementers ohne Wirkung, ein gleichartiger des Planners der
Closure-Sitzung (Angabe des Auftraggebers) und die Mutationsläufe des Verifiers mit einem
Host-`python3`-Skript, dessen Ziel der Bericht widersprüchlich nennt (Arbeitskopie und Kopie). Aus
`slice-harness-guard-inplace-textwerkzeug` selbst fällt kein Beleg an: die Mutationsläufe von Reviewer
und Verifier liefen mit `sed` ohne `-i` und `awk` nach stdout an Kopien. Wirkung des Guards in der
Closure-Sitzung gemessen: ein versehentlicher `python3 --version`-Aufruf mit einem Repo-Pfad im
Befehlsstring wurde mit der Meldung der Klasse `interp` geblockt und lief nicht (Falsch-Positiv-Rand,
`MR-003`); ein Beleg für den Neubewertungs-Trigger der Kopf-Liste `tools/harness/blocked/go` ist er nicht:
der Guard hat gegriffen, der Aufruf hatte keine Wirkung auf eine Repo-Datei.
