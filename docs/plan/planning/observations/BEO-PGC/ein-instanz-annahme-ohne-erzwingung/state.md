Zustand: offen — Ausgang: **weiter offen**, adressiert. Die Annahme steht im Godoc von
`NewBackfillAdmission` und im Plan von `slice-backfill-run-store` (§6); Adresse ist der
Re-Evaluierungs-Trigger in `ADR-0113` §Re-Evaluierungs-Trigger („Mehr als eine Instanz
je Quelle nimmt Anträge an"). Der sequenzielle Fall trägt der DoD „Verarbeitung" in
`slice-backfill-sql-administration` (zwei aufeinanderfolgende `backfill`-Anträge derselben
Tabelle, der zweite endet `failed`; belegt im Login-Test `TestAdministrationPathRunsUnderLeastPrivilegeLogins`,
Verifikation §2 Zeile 2); der nebenläufige Fall bleibt ungetestet, eine Sperre
oder eine Unique-Kante auf aktive Runs ist Sache der Re-Evaluierung, nicht dieses Eintrags.
Der zweite Beleg (`slice-backfill-sql-administration`, F-10) betrifft die andere Hälfte der
Annahme, „je Quelle nimmt eine Instanz Anträge an": für die Antragsart `backfill` bindet seit
der Fixrunde der Quellvergleich in `processAdministrationRequests` sie am Code (Test mit zwei
Quellen, Mutation rot); für die vier übrigen Antragsarten liest die Queue weiter ohne
Quellbezug (Bestand aus `ADR-0050`), und zwei Instanzen derselben Quelle bleiben ohne Schutz.
Zähler (abgeleitet): 2× (evidence/slice-backfill-run-store.md,
evidence/slice-backfill-sql-administration.md).
