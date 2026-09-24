package mapper

import (
	stderrors "errors"
	"math"
	"time"

	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// ErrUnknownBackfillStatus: der Status einer `cdc.backfill_run`-Zeile liegt
// außerhalb der geschlossenen Menge (`SPEC-029`) — die CHECK-Kante der
// Tabelle hält sie, der Lesepfad prüft sie zusätzlich, statt einen
// unbekannten Wert als Domänen-Zustand zu führen.
var ErrUnknownBackfillStatus = stderrors.New("Status der Backfill-Run-Zeile liegt außerhalb der geschlossenen Menge")

// BackfillRunRow trägt eine Zeile aus `cdc.backfill_run` (`SPEC-029`).
// Die nullbaren Spalten sind Zeiger; `ErrorMessage` liest als leerer Text,
// wenn die Spalte NULL trägt (die Abfrage normalisiert über `COALESCE`).
type BackfillRunRow struct {
	RunID             string
	SourceID          string
	Schema            string
	Table             string
	Status            string
	RequestedAt       time.Time
	StartedAt         *time.Time
	FinishedAt        *time.Time
	SnapshotPosition  *int64
	RowsCopied        int64
	EstimatedRows     *int64
	ErrorMessage      string
	WarnEstimatedSize bool
	WarnDuration      bool
}

// ToBackfillRun trägt den Run aus einer `cdc.backfill_run`-Zeile. Die
// Kennungen laufen durch den Domänen-Konstruktor (`ADR-0029`), der Status
// durch die geschlossene Menge; eine NULL-Schätzung liest als „unbekannt“
// (`model.UnknownRowEstimate`), nie als 0 (`SPEC-029`), eine NULL-Position
// als Nullwert.
func ToBackfillRun(row BackfillRunRow) (model.BackfillRun, error) {
	run, err := model.NewQueuedBackfillRun(
		model.BackfillRunID(row.RunID), model.SourceID(row.SourceID), row.Schema, row.Table,
		model.NewTimePoint(row.RequestedAt.UnixNano()),
	)
	if err != nil {
		return model.BackfillRun{}, err
	}
	status := model.BackfillRunStatus(row.Status)
	switch status {
	case model.BackfillRunQueued, model.BackfillRunRunning, model.BackfillRunCompleted,
		model.BackfillRunFailed, model.BackfillRunInterrupted:
		run.Status = status
	default:
		return model.BackfillRun{}, ErrUnknownBackfillStatus
	}
	if row.RowsCopied < 0 {
		return model.BackfillRun{}, domainerrors.ErrNegativeRowCount
	}
	run.RowsCopied = row.RowsCopied
	run.ErrorMessage = row.ErrorMessage
	run.WarnEstimatedSize = row.WarnEstimatedSize
	run.WarnDuration = row.WarnDuration
	if row.StartedAt != nil {
		run.StartedAt = model.NewTimePoint(row.StartedAt.UnixNano())
	}
	if row.FinishedAt != nil {
		run.FinishedAt = model.NewTimePoint(row.FinishedAt.UnixNano())
	}
	if row.SnapshotPosition != nil {
		position, err := ToPosition(row.SourceID, *row.SnapshotPosition)
		if err != nil {
			return model.BackfillRun{}, err
		}
		run.SnapshotPosition = position
	}
	if row.EstimatedRows != nil {
		estimate, err := model.NewRowEstimate(*row.EstimatedRows)
		if err != nil {
			return model.BackfillRun{}, err
		}
		run.EstimatedRows = estimate
	}
	return run, nil
}

// TimeArgument trägt einen Zeitpunkt als `timestamptz`-Argument; der Nullwert
// (`TimePoint.Unset`) geht als NULL.
func TimeArgument(at model.TimePoint) any {
	if at.Unset() {
		return nil
	}
	return time.Unix(0, at.UnixNanos).UTC()
}

// EstimateArgument trägt die geschätzte Zeilenzahl als `bigint`-Argument:
// „unbekannt“ geht als NULL, nie als 0 (`SPEC-029`).
func EstimateArgument(estimate model.RowEstimate) any {
	rows, known := estimate.Rows()
	if !known {
		return nil
	}
	return rows
}

// PositionArgument trägt die Snapshot-Position `X` als `bigint`-Argument:
// der Nullwert (Offset 0, der Snapshot ist noch nicht angelegt) geht als
// NULL. Eine Position außerhalb des bigint-Bereichs ist eine Grenze der
// PostgreSQL-Abbildung (`SPEC-003`).
func PositionArgument(position model.SourcePosition) (any, error) {
	if position.Offset == 0 {
		return nil, nil
	}
	if position.Offset > math.MaxInt64 {
		return nil, ErrPositionOutOfRange
	}
	return int64(position.Offset), nil
}
