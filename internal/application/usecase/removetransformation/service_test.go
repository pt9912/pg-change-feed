package removetransformation_test

import (
	"context"
	stderrors "errors"
	"testing"

	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/removetransformation"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// fakeTransformationPort trägt den `TransformationPort` als Fake (`ADR-0030`):
// der Regelstand ist an Quelle und Tabelle gebunden, die der Aufruf nennt.
type fakeTransformationPort struct {
	source model.SourceID
	rules  map[string][]model.Transformation
	err    error
	calls  int
}

func (f *fakeTransformationPort) TransformationRules(ctx context.Context, source model.SourceID) (map[string][]model.Transformation, error) {
	f.calls++
	if f.err != nil {
		return nil, f.err
	}
	if source != f.source {
		return map[string][]model.Transformation{}, nil
	}
	return f.rules, nil
}

func (f *fakeTransformationPort) SourceColumns(ctx context.Context, schema, table string) ([]string, error) {
	return nil, nil
}

var _ outbound.TransformationPort = (*fakeTransformationPort)(nil)

func newPort(t *testing.T) *fakeTransformationPort {
	kunde, err := model.NewRenameColumn("kundenname", "name", "customer_name")
	if err != nil {
		t.Fatalf("NewRenameColumn: %v", err)
	}
	status, err := model.NewRenameColumn("statusname", "status", "state")
	if err != nil {
		t.Fatalf("NewRenameColumn: %v", err)
	}
	return &fakeTransformationPort{
		source: "src-1",
		rules: map[string][]model.Transformation{
			"public.orders": {kunde},
			"public.other":  {status},
		},
	}
}

func remove(port *fakeTransformationPort, source model.SourceID, ruleName string) error {
	return removetransformation.NewRemoveTransformationService(port).Remove(context.Background(), removetransformation.RemoveTransformationCommand{
		Source: source, Schema: "public", Table: "orders", RuleName: ruleName,
	})
}

// TestRemoveTransformationAcceptsAKeptRule trägt den Happy Path
// (`LH-FA-CFG-007`): der Regelname gehört zum Regelstand der Tabelle.
func TestRemoveTransformationAcceptsAKeptRule(t *testing.T) {
	if err := remove(newPort(t), "src-1", "kundenname"); err != nil {
		t.Fatalf("Remove = %v, wollen nil", err)
	}
}

// TestRemoveTransformationRejectsWithTheSpecTexts trägt die beiden Zeilen, die
// `remove_transformation` durchläuft (`SPEC-019`): die Zeile zum Regelnamen
// und K4 — jeweils mit dem exakten Fehlertext (Klartext, Doppelpunkt,
// Adresse `schema.table.rule_name`) und dem Grund; ein Name, den nur die
// Nachbartabelle bzw. die Nachbarquelle führt, ist an dieser Tabelle nicht
// geführt. Rot färbende Mutation: die Namensprüfung `rule.Name() ==
// command.RuleName` durch `true` ersetzen (der nicht geführte Name wird
// angenommen), `state[command.Schema+"."+command.Table]` durch
// `state["public.other"]` ersetzen (der Nachbarname wird angenommen, der
// geführte abgelehnt), `CheckRuleName` streichen (leerer Name endet als
// `Regelname nicht geführt`).
func TestRemoveTransformationRejectsWithTheSpecTexts(t *testing.T) {
	for _, tc := range []struct {
		label    string
		source   model.SourceID
		ruleName string
		text     string
		reason   error
	}{
		{"Regelname leer", "src-1", "", "Regelname ist ungültig: public.orders.", domainerrors.ErrInvalidRuleName},
		{"Regelname mit Großbuchstaben", "src-1", "Kundenname", "Regelname ist ungültig: public.orders.Kundenname", domainerrors.ErrInvalidRuleName},
		{"K4 Name nicht geführt", "src-1", "gibt_es_nicht", "Regelname nicht geführt: public.orders.gibt_es_nicht", domainerrors.ErrRuleNotKept},
		{"K4 Name nur an der Nachbartabelle", "src-1", "statusname", "Regelname nicht geführt: public.orders.statusname", domainerrors.ErrRuleNotKept},
		{"K4 Name nur an der Nachbarquelle", "src-2", "kundenname", "Regelname nicht geführt: public.orders.kundenname", domainerrors.ErrRuleNotKept},
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

// TestRemoveTransformationChecksTheNameBeforeReadingTheStore trägt: ein
// Name außerhalb des Alphabets endet, ohne dass der Use Case den Regelstand
// liest. Rot färbende Mutation: `CheckRuleName` hinter das Lesen stellen.
func TestRemoveTransformationChecksTheNameBeforeReadingTheStore(t *testing.T) {
	port := newPort(t)
	if err := remove(port, "src-1", ""); err == nil {
		t.Fatal("Remove mit leerem Namen = nil, wollen einen Fehler")
	}
	if port.calls != 0 {
		t.Fatalf("Port-Aufrufe = %d, wollen 0", port.calls)
	}
}

// TestRemoveTransformationReturnsPortErrors trägt: ein Lesefehler des Ports
// endet unverändert.
func TestRemoveTransformationReturnsPortErrors(t *testing.T) {
	wantErr := stderrors.New("Antrags-Historie nicht lesbar")
	port := newPort(t)
	port.err = wantErr
	if err := remove(port, "src-1", "kundenname"); !stderrors.Is(err, wantErr) {
		t.Fatalf("Fehler = %v, wollen %v", err, wantErr)
	}
}
