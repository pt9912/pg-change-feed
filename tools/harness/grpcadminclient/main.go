// Command grpcadminclient ist ein Wegwerf-Testclient für den
// gRPC-Administration-E2E-Beleg (ADR-0132): er ruft real gegen den laufenden
// Feed-Container mindestens eine RPC je Token-Klasse auf (ListTables,
// ReadChanges und Diagnose mit dem reader-Token, RegisterConsumer mit dem
// admin-Token), meldet jedes Ergebnis auf stdout und belegt abschließend die
// beiden Negative-Pfade — ein Aufruf ohne Token endet mit gRPC-Status
// `Unauthenticated`, ein reader-Token gegen die admin-RPC RegisterConsumer
// mit `PermissionDenied`. Träger ist tools/harness/run-integration-tests.sh
// — der Aufrufer liest die stdout-Zeilen dieses Prozesses über
// `docker logs`, analog zu tools/harness/grpcclient.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
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
	target := flag.String("target", "", "Zustellziel-Filter des RPC ReadChanges (leer = kein Filter)")
	flag.Parse()
	args := flag.Args()
	if len(args) != 11 {
		fmt.Fprintln(os.Stderr, "usage: grpcadminclient [-target <ziel>] <addr> <reader-token> <admin-token> <consumer-id> <source> <publication> <read-schema> <read-table> <read-from> <read-to> <read-limit>")
		os.Exit(2)
	}
	addr, readerToken, adminToken, consumerID, source, publication := args[0], args[1], args[2], args[3], args[4], args[5]
	readSchema, readTable, readFrom, readTo, readLimit := args[6], args[7], args[8], args[9], args[10]

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
	if err := readChanges(client, readerToken, source, readSchema, readTable, *target, readFrom, readTo, readLimit); err != nil {
		fmt.Fprintf(os.Stderr, "grpcadminclient: %v\n", err)
		os.Exit(1)
	}
	if err := diagnose(client, readerToken, source); err != nil {
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

// readChanges ruft die reader-RPC `ReadChanges` auf und prüft die Antwort
// inhaltlich (`ADR-0131`): eine gesetzte, nicht leere Changes-Liste (ein
// gesetztes `target` filtert nach Zustellziel, `ChangeRecord` trägt es nicht);
// je Eintrag eine nicht leere Kennung, ein bekannter Operationswert, eine
// Position ≥ 1 und die zum Filter passende Klartext-Identität. Ein leerer
// Bereichs-/Limit-Parameter trägt `0` (nicht gesetzt) — dieselbe Semantik
// wie bei `httpclient`s `GET /changes`.
func readChanges(client administrationv1.AdministrationClient, readerToken, source, schema, table, target, from, to, limit string) error {
	fromVal, err := parseUint64(from)
	if err != nil {
		return fmt.Errorf("from %q ist keine Ganzzahl: %w", from, err)
	}
	toVal, err := parseUint64(to)
	if err != nil {
		return fmt.Errorf("to %q ist keine Ganzzahl: %w", to, err)
	}
	limitVal, err := parseInt64(limit)
	if err != nil {
		return fmt.Errorf("limit %q ist keine Ganzzahl: %w", limit, err)
	}

	ctx, cancel := callCtx(readerToken)
	defer cancel()
	resp, err := client.ReadChanges(ctx, &administrationv1.ReadChangesRequest{
		Source: source, Schema: schema, Table: table, Target: target, From: fromVal, To: toVal, Limit: limitVal,
	})
	if err != nil {
		return fmt.Errorf("ReadChanges (reader) fehlgeschlagen: %w", err)
	}
	changes := resp.GetChanges()
	if len(changes) == 0 {
		return fmt.Errorf("leere Changes-Liste für den erwarteten Bestand (%s.%s, Bereich [%s,%s))", schema, table, from, to)
	}
	var lastPosition int64
	for i, change := range changes {
		if change.GetChangeId() == "" {
			return fmt.Errorf("Eintrag %d trägt keine change_id", i)
		}
		if schema != "" && change.GetSchema() != schema {
			return fmt.Errorf("Eintrag %d trägt schema=%q, der Filter verlangt %q", i, change.GetSchema(), schema)
		}
		if table != "" && change.GetTable() != table {
			return fmt.Errorf("Eintrag %d trägt table=%q, der Filter verlangt %q", i, change.GetTable(), table)
		}
		switch change.GetOperation() {
		case "INSERT", "UPDATE", "DELETE":
		default:
			return fmt.Errorf("Eintrag %d trägt die unbekannte Operation %q", i, change.GetOperation())
		}
		if change.GetCommitPosition() < 1 {
			return fmt.Errorf("Eintrag %d trägt die Commit-Position %d (Position 0 existiert nicht)", i, change.GetCommitPosition())
		}
		if i > 0 && change.GetCommitPosition() < lastPosition {
			return fmt.Errorf("Eintrag %d trägt die Commit-Position %d nach %d — die Ordnung ist absteigend", i, change.GetCommitPosition(), lastPosition)
		}
		switch change.GetOrigin() {
		case "wal", "backfill":
		default:
			return fmt.Errorf("Eintrag %d trägt die unbekannte Herkunft %q", i, change.GetOrigin())
		}
		lastPosition = change.GetCommitPosition()
		fmt.Printf("READ changes=%d table=%s schema=%s change_id=%s operation=%s commit_position=%d new_image=%s origin=%s\n",
			len(changes), change.GetTable(), change.GetSchema(), change.GetChangeId(), change.GetOperation(), change.GetCommitPosition(), string(change.GetNewImage()), change.GetOrigin())
	}
	return nil
}

// diagnose ruft die reader-RPC `Diagnose` auf und meldet auf stdout
// ("DIAGNOSED") den bekannten Betriebsstatus — ob ein Lebenszeichen bekannt
// ist, dessen Fehlerzustand (leer = Normalbetrieb) und den gelesenen
// CDC-Abstand. Der Aufrufer hält `heartbeat_known`/`heartbeat_error_class`
// gegen eine kontemporäre `docker exec … diagnose`-Ausgabe (Querabgleich,
// nicht bloße Feld-Präsenz).
func diagnose(client administrationv1.AdministrationClient, readerToken, source string) error {
	ctx, cancel := callCtx(readerToken)
	defer cancel()
	resp, err := client.Diagnose(ctx, &administrationv1.DiagnoseRequest{Source: source})
	if err != nil {
		return fmt.Errorf("Diagnose (reader) fehlgeschlagen: %w", err)
	}
	fmt.Printf("DIAGNOSED heartbeat_known=%v heartbeat_error_class=%s capture_lag=%f consumer_lags=%d backfill=%d\n",
		resp.GetHeartbeat().GetKnown(), resp.GetHeartbeat().GetErrorClass(), resp.GetCaptureLag(), len(resp.GetConsumerLags()), len(resp.GetBackfill()))
	return nil
}

// parseUint64 liest einen optionalen `uint64`-Wert; ein leerer String trägt
// `0` (nicht gesetzt, ADR-0131 Teilfrage 3).
func parseUint64(raw string) (uint64, error) {
	if raw == "" {
		return 0, nil
	}
	return strconv.ParseUint(raw, 10, 64)
}

// parseInt64 liest einen optionalen `int64`-Wert; ein leerer String trägt
// `0` (unbegrenzt, ADR-0131 Teilfrage 3).
func parseInt64(raw string) (int64, error) {
	if raw == "" {
		return 0, nil
	}
	return strconv.ParseInt(raw, 10, 64)
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
