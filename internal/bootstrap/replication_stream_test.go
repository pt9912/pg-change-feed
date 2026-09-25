package bootstrap_test

import (
	"context"
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
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/capture"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Der Verdrahtungs-Test läuft im Composition-Root-Layer (`ADR-0026`):
// er verdrahtet die konkreten Adapter — Replication-Stream, ChangeStore
// und ACK — und trägt die Persist-before-ACK-Ordnung am realen Treiber
// (`LH-QA-REL-001.a`). Die Instanz gehört dem gepinnten Testcontainer
// (`make test-replication`, `ADR-0030`); ohne DSN überspringt der Test.

const (
	wireSource  = "src-1"
	wireTableID = "tbl-1"
	wireSchemaV = "sv-1"
	wireFeed    = "public.feed_wire_test"
)

// TestRealPersistBeforeAck trägt die Persist-before-ACK-Ordnung am
// realen Treiber (`LH-QA-REL-001.a`): der echte Capture Service
// persistiert über den ChangeStore und bestätigt die Position über den
// realen ACK-Adapter an der Stream-Verbindung — confirmed_flush_lsn
// trägt die bestätigte Position und die persistierten Changes sind
// lesbar (`ADR-0027`, `ADR-0007`).
func TestRealPersistBeforeAck(t *testing.T) {
	dsn := os.Getenv("CDC_REPLICATION_TEST_DSN")
	if dsn == "" {
		t.Skip("CDC_REPLICATION_TEST_DSN nicht gesetzt — reale Verdrahtungs-Tests laufen über make test-replication")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("Verbindungsaufbau: %v", err)
	}
	t.Cleanup(pool.Close)

	publication := "pub_pgc_test_wire"
	slot := "slot_pgc_test_wire"
	for _, statement := range []string{
		fmt.Sprintf("CREATE TABLE %s (id int PRIMARY KEY, name text)", wireFeed),
		fmt.Sprintf("CREATE PUBLICATION %s FOR TABLE %s", publication, wireFeed),
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
			_, err := pool.Exec(dropCtx, fmt.Sprintf("SELECT pg_drop_replication_slot('%s')", slot))
			if err == nil {
				break
			}
			time.Sleep(200 * time.Millisecond)
		}
		for _, statement := range []string{
			"DROP PUBLICATION IF EXISTS " + publication,
			"DROP TABLE IF EXISTS " + wireFeed,
		} {
			_, _ = pool.Exec(dropCtx, statement)
		}
	})

	// Den Schema-Stand dieses Laufs trägt die eine Schema-Anwendung der
	// Test-Läufe (`tools/schema/apply-rollout.sh`), die der Lauf-Aufruf vor
	// dem Tier-Lauf ausführt — derselbe d-migrate-Rollout wie im Betrieb und
	// in `make test-store`. Dieser Test trägt darin nur seine Referenz-Zeilen.
	for _, statement := range []string{
		fmt.Sprintf("INSERT INTO cdc.source (source_id, name) VALUES ('%s', 'Quelle')", wireSource),
		fmt.Sprintf("INSERT INTO cdc.source_table (source_table_id, source_id, schema_name, table_name) VALUES ('%s', '%s', 'public', 'feed_wire_test')", wireTableID, wireSource),
		fmt.Sprintf("INSERT INTO cdc.schema_version (schema_version_id, source_table_id, version) VALUES ('%s', '%s', 1)", wireSchemaV, wireTableID),
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

	// Die Verdrahtung läuft in dieser Reihenfolge: der Stream baut die
	// Verbindung, der ACK-Adapter trägt seine Bestätigung über dieselbe
	// Verbindung (`ADR-0007`, Option C), der Service orchestriert.
	stream, err := receive.NewStream(ctx, receive.Config{
		DSN:         dsn,
		Source:      wireSource,
		Publication: publication,
		Slot:        slot,
		Tables: map[string]mapper.TableBinding{
			wireFeed: {TableID: wireTableID, SchemaVersion: wireSchemaV},
		},
	})
	if err != nil {
		t.Fatalf("NewStream: %v", err)
	}
	ack, err := postgresack.New(stream.Conn())
	if err != nil {
		t.Fatalf("ACK-Adapter: %v", err)
	}
	service := capture.NewCaptureService(store, ack)
	if err := stream.BindCapture(service); err != nil {
		t.Fatalf("BindCapture: %v", err)
	}
	if err := stream.BindIdleConfirmation(service); err != nil {
		t.Fatalf("BindIdleConfirmation: %v", err)
	}
	go func() {
		_ = stream.Run(ctx)
	}()

	if _, err := pool.Exec(ctx, "INSERT INTO "+wireFeed+" (id, name) VALUES (1, 'Alpha')"); err != nil {
		t.Fatalf("INSERT: %v", err)
	}
	if _, err := pool.Exec(ctx, "INSERT INTO "+wireFeed+" (id, name) VALUES (2, 'Beta')"); err != nil {
		t.Fatalf("INSERT: %v", err)
	}

	records, err := awaitPersistedChanges(t, store, 2)
	if err != nil {
		t.Fatalf("Persistierte Changes: %v", err)
	}
	if len(records) != 2 {
		t.Fatalf("Persistierte Changes: %d", len(records))
	}
	lastPosition := records[len(records)-1].Position

	// Die bestätigte Position trägt confirmed_flush_lsn; sie kann die
	// letzte Commit-Position nicht unterschreiten — der ACK lief erst
	// nach der Persistenz (`LH-QA-REL-001.a`). Der Slot-Stand trägt den
	// Feedback-Zug asynchron; der Test pollt mit Zeitgrenze.
	var confirmed string
	deadline := time.Now().Add(15 * time.Second)
	for {
		values, err := pool.Query(ctx,
			"SELECT confirmed_flush_lsn::text FROM pg_replication_slots WHERE slot_name = $1", slot)
		if err != nil {
			t.Fatalf("Slot-Stand: %v", err)
		}
		if !values.Next() {
			values.Close()
			t.Fatalf("Slot %q fehlt", slot)
		}
		if err := values.Scan(&confirmed); err != nil {
			values.Close()
			t.Fatalf("Slot-Stand lesen: %v", err)
		}
		values.Close()
		confirmedLSN, err := pglogrepl.ParseLSN(confirmed)
		if err != nil {
			t.Fatalf("confirmed_flush_lsn %q: %v", confirmed, err)
		}
		if uint64(confirmedLSN) >= lastPosition.Offset {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("confirmed_flush_lsn %x unterschreitet die bestätigte Position %x", uint64(confirmedLSN), lastPosition.Offset)
		}
		time.Sleep(200 * time.Millisecond)
	}

	// Der Standby-Status-Aufruf des ACK-Adapters trägt die Position
	// direkt: der Test bestätigt eine Position 64 KiB hinter dem aktuellen
	// WAL-Ende der Instanz — nur der ACK-Aufruf kann confirmed_flush_lsn
	// dorthin schieben, weil jede Bestätigung des Streams (Capture-Ergebnis
	// und Leerlauf-Bestätigung) höchstens das WAL-Ende erreicht. Die
	// synthetische Position trägt kein persistiertes Change-Guthaben —
	// sie belegt den Transport-Zug, nicht die Persistenz-Ordnung.
	var walEnd string
	if err := pool.QueryRow(ctx, "SELECT pg_current_wal_insert_lsn()::text").Scan(&walEnd); err != nil {
		t.Fatalf("WAL-Ende: %v", err)
	}
	walEndLSN, err := pglogrepl.ParseLSN(walEnd)
	if err != nil {
		t.Fatalf("WAL-Ende %q: %v", walEnd, err)
	}
	position, err := model.NewSourcePosition(wireSource, uint64(walEndLSN)+0x10000)
	if err != nil {
		t.Fatalf("Bestätigungs-Position: %v", err)
	}
	if err := ack.Acknowledge(ctx, position); err != nil {
		t.Fatalf("Acknowledge: %v", err)
	}
	ackDeadline := time.Now().Add(15 * time.Second)
	for {
		values, err := pool.Query(ctx,
			"SELECT confirmed_flush_lsn::text FROM pg_replication_slots WHERE slot_name = $1", slot)
		if err != nil {
			t.Fatalf("Slot-Stand: %v", err)
		}
		if !values.Next() {
			values.Close()
			t.Fatalf("Slot %q fehlt", slot)
		}
		var confirmedAfter string
		if err := values.Scan(&confirmedAfter); err != nil {
			values.Close()
			t.Fatalf("Slot-Stand lesen: %v", err)
		}
		values.Close()
		afterLSN, err := pglogrepl.ParseLSN(confirmedAfter)
		if err != nil {
			t.Fatalf("confirmed_flush_lsn %q: %v", confirmedAfter, err)
		}
		if uint64(afterLSN) >= position.Offset {
			return
		}
		if time.Now().After(ackDeadline) {
			t.Fatalf("confirmed_flush_lsn %x trägt die bestätigte Position %x nicht", uint64(afterLSN), position.Offset)
		}
		time.Sleep(200 * time.Millisecond)
	}
}

