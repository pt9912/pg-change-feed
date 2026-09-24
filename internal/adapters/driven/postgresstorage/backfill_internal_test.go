package postgresstorage

import (
	"context"
	stderrors "errors"
	"math"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage/mapper"
	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage/queries"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Diese Tests laufen netzlos an der Ausführungs-Naht (`sqlexec`,
// `ADR-0071` Punkt 5): sie prüfen die Fristen der drei Backfill-Adapter und
// den Vertrag der Adapter-Eingaben, ohne PostgreSQL. Das SQL selbst und die
// Atomarität prüft der Store-Tier (`make test-store`,
// `backfilladmission_test.go`, `backfillrun_test.go`,
// `backfillwriter_test.go`).

// errNoDeadline ist der Fehler einer Operation, die auch nach der Obergrenze
// des Fakes keine Frist gesehen hat.
var errNoDeadline = stderrors.New("fake: Operation ohne Frist blockiert")

// blockingSeam erfüllt `sqlexec.DB` und blockiert jede Operation, bis ihr
// Kontext endet; ohne Kontext-Ende gibt sie nach `hardCap` mit
// `errNoDeadline` auf, damit ein fehlendes Zeitlimit den Test rot färbt statt
// ihn zu hängen.
type blockingSeam struct {
	hardCap time.Duration
	tx      *blockingTx
}

func (s *blockingSeam) wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(s.hardCap):
		return errNoDeadline
	}
}

func (s *blockingSeam) Query(ctx context.Context, _ string, _ ...any) (pgx.Rows, error) {
	return nil, s.wait(ctx)
}

func (s *blockingSeam) QueryRow(ctx context.Context, _ string, _ ...any) pgx.Row {
	return blockingRow{err: s.wait(ctx)}
}

func (s *blockingSeam) Exec(ctx context.Context, _ string, _ ...any) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, s.wait(ctx)
}

func (s *blockingSeam) Begin(context.Context) (pgx.Tx, error) {
	if s.tx == nil {
		return nil, stderrors.New("fake: kein Transaktions-Träger")
	}
	return s.tx, nil
}

func (s *blockingSeam) Close() {}

type blockingRow struct{ err error }

func (r blockingRow) Scan(...any) error { return r.err }

// blockingTx erfüllt `pgx.Tx` für die Operationen der Schreibtransaktion und
// blockiert sie wie `blockingSeam`.
type blockingTx struct {
	pgx.Tx
	seam *blockingSeam
}

func (t *blockingTx) Exec(ctx context.Context, _ string, _ ...any) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, t.seam.wait(ctx)
}

func (t *blockingTx) QueryRow(ctx context.Context, _ string, _ ...any) pgx.Row {
	return blockingRow{err: t.seam.wait(ctx)}
}

func (t *blockingTx) Commit(ctx context.Context) error   { return t.seam.wait(ctx) }
func (t *blockingTx) Rollback(ctx context.Context) error { return t.seam.wait(ctx) }

func newBlockingSeam() *blockingSeam {
	seam := &blockingSeam{hardCap: 3 * time.Second}
	seam.tx = &blockingTx{seam: seam}
	return seam
}

// detachedFromCancelledParent liefert den Kontext, den der Use Case für
// Endzustand und Rollback benutzt: vom Abbruch gelöst und ohne Frist.
func detachedFromCancelledParent() context.Context {
	parent, cancel := context.WithCancel(context.Background())
	cancel()
	return context.WithoutCancel(parent)
}

// assertBoundedByOwnTimeout prüft: die Operation endete an der eigenen Frist
// — nicht vorzeitig (der Kontext war vom Abbruch gelöst), nicht erst an der
// Obergrenze des Fakes — mit dem Frist-Fehler in der Ursache.
func assertBoundedByOwnTimeout(t *testing.T, name string, timeout time.Duration, elapsed time.Duration, err error) {
	t.Helper()
	if !stderrors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("%s: Fehler = %v, erwartet die eigene Frist (context.DeadlineExceeded)", name, err)
	}
	if elapsed < timeout {
		t.Fatalf("%s: endete nach %v, vor der eigenen Frist %v", name, elapsed, timeout)
	}
	if elapsed >= 2*time.Second {
		t.Fatalf("%s: endete erst nach %v, die Frist %v greift nicht", name, elapsed, timeout)
	}
}

