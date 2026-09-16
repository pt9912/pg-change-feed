package mapper

import (
	"testing"

	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Whitebox-Test (`package mapper`, nicht `mapper_test`): `setSchemaVersion`
// ist ein unexportiertes Nach trag-Detail der dynamischen
// Re-Versionierung (`ADR-0015` Folgepflicht, `ADR-0059` Teilfrage 3). Sein
// Vertrag für eine **nicht** getragene Bindung — der Nachtrag belebt sie
// nicht neu — ist über `observeRelation` nicht erreichbar: der Aufrufer
// löst die Bindung unmittelbar davor auf. Der Test greift deshalb direkt
// zu, statt die Methode für den Test zu exportieren.
//
// Rot färbende Mutation: in `setSchemaVersion` den Zug
// `if !activated { return }` entfernen — dann legt der Nachtrag eine
// Bindung mit leerer Tabellen-Kennung an, und dieser Test fällt.
func TestSetSchemaVersionLeavesUnboundBindingUntouched(t *testing.T) {
	assembler, err := NewAssembler("src-1", map[string]TableBinding{
		"public.feed": {TableID: "tbl-1", SchemaVersion: "sv-1"},
	}, nil)
	if err != nil {
		t.Fatalf("NewAssembler: %v", err)
	}

	assembler.setSchemaVersion("public.orders", model.SchemaVersionID("sv-9"))

	if len(assembler.tables) != 1 {
		t.Fatalf("Bindungen = %+v, wollen nur public.feed (der Nachtrag belebt keine Bindung)", assembler.tables)
	}
	if got := assembler.tables["public.feed"].SchemaVersion; got != "sv-1" {
		t.Fatalf("Schema-Version von public.feed = %q, wollen sv-1 (unberührt)", got)
	}
}
