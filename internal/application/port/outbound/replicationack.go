package outbound

import (
	"context"

	"github.com/pt9912/pg-change-feed/internal/domain/messagecode"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// ErrReplication trägt die Fehlerklasse `replication` des
// PostgresReplicationAckAdapters (`ADR-0023`): eine
// Störung am Replication-Stream/Slot endet als sichtbarer Fehler — die
// Bestätigung gilt erst mit der Rückkehr ohne Fehler
// (Schritt ACK Source). Application und Betrieb
// klassifizieren über `errors.Is` und kennen keinen Treibertyp; die
// technische Ursache bleibt über die zweite Wrappung lesbar.
var ErrReplication = messagecode.New(messagecode.ReplicationAckFailed, "Bestätigung am Replication-Stream fehlgeschlagen")

// ReplicationAckPort trägt die Quell-Bestätigung als Core-gesteuerte
// Wirkung (`ADR-0007`): der Replication-Stream-Adapter liest denselben
// Quell-Stream, entscheidet aber nicht selbst, wann eine Position
// dauerhaft verarbeitet ist. Die Application ruft ihn erst,
// wenn die abhängigen Changes dauerhaft gespeichert sind.
// Die Driven-Implementierung ist der
// `PostgresReplicationAckAdapter`.
type ReplicationAckPort interface {
	// Acknowledge bestätigt die Position gegenüber der Quelle; die
	// Bestätigung gilt innerhalb der Quelle der Position
	// (`LH-QA-REL-001.a`, Schritt ACK Source).
	Acknowledge(ctx context.Context, position model.SourcePosition) error
}
