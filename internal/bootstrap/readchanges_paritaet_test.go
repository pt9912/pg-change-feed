package bootstrap

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	administrationv1 "github.com/pt9912/pg-change-feed/gen/cdc/administration/v1"
	grpcadapter "github.com/pt9912/pg-change-feed/internal/adapters/driving/grpc"
	httpadapter "github.com/pt9912/pg-change-feed/internal/adapters/driving/http"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/readchanges"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

const (
	parityReaderToken = "reader-token"
	parityAdminToken  = "admin-token"
)

// parityStore trägt den `ChangeStorePort` als Fälschung für die Parität der
// Lesewege: er filtert die vorgelegten Changes mit derselben Gleichheit auf
// Quelle, Schema, Tabelle und Ziel, die die Abfrage des realen Stores trägt
// (`queries.SelectChanges`), und zählt seine Aufrufe.
type parityStore struct {
	records []outbound.ChangeRecord
	calls   int
}

func (s *parityStore) PersistTransaction(context.Context, *model.ChangeTransaction) error {
	return nil
}

func (s *parityStore) ReadChanges(_ context.Context, query outbound.ChangeQuery) ([]outbound.ChangeRecord, error) {
	s.calls++
	if err := query.Validate(); err != nil {
		return nil, err
	}
	var out []outbound.ChangeRecord
	for _, record := range s.records {
		if query.Source != record.Position.SourceID {
			continue
		}
		if query.Schema != "" && record.Change.Schema != query.Schema {
			continue
		}
		if query.Table != "" && record.Change.Table != query.Table {
			continue
		}
		if query.Target != "" && string(record.Change.RouteTarget) != query.Target {
			continue
		}
		out = append(out, record)
	}
	return out, nil
}

func (s *parityStore) ReadRetentionCandidates(context.Context, model.SourceID, model.ChangeID, int) ([]outbound.RetentionCandidate, error) {
	return nil, nil
}

func (s *parityStore) DeleteChanges(context.Context, []model.ChangeID) error { return nil }

func parityRecord(t *testing.T, id string, offset uint64, schema, table string, target model.RouteTarget) outbound.ChangeRecord {
	t.Helper()
	position, err := model.NewSourcePosition("src-1", offset)
	if err != nil {
		t.Fatalf("NewSourcePosition: %v", err)
	}
	change, err := model.NewChange(
		model.ChangeID(id), "tx-1", "tbl-1", 1, model.OperationInsert, nil, []byte(`{"n":1}`), "sv-1",
	)
	if err != nil {
		t.Fatalf("NewChange: %v", err)
	}
	change.Schema = schema
	change.Table = table
	change.RouteTarget = target
	return outbound.ChangeRecord{Position: position, Change: change, CommittedAt: model.NewTimePoint(42)}
}

// freeLoopbackAddr liefert eine freie Loopback-Adresse für den gRPC-Server.
// Der Server bindet die Adresse selbst; ein Start-Fehler (etwa ein zwischen
// Wahl und Start belegter Port) kommt über den Kanal von `startGRPC` zurück.
func freeLoopbackAddr(t *testing.T) string {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Loopback-Adresse: %v", err)
	}
	addr := listener.Addr().String()
	if err := listener.Close(); err != nil {
		t.Fatalf("Loopback-Adresse freigeben: %v", err)
	}
	return addr
}

// startGRPC startet den Server und liefert den Kanal seines Start-Ergebnisses;
// ein Aufrufer, dessen Anfrage scheitert, liest daraus den benannten
// Start-Fehler statt einer Zeitüberschreitung.
func startGRPC(t *testing.T, server *grpcadapter.Server) <-chan error {
	t.Helper()
	errCh := make(chan error, 1)
	go func() { errCh <- server.Start() }()
	t.Cleanup(server.Shutdown)
	return errCh
}

// failWithStartError meldet einen gescheiterten Aufruf; ist der Server nicht
// gestartet, nennt die Meldung den Start-Fehler.
func failWithStartError(t *testing.T, startErr <-chan error, call error) {
	t.Helper()
	select {
	case err := <-startErr:
		t.Fatalf("gRPC-Server nicht gestartet: %v (Aufruf: %v)", err, call)
	default:
		t.Fatalf("Aufruf: %v", call)
	}
}

