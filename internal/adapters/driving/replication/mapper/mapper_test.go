package mapper_test

import (
	"context"
	stderrors "errors"
	"strconv"
	"testing"
	"time"

	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/decode"
	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/mapper"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"

	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
)

// Die Mapper-Tests tragen die Übersetzung in die Domänen-Transaktion:
// Sequenz-Ordnung, Row Images (`SPEC-002`), Aktivierungs-Filter
// (`LH-FA-CAP-001`…003), Commit-Position (`SPEC-003`), die Fehlergrenzen
// von TRUNCATE und Stream-Vertragsverstoß sowie die dynamische
// Re-Versionierung über eine `*decode.Relation`-Nachricht (`ADR-0015`
// Folgepflicht, `LH-FA-SCH-005`).

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

// newAssembler legt den Übersetzer mit den Test-Aktivierungen an, ohne
// SchemaStorePort (`nil`) — die dynamische Re-Versionierung bleibt für
// diese Tests wirkungslos; sie prüfen die Transaktions-Übersetzung.
func newAssembler(t *testing.T, tables map[string]mapper.TableBinding) *mapper.Assembler {
	t.Helper()
	assembler, err := mapper.NewAssembler("src-1", tables, nil)
	if err != nil {
		t.Fatalf("NewAssembler: %v", err)
	}
	return assembler
}

// fakeSchemaStore trägt einen In-Memory-Stub des `outbound.SchemaStorePort`
// für die Assembler-Tests der dynamischen Re-Versionierung (`ADR-0015`
// Folgepflicht) — der reale Adapter gegen PostgreSQL liegt bei
// `schemastore_test.go` (`make test-store`). `registrations` zählt die
// tatsächlich neu geschriebenen Spaltenformen (`RegisterVersion` liefert
// `false`, wenn eine `SchemaVersionID` bereits eine Spaltenform trägt).
type fakeSchemaStore struct {
	versions      map[model.SourceTableID]model.SchemaVersion
	schemas       map[model.SchemaVersionID]model.TableSchema
	registrations int
}

func newFakeSchemaStore() *fakeSchemaStore {
	return &fakeSchemaStore{
		versions: map[model.SourceTableID]model.SchemaVersion{},
		schemas:  map[model.SchemaVersionID]model.TableSchema{},
	}
}

func (f *fakeSchemaStore) CurrentVersion(_ context.Context, table model.SourceTableID) (model.SchemaVersion, bool, error) {
	version, found := f.versions[table]
	return version, found, nil
}

func (f *fakeSchemaStore) RegisterVersion(_ context.Context, version model.SchemaVersion, schema model.TableSchema) (bool, error) {
	if version.ID != schema.VersionID {
		return false, outbound.ErrSchemaVersionMismatch
	}
	f.versions[version.SourceTableID] = version
	if _, exists := f.schemas[schema.VersionID]; exists {
		return false, nil
	}
	f.schemas[schema.VersionID] = schema
	f.registrations++
	return true, nil
}

