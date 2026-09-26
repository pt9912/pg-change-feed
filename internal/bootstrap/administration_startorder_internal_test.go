package bootstrap

import (
	"context"
	stderrors "errors"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"testing"
	"time"

	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/mapper"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Whitebox-Tests der Start-Reihenfolge (`runStreamAfterAdministrationPass`,
// `ADR-0112`): die Sequenz Vorlauf, Goroutinen-Start, Stream-Lauf steht an einer
// Stelle, die diese Tests mit Fakes aufrufen. Ein Lauf am Feed-Container zeigt
// das Fehlen einer Race nicht; die Ordnung selbst belegen diese Tests.

const renamedImage = `{"id":"1","renamed_secret":"geheim"}`

// TestRunStreamAfterAdministrationPassAppliesTheRuleRemovalBeforeTheStreamAssemblesTheFirstTransaction
// trägt die Ordnung an dem Antrag, der sie braucht: ein `pending` stehender
// `remove_transformation`-Antrag ist vermerkt (`applied`) und aus der
// `Assembler`-Bindung genommen, wenn der Stream-Start-Fake aufgerufen wird; die
// erste Transaktion, die der Stream-Fake assembliert, trägt die Rohform. Der
// Goroutinen-Start-Fake sieht den Antrag ebenfalls schon vermerkt.
//
// Rot färbende Mutationen (je gesehen): in `runStreamAfterAdministrationPass`
// `runStream(ctx)` vor `processAdministrationRequests(ctx, deps)` stellen — der
// Stream-Fake sieht den Antrag `pending` und das Bild mit dem Zielnamen;
// `startLoop()` vor den Vorlauf stellen — der Goroutinen-Start-Fake sieht den
// Antrag `pending`; den Vorlauf streichen — beide Fakes sehen ihn `pending`.
// Eingabeseite: der Antrag nennt eine Regel, die der Regelstand nicht führt
// (`remove_transformation` gegen einen fremden Namen) — der Antrag endet `failed`,
// die Regel bleibt in der Bindung, und der Test färbt sich rot.
func TestRunStreamAfterAdministrationPassAppliesTheRuleRemovalBeforeTheStreamAssemblesTheFirstTransaction(t *testing.T) {
	set := setRuleRequest("req-set", "geheimname", ruleSpecText("secret", "renamed_secret"))
	remove := removeRuleRequest("req-remove", "geheimname")
	deps, queue, assembler := ruleFixture(t, set, remove)
	// Der Stand vor dem Prozessstart: die Regel ist vermerkt und trägt in der
	// Bindung, das Herausnehmen wartet.
	queue.mu.Lock()
	queue.pending = []model.AdministrationRequest{remove}
	queue.applied = []model.AdministrationRequestID{set.ID}
	queue.mu.Unlock()
	assembler.SetTransformation("public."+ruleTable, mustRenameRule(t, "geheimname", "secret", "renamed_secret"))
	if image := assemblerRowImage(t, assembler, 1, "public", ruleTable); image != renamedImage {
		t.Fatalf("Vorbedingung verletzt: Row Image vor dem Lauf = %s, wollen %s", image, renamedImage)
	}

	var events []string
	var loopSawApplied, streamSawApplied bool
	var streamImage string
	startLoop := func() {
		events = append(events, "loop")
		loopSawApplied = isApplied(queue, "req-remove")
	}
	runStream := func(context.Context) error {
		events = append(events, "stream")
		streamSawApplied = isApplied(queue, "req-remove")
		streamImage = assemblerRowImage(t, assembler, 2, "public", ruleTable)
		return nil
	}

	if err := runStreamAfterAdministrationPass(context.Background(), deps, startLoop, runStream); err != nil {
		t.Fatalf("runStreamAfterAdministrationPass: %v", err)
	}

	if !reflect.DeepEqual(events, []string{"loop", "stream"}) {
		t.Fatalf("Aufruffolge = %v, wollen [loop stream]", events)
	}
	if !loopSawApplied {
		message, _ := failureOf(queue, "req-remove")
		t.Fatalf("der Goroutinen-Start sah den remove_transformation-Antrag nicht vermerkt (Fehlertext %q)", message)
	}
	if !streamSawApplied {
		t.Fatal("der Stream-Start sah den remove_transformation-Antrag nicht vermerkt")
	}
	if streamImage != rawImage {
		t.Fatalf("erste Transaktion des Streams = %s, wollen die Rohform %s — die Regel trug noch in der Bindung", streamImage, rawImage)
	}
}

// TestRunStreamAfterAdministrationPassBindsAnEnabledTableBeforeTheStreamStarts
// trägt dieselbe Ordnung für eine `enable`-Anfrage: die Bindung entsteht vor dem
// Stream-Start (`ADR-0112`) — eine Tabelle, die beim
// Prozessstart per SQL beantragt `pending` steht, erfasst die erste Transaktion
// des Streams. Rot färbende Mutation: der Aufruf `runStream` vor den Vorlauf —
// der Stream-Fake sieht die Tabelle ungebunden. Eingabeseite: der Antrag trägt
// die Art `disable` — die Bindung entsteht nicht, und der Test färbt sich rot.
func TestRunStreamAfterAdministrationPassBindsAnEnabledTableBeforeTheStreamStarts(t *testing.T) {
	const table = "orders_startorder"
	tableID := administrationTableID("public", table)
	deps, queue, _ := ruleFixture(t)
	assembler, err := mapper.NewAssembler("src-admin", map[string]mapper.TableBinding{}, nil)
	if err != nil {
		t.Fatalf("NewAssembler: %v", err)
	}
	deps.assembler = assembler
	deps.activation = &fakeTableActivationPort{registered: map[string]model.SourceTable{
		"public." + table: {ID: tableID, SourceID: "src-admin", Schema: "public", Table: table},
	}}
	deps.schemaStore = &fakeSchemaStorePort{versions: map[model.SourceTableID]model.SchemaVersion{
		tableID: {ID: administrationSchemaVersionID(tableID), SourceTableID: tableID, Version: 1},
	}}
	queue.mu.Lock()
	queue.pending = []model.AdministrationRequest{
		{ID: "req-enable", Source: "src-admin", Schema: "public", Table: table, Kind: model.AdministrationRequestEnable},
	}
	queue.mu.Unlock()

	var boundAtStream bool
	runStream := func(context.Context) error {
		boundAtStream = assemblerCapturesQualified(t, assembler, 1, "public", table)
		return nil
	}

	if err := runStreamAfterAdministrationPass(context.Background(), deps, func() {}, runStream); err != nil {
		t.Fatalf("runStreamAfterAdministrationPass: %v", err)
	}

	if !boundAtStream {
		t.Fatal("die Tabelle war beim Stream-Start nicht gebunden — der enable-Antrag ist zu diesem Zeitpunkt nicht in der Assembler-Bindung nachgetragen")
	}
	if !isApplied(queue, "req-enable") {
		t.Fatal("der enable-Antrag ist nicht vermerkt")
	}
}

// TestRunStreamAfterAdministrationPassKeepsTheStartOnAReadFailureAndReturnsTheStreamOutcome
// trägt zwei Zusagen: ein Lesefehler der Queue im Vorlauf wird protokolliert und
// hält den Start nicht an (Stream-Start und Goroutinen-Start laufen je einmal),
// und der Rückgabewert ist der des Stream-Laufs unverändert. Rot färbende
// Mutationen: `return runStream(ctx)` durch `runStream(ctx); return nil` ersetzen
// — der Ausgang des Streams geht verloren; den Warn-Zweig in
// `processAdministrationRequests` streichen — die Warnung fehlt. Eingabeseite:
// die Queue liest ohne Fehler (`listErr` nil) — kein Warn-Log, und der Test färbt
// sich rot.
func TestRunStreamAfterAdministrationPassKeepsTheStartOnAReadFailureAndReturnsTheStreamOutcome(t *testing.T) {
	deps, queue, _ := ruleFixture(t)
	queue.listErr = stderrors.New("Queue nicht lesbar")
	log := &recordingLog{}
	deps.log = log
	streamOutcome := stderrors.New("Stream-Ausgang")
	var loops, streams int

	err := runStreamAfterAdministrationPass(context.Background(), deps,
		func() { loops++ },
		func(context.Context) error { streams++; return streamOutcome },
	)

	if err != streamOutcome {
		t.Fatalf("Rückgabewert = %v, wollen den Stream-Ausgang %v", err, streamOutcome)
	}
	if loops != 1 || streams != 1 {
		t.Fatalf("Goroutinen-Start %d, Stream-Start %d, wollen je 1 — der Lesefehler hielt den Start an", loops, streams)
	}
	log.mu.Lock()
	warns := log.warns
	log.mu.Unlock()
	if warns != 1 {
		t.Fatalf("Warnungen = %d, wollen 1 (der Lesefehler wird protokolliert)", warns)
	}
}

// blockingListRequestPort liest die Queue nicht, bis der Kontext endet: das Bild
// einer hängenden Datenbank-Abfrage im Vorlauf.
type blockingListRequestPort struct {
	fakeAdministrationRequestPort
}

func (b *blockingListRequestPort) ListPending(ctx context.Context) ([]outbound.PendingAdministrationRequest, error) {
	<-ctx.Done()
	return nil, ctx.Err()
}

var _ outbound.AdministrationRequestPort = (*blockingListRequestPort)(nil)

// TestRunStreamAfterAdministrationPassHoldsTheStreamUntilThePassEndsAndContextCancelEndsIt
// trägt das Risiko des Vorlaufs: eine hängende Abfrage hält den Stream-Start an,
// solange der Kontext lebt, und der Abbruch des Kontexts beendet den Vorlauf —
// der Stream-Lauf startet danach mit dem beendeten Kontext (er schließt seine
// Verbindung selbst). Rot färbende Mutationen: `runStream` vor den Vorlauf stellen
// — der Stream-Fake läuft schon, während die Abfrage hängt; `ctx` im Aufruf von
// `processAdministrationRequests` durch `context.Background()` ersetzen — der
// Abbruch beendet den Vorlauf nicht, der Test läuft in die Frist.
func TestRunStreamAfterAdministrationPassHoldsTheStreamUntilThePassEndsAndContextCancelEndsIt(t *testing.T) {
	deps, _, _ := ruleFixture(t)
	deps.requests = &blockingListRequestPort{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	streamStarted := make(chan struct{})
	var streamCtxErr error
	done := make(chan error, 1)
	go func() {
		done <- runStreamAfterAdministrationPass(ctx, deps, func() {}, func(streamCtx context.Context) error {
			streamCtxErr = streamCtx.Err()
			close(streamStarted)
			return nil
		})
	}()

	select {
	case <-streamStarted:
		t.Fatal("der Stream startete, während der Vorlauf an der Abfrage hing")
	case <-time.After(100 * time.Millisecond):
	}
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Rückgabewert = %v, wollen nil", err)
		}
	case <-time.After(time.Second):
		t.Fatal("der Vorlauf endete nach dem Abbruch des Kontexts nicht innerhalb 1s")
	}
	if streamCtxErr == nil {
		t.Fatal("der Stream startete mit lebendem Kontext, wollen den beendeten Kontext des Abbruchs")
	}
}

