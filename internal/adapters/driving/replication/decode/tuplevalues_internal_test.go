package decode

import "testing"

// Whitebox-Test (`package decode`, nicht `decode_test`): die beiden
// Tupel-Übersetzungen sind unexportierte Übersetzungsdetails; die
// Abwesenheits-Grenze (`LH-FA-CAP-008` Boundary) gehört zu ihrem
// Vertrag, der Dekodier-Weg über `Decode` reicht sie aber nicht her —
// die Relation-Nachricht muß vorausgehen, und die Tupel des Treibers
// sind nach dem Parsen nie nil. Der Test greift deshalb direkt zu,
// statt die Funktion für den Test zu exportieren.
//
// Rot färbende Mutation: in `tupleValues`/`oldTupleValues` den Zug
// `if tuple == nil { return nil, nil }` auf `return []*string{}, nil`
// umstellen — dann tragen beide Übersetzungen ein leeres, gesetztes
// Ergebnis statt der Abwesenheit, und dieser Test fällt.
func TestTupleTranslationsCarryMissingTupleAsAbsence(t *testing.T) {
	relation := &Relation{
		Schema:  "public",
		Name:    "feed",
		Columns: []Column{{Name: "id", Key: true}},
	}

	values, err := tupleValues(relation, nil)
	if err != nil || values != nil {
		t.Fatalf("tupleValues(nil) = %v, %v, wollen nil, nil (Abwesenheit)", values, err)
	}

	old, err := oldTupleValues(relation, 0, nil)
	if err != nil || old != nil {
		t.Fatalf("oldTupleValues(nil) = %v, %v, wollen nil, nil (Abwesenheit)", old, err)
	}
}
