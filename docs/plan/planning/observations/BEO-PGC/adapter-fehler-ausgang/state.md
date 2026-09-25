Zustand: **geplant** — Ausgang: **geplant** →
[`slice-capture-transient-wiederholung`](../../../open/slice-capture-transient-wiederholung.md)
(Start: eine ADR des Architects zur Wiederholungsform; Klasse `transient`,
begrenzter Backoff, Ausgang bei Erschöpfung) · seit welle-backfill-bestand
(Architect-Verdikt `architect-verdict-welle-backfill-bestand-lese-schritt`
§4.1 und §8).

Gegenstand: der Erfassungspfad endet auf jeden Adapter-Fehler mit Prozess-Ausgang
1; die Aktion der Klasse `transient` (`SPEC-008`) trägt kein Element des Pfads. Eine
Ausprägung: der Adapter `receive` wiederholt `START_REPLICATION` bei SQLSTATE 55006
(Slot noch aktiv) nicht; der Test `TestStreamRestartsOnExistingSlot` wartet auf die
Freigabe des Slots.

Zähler: 3× (Dateien unter `evidence/`).
