**Vorgang:** slice-harness-baseline-v6-13-0 — Baseline `v6.13.0` vendored,
`v6.9.0` entfernt (`88cea828`).

**Fund (später gefunden, in der Closure von `slice-harness-baseline-v6-14-1`
gemessen):** Am Commit, der den Plan nach `done/` bringt (`3c9d9483`), liegt
unter `.harness/baseline/` nur `v6.13.0`; das Ziel von
`.claude/rules/modul-05-planning-harness.md` liegt weiter im
Tag-Verzeichnis `v6.5.0`. Der
Bump ließ die seit `slice-105` hängenden Symlinks unbemerkt stehen.
