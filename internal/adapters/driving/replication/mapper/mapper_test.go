package mapper_test

import (
	stderrors "errors"
	"strconv"
	"testing"

	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/decode"
	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/mapper"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Die Mapper-Tests tragen die Übersetzung in die Domänen-Transaktion:
// Sequenz-Ordnung, Row Images (`SPEC-002`), Aktivierungs-Filter
// (`LH-FA-CAP-001`…003), Commit-Position (`SPEC-003`) und die
// Fehlergrenzen von TRUNCATE und Stream-Vertragsverstoß.

// relation trägt eine Dekodier-Relation für die Tests.
func relation(schema, name string, columns ...decode.Column) *decode.Relation {
	return &decode.Relation{Schema: schema, Name: name, Columns: columns}
}

// pointer trägt einen Text-Wert-Zeiger.
func pointer(value string) *string { return &value }

// transactionID trägt die Domänen-Transaktions-Kennung einer XID.
func transactionID(xid uint32) model.TransactionID {
	return model.TransactionID(strconv.FormatUint(uint64(xid), 10))
}

// newAssembler legt den Übersetzer mit den Test-Aktivierungen an.
func newAssembler(t *testing.T, tables map[string]mapper.TableBinding) *mapper.Assembler {
	t.Helper()
	assembler, err := mapper.NewAssembler("src-1", tables)
	if err != nil {
		t.Fatalf("NewAssembler: %v", err)
	}
	return assembler
}

// testTables trägt die Aktivierungen der Tests: `public.feed` ist
// aktiviert, die übrigen Tabellen nicht.
func testTables() map[string]mapper.TableBinding {
	return map[string]mapper.TableBinding{
		"public.feed": {
			TableID:       "tbl-1",
			SchemaVersion: "sv-1",
		},
	}
}

// TestConsumeFullTransaction trägt einen vollständigen Durchlauf: BEGIN,
// Relation, Insert, Update, Delete, COMMIT — der Commit meldet die
// Transaktion als CaptureCommand mit Sequenz-Reihenfolge und
// Commit-Position (`LH-FA-CAP-001`…003, `LH-FA-CAP-006.a`).
func TestConsumeFullTransaction(t *testing.T) {
	assembler := newAssembler(t, testTables())
	relationEvent := &decode.Relation{
		Schema:  "public",
		Name:    "feed",
		Columns: []decode.Column{{Name: "id", Key: true}, {Name: "name"}},
	}

	if command, err := assembler.Consume(decode.Begin{XID: 42}); err != nil || command != nil {
		t.Fatalf("Begin: command=%v err=%v", command, err)
	}
	if _, err := assembler.Consume(decode.Change{
		XID:       42,
		Relation:  relationEvent,
		Operation: decode.OpInsert,
		New:       []*string{pointer("1"), pointer("Wert")},
	}); err != nil {
		t.Fatalf("Insert: %v", err)
	}
	if _, err := assembler.Consume(decode.Change{
		XID:       42,
		Relation:  &decode.Relation{Schema: "public", Name: "other", Columns: []decode.Column{{Name: "id", Key: true}}},
		Operation: decode.OpInsert,
		New:       []*string{pointer("9")},
	}); err != nil {
		t.Fatalf("Insert an nicht aktivierter Tabelle: %v", err)
	}
	if _, err := assembler.Consume(decode.Change{
		XID:       42,
		Relation:  relationEvent,
		Operation: decode.OpUpdate,
		Old:       []*string{pointer("1"), nil},
		New:       []*string{pointer("1"), nil},
	}); err != nil {
		t.Fatalf("Update: %v", err)
	}
	if _, err := assembler.Consume(decode.Change{
		XID:       42,
		Relation:  relationEvent,
		Operation: decode.OpDelete,
		Old:       []*string{pointer("2"), nil},
	}); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	command, err := assembler.Consume(decode.Commit{CommitLSN: 0xABCDEF01, EndLSN: 0xABCDFFFF})
	if err != nil {
		t.Fatalf("Commit: %v", err)
	}
	if command == nil {
		t.Fatalf("Commit meldet kein CaptureCommand")
	}
	transaction := command.Transaction
	if transaction.ID != "42" || transaction.SourceID != "src-1" {
		t.Fatalf("Transaktion: %+v", transaction)
	}
	position, committed := transaction.CommitPosition()
	if !committed || position.Offset != 0xABCDEF01 || position.SourceID != "src-1" {
		t.Fatalf("Commit-Position: %+v committed=%v", position, committed)
	}
	changes, err := transaction.Changes()
	if err != nil {
		t.Fatalf("Changes: %v", err)
	}
	if len(changes) != 3 {
		t.Fatalf("Change-Anzahl: %d (nicht aktivierte Tabelle fließt nicht ein)", len(changes))
	}
	if changes[0].ID != "42-1" || changes[0].Sequence != 1 || changes[0].Operation != model.OperationInsert {
		t.Fatalf("Change 1: %+v", changes[0])
	}
	if string(changes[0].NewImage) != `{"id":"1","name":"Wert"}` {
		t.Fatalf("Insert-Neu-Image: %s", changes[0].NewImage)
	}
	if changes[1].ID != "42-2" || changes[1].Sequence != 2 || changes[1].Operation != model.OperationUpdate {
		t.Fatalf("Change 2: %+v", changes[1])
	}
	if string(changes[1].NewImage) != `{"id":"1"}` {
		t.Fatalf("Update-Neu-Image (NULL-Spalte als Abwesenheit): %s", changes[1].NewImage)
	}
	if string(changes[1].OldImage) != `{"id":"1"}` {
		t.Fatalf("Update-Alt-Image (Schlüssel-Spalte): %s", changes[1].OldImage)
	}
	if changes[2].ID != "42-3" || changes[2].Operation != model.OperationDelete {
		t.Fatalf("Change 3: %+v", changes[2])
	}
	if string(changes[2].OldImage) != `{"id":"2"}` {
		t.Fatalf("Delete-Alt-Image: %s", changes[2].OldImage)
	}
}

// TestConsumeEmptyTransaction trägt ein BEGIN/COMMIT-Paar ohne
// relevante Relation: die leere committed Transaktion wird als Command
// gemeldet (der Port-Kontrakt persistiert sie mit ihrer Commit-Position).
func TestConsumeEmptyTransaction(t *testing.T) {
	assembler := newAssembler(t, testTables())
	if _, err := assembler.Consume(decode.Begin{XID: 7}); err != nil {
		t.Fatalf("Begin: %v", err)
	}
	command, err := assembler.Consume(decode.Commit{CommitLSN: 100})
	if err != nil {
		t.Fatalf("Commit: %v", err)
	}
	changes, err := command.Transaction.Changes()
	if err != nil {
		t.Fatalf("Changes: %v", err)
	}
	if len(changes) != 0 {
		t.Fatalf("Leere Transaktion trägt %d Changes", len(changes))
	}
}

// TestConsumeWithoutBegin trägt die Stream-Vertragsverstöße als
// sichtbare Fehler.
func TestConsumeWithoutBegin(t *testing.T) {
	assembler := newAssembler(t, testTables())
	relationEvent := &decode.Relation{Schema: "public", Name: "feed", Columns: []decode.Column{{Name: "id", Key: true}}}

	if _, err := assembler.Consume(decode.Change{Relation: relationEvent, Operation: decode.OpInsert}); !stderrors.Is(err, mapper.ErrChangeWithoutBegin) {
		t.Fatalf("Änderung ohne BEGIN: %v", err)
	}
	if _, err := assembler.Consume(decode.Commit{CommitLSN: 100}); !stderrors.Is(err, mapper.ErrCommitWithoutBegin) {
		t.Fatalf("Commit ohne BEGIN: %v", err)
	}
}

// TestConsumeTruncate trägt TRUNCATE als sichtbaren Fehler der
// Nicht-Unterstützung (`LH-FA-CFG-001.a`).
func TestConsumeTruncate(t *testing.T) {
	assembler := newAssembler(t, testTables())
	if _, err := assembler.Consume(decode.Begin{XID: 42}); err != nil {
		t.Fatalf("Begin: %v", err)
	}
	_, err := assembler.Consume(decode.Truncate{
		Relations: []*decode.Relation{{Schema: "public", Name: "feed"}},
	})
	if !stderrors.Is(err, mapper.ErrTruncateUnsupported) {
		t.Fatalf("TRUNCATE: %v", err)
	}
}

// TestNewAssemblerLimits trägt die Konfigurationsgrenzen: leere Quelle
// und Bindungen ohne Kennungen enden über den Domänen-Fehler
// (`ADR-0029`).
func TestNewAssemblerLimits(t *testing.T) {
	if _, err := mapper.NewAssembler("", nil); !stderrors.Is(err, domainerrors.ErrEmptyIdentifier) {
		t.Fatalf("Leere Quelle: %v", err)
	}
	if _, err := mapper.NewAssembler("src-1", map[string]mapper.TableBinding{
		"public.feed": {TableID: "tbl-1"},
	}); !stderrors.Is(err, domainerrors.ErrEmptyIdentifier) {
		t.Fatalf("Bindung ohne Schema-Version: %v", err)
	}
}

// TestConsumeSecondTransaction trägt zwei Transaktionen nacheinander:
// der Sequenz-Zähler startet je Transaktion bei 1 (`SPEC-002`).
func TestConsumeTwoTransactions(t *testing.T) {
	assembler := newAssembler(t, testTables())
	relationEvent := &decode.Relation{Schema: "public", Name: "feed", Columns: []decode.Column{{Name: "id", Key: true}}}

	for _, xid := range []uint32{1, 2} {
		if _, err := assembler.Consume(decode.Begin{XID: xid}); err != nil {
			t.Fatalf("Begin %d: %v", xid, err)
		}
		if _, err := assembler.Consume(decode.Change{XID: xid, Relation: relationEvent, Operation: decode.OpInsert, New: []*string{pointer("1")}}); err != nil {
			t.Fatalf("Insert %d: %v", xid, err)
		}
		command, err := assembler.Consume(decode.Commit{CommitLSN: uint64(xid) * 100})
		if err != nil {
			t.Fatalf("Commit %d: %v", xid, err)
		}
		changes, err := command.Transaction.Changes()
		if err != nil {
			t.Fatalf("Changes %d: %v", xid, err)
		}
		if changes[0].ID != model.ChangeID(transactionID(xid)+"-1") {
			t.Fatalf("Change-Kennung %d: %s", xid, changes[0].ID)
		}
		if changes[0].TransactionID != transactionID(xid) {
			t.Fatalf("Transaktions-Kennung %d: %s", xid, changes[0].TransactionID)
		}
	}
}