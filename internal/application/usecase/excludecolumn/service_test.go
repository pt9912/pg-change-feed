package excludecolumn_test

import (
	"context"
	stderrors "errors"
	"testing"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/excludecolumn"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// fakeColumnExclusion trägt den `ColumnExclusionPort` als Fake (`ADR-0030`):
// er meldet die vorgegebene Spaltenexistenz und merkt sich die angefragte
// Adresse, damit der Test belegt, dass der Use Case genau die Spalte des
// Kommandos prüft. `existsNotFor` bindet die Ablehnung an genau eine
// Spaltenadresse (LP2): jede andere Adresse trägt derselbe Fake.
type fakeColumnExclusion struct {
	exists bool
	err    error

	existsNotFor map[string]bool // schema.table.column → existiert nicht

	calls  int
	schema string
	table  string
	column string
}

func (f *fakeColumnExclusion) ColumnExists(ctx context.Context, schema, table, column string) (bool, error) {
	f.calls++
	f.schema, f.table, f.column = schema, table, column
	if f.existsNotFor[schema+"."+table+"."+column] {
		return false, nil
	}
	return f.exists, f.err
}

// ExcludedColumns bleibt ungenutzt: der Use Case liest die Spaltenexistenz;
// den dauerhaften Ausschlussstand trägt die Verdrahtung (`ADR-0065`).
func (f *fakeColumnExclusion) ExcludedColumns(ctx context.Context, source model.SourceID) (map[string][]string, error) {
	return nil, nil
}

// TestExcludeColumnChecksSourceColumn trägt den Happy Path (`LH-FA-CFG-005`):
// eine an der Quelle vorhandene Spalte wird geprüft, der Aufruf endet ohne
// Fehler.
func TestExcludeColumnChecksSourceColumn(t *testing.T) {
	columns := &fakeColumnExclusion{exists: true}
	service := excludecolumn.NewExcludeColumnService(columns)

	if err := service.Exclude(context.Background(), excludecolumn.ExcludeColumnCommand{
		Source: "src-1", Schema: "public", Table: "orders", Column: "secret",
	}); err != nil {
		t.Fatalf("Exclude = %v, wollen nil", err)
	}
	if columns.calls != 1 {
		t.Fatalf("ColumnExists-Aufrufe = %d, wollen 1", columns.calls)
	}
	if columns.schema != "public" || columns.table != "orders" || columns.column != "secret" {
		t.Fatalf("geprüfte Adresse = %s.%s.%s, wollen public.orders.secret", columns.schema, columns.table, columns.column)
	}
}

// TestExcludeColumnRejectsMissingSourceColumn trägt den Negative-Pfad
// (`LH-FA-CFG-005` Negative: „folgt ein expliziter Fehlerpfad"): eine nicht
// existierende Spalte endet über `ErrSourceColumnMissing` — ohne dieses
// Sentinel bliebe der Antrag still erfolgreich und wirkte nie. Die Ablehnung
// ist an die Spaltenadresse des Kommandos gebunden: der Fake kennt genau die
// fehlende Adresse, dieselbe Anlage trägt die vorhandene Spalte durch.
func TestExcludeColumnRejectsMissingSourceColumn(t *testing.T) {
	columns := &fakeColumnExclusion{exists: true, existsNotFor: map[string]bool{"public.orders.does_not_exist": true}}
	service := excludecolumn.NewExcludeColumnService(columns)

	err := service.Exclude(context.Background(), excludecolumn.ExcludeColumnCommand{
		Source: "src-1", Schema: "public", Table: "orders", Column: "does_not_exist",
	})
	if !stderrors.Is(err, inbound.ErrSourceColumnMissing) {
		t.Fatalf("Fehler = %v, wollen ErrSourceColumnMissing", err)
	}
	if columns.schema != "public" || columns.table != "orders" || columns.column != "does_not_exist" {
		t.Fatalf("geprüfte Adresse = %s.%s.%s, wollen public.orders.does_not_exist", columns.schema, columns.table, columns.column)
	}

	if err := service.Exclude(context.Background(), excludecolumn.ExcludeColumnCommand{
		Source: "src-1", Schema: "public", Table: "orders", Column: "secret",
	}); err != nil {
		t.Fatalf("vorhandene Spalte: %v", err)
	}
}

// TestExcludeColumnPropagatesPortError trägt den Adapter-Fehlerpfad: ein
// Fehler der Katalog-Prüfung wird unverändert durchgereicht, nicht als
// fehlende Spalte fehlinterpretiert.
func TestExcludeColumnPropagatesPortError(t *testing.T) {
	wantErr := stderrors.New("Katalog nicht lesbar")
	service := excludecolumn.NewExcludeColumnService(&fakeColumnExclusion{err: wantErr})

	err := service.Exclude(context.Background(), excludecolumn.ExcludeColumnCommand{
		Source: "src-1", Schema: "public", Table: "orders", Column: "secret",
	})
	if !stderrors.Is(err, wantErr) {
		t.Fatalf("Fehler = %v, wollen %v", err, wantErr)
	}
	if stderrors.Is(err, inbound.ErrSourceColumnMissing) {
		t.Fatalf("Fehler = %v, darf nicht als fehlende Spalte gelesen werden", err)
	}
}
