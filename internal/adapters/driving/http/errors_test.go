package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

// TestWriteDomainErrorBildetEinheitlichAb trägt das einheitliche
// Fehler-Mapping direkt an der Funktion, die alle acht Handler dieses
// Adapters gemeinsam nutzen: jede benannte Domänen-Invariante (`ADR-0029`)
// auf `400`, die fehlende physische Tabelle (`inbound.ErrSourceTableMissing`)
// auf `404`, jeder übrige Fehler auf `500`.
func TestWriteDomainErrorBildetEinheitlichAb(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{"ErrEmptyIdentifier", domainerrors.ErrEmptyIdentifier, http.StatusBadRequest},
		{"ErrInvalidPosition", domainerrors.ErrInvalidPosition, http.StatusBadRequest},
		{"ErrSourceMismatch", domainerrors.ErrSourceMismatch, http.StatusBadRequest},
		{"ErrPositionRegression", domainerrors.ErrPositionRegression, http.StatusBadRequest},
		{"ErrNegativeDuration", domainerrors.ErrNegativeDuration, http.StatusBadRequest},
		{"ErrNonPositiveVersion", domainerrors.ErrNonPositiveVersion, http.StatusBadRequest},
		{"ErrSourceTableMissing (gewrappt)", fmt.Errorf("%w: public.orders", inbound.ErrSourceTableMissing), http.StatusNotFound},
		{"unbekannter Fehler", errors.New("speicher nicht erreichbar"), http.StatusInternalServerError},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			writeDomainError(context.Background(), rec, outbound.NoopLog, "Test", tc.err)
			if rec.Code != tc.want {
				t.Fatalf("Status: %d (Erwartung: %d)", rec.Code, tc.want)
			}
			var decoded errorResponse
			if err := json.NewDecoder(rec.Body).Decode(&decoded); err != nil {
				t.Fatalf("Antwort dekodieren: %v", err)
			}
			if decoded.Error == "" {
				t.Fatalf("Antwort trägt keine Fehlermeldung: %+v", decoded)
			}
		})
	}
}
