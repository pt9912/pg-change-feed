package bootstrap

import (
	"context"
	stderrors "errors"
	"sync"
	"testing"
	"time"

	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/mapper"
	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Whitebox-Tests (`package bootstrap`) des Backfill-Workers, des
// Start-Abgleichs und des Antragszweigs `backfill` (`ADR-0111` Teilfrage 5,
// `ADR-0113` Festlegung 2): die Fakes tragen die Port-Verträge nach; der
// reale Weg gegen PostgreSQL liegt in `backfill_endtoend_test.go`
// (`make test-store`) und in `administration_roles_internal_test.go`.

// fakeBackfillRunPort trägt einen In-Memory-`BackfillRunPort`: `Queued` liefert
// die Runs im Zustand `queued` in der gehaltenen Reihenfolge (der Vertrag des
// Ports), jede Operation außer `Queued` wird gezählt — der Worker ruft
// keine.
type fakeBackfillRunPort struct {
	mu         sync.Mutex
	runs       []model.BackfillRun
	queuedErr  error
	queuedCall int
	// afterQueued läuft nach der Beantwortung der Lesung Nr. `call`, ohne
	// die Sperre des Fakes.
	afterQueued func(call int)
	otherCalls  []string

	interruptCalls []interruptCall
	interruptCount int
	interruptErr   error
}

type interruptCall struct {
	source model.SourceID
	at     model.TimePoint
}

func (f *fakeBackfillRunPort) Queued(ctx context.Context, source model.SourceID) ([]model.BackfillRun, error) {
	f.mu.Lock()
	f.queuedCall++
	call := f.queuedCall
	hook := f.afterQueued
	if f.queuedErr != nil {
		err := f.queuedErr
		f.mu.Unlock()
		return nil, err
	}
	var queued []model.BackfillRun
	for _, run := range f.runs {
		if run.Source == source && run.Status == model.BackfillRunQueued {
			queued = append(queued, run)
		}
	}
	f.mu.Unlock()
	if hook != nil {
		hook(call)
	}
	return queued, nil
}

func (f *fakeBackfillRunPort) MarkRunning(ctx context.Context, run model.BackfillRun) error {
	f.record("MarkRunning")
	return nil
}

func (f *fakeBackfillRunPort) RecordProgress(ctx context.Context, run model.BackfillRun) error {
	f.record("RecordProgress")
	return nil
}

func (f *fakeBackfillRunPort) Finish(ctx context.Context, run model.BackfillRun) error {
	f.record("Finish")
	return nil
}

func (f *fakeBackfillRunPort) InterruptRunning(ctx context.Context, source model.SourceID, at model.TimePoint) (int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.interruptCalls = append(f.interruptCalls, interruptCall{source: source, at: at})
	return f.interruptCount, f.interruptErr
}

func (f *fakeBackfillRunPort) record(name string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.otherCalls = append(f.otherCalls, name)
}

// setStatus trägt den Zustandswechsel nach, den der reale Use Case über den
// Run-Zustand festhält.
func (f *fakeBackfillRunPort) setStatus(id model.BackfillRunID, status model.BackfillRunStatus) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for i := range f.runs {
		if f.runs[i].ID == id {
			f.runs[i].Status = status
		}
	}
}

func (f *fakeBackfillRunPort) add(run model.BackfillRun) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.runs = append(f.runs, run)
}

func (f *fakeBackfillRunPort) other() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]string(nil), f.otherCalls...)
}

var _ outbound.BackfillRunPort = (*fakeBackfillRunPort)(nil)

// fakeBackfillUseCase trägt einen `BackfillTableUseCase`: `Execute` ruft
// `onExecute` (Zustandswechsel, Blockieren) und liefert `result`/`err`;
// `Request` trägt jede Eingabe zur Prüfung.
type fakeBackfillUseCase struct {
	mu        sync.Mutex
	executed  []model.BackfillRunID
	commands  []inbound.BackfillExecuteCommand
	requested []inbound.BackfillRequestCommand
	onExecute func(command inbound.BackfillExecuteCommand) (model.BackfillRun, error)
	requestFn func(command inbound.BackfillRequestCommand) error
}

