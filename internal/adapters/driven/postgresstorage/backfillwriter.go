package postgresstorage

import (
	"context"
	stderrors "errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage/mapper"
	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage/queries"
	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage/sqlexec"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// BackfillWriterAdapter implementiert den `BackfillWriterPort` (`outbound`,
// `ARC-004`, `ADR-0111` Teilfrage 4) über den Verbindungspool des
// Backfill-Workers (`cdc_capture`-Rolle, `ADR-0047`): **eine**
// Store-Transaktion über alle Blöcke eines Runs und ein Commit am Ende,
// zusammen mit der Run-Zeile `completed`. Die Transaktion berührt die
// Run-Zeile erst mit der letzten Anweisung vor dem Commit — der Fortschritt
// des Runs (`BackfillRunAdapter.RecordProgress`, andere Verbindung) wartet
// damit auf keine Zeilensperre der offenen Transaktion.
type BackfillWriterAdapter struct {
	db           sqlexec.DB
	log          outbound.LogPort
	stateTimeout time.Duration
	blockTimeout time.Duration
}

// NewBackfillWriter baut den Verbindungspool gegen die Instanz und meldet
// eine nicht erreichbare Instanz als Fehler der Klasse `storage`
// (`outbound.ErrBackfillStorage`, `SPEC-008`).
func NewBackfillWriter(ctx context.Context, dsn string, opts ...Option) (*BackfillWriterAdapter, error) {
	o := newOptions(opts)
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, backfillStorageFailure(ctx, o.log, err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, backfillStorageFailure(ctx, o.log, err)
	}
	o.log.Info(ctx, "backfillwriter: verbunden")
	return newBackfillWriter(pool, o.log, backfillStateTimeout, backfillBlockTimeout), nil
}

// newBackfillWriter baut den Adapter über eine Ausführungs-Naht; die Tests
// setzen hier kurze Fristen.
func newBackfillWriter(db sqlexec.DB, log outbound.LogPort, stateTimeout, blockTimeout time.Duration) *BackfillWriterAdapter {
	return &BackfillWriterAdapter{db: db, log: log, stateTimeout: stateTimeout, blockTimeout: blockTimeout}
}

// Close schließt den Verbindungspool.
func (a *BackfillWriterAdapter) Close() {
	a.db.Close()
}

var _ outbound.BackfillWriterPort = (*BackfillWriterAdapter)(nil)

// Begin öffnet die Schreibtransaktion des Runs. Der Run muss `running` sein
// und die Position `X` tragen; die Blöcke tragen sie als Commit-Position.
func (a *BackfillWriterAdapter) Begin(ctx context.Context, run model.BackfillRun) (outbound.BackfillTransaction, error) {
	if run.ID == "" || run.Status != model.BackfillRunRunning || run.SnapshotPosition.Offset == 0 {
		return nil, outbound.ErrBackfillRunInvalid
	}
	ctx, cancel := boundedContext(ctx, a.stateTimeout)
	defer cancel()
	tx, err := a.db.Begin(ctx)
	if err != nil {
		return nil, backfillStorageFailure(ctx, a.log, err)
	}
	return &backfillTransaction{adapter: a, tx: tx, run: run}, nil
}

// backfillTransaction ist die eine offene Schreibtransaktion eines Runs; sie
// wird von einer Goroutine benutzt. `blocks` zählt die angehängten Blöcke,
// `written` ihre Changes.
type backfillTransaction struct {
	adapter *BackfillWriterAdapter
	tx      pgx.Tx
	run     model.BackfillRun
	blocks  int
	written int64
	closed  bool
}

var _ outbound.BackfillTransaction = (*backfillTransaction)(nil)

// AppendBlock hängt einen Block an. Die Prüfung des Vertrags läuft vor dem
// ersten Schreibvorgang: eine verletzte Zusage schreibt nichts.
func (t *backfillTransaction) AppendBlock(ctx context.Context, block *model.ChangeTransaction) error {
	if t.closed {
		return t.adapter.failure(ctx, pgx.ErrTxClosed)
	}
	transactionRow, changeRows, err := t.blockRows(block)
	if err != nil {
		return err
	}

	ctx, cancel := boundedContext(ctx, t.adapter.blockTimeout)
	defer cancel()
	if _, err := t.tx.Exec(ctx, queries.InsertBackfillTransaction,
		transactionRow.TransactionID,
		transactionRow.SourceID,
		transactionRow.CommitPosition,
		transactionRow.CommittedAt,
	); err != nil {
		return t.adapter.failure(ctx, err)
	}
	for _, row := range changeRows {
		if _, err := t.tx.Exec(ctx, queries.InsertBackfillChange,
			row.ChangeID,
			row.TransactionID,
			row.SourceTableID,
			row.Sequence,
			row.Operation,
			mapper.JSONImage(row.OldData),
			mapper.JSONImage(row.NewData),
			row.SchemaVersion,
			row.Origin,
		); err != nil {
			return t.adapter.failure(ctx, err)
		}
	}
	t.blocks++
	t.written += int64(len(changeRows))
	return nil
}

// blockRows prüft den Block gegen den Vertrag des Schreibers und bildet
// seine Zeilen: committed, auf der Position `X` des Runs, die Kennung
// `model.BackfillTransactionID` in der Reihenfolge 1, 2, 3, …, mindestens ein
// Change, jeder Change `model.ChangeIDFor`, `INSERT`, Herkunft `backfill`,
// ohne `old_data`.
func (t *backfillTransaction) blockRows(block *model.ChangeTransaction) (mapper.TransactionRow, []mapper.ChangeRow, error) {
	invalid := func(reason string) (mapper.TransactionRow, []mapper.ChangeRow, error) {
		return mapper.TransactionRow{}, nil, fmt.Errorf("%w: %s", outbound.ErrBackfillBlockInvalid, reason)
	}
	if block == nil {
		return invalid("kein Block")
	}
	position, committed := block.CommitPosition()
	if !committed {
		return mapper.TransactionRow{}, nil, domainerrors.ErrTransactionNotCommitted
	}
	if block.SourceID != t.run.Source || position != t.run.SnapshotPosition {
		return invalid("Quelle oder Position weichen von der Position X des Runs ab")
	}
	expected, err := model.BackfillTransactionID(t.run.ID, t.blocks+1)
	if err != nil {
		return mapper.TransactionRow{}, nil, err
	}
	if block.ID != expected {
		return invalid(fmt.Sprintf("Transaktions-Kennung %s, erwartet %s", block.ID, expected))
	}
	changes, err := block.Changes()
	if err != nil {
		return mapper.TransactionRow{}, nil, err
	}
	if len(changes) == 0 {
		return invalid("Block ohne Change")
	}
	for _, change := range changes {
		switch {
		case change.Origin != model.ChangeOriginBackfill:
			return invalid(fmt.Sprintf("Change %s trägt die Herkunft %q", change.ID, change.Origin))
		case change.Operation != model.OperationInsert:
			return invalid(fmt.Sprintf("Change %s trägt die Operation %s", change.ID, change.Operation))
		case len(change.OldImage) != 0:
			return invalid(fmt.Sprintf("Change %s trägt old_data", change.ID))
		case change.TransactionID != block.ID || change.ID != model.ChangeIDFor(block.ID, change.Sequence):
			return invalid(fmt.Sprintf("Change %s trägt nicht die Kennung des Blocks", change.ID))
		}
	}
	transactionRow, err := mapper.NewTransactionRow(block, position)
	if err != nil {
		return mapper.TransactionRow{}, nil, err
	}
	changeRows, err := mapper.NewChangeRows(changes)
	if err != nil {
		return mapper.TransactionRow{}, nil, err
	}
	return transactionRow, changeRows, nil
}

// Commit schreibt die Run-Zeile im Endzustand `completed` als letzte
// Anweisung und committet die Transaktion. Der Zähler des Runs trägt die
// Zahl der angehängten Changes; eine Abweichung, ein anderer Run oder ein Run
// außerhalb von `running` (etwa ein zwischenzeitlich `interrupted` gesetzter)
// committet nichts, die Transaktion bleibt für `Rollback` offen.
func (t *backfillTransaction) Commit(ctx context.Context, run model.BackfillRun) error {
	if t.closed {
		return t.adapter.failure(ctx, pgx.ErrTxClosed)
	}
	if run.ID != t.run.ID || run.Status != model.BackfillRunCompleted {
		return fmt.Errorf("%w: Commit trägt nicht den Run %s im Zustand completed", outbound.ErrBackfillBlockInvalid, t.run.ID)
	}
	if run.RowsCopied != t.written {
		return fmt.Errorf("%w: Zähler %d, angehängte Changes %d", outbound.ErrBackfillBlockInvalid, run.RowsCopied, t.written)
	}
	position, err := mapper.PositionArgument(run.SnapshotPosition)
	if err != nil {
		return err
	}

	ctx, cancel := boundedContext(ctx, t.adapter.blockTimeout)
	defer cancel()
	tag, err := t.tx.Exec(ctx, queries.UpdateBackfillRunCompleted,
		string(run.ID),
		mapper.TimeArgument(run.FinishedAt),
		position,
		run.RowsCopied,
		run.WarnDuration,
	)
	if err != nil {
		return t.adapter.failure(ctx, err)
	}
	if tag.RowsAffected() != 1 {
		return fmt.Errorf("%w: Run %s ist nicht running", domainerrors.ErrInvalidBackfillTransition, run.ID)
	}
	err = t.tx.Commit(ctx)
	t.closed = true
	if err != nil {
		return t.adapter.failure(ctx, err)
	}
	return nil
}

// Rollback verwirft alle angehängten Blöcke. Ein wiederholter Aufruf und ein
// Aufruf nach dem Commit bleiben ohne Wirkung. Die Operation ist
// zeitbegrenzt und beendet sich bei einem vom Abbruch gelösten Kontext nicht
// vorzeitig.
func (t *backfillTransaction) Rollback(ctx context.Context) error {
	if t.closed {
		return nil
	}
	t.closed = true
	ctx, cancel := boundedContext(ctx, t.adapter.stateTimeout)
	defer cancel()
	if err := t.tx.Rollback(ctx); err != nil && !stderrors.Is(err, pgx.ErrTxClosed) {
		return t.adapter.failure(ctx, err)
	}
	return nil
}

// failure übersetzt einen Treiber-Fehler in die Klasse `storage`.
func (a *BackfillWriterAdapter) failure(ctx context.Context, cause error) error {
	return backfillStorageFailure(ctx, a.log, cause)
}
