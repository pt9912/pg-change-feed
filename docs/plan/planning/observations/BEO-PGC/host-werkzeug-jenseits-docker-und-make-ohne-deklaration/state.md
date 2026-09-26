Zustand: **verkörpert** — Ausgang: **verkörpert** → `AGENTS.md` §3.1, Absatz „Host-Werkzeug
ohne Installation“ (die eine Klasse neben Docker und `make`: `bash`, `git`, POSIX-/coreutils-Basis;
Werkzeuge jenseits davon laufen im Container; Sensor-/Target-Verträge nennen, was über `bash` und
`git` hinausgeht). Die Klasse beschreibt den Bestand und lockert kein Gate — keine ADR nötig · seit
`AGENTS.md` §3.1 (Beleg-Anker: `git grep -n 'Host-Werkzeug ohne Installation' -- AGENTS.md`).
Trigger für ein weiteres Auftreten: ein Skript, das ein Host-Werkzeug außerhalb der Klasse aufruft
(Reviewer-HIGH „Docker-only-Verstoß“, `.harness/skills/reviewer.md`); ein Plan, dessen §6 die Frage
erneut als offenes Risiko führt, ist ein Zeichen, dass der Absatz sie nicht beantwortet.
Geprüft an `slice-harness-guard-inplace-textwerkzeug`: der Plan führte die Frage als Risiko 7, und der
Absatz beantwortet sie (die Host-Werkzeuge des Tabellentests — `bash`, `awk`, `mktemp`, `grep`, `cp`,
`mkdir`, `ln` — und `ls` im Guard stehen in der Klasse; der Kopf des Tabellentests und die Zeile zu
`make test-command-guard` in `harness/README.md` nennen sie): der Trigger ist nicht eingetreten,
kein neuer Beleg.
Zähler (abgeleitet): 2× (evidence/slice-harness-suchlauf-nachmessen.md,
evidence/slice-code-kommentare-kennungen.md).
