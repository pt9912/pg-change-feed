package receive_test

import (
	"context"
	stderrors "errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pglogrepl"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/mapper"
	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/receive"
	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Die Stream-Tests laufen gegen eine reale PostgreSQL-Instanz mit
// Publication und Logical Replication Slot im Testcontainer (`ADR-0030`,
// `LH-FA-CFG-001.a`), gepinnt über `make test-replication`; ohne DSN
// überspringen sie — die Dekodier- und Mapper-Seite tragen die
// Unit-Tests gegen die `pgoutput`-Binärcodes. Die
// Verdrahtung mit ChangeStore und ACK-Adapter trägt der
// Verdrahtungs-Test in der Composition-Root (`internal/bootstrap`,
// `ADR-0026`). Die Instanz gehört dem Container: Tabellen, Publication
// und Slot werden je Test frisch aufgesetzt, Daten bleiben im Container
// und landen nicht im Arbeitsbaum.

const (
	testSource  = "src-1"
	testTableID = "tbl-1"
	testSchemaV = "sv-1"
	testFeed    = "public.feed_stream_test"
	testOther   = "public.other_stream_test"
)

// fakeCapture nimmt die Commands des Streams auf und meldet die
// Commit-Position als bestätigt — die Stand-in-Application des
// Adapter-Tests; die Persist-before-ACK-Ordnung am realen Treiber trägt
// der Verdrahtungs-Test in der Composition-Root. Mit ackFirstOnly
// bestätigt der Stand-in nur die erste Transaktion und mit declineIdle
// keine Leerlauf-Position — die Keepalive-Position-Regel (F-3) trägt der
// Test darüber: die bestätigte Position bleibt hinter dem Empfangsstand
// zurück. Ohne declineIdle bestätigt der Stand-in jede gemeldete
// Leerlauf-Position, wie der `ReplicationAckPort` der Application: das
// Standby-Status-Update geht über die Verbindung des Streams (`conn`) an
// den Slot.
type fakeCapture struct {
	commands     chan *inbound.CaptureCommand
	ackFirstOnly bool
	acked        bool
	declineIdle  bool
	conn         *pgconn.PgConn

	idleMu        sync.Mutex
	idleConfirmed []uint64
}

// ConfirmIdle bestätigt die gemeldete Position über die Verbindung des
// Streams, sofern der Stand-in nicht ablehnt, und hält sie fest.
func (f *fakeCapture) ConfirmIdle(ctx context.Context, command inbound.IdleConfirmationCommand) (inbound.IdleConfirmationResult, error) {
	if f.declineIdle {
		return inbound.IdleConfirmationResult{}, nil
	}
	lsn := pglogrepl.LSN(command.Position.Offset)
	err := pglogrepl.SendStandbyStatusUpdate(ctx, f.conn, pglogrepl.StandbyStatusUpdate{
		WALWritePosition: lsn,
		WALFlushPosition: lsn,
		WALApplyPosition: lsn,
	})
	if err != nil {
		return inbound.IdleConfirmationResult{}, err
	}
	f.idleMu.Lock()
	f.idleConfirmed = append(f.idleConfirmed, command.Position.Offset)
	f.idleMu.Unlock()
	return inbound.IdleConfirmationResult{Acknowledged: command.Position}, nil
}

// idleConfirmations trägt die Zahl der bestätigten Leerlauf-Positionen.
func (f *fakeCapture) idleConfirmations() int {
	f.idleMu.Lock()
	defer f.idleMu.Unlock()
	return len(f.idleConfirmed)
}

// Capture nimmt ein Command auf und meldet die Commit-Position als
// bestätigt — im ackFirstOnly-Modus nur für die erste Transaktion.
func (f *fakeCapture) Capture(_ context.Context, command inbound.CaptureCommand) (inbound.CaptureResult, error) {
	f.commands <- &command
	if f.ackFirstOnly && f.acked {
		return inbound.CaptureResult{}, nil
	}
	f.acked = true
	position, _ := command.Transaction.CommitPosition()
	return inbound.CaptureResult{Acknowledged: position}, nil
}

// newPool baut den normalen Verbindungspool gegen die Test-Instanz.
func newPool(t *testing.T) (*pgxpool.Pool, context.Context) {
	t.Helper()
	return newPoolOn(t, "CDC_REPLICATION_TEST_DSN")
}

// newPoolOn baut den Verbindungspool gegen die Instanz, deren DSN die
// Umgebungsvariable trägt; ohne DSN überspringt der Test.
func newPoolOn(t *testing.T, dsnVariable string) (*pgxpool.Pool, context.Context) {
	t.Helper()
	dsn := os.Getenv(dsnVariable)
	if dsn == "" {
		t.Skipf("%s nicht gesetzt — reale Replication-Tests laufen über make test-replication", dsnVariable)
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("Verbindungsaufbau: %v", err)
	}
	t.Cleanup(pool.Close)
	return pool, ctx
}

// testEnv trägt die Umgebung eines Tests: eigener Pool, eigene Tabelle,
// eigene Publication und eigener Slot (`LH-FA-CFG-001.a`), damit die
// Tests voneinander isoliert sind.
type testEnv struct {
	pool        *pgxpool.Pool
	ctx         context.Context
	publication string
	slot        string
}

// newTestEnv setzt Quelle, Tabellen und Publication je Test auf; der
// Slot wird mit dem Stream-Adapter angelegt. Die Publication trägt
// genau die Test-Tabellen — parallel laufende Test-Pakete teilen denselben
// Container, eine FOR ALL TABLES-Publication ließe jedes Test-Paket die
// Changes der anderen sehen. Der Cleanup stoppt den Laufkontext vor dem
// Slot-Rückbau — ein belegter Slot droppt nicht.
func newTestEnv(t *testing.T, name string) *testEnv {
	t.Helper()
	return newTestEnvOn(t, name, "CDC_REPLICATION_TEST_DSN")
}

// newTestEnvOn setzt die Test-Umgebung gegen die Instanz auf, deren DSN die
// Umgebungsvariable trägt.
func newTestEnvOn(t *testing.T, name, dsnVariable string) *testEnv {
	t.Helper()
	pool, ctx := newPoolOn(t, dsnVariable)
	runCtx, cancel := context.WithCancel(ctx)
	env := &testEnv{
		pool:        pool,
		ctx:         runCtx,
		publication: "pub_pgc_test_" + name,
		slot:        "slot_pgc_test_" + name,
	}
	for _, statement := range []string{
		fmt.Sprintf("CREATE TABLE %s (id int PRIMARY KEY, name text)", testFeed),
		fmt.Sprintf("CREATE TABLE %s (id int PRIMARY KEY)", testOther),
		fmt.Sprintf("CREATE PUBLICATION %s FOR TABLE %s, %s", env.publication, testFeed, testOther),
	} {
		if _, err := pool.Exec(ctx, statement); err != nil {
			t.Fatalf("Test-Umgebung: %v", err)
		}
	}
	t.Cleanup(func() {
		cancel()
		dropCtx, dropCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer dropCancel()
		deadline := time.Now().Add(8 * time.Second)
		for time.Now().Before(deadline) {
			_, err := pool.Exec(dropCtx, fmt.Sprintf("SELECT pg_drop_replication_slot('%s')", env.slot))
			if err == nil {
				break
			}
			time.Sleep(200 * time.Millisecond)
		}
		for _, statement := range []string{
			"DROP PUBLICATION IF EXISTS " + env.publication,
			"DROP TABLE IF EXISTS " + testFeed,
			"DROP TABLE IF EXISTS " + testOther,
		} {
			_, _ = pool.Exec(dropCtx, statement)
		}
	})
	return env
}

// startStream legt den Stream-Adapter an und läuft im Hintergrund; der
// Test liest über das Command-Feld, der zweite Ausgang trägt das
// Lauf-Ende (F-4: der Restart wartet auf das Verbindungs-Ende).
func startStream(t *testing.T, env *testEnv, capturePort *fakeCapture) (*receive.Stream, <-chan error) {
	t.Helper()
	runDone := make(chan error, 1)
	runCtx, cancel := context.WithCancel(env.ctx)
	t.Cleanup(cancel)
	stream, err := newStream(runCtx, env, capturePort)
	if err != nil {
		t.Fatalf("NewStream: %v", err)
	}
	go func() {
		runDone <- stream.Run(runCtx)
	}()
	return stream, runDone
}

// newStream legt den Stream-Adapter mit dem Stand-in als Capture- und
// Leerlauf-Port an und gibt dem Stand-in die Verbindung des Streams, über die
// er seine Bestätigungen sendet.
func newStream(ctx context.Context, env *testEnv, standIn *fakeCapture) (*receive.Stream, error) {
	stream, err := receive.NewStream(ctx, receive.Config{
		DSN:         env.pool.Config().ConnString(),
		Source:      testSource,
		Publication: env.publication,
		Slot:        env.slot,
		Tables: map[string]mapper.TableBinding{
			testFeed: {TableID: testTableID, SchemaVersion: testSchemaV},
		},
		Capture:          standIn,
		IdleConfirmation: standIn,
	})
	if err != nil {
		return nil, err
	}
	standIn.conn = stream.Conn()
	return stream, nil
}

// awaitCommand liest das nächste Command mit Test-Zeitgrenze.
func awaitCommand(t *testing.T, commands chan *inbound.CaptureCommand, timeout time.Duration) *inbound.CaptureCommand {
	t.Helper()
	select {
	case command := <-commands:
		return command
	case <-time.After(timeout):
		t.Fatalf("kein CaptureCommand innerhalb %s", timeout)
		return nil
	}
}

// awaitNoCommand meldet, dass innerhalb der Zeitspanne kein Command
// auftritt.
func awaitNoCommand(t *testing.T, commands chan *inbound.CaptureCommand, duration time.Duration) {
	t.Helper()
	select {
	case command := <-commands:
		t.Fatalf("unerwartetes CaptureCommand: %+v", command)
	case <-time.After(duration):
	}
}

// TestStreamTranslatesRealChanges trägt die Translation am realen
// Pfad: INSERT, UPDATE und DELETE auf der aktivierten Tabelle erzeugen
// committed Quelltransaktionen mit ihren Row Images; die Transaktion
// über die nicht aktivierte Tabelle kommt ohne Change
// (`LH-FA-CAP-001`…003, `LH-FA-CFG-001`).
func TestStreamTranslatesRealChanges(t *testing.T) {
	pool, ctx := newPool(t)
	env := newTestEnv(t, "changes")
	commands := make(chan *inbound.CaptureCommand, 32)
	_, _ = startStream(t, env, &fakeCapture{commands: commands})

	if _, err := pool.Exec(ctx, "INSERT INTO "+testFeed+" (id, name) VALUES (1, 'Alpha')"); err != nil {
		t.Fatalf("INSERT: %v", err)
	}
	if _, err := pool.Exec(ctx, "INSERT INTO "+testOther+" (id) VALUES (9)"); err != nil {
		t.Fatalf("INSERT andere Tabelle: %v", err)
	}
	if _, err := pool.Exec(ctx, "UPDATE "+testFeed+" SET name = 'Beta' WHERE id = 1"); err != nil {
		t.Fatalf("UPDATE: %v", err)
	}
	if _, err := pool.Exec(ctx, "DELETE FROM "+testFeed+" WHERE id = 1"); err != nil {
		t.Fatalf("DELETE: %v", err)
	}

	// Die vier Quelltransaktionen laufen ein; die Zuordnung läuft über
	// den Inhalt, die Ordnung der Feed-Transaktionen über ihre
	// Commit-Positionen (`LH-FA-CAP-004`).
	var delivered []*inbound.CaptureCommand
	for i := 0; i < 4; i++ {
		delivered = append(delivered, awaitCommand(t, commands, 15*time.Second))
	}

	var emptyTransaction *inbound.CaptureCommand
	feedTransactions := map[model.Operation]*inbound.CaptureCommand{}
	for _, command := range delivered {
		changes, err := command.Transaction.Changes()
		if err != nil {
			t.Fatalf("Changes: %v", err)
		}
		if len(changes) == 0 {
			emptyTransaction = command
			continue
		}
		feedTransactions[changes[0].Operation] = command
	}
	if emptyTransaction == nil {
		t.Fatalf("Transaktion der nicht aktivierten Tabelle fehlt")
	}
	for _, operation := range []model.Operation{model.OperationInsert, model.OperationUpdate, model.OperationDelete} {
		if feedTransactions[operation] == nil {
			t.Fatalf("Transaktion mit %s fehlt", operation)
		}
	}
	insertChanges, err := feedTransactions[model.OperationInsert].Transaction.Changes()
	if err != nil {
		t.Fatalf("Insert-Changes: %v", err)
	}
	if len(insertChanges) != 1 {
		t.Fatalf("Insert-Change-Anzahl: %d", len(insertChanges))
	}
	if string(insertChanges[0].NewImage) != `{"id":"1","name":"Alpha"}` {
		t.Fatalf("Insert-Image: %s", insertChanges[0].NewImage)
	}
	updateChanges, err := feedTransactions[model.OperationUpdate].Transaction.Changes()
	if err != nil {
		t.Fatalf("Update-Changes: %v", err)
	}
	if len(updateChanges) != 1 || updateChanges[0].Operation != model.OperationUpdate {
		t.Fatalf("Update-Change: %+v", updateChanges)
	}
	deleteChanges, err := feedTransactions[model.OperationDelete].Transaction.Changes()
	if err != nil {
		t.Fatalf("Delete-Changes: %v", err)
	}
	if len(deleteChanges) != 1 {
		t.Fatalf("Delete-Change-Anzahl: %d", len(deleteChanges))
	}
	if string(deleteChanges[0].OldImage) != `{"id":"1"}` {
		t.Fatalf("Delete-Alt-Image: %s", deleteChanges[0].OldImage)
	}

	// Die Commit-Positionen der Feed-Transaktionen laufen in
	// Commit-Reihenfolge (`LH-FA-CAP-004`, `LH-FA-DAT-004`).
	var feedPosition []model.SourcePosition
	for _, command := range delivered {
		position, committed := command.Transaction.CommitPosition()
		if !committed {
			t.Fatalf("Transaktion ohne Commit-Position")
		}
		if _, err := command.Transaction.Changes(); err != nil {
			t.Fatalf("Changes: %v", err)
		}
		changes, _ := command.Transaction.Changes()
		if len(changes) > 0 {
			feedPosition = append(feedPosition, position)
		}
	}
	for i := 1; i < len(feedPosition); i++ {
		if !feedPosition[i-1].Before(feedPosition[i]) {
			t.Fatalf("Commit-Positionen: %v läuft nicht vor %v", feedPosition[i-1], feedPosition[i])
		}
	}
}

// TestStreamOpenTransactionNotConsumable trägt die
// Transaktions-Sichtbarkeit am realen Pfad (`LH-FA-CAP-006.a`): eine
// offene Quelltransaktion liefert kein Command; erst der Commit meldet
// die Transaktion konsumierbar.
func TestStreamOpenTransactionNotConsumable(t *testing.T) {
	pool, ctx := newPool(t)
	env := newTestEnv(t, "visibility")
	commands := make(chan *inbound.CaptureCommand, 32)
	_, _ = startStream(t, env, &fakeCapture{commands: commands})

	sourceTx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatalf("Quelltransaktion: %v", err)
	}
	// Der Rollback-Grenze: ein Test-Abbruch vor dem Commit gibt die
	// Pool-Verbindung wieder frei (pgx-Tx lebt über den Test-Ausgang hinaus).
	t.Cleanup(func() { _ = sourceTx.Rollback(context.Background()) })
	if _, err := sourceTx.Exec(ctx, "INSERT INTO "+testFeed+" (id, name) VALUES (2, 'Offen')"); err != nil {
		t.Fatalf("INSERT: %v", err)
	}
	awaitNoCommand(t, commands, 2*time.Second)
	if err := sourceTx.Commit(ctx); err != nil {
		t.Fatalf("Commit: %v", err)
	}
	command := awaitCommand(t, commands, 15*time.Second)
	changes, err := command.Transaction.Changes()
	if err != nil {
		t.Fatalf("Changes: %v", err)
	}
	if len(changes) != 1 || string(changes[0].NewImage) != `{"id":"2","name":"Offen"}` {
		t.Fatalf("Transaktion nach Commit: %+v", changes)
	}
}

