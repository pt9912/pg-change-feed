package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestDoRequestJSONDecodesSuccessBody prüft den Erfolgspfad: der Body wird
// kodiert, das Bearer-Token gesetzt, und die Antwort bei `wantStatus` nach
// `out` dekodiert.
func TestDoRequestJSONDecodesSuccessBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Fatalf("unerwarteter Authorization-Header: %q", got)
		}
		var body map[string]string
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("Request-Body: %v", err)
		}
		if body["field"] != "wert" {
			t.Fatalf("unerwarteter Request-Body: %+v", body)
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"echo": body["field"]})
	}))
	defer server.Close()

	var out map[string]string
	err := doRequestJSON(server.Client(), http.MethodPost, server.URL, "test-token", map[string]string{"field": "wert"}, &out, http.StatusOK)
	if err != nil {
		t.Fatalf("doRequestJSON: %v", err)
	}
	if out["echo"] != "wert" {
		t.Fatalf("unerwartete Antwort: %+v", out)
	}
}

// TestDoRequestJSONFailsOnUnexpectedStatus prüft den Fehlerpfad: ein von
// `wantStatus` abweichender Statuscode liefert Statuscode und Antworttext
// im Fehler, statt still zu dekodieren.
func TestDoRequestJSONFailsOnUnexpectedStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "kein Zugriff"})
	}))
	defer server.Close()

	var out map[string]string
	err := doRequestJSON(server.Client(), http.MethodGet, server.URL, "test-token", nil, &out, http.StatusOK)
	if err == nil {
		t.Fatal("doRequestJSON: erwarteter Fehler blieb aus")
	}
	if !strings.Contains(err.Error(), "403") || !strings.Contains(err.Error(), "kein Zugriff") {
		t.Fatalf("Fehlertext trägt weder Status noch Antworttext: %v", err)
	}
}

// TestDoRequestJSONFailsOnUnreachableHost prüft den Fehlerpfad ohne
// erreichbaren Host — netzlos: die Verbindung zu einem Loopback-Port ohne
// Listener scheitert sofort, ohne echtes Netz zu brauchen (dieselbe Grenze
// wie beim `natsnotify`-Adapter, `harness/README.md` §Sensors,
// `make test-notify`).
func TestDoRequestJSONFailsOnUnreachableHost(t *testing.T) {
	client := &http.Client{Timeout: requestTimeout}
	var out map[string]string
	err := doRequestJSON(client, http.MethodGet, "http://127.0.0.1:1/nichts", "test-token", nil, &out, http.StatusOK)
	if err == nil {
		t.Fatal("doRequestJSON: erwarteter Fehler blieb aus")
	}
}

// TestDoRequestJSONSkipsDecodeWithoutOut prüft, dass `out == nil` keinen
// Dekodier-Versuch auslöst — nicht jede Fähigkeit braucht eine typisierte
// Antwort.
func TestDoRequestJSONSkipsDecodeWithoutOut(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("kein gültiges JSON"))
	}))
	defer server.Close()

	err := doRequestJSON(server.Client(), http.MethodGet, server.URL, "test-token", nil, nil, http.StatusOK)
	if err != nil {
		t.Fatalf("doRequestJSON: %v", err)
	}
}
