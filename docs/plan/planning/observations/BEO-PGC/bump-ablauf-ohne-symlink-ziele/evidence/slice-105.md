**Vorgang:** slice-105 — Baseline `v6.9.0` materialisiert, das Verzeichnis des
alten Tags entfernt (`47d86b26`).

**Fund (später gefunden, in der Closure von `slice-harness-baseline-v6-14-1`
gemessen):** Am Move-Commit nach `done/` (`9b9e9b84`) liegt unter
`.harness/baseline/` nur `v6.9.0` (`git ls-tree 9b9e9b84 .harness/baseline/`),
während `git cat-file -p 9b9e9b84:.claude/rules/modul-05-planning-harness.md`
ein Ziel im Tag-Verzeichnis `v6.5.0` druckt
(`…/regelwerk/modul-05-planning-harness.md`). Die Regel-Symlinks hingen ab diesem Bump; weder Plan noch Review des
Vorgangs nennen sie.
