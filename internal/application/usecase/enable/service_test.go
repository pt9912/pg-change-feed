package enable_test

import (
	"context"
	stderrors "errors"
	"testing"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/enable"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// fakeActivation trägt den Aktivierungs-Port als Fake (`ADR-0030`); die
// Aktivierungs- und Publication-Züge melden ihre Aufrufe, die Stubs der
// übrigen Operationen tragen die Schnittstelle.
type fakeActivation struct {
	exists        bool
	existsErr     error
	registerErr   error
	alreadyActive bool

	registerCalls int
	publishCalls  int
	lastTable     model.SourceTable
	lastVersion   model.SchemaVersion
}

func (f *fakeActivation) TableExists(ctx context.Context, schema, table string) (bool, error) {
	return f.exists, f.existsErr
}

func (f *fakeActivation) Registered(ctx context.Context, source model.SourceID, schema, table string) (model.SourceTable, bool, error) {
	return model.SourceTable{}, f.alreadyActive, nil
}

func (f *fakeActivation) Register(ctx context.Context, table model.SourceTable, version model.SchemaVersion) (bool, error) {
	f.registerCalls++
	f.lastTable = table
	f.lastVersion = version
	if f.registerErr != nil {
		return false, f.registerErr
	}
	return !f.alreadyActive, nil
}

func (f *fakeActivation) Unregister(ctx context.Context, table model.SourceTable) (outbound.ActivationRemoval, error) {
	return outbound.ActivationRemoved, nil
}

func (f *fakeActivation) List(ctx context.Context, source model.SourceID) ([]model.SourceTable, error) {
	return nil, nil
}

func (f *fakeActivation) Publish(ctx context.Context, publication, schema, table string) error {
	f.publishCalls++
	return nil
}

func (f *fakeActivation) Unpublish(ctx context.Context, publication, schema, table string) error {
	return nil
}

// TestEnableHappyPath trägt die Aktivierung (`LH-FA-CFG-001` Happy Path):
// der Use Case trägt die Bindungs-Zeilen und die Publication — die
// Rückkehr meldet die aktivierte Tabelle.
func TestEnableHappyPath(t *testing.T) {
	fake := &fakeActivation{exists: true}
	service := enable.NewEnableTableService(fake)

	result, err := service.Enable(context.Background(), inbound.EnableTableCommand{
		Source:          "src-1",
		Schema:          "public",
		Table:           "t1",
		TableID:         "tbl-1",
		SchemaVersionID: "sv-1",
		Version:         1,
		Publication:     "pub-1",
	})
	if err != nil {
		t.Fatalf("Enable: %v", err)
	}
	if result.AlreadyEnabled {
		t.Fatalf("neue Aktivierung meldet AlreadyEnabled")
	}
	if result.Table.ID != "tbl-1" || result.Table.QualifiedName() != "public.t1" {
		t.Fatalf("aktivierte Tabelle: %+v", result.Table)
	}
	if fake.registerCalls != 1 || fake.publishCalls != 1 {
		t.Fatalf("Aktivierungs-Züge: Register %d, Publish %d (Erwartung: je 1)", fake.registerCalls, fake.publishCalls)
	}
	if fake.lastTable.ID != "tbl-1" || fake.lastVersion.ID != "sv-1" || fake.lastVersion.Version != 1 {
		t.Fatalf("Bindungs-Zeilen: Tabelle %v, Version %v", fake.lastTable, fake.lastVersion)
	}
}

// TestEnableIdempotent trägt den Boundary-Pfad (`LH-FA-CFG-001`): die
// bereits aktivierte Tabelle meldet das Ergebnis, ohne den Stand zu
// ändern — die Publication läuft idempotent nach.
func TestEnableIdempotent(t *testing.T) {
	fake := &fakeActivation{exists: true, alreadyActive: true}
	service := enable.NewEnableTableService(fake)

	result, err := service.Enable(context.Background(), inbound.EnableTableCommand{
		Source:          "src-1",
		Schema:          "public",
		Table:           "t1",
		TableID:         "tbl-1",
		SchemaVersionID: "sv-1",
		Version:         1,
		Publication:     "pub-1",
	})
	if err != nil {
		t.Fatalf("Enable: %v", err)
	}
	if !result.AlreadyEnabled {
		t.Fatalf("Erneut aktivierung meldet AlreadyEnabled=false")
	}
	if fake.registerCalls != 1 || fake.publishCalls != 1 {
		t.Fatalf("Aktivierungs-Züge: Register %d, Publish %d (Erwartung: je 1)", fake.registerCalls, fake.publishCalls)
	}
}

// TestEnableMissingTable trägt den Negative-Pfad (`LH-FA-CFG-001`):
// eine nicht existierende Tabelle endet sichtbar, ohne Bindungs-Zeilen
// oder Publication zu tragen.
func TestEnableMissingTable(t *testing.T) {
	fake := &fakeActivation{exists: false}
	service := enable.NewEnableTableService(fake)

	_, err := service.Enable(context.Background(), inbound.EnableTableCommand{
		Source:          "src-1",
		Schema:          "public",
		Table:           "t1",
		TableID:         "tbl-1",
		SchemaVersionID: "sv-1",
		Version:         1,
		Publication:     "pub-1",
	})
	if !stderrors.Is(err, inbound.ErrSourceTableMissing) {
		t.Fatalf("fehlende Tabelle: %v (Erwartung: ErrSourceTableMissing)", err)
	}
	if fake.registerCalls != 0 || fake.publishCalls != 0 {
		t.Fatalf("fehlende Tabelle trägt Aktivierungs-Züge: Register %d, Publish %d", fake.registerCalls, fake.publishCalls)
	}
}

// TestEnableInvalidCommand trägt die Kennungs-Grenzen der Aktivierung:
// leere Kennungen und eine Version unter 1 enden über die
// Domänen-Invarianten (`ADR-0029`), ohne die Quelle zu berühren.
func TestEnableInvalidCommand(t *testing.T) {
	fake := &fakeActivation{exists: true}
	service := enable.NewEnableTableService(fake)

	_, err := service.Enable(context.Background(), inbound.EnableTableCommand{
		Source:          "src-1",
		Schema:          "public",
		Table:           "t1",
		SchemaVersionID: "sv-1",
		Version:         1,
		Publication:     "pub-1",
	})
	if !stderrors.Is(err, domainerrors.ErrEmptyIdentifier) {
		t.Fatalf("leere Tabellen-Kennung: %v (Erwartung: ErrEmptyIdentifier)", err)
	}

	_, err = service.Enable(context.Background(), inbound.EnableTableCommand{
		Source:          "src-1",
		Schema:          "public",
		Table:           "t1",
		TableID:         "tbl-1",
		SchemaVersionID: "sv-1",
		Version:         0,
		Publication:     "pub-1",
	})
	if !stderrors.Is(err, domainerrors.ErrNonPositiveVersion) {
		t.Fatalf("Version 0: %v (Erwartung: ErrNonPositiveVersion)", err)
	}
	if fake.registerCalls != 0 || fake.publishCalls != 0 {
		t.Fatalf("ungültige Eingabe trägt Aktivierungs-Züge: Register %d, Publish %d", fake.registerCalls, fake.publishCalls)
	}
}
