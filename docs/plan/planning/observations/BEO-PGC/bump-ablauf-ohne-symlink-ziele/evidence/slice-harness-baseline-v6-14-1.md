**Vorgang:** slice-harness-baseline-v6-14-1 — Baseline `v6.14.1` vendored,
`v6.14.0` entfernt (`e2666499`).

**Fund:** Am Stand `4045dc4f` (v6.14.0 schon entfernt) zeigen vier der sieben
Regel-Symlinks ins Leere (Blob-Scan über `git ls-tree -r`: 4 Ziele mit
`/v6.14.0/`). Der Suchlauf des Plans sah sie nicht, der Haken von
Liefer-Punkt 2 stand auf einem Beleg, der diese Träger nicht liest (Review
`review-slice-harness-baseline-v6-14-1` F-1, HIGH; F-4, LOW, für den
Bump-Ablauf). Den Gegenstand fand der Auftraggeber; umgestellt in `b6c5b419`,
die Prüfung steht seit `fca136d4` in `harness/targets/pin-stale.md`
§Bump-Ablauf Schritt 1.
