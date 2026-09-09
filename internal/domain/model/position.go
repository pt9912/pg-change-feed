package model

import (
	"cmp"
	"strings"

	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

// SourcePosition trägt die technologieunabhängige Quellposition
// (`SPEC-003`): der Quelladapter mappt seine Positionsgröße — beim
// PostgreSQL-Adapter die LSN — auf den Offset (`ADR-0005`). Die fachliche
// Ordnung gilt innerhalb einer Quelle (`LH-FA-DAT-004`); der Vergleich über
// Quellen ist ein stabiler Sortier-Schlüssel (`LH-FA-REA-004`), keine
// fachliche Aussage über Quellen hinweg.
type SourcePosition struct {
	SourceID SourceID
	Offset   uint64
}

// NewSourcePosition legt eine Position an und verlangt eine nichtleere
// Quelle und einen Offset größer als 0.
func NewSourcePosition(source SourceID, offset uint64) (SourcePosition, error) {
	if source == "" {
		return SourcePosition{}, domainerrors.ErrEmptyIdentifier
	}
	if offset == 0 {
		return SourcePosition{}, domainerrors.ErrInvalidPosition
	}
	return SourcePosition{SourceID: source, Offset: offset}, nil
}

// IsZero meldet den Nullwert; er steht für „keine Position bestätigt“.
func (p SourcePosition) IsZero() bool {
	return p.SourceID == "" && p.Offset == 0
}

// Compare ordnet zwei Positionen deterministisch: gleiche Quelle nach
// Offset, verschiedene Quellen nach Quell-Kennung als Tiebreak.
func (p SourcePosition) Compare(other SourcePosition) int {
	if c := strings.Compare(string(p.SourceID), string(other.SourceID)); c != 0 {
		return c
	}
	return cmp.Compare(p.Offset, other.Offset)
}

// Before meldet, dass p vor other liegt (deterministische Ordnung,
// Compare).
func (p SourcePosition) Before(other SourcePosition) bool {
	return p.Compare(other) < 0
}

// After meldet, dass p nach other liegt (deterministische Ordnung,
// Compare).
func (p SourcePosition) After(other SourcePosition) bool {
	return p.Compare(other) > 0
}
