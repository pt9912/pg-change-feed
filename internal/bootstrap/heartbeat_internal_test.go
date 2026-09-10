package bootstrap

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/capture"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Whitebox-Test (`package bootstrap`, nicht `bootstrap_test`): `runHeartbeat`
// und `healthcheckVerdict` sind unexportierte Verdrahtungsdetails
// (`ADR-0026`), das §6-Risiko dieses Slice verlangt aber einen Beleg an
// genau diesem Mechanismus — der Test greift deshalb direkt zu, statt die
// Funktion für den Test zu exportieren.

// blockingHeartbeat blockiert jeden Beat-Aufruf, bis der Test ihn
// freigibt — Stellvertreter für einen langsamen Heartbeat-Schreib-Zug
// (slice-012 §6-Risiko: der Timer-Zug darf die Persist-before-ACK-Ordnung,
// `LH-QA-REL-001.a`, nicht stören).
type blockingHeartbeat struct {
	release chan struct{}
	calls   chan struct{}
}

func (b *blockingHeartbeat) Beat(ctx context.Context, source model.SourceID) error {
	select {
	case b.calls <- struct{}{}:
	default:
	}
	<-b.release
	return nil
}

var _ outbound.HeartbeatPort = (*blockingHeartbeat)(nil)

// loggingStore und loggingAck tragen die Ordnungs-Logik des
// Capture-Persist-ACK-Pfads für diesen Test — dieselbe Fake-Form wie
// `capture.service_test.go` (dort paketintern, hier erneut minimal
// gehalten, weil `bootstrap` keinen Export dieser Test-Helfer besitzt).
type loggingStore struct{ events *[]string }

func (f *loggingStore) PersistTransaction(ctx context.Context, transaction *model.ChangeTransaction) error {
	*f.events = append(*f.events, "persist:"+string(transaction.ID))
	return nil
}

func (f *loggingStore) ReadChanges(ctx context.Context, query outbound.ChangeQuery) ([]outbound.ChangeRecord, error) {
	return nil, nil
}

var _ outbound.ChangeStorePort = (*loggingStore)(nil)

type loggingAck struct{ events *[]string }

func (f *loggingAck) Acknowledge(ctx context.Context, position model.SourcePosition) error {
	*f.events = append(*f.events, "ack")
	return nil
}

var _ outbound.ReplicationAckPort = (*loggingAck)(nil)

// TestHeartbeatDoesNotBlockCapturePersistAck belegt die §6-Risiko-Zusage
// (slice-012): `runHeartbeat` läuft in einer eigenen Goroutine über einen
// eigenen Port — ein blockierender Heartbeat-Schreib-Zug hält den
// Capture-Persist-ACK-Pfad (`LH-QA-REL-001.a`) nicht auf. Rot färbende
// Mutation: `runHeartbeat` unter demselben Lock/derselben Goroutine wie
// die Capture-Persist-ACK-Schleife laufen lassen — dann blockiert
// `Capture` unten, bis `hb.release` schließt, und der Test läuft in den
// 1s-Timeout.
func TestHeartbeatDoesNotBlockCapturePersistAck(t *testing.T) {
	hb := &blockingHeartbeat{release: make(chan struct{}), calls: make(chan struct{}, 1)}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		runHeartbeat(ctx, hb, "src-1", time.Millisecond)
	}()

	select {
	case <-hb.calls:
	case <-time.After(time.Second):
		t.Fatal("Heartbeat-Timer hat nicht innerhalb 1s aufgerufen")
	}
	// hb.Beat blockiert jetzt in <-b.release. Die Capture-Persist-ACK-
	// Schleife läuft unabhängig weiter — genau das behauptet slice-012 §6.

	events := []string{}
	captureSvc := capture.NewCaptureService(&loggingStore{events: &events}, &loggingAck{events: &events})
	tx, err := model.NewOpenTransaction("t-1", "src-1")
	if err != nil {
		t.Fatalf("NewOpenTransaction: %v", err)
	}
	position, err := model.NewSourcePosition("src-1", 100)
	if err != nil {
		t.Fatalf("NewSourcePosition: %v", err)
	}
	if err := tx.Commit(position); err != nil {
		t.Fatalf("Commit: %v", err)
	}

	done := make(chan error, 1)
	go func() {
		_, err := captureSvc.Capture(context.Background(), capture.CaptureCommand{Transaction: tx})
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Capture: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Capture blockierte, während der Heartbeat-Schreib-Zug lief — Timer-Zug und Persist-before-ACK teilen sich eine kritische Sektion")
	}
	if got, want := events, []string{"persist:t-1", "ack"}; got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("Ereignisse = %v, wollen %v", got, want)
	}

	close(hb.release)
	cancel()
	wg.Wait()
}

// TestHealthcheckVerdictThreshold belegt die Schwelle des
// `--healthcheck`-Ausgangs (`healthcheckVerdict`): unterhalb von
// `heartbeatStaleAfter` healthy (0), ab der Schwelle unhealthy (1). Rot
// färbende Mutation: die Schwelle von `>` auf `>=` ändern — dann meldet
// der Grenzwert selbst bereits unhealthy und dieser Test schlägt fehl.
func TestHealthcheckVerdictThreshold(t *testing.T) {
	cases := []struct {
		name string
		age  time.Duration
		want int
	}{
		{"frisch", 0, 0},
		{"unter der Schwelle", heartbeatStaleAfter - time.Second, 0},
		{"an der Schwelle", heartbeatStaleAfter, 0},
		{"über der Schwelle", heartbeatStaleAfter + time.Second, 1},
	}
	for _, c := range cases {
		if got := healthcheckVerdict(c.age); got != c.want {
			t.Errorf("healthcheckVerdict(%s) = %d, wollen %d (%s)", c.age, got, c.want, c.name)
		}
	}
}
