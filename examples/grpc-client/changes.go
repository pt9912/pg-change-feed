package main

import (
	"fmt"
	"strings"

	administrationv1 "github.com/pt9912/pg-change-feed/gen/cdc/administration/v1"
)

// readChangesRequest bildet die Anfrage von `ReadChanges` aus der
// Konfiguration; ein leerer `target` ist kein Filter.
func readChangesRequest(cfg config) *administrationv1.ReadChangesRequest {
	return &administrationv1.ReadChangesRequest{
		Source: cfg.source, Schema: cfg.schema, Table: cfg.table, Target: cfg.target,
		From: cfg.from, To: cfg.to, Limit: cfg.limit,
	}
}

// readChanges ruft die reader-RPC `ReadChanges` auf (`LH-FA-SST-006`): einen
// begrenzten Bereich persistierter Änderungen, dieselbe Filter- und
// Bereichs-Semantik wie `GET /changes` — ein leerer `from`/`to`/`limit`-Wert
// trägt `0` (nicht gesetzt).
func readChanges(client administrationv1.AdministrationClient, cfg config) (string, error) {
	ctx, cancel := callCtx(cfg.token)
	defer cancel()
	resp, err := client.ReadChanges(ctx, readChangesRequest(cfg))
	if err != nil {
		return "", err
	}
	changes := resp.GetChanges()
	lines := make([]string, 0, len(changes)+1)
	lines = append(lines, fmt.Sprintf("grpc-client: changes=%d", len(changes)))
	for _, c := range changes {
		lines = append(lines, fmt.Sprintf("  change_id=%s table=%s.%s operation=%s commit_position=%d origin=%s new_image=%s",
			c.GetChangeId(), c.GetSchema(), c.GetTable(), c.GetOperation(), c.GetCommitPosition(), c.GetOrigin(), c.GetNewImage()))
	}
	return strings.Join(lines, "\n"), nil
}
