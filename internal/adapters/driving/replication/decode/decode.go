// Package decode trägt die `pgoutput`-Dekodierung des
// Replication-Stream-Adapters (`ADR-0008`, `SPEC-010`): die binären
// Output-Plugin-Nachrichten werden in die typisierten Ereignisse dieses
// Pakets übersetzt — BEGIN, COMMIT, Relation, INSERT, UPDATE, DELETE und
// TRUNCATE (`LH-FA-CFG-001.a`), alle übrigen Nachrichtentypen als
// Abwesenheit. Die Protokoll-Details liegen am Treiber (pglogrepl,
// `ADR-0032`); dieser Träger hält die Katalog- und Auflösungs-Arbeit der
// Adapter-Verantwortung (`ADR-0023`): eine Nachricht, die nicht sicher
// interpretierbar ist, endet als sichtbarer Fehler der Klasse `schema`
// statt als stilles Überspringen (`SPEC-008`,
// `LH-QA-REL-001.a` Fehlermodi).
package decode

import (
	stderrors "errors"
	"fmt"
	"time"

	"github.com/jackc/pglogrepl"
)

// ErrSchema trägt die Fehlerklasse `schema` der Dekodierung (`SPEC-008`,
// `ADR-0023`): Dekodierfehler sind sichtbare Fehler, kein stilles
// Überspringen (`LH-QA-REL-001.a` Fehlermodi, `LH-FA-SCH-004.a`). Der
// Aufrufer klassifiziert über `errors.Is`; die technische Ursache bleibt
// über die zweite Wrappung lesbar.
var ErrSchema = stderrors.New("Fehlerklasse schema: pgoutput-Nachricht nicht sicher interpretierbar")

// Event ist ein dekodiertes `pgoutput`-Ereignis. Die Konkretisierungen
// tragen die Nachrichten des Capture-Pfads; nicht relevante
// Nachrichtentypen liefern kein Ereignis.
type Event interface {
	isEvent()
}

// Begin trägt den Transaktionsbeginn (`pgoutput`, Typ `B`).
type Begin struct {
	XID uint32
}

func (Begin) isEvent() {}

// Commit trägt den Commit einer Quelltransaktion (Typ `C`): CommitLSN
// ist die Commit-Position (`SPEC-003`, der Adapter mappt die LSN auf
// den Offset), EndLSN ist das Transaktionsende. CommitTime trägt den
// realen Quell-Commit-Zeitpunkt (`LH-FA-ADM-004`) aus
// `pglogrepl.CommitMessage.CommitTime`; der Mapper übersetzt ihn beim
// Domänen-Commit in `model.TimePoint` (`ADR-0040` — Domain importiert
// `time` nicht).
type Commit struct {
	CommitLSN  uint64
	EndLSN     uint64
	CommitTime time.Time
}

func (Commit) isEvent() {}

// Column ist eine Relation-Spalte; Key markiert eine Spalte der
// Replica-Identity. TypeOID trägt die rohe PostgreSQL-Typ-OID der Spalte
// (`ADR-0015` Folgepflicht, `SPEC-004`) — die Grundlage, auf der der
// Mapper unverändert von kompatibel erweitert unterscheidet
// (`LH-FA-SCH-005`); die Übersetzung dieser OID in eine
// Vergleichsentscheidung trägt dieses Paket nicht.
type Column struct {
	Name    string
	Key     bool
	TypeOID uint32
}

// Relation trägt die Relation-Metadaten einer `pgoutput`-Relation (Typ
// `R`): Schema und Tabellenname identifizieren die Tabelle innerhalb
// der Quelle (`LH-FA-DAT-002`); die Spalten tragen die Namen und die
// Replica-Identity-Zugehörigkeit der Tupel-Werte.
type Relation struct {
	Schema  string
	Name    string
	Columns []Column
}

func (*Relation) isEvent() {}

// QualifiedName trägt den schema-qualifizierten Tabellennamen.
func (r *Relation) QualifiedName() string {
	return r.Schema + "." + r.Name
}

// Operation trägt den Operationstyp einer pgoutput-Änderung; die
// Übersetzung in die Domänen-Operationen (`SPEC-002`) trägt der Mapper.
type Operation uint8

const (
	OpInsert Operation = iota
	OpUpdate
	OpDelete
)

