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
Kopf-Liste `tools/harness/blocked/go`: **entschieden** — der Nutzer hat „Weg 3“ gewählt (Sitzung
2026-09-27; Frage und Wege: `done/slice-transformationen-map-value.md` §7, Punkt „Frage an den Nutzer“):
nicht die volle Liste (`go gofmt python python3 node dotnet java gradle uv`), sondern nur `python python3`;
die übrigen Namen bleiben ungesperrt, weil für sie kein Beleg vorliegt. Ausgang dieses Teils: **verkörpert** →
`slice-harness-guard-blocked-python` (Fragment `tools/harness/blocked/python` — der Name folgt der Konvention
`blocked/<sprache>` des Guards —, `MR-004`, Tabellentest `make test-command-guard` mit 368 Fällen). Der Guard
blockt `python`/`python3` am Kopf jetzt unbedingt; `python3.<N>`/`perl` bleiben unter der Repo-Pfad-Regel.
Ausgelöst hat die Neubewertung das zweite Trigger-Kriterium von `MR-003` (Auflösungs-Trigger: die Closure
des nächsten Slice, dessen Läufe unter diesem Guard liefen; eingetreten mit der Closure von
`slice-transformationen-map-value`); das erste (eine Beleg-Datei mit einem Host-Interpreter-Aufruf ohne
Repo-Pfad **mit Wirkung** auf eine Repo-Datei) ist nicht eingetreten. Die Neubewertung der nicht gelisteten
Namen trägt der Auflösungs-Trigger von `MR-004`. Die Scratchpad-Ausnahme für
`sed -i` ist entschieden: keine, der Guard blockt unbedingt.
Zähler (abgeleitet): 8× (evidence/slice-backfill-speicher-untersuchung.md,
evidence/slice-transformationen-antragsweg-usecase.md,
evidence/slice-transformationen-backfill-pfad.md,
evidence/slice-antragsqueue-lesefehler-failed.md,
evidence/slice-transformationen-map-value.md,
evidence/adr-0129-capture-quellseite-keepalive-test-an-beiden-pins.md,
evidence/slice-leerlauf-phase-last-in-stuecken.md,
evidence/slice-harness-guard-blocked-python.md). Der siebte Beleg
(`slice-leerlauf-phase-last-in-stuecken`) trägt einen Host-Aufruf `python3 -c 1` des Implementers
ohne im Bericht genannten Guard-Ausgang und ohne Wirkung auf eine Repo-Datei im Diff —
**unvollständig belegt**: weder ein Beleg noch ein Gegenbeleg für das erste
Neubewertungs-Kriterium der Kopf-Liste. Der achte Beleg
(`slice-harness-guard-blocked-python`, derselbe Vorgang wie die Kopf-Liste selbst) trägt eine
Umleitung (`>`), die eine neue Repo-Datei anlegte, statt Edit/Write zu nutzen — vom Implementer
selbst im selben Werkzeugaufruf-Batch bemerkt und zurückgenommen, bevor sie einen Commit erreichte.
Der fünfte Beleg ist eine **Regelgrenze**, kein neuer
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
`MR-003`); ein Beleg für das erste Neubewertungs-Kriterium der Kopf-Liste ist er nicht:
der Guard hat gegriffen, der Aufruf hatte keine Wirkung auf eine Repo-Datei.
Der sechste Beleg trägt zwei Aufrufe, beide ohne Wirkung auf eine Repo-Datei: ein Architect-Aufruf
`python3 --version` ohne Repo-Pfad passierte den Guard, ein Reviewer-Aufruf `python3 --version` im selben
Vorgang wurde geblockt. Ursprung beider: aus Agenten-Berichten übernommen (Angabe des Auftraggebers), nicht am
Guard gemessen; gemessen ist allein, dass die Hook-Eingabe `python3 --version` am Guard des Standes
`cea198fb` mit Exit 0 ohne Ausgabe endet. **Nicht gezählt** (kein abgeschlossener Vorgang mit Beleg-Datei-Namen;
Ursprung wie oben, keine Wirkung auf eine Repo-Datei): ein Planner-Aufruf `python3` mit leerem Heredoc bei
einer Slice-Anlage (geblockt; der Vorgang ist im Auftrag nicht benannt).
Der neunte Beleg (`slice-routing-sdk-beispiel-target`, evidence/slice-routing-sdk-beispiel-target.md)
ist eine Mutationsprobe des Reviewers, die nach einem fehlgeschlagenen `cd` an drei Repo-Dateien statt an
der Kopie im Scratchpad lief; der Reviewer bemerkte es an `git status`, nahm sie mit `git checkout` zurück
und prüfte danach `git status --short` leer. Die Fassung `AGENTS.md` §3.1 („eine Mutationsprobe arbeitet
auf einer Kopie“) galt; was fehlte, war die Absicherung des Verzeichniswechsels. Beobachtung, keine neue
Regel: eine Mutationsreihe nennt absolute Pfade und prüft `git status` nach jedem Lauf. Kein Beleg für
das erste Neubewertungs-Kriterium (kein Host-Interpreter am Kopf).
