package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresack"
	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage"
	"github.com/pt9912/pg-change-feed/internal/adapters/driven/systemclock"
	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/mapper"
	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/receive"
	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/capture"
)

// TestRunStreamWithRetrySlotStillActive belegt den Fall „Slot noch aktiv"
// der Fitness Function (`ADR-0135`): der erste Lauf hält den Slot, der
// Retry-Zyklus scheitert am SQLSTATE 55006 (Transport-/Verbindungsstörung,
// wiederholbar), der Wartezug gibt den Slot frei und der zweite Versuch
// liefert die danach committete Change, Persist-before-ACK unberührt. Der Test
// bindet die Ursache des ersten Fehlschlags (SQLSTATE 55006), genau zwei
// Versuche, die Lieferung der Retry-Change hinter dem Slot-Stand vor dem
// Aufbau des zweiten Versuchs und genau zwei persistierte Changes (die
// Halter-Change und die Retry-Change je einmal); der Zyklus ist
// `runStreamCycle`. Die Zustellung ist At-least-once: der zweite Versuch darf
// die bereits persistierte Halter-Transaktion erneut liefern (die Persistierung
// ist idempotent), das färbt den Test nicht rot. Grenze: die Start-Position des
// Adapters ist nicht gebunden — setzt er vor `confirmed_flush_lsn` an, setzt
// der Server selbst dort an, der Test bleibt grün.
func TestRunStreamWithRetrySlotStillActive(t *testing.T) {
	dsn := os.Getenv("CDC_REPLICATION_TEST_DSN")
	if dsn == "" {
		t.Skip("CDC_REPLICATION_TEST_DSN nicht gesetzt — reale Verdrahtungs-Tests laufen über make test-replication")
	}
	const (
		source      = "src-retry"
		tableID     = "tbl-retry"
		schemaV     = "sv-retry"
		feed        = "public.feed_wire_retry_test"
		publication = "pub_pgc_test_wire_retry"
		slot        = "slot_pgc_test_wire_retry"
	)
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("Verbindungsaufbau: %v", err)
	}
	t.Cleanup(pool.Close)

	for _, statement := range []string{
		fmt.Sprintf("CREATE TABLE %s (id int PRIMARY KEY, name text)", feed),
		fmt.Sprintf("CREATE PUBLICATION %s FOR TABLE %s", publication, feed),
	} {
		if _, err := pool.Exec(ctx, statement); err != nil {
			t.Fatalf("Test-Umgebung: %v", err)
		}
	}
	t.Cleanup(func() {
		dropCtx, dropCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer dropCancel()
		deadline := time.Now().Add(8 * time.Second)
		for time.Now().Before(deadline) {
			if _, err := pool.Exec(dropCtx, fmt.Sprintf("SELECT pg_drop_replication_slot('%s')", slot)); err == nil {
				break
			}
			time.Sleep(200 * time.Millisecond)
		}
		for _, statement := range []string{
			"DROP PUBLICATION IF EXISTS " + publication,
			"DROP TABLE IF EXISTS " + feed,
		} {
			_, _ = pool.Exec(dropCtx, statement)
		}
	})

	for _, statement := range []string{
		fmt.Sprintf("INSERT INTO cdc.source (source_id, name) VALUES ('%s', 'Quelle')", source),
		fmt.Sprintf("INSERT INTO cdc.source_table (source_table_id, source_id, schema_name, table_name) VALUES ('%s', '%s', 'public', 'feed_wire_retry_test')", tableID, source),
		fmt.Sprintf("INSERT INTO cdc.schema_version (schema_version_id, source_table_id, version) VALUES ('%s', '%s', 1)", schemaV, tableID),
	} {
		if _, err := pool.Exec(ctx, statement); err != nil {
			t.Fatalf("Referenz-Zeilen: %v", err)
		}
	}
	store, err := postgresstorage.New(ctx, dsn)
	if err != nil {
		t.Fatalf("Store: %v", err)
	}
	t.Cleanup(store.Close)

	// Der erste Lauf hält den Slot und belegt die Erfassung bis zur
	// Persistierung — der Slot ist danach aktiv, solange dieser Lauf läuft.
	holderCtx, stopHolder := context.WithCancel(ctx)
	holderDone := make(chan struct{})
	holderStream, err := receive.NewStream(holderCtx, receive.Config{
		DSN:         dsn,
		Source:      source,
		Publication: publication,
		Slot:        slot,
		Tables: map[string]mapper.TableBinding{
			feed: {TableID: tableID, SchemaVersion: schemaV},
		},
	})
	if err != nil {
		t.Fatalf("erster Stream: %v", err)
	}
	holderAck, err := postgresack.New(holderStream.Conn())
	if err != nil {
		t.Fatalf("ACK-Adapter (Halter): %v", err)
	}
	holderService := capture.NewCaptureService(store, holderAck)
	if err := holderStream.BindCapture(holderService); err != nil {
		t.Fatalf("BindCapture (Halter): %v", err)
	}
	if err := holderStream.BindIdleConfirmation(holderService); err != nil {
		t.Fatalf("BindIdleConfirmation (Halter): %v", err)
	}
	go func() {
		defer close(holderDone)
		_ = holderStream.Run(holderCtx)
	}()
	if _, err := pool.Exec(ctx, "INSERT INTO "+feed+" (id, name) VALUES (1, 'Halter')"); err != nil {
		t.Fatalf("INSERT (Halter): %v", err)
	}
	holderDeadline := time.Now().Add(15 * time.Second)
	for {
		var count int
		if err := pool.QueryRow(ctx,
			"SELECT count(*) FROM cdc.change WHERE transaction_id IN (SELECT transaction_id FROM cdc.transaction WHERE source_id = $1)", source).Scan(&count); err != nil {
			t.Fatalf("cdc.change zählen: %v", err)
		}
		if count >= 1 {
			break
		}
		if time.Now().After(holderDeadline) {
			t.Fatalf("der Halter persistiert die Change nicht")
		}
		time.Sleep(100 * time.Millisecond)
	}

	var haltedPosition uint64
	if err := pool.QueryRow(ctx,
		"SELECT max(commit_position) FROM cdc.transaction WHERE source_id = $1", source).Scan(&haltedPosition); err != nil {
		t.Fatalf("Position der Halter-Transaktion: %v", err)
	}
	// Der Server nimmt die Bestätigung asynchron an: der Slot-Stand rückt
	// auf die Halter-Position, bevor der Retry beginnt.
	flushDeadline := time.Now().Add(10 * time.Second)
	for {
		flush, err := readConfirmedFlush(pool, slot)
		if err != nil {
			t.Fatalf("%v", err)
		}
		if flush >= haltedPosition {
			break
		}
		if time.Now().After(flushDeadline) {
			t.Fatalf("confirmed_flush_lsn erreicht die Halter-Position %d nicht", haltedPosition)
		}
		time.Sleep(50 * time.Millisecond)
	}

	// Der Retry: der Zyklus ist `runStreamCycle` (Aufbau unter einem Kontext,
	// der vor `Run` endet). Der erste Versuch scheitert am SQLSTATE 55006
	// (der Slot ist aktiv), der Wartezug gibt den Slot frei, und der zweite
	// Versuch liefert die danach committete Change. Der Mitschnitt hält je
	// Versuch den Fehler, den Slot-Stand vor dem Aufbau und die Lieferungen.
	recorder := &recordingCycleService{}
	ackPort := &streamCycleAck{}
	recorder.cycleService = capture.NewCaptureService(store, ackPort)
	var trace attemptTrace
	cycle := func(attemptCtx context.Context, streaming func()) error {
		flushAtStart, err := readConfirmedFlush(pool, slot)
		if err != nil {
			t.Errorf("%v", err)
			return err
		}
		attempt := trace.begin(flushAtStart)
		recorder.setAttempt(attempt)
		open := func(setupCtx context.Context) (cycleStream, error) {
			return receive.NewStream(setupCtx, receive.Config{
				DSN:         dsn,
				Source:      source,
				Publication: publication,
				Slot:        slot,
				OnStreaming: streaming,
				Tables: map[string]mapper.TableBinding{
					feed: {TableID: tableID, SchemaVersion: schemaV},
				},
			})
		}
		newAck := func(conn *pgconn.PgConn) (outbound.ReplicationAckPort, error) {
			return postgresack.New(conn)
		}
		err = runStreamCycle(attemptCtx, streamSetupTimeout, open, newAck, ackPort, recorder)
		trace.end(attempt, err)
		return err
	}

	var insertRetry sync.Once
	retryCtx, stopRetry := context.WithCancel(ctx)
	retryDone := make(chan error, 1)
	go func() {
		retryDone <- runStreamWithRetry(retryCtx, outbound.NoopLog, systemclock.New(), cycle, func(time.Duration) {
			// Der Wartezug gibt den Slot frei: der Halter endet regulär,
			// der Slot steht danach für den zweiten Versuch bereit. Die
			// Change „Retry" entsteht erst danach — kein Halter kann sie mehr
			// liefern.
			stopHolder()
			<-holderDone
			awaitSlotInactive(t, pool, slot, 10*time.Second)
			insertRetry.Do(func() {
				if _, err := pool.Exec(ctx, "INSERT INTO "+feed+" (id, name) VALUES (2, 'Retry')"); err != nil {
					t.Errorf("INSERT (Retry): %v", err)
				}
			})
		})
	}()

	// Die Change „Retry" erscheint nach dem erfolgreichen zweiten Versuch.
	retryDeadline := time.Now().Add(30 * time.Second)
	for {
		var count int
		if err := pool.QueryRow(ctx,
			"SELECT count(*) FROM cdc.change WHERE transaction_id IN (SELECT transaction_id FROM cdc.transaction WHERE source_id = $1)", source).Scan(&count); err != nil {
			t.Fatalf("cdc.change zählen: %v", err)
		}
		if count >= 2 {
			break
		}
		if time.Now().After(retryDeadline) {
			t.Fatalf("die Retry-Change wird nicht persistiert")
		}
		time.Sleep(100 * time.Millisecond)
	}
	// Nachlauf: eine doppelte oder weitere Lieferung erreicht in dieser Zeit
	// den Mitschnitt und die Zählung.
	time.Sleep(1500 * time.Millisecond)
	stopRetry()
	select {
	case err := <-retryDone:
		if err != nil {
			t.Fatalf("Retry-Zyklus: %v", err)
		}
	case <-time.After(10 * time.Second):
		t.Fatalf("Retry-Zyklus endet nach Kontext-Ende nicht")
	}

	// Persistiert: je eine Zeile für die Halter-Transaktion und für die
	// Retry-Change, geordnet nach der Commit-Position ihrer Transaktion.
	persisted, err := pool.Query(ctx,
		"SELECT t.commit_position, c.change_id FROM cdc.change c JOIN cdc.transaction t ON t.transaction_id = c.transaction_id WHERE t.source_id = $1 ORDER BY t.commit_position, c.sequence", source)
	if err != nil {
		t.Fatalf("cdc.change lesen: %v", err)
	}
	type persistedChange struct {
		position uint64
		changeID string
	}
	var changes []persistedChange
	for persisted.Next() {
		var position int64
		var changeID string
		if err := persisted.Scan(&position, &changeID); err != nil {
			t.Fatalf("cdc.change-Zeile: %v", err)
		}
		changes = append(changes, persistedChange{position: uint64(position), changeID: changeID})
	}
	persisted.Close()
	if err := persisted.Err(); err != nil {
		t.Fatalf("cdc.change lesen: %v", err)
	}
	if len(changes) != 2 {
		t.Fatalf("cdc.change-Zeilen = %v, erwartet genau 2 (Halter-Change und Retry-Change)", changes)
	}
	if changes[0].position != haltedPosition {
		t.Fatalf("erste Change auf Position %d, erwartet die Halter-Position %d", changes[0].position, haltedPosition)
	}
	retryPosition := changes[1].position
	if retryPosition <= haltedPosition {
		t.Fatalf("Retry-Change %s auf Position %d liegt nicht hinter der Halter-Position %d", changes[1].changeID, retryPosition, haltedPosition)
	}

	// Ursache des ersten Fehlschlags: SQLSTATE 55006 an START_REPLICATION,
	// wiederholbare Klasse `replication` ohne `ErrRejected`/`ErrPermission`.
	// `serverFault` trägt den SQLSTATE im Text der Ursache, nicht als
	// `*pgconn.PgError` in der Kette.
	attempts := trace.snapshot()
	if len(attempts) != 2 {
		t.Fatalf("Versuche = %d, erwartet genau 2", len(attempts))
	}
	first := attempts[0].err
	if first == nil || !strings.Contains(first.Error(), "SQLSTATE 55006") || !strings.Contains(first.Error(), "START_REPLICATION") {
		t.Fatalf("Ursache des ersten Fehlschlags = %v, erwartet START_REPLICATION mit SQLSTATE 55006", first)
	}
	if !errors.Is(first, receive.ErrReplication) || errors.Is(first, receive.ErrRejected) || errors.Is(first, receive.ErrPermission) {
		t.Fatalf("Klasse des ersten Fehlschlags = %v, erwartet wiederholbares ErrReplication", first)
	}
	if attempts[1].err != nil {
		t.Fatalf("zweiter Versuch endet mit Fehler: %v", attempts[1].err)
	}

	// Lieferungen: der erste Versuch liefert nichts; der zweite liefert die
	// Retry-Transaktion, hinter dem Slot-Stand vor seinem Aufbau. Die
	// Halter-Transaktion liegt auf oder vor diesem Stand und darf erneut
	// geliefert werden; jede andere Position ist unbekannt.
	delivered := recorder.snapshot()
	if len(delivered[1]) != 0 {
		t.Fatalf("Lieferungen des ersten Versuchs = %v, erwartet keine", delivered[1])
	}
	flushAtSecondStart := attempts[1].flushAtStart
	t.Logf("zweiter Versuch: Lieferungen %v, Halter-Position %d, Retry-Position %d, confirmed_flush_lsn vor dem Aufbau %d",
		delivered[2], haltedPosition, retryPosition, flushAtSecondStart)
	if haltedPosition > flushAtSecondStart {
		t.Fatalf("Halter-Position %d liegt hinter confirmed_flush_lsn %d vor dem zweiten Versuch", haltedPosition, flushAtSecondStart)
	}
	if retryPosition <= flushAtSecondStart {
		t.Fatalf("Retry-Position %d liegt nicht hinter confirmed_flush_lsn %d vor dem zweiten Versuch", retryPosition, flushAtSecondStart)
	}
	retryDeliveries := 0
	for _, position := range delivered[2] {
		switch position {
		case retryPosition:
			retryDeliveries++
		case haltedPosition:
		default:
			t.Fatalf("Lieferungen des zweiten Versuchs = %v, enthalten die unbekannte Position %d", delivered[2], position)
		}
	}
	if retryDeliveries == 0 {
		t.Fatalf("Lieferungen des zweiten Versuchs = %v, erwartet die Retry-Position %d", delivered[2], retryPosition)
	}
}

