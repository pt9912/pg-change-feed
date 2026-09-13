package bootstrap

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Whitebox-Test (`package bootstrap`, nicht `bootstrap_test`): der
// periodische Auslöse-Zug ist ein unexportiertes Verdrahtungsdetail
// (`slice-044`, `ADR-0014`) — der reale Ende-zu-Ende-Beleg (Löschung
// erfolgt/unterbleibt am laufenden Feed-Container) liegt in
// `tools/harness/run-integration-tests.sh` (`make test-integration`);
// dieser Test belegt die Auslöse-/Fehlerbehandlungs-Logik selbst, ohne
// reale PostgreSQL-Instanz.

// fakeRunRetentionUseCase zeichnet jeden `Run`-Aufruf auf und liefert eine
// vorab festgelegte Folge von Ergebnissen/Fehlern — Stellvertreter für
// `*retention.RunRetentionService` über `inbound.RunRetentionUseCase`.
type fakeRunRetentionUseCase struct {
	mu      sync.Mutex
	calls   []inbound.RunRetentionCommand
	results []inbound.RunRetentionResult
	errs    []error
}

func (f *fakeRunRetentionUseCase) Run(_ context.Context, command inbound.RunRetentionCommand) (inbound.RunRetentionResult, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	idx := len(f.calls)
	f.calls = append(f.calls, command)
	var err error
	if idx < len(f.errs) {
		err = f.errs[idx]
	}
	if err != nil {
		return inbound.RunRetentionResult{}, err
	}
	if idx < len(f.results) {
		return f.results[idx], nil
	}
	return inbound.RunRetentionResult{}, nil
}

func (f *fakeRunRetentionUseCase) callCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.calls)
}

func (f *fakeRunRetentionUseCase) call(i int) inbound.RunRetentionCommand {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.calls[i]
}

func waitForCallCount(t *testing.T, useCase *fakeRunRetentionUseCase, want int) {
	t.Helper()
	deadline := time.After(time.Second)
	for {
		if useCase.callCount() >= want {
			return
		}
		select {
		case <-deadline:
			t.Fatalf("erwartete mindestens %d Aufrufe nicht innerhalb 1s (erhalten: %d)", want, useCase.callCount())
		case <-time.After(time.Millisecond):
		}
	}
}

// TestRunRetentionCleanupCallsUseCaseWithSourceAndPolicy belegt, dass jeder
// Tick `Run` mit der konfigurierten Quelle und Policy aufruft — die
// Freigabe je Change trägt der Use Case selbst
// (`retention.RunRetentionService.Run`), diese Schleife nur den Auslöser.
// Rot färbende Mutation: `inbound.RunRetentionCommand{Source: source,
// Policy: policy}` durch eine Zero-Value-Command ersetzen — dann trüge
// jeder aufgezeichnete Aufruf eine leere Quelle statt der übergebenen.
func TestRunRetentionCleanupCallsUseCaseWithSourceAndPolicy(t *testing.T) {
	minAge, err := model.NewDuration(int64(time.Hour))
	if err != nil {
		t.Fatalf("model.NewDuration: %v", err)
	}
	policy, err := model.NewRetentionPolicy(minAge)
	if err != nil {
		t.Fatalf("model.NewRetentionPolicy: %v", err)
	}
	const source model.SourceID = "src-retention-cleanup-test"

	useCase := &fakeRunRetentionUseCase{}
	log := &recordingLog{}
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		defer close(done)
		runRetentionCleanup(ctx, useCase, source, time.Millisecond, policy, log)
	}()

	waitForCallCount(t, useCase, 2)
	cancel()
	<-done

	got := useCase.call(0)
	if got.Source != source {
		t.Fatalf("RunRetentionCommand.Source = %q, wollen %q", got.Source, source)
	}
	if got.Policy != policy {
		t.Fatalf("RunRetentionCommand.Policy = %+v, wollen %+v", got.Policy, policy)
	}
}

// TestRunRetentionCleanupContinuesAfterError belegt die best-effort-Haltung
// (`SPEC-008`): ein Fehler eines Ticks bricht die Schleife nicht ab, der
// nächste Tick ruft `Run` erneut auf, und der Fehler wird protokolliert
// (nicht verworfen). Rot färbende Mutation: `continue` durch `return`
// ersetzen — dann bliebe der Aufrufzähler nach dem ersten (fehlerhaften)
// Tick bei 1 stehen und dieser Test liefe in den 1s-Timeout.
func TestRunRetentionCleanupContinuesAfterError(t *testing.T) {
	minAge, err := model.NewDuration(0)
	if err != nil {
		t.Fatalf("model.NewDuration: %v", err)
	}
	policy, err := model.NewRetentionPolicy(minAge)
	if err != nil {
		t.Fatalf("model.NewRetentionPolicy: %v", err)
	}

	useCase := &fakeRunRetentionUseCase{errs: []error{errors.New("simulierter Bereinigungsfehler")}}
	log := &recordingLog{}
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		defer close(done)
		runRetentionCleanup(ctx, useCase, "src-retention-error-test", time.Millisecond, policy, log)
	}()

	waitForCallCount(t, useCase, 2)
	cancel()
	<-done

	log.mu.Lock()
	warns := log.warns
	log.mu.Unlock()
	if warns < 1 {
		t.Fatalf("Warn-Log-Aufrufe = %d, wollen mindestens 1 (der simulierte Fehler)", warns)
	}
}

// TestRunRetentionCleanupStopsOnContextCancel belegt, dass die Schleife auf
// `ctx.Done()` zurückkehrt, statt weiter zu ticken — dieselbe Eigenschaft
// wie bei `runHeartbeat`/`runWALRetentionCheck`.
func TestRunRetentionCleanupStopsOnContextCancel(t *testing.T) {
	minAge, err := model.NewDuration(0)
	if err != nil {
		t.Fatalf("model.NewDuration: %v", err)
	}
	policy, err := model.NewRetentionPolicy(minAge)
	if err != nil {
		t.Fatalf("model.NewRetentionPolicy: %v", err)
	}

	useCase := &fakeRunRetentionUseCase{}
	log := &recordingLog{}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	done := make(chan struct{})
	go func() {
		defer close(done)
		runRetentionCleanup(ctx, useCase, "src-retention-cancel-test", time.Hour, policy, log)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("runRetentionCleanup ist nach sofortigem Kontext-Abbruch nicht innerhalb 1s zurückgekehrt")
	}
}