// Change trägt eine INSERT-, UPDATE- oder DELETE-Änderung (Typen `I`,
// `U`, `D`): die Tupel-Werte sind positionsgleich zu Relation.Columns.
// Ein nil-Wert ist Abwesenheit — NULL-Wert oder unverändertes
// TOAST-Stützwert-Platzhalter-Verhalten — und kein Fehler
// (`LH-FA-CAP-008` Boundary). Old trägt den alten Tupel-Stand nur, wenn
// die Quelle ihn gesendet hat (Replica Identity `O` bzw. Schlüssel `K`).
type Change struct {
	Relation  *Relation
	Operation Operation
	Old       []*string
	New       []*string
}

func (Change) isEvent() {}

// Truncate trägt eine TRUNCATE-Meldung (Typ `T`) über die betroffenen
// Relationen; die Nicht-Unterstützung entscheidet der Mapper
// (`LH-FA-CFG-001.a`).
type Truncate struct {
	Relations []*Relation
}

func (Truncate) isEvent() {}

// Decoder dekodiert `pgoutput`-Payloads und trägt den Relation-Katalog
// der laufenden Verbindung: die Relation-Nachrichten laufen vor den
// Änderungen und füllen ihn, die Änderungen lösen ihre Relation darüber
// auf.
type Decoder struct {
	relations map[uint32]*Relation
}

// NewDecoder legt einen Decoder ohne Katalog an.
func NewDecoder() *Decoder {
	return &Decoder{relations: map[uint32]*Relation{}}
}

// Decode dekodiert einen `pgoutput`-Payload in ein Ereignis; nicht
// relevante Nachrichtentypen liefern kein Ereignis. Ein Fehler trägt
// ErrSchema (`SPEC-008`, Klasse `schema`) mit der technischen Ursache.
func (d *Decoder) Decode(payload []byte) (Event, error) {
	message, err := pglogrepl.Parse(payload)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrSchema, err)
	}
	switch message := message.(type) {
	case *pglogrepl.BeginMessage:
		return Begin{XID: message.Xid}, nil
	case *pglogrepl.CommitMessage:
		return Commit{
			CommitLSN:  uint64(message.CommitLSN),
			EndLSN:     uint64(message.TransactionEndLSN),
			CommitTime: message.CommitTime,
		}, nil
	case *pglogrepl.RelationMessage:
		relation := &Relation{
			Schema:  message.Namespace,
			Name:    message.RelationName,
			Columns: relationColumns(message),
		}
		d.relations[message.RelationID] = relation
		return relation, nil
	case *pglogrepl.InsertMessage:
		relation, err := d.lookup(message.RelationID)
		if err != nil {
			return nil, err
		}
		values, err := tupleValues(relation, message.Tuple)
		if err != nil {
			return nil, err
		}
		return Change{
			Relation:  relation,
			Operation: OpInsert,
			New:       values,
		}, nil
	case *pglogrepl.UpdateMessage:
		relation, err := d.lookup(message.RelationID)
		if err != nil {
			return nil, err
		}
		values, err := tupleValues(relation, message.NewTuple)
		if err != nil {
			return nil, err
		}
		change := Change{
			Relation:  relation,
			Operation: OpUpdate,
			New:       values,
		}
		if message.OldTuple != nil {
			oldValues, err := oldTupleValues(relation, message.OldTupleType, message.OldTuple)
			if err != nil {
				return nil, err
			}
			change.Old = oldValues
		}
		return change, nil
	case *pglogrepl.DeleteMessage:
		relation, err := d.lookup(message.RelationID)
		if err != nil {
			return nil, err
		}
		oldValues, err := oldTupleValues(relation, message.OldTupleType, message.OldTuple)
		if err != nil {
			return nil, err
		}
		return Change{
			Relation:  relation,
			Operation: OpDelete,
			Old:       oldValues,
		}, nil
	case *pglogrepl.TruncateMessage:
		relationList := make([]*Relation, 0, len(message.RelationIDs))
		for _, relationID := range message.RelationIDs {
			relation, err := d.lookup(relationID)
			if err != nil {
				return nil, err
			}
			relationList = append(relationList, relation)
		}
		return Truncate{Relations: relationList}, nil
	case *pglogrepl.TypeMessage, *pglogrepl.OriginMessage, *pglogrepl.LogicalDecodingMessage:
		// Nachrichtentypen außerhalb des Capture-Pfads: sie tragen
		// keine Quelländerung und liefern kein Ereignis.
		return nil, nil
	default:
		return nil, fmt.Errorf("%w: unbekannter Nachrichtentyp %q", ErrSchema, message.Type())
	}
}