func (f *fakeBackfillUseCase) Request(ctx context.Context, command inbound.BackfillRequestCommand) (inbound.BackfillRequestResult, error) {
	f.mu.Lock()
	f.requested = append(f.requested, command)
	f.mu.Unlock()
	if f.requestFn != nil {
		if err := f.requestFn(command); err != nil {
			return inbound.BackfillRequestResult{}, err
		}
	}
	return inbound.BackfillRequestResult{}, nil
}

func (f *fakeBackfillUseCase) Execute(ctx context.Context, command inbound.BackfillExecuteCommand) (inbound.BackfillExecuteResult, error) {
	f.mu.Lock()
	f.executed = append(f.executed, command.Run.ID)
	f.commands = append(f.commands, command)
	f.mu.Unlock()
	if f.onExecute == nil {
		return inbound.BackfillExecuteResult{Run: command.Run}, nil
	}
	run, err := f.onExecute(command)
	return inbound.BackfillExecuteResult{Run: run}, err
}

func (f *fakeBackfillUseCase) executedIDs() []model.BackfillRunID {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]model.BackfillRunID(nil), f.executed...)
}

var _ inbound.BackfillTableUseCase = (*fakeBackfillUseCase)(nil)

func queuedRunOf(t *testing.T, id, table string) model.BackfillRun {
	t.Helper()
	run, err := model.NewQueuedBackfillRun(model.BackfillRunID(id), "src-backfill-worker", "public", table, model.NewTimePoint(1))
	if err != nil {
		t.Fatalf("NewQueuedBackfillRun: %v", err)
	}
	return run
}

// completesRuns liefert eine `onExecute`-Funktion, die den Run als
// abgeschlossen im Run-Zustand festhält, wie es der reale Use Case über
// `Finish` bzw. der Schreiber im Commit tut.
func completesRuns(runs *fakeBackfillRunPort) func(inbound.BackfillExecuteCommand) (model.BackfillRun, error) {
	return func(command inbound.BackfillExecuteCommand) (model.BackfillRun, error) {
		runs.setStatus(command.Run.ID, model.BackfillRunCompleted)
		done := command.Run
		done.Status = model.BackfillRunCompleted
		return done, nil
	}
}

func startWorker(t *testing.T, deps backfillWorkerDeps) (cancel context.CancelFunc, stopped <-chan struct{}) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		defer close(done)
		runBackfillWorker(ctx, deps)
	}()
	t.Cleanup(func() {
		cancel()
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Error("der Backfill-Worker endete nicht nach dem Abbruch seines Kontexts")
		}
	})
	return cancel, done
}

func waitUntil(t *testing.T, what string, condition func() bool) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for !condition() {
		if time.Now().After(deadline) {
			t.Fatalf("Zeitüberschreitung beim Warten auf: %s", what)
		}
		time.Sleep(2 * time.Millisecond)
	}
}

// Der Worker führt beim Start alle `queued`-Runs aus, ohne dass ein
// Wecksignal eintrifft, und zwar in der Reihenfolge, in der der Run-Zustands-Port
// sie liefert (`ADR-0113` Festlegung 2: Aufnahme beim Prozessstart;
// `(requested_at, run_id)` ist die Ordnung des Ports); ein Run im Zustand
// `interrupted`, `completed`, `failed` oder `running` wird nicht aufgenommen.
// Rot färbende Mutationen (je eine): den ersten Lesegang vor dem Warten
// entfernen (der Worker wartet auf ein Signal, das nie kommt) — der Test läuft
// in die Zeitüberschreitung; statt `queued[0]` das letzte Element ausführen —
// die Reihenfolge kehrt sich um.
func TestBackfillWorkerRunsQueuedRunsAtStartInOrderAndSkipsTheOthers(t *testing.T) {
	runs := &fakeBackfillRunPort{}
	for _, id := range []string{"run-a", "run-b", "run-c"} {
		runs.add(queuedRunOf(t, id, "orders_"+id))
	}
	for id, status := range map[string]model.BackfillRunStatus{
		"run-int": model.BackfillRunInterrupted, "run-done": model.BackfillRunCompleted,
		"run-failed": model.BackfillRunFailed, "run-running": model.BackfillRunRunning,
	} {
		run := queuedRunOf(t, id, "orders_"+id)
		run.Status = status
		runs.add(run)
	}
	useCase := &fakeBackfillUseCase{onExecute: completesRuns(runs)}
	wake := newBackfillWake()
	startWorker(t, backfillWorkerDeps{
		runs: runs, useCase: useCase, source: "src-backfill-worker", publication: "cdc_pub",
		wake: wake, retryAfter: time.Hour, log: &recordingLog{},
	})

	waitUntil(t, "drei ausgeführte Runs", func() bool { return len(useCase.executedIDs()) >= 3 })
	got := useCase.executedIDs()
	if len(got) != 3 || got[0] != "run-a" || got[1] != "run-b" || got[2] != "run-c" {
		t.Fatalf("ausgeführte Runs = %v, erwartet [run-a run-b run-c] (nur queued, in der Reihenfolge des Ports)", got)
	}
	for _, command := range useCase.commands {
		if command.Publication != "cdc_pub" {
			t.Fatalf("Execute trägt die Publication %q, erwartet cdc_pub", command.Publication)
		}
	}
}

