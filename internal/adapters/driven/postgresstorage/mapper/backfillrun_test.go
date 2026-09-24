package mapper_test

import (
	stderrors "errors"
	"math"
	"testing"
	"time"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage/mapper"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// backfillRow trägt eine vollständige, gültige Run-Zeile; die Tests ändern
// je ein Feld.
func backfillRow() mapper.BackfillRunRow {
	requested := time.Unix(100, 5).UTC()
	started := time.Unix(200, 0).UTC()
	finished := time.Unix(300, 0).UTC()
	position := int64(7000)
	estimate := int64(1234)
	return mapper.BackfillRunRow{
		RunID:             "req-1",
		SourceID:          "src-1",
		Schema:            "public",
		Table:             "orders",
		Status:            "failed",
		RequestedAt:       requested,
		StartedAt:         &started,
		FinishedAt:        &finished,
		SnapshotPosition:  &position,
		RowsCopied:        42,
		EstimatedRows:     &estimate,
		ErrorMessage:      "storage: Verbindung verloren",
		WarnEstimatedSize: true,
		WarnDuration:      true,
	}
}

// Die Zeile trägt jedes Feld in den Run; die Zeit- und Positionsspalten
// laufen in ihre Domänen-Form, die Warn-Spalten unverändert.
func TestToBackfillRunCarriesEveryColumn(t *testing.T) {
	run, err := mapper.ToBackfillRun(backfillRow())
	if err != nil {
		t.Fatalf("ToBackfillRun: %v", err)
	}
	if run.ID != "req-1" || run.Source != "src-1" || run.Schema != "public" || run.Table != "orders" {
		t.Fatalf("Kennungen = %+v", run)
	}
	if run.Status != model.BackfillRunFailed {
		t.Fatalf("Status = %q", run.Status)
	}
	if run.RequestedAt.UnixNanos != time.Unix(100, 5).UnixNano() ||
		run.StartedAt.UnixNanos != time.Unix(200, 0).UnixNano() ||
		run.FinishedAt.UnixNanos != time.Unix(300, 0).UnixNano() {
		t.Fatalf("Zeiten = %+v %+v %+v", run.RequestedAt, run.StartedAt, run.FinishedAt)
	}
	if run.SnapshotPosition != (model.SourcePosition{SourceID: "src-1", Offset: 7000}) {
		t.Fatalf("Position = %+v", run.SnapshotPosition)
	}
	if run.RowsCopied != 42 || run.ErrorMessage != "storage: Verbindung verloren" || !run.WarnEstimatedSize || !run.WarnDuration {
		t.Fatalf("Zähler/Text/Warnungen = %d %q %v %v", run.RowsCopied, run.ErrorMessage, run.WarnEstimatedSize, run.WarnDuration)
	}
	if rows, known := run.EstimatedRows.Rows(); !known || rows != 1234 {
		t.Fatalf("Schätzung = %d, bekannt %v", rows, known)
	}
}

// NULL heißt „unbekannt“, nie 0 (`SPEC-029`): die NULL-Schätzung liest als
// unbekannt, die Schätzung 0 als bekannte Zahl 0. Die NULL-Zeiten und die
// NULL-Position lesen als Nullwert. Rot färbende Mutation: `ToBackfillRun` liest die
// NULL-Schätzung als `NewRowEstimate(0)`.
func TestToBackfillRunReadsNullAsUnknown(t *testing.T) {
	row := backfillRow()
	row.StartedAt, row.FinishedAt, row.SnapshotPosition, row.EstimatedRows = nil, nil, nil, nil
	row.Status, row.ErrorMessage = "queued", ""

	run, err := mapper.ToBackfillRun(row)
	if err != nil {
		t.Fatalf("ToBackfillRun: %v", err)
	}
	if _, known := run.EstimatedRows.Rows(); known {
		t.Fatal("NULL-Schätzung liest als bekannt")
	}
	if !run.StartedAt.Unset() || !run.FinishedAt.Unset() || run.SnapshotPosition != (model.SourcePosition{}) {
		t.Fatalf("NULL-Spalten lesen nicht als Nullwert: %+v", run)
	}

	zero := int64(0)
	row.EstimatedRows = &zero
	run, err = mapper.ToBackfillRun(row)
	if err != nil {
		t.Fatalf("ToBackfillRun: %v", err)
	}
	if rows, known := run.EstimatedRows.Rows(); !known || rows != 0 {
		t.Fatalf("Schätzung 0 = %d, bekannt %v", rows, known)
	}
}

// Jeder Status der geschlossenen Menge liest, jeder andere Wert endet
// sichtbar.
func TestToBackfillRunChecksTheClosedStatusSet(t *testing.T) {
	for _, status := range []string{"queued", "running", "completed", "failed", "interrupted"} {
		row := backfillRow()
		row.Status = status
		run, err := mapper.ToBackfillRun(row)
		if err != nil || string(run.Status) != status {
			t.Fatalf("Status %s: %+v, %v", status, run.Status, err)
		}
	}
	row := backfillRow()
	row.Status = "paused"
	if _, err := mapper.ToBackfillRun(row); !stderrors.Is(err, mapper.ErrUnknownBackfillStatus) {
		t.Fatalf("Status außerhalb der Menge: %v", err)
	}
}

// Die Zeile läuft durch die Domänen-Invarianten: leere Kennung, negative
// Zahlen und eine Position unter 1 enden als Domänen-Fehler.
func TestToBackfillRunRejectsRowsOutsideTheInvariants(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*mapper.BackfillRunRow)
		want   error
	}{
		{"leere Run-Kennung", func(r *mapper.BackfillRunRow) { r.RunID = "" }, domainerrors.ErrEmptyIdentifier},
		{"negativer Zähler", func(r *mapper.BackfillRunRow) { r.RowsCopied = -1 }, domainerrors.ErrNegativeRowCount},
		{"negative Schätzung", func(r *mapper.BackfillRunRow) { v := int64(-1); r.EstimatedRows = &v }, domainerrors.ErrNegativeRowCount},
		{"Position 0", func(r *mapper.BackfillRunRow) { v := int64(0); r.SnapshotPosition = &v }, domainerrors.ErrInvalidPosition},
	}
	for _, c := range cases {
		row := backfillRow()
		c.mutate(&row)
		if _, err := mapper.ToBackfillRun(row); !stderrors.Is(err, c.want) {
			t.Fatalf("%s: Fehler = %v, wollen %v", c.name, err, c.want)
		}
	}
}

