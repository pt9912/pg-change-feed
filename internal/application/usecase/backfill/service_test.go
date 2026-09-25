package backfill_test

import (
	"context"
	stderrors "errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/backfill"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

const (
	testSource      = model.SourceID("src-1")
	testSchema      = "public"
	testTable       = "orders"
	testPublication = "cdc_pub"
	testRequestID   = model.AdministrationRequestID("run-1")
	testTableID     = model.SourceTableID("tbl-1")
	testVersionID   = model.SchemaVersionID("tbl-1-v1")
	testOffset      = uint64(0x100)
)

// trace hält die Aufruf-Reihenfolge aller Fakes; Tests prüfen Reihenfolgen
// und Aufrufzahlen an ihr.
type trace struct{ events []string }

func (t *trace) add(format string, args ...any) {
	t.events = append(t.events, fmt.Sprintf(format, args...))
}

func (t *trace) count(prefix string) int {
	n := 0
	for _, e := range t.events {
		if strings.HasPrefix(e, prefix) {
			n++
		}
	}
	return n
}

// only liefert die Ereignisse, die mit einem der Präfixe beginnen, in
// ihrer Reihenfolge.
func (t *trace) only(prefixes ...string) []string {
	var out []string
	for _, e := range t.events {
		for _, p := range prefixes {
			if strings.HasPrefix(e, p) {
				out = append(out, e)
				break
			}
		}
	}
	return out
}

func str(v string) *string { return &v }

// --- Fakes ---------------------------------------------------------------

type fakeActivation struct {
	trace *trace
	table model.SourceTable

	registeredFn func(call int) (model.SourceTable, bool, error)
	registered   int
	published    bool
	publishedErr error

	gotSource      model.SourceID
	gotSchema      string
	gotTable       string
	gotPublication string
	publishedCalls int
}

func (f *fakeActivation) Registered(ctx context.Context, source model.SourceID, schema, table string) (model.SourceTable, bool, error) {
	f.registered++
	f.trace.add("Registered")
	f.gotSource, f.gotSchema, f.gotTable = source, schema, table
	if f.registeredFn != nil {
		return f.registeredFn(f.registered)
	}
	return f.table, true, nil
}

func (f *fakeActivation) Published(ctx context.Context, publication, schema, table string) (bool, error) {
	f.publishedCalls++
	f.trace.add("Published")
	f.gotPublication = publication
	if schema != f.gotSchema || table != f.gotTable {
		return false, fmt.Errorf("Published(%s.%s) weicht von Registered(%s.%s) ab", schema, table, f.gotSchema, f.gotTable)
	}
	return f.published, f.publishedErr
}

func (f *fakeActivation) TableExists(context.Context, string, string) (bool, error) {
	return true, nil
}
func (f *fakeActivation) Register(context.Context, model.SourceTable, model.SchemaVersion) (bool, error) {
	return false, nil
}
func (f *fakeActivation) Unregister(context.Context, model.SourceTable) (outbound.ActivationRemoval, error) {
	return outbound.ActivationAbsent, nil
}
func (f *fakeActivation) List(context.Context, model.SourceID) ([]model.SourceTable, error) {
	return nil, nil
}
func (f *fakeActivation) Publish(context.Context, string, string, string) error   { return nil }
func (f *fakeActivation) Unpublish(context.Context, string, string, string) error { return nil }

type fakeExclusion struct {
	trace *trace
	// stateFn liefert den Stand des n-ten Aufrufs (ab 1).
	stateFn func(call int) map[string][]string
	// err lässt jeden Aufruf scheitern; ist errCall gesetzt, scheitert nur
	// der errCall-te Aufruf (ab 1).
	err       error
	errCall   int
	calls     int
	gotSource model.SourceID
}

func (f *fakeExclusion) ExcludedColumns(ctx context.Context, source model.SourceID) (map[string][]string, error) {
	f.calls++
	f.trace.add("ExcludedColumns")
	f.gotSource = source
	if f.err != nil && (f.errCall == 0 || f.errCall == f.calls) {
		return nil, f.err
	}
	if f.stateFn == nil {
		return nil, nil
	}
	return f.stateFn(f.calls), nil
}

func (f *fakeExclusion) ColumnExists(context.Context, string, string, string) (bool, error) {
	return true, nil
}

type fakeSchemas struct {
	version model.SchemaVersion
	found   bool
	err     error
	gotID   model.SourceTableID
}

func (f *fakeSchemas) CurrentVersion(ctx context.Context, table model.SourceTableID) (model.SchemaVersion, bool, error) {
	f.gotID = table
	return f.version, f.found, f.err
}
func (f *fakeSchemas) RegisterVersion(context.Context, model.SchemaVersion, model.TableSchema) (bool, error) {
	return false, nil
}
func (f *fakeSchemas) TableSchema(context.Context, model.SchemaVersionID) (model.TableSchema, error) {
	return model.TableSchema{}, nil
}

type fakeSnapshotPort struct {
	trace    *trace
	snapshot *fakeSnapshot
	openErr  error
	// beforeOpen läuft am Anfang jedes OpenSnapshot-Aufrufs, onOpen nach
	// einem erfolgreichen Öffnen, bevor der Snapshot zurückgegeben wird.
	beforeOpen func()
	onOpen     func()

	estimateRows  int64
	estimateKnown bool
	estimateErr   error

	openCalls int
	gotRunID  string
	gotSchema string
	gotTable  string
}

func (f *fakeSnapshotPort) OpenSnapshot(ctx context.Context, runID, schema, table string) (outbound.TableSnapshot, error) {
	f.openCalls++
	f.trace.add("OpenSnapshot")
	f.gotRunID, f.gotSchema, f.gotTable = runID, schema, table
	if f.beforeOpen != nil {
		f.beforeOpen()
	}
	if f.openErr != nil {
		return nil, f.openErr
	}
	if f.onOpen != nil {
		f.onOpen()
	}
	return f.snapshot, nil
}

func (f *fakeSnapshotPort) EstimatedRows(ctx context.Context, schema, table string) (int64, bool, error) {
	f.trace.add("EstimatedRows")
	return f.estimateRows, f.estimateKnown, f.estimateErr
}

type fakeSnapshot struct {
	trace   *trace
	offset  uint64
	columns []string
	blocks  [][][]*string

	// nextErrAt lässt den n-ten NextBlock-Aufruf (ab 1) mit nextErr scheitern.
	nextErrAt int
	nextErr   error
	// onNext läuft am Anfang jedes NextBlock-Aufrufs (ab 1).
	onNext func(call int)

	closeErr error

	nextCalls int
	closed    int
	closeCtx  []error
}

func (f *fakeSnapshot) Offset() uint64    { return f.offset }
func (f *fakeSnapshot) Columns() []string { return f.columns }

func (f *fakeSnapshot) NextBlock(ctx context.Context) ([][]*string, error) {
	f.nextCalls++
	f.trace.add("NextBlock")
	if f.onNext != nil {
		f.onNext(f.nextCalls)
	}
	if f.nextErrAt == f.nextCalls {
		return nil, f.nextErr
	}
	if ctx.Err() != nil {
		return nil, fmt.Errorf("%w: Kontext beendet", outbound.ErrSnapshotTransient)
	}
	if f.nextCalls > len(f.blocks) {
		return nil, nil
	}
	return f.blocks[f.nextCalls-1], nil
}

func (f *fakeSnapshot) Close(ctx context.Context) error {
	f.closed++
	f.trace.add("Close")
	f.closeCtx = append(f.closeCtx, ctx.Err())
	return f.closeErr
}

type fakeAdmission struct {
	trace *trace
	err   error

	calls  int
	gotID  model.AdministrationRequestID
	gotRun model.BackfillRun
}

func (f *fakeAdmission) Admit(ctx context.Context, requestID model.AdministrationRequestID, run model.BackfillRun) error {
	f.calls++
	f.trace.add("Admit")
	f.gotID, f.gotRun = requestID, run
	return f.err
}

type fakeRuns struct {
	trace *trace

	markRunningErr error
	// progressErr lässt jeden RecordProgress-Aufruf scheitern; ist
	// progressErrCall gesetzt, scheitert nur der progressErrCall-te Aufruf
	// (ab 1). finishErr und finishErrCall verhalten sich für Finish ebenso.
	progressErr     error
	progressErrCall int
	finishErr       error
	finishErrCall   int

	marked     []model.BackfillRun
	progresses []model.BackfillRun
	finished   []model.BackfillRun
	// finishCtxErr hält ctx.Err() je Finish-Aufruf.
	finishCtxErr []error
}

func (f *fakeRuns) Queued(context.Context, model.SourceID) ([]model.BackfillRun, error) {
	return nil, nil
}

func (f *fakeRuns) MarkRunning(ctx context.Context, run model.BackfillRun) error {
	f.trace.add("MarkRunning")
	f.marked = append(f.marked, run)
	return f.markRunningErr
}

func (f *fakeRuns) RecordProgress(ctx context.Context, run model.BackfillRun) error {
	f.trace.add("Progress")
	f.progresses = append(f.progresses, run)
	if f.progressErr != nil && (f.progressErrCall == 0 || f.progressErrCall == len(f.progresses)) {
		return f.progressErr
	}
	return nil
}

func (f *fakeRuns) Finish(ctx context.Context, run model.BackfillRun) error {
	f.trace.add("Finish(%s)", run.Status)
	f.finished = append(f.finished, run)
	f.finishCtxErr = append(f.finishCtxErr, ctx.Err())
	if f.finishErr != nil && (f.finishErrCall == 0 || f.finishErrCall == len(f.finished)) {
		return f.finishErr
	}
	return nil
}

func (f *fakeRuns) InterruptRunning(context.Context, model.SourceID, model.TimePoint) (int, error) {
	return 0, nil
}

type blockRecord struct {
	id          model.TransactionID
	source      model.SourceID
	position    model.SourcePosition
	committedAt model.TimePoint
	changes     []model.Change
}

type fakeWriter struct {
	trace *trace

	beginErr   error
	appendErr  error
	appendErrN int // scheitert der n-te AppendBlock-Aufruf (ab 1)
	commitErr  error
	// onAppend läuft am Anfang jedes AppendBlock-Aufrufs (ab 1).
	onAppend func(call int)
	// rollbackErr ist der Fehler von Rollback.
	rollbackErr error

	begun       int
	appendCalls int
	blocks      []blockRecord
	committed   []model.BackfillRun
	rollbacks   int
	rollCtx     []error
}

func (f *fakeWriter) Begin(ctx context.Context, run model.BackfillRun) (outbound.BackfillTransaction, error) {
	f.trace.add("Begin")
	if f.beginErr != nil {
		return nil, f.beginErr
	}
	f.begun++
	return &fakeTx{writer: f}, nil
}

type fakeTx struct{ writer *fakeWriter }

func (t *fakeTx) AppendBlock(ctx context.Context, block *model.ChangeTransaction) error {
	f := t.writer
	f.trace.add("Append")
	f.appendCalls++
	if f.onAppend != nil {
		f.onAppend(f.appendCalls)
	}
	if f.appendErrN == f.appendCalls && f.appendErr != nil {
		return f.appendErr
	}
	changes, err := block.Changes()
	if err != nil {
		return fmt.Errorf("Block nicht committed: %w", err)
	}
	position, _ := block.CommitPosition()
	at, _ := block.SourceCommittedAt()
	f.blocks = append(f.blocks, blockRecord{id: block.ID, source: block.SourceID, position: position, committedAt: at, changes: changes})
	return nil
}

func (t *fakeTx) Commit(ctx context.Context, run model.BackfillRun) error {
	f := t.writer
	f.trace.add("Commit")
	if f.commitErr != nil {
		return f.commitErr
	}
	f.committed = append(f.committed, run)
	return nil
}

func (t *fakeTx) Rollback(ctx context.Context) error {
	f := t.writer
	f.trace.add("Rollback")
	f.rollbacks++
	f.rollCtx = append(f.rollCtx, ctx.Err())
	return f.rollbackErr
}

// fakeClock liefert `at`; mit `step` > 0 läuft die Uhr je Aufruf um `step`
// weiter, sodass jede Lesung einen erkennbar anderen Zeitpunkt trägt. `last`
// hält die zuletzt gelieferte Lesung.
type fakeClock struct {
	at   int64
	step int64
	last model.TimePoint
}

func (c *fakeClock) Now() model.TimePoint {
	c.last = model.NewTimePoint(c.at)
	c.at += c.step
	return c.last
}

type fakeNotifier struct {
	trace *trace
	err   error
	calls []string
}

func (n *fakeNotifier) Notify(ctx context.Context, sourceID, schema, table string) error {
	n.trace.add("Notify")
	n.calls = append(n.calls, sourceID+"/"+schema+"."+table)
	return n.err
}

type fakeLog struct{ warns []string }

func (l *fakeLog) Debug(context.Context, string, ...any) {}
func (l *fakeLog) Info(context.Context, string, ...any)  {}
func (l *fakeLog) Warn(_ context.Context, msg string, _ ...any) {
	l.warns = append(l.warns, msg)
}
func (l *fakeLog) Error(context.Context, string, ...any) {}

// --- Rig -----------------------------------------------------------------

// rig verdrahtet den Use Case mit Fakes im Normalzustand: Tabelle
// aktiviert und in der Publication, drei Blöcke zu 2/2/1 Zeilen über die
// Spalten id/name/secret, kein Ausschluss.
type rig struct {
	trace      *trace
	activation *fakeActivation
	exclusion  *fakeExclusion
	schemas    *fakeSchemas
	snapshotP  *fakeSnapshotPort
	snapshot   *fakeSnapshot
	admission  *fakeAdmission
	runs       *fakeRuns
	writer     *fakeWriter
	clock      *fakeClock
	notifier   *fakeNotifier
	log        *fakeLog
	service    *backfill.BackfillTableService
}

func newRig() *rig {
	tr := &trace{}
	table, _ := model.NewSourceTable(testTableID, testSource, testSchema, testTable)
	version, _ := model.NewSchemaVersion(testVersionID, testTableID, 1)
	r := &rig{
		trace:      tr,
		activation: &fakeActivation{trace: tr, table: table, published: true},
		exclusion:  &fakeExclusion{trace: tr},
		schemas:    &fakeSchemas{version: version, found: true},
		snapshot: &fakeSnapshot{
			trace:   tr,
			offset:  testOffset,
			columns: []string{"id", "name", "secret"},
			blocks: [][][]*string{
				{{str("1"), str("a"), str("s1")}, {str("2"), nil, str("s2")}},
				{{str("3"), str("c"), str("s3")}, {str("4"), str("d"), str("s4")}},
				{{str("5"), str("e"), str("s5")}},
			},
		},
		admission: &fakeAdmission{trace: tr},
		runs:      &fakeRuns{trace: tr},
		writer:    &fakeWriter{trace: tr},
		clock:     &fakeClock{at: 1000},
		notifier:  &fakeNotifier{trace: tr},
		log:       &fakeLog{},
	}
	r.snapshotP = &fakeSnapshotPort{trace: tr, snapshot: r.snapshot, estimateRows: 1234, estimateKnown: true}
	r.service = backfill.NewBackfillTableService(backfill.Ports{
		Activation: r.activation,
		Exclusion:  r.exclusion,
		Schemas:    r.schemas,
		Snapshot:   r.snapshotP,
		Admission:  r.admission,
		Runs:       r.runs,
		Writer:     r.writer,
		Clock:      r.clock,
	}, backfill.WithChangeNotification(r.notifier), backfill.WithLog(r.log))
	return r
}

func requestCommand() backfill.BackfillRequestCommand {
	return backfill.BackfillRequestCommand{
		RequestID: testRequestID, Source: testSource, Schema: testSchema, Table: testTable, Publication: testPublication,
	}
}

func queuedRun(t *testing.T) model.BackfillRun {
	t.Helper()
	run, err := model.NewQueuedBackfillRun(model.BackfillRunID(testRequestID), testSource, testSchema, testTable, model.NewTimePoint(500))
	if err != nil {
		t.Fatal(err)
	}
	return run
}

func (r *rig) execute(ctx context.Context, t *testing.T) (backfill.BackfillExecuteResult, error) {
	t.Helper()
	return r.service.Execute(ctx, backfill.BackfillExecuteCommand{Run: queuedRun(t), Publication: testPublication})
}

func mustExecute(t *testing.T, r *rig) model.BackfillRun {
	t.Helper()
	result, err := r.execute(context.Background(), t)
	if err != nil {
		t.Fatalf("Execute = %v, will nil", err)
	}
	return result.Run
}

// --- Request -------------------------------------------------------------

// TestRequestAdmitsLastAfterPreconditionsAndEstimate trägt die Reihenfolge der
// Annahme (`ADR-0113` Festlegung 1): Bindung, Publication, Schätzung, dann
// `Admit` als letzter Schritt — und ohne geöffneten Snapshot. Der Run, den
// `Admit` erhält, trägt die Antrags-Kennung, `queued`, die Tabelle und den
// Anlage-Zeitpunkt der Uhr.
func TestRequestAdmitsLastAfterPreconditionsAndEstimate(t *testing.T) {
	r := newRig()
	result, err := r.service.Request(context.Background(), requestCommand())
	if err != nil {
		t.Fatalf("Request = %v", err)
	}
	want := []string{"Registered", "Published", "EstimatedRows", "Admit"}
	if got := r.trace.events; !reflect.DeepEqual(got, want) {
		t.Fatalf("Aufrufe = %v, will %v", got, want)
	}
	if r.snapshotP.openCalls != 0 {
		t.Fatalf("OpenSnapshot-Aufrufe = %d, will 0", r.snapshotP.openCalls)
	}
	run := r.admission.gotRun
	if r.admission.gotID != testRequestID || run.ID != model.BackfillRunID(testRequestID) {
		t.Fatalf("Admit erhielt Antrag %q / Run %q", r.admission.gotID, run.ID)
	}
	if run.Status != model.BackfillRunQueued || run.Source != testSource || run.Schema != testSchema || run.Table != testTable {
		t.Fatalf("Admit erhielt %+v", run)
	}
	if run.RequestedAt != model.NewTimePoint(1000) {
		t.Fatalf("RequestedAt = %v, will die Uhr 1000", run.RequestedAt)
	}
	if run.WarnEstimatedSize || run.WarnDuration {
		t.Fatalf("Warn-Kennzeichnungen gesetzt: %+v", run)
	}
	if result.Run != run {
		t.Fatalf("Ergebnis-Run %+v ≠ angenommener Run %+v", result.Run, run)
	}
	if r.activation.gotSource != testSource || r.activation.gotSchema != testSchema || r.activation.gotTable != testTable || r.activation.gotPublication != testPublication {
		t.Fatalf("Vorbedingungen adressierten %s %s.%s in %q", r.activation.gotSource, r.activation.gotSchema, r.activation.gotTable, r.activation.gotPublication)
	}
}

// TestRequestCarriesEstimate trägt die drei Zustände der Schätzung
// (`ADR-0113` Festlegung 3, Punkt 5): eine Zahl bleibt die Zahl, eine
// unbekannte Schätzung (der Katalog führt `−1`) erreicht `Admit` als
// „unbekannt" und nicht als `0`, eine bekannte `0` bleibt bekannt `0`.
func TestRequestCarriesEstimate(t *testing.T) {
	cases := []struct {
		name      string
		rows      int64
		known     bool
		wantRows  int64
		wantKnown bool
	}{
		{"Zahl", 1234, true, 1234, true},
		{"unbekannt (-1)", -1, false, 0, false},
		{"unbekannt (0, nicht bekannt)", 0, false, 0, false},
		{"bekannte Null", 0, true, 0, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := newRig()
			r.snapshotP.estimateRows, r.snapshotP.estimateKnown = tc.rows, tc.known
			if _, err := r.service.Request(context.Background(), requestCommand()); err != nil {
				t.Fatalf("Request = %v", err)
			}
			rows, known := r.admission.gotRun.EstimatedRows.Rows()
			if known != tc.wantKnown || (known && rows != tc.wantRows) {
				t.Fatalf("Admit erhielt (%d, known=%t), will (%d, known=%t)", rows, known, tc.wantRows, tc.wantKnown)
			}
		})
	}
}

