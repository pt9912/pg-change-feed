package model

import (
	stderrors "errors"
	"testing"

	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

// NewAdministrationRequest trägt die Tabellen-Antragsarten der geschlossenen
// Menge (`LH-FA-ADM-001`, `LH-FA-CAP-009`); `enable`, `disable` und
// `backfill` tragen weder Spalte noch Regel. Rot färbende Mutation:
// `AdministrationRequestBackfill` aus dem `switch` des Konstruktors
// streichen — die Art endet als `ErrInvalidAdministrationRequestKind`.
func TestNewAdministrationRequestAcceptsTabellenAntragsarten(t *testing.T) {
	for _, kind := range []AdministrationRequestKind{
		AdministrationRequestEnable,
		AdministrationRequestDisable,
		AdministrationRequestBackfill,
	} {
		request, err := NewAdministrationRequest("req-1", "src-1", "public", "orders", "", "", "", kind)
		if err != nil {
			t.Fatalf("NewAdministrationRequest(%q): %v", kind, err)
		}
		if request.Kind != kind || request.Column != "" || request.RuleName != "" || request.RuleSpec != "" {
			t.Fatalf("Antrag = %+v, wollen Art %q ohne Spalte und ohne Regel", request, kind)
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
		request, err := NewAdministrationRequest("req-1", "src-1", "public", "orders", "secret", "", "", kind)
		if err != nil {
			t.Fatalf("NewAdministrationRequest(%q): %v", kind, err)
		}
		if request.Kind != kind || request.Column != "secret" {
			t.Fatalf("Antrag = %+v, wollen Art %q mit Spalte secret", request, kind)
		}
	}
}

// Die beiden Transformations-Antragsarten tragen Regelname und — nur
// `set_transformation` — die Regelform (`LH-FA-CFG-007`, `ADR-0112`
// Teilfrage 1). Rot färbende Mutation: `AdministrationRequestSetTransformation`
// bzw. `AdministrationRequestRemoveTransformation` aus dem `switch` des
// Konstruktors streichen — die Art endet als
// `ErrInvalidAdministrationRequestKind`.
func TestNewAdministrationRequestAcceptsTransformationsAntragsarten(t *testing.T) {
	set, err := NewAdministrationRequest("req-1", "src-1", "public", "orders", "", "umbenennung", `{"kind": "rename_column"}`, AdministrationRequestSetTransformation)
	if err != nil {
		t.Fatalf("NewAdministrationRequest(set_transformation): %v", err)
	}
	if set.Kind != AdministrationRequestSetTransformation || set.RuleName != "umbenennung" || set.RuleSpec != `{"kind": "rename_column"}` || set.Column != "" {
		t.Fatalf("Antrag = %+v, wollen set_transformation mit Regelname und Regelform, ohne Spalte", set)
	}
	remove, err := NewAdministrationRequest("req-2", "src-1", "public", "orders", "", "umbenennung", "", AdministrationRequestRemoveTransformation)
	if err != nil {
		t.Fatalf("NewAdministrationRequest(remove_transformation): %v", err)
	}
	if remove.Kind != AdministrationRequestRemoveTransformation || remove.RuleName != "umbenennung" || remove.RuleSpec != "" {
		t.Fatalf("Antrag = %+v, wollen remove_transformation mit Regelname, ohne Regelform", remove)
	}
}

// Die Ablehnungszweige des Konstruktors (`LH-FA-ADM-001`,
// `LH-FA-CFG-005`, `LH-FA-CFG-007`): leere Kennungen, die Spalten-Antragsart
// ohne Spalte, die Transformations-Antragsarten ohne Regelname bzw.
// `set_transformation` ohne Regelform und eine Antragsart außerhalb der
// geschlossenen Menge.
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
			if _, err := NewAdministrationRequest(tc.id, tc.source, tc.schema, tc.table, "", "", "", AdministrationRequestEnable); !stderrors.Is(err, domainerrors.ErrEmptyIdentifier) {
				t.Fatalf("leere %s: Fehler = %v, wollen ErrEmptyIdentifier", tc.name, err)
			}
		}
	})
	t.Run("Spalten-Antragsart ohne Spalte", func(t *testing.T) {
		for _, kind := range []AdministrationRequestKind{
			AdministrationRequestExcludeColumn,
			AdministrationRequestIncludeColumn,
		} {
			if _, err := NewAdministrationRequest("req-1", "src-1", "public", "orders", "", "", "", kind); !stderrors.Is(err, domainerrors.ErrEmptyIdentifier) {
				t.Fatalf("Art %q ohne Spalte: Fehler = %v, wollen ErrEmptyIdentifier", kind, err)
			}
		}
	})
	t.Run("Transformations-Antragsart ohne Regelname", func(t *testing.T) {
		for _, kind := range []AdministrationRequestKind{
			AdministrationRequestSetTransformation,
			AdministrationRequestRemoveTransformation,
		} {
			if _, err := NewAdministrationRequest("req-1", "src-1", "public", "orders", "", "", `{"kind": "rename_column"}`, kind); !stderrors.Is(err, domainerrors.ErrEmptyIdentifier) {
				t.Fatalf("Art %q ohne Regelname: Fehler = %v, wollen ErrEmptyIdentifier", kind, err)
			}
		}
	})
	t.Run("set_transformation ohne Regelform", func(t *testing.T) {
		if _, err := NewAdministrationRequest("req-1", "src-1", "public", "orders", "", "umbenennung", "", AdministrationRequestSetTransformation); !stderrors.Is(err, domainerrors.ErrEmptyIdentifier) {
			t.Fatalf("set_transformation ohne Regelform: Fehler = %v, wollen ErrEmptyIdentifier", err)
		}
	})
	t.Run("Antragsart außerhalb der geschlossenen Menge", func(t *testing.T) {
		for _, kind := range []AdministrationRequestKind{"", "truncate"} {
			if _, err := NewAdministrationRequest("req-1", "src-1", "public", "orders", "secret", "umbenennung", `{"kind": "rename_column"}`, kind); !stderrors.Is(err, domainerrors.ErrInvalidAdministrationRequestKind) {
				t.Fatalf("Art %q: Fehler = %v, wollen ErrInvalidAdministrationRequestKind", kind, err)
			}
		}
	})
}
