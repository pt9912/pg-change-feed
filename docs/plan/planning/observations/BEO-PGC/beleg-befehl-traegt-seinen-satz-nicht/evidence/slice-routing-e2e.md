**Vorgang:** slice-routing-e2e (Review F-1, MEDIUM)

**Fund:** Der Satz „Ein auf ein Ziel gewählter Stream-Leser sieht nur dieses Ziel" stand in der Abdeckungszeile und im DoD. Die sechs Stream-Clients mit Ziel liefen mit `-count 1` und endeten nach der ersten empfangenen Change; der Runner prüfte nur diese eine Zeile. Getragen war „der Filter weist die vorangestellte fremde Change ab", nicht „keine fremde Change danach" — die Mutation, die den Filter erst hinter der ersten Treffer-Change lecken lässt, blieb ungeprüft. Die Fixrunde (`3df6db30`) ließ alle neun Clients mit einem Ruhefenster (`-window 15s`) weiterzählen, setzte eine feste Menge gemischter Changes (Region NULL, asia, us, eu, zweimal) hinter den Empfang der ersten Treffer-Change und hielt jede empfangene Zeile gegen die persistierte Change derselben `change_id`; die Zahl der Changes des eigenen Ziels aus der Menge ist gleich der SQL-Zählung. Der Verifier mutierte den `grpcclient` so, dass der erste Treffer stimmt und alles danach ungefiltert ist: rot. Die Grenze (15 s Ruhe, Menge NULL/asia/us/eu, nur `grpcclient` mutiert) steht als weiter offen im Plan.

**Form (Ausprägung):** Assertion — „nur" ist breiter als die Messung „erste Zeile". Schwere MEDIUM; vor dem Merge vom Reviewer gefunden, vom Verifier an Runner und Mutation geschlossen.

Quelle: `docs/reviews/review-slice-routing-e2e.md` (F-1) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-routing-e2e.md` (§4 M1, M2, §5 Zeile F-1). <!-- d-check:status-provenance -->
