package main

import "testing"

// TestTablesURLCarriesBothRequiredFilters prüft den Aufbau der HTTP-Abfrage
// (`LH-FA-SST-006`): die Lese-Adresse trägt beide Pflichtfelder des Endpunkts
// als Query-Parameter — ohne sie endet `GET /tables` mit `400`.
func TestTablesURLCarriesBothRequiredFilters(t *testing.T) {
	got := TablesURL("feed:8080", "quelle-1", "pub_quelle_1")
	want := "http://feed:8080/tables?publication=pub_quelle_1&source=quelle-1"
	if got != want {
		t.Fatalf("TablesURL = %q, want %q", got, want)
	}
}

// TestTablesURLKeepsHostBoundary prüft eine Grenze der reinen Funktion: die
// Adresse ist der Horch-Host aus `CDC_HTTP_ADDR`; die Funktion setzt kein
// eigenes Schema auf einen bereits absoluten Wert.
func TestTablesURLKeepsHostBoundary(t *testing.T) {
	got := TablesURL("localhost:9090", "src-e2e", "cdc_pub")
	want := "http://localhost:9090/tables?publication=cdc_pub&source=src-e2e"
	if got != want {
		t.Fatalf("TablesURL = %q, want %q", got, want)
	}
}

// TestTablesURLEscapesReservedCharacters prüft, dass die Bestandteile als
// Query-Werte kodiert werden statt aneinandergereiht: eine `source_id` mit
// reservierten Zeichen bleibt ein einziger Parameterwert und zerfällt nicht
// in mehrere Parameter.
func TestTablesURLEscapesReservedCharacters(t *testing.T) {
	got := TablesURL("feed:8080", "quelle & test", "pub/1")
	want := "http://feed:8080/tables?publication=pub%2F1&source=quelle+%26+test"
	if got != want {
		t.Fatalf("TablesURL = %q, want %q", got, want)
	}
}