// Die Zusage: jede Operation des Run-Zustands-Adapters endet an der
// adapterseitigen Frist, auch auf einem vom Abbruch gelösten Kontext ohne
// Frist des Aufrufers. Rot färbende Mutation: in `Finish`,
// `InterruptRunning` oder `transition` den Aufruf von `boundedContext`
// entfernen — die Operation blockiert bis zur Obergrenze des Fakes und
// endet mit `errNoDeadline`.
func TestBackfillRunOperationsEndAtTheirOwnTimeoutOnDetachedContext(t *testing.T) {
	const timeout = 40 * time.Millisecond
	adapter := newBackfillRun(newBlockingSeam(), outbound.NoopLog, timeout)

	running := model.BackfillRun{ID: "run-1", Source: "src-1", Status: model.BackfillRunRunning}
	failed := model.BackfillRun{ID: "run-1", Source: "src-1", Status: model.BackfillRunFailed}
	operations := map[string]func(context.Context) error{
		"Finish": func(ctx context.Context) error { return adapter.Finish(ctx, failed) },
		"InterruptRunning": func(ctx context.Context) error {
			_, err := adapter.InterruptRunning(ctx, "src-1", model.NewTimePoint(1))
			return err
		},
		"MarkRunning":    func(ctx context.Context) error { return adapter.MarkRunning(ctx, running) },
		"RecordProgress": func(ctx context.Context) error { return adapter.RecordProgress(ctx, running) },
		"Queued": func(ctx context.Context) error {
			_, err := adapter.Queued(ctx, "src-1")
			return err
		},
	}
	for name, operation := range operations {
		start := time.Now()
		err := operation(detachedFromCancelledParent())
		assertBoundedByOwnTimeout(t, name, timeout, time.Since(start), err)
		if name != "Queued" && !stderrors.Is(err, outbound.ErrBackfillStorage) {
			t.Fatalf("%s: Fehler trägt nicht die Klasse storage: %v", name, err)
		}
	}
}

// Dieselbe Zusage für die Annahme und die Schreibtransaktion: Annahme,
// Anhängen eines Blocks, Commit und Rollback enden an der adapterseitigen
// Frist. Rot färbende Mutation: in `Rollback` den Aufruf von
// `boundedContext` entfernen.
func TestBackfillAdmissionAndWriterEndAtTheirOwnTimeoutOnDetachedContext(t *testing.T) {
	const timeout = 40 * time.Millisecond

	admission := newBackfillAdmission(newBlockingSeam(), outbound.NoopLog, timeout)
	queued := model.BackfillRun{ID: "req-1", Source: "src-1", Schema: "public", Table: "orders",
		Status: model.BackfillRunQueued, RequestedAt: model.NewTimePoint(5)}
	// Der blockierende Träger lässt schon `Begin` durch; die Annahme blockiert
	// an der ersten Anweisung ihrer Transaktion.
	start := time.Now()
	err := admission.Admit(detachedFromCancelledParent(), "req-1", queued)
	assertBoundedByOwnTimeout(t, "Admit", timeout, time.Since(start), err)

	seam := newBlockingSeam()
	writer := newBackfillWriter(seam, outbound.NoopLog, timeout, timeout)
	position := model.SourcePosition{SourceID: "src-1", Offset: 100}
	run := model.BackfillRun{ID: "run-1", Source: "src-1", Status: model.BackfillRunRunning, SnapshotPosition: position}
	tx, err := writer.Begin(detachedFromCancelledParent(), run)
	if err != nil {
		t.Fatalf("Begin: %v", err)
	}

	start = time.Now()
	err = tx.AppendBlock(detachedFromCancelledParent(), validBlock(t, run, 1, 1))
	assertBoundedByOwnTimeout(t, "AppendBlock", timeout, time.Since(start), err)

	completed := run
	completed.Status = model.BackfillRunCompleted
	start = time.Now()
	err = tx.Commit(detachedFromCancelledParent(), completed)
	assertBoundedByOwnTimeout(t, "Commit", timeout, time.Since(start), err)

	// Der Commit endete mit Fehler, die Transaktion ist geschlossen; ein
	// Rollback auf einer noch offenen Transaktion trägt seine eigene Frist.
	open, err := writer.Begin(detachedFromCancelledParent(), run)
	if err != nil {
		t.Fatalf("Begin: %v", err)
	}
	start = time.Now()
	err = open.Rollback(detachedFromCancelledParent())
	assertBoundedByOwnTimeout(t, "Rollback", timeout, time.Since(start), err)
}