// Ein Wecksignal nach der Annahme eines Antrags lässt den wartenden Worker den
// neuen Run lesen und ausführen. Rot färbende Mutation: das Warten im
// `select` von `runBackfillWorker` ohne `case <-deps.wake` — der Run bleibt
// unausgeführt.
func TestBackfillWorkerWakesForARunAdmittedAfterItsStart(t *testing.T) {
	runs := &fakeBackfillRunPort{}
	useCase := &fakeBackfillUseCase{onExecute: completesRuns(runs)}
	wake := newBackfillWake()
	startWorker(t, backfillWorkerDeps{
		runs: runs, useCase: useCase, source: "src-backfill-worker", publication: "cdc_pub",
		wake: wake, retryAfter: time.Hour, log: &recordingLog{},
	})
	waitUntil(t, "erste Lesung des Workers", func() bool {
		runs.mu.Lock()
		defer runs.mu.Unlock()
		return runs.queuedCall >= 1
	})
	if n := len(useCase.executedIDs()); n != 0 {
		t.Fatalf("der Worker führte %d Run(s) ohne queued-Zeile aus", n)
	}

	runs.add(queuedRunOf(t, "run-late", "orders_late"))
	signalBackfillWorker(wake)

	waitUntil(t, "Ausführung des spät angenommenen Runs", func() bool { return len(useCase.executedIDs()) == 1 })
	if got := useCase.executedIDs(); got[0] != "run-late" {
		t.Fatalf("ausgeführt: %v, erwartet [run-late]", got)
	}
}

// Ein Run, der während eines anderen Runs angenommen wird, läuft im selben
// Durchgang danach: der Worker liest die `queued`-Zeilen nach jedem Run erneut
// (`ADR-0113` Festlegung 2), auch ohne dass ein Signal ihn weckt. Rot
// färbende Mutation: `drainBackfillQueue` kehrt nach dem ersten Run zurück
// (statt erneut zu lesen) — der zweite Run bleibt aus, weil kein Signal
// kommt.
func TestBackfillWorkerReadsAgainAfterEachRun(t *testing.T) {
	runs := &fakeBackfillRunPort{}
	runs.add(queuedRunOf(t, "run-first", "orders_first"))
	release := make(chan struct{})
	entered := make(chan struct{})
	useCase := &fakeBackfillUseCase{}
	useCase.onExecute = func(command inbound.BackfillExecuteCommand) (model.BackfillRun, error) {
		if command.Run.ID == "run-first" {
			close(entered)
			<-release
		}
		return completesRuns(runs)(command)
	}
	startWorker(t, backfillWorkerDeps{
		runs: runs, useCase: useCase, source: "src-backfill-worker", publication: "cdc_pub",
		wake: newBackfillWake(), retryAfter: time.Hour, log: &recordingLog{},
	})

	select {
	case <-entered:
	case <-time.After(5 * time.Second):
		t.Fatal("der erste Run begann nicht")
	}
	runs.add(queuedRunOf(t, "run-second", "orders_second"))
	close(release)

	waitUntil(t, "Ausführung des zweiten Runs", func() bool { return len(useCase.executedIDs()) == 2 })
	if got := useCase.executedIDs(); got[0] != "run-first" || got[1] != "run-second" {
		t.Fatalf("ausgeführt: %v, erwartet [run-first run-second]", got)
	}
}

