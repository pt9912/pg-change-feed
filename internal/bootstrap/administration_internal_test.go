package bootstrap

import (
	"context"
	stderrors "errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/decode"
	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/mapper"
	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/removetransformation"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/settransformation"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Whitebox-Test (`package bootstrap`, nicht `bootstrap_test`): die
// Administrations-Goroutine (`runAdministration`/
// `processAdministrationRequests`/`applyAdministrationRequest`, `ADR-0050`)
// ist ein unexportiertes Verdrahtungsdetail — derselbe Aufbau wie
// `walretention_internal_test.go`/`heartbeat_internal_test.go`. Der reale
// Ende-zu-Ende-Beleg (`cdc.enable_table`/`cdc.disable_table` → laufender
// Feed-Container erfasst/stoppt ohne Neustart) liegt in
// `tools/harness/run-integration-tests.sh` (`make test-integration`); diese
// Datei belegt die Fehlerpfade und die Assembler-Bindungs-Nachtragung ohne
// reale PostgreSQL-Instanz (Review-Finding F-2, Review zu `slice-037`).

// fakeAdministrationRequestPort trägt einen In-Memory-Stub des
// `outbound.AdministrationRequestPort`: `pending` wird bei `ListPending`
// einmalig geliefert und danach geleert — derselbe Konsum-einmal-Vertrag wie
// die reale Antrags-Queue (ein bereits gelesener Antrag bleibt `pending`,
// bis `MarkApplied`/`MarkFailed` ihn vermerkt; dieser Stub bildet nur den
// Lese-Teil nach, den `processAdministrationRequests` je Durchlauf braucht).
type fakeAdministrationRequestPort struct {
	mu      sync.Mutex
	pending []model.AdministrationRequest
	// rows trägt Zeilen mit verworfenen Einträgen in der Ordnung der Queue;
	// ist es gesetzt, liefert `ListPending` es zusätzlich hinter `pending`.
	rows    []outbound.PendingAdministrationRequest
	applied []model.AdministrationRequestID
	failed  map[model.AdministrationRequestID]string
	// marked trägt die Kennungen in der Reihenfolge der erfolgreichen
	// Vermerke, `applied` und `failed` zusammen.
	marked []model.AdministrationRequestID
	// listErr/markAppliedErr/markFailedErr tragen die drei Fehlerzweige
	// des Ports (`processAdministrationRequests`): der Lesefehler beendet
	// den Durchlauf, die beiden Vermerk-Fehler bleiben best-effort und
	// werden protokolliert.
	listErr        error
	markAppliedErr error
	markFailedErr  error
}

func (f *fakeAdministrationRequestPort) ListPending(ctx context.Context) ([]outbound.PendingAdministrationRequest, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.listErr != nil {
		return nil, f.listErr
	}
	pending := make([]outbound.PendingAdministrationRequest, 0, len(f.pending)+len(f.rows))
	for _, request := range f.pending {
		pending = append(pending, outbound.PendingAdministrationRequest{Request: request})
	}
	pending = append(pending, f.rows...)
	f.pending = nil
	f.rows = nil
	return pending, nil
}

func (f *fakeAdministrationRequestPort) MarkApplied(ctx context.Context, id model.AdministrationRequestID) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.markAppliedErr != nil {
		return f.markAppliedErr
	}
	f.applied = append(f.applied, id)
	f.marked = append(f.marked, id)
	return nil
}

func (f *fakeAdministrationRequestPort) MarkFailed(ctx context.Context, id model.AdministrationRequestID, message string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.markFailedErr != nil {
		return f.markFailedErr
	}
	if f.failed == nil {
		f.failed = map[model.AdministrationRequestID]string{}
	}
	f.failed[id] = message
	f.marked = append(f.marked, id)
	return nil
}

func (f *fakeAdministrationRequestPort) appliedCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.applied)
}

var _ outbound.AdministrationRequestPort = (*fakeAdministrationRequestPort)(nil)

// fakeAdministrationListener erfüllt `administrationListener` — dasselbe
// Whitebox-Test-Muster wie `fakeWALRetentionMeasurer`
// (`walretention_internal_test.go`): blockiert bis `ctx` endet (derselbe
// Timeout-Vertrag wie die reale `AdministrationListener.WaitForNotification`,
// deren Aufrufer den Fallback-Poll-Takt als `ctx`-Timeout überträgt) und
// liefert dann `err`, falls gesetzt — sonst `ctx.Err()`.
type fakeAdministrationListener struct {
	err error
}

func (f *fakeAdministrationListener) WaitForNotification(ctx context.Context) error {
	<-ctx.Done()
	if f.err != nil {
		return f.err
	}
	return ctx.Err()
}

var _ administrationListener = (*fakeAdministrationListener)(nil)

// fakeEnableTableUseCase trägt jeden Aufruf zur Prüfung — `err` gesetzt
// bildet den Fehlschlag-Pfad nach (`MarkFailed`-Zweig).
type fakeEnableTableUseCase struct {
	err   error
	calls []inbound.EnableTableCommand
}

func (f *fakeEnableTableUseCase) Enable(ctx context.Context, command inbound.EnableTableCommand) (inbound.EnableTableResult, error) {
	f.calls = append(f.calls, command)
	if f.err != nil {
		return inbound.EnableTableResult{}, f.err
	}
	return inbound.EnableTableResult{Table: model.SourceTable{ID: command.TableID}}, nil
}

var _ inbound.EnableTableUseCase = (*fakeEnableTableUseCase)(nil)

// fakeDisableTableUseCase spiegelt `fakeEnableTableUseCase` für den
// Deaktivierungs-Pfad.
type fakeDisableTableUseCase struct {
	err   error
	calls []inbound.DisableTableCommand
}

func (f *fakeDisableTableUseCase) Disable(ctx context.Context, command inbound.DisableTableCommand) (inbound.DisableTableResult, error) {
	f.calls = append(f.calls, command)
	if f.err != nil {
		return inbound.DisableTableResult{}, f.err
	}
	return inbound.DisableTableResult{Removed: true}, nil
}

var _ inbound.DisableTableUseCase = (*fakeDisableTableUseCase)(nil)

// fakeExcludeColumnUseCase trägt jeden Aufruf zur Prüfung — `err` gesetzt
// bildet den Fehlschlag-Pfad nach (`MarkFailed`-Zweig).
type fakeExcludeColumnUseCase struct {
	err   error
	calls []inbound.ExcludeColumnCommand
}

func (f *fakeExcludeColumnUseCase) Exclude(ctx context.Context, command inbound.ExcludeColumnCommand) error {
	f.calls = append(f.calls, command)
	return f.err
}

var _ inbound.ExcludeColumnUseCase = (*fakeExcludeColumnUseCase)(nil)

// fakeIncludeColumnUseCase spiegelt `fakeExcludeColumnUseCase` für den
// Einschluss-Pfad.
type fakeIncludeColumnUseCase struct {
	err   error
	calls []inbound.IncludeColumnCommand
}

func (f *fakeIncludeColumnUseCase) Include(ctx context.Context, command inbound.IncludeColumnCommand) error {
	f.calls = append(f.calls, command)
	return f.err
}

var _ inbound.IncludeColumnUseCase = (*fakeIncludeColumnUseCase)(nil)

// fakeTableActivationPort trägt `Registered` und `List` mit echtem
// Verhalten — `List` bedient den Bindungs-Neuaufbau
// (`activatedTableBindings`), `Registered` den Enable-Zweig von
// `applyAdministrationRequest` (dort frisch zurückgelesene Bindung,
// `wiring.go`-Kommentar).
type fakeTableActivationPort struct {
	registered map[string]model.SourceTable
	listed     []model.SourceTable
	// registeredErr/listErr tragen die beiden Lese-Fehlerzweige
	// (`applyAdministrationRequest`s Rücklesen der Bindung,
	// `activatedTableBindings`s committed Stand).
	registeredErr error
	listErr       error
}

func (f *fakeTableActivationPort) TableExists(ctx context.Context, schema, table string) (bool, error) {
	return true, nil
}

func (f *fakeTableActivationPort) Registered(ctx context.Context, source model.SourceID, schema, table string) (model.SourceTable, bool, error) {
	if f.registeredErr != nil {
		return model.SourceTable{}, false, f.registeredErr
	}
	entry, found := f.registered[schema+"."+table]
	return entry, found, nil
}

func (f *fakeTableActivationPort) Register(ctx context.Context, table model.SourceTable, version model.SchemaVersion) (bool, error) {
	return true, nil
}

func (f *fakeTableActivationPort) Unregister(ctx context.Context, table model.SourceTable) (outbound.ActivationRemoval, error) {
	return outbound.ActivationRemoved, nil
}

func (f *fakeTableActivationPort) List(ctx context.Context, source model.SourceID) ([]model.SourceTable, error) {
	if f.listErr != nil {
		return nil, f.listErr
	}
	return f.listed, nil
}

func (f *fakeTableActivationPort) Publish(ctx context.Context, publication, schema, table string) error {
	return nil
}

func (f *fakeTableActivationPort) Unpublish(ctx context.Context, publication, schema, table string) error {
	return nil
}

func (f *fakeTableActivationPort) Published(ctx context.Context, publication, schema, table string) (bool, error) {
	return true, nil
}

var _ outbound.TableActivationPort = (*fakeTableActivationPort)(nil)

// fakeSchemaStorePort trägt nur `CurrentVersion` mit echtem Verhalten —
// `applyAdministrationRequest` ruft ausschließlich diese Methode auf.
type fakeSchemaStorePort struct {
	versions map[model.SourceTableID]model.SchemaVersion
	// err trägt den Lese-Fehlerzweig (beide Aufrufer:
	// `applyAdministrationRequest`s Enable-Zweig und
	// `activatedTableBindings`).
	err error
}

func (f *fakeSchemaStorePort) CurrentVersion(ctx context.Context, table model.SourceTableID) (model.SchemaVersion, bool, error) {
	if f.err != nil {
		return model.SchemaVersion{}, false, f.err
	}
	version, found := f.versions[table]
	return version, found, nil
}

