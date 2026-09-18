// Command rolloutguard entscheidet, ob ein zweiter make schema-rollout-Lauf
// gegen ein bereits migriertes Ziel ausschließlich auf den bekannten
// Fremdobjekten blockiert (ADR-0043) — zentral, für jeden
// Aufrufer gleich, statt bei jedem Aufrufer einzeln
// (docs/plan/planning/observations/BEO-PGC/schema-rollout-fremdobjekte/).
// Es liest den strukturierten JSON-Report eines vorgelagerten
// `schema migrate --plan-only`-Laufs und meldet über den Exit-Code, ob das
// Makefile-Target den nachfolgenden `--execute`-Schritt zusätzlich mit
// `--allow-destructive` laufen lassen darf.
package main

import "encoding/json"

// report bildet nur die Felder des d-migrate-Plan-Reports ab, die dieser
// Guard braucht (reales Schema, gemessen am gepinnten Image via
// --plan-only --report, siehe Slice-Plan §3). d-migrate trägt weitere
// Felder (summary, diagnostics, …), die hier bewusst nicht dekodiert
// werden.
type report struct {
	Status   string `json:"status"`
	Blockers []struct {
		Reason       string   `json:"reason"`
		OperationIDs []string `json:"operationIds"`
	} `json:"blockers"`
	Operations []struct {
		ID         string   `json:"id"`
		Kind       string   `json:"kind"`
		ObjectType string   `json:"objectType"`
		Path       []string `json:"path"`
	} `json:"operations"`
}

// foreignObject identifiziert eine Blocker-Operation stabil über Kind,
// Objekttyp und Pfad — nicht über die id-Hashes des Reports, die sich mit
// dem Inhalt der betroffenen Routine/View ändern können. Grenze: path
// trägt strings.Join(op.Path, ".") — kollisionsfrei nur, solange kein
// Pfadsegment selbst einen Punkt enthält und alle Objekte im einen
// Ziel-Schema `cdc` liegen (beides für den aktuellen Bestand der Fall).
type foreignObject struct {
	kind       string
	objectType string
	path       string
}

func parseReport(data []byte) (report, error) {
	var r report
	err := json.Unmarshal(data, &r)
	return r, err
}
