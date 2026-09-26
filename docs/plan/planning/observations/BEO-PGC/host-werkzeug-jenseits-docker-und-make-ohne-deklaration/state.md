Zustand: **verkörpert** — Ausgang: **verkörpert** → `AGENTS.md` §3.1, Absatz „Host-Werkzeug
ohne Installation“ (die eine Klasse neben Docker und `make`: `bash`, `git`, POSIX-/coreutils-Basis;
Werkzeuge jenseits davon laufen im Container; Sensor-/Target-Verträge nennen, was über `bash` und
`git` hinausgeht). Die Klasse beschreibt den Bestand und lockert kein Gate — keine ADR nötig · seit
`AGENTS.md` §3.1 (Beleg-Anker: `git grep -n 'Host-Werkzeug ohne Installation' -- AGENTS.md`).
Trigger für ein weiteres Auftreten: ein Skript, das ein Host-Werkzeug außerhalb der Klasse aufruft
(Reviewer-HIGH „Docker-only-Verstoß“, `.harness/skills/reviewer.md`); ein Plan, dessen §6 die Frage
erneut als offenes Risiko führt, ist ein Zeichen, dass der Absatz sie nicht beantwortet.
Zähler (abgeleitet): 2× (evidence/slice-harness-suchlauf-nachmessen.md,
evidence/slice-code-kommentare-kennungen.md).
