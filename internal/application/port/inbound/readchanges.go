// Der Changes-Lese-Use-Case trägt diese Datei an einer Stelle (`ARC-003`):
// die Driving-Adapter (`ARC-005`) lesen persistierte Changes über ihn
// (`ADR-0028`, `ADR-0081`); die Orchestrierung liegt im Application Service
// (`ARC-002`).

package inbound

import (
	"context"

	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// ReadChangesQuery trägt die Eingabe des Changes-Lesens: die Quelle
// (**Pflicht**), die Klartext-Filter `Schema`/`Table` (je optional und
// **unabhängig**, `LH-FA-REA-006`), den Positionsbereich `Start`/`End`
// (optional, Start **inklusiv**, End **exklusiv**, `LH-FA-REA-001`) und
// `Limit` (optional, gesetzt ≥ 1, `LH-FA-REA-003`). Ein nicht gesetztes
// Limit liest unbegrenzt — es gibt **kein** Default-Limit (`ADR-0081`
// Teilfrage 2): die Antwort bleibt damit deckungsgleich mit dem
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
// `LH-FA-REA-004.a`, `ADR-0081` Teilfrage 1).
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
// `ChangeStorePort` (`LH-FA-REA-001` ff., `ADR-0081`): er ist eine
// Übersetzung, keine eigene Abfrage — dieselbe Delegation, dieselbe
// Sortierung wie der View-Direktzugriff, kein zweiter Lesepfad
// (`LH-FA-SST-006` Boundary).
type ReadChangesUseCase interface {
	ReadChanges(ctx context.Context, query ReadChangesQuery) (ReadChangesResult, error)
}
