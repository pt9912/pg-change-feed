Zustand: offen — Ausgang: **weiter offen** → kein Wächter prüft, dass
jede Testfunktion in `test/integration/integration_test.go` von
mindestens einem `-run`-Muster in `run-integration-tests.sh` erfasst
wird; wird beobachtet, bis ein Sensor/eine Vollständigkeitsprüfung dafür
gebaut wird. Zähler (abgeleitet): **2×** (evidence/slice-033.md,
evidence/slice-074.md, evidence/slice-routing-nats-subjekt.md) — **3×, Schwelle
erreicht**; der Ausgang gehört zum Lese-Schritt der Closure von `welle-routing`. Der
zweite Träger ist die **Deklarations-Hälfte** der E2E-Abdeckungstabelle (`slice-074`): eine
neue Runner-Phase, die niemand deklariert, fehlt still in der Tabelle. Der dritte Träger
(`slice-routing-nats-subjekt`, Review F-3, LOW) ist `tools/harness/run-notify-tests.sh` mit
einem `-run`-Namenspräfix: dort ist die Lücke geschlossen (das Skript fährt das ganze Paket und
endet bei `--- SKIP` mit Exit 1, vom Verifier in drei Zuständen ausgeführt). Für
`run-integration-tests.sh` bleibt sie offen.
