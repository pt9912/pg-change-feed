**Vorgang:** slice-100

**Fund:** Die Einführung des zweiten Runtime-Images je Sprache
(`examples/csharp/sse-client`, `examples/kotlin/sse-client` — je eine neue
`runtime-sse`-Stufe neben der unveränderten `runtime`-Stufe) machte zwei
Sätze in `harness/README.md` §Sensors falsch, die nicht im Diff standen:
beide `make examples-csharp`/`make examples-kotlin`-Zeilen beschrieben „das
Werkzeugketten-Image" (Singular) je Sprache — bei Niederschrift
(`slice-098`/`-099`) wahr, seit diesem Slice falsch, weil jetzt zwei Images
je Sprache gebaut werden.

Der Implementer-eigene §3.13-Suchlauf fand beide Sätze vollständig und
korrigierte sie im selben Commit (`7484fc3`) auf Plural mit benannten
Image-Tags (`:csharp-sse`/`:kotlin-sse`) und erweitertem Herkunfts-Anker
(„· seit slice-098, erweitert seit slice-100" bzw. „· seit slice-099,
erweitert seit slice-100"). Reviewer und Verifier bestätigen die Korrektur
unabhängig als vollständig und nicht-überschießend. Anders als bei
`slice-097` (Fund erst durch einen zweiten, unabhängigen Suchlauf) und wie
bei `slice-093`/`slice-094` ist dies ein Fall, in dem der vorgeschriebene
Suchlauf selbst den Fund lieferte — das zählt als Auftreten der zugrunde
liegenden Beobachtung (ein stehender Träger wurde real überholt), nicht nur
als Beleg dafür, dass die Regel wirkt: Gezählt wird, wie oft Arbeit einen
stehenden Träger überholt, nicht, wer den Fund machte.

Quelle: `docs/reviews/review-slice-100.md` (Negativbefund „`harness/
README.md`-Korrektur") · `docs/reviews/verify-slice-100.md` (#15) ·
`AGENTS.md` §3.13.