// Ein Signal, das zwischen der letzten Lesung und dem Warten eintrifft, geht
// nicht verloren: das Wecksignal trägt Kapazität 1 (`newBackfillWake`) und
// hält das Signal, bis der Worker wartet. Der Fake nimmt den Run und sendet das
// Signal in der zweiten Lesung, nachdem er sie beantwortet hat (der Worker
// steht noch nicht im Warten). Rot färbende Mutation: `newBackfillWake` legt
// einen ungepufferten Kanal an — der nicht blockierende Sender findet keinen
// Empfänger und verwirft das Signal, der zweite Run bleibt aus.
func TestBackfillWorkerKeepsASignalThatArrivesBetweenReadAndWait(t *testing.T) {
	runs := &fakeBackfillRunPort{}
	runs.add(queuedRunOf(t, "run-first", "orders_first"))
	wake := newBackfillWake()
	runs.afterQueued = func(call int) {
		if call == 2 {
			runs.add(queuedRunOf(t, "run-second", "orders_second"))
			signalBackfillWorker(wake)
		}
	}
	useCase := &fakeBackfillUseCase{onExecute: completesRuns(runs)}
	startWorker(t, backfillWorkerDeps{
		runs: runs, useCase: useCase, source: "src-backfill-worker", publication: "cdc_pub",
		wake: wake, retryAfter: time.Hour, log: &recordingLog{},
	})

	waitUntil(t, "Ausführung des zweiten Runs", func() bool { return len(useCase.executedIDs()) == 2 })
	if got := useCase.executedIDs(); got[0] != "run-first" || got[1] != "run-second" {
		t.Fatalf("ausgeführt: %v, erwartet [run-first run-second]", got)
	}
}

// Das Wecksignal blockiert den Sender nie und verschmilzt: drei Signale ohne
// Empfänger lassen genau eines offen, und der Aufruf kehrt jedes Mal zurück
// (die Administrations-Goroutine blockiert nicht, `ADR-0113` Festlegung 2).
// Rot färbende Mutation: den `default`-Zweig in `signalBackfillWorker`
// entfernen — der zweite Aufruf blockiert, der Test läuft in die
// Zeitüberschreitung.
func TestSignalBackfillWorkerNeverBlocksAndCoalesces(t *testing.T) {
	wake := newBackfillWake()
	returned := make(chan struct{})
	go func() {
		defer close(returned)
		for range 3 {
			signalBackfillWorker(wake)
		}
	}()
	select {
	case <-returned:
	case <-time.After(2 * time.Second):
		t.Fatal("signalBackfillWorker blockierte ohne Empfänger")
	}
	if len(wake) != 1 {
		t.Fatalf("offene Signale = %d, erwartet 1 (Signale verschmelzen)", len(wake))
	}
	signalBackfillWorker(nil)
}

