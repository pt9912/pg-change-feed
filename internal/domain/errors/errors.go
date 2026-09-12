// Package errors trägt die Domänen-Fehler der Invarianten-Verletzungen
// (`ADR-0029`). Die Konstruktoren in `internal/domain/model` liefern sie
// als Fehlerursache; Adapter übersetzen sie an ihrer Grenze in die
// Fehlerklassen des Pflichtenhefts (§4, `SPEC-008`).
package errors

import (
	stderrors "errors"
)

// Jede Invariante aus `ADR-0029`, die ein Domänentyp erzwingt, hat hier
// einen Sentinel: die Konstruktoren wickeln ihn ein (`%w`), der Aufrufer
// klassifiziert über `errors.Is`.
var (
	// ErrEmptyIdentifier: Kennungen (Source, Tabelle, Change,
	// Transaktion, Consumer, Schema-Version) sind nicht leer.
	ErrEmptyIdentifier = stderrors.New("leere Kennung")

	// ErrInvalidOperation: Operation ist INSERT, UPDATE oder DELETE
	// (`SPEC-002`).
	ErrInvalidOperation = stderrors.New("unbekannter Operationstyp")

	// ErrNonPositiveSequence: die Sequenz innerhalb der Transaktion ist
	// mindestens 1 (`SPEC-002`, eindeutige Sequenz).
	ErrNonPositiveSequence = stderrors.New("Sequenz ist kleiner als 1")

	// ErrNonPositiveVersion: die Versionsnummer einer Schema-Version ist
	// mindestens 1 (`SPEC-004`).
	ErrNonPositiveVersion = stderrors.New("Versionsnummer ist kleiner als 1")

	// ErrInvalidPosition: eine Quellposition trägt einen Offset größer
	// als 0 (`SPEC-003`, sortierbare Position `LH-FA-DAT-004`).
	ErrInvalidPosition = stderrors.New("Position ohne Offset")

	// ErrPositionRegression: der Consumer-ACK verläuft regulär nur
	// vorwärts (`ADR-0029`, Regel 2); Wiederholung derselben Position ist
	// idempotent (`LH-FA-CON-004`).
	ErrPositionRegression = stderrors.New("Position rückt nicht vor")

	// ErrSourceMismatch: Positionen sind nur innerhalb ihrer Quelle
	// fachlich geordnet (`ADR-0005`); eine Position einer anderen Quelle
	// ist keine Bestätigungs- oder Commit-Größe für diese.
	ErrSourceMismatch = stderrors.New("Position einer anderen Quelle")

	// ErrTransactionMismatch: ein Change gehört genau einer Transaktion
	// an (`ADR-0029`, Regel 6).
	ErrTransactionMismatch = stderrors.New("Change gehört nicht zu dieser Transaktion")

	// ErrDuplicateSequence: die Sequenz ist innerhalb der Transaktion
	// eindeutig (`ADR-0029`, Regel 6).
	ErrDuplicateSequence = stderrors.New("Sequenz innerhalb der Transaktion bereits vergeben")

	// ErrTransactionAlreadyCommitted: ein Commit ist einmalig; einer
	// committed Transaktion fügt keine Changes mehr hinzu (`ADR-0029`,
	// Regel 3).
	ErrTransactionAlreadyCommitted = stderrors.New("Transaktion ist bereits committed")

	// ErrTransactionNotCommitted: offene Transaktionen sind nicht
	// konsumierbar (`ADR-0029`, Regel 3); ihre Changes sind erst nach dem
	// Commit lesbar (`LH-FA-CAP-006`).
	ErrTransactionNotCommitted = stderrors.New("Transaktion ist nicht committed")

	// ErrNegativeDuration: Zeiträume der Retention sind nicht negativ
	// (`LH-FA-RET-003`).
	ErrNegativeDuration = stderrors.New("negative Dauer")

	// ErrInvalidErrorClass: eine Fehlerklasse ist eine der sieben stabilen
	// Kategorien aus `ADR-0023` (`SPEC-008`) — eine leere oder unbekannte
	// Klasse verletzt die Invariante (`LH-FA-ADM-003`).
	ErrInvalidErrorClass = stderrors.New("unbekannte Fehlerklasse")
)
