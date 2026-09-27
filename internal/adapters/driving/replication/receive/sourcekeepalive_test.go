package receive_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pglogrepl"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgproto3"
)

// Der Beleg der Quellseite läuft gegen die Instanz mit dem Standardwert von
// `wal_sender_timeout` (60 s) und ohne fremden Schreiber; `run-replication-tests.sh`
// startet ihn gesondert im Tier-Lauf und verlangt sein `--- PASS`. Ohne DSN
// überspringt der Test.
const (
	sourceKeepaliveDSNVariable = "CDC_SOURCE_KEEPALIVE_TEST_DSN"

	// sourceKeepaliveRows ist die Größe der Quelltransaktion: sie übersteigt
	// die Puffer der Verbindung, damit der Walsender beim Senden blockiert und
	// seinen Keepalive zwischen die Nachrichten der Transaktion stellt.
	sourceKeepaliveRows = 400000

	// sourceKeepaliveSilence ist die Zeit ab START_REPLICATION, in der der
	// Client nichts liest und nichts antwortet: der Walsender sendet seinen
	// Keepalive nach der Hälfte von `wal_sender_timeout` (30 s) und bricht die
	// Verbindung nach dem ganzen Wert (60 s) ab.
	sourceKeepaliveSilence = 38 * time.Second
)

