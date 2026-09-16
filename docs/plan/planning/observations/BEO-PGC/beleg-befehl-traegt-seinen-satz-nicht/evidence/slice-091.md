# Beleg: slice-091

Vorgang: `slice-091` — Coverage Cluster C (Zustell- und Betriebs-Rand).

Fund: **Zwei** Mutationen, die die Kommentare der neuen Tests als **rot färbend
benennen**, färbten nicht rot — gefunden vom Review (`review-slice-091` F-1, F-2)
durch Setzen genau der genannten Mutation.

- **F-1:** `notify_test.go` behauptete im Namen die Weitergabe des `LogPort`;
  gebunden war nur das **erste** Kettenglied. Die im Kommentar genannte Mutation
  (`notify.go:81` `log: o.log` → `outbound.NoopLog`) ließ Test **und Paket**
  grün (EC=0). Die Zusage „der Adapter führt seine Betriebszeilen über den
  injizierten Port" hatte damit keinen Träger am mittleren Glied.
- **F-2:** `interceptor_test.go` nannte als rot färbende Mutation „einen Wert
  statt `""`" — rot färbt erst ein **gültiger** Token (`"reader-token"`, EC=1);
  ein unbekannter Wert (`"probe"`) bleibt grün (EC=0), weil `classifyToken` ihn
  wie den leeren Token `roleNone` zuordnet.

**Ein dritter Fall, aus dem Delta-Review (D-1):** die **Behebung** von F-1/F-2
machte eine **Zählung** tragend, die ihre Gruppe nicht erzeugt —
`go list -f '{{len .TestGoFiles}}'` trifft **23** der 31 Pakete des Gegenstands
und ist damit keine Gruppierungsregel, während der Aufzählungspunkt „fünf" sagt.

**Ein Vorgang, eine Zählung.** Alle drei fallen in dieselbe Klasse und in
**denselben** Vorgang — der Zähler bewegt sich durch diesen Beleg **einmal**.
Ursprung und Vorkommen sind getrennt: die Form der beiden Testkommentare stammt
aus `slice-091` selbst (sie sind Teil dieses Diff), gefunden wurden sie in
`review-slice-091` und `review-slice-091-delta`; die zitierte Zählung ist
**älter** als dieser Vorgang (sie steht seit `slice-085`/`slice-079` im
Sensor-Dokument).

**Warum das zählt:** Ein Beleg ist die Prüf-Form einer Aussage. Hier stand der
Beleg **im Kommentar des Tests selbst** — wer ihn las, glaubte die Prüfung, statt
sie zu fahren. Die Mutation ist die Gegenprobe, und sie ist die einzige, die
diesen Fall aufdeckt.

Quelle: `docs/reviews/review-slice-091.md` (F-1, F-2) ·
`docs/reviews/review-slice-091-delta.md` (D-1) ·
`internal/adapters/driven/natsnotify/notify_test.go`,
`internal/adapters/driving/grpc/interceptor_test.go` (berichtigt in `f9cd5e4`) ·
`harness/sensors/coverage-gate.md` (berichtigt in `af3ea9f`).
