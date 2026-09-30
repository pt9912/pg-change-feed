Zustand: **umgesetzt** — Ausgang: der Slice
`slice-capture-transient-wiederholung` (`ADR-0135` legt die
Wiederholungsform fest: begrenzter Backoff am Stream-Zyklus in der
Composition Root, Klasse `transient` bei Erschöpfung) · seit
slice-capture-transient-wiederholung
(Architect-Verdikt `architect-verdict-welle-backfill-bestand-lese-schritt`
§4.1 und §8).

Gegenstand: der Erfassungspfad endet auf jeden Adapter-Fehler mit Prozess-Ausgang
1; die Aktion der Klasse `transient` (`SPEC-008`) trägt kein Element des Pfads. Eine
Ausprägung: der Adapter `receive` wiederholt `START_REPLICATION` bei SQLSTATE 55006
(Slot noch aktiv) nicht; der Test `TestStreamRestartsOnExistingSlot` wartet auf die
Freigabe des Slots.

Zähler: 3× (Dateien unter `evidence/`).

Beleg des Ausgangs: Verifikations-Report `docs/reviews/verifikation-slice-capture-transient-wiederholung.md` <!-- d-check:status-provenance -->
(DoD-Zeilen 1–4 getragen; Wiederholung am Stream-Zyklus, Klasse `transient` bei Erschöpfung).
