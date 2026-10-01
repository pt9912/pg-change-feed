// Diese Tests prüfen den erzeugten Protokoll-Stub an den Stellen, die eine
// Zusage des Draht-Vertrags tragen (`SPEC-020`, `ADR-0060` Teilfrage 1): die
// Feld-Getter, die Fehler-Weitergabe des erzeugten Clients und die
// Handler-Verklebung des Servers. Der erzeugte Code ist nicht von Hand
// änderbar — `make generated-sync` hält ihn byte-gleich zur Ausgabe des
// gepinnten Generators; geprüft wird deshalb sein Verhalten, nicht sein Text.
//
// Nicht Gegenstand sind die Proto-Runtime-Interna (`String`, `ProtoMessage`,
// `ProtoReflect`, `Descriptor`, die Deskriptor-Kompression und der
// Wiedereintritt von `init`). Sie sind aufrufbar, tragen aber keine Aussage
// über den Change-Stream: sie beschreiben den erzeugten Code gegen sich
// selbst (`ADR-0082` §Konsequenzen).
package streamv1_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/pt9912/pg-change-feed/gen/cdc/stream/v1"
)

// fakeClientConn trägt einen `grpc.ClientConnInterface`-Doppel: er liefert
// den injizierten Stream und den injizierten `NewStream`-Fehler, ohne dass
// ein Transport im Spiel ist.
type fakeClientConn struct {
	stream       grpc.ClientStream
	newStreamErr error
}

func (f fakeClientConn) Invoke(context.Context, string, any, any, ...grpc.CallOption) error {
	return nil
}

func (f fakeClientConn) NewStream(context.Context, *grpc.StreamDesc, string, ...grpc.CallOption) (grpc.ClientStream, error) {
	return f.stream, f.newStreamErr
}

var _ grpc.ClientConnInterface = fakeClientConn{}

// fakeClientStream trägt einen `grpc.ClientStream`-Doppel mit injizierten
// Sende- und Schließfehlern.
type fakeClientStream struct {
	sendErr  error
	closeErr error
}

func (fakeClientStream) Header() (metadata.MD, error) { return nil, nil }
func (fakeClientStream) Trailer() metadata.MD         { return nil }
func (fakeClientStream) Context() context.Context     { return context.Background() }
func (fakeClientStream) RecvMsg(any) error            { return nil }
func (f fakeClientStream) CloseSend() error           { return f.closeErr }
func (f fakeClientStream) SendMsg(any) error          { return f.sendErr }

var _ grpc.ClientStream = fakeClientStream{}

// fakeServerStream trägt einen `grpc.ServerStream`-Doppel mit injiziertem
// Empfangsfehler.
type fakeServerStream struct{ recvErr error }

func (fakeServerStream) SetHeader(metadata.MD) error  { return nil }
func (fakeServerStream) SendHeader(metadata.MD) error { return nil }
func (fakeServerStream) SetTrailer(metadata.MD)       {}
func (fakeServerStream) Context() context.Context     { return context.Background() }
func (fakeServerStream) SendMsg(any) error            { return nil }
func (f fakeServerStream) RecvMsg(any) error          { return f.recvErr }

var _ grpc.ServerStream = fakeServerStream{}

