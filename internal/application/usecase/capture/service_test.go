package capture_test

import (
	"context"
	stderrors "errors"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

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
	_ inbound.CaptureInboundPort      = (*capture.CaptureService)(nil)
	_ outbound.ChangeStorePort        = (*fakeStore)(nil)
	_ outbound.ReplicationAckPort     = (*fakeAck)(nil)
	_ outbound.ChangeNotificationPort = (*fakeNotify)(nil)
	_ outbound.ChangeStreamPort       = (*fakeStream)(nil)
	_ outbound.ChangeStreamPort       = (*fakeBlockingStream)(nil)
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

// DeleteChanges trägt die Löschfähigkeit am Port (`ADR-0014`); der Fake
// implementiert sie ohne Stand — der Capture-Pfad löscht nicht, die reale
// Löschfähigkeit trägt der Adapter-Test gegen PostgreSQL und der
// RunRetention-Use-Case-Test seine eigenen Fakes.
func (f *fakeStore) DeleteChanges(ctx context.Context, changeIDs []model.ChangeID) error {
	return stderrors.New("Fake-Store trägt nur die Persist-Seite")
}

// fakeAck trägt den `ReplicationAckPort`: `ackErr` lässt jeden Aufruf
// scheitern; mit `failFrom` (1-basiert) scheitern erst die Aufrufe ab dem
// n-ten.
type fakeAck struct {
	events   *[]string
	ackErr   error
	failFrom int
	calls    int
	acked    []model.SourcePosition
}

func (f *fakeAck) Acknowledge(ctx context.Context, position model.SourcePosition) error {
	f.calls++
	*f.events = append(*f.events, "ack:"+strconv.FormatUint(position.Offset, 10))
	if f.ackErr != nil && f.calls >= f.failFrom {
		return f.ackErr
	}
	f.acked = append(f.acked, position)
	return nil
}

// fakeNotify trägt den Regressionsbeleg für `ADR-0055`: ein fehlschlagender
// `ChangeNotificationPort` darf `Capture()` nicht scheitern lassen, wenn
// `store`/`ack` erfolgreich waren — `notifyErr` ist standardmäßig gesetzt,
// damit der Test den ungünstigsten Fall trägt. `notified` trägt die
// distinkten `(schema, table)`-Aufrufe für den Deduplizierungs-Beleg
// (`ADR-0056`).
type fakeNotify struct {
	events    *[]string
	notifyErr error
	notified  []string
}

func (f *fakeNotify) Notify(ctx context.Context, sourceID, schema, table string) error {
	*f.events = append(*f.events, fmt.Sprintf("notify:%s.%s.%s", sourceID, schema, table))
	if f.notifyErr != nil {
		return f.notifyErr
	}
	f.notified = append(f.notified, schema+"."+table)
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
// Changes an derselben Tabelle `public.tbl-1` an (`LH-FA-CAP-005`,
// `LH-FA-DAT-004`).
func committedTransaction(t *testing.T, id model.TransactionID, offset uint64, count int) *model.ChangeTransaction {
	t.Helper()
	return committedTransactionOverTables(t, id, offset, tableSpec{schema: "public", table: "tbl-1", count: count})
}

// tableSpec trägt die Change-Zahl einer Tabelle für
// `committedTransactionOverTables` (`ADR-0056`-Deduplizierungsbelege).
type tableSpec struct {
	schema string
	table  string
	count  int
}

// committedTransactionOverTables legt eine committed Quelltransaktion an,
// deren Changes über mehrere `(schema, table)`-Paare verteilt sein können
// (`ADR-0056`, Notify-Deduplizierung je distinktem Paar). Jeder Change
// trägt `SourceTableID` als `<schema>.<table>` (analog zum
// Driving-Adapter-Mapper) und die Klartext-Felder `Schema`/`Table`.
func committedTransactionOverTables(t *testing.T, id model.TransactionID, offset uint64, specs ...tableSpec) *model.ChangeTransaction {
	t.Helper()
	tx, err := model.NewOpenTransaction(id, "src-1")
	if err != nil {
		t.Fatalf("NewOpenTransaction: %v", err)
	}
	sequence := int64(0)
	for _, spec := range specs {
		for i := 1; i <= spec.count; i++ {
			sequence++
			change, err := model.NewChange(
				model.ChangeID(fmt.Sprintf("%s-%d", id, sequence)),
				id,
				model.SourceTableID(spec.schema+"."+spec.table),
				sequence,
				model.OperationInsert,
				nil,
				[]byte(`{"x":1}`),
				"sv-1",
			)
			if err != nil {
				t.Fatalf("NewChange %d: %v", sequence, err)
			}
			change.Schema = spec.schema
			change.Table = spec.table
			if err := tx.AppendChange(change); err != nil {
				t.Fatalf("AppendChange %d: %v", sequence, err)
			}
		}
	}
	position, err := model.NewSourcePosition("src-1", offset)
	if err != nil {
		t.Fatalf("NewSourcePosition: %v", err)
	}
	if err := tx.Commit(position, model.NewTimePoint(1)); err != nil {
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

// Ohne konfigurierten `ChangeNotificationPort` bleibt das Bestandsverhalten
// bit-identisch: kein Notify-Versuch, kein zusätzliches Ereignis
// (`ADR-0055` Punkt 5, Boundary: ungesetzt bedeutet deaktiviert).
func TestCaptureWithoutNotificationPortLeavesBehaviourUnchanged(t *testing.T) {
	service, events, store, ack := newService(t)
	tx := committedTransaction(t, "t-1", 100, 1)

	result, err := service.Capture(context.Background(), capture.CaptureCommand{Transaction: tx})
	if err != nil {
		t.Fatalf("Capture: %v", err)
	}
	if got, want := *events, []string{"persist:t-1", "ack:100"}; fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("Ereignisse = %v, wollen %v (kein Notify ohne konfigurierten Port)", got, want)
	}
	if result.Acknowledged.Offset != 100 {
		t.Fatalf("Ergebnis trägt Offset %d, wollen 100", result.Acknowledged.Offset)
	}
	_ = store
	_ = ack
}

// Erfolgreicher Notify-Aufruf reiht sich NACH ACK Source ein (`ADR-0055`:
// Receive → Decode → Persist → COMMIT Store → ACK Source → Notify).
func TestCaptureNotifiesAfterAckOnSuccess(t *testing.T) {
	events := []string{}
	store := &fakeStore{events: &events}
	ack := &fakeAck{events: &events}
	notify := &fakeNotify{events: &events}
	service := capture.NewCaptureService(store, ack, capture.WithChangeNotification(notify))
	tx := committedTransaction(t, "t-1", 100, 1)

	result, err := service.Capture(context.Background(), capture.CaptureCommand{Transaction: tx})
	if err != nil {
		t.Fatalf("Capture: %v", err)
	}
	if got, want := events, []string{"persist:t-1", "ack:100", "notify:src-1.public.tbl-1"}; fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("Ereignisse = %v, wollen %v", got, want)
	}
	if len(notify.notified) != 1 || notify.notified[0] != "public.tbl-1" {
		t.Fatalf("Notify trägt %v, wollen [\"public.tbl-1\"]", notify.notified)
	}
	if result.Acknowledged.Offset != 100 {
		t.Fatalf("Ergebnis trägt Offset %d, wollen 100", result.Acknowledged.Offset)
	}
}

// Regressionstest — die wichtigste Einzeleigenschaft aus `ADR-0055`: ein
// fehlschlagender `ChangeNotificationPort` darf `Capture()` nicht scheitern
// lassen, wenn `store`/`ack` bereits erfolgreich waren: entfernt man das
// Abfangen in `CaptureService.Capture()` (`if err := s.notify.Notify(...);
// err != nil { … }` durch `return CaptureResult{}, err` ersetzt), schlägt
// genau dieser Test fehl.
func TestCaptureSucceedsDespiteFailingNotification(t *testing.T) {
	events := []string{}
	store := &fakeStore{events: &events}
	ack := &fakeAck{events: &events}
	notifyErr := stderrors.New("NATS nicht erreichbar")
	notify := &fakeNotify{events: &events, notifyErr: notifyErr}
	service := capture.NewCaptureService(store, ack, capture.WithChangeNotification(notify))
	tx := committedTransaction(t, "t-1", 100, 1)

	result, err := service.Capture(context.Background(), capture.CaptureCommand{Transaction: tx})
	if err != nil {
		t.Fatalf("Capture: %v, wollen keinen Fehler trotz fehlschlagendem Notify", err)
	}
	if result.Acknowledged.Offset != 100 {
		t.Fatalf("Ergebnis trägt Offset %d, wollen 100", result.Acknowledged.Offset)
	}
	if len(store.persisted) != 1 || len(ack.acked) != 1 {
		t.Fatalf("Store/ACK tragen %d/%d, wollen 1/1 trotz fehlschlagendem Notify", len(store.persisted), len(ack.acked))
	}
	if len(notify.notified) != 0 {
		t.Fatalf("Notify trägt %v, wollen keine erfolgreiche Zustellung", notify.notified)
	}
	if got, want := events, []string{"persist:t-1", "ack:100", "notify:src-1.public.tbl-1"}; fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("Ereignisse = %v, wollen %v (Notify-Versuch trotz Fehlschlag geloggt)", got, want)
	}
}

// TestCaptureNotifiesOnceForSameTableMultipleChanges trägt die
// Deduplizierung (`ADR-0056`): mehrere Changes derselben Tabelle in einer
// Transaktion lösen genau ein Notify für dieses `(schema, table)`-Paar aus.
func TestCaptureNotifiesOnceForSameTableMultipleChanges(t *testing.T) {
	events := []string{}
	store := &fakeStore{events: &events}
	ack := &fakeAck{events: &events}
	notify := &fakeNotify{events: &events}
	service := capture.NewCaptureService(store, ack, capture.WithChangeNotification(notify))
	tx := committedTransactionOverTables(t, "t-1", 100, tableSpec{schema: "public", table: "tbl-1", count: 3})

	if _, err := service.Capture(context.Background(), capture.CaptureCommand{Transaction: tx}); err != nil {
		t.Fatalf("Capture: %v", err)
	}
	if got, want := notify.notified, []string{"public.tbl-1"}; fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("Notify trägt %v, wollen %v (genau ein Aufruf trotz drei Changes derselben Tabelle)", got, want)
	}
}

