package bootstrap

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/mapper"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Diese Datei trägt den netzlos prüfbaren Rest der Verdrahtung, der neben
// den Happy Paths in `administration_internal_test.go`,
// `walretention_internal_test.go` und `config_file_internal_test.go`
// offenblieb (`ADR-0082` Cluster D2): die Lese- und Vermerk-Fehlerzweige
// der Administrations-Schleife, die Fehlerzweige des Bindungs-Neuaufbaus,
// die Namens-Teilung und der Mess-Fehlerzweig der WAL-Retention-Prüfung.
// Die reale Ende-zu-Ende-Seite derselben Zusagen liegt in
// `tools/harness/run-integration-tests.sh` (`make test-integration`, u. a.
// das `failed`-Ergebnis einer Spalten-Anfrage gegen eine nicht
// existierende Spalte) und in `make test-replication`.

// TestSplitQualifiedNameTrenntSchemaUndTabelle trägt die Form-Regel der
// Aktivierungs-Zeile (`LH-FA-DAT-002`): der qualifizierte Name trägt Schema
// und Tabellenname durch genau einen Punkt getrennt; beide Teile müssen
// nicht leer sein. Ein `Run`-Start über eine `CDC_TABLES`-Zeile ohne diese
// Form endet in der Startfehlerklasse `configuration`, nicht in einer
// Aktivierung mit leerem Schema oder leerem Tabellennamen.
//
// Rot färbende Mutation: die Bedingung `!found || schema == "" ||
// table == ""` auf `!found` verkürzen — dann liefert `public.` (leerer
// Tabellenname) keinen Fehler mehr und dieser Test bricht.
func TestSplitQualifiedNameTrenntSchemaUndTabelle(t *testing.T) {
	schema, table, err := splitQualifiedName("public.orders")
	if err != nil {
		t.Fatalf("splitQualifiedName(public.orders): unerwarteter Fehler: %v", err)
	}
	if schema != "public" || table != "orders" {
		t.Fatalf("splitQualifiedName(public.orders) = (%q, %q), wollen (public, orders)", schema, table)
	}

	for _, qualified := range []string{"orders", "public.", ".orders"} {
		_, _, err := splitQualifiedName(qualified)
		if !errors.Is(err, ErrConfiguration) {
			t.Fatalf("splitQualifiedName(%q): Fehlerklasse configuration erwartet, erhalten: %v", qualified, err)
		}
	}
}

// TestActivatedTableBindingsFehlerpfade trägt den Bindungs-Neuaufbau am
// Prozessstart (`ADR-0050`, `ADR-0065`): der committed Stand ist die
// Grundlage des Stream-Starts — ein Lesefehler einer der vier Quellen
// (`TableActivationPort.List`, `ColumnExclusionPort.ExcludedColumns`,
// `TransformationPort.TransformationRules`,
// `SchemaStorePort.CurrentVersion`) bricht den Aufbau ab, statt eine
// unvollständige Bindung zu tragen; eine aktivierte Tabelle **ohne**
// registrierte Schema-Version bleibt ohne Bindung (keine Erfassung ohne
// Interpretationsgrundlage).
//
// Rot färbende Mutation (letzter Fall): `if !found { continue }` entfernen
// — dann erscheint die Tabelle mit leerer `SchemaVersion` in der Bindung
// und der Test bricht.
func TestActivatedTableBindingsFehlerpfade(t *testing.T) {
	ctx := context.Background()
	const source = model.SourceID("src-bindungen")
	const tableID = model.SourceTableID("public.bindungs_tabelle")

	listed := []model.SourceTable{{ID: tableID, SourceID: source, Schema: "public", Table: "bindungs_tabelle"}}
	readErr := errors.New("Lesefehler der Quelle")

	cases := []struct {
		name          string
		activation    *fakeTableActivationPort
		schemaStore   *fakeSchemaStorePort
		columns       *fakeColumnExclusionPort
		rules         *fakeTransformationPort
		wantErr       error
		wantBoundKeys []string
	}{
		{
			name:        "List-Fehler",
			activation:  &fakeTableActivationPort{listErr: readErr},
			schemaStore: &fakeSchemaStorePort{},
			columns:     &fakeColumnExclusionPort{},
			rules:       &fakeTransformationPort{},
			wantErr:     readErr,
		},
		{
			name:        "ExcludedColumns-Fehler",
			activation:  &fakeTableActivationPort{listed: listed},
			schemaStore: &fakeSchemaStorePort{},
			columns:     &fakeColumnExclusionPort{err: readErr},
			rules:       &fakeTransformationPort{},
			wantErr:     readErr,
		},
		{
			name:        "TransformationRules-Fehler",
			activation:  &fakeTableActivationPort{listed: listed},
			schemaStore: &fakeSchemaStorePort{},
			columns:     &fakeColumnExclusionPort{},
			rules:       &fakeTransformationPort{err: readErr},
			wantErr:     readErr,
		},
		{
			name:        "CurrentVersion-Fehler",
			activation:  &fakeTableActivationPort{listed: listed},
			schemaStore: &fakeSchemaStorePort{err: readErr},
			columns:     &fakeColumnExclusionPort{},
			rules:       &fakeTransformationPort{},
			wantErr:     readErr,
		},
		{
			name:        "Tabelle ohne registrierte Schema-Version bleibt ungebunden",
			activation:  &fakeTableActivationPort{listed: listed},
			schemaStore: &fakeSchemaStorePort{versions: map[model.SourceTableID]model.SchemaVersion{}},
			columns:     &fakeColumnExclusionPort{},
			rules:       &fakeTransformationPort{},
			wantErr:     nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tables, err := activatedTableBindings(ctx, c.activation, c.schemaStore, c.columns, c.rules, source)
			if !errors.Is(err, c.wantErr) {
				t.Fatalf("activatedTableBindings: Fehler %v, wollen %v", err, c.wantErr)
			}
			if err != nil {
				if tables != nil {
					t.Fatalf("activatedTableBindings bei Fehler: Bindung %v statt nil — eine unvollständige Bindung darf den Stream-Start nicht tragen", tables)
				}
				return
			}
			if len(tables) != len(c.wantBoundKeys) {
				t.Fatalf("activatedTableBindings = %v, wollen %d Bindungen", tables, len(c.wantBoundKeys))
			}
		})
	}
}

