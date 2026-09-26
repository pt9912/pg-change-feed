Zustand: offen — **Schwelle erreicht** (3×; der Lese-Schritt der Closure von
`welle-transformationen` liest den Eintrag und weist den Ausgang zu). Ausgang-Vorschlag: die
Regel bekommt einen committeten Träger — `AGENTS.md` §3.1 nennt neben „kein lokaler Install“
das in-place schreibende Host-Textwerkzeug (`sed -i`, `perl -pi`, ein Host-Interpreter auf der
Repo-Datei) als Verstoß, und die Agenten-Prompts unter `.claude/agents/` verweisen darauf;
Mutationsläufe von Reviewer und Verifier gehören auf eine Kopie im Scratchpad, die Rücknahme
ist `cp` bzw. `git checkout`. Adresse: der nächste Architect-Zug, der `AGENTS.md` §3.1
berührt, gemeinsam mit der Ausnahme-Klasse aus
`BEO-PGC/host-werkzeug-jenseits-docker-und-make-ohne-deklaration`. Ein Sensor ist
ausgeschlossen: ein Werkzeugaufruf hinterlässt in einer Datei keine Signatur, die `docs-check`
oder ein anderes Gate liest.

Zähler (abgeleitet): 3× (evidence/slice-backfill-speicher-untersuchung.md,
evidence/slice-transformationen-antragsweg-usecase.md,
evidence/slice-transformationen-backfill-pfad.md). Der dritte Beleg trägt alle drei Rollen
in einem Vorgang, jeweils ohne Wirkung auf das Repo.