// TestCaptureNotifiesDistinctlyForTwoTables trägt die Kehrseite: Changes
// über zwei Tabellen lösen zwei distinkte Notify-Aufrufe aus (`ADR-0056`).
func TestCaptureNotifiesDistinctlyForTwoTables(t *testing.T) {
	events := []string{}
	store := &fakeStore{events: &events}
	ack := &fakeAck{events: &events}
	notify := &fakeNotify{events: &events}
	service := capture.NewCaptureService(store, ack, capture.WithChangeNotification(notify))
	tx := committedTransactionOverTables(t, "t-1", 100,
		tableSpec{schema: "public", table: "tbl-1", count: 2},
		tableSpec{schema: "public", table: "tbl-2", count: 1},
	)

	if _, err := service.Capture(context.Background(), capture.CaptureCommand{Transaction: tx}); err != nil {
		t.Fatalf("Capture: %v", err)
	}
	if got, want := notify.notified, []string{"public.tbl-1", "public.tbl-2"}; fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("Notify trägt %v, wollen %v (zwei distinkte Aufrufe für zwei Tabellen)", got, want)
	}
}

// fakeStream trägt den optionalen `ChangeStreamPort` (`ADR-0060` Teilfrage 2)
// und spiegelt die Zustellsemantik des realen Ports (`ADR-0066`): er kehrt
// unabhängig davon zurück, ob ein Empfänger liest. `publishErr` trägt den
// Regressionsbeleg — er ist standardmäßig gesetzt, damit der ungünstigste
// Fall im Test steht. `nichtLesenderEmpfaenger` modelliert einen
// registrierten, nie lesenden Empfänger über dessen begrenzte
// Empfangs-Warteschlange; ist sie voll, verwirft der Aufruf den eintreffenden
// Change für diesen Empfänger (Drop-Newest) statt zu warten.
type fakeStream struct {
	events                  *[]string
	publishErr              error
	published               []model.ChangeID
	nichtLesenderEmpfaenger chan *model.Change
}

