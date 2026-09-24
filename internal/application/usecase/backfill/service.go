// Package backfill trägt den BackfillTable Use Case (`ARC-002`,
// `ADR-0028`): die Annahme eines Backfill-Antrags und die Ausführung des
// angenommenen Runs (`LH-FA-CAP-009`, `ADR-0111`, `ADR-0113`).
package backfill

import (
	"context"
	"errors"
	"fmt"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Transport-Typen des Use Cases (`ADR-0042`): ihre Definition liegt am
// Inbound-Port, die Adressen hier sind Aliase desselben Typs.
type (
	BackfillRequestCommand = inbound.BackfillRequestCommand
	BackfillRequestResult  = inbound.BackfillRequestResult
	BackfillExecuteCommand = inbound.BackfillExecuteCommand
	BackfillExecuteResult  = inbound.BackfillExecuteResult
)

// Ports bündelt die acht obligatorischen Outbound-Ports des Use Cases; alle
// sind Fähigkeits-Ports (`ADR-0034`): Bindung und Publication
// (`TableActivationPort`), Ausschlussstand (`ColumnExclusionPort`),
// Schema-Version der Tabelle (`SchemaStorePort`), Snapshot
// (`TableSnapshotPort`), Annahme, Run-Zustand und Schreiber sowie die Uhr
// (`ClockPort`, `ADR-0040`).
type Ports struct {
	Activation outbound.TableActivationPort
	Exclusion  outbound.ColumnExclusionPort
	Schemas    outbound.SchemaStorePort
	Snapshot   outbound.TableSnapshotPort
	Admission  outbound.BackfillAdmissionPort
	Runs       outbound.BackfillRunPort
	Writer     outbound.BackfillWriterPort
	Clock      outbound.ClockPort
}

// BackfillTableService implementiert `inbound.BackfillTableUseCase`.
//
// `Execute` hält höchstens einen Block im Speicher: der Schreiber erhält
// Block `n`, bevor der Leser Block `n+1` liefert (`LH-FA-CAP-006.a`). Die
// Grenze zählt Zeilen, nicht Bytes — der Speicherbedarf eines Blocks ist
// seine Zeilenzahl mal die Zeilenbreite (`TableSnapshot.NextBlock`).
//
// Fehler eines Runs sind run-lokal: sie enden den Run, sie berühren weder
// den Heartbeat-Fehlerzustand noch den Capture-Pfad (`SPEC-008`).
type BackfillTableService struct {
	ports  Ports
	notify outbound.ChangeNotificationPort
	log    outbound.LogPort
}

// Option konfiguriert die optionalen Ports des Use Cases.
type Option func(*BackfillTableService)

// WithChangeNotification injiziert den optionalen `ChangeNotificationPort`
// (`ADR-0055`): ungesetzt unterbleibt das Wecksignal nach dem Commit.
func WithChangeNotification(notify outbound.ChangeNotificationPort) Option {
	return func(s *BackfillTableService) { s.notify = notify }
}

// WithLog injiziert den `LogPort` (`ADR-0024`) für die Fehlerpfade, die
// den Run nicht beenden (Wecksignal, Rollback, Schließen des Snapshots).
// Ungesetzt bleibt die Protokollierung beim No-Op (`outbound.NoopLog`).
func WithLog(log outbound.LogPort) Option {
	return func(s *BackfillTableService) { s.log = log }
}

// NewBackfillTableService verdrahtet den Use Case mit seinen
// obligatorischen Ports.
func NewBackfillTableService(ports Ports, opts ...Option) *BackfillTableService {
	s := &BackfillTableService{ports: ports, log: outbound.NoopLog}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

var _ inbound.BackfillTableUseCase = (*BackfillTableService)(nil)

// Request nimmt einen Backfill-Antrag an (`ADR-0113` Festlegung 1): die
// Vorbedingungen — Tabelle mit laufender Bindung und Mitglied der
// Publication —, die geschätzte Zeilenzahl, dann als **letzter** Schritt
// `Admit`. Vor `Admit` entsteht weder eine Run-Zeile noch ein Snapshot; die
// Schätzung ist ein Katalog-Lesezugriff. Die Prüfung „kein aktiver Run"
// liegt in `Admit`.
func (s *BackfillTableService) Request(ctx context.Context, command BackfillRequestCommand) (BackfillRequestResult, error) {
	run, err := model.NewQueuedBackfillRun(model.BackfillRunID(command.RequestID), command.Source, command.Schema, command.Table, s.ports.Clock.Now())
	if err != nil {
		return BackfillRequestResult{}, err
	}
	if _, err := s.publishedTable(ctx, command.Source, command.Schema, command.Table, command.Publication); err != nil {
		return BackfillRequestResult{}, err
	}
	rows, known, err := s.ports.Snapshot.EstimatedRows(ctx, command.Schema, command.Table)
	if err != nil {
		return BackfillRequestResult{}, err
	}
	estimate := model.UnknownRowEstimate()
	if known {
		if estimate, err = model.NewRowEstimate(rows); err != nil {
			return BackfillRequestResult{}, err
		}
	}
	run = run.WithEstimatedRows(estimate)
	if err := s.ports.Admission.Admit(ctx, command.RequestID, run); err != nil {
		return BackfillRequestResult{}, err
	}
	return BackfillRequestResult{Run: run}, nil
}

// Execute führt einen angenommenen Run aus. Zuerst prüft er Bindung und
// Publication-Mitgliedschaft erneut: eine Abweichung endet den Run `failed`
// (Klasse `configuration`) ohne Snapshot und ohne Kopie. Danach kopiert er
// blockweise in **eine** Schreibtransaktion (`copyBlocks`) und committet sie
// einmal. Jeder Fehler rollt die Transaktion zurück und endet den Run
// (`conclude`); nach dem Commit sendet er ein Wecksignal je Tabelle.
func (s *BackfillTableService) Execute(ctx context.Context, command BackfillExecuteCommand) (BackfillExecuteResult, error) {
	run := command.Run
	if run.Status != model.BackfillRunQueued {
		return BackfillExecuteResult{Run: run}, domainerrors.ErrInvalidBackfillTransition
	}
	if err := ctx.Err(); err != nil {
		return BackfillExecuteResult{Run: run}, err
	}

	table, err := s.publishedTable(ctx, run.Source, run.Schema, run.Table, command.Publication)
	if err != nil {
		return s.conclude(ctx, run, err)
	}
	version, err := s.currentVersion(ctx, table)
	if err != nil {
		return s.conclude(ctx, run, err)
	}

	started, err := run.Start(s.ports.Clock.Now())
	if err != nil {
		return s.conclude(ctx, run, err)
	}
	if err := s.ports.Runs.MarkRunning(ctx, started); err != nil {
		return s.conclude(ctx, run, err)
	}

	finished, err := s.copyBlocks(ctx, started, table, version)
	if err != nil {
		return s.conclude(ctx, finished, err)
	}
	if finished.RowsCopied > 0 {
		s.wake(ctx, finished)
	}
	return BackfillExecuteResult{Run: finished}, nil
}

// copyBlocks liest den Snapshot blockweise und schreibt jeden Block an den
// Schreiber, bevor der nächste gelesen wird. Es liefert den Run im Endzustand
// `completed` oder, mit Fehler, den zuletzt fortgeschriebenen Run. Snapshot
// und Schreibtransaktion sind auf jedem Pfad geschlossen bzw. zurückgerollt,
// bevor die Funktion zurückkehrt; eine leere Tabelle öffnet keine
// Schreibtransaktion.
//
// Alle Blöcke tragen die Position `X` des Snapshots und den
// Snapshot-Zeitpunkt (die Uhr wird nach dem Öffnen des Snapshots gelesen,
// `ADR-0111` Teilfrage 7). Jeder Block liest den Ausschlussstand neu und
// baut das Bild mit ihm. Fail-closed (`ADR-0111` Teilfrage 4,
// `LH-QA-SEC-004`): jeder weitere Block und der Zustand unmittelbar vor dem
// Commit tragen denselben Ausschlussstand wie der erste Block, und die
// Bindung besteht unter derselben Tabellen-Kennung.
func (s *BackfillTableService) copyBlocks(ctx context.Context, run model.BackfillRun, table model.SourceTable, version model.SchemaVersion) (model.BackfillRun, error) {
	snapshot, err := s.ports.Snapshot.OpenSnapshot(ctx, string(run.ID), run.Schema, run.Table)
	if err != nil {
		return run, err
	}
	defer s.closeSnapshot(ctx, snapshot)

	position, err := model.NewSourcePosition(run.Source, snapshot.Offset())
	if err != nil {
		return run, err
	}
	builder := blockBuilder{
		run:      run,
		table:    table,
		version:  version,
		position: position,
		at:       s.ports.Clock.Now(),
		columns:  snapshot.Columns(),
	}
	if run, err = s.progress(ctx, run, position, 0); err != nil {
		return run, err
	}

	var (
		writer    outbound.BackfillTransaction
		committed bool
		baseline  []string
		copied    int64
	)
	defer func() {
		if !committed {
			s.rollback(ctx, writer)
		}
	}()

	for blockNumber := 1; ; blockNumber++ {
		rows, err := snapshot.NextBlock(ctx)
		if err != nil {
			return run, err
		}
		if len(rows) == 0 {
			break
		}
		excluded, err := s.excludedColumns(ctx, run)
		if err != nil {
			return run, err
		}
		if writer == nil {
			baseline = excluded
			if writer, err = s.ports.Writer.Begin(ctx, run); err != nil {
				return run, err
			}
		} else if !sameNames(baseline, excluded) {
			return run, domainerrors.ErrExclusionStateChanged
		}

		block, err := builder.build(blockNumber, rows, excluded)
		if err != nil {
			return run, err
		}
		if err := writer.AppendBlock(ctx, block); err != nil {
			return run, err
		}
		copied += int64(len(rows))
		if run, err = s.progress(ctx, run, position, copied); err != nil {
			return run, err
		}
	}

	if writer == nil {
		completed, err := run.Complete(s.ports.Clock.Now(), 0)
		if err != nil {
			return run, err
		}
		if err := s.ports.Runs.Finish(ctx, completed); err != nil {
			return run, err
		}
		return completed, nil
	}

	if err := s.stillBound(ctx, table); err != nil {
		return run, err
	}
	excluded, err := s.excludedColumns(ctx, run)
	if err != nil {
		return run, err
	}
	if !sameNames(baseline, excluded) {
		return run, domainerrors.ErrExclusionStateChanged
	}
	completed, err := run.Complete(s.ports.Clock.Now(), copied)
	if err != nil {
		return run, err
	}
	if err := writer.Commit(ctx, completed); err != nil {
		return run, err
	}
	committed = true
	return completed, nil
}

// conclude beendet den Run, dessen Ausführung mit `cause` endete: `failed`
// mit der Fehlerklasse der Ursache (`classifyError`) — oder `interrupted`,
// wenn der eigene Kontext endete, denn ein Fehler bei beendetem Kontext ist
// dessen Folge. Ein noch `queued` Run bleibt bei beendetem Kontext `queued`
// und wird nach einem Neustart aufgenommen. Der Endzustand wird auf einem vom
// Abbruch gelösten Kontext festgehalten; ohne Festhalten meldet der Aufruf
// einen Fehler, und der Run bleibt bis zum Abgleich beim Prozessstart
// `running`.
func (s *BackfillTableService) conclude(ctx context.Context, run model.BackfillRun, cause error) (BackfillExecuteResult, error) {
	var (
		final model.BackfillRun
		err   error
	)
	switch {
	case ctx.Err() != nil && run.Status == model.BackfillRunQueued:
		return BackfillExecuteResult{Run: run}, cause
	case ctx.Err() != nil:
		final, err = run.Interrupt(s.ports.Clock.Now())
	default:
		final, err = run.Fail(s.ports.Clock.Now(), classifyError(cause), cause.Error())
	}
	if err != nil {
		return BackfillExecuteResult{Run: run}, fmt.Errorf("%w (Ursache des Runs: %v)", err, cause)
	}
	if err := s.ports.Runs.Finish(context.WithoutCancel(ctx), final); err != nil {
		return BackfillExecuteResult{Run: run}, fmt.Errorf("Run-Zustand nicht festgehalten: %w (Ursache des Runs: %v)", err, cause)
	}
	return BackfillExecuteResult{Run: final}, nil
}

// classifyError ordnet die Ursache eines Run-Fehlers einer der Klassen des
// Run-Vertrags zu — `permission`, `configuration`, `storage`, `transient`,
// `replication` (`SPEC-008`, `ADR-0111` Teilfrage 5); ein nicht erkannter
// Fehler bleibt `internal`. Die Klasse `schema` vergibt der Run nicht. Die
// Abbildung gilt dem Run und ist nicht die des Capture-Pfads.
func classifyError(err error) model.ErrorClass {
	switch {
	case errors.Is(err, outbound.ErrSnapshotPermission):
		return model.ErrorClassPermission
	case errors.Is(err, outbound.ErrSnapshotConfiguration),
		errors.Is(err, domainerrors.ErrTableNotActivated),
		errors.Is(err, domainerrors.ErrExclusionStateChanged),
		errors.Is(err, outbound.ErrSchemaVersionUnknown):
		return model.ErrorClassConfiguration
	case errors.Is(err, outbound.ErrSnapshotTransient):
		return model.ErrorClassTransient
	case errors.Is(err, outbound.ErrSnapshotReplication):
		return model.ErrorClassReplication
	case errors.Is(err, outbound.ErrSnapshotStorage),
		errors.Is(err, outbound.ErrBackfillStorage),
		errors.Is(err, outbound.ErrStorage),
		errors.Is(err, outbound.ErrSchemaStoreStorage):
		return model.ErrorClassStorage
	default:
		return model.ErrorClassInternal
	}
}

// publishedTable prüft die Vorbedingungen des Backfills: die Tabelle trägt
// eine Bindung und ist Mitglied der Publication der Quelle. Eine
// Abweichung endet als `domainerrors.ErrTableNotActivated`.
func (s *BackfillTableService) publishedTable(ctx context.Context, source model.SourceID, schema, table, publication string) (model.SourceTable, error) {
	bound, found, err := s.ports.Activation.Registered(ctx, source, schema, table)
	if err != nil {
		return model.SourceTable{}, err
	}
	if !found {
		return model.SourceTable{}, fmt.Errorf("%w: %s.%s trägt keine Bindung", domainerrors.ErrTableNotActivated, schema, table)
	}
	published, err := s.ports.Activation.Published(ctx, publication, schema, table)
	if err != nil {
		return model.SourceTable{}, err
	}
	if !published {
		return model.SourceTable{}, fmt.Errorf("%w: %s.%s ist nicht Mitglied der Publication", domainerrors.ErrTableNotActivated, schema, table)
	}
	return bound, nil
}

// stillBound prüft unmittelbar vor dem Commit, dass die Bindung der Tabelle
// besteht und dieselbe Kennung trägt wie beim Bau der Blöcke.
func (s *BackfillTableService) stillBound(ctx context.Context, table model.SourceTable) error {
	current, found, err := s.ports.Activation.Registered(ctx, table.SourceID, table.Schema, table.Table)
	if err != nil {
		return err
	}
	if !found || current.ID != table.ID {
		return fmt.Errorf("%w: die Bindung von %s besteht nicht mehr", domainerrors.ErrTableNotActivated, table.QualifiedName())
	}
	return nil
}

// currentVersion liest die aktuelle Schema-Version der Tabelle, auf die
// jeder Change des Runs verweist (`ADR-0029`, Regel 7).
func (s *BackfillTableService) currentVersion(ctx context.Context, table model.SourceTable) (model.SchemaVersion, error) {
	version, found, err := s.ports.Schemas.CurrentVersion(ctx, table.ID)
	if err != nil {
		return model.SchemaVersion{}, err
	}
	if !found {
		return model.SchemaVersion{}, fmt.Errorf("%w: %s", outbound.ErrSchemaVersionUnknown, table.QualifiedName())
	}
	return version, nil
}

// excludedColumns liest den dauerhaften Ausschlussstand der Tabelle des Runs
// (`ADR-0065`).
func (s *BackfillTableService) excludedColumns(ctx context.Context, run model.BackfillRun) ([]string, error) {
	byTable, err := s.ports.Exclusion.ExcludedColumns(ctx, run.Source)
	if err != nil {
		return nil, err
	}
	return byTable[run.QualifiedName()], nil
}

// progress schreibt Snapshot-Position und Fortschrittszähler fort. Bei einem
// Fehler des Ports bleibt der Run im zuletzt festgehaltenen Stand.
func (s *BackfillTableService) progress(ctx context.Context, run model.BackfillRun, position model.SourcePosition, rowsCopied int64) (model.BackfillRun, error) {
	updated, err := run.RecordProgress(position, rowsCopied)
	if err != nil {
		return run, err
	}
	if err := s.ports.Runs.RecordProgress(ctx, updated); err != nil {
		return run, err
	}
	return updated, nil
}

// wake sendet nach dem Commit ein Wecksignal je Tabelle (`ADR-0055`, best
// effort): ein Fehler ändert den abgeschlossenen Run nicht.
func (s *BackfillTableService) wake(ctx context.Context, run model.BackfillRun) {
	if s.notify == nil {
		return
	}
	if err := s.notify.Notify(ctx, string(run.Source), run.Schema, run.Table); err != nil {
		s.log.Warn(ctx, "backfill: Wecksignal fehlgeschlagen", "error", err, "run_id", run.ID, "schema", run.Schema, "table", run.Table)
	}
}

// closeSnapshot schließt den Snapshot auf einem vom Abbruch gelösten
// Kontext (Port-Vertrag `OpenSnapshot`: der Aufrufer schließt auf jedem
// Pfad).
func (s *BackfillTableService) closeSnapshot(ctx context.Context, snapshot outbound.TableSnapshot) {
	if err := snapshot.Close(context.WithoutCancel(ctx)); err != nil {
		s.log.Warn(ctx, "backfill: Schließen des Snapshots fehlgeschlagen", "error", err)
	}
}

// rollback verwirft die offene Schreibtransaktion; ohne geöffnete
// Transaktion (`nil`) ist nichts zu tun.
func (s *BackfillTableService) rollback(ctx context.Context, writer outbound.BackfillTransaction) {
	if writer == nil {
		return
	}
	if err := writer.Rollback(context.WithoutCancel(ctx)); err != nil {
		s.log.Warn(ctx, "backfill: Rollback der Schreibtransaktion fehlgeschlagen", "error", err)
	}
}

// blockBuilder baut die synthetische Transaktion eines Blocks. Er ist die
// eine Stelle des Use Cases, an der ein Row Image entsteht: über
// `model.BuildRowImage`, dieselbe Funktion wie der WAL-Pfad.
type blockBuilder struct {
	run      model.BackfillRun
	table    model.SourceTable
	version  model.SchemaVersion
	position model.SourcePosition
	at       model.TimePoint
	columns  []string
}

// build legt die committed Transaktion `0bf-<Run-Kennung>-<Blocknummer>` an:
// je Zeile ein `INSERT`-Change der Herkunft `backfill` ohne `old_data`, die
// Sequenz zählt ab 1 (`ADR-0111` Teilfrage 6).
func (b blockBuilder) build(blockNumber int, rows [][]*string, excluded []string) (*model.ChangeTransaction, error) {
	id, err := model.BackfillTransactionID(b.run.ID, blockNumber)
	if err != nil {
		return nil, err
	}
	transaction, err := model.NewOpenTransaction(id, b.run.Source)
	if err != nil {
		return nil, err
	}
	for i, row := range rows {
		sequence := int64(i + 1)
		image, err := model.BuildRowImage(b.columns, row, excluded)
		if err != nil {
			return nil, err
		}
		change, err := model.NewChange(model.ChangeIDFor(id, sequence), id, b.table.ID, sequence, model.OperationInsert, nil, image, b.version.ID)
		if err != nil {
			return nil, err
		}
		change.Schema = b.run.Schema
		change.Table = b.run.Table
		if change, err = change.WithOrigin(model.ChangeOriginBackfill); err != nil {
			return nil, err
		}
		if err := transaction.AppendChange(change); err != nil {
			return nil, err
		}
	}
	if err := transaction.Commit(b.position, b.at); err != nil {
		return nil, err
	}
	return transaction, nil
}

// sameNames vergleicht zwei Spaltenlisten als Mengen, unabhängig von
// Reihenfolge und Doppelungen.
func sameNames(a, b []string) bool {
	left := nameSet(a)
	right := nameSet(b)
	if len(left) != len(right) {
		return false
	}
	for name := range left {
		if _, ok := right[name]; !ok {
			return false
		}
	}
	return true
}

func nameSet(names []string) map[string]struct{} {
	set := make(map[string]struct{}, len(names))
	for _, name := range names {
		set[name] = struct{}{}
	}
	return set
}
