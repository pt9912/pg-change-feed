package main

import "testing"

// TestSubscribeSubjectDefaultsToRootWildcard prüft den unveränderten
// Default: ohne `-source` und `-target` abonniert das Beispiel den
// Wurzel-Wildcard.
func TestSubscribeSubjectDefaultsToRootWildcard(t *testing.T) {
	got, err := SubscribeSubject("", "")
	if err != nil || got != "cdc.stream.>" {
		t.Fatalf("SubscribeSubject = %q, %v, want cdc.stream.>", got, err)
	}
}

// TestSubscribeSubjectDerivesTargetSubject prüft die Bindung Flag → Subjekt:
// Quelle und Ziel ergeben das Zusatz-Subjekt des Zustellziels.
func TestSubscribeSubjectDerivesTargetSubject(t *testing.T) {
	got, err := SubscribeSubject("quelle-1", "eu")
	if err != nil || got != "cdc.route.quelle-1.eu" {
		t.Fatalf("SubscribeSubject = %q, %v, want cdc.route.quelle-1.eu", got, err)
	}
}

// TestSubscribeSubjectRejectsHalfAndInvalidInput prüft die Eingabefehler: eines
// der beiden Flags allein, leerer Leerraum, Trennzeichen und Platzhalter.
func TestSubscribeSubjectRejectsHalfAndInvalidInput(t *testing.T) {
	for _, in := range [][2]string{
		{"quelle-1", ""}, {"", "eu"},
		{"quelle-1", "   "}, {"quelle-1", "a.b"}, {"quelle-1", "*"}, {"quelle-1", ">"}, {"quelle-1", "a b"},
		{"a.b", "eu"}, {"*", "eu"}, {">", "eu"},
	} {
		if got, err := SubscribeSubject(in[0], in[1]); err == nil {
			t.Fatalf("SubscribeSubject(%q, %q) = %q, erwarteter Fehler blieb aus", in[0], in[1], got)
		}
	}
}