// Das Ergebnis von `Execute` ist der Zustand des Aufrufs, nicht der der Zeile
// (`V-1`): meldet ein Aufruf `failed`, schreibt der Worker daraus keinen
// weiteren Endzustand und stellt keinen neuen Antrag — der Run-Zustands-Port
// erhält keinen `MarkRunning`-, `RecordProgress`- oder `Finish`-Aufruf, der
// Use Case keinen `Request`; der Worker protokolliert den Ausgang und nimmt
// den nächsten Run auf. Rot färbende Mutationen (je eine): in
// `drainBackfillQueue` bei `failed` `deps.runs.Finish(...)` aufrufen — die
// Aufruf-Liste des Ports ist nicht leer; `deps.useCase.Request(...)`
// aufrufen — der Use Case trägt einen Antrag.
func TestBackfillWorkerWritesNoEndStateAndFilesNoRequestFromAFailedResult(t *testing.T) {
	runs := &fakeBackfillRunPort{}
	runs.add(queuedRunOf(t, "run-failing", "orders_failing"))
	runs.add(queuedRunOf(t, "run-after", "orders_after"))
	useCase := &fakeBackfillUseCase{}
	useCase.onExecute = func(command inbound.BackfillExecuteCommand) (model.BackfillRun, error) {
		if command.Run.ID == "run-failing" {
			// Die Zeile trägt `completed` (der Commit wirkte), das Ergebnis
			// meldet `failed` (Commit mit unbekanntem Ausgang).
			runs.setStatus(command.Run.ID, model.BackfillRunCompleted)
			failed := command.Run
			failed.Status = model.BackfillRunFailed
			failed.ErrorMessage = "storage: Commit mit unbekanntem Ausgang"
			return failed, nil
		}
		return completesRuns(runs)(command)
	}
	log := &recordingLog{}
	startWorker(t, backfillWorkerDeps{
		runs: runs, useCase: useCase, source: "src-backfill-worker", publication: "cdc_pub",
		wake: newBackfillWake(), retryAfter: time.Hour, log: log,
	})

	waitUntil(t, "beide Runs ausgeführt", func() bool { return len(useCase.executedIDs()) == 2 })
	if calls := runs.other(); len(calls) != 0 {
		t.Fatalf("der Worker rief %v am Run-Zustands-Port auf, erwartet nur Queued", calls)
	}
	useCase.mu.Lock()
	requested := len(useCase.requested)
	useCase.mu.Unlock()
	if requested != 0 {
		t.Fatalf("der Worker stellte %d Antrag/Anträge, erwartet keinen", requested)
	}
	log.mu.Lock()
	defer log.mu.Unlock()
	if log.warns < 1 {
		t.Fatalf("der fehlgeschlagene Run ist nicht protokolliert: %v", log.messages)
	}
}

// Ein Fehler des Aufrufs von `Execute` („Endzustand nicht festgehalten“) beendet
// den Durchgang, ohne zu schleifen; der Worker liest nach `retryAfter` erneut
// und führt den Run dann aus. Ein Lesefehler des Run-Zustands wirkt ebenso.
// Rot färbende Mutation: `drainBackfillQueue` gibt bei einem Fehler des Aufrufs
// `false` statt `true` zurück — der Worker liest nicht nach `retryAfter`
// erneut (der zweite Versuch bleibt aus).
func TestBackfillWorkerRetriesAfterAFailedPass(t *testing.T) {
	runs := &fakeBackfillRunPort{}
	runs.add(queuedRunOf(t, "run-retry", "orders_retry"))
	var calls int
	useCase := &fakeBackfillUseCase{}
	useCase.onExecute = func(command inbound.BackfillExecuteCommand) (model.BackfillRun, error) {
		calls++
		if calls == 1 {
			return command.Run, stderrors.New("Run-Zustand nicht festgehalten")
		}
		return completesRuns(runs)(command)
	}
	startWorker(t, backfillWorkerDeps{
		runs: runs, useCase: useCase, source: "src-backfill-worker", publication: "cdc_pub",
		wake: newBackfillWake(), retryAfter: 20 * time.Millisecond, log: &recordingLog{},
	})

	waitUntil(t, "zweiter Versuch nach retryAfter", func() bool { return len(useCase.executedIDs()) >= 2 })
	got := useCase.executedIDs()
	if got[0] != "run-retry" || got[1] != "run-retry" {
		t.Fatalf("Versuche = %v, erwartet zweimal run-retry", got)
	}

	failing := &fakeBackfillRunPort{queuedErr: stderrors.New("Verbindung verloren")}
	if failed := drainBackfillQueue(context.Background(), backfillWorkerDeps{
		runs: failing, useCase: &fakeBackfillUseCase{}, source: "src-backfill-worker", log: &recordingLog{},
	}); !failed {
		t.Fatal("drainBackfillQueue meldete keinen Fehler bei einem Lesefehler des Run-Zustands")
	}
}

