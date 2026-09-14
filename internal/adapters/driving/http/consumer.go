package http

import (
	"encoding/json"
	"net/http"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// acknowledgeConsumerRequest trägt den JSON-Request-Body (`SPEC-018`): alle
// drei Felder Pflicht, Übersetzung in `inbound.AcknowledgeConsumerCommand`.
type acknowledgeConsumerRequest struct {
	ConsumerID string `json:"consumer_id"`
	SourceID   string `json:"source_id"`
	Offset     uint64 `json:"offset"`
}

// acknowledgeConsumerResponse trägt den JSON-Response-Body bei Erfolg
// (`SPEC-018`): die fortgeführte Position (`LH-FA-CON-004`).
type acknowledgeConsumerResponse struct {
	ConsumerID string `json:"consumer_id"`
	SourceID   string `json:"source_id"`
	Offset     uint64 `json:"offset"`
}

// acknowledgeConsumerHandler übersetzt den JSON-Request in
// `inbound.AcknowledgeConsumerCommand` (`LH-FA-CON-004`, `ADR-0057`
// Teilfrage 4): der Adapter importiert ausschließlich den Inbound Port und
// Domain-Typen zur Übersetzung, keine Application-Interna.
func acknowledgeConsumerHandler(useCase inbound.AcknowledgeConsumerUseCase, log outbound.LogPort) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req acknowledgeConsumerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "Request-Body ist kein gültiges JSON")
			return
		}
		result, err := useCase.Acknowledge(r.Context(), inbound.AcknowledgeConsumerCommand{
			Consumer: model.ConsumerID(req.ConsumerID),
			Position: model.SourcePosition{SourceID: model.SourceID(req.SourceID), Offset: req.Offset},
		})
		if err != nil {
			writeDomainError(r.Context(), w, log, "AcknowledgeConsumer", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(acknowledgeConsumerResponse{
			ConsumerID: string(result.Position.ConsumerID),
			SourceID:   string(result.Position.Position.SourceID),
			Offset:     result.Position.Position.Offset,
		})
	})
}

// getConsumerPositionResponse trägt den JSON-Response-Body bei Erfolg
// (`SPEC-018`): `Acknowledged` trägt `LH-FA-CON-005`s Boundary — der
// Nullwert liest die definierte Anfangsposition eines Consumers ohne
// Bestätigung, statt sie von einer echten Bestätigung mit Offset 0 zu
// unterscheiden (Letzteres ist keine gültige Position, `SPEC-003`).
type getConsumerPositionResponse struct {
	ConsumerID   string `json:"consumer_id"`
	SourceID     string `json:"source_id"`
	Offset       uint64 `json:"offset"`
	Acknowledged bool   `json:"acknowledged"`
}

// getConsumerPositionHandler liest die Kennung aus dem Query-Parameter
// `consumer_id` — ein `GET`-Request trägt keinen Body (`ADR-0057`
// Teilfrage 4, lesender Endpunkt).
func getConsumerPositionHandler(useCase inbound.GetConsumerPositionUseCase, log outbound.LogPort) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		consumerID := r.URL.Query().Get("consumer_id")
		if consumerID == "" {
			writeError(w, http.StatusBadRequest, "consumer_id ist Pflichtfeld")
			return
		}
		result, err := useCase.Position(r.Context(), inbound.GetConsumerPositionQuery{
			Consumer: model.ConsumerID(consumerID),
		})
		if err != nil {
			writeDomainError(r.Context(), w, log, "GetConsumerPosition", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(getConsumerPositionResponse{
			ConsumerID:   consumerID,
			SourceID:     string(result.Position.Position.SourceID),
			Offset:       result.Position.Position.Offset,
			Acknowledged: result.Position.Acknowledged(),
		})
	})
}

// removeConsumerRequest trägt den JSON-Request-Body (`SPEC-018`):
// `consumer_id` Pflicht.
type removeConsumerRequest struct {
	ConsumerID string `json:"consumer_id"`
}

// removeConsumerResponse trägt den JSON-Response-Body bei Erfolg
// (`SPEC-018`): `Removed` trägt `LH-FA-CON-006`s Idempotenz-Ausgang — ein
// nie registrierter Consumer meldet `false`, kein `404`: der Use Case
// behandelt die Entfernung als Idempotenz, nicht als Fehler gegen eine
// unbekannte Ressource, dieselbe fachliche Gleichwertigkeit über alle
// Zugriffswege wie CLI/SQL (`LH-FA-SST-006` Boundary).
type removeConsumerResponse struct {
	ConsumerID string `json:"consumer_id"`
	Removed    bool   `json:"removed"`
}

func removeConsumerHandler(useCase inbound.RemoveConsumerUseCase, log outbound.LogPort) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req removeConsumerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "Request-Body ist kein gültiges JSON")
			return
		}
		result, err := useCase.Remove(r.Context(), inbound.RemoveConsumerCommand{
			Consumer: model.ConsumerID(req.ConsumerID),
		})
		if err != nil {
			writeDomainError(r.Context(), w, log, "RemoveConsumer", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(removeConsumerResponse{
			ConsumerID: req.ConsumerID,
			Removed:    result.Removed,
		})
	})
}
