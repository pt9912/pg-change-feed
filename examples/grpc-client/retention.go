package main

import (
	"fmt"

	administrationv1 "github.com/pt9912/pg-change-feed/gen/cdc/administration/v1"
)

// runRetention ruft die admin-RPC `RunRetention` auf (`LH-FA-RET-002`):
// `MinAgeNanos` 0 heißt „kein zeitliches Mindestalter" und ist gültig.
func runRetention(client administrationv1.AdministrationClient, cfg config) (string, error) {
	ctx, cancel := callCtx(cfg.adminToken)
	defer cancel()
	resp, err := client.RunRetention(ctx, &administrationv1.RunRetentionRequest{Source: cfg.source, MinAgeNanos: cfg.minAgeNanos})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("grpc-client: deleted=%d", resp.GetDeleted()), nil
}