// Bleibt ein Run nach seinem Versuch `queued`, endet der Durchgang (kein
// Schleifen auf derselben Zeile): `Execute` läuft je Run und Durchgang einmal.
// Rot färbende Mutation: die `attempted`-Prüfung in `drainBackfillQueue`
// entfernen — der Worker ruft `Execute` in einer Schleife auf.
func TestBackfillWorkerDoesNotSpinOnARunThatStaysQueued(t *testing.T) {
	runs := &fakeBackfillRunPort{}
	runs.add(queuedRunOf(t, "run-stuck", "orders_stuck"))
	useCase := &fakeBackfillUseCase{onExecute: func(command inbound.BackfillExecuteCommand) (model.BackfillRun, error) {
		return command.Run, nil
	}}
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	failed := drainBackfillQueue(ctx, backfillWorkerDeps{
		runs: runs, useCase: useCase, source: "src-backfill-worker", publication: "cdc_pub", log: &recordingLog{},
	})
	if got := len(useCase.executedIDs()); got != 1 || !failed {
		t.Fatalf("Execute-Aufrufe = %d, Durchgang mit Fehler = %v — erwartet genau ein Aufruf je Durchgang und ein Durchgang mit Fehler, obwohl der Run queued blieb", got, failed)
	}
}

// Der Abgleich beim Prozessstart ruft `InterruptRunning` für die Quelle mit
// der Uhrzeit des `ClockPort`; ein Fehler des Ports geht an den Aufrufer
// (`Run` endet dann vor dem Worker-Start). Rot färbende Mutationen (je
// eine): eine andere Quelle übergeben; den Zeitpunkt `TimePoint{}` statt
// `clock.Now()` übergeben; den Fehler verwerfen.
func TestReconcileBackfillRunsInterruptsRunningRunsOfTheSource(t *testing.T) {
	runs := &fakeBackfillRunPort{interruptCount: 2}
	log := &recordingLog{}
	clock := fixedClock{at: model.NewTimePoint(42)}
	if err := reconcileBackfillRuns(context.Background(), runs, clock, "src-backfill-worker", log); err != nil {
		t.Fatalf("reconcileBackfillRuns: %v", err)
	}
	if len(runs.interruptCalls) != 1 || runs.interruptCalls[0].source != "src-backfill-worker" || runs.interruptCalls[0].at != clock.at {
		t.Fatalf("InterruptRunning-Aufrufe = %+v, erwartet einer für src-backfill-worker zum Zeitpunkt der Uhr", runs.interruptCalls)
	}
	if log.warns != 1 {
		t.Fatalf("der Abgleich von 2 Runs ist nicht protokolliert: %v", log.messages)
	}

	broken := &fakeBackfillRunPort{interruptErr: stderrors.New("Speicher nicht erreichbar")}
	if err := reconcileBackfillRuns(context.Background(), broken, clock, "src-backfill-worker", log); err == nil {
		t.Fatal("reconcileBackfillRuns verwarf den Fehler des Ports")
	}
}

type fixedClock struct{ at model.TimePoint }

func (c fixedClock) Now() model.TimePoint { return c.at }

var _ outbound.ClockPort = fixedClock{}