// awaitPersistedChanges liest den Store, bis die erwarteten Changes
// persistiert sind (Polling mit Test-Zeitgrenze).
func awaitPersistedChanges(t *testing.T, store *postgresstorage.PostgresChangeStoreAdapter, limit int) ([]outbound.ChangeRecord, error) {
	t.Helper()
	deadline := time.Now().Add(15 * time.Second)
	source := model.SourceID(wireSource)
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

// TestRealIdleConfirmationReleasesForeignWAL trägt die Leerlauf-Bestätigung
// am zusammengesetzten System (`ADR-0120`, `ADR-0007`): der echte Capture
// Service bestätigt die vom Stream gemeldete Leerlauf-Position über den
// realen ACK-Adapter an der Stream-Verbindung. Schreiblast außerhalb der
// Publication hält den Slot nicht zurück — `confirmed_flush_lsn` erreicht das
// WAL-Ende hinter der Last (`wal_sender_timeout=2s` des Testcontainers), und
// die Bestätigung schreibt nichts in `cdc.transaction` oder `cdc.change`
// (`ADR-0120` Festlegung 2).
func TestRealIdleConfirmationReleasesForeignWAL(t *testing.T) {
	dsn := os.Getenv("CDC_REPLICATION_TEST_DSN")
	if dsn == "" {
		t.Skip("CDC_REPLICATION_TEST_DSN nicht gesetzt — reale Verdrahtungs-Tests laufen über make test-replication")
	}
	const (
		source      = "src-idle"
		tableID     = "tbl-idle"
		schemaV     = "sv-idle"
		feed        = "public.feed_wire_idle_test"
		foreign     = "public.foreign_wire_idle_test"
		publication = "pub_pgc_test_wire_idle"
		slot        = "slot_pgc_test_wire_idle"
	)
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("Verbindungsaufbau: %v", err)
	}
	t.Cleanup(pool.Close)

	for _, statement := range []string{
		fmt.Sprintf("CREATE TABLE %s (id int PRIMARY KEY, name text)", feed),
		fmt.Sprintf("CREATE TABLE %s (id bigint PRIMARY KEY, payload text)", foreign),
		fmt.Sprintf("CREATE PUBLICATION %s FOR TABLE %s", publication, feed),
	} {
		if _, err := pool.Exec(ctx, statement); err != nil {
			t.Fatalf("Test-Umgebung: %v", err)
		}
	}
	runCtx, cancelRun := context.WithCancel(ctx)
	runDone := make(chan struct{})
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
			"DROP TABLE IF EXISTS " + foreign,
		} {
			_, _ = pool.Exec(dropCtx, statement)
		}
	})
	t.Cleanup(func() {
		cancelRun()
		select {
		case <-runDone:
		case <-time.After(10 * time.Second):
		}
	})

	for _, statement := range []string{
		fmt.Sprintf("INSERT INTO cdc.source (source_id, name) VALUES ('%s', 'Quelle')", source),
		fmt.Sprintf("INSERT INTO cdc.source_table (source_table_id, source_id, schema_name, table_name) VALUES ('%s', '%s', 'public', 'feed_wire_idle_test')", tableID, source),
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

	stream, err := receive.NewStream(runCtx, receive.Config{
		DSN:         dsn,
		Source:      source,
		Publication: publication,
		Slot:        slot,
		Tables: map[string]mapper.TableBinding{
			feed: {TableID: tableID, SchemaVersion: schemaV},
		},
	})
	if err != nil {
		t.Fatalf("NewStream: %v", err)
	}
	ack, err := postgresack.New(stream.Conn())
	if err != nil {
		t.Fatalf("ACK-Adapter: %v", err)
	}
	service := capture.NewCaptureService(store, ack)
	if err := stream.BindCapture(service); err != nil {
		t.Fatalf("BindCapture: %v", err)
	}
	if err := stream.BindIdleConfirmation(service); err != nil {
		t.Fatalf("BindIdleConfirmation: %v", err)
	}
	go func() {
		defer close(runDone)
		_ = stream.Run(runCtx)
	}()

	countRows := func(table string) int {
		var count int
		query := fmt.Sprintf("SELECT count(*) FROM cdc.%s WHERE ", table)
		if table == "transaction" {
			query += "source_id = $1"
		} else {
			query += "transaction_id IN (SELECT transaction_id FROM cdc.transaction WHERE source_id = $1)"
		}
		if err := pool.QueryRow(ctx, query, source).Scan(&count); err != nil {
			t.Fatalf("cdc.%s zählen: %v", table, err)
		}
		return count
	}

	if _, err := pool.Exec(ctx, "INSERT INTO "+feed+" (id, name) VALUES (1, 'Start')"); err != nil {
		t.Fatalf("INSERT: %v", err)
	}
	persistDeadline := time.Now().Add(15 * time.Second)
	for countRows("change") < 1 {
		if time.Now().After(persistDeadline) {
			t.Fatalf("die Change der aktivierten Tabelle wird nicht persistiert")
		}
		time.Sleep(100 * time.Millisecond)
	}

	var before string
	if err := pool.QueryRow(ctx, "SELECT pg_current_wal_lsn()::text").Scan(&before); err != nil {
		t.Fatalf("WAL-Position: %v", err)
	}
	if _, err := pool.Exec(ctx,
		fmt.Sprintf("INSERT INTO %s SELECT g, repeat('x', 130) FROM generate_series(1, 100000) g", foreign)); err != nil {
		t.Fatalf("Schreiblast außerhalb der Publication: %v", err)
	}
	var after string
	var load int64
	if err := pool.QueryRow(ctx,
		"SELECT pg_current_wal_lsn()::text, pg_wal_lsn_diff(pg_current_wal_lsn(), $1::pg_lsn)::bigint", before).Scan(&after, &load); err != nil {
		t.Fatalf("Last: %v", err)
	}
	if load < 8*1024*1024 {
		t.Fatalf("Last trägt %d B WAL, erwartet mindestens 8 MiB", load)
	}

	// Die Bedingung liegt an der Position: `confirmed_flush_lsn` erreicht das
	// WAL-Ende hinter der Last — unabhängig von WAL, das andere Test-Pakete
	// dieser Instanz gleichzeitig erzeugen.
	deadline := time.Now().Add(30 * time.Second)
	for {
		var reached bool
		if err := pool.QueryRow(ctx,
			"SELECT confirmed_flush_lsn >= $1::pg_lsn FROM pg_replication_slots WHERE slot_name = $2", after, slot).Scan(&reached); err != nil {
			t.Fatalf("Slot-Stand: %v", err)
		}
		if reached {
			t.Logf("Last %d B WAL, confirmed_flush_lsn hat das WAL-Ende der Last (%s) erreicht", load, after)
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("confirmed_flush_lsn erreicht das WAL-Ende der Last (%s) nach 30 s nicht", after)
		}
		time.Sleep(250 * time.Millisecond)
	}
	if got := countRows("transaction"); got != 1 {
		t.Fatalf("cdc.transaction trägt %d Zeilen der Quelle, erwartet 1 — die Leerlauf-Bestätigung schreibt nichts", got)
	}
	if got := countRows("change"); got != 1 {
		t.Fatalf("cdc.change trägt %d Zeilen der Quelle, erwartet 1", got)
	}
}
