Zustand: offen — der Ausgang wird im Lese-Schritt dieser Closure zugewiesen
(Architect-Zug, Modul 6 Schritt 3b). Ein Träger wäre ein `make`-Ziel, das den
Generator laufen lässt und `git diff --exit-code` gegen das committete Erzeugnis
prüft (für den Protobuf-Code: `make proto-generate` plus Diff auf
`streamv1/`); das wäre ein eigener Sensor mit eigener Bindung und eigenem Aufwand
pro `make gates` — deshalb keine spontane Ergänzung, sondern eine Entscheidung.
Zähler (abgeleitet): **3×** (evidence/slice-069.md, evidence/slice-074.md,
evidence/slice-082.md) — **Schwelle erreicht**. Die drei Träger: der
Protobuf-Code (`slice-069`), die E2E-Abdeckungstabelle (`slice-074`) und
`tools/schema/plan.yaml` (`slice-082`, dessen Rollout bei **jedem** Lauf die
committete Datei verändert).