func (f *fakeStream) Publish(ctx context.Context, change *model.Change) error {
	*f.events = append(*f.events, "stream:"+string(change.ID))
	if f.nichtLesenderEmpfaenger != nil {
		select {
		case f.nichtLesenderEmpfaenger <- change:
		default:
		}
	}
	if f.publishErr != nil {
		return f.publishErr
	}
	f.published = append(f.published, change.ID)
	return nil
}

// fakeBlockingStream trägt den ungünstigsten Port-Zustand: `Publish` meldet
// seinen Eintritt und kehrt erst zurück, wenn der Test die Freigabe schließt —
// ein Port, der entgegen seinem Vertrag nicht zurückkehrt. `eingetreten` ist
// gepuffert, damit die Meldung selbst nicht zum Warten wird.
type fakeBlockingStream struct {
	eingetreten chan model.ChangeID
	freigabe    chan struct{}
}

func (f *fakeBlockingStream) Publish(ctx context.Context, change *model.Change) error {
	select {
	case f.eingetreten <- change.ID:
	default:
	}
	<-f.freigabe
	return nil
}

// TestCapturePublishesEachChangeOnceAfterAckAndNotify trägt den Happy Path
// aus `LH-FA-SST-008` und die Ordnung aus `ADR-0060` Teilfrage 2: der
// Stream-Publish reiht sich als letzter Best-Effort-Schritt ein — nach
// `ACK Source` und nach dem Wecksignal — und ruft je Change der Transaktion
// genau einmal auf, in der Reihenfolge von `tx.Changes()` und ohne
// Deduplizierung nach Tabelle (zwei Changes derselben Tabelle ergeben zwei
// Aufrufe, anders als beim tabellen-granularen Notify).
func TestCapturePublishesEachChangeOnceAfterAckAndNotify(t *testing.T) {
	events := []string{}
	store := &fakeStore{events: &events}
	ack := &fakeAck{events: &events}
	notify := &fakeNotify{events: &events}
	stream := &fakeStream{events: &events}
	service := capture.NewCaptureService(store, ack, capture.WithChangeNotification(notify), capture.WithChangeStream(stream))
	tx := committedTransactionOverTables(t, "t-1", 100, tableSpec{schema: "public", table: "tbl-1", count: 2})

	result, err := service.Capture(context.Background(), capture.CaptureCommand{Transaction: tx})
	if err != nil {
		t.Fatalf("Capture: %v", err)
	}
	want := []string{"persist:t-1", "ack:100", "notify:src-1.public.tbl-1", "stream:t-1-1", "stream:t-1-2"}
	if fmt.Sprint(events) != fmt.Sprint(want) {
		t.Fatalf("Ereignisse = %v, wollen %v", events, want)
	}
	if got, wantIDs := stream.published, []model.ChangeID{"t-1-1", "t-1-2"}; fmt.Sprint(got) != fmt.Sprint(wantIDs) {
		t.Fatalf("Stream trägt %v, wollen %v", got, wantIDs)
	}
	if result.Acknowledged.Offset != 100 {
		t.Fatalf("Ergebnis trägt Offset %d, wollen 100", result.Acknowledged.Offset)
	}
}

