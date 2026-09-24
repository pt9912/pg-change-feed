package postgresstorage_test

import (
	"context"
	stderrors "errors"
	"fmt"
	"testing"
	"time"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// writerRun trägt einen `running`-Run mit Position `X` und die Bindung seiner
// Tabelle: die Fremdschlüssel von `cdc.change` verlangen `source_table` und
// `schema_version`.
type writerRun struct {
	run       model.BackfillRun
	tableID   model.SourceTableID
	versionID model.SchemaVersionID
	runs      *postgresstorage.BackfillRunAdapter
	writer    *postgresstorage.BackfillWriterAdapter
}

// newWriterRun legt Bindung, Antrag und Run an und führt ihn bis `running`
// mit Snapshot-Position `offset`. Die Bereinigung entfernt genau die
// Zeilen dieses Präfixes: Changes, Transaktionen, Run und Antrag, dann
// Version und Bindung.
func newWriterRun(t *testing.T, f *backfillFixture, prefix string, offset uint64) *writerRun {
	t.Helper()
	admission := newBackfillAdmission(t, f)
	runs := newBackfillRunAdapter(t, f)
	writer, err := postgresstorage.NewBackfillWriter(context.Background(), f.dsn)
	if err != nil {
		t.Fatalf("NewBackfillWriter: %v", err)
	}
	t.Cleanup(writer.Close)

	tableID, versionID := prefix+"-tbl", prefix+"-sv"
	f.exec("INSERT INTO cdc.source_table (source_table_id, source_id, schema_name, table_name) VALUES ($1, $2, 'public', $3)",
		tableID, backfillTestSource, prefix)
	f.exec("INSERT INTO cdc.schema_version (schema_version_id, source_table_id, version) VALUES ($1, $2, 1)", versionID, tableID)
	f.t.Cleanup(func() {
		ctx := context.Background()
		_, _ = f.pool.Exec(ctx, "DELETE FROM cdc.change WHERE source_table_id = $1", tableID)
		_, _ = f.pool.Exec(ctx, "DELETE FROM cdc.transaction WHERE transaction_id LIKE $1", "%"+prefix+"%")
		_, _ = f.pool.Exec(ctx, "DELETE FROM cdc.schema_version WHERE schema_version_id = $1", versionID)
		_, _ = f.pool.Exec(ctx, "DELETE FROM cdc.source_table WHERE source_table_id = $1", tableID)
	})

	base := time.Date(2026, 9, 24, 18, 0, 0, 0, time.UTC)
	run := f.admit(admission, prefix, prefix, base, model.UnknownRowEstimate())
	run = mustRunning(t, runs, run, base)
	run, err = run.RecordProgress(backfillPosition(t, offset), 0)
	if err != nil {
		t.Fatalf("RecordProgress (Domäne): %v", err)
	}
	if err := runs.RecordProgress(context.Background(), run); err != nil {
		t.Fatalf("RecordProgress: %v", err)
	}
	return &writerRun{run: run, tableID: model.SourceTableID(tableID), versionID: model.SchemaVersionID(versionID), runs: runs, writer: writer}
}

// block baut den Block `number` mit `rows` Changes, wie ihn der Use Case
// baut: committed auf der Position `X` zum Snapshot-Zeitpunkt `at`.
func (w *writerRun) block(t *testing.T, number, rows int, at time.Time) *model.ChangeTransaction {
	t.Helper()
	id, err := model.BackfillTransactionID(w.run.ID, number)
	if err != nil {
		t.Fatalf("BackfillTransactionID: %v", err)
	}
	block, err := model.NewOpenTransaction(id, w.run.Source)
	if err != nil {
		t.Fatalf("NewOpenTransaction: %v", err)
	}
	for sequence := int64(1); sequence <= int64(rows); sequence++ {
		change, err := model.NewChange(model.ChangeIDFor(id, sequence), id, w.tableID, sequence, model.OperationInsert, nil,
			[]byte(fmt.Sprintf(`{"block":%d,"row":%d}`, number, sequence)), w.versionID)
		if err != nil {
			t.Fatalf("NewChange: %v", err)
		}
		if change, err = change.WithOrigin(model.ChangeOriginBackfill); err != nil {
			t.Fatalf("WithOrigin: %v", err)
		}
		if err := block.AppendChange(change); err != nil {
			t.Fatalf("AppendChange: %v", err)
		}
	}
	if err := block.Commit(w.run.SnapshotPosition, model.NewTimePoint(at.UnixNano())); err != nil {
		t.Fatalf("Commit: %v", err)
	}
	return block
}

// completed trägt den Run im Endzustand `completed` mit `rows` Changes.
func (w *writerRun) completed(t *testing.T, rows int64, at time.Time) model.BackfillRun {
	t.Helper()
	completed, err := w.run.Complete(model.NewTimePoint(at.UnixNano()), rows)
	if err != nil {
		t.Fatalf("Complete (Domäne): %v", err)
	}
	return completed
}

func (f *backfillFixture) backfillRows(prefix string) (transactions, changes int) {
	f.t.Helper()
	pattern := "0bf-" + prefix + "-%"
	return f.count("SELECT count(*) FROM cdc.transaction WHERE transaction_id LIKE $1", pattern),
		f.count("SELECT count(*) FROM cdc.change WHERE transaction_id LIKE $1", pattern)
}

// Atomarität: ein zweiter Leser sieht **vor** dem Commit keine Zeile des Runs
// und **danach** alle Blöcke zugleich; der Commit schreibt die Run-Zeile im
// selben Schritt `completed`. Jeder Change trägt `origin = 'backfill'`,
// `operation = 'INSERT'` und `old_data IS NULL`, die Transaktion die Position
// `X` und den Snapshot-Zeitpunkt. Rot färbende Mutation: `AppendBlock`
// committet je Block (eine Transaktion je Block) — die Vor-Commit-Zählung
// findet Zeilen.
func TestBackfillWriterIsInvisibleBeforeTheCommitAndCompleteAfterIt(t *testing.T) {
	f := newBackfillFixture(t)
	w := newWriterRun(t, f, "wr-atomic", 3_000_100)
	snapshotAt := time.Date(2026, 9, 24, 18, 30, 0, 0, time.UTC)

	tx, err := w.writer.Begin(context.Background(), w.run)
	if err != nil {
		t.Fatalf("Begin: %v", err)
	}
	t.Cleanup(func() { _ = tx.Rollback(context.Background()) })
	if err := tx.AppendBlock(context.Background(), w.block(t, 1, 2, snapshotAt)); err != nil {
		t.Fatalf("AppendBlock 1: %v", err)
	}
	if err := tx.AppendBlock(context.Background(), w.block(t, 2, 3, snapshotAt)); err != nil {
		t.Fatalf("AppendBlock 2: %v", err)
	}

	if transactions, changes := f.backfillRows("wr-atomic"); transactions != 0 || changes != 0 {
		t.Fatalf("vor dem Commit sichtbar: %d Transaktionen, %d Changes", transactions, changes)
	}
	if status, _ := f.runStatus("wr-atomic"); status != "running" {
		t.Fatalf("Run-Status vor dem Commit = %s", status)
	}

	if err := tx.Commit(context.Background(), w.completed(t, 5, snapshotAt.Add(time.Minute))); err != nil {
		t.Fatalf("Commit: %v", err)
	}

	if transactions, changes := f.backfillRows("wr-atomic"); transactions != 2 || changes != 5 {
		t.Fatalf("nach dem Commit: %d Transaktionen, %d Changes, erwartet 2/5", transactions, changes)
	}
	var status string
	var rowsCopied int64
	var finished time.Time
	f.scan("SELECT status, rows_copied, finished_at FROM cdc.backfill_run WHERE run_id = 'wr-atomic'", nil, &status, &rowsCopied, &finished)
	if status != "completed" || rowsCopied != 5 || !finished.Equal(snapshotAt.Add(time.Minute)) {
		t.Fatalf("Run-Zeile: %s, Zähler %d, finished_at %v", status, rowsCopied, finished)
	}
	if n := f.count(`SELECT count(*) FROM cdc.change WHERE transaction_id LIKE '0bf-wr-atomic-%'
	                 AND (origin IS DISTINCT FROM 'backfill' OR operation <> 'INSERT' OR old_data IS NOT NULL OR new_data IS NULL)`); n != 0 {
		t.Fatalf("%d Changes tragen nicht origin backfill / INSERT / ohne old_data", n)
	}
	if n := f.count(`SELECT count(*) FROM cdc.transaction WHERE transaction_id LIKE '0bf-wr-atomic-%'
	                 AND commit_position = 3000100 AND committed_at = $1 AND source_id = $2`, snapshotAt, backfillTestSource); n != 2 {
		t.Fatalf("%d von 2 Transaktionen tragen Position X und den Snapshot-Zeitpunkt", n)
	}
	var changeID, newData string
	f.scan(`SELECT change_id, new_data::text FROM cdc.change WHERE transaction_id = '0bf-wr-atomic-00000002' AND sequence = 3`, nil, &changeID, &newData)
	if changeID != "0bf-wr-atomic-00000002-3" || newData != `{"row": 3, "block": 2}` {
		t.Fatalf("Change = %s %s", changeID, newData)
	}
}

// Rollback: ein verworfener Run hinterlässt weder `cdc.transaction`- noch
// `cdc.change`-Zeile (`LH-FA-CAP-009` Negative), die Run-Zeile bleibt
// `running` und kann `failed` werden; ein wiederholter Rollback ist ohne
// Wirkung.
func TestBackfillWriterRollbackLeavesNoRow(t *testing.T) {
	f := newBackfillFixture(t)
	w := newWriterRun(t, f, "wr-rollback", 3_000_200)
	snapshotAt := time.Date(2026, 9, 24, 18, 30, 0, 0, time.UTC)

	tx, err := w.writer.Begin(context.Background(), w.run)
	if err != nil {
		t.Fatalf("Begin: %v", err)
	}
	if err := tx.AppendBlock(context.Background(), w.block(t, 1, 4, snapshotAt)); err != nil {
		t.Fatalf("AppendBlock: %v", err)
	}
	if err := tx.Rollback(context.Background()); err != nil {
		t.Fatalf("Rollback: %v", err)
	}
	if err := tx.Rollback(context.Background()); err != nil {
		t.Fatalf("wiederholter Rollback: %v", err)
	}

	if transactions, changes := f.backfillRows("wr-rollback"); transactions != 0 || changes != 0 {
		t.Fatalf("nach dem Rollback: %d Transaktionen, %d Changes", transactions, changes)
	}
	if status, _ := f.runStatus("wr-rollback"); status != "running" {
		t.Fatalf("Run-Status nach dem Rollback = %s, erwartet running", status)
	}
	failed, err := w.run.Fail(model.NewTimePoint(snapshotAt.UnixNano()), model.ErrorClassInternal, "verworfen")
	if err != nil {
		t.Fatalf("Fail: %v", err)
	}
	if err := w.runs.Finish(context.Background(), failed); err != nil {
		t.Fatalf("Finish: %v", err)
	}
	if err := tx.Commit(context.Background(), w.completed(t, 4, snapshotAt)); err == nil {
		t.Fatal("Commit nach dem Rollback gelingt")
	}
}

// Ein Commit auf einen Run, der zwischenzeitlich nicht mehr `running` ist
// (hier: `interrupted` durch den Abgleich eines anderen Prozesses), committet
// nichts: die Run-Zeile trifft `WHERE status = 'running'` nicht, der Commit
// endet als `ErrInvalidBackfillTransition`, und der Rollback hinterlässt keine
// Zeile. Rot färbende Mutation: in `UpdateBackfillRunCompleted` die Klausel
// `AND status = 'running'` entfernen — der beendete Run wird überschrieben und
// die Daten werden committet.
func TestBackfillWriterDoesNotCommitAnEndedRun(t *testing.T) {
	f := newBackfillFixture(t)
	w := newWriterRun(t, f, "wr-ended", 3_000_300)
	snapshotAt := time.Date(2026, 9, 24, 18, 30, 0, 0, time.UTC)

	tx, err := w.writer.Begin(context.Background(), w.run)
	if err != nil {
		t.Fatalf("Begin: %v", err)
	}
	if err := tx.AppendBlock(context.Background(), w.block(t, 1, 2, snapshotAt)); err != nil {
		t.Fatalf("AppendBlock: %v", err)
	}
	// Der Abgleich läuft auf einer anderen Verbindung und wartet auf keine
	// Sperre der offenen Schreibtransaktion: sie hat die Run-Zeile nicht berührt.
	if _, err := w.runs.InterruptRunning(context.Background(), backfillTestSource, model.NewTimePoint(snapshotAt.UnixNano())); err != nil {
		t.Fatalf("InterruptRunning: %v", err)
	}

	err = tx.Commit(context.Background(), w.completed(t, 2, snapshotAt))
	if !stderrors.Is(err, domainerrors.ErrInvalidBackfillTransition) {
		t.Fatalf("Commit: Fehler = %v, erwartet ErrInvalidBackfillTransition", err)
	}
	if err := tx.Rollback(context.Background()); err != nil {
		t.Fatalf("Rollback: %v", err)
	}
	if transactions, changes := f.backfillRows("wr-ended"); transactions != 0 || changes != 0 {
		t.Fatalf("Zeilen eines beendeten Runs: %d Transaktionen, %d Changes", transactions, changes)
	}
	if status, _ := f.runStatus("wr-ended"); status != "interrupted" {
		t.Fatalf("Run-Status = %s, erwartet interrupted", status)
	}
}

// Ein Block, der den Vertrag des Schreibers verletzt, schreibt nichts (die
// Prüfung im Einzelnen belegt `backfill_internal_test.go` netzlos); an der
// realen Datenbank bleibt die Transaktion nach der Ablehnung benutzbar.
func TestBackfillWriterRejectsAContractViolationAndStaysUsable(t *testing.T) {
	f := newBackfillFixture(t)
	w := newWriterRun(t, f, "wr-reject", 3_000_400)
	snapshotAt := time.Date(2026, 9, 24, 18, 30, 0, 0, time.UTC)

	tx, err := w.writer.Begin(context.Background(), w.run)
	if err != nil {
		t.Fatalf("Begin: %v", err)
	}
	t.Cleanup(func() { _ = tx.Rollback(context.Background()) })
	if err := tx.AppendBlock(context.Background(), w.block(t, 2, 1, snapshotAt)); !stderrors.Is(err, outbound.ErrBackfillBlockInvalid) {
		t.Fatalf("Block 2 zuerst: %v", err)
	}
	if err := tx.AppendBlock(context.Background(), w.block(t, 1, 1, snapshotAt)); err != nil {
		t.Fatalf("Block 1 nach der Ablehnung: %v", err)
	}
	if err := tx.Commit(context.Background(), w.completed(t, 1, snapshotAt)); err != nil {
		t.Fatalf("Commit: %v", err)
	}
	if transactions, changes := f.backfillRows("wr-reject"); transactions != 1 || changes != 1 {
		t.Fatalf("Zeilen: %d Transaktionen, %d Changes, erwartet 1/1", transactions, changes)
	}
}

// Der Fortschritt läuft außerhalb der Daten-Transaktion auf einer zweiten
// Verbindung, ist sofort für Leser sichtbar und wartet auf keine Sperre der
// offenen Schreibtransaktion; der Commit derselben Run-Zeile danach kollidiert
// nicht mit ihm. Rot färbende Mutation: `AppendBlock` schreibt die Run-Zeile
// (etwa `rows_copied`) in der Schreibtransaktion — das Fortschritts-Update
// wartet dann auf deren Zeilensperre, der Aufruf endet an der Zeitgrenze des
// Tests.
func TestBackfillProgressDoesNotWaitForTheOpenWriteTransaction(t *testing.T) {
	f := newBackfillFixture(t)
	w := newWriterRun(t, f, "wr-progress", 3_000_500)
	snapshotAt := time.Date(2026, 9, 24, 18, 30, 0, 0, time.UTC)

	tx, err := w.writer.Begin(context.Background(), w.run)
	if err != nil {
		t.Fatalf("Begin: %v", err)
	}
	t.Cleanup(func() { _ = tx.Rollback(context.Background()) })
	if err := tx.AppendBlock(context.Background(), w.block(t, 1, 2, snapshotAt)); err != nil {
		t.Fatalf("AppendBlock: %v", err)
	}

	progressed, err := w.run.RecordProgress(w.run.SnapshotPosition, 2)
	if err != nil {
		t.Fatalf("RecordProgress (Domäne): %v", err)
	}
	progressed.WarnDuration = true
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	start := time.Now()
	if err := w.runs.RecordProgress(ctx, progressed); err != nil {
		t.Fatalf("RecordProgress während der offenen Schreibtransaktion: %v", err)
	}
	if elapsed := time.Since(start); elapsed > 2*time.Second {
		t.Fatalf("RecordProgress wartete %v auf die Schreibtransaktion", elapsed)
	}

	var rowsCopied int64
	var warn bool
	f.scan("SELECT rows_copied, warn_duration FROM cdc.backfill_run WHERE run_id = 'wr-progress'", nil, &rowsCopied, &warn)
	if rowsCopied != 2 || !warn {
		t.Fatalf("sichtbarer Fortschritt: Zähler %d, warn_duration %v", rowsCopied, warn)
	}
	if transactions, changes := f.backfillRows("wr-progress"); transactions != 0 || changes != 0 {
		t.Fatalf("Daten vor dem Commit sichtbar: %d/%d", transactions, changes)
	}

	completed, err := progressed.Complete(model.NewTimePoint(snapshotAt.UnixNano()), 2)
	if err != nil {
		t.Fatalf("Complete (Domäne): %v", err)
	}
	if err := tx.Commit(context.Background(), completed); err != nil {
		t.Fatalf("Commit nach dem Fortschritt: %v", err)
	}
	var status string
	f.scan("SELECT status, warn_duration FROM cdc.backfill_run WHERE run_id = 'wr-progress'", nil, &status, &warn)
	if status != "completed" || !warn {
		t.Fatalf("nach dem Commit: %s, warn_duration %v", status, warn)
	}
}

// Ordnung (`ADR-0111` Teilfrage 6, `LH-FA-CAP-004`, `LH-FA-REA-004`): mit
// einem WAL-Commit **auf derselben Position `X`** liest `cdc.changes` die
// Backfill-Blöcke `0bf-…` vor dem WAL-Commit, nach `(commit_position,
// transaction_id, sequence)` — im Lesezugriff des Store-Adapters ebenso wie in
// der View. Die Ordnung ist eine Textsortierung der Datenbank; der Test hält
// die Kollation der Datenbank im Log fest. Rot färbende Mutation: der Präfix
// `backfillTransactionPrefix` (`internal/domain/model/backfillrun.go`) beginnt
// mit einem Zeichen, das hinter Ziffern sortiert (`bf-` statt `0bf-`) — der
// WAL-Commit steht vor den Blöcken.
func TestBackfillBlocksSortBeforeTheWALTransactionOnTheSamePosition(t *testing.T) {
	f := newBackfillFixture(t)
	const offset = 9_100_000_000
	w := newWriterRun(t, f, "wr-order", offset)
	snapshotAt := time.Date(2026, 9, 24, 18, 30, 0, 0, time.UTC)

	var collation string
	f.scan("SELECT datcollate FROM pg_database WHERE datname = current_database()", nil, &collation)
	t.Logf("datcollate der Test-Datenbank: %s", collation)

	tx, err := w.writer.Begin(context.Background(), w.run)
	if err != nil {
		t.Fatalf("Begin: %v", err)
	}
	t.Cleanup(func() { _ = tx.Rollback(context.Background()) })
	if err := tx.AppendBlock(context.Background(), w.block(t, 1, 2, snapshotAt)); err != nil {
		t.Fatalf("AppendBlock 1: %v", err)
	}
	if err := tx.AppendBlock(context.Background(), w.block(t, 2, 1, snapshotAt)); err != nil {
		t.Fatalf("AppendBlock 2: %v", err)
	}
	if err := tx.Commit(context.Background(), w.completed(t, 3, snapshotAt)); err != nil {
		t.Fatalf("Commit: %v", err)
	}

	// WAL-Kennungen beginnen mit einer Ziffer ungleich 0; die Werte decken
	// kurze, lange sowie lexikographisch kleine und große ab, das Suffix
	// ordnet sie dem Test zu.
	store, err := postgresstorage.New(context.Background(), f.dsn)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	t.Cleanup(store.Close)
	walIDs := []string{"1", "748", "9999999999"}
	for _, id := range walIDs {
		wal, err := model.NewOpenTransaction(model.TransactionID(id+"wr-order"), backfillTestSource)
		if err != nil {
			t.Fatalf("NewOpenTransaction: %v", err)
		}
		change, err := model.NewChange(model.ChangeID(id+"wr-order-1"), wal.ID, w.tableID, 1, model.OperationInsert, nil, []byte(`{"wal":true}`), w.versionID)
		if err != nil {
			t.Fatalf("NewChange: %v", err)
		}
		if err := wal.AppendChange(change); err != nil {
			t.Fatalf("AppendChange: %v", err)
		}
		if err := wal.Commit(w.run.SnapshotPosition, model.NewTimePoint(snapshotAt.UnixNano())); err != nil {
			t.Fatalf("Commit: %v", err)
		}
		if err := store.PersistTransaction(context.Background(), wal); err != nil {
			t.Fatalf("PersistTransaction %s: %v", id, err)
		}
	}

	start, end := w.run.SnapshotPosition, model.SourcePosition{SourceID: backfillTestSource, Offset: offset + 1}
	records, err := store.ReadChanges(context.Background(), outbound.ChangeQuery{Source: backfillTestSource, Start: &start, End: &end, Schema: "public", Table: "wr-order"})
	if err != nil {
		t.Fatalf("ReadChanges: %v", err)
	}
	var order []string
	for _, record := range records {
		order = append(order, fmt.Sprintf("%s/%d/%s", record.Change.TransactionID, record.Change.Sequence, record.Change.Origin))
	}
	block1, err := model.BackfillTransactionID(w.run.ID, 1)
	if err != nil {
		t.Fatalf("BackfillTransactionID: %v", err)
	}
	block2, err := model.BackfillTransactionID(w.run.ID, 2)
	if err != nil {
		t.Fatalf("BackfillTransactionID: %v", err)
	}
	// Die Soll-Ordnung nennt die Block-Kennungen über `model.BackfillTransactionID`,
	// nicht als Literal: ein Präfix, der hinter Ziffern sortiert, färbt die
	// Ordnung rot und nicht bloß den Text-Vergleich.
	want := []string{
		fmt.Sprintf("%s/1/backfill", block1), fmt.Sprintf("%s/2/backfill", block1), fmt.Sprintf("%s/1/backfill", block2),
		"1wr-order/1/wal", "748wr-order/1/wal", "9999999999wr-order/1/wal",
	}
	if fmt.Sprint(order) != fmt.Sprint(want) {
		t.Fatalf("Ordnung des Lesezugriffs (datcollate %s):\n got %v\nwant %v", collation, order, want)
	}

	rows, err := f.pool.Query(context.Background(),
		"SELECT transaction_id, sequence, origin FROM cdc.changes WHERE source_id = $1 AND commit_position = $2 AND table_name = 'wr-order' ORDER BY commit_position, transaction_id, sequence",
		backfillTestSource, int64(offset))
	if err != nil {
		t.Fatalf("View cdc.changes: %v", err)
	}
	defer rows.Close()
	var viewOrder []string
	for rows.Next() {
		var id, origin string
		var sequence int64
		if err := rows.Scan(&id, &sequence, &origin); err != nil {
			t.Fatalf("Scan: %v", err)
		}
		viewOrder = append(viewOrder, fmt.Sprintf("%s/%d/%s", id, sequence, origin))
	}
	if fmt.Sprint(viewOrder) != fmt.Sprint(want) {
		t.Fatalf("Ordnung der View (datcollate %s):\n got %v\nwant %v", collation, viewOrder, want)
	}
}