// Der Antragszweig `backfill` ruft `Request` mit den Werten des Antrags und der
// Publication der Verdrahtung und weckt danach den Worker genau einmal;
// die Kopie läuft nicht in dieser Goroutine (kein `Execute`-Aufruf). Ein
// Fehler von `Request` (nicht aktivierte Tabelle, aktiver Run) endet den
// Antrag als `failed` samt Text und weckt den Worker nicht. Rot färbende
// Mutationen (je eine): die Publication nicht übergeben; Schema und Tabelle
// vertauschen; das Wecksignal auch bei einem Fehler senden; `Execute` im
// Zweig aufrufen.
func TestAdministrationBackfillBranchRequestsAndWakesTheWorker(t *testing.T) {
	ctx := context.Background()
	useCase := &fakeBackfillUseCase{}
	wake := newBackfillWake()
	requests := &fakeAdministrationRequestPort{pending: []model.AdministrationRequest{
		{ID: "req-ok", Source: "src-admin", Schema: "sales", Table: "orders", Kind: model.AdministrationRequestBackfill},
	}}
	deps := administrationDeps{requests: requests, backfill: useCase, backfillWake: wake, publication: "cdc_pub", log: &recordingLog{}}

	processAdministrationRequests(ctx, deps)

	if len(useCase.requested) != 1 {
		t.Fatalf("Request-Aufrufe = %d, erwartet 1", len(useCase.requested))
	}
	got := useCase.requested[0]
	if got.RequestID != "req-ok" || got.Source != "src-admin" || got.Schema != "sales" || got.Table != "orders" || got.Publication != "cdc_pub" {
		t.Fatalf("Request-Eingabe = %+v, erwartet die Werte des Antrags und die Publication cdc_pub", got)
	}
	if len(wake) != 1 {
		t.Fatalf("Wecksignale nach der Annahme = %d, erwartet 1", len(wake))
	}
	if requests.appliedCount() != 1 {
		t.Fatalf("Vermerke applied = %d, erwartet 1 (das nachgelagerte MarkApplied ist wirkungslos, kein Fehler)", requests.appliedCount())
	}

	// Ablehnung: kein Signal, Vermerk failed samt Text.
	<-wake
	useCase.requestFn = func(inbound.BackfillRequestCommand) error {
		return stderrors.New("Tabelle nicht aktiviert")
	}
	requests.pending = []model.AdministrationRequest{
		{ID: "req-refused", Source: "src-admin", Schema: "sales", Table: "orders", Kind: model.AdministrationRequestBackfill},
	}
	processAdministrationRequests(ctx, deps)
	if len(wake) != 0 {
		t.Fatal("ein abgelehnter Antrag weckte den Worker")
	}
	if message := requests.failed["req-refused"]; message != "Tabelle nicht aktiviert" {
		t.Fatalf("Fehlertext des abgelehnten Antrags = %q, erwartet der Text von Request", message)
	}
}

// Ohne verdrahteten Use Case endet ein `backfill`-Antrag als Fehler
// (`errBackfillNotWired`).
func TestAdministrationBackfillBranchWithoutUseCaseFailsTheRequest(t *testing.T) {
	err := applyAdministrationRequest(context.Background(), administrationDeps{log: &recordingLog{}}, model.AdministrationRequest{
		ID: "req-x", Source: "src-admin", Schema: "sales", Table: "orders", Kind: model.AdministrationRequestBackfill,
	})
	if !stderrors.Is(err, errBackfillNotWired) {
		t.Fatalf("Fehler = %v, erwartet errBackfillNotWired", err)
	}
}

