package bootstrap

import (
	"context"
	"errors"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/receive"
	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// fakeCycleStream ist der Stream des Zyklus-Tests: er zählt Schließen und
// Lauf und hält den Kontext von `Run`.
type fakeCycleStream struct {
	bindCaptureErr error
	bindIdleErr    error
	runErr         error
	closes         int
	runs           int
	runCtx         context.Context
}

func (f *fakeCycleStream) Conn() *pgconn.PgConn { return nil }

func (f *fakeCycleStream) BindCapture(inbound.CaptureInboundPort) error { return f.bindCaptureErr }

func (f *fakeCycleStream) BindIdleConfirmation(inbound.IdleConfirmationInboundPort) error {
	return f.bindIdleErr
}

func (f *fakeCycleStream) Close(context.Context) { f.closes++ }

func (f *fakeCycleStream) Run(ctx context.Context) error {
	f.runs++
	f.runCtx = ctx
	return f.runErr
}

// fakeCycleService erfüllt beide Ports, die der Zyklus bindet.
type fakeCycleService struct{}

func (fakeCycleService) Capture(context.Context, inbound.CaptureCommand) (inbound.CaptureResult, error) {
	return inbound.CaptureResult{}, nil
}

func (fakeCycleService) ConfirmIdle(context.Context, inbound.IdleConfirmationCommand) (inbound.IdleConfirmationResult, error) {
	return inbound.IdleConfirmationResult{}, nil
}

// fakeCycleAck zeichnet die bestätigte Position auf.
type fakeCycleAck struct{ positions []model.SourcePosition }

func (a *fakeCycleAck) Acknowledge(_ context.Context, p model.SourcePosition) error {
	a.positions = append(a.positions, p)
	return nil
}

func okAck(ack outbound.ReplicationAckPort) func(*pgconn.PgConn) (outbound.ReplicationAckPort, error) {
	return func(*pgconn.PgConn) (outbound.ReplicationAckPort, error) { return ack, nil }
}

// TestRunStreamCycleClosesTheStreamOnEveryPathBeforeRun bindet das Schließen
// des aufgebauten Streams an jeden Fehlerpfad zwischen Aufbau und Lauf:
// ACK-Adapter, Capture-Bindung, Bindung der Leerlauf-Bestätigung. Der Lauf
// findet dann nicht statt.
func TestRunStreamCycleClosesTheStreamOnEveryPathBeforeRun(t *testing.T) {
	boom := errors.New("boom")
	cases := []struct {
		name   string
		stream *fakeCycleStream
		newAck func(*pgconn.PgConn) (outbound.ReplicationAckPort, error)
	}{
		{"ACK-Adapter", &fakeCycleStream{}, func(*pgconn.PgConn) (outbound.ReplicationAckPort, error) { return nil, boom }},
		{"BindCapture", &fakeCycleStream{bindCaptureErr: boom}, okAck(&fakeCycleAck{})},
		{"BindIdleConfirmation", &fakeCycleStream{bindIdleErr: boom}, okAck(&fakeCycleAck{})},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			open := func(context.Context) (cycleStream, error) { return c.stream, nil }
			err := runStreamCycle(context.Background(), time.Second, open, c.newAck, &streamCycleAck{}, fakeCycleService{})
			if !errors.Is(err, boom) {
				t.Fatalf("Fehler = %v, erwartet boom", err)
			}
			if c.stream.closes != 1 {
				t.Fatalf("Close-Aufrufe = %d, erwartet 1", c.stream.closes)
			}
			if c.stream.runs != 0 {
				t.Fatalf("Run-Aufrufe = %d, erwartet 0", c.stream.runs)
			}
		})
	}
}

// TestRunStreamCycleOpenErrorHasNoStreamToClose trägt den Aufbau-Fehler: er
// wird durchgereicht, es gibt nichts zu schließen und nichts zu binden.
func TestRunStreamCycleOpenErrorHasNoStreamToClose(t *testing.T) {
	boom := errors.New("boom")
	acked := false
	newAck := func(*pgconn.PgConn) (outbound.ReplicationAckPort, error) { acked = true; return &fakeCycleAck{}, nil }
	open := func(context.Context) (cycleStream, error) { return nil, boom }
	err := runStreamCycle(context.Background(), time.Second, open, newAck, &streamCycleAck{}, fakeCycleService{})
	if !errors.Is(err, boom) || acked {
		t.Fatalf("Fehler = %v, ACK-Adapter gebaut = %v", err, acked)
	}
}

