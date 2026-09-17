package grpc

import (
	"context"
	"errors"
	"net"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	"github.com/pt9912/pg-change-feed/gen/cdc/stream/v1"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

const (
	testReaderToken = "reader-token"
	testAdminToken  = "admin-token"
)

// fakeSubscriber trägt eine In-Memory-Fälschung des `changeSubscriber`
// (`ADR-0030`): Whitebox-Test des Adapters ohne den Driven-Broadcaster —
// der Adapter-Test darf kein Adapter-Paket importieren (Richtungskonvention
// `spec/architecture.md` §1, maschinell über `a-check`). `ready` schließt
// `Subscribe` einmalig und macht die Registrierung für den Test beobachtbar.
type fakeSubscriber struct {
	changes    chan *model.Change
	ready      chan struct{}
	once       sync.Once
	abgemeldet bool
}

func newFakeSubscriber() *fakeSubscriber {
	return &fakeSubscriber{changes: make(chan *model.Change), ready: make(chan struct{})}
}

func (f *fakeSubscriber) Subscribe() (<-chan *model.Change, func()) {
	f.once.Do(func() { close(f.ready) })
	return f.changes, func() { f.abgemeldet = true }
}

// startTestServer verdrahtet den Adapter auf einem In-Memory-Listener
// (`bufconn`) — kein realer Port, kein Netz (`AGENTS.md` §3.1 im Testlauf).
func startTestServer(t *testing.T, subscriber changeSubscriber) streamv1.ChangeStreamClient {
	t.Helper()
	return startTestServerMitTokenKonfiguration(t, subscriber, testReaderToken, testAdminToken)
}

// startTestServerMitTokenKonfiguration verdrahtet den Adapter mit den
// übergebenen Token-Klassen; ein leerer Wert trägt die ungesetzte
// Konfiguration (`CDC_API_TOKEN_READER`/`CDC_API_TOKEN_ADMIN`).
func startTestServerMitTokenKonfiguration(t *testing.T, subscriber changeSubscriber, readerToken, adminToken string) streamv1.ChangeStreamClient {
	t.Helper()
	srv := New(Config{
		Addr:        "unused:0",
		TokenReader: readerToken,
		TokenAdmin:  adminToken,
		Subscriber:  subscriber,
	})
	listener := bufconn.Listen(1024 * 1024)
	go func() { _ = srv.serve(listener) }()

	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			return listener.DialContext(ctx)
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Client-Verbindung bauen: %v", err)
	}
	t.Cleanup(func() {
		_ = conn.Close()
		srv.Shutdown()
	})
	return streamv1.NewChangeStreamClient(conn)
}

// streamMitToken öffnet den Stream mit dem übergebenen Metadata-Wert; ein
// leerer Wert lässt die Metadata weg.
func streamMitToken(t *testing.T, client streamv1.ChangeStreamClient, authorization string) grpc.ServerStreamingClient[streamv1.Change] {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	if authorization != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, authorizationMetadataKey, authorization)
	}
	stream, err := client.StreamChanges(ctx, &streamv1.StreamChangesRequest{})
	if err != nil {
		t.Fatalf("Stream öffnen: %v", err)
	}
	return stream
}

// recvChange liest einen Change mit Frist — ein hängender Stream ist ein
// Fehlschlag, kein Warten.
func recvChange(t *testing.T, stream grpc.ServerStreamingClient[streamv1.Change]) (*streamv1.Change, error) {
	t.Helper()
	type ergebnis struct {
		msg *streamv1.Change
		err error
	}
	out := make(chan ergebnis, 1)
	go func() {
		msg, err := stream.Recv()
		out <- ergebnis{msg, err}
	}()
	select {
	case r := <-out:
		return r.msg, r.err
	case <-time.After(3 * time.Second):
		t.Fatal("kein Stream-Ergebnis innerhalb der Frist")
		return nil, nil
	}
}

// TestStreamChangesOhneMetadataEndetMitUnauthenticated trägt die erste
// Hälfte der Fitness Function aus `ADR-0060` Teilfrage 4: ein
// Stream-Öffnungsversuch ohne `authorization`-Metadata endet mit dem
// gRPC-Status `Unauthenticated` (`LH-FA-SST-008` Negative).
func TestStreamChangesOhneMetadataEndetMitUnauthenticated(t *testing.T) {
	client := startTestServer(t, newFakeSubscriber())
	stream := streamMitToken(t, client, "")
	_, err := recvChange(t, stream)
	if status.Code(err) != codes.Unauthenticated {
		t.Fatalf("Status: %v (Erwartung: %v)", status.Code(err), codes.Unauthenticated)
	}
}

// TestStreamChangesUnbekannterWertEndetMitUnauthenticated trägt denselben
// Pfad für einen `authorization`-Wert, der keiner konfigurierten Klasse
// entspricht — nicht nur ein fehlender Eintrag.
func TestStreamChangesUnbekannterWertEndetMitUnauthenticated(t *testing.T) {
	client := startTestServer(t, newFakeSubscriber())
	stream := streamMitToken(t, client, bearerPrefix+"unbekannt")
	_, err := recvChange(t, stream)
	if status.Code(err) != codes.Unauthenticated {
		t.Fatalf("Status: %v (Erwartung: %v)", status.Code(err), codes.Unauthenticated)
	}
}