func (f *fakeSchemaStorePort) RegisterVersion(ctx context.Context, version model.SchemaVersion, schema model.TableSchema) (bool, error) {
	return true, nil
}

func (f *fakeSchemaStorePort) TableSchema(ctx context.Context, versionID model.SchemaVersionID) (model.TableSchema, error) {
	return model.TableSchema{}, nil
}

var _ outbound.SchemaStorePort = (*fakeSchemaStorePort)(nil)

// fakeColumnExclusionPort trägt den dauerhaften Ausschlussstand als
// In-Memory-Stub des `outbound.ColumnExclusionPort` (`ADR-0065`): beide
// Lesepfade — Bindungs-Neuaufbau und Aktivierungs-Zweig — greifen auf
// dieselbe Rückgabe zu. `ColumnExists` bleibt ungenutzt; die Spaltenprüfung
// des Spalten-Zweigs deckt der reale Adapter-Test ab.
type fakeColumnExclusionPort struct {
	excluded map[string][]string
	err      error
}

func (f *fakeColumnExclusionPort) ColumnExists(ctx context.Context, schema, table, column string) (bool, error) {
	return true, nil
}

func (f *fakeColumnExclusionPort) ExcludedColumns(ctx context.Context, source model.SourceID) (map[string][]string, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.excluded, nil
}

var _ outbound.ColumnExclusionPort = (*fakeColumnExclusionPort)(nil)

// assemblerCapturesQualified belegt den Bindungs-Zustand eines Assemblers
// über seine öffentliche `Consume`-Schnittstelle (kein Zugriff auf
// unexportierte Felder aus einem anderen Paket): eine vollständige
// BEGIN/Change/COMMIT-Folge für `schema.table` trägt genau dann einen
// Change, wenn die Tabelle gebunden ist (`LH-FA-CFG-001`).
func assemblerCapturesQualified(t *testing.T, assembler *mapper.Assembler, xid uint32, schema, table string) bool {
	t.Helper()
	ctx := context.Background()
	rel := &decode.Relation{Schema: schema, Name: table, Columns: []decode.Column{{Name: "id", Key: true}}}
	if _, err := assembler.Consume(ctx, decode.Begin{XID: xid}); err != nil {
		t.Fatalf("Begin: %v", err)
	}
	value := "1"
	if _, err := assembler.Consume(ctx, decode.Change{Relation: rel, Operation: decode.OpInsert, New: []*string{&value}}); err != nil {
		t.Fatalf("Change: %v", err)
	}
	command, err := assembler.Consume(ctx, decode.Commit{CommitLSN: uint64(xid), CommitTime: time.Now()})
	if err != nil {
		t.Fatalf("Commit: %v", err)
	}
	changes, err := command.Transaction.Changes()
	if err != nil {
		t.Fatalf("Changes: %v", err)
	}
	return len(changes) == 1
}

// TestProcessAdministrationRequestsEnableAppliesAndBindsAssembler trägt den
// Happy Path des Enable-Zweigs: der Antrag wird über den Inbound Port
// verarbeitet, als `applied` vermerkt, und die laufende `Assembler`-Bindung
// wird nachgetragen — ohne sie bliebe eine über SQL aktivierte Tabelle im
// laufenden Prozess unerfasst, bis zum nächsten Neustart.
func TestProcessAdministrationRequestsEnableAppliesAndBindsAssembler(t *testing.T) {
	ctx := context.Background()
	const requestID = model.AdministrationRequestID("req-enable-1")
	tableID := administrationTableID("public", "orders_admin_enable")
	requests := &fakeAdministrationRequestPort{pending: []model.AdministrationRequest{
		{ID: requestID, Source: "src-admin", Schema: "public", Table: "orders_admin_enable", Kind: model.AdministrationRequestEnable},
	}}
	activation := &fakeTableActivationPort{registered: map[string]model.SourceTable{
		"public.orders_admin_enable": {ID: tableID, SourceID: "src-admin", Schema: "public", Table: "orders_admin_enable"},
	}}
	schemaStore := &fakeSchemaStorePort{versions: map[model.SourceTableID]model.SchemaVersion{
		tableID: {ID: administrationSchemaVersionID(tableID), SourceTableID: tableID, Version: 1},
	}}
	assembler, err := mapper.NewAssembler("src-admin", map[string]mapper.TableBinding{}, nil)
	if err != nil {
		t.Fatalf("NewAssembler: %v", err)
	}
	enableTables := &fakeEnableTableUseCase{}
	deps := administrationDeps{
		requests:        requests,
		activation:      activation,
		enableTables:    enableTables,
		disableTables:   &fakeDisableTableUseCase{},
		schemaStore:     schemaStore,
		columnExclusion: &fakeColumnExclusionPort{},
		transformations: &fakeTransformationPort{},
		assembler:       assembler,
		publication:     "cdc_pub",
		log:             &recordingLog{},
	}

	processAdministrationRequests(ctx, deps)

	if len(enableTables.calls) != 1 {
		t.Fatalf("Enable-Aufrufe = %d, wollen 1", len(enableTables.calls))
	}
	if got := enableTables.calls[0].TableID; got != tableID {
		t.Fatalf("Enable TableID = %q, wollen %q", got, tableID)
	}
	if requests.appliedCount() != 1 {
		t.Fatalf("applied = %d, wollen 1", requests.appliedCount())
	}
	if !assemblerCapturesQualified(t, assembler, 1, "public", "orders_admin_enable") {
		t.Fatal("Assembler trägt nach Enable keine Bindung für orders_admin_enable — AddBinding hat nicht nachgetragen")
	}
}

// TestProcessAdministrationRequestsDisableAppliesAndUnbindsAssembler trägt
// den Happy Path des Disable-Zweigs — die Gegenseite: eine zuvor gebundene
// Tabelle wird nach der Verarbeitung nicht mehr erfasst.
func TestProcessAdministrationRequestsDisableAppliesAndUnbindsAssembler(t *testing.T) {
	ctx := context.Background()
	const requestID = model.AdministrationRequestID("req-disable-1")
	requests := &fakeAdministrationRequestPort{pending: []model.AdministrationRequest{
		{ID: requestID, Source: "src-admin", Schema: "public", Table: "orders_admin_disable", Kind: model.AdministrationRequestDisable},
	}}
	assembler, err := mapper.NewAssembler("src-admin", map[string]mapper.TableBinding{
		"public.orders_admin_disable": {TableID: "tbl-disable", SchemaVersion: "sv-disable"},
	}, nil)
	if err != nil {
		t.Fatalf("NewAssembler: %v", err)
	}
	if !assemblerCapturesQualified(t, assembler, 1, "public", "orders_admin_disable") {
		t.Fatal("Vorbedingung verletzt: Assembler sollte vor der Deaktivierung gebunden sein")
	}
	disableTables := &fakeDisableTableUseCase{}
	deps := administrationDeps{
		requests:      requests,
		activation:    &fakeTableActivationPort{},
		enableTables:  &fakeEnableTableUseCase{},
		disableTables: disableTables,
		schemaStore:   &fakeSchemaStorePort{},
		assembler:     assembler,
		publication:   "cdc_pub",
		log:           &recordingLog{},
	}

	processAdministrationRequests(ctx, deps)

	if len(disableTables.calls) != 1 {
		t.Fatalf("Disable-Aufrufe = %d, wollen 1", len(disableTables.calls))
	}
	if requests.appliedCount() != 1 {
		t.Fatalf("applied = %d, wollen 1", requests.appliedCount())
	}
	if assemblerCapturesQualified(t, assembler, 2, "public", "orders_admin_disable") {
		t.Fatal("Assembler trägt nach Disable weiterhin eine Bindung für orders_admin_disable — RemoveBinding hat nicht nachgetragen")
	}
}

// TestProcessAdministrationRequestsExcludeColumnAppliesWithoutBinding trägt
// den Happy Path des Spaltenausschlusses (`LH-FA-CFG-005`, `ADR-0059`): der
// Antrag läuft über den Inbound Port mit der Spalte aus dem Antrags-Datensatz
// und wird als `applied` vermerkt — die beiden Spalten-Antragsarten tragen
// keine `Assembler`-Bindung nach, ihr Ziel ist der Filterzustand einer
// bereits getragenen Bindung. Eine nicht gebundene Tabelle bleibt deshalb
// ohne Filterwirkung (`mapper.Assembler.ExcludeColumn`, siehe
// `TestProcessAdministrationRequestsExcludeColumnFiltersAssemblerRowImage`).
func TestProcessAdministrationRequestsExcludeColumnAppliesWithoutBinding(t *testing.T) {
	ctx := context.Background()
	const requestID = model.AdministrationRequestID("req-exclude-1")
	requests := &fakeAdministrationRequestPort{pending: []model.AdministrationRequest{
		{
			ID: requestID, Source: "src-admin", Schema: "public", Table: "orders_admin_exclude",
			Column: "secret", Kind: model.AdministrationRequestExcludeColumn,
		},
	}}
	assembler, err := mapper.NewAssembler("src-admin", map[string]mapper.TableBinding{}, nil)
	if err != nil {
		t.Fatalf("NewAssembler: %v", err)
	}
	excludeColumns := &fakeExcludeColumnUseCase{}
	deps := administrationDeps{
		requests:       requests,
		activation:     &fakeTableActivationPort{},
		enableTables:   &fakeEnableTableUseCase{},
		disableTables:  &fakeDisableTableUseCase{},
		excludeColumns: excludeColumns,
		includeColumns: &fakeIncludeColumnUseCase{},
		schemaStore:    &fakeSchemaStorePort{},
		assembler:      assembler,
		publication:    "cdc_pub",
		log:            &recordingLog{},
	}

	processAdministrationRequests(ctx, deps)

	if len(excludeColumns.calls) != 1 {
		t.Fatalf("Exclude-Aufrufe = %d, wollen 1", len(excludeColumns.calls))
	}
	if got := excludeColumns.calls[0]; got.Table != "orders_admin_exclude" || got.Column != "secret" {
		t.Fatalf("Exclude-Kommando = %+v, wollen Tabelle orders_admin_exclude mit Spalte secret", got)
	}
	if requests.appliedCount() != 1 {
		t.Fatalf("applied = %d, wollen 1", requests.appliedCount())
	}
	if assemblerCapturesQualified(t, assembler, 1, "public", "orders_admin_exclude") {
		t.Fatal("Assembler trägt nach dem Spaltenausschluss eine Bindung für orders_admin_exclude — diese Antragsart trägt keine Bindung nach")
	}
}

