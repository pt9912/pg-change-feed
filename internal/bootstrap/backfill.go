package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// backfillRetryInterval trägt die Wartezeit des Backfill-Workers nach einem
// Durchgang, der die Run-Zeilen nicht lesen oder einen Run nicht abschließen
// konnte: nach ihr liest er die `queued`-Zeilen erneut, ohne auf ein
// Wecksignal zu warten. Ein Durchgang ohne Fehler wartet allein auf das
// Wecksignal (`ADR-0113` Festlegung 2). Derselbe Takt wie der Fallback-Poll
// der Administrations-Goroutine.
const backfillRetryInterval = administrationPollInterval

// backfillWorkerDeps bündelt die Abhängigkeiten des Backfill-Workers
// (`runBackfillWorker`, `ADR-0111` Teilfrage 5, `ADR-0113` Festlegung 2):
// der Run-Zustands-Port liefert die `queued`-Runs der Quelle, der Use Case
// führt sie aus.
type backfillWorkerDeps struct {
	runs        outbound.BackfillRunPort
	useCase     inbound.BackfillTableUseCase
	source      model.SourceID
	publication string
	// wake ist das Wecksignal der Administrations-Goroutine
	// (`newBackfillWake`, `signalBackfillWorker`).
	wake       <-chan struct{}
	retryAfter time.Duration
	log        outbound.LogPort
}

// newBackfillWake legt das Wecksignal zwischen Administrations-Goroutine und
// Backfill-Worker an: Kapazität 1, Signale verschmelzen (`ADR-0113`
// Festlegung 2).
func newBackfillWake() chan struct{} {
	return make(chan struct{}, 1)
}

// signalBackfillWorker weckt den Backfill-Worker, ohne zu blockieren: ist
// bereits ein Signal offen, verschmilzt das neue mit ihm; der Worker liest
// nach jedem Run erneut, ein offenes Signal trägt also jeden Antrag, der
// während eines Runs angenommen wurde. Ein `nil`-Kanal (kein Worker
// verdrahtet) ist kein Fehler.
func signalBackfillWorker(wake chan<- struct{}) {
	select {
	case wake <- struct{}{}:
	default:
	}
}

// runBackfillWorker führt die angenommenen Runs der Quelle aus, bis `ctx`
// endet: ein Worker, ein Run zugleich. Jeder Durchgang arbeitet **zuerst**
// alle `queued`-Zeilen in der Reihenfolge `(requested_at, run_id)` ab und
// wartet **dann** auf das nächste Wecksignal — ein Signal, das während eines
// Runs eintrifft, geht nicht verloren, weil der Durchgang nach jedem Run
// erneut liest; die erste Lesung geschieht beim Start (`ADR-0113`
// Festlegung 2). Nach einem Durchgang mit Fehler liest der Worker nach
// `retryAfter` erneut. Eine `interrupted`-Zeile nimmt der Worker nicht auf:
// nur `queued` wird gelesen (`ADR-0111` Teilfrage 4).
func runBackfillWorker(ctx context.Context, deps backfillWorkerDeps) {
	for {
		if ctx.Err() != nil {
			return
		}
		failed := drainBackfillQueue(ctx, deps)
		if ctx.Err() != nil {
			return
		}
		var retry <-chan time.Time
		var timer *time.Timer
		if failed {
			timer = time.NewTimer(deps.retryAfter)
			retry = timer.C
		}
		select {
		case <-ctx.Done():
		case <-deps.wake:
		case <-retry:
		}
		if timer != nil {
			timer.Stop()
		}
	}
}

// drainBackfillQueue führt die `queued`-Runs nacheinander aus, bis keiner
// mehr wartet, und meldet, ob der Durchgang an einem Fehler endete. Der Worker
// liest nach jedem Run erneut. Ein Run, der nach seinem Versuch weiter
// `queued` ist, beendet den Durchgang (kein Schleifen auf derselben Zeile).
//
// Das Ergebnis von `Execute` trägt den Run im Endzustand des Aufrufs: der
// Worker verwendet es für Log und Fehlerklasse, schreibt daraus keinen
// weiteren Endzustand und stellt keinen neuen Antrag — die Zeile in
// `cdc.backfill_run` ist der dauerhafte Zustand (`cdc.backfill_status` und
// `diagnose` lesen sie). Ein Fehler des Aufrufs meldet „Endzustand nicht
// festgehalten“ oder „Kontext endete vor dem Beginn“.
func drainBackfillQueue(ctx context.Context, deps backfillWorkerDeps) (failed bool) {
	attempted := map[model.BackfillRunID]bool{}
	for ctx.Err() == nil {
		queued, err := deps.runs.Queued(ctx, deps.source)
		if err != nil {
			deps.log.Warn(ctx, "backfill: queued-Runs lesen fehlgeschlagen", "source", string(deps.source), "error", err)
			return true
		}
		if len(queued) == 0 {
			return false
		}
		run := queued[0]
		if attempted[run.ID] {
			deps.log.Warn(ctx, "backfill: Run nach seinem Versuch weiter queued — Durchgang endet", "run_id", run.ID, "table", run.QualifiedName())
			return true
		}
		attempted[run.ID] = true

		result, err := deps.useCase.Execute(ctx, inbound.BackfillExecuteCommand{Run: run, Publication: deps.publication})
		if err != nil {
			deps.log.Warn(ctx, "backfill: Run nicht abgeschlossen festgehalten", "run_id", run.ID, "table", run.QualifiedName(), "error", err)
			return true
		}
		logBackfillResult(ctx, deps.log, result.Run)
	}
	return false
}

