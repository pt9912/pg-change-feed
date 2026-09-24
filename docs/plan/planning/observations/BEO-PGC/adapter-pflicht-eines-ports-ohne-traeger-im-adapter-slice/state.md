Zustand: offen — Ausgang noch nicht zugewiesen (unter der 3×-Schwelle). Die konkrete
Manifestation ist getragen: die beiden Pflichten stehen als Port-Vertrag in den
Doc-Kommentaren von `BackfillRunPort`, `BackfillTransaction.Rollback` und
`TableSnapshot.Close` und als DoD-Punkt „Adapter-Pflichten" im Plan von
`slice-backfill-run-store`; der Beleg im Adapter folgt mit dessen Store-Test. Bei 3× wäre zu
prüfen, ob die Slice-Planung für einen Use Case, der eine Pflicht an einen späteren Adapter
stellt, den DoD-Punkt im Adapter-Plan **im selben Zug** verlangen soll.

Zähler (abgeleitet): 1× (evidence/slice-backfill-run-usecase.md).