// assemblerRowImage liefert das neue Row Image eines erfassten Changes für
// `schema.table` über die öffentliche `Consume`-Schnittstelle (derselbe
// Zugriffsweg wie `assemblerCapturesQualified`): die Relation trägt zwei
// Spalten, das Bild lässt den ausgeschlossenen Schlüssel bei gebundener
// Tabelle und geführtem Ausschluss weg.
func assemblerRowImage(t *testing.T, assembler *mapper.Assembler, xid uint32, schema, table string) string {
	t.Helper()
	ctx := context.Background()
	rel := &decode.Relation{Schema: schema, Name: table, Columns: []decode.Column{
		{Name: "id", Key: true}, {Name: "secret"},
	}}
	if _, err := assembler.Consume(ctx, decode.Begin{XID: xid}); err != nil {
		t.Fatalf("Begin: %v", err)
	}
	id, secret := "1", "geheim"
	if _, err := assembler.Consume(ctx, decode.Change{Relation: rel, Operation: decode.OpInsert, New: []*string{&id, &secret}}); err != nil {
		t.Fatalf("Change: %v", err)
	}
	command, err := assembler.Consume(ctx, decode.Commit{CommitLSN: uint64(xid), CommitTime: time.Now()})
	if err != nil {
		t.Fatalf("Commit: %v", err)
	}
	changes, err := command.Transaction.Changes()
	if err != nil {
		t.Fatalf("Changes: %v", err)
	}
	if len(changes) != 1 {
		t.Fatalf("Change-Anzahl = %d, wollen 1 (Tabelle gebunden)", len(changes))
	}
	return string(changes[0].NewImage)
}

// TestProcessAdministrationRequestsExcludeColumnFiltersAssemblerRowImage
// trägt die Live-Reload-Wirkung des Spaltenausschlusses über die
// Verdrahtung (`LH-FA-CFG-005`, `ADR-0059` Bestätigung): ein verarbeiteter
// `exclude_column`-Antrag trägt den Spaltennamen in den Ausschlussstand der
// laufenden `Assembler`-Bindung nach, ein danach erfasster Change führt ihn
// nicht mehr — ohne diesen Nachtrag bliebe der Antrag `applied` und ohne
// Filterwirkung. Der spiegelbildliche `include_column`-Fall nimmt denselben
// Ausschluss über denselben Weg wieder zurück.
func TestProcessAdministrationRequestsExcludeColumnFiltersAssemblerRowImage(t *testing.T) {
	ctx := context.Background()
	const table = "orders_admin_filter"
	requests := &fakeAdministrationRequestPort{pending: []model.AdministrationRequest{
		{
			ID: "req-exclude-filter-1", Source: "src-admin", Schema: "public", Table: table,
			Column: "secret", Kind: model.AdministrationRequestExcludeColumn,
		},
	}}
	assembler, err := mapper.NewAssembler("src-admin", map[string]mapper.TableBinding{
		"public." + table: {TableID: "tbl-filter", SchemaVersion: "sv-filter"},
	}, nil)
	if err != nil {
		t.Fatalf("NewAssembler: %v", err)
	}
	deps := administrationDeps{
		requests:       requests,
		activation:     &fakeTableActivationPort{},
		enableTables:   &fakeEnableTableUseCase{},
		disableTables:  &fakeDisableTableUseCase{},
		excludeColumns: &fakeExcludeColumnUseCase{},
		includeColumns: &fakeIncludeColumnUseCase{},
		schemaStore:    &fakeSchemaStorePort{},
		assembler:      assembler,
		publication:    "cdc_pub",
		log:            &recordingLog{},
	}

	processAdministrationRequests(ctx, deps)

	if requests.appliedCount() != 1 {
		t.Fatalf("applied = %d, wollen 1", requests.appliedCount())
	}
	excluded := assemblerRowImage(t, assembler, 1, "public", table)
	if excluded != `{"id":"1"}` {
		t.Fatalf("Neu-Image nach dem Spaltenausschluss: %s, wollen ohne den ausgeschlossenen Schlüssel secret", excluded)
	}

	requests.mu.Lock()
	requests.pending = []model.AdministrationRequest{{
		ID: "req-include-filter-1", Source: "src-admin", Schema: "public", Table: table,
		Column: "secret", Kind: model.AdministrationRequestIncludeColumn,
	}}
	requests.mu.Unlock()
	processAdministrationRequests(ctx, deps)

	restored := assemblerRowImage(t, assembler, 2, "public", table)
	if restored != `{"id":"1","secret":"geheim"}` {
		t.Fatalf("Neu-Image nach dem Spalteneinschluss: %s, wollen beide Spalten", restored)
	}
}

// TestProcessAdministrationRequestsMarksFailedWhenColumnUseCaseErrors trägt
// den `MarkFailed`-Zweig des Spalteneinschlusses: eine nicht existierende
// Spalte endet über `ErrSourceColumnMissing` und wird als `failed` samt
// Fehlertext vermerkt (`LH-FA-CFG-005` Negative) — ohne diesen Vermerk
// bliebe der Antrag `pending` und der Fehlschlag unsichtbar.
func TestProcessAdministrationRequestsMarksFailedWhenColumnUseCaseErrors(t *testing.T) {
	ctx := context.Background()
	const requestID = model.AdministrationRequestID("req-include-fail-1")
	wantErr := fmt.Errorf("%w: public.orders_admin_include.missing", inbound.ErrSourceColumnMissing)
	requests := &fakeAdministrationRequestPort{pending: []model.AdministrationRequest{
		{
			ID: requestID, Source: "src-admin", Schema: "public", Table: "orders_admin_include",
			Column: "missing", Kind: model.AdministrationRequestIncludeColumn,
		},
	}}
	assembler, err := mapper.NewAssembler("src-admin", map[string]mapper.TableBinding{}, nil)
	if err != nil {
		t.Fatalf("NewAssembler: %v", err)
	}
	deps := administrationDeps{
		requests:       requests,
		activation:     &fakeTableActivationPort{},
		enableTables:   &fakeEnableTableUseCase{},
		disableTables:  &fakeDisableTableUseCase{},
		excludeColumns: &fakeExcludeColumnUseCase{},
		includeColumns: &fakeIncludeColumnUseCase{err: wantErr},
		schemaStore:    &fakeSchemaStorePort{},
		assembler:      assembler,
		publication:    "cdc_pub",
		log:            &recordingLog{},
	}

	processAdministrationRequests(ctx, deps)

	if requests.appliedCount() != 0 {
		t.Fatalf("applied = %d, wollen 0 (Include ist gescheitert)", requests.appliedCount())
	}
	requests.mu.Lock()
	failedMessage, failedFound := requests.failed[requestID]
	requests.mu.Unlock()
	if !failedFound {
		t.Fatal("MarkFailed wurde nicht aufgerufen")
	}
	if failedMessage != wantErr.Error() {
		t.Fatalf("Fehlertext = %q, wollen %q", failedMessage, wantErr.Error())
	}
}

// TestProcessAdministrationRequestsMarksFailedWhenUseCaseErrors trägt den
// `MarkFailed`-Zweig (Review-Finding F-2): ein gescheiterter Inbound-Port-
// Aufruf wird als `failed` samt Fehlertext vermerkt, nicht als `applied`.
func TestProcessAdministrationRequestsMarksFailedWhenUseCaseErrors(t *testing.T) {
	ctx := context.Background()
	const requestID = model.AdministrationRequestID("req-fail-1")
	wantErr := stderrors.New("Quelltabelle existiert nicht")
	requests := &fakeAdministrationRequestPort{pending: []model.AdministrationRequest{
		{ID: requestID, Source: "src-admin", Schema: "public", Table: "orders_admin_fail", Kind: model.AdministrationRequestEnable},
	}}
	assembler, err := mapper.NewAssembler("src-admin", map[string]mapper.TableBinding{}, nil)
	if err != nil {
		t.Fatalf("NewAssembler: %v", err)
	}
	deps := administrationDeps{
		requests:      requests,
		activation:    &fakeTableActivationPort{},
		enableTables:  &fakeEnableTableUseCase{err: wantErr},
		disableTables: &fakeDisableTableUseCase{},
		schemaStore:   &fakeSchemaStorePort{},
		assembler:     assembler,
		publication:   "cdc_pub",
		log:           &recordingLog{},
	}

	processAdministrationRequests(ctx, deps)

	if requests.appliedCount() != 0 {
		t.Fatalf("applied = %d, wollen 0 (Enable ist gescheitert)", requests.appliedCount())
	}
	requests.mu.Lock()
	failedMessage, failedFound := requests.failed[requestID]
	requests.mu.Unlock()
	if !failedFound {
		t.Fatal("MarkFailed wurde nicht aufgerufen")
	}
	if failedMessage != wantErr.Error() {
		t.Fatalf("Fehlertext = %q, wollen %q", failedMessage, wantErr.Error())
	}
}

// fakeTransformationPort trägt den `outbound.TransformationPort` als
// In-Memory-Stub. Ohne `derivedFrom` liefert er den festen Stand `rules`;
// mit `derivedFrom` leitet er den Regelstand wie der Store aus den
// **vermerkten** (`applied`) Anträgen des Request-Fakes ab — in der Ordnung
// von `history`, gefaltet je Tabelle —, ein noch `pending` stehender oder
// `failed` vermerkter Antrag trägt keinen Stand. `columns` sind die
// Katalog-Spalten je Tabelle.
type fakeTransformationPort struct {
	rules       map[string][]model.Transformation
	columns     map[string][]string
	derivedFrom *fakeAdministrationRequestPort
	history     []model.AdministrationRequest
	err         error
}

