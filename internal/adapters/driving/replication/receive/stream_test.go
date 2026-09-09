package receive_test

import (
	"context"
	stderrors "errors"
	"fmt"
	"os"
	"testing"
	"time"

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
// der Verdrahtungs-Test in der Composition-Root.
type fakeCapture struct {
	commands chan *inbound.CaptureCommand
}

// Capture nimmt ein Command auf und meldet die Commit-Position als
// bestätigt.
func (f *fakeCapture) Capture(_ context.Context, command inbound.CaptureCommand) (inbound.CaptureResult, error) {
	f.commands <- &command
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
// Slot wird mit dem Stream-Adapter angelegt. Der Cleanup stoppt den
// Laufkontext vor dem Slot-Rückbau — ein belegter Slot droppt nicht.
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
		fmt.Sprintf("CREATE PUBLICATION %s FOR ALL TABLES", env.publication),
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
// Test liest über das Command-Feld.
func startStream(t *testing.T, env *testEnv, capturePort inbound.CaptureInboundPort) *receive.Stream {
	t.Helper()
	stream, err := receive.NewStream(env.ctx, receive.Config{
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
		_ = stream.Run(env.ctx)
	}()
	return stream
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
	startStream(t, env, &fakeCapture{commands: commands})

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
	startStream(t, env, &fakeCapture{commands: commands})

	sourceTx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatalf("Quelltransaktion: %v", err)
	}
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
