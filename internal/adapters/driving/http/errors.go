package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/pt9912/pg-change-feed/internal/application/port/apifault"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/domain/messagecode"
)

// writeDomainError bildet einen Use-Case-Fehler auf einen HTTP-Statuscode ab;
// Art und Meldungscode liefert `apifault.Classify`, dieselbe Stelle wie der
// gRPC-Adapter: eine unbekannte Ressource an der Quelle (`404`) und eine
// abgelehnte Eingabe (`400`) tragen einen Code des Bereichs `E8`; jeder
// übrige Fehler ist ein unerwarteter interner Fehler (`500`) und wird
// geloggt, nicht nur verworfen.
func writeDomainError(ctx context.Context, w http.ResponseWriter, log outbound.LogPort, action string, err error) {
	switch kind, code := apifault.Classify(err); kind {
	case apifault.NotFound:
		writeError(w, http.StatusNotFound, err.Error(), code)
	case apifault.InvalidInput:
		writeError(w, http.StatusBadRequest, err.Error(), code)
	default:
		writeInternalError(ctx, w, log, action, err)
	}
}

// writeInternalError antwortet mit `500` und dem Klartext „interner Fehler“;
// der Code im Körper ist der der Ursache (`apifault.Classify`), nie ein Code
// des Bereichs `E8`. Die Warnung trägt `PCF-W4008` und die Ursache im
// Attribut `error`.
func writeInternalError(ctx context.Context, w http.ResponseWriter, log outbound.LogPort, action string, err error) {
	log.Warn(ctx, "http: "+action+" fehlgeschlagen", messagecode.LogKey, messagecode.WarnAPIRequestFailed, "error", err)
	kind, code := apifault.Classify(err)
	if kind != apifault.Internal {
		code = messagecode.InternalFallback
	}
	writeError(w, http.StatusInternalServerError, "interner Fehler", code)
}

// paramError ist ein abgelehnter Query-Parameter mit seinem Meldungscode.
type paramError struct {
	code messagecode.Code
	text string
}

func (e *paramError) Error() string { return e.text }

// rejectParam legt den Fehler eines abgelehnten Query-Parameters an.
func rejectParam(code messagecode.Code, format string, args ...any) error {
	return &paramError{code: code, text: fmt.Sprintf(format, args...)}
}

// writeBadRequest antwortet mit `400` auf den Fehler einer Parameter-Prüfung:
// ein `paramError` trägt seinen Code, eine verletzte Domänen-Invariante den
// Code von `apifault.Classify`, jeder andere Fehler den Rückfall der
// Ablehnungen.
func writeBadRequest(w http.ResponseWriter, err error) {
	var param *paramError
	if errors.As(err, &param) {
		writeError(w, http.StatusBadRequest, param.text, param.code)
		return
	}
	kind, code := apifault.Classify(err)
	if kind != apifault.InvalidInput {
		code = messagecode.RejectedFallback
	}
	writeError(w, http.StatusBadRequest, err.Error(), code)
}