// TestStreamChangesFalscheWertformEndetMitUnauthenticated trägt die
// Grenze der Wertform: ein Wert ohne `Bearer `-Vorsprung trägt keinen
// Token (`SPEC-020`) und endet wie ein fehlender.
func TestStreamChangesFalscheWertformEndetMitUnauthenticated(t *testing.T) {
	client := startTestServer(t, newFakeSubscriber())
	stream := streamMitToken(t, client, testReaderToken)
	_, err := recvChange(t, stream)
	if status.Code(err) != codes.Unauthenticated {
		t.Fatalf("Status: %v (Erwartung: %v)", status.Code(err), codes.Unauthenticated)
	}
}

// TestStreamChangesReaderTokenOeffnetTraegtChange trägt die zweite Hälfte
// der Fitness Function: ein gültiges `reader`-Token öffnet den Stream, und
// ein vom Broadcaster gelieferter Change erreicht den Client mit
// vollständigem Inhalt (`LH-FA-SST-008` Happy Path, `SPEC-020`).
func TestStreamChangesReaderTokenOeffnetTraegtChange(t *testing.T) {
	traegtTokenOeffnetStreamUndTraegtChange(t, testReaderToken)
}

// TestStreamChangesAdminTokenOeffnetTraegtChange trägt denselben Pfad für
// die administrative Klasse — ein `admin`-Token deckt die lesende
// Streaming-Fähigkeit ab (`ADR-0060` Teilfrage 4).
func TestStreamChangesAdminTokenOeffnetTraegtChange(t *testing.T) {
	traegtTokenOeffnetStreamUndTraegtChange(t, testAdminToken)
}

func traegtTokenOeffnetStreamUndTraegtChange(t *testing.T, token string) {
	t.Helper()
	subscriber := newFakeSubscriber()
	client := startTestServer(t, subscriber)
	stream := streamMitToken(t, client, bearerPrefix+token)

	select {
	case <-subscriber.ready:
	case <-time.After(3 * time.Second):
		t.Fatal("der Stream hat sich nicht am Broadcaster registriert")
	}

	change, err := model.NewChange("change-1", "tx-1", "table-1", 2, model.OperationUpdate,
		[]byte(`{"id":1}`), []byte(`{"id":1,"bestellstatus":"bezahlt"}`), "table-1-v1")
	if err != nil {
		t.Fatalf("Change bauen: %v", err)
	}
	change.Schema = "public"
	change.Table = "orders"
	subscriber.changes <- &change

	msg, err := recvChange(t, stream)
	if err != nil {
		t.Fatalf("Recv: %v", err)
	}
	if msg.GetChangeId() != "change-1" || msg.GetTransactionId() != "tx-1" ||
		msg.GetSourceTableId() != "table-1" || msg.GetSequence() != 2 ||
		msg.GetOperation() != string(model.OperationUpdate) ||
		msg.GetSchema() != "public" || msg.GetTable() != "orders" ||
		msg.GetSchemaVersion() != "table-1-v1" {
		t.Fatalf("übertragener Change trägt nicht die erwarteten Felder: %+v", msg)
	}
	if string(msg.GetOldImage()) != `{"id":1}` || string(msg.GetNewImage()) != `{"id":1,"bestellstatus":"bezahlt"}` {
		t.Fatalf("übertragener Change trägt nicht die erwarteten Row Images: old=%q new=%q", msg.GetOldImage(), msg.GetNewImage())
	}
}

// TestStreamChangesOhneBroadcasterEndetMitInternal trägt den Fehler-Ausgang
// einer fehlenden Verdrahtung: ohne Broadcaster endet der RPC sichtbar,
// statt still einen leeren, nie endenden Stream zu liefern.
func TestStreamChangesOhneBroadcasterEndetMitInternal(t *testing.T) {
	client := startTestServer(t, nil)
	stream := streamMitToken(t, client, bearerPrefix+testReaderToken)
	_, err := recvChange(t, stream)
	if status.Code(err) != codes.Internal {
		t.Fatalf("Status: %v (Erwartung: %v)", status.Code(err), codes.Internal)
	}
}

