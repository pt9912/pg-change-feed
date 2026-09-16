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
// übrigen Operationen tragen die Schnittstelle. Die drei `…ErrFor`-Tabellen
// binden einen Port-Fehler an genau den Eingabewert, auf den er antwortet
// (LP2): jede andere Adresse, Kennung oder Publication trägt derselbe Fake.
type fakeActivation struct {
	exists        bool
	alreadyActive bool

	existsNotFor   map[string]bool  // schema.table → existiert nicht
	existsErrFor   map[string]error // schema.table → Fehler
	registerErrFor map[string]error // Bindungs-Zeilen-Kennung → Fehler
	publishErrFor  map[string]error // Publication → Fehler

	registerCalls int
	publishCalls  int
	lastTable     model.SourceTable
	lastVersion   model.SchemaVersion
}

func (f *fakeActivation) TableExists(ctx context.Context, schema, table string) (bool, error) {
	address := schema + "." + table
	if err, ok := f.existsErrFor[address]; ok {
		return false, err
	}
	if f.existsNotFor[address] {
		return false, nil
	}
	return f.exists, nil
}

func (f *fakeActivation) Registered(ctx context.Context, source model.SourceID, schema, table string) (model.SourceTable, bool, error) {
	return model.SourceTable{}, f.alreadyActive, nil
}