// TestProcessAdministrationRequestsLesefehlerBeendetDurchlauf trägt die
// best-effort-Haltung des Antrags-Durchlaufs (`ADR-0050`): ein
// Lesefehler der Antrags-Queue beendet diesen Durchlauf sichtbar
// (Log-Zeile) und wird nicht als Absturz behandelt — der nächste Durchlauf
// versucht erneut. Nichts wird in diesem Durchlauf als verarbeitet
// vermerkt.
//
// Rot färbende Mutation: den Fehlerzweig auf `if err != nil { }` (ohne
// Log und ohne `return`) verkürzen — dann verstummt die Zeile und der Test
// bricht.
func TestProcessAdministrationRequestsLesefehlerBeendetDurchlauf(t *testing.T) {
	ctx := context.Background()
	requests := &fakeAdministrationRequestPort{listErr: errors.New("Queue nicht erreichbar")}
	enableTables := &fakeEnableTableUseCase{}
	deps := neueAdministrationsDeps(t, requests, &fakeTableActivationPort{}, &fakeSchemaStorePort{}, &fakeColumnExclusionPort{}, enableTables, &fakeDisableTableUseCase{})

	processAdministrationRequests(ctx, deps)

	if !deps.log.(*recordingLog).contains("WARN", "Anträge lesen fehlgeschlagen") {
		t.Fatal("Lesefehler der Antrags-Queue wird nicht protokolliert — der Durchlauf endet stumm")
	}
	if len(enableTables.calls) != 0 {
		t.Fatalf("Enable-Aufrufe = %d, wollen 0 (kein Antrag gelesen)", len(enableTables.calls))
	}
	if requests.appliedCount() != 0 {
		t.Fatalf("applied = %d, wollen 0", requests.appliedCount())
	}
}