// TestRequestPreconditionFailuresLeaveNoRun trägt die Vorbedingungs-Fehler
// (`LH-FA-CAP-009`, `ADR-0111` Teilfrage 5): jeder endet vor `Admit`, ohne
// Run-Zeile, ohne geöffneten Snapshot und ohne Schätzung; eine fehlende
// Bindung und eine fehlende Publication-Mitgliedschaft tragen
// `ErrTableNotActivated`, ein Port-Fehler bleibt sein Fehler.
func TestRequestPreconditionFailuresLeaveNoRun(t *testing.T) {
	portErr := stderrors.New("Katalog nicht lesbar")
	cases := []struct {
		name  string
		setup func(*rig)
		want  error
	}{
		{"keine Bindung", func(r *rig) {
			r.activation.registeredFn = func(int) (model.SourceTable, bool, error) { return model.SourceTable{}, false, nil }
		}, domainerrors.ErrTableNotActivated},
		{"nicht in der Publication", func(r *rig) { r.activation.published = false }, domainerrors.ErrTableNotActivated},
		{"Fehler beim Lesen der Bindung", func(r *rig) {
			r.activation.registeredFn = func(int) (model.SourceTable, bool, error) { return model.SourceTable{}, false, portErr }
		}, portErr},
		{"Fehler beim Lesen der Publication", func(r *rig) { r.activation.publishedErr = portErr }, portErr},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := newRig()
			tc.setup(r)
			_, err := r.service.Request(context.Background(), requestCommand())
			if !stderrors.Is(err, tc.want) {
				t.Fatalf("Request = %v, will %v", err, tc.want)
			}
			if r.admission.calls != 0 || r.snapshotP.openCalls != 0 || r.trace.count("EstimatedRows") != 0 {
				t.Fatalf("nach dem Vorbedingungs-Fehler liefen Admit/OpenSnapshot/EstimatedRows: %v", r.trace.events)
			}
		})
	}
}

