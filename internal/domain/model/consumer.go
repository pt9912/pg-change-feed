package model

import (
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

// ConsumerID identifiziert einen registrierten Consumer (`SPEC-001`,
// Tabelle `cdc.consumer`).
type ConsumerID string

// Consumer ist ein registrierter Abnehmer der Changes (`ADR-0013`).
type Consumer struct {
	ID   ConsumerID
	Name string
}

// NewConsumer legt einen Consumer an und verlangt nichtleere Kennung und
// Name.
func NewConsumer(id ConsumerID, name string) (Consumer, error) {
	if id == "" || name == "" {
		return Consumer{}, domainerrors.ErrEmptyIdentifier
	}
	return Consumer{ID: id, Name: name}, nil
}

// ConsumerPosition trägt die zuletzt bestätigte Position eines Consumers
// (`SPEC-001`, Tabelle `cdc.consumer_position`; `LH-FA-CON-003`). Der
// Nullwert der Position steht für „noch nichts bestätigt“; die Quelle
// trägt die Position selbst (`SPEC-003`) — die erste Bestätigung bindet
// sie, jede weitere muss dieselbe Quelle wiederholen.
type ConsumerPosition struct {
	ConsumerID ConsumerID
	Position   SourcePosition
}

// NewConsumerPosition legt eine Consumer-Position ohne Bestätigung an;
// die Kennung ist nichtleer.
func NewConsumerPosition(consumer ConsumerID) (ConsumerPosition, error) {
	if consumer == "" {
		return ConsumerPosition{}, domainerrors.ErrEmptyIdentifier
	}
	return ConsumerPosition{ConsumerID: consumer}, nil
}

// Acknowledged meldet, dass der Consumer eine Position bestätigt hat.
func (c ConsumerPosition) Acknowledged() bool {
	return !c.Position.IsZero()
}

// Advance bestätigt eine Position und trägt den neuen Fortschritt. Der
// Consumer-ACK verläuft regulär nur vorwärts (`ADR-0029`, Regel 2); die
// Wiederholung derselben Position ist idempotent (`LH-FA-CON-004`), eine
// frühere Position ist ein Fehler, und eine bestätigte Quelle bleibt
// gebunden — die Position einer anderen Quelle ist keine
// Bestätigungs-Größe (`ADR-0005`).
func (c ConsumerPosition) Advance(position SourcePosition) (ConsumerPosition, error) {
	if c.Acknowledged() {
		if position.SourceID != c.Position.SourceID {
			return ConsumerPosition{}, domainerrors.ErrSourceMismatch
		}
		if position.Before(c.Position) {
			return ConsumerPosition{}, domainerrors.ErrPositionRegression
		}
	}
	return ConsumerPosition{ConsumerID: c.ConsumerID, Position: position}, nil
}
