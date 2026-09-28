package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestTableStatusURLCarriesAllFourFields prüft den Aufbau der Lese-Adresse
// von `GET /tables/status`: alle vier Pflichtfelder stehen als
// Query-Parameter.
func TestTableStatusURLCarriesAllFourFields(t *testing.T) {
	got := TableStatusURL("feed:8080", "quelle-1", "public", "orders", "pub_quelle_1")
	want := "http://feed:8080/tables/status?publication=pub_quelle_1&schema=public&source=quelle-1&table=orders"
	if got != want {
		t.Fatalf("TableStatusURL = %q, want %q", got, want)
	}
}

// TestEnableTableParsesSuccessResponse prüft die Erfolgs-Antwort von
// `POST /tables/enable`, inklusive der sieben Request-Felder.
func TestEnableTableParsesSuccessResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req enableTableRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Request-Body: %v", err)
		}
		if req.TableID != "public.orders" || req.SchemaVersionID != "public.orders-v1" || req.Version != 1 {
			t.Fatalf("unerwarteter Request-Body: %+v", req)
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(enableTableResponse{TableID: req.TableID, Source: req.Source, Schema: req.Schema, Table: req.Table, AlreadyEnabled: false})
	}))
	defer server.Close()

	cfg := config{
		addr: strings.TrimPrefix(server.URL, "http://"), adminToken: "admin-token",
		source: "quelle-1", schema: "public", table: "orders",
		tableID: "public.orders", schemaVersionID: "public.orders-v1", version: 1,
		publication: "pub_quelle_1",
	}
	resp, err := enableTable(server.Client(), cfg)
	if err != nil {
		t.Fatalf("enableTable: %v", err)
	}
	if resp.AlreadyEnabled {
		t.Fatalf("unerwartete Antwort: %+v", resp)
	}
}

// TestEnableTableFailsOnNon2xx prüft den Fehlerpfad.
func TestEnableTableFailsOnNon2xx(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Tabelle existiert an der Quelle nicht"})
	}))
	defer server.Close()

	cfg := config{addr: strings.TrimPrefix(server.URL, "http://"), adminToken: "admin-token"}
	_, err := enableTable(server.Client(), cfg)
	if err == nil {
		t.Fatal("enableTable: erwarteter Fehler blieb aus")
	}
	if !strings.Contains(err.Error(), "404") {
		t.Fatalf("Fehlertext trägt keinen Statuscode: %v", err)
	}
}

// TestDisableTableParsesBothOutcomeFields prüft die Erfolgs-Antwort von
// `POST /tables/disable`.
func TestDisableTableParsesBothOutcomeFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(disableTableResponse{Removed: true, Retained: true})
	}))
	defer server.Close()

	cfg := config{addr: strings.TrimPrefix(server.URL, "http://"), adminToken: "admin-token", source: "quelle-1", schema: "public", table: "orders", publication: "pub_quelle_1"}
	resp, err := disableTable(server.Client(), cfg)
	if err != nil {
		t.Fatalf("disableTable: %v", err)
	}
	if !resp.Removed || !resp.Retained {
		t.Fatalf("unerwartete Antwort: %+v", resp)
	}
}

// TestDisableTableFailsOnNon2xx prüft den Fehlerpfad.
func TestDisableTableFailsOnNon2xx(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "source, schema, table und publication sind Pflichtfelder"})
	}))
	defer server.Close()

	cfg := config{addr: strings.TrimPrefix(server.URL, "http://"), adminToken: "admin-token"}
	_, err := disableTable(server.Client(), cfg)
	if err == nil {
		t.Fatal("disableTable: erwarteter Fehler blieb aus")
	}
	if !strings.Contains(err.Error(), "400") || !strings.Contains(err.Error(), "Pflichtfelder") {
		t.Fatalf("Fehlertext trägt weder Status noch Antworttext: %v", err)
	}
}

// TestTableStatusParsesNeverEnabledBoundary prüft `LH-FA-CFG-003`s
// Boundary: eine nie registrierte Tabelle trägt beide Felder `false`.
func TestTableStatusParsesNeverEnabledBoundary(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(tableStatusResponse{Enabled: false, Retained: false})
	}))
	defer server.Close()

	cfg := config{addr: strings.TrimPrefix(server.URL, "http://"), token: "reader-token", source: "quelle-1", schema: "public", table: "nie_aktiviert", publication: "pub_quelle_1"}
	resp, err := tableStatus(server.Client(), cfg)
	if err != nil {
		t.Fatalf("tableStatus: %v", err)
	}
	if resp.Enabled || resp.Retained {
		t.Fatalf("unerwartete Antwort: %+v", resp)
	}
}

// TestTableStatusFailsOnNon2xx prüft den Fehlerpfad.
func TestTableStatusFailsOnNon2xx(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "source, schema, table und publication sind Pflichtfelder"})
	}))
	defer server.Close()

	cfg := config{addr: strings.TrimPrefix(server.URL, "http://"), token: "reader-token"}
	_, err := tableStatus(server.Client(), cfg)
	if err == nil {
		t.Fatal("tableStatus: erwarteter Fehler blieb aus")
	}
	if !strings.Contains(err.Error(), "400") || !strings.Contains(err.Error(), "Pflichtfelder") {
		t.Fatalf("Fehlertext trägt weder Status noch Antworttext: %v", err)
	}
}
