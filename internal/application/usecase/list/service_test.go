package list_test

import (
	"context"
	stderrors "errors"
	"testing"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/list"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// fakeActivation trägt den Aktivierungs-Port als Fake (`ADR-0030`); die
// Liste liest den Bindungs-Zeilen-Bestand, die übrigen Operationen tragen
// die Schnittstelle.
type fakeActivation struct {
	tables  []model.SourceTable
	listErr error
}

func (f *fakeActivation) TableExists(ctx context.Context, schema, table string) (bool, error) {
	return true, nil
}

func (f *fakeActivation) Registered(ctx context.Context, source model.SourceID, schema, table string) (model.SourceTable, bool, error) {
	return model.SourceTable{}, false, nil
}

func (f *fakeActivation) Register(ctx context.Context, table model.SourceTable, version model.SchemaVersion) (bool, error) {
	return true, nil
}

func (f *fakeActivation) Unregister(ctx context.Context, table model.SourceTable) (outbound.ActivationRemoval, error) {
	return outbound.ActivationRemoved, nil
}

func (f *fakeActivation) List(ctx context.Context, source model.SourceID) ([]model.SourceTable, error) {
	return f.tables, f.listErr
}

func (f *fakeActivation) Publish(ctx context.Context, publication, schema, table string) error {
	return nil
}

func (f *fakeActivation) Unpublish(ctx context.Context, publication, schema, table string) error {
	return nil
}

// TestListTables trägt den Happy Path (`LH-FA-CFG-004`): die Liste trägt
// die aktivierten Tabellen der Quelle.
func TestListTables(t *testing.T) {
	first, err := model.NewSourceTable("tbl-1", "src-1", "public", "t1")
	if err != nil {
		t.Fatalf("Tabelle: %v", err)
	}
	second, err := model.NewSourceTable("tbl-2", "src-1", "public", "t2")
	if err != nil {
		t.Fatalf("Tabelle: %v", err)
	}
	service := list.NewListTablesService(&fakeActivation{tables: []model.SourceTable{first, second}})

	result, err := service.ListTables(context.Background(), inbound.ListTablesQuery{Source: "src-1"})
	if err != nil {
		t.Fatalf("ListTables: %v", err)
	}
	if len(result.Tables) != 2 {
		t.Fatalf("Tabellen-Liste: %d Tabellen (Erwartung: 2)", len(result.Tables))
	}
	if result.Tables[0].QualifiedName() != "public.t1" || result.Tables[1].QualifiedName() != "public.t2" {
		t.Fatalf("Tabellen-Liste: %v", result.Tables)
	}
}

// TestListTablesEmpty trägt den Boundary-Pfad (`LH-FA-CFG-004`): ohne
// Aktivierung trägt die Rückkehr eine leere Liste — auch dann, wenn der
// Port keinen Bestand trägt.
func TestListTablesEmpty(t *testing.T) {
	service := list.NewListTablesService(&fakeActivation{})

	result, err := service.ListTables(context.Background(), inbound.ListTablesQuery{Source: "src-1"})
	if err != nil {
		t.Fatalf("ListTables: %v", err)
	}
	if result.Tables == nil || len(result.Tables) != 0 {
		t.Fatalf("Tabellen-Liste ohne Bestand: %v (Erwartung: leere Liste)", result.Tables)
	}
}

// TestListTablesEmptySource trägt die Kennungs-Grenze der Liste.
func TestListTablesEmptySource(t *testing.T) {
	service := list.NewListTablesService(&fakeActivation{})

	_, err := service.ListTables(context.Background(), inbound.ListTablesQuery{})
	if !stderrors.Is(err, domainerrors.ErrEmptyIdentifier) {
		t.Fatalf("leere Eingabe: %v (Erwartung: ErrEmptyIdentifier)", err)
	}
}
