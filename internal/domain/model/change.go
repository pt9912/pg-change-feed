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
// Konstanten tragen die Werte aus `SPEC-002`. Die Menge bleibt bei drei
// Werten: ein Bestands-Change eines Backfills (`LH-FA-CAP-009`) ist
// `OperationInsert` mit der Herkunft `ChangeOriginBackfill`, kein vierter
// Operationswert (`ADR-0111`).
type Operation string

const (
	OperationInsert Operation = "INSERT"
	OperationUpdate Operation = "UPDATE"
	OperationDelete Operation = "DELETE"
)

// ChangeOrigin trägt die Herkunft eines Changes (`SPEC-002`, Feld
// `origin`): `wal` für einen über den Replication Stream erfassten Change,
// `backfill` für einen Bestands-Change eines Backfill-Runs
// (`LH-FA-CAP-009`). Die Menge ist geschlossen; die Datenbankspalte
// `cdc.change.origin` trägt keinen CHECK, sie erzwingt allein diese
// Domäne (`ADR-0111` Teilfrage 2). Ein fehlender Wert — die leere
// Zeichenkette, in der Datenbank `NULL` — liest als `wal`
// (`LH-FA-DAT-006`, Boundary).
type ChangeOrigin string

const (
	ChangeOriginWAL      ChangeOrigin = "wal"
	ChangeOriginBackfill ChangeOrigin = "backfill"
)

// NewChangeOrigin legt eine Herkunft aus ihrer Text-Form an: die leere
// Zeichenkette (fehlender Wert, `NULL`) ist `wal`, jeder Wert außerhalb
// von `wal`/`backfill` wird abgelehnt.
func NewChangeOrigin(raw string) (ChangeOrigin, error) {
	switch ChangeOrigin(raw) {
	case "", ChangeOriginWAL:
		return ChangeOriginWAL, nil
	case ChangeOriginBackfill:
		return ChangeOriginBackfill, nil
	default:
		return "", domainerrors.ErrInvalidChangeOrigin
	}
}

// OrDefault liest einen fehlenden Wert (die leere Zeichenkette) als `wal`;
// jeder gesetzte Wert bleibt unverändert.
func (o ChangeOrigin) OrDefault() ChangeOrigin {
	if o == "" {
		return ChangeOriginWAL
	}
	return o
}

// Change trägt einen einzelnen erfassten Change nach `SPEC-002`
// (`cdc.change`). Row Images sind JSON-Bytes; ihre Zusammensetzung bleibt
// beim Adapter-Mapper. Ein fehlendes Bild ist Abwesenheit, kein Fehler
// (`LH-FA-CAP-008`, Boundary) — der Konstruktor fordert deshalb kein Bild.
//
// Origin trägt die Herkunft (`ChangeOrigin`); `NewChange` setzt `wal`,
// `WithOrigin` setzt `backfill` (`LH-FA-CAP-009`). Nur die lesenden
// Wege `cdc.changes` und `GET /changes` tragen das Feld; die drei
// Live-Wege (gRPC, SSE, NATS-Vollinhalt) bilden weiter zehn Felder ab und
// lassen es aus (`SPEC-020`, `SPEC-021`, `SPEC-024`).
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
	Origin        ChangeOrigin
}

// NewChange legt einen Change an und erzwingt die Change-Invarianten:
// nichtleere Kennungen (`LH-FA-DAT-001`), Sequenz mindestens 1
// (`SPEC-002`, eindeutige Sequenz innerhalb der Transaktion),
// gültiger Operationstyp (`LH-FA-DAT-003`) und eine Schema-Version-Referenz
// (`ADR-0029`, Regel 7). Die Herkunft ist `wal` (`ChangeOriginWAL`); eine
// andere setzt `WithOrigin`.
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
		Origin:        ChangeOriginWAL,
	}, nil
}

// WithOrigin liefert den Change mit der übergebenen Herkunft; ein Wert
// außerhalb der geschlossenen Menge (`ChangeOrigin`) wird abgelehnt, der
// Change bleibt dann unverändert.
func (c Change) WithOrigin(origin ChangeOrigin) (Change, error) {
	checked, err := NewChangeOrigin(string(origin))
	if err != nil {
		return c, err
	}
	c.Origin = checked
	return c, nil
}