// recordingTx erfüllt `pgx.Tx` und hält jede Anweisung fest; der Test prüft
// damit, was der Adapter an die Naht gibt und was er **nicht** gibt.
type recordingTx struct {
	pgx.Tx
	execs      []recordedExec
	tag        pgconn.CommandTag
	committed  bool
	rolledBack bool
}

type recordedExec struct {
	sql  string
	args []any
}

func (t *recordingTx) Exec(_ context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	t.execs = append(t.execs, recordedExec{sql: sql, args: args})
	return t.tag, nil
}

func (t *recordingTx) Commit(context.Context) error   { t.committed = true; return nil }
func (t *recordingTx) Rollback(context.Context) error { t.rolledBack = true; return nil }

// recordingSeam erfüllt `sqlexec.DB` und reicht `recordingTx` als
// Transaktion.
type recordingSeam struct {
	blockingSeam
	tx *recordingTx
}

func (s *recordingSeam) Begin(context.Context) (pgx.Tx, error) { return s.tx, nil }

func newRecordingWriter(t *testing.T, run model.BackfillRun) (outbound.BackfillTransaction, *recordingTx) {
	t.Helper()
	tx := &recordingTx{tag: pgconn.NewCommandTag("UPDATE 1")}
	writer := newBackfillWriter(&recordingSeam{tx: tx}, outbound.NoopLog, time.Second, time.Second)
	transaction, err := writer.Begin(context.Background(), run)
	if err != nil {
		t.Fatalf("Begin: %v", err)
	}
	return transaction, tx
}

func runningRun() model.BackfillRun {
	return model.BackfillRun{
		ID: "run-1", Source: "src-1", Schema: "public", Table: "orders", Status: model.BackfillRunRunning,
		SnapshotPosition: model.SourcePosition{SourceID: "src-1", Offset: 100},
	}
}

// validBlock baut den Block `number` des Runs mit `changes` Changes so, wie
// ihn der Use Case baut: committed auf der Position `X`, Kennungen aus
// `model.BackfillTransactionID`/`model.ChangeIDFor`, `INSERT`, Herkunft
// `backfill`, ohne `old_data`.
func validBlock(t *testing.T, run model.BackfillRun, number, changes int) *model.ChangeTransaction {
	t.Helper()
	return buildBlock(t, run, number, changes, func(change model.Change) model.Change { return change })
}

func buildBlock(t *testing.T, run model.BackfillRun, number, changes int, adjust func(model.Change) model.Change) *model.ChangeTransaction {
	t.Helper()
	id, err := model.BackfillTransactionID(run.ID, number)
	if err != nil {
		t.Fatalf("BackfillTransactionID: %v", err)
	}
	block, err := model.NewOpenTransaction(id, run.Source)
	if err != nil {
		t.Fatalf("NewOpenTransaction: %v", err)
	}
	for sequence := int64(1); sequence <= int64(changes); sequence++ {
		change, err := model.NewChange(model.ChangeIDFor(id, sequence), id, "tbl-1", sequence,
			model.OperationInsert, nil, []byte(`{"n":1}`), "sv-1")
		if err != nil {
			t.Fatalf("NewChange: %v", err)
		}
		if change, err = change.WithOrigin(model.ChangeOriginBackfill); err != nil {
			t.Fatalf("WithOrigin: %v", err)
		}
		if err := block.AppendChange(adjust(change)); err != nil {
			t.Fatalf("AppendChange: %v", err)
		}
	}
	if err := block.Commit(run.SnapshotPosition, model.NewTimePoint(1_000)); err != nil {
		t.Fatalf("Commit: %v", err)
	}
	return block
}

