package postgresstorage_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Die Backfill-Store-Tests (`backfilladmission_test.go`,
// `backfillrun_test.go`, `backfillwriter_test.go`) laufen gegen dieselbe
// reale PostgreSQL-Instanz wie die übrigen Store-Tests (`make test-store`,
// `ADR-0030`); ohne DSN überspringen sie. `cdc.backfill_run` trägt der
// d-migrate-Rollout (`tools/schema/schema.yaml`, `ADR-0043`), die Grants die
// Rollen-Datei (`tools/schema/nacharbeit-roles.sql`).
//
// Kopplung: diese Dateien und `roles_test.go` müssen vor `store_test.go`
// (`newTestStore`, `DROP SCHEMA cdc CASCADE` samt hand-DDL-Neuaufbau) und
// `tableactivation_test.go` (`newTestActivation`, derselbe Neuaufbau) laufen
// — der Neuaufbau trägt weder `cdc.backfill_run` noch
// `cdc.administration_request` mit. Die Dateinamen `backfill*_test.go` und
// `roles_test.go` sortieren vor ihnen (`b`, `r` < `s`, `t`), dieselbe
// Ordnung wie bei `administrationrequest_test.go`.
//
// Isolation (`BEO-PGC/test-isolation-geteilter-zustand`): jeder Test führt
// seine Zeilen unter eigenen Kennungen (Run-, Antrags-, Tabellen- und
// Transaktions-Kennungen mit dem Testnamen als Präfix) und räumt genau diese
// wieder ab.
const backfillTestSource = "src-backfill"

// backfillFixture trägt den Pool der Test-Instanz (Superuser des
// Testcontainers) und die Bereinigung der Zeilen eines Tests.
type backfillFixture struct {
	t    *testing.T
	pool *pgxpool.Pool
	dsn  string
}

func newBackfillFixture(t *testing.T) *backfillFixture {
	t.Helper()
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

	var tables int
	if err := pool.QueryRow(ctx,
		"SELECT count(*) FROM information_schema.tables WHERE table_schema = 'cdc' AND table_name = 'backfill_run'",
	).Scan(&tables); err != nil {
		t.Fatalf("Tabellen-Prüfung: %v", err)
	}
	if tables != 1 {
		t.Fatalf("cdc.backfill_run fehlt — der Schema-Rollout über make schema-rollout trägt sie (SPEC-029, tools/schema/schema.yaml); der test-store-Lauf rollt sie vor dem Testlauf aus")
	}
	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.source (source_id, name) VALUES ($1, 'Backfill-Quelle') ON CONFLICT (source_id) DO NOTHING", backfillTestSource,
	); err != nil {
		t.Fatalf("Quelle-Zeile: %v", err)
	}
	return &backfillFixture{t: t, pool: pool, dsn: dsn}
}

// exec führt eine Anweisung als Superuser aus und bricht den Test bei einem
// Fehler ab.
func (f *backfillFixture) exec(sql string, args ...any) {
	f.t.Helper()
	if _, err := f.pool.Exec(context.Background(), sql, args...); err != nil {
		f.t.Fatalf("%s: %v", sql, err)
	}
}

// scan liest eine Zeile als Superuser in die Ziele.
func (f *backfillFixture) scan(sql string, args []any, dest ...any) {
	f.t.Helper()
	if err := f.pool.QueryRow(context.Background(), sql, args...).Scan(dest...); err != nil {
		f.t.Fatalf("%s: %v", sql, err)
	}
}

// count zählt die Zeilen einer Abfrage `SELECT count(*) …`.
func (f *backfillFixture) count(sql string, args ...any) int {
	f.t.Helper()
	var n int
	f.scan(sql, args, &n)
	return n
}

// pendingRequest legt einen offenen Antrag an und räumt ihn samt der Run-Zeile
// gleicher Kennung ab. Die Antragsart ist `backfill`, die Annahme prüft sie.
func (f *backfillFixture) pendingRequest(id, schema, table string) {
	f.t.Helper()
	f.pendingRequestOfKind(id, schema, table, "backfill")
}

// pendingRequestOfKind legt einen offenen Antrag einer Antragsart an.
func (f *backfillFixture) pendingRequestOfKind(id, schema, table, kind string) {
	f.t.Helper()
	f.exec(`INSERT INTO cdc.administration_request
	    (administration_request_id, source_id, schema_name, table_name, request_kind, status)
	    VALUES ($1, $2, $3, $4, $5, 'pending')`, id, backfillTestSource, schema, table, kind)
	f.t.Cleanup(func() {
		_, _ = f.pool.Exec(context.Background(), "DELETE FROM cdc.backfill_run WHERE run_id = $1", id)
		_, _ = f.pool.Exec(context.Background(), "DELETE FROM cdc.administration_request WHERE administration_request_id = $1", id)
	})
}

// requestStatus liest den Status eines Antrags.
func (f *backfillFixture) requestStatus(id string) string {
	f.t.Helper()
	var status string
	f.scan("SELECT status FROM cdc.administration_request WHERE administration_request_id = $1", []any{id}, &status)
	return status
}

// runStatus liest den Status einer Run-Zeile; `ok` ist falsch, wenn die Zeile
// nicht besteht.
func (f *backfillFixture) runStatus(id string) (status string, ok bool) {
	f.t.Helper()
	rows, err := f.pool.Query(context.Background(), "SELECT status FROM cdc.backfill_run WHERE run_id = $1", id)
	if err != nil {
		f.t.Fatalf("Run-Status: %v", err)
	}
	defer rows.Close()
	if !rows.Next() {
		return "", false
	}
	if err := rows.Scan(&status); err != nil {
		f.t.Fatalf("Run-Status: %v", err)
	}
	return status, true
}

// queuedRun baut den Run, den die Annahme einträgt.
func queuedRun(t *testing.T, id, table string, requestedAt time.Time, estimate model.RowEstimate) model.BackfillRun {
	t.Helper()
	run, err := model.NewQueuedBackfillRun(model.BackfillRunID(id), backfillTestSource, "public", table, model.NewTimePoint(requestedAt.UnixNano()))
	if err != nil {
		t.Fatalf("NewQueuedBackfillRun: %v", err)
	}
	return run.WithEstimatedRows(estimate)
}

func knownEstimate(t *testing.T, rows int64) model.RowEstimate {
	t.Helper()
	estimate, err := model.NewRowEstimate(rows)
	if err != nil {
		t.Fatalf("NewRowEstimate: %v", err)
	}
	return estimate
}

// admit legt Antrag und Run an und nimmt ihn über den Annahme-Adapter an.
func (f *backfillFixture) admit(admission outbound.BackfillAdmissionPort, id, table string, requestedAt time.Time, estimate model.RowEstimate) model.BackfillRun {
	f.t.Helper()
	f.pendingRequest(id, "public", table)
	run := queuedRun(f.t, id, table, requestedAt, estimate)
	if err := admission.Admit(context.Background(), model.AdministrationRequestID(id), run); err != nil {
		f.t.Fatalf("Admit %s: %v", id, err)
	}
	return run
}

// position trägt die Snapshot-Position `X` der Quelle.
func backfillPosition(t *testing.T, offset uint64) model.SourcePosition {
	t.Helper()
	position, err := model.NewSourcePosition(backfillTestSource, offset)
	if err != nil {
		t.Fatalf("NewSourcePosition: %v", err)
	}
	return position
}

func runID(prefix string, n int) string { return fmt.Sprintf("%s-%d", prefix, n) }
