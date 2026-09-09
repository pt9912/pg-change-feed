// Package mapper trägt die Zeilen-Übersetzung des PostgresChangeStore-
// Adapters (Paketstruktur je `ADR-0042`, der die Struktur-Regeln des
// abgelösten `ADR-0039` als Rest fortgilt): die Row-Typen sind die
// Zeilenabbilder der `SPEC-001`/`SPEC-002`-Tabellen, die Funktionen
// übersetzen sie gegen das Domänenmodell — PostgreSQL-Berührungen bleiben
// im Adapter (`ADR-0032`).
package mapper

import (
	stderrors "errors"
	"math"

	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// ErrPositionOutOfRange: die Commit-Position liegt außerhalb des
// bigint-Bereichs der Spalte — eine Grenze der PostgreSQL-Abbildung
// (`SPEC-003`: der Adapter mappt die Positionsgröße), keine
// Domänen-Invariante.
var ErrPositionOutOfRange = stderrors.New("Position liegt außerhalb des bigint-Bereichs")

// TransactionRow trägt eine Zeile aus `cdc.transaction` (`SPEC-001`).
type TransactionRow struct {
	TransactionID  string
	SourceID       string
	CommitPosition int64
}

// ChangeRow trägt eine Zeile aus `cdc.change` (`SPEC-002`); die Row Images
// sind JSON-Bytes, NULL liest sich als nil (Abwesenheit).
type ChangeRow struct {
	ChangeID      string
	TransactionID string
	SourceTableID string
	Sequence      int64
	Operation     string
	OldData       []byte
	NewData       []byte
	SchemaVersion string
}

// NewTransactionRow trägt die `cdc.transaction`-Zeile einer committed
// Quelltransaktion. Die Position muss in den bigint-Bereich der Spalte
// passen (`SPEC-003`, PostgreSQL-Abbildung).
func NewTransactionRow(transaction *model.ChangeTransaction, position model.SourcePosition) (TransactionRow, error) {
	if position.Offset > math.MaxInt64 {
		return TransactionRow{}, ErrPositionOutOfRange
	}
	return TransactionRow{
		TransactionID:  string(transaction.ID),
		SourceID:       string(transaction.SourceID),
		CommitPosition: int64(position.Offset),
	}, nil
}

// NewChangeRows trägt die `cdc.change`-Zeilen der Changes; ein fehlendes
// Bild trägt die Zeile als nil (NULL, Abwesenheit).
func NewChangeRows(changes []model.Change) ([]ChangeRow, error) {
	rows := make([]ChangeRow, 0, len(changes))
	for _, change := range changes {
		rows = append(rows, ChangeRow{
			ChangeID:      string(change.ID),
			TransactionID: string(change.TransactionID),
			SourceTableID: string(change.SourceTableID),
			Sequence:      change.Sequence,
			Operation:     string(change.Operation),
			OldData:       change.OldImage,
			NewData:       change.NewImage,
			SchemaVersion: string(change.SchemaVersion),
		})
	}
	return rows, nil
}

// JSONImage trägt ein Row Image in seine Parameter-Form für die
// `jsonb`-Spalte: nil bleibt NULL (Abwesenheit, kein Fehler —
// `LH-FA-CAP-008`), ein leeres Byte-Slice ist kein gültiges JSON und liest
// sich ebenfalls als NULL; ein nicht-leeres Bild geht als JSON-Text.
func JSONImage(raw []byte) any {
	if len(raw) == 0 {
		return nil
	}
	return string(raw)
}

// ToPosition trägt die Quellposition aus einer `cdc.transaction`-Zeile
// (`SPEC-003`): Quelle und Commit-Position; die Spalte trägt bigint, ihre
// CHECK-Kante hält sie positiv, der Wert passt damit in den
// Domänen-Offset — die Bereichs-Grenze des Schreibpfads liegt allein in
// NewTransactionRow.
func ToPosition(source string, commitPosition int64) (model.SourcePosition, error) {
	if commitPosition < 1 {
		return model.SourcePosition{}, domainerrors.ErrInvalidPosition
	}
	return model.NewSourcePosition(model.SourceID(source), uint64(commitPosition))
}

// ToChange trägt den Change aus einer `cdc.change`-Zeile; die Domänen-
// Konstruktoren prüfen die Zeile über die Change-Invarianten (`ADR-0029`).
func ToChange(row ChangeRow) (model.Change, error) {
	return model.NewChange(
		model.ChangeID(row.ChangeID),
		model.TransactionID(row.TransactionID),
		model.SourceTableID(row.SourceTableID),
		row.Sequence,
		model.Operation(row.Operation),
		row.OldData,
		row.NewData,
		model.SchemaVersionID(row.SchemaVersion),
	)
}
