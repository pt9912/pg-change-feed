# BEO-PGC/gate-modul-abgeschaltet-trotz-regel

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft das Verhältnis
zwischen einer gewollten Regel und dem Sensor, der sie durchsetzen soll, keine
eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Eine Regel ist als Absicht vorhanden — sie steht in
Entscheidungen, wird zitiert, ist Gegenstand von Verdikten —, aber **ihr Sensor
ist abgeschaltet**: das Modul fehlt in der Modul-Liste, das Target hängt nicht
im Bündel, die Prüfung steht auskommentiert. Der Bestand verletzt sie dann über
viele Vorgänge hinweg, ohne dass etwas meldet; auffindbar wird das erst, wenn
jemand die Regel einfodert und die Messung nachholt. Belegt an `slice-078`:
`hostpaths` fehlte in der `modules`-Liste, während dieses Repo die Regel gegen
host-lokale Pfade längst wollte — **42 Vorkommen in 15 Dateien** hatten sich
angesammelt.
