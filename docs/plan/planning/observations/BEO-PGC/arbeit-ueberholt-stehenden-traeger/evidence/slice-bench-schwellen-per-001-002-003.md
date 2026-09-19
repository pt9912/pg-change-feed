# Beleg: slice-bench-schwellen-per-001-002-003

Vorgang: `slice-bench-schwellen-per-001-002-003` — `LH-QA-PER-001`/`002`/`003`
über `docs/user/bench-abdeckung.md` (neu) als `trace.coverage`-Dimension für
`make doc-trace` sichtbar gemacht; die RTM-Waisenzahl bewegt sich dadurch von
21 auf 18 (mit `trace.coverage`).

Fund: Die Arbeit bewegt dieselbe Eigenschaft, die
`harness/README.md`s `make doc-trace`-Zeile beschreibt — die real gemessene
Waisenzahl —, und überholt damit die dort stehende Zahl. Der Implementer war
in derselben Datei bereits aktiv (eine Zeile darunter, `make bench`-Zeile
aktualisiert) und übersah die `make doc-trace`-Zeile trotzdem: keine
Enumerations-Lücke über mehrere Dateien wie in früheren Belegen dieses
Registers, sondern eine **Nachbarzeilen-Lücke** in derselben Datei — der
eigene §3.13-Suchlauf griff nicht, weil er auf Symbolnamen/Pfade zielt und
keine Zahlenwerte zuverlässig trifft (dieselbe Grenze, die `AGENTS.md` §3.13
selbst für Zahlen benennt). Nicht vom Implementer, sondern vom unabhängigen
Reviewer gefunden (`review-slice-bench-schwellen-per-001-002-003.md` F-1) —
über ein reales `make doc-trace`-Nachmessen, nicht über Diff-Lesen allein.
Derselbe Reviewer-Lauf fand zwei weitere, verwandte Fälle in
`tools/bench-scaling.sh`/`tools/bench-batch-vs-single.sh`: deren
Kopfkommentare behaupteten weiterhin „kein Pass/Fail", obwohl derselbe Diff
echtes Pass/Fail (`exit 1` gegen `THRESHOLD_*`) einführte — dieselbe
Fehlerklasse, diesmal am Kommentar statt an einem Doku-Träger.

Quelle: `docs/reviews/review-slice-bench-schwellen-per-001-002-003.md` <!-- d-check:status-provenance -->
(Findings F-1/F-2/F-3) · Fixrunden-Commit zu
`docs/plan/planning/in-progress/bench-schwellen-per-001-002-003.md`.
