Zustand: **verkörpert** — Ausgang: **verkörpert** → `AGENTS.md` §3.1, Absatz „Text-Umschreiben im
Repo ist Sache der Datei-Werkzeuge des Laufs“ (`sed -i`, `perl -pi`, `awk -i inplace`, Host-Interpreter
auf einer Repo-Datei sind verboten; Mutationsproben laufen auf einer Kopie im Scratchpad),
`.harness/skills/reviewer.md` (Unterpunkt „Docker-only-Verstoß“) und `.claude/commands/implement-slice.md`
(Verweis) · seit `AGENTS.md` §3.1 (Beleg-Anker: `git grep -n 'awk -i inplace' -- AGENTS.md`).
Durchsetzung heute: das Review — der PreToolUse-Guard blockt Paketmanager, kein in-place
Textwerkzeug. Adresse des Guard-Ausbaus: `slice-harness-guard-inplace-textwerkzeug` (Guard blockt
die Formen der Regel; noch kein Plan angelegt). Ein Sensor über Dateiinhalte ist ausgeschlossen
(ein Werkzeugaufruf hinterlässt keine Signatur in der Datei).
Zähler (abgeleitet): 3× (evidence/slice-backfill-speicher-untersuchung.md,
evidence/slice-transformationen-antragsweg-usecase.md,
evidence/slice-transformationen-backfill-pfad.md).