// TestRequestActiveRunEndsWithoutSecondRow trägt „kein zweiter aktiver Run":
// meldet `Admit` `ErrBackfillRunActive`, endet der Antrag mit demselben
// Sentinel, ohne zweiten Versuch und ohne Ergebnis-Run.
func TestRequestActiveRunEndsWithoutSecondRow(t *testing.T) {
	r := newRig()
	r.admission.err = fmt.Errorf("%w: public.orders", domainerrors.ErrBackfillRunActive)
	result, err := r.service.Request(context.Background(), requestCommand())
	if !stderrors.Is(err, domainerrors.ErrBackfillRunActive) {
		t.Fatalf("Request = %v, will ErrBackfillRunActive", err)
	}
	if r.admission.calls != 1 {
		t.Fatalf("Admit-Aufrufe = %d, will 1", r.admission.calls)
	}
	if (result != backfill.BackfillRequestResult{}) {
		t.Fatalf("Ergebnis = %+v, will leer", result)
	}
}

// TestRequestErrorsBeforeAdmit trägt die übrigen Fehler vor `Admit`: die
// Schätzung scheitert, der Katalog meldet eine negative bekannte Zahl, ein
// Argument ist leer — jeweils ohne `Admit`.
func TestRequestErrorsBeforeAdmit(t *testing.T) {
	portErr := stderrors.New("Katalog nicht lesbar")

	r := newRig()
	r.snapshotP.estimateErr = portErr
	if _, err := r.service.Request(context.Background(), requestCommand()); !stderrors.Is(err, portErr) || r.admission.calls != 0 {
		t.Fatalf("Schätzungs-Fehler: err=%v admit=%d", err, r.admission.calls)
	}

	r = newRig()
	r.snapshotP.estimateRows, r.snapshotP.estimateKnown = -5, true
	if _, err := r.service.Request(context.Background(), requestCommand()); !stderrors.Is(err, domainerrors.ErrNegativeRowCount) || r.admission.calls != 0 {
		t.Fatalf("negative Schätzung: err=%v admit=%d", err, r.admission.calls)
	}

	r = newRig()
	command := requestCommand()
	command.RequestID = ""
	if _, err := r.service.Request(context.Background(), command); !stderrors.Is(err, domainerrors.ErrEmptyIdentifier) || len(r.trace.events) != 0 {
		t.Fatalf("leere Antrags-Kennung: err=%v Aufrufe=%v", err, r.trace.events)
	}
}

// TestPortSchnitt trägt die Port-Definitionen (`ADR-0113` Festlegung 1): der
// Annahme-Port hat genau `Admit`, der Run-Zustands-Port keine Anlage-Operation
// — die Run-Zeile entsteht auf genau einem Weg.
func TestPortSchnitt(t *testing.T) {
	methods := func(port any) []string {
		typ := reflect.TypeOf(port).Elem()
		names := make([]string, 0, typ.NumMethod())
		for i := 0; i < typ.NumMethod(); i++ {
			names = append(names, typ.Method(i).Name)
		}
		sort.Strings(names)
		return names
	}
	if got, want := methods((*outbound.BackfillAdmissionPort)(nil)), []string{"Admit"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("BackfillAdmissionPort = %v, will %v", got, want)
	}
	want := []string{"Finish", "InterruptRunning", "MarkRunning", "Queued", "RecordProgress"}
	if got := methods((*outbound.BackfillRunPort)(nil)); !reflect.DeepEqual(got, want) {
		t.Fatalf("BackfillRunPort = %v, will %v", got, want)
	}
}

// TestPortsCarryNoCapturePathPort trägt „run-lokal": der Use Case hält weder
// den Heartbeat- noch den Bestätigungs-, Store- oder Stream-Port des
// Capture-Pfads; ein Run-Fehler kann den Heartbeat-Fehlerzustand deshalb nicht
// setzen.
func TestPortsCarryNoCapturePathPort(t *testing.T) {
	typ := reflect.TypeOf(backfill.Ports{})
	for i := 0; i < typ.NumField(); i++ {
		name := typ.Field(i).Type.Name()
		for _, forbidden := range []string{"Heartbeat", "ReplicationAck", "ChangeStore", "ChangeStream"} {
			if strings.Contains(name, forbidden) {
				t.Fatalf("Ports.%s trägt den Capture-Pfad-Port %s", typ.Field(i).Name, name)
			}
		}
	}
}

// --- Execute: Happy Path --------------------------------------------------

// TestExecuteWritesAllBlocksInOneTransaction trägt den Happy Path
// (`LH-FA-CAP-009`, `ADR-0111` Teilfrage 2/6): alle Blöcke in **einer**
// Transaktion, ein Commit mit der Run-Zeile `completed`; jeder Change ist ein
// `INSERT` der Herkunft `backfill` ohne `old_data`, sein Bild trägt nur die
// nicht ausgeschlossenen, nicht-NULL-Spalten; Position `X` an jedem Block;
// Transaktions- und Change-Kennungen nach der Bildungsregel; genau ein
// Wecksignal.
func TestExecuteWritesAllBlocksInOneTransaction(t *testing.T) {
	r := newRig()
	// Der Ausschluss der eigenen Tabelle wirkt; der Eintrag einer anderen
	// Tabelle bleibt ohne Wirkung.
	r.exclusion.stateFn = func(int) map[string][]string {
		return map[string][]string{"public.orders": {"secret"}, "public.other": {"name"}}
	}
	run := mustExecute(t, r)

	if r.writer.begun != 1 || r.trace.count("Commit") != 1 || r.writer.rollbacks != 0 {
		t.Fatalf("Begin=%d Commit=%d Rollback=%d, will 1/1/0", r.writer.begun, r.trace.count("Commit"), r.writer.rollbacks)
	}
	if len(r.writer.blocks) != 3 {
		t.Fatalf("Blöcke = %d, will 3", len(r.writer.blocks))
	}
	wantImages := [][]string{
		{`{"id":"1","name":"a"}`, `{"id":"2"}`},
		{`{"id":"3","name":"c"}`, `{"id":"4","name":"d"}`},
		{`{"id":"5","name":"e"}`},
	}
	position, _ := model.NewSourcePosition(testSource, testOffset)
	for i, block := range r.writer.blocks {
		wantID := model.TransactionID(fmt.Sprintf("0bf-run-1-%08d", i+1))
		if block.id != wantID || block.source != testSource {
			t.Fatalf("Block %d: Transaktion %q/%q, will %q/%q", i+1, block.id, block.source, wantID, testSource)
		}
		if block.position != position {
			t.Fatalf("Block %d: Position %+v, will %+v", i+1, block.position, position)
		}
		if block.committedAt != model.NewTimePoint(1000) {
			t.Fatalf("Block %d: committed_at %v", i+1, block.committedAt)
		}
		if len(block.changes) != len(wantImages[i]) {
			t.Fatalf("Block %d: %d Changes, will %d", i+1, len(block.changes), len(wantImages[i]))
		}
		for j, change := range block.changes {
			sequence := int64(j + 1)
			if change.Sequence != sequence || change.ID != model.ChangeID(fmt.Sprintf("%s-%d", wantID, sequence)) || change.TransactionID != wantID {
				t.Fatalf("Block %d Change %d: Kennung/Sequenz %q/%d", i+1, j+1, change.ID, change.Sequence)
			}
			if change.Operation != model.OperationInsert || change.Origin != model.ChangeOriginBackfill || change.OldImage != nil {
				t.Fatalf("Block %d Change %d: Operation=%s Origin=%s OldImage=%q", i+1, j+1, change.Operation, change.Origin, change.OldImage)
			}
			if string(change.NewImage) != wantImages[i][j] {
				t.Fatalf("Block %d Change %d: Bild %s, will %s", i+1, j+1, change.NewImage, wantImages[i][j])
			}
			if change.SourceTableID != testTableID || change.SchemaVersion != testVersionID || change.Schema != testSchema || change.Table != testTable {
				t.Fatalf("Block %d Change %d: Tabelle/Version/Klartext %q %q %s.%s", i+1, j+1, change.SourceTableID, change.SchemaVersion, change.Schema, change.Table)
			}
		}
	}
	if len(r.writer.committed) != 1 {
		t.Fatalf("Commits = %d", len(r.writer.committed))
	}
	committed := r.writer.committed[0]
	if committed.Status != model.BackfillRunCompleted || committed.RowsCopied != 5 || committed.FinishedAt != model.NewTimePoint(1000) {
		t.Fatalf("Commit trug %+v", committed)
	}
	if run != committed {
		t.Fatalf("Ergebnis-Run %+v ≠ committeter Run %+v", run, committed)
	}
	if r.trace.count("Finish") != 0 {
		t.Fatalf("Finish neben dem Commit gerufen: %v", r.trace.events)
	}
	if !reflect.DeepEqual(r.notifier.calls, []string{"src-1/public.orders"}) {
		t.Fatalf("Wecksignale = %v, will genau eines je Tabelle", r.notifier.calls)
	}
	if r.exclusion.gotSource != testSource || r.schemas.gotID != testTableID {
		t.Fatalf("Ausschluss/Schema-Version lasen Quelle %q / Tabelle %q", r.exclusion.gotSource, r.schemas.gotID)
	}
	if r.snapshotP.gotRunID != string(testRequestID) || r.snapshotP.gotSchema != testSchema || r.snapshotP.gotTable != testTable {
		t.Fatalf("OpenSnapshot(%q, %q, %q)", r.snapshotP.gotRunID, r.snapshotP.gotSchema, r.snapshotP.gotTable)
	}
	if r.snapshot.closed != 1 {
		t.Fatalf("Snapshot %d-mal geschlossen, will 1", r.snapshot.closed)
	}
}

