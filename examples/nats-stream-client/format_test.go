package main

import (
	"encoding/json"
	"testing"
)

// TestFormatChangeCarriesIdentityAndPayload prüft den Aufbau der Ausgabezeile
// (LH-FA-SST-008): Tabelle, Operation und change_id stehen darin, das neue
// Row Image wird unverändert angehängt.
func TestFormatChangeCarriesIdentityAndPayload(t *testing.T) {
	c := streamMessage{
		ChangeID:  "c-1",
		Schema:    "public",
		Table:     "orders",
		Operation: "INSERT",
		NewImage:  []byte(`{"id":1}`),
	}
	got := formatChange(c)
	want := `nats-stream-client: change_id=c-1 table=public.orders operation=INSERT new_image={"id":1}`
	if got != want {
		t.Fatalf("formatChange = %q, want %q", got, want)
	}
}

// TestFormatChangeHandlesMissingNewImage prüft den DELETE-Fall (SPEC-024):
// ein fehlendes Row Image trägt auf dem Draht das JSON-Literal `null`
// (dieselbe Übersetzung wie `natsstream.rowImage` auf der Erzeugerseite) —
// nach dem Dekodieren stehen genau diese vier Bytes im Feld, kein leerer
// Wert.
func TestFormatChangeHandlesMissingNewImage(t *testing.T) {
	c := streamMessage{
		ChangeID:  "c-2",
		Schema:    "public",
		Table:     "orders",
		Operation: "DELETE",
		NewImage:  json.RawMessage("null"),
	}
	got := formatChange(c)
	want := `nats-stream-client: change_id=c-2 table=public.orders operation=DELETE new_image=null`
	if got != want {
		t.Fatalf("formatChange = %q, want %q", got, want)
	}
}
