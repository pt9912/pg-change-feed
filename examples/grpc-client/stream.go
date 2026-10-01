package main

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	streamv1 "github.com/pt9912/pg-change-feed/gen/cdc/stream/v1"
)

// streamRequest bildet die Anfrage des Streams aus der Konfiguration:
// `cfg.schema`/`cfg.table`/`cfg.target` tragen den optionalen, unabhängig
// setzbaren Filter (`ADR-0133`) — alle leer liefert jeden Change aller
// aktivierten Tabellen.
func streamRequest(cfg config) *streamv1.StreamChangesRequest {
	return &streamv1.StreamChangesRequest{Schema: cfg.schema, Table: cfg.table, Target: cfg.target}
}

// runStream öffnet den Server-Streaming-RPC `ChangeStream/StreamChanges`
// real gegen den laufenden Feed-Container und gibt jede empfangene Nachricht
// aus, bis die Verbindung endet.
func runStream(conn *grpc.ClientConn, cfg config) {
	client := streamv1.NewChangeStreamClient(conn)

	// Der Stream bleibt offen, bis die Verbindung endet: der Aufruf läuft
	// ohne eigene Frist über `context.Background()`.
	ctx := metadata.AppendToOutgoingContext(context.Background(), authorizationMetadataKey, bearerPrefix+cfg.token)
	stream, err := client.StreamChanges(ctx, streamRequest(cfg))
	if err != nil {
		fmt.Fprintf(os.Stderr, "grpc-client: StreamChanges fehlgeschlagen: %v\n", err)
		os.Exit(1)
	}

	for {
		change, err := stream.Recv()
		if err != nil {
			fmt.Fprintf(os.Stderr, "grpc-client: Stream endete: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(formatChange(change))
	}
}
