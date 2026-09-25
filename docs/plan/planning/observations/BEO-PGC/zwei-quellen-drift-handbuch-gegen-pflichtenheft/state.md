Zustand: offen (**2×**) — unter der Schwelle, kein Ausgang zugewiesen. Beide Vorgänge
behoben; der zweite (`slice-backfill-sdk-origin`, Review F-5) liegt innerhalb des
Pflichtenhefts (Aufzählung in `LH-FA-SST-009.a` gegen die §6-Zeilen); der Eintrag zählt
die Klasse Zwei-Quellen-Drift, sein Name trägt die erste Ausprägung. Erster Vorgang: behoben im
Vorgang (beide Stellen nennen den Startzeitpunkt des Runs; `grep -n 'Zeitpunkt des Antrags'
docs/user/benutzerhandbuch.md` druckt 0 Zeilen, vom Verifier nachgemessen). Der Reviewer-Skill
trägt die Klasse als Punkt „Zwei-Quellen-Drift" (Abschnitt der Kategorien-Regeln); dieser
Eintrag zählt die Ausprägung Handbuch gegen Pflichtenheft.

Zähler (abgeleitet): **2×** (evidence/slice-backfill-sql-administration.md,
evidence/slice-backfill-sdk-origin.md).

**Verwandt, nicht gleich:** `BEO-PGC/lese-doppelquelle` (dort driftet dieselbe
Lese-Semantik zwischen SQL-View und Go-Use-Case, beides Code) und
`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (dort driftet ein Zahlenwert gegen
eine Messung).
