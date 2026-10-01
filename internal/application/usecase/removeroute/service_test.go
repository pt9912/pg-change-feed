package removeroute_test

import (
	"context"
	stderrors "errors"
	"testing"

	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/removeroute"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// fakeRouting trägt den `RoutingPort` als Fake (`ADR-0030`): der Regelstand ist
// an Quelle und Tabelle gebunden, die der Aufruf nennt.
type fakeRouting struct {
	source model.SourceID
	rules  map[string][]model.RouteRule
	err    error
	calls  int
}

func (f *fakeRouting) RoutingRules(ctx context.Context, source model.SourceID) (map[string][]model.RouteRule, error) {
	f.calls++
	if f.err != nil {
		return nil, f.err
	}
	if source != f.source {
		return map[string][]model.RouteRule{}, nil
	}
	return f.rules, nil
}

var _ outbound.RoutingPort = (*fakeRouting)(nil)

func newPort(t *testing.T) *fakeRouting {
	eu, err := model.NewRouteRule("eu_orders", "eu", 10, &model.RouteCondition{Column: "region", Equals: "eu"})
	if err != nil {
		t.Fatalf("NewRouteRule: %v", err)
	}
	other, err := model.NewRouteRule("other_rule", "o", 5, nil)
	if err != nil {
		t.Fatalf("NewRouteRule: %v", err)
	}
	return &fakeRouting{
		source: "src-1",
		rules: map[string][]model.RouteRule{
			"public.orders": {eu},
			"public.other":  {other},
		},
	}
}

func remove(port *fakeRouting, source model.SourceID, ruleName string) error {
	return removeroute.NewRemoveRouteService(port).Remove(context.Background(), removeroute.RemoveRouteCommand{
		Source: source, Schema: "public", Table: "orders", RuleName: ruleName,
	})
}

// TestRemoveRouteAcceptsAKeptRule trägt den Happy Path (`LH-FA-CFG-008`): der
// Regelname gehört zum Routing-Regelstand der Tabelle.
func TestRemoveRouteAcceptsAKeptRule(t *testing.T) {
	if err := remove(newPort(t), "src-1", "eu_orders"); err != nil {
		t.Fatalf("Remove = %v, wollen nil", err)
	}
}

// TestRemoveRouteRejectsWithTheSpecTexts trägt die beiden Zeilen, die
// `remove_route` durchläuft (`SPEC-019`): die Zeile zum Regelnamen und R6 —
// jeweils mit dem exakten Fehlertext (Klartext, Doppelpunkt, Adresse
// `schema.table.rule_name`) und dem Grund; ein Name, den nur die
// Nachbartabelle bzw. die Nachbarquelle führt, ist an dieser Tabelle nicht
// geführt. Rot färbende Mutation: die Namensprüfung `rule.Name() ==
// command.RuleName` durch `true` ersetzen (der nicht geführte Name wird
// angenommen), `state[command.Schema+"."+command.Table]` durch
// `state["public.other"]` ersetzen (der Nachbarname wird angenommen, der
// geführte abgelehnt), `CheckRuleName` streichen (leerer Name endet als
// `Regelname nicht geführt`).
func TestRemoveRouteRejectsWithTheSpecTexts(t *testing.T) {
	for _, tc := range []struct {
		label    string
		source   model.SourceID
		ruleName string
		text     string
		reason   error
	}{
		{"Regelname leer", "src-1", "", "Regelname ist ungültig: public.orders.", domainerrors.ErrInvalidRuleName},
		{"Regelname mit Großbuchstaben", "src-1", "Eu_Orders", "Regelname ist ungültig: public.orders.Eu_Orders", domainerrors.ErrInvalidRuleName},
		{"R6 Name nicht geführt", "src-1", "gibt_es_nicht", "Regelname nicht geführt: public.orders.gibt_es_nicht", domainerrors.ErrRuleNotKept},
		{"R6 Name nur an der Nachbartabelle", "src-1", "other_rule", "Regelname nicht geführt: public.orders.other_rule", domainerrors.ErrRuleNotKept},
		{"R6 Name nur an der Nachbarquelle", "src-2", "eu_orders", "Regelname nicht geführt: public.orders.eu_orders", domainerrors.ErrRuleNotKept},
	} {
		err := remove(newPort(t), tc.source, tc.ruleName)
		if err == nil {
			t.Errorf("%s: Remove = nil, wollen %q", tc.label, tc.text)
			continue
		}
		if err.Error() != tc.text {
			t.Errorf("%s: Fehlertext = %q, wollen %q", tc.label, err.Error(), tc.text)
		}
		if !stderrors.Is(err, tc.reason) {
			t.Errorf("%s: Fehler = %v, wollen Grund %v", tc.label, err, tc.reason)
		}
	}
}

// TestRemoveRouteChecksTheNameBeforeReadingTheStore trägt: ein Name außerhalb
// des Alphabets endet, ohne dass der Use Case den Regelstand liest. Rot
// färbende Mutation: `CheckRuleName` hinter das Lesen stellen.
func TestRemoveRouteChecksTheNameBeforeReadingTheStore(t *testing.T) {
	port := newPort(t)
	if err := remove(port, "src-1", ""); err == nil {
		t.Fatal("Remove mit leerem Namen = nil, wollen einen Fehler")
	}
	if port.calls != 0 {
		t.Fatalf("Port-Aufrufe = %d, wollen 0", port.calls)
	}
}

// TestRemoveRouteReturnsPortErrors trägt: ein Lesefehler des Ports endet
// unverändert.
func TestRemoveRouteReturnsPortErrors(t *testing.T) {
	wantErr := stderrors.New("Antrags-Historie nicht lesbar")
	port := newPort(t)
	port.err = wantErr
	if err := remove(port, "src-1", "eu_orders"); !stderrors.Is(err, wantErr) {
		t.Fatalf("Fehler = %v, wollen %v", err, wantErr)
	}
}
