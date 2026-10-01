package model

import (
	"fmt"

	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

// ChangeID identifiziert einen Change eindeutig (`LH-FA-DAT-001`); die
// Eindeutigkeit ist Domänen-Eigenschaft, die Bildungsregel der erfassenden
// Pfade trägt `ChangeIDFor`.
type ChangeID string

// ChangeIDFor bildet die Kennung eines Changes aus der Kennung seiner
// Transaktion und seiner Sequenz: `<Transaktions-ID>-<Sequenz>`
// (`ADR-0111` Teilfrage 6). Der WAL-Pfad und der Backfill-Pfad rufen diese
// eine Funktion; sie ist eindeutig, solange die Transaktions-Kennung es ist
// und die Sequenz innerhalb der Transaktion eindeutig ist.
func ChangeIDFor(tx TransactionID, sequence int64) ChangeID {
	return ChangeID(fmt.Sprintf("%s-%d", tx, sequence))
}

// TransactionID identifiziert eine Quelltransaktion (`SPEC-001`, Tabelle
// `cdc.transaction`).
type TransactionID string

// Operation trägt den Operationstyp eines Changes (`LH-FA-DAT-003`); die
// Konstanten tragen die Werte des Schemas. Die Menge bleibt bei drei
// Werten: ein Bestands-Change eines Backfills ist
// `OperationInsert` mit der Herkunft `ChangeOriginBackfill`, kein vierter
// Operationswert.
type Operation string

const (
	OperationInsert Operation = "INSERT"
	OperationUpdate Operation = "UPDATE"
	OperationDelete Operation = "DELETE"
)

// ChangeOrigin trägt die Herkunft eines Changes (`LH-FA-CAP-009`, Feld
// `origin`): `wal` für einen über den Replication Stream erfassten Change,
// `backfill` für einen Bestands-Change eines Backfill-Runs.
// Die Menge ist geschlossen; die Datenbankspalte
// `cdc.change.origin` trägt keinen CHECK, sie erzwingt allein diese
// Domäne. Ein fehlender Wert — die leere
// Zeichenkette, in der Datenbank `NULL` — liest als `wal`
// (Boundary).
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
// (Boundary) — der Konstruktor fordert deshalb kein Bild.
//
// Origin trägt die Herkunft (`ChangeOrigin`); `NewChange` setzt `wal`,
// `WithOrigin` setzt `backfill`. Nur die lesenden
// Wege `cdc.changes` und `GET /changes` tragen das Feld; die drei
// Live-Wege (gRPC, SSE, NATS-Vollinhalt) bilden weiter zehn Felder ab und
// lassen es aus.
//
// RouteTarget trägt das Zustellziel (`RouteTarget`); die leere Zeichenkette
// ist „kein Ziel" (in der Datenbank `NULL`). `NewChange` lässt es leer,
// `WithRouteTarget` setzt es. Das Ziel ist nicht Teil der Nachrichten der
// Live-Wege; gelesen wird es über die View `cdc.changes`.
//
// Schema und Table tragen die Klartext-Bezeichner der betroffenen Tabelle
// zusätzlich zur opaken SourceTableID: das tabellen-granulare Subjekt des
// NATS-Wecksignals (`ChangeNotificationPort`) und die Antwort
// des lesenden API-Endpunkts (`GET /changes`) adressieren die
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
	RouteTarget   RouteTarget
}

// NewChange legt einen Change an und erzwingt die Change-Invarianten:
// nichtleere Kennungen (`LH-FA-DAT-001`), Sequenz mindestens 1
// (eindeutige Sequenz innerhalb der Transaktion),
// gültiger Operationstyp und eine referenzierte Schema-Version.
// Die Herkunft ist `wal` (`ChangeOriginWAL`); eine
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

// WithRouteTarget liefert den Change mit dem übergebenen Zustellziel; die
// leere Zeichenkette ist „kein Ziel", jeder andere Wert außerhalb des
// Alphabets (`IsValidRouteTarget`) wird abgelehnt, der Change bleibt dann
// unverändert.
func (c Change) WithRouteTarget(target RouteTarget) (Change, error) {
	if target != "" {
		checked, err := NewRouteTarget(string(target))
		if err != nil {
			return c, err
		}
		target = checked
	}
	c.RouteTarget = target
	return c, nil
}

// MatchesFilter prüft den Change gegen ein optionales, unabhängig
// setzbares `schema`/`table`-Filterpaar (`ADR-0133`): ein leeres Feld
// trägt keinen Filter auf dieser Dimension, beide leer lässt jeden Change
// passieren. Der gRPC- und der SSE-Stream-Handler rufen diese eine
// Funktion, damit ein gefilterter Change nie zwischen den beiden
// Zugriffswegen auseinanderläuft.
func (c Change) MatchesFilter(schema, table string) bool {
	if schema != "" && c.Schema != schema {
		return false
	}
	if table != "" && c.Table != table {
		return false
	}
	return true
}