// Ein blockierender Run hält die Administrations-Goroutine nicht auf: während
// der Worker in `Execute` steht, verarbeitet `processAdministrationRequests`
// einen weiteren Antrag (`enable`). Rot färbende Mutation: den
// `Execute`-Aufruf in den `backfill`-Zweig verlegen — die Verarbeitung des
// zweiten Antrags wartet auf den Run, der Test läuft in die
// Zeitüberschreitung.
func TestBackfillRunDoesNotBlockTheAdministrationGoroutine(t *testing.T) {
	var enteredOnce sync.Once
	runs := &fakeBackfillRunPort{}
	runs.add(queuedRunOf(t, "run-blocking", "orders_blocking"))
	entered := make(chan struct{})
	release := make(chan struct{})
	useCase := &fakeBackfillUseCase{}
	useCase.onExecute = func(command inbound.BackfillExecuteCommand) (model.BackfillRun, error) {
		enteredOnce.Do(func() { close(entered) })
		<-release
		return completesRuns(runs)(command)
	}
	startWorker(t, backfillWorkerDeps{
		runs: runs, useCase: useCase, source: "src-backfill-worker", publication: "cdc_pub",
		wake: newBackfillWake(), retryAfter: time.Hour, log: &recordingLog{},
	})
	select {
	case <-entered:
	case <-time.After(5 * time.Second):
		t.Fatal("der Run begann nicht")
	}
	t.Cleanup(func() {
		select {
		case <-release:
		default:
			close(release)
		}
	})

	tableID := administrationTableID("public", "orders_enable")
	enableCalls := &fakeEnableTableUseCase{}
	requests := &fakeAdministrationRequestPort{pending: []model.AdministrationRequest{
		{ID: "req-enable", Source: "src-admin", Schema: "public", Table: "orders_enable", Kind: model.AdministrationRequestEnable},
		{ID: "req-backfill", Source: "src-admin", Schema: "public", Table: "orders_enable", Kind: model.AdministrationRequestBackfill},
	}}
	assembler, err := mapper.NewAssembler("src-admin", map[string]mapper.TableBinding{}, nil)
	if err != nil {
		t.Fatalf("NewAssembler: %v", err)
	}
	deps := administrationDeps{
		requests: requests,
		activation: &fakeTableActivationPort{registered: map[string]model.SourceTable{
			"public.orders_enable": {ID: tableID, SourceID: "src-admin", Schema: "public", Table: "orders_enable"},
		}},
		enableTables: enableCalls,
		schemaStore: &fakeSchemaStorePort{versions: map[model.SourceTableID]model.SchemaVersion{
			tableID: {ID: administrationSchemaVersionID(tableID), SourceTableID: tableID, Version: 1},
		}},
		columnExclusion: &fakeColumnExclusionPort{},
		assembler:       assembler,
		backfill:        useCase,
		backfillWake:    newBackfillWake(),
		publication:     "cdc_pub",
		log:             &recordingLog{},
	}
	done := make(chan struct{})
	go func() {
		defer close(done)
		processAdministrationRequests(context.Background(), deps)
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("die Administrations-Goroutine blockierte, während ein Run lief")
	}
	if len(enableCalls.calls) != 1 || requests.appliedCount() != 2 {
		t.Fatalf("enable-Aufrufe %d, Vermerke applied %d — erwartet 1 und 2", len(enableCalls.calls), requests.appliedCount())
	}
}

// Die Ausgabe eines Runs in `diagnose` nennt die Zeilenzahl als **geschätzt**:
// eine bekannte Schätzung als Zahl (auch 0), eine unbekannte als „unbekannt“,
// nie als 0; die beiden Kennzeichnungen stehen als `true`/`false`, ein
// Fehlertext in einer zweiten Zeile. Rot färbende Mutationen (je eine): für
// `estimated == nil` die Zahl `0` ausgeben; die Kennzeichnungen vertauschen;
// den Fehlertext weglassen.
func TestFormatBackfillRunNamesTheEstimateAsEstimatedAndUnknownAsUnknown(t *testing.T) {
	zero, ten := int64(0), int64(10)
	cases := []struct {
		name      string
		estimated *int64
		warnSize  bool
		warnDur   bool
		message   string
		want      string
	}{
		{"bekannt", &ten, false, false, "", "    public.orders: completed, 12 Zeilen kopiert, geschätzt 10, Warnung Größe false, Warnung Dauer false\n"},
		{"bekannt null", &zero, false, false, "", "    public.orders: completed, 12 Zeilen kopiert, geschätzt 0, Warnung Größe false, Warnung Dauer false\n"},
		{"unbekannt", nil, false, false, "", "    public.orders: completed, 12 Zeilen kopiert, geschätzt unbekannt, Warnung Größe false, Warnung Dauer false\n"},
		{"Kennzeichnungen", nil, true, false, "", "    public.orders: completed, 12 Zeilen kopiert, geschätzt unbekannt, Warnung Größe true, Warnung Dauer false\n"},
		{"Dauer", nil, false, true, "", "    public.orders: completed, 12 Zeilen kopiert, geschätzt unbekannt, Warnung Größe false, Warnung Dauer true\n"},
		{"Fehlertext", nil, false, false, "permission: kein SELECT", "    public.orders: completed, 12 Zeilen kopiert, geschätzt unbekannt, Warnung Größe false, Warnung Dauer false\n      Fehler: permission: kein SELECT\n"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := formatBackfillRun("public", "orders", "completed", 12, c.estimated, c.warnSize, c.warnDur, c.message); got != c.want {
				t.Fatalf("Ausgabe = %q, erwartet %q", got, c.want)
			}
		})
	}
}
