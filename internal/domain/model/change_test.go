package model

import (
	"testing"

	stderrors "errors"

	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

// changeArgs trägt die Argumente eines gültigen Changes; je Test weicht
// genau ein Feld ab.
type changeArgs struct {
	id      ChangeID
	tx      TransactionID
	table   SourceTableID
	seq     int64
	op      Operation
	oldData []byte
	newData []byte
	sv      SchemaVersionID
}

func validChangeArgs() changeArgs {
	return changeArgs{
		id:      "ch-1",
		tx:      "tx-1",
		table:   "tbl-1",
		seq:     1,
		op:      OperationUpdate,
		oldData: []byte(`{"a":1}`),
		newData: []byte(`{"a":2}`),
		sv:      "sv-1",
	}
}

func buildChange(t *testing.T, args changeArgs) Change {
	t.Helper()
	change, err := NewChange(args.id, args.tx, args.table, args.seq, args.op, args.oldData, args.newData, args.sv)
	if err != nil {
		t.Fatalf("Change-Konstruktor: %v", err)
	}
	return change
}

// LH-FA-DAT-001, Happy Path: zwei erfasste Changes sind über ihre
// Identifikatoren unterscheidbar.
func TestLHFADAT001ChangesAreDistinguishable(t *testing.T) {
	first := buildChange(t, validChangeArgs())
	second := buildChange(t, changeArgs{id: "ch-2", tx: "tx-1", table: "tbl-1", seq: 2, op: OperationUpdate, oldData: []byte(`{"a":1}`), newData: []byte(`{"a":2}`), sv: "sv-1"})
	if first.ID == second.ID {
		t.Fatalf("zwei erfasste Changes tragen dieselbe Kennung: %q", first.ID)
	}
}

// LH-FA-DAT-001, Boundary: derselbe Datensatz, zweimal identisch geändert,
// ergibt zwei unterscheidbare Changes — identische Bilder, verschiedene
// Kennung und Sequenz.
func TestLHFADAT001SameRowChangedTwiceYieldsDistinctChanges(t *testing.T) {
	first := buildChange(t, validChangeArgs())
	second := buildChange(t, changeArgs{id: "ch-2", tx: "tx-1", table: "tbl-1", seq: 2, op: OperationUpdate, oldData: []byte(`{"a":1}`), newData: []byte(`{"a":2}`), sv: "sv-1"})
	if string(first.NewImage) != string(second.NewImage) || string(first.OldImage) != string(second.OldImage) {
		t.Fatalf("Boundary verletzt: die Änderung ist identisch, Bilder %s/%s vs %s/%s", first.OldImage, first.NewImage, second.OldImage, second.NewImage)
	}
	if first.ID == second.ID || first.Sequence == second.Sequence {
		t.Fatalf("zwei identische Änderungen sind nicht unterscheidbar: Kennung %q und Sequenz %d doppelt", first.ID, first.Sequence)
	}
}

// Die Konstruktoren erzwingen die Change-Invarianten (`ADR-0029`, Regel 6
// und 7): illegale Zustände scheitern am Konstruktor.
func TestNewChangeRejectsInvariantViolations(t *testing.T) {
	cases := []struct {
		name      string
		mutate    func(args *changeArgs)
		wantError error
	}{
		{"leere Change-Kennung", func(a *changeArgs) { a.id = "" }, domainerrors.ErrEmptyIdentifier},
		{"leere Transaktions-Kennung", func(a *changeArgs) { a.tx = "" }, domainerrors.ErrEmptyIdentifier},
		{"leere Tabellen-Kennung", func(a *changeArgs) { a.table = "" }, domainerrors.ErrEmptyIdentifier},
		{"leere Schema-Version", func(a *changeArgs) { a.sv = "" }, domainerrors.ErrEmptyIdentifier},
		{"Sequenz 0", func(a *changeArgs) { a.seq = 0 }, domainerrors.ErrNonPositiveSequence},
		{"negative Sequenz", func(a *changeArgs) { a.seq = -3 }, domainerrors.ErrNonPositiveSequence},
		{"unbekannte Operation", func(a *changeArgs) { a.op = Operation("TRUNCATE") }, domainerrors.ErrInvalidOperation},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			args := validChangeArgs()
			tc.mutate(&args)
			_, err := NewChange(args.id, args.tx, args.table, args.seq, args.op, args.oldData, args.newData, args.sv)
			if !stderrors.Is(err, tc.wantError) {
				t.Fatalf("Fehler = %v, wollen %v", err, tc.wantError)
			}
		})
	}
}