// lookup löst eine Relation über den Katalog auf; eine Änderung vor der
// zugehörigen Relation-Nachricht ist nicht sicher interpretierbar
// (`LH-FA-SCH-004.a`).
func (d *Decoder) lookup(relationID uint32) (*Relation, error) {
	relation, known := d.relations[relationID]
	if !known {
		return nil, fmt.Errorf("%w: Änderung ohne Relation-Nachricht für Relation %d", ErrSchema, relationID)
	}
	return relation, nil
}

// relationColumns trägt die Spalten einer Relation-Nachricht mit ihren
// Replica-Identity-Markierungen und ihrer Typ-OID; der Protokollwert 1
// markiert eine Schlüssel-Spalte, `DataType` trägt die PostgreSQL-Typ-OID
// (`ADR-0015` Folgepflicht).
func relationColumns(message *pglogrepl.RelationMessage) []Column {
	const columnFlagKey = uint8(1)
	columns := make([]Column, 0, len(message.Columns))
	for _, source := range message.Columns {
		columns = append(columns, Column{Name: source.Name, Key: source.Flags == columnFlagKey, TypeOID: source.DataType})
	}
	return columns
}

// tupleValues trägt die Text-Werte eines vollen Tupels: positionsgleich
// zu den Relation-Spalten. NULL-Werte und unveränderte TOAST-Stützwerte
// sind Abwesenheit (`LH-FA-CAP-008` Boundary); ein Binär-Tupel ist
// nicht sicher interpretierbar, weil die Verbindung im Text-Format
// streamt (`SPEC-008`, Klasse `schema`).
func tupleValues(relation *Relation, tuple *pglogrepl.TupleData) ([]*string, error) {
	if tuple == nil {
		return nil, nil
	}
	if len(tuple.Columns) != len(relation.Columns) {
		return nil, fmt.Errorf("%w: Tupel trägt %d Werte für %d Spalten", ErrSchema, len(tuple.Columns), len(relation.Columns))
	}
	values := make([]*string, len(relation.Columns))
	for i, source := range tuple.Columns {
		switch source.DataType {
		case pglogrepl.TupleDataTypeText:
			value := string(source.Data)
			values[i] = &value
		case pglogrepl.TupleDataTypeNull, pglogrepl.TupleDataTypeToast:
			// Abwesenheit: der Wert trägt NULL oder wurde nicht
			// gesendet (unverändertes TOAST).
		default:
			return nil, fmt.Errorf("%w: Binär-Tupel-Wert in Spalte %q", ErrSchema, relation.Columns[i].Name)
		}
	}
	return values, nil
}

// oldTupleValues trägt die Werte eines alten Tupels: ein `O`-Tupel
// (Replica Identity FULL) umfasst alle Spalten positionsgleich, ein
// `K`-Tupel nur die Schlüssel-Spalten in deren Reihenfolge. Ein nil-
// Ergebnis ist Abwesenheit — die Quelle hat keinen alten Tupel gesendet
// (`LH-FA-CAP-008` Boundary).
func oldTupleValues(relation *Relation, oldTupleType uint8, tuple *pglogrepl.TupleData) ([]*string, error) {
	if tuple == nil {
		return nil, nil
	}
	if oldTupleType == pglogrepl.UpdateMessageTupleTypeOld || oldTupleType == pglogrepl.DeleteMessageTupleTypeOld {
		return tupleValues(relation, tuple)
	}
	if len(tuple.Columns) > len(relation.Columns) {
		return nil, fmt.Errorf("%w: Schlüssel-Tupel trägt %d Werte für %d Spalten", ErrSchema, len(tuple.Columns), len(relation.Columns))
	}
	keyValues := make([]*string, len(relation.Columns))
	keyIndex := 0
	for i, column := range relation.Columns {
		if !column.Key || keyIndex >= len(tuple.Columns) {
			continue
		}
		switch tuple.Columns[keyIndex].DataType {
		case pglogrepl.TupleDataTypeText:
			value := string(tuple.Columns[keyIndex].Data)
			keyValues[i] = &value
		case pglogrepl.TupleDataTypeNull:
			// Abwesenheit.
		default:
			return nil, fmt.Errorf("%w: Binär-Tupel-Wert in Schlüssel-Spalte %q", ErrSchema, column.Name)
		}
		keyIndex++
	}
	return keyValues, nil
}