func (f *fakeTransformationPort) TransformationRules(ctx context.Context, source model.SourceID) (map[string][]model.Transformation, error) {
	if f.err != nil {
		return nil, f.err
	}
	if f.derivedFrom == nil {
		return f.rules, nil
	}
	f.derivedFrom.mu.Lock()
	applied := map[model.AdministrationRequestID]bool{}
	for _, id := range f.derivedFrom.applied {
		applied[id] = true
	}
	f.derivedFrom.mu.Unlock()
	records := map[string][]model.TransformationRecord{}
	for _, request := range f.history {
		if !applied[request.ID] || request.Source != source {
			continue
		}
		if request.Kind != model.AdministrationRequestSetTransformation && request.Kind != model.AdministrationRequestRemoveTransformation {
			continue
		}
		qualified := request.Schema + "." + request.Table
		records[qualified] = append(records[qualified], model.TransformationRecord{Kind: request.Kind, Name: request.RuleName, Spec: request.RuleSpec})
	}
	state := map[string][]model.Transformation{}
	for qualified, list := range records {
		folded, err := model.FoldTransformations(list)
		if err != nil {
			return nil, err
		}
		if len(folded) > 0 {
			state[qualified] = folded
		}
	}
	return state, nil
}

func (f *fakeTransformationPort) SourceColumns(ctx context.Context, schema, table string) ([]string, error) {
	return f.columns[schema+"."+table], nil
}

var _ outbound.TransformationPort = (*fakeTransformationPort)(nil)

// ruleTable ist die Tabelle der Regel-Tests: die Spalten `id` und `secret`
// entsprechen dem Relation-Aufbau von `assemblerRowImage`.
const ruleTable = "orders_rules"

// ruleColumns sind die Katalog-Spalten von `ruleTable`; `renamed_secret` und
// `note` gehören nur zum Katalog (K3, K2).
var ruleColumns = map[string][]string{"public." + ruleTable: {"id", "secret", "note"}}

func mustRenameRule(t *testing.T, name, column, to string) model.Transformation {
	t.Helper()
	rule, err := model.NewRenameColumn(name, column, to)
	if err != nil {
		t.Fatalf("NewRenameColumn(%q, %q, %q): %v", name, column, to, err)
	}
	return rule
}

func ruleSpecText(column, to string) string {
	return `{"kind": "rename_column", "column": "` + column + `", "to": "` + to + `"}`
}

func setRuleRequest(id, name, spec string) model.AdministrationRequest {
	return model.AdministrationRequest{
		ID: model.AdministrationRequestID(id), Source: "src-admin", Schema: "public", Table: ruleTable,
		RuleName: name, RuleSpec: spec, Kind: model.AdministrationRequestSetTransformation,
	}
}

func removeRuleRequest(id, name string) model.AdministrationRequest {
	return model.AdministrationRequest{
		ID: model.AdministrationRequestID(id), Source: "src-admin", Schema: "public", Table: ruleTable,
		RuleName: name, Kind: model.AdministrationRequestRemoveTransformation,
	}
}

// ruleFixture baut die Verdrahtung der Regel-Tests: die Antrags-Queue mit den
// übergebenen Anträgen, den Regelstand, den der Store aus den vermerkten
// Anträgen ableitet, die beiden realen Use Cases und einen Assembler, dessen
// Tabelle `public.orders_rules` gebunden ist.
func ruleFixture(t *testing.T, requests ...model.AdministrationRequest) (administrationDeps, *fakeAdministrationRequestPort, *mapper.Assembler) {
	t.Helper()
	queue := &fakeAdministrationRequestPort{pending: append([]model.AdministrationRequest(nil), requests...)}
	port := &fakeTransformationPort{derivedFrom: queue, history: requests, columns: ruleColumns}
	assembler, err := mapper.NewAssembler("src-admin", map[string]mapper.TableBinding{
		"public." + ruleTable: {TableID: "tbl-rules", SchemaVersion: "sv-rules"},
	}, nil)
	if err != nil {
		t.Fatalf("NewAssembler: %v", err)
	}
	deps := administrationDeps{
		requests:              queue,
		activation:            &fakeTableActivationPort{},
		enableTables:          &fakeEnableTableUseCase{},
		disableTables:         &fakeDisableTableUseCase{},
		schemaStore:           &fakeSchemaStorePort{},
		columnExclusion:       &fakeColumnExclusionPort{},
		transformations:       port,
		setTransformations:    settransformation.NewSetTransformationService(port),
		removeTransformations: removetransformation.NewRemoveTransformationService(port),
		assembler:             assembler,
		source:                "src-admin",
		publication:           "cdc_pub",
		log:                   &recordingLog{},
	}
	return deps, queue, assembler
}

func failureOf(queue *fakeAdministrationRequestPort, id string) (string, bool) {
	queue.mu.Lock()
	defer queue.mu.Unlock()
	message, found := queue.failed[model.AdministrationRequestID(id)]
	return message, found
}

func isApplied(queue *fakeAdministrationRequestPort, id string) bool {
	queue.mu.Lock()
	defer queue.mu.Unlock()
	for _, applied := range queue.applied {
		if string(applied) == id {
			return true
		}
	}
	return false
}

const rawImage = `{"id":"1","secret":"geheim"}`

// TestApplyAdministrationRequestRejectsKindOutsideTheClosedSet trägt den
// `default`-Zweig: eine Antragsart außerhalb der geschlossenen Menge der
// Domäne endet über einen sichtbaren Fehler, statt still übersprungen zu
// werden. Der Fehlertext nennt die Antragsart des Antrags und alle sieben
// verarbeiteten Antragsarten.
func TestApplyAdministrationRequestRejectsKindOutsideTheClosedSet(t *testing.T) {
	deps, _, _ := ruleFixture(t)
	err := applyAdministrationRequest(context.Background(), deps, model.AdministrationRequest{
		ID: "req-unprocessed", Source: "src-admin", Schema: "public", Table: ruleTable, Kind: "truncate",
	})
	const want = `Antragsart "truncate" gehört nicht zu den verarbeiteten Antragsarten enable/disable/exclude_column/include_column/backfill/set_transformation/remove_transformation`
	if err == nil || err.Error() != want {
		t.Fatalf("Fehler = %v, wollen %q", err, want)
	}
}

// TestApplyAdministrationRequestHandlesEveryKindOfTheClosedSet bindet die
// Aufzählung der verarbeiteten Antragsarten an die Fälle des `switch`: jede
// Art der geschlossenen Menge der Domäne trägt einen Zweig und endet nicht im
// `default`. Rot färbende Mutation: einen `case` aus `applyAdministrationRequest`
// streichen (z. B. `AdministrationRequestSetTransformation`) — die Art endet im
// `default`-Fehler, statt zu wirken (`LH-FA-CFG-007`). Die Grenze: eine Art, die der
// Konstruktor annimmt und `AdministrationRequestKinds` nicht aufzählt, sieht
// dieser Test nicht; der Domänen-Test bindet den Konstruktor an die Aufzählung
// nur in die eine Richtung.
func TestApplyAdministrationRequestHandlesEveryKindOfTheClosedSet(t *testing.T) {
	tableID := administrationTableID("public", ruleTable)
	for _, kind := range model.AdministrationRequestKinds() {
		deps, _, _ := ruleFixture(t)
		deps.activation = &fakeTableActivationPort{registered: map[string]model.SourceTable{
			"public." + ruleTable: {ID: tableID, SourceID: "src-admin", Schema: "public", Table: ruleTable},
		}}
		deps.schemaStore = &fakeSchemaStorePort{versions: map[model.SourceTableID]model.SchemaVersion{
			tableID: {ID: administrationSchemaVersionID(tableID), SourceTableID: tableID, Version: 1},
		}}
		deps.excludeColumns = &fakeExcludeColumnUseCase{}
		deps.includeColumns = &fakeIncludeColumnUseCase{}
		err := applyAdministrationRequest(context.Background(), deps, model.AdministrationRequest{
			ID: "req-kind", Source: "src-admin", Schema: "public", Table: ruleTable,
			Column: "secret", RuleName: "regel", RuleSpec: ruleSpecText("secret", "renamed_secret"), Kind: kind,
		})
		if err != nil && strings.Contains(err.Error(), "gehört nicht zu den verarbeiteten Antragsarten") {
			t.Fatalf("Antragsart %q endet im default-Zweig: %v", kind, err)
		}
	}
}

// TestProcessAdministrationRequestsSetAndRemoveTransformationTakeEffectLive
// trägt den Happy Path über die Verdrahtung (`LH-FA-CFG-007`): ein
// `set_transformation`-Antrag wird `applied`, und die laufende
// `Assembler`-Bindung trägt die Regel ohne Neustart — der nächste Change führt
// den Zielnamen statt der Quellspalte; ein `remove_transformation`-Antrag
// nimmt sie wieder heraus. Rot färbende Mutation: den Nachtrag
// `deps.assembler.SetTransformation` bzw. `RemoveTransformation` streichen — der
// Antrag bleibt `applied`, das Bild bleibt roh bzw. die Regel bleibt.
func TestProcessAdministrationRequestsSetAndRemoveTransformationTakeEffectLive(t *testing.T) {
	ctx := context.Background()
	deps, queue, assembler := ruleFixture(t,
		setRuleRequest("req-set", "geheimname", ruleSpecText("secret", "renamed_secret")),
		removeRuleRequest("req-remove", "geheimname"),
	)
	if image := assemblerRowImage(t, assembler, 1, "public", ruleTable); image != rawImage {
		t.Fatalf("Row Image vor dem Antrag = %s, wollen %s", image, rawImage)
	}

	// Ein Durchlauf je Antrag: das zweite Lesen der Queue sieht den Stand des
	// ersten (`applied`), wie der Fallback-Poll der Goroutine.
	queue.mu.Lock()
	remove := queue.pending[1]
	queue.pending = queue.pending[:1]
	queue.mu.Unlock()
	processAdministrationRequests(ctx, deps)
	if !isApplied(queue, "req-set") {
		message, _ := failureOf(queue, "req-set")
		t.Fatalf("set_transformation nicht applied, Fehlertext %q", message)
	}
	if image := assemblerRowImage(t, assembler, 2, "public", ruleTable); image != `{"id":"1","renamed_secret":"geheim"}` {
		t.Fatalf("Row Image nach set_transformation = %s, wollen den Schlüssel renamed_secret ohne Neustart", image)
	}

	queue.mu.Lock()
	queue.pending = []model.AdministrationRequest{remove}
	queue.mu.Unlock()
	processAdministrationRequests(ctx, deps)
	if !isApplied(queue, "req-remove") {
		message, _ := failureOf(queue, "req-remove")
		t.Fatalf("remove_transformation nicht applied, Fehlertext %q", message)
	}
	if image := assemblerRowImage(t, assembler, 3, "public", ruleTable); image != rawImage {
		t.Fatalf("Row Image nach remove_transformation = %s, wollen %s", image, rawImage)
	}
}

