package main

import "net/http"

// runRetentionRequest/-Response spiegeln `POST /retention/run`
// (`LH-FA-RET-002`, `internal/adapters/driving/http/retention.go`):
// `MinAgeNanos` 0 heißt „kein zeitliches Mindestalter" und ist gültig.
type runRetentionRequest struct {
	Source      string `json:"source"`
	MinAgeNanos int64  `json:"min_age_nanos"`
}

type runRetentionResponse struct {
	Deleted int `json:"deleted"`
}

// runRetention ruft den `admin`-Endpunkt `POST /retention/run` auf.
func runRetention(client *http.Client, cfg config) (runRetentionResponse, error) {
	var resp runRetentionResponse
	err := doRequestJSON(client, http.MethodPost, "http://"+cfg.addr+"/retention/run", cfg.adminToken,
		runRetentionRequest{Source: cfg.source, MinAgeNanos: cfg.minAgeNanos}, &resp, http.StatusOK)
	return resp, err
}