// openRawReplication öffnet eine Replication-Verbindung ohne den
// Stream-Adapter: der Test liest das Protokoll selbst und entscheidet als
// Client, welche Position er bestätigt.
func openRawReplication(ctx context.Context, t *testing.T, dsn string) *pgconn.PgConn {
	t.Helper()
	config, err := pgconn.ParseConfig(dsn)
	if err != nil {
		t.Fatalf("DSN: %v", err)
	}
	config.RuntimeParams["replication"] = "database"
	conn, err := pgconn.ConnectConfig(ctx, config)
	if err != nil {
		t.Fatalf("Replication-Verbindung: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close(context.Background()) })
	return conn
}

// startRawStream beginnt den Stream auf dem Slot ab der gegebenen Position.
func startRawStream(ctx context.Context, t *testing.T, conn *pgconn.PgConn, env *testEnv, from pglogrepl.LSN) {
	t.Helper()
	if err := pglogrepl.StartReplication(ctx, conn, env.slot, from, pglogrepl.StartReplicationOptions{
		Mode: pglogrepl.LogicalReplication,
		PluginArgs: []string{
			"proto_version '1'",
			"publication_names '" + env.publication + "'",
		},
	}); err != nil {
		t.Fatalf("START_REPLICATION: %v", err)
	}
}

// confirmRaw sendet das Standby-Status-Update mit der Position.
func confirmRaw(ctx context.Context, t *testing.T, conn *pgconn.PgConn, position pglogrepl.LSN) {
	t.Helper()
	if err := pglogrepl.SendStandbyStatusUpdate(ctx, conn, pglogrepl.StandbyStatusUpdate{
		WALWritePosition: position,
		WALFlushPosition: position,
		WALApplyPosition: position,
	}); err != nil {
		t.Fatalf("Standby-Status-Update %s: %v", position, err)
	}
}

// keepaliveInsideTransaction trägt, was der Client am ersten Keepalive
// zwischen BEGIN und COMMIT einer Quelltransaktion liest.
type keepaliveInsideTransaction struct {
	serverWALEnd pglogrepl.LSN
	finalLSN     pglogrepl.LSN
	changesSoFar int
}

// readUntilKeepaliveInsideTransaction liest, bis ein Keepalive zwischen BEGIN
// und COMMIT eintrifft; ein COMMIT davor beendet den Test.
func readUntilKeepaliveInsideTransaction(ctx context.Context, t *testing.T, conn *pgconn.PgConn) keepaliveInsideTransaction {
	t.Helper()
	var current keepaliveInsideTransaction
	inTransaction := false
	for {
		raw, err := conn.ReceiveMessage(ctx)
		if err != nil {
			t.Fatalf("Empfang: %v", err)
		}
		copyData, ok := raw.(*pgproto3.CopyData)
		if !ok || len(copyData.Data) == 0 {
			continue
		}
		switch copyData.Data[0] {
		case pglogrepl.XLogDataByteID:
			xlog, err := pglogrepl.ParseXLogData(copyData.Data[1:])
			if err != nil {
				t.Fatalf("XLogData: %v", err)
			}
			message, err := pglogrepl.Parse(xlog.WALData)
			if err != nil {
				t.Fatalf("pgoutput: %v", err)
			}
			switch message := message.(type) {
			case *pglogrepl.BeginMessage:
				inTransaction = true
				current = keepaliveInsideTransaction{finalLSN: message.FinalLSN}
			case *pglogrepl.InsertMessage:
				current.changesSoFar++
			case *pglogrepl.CommitMessage:
				t.Fatalf("COMMIT nach %d Änderungen ohne Keepalive zwischen BEGIN und COMMIT", current.changesSoFar)
			}
		case pglogrepl.PrimaryKeepaliveMessageByteID:
			keepalive, err := pglogrepl.ParsePrimaryKeepaliveMessage(copyData.Data[1:])
			if err != nil {
				t.Fatalf("Keepalive: %v", err)
			}
			if inTransaction {
				current.serverWALEnd = keepalive.ServerWALEnd
				return current
			}
		}
	}
}

// deliveredTransaction trägt die Quelltransaktion, die der Client nach dem
// Neustart des Streams vollständig liest.
type deliveredTransaction struct {
	finalLSN  pglogrepl.LSN
	commitLSN pglogrepl.LSN
	changes   int
}

// readOneTransaction liest bis zum COMMIT der ersten Transaktion mit
// Änderungen und beantwortet die Keepalives mit der bestätigten Position.
func readOneTransaction(ctx context.Context, t *testing.T, conn *pgconn.PgConn, confirmed pglogrepl.LSN) deliveredTransaction {
	t.Helper()
	var current deliveredTransaction
	for {
		raw, err := conn.ReceiveMessage(ctx)
		if err != nil {
			t.Fatalf("Der Neustart des Streams lieferte die Transaktion nicht bis zum COMMIT (%d Änderungen gelesen): %v", current.changes, err)
		}
		copyData, ok := raw.(*pgproto3.CopyData)
		if !ok || len(copyData.Data) == 0 {
			continue
		}
		switch copyData.Data[0] {
		case pglogrepl.XLogDataByteID:
			xlog, err := pglogrepl.ParseXLogData(copyData.Data[1:])
			if err != nil {
				t.Fatalf("XLogData: %v", err)
			}
			message, err := pglogrepl.Parse(xlog.WALData)
			if err != nil {
				t.Fatalf("pgoutput: %v", err)
			}
			switch message := message.(type) {
			case *pglogrepl.BeginMessage:
				current = deliveredTransaction{finalLSN: message.FinalLSN}
			case *pglogrepl.InsertMessage:
				current.changes++
			case *pglogrepl.CommitMessage:
				current.commitLSN = message.CommitLSN
				if current.changes > 0 {
					return current
				}
			}
		case pglogrepl.PrimaryKeepaliveMessageByteID:
			keepalive, err := pglogrepl.ParsePrimaryKeepaliveMessage(copyData.Data[1:])
			if err != nil {
				t.Fatalf("Keepalive: %v", err)
			}
			if keepalive.ReplyRequested {
				confirmRaw(ctx, t, conn, confirmed)
			}
		}
	}
}

// TestSourceKeepaliveInsideTransactionDeliversWholeTransaction belegt die
// Quellseite der Leerlauf-Bestätigung am realen Protokoll (`ADR-0121`): ein
// roher Client liest eine große Quelltransaktion erst nach der Hälfte von
// `wal_sender_timeout`, der Keepalive des Walsenders steht dann zwischen BEGIN
// und COMMIT; der Test bestätigt dessen `ServerWALEnd`, beendet die
// Verbindung mitten in der Transaktion und startet den Stream ab
// `confirmed_flush_lsn` neu. Gelesen wird, dass die Position des Keepalives die
// Commit-LSN der Transaktion ist (`BEGIN`-Kopf und `COMMIT` nennen dieselbe
// LSN), dass der Slot sie annimmt und dass der Neustart alle
// `sourceKeepaliveRows` Änderungen der Transaktion liefert. Der Adapter
// bestätigt inmitten einer Transaktion nicht (`TestRunNoConfirmationInsideOpenTransaction`);
// der Test liest deshalb das Protokoll selbst. Bestätigt der Client eine
// Position hinter der Commit-LSN (Mutation `+ 1 GiB`), liefert der Neustart die
// Transaktion nicht, und der Test endet an der Lieferung.
func TestSourceKeepaliveInsideTransactionDeliversWholeTransaction(t *testing.T) {
	env := newTestEnvOn(t, "srckeepalive", sourceKeepaliveDSNVariable)
	dsn := env.pool.Config().ConnString()

	if _, err := env.pool.Exec(env.ctx, "SELECT pg_create_logical_replication_slot($1, 'pgoutput')", env.slot); err != nil {
		t.Fatalf("Slot: %v", err)
	}
	var serverVersion string
	if err := env.pool.QueryRow(env.ctx, "SHOW server_version").Scan(&serverVersion); err != nil {
		t.Fatalf("server_version: %v", err)
	}

	// Erster Lauf: der Client liest nichts, während die Transaktion committet.
	readCtx, cancelRead := context.WithTimeout(env.ctx, 2*time.Minute)
	defer cancelRead()
	first := openRawReplication(readCtx, t, dsn)
	startedAt := time.Now()
	startRawStream(readCtx, t, first, env, 0)
	if _, err := env.pool.Exec(env.ctx, fmt.Sprintf(
		"INSERT INTO %s SELECT g, repeat('x', 20) FROM generate_series(1, %d) g", testFeed, sourceKeepaliveRows)); err != nil {
		t.Fatalf("Quelltransaktion: %v", err)
	}
	time.Sleep(time.Until(startedAt.Add(sourceKeepaliveSilence)))

	inside := readUntilKeepaliveInsideTransaction(readCtx, t, first)
	if inside.changesSoFar == 0 || inside.changesSoFar >= sourceKeepaliveRows {
		t.Fatalf("Keepalive nach %d von %d Änderungen liegt nicht inmitten der Transaktion", inside.changesSoFar, sourceKeepaliveRows)
	}
	if inside.serverWALEnd != inside.finalLSN {
		t.Fatalf("ServerWALEnd %s des Keepalives ist nicht die Commit-LSN %s der Transaktion", inside.serverWALEnd, inside.finalLSN)
	}

	// Der Client bestätigt die Position des Keepalives und beendet die
	// Verbindung, bevor der Rest der Transaktion gelesen ist.
	confirmedPosition := inside.serverWALEnd
	confirmRaw(readCtx, t, first, confirmedPosition)
	deadline := time.Now().Add(15 * time.Second)
	for readConfirmedFlush(t, env.pool, env.slot) != uint64(confirmedPosition) {
		if time.Now().After(deadline) {
			t.Fatalf("confirmed_flush_lsn %x erreicht die bestätigte Position %s nicht", readConfirmedFlush(t, env.pool, env.slot), confirmedPosition)
		}
		time.Sleep(100 * time.Millisecond)
	}
	if err := first.Close(context.Background()); err != nil {
		t.Fatalf("Verbindung schließen: %v", err)
	}
	awaitSlotInactive(t, env.pool, env.slot, 15*time.Second)

	// Zweiter Lauf: der Stream setzt bei confirmed_flush_lsn an.
	restartFrom := pglogrepl.LSN(readConfirmedFlush(t, env.pool, env.slot))
	deliverCtx, cancelDeliver := context.WithTimeout(env.ctx, 40*time.Second)
	defer cancelDeliver()
	second := openRawReplication(deliverCtx, t, dsn)
	startRawStream(deliverCtx, t, second, env, restartFrom)
	delivered := readOneTransaction(deliverCtx, t, second, confirmedPosition)

	t.Logf("PostgreSQL %s: Keepalive inmitten der Transaktion nach %d von %d Änderungen, ServerWALEnd %s gleich Commit-LSN %s, bestätigt %s; Neustart ab %s lieferte %d Änderungen mit Commit-LSN %s",
		serverVersion, inside.changesSoFar, sourceKeepaliveRows, inside.serverWALEnd, inside.finalLSN, confirmedPosition,
		restartFrom, delivered.changes, delivered.commitLSN)
	if delivered.changes != sourceKeepaliveRows {
		t.Fatalf("Neustart lieferte %d von %d Änderungen der Transaktion", delivered.changes, sourceKeepaliveRows)
	}
	if delivered.commitLSN != inside.finalLSN || delivered.finalLSN != inside.finalLSN {
		t.Fatalf("Neustart lieferte die Transaktion mit BEGIN-Kopf %s und COMMIT %s, erwartet die Commit-LSN %s", delivered.finalLSN, delivered.commitLSN, inside.finalLSN)
	}
}