func (f *fakeActivation) Register(ctx context.Context, table model.SourceTable, version model.SchemaVersion) (bool, error) {
	f.registerCalls++
	f.lastTable = table
	f.lastVersion = version
	if err, ok := f.registerErrFor[string(table.ID)]; ok {
		return false, err
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
	if err, ok := f.publishErrFor[publication]; ok {
		return err
	}
	return nil
}

func (f *fakeActivation) Unpublish(ctx context.Context, publication, schema, table string) error {
	return nil
}

func (f *fakeActivation) Published(ctx context.Context, publication, schema, table string) (bool, error) {
	return true, nil
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
// oder Publication zu tragen. Die Ablehnung ist an den Eingabewert
// gebunden: der Fake kennt genau die fehlende Adresse, dieselbe Anlage
// trägt die vorhandene Adresse durch.
func TestEnableMissingTable(t *testing.T) {
	fake := &fakeActivation{exists: true, existsNotFor: map[string]bool{"public.missing": true}}
	service := enable.NewEnableTableService(fake)

	_, err := service.Enable(context.Background(), enableCommand("missing"))
	if !stderrors.Is(err, inbound.ErrSourceTableMissing) {
		t.Fatalf("fehlende Tabelle: %v (Erwartung: ErrSourceTableMissing)", err)
	}
	if fake.registerCalls != 0 || fake.publishCalls != 0 {
		t.Fatalf("fehlende Tabelle trägt Aktivierungs-Züge: Register %d, Publish %d", fake.registerCalls, fake.publishCalls)
	}

	if _, err := service.Enable(context.Background(), enableCommand("t1")); err != nil {
		t.Fatalf("vorhandene Tabelle: %v", err)
	}
	if fake.registerCalls != 1 || fake.publishCalls != 1 {
		t.Fatalf("Aktivierungs-Züge der vorhandenen Tabelle: Register %d, Publish %d (Erwartung: je 1)", fake.registerCalls, fake.publishCalls)
	}
}

// TestEnableTableExistsErrorFollowsAddress trägt den Fehlerpfad der
// Existenz-Prüfung: ein Fehler des Aktivierungs-Ports wird unverändert
// durchgereicht, ohne Bindungs-Zeilen oder Publication zu tragen. Die
// Ablehnung ist an die Kommando-Adresse gebunden — derselbe Fake trägt eine
// andere Adresse durch.
func TestEnableTableExistsErrorFollowsAddress(t *testing.T) {
	wantErr := stderrors.New("Katalog nicht lesbar")
	fake := &fakeActivation{exists: true, existsErrFor: map[string]error{"public.kaputt": wantErr}}
	service := enable.NewEnableTableService(fake)

	_, err := service.Enable(context.Background(), enableCommand("kaputt"))
	if !stderrors.Is(err, wantErr) {
		t.Fatalf("Port-Fehler: %v (Erwartung: %v)", err, wantErr)
	}
	if fake.registerCalls != 0 || fake.publishCalls != 0 {
		t.Fatalf("Existenz-Fehler trägt Aktivierungs-Züge: Register %d, Publish %d", fake.registerCalls, fake.publishCalls)
	}

	if _, err := service.Enable(context.Background(), enableCommand("t1")); err != nil {
		t.Fatalf("Adresse public.t1: %v", err)
	}
}

// TestEnableRegisterErrorFollowsTableID trägt den Fehlerpfad des
// Bindungs-Zeilen-Schreibens: ein Fehler von `Register` wird unverändert
// durchgereicht, die Publication läuft danach nicht. Die Ablehnung ist an
// die Tabellen-Kennung des Kommandos gebunden — derselbe Fake schreibt eine
// andere Kennung.
func TestEnableRegisterErrorFollowsTableID(t *testing.T) {
	wantErr := stderrors.New("Bindungs-Zeile nicht schreibbar")
	fake := &fakeActivation{exists: true, registerErrFor: map[string]error{"tbl-kaputt": wantErr}}
	service := enable.NewEnableTableService(fake)

	command := enableCommand("t1")
	command.TableID = "tbl-kaputt"
	_, err := service.Enable(context.Background(), command)
	if !stderrors.Is(err, wantErr) {
		t.Fatalf("Port-Fehler: %v (Erwartung: %v)", err, wantErr)
	}
	if fake.publishCalls != 0 {
		t.Fatalf("fehlgeschlagenes Register trägt die Publication: %d", fake.publishCalls)
	}

	if _, err := service.Enable(context.Background(), enableCommand("t1")); err != nil {
		t.Fatalf("Tabellen-Kennung tbl-1: %v", err)
	}
	if fake.lastTable.ID != "tbl-1" {
		t.Fatalf("Bindungs-Zeile trägt Kennung %q (Erwartung: tbl-1)", fake.lastTable.ID)
	}
}

// TestEnablePublishErrorFollowsPublication trägt den Fehlerpfad der
// Publication (`LH-FA-CFG-001.a`): ein Fehler von `Publish` wird
// unverändert durchgereicht; die zuvor geschriebene Bindungs-Zeile bleibt
// stehen (ein erneuter Aufruf trägt die Publication nach). Die Ablehnung ist
// an die Publication des Kommandos gebunden — derselbe Fake trägt eine
// andere Publication durch.
func TestEnablePublishErrorFollowsPublication(t *testing.T) {
	wantErr := stderrors.New("Publication nicht schreibbar")
	fake := &fakeActivation{exists: true, publishErrFor: map[string]error{"pub-kaputt": wantErr}}
	service := enable.NewEnableTableService(fake)

	command := enableCommand("t1")
	command.Publication = "pub-kaputt"
	_, err := service.Enable(context.Background(), command)
	if !stderrors.Is(err, wantErr) {
		t.Fatalf("Port-Fehler: %v (Erwartung: %v)", err, wantErr)
	}
	if fake.registerCalls != 1 {
		t.Fatalf("Publication-Fehler ohne Bindungs-Zeile: Register %d (Erwartung: 1)", fake.registerCalls)
	}

	if _, err := service.Enable(context.Background(), enableCommand("t1")); err != nil {
		t.Fatalf("Publication pub-1: %v", err)
	}
}

// enableCommand trägt die Aktivierung der Tabelle `public.<table>` an der
// Quelle `src-1`; der Test variiert über die Adresse.
func enableCommand(table string) inbound.EnableTableCommand {
	return inbound.EnableTableCommand{
		Source:          "src-1",
		Schema:          "public",
		Table:           table,
		TableID:         "tbl-1",
		SchemaVersionID: "sv-1",
		Version:         1,
		Publication:     "pub-1",
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
