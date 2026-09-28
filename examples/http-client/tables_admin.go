package main

import (
	"net/http"
	"net/url"
)

// enableTableRequest/-Response spiegeln `POST /tables/enable`
// (`LH-FA-CFG-001`, `internal/adapters/driving/http/verwaltung.go`): alle
// sieben Request-Felder Pflicht, `AlreadyEnabled` trägt die
// Idempotenz-Antwort ohne eigenen Statuscode.
type enableTableRequest struct {
	Source          string `json:"source"`
	Schema          string `json:"schema"`
	Table           string `json:"table"`
	TableID         string `json:"table_id"`
	SchemaVersionID string `json:"schema_version_id"`
	Version         int64  `json:"version"`
	Publication     string `json:"publication"`
}

type enableTableResponse struct {
	TableID        string `json:"table_id"`
	Source         string `json:"source"`
	Schema         string `json:"schema"`
	Table          string `json:"table"`
	AlreadyEnabled bool   `json:"already_enabled"`
}

// enableTable ruft den `admin`-Endpunkt `POST /tables/enable` auf.
func enableTable(client *http.Client, cfg config) (enableTableResponse, error) {
	var resp enableTableResponse
	err := doRequestJSON(client, http.MethodPost, "http://"+cfg.addr+"/tables/enable", cfg.adminToken,
		enableTableRequest{
			Source:          cfg.source,
			Schema:          cfg.schema,
			Table:           cfg.table,
			TableID:         cfg.tableID,
			SchemaVersionID: cfg.schemaVersionID,
			Version:         cfg.version,
			Publication:     cfg.publication,
		}, &resp, http.StatusCreated)
	return resp, err
}

// disableTableRequest/-Response spiegeln `POST /tables/disable`
// (`LH-FA-CFG-002`).
type disableTableRequest struct {
	Source      string `json:"source"`
	Schema      string `json:"schema"`
	Table       string `json:"table"`
	Publication string `json:"publication"`
}

type disableTableResponse struct {
	Removed  bool `json:"removed"`
	Retained bool `json:"retained"`
}

// disableTable ruft den `admin`-Endpunkt `POST /tables/disable` auf.
func disableTable(client *http.Client, cfg config) (disableTableResponse, error) {
	var resp disableTableResponse
	err := doRequestJSON(client, http.MethodPost, "http://"+cfg.addr+"/tables/disable", cfg.adminToken,
		disableTableRequest{Source: cfg.source, Schema: cfg.schema, Table: cfg.table, Publication: cfg.publication}, &resp, http.StatusOK)
	return resp, err
}

// tableStatusResponse spiegelt `GET /tables/status` (`LH-FA-CFG-003`):
// `Enabled`/`Retained` trennen Erfassungs-Zustand und Herkunft.
type tableStatusResponse struct {
	Enabled  bool `json:"enabled"`
	Retained bool `json:"retained"`
}

// TableStatusURL baut die Lese-Adresse von `GET /tables/status`: alle vier
// Parameter sind Pflicht.
func TableStatusURL(addr, source, schema, table, publication string) string {
	u := url.URL{Scheme: "http", Host: addr, Path: "/tables/status"}
	q := u.Query()
	q.Set("source", source)
	q.Set("schema", schema)
	q.Set("table", table)
	q.Set("publication", publication)
	u.RawQuery = q.Encode()
	return u.String()
}

// tableStatus ruft den `reader`-Endpunkt `GET /tables/status` auf.
func tableStatus(client *http.Client, cfg config) (tableStatusResponse, error) {
	var resp tableStatusResponse
	err := doRequestJSON(client, http.MethodGet, TableStatusURL(cfg.addr, cfg.source, cfg.schema, cfg.table, cfg.publication), cfg.token, nil, &resp, http.StatusOK)
	return resp, err
}
