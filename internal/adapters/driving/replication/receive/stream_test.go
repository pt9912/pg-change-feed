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

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresack"
	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage"
	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/mapper"
	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/receive"
	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/capture"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Die Stream-Tests laufen gegen eine reale PostgreSQL-Instanz mit
// Publication und Logical Replication Slot im Testcontainer (`ADR-0030`,
// `LH-FA-CFG-001.a`), gepinnt über `make test-replication`; ohne DSN
// überspringen sie — die Dekodier- und Mapper-Seite tragen die
// Unit-Tests gegen die `pgoutput`-Binärcodes. Die Instanz gehört dem
// Container: Tabellen, Publication und Slot werden je Test frisch
// aufgesetzt, Daten bleiben im Container und landen nicht im
// Arbeitsbaum.

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
// TestRealPersistBeforeAck mit dem echten Service.
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

	insert := awaitCommand(t, commands, 15*time.Second)
	insertChanges, err := insert.Transaction.Changes()
	if err != nil {
		t.Fatalf("Insert-Changes: %v", err)
	}
	if len(insertChanges) != 1 || insertChanges[0].Operation != model.OperationInsert {
		t.Fatalf("Insert-Change: %+v", insertChanges)
	}
	if string(insertChanges[0].NewImage) != `{"id":"1","name":"Alpha"}` {
		t.Fatalf("Insert-Image: %s", insertChanges[0].NewImage)
	}

	other := awaitCommand(t, commands, 15*time.Second)
	otherChanges, err := other.Transaction.Changes()
	if err != nil {
		t.Fatalf("Other-Changes: %v", err)
	}
	if len(otherChanges) != 0 {
		t.Fatalf("Nicht aktivierte Tabelle trägt %d Changes", len(otherChanges))
	}

	update := awaitCommand(t, commands, 15*time.Second)
	updateChanges, err := update.Transaction.Changes()
	if err != nil {
		t.Fatalf("Update-Changes: %v", err)
	}
	if len(updateChanges) != 1 || updateChanges[0].Operation != model.OperationUpdate {
		t.Fatalf("Update-Change: %+v", updateChanges)
	}

	deleteCommand := awaitCommand(t, commands, 15*time.Second)
	deleteChanges, err := deleteCommand.Transaction.Changes()
	if err != nil {
		t.Fatalf("Delete-Changes: %v", err)
	}
	if len(deleteChanges) != 1 || deleteChanges[0].Operation != model.OperationDelete {
		t.Fatalf("Delete-Change: %+v", deleteChanges)
	}
	if string(deleteChanges[0].OldImage) != `{"id":"1"}` {
		t.Fatalf("Delete-Alt-Image: %s", deleteChanges[0].OldImage)
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

// TestRealPersistBeforeAck trägt die Persist-before-ACK-Ordnung am
// realen Treiber (`LH-QA-REL-001.a`): der echte Capture Service
// persistiert über den ChangeStore und bestätigt die Position über den
// realen ACK-Adapter an der Stream-Verbindung — confirmed_flush_lsn
// trägt die bestätigte Position und die persistierten Changes sind
// lesbar (`ADR-0027`, `ADR-0007`).
func TestRealPersistBeforeAck(t *testing.T) {
	pool, ctx := newPool(t)
	env := newTestEnv(t, "ack")
	if _, err := pool.Exec(ctx, "DROP SCHEMA IF EXISTS cdc CASCADE"); err != nil {
		t.Fatalf("Schema-Rückbau: %v", err)
	}
	if err := postgresstorage.ApplySchema(ctx, pool); err != nil {
		t.Fatalf("ApplySchema: %v", err)
	}
	for _, statement := range []string{
		fmt.Sprintf("INSERT INTO cdc.source (source_id, name) VALUES ('%s', 'Quelle')", testSource),
		fmt.Sprintf("INSERT INTO cdc.source_table (source_table_id, source_id, schema_name, table_name) VALUES ('%s', '%s', 'public', 'feed_stream_test')", testTableID, testSource),
		fmt.Sprintf("INSERT INTO cdc.schema_version (schema_version_id, source_table_id, version) VALUES ('%s', '%s', 1)", testSchemaV, testTableID),
	} {
		if _, err := pool.Exec(ctx, statement); err != nil {
			t.Fatalf("Referenz-Zeilen: %v", err)
		}
	}

	store, err := postgresstorage.New(ctx, env.pool.Config().ConnString())
	if err != nil {
		t.Fatalf("Store: %v", err)
	}
	t.Cleanup(store.Close)

	// Die Verdrahtung läuft in dieser Reihenfolge: der Stream baut die
	// Verbindung, der ACK-Adapter trägt seine Bestätigung über dieselbe
	// Verbindung (`ADR-0007`, Option C), der Service orchestriert.
	stream, err := receive.NewStream(env.ctx, receive.Config{
		DSN:         env.pool.Config().ConnString(),
		Source:      testSource,
		Publication: env.publication,
		Slot:        env.slot,
		Tables: map[string]mapper.TableBinding{
			testFeed: {TableID: testTableID, SchemaVersion: testSchemaV},
		},
	})
	if err != nil {
		t.Fatalf("NewStream: %v", err)
	}
	ack, err := postgresack.New(stream.Conn())
	if err != nil {
		t.Fatalf("ACK-Adapter: %v", err)
	}
	if err := stream.BindCapture(capture.NewCaptureService(store, ack)); err != nil {
		t.Fatalf("BindCapture: %v", err)
	}
	runDone := make(chan error, 1)
	go func() { runDone <- stream.Run(env.ctx) }()

	if _, err := pool.Exec(ctx, "INSERT INTO "+testFeed+" (id, name) VALUES (1, 'Alpha')"); err != nil {
		t.Fatalf("INSERT: %v", err)
	}
	if _, err := pool.Exec(ctx, "INSERT INTO "+testFeed+" (id, name) VALUES (2, 'Beta')"); err != nil {
		t.Fatalf("INSERT: %v", err)
	}

	records, err := awaitPersistedChanges(t, store, 2)
	if err != nil {
		// Der Ausgang des Laufs ist der nächstgelegene Beleg — ein
		// Stream-Fehler endet ohne Persistenz.
		select {
		case runErr := <-runDone:
			t.Fatalf("Persistierte Changes: %v — Stream-Ausgang: %v", err, runErr)
		default:
			t.Fatalf("Persistierte Changes: %v", err)
		}
	}
	if len(records) != 2 {
		t.Fatalf("Persistierte Changes: %d", len(records))
	}
	lastPosition := records[len(records)-1].Position

	// Die bestätigte Position trägt confirmed_flush_lsn; sie kann die
	// letzte Commit-Position nicht unterschreiten — der ACK lief erst
	// nach der Persistenz (`LH-QA-REL-001.a`).
	values, err := pool.Query(ctx,
		"SELECT confirmed_flush_lsn::text FROM pg_replication_slots WHERE slot_name = $1", env.slot)
	if err != nil {
		t.Fatalf("Slot-Stand: %v", err)
	}
	defer values.Close()
	if !values.Next() {
		t.Fatalf("Slot %q fehlt", env.slot)
	}
	var confirmed string
	if err := values.Scan(&confirmed); err != nil {
		t.Fatalf("Slot-Stand lesen: %v", err)
	}
	confirmedLSN, err := pglogrepl.ParseLSN(confirmed)
	if err != nil {
		t.Fatalf("confirmed_flush_lsn %q: %v", confirmed, err)
	}
	if uint64(confirmedLSN) < lastPosition.Offset {
		t.Fatalf("confirmed_flush_lsn %x unterschreitet die bestätigte Position %x", uint64(confirmedLSN), lastPosition.Offset)
	}
}

// awaitPersistedChanges liest den Store, bis die erwarteten Changes
// persistiert sind (Polling mit Test-Zeitgrenze).
func awaitPersistedChanges(t *testing.T, store *postgresstorage.PostgresChangeStoreAdapter, limit int) ([]outbound.ChangeRecord, error) {
	t.Helper()
	deadline := time.Now().Add(15 * time.Second)
	source := model.SourceID(testSource)
	for time.Now().Before(deadline) {
		query := outbound.ChangeQuery{Source: source}
		records, err := store.ReadChanges(context.Background(), query)
		if err != nil {
			return nil, err
		}
		if len(records) >= limit {
			return records, nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return nil, fmt.Errorf("persistierte Changes innerhalb der Zeitspanne fehlgeschlagen")
}