func (f *fakeSchemaStore) TableSchema(_ context.Context, versionID model.SchemaVersionID) (model.TableSchema, error) {
	schema, found := f.schemas[versionID]
	if !found {
		return model.TableSchema{}, outbound.ErrSchemaVersionUnknown
	}
	return schema, nil
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
	ctx := context.Background()
	assembler := newAssembler(t, testTables())
	relationEvent := &decode.Relation{
		Schema:  "public",
		Name:    "feed",
		Columns: []decode.Column{{Name: "id", Key: true}, {Name: "name"}},
	}

	if command, err := assembler.Consume(ctx, decode.Begin{XID: 42}); err != nil || command != nil {
		t.Fatalf("Begin: command=%v err=%v", command, err)
	}
	if _, err := assembler.Consume(ctx, decode.Change{
		Relation:  relationEvent,
		Operation: decode.OpInsert,
		New:       []*string{pointer("1"), pointer("Wert")},
	}); err != nil {
		t.Fatalf("Insert: %v", err)
	}
	if _, err := assembler.Consume(ctx, decode.Change{
		Relation:  &decode.Relation{Schema: "public", Name: "other", Columns: []decode.Column{{Name: "id", Key: true}}},
		Operation: decode.OpInsert,
		New:       []*string{pointer("9")},
	}); err != nil {
		t.Fatalf("Insert an nicht aktivierter Tabelle: %v", err)
	}
	if _, err := assembler.Consume(ctx, decode.Change{
		Relation:  relationEvent,
		Operation: decode.OpUpdate,
		Old:       []*string{pointer("1"), nil},
		New:       []*string{pointer("1"), nil},
	}); err != nil {
		t.Fatalf("Update: %v", err)
	}
	if _, err := assembler.Consume(ctx, decode.Change{
		Relation:  relationEvent,
		Operation: decode.OpDelete,
		Old:       []*string{pointer("2"), nil},
	}); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	commitTime := time.Date(2026, 2, 4, 12, 30, 0, 0, time.UTC)
	command, err := assembler.Consume(ctx, decode.Commit{CommitLSN: 0xABCDEF01, EndLSN: 0xABCDFFFF, CommitTime: commitTime})
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
	sourceCommittedAt, committedAt := transaction.SourceCommittedAt()
	if !committedAt || sourceCommittedAt != model.NewTimePoint(commitTime.UnixNano()) {
		t.Fatalf("SourceCommittedAt = %+v committed=%v, wollen %+v (LH-FA-ADM-004)", sourceCommittedAt, committedAt, model.NewTimePoint(commitTime.UnixNano()))
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
	ctx := context.Background()
	assembler := newAssembler(t, testTables())
	if _, err := assembler.Consume(ctx, decode.Begin{XID: 7}); err != nil {
		t.Fatalf("Begin: %v", err)
	}
	command, err := assembler.Consume(ctx, decode.Commit{CommitLSN: 100})
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
	ctx := context.Background()
	assembler := newAssembler(t, testTables())
	relationEvent := &decode.Relation{Schema: "public", Name: "feed", Columns: []decode.Column{{Name: "id", Key: true}}}

	if _, err := assembler.Consume(ctx, decode.Change{Relation: relationEvent, Operation: decode.OpInsert}); !stderrors.Is(err, mapper.ErrChangeWithoutBegin) {
		t.Fatalf("Änderung ohne BEGIN: %v", err)
	}
	if _, err := assembler.Consume(ctx, decode.Commit{CommitLSN: 100}); !stderrors.Is(err, mapper.ErrCommitWithoutBegin) {
		t.Fatalf("Commit ohne BEGIN: %v", err)
	}
}

// TestConsumeTruncate trägt TRUNCATE als sichtbaren Fehler der
// Nicht-Unterstützung (`LH-FA-CFG-001.a`).
func TestConsumeTruncate(t *testing.T) {
	ctx := context.Background()
	assembler := newAssembler(t, testTables())
	if _, err := assembler.Consume(ctx, decode.Begin{XID: 42}); err != nil {
		t.Fatalf("Begin: %v", err)
	}
	_, err := assembler.Consume(ctx, decode.Truncate{
		Relations: []*decode.Relation{{Schema: "public", Name: "feed"}},
	})
	if !stderrors.Is(err, mapper.ErrTruncateUnsupported) {
		t.Fatalf("TRUNCATE: %v", err)
	}
}

// TestConsumeBeginWithoutCommit trägt den Doppel-BEGIN als sichtbaren
// Fehler: ein zweites BEGIN überschreibt die offene Transaktion nicht
// still (F-8-Klasse, `LH-QA-REL-001.a` Fehlermodi).
func TestConsumeBeginWithoutCommit(t *testing.T) {
	ctx := context.Background()
	assembler := newAssembler(t, testTables())
	if _, err := assembler.Consume(ctx, decode.Begin{XID: 42}); err != nil {
		t.Fatalf("Begin: %v", err)
	}
	if _, err := assembler.Consume(ctx, decode.Begin{XID: 43}); !stderrors.Is(err, mapper.ErrBeginWithoutCommit) {
		t.Fatalf("Zweites BEGIN: %v", err)
	}
}

// TestNewAssemblerLimits trägt die Konfigurationsgrenzen: leere Quelle
// und Bindungen ohne Kennungen enden über den Domänen-Fehler
// (`ADR-0029`).
func TestNewAssemblerLimits(t *testing.T) {
	if _, err := mapper.NewAssembler("", nil, nil); !stderrors.Is(err, domainerrors.ErrEmptyIdentifier) {
		t.Fatalf("Leere Quelle: %v", err)
	}
	if _, err := mapper.NewAssembler("src-1", map[string]mapper.TableBinding{
		"public.feed": {TableID: "tbl-1"},
	}, nil); !stderrors.Is(err, domainerrors.ErrEmptyIdentifier) {
		t.Fatalf("Bindung ohne Schema-Version: %v", err)
	}
}

// TestConsumeSecondTransaction trägt zwei Transaktionen nacheinander:
// der Sequenz-Zähler startet je Transaktion bei 1 (`SPEC-002`).
func TestConsumeTwoTransactions(t *testing.T) {
	ctx := context.Background()
	assembler := newAssembler(t, testTables())
	relationEvent := &decode.Relation{Schema: "public", Name: "feed", Columns: []decode.Column{{Name: "id", Key: true}}}

	for _, xid := range []uint32{1, 2} {
		if _, err := assembler.Consume(ctx, decode.Begin{XID: xid}); err != nil {
			t.Fatalf("Begin %d: %v", xid, err)
		}
		if _, err := assembler.Consume(ctx, decode.Change{Relation: relationEvent, Operation: decode.OpInsert, New: []*string{pointer("1")}}); err != nil {
			t.Fatalf("Insert %d: %v", xid, err)
		}
		command, err := assembler.Consume(ctx, decode.Commit{CommitLSN: uint64(xid) * 100})
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

// consumedChangeSchemaVersion trägt ein Begin/Change/Commit über den
// Assembler und liest die Schema-Version der resultierenden Change — die
// gemeinsame Prüf-Form der Relation-Klassifikationstests unten.
func consumedChangeSchemaVersion(t *testing.T, ctx context.Context, assembler *mapper.Assembler, xid uint32, relationEvent *decode.Relation) model.SchemaVersionID {
	t.Helper()
	if _, err := assembler.Consume(ctx, decode.Begin{XID: xid}); err != nil {
		t.Fatalf("Begin: %v", err)
	}
	if _, err := assembler.Consume(ctx, decode.Change{Relation: relationEvent, Operation: decode.OpInsert, New: []*string{pointer("1")}}); err != nil {
		t.Fatalf("Insert: %v", err)
	}
	command, err := assembler.Consume(ctx, decode.Commit{CommitLSN: uint64(xid) * 100})
	if err != nil {
		t.Fatalf("Commit: %v", err)
	}
	changes, err := command.Transaction.Changes()
	if err != nil {
		t.Fatalf("Changes: %v", err)
	}
	if len(changes) != 1 {
		t.Fatalf("Change-Anzahl: %d", len(changes))
	}
	return changes[0].SchemaVersion
}

// TestConsumeRelationBackfillsMissingTableSchema trägt den Fall der
// statischen Erstaktivierung (`slice-031`s `EnableTableService` registriert
// nur die Versions-Zeile, keine Spaltenform): die erste real eintreffende
// Relation-Nachricht trägt die Spaltenform zur bestehenden Version nach,
// ohne eine neue Version zu registrieren (`ADR-0015` Folgepflicht).
func TestConsumeRelationBackfillsMissingTableSchema(t *testing.T) {
	ctx := context.Background()
	store := newFakeSchemaStore()
	store.versions["tbl-1"] = model.SchemaVersion{ID: "sv-1", SourceTableID: "tbl-1", Version: 1}
	assembler, err := mapper.NewAssembler("src-1", testTables(), store)
	if err != nil {
		t.Fatalf("NewAssembler: %v", err)
	}
	relationEvent := relation("public", "feed",
		decode.Column{Name: "id", Key: true, TypeOID: 23},
		decode.Column{Name: "name", TypeOID: 25})

	if _, err := assembler.Consume(ctx, relationEvent); err != nil {
		t.Fatalf("Relation: %v", err)
	}
	schema, err := store.TableSchema(ctx, "sv-1")
	if err != nil {
		t.Fatalf("TableSchema nach Backfill: %v", err)
	}
	if len(schema.Columns) != 2 || schema.Columns[0].OID != 23 || schema.Columns[1].OID != 25 {
		t.Fatalf("Backfill-Spaltenform: %+v", schema.Columns)
	}
	if store.versions["tbl-1"].Version != 1 {
		t.Fatalf("Backfill ändert die Versionsnummer: %+v", store.versions["tbl-1"])
	}
	if store.registrations != 1 {
		t.Fatalf("registrations = %d, wollen 1 (genau ein Backfill)", store.registrations)
	}

	version := consumedChangeSchemaVersion(t, ctx, assembler, 1, relationEvent)
	if version != "sv-1" {
		t.Fatalf("Schema-Version nach Backfill: %s, wollen sv-1 (Bindung unverändert)", version)
	}
}

// TestConsumeRelationUnchanged trägt die unveränderte Spaltenform: keine
// neue Version, die Bindung bleibt auf der bekannten Version
// (`LH-FA-SCH-005`).
func TestConsumeRelationUnchanged(t *testing.T) {
	ctx := context.Background()
	store := newFakeSchemaStore()
	store.versions["tbl-1"] = model.SchemaVersion{ID: "sv-1", SourceTableID: "tbl-1", Version: 1}
	store.schemas["sv-1"] = model.TableSchema{VersionID: "sv-1", Columns: []model.Column{
		{Name: "id", OID: 23}, {Name: "name", OID: 25},
	}}
	assembler, err := mapper.NewAssembler("src-1", testTables(), store)
	if err != nil {
		t.Fatalf("NewAssembler: %v", err)
	}
	relationEvent := relation("public", "feed",
		decode.Column{Name: "id", Key: true, TypeOID: 23},
		decode.Column{Name: "name", TypeOID: 25})

	if _, err := assembler.Consume(ctx, relationEvent); err != nil {
		t.Fatalf("Relation: %v", err)
	}
	if store.registrations != 0 {
		t.Fatalf("registrations = %d, wollen 0 (unverändert schreibt nicht)", store.registrations)
	}
	version := consumedChangeSchemaVersion(t, ctx, assembler, 1, relationEvent)
	if version != "sv-1" {
		t.Fatalf("Schema-Version nach unveränderter Relation: %s, wollen sv-1", version)
	}
}

// TestConsumeRelationCompatibleExtension trägt die kompatible Erweiterung:
// alle bekannten Spalten bleiben unverändert, eine neue Spalte kommt
// hinzu — eine neue Version wird registriert und die Bindung darauf
// gehoben (`LH-FA-SCH-005` Boundary).
func TestConsumeRelationCompatibleExtension(t *testing.T) {
	ctx := context.Background()
	store := newFakeSchemaStore()
	store.versions["tbl-1"] = model.SchemaVersion{ID: "sv-1", SourceTableID: "tbl-1", Version: 1}
	store.schemas["sv-1"] = model.TableSchema{VersionID: "sv-1", Columns: []model.Column{
		{Name: "id", OID: 23}, {Name: "name", OID: 25},
	}}
	assembler, err := mapper.NewAssembler("src-1", testTables(), store)
	if err != nil {
		t.Fatalf("NewAssembler: %v", err)
	}
	extendedRelation := relation("public", "feed",
		decode.Column{Name: "id", Key: true, TypeOID: 23},
		decode.Column{Name: "name", TypeOID: 25},
		decode.Column{Name: "extra", TypeOID: 25})

	if _, err := assembler.Consume(ctx, extendedRelation); err != nil {
		t.Fatalf("Relation: %v", err)
	}
	if store.registrations != 1 {
		t.Fatalf("registrations = %d, wollen 1 (eine neue Version)", store.registrations)
	}
	nextVersion := store.versions["tbl-1"]
	if nextVersion.Version != 2 || nextVersion.ID == "sv-1" {
		t.Fatalf("neue Version: %+v", nextVersion)
	}
	version := consumedChangeSchemaVersion(t, ctx, assembler, 1, extendedRelation)
	if version != nextVersion.ID {
		t.Fatalf("Schema-Version nach Erweiterung: %s, wollen %s (Bindung angehoben)", version, nextVersion.ID)
	}
	schema, err := store.TableSchema(ctx, nextVersion.ID)
	if err != nil {
		t.Fatalf("TableSchema der neuen Version: %v", err)
	}
	if len(schema.Columns) != 3 {
		t.Fatalf("Spaltenform der neuen Version: %+v", schema.Columns)
	}
}

// TestConsumeRelationOtherChangeReportsSchemaError trägt jede nicht sicher
// als Obermenge erkennbare Änderung (hier: eine bekannte Spalte trägt eine
// andere Typ-OID) als sichtbaren Fehler der Fehlerklasse `schema`
// (`ErrIncompatibleSchemaChange`, `LH-FA-SCH-004.a`): kein
// Store-Schreibzugriff, und die Bindung bleibt auf der bekannten Version
// stehen.
func TestConsumeRelationOtherChangeReportsSchemaError(t *testing.T) {
	ctx := context.Background()
	store := newFakeSchemaStore()
	store.versions["tbl-1"] = model.SchemaVersion{ID: "sv-1", SourceTableID: "tbl-1", Version: 1}
	store.schemas["sv-1"] = model.TableSchema{VersionID: "sv-1", Columns: []model.Column{
		{Name: "id", OID: 23}, {Name: "amount", OID: 25},
	}}
	assembler, err := mapper.NewAssembler("src-1", testTables(), store)
	if err != nil {
		t.Fatalf("NewAssembler: %v", err)
	}
	typeChangedRelation := relation("public", "feed",
		decode.Column{Name: "id", Key: true, TypeOID: 23},
		decode.Column{Name: "amount", TypeOID: 23}) // Typ geändert: 25 -> 23

	if _, err := assembler.Consume(ctx, typeChangedRelation); !stderrors.Is(err, mapper.ErrIncompatibleSchemaChange) {
		t.Fatalf("Relation mit geändertem Spaltentyp: %v, wollen ErrIncompatibleSchemaChange", err)
	}
	if store.registrations != 0 {
		t.Fatalf("registrations = %d, wollen 0 (kein Store-Schreibzugriff bei relationOther)", store.registrations)
	}
	version := consumedChangeSchemaVersion(t, ctx, assembler, 1, typeChangedRelation)
	if version != "sv-1" {
		t.Fatalf("Schema-Version nach dem gemeldeten Fehler: %s, wollen sv-1 (Bindung unverändert)", version)
	}
}

// TestConsumeRelationNilSchemaStoreStaysNoop trägt den Default ohne
// verdrahteten SchemaStorePort: die Relation-Behandlung bleibt
// wirkungslos, kein Fehler.
func TestConsumeRelationNilSchemaStoreStaysNoop(t *testing.T) {
	ctx := context.Background()
	assembler := newAssembler(t, testTables())
	relationEvent := relation("public", "feed", decode.Column{Name: "id", Key: true, TypeOID: 23})
	if _, err := assembler.Consume(ctx, relationEvent); err != nil {
		t.Fatalf("Relation ohne SchemaStorePort: %v", err)
	}
}
