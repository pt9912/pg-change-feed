# Beleg: slice-026 (Schwellen-Überwachung mit kontrollierter Fortsetzung)

Vorgang: slice-026.

Fund: `slice-024` (`ADR-0049`, Accepted) trennte die `replication`-Klasse in
Stream-Ordnungs-Verletzung (bleibt hart abbrechend) und Transport-/
Verbindungsstörung (Schwellen-Kandidat) und setzte WAL-Rückstand-Schwellen
(Warn 100 MiB, Fehler 1 GiB). `slice-025` lieferte die Metrik
`cdc_wal_retention_bytes`. Dieser Slice verdrahtet den Schwellen-Vergleich
in `internal/bootstrap/wiring.go` (`classifyWALRetention`,
`runWALRetentionCheck`, `mergeStreamAndWALFaultOutcome`): unterhalb der
Warnschwelle unveränderte Fortsetzung, zwischen Warn-/Fehlerschwelle
sichtbare Fortsetzung (Log-Warnung), oberhalb der Fehlerschwelle
kontrollierter Abbruch über den bestehenden `replication`-Klassifikations-
pfad. `TestWALRetentionThresholdEndToEnd` belegt real, mit künstlich
erzeugtem WAL-Rückstand, beide Seiten der Schwelle in einem Lauf — dreimal
in Folge grün. Ein Regressionstest und unabhängige Prüfungen durch Reviewer
und Verifier (eigene Race-Analyse, `-race`-Läufe, eigene Rot-Grün-Gegenprobe)
bestätigen: Stream-Ordnungs-Verletzungen bleiben davon strukturell
unberührt und weiterhin hart abbrechend.

**Ausgang: eingetreten.** `SPEC-008`s Ziel-Zustand für die `replication`-
Fehlerklasse ist damit real gebaut, nicht mehr nur beschrieben — die Lücke
zwischen Spec und Code aus `slice-020`s Fund ist geschlossen.

Quelle: `docs/reviews/review-slice-026.md`, `docs/reviews/verify-slice-026.md`,
[`ADR-0049`](../../../../../adr/0049-replication-fehlerklassen-schwellen.md).