// TestStreamTruncateUnsupported trägt die TRUNCATE-Erkennung am realen
// Pfad: die Operation endet als sichtbarer Fehler der
// Nicht-Unterstützung (`LH-FA-CFG-001.a`).
func TestStreamTruncateUnsupported(t *testing.T) {
	pool, ctx := newPool(t)
	env := newTestEnv(t, "truncate")
	commands := make(chan *inbound.CaptureCommand, 32)
	stream, err := newStream(env.ctx, env, &fakeCapture{commands: commands})
	if err != nil {
		t.Fatalf("NewStream: %v", err)
	}
	runDone := make(chan error, 1)
	go func() { runDone <- stream.Run(env.ctx) }()

	if _, err := pool.Exec(ctx, "TRUNCATE "+testFeed); err != nil {
		t.Fatalf("TRUNCATE: %v", err)
	}
	select {
	case err := <-runDone:
		if !stderrors.Is(err, mapper.ErrTruncateUnsupported) {
			t.Fatalf("TRUNCATE-Ausgang: %v", err)
		}
	case <-time.After(15 * time.Second):
		t.Fatalf("TRUNCATE endete nicht als Stream-Fehler")
	}
}

// readConfirmedFlush trägt den confirmed_flush_lsn-Stand eines Slots
// (pg_replication_slots) als Offset.
func readConfirmedFlush(t *testing.T, pool *pgxpool.Pool, slot string) uint64 {
	t.Helper()
	values, err := pool.Query(context.Background(),
		"SELECT confirmed_flush_lsn::text FROM pg_replication_slots WHERE slot_name = $1", slot)
	if err != nil {
		t.Fatalf("Slot-Stand: %v", err)
	}
	defer values.Close()
	if !values.Next() {
		t.Fatalf("Slot %q fehlt", slot)
	}
	var confirmed string
	if err := values.Scan(&confirmed); err != nil {
		t.Fatalf("Slot-Stand lesen: %v", err)
	}
	lsn, err := pglogrepl.ParseLSN(confirmed)
	if err != nil {
		t.Fatalf("confirmed_flush_lsn %q: %v", confirmed, err)
	}
	return uint64(lsn)
}