// TestCaptureWithoutStreamPortLeavesBehaviourUnchanged trägt die Boundary aus
// `ADR-0060` Teilfrage 6: ohne konfigurierten `ChangeStreamPort` bleibt das
// Bestandsverhalten bit-identisch — kein Publish-Versuch, kein zusätzliches
// Ereignis.
func TestCaptureWithoutStreamPortLeavesBehaviourUnchanged(t *testing.T) {
	events := []string{}
	store := &fakeStore{events: &events}
	ack := &fakeAck{events: &events}
	notify := &fakeNotify{events: &events}
	service := capture.NewCaptureService(store, ack, capture.WithChangeNotification(notify))
	tx := committedTransaction(t, "t-1", 100, 2)

	result, err := service.Capture(context.Background(), capture.CaptureCommand{Transaction: tx})
	if err != nil {
		t.Fatalf("Capture: %v", err)
	}
	want := []string{"persist:t-1", "ack:100", "notify:src-1.public.tbl-1"}
	if fmt.Sprint(events) != fmt.Sprint(want) {
		t.Fatalf("Ereignisse = %v, wollen %v (kein Stream-Publish ohne konfigurierten Port)", events, want)
	}
	if result.Acknowledged.Offset != 100 {
		t.Fatalf("Ergebnis trägt Offset %d, wollen 100", result.Acknowledged.Offset)
	}
}

