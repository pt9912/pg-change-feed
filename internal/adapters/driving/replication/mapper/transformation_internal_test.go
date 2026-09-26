package mapper

import (
	"testing"

	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Whitebox-Test (`package mapper`): der Vertrag „ein Leser hält seinen
// Schnappschuss ohne eigene Sperre“ hängt daran, dass `SetTransformation` und
// `RemoveTransformation` die Regelliste neu aufbauen und den Speicher der
// gelesenen Liste nicht anfassen — über `lookupBinding` ist dieser
// Schnappschuss greifbar, über `Consume` nicht. Jede Operation läuft auf einem
// frischen Assembler mit einer Liste, deren Kapazität über ihre Länge
// hinausgeht.
//
// Rot färbende Mutationen: in `withTransformation` die Ersetzung an Ort und
// Stelle (`rules[i] = rule`) bzw. das Anhängen an die übergebene Liste
// (`append(rules, rule)` bei freier Kapazität), in `withoutTransformation`
// das Entfernen durch Verschieben in der übergebenen Liste — dann ändert
// sich der Schnappschuss des jeweiligen Falls.
func TestTransformationListsAreReplacedNotMutated(t *testing.T) {
	rule := func(name, column, to string) model.Transformation {
		t.Helper()
		created, err := model.NewRenameColumn(name, column, to)
		if err != nil {
			t.Fatalf("NewRenameColumn: %v", err)
		}
		return created
	}
	a, b, c := rule("a", "col_a", "to_a"), rule("b", "col_b", "to_b"), rule("c", "col_c", "to_c")

	cases := []struct {
		name  string
		apply func(*Assembler)
		want  []string
	}{
		{"Ersetzen unter demselben Namen", func(asm *Assembler) { asm.SetTransformation("public.feed", rule("b", "col_b", "to_b2")) }, []string{"a:to_a", "b:to_b2", "c:to_c"}},
		{"Anhängen", func(asm *Assembler) { asm.SetTransformation("public.feed", rule("d", "col_d", "to_d")) }, []string{"a:to_a", "b:to_b", "c:to_c", "d:to_d"}},
		{"Entfernen", func(asm *Assembler) { asm.RemoveTransformation("public.feed", "a") }, []string{"b:to_b", "c:to_c"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			initial := make([]model.Transformation, 0, 8)
			initial = append(initial, a, b, c)
			assembler, err := NewAssembler("src-1", map[string]TableBinding{
				"public.feed": {TableID: "tbl-1", SchemaVersion: "sv-1", Transformations: initial},
			}, nil)
			if err != nil {
				t.Fatalf("NewAssembler: %v", err)
			}
			binding, _ := assembler.lookupBinding("public.feed")
			snapshot := binding.Transformations

			tc.apply(assembler)

			if len(snapshot) != 3 || snapshot[0] != a || snapshot[1] != b || snapshot[2] != c {
				t.Fatalf("Schnappschuss %+v verändert, wollen [a b c]", snapshot)
			}
			if beyond := snapshot[:4][3]; beyond != (model.Transformation{}) {
				t.Fatalf("der Speicher hinter dem Schnappschuss trägt %+v", beyond)
			}
			live, _ := assembler.lookupBinding("public.feed")
			got := make([]string, len(live.Transformations))
			for i, ruleValue := range live.Transformations {
				got[i] = ruleValue.Name() + ":" + ruleValue.To()
			}
			if len(got) != len(tc.want) {
				t.Fatalf("Regelstand der Bindung = %v, wollen %v", got, tc.want)
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Fatalf("Regelstand der Bindung = %v, wollen %v", got, tc.want)
				}
			}
		})
	}
}