// Die Schreib-Argumente: unbekannte Schätzung und Nullwerte gehen als NULL,
// bekannte Werte in ihrer Spalten-Form.
func TestBackfillArgumentsCarryNullAndValues(t *testing.T) {
	if got := mapper.EstimateArgument(model.UnknownRowEstimate()); got != nil {
		t.Fatalf("unbekannte Schätzung = %v, wollen nil", got)
	}
	known, err := model.NewRowEstimate(0)
	if err != nil {
		t.Fatalf("NewRowEstimate: %v", err)
	}
	if got := mapper.EstimateArgument(known); got != int64(0) {
		t.Fatalf("Schätzung 0 = %v, wollen 0 (bekannt)", got)
	}

	if got := mapper.TimeArgument(model.TimePoint{}); got != nil {
		t.Fatalf("Nullzeit = %v, wollen nil", got)
	}
	want := time.Unix(0, 1_700_000_000_123_456_789).UTC()
	if got := mapper.TimeArgument(model.NewTimePoint(want.UnixNano())); got != want {
		t.Fatalf("Zeit = %v, wollen %v", got, want)
	}

	if got, err := mapper.PositionArgument(model.SourcePosition{}); got != nil || err != nil {
		t.Fatalf("Nullposition = %v, %v", got, err)
	}
	if got, err := mapper.PositionArgument(model.SourcePosition{SourceID: "s", Offset: 99}); got != int64(99) || err != nil {
		t.Fatalf("Position = %v, %v", got, err)
	}
	if _, err := mapper.PositionArgument(model.SourcePosition{SourceID: "s", Offset: math.MaxInt64 + 1}); !stderrors.Is(err, mapper.ErrPositionOutOfRange) {
		t.Fatalf("Position jenseits bigint: %v", err)
	}
}