// TestRunSourceTextPassesStreamRunOnlyAsArgumentOfTheSequence bindet die
// Aufrufstelle in `Run`, die kein netzloser Lauf erreicht (`Run` braucht die
// Datenbank): im Quelltext des Rumpfes von `Run` ist `stream.Run`
// ausschließlich das vierte Argument von `runStreamAfterAdministrationPass`, nie
// selbst aufgerufen. Grenze: der Test liest die Gestalt des Quelltexts, nicht das
// Verhalten von `Run`; er bindet weder die Reihenfolge der Argumente noch den
// Inhalt von `administration` noch den Rumpf von `startAdministration` (dessen
// Goroutinen-Start deckt am laufenden Prozess `make test-integration`). Rot
// färbende Mutation: den Aufruf durch `stream.Run(streamCtx)` ersetzen.
func TestRunSourceTextPassesStreamRunOnlyAsArgumentOfTheSequence(t *testing.T) {
	run := runFunction(t)
	isStreamRun := func(expr ast.Expr) bool {
		selector, ok := expr.(*ast.SelectorExpr)
		if !ok || selector.Sel.Name != "Run" {
			return false
		}
		receiver, ok := selector.X.(*ast.Ident)
		return ok && receiver.Name == "stream"
	}
	var direct, sequenced int
	ast.Inspect(run.Body, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		if isStreamRun(call.Fun) {
			direct++
		}
		if name, ok := call.Fun.(*ast.Ident); ok && name.Name == "runStreamAfterAdministrationPass" &&
			len(call.Args) == 4 && isStreamRun(call.Args[3]) {
			sequenced++
		}
		return true
	})
	if direct != 0 {
		t.Fatalf("Run ruft stream.Run %d-mal direkt auf, wollen 0 — der Aufruf steht ohne den Vorlauf", direct)
	}
	if sequenced != 1 {
		t.Fatalf("Run übergibt stream.Run %d-mal an runStreamAfterAdministrationPass, wollen 1", sequenced)
	}
}

