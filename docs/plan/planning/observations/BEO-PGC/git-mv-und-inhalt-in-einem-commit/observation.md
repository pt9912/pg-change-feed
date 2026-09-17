# BEO-PGC/git-mv-und-inhalt-in-einem-commit

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die
Commit-Disziplin bei Umzügen erzeugter oder verschobener Dateien, keine
eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Ein Commit verschiebt Dateien (Git erkennt die Umbenennung,
hohe Similarity) **und** ändert im selben Commit ihren Inhalt sowie weitere
Dateien — statt, wie `AGENTS.md` §3.3 verlangt, den Move-Commit rein zu
halten und die Inhaltsänderung als eigenen Folge-Commit zu führen. Anders als
bei der in §3.3 selbst benannten Lifecycle-`done/`-Ausnahme (dort ist die
Reihenfolge *Inhalt vor* `git mv` vorgeschrieben, mit eigener Begründung)
geht es hier um den **Regelfall**: Datei verschoben und Inhalt umgeschrieben
sind zwei Commits, unabhängig von der Reihenfolge. Die Regel selbst ist
bereits eine Hard Rule (`AGENTS.md` §3.3) und bereits verkörpert — dieser
Eintrag zählt **Verstöße gegen eine bestehende Regel**, nicht Wiederholungen,
aus denen erst eine Regel entstehen soll. Das ist ein anderer
Beobachtungstyp als die 3×-Schwelle für neue Regeln.

Deklaration: `slice-097`, Review-Report F-1 (`docs/reviews/review-slice-097.md`)
und Verifikationsbericht (`docs/reviews/verify-slice-097.md` §4) — beide
bestätigen den Verstoß unabhängig voneinander, beide stufen ihn als
nicht-merge-blockierend und nicht-fix-pflichtig ein (Git erkannte die
Umbenennung trotzdem, `git log --follow` bleibt funktionsfähig).
