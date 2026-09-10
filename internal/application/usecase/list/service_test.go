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
// Liste liest den Bindungs-Zeilen-Bestand und die
// Publication-Mitgliedschaft je Tabelle, die übrigen Operationen tragen
// die Schnittstelle.
type fakeActivation struct {
	tables    []model.SourceTable
	published map[string]bool
	listErr   error
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

func (f *fakeActivation) Published(ctx context.Context, publication, schema, table string) (bool, error) {
	return f.published[schema+"."+table], nil
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
	service := list.NewListTablesService(&fakeActivation{
		tables:    []model.SourceTable{first, second},
		published: map[string]bool{"public.t1": true, "public.t2": true},
	})

	result, err := service.ListTables(context.Background(), listQuery())
	if err != nil {
		t.Fatalf("ListTables: %v", err)
	}
	if len(result.Tables) != 2 || len(result.Retained) != 0 {
		t.Fatalf("Tabellen-Liste: %d aktiviert, %d Herkunft (Erwartung: 2 und 0)", len(result.Tables), len(result.Retained))
	}
	if result.Tables[0].QualifiedName() != "public.t1" || result.Tables[1].QualifiedName() != "public.t2" {
		t.Fatalf("Tabellen-Liste: %v", result.Tables)
	}
}

// TestListTablesRetainedPartition trägt die Trennung der Herkunfts-Zeilen:
// eine Bindungs-Zeile ohne Publication-Mitgliedschaft liest sich in der
// Retained-Liste, nicht in der aktivierten (`LH-FA-CFG-002` Out-of-Scope).
func TestListTablesRetainedPartition(t *testing.T) {
	first, err := model.NewSourceTable("tbl-1", "src-1", "public", "t1")
	if err != nil {
		t.Fatalf("Tabelle: %v", err)
	}
	second, err := model.NewSourceTable("tbl-2", "src-1", "public", "t2")
	if err != nil {
		t.Fatalf("Tabelle: %v", err)
	}
	service := list.NewListTablesService(&fakeActivation{
		tables:    []model.SourceTable{first, second},
		published: map[string]bool{"public.t1": true},
	})

	result, err := service.ListTables(context.Background(), listQuery())
	if err != nil {
		t.Fatalf("ListTables: %v", err)
	}
	if len(result.Tables) != 1 || result.Tables[0].QualifiedName() != "public.t1" {
		t.Fatalf("aktivierte Tabellen: %v (Erwartung: public.t1)", result.Tables)
	}
	if len(result.Retained) != 1 || result.Retained[0].QualifiedName() != "public.t2" {
		t.Fatalf("Herkunfts-Zeilen: %v (Erwartung: public.t2)", result.Retained)
	}
}

// TestListTablesEmpty trägt den Boundary-Pfad (`LH-FA-CFG-004`): ohne
// Aktivierung tragen beide Rückgaben leere Listen — auch dann, wenn der
// Port keinen Bestand trägt.
func TestListTablesEmpty(t *testing.T) {
	service := list.NewListTablesService(&fakeActivation{})

	result, err := service.ListTables(context.Background(), listQuery())
	if err != nil {
		t.Fatalf("ListTables: %v", err)
	}
	if result.Tables == nil || len(result.Tables) != 0 ||
		result.Retained == nil || len(result.Retained) != 0 {
		t.Fatalf("Tabellen-Liste ohne Bestand: %+v (Erwartung: leere Listen)", result)
	}
}

// TestListTablesEmptyQuery trägt die Kennungs-Grenze der Liste.
func TestListTablesEmptyQuery(t *testing.T) {
	service := list.NewListTablesService(&fakeActivation{})

	_, err := service.ListTables(context.Background(), inbound.ListTablesQuery{})
	if !stderrors.Is(err, domainerrors.ErrEmptyIdentifier) {
		t.Fatalf("leere Eingabe: %v (Erwartung: ErrEmptyIdentifier)", err)
	}
}

// listQuery trägt die Abfrage der Quelle samt Publication.
func listQuery() inbound.ListTablesQuery {
	return inbound.ListTablesQuery{Source: "src-1", Publication: "pub-1"}
}
