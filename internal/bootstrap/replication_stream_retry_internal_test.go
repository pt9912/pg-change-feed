package bootstrap

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresack"
	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage"
	"github.com/pt9912/pg-change-feed/internal/adapters/driven/systemclock"
	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/mapper"
	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/receive"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/capture"
)

// TestRunStreamWithRetrySlotStillActive belegt den Fall „Slot noch aktiv"
// der Fitness Function (`ADR-0135`): der erste Lauf hält den Slot, der
// Retry-Zyklus scheitert am SQLSTATE 55006 (Transport-/Verbindungsstörung,
// wiederholbar), der Wartezug gibt den Slot frei und der zweite Versuch
// liefert die danach committete Change — Fortsetzung an
// `confirmed_flush_lsn`, Persist-before-ACK unberührt.
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
		Log: diagLogger{t},
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
	go func() {
		defer close(holderDone)
		_ = holderStream.Run(holderCtx)
	}()
	t.Logf("diagnose: Halter-Goroutine gestartet")
	if _, err := pool.Exec(ctx, "INSERT INTO "+feed+" (id, name) VALUES (1, 'Halter')"); err != nil {
		t.Fatalf("INSERT (Halter): %v", err)
	}
	t.Logf("diagnose: Halter-INSERT committet")
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

	// Der Retry: der Zyklus scheitert am ersten Versuch mit SQLSTATE 55006
	// (der Slot ist aktiv), der Wartezug gibt den Slot frei, und der zweite
	// Versuch liefert die danach committete Change.
	cycle := func(attemptCtx context.Context) error {
		t.Logf("diagnose: Zyklus-Versuch startet (attemptCtx: %v)", attemptCtx.Err())
		cycleStream, err := receive.NewStream(attemptCtx, receive.Config{
			DSN:         dsn,
			Source:      source,
			Publication: publication,
			Slot:        slot,
			Tables: map[string]mapper.TableBinding{
				feed: {TableID: tableID, SchemaVersion: schemaV},
			},
		})
		if err != nil {
			return err
		}
		cycleAck, err := postgresack.New(cycleStream.Conn())
		if err != nil {
			return err
		}
		cycleService := capture.NewCaptureService(store, cycleAck)
		if err := cycleStream.BindCapture(cycleService); err != nil {
			return err
		}
		t.Logf("diagnose: Zyklus gebunden, Run startet")
		return cycleStream.Run(attemptCtx)
	}

	if _, err := pool.Exec(ctx, "INSERT INTO "+feed+" (id, name) VALUES (2, 'Retry')"); err != nil {
		t.Fatalf("INSERT (Retry): %v", err)
	}
	retryCtx, stopRetry := context.WithCancel(ctx)
	retryDone := make(chan error, 1)
	go func() {
		retryDone <- runStreamWithRetry(retryCtx, outbound.NoopLog, systemclock.New(), cycle, func(time.Duration) {
			// Der Wartezug gibt den Slot frei: der Halter endet regulär,
			// der Slot steht danach für den zweiten Versuch bereit.
			stopHolder()
			<-holderDone
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
	stopRetry()
	select {
	case err := <-retryDone:
		if err != nil {
			t.Fatalf("Retry-Zyklus: %v", err)
		}
	case <-time.After(10 * time.Second):
		t.Fatalf("Retry-Zyklus endet nach Kontext-Ende nicht")
	}
}

// diagLogger trägt die Log-Ausgabe des untersuchten Streams in den Test-Log.
type diagLogger struct{ t *testing.T }

func (l diagLogger) Debug(_ context.Context, msg string, args ...any) {
	l.t.Logf("DEBUG %s %v", msg, args)
}

func (l diagLogger) Info(_ context.Context, msg string, args ...any) {
	l.t.Logf("INFO %s %v", msg, args)
}

func (l diagLogger) Warn(_ context.Context, msg string, args ...any) {
	l.t.Logf("WARN %s %v", msg, args)
}

func (l diagLogger) Error(_ context.Context, msg string, args ...any) {
	l.t.Logf("ERROR %s %v", msg, args)
}
