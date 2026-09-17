// Command grpc-client ist ein öffentliches Beispiel für den Live-Change-Stream
// über gRPC (`LH-FA-SST-008`, `ADR-0060`, `ADR-0076`): es öffnet den
// Server-Streaming-RPC `ChangeStream/StreamChanges` real gegen den laufenden
// Feed-Container und gibt jede empfangene Nachricht aus. Startform ist
// `go run ./examples/grpc-client`; der Zugriffs-Abschnitt des
// Benutzerhandbuchs ist `### Zugriff über den gRPC-Change-Stream`
// (`docs/user/benutzerhandbuch.md`).
//
// Dieses Programm ist das Vorbild (minimal, lesbar); der E2E-Belegträger zu
// LH-FA-SST-008 ist der Wegwerf-Client `tools/harness/grpcclient`. Das
// Beispiel trägt keine Zustandsmaschine: der Stream kennt kein Replay
// (`ADR-0060`), verpasste Changes holt der bestehende Lesezugriffsweg nach.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	streamv1 "github.com/pt9912/pg-change-feed/gen/cdc/stream/v1"
)

// authorizationMetadataKey und bearerPrefix tragen dieselbe Wertform wie der
// Auth-Interceptor des Adapters (`SPEC-020`): der `authorization`-Metadata-
// Wert trägt den Token hinter dem `Bearer `-Vorsprung.
const (
	authorizationMetadataKey = "authorization"
	bearerPrefix             = "Bearer "
)

// config trägt die Laufzeit-Eingabe des Beispiels: die Horch-Adresse des
// gRPC-Streaming-Servers und das Token der lesenden Rechtsklasse. Beide
// kommen aus denselben Umgebungsvariablen, die das Benutzerhandbuch führt
// (`CDC_GRPC_ADDR`, `CDC_API_TOKEN_READER`), und lassen sich per Flag
// übersteuern (`ADR-0076` Festlegung 1).
type config struct {
	addr  string
	token string
}

func main() {
	cfg := parseFlags()
	if cfg.addr == "" {
		fmt.Fprintln(os.Stderr, "grpc-client: keine gRPC-Adresse gesetzt — CDC_GRPC_ADDR (oder -addr) ist nötig, um den Stream zu öffnen")
		os.Exit(2)
	}
	if cfg.token == "" {
		fmt.Fprintln(os.Stderr, "grpc-client: kein Token gesetzt — CDC_API_TOKEN_READER (oder -token) ist nötig, um über die reader-Rechtsklasse zu lesen")
		os.Exit(2)
	}

	conn, err := grpc.NewClient(cfg.addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "grpc-client: Verbindung (%s) fehlgeschlagen: %v\n", cfg.addr, err)
		os.Exit(1)
	}
	defer func() { _ = conn.Close() }()

	client := streamv1.NewChangeStreamClient(conn)

	// Der Stream bleibt offen, bis die Verbindung endet: der Aufruf läuft
	// ohne eigene Frist über `context.Background()`.
	ctx := metadata.AppendToOutgoingContext(context.Background(), authorizationMetadataKey, bearerPrefix+cfg.token)
	stream, err := client.StreamChanges(ctx, &streamv1.StreamChangesRequest{})
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

// parseFlags liest die Flag-Werte und füllt fehlende Felder aus den im
// Handbuch dokumentierten Umgebungsvariablen (`ADR-0076` Festlegung 1).
func parseFlags() config {
	var cfg config
	flag.StringVar(&cfg.addr, "addr", os.Getenv("CDC_GRPC_ADDR"), "Horch-Adresse des gRPC-Streaming-Servers, host:port (Default: CDC_GRPC_ADDR)")
	flag.StringVar(&cfg.token, "token", os.Getenv("CDC_API_TOKEN_READER"), "Bearer-Token der lesenden Rechtsklasse (Default: CDC_API_TOKEN_READER)")
	flag.Parse()
	return cfg
}