// TestRunStreamCycleRunsAndBindsTheAck trägt den Erfolgspfad: der ACK-Adapter
// des Zyklus steht hinter der Weiterleitung, `Run` läuft, ein Fehler von
// `Run` kommt unverändert zurück, und der Zyklus schließt den Stream nicht
// selbst (`Run` schließt bei seiner Rückkehr).
func TestRunStreamCycleRunsAndBindsTheAck(t *testing.T) {
	runErr := errors.New("Run-Fehler")
	stream := &fakeCycleStream{runErr: runErr}
	ack := &fakeCycleAck{}
	port := &streamCycleAck{}
	open := func(context.Context) (cycleStream, error) { return stream, nil }
	err := runStreamCycle(context.Background(), time.Second, open, okAck(ack), port, fakeCycleService{})
	if !errors.Is(err, runErr) {
		t.Fatalf("Fehler = %v, erwartet Run-Fehler", err)
	}
	if stream.runs != 1 || stream.closes != 0 {
		t.Fatalf("Run = %d, Close = %d, erwartet 1 und 0", stream.runs, stream.closes)
	}
	if err := port.Acknowledge(context.Background(), model.SourcePosition{Offset: 7}); err != nil {
		t.Fatalf("Weiterleitung: %v", err)
	}
	if len(ack.positions) != 1 || ack.positions[0].Offset != 7 {
		t.Fatalf("Der ACK-Adapter des Zyklus erhielt %v", ack.positions)
	}
}

// TestRunStreamCycleSetupDeadlineBindsOnlyTheSetup bindet die Frist an den
// Aufbau: `open` erhält einen Kontext mit der Frist, `Run` den Kontext des
// Zyklus ohne Frist — der Aufbau-Kontext ist bei `Run` bereits beendet, der
// Lauf-Kontext nicht.
func TestRunStreamCycleSetupDeadlineBindsOnlyTheSetup(t *testing.T) {
	const timeout = 40 * time.Second
	var setupCtx context.Context
	stream := &fakeCycleStream{}
	open := func(ctx context.Context) (cycleStream, error) {
		setupCtx = ctx
		return stream, nil
	}
	before := time.Now()
	if err := runStreamCycle(context.Background(), timeout, open, okAck(&fakeCycleAck{}), &streamCycleAck{}, fakeCycleService{}); err != nil {
		t.Fatalf("Zyklus: %v", err)
	}
	deadline, ok := setupCtx.Deadline()
	if !ok {
		t.Fatal("der Aufbau-Kontext trägt keine Frist")
	}
	if got := deadline.Sub(before); got < timeout-time.Second || got > timeout+time.Second {
		t.Fatalf("Frist des Aufbaus = %v, erwartet %v", got, timeout)
	}
	if setupCtx.Err() == nil {
		t.Fatal("der Aufbau-Kontext ist nach dem Aufbau nicht beendet")
	}
	if _, ok := stream.runCtx.Deadline(); ok {
		t.Fatal("der Lauf-Kontext trägt die Frist des Aufbaus")
	}
	if stream.runCtx.Err() != nil {
		t.Fatalf("der Lauf-Kontext ist beendet: %v", stream.runCtx.Err())
	}
}

// TestRunStreamCycleSetupDeadlineAgainstSilentListener trägt die Frist am
// echten Aufbau: ein Listener nimmt TCP an und antwortet nie; der Zyklus
// endet nach der Frist mit `ErrReplication`, der Ausgang ist wiederholbar.
func TestRunStreamCycleSetupDeadlineAgainstSilentListener(t *testing.T) {
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
	open := func(ctx context.Context) (cycleStream, error) {
		stream, err := receive.NewStream(ctx, receive.Config{
			DSN:         "postgres://cdc:pw@" + listener.Addr().String() + "/db?sslmode=disable",
			Source:      "src-1",
			Publication: "pub_x",
			Slot:        "slot_x",
		})
		if err != nil {
			return nil, err
		}
		return stream, nil
	}
	const timeout = 300 * time.Millisecond
	started := time.Now()
	err = runStreamCycle(context.Background(), timeout, open, okAck(&fakeCycleAck{}), &streamCycleAck{}, fakeCycleService{})
	if elapsed := time.Since(started); elapsed > 10*timeout {
		t.Fatalf("Zyklus kehrt nach %v zurück, Frist %v", elapsed, timeout)
	}
	if !retryableStreamError(err) {
		t.Fatalf("Fristablauf ist nicht wiederholbar: %v", err)
	}
}