// Regressionstest — die wichtigste Einzeleigenschaft aus `ADR-0060`
// Teilfrage 2 (`SPEC-020`, Zeile *Fehler bei Publish-Fehlschlag*): ein
// fehlschlagender `ChangeStreamPort` darf `Capture()` nicht scheitern lassen,
// wenn `store`/`ack` bereits erfolgreich waren. Der Aufruf läuft über alle
// Changes weiter — ein Fehler auf dem ersten verwirft die übrigen nicht.
// Entfernt man das Abfangen in `CaptureService.Capture()` (die
// `s.log.Warn`-Verzweigung durch `return CaptureResult{}, err` ersetzt),
// schlägt genau dieser Test fehl.
func TestCaptureSucceedsDespiteFailingStream(t *testing.T) {
	events := []string{}
	store := &fakeStore{events: &events}
	ack := &fakeAck{events: &events}
	streamErr := stderrors.New("Stream nicht verfügbar")
	stream := &fakeStream{events: &events, publishErr: streamErr}
	service := capture.NewCaptureService(store, ack, capture.WithChangeStream(stream))
	tx := committedTransactionOverTables(t, "t-1", 100, tableSpec{schema: "public", table: "tbl-1", count: 2})

	result, err := service.Capture(context.Background(), capture.CaptureCommand{Transaction: tx})
	if err != nil {
		t.Fatalf("Capture: %v, wollen keinen Fehler trotz fehlschlagendem Stream-Publish", err)
	}
	if result.Acknowledged.Offset != 100 {
		t.Fatalf("Ergebnis trägt Offset %d, wollen 100", result.Acknowledged.Offset)
	}
	if len(store.persisted) != 1 || len(ack.acked) != 1 {
		t.Fatalf("Store/ACK tragen %d/%d, wollen 1/1 trotz fehlschlagendem Stream-Publish", len(store.persisted), len(ack.acked))
	}
	if len(stream.published) != 0 {
		t.Fatalf("Stream trägt %v, wollen keine erfolgreiche Zustellung", stream.published)
	}
	want := []string{"persist:t-1", "ack:100", "stream:t-1-1", "stream:t-1-2"}
	if fmt.Sprint(events) != fmt.Sprint(want) {
		t.Fatalf("Ereignisse = %v, wollen %v (Publish-Versuch je Change trotz Fehlschlag)", events, want)
	}
}

