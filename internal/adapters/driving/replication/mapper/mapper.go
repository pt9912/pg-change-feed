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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"

	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/decode"
	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
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

// ErrBeginWithoutCommit trägt denselben Vertragsverstoß für ein BEGIN bei
// bereits offener Transaktion: ein zweites BEGIN überschreibt die offene
// Transaktion samt ihrer gesammelten Changes nicht still — der Stream
// serialisiert die Quelltransaktionen (`ADR-0021`, Proposed); ein
// Doppel-BEGIN ist eine Störung der Stream-Ordnung und endet sichtbar
// (`LH-QA-REL-001.a` Fehlermodi, `SPEC-008` Klasse `replication`).
var ErrBeginWithoutCommit = errors.New("Fehlerklasse replication: BEGIN während offener Quelltransaktion")

// ErrIncompatibleSchemaChange trägt jede nicht sicher als Obermenge
// erkennbare Relation-Änderung (`relationOther`: Spalte entfernt, Typ
// einer bestehenden Spalte geändert, Spalte umbenannt) als sichtbaren
// Fehler der Klasse `schema` (`LH-FA-SCH-004.a`, `SPEC-008`,
// `observeRelation`) — der Erfassungspfad endet darüber, statt die
// Änderung stillschweigend zu übernehmen.
var ErrIncompatibleSchemaChange = errors.New("Fehlerklasse schema: Relation-Änderung nicht sicher als Obermenge interpretierbar")

// TableBinding trägt die am Port getragenen Kennungen einer aktivierten
// Tabelle (`SPEC-001`): die Tabelle und die Schema-Version, die die
// Changes dieser Tabelle referenzieren (`LH-FA-SCH-005`). Die
// Schema-Version liegt initial bei der Konfiguration und wird von
// `Assembler.Consume` bei einer real erkannten kompatiblen Erweiterung
// aktualisiert (`ADR-0015` Folgepflicht) — die Erkennung nicht sicher
// interpretierbarer Änderungen als Fehlerklasse `schema` trägt
// `LH-FA-SCH-004.a` über denselben Metadata-Pfad.
type TableBinding struct {
	TableID       model.SourceTableID
	SchemaVersion model.SchemaVersionID
}

