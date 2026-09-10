package model

import (
	stderrors "errors"
	"testing"

	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

// LH-FA-ADM-003, Happy Path: jede der sieben stabilen Kategorien aus
// `ADR-0023` ist eine gültige Fehlerklasse.
func TestNewErrorClassAcceptsTheSevenADR0023Categories(t *testing.T) {
	cases := []ErrorClass{
		ErrorClassTransient, ErrorClassConfiguration, ErrorClassPermission,
		ErrorClassSchema, ErrorClassStorage, ErrorClassReplication, ErrorClassInternal,
	}
	for _, want := range cases {
		t.Run(string(want), func(t *testing.T) {
			got, err := NewErrorClass(string(want))
			if err != nil {
				t.Fatalf("NewErrorClass(%q): %v", want, err)
			}
			if got != want {
				t.Fatalf("NewErrorClass(%q) = %q, wollen %q", want, got, want)
			}
		})
	}
}

// LH-FA-ADM-003, Negative: eine leere oder unbekannte Klasse verletzt die
// Invariante — die sieben Kategorien sind eine geschlossene Menge, kein
// Freitext.
func TestNewErrorClassRejectsUnknownClass(t *testing.T) {
	cases := []string{"", "unbekannt", "STORAGE", " storage"}
	for _, raw := range cases {
		t.Run(raw, func(t *testing.T) {
			if _, err := NewErrorClass(raw); !stderrors.Is(err, domainerrors.ErrInvalidErrorClass) {
				t.Fatalf("NewErrorClass(%q) = %v, wollen %v", raw, err, domainerrors.ErrInvalidErrorClass)
			}
		})
	}
}
