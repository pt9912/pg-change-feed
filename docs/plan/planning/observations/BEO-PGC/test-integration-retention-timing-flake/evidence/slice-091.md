# Beleg: slice-091

Vorgang: `slice-091` — Coverage Cluster C (Zustell- und Betriebs-Rand).

Fund: Über denselben Vorgang sind **zwei** Läufe desselben Commits unterschiedlich
ausgefallen. Der Implementer berichtete **77,30 %**, mein eigener Lauf und das
Gate druckten **77,20 %**; die Verifikation sah über 47 Läufe ein Band von
**77,1–77,3 %**. Die Differenz zwischen den beiden genannten Enden ist
**0,105 pp** — und **2 von 1903** Statements sind genau 0,105 pp: der
Takt-Zweig von `runWALRetentionCheck` (`internal/bootstrap/wiring.go:991.5,992.13`),
den [`ADR-0082`](../../../../../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)
§Kontext (2) als Zähler-Schwankung ±2 führt.

**Derselbe Gegenstand, dritter Träger.** `evidence/slice-057.md` dateit einen
`make test-integration`-Lauf, der real mit Exit 2 abbrach und bei Wiederholung
grün durchlief; `evidence/slice-090.md` dieselbe Erscheinung an der
Coverage-Zahl. Hier ist der Träger der **Vergleich zweier Läufe desselben
Commits** — die Beobachtung ist dieselbe: *ein zeitabhängiger Test liefert bei
unverändertem Stand wechselnde Ergebnisse.*

**Was dieser Vorgang zusätzlich gemessen hat:** Der Flapper ist **nicht eine
Stelle**. Der Delta-Review und die vierte Runde haben über 24 bzw. 8 Läufe
**zwei** unruhige Blöcke isoliert — den Takt-Zweig von `runWALRetentionCheck`
(2 Statements) und den Kontext-Ende-Zweig von `runAdministration`
(`:1091.4,1092.1`, 1 Statement). Der zweite erklärt das dritte Ende, das die
Verifikation zunächst nicht lokalisieren konnte: ein Lauf mit 435 ungedeckten
Statements, eins unter dem dokumentierten Band. Die **alte** Zuordnung im
Sensor-Dokument (`grpcstream/broadcaster.go` `Publish`) ist widerlegt — sie steht
in **8 von 8** Läufen auf `100,0 %`.

**Warum das zählt:** Die Zahl, an der `welle-20` ihr Ziel misst, ist **selbst
nicht vollständig stabil** — und das Band ist breiter als die dokumentierte
Einzelquelle. Solange `THRESHOLD = 70` und rund 7 pp Abstand bleiben, entscheidet
das nichts; bei der Endstufe 80 rückt die gemessene Zahl an die Schwelle. Die
Verifikation hat dafür **47 Läufe** gefahren — eine Arbeit, die keine DoD
verlangt und kein Gate leisten kann.

Quelle: Verifikationsbericht zu `slice-091` (V-1) ·
Delta-Review zu `slice-091` (D-2, Negativbefunde) ·
`harness/sensors/coverage-gate.md` §Zählbasis (berichtigt in `a7d7f7b`) ·
`internal/bootstrap/wiring.go:991.5,992.13` und `:1091.4,1092.1`.