// TestStartUndShutdown trägt den Lebenszyklus über eine reale
// Horch-Adresse: `Start` bindet und kehrt nach `Shutdown` ohne Fehler
// zurück (`grpc.ErrServerStopped` ist der reguläre Ausgang, keine
// Fehlerklasse).
func TestStartUndShutdown(t *testing.T) {
	srv := New(Config{Addr: "127.0.0.1:0", Subscriber: newFakeSubscriber()})
	fertig := make(chan error, 1)
	go func() { fertig <- srv.Start() }()

	srv.Shutdown()
	select {
	case err := <-fertig:
		if err != nil {
			t.Fatalf("Start/Shutdown: %v (Erwartung: kein Fehler)", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Start kehrt nach Shutdown nicht zurück")
	}
}

// TestStartLiefertFehlerBeiUngueltigerAdresse trägt den Bind-Fehlerpfad:
// eine unbrauchbare Horch-Adresse endet als Fehler von `Start`, nicht in
// einem stillen Lauf ohne Listener.
func TestStartLiefertFehlerBeiUngueltigerAdresse(t *testing.T) {
	srv := New(Config{Addr: "127.0.0.1:-1", Subscriber: newFakeSubscriber()})
	if err := srv.Start(); err == nil {
		t.Fatal("Start mit ungültiger Horch-Adresse liefert keinen Fehler")
	}
}

// fakeServerStream trägt einen `grpc.ServerStream`-Doppel: er liefert den
// übergebenen Kontext und die injizierten Sende-/Empfangsfehler, ohne dass
// ein Transport im Spiel ist.
type fakeServerStream struct {
	ctx     context.Context
	sendErr error
	recvErr error
}

func (f *fakeServerStream) SetHeader(metadata.MD) error  { return nil }
func (f *fakeServerStream) SendHeader(metadata.MD) error { return nil }
func (f *fakeServerStream) SetTrailer(metadata.MD)       {}
func (f *fakeServerStream) Context() context.Context     { return f.ctx }
func (f *fakeServerStream) SendMsg(any) error            { return nil }
func (f *fakeServerStream) RecvMsg(any) error            { return f.recvErr }

var _ grpc.ServerStream = (*fakeServerStream)(nil)

// fakeChangeServerStream erfüllt `grpc.ServerStreamingServer[streamv1.Change]`
// auf demselben Doppel.
type fakeChangeServerStream struct{ *fakeServerStream }

func (f fakeChangeServerStream) Send(*streamv1.Change) error { return f.sendErr }

var _ grpc.ServerStreamingServer[streamv1.Change] = fakeChangeServerStream{}

// haengeFrist trägt den Hänge-Schutz des Tests, dessen Server-Stream sein
// Ende nur über ein Ergebnis-Signal meldet: sie unterscheidet
// „hängengeblieben" von „fertig" und ist kein Urteil über eine Uhr
// (`BEO-PGC/test-integration-retention-timing-flake`). Ihre Größe ist gegen
// die Laufzeit abgesetzt — `go test -race -count=20` läuft über dieses Paket
// für alle zwanzig Iterationen zusammen in wenigen Sekunden durch, die Frist
// steht bei 30 s und greift deshalb nur bei einem realen Hänger.
const haengeFrist = 30 * time.Second

// TestServeMeldetListenerFehler trägt den Fehlerausgang von `serve`: ein
// Listener, der keine Verbindungen mehr annimmt, endet sichtbar als Fehler —
// nicht in einem stillen Lauf ohne Empfänger. `grpc.ErrServerStopped` bleibt
// davon getrennt: es trägt den regulären Ausgang eines `Shutdown`
// (`TestStartUndShutdown`), dieser Fall trägt den unerwarteten.
// Rot färbende Mutation: in `serve` den `err != nil`-Zweig zu `return nil`
// machen.
func TestServeMeldetListenerFehler(t *testing.T) {
	srv := New(Config{Addr: "unused:0", Subscriber: newFakeSubscriber()})
	listener := bufconn.Listen(1024)
	if err := listener.Close(); err != nil {
		t.Fatalf("Listener schließen: %v", err)
	}
	if err := srv.serve(listener); err == nil {
		t.Fatal("serve auf einem geschlossenen Listener liefert keinen Fehler")
	}
}

// TestStreamChangesSendefehlerWirdWeitergereicht trägt den Fehlerausgang des
// Streams: scheitert `Send`, endet der RPC mit **genau diesem** Fehler statt
// still weiterzulaufen (`LH-FA-SST-008` Negative). Die Ablehnung hängt am
// injizierten Sendefehler — der zurückgegebene Wert wird gegen ihn geprüft,
// nicht gegen „irgendein Fehler".
// Rot färbende Mutation: den `Send`-Fehler verwerfen und `nil` zurückgeben.
func TestStreamChangesSendefehlerWirdWeitergereicht(t *testing.T) {
	sendefehler := errors.New("transport weg")
	subscriber := newFakeSubscriber()
	service := &changeStreamService{subscriber: subscriber, log: outbound.NoopLog}
	stream := fakeChangeServerStream{&fakeServerStream{ctx: context.Background(), sendErr: sendefehler}}

	ergebnis := make(chan error, 1)
	go func() {
		ergebnis <- service.StreamChanges(&streamv1.StreamChangesRequest{}, stream)
	}()

	change, err := model.NewChange("change-1", "tx-1", "table-1", 1, model.OperationInsert, nil, nil, "table-1-v1")
	if err != nil {
		t.Fatalf("Change bauen: %v", err)
	}
	subscriber.changes <- &change

	select {
	case err := <-ergebnis:
		if !errors.Is(err, sendefehler) {
			t.Fatalf("StreamChanges: %v (Erwartung: der injizierte Sendefehler %v)", err, sendefehler)
		}
	case <-time.After(haengeFrist):
		t.Fatal("StreamChanges endete nach einem Sendefehler nicht")
	}
}
