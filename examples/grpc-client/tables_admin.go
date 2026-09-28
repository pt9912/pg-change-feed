package main

import (
	"fmt"
	"strings"

	administrationv1 "github.com/pt9912/pg-change-feed/gen/cdc/administration/v1"
)

// enableTable ruft die admin-RPC `EnableTable` auf (`LH-FA-CFG-001`).
func enableTable(client administrationv1.AdministrationClient, cfg config) (string, error) {
	ctx, cancel := callCtx(cfg.adminToken)
	defer cancel()
	resp, err := client.EnableTable(ctx, &administrationv1.EnableTableRequest{
		Source:          cfg.source,
		Schema:          cfg.schema,
		Table:           cfg.table,
		TableId:         cfg.tableID,
		SchemaVersionId: cfg.schemaVersionID,
		Version:         cfg.version,
		Publication:     cfg.publication,
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("grpc-client: table_id=%s source=%s schema=%s table=%s already_enabled=%v",
		resp.GetTableId(), resp.GetSource(), resp.GetSchema(), resp.GetTable(), resp.GetAlreadyEnabled()), nil
}

// disableTable ruft die admin-RPC `DisableTable` auf (`LH-FA-CFG-002`).
func disableTable(client administrationv1.AdministrationClient, cfg config) (string, error) {
	ctx, cancel := callCtx(cfg.adminToken)
	defer cancel()
	resp, err := client.DisableTable(ctx, &administrationv1.DisableTableRequest{
		Source: cfg.source, Schema: cfg.schema, Table: cfg.table, Publication: cfg.publication,
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("grpc-client: removed=%v retained=%v", resp.GetRemoved(), resp.GetRetained()), nil
}

// getTableStatus ruft die reader-RPC `GetTableStatus` auf (`LH-FA-CFG-003`).
func getTableStatus(client administrationv1.AdministrationClient, cfg config) (string, error) {
	ctx, cancel := callCtx(cfg.token)
	defer cancel()
	resp, err := client.GetTableStatus(ctx, &administrationv1.GetTableStatusRequest{
		Source: cfg.source, Schema: cfg.schema, Table: cfg.table, Publication: cfg.publication,
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("grpc-client: enabled=%v retained=%v", resp.GetEnabled(), resp.GetRetained()), nil
}

// listTables ruft die reader-RPC `ListTables` auf (`LH-FA-CFG-004`) und
// listet jede Tabelle der Antwort — `tables` und `retained` je eine
// `SourceTable`-Zeile.
func listTables(client administrationv1.AdministrationClient, cfg config) (string, error) {
	ctx, cancel := callCtx(cfg.token)
	defer cancel()
	resp, err := client.ListTables(ctx, &administrationv1.ListTablesRequest{Source: cfg.source, Publication: cfg.publication})
	if err != nil {
		return "", err
	}
	lines := []string{fmt.Sprintf("grpc-client: tables=%d retained=%d", len(resp.GetTables()), len(resp.GetRetained()))}
	for _, t := range resp.GetTables() {
		lines = append(lines, fmt.Sprintf("  table: table_id=%s source=%s schema=%s table=%s", t.GetTableId(), t.GetSource(), t.GetSchema(), t.GetTable()))
	}
	for _, t := range resp.GetRetained() {
		lines = append(lines, fmt.Sprintf("  retained: table_id=%s source=%s schema=%s table=%s", t.GetTableId(), t.GetSource(), t.GetSchema(), t.GetTable()))
	}
	return strings.Join(lines, "\n"), nil
}
