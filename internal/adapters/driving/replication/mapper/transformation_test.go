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

// Die Tests dieser Datei tragen den Regelstand einer Bindung
// (`LH-FA-CFG-007`, `ADR-0112`): Wirkung von `rename_column` auf beide Row
// Images, Anwendbarkeits-Prüfung vor jeder Serialisierung, Live-Reload des
// Regelstands und die Fitness Function der Entscheidung.

// renameRule legt eine `rename_column`-Regel an oder bricht den Test ab.
func renameRule(t testing.TB, name, column, to string) model.Transformation {
	t.Helper()
	rule, err := model.NewRenameColumn(name, column, to)
	if err != nil {
		t.Fatalf("NewRenameColumn(%q, %q, %q): %v", name, column, to, err)
	}
	return rule
}

func mapValueRule(t testing.TB, name, column string, values map[string]string) model.Transformation {
	t.Helper()
	rule, err := model.NewMapValue(name, column, values)
	if err != nil {
		t.Fatalf("NewMapValue(%q, %q, %v): %v", name, column, values, err)
	}
	return rule
}

// ruleTables trägt `public.feed` mit einem Ausschluss und einem Regelstand.
func ruleTables(excluded []string, rules ...model.Transformation) map[string]mapper.TableBinding {
	return map[string]mapper.TableBinding{
		"public.feed": {TableID: "tbl-1", SchemaVersion: "sv-1", ExcludedColumns: excluded, Transformations: rules},
	}
}

// feedRelation trägt die Relation der Regel-Tests: id, secret, name.
func feedRelation() *decode.Relation {
	return relation("public", "feed",
		decode.Column{Name: "id", Key: true},
		decode.Column{Name: "secret"},
		decode.Column{Name: "name"})
}

// Die Regel wirkt auf beide Images aller drei Operationen: der Schlüssel
// `name` steht als `customer_name` an der Position seiner Quellspalte, der
// Wert bleibt, ein fehlendes Bild bleibt fehlend, und der ausgeschlossene
// Schlüssel erscheint weder unter dem Quell- noch unter einem Zielnamen.
// Rot färbende Mutation: in `Assembler.change` `binding.Transformations`
// durch `nil` ersetzen — dann bleibt der Schlüssel `name`.
func TestConsumeRenameColumnAppliesToBothImages(t *testing.T) {
	ctx := context.Background()
	assembler := newAssembler(t, ruleTables([]string{"secret"}, renameRule(t, "kundenname", "name", "customer_name")))
	relationEvent := feedRelation()

	insert := consumedChange(t, ctx, assembler, 1, decode.Change{
		Relation: relationEvent, Operation: decode.OpInsert,
		New: []*string{pointer("1"), pointer("geheim"), pointer("Ada")},
	})
	if string(insert.NewImage) != `{"id":"1","customer_name":"Ada"}` {
		t.Fatalf("Insert-Neu-Image: %s", insert.NewImage)
	}
	if insert.OldImage != nil {
		t.Fatalf("Insert-Alt-Image: %s, wollen kein Bild", insert.OldImage)
	}

	update := consumedChange(t, ctx, assembler, 2, decode.Change{
		Relation: relationEvent, Operation: decode.OpUpdate,
		Old: []*string{pointer("1"), pointer("geheim"), pointer("alt")},
		New: []*string{pointer("1"), pointer("geheimer"), pointer("neu")},
	})
	if string(update.OldImage) != `{"id":"1","customer_name":"alt"}` || string(update.NewImage) != `{"id":"1","customer_name":"neu"}` {
		t.Fatalf("Update-Images: alt %s, neu %s", update.OldImage, update.NewImage)
	}

	deletion := consumedChange(t, ctx, assembler, 3, decode.Change{
		Relation: relationEvent, Operation: decode.OpDelete,
		Old: []*string{pointer("2"), pointer("geheim"), pointer("weg")},
	})
	if string(deletion.OldImage) != `{"id":"2","customer_name":"weg"}` {
		t.Fatalf("Delete-Alt-Image: %s", deletion.OldImage)
	}
	if deletion.NewImage != nil {
		t.Fatalf("Delete-Neu-Image: %s, wollen kein Bild", deletion.NewImage)
	}
	for _, image := range []string{string(insert.NewImage), string(update.OldImage), string(update.NewImage), string(deletion.OldImage)} {
		if strings.Contains(image, "geheim") || strings.Contains(image, `"secret"`) {
			t.Fatalf("Row Image %s trägt die ausgeschlossene Spalte", image)
		}
	}
}

