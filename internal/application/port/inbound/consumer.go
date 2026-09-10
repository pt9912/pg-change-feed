// Die Consumer-Use-Cases trägt diese Datei an einer Stelle (`ARC-003`):
// die Driving-Adapter (`ARC-005`) rufen Registrierung, Bestätigung,
// Positions-Lese und Entfernung über sie auf (`ADR-0028`, `ADR-0013`); die
// Orchestrierung liegt in den Application Services (`ARC-002`).

package inbound

import (
	"context"

	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// RegisterConsumerCommand trägt die Eingabe der Registrierung
// (`LH-FA-CON-001`): Kennung und Name des Consumers.
type RegisterConsumerCommand struct {
	Consumer model.ConsumerID
	Name     string
}

// RegisterConsumerResult trägt den Ausgang der Registrierung: den
// Consumer und den idempotenten Ausgang (`LH-FA-CON-001` Boundary) — eine
// bereits registrierte Kennung meldet das Ergebnis, ohne den Stand zu
// ändern.
type RegisterConsumerResult struct {
	Consumer          model.Consumer
	AlreadyRegistered bool
}

// AcknowledgeConsumerCommand trägt die Eingabe der Bestätigung
// (`LH-FA-CON-004`): der Consumer und die verarbeitete Position — der
// reguläre ACK verläuft nur vorwärts (`ADR-0029`, Regel 2).
type AcknowledgeConsumerCommand struct {
	Consumer model.ConsumerID
	Position model.SourcePosition
}

// AcknowledgeConsumerResult trägt den fortgeführten Stand der Position:
// der gespeicherte Fortschritt des Consumers nach der Bestätigung
// (`LH-FA-CON-003`).
type AcknowledgeConsumerResult struct {
	Position model.ConsumerPosition
}

// GetConsumerPositionQuery trägt die Eingabe des Positions-Lese:
// der Consumer, dessen bestätigte Position die Abfrage liest
// (`LH-FA-CON-003`, `LH-FA-CON-005`).
type GetConsumerPositionQuery struct {
	Consumer model.ConsumerID
}

// GetConsumerPositionResult trägt die bestätigte Position des Consumers:
// der Nullwert liest die definierte Anfangsposition eines Consumers ohne
// Bestätigung (`LH-FA-CON-005` Boundary).
type GetConsumerPositionResult struct {
	Position model.ConsumerPosition
}

// RemoveConsumerCommand trägt die Eingabe der administrativen Entfernung
// (`LH-FA-CON-006`): der Consumer, der keine Rolle mehr in der Retention
// trägt.
type RemoveConsumerCommand struct {
	Consumer model.ConsumerID
}

// RemoveConsumerResult trägt den Ausgang der Entfernung: die Rückkehr
// meldet, ob der Consumer entfernt war; ein erneuter Aufruf bleibt ohne
// Wirkung (`LH-FA-CON-006` Idempotenz).
type RemoveConsumerResult struct {
	Removed bool
}

// RegisterConsumerUseCase registriert einen benannten Consumer
// (`LH-FA-CON-001`, `ADR-0028`).
type RegisterConsumerUseCase interface {
	Register(ctx context.Context, command RegisterConsumerCommand) (RegisterConsumerResult, error)
}

// AcknowledgeConsumerUseCase bestätigt die verarbeitete Position eines
// Consumers (`LH-FA-CON-004`, `ADR-0028`); die Monotonie trägt die
// Domänen-Ordnung, kein administrativer Reset läuft als ACK durch
// (`ADR-0013`).
type AcknowledgeConsumerUseCase interface {
	Acknowledge(ctx context.Context, command AcknowledgeConsumerCommand) (AcknowledgeConsumerResult, error)
}

// GetConsumerPositionUseCase liest die bestätigte Position eines Consumers
// (`LH-FA-CON-003`, `LH-FA-CON-005`, `ADR-0028`).
type GetConsumerPositionUseCase interface {
	Position(ctx context.Context, query GetConsumerPositionQuery) (GetConsumerPositionResult, error)
}

// RemoveConsumerUseCase entfernt einen Consumer administrativ
// (`LH-FA-CON-006`): das Verhalten der persistierten Changes trägt die
// Retention, nicht die Entfernung.
type RemoveConsumerUseCase interface {
	Remove(ctx context.Context, command RemoveConsumerCommand) (RemoveConsumerResult, error)
}