// TestStreamKeepaliveReportsAcknowledgedPosition trägt die
// Keepalive-Position-Regel am realen Pfad (`LH-QA-REL-001.a`): bestätigt die
// Application keine Leerlauf-Position (`declineIdle`), meldet die
// Keepalive-Antwort ausschließlich die vom Capture-Ergebnis bestätigte
// Position — confirmed_flush_lsn rückt auf sie und läuft nie über sie
// hinaus, obwohl der Stream weiteren WAL-Stand empfangen hat. Die Bestätigung
// im Leerlauf trägt `TestStreamIdleConfirmationReleasesForeignWAL`. Der
// Testcontainer trägt wal_sender_timeout=2s: der Walsender verlangt die
// Antwort nach der Hälfte der Zeit, die Keepalive-Antwort trägt lastAcked.
func TestStreamKeepaliveReportsAcknowledgedPosition(t *testing.T) {
	pool, ctx := newPool(t)
	env := newTestEnv(t, "keepalive")
	commands := make(chan *inbound.CaptureCommand, 32)
	startStream(t, env, &fakeCapture{commands: commands, ackFirstOnly: true, declineIdle: true})

	if _, err := pool.Exec(ctx, "INSERT INTO "+testFeed+" (id, name) VALUES (1, 'Bestätigt')"); err != nil {
		t.Fatalf("INSERT: %v", err)
	}
	first := awaitCommand(t, commands, 15*time.Second)
	ackedPosition, committed := first.Transaction.CommitPosition()
	if !committed {
		t.Fatalf("bestätigte Transaktion ohne Commit-Position")
	}

	// Weitere WAL-Ereignisse laufen ein und werden empfangen, aber nicht
	// bestätigt: die zweite Transaktion läuft über die nicht aktivierte
	// Tabelle (leere Transaktion), die dritte über die aktivierte. Der
	// Empfangsstand liegt damit hinter der bestätigten Position.
	if _, err := pool.Exec(ctx, "INSERT INTO "+testOther+" (id) VALUES (9)"); err != nil {
		t.Fatalf("INSERT andere Tabelle: %v", err)
	}
	if _, err := pool.Exec(ctx, "INSERT INTO "+testFeed+" (id, name) VALUES (2, 'Unbestätigt')"); err != nil {
		t.Fatalf("INSERT: %v", err)
	}
	awaitCommand(t, commands, 15*time.Second)
	unacked := awaitCommand(t, commands, 15*time.Second)
	unackedPosition, committed := unacked.Transaction.CommitPosition()
	if !committed || unackedPosition.Offset <= ackedPosition.Offset {
		t.Fatalf("Empfangsstand läuft nicht hinter der bestätigten Position")
	}

	// Die Keepalive-Antwort trägt die bestätigte Position: der
	// confirmed_flush_lsn rückt auf sie und nie über sie hinaus — eine
	// Antwort mit dem Empfangsstand (Mutation) würde confirmed_flush_lsn
	// über die bestätigte Position schieben.
	deadline := time.Now().Add(20 * time.Second)
	for {
		confirmed := readConfirmedFlush(t, pool, env.slot)
		if confirmed > ackedPosition.Offset {
			t.Fatalf("confirmed_flush_lsn %x läuft über die bestätigte Position %x hinaus", confirmed, ackedPosition.Offset)
		}
		if confirmed == ackedPosition.Offset {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("confirmed_flush_lsn erreicht die bestätigte Position nicht (%x != %x)", confirmed, ackedPosition.Offset)
		}
		time.Sleep(200 * time.Millisecond)
	}
}

