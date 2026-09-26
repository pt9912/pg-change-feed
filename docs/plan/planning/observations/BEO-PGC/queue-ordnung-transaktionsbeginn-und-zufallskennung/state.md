Zustand: offen (**1×**) — unter der Schwelle, kein Ausgang zugewiesen. Die Instanz ist behoben
(seit slice-transformationen-antragsweg-usecase):
die sieben SQL-Funktionen schreiben `requested_at` je Aufruf mit `clock_timestamp()`, die
Verarbeitung ordnet wie die Ableitung nach `(requested_at, administration_request_id)`; die
Entscheidung steht in `ADR-0127`, der Wortlaut in `SPEC-019` (Absatz „Ordnung der
Verarbeitung“), gebunden durch `TestAdministrationRequestListPendingOrdersTiesByRequestID`,
`TestAdministrationRequestSameTransactionCallsKeepCallOrder` (alle sieben Funktionen in beiden
Richtungen) und `TestAdministrationSameTransactionRemoveThenSetLeavesTheNewRuleLiveAndDerived`.
Die Klasse bleibt offen, bis ein zweites Vorkommen sie bestätigt: eine weitere Ordnung, deren
Schlüssel ein Transaktionszeitstempel oder eine Zufallskennung ist und die zwei Leser tragen
(Kandidat: `cdc.backfill_run.requested_at` mit `run_id`, dort aus der Uhr des Prozesses, `ADR-0113`
Festlegung 2 — gelesen, nicht als Auftreten gezählt). Grenzen der Behebung (in `ADR-0127`
Festlegung 3 benannt, hergeleitet, nicht erprobt): Aufrufzeitpunkt statt Festschreibung bei
überlappenden Transaktionen mehrerer Sitzungen, Rückwärtssprung der Serveruhr, Zeilen vor der
Änderung tragen den Transaktionsbeginn — **akzeptiertes Negativ**, kein Register-Eintrag je Grenze;
das Handbuch trägt die Betreiberregel (`slice-transformationen-betriebsdoku` §2).
Zähler (abgeleitet): **1×** (evidence/slice-transformationen-antragsweg-usecase.md).
