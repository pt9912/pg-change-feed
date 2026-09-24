# BEO-PGC/execute-ergebnis-widerspricht-persistierter-zeile

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft das Ergebnis eines Use-Case-Aufrufs
gegenüber dem dauerhaften Zustand, den es beschreibt, keine eigene Sub-Area im Sinn der
Modus-Deklaration).

Die Beobachtung: Das Ergebnis eines Aufrufs (`BackfillExecuteResult.Run`) und die dauerhafte
Zeile (`cdc.backfill_run`) können einander widersprechen. Wirkt der Commit eines Runs
serverseitig und erhält der Client einen Verbindungsfehler, rollt der Use Case zurück und ruft
`Finish(failed)`; nach dem Port-Vertrag ist das auf einen beendeten Run ein wirkungsloser
Erfolg, die Zeile bleibt `completed`. Das Ergebnis meldet dennoch `failed` mit Fehlertext, und
das Wecksignal entfällt, obwohl die Daten sichtbar sind. Die Zeile ist der dauerhafte Zustand,
das Ergebnis der Zustand des **Aufrufs**; der Use Case liest die Zeile nicht zurück.

**Warum das zählt:** Ein Aufrufer, der `result.Run.Status` als Zeilenzustand liest, würde
neu beantragen und den Bestand doppelt kopieren; und die Richtigkeit hängt an einer
Adapter-Pflicht (`Finish` als wirkungsloser Erfolg), die vor dem Adapter-Slice nur ein Plan
trägt. Die Daten sind atomar committet, es entsteht keine Lücke; der Befund ist ein
Aufrufer-Vertrag, kein Datenfehler.

Deklaration: `slice-backfill-run-usecase`, Risiko §6 („Commit mit unbekanntem Ausgang"),
Ausgang *weiter offen*; Verifikation V-1 (LOW).
