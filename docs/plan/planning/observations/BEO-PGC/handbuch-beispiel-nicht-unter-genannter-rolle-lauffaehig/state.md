Zustand: offen (**1×**) — unter der Schwelle, kein Ausgang zugewiesen. Behoben im
Vorgang (Handbuch 1.49: Antrag und Vermerk unter `cdc_admin`, Lesen von
`cdc.backfill_status` unter `cdc_reader`, die Rolle `cdc_admin` trägt kein `SELECT`
auf die View); der Verifier hat die Behebung am Handbuch-Diff und an Lauf 5 des
Guard-Test-Skripts (`SELECT` auf die View allein für `cdc_reader`) nachgemessen.
Ein Träger ist nicht vorgeschlagen: die Prüfung „Beispiel unter der genannten Rolle
ausführen" ist die Handlung des Reviewers, ein Sensor müsste Prosa-Beispiele einer
Rolle zuordnen.

Zähler (abgeleitet): **1×** (evidence/slice-backfill-sql-administration.md).

**Nicht zu verwechseln** mit `BEO-PGC/rollen-test-abdeckungsluecken` (dort fehlt der
Test einer Rollen-Zusage im Code; hier steht die Rollen-Aussage im Handbuch).