// TestExecuteRecordsProgressPerBlock trägt den Fortschritt (`SPEC-029`):
// `running` mit `started_at`, dann die Position `X` mit dem Zähler 0 und je
// Block der wachsende Zähler.
func TestExecuteRecordsProgressPerBlock(t *testing.T) {
	r := newRig()
	mustExecute(t, r)
	if len(r.runs.marked) != 1 || r.runs.marked[0].Status != model.BackfillRunRunning || r.runs.marked[0].StartedAt != model.NewTimePoint(1000) {
		t.Fatalf("MarkRunning = %+v", r.runs.marked)
	}
	var counts []int64
	position, _ := model.NewSourcePosition(testSource, testOffset)
	for _, progress := range r.runs.progresses {
		if progress.SnapshotPosition != position {
			t.Fatalf("Fortschritt trägt Position %+v", progress.SnapshotPosition)
		}
		counts = append(counts, progress.RowsCopied)
	}
	if want := []int64{0, 2, 4, 5}; !reflect.DeepEqual(counts, want) {
		t.Fatalf("Fortschrittszähler = %v, will %v", counts, want)
	}
}

// TestExecuteMarksRunningBeforeOpeningSnapshot trägt die Reihenfolge des
// Zustandswechsels (`ADR-0111` Teilfrage 4): die Slot-Anlage im Snapshot-Port
// wartet auf laufende Schreibtransaktionen der Quelle, der Run steht dabei
// schon `running`. Der Snapshot-Port sieht beim Öffnen den festgehaltenen
// Zustand.
func TestExecuteMarksRunningBeforeOpeningSnapshot(t *testing.T) {
	r := newRig()
	var runningAtOpen []model.BackfillRunStatus
	r.snapshotP.beforeOpen = func() {
		for _, run := range r.runs.marked {
			runningAtOpen = append(runningAtOpen, run.Status)
		}
	}
	mustExecute(t, r)
	if want := []model.BackfillRunStatus{model.BackfillRunRunning}; !reflect.DeepEqual(runningAtOpen, want) {
		t.Fatalf("beim Öffnen des Snapshots festgehaltene Zustände = %v, will %v", runningAtOpen, want)
	}
	if got, want := r.trace.only("MarkRunning", "OpenSnapshot"), []string{"MarkRunning", "OpenSnapshot"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Reihenfolge = %v, will %v", got, want)
	}
}

// TestExecuteStreamsBlocks trägt `LH-FA-CAP-006.a`: der Schreiber erhält
// Block `n`, bevor der Leser Block `n+1` liefert — höchstens ein Block im
// Speicher des Use Cases.
func TestExecuteStreamsBlocks(t *testing.T) {
	r := newRig()
	mustExecute(t, r)
	want := []string{"NextBlock", "Append", "NextBlock", "Append", "NextBlock", "Append", "NextBlock"}
	if got := r.trace.only("NextBlock", "Append"); !reflect.DeepEqual(got, want) {
		t.Fatalf("Reihenfolge = %v, will %v", got, want)
	}
}

// TestExecuteEmptyTable trägt die leere Tabelle (`LH-FA-CAP-009` Boundary):
// `completed` mit 0 Zeilen, keine Schreibtransaktion, kein Wecksignal; der
// Endzustand geht über den Run-Zustands-Port.
func TestExecuteEmptyTable(t *testing.T) {
	r := newRig()
	r.snapshot.blocks = nil
	run := mustExecute(t, r)
	if run.Status != model.BackfillRunCompleted || run.RowsCopied != 0 {
		t.Fatalf("Run = %+v", run)
	}
	if r.trace.count("Begin") != 0 || r.trace.count("Append") != 0 || r.trace.count("Commit") != 0 {
		t.Fatalf("Schreiber gerufen: %v", r.trace.events)
	}
	if len(r.runs.finished) != 1 || r.runs.finished[0] != run {
		t.Fatalf("Finish = %+v", r.runs.finished)
	}
	if r.trace.count("Notify") != 0 {
		t.Fatalf("Wecksignal für eine leere Tabelle")
	}
	if r.snapshot.closed != 1 {
		t.Fatalf("Snapshot %d-mal geschlossen", r.snapshot.closed)
	}
}

// TestExecuteCommittedAtIsSnapshotTime trägt die Retention-Grundlage
// (`ADR-0111` Teilfrage 7): `committed_at` der Blöcke ist die Uhr **nach**
// dem Öffnen des Snapshots, `started_at` die Uhr davor.
func TestExecuteCommittedAtIsSnapshotTime(t *testing.T) {
	r := newRig()
	r.snapshotP.onOpen = func() { r.clock.at = 2000 }
	mustExecute(t, r)
	if r.runs.marked[0].StartedAt != model.NewTimePoint(1000) {
		t.Fatalf("started_at = %v, will die Uhr vor dem Snapshot", r.runs.marked[0].StartedAt)
	}
	for i, block := range r.writer.blocks {
		if block.committedAt != model.NewTimePoint(2000) {
			t.Fatalf("Block %d: committed_at = %v, will 2000 (Snapshot-Zeitpunkt)", i+1, block.committedAt)
		}
	}
}

// TestExecuteWithoutNotifier trägt den optionalen Port: ohne
// `ChangeNotificationPort` läuft der Run ohne Wecksignal durch.
func TestExecuteWithoutNotifier(t *testing.T) {
	r := newRig()
	service := backfill.NewBackfillTableService(backfill.Ports{
		Activation: r.activation, Exclusion: r.exclusion, Schemas: r.schemas, Snapshot: r.snapshotP,
		Admission: r.admission, Runs: r.runs, Writer: r.writer, Clock: r.clock,
	})
	result, err := service.Execute(context.Background(), backfill.BackfillExecuteCommand{Run: queuedRun(t), Publication: testPublication})
	if err != nil || result.Run.Status != model.BackfillRunCompleted {
		t.Fatalf("Execute = %+v, %v", result, err)
	}
	if r.trace.count("Notify") != 0 {
		t.Fatalf("Wecksignal ohne Port")
	}
}

// TestExecuteNotifyIsBestEffort trägt `ADR-0055`: ein fehlgeschlagenes
// Wecksignal ändert den abgeschlossenen Run nicht und wird protokolliert.
func TestExecuteNotifyIsBestEffort(t *testing.T) {
	r := newRig()
	r.notifier.err = outbound.ErrNotify
	run := mustExecute(t, r)
	if run.Status != model.BackfillRunCompleted || len(r.writer.committed) != 1 {
		t.Fatalf("Run = %+v", run)
	}
	if len(r.log.warns) != 1 {
		t.Fatalf("Warnungen = %v, will genau eine", r.log.warns)
	}
}

// --- Execute: Vorbedingungen ---------------------------------------------

// TestExecuteRechecksPreconditions trägt `ADR-0113` Festlegung 2 Punkt 4:
// fehlen Bindung oder Publication-Mitgliedschaft erst in `Execute`, endet der
// Run `failed` mit der Klasse `configuration` — ohne Snapshot, ohne Schreiber,
// ohne `running`.
func TestExecuteRechecksPreconditions(t *testing.T) {
	cases := []struct {
		name  string
		setup func(*rig)
	}{
		{"Bindung fehlt", func(r *rig) {
			r.activation.registeredFn = func(int) (model.SourceTable, bool, error) { return model.SourceTable{}, false, nil }
		}},
		{"nicht in der Publication", func(r *rig) { r.activation.published = false }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := newRig()
			tc.setup(r)
			run := mustExecute(t, r)
			if run.Status != model.BackfillRunFailed || !strings.HasPrefix(run.ErrorMessage, "configuration: ") {
				t.Fatalf("Run = %+v", run)
			}
			if !run.StartedAt.Unset() {
				t.Fatalf("started_at gesetzt: %v", run.StartedAt)
			}
			if r.snapshotP.openCalls != 0 || r.trace.count("Begin") != 0 || r.trace.count("MarkRunning") != 0 || r.trace.count("Notify") != 0 {
				t.Fatalf("Ports gerufen: %v", r.trace.events)
			}
			if len(r.runs.finished) != 1 || r.runs.finished[0] != run {
				t.Fatalf("Finish = %+v", r.runs.finished)
			}
			if r.activation.publishedCalls > 0 && r.activation.gotPublication != testPublication {
				t.Fatalf("Publication %q geprüft, will %q", r.activation.gotPublication, testPublication)
			}
		})
	}
}

// TestExecuteRejectsNonQueuedRun trägt den Eintritt: nur ein `queued`-Run
// wird ausgeführt.
func TestExecuteRejectsNonQueuedRun(t *testing.T) {
	r := newRig()
	running, err := queuedRun(t).Start(model.NewTimePoint(1))
	if err != nil {
		t.Fatal(err)
	}
	_, err = r.service.Execute(context.Background(), backfill.BackfillExecuteCommand{Run: running, Publication: testPublication})
	if !stderrors.Is(err, domainerrors.ErrInvalidBackfillTransition) || len(r.trace.events) != 0 {
		t.Fatalf("Execute = %v, Aufrufe %v", err, r.trace.events)
	}
}