// Ein gültiger Block schreibt eine Transaktions-Zeile und je Change eine
// Change-Zeile mit den Werten des Blocks: Kennung, Position `X`, Herkunft
// `backfill`, Operation `INSERT`, kein `old_data`.
func TestBackfillWriterAppendsTheBlockAsRows(t *testing.T) {
	run := runningRun()
	transaction, tx := newRecordingWriter(t, run)

	if err := transaction.AppendBlock(context.Background(), validBlock(t, run, 1, 2)); err != nil {
		t.Fatalf("AppendBlock: %v", err)
	}

	if len(tx.execs) != 3 {
		t.Fatalf("Anweisungen = %d, erwartet 3 (Transaktion + 2 Changes)", len(tx.execs))
	}
	if tx.execs[0].sql != queries.InsertBackfillTransaction ||
		tx.execs[0].args[0] != "0bf-run-1-00000001" || tx.execs[0].args[1] != "src-1" || tx.execs[0].args[2] != int64(100) {
		t.Fatalf("Transaktions-Zeile = %+v", tx.execs[0])
	}
	change := tx.execs[2]
	if change.sql != queries.InsertBackfillChange || change.args[0] != "0bf-run-1-00000001-2" ||
		change.args[1] != "0bf-run-1-00000001" || change.args[3] != int64(2) || change.args[4] != "INSERT" ||
		change.args[5] != nil || change.args[6] != `{"n":1}` || change.args[8] != "backfill" {
		t.Fatalf("Change-Zeile = %+v", change)
	}
}

// Jede Abweichung vom Vertrag des Schreibers endet als
// `ErrBackfillBlockInvalid`, **bevor** eine Anweisung die Naht erreicht: die
// mutierte Eingabe ist je Zeile der Tabelle benannt. Rot färbende Mutation je
// Zeile: die zugehörige Prüfung in `blockRows` entfernen — der Block erreicht
// die Naht, der Test meldet die geschriebenen Anweisungen.
func TestBackfillWriterRejectsBlocksOutsideTheContract(t *testing.T) {
	run := runningRun()
	otherPosition := model.SourcePosition{SourceID: "src-1", Offset: 101}
	otherSource := model.BackfillRun{ID: "run-1", Source: "src-2", SnapshotPosition: model.SourcePosition{SourceID: "src-2", Offset: 100}}
	cases := []struct {
		name  string
		block func() *model.ChangeTransaction
	}{
		{"falsche Blocknummer (2 statt 1)", func() *model.ChangeTransaction { return validBlock(t, run, 2, 1) }},
		{"Kennung eines anderen Runs", func() *model.ChangeTransaction {
			other := run
			other.ID = "run-9"
			return validBlock(t, other, 1, 1)
		}},
		{"Herkunft wal", func() *model.ChangeTransaction {
			return buildBlock(t, run, 1, 1, func(c model.Change) model.Change { c.Origin = model.ChangeOriginWAL; return c })
		}},
		{"Operation UPDATE", func() *model.ChangeTransaction {
			return buildBlock(t, run, 1, 1, func(c model.Change) model.Change { c.Operation = model.OperationUpdate; return c })
		}},
		{"old_data gesetzt", func() *model.ChangeTransaction {
			return buildBlock(t, run, 1, 1, func(c model.Change) model.Change { c.OldImage = []byte(`{"n":0}`); return c })
		}},
		{"Change-Kennung nicht ChangeIDFor", func() *model.ChangeTransaction {
			return buildBlock(t, run, 1, 1, func(c model.Change) model.Change { c.ID = "irgendeine-kennung"; return c })
		}},
		{"Block ohne Change", func() *model.ChangeTransaction { return validBlock(t, run, 1, 0) }},
		{"andere Position als X", func() *model.ChangeTransaction {
			shifted := run
			shifted.SnapshotPosition = otherPosition
			return validBlock(t, shifted, 1, 1)
		}},
		{"andere Quelle", func() *model.ChangeTransaction { return validBlock(t, otherSource, 1, 1) }},
		{"kein Block", func() *model.ChangeTransaction { return nil }},
	}
	for _, c := range cases {
		transaction, tx := newRecordingWriter(t, run)
		err := transaction.AppendBlock(context.Background(), c.block())
		if !stderrors.Is(err, outbound.ErrBackfillBlockInvalid) {
			t.Fatalf("%s: Fehler = %v, erwartet ErrBackfillBlockInvalid", c.name, err)
		}
		if len(tx.execs) != 0 {
			t.Fatalf("%s: %d Anweisungen erreichten die Naht", c.name, len(tx.execs))
		}
	}
}

// Ein Block, der nicht committed ist, endet als Domänen-Fehler ohne
// Schreibvorgang.
func TestBackfillWriterRejectsAnOpenBlock(t *testing.T) {
	run := runningRun()
	transaction, tx := newRecordingWriter(t, run)
	id, _ := model.BackfillTransactionID(run.ID, 1)
	open, err := model.NewOpenTransaction(id, run.Source)
	if err != nil {
		t.Fatalf("NewOpenTransaction: %v", err)
	}
	if err := transaction.AppendBlock(context.Background(), open); !stderrors.Is(err, domainerrors.ErrTransactionNotCommitted) {
		t.Fatalf("Fehler = %v", err)
	}
	if len(tx.execs) != 0 {
		t.Fatalf("Anweisungen = %d", len(tx.execs))
	}
}

