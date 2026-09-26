package sqlexec

import (
	stderrors "errors"
	"testing"

	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

// rejectionMessage trägt für einen Konstruktor-Grund ohne eigenen Klartext den
// allgemeinen Text `Antrag ist ungültig` mit der Kennung. Der Fall entsteht
// über die Zeilen der Abfrage nicht (jeder Grund des Konstruktors hat einen
// Klartext, `TestReadPendingRequestsPassesRejectedRowsThrough`); der Test ruft
// die Funktion mit einem Fehler außerhalb dieser Gründe auf. Rot färbende
// Mutation: den Anfangswert von `clear` durch die leere Zeichenkette ersetzen —
// der Text beginnt mit dem Doppelpunkt.
func TestRejectionMessageFallsBackToGeneralText(t *testing.T) {
	got := rejectionMessage(stderrors.New("künftiger Grund"), "req-1", "src-1", "public", "feed", "spalte")

	if want := "Antrag ist ungültig: req-1"; got != want {
		t.Fatalf("rejectionMessage = %q, wollen %q", got, want)
	}
}

// Der Klartext `Spaltenname ist leer` gehört zu einer leeren Spalte mit dem
// Grund `ErrEmptyIdentifier`; derselbe Grund bei gesetzter Spalte und jeder
// andere Grund bei leerer Spalte tragen den allgemeinen Text — das letzte
// Paar der Reihenfolge von `rejectionMessage` (Antragsart, Spalte, allgemein).
// Die Fälle sind über die Zeilen der Abfrage nicht herstellbar (der
// Konstruktor nennt `ErrEmptyIdentifier` bei gesetzter Spalte nur mit leerer
// Kennung, Quelle, Schema oder Tabelle, und diese Fälle stehen davor); der
// Test ruft die Funktion mit den Grund-Feld-Paaren auf. Rot färbende Mutationen:
// die Bedingung `column == ""` aus dem Fall der Spalte entfernen — der Fall mit
// gesetzter Spalte nennt die Spalte; die Bedingung auf den Grund
// (`ErrEmptyIdentifier`) entfernen — der Fall mit fremdem Grund nennt die Spalte.
func TestRejectionMessageColumnCaseNeedsEmptyColumnAndEmptyIdentifierCause(t *testing.T) {
	for _, tc := range []struct {
		name   string
		cause  error
		column string
		want   string
	}{
		{"leere Spalte, Grund leerer Bezeichner", domainerrors.ErrEmptyIdentifier, "", "Spaltenname ist leer: req-1"},
		{"gesetzte Spalte, Grund leerer Bezeichner", domainerrors.ErrEmptyIdentifier, "spalte", "Antrag ist ungültig: req-1"},
		{"leere Spalte, fremder Grund", stderrors.New("künftiger Grund"), "", "Antrag ist ungültig: req-1"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := rejectionMessage(tc.cause, "req-1", "src-1", "public", "feed", tc.column); got != tc.want {
				t.Fatalf("rejectionMessage = %q, wollen %q", got, tc.want)
			}
		})
	}
}