// TestProcessAdministrationRequestsMapValueTakesEffectLive trägt den zweiten
// Regeltyp über dieselbe Verdrahtung (`LH-FA-CFG-007`): ein
// `map_value`-Antrag wird `applied`, die laufende `Assembler`-Bindung trägt den
// abgebildeten Wert ohne Neustart und behält den Schlüssel; ein Wert ohne
// Zuordnung bleibt, und ein zweiter `map_value`-Antrag auf dieselbe Spalte
// endet `failed` mit dem Text von K2. Rot färbende Mutation: der Zweig
// `TransformationMapValue` in `TransformationSpec.Build` entfällt — der Antrag
// endet `failed` statt `applied`.
func TestProcessAdministrationRequestsMapValueTakesEffectLive(t *testing.T) {
	ctx := context.Background()
	deps, queue, assembler := ruleFixture(t,
		setRuleRequest("req-map", "geheimwert", `{"kind": "map_value", "column": "secret", "values": {"geheim": "GEHEIM", "offen": "OFFEN"}}`),
	)
	processAdministrationRequests(ctx, deps)
	if !isApplied(queue, "req-map") {
		message, _ := failureOf(queue, "req-map")
		t.Fatalf("map_value-Antrag nicht applied, Fehlertext %q", message)
	}
	if image := assemblerRowImage(t, assembler, 1, "public", ruleTable); image != `{"id":"1","secret":"GEHEIM"}` {
		t.Fatalf("Row Image nach map_value = %s, wollen den abgebildeten Wert unter dem Schlüssel secret", image)
	}

	queue.mu.Lock()
	queue.pending = []model.AdministrationRequest{
		setRuleRequest("req-again", "zweiter", `{"kind": "map_value", "column": "secret", "values": {"x": "y"}}`),
	}
	queue.mu.Unlock()
	processAdministrationRequests(ctx, deps)
	message, failed := failureOf(queue, "req-again")
	if !failed || message != "Spalte trägt bereits eine Regel: public."+ruleTable+".secret" {
		t.Fatalf("zweiter map_value-Antrag auf dieselbe Spalte: failed = %t, Fehlertext %q, wollen K2", failed, message)
	}
	if image := assemblerRowImage(t, assembler, 2, "public", ruleTable); image != `{"id":"1","secret":"GEHEIM"}` {
		t.Fatalf("Row Image nach dem abgelehnten Antrag = %s, wollen unverändert", image)
	}
}

// TestProcessAdministrationRequestsRuleViolationsFailWithSpecTexts trägt die
// Konfliktfreiheit durch die Verdrahtung: jeder Verstoß endet als `failed` mit
// dem Fehlertext der Spec, der Antrag wird nicht `applied`, und die
// laufende Bindung trägt weiter nur die erste Regel (das Bild bleibt, wie es
// nach ihr war). Die Antragsfolge trägt `applied` zwischen den Anträgen: K1
// prüft gegen vermerkte Anträge.
func TestProcessAdministrationRequestsRuleViolationsFailWithSpecTexts(t *testing.T) {
	ctx := context.Background()
	first := setRuleRequest("req-1", "geheimname", ruleSpecText("secret", "renamed_secret"))
	violations := []struct {
		request model.AdministrationRequest
		text    string
	}{
		{setRuleRequest("req-k1", "geheimname", ruleSpecText("note", "notiz")), "Regelname bereits vergeben: public." + ruleTable + ".geheimname"},
		{setRuleRequest("req-k2", "zweite", ruleSpecText("secret", "anders")), "Spalte trägt bereits eine Regel: public." + ruleTable + ".secret"},
		{setRuleRequest("req-k3a", "dritte", ruleSpecText("note", "renamed_secret")), "Zielname kollidiert mit einer anderen Regel: public." + ruleTable + ".renamed_secret"},
		{setRuleRequest("req-k3b", "vierte", ruleSpecText("note", "id")), "Zielname kollidiert mit einer Spalte der Tabelle: public." + ruleTable + ".id"},
		{setRuleRequest("req-k3c", "fuenfte", ruleSpecText("note", "note")), "Zielname kollidiert mit einer Spalte der Tabelle: public." + ruleTable + ".note"},
		{setRuleRequest("req-k4", "sechste", ruleSpecText("gibt_es_nicht", "x")), "Spalte existiert nicht an der Quelle: public." + ruleTable + ".gibt_es_nicht"},
		{removeRuleRequest("req-k4r", "gibt_es_nicht"), "Regelname nicht geführt: public." + ruleTable + ".gibt_es_nicht"},
		{setRuleRequest("req-kind", "siebte", `{"kind": "explode"}`), "unbekannter Regeltyp: explode"},
		{setRuleRequest("req-key", "achte", `{"kind": "rename_column", "column": "note", "to": "n", "extra": 1}`), "unbekannter Schlüssel in rule_spec: extra"},
		{setRuleRequest("req-form", "neunte", `{"kind": "rename_column", "column": "note"}`), "rule_spec ist ungültig: public." + ruleTable + ".neunte"},
		{setRuleRequest("req-name", "Zehnte", ruleSpecText("note", "n")), "Regelname ist ungültig: public." + ruleTable + ".Zehnte"},
	}
	all := []model.AdministrationRequest{first}
	for _, violation := range violations {
		all = append(all, violation.request)
	}
	deps, queue, assembler := ruleFixture(t, all...)
	queue.mu.Lock()
	queue.pending = []model.AdministrationRequest{first}
	queue.mu.Unlock()
	processAdministrationRequests(ctx, deps)
	if !isApplied(queue, "req-1") {
		t.Fatal("die erste Regel ist nicht applied")
	}
	wantImage := `{"id":"1","renamed_secret":"geheim"}`
	xid := uint32(10)
	for _, violation := range violations {
		queue.mu.Lock()
		queue.pending = []model.AdministrationRequest{violation.request}
		queue.mu.Unlock()
		processAdministrationRequests(ctx, deps)
		message, found := failureOf(queue, string(violation.request.ID))
		if !found || message != violation.text {
			t.Fatalf("Antrag %q: Fehlertext = %q (vermerkt %v), wollen %q", violation.request.ID, message, found, violation.text)
		}
		if isApplied(queue, string(violation.request.ID)) {
			t.Fatalf("Antrag %q ist applied", violation.request.ID)
		}
		xid++
		if image := assemblerRowImage(t, assembler, xid, "public", ruleTable); image != wantImage {
			t.Fatalf("Row Image nach dem abgelehnten Antrag %q = %s, wollen unverändert %s", violation.request.ID, image, wantImage)
		}
	}
}

// TestProcessAdministrationRequestsInvalidRuleRowsDoNotStallTheQueue trägt die
// Stelle der Prüfung (Verarbeiten, nicht Lesen): Zeilen mit leerem
// Regelnamen, leerer Regelform (SQL-NULL), JSON-`null` und einem Wert ohne
// Objekt erreichen die Verarbeitung als Anträge und enden `failed` mit dem
// Fehlertext der Spec, und die gültige Zeile dahinter wird verarbeitet und
// `applied`. Rot färbende Mutation: `CheckRuleName` bzw.
// `ParseTransformationSpec` im Use Case streichen färbt die Fehlertexte rot;
// dass die Zeilen die Verarbeitung als Antrag erreichen, trägt
// `TestReadPendingRequestsCarriesRuleRowsWithEmptyRuleFields`.
func TestProcessAdministrationRequestsInvalidRuleRowsDoNotStallTheQueue(t *testing.T) {
	ctx := context.Background()
	requests := []model.AdministrationRequest{
		setRuleRequest("req-empty-name", "", ruleSpecText("note", "n")),
		setRuleRequest("req-null-spec", "regel_a", ""),
		setRuleRequest("req-json-null", "regel_b", "null"),
		setRuleRequest("req-scalar", "regel_c", "1"),
		removeRuleRequest("req-remove-empty", ""),
		setRuleRequest("req-valid", "geheimname", ruleSpecText("secret", "renamed_secret")),
	}
	deps, queue, assembler := ruleFixture(t, requests...)

	processAdministrationRequests(ctx, deps)

	for id, want := range map[string]string{
		"req-empty-name":   "Regelname ist ungültig: public." + ruleTable + ".",
		"req-null-spec":    "rule_spec ist ungültig: public." + ruleTable + ".regel_a",
		"req-json-null":    "rule_spec ist ungültig: public." + ruleTable + ".regel_b",
		"req-scalar":       "rule_spec ist ungültig: public." + ruleTable + ".regel_c",
		"req-remove-empty": "Regelname ist ungültig: public." + ruleTable + ".",
	} {
		if message, found := failureOf(queue, id); !found || message != want {
			t.Fatalf("Antrag %q: Fehlertext = %q (vermerkt %v), wollen %q", id, message, found, want)
		}
	}
	if !isApplied(queue, "req-valid") {
		t.Fatal("die gültige Zeile hinter den ungültigen ist nicht applied")
	}
	if image := assemblerRowImage(t, assembler, 1, "public", ruleTable); image != `{"id":"1","renamed_secret":"geheim"}` {
		t.Fatalf("Row Image = %s, wollen die Regel der gültigen Zeile", image)
	}
}

