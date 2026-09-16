package includecolumn_test

import (
	"context"
	stderrors "errors"
	"testing"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/includecolumn"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// fakeColumnExclusion trägt den `ColumnExclusionPort` als Fake (`ADR-0030`)
// — dieselbe Form wie im Ausschluss-Test, hier für die Gegenrichtung;
// `existsNotFor` bindet die Ablehnung an genau eine Spaltenadresse (LP2).
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

// TestIncludeColumnChecksSourceColumn trägt den Happy Path (`LH-FA-CFG-005`):
// die Spalte existiert an der Quelle, der Aufruf endet ohne Fehler.
func TestIncludeColumnChecksSourceColumn(t *testing.T) {
	columns := &fakeColumnExclusion{exists: true}
	service := includecolumn.NewIncludeColumnService(columns)

	if err := service.Include(context.Background(), includecolumn.IncludeColumnCommand{
		Source: "src-1", Schema: "public", Table: "orders", Column: "secret",
	}); err != nil {
		t.Fatalf("Include = %v, wollen nil", err)
	}
	if columns.calls != 1 {
		t.Fatalf("ColumnExists-Aufrufe = %d, wollen 1", columns.calls)
	}
}

// TestIncludeColumnRejectsMissingSourceColumn trägt den Negative-Pfad
// (`LH-FA-CFG-005` Negative): dieselbe Vorbedingung wie der Ausschluss — die
// Spalte muss an der Quelle existieren. Die Ablehnung ist an die
// Spaltenadresse des Kommandos gebunden: der Fake kennt genau die fehlende
// Adresse, dieselbe Anlage trägt die vorhandene Spalte durch.
func TestIncludeColumnRejectsMissingSourceColumn(t *testing.T) {
	columns := &fakeColumnExclusion{exists: true, existsNotFor: map[string]bool{"public.orders.does_not_exist": true}}
	service := includecolumn.NewIncludeColumnService(columns)

	err := service.Include(context.Background(), includecolumn.IncludeColumnCommand{
		Source: "src-1", Schema: "public", Table: "orders", Column: "does_not_exist",
	})
	if !stderrors.Is(err, inbound.ErrSourceColumnMissing) {
		t.Fatalf("Fehler = %v, wollen ErrSourceColumnMissing", err)
	}
	if columns.schema != "public" || columns.table != "orders" || columns.column != "does_not_exist" {
		t.Fatalf("geprüfte Adresse = %s.%s.%s, wollen public.orders.does_not_exist", columns.schema, columns.table, columns.column)
	}

	if err := service.Include(context.Background(), includecolumn.IncludeColumnCommand{
		Source: "src-1", Schema: "public", Table: "orders", Column: "secret",
	}); err != nil {
		t.Fatalf("vorhandene Spalte: %v", err)
	}
}

// TestIncludeColumnPropagatesPortError trägt den Adapter-Fehlerpfad: ein
// Fehler der Katalog-Prüfung wird unverändert durchgereicht, nicht als
// fehlende Spalte fehlinterpretiert.
func TestIncludeColumnPropagatesPortError(t *testing.T) {
	wantErr := stderrors.New("Katalog nicht lesbar")
	service := includecolumn.NewIncludeColumnService(&fakeColumnExclusion{err: wantErr})

	err := service.Include(context.Background(), includecolumn.IncludeColumnCommand{
		Source: "src-1", Schema: "public", Table: "orders", Column: "secret",
	})
	if !stderrors.Is(err, wantErr) {
		t.Fatalf("Fehler = %v, wollen %v", err, wantErr)
	}
	if stderrors.Is(err, inbound.ErrSourceColumnMissing) {
		t.Fatalf("Fehler = %v, darf nicht als fehlende Spalte gelesen werden", err)
	}
}