// runFunction liest die Funktion `Run` aus `wiring.go` (Quelltext-Gestalt, kein
// Lauf).
func runFunction(t *testing.T) *ast.FuncDecl {
	t.Helper()
	file, err := parser.ParseFile(token.NewFileSet(), "wiring.go", nil, 0)
	if err != nil {
		t.Fatalf("wiring.go lesen: %v", err)
	}
	for _, decl := range file.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok && fn.Recv == nil && fn.Name.Name == "Run" {
			return fn
		}
	}
	t.Fatal("wiring.go trägt keine Funktion Run")
	return nil
}

// TestRunSourceTextOrdersReconcileAndWorkerStartBeforeTheSequenceCall bindet die
// Stellung des Vorlaufs zum Backfill in `Run` (`ADR-0113` Festlegung 2): der
// Aufruf des Abgleichs `running` → `interrupted` und der Start des Workers stehen
// im Quelltext von `Run` vor dem Aufruf von `runStreamAfterAdministrationPass`,
// in dem der Vorlauf einen `backfill`-Antrag annimmt — die Annahme sieht keinen
// `running`-Run einer früheren Prozessinstanz. Grenze: die Positionen der
// Aufrufe im Quelltext (bei `runBackfillWorker` die Stelle des `go func()`-Literals,
// in dem er steht), kein Lauf; ein Aufruf in einer toten Verzweigung
// (`if false { … }`) färbt den Test nicht rot. Rot färbende Mutation: den
// Aufruf `reconcileBackfillRuns` hinter den Aufruf von
// `runStreamAfterAdministrationPass` verschieben.
func TestRunSourceTextOrdersReconcileAndWorkerStartBeforeTheSequenceCall(t *testing.T) {
	run := runFunction(t)
	positions := map[string]token.Pos{}
	ast.Inspect(run.Body, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		if name, ok := call.Fun.(*ast.Ident); ok {
			positions[name.Name] = call.Pos()
		}
		return true
	})
	order := []string{"reconcileBackfillRuns", "runBackfillWorker", "runStreamAfterAdministrationPass"}
	for _, name := range order {
		if _, found := positions[name]; !found {
			t.Fatalf("Run ruft %s nicht auf", name)
		}
	}
	for i := 1; i < len(order); i++ {
		if positions[order[i-1]] >= positions[order[i]] {
			t.Fatalf("Run ruft %s nicht vor %s auf", order[i-1], order[i])
		}
	}
}
