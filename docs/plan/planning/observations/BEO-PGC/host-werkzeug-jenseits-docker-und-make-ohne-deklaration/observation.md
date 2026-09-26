
# BEO-PGC/host-werkzeug-jenseits-docker-und-make-ohne-deklaration

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die Hard Rule `AGENTS.md`
§3.1 und die Skripte unter `tools/`, keine eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: `AGENTS.md` §3.1 sagt „Host braucht nur Docker und GNU `make`“. Skripte unter
`tools/harness/` und `tools/schema/` rufen Host-Werkzeuge jenseits davon auf, ohne dass die Regel
sie nennt: `git rev-parse --show-toplevel` steht in 42 Dateien unter `tools/` (gemessen am Stand
`142ca4b5` mit `git grep -c 'git rev-parse --show-toplevel' -- tools`: 42 Dateien, 43 Zeilen),
das Nachmess-Werkzeug braucht zusätzlich `bash` und `realpath`. Die Praxis ist Präzedenz
(„es wird nichts installiert, `git` ist Voraussetzung jeder Arbeit an diesem Repo“); die
Regel führt keine Ausnahme-Klasse „Host-Werkzeuge, die auf jedem Entwicklerrechner des Repos
liegen“. Reviewer (F-6) und Verifier (V-3) nannten die Vereinbarkeit unabhängig als
Architect-Aussage, nicht als Befund.

**Warum das zählt:** Ein Slice, dessen Plan „Docker-only für `git`“ als Risiko führt, kann
das Risiko nur mit Präzedenz auflösen; die Regel selbst gibt ihm kein Kriterium, welche
Host-Werkzeuge zulässig sind.

**Abgrenzung.** Die Beobachtung betrifft die **Deklaration** der Klasse, nicht das
Werkzeug: kein Skript installiert etwas auf dem Host.

Deklaration: `slice-harness-suchlauf-nachmessen`.