// Die Blöcke zählen ab 1 lückenlos: nach Block 1 verlangt der Schreiber
// Block 2, ein Wiederholen von Block 1 endet als Vertragsverletzung.
func TestBackfillWriterRequiresConsecutiveBlockNumbers(t *testing.T) {
	run := runningRun()
	transaction, tx := newRecordingWriter(t, run)
	if err := transaction.AppendBlock(context.Background(), validBlock(t, run, 1, 1)); err != nil {
		t.Fatalf("Block 1: %v", err)
	}
	if err := transaction.AppendBlock(context.Background(), validBlock(t, run, 1, 1)); !stderrors.Is(err, outbound.ErrBackfillBlockInvalid) {
		t.Fatalf("Block 1 erneut: %v", err)
	}
	if err := transaction.AppendBlock(context.Background(), validBlock(t, run, 2, 1)); err != nil {
		t.Fatalf("Block 2: %v", err)
	}
	if len(tx.execs) != 4 {
		t.Fatalf("Anweisungen = %d, erwartet 4 (zwei Blöcke à Transaktion + Change)", len(tx.execs))
	}
}

// Der Commit schreibt die Run-Zeile als letzte Anweisung und committet; der
// Zähler muss die Zahl der angehängten Changes tragen. Rot färbende Mutation:
// die Zähler-Prüfung in `Commit` entfernen.
func TestBackfillWriterCommitWritesTheRunRowLastAndChecksTheCounter(t *testing.T) {
	run := runningRun()
	transaction, tx := newRecordingWriter(t, run)
	if err := transaction.AppendBlock(context.Background(), validBlock(t, run, 1, 3)); err != nil {
		t.Fatalf("AppendBlock: %v", err)
	}

	completed := run
	completed.Status = model.BackfillRunCompleted
	completed.RowsCopied = 2
	completed.FinishedAt = model.NewTimePoint(9)
	before := len(tx.execs)
	if err := transaction.Commit(context.Background(), completed); !stderrors.Is(err, outbound.ErrBackfillBlockInvalid) {
		t.Fatalf("Commit mit falschem Zähler: %v", err)
	}
	if len(tx.execs) != before || tx.committed {
		t.Fatalf("der abgelehnte Commit erreichte die Naht (Anweisungen %d, committed %v)", len(tx.execs)-before, tx.committed)
	}

	completed.RowsCopied = 3
	if err := transaction.Commit(context.Background(), completed); err != nil {
		t.Fatalf("Commit: %v", err)
	}
	last := tx.execs[len(tx.execs)-1]
	if last.sql != queries.UpdateBackfillRunCompleted || last.args[0] != "run-1" || last.args[3] != int64(3) {
		t.Fatalf("letzte Anweisung = %+v", last)
	}
	if !tx.committed {
		t.Fatal("die Transaktion ist nicht committet")
	}
	if err := transaction.Rollback(context.Background()); err != nil || tx.rolledBack {
		t.Fatalf("Rollback nach dem Commit: %v, rolledBack %v", err, tx.rolledBack)
	}
}

// Trifft die Run-Zeile keinen `running`-Run (etwa ein zwischenzeitlich
// `interrupted` gesetzter), committet der Schreiber nicht.
func TestBackfillWriterDoesNotCommitWhenTheRunIsNotRunning(t *testing.T) {
	run := runningRun()
	transaction, tx := newRecordingWriter(t, run)
	if err := transaction.AppendBlock(context.Background(), validBlock(t, run, 1, 1)); err != nil {
		t.Fatalf("AppendBlock: %v", err)
	}
	tx.tag = pgconn.NewCommandTag("UPDATE 0")
	completed := run
	completed.Status = model.BackfillRunCompleted
	completed.RowsCopied = 1

	if err := transaction.Commit(context.Background(), completed); !stderrors.Is(err, domainerrors.ErrInvalidBackfillTransition) {
		t.Fatalf("Commit: %v", err)
	}
	if tx.committed {
		t.Fatal("die Transaktion committet trotz nicht getroffener Run-Zeile")
	}
	if err := transaction.Rollback(context.Background()); err != nil || !tx.rolledBack {
		t.Fatalf("Rollback: %v, rolledBack %v", err, tx.rolledBack)
	}
	if err := transaction.Rollback(context.Background()); err != nil {
		t.Fatalf("wiederholter Rollback: %v", err)
	}
}

