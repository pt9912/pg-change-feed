package sqlexec

import (
	stderrors "errors"
	"testing"
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
