package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestRunRetentionParsesZeroDeletedAsValidOutcome prüft, dass `deleted=0`
// (nichts abzuräumen) ein gültiges, unterscheidbares Ergebnis ist — kein
// Fehler und kein Leerwert.
func TestRunRetentionParsesZeroDeletedAsValidOutcome(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req runRetentionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Request-Body: %v", err)
		}
		if req.Source != "quelle-1" || req.MinAgeNanos != 3600000000000 {
			t.Fatalf("unerwarteter Request-Body: %+v", req)
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(runRetentionResponse{Deleted: 0})
	}))
	defer server.Close()

	cfg := config{addr: strings.TrimPrefix(server.URL, "http://"), adminToken: "admin-token", source: "quelle-1", minAgeNanos: 3600000000000}
	resp, err := runRetention(server.Client(), cfg)
	if err != nil {
		t.Fatalf("runRetention: %v", err)
	}
	if resp.Deleted != 0 {
		t.Fatalf("unerwartete Antwort: %+v", resp)
	}
}

// TestRunRetentionFailsOnNon2xx prüft den Fehlerpfad — eine negative Dauer
// endet serverseitig mit `400`.
func TestRunRetentionFailsOnNon2xx(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "min_age_nanos darf nicht negativ sein"})
	}))
	defer server.Close()

	cfg := config{addr: strings.TrimPrefix(server.URL, "http://"), adminToken: "admin-token", source: "quelle-1", minAgeNanos: -1}
	_, err := runRetention(server.Client(), cfg)
	if err == nil {
		t.Fatal("runRetention: erwarteter Fehler blieb aus")
	}
	if !strings.Contains(err.Error(), "400") {
		t.Fatalf("Fehlertext trägt keinen Statuscode: %v", err)
	}
}
