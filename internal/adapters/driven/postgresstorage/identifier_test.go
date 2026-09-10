package postgresstorage

import (
	stderrors "errors"
	"strings"
	"testing"
)

// Die Bezeichner-Verweigerung läuft netzlos: die Prüfung endet vor dem
// ersten SQL-Aufruf (`SPEC-008`, Klasse `configuration` — kein Start im
// falschen Stand).
func TestValidateIdentifierVerweigert(t *testing.T) {
	for _, name := range []string{
		"",
		"Upper",
		"mit-Bindestrich",
		"mit punkt",
		"mitpunkt.",
		strings.Repeat("a", 64),
	} {
		if err := validateIdentifier(name); !stderrors.Is(err, ErrActivationConfiguration) {
			t.Fatalf("Bezeichner %q: %v (Erwartung: Fehlerklasse configuration)", name, err)
		}
	}
}

// TestValidateIdentifierTraegt trägt den gültigen Bereich des Alphabets
// samt seiner Längen-Grenze.
func TestValidateIdentifierTraegt(t *testing.T) {
	for _, name := range []string{"pub_1", "t9", strings.Repeat("a", 63)} {
		if err := validateIdentifier(name); err != nil {
			t.Fatalf("Bezeichner %q: %v (Erwartung: gültig)", name, err)
		}
	}
}

// TestPublicationIdentifiersVerweigert trägt die Verweigerung der
// Publication-DDL vor der Interpolation: ein Bezeichner außerhalb des
// Alphabets endet ohne geprüfte Literale.
func TestPublicationIdentifiersVerweigert(t *testing.T) {
	if _, _, err := publicationIdentifiers("pub_1", "Public", "t1"); !stderrors.Is(err, ErrActivationConfiguration) {
		t.Fatalf("Schema außerhalb des Alphabets: %v (Erwartung: Fehlerklasse configuration)", err)
	}
	if _, _, err := publicationIdentifiers("mit-Bindestrich", "public", "t1"); !stderrors.Is(err, ErrActivationConfiguration) {
		t.Fatalf("Publication außerhalb des Alphabets: %v (Erwartung: Fehlerklasse configuration)", err)
	}
}
