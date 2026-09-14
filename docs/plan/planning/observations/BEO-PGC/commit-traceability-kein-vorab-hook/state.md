Zustand: offen — Schwelle (3×) mit `evidence/slice-059.md` erreicht.
Zähler (abgeleitet): 3× (evidence/slice-038.md, evidence/review-slice-041.md,
evidence/slice-059.md). Ausgang noch nicht zugewiesen: `slice-059` trägt
`Welle: welle-16` (nicht wellenlos) — der Lese-Schritt, der den Ausgang
zuweist, läuft bei der nächsten Welle-Closure (`welle-16`), nicht bei
dieser Slice-Closure (Baseline-Regelwerk `modul-08-agentenrollen.md`
§Rollen-Sequenz für eine Welle). Vorschlag für den Lese-Schritt: ein
lokaler `commit-msg`-Git-Hook (z. B. `.githooks/commit-msg`, per `git
config core.hooksPath` aktiviert, oder als Teil des Docker-only-Wrappers),
der dieselben zwei Prüfungen aus `tools/harness/commit-traceability.sh`
und dem d-check-Modul `commits` bereits vor dem `git commit`-Abschluss
anwendet, statt sie erst nachträglich über `make gates` zu melden.
