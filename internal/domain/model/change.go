package model

import (
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

// ChangeID identifiziert einen Change eindeutig (`LH-FA-DAT-001`); die
// konkrete Form legt der Adapter fest, die Eindeutigkeit ist Domänen-Eigenschaft.
type ChangeID string

// TransactionID identifiziert eine Quelltransaktion (`SPEC-001`, Tabelle
// `cdc.transaction`).
type TransactionID string

// Operation trägt den Operationstyp eines Changes (`LH-FA-DAT-003`); die
// Konstanten tragen die Werte aus `SPEC-002`.
type Operation string

const (
	OperationInsert Operation = "INSERT"
	OperationUpdate Operation = "UPDATE"
	OperationDelete Operation = "DELETE"
)

// Change trägt einen einzelnen erfassten Change nach `SPEC-002`
// (`cdc.change`). Row Images sind JSON-Bytes; ihre Zusammensetzung bleibt
// beim Adapter-Mapper. Ein fehlendes Bild ist Abwesenheit, kein Fehler
// (`LH-FA-CAP-008`, Boundary) — der Konstruktor fordert deshalb kein Bild.
type Change struct {
	ID            ChangeID
	TransactionID TransactionID
	SourceTableID SourceTableID
	Sequence      int64
	Operation     Operation
	OldImage      []byte
	NewImage      []byte
	SchemaVersion SchemaVersionID
}

// NewChange legt einen Change an und erzwingt die Change-Invarianten:
// nichtleere Kennungen (`LH-FA-DAT-001`), Sequenz mindestens 1
// (`SPEC-002`, eindeutige Sequenz innerhalb der Transaktion),
// gültiger Operationstyp (`LH-FA-DAT-003`) und eine Schema-Version-Referenz
// (`ADR-0029`, Regel 7).
func NewChange(id ChangeID, tx TransactionID, table SourceTableID, sequence int64, op Operation, oldImage, newImage []byte, schemaVersion SchemaVersionID) (Change, error) {
	if id == "" || tx == "" || table == "" || schemaVersion == "" {
		return Change{}, domainerrors.ErrEmptyIdentifier
	}
	if sequence < 1 {
		return Change{}, domainerrors.ErrNonPositiveSequence
	}
	if op != OperationInsert && op != OperationUpdate && op != OperationDelete {
		return Change{}, domainerrors.ErrInvalidOperation
	}
	return Change{
		ID:            id,
		TransactionID: tx,
		SourceTableID: table,
		Sequence:      sequence,
		Operation:     op,
		OldImage:      oldImage,
		NewImage:      newImage,
		SchemaVersion: schemaVersion,
	}, nil
}
