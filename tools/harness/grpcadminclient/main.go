// Command grpcadminclient ist ein Wegwerf-Testclient für den
// gRPC-Administration-E2E-Beleg (ADR-0130): er ruft real gegen den
// laufenden Feed-Container mindestens eine RPC je Token-Klasse auf
// (ListTables mit dem reader-Token, RegisterConsumer mit dem admin-Token),
// meldet jedes Ergebnis auf stdout und belegt abschließend die beiden
// Negative-Pfade — ein Aufruf ohne Token endet mit gRPC-Status
// `Unauthenticated`, ein reader-Token gegen die admin-RPC RegisterConsumer
// mit `PermissionDenied`. Träger ist tools/harness/run-integration-tests.sh
// — der Aufrufer liest die stdout-Zeilen dieses Prozesses über `docker
// logs`, analog zu tools/harness/grpcclient.
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	administrationv1 "github.com/pt9912/pg-change-feed/gen/cdc/administration/v1"
)

// authorizationMetadataKey und bearerPrefix tragen dieselbe Wertform wie
// der Auth-Interceptor des Adapters (SPEC-031).
const (
	authorizationMetadataKey = "authorization"
	bearerPrefix             = "Bearer "
)

func main() {
	if len(os.Args) != 7 {
		fmt.Fprintln(os.Stderr, "usage: grpcadminclient <addr> <reader-token> <admin-token> <consumer-id> <source> <publication>")
		os.Exit(2)
	}
	addr, readerToken, adminToken, consumerID, source, publication := os.Args[1], os.Args[2], os.Args[3], os.Args[4], os.Args[5], os.Args[6]

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "grpcadminclient: Verbindung (%s) fehlgeschlagen: %v\n", addr, err)
		os.Exit(1)
	}
	defer func() { _ = conn.Close() }()

	client := administrationv1.NewAdministrationClient(conn)

	if err := listTables(client, readerToken, source, publication); err != nil {
		fmt.Fprintf(os.Stderr, "grpcadminclient: %v\n", err)
		os.Exit(1)
	}
	if err := registerConsumer(client, adminToken, consumerID); err != nil {
		fmt.Fprintf(os.Stderr, "grpcadminclient: %v\n", err)
		os.Exit(1)
	}
	if err := assertRejected(client, "", codes.Unauthenticated); err != nil {
		fmt.Fprintf(os.Stderr, "grpcadminclient: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("REJECTED code=Unauthenticated")
	if err := assertRejected(client, readerToken, codes.PermissionDenied); err != nil {
		fmt.Fprintf(os.Stderr, "grpcadminclient: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("REJECTED code=PermissionDenied")
}

// callCtx trägt die Aufruf-Frist und, sofern ein Token übergeben ist, den
// `authorization`-Metadata-Eintrag in der `Bearer`-Wertform.
func callCtx(token string) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	if token == "" {
		return ctx, cancel
	}
	return metadata.AppendToOutgoingContext(ctx, authorizationMetadataKey, bearerPrefix+token), cancel
}

// listTables ruft die reader-RPC `ListTables` auf und meldet das Ergebnis
// auf stdout ("LISTED").
func listTables(client administrationv1.AdministrationClient, readerToken, source, publication string) error {
	ctx, cancel := callCtx(readerToken)
	defer cancel()
	resp, err := client.ListTables(ctx, &administrationv1.ListTablesRequest{Source: source, Publication: publication})
	if err != nil {
		return fmt.Errorf("ListTables (reader) fehlgeschlagen: %w", err)
	}
	fmt.Printf("LISTED tables=%d retained=%d\n", len(resp.GetTables()), len(resp.GetRetained()))
	return nil
}

// registerConsumer ruft die admin-RPC `RegisterConsumer` auf und meldet
// das Ergebnis auf stdout ("REGISTERED").
func registerConsumer(client administrationv1.AdministrationClient, adminToken, consumerID string) error {
	ctx, cancel := callCtx(adminToken)
	defer cancel()
	resp, err := client.RegisterConsumer(ctx, &administrationv1.RegisterConsumerRequest{ConsumerId: consumerID, Name: "grpcadminclient"})
	if err != nil {
		return fmt.Errorf("RegisterConsumer (admin) fehlgeschlagen: %w", err)
	}
	fmt.Printf("REGISTERED consumer_id=%s already_registered=%v\n", resp.GetConsumerId(), resp.GetAlreadyRegistered())
	return nil
}

// assertRejected ruft `RegisterConsumer` (eine admin-RPC) mit dem
// übergebenen Token auf (leer heißt: kein `authorization`-Metadata-Eintrag)
// und verlangt den übergebenen gRPC-Status.
func assertRejected(client administrationv1.AdministrationClient, token string, want codes.Code) error {
	ctx, cancel := callCtx(token)
	defer cancel()
	_, err := client.RegisterConsumer(ctx, &administrationv1.RegisterConsumerRequest{ConsumerId: "grpcadminclient-rejected", Name: "grpcadminclient"})
	if status.Code(err) != want {
		return fmt.Errorf("RegisterConsumer mit Token %q endete mit Status %s, wollen %s: %v", token, status.Code(err), want, err)
	}
	return nil
}
