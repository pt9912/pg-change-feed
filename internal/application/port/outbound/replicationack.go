package outbound

import (
	"context"

	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// ReplicationAckPort trägt die Quell-Bestätigung als Core-gesteuerte
// Wirkung (`ADR-0007`): der Replication-Stream-Adapter liest denselben
// Quell-Stream, entscheidet aber nicht selbst, wann eine Position
// dauerhaft verarbeitet ist (`ADR-0027`). Die Application ruft ihn erst,
// wenn die abhängigen Changes dauerhaft gespeichert sind
// (`LH-QA-REL-001.a`). Die Driven-Implementierung ist der
// `PostgresReplicationAckAdapter`.
type ReplicationAckPort interface {
	// Acknowledge bestätigt die Position gegenüber der Quelle; die
	// Bestätigung gilt innerhalb der Quelle der Position
	// (`LH-QA-REL-001.a`, Schritt ACK Source).
	Acknowledge(ctx context.Context, position model.SourcePosition) error
}
