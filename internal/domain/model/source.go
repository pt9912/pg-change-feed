// Package model trägt die Domänenobjekte des CDC-Modells
// (`ARC-001`): pur, ohne Treiber, Frameworks oder PostgreSQL-Typen
// (`ADR-0004`, `ADR-0039`). Die Invarianten (`ADR-0029`) tragen die
// Konstruktoren und Methoden; die exportierten Felder lassen
// Struktur-Literale zu, die diese Prüfung umgehen — der
// Konstruktor-/Methoden-Pfad ist der einzige geprüfte.
package model

import (
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

// SourceID identifiziert eine erfasste Quelle (`SPEC-001`, Tabelle
// `cdc.source`).
type SourceID string

// Source ist eine erfasste Quelle.
type Source struct {
	ID   SourceID
	Name string
}

// NewSource legt eine Quelle an und verlangt nichtleere Kennung und Name.
func NewSource(id SourceID, name string) (Source, error) {
	if id == "" {
		return Source{}, domainerrors.ErrEmptyIdentifier
	}
	if name == "" {
		return Source{}, domainerrors.ErrEmptyIdentifier
	}
	return Source{ID: id, Name: name}, nil
}
