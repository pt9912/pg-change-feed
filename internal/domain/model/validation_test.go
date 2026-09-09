package model

import (
	"testing"

	stderrors "errors"

	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

// Die verbleibenden Konstruktoren erzwingen die Nichtleere ihrer Kennungen
// (`LH-FA-DAT-002`, `LH-FA-SCH-005`, `ADR-0029`, Regel 7).
func TestValueObjectConstructorsRejectInvariantViolations(t *testing.T) {
	t.Run("Source ohne Kennung", func(t *testing.T) {
		if _, err := NewSource("", "quelle"); !stderrors.Is(err, domainerrors.ErrEmptyIdentifier) {
			t.Fatalf("Fehler = %v, wollen ErrEmptyIdentifier", err)
		}
		if _, err := NewSource("src-1", ""); !stderrors.Is(err, domainerrors.ErrEmptyIdentifier) {
			t.Fatalf("Fehler = %v, wollen ErrEmptyIdentifier", err)
		}
	})
	t.Run("SourceTable ohne Kennung", func(t *testing.T) {
		if _, err := NewSourceTable("", "src-1", "public", "orders"); !stderrors.Is(err, domainerrors.ErrEmptyIdentifier) {
			t.Fatalf("Fehler = %v, wollen ErrEmptyIdentifier", err)
		}
		if _, err := NewSourceTable("tbl-1", "src-1", "", "orders"); !stderrors.Is(err, domainerrors.ErrEmptyIdentifier) {
			t.Fatalf("Fehler = %v, wollen ErrEmptyIdentifier", err)
		}
	})
	t.Run("SchemaVersion ohne Versionsnummer", func(t *testing.T) {
		if _, err := NewSchemaVersion("sv-1", "tbl-1", 0); !stderrors.Is(err, domainerrors.ErrNonPositiveVersion) {
			t.Fatalf("Fehler = %v, wollen ErrNonPositiveVersion", err)
		}
	})
	t.Run("Consumer ohne Kennung", func(t *testing.T) {
		if _, err := NewConsumer("", "reports"); !stderrors.Is(err, domainerrors.ErrEmptyIdentifier) {
			t.Fatalf("Fehler = %v, wollen ErrEmptyIdentifier", err)
		}
	})
}

// LH-FA-DAT-002: die Quelltabelle ist über Schema und Tabellenname
// identifizierbar; gleichnamige Tabellen in verschiedenen Schemata sind
// unterscheidbar.
func TestLHFADAT002SourceTableIsIdentifiable(t *testing.T) {
	publicOrders, err := NewSourceTable("tbl-1", "src-1", "public", "orders")
	if err != nil {
		t.Fatalf("Tabelle public.orders: %v", err)
	}
	salesOrders, err := NewSourceTable("tbl-2", "src-1", "sales", "orders")
	if err != nil {
		t.Fatalf("Tabelle sales.orders: %v", err)
	}
	if publicOrders.QualifiedName() == salesOrders.QualifiedName() {
		t.Fatal("gleichnamige Tabellen in verschiedenen Schemata sind nicht unterscheidbar")
	}
}
