package apiauth_test

import (
	"testing"

	"github.com/pt9912/pg-change-feed/internal/application/port/apiauth"
)

// TestClassify trägt die Zuordnung `Token → Rolle` je Fall mit Meldungstext.
//
// Rot färbende Mutationen am geprüften Code: die zwei `if`-Zweige am Ende von
// `Classify` vertauschen (Fall „derselbe Wert in beiden Klassen“ liefert
// `Reader`); die Wache `token == ""` entfernen (Fälle „leerer Wert“ liefern
// `Admin` bzw. `Reader`); `matchAny` nach dem ersten Treffer abbrechen
// ändert kein Ergebnis dieser Tabelle, die Zeitkonstanz bleibt eine
// Erwartung.
func TestClassify(t *testing.T) {
	cases := []struct {
		name   string
		reader []string
		admin  []string
		token  string
		want   apiauth.Role
	}{
		{"erstes Reader-Token", []string{"r1", "r2"}, []string{"a1"}, "r1", apiauth.Reader},
		{"zweites Reader-Token", []string{"r1", "r2"}, []string{"a1"}, "r2", apiauth.Reader},
		{"erstes Admin-Token", []string{"r1"}, []string{"a1", "a2"}, "a1", apiauth.Admin},
		{"zweites Admin-Token", []string{"r1"}, []string{"a1", "a2"}, "a2", apiauth.Admin},
		{"entferntes Token", []string{"r2"}, []string{"a1"}, "r1", apiauth.None},
		{"unbekanntes Token", []string{"r1"}, []string{"a1"}, "x", apiauth.None},
		{"Präfix eines Tokens", []string{"reader-token"}, []string{"a1"}, "reader", apiauth.None},
		{"Verlängerung eines Tokens", []string{"reader"}, []string{"a1"}, "reader-token", apiauth.None},
		{"derselbe Wert in beiden Klassen", []string{"gleich"}, []string{"gleich"}, "gleich", apiauth.Admin},
		{"derselbe Wert, zweites Element", []string{"r1", "gleich"}, []string{"a1", "gleich"}, "gleich", apiauth.Admin},
		{"leerer Wert, beide Klassen leer", nil, nil, "", apiauth.None},
		{"leerer Wert gegen leeres Element der Reader-Klasse", []string{""}, nil, "", apiauth.None},
		{"leerer Wert gegen leeres Element der Admin-Klasse", nil, []string{""}, "", apiauth.None},
		{"leerer Wert, gesetzte Klassen", []string{"r1"}, []string{"a1"}, "", apiauth.None},
		{"beliebiges Token, beide Klassen ungesetzt", nil, nil, "beliebig", apiauth.None},
		{"Komma bleibt Teil des Tokens", []string{"a,b"}, nil, "a,b", apiauth.Reader},
		{"Komma-Token wird nicht zerschnitten", []string{"a,b"}, nil, "a", apiauth.None},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := apiauth.New(tc.reader, tc.admin).Classify(tc.token)
			if got != tc.want {
				t.Fatalf("Classify(%q) = %v, erwartet %v (Reader %q, Admin %q)", tc.token, got, tc.want, tc.reader, tc.admin)
			}
		})
	}
}

// TestRoleOrdnung trägt die Hierarchie: `Admin` deckt `Reader` ab, `None` ist
// die niedrigste Klasse.
func TestRoleOrdnung(t *testing.T) {
	if !(apiauth.Admin > apiauth.Reader && apiauth.Reader > apiauth.None) {
		t.Fatalf("Ordnung Admin %d > Reader %d > None %d verletzt", apiauth.Admin, apiauth.Reader, apiauth.None)
	}
}
