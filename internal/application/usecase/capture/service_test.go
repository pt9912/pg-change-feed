package capture_test

import (
	"context"
	stderrors "errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/capture"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Die Fakes tragen die Ordnungs-Logik des Capture-Pfads (`ADR-0030`,
// Fake-Ports statt realer Treiber): beide loggen ihre Aufrufe in dieselbe
// Ereignisliste, damit die Reihenfolge Persist vor ACK an einem Beleg
// prüfbar ist.

var (
	_ inbound.CaptureInboundPort  = (*capture.CaptureService)(nil)
	_ outbound.ChangeStorePort    = (*fakeStore)(nil)
	_ outbound.ReplicationAckPort = (*fakeAck)(nil)
)

type fakeStore struct {
	events     *[]string
	persistErr error
	persisted  []*model.ChangeTransaction
}

func (f *fakeStore) PersistTransaction(ctx context.Context, transaction *model.ChangeTransaction) error {
	*f.events = append(*f.events, "persist:"+string(transaction.ID))
	if f.persistErr != nil {
		return f.persistErr
	}
	f.persisted = append(f.persisted, transaction)
	return nil
}

// ReadChanges trägt die Lesefähigkeit am Port (`ADR-0009`); der Fake
// implementiert sie ohne Stand — der Capture-Pfad liest nicht, die reale
// Lesefähigkeit trägt der Adapter-Test gegen PostgreSQL.
func (f *fakeStore) ReadChanges(ctx context.Context, query outbound.ChangeQuery) ([]outbound.ChangeRecord, error) {
	return nil, stderrors.New("Fake-Store trägt nur die Persist-Seite")
}

type fakeAck struct {
	events *[]string
	ackErr error
	acked  []model.SourcePosition
}

func (f *fakeAck) Acknowledge(ctx context.Context, position model.SourcePosition) error {
	*f.events = append(*f.events, "ack:"+strconv.FormatUint(position.Offset, 10))
	if f.ackErr != nil {
		return f.ackErr
	}
	f.acked = append(f.acked, position)
	return nil
}

// service baut den CaptureService mit den Fakes und liefert die
// Ereignisliste als Ordnungs-Beleg.
func newService(t *testing.T) (*capture.CaptureService, *[]string, *fakeStore, *fakeAck) {
	t.Helper()
	events := []string{}
	store := &fakeStore{events: &events}
	ack := &fakeAck{events: &events}
	return capture.NewCaptureService(store, ack), &events, store, ack
}

// committedTransaction legt eine committed Quelltransaktion mit count
// Changes an (`LH-FA-CAP-005`, `LH-FA-DAT-004`).
func committedTransaction(t *testing.T, id model.TransactionID, offset uint64, count int) *model.ChangeTransaction {
	t.Helper()
	tx, err := model.NewOpenTransaction(id, "src-1")
	if err != nil {
		t.Fatalf("NewOpenTransaction: %v", err)
	}
	for i := 1; i <= count; i++ {
		change, err := model.NewChange(
			model.ChangeID(fmt.Sprintf("%s-%d", id, i)),
			id,
			"tbl-1",
			int64(i),
			model.OperationInsert,
			nil,
			[]byte(`{"x":1}`),
			"sv-1",
		)
		if err != nil {
			t.Fatalf("NewChange %d: %v", i, err)
		}
		if err := tx.AppendChange(change); err != nil {
			t.Fatalf("AppendChange %d: %v", i, err)
		}
	}
	position, err := model.NewSourcePosition("src-1", offset)
	if err != nil {
		t.Fatalf("NewSourcePosition: %v", err)
	}
	if err := tx.Commit(position); err != nil {
		t.Fatalf("Commit: %v", err)
	}
	return tx
}

// Ordnung (`LH-QA-REL-001.a`): der ACK trägt die Commit-Position erst
// nach dem Store-Commit; der Ereignisbeleg trägt Persist vor ACK.
func TestCapturePersistsBeforeAck(t *testing.T) {
	service, events, store, ack := newService(t)
	tx := committedTransaction(t, "t-1", 100, 2)

	result, err := service.Capture(context.Background(), capture.CaptureCommand{Transaction: tx})
	if err != nil {
		t.Fatalf("Capture: %v", err)
	}
	if got, want := *events, []string{"persist:t-1", "ack:100"}; fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("Ereignisse = %v, wollen %v", got, want)
	}
	if len(store.persisted) != 1 || store.persisted[0].ID != tx.ID {
		t.Fatalf("Store trägt %d Transaktionen, wollen die von %s", len(store.persisted), tx.ID)
	}
	if len(ack.acked) != 1 || ack.acked[0].Offset != 100 {
		t.Fatalf("ACK trägt %v, wollen Offset 100", ack.acked)
	}
	if result.Acknowledged.Offset != 100 {
		t.Fatalf("Ergebnis trägt Offset %d, wollen 100", result.Acknowledged.Offset)
	}
}

// Eine committed Quelltransaktion ohne Changes passiert den Pfad: sie
// wird mit ihrer Commit-Position persistiert und geackt (`LH-FA-CAP-006.a`,
// Grenze am `CaptureInboundPort`) — das Ablehnen würde dieselbe leere
// Transaktion in einer Wiederholung (`ADR-0012`) endlos neu liefern.
func TestCapturePersistsEmptyCommittedTransaction(t *testing.T) {
	service, events, store, ack := newService(t)
	tx := committedTransaction(t, "t-1", 100, 0)

	result, err := service.Capture(context.Background(), capture.CaptureCommand{Transaction: tx})
	if err != nil {
		t.Fatalf("Capture: %v", err)
	}
	if len(store.persisted) != 1 || store.persisted[0].ID != tx.ID {
		t.Fatalf("Store trägt %d Transaktionen, wollen die von %s", len(store.persisted), tx.ID)
	}
	if len(ack.acked) != 1 || ack.acked[0].Offset != 100 {
		t.Fatalf("ACK trägt %v, wollen Offset 100", ack.acked)
	}
	if result.Acknowledged.Offset != 100 {
		t.Fatalf("Ergebnis trägt Offset %d, wollen 100", result.Acknowledged.Offset)
	}
	if got, want := *events, []string{"persist:t-1", "ack:100"}; fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("Ereignisse = %v, wollen %v", got, want)
	}
}

