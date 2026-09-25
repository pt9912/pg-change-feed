package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/mapper"
	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/receive"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
)

// TestWALRetentionThresholdEndToEnd trägt den Ende-zu-Ende-Beleg der
// WAL-Rückstand-Schwellen (`ADR-0049`, `SPEC-013`) gegen reale PostgreSQL: ein
// real wachsender Rückstand durchläuft beide Seiten der Schwelle, gemessen
// vom echten `WALRetentionChecker` und bewertet von `runWALRetentionCheck`
// samt Rückgabewert-Priorität. Der Rückstand wächst bei einem Slot, den kein
// Stream bestätigt (`LH-QA-REL-003`: ein Feed, der nicht antwortet, lässt ihn
// wachsen) — WAL ohne Inhalt für die Publication allein lässt ihn bei einem
// antwortenden Feed nicht wachsen (Leerlauf-Bestätigung, `ADR-0120`). Die
// Schwellen sind ein Test-Override (`Config.WALRetentionWarnBytes`/
// `WALRetentionErrorBytes`), klein gegen `SPEC-013`s 100 MiB/1 GiB und groß
// gegen das Grundrauschen der Instanz.
//
// Die Instanz gehört dem Test allein (`CDC_WALRETENTION_TEST_DSN`,
// `run-replication-tests.sh` startet den Lauf gesondert nach dem Tier-Lauf):
// der Rückstand misst das WAL der ganzen Instanz, ein gleichzeitiger Schreiber
// eines anderen Test-Pakets läge sonst in derselben Größenordnung wie die
// Schwellen.
func TestWALRetentionThresholdEndToEnd(t *testing.T) {
	dsn := os.Getenv("CDC_WALRETENTION_TEST_DSN")
	if dsn == "" {
		t.Skip("CDC_WALRETENTION_TEST_DSN nicht gesetzt — der Beleg läuft auf einer exklusiven Instanz über make test-replication")
	}
	ctx, cancelRun := context.WithCancel(context.Background())
	defer cancelRun()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("Verbindungsaufbau: %v", err)
	}
	t.Cleanup(pool.Close)

	const (
		noise = "public.wal_e2e_noise"
		slot  = "slot_wal_e2e"
		// warnBytes/errorBytes sind ein Test-Fixture (Testkommentar oben) —
		// nicht die SPEC-013-Startwerte (100 MiB/1 GiB, `ADR-0049`(b)).
		warnBytes  int64 = 32 * 1024
		errorBytes int64 = 2 * 1024 * 1024
	)
	for _, stmt := range []string{
		"DROP TABLE IF EXISTS " + noise,
		fmt.Sprintf("CREATE TABLE %s (id int PRIMARY KEY, payload text)", noise),
		fmt.Sprintf("SELECT pg_create_logical_replication_slot('%s', 'pgoutput')", slot),
	} {
		if _, err := pool.Exec(ctx, stmt); err != nil {
			t.Fatalf("Test-Umgebung (%s): %v", stmt, err)
		}
	}
	t.Cleanup(func() {
		dropCtx, dropCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer dropCancel()
		_, _ = pool.Exec(dropCtx, fmt.Sprintf("SELECT pg_drop_replication_slot('%s')", slot))
		_, _ = pool.Exec(dropCtx, "DROP TABLE IF EXISTS "+noise)
	})

	checker, err := receive.NewWALRetentionChecker(ctx, dsn, slot)
	if err != nil {
		t.Fatalf("NewWALRetentionChecker: %v", err)
	}
	t.Cleanup(func() { _ = checker.Close(context.Background()) })
	if baseline, err := checker.Measure(ctx); err != nil || baseline >= warnBytes {
		t.Fatalf("Rückstand vor der Last %d B (Fehler %v), erwartet unter der Warnschwelle %d B — die Instanz ist nicht exklusiv", baseline, err, warnBytes)
	}

	log := &recordingLog{}
	streamCtx, stopStream := context.WithCancel(ctx)
	defer stopStream()
	var fault walRetentionFault
	checkDone := make(chan struct{})
	go func() {
		defer close(checkDone)
		runWALRetentionCheck(ctx, checker, log, 100*time.Millisecond, warnBytes, errorBytes, stopStream, &fault)
	}()

	// Seite 1: der Rückstand liegt über der Warnschwelle und deutlich unter
	// der Fehlerschwelle (Größenordnung ~170 KiB gegen 32 KiB Warn- und
	// 2 MiB Fehlerschwelle).
	if _, err := pool.Exec(ctx,
		fmt.Sprintf("INSERT INTO %s (id, payload) SELECT g, repeat('a', 800) FROM generate_series(1, 200) g", noise)); err != nil {
		t.Fatalf("Seite 1 (Warnschwelle): %v", err)
	}
	warnDeadline := time.Now().Add(10 * time.Second)
	for !log.contains("WARN", "über Warnschwelle") {
		if time.Now().After(warnDeadline) {
			t.Fatalf("keine Warnung über der Warnschwelle nach 10 s (Log: %v)", logLines(log))
		}
		time.Sleep(50 * time.Millisecond)
	}
	// „Kontrollierte Fortsetzung" (`SPEC-008`): über mehrere weitere Ticks
	// endet weder die Schleife noch der Stream-Lauf.
	time.Sleep(600 * time.Millisecond)
	select {
	case <-checkDone:
		t.Fatalf("die Prüf-Schleife endete zwischen Warn- und Fehlerschwelle (%v) — SPEC-008 verlangt kontrollierte Fortsetzung", fault.get())
	case <-streamCtx.Done():
		t.Fatalf("der Stream-Lauf wurde zwischen Warn- und Fehlerschwelle beendet")
	default:
	}

	// Seite 2: der Rückstand liegt über der Fehlerschwelle (Größenordnung
	// ~6,6 MiB gegen 2 MiB) — derselbe Mechanismus, deutlich größer.
	if _, err := pool.Exec(ctx,
		fmt.Sprintf("INSERT INTO %s (id, payload) SELECT g, repeat('a', 800) FROM generate_series(1000, 9000) g", noise)); err != nil {
		t.Fatalf("Seite 2 (Fehlerschwelle): %v", err)
	}
	select {
	case <-streamCtx.Done():
	case <-time.After(20 * time.Second):
		t.Fatal("der Stream-Lauf wurde nach Überschreiten der Fehlerschwelle nicht innerhalb 20 s beendet")
	}
	<-checkDone

	err = fault.get()
	if err == nil || !errors.Is(err, outbound.ErrReplication) {
		t.Fatalf("Schwellen-Fehler %v, wollen einen Fehler der Klasse replication (outbound.ErrReplication) — derselbe Pfad wie classifyRunError/os.Exit(1) in main.go", err)
	}
	if errors.Is(err, mapper.ErrChangeWithoutBegin) || errors.Is(err, mapper.ErrCommitWithoutBegin) || errors.Is(err, mapper.ErrBeginWithoutCommit) {
		t.Fatalf("Schwellen-Fehler ist eine Stream-Ordnungs-Verletzung (%v) — Sentinel-Trennung verletzt (ADR-0049(a))", err)
	}
	if merged := mergeStreamAndWALFaultOutcome(nil, &fault); !errors.Is(merged, outbound.ErrReplication) {
		t.Fatalf("Rückgabewert von Run bei regulärem Stream-Ende %v, wollen die Klasse replication", merged)
	}
	if !log.contains("ERROR", "über Fehlerschwelle") {
		t.Fatalf("keine Fehler-Zeile über der Fehlerschwelle (Log: %v)", logLines(log))
	}
}

// logLines trägt die protokollierten Zeilen des `recordingLog`.
func logLines(log *recordingLog) []string {
	log.mu.Lock()
	defer log.mu.Unlock()
	return append([]string(nil), log.messages...)
}
