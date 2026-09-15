Zustand: offen — Ausgang noch nicht zugewiesen. Ein Träger wäre ein
`make`-Ziel, das den Generator laufen lässt und `git diff --exit-code` gegen
das committete Erzeugnis prüft (für den Protobuf-Code: `make proto-generate`
plus Diff auf `streamv1/`); das wäre ein eigener Sensor mit eigener Bindung
und eigenem Aufwand pro `make gates` — deshalb keine spontane Ergänzung,
sondern eine Entscheidung. Zähler (abgeleitet): **2×** (evidence/slice-069.md,
evidence/slice-074.md) — der zweite Träger ist das committete Erzeugnis
`docs/user/e2e-abdeckung.md` (`slice-074`), das kein Sensor ohne den vollen
E2E-Lauf gegen seine Quelle hält; unter der 3×-Schwelle.
