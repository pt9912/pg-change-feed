// Command rolloutguard entscheidet, was ein make schema-rollout-Lauf gegen
// ein bereits migriertes Ziel zusätzlich zum regulären `--execute` tun darf
// (ADR-0043, ADR-0114): `--allow-destructive`, wenn ausschließlich die
// bekannten Fremdobjekte blockieren, und einen Vorlauf `DROP VIEW`, wenn
// eine im neutralen Modell deklarierte View ihre Signatur ändert — zentral,
// für jeden Aufrufer gleich, statt bei jedem Aufrufer einzeln
// (docs/plan/planning/observations/BEO-PGC/schema-rollout-fremdobjekte/).
// Es liest den strukturierten JSON-Report eines vorgelagerten
// `schema migrate --plan-only`-Laufs und meldet über Exit-Code und
// stdout-Zeilen, was das Makefile-Target tun darf.
package main

import "encoding/json"

// report bildet nur die Felder des d-migrate-Plan-Reports ab, die dieser
// Guard braucht (reales Schema, gemessen am gepinnten Image via
// --plan-only --report). d-migrate trägt weitere Felder (summary,
// statements, …), die hier bewusst nicht dekodiert werden.
type report struct {
	Status      string       `json:"status"`
	Blockers    []blocker    `json:"blockers"`
	Operations  []operation  `json:"operations"`
	Diagnostics []diagnostic `json:"diagnostics"`
}

type blocker struct {
	Reason       string   `json:"reason"`
	OperationIDs []string `json:"operationIds"`
}

type operation struct {
	ID         string   `json:"id"`
	Kind       string   `json:"kind"`
	ObjectType string   `json:"objectType"`
	Path       []string `json:"path"`
}

// diagnostic trägt den Diagnosecode zu einer Operation. Das Feld
// blockers[].diagnosticCodes ist im realen Report leer (gemessen); die
// Zuordnung von Code zu Operation läuft ausschließlich über operationId.
type diagnostic struct {
	Code        string `json:"code"`
	OperationID string `json:"operationId"`
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