// Persistenzfehler → kein Source-ACK (`LH-QA-REL-001.a`; `SPEC-008`,
// Klasse `storage`): der Store-Commit endet mit Fehler, der Fake-ACK
// trägt keine Bestätigung.
func TestCaptureDoesNotAckOnPersistenceError(t *testing.T) {
	service, events, store, ack := newService(t)
	storeErr := stderrors.New("storage defekt")
	store.persistErr = storeErr
	tx := committedTransaction(t, "t-1", 100, 1)

	_, err := service.Capture(context.Background(), capture.CaptureCommand{Transaction: tx})
	if !stderrors.Is(err, storeErr) {
		t.Fatalf("Capture-Fehler = %v, wollen %v", err, storeErr)
	}
	if len(ack.acked) != 0 {
		t.Fatalf("ACK trägt %v nach Persistenzfehler, wollen keine Bestätigung", ack.acked)
	}
	for _, event := range *events {
		if event[:3] == "ack" {
			t.Fatalf("Ereignisbeleg trägt %q nach Persistenzfehler", event)
		}
	}
}

// Offene Transaktionen sind nicht konsumierbar (`LH-FA-CAP-006`);
// zurückgerollte Transaktionen erreichen den Commit nicht (`LH-FA-CAP-007`)
// — der Capture-Pfad persistiert sie nicht und bestätigt nichts.
func TestCaptureRejectsOpenTransaction(t *testing.T) {
	service, _, store, ack := newService(t)
	tx, err := model.NewOpenTransaction("t-1", "src-1")
	if err != nil {
		t.Fatalf("NewOpenTransaction: %v", err)
	}

	_, err = service.Capture(context.Background(), capture.CaptureCommand{Transaction: tx})
	if !stderrors.Is(err, domainerrors.ErrTransactionNotCommitted) {
		t.Fatalf("Capture-Fehler = %v, wollen %v", err, domainerrors.ErrTransactionNotCommitted)
	}
	if len(store.persisted) != 0 {
		t.Fatalf("Store trägt %d Transaktionen, wollen 0", len(store.persisted))
	}
	if len(ack.acked) != 0 {
		t.Fatalf("ACK trägt %v, wollen keine Bestätigung", ack.acked)
	}
}

