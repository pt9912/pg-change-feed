# BEO-PGC/handbuch-beispiel-nicht-unter-genannter-rolle-lauffaehig

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft das Benutzerhandbuch
als Betreiber-Oberfläche, keine eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Ein Handbuch-Abschnitt nennt eine Voraussetzung (eine Login-Identität
mit einer bestimmten Rollen-Mitgliedschaft) und führt danach ein SQL-Beispiel, das unter
genau dieser Identität mit „permission denied" endet, weil das Recht bei einer anderen
Rolle liegt. Der Leser, der die Voraussetzung erfüllt, erhält für das Beispiel einen
Fehler. Die Aussage ist am Schema wahr (das Recht existiert), nur nicht unter der Rolle,
die der Abschnitt nennt; kein Gate führt ein Handbuch-Beispiel unter einem Login aus.

**Warum das zählt:** Ein Handbuch-Beispiel ist ein Beleg-Satz für den Betreiber
(`AGENTS.md` §3.12 Instanz B); er trägt nur, wenn jemand das Beispiel unter der
genannten Rolle ausgeführt hat. Der Wächter ist der Reviewer, der es fährt.

Deklaration: `slice-backfill-sql-administration`, Review F-4 (MEDIUM).
