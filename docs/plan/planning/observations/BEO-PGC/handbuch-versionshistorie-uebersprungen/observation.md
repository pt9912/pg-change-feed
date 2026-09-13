# BEO-PGC/handbuch-versionshistorie-uebersprungen

**Sub-Area:** Planning-Harness / Dokumentation (Sub-Area-Kürzel `PGC` aus
der Modus-Deklaration)

Die Beobachtung: `docs/user/benutzerhandbuch.md` trägt einen `Version:`-
Kopf plus eine `### Änderungshistorie`-Tabelle, die jede inhaltliche
Änderung als eigene, fortlaufend nummerierte Zeile führt — jeder
DoD-Punkt „Doku-Update" in einem Slice-Plan hat das bislang konsistent
getan (`slice-020`, `-023`, `-036`/`037`/`042`, `-038`, `-041`). Zwei
Slices in Folge (`slice-045`, `slice-046`) haben `benutzerhandbuch.md`
real mit neuem Inhalt erweitert (§4 „Blockierende Consumer erkennen" bzw.
„Metriken lesen"/`cdc_storage_bytes`), dabei aber weder den `Version:`-
Kopf noch die Änderungshistorie-Tabelle fortgeschrieben — der DoD-Punkt
„Doku-Update" verlangt inhaltlich nur, dass die neue Fähigkeit genannt
wird, nicht ausdrücklich eine Versionszeile. Der Lückenschluss fiel erst
bei einem `slice-047`-Review auf (der `slice-047`-Implementer nummerierte
seine eigene neue Zeile als „1.8", direkt auf die letzte tatsächlich
vorhandene Zeile „1.7" folgend, ohne die beiden übersprungenen
Änderungen zu bemerken).