// TestChangeGetterTragenDenNullwertNilSicher trägt die Getter-Grenze des
// erzeugten Nachrichtentyps: auf dem Nullwert-Zeiger antworten alle zehn
// Getter mit ihrem Nullwert, statt zu dereferenzieren. Der Nullwert ist die
// Form, in der ein nicht gesetztes Feld über den Draht ankommt.
// Rot färbende Mutation: in einem Getter den `x != nil`-Zweig fallenlassen
// (Dereferenzierung des Nullwerts) oder einen anderen Nullwert zurückgeben.
func TestChangeGetterTragenDenNullwertNilSicher(t *testing.T) {
	var change *streamv1.Change
	if got := change.GetChangeId(); got != "" {
		t.Fatalf("GetChangeId auf dem Nullwert = %q (Erwartung: leer)", got)
	}
	if got := change.GetTransactionId(); got != "" {
		t.Fatalf("GetTransactionId auf dem Nullwert = %q (Erwartung: leer)", got)
	}
	if got := change.GetSourceTableId(); got != "" {
		t.Fatalf("GetSourceTableId auf dem Nullwert = %q (Erwartung: leer)", got)
	}
	if got := change.GetSequence(); got != 0 {
		t.Fatalf("GetSequence auf dem Nullwert = %d (Erwartung: 0)", got)
	}
	if got := change.GetOperation(); got != "" {
		t.Fatalf("GetOperation auf dem Nullwert = %q (Erwartung: leer)", got)
	}
	if got := change.GetOldImage(); got != nil {
		t.Fatalf("GetOldImage auf dem Nullwert = %q (Erwartung: nil)", got)
	}
	if got := change.GetNewImage(); got != nil {
		t.Fatalf("GetNewImage auf dem Nullwert = %q (Erwartung: nil)", got)
	}
	if got := change.GetSchemaVersion(); got != "" {
		t.Fatalf("GetSchemaVersion auf dem Nullwert = %q (Erwartung: leer)", got)
	}
	if got := change.GetSchema(); got != "" {
		t.Fatalf("GetSchema auf dem Nullwert = %q (Erwartung: leer)", got)
	}
	if got := change.GetTable(); got != "" {
		t.Fatalf("GetTable auf dem Nullwert = %q (Erwartung: leer)", got)
	}

	// Gegenprobe an derselben Aufrufstelle: ein gesetzter Change trägt seine
	// Werte — der Nullwert oben ist damit das Ergebnis der Eingabe, nicht des
	// Getter-Aufrufs.
	gesetzt := &streamv1.Change{ChangeId: "change-1", Sequence: 7, Schema: "public", Table: "orders"}
	if gesetzt.GetChangeId() != "change-1" || gesetzt.GetSequence() != 7 ||
		gesetzt.GetSchema() != "public" || gesetzt.GetTable() != "orders" {
		t.Fatalf("gesetzter Change trägt nicht seine Werte: %+v", gesetzt)
	}
}

// TestStreamChangesRequestTraegtTargetAlsFeldDrei trägt den Draht-Vertrag
// des Feldes `target` (`SPEC-020`, Feldnummer 3): der Getter liest den
// gesetzten Wert und auf dem Nullwert leer; auf dem Draht steht der Wert
// unter dem Tag von Feld 3 (`0x1a`, Wire-Typ 2), neben `schema` (Feld 1)
// und `table` (Feld 2); Bytes eines Clients ohne das Feld (nur Feld 1 und 2)
// lesen mit leerem `target` — der Filter bleibt aus.
// Rot färbende Mutation: in der `.proto` die Feldnummer von `target` ändern
// (dann ändert sich der erwartete Tag, die Draht-Bytes und der Getter
// weichen von den festen Bytes ab).
func TestStreamChangesRequestTraegtTargetAlsFeldDrei(t *testing.T) {
	var nullwert *streamv1.StreamChangesRequest
	if got := nullwert.GetTarget(); got != "" {
		t.Fatalf("GetTarget auf dem Nullwert = %q (Erwartung: leer)", got)
	}
	gesetzt := &streamv1.StreamChangesRequest{Schema: "s", Table: "t", Target: "eu"}
	if got := gesetzt.GetTarget(); got != "eu" {
		t.Fatalf("GetTarget = %q (Erwartung: eu)", got)
	}

	draht, err := proto.Marshal(gesetzt)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	wollen := []byte{0x0a, 0x01, 's', 0x12, 0x01, 't', 0x1a, 0x02, 'e', 'u'}
	if !bytes.Equal(draht, wollen) {
		t.Fatalf("Draht-Bytes = %x (Erwartung: %x)", draht, wollen)
	}

	ohneFeld := &streamv1.StreamChangesRequest{}
	if err := proto.Unmarshal([]byte{0x0a, 0x01, 's', 0x12, 0x01, 't'}, ohneFeld); err != nil {
		t.Fatalf("Unmarshal ohne Feld 3: %v", err)
	}
	if ohneFeld.GetSchema() != "s" || ohneFeld.GetTable() != "t" || ohneFeld.GetTarget() != "" {
		t.Fatalf("Request ohne Feld 3 = %+v (Erwartung: schema/table gesetzt, target leer)", ohneFeld)
	}
}

