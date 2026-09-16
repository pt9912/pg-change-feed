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
// Operationen tragen die Schnittstelle. Die `…ErrFor`-Tabellen binden einen
// Port-Fehler an genau den Eingabewert, auf den er antwortet (LP2): jede
// andere Adresse, Quelle oder Publication trägt derselbe Fake.
type fakeActivation struct {
	exists     bool
	registered bool
	removal    outbound.ActivationRemoval

	existsNotFor     map[string]bool  // schema.table → existiert nicht
	existsErrFor     map[string]error // schema.table → Fehler
	unpublishErrFor  map[string]error // Publication → Fehler
	registeredErrFor map[string]error // Quell-Kennung → Fehler
	unregisterErrFor map[string]error // schema.table der Bindungs-Zeile → Fehler

	unpublishCalls  int
	unregisterCalls int
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
	if err, ok := f.registeredErrFor[string(source)]; ok {
		return model.SourceTable{}, false, err
	}
	registered, err := model.NewSourceTable("tbl-1", source, schema, table)
	if err != nil {
		return model.SourceTable{}, false, err
	}
	return registered, f.registered, nil
}

func (f *fakeActivation) Register(ctx context.Context, table model.SourceTable, version model.SchemaVersion) (bool, error) {
	return true, nil
}

func (f *fakeActivation) Unregister(ctx context.Context, table model.SourceTable) (outbound.ActivationRemoval, error) {
	f.unregisterCalls++
	if err, ok := f.unregisterErrFor[table.QualifiedName()]; ok {
		return "", err
	}
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
	if err, ok := f.unpublishErrFor[publication]; ok {
		return err
	}
	return nil
}

