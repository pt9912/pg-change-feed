package model

import (
	stderrors "errors"
	"testing"

	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

// NewTableSchema verlangt eine nichtleere Versions-Kennung, mindestens
// eine Spalte und nichtleere Spaltennamen (`SPEC-004`, `ADR-0015`
// Folgepflicht).
func TestNewTableSchemaHappyPath(t *testing.T) {
	schema, err := NewTableSchema("sv-1", []Column{
		{Name: "id", OID: 23},
		{Name: "name", OID: 25},
	})
	if err != nil {
		t.Fatalf("NewTableSchema: %v", err)
	}
	if schema.VersionID != "sv-1" {
		t.Fatalf("VersionID = %q, wollen sv-1", schema.VersionID)
	}
	if len(schema.Columns) != 2 || schema.Columns[0].Name != "id" || schema.Columns[1].OID != 25 {
		t.Fatalf("Columns = %+v", schema.Columns)
	}
}

func TestNewTableSchemaRejectsEmptyVersionID(t *testing.T) {
	if _, err := NewTableSchema("", []Column{{Name: "id", OID: 23}}); !stderrors.Is(err, domainerrors.ErrEmptyIdentifier) {
		t.Fatalf("Fehler = %v, wollen ErrEmptyIdentifier", err)
	}
}

func TestNewTableSchemaRejectsEmptyColumns(t *testing.T) {
	if _, err := NewTableSchema("sv-1", nil); !stderrors.Is(err, domainerrors.ErrEmptyColumns) {
		t.Fatalf("Fehler = %v, wollen ErrEmptyColumns", err)
	}
	if _, err := NewTableSchema("sv-1", []Column{}); !stderrors.Is(err, domainerrors.ErrEmptyColumns) {
		t.Fatalf("Fehler = %v, wollen ErrEmptyColumns", err)
	}
}

func TestNewTableSchemaRejectsEmptyColumnName(t *testing.T) {
	if _, err := NewTableSchema("sv-1", []Column{{Name: "", OID: 23}}); !stderrors.Is(err, domainerrors.ErrEmptyIdentifier) {
		t.Fatalf("Fehler = %v, wollen ErrEmptyIdentifier", err)
	}
}

// Der Konstruktor kopiert die Spaltenmenge — eine spätere Änderung am
// übergebenen Slice trägt das bereits konstruierte TableSchema nicht.
func TestNewTableSchemaCopiesColumns(t *testing.T) {
	columns := []Column{{Name: "id", OID: 23}}
	schema, err := NewTableSchema("sv-1", columns)
	if err != nil {
		t.Fatalf("NewTableSchema: %v", err)
	}
	columns[0].Name = "geändert"
	if schema.Columns[0].Name != "id" {
		t.Fatalf("Columns[0].Name = %q, wollen unverändert id", schema.Columns[0].Name)
	}
}
