package outbound

import (
	"context"

	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// MetricBatch ist eine Übertragung: die Quellkennung des Prozesses und alle
// Kennzahlen eines Zyklus samt Messzeitpunkt je Zeile.
type MetricBatch struct {
	Source  model.SourceID
	Samples []MetricSample
}

// MetricExportPort überträgt einen Zyklus an den Empfänger. Ein Adapter
// beachtet die Frist des übergebenen Kontexts, hält keine Warteschlange und
// wiederholt nicht: der Takt des Aufrufers ist die Wiederholung. Jeder
// Fehlschlag geht über `ErrMetricExport` zurück; der Fehlertext trägt weder
// die URL noch einen Header-Wert, und der Fehler liefert `FailureDetail`.
type MetricExportPort interface {
	Export(ctx context.Context, batch MetricBatch) error
}

// WALRetentionSource liefert den zuletzt im Prozess gemessenen
// `cdc_wal_retention_bytes` (`SPEC-009`) samt Messzeitpunkt; `ok` ist falsch,
// solange keiner gemessen ist. Der Wert gehört nicht zur Sicht der
// Betriebsschnittstelle (er liegt außerhalb der Leserechte), ihn schreibt der
// periodische Prüfzug des Prozesses.
type WALRetentionSource interface {
	LastWALRetention() (bytes int64, at model.TimePoint, ok bool)
}
