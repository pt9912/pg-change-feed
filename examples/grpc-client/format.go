package main

import (
	"fmt"

	streamv1 "github.com/pt9912/pg-change-feed/gen/cdc/stream/v1"
)

// formatChange baut die Ausgabezeile einer empfangenen Stream-Nachricht
// (`LH-FA-SST-008`): Tabelle, Operation und das neue Row Image, dazu die
// Kennung, an der sich die Nachricht gegen den Lesezugriffsweg
// `cdc.changes` halten lässt (`change_id`).
func formatChange(c *streamv1.Change) string {
	return fmt.Sprintf("grpc-client: change_id=%s table=%s.%s operation=%s new_image=%s",
		c.GetChangeId(), c.GetSchema(), c.GetTable(), c.GetOperation(), c.GetNewImage())
}
