package bootstrap

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgressnapshot"
	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage"
	"github.com/pt9912/pg-change-feed/internal/adapters/driven/systemclock"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/backfill"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// TestReconcileBackfillRunsAgainstPostgreSQL trägt den Abgleich beim
// Prozessstart gegen die reale Run-Tabelle (`ADR-0111` Teilfrage 4,
// `ADR-0113` Festlegung 2): ein `running`-Run der eigenen Quelle wird
// `interrupted`; ein `queued`-Run derselben Quelle bleibt `queued` (der Worker
// nimmt ihn beim Start auf), ein bereits `interrupted`, `completed` oder
// `failed` beendeter Run bleibt unverändert, und ein `running`-Run einer
// anderen Quelle bleibt unberührt. Rot färbende Mutationen (je eine): eine
// andere Quelle an `InterruptRunning` übergeben — der `running`-Run der
// eigenen Quelle bleibt `running`; die Quelle nicht übergeben (Abgleich über
// alle Quellen) — der Run der fremden Quelle wird `interrupted`.
func TestReconcileBackfillRunsAgainstPostgreSQL(t *testing.T) {
	dsn := os.Getenv("CDC_STORE_TEST_DSN")
	if dsn == "" {
		t.Skip("CDC_STORE_TEST_DSN nicht gesetzt — reale PostgreSQL-Tests laufen über make test-store")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("Verbindungsaufbau: %v", err)
	}
	t.Cleanup(pool.Close)

	const (
		own   = model.SourceID("src-reconcile-own")
		other = model.SourceID("src-reconcile-other")
	)
	for _, source := range []model.SourceID{own, other} {
		if _, err := pool.Exec(ctx,
			"INSERT INTO cdc.source (source_id, name) VALUES ($1, 'Abgleich-Quelle') ON CONFLICT (source_id) DO NOTHING", string(source),
		); err != nil {
			t.Fatalf("Quelle-Zeile: %v", err)
		}
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM cdc.backfill_run WHERE run_id LIKE 'reconcile-%'")
		_, _ = pool.Exec(context.Background(), "DELETE FROM cdc.source WHERE source_id IN ($1, $2)", string(own), string(other))
	})
	for _, run := range []struct{ id, source, status string }{
		{"reconcile-running", string(own), "running"},
		{"reconcile-queued", string(own), "queued"},
		{"reconcile-interrupted", string(own), "interrupted"},
		{"reconcile-completed", string(own), "completed"},
		{"reconcile-failed", string(own), "failed"},
		{"reconcile-other-running", string(other), "running"},
	} {
		if _, err := pool.Exec(ctx,
			`INSERT INTO cdc.backfill_run (run_id, source_id, schema_name, table_name, status) VALUES ($1, $2, 'public', $1, $3)`,
			run.id, run.source, run.status,
		); err != nil {
			t.Fatalf("Run-Zeile %s: %v", run.id, err)
		}
	}

	runs, err := postgresstorage.NewBackfillRun(ctx, dsn)
	if err != nil {
		t.Fatalf("NewBackfillRun: %v", err)
	}
	t.Cleanup(runs.Close)
	if err := reconcileBackfillRuns(ctx, runs, systemclock.New(), own, &recordingLog{}); err != nil {
		t.Fatalf("reconcileBackfillRuns: %v", err)
	}

	for id, want := range map[string]string{
		"reconcile-running":       "interrupted",
		"reconcile-queued":        "queued",
		"reconcile-interrupted":   "interrupted",
		"reconcile-completed":     "completed",
		"reconcile-failed":        "failed",
		"reconcile-other-running": "running",
	} {
		var status string
		if err := pool.QueryRow(ctx, "SELECT status FROM cdc.backfill_run WHERE run_id = $1", id).Scan(&status); err != nil {
			t.Fatalf("Status %s lesen: %v", id, err)
		}
		if status != want {
			t.Fatalf("Run %s nach dem Abgleich: %q, erwartet %q", id, status, want)
		}
	}
	var finished *string
	if err := pool.QueryRow(ctx, "SELECT finished_at::text FROM cdc.backfill_run WHERE run_id = 'reconcile-running'").Scan(&finished); err != nil || finished == nil {
		t.Fatalf("finished_at des abgeglichenen Runs: %v (%v), erwartet gesetzt", finished, err)
	}
}