// Assembler baut aus den dekodierten Ereignissen committed
// Quelltransaktionen. Ohne Streaming-Option (`ADR-0021`, Proposed)
// serialisiert der Stream die Quelltransaktionen — der Träger hält
// genau eine offene Transaktion. `schemaStore` trägt die dynamische
// Re-Versionierung (`ADR-0015` Folgepflicht, `LH-FA-SCH-005`); ohne ihn
// (`nil`) bleibt eine Relation-Nachricht wirkungslos — der Regelfall für
// Tests, die diesen Pfad nicht prüfen, die reale Verdrahtung übergibt
// immer eine Instanz (`internal/bootstrap/wiring.go`).
type Assembler struct {
	source      model.SourceID
	tables      map[string]TableBinding
	schemaStore outbound.SchemaStorePort
	open        *openTransaction
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
// den Domänen-Fehler der leeren Kennung (`ADR-0029`). `schemaStore` trägt
// die dynamische Re-Versionierung (`Consume`, `ADR-0015` Folgepflicht);
// `nil` ist zulässig — die Relation-Behandlung bleibt dann wirkungslos.
// Die Aktivierungen gehen kopiert in den Assembler ein: eine kompatible
// Erweiterung schreibt die aktualisierte Bindung in die Kopie
// (`observeRelation`), nicht in die Aufrufer-Map.
func NewAssembler(source model.SourceID, tables map[string]TableBinding, schemaStore outbound.SchemaStorePort) (*Assembler, error) {
	if source == "" {
		return nil, fmt.Errorf("Quelle ohne Kennung: %w", domainerrors.ErrEmptyIdentifier)
	}
	copiedTables := make(map[string]TableBinding, len(tables))
	for qualified, binding := range tables {
		if qualified == "" || binding.TableID == "" || binding.SchemaVersion == "" {
			return nil, fmt.Errorf("Tabellen-Bindung %q trägt eine leere Kennung: %w", qualified, domainerrors.ErrEmptyIdentifier)
		}
		copiedTables[qualified] = binding
	}
	return &Assembler{source: source, tables: copiedTables, schemaStore: schemaStore}, nil
}

// Consume übersetzt ein Ereignis in höchstens einen Capture-Aufruf:
// BEGIN öffnet, Änderungen hängen an, COMMIT bringt die Transaktion an
// ihre Commit-Position (`SPEC-001`, `cdc.transaction`) samt dem realen
// Quell-Commit-Zeitpunkt (`LH-FA-ADM-004`, als `model.TimePoint` —
// `ADR-0040`) und meldet sie konsumierbar (`LH-FA-CAP-006`). Änderungen
// an nicht aktivierten Tabellen fließen nicht in die Transaktion — CDC
// erfasst nur aktivierte Tabellen (`LH-FA-CFG-001`). Eine
// `*decode.Relation`-Nachricht löst die dynamische Re-Versionierung aus
// (`observeRelation`, `ADR-0015` Folgepflicht) — sie meldet nie ein
// CaptureCommand.
func (a *Assembler) Consume(ctx context.Context, event decode.Event) (*inbound.CaptureCommand, error) {
	switch event := event.(type) {
	case decode.Begin:
		if a.open != nil {
			return nil, ErrBeginWithoutCommit
		}
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
		sourceCommittedAt := model.NewTimePoint(event.CommitTime.UnixNano())
		if err := a.open.tx.Commit(position, sourceCommittedAt); err != nil {
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
	case *decode.Relation:
		return nil, a.observeRelation(ctx, event)
	default:
		// Nachrichtentypen außerhalb des Capture-Pfads (`decode.Decode`
		// liefert für sie bereits kein Ereignis): keine Änderung.
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

// relationComparison trägt das Ergebnis von `classifyRelationColumns`: die
// drei Fälle, die `observeRelation` unterscheidet.
type relationComparison int

const (
	// relationUnchanged: dieselbe Spaltenmenge (Name und Typ-OID) wie die
	// bekannte TableSchema — kein Store-Schreibzugriff.
	relationUnchanged relationComparison = iota
	// relationCompatibleExtension: jede bekannte Spalte kommt unverändert
	// in der eingehenden Relation vor, mindestens eine neue Spalte ist
	// hinzugekommen (`LH-FA-SCH-005`).
	relationCompatibleExtension
	// relationOther trägt jede nicht sicher als Obermenge erkennbare
	// Änderung (Spalte entfernt, Typ einer bestehenden Spalte geändert,
	// Spalte umbenannt) — ohne Store-Schreibzugriff, aber sichtbar als
	// `ErrIncompatibleSchemaChange` gemeldet (Fehlerklasse `schema`,
	// `LH-FA-SCH-004.a`, `observeRelation`).
	relationOther
)

// classifyRelationColumns vergleicht die bekannte Spaltenform gegen die
// einer eingehenden Relation-Nachricht (`ADR-0015` Folgepflicht): fehlt
// eine bekannte Spalte in der eingehenden Relation oder trägt sie dort
// eine andere Typ-OID, ist das Ergebnis `relationOther` — unabhängig von
// der Spaltenzahl. Andernfalls trägt gleiche Spaltenzahl `relationUnchanged`,
// eine größere Spaltenzahl `relationCompatibleExtension`.
func classifyRelationColumns(known []model.Column, incoming []decode.Column) relationComparison {
	incomingByName := make(map[string]uint32, len(incoming))
	for _, column := range incoming {
		incomingByName[column.Name] = column.TypeOID
	}
	for _, column := range known {
		oid, present := incomingByName[column.Name]
		if !present || oid != uint32(column.OID) {
			return relationOther
		}
	}
	switch {
	case len(incoming) == len(known):
		return relationUnchanged
	case len(incoming) > len(known):
		return relationCompatibleExtension
	default:
		return relationOther
	}
}

// tableSchemaFromRelation trägt die Spaltenform einer Relation-Nachricht
// als `model.TableSchema` unter der übergebenen Versions-Kennung
// (`ADR-0015` Folgepflicht).
func tableSchemaFromRelation(versionID model.SchemaVersionID, columns []decode.Column) (model.TableSchema, error) {
	relationColumns := make([]model.Column, len(columns))
	for i, column := range columns {
		relationColumns[i] = model.Column{Name: column.Name, OID: model.ColumnOID(column.TypeOID)}
	}
	return model.NewTableSchema(versionID, relationColumns)
}

// nextSchemaVersionID trägt das Kennungsformat neu registrierter
// Schema-Versionen der dynamischen Re-Versionierung: Tabellen-Kennung und
// Versionsnummer, getrennt durch `-v` — eindeutig je Tabelle, weil die
// Versionsnummer je Tabelle monoton steigt (`SchemaStorePort.CurrentVersion`).
func nextSchemaVersionID(table model.SourceTableID, version int64) model.SchemaVersionID {
	return model.SchemaVersionID(fmt.Sprintf("%s-v%d", table, version))
}

// observeRelation trägt die dynamische Re-Versionierung im laufenden
// Erfassungspfad (`ADR-0015` Folgepflicht, `LH-FA-SCH-005`): ohne
// verdrahteten SchemaStorePort bleibt eine Relation-Nachricht
// wirkungslos. Eine nicht aktivierte Tabelle trägt keine Bindung und
// bleibt ebenfalls wirkungslos (`LH-FA-CFG-001`). Trägt die aktuelle
// Version noch keine TableSchema (die statische Erstaktivierung
// registriert nur die Versions-Zeile, keine Spaltenform —
// `TableActivationPort.Register`), bekommt sie ihre Spaltenform aus der
// eingehenden Relation nachgetragen (Backfill), ohne die Version zu
// wechseln. Eine kompatible Erweiterung registriert eine neue Version
// und hebt die `TableBinding` auf sie; unverändert bleibt ohne Wirkung.
// Jede andere Änderung (`relationOther`) meldet `ErrIncompatibleSchemaChange`
// (Fehlerklasse `schema`, `LH-FA-SCH-004.a`) — der Erfassungspfad endet
// darüber sichtbar, statt die Änderung still zu verwerfen.
func (a *Assembler) observeRelation(ctx context.Context, relation *decode.Relation) error {
	if a.schemaStore == nil {
		return nil
	}
	binding, activated := a.tables[relation.QualifiedName()]
	if !activated {
		return nil
	}
	current, found, err := a.schemaStore.CurrentVersion(ctx, binding.TableID)
	if err != nil {
		return err
	}
	if !found {
		return nil
	}
	known, err := a.schemaStore.TableSchema(ctx, current.ID)
	if err != nil {
		if errors.Is(err, outbound.ErrSchemaVersionUnknown) {
			schema, buildErr := tableSchemaFromRelation(current.ID, relation.Columns)
			if buildErr != nil {
				return buildErr
			}
			_, err := a.schemaStore.RegisterVersion(ctx, current, schema)
			return err
		}
		return err
	}
	comparison := classifyRelationColumns(known.Columns, relation.Columns)
	if comparison == relationOther {
		return fmt.Errorf("%w: %s", ErrIncompatibleSchemaChange, relation.QualifiedName())
	}
	if comparison != relationCompatibleExtension {
		return nil
	}
	nextID := nextSchemaVersionID(binding.TableID, current.Version+1)
	nextVersion, err := model.NewSchemaVersion(nextID, binding.TableID, current.Version+1)
	if err != nil {
		return err
	}
	schema, err := tableSchemaFromRelation(nextID, relation.Columns)
	if err != nil {
		return err
	}
	if _, err := a.schemaStore.RegisterVersion(ctx, nextVersion, schema); err != nil {
		return err
	}
	a.tables[relation.QualifiedName()] = TableBinding{TableID: binding.TableID, SchemaVersion: nextID}
	return nil
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