// TestCaptureKehrtOhneUndMitNichtLesendemStreamEmpfaengerZurueck trägt die
// Fire-and-Forget-Hälfte aus `ADR-0066`: `Capture()` kehrt ohne eigene
// Zeit-Isolation zurück — bei getrenntem Client (kein Abonnent registriert)
// ebenso wie bei einem registrierten, nie lesenden Empfänger. Die
// Entkopplung trägt der Port (nicht-blockierender Send mit Drop-Newest),
// nicht diese Aufrufstelle. Der Frist-Beleg macht einen blockierenden Aufruf
// sofort sichtbar, statt den Testlauf bis zum globalen Timeout hängen zu
// lassen.
func TestCaptureKehrtOhneUndMitNichtLesendemStreamEmpfaengerZurueck(t *testing.T) {
	faelle := []struct {
		name                         string
		empfaenger                   chan *model.Change
		erwarteteErfolgreicheAufrufe int
	}{
		{name: "getrennter Client", empfaenger: nil, erwarteteErfolgreicheAufrufe: 2},
		{name: "nicht lesender Empfänger", empfaenger: make(chan *model.Change, 1), erwarteteErfolgreicheAufrufe: 2},
	}
	for _, fall := range faelle {
		t.Run(fall.name, func(t *testing.T) {
			events := []string{}
			store := &fakeStore{events: &events}
			ack := &fakeAck{events: &events}
			stream := &fakeStream{events: &events, nichtLesenderEmpfaenger: fall.empfaenger}
			service := capture.NewCaptureService(store, ack, capture.WithChangeStream(stream))
			tx := committedTransactionOverTables(t, "t-1", 100, tableSpec{schema: "public", table: "tbl-1", count: 2})

			type ergebnis struct {
				result capture.CaptureResult
				err    error
			}
			fertig := make(chan ergebnis, 1)
			go func() {
				result, err := service.Capture(context.Background(), capture.CaptureCommand{Transaction: tx})
				fertig <- ergebnis{result: result, err: err}
			}()

			select {
			case got := <-fertig:
				if got.err != nil {
					t.Fatalf("Capture: %v", got.err)
				}
				if got.result.Acknowledged.Offset != 100 {
					t.Fatalf("Ergebnis trägt Offset %d, wollen 100", got.result.Acknowledged.Offset)
				}
			case <-time.After(2 * time.Second):
				t.Fatal("Capture kehrt nicht zurück")
			}

			if got, want := len(stream.published), fall.erwarteteErfolgreicheAufrufe; got != want {
				t.Fatalf("Stream trägt %d erfolgreiche Aufrufe, wollen %d (der Port kehrt je Change zurück)", got, want)
			}
			if fall.empfaenger == nil {
				return
			}
			select {
			case got := <-fall.empfaenger:
				if got.ID != "t-1-1" {
					t.Fatalf("Empfangswarteschlange trägt %q, wollen die erste Change-ID", got.ID)
				}
			default:
				t.Fatalf("Empfangswarteschlange des nicht lesenden Empfängers ist leer")
			}
			select {
			case got := <-fall.empfaenger:
				t.Fatalf("über die Kapazität hinausgehender Change wurde eingereiht: %q", got.ID)
			default:
			}
		})
	}
}

// TestCapturePublishOhneRueckkehrHaeltKritischeKetteNichtAn pinnt die Grenze
// aus `ADR-0066`: die Entkopplung liegt im Port — `Publish` kehrt nicht auf
// einen Abonnenten wartend zurück —, diese Aufrufstelle ruft ihn synchron in
// der Best-Effort-Kette. Ein Port, der dennoch nicht zurückkehrt, kann die
// Capture-kritische Kette nicht mehr berühren: Persistenz und Source-ACK sind
// geschrieben, bevor der erste Change das Haus verlässt (`ADR-0011`,
// `ADR-0027`), und der Rückgabewert von `Capture()` hängt an einem Port, der
// seinen Vertrag nicht verletzt.
func TestCapturePublishOhneRueckkehrHaeltKritischeKetteNichtAn(t *testing.T) {
	events := []string{}
	store := &fakeStore{events: &events}
	ack := &fakeAck{events: &events}
	stream := &fakeBlockingStream{eingetreten: make(chan model.ChangeID, 1), freigabe: make(chan struct{})}
	service := capture.NewCaptureService(store, ack, capture.WithChangeStream(stream))
	tx := committedTransaction(t, "t-1", 100, 2)

	fertig := make(chan error, 1)
	go func() {
		_, err := service.Capture(context.Background(), capture.CaptureCommand{Transaction: tx})
		fertig <- err
	}()

	select {
	case <-stream.eingetreten:
	case <-time.After(2 * time.Second):
		t.Fatal("Capture erreicht den Stream-Publish nicht")
	}
	if len(store.persisted) != 1 {
		t.Fatalf("Store trägt %d Transaktionen beim Eintritt in den Stream-Publish, wollen 1", len(store.persisted))
	}
	if len(ack.acked) != 1 || ack.acked[0].Offset != 100 {
		t.Fatalf("ACK trägt %v beim Eintritt in den Stream-Publish, wollen Offset 100", ack.acked)
	}
	if got, want := fmt.Sprint(events), fmt.Sprint([]string{"persist:t-1", "ack:100"}); got != want {
		t.Fatalf("Ereignisse = %v, wollen %v vor dem ersten Stream-Publish", got, want)
	}

	close(stream.freigabe)
	if err := <-fertig; err != nil {
		t.Fatalf("Capture: %v, wollen keinen Fehler", err)
	}
}

