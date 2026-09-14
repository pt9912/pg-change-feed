package http

import (
	"encoding/json"
	"net/http"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// enableTableRequest trägt den JSON-Request-Body (`SPEC-018`): alle sieben
// Felder Pflicht, Übersetzung in `inbound.EnableTableCommand`
// (`LH-FA-CFG-001`).
type enableTableRequest struct {
	Source          string `json:"source"`
	Schema          string `json:"schema"`
	Table           string `json:"table"`
	TableID         string `json:"table_id"`
	SchemaVersionID string `json:"schema_version_id"`
	Version         int64  `json:"version"`
	Publication     string `json:"publication"`
}

// enableTableResponse trägt den JSON-Response-Body bei Erfolg (`SPEC-018`):
// `AlreadyEnabled` trägt `LH-FA-CFG-001`s Idempotenz-Ausgang fort — derselbe
// `201`-Status trägt beide Ausgänge, analog zu `registerConsumerResponse`.
type enableTableResponse struct {
	TableID        string `json:"table_id"`
	Source         string `json:"source"`
	Schema         string `json:"schema"`
	Table          string `json:"table"`
	AlreadyEnabled bool   `json:"already_enabled"`
}

// enableTableHandler übersetzt den JSON-Request in
// `inbound.EnableTableCommand` (`LH-FA-CFG-001`, `ADR-0057` Teilfrage 4):
// der Adapter importiert ausschließlich den Inbound Port und Domain-Typen
// zur Übersetzung, keine Application-Interna. Eine an der Quelle fehlende
// physische Tabelle trägt `inbound.ErrSourceTableMissing` (`404`,
// `writeDomainError`).
func enableTableHandler(useCase inbound.EnableTableUseCase, log outbound.LogPort) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req enableTableRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "Request-Body ist kein gültiges JSON")
			return
		}
		result, err := useCase.Enable(r.Context(), inbound.EnableTableCommand{
			Source:          model.SourceID(req.Source),
			Schema:          req.Schema,
			Table:           req.Table,
			TableID:         model.SourceTableID(req.TableID),
			SchemaVersionID: model.SchemaVersionID(req.SchemaVersionID),
			Version:         req.Version,
			Publication:     req.Publication,
		})
		if err != nil {
			writeDomainError(r.Context(), w, log, "EnableTable", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(enableTableResponse{
			TableID:        string(result.Table.ID),
			Source:         string(result.Table.SourceID),
			Schema:         result.Table.Schema,
			Table:          result.Table.Table,
			AlreadyEnabled: result.AlreadyEnabled,
		})
	})
}

// disableTableRequest trägt den JSON-Request-Body (`SPEC-018`): alle vier
// Felder Pflicht, Übersetzung in `inbound.DisableTableCommand`
// (`LH-FA-CFG-002`).
type disableTableRequest struct {
	Source      string `json:"source"`
	Schema      string `json:"schema"`
	Table       string `json:"table"`
	Publication string `json:"publication"`
}

// disableTableResponse trägt den JSON-Response-Body bei Erfolg
// (`SPEC-018`): `Removed`/`Retained` tragen den Bindungs-Zeilen-Ausgang
// (`LH-FA-CFG-002` Out-of-Scope: das Verhalten persistierter Changes folgt
// der Retention, nicht dieser Deaktivierung).
type disableTableResponse struct {
	Removed  bool `json:"removed"`
	Retained bool `json:"retained"`
}

func disableTableHandler(useCase inbound.DisableTableUseCase, log outbound.LogPort) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req disableTableRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "Request-Body ist kein gültiges JSON")
			return
		}
		result, err := useCase.Disable(r.Context(), inbound.DisableTableCommand{
			Source:      model.SourceID(req.Source),
			Schema:      req.Schema,
			Table:       req.Table,
			Publication: req.Publication,
		})
		if err != nil {
			writeDomainError(r.Context(), w, log, "DisableTable", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(disableTableResponse{
			Removed:  result.Removed,
			Retained: result.Retained,
		})
	})
}

