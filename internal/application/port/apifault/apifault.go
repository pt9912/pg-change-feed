// Package apifault ordnet den Fehler eines Aufrufs der HTTP- oder gRPC-API
// seiner Art und seinem Meldungscode zu. Beide Driving-Adapter fragen
// dieselbe Stelle, deshalb meldet ein Fehler auf beiden Wegen denselben Code
// (`ADR-0144`); der Statuscode der Antwort bleibt Sache des Adapters.
package apifault

import (
	"errors"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/messagecode"
)

// Kind ist die Art eines Fehlers aus Sicht des Aufrufers.
type Kind int

const (
	// Internal ist ein unerwarteter Fehler (HTTP `500`, gRPC `Internal`).
	Internal Kind = iota
	// InvalidInput ist eine abgelehnte Aufrufer-Eingabe (HTTP `400`, gRPC
	// `InvalidArgument`).
	InvalidInput
	// NotFound ist eine unbekannte Ressource an der Quelle (HTTP `404`, gRPC
	// `NotFound`).
	NotFound
)

// Classify liefert Art und Meldungscode des Fehlers eines Use-Case-Aufrufs.
// Eine Ablehnung trägt einen Code des Bereichs `E8`; ein unerwarteter Fehler
// trägt den Code seiner Ursache, wenn die Kette einen klassifizierten
// Fehlerwert führt (`messagecode.From`), sonst den Rückfall der Klasse
// `internal`.
func Classify(err error) (Kind, messagecode.Code) {
	switch {
	case errors.Is(err, inbound.ErrSourceTableMissing):
		return NotFound, messagecode.RejectedTableMissing
	case errors.Is(err, domainerrors.ErrEmptyIdentifier):
		return InvalidInput, messagecode.RejectedRequiredField
	case errors.Is(err, domainerrors.ErrInvalidPosition):
		return InvalidInput, messagecode.RejectedPositionInvalid
	case errors.Is(err, domainerrors.ErrPositionRegression):
		return InvalidInput, messagecode.RejectedPositionRegressed
	case errors.Is(err, domainerrors.ErrSourceMismatch):
		return InvalidInput, messagecode.RejectedPositionSource
	case errors.Is(err, outbound.ErrRangeInverted):
		return InvalidInput, messagecode.RejectedRangeInverted
	case errors.Is(err, domainerrors.ErrNegativeDuration),
		errors.Is(err, domainerrors.ErrNonPositiveVersion),
		errors.Is(err, outbound.ErrNonPositiveLimit):
		return InvalidInput, messagecode.RejectedValueInvalid
	}
	if code, ok := messagecode.From(err); ok {
		return Internal, code
	}
	return Internal, messagecode.InternalFallback
}
