package inbound

import (
	"context"

	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// IdleConfirmationCommand trägt die Eingabe der Leerlauf-Bestätigung: die
// Position, bis zu der die Quelle in ihrer Keepalive-Nachricht ihr WAL-Ende
// nennt, während der Stream keine Quelltransaktion offen hat
// (`ADR-0120` Festlegung 1). Die Position ist keine Position eines
// Changes: sie trägt kein persistiertes Guthaben.
type IdleConfirmationCommand struct {
	Position model.SourcePosition
}

// IdleConfirmationResult trägt den Ausgang der Leerlauf-Bestätigung: die
// bestätigte Quellposition. Eine leere Position bedeutet, dass die
// Application nicht bestätigt hat.
type IdleConfirmationResult struct {
	Acknowledged model.SourcePosition
}

// IdleConfirmationInboundPort nimmt die Meldung „Leerlauf bis P“ des
// Replication-Stream-Adapters auf (`ADR-0120`, `ADR-0034`): die Application
// entscheidet über die Bestätigung, der Adapter meldet nur den Zustand
// (`ADR-0007`). Der Port ist von `CaptureInboundPort` getrennt, weil er eine
// eigene Konsistenzgrenze trägt — eine Bestätigung ohne Persistenz, gültig
// allein im Leerlauf, in dem kein Change der Publication zu speichern ist.
type IdleConfirmationInboundPort interface {
	ConfirmIdle(ctx context.Context, command IdleConfirmationCommand) (IdleConfirmationResult, error)
}
