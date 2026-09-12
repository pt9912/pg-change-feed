package model

import (
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

// ChangeTransaction bündelt die Changes einer Quelltransaktion
// (`LH-FA-CAP-005`). Eine Transaktion entsteht offen und wird beim Commit
// an ihre Commit-Position gebracht (`SPEC-001`, `cdc.transaction`):
// zurückgerollte Transaktionen erreichen den Commit nicht und erzeugen
// keine committed Changes (`ADR-0029`, Regel 4). Regel 3 („offene
// Transaktionen sind nicht konsumierbar") trägt `Changes()`: an einer
// offenen Transaktion liefert sie einen Fehler; konsumierbar ist die
// Transaktion erst nach dem Commit (`LH-FA-CAP-006`).
type ChangeTransaction struct {
	ID                TransactionID
	SourceID          SourceID
	commitPosition    SourcePosition
	sourceCommittedAt TimePoint
	committed         bool
	changes           []Change
	usedSequences     map[int64]struct{}
}

// NewOpenTransaction legt eine offene Transaktion an; ihr Commit-Position
// fehlt noch.
func NewOpenTransaction(id TransactionID, source SourceID) (*ChangeTransaction, error) {
	if id == "" || source == "" {
		return nil, domainerrors.ErrEmptyIdentifier
	}
	return &ChangeTransaction{
		ID:            id,
		SourceID:      source,
		usedSequences: map[int64]struct{}{},
	}, nil
}

// IsCommitted meldet, dass die Transaktion einen Commit trägt und ihre
// Changes damit konsumierbar sind (`LH-FA-CAP-006`).
func (t *ChangeTransaction) IsCommitted() bool {
	return t.committed
}

// CommitPosition trägt die Commit-Position und meldet über das zweite
// Ergebnis, ob die Transaktion committed ist (`SPEC-001`, Commit-Position).
func (t *ChangeTransaction) CommitPosition() (SourcePosition, bool) {
	return t.commitPosition, t.committed
}

// SourceCommittedAt trägt den realen Quell-Commit-Zeitpunkt
// (`LH-FA-ADM-004`) und meldet über das zweite Ergebnis, ob die
// Transaktion committed ist — analog zu `CommitPosition`.
func (t *ChangeTransaction) SourceCommittedAt() (TimePoint, bool) {
	return t.sourceCommittedAt, t.committed
}

// Commit bringt die Transaktion an ihre Commit-Position — einmalig, und
// nur an einer Position der eigenen Quelle (`LH-FA-DAT-004`).
// sourceCommittedAt trägt den realen Quell-Commit-Zeitpunkt
// (`LH-FA-ADM-004`), abrufbar über `SourceCommittedAt`.
func (t *ChangeTransaction) Commit(position SourcePosition, sourceCommittedAt TimePoint) error {
	if t.committed {
		return domainerrors.ErrTransactionAlreadyCommitted
	}
	if position.SourceID != t.SourceID {
		return domainerrors.ErrSourceMismatch
	}
	t.commitPosition = position
	t.sourceCommittedAt = sourceCommittedAt
	t.committed = true
	return nil
}

// AppendChange hängt einen Change an und erzwingt die Zuordnungs- und
// Reihenfolge-Invariante (`ADR-0029`, Regel 6): der Change gehört zu
// dieser Transaktion, und seine Sequenz ist innerhalb der Transaktion
// eindeutig. An einer committed Transaktion hängt kein Change mehr.
func (t *ChangeTransaction) AppendChange(change Change) error {
	if t.committed {
		return domainerrors.ErrTransactionAlreadyCommitted
	}
	if change.TransactionID != t.ID {
		return domainerrors.ErrTransactionMismatch
	}
	if _, used := t.usedSequences[change.Sequence]; used {
		return domainerrors.ErrDuplicateSequence
	}
	t.usedSequences[change.Sequence] = struct{}{}
	t.changes = append(t.changes, change)
	return nil
}

// Changes trägt die Changes der committed Transaktion in Anhang-Reihenfolge;
// die Reihenfolge innerhalb der Transaktion liegt in der Sequenz
// (`LH-FA-DAT-004`, Boundary). An einer offenen Transaktion liefert sie
// einen Fehler — offene Transaktionen sind nicht konsumierbar
// (`ADR-0029`, Regel 3).
func (t *ChangeTransaction) Changes() ([]Change, error) {
	if !t.committed {
		return nil, domainerrors.ErrTransactionNotCommitted
	}
	return append([]Change(nil), t.changes...), nil
}
