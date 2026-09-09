package model

import (
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

// TimePoint trägt einen Zeitpunkt als Unix-Nanosekunden. Die Domain
// importiert das Paket `time` nicht (`ADR-0040`, Fitness-Funktion): der
// Produktions-System-Clock-Adapter mappt `time.Time` auf `UnixNanos`.
type TimePoint struct {
	UnixNanos int64
}

// NewTimePoint trägt die Unix-Nanosekunden in einen Zeitpunkt.
func NewTimePoint(unixNanos int64) TimePoint {
	return TimePoint{UnixNanos: unixNanos}
}

// IsZero meldet den Nullwert.
func (t TimePoint) IsZero() bool {
	return t.UnixNanos == 0
}

// Before meldet, dass t vor other liegt.
func (t TimePoint) Before(other TimePoint) bool {
	return t.UnixNanos < other.UnixNanos
}

// After meldet, dass t nach other liegt.
func (t TimePoint) After(other TimePoint) bool {
	return t.UnixNanos > other.UnixNanos
}

// Add verschiebt t um d; das Ergebnis ist ein neuer Wert.
func (t TimePoint) Add(d Duration) TimePoint {
	return TimePoint{UnixNanos: t.UnixNanos + d.Nanos}
}

// Sub trägt den Abstand zwischen t und other als Dauer.
func (t TimePoint) Sub(other TimePoint) Duration {
	return Duration{Nanos: t.UnixNanos - other.UnixNanos}
}

// Duration trägt einen Zeitraum in Nanosekunden; sie ist nicht negativ
// (`NewDuration`).
type Duration struct {
	Nanos int64
}

// NewDuration legt eine Dauer an und verlangt einen nichtnegativen Wert.
func NewDuration(nanos int64) (Duration, error) {
	if nanos < 0 {
		return Duration{}, domainerrors.ErrNegativeDuration
	}
	return Duration{Nanos: nanos}, nil
}
