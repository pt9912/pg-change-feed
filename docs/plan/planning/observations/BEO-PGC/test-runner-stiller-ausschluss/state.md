Zustand: geplant (Ausgang unten) — Ausgangslage: kein Wächter prüfte, dass
jede Testfunktion in `test/integration/integration_test.go` von
mindestens einem `-run`-Muster in `run-integration-tests.sh` erfasst
wird. Zähler (abgeleitet): **2×** (evidence/slice-033.md,
evidence/slice-074.md, evidence/slice-routing-nats-subjekt.md) — **3×, Schwelle
erreicht**; der Ausgang gehört zum Lese-Schritt der Closure von `welle-routing`. Der
zweite Träger ist die **Deklarations-Hälfte** der E2E-Abdeckungstabelle (`slice-074`): eine
neue Runner-Phase, die niemand deklariert, fehlt still in der Tabelle. Der dritte Träger
(`slice-routing-nats-subjekt`, Review F-3, LOW) ist `tools/harness/run-notify-tests.sh` mit
einem `-run`-Namenspräfix: dort ist die Lücke geschlossen (das Skript fährt das ganze Paket und
endet bei `--- SKIP` mit Exit 1, vom Verifier in drei Zuständen ausgeführt). Für
`run-integration-tests.sh` bleibt sie offen.

**Ausgang (Architect-Verdikt zur Closure von `welle-routing`, 2026-10-02): geplant → `slice-harness-integration-runner-vollstaendigkeit`** (Datei in `open/`). Ein Vollständigkeits-Sensor ohne Regeländerung: jede `func TestE2E*` steht in einem `-run`-Argument von `tools/harness/run-integration-tests.sh`; gelesen werden die `-run`-Werte, nicht Kommentare. Heute ist die Lücke leer (21 Funktionen, 0 ohne Namenstreffer im Skript, am Stand `cff48b65` gemessen). Die Wirksamkeit des Wächters ist *hergeleitet*, nicht erprobt (Verdikt `architect-verdict-welle-routing-lese-schritt` §3.2). Der Eintrag wird mit dem Slice `verkörpert` (Zielort: der Test).

Lese-Schritt der Closure von `welle-routing` (2026-10-02): gelesen. Für `tools/harness/run-notify-tests.sh` ist die Lücke geschlossen (ganzes Paket, Exit 1 bei `--- SKIP`); für `tools/harness/run-integration-tests.sh` (`go test -run '^(…)$'`-Muster, ab Zeile 466) bleibt sie: die Deklarations-Hälfte deckt `TestAbdeckungstabelleZeilen`, die Vollständigkeits-Hälfte (jede `func TestE2E*` wird von einem Muster oder einer Phase erfasst) trägt kein Wächter. Die Entscheidung, ob ein Vollständigkeits-Sensor gebaut wird, steht beim Architect (ein Sensor wäre ein eigener Slice); bis dahin ist der Trigger das nächste Auftreten bei `run-integration-tests.sh`. Adresse: `welle-routing-results.md`, Lese-Schritt.
