# BEO-PGC/zwei-quellen-drift-handbuch-gegen-pflichtenheft

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft dieselbe Aussage in
Benutzerhandbuch und Pflichtenheft, keine eigene Sub-Area im Sinn der
Modus-Deklaration).

Die Beobachtung: Das Benutzerhandbuch nennt für einen Mechanismus einen anderen Wert
oder Zeitpunkt als das Pflichtenheft und der Mechanismus selbst. Belegt am Stichtag
des Bestands eines Backfills: das Pflichtenheft (`LH-FA-CAP-009.a`) und der Mechanismus
tragen den Startzeitpunkt des Runs (der Snapshot entsteht mit dem Slot beim Beginn der
Ausführung), das Handbuch nannte an zwei Stellen den Zeitpunkt des Antrags. Ein Run wartet
`queued` hinter einem anderen Run oder über einen Neustart; beide Zeitpunkte liegen dann
auseinander. Das Handbuch steht im Rang unter dem Pflichtenheft (Source Precedence); kein
Gate hält beide gegeneinander.

**Warum das zählt:** Ein Betreiber liest den Stichtag im Handbuch und plant danach; eine
Abweichung fällt erst auf, wenn jemand Handbuch, Spec und Code nebeneinander liest — der
Reviewer als Leser der Zwei-Quellen-Drift-Regel.

Deklaration: `slice-backfill-sql-administration`, Review F-5 (MEDIUM).