// --- Execute: Fehlerklassen und Rollback -----------------------------------

// TestExecuteClassifiesFailures trägt die Abbildung der Fehlerursachen auf die
// `SPEC-008`-Klassen, je Sentinel ein Fall an seiner Eingabe: die fünf
// Sentinels des Snapshot-Ports, die Storage-Sentinels der Backfill-Ports und
// des Schema Stores, `ErrSchemaVersionUnknown` als `configuration`; ein nicht
// erkannter Fehler bleibt `internal`, die Klasse `schema` vergibt der Run nicht. Jeder Fall endet den Run `failed` mit der Klasse
// vor dem Text, ohne Commit und ohne Wecksignal.
func TestExecuteClassifiesFailures(t *testing.T) {
	cause := func(sentinel error) error { return fmt.Errorf("%w: technische Ursache", sentinel) }
	cases := []struct {
		name  string
		setup func(*rig)
		class string
	}{
		{"Snapshot permission", func(r *rig) { r.snapshotP.openErr = cause(outbound.ErrSnapshotPermission) }, "permission"},
		{"Snapshot configuration", func(r *rig) { r.snapshotP.openErr = cause(outbound.ErrSnapshotConfiguration) }, "configuration"},
		{"Snapshot transient", func(r *rig) { r.snapshotP.openErr = cause(outbound.ErrSnapshotTransient) }, "transient"},
		{"Snapshot replication", func(r *rig) { r.snapshotP.openErr = cause(outbound.ErrSnapshotReplication) }, "replication"},
		{"Snapshot storage", func(r *rig) { r.snapshotP.openErr = cause(outbound.ErrSnapshotStorage) }, "storage"},
		{"Lesefehler im Block", func(r *rig) {
			r.snapshot.nextErrAt, r.snapshot.nextErr = 2, cause(outbound.ErrSnapshotStorage)
		}, "storage"},
		{"Schreiber Begin", func(r *rig) { r.writer.beginErr = cause(outbound.ErrBackfillStorage) }, "storage"},
		{"Schreiber AppendBlock", func(r *rig) {
			r.writer.appendErr, r.writer.appendErrN = cause(outbound.ErrBackfillStorage), 2
		}, "storage"},
		{"Schreiber Commit", func(r *rig) { r.writer.commitErr = cause(outbound.ErrBackfillStorage) }, "storage"},
		{"ChangeStore-Fehler am Schreiber", func(r *rig) { r.writer.commitErr = cause(outbound.ErrStorage) }, "storage"},
		{"Run-Zustand MarkRunning", func(r *rig) { r.runs.markRunningErr = cause(outbound.ErrBackfillStorage) }, "storage"},
		{"Run-Zustand Fortschritt", func(r *rig) { r.runs.progressErr = cause(outbound.ErrBackfillStorage) }, "storage"},
		{"Schema Store", func(r *rig) { r.schemas.err = cause(outbound.ErrSchemaStoreStorage) }, "storage"},
		{"Schema-Version unbekannt", func(r *rig) { r.schemas.found = false }, "configuration"},
		{"Ausschlussstand nicht lesbar", func(r *rig) { r.exclusion.err = stderrors.New("unbekannter Fehler") }, "internal"},
		{"Position ohne Offset", func(r *rig) { r.snapshot.offset = 0 }, "internal"},
		{"unbekannter Fehler", func(r *rig) { r.snapshotP.openErr = stderrors.New("unbekannt") }, "internal"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := newRig()
			tc.setup(r)
			run := mustExecute(t, r)
			if run.Status != model.BackfillRunFailed {
				t.Fatalf("Status = %q, will failed (%+v)", run.Status, run)
			}
			if !strings.HasPrefix(run.ErrorMessage, tc.class+": ") {
				t.Fatalf("Fehlertext = %q, will Klasse %q vorn", run.ErrorMessage, tc.class)
			}
			if len(r.runs.finished) != 1 || r.runs.finished[0] != run {
				t.Fatalf("Finish = %+v", r.runs.finished)
			}
			if r.trace.count("Notify") != 0 || len(r.writer.committed) != 0 {
				t.Fatalf("Commit/Wecksignal nach dem Fehler: %v", r.trace.events)
			}
			// Snapshot und Schreibtransaktion sind auf jedem Pfad beendet.
			if r.snapshotP.openCalls == 1 && r.snapshotP.openErr == nil && r.snapshot.closed != 1 {
				t.Fatalf("Snapshot %d-mal geschlossen, will 1", r.snapshot.closed)
			}
			if r.writer.begun > 0 && r.writer.rollbacks != 1 {
				t.Fatalf("Rollbacks = %d, will 1", r.writer.rollbacks)
			}
			if r.writer.begun == 0 && r.writer.rollbacks != 0 {
				t.Fatalf("Rollback ohne geöffnete Transaktion")
			}
		})
	}
}

// TestExecuteFailureTextCarriesClassOnce trägt: der Fehlertext nennt die
// Klasse einmal vor der Ursache, auch wenn der Fehlerwert des Ports sie
// selbst trägt („Fehlerklasse transient: …“).
func TestExecuteFailureTextCarriesClassOnce(t *testing.T) {
	r := newRig()
	r.snapshotP.openErr = fmt.Errorf("%w: Tabelle umgeschrieben", outbound.ErrSnapshotTransient)
	run := mustExecute(t, r)
	want := "transient: " + strings.TrimPrefix(outbound.ErrSnapshotTransient.Error(), "Fehlerklasse transient: ") + ": Tabelle umgeschrieben"
	if run.Status != model.BackfillRunFailed || run.ErrorMessage != want {
		t.Fatalf("Fehlertext = %q, will %q", run.ErrorMessage, want)
	}
	if strings.Contains(run.ErrorMessage, "Fehlerklasse") {
		t.Fatalf("Fehlertext trägt die Klasse doppelt: %q", run.ErrorMessage)
	}
}

// TestExecuteMidCopyFailureKeepsProgressAndWritesNothing trägt „keine Zeile
// geschrieben": scheitert das Lesen des zweiten Blocks, ist der erste Block
// zurückgerollt und nichts committet; der Endzustand trägt den zuletzt
// festgehaltenen Zähler.
func TestExecuteMidCopyFailureKeepsProgressAndWritesNothing(t *testing.T) {
	r := newRig()
	r.snapshot.nextErrAt, r.snapshot.nextErr = 2, fmt.Errorf("%w: Verbindung weg", outbound.ErrSnapshotStorage)
	run := mustExecute(t, r)
	if run.Status != model.BackfillRunFailed || run.RowsCopied != 2 {
		t.Fatalf("Run = %+v, will failed mit 2 kopierten Zeilen", run)
	}
	if len(r.writer.blocks) != 1 || len(r.writer.committed) != 0 || r.writer.rollbacks != 1 {
		t.Fatalf("Blöcke=%d Commits=%d Rollbacks=%d", len(r.writer.blocks), len(r.writer.committed), r.writer.rollbacks)
	}
	if got, want := r.trace.only("Append", "Rollback", "Finish"), []string{"Append", "Rollback", "Finish(failed)"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Reihenfolge = %v, will %v", got, want)
	}
}

// --- Execute: Fail-closed --------------------------------------------------

// TestExecuteFailClosed trägt `LH-QA-SEC-004` (`ADR-0111` Teilfrage 4): jede
// Abweichung des Ausschlussstands oder der Bindung zwischen dem Bau der
// Blöcke und dem Commit rollt zurück und endet den Run `failed`
// (`configuration`) — auch eine Abweichung in einem Zwischenblock, die am Ende
// wieder gleich ist. Ein reiner Ordnungsunterschied ist keine Abweichung.
func TestExecuteFailClosed(t *testing.T) {
	state := func(cols ...string) map[string][]string { return map[string][]string{"public.orders": cols} }
	// Zählung der Ausschluss-Lesungen: Aufruf 1..3 je ein Block, Aufruf 4 vor dem Commit.
	sequence := func(states ...map[string][]string) func(int) map[string][]string {
		return func(call int) map[string][]string { return states[call-1] }
	}
	cases := []struct {
		name    string
		setup   func(*rig)
		failed  bool
		wantErr error
	}{
		{"Ausschluss ab Block 2", func(r *rig) {
			r.exclusion.stateFn = sequence(state(), state("secret"), state("secret"), state("secret"))
		}, true, domainerrors.ErrExclusionStateChanged},
		{"Zwischenblock weicht ab, am Ende wieder gleich", func(r *rig) {
			r.exclusion.stateFn = sequence(state(), state("secret"), state(), state())
		}, true, domainerrors.ErrExclusionStateChanged},
		{"Abweichung erst vor dem Commit", func(r *rig) {
			r.exclusion.stateFn = sequence(state(), state(), state(), state("secret"))
		}, true, domainerrors.ErrExclusionStateChanged},
		{"Einschluss statt Ausschluss", func(r *rig) {
			r.exclusion.stateFn = sequence(state("secret"), state("secret"), state("secret"), state())
		}, true, domainerrors.ErrExclusionStateChanged},
		{"Ordnung und Doppelung sind keine Abweichung", func(r *rig) {
			r.exclusion.stateFn = sequence(state("secret", "name"), state("name", "secret"), state("name", "secret", "name"), state("secret", "name"))
		}, false, nil},
		{"Bindung fehlt vor dem Commit", func(r *rig) {
			r.activation.registeredFn = func(call int) (model.SourceTable, bool, error) {
				if call >= 2 {
					return model.SourceTable{}, false, nil
				}
				return r.activation.table, true, nil
			}
		}, true, domainerrors.ErrTableNotActivated},
		{"Bindung mit anderer Kennung vor dem Commit", func(r *rig) {
			r.activation.registeredFn = func(call int) (model.SourceTable, bool, error) {
				if call >= 2 {
					other := r.activation.table
					other.ID = "tbl-2"
					return other, true, nil
				}
				return r.activation.table, true, nil
			}
		}, true, domainerrors.ErrTableNotActivated},
		{"Bindung nicht lesbar vor dem Commit", func(r *rig) {
			r.activation.registeredFn = func(call int) (model.SourceTable, bool, error) {
				if call >= 2 {
					return model.SourceTable{}, false, outbound.ErrBackfillStorage
				}
				return r.activation.table, true, nil
			}
		}, true, outbound.ErrBackfillStorage},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := newRig()
			tc.setup(r)
			run := mustExecute(t, r)
			if !tc.failed {
				if run.Status != model.BackfillRunCompleted || len(r.writer.committed) != 1 {
					t.Fatalf("Run = %+v, will completed", run)
				}
				return
			}
			if run.Status != model.BackfillRunFailed {
				t.Fatalf("Status = %q, will failed (%+v)", run.Status, run)
			}
			class := "configuration: "
			if tc.wantErr == outbound.ErrBackfillStorage {
				class = "storage: "
			}
			wantText := strings.TrimPrefix(tc.wantErr.Error(), "Fehlerklasse "+class)
			if !strings.HasPrefix(run.ErrorMessage, class) || !strings.Contains(run.ErrorMessage, wantText) {
				t.Fatalf("Fehlertext = %q", run.ErrorMessage)
			}
			if len(r.writer.committed) != 0 || r.trace.count("Commit") != 0 {
				t.Fatalf("Commit trotz Abweichung: %v", r.trace.events)
			}
			if r.writer.rollbacks != 1 || r.trace.count("Notify") != 0 {
				t.Fatalf("Rollbacks=%d Wecksignale=%d", r.writer.rollbacks, r.trace.count("Notify"))
			}
		})
	}
}

