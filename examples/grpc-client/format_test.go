package main

import (
	"testing"

	streamv1 "github.com/pt9912/pg-change-feed/gen/cdc/stream/v1"
)

// TestFormatChangeCarriesIdentityAndPayload prüft den Aufbau der Ausgabezeile
// (`LH-FA-SST-008`): Tabelle, Operation und `change_id` stehen darin, das
// neue Row Image wird unverändert angehängt.
func TestFormatChangeCarriesIdentityAndPayload(t *testing.T) {
	c := &streamv1.Change{
		ChangeId:  "c-1",
		Schema:    "public",
		Table:     "orders",
		Operation: "INSERT",
		NewImage:  []byte(`{"id":1}`),
	}
	got := formatChange(c)
	want := `grpc-client: change_id=c-1 table=public.orders operation=INSERT new_image={"id":1}`
	if got != want {
		t.Fatalf("formatChange = %q, want %q", got, want)
	}
}

// TestFormatChangeHandlesEmptyNewImage prüft den `DELETE`-Fall (`SPEC-020`):
// das neue Row Image ist bei einer Löschung leer, die Ausgabezeile trägt es
// als leeren Wert statt einen Platzhalter zu erfinden.
func TestFormatChangeHandlesEmptyNewImage(t *testing.T) {
	c := &streamv1.Change{
		ChangeId:  "c-2",
		Schema:    "public",
		Table:     "orders",
		Operation: "DELETE",
	}
	got := formatChange(c)
	want := `grpc-client: change_id=c-2 table=public.orders operation=DELETE new_image=`
	if got != want {
		t.Fatalf("formatChange = %q, want %q", got, want)
	}
}
