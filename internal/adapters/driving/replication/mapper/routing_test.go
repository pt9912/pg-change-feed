package mapper_test

import (
	"context"
	stderrors "errors"
	"sync"
	"testing"
	"time"

	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/decode"
	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/mapper"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Die Tests dieser Datei tragen den Routing-Regelstand einer Bindung
// (`LH-FA-CFG-008`): Auswertung je Operation und Bildbasis, Reihenfolge,
// Abwesenheit, Quellwert vor der Transformation, Anwendbarkeits-Prüfung vor
// jeder Serialisierung, Live-Reload des Regelstands. Ohne Antragsweg entsteht
// der Regelstand hier über die Bindung und `SetRoute`.

// routeRule legt eine Routing-Regel an oder bricht den Test ab.
func routeRule(t testing.TB, name, target string, order int64, column, equals string) model.RouteRule {
	t.Helper()
	var when *model.RouteCondition
	if column != "" {
		when = &model.RouteCondition{Column: column, Equals: equals}
	}
	rule, err := model.NewRouteRule(name, target, order, when)
	if err != nil {
		t.Fatalf("NewRouteRule(%q, %q, %d, %q): %v", name, target, order, column, err)
	}
	return rule
}

// routeTables trägt `public.orders` mit Ausschluss, Transformationen und
// Routing-Regelstand.
func routeTables(excluded []string, transformations []model.Transformation, routes ...model.RouteRule) map[string]mapper.TableBinding {
	return map[string]mapper.TableBinding{
		"public.orders": {TableID: "tbl-1", SchemaVersion: "sv-1", ExcludedColumns: excluded, Transformations: transformations, Routes: routes},
	}
}

// ordersRelation trägt die Relation der Routing-Tests: id, name, region.
func ordersRelation() *decode.Relation {
	return relation("public", "orders",
		decode.Column{Name: "id", Key: true},
		decode.Column{Name: "name"},
		decode.Column{Name: "region"})
}

// Die Auswertung liest bei INSERT und UPDATE das Neu-Bild, bei DELETE das
// Alt-Bild; die kleinere `order` gewinnt vor der Listenposition; keine
// treffende Regel lässt das Ziel leer. Rot färbende Mutationen:
// `routeValues` in `Assembler.change` bei DELETE nicht auf `event.Old` setzen
// (DELETE-Fall), `binding.Routes` durch `nil` ersetzen (alle Fälle mit Ziel),
// in `model.EvaluateRoute` den `order`-Vergleich entfernen (Reihenfolge-Fall).
func TestConsumeRoutesByBildbasisAndOrder(t *testing.T) {
	ctx := context.Background()
	eu := routeRule(t, "eu", "eu", 10, "region", "eu")
	euLate := routeRule(t, "eu-spaet", "spaet", 30, "region", "eu")
	rest := routeRule(t, "rest", "sonstige", 100, "", "")
	assembler := newAssembler(t, routeTables(nil, nil, euLate, rest, eu))
	relationEvent := ordersRelation()

	insert := consumedChange(t, ctx, assembler, 1, decode.Change{
		Relation: relationEvent, Operation: decode.OpInsert,
		New: []*string{pointer("1"), pointer("Ada"), pointer("eu")},
	})
	if insert.RouteTarget != "eu" {
		t.Fatalf("INSERT: RouteTarget = %q, wollen eu (kleinere order vor Listenposition)", insert.RouteTarget)
	}

	update := consumedChange(t, ctx, assembler, 2, decode.Change{
		Relation: relationEvent, Operation: decode.OpUpdate,
		Old: []*string{pointer("1"), pointer("Ada"), pointer("eu")},
		New: []*string{pointer("1"), pointer("Ada"), pointer("us")},
	})
	if update.RouteTarget != "sonstige" {
		t.Fatalf("UPDATE: RouteTarget = %q, wollen sonstige (Bildbasis ist das Neu-Bild us, nicht das Alt-Bild eu)", update.RouteTarget)
	}

	deletion := consumedChange(t, ctx, assembler, 3, decode.Change{
		Relation: relationEvent, Operation: decode.OpDelete,
		Old: []*string{pointer("1"), pointer("Ada"), pointer("eu")},
	})
	if deletion.RouteTarget != "eu" {
		t.Fatalf("DELETE: RouteTarget = %q, wollen eu (Bildbasis ist das Alt-Bild)", deletion.RouteTarget)
	}

	unrouted := consumedChange(t, ctx, newAssembler(t, routeTables(nil, nil, eu)), 4, decode.Change{
		Relation: relationEvent, Operation: decode.OpInsert,
		New: []*string{pointer("1"), pointer("Ada"), pointer("us")},
	})
	if unrouted.RouteTarget != "" {
		t.Fatalf("ohne treffende Regel: RouteTarget = %q, wollen leer", unrouted.RouteTarget)
	}
}

// Ein abwesender Wert (NULL, unverändertes TOAST, Bild ohne die Spalte, bei
// DELETE ohne volle Replica-Identität jede Nicht-Schlüsselspalte) ist ein
// Nicht-Treffer ohne Fehler; die nächste Regel wird geprüft. Rot färbende
// Mutation: in `RouteRule.matches` die Prüfung `values[i] != nil` entfernen
// (Nullzeiger-Zugriff).
func TestConsumeRoutingTreatsAbsentValueAsNoMatch(t *testing.T) {
	ctx := context.Background()
	eu := routeRule(t, "eu", "eu", 10, "region", "eu")
	rest := routeRule(t, "rest", "sonstige", 100, "", "")
	assembler := newAssembler(t, routeTables(nil, nil, eu, rest))
	relationEvent := ordersRelation()

	nullValue := consumedChange(t, ctx, assembler, 1, decode.Change{
		Relation: relationEvent, Operation: decode.OpInsert,
		New: []*string{pointer("1"), pointer("Ada"), nil},
	})
	if nullValue.RouteTarget != "sonstige" {
		t.Fatalf("NULL: RouteTarget = %q, wollen sonstige", nullValue.RouteTarget)
	}
	keyOnlyDelete := consumedChange(t, ctx, assembler, 2, decode.Change{
		Relation: relationEvent, Operation: decode.OpDelete,
		Old: []*string{pointer("1"), nil, nil},
	})
	if keyOnlyDelete.RouteTarget != "sonstige" {
		t.Fatalf("DELETE nur mit Schlüssel: RouteTarget = %q, wollen sonstige", keyOnlyDelete.RouteTarget)
	}
	short := consumedChange(t, ctx, assembler, 3, decode.Change{
		Relation: relationEvent, Operation: decode.OpInsert,
		New: []*string{pointer("1")},
	})
	if short.RouteTarget != "sonstige" {
		t.Fatalf("Bild ohne die Spalte: RouteTarget = %q, wollen sonstige", short.RouteTarget)
	}

	withoutFallback := newAssembler(t, routeTables(nil, nil, eu))
	none := consumedChange(t, ctx, withoutFallback, 4, decode.Change{
		Relation: relationEvent, Operation: decode.OpInsert,
		New: []*string{pointer("1"), pointer("Ada"), nil},
	})
	if none.RouteTarget != "" {
		t.Fatalf("NULL ohne Abschlussregel: RouteTarget = %q, wollen leer", none.RouteTarget)
	}
}

// Die Bedingung liest den Quellwert vor `rename_column` und `map_value`: eine
// Regel auf den Quellnamen und den Quellwert trifft, obwohl das Bild den
// Zielnamen und den abgebildeten Wert trägt; eine Regel auf den Zielnamen oder
// den abgebildeten Wert trifft nicht, und das Bild bleibt von der Regel
// unberührt. Rot färbende Mutation: in `Assembler.change` die Bedingung
// gegen das serialisierte Bild statt gegen `event.New` auswerten.
func TestConsumeRoutingReadsSourceValueBeforeTransformations(t *testing.T) {
	ctx := context.Background()
	transformations := []model.Transformation{
		renameRule(t, "gebiet", "region", "area"),
		mapValueRule(t, "namen", "name", map[string]string{"Ada": "Lovelace"}),
	}
	bySource := routeRule(t, "quelle", "quelle", 10, "region", "eu")
	byMappedValue := routeRule(t, "abgebildet", "abgebildet", 30, "name", "Lovelace")
	event := decode.Change{
		Relation: ordersRelation(), Operation: decode.OpInsert,
		New: []*string{pointer("1"), pointer("Ada"), pointer("eu")},
	}

	routed := consumedChange(t, ctx, newAssembler(t, routeTables(nil, transformations, bySource)), 1, event)
	if routed.RouteTarget != "quelle" {
		t.Fatalf("Regel auf Quellname und Quellwert: RouteTarget = %q, wollen quelle", routed.RouteTarget)
	}
	if string(routed.NewImage) != `{"id":"1","name":"Lovelace","area":"eu"}` {
		t.Fatalf("Neu-Image: %s, wollen die transformierte Form unverändert", routed.NewImage)
	}
	unrouted := consumedChange(t, ctx, newAssembler(t, routeTables(nil, transformations, byMappedValue)), 2, event)
	if unrouted.RouteTarget != "" {
		t.Fatalf("Regel auf den abgebildeten Wert: RouteTarget = %q, wollen leer", unrouted.RouteTarget)
	}
	// Der Zielname einer Umbenennung ist keine Spalte der Relation: eine
	// Bedingung darauf ist nicht anwendbar.
	byTargetName := routeRule(t, "zielname", "zielname", 20, "area", "eu")
	assembler := newAssembler(t, routeTables(nil, transformations, byTargetName))
	if _, err := assembler.Consume(ctx, decode.Begin{XID: 3}); err != nil {
		t.Fatalf("Begin: %v", err)
	}
	if _, err := assembler.Consume(ctx, event); !stderrors.Is(err, mapper.ErrRoutingNotApplicable) {
		t.Fatalf("Regel auf den Zielnamen: %v, wollen ErrRoutingNotApplicable", err)
	}
}

// Eine Regel mit Bedingung auf eine Spalte, die in der Relation fehlt, ist
// nicht anwendbar: `Consume` meldet `ErrRoutingNotApplicable` statt eines
// Changes, die Sequenz rückt nicht vor, und nach dem Entfernen der Regel läuft
// dieselbe Transaktion mit Sequenz 1 weiter. Rot färbende Mutation: den Aufruf
// `checkRoutes` in `Assembler.change` entfernen.
func TestConsumeRoutingRuleOnMissingColumnIsNotApplicable(t *testing.T) {
	ctx := context.Background()
	missing := routeRule(t, "fehlt", "x", 10, "fehlende_spalte", "eu")
	assembler := newAssembler(t, routeTables(nil, nil, missing))

	if _, err := assembler.Consume(ctx, decode.Begin{XID: 7}); err != nil {
		t.Fatalf("Begin: %v", err)
	}
	_, err := assembler.Consume(ctx, decode.Change{
		Relation: ordersRelation(), Operation: decode.OpInsert,
		New: []*string{pointer("1"), pointer("Ada"), pointer("eu")},
	})
	if !stderrors.Is(err, mapper.ErrRoutingNotApplicable) {
		t.Fatalf("Fehler = %v, wollen ErrRoutingNotApplicable", err)
	}

	assembler.RemoveRoute("public.orders", "fehlt")
	if _, err := assembler.Consume(ctx, decode.Change{
		Relation: ordersRelation(), Operation: decode.OpInsert,
		New: []*string{pointer("1"), pointer("Ada"), pointer("eu")},
	}); err != nil {
		t.Fatalf("Change nach dem Entfernen der Regel: %v", err)
	}
	command, err := assembler.Consume(ctx, decode.Commit{CommitLSN: 700})
	if err != nil {
		t.Fatalf("Commit: %v", err)
	}
	changes, err := command.Transaction.Changes()
	if err != nil || len(changes) != 1 || changes[0].Sequence != 1 {
		t.Fatalf("Changes = %+v, %v, wollen genau eine Change mit Sequenz 1 (die fehlgeschlagene hat die Sequenz nicht belegt)", changes, err)
	}
}

// Eine Regel ohne Bedingung ist auf jede Relation anwendbar, auch auf eine,
// die keine Spalte der anderen Regeln trägt.
func TestConsumeRoutingRuleWithoutConditionIsAlwaysApplicable(t *testing.T) {
	ctx := context.Background()
	assembler := newAssembler(t, routeTables(nil, nil, routeRule(t, "rest", "alle", 1, "", "")))
	change := consumedChange(t, ctx, assembler, 1, decode.Change{
		Relation:  relation("public", "orders", decode.Column{Name: "id", Key: true}),
		Operation: decode.OpInsert, New: []*string{pointer("1")},
	})
	if change.RouteTarget != "alle" {
		t.Fatalf("RouteTarget = %q, wollen alle", change.RouteTarget)
	}
}

// V3 auf Unit-Ebene: ob `ErrRoutingNotApplicable` über den Erfassungspfad
// erreichbar ist, ohne die Relation-Prüfung zu umgehen. Eine entfernte Spalte
// endet an einer bekannten Spaltenform als `ErrIncompatibleSchemaChange`,
// bevor eine Change assembliert wird; trägt die aktuelle Version noch keine
// Spaltenform (Erstaktivierung), registriert die erste Relation-Nachricht ohne
// Vergleich, und die Change endet danach als `ErrRoutingNotApplicable`.
func TestRoutingNotApplicableReachabilityThroughRelationCheck(t *testing.T) {
	ctx := context.Background()
	rule := routeRule(t, "eu", "eu", 10, "region", "eu")
	withoutRegion := relation("public", "orders",
		decode.Column{Name: "id", Key: true, TypeOID: 23},
		decode.Column{Name: "name", TypeOID: 25})
	event := decode.Change{Relation: withoutRegion, Operation: decode.OpInsert, New: []*string{pointer("1"), pointer("Ada")}}

	t.Run("bekannte Spaltenform: die Relation-Prüfung meldet zuerst", func(t *testing.T) {
		store := newFakeSchemaStore()
		store.versions["tbl-1"] = model.SchemaVersion{ID: "sv-1", SourceTableID: "tbl-1", Version: 1}
		store.schemas["sv-1"] = model.TableSchema{VersionID: "sv-1", Columns: []model.Column{
			{Name: "id", OID: 23}, {Name: "name", OID: 25}, {Name: "region", OID: 25},
		}}
		assembler, err := mapper.NewAssembler("src-1", routeTables(nil, nil, rule), store)
		if err != nil {
			t.Fatalf("NewAssembler: %v", err)
		}
		if _, err := assembler.Consume(ctx, withoutRegion); !stderrors.Is(err, mapper.ErrIncompatibleSchemaChange) {
			t.Fatalf("Relation ohne die Spalte: %v, wollen ErrIncompatibleSchemaChange", err)
		}
	})

	t.Run("Erstaktivierung ohne Spaltenform: die Change meldet die Nichtanwendbarkeit", func(t *testing.T) {
		store := newFakeSchemaStore()
		store.versions["tbl-1"] = model.SchemaVersion{ID: "sv-1", SourceTableID: "tbl-1", Version: 1}
		assembler, err := mapper.NewAssembler("src-1", routeTables(nil, nil, rule), store)
		if err != nil {
			t.Fatalf("NewAssembler: %v", err)
		}
		if _, err := assembler.Consume(ctx, withoutRegion); err != nil {
			t.Fatalf("Relation ohne die Spalte bei unbekannter Spaltenform: %v, wollen keinen Fehler", err)
		}
		if _, err := assembler.Consume(ctx, decode.Begin{XID: 1}); err != nil {
			t.Fatalf("Begin: %v", err)
		}
		if _, err := assembler.Consume(ctx, event); !stderrors.Is(err, mapper.ErrRoutingNotApplicable) {
			t.Fatalf("Change: %v, wollen ErrRoutingNotApplicable", err)
		}
	})
}

// `SetRoute` ersetzt eine Regel gleichen Namens an ihrer Stelle, hängt eine
// neue an und wirkt auf die nächste Change; `RemoveRoute` nimmt sie heraus.
// Eine nicht getragene Bindung und ein nicht geführter Name bleiben ohne
// Wirkung. Die bei der Anlage übergebene Liste bleibt unverändert (Schnappschuss
// eines Lesers). Rot färbende Mutation: in `withRoute` den Zweig `replaced`
// entfernen — dann trägt die Liste nach dem Ersetzen zwei Regeln gleichen
// Namens.
func TestAssemblerSetAndRemoveRoute(t *testing.T) {
	ctx := context.Background()
	initial := []model.RouteRule{routeRule(t, "eu", "eu", 10, "region", "eu")}
	assembler := newAssembler(t, map[string]mapper.TableBinding{
		"public.orders": {TableID: "tbl-1", SchemaVersion: "sv-1", Routes: initial},
	})
	event := decode.Change{
		Relation: ordersRelation(), Operation: decode.OpInsert,
		New: []*string{pointer("1"), pointer("Ada"), pointer("eu")},
	}
	target := func(xid uint32) model.RouteTarget { return consumedChange(t, ctx, assembler, xid, event).RouteTarget }

	if got := target(1); got != "eu" {
		t.Fatalf("Ausgangsstand: %q, wollen eu", got)
	}
	assembler.SetRoute("public.orders", routeRule(t, "eu", "europa", 10, "region", "eu"))
	if got := target(2); got != "europa" {
		t.Fatalf("nach dem Ersetzen: %q, wollen europa", got)
	}
	assembler.SetRoute("public.orders", routeRule(t, "frueher", "frueh", 5, "region", "eu"))
	if got := target(3); got != "frueh" {
		t.Fatalf("nach dem Anhängen einer Regel mit kleinerer order: %q, wollen frueh", got)
	}
	assembler.RemoveRoute("public.orders", "frueher")
	assembler.RemoveRoute("public.orders", "unbekannt")
	if got := target(4); got != "europa" {
		t.Fatalf("nach dem Entfernen: %q, wollen europa", got)
	}
	assembler.SetRoute("public.unbekannt", routeRule(t, "x", "x", 1, "", ""))
	assembler.RemoveRoute("public.unbekannt", "x")
	assembler.RemoveRoute("public.orders", "eu")
	if got := target(5); got != "" {
		t.Fatalf("ohne Regel: %q, wollen leer", got)
	}
	if len(initial) != 1 || initial[0].Target() != "eu" {
		t.Fatalf("die bei der Anlage übergebene Liste wurde verändert: %+v", initial)
	}
}

// `AddBinding` bei getragener Bindung und der Nachtrag der Schema-Version
// lassen den Routing-Regelstand stehen. Rot färbende Mutation: in
// `AddBinding` die Zeile `binding.Routes = existing.Routes` entfernen.
func TestAssemblerRoutesSurviveAddBindingAndSchemaVersionBump(t *testing.T) {
	ctx := context.Background()
	store := newFakeSchemaStore()
	store.versions["tbl-1"] = model.SchemaVersion{ID: "sv-1", SourceTableID: "tbl-1", Version: 1}
	store.schemas["sv-1"] = model.TableSchema{VersionID: "sv-1", Columns: []model.Column{
		{Name: "id", OID: 23}, {Name: "name", OID: 25}, {Name: "region", OID: 25},
	}}
	assembler, err := mapper.NewAssembler("src-1", routeTables(nil, nil, routeRule(t, "eu", "eu", 10, "region", "eu")), store)
	if err != nil {
		t.Fatalf("NewAssembler: %v", err)
	}
	assembler.AddBinding("public.orders", mapper.TableBinding{TableID: "tbl-1", SchemaVersion: "sv-1"})

	extended := relation("public", "orders",
		decode.Column{Name: "id", Key: true, TypeOID: 23},
		decode.Column{Name: "name", TypeOID: 25},
		decode.Column{Name: "region", TypeOID: 25},
		decode.Column{Name: "extra", TypeOID: 25})
	if _, err := assembler.Consume(ctx, extended); err != nil {
		t.Fatalf("Relation mit Erweiterung: %v", err)
	}
	change := consumedChange(t, ctx, assembler, 1, decode.Change{
		Relation: extended, Operation: decode.OpInsert,
		New: []*string{pointer("1"), pointer("Ada"), pointer("eu"), nil},
	})
	if change.RouteTarget != "eu" || change.SchemaVersion == "sv-1" {
		t.Fatalf("RouteTarget = %q, Schema-Version = %q, wollen eu und eine angehobene Version", change.RouteTarget, change.SchemaVersion)
	}
}

// Die Fitness Function der Nebenläufigkeit (`LH-FA-ADM-003`): der Regelstand
// gehört einer zweiten Goroutine (Administration) zu schreiben, während die
// Capture-Goroutine liest. Jede erfasste Change trägt eines der beiden Ziele
// der vollständigen Stände, nie ein drittes; `go test -race` deckt den
// gleichzeitigen Zugriff. Rot färbende Mutation: in `SetRoute` das Sperren
// (`a.tablesMu.Lock()`) entfernen.
func TestAssemblerRoutesAreRaceFree(t *testing.T) {
	ctx := context.Background()
	assembler := newAssembler(t, routeTables(nil, nil))
	relationEvent := ordersRelation()
	rule := routeRule(t, "eu", "eu", 10, "region", "eu")

	const iterations = 300
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := uint32(1); i <= iterations; i++ {
			if _, err := assembler.Consume(ctx, decode.Begin{XID: i}); err != nil {
				t.Errorf("Begin: %v", err)
				return
			}
			if _, err := assembler.Consume(ctx, decode.Change{
				Relation: relationEvent, Operation: decode.OpInsert,
				New: []*string{pointer("1"), pointer("Ada"), pointer("eu")},
			}); err != nil {
				t.Errorf("Change: %v", err)
				return
			}
			command, err := assembler.Consume(ctx, decode.Commit{CommitLSN: uint64(i), CommitTime: time.Now()})
			if err != nil {
				t.Errorf("Commit: %v", err)
				return
			}
			changes, err := command.Transaction.Changes()
			if err != nil || len(changes) != 1 {
				t.Errorf("Changes: %v (%d)", err, len(changes))
				return
			}
			if got := changes[0].RouteTarget; got != "" && got != "eu" {
				t.Errorf("RouteTarget %q trägt keinen der beiden Stände", got)
				return
			}
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			assembler.SetRoute("public.orders", rule)
			assembler.RemoveRoute("public.orders", "eu")
		}
	}()
	wg.Wait()
}
