package outbound

import (
	"context"

	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// ChangeStorePort trägt die Persistenz-Fähigkeit (`ARC-009`): er
// persistiert committed CDC-Transaktionen dauerhaft und ist die einzige
// Persistenz-Grenze des Capture-Pfads (`ADR-0009`). Die erste
// Produktionsimplementierung ist der `PostgresChangeStoreAdapter`;
// Tests setzen Fake-Ports ein (`ADR-0030`).
//
// Der Port bestätigt keine Quellpositionen. Bestätigt wird eine Position
// erst, wenn ihre abhängigen Changes dauerhaft gespeichert sind — ACK nur
// über nach Persistenz gefragte Positionen (`LH-QA-REL-001.a`,
// `ADR-0029` Regel 1). Die Quell-Bestätigung ist die Wirkung des
// `ReplicationAckPort` (`ADR-0007`), keine Speicherwirkung dieses Ports.
type ChangeStorePort interface {
	// PersistTransaction persistiert eine committed Quelltransaktion
	// dauerhaft — mit interner ID, Commit-Position und Changes in
	// Sequenz-Reihenfolge (`SPEC-001`, `cdc.transaction`). Die Rückkehr
	// ohne Fehler meldet den abgeschlossenen Store-Commit
	// (`LH-QA-REL-001.a`, Schritt COMMIT Store).
	PersistTransaction(ctx context.Context, transaction *model.ChangeTransaction) error
}