// logBackfillResult protokolliert den Endzustand des Aufrufs; ein `failed`-
// oder `interrupted`-Run ist Berichtsinhalt, kein Fehler des Workers.
func logBackfillResult(ctx context.Context, log outbound.LogPort, run model.BackfillRun) {
	switch run.Status {
	case model.BackfillRunCompleted:
		log.Info(ctx, "backfill: Run abgeschlossen", "run_id", run.ID, "table", run.QualifiedName(), "rows_copied", run.RowsCopied)
	case model.BackfillRunFailed:
		log.Warn(ctx, "backfill: Run fehlgeschlagen", "run_id", run.ID, "table", run.QualifiedName(), "error", run.ErrorMessage)
	default:
		log.Warn(ctx, "backfill: Run nicht abgeschlossen", "run_id", run.ID, "table", run.QualifiedName(), "status", run.Status)
	}
}

// reconcileBackfillRuns setzt beim Prozessstart jeden `running`-Run der
// Quelle auf `interrupted` (`ADR-0111` Teilfrage 4, `ADR-0113` Festlegung 2):
// eine Instanz je Quelle trägt den Slot, ein `running`-Run gehört also zu
// einem beendeten Prozess. `queued`-Runs bleiben unberührt; die Aufnahme
// übernimmt der Worker.
func reconcileBackfillRuns(ctx context.Context, runs outbound.BackfillRunPort, clock outbound.ClockPort, source model.SourceID, log outbound.LogPort) error {
	interrupted, err := runs.InterruptRunning(ctx, source, clock.Now())
	if err != nil {
		return err
	}
	if interrupted > 0 {
		log.Warn(ctx, "backfill: laufende Runs beim Start auf interrupted gesetzt", "source", string(source), "runs", interrupted)
	}
	return nil
}

// diagnoseBackfillStatus gibt den zuletzt beantragten Backfill-Run je Tabelle
// der Quelle aus `cdc.backfill_status` aus (`LH-FA-SST-003`, `SPEC-029`):
// Status und Fortschritt, die **geschätzte** Zeilenzahl (unbekannt bleibt
// „unbekannt“, nie `0`), die zwei Warn-Kennzeichnungen und bei einem
// `failed`-Run der Fehlertext. Ein `failed`- oder `interrupted`-Run ist
// Berichtsinhalt, kein Befehlsfehler; der Rückgabewert meldet nur, dass die
// View lesbar war.
func diagnoseBackfillStatus(ctx context.Context, pool *pgxpool.Pool, source model.SourceID) error {
	rows, err := pool.Query(ctx,
		`SELECT schema_name, table_name, status, rows_copied, estimated_rows, warn_estimated_size, warn_duration, COALESCE(error_message, '')
		 FROM cdc.backfill_status WHERE source_id = $1 ORDER BY schema_name, table_name`, string(source))
	if err != nil {
		return err
	}
	defer rows.Close()

	fmt.Println("  Backfill je Tabelle (LH-FA-CAP-009, letzter Run; die Zeilenzahl ist geschätzt):")
	found := false
	for rows.Next() {
		var (
			schema, table, status, message string
			copied                         int64
			estimated                      *int64
			warnSize, warnDuration         bool
		)
		if err := rows.Scan(&schema, &table, &status, &copied, &estimated, &warnSize, &warnDuration, &message); err != nil {
			return err
		}
		fmt.Print(formatBackfillRun(schema, table, status, copied, estimated, warnSize, warnDuration, message))
		found = true
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if !found {
		fmt.Println("    (keiner — kein Backfill beantragt)")
	}
	return nil
}

// errBackfillNotWired meldet einen `backfill`-Antrag an eine Verdrahtung ohne
// Backfill-Use-Case; die Produktion verdrahtet ihn immer.
var errBackfillNotWired = errors.New("Backfill-Antrag ohne verdrahteten BackfillTableUseCase")

// formatBackfillRun bildet die Ausgabe eines Runs in `diagnose`: eine Zeile mit
// Status, Fortschritt, der **geschätzten** Zeilenzahl — eine unbekannte
// Schätzung (`nil`) heißt „unbekannt“, nie `0` — und den beiden
// Kennzeichnungen, bei einem Fehlertext eine zweite Zeile.
func formatBackfillRun(schema, table, status string, copied int64, estimated *int64, warnSize, warnDuration bool, message string) string {
	estimate := "unbekannt"
	if estimated != nil {
		estimate = strconv.FormatInt(*estimated, 10)
	}
	line := fmt.Sprintf("    %s.%s: %s, %d Zeilen kopiert, geschätzt %s, Warnung Größe %t, Warnung Dauer %t\n",
		schema, table, status, copied, estimate, warnSize, warnDuration)
	if message != "" {
		line += fmt.Sprintf("      Fehler: %s\n", message)
	}
	return line
}