// TestExecuteExclusionReadFailure trägt den Lesefehler-Zweig der
// Fail-closed-Prüfung (`LH-QA-SEC-004`, `ADR-0111` Teilfrage 4): ein Stand, der
// nicht gelesen werden kann, ist kein bestätigter Stand. Der n-te Lesefehler
// (Lesung 1..3 je Block, Lesung 4 unmittelbar vor dem Commit) endet den Run
// `failed` mit der Klasse der Ursache, ohne Commit und ohne Wecksignal; die
// Lesung bricht den Lauf ab (keine weitere Lesung, kein weiterer Block), die
// Schreibtransaktion ist zurückgerollt, der Snapshot geschlossen. Weil nur
// der n-te Aufruf scheitert, färbt sich der Test rot, sobald ein Lesefehler
// verworfen wird — dann läuft der Run bis zum Commit durch.
func TestExecuteExclusionReadFailure(t *testing.T) {
	cases := []struct {
		name      string
		call      int
		appended  int
		rollbacks int
		rows      int64
	}{
		{"vor dem ersten Block", 1, 0, 0, 0},
		{"vor dem zweiten Block", 2, 1, 1, 2},
		{"vor dem letzten Block", 3, 2, 1, 4},
		{"unmittelbar vor dem Commit", 4, 3, 1, 5},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := newRig()
			r.exclusion.err = fmt.Errorf("%w: Ausschlussstand nicht lesbar", outbound.ErrStorage)
			r.exclusion.errCall = tc.call
			run := mustExecute(t, r)
			if run.Status != model.BackfillRunFailed || !strings.HasPrefix(run.ErrorMessage, "storage: ") || !strings.Contains(run.ErrorMessage, "Ausschlussstand nicht lesbar") {
				t.Fatalf("Run = %+v, will failed mit storage-Klasse und Ursache", run)
			}
			if run.RowsCopied != tc.rows {
				t.Fatalf("RowsCopied = %d, will %d (zuletzt festgehaltener Zähler)", run.RowsCopied, tc.rows)
			}
			if r.exclusion.calls != tc.call {
				t.Fatalf("Ausschluss-Lesungen = %d, will %d (Abbruch mit dem Lesefehler)", r.exclusion.calls, tc.call)
			}
			if len(r.writer.blocks) != tc.appended {
				t.Fatalf("angehängte Blöcke = %d, will %d", len(r.writer.blocks), tc.appended)
			}
			if r.trace.count("Commit") != 0 || len(r.writer.committed) != 0 || r.trace.count("Notify") != 0 {
				t.Fatalf("Commit/Wecksignal trotz Lesefehler: %v", r.trace.events)
			}
			if r.writer.rollbacks != tc.rollbacks {
				t.Fatalf("Rollbacks = %d, will %d", r.writer.rollbacks, tc.rollbacks)
			}
			if r.snapshot.closed != 1 {
				t.Fatalf("Snapshot %d-mal geschlossen, will 1", r.snapshot.closed)
			}
			if len(r.runs.finished) != 1 || r.runs.finished[0] != run {
				t.Fatalf("Finish = %+v", r.runs.finished)
			}
		})
	}
}

// TestExecuteProgressFailure trägt die Fehlerzweige der Fortschritts-Schreibung
// (`SPEC-029`): der n-te Schreibfehler des Run-Zustands (Aufruf 1 mit der
// Anfangsposition, Aufruf 2..4 je Block) endet den Run `failed` (`storage`),
// ohne Commit, mit dem zuletzt **festgehaltenen** Zähler, zurückgerollter
// Transaktion und geschlossenem Snapshot; nach dem Fehler folgt kein weiterer
// Fortschritt. Nur der n-te Aufruf scheitert: ein verworfener Fehler lässt den
// Run bis zum Commit durchlaufen.
func TestExecuteProgressFailure(t *testing.T) {
	cases := []struct {
		name      string
		call      int
		appended  int
		rollbacks int
		rows      int64
	}{
		{"mit der Anfangsposition", 1, 0, 0, 0},
		{"nach dem ersten Block", 2, 1, 1, 0},
		{"nach dem zweiten Block", 3, 2, 1, 2},
		{"nach dem letzten Block", 4, 3, 1, 4},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := newRig()
			r.runs.progressErr = fmt.Errorf("%w: Fortschritt nicht schreibbar", outbound.ErrBackfillStorage)
			r.runs.progressErrCall = tc.call
			run := mustExecute(t, r)
			if run.Status != model.BackfillRunFailed || !strings.HasPrefix(run.ErrorMessage, "storage: ") || !strings.Contains(run.ErrorMessage, "Fortschritt nicht schreibbar") {
				t.Fatalf("Run = %+v, will failed mit storage-Klasse und Ursache", run)
			}
			if run.RowsCopied != tc.rows {
				t.Fatalf("RowsCopied = %d, will %d (zuletzt festgehaltener Zähler)", run.RowsCopied, tc.rows)
			}
			if len(r.runs.progresses) != tc.call {
				t.Fatalf("Fortschritts-Aufrufe = %d, will %d (kein Fortschritt nach dem Fehler)", len(r.runs.progresses), tc.call)
			}
			if len(r.writer.blocks) != tc.appended || r.trace.count("Commit") != 0 || r.trace.count("Notify") != 0 {
				t.Fatalf("Blöcke=%d, Ereignisse %v", len(r.writer.blocks), r.trace.events)
			}
			if r.writer.rollbacks != tc.rollbacks || r.snapshot.closed != 1 {
				t.Fatalf("Rollbacks=%d (will %d), Snapshot %d-mal geschlossen", r.writer.rollbacks, tc.rollbacks, r.snapshot.closed)
			}
			if len(r.runs.finished) != 1 || r.runs.finished[0] != run {
				t.Fatalf("Finish = %+v", r.runs.finished)
			}
		})
	}
}

// TestExecuteEmptyTableFinishFailure trägt den Endzustand der leeren Tabelle:
// `completed` gilt erst, wenn `Finish` ihn festhält. Scheitert `Finish(completed)`,
// endet der Run `failed` (`storage`) — nicht `completed` —, und scheitert auch
// dieses `Finish`, meldet der Aufruf den Fehler mit dem zuletzt festgehaltenen
// Zustand `running`.
func TestExecuteEmptyTableFinishFailure(t *testing.T) {
	finishedStatuses := func(r *rig) []model.BackfillRunStatus {
		var out []model.BackfillRunStatus
		for _, run := range r.runs.finished {
			out = append(out, run.Status)
		}
		return out
	}
	want := []model.BackfillRunStatus{model.BackfillRunCompleted, model.BackfillRunFailed}

	t.Run("erstes Finish scheitert", func(t *testing.T) {
		r := newRig()
		r.snapshot.blocks = nil
		r.runs.finishErr = fmt.Errorf("%w: Endzustand nicht schreibbar", outbound.ErrBackfillStorage)
		r.runs.finishErrCall = 1
		run := mustExecute(t, r)
		if run.Status != model.BackfillRunFailed || !strings.HasPrefix(run.ErrorMessage, "storage: ") || !strings.Contains(run.ErrorMessage, "Endzustand nicht schreibbar") {
			t.Fatalf("Run = %+v, will failed mit storage-Klasse und Ursache", run)
		}
		if got := finishedStatuses(r); !reflect.DeepEqual(got, want) {
			t.Fatalf("Finish-Zustände = %v, will %v", got, want)
		}
		if r.trace.count("Notify") != 0 || r.snapshot.closed != 1 {
			t.Fatalf("Ereignisse %v", r.trace.events)
		}
	})

	t.Run("jedes Finish scheitert", func(t *testing.T) {
		r := newRig()
		r.snapshot.blocks = nil
		r.runs.finishErr = fmt.Errorf("%w: Endzustand nicht schreibbar", outbound.ErrBackfillStorage)
		result, err := r.execute(context.Background(), t)
		if !stderrors.Is(err, outbound.ErrBackfillStorage) || !strings.Contains(err.Error(), "Endzustand nicht schreibbar") {
			t.Fatalf("Execute = %v, will den Persistenzfehler", err)
		}
		if result.Run.Status != model.BackfillRunRunning {
			t.Fatalf("Ergebnis-Run = %+v, will den zuletzt festgehaltenen Zustand running", result.Run)
		}
		if got := finishedStatuses(r); !reflect.DeepEqual(got, want) {
			t.Fatalf("Finish-Zustände = %v, will %v", got, want)
		}
		if r.trace.count("Notify") != 0 {
			t.Fatalf("Wecksignal ohne festgehaltenen Endzustand")
		}
	})
}

