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
// Status-Abfrage liest Existenz, Bindungs-Zeile und
// Publication-Mitgliedschaft, die übrigen Operationen tragen die
// Schnittstelle. Die `…ErrFor`-Tabellen binden einen Port-Fehler an genau den
// Eingabewert, auf den er antwortet (LP2): jede andere Adresse, Quelle oder
// Publication trägt derselbe Fake.
type fakeActivation struct {
	exists     bool
	registered bool
	published  bool

	existsNotFor     map[string]bool  // schema.table → existiert nicht
	existsErrFor     map[string]error // schema.table → Fehler
	registeredErrFor map[string]error // Quell-Kennung → Fehler
	publishedErrFor  map[string]error // Publication → Fehler
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
	return model.SourceTable{}, f.registered, nil
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

func (f *fakeActivation) Published(ctx context.Context, publication, schema, table string) (bool, error) {
	if err, ok := f.publishedErrFor[publication]; ok {
		return false, err
	}
	return f.published, nil
}

// TestStatusEnabled trägt den Happy Path (`LH-FA-CFG-003`): die aktivierte
// Tabelle — Bindungs-Zeile samt Publication-Mitgliedschaft — meldet den
// Zustand „aktiviert".
func TestStatusEnabled(t *testing.T) {
	service := status.NewGetStatusService(&fakeActivation{exists: true, registered: true, published: true})

	result, err := service.Status(context.Background(), statusQuery("t1"))
	if err != nil {
		t.Fatalf("Status: %v", err)
	}
	if !result.Enabled || result.Retained {
		t.Fatalf("aktivierte Tabelle: %+v (Erwartung: Enabled)", result)
	}
}

// TestStatusRetained trägt den Zustand nach der Deaktivierung mit
// Change-Bestand: die Bindungs-Zeile trägt die Herkunft der persistierten
// Changes (`LH-FA-CFG-002` Out-of-Scope), die Publication trägt die
// Tabelle nicht mehr — die Abfrage trennt den Erfassungs-Zustand von der
// Herkunft statt die Zeile doppelt zu lesen.
func TestStatusRetained(t *testing.T) {
	service := status.NewGetStatusService(&fakeActivation{exists: true, registered: true, published: false})

	result, err := service.Status(context.Background(), statusQuery("t1"))
	if err != nil {
		t.Fatalf("Status: %v", err)
	}
	if result.Enabled || !result.Retained {
		t.Fatalf("Bindungs-Zeile ohne Erfassung: %+v (Erwartung: Retained)", result)
	}
}

// TestStatusNotActivated trägt den Boundary-Pfad (`LH-FA-CFG-003`): eine
// nie aktivierte Tabelle meldet den Zustand „nicht aktiviert".
func TestStatusNotActivated(t *testing.T) {
	service := status.NewGetStatusService(&fakeActivation{exists: true, registered: false})

	result, err := service.Status(context.Background(), statusQuery("t1"))
	if err != nil {
		t.Fatalf("Status: %v", err)
	}
	if result.Enabled || result.Retained {
		t.Fatalf("nie aktivierte Tabelle: %+v (Erwartung: ohne Zustand)", result)
	}
}

// TestStatusMissingTable trägt den Negative-Pfad (`LH-FA-CFG-003`): eine
// nicht existierende Tabelle endet sichtbar. Die Ablehnung ist an den
// Eingabewert gebunden: der Fake kennt genau die fehlende Adresse, dieselbe
// Anlage trägt die vorhandene Adresse durch.
func TestStatusMissingTable(t *testing.T) {
	fake := &fakeActivation{exists: true, existsNotFor: map[string]bool{"public.missing": true}, registered: true, published: true}
	service := status.NewGetStatusService(fake)

	_, err := service.Status(context.Background(), statusQuery("missing"))
	if !stderrors.Is(err, inbound.ErrSourceTableMissing) {
		t.Fatalf("fehlende Tabelle: %v (Erwartung: ErrSourceTableMissing)", err)
	}

	result, err := service.Status(context.Background(), statusQuery("t1"))
	if err != nil {
		t.Fatalf("vorhandene Tabelle: %v", err)
	}
	if !result.Enabled {
		t.Fatalf("Adresse public.t1: %+v (Erwartung: Enabled)", result)
	}
}

// TestStatusTableExistsErrorFollowsAddress trägt den Fehlerpfad der
// Existenz-Prüfung: ein Fehler des Aktivierungs-Ports wird unverändert
// durchgereicht. Die Ablehnung ist an die Kommando-Adresse gebunden —
// derselbe Fake trägt eine andere Adresse durch.
func TestStatusTableExistsErrorFollowsAddress(t *testing.T) {
	wantErr := stderrors.New("Katalog nicht lesbar")
	fake := &fakeActivation{exists: true, existsErrFor: map[string]error{"public.kaputt": wantErr}, registered: true, published: true}
	service := status.NewGetStatusService(fake)

	_, err := service.Status(context.Background(), statusQuery("kaputt"))
	if !stderrors.Is(err, wantErr) {
		t.Fatalf("Port-Fehler: %v (Erwartung: %v)", err, wantErr)
	}

	if _, err := service.Status(context.Background(), statusQuery("t1")); err != nil {
		t.Fatalf("Adresse public.t1: %v", err)
	}
}

// TestStatusRegisteredErrorFollowsSource trägt den Fehlerpfad des
// Bindungs-Zeilen-Lesens: ein Fehler von `Registered` wird unverändert
// durchgereicht. Die Ablehnung ist an die Quell-Kennung gebunden — derselbe
// Fake liest eine andere Quelle ohne Fehler.
func TestStatusRegisteredErrorFollowsSource(t *testing.T) {
	wantErr := stderrors.New("Bindungs-Zeile nicht lesbar")
	fake := &fakeActivation{exists: true, registered: true, published: true, registeredErrFor: map[string]error{"src-kaputt": wantErr}}
	service := status.NewGetStatusService(fake)

	query := statusQuery("t1")
	query.Source = "src-kaputt"
	_, err := service.Status(context.Background(), query)
	if !stderrors.Is(err, wantErr) {
		t.Fatalf("Port-Fehler: %v (Erwartung: %v)", err, wantErr)
	}

	if _, err := service.Status(context.Background(), statusQuery("t1")); err != nil {
		t.Fatalf("Quelle src-1: %v", err)
	}
}

// TestStatusPublishedErrorFollowsPublication trägt den Fehlerpfad der
// Publication-Mitgliedschaft: ein Fehler von `Published` wird unverändert
// durchgereicht. Die Ablehnung ist an die Publication des Kommandos
// gebunden — derselbe Fake liest eine andere Publication ohne Fehler.
func TestStatusPublishedErrorFollowsPublication(t *testing.T) {
	wantErr := stderrors.New("Publication nicht lesbar")
	fake := &fakeActivation{exists: true, registered: true, published: true, publishedErrFor: map[string]error{"pub-kaputt": wantErr}}
	service := status.NewGetStatusService(fake)

	query := statusQuery("t1")
	query.Publication = "pub-kaputt"
	_, err := service.Status(context.Background(), query)
	if !stderrors.Is(err, wantErr) {
		t.Fatalf("Port-Fehler: %v (Erwartung: %v)", err, wantErr)
	}

	if _, err := service.Status(context.Background(), statusQuery("t1")); err != nil {
		t.Fatalf("Publication pub-1: %v", err)
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

// statusQuery trägt die Abfrage einer aktivierten Tabelle; der Test
// variiert nur den Tabellennamen.
func statusQuery(table string) inbound.GetStatusQuery {
	return inbound.GetStatusQuery{
		Source: "src-1", Schema: "public", Table: table, Publication: "pub-1",
	}
}
