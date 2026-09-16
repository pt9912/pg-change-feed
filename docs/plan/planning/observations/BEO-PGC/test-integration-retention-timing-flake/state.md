Zustand: offen — Ausgang: **weiter offen**. Zähler (abgeleitet): **3×**
(evidence/slice-057.md, evidence/slice-090.md, evidence/slice-091.md) —
**Schwelle erreicht**; den Ausgang weist der **Lese-Schritt der
`welle-20`-Closure** zu (Modul 6), nicht die Slice-Closure.

Die drei Belege treffen **denselben Gegenstand über verschiedene Träger**:
`slice-057` den `make test-integration`-Exit, `slice-090` die **Coverage-Zahl**
desselben Test-Objekts (`WAL-Retention`, `runWALRetentionCheck`) — dort brach ein
Lauf real ab, hier liefert `go tool cover -func` in 7 von 8 Läufen 87,5 % und in
einem 100,0 %, bei durchweg grünen Tests. Beide Male wechselt das Ergebnis bei
**unverändertem Stand**.