// Ein fehlender Wert bleibt fehlend: NULL und unverändertes TOAST (nil-Wert)
// erzeugen keinen Zielschlüssel, ebenso ein Bild ohne die Spalte.
func TestConsumeRenameColumnKeepsAbsence(t *testing.T) {
	ctx := context.Background()
	assembler := newAssembler(t, ruleTables(nil, renameRule(t, "kundenname", "name", "customer_name")))
	relationEvent := feedRelation()

	nullValue := consumedChange(t, ctx, assembler, 1, decode.Change{
		Relation: relationEvent, Operation: decode.OpInsert,
		New: []*string{pointer("1"), pointer("s"), nil},
	})
	if string(nullValue.NewImage) != `{"id":"1","secret":"s"}` {
		t.Fatalf("Neu-Image mit NULL: %s, wollen weder name noch customer_name", nullValue.NewImage)
	}
	short := consumedChange(t, ctx, assembler, 2, decode.Change{
		Relation: relationEvent, Operation: decode.OpInsert,
		New: []*string{pointer("1")},
	})
	if string(short.NewImage) != `{"id":"1"}` {
		t.Fatalf("Neu-Image ohne die Spalte: %s", short.NewImage)
	}
}

// Die Regel ändert kein Feld außer den Images: derselbe Change mit und ohne
// Regelstand trägt gleiche Kennung, Transaktion, Tabelle, Sequenz, Operation,
// Schema-Version, Schema und Tabellenname.
func TestConsumeRenameColumnLeavesMetadataUnchanged(t *testing.T) {
	ctx := context.Background()
	event := decode.Change{
		Relation: feedRelation(), Operation: decode.OpUpdate,
		Old: []*string{pointer("1"), pointer("s"), pointer("alt")},
		New: []*string{pointer("1"), pointer("s"), pointer("neu")},
	}
	plain := consumedChange(t, ctx, newAssembler(t, ruleTables(nil)), 7, event)
	ruled := consumedChange(t, ctx, newAssembler(t, ruleTables(nil, renameRule(t, "kundenname", "name", "customer_name"))), 7, event)

	if string(ruled.NewImage) == string(plain.NewImage) {
		t.Fatalf("Regel ohne Wirkung auf das Bild: %s", ruled.NewImage)
	}
	ruled.OldImage, ruled.NewImage = plain.OldImage, plain.NewImage
	if !reflect.DeepEqual(ruled, plain) {
		t.Fatalf("Change mit Regel %+v weicht ohne die Images von dem ohne Regel %+v ab", ruled, plain)
	}
	if plain.ID == "" || plain.SourceTableID != "tbl-1" || plain.SchemaVersion != "sv-1" || plain.Schema != "public" || plain.Table != "feed" || plain.Sequence != 1 || plain.Operation != model.OperationUpdate {
		t.Fatalf("Vergleichs-Change trägt nicht die erwarteten Felder: %+v", plain)
	}
}

// consumeChangeError trägt Begin und Change und liefert den Fehler der
// Änderung samt dem Assembler-Zustand danach: der Commit der offenen
// Transaktion und die Zahl ihrer Changes.
func consumeChangeError(t *testing.T, assembler *mapper.Assembler, event decode.Change) (count int, changeErr error) {
	t.Helper()
	ctx := context.Background()
	if _, err := assembler.Consume(ctx, decode.Begin{XID: 1}); err != nil {
		t.Fatalf("Begin: %v", err)
	}
	_, changeErr = assembler.Consume(ctx, event)
	command, err := assembler.Consume(ctx, decode.Commit{CommitLSN: 100})
	if err != nil {
		t.Fatalf("Commit: %v", err)
	}
	changes, err := command.Transaction.Changes()
	if err != nil {
		t.Fatalf("Changes: %v", err)
	}
	return len(changes), changeErr
}

// Eine Regel, deren Spalte in der Relation fehlt, ist nicht anwendbar: der
// Assembler meldet `ErrTransformationNotApplicable` mit dem Grund, es
// entsteht kein Change, und die Sequenz rückt nicht vor. Rot färbende
// Mutation: in `Transformation.CheckApplicable` die Prüfung der Spalte
// entfernen — dann liefert der Assembler ein ungeändertes Bild ohne Fehler.
func TestConsumeRuleColumnMissingInRelationIsNotApplicable(t *testing.T) {
	assembler := newAssembler(t, ruleTables(nil, renameRule(t, "kundenname", "name", "customer_name")))
	withoutName := relation("public", "feed", decode.Column{Name: "id", Key: true}, decode.Column{Name: "secret"})

	count, err := consumeChangeError(t, assembler, decode.Change{
		Relation: withoutName, Operation: decode.OpInsert, New: []*string{pointer("1"), pointer("s")},
	})
	if !stderrors.Is(err, mapper.ErrTransformationNotApplicable) || !stderrors.Is(err, domainerrors.ErrTransformationColumnMissing) {
		t.Fatalf("Fehler = %v, wollen ErrTransformationNotApplicable mit ErrTransformationColumnMissing", err)
	}
	if !strings.Contains(err.Error(), "kundenname") || !strings.Contains(err.Error(), "public.feed") {
		t.Fatalf("Fehler %q nennt weder Regelname noch Tabelle", err)
	}
	if count != 0 {
		t.Fatalf("Change-Anzahl = %d, wollen 0 (kein Change ohne anwendbare Regel)", count)
	}
}