// TestProcessAdministrationRequestsVermerkfehlerBleibenBestEffort trägt
// die beiden Vermerk-Zweige (`ADR-0050`): der Durchlauf ist best-effort —
// ein gescheiterter `MarkFailed` bzw. `MarkApplied` wird protokolliert,
// ohne die bereits geleistete Arbeit rückgängig zu machen. Der erste Fall
// belegt zugleich, dass der gescheiterte Antrag als `failed` **versucht**
// wird (nicht `pending` bleibt und jeden Durchlauf erneut läuft).
//
// Rot färbende Mutation: den jeweiligen inneren Fehlerzweig
// (`if markErr != nil { log.Warn(...) }` bzw. `if err != nil { log.Warn(...)
// }`) entfernen — dann fehlt die Zeile und der Test bricht.
func TestProcessAdministrationRequestsVermerkfehlerBleibenBestEffort(t *testing.T) {
	ctx := context.Background()

	t.Run("MarkFailed scheitert", func(t *testing.T) {
		const requestID = model.AdministrationRequestID("req-vermerk-failed")
		requests := &fakeAdministrationRequestPort{
			pending: []model.AdministrationRequest{{
				ID: requestID, Source: "src-admin", Schema: "public", Table: "orders_vermerk_failed",
				Kind: model.AdministrationRequestEnable,
			}},
			markFailedErr: errors.New("Vermerk nicht schreibbar"),
		}
		enableTables := &fakeEnableTableUseCase{err: errors.New("Aktivierung gescheitert")}
		deps := neueAdministrationsDeps(t, requests, &fakeTableActivationPort{}, &fakeSchemaStorePort{}, &fakeColumnExclusionPort{}, enableTables, &fakeDisableTableUseCase{})

		processAdministrationRequests(ctx, deps)

		if !deps.log.(*recordingLog).contains("WARN", "Antrag fehlgeschlagen") {
			t.Fatal("gescheiterter Antrag wird nicht protokolliert")
		}
		if !deps.log.(*recordingLog).contains("WARN", "Fehlschlag nicht vermerkt") {
			t.Fatal("gescheiterter MarkFailed wird nicht protokolliert — der Antrag bliebe pending und liefe jeden Durchlauf erneut")
		}
		if requests.appliedCount() != 0 {
			t.Fatalf("applied = %d, wollen 0 (der Antrag ist gescheitert)", requests.appliedCount())
		}
	})

	t.Run("MarkApplied scheitert", func(t *testing.T) {
		const requestID = model.AdministrationRequestID("req-vermerk-applied")
		tableID := administrationTableID("public", "orders_vermerk_applied")
		requests := &fakeAdministrationRequestPort{
			pending: []model.AdministrationRequest{{
				ID: requestID, Source: "src-admin", Schema: "public", Table: "orders_vermerk_applied",
				Kind: model.AdministrationRequestEnable,
			}},
			markAppliedErr: errors.New("Vermerk nicht schreibbar"),
		}
		activation := &fakeTableActivationPort{registered: map[string]model.SourceTable{
			"public.orders_vermerk_applied": {ID: tableID, SourceID: "src-admin", Schema: "public", Table: "orders_vermerk_applied"},
		}}
		schemaStore := &fakeSchemaStorePort{versions: map[model.SourceTableID]model.SchemaVersion{
			tableID: {ID: administrationSchemaVersionID(tableID), SourceTableID: tableID, Version: 1},
		}}
		enableTables := &fakeEnableTableUseCase{}
		deps := neueAdministrationsDeps(t, requests, activation, schemaStore, &fakeColumnExclusionPort{}, enableTables, &fakeDisableTableUseCase{})

		processAdministrationRequests(ctx, deps)

		if len(enableTables.calls) != 1 {
			t.Fatalf("Enable-Aufrufe = %d, wollen 1 (die Arbeit bleibt geleistet, nur der Vermerk scheitert)", len(enableTables.calls))
		}
		if !deps.log.(*recordingLog).contains("WARN", "Erfolg nicht vermerkt") {
			t.Fatal("gescheiterter MarkApplied wird nicht protokolliert")
		}
	})
}