// TestBackfillWorkerEndsARunOfAnUnboundTableAsConfigurationFailure trägt den
// Weg Aufnahme → Ausführung gegen die reale Run-Tabelle (`ADR-0113`
// Festlegung 2): eine `queued`-Zeile, deren Tabelle keine Bindung trägt, nimmt
// der Worker beim Start auf; die erneute Prüfung der Vorbedingungen vor der
// Ausführung endet den Run `failed` mit der Fehlerklasse `configuration`, ohne
// Slot und ohne Kopie. Die Adapter laufen unter dem Superuser des
// Testcontainers (kein Replication-Recht nötig: der Run endet vor dem Slot).
// Rot färbende Mutationen (je eine): den Aufruf `deps.useCase.Execute(...)` aus
// `drainBackfillQueue` entfernen — der Run bleibt `queued`, der Test meldet, dass
// der Worker die Zeile nicht aufnahm; `deps.source` beim Lesen der `queued`-Zeilen
// durch eine andere Quelle ersetzen — dieselbe Meldung.
func TestBackfillWorkerEndsARunOfAnUnboundTableAsConfigurationFailure(t *testing.T) {
	dsn := os.Getenv("CDC_STORE_TEST_DSN")
	if dsn == "" {
		t.Skip("CDC_STORE_TEST_DSN nicht gesetzt — reale PostgreSQL-Tests laufen über make test-store")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("Verbindungsaufbau: %v", err)
	}
	t.Cleanup(pool.Close)

	const (
		source = model.SourceID("src-worker-unbound")
		runID  = "worker-unbound-run"
	)
	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.source (source_id, name) VALUES ($1, 'Worker-Quelle') ON CONFLICT (source_id) DO NOTHING", string(source),
	); err != nil {
		t.Fatalf("Quelle-Zeile: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM cdc.backfill_run WHERE run_id = $1", runID)
		_, _ = pool.Exec(context.Background(), "DELETE FROM cdc.source WHERE source_id = $1", string(source))
	})
	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.backfill_run (run_id, source_id, schema_name, table_name, status) VALUES ($1, $2, 'public', 'worker_unbound', 'queued')",
		runID, string(source),
	); err != nil {
		t.Fatalf("Run-Zeile: %v", err)
	}

	log := &recordingLog{}
	activation, err := postgresstorage.NewTableActivation(ctx, dsn, postgresstorage.WithLog(log))
	if err != nil {
		t.Fatalf("NewTableActivation: %v", err)
	}
	t.Cleanup(activation.Close)
	schemaStore, err := postgresstorage.NewSchemaStore(ctx, dsn, postgresstorage.WithLog(log))
	if err != nil {
		t.Fatalf("NewSchemaStore: %v", err)
	}
	t.Cleanup(schemaStore.Close)
	admission, err := postgresstorage.NewBackfillAdmission(ctx, dsn, postgresstorage.WithLog(log))
	if err != nil {
		t.Fatalf("NewBackfillAdmission: %v", err)
	}
	t.Cleanup(admission.Close)
	runs, err := postgresstorage.NewBackfillRun(ctx, dsn, postgresstorage.WithLog(log))
	if err != nil {
		t.Fatalf("NewBackfillRun: %v", err)
	}
	t.Cleanup(runs.Close)
	writer, err := postgresstorage.NewBackfillWriter(ctx, dsn, postgresstorage.WithLog(log))
	if err != nil {
		t.Fatalf("NewBackfillWriter: %v", err)
	}
	t.Cleanup(writer.Close)
	snapshot, err := postgressnapshot.New(dsn)
	if err != nil {
		t.Fatalf("postgressnapshot.New: %v", err)
	}
	useCase := backfill.NewBackfillTableService(backfill.Ports{
		Activation: activation, Exclusion: activation, Schemas: schemaStore, Snapshot: snapshot,
		Admission: admission, Runs: runs, Writer: writer, Clock: systemclock.New(),
	}, backfill.WithLog(log))

	workerCtx, cancel := context.WithCancel(ctx)
	done := make(chan struct{})
	go func() {
		defer close(done)
		runBackfillWorker(workerCtx, backfillWorkerDeps{
			runs: runs, useCase: useCase, source: source, publication: "pub_worker_unbound",
			wake: newBackfillWake(), retryAfter: time.Hour, log: log,
		})
	}()
	t.Cleanup(func() {
		cancel()
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Error("der Worker endete nicht nach dem Abbruch seines Kontexts")
		}
	})

	var status, message string
	var finished *time.Time
	deadline := time.Now().Add(10 * time.Second)
	for {
		if err := pool.QueryRow(ctx,
			"SELECT status, COALESCE(error_message, ''), finished_at FROM cdc.backfill_run WHERE run_id = $1", runID,
		).Scan(&status, &message, &finished); err != nil {
			t.Fatalf("Run-Zeile lesen: %v", err)
		}
		if status != "queued" {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("der Worker nahm die queued-Zeile beim Start nicht auf — Log: %v", log.messages)
		}
		time.Sleep(20 * time.Millisecond)
	}
	if status != "failed" || !strings.HasPrefix(message, "configuration: ") || !strings.Contains(message, "trägt keine Bindung") || finished == nil {
		t.Fatalf("Run nach der Aufnahme: Status %q, Fehlertext %q, finished_at %v — erwartet failed, „configuration: … trägt keine Bindung“ und gesetztes finished_at", status, message, finished)
	}
}
