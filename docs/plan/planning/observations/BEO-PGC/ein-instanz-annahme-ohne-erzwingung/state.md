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
Zähler (abgeleitet, gemessen mit `ls evidence | wc -l`): **3×**
(evidence/slice-backfill-run-store.md,
evidence/slice-backfill-sql-administration.md,
evidence/slice-start-vorlauf-grenze.md). Der dritte Beleg
(`slice-start-vorlauf-grenze`) trifft dieselbe Annahme aus einer neuen
Richtung: die Verschiebung von `START_REPLICATION` nach `Stream.Run`
(`ADR-0128`) verschiebt den Zeitpunkt, an dem eine zweite Instanz derselben
Quelle am Slot scheitert, nach hinten (Vorlauf plus Goroutine-Start statt
sofort beim Verbindungsaufbau); am realen Container beim Debuggen jenes
Slice trat einmalig ein verwandter, aber nicht identischer Fall auf
(„context canceled“ statt SQLSTATE 55006 bei zwei überlappenden
Vorläufen) — kein committeter Test. Schwelle erreicht: dieser Eintrag ist
ab jetzt **Lese-Schritt-Kandidat** für die Closure von
`welle-transformationen` (die Roadmap führt sie unter *Offene Wellen*) —
nur benannt, kein Architect-Zug und keine ADR in diesem Zug; der Architect
entscheidet bei jenem Lese-Schritt, ob ein Ein-Instanz-Wächter für den
Prozessstart eine Entscheidung braucht.