// Ein Zielname, der einer Spalte der Relation gleicht, macht die Regel
// nicht anwendbar — auch dann, wenn die Spalte ausgeschlossen ist (K3). Rot
// färbende Mutation: in `Transformation.CheckApplicable` die Prüfung des
// Zielnamens entfernen — dann trägt das Bild zwei gleichnamige Schlüssel.
func TestConsumeRuleTargetCollidesWithRelationColumnIsNotApplicable(t *testing.T) {
	cases := []struct {
		name     string
		excluded []string
		relation *decode.Relation
		values   []*string
	}{
		{
			name: "Zielname gleicht einer Spalte",
			relation: relation("public", "feed", decode.Column{Name: "id", Key: true},
				decode.Column{Name: "name"}, decode.Column{Name: "customer_name"}),
			values: []*string{pointer("1"), pointer("Ada"), pointer("Kunde")},
		},
		{
			name:     "Zielname gleicht einer ausgeschlossenen Spalte",
			excluded: []string{"customer_name"},
			relation: relation("public", "feed", decode.Column{Name: "id", Key: true},
				decode.Column{Name: "name"}, decode.Column{Name: "customer_name"}),
			values: []*string{pointer("1"), pointer("Ada"), pointer("Kunde")},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assembler := newAssembler(t, ruleTables(c.excluded, renameRule(t, "kundenname", "name", "customer_name")))
			count, err := consumeChangeError(t, assembler, decode.Change{Relation: c.relation, Operation: decode.OpInsert, New: c.values})
			if !stderrors.Is(err, mapper.ErrTransformationNotApplicable) || !stderrors.Is(err, domainerrors.ErrTransformationTargetCollides) {
				t.Fatalf("Fehler = %v, wollen ErrTransformationNotApplicable mit ErrTransformationTargetCollides", err)
			}
			if count != 0 {
				t.Fatalf("Change-Anzahl = %d, wollen 0", count)
			}
		})
	}
}

// Der Zielname wird zeichengenau verglichen: `Customer_Name` kollidiert
// nicht mit `customer_name`.
func TestConsumeRuleTargetComparisonIsCaseSensitive(t *testing.T) {
	ctx := context.Background()
	assembler := newAssembler(t, ruleTables(nil, renameRule(t, "kundenname", "name", "Customer_Name")))
	relationEvent := relation("public", "feed", decode.Column{Name: "id", Key: true},
		decode.Column{Name: "name"}, decode.Column{Name: "customer_name"})
	change := consumedChange(t, ctx, assembler, 1, decode.Change{
		Relation: relationEvent, Operation: decode.OpInsert, New: []*string{pointer("1"), pointer("Ada"), pointer("Kunde")},
	})
	if string(change.NewImage) != `{"id":"1","Customer_Name":"Ada","customer_name":"Kunde"}` {
		t.Fatalf("Neu-Image: %s", change.NewImage)
	}
}

// Der Fall, den erst diese Prüfung fängt (`ADR-0112` Teilfrage 4): eine
// kompatible Spalten-Erweiterung um den Zielnamen passiert `observeRelation`
// ohne Fehler, die nächste Änderung der Tabelle endet an der Regel. Die
// spalten-entfernenden Fälle enden vorher an `relationOther`, mit
// `ErrIncompatibleSchemaChange` und ohne `ErrTransformationNotApplicable`.
func TestConsumeRuleTargetCollisionAfterCompatibleExtension(t *testing.T) {
	ctx := context.Background()
	store := newFakeSchemaStore()
	store.versions["tbl-1"] = model.SchemaVersion{ID: "sv-1", SourceTableID: "tbl-1", Version: 1}
	store.schemas["sv-1"] = model.TableSchema{VersionID: "sv-1", Columns: []model.Column{{Name: "id", OID: 23}, {Name: "name", OID: 25}}}
	tables := map[string]mapper.TableBinding{
		"public.feed": {TableID: "tbl-1", SchemaVersion: "sv-1", Transformations: []model.Transformation{renameRule(t, "kundenname", "name", "customer_name")}},
	}
	assembler, err := mapper.NewAssembler("src-1", tables, store)
	if err != nil {
		t.Fatalf("NewAssembler: %v", err)
	}

	original := relation("public", "feed", decode.Column{Name: "id", Key: true, TypeOID: 23}, decode.Column{Name: "name", TypeOID: 25})
	before := consumedChange(t, ctx, assembler, 1, decode.Change{
		Relation: original, Operation: decode.OpInsert, New: []*string{pointer("1"), pointer("Ada")},
	})
	if string(before.NewImage) != `{"id":"1","customer_name":"Ada"}` {
		t.Fatalf("Neu-Image vor der Erweiterung: %s", before.NewImage)
	}

	extended := relation("public", "feed", decode.Column{Name: "id", Key: true, TypeOID: 23},
		decode.Column{Name: "name", TypeOID: 25}, decode.Column{Name: "customer_name", TypeOID: 25})
	if _, err := assembler.Consume(ctx, extended); err != nil {
		t.Fatalf("kompatible Erweiterung um den Zielnamen: %v, wollen keinen Fehler in observeRelation", err)
	}
	count, changeErr := consumeChangeError(t, assembler, decode.Change{
		Relation: extended, Operation: decode.OpInsert, New: []*string{pointer("2"), pointer("Bob"), pointer("Kunde")},
	})
	if !stderrors.Is(changeErr, mapper.ErrTransformationNotApplicable) || !stderrors.Is(changeErr, domainerrors.ErrTransformationTargetCollides) {
		t.Fatalf("Fehler = %v, wollen ErrTransformationNotApplicable mit ErrTransformationTargetCollides", changeErr)
	}
	if count != 0 {
		t.Fatalf("Change-Anzahl = %d, wollen 0", count)
	}

	removed := relation("public", "feed", decode.Column{Name: "id", Key: true, TypeOID: 23})
	relationErr := func() error { _, err := assembler.Consume(ctx, removed); return err }()
	if !stderrors.Is(relationErr, mapper.ErrIncompatibleSchemaChange) || stderrors.Is(relationErr, mapper.ErrTransformationNotApplicable) {
		t.Fatalf("spalten-entfernende Relation: %v, wollen ErrIncompatibleSchemaChange ohne ErrTransformationNotApplicable", relationErr)
	}
}

