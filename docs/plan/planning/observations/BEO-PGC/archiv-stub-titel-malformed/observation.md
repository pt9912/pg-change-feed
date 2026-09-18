# BEO-PGC/archiv-stub-titel-malformed

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft das
Zusammenspiel zwischen dem externen Archivierungs-Werkzeug und den
Titelformen der von ihm archivierten Zeitdokumente, keine eigene Sub-Area
im Sinn der Modus-Deklaration).

Die Beobachtung: Das externe, vendored Werkzeug (`ai-harness-init
archive-welle`) erzeugt für jeden archivierten Slice/jede archivierte Welle
einen gekürzten Stub aus einer Titel-Vorlage, die eine bestimmte
Quelltitelform voraussetzt. Trifft die reale Quelltitelform diese Annahme
nicht, entsteht ein malformter Stub-Titel — bisher in zwei voneinander
unabhängigen Auslösebedingungen beobachtet:

1. **Namensbasierte Slice-Kennung** (`MR-002`, ab `slice-105` exklusiv):
   das Werkzeug nimmt einen rein numerischen `slice-<NNN>` an und rendert
   z. B. `# slice- — slice-d-check-tracked-modul: Titel` statt
   `# slice-d-check-tracked-modul — Titel` (gefunden in
   `slice-archive-altbestand-vollzug`, Review-Finding F-1).
2. **Numerische Welle-Quelltitelform ohne `welle-`-Präfix** (`# Welle
   <N>: Titel` statt `# Welle welle-<N>: Titel`): das Werkzeug verdoppelt
   die Nummer, z. B. `# welle-14 — 14: Performance-Benchmarks &
   Test-Coverage-Gate` statt `# welle-14 — Performance-Benchmarks &
   Test-Coverage-Gate` (gefunden bei 13 von 17 Wellen in
   `slice-archive-wellen-verbleibend`, Review-Finding F-1).

Beide Auslöser sind Instanzen derselben zugrunde liegenden
Werkzeug-Schwäche (starre Titel-Parsing-Annahme, kein generischer
Titel-Extraktor), aber strukturell verschieden (Slice- vs. Welle-Ebene,
Namens- vs. Nummernform) — deshalb ein eigener Register-Eintrag statt einer
Zusammenlegung mit `BEO-PGC/externes-werkzeug-committet-ohne-kennung`
(dort: Commit-Message-Kennung, hier: Stub-Titel-Form — verschiedene
Werkzeug-Facetten).
