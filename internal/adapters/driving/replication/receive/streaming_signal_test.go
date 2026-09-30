package receive

import (
	"context"
	stderrors "errors"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgproto3"
)

// TestRunSignalsStreamingAfterStartReplication trägt die Meldung des
// Streaming-Beginns: `OnStreaming` läuft genau einmal, nachdem der Server
// `START_REPLICATION` bestätigt hat.
func TestRunSignalsStreamingAfterStartReplication(t *testing.T) {
	session := &fakeSession{messages: []pgproto3.BackendMessage{&pgproto3.CopyDone{}}}
	stream := newTestStream(session, noopCapture{})
	calls := 0
	startCallsAtSignal := 0
	stream.onStreaming = func() {
		calls++
		startCallsAtSignal = session.startCalls
	}
	if err := stream.Run(context.Background()); err != nil {
		t.Fatalf("Lauf: %v", err)
	}
	if calls != 1 {
		t.Fatalf("Streaming-Signale = %d, erwartet 1", calls)
	}
	if startCallsAtSignal != 1 {
		t.Fatalf("START_REPLICATION-Aufrufe beim Signal = %d, erwartet 1", startCallsAtSignal)
	}
}

// TestRunSignalsNothingWhenStartReplicationFails trägt die Grenze des
// Signals: ein abgewiesenes `START_REPLICATION` meldet keinen Streaming-Beginn.
func TestRunSignalsNothingWhenStartReplicationFails(t *testing.T) {
	session := &fakeSession{startErr: stderrors.New("55006 Slot aktiv")}
	stream := newTestStream(session, noopCapture{})
	calls := 0
	stream.onStreaming = func() { calls++ }
	if err := stream.Run(context.Background()); !stderrors.Is(err, ErrReplication) {
		t.Fatalf("Lauf: %v, erwartet ErrReplication", err)
	}
	if calls != 0 {
		t.Fatalf("Streaming-Signale = %d, erwartet 0", calls)
	}
}

// TestRunWithoutOnStreamingRuns trägt den ungesetzten Rückruf: der Lauf
// verlangt ihn nicht.
func TestRunWithoutOnStreamingRuns(t *testing.T) {
	session := &fakeSession{messages: []pgproto3.BackendMessage{&pgproto3.CopyDone{}}}
	stream := newTestStream(session, noopCapture{})
	if err := stream.Run(context.Background()); err != nil {
		t.Fatalf("Lauf: %v", err)
	}
	if session.startCalls != 1 {
		t.Fatalf("START_REPLICATION-Aufrufe = %d, erwartet 1", session.startCalls)
	}
}

// TestStreamCloseClosesTheSession trägt `Stream.Close`: ein Stream, dessen
// `Run` nicht mehr aufgerufen wird, schließt seine Sitzung genau einmal.
func TestStreamCloseClosesTheSession(t *testing.T) {
	session := &fakeSession{}
	stream := newTestStream(session, noopCapture{})
	stream.Close(context.Background())
	if session.closeCalls != 1 {
		t.Fatalf("Close-Aufrufe = %d, erwartet 1", session.closeCalls)
	}
}

// TestNewStreamSetupDeadlineEndsWithReplicationError trägt die Frist des
// Aufbaus gegen den Treiber: ein Listener nimmt TCP an und antwortet nie;
// `NewStream` kehrt nach Ablauf der Kontext-Frist mit `ErrReplication`
// zurück und ohne `ErrRejected` (ein Fehler ohne SQLSTATE).
func TestNewStreamSetupDeadlineEndsWithReplicationError(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Listener: %v", err)
	}
	var mu sync.Mutex
	var held []net.Conn
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			mu.Lock()
			held = append(held, conn)
			mu.Unlock()
		}
	}()
	t.Cleanup(func() {
		listener.Close()
		mu.Lock()
		defer mu.Unlock()
		for _, c := range held {
			c.Close()
		}
	})

	const frist = 300 * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), frist)
	defer cancel()
	started := time.Now()
	_, err = NewStream(ctx, Config{
		DSN:         "postgres://cdc:pw@" + listener.Addr().String() + "/db?sslmode=disable",
		Source:      testSource,
		Publication: "pub_x",
		Slot:        "slot_x",
	})
	elapsed := time.Since(started)
	if !stderrors.Is(err, ErrReplication) {
		t.Fatalf("NewStream: %v, erwartet ErrReplication", err)
	}
	if stderrors.Is(err, ErrRejected) || stderrors.Is(err, ErrPermission) {
		t.Fatalf("Fristablauf trägt ErrRejected/ErrPermission: %v", err)
	}
	if elapsed > 10*frist {
		t.Fatalf("NewStream kehrt nach %v zurück, Frist %v", elapsed, frist)
	}
}