// TestWALRetentionMeasuresGrowingBytes trägt die reale Byte-Differenz
// zwischen dem aktuellen WAL-Schreibstand und `confirmed_flush_lsn` eines
// inaktiven Slots (`SPEC-009` `cdc_wal_retention_bytes`, `ADR-0049`): der
// Stream bestätigt eine erste Transaktion und endet danach — der Slot
// bleibt mit seinem Bestand liegen (derselbe inaktive Zustand wie in
// `TestStreamRestartsOnExistingSlot`). Weitere, unbestätigte Transaktionen
// lassen den gemessenen Rückstand real wachsen — eine vertauschte
// Subtraktion in `WALRetentionChecker.Measure` (Mutation) würde hier einen
// fallenden oder negativen Wert liefern statt eines steigenden.
func TestWALRetentionMeasuresGrowingBytes(t *testing.T) {
	pool, ctx := newPool(t)
	env := newTestEnv(t, "retention")
	commands := make(chan *inbound.CaptureCommand, 32)

	runCtx, cancelRun := context.WithCancel(ctx)
	stream, err := newStream(runCtx, env, &fakeCapture{commands: commands})
	if err != nil {
		t.Fatalf("NewStream: %v", err)
	}
	runDone := make(chan error, 1)
	go func() { runDone <- stream.Run(runCtx) }()

	if _, err := pool.Exec(ctx, "INSERT INTO "+testFeed+" (id, name) VALUES (1, 'Bestätigt')"); err != nil {
		t.Fatalf("INSERT: %v", err)
	}
	awaitCommand(t, commands, 15*time.Second)

	// confirmed_flush_lsn erreicht die bestätigte Position, bevor der Slot
	// inaktiv wird — sonst trüge die spätere Messung einen zufällig
	// kleineren Anfangsstand.
	deadline := time.Now().Add(15 * time.Second)
	for readConfirmedFlush(t, pool, env.slot) == 0 {
		if time.Now().After(deadline) {
			t.Fatalf("confirmed_flush_lsn bleibt 0")
		}
		time.Sleep(200 * time.Millisecond)
	}

	// Der Stream endet — der Slot bleibt mit seinem Bestand liegen
	// (inaktiv, wie `TestStreamRestartsOnExistingSlot`).
	cancelRun()
	select {
	case <-runDone:
	case <-time.After(15 * time.Second):
		t.Fatalf("Stream endet nicht")
	}

	checker, err := receive.NewWALRetentionChecker(context.Background(), env.pool.Config().ConnString(), env.slot)
	if err != nil {
		t.Fatalf("NewWALRetentionChecker: %v", err)
	}
	t.Cleanup(func() { _ = checker.Close(context.Background()) })

	before, err := checker.Measure(context.Background())
	if err != nil {
		t.Fatalf("Measure (vorher): %v", err)
	}

	// Weitere, unbestätigte Transaktionen — der Slot bleibt inaktiv, der
	// WAL-Rückstand wächst real.
	for i := 0; i < 20; i++ {
		if _, err := pool.Exec(ctx, "INSERT INTO "+testFeed+" (id, name) VALUES ($1, repeat('x', 4000))", 100+i); err != nil {
			t.Fatalf("INSERT (Rückstand): %v", err)
		}
	}

	after, err := checker.Measure(context.Background())
	if err != nil {
		t.Fatalf("Measure (nachher): %v", err)
	}
	if after <= before {
		t.Fatalf("WAL-Rückstand ist nicht real gestiegen: vorher %d, nachher %d", before, after)
	}
}

// TestWALRetentionMeasureInvalidSlotName trägt einen der drei
// Fehlerpfade aus dem Review zu `slice-025` F-1 (ungültiger Slot):
// `NewWALRetentionChecker` weist einen Slot-Namen außerhalb des
// Bezeichner-Alphabets ab, bevor überhaupt eine Verbindung versucht wird
// — kein Testcontainer nötig, dieser Zweig läuft vor jeder DB-Interaktion.
func TestWALRetentionMeasureInvalidSlotName(t *testing.T) {
	_, err := receive.NewWALRetentionChecker(context.Background(), "postgres://ignored/db", "Ungueltiger Slot!")
	if !stderrors.Is(err, receive.ErrConfiguration) {
		t.Fatalf("NewWALRetentionChecker mit ungültigem Slot-Namen = %v, wollen ErrConfiguration", err)
	}
}

