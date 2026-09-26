**Vorgang:** slice-transformationen-map-value (Review F-6, INFO; Verifikation V-3, INFO)

**Fund:** Die Kosten der Zuordnung in `lookupMappedValue`
(`internal/domain/model/transformation.go`) stehen im Doc-Kommentar und im Plan als „linear in
der Zahl der Paare“ — benannt, hergeleitet aus dem Quelltext. Für den Betreiber blieb die
Größenordnung ohne Empfänger: die Spec bindet die Zahl der Paare von `values` nicht nach oben,
und `slice-transformationen-betriebsdoku` nannte weder „linear“ noch „Paare“
(Suche der Verifikation `git grep -n -i linear -- docs/plan/planning/open`: kein Treffer, **übernommen**).
Der Verifier maß sie (Wegwerf-Benchmark in einer Kopie, Schlüssel am Ende der Zuordnung, ohne
`-race`, **übernommen**): 80 ns bei 10 Paaren, 6,8 µs bei 1000 und 0,72 ms bei 100 000 je Wert
und Regel.

**Form (Ausprägung):** keine Ungenauigkeit in einer Datei, sondern eine **benannte Grenze
ohne Adresse**: der Träger (Kommentar, Plan) nennt die Grenze richtig, aber keine Stelle nimmt
die Betreiber-Aussage an. Die Closure gab sie mit Adresse ab: ein Übergabe-Block in §2 von
`slice-transformationen-betriebsdoku` und die Spec-Frage nach einer Obergrenze in
`welle-transformationen` §5 (Fragen für den nächsten Architect-Zug, Punkt (d)).

Quelle: `docs/reviews/review-slice-transformationen-map-value.md` (F-6) <!-- d-check:status-provenance -->
· `docs/reviews/verify-slice-transformationen-map-value.md` (V-3, §5 Zeile F-6). <!-- d-check:status-provenance -->
