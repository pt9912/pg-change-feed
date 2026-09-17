Zustand: **verkörpert** — Ausgang: **verkörpert** → `make generated-sync`
(harness/mk/generated-sync.mk, tools/harness/generated-sync.sh,
harness/sensors/generated-sync.md), in `make gates` · seit slice-090. Der
**Lese-Schritt der `welle-20`-Closure** hat den Ausgang zugewiesen; in seinem
Lauf ist das Gate grün (das committete Erzeugnis ist byte-gleich). Ein Träger wäre ein `make`-Ziel, das den
Generator laufen lässt und `git diff --exit-code` gegen das committete Erzeugnis
prüft (für den Protobuf-Code: `make proto-generate` plus Diff auf
`streamv1/`); das wäre ein eigener Sensor mit eigener Bindung und eigenem Aufwand
pro `make gates` — deshalb keine spontane Ergänzung, sondern eine Entscheidung.
Zähler (abgeleitet): **4×** (evidence/slice-069.md, evidence/slice-074.md,
evidence/slice-082.md, evidence/slice-086.md) — **Schwelle erreicht**. Die drei
Träger: der Protobuf-Code (`slice-069`), die E2E-Abdeckungstabelle (`slice-074`)
und `tools/schema/plan.yaml` (`slice-082`, dessen Rollout bei **jedem** Lauf die
committete Datei verändert). Der vierte Beleg ist ein weiterer Vorgang derselben
Klasse und kein neuer Handlungsbedarf: `slice-086` **erzeugt** die
Abdeckungstabelle neu (Teil seiner DoD) und nimmt `plan.yaml` wieder von Hand
zurück — die Bindung ist weiterhin Disziplin, kein Sensor.
