package inbound

import "context"

// ExportMetricsUseCase führt einen Zyklus des Metrik-Exports aus
// (`LH-FA-SST-010`): lesen, den zuletzt gemessenen WAL-Rückstand ergänzen,
// übertragen. Der Aufrufer — der Takt des Prozesses — ruft ihn in einer
// eigenen Goroutine und wertet den Rückgabewert nicht aus: der Use Case
// protokolliert jeden Fehlschlag selbst, und ein Fehlschlag ändert weder
// Erfassung noch Health.
type ExportMetricsUseCase interface {
	Export(ctx context.Context) error
}