// rejectedRow ist eine vom Antrags-Konstruktor verworfene Zeile der Queue.
func rejectedRow(id, message string) outbound.PendingAdministrationRequest {
	return outbound.PendingAdministrationRequest{Rejected: &outbound.RejectedAdministrationRequest{
		ID: model.AdministrationRequestID(id), Message: message,
	}}
}

func validRow(request model.AdministrationRequest) outbound.PendingAdministrationRequest {
	return outbound.PendingAdministrationRequest{Request: request}
}

// TestProcessAdministrationRequestsFailsRejectedRowsInQueueOrder trägt den
// Umgang mit einer vom Antrags-Konstruktor verworfenen Zeile (`SPEC-019`): sie
// endet `failed` mit dem Fehlertext der Lesung, und die Zeilen davor und
// dahinter werden in der Ordnung der Queue verarbeitet — die verworfene Zeile
// wird an ihrer Stelle vermerkt, nicht am Ende. Eine Zeile ohne Kennung lässt
// sich nicht vermerken: sie wird mit einer Warnung übersprungen und hält die
// Zeile dahinter nicht an. Rot färbende Mutationen (Eingabeseite): die
// Schleife bricht an der verworfenen Zeile ab (`continue` durch `return`
// ersetzen) — die Zeilen dahinter bleiben ohne Vermerk; der Text des Vermerks
// wird durch einen festen Text ersetzt — der Fehlertext je Zeile weicht ab;
// die verworfene Zeile wird erst nach der Schleife vermerkt — die Ordnung der
// Vermerke weicht ab.
func TestProcessAdministrationRequestsFailsRejectedRowsInQueueOrder(t *testing.T) {
	ctx := context.Background()
	deps, queue, _ := ruleFixture(t)
	log := deps.log.(*recordingLog)
	queue.rows = []outbound.PendingAdministrationRequest{
		validRow(setRuleRequest("req-1", "regel_a", ruleSpecText("secret", "renamed_secret"))),
		rejectedRow("req-2", "Schemaname ist leer: req-2"),
		validRow(setRuleRequest("req-3", "regel_b", ruleSpecText("note", "renamed_note"))),
		rejectedRow("", "Kennung ist leer: public.orders_rules"),
		rejectedRow("req-5", "Spaltenname ist leer: req-5"),
		validRow(setRuleRequest("req-6", "regel_c", ruleSpecText("id", "renamed_id"))),
	}

	processAdministrationRequests(ctx, deps)

	queue.mu.Lock()
	marked := append([]model.AdministrationRequestID(nil), queue.marked...)
	queue.mu.Unlock()
	if want := []model.AdministrationRequestID{"req-1", "req-2", "req-3", "req-5", "req-6"}; !reflect.DeepEqual(marked, want) {
		t.Fatalf("Vermerke in der Reihenfolge %v, wollen %v (jede Zeile an ihrer Stelle der Queue)", marked, want)
	}
	for id, want := range map[string]string{
		"req-2": "Schemaname ist leer: req-2",
		"req-5": "Spaltenname ist leer: req-5",
	} {
		if message, found := failureOf(queue, id); !found || message != want {
			t.Fatalf("Antrag %q: Fehlertext = %q (vermerkt %v), wollen %q", id, message, found, want)
		}
	}
	for _, id := range []string{"req-1", "req-3", "req-6"} {
		if !isApplied(queue, id) {
			t.Fatalf("die gültige Zeile %q ist nicht applied", id)
		}
	}
	if _, found := failureOf(queue, ""); found {
		t.Fatal("eine Zeile ohne Kennung darf nicht vermerkt werden")
	}
	if !log.contains("WARN", "Zeile ohne Kennung übersprungen") {
		t.Fatalf("die Zeile ohne Kennung trägt keine Warnung: %v", log.messages)
	}
}

// TestProcessAdministrationRequestsRejectedRowSurvivesMarkFailedError trägt den
// Fehler des Vermerks einer verworfenen Zeile: er wird protokolliert und
// bricht den Durchlauf nicht ab, die Zeilen dahinter laufen weiter. Rot
// färbende Mutation (Eingabeseite: die verworfene Zeile vor der gültigen): in
// `processAdministrationRequests` nach der verworfenen Zeile die Funktion
// verlassen statt fortzufahren (`continue` durch `return`) — die Zeile dahinter
// bleibt `pending`.
func TestProcessAdministrationRequestsRejectedRowSurvivesMarkFailedError(t *testing.T) {
	ctx := context.Background()
	deps, queue, _ := ruleFixture(t)
	log := deps.log.(*recordingLog)
	queue.markFailedErr = stderrors.New("Vermerk nicht schreibbar")
	queue.rows = []outbound.PendingAdministrationRequest{
		rejectedRow("req-1", "Tabellenname ist leer: req-1"),
		validRow(setRuleRequest("req-2", "regel_a", ruleSpecText("secret", "renamed_secret"))),
	}

	processAdministrationRequests(ctx, deps)

	if !isApplied(queue, "req-2") {
		t.Fatal("die Zeile hinter dem nicht vermerkten Verwurf ist nicht applied")
	}
	if !log.contains("WARN", "Fehlschlag nicht vermerkt") {
		t.Fatalf("der Fehler des Vermerks ist nicht protokolliert: %v", log.messages)
	}
}

// TestProcessAdministrationRequestsSetTransformationIsIdempotent trägt die
// Wiederholung eines bereits nachgetragenen, noch `pending` stehenden Antrags
// (`ADR-0065`-Muster): scheitert der Vermerk `applied` nach dem Nachtrag,
// verarbeitet der nächste Durchlauf denselben Antrag erneut — K1 prüft nur
// gegen vermerkte Anträge, der Nachtrag ersetzt die Regel nach ihrem Namen —,
// ohne Fehler und ohne Änderung des Bildes. Der Test bindet das Verhalten an
// den Store-Fake, der nur vermerkte Anträge ableitet; am Produktionscode dieses
// Pakets lässt sich die Wiederholung nicht rot färben, weil der Use Case den
// Regelstand nur über den Port liest. Die tragende Bindung ist die Abfrage des
// Stores: `TestAdministrationPathRunsUnderLeastPrivilegeLogins` (`make
// test-store`) färbt sich rot, wenn der Statusfilter `status = 'applied'` von
// `SelectAppliedTransformationRequests` durch `status <> 'failed'` ersetzt
// wird (der noch `pending` stehende eigene Antrag zählte, K1 endete `Regelname
// bereits vergeben`).
func TestProcessAdministrationRequestsSetTransformationIsIdempotent(t *testing.T) {
	ctx := context.Background()
	request := setRuleRequest("req-idem", "geheimname", ruleSpecText("secret", "renamed_secret"))
	deps, queue, assembler := ruleFixture(t, request)
	queue.markAppliedErr = stderrors.New("Vermerk nicht schreibbar")

	processAdministrationRequests(ctx, deps)
	if isApplied(queue, "req-idem") {
		t.Fatal("der Vermerk war gestört, der Antrag darf nicht applied sein")
	}
	want := `{"id":"1","renamed_secret":"geheim"}`
	if image := assemblerRowImage(t, assembler, 1, "public", ruleTable); image != want {
		t.Fatalf("Row Image nach dem ersten Durchlauf = %s, wollen %s", image, want)
	}

	queue.mu.Lock()
	queue.markAppliedErr = nil
	queue.pending = []model.AdministrationRequest{request}
	queue.mu.Unlock()
	processAdministrationRequests(ctx, deps)
	if message, found := failureOf(queue, "req-idem"); found {
		t.Fatalf("die Wiederholung endet failed: %q", message)
	}
	if !isApplied(queue, "req-idem") {
		t.Fatal("die Wiederholung ist nicht applied")
	}
	if image := assemblerRowImage(t, assembler, 2, "public", ruleTable); image != want {
		t.Fatalf("Row Image nach der Wiederholung = %s, wollen unverändert %s", image, want)
	}
}

// TestProcessAdministrationRequestsRuleAgainstTableWithoutBindingIsApplied
// trägt den Antrag gegen eine Tabelle ohne laufende Bindung (`SPEC-019`): er
// endet `applied`, weil er wirkt, sobald die Tabelle erfasst wird; der
// Nachtrag in den Assembler ist wirkungslos und legt keine Bindung an.
func TestProcessAdministrationRequestsRuleAgainstTableWithoutBindingIsApplied(t *testing.T) {
	ctx := context.Background()
	deps, queue, assembler := ruleFixture(t, setRuleRequest("req-unbound", "geheimname", ruleSpecText("secret", "renamed_secret")))
	assembler.RemoveBinding("public." + ruleTable)

	processAdministrationRequests(ctx, deps)

	if !isApplied(queue, "req-unbound") {
		message, _ := failureOf(queue, "req-unbound")
		t.Fatalf("Antrag ohne laufende Bindung nicht applied, Fehlertext %q", message)
	}
	if assemblerCapturesQualified(t, assembler, 1, "public", ruleTable) {
		t.Fatal("der Nachtrag hat eine Bindung angelegt")
	}
}

