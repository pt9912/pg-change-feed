package receive_test

import (
	"context"
	stderrors "errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pglogrepl"
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
// bestätigt der Stand-in nur die erste Transaktion — die
// Keepalive-Position-Regel (F-3) trägt der Test darüber: die bestätigte
// Position bleibt hinter dem Empfangsstand zurück.
type fakeCapture struct {
	commands     chan *inbound.CaptureCommand
	ackFirstOnly bool
	acked        bool
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
	dsn := os.Getenv("CDC_REPLICATION_TEST_DSN")
	if dsn == "" {
		t.Skip("CDC_REPLICATION_TEST_DSN nicht gesetzt — reale Replication-Tests laufen über make test-replication")
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
	pool, ctx := newPool(t)
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
func startStream(t *testing.T, env *testEnv, capturePort inbound.CaptureInboundPort) (*receive.Stream, <-chan error) {
	t.Helper()
	runDone := make(chan error, 1)
	runCtx, cancel := context.WithCancel(env.ctx)
	t.Cleanup(cancel)
	stream, err := receive.NewStream(runCtx, receive.Config{
		DSN:         env.pool.Config().ConnString(),
		Source:      testSource,
		Publication: env.publication,
		Slot:        env.slot,
		Tables: map[string]mapper.TableBinding{
			testFeed: {TableID: testTableID, SchemaVersion: testSchemaV},
		},
		Capture: capturePort,
	})
	if err != nil {
		t.Fatalf("NewStream: %v", err)
	}
	go func() {
		runDone <- stream.Run(runCtx)
	}()
	return stream, runDone
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
	stream, err := receive.NewStream(env.ctx, receive.Config{
		DSN:         env.pool.Config().ConnString(),
		Source:      testSource,
		Publication: env.publication,
		Slot:        env.slot,
		Tables: map[string]mapper.TableBinding{
			testFeed: {TableID: testTableID, SchemaVersion: testSchemaV},
		},
		Capture: &fakeCapture{commands: commands},
	})
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
// Keepalive-Position-Regel am realen Pfad (`LH-QA-REL-001.a`): die
// Keepalive-Antwort meldet ausschließlich die vom Capture-Ergebnis
// bestätigte Position — confirmed_flush_lsn rückt auf sie und läuft nie
// über sie hinaus, obwohl der Stream weiteren WAL-Stand empfangen hat.
// Der Testcontainer trägt wal_sender_timeout=2s: der Walsender verlangt
// die Antwort nach der Hälfte der Zeit, die Keepalive-Antwort trägt
// lastAcked.
func TestStreamKeepaliveReportsAcknowledgedPosition(t *testing.T) {
	pool, ctx := newPool(t)
	env := newTestEnv(t, "keepalive")
	commands := make(chan *inbound.CaptureCommand, 32)
	startStream(t, env, &fakeCapture{commands: commands, ackFirstOnly: true})

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
	stream1, err := receive.NewStream(ctx1, receive.Config{
		DSN:         env.pool.Config().ConnString(),
		Source:      testSource,
		Publication: env.publication,
		Slot:        env.slot,
		Tables: map[string]mapper.TableBinding{
			testFeed: {TableID: testTableID, SchemaVersion: testSchemaV},
		},
		Capture: &fakeCapture{commands: commands1},
	})
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

	if _, err := pool.Exec(ctx, "INSERT INTO "+testFeed+" (id, name) VALUES (2, 'Zweite')"); err != nil {
		t.Fatalf("INSERT nach Stream-Ende: %v", err)
	}

	commands2 := make(chan *inbound.CaptureCommand, 32)
	ctx2, cancel2 := context.WithCancel(ctx)
	t.Cleanup(cancel2)
	stream2, err := receive.NewStream(ctx2, receive.Config{
		DSN:         env.pool.Config().ConnString(),
		Source:      testSource,
		Publication: env.publication,
		Slot:        env.slot,
		Tables: map[string]mapper.TableBinding{
			testFeed: {TableID: testTableID, SchemaVersion: testSchemaV},
		},
		Capture: &fakeCapture{commands: commands2},
	})
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
