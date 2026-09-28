package main

import (
	"strings"
	"testing"
)

// TestValidateRejectsUnknownVerb prüft, dass ein unbekannter `-verb`-Wert
// vor jedem Netzwerkaufruf abbricht.
func TestValidateRejectsUnknownVerb(t *testing.T) {
	err := validate(config{addr: "feed:8080", verb: "unbekannt"})
	if err == nil {
		t.Fatal("validate: erwarteter Fehler blieb aus")
	}
	if !strings.Contains(err.Error(), "unbekanntes -verb") {
		t.Fatalf("unerwarteter Fehlertext: %v", err)
	}
}

// TestValidateRequiresAddrBeforeVerbFields prüft, dass die fehlende
// HTTP-Adresse vor jeder verb-spezifischen Prüfung erkannt wird.
func TestValidateRequiresAddrBeforeVerbFields(t *testing.T) {
	err := validate(config{verb: "tables"})
	if err == nil || !strings.Contains(err.Error(), "CDC_HTTP_ADDR") {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
}

// TestValidateEnableTableAllowsDefaultTableID prüft, dass `enable-table`
// ohne `-table-id`/`-schema-version-id` gültig ist — beide Felder tragen
// einen aus Schema und Tabelle abgeleiteten Default (`dispatch`,
// `internal/bootstrap/wiring.go` `administrationTableID`).
func TestValidateEnableTableAllowsDefaultTableID(t *testing.T) {
	cfg := config{addr: "feed:8080", adminToken: "admin-token", verb: "enable-table", source: "quelle-1", schema: "public", table: "orders", publication: "pub_quelle_1"}
	if err := validate(cfg); err != nil {
		t.Fatalf("validate: %v", err)
	}
}

// TestValidateRetentionRunAllowsZeroMinAge prüft, dass `min-age-nanos=0`
// (kein zeitliches Mindestalter) kein Pflichtfeld-Fehler ist.
func TestValidateRetentionRunAllowsZeroMinAge(t *testing.T) {
	cfg := config{addr: "feed:8080", adminToken: "admin-token", verb: "retention-run", source: "quelle-1"}
	if err := validate(cfg); err != nil {
		t.Fatalf("validate: %v", err)
	}
}

// TestValidateAcknowledgeRequiresAdminTokenNotReaderToken prüft die
// Rechtsklassen-Bindung: `acknowledge` ist ein `admin`-Endpunkt — ein
// gesetztes reader-Token allein genügt nicht.
func TestValidateAcknowledgeRequiresAdminTokenNotReaderToken(t *testing.T) {
	cfg := config{addr: "feed:8080", token: "reader-token", verb: "acknowledge", consumerID: "consumer-1", source: "quelle-1"}
	err := validate(cfg)
	if err == nil || !strings.Contains(err.Error(), "Admin-Token") {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
}

// TestDispatchRejectsUnknownVerb prüft, dass `dispatch` denselben Schutz
// wie `validate` trägt — verteidigend, für den Fall eines Aufrufs ohne
// vorherige Validierung.
func TestDispatchRejectsUnknownVerb(t *testing.T) {
	_, err := dispatch(nil, config{verb: "unbekannt"})
	if err == nil || !strings.Contains(err.Error(), "unbekanntes -verb") {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
}
