# BEO-PGC/bump-ablauf-ohne-symlink-ziele

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft den
Bump-Ablauf der vendored Baseline, keine eigene Sub-Area im Sinn der
Modus-Deklaration).

Die Beobachtung: Ein Baseline-Bump entfernt das Verzeichnis des alten Tags,
die getrackten Symlinks unter `.claude/rules/` (Modus 120000) zeigen aber
weiter in dieses Verzeichnis und hängen danach. Kein Werkzeug des Ablaufs sieht
sie: `git grep` liest kein Symlink-Ziel (`git grep -c 'v6\.14\.0' 4045dc4f --
.claude/rules` endet mit Exit 1, obwohl vier Ziele dort den Tag nennen), und
`make docs-check` bleibt bei hängenden Symlinks grün (Review
`review-slice-harness-baseline-v6-14-1` F-4, an einem Klon von `4045dc4f`
gemessen). Der Suchlauf nach `AGENTS.md` §3.13 deckt diese Träger deshalb
nicht, obwohl sein Suchraum der ganze Baum ist.

Benannt, nicht gezählt: `2d51b380` stellte die vier Symlinks von `v6.5.0` auf
`v6.14.0` um — ein einzelner Commit außerhalb eines Slice, in keinem Plan
genannt.
