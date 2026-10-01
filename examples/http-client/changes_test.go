package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestChangesURLCarriesOnlyRequiredField prüft, dass ohne die sechs
// optionalen Parameter nur `source` in der Adresse steht (kein leeres
// `target=`).
func TestChangesURLCarriesOnlyRequiredField(t *testing.T) {
	got := ChangesURL("feed:8080", "quelle-1", "", "", "", "", "", "")
	want := "http://feed:8080/changes?source=quelle-1"
	if got != want {
		t.Fatalf("ChangesURL = %q, want %q", got, want)
	}
}

// TestChangesURLCarriesAllOptionalFields prüft, dass alle sechs optionalen
// Parameter unabhängig voneinander übernommen werden.
func TestChangesURLCarriesAllOptionalFields(t *testing.T) {
	got := ChangesURL("feed:8080", "quelle-1", "public", "orders", "eu", "10", "20", "5")
	want := "http://feed:8080/changes?from=10&limit=5&schema=public&source=quelle-1&table=orders&target=eu&to=20"
	if got != want {
		t.Fatalf("ChangesURL = %q, want %q", got, want)
	}
}

// TestChangesURLCarriesTargetAlone prüft, dass `target` ohne `schema`/`table`
// als eigener Query-Parameter erscheint und ein Wert mit Sonderzeichen
// maskiert wird.
func TestChangesURLCarriesTargetAlone(t *testing.T) {
	got := ChangesURL("feed:8080", "quelle-1", "", "", "a&b=c", "", "", "")
	want := "http://feed:8080/changes?source=quelle-1&target=a%26b%3Dc"
	if got != want {
		t.Fatalf("ChangesURL = %q, want %q", got, want)
	}
}

// TestReadChangesParsesSuccessResponse prüft die Erfolgs-Antwort von
// `GET /changes`, inklusive `old_image`/`new_image` als eingebettete
// JSON-Werte und der Origin-Kennzeichnung.
func TestReadChangesParsesSuccessResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("source"); got != "quelle-1" {
			t.Fatalf("unerwartete source: %q", got)
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(readChangesResponse{Changes: []readChangeResponse{{
			CommitPosition: 42,
			ChangeID:       "change-1",
			Schema:         "public",
			Table:          "orders",
			Operation:      "INSERT",
			OldImage:       nil,
			NewImage:       json.RawMessage(`{"name":"erste Zeile"}`),
			Origin:         "wal",
		}}})
	}))
	defer server.Close()

	cfg := config{addr: strings.TrimPrefix(server.URL, "http://"), token: "reader-token", source: "quelle-1"}
	resp, err := readChanges(server.Client(), cfg)
	if err != nil {
		t.Fatalf("readChanges: %v", err)
	}
	if len(resp.Changes) != 1 || resp.Changes[0].ChangeID != "change-1" || resp.Changes[0].Origin != "wal" {
		t.Fatalf("unerwartete Antwort: %+v", resp)
	}
	if string(resp.Changes[0].OldImage) != "null" {
		t.Fatalf("unerwartetes old_image: %s", resp.Changes[0].OldImage)
	}
}

// TestReadChangesSendsTargetFromConfig prüft die Bindung Flag → Anfrage: das
// Feld `target` der Konfiguration erscheint als Query-Parameter, ohne Wert
// bleibt der Parameter weg.
func TestReadChangesSendsTargetFromConfig(t *testing.T) {
	var gotTarget []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotTarget = r.URL.Query()["target"]
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(readChangesResponse{Changes: []readChangeResponse{}})
	}))
	defer server.Close()
	addr := strings.TrimPrefix(server.URL, "http://")

	if _, err := readChanges(server.Client(), config{addr: addr, token: "reader-token", source: "quelle-1", target: "eu"}); err != nil {
		t.Fatalf("readChanges: %v", err)
	}
	if len(gotTarget) != 1 || gotTarget[0] != "eu" {
		t.Fatalf("target = %v, want [eu]", gotTarget)
	}

	if _, err := readChanges(server.Client(), config{addr: addr, token: "reader-token", source: "quelle-1"}); err != nil {
		t.Fatalf("readChanges: %v", err)
	}
	if len(gotTarget) != 0 {
		t.Fatalf("target = %v, want absent", gotTarget)
	}
}

// TestReadChangesFailsOnNon2xx prüft den Fehlerpfad — ein unbekannter
// Query-Parameter endet serverseitig mit `400` (`ADR-0081` Teilfrage 4).
func TestReadChangesFailsOnNon2xx(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "source ist Pflichtfeld"})
	}))
	defer server.Close()

	cfg := config{addr: strings.TrimPrefix(server.URL, "http://"), token: "reader-token"}
	_, err := readChanges(server.Client(), cfg)
	if err == nil {
		t.Fatal("readChanges: erwarteter Fehler blieb aus")
	}
	if !strings.Contains(err.Error(), "400") {
		t.Fatalf("Fehlertext trägt keinen Statuscode: %v", err)
	}
}