// getStatusResponse trägt den JSON-Response-Body bei Erfolg (`SPEC-018`):
// `Enabled`/`Retained` trennen Erfassungs-Zustand und Herkunft
// (`LH-FA-CFG-003`); beide `false` liest eine nie aktivierte Tabelle.
type getStatusResponse struct {
	Enabled  bool `json:"enabled"`
	Retained bool `json:"retained"`
}

// getStatusHandler liest die vier Pflichtfelder aus Query-Parametern — ein
// `GET`-Request trägt keinen Body (`ADR-0057` Teilfrage 4, lesender
// Endpunkt). Eine an der Quelle fehlende physische Tabelle trägt
// `inbound.ErrSourceTableMissing` (`404`); eine nie registrierte Tabelle
// trägt `enabled=false, retained=false` (`200`, `LH-FA-CFG-003` Boundary).
func getStatusHandler(useCase inbound.GetStatusUseCase, log outbound.LogPort) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		source, schema, table, publication := q.Get("source"), q.Get("schema"), q.Get("table"), q.Get("publication")
		if source == "" || schema == "" || table == "" || publication == "" {
			writeError(w, http.StatusBadRequest, "source, schema, table und publication sind Pflichtfelder")
			return
		}
		result, err := useCase.Status(r.Context(), inbound.GetStatusQuery{
			Source:      model.SourceID(source),
			Schema:      schema,
			Table:       table,
			Publication: publication,
		})
		if err != nil {
			writeDomainError(r.Context(), w, log, "GetStatus", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(getStatusResponse{
			Enabled:  result.Enabled,
			Retained: result.Retained,
		})
	})
}

// sourceTableResponse trägt eine aktivierte oder als Herkunft erhaltene
// Tabelle (`SPEC-018`, `listTablesResponse`).
type sourceTableResponse struct {
	TableID string `json:"table_id"`
	Source  string `json:"source"`
	Schema  string `json:"schema"`
	Table   string `json:"table"`
}

// listTablesResponse trägt den JSON-Response-Body bei Erfolg (`SPEC-018`):
// `Tables`/`Retained` trennen Erfassungs-Zustand und Herkunft
// (`LH-FA-CFG-004`); ohne Aktivierung tragen beide leere Listen, nie `null`
// (`toSourceTableResponse` initialisiert mit `make(…, 0, …)`).
type listTablesResponse struct {
	Tables   []sourceTableResponse `json:"tables"`
	Retained []sourceTableResponse `json:"retained"`
}

func toSourceTableResponse(tables []model.SourceTable) []sourceTableResponse {
	out := make([]sourceTableResponse, 0, len(tables))
	for _, table := range tables {
		out = append(out, sourceTableResponse{
			TableID: string(table.ID),
			Source:  string(table.SourceID),
			Schema:  table.Schema,
			Table:   table.Table,
		})
	}
	return out
}

// listTablesHandler liest die zwei Pflichtfelder aus Query-Parametern — ein
// `GET`-Request trägt keinen Body (`ADR-0057` Teilfrage 4, lesender
// Endpunkt).
func listTablesHandler(useCase inbound.ListTablesUseCase, log outbound.LogPort) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		source, publication := q.Get("source"), q.Get("publication")
		if source == "" || publication == "" {
			writeError(w, http.StatusBadRequest, "source und publication sind Pflichtfelder")
			return
		}
		result, err := useCase.ListTables(r.Context(), inbound.ListTablesQuery{
			Source:      model.SourceID(source),
			Publication: publication,
		})
		if err != nil {
			writeDomainError(r.Context(), w, log, "ListTables", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(listTablesResponse{
			Tables:   toSourceTableResponse(result.Tables),
			Retained: toSourceTableResponse(result.Retained),
		})
	})
}