// TestWALRetentionMeasureConnectionRefusedFails trägt den
// Verbindungsfehler-Pfad aus dem Review zu `slice-025` F-1: ein nicht
// erreichbarer Host scheitert am Verbindungsaufbau selbst, sichtbar als
// `ErrReplication` — kein Testcontainer nötig, der Verbindungsversuch
// scheitert bereits am Transport (Port 1 trägt keinen Listener).
func TestWALRetentionMeasureConnectionRefusedFails(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := receive.NewWALRetentionChecker(ctx,
		"postgres://cdc:cdc@127.0.0.1:1/cdc_test?sslmode=disable&connect_timeout=2", "slot_pgc_test_unreachable")
	if !stderrors.Is(err, receive.ErrReplication) {
		t.Fatalf("NewWALRetentionChecker gegen nicht erreichbaren Host = %v, wollen ErrReplication", err)
	}
}

// TestWALRetentionMeasureMissingSlot trägt den `!exists`-Fehlerpfad aus
// dem Review zu `slice-025` F-1 (fehlender Slot): ein syntaktisch gültiger,
// aber nie angelegter Slot-Name liefert einen sichtbaren
// `ErrReplication`-Fehler statt eines stillen Nullwerts.
func TestWALRetentionMeasureMissingSlot(t *testing.T) {
	_, ctx := newPool(t)
	dsn := os.Getenv("CDC_REPLICATION_TEST_DSN")
	checker, err := receive.NewWALRetentionChecker(ctx, dsn, "slot_pgc_test_missing_retention")
	if err != nil {
		t.Fatalf("NewWALRetentionChecker: %v", err)
	}
	t.Cleanup(func() { _ = checker.Close(context.Background()) })

	if _, err := checker.Measure(ctx); !stderrors.Is(err, receive.ErrReplication) {
		t.Fatalf("Measure gegen fehlenden Slot = %v, wollen ErrReplication", err)
	}
}

// dsnWithApplicationName markiert eine Verbindungszeichenkette mit einem
// `application_name`, damit ein Test das serverseitige Backend eindeutig
// über `pg_stat_activity` wiederfindet — unabhängig davon, ob die DSN als
// URL oder als Schlüssel/Wert-Zeichenkette vorliegt (beide Formen
// akzeptiert `pgconn.ParseConfig`).
func dsnWithApplicationName(dsn, name string) string {
	if strings.HasPrefix(dsn, "postgres://") || strings.HasPrefix(dsn, "postgresql://") {
		sep := "?"
		if strings.Contains(dsn, "?") {
			sep = "&"
		}
		return dsn + sep + "application_name=" + name
	}
	return dsn + " application_name=" + name
}

// countBackends trägt die Anzahl aktiver Backends mit dem gegebenen
// `application_name`.
func countBackends(t *testing.T, pool *pgxpool.Pool, applicationName string) int {
	t.Helper()
	var count int
	if err := pool.QueryRow(context.Background(),
		"SELECT count(*) FROM pg_stat_activity WHERE application_name = $1", applicationName).Scan(&count); err != nil {
		t.Fatalf("pg_stat_activity zählen: %v", err)
	}
	return count
}

