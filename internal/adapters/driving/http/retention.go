package http

import (
	"encoding/json"
	"net/http"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// runRetentionRequest trägt den JSON-Request-Body (`SPEC-018`): `source`
// Pflicht, `min_age_nanos` das Mindestalter der Policy in Nanosekunden
// (`LH-FA-RET-003`) — 0 heißt „kein zeitliches Mindestalter" und ist gültig,
// negativ trägt `model.ErrNegativeDuration` (`400`).
type runRetentionRequest struct {
	Source      string `json:"source"`
	MinAgeNanos int64  `json:"min_age_nanos"`
}

// runRetentionResponse trägt den JSON-Response-Body bei Erfolg
// (`SPEC-018`): die Anzahl real gelöschter Changes.
type runRetentionResponse struct {
	Deleted int `json:"deleted"`
}

// runRetentionHandler übersetzt den JSON-Request in
// `inbound.RunRetentionCommand` (`LH-FA-RET-002`…`004`, `ADR-0057`
// Teilfrage 4): der Adapter importiert ausschließlich den Inbound Port und
// Domain-Typen zur Übersetzung, keine Application-Interna. Die Policy
// entsteht über den Domänen-Konstruktor — eine negative Dauer endet über die
// Invariante, bevor der Use Case aufgerufen wird.
func runRetentionHandler(useCase inbound.RunRetentionUseCase, log outbound.LogPort) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req runRetentionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "Request-Body ist kein gültiges JSON")
			return
		}
		minAge, err := model.NewDuration(req.MinAgeNanos)
		if err != nil {
			writeDomainError(r.Context(), w, log, "RunRetention", err)
			return
		}
		policy, err := model.NewRetentionPolicy(minAge)
		if err != nil {
			writeDomainError(r.Context(), w, log, "RunRetention", err)
			return
		}
		result, err := useCase.Run(r.Context(), inbound.RunRetentionCommand{
			Source: model.SourceID(req.Source),
			Policy: policy,
		})
		if err != nil {
			writeDomainError(r.Context(), w, log, "RunRetention", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(runRetentionResponse{Deleted: result.Deleted})
	})
}
