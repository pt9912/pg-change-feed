package disable_test

import (
	"context"
	stderrors "errors"
	"testing"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/disable"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// fakeActivation trägt den Aktivierungs-Port als Fake (`ADR-0030`); die
// drei Entzugs-Züge melden ihre Aufrufe, die Stubs der übrigen
// Operationen tragen die Schnittstelle.
type fakeActivation struct {
	exists      bool
	existsErr   error
	registered  bool
	registerErr error
	removal     outbound.ActivationRemoval

	unpublishCalls  int
	unregisterCalls int
}

func (f *fakeActivation) TableExists(ctx context.Context, schema, table string) (bool, error) {
	return f.exists, f.existsErr
}

func (f *fakeActivation) Registered(ctx context.Context, source model.SourceID, schema, table string) (model.SourceTable, bool, error) {
	registered, err := model.NewSourceTable("tbl-1", source, schema, table)
	if err != nil {
		return model.SourceTable{}, false, err
	}
	return registered, f.registered, f.registerErr
}

func (f *fakeActivation) Register(ctx context.Context, table model.SourceTable, version model.SchemaVersion) (bool, error) {
	return true, nil
}

func (f *fakeActivation) Unregister(ctx context.Context, table model.SourceTable) (outbound.ActivationRemoval, error) {
	f.unregisterCalls++
	return f.removal, nil
}

func (f *fakeActivation) List(ctx context.Context, source model.SourceID) ([]model.SourceTable, error) {
	return nil, nil
}

func (f *fakeActivation) Publish(ctx context.Context, publication, schema, table string) error {
	return nil
}

func (f *fakeActivation) Unpublish(ctx context.Context, publication, schema, table string) error {
	f.unpublishCalls++
	return nil
}

// TestDisableHappyPath trägt die Deaktivierung ohne Change-Bestand
// (`LH-FA-CFG-002` Happy Path): der Publication-Entzug trägt den Stopp
// der Erfassung, die Bindungs-Zeile wird entfernt.
func TestDisableHappyPath(t *testing.T) {
	fake := &fakeActivation{exists: true, registered: true, removal: outbound.ActivationRemoved}
	service := disable.NewDisableTableService(fake)

	result, err := service.Disable(context.Background(), inbound.DisableTableCommand{
		Source: "src-1", Schema: "public", Table: "t1", Publication: "pub-1",
	})
	if err != nil {
		t.Fatalf("Disable: %v", err)
	}
	if !result.Removed || result.Retained {
		t.Fatalf("Deaktivierung ohne Bestand: %+v (Erwartung: Removed)", result)
	}
	if fake.unpublishCalls != 1 || fake.unregisterCalls != 1 {
		t.Fatalf("Entzugs-Züge: Unpublish %d, Unregister %d (Erwartung: je 1)", fake.unpublishCalls, fake.unregisterCalls)
	}
}

// TestDisableRetained trägt die Deaktivierung mit Change-Bestand
// (`LH-FA-CFG-002` Out-of-Scope): die Bindungs-Zeile bleibt als Herkunft
// der persistierten Changes bestehen — kein Change wird gelöscht.
func TestDisableRetained(t *testing.T) {
	fake := &fakeActivation{exists: true, registered: true, removal: outbound.ActivationRetained}
	service := disable.NewDisableTableService(fake)

	result, err := service.Disable(context.Background(), inbound.DisableTableCommand{
		Source: "src-1", Schema: "public", Table: "t1", Publication: "pub-1",
	})
	if err != nil {
		t.Fatalf("Disable: %v", err)
	}
	if result.Removed || !result.Retained {
		t.Fatalf("Deaktivierung mit Bestand: %+v (Erwartung: Retained)", result)
	}
}

// TestDisableIdempotent trägt den Boundary-Pfad (`LH-FA-CFG-002`): eine
// nicht aktivierte Tabelle bleibt ohne Bindungs-Zeilen-Entzug — der
// Publication-Entzug läuft idempotent.
func TestDisableIdempotent(t *testing.T) {
	fake := &fakeActivation{exists: true, registered: false}
	service := disable.NewDisableTableService(fake)

	result, err := service.Disable(context.Background(), inbound.DisableTableCommand{
		Source: "src-1", Schema: "public", Table: "t1", Publication: "pub-1",
	})
	if err != nil {
		t.Fatalf("Disable: %v", err)
	}
	if result.Removed || result.Retained {
		t.Fatalf("Deaktivierung einer nie aktivierten Tabelle: %+v (Erwartung: ohne Zeilen-Ausgang)", result)
	}
	if fake.unpublishCalls != 1 || fake.unregisterCalls != 0 {
		t.Fatalf("Entzugs-Züge: Unpublish %d, Unregister %d (Erwartung: 1 und 0)", fake.unpublishCalls, fake.unregisterCalls)
	}
}

// TestDisableMissingTable trägt den Negative-Pfad (`LH-FA-CFG-002`):
// eine nicht existierende Tabelle endet sichtbar, ohne die Publication
// zu berühren.
func TestDisableMissingTable(t *testing.T) {
	fake := &fakeActivation{exists: false}
	service := disable.NewDisableTableService(fake)

	_, err := service.Disable(context.Background(), inbound.DisableTableCommand{
		Source: "src-1", Schema: "public", Table: "t1", Publication: "pub-1",
	})
	if !stderrors.Is(err, inbound.ErrSourceTableMissing) {
		t.Fatalf("fehlende Tabelle: %v (Erwartung: ErrSourceTableMissing)", err)
	}
	if fake.unpublishCalls != 0 {
		t.Fatalf("fehlende Tabelle trägt den Publication-Entzug: %d", fake.unpublishCalls)
	}
}

// TestDisableEmptyCommand trägt die Kennungs-Grenze der Deaktivierung.
func TestDisableEmptyCommand(t *testing.T) {
	service := disable.NewDisableTableService(&fakeActivation{})

	_, err := service.Disable(context.Background(), inbound.DisableTableCommand{})
	if !stderrors.Is(err, domainerrors.ErrEmptyIdentifier) {
		t.Fatalf("leere Eingabe: %v (Erwartung: ErrEmptyIdentifier)", err)
	}
}
