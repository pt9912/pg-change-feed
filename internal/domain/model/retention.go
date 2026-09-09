package model

import (
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

// RetentionPolicy trägt die Grenzen der Bereinigung (`ADR-0014`). Die
// Zusage „Retention löscht keine benötigten Changes“ (`ADR-0029`, Regel 5)
// trägt `AllowsDeletion` am Beleg der Consumer-Positionen: ohne bestätigte
// Position an oder hinter der Change-Position gibt die Policy keine
// Bereinigung frei; das Mindestalter ist die zeitbasierte Grenze
// (`LH-FA-RET-003`). Die Ausführung liegt im Run-Retention-Use-Case, nicht
// in diesem Wertobjekt.
type RetentionPolicy struct {
	// MinAge ist das Mindestalter eines Changes für die Bereinigung;
	// 0 heißt „kein zeitliches Mindestalter“.
	MinAge Duration
}

// NewRetentionPolicy legt eine Policy an; das Mindestalter ist nicht
// negativ (`LH-FA-RET-003`).
func NewRetentionPolicy(minAge Duration) (RetentionPolicy, error) {
	if minAge.Nanos < 0 {
		return RetentionPolicy{}, domainerrors.ErrNegativeDuration
	}
	return RetentionPolicy{MinAge: minAge}, nil
}

// AllowsDeletion meldet, ob die Policy die Bereinigung eines Changes an
// seiner Position freigibt: der Change trägt das Mindestalter, und jede
// übergebene Consumer-Position ist bestätigt und liegt an oder hinter der
// Change-Position (`LH-FA-RET-004`, `ADR-0029`, Regel 5). Ohne
// übergebene Positionen blockiert keine Position die Bereinigung; das
// Mindestalter gilt weiter.
func (p RetentionPolicy) AllowsDeletion(age Duration, changePosition SourcePosition, consumerPositions []ConsumerPosition) bool {
	if age.Nanos < p.MinAge.Nanos {
		return false
	}
	for _, position := range consumerPositions {
		if !position.Acknowledged() || position.Position.Before(changePosition) {
			return false
		}
	}
	return true
}
