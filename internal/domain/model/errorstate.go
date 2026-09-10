package model

import (
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

// ErrorClass trägt eine der sieben stabilen Fehlerkategorien aus
// `SPEC-008`/`ADR-0023` (`transient`, `configuration`, `permission`,
// `schema`, `storage`, `replication`, `internal`): der Capture-Prozess
// meldet seinen zuletzt beobachteten Fehlerzustand darüber, erkennbar und
// unterscheidbar vom Normalbetrieb (`LH-FA-ADM-003`, `LH-QA-REL-003`).
type ErrorClass string

const (
	ErrorClassTransient     ErrorClass = "transient"
	ErrorClassConfiguration ErrorClass = "configuration"
	ErrorClassPermission    ErrorClass = "permission"
	ErrorClassSchema        ErrorClass = "schema"
	ErrorClassStorage       ErrorClass = "storage"
	ErrorClassReplication   ErrorClass = "replication"
	ErrorClassInternal      ErrorClass = "internal"
)

// NewErrorClass validiert gegen die geschlossene Menge aus `ADR-0023`: eine
// leere oder unbekannte Klasse endet über den Domänenfehler
// `ErrInvalidErrorClass` — der Aufrufer (Adapter-Grenze) klassifiziert nie
// frei, sondern immer gegen diese sieben Kategorien.
func NewErrorClass(raw string) (ErrorClass, error) {
	class := ErrorClass(raw)
	switch class {
	case ErrorClassTransient, ErrorClassConfiguration, ErrorClassPermission,
		ErrorClassSchema, ErrorClassStorage, ErrorClassReplication, ErrorClassInternal:
		return class, nil
	default:
		return "", domainerrors.ErrInvalidErrorClass
	}
}
