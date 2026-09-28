package main

import (
	"strings"
	"testing"
)

// TestValidateRejectsUnknownVerb prüft, dass ein unbekannter `-verb`-Wert
// vor jedem Netzwerkaufruf abbricht.
func TestValidateRejectsUnknownVerb(t *testing.T) {
	err := validate(config{addr: "feed:9090", verb: "unbekannt"})
	if err == nil {
		t.Fatal("validate: erwarteter Fehler blieb aus")
	}
	if !strings.Contains(err.Error(), "unbekanntes -verb") {
		t.Fatalf("unerwarteter Fehlertext: %v", err)
	}
}

// TestValidateRequiresAddrBeforeVerbFields prüft, dass die fehlende
// gRPC-Adresse vor jeder verb-spezifischen Prüfung erkannt wird.
func TestValidateRequiresAddrBeforeVerbFields(t *testing.T) {
	err := validate(config{verb: "list-tables"})
	if err == nil || !strings.Contains(err.Error(), "CDC_GRPC_ADDR") {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
}

// TestValidateStreamDefaultAcceptsReaderTokenWithoutFilter prüft das
// unveränderte Default-Verhalten: `stream` ohne Filter (`-schema`/`-table`
// leer) ist mit einem gesetzten Reader-Token gültig (`ADR-0133`).
func TestValidateStreamDefaultAcceptsReaderTokenWithoutFilter(t *testing.T) {
	cfg := config{addr: "feed:9090", token: "reader-token", verb: "stream"}
	if err := validate(cfg); err != nil {
		t.Fatalf("validate: %v", err)
	}
}

// TestValidateStreamRequiresReaderTokenNotAdminToken prüft die
// Rechtsklassen-Bindung des Streams: ein gesetztes Admin-Token allein genügt
// nicht.
func TestValidateStreamRequiresReaderTokenNotAdminToken(t *testing.T) {
	cfg := config{addr: "feed:9090", adminToken: "admin-token", verb: "stream"}
	err := validate(cfg)
	if err == nil || !strings.Contains(err.Error(), "CDC_API_TOKEN_READER") {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
}

// TestValidateEnableTableAllowsDefaultTableID prüft, dass `enable-table`
// ohne `-table-id`/`-schema-version-id` gültig ist — beide Felder tragen
// einen aus Schema und Tabelle abgeleiteten Default (`dispatchAdmin`).
func TestValidateEnableTableAllowsDefaultTableID(t *testing.T) {
	cfg := config{addr: "feed:9090", adminToken: "admin-token", verb: "enable-table", source: "quelle-1", schema: "public", table: "orders", publication: "pub_quelle_1"}
	if err := validate(cfg); err != nil {
		t.Fatalf("validate: %v", err)
	}
}

// TestValidateRunRetentionAllowsZeroMinAge prüft, dass `min-age-nanos=0`
// (kein zeitliches Mindestalter) kein Pflichtfeld-Fehler ist.
func TestValidateRunRetentionAllowsZeroMinAge(t *testing.T) {
	cfg := config{addr: "feed:9090", adminToken: "admin-token", verb: "run-retention", source: "quelle-1"}
	if err := validate(cfg); err != nil {
		t.Fatalf("validate: %v", err)
	}
}

// TestValidateAcknowledgeConsumerRequiresAdminTokenNotReaderToken prüft die
// Rechtsklassen-Bindung: `acknowledge-consumer` ist ein `admin`-Endpunkt —
// ein gesetztes Reader-Token allein genügt nicht.
func TestValidateAcknowledgeConsumerRequiresAdminTokenNotReaderToken(t *testing.T) {
	cfg := config{addr: "feed:9090", token: "reader-token", verb: "acknowledge-consumer", consumerID: "consumer-1", source: "quelle-1"}
	err := validate(cfg)
	if err == nil || !strings.Contains(err.Error(), "Admin-Token") {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
}

// TestValidateRegisterConsumerRequiresAdminTokenNotReaderToken prüft dieselbe
// Rechtsklassen-Bindung für `register-consumer` (`admin`-Endpunkt).
func TestValidateRegisterConsumerRequiresAdminTokenNotReaderToken(t *testing.T) {
	cfg := config{addr: "feed:9090", token: "reader-token", verb: "register-consumer", consumerID: "consumer-1", name: "Consumer"}
	err := validate(cfg)
	if err == nil || !strings.Contains(err.Error(), "Admin-Token") {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
}

// TestValidateRemoveConsumerRequiresAdminTokenNotReaderToken prüft dieselbe
// Rechtsklassen-Bindung für `remove-consumer` (`admin`-Endpunkt).
func TestValidateRemoveConsumerRequiresAdminTokenNotReaderToken(t *testing.T) {
	cfg := config{addr: "feed:9090", token: "reader-token", verb: "remove-consumer", consumerID: "consumer-1"}
	err := validate(cfg)
	if err == nil || !strings.Contains(err.Error(), "Admin-Token") {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
}

// TestValidateGetConsumerPositionRequiresReaderTokenNotAdminToken prüft die
// Rechtsklassen-Bindung: `get-consumer-position` ist ein `reader`-Endpunkt —
// ein gesetztes Admin-Token allein genügt nicht.
func TestValidateGetConsumerPositionRequiresReaderTokenNotAdminToken(t *testing.T) {
	cfg := config{addr: "feed:9090", adminToken: "admin-token", verb: "get-consumer-position", consumerID: "consumer-1"}
	err := validate(cfg)
	if err == nil || !strings.Contains(err.Error(), "CDC_API_TOKEN_READER") {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
}

// TestValidateEnableTableRequiresAdminTokenNotReaderToken prüft dieselbe
// Rechtsklassen-Bindung für `enable-table` (`admin`-Endpunkt).
func TestValidateEnableTableRequiresAdminTokenNotReaderToken(t *testing.T) {
	cfg := config{addr: "feed:9090", token: "reader-token", verb: "enable-table", source: "quelle-1", schema: "public", table: "orders", publication: "pub_quelle_1"}
	err := validate(cfg)
	if err == nil || !strings.Contains(err.Error(), "Admin-Token") {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
}

// TestValidateDisableTableRequiresAdminTokenNotReaderToken prüft dieselbe
// Rechtsklassen-Bindung für `disable-table` (`admin`-Endpunkt).
func TestValidateDisableTableRequiresAdminTokenNotReaderToken(t *testing.T) {
	cfg := config{addr: "feed:9090", token: "reader-token", verb: "disable-table", source: "quelle-1", schema: "public", table: "orders", publication: "pub_quelle_1"}
	err := validate(cfg)
	if err == nil || !strings.Contains(err.Error(), "Admin-Token") {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
}

// TestValidateGetTableStatusRequiresReaderTokenNotAdminToken prüft dieselbe
// Rechtsklassen-Bindung für `get-table-status` (`reader`-Endpunkt).
func TestValidateGetTableStatusRequiresReaderTokenNotAdminToken(t *testing.T) {
	cfg := config{addr: "feed:9090", adminToken: "admin-token", verb: "get-table-status", source: "quelle-1", schema: "public", table: "orders", publication: "pub_quelle_1"}
	err := validate(cfg)
	if err == nil || !strings.Contains(err.Error(), "CDC_API_TOKEN_READER") {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
}

// TestValidateListTablesRequiresReaderTokenNotAdminToken prüft dieselbe
// Rechtsklassen-Bindung für `list-tables` (`reader`-Endpunkt).
func TestValidateListTablesRequiresReaderTokenNotAdminToken(t *testing.T) {
	cfg := config{addr: "feed:9090", adminToken: "admin-token", verb: "list-tables", source: "quelle-1", publication: "pub_quelle_1"}
	err := validate(cfg)
	if err == nil || !strings.Contains(err.Error(), "CDC_API_TOKEN_READER") {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
}

// TestValidateRunRetentionRequiresAdminTokenNotReaderToken prüft dieselbe
// Rechtsklassen-Bindung für `run-retention` (`admin`-Endpunkt).
func TestValidateRunRetentionRequiresAdminTokenNotReaderToken(t *testing.T) {
	cfg := config{addr: "feed:9090", token: "reader-token", verb: "run-retention", source: "quelle-1"}
	err := validate(cfg)
	if err == nil || !strings.Contains(err.Error(), "Admin-Token") {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
}

// TestValidateDiagnoseRequiresReaderTokenNotAdminToken prüft dieselbe
// Rechtsklassen-Bindung für `diagnose` (`reader`-Endpunkt).
func TestValidateDiagnoseRequiresReaderTokenNotAdminToken(t *testing.T) {
	cfg := config{addr: "feed:9090", adminToken: "admin-token", verb: "diagnose", source: "quelle-1"}
	err := validate(cfg)
	if err == nil || !strings.Contains(err.Error(), "CDC_API_TOKEN_READER") {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
}

// TestValidateReadChangesRequiresReaderTokenNotAdminToken prüft die
// Rechtsklassen-Bindung: `read-changes` ist ein `reader`-Endpunkt — ein
// gesetztes Admin-Token allein genügt nicht.
func TestValidateReadChangesRequiresReaderTokenNotAdminToken(t *testing.T) {
	cfg := config{addr: "feed:9090", adminToken: "admin-token", verb: "read-changes", source: "quelle-1"}
	err := validate(cfg)
	if err == nil || !strings.Contains(err.Error(), "CDC_API_TOKEN_READER") {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
}

// TestValidateDiagnoseRequiresSource prüft, dass `diagnose` ohne `-source`
// abbricht, bevor ein Netzwerkaufruf versucht wird.
func TestValidateDiagnoseRequiresSource(t *testing.T) {
	cfg := config{addr: "feed:9090", token: "reader-token", verb: "diagnose"}
	err := validate(cfg)
	if err == nil || !strings.Contains(err.Error(), "-source ist Pflicht") {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
}

// TestDispatchAdminRejectsUnknownVerb prüft, dass `dispatchAdmin` denselben
// Schutz wie `validate` trägt — verteidigend, für den Fall eines Aufrufs
// ohne vorherige Validierung. Ein `nil`-Client wird dabei nie erreicht.
func TestDispatchAdminRejectsUnknownVerb(t *testing.T) {
	_, err := dispatchAdmin(nil, config{verb: "unbekannt"})
	if err == nil || !strings.Contains(err.Error(), "unbekanntes -verb") {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
}