// Nach dem Fehler der Nichtanwendbarkeit trägt die Bindung ihren Stand: die
// Abhilfe `RemoveTransformation` setzt die Erfassung fort, die erste Change
// danach trägt die Sequenz 1 (die fehlgeschlagene Änderung hat sie nicht
// verbraucht).
func TestConsumeRuleRemedyRestoresCapture(t *testing.T) {
	ctx := context.Background()
	assembler := newAssembler(t, ruleTables(nil, renameRule(t, "kundenname", "name", "customer_name")))
	withoutName := relation("public", "feed", decode.Column{Name: "id", Key: true})

	if _, err := assembler.Consume(ctx, decode.Begin{XID: 1}); err != nil {
		t.Fatalf("Begin: %v", err)
	}
	if _, err := assembler.Consume(ctx, decode.Change{Relation: withoutName, Operation: decode.OpInsert, New: []*string{pointer("1")}}); !stderrors.Is(err, mapper.ErrTransformationNotApplicable) {
		t.Fatalf("Fehler = %v, wollen ErrTransformationNotApplicable", err)
	}
	assembler.RemoveTransformation("public.feed", "kundenname")
	if _, err := assembler.Consume(ctx, decode.Change{Relation: withoutName, Operation: decode.OpInsert, New: []*string{pointer("1")}}); err != nil {
		t.Fatalf("Change nach der Abhilfe: %v", err)
	}
	command, err := assembler.Consume(ctx, decode.Commit{CommitLSN: 100})
	if err != nil {
		t.Fatalf("Commit: %v", err)
	}
	changes, err := command.Transaction.Changes()
	if err != nil || len(changes) != 1 || changes[0].Sequence != 1 {
		t.Fatalf("Changes nach der Abhilfe: %+v (%v), wollen eine Change mit Sequenz 1", changes, err)
	}
}

// Eine Regel an einer Bindung wirkt nur an ihrer Tabelle, und eine nicht
// aktivierte Tabelle bleibt ohne Erfassung, auch ohne anwendbare Regel.
func TestConsumeRuleActsOnlyOnItsTable(t *testing.T) {
	ctx := context.Background()
	tables := ruleTables(nil, renameRule(t, "kundenname", "name", "customer_name"))
	tables["public.other"] = mapper.TableBinding{TableID: "tbl-2", SchemaVersion: "sv-2"}
	assembler := newAssembler(t, tables)
	other := relation("public", "other", decode.Column{Name: "id", Key: true}, decode.Column{Name: "name"})

	change := consumedChange(t, ctx, assembler, 1, decode.Change{
		Relation: other, Operation: decode.OpInsert, New: []*string{pointer("1"), pointer("Ada")},
	})
	if string(change.NewImage) != `{"id":"1","name":"Ada"}` {
		t.Fatalf("Neu-Image der anderen Tabelle: %s, wollen unverändert", change.NewImage)
	}

	unbound := relation("public", "unbound", decode.Column{Name: "id", Key: true})
	if _, err := assembler.Consume(ctx, decode.Begin{XID: 2}); err != nil {
		t.Fatalf("Begin: %v", err)
	}
	if _, err := assembler.Consume(ctx, decode.Change{Relation: unbound, Operation: decode.OpInsert, New: []*string{pointer("1")}}); err != nil {
		t.Fatalf("Change an nicht aktivierter Tabelle: %v", err)
	}
}

