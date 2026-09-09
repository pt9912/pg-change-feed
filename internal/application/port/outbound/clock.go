// Package outbound bündelt die Outbound-Ports der Fähigkeiten
// (`ARC-004`): die Application Services fordern technische Wirkungen über
// sie an; die Driven-Adapter (`ARC-006`) setzen sie.
package outbound

import (
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// ClockPort reicht die Wanduhr als Fähigkeit an die Application
// (`ADR-0040`): Domain und Application rufen die Systemzeit nicht direkt
// auf; der System-Clock-Adapter ist die einzige
// Produktionsimplementierung, Tests setzen eine Fake-Uhr ein.
type ClockPort interface {
	Now() model.TimePoint
}