// Der Aufruf ohne Quelltransaktion endet an der Port-Grenze.
func TestCaptureRejectsMissingTransaction(t *testing.T) {
	service, events, store, ack := newService(t)

	_, err := service.Capture(context.Background(), capture.CaptureCommand{})
	if !stderrors.Is(err, capture.ErrMissingTransaction) {
		t.Fatalf("Capture-Fehler = %v, wollen %v", err, capture.ErrMissingTransaction)
	}
	if len(store.persisted) != 0 || len(ack.acked) != 0 || len(*events) != 0 {
		t.Fatalf("Store oder ACK tragen %d/%d Aufrufe, wollen 0/0", len(store.persisted), len(ack.acked))
	}
}

// Crash zwischen Persistenz und ACK (`LH-QA-REL-002`, `LH-QA-REL-001.a`):
// die Persistenz besteht, der ACK fehlt; der erneute Aufruf mit derselben
// Transaktion führt die Ordnung zu Ende — Wiederholung statt Datenverlust
// (`ADR-0012`). Der Fake-Store trägt die Transaktion dann zweimal; die
// reale Idempotenz des Stores trägt die Deduplizierungsbasis (`SPEC-002`,
// change_id) und bleibt Beleg des realen Adapters.
func TestCrashBetweenPersistAndAckLeadsToReprocessing(t *testing.T) {
	service, events, store, ack := newService(t)
	ackErr := stderrors.New("Quelle zwischen Persistenz und ACK nicht erreichbar")
	ack.ackErr = ackErr
	tx := committedTransaction(t, "t-1", 100, 1)

	if _, err := service.Capture(context.Background(), capture.CaptureCommand{Transaction: tx}); !stderrors.Is(err, ackErr) {
		t.Fatalf("erster Capture-Fehler = %v, wollen %v", err, ackErr)
	}
	if len(store.persisted) != 1 {
		t.Fatalf("Store trägt %d Transaktionen nach fehlgeschlagenem ACK, wollen 1", len(store.persisted))
	}
	if len(ack.acked) != 0 {
		t.Fatalf("ACK trägt %v nach fehlgeschlagenem ACK, wollen keine Bestätigung", ack.acked)
	}

	// Der Neustart findet die Quelle wieder erreichbar: die Störung des
	// ersten Laufs trägt der erneute Aufruf nicht mit.
	ack.ackErr = nil
	result, err := service.Capture(context.Background(), capture.CaptureCommand{Transaction: tx})
	if err != nil {
		t.Fatalf("erneuter Capture: %v", err)
	}
	if result.Acknowledged.Offset != 100 {
		t.Fatalf("Ergebnis trägt Offset %d, wollen 100", result.Acknowledged.Offset)
	}
	if len(store.persisted) != 2 {
		t.Fatalf("Store trägt %d Persistierungen, wollen 2 (erneute Verarbeitung)", len(store.persisted))
	}
	if got, want := *events, []string{"persist:t-1", "ack:100", "persist:t-1", "ack:100"}; fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("Ereignisse = %v, wollen %v", got, want)
	}
}
