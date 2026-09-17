// Command grpcclient ist ein Wegwerf-Testclient für den gRPC-Stream-E2E-Beleg
// (LH-FA-SST-008, ADR-0060): er öffnet den Server-Stream real gegen den
// laufenden Feed-Container, meldet die Bereitschaft über die Zeile "READY"
// auf stdout, empfängt danach einen committeten Change mit vollständigem
// Inhalt ("RECEIVED") und belegt abschließend, dass ein Öffnungsversuch ohne
// gültiges Token über gRPC-Status `Unauthenticated` abgelehnt wird
// ("REJECTED"). Träger ist tools/harness/run-integration-tests.sh — der
// Aufrufer liest die stdout-Zeilen dieses Prozesses über `docker logs`, nicht
// über einen Exit-Code allein, weil "READY" vor der auslösenden Change
// beobachtbar sein muss.
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

	"github.com/pt9912/pg-change-feed/gen/cdc/stream/v1"
)

// authorizationMetadataKey und bearerPrefix tragen dieselbe Wertform wie der
// Auth-Interceptor des Adapters (SPEC-020): der `authorization`-Metadata-Wert
// trägt den Token hinter dem `Bearer `-Vorsprung.
const (
	authorizationMetadataKey = "authorization"
	bearerPrefix             = "Bearer "
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: grpcclient <addr> <token>")
		os.Exit(2)
	}
	addr, token := os.Args[1], os.Args[2]

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "grpcclient: Verbindung (%s) fehlgeschlagen: %v\n", addr, err)
		os.Exit(1)
	}
	defer func() { _ = conn.Close() }()

	client := streamv1.NewChangeStreamClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	stream, err := client.StreamChanges(
		metadata.AppendToOutgoingContext(ctx, authorizationMetadataKey, bearerPrefix+token),
		&streamv1.StreamChangesRequest{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "grpcclient: StreamChanges (authentifiziert) fehlgeschlagen: %v\n", err)
		os.Exit(1)
	}

	// Die Bereitschafts-Zeile folgt dem Aufbau der Streaming-Verbindung; der
	// Server registriert den Empfänger am Broadcaster asynchron dazu.
	fmt.Println("READY")

	change, err := stream.Recv()
	if err != nil {
		fmt.Fprintf(os.Stderr, "grpcclient: kein Change auf dem Stream innerhalb der Frist: %v\n", err)
		fmt.Println("TIMEOUT")
		os.Exit(1)
	}
	fmt.Printf("RECEIVED change_id=%s table=%s operation=%s new_image=%s\n",
		change.GetChangeId(), change.GetTable(), change.GetOperation(), change.GetNewImage())

	if err := assertUnauthenticated(client); err != nil {
		fmt.Fprintf(os.Stderr, "grpcclient: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("REJECTED code=Unauthenticated")
}

// assertUnauthenticated öffnet einen zweiten Stream ohne
// `authorization`-Metadata und erwartet die Ablehnung über gRPC-Status
// `Unauthenticated` — sichtbar, nicht still mit leeren Daten fortgesetzt
// (LH-FA-SST-008 Negative, ADR-0060 Teilfrage 4). Der Status kann am Aufruf
// oder erst am ersten Empfang eintreffen, deshalb deckt die Prüfung beide
// Ausgänge.
func assertUnauthenticated(client streamv1.ChangeStreamClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	stream, err := client.StreamChanges(ctx, &streamv1.StreamChangesRequest{})
	if err != nil {
		return checkUnauthenticated(err)
	}
	if _, err := stream.Recv(); err != nil {
		return checkUnauthenticated(err)
	}
	return fmt.Errorf("Stream-Öffnungsversuch ohne Token wurde akzeptiert (weder der Aufruf noch der Empfang endete mit Unauthenticated)")
}

// checkUnauthenticated verlangt den gRPC-Status `Unauthenticated` und reicht
// jeden anderen Ausgang als Fehler weiter.
func checkUnauthenticated(err error) error {
	if status.Code(err) != codes.Unauthenticated {
		return fmt.Errorf("Stream-Öffnungsversuch ohne Token endete mit Status %s, wollen Unauthenticated: %v", status.Code(err), err)
	}
	return nil
}
