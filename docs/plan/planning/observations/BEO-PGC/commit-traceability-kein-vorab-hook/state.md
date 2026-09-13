Zustand: offen — kein Träger bislang; eine Lösung wäre ein lokaler
`commit-msg`-Git-Hook (z. B. `.githooks/commit-msg`, per `git config
core.hooksPath` aktiviert, oder als Teil des Docker-only-Wrappers), der
dieselben zwei Prüfungen aus `tools/harness/commit-traceability.sh` und
dem d-check-Modul `commits` bereits vor dem `git commit`-Abschluss
anwendet, statt sie erst nachträglich über `make gates` zu melden.
Zähler (abgeleitet): 2× (evidence/slice-038.md, evidence/review-slice-041.md)
— unter der 3×-Schwelle, kein Ausgang fällig.
