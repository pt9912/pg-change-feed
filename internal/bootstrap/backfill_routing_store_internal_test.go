package bootstrap

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage"
	"github.com/pt9912/pg-change-feed/internal/adapters/driven/systemclock"
	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/backfill"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// routeStoreSnapshot liefert die Zeilen in Blöcken; `beforeBlock` läuft am
// Anfang jedes `NextBlock`-Aufrufs (ab 1).
type routeStoreSnapshot struct {
	columns     []string
	blocks      [][][]*string
	beforeBlock func(call int)
	calls       int
	closed      int
}

func (s *routeStoreSnapshot) Offset() uint64    { return 0x200 }
func (s *routeStoreSnapshot) Columns() []string { return s.columns }
func (s *routeStoreSnapshot) NextBlock(context.Context) ([][]*string, error) {
	s.calls++
	if s.beforeBlock != nil {
		s.beforeBlock(s.calls)
	}
	if s.calls > len(s.blocks) {
		return nil, nil
	}
	return s.blocks[s.calls-1], nil
}
func (s *routeStoreSnapshot) Close(context.Context) error { s.closed++; return nil }

type routeStoreSnapshotPort struct{ snapshot *routeStoreSnapshot }

func (p *routeStoreSnapshotPort) OpenSnapshot(context.Context, string, string, string) (outbound.TableSnapshot, error) {
	return p.snapshot, nil
}
func (p *routeStoreSnapshotPort) EstimatedRows(context.Context, string, string) (int64, bool, error) {
	return 0, false, nil
}

