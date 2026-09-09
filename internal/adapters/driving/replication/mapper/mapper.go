// Package mapper übersetzt die dekodierten `pgoutput`-Ereignisse des
// Replication-Stream-Adapters in Aufrufe des `CaptureInboundPort`
// (`ARC-003`, `ADR-0006`, `ADR-0042`): der Träger baut aus BEGIN/Commit
// und den Änderungen eine vollständige, committed Quelltransaktion
// (`LH-FA-CAP-006.a`) und meldet sie als CaptureCommand. Er entscheidet
// nicht über Persistenz oder ACK — Persist-before-ACK liegt am Capture
// Use Case (`ADR-0027`, `LH-QA-REL-001.a`).
package mapper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"

	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/decode"
	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
)

// ErrTruncateUnsupported trägt die Nicht-Unterstützung von TRUNCATE im
// MVP (`LH-FA-CFG-001.a`): die Operation wird erkannt und als sichtbarer
// Fehler der Klasse `schema` behandelt, nicht still übersprungen
// (`SPEC-008`, Lastenheft §5 Out-of-Scope).
var ErrTruncateUnsupported = errors.New("Fehlerklasse schema: TRUNCATE wird im MVP nicht unterstützt")

// ErrChangeWithoutBegin trägt einen Stream-Vertragsverstoß: eine
// Änderung vor dem zugehörigen BEGIN ist nicht interpretierbar; die
// Klasse ist `replication` (`SPEC-008`), weil sie eine Störung der
// Stream-Ordnung trägt.
var ErrChangeWithoutBegin = errors.New("Fehlerklasse replication: Änderung ohne offene Quelltransaktion")

// ErrCommitWithoutBegin trägt denselben Vertragsverstoß für einen Commit
// ohne zugehöriges BEGIN.
var ErrCommitWithoutBegin = errors.New("Fehlerklasse replication: Commit ohne offene Quelltransaktion")

// TableBinding trägt die am Port getragenen Kennungen einer aktivierten
// Tabelle (`SPEC-001`): die Tabelle und die Schema-Version, die die
// Changes dieser Tabelle referenzieren (`LH-FA-SCH-005`). Die
// Schema-Version liegt in diesem Slice statisch bei der Konfiguration;
// die Relation-Metadaten-Übersetzung der Schema-Evolution trägt
// LH-FA-SCH-004.a über den Metadata-Pfad.
type TableBinding struct {
	TableID       model.SourceTableID
	SchemaVersion model.SchemaVersionID
}

// Assembler baut aus den dekodierten Ereignissen committed
// Quelltransaktionen. Ohne Streaming-Option (`ADR-0021`, Proposed)
// serialisiert der Stream die Quelltransaktionen — der Träger hält
// genau eine offene Transaktion.
type Assembler struct {
	source model.SourceID
	tables map[string]TableBinding
	open   *openTransaction
}

// openTransaction trägt die offene Transaktion und den Sequenz-Zähler
// ihrer Änderungen: die Sequenz ist innerhalb der Transaktion eindeutig
// und trägt die Ausführungsreihenfolge (`LH-FA-DAT-004`).
type openTransaction struct {
	tx       *model.ChangeTransaction
	sequence int64
}

// NewAssembler legt den Übersetzer für eine Quelle an; die Aktivierungen
// tragen die qualifizierten Tabellennamen (`schema.table`) als
// Schlüssel. Quelle ohne Kennung und Bindungen ohne Kennungen enden über
// den Domänen-Fehler der leeren Kennung (`ADR-0029`).
func NewAssembler(source model.SourceID, tables map[string]TableBinding) (*Assembler, error) {
	if source == "" {
		return nil, fmt.Errorf("Quelle ohne Kennung: %w", domainerrors.ErrEmptyIdentifier)
	}
	for qualified, binding := range tables {
		if qualified == "" || binding.TableID == "" || binding.SchemaVersion == "" {
			return nil, fmt.Errorf("Tabellen-Bindung %q trägt eine leere Kennung: %w", qualified, domainerrors.ErrEmptyIdentifier)
		}
	}
	return &Assembler{source: source, tables: tables}, nil
}

