# BEO-PGC/lesesperre-ohne-zeitgrenze

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die Lesesperre des
Snapshot-Imports beim Backfill, keine eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Der Snapshot-Import des Backfills sperrt die Quelltabelle mit
`LOCK TABLE … IN ACCESS SHARE MODE` (`internal/adapters/driven/postgressnapshot`,
`lockAndVerify`). Die Wartezeit an dieser Anweisung trägt allein der Kontext des Aufrufs
(`ADR-0118` Festlegung 1); der Kontext, den der Worker übergibt, trägt kein Zeitlimit.
Hält eine fremde Transaktion `ACCESS EXCLUSIVE` auf der Tabelle, steht der Run in
`running` mit `rows_copied` 0, solange diese Transaktion offen ist, und hält in dieser
Zeit seinen Snapshot. Gemessen ist das Ende des Wartens durch einen ablaufenden Kontext:
der Import kehrt mit der Klasse `transient` zurück und hinterlässt keine Sitzung
(`TestImportLockWaitEndsWithTheContext`); eine Zeitgrenze des Runs selbst gibt es nicht.

**Warum das zählt:** Ein Betreiber sieht `running` ohne Fortschritt und ohne Ablauf, solange
die fremde DDL-Transaktion offen ist; das Handbuch nennt die Gegenrichtung (Abschnitt „Sperre
der Tabelle“), aber kein Gate und keine Konfiguration begrenzt das Warten.

Deklaration: `slice-backfill-e2e` (Review F-6), Risiko §6 „Wartegrenze der Lesesperre“,
Ausgang *weiter offen*.
