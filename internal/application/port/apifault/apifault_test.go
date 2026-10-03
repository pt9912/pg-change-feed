package apifault_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/pt9912/pg-change-feed/internal/application/port/apifault"
	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/messagecode"
)

// TestClassifyOrdnetJedenFehlerSeinerArtUndSeinemCodeZu trägt: jede Ablehnung
// trägt den Code ihrer Ursache im Bereich `E8`, die unbekannte Ressource den
// Code der fehlenden Tabelle, ein unerwarteter Fehler den Code seiner
// klassifizierten Ursache oder den Rückfall der Klasse `internal`. Färbt rot,
// sobald eine Zeile einen anderen Code oder eine andere Art liefert.
func TestClassifyOrdnetJedenFehlerSeinerArtUndSeinemCodeZu(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantKind apifault.Kind
		wantCode messagecode.Code
	}{
		{"Tabelle fehlt (gewrappt)", fmt.Errorf("%w: public.orders", inbound.ErrSourceTableMissing), apifault.NotFound, messagecode.RejectedTableMissing},
		{"leere Kennung", domainerrors.ErrEmptyIdentifier, apifault.InvalidInput, messagecode.RejectedRequiredField},
		{"Position ohne Offset", domainerrors.ErrInvalidPosition, apifault.InvalidInput, messagecode.RejectedPositionInvalid},
		{"Position rückt nicht vor", domainerrors.ErrPositionRegression, apifault.InvalidInput, messagecode.RejectedPositionRegressed},
		{"Position einer anderen Quelle", domainerrors.ErrSourceMismatch, apifault.InvalidInput, messagecode.RejectedPositionSource},
		{"Endposition vor Startposition", outbound.ErrRangeInverted, apifault.InvalidInput, messagecode.RejectedRangeInverted},
		{"negative Dauer", domainerrors.ErrNegativeDuration, apifault.InvalidInput, messagecode.RejectedValueInvalid},
		{"Version kleiner 1", domainerrors.ErrNonPositiveVersion, apifault.InvalidInput, messagecode.RejectedValueInvalid},
		{"Limit kleiner 1", outbound.ErrNonPositiveLimit, apifault.InvalidInput, messagecode.RejectedValueInvalid},
		{"klassifizierte Ursache", fmt.Errorf("lesen: %w", outbound.ErrStorage), apifault.Internal, messagecode.ChangeStoreFailed},
		{"unklassifizierter Fehler", errors.New("speicher nicht erreichbar"), apifault.Internal, messagecode.InternalFallback},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			kind, code := apifault.Classify(tc.err)
			if kind != tc.wantKind || code != tc.wantCode {
				t.Fatalf("Classify = (%d, %s), Erwartung (%d, %s)", kind, code, tc.wantKind, tc.wantCode)
			}
		})
	}
}

// TestClassifyAblehnungTraegtKeineFehlerklasse trägt: jeder Code einer
// Ablehnung oder einer unbekannten Ressource liegt im Bereich `E8` (ohne
// Fehlerklasse), jeder Code eines unerwarteten Fehlers in einer Fehlerklasse.
func TestClassifyAblehnungTraegtKeineFehlerklasse(t *testing.T) {
	for _, err := range []error{
		inbound.ErrSourceTableMissing, domainerrors.ErrEmptyIdentifier, domainerrors.ErrInvalidPosition,
		domainerrors.ErrPositionRegression, domainerrors.ErrSourceMismatch, outbound.ErrRangeInverted,
		domainerrors.ErrNegativeDuration, domainerrors.ErrNonPositiveVersion, outbound.ErrNonPositiveLimit,
	} {
		_, code := apifault.Classify(err)
		if class, ok := messagecode.DigitClass(code); !ok || class != messagecode.ClassNone {
			t.Errorf("%v: Code %s liegt nicht im Bereich E8", err, code)
		}
	}
	_, code := apifault.Classify(outbound.ErrStorage)
	if messagecode.ClassOf(code) == messagecode.ClassNone {
		t.Errorf("unerwarteter Fehler: Code %s trägt keine Fehlerklasse", code)
	}
}
