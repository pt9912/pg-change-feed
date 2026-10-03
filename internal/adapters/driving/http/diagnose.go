package http

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/domain/messagecode"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// diagnoseParamSource trägt den einzigen Query-Parameter des lesenden
// Endpunkts `GET /diagnose` (`ADR-0132`).
const diagnoseParamSource = "source"

// diagnoseParams trägt die geschlossene Parameter-Menge dieses Endpunkts —
// dieselbe Verschärfung wie bei `GET /changes` (`ADR-0132` Teilfrage 4): ein
// Parameter außerhalb der Menge endet mit `400`.
var diagnoseParams = map[string]bool{diagnoseParamSource: true}

// diagnoseConsumerLagResponse, diagnoseRetentionBlockerResponse und
// diagnoseBackfillTableResponse spiegeln `inbound.ConsumerLag`/
// `RetentionBlocker`/`BackfillTableStatus` 1:1 in JSON-Form (`ADR-0132`
// Teilfrage 4); ein abwesender Wert steht als `null`, dieselbe Konvention
// wie `readChangeResponse.OldImage`/`NewImage`.
type diagnoseConsumerLagResponse struct {
	ConsumerID string   `json:"consumer_id"`
	Lag        *float64 `json:"lag"`
}

type diagnoseRetentionBlockerResponse struct {
	ConsumerID           string `json:"consumer_id"`
	Name                 string `json:"name"`
	AcknowledgedPosition int64  `json:"acknowledged_position"`
	Backlog              *int64 `json:"backlog"`
}

type diagnoseBackfillTableResponse struct {
	Schema            string `json:"schema"`
	Table             string `json:"table"`
	Status            string `json:"status"`
	RowsCopied        int64  `json:"rows_copied"`
	EstimatedRows     *int64 `json:"estimated_rows"`
	WarnEstimatedSize bool   `json:"warn_estimated_size"`
	WarnDuration      bool   `json:"warn_duration"`
	ErrorMessage      string `json:"error_message"`
}

// diagnoseResponse trägt den JSON-Response-Body von `GET /diagnose`
// (`ADR-0132`): dieselben sechs Signalgruppen wie der bestehende
// CLI-Text-Bericht, strukturiert statt gedruckt. Leere Mengen
// (`ConsumerLags`, `Backfill`) tragen eine leere, gesetzte Liste, nie
// `null` — dieselbe Zusage wie `readChangesResponse.Changes`.
type diagnoseResponse struct {
	HeartbeatAgeSeconds *float64                          `json:"heartbeat_age_seconds"`
	ErrorClass          *string                           `json:"error_class"`
	ErrorCode           *string                           `json:"error_code"`
	CaptureLag          float64                           `json:"capture_lag"`
	ConsumerLags        []diagnoseConsumerLagResponse     `json:"consumer_lags"`
	RetentionBlocker    *diagnoseRetentionBlockerResponse `json:"retention_blocker"`
	StorageBytes        float64                           `json:"storage_bytes"`
	Backfill            []diagnoseBackfillTableResponse   `json:"backfill"`
}

// toDiagnoseResponse übersetzt das Use-Case-Ergebnis in seine Antwortform.
func toDiagnoseResponse(result inbound.DiagnoseResult) diagnoseResponse {
	consumerLags := make([]diagnoseConsumerLagResponse, 0, len(result.ConsumerLags))
	for _, lag := range result.ConsumerLags {
		consumerLags = append(consumerLags, diagnoseConsumerLagResponse{ConsumerID: lag.ConsumerID, Lag: lag.Lag})
	}
	var blocker *diagnoseRetentionBlockerResponse
	if result.RetentionBlocker != nil {
		blocker = &diagnoseRetentionBlockerResponse{
			ConsumerID:           result.RetentionBlocker.ConsumerID,
			Name:                 result.RetentionBlocker.Name,
			AcknowledgedPosition: result.RetentionBlocker.AcknowledgedPosition,
			Backlog:              result.RetentionBlocker.Backlog,
		}
	}
	backfill := make([]diagnoseBackfillTableResponse, 0, len(result.Backfill))
	for _, table := range result.Backfill {
		backfill = append(backfill, diagnoseBackfillTableResponse{
			Schema:            table.Schema,
			Table:             table.Table,
			Status:            table.Status,
			RowsCopied:        table.RowsCopied,
			EstimatedRows:     table.EstimatedRows,
			WarnEstimatedSize: table.WarnEstimatedSize,
			WarnDuration:      table.WarnDuration,
			ErrorMessage:      table.ErrorMessage,
		})
	}
	return diagnoseResponse{
		HeartbeatAgeSeconds: result.HeartbeatAgeSeconds,
		ErrorClass:          result.ErrorClass,
		ErrorCode:           result.ErrorCode,
		CaptureLag:          result.CaptureLag,
		ConsumerLags:        consumerLags,
		RetentionBlocker:    blocker,
		StorageBytes:        result.StorageBytes,
		Backfill:            backfill,
	}
}

// diagnoseHandler trägt den lesenden Endpunkt `GET /diagnose` (`ADR-0132`):
// er übersetzt den Query-Parameter in `inbound.DiagnoseQuery` und ruft den
// Inbound Port — derselbe Use Case, den auch der CLI-Sondermodus und der
// gRPC-Handler aufrufen, kein zweiter Domänenpfad. Die Authn trägt die
// vorgelagerte `withToken`-Middleware (`reader`-Rechtsklasse).
func diagnoseHandler(useCase inbound.DiagnoseUseCase, log outbound.LogPort) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		source, err := parseDiagnoseQuery(r.URL.Query())
		if err != nil {
			writeBadRequest(w, err)
			return
		}
		result, err := useCase.Diagnose(r.Context(), inbound.DiagnoseQuery{Source: source})
		if err != nil {
			writeDomainError(r.Context(), w, log, "Diagnose", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(toDiagnoseResponse(result))
	})
}

// parseDiagnoseQuery liest die geschlossene Parameter-Menge (`source`
// Pflicht) — ein Parameter außerhalb der Menge endet sichtbar, kein
// unbemerkter Tippfehler.
func parseDiagnoseQuery(values url.Values) (model.SourceID, error) {
	for name := range values {
		if !diagnoseParams[name] {
			return "", rejectParam(messagecode.RejectedParameterUnknown, "unbekannter Query-Parameter %q", name)
		}
	}
	source := values.Get(diagnoseParamSource)
	if source == "" {
		return "", rejectParam(messagecode.RejectedRequiredField, "%s ist Pflichtfeld", diagnoseParamSource)
	}
	return model.SourceID(source), nil
}
