package main

import "testing"

// TestSubjectDerivesTableGranularSubject prüft die Ableitung des Subjekts aus
// den drei Bestandteilen (`SPEC-017`): Quelle, Schema und Tabelle werden in
// das vierstufige Subjekt übersetzt, nicht handgetippt.
func TestSubjectDerivesTableGranularSubject(t *testing.T) {
	got := Subject("quelle-1", "public", "orders")
	want := "cdc.changes.quelle-1.public.orders"
	if got != want {
		t.Fatalf("Subject = %q, want %q", got, want)
	}
}

// TestSubjectKeepsDistinctTables prüft die Tabellen-Granularität: zwei
// Tabellen derselben Quelle ergeben zwei verschiedene Subjekte (`ADR-0056`).
func TestSubjectKeepsDistinctTables(t *testing.T) {
	orders := Subject("quelle-1", "public", "orders")
	invoices := Subject("quelle-1", "public", "invoices")
	if orders == invoices {
		t.Fatalf("zwei Tabellen derselben Quelle liefern dasselbe Subjekt: %q", orders)
	}
}

// TestChangesURLCarriesFiltersOnReadPath prüft den Aufbau der HTTP-Abfrage
// (`LH-FA-SST-006`): die Lese-Adresse trägt die drei Filter aus dem
// Wecksignal-Subjekt als Query-Parameter.
func TestChangesURLCarriesFiltersOnReadPath(t *testing.T) {
	got := ChangesURL("feed:8080", "quelle-1", "public", "orders")
	want := "http://feed:8080/changes?schema=public&source=quelle-1&table=orders"
	if got != want {
		t.Fatalf("ChangesURL = %q, want %q", got, want)
	}
}

// TestChangesURLKeepsHostBoundary prüft eine Grenze der reinen Funktion: die
// Adresse ist der Horch-Host aus `CDC_HTTP_ADDR`; die Funktion setzt kein
// eigenes Schema auf einen bereits absoluten Wert.
func TestChangesURLKeepsHostBoundary(t *testing.T) {
	got := ChangesURL("localhost:9090", "src-e2e", "public", "feed_e2e")
	want := "http://localhost:9090/changes?schema=public&source=src-e2e&table=feed_e2e"
	if got != want {
		t.Fatalf("ChangesURL = %q, want %q", got, want)
	}
}
