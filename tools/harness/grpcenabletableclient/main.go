// Command grpcenabletableclient ist ein Wegwerf-Testclient für den direkten
// gRPC-Zugriffsweg auf EnableTable/DisableTable (`LH-FA-CFG-001`): er ruft
// real gegen den laufenden Feed-Container genau eine der beiden RPCs mit dem
// Admin-Token auf, ohne die SQL-Antragsqueue (`cdc.enable_table`/
// `cdc.disable_table`) zu durchlaufen — der Beleg dafür, dass dieser Weg den
// laufenden Prozess ohne Neustart aktualisiert. Träger ist
// tools/harness/run-integration-tests.sh — der Aufrufer liest die
// stdout-Zeilen dieses Prozesses über `docker logs`, analog zu
// tools/harness/grpcadminclient.
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	administrationv1 "github.com/pt9912/pg-change-feed/gen/cdc/administration/v1"
)

// authorizationMetadataKey und bearerPrefix tragen dieselbe Wertform wie
// der Auth-Interceptor des Adapters (`SPEC-031`).
const (
	authorizationMetadataKey = "authorization"
	bearerPrefix             = "Bearer "
)

func main() {
	if len(os.Args) < 3 {
		usage()
	}
	addr, mode := os.Args[1], os.Args[2]

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "grpcenabletableclient: Verbindung (%s) fehlgeschlagen: %v\n", addr, err)
		os.Exit(1)
	}
	defer func() { _ = conn.Close() }()
	client := administrationv1.NewAdministrationClient(conn)

	switch mode {
	case "enable-table":
		runEnable(client, os.Args[3:])
	case "disable-table":
		runDisable(client, os.Args[3:])
	default:
		usage()
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: grpcenabletableclient <addr> enable-table <admin-token> <source> <schema> <table> <table-id> <schema-version-id> <publication>")
	fmt.Fprintln(os.Stderr, "   or: grpcenabletableclient <addr> disable-table <admin-token> <source> <schema> <table> <publication>")
	os.Exit(2)
}

// callCtx trägt die Aufruf-Frist und den `authorization`-Metadata-Eintrag in
// der `Bearer`-Wertform.
func callCtx(token string) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	return metadata.AppendToOutgoingContext(ctx, authorizationMetadataKey, bearerPrefix+token), cancel
}

// runEnable ruft die admin-RPC `EnableTable` real gegen den laufenden
// Feed-Container auf — derselbe Inbound Use Case wie der SQL-Antragsqueue-Pfad,
// aber ohne die Queue zu durchlaufen.
func runEnable(client administrationv1.AdministrationClient, args []string) {
	if len(args) != 7 {
		usage()
	}
	adminToken, source, schema, table, tableID, schemaVersionID, publication := args[0], args[1], args[2], args[3], args[4], args[5], args[6]

	ctx, cancel := callCtx(adminToken)
	defer cancel()
	resp, err := client.EnableTable(ctx, &administrationv1.EnableTableRequest{
		Source: source, Schema: schema, Table: table,
		TableId: tableID, SchemaVersionId: schemaVersionID, Version: 1, Publication: publication,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "grpcenabletableclient: EnableTable (admin) fehlgeschlagen: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("ENABLED table_id=%s already_enabled=%v\n", resp.GetTableId(), resp.GetAlreadyEnabled())
}

// runDisable ruft die admin-RPC `DisableTable` real gegen den laufenden
// Feed-Container auf — das Gegenstück zu `runEnable`.
func runDisable(client administrationv1.AdministrationClient, args []string) {
	if len(args) != 5 {
		usage()
	}
	adminToken, source, schema, table, publication := args[0], args[1], args[2], args[3], args[4]

	ctx, cancel := callCtx(adminToken)
	defer cancel()
	resp, err := client.DisableTable(ctx, &administrationv1.DisableTableRequest{
		Source: source, Schema: schema, Table: table, Publication: publication,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "grpcenabletableclient: DisableTable (admin) fehlgeschlagen: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("DISABLED removed=%v retained=%v\n", resp.GetRemoved(), resp.GetRetained())
}
