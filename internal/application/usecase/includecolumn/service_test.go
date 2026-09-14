package includecolumn_test

import (
	"context"
	stderrors "errors"
	"testing"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/includecolumn"
)

// fakeColumnExclusion trägt den `ColumnExclusionPort` als Fake (`ADR-0030`)
// — dieselbe Form wie im Ausschluss-Test, hier für die Gegenrichtung.
type fakeColumnExclusion struct {
	exists bool
	err    error

	calls int
}

func (f *fakeColumnExclusion) ColumnExists(ctx context.Context, schema, table, column string) (bool, error) {
	f.calls++
	return f.exists, f.err
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
// Spalte muss an der Quelle existieren.
func TestIncludeColumnRejectsMissingSourceColumn(t *testing.T) {
	service := includecolumn.NewIncludeColumnService(&fakeColumnExclusion{exists: false})

	err := service.Include(context.Background(), includecolumn.IncludeColumnCommand{
		Source: "src-1", Schema: "public", Table: "orders", Column: "does_not_exist",
	})
	if !stderrors.Is(err, inbound.ErrSourceColumnMissing) {
		t.Fatalf("Fehler = %v, wollen ErrSourceColumnMissing", err)
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