// Der Regelstand überlebt die Schema-Version-Hebung einer kompatiblen
// Erweiterung (`setSchemaVersion`): die Regel wirkt danach weiter.
func TestConsumeRuleSurvivesSchemaBump(t *testing.T) {
	ctx := context.Background()
	store := newFakeSchemaStore()
	store.versions["tbl-1"] = model.SchemaVersion{ID: "sv-1", SourceTableID: "tbl-1", Version: 1}
	store.schemas["sv-1"] = model.TableSchema{VersionID: "sv-1", Columns: []model.Column{{Name: "id", OID: 23}, {Name: "name", OID: 25}}}
	tables := map[string]mapper.TableBinding{
		"public.feed": {TableID: "tbl-1", SchemaVersion: "sv-1", Transformations: []model.Transformation{renameRule(t, "kundenname", "name", "customer_name")}},
	}
	assembler, err := mapper.NewAssembler("src-1", tables, store)
	if err != nil {
		t.Fatalf("NewAssembler: %v", err)
	}
	extended := relation("public", "feed", decode.Column{Name: "id", Key: true, TypeOID: 23},
		decode.Column{Name: "name", TypeOID: 25}, decode.Column{Name: "extra", TypeOID: 25})
	if _, err := assembler.Consume(ctx, extended); err != nil {
		t.Fatalf("Relation: %v", err)
	}
	if store.versions["tbl-1"].Version != 2 {
		t.Fatalf("Schema-Version: %+v, wollen Version 2", store.versions["tbl-1"])
	}
	change := consumedChange(t, ctx, assembler, 1, decode.Change{
		Relation: extended, Operation: decode.OpInsert, New: []*string{pointer("1"), pointer("Ada"), pointer("z")},
	})
	if change.SchemaVersion != store.versions["tbl-1"].ID {
		t.Fatalf("Schema-Version des Changes: %s, wollen %s", change.SchemaVersion, store.versions["tbl-1"].ID)
	}
	if string(change.NewImage) != `{"id":"1","customer_name":"Ada","extra":"z"}` {
		t.Fatalf("Neu-Image nach dem Schema-Bump: %s, wollen die Regel weiter wirksam", change.NewImage)
	}
}

// `AddBinding` bei getragener Bindung erhält deren Regelstand — auch gegen
// eine übergebene Bindung, die einen anderen trägt (dieselbe Regel wie beim
// Ausschlussstand). Bei einer neuen Tabelle geht der übergebene Stand ein.
// Rot färbende Mutation: in `AddBinding` die Zeile
// `binding.Transformations = existing.Transformations` entfernen.
func TestAddBindingKeepsRuleState(t *testing.T) {
	ctx := context.Background()
	assembler := newAssembler(t, ruleTables(nil, renameRule(t, "kundenname", "name", "customer_name")))

	assembler.AddBinding("public.feed", mapper.TableBinding{TableID: "tbl-2", SchemaVersion: "sv-2",
		Transformations: []model.Transformation{renameRule(t, "fremd", "id", "pk")}})

	change := consumedChange(t, ctx, assembler, 1, decode.Change{
		Relation: feedRelation(), Operation: decode.OpInsert, New: []*string{pointer("1"), pointer("s"), pointer("Ada")},
	})
	if change.SourceTableID != "tbl-2" || change.SchemaVersion != "sv-2" {
		t.Fatalf("Change trägt nicht die nachgetragene Bindung: %+v", change)
	}
	if string(change.NewImage) != `{"id":"1","secret":"s","customer_name":"Ada"}` {
		t.Fatalf("Neu-Image nach dem erneuten AddBinding: %s, wollen den erhaltenen Regelstand", change.NewImage)
	}

	assembler.AddBinding("public.new", mapper.TableBinding{TableID: "tbl-3", SchemaVersion: "sv-3",
		Transformations: []model.Transformation{renameRule(t, "kundenname", "name", "customer_name")}})
	created := consumedChange(t, ctx, assembler, 2, decode.Change{
		Relation:  relation("public", "new", decode.Column{Name: "id", Key: true}, decode.Column{Name: "name"}),
		Operation: decode.OpInsert, New: []*string{pointer("1"), pointer("Ada")},
	})
	if string(created.NewImage) != `{"id":"1","customer_name":"Ada"}` {
		t.Fatalf("Neu-Image einer neuen Bindung: %s, wollen den übergebenen Regelstand", created.NewImage)
	}
}

