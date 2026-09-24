package model

import (
	stderrors "errors"
	"testing"

	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

func queuedRun(t *testing.T) BackfillRun {
	t.Helper()
	run, err := NewQueuedBackfillRun("run-1", "src-1", "public", "orders", NewTimePoint(100))
	if err != nil {
		t.Fatalf("NewQueuedBackfillRun: %v", err)
	}
	return run
}

func runningRun(t *testing.T) BackfillRun {
	t.Helper()
	run, err := queuedRun(t).Start(NewTimePoint(200))
	if err != nil {
		t.Fatalf("Start: %v", err)
	}
	return run
}

// TestNewQueuedBackfillRun trägt den Anfangszustand (`SPEC-029`): `queued`,
// die Schätzung unbekannt, kein Fortschritt, beide Warn-Kennzeichnungen
// `false`; eine leere Kennung endet mit `ErrEmptyIdentifier`.
func TestNewQueuedBackfillRun(t *testing.T) {
	run := queuedRun(t)
	if run.Status != BackfillRunQueued || run.RowsCopied != 0 || !run.StartedAt.Unset() || !run.FinishedAt.Unset() {
		t.Fatalf("Anfangszustand = %+v", run)
	}
	if _, known := run.EstimatedRows.Rows(); known {
		t.Fatalf("Schätzung des neuen Runs ist bekannt, will unbekannt")
	}
	if run.WarnEstimatedSize || run.WarnDuration || !run.SnapshotPosition.IsZero() || run.ErrorMessage != "" {
		t.Fatalf("Warn-Kennzeichnungen, Position oder Fehlertext gesetzt: %+v", run)
	}
	if run.QualifiedName() != "public.orders" {
		t.Fatalf("QualifiedName = %q", run.QualifiedName())
	}
	for _, args := range [][4]string{{"", "src", "s", "t"}, {"r", "", "s", "t"}, {"r", "src", "", "t"}, {"r", "src", "s", ""}} {
		_, err := NewQueuedBackfillRun(BackfillRunID(args[0]), SourceID(args[1]), args[2], args[3], NewTimePoint(1))
		if !stderrors.Is(err, domainerrors.ErrEmptyIdentifier) {
			t.Fatalf("NewQueuedBackfillRun(%q) = %v, wollen ErrEmptyIdentifier", args, err)
		}
	}
}

// TestRowEstimate trägt die Unterscheidung „unbekannt" gegen `0`
// (`SPEC-029`): der Nullwert ist unbekannt, `NewRowEstimate(0)` ist eine
// bekannte 0, eine negative Zahl ist keine Schätzung.
func TestRowEstimate(t *testing.T) {
	if _, known := (RowEstimate{}).Rows(); known {
		t.Fatal("Nullwert ist bekannt")
	}
	if _, known := UnknownRowEstimate().Rows(); known {
		t.Fatal("UnknownRowEstimate ist bekannt")
	}
	zero, err := NewRowEstimate(0)
	if err != nil {
		t.Fatalf("NewRowEstimate(0): %v", err)
	}
	if rows, known := zero.Rows(); !known || rows != 0 {
		t.Fatalf("NewRowEstimate(0) = (%d, %t), will (0, true)", rows, known)
	}
	if rows, known := mustEstimate(t, 42).Rows(); !known || rows != 42 {
		t.Fatalf("NewRowEstimate(42) = (%d, %t)", rows, known)
	}
	if _, err := NewRowEstimate(-1); !stderrors.Is(err, domainerrors.ErrNegativeRowCount) {
		t.Fatalf("NewRowEstimate(-1) = %v, will ErrNegativeRowCount", err)
	}
	withEstimate := queuedRun(t).WithEstimatedRows(mustEstimate(t, 7))
	if rows, known := withEstimate.EstimatedRows.Rows(); !known || rows != 7 {
		t.Fatalf("WithEstimatedRows trägt (%d, %t)", rows, known)
	}
}

func mustEstimate(t *testing.T, rows int64) RowEstimate {
	t.Helper()
	estimate, err := NewRowEstimate(rows)
	if err != nil {
		t.Fatalf("NewRowEstimate(%d): %v", rows, err)
	}
	return estimate
}

// TestBackfillRunTransitions trägt die zulässigen und unzulässigen Übergänge
// (`SPEC-029`) und den Wert-Vertrag: jede Methode lässt den Empfänger
// unverändert.
func TestBackfillRunTransitions(t *testing.T) {
	position, err := NewSourcePosition("src-1", 0x100)
	if err != nil {
		t.Fatal(err)
	}
	start := func(r BackfillRun) (BackfillRun, error) { return r.Start(NewTimePoint(200)) }
	progress := func(r BackfillRun) (BackfillRun, error) { return r.RecordProgress(position, 3) }
	complete := func(r BackfillRun) (BackfillRun, error) { return r.Complete(NewTimePoint(300), 3) }
	fail := func(r BackfillRun) (BackfillRun, error) {
		return r.Fail(NewTimePoint(300), ErrorClassStorage, "Schreibfehler")
	}
	interrupt := func(r BackfillRun) (BackfillRun, error) { return r.Interrupt(NewTimePoint(300)) }

	queued := queuedRun(t)
	running := runningRun(t)
	completed, _ := complete(running)
	failed, _ := fail(running)
	interrupted, _ := interrupt(running)

	cases := []struct {
		name    string
		from    BackfillRun
		do      func(BackfillRun) (BackfillRun, error)
		want    BackfillRunStatus
		invalid bool
	}{
		{"queued -> running", queued, start, BackfillRunRunning, false},
		{"running -> running (Start)", running, start, "", true},
		{"terminal -> running (Start)", completed, start, "", true},
		{"queued Fortschritt", queued, progress, "", true},
		{"running Fortschritt", running, progress, BackfillRunRunning, false},
		{"completed Fortschritt", completed, progress, "", true},
		{"queued -> completed", queued, complete, "", true},
		{"running -> completed", running, complete, BackfillRunCompleted, false},
		{"failed -> completed", failed, complete, "", true},
		{"queued -> failed", queued, fail, BackfillRunFailed, false},
		{"running -> failed", running, fail, BackfillRunFailed, false},
		{"completed -> failed", completed, fail, "", true},
		{"failed -> failed", failed, fail, "", true},
		{"interrupted -> failed", interrupted, fail, "", true},
		{"queued -> interrupted", queued, interrupt, "", true},
		{"running -> interrupted", running, interrupt, BackfillRunInterrupted, false},
		{"completed -> interrupted", completed, interrupt, "", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			before := tc.from
			got, err := tc.do(tc.from)
			if tc.from != before {
				t.Fatalf("der Empfänger wurde verändert: %+v → %+v", before, tc.from)
			}
			if tc.invalid {
				if !stderrors.Is(err, domainerrors.ErrInvalidBackfillTransition) {
					t.Fatalf("Fehler = %v, will ErrInvalidBackfillTransition", err)
				}
				if got != before {
					t.Fatalf("ein abgelehnter Übergang lieferte einen anderen Run: %+v", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("Übergang = %v, will nil", err)
			}
			if got.Status != tc.want {
				t.Fatalf("Status = %q, will %q", got.Status, tc.want)
			}
		})
	}
}

// TestBackfillRunFields trägt, welches Feld welcher Übergang setzt:
// `started_at` beim Start, `finished_at` und Fehlertext im Endzustand.
func TestBackfillRunFields(t *testing.T) {
	running := runningRun(t)
	if running.StartedAt != NewTimePoint(200) || !running.FinishedAt.Unset() {
		t.Fatalf("Start setzt StartedAt/FinishedAt = %v/%v", running.StartedAt, running.FinishedAt)
	}
	position, _ := NewSourcePosition("src-1", 0x100)
	progressed, err := running.RecordProgress(position, 5)
	if err != nil || progressed.RowsCopied != 5 || progressed.SnapshotPosition != position {
		t.Fatalf("RecordProgress = %+v, %v", progressed, err)
	}
	completed, err := progressed.Complete(NewTimePoint(300), 9)
	if err != nil || completed.RowsCopied != 9 || completed.FinishedAt != NewTimePoint(300) || completed.ErrorMessage != "" {
		t.Fatalf("Complete = %+v, %v", completed, err)
	}
	failed, err := progressed.Fail(NewTimePoint(300), ErrorClassPermission, "kein SELECT")
	if err != nil || failed.ErrorMessage != "permission: kein SELECT" || failed.FinishedAt != NewTimePoint(300) || failed.RowsCopied != 5 {
		t.Fatalf("Fail = %+v, %v", failed, err)
	}
	interrupted, err := progressed.Interrupt(NewTimePoint(300))
	if err != nil || interrupted.ErrorMessage != "" || interrupted.FinishedAt != NewTimePoint(300) || interrupted.RowsCopied != 5 {
		t.Fatalf("Interrupt = %+v, %v", interrupted, err)
	}
	for _, run := range []BackfillRun{queuedRun(t), running, progressed} {
		if !run.IsActive() {
			t.Fatalf("%q ist nicht aktiv", run.Status)
		}
	}
	for _, run := range []BackfillRun{completed, failed, interrupted} {
		if run.IsActive() {
			t.Fatalf("%q ist aktiv", run.Status)
		}
	}
}

// TestBackfillRunProgressGuards trägt die Eingabe-Grenzen: der Zähler wächst
// nur (Fortschritt und Abschluss), die Position gehört zur Quelle des Runs,
// eine unbekannte Fehlerklasse wird abgelehnt.
func TestBackfillRunProgressGuards(t *testing.T) {
	position, _ := NewSourcePosition("src-1", 0x100)
	other, _ := NewSourcePosition("src-2", 0x100)
	progressed, err := runningRun(t).RecordProgress(position, 5)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := progressed.RecordProgress(position, 4); !stderrors.Is(err, domainerrors.ErrBackfillProgressRegression) {
		t.Fatalf("Rückschritt = %v, will ErrBackfillProgressRegression", err)
	}
	if _, err := progressed.Complete(NewTimePoint(300), 4); !stderrors.Is(err, domainerrors.ErrBackfillProgressRegression) {
		t.Fatalf("Complete mit kleinerem Zähler = %v, will ErrBackfillProgressRegression", err)
	}
	if _, err := progressed.RecordProgress(other, 6); !stderrors.Is(err, domainerrors.ErrSourceMismatch) {
		t.Fatalf("Position einer anderen Quelle = %v, will ErrSourceMismatch", err)
	}
	if _, err := progressed.Fail(NewTimePoint(300), ErrorClass("bogus"), "x"); !stderrors.Is(err, domainerrors.ErrInvalidErrorClass) {
		t.Fatalf("unbekannte Klasse = %v, will ErrInvalidErrorClass", err)
	}
	if _, err := progressed.RecordProgress(position, 5); err != nil {
		t.Fatalf("gleicher Zähler ist zulässig: %v", err)
	}
}

// TestBackfillTransactionID trägt die Bildungsregel der Transaktions-Kennung
// (`ADR-0111` Teilfrage 6): Präfix `0bf-`, Run-Kennung, achtstellige
// null-aufgefüllte Blocknummer; die Kennung sortiert vor jeder WAL-Kennung
// (Ziffernfolge ohne führende Null) und die Blöcke untereinander in
// Blockreihenfolge.
func TestBackfillTransactionID(t *testing.T) {
	first, err := BackfillTransactionID("run-1", 1)
	if err != nil || first != "0bf-run-1-00000001" {
		t.Fatalf("Block 1 = %q, %v", first, err)
	}
	tenth, _ := BackfillTransactionID("run-1", 10)
	last, err := BackfillTransactionID("run-1", 99_999_999)
	if err != nil || last != "0bf-run-1-99999999" {
		t.Fatalf("Block 99999999 = %q, %v", last, err)
	}
	if !(first < tenth && tenth < last) {
		t.Fatalf("die Blockreihenfolge ist nicht die Sortierreihenfolge: %q %q %q", first, tenth, last)
	}
	for _, wal := range []string{"1", "745", "4294967295"} {
		if !(string(last) < wal) {
			t.Fatalf("%q sortiert nicht vor der WAL-Kennung %q", last, wal)
		}
	}
	for _, block := range []int{0, -1, 100_000_000} {
		if _, err := BackfillTransactionID("run-1", block); !stderrors.Is(err, domainerrors.ErrBackfillBlockOverflow) {
			t.Fatalf("Block %d = %v, will ErrBackfillBlockOverflow", block, err)
		}
	}
	if _, err := BackfillTransactionID("", 1); !stderrors.Is(err, domainerrors.ErrEmptyIdentifier) {
		t.Fatalf("leere Run-Kennung = %v, will ErrEmptyIdentifier", err)
	}
}

// TestChangeIDFor trägt die Bildungsregel der Change-Kennung
// `<Transaktions-ID>-<Sequenz>` (`ADR-0111` Teilfrage 6).
func TestChangeIDFor(t *testing.T) {
	if got := ChangeIDFor("0bf-run-1-00000001", 3); got != "0bf-run-1-00000001-3" {
		t.Fatalf("ChangeIDFor = %q", got)
	}
	if a, b := ChangeIDFor("tx", 1), ChangeIDFor("tx", 11); a == b {
		t.Fatalf("verschiedene Sequenzen tragen dieselbe Kennung: %q %q", a, b)
	}
}