// fakeLog trägt den `LogPort` (`ADR-0024`) als Aufzeichnung: er merkt sich
// die Warn-Zeilen samt Nachricht und Attributen, damit der Test belegt, über
// welchen Port die Best-Effort-Fehlerpfade melden.
type fakeLog struct {
	warnings []string
}

func (f *fakeLog) Debug(ctx context.Context, msg string, attrs ...any) {}
func (f *fakeLog) Info(ctx context.Context, msg string, attrs ...any)  {}

func (f *fakeLog) Warn(ctx context.Context, msg string, attrs ...any) {
	f.warnings = append(f.warnings, fmt.Sprintf("%s %v", msg, attrs))
}

func (f *fakeLog) Error(ctx context.Context, msg string, attrs ...any) {}

var _ outbound.LogPort = (*fakeLog)(nil)

// TestCaptureLoggtFehlschlaegeUeberDenInjiziertenPort trägt `ADR-0024`
// (strukturiertes Logging bleibt Infrastruktur und wird über den Outbound
// Port substituiert): beide Best-Effort-Fehlerpfade — Wecksignal (`ADR-0055`)
// und Stream-Publish (`ADR-0060` Teilfrage 2) — melden ihren Fehlschlag über
// den injizierten `LogPort`; das Attribut der Warn-Zeile trägt die
// gescheiterte Tabelle bzw. den gescheiterten Change.
func TestCaptureLoggtFehlschlaegeUeberDenInjiziertenPort(t *testing.T) {
	events := []string{}
	store := &fakeStore{events: &events}
	ack := &fakeAck{events: &events}
	notify := &fakeNotify{events: &events, notifyErr: stderrors.New("NATS nicht erreichbar")}
	stream := &fakeStream{events: &events, publishErr: stderrors.New("Stream nicht verfügbar")}
	log := &fakeLog{}
	service := capture.NewCaptureService(store, ack,
		capture.WithChangeNotification(notify),
		capture.WithChangeStream(stream),
		capture.WithLog(log),
	)
	tx := committedTransaction(t, "t-1", 100, 1)

	if _, err := service.Capture(context.Background(), capture.CaptureCommand{Transaction: tx}); err != nil {
		t.Fatalf("Capture: %v, wollen keinen Fehler trotz fehlschlagender Best-Effort-Ports", err)
	}
	if len(log.warnings) != 2 {
		t.Fatalf("Warn-Zeilen = %v, wollen 2 (Wecksignal und Stream-Publish)", log.warnings)
	}
	if !strings.Contains(log.warnings[0], "Wecksignal fehlgeschlagen") || !strings.Contains(log.warnings[0], "tbl-1") {
		t.Fatalf("Wecksignal-Warnzeile = %q, wollen Nachricht und Tabelle", log.warnings[0])
	}
	if !strings.Contains(log.warnings[1], "Stream-Publish fehlgeschlagen") || !strings.Contains(log.warnings[1], "t-1-1") {
		t.Fatalf("Stream-Warnzeile = %q, wollen Nachricht und Change-Kennung", log.warnings[1])
	}
}