// `SetTransformation`/`RemoveTransformation` wirken ab dem Aufruf auf die
// laufende Bindung: setzen, ersetzen unter demselben Namen, gezielt
// entfernen, unbekannten Namen und nicht getragene Bindung ohne Wirkung.
func TestSetAndRemoveTransformationOnLiveBinding(t *testing.T) {
	ctx := context.Background()
	assembler := newAssembler(t, testTables())
	relationEvent := feedRelation()
	newImage := func(xid uint32) string {
		change := consumedChange(t, ctx, assembler, xid, decode.Change{
			Relation: relationEvent, Operation: decode.OpInsert, New: []*string{pointer("1"), pointer("s"), pointer("Ada")},
		})
		return string(change.NewImage)
	}

	if got := newImage(1); got != `{"id":"1","secret":"s","name":"Ada"}` {
		t.Fatalf("ohne Regel: %s", got)
	}
	assembler.SetTransformation("public.feed", renameRule(t, "kundenname", "name", "customer_name"))
	if got := newImage(2); got != `{"id":"1","secret":"s","customer_name":"Ada"}` {
		t.Fatalf("nach SetTransformation: %s", got)
	}
	assembler.SetTransformation("public.feed", renameRule(t, "kundenname", "name", "kunde"))
	if got := newImage(3); got != `{"id":"1","secret":"s","kunde":"Ada"}` {
		t.Fatalf("nach dem Ersetzen unter demselben Namen: %s", got)
	}
	assembler.SetTransformation("public.feed", renameRule(t, "schluessel", "id", "pk"))
	if got := newImage(4); got != `{"pk":"1","secret":"s","kunde":"Ada"}` {
		t.Fatalf("mit zwei Regeln: %s", got)
	}
	assembler.RemoveTransformation("public.feed", "unbekannt")
	if got := newImage(5); got != `{"pk":"1","secret":"s","kunde":"Ada"}` {
		t.Fatalf("nach Entfernen eines nicht geführten Namens: %s", got)
	}
	assembler.RemoveTransformation("public.feed", "kundenname")
	if got := newImage(6); got != `{"pk":"1","secret":"s","name":"Ada"}` {
		t.Fatalf("nach gezieltem Entfernen: %s, wollen nur die zweite Regel", got)
	}

	// Eine nicht getragene Bindung bleibt ohne Wirkung und wird nicht
	// belebt: `SetTransformation` hinterlässt keine Regel, die eine spätere
	// Bindung erbt, `RemoveTransformation` aktiviert die Tabelle nicht.
	assembler.SetTransformation("public.unbound", renameRule(t, "r", "id", "pk"))
	assembler.AddBinding("public.unbound", mapper.TableBinding{TableID: "tbl-9", SchemaVersion: "sv-9"})
	unbound := consumedChange(t, ctx, assembler, 7, decode.Change{
		Relation: relation("public", "unbound", decode.Column{Name: "id", Key: true}), Operation: decode.OpInsert, New: []*string{pointer("1")},
	})
	if string(unbound.NewImage) != `{"id":"1"}` {
		t.Fatalf("Regel gegen eine nicht getragene Bindung: %s, wollen keine Wirkung", unbound.NewImage)
	}
	assembler.RemoveTransformation("public.never", "r")
	count, err := consumeChangeError(t, assembler, decode.Change{
		Relation: relation("public", "never", decode.Column{Name: "id", Key: true}), Operation: decode.OpInsert, New: []*string{pointer("1")},
	})
	if err != nil || count != 0 {
		t.Fatalf("Change an einer Tabelle ohne Bindung: %d Changes, Fehler %v, wollen keine Erfassung", count, err)
	}
}

// Die Fitness Function von `ADR-0112` (Eigenschaftstest): für jeden Regeltyp der Domänen-Menge und jede Spalte der Relation als
// ausgeschlossene Spalte trägt das Bild weder den Quellschlüssel noch einen
// Zielnamen noch den Quellwert noch den abgebildeten Wert — die Regel sitzt
// dabei genau auf der ausgeschlossenen Spalte, eine zweite Regel auf einer
// anderen Spalte wirkt weiter. Rot färbende Mutationen: in `BuildRowImage` die
// Prüfung `containsName(excluded, column)` entfernen bzw. hinter die Auswertung
// der Regeln setzen; den Fall eines Regeltyps aus `ruleFor` streichen (der Test
// bricht mit „nicht abgedeckt“ ab).
func TestExcludedColumnIsUnreachableForEveryRuleKind(t *testing.T) {
	ctx := context.Background()
	columnNames := []string{"id", "secret", "name", "status"}
	columns := make([]decode.Column, len(columnNames))
	for i, name := range columnNames {
		columns[i] = decode.Column{Name: name, Key: i == 0}
	}
	relationEvent := relation("public", "feed", columns...)
	sentinel := func(column string) string { return "SENTINEL-" + column }

	// ruled trägt eine Regel und den Schlüssel und den Wert, unter denen die
	// Spalte mit dem Sentinel-Wert im Bild steht.
	type ruled struct {
		rule       model.Transformation
		key, value string
	}
	// ruleFor liefert je Regeltyp eine Regel auf der Spalte; ein Regeltyp,
	// den diese Tabelle nicht kennt, bricht den Test ab — ein neuer Typ
	// braucht seinen Fall hier.
	ruleFor := func(kind model.TransformationKind, column string) ruled {
		switch kind {
		case model.TransformationRenameColumn:
			return ruled{renameRule(t, "regel_"+column, column, "ziel_"+column), "ziel_" + column, sentinel(column)}
		case model.TransformationMapValue:
			mapped := "ABGEBILDET-" + column
			return ruled{mapValueRule(t, "regel_"+column, column, map[string]string{sentinel(column): mapped}), column, mapped}
		default:
			t.Fatalf("Regeltyp %q ist in diesem Test nicht abgedeckt", kind)
			return ruled{}
		}
	}

	xid := uint32(0)
	for _, kind := range model.TransformationKinds() {
		for _, excludedColumn := range columnNames {
			otherColumn := "id"
			if excludedColumn == "id" {
				otherColumn = "status"
			}
			excludedRule, otherRule := ruleFor(kind, excludedColumn), ruleFor(kind, otherColumn)
			rules := []model.Transformation{excludedRule.rule, otherRule.rule}
			assembler := newAssembler(t, ruleTables([]string{excludedColumn}, rules...))
			values := make([]*string, len(columnNames))
			for i, name := range columnNames {
				values[i] = pointer(sentinel(name))
			}
			xid++
			change := consumedChange(t, ctx, assembler, xid, decode.Change{
				Relation: relationEvent, Operation: decode.OpUpdate, Old: values, New: values,
			})
			for label, image := range map[string][]byte{"Alt-Image": change.OldImage, "Neu-Image": change.NewImage} {
				text := string(image)
				for _, forbidden := range []string{
					`"` + excludedColumn + `"`,
					`"` + excludedRule.key + `"`,
					sentinel(excludedColumn),
					excludedRule.value,
				} {
					if strings.Contains(text, forbidden) {
						t.Fatalf("Regeltyp %s, ausgeschlossen %s: %s %s trägt %s", kind, excludedColumn, label, text, forbidden)
					}
				}
				if !strings.Contains(text, `"`+otherRule.key+`":"`+otherRule.value+`"`) {
					t.Fatalf("Regeltyp %s, ausgeschlossen %s: %s %s trägt die Regel auf %s nicht", kind, excludedColumn, label, text, otherColumn)
				}
			}
		}
	}
}

