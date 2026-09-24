package postgresstorage

import (
	"context"
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

// BackfillAdmissionAdapter implementiert den `BackfillAdmissionPort`
// (`outbound`, `ARC-004`, `ADR-0113` Festlegung 1) über den Verbindungspool
// der Administrations-Goroutine (`cdc_admin`-Rolle, `ADR-0047`): die Annahme
// ist **eine** Transaktion über `cdc.backfill_run` und
// `cdc.administration_request`. `db` trägt die Ausführung über die schmale
// Naht (`sqlexec`, `ADR-0071` Punkt 5).
type BackfillAdmissionAdapter struct {
	db      sqlexec.DB
	log     outbound.LogPort
	timeout time.Duration
}

// NewBackfillAdmission baut den Verbindungspool gegen die Instanz und meldet
// eine nicht erreichbare Instanz als Fehler der Klasse `storage`
// (`outbound.ErrBackfillStorage`, `SPEC-008`).
func NewBackfillAdmission(ctx context.Context, dsn string, opts ...Option) (*BackfillAdmissionAdapter, error) {
	o := newOptions(opts)
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, backfillStorageFailure(ctx, o.log, err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, backfillStorageFailure(ctx, o.log, err)
	}
	o.log.Info(ctx, "backfilladmission: verbunden")
	return newBackfillAdmission(pool, o.log, backfillStateTimeout), nil
}

// newBackfillAdmission baut den Adapter über eine Ausführungs-Naht; die Tests
// setzen hier eine kurze Frist.
func newBackfillAdmission(db sqlexec.DB, log outbound.LogPort, timeout time.Duration) *BackfillAdmissionAdapter {
	return &BackfillAdmissionAdapter{db: db, log: log, timeout: timeout}
}

// Close schließt den Verbindungspool.
func (a *BackfillAdmissionAdapter) Close() {
	a.db.Close()
}

var _ outbound.BackfillAdmissionPort = (*BackfillAdmissionAdapter)(nil)

// Admit nimmt den Antrag in **einer** Transaktion an: Prüfung „kein aktiver
// Run derselben Tabelle“, Anlage der Run-Zeile `queued`, Antragsvermerk
// `applied`, Commit. Der Vermerk muss genau eine `pending`-Zeile treffen;
// jede Abweichung verwirft die Transaktion und hinterlässt weder Run-Zeile
// noch Vermerk. Zwischen mehreren Annehmenden sperrt die Transaktion
// nicht (`ADR-0113` Festlegung 1 Punkt 2): sie beruht darauf, dass eine
// Goroutine je Quelle Anträge nacheinander annimmt.
func (a *BackfillAdmissionAdapter) Admit(ctx context.Context, requestID model.AdministrationRequestID, run model.BackfillRun) error {
	if requestID == "" {
		return domainerrors.ErrEmptyIdentifier
	}
	if run.ID != model.BackfillRunID(requestID) || run.Status != model.BackfillRunQueued || run.RequestedAt.Unset() ||
		run.Source == "" || run.Schema == "" || run.Table == "" {
		return outbound.ErrBackfillRunInvalid
	}

	ctx, cancel := boundedContext(ctx, a.timeout)
	defer cancel()

	tx, err := a.db.Begin(ctx)
	if err != nil {
		return backfillStorageFailure(ctx, a.log, err)
	}
	defer func() {
		rollbackCtx, rollbackCancel := boundedContext(context.WithoutCancel(ctx), a.timeout)
		defer rollbackCancel()
		_ = tx.Rollback(rollbackCtx)
	}()

	var active bool
	if err := tx.QueryRow(ctx, queries.SelectActiveBackfillRun, string(run.Source), run.Schema, run.Table).Scan(&active); err != nil {
		return backfillStorageFailure(ctx, a.log, err)
	}
	if active {
		return fmt.Errorf("%w: %s", domainerrors.ErrBackfillRunActive, run.QualifiedName())
	}

	if _, err := tx.Exec(ctx, queries.InsertBackfillRun,
		string(run.ID),
		string(run.Source),
		run.Schema,
		run.Table,
		mapper.TimeArgument(run.RequestedAt),
		mapper.EstimateArgument(run.EstimatedRows),
		run.WarnEstimatedSize,
	); err != nil {
		return backfillStorageFailure(ctx, a.log, err)
	}

	tag, err := tx.Exec(ctx, queries.UpdateAdministrationRequestApplied, string(requestID))
	if err != nil {
		return backfillStorageFailure(ctx, a.log, err)
	}
	if tag.RowsAffected() != 1 {
		return fmt.Errorf("%w: %s", outbound.ErrBackfillRequestNotPending, requestID)
	}

	if err := tx.Commit(ctx); err != nil {
		return backfillStorageFailure(ctx, a.log, err)
	}
	a.log.Info(ctx, "backfilladmission: Antrag angenommen", "request_id", requestID, "table", run.QualifiedName())
	return nil
}
