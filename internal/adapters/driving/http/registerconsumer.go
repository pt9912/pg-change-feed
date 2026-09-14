package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// registerConsumerRequest trägt den JSON-Request-Body (`SPEC-018`): beide
// Felder Pflicht, Übersetzung in `inbound.RegisterConsumerCommand`.
type registerConsumerRequest struct {
	ConsumerID string `json:"consumer_id"`
	Name       string `json:"name"`
}

// registerConsumerResponse trägt den JSON-Response-Body bei Erfolg
// (`SPEC-018`): `AlreadyRegistered` trägt `LH-FA-CON-001`s
// Idempotenz-Ausgang fort, ohne eine gesonderte Statuscode-Unterscheidung
// — derselbe `201`-Status trägt beide Ausgänge.
type registerConsumerResponse struct {
	ConsumerID        string `json:"consumer_id"`
	Name              string `json:"name"`
	AlreadyRegistered bool   `json:"already_registered"`
}

// registerConsumerHandler übersetzt den JSON-Request in
// `inbound.RegisterConsumerCommand` (`LH-FA-CON-001`, `ADR-0057`
// Teilfrage 4): der Adapter importiert ausschließlich den Inbound Port
// und Domain-Typen zur Übersetzung, keine Application-Interna.
func registerConsumerHandler(useCase inbound.RegisterConsumerUseCase, log outbound.LogPort) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req registerConsumerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "Request-Body ist kein gültiges JSON")
			return
		}
		result, err := useCase.Register(r.Context(), inbound.RegisterConsumerCommand{
			Consumer: model.ConsumerID(req.ConsumerID),
			Name:     req.Name,
		})
		if err != nil {
			if errors.Is(err, domainerrors.ErrEmptyIdentifier) {
				writeError(w, http.StatusBadRequest, "consumer_id und name sind Pflichtfelder")
				return
			}
			log.Warn(r.Context(), "http: RegisterConsumer fehlgeschlagen", "error", err)
			writeError(w, http.StatusInternalServerError, "interner Fehler")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(registerConsumerResponse{
			ConsumerID:        string(result.Consumer.ID),
			Name:              result.Consumer.Name,
			AlreadyRegistered: result.AlreadyRegistered,
		})
	})
}