// Consume übersetzt ein Ereignis in höchstens einen Capture-Aufruf:
// BEGIN öffnet, Änderungen hängen an, COMMIT bringt die Transaktion an
// ihre Commit-Position (`SPEC-001`, `cdc.transaction`) und meldet sie
// konsumierbar (`LH-FA-CAP-006`). Änderungen an nicht aktivierten
// Tabellen fließen nicht in die Transaktion — CDC erfasst nur
// aktivierte Tabellen (`LH-FA-CFG-001`).
func (a *Assembler) Consume(event decode.Event) (*inbound.CaptureCommand, error) {
	switch event := event.(type) {
	case decode.Begin:
		tx, err := model.NewOpenTransaction(model.TransactionID(strconv.FormatUint(uint64(event.XID), 10)), a.source)
		if err != nil {
			return nil, err
		}
		a.open = &openTransaction{tx: tx, sequence: 0}
		return nil, nil
	case decode.Commit:
		if a.open == nil {
			return nil, ErrCommitWithoutBegin
		}
		position, err := model.NewSourcePosition(a.source, event.CommitLSN)
		if err != nil {
			return nil, err
		}
		if err := a.open.tx.Commit(position); err != nil {
			return nil, err
		}
		command := &inbound.CaptureCommand{Transaction: a.open.tx}
		a.open = nil
		return command, nil
	case decode.Change:
		if a.open == nil {
			return nil, ErrChangeWithoutBegin
		}
		change, err := a.change(event)
		if err != nil {
			return nil, err
		}
		if change == nil {
			// Nicht aktivierte Tabelle: keine Erfassung
			// (`LH-FA-CFG-001`), kein Change.
			return nil, nil
		}
		if err := a.open.tx.AppendChange(*change); err != nil {
			return nil, err
		}
		return nil, nil
	case decode.Truncate:
		return nil, fmt.Errorf("%w: TRUNCATE an %s", ErrTruncateUnsupported, qualifiedNames(event.Relations))
	default:
		// Relation-Ereignisse tragen nur den Katalog der Dekodierung;
		// sie erzeugen keinen Change.
		return nil, nil
	}
}

// change übersetzt eine Änderung in einen Domänen-Change: Kennung aus
// Transaktion und Sequenz, Operation, Row Images als JSON
// (`SPEC-002`), Tabelle und Schema-Version über die Aktivierung. Die
// Sequenz trägt die Anhang-Reihenfolge und startet je Transaktion bei 1
// (`SPEC-002`).
func (a *Assembler) change(event decode.Change) (*model.Change, error) {
	binding, activated := a.tables[event.Relation.QualifiedName()]
	if !activated {
		return nil, nil
	}
	a.open.sequence++
	sequence := a.open.sequence

	var operation model.Operation
	switch event.Operation {
	case decode.OpInsert:
		operation = model.OperationInsert
	case decode.OpUpdate:
		operation = model.OperationUpdate
	case decode.OpDelete:
		operation = model.OperationDelete
	default:
		return nil, fmt.Errorf("%w: unbekannte Operation %d", domainerrors.ErrInvalidOperation, event.Operation)
	}

	newImage, err := rowImage(event.Relation, event.New)
	if err != nil {
		return nil, err
	}
	oldImage, err := rowImage(event.Relation, event.Old)
	if err != nil {
		return nil, err
	}

	change, err := model.NewChange(
		model.ChangeID(fmt.Sprintf("%s-%d", a.open.tx.ID, sequence)),
		a.open.tx.ID,
		binding.TableID,
		sequence,
		operation,
		oldImage,
		newImage,
		binding.SchemaVersion,
	)
	if err != nil {
		return nil, err
	}
	return &change, nil
}

// rowImage trägt das JSON-Row-Image (`SPEC-002`, `ADR-0016`) einer
// Änderung: ein JSON-Objekt über die gesendeten Spalten-Werte in
// Relation-Reihenfolge. Werte sind JSON-Strings — der Text-Stand der
// Quelle geht unverändert in das Bild, ohne Typ-Interpretation; ein
// nil-Wert trägt NULL oder unverändertes TOAST und ist Abwesenheit
// (`LH-FA-CAP-008` Boundary).
func rowImage(relation *decode.Relation, values []*string) ([]byte, error) {
	if values == nil {
		return nil, nil
	}
	var image bytes.Buffer
	image.WriteByte('{')
	first := true
	for i, column := range relation.Columns {
		if i >= len(values) || values[i] == nil {
			continue
		}
		name, err := json.Marshal(column.Name)
		if err != nil {
			return nil, err
		}
		value, err := json.Marshal(*values[i])
		if err != nil {
			return nil, err
		}
		if !first {
			image.WriteByte(',')
		}
		first = false
		image.Write(name)
		image.WriteByte(':')
		image.Write(value)
	}
	image.WriteByte('}')
	return image.Bytes(), nil
}

// qualifiedNames trägt die qualifizierten Namen einer Relation-Liste.
func qualifiedNames(relations []*decode.Relation) string {
	var names bytes.Buffer
	for i, relation := range relations {
		if i > 0 {
			names.WriteString(", ")
		}
		names.WriteString(relation.QualifiedName())
	}
	return names.String()
}