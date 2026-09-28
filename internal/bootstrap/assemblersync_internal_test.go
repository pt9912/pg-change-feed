package bootstrap

import (
	"context"
	stderrors "errors"
	"testing"

	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/mapper"
	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Whitebox-Test (`package bootstrap`): `enableTableWithAssemblerSync`/
// `disableTableWithAssemblerSync` (`assemblersync.go`) sind der Nachtrag-
// Mechanismus des direkten HTTP-/gRPC-Zugriffswegs — dieselben Fakes wie
// `administration_internal_test.go` (`fakeEnableTableUseCase`,
// `fakeTableActivationPort`, `fakeSchemaStorePort`, `fakeColumnExclusionPort`,
// `fakeTransformationPort`, `assemblerCapturesQualified`, `assemblerRowImage`),
// eingebunden auf der Eingabeseite (Erfolg/Fehlschlag des inneren Use Cases),
// nicht nur auf der Ausgabeseite.

// TestEnableTableWithAssemblerSyncAddsBindingOnSuccess trägt den Happy Path:
// ein erfolgreicher innerer `Enable`-Aufruf trägt die Bindung über
// `syncAssemblerAddBinding` in den laufenden Assembler nach — Tabellen-
// Kennung, Schema-Version, Ausschluss- und Regelstand der Quelle
// eingeschlossen. Ohne diesen Nachtrag bliebe eine über den direkten
// HTTP-/gRPC-Zugriffsweg aktivierte Tabelle bis zum nächsten Neustart
// unerfasst.
func TestEnableTableWithAssemblerSyncAddsBindingOnSuccess(t *testing.T) {
	ctx := context.Background()
	const source = model.SourceID("src-api")
	tableID := administrationTableID("public", "orders_api_enable")
	activation := &fakeTableActivationPort{registered: map[string]model.SourceTable{
		"public.orders_api_enable": {ID: tableID, SourceID: source, Schema: "public", Table: "orders_api_enable"},
	}}
	schemaStore := &fakeSchemaStorePort{versions: map[model.SourceTableID]model.SchemaVersion{
		tableID: {ID: administrationSchemaVersionID(tableID), SourceTableID: tableID, Version: 1},
	}}
	columnExclusion := &fakeColumnExclusionPort{excluded: map[string][]string{
		"public.orders_api_enable": {"secret"},
	}}
	transformations := &fakeTransformationPort{rules: map[string][]model.Transformation{}}
	assembler, err := mapper.NewAssembler(source, map[string]mapper.TableBinding{}, nil)
	if err != nil {
		t.Fatalf("NewAssembler: %v", err)
	}
	enableUseCase := &fakeEnableTableUseCase{}
	decorator := enableTableWithAssemblerSync{
		EnableTableUseCase: enableUseCase,
		assembler:          assembler,
		activation:         activation,
		schemaStore:        schemaStore,
		columnExclusion:    columnExclusion,
		transformations:    transformations,
	}

	result, err := decorator.Enable(ctx, inbound.EnableTableCommand{
		Source:      source,
		Schema:      "public",
		Table:       "orders_api_enable",
		TableID:     tableID,
		Publication: "cdc_pub",
	})
	if err != nil {
		t.Fatalf("Enable = %v, wollen nil", err)
	}
	if len(enableUseCase.calls) != 1 {
		t.Fatalf("innere Enable-Aufrufe = %d, wollen 1", len(enableUseCase.calls))
	}
	if got := result.Table.ID; got != tableID {
		t.Fatalf("Ergebnis-TableID = %q, wollen %q", got, tableID)
	}
	if !assemblerCapturesQualified(t, assembler, 1, "public", "orders_api_enable") {
		t.Fatal("Assembler trägt nach Enable keine Bindung — AddBinding hat nicht nachgetragen")
	}
	if image := assemblerRowImage(t, assembler, 2, "public", "orders_api_enable"); image != `{"id":"1"}` {
		t.Fatalf("Row Image = %s, wollen ohne den ausgeschlossenen Schlüssel secret (Ausschlussstand nicht übernommen)", image)
	}
}

// TestEnableTableWithAssemblerSyncSkipsBindingWhenInnerEnableFails trägt den
// Fehlschlag-Pfad auf der Eingabeseite: scheitert der innere Use Case, bleibt
// der Assembler unangetastet — kein `AddBinding` auf Basis eines Aufrufs, der
// nie wirklich aktiviert hat.
func TestEnableTableWithAssemblerSyncSkipsBindingWhenInnerEnableFails(t *testing.T) {
	ctx := context.Background()
	wantErr := stderrors.New("Quelltabelle existiert nicht")
	assembler, err := mapper.NewAssembler("src-api", map[string]mapper.TableBinding{}, nil)
	if err != nil {
		t.Fatalf("NewAssembler: %v", err)
	}
	decorator := enableTableWithAssemblerSync{
		EnableTableUseCase: &fakeEnableTableUseCase{err: wantErr},
		assembler:          assembler,
		activation:         &fakeTableActivationPort{},
		schemaStore:        &fakeSchemaStorePort{},
		columnExclusion:    &fakeColumnExclusionPort{},
		transformations:    &fakeTransformationPort{},
	}

	_, err = decorator.Enable(ctx, inbound.EnableTableCommand{
		Source: "src-api", Schema: "public", Table: "orders_api_fail", Publication: "cdc_pub",
	})
	if !stderrors.Is(err, wantErr) {
		t.Fatalf("Enable-Fehler = %v, wollen %v", err, wantErr)
	}
	if assemblerCapturesQualified(t, assembler, 1, "public", "orders_api_fail") {
		t.Fatal("Assembler trägt eine Bindung, obwohl der innere Enable-Aufruf gescheitert ist")
	}
}

// TestEnableTableWithAssemblerSyncFailsWhenSyncLookupFails trägt den
// Fehlschlag-Pfad des Nachtrags selbst: der innere Use Case gelingt, aber die
// Rücklesung der Bindungs-Zeile (`TableActivationPort.Registered`) findet
// nichts — derselbe Fehlertext wie `applyAdministrationRequest`s
// Enable-Zweig. Der Aufruf meldet den Fehler statt eines Erfolgs mit leerem
// Nachtrag; die reale Aktivierung bleibt in der Datenbank stehen, aber der
// laufende Assembler trägt sie nicht.
func TestEnableTableWithAssemblerSyncFailsWhenSyncLookupFails(t *testing.T) {
	ctx := context.Background()
	assembler, err := mapper.NewAssembler("src-api", map[string]mapper.TableBinding{}, nil)
	if err != nil {
		t.Fatalf("NewAssembler: %v", err)
	}
	decorator := enableTableWithAssemblerSync{
		EnableTableUseCase: &fakeEnableTableUseCase{},
		assembler:          assembler,
		activation:         &fakeTableActivationPort{}, // registered bleibt leer: Registered findet nichts
		schemaStore:        &fakeSchemaStorePort{},
		columnExclusion:    &fakeColumnExclusionPort{},
		transformations:    &fakeTransformationPort{},
	}

	_, err = decorator.Enable(ctx, inbound.EnableTableCommand{
		Source: "src-api", Schema: "public", Table: "orders_api_nobind", Publication: "cdc_pub",
	})
	if err == nil {
		t.Fatal("Enable = nil, wollen einen Fehler (keine Bindungs-Zeile gefunden)")
	}
	if assemblerCapturesQualified(t, assembler, 1, "public", "orders_api_nobind") {
		t.Fatal("Assembler trägt eine Bindung, obwohl der Nachtrag gescheitert ist")
	}
}

// TestDisableTableWithAssemblerSyncRemovesBindingOnSuccess trägt den Happy
// Path der Gegenseite: ein erfolgreicher innerer `Disable`-Aufruf entfernt die
// laufende Bindung — ohne diesen Nachtrag würde eine über den direkten
// HTTP-/gRPC-Zugriffsweg deaktivierte Tabelle bis zum nächsten Neustart
// weiter erfasst.
func TestDisableTableWithAssemblerSyncRemovesBindingOnSuccess(t *testing.T) {
	ctx := context.Background()
	assembler, err := mapper.NewAssembler("src-api", map[string]mapper.TableBinding{
		"public.orders_api_disable": {TableID: "tbl-api-disable", SchemaVersion: "sv-api-disable"},
	}, nil)
	if err != nil {
		t.Fatalf("NewAssembler: %v", err)
	}
	if !assemblerCapturesQualified(t, assembler, 1, "public", "orders_api_disable") {
		t.Fatal("Vorbedingung verletzt: Assembler sollte vor der Deaktivierung gebunden sein")
	}
	disableUseCase := &fakeDisableTableUseCase{}
	decorator := disableTableWithAssemblerSync{DisableTableUseCase: disableUseCase, assembler: assembler}

	if _, err := decorator.Disable(ctx, inbound.DisableTableCommand{
		Source: "src-api", Schema: "public", Table: "orders_api_disable", Publication: "cdc_pub",
	}); err != nil {
		t.Fatalf("Disable = %v, wollen nil", err)
	}
	if len(disableUseCase.calls) != 1 {
		t.Fatalf("innere Disable-Aufrufe = %d, wollen 1", len(disableUseCase.calls))
	}
	if assemblerCapturesQualified(t, assembler, 2, "public", "orders_api_disable") {
		t.Fatal("Assembler trägt nach Disable weiterhin eine Bindung — RemoveBinding hat nicht nachgetragen")
	}
}

// TestDisableTableWithAssemblerSyncKeepsBindingWhenInnerDisableFails trägt
// den Fehlschlag-Pfad auf der Eingabeseite: scheitert der innere Use Case,
// bleibt die laufende Bindung unangetastet — kein `RemoveBinding` auf Basis
// eines Aufrufs, der nie wirklich deaktiviert hat.
func TestDisableTableWithAssemblerSyncKeepsBindingWhenInnerDisableFails(t *testing.T) {
	ctx := context.Background()
	wantErr := stderrors.New("Publication-Entzug fehlgeschlagen")
	assembler, err := mapper.NewAssembler("src-api", map[string]mapper.TableBinding{
		"public.orders_api_disable_fail": {TableID: "tbl-api-disable-fail", SchemaVersion: "sv-api-disable-fail"},
	}, nil)
	if err != nil {
		t.Fatalf("NewAssembler: %v", err)
	}
	decorator := disableTableWithAssemblerSync{
		DisableTableUseCase: &fakeDisableTableUseCase{err: wantErr},
		assembler:           assembler,
	}

	_, err = decorator.Disable(ctx, inbound.DisableTableCommand{
		Source: "src-api", Schema: "public", Table: "orders_api_disable_fail", Publication: "cdc_pub",
	})
	if !stderrors.Is(err, wantErr) {
		t.Fatalf("Disable-Fehler = %v, wollen %v", err, wantErr)
	}
	if !assemblerCapturesQualified(t, assembler, 1, "public", "orders_api_disable_fail") {
		t.Fatal("Assembler trägt keine Bindung mehr, obwohl der innere Disable-Aufruf gescheitert ist")
	}
}
