package model

import (
	"fmt"
	"sort"
	"testing"

	stderrors "errors"

	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

// openTransaction liefert eine gültige offene Transaktion für Tests.
func openTransaction(t *testing.T, id TransactionID, source SourceID) *ChangeTransaction {
	t.Helper()
	tx, err := NewOpenTransaction(id, source)
	if err != nil {
		t.Fatalf("Transaktion-Konstruktor: %v", err)
	}
	return tx
}

// LH-FA-CAP-005, Happy Path: drei Änderungen derselben Quelltransaktion
// tragen dieselbe Transaktionskennung.
func TestLHFACAP005ChangesCarrySameTransactionID(t *testing.T) {
	tx := openTransaction(t, "tx-1", "src-1")
	for i := int64(1); i <= 3; i++ {
		change, err := NewChange(ChangeID(fmt.Sprintf("ch-%d", i)), "tx-1", "tbl-1", i, OperationInsert, nil, []byte(`{}`), "sv-1")
		if err != nil {
			t.Fatalf("Change %d: %v", i, err)
		}
		if err := tx.AppendChange(change); err != nil {
			t.Fatalf("Change %d anhängen: %v", i, err)
		}
	}
	changes := tx.Changes()
	if len(changes) != 3 {
		t.Fatalf("Transaktion trägt %d Changes, wollen 3", len(changes))
	}
	for _, change := range changes {
		if change.TransactionID != "tx-1" {
			t.Fatalf("Change %q trägt Transaktion %q, wollen \"tx-1\"", change.ID, change.TransactionID)
		}
	}
}

// LH-FA-CAP-005, Boundary: zwei Transaktionen schreiben zeitgleich — ihre
// Changes sind unterscheidbar zugeordnet.
func TestLHFACAP005ConcurrentTransactionsAreDistinguishable(t *testing.T) {
	first := openTransaction(t, "tx-1", "src-1")
	second := openTransaction(t, "tx-2", "src-1")
	for _, pair := range []struct {
		tx  *ChangeTransaction
		id  ChangeID
		seq int64
	}{
		{first, "ch-1", 1},
		{second, "ch-2", 1},
	} {
		change, err := NewChange(pair.id, pair.tx.ID, "tbl-1", pair.seq, OperationInsert, nil, []byte(`{}`), "sv-1")
		if err != nil {
			t.Fatalf("Change %q: %v", pair.id, err)
		}
		if err := pair.tx.AppendChange(change); err != nil {
			t.Fatalf("Change %q anhängen: %v", pair.id, err)
		}
	}
	for _, tx := range []*ChangeTransaction{first, second} {
		for _, change := range tx.Changes() {
			if change.TransactionID != tx.ID {
				t.Fatalf("Change %q hängt an Transaktion %q, trägt aber %q", change.ID, tx.ID, change.TransactionID)
			}
		}
	}
}

// LH-FA-CAP-004, Happy Path: zwei Änderungen an derselben Zeile in zwei
// aufeinanderfolgenden Transaktionen — die Ausführungsreihenfolge
// rekonstruiert sich aus den Commit-Positionen.
func TestLHFACAP004ExecutionOrderReconstructibleAcrossTransactions(t *testing.T) {
	first := openTransaction(t, "tx-1", "src-1")
	second := openTransaction(t, "tx-2", "src-1")
	firstChange := buildChange(t, changeArgs{id: "ch-1", tx: "tx-1", table: "tbl-1", seq: 1, op: OperationUpdate, oldData: []byte(`{"v":1}`), newData: []byte(`{"v":2}`), sv: "sv-1"})
	secondChange := buildChange(t, changeArgs{id: "ch-2", tx: "tx-2", table: "tbl-1", seq: 1, op: OperationUpdate, oldData: []byte(`{"v":2}`), newData: []byte(`{"v":3}`), sv: "sv-1"})
	if err := first.AppendChange(firstChange); err != nil {
		t.Fatalf("Change 1 anhängen: %v", err)
	}
	if err := second.AppendChange(secondChange); err != nil {
		t.Fatalf("Change 2 anhängen: %v", err)
	}
	firstPosition, err := NewSourcePosition("src-1", 100)
	if err != nil {
		t.Fatalf("Position 1: %v", err)
	}
	secondPosition, err := NewSourcePosition("src-1", 200)
	if err != nil {
		t.Fatalf("Position 2: %v", err)
	}
	if err := first.Commit(firstPosition); err != nil {
		t.Fatalf("Commit 1: %v", err)
	}
	if err := second.Commit(secondPosition); err != nil {
		t.Fatalf("Commit 2: %v", err)
	}
	firstCommit, _ := first.CommitPosition()
	secondCommit, _ := second.CommitPosition()
	if c := secondCommit.Compare(firstCommit); c <= 0 {
		t.Fatalf("Ausführungsreihenfolge nicht rekonstruierbar: Commit-Positionen %d vs %d (Compare %d)", firstCommit.Offset, secondCommit.Offset, c)
	}
}

// LH-FA-DAT-004, Boundary: mehrere Changes derselben Transaktion erhalten
// ihre Reihenfolge innerhalb der Transaktion (Sequenz).
func TestLHFADAT004IntraTransactionOrderPreservedBySequence(t *testing.T) {
	tx := openTransaction(t, "tx-1", "src-1")
	for _, seq := range []int64{3, 1, 2} {
		change, err := NewChange(ChangeID(fmt.Sprintf("ch-%d", seq)), "tx-1", "tbl-1", seq, OperationInsert, nil, []byte(`{}`), "sv-1")
		if err != nil {
			t.Fatalf("Change %d: %v", seq, err)
		}
		if err := tx.AppendChange(change); err != nil {
			t.Fatalf("Change %d anhängen: %v", seq, err)
		}
	}
	changes := tx.Changes()
	sort.Slice(changes, func(i, j int) bool { return changes[i].Sequence < changes[j].Sequence })
	for i, change := range changes {
		if change.Sequence != int64(i+1) {
			t.Fatalf("Reihenfolge innerhalb der Transaktion nicht erhalten: Stelle %d trägt Sequenz %d", i+1, change.Sequence)
		}
	}
}

// Die Transaktions-Invarianten (`ADR-0029`, Regel 3 und 6) scheitern am
// Aggregat: falsche Zuordnung, doppelte Sequenz, Anhang nach Commit,
// doppelter Commit, Commit an fremder Quelle.
func TestChangeTransactionRejectsInvariantViolations(t *testing.T) {
	t.Run("fremde Transaktion", func(t *testing.T) {
		tx := openTransaction(t, "tx-1", "src-1")
		change := buildChange(t, changeArgs{id: "ch-1", tx: "tx-2", table: "tbl-1", seq: 1, op: OperationInsert, newData: []byte(`{}`), sv: "sv-1"})
		err := tx.AppendChange(change)
		if !stderrors.Is(err, domainerrors.ErrTransactionMismatch) {
			t.Fatalf("Fehler = %v, wollen ErrTransactionMismatch", err)
		}
	})
	t.Run("doppelte Sequenz", func(t *testing.T) {
		tx := openTransaction(t, "tx-1", "src-1")
		first := buildChange(t, changeArgs{id: "ch-1", tx: "tx-1", table: "tbl-1", seq: 1, op: OperationInsert, newData: []byte(`{}`), sv: "sv-1"})
		second := buildChange(t, changeArgs{id: "ch-2", tx: "tx-1", table: "tbl-1", seq: 1, op: OperationInsert, newData: []byte(`{}`), sv: "sv-1"})
		if err := tx.AppendChange(first); err != nil {
			t.Fatalf("erster Anhang: %v", err)
		}
		err := tx.AppendChange(second)
		if !stderrors.Is(err, domainerrors.ErrDuplicateSequence) {
			t.Fatalf("Fehler = %v, wollen ErrDuplicateSequence", err)
		}
	})
	t.Run("Anhang nach Commit", func(t *testing.T) {
		tx := openTransaction(t, "tx-1", "src-1")
		change := buildChange(t, changeArgs{id: "ch-1", tx: "tx-1", table: "tbl-1", seq: 1, op: OperationInsert, newData: []byte(`{}`), sv: "sv-1"})
		if err := tx.AppendChange(change); err != nil {
			t.Fatalf("Anhang: %v", err)
		}
		position, err := NewSourcePosition("src-1", 100)
		if err != nil {
			t.Fatalf("Position: %v", err)
		}
		if err := tx.Commit(position); err != nil {
			t.Fatalf("Commit: %v", err)
		}
		late := buildChange(t, changeArgs{id: "ch-2", tx: "tx-1", table: "tbl-1", seq: 2, op: OperationInsert, newData: []byte(`{}`), sv: "sv-1"})
		err = tx.AppendChange(late)
		if !stderrors.Is(err, domainerrors.ErrTransactionAlreadyCommitted) {
			t.Fatalf("Fehler = %v, wollen ErrTransactionAlreadyCommitted", err)
		}
	})
	t.Run("doppelter Commit", func(t *testing.T) {
		tx := openTransaction(t, "tx-1", "src-1")
		position, err := NewSourcePosition("src-1", 100)
		if err != nil {
			t.Fatalf("Position: %v", err)
		}
		if err := tx.Commit(position); err != nil {
			t.Fatalf("erster Commit: %v", err)
		}
		err = tx.Commit(position)
		if !stderrors.Is(err, domainerrors.ErrTransactionAlreadyCommitted) {
			t.Fatalf("Fehler = %v, wollen ErrTransactionAlreadyCommitted", err)
		}
	})
	t.Run("Commit an fremder Quelle", func(t *testing.T) {
		tx := openTransaction(t, "tx-1", "src-1")
		position, err := NewSourcePosition("src-2", 100)
		if err != nil {
			t.Fatalf("Position: %v", err)
		}
		err = tx.Commit(position)
		if !stderrors.Is(err, domainerrors.ErrSourceMismatch) {
			t.Fatalf("Fehler = %v, wollen ErrSourceMismatch", err)
		}
	})
}

// LH-FA-CAP-006, Happy Path: solange die Transaktion offen ist, trägt sie
// keine Commit-Position und ihre Changes sind nicht konsumierbar.
func TestLHFACAP006OpenTransactionHasNoCommitPosition(t *testing.T) {
	tx := openTransaction(t, "tx-1", "src-1")
	if tx.IsCommitted() {
		t.Fatal("offene Transaktion meldet committed")
	}
	if _, committed := tx.CommitPosition(); committed {
		t.Fatal("offene Transaktion trägt eine Commit-Position")
	}
}
