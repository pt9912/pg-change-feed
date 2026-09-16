// Die Tests dieses Pakets laufen netzlos: sie fahren die Verklebung des
// Adapters gegen einen Fake, der die Naht erfüllt — den abgesetzten
// Standby-Status-Update, die Null-Positions-Grenze und das
// Fehlerklassen-Wrapping. Sie sind ausdrücklich **kein** Ersatz der
// realen Tests: dass der Treiber die Meldung als confirmed_flush_lsn
// trägt, prüft weiterhin der PostgreSQL der Adapter-Tests
// (`make test-replication`, `ADR-0030`).

package postgresack

import (
	"context"
	stderrors "errors"
	"strings"
	"testing"

	"github.com/jackc/pglogrepl"

	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// fakeSender erfüllt dieselbe Naht wie die Treiber-Hülle und hält fest,
// was der Adapter ihm übergibt — die Verklebung, nicht das Protokoll.
type fakeSender struct {
	sent  []pglogrepl.StandbyStatusUpdate
	calls int
	err   error
}

func (f *fakeSender) SendStandbyStatusUpdate(_ context.Context, ssu pglogrepl.StandbyStatusUpdate) error {
	f.calls++
	if f.err != nil {
		return f.err
	}
	f.sent = append(f.sent, ssu)
	return nil
}

var _ standbySender = (*fakeSender)(nil)

// recordingLog trägt den `LogPort` der Tests (`LH-QA-OPS-004`,
// `ADR-0024`): er hält die Fehler-Aufrufe fest, ohne sie auszugeben.
type recordingLog struct {
	errors []string
}

func (l *recordingLog) Debug(context.Context, string, ...any) {}
func (l *recordingLog) Info(context.Context, string, ...any)  {}
func (l *recordingLog) Warn(context.Context, string, ...any)  {}
func (l *recordingLog) Error(_ context.Context, msg string, _ ...any) {
	l.errors = append(l.errors, msg)
}

var _ outbound.LogPort = (*recordingLog)(nil)

// TestStandbyStatusCarriesPositionInAllThreeLSNs trägt die
// Standby-Status-Form: die Meldung trägt die Position als Write-, Flush-
// und Apply-Stand, damit der Slot sie als confirmed_flush_lsn trägt
// (`LH-QA-REL-001.a`).
func TestStandbyStatusCarriesPositionInAllThreeLSNs(t *testing.T) {
	status := standbyStatus(pglogrepl.LSN(100))
	fields := map[string]pglogrepl.LSN{
		"WALWritePosition": status.WALWritePosition,
		"WALFlushPosition": status.WALFlushPosition,
		"WALApplyPosition": status.WALApplyPosition,
	}
	for name, got := range fields {
		if got != pglogrepl.LSN(100) {
			t.Fatalf("%s: %d, erwartet 100", name, got)
		}
	}
}

// TestAckLSNRejectsZeroPosition trägt die Null-Positions-Grenze: ohne
// Position gibt es nichts zu bestätigen (`SPEC-008`, Klasse
// `replication`).
func TestAckLSNRejectsZeroPosition(t *testing.T) {
	if _, err := ackLSN(model.SourcePosition{}); !stderrors.Is(err, outbound.ErrReplication) {
		t.Fatalf("Null-Position: %v", err)
	}
}

// TestAckLSNCarriesOffset trägt die LSN-Form (`ADR-0005`): der Offset
// wandert unverändert in die LSN.
func TestAckLSNCarriesOffset(t *testing.T) {
	position, err := model.NewSourcePosition("src-1", 4711)
	if err != nil {
		t.Fatalf("Position: %v", err)
	}
	lsn, err := ackLSN(position)
	if err != nil {
		t.Fatalf("ackLSN: %v", err)
	}
	if lsn != pglogrepl.LSN(4711) {
		t.Fatalf("LSN: %d, erwartet 4711", lsn)
	}
}

// TestReplicationClassWrapsCause trägt das Fehlerklassen-Wrapping am
// Treiber-Fehler: er endet über `outbound.ErrReplication`, die technische
// Ursache bleibt über die zweite Wrappung lesbar (`ADR-0023`,
// `SPEC-008`).
func TestReplicationClassWrapsCause(t *testing.T) {
	err := replicationClass(stderrors.New("verbindung weg"))
	if !stderrors.Is(err, outbound.ErrReplication) {
		t.Fatalf("Klasse replication fehlt: %v", err)
	}
	if !strings.Contains(err.Error(), "verbindung weg") {
		t.Fatalf("technische Ursache nicht lesbar: %v", err)
	}
}

// TestAcknowledgeSendsStandbyStatus trägt die Verklebung: die Bestätigung
// setzt genau eine Meldung ab, und sie trägt die Position.
func TestAcknowledgeSendsStandbyStatus(t *testing.T) {
	sender := &fakeSender{}
	ack, err := newOnSender(sender, WithLog(&recordingLog{}))
	if err != nil {
		t.Fatalf("newOnSender: %v", err)
	}
	position, err := model.NewSourcePosition("src-1", 100)
	if err != nil {
		t.Fatalf("Position: %v", err)
	}
	if err := ack.Acknowledge(context.Background(), position); err != nil {
		t.Fatalf("Acknowledge: %v", err)
	}
	if sender.calls != 1 || len(sender.sent) != 1 {
		t.Fatalf("Aufrufe: %d, Meldungen: %d, erwartet je 1", sender.calls, len(sender.sent))
	}
	got := sender.sent[0]
	if got.WALWritePosition != pglogrepl.LSN(100) ||
		got.WALFlushPosition != pglogrepl.LSN(100) ||
		got.WALApplyPosition != pglogrepl.LSN(100) {
		t.Fatalf("Meldung trägt %d/%d/%d, erwartet 100/100/100",
			got.WALWritePosition, got.WALFlushPosition, got.WALApplyPosition)
	}
}

// TestAcknowledgeZeroPositionSendsNothing trägt die Null-Positions-Grenze
// an der Verklebung: die Bestätigung ohne Position endet sichtbar und
// setzt nichts ab.
func TestAcknowledgeZeroPositionSendsNothing(t *testing.T) {
	sender := &fakeSender{}
	ack, err := newOnSender(sender)
	if err != nil {
		t.Fatalf("newOnSender: %v", err)
	}
	err = ack.Acknowledge(context.Background(), model.SourcePosition{})
	if !stderrors.Is(err, outbound.ErrReplication) {
		t.Fatalf("Null-Position: %v", err)
	}
	if sender.calls != 0 {
		t.Fatalf("Aufrufe: %d, erwartet 0", sender.calls)
	}
}

// TestAcknowledgeWrapsSenderFailure trägt den Fehlerpfad der Verklebung:
// ein Treiber-Fehler endet über `outbound.ErrReplication`, trägt die
// technische Ursache und geht durch den `LogPort` (`ADR-0023`,
// `LH-QA-OPS-004`).
func TestAcknowledgeWrapsSenderFailure(t *testing.T) {
	sender := &fakeSender{err: stderrors.New("verbindung weg")}
	log := &recordingLog{}
	ack, err := newOnSender(sender, WithLog(log))
	if err != nil {
		t.Fatalf("newOnSender: %v", err)
	}
	position, err := model.NewSourcePosition("src-1", 100)
	if err != nil {
		t.Fatalf("Position: %v", err)
	}
	err = ack.Acknowledge(context.Background(), position)
	if !stderrors.Is(err, outbound.ErrReplication) {
		t.Fatalf("Klasse replication fehlt: %v", err)
	}
	if !strings.Contains(err.Error(), "verbindung weg") {
		t.Fatalf("technische Ursache nicht lesbar: %v", err)
	}
	if len(log.errors) != 1 {
		t.Fatalf("Fehler-Log-Aufrufe: %d, erwartet 1", len(log.errors))
	}
	if sender.calls != 1 {
		t.Fatalf("Aufrufe: %d, erwartet 1", sender.calls)
	}
}

// TestNewOnSenderRejectsNilSender trägt die Konstruktions-Grenze des
// paket-internen Einstiegs: auch ohne Naht gibt es keinen ACK-Stand
// (`SPEC-008`).
func TestNewOnSenderRejectsNilSender(t *testing.T) {
	if _, err := newOnSender(nil); !stderrors.Is(err, outbound.ErrReplication) {
		t.Fatalf("newOnSender(nil): %v", err)
	}
}
