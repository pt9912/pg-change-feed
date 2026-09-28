package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestConsumerPositionURLCarriesConsumerID prüft den Aufbau der
// Lese-Adresse von `GET /consumers/position`: ihr einziges Pflichtfeld
// steht als Query-Parameter.
func TestConsumerPositionURLCarriesConsumerID(t *testing.T) {
	got := ConsumerPositionURL("feed:8080", "consumer-1")
	want := "http://feed:8080/consumers/position?consumer_id=consumer-1"
	if got != want {
		t.Fatalf("ConsumerPositionURL = %q, want %q", got, want)
	}
}

// TestConsumerPositionURLEscapesReservedCharacters prüft, dass die
// Consumer-Kennung als Query-Wert kodiert wird.
func TestConsumerPositionURLEscapesReservedCharacters(t *testing.T) {
	got := ConsumerPositionURL("feed:8080", "consumer & 1")
	want := "http://feed:8080/consumers/position?consumer_id=consumer+%26+1"
	if got != want {
		t.Fatalf("ConsumerPositionURL = %q, want %q", got, want)
	}
}

// TestRegisterConsumerParsesSuccessResponse prüft die Erfolgs-Antwort von
// `POST /consumers` gegen einen netzlosen In-Prozess-Server (loopback,
// dieselbe Grenze wie beim `natsnotify`-Adapter, `harness/README.md`
// §Sensors, `make test-notify`).
func TestRegisterConsumerParsesSuccessResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/consumers" {
			t.Fatalf("unerwarteter Aufruf: %s %s", r.Method, r.URL.Path)
		}
		var req registerConsumerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Request-Body: %v", err)
		}
		if req.ConsumerID != "consumer-1" || req.Name != "Consumer Eins" {
			t.Fatalf("unerwarteter Request-Body: %+v", req)
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(registerConsumerResponse{ConsumerID: "consumer-1", Name: "Consumer Eins", AlreadyRegistered: false})
	}))
	defer server.Close()

	cfg := config{addr: strings.TrimPrefix(server.URL, "http://"), adminToken: "admin-token", consumerID: "consumer-1", name: "Consumer Eins"}
	resp, err := registerConsumer(server.Client(), cfg)
	if err != nil {
		t.Fatalf("registerConsumer: %v", err)
	}
	if resp.ConsumerID != "consumer-1" || resp.Name != "Consumer Eins" || resp.AlreadyRegistered {
		t.Fatalf("unerwartete Antwort: %+v", resp)
	}
}

// TestRegisterConsumerFailsOnNon2xx prüft den Fehlerpfad: ein Nicht-201-
// Status wird als sichtbarer Fehler mit Statuscode und Antworttext
// weitergegeben, kein stiller Leerwert.
func TestRegisterConsumerFailsOnNon2xx(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "consumer_id und name sind Pflichtfelder"})
	}))
	defer server.Close()

	cfg := config{addr: strings.TrimPrefix(server.URL, "http://"), adminToken: "admin-token"}
	_, err := registerConsumer(server.Client(), cfg)
	if err == nil {
		t.Fatal("registerConsumer: erwarteter Fehler blieb aus")
	}
	if !strings.Contains(err.Error(), "400") || !strings.Contains(err.Error(), "Pflichtfelder") {
		t.Fatalf("Fehlertext trägt weder Status noch Antworttext: %v", err)
	}
}

// TestAcknowledgeConsumerParsesSuccessResponse prüft die Erfolgs-Antwort von
// `POST /consumers/acknowledge`.
func TestAcknowledgeConsumerParsesSuccessResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req acknowledgeConsumerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Request-Body: %v", err)
		}
		if req.Offset != 42 {
			t.Fatalf("unerwarteter Offset: %d", req.Offset)
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(acknowledgeConsumerResponse{ConsumerID: req.ConsumerID, SourceID: req.SourceID, Offset: req.Offset})
	}))
	defer server.Close()

	cfg := config{addr: strings.TrimPrefix(server.URL, "http://"), adminToken: "admin-token", consumerID: "consumer-1", source: "quelle-1", offset: 42}
	resp, err := acknowledgeConsumer(server.Client(), cfg)
	if err != nil {
		t.Fatalf("acknowledgeConsumer: %v", err)
	}
	if resp.Offset != 42 {
		t.Fatalf("unerwartete Antwort: %+v", resp)
	}
}

// TestConsumerPositionParsesSuccessResponse prüft die Erfolgs-Antwort von
// `GET /consumers/position`, inklusive der Boundary `acknowledged=false`.
func TestConsumerPositionParsesSuccessResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("consumer_id"); got != "consumer-1" {
			t.Fatalf("unerwartete consumer_id: %q", got)
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(consumerPositionResponse{ConsumerID: "consumer-1", Acknowledged: false})
	}))
	defer server.Close()

	cfg := config{addr: strings.TrimPrefix(server.URL, "http://"), token: "reader-token", consumerID: "consumer-1"}
	resp, err := consumerPosition(server.Client(), cfg)
	if err != nil {
		t.Fatalf("consumerPosition: %v", err)
	}
	if resp.Acknowledged {
		t.Fatalf("unerwartete Antwort: %+v", resp)
	}
}

// TestRemoveConsumerParsesIdempotentOutcome prüft `LH-FA-CON-006`s
// Idempotenz-Ausgang: ein nie registrierter Consumer meldet `removed=false`
// über `200`, keinen `404`.
func TestRemoveConsumerParsesIdempotentOutcome(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(removeConsumerResponse{ConsumerID: "nie-registriert", Removed: false})
	}))
	defer server.Close()

	cfg := config{addr: strings.TrimPrefix(server.URL, "http://"), adminToken: "admin-token", consumerID: "nie-registriert"}
	resp, err := removeConsumer(server.Client(), cfg)
	if err != nil {
		t.Fatalf("removeConsumer: %v", err)
	}
	if resp.Removed {
		t.Fatalf("unerwartete Antwort: %+v", resp)
	}
}