// TestActivatedTableBindingsCarriesTransformations trägt den Startpfad des
// dauerhaften Regelstandes (`ADR-0112` Teilfrage 6): der Bindungs-Neuaufbau
// liest den Regelstand aus seiner dauerhaften Herkunft und übergibt ihn je
// Tabelle an die Bindung — ohne diesen Schritt liefert ein neu gestarteter
// Prozess die Rohform, obwohl der Antrag `applied` trägt. Der Test führt den
// Neuaufbau bis in die Wirkung: das Row Image der neu gebauten Bindung trägt
// den Zielnamen, die Nachbartabelle bleibt roh. Rot färbende Mutation:
// `Transformations: rules[table.QualifiedName()]` in `activatedTableBindings`
// streichen bzw. den Schlüssel gegen einen festen ersetzen.
func TestActivatedTableBindingsCarriesTransformations(t *testing.T) {
	ctx := context.Background()
	const (
		source  = model.SourceID("src-restart")
		tableID = model.SourceTableID("tbl-restart")
		otherID = model.SourceTableID("tbl-restart-other")
	)
	activation := &fakeTableActivationPort{listed: []model.SourceTable{
		{ID: tableID, SourceID: source, Schema: "public", Table: "orders_restart"},
		{ID: otherID, SourceID: source, Schema: "public", Table: "orders_restart_other"},
	}}
	schemaStore := &fakeSchemaStorePort{versions: map[model.SourceTableID]model.SchemaVersion{
		tableID: {ID: "sv-restart", SourceTableID: tableID, Version: 1},
		otherID: {ID: "sv-restart-other", SourceTableID: otherID, Version: 1},
	}}
	rules := &fakeTransformationPort{rules: map[string][]model.Transformation{
		"public.orders_restart": {mustRenameRule(t, "geheimname", "secret", "renamed_secret")},
	}}

	tables, err := activatedTableBindings(ctx, activation, schemaStore, &fakeColumnExclusionPort{}, rules, source)
	if err != nil {
		t.Fatalf("activatedTableBindings = %v, wollen nil", err)
	}
	if got := tables["public.orders_restart"].Transformations; len(got) != 1 || got[0].Name() != "geheimname" {
		t.Fatalf("Transformations = %v, wollen die Regel geheimname", got)
	}
	if got := tables["public.orders_restart_other"].Transformations; len(got) != 0 {
		t.Fatalf("Transformations der Nachbartabelle = %v, wollen keine", got)
	}

	assembler, err := mapper.NewAssembler(source, tables, nil)
	if err != nil {
		t.Fatalf("NewAssembler: %v", err)
	}
	if image := assemblerRowImage(t, assembler, 1, "public", "orders_restart"); image != `{"id":"1","renamed_secret":"geheim"}` {
		t.Fatalf("Row Image der neu gebauten Bindung = %s, wollen den Zielnamen renamed_secret", image)
	}
	if image := assemblerRowImage(t, assembler, 2, "public", "orders_restart_other"); image != rawImage {
		t.Fatalf("Row Image der Nachbartabelle = %s, wollen die Rohform", image)
	}
}

// TestProcessAdministrationRequestsDisableEnableCycleRestoresTransformations
// trägt den zweiten Auslöser des Stand-Verlusts (`ADR-0112` Teilfrage 6): die
// Deaktivierung entfernt den Bindungs-Eintrag samt Regelstand, der
// Aktivierungs-Zweig legt ihn über die dauerhafte Herkunft neu an. Rot
// färbende Mutation: `Transformations: rules[qualified]` im Aktivierungs-Zweig
// streichen — nach `disable` → `enable` steht eine Bindung ohne Regel.
func TestProcessAdministrationRequestsDisableEnableCycleRestoresTransformations(t *testing.T) {
	ctx := context.Background()
	tableID := administrationTableID("public", ruleTable)
	deps, queue, assembler := ruleFixture(t)
	deps.activation = &fakeTableActivationPort{registered: map[string]model.SourceTable{
		"public." + ruleTable: {ID: tableID, SourceID: "src-admin", Schema: "public", Table: ruleTable},
	}}
	deps.schemaStore = &fakeSchemaStorePort{versions: map[model.SourceTableID]model.SchemaVersion{
		tableID: {ID: administrationSchemaVersionID(tableID), SourceTableID: tableID, Version: 1},
	}}
	deps.transformations = &fakeTransformationPort{rules: map[string][]model.Transformation{
		"public." + ruleTable: {mustRenameRule(t, "geheimname", "secret", "renamed_secret")},
	}}
	assembler.SetTransformation("public."+ruleTable, mustRenameRule(t, "geheimname", "secret", "renamed_secret"))
	if image := assemblerRowImage(t, assembler, 1, "public", ruleTable); image != `{"id":"1","renamed_secret":"geheim"}` {
		t.Fatalf("Vorbedingung verletzt: Row Image mit Regel = %s", image)
	}

	queue.mu.Lock()
	queue.pending = []model.AdministrationRequest{{
		ID: "req-cycle-disable", Source: "src-admin", Schema: "public", Table: ruleTable, Kind: model.AdministrationRequestDisable,
	}}
	queue.mu.Unlock()
	processAdministrationRequests(ctx, deps)
	if assemblerCapturesQualified(t, assembler, 2, "public", ruleTable) {
		t.Fatal("Assembler trägt nach Disable weiterhin eine Bindung")
	}

	queue.mu.Lock()
	queue.pending = []model.AdministrationRequest{{
		ID: "req-cycle-enable", Source: "src-admin", Schema: "public", Table: ruleTable, Kind: model.AdministrationRequestEnable,
	}}
	queue.mu.Unlock()
	processAdministrationRequests(ctx, deps)
	if image := assemblerRowImage(t, assembler, 3, "public", ruleTable); image != `{"id":"1","renamed_secret":"geheim"}` {
		t.Fatalf("Row Image nach dem disable/enable-Zyklus = %s, wollen den Zielnamen renamed_secret (Regelstand aus der Herkunft)", image)
	}
}

// TestProcessAdministrationRequestsMarksFailedWhenRuleStateReadFails trägt den
// Fehlerpfad der dauerhaften Herkunft im Aktivierungs-Zweig: ein Lesefehler
// des Regelstandes endet im `failed`-Vermerk, statt eine Bindung ohne den
// geführten Stand anzulegen.
func TestProcessAdministrationRequestsMarksFailedWhenRuleStateReadFails(t *testing.T) {
	ctx := context.Background()
	tableID := administrationTableID("public", ruleTable)
	wantErr := stderrors.New("Antrags-Historie nicht lesbar")
	deps, queue, assembler := ruleFixture(t)
	assembler.RemoveBinding("public." + ruleTable)
	deps.activation = &fakeTableActivationPort{registered: map[string]model.SourceTable{
		"public." + ruleTable: {ID: tableID, SourceID: "src-admin", Schema: "public", Table: ruleTable},
	}}
	deps.schemaStore = &fakeSchemaStorePort{versions: map[model.SourceTableID]model.SchemaVersion{
		tableID: {ID: administrationSchemaVersionID(tableID), SourceTableID: tableID, Version: 1},
	}}
	deps.transformations = &fakeTransformationPort{err: wantErr}
	queue.mu.Lock()
	queue.pending = []model.AdministrationRequest{{
		ID: "req-rules-read-fail", Source: "src-admin", Schema: "public", Table: ruleTable, Kind: model.AdministrationRequestEnable,
	}}
	queue.mu.Unlock()

	processAdministrationRequests(ctx, deps)

	if isApplied(queue, "req-rules-read-fail") {
		t.Fatal("der Antrag ist applied, obwohl der Regelstand nicht lesbar war")
	}
	if message, found := failureOf(queue, "req-rules-read-fail"); !found || message != wantErr.Error() {
		t.Fatalf("Fehlertext = %q (vermerkt %v), wollen %q", message, found, wantErr.Error())
	}
	if assemblerCapturesQualified(t, assembler, 1, "public", ruleTable) {
		t.Fatal("Assembler trägt nach gescheitertem Regelstand-Lesen eine Bindung")
	}
}

// TestRunAdministrationExitsOnContextCancel trägt den ctx-Cancel-Austritt
// aus `runAdministration` (Review-Finding F-2): ein bereits abgebrochener
// `ctx` beendet die Schleife nach dem ersten Verarbeitungs-Durchlauf, ohne
// auf ein Wecksignal zu warten.
func TestRunAdministrationExitsOnContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	assembler, err := mapper.NewAssembler("src-admin", map[string]mapper.TableBinding{}, nil)
	if err != nil {
		t.Fatalf("NewAssembler: %v", err)
	}
	deps := administrationDeps{
		requests:      &fakeAdministrationRequestPort{},
		listener:      &fakeAdministrationListener{},
		activation:    &fakeTableActivationPort{},
		enableTables:  &fakeEnableTableUseCase{},
		disableTables: &fakeDisableTableUseCase{},
		schemaStore:   &fakeSchemaStorePort{},
		assembler:     assembler,
		publication:   "cdc_pub",
		pollInterval:  time.Second,
		log:           &recordingLog{},
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		runAdministration(ctx, deps)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("runAdministration ist nach bereits abgebrochenem ctx nicht innerhalb 1s zurückgekehrt")
	}
}

