Zustand: **verkörpert** — Ausgang: **verkörpert** → `AGENTS.md` §3.1, Absatz „Text-Umschreiben im
Repo ist Sache der Datei-Werkzeuge des Laufs“ (`sed -i`, `perl -pi`, `awk -i inplace`, Host-Interpreter
auf einer Repo-Datei sind verboten; Mutationsproben laufen auf einer Kopie im Scratchpad),
`.harness/skills/reviewer.md` (Unterpunkt „Docker-only-Verstoß“) und `.claude/commands/implement-slice.md`
(Verweis) · seit `AGENTS.md` §3.1 (Beleg-Anker: `git grep -n 'awk -i inplace' -- AGENTS.md`).
Durchsetzung heute: das Review — der PreToolUse-Guard blockt Paketmanager, kein in-place
Textwerkzeug. Adresse des Guard-Ausbaus: `slice-harness-guard-inplace-textwerkzeug` (Guard blockt
die Formen der Regel; der Plan liegt in `open/`). Ein Sensor über Dateiinhalte ist ausgeschlossen
(ein Werkzeugaufruf hinterlässt keine Signatur in der Datei).
Zähler (abgeleitet): 4× (evidence/slice-backfill-speicher-untersuchung.md,
evidence/slice-transformationen-antragsweg-usecase.md,
evidence/slice-transformationen-backfill-pfad.md,
evidence/slice-antragsqueue-lesefehler-failed.md). Der vierte Beleg trifft zwei Rollen: ein
Host-`python3`-Aufruf des Implementers ohne Wirkung und die Mutationsläufe des Verifiers mit einem
Host-`python3`-Skript, dessen Ziel der Bericht widersprüchlich nennt (Arbeitskopie und Kopie); die
Durchsetzung ist unverändert das Review, der Guard-Ausbau liegt in `slice-harness-guard-inplace-textwerkzeug`.
