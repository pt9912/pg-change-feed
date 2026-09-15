package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

// writeDomainError bildet einen Use-Case-Fehler auf einen HTTP-Statuscode ab:
// `inbound.ErrSourceTableMissing` trägt eine unbekannte Ressource an der
// Quelle (`404`); die benannten Domänen-Invarianten (`ADR-0029`) und die
// Kontrakt-Sentinels des Leseports (`outbound.ErrNonPositiveLimit`,
// `outbound.ErrRangeInverted`, `LH-FA-REA-001`/`003` Negative) tragen eine
// ungültige Eingabe (`400`) — dieselbe Klasse wie
// `domainerrors.ErrEmptyIdentifier` in `registerconsumer.go`; jeder übrige
// (unbekannte) Fehler ist ein unerwarteter interner Fehler (`500`) und wird
// geloggt, nicht nur verworfen.
func writeDomainError(ctx context.Context, w http.ResponseWriter, log outbound.LogPort, action string, err error) {
	switch {
	case errors.Is(err, inbound.ErrSourceTableMissing):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domainerrors.ErrEmptyIdentifier),
		errors.Is(err, domainerrors.ErrInvalidPosition),
		errors.Is(err, domainerrors.ErrSourceMismatch),
		errors.Is(err, domainerrors.ErrPositionRegression),
		errors.Is(err, domainerrors.ErrNegativeDuration),
		errors.Is(err, domainerrors.ErrNonPositiveVersion),
		errors.Is(err, outbound.ErrNonPositiveLimit),
		errors.Is(err, outbound.ErrRangeInverted):
		writeError(w, http.StatusBadRequest, err.Error())
	default:
		log.Warn(ctx, "http: "+action+" fehlgeschlagen", "error", err)
		writeError(w, http.StatusInternalServerError, "interner Fehler")
	}
}
