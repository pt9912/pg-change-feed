package postgresstorage

import "testing"

// TestParseMetricValue belegt die zwei Lesarten eines `numeric`: eine ganze
// Zahl im Bereich von `int64` trägt `Int` exakt (auch jenseits von 2^53, wo
// ein `float64` rundet), ein Bruch oder eine Zahl außerhalb von `int64` nur
// `Float`, ein Nicht-Zahl-Text ist ein Fehler.
// Rot färbende Mutation: in `parseMetricValue` die Ganzzahl-Lesung streichen —
// `9007199254740993` kommt als `IsInt == false` und gerundeter Wert an.
func TestParseMetricValue(t *testing.T) {
	cases := []struct {
		text    string
		isInt   bool
		wantInt int64
		wantF   float64
		wantErr bool
	}{
		{text: "0", isInt: true, wantInt: 0, wantF: 0},
		{text: "42", isInt: true, wantInt: 42, wantF: 42},
		{text: "9007199254740993", isInt: true, wantInt: 9007199254740993, wantF: 9007199254740992},
		{text: "12.5", wantF: 12.5},
		{text: "5.0", wantF: 5},
		{text: "99999999999999999999", wantF: 1e20},
		{text: "abc", wantErr: true},
		{text: "", wantErr: true},
	}
	for _, tc := range cases {
		got, err := parseMetricValue(tc.text)
		if tc.wantErr {
			if err == nil {
				t.Errorf("%q: kein Fehler", tc.text)
			}
			continue
		}
		if err != nil {
			t.Errorf("%q: %v", tc.text, err)
			continue
		}
		if got.IsInt != tc.isInt || got.Int != tc.wantInt || got.Float != tc.wantF {
			t.Errorf("%q: %+v, wollen IsInt %v, Int %d, Float %v", tc.text, got, tc.isInt, tc.wantInt, tc.wantF)
		}
	}
}
