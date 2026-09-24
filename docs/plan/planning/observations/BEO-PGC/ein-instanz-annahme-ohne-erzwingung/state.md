Zustand: offen — Ausgang: **weiter offen**, adressiert. Die Annahme steht im Godoc von
`NewBackfillAdmission` und im Plan von `slice-backfill-run-store` (§6); Adresse ist der
Re-Evaluierungs-Trigger in `ADR-0113` §Re-Evaluierungs-Trigger („Mehr als eine Instanz
je Quelle nimmt Anträge an"). Der sequenzielle Fall trägt der DoD „Verarbeitung" in
`slice-backfill-sql-administration` (zwei aufeinanderfolgende `backfill`-Anträge derselben
Tabelle, der zweite endet `failed`); der nebenläufige Fall bleibt ungetestet, eine Sperre
oder eine Unique-Kante auf aktive Runs ist Sache der Re-Evaluierung, nicht dieses Eintrags.
Zähler (abgeleitet): 1× (evidence/slice-backfill-run-store.md).
