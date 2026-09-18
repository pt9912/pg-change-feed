package main

import (
	"encoding/json"
	"fmt"
)

// streamMessage trägt dasselbe Nachrichtenschema wie der Publisher
// (internal/adapters/driven/natsstream) und der SSE-Adapter (SPEC-021/
// SPEC-024) — hier eigenständig geführt: dieses Beispiel importiert keinen
// privaten Paketbaum dieses Repositories (SPEC-023).
type streamMessage struct {
	ChangeID      string          `json:"change_id"`
	TransactionID string          `json:"transaction_id"`
	SourceTableID string          `json:"source_table_id"`
	Sequence      int64           `json:"sequence"`
	Operation     string          `json:"operation"`
	OldImage      json.RawMessage `json:"old_image"`
	NewImage      json.RawMessage `json:"new_image"`
	SchemaVersion string          `json:"schema_version"`
	Schema        string          `json:"schema"`
	Table         string          `json:"table"`
}

// formatChange baut die Ausgabezeile einer empfangenen Stream-Nachricht
// (LH-FA-SST-008): Tabelle, Operation und das neue Row Image, dazu die
// Kennung, an der sich die Nachricht gegen den Lesezugriffsweg cdc.changes
// halten lässt (change_id) — dieselbe Form wie beim gRPC-Beispiel
// (examples/grpc-client/format.go).
func formatChange(c streamMessage) string {
	return fmt.Sprintf("nats-stream-client: change_id=%s table=%s.%s operation=%s new_image=%s",
		c.ChangeID, c.Schema, c.Table, c.Operation, c.NewImage)
}
