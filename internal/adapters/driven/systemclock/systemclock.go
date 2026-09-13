// Package systemclock trägt den `SystemClockAdapter` als Driven-
// Implementierung des `ClockPort` (`ARC-006`, `ADR-0040`): die Wanduhr
// bleibt Infrastruktur und lebt ausschließlich hier — Domain- und
// Application-Code rufen `time.Now()` nicht direkt auf, sie beziehen den
// aktuellen Zeitpunkt über den Port; Tests setzen eine Fake-Uhr ein
// (`outbound.ClockPort`).
package systemclock

import (
	"time"

	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Adapter implementiert `outbound.ClockPort` über die Wanduhr des
// Prozesses (`time.Now`); er trägt keinen eigenen Zustand.
type Adapter struct{}

var _ outbound.ClockPort = Adapter{}

// New baut den System-Clock-Adapter.
func New() Adapter {
	return Adapter{}
}

// Now liefert den aktuellen Zeitpunkt als Domänenwert (`model.TimePoint`,
// Unix-Nanosekunden) — die einzige Stelle dieses Pakets, die `time.Now()`
// aufruft.
func (Adapter) Now() model.TimePoint {
	return model.NewTimePoint(time.Now().UnixNano())
}
