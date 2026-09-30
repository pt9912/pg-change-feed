**Vorgang:** slice-capture-transient-wiederholung (Review F-3, MEDIUM; Fixrunde 2, N-5, LOW; Verifikation §4 und §5 F-8)

**Fund:** Die Grenzen der Wiederholung (Fenster, Anfangsverzögerung, Obergrenze) waren nicht an ihre Werte gebunden (F-3, in der Fixrunde gebunden). Der Wert der Aufbau-Frist im Aufruf in `Run` blieb unter Mutation grün (N-5, vom Verifier unabhängig reproduziert; Folge-Slice `slice-capture-retry-aufbau-frist-bindung`). Der Realtest „Slot noch aktiv" erzwingt den Fehlschlag nur indirekt und prüft weder SQLSTATE 55006 noch `count == 2` (F-8, LOW, nach der Deckel-Regel ohne eigene Datei; Folge-Slice `slice-capture-retry-realtest-belege-schaerfen`). Ausgang unverändert **verkörpert**; die Regel hat mit Reviewer und Verifier gewirkt.

Quelle: `docs/reviews/review-slice-capture-transient-wiederholung.md` (F-3, N-5) · `docs/reviews/verifikation-slice-capture-transient-wiederholung.md` (§4, §5). <!-- d-check:status-provenance -->
