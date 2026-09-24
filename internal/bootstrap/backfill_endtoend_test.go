package bootstrap

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage"
	"github.com/pt9912/pg-change-feed/internal/adapters/driven/systemclock"
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