// TestStreamChangesRequestIgnoriertUnbekannteFelder trägt, am Laufzeit-
// Verhalten der eingesetzten Protobuf-Bibliothek gemessen, die
// Voraussetzung der additiven Erweiterung: Bytes mit einem Feld, das der
// Empfänger nicht kennt (hier Feld 9), lesen ohne Fehler, die bekannten
// Felder bleiben erhalten. Das ist die Messung der Bibliothek, nicht eines
// ausgelieferten Altservers.
func TestStreamChangesRequestIgnoriertUnbekannteFelder(t *testing.T) {
	request := &streamv1.StreamChangesRequest{}
	if err := proto.Unmarshal([]byte{0x0a, 0x01, 's', 0x4a, 0x01, 'x'}, request); err != nil {
		t.Fatalf("Unmarshal mit unbekanntem Feld: %v", err)
	}
	if request.GetSchema() != "s" || request.GetTarget() != "" {
		t.Fatalf("Request = %+v (Erwartung: schema s, target leer)", request)
	}
}

// TestStreamChangesClientReichtFehlerWeiter trägt die Fehler-Weitergabe des
// erzeugten Clients (`LH-FA-SST-008` Negative): scheitert der
// Verbindungsaufbau, das Senden der Anfrage oder das Schließen der
// Sendeseite, liefert `StreamChanges` **genau diesen** Fehler und keinen
// Stream — ein Öffnungsversuch endet sichtbar, nicht in einem stillen,
// leeren Stream.
// Rot färbende Mutation: einen der drei Fehlerzweige auf `return x, nil`
// legen oder den Fehler durch einen anderen ersetzen.
func TestStreamChangesClientReichtFehlerWeiter(t *testing.T) {
	newStreamFehler := errors.New("kein Transport")
	sendeFehler := errors.New("Senden scheiterte")
	schliessFehler := errors.New("Schliessen scheiterte")

	cases := []struct {
		name string
		conn grpc.ClientConnInterface
		want error
	}{
		{"NewStream scheitert", fakeClientConn{newStreamErr: newStreamFehler}, newStreamFehler},
		{"SendMsg scheitert", fakeClientConn{stream: fakeClientStream{sendErr: sendeFehler}}, sendeFehler},
		{"CloseSend scheitert", fakeClientConn{stream: fakeClientStream{closeErr: schliessFehler}}, schliessFehler},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			client := streamv1.NewChangeStreamClient(tc.conn)
			stream, err := client.StreamChanges(context.Background(), &streamv1.StreamChangesRequest{})
			if !errors.Is(err, tc.want) {
				t.Fatalf("StreamChanges: %v (Erwartung: der injizierte Fehler %v)", err, tc.want)
			}
			if stream != nil {
				t.Fatalf("StreamChanges liefert zu einem Fehler einen Stream: %v", stream)
			}
		})
	}
}

// TestUnimplementedServerEndetMitUnimplemented trägt den Ausgang des
// Vorwärtskompatibilitäts-Stubs (`SPEC-020`): eine Implementierung, die
// `StreamChanges` nicht trägt, endet mit dem gRPC-Status `Unimplemented` —
// sichtbar, nicht mit einem stillen leeren Stream.
// Rot färbende Mutation: den Statuscode des Stubs ändern.
func TestUnimplementedServerEndetMitUnimplemented(t *testing.T) {
	unimplemented := streamv1.UnimplementedChangeStreamServer{}
	err := unimplemented.StreamChanges(&streamv1.StreamChangesRequest{}, nil)
	if status.Code(err) != codes.Unimplemented {
		t.Fatalf("Status: %v (Erwartung: %v)", status.Code(err), codes.Unimplemented)
	}
}

// TestStreamChangesHandlerReichtEmpfangsfehlerWeiter trägt die
// Handler-Verklebung des erzeugten Servers: scheitert `RecvMsg`, endet der
// RPC mit **genau diesem** Fehler, statt die Service-Methode mit einem
// unvollständigen Request aufzurufen. Der Handler wird über den
// veröffentlichten Service-Deskriptor angesprochen — derselbe Wert, den
// `RegisterChangeStreamServer` registriert.
// Rot färbende Mutation: den Fehler verwerfen und die Service-Methode
// dennoch aufrufen.
func TestStreamChangesHandlerReichtEmpfangsfehlerWeiter(t *testing.T) {
	empfangsfehler := errors.New("Nachricht nicht lesbar")
	handler := streamv1.ChangeStream_ServiceDesc.Streams[0].Handler
	err := handler(nil, fakeServerStream{recvErr: empfangsfehler})
	if !errors.Is(err, empfangsfehler) {
		t.Fatalf("Handler: %v (Erwartung: der injizierte Empfangsfehler %v)", err, empfangsfehler)
	}
}
