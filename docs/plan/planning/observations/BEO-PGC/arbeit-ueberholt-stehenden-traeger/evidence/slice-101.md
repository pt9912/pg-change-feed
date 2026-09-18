**Vorgang:** slice-101

**Fund:** Die Einführung des dritten Programms je Sprache
(`examples/csharp/nats-client`, `examples/kotlin/nats-client` — je eine
neue `runtime-nats`-Stufe neben den unveränderten `runtime`-/`runtime-sse`-
Stufen) machte zwei Sätze in `harness/README.md` §Sensors falsch, die nicht
im Diff standen: beide `make examples-csharp`/`make examples-kotlin`-Zeilen
beschrieben seit `slice-100` „zwei Images" je Sprache — bei Niederschrift
(`slice-100`) wahr, seit diesem Slice falsch, weil jetzt drei Images je
Sprache gebaut werden.

Der Implementer-eigene §3.13-Suchlauf fand beide Sätze und korrigierte sie
im selben Commit (`79dbd5d`) auf „drei Images" mit benannten Image-Tags
(`:csharp-nats`/`:kotlin-nats`). Reviewer und Verifier bestätigen die
Korrektur unabhängig als vollständig und nicht-überschießend
(Review zu `slice-101`, Negativbefund „`harness/README.md`" ·
Verifikationsbericht zu `slice-101`, #18). Wie bei `slice-093`/`slice-094`/
`slice-100` ist dies ein Fall, in dem der vorgeschriebene Suchlauf selbst
den Fund lieferte — das zählt als Auftreten der zugrunde liegenden
Beobachtung, nicht nur als Beleg dafür, dass die Regel wirkt.

**Nicht hierher gehört** ein zweiter, im selben Slice gefundener Defekt:
Der Slice-Kopf/§6/§8 dieses Plans zitierten `BEO-PGC/github-actions-
unverifizierbar-lokal` mit „5×" — real bereits 7× zum
Niederschrift-Zeitpunkt. Diese Zahl driftete nicht *durch* die Arbeit
dieses Slice (wie hier), sie war bereits *bei Niederschrift* veraltet —
das ist die Form von `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`,
eigene Evidenzdatei dort.

Quelle: Review zu `slice-101` (Negativbefund) ·
Verifikationsbericht zu `slice-101` (#18) · `AGENTS.md` §3.13.
