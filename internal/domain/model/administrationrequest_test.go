package model

import (
	stderrors "errors"
	"testing"

	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

// NewAdministrationRequest trägt die fünf Antragsarten der geschlossenen
// Menge (`LH-FA-ADM-001`, `LH-FA-CFG-005`, `LH-FA-CAP-009`); die beiden
// Tabellen-Antragsarten und `backfill` tragen keine Spalte. Rot färbende
// Mutation: `AdministrationRequestBackfill` aus dem `switch` des
// Konstruktors streichen — die Art endet als `ErrInvalidAdministrationRequestKind`.
func TestNewAdministrationRequestAcceptsTabellenAntragsarten(t *testing.T) {
	for _, kind := range []AdministrationRequestKind{
		AdministrationRequestEnable,
		AdministrationRequestDisable,
		AdministrationRequestBackfill,
	} {
		request, err := NewAdministrationRequest("req-1", "src-1", "public", "orders", "", kind)
		if err != nil {
			t.Fatalf("NewAdministrationRequest(%q): %v", kind, err)
		}
		if request.Kind != kind || request.Column != "" {
			t.Fatalf("Antrag = %+v, wollen Art %q ohne Spalte", request, kind)
		}
	}
}

// Die beiden Spalten-Antragsarten tragen den Ziel-Spaltennamen
// (`LH-FA-CFG-005`).
func TestNewAdministrationRequestAcceptsSpaltenAntragsarten(t *testing.T) {
	for _, kind := range []AdministrationRequestKind{
		AdministrationRequestExcludeColumn,
		AdministrationRequestIncludeColumn,
	} {
		request, err := NewAdministrationRequest("req-1", "src-1", "public", "orders", "secret", kind)
		if err != nil {
			t.Fatalf("NewAdministrationRequest(%q): %v", kind, err)
		}
		if request.Kind != kind || request.Column != "secret" {
			t.Fatalf("Antrag = %+v, wollen Art %q mit Spalte secret", request, kind)
		}
	}
}

// Die Ablehnungszweige des Konstruktors (`LH-FA-ADM-001`,
// `LH-FA-CFG-005`): leere Kennungen, die Spalten-Antragsart ohne Spalte und
// eine Antragsart außerhalb der geschlossenen Menge.
func TestNewAdministrationRequestRejectsInvariantViolations(t *testing.T) {
	t.Run("leere Kennung", func(t *testing.T) {
		for _, tc := range []struct {
			name   string
			id     AdministrationRequestID
			source SourceID
			schema string
			table  string
		}{
			{"ID", "", "src-1", "public", "orders"},
			{"Quelle", "req-1", "", "public", "orders"},
			{"Schema", "req-1", "src-1", "", "orders"},
			{"Tabelle", "req-1", "src-1", "public", ""},
		} {
			if _, err := NewAdministrationRequest(tc.id, tc.source, tc.schema, tc.table, "", AdministrationRequestEnable); !stderrors.Is(err, domainerrors.ErrEmptyIdentifier) {
				t.Fatalf("leere %s: Fehler = %v, wollen ErrEmptyIdentifier", tc.name, err)
			}
		}
	})
	t.Run("Spalten-Antragsart ohne Spalte", func(t *testing.T) {
		for _, kind := range []AdministrationRequestKind{
			AdministrationRequestExcludeColumn,
			AdministrationRequestIncludeColumn,
		} {
			if _, err := NewAdministrationRequest("req-1", "src-1", "public", "orders", "", kind); !stderrors.Is(err, domainerrors.ErrEmptyIdentifier) {
				t.Fatalf("Art %q ohne Spalte: Fehler = %v, wollen ErrEmptyIdentifier", kind, err)
			}
		}
	})
	t.Run("Antragsart außerhalb der geschlossenen Menge", func(t *testing.T) {
		for _, kind := range []AdministrationRequestKind{"", "truncate"} {
			if _, err := NewAdministrationRequest("req-1", "src-1", "public", "orders", "secret", kind); !stderrors.Is(err, domainerrors.ErrInvalidAdministrationRequestKind) {
				t.Fatalf("Art %q: Fehler = %v, wollen ErrInvalidAdministrationRequestKind", kind, err)
			}
		}
	})
}
