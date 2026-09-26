# BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die Arbeitsweise aller
Agenten-Rollen am Repo, keine eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Die Nutzerregel „kein `sed -i`/`perl -pi`“ (Textänderungen an Repo-Dateien
laufen über die Edit-Werkzeuge des Laufs, nicht über ein in-place schreibendes Host-Werkzeug)
steht in keinem committeten Text dieses Repos — `AGENTS.md` §3.1 (Docker-only) nennt sie nicht,
kein Agenten-Skill und kein Befehl unter `.claude/` (gemessen mit
`git grep -n -i -E 'sed -i|perl -pi' -- .claude harness AGENTS.md .harness/skills`: kein
Treffer) — und wird über Läufe hinweg von jeder Rolle verletzt: Implementer, Reviewer und
Verifier melden je einen Fehlgriff, mal ohne Wirkung (Ziel `/dev/null` oder eine
Scratch-Datei), mal mit Wirkung auf eine Repo-Datei (die Spur ist eine
`gofmt`-Abweichung im Diff).

**Warum das zählt:** Ein in-place Textwerkzeug am Repo hinterlässt keine Spur, die ein
Sensor liest; die Wirkung zeigt sich nur, wenn der Diff sie zufällig trägt. Die Regel lebt
allein im Prompt des Orchestrators — ein Lauf ohne diesen Prompt kennt sie nicht.

**Abgrenzung.** `BEO-PGC/host-werkzeug-jenseits-docker-und-make-ohne-deklaration` betrifft die
Deklaration der Host-Werkzeuge, die jede Arbeit voraussetzt (`git`, `bash`); hier geht es um
ein schreibendes Werkzeug, das die Edit-Werkzeuge des Laufs umgeht.

Deklaration: `slice-backfill-speicher-untersuchung`, Review F-14 (INFO).