// readConfirmedFlush trägt `confirmed_flush_lsn` des Slots als Offset.
// Der Fehler geht an den Aufrufer zurück: die Funktion läuft auch in der
// Zyklus-Closure, also nicht in der Test-Goroutine.
func readConfirmedFlush(pool *pgxpool.Pool, slot string) (uint64, error) {
	var flush int64
	if err := pool.QueryRow(context.Background(),
		"SELECT (confirmed_flush_lsn - '0/0'::pg_lsn)::bigint FROM pg_replication_slots WHERE slot_name = $1", slot).Scan(&flush); err != nil {
		return 0, fmt.Errorf("confirmed_flush_lsn von %s: %w", slot, err)
	}
	return uint64(flush), nil
}

// awaitSlotInactive wartet, bis kein Prozess den Slot mehr hält; ein
// Versuch davor endet mit SQLSTATE 55006.
func awaitSlotInactive(t *testing.T, pool *pgxpool.Pool, slot string, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for {
		var active bool
		if err := pool.QueryRow(context.Background(),
			"SELECT active FROM pg_replication_slots WHERE slot_name = $1", slot).Scan(&active); err != nil {
			t.Errorf("Slot %s im Katalog: %v", slot, err)
			return
		}
		if !active {
			return
		}
		if time.Now().After(deadline) {
			t.Errorf("Slot %s ist %s nach dem Halter-Ende noch aktiv", slot, timeout)
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
}

// attemptRecord trägt einen Versuch des Zyklus: den Slot-Stand vor seinem
// Aufbau und sein Ende.
type attemptRecord struct {
	flushAtStart uint64
	err          error
}

// attemptTrace schneidet die Versuche des Zyklus mit.
type attemptTrace struct {
	mu       sync.Mutex
	attempts []attemptRecord
}

// begin legt einen Versuch an und trägt seine laufende Nummer (ab 1).
func (a *attemptTrace) begin(flushAtStart uint64) int {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.attempts = append(a.attempts, attemptRecord{flushAtStart: flushAtStart})
	return len(a.attempts)
}

// end trägt den Ausgang des Versuchs `attempt`.
func (a *attemptTrace) end(attempt int, err error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.attempts[attempt-1].err = err
}

func (a *attemptTrace) snapshot() []attemptRecord {
	a.mu.Lock()
	defer a.mu.Unlock()
	return append([]attemptRecord(nil), a.attempts...)
}

// recordingCycleService schneidet die Commit-Positionen der gelieferten
// Transaktionen je Versuch mit und reicht an den Capture Service weiter.
type recordingCycleService struct {
	cycleService

	mu        sync.Mutex
	attempt   int
	delivered map[int][]uint64
}

func (r *recordingCycleService) setAttempt(attempt int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.attempt = attempt
}

func (r *recordingCycleService) Capture(ctx context.Context, command inbound.CaptureCommand) (inbound.CaptureResult, error) {
	position, _ := command.Transaction.CommitPosition()
	r.mu.Lock()
	if r.delivered == nil {
		r.delivered = map[int][]uint64{}
	}
	r.delivered[r.attempt] = append(r.delivered[r.attempt], position.Offset)
	r.mu.Unlock()
	return r.cycleService.Capture(ctx, command)
}

func (r *recordingCycleService) snapshot() map[int][]uint64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := map[int][]uint64{}
	for attempt, positions := range r.delivered {
		out[attempt] = append([]uint64(nil), positions...)
	}
	return out
}
