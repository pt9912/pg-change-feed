Zustand: offen — der Ausgang wird im Lese-Schritt der **laufenden
Welle-Closure** (`welle-20`) zugewiesen; dort trägt ihn der Planner → Architect
→ Planner-Zug (Modul 6, „Bei 3× … laufende Welle-Closure"; Modul 8 Schritt 3b).
Bis dahin ist `offen` der zulässige, vorübergehende Stand. Ein Träger wäre ein `make`-Ziel, das den
Generator laufen lässt und `git diff --exit-code` gegen das committete Erzeugnis
prüft (für den Protobuf-Code: `make proto-generate` plus Diff auf
`streamv1/`); das wäre ein eigener Sensor mit eigener Bindung und eigenem Aufwand
pro `make gates` — deshalb keine spontane Ergänzung, sondern eine Entscheidung.
Zähler (abgeleitet): **3×** (evidence/slice-069.md, evidence/slice-074.md,
evidence/slice-082.md) — **Schwelle erreicht**. Die drei Träger: der
Protobuf-Code (`slice-069`), die E2E-Abdeckungstabelle (`slice-074`) und
`tools/schema/plan.yaml` (`slice-082`, dessen Rollout bei **jedem** Lauf die
committete Datei verändert).
