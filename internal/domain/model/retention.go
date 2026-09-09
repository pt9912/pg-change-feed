package model

import (
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

// RetentionPolicy trägt die Grenzen der Bereinigung (`ADR-0014`). Die
// Zusage „Retention löscht keine benötigten Changes“ (`ADR-0029`, Regel 5)
// ist hier als Freigabe-Bedingung modelliert: ohne die Bestätigung aller
// Consumer gibt die Policy keine Bereinigung frei; das Mindestalter ist die
// zeitbasierte Grenze (`LH-FA-RET-003`). Die Ausführung liegt im
// Run-Retention-Use-Case, nicht in diesem Wertobjekt.
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

// AllowsDeletion meldet, ob die Policy eine Bereinigung freigibt: der
// Change trägt das Mindestalter und alle relevanten Consumer haben die
// Position bestätigt (`LH-FA-RET-004`, `ADR-0029`, Regel 5).
func (p RetentionPolicy) AllowsDeletion(age Duration, allConsumersAcknowledged bool) bool {
	if !allConsumersAcknowledged {
		return false
	}
	return age.Nanos >= p.MinAge.Nanos
}