// idlePosition trägt die Leerlauf-Position der Tests: eine Position der
// Quelle, die kein Change trägt.
func idlePosition(t *testing.T, offset uint64) model.SourcePosition {
	t.Helper()
	position, err := model.NewSourcePosition("src-1", offset)
	if err != nil {
		t.Fatalf("NewSourcePosition: %v", err)
	}
	return position
}

// Leerlauf-Bestätigung (`ADR-0120` Festlegung 1): die Application bestätigt
// die gemeldete Position ausschließlich über den `ReplicationAckPort` — der
// Ereignisbeleg trägt genau ein `ack` mit dieser Position und weder
// Persistenz noch Wecksignal noch Stream-Publish, obwohl beide
// Best-Effort-Ports verdrahtet sind.
func TestConfirmIdleAcknowledgesOnlyThroughTheAckPort(t *testing.T) {
	events := []string{}
	store := &fakeStore{events: &events}
	ack := &fakeAck{events: &events}
	notify := &fakeNotify{events: &events}
	stream := &fakeStream{events: &events}
	service := capture.NewCaptureService(store, ack,
		capture.WithChangeNotification(notify), capture.WithChangeStream(stream))

	result, err := service.ConfirmIdle(context.Background(), inbound.IdleConfirmationCommand{Position: idlePosition(t, 500)})
	if err != nil {
		t.Fatalf("ConfirmIdle: %v", err)
	}
	if got, want := events, []string{"ack:500"}; fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("Ereignisse = %v, wollen %v", got, want)
	}
	if len(store.persisted) != 0 {
		t.Fatalf("Store trägt %d Transaktionen, wollen keine", len(store.persisted))
	}
	if len(ack.acked) != 1 || ack.acked[0].Offset != 500 {
		t.Fatalf("ACK trägt %v, wollen Offset 500", ack.acked)
	}
	if result.Acknowledged.Offset != 500 {
		t.Fatalf("Ergebnis trägt Offset %d, wollen 500", result.Acknowledged.Offset)
	}
}

// Ein Fehler des Ports geht unverändert durch (`ADR-0120` Festlegung 3):
// die Klasse `replication` des Ports bleibt über `errors.Is` lesbar, das
// Ergebnis trägt keine Bestätigung. Der erste Aufruf gelingt, ab dem zweiten
// scheitert der Port.
func TestConfirmIdleForwardsPortErrorUnchanged(t *testing.T) {
	events := []string{}
	cause := stderrors.New("Standby-Status-Update fehlgeschlagen")
	portErr := fmt.Errorf("%w: %v", outbound.ErrReplication, cause)
	store := &fakeStore{events: &events}
	ack := &fakeAck{events: &events, ackErr: portErr, failFrom: 2}
	service := capture.NewCaptureService(store, ack)

	if _, err := service.ConfirmIdle(context.Background(), inbound.IdleConfirmationCommand{Position: idlePosition(t, 500)}); err != nil {
		t.Fatalf("erster Aufruf: %v", err)
	}
	result, err := service.ConfirmIdle(context.Background(), inbound.IdleConfirmationCommand{Position: idlePosition(t, 600)})
	if !stderrors.Is(err, outbound.ErrReplication) {
		t.Fatalf("Klasse replication fehlt: %v", err)
	}
	if !result.Acknowledged.IsZero() {
		t.Fatalf("Ergebnis trägt %v trotz Fehler, wollen keine Bestätigung", result.Acknowledged)
	}
	if got, want := events, []string{"ack:500", "ack:600"}; fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("Ereignisse = %v, wollen %v", got, want)
	}
}

// Eine Meldung ohne Position erreicht den Port nicht (`ADR-0120`): die
// Application bestätigt nie eine leere Position.
func TestConfirmIdleRejectsMissingPosition(t *testing.T) {
	service, events, _, _ := newService(t)

	_, err := service.ConfirmIdle(context.Background(), inbound.IdleConfirmationCommand{})
	if !stderrors.Is(err, capture.ErrMissingIdlePosition) {
		t.Fatalf("Fehler = %v, wollen ErrMissingIdlePosition", err)
	}
	if len(*events) != 0 {
		t.Fatalf("Ereignisse = %v, wollen keine", *events)
	}
}
