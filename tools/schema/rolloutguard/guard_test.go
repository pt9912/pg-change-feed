package main

import "testing"

// knownBlockedReport spiegelt real gemessene Report-Felder (--plan-only
// gegen ein bereits vollständig migriertes Ziel, siehe Slice-Plan §3):
// genau die sechs bekannten Fremdobjekt-Blocker, sonst nichts.
func knownBlockedReport() report {
	r := report{Status: "blocked"}
	r.Blockers = []struct {
		Reason       string   `json:"reason"`
		OperationIDs []string `json:"operationIds"`
	}{
		{Reason: destructiveConfirmationReason, OperationIDs: []string{
			"DropFunction:FUNCTION:a1:a2",
			"DropFunction:FUNCTION:b1:b2",
			"DropFunction:FUNCTION:c1:c2",
			"DropFunction:FUNCTION:d1:d2",
			"DropView:VIEW:e1:e2",
			"DropView:VIEW:f1:f2",
		}},
	}
	r.Operations = []struct {
		ID         string   `json:"id"`
		Kind       string   `json:"kind"`
		ObjectType string   `json:"objectType"`
		Path       []string `json:"path"`
	}{
		{ID: "DropFunction:FUNCTION:a1:a2", Kind: "DropFunction", ObjectType: "FUNCTION", Path: []string{"enable_table(in:text,in:text,in:text)"}},
		{ID: "DropFunction:FUNCTION:b1:b2", Kind: "DropFunction", ObjectType: "FUNCTION", Path: []string{"disable_table(in:text,in:text,in:text)"}},
		{ID: "DropFunction:FUNCTION:c1:c2", Kind: "DropFunction", ObjectType: "FUNCTION", Path: []string{"exclude_column(in:text,in:text,in:text,in:text)"}},
		{ID: "DropFunction:FUNCTION:d1:d2", Kind: "DropFunction", ObjectType: "FUNCTION", Path: []string{"include_column(in:text,in:text,in:text,in:text)"}},
		{ID: "DropView:VIEW:e1:e2", Kind: "DropView", ObjectType: "VIEW", Path: []string{"heartbeat"}},
		{ID: "DropView:VIEW:f1:f2", Kind: "DropView", ObjectType: "VIEW", Path: []string{"metrics"}},
	}
	return r
}

// TestDecideSkipsWhenAllBlockersKnown prüft den Regelfall: ein zweiter Lauf
// gegen ein bereits vollständig migriertes Ziel trägt ausschließlich die
// sechs bekannten Fremdobjekt-Blocker — decide erlaubt den Skip.
func TestDecideSkipsWhenAllBlockersKnown(t *testing.T) {
	skip, reason := decide(knownBlockedReport())
	if !skip {
		t.Fatalf("decide() skip = false, reason %q — wollte true", reason)
	}
}

// TestDecideRefusesUnknownDestructiveBlocker prüft den DoD-Negativfall
// (Slice-Plan §2 zweiter Punkt): eine künstlich per ALTER TABLE … ADD
// COLUMN hinzugefügte, nicht deklarierte Spalte erzeugt real einen
// zusätzlichen, unbekannten DropColumn-Blocker neben den sechs bekannten
// (gemessen gegen das gepinnte d-migrate-Image) — decide verweigert den
// Skip, auch im Mischfall mit sechs sonst bekannten Blockern.
func TestDecideRefusesUnknownDestructiveBlocker(t *testing.T) {
	r := knownBlockedReport()
	r.Blockers[0].OperationIDs = append(r.Blockers[0].OperationIDs, "DropColumn:COLUMN:g1:g2")
	r.Operations = append(r.Operations, struct {
		ID         string   `json:"id"`
		Kind       string   `json:"kind"`
		ObjectType string   `json:"objectType"`
		Path       []string `json:"path"`
	}{ID: "DropColumn:COLUMN:g1:g2", Kind: "DropColumn", ObjectType: "COLUMN", Path: []string{"source", "_scratch_test_col"}})

	skip, _ := decide(r)
	if skip {
		t.Fatal("decide() skip = true — wollte false, weil ein Blocker nicht auf der bekannten Liste steht")
	}
}

// TestDecideRefusesUnknownBlockerReason prüft, dass decide nur die eine
// bekannte Blocker-Klasse automatisch auflöst — jede andere gemeldete
// Blocker-Ursache bleibt unaufgelöst.
func TestDecideRefusesUnknownBlockerReason(t *testing.T) {
	r := knownBlockedReport()
	r.Blockers[0].Reason = "SOME_OTHER_REASON"

	skip, _ := decide(r)
	if skip {
		t.Fatal("decide() skip = true — wollte false bei unbekannter Blocker-Klasse")
	}
}

// TestDecideRefusesMissingOperation prüft die Verteidigungsgrenze: eine
// Blocker-operationId ohne zugehörigen Eintrag in operations[] gilt als
// unbekannt, nicht als übersehbar.
func TestDecideRefusesMissingOperation(t *testing.T) {
	r := knownBlockedReport()
	r.Blockers[0].OperationIDs = append(r.Blockers[0].OperationIDs, "DropTable:TABLE:zz:zz")

	skip, _ := decide(r)
	if skip {
		t.Fatal("decide() skip = true — wollte false bei einer operationId ohne Eintrag in operations[]")
	}
}

// TestDecideRefusesEmptyBlockers prüft die Randbedingung: ein Report ohne
// jeden Blocker trägt keinen Skip-Grund — dieser Fall entsteht in der
// Praxis nicht (der Aufrufer ruft decide nur bei Exit 8 auf), bleibt aber
// eine explizite Zusage der Funktion selbst.
func TestDecideRefusesEmptyBlockers(t *testing.T) {
	skip, _ := decide(report{Status: "ok"})
	if skip {
		t.Fatal("decide() skip = true — wollte false bei leerem Blockers[]")
	}
}
