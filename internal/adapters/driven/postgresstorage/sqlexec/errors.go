package sqlexec

import (
	stderrors "errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// Classify trägt die Fehlerklassen-Übersetzung an der Adapter-Grenze
// (`ADR-0023`, `SPEC-008`): die technische Ursache geht in die Klasse des
// Ports, bleibt über die zweite Wrappung aber lesbar — Application und
// Betrieb klassifizieren über `errors.Is`, ohne einen Treibertyp zu kennen.
// Die Funktion ist rein: sie liest ihre beiden Eingänge und berührt keinen
// Zustand.
func Classify(class error, cause error) error {
	return fmt.Errorf("%w: %w", class, cause)
}

// IsAbsent meldet die Abwesenheit einer Zeile als definierten Zustand: der
// Treiber führt sie als `pgx.ErrNoRows`, die Adapter lesen sie als Nullwert
// bzw. als „nicht vorhanden" (`LH-FA-CON-005` Boundary, `LH-FA-CFG-003`
// Boundary) — kein Datenbankfehler. Rein: `errors.Is` über die Ursache.
func IsAbsent(err error) bool {
	return stderrors.Is(err, pgx.ErrNoRows)
}