// TestExecuteFinishedAt trägt `finished_at` in jedem Endzustand (`SPEC-029`,
// `ADR-0113` Festlegung 3: die Kopierdauer ist `finished_at − started_at`):
// der Zeitpunkt ist die letzte Lesung der Uhr, die je Aufruf weiterläuft, und
// steht sowohl im Ergebnis als auch an dem Port, der den Endzustand festhält
// (Commit oder `Finish`). Ein gestarteter Run trägt ihn nach `started_at`.
func TestExecuteFinishedAt(t *testing.T) {
	cases := []struct {
		name    string
		setup   func(*rig, context.CancelFunc)
		status  model.BackfillRunStatus
		started bool
		held    func(*rig) model.BackfillRun
	}{
		{"completed mit Daten", func(*rig, context.CancelFunc) {}, model.BackfillRunCompleted, true,
			func(r *rig) model.BackfillRun { return r.writer.committed[0] }},
		{"completed, leere Tabelle", func(r *rig, _ context.CancelFunc) { r.snapshot.blocks = nil }, model.BackfillRunCompleted, true,
			func(r *rig) model.BackfillRun { return r.runs.finished[0] }},
		{"failed beim Kopieren", func(r *rig, _ context.CancelFunc) {
			r.snapshot.nextErrAt, r.snapshot.nextErr = 2, fmt.Errorf("%w: Verbindung weg", outbound.ErrSnapshotStorage)
		}, model.BackfillRunFailed, true,
			func(r *rig) model.BackfillRun { return r.runs.finished[0] }},
		{"failed bei der erneuten Prüfung", func(r *rig, _ context.CancelFunc) { r.activation.published = false }, model.BackfillRunFailed, false,
			func(r *rig) model.BackfillRun { return r.runs.finished[0] }},
		{"interrupted", func(r *rig, cancel context.CancelFunc) {
			r.snapshot.onNext = func(call int) {
				if call == 2 {
					cancel()
				}
			}
		}, model.BackfillRunInterrupted, true,
			func(r *rig) model.BackfillRun { return r.runs.finished[0] }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := newRig()
			r.clock.at, r.clock.step = 100, 10
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			tc.setup(r, cancel)
			result, err := r.service.Execute(ctx, backfill.BackfillExecuteCommand{Run: queuedRun(t), Publication: testPublication})
			if err != nil {
				t.Fatalf("Execute = %v", err)
			}
			run := result.Run
			if run.Status != tc.status {
				t.Fatalf("Status = %q, will %q (%+v)", run.Status, tc.status, run)
			}
			if run.FinishedAt.Unset() || run.FinishedAt != r.clock.last {
				t.Fatalf("finished_at = %v, will die letzte Lesung der Uhr %v", run.FinishedAt, r.clock.last)
			}
			if tc.started && !run.FinishedAt.After(run.StartedAt) {
				t.Fatalf("finished_at %v nicht nach started_at %v", run.FinishedAt, run.StartedAt)
			}
			if held := tc.held(r); held.FinishedAt != run.FinishedAt {
				t.Fatalf("der Port hielt finished_at %v, das Ergebnis trägt %v", held.FinishedAt, run.FinishedAt)
			}
		})
	}
}

// --- Execute: Kontext ---------------------------------------------------------

// TestExecuteContextEndedInterrupts trägt die Unterscheidung `interrupted`
// gegen `failed` am eigenen Kontext: endet der Kontext des Aufrufs während
// des Laufs, endet der Run `interrupted` (ohne Fehlertext), gleich ob der
// Fehler aus dem Lesen, dem Schreiben oder dem Öffnen kommt; derselbe
// `transient`-Fehler bei lebendem Kontext endet `failed` (siehe
// `TestExecuteClassifiesFailures`). Snapshot und Transaktion sind beendet, der
// Endzustand geht über einen vom Abbruch gelösten Kontext, es entsteht keine
// Zeile.
func TestExecuteContextEndedInterrupts(t *testing.T) {
	cases := []struct {
		name  string
		setup func(*rig, context.CancelFunc)
	}{
		{"beim Lesen des zweiten Blocks", func(r *rig, cancel context.CancelFunc) {
			r.snapshot.onNext = func(call int) {
				if call == 2 {
					cancel()
				}
			}
		}},
		{"beim Schreiben", func(r *rig, cancel context.CancelFunc) {
			r.writer.appendErr, r.writer.appendErrN = fmt.Errorf("%w: Kontext beendet", outbound.ErrBackfillStorage), 2
			r.writer.onAppend = func(call int) {
				if call == 2 {
					cancel()
				}
			}
		}},
		{"beim Öffnen des Snapshots", func(r *rig, cancel context.CancelFunc) {
			r.snapshotP.openErr = fmt.Errorf("%w: Kontext beendet", outbound.ErrSnapshotTransient)
			r.snapshotP.beforeOpen = cancel
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := newRig()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			tc.setup(r, cancel)

			result, err := r.service.Execute(ctx, backfill.BackfillExecuteCommand{Run: queuedRun(t), Publication: testPublication})
			if err != nil {
				t.Fatalf("Execute = %v, will nil", err)
			}
			run := result.Run
			if run.Status != model.BackfillRunInterrupted || run.ErrorMessage != "" {
				t.Fatalf("Run = %+v, will interrupted ohne Fehlertext", run)
			}
			if len(r.runs.finished) != 1 || r.runs.finished[0] != run {
				t.Fatalf("Finish = %+v", r.runs.finished)
			}
			if r.runs.finishCtxErr[0] != nil {
				t.Fatalf("Finish lief auf dem beendeten Kontext: %v", r.runs.finishCtxErr[0])
			}
			if len(r.writer.committed) != 0 || r.trace.count("Commit") != 0 || r.trace.count("Notify") != 0 {
				t.Fatalf("Commit/Wecksignal nach dem Abbruch: %v", r.trace.events)
			}
			if r.writer.begun > 0 {
				if r.writer.rollbacks != 1 || r.writer.rollCtx[0] != nil {
					t.Fatalf("Rollbacks=%d auf Kontext-Fehler %v", r.writer.rollbacks, r.writer.rollCtx)
				}
			}
			if r.snapshotP.openErr == nil && (r.snapshot.closed != 1 || r.snapshot.closeCtx[0] != nil) {
				t.Fatalf("Snapshot geschlossen %d-mal auf Kontext-Fehler %v", r.snapshot.closed, r.snapshot.closeCtx)
			}
		})
	}
}

// TestExecuteContextEndedBeforeStartLeavesQueued trägt die Aufnahme nach dem
// Neustart (`ADR-0113` Festlegung 2): endet der Kontext, bevor der Run
// `running` ist, bleibt er `queued`, es wird nichts geschrieben, und der Aufruf
// meldet den Kontext-Fehler.
func TestExecuteContextEndedBeforeStartLeavesQueued(t *testing.T) {
	t.Run("vor dem Aufruf", func(t *testing.T) {
		r := newRig()
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		result, err := r.service.Execute(ctx, backfill.BackfillExecuteCommand{Run: queuedRun(t), Publication: testPublication})
		if !stderrors.Is(err, context.Canceled) || result.Run.Status != model.BackfillRunQueued || len(r.trace.events) != 0 {
			t.Fatalf("Execute = %+v, %v, Aufrufe %v", result, err, r.trace.events)
		}
	})
	t.Run("bei der erneuten Vorbedingungs-Prüfung", func(t *testing.T) {
		r := newRig()
		ctx, cancel := context.WithCancel(context.Background())
		r.activation.registeredFn = func(int) (model.SourceTable, bool, error) {
			cancel()
			return model.SourceTable{}, false, context.Canceled
		}
		result, err := r.service.Execute(ctx, backfill.BackfillExecuteCommand{Run: queuedRun(t), Publication: testPublication})
		if !stderrors.Is(err, context.Canceled) || result.Run.Status != model.BackfillRunQueued {
			t.Fatalf("Execute = %+v, %v", result, err)
		}
		if len(r.runs.finished) != 0 || r.snapshotP.openCalls != 0 {
			t.Fatalf("Zustand geschrieben oder Snapshot geöffnet: %v", r.trace.events)
		}
	})
}

// TestExecuteFinishFailureIsReported trägt die Grenze des Ergebnisses: kann der
// Endzustand nicht festgehalten werden, meldet der Aufruf einen Fehler, der die
// Ursache des Runs nennt; das Ergebnis trägt den zuletzt festgehaltenen Zustand.
func TestExecuteFinishFailureIsReported(t *testing.T) {
	r := newRig()
	r.snapshot.nextErrAt, r.snapshot.nextErr = 2, fmt.Errorf("%w: Verbindung weg", outbound.ErrSnapshotStorage)
	r.runs.finishErr = fmt.Errorf("%w: Verbindung zur Run-Tabelle weg", outbound.ErrBackfillStorage)
	result, err := r.execute(context.Background(), t)
	if !stderrors.Is(err, outbound.ErrBackfillStorage) || !strings.Contains(err.Error(), "Verbindung weg") {
		t.Fatalf("Execute = %v, will den Persistenzfehler und die Ursache des Runs", err)
	}
	if result.Run.Status != model.BackfillRunRunning || result.Run.RowsCopied != 2 {
		t.Fatalf("Ergebnis-Run = %+v, will den zuletzt festgehaltenen Zustand (running, 2 Zeilen)", result.Run)
	}
}

// TestExecuteCleanupFailuresAreLoggedNotFatal trägt die Fehlerpfade, die den
// Run nicht beenden: schlagen Rollback und Schließen des Snapshots fehl, bleibt
// der Run `failed` mit der Ursache des ersten Fehlers, und beide Fehler landen im
// Log.
func TestExecuteCleanupFailuresAreLoggedNotFatal(t *testing.T) {
	r := newRig()
	r.snapshot.nextErrAt, r.snapshot.nextErr = 2, fmt.Errorf("%w: Verbindung weg", outbound.ErrSnapshotStorage)
	r.snapshot.closeErr = stderrors.New("Close fehlgeschlagen")
	r.writer.rollbackErr = stderrors.New("Rollback fehlgeschlagen")
	run := mustExecute(t, r)
	if run.Status != model.BackfillRunFailed || !strings.HasPrefix(run.ErrorMessage, "storage: ") || !strings.Contains(run.ErrorMessage, "Verbindung weg") {
		t.Fatalf("Run = %+v", run)
	}
	if len(r.log.warns) != 2 {
		t.Fatalf("Warnungen = %v, will Rollback und Schließen", r.log.warns)
	}
}