// TestBackfillRunAgainstPostgreSQLCarriesRoutingStateAndLabel belegt den Run
// mit dem Routing-Regelstand gegen die reale Run-Tabelle, Antrags-Tabelle und
// Change-Tabelle (`LH-FA-CFG-008`): Bindung, Schema-Version, Run-Zustand,
// Schreiber und der Regelstand aus den `applied`-Zeilen von
// `cdc.administration_request` sind die echten Adapter; der Snapshot ist ein
// Test-Port mit zwei Blöcken (die Snapshot-Adapter-Tests laufen unter
// `make test-replication`). Die Fälle:
//   - ohne Wechsel trägt jede Change in `cdc.changes` das Ziel der Regel
//     (`route_target`, Treffer, abwesender Wert, Nicht-Treffer) — die
//     Persistenz des Labels im Backfill-Insert;
//   - ein `set_route` bzw. `remove_route` (SQL-Funktion, vermerkt `applied`)
//     zwischen den beiden Blöcken, der Run endet `failed` mit der Klasse
//     `configuration`, es steht keine Change des Runs in `cdc.change`;
//   - eine Regel auf eine Spalte, die der Snapshot nicht trägt, endet den Run
//     `failed` mit der Klasse `schema` vor der ersten Zeile, ohne Change.
//
// Rot färbende Mutationen (je eine, an der Eingabeseite): die Spalte
// `route_target` aus `InsertBackfillChange` streichen — die Ziele in
// `cdc.changes` sind leer; den Vergleich des Routing-Standes je Block im Use
// Case entfernen — die beiden Wechsel-Fälle enden `completed`; die Prüfung der
// Anwendbarkeit entfernen — der Fall „Spalte fehlt“ endet `completed`;
// `'applied'` im `WHERE` von `SelectAppliedRoutingRequests` durch `'pending'`
// ersetzen — der Stand ist leer, alle Ziele fehlen.
func TestBackfillRunAgainstPostgreSQLCarriesRoutingStateAndLabel(t *testing.T) {
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
		source      = model.SourceID("src-backfill-route-store")
		table       = "bf_route_store"
		publication = "pgc_bf_route_store_pub"
		tableID     = model.SourceTableID("tbl-bf-route-store")
	)
	exec := func(sql string, args ...any) {
		t.Helper()
		if _, err := pool.Exec(ctx, sql, args...); err != nil {
			t.Fatalf("%s: %v", sql, err)
		}
	}
	exec("INSERT INTO cdc.source (source_id, name) VALUES ($1, 'Backfill-Routing') ON CONFLICT (source_id) DO NOTHING", string(source))
	exec("CREATE TABLE public." + table + " (id integer PRIMARY KEY, region text)")
	exec("CREATE PUBLICATION " + publication + " FOR TABLE public." + table)
	t.Cleanup(func() {
		bg := context.Background()
		_, _ = pool.Exec(bg, "DELETE FROM cdc.change WHERE transaction_id LIKE '0bf-bfroute-%'")
		_, _ = pool.Exec(bg, "DELETE FROM cdc.transaction WHERE transaction_id LIKE '0bf-bfroute-%'")
		_, _ = pool.Exec(bg, "DELETE FROM cdc.backfill_run WHERE source_id = $1", string(source))
		_, _ = pool.Exec(bg, "DELETE FROM cdc.administration_request WHERE source_id = $1", string(source))
		_, _ = pool.Exec(bg, "DELETE FROM cdc.schema_version WHERE source_table_id = $1", string(tableID))
		_, _ = pool.Exec(bg, "DELETE FROM cdc.source_table WHERE source_id = $1", string(source))
		_, _ = pool.Exec(bg, "DELETE FROM cdc.source WHERE source_id = $1", string(source))
		_, _ = pool.Exec(bg, "DROP PUBLICATION IF EXISTS "+publication)
		_, _ = pool.Exec(bg, "DROP TABLE IF EXISTS public."+table)
	})

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

	sourceTable, err := model.NewSourceTable(tableID, source, "public", table)
	if err != nil {
		t.Fatal(err)
	}
	version, err := model.NewSchemaVersion("tbl-bf-route-store-v1", tableID, 1)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := activation.Register(ctx, sourceTable, version); err != nil {
		t.Fatalf("Register: %v", err)
	}

	// applyRoute legt den Antrag über die SQL-Funktion an und vermerkt ihn
	// `applied`, wie es die Antragsverarbeitung tut.
	applyRoute := func(call string, args ...any) {
		t.Helper()
		var id string
		if err := pool.QueryRow(ctx, call, args...).Scan(&id); err != nil {
			t.Fatalf("%s: %v", call, err)
		}
		exec("UPDATE cdc.administration_request SET status = 'applied' WHERE administration_request_id = $1", id)
	}
	setRoute := func(name, spec string) {
		t.Helper()
		applyRoute("SELECT cdc.set_route($1, 'public', $2, $3, $4::json)", string(source), table, name, spec)
	}
	removeRoute := func(name string) {
		t.Helper()
		applyRoute("SELECT cdc.remove_route($1, 'public', $2, $3)", string(source), table, name)
	}
	clearRoutes := func() {
		exec("DELETE FROM cdc.administration_request WHERE source_id = $1", string(source))
	}

	columns := []string{"id", "region"}
	rows := func() [][][]*string {
		return [][][]*string{
			{{parityValue("1"), parityValue("eu")}, {parityValue("2"), parityValue("us")}},
			{{parityValue("3"), parityValue("eu")}, {parityValue("4"), nil}},
		}
	}
	euRule := `{"target": "eu_ziel", "order": 10, "when": {"column": "region", "equals": "eu"}}`
	runCounter := 0
	execute := func(snapshot *routeStoreSnapshot) (runID string, run model.BackfillRun) {
		t.Helper()
		runCounter++
		runID = fmt.Sprintf("bfroute-%d", runCounter)
		exec("INSERT INTO cdc.backfill_run (run_id, source_id, schema_name, table_name, status) VALUES ($1, $2, 'public', $3, 'queued')", runID, string(source), table)
		service := backfill.NewBackfillTableService(backfill.Ports{
			Activation: activation, Exclusion: activation, Transformations: activation, Routing: activation, Schemas: schemaStore,
			Snapshot: &routeStoreSnapshotPort{snapshot: snapshot}, Admission: admission, Runs: runs, Writer: writer, Clock: systemclock.New(),
		}, backfill.WithLog(log))
		queued, err := model.NewQueuedBackfillRun(model.BackfillRunID(runID), source, "public", table, systemclock.New().Now())
		if err != nil {
			t.Fatal(err)
		}
		result, err := service.Execute(ctx, inbound.BackfillExecuteCommand{Run: queued, Publication: publication})
		if err != nil {
			t.Fatalf("Execute: %v — Log: %v", err, log.messages)
		}
		return runID, result.Run
	}
	targetsOf := func(runID string) []string {
		t.Helper()
		found, err := pool.Query(ctx,
			"SELECT COALESCE(route_target, '-') FROM cdc.changes WHERE transaction_id LIKE '0bf-' || $1 || '-%' ORDER BY transaction_id, sequence", runID)
		if err != nil {
			t.Fatalf("cdc.changes lesen: %v", err)
		}
		defer found.Close()
		var targets []string
		for found.Next() {
			var target string
			if err := found.Scan(&target); err != nil {
				t.Fatal(err)
			}
			targets = append(targets, target)
		}
		return targets
	}
	persisted := func(runID string) (status, message string) {
		t.Helper()
		if err := pool.QueryRow(ctx, "SELECT status, COALESCE(error_message, '') FROM cdc.backfill_run WHERE run_id = $1", runID).Scan(&status, &message); err != nil {
			t.Fatalf("Run-Zeile lesen: %v", err)
		}
		return status, message
	}

	t.Run("ohne Wechsel trägt jede Change das Ziel", func(t *testing.T) {
		clearRoutes()
		setRoute("eu", euRule)
		runID, run := execute(&routeStoreSnapshot{columns: columns, blocks: rows()})
		if run.Status != model.BackfillRunCompleted {
			t.Fatalf("Run = %+v, will completed", run)
		}
		if got := fmt.Sprint(targetsOf(runID)); got != "[eu_ziel - eu_ziel -]" {
			t.Fatalf("route_target in cdc.changes = %s, will [eu_ziel - eu_ziel -]", got)
		}
	})

	t.Run("set_route zwischen den Blöcken", func(t *testing.T) {
		clearRoutes()
		snapshot := &routeStoreSnapshot{columns: columns, blocks: rows()}
		snapshot.beforeBlock = func(call int) {
			if call == 2 {
				setRoute("eu", euRule)
			}
		}
		runID, run := execute(snapshot)
		status, message := persisted(runID)
		if run.Status != model.BackfillRunFailed || status != "failed" || !strings.HasPrefix(message, "configuration: ") || !strings.Contains(message, "Routing-Regelstand") {
			t.Fatalf("Run = %+v, Zeile %q/%q, will failed mit Klasse configuration und dem Routing-Regelstand", run, status, message)
		}
		if got := targetsOf(runID); len(got) != 0 {
			t.Fatalf("Changes des fehlgeschlagenen Runs in cdc.changes: %v", got)
		}
		if snapshot.closed != 1 {
			t.Fatalf("Snapshot %d-mal geschlossen, will 1", snapshot.closed)
		}
	})

	t.Run("remove_route zwischen den Blöcken", func(t *testing.T) {
		clearRoutes()
		setRoute("eu", euRule)
		snapshot := &routeStoreSnapshot{columns: columns, blocks: rows()}
		snapshot.beforeBlock = func(call int) {
			if call == 2 {
				removeRoute("eu")
			}
		}
		runID, run := execute(snapshot)
		status, message := persisted(runID)
		if run.Status != model.BackfillRunFailed || status != "failed" || !strings.HasPrefix(message, "configuration: ") || !strings.Contains(message, "Routing-Regelstand") {
			t.Fatalf("Run = %+v, Zeile %q/%q, will failed mit Klasse configuration und dem Routing-Regelstand", run, status, message)
		}
		if got := targetsOf(runID); len(got) != 0 {
			t.Fatalf("Changes des fehlgeschlagenen Runs in cdc.changes: %v", got)
		}
	})

	t.Run("Spalte fehlt im Snapshot", func(t *testing.T) {
		clearRoutes()
		setRoute("fehlt", `{"target": "x_ziel", "order": 10, "when": {"column": "gibt_es_nicht", "equals": "v"}}`)
		snapshot := &routeStoreSnapshot{columns: columns, blocks: rows()}
		runID, run := execute(snapshot)
		status, message := persisted(runID)
		if run.Status != model.BackfillRunFailed || status != "failed" || !strings.HasPrefix(message, "schema: ") ||
			!strings.Contains(message, `"fehlt"`) || !strings.Contains(message, `"gibt_es_nicht"`) {
			t.Fatalf("Run = %+v, Zeile %q/%q, will failed mit Klasse schema, Regelname und Spalte", run, status, message)
		}
		if snapshot.calls != 0 || snapshot.closed != 1 {
			t.Fatalf("NextBlock %d-mal, Snapshot %d-mal geschlossen, will 0 und 1", snapshot.calls, snapshot.closed)
		}
		if got := targetsOf(runID); len(got) != 0 {
			t.Fatalf("Changes des fehlgeschlagenen Runs in cdc.changes: %v", got)
		}
	})
}
