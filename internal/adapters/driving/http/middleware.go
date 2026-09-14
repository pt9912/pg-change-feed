package http

import (
	"encoding/json"
	"net/http"
	"strings"
)

// role trägt die zwei Rechtsklassen der Token-Middleware (`ADR-0057`
// Teilfrage 3): `roleAdmin` deckt implizit `roleReader` ab — dieselbe
// Hierarchie wie zwischen `cdc_admin` und `cdc_reader` (`ADR-0047`).
// `roleNone` trägt sowohl den fehlenden als auch den unbekannten Token —
// beide enden über denselben `401`-Pfad (`withToken`).
type role int

const (
	roleNone role = iota
	roleReader
	roleAdmin
)

// classifyToken ordnet einen Bearer-Token einer Rechtsklasse zu. Ein
// leer konfiguriertes Token (`readerToken`/`adminToken` ungesetzt) trifft
// nie ein Aufruf-Token — sonst würde ein fehlender Header (leerer Token)
// gegen eine ebenfalls ungesetzte Token-Klasse eine dritte, implizite
// Rechtsklasse eröffnen.
//
// Diese Funktion steht als zweite, wortgleiche Fassung in
// `internal/adapters/driving/grpc/interceptor.go` (`role`, die drei
// Konstanten und `classifyToken`). Das `.a-check.yml`-Schichtenmodell führt
// keine `adapters→adapters`-Kante, deshalb trägt jeder Driving-Adapter seine
// eigene Fassung derselben Zuordnung; beide Fassungen sind zusammen zu
// ändern.
func classifyToken(token, readerToken, adminToken string) role {
	if token == "" {
		return roleNone
	}
	if adminToken != "" && token == adminToken {
		return roleAdmin
	}
	if readerToken != "" && token == readerToken {
		return roleReader
	}
	return roleNone
}

// bearerToken liest den Token aus dem `Authorization`-Header; ein
// fehlendes oder falsch geformtes Schema trägt einen leeren Token zurück
// — `classifyToken` behandelt ihn wie einen fehlenden Token.
func bearerToken(r *http.Request) string {
	const prefix = "Bearer "
	header := r.Header.Get("Authorization")
	if !strings.HasPrefix(header, prefix) {
		return ""
	}
	return strings.TrimPrefix(header, prefix)
}

// errorResponse trägt den JSON-Fehler-Body (`SPEC-018`) für
// `400`/`401`/`403`/`500`.
type errorResponse struct {
	Error string `json:"error"`
}

// writeError schreibt einen JSON-Fehler-Body mit dem übergebenen
// Statuscode.
func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(errorResponse{Error: message})
}

// withToken schützt einen Handler mit der Token-Middleware (`ADR-0057`
// Teilfrage 3, Fitness Function): ein Aufruf ohne oder mit unbekanntem
// Bearer-Token endet mit `401`; ein bekanntes Token unterhalb der
// geforderten Rechtsklasse endet mit `403` — beides sichtbar, nicht still
// verworfen (`LH-FA-SST-006` Negative-Kriterium).
func withToken(readerToken, adminToken string, required role, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callerRole := classifyToken(bearerToken(r), readerToken, adminToken)
		if callerRole == roleNone {
			writeError(w, http.StatusUnauthorized, "fehlender oder unbekannter Bearer-Token")
			return
		}
		if callerRole < required {
			writeError(w, http.StatusForbidden, "Rechtsklasse unzureichend für diesen Endpunkt")
			return
		}
		next.ServeHTTP(w, r)
	})
}
