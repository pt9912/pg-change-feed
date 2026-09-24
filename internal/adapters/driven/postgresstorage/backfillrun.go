package postgresstorage

import (
	"context"
	stderrors "errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage/mapper"
	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage/queries"
	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage/sqlexec"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// errBackfillRunAbsent: eine Operation trifft eine Run-Zeile, die nicht
// besteht — ein inkonsistenter Zustand, kein Übergangsfehler.
var errBackfillRunAbsent = stderrors.New("Run-Zeile besteht nicht")

// BackfillRunAdapter implementiert den `BackfillRunPort` (`outbound`,
// `ARC-004`, `SPEC-029`) über den Verbindungspool des Backfill-Workers
// (`cdc_capture`-Rolle, `ADR-0047`, `ADR-0113` Festlegung 1). Er legt keine
// Run-Zeile an: die Anlage trägt allein `BackfillAdmissionAdapter`. Jede
// schreibende Anweisung trägt den erlaubten Ausgangszustand in ihrer
// `WHERE`-Klausel, ein beendeter Run bleibt damit unverändert.
type BackfillRunAdapter struct {
	db      sqlexec.DB
	log     outbound.LogPort
	timeout time.Duration
}

// NewBackfillRun baut den Verbindungspool gegen die Instanz und meldet eine
// nicht erreichbare Instanz als Fehler der Klasse `storage`
// (`outbound.ErrBackfillStorage`, `SPEC-008`).
func NewBackfillRun(ctx context.Context, dsn string, opts ...Option) (*BackfillRunAdapter, error) {
	o := newOptions(opts)
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, backfillStorageFailure(ctx, o.log, err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, backfillStorageFailure(ctx, o.log, err)
	}
	o.log.Info(ctx, "backfillrun: verbunden")
	return newBackfillRun(pool, o.log, backfillStateTimeout), nil
}

// newBackfillRun baut den Adapter über eine Ausführungs-Naht; die Tests
// setzen hier eine kurze Frist.
func newBackfillRun(db sqlexec.DB, log outbound.LogPort, timeout time.Duration) *BackfillRunAdapter {
	return &BackfillRunAdapter{db: db, log: log, timeout: timeout}
}

// Close schließt den Verbindungspool.
func (a *BackfillRunAdapter) Close() {
	a.db.Close()
}

var _ outbound.BackfillRunPort = (*BackfillRunAdapter)(nil)

// Queued liest die `queued`-Runs der Quelle in der Aufnahme-Ordnung
// `(requested_at, run_id)`.
func (a *BackfillRunAdapter) Queued(ctx context.Context, source model.SourceID) ([]model.BackfillRun, error) {
	ctx, cancel := boundedContext(ctx, a.timeout)
	defer cancel()
	return sqlexec.ReadBackfillRuns(ctx, a.db, sqlexec.Statement{
		SQL:  queries.SelectQueuedBackfillRuns,
		Args: []any{string(source)},
		Fail: func(cause error) error { return backfillStorageFailure(ctx, a.log, cause) },
	})
}

// MarkRunning setzt den Run auf `running` und hält `started_at` fest; ein Run,
// der nicht `queued` ist, endet als `ErrInvalidBackfillTransition`.
func (a *BackfillRunAdapter) MarkRunning(ctx context.Context, run model.BackfillRun) error {
	if run.Status != model.BackfillRunRunning {
		return domainerrors.ErrInvalidBackfillTransition
	}
	return a.transition(ctx, run, queries.UpdateBackfillRunRunning, string(run.ID), mapper.TimeArgument(run.StartedAt))
}

// RecordProgress schreibt Snapshot-Position, Fortschrittszähler und die
// Warnung „Kopierdauer“ des Runs fort; sie steht außerhalb der
// Daten-Transaktion des Schreibers und ist damit für Leser sichtbar. Ein Run,
// der nicht `running` ist, endet als `ErrInvalidBackfillTransition`.
func (a *BackfillRunAdapter) RecordProgress(ctx context.Context, run model.BackfillRun) error {
	if run.Status != model.BackfillRunRunning {
		return domainerrors.ErrInvalidBackfillTransition
	}
	position, err := mapper.PositionArgument(run.SnapshotPosition)
	if err != nil {
		return err
	}
	return a.transition(ctx, run, queries.UpdateBackfillRunProgress, string(run.ID), position, run.RowsCopied, run.WarnDuration)
}

// transition führt eine Übergangs-Anweisung aus, die genau eine Zeile treffen
// muss: keine betroffene Zeile heißt, dass der Run nicht im erlaubten
// Ausgangszustand steht.
func (a *BackfillRunAdapter) transition(ctx context.Context, run model.BackfillRun, sql string, args ...any) error {
	ctx, cancel := boundedContext(ctx, a.timeout)
	defer cancel()
	tag, err := a.db.Exec(ctx, sql, args...)
	if err != nil {
		return backfillStorageFailure(ctx, a.log, err)
	}
	if tag.RowsAffected() != 1 {
		return fmt.Errorf("%w: Run %s steht nicht im Ausgangszustand des Übergangs nach %s", domainerrors.ErrInvalidBackfillTransition, run.ID, run.Status)
	}
	return nil
}

// Finish hält den Endzustand `failed`, `interrupted` oder `completed` (leere
// Tabelle) fest: `failed` folgt auf `queued` und `running`, `interrupted` und
// `completed` auf `running`. Ein bereits beendeter Run bleibt unverändert und
// die Rückkehr ist ein Erfolg. Die Operation ist zeitbegrenzt und beendet sich
// bei einem vom Abbruch gelösten Kontext nicht vorzeitig.
func (a *BackfillRunAdapter) Finish(ctx context.Context, run model.BackfillRun) error {
	switch run.Status {
	case model.BackfillRunFailed, model.BackfillRunInterrupted, model.BackfillRunCompleted:
	default:
		return domainerrors.ErrInvalidBackfillTransition
	}
	position, err := mapper.PositionArgument(run.SnapshotPosition)
	if err != nil {
		return err
	}

	ctx, cancel := boundedContext(ctx, a.timeout)
	defer cancel()
	tag, err := a.db.Exec(ctx, queries.UpdateBackfillRunFinish,
		string(run.ID),
		string(run.Status),
		mapper.TimeArgument(run.FinishedAt),
		position,
		run.RowsCopied,
		run.ErrorMessage,
		run.WarnDuration,
	)
	if err != nil {
		return backfillStorageFailure(ctx, a.log, err)
	}
	if tag.RowsAffected() == 1 {
		return nil
	}

	var current string
	if err := a.db.QueryRow(ctx, queries.SelectBackfillRunStatus, string(run.ID)).Scan(&current); err != nil {
		if sqlexec.IsAbsent(err) {
			return backfillStorageFailure(ctx, a.log, fmt.Errorf("%w: %s", errBackfillRunAbsent, run.ID))
		}
		return backfillStorageFailure(ctx, a.log, err)
	}
	switch model.BackfillRunStatus(current) {
	case model.BackfillRunCompleted, model.BackfillRunFailed, model.BackfillRunInterrupted:
		return nil
	default:
		return fmt.Errorf("%w: Run %s ist %s, Endzustand %s unzulässig", domainerrors.ErrInvalidBackfillTransition, run.ID, current, run.Status)
	}
}

// InterruptRunning setzt jeden `running`-Run der Quelle auf `interrupted` und
// meldet ihre Zahl; `queued`-Runs bleiben unberührt. Die Operation ist
// zeitbegrenzt und beendet sich bei einem vom Abbruch gelösten Kontext nicht
// vorzeitig.
func (a *BackfillRunAdapter) InterruptRunning(ctx context.Context, source model.SourceID, at model.TimePoint) (int, error) {
	ctx, cancel := boundedContext(ctx, a.timeout)
	defer cancel()
	tag, err := a.db.Exec(ctx, queries.UpdateBackfillRunInterrupted, string(source), mapper.TimeArgument(at))
	if err != nil {
		return 0, backfillStorageFailure(ctx, a.log, err)
	}
	count := int(tag.RowsAffected())
	a.log.Info(ctx, "backfillrun: laufende Runs unterbrochen", "source", source, "runs", count)
	return count, nil
}