// --- Warnungen (`ADR-0113` Festlegung 3) -----------------------------------

// scriptClock liefert die Lesungen der Reihe nach und danach die letzte.
type scriptClock struct {
	readings []int64
	calls    int
}

func (c *scriptClock) Now() model.TimePoint {
	i := c.calls
	c.calls++
	if i >= len(c.readings) {
		i = len(c.readings) - 1
	}
	return model.NewTimePoint(c.readings[i])
}

func (r *rig) executeWithClock(ctx context.Context, t *testing.T, clock outbound.ClockPort) (backfill.BackfillExecuteResult, error) {
	t.Helper()
	service := backfill.NewBackfillTableService(backfill.Ports{
		Activation: r.activation, Exclusion: r.exclusion, Schemas: r.schemas, Snapshot: r.snapshotP,
		Admission: r.admission, Runs: r.runs, Writer: r.writer, Clock: clock,
	}, backfill.WithChangeNotification(r.notifier), backfill.WithLog(r.log))
	return service.Execute(ctx, backfill.BackfillExecuteCommand{Run: queuedRun(t), Publication: testPublication})
}

// TestRequestWarnsAtEstimatedSizeGuideline trägt Warnung (1): die
// **geschätzte** Zeilenzahl über der Richtgröße setzt die Kennzeichnung im Run,
// den `Admit` erhält; genau auf der Richtgröße, bei einer bekannten `0` und bei
// einer unbekannten Schätzung bleibt sie aus — auch wenn die Zahl neben
// „unbekannt“ groß ist. Kein Fall lehnt den Antrag ab: `Admit` läuft in jedem
// als letzter Schritt.
func TestRequestWarnsAtEstimatedSizeGuideline(t *testing.T) {
	guideline := backfill.EstimatedRowsGuidelineForTest
	cases := []struct {
		name  string
		rows  int64
		known bool
		want  bool
	}{
		{"genau auf der Richtgröße", guideline, true, false},
		{"eine Zeile über der Richtgröße", guideline + 1, true, true},
		{"weit über der Richtgröße", guideline * 10, true, true},
		{"unbekannt neben großer Zahl", guideline * 10, false, false},
		{"bekannte Null", 0, true, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := newRig()
			r.snapshotP.estimateRows, r.snapshotP.estimateKnown = tc.rows, tc.known
			result, err := r.service.Request(context.Background(), requestCommand())
			if err != nil {
				t.Fatalf("Request = %v, will nil (keine Ablehnung wegen der Größe)", err)
			}
			if r.admission.calls != 1 || r.trace.events[len(r.trace.events)-1] != "Admit" {
				t.Fatalf("Aufrufe = %v, will Admit als letzten Schritt", r.trace.events)
			}
			if r.admission.gotRun.WarnEstimatedSize != tc.want || result.Run.WarnEstimatedSize != tc.want {
				t.Fatalf("Warnung Größe: Admit %t, Ergebnis %t, will %t", r.admission.gotRun.WarnEstimatedSize, result.Run.WarnEstimatedSize, tc.want)
			}
			if r.admission.gotRun.WarnDuration {
				t.Fatal("Warnung Dauer beim Antrag gesetzt")
			}
			if r.admission.gotRun.Status != model.BackfillRunQueued {
				t.Fatalf("Status = %q, will queued", r.admission.gotRun.Status)
			}
		})
	}
}

// TestExecuteWarnsCopyDurationAtTolerance trägt Warnung (2) im Lauf: das
// Fortschritts-Update je Block und der Abschluss werten die Kopierdauer
// (`finished_at − started_at` an der Uhr) gegen die Toleranz aus. Die Lesungen
// der Uhr eines Laufs mit drei Blöcken sind: 1 Beginn, 2 Snapshot-Zeitpunkt,
// 3 Anfangsposition, 4 bis 6 die Blöcke, 7 Abschluss. Darunter und genau auf der
// Toleranz keine Warnung; darüber gesetzt, ab dem Update, das sie sieht, und bis
// zum Commit weitergetragen. Die Warnung ändert weder Status noch Ablauf.
func TestExecuteWarnsCopyDurationAtTolerance(t *testing.T) {
	const t0 = int64(1000)
	tol := backfill.CopyDurationToleranceNanosForTest
	baseline := newRig()
	mustExecute(t, baseline)

	cases := []struct {
		name         string
		readings     []int64
		wantProgress []bool // je RecordProgress-Aufruf: Anfangsposition, Block 1, 2, 3
		wantFinal    bool
	}{
		{"eine Nanosekunde unter der Toleranz", []int64{t0, t0, t0 + tol - 1}, []bool{false, false, false, false}, false},
		{"genau auf der Toleranz", []int64{t0, t0, t0 + tol}, []bool{false, false, false, false}, false},
		{"über der Toleranz ab dem zweiten Block", []int64{t0, t0, t0, t0, t0 + tol + 1}, []bool{false, false, true, true}, true},
		{"über der Toleranz erst beim Abschluss", []int64{t0, t0, t0, t0, t0, t0, t0 + tol + 1}, []bool{false, false, false, false}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := newRig()
			result, err := r.executeWithClock(context.Background(), t, &scriptClock{readings: tc.readings})
			if err != nil {
				t.Fatalf("Execute = %v", err)
			}
			run := result.Run
			if run.Status != model.BackfillRunCompleted || run.RowsCopied != 5 {
				t.Fatalf("Run = %+v, will completed mit 5 Zeilen (die Warnung ändert den Status nicht)", run)
			}
			if !reflect.DeepEqual(r.trace.events, baseline.trace.events) {
				t.Fatalf("Ablauf = %v, will den Ablauf ohne Warnung %v", r.trace.events, baseline.trace.events)
			}
			var got []bool
			for _, p := range r.runs.progresses {
				got = append(got, p.WarnDuration)
			}
			if !reflect.DeepEqual(got, tc.wantProgress) {
				t.Fatalf("Warnung Dauer je Fortschritts-Update = %v, will %v", got, tc.wantProgress)
			}
			if len(r.writer.committed) != 1 || r.writer.committed[0].WarnDuration != tc.wantFinal || run.WarnDuration != tc.wantFinal {
				t.Fatalf("Warnung Dauer beim Abschluss: Commit %+v, Ergebnis %t, will %t", r.writer.committed, run.WarnDuration, tc.wantFinal)
			}
			if run.WarnEstimatedSize {
				t.Fatal("Warnung Größe im Lauf gesetzt")
			}
		})
	}
}

// TestExecuteWarnsCopyDurationOnEmptyTable trägt den Abschluss ohne
// Schreibtransaktion: eine leere Tabelle, deren Kopierdauer die Toleranz
// überschreitet, hält die Warnung im Endzustand fest.
func TestExecuteWarnsCopyDurationOnEmptyTable(t *testing.T) {
	r := newRig()
	r.snapshot.blocks = nil
	t0, tol := int64(1000), backfill.CopyDurationToleranceNanosForTest
	result, err := r.executeWithClock(context.Background(), t, &scriptClock{readings: []int64{t0, t0, t0, t0 + tol + 1}})
	if err != nil || result.Run.Status != model.BackfillRunCompleted {
		t.Fatalf("Execute = %+v, %v", result, err)
	}
	if len(r.runs.finished) != 1 || !r.runs.finished[0].WarnDuration || !result.Run.WarnDuration {
		t.Fatalf("Finish = %+v, will die Warnung Dauer", r.runs.finished)
	}
	if r.trace.count("Begin") != 0 {
		t.Fatalf("Schreibtransaktion für eine leere Tabelle: %v", r.trace.events)
	}
}

// TestExecuteWarnsCopyDurationOnFailureAndInterrupt trägt den Abschluss eines
// Runs, der endet, bevor er fertig ist: `failed` und `interrupted` halten die
// Warnung fest, wenn die Kopierdauer die Toleranz bis dahin überschritten hat;
// ein Run, der vor seinem Beginn `failed` endet (`queued`, `started_at` nicht
// gesetzt), trägt keine Warnung, auch wenn die Uhr weiter als die Toleranz vom
// Nullwert der Startzeit entfernt ist.
func TestExecuteWarnsCopyDurationOnFailureAndInterrupt(t *testing.T) {
	t0, tol := int64(1000), backfill.CopyDurationToleranceNanosForTest
	late := []int64{t0, t0, t0, t0, t0 + tol + 1}
	cases := []struct {
		name     string
		readings []int64
		setup    func(*rig, context.CancelFunc)
		status   model.BackfillRunStatus
		want     bool
	}{
		{"failed nach der Toleranz", late, func(r *rig, _ context.CancelFunc) {
			r.snapshot.nextErrAt, r.snapshot.nextErr = 2, fmt.Errorf("%w: Verbindung weg", outbound.ErrSnapshotStorage)
		}, model.BackfillRunFailed, true},
		{"failed unter der Toleranz", []int64{t0}, func(r *rig, _ context.CancelFunc) {
			r.snapshot.nextErrAt, r.snapshot.nextErr = 2, fmt.Errorf("%w: Verbindung weg", outbound.ErrSnapshotStorage)
		}, model.BackfillRunFailed, false},
		{"interrupted nach der Toleranz", late, func(r *rig, cancel context.CancelFunc) {
			r.snapshot.onNext = func(call int) {
				if call == 2 {
					cancel()
				}
			}
		}, model.BackfillRunInterrupted, true},
		{"failed vor dem Beginn, Uhr weit von null", []int64{t0 + 2*tol}, func(r *rig, _ context.CancelFunc) {
			r.activation.published = false
		}, model.BackfillRunFailed, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := newRig()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			tc.setup(r, cancel)
			result, err := r.executeWithClock(ctx, t, &scriptClock{readings: tc.readings})
			if err != nil {
				t.Fatalf("Execute = %v", err)
			}
			if result.Run.Status != tc.status || result.Run.WarnDuration != tc.want {
				t.Fatalf("Run = %+v, will %s mit Warnung Dauer %t", result.Run, tc.status, tc.want)
			}
			if len(r.runs.finished) != 1 || r.runs.finished[0].WarnDuration != tc.want {
				t.Fatalf("Finish = %+v, will Warnung Dauer %t", r.runs.finished, tc.want)
			}
		})
	}
}
