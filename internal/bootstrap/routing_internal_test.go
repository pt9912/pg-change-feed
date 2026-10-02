package bootstrap

import (
	"context"
	stderrors "errors"
	"strconv"
	"testing"
	"time"

	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/decode"
	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/mapper"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/excludecolumn"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/removeroute"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/setroute"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Whitebox-Tests des Routing-Antragswegs (`LH-FA-CFG-008`) in der
// Verdrahtung: `applyAdministrationRequest` prüft über die realen Use Cases,
// trägt die Regel in die laufende `Assembler`-Bindung nach, und der
// Regelstand entsteht beim Anlegen einer Bindung aus seiner dauerhaften
// Herkunft. Die Herkunft ist hier ein Stub, der den Stand wie der Store aus
// den vermerkten Anträgen faltet; der reale Store-Beleg (Ordnung bei gleichem
// Zeitstempel, Rollen) liegt in `internal/adapters/driven/postgresstorage`
// (`make test-store`).

// fakeRoutingPort trägt den `outbound.RoutingPort` als In-Memory-Stub. Ohne
// `derivedFrom` liefert er den festen Stand `rules`; mit `derivedFrom` leitet er
// den Regelstand wie der Store aus den **vermerkten** (`applied`) Anträgen des
// Request-Fakes ab, in der Ordnung von `history`, gefaltet je Tabelle — ein
// `pending` oder `failed` vermerkter Antrag trägt keinen Stand.
type fakeRoutingPort struct {
	rules       map[string][]model.RouteRule
	derivedFrom *fakeAdministrationRequestPort
	history     []model.AdministrationRequest
	err         error
}

func (f *fakeRoutingPort) RoutingRules(ctx context.Context, source model.SourceID) (map[string][]model.RouteRule, error) {
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
	records := map[string][]model.RouteRecord{}
	for _, request := range f.history {
		if !applied[request.ID] || request.Source != source {
			continue
		}
		if request.Kind != model.AdministrationRequestSetRoute && request.Kind != model.AdministrationRequestRemoveRoute {
			continue
		}
		qualified := request.Schema + "." + request.Table
		records[qualified] = append(records[qualified], model.RouteRecord{Kind: request.Kind, Name: request.RuleName, Spec: request.RuleSpec})
	}
	state := map[string][]model.RouteRule{}
	for qualified, list := range records {
		folded, err := model.FoldRoutes(list)
		if err != nil {
			return nil, err
		}
		if len(folded) > 0 {
			state[qualified] = folded
		}
	}
	return state, nil
}

var _ outbound.RoutingPort = (*fakeRoutingPort)(nil)

// routeColumnsPort trägt den `ColumnExclusionPort` mit Spaltenexistenz je
// Adresse und festem Ausschlussstand.
type routeColumnsPort struct {
	existing map[string]bool
	excluded map[string][]string
}

func (f *routeColumnsPort) ColumnExists(ctx context.Context, schema, table, column string) (bool, error) {
	return f.existing[schema+"."+table+"."+column], nil
}

func (f *routeColumnsPort) ExcludedColumns(ctx context.Context, source model.SourceID) (map[string][]string, error) {
	return f.excluded, nil
}

var _ outbound.ColumnExclusionPort = (*routeColumnsPort)(nil)

// routeTable ist die Tabelle der Routing-Tests; ihre Relation im Assembler
// trägt `id` und `region`, der Katalog zusätzlich `secret` (ausgeschlossen)
// und `note`.
const routeTable = "orders_routes"

func routeRequest(id, name, spec string) model.AdministrationRequest {
	return model.AdministrationRequest{
		ID: model.AdministrationRequestID(id), Source: "src-admin", Schema: "public", Table: routeTable,
		RuleName: name, RuleSpec: spec, Kind: model.AdministrationRequestSetRoute,
	}
}

func unrouteRequest(id, name string) model.AdministrationRequest {
	return model.AdministrationRequest{
		ID: model.AdministrationRequestID(id), Source: "src-admin", Schema: "public", Table: routeTable,
		RuleName: name, Kind: model.AdministrationRequestRemoveRoute,
	}
}

func excludeRequest(id, column string) model.AdministrationRequest {
	return model.AdministrationRequest{
		ID: model.AdministrationRequestID(id), Source: "src-admin", Schema: "public", Table: routeTable,
		Column: column, Kind: model.AdministrationRequestExcludeColumn,
	}
}

func routeSpecText(target string, order int, column, equals string) string {
	spec := `{"target": "` + target + `", "order": ` + strconv.Itoa(order)
	if column != "" {
		spec += `, "when": {"column": "` + column + `", "equals": "` + equals + `"}`
	}
	return spec + "}"
}

// routeFixture baut die Verdrahtung der Routing-Tests: die Antrags-Queue mit den
// übergebenen Anträgen, den Routing-Regelstand, den der Store aus den
// vermerkten Anträgen ableitet, die vier realen Use Cases (Routing und
// Spaltenausschluss) und einen Assembler, dessen Tabelle
// `public.orders_routes` gebunden ist.
func routeFixture(t *testing.T, requests ...model.AdministrationRequest) (administrationDeps, *fakeAdministrationRequestPort, *mapper.Assembler) {
	t.Helper()
	queue := &fakeAdministrationRequestPort{pending: append([]model.AdministrationRequest(nil), requests...)}
	routing := &fakeRoutingPort{derivedFrom: queue, history: requests}
	columns := &routeColumnsPort{
		existing: map[string]bool{
			"public." + routeTable + ".id": true, "public." + routeTable + ".region": true,
			"public." + routeTable + ".secret": true, "public." + routeTable + ".note": true,
		},
		excluded: map[string][]string{"public." + routeTable: {"secret"}},
	}
	assembler, err := mapper.NewAssembler("src-admin", map[string]mapper.TableBinding{
		"public." + routeTable: {TableID: "tbl-routes", SchemaVersion: "sv-routes"},
	}, nil)
	if err != nil {
		t.Fatalf("NewAssembler: %v", err)
	}
	deps := administrationDeps{
		requests:        queue,
		activation:      &fakeTableActivationPort{},
		enableTables:    &fakeEnableTableUseCase{},
		disableTables:   &fakeDisableTableUseCase{},
		schemaStore:     &fakeSchemaStorePort{},
		columnExclusion: columns,
		transformations: &fakeTransformationPort{},
		routing:         routing,
		excludeColumns:  excludecolumn.NewExcludeColumnService(columns, routing),
		setRoutes:       setroute.NewSetRouteService(routing, columns),
		removeRoutes:    removeroute.NewRemoveRouteService(routing),
		assembler:       assembler,
		source:          "src-admin",
		publication:     "cdc_pub",
		log:             &recordingLog{},
	}
	return deps, queue, assembler
}

// routedTarget erfasst eine Änderung der Tabelle `routeTable` mit dem Wert
// `region` über die öffentliche `Consume`-Schnittstelle und liefert das
// Zustellziel der erzeugten Change (leer: kein Ziel).
func routedTarget(t *testing.T, assembler *mapper.Assembler, xid uint32, region string) model.RouteTarget {
	t.Helper()
	return routedTargetIn(t, assembler, xid, routeTable, "region", region)
}

// routedTargetIn erfasst eine Änderung der Tabelle `table` (Relation mit dem
// Schlüssel `id` und der Spalte `column`) mit dem Wert `value` und liefert das
// Zustellziel der erzeugten Change.
func routedTargetIn(t *testing.T, assembler *mapper.Assembler, xid uint32, table, column, value string) model.RouteTarget {
	t.Helper()
	ctx := context.Background()
	region := value
	rel := &decode.Relation{Schema: "public", Name: table, Columns: []decode.Column{
		{Name: "id", Key: true}, {Name: column},
	}}
	if _, err := assembler.Consume(ctx, decode.Begin{XID: xid}); err != nil {
		t.Fatalf("Begin: %v", err)
	}
	id := "1"
	if _, err := assembler.Consume(ctx, decode.Change{Relation: rel, Operation: decode.OpInsert, New: []*string{&id, &region}}); err != nil {
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
	return changes[0].RouteTarget
}

func process(ctx context.Context, deps administrationDeps, queue *fakeAdministrationRequestPort, request model.AdministrationRequest) {
	queue.mu.Lock()
	queue.pending = []model.AdministrationRequest{request}
	queue.mu.Unlock()
	processAdministrationRequests(ctx, deps)
}

// TestProcessAdministrationRequestsSetAndRemoveRouteTakeEffectLive trägt den
// Happy Path über die Verdrahtung (`LH-FA-CFG-008`): ein `set_route`-Antrag
// wird `applied`, und die laufende `Assembler`-Bindung trägt die Regel ohne
// Neustart — die nächste Change führt das Ziel; ein `remove_route`-Antrag
// nimmt sie wieder heraus. Rot färbende Mutation: den Nachtrag
// `deps.assembler.SetRoute` bzw. `RemoveRoute` streichen — der Antrag bleibt
// `applied`, das Ziel bleibt leer bzw. die Regel bleibt.
func TestProcessAdministrationRequestsSetAndRemoveRouteTakeEffectLive(t *testing.T) {
	ctx := context.Background()
	set := routeRequest("req-set", "eu_orders", routeSpecText("eu", 10, "region", "eu"))
	remove := unrouteRequest("req-remove", "eu_orders")
	deps, queue, assembler := routeFixture(t, set, remove)
	if target := routedTarget(t, assembler, 1, "eu"); target != "" {
		t.Fatalf("Ziel vor dem Antrag = %q, wollen keines", target)
	}

	process(ctx, deps, queue, set)
	if !isApplied(queue, "req-set") {
		message, _ := failureOf(queue, "req-set")
		t.Fatalf("set_route nicht applied, Fehlertext %q", message)
	}
	if target := routedTarget(t, assembler, 2, "eu"); target != "eu" {
		t.Fatalf("Ziel nach set_route = %q, wollen eu ohne Neustart", target)
	}
	if target := routedTarget(t, assembler, 3, "us"); target != "" {
		t.Fatalf("Ziel einer nicht treffenden Change = %q, wollen keines", target)
	}

	process(ctx, deps, queue, remove)
	if !isApplied(queue, "req-remove") {
		message, _ := failureOf(queue, "req-remove")
		t.Fatalf("remove_route nicht applied, Fehlertext %q", message)
	}
	if target := routedTarget(t, assembler, 4, "eu"); target != "" {
		t.Fatalf("Ziel nach remove_route = %q, wollen keines", target)
	}
}

// TestProcessAdministrationRequestsRouteViolationsFailWithSpecTexts trägt die
// Konfliktfreiheit durch die Verdrahtung: jeder Verstoß endet als `failed` mit
// dem Fehlertext der Spec, der Antrag wird nicht `applied`, und die laufende
// Bindung trägt weiter nur die beiden ersten Regeln (die Ziele bleiben, wie sie
// nach ihnen waren). Die Antragsfolge trägt `applied` zwischen den Anträgen: R1
// prüft gegen vermerkte Anträge.
func TestProcessAdministrationRequestsRouteViolationsFailWithSpecTexts(t *testing.T) {
	ctx := context.Background()
	first := routeRequest("req-1", "eu_orders", routeSpecText("eu", 10, "region", "eu"))
	closing := routeRequest("req-2", "rest", routeSpecText("sonstige", 100, "", ""))
	table := "public." + routeTable
	violations := []struct {
		request model.AdministrationRequest
		text    string
	}{
		{routeRequest("req-r1", "eu_orders", routeSpecText("x", 5, "region", "q")), "abgelehnt [PCF-E8020]: Regelname bereits vergeben: " + table + ".eu_orders"},
		{routeRequest("req-r2", "zweite", routeSpecText("x", 10, "region", "q")), "abgelehnt [PCF-E8030]: order bereits vergeben: " + table + ".10"},
		{routeRequest("req-r3a", "dritte", routeSpecText("x", 5, "gibt_es_nicht", "q")), "abgelehnt [PCF-E8024]: Spalte existiert nicht an der Quelle: " + table + ".gibt_es_nicht"},
		{routeRequest("req-r3b", "vierte", routeSpecText("x", 5, "secret", "q")), "abgelehnt [PCF-E8033]: Spalte ist ausgeschlossen: " + table + ".secret"},
		{routeRequest("req-r4a", "fuenfte", routeSpecText("x", 200, "", "")), "abgelehnt [PCF-E8032]: Regel ohne when bereits vorhanden: " + table + ".rest"},
		{routeRequest("req-r4c", "sechste", routeSpecText("x", 150, "region", "q")), "abgelehnt [PCF-E8032]: order liegt hinter der Regel ohne when: " + table + ".rest"},
		{routeRequest("req-r5", "siebte", routeSpecText("x", 15, "region", "eu")), "abgelehnt [PCF-E8031]: Bedingung bereits vergeben: " + table + ".region"},
		{unrouteRequest("req-r6", "gibt_es_nicht"), "abgelehnt [PCF-E8023]: Regelname nicht geführt: " + table + ".gibt_es_nicht"},
		{routeRequest("req-key", "achte", `{"target": "x", "order": 5, "extra": 1}`), "abgelehnt [PCF-E8011]: unbekannter Schlüssel in rule_spec: extra"},
		{routeRequest("req-form", "neunte", `{"target": "x"}`), "abgelehnt [PCF-E8011]: rule_spec ist ungültig: " + table + ".neunte"},
		{routeRequest("req-target", "zehnte", routeSpecText("EU", 5, "", "")), "abgelehnt [PCF-E8012]: Zielname ist ungültig: " + table + ".EU"},
		{routeRequest("req-name", "Elfte", routeSpecText("x", 5, "", "")), "abgelehnt [PCF-E8010]: Regelname ist ungültig: " + table + ".Elfte"},
		// Zeilen mit fehlendem Regelnamen und fehlender Regelform (SQL-NULL):
		// sie erreichen den Zweig als Antrag und enden `failed`, ohne die Queue
		// anzuhalten.
		{routeRequest("req-empty-name", "", routeSpecText("x", 5, "", "")), "abgelehnt [PCF-E8010]: Regelname ist ungültig: " + table + "."},
		{routeRequest("req-empty-spec", "zwoelfte", ""), "abgelehnt [PCF-E8011]: rule_spec ist ungültig: " + table + ".zwoelfte"},
	}
	all := []model.AdministrationRequest{first, closing}
	for _, violation := range violations {
		all = append(all, violation.request)
	}
	deps, queue, assembler := routeFixture(t, all...)
	process(ctx, deps, queue, first)
	process(ctx, deps, queue, closing)
	if !isApplied(queue, "req-1") || !isApplied(queue, "req-2") {
		t.Fatal("die beiden ersten Regeln sind nicht applied")
	}
	xid := uint32(10)
	for _, violation := range violations {
		process(ctx, deps, queue, violation.request)
		message, found := failureOf(queue, string(violation.request.ID))
		if !found || message != violation.text {
			t.Fatalf("%s: failed = %t, Fehlertext %q, wollen %q", violation.request.ID, found, message, violation.text)
		}
		if isApplied(queue, string(violation.request.ID)) {
			t.Fatalf("%s: der Antrag ist applied", violation.request.ID)
		}
		if target := routedTarget(t, assembler, xid, "eu"); target != "eu" {
			t.Fatalf("%s: Ziel von region=eu = %q, wollen eu (Regelstand unverändert)", violation.request.ID, target)
		}
		if target := routedTarget(t, assembler, xid+1, "zz"); target != "sonstige" {
			t.Fatalf("%s: Ziel von region=zz = %q, wollen sonstige (Regelstand unverändert)", violation.request.ID, target)
		}
		xid += 2
	}
}

// TestProcessAdministrationRequestsExcludeColumnAndRouteBlockEachOther trägt R3
// in beide Richtungen durch die Verdrahtung (`LH-QA-SEC-004`): `exclude_column`
// gegen die Spalte der Routing-Bedingung endet `failed`, und die Spalte bleibt
// in der Erfassung; `exclude_column` gegen eine Spalte ohne Bedingung bleibt
// `applied` (Regressionsfall des Bestandspfads). Die Gegenrichtung (Regel gegen
// ausgeschlossene Spalte) trägt der Fall `req-r3b` oben. Rot färbende Mutation:
// den `routing`-Port aus `excludecolumn.NewExcludeColumnService` durch einen
// leeren Stand ersetzen — `req-exclude-region` wird `applied`.
func TestProcessAdministrationRequestsExcludeColumnAndRouteBlockEachOther(t *testing.T) {
	ctx := context.Background()
	rule := routeRequest("req-rule", "eu_orders", routeSpecText("eu", 10, "region", "eu"))
	blocked := excludeRequest("req-exclude-region", "region")
	free := excludeRequest("req-exclude-note", "note")
	deps, queue, assembler := routeFixture(t, rule, blocked, free)
	process(ctx, deps, queue, rule)
	if !isApplied(queue, "req-rule") {
		t.Fatal("die Regel ist nicht applied")
	}

	process(ctx, deps, queue, blocked)
	message, found := failureOf(queue, "req-exclude-region")
	if !found || message != "abgelehnt [PCF-E8034]: Spalte trägt eine Routing-Bedingung: public."+routeTable+".region" {
		t.Fatalf("exclude_column gegen die Routing-Spalte: failed = %t, Fehlertext %q", found, message)
	}
	if target := routedTarget(t, assembler, 1, "eu"); target != "eu" {
		t.Fatalf("Ziel nach dem abgelehnten Ausschluss = %q, wollen eu", target)
	}

	process(ctx, deps, queue, free)
	if !isApplied(queue, "req-exclude-note") {
		message, _ := failureOf(queue, "req-exclude-note")
		t.Fatalf("exclude_column gegen eine Spalte ohne Routing-Bedingung nicht applied, Fehlertext %q", message)
	}
}

// TestProcessAdministrationRequestsRouteOrderNotation trägt `SPEC-032` über den
// Antragsweg (Queue, Use Case, Assembler): `order` gilt als JSON-Zahl mit dem
// Wert einer positiven ganzen Zahl, gleich in welcher Schreibweise — der Antrag
// wird `applied` und die Regel weist ihr Ziel zu; ein Bruchteil, 0, ein
// negativer Wert, ein Wert über `MaxRouteOrder` und eine Zeichenkette enden
// `failed` mit `rule_spec ist ungültig`, und die Bindung bleibt ohne Ziel. Die
// Eingabe jedes Falls ist das `order`-Literal. Rot färbende Mutation: die
// Zahlenprüfung von `jsonRouteOrder` auf die Ziffernform des Rohtexts
// zurückführen — `10.0` und `1e1` enden `failed`.
func TestProcessAdministrationRequestsRouteOrderNotation(t *testing.T) {
	ctx := context.Background()
	table := "public." + routeTable
	for _, tc := range []struct {
		literal string
		applied bool
	}{
		{"10", true},
		{"10.0", true},
		{"1e1", true},
		{"1E+1", true},
		{"1.5", false},
		{"0", false},
		{"-1", false},
		{"2147483648", false},
		{`"5"`, false},
	} {
		request := routeRequest("req-order", "eu_orders", `{"target": "eu", "order": `+tc.literal+`, "when": {"column": "region", "equals": "eu"}}`)
		deps, queue, assembler := routeFixture(t, request)
		process(ctx, deps, queue, request)
		message, failed := failureOf(queue, "req-order")
		target := routedTarget(t, assembler, 1, "eu")
		if tc.applied {
			if !isApplied(queue, "req-order") || target != "eu" {
				t.Fatalf("order %s: applied = %t, Fehlertext %q, Ziel %q, wollen applied und Ziel eu", tc.literal, isApplied(queue, "req-order"), message, target)
			}
			continue
		}
		if !failed || message != "abgelehnt [PCF-E8011]: rule_spec ist ungültig: "+table+".eu_orders" || target != "" {
			t.Fatalf("order %s: failed = %t, Fehlertext %q, Ziel %q, wollen failed mit rule_spec ist ungültig und ohne Ziel", tc.literal, failed, message, target)
		}
	}
}

// TestProcessAdministrationRequestsRouteWithoutBindingIsApplied trägt die
// Zusage des Antrags gegen eine Tabelle ohne laufende Bindung (`SPEC-019`): er
// endet `applied`, der Nachtrag ist wirkungslos, und die Regel entsteht beim
// Anlegen der Bindung aus der dauerhaften Herkunft. Rot färbende Mutation:
// im Zweig `set_route` auf eine fehlende Bindung mit einem Fehler antworten —
// der Antrag endet `failed`.
func TestProcessAdministrationRequestsRouteWithoutBindingIsApplied(t *testing.T) {
	ctx := context.Background()
	rule := routeRequest("req-unbound", "eu_orders", routeSpecText("eu", 10, "region", "eu"))
	deps, queue, assembler := routeFixture(t, rule)
	assembler.RemoveBinding("public." + routeTable)

	process(ctx, deps, queue, rule)
	if !isApplied(queue, "req-unbound") {
		message, _ := failureOf(queue, "req-unbound")
		t.Fatalf("set_route gegen eine Tabelle ohne Bindung nicht applied, Fehlertext %q", message)
	}
	if assemblerCapturesQualified(t, assembler, 1, "public", routeTable) {
		t.Fatal("der Nachtrag hat eine Bindung angelegt")
	}
}

// TestActivatedTableBindingsCarriesRoutes trägt den Startpfad des dauerhaften
// Routing-Regelstandes: der Bindungs-Neuaufbau liest den Regelstand aus seiner
// dauerhaften Herkunft und übergibt ihn je Tabelle an die Bindung — ohne diesen
// Schritt liefert ein neu gestarteter Prozess kein Ziel, obwohl der Antrag
// `applied` trägt. Der Test führt den Neuaufbau bis in die Wirkung: die neu
// gebaute Bindung weist das Ziel zu, die Nachbartabelle bleibt ohne Ziel. Rot
// färbende Mutation: `Routes: routes[table.QualifiedName()]` in
// `activatedTableBindings` streichen bzw. den Schlüssel gegen einen festen
// ersetzen.
func TestActivatedTableBindingsCarriesRoutes(t *testing.T) {
	ctx := context.Background()
	const (
		source  = model.SourceID("src-restart")
		tableID = model.SourceTableID("tbl-restart")
		otherID = model.SourceTableID("tbl-restart-other")
	)
	activation := &fakeTableActivationPort{listed: []model.SourceTable{
		{ID: tableID, SourceID: source, Schema: "public", Table: routeTable},
		{ID: otherID, SourceID: source, Schema: "public", Table: routeTable + "_other"},
	}}
	schemaStore := &fakeSchemaStorePort{versions: map[model.SourceTableID]model.SchemaVersion{
		tableID: {ID: "sv-restart", SourceTableID: tableID, Version: 1},
		otherID: {ID: "sv-restart-other", SourceTableID: otherID, Version: 1},
	}}
	rule, err := model.NewRouteRule("eu_orders", "eu", 10, &model.RouteCondition{Column: "region", Equals: "eu"})
	if err != nil {
		t.Fatal(err)
	}
	routing := &fakeRoutingPort{rules: map[string][]model.RouteRule{"public." + routeTable: {rule}}}

	tables, err := activatedTableBindings(ctx, activation, schemaStore, &fakeColumnExclusionPort{}, &fakeTransformationPort{}, routing, source)
	if err != nil {
		t.Fatalf("activatedTableBindings = %v, wollen nil", err)
	}
	if got := tables["public."+routeTable].Routes; len(got) != 1 || got[0].Name() != "eu_orders" {
		t.Fatalf("Routes = %v, wollen die Regel eu_orders", got)
	}
	if got := tables["public."+routeTable+"_other"].Routes; len(got) != 0 {
		t.Fatalf("Routes der Nachbartabelle = %v, wollen keine", got)
	}
	assembler, err := mapper.NewAssembler(source, tables, nil)
	if err != nil {
		t.Fatalf("NewAssembler: %v", err)
	}
	if target := routedTarget(t, assembler, 1, "eu"); target != "eu" {
		t.Fatalf("Ziel der neu gebauten Bindung = %q, wollen eu", target)
	}
}

// TestProcessAdministrationRequestsDisableEnableCycleRestoresRoutes trägt den
// zweiten Auslöser des Stand-Verlusts: die Deaktivierung entfernt den
// Bindungs-Eintrag samt Regelstand, der Aktivierungs-Zweig legt ihn über die
// dauerhafte Herkunft neu an. Rot färbende Mutation: `Routes: routes[qualified]`
// im Aktivierungs-Zweig streichen — nach `disable` → `enable` steht eine
// Bindung ohne Regel.
func TestProcessAdministrationRequestsDisableEnableCycleRestoresRoutes(t *testing.T) {
	ctx := context.Background()
	tableID := administrationTableID("public", routeTable)
	deps, queue, assembler := routeFixture(t)
	deps.activation = &fakeTableActivationPort{registered: map[string]model.SourceTable{
		"public." + routeTable: {ID: tableID, SourceID: "src-admin", Schema: "public", Table: routeTable},
	}}
	deps.schemaStore = &fakeSchemaStorePort{versions: map[model.SourceTableID]model.SchemaVersion{
		tableID: {ID: administrationSchemaVersionID(tableID), SourceTableID: tableID, Version: 1},
	}}
	rule, err := model.NewRouteRule("eu_orders", "eu", 10, &model.RouteCondition{Column: "region", Equals: "eu"})
	if err != nil {
		t.Fatal(err)
	}
	deps.routing = &fakeRoutingPort{rules: map[string][]model.RouteRule{"public." + routeTable: {rule}}}
	assembler.SetRoute("public."+routeTable, rule)
	if target := routedTarget(t, assembler, 1, "eu"); target != "eu" {
		t.Fatalf("Vorbedingung verletzt: Ziel mit Regel = %q", target)
	}

	process(ctx, deps, queue, model.AdministrationRequest{
		ID: "req-cycle-disable", Source: "src-admin", Schema: "public", Table: routeTable, Kind: model.AdministrationRequestDisable,
	})
	if assemblerCapturesQualified(t, assembler, 2, "public", routeTable) {
		t.Fatal("Assembler trägt nach Disable weiterhin eine Bindung")
	}
	process(ctx, deps, queue, model.AdministrationRequest{
		ID: "req-cycle-enable", Source: "src-admin", Schema: "public", Table: routeTable, Kind: model.AdministrationRequestEnable,
	})
	if target := routedTarget(t, assembler, 3, "eu"); target != "eu" {
		t.Fatalf("Ziel nach dem disable/enable-Zyklus = %q, wollen eu (Regelstand aus der Herkunft)", target)
	}
}

// TestProcessAdministrationRequestsMarksFailedWhenRoutingStateReadFails trägt
// den Fehlerpfad der dauerhaften Herkunft im Aktivierungs-Zweig: ein Lesefehler
// des Routing-Regelstandes endet im `failed`-Vermerk, statt eine Bindung ohne
// den geführten Stand anzulegen.
func TestProcessAdministrationRequestsMarksFailedWhenRoutingStateReadFails(t *testing.T) {
	ctx := context.Background()
	tableID := administrationTableID("public", routeTable)
	wantErr := stderrors.New("Antrags-Historie nicht lesbar")
	deps, queue, assembler := routeFixture(t)
	assembler.RemoveBinding("public." + routeTable)
	deps.activation = &fakeTableActivationPort{registered: map[string]model.SourceTable{
		"public." + routeTable: {ID: tableID, SourceID: "src-admin", Schema: "public", Table: routeTable},
	}}
	deps.schemaStore = &fakeSchemaStorePort{versions: map[model.SourceTableID]model.SchemaVersion{
		tableID: {ID: administrationSchemaVersionID(tableID), SourceTableID: tableID, Version: 1},
	}}
	deps.routing = &fakeRoutingPort{err: wantErr}

	process(ctx, deps, queue, model.AdministrationRequest{
		ID: "req-routes-read-fail", Source: "src-admin", Schema: "public", Table: routeTable, Kind: model.AdministrationRequestEnable,
	})
	if isApplied(queue, "req-routes-read-fail") {
		t.Fatal("der Antrag ist applied, obwohl der Routing-Regelstand nicht lesbar war")
	}
	if message, found := failureOf(queue, "req-routes-read-fail"); !found || message != "Fehlerklasse internal [PCF-E7000]: "+wantErr.Error() {
		t.Fatalf("Fehlertext = %q (vermerkt %v), wollen den Kopf der Klasse internal vor %q", message, found, wantErr.Error())
	}
	if assemblerCapturesQualified(t, assembler, 1, "public", routeTable) {
		t.Fatal("Assembler trägt nach gescheitertem Regelstand-Lesen eine Bindung")
	}
}
