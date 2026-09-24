Zustand: offen — Ausgang: **weiter offen** → Retry/Backoff folgt mit der
Konfigurationsschicht (späterer Slice); die Grenze trägt der Kommentar in
`wiring.go` und die Container-Vertrags-Zeile. Zähler (abgeleitet): **3×**
(evidence/slice-007.md, evidence/slice-038.md,
evidence/slice-backfill-snapshot-reader.md) — Schwelle erreicht mit
`slice-backfill-snapshot-reader`, Ausgang noch **nicht** zugewiesen; er gehört
dem Lese-Schritt der Closure von `welle-backfill-bestand` (Modul 6). Der dritte
Beleg trägt eine andere Ausprägung desselben Ausgangs: der Adapter im Paket
`internal/adapters/driving/replication/receive` wiederholt `START_REPLICATION`
bei SQLSTATE 55006 (Slot noch aktiv) nicht, der Prozess endet mit Ausgang 1
(Klasse `replication`). **Adresse dieser Ausprägung:** der Slice, der den
Wiederholungsversuch der Konfigurationsschicht liefert, oder ein eigener kleiner
Zug am `receive`-Adapter; bis dahin trägt die Fortsetzung der Prozess-Neustart,
und der Test `TestStreamRestartsOnExistingSlot` wartet auf die Freigabe des
Slots.
