package status_test

import (
	"context"
	stderrors "errors"
	"testing"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/status"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// fakeActivation trägt den Aktivierungs-Port als Fake (`ADR-0030`); die
// Status-Abfrage liest Existenz und Bindungs-Zeile, die übrigen
// Operationen tragen die Schnittstelle.
type fakeActivation struct {
	exists      bool
	existsErr   error
	registered  bool
	registerErr error
}

func (f *fakeActivation) TableExists(ctx context.Context, schema, table string) (bool, error) {
	return f.exists, f.existsErr
}

func (f *fakeActivation) Registered(ctx context.Context, source model.SourceID, schema, table string) (model.SourceTable, bool, error) {
	return model.SourceTable{}, f.registered, f.registerErr
}

func (f *fakeActivation) Register(ctx context.Context, table model.SourceTable, version model.SchemaVersion) (bool, error) {
	return true, nil
}

func (f *fakeActivation) Unregister(ctx context.Context, table model.SourceTable) (outbound.ActivationRemoval, error) {
	return outbound.ActivationRemoved, nil
}

func (f *fakeActivation) List(ctx context.Context, source model.SourceID) ([]model.SourceTable, error) {
	return nil, nil
}

func (f *fakeActivation) Publish(ctx context.Context, publication, schema, table string) error {
	return nil
}

func (f *fakeActivation) Unpublish(ctx context.Context, publication, schema, table string) error {
	return nil
}

// TestStatusEnabled trägt den Happy Path (`LH-FA-CFG-003`): die
// aktivierte Tabelle meldet den Zustand „aktiviert".
func TestStatusEnabled(t *testing.T) {
	service := status.NewGetStatusService(&fakeActivation{exists: true, registered: true})

	result, err := service.Status(context.Background(), inbound.GetStatusQuery{
		Source: "src-1", Schema: "public", Table: "t1",
	})
	if err != nil {
		t.Fatalf("Status: %v", err)
	}
	if !result.Enabled {
		t.Fatalf("aktivierte Tabelle meldet „nicht aktiviert“")
	}
}

// TestStatusNotActivated trägt den Boundary-Pfad (`LH-FA-CFG-003`): eine
// nie aktivierte Tabelle meldet den Zustand „nicht aktiviert".
func TestStatusNotActivated(t *testing.T) {
	service := status.NewGetStatusService(&fakeActivation{exists: true, registered: false})

	result, err := service.Status(context.Background(), inbound.GetStatusQuery{
		Source: "src-1", Schema: "public", Table: "t1",
	})
	if err != nil {
		t.Fatalf("Status: %v", err)
	}
	if result.Enabled {
		t.Fatalf("nie aktivierte Tabelle meldet „aktiviert“")
	}
}

// TestStatusMissingTable trägt den Negative-Pfad (`LH-FA-CFG-003`): eine
// nicht existierende Tabelle endet sichtbar.
func TestStatusMissingTable(t *testing.T) {
	service := status.NewGetStatusService(&fakeActivation{exists: false})

	_, err := service.Status(context.Background(), inbound.GetStatusQuery{
		Source: "src-1", Schema: "public", Table: "t1",
	})
	if !stderrors.Is(err, inbound.ErrSourceTableMissing) {
		t.Fatalf("fehlende Tabelle: %v (Erwartung: ErrSourceTableMissing)", err)
	}
}

// TestStatusEmptyQuery trägt die Kennungs-Grenze der Status-Abfrage.
func TestStatusEmptyQuery(t *testing.T) {
	service := status.NewGetStatusService(&fakeActivation{})

	_, err := service.Status(context.Background(), inbound.GetStatusQuery{})
	if !stderrors.Is(err, domainerrors.ErrEmptyIdentifier) {
		t.Fatalf("leere Eingabe: %v (Erwartung: ErrEmptyIdentifier)", err)
	}
}
