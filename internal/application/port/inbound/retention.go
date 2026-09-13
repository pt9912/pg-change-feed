// Die Retention-Use-Cases trägt diese Datei an einer Stelle (`ARC-002`):
// die Driving-Adapter (`ARC-005`) rufen die Bereinigung über sie auf
// (`ADR-0028`, `ADR-0014`); die Orchestrierung liegt im Application
// Service.

package inbound

import (
	"context"

	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// RunRetentionCommand trägt die Eingabe des RunRetention Use Cases: die
// Quelle, deren Changes bereinigt werden, und die anzuwendende Policy
// (`LH-FA-RET-002`…`004`, `ADR-0014`).
type RunRetentionCommand struct {
	Source model.SourceID
	Policy model.RetentionPolicy
}

// RunRetentionResult trägt den Ausgang der Bereinigung: die Anzahl der
// real gelöschten Changes.
type RunRetentionResult struct {
	Deleted int
}

// RunRetentionUseCase führt eine Bereinigung für eine Quelle aus
// (`LH-FA-RET-002`…`004`, `ADR-0014`): die Freigabe je Change trägt
// `RetentionPolicy.AllowsDeletion` — Alter, Change-Position und alle
// bestätigten Consumer-Positionen der Quelle —, die physische Löschung der
// freigegebenen Changes trägt der `ChangeStorePort`. Kein blockierender
// Consumer wird über diesen Use Case sichtbar gemacht (`LH-FA-RET-005`,
// eigener Folge-Slice).
type RunRetentionUseCase interface {
	Run(ctx context.Context, command RunRetentionCommand) (RunRetentionResult, error)
}