func (f *fakeActivation) Published(ctx context.Context, publication, schema, table string) (bool, error) {
	return false, nil
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
// zu berühren. Die Ablehnung ist an den Eingabewert gebunden: der Fake
// kennt genau die fehlende Adresse, dieselbe Anlage trägt die vorhandene
// Adresse durch.
func TestDisableMissingTable(t *testing.T) {
	fake := &fakeActivation{exists: true, existsNotFor: map[string]bool{"public.missing": true}, registered: true, removal: outbound.ActivationRemoved}
	service := disable.NewDisableTableService(fake)

	_, err := service.Disable(context.Background(), disableCommand("missing"))
	if !stderrors.Is(err, inbound.ErrSourceTableMissing) {
		t.Fatalf("fehlende Tabelle: %v (Erwartung: ErrSourceTableMissing)", err)
	}
	if fake.unpublishCalls != 0 {
		t.Fatalf("fehlende Tabelle trägt den Publication-Entzug: %d", fake.unpublishCalls)
	}

	if _, err := service.Disable(context.Background(), disableCommand("t1")); err != nil {
		t.Fatalf("vorhandene Tabelle: %v", err)
	}
	if fake.unpublishCalls != 1 || fake.unregisterCalls != 1 {
		t.Fatalf("Entzugs-Züge der vorhandenen Tabelle: Unpublish %d, Unregister %d (Erwartung: je 1)", fake.unpublishCalls, fake.unregisterCalls)
	}
}

// TestDisableTableExistsErrorFollowsAddress trägt den Fehlerpfad der
// Existenz-Prüfung: ein Fehler des Aktivierungs-Ports wird unverändert
// durchgereicht, ohne die Publication zu berühren. Die Ablehnung ist an die
// Kommando-Adresse gebunden — derselbe Fake trägt eine andere Adresse durch.
func TestDisableTableExistsErrorFollowsAddress(t *testing.T) {
	wantErr := stderrors.New("Katalog nicht lesbar")
	fake := &fakeActivation{exists: true, existsErrFor: map[string]error{"public.kaputt": wantErr}, registered: true, removal: outbound.ActivationRemoved}
	service := disable.NewDisableTableService(fake)

	_, err := service.Disable(context.Background(), disableCommand("kaputt"))
	if !stderrors.Is(err, wantErr) {
		t.Fatalf("Port-Fehler: %v (Erwartung: %v)", err, wantErr)
	}
	if fake.unpublishCalls != 0 {
		t.Fatalf("Existenz-Fehler trägt den Publication-Entzug: %d", fake.unpublishCalls)
	}

	if _, err := service.Disable(context.Background(), disableCommand("t1")); err != nil {
		t.Fatalf("Adresse public.t1: %v", err)
	}
}

// TestDisableUnpublishErrorFollowsPublication trägt den Fehlerpfad des
// Publication-Entzugs (`LH-FA-CFG-002`): ein Fehler von `Unpublish` wird
// unverändert durchgereicht, die Bindungs-Zeile bleibt danach unberührt. Die
// Ablehnung ist an die Publication des Kommandos gebunden — derselbe Fake
// entzieht eine andere Publication ohne Fehler.
func TestDisableUnpublishErrorFollowsPublication(t *testing.T) {
	wantErr := stderrors.New("Publication nicht entziehbar")
	fake := &fakeActivation{exists: true, registered: true, removal: outbound.ActivationRemoved, unpublishErrFor: map[string]error{"pub-kaputt": wantErr}}
	service := disable.NewDisableTableService(fake)

	command := disableCommand("t1")
	command.Publication = "pub-kaputt"
	_, err := service.Disable(context.Background(), command)
	if !stderrors.Is(err, wantErr) {
		t.Fatalf("Port-Fehler: %v (Erwartung: %v)", err, wantErr)
	}
	if fake.unregisterCalls != 0 {
		t.Fatalf("fehlgeschlagener Entzug trägt den Zeilen-Entzug: %d", fake.unregisterCalls)
	}

	if _, err := service.Disable(context.Background(), disableCommand("t1")); err != nil {
		t.Fatalf("Publication pub-1: %v", err)
	}
}

// TestDisableRegisteredErrorFollowsSource trägt den Fehlerpfad des
// Bindungs-Zeilen-Lesens: ein Fehler von `Registered` wird unverändert
// durchgereicht. Die Ablehnung ist an die Quell-Kennung des Kommandos
// gebunden — derselbe Fake liest eine andere Quelle ohne Fehler.
func TestDisableRegisteredErrorFollowsSource(t *testing.T) {
	wantErr := stderrors.New("Bindungs-Zeile nicht lesbar")
	fake := &fakeActivation{exists: true, registered: true, removal: outbound.ActivationRemoved, registeredErrFor: map[string]error{"src-kaputt": wantErr}}
	service := disable.NewDisableTableService(fake)

	command := disableCommand("t1")
	command.Source = "src-kaputt"
	_, err := service.Disable(context.Background(), command)
	if !stderrors.Is(err, wantErr) {
		t.Fatalf("Port-Fehler: %v (Erwartung: %v)", err, wantErr)
	}
	if fake.unregisterCalls != 0 {
		t.Fatalf("fehlgeschlagenes Lesen trägt den Zeilen-Entzug: %d", fake.unregisterCalls)
	}

	if _, err := service.Disable(context.Background(), disableCommand("t1")); err != nil {
		t.Fatalf("Quelle src-1: %v", err)
	}
}

// TestDisableUnregisterErrorFollowsBinding trägt den Fehlerpfad des
// Bindungs-Zeilen-Entzugs: ein Fehler von `Unregister` wird unverändert
// durchgereicht. Die Ablehnung ist an die Kommando-Adresse gebunden — der
// Fake entzieht genau die Bindungs-Zeile, die er zuvor aus dem Kommando
// gebaut hat, und trägt jede andere.
func TestDisableUnregisterErrorFollowsBinding(t *testing.T) {
	wantErr := stderrors.New("Bindungs-Zeile nicht entziehbar")
	fake := &fakeActivation{exists: true, registered: true, removal: outbound.ActivationRemoved, unregisterErrFor: map[string]error{"public.boom": wantErr}}
	service := disable.NewDisableTableService(fake)

	_, err := service.Disable(context.Background(), disableCommand("boom"))
	if !stderrors.Is(err, wantErr) {
		t.Fatalf("Port-Fehler: %v (Erwartung: %v)", err, wantErr)
	}

	result, err := service.Disable(context.Background(), disableCommand("t1"))
	if err != nil {
		t.Fatalf("Adresse public.t1: %v", err)
	}
	if !result.Removed {
		t.Fatalf("Adresse public.t1: %+v (Erwartung: Removed)", result)
	}
}

// TestDisableAbsentReport trägt den dritten Ausgang des Ports
// (`outbound.ActivationAbsent`): er wird nicht als entzogene oder belassene
// Bindungs-Zeile ausgegeben — die Rückkehr trägt keinen der beiden Ausgänge.
func TestDisableAbsentReport(t *testing.T) {
	fake := &fakeActivation{exists: true, registered: true, removal: outbound.ActivationAbsent}
	service := disable.NewDisableTableService(fake)

	result, err := service.Disable(context.Background(), disableCommand("t1"))
	if err != nil {
		t.Fatalf("Disable: %v", err)
	}
	if result.Removed || result.Retained {
		t.Fatalf("Ausgang nicht-aktiviert: %+v (Erwartung: ohne Zeilen-Ausgang)", result)
	}
	if fake.unregisterCalls != 1 {
		t.Fatalf("Entzugs-Züge: %d (Erwartung: 1)", fake.unregisterCalls)
	}
}

// disableCommand trägt die Deaktivierung der Tabelle `public.<table>` an der
// Quelle `src-1` in der Publication `pub-1`; der Test variiert über die
// Adresse.
func disableCommand(table string) inbound.DisableTableCommand {
	return inbound.DisableTableCommand{
		Source: "src-1", Schema: "public", Table: table, Publication: "pub-1",
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
