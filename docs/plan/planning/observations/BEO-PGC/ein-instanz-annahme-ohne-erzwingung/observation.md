# BEO-PGC/ein-instanz-annahme-ohne-erzwingung

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die Annahme eines
Backfill-Antrags im Postgres-Adapter, keine eigene Sub-Area im Sinn der
Modus-Deklaration).

Die Beobachtung: Der Annahme-Adapter des Backfills (`NewBackfillAdmission`, Paket
`postgresstorage`) prüft „kein aktiver Run derselben Tabelle" (Lesen) und legt danach die
Run-Zeile `queued` an, beides in **einer** Transaktion unter `READ COMMITTED`, ohne Sperre
und ohne Unique-Kante auf aktive Runs: `cdc.backfill_run` trägt nur den Primärschlüssel
`run_id`. Zwei gleichzeitig annehmende Verbindungen sähen beide „kein aktiver Run" und
legten zwei `queued`-Zeilen derselben Tabelle an; der Worker arbeitet jede `queued`-Zeile
ab, der Bestand würde zweimal kopiert (abgeleitet, nicht gemessen). Die Zusage trägt die
Annahme, dass `processAdministrationRequests` Anträge sequenziell in **einer** Goroutine
verarbeitet und je Quelle **eine** Instanz Anträge annimmt (`ADR-0113` Festlegung 1
Punkt 2; im Godoc des Adapters und im Plan des Slice benannt). Sie ist weder erzwungen
noch getestet: kein Test fährt zwei Annehmende gleichzeitig, kein Schema-Objekt
verhindert die zweite Zeile.

**Warum das zählt:** Eine benannte Annahme ohne Erzwingung trägt nur, solange der Betrieb
sie hält. Ein Betreiber, der eine zweite Instanz je Quelle startet — die Bedingung des
Re-Evaluierungs-Triggers von `ADR-0113` —, bekommt keinen Fehler, sondern zwei Runs
derselben Tabelle. Kein Gate liest die Zahl der Instanzen; der Wächter ist der Leser, der
die Annahme gegen den Betrieb hält.

Deklaration: `slice-backfill-run-store`, Risiko §6 fünfter Punkt („Die Annahme-Transaktion
ist nicht atomar oder greift zu weit"), Review F-5 (INFO), Ausgang *weiter offen*.
