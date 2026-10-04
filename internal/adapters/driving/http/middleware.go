package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/pt9912/pg-change-feed/internal/application/port/apiauth"
	"github.com/pt9912/pg-change-feed/internal/domain/messagecode"
)

// bearerToken liest den Token aus dem `Authorization`-Header; ein
// fehlendes oder falsch geformtes Schema trägt einen leeren Token zurück
// — `apiauth.Classifier.Classify` behandelt ihn wie einen fehlenden Token.
func bearerToken(r *http.Request) string {
	const prefix = "Bearer "
	header := r.Header.Get("Authorization")
	if !strings.HasPrefix(header, prefix) {
		return ""
	}
	return strings.TrimPrefix(header, prefix)
}

// errorResponse trägt den JSON-Fehler-Body (`SPEC-018`) für
// `400`/`401`/`403`/`404`/`500`/`503`: den Klartext und, wo eine
// Code-Zuordnung besteht, den Meldungscode. Ein leerer Code erscheint nicht
// im Body (`401`, `403`).
type errorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code,omitempty"`
}

// writeError schreibt einen JSON-Fehler-Body mit dem übergebenen
// Statuscode; `code` ist der Meldungscode der Ursache, leer ohne Zuordnung.
func writeError(w http.ResponseWriter, status int, message string, code messagecode.Code) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(errorResponse{Error: message, Code: string(code)})
}

// withToken schützt einen Handler mit der Token-Middleware (`ADR-0057`
// Teilfrage 3, Fitness Function): ein Aufruf ohne oder mit unbekanntem
// Bearer-Token endet mit `401`; ein bekanntes Token unterhalb der
// geforderten Rechtsklasse endet mit `403` — beides sichtbar, nicht still
// verworfen.
func withToken(tokens *apiauth.Classifier, required apiauth.Role, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callerRole := tokens.Classify(bearerToken(r))
		if callerRole == apiauth.None {
			writeError(w, http.StatusUnauthorized, "fehlender oder unbekannter Bearer-Token", "")
			return
		}
		if callerRole < required {
			writeError(w, http.StatusForbidden, "Rechtsklasse unzureichend für diesen Endpunkt", "")
			return
		}
		next.ServeHTTP(w, r)
	})
}
