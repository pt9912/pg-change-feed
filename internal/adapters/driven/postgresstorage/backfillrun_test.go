package postgresstorage_test

import (
	"context"
	stderrors "errors"
	"testing"
	"time"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

func newBackfillRunAdapter(t *testing.T, f *backfillFixture) *postgresstorage.BackfillRunAdapter {
	t.Helper()
	runs, err := postgresstorage.NewBackfillRun(context.Background(), f.dsn)
	if err != nil {
		t.Fatalf("NewBackfillRun: %v", err)
	}
	t.Cleanup(runs.Close)
	return runs
}

// mustRunning führt den Run über `MarkRunning` in `running`.
func mustRunning(t *testing.T, runs outbound.BackfillRunPort, run model.BackfillRun, at time.Time) model.BackfillRun {
	t.Helper()
	started, err := run.Start(model.NewTimePoint(at.UnixNano()))
	if err != nil {
		t.Fatalf("Start: %v", err)
	}
	if err := runs.MarkRunning(context.Background(), started); err != nil {
		t.Fatalf("MarkRunning: %v", err)
	}
	return started
}

// Der Weg eines Runs durch den Run-Zustands-Adapter: `queued` lesen,
// `running` mit `started_at`, Fortschritt (sichtbar für einen zweiten Leser),
// Endzustand `completed` einer leeren Tabelle mit `finished_at`.
func TestBackfillRunFollowsTheRunThroughItsTransitions(t *testing.T) {
	f := newBackfillFixture(t)
	admission := newBackfillAdmission(t, f)
	runs := newBackfillRunAdapter(t, f)
	base := time.Date(2026, 9, 24, 11, 0, 0, 0, time.UTC)

	run := f.admit(admission, "run-flow", "run_flow", base, knownEstimate(t, 10))
	queued, err := runs.Queued(context.Background(), backfillTestSource)
	if err != nil {
		t.Fatalf("Queued: %v", err)
	}
	var listed model.BackfillRun
	for _, candidate := range queued {
		if candidate.ID == "run-flow" {
			listed = candidate
		}
	}
	if listed.ID == "" || listed.Status != model.BackfillRunQueued || listed.Table != "run_flow" || listed.RequestedAt.UnixNanos != base.UnixNano() {
		t.Fatalf("Queued lieferte %+v", listed)
	}

	running := mustRunning(t, runs, run, base.Add(time.Second))
	var status string
	var startedAt time.Time
	f.scan("SELECT status, started_at FROM cdc.backfill_run WHERE run_id = 'run-flow'", nil, &status, &startedAt)
	if status != "running" || !startedAt.Equal(base.Add(time.Second)) {
		t.Fatalf("nach MarkRunning: %s, started_at %v", status, startedAt)
	}
	if again, err := runs.Queued(context.Background(), backfillTestSource); err != nil {
		t.Fatalf("Queued: %v", err)
	} else {
		for _, candidate := range again {
			if candidate.ID == "run-flow" {
				t.Fatal("ein running-Run erscheint in Queued")
			}
		}
	}

	progressed, err := running.RecordProgress(backfillPosition(t, 4242), 7)
	if err != nil {
		t.Fatalf("RecordProgress (Domäne): %v", err)
	}
	if err := runs.RecordProgress(context.Background(), progressed); err != nil {
		t.Fatalf("RecordProgress: %v", err)
	}
	var position, rowsCopied int64
	f.scan("SELECT snapshot_position, rows_copied FROM cdc.backfill_run WHERE run_id = 'run-flow'", nil, &position, &rowsCopied)
	if position != 4242 || rowsCopied != 7 {
		t.Fatalf("nach RecordProgress: Position %d, Zähler %d", position, rowsCopied)
	}

	completed, err := progressed.Complete(model.NewTimePoint(base.Add(2*time.Second).UnixNano()), 7)
	if err != nil {
		t.Fatalf("Complete (Domäne): %v", err)
	}
	if err := runs.Finish(context.Background(), completed); err != nil {
		t.Fatalf("Finish: %v", err)
	}
	var finishedAt time.Time
	f.scan("SELECT status, finished_at, rows_copied, snapshot_position FROM cdc.backfill_run WHERE run_id = 'run-flow'", nil,
		&status, &finishedAt, &rowsCopied, &position)
	if status != "completed" || !finishedAt.Equal(base.Add(2*time.Second)) || rowsCopied != 7 || position != 4242 {
		t.Fatalf("nach Finish: %s, finished_at %v, Zähler %d, Position %d", status, finishedAt, rowsCopied, position)
	}
}

// Die Aufnahme-Ordnung ist `(requested_at, run_id)`: ein früherer Antrag
// zuerst, bei gleichem Zeitstempel die kleinere Run-Kennung; `running`-Runs
// und Runs einer anderen Quelle erscheinen nicht. Rot färbende Mutation:
// `ORDER BY requested_at` ohne `run_id` — die zwei gleichzeitigen Anträge
// kämen in Einfügereihenfolge statt in Kennungs-Reihenfolge.
func TestBackfillRunQueuedIsOrderedByRequestedAtThenRunID(t *testing.T) {
	f := newBackfillFixture(t)
	admission := newBackfillAdmission(t, f)
	runs := newBackfillRunAdapter(t, f)
	base := time.Date(2026, 9, 24, 12, 0, 0, 0, time.UTC)

	// Einfüge-Reihenfolge absichtlich gegen die erwartete Ordnung.
	f.admit(admission, "ord-b", "ord_b", base, model.UnknownRowEstimate())
	f.admit(admission, "ord-c", "ord_c", base.Add(time.Second), model.UnknownRowEstimate())
	f.admit(admission, "ord-a", "ord_a", base, model.UnknownRowEstimate())
	f.admit(admission, "ord-early", "ord_early", base.Add(-time.Second), model.UnknownRowEstimate())
	running := f.admit(admission, "ord-running", "ord_running", base, model.UnknownRowEstimate())
	mustRunning(t, runs, running, base)

	f.exec("INSERT INTO cdc.source (source_id, name) VALUES ('src-backfill-other', 'andere') ON CONFLICT (source_id) DO NOTHING")
	f.t.Cleanup(func() {
		_, _ = f.pool.Exec(context.Background(), "DELETE FROM cdc.backfill_run WHERE run_id = 'ord-other'")
		_, _ = f.pool.Exec(context.Background(), "DELETE FROM cdc.source WHERE source_id = 'src-backfill-other'")
	})
	f.exec(`INSERT INTO cdc.backfill_run (run_id, source_id, schema_name, table_name, status, requested_at)
	        VALUES ('ord-other', 'src-backfill-other', 'public', 'ord_other', 'queued', $1)`, base)

	queued, err := runs.Queued(context.Background(), backfillTestSource)
	if err != nil {
		t.Fatalf("Queued: %v", err)
	}
	var got []model.BackfillRunID
	for _, run := range queued {
		switch run.ID {
		case "ord-early", "ord-a", "ord-b", "ord-c", "ord-running", "ord-other":
			got = append(got, run.ID)
		}
	}
	want := []model.BackfillRunID{"ord-early", "ord-a", "ord-b", "ord-c"}
	if len(got) != len(want) {
		t.Fatalf("Queued = %v, erwartet %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("Queued = %v, erwartet %v", got, want)
		}
	}
}

// Die Endzustände und ihre Ausgangszustände: `failed` folgt auf `queued` und
// `running`, `interrupted` und `completed` auf `running`; `queued` →
// `completed`/`interrupted` ist unzulässig und lässt die Zeile `queued`.
// `error_message` trägt bei `failed` die Klasse vor dem Text, sonst NULL;
// `started_at` bleibt NULL, wenn der Run nie `running` war. Rot färbende
// Mutation: in `UpdateBackfillRunFinish` die Klausel `status = 'queued' AND
// $2 = 'failed'` durch `status = 'queued'` ersetzen — `queued` →
// `completed` gelänge.
func TestBackfillRunFinishAllowsExactlyTheSpecifiedTransitions(t *testing.T) {
	f := newBackfillFixture(t)
	admission := newBackfillAdmission(t, f)
	runs := newBackfillRunAdapter(t, f)
	base := time.Date(2026, 9, 24, 13, 0, 0, 0, time.UTC)
	at := model.NewTimePoint(base.Add(time.Minute).UnixNano())

	newRun := func(id string, running bool) model.BackfillRun {
		run := f.admit(admission, id, id, base, model.UnknownRowEstimate())
		if running {
			run = mustRunning(t, runs, run, base)
		}
		return run
	}

	// queued -> failed
	queuedRun := newRun("fin-q-failed", false)
	failed, err := queuedRun.Fail(at, model.ErrorClassConfiguration, "Bindung fehlt")
	if err != nil {
		t.Fatalf("Fail: %v", err)
	}
	if err := runs.Finish(context.Background(), failed); err != nil {
		t.Fatalf("Finish queued->failed: %v", err)
	}
	var status string
	var started, finished *time.Time
	var message *string
	f.scan("SELECT status, started_at, finished_at, error_message FROM cdc.backfill_run WHERE run_id = 'fin-q-failed'", nil,
		&status, &started, &finished, &message)
	if status != "failed" || started != nil || finished == nil || message == nil || *message != "configuration: Bindung fehlt" {
		t.Fatalf("queued->failed: %s started %v finished %v message %v", status, started, finished, message)
	}

	// running -> failed, interrupted, completed
	runningRun := newRun("fin-r-failed", true)
	failed, _ = runningRun.Fail(at, model.ErrorClassStorage, "Verbindung verloren")
	if err := runs.Finish(context.Background(), failed); err != nil {
		t.Fatalf("Finish running->failed: %v", err)
	}
	f.scan("SELECT status, started_at, error_message FROM cdc.backfill_run WHERE run_id = 'fin-r-failed'", nil, &status, &started, &message)
	if status != "failed" || started == nil || message == nil || *message != "storage: Verbindung verloren" {
		t.Fatalf("running->failed: %s started %v message %v", status, started, message)
	}

	interruptedRun := newRun("fin-r-interrupted", true)
	interrupted, _ := interruptedRun.Interrupt(at)
	if err := runs.Finish(context.Background(), interrupted); err != nil {
		t.Fatalf("Finish running->interrupted: %v", err)
	}
	f.scan("SELECT status, error_message FROM cdc.backfill_run WHERE run_id = 'fin-r-interrupted'", nil, &status, &message)
	if status != "interrupted" || message != nil {
		t.Fatalf("running->interrupted: %s message %v", status, message)
	}

	completingRun := newRun("fin-r-completed", true)
	completed, _ := completingRun.Complete(at, 0)
	if err := runs.Finish(context.Background(), completed); err != nil {
		t.Fatalf("Finish running->completed: %v", err)
	}
	if got, _ := f.runStatus("fin-r-completed"); got != "completed" {
		t.Fatalf("running->completed: %s", got)
	}

	// queued -> completed / interrupted sind unzulässig.
	for _, target := range []model.BackfillRunStatus{model.BackfillRunCompleted, model.BackfillRunInterrupted} {
		id := "fin-q-" + string(target)
		queued := newRun(id, false)
		attempt := queued
		attempt.Status = target
		attempt.FinishedAt = at
		if err := runs.Finish(context.Background(), attempt); !stderrors.Is(err, domainerrors.ErrInvalidBackfillTransition) {
			t.Fatalf("Finish queued->%s: Fehler = %v, erwartet ErrInvalidBackfillTransition", target, err)
		}
		if got, _ := f.runStatus(id); got != "queued" {
			t.Fatalf("queued->%s: Status = %s, erwartet queued", target, got)
		}
	}

	// Eine fehlende Run-Zeile ist ein Speicherfehler.
	missing := model.BackfillRun{ID: "fin-missing", Source: backfillTestSource, Status: model.BackfillRunFailed, FinishedAt: at}
	if err := runs.Finish(context.Background(), missing); !stderrors.Is(err, outbound.ErrBackfillStorage) {
		t.Fatalf("Finish ohne Zeile: %v", err)
	}
}

// Jeder Endzustand ist endgültig: `Finish` auf einen beendeten Run ist ein
// wirkungsloser Erfolg — der Aufruf meldet nil, Status, `finished_at`,
// Zähler und Fehlertext bleiben. Rot färbende Mutation: in
// `UpdateBackfillRunFinish` die `WHERE`-Klausel auf `status IN ('queued',
// 'running', 'completed')` erweitern — der zweite Aufruf überschriebe den
// `completed`-Run mit `failed`.
func TestBackfillRunFinishOnAnEndedRunIsANoOpSuccess(t *testing.T) {
	f := newBackfillFixture(t)
	admission := newBackfillAdmission(t, f)
	runs := newBackfillRunAdapter(t, f)
	base := time.Date(2026, 9, 24, 14, 0, 0, 0, time.UTC)
	at := model.NewTimePoint(base.Add(time.Minute).UnixNano())
	later := model.NewTimePoint(base.Add(time.Hour).UnixNano())

	for _, ended := range []model.BackfillRunStatus{model.BackfillRunCompleted, model.BackfillRunFailed, model.BackfillRunInterrupted} {
		id := "end-" + string(ended)
		run := mustRunning(t, runs, f.admit(admission, id, id, base, model.UnknownRowEstimate()), base)
		var final model.BackfillRun
		var err error
		switch ended {
		case model.BackfillRunCompleted:
			final, err = run.Complete(at, 0)
		case model.BackfillRunFailed:
			final, err = run.Fail(at, model.ErrorClassInternal, "erster Ausgang")
		default:
			final, err = run.Interrupt(at)
		}
		if err != nil {
			t.Fatalf("%s: %v", ended, err)
		}
		if err := runs.Finish(context.Background(), final); err != nil {
			t.Fatalf("%s: erster Finish: %v", ended, err)
		}

		for _, second := range []model.BackfillRunStatus{model.BackfillRunFailed, model.BackfillRunCompleted, model.BackfillRunInterrupted} {
			attempt := run
			attempt.Status = second
			attempt.FinishedAt = later
			attempt.RowsCopied = 99
			attempt.ErrorMessage = "zweiter Ausgang"
			if err := runs.Finish(context.Background(), attempt); err != nil {
				t.Fatalf("%s: Finish(%s) auf beendeten Run: %v, erwartet nil", ended, second, err)
			}
		}
		var status string
		var finished time.Time
		var rowsCopied int64
		var message *string
		f.scan("SELECT status, finished_at, rows_copied, error_message FROM cdc.backfill_run WHERE run_id = $1", []any{id},
			&status, &finished, &rowsCopied, &message)
		if status != string(ended) || !finished.Equal(base.Add(time.Minute)) || rowsCopied != 0 ||
			(ended == model.BackfillRunFailed) != (message != nil && *message == "internal: erster Ausgang") {
			t.Fatalf("%s: die Zeile änderte sich: %s finished %v rows %d message %v", ended, status, finished, rowsCopied, message)
		}
	}
}

// `MarkRunning` wirkt nur auf `queued`, `RecordProgress` nur auf `running`:
// jede andere Ausgangslage endet als `ErrInvalidBackfillTransition` und
// ändert die Zeile nicht.
func TestBackfillRunMarkRunningAndProgressRequireTheirSourceState(t *testing.T) {
	f := newBackfillFixture(t)
	admission := newBackfillAdmission(t, f)
	runs := newBackfillRunAdapter(t, f)
	base := time.Date(2026, 9, 24, 15, 0, 0, 0, time.UTC)
	at := model.NewTimePoint(base.UnixNano())

	run := f.admit(admission, "src-state", "src_state", base, model.UnknownRowEstimate())
	running := mustRunning(t, runs, run, base)

	// MarkRunning auf einen laufenden Run.
	if err := runs.MarkRunning(context.Background(), running); !stderrors.Is(err, domainerrors.ErrInvalidBackfillTransition) {
		t.Fatalf("MarkRunning auf running: %v", err)
	}

	// RecordProgress auf einen queued-Run.
	queued := f.admit(admission, "src-state-q", "src_state_q", base, model.UnknownRowEstimate())
	asRunning := queued
	asRunning.Status = model.BackfillRunRunning
	asRunning.SnapshotPosition = backfillPosition(t, 77)
	asRunning.RowsCopied = 5
	if err := runs.RecordProgress(context.Background(), asRunning); !stderrors.Is(err, domainerrors.ErrInvalidBackfillTransition) {
		t.Fatalf("RecordProgress auf queued: %v", err)
	}
	if n := f.count("SELECT count(*) FROM cdc.backfill_run WHERE run_id = 'src-state-q' AND status = 'queued' AND rows_copied = 0 AND snapshot_position IS NULL"); n != 1 {
		t.Fatal("der abgelehnte Fortschritt änderte die queued-Zeile")
	}

	// RecordProgress und MarkRunning auf einen beendeten Run.
	completed, _ := running.Complete(at, 0)
	if err := runs.Finish(context.Background(), completed); err != nil {
		t.Fatalf("Finish: %v", err)
	}
	progressed, _ := running.RecordProgress(backfillPosition(t, 88), 9)
	if err := runs.RecordProgress(context.Background(), progressed); !stderrors.Is(err, domainerrors.ErrInvalidBackfillTransition) {
		t.Fatalf("RecordProgress auf completed: %v", err)
	}
	if err := runs.MarkRunning(context.Background(), running); !stderrors.Is(err, domainerrors.ErrInvalidBackfillTransition) {
		t.Fatalf("MarkRunning auf completed: %v", err)
	}
	if n := f.count("SELECT count(*) FROM cdc.backfill_run WHERE run_id = 'src-state' AND status = 'completed' AND rows_copied = 0"); n != 1 {
		t.Fatal("ein beendeter Run wurde verändert")
	}
}

// `warn_duration` geht nur von `false` nach `true`: ein späterer Fortschritt
// mit `false` löscht die Warnung nicht (`SPEC-029`). Rot färbende Mutation:
// in `UpdateBackfillRunProgress` `warn_duration = (warn_duration OR $4)`
// durch `warn_duration = $4` ersetzen.
func TestBackfillRunProgressKeepsAWarningOnceSet(t *testing.T) {
	f := newBackfillFixture(t)
	admission := newBackfillAdmission(t, f)
	runs := newBackfillRunAdapter(t, f)
	base := time.Date(2026, 9, 24, 16, 0, 0, 0, time.UTC)

	running := mustRunning(t, runs, f.admit(admission, "warn-run", "warn_run", base, model.UnknownRowEstimate()), base)
	warned, _ := running.RecordProgress(backfillPosition(t, 10), 1)
	warned.WarnDuration = true
	if err := runs.RecordProgress(context.Background(), warned); err != nil {
		t.Fatalf("RecordProgress mit Warnung: %v", err)
	}
	quiet, _ := warned.RecordProgress(backfillPosition(t, 10), 2)
	quiet.WarnDuration = false
	if err := runs.RecordProgress(context.Background(), quiet); err != nil {
		t.Fatalf("RecordProgress ohne Warnung: %v", err)
	}
	var warn bool
	var rowsCopied int64
	f.scan("SELECT warn_duration, rows_copied FROM cdc.backfill_run WHERE run_id = 'warn-run'", nil, &warn, &rowsCopied)
	if !warn || rowsCopied != 2 {
		t.Fatalf("warn_duration = %v, Zähler = %d — erwartet true/2", warn, rowsCopied)
	}
}

// Der Abgleich beim Prozessstart setzt jeden `running`-Run der Quelle auf
// `interrupted` und meldet ihre Zahl; `queued`, beendete Runs und Runs einer
// anderen Quelle bleiben. Rot färbende Mutation: in
// `UpdateBackfillRunInterrupted` `status = 'running'` durch
// `status IN ('queued', 'running')` ersetzen — der `queued`-Run ginge mit.
func TestBackfillRunInterruptRunningTouchesOnlyRunningRunsOfTheSource(t *testing.T) {
	f := newBackfillFixture(t)
	admission := newBackfillAdmission(t, f)
	runs := newBackfillRunAdapter(t, f)
	base := time.Date(2026, 9, 24, 17, 0, 0, 0, time.UTC)

	f.exec("INSERT INTO cdc.source (source_id, name) VALUES ('src-backfill-intr', 'andere') ON CONFLICT (source_id) DO NOTHING")
	f.t.Cleanup(func() {
		_, _ = f.pool.Exec(context.Background(), "DELETE FROM cdc.backfill_run WHERE run_id = 'int-foreign'")
		_, _ = f.pool.Exec(context.Background(), "DELETE FROM cdc.source WHERE source_id = 'src-backfill-intr'")
	})
	f.exec(`INSERT INTO cdc.backfill_run (run_id, source_id, schema_name, table_name, status, requested_at, started_at)
	        VALUES ('int-foreign', 'src-backfill-intr', 'public', 'int_foreign', 'running', $1, $1)`, base)

	// Fremde, in derselben Quelle laufende Runs anderer Tests bestehen nicht:
	// jeder Test räumt seine Zeilen ab und beendet seine Runs.
	f.admit(admission, "int-queued", "int_queued", base, model.UnknownRowEstimate())
	mustRunning(t, runs, f.admit(admission, "int-running-1", "int_running_1", base, model.UnknownRowEstimate()), base)
	mustRunning(t, runs, f.admit(admission, "int-running-2", "int_running_2", base, model.UnknownRowEstimate()), base)
	done := mustRunning(t, runs, f.admit(admission, "int-done", "int_done", base, model.UnknownRowEstimate()), base)
	completed, _ := done.Complete(model.NewTimePoint(base.UnixNano()), 0)
	if err := runs.Finish(context.Background(), completed); err != nil {
		t.Fatalf("Finish: %v", err)
	}

	at := model.NewTimePoint(base.Add(time.Hour).UnixNano())
	count, err := runs.InterruptRunning(context.Background(), backfillTestSource, at)
	if err != nil {
		t.Fatalf("InterruptRunning: %v", err)
	}
	if count != 2 {
		t.Fatalf("InterruptRunning meldet %d Runs, erwartet 2", count)
	}
	for id, want := range map[string]string{
		"int-running-1": "interrupted", "int-running-2": "interrupted",
		"int-queued": "queued", "int-done": "completed", "int-foreign": "running",
	} {
		if got, _ := f.runStatus(id); got != want {
			t.Fatalf("Run %s: Status %s, erwartet %s", id, got, want)
		}
	}
	var finished time.Time
	var message *string
	f.scan("SELECT finished_at, error_message FROM cdc.backfill_run WHERE run_id = 'int-running-1'", nil, &finished, &message)
	if !finished.Equal(base.Add(time.Hour)) || message != nil {
		t.Fatalf("interrupted: finished_at %v, Text %v", finished, message)
	}
}
