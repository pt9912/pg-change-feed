# Beleg: slice-release-image-scan

Der Fixrunden-Commit `1a8c0efb` ("fix(release): Reviewer-Fixrunde für
release-image-scan (2 MEDIUM, 2 LOW)") trug keine `LH-*`-/`ADR-*`-Kennung
— gefunden vom Verifier (Verifikationsbericht zu `release-image-scan`, §6)
über einen realen, ungepipten `make gates`-Lauf (Exit 2,
`commit-traceability` rot), nicht vom Review, das gegen den *vorherigen*
Commit lief und die durch den Fixrunden-Commit selbst entstehende
Verletzung strukturell nicht sehen konnte.

**Fünfter Beleg für denselben Mechanismus:** `core.hooksPath` war in
dieser Arbeitskopie nicht gesetzt (`git config --get core.hooksPath` →
leer) — der lokale, opt-in `commit-msg`-Hook (`.githooks/commit-msg`,
[`ADR-0062`](../../../../../adr/0062-lokaler-commit-msg-hook-ergaenzt-standing-gate.md))
hätte den Verstoß vor dem Commit-Abschluss gemeldet, war aber nicht
aktiviert. Trifft exakt die in `state.md` bereits benannte schwächste
Stelle der Verkörperung — kein neuer Mechanismus, derselbe.

**Behoben** über das etablierte Muster: `git reset --soft
ea631152` + Wiederherstellung beider Datei-Stände (`1141786c`,
`1a8c0efb`) aus dem Objekt-Store + Neu-Commit mit korrigierter,
`ADR-0051`-tragender Message (`eb76267f`) — kein `git commit --amend`,
kein `git rebase -i` (unpushed, lokale Historie). Erneuter, ungepipter
`make gates`-Lauf danach: `EXIT=0`.