// Die Eingaben der Operationen tragen ihren Vertrag: `Begin` verlangt einen
// `running`-Run mit Position, `Admit` einen `queued`-Run mit der Kennung des
// Antrags. Die Naht wird nicht erreicht (`nil`-Träger würde abstürzen).
func TestBackfillAdaptersRejectInputsOutsideTheirContract(t *testing.T) {
	writer := newBackfillWriter(nil, outbound.NoopLog, time.Second, time.Second)
	for name, run := range map[string]model.BackfillRun{
		"Begin ohne Kennung":    {Status: model.BackfillRunRunning, SnapshotPosition: model.SourcePosition{SourceID: "s", Offset: 1}},
		"Begin queued":          {ID: "r", Status: model.BackfillRunQueued, SnapshotPosition: model.SourcePosition{SourceID: "s", Offset: 1}},
		"Begin ohne Position X": {ID: "r", Status: model.BackfillRunRunning},
	} {
		if _, err := writer.Begin(context.Background(), run); !stderrors.Is(err, outbound.ErrBackfillRunInvalid) {
			t.Fatalf("%s: %v", name, err)
		}
	}

	admission := newBackfillAdmission(nil, outbound.NoopLog, time.Second)
	base := model.BackfillRun{ID: "req-1", Source: "s", Schema: "public", Table: "t", Status: model.BackfillRunQueued, RequestedAt: model.NewTimePoint(1)}
	mutate := func(edit func(*model.BackfillRun)) model.BackfillRun { r := base; edit(&r); return r }
	for name, run := range map[string]model.BackfillRun{
		"Kennung weicht vom Antrag ab": mutate(func(r *model.BackfillRun) { r.ID = "req-2" }),
		"Run running":                  mutate(func(r *model.BackfillRun) { r.Status = model.BackfillRunRunning }),
		"ohne requested_at":            mutate(func(r *model.BackfillRun) { r.RequestedAt = model.TimePoint{} }),
		"ohne Quelle":                  mutate(func(r *model.BackfillRun) { r.Source = "" }),
	} {
		if err := admission.Admit(context.Background(), "req-1", run); !stderrors.Is(err, outbound.ErrBackfillRunInvalid) {
			t.Fatalf("Admit %s: %v", name, err)
		}
	}
	if err := admission.Admit(context.Background(), "", base); !stderrors.Is(err, domainerrors.ErrEmptyIdentifier) {
		t.Fatalf("Admit ohne Antrags-Kennung: %v", err)
	}

	runs := newBackfillRun(nil, outbound.NoopLog, time.Second)
	if err := runs.Finish(context.Background(), model.BackfillRun{ID: "r", Status: model.BackfillRunQueued}); !stderrors.Is(err, domainerrors.ErrInvalidBackfillTransition) {
		t.Fatalf("Finish queued: %v", err)
	}
	if err := runs.Finish(context.Background(), model.BackfillRun{ID: "r", Status: model.BackfillRunRunning}); !stderrors.Is(err, domainerrors.ErrInvalidBackfillTransition) {
		t.Fatalf("Finish running: %v", err)
	}
	if err := runs.MarkRunning(context.Background(), model.BackfillRun{ID: "r", Status: model.BackfillRunQueued}); !stderrors.Is(err, domainerrors.ErrInvalidBackfillTransition) {
		t.Fatalf("MarkRunning queued: %v", err)
	}
	if err := runs.RecordProgress(context.Background(), model.BackfillRun{ID: "r", Status: model.BackfillRunQueued}); !stderrors.Is(err, domainerrors.ErrInvalidBackfillTransition) {
		t.Fatalf("RecordProgress queued: %v", err)
	}
	beyond := model.BackfillRun{ID: "r", Status: model.BackfillRunFailed, SnapshotPosition: model.SourcePosition{SourceID: "s", Offset: math.MaxInt64 + 1}}
	if err := runs.Finish(context.Background(), beyond); !stderrors.Is(err, mapper.ErrPositionOutOfRange) {
		t.Fatalf("Finish mit Position jenseits bigint: %v", err)
	}
}
