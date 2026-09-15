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
//
// Schema und Table tragen die Klartext-Bezeichner der betroffenen Tabelle
// zusätzlich zur opaken SourceTableID: das tabellen-granulare Subjekt des
// NATS-Wecksignals (`ChangeNotificationPort`, `ADR-0056`) und die Antwort
// des lesenden API-Endpunkts (`GET /changes`, `ADR-0081`) adressieren die
// Tabelle in Klartext. Sie sind bewusst keine Konstruktor-Invariante: der
// füllende Pfad kennt beide bereits vor dem Aufruf von `NewChange` und
// setzt sie am Ergebnis — der Driving-Adapter-Mapper aus dem
// Replication-Ereignis, die Rekonstruktion eines persistierten Change
// (`postgresstorage/mapper.ToChange`) aus dem Join auf `cdc.source_table`.
type Change struct {
	ID            ChangeID
	TransactionID TransactionID
	SourceTableID SourceTableID
	Sequence      int64
	Operation     Operation
	OldImage      []byte
	NewImage      []byte
	SchemaVersion SchemaVersionID
	Schema        string
	Table         string
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
