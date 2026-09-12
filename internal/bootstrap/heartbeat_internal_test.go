package bootstrap

import (
	"context"
	stderrors "errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage"
	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/decode"
	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/mapper"
	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/receive"
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
// freigibt — Stellvertreter für einen langsamen Heartbeat-Schreib-Zug:
// der Timer-Zug darf die Persist-before-ACK-Ordnung (`LH-QA-REL-001.a`)
// nicht stören.
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

// Fault trägt hier nur die Interface-Erfüllung — dieser Test belegt den
// Timer-Zug (`runHeartbeat`/`Beat`), nicht den Fehlerzustand-Pfad
// (`reportFault`/`classifyRunError`, eigene Tests unten).
func (b *blockingHeartbeat) Fault(ctx context.Context, source model.SourceID, class model.ErrorClass) error {
	return nil
}

var _ outbound.HeartbeatPort = (*blockingHeartbeat)(nil)

// recordingHeartbeat trägt jeden Fault-Aufruf zur Prüfung von reportFault
// (`LH-FA-ADM-003`) — Beat bleibt hier ungenutzt.
type recordingHeartbeat struct {
	faultClass model.ErrorClass
	faultCalls int
}

func (r *recordingHeartbeat) Beat(ctx context.Context, source model.SourceID) error { return nil }

func (r *recordingHeartbeat) Fault(ctx context.Context, source model.SourceID, class model.ErrorClass) error {
	r.faultCalls++
	r.faultClass = class
	return nil
}

var _ outbound.HeartbeatPort = (*recordingHeartbeat)(nil)

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

// TestHeartbeatDoesNotBlockCapturePersistAck belegt: `runHeartbeat`
// läuft in einer eigenen Goroutine über einen eigenen Port — ein
// blockierender Heartbeat-Schreib-Zug hält den
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
	// hb.Beat blockiert in <-b.release. Die Capture-Persist-ACK-
	// Schleife läuft trotzdem unabhängig weiter.

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
	if err := tx.Commit(position, model.NewTimePoint(1)); err != nil {
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

// TestClassifyRunErrorMapsKnownSentinelsToADR0023Classes belegt die
// Übersetzung der Composition Root (`LH-FA-ADM-003`, `ADR-0023`): jeder
// bekannte Sentinel trägt die dokumentierte Klasse, ein
// unbekannter Fehler bleibt `internal`.
func TestClassifyRunErrorMapsKnownSentinelsToADR0023Classes(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want model.ErrorClass
	}{
		{"Verdrahtungs-Konfiguration", ErrConfiguration, model.ErrorClassConfiguration},
		{"Stream-Konfiguration", receive.ErrConfiguration, model.ErrorClassConfiguration},
		{"Aktivierungs-Konfiguration", postgresstorage.ErrActivationConfiguration, model.ErrorClassConfiguration},
		{"Dekodierfehler", decode.ErrSchema, model.ErrorClassSchema},
		{"TRUNCATE nicht unterstützt", mapper.ErrTruncateUnsupported, model.ErrorClassSchema},
		{"Replication-Stream-Störung", receive.ErrReplication, model.ErrorClassReplication},
		{"ACK fehlgeschlagen", outbound.ErrReplication, model.ErrorClassReplication},
		{"Change ohne Begin", mapper.ErrChangeWithoutBegin, model.ErrorClassReplication},
		{"ChangeStore-Persistenzfehler", outbound.ErrStorage, model.ErrorClassStorage},
		{"Heartbeat-Persistenzfehler", outbound.ErrHeartbeatStorage, model.ErrorClassStorage},
		{"unbekannter Fehler", stderrors.New("überraschung"), model.ErrorClassInternal},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := classifyRunError(tc.err); got != tc.want {
				t.Fatalf("classifyRunError(%v) = %q, wollen %q", tc.err, got, tc.want)
			}
		})
	}
	// Die Wrappung (fmt.Errorf("%w: …", sentinel, cause)) bleibt über
	// errors.Is auflösbar — derselbe Pfad, den jeder Adapter trägt
	// (`ADR-0023`).
	wrapped := fmt.Errorf("%w: technischer Grund", outbound.ErrStorage)
	if got := classifyRunError(wrapped); got != model.ErrorClassStorage {
		t.Fatalf("classifyRunError(gewrappter ErrStorage) = %q, wollen %q", got, model.ErrorClassStorage)
	}
}

// TestReportFaultWritesClassifiedFaultOnNonNilError belegt reportFault
// (`LH-FA-ADM-003`): ein Lauf-Fehler erreicht den
// Heartbeat-Port mit der klassifizierten Kategorie, ein regulärer
// Lauf-Abschluss (`nil`) schreibt nichts. Rot färbende Mutation: die
// nil-Prüfung entfernen — dann trägt jeder reguläre Lauf-Abschluss
// fälschlich einen Fehlerzustand.
func TestReportFaultWritesClassifiedFaultOnNonNilError(t *testing.T) {
	rec := &recordingHeartbeat{}
	runErr := error(decode.ErrSchema)
	reportFault(rec, "src-1", &runErr)
	if rec.faultCalls != 1 {
		t.Fatalf("Fault-Aufrufe = %d, wollen 1", rec.faultCalls)
	}
	if rec.faultClass != model.ErrorClassSchema {
		t.Fatalf("Fault-Klasse = %q, wollen %q", rec.faultClass, model.ErrorClassSchema)
	}

	recNil := &recordingHeartbeat{}
	var nilErr error
	reportFault(recNil, "src-1", &nilErr)
	if recNil.faultCalls != 0 {
		t.Fatalf("Fault-Aufrufe bei nil-Fehler = %d, wollen 0", recNil.faultCalls)
	}
}
