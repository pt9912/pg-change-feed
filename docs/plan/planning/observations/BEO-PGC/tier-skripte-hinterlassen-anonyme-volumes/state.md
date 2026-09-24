Zustand: offen (**1×**) — unter der Schwelle, kein Ausgang zugewiesen. Ein Fix ist nicht Teil
des Vorgangs. Ausgangs-Kandidat (nicht entschieden): `docker rm -fv` in den `cleanup`-Zweigen
der Tier-Skripte, die einen PostgreSQL-Container mit anonymem Datenvolume anlegen; die
Läufe, die ausschließlich benannte Volumes nutzen, bleiben unberührt.

Zähler (abgeleitet): **1×** (evidence/slice-backfill-sql-administration.md).

**Verwandt, nicht doppelt gezählt:** `BEO-PGC/test-schreibt-in-committete-datei` (dort ein
committetes Erzeugnis als Tier-Nebenwirkung; hier ein Docker-Objekt außerhalb des Repos).
