package mapper_test

import (
	"context"
	stderrors "errors"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/decode"
	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/mapper"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Die Tests dieser Datei tragen die Wirkung von `map_value` am Ausgang des
// `Assembler` (`LH-FA-CFG-007`): die Relation der Regel-Tests (`feedRelation`)
// trägt `id`, `secret`, `name`; die Regel bildet Werte der Spalte `name` ab.

var mappedNames = map[string]string{"Ada": "Ada Lovelace", "Alan": "Alan Turing"}

// Die Regel wirkt auf beide Images aller drei Operationen: ein zugeordneter
// Wert steht abgebildet unter dem Schlüssel der Quellspalte, ein nicht
// zugeordneter bleibt, ein fehlendes Bild bleibt fehlend, und die
// ausgeschlossene Spalte erscheint weder mit Schlüssel noch mit Wert. Rot
// färbende Mutation: in `Assembler.change` `binding.Transformations` durch
// `nil` ersetzen — dann bleibt der Wert `Ada`.
func TestConsumeMapValueAppliesToBothImages(t *testing.T) {
	ctx := context.Background()
	assembler := newAssembler(t, ruleTables([]string{"secret"}, mapValueRule(t, "namen", "name", mappedNames)))
	relationEvent := feedRelation()

	insert := consumedChange(t, ctx, assembler, 1, decode.Change{
		Relation: relationEvent, Operation: decode.OpInsert,
		New: []*string{pointer("1"), pointer("geheim"), pointer("Ada")},
	})
	if string(insert.NewImage) != `{"id":"1","name":"Ada Lovelace"}` || insert.OldImage != nil {
		t.Fatalf("Insert: alt %s, neu %s", insert.OldImage, insert.NewImage)
	}

	update := consumedChange(t, ctx, assembler, 2, decode.Change{
		Relation: relationEvent, Operation: decode.OpUpdate,
		Old: []*string{pointer("1"), pointer("geheim"), pointer("Ada")},
		New: []*string{pointer("1"), pointer("geheimer"), pointer("Grace")},
	})
	if string(update.OldImage) != `{"id":"1","name":"Ada Lovelace"}` || string(update.NewImage) != `{"id":"1","name":"Grace"}` {
		t.Fatalf("Update-Images: alt %s, neu %s", update.OldImage, update.NewImage)
	}

	deletion := consumedChange(t, ctx, assembler, 3, decode.Change{
		Relation: relationEvent, Operation: decode.OpDelete,
		Old: []*string{pointer("2"), pointer("geheim"), pointer("Alan")},
	})
	if string(deletion.OldImage) != `{"id":"2","name":"Alan Turing"}` || deletion.NewImage != nil {
		t.Fatalf("Delete: alt %s, neu %s", deletion.OldImage, deletion.NewImage)
	}
	for _, image := range []string{string(insert.NewImage), string(update.OldImage), string(update.NewImage), string(deletion.OldImage)} {
		if strings.Contains(image, "geheim") || strings.Contains(image, `"secret"`) {
			t.Fatalf("Row Image %s trägt die ausgeschlossene Spalte", image)
		}
	}
}

// Ein fehlender Wert bleibt fehlend: NULL und unverändertes TOAST (nil-Wert)
// erzeugen weder einen Schlüssel noch einen abgebildeten Wert, ebenso ein
// Bild ohne die Spalte.
func TestConsumeMapValueKeepsAbsence(t *testing.T) {
	ctx := context.Background()
	assembler := newAssembler(t, ruleTables(nil, mapValueRule(t, "namen", "name", mappedNames)))
	relationEvent := feedRelation()

	nullValue := consumedChange(t, ctx, assembler, 1, decode.Change{
		Relation: relationEvent, Operation: decode.OpInsert,
		New: []*string{pointer("1"), pointer("s"), nil},
	})
	if string(nullValue.NewImage) != `{"id":"1","secret":"s"}` {
		t.Fatalf("Neu-Image mit NULL: %s", nullValue.NewImage)
	}
	short := consumedChange(t, ctx, assembler, 2, decode.Change{
		Relation: relationEvent, Operation: decode.OpInsert,
		New: []*string{pointer("1")},
	})
	if string(short.NewImage) != `{"id":"1"}` {
		t.Fatalf("Neu-Image ohne die Spalte: %s", short.NewImage)
	}
}

// Die Regel ändert kein Feld außer den Images.
func TestConsumeMapValueLeavesMetadataUnchanged(t *testing.T) {
	ctx := context.Background()
	event := decode.Change{
		Relation: feedRelation(), Operation: decode.OpUpdate,
		Old: []*string{pointer("1"), pointer("s"), pointer("Ada")},
		New: []*string{pointer("1"), pointer("s"), pointer("Alan")},
	}
	plain := consumedChange(t, ctx, newAssembler(t, ruleTables(nil)), 7, event)
	ruled := consumedChange(t, ctx, newAssembler(t, ruleTables(nil, mapValueRule(t, "namen", "name", mappedNames))), 7, event)

	if string(ruled.NewImage) == string(plain.NewImage) {
		t.Fatalf("Regel ohne Wirkung auf das Bild: %s", ruled.NewImage)
	}
	ruled.OldImage, ruled.NewImage = plain.OldImage, plain.NewImage
	if !reflect.DeepEqual(ruled, plain) {
		t.Fatalf("Change mit Regel %+v weicht ohne die Images von dem ohne Regel %+v ab", ruled, plain)
	}
}

// Anwendbarkeit hängt an der Spalte, nie am Wert einer Zeile (`SPEC-030`):
// eine Regel, deren Spalte in der Relation fehlt, ist nicht anwendbar (kein
// Change); jeder Wert einer vorhandenen Spalte — auch einer, der einer
// Spalte der Relation oder dem Zielwert einer anderen Zuordnung gleicht —
// lässt die Regel anwendbar. Rot färbende Mutationen: in
// `Transformation.CheckApplicable` die Prüfung der Spalte entfernen (Fall
// „Spalte fehlt“); die Prüfung des Zielnamens für `map_value` einschalten
// (`t.kind == TransformationRenameColumn` entfernen, die Relation enthält
// eine Spalte mit leerem Namen: Fall „leerer Spaltenname“).
func TestConsumeMapValueApplicabilityHangsOnTheColumnNotTheValue(t *testing.T) {
	assembler := newAssembler(t, ruleTables(nil, mapValueRule(t, "namen", "name", map[string]string{"Ada": "id", "id": "name", "": "secret"})))
	withoutName := relation("public", "feed", decode.Column{Name: "id", Key: true}, decode.Column{Name: "secret"})
	count, err := consumeChangeError(t, assembler, decode.Change{
		Relation: withoutName, Operation: decode.OpInsert, New: []*string{pointer("1"), pointer("s")},
	})
	if !stderrors.Is(err, mapper.ErrTransformationNotApplicable) || !stderrors.Is(err, domainerrors.ErrTransformationColumnMissing) || count != 0 {
		t.Fatalf("Spalte fehlt: Fehler = %v, Changes = %d, wollen ErrTransformationNotApplicable mit ErrTransformationColumnMissing und keinen Change", err, count)
	}
	if !strings.Contains(err.Error(), "namen") || !strings.Contains(err.Error(), "public.feed") {
		t.Fatalf("Fehler %q nennt weder Regelname noch Tabelle", err)
	}

	ctx := context.Background()
	for _, tc := range []struct {
		name     string
		relation *decode.Relation
		row      []*string
		want     string
	}{
		{"Wert gleicht einer Spalte der Relation", feedRelation(), []*string{pointer("1"), pointer("s"), pointer("Ada")}, `{"id":"1","secret":"s","name":"id"}`},
		{"Wert gleicht dem Namen der Quellspalte", feedRelation(), []*string{pointer("1"), pointer("s"), pointer("id")}, `{"id":"1","secret":"s","name":"name"}`},
		{"leerer Wert", feedRelation(), []*string{pointer("1"), pointer("s"), pointer("")}, `{"id":"1","secret":"s","name":"secret"}`},
		{"leerer Spaltenname", relation("public", "feed", decode.Column{Name: "id", Key: true}, decode.Column{Name: ""}, decode.Column{Name: "name"}), []*string{pointer("1"), pointer("s"), pointer("Ada")}, `{"id":"1","":"s","name":"id"}`},
	} {
		got := consumedChange(t, ctx, newAssembler(t, ruleTables(nil, mapValueRule(t, "namen", "name", map[string]string{"Ada": "id", "id": "name", "": "secret"}))), 1, decode.Change{
			Relation: tc.relation, Operation: decode.OpInsert, New: tc.row,
		})
		if string(got.NewImage) != tc.want {
			t.Fatalf("%s: Neu-Image %s, wollen %s", tc.name, got.NewImage, tc.want)
		}
	}
}

// Die Auswertung ist deterministisch: gleiche Regelmenge und Relation
// ergeben am Ausgang des Assemblers byte-gleiche Images, über getrennte
// Assembler-Instanzen, wiederholte Aufrufe, die Ordnung der Regeln und die
// Einfüge-Ordnung der Zuordnung. Rot färbende Mutation: in
// `applyTransformations` nur die erste Regel der Liste prüfen (`rules[:1]`) —
// die vertauschte Regelliste bildet den Namen nicht mehr ab.
func TestMapValueImagesAreDeterministic(t *testing.T) {
	ctx := context.Background()
	event := decode.Change{
		Relation: feedRelation(), Operation: decode.OpUpdate,
		Old: []*string{pointer("1"), pointer("s"), pointer("Ada")},
		New: []*string{pointer("1"), pointer("s"), pointer("Alan")},
	}
	build := func(names map[string]string) []model.Transformation {
		return []model.Transformation{
			mapValueRule(t, "a", "name", names),
			mapValueRule(t, "b", "id", map[string]string{"1": "eins"}),
		}
	}
	rules := build(mappedNames)
	first := consumedChange(t, ctx, newAssembler(t, ruleTables([]string{"secret"}, rules...)), 1, event)
	if string(first.OldImage) != `{"id":"eins","name":"Ada Lovelace"}` || string(first.NewImage) != `{"id":"eins","name":"Alan Turing"}` {
		t.Fatalf("Images: alt %s, neu %s", first.OldImage, first.NewImage)
	}
	permuted := []model.Transformation{rules[1], rules[0]}
	for i := uint32(2); i < 10; i++ {
		for _, again := range [][]model.Transformation{rules, permuted, build(map[string]string{"Alan": "Alan Turing", "Ada": "Ada Lovelace"})} {
			got := consumedChange(t, ctx, newAssembler(t, ruleTables([]string{"secret"}, again...)), 1, event)
			if string(got.OldImage) != string(first.OldImage) || string(got.NewImage) != string(first.NewImage) {
				t.Fatalf("Wiederholung %d: alt %s, neu %s, wollen alt %s, neu %s", i, got.OldImage, got.NewImage, first.OldImage, first.NewImage)
			}
		}
	}
}

// Der Regelstand einer laufenden Bindung folgt `SetTransformation` und
// `RemoveTransformation` auch für `map_value`: nach dem Setzen steht der
// abgebildete Wert, nach dem Entfernen wieder die Rohform.
func TestSetAndRemoveMapValueOnLiveBinding(t *testing.T) {
	ctx := context.Background()
	assembler := newAssembler(t, ruleTables(nil))
	image := func(xid uint32) string {
		change := consumedChange(t, ctx, assembler, xid, decode.Change{
			Relation: feedRelation(), Operation: decode.OpInsert, New: []*string{pointer("1"), pointer("s"), pointer("Ada")},
		})
		return string(change.NewImage)
	}
	if got := image(1); got != `{"id":"1","secret":"s","name":"Ada"}` {
		t.Fatalf("ohne Regel: %s", got)
	}
	assembler.SetTransformation("public.feed", mapValueRule(t, "namen", "name", mappedNames))
	if got := image(2); got != `{"id":"1","secret":"s","name":"Ada Lovelace"}` {
		t.Fatalf("nach dem Setzen: %s", got)
	}
	assembler.RemoveTransformation("public.feed", "namen")
	if got := image(3); got != `{"id":"1","secret":"s","name":"Ada"}` {
		t.Fatalf("nach dem Entfernen: %s", got)
	}
}

// Nebenläufigkeit (`-race`): der Regelstand gehört der Administration zu
// schreiben, während die Capture-Goroutine liest; jede erfasste Change trägt
// einen der beiden vollständigen Stände. Rot färbende Mutation: in
// `SetTransformation` das Sperren (`a.tablesMu.Lock()`) entfernen.
func TestAssemblerMapValueIsRaceFree(t *testing.T) {
	ctx := context.Background()
	assembler := newAssembler(t, testTables())
	relationEvent := relation("public", "feed", decode.Column{Name: "id", Key: true}, decode.Column{Name: "name"})
	rule := mapValueRule(t, "namen", "name", mappedNames)

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
				Relation: relationEvent, Operation: decode.OpInsert, New: []*string{pointer("1"), pointer("Ada")},
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
			if image := string(changes[0].NewImage); image != `{"id":"1","name":"Ada"}` && image != `{"id":"1","name":"Ada Lovelace"}` {
				t.Errorf("Neu-Image %s trägt keinen der beiden vollständigen Stände", image)
				return
			}
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			assembler.SetTransformation("public.feed", rule)
			assembler.RemoveTransformation("public.feed", "namen")
		}
	}()
	wg.Wait()
}
