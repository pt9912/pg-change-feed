package model

import (
	"fmt"

	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

// BackfillRunID identifiziert einen Backfill-Run: die Kennung des Antrags,
// der ihn annimmt (`SPEC-029`, Spalte `run_id`).
type BackfillRunID string

// BackfillRunStatus trägt den Zustand eines Runs; die Menge ist
// geschlossen (`SPEC-029`, Spalte `status`).
type BackfillRunStatus string

const (
	BackfillRunQueued      BackfillRunStatus = "queued"
	BackfillRunRunning     BackfillRunStatus = "running"
	BackfillRunCompleted   BackfillRunStatus = "completed"
	BackfillRunFailed      BackfillRunStatus = "failed"
	BackfillRunInterrupted BackfillRunStatus = "interrupted"
)

// RowEstimate trägt die vom Katalog **geschätzte** Zeilenzahl einer
// Tabelle: eine Zahl oder „unbekannt" (`SPEC-029`, `estimated_rows`). Der
// Nullwert ist „unbekannt", nie 0 — die Zahl 0 einer analysierten, leeren
// Tabelle ist eine bekannte Schätzung (`NewRowEstimate(0)`). Die Zahl ist
// eine Orientierung, keine Grenze.
type RowEstimate struct {
	rows  int64
	known bool
}

// UnknownRowEstimate trägt „unbekannt": der Katalog führt keine Schätzung.
func UnknownRowEstimate() RowEstimate {
	return RowEstimate{}
}

// NewRowEstimate trägt eine bekannte, geschätzte Zeilenzahl; eine negative
// Zahl ist keine Schätzung.
func NewRowEstimate(rows int64) (RowEstimate, error) {
	if rows < 0 {
		return RowEstimate{}, domainerrors.ErrNegativeRowCount
	}
	return RowEstimate{rows: rows, known: true}, nil
}

// Rows trägt die geschätzte Zahl und meldet über `known`, ob sie bekannt
// ist; ist sie unbekannt, ist `rows` ohne Aussage.
func (e RowEstimate) Rows() (rows int64, known bool) {
	return e.rows, e.known
}

// BackfillRun trägt den Zustand eines Backfill-Runs (`SPEC-029`,
// `LH-FA-CAP-009`): eine Zeile je Run. Ein Run entsteht `queued`
// (`NewQueuedBackfillRun`) und wechselt nur über die Methoden dieses Typs;
// jede liefert einen neuen Wert und lässt den Empfänger unverändert.
//
// `StartedAt` ist der Beginn der Kopierdauer (Übergang nach `running`), die
// Wartezeit in `queued` zählt nicht. `SnapshotPosition` ist die Position `X`
// des Runs, der Nullwert bis zur Anlage des Snapshots. `ErrorMessage` trägt
// bei `failed` die Fehlerklasse (`SPEC-008`) vor dem Text, sonst ist er
// leer. Die beiden Warn-Kennzeichnungen sind `false`, bis die Auswertung des
// Use Cases sie setzt (`SPEC-029`, `ADR-0113`); sie ändern weder Status noch
// Ablauf des Runs.
type BackfillRun struct {
	ID                BackfillRunID
	Source            SourceID
	Schema            string
	Table             string
	Status            BackfillRunStatus
	RequestedAt       TimePoint
	StartedAt         TimePoint
	FinishedAt        TimePoint
	SnapshotPosition  SourcePosition
	RowsCopied        int64
	EstimatedRows     RowEstimate
	ErrorMessage      string
	WarnEstimatedSize bool
	WarnDuration      bool
}

// NewQueuedBackfillRun legt einen Run im Zustand `queued` an und verlangt
// nichtleere Kennungen; die Schätzung ist zunächst „unbekannt"
// (`WithEstimatedRows`).
func NewQueuedBackfillRun(id BackfillRunID, source SourceID, schema, table string, requestedAt TimePoint) (BackfillRun, error) {
	if id == "" || source == "" || schema == "" || table == "" {
		return BackfillRun{}, domainerrors.ErrEmptyIdentifier
	}
	return BackfillRun{
		ID:          id,
		Source:      source,
		Schema:      schema,
		Table:       table,
		Status:      BackfillRunQueued,
		RequestedAt: requestedAt,
	}, nil
}

// WithEstimatedRows liefert den Run mit der übergebenen Schätzung.
func (r BackfillRun) WithEstimatedRows(estimate RowEstimate) BackfillRun {
	r.EstimatedRows = estimate
	return r
}

// QualifiedName trägt den schema-qualifizierten Tabellennamen in der Form
// `schema.table`.
func (r BackfillRun) QualifiedName() string {
	return r.Schema + "." + r.Table
}

// IsActive meldet, dass der Run `queued` oder `running` ist — die Zustände,
// in denen ein zweiter Run derselben Tabelle ausgeschlossen ist
// (`ADR-0111` Teilfrage 4).
func (r BackfillRun) IsActive() bool {
	return r.Status == BackfillRunQueued || r.Status == BackfillRunRunning
}

// Start wechselt `queued` → `running` und setzt den Beginn der
// Kopierdauer.
func (r BackfillRun) Start(at TimePoint) (BackfillRun, error) {
	if r.Status != BackfillRunQueued {
		return r, domainerrors.ErrInvalidBackfillTransition
	}
	r.Status = BackfillRunRunning
	r.StartedAt = at
	return r, nil
}

// RecordProgress trägt die Snapshot-Position `X` und den Fortschrittszähler
// eines laufenden Runs fort; der Zähler wächst nur, die Position gehört zur
// Quelle des Runs.
func (r BackfillRun) RecordProgress(position SourcePosition, rowsCopied int64) (BackfillRun, error) {
	if r.Status != BackfillRunRunning {
		return r, domainerrors.ErrInvalidBackfillTransition
	}
	if rowsCopied < r.RowsCopied {
		return r, domainerrors.ErrBackfillProgressRegression
	}
	if position.SourceID != r.Source {
		return r, domainerrors.ErrSourceMismatch
	}
	r.SnapshotPosition = position
	r.RowsCopied = rowsCopied
	return r, nil
}

// Complete wechselt `running` → `completed` mit der Zahl der geschriebenen
// Backfill-Changes.
func (r BackfillRun) Complete(at TimePoint, rowsCopied int64) (BackfillRun, error) {
	if r.Status != BackfillRunRunning {
		return r, domainerrors.ErrInvalidBackfillTransition
	}
	if rowsCopied < r.RowsCopied {
		return r, domainerrors.ErrBackfillProgressRegression
	}
	r.Status = BackfillRunCompleted
	r.FinishedAt = at
	r.RowsCopied = rowsCopied
	return r, nil
}

// Fail wechselt `queued` oder `running` → `failed`; der Fehlertext trägt die
// Fehlerklasse (`SPEC-008`) vor der Ursache: `<Klasse>: <Ursache>`.
func (r BackfillRun) Fail(at TimePoint, class ErrorClass, cause string) (BackfillRun, error) {
	if !r.IsActive() {
		return r, domainerrors.ErrInvalidBackfillTransition
	}
	if _, err := NewErrorClass(string(class)); err != nil {
		return r, err
	}
	r.Status = BackfillRunFailed
	r.FinishedAt = at
	r.ErrorMessage = fmt.Sprintf("%s: %s", class, cause)
	return r, nil
}

// Interrupt wechselt `running` → `interrupted`: der Run endete, weil sein
// Prozess oder sein Kontext endete; ein `queued`-Run bleibt `queued`.
func (r BackfillRun) Interrupt(at TimePoint) (BackfillRun, error) {
	if r.Status != BackfillRunRunning {
		return r, domainerrors.ErrInvalidBackfillTransition
	}
	r.Status = BackfillRunInterrupted
	r.FinishedAt = at
	return r, nil
}

// backfillTransactionPrefix trägt die synthetischen Transaktionen eines
// Backfills; die führende `0` sortiert sie vor jede WAL-Transaktions-Kennung
// (Ziffernfolgen ohne führende Null) derselben Position (`ADR-0111`
// Teilfrage 6).
const backfillTransactionPrefix = "0bf-"

// maxBackfillBlock ist die größte achtstellige Blocknummer.
const maxBackfillBlock = 99_999_999

// BackfillTransactionID bildet die Kennung der synthetischen Transaktion
// eines Blocks: `0bf-<Run-Kennung>-<Blocknummer>` mit achtstelliger,
// null-aufgefüllter Blocknummer (`ADR-0111` Teilfrage 6). Die Blöcke eines
// Runs zählen ab 1.
func BackfillTransactionID(run BackfillRunID, block int) (TransactionID, error) {
	if run == "" {
		return "", domainerrors.ErrEmptyIdentifier
	}
	if block < 1 || block > maxBackfillBlock {
		return "", domainerrors.ErrBackfillBlockOverflow
	}
	return TransactionID(fmt.Sprintf("%s%s-%08d", backfillTransactionPrefix, run, block)), nil
}
