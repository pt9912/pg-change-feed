package bootstrap_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/mapper"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/bootstrap"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// TestWALRetentionThresholdEndToEnd trägt den Ende-zu-Ende-Beleg der
// WAL-Rückstand-Schwellen (`ADR-0049`, `SPEC-013`): ein real wachsender
// WAL-Rückstand durchläuft beide Seiten der Schwelle in einem Lauf von
// `bootstrap.Run`. Der Rückstand wächst bei einem Feed, der nicht
// antwortet (`LH-QA-REL-003`): eine Zeilensperre des Tests auf der
// Schema-Version der aktivierten Tabelle lässt die Persistenz des
// Capture-Aufrufs warten, der Stream bestätigt dann nichts und liest nichts,
// während Transaktionen auf einer *nicht* publizierten Tabelle das WAL
// wachsen lassen — `WALRetentionChecker.Measure` (die Differenz aus
// `IDENTIFY_SYSTEM` und `confirmed_flush_lsn`) sieht genau diesen
// Rückstand. WAL ohne Inhalt für die Publication allein lässt ihn bei einem
// antwortenden Feed nicht wachsen (Leerlauf-Bestätigung, `ADR-0120`). Die
// Schwellen selbst sind ein Test-Override
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

	// Den Schema-Stand dieses Laufs trägt die eine Schema-Anwendung der
	// Test-Läufe (`tools/schema/apply-rollout.sh`), die der Lauf-Aufruf vor
	// dem Tier-Lauf ausführt — derselbe d-migrate-Rollout wie im Betrieb und
	// in `make test-store`. Er trägt die Tabellen, die `bootstrap.Run` liest:
	// `cdc.table_schema` (`SchemaStorePort.CurrentVersion`),
	// `cdc.administration_request`
	// (`ColumnExclusionPort.ExcludedColumns`, `ADR-0050`) und
	// `cdc.process_heartbeat` (`HeartbeatPort`).
	//
	// Die Zeile `cdc.source` trägt der Aufrufer vor der Aktivierung
	// (`TableActivationAdapter.Register`); dieser Lauf legt sie für seine
	// Quelle an.
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

	// holdCapture hält eine Zeilensperre auf der Schema-Version der
	// aktivierten Tabelle: die Persistenz jeder Transaktion mit einem Change
	// dieser Tabelle wartet auf die Fremdschlüssel-Prüfung, der Feed
	// antwortet nicht mehr. Der Rückgabewert gibt die Sperre frei.
	holdCapture := func() func() {
		lockTx, err := pool.Begin(ctx)
		if err != nil {
			t.Fatalf("Sperr-Transaktion: %v", err)
		}
		release := func() { _ = lockTx.Rollback(context.Background()) }
		t.Cleanup(release)
		if _, err := lockTx.Exec(ctx,
			"SELECT 1 FROM cdc.schema_version WHERE schema_version_id = $1 FOR UPDATE", schemaVerID,
		); err != nil {
			t.Fatalf("Sperre auf der Schema-Version: %v", err)
		}
		return release
	}

	// Phase 1 (Seite 1): WAL-Rückstand real über die Warnschwelle wachsen
	// lassen, aber deutlich unter der Fehlerschwelle halten (Größenordnung:
	// ~150-200 KiB gegen 32 KiB Warn-/512 KiB Fehlerschwelle). Der Feed
	// wartet in der Persistenz einer Change der publizierten Tabelle; die
	// Transaktion auf der nicht publizierten Tabelle liegt hinter ihr im WAL
	// und wird nicht bestätigt.
	releaseCapture := holdCapture()
	if _, err := pool.Exec(ctx, fmt.Sprintf("INSERT INTO %s (id) VALUES (1)", published)); err != nil {
		t.Fatalf("zweite Change auf der publizierten Tabelle: %v", err)
	}
	if _, err := pool.Exec(ctx,
		fmt.Sprintf("INSERT INTO %s (id, payload) SELECT g, repeat('a', 800) FROM generate_series(1, 200) g", noise),
	); err != nil {
		t.Fatalf("Phase 1 (Warnschwelle) Noise-Wachstum: %v", err)
	}

	// Belegt „kontrollierte Fortsetzung" (SPEC-008): über mehrere Sekunden
	// hinweg kehrt `Run` nicht zurück: `heartbeatInterval` (5s,
	// `internal/bootstrap`) ist der Takt des periodischen
	// Schwellen-Vergleichs — dieses Warten deckt mindestens einen Tick im
	// Warn-Bereich ab, ohne dass er den Lauf beendet.
	stillRunningDeadline := time.Now().Add(8 * time.Second)
	for time.Now().Before(stillRunningDeadline) {
		select {
		case err := <-runDone:
			t.Fatalf("Run() ist zwischen Warn- und Fehlerschwelle beendet (%v) — SPEC-008 verlangt kontrollierte Fortsetzung, keinen Abbruch", err)
		case <-time.After(200 * time.Millisecond):
		}
	}
	// Gibt der Test die Sperre frei, nimmt der Capture-Betrieb die
	// wartende Change nach dem Warten auf: der Capture-Pfad lief über der
	// Warnschwelle weiter.
	releaseCapture()
	awaitPublishedChangeCount(t, pool, source, 2, 30*time.Second)

	// Phase 2 (Seite 2): WAL-Rückstand real über die Fehlerschwelle
	// wachsen lassen (Größenordnung: ~2 MiB gegen 512 KiB Fehlerschwelle,
	// oben drauf) — derselbe Mechanismus, deutlich größer.
	holdCapture()
	if _, err := pool.Exec(ctx, fmt.Sprintf("INSERT INTO %s (id) VALUES (2)", published)); err != nil {
		t.Fatalf("dritte Change auf der publizierten Tabelle: %v", err)
	}
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
// (`SPEC-001`/`SPEC-002`), bis mindestens `want` Changes der Quelle
// sichtbar sind — der Beleg, dass der Capture-Betrieb tatsächlich
// weiterläuft. Der Zählstand kommt damit aus den Persistenz-Tabellen
// selbst und nicht aus der Projektion der Lese-Sicht `cdc.changes`
// (`LH-FA-SST-002`).
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