// Jede Regel des Regelstands wird gegen die Relation geprüft, an jeder
// Position: die nicht anwendbare Regel steht an erster, an mittlerer und an
// letzter Stelle, und der Fehler nennt genau sie. Rot färbende Mutationen:
// `checkTransformations` prüft nur `rules[:1]` (Fälle 1, 3 und 4), nur
// `rules[1:]` (Fall 2) bzw. `rules[:len(rules)-1]` (Fälle 1 und 3).
func TestConsumeEveryRuleIsCheckedForApplicability(t *testing.T) {
	cases := []struct {
		name     string
		rules    []model.Transformation
		wantErr  error
		wantRule string
	}{
		{
			name: "zweite Regel nicht anwendbar, erste anwendbar",
			rules: []model.Transformation{
				renameRule(t, "erste", "name", "customer_name"),
				renameRule(t, "zweite", "unbekannt", "x"),
			},
			wantErr:  domainerrors.ErrTransformationColumnMissing,
			wantRule: `"zweite"`,
		},
		{
			name: "erste Regel nicht anwendbar, zweite anwendbar",
			rules: []model.Transformation{
				renameRule(t, "erste", "unbekannt", "x"),
				renameRule(t, "zweite", "name", "customer_name"),
			},
			wantErr:  domainerrors.ErrTransformationColumnMissing,
			wantRule: `"erste"`,
		},
		{
			name: "letzte von drei Regeln nicht anwendbar",
			rules: []model.Transformation{
				renameRule(t, "erste", "name", "customer_name"),
				renameRule(t, "zweite", "id", "pk"),
				renameRule(t, "dritte", "secret", "id"),
			},
			wantErr:  domainerrors.ErrTransformationTargetCollides,
			wantRule: `"dritte"`,
		},
		{
			name: "mittlere von drei Regeln nicht anwendbar",
			rules: []model.Transformation{
				renameRule(t, "erste", "id", "pk"),
				renameRule(t, "zweite", "name", "secret"),
				renameRule(t, "dritte", "secret", "hidden"),
			},
			wantErr:  domainerrors.ErrTransformationTargetCollides,
			wantRule: `"zweite"`,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assembler := newAssembler(t, ruleTables(nil, c.rules...))
			count, err := consumeChangeError(t, assembler, decode.Change{
				Relation: feedRelation(), Operation: decode.OpInsert,
				New: []*string{pointer("1"), pointer("s"), pointer("Ada")},
			})
			if !stderrors.Is(err, mapper.ErrTransformationNotApplicable) || !stderrors.Is(err, c.wantErr) {
				t.Fatalf("Fehler = %v, wollen ErrTransformationNotApplicable mit %v", err, c.wantErr)
			}
			if !strings.Contains(err.Error(), c.wantRule) {
				t.Fatalf("Fehler %q nennt nicht die Regel %s", err, c.wantRule)
			}
			if count != 0 {
				t.Fatalf("Change-Anzahl = %d, wollen 0", count)
			}
		})
	}
}

