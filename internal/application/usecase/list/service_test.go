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
// die Schnittstelle. Die beiden `…ErrFor`-Tabellen binden einen Port-Fehler
// an genau den Eingabewert, auf den er antwortet (LP2): jede andere Quelle
// bzw. Publication-plus-Adresse trägt derselbe Fake.
type fakeActivation struct {
	tables    []model.SourceTable
	published map[string]bool

	listErrFor      map[model.SourceID]error // Quell-Kennung → Fehler
	publishedErrFor map[string]error         // publication + " " + schema.table → Fehler
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
	if err, ok := f.listErrFor[source]; ok {
		return nil, err
	}
	return f.tables, nil
}

func (f *fakeActivation) Publish(ctx context.Context, publication, schema, table string) error {
	return nil
}

func (f *fakeActivation) Unpublish(ctx context.Context, publication, schema, table string) error {
	return nil
}

func (f *fakeActivation) Published(ctx context.Context, publication, schema, table string) (bool, error) {
	if err, ok := f.publishedErrFor[publication+" "+schema+"."+table]; ok {
		return false, err
	}
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

// TestListTablesListErrorFollowsSource trägt den Fehlerpfad des
// Bindungs-Zeilen-Lesens: ein Fehler von `List` wird unverändert
// durchgereicht. Die Ablehnung ist an die Quell-Kennung gebunden — derselbe
// Fake liest eine andere Quelle ohne Fehler.
func TestListTablesListErrorFollowsSource(t *testing.T) {
	wantErr := stderrors.New("Bindungs-Zeilen nicht lesbar")
	fake := &fakeActivation{listErrFor: map[model.SourceID]error{"src-kaputt": wantErr}}
	service := list.NewListTablesService(fake)

	query := listQuery()
	query.Source = "src-kaputt"
	_, err := service.ListTables(context.Background(), query)
	if !stderrors.Is(err, wantErr) {
		t.Fatalf("Port-Fehler: %v (Erwartung: %v)", err, wantErr)
	}

	result, err := service.ListTables(context.Background(), listQuery())
	if err != nil {
		t.Fatalf("Quelle src-1: %v", err)
	}
	if result.Tables == nil || len(result.Tables) != 0 || result.Retained == nil || len(result.Retained) != 0 {
		t.Fatalf("Tabellen-Liste: %+v (Erwartung: leere Listen)", result)
	}
}

// TestListTablesPublishedErrorFollowsPublication trägt den Fehlerpfad der
// Publication-Mitgliedschaft: ein Fehler von `Published` wird unverändert
// durchgereicht. Die Ablehnung ist an die Publication und die Adresse der
// gelesenen Bindungs-Zeile gebunden — derselbe Fake trägt dieselbe
// Bindungs-Zeile in einer anderen Publication durch.
func TestListTablesPublishedErrorFollowsPublication(t *testing.T) {
	first, err := model.NewSourceTable("tbl-1", "src-1", "public", "t1")
	if err != nil {
		t.Fatalf("Tabelle: %v", err)
	}
	second, err := model.NewSourceTable("tbl-2", "src-1", "public", "t2")
	if err != nil {
		t.Fatalf("Tabelle: %v", err)
	}
	wantErr := stderrors.New("Publication nicht lesbar")
	fake := &fakeActivation{
		tables:          []model.SourceTable{first, second},
		published:       map[string]bool{"public.t1": true, "public.t2": true},
		publishedErrFor: map[string]error{"pub-1 public.t2": wantErr},
	}
	service := list.NewListTablesService(fake)

	_, err = service.ListTables(context.Background(), listQuery())
	if !stderrors.Is(err, wantErr) {
		t.Fatalf("Port-Fehler: %v (Erwartung: %v)", err, wantErr)
	}

	other := listQuery()
	other.Publication = "pub-2"
	result, err := service.ListTables(context.Background(), other)
	if err != nil {
		t.Fatalf("Publication pub-2: %v", err)
	}
	if len(result.Tables) != 2 || len(result.Retained) != 0 {
		t.Fatalf("Tabellen-Liste: %d aktiviert, %d Herkunft (Erwartung: 2 und 0)", len(result.Tables), len(result.Retained))
	}
}

// listQuery trägt die Abfrage der Quelle samt Publication.
func listQuery() inbound.ListTablesQuery {
	return inbound.ListTablesQuery{Source: "src-1", Publication: "pub-1"}
}
