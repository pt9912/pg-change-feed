# BEO-PGC/nicht-reproduzierbarer-test-ausfall

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die Test-Infrastruktur
der Go-Tests gegen reale PostgreSQL, keine eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Ein Test (`TestWALRetentionThresholdEndToEnd`, Paket `internal/bootstrap`)
endete im Lauf des Implementers einmal rot („cdc.changes trägt 0 Zeilen nach 30s",
`awaitPublishedChangeCount`) und dreimal grün bei unverändertem Stand. Der Reviewer maß
gegen eine Wegwerf-PostgreSQL 18 in der Tier-Konfiguration (`wal_level=logical`,
`wal_sender_timeout=2000`, `max_replication_slots=10`, je Lauf frische Datenbank): der Test
isoliert 30 von 30 grün am Slice-Stand und 30 von 30 grün am Parent `2d47d8a7`; das ganze
Paket `internal/bootstrap` 12 von 12 grün am Slice-Stand und 12 von 12 am Parent. Weder eine
Beteiligung der neuen Verdrahtung (drei Pools, ein `UPDATE` des Abgleichs, ein Worker vor
`stream.Run`) noch eine Anfälligkeit am Parent ist gezeigt oder ausgeschlossen. Der Ausfall
ist mit dem Repo-Skript nicht reproduzierbar; ein Fix ist nicht möglich, weil keine Ursache
vorliegt.

**Warum das zählt:** Ein einmal roter, nie reproduzierter Test entwertet das Vertrauen in
das Rot des nächsten Laufs; ohne Zähler bleibt jede Wiederholung ein Einzelfall.

Deklaration: `slice-backfill-sql-administration`, Review F-11 (INFO).
