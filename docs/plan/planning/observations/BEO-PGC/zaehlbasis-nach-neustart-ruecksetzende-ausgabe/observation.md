# BEO-PGC/zaehlbasis-nach-neustart-ruecksetzende-ausgabe

**Sub-Area:** `*`/`PGC` (Realserver-Phasen des Integrationsrunners unter `tools/harness/`).

Die Beobachtung: Eine Runner-Phase belegt „nach dem Neustart eines Dienstes kommt wieder Ausgabe
an“ über das Wachstum einer Ausgabedatei und nimmt die Basis (Zeilenzahl) vor dem Stopp. Der
Dienst legt die Datei bei jedem Start leer neu an; die Basis aus der Zeit vor dem Stopp liegt
dann über jedem Wert, den ein Fenster nach dem Start erreicht. Die Phase wartet auf einen Wert,
den die neue Datei erst nach mehr Takten als das Fenster trägt, und endet rot ohne Aussage über das
Produkt. Lokal lief sie grün (anderes Zeitverhalten), am gehosteten Runner rot.

**Form (Ausprägung):** Zählbasis überlebt den Neustart des gezählten Dienstes. Verwandt mit
`BEO-PGC/ready-ist-nicht-verbunden` (Beleg wartet auf den falschen Zeitpunkt) und
`BEO-PGC/github-actions-unverifizierbar-lokal` (lokaler Grünlauf sagt nichts über den Runner).

Erstes Auftreten: `slice-otlp-metrik-export-e2e`, Phase „Wiederaufnahme nach docker start“,
e2e.yml Lauf `37225751324`, Leg PostgreSQL 18; Behebung `a19c28cc` (der Runner wartet auf
„Everything is ready“ im Log des Collectors seit dem Start und nimmt die Basis an der neuen Datei).