// Zwei Regeln auf verschiedene Spalten mit gleichem Zielnamen sind je einzeln
// anwendbar; treffen beide im selben Bild, endet die Änderung als nicht
// anwendbare Regel statt mit zwei gleichnamigen Schlüsseln — an beiden Images
// (Insert: Neu-Image, Delete: Alt-Image). Trägt eine der beiden Spalten
// keinen Wert, steht nur ein Schlüssel im Bild und die Änderung geht durch.
// Rot färbende Mutationen: in `Assembler.change` den Fehler eines der beiden
// `BuildRowImage`-Aufrufe ohne `imageError` zurückgeben (der Fehler trägt dann
// `ErrTransformationNotApplicable` nicht); in `BuildRowImage` die Prüfung
// gegen die umbenannten Schlüssel entfernen.
func TestConsumeTwoRulesWithSameTargetAreNotApplicableWhenBothMatch(t *testing.T) {
	ctx := context.Background()
	rules := []model.Transformation{renameRule(t, "a", "name", "z"), renameRule(t, "b", "secret", "z")}
	for _, c := range []struct {
		name  string
		event decode.Change
	}{
		{"Neu-Image", decode.Change{Relation: feedRelation(), Operation: decode.OpInsert,
			New: []*string{pointer("1"), pointer("s"), pointer("Ada")}}},
		{"Alt-Image", decode.Change{Relation: feedRelation(), Operation: decode.OpDelete,
			Old: []*string{pointer("1"), pointer("s"), pointer("Ada")}}},
	} {
		t.Run(c.name, func(t *testing.T) {
			count, err := consumeChangeError(t, newAssembler(t, ruleTables(nil, rules...)), c.event)
			if !stderrors.Is(err, mapper.ErrTransformationNotApplicable) || !stderrors.Is(err, domainerrors.ErrTransformationTargetCollides) {
				t.Fatalf("Fehler = %v, wollen ErrTransformationNotApplicable mit ErrTransformationTargetCollides", err)
			}
			if count != 0 {
				t.Fatalf("Change-Anzahl = %d, wollen 0", count)
			}
		})
	}

	single := consumedChange(t, ctx, newAssembler(t, ruleTables(nil, rules...)), 1, decode.Change{
		Relation: feedRelation(), Operation: decode.OpInsert,
		New: []*string{pointer("1"), nil, pointer("Ada")},
	})
	if string(single.NewImage) != `{"id":"1","z":"Ada"}` {
		t.Fatalf("Neu-Image mit einer wertlosen Quellspalte: %s", single.NewImage)
	}
}

// Die Auswertung ist deterministisch: gleiche Regelmenge und Relation
// ergeben am Ausgang des Assemblers byte-gleiche Images, über getrennte
// Assembler-Instanzen und wiederholte Aufrufe.
func TestRenameColumnImagesAreDeterministic(t *testing.T) {
	ctx := context.Background()
	rules := []model.Transformation{renameRule(t, "a", "name", "customer_name"), renameRule(t, "b", "id", "pk")}
	event := decode.Change{
		Relation: feedRelation(), Operation: decode.OpUpdate,
		Old: []*string{pointer("1"), pointer("s"), pointer("alt")},
		New: []*string{pointer("1"), pointer("s"), pointer("neu")},
	}
	first := consumedChange(t, ctx, newAssembler(t, ruleTables([]string{"secret"}, rules...)), 1, event)
	for i := uint32(2); i < 10; i++ {
		again := consumedChange(t, ctx, newAssembler(t, ruleTables([]string{"secret"}, rules...)), 1, event)
		if string(again.OldImage) != string(first.OldImage) || string(again.NewImage) != string(first.NewImage) {
			t.Fatalf("Wiederholung %d: alt %s, neu %s, wollen alt %s, neu %s", i, again.OldImage, again.NewImage, first.OldImage, first.NewImage)
		}
	}
	if string(first.NewImage) != `{"pk":"1","customer_name":"neu"}` {
		t.Fatalf("Neu-Image: %s", first.NewImage)
	}
}

// Die Fitness Function der Nebenläufigkeit (`ADR-0112`): der Regelstand
// gehört einer zweiten Goroutine (Administration) zu schreiben, während die
// Capture-Goroutine liest. Jede erfasste Change trägt einen der beiden
// vollständigen Stände, nie einen gemischten; `go test -race` deckt den
// gleichzeitigen Zugriff. Rot färbende Mutation: in `SetTransformation` das
// Sperren (`a.tablesMu.Lock()`) entfernen.
func TestAssemblerTransformationsAreRaceFree(t *testing.T) {
	ctx := context.Background()
	assembler := newAssembler(t, testTables())
	relationEvent := relation("public", "feed", decode.Column{Name: "id", Key: true}, decode.Column{Name: "name"})
	rule := renameRule(t, "kundenname", "name", "customer_name")

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
			if image := string(changes[0].NewImage); image != `{"id":"1","name":"Ada"}` && image != `{"id":"1","customer_name":"Ada"}` {
				t.Errorf("Neu-Image %s trägt keinen der beiden vollständigen Stände", image)
				return
			}
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			assembler.SetTransformation("public.feed", rule)
			assembler.RemoveTransformation("public.feed", "kundenname")
		}
	}()
	wg.Wait()
}