// TestRunAdministrationLogsWarnOnListenerErrorAndKeepsPolling trägt den
// Listener-Fehler-Zweig (Review-Finding F-2): ein gestörtes Wecksignal
// protokolliert eine Warnung, bricht die Schleife aber nicht ab — der
// Fallback-Poll verarbeitet den offenen Antrag trotzdem.
func TestRunAdministrationLogsWarnOnListenerErrorAndKeepsPolling(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	const requestID = model.AdministrationRequestID("req-listener-1")
	tableID := administrationTableID("public", "orders_admin_listener")
	requests := &fakeAdministrationRequestPort{pending: []model.AdministrationRequest{
		{ID: requestID, Source: "src-admin", Schema: "public", Table: "orders_admin_listener", Kind: model.AdministrationRequestEnable},
	}}
	activation := &fakeTableActivationPort{registered: map[string]model.SourceTable{
		"public.orders_admin_listener": {ID: tableID},
	}}
	schemaStore := &fakeSchemaStorePort{versions: map[model.SourceTableID]model.SchemaVersion{
		tableID: {ID: administrationSchemaVersionID(tableID), SourceTableID: tableID, Version: 1},
	}}
	assembler, err := mapper.NewAssembler("src-admin", map[string]mapper.TableBinding{}, nil)
	if err != nil {
		t.Fatalf("NewAssembler: %v", err)
	}
	log := &recordingLog{}
	deps := administrationDeps{
		requests:        requests,
		listener:        &fakeAdministrationListener{err: stderrors.New("connection refused")},
		activation:      activation,
		enableTables:    &fakeEnableTableUseCase{},
		disableTables:   &fakeDisableTableUseCase{},
		schemaStore:     schemaStore,
		columnExclusion: &fakeColumnExclusionPort{},
		transformations: &fakeTransformationPort{},
		assembler:       assembler,
		publication:     "cdc_pub",
		pollInterval:    time.Millisecond,
		log:             log,
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		runAdministration(ctx, deps)
	}()

	// Der Antrag wird bereits im ersten Durchlauf verarbeitet (noch vor dem
	// ersten `WaitForNotification`-Aufruf) — abgewartet wird deshalb
	// zusätzlich mindestens ein Warn-Log des gestörten Listeners, damit der
	// Test den Fehlerpfad tatsächlich durchlaufen sieht, bevor er abbricht.
	deadline := time.After(time.Second)
	for requests.appliedCount() < 1 || func() int { log.mu.Lock(); defer log.mu.Unlock(); return log.warns }() < 1 {
		select {
		case <-deadline:
			t.Fatal("der offene Antrag wurde nicht innerhalb 1s über den Fallback-Poll verarbeitet, oder der gestörte Listener hat kein Warn-Log erzeugt")
		case <-time.After(time.Millisecond):
		}
	}
	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("runAdministration ist nach ctx.Cancel nicht innerhalb 1s zurückgekehrt")
	}

	log.mu.Lock()
	warns := log.warns
	log.mu.Unlock()
	if warns == 0 {
		t.Fatal("kein Warn-Log für den gestörten Listener trotz Fehler auf jedem Wecksignal-Versuch")
	}
}

// TestActivatedTableBindingsCarriesExcludedColumns trägt den Startpfad des
// dauerhaften Ausschlussstandes (`ADR-0065`): der Bindungs-Neuaufbau liest
// den Stand aus seiner dauerhaften Herkunft und übergibt ihn je Tabelle an
// die Bindung — ohne diesen Schritt erfasst ein neu gestarteter Prozess die
// zuvor ausgeschlossene Spalte wieder, obwohl der Antrag `applied` trägt.
// Der Test führt den Neuaufbau bis in die Filterwirkung: das Row Image der
// neu gebauten Bindung lässt den geführten Spaltennamen weg.
func TestActivatedTableBindingsCarriesExcludedColumns(t *testing.T) {
	ctx := context.Background()
	const (
		source  = model.SourceID("src-restart")
		tableID = model.SourceTableID("tbl-restart")
	)
	activation := &fakeTableActivationPort{listed: []model.SourceTable{
		{ID: tableID, SourceID: source, Schema: "public", Table: "orders_restart"},
	}}
	schemaStore := &fakeSchemaStorePort{versions: map[model.SourceTableID]model.SchemaVersion{
		tableID: {ID: "sv-restart", SourceTableID: tableID, Version: 1},
	}}
	exclusions := &fakeColumnExclusionPort{excluded: map[string][]string{
		"public.orders_restart": {"secret"},
	}}

	tables, err := activatedTableBindings(ctx, activation, schemaStore, exclusions, &fakeTransformationPort{}, source)
	if err != nil {
		t.Fatalf("activatedTableBindings = %v, wollen nil", err)
	}
	binding, found := tables["public.orders_restart"]
	if !found {
		t.Fatalf("Bindungs-Stand trägt public.orders_restart nicht: %+v", tables)
	}
	if len(binding.ExcludedColumns) != 1 || binding.ExcludedColumns[0] != "secret" {
		t.Fatalf("ExcludedColumns = %v, wollen [secret]", binding.ExcludedColumns)
	}

	assembler, err := mapper.NewAssembler(source, tables, nil)
	if err != nil {
		t.Fatalf("NewAssembler: %v", err)
	}
	if image := assemblerRowImage(t, assembler, 1, "public", "orders_restart"); image != `{"id":"1"}` {
		t.Fatalf("Row Image der neu gebauten Bindung = %s, wollen ohne den ausgeschlossenen Schlüssel secret", image)
	}
}

// TestProcessAdministrationRequestsDisableEnableCycleRestoresExclusion trägt
// den zweiten Auslöser des Stand-Verlusts (`ADR-0065`): die Deaktivierung
// entfernt den Bindungs-Eintrag samt Ausschlussstand, der Aktivierungs-Zweig
// legt ihn über die dauerhafte Herkunft neu an. Ohne die Ableitung im
// Aktivierungs-Zweig stünde nach `disable` → `enable` eine Bindung ohne
// Ausschluss, und die Spalte würde wieder erfasst — derselbe Verlust wie
// beim Prozess-Neustart, ohne Neustart.
func TestProcessAdministrationRequestsDisableEnableCycleRestoresExclusion(t *testing.T) {
	ctx := context.Background()
	const table = "orders_cycle"
	tableID := administrationTableID("public", table)
	requests := &fakeAdministrationRequestPort{pending: []model.AdministrationRequest{
		{ID: "req-cycle-disable", Source: "src-admin", Schema: "public", Table: table, Kind: model.AdministrationRequestDisable},
	}}
	activation := &fakeTableActivationPort{registered: map[string]model.SourceTable{
		"public." + table: {ID: tableID, SourceID: "src-admin", Schema: "public", Table: table},
	}}
	schemaStore := &fakeSchemaStorePort{versions: map[model.SourceTableID]model.SchemaVersion{
		tableID: {ID: administrationSchemaVersionID(tableID), SourceTableID: tableID, Version: 1},
	}}
	// Der laufende Prozess trägt den Ausschluss bis zur Deaktivierung.
	assembler, err := mapper.NewAssembler("src-admin", map[string]mapper.TableBinding{
		"public." + table: {TableID: tableID, SchemaVersion: administrationSchemaVersionID(tableID), ExcludedColumns: []string{"secret"}},
	}, nil)
	if err != nil {
		t.Fatalf("NewAssembler: %v", err)
	}
	deps := administrationDeps{
		requests:        requests,
		activation:      activation,
		enableTables:    &fakeEnableTableUseCase{},
		disableTables:   &fakeDisableTableUseCase{},
		schemaStore:     schemaStore,
		columnExclusion: &fakeColumnExclusionPort{excluded: map[string][]string{"public." + table: {"secret"}}},
		transformations: &fakeTransformationPort{},
		assembler:       assembler,
		publication:     "cdc_pub",
		log:             &recordingLog{},
	}

	processAdministrationRequests(ctx, deps)
	if assemblerCapturesQualified(t, assembler, 1, "public", table) {
		t.Fatal("Assembler trägt nach Disable weiterhin eine Bindung für " + table)
	}

	requests.mu.Lock()
	requests.pending = []model.AdministrationRequest{{
		ID: "req-cycle-enable", Source: "src-admin", Schema: "public", Table: table, Kind: model.AdministrationRequestEnable,
	}}
	requests.mu.Unlock()
	processAdministrationRequests(ctx, deps)

	if !assemblerCapturesQualified(t, assembler, 2, "public", table) {
		t.Fatal("Assembler trägt nach Enable keine Bindung für " + table + " — AddBinding hat nicht nachgetragen")
	}
	if image := assemblerRowImage(t, assembler, 3, "public", table); image != `{"id":"1"}` {
		t.Fatalf("Row Image nach dem disable/enable-Zyklus = %s, wollen ohne den ausgeschlossenen Schlüssel secret", image)
	}
}

// TestProcessAdministrationRequestsMarksFailedWhenExclusionReadFails trägt
// den Fehlerpfad der dauerhaften Herkunft im Aktivierungs-Zweig: ein
// Lesefehler des Ausschlussstandes endet im `failed`-Vermerk, statt eine
// Bindung ohne den geführten Stand anzulegen — der stille Verlust, den
// `ADR-0065` ausschließt.
func TestProcessAdministrationRequestsMarksFailedWhenExclusionReadFails(t *testing.T) {
	ctx := context.Background()
	const table = "orders_read_fail"
	tableID := administrationTableID("public", table)
	wantErr := stderrors.New("Antrags-Historie nicht lesbar")
	requests := &fakeAdministrationRequestPort{pending: []model.AdministrationRequest{
		{ID: "req-exclusion-read-fail", Source: "src-admin", Schema: "public", Table: table, Kind: model.AdministrationRequestEnable},
	}}
	activation := &fakeTableActivationPort{registered: map[string]model.SourceTable{
		"public." + table: {ID: tableID, SourceID: "src-admin", Schema: "public", Table: table},
	}}
	schemaStore := &fakeSchemaStorePort{versions: map[model.SourceTableID]model.SchemaVersion{
		tableID: {ID: administrationSchemaVersionID(tableID), SourceTableID: tableID, Version: 1},
	}}
	assembler, err := mapper.NewAssembler("src-admin", map[string]mapper.TableBinding{}, nil)
	if err != nil {
		t.Fatalf("NewAssembler: %v", err)
	}
	deps := administrationDeps{
		requests:        requests,
		activation:      activation,
		enableTables:    &fakeEnableTableUseCase{},
		disableTables:   &fakeDisableTableUseCase{},
		schemaStore:     schemaStore,
		columnExclusion: &fakeColumnExclusionPort{err: wantErr},
		transformations: &fakeTransformationPort{},
		assembler:       assembler,
		publication:     "cdc_pub",
		log:             &recordingLog{},
	}

	processAdministrationRequests(ctx, deps)

	if requests.appliedCount() != 0 {
		t.Fatalf("applied = %d, wollen 0 (Ausschlussstand nicht lesbar)", requests.appliedCount())
	}
	if assemblerCapturesQualified(t, assembler, 1, "public", table) {
		t.Fatal("Assembler trägt nach gescheitertem Ausschlussstand-Lesen eine Bindung — sie trüge den geführten Ausschluss nicht")
	}
}