// TestApplyAdministrationRequestFehlerpfade trägt die Fehlerzweige des
// einzelnen Antrags: jeder Zweig bricht ab, **bevor** er die laufende
// `Assembler`-Bindung nachträgt — eine Bindung ohne die zugehörige
// registrierte Zeile wäre eine Erfassung ohne Interpretationsgrundlage
// (`wiring.go` an `applyAdministrationRequest`).
//
// Rot färbende Mutation (Fall „Registered findet keine Zeile"): die
// Prüfung `if !found { return ... }` entfernen — dann trägt der Aufruf die
// Bindung mit leerer `TableID` nach und der Test bricht.
func TestApplyAdministrationRequestFehlerpfade(t *testing.T) {
	ctx := context.Background()
	readErr := errors.New("Lese-Fehler")
	writeErr := errors.New("Schreib-Fehler")

	const (
		schema = "public"
		table  = "orders_admin_fehler"
	)
	qualified := schema + "." + table
	tableID := administrationTableID(schema, table)

	enableRequest := model.AdministrationRequest{ID: "req-fehler", Source: "src-admin", Schema: schema, Table: table, Kind: model.AdministrationRequestEnable}
	disableRequest := model.AdministrationRequest{ID: "req-fehler", Source: "src-admin", Schema: schema, Table: table, Kind: model.AdministrationRequestDisable}
	excludeRequest := model.AdministrationRequest{ID: "req-fehler", Source: "src-admin", Schema: schema, Table: table, Column: "secret", Kind: model.AdministrationRequestExcludeColumn}

	registered := map[string]model.SourceTable{
		qualified: {ID: tableID, SourceID: "src-admin", Schema: schema, Table: table},
	}
	versions := map[model.SourceTableID]model.SchemaVersion{
		tableID: {ID: administrationSchemaVersionID(tableID), SourceTableID: tableID, Version: 1},
	}

	cases := []struct {
		name        string
		request     model.AdministrationRequest
		activation  *fakeTableActivationPort
		schemaStore *fakeSchemaStorePort
		disable     *fakeDisableTableUseCase
		exclude     *fakeExcludeColumnUseCase
		wantFehler  string
	}{
		{
			name:        "Enable: Registered-Lesefehler",
			request:     enableRequest,
			activation:  &fakeTableActivationPort{registeredErr: readErr},
			schemaStore: &fakeSchemaStorePort{versions: versions},
			disable:     &fakeDisableTableUseCase{},
			exclude:     &fakeExcludeColumnUseCase{},
			wantFehler:  readErr.Error(),
		},
		{
			name:        "Enable: keine registrierte Bindungs-Zeile",
			request:     enableRequest,
			activation:  &fakeTableActivationPort{registered: map[string]model.SourceTable{}},
			schemaStore: &fakeSchemaStorePort{versions: versions},
			disable:     &fakeDisableTableUseCase{},
			exclude:     &fakeExcludeColumnUseCase{},
			wantFehler:  "Aktivierung ohne Bindungs-Zeile",
		},
		{
			name:        "Enable: CurrentVersion-Lesefehler",
			request:     enableRequest,
			activation:  &fakeTableActivationPort{registered: registered},
			schemaStore: &fakeSchemaStorePort{err: readErr},
			disable:     &fakeDisableTableUseCase{},
			exclude:     &fakeExcludeColumnUseCase{},
			wantFehler:  readErr.Error(),
		},
		{
			name:        "Enable: keine registrierte Schema-Version",
			request:     enableRequest,
			activation:  &fakeTableActivationPort{registered: registered},
			schemaStore: &fakeSchemaStorePort{versions: map[model.SourceTableID]model.SchemaVersion{}},
			disable:     &fakeDisableTableUseCase{},
			exclude:     &fakeExcludeColumnUseCase{},
			wantFehler:  "Aktivierung ohne registrierte Schema-Version",
		},
		{
			name:        "Disable: Use-Case-Fehler",
			request:     disableRequest,
			activation:  &fakeTableActivationPort{registered: registered},
			schemaStore: &fakeSchemaStorePort{versions: versions},
			disable:     &fakeDisableTableUseCase{err: writeErr},
			exclude:     &fakeExcludeColumnUseCase{},
			wantFehler:  writeErr.Error(),
		},
		{
			name:        "ExcludeColumn: Use-Case-Fehler",
			request:     excludeRequest,
			activation:  &fakeTableActivationPort{registered: registered},
			schemaStore: &fakeSchemaStorePort{versions: versions},
			disable:     &fakeDisableTableUseCase{},
			exclude:     &fakeExcludeColumnUseCase{err: writeErr},
			wantFehler:  writeErr.Error(),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assembler, err := mapper.NewAssembler("src-admin", map[string]mapper.TableBinding{}, nil)
			if err != nil {
				t.Fatalf("NewAssembler: %v", err)
			}
			deps := administrationDeps{
				requests:        &fakeAdministrationRequestPort{},
				activation:      c.activation,
				enableTables:    &fakeEnableTableUseCase{},
				disableTables:   c.disable,
				excludeColumns:  c.exclude,
				includeColumns:  &fakeIncludeColumnUseCase{},
				schemaStore:     c.schemaStore,
				columnExclusion: &fakeColumnExclusionPort{},
				transformations: &fakeTransformationPort{},
				assembler:       assembler,
				publication:     "cdc_pub",
				log:             &recordingLog{},
			}

			err = applyAdministrationRequest(ctx, deps, c.request)
			if err == nil {
				t.Fatalf("applyAdministrationRequest: Fehler erwartet (%s), erhalten nil", c.wantFehler)
			}
			if !strings.Contains(err.Error(), c.wantFehler) {
				t.Fatalf("applyAdministrationRequest: Fehler %q trägt %q nicht", err.Error(), c.wantFehler)
			}
			if assemblerCapturesQualified(t, assembler, 1, schema, table) {
				t.Fatal("gescheiterter Antrag hat die Assembler-Bindung nachgetragen — die Bindung hätte keine registrierte Zeile")
			}
		})
	}
}

