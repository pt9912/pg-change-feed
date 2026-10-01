package excludecolumn_test

import (
	"context"
	stderrors "errors"
	"testing"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/excludecolumn"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// fakeColumnExclusion trägt den `ColumnExclusionPort` als Fake (`ADR-0030`):
// er meldet die vorgegebene Spaltenexistenz und merkt sich die angefragte
// Adresse, damit der Test belegt, dass der Use Case genau die Spalte des
// Kommandos prüft. `existsNotFor` bindet die Ablehnung an genau eine
// Spaltenadresse (LP2): jede andere Adresse trägt derselbe Fake.
type fakeColumnExclusion struct {
	exists bool
	err    error

	existsNotFor map[string]bool // schema.table.column → existiert nicht

	calls  int
	schema string
	table  string
	column string
}

func (f *fakeColumnExclusion) ColumnExists(ctx context.Context, schema, table, column string) (bool, error) {
	f.calls++
	f.schema, f.table, f.column = schema, table, column
	if f.existsNotFor[schema+"."+table+"."+column] {
		return false, nil
	}
	return f.exists, f.err
}

// ExcludedColumns bleibt ungenutzt: der Use Case liest die Spaltenexistenz;
// den dauerhaften Ausschlussstand trägt die Verdrahtung (`ADR-0065`).
func (f *fakeColumnExclusion) ExcludedColumns(ctx context.Context, source model.SourceID) (map[string][]string, error) {
	return nil, nil
}

// fakeRouting trägt den `RoutingPort` als Fake (`ADR-0030`): der Regelstand ist
// an die Quelle gebunden, die der Aufruf nennt, und zählt die Lesungen.
type fakeRouting struct {
	source model.SourceID
	rules  map[string][]model.RouteRule // schema.table → Regelstand
	err    error
	calls  int
}

func (f *fakeRouting) RoutingRules(ctx context.Context, source model.SourceID) (map[string][]model.RouteRule, error) {
	f.calls++
	if f.err != nil {
		return nil, f.err
	}
	if source != f.source {
		return nil, nil
	}
	return f.rules, nil
}

// TestExcludeColumnChecksSourceColumn trägt den Happy Path (`LH-FA-CFG-005`):
// eine an der Quelle vorhandene Spalte wird geprüft, der Aufruf endet ohne
// Fehler.
func TestExcludeColumnChecksSourceColumn(t *testing.T) {
	columns := &fakeColumnExclusion{exists: true}
	service := excludecolumn.NewExcludeColumnService(columns, &fakeRouting{})

	if err := service.Exclude(context.Background(), excludecolumn.ExcludeColumnCommand{
		Source: "src-1", Schema: "public", Table: "orders", Column: "secret",
	}); err != nil {
		t.Fatalf("Exclude = %v, wollen nil", err)
	}
	if columns.calls != 1 {
		t.Fatalf("ColumnExists-Aufrufe = %d, wollen 1", columns.calls)
	}
	if columns.schema != "public" || columns.table != "orders" || columns.column != "secret" {
		t.Fatalf("geprüfte Adresse = %s.%s.%s, wollen public.orders.secret", columns.schema, columns.table, columns.column)
	}
}

// TestExcludeColumnRejectsMissingSourceColumn trägt den Negative-Pfad
// (`LH-FA-CFG-005` Negative: „folgt ein expliziter Fehlerpfad"): eine nicht
// existierende Spalte endet über `ErrSourceColumnMissing` — ohne dieses
// Sentinel bliebe der Antrag still erfolgreich und wirkte nie. Die Ablehnung
// ist an die Spaltenadresse des Kommandos gebunden: der Fake kennt genau die
// fehlende Adresse, dieselbe Anlage trägt die vorhandene Spalte durch.
func TestExcludeColumnRejectsMissingSourceColumn(t *testing.T) {
	columns := &fakeColumnExclusion{exists: true, existsNotFor: map[string]bool{"public.orders.does_not_exist": true}}
	service := excludecolumn.NewExcludeColumnService(columns, &fakeRouting{})

	err := service.Exclude(context.Background(), excludecolumn.ExcludeColumnCommand{
		Source: "src-1", Schema: "public", Table: "orders", Column: "does_not_exist",
	})
	if !stderrors.Is(err, inbound.ErrSourceColumnMissing) {
		t.Fatalf("Fehler = %v, wollen ErrSourceColumnMissing", err)
	}
	if columns.schema != "public" || columns.table != "orders" || columns.column != "does_not_exist" {
		t.Fatalf("geprüfte Adresse = %s.%s.%s, wollen public.orders.does_not_exist", columns.schema, columns.table, columns.column)
	}

	if err := service.Exclude(context.Background(), excludecolumn.ExcludeColumnCommand{
		Source: "src-1", Schema: "public", Table: "orders", Column: "secret",
	}); err != nil {
		t.Fatalf("vorhandene Spalte: %v", err)
	}
}

// TestExcludeColumnPropagatesPortError trägt den Adapter-Fehlerpfad: ein
// Fehler der Katalog-Prüfung wird unverändert durchgereicht, nicht als
// fehlende Spalte fehlinterpretiert.
func TestExcludeColumnPropagatesPortError(t *testing.T) {
	wantErr := stderrors.New("Katalog nicht lesbar")
	service := excludecolumn.NewExcludeColumnService(&fakeColumnExclusion{err: wantErr}, &fakeRouting{})

	err := service.Exclude(context.Background(), excludecolumn.ExcludeColumnCommand{
		Source: "src-1", Schema: "public", Table: "orders", Column: "secret",
	})
	if !stderrors.Is(err, wantErr) {
		t.Fatalf("Fehler = %v, wollen %v", err, wantErr)
	}
	if stderrors.Is(err, inbound.ErrSourceColumnMissing) {
		t.Fatalf("Fehler = %v, darf nicht als fehlende Spalte gelesen werden", err)
	}
}

func routingState(t *testing.T) *fakeRouting {
	t.Helper()
	eu, err := model.NewRouteRule("eu_orders", "eu", 10, &model.RouteCondition{Column: "region", Equals: "eu"})
	if err != nil {
		t.Fatal(err)
	}
	closing, err := model.NewRouteRule("rest", "sonstige", 100, nil)
	if err != nil {
		t.Fatal(err)
	}
	other, err := model.NewRouteRule("other_rule", "o", 5, &model.RouteCondition{Column: "secret", Equals: "x"})
	if err != nil {
		t.Fatal(err)
	}
	return &fakeRouting{
		source: "src-1",
		rules: map[string][]model.RouteRule{
			"public.orders": {eu, closing},
			"public.other":  {other},
		},
	}
}

func exclude(routes *fakeRouting, source model.SourceID, column string) error {
	return excludecolumn.NewExcludeColumnService(&fakeColumnExclusion{exists: true}, routes).Exclude(context.Background(), excludecolumn.ExcludeColumnCommand{
		Source: source, Schema: "public", Table: "orders", Column: column,
	})
}

// TestExcludeColumnRejectsAColumnWithARouteCondition trägt R3 in der
// Gegenrichtung (`SPEC-019`): eine Spalte, die eine
// Routing-Bedingung der Tabelle trägt, wird nicht ausgeschlossen; der
// Fehlertext ist der der Spec (Klartext, Doppelpunkt, Adresse
// `schema.table.column`). Die Eingabe ist die Verletzung: die Spalte `region`
// der Bedingung. Die Gegenproben tragen die Grenzen — eine Spalte ohne
// Bedingung (der Bestandspfad), die Bedingung einer Nachbartabelle auf
// dieselbe Spalte, der Regelstand einer Nachbarquelle, eine Regel ohne `when`
// (trägt keine Spalte). Rot färbende Mutation: `RoutesConditionOn` im Use Case
// streichen (die Verletzung wird angenommen); `state[command.Schema+"."+
// command.Table]` gegen `state["public.other"]` tauschen (die Gegenprobe
// `secret` wird abgelehnt, `region` angenommen).
func TestExcludeColumnRejectsAColumnWithARouteCondition(t *testing.T) {
	err := exclude(routingState(t), "src-1", "region")
	if err == nil || err.Error() != "Spalte trägt eine Routing-Bedingung: public.orders.region" {
		t.Fatalf("Exclude(region) = %v, wollen die Ablehnung mit dem Text der Spec", err)
	}
	if !stderrors.Is(err, domainerrors.ErrColumnHasRouteCondition) {
		t.Fatalf("Fehler %v löst nicht auf ErrColumnHasRouteCondition auf", err)
	}
	for _, tc := range []struct {
		label  string
		source model.SourceID
		column string
	}{
		{"Spalte ohne Routing-Bedingung (Bestandspfad)", "src-1", "name"},
		{"Bedingung nur an der Nachbartabelle", "src-1", "secret"},
		{"Regelstand einer Nachbarquelle", "src-2", "region"},
	} {
		if err := exclude(routingState(t), tc.source, tc.column); err != nil {
			t.Fatalf("%s: Exclude(%s) = %v, wollen nil", tc.label, tc.column, err)
		}
	}
}

// TestExcludeColumnReadsTheRoutingStateOnlyForAnExistingColumn trägt: eine
// nicht vorhandene Spalte endet vor der Lesung des Regelstands mit dem
// Bestandstext, ein Lesefehler des Regelstands endet unverändert. Rot färbende
// Mutation: die Regelstand-Lesung vor `ColumnExists` ziehen (der Zähler färbt
// rot); den Lesefehler verwerfen (der Fall endet angenommen).
func TestExcludeColumnReadsTheRoutingStateOnlyForAnExistingColumn(t *testing.T) {
	routes := routingState(t)
	service := excludecolumn.NewExcludeColumnService(&fakeColumnExclusion{exists: false}, routes)
	err := service.Exclude(context.Background(), excludecolumn.ExcludeColumnCommand{
		Source: "src-1", Schema: "public", Table: "orders", Column: "region",
	})
	if !stderrors.Is(err, inbound.ErrSourceColumnMissing) || routes.calls != 0 {
		t.Fatalf("Fehler = %v, Lesungen = %d, wollen ErrSourceColumnMissing ohne Lesung", err, routes.calls)
	}
	wantErr := stderrors.New("Antrags-Historie nicht lesbar")
	routes = routingState(t)
	routes.err = wantErr
	if err := exclude(routes, "src-1", "name"); !stderrors.Is(err, wantErr) {
		t.Fatalf("Fehler = %v, wollen %v", err, wantErr)
	}
}