// TestReadChangesWegeLiefernFuerDieselbeEingabeDieselbenChanges trägt die
// Gleichwertigkeit der Zugriffswege für den Ziel-Filter (`ADR-0138`
// Festlegung 1): `GET /changes` und der RPC `ReadChanges` rufen
// über den echten Use Case denselben Store; für dieselbe Eingabe — Ziel,
// Konjunktion mit Schema/Tabelle, Ziel ohne Treffer, Ziel außerhalb des
// Alphabets (Großbuchstabe, 64 Zeichen, U+0000) — liefern beide dieselben
// Changes, und ein Ziel außerhalb des Alphabets erreicht den Store nie.
// Rot färbende Mutation: in `readChangesHandler`s Parser (HTTP) oder in
// `administrationService.ReadChanges` (gRPC) das Ziel nicht an die Abfrage
// reichen — der Weg liefert dann alle Changes, der Vergleich mit dem
// erwarteten Satz und mit dem anderen Weg schlägt fehl; die
// Alphabet-Prüfung im Use Case streichen — der Store-Zähler steigt; die
// Alphabet-Prüfung vor die Validierung des Lese-Kontrakts setzen — die
// Fehlerfälle mit ungültigem Ziel antworten leer statt mit 400/InvalidArgument.
func TestReadChangesWegeLiefernFuerDieselbeEingabeDieselbenChanges(t *testing.T) {
	store := &parityStore{records: []outbound.ChangeRecord{
		parityRecord(t, "c-eu-orders", 10, "public", "orders", "eu"),
		parityRecord(t, "c-eu-customers", 20, "public", "customers", "eu"),
		parityRecord(t, "c-us-orders", 30, "public", "orders", "us"),
		parityRecord(t, "c-roh", 40, "public", "orders", ""),
	}}
	useCase := readchanges.NewReadChangesService(store)

	httpServer := httptest.NewServer(httpadapter.New(httpadapter.Config{
		Addr: "unused:0", TokenReader: parityReaderToken, TokenAdmin: parityAdminToken, ReadChanges: useCase,
	}).Handler())
	t.Cleanup(httpServer.Close)

	grpcAddr := freeLoopbackAddr(t)
	grpcServer := grpcadapter.New(grpcadapter.Config{
		Addr: grpcAddr, TokenReader: parityReaderToken, TokenAdmin: parityAdminToken, ReadChanges: useCase,
	})
	startErr := startGRPC(t, grpcServer)
	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("gRPC-Client: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	client := administrationv1.NewAdministrationClient(conn)

	httpIDs := func(schema, table, target string) []string {
		t.Helper()
		values := url.Values{"source": {"src-1"}}
		if schema != "" {
			values.Set("schema", schema)
		}
		if table != "" {
			values.Set("table", table)
		}
		if target != "" {
			values.Set("target", target)
		}
		req, err := http.NewRequest(http.MethodGet, httpServer.URL+"/changes?"+values.Encode(), nil)
		if err != nil {
			t.Fatalf("Request bauen: %v", err)
		}
		req.Header.Set("Authorization", "Bearer "+parityReaderToken)
		resp, err := httpServer.Client().Do(req)
		if err != nil {
			t.Fatalf("GET /changes: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("GET /changes (target %q): Status %d, wollen 200", target, resp.StatusCode)
		}
		var body struct {
			Changes []struct {
				ChangeID string `json:"change_id"`
			} `json:"changes"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
			t.Fatalf("Antwort lesen: %v", err)
		}
		if body.Changes == nil {
			t.Fatalf("GET /changes (target %q): changes ist null, wollen eine gesetzte Liste", target)
		}
		ids := []string{}
		for _, change := range body.Changes {
			ids = append(ids, change.ChangeID)
		}
		return ids
	}
	grpcIDs := func(schema, table, target string) []string {
		t.Helper()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+parityReaderToken)
		resp, err := client.ReadChanges(ctx, &administrationv1.ReadChangesRequest{
			Source: "src-1", Schema: schema, Table: table, Target: target,
		}, grpc.WaitForReady(true))
		if err != nil {
			failWithStartError(t, startErr, err)
		}
		ids := []string{}
		for _, change := range resp.GetChanges() {
			ids = append(ids, change.GetChangeId())
		}
		return ids
	}

	cases := []struct {
		name                  string
		schema, table, target string
		want                  []string
		storeCalls            int
	}{
		{"ohne Ziel liefert alles", "", "", "", []string{"c-eu-orders", "c-eu-customers", "c-us-orders", "c-roh"}, 1},
		{"Ziel eu", "", "", "eu", []string{"c-eu-orders", "c-eu-customers"}, 1},
		{"Ziel us", "", "", "us", []string{"c-us-orders"}, 1},
		{"Konjunktion mit Tabelle", "", "orders", "eu", []string{"c-eu-orders"}, 1},
		{"Konjunktion mit Schema und Tabelle", "public", "customers", "eu", []string{"c-eu-customers"}, 1},
		{"Ziel ohne Treffer", "", "", "zz", []string{}, 1},
		{"Großbuchstabe", "", "", "EU", []string{}, 0},
		{"64 Zeichen", "", "", strings.Repeat("a", 64), []string{}, 0},
		{"U+0000", "", "", "eu\x00", []string{}, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			store.calls = 0
			gotHTTP := httpIDs(tc.schema, tc.table, tc.target)
			callsAfterHTTP := store.calls
			gotGRPC := grpcIDs(tc.schema, tc.table, tc.target)
			if !reflect.DeepEqual(gotHTTP, tc.want) {
				t.Fatalf("GET /changes = %v, wollen %v", gotHTTP, tc.want)
			}
			if !reflect.DeepEqual(gotGRPC, tc.want) {
				t.Fatalf("ReadChanges = %v, wollen %v", gotGRPC, tc.want)
			}
			if callsAfterHTTP != tc.storeCalls || store.calls != 2*tc.storeCalls {
				t.Fatalf("Store-Aufrufe = %d nach HTTP, %d nach beiden Wegen, wollen %d und %d",
					callsAfterHTTP, store.calls, tc.storeCalls, 2*tc.storeCalls)
			}
		})
	}

	// Ein Fehler des Lese-Kontrakts endet auf beiden Wegen mit dem Fehler des
	// Bestands (HTTP 400, gRPC InvalidArgument), auch bei einem Ziel außerhalb
	// des Alphabets; kein Weg antwortet dort „leer".
	errorCases := []struct {
		name   string
		target string
		query  url.Values
		req    *administrationv1.ReadChangesRequest
	}{
		{"Limit unter 1, gültiges Ziel", "eu", url.Values{"limit": {"-1"}}, &administrationv1.ReadChangesRequest{Limit: -1}},
		{"Limit unter 1, ungültiges Ziel", "EU", url.Values{"limit": {"-1"}}, &administrationv1.ReadChangesRequest{Limit: -1}},
		{"invertierter Bereich, gültiges Ziel", "eu", url.Values{"from": {"300"}, "to": {"100"}}, &administrationv1.ReadChangesRequest{From: 300, To: 100}},
		{"invertierter Bereich, ungültiges Ziel", "EU", url.Values{"from": {"300"}, "to": {"100"}}, &administrationv1.ReadChangesRequest{From: 300, To: 100}},
	}
	for _, tc := range errorCases {
		t.Run(tc.name, func(t *testing.T) {
			values := tc.query
			values.Set("source", "src-1")
			values.Set("target", tc.target)
			req, err := http.NewRequest(http.MethodGet, httpServer.URL+"/changes?"+values.Encode(), nil)
			if err != nil {
				t.Fatalf("Request bauen: %v", err)
			}
			req.Header.Set("Authorization", "Bearer "+parityReaderToken)
			resp, err := httpServer.Client().Do(req)
			if err != nil {
				t.Fatalf("GET /changes: %v", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusBadRequest {
				t.Fatalf("GET /changes: Status %d, wollen 400", resp.StatusCode)
			}

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+parityReaderToken)
			tc.req.Source = "src-1"
			tc.req.Target = tc.target
			_, err = client.ReadChanges(ctx, tc.req, grpc.WaitForReady(true))
			if status.Code(err) != codes.InvalidArgument {
				t.Fatalf("ReadChanges: Code %v (%v), wollen InvalidArgument", status.Code(err), err)
			}
		})
	}
}
