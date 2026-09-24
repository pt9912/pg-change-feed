package outbound

import (
	"context"

	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// BackfillAdmissionPort trägt die Fähigkeit, einen Backfill-Antrag
// anzunehmen (`ARC-004`, Fähigkeits-Port je `ADR-0034`; `ADR-0113`
// Festlegung 1). Er ist der einzige Weg, eine Run-Zeile anzulegen: der
// `BackfillRunPort` trägt keine Anlage-Operation, damit es keinen zweiten,
// nicht atomaren Weg gibt.
type BackfillAdmissionPort interface {
	// Admit nimmt den Antrag `requestID` an: die Prüfung „kein aktiver Run
	// (`queued`/`running`) derselben Tabelle", die Anlage der Run-Zeile im
	// Zustand `queued` und der Antragsvermerk `applied` — „angenommen" —
	// sind **eine** Einheit. Ein aktiver Run endet als
	// `domainerrors.ErrBackfillRunActive`; jede Abweichung hinterlässt
	// weder Run-Zeile noch Antragsvermerk. `run.ID` trägt die Kennung des
	// Antrags. Persistenzfehler tragen `ErrBackfillStorage`.
	Admit(ctx context.Context, requestID model.AdministrationRequestID, run model.BackfillRun) error
}
