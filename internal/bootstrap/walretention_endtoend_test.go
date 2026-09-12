package bootstrap_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage"
	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/mapper"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/bootstrap"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// TestWALRetentionThresholdEndToEnd trägt den in `welle-7.md` §3 geforderten
// Ende-zu-Ende-Beleg (`slice-026`, `ADR-0049`): ein real wachsender
// WAL-Rückstand durchläuft beide Seiten der Schwelle in einem Lauf von
// `bootstrap.Run`. Das Wachstum kommt über Transaktionen auf einer
// *nicht* publizierten Tabelle zustande — dieselbe Ursache, die
// `ADR-0049`s Kontext benennt („WAL-Wachstum inaktiver Slots"): `pgoutput`
// meldet für eine Transaktion ohne Änderung an einer publizierten Tabelle
// weder BEGIN noch COMMIT, der Slot bestätigt sie also nie, während die
// physische WAL trotzdem wächst — `WALRetentionChecker.Measure` (die
// Differenz aus `IDENTIFY_SYSTEM` und `confirmed_flush_lsn`) sieht genau
// diesen Rückstand. Die Schwellen selbst sind ein Test-Override
// (`Config.WALRetentionWarnBytes`/`WALRetentionErrorBytes`) — deutlich
// kleiner als `SPEC-013`s 100 MiB/1 GiB, damit beide Seiten in
// vertretbarer Testzeit real erreicht werden; die Produktions-Startwerte
// selbst bleiben unverändert.
func TestWALRetentionThresholdEndToEnd(t *testing.T) {
	dsn := os.Getenv("CDC_REPLICATION_TEST_DSN")
	if dsn == "" {
		t.Skip("CDC_REPLICATION_TEST_DSN nicht gesetzt — reale Verdrahtungs-Tests laufen über make test-replication")
	}
	ctx, cancelRun := context.WithCancel(context.Background())
	defer cancelRun()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("Verbindungsaufbau: %v", err)
	}
	t.Cleanup(pool.Close)

	const (
		source      = model.SourceID("wal-e2e-src")
		tableID     = "wal-e2e-tbl"
		schemaVerID = "wal-e2e-sv1"
		published   = "public.wal_e2e_published"
		noise       = "public.wal_e2e_noise"
		publication = "pub_wal_e2e"
		slot        = "slot_wal_e2e"
		// warnBytes/errorBytes sind ein Test-Fixture (Dateikommentar oben)
		// — nicht die SPEC-013-Startwerte (100 MiB/1 GiB, `ADR-0049`(b)).
		warnBytes  int64 = 32 * 1024
		errorBytes int64 = 512 * 1024
	)

	for _, stmt := range []string{
		"DROP TABLE IF EXISTS " + noise,
		"DROP TABLE IF EXISTS " + published,
		fmt.Sprintf("CREATE TABLE %s (id int PRIMARY KEY)", published),
		fmt.Sprintf("CREATE TABLE %s (id int PRIMARY KEY, payload text)", noise),
	} {
		if _, err := pool.Exec(ctx, stmt); err != nil {
			t.Fatalf("Test-Umgebung (%s): %v", stmt, err)
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
		for _, stmt := range []string{
			"DROP PUBLICATION IF EXISTS " + publication,
			"DROP TABLE IF EXISTS " + noise,
			"DROP TABLE IF EXISTS " + published,
		} {
			_, _ = pool.Exec(dropCtx, stmt)
		}
	})

	if _, err := pool.Exec(ctx, "DROP SCHEMA IF EXISTS cdc CASCADE"); err != nil {
		t.Fatalf("Schema-Rückbau: %v", err)
	}
	if err := postgresstorage.ApplySchema(ctx, pool); err != nil {
		t.Fatalf("ApplySchema: %v", err)
	}
	// `cdc.table_schema` trägt `ApplySchema` (schema.sql) nicht — ihr Port
	// (SchemaStorePort) liegt außerhalb dieses Store-Adapters, dieselbe
	// Abgrenzung wie bei den Consumer-State-Tabellen
	// (`schemastore_test.go`, `ADR-0015` Folgepflicht). Dieser Lauf
	// verdrahtet `bootstrap.Run` real und durchläuft damit
	// `mapper.Assembler.Consume`s Relation-Behandlung — anders als
	// `cdc.consumer`/`cdc.consumer_position` liegt dieser Port im
	// laufenden Erfassungspfad, die Tabelle muss deshalb hier bestehen.
	if _, err := pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS cdc.table_schema (
		schema_version_id text NOT NULL REFERENCES cdc.schema_version (schema_version_id),
		ordinal_position   bigint NOT NULL CHECK (ordinal_position >= 1),
		column_name        text NOT NULL,
		column_oid         bigint NOT NULL,
		PRIMARY KEY (schema_version_id, ordinal_position),
		UNIQUE (schema_version_id, column_name)
	)`); err != nil {
		t.Fatalf("cdc.table_schema anlegen: %v", err)
	}
	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.source (source_id, name) VALUES ($1, 'WAL-E2E-Quelle')", string(source),
	); err != nil {
		t.Fatalf("Quelle-Zeile: %v", err)
	}

	cfg := bootstrap.Config{
		CaptureDSN:  dsn,
		AdminDSN:    dsn,
		ReaderDSN:   dsn,
		Source:      source,
		Publication: publication,
		Slot:        slot,
		Tables: map[string]mapper.TableBinding{
			published: {TableID: tableID, SchemaVersion: schemaVerID},
		},
		WALRetentionWarnBytes:  warnBytes,
		WALRetentionErrorBytes: errorBytes,
	}

	runDone := make(chan error, 1)
	go func() {
		runDone <- bootstrap.Run(ctx, cfg)
	}()

	// Vorbedingung: die Verdrahtung ist tatsächlich hochgefahren, bevor der
	// Test Wachstum erzeugt — `Run` legt Tabellenbindung, Publication (vor
	// dem Slot, `Run`s Reihenfolge) und zuletzt den Slot an; der Slot in
	// `pg_replication_slots` ist deshalb das späteste, billig abfragbare
	// Bereitschafts-Signal. Ohne dieses Warten liefe die erste Change unten
	// gegen ein Zeitfenster, in dem die Publication ihre Mitgliedschaft
	// noch nicht committet hat — `pgoutput` ließe sie dann lautlos aus
	// (historische Snapshot-Sicht der Katalog-Mitgliedschaft je Transaktion),
	// ohne jeden Fehler.
	awaitSlotExists(t, pool, slot, 30*time.Second)

	// Die erste persistierte Change auf der publizierten Tabelle belegt
	// Aktivierung, Slot und Stream in einem Zug.
	if _, err := pool.Exec(ctx, fmt.Sprintf("INSERT INTO %s (id) VALUES (0)", published)); err != nil {
		t.Fatalf("erste Change auf der publizierten Tabelle: %v", err)
	}
	awaitPublishedChangeCount(t, pool, source, 1, 30*time.Second)

	// Phase 1 (welle-7 §3, Seite 1): WAL-Rückstand real über die
	// Warnschwelle wachsen lassen, aber deutlich unter der Fehlerschwelle
	// halten (Größenordnung: ~150-200 KiB gegen 32 KiB Warn-/512 KiB
	// Fehlerschwelle) — die Transaktion trifft ausschließlich die nicht
	// publizierte Tabelle, `pgoutput` meldet dafür kein BEGIN/COMMIT.
	if _, err := pool.Exec(ctx,
		fmt.Sprintf("INSERT INTO %s (id, payload) SELECT g, repeat('a', 800) FROM generate_series(1, 200) g", noise),
	); err != nil {
		t.Fatalf("Phase 1 (Warnschwelle) Noise-Wachstum: %v", err)
	}

	// Belegt „kontrollierte Fortsetzung" (SPEC-008): Der Capture-Betrieb
	// läuft über der Warnschwelle nachweislich weiter — eine weitere
	// Change auf der publizierten Tabelle wird noch verarbeitet, und
	// `Run` ist währenddessen nicht zurückgekehrt.
	if _, err := pool.Exec(ctx, fmt.Sprintf("INSERT INTO %s (id) VALUES (1)", published)); err != nil {
		t.Fatalf("zweite Change auf der publizierten Tabelle: %v", err)
	}
	awaitPublishedChangeCount(t, pool, source, 2, 30*time.Second)
	// Über mehrere Sekunden hinweg nicht zurückkehren: `heartbeatInterval`
	// (5s, `internal/bootstrap`) ist der Takt des periodischen
	// Schwellen-Vergleichs — dieses Warten deckt mindestens einen weiteren
	// Tick im Warn-Bereich ab, ohne dass er den Lauf beendet.
	stillRunningDeadline := time.Now().Add(8 * time.Second)
	for time.Now().Before(stillRunningDeadline) {
		select {
		case err := <-runDone:
			t.Fatalf("Run() ist zwischen Warn- und Fehlerschwelle beendet (%v) — SPEC-008 verlangt kontrollierte Fortsetzung, keinen Abbruch", err)
		case <-time.After(200 * time.Millisecond):
		}
	}

	// Phase 2 (welle-7 §3, Seite 2): WAL-Rückstand real über die
	// Fehlerschwelle wachsen lassen (Größenordnung: ~2 MiB gegen 512 KiB
	// Fehlerschwelle, oben drauf) — derselbe Mechanismus, deutlich größer.
	if _, err := pool.Exec(ctx,
		fmt.Sprintf("INSERT INTO %s (id, payload) SELECT g, repeat('a', 800) FROM generate_series(1000, 3500) g", noise),
	); err != nil {
		t.Fatalf("Phase 2 (Fehlerschwelle) Noise-Wachstum: %v", err)
	}

	select {
	case err := <-runDone:
		if err == nil {
			t.Fatal("Run() ist regulär (nil) beendet — wollen einen replication-klassifizierten Fehler oberhalb der Fehlerschwelle")
		}
		if !errors.Is(err, outbound.ErrReplication) {
			t.Fatalf("Run() endete mit %v, wollen einen Fehler der Klasse replication (outbound.ErrReplication) — derselbe Pfad wie classifyRunError/os.Exit(1) in main.go", err)
		}
		if errors.Is(err, mapper.ErrChangeWithoutBegin) || errors.Is(err, mapper.ErrCommitWithoutBegin) || errors.Is(err, mapper.ErrBeginWithoutCommit) {
			t.Fatalf("Run() endete mit einer Stream-Ordnungs-Verletzung (%v) statt dem WAL-Schwellen-Fehler — Sentinel-Trennung verletzt (ADR-0049(a))", err)
		}
	case <-time.After(60 * time.Second):
		t.Fatal("Run() hat nach Überschreiten der Fehlerschwelle nicht innerhalb 60s beendet — Hängenbleiben statt kontrollierten Abbruchs")
	}
}

// awaitSlotExists pollt `pg_replication_slots` auf den benannten Slot —
// das Bereitschafts-Signal für `bootstrap.Run`s Verdrahtung (Dateikommentar
// am Aufrufort).
func awaitSlotExists(t *testing.T, pool *pgxpool.Pool, slot string, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for {
		var exists int
		if err := pool.QueryRow(context.Background(),
			"SELECT 1 FROM pg_replication_slots WHERE slot_name = $1", slot,
		).Scan(&exists); err == nil {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("Slot %q ist nach %s nicht in pg_replication_slots erschienen — Run() nicht hochgefahren", slot, timeout)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// awaitPublishedChangeCount pollt `cdc.change`/`cdc.transaction` direkt
// (`schema.sql`, `SPEC-001`/`SPEC-002`), bis mindestens `want` Changes der
// Quelle sichtbar sind — der Beleg, dass der Capture-Betrieb tatsächlich
// weiterläuft. Die Lese-Sicht `cdc.changes` (`LH-FA-SST-002`, genutzt in
// `welle6_endtoend_test.go`) trägt `make schema-rollout` (d-migrate) und
// steht unter `make test-replication` nicht zur Verfügung — dieser
// Testlauf trägt nur `postgresstorage.ApplySchema` (Dateikommentar oben).
func awaitPublishedChangeCount(t *testing.T, pool *pgxpool.Pool, source model.SourceID, want int, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for {
		var count int
		if err := pool.QueryRow(context.Background(),
			`SELECT count(*) FROM cdc.change c
			 JOIN cdc.transaction t ON t.transaction_id = c.transaction_id
			 WHERE t.source_id = $1`, string(source),
		).Scan(&count); err != nil {
			t.Fatalf("cdc.change zählen: %v", err)
		}
		if count >= want {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("cdc.changes trägt %d Zeilen nach %s, wollen mindestens %d — Capture-Betrieb kam nicht nach", count, timeout, want)
		}
		time.Sleep(200 * time.Millisecond)
	}
}
