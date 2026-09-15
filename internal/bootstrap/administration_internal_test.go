package bootstrap

import (
	"context"
	stderrors "errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/decode"
	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/mapper"
	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
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
// reale PostgreSQL-Instanz (Review-Finding F-2, `review-slice-037.md`).

// fakeAdministrationRequestPort trägt einen In-Memory-Stub des
// `outbound.AdministrationRequestPort`: `pending` wird bei `ListPending`
// einmalig geliefert und danach geleert — derselbe Konsum-einmal-Vertrag wie
// die reale Antrags-Queue (ein bereits gelesener Antrag bleibt `pending`,
// bis `MarkApplied`/`MarkFailed` ihn vermerkt; dieser Stub bildet nur den
// Lese-Teil nach, den `processAdministrationRequests` je Durchlauf braucht).
type fakeAdministrationRequestPort struct {
	mu      sync.Mutex
	pending []model.AdministrationRequest
	applied []model.AdministrationRequestID
	failed  map[model.AdministrationRequestID]string
}

func (f *fakeAdministrationRequestPort) ListPending(ctx context.Context) ([]model.AdministrationRequest, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	pending := f.pending
	f.pending = nil
	return pending, nil
}

func (f *fakeAdministrationRequestPort) MarkApplied(ctx context.Context, id model.AdministrationRequestID) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.applied = append(f.applied, id)
	return nil
}

func (f *fakeAdministrationRequestPort) MarkFailed(ctx context.Context, id model.AdministrationRequestID, message string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.failed == nil {
		f.failed = map[model.AdministrationRequestID]string{}
	}
	f.failed[id] = message
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
}

func (f *fakeTableActivationPort) TableExists(ctx context.Context, schema, table string) (bool, error) {
	return true, nil
}

func (f *fakeTableActivationPort) Registered(ctx context.Context, source model.SourceID, schema, table string) (model.SourceTable, bool, error) {
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
}

func (f *fakeSchemaStorePort) CurrentVersion(ctx context.Context, table model.SourceTableID) (model.SchemaVersion, bool, error) {
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

// TestApplyAdministrationRequestRejectsUnknownKind trägt den `default`-Zweig
// (Review-Finding F-2): eine Antragsart außerhalb der geschlossenen Menge
// `enable`/`disable` endet über einen sichtbaren Fehler, statt still
// übersprungen zu werden.
func TestApplyAdministrationRequestRejectsUnknownKind(t *testing.T) {
	ctx := context.Background()
	assembler, err := mapper.NewAssembler("src-admin", map[string]mapper.TableBinding{}, nil)
	if err != nil {
		t.Fatalf("NewAssembler: %v", err)
	}
	deps := administrationDeps{
		activation:    &fakeTableActivationPort{},
		enableTables:  &fakeEnableTableUseCase{},
		disableTables: &fakeDisableTableUseCase{},
		schemaStore:   &fakeSchemaStorePort{},
		assembler:     assembler,
		publication:   "cdc_pub",
		log:           &recordingLog{},
	}
	request := model.AdministrationRequest{
		ID: "req-unknown", Source: "src-admin", Schema: "public", Table: "orders_admin_unknown",
		Kind: model.AdministrationRequestKind("truncate"),
	}

	if err := applyAdministrationRequest(ctx, deps, request); err == nil {
		t.Fatal("applyAdministrationRequest(unbekannte Kind) = nil, wollen einen Fehler (geschlossene Menge enable/disable)")
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

	tables, err := activatedTableBindings(ctx, activation, schemaStore, exclusions, source)
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
