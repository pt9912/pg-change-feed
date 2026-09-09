# BEO-PGC/a-check-null-abdeckung

**Sub-Area:** Application-/Domain-Layer-Gate-Abdeckung (`internal/**`-Globs
der `.a-check.yml`; Sub-Area-Kürzel `PGC` aus der Modus-Deklaration)

Die Beobachtung: Ein a-check-Lauf ist grün, während die Layer-Globs null
Dateien matchen (Abdeckungs-/Auflösungs-Hinweis statt Prüfbereich) — das
Grün sagt dann nichts über die Schichten-Regeln. Der Lauf meldet den
Zustand als Hinweis auf stderr, wechselt aber den Exit nicht
(`harness/sensors/a-check.md`, Grenze 1 und 2).