// TestRunWALRetentionCheckProtokolliertMessfehlerUndLaeuftWeiter trägt die
// best-effort-Haltung der WAL-Retention-Prüfung (`SPEC-008`/`ADR-0049`):
// eine fehlgeschlagene Messung selbst beendet den Erfassungspfad nicht —
// sie wird protokolliert, der Durchlauf läuft weiter, `stopStream` und
// `fault` bleiben unberührt. Die Prüfung dieses Zweigs ist zugleich die
// Ursache dafür, dass dieser Block `count > 0` trägt: er ist nicht mehr
// vom Zufall des Kontext-Endes abhängig (siehe
// `harness/sensors/coverage-gate.md` §Zählbasis).
//
// Rot färbende Mutation: den `log.Warn`-Aufruf im Fehlerzweig entfernen —
// dann verstummt der Fehlschlag und der Test bricht.
func TestRunWALRetentionCheckProtokolliertMessfehlerUndLaeuftWeiter(t *testing.T) {
	measurer := &fakeFehlerMessung{err: errors.New("Verbindung gestört")}
	log := &recordingLog{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	streamCtx, stopStream := context.WithCancel(context.Background())
	defer stopStream()
	var fault walRetentionFault

	done := make(chan struct{})
	go func() {
		defer close(done)
		runWALRetentionCheck(ctx, measurer, log, time.Millisecond, 1024, 8192, stopStream, &fault)
	}()

	deadline := time.After(2 * time.Second)
	for !log.contains("WARN", "WAL-Rückstand-Messung fehlgeschlagen") {
		select {
		case <-deadline:
			t.Fatal("Messfehler wurde nicht innerhalb 2s protokolliert")
		case <-time.After(time.Millisecond):
		}
	}

	if streamCtx.Err() != nil {
		t.Fatal("stopStream wurde wegen eines Messfehlers aufgerufen — die Messung selbst bricht den Erfassungspfad nicht ab")
	}
	if err := fault.get(); err != nil {
		t.Fatalf("fault = %v, wollen nil (ein Messfehler setzt den Schwellen-Fehler nicht)", err)
	}

	cancel()
	<-done
}

// fakeFehlerMessung liefert bei jedem `Measure` denselben Fehler — die
// Gegenprobe zu `fakeWALRetentionMeasurer`, der Byte-Werte liefert. Der
// Fehler kommt sofort zurück (kein Blockieren): die Prüf-Goroutine soll
// den Fehlerzweig bei **jedem** Tick fahren, nicht am Kontext-Ende hängen.
type fakeFehlerMessung struct {
	err   error
	calls int
}

func (f *fakeFehlerMessung) Measure(ctx context.Context) (int64, error) {
	f.calls++
	return 0, f.err
}

// neueAdministrationsDeps baut die Verdrahtung der Antrags-Schleife mit
// den übergebenen Fakes und einem frischen Assembler — dieselbe Form wie
// in `administration_internal_test.go`, hier als Helfer, weil die
// Fehlerzweige denselben Aufbau in sechs Varianten brauchen.
func neueAdministrationsDeps(
	t *testing.T,
	requests *fakeAdministrationRequestPort,
	activation *fakeTableActivationPort,
	schemaStore *fakeSchemaStorePort,
	columnExclusion *fakeColumnExclusionPort,
	enableTables *fakeEnableTableUseCase,
	disableTables *fakeDisableTableUseCase,
) administrationDeps {
	t.Helper()
	assembler, err := mapper.NewAssembler("src-admin", map[string]mapper.TableBinding{}, nil)
	if err != nil {
		t.Fatalf("NewAssembler: %v", err)
	}
	return administrationDeps{
		requests:        requests,
		activation:      activation,
		enableTables:    enableTables,
		disableTables:   disableTables,
		excludeColumns:  &fakeExcludeColumnUseCase{},
		includeColumns:  &fakeIncludeColumnUseCase{},
		schemaStore:     schemaStore,
		columnExclusion: columnExclusion,
		transformations: &fakeTransformationPort{},
		assembler:       assembler,
		publication:     "cdc_pub",
		log:             &recordingLog{},
	}
}
