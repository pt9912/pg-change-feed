// Der Changes-Lese-Use-Case trägt diese Datei an einer Stelle (`ARC-003`):
// die Driving-Adapter lesen persistierte Changes über ihn; die
// Orchestrierung liegt im Application Service.

package inbound

import (
	"context"

	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// ReadChangesQuery trägt die Eingabe des Changes-Lesens (`LH-FA-REA-001`):
// die Quelle (**Pflicht**), die Klartext-Filter `Schema`/`Table` (je
// optional und **unabhängig**), den Positionsbereich `Start`/`End`
// (optional, Start **inklusiv**, End **exklusiv**) und
// `Limit` (optional, gesetzt ≥ 1). Ein nicht gesetztes
// Limit liest unbegrenzt — es gibt **kein** Default-Limit: die Antwort
// bleibt damit deckungsgleich mit dem
// View-Direktzugriff auf `cdc.changes`.
type ReadChangesQuery struct {
	Source model.SourceID
	Schema string
	Table  string
	Start  *model.SourcePosition
	End    *model.SourcePosition
	Limit  *int
}

// ReadChange bündelt einen gelesenen Change mit der Commit-Position seiner
// Quelltransaktion und ihrem Commit-Zeitpunkt — dieselbe Bündelung, die der
// Outbound-Port als `ChangeRecord` führt (Ordnungsgröße an der Transaktion,
// `LH-FA-REA-004`).
type ReadChange struct {
	Position    model.SourcePosition
	Change      model.Change
	CommittedAt model.TimePoint
}

// ReadChangesResult trägt die gelesenen Changes in der Ordnung des
// Lesepfads (`LH-FA-REA-004`); eine leere Menge trägt eine leere, gesetzte
// Liste.
type ReadChangesResult struct {
	Changes []ReadChange
}

// ReadChangesUseCase liest persistierte Changes über den bestehenden
// `ChangeStorePort` (`LH-FA-REA-001`): er ist eine
// Übersetzung, keine eigene Abfrage — dieselbe Delegation, dieselbe
// Sortierung wie der View-Direktzugriff, kein zweiter Lesepfad.
type ReadChangesUseCase interface {
	ReadChanges(ctx context.Context, query ReadChangesQuery) (ReadChangesResult, error)
}
