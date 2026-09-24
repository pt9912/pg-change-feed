package outbound

import (
	"context"
	stderrors "errors"

	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// ErrBackfillRequestNotPending: der Antrag, den `Admit` annehmen soll, ist
// nicht offen (`pending`) oder besteht nicht — die Annahme hinterlässt weder
// Run-Zeile noch Antragsvermerk.
var ErrBackfillRequestNotPending = stderrors.New("Backfill-Antrag ist nicht offen: die Annahme trifft keine pending-Zeile")

// ErrBackfillRunInvalid: der übergebene Run passt nicht zur Operation — er
// trägt nicht die Kennung des Antrags, ist nicht `queued` (`Admit`) bzw. nicht
// `running` (`BackfillWriterPort.Begin`) oder trägt nicht die Angaben, die
// die Operation voraussetzt.
var ErrBackfillRunInvalid = stderrors.New("Backfill-Run passt nicht zur Operation")

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
	// `domainerrors.ErrBackfillRunActive`; ein Antrag, der nicht `pending`
	// ist, als `ErrBackfillRequestNotPending`; ein Run, der nicht die Kennung
	// des Antrags trägt oder nicht `queued` ist, als `ErrBackfillRunInvalid` —
	// ebenso ein Antrag, der nicht die Art `backfill` trägt oder nicht Quelle,
	// Schema und Tabelle des Runs adressiert.
	// Jede Abweichung hinterlässt weder Run-Zeile noch Antragsvermerk.
	// `run.ID` trägt die Kennung des Antrags. Persistenzfehler tragen
	// `ErrBackfillStorage`.
	Admit(ctx context.Context, requestID model.AdministrationRequestID, run model.BackfillRun) error
}