// terminateBackend beendet serverseitig die Backends mit dem gegebenen
// `application_name` und wartet, bis sie tatsächlich beendet sind —
// simuliert einen dauerhaften Verbindungsabbruch (Idle-Timeout zwischen
// zwei Ticks) für den Reconnect-Test unten.
func terminateBackend(t *testing.T, pool *pgxpool.Pool, applicationName string) {
	t.Helper()
	ctx := context.Background()
	if _, err := pool.Exec(ctx,
		"SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE application_name = $1", applicationName); err != nil {
		t.Fatalf("pg_terminate_backend: %v", err)
	}
	deadline := time.Now().Add(10 * time.Second)
	for countBackends(t, pool, applicationName) > 0 {
		if time.Now().After(deadline) {
			t.Fatalf("Backend %q endet nicht nach pg_terminate_backend", applicationName)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// TestWALRetentionMeasureReconnectsAfterConnectionLoss trägt sowohl den
// IDENTIFY_SYSTEM-Fehlerpfad aus dem Review zu `slice-025` F-1 als auch den
// Reconnect-Pfad aus F-2: Ein serverseitig beendetes Backend (simuliert
// einen dauerhaften Idle-Timeout zwischen zwei Ticks, anders als die
// Stream-Verbindung mit ihrem Keepalive-Verkehr) lässt die laufende
// Messung sichtbar scheitern; die eigene Verbindung des
// `WALRetentionChecker` wird dabei intern ersetzt (`reconnectAfterError`)
// — der übernächste Aufruf (der nächste Tick in `runWALRetentionCheck`)
// misst wieder erfolgreich, statt dass die Metrik für den Rest des
// Prozesslaufs verstummt.
func TestWALRetentionMeasureReconnectsAfterConnectionLoss(t *testing.T) {
	pool, ctx := newPool(t)
	env := newTestEnv(t, "reconnect")
	commands := make(chan *inbound.CaptureCommand, 32)

	// Der Slot muss bestehen, bevor der Checker misst — derselbe
	// kurzlebige Stream-Aufbau wie in TestWALRetentionMeasuresGrowingBytes.
	runCtx, cancelRun := context.WithCancel(ctx)
	stream, err := newStream(runCtx, env, &fakeCapture{commands: commands})
	if err != nil {
		t.Fatalf("NewStream: %v", err)
	}
	runDone := make(chan error, 1)
	go func() { runDone <- stream.Run(runCtx) }()

	if _, err := pool.Exec(ctx, "INSERT INTO "+testFeed+" (id, name) VALUES (1, 'Bestätigt')"); err != nil {
		t.Fatalf("INSERT: %v", err)
	}
	awaitCommand(t, commands, 15*time.Second)
	cancelRun()
	select {
	case <-runDone:
	case <-time.After(15 * time.Second):
		t.Fatalf("Stream endet nicht")
	}

	const applicationName = "walretention_reconnect_test"
	checker, err := receive.NewWALRetentionChecker(context.Background(),
		dsnWithApplicationName(env.pool.Config().ConnString(), applicationName), env.slot)
	if err != nil {
		t.Fatalf("NewWALRetentionChecker: %v", err)
	}
	t.Cleanup(func() { _ = checker.Close(context.Background()) })

	if _, err := checker.Measure(context.Background()); err != nil {
		t.Fatalf("Measure (vor Verbindungsabbruch): %v", err)
	}
	if countBackends(t, pool, applicationName) == 0 {
		t.Fatalf("Checker-Backend %q nicht in pg_stat_activity sichtbar", applicationName)
	}

	terminateBackend(t, pool, applicationName)

	measureCtx, cancelMeasure := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelMeasure()
	if _, err := checker.Measure(measureCtx); !stderrors.Is(err, receive.ErrReplication) {
		t.Fatalf("Measure nach Verbindungsabbruch = %v, wollen ErrReplication", err)
	}

	if _, err := checker.Measure(context.Background()); err != nil {
		t.Fatalf("Measure nach Reconnect: %v", err)
	}
}

// awaitSlotInactive wartet, bis der Walsender des Slots ihn freigegeben hat:
// das Ende von `Stream.Run` schließt die Verbindung des Clients, der Server
// gibt den Slot asynchron dazu frei (`active = false`); ein Wiederaufsetzen
// davor endet mit SQLSTATE 55006.
func awaitSlotInactive(t *testing.T, pool *pgxpool.Pool, slot string, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for {
		var active bool
		if err := pool.QueryRow(context.Background(),
			"SELECT active FROM pg_replication_slots WHERE slot_name = $1", slot).Scan(&active); err != nil {
			t.Fatalf("Slot %s im Katalog: %v", slot, err)
		}
		if !active {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("Slot %s ist %s nach dem Stream-Ende noch aktiv", slot, timeout)
		}
		time.Sleep(20 * time.Millisecond)
	}
}

// TestStreamRestartsOnExistingSlot trägt den bestehende-Slot-Zweig am
// realen Pfad (`ADR-0012`): der Stream endet, weitere Changes entstehen,
// der Restart legt denselben Slot wieder auf — `ensureSlot` liest den
// Slot-Bestand (`confirmed_flush_lsn`) statt still neu anzulegen, und
// der Stream setzt dort fort; die wiederholte Lieferung bestätigter
// Transaktionen ist der At-Least-Once-Fall (`ADR-0011`, `ADR-0012`).
func TestStreamRestartsOnExistingSlot(t *testing.T) {
	pool, ctx := newPool(t)
	env := newTestEnv(t, "restart")

	ctx1, cancel1 := context.WithCancel(ctx)
	defer cancel1()
	commands1 := make(chan *inbound.CaptureCommand, 32)
	stream1, err := newStream(ctx1, env, &fakeCapture{commands: commands1})
	if err != nil {
		t.Fatalf("NewStream (erster Lauf): %v", err)
	}
	runDone1 := make(chan error, 1)
	go func() { runDone1 <- stream1.Run(ctx1) }()

	if _, err := pool.Exec(ctx, "INSERT INTO "+testFeed+" (id, name) VALUES (1, 'Erste')"); err != nil {
		t.Fatalf("INSERT: %v", err)
	}
	firstCommand := awaitCommand(t, commands1, 15*time.Second)
	firstPosition, committed := firstCommand.Transaction.CommitPosition()
	if !committed {
		t.Fatalf("erste Transaktion ohne Commit-Position")
	}

	// Der erste Lauf endet regulär; die Verbindung schließt, der Slot
	// bleibt mit seinem Bestand zurück.
	cancel1()
	select {
	case err := <-runDone1:
		if err != nil {
			t.Fatalf("erster Lauf: %v", err)
		}
	case <-time.After(15 * time.Second):
		t.Fatalf("erster Lauf endet nicht")
	}
	awaitSlotInactive(t, pool, env.slot, 10*time.Second)

	if _, err := pool.Exec(ctx, "INSERT INTO "+testFeed+" (id, name) VALUES (2, 'Zweite')"); err != nil {
		t.Fatalf("INSERT nach Stream-Ende: %v", err)
	}

	commands2 := make(chan *inbound.CaptureCommand, 32)
	ctx2, cancel2 := context.WithCancel(ctx)
	t.Cleanup(cancel2)
	stream2, err := newStream(ctx2, env, &fakeCapture{commands: commands2})
	if err != nil {
		t.Fatalf("NewStream (Restart): %v", err)
	}
	runDone2 := make(chan error, 1)
	go func() { runDone2 <- stream2.Run(ctx2) }()

	// Der Restart liefert die Change nach dem Stream-Ende; die
	// bestätigte Transaktion kann als Wiederholung erneut kommen
	// (At-Least-Once, `ADR-0012`).
	deadline := time.Now().Add(20 * time.Second)
	for {
		command := awaitCommand(t, commands2, time.Until(deadline))
		position, committed := command.Transaction.CommitPosition()
		if !committed {
			t.Fatalf("Transaktion ohne Commit-Position")
		}
		changes, err := command.Transaction.Changes()
		if err != nil {
			t.Fatalf("Changes: %v", err)
		}
		if len(changes) == 1 && string(changes[0].NewImage) == `{"id":"2","name":"Zweite"}` {
			if position.Offset <= firstPosition.Offset {
				t.Fatalf("Restart-Position %x läuft nicht hinter der ersten Position %x", position.Offset, firstPosition.Offset)
			}
			return
		}
		if len(changes) == 1 && string(changes[0].NewImage) == `{"id":"1","name":"Erste"}` {
			// Wiederholung der bestätigten Transaktion: der
			// At-Least-Once-Fall, die Deduplizierung trägt der
			// Store (`ADR-0011`).
			if position.Offset != firstPosition.Offset {
				t.Fatalf("Wiederholte Transaktion trägt andere Position %x (erste %x)", position.Offset, firstPosition.Offset)
			}
			continue
		}
		t.Fatalf("unerwartete Transaktion: %+v", changes)
	}
}

// testForeign trägt die Tabelle außerhalb der Publication: ihr WAL trägt
// keinen Inhalt für den Stream (`ADR-0120`).
const testForeign = "public.foreign_stream_test"

// createForeignTable legt die Tabelle außerhalb der Publication an und
// räumt sie am Test-Ende ab.
func createForeignTable(t *testing.T, env *testEnv) {
	t.Helper()
	if _, err := env.pool.Exec(context.Background(),
		fmt.Sprintf("CREATE TABLE %s (id bigint PRIMARY KEY, payload text)", testForeign)); err != nil {
		t.Fatalf("Tabelle außerhalb der Publication: %v", err)
	}
	t.Cleanup(func() {
		_, _ = env.pool.Exec(context.Background(), "DROP TABLE IF EXISTS "+testForeign)
	})
}

// walLSNDiff trägt die Differenz zweier WAL-Positionen in Bytes.
func walLSNDiff(t *testing.T, pool *pgxpool.Pool, from, to string) int64 {
	t.Helper()
	var diff int64
	if err := pool.QueryRow(context.Background(),
		"SELECT pg_wal_lsn_diff($1::pg_lsn, $2::pg_lsn)::bigint", to, from).Scan(&diff); err != nil {
		t.Fatalf("pg_wal_lsn_diff: %v", err)
	}
	return diff
}

// currentWALLSN trägt die aktuelle WAL-Schreibposition der Instanz.
func currentWALLSN(t *testing.T, pool *pgxpool.Pool) string {
	t.Helper()
	var lsn string
	if err := pool.QueryRow(context.Background(), "SELECT pg_current_wal_lsn()::text").Scan(&lsn); err != nil {
		t.Fatalf("pg_current_wal_lsn: %v", err)
	}
	return lsn
}

// foreignWALRun trägt das Ergebnis einer Schreiblast außerhalb der
// Publication gegen einen laufenden Stream.
type foreignWALRun struct {
	load      int64
	backlog   int64
	confirmed uint64
	end       uint64
	elapsed   time.Duration
}

// runForeignWAL führt die Form X1 des Verdikts gegen einen laufenden Stream
// aus: 200.000 Zeilen in **eine** Transaktion auf eine Tabelle außerhalb der
// Publication. `load` ist die WAL-Menge der Last (Differenz der
// WAL-Schreibpositionen um den INSERT), `end` die WAL-Schreibposition hinter
// ihr. Gepollt wird bis `wait` verstreicht oder die Bedingung gilt: mit
// `byPosition` liegt `confirmed_flush_lsn` bei oder hinter `end` — unabhängig
// von WAL, das andere Schreiber der Instanz erzeugen —, sonst liegt der mit
// `WALRetentionChecker.Measure` gemessene Rückstand unter einem Zehntel der
// Last (nur auf einer Instanz ohne fremde Schreiber aussagekräftig).
func runForeignWAL(t *testing.T, dsnVariable, name string, standIn *fakeCapture, wait time.Duration, byPosition bool) foreignWALRun {
	t.Helper()
	env := newTestEnvOn(t, name, dsnVariable)
	createForeignTable(t, env)
	standIn.commands = make(chan *inbound.CaptureCommand, 32)
	startStream(t, env, standIn)

	// Der Stream läuft, sobald die erste Transaktion der aktivierten Tabelle
	// eingetroffen ist.
	if _, err := env.pool.Exec(env.ctx, "INSERT INTO "+testFeed+" (id, name) VALUES (1, 'Start')"); err != nil {
		t.Fatalf("INSERT: %v", err)
	}
	awaitCommand(t, standIn.commands, 15*time.Second)

	checker, err := receive.NewWALRetentionChecker(context.Background(), env.pool.Config().ConnString(), env.slot)
	if err != nil {
		t.Fatalf("NewWALRetentionChecker: %v", err)
	}
	t.Cleanup(func() { _ = checker.Close(context.Background()) })

	before := currentWALLSN(t, env.pool)
	if _, err := env.pool.Exec(env.ctx,
		fmt.Sprintf("INSERT INTO %s SELECT g, repeat('x', 130) FROM generate_series(1, 200000) g", testForeign)); err != nil {
		t.Fatalf("Schreiblast außerhalb der Publication: %v", err)
	}
	loadedAt := time.Now()
	afterLoad := currentWALLSN(t, env.pool)
	load := walLSNDiff(t, env.pool, before, afterLoad)
	if load < 16*1024*1024 {
		t.Fatalf("Last trägt %d B WAL, erwartet mindestens 16 MiB", load)
	}
	end, err := pglogrepl.ParseLSN(afterLoad)
	if err != nil {
		t.Fatalf("LSN %q: %v", afterLoad, err)
	}

	limit := load / 10
	deadline := loadedAt.Add(wait)
	for {
		backlog, err := checker.Measure(context.Background())
		if err != nil {
			t.Fatalf("Measure: %v", err)
		}
		confirmed := readConfirmedFlush(t, env.pool, env.slot)
		reached := backlog < limit
		if byPosition {
			reached = confirmed >= uint64(end)
		}
		if reached || time.Now().After(deadline) {
			return foreignWALRun{load: load, backlog: backlog, confirmed: confirmed, end: uint64(end), elapsed: time.Since(loadedAt)}
		}
		time.Sleep(250 * time.Millisecond)
	}
}

// TestStreamIdleConfirmationReleasesForeignWAL trägt die Form X1 des
// Verdikts am realen Stream (`ADR-0120`, `LH-QA-REL-001.a`): WAL ohne Inhalt
// für die Publication hält den Rückstand des Slots nicht — nach wenigen
// Keepalive-Takten liegt er unter einem Zehntel der Last, und der Stand-in
// hat Leerlauf-Positionen bestätigt. Die Instanz trägt den Standardwert von
// `wal_sender_timeout` (60 s) und hat keinen anderen Schreiber
// (`CDC_REPLICATION_TEST_STANDARD_DSN`, `run-replication-tests.sh`): der
// Rückstand der Instanz ist der des Slots.
func TestStreamIdleConfirmationReleasesForeignWAL(t *testing.T) {
	standIn := &fakeCapture{}
	run := runForeignWAL(t, "CDC_REPLICATION_TEST_STANDARD_DSN", "foreign", standIn, 150*time.Second, false)
	t.Logf("Last %d B WAL, Rückstand %d B nach %s, %d Leerlauf-Bestätigungen",
		run.load, run.backlog, run.elapsed.Round(time.Millisecond), standIn.idleConfirmations())
	if run.backlog >= run.load/10 {
		t.Fatalf("Rückstand %d B liegt nach %s nicht unter einem Zehntel der Last (%d B)", run.backlog, run.elapsed, run.load)
	}
	if standIn.idleConfirmations() == 0 {
		t.Fatalf("keine Leerlauf-Bestätigung des Stand-ins")
	}
}

// TestStreamWithoutIdleConfirmationKeepsForeignWAL ist die Nullprobe zu
// `TestStreamIdleConfirmationReleasesForeignWAL`: bestätigt die Application
// keine Leerlauf-Position (`declineIdle`), bleibt der Rückstand über die
// gleiche Wartezeit bei mindestens der halben Last — die Last ist real, und
// nur die Leerlauf-Bestätigung senkt den Rückstand.
func TestStreamWithoutIdleConfirmationKeepsForeignWAL(t *testing.T) {
	standIn := &fakeCapture{declineIdle: true}
	run := runForeignWAL(t, "CDC_REPLICATION_TEST_STANDARD_DSN", "foreigncontrol", standIn, 8*time.Second, false)
	t.Logf("Last %d B WAL, Rückstand %d B nach %s", run.load, run.backlog, run.elapsed.Round(time.Millisecond))
	if run.backlog < run.load/2 {
		t.Fatalf("Rückstand %d B fiel ohne Leerlauf-Bestätigung unter die Hälfte der Last (%d B)", run.backlog, run.load)
	}
}

// TestStreamIdleConfirmationReachesEndOfForeignWAL trägt dieselbe Form X1 auf
// der Instanz der übrigen Stream-Tests, die fremde Schreiber teilt und
// `wal_sender_timeout=2s` trägt: die Bedingung liegt an der Position, nicht am
// Rückstand — `confirmed_flush_lsn` erreicht das WAL-Ende hinter der Last,
// unabhängig von WAL, das andere Test-Pakete gleichzeitig erzeugen.
func TestStreamIdleConfirmationReachesEndOfForeignWAL(t *testing.T) {
	standIn := &fakeCapture{}
	run := runForeignWAL(t, "CDC_REPLICATION_TEST_DSN", "foreignend", standIn, 30*time.Second, true)
	t.Logf("Last %d B WAL, confirmed_flush_lsn %x gegen WAL-Ende der Last %x nach %s, %d Leerlauf-Bestätigungen",
		run.load, run.confirmed, run.end, run.elapsed.Round(time.Millisecond), standIn.idleConfirmations())
	if run.confirmed < run.end {
		t.Fatalf("confirmed_flush_lsn %x erreicht das WAL-Ende der Last (%x) nach %s nicht", run.confirmed, run.end, run.elapsed)
	}
}

// TestStreamIdleConfirmationKeepsOpenTransactionDeliverable ist der
// Sicherheits-Test der Leerlauf-Bestätigung (`ADR-0120` Festlegung 1,
// `LH-QA-REL-001`): eine Quelltransaktion auf einer veröffentlichten Tabelle
// bleibt offen, währenddessen bestätigt der Stream im Leerlauf ein WAL-Ende
// hinter dem ersten Change dieser Transaktion (`confirmed_flush_lsn` liegt
// danach hinter ihm). Endet der Stream und committet die Transaktion danach,
// liefert der Neustart ihren Change — die Quelle sendet jede Transaktion,
// deren Commit hinter `confirmed_flush_lsn` liegt, vollständig.
func TestStreamIdleConfirmationKeepsOpenTransactionDeliverable(t *testing.T) {
	pool, ctx := newPool(t)
	env := newTestEnv(t, "opentx")
	createForeignTable(t, env)

	ctx1, cancel1 := context.WithCancel(ctx)
	defer cancel1()
	commands1 := make(chan *inbound.CaptureCommand, 32)
	standIn1 := &fakeCapture{commands: commands1}
	stream1, err := newStream(ctx1, env, standIn1)
	if err != nil {
		t.Fatalf("NewStream (erster Lauf): %v", err)
	}
	runDone1 := make(chan error, 1)
	go func() { runDone1 <- stream1.Run(ctx1) }()

	// Der Stream läuft, sobald die erste Transaktion der aktivierten Tabelle
	// eingetroffen ist.
	if _, err := pool.Exec(ctx, "INSERT INTO "+testFeed+" (id, name) VALUES (1, 'Start')"); err != nil {
		t.Fatalf("INSERT: %v", err)
	}
	awaitCommand(t, commands1, 15*time.Second)

	// Die offene Quelltransaktion schreibt auf die veröffentlichte Tabelle;
	// `openEnd` liegt hinter ihrem WAL-Eintrag.
	sourceTx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatalf("Quelltransaktion: %v", err)
	}
	t.Cleanup(func() { _ = sourceTx.Rollback(context.Background()) })
	if _, err := sourceTx.Exec(ctx, "INSERT INTO "+testFeed+" (id, name) VALUES (2, 'Offen')"); err != nil {
		t.Fatalf("INSERT (offene Transaktion): %v", err)
	}
	var openEnd string
	if err := pool.QueryRow(ctx, "SELECT pg_current_wal_insert_lsn()::text").Scan(&openEnd); err != nil {
		t.Fatalf("WAL-Ende der offenen Transaktion: %v", err)
	}
	openEndLSN, err := pglogrepl.ParseLSN(openEnd)
	if err != nil {
		t.Fatalf("LSN %q: %v", openEnd, err)
	}

	// WAL ohne Inhalt für die Publication nach der offenen Transaktion —
	// der Stream ist im Leerlauf und bestätigt bis hinter sie.
	if _, err := pool.Exec(ctx,
		fmt.Sprintf("INSERT INTO %s SELECT g, repeat('x', 130) FROM generate_series(1, 20000) g", testForeign)); err != nil {
		t.Fatalf("WAL nach der offenen Transaktion: %v", err)
	}
	deadline := time.Now().Add(30 * time.Second)
	for readConfirmedFlush(t, pool, env.slot) <= uint64(openEndLSN) {
		if time.Now().After(deadline) {
			t.Fatalf("confirmed_flush_lsn %x rückt nicht hinter den Change der offenen Transaktion (%x)",
				readConfirmedFlush(t, pool, env.slot), uint64(openEndLSN))
		}
		time.Sleep(200 * time.Millisecond)
	}
	if standIn1.idleConfirmations() == 0 {
		t.Fatalf("confirmed_flush_lsn liegt hinter der offenen Transaktion ohne Leerlauf-Bestätigung")
	}

	// Der Stream endet, die Transaktion committet, der Neustart liefert sie.
	cancel1()
	select {
	case <-runDone1:
	case <-time.After(15 * time.Second):
		t.Fatalf("erster Lauf endet nicht")
	}
	awaitSlotInactive(t, pool, env.slot, 10*time.Second)
	if err := sourceTx.Commit(ctx); err != nil {
		t.Fatalf("Commit der offenen Transaktion: %v", err)
	}

	commands2 := make(chan *inbound.CaptureCommand, 32)
	ctx2, cancel2 := context.WithCancel(ctx)
	t.Cleanup(cancel2)
	stream2, err := newStream(ctx2, env, &fakeCapture{commands: commands2})
	if err != nil {
		t.Fatalf("NewStream (Neustart): %v", err)
	}
	go func() { _ = stream2.Run(ctx2) }()

	limit := time.Now().Add(20 * time.Second)
	for {
		command := awaitCommand(t, commands2, time.Until(limit))
		changes, err := command.Transaction.Changes()
		if err != nil {
			t.Fatalf("Changes: %v", err)
		}
		if len(changes) == 1 && string(changes[0].NewImage) == `{"id":"2","name":"Offen"}` {
			return
		}
	}
}
