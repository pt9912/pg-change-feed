package outbound

import (
	"context"
	stderrors "errors"

	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// ErrConsumerStateStorage trägt die Fehlerklasse `storage` dieses Ports
// (`SPEC-008`, `ADR-0023`): ein Persistenzfehler am Consumer-State-Speicher
// endet sichtbar, Application und Betrieb klassifizieren über `errors.Is`
// und kennen keinen Treibertyp. Abgrenzung: der ChangeStore-Sentinel
// (`ErrStorage`) trägt dieselbe Klasse am `ChangeStorePort` samt seiner
// Klasse-Aktion (kein Source-ACK, `LH-QA-REL-001.a`) — diese Aktion hat am
// Consumer-ACK keinen Träger, der Consumer-ACK bestätigt keine
// Quellposition; deshalb führt dieser Port seine Klasse als eigenen
// Sentinel.
var ErrConsumerStateStorage = stderrors.New("Fehlerklasse storage: Persistenzfehler im Consumer-State-Speicher")

// ErrConsumerUnregistered trägt die Registrierungs-Grenze der Bestätigung
// (`LH-FA-CON-004`): ein ACK ohne registrierte Kennung ist ein
// Aufrufvertrags-Verstoß und endet sichtbar über diesen Sentinel, nicht
// still — die Rest-Grenze zwischen Prüfung und Schreiben trägt der
// Fremdschlüssel der DDL (`SPEC-001`) über `ErrConsumerStateStorage`.
var ErrConsumerUnregistered = stderrors.New("Consumer ist nicht registriert")

// ConsumerStatePort trägt die Zustands-Fähigkeit der Consumer
// (`ARC-004`, Fähigkeits-Port je `ADR-0034`): die unabhängigen
// persistierten Positionen benannter Consumer (`ADR-0013`, `LH-FA-CON-002`,
// `LH-FA-CON-003`) — Registrierung, Bestätigung, Fortsetzungs-Lese und
// Entfernung liegen an einer Konsistenzgrenze.
//
// Die Zeile `cdc.consumer_position` (`SPEC-001`) trägt nur bestätigte
// Positionen: ihr Bestehen liest sich als Bestätigung, ihre Abwesenheit
// als „noch nichts bestätigt" — die definierte Anfangsposition des
// Fortsetzungs-Lese (`LH-FA-CON-005` Boundary).
//
// Der Port bestätigt keine Quellpositionen; die Quell-Bestätigung ist die
// Wirkung des `ReplicationAckPort` (`ADR-0007`). Consumer-ACK verläuft
// regulär nur vorwärts (`ADR-0029`, Regel 2): die Bestätigung liest den
// gespeicherten Fortschritt und führt ihn über die Domänen-Ordnung —
// eine frühere Position und eine Position einer anderen Quelle tragen
// ihre Invarianten-Sentinels (`ErrPositionRegression`,
// `ErrSourceMismatch`), die Wiederholung derselben Position ist
// idempotent (`LH-FA-CON-004` Boundary).
//
// Die Bestätigung setzt die Registrierung voraus (`SPEC-001`
// Fremdschlüssel); ein Aufruf ohne Registrierung endet sichtbar über
// `ErrConsumerUnregistered`, nicht still — zwischen Prüfung und Schreiben
// trägt der Fremdschlüssel die Rest-Grenze über `ErrConsumerStateStorage`.
type ConsumerStatePort interface {
	// Register trägt die Consumer-Zeile ein; die Rückkehr meldet, ob der
	// Aufruf neu registriert hat. Eine bereits registrierte Kennung
	// bleibt unverändert (`LH-FA-CON-001` Boundary: idempotentes
	// Verhalten).
	Register(ctx context.Context, consumer model.Consumer) (bool, error)

	// Position liest die bestätigte Position des Consumers
	// (`LH-FA-CON-003` Happy Path); der Nullwert liest die definierte
	// Anfangsposition eines Consumers ohne Bestätigung
	// (`LH-FA-CON-005` Boundary — auch nach einem Neustart, der Lese
	// läuft gegen denselben Bestand).
	Position(ctx context.Context, consumer model.ConsumerID) (model.ConsumerPosition, error)

	// Acknowledge trägt die bestätigte Position fort (`LH-FA-CON-004`):
	// der gespeicherte Fortschritt liest vor dem Schreiben, die Monotonie
	// und die Quell-Bindung trägt der Domänen-Vergleich
	// (`model.ConsumerPosition.Advance`); die Rückkehr liest den
	// fortgeführten Stand.
	Acknowledge(ctx context.Context, position model.ConsumerPosition) (model.ConsumerPosition, error)

	// Remove entzieht die Consumer-Zeile administrativ
	// (`LH-FA-CON-006`): die bestätigte Position geht mit — ihr
	// dokumentiertes Verhalten ist der Entzug mit der Zeile, die
	// Entfernung löscht keine Changes. Die Rückkehr meldet, ob der
	// Consumer entfernt war; ein erneuter Aufruf bleibt ohne Wirkung
	// (Idempotenz). Der Rolle des Consumers in der Retention
	// (`LH-FA-RET-004`) trägt die Zeilen-Abwesenheit Rechnung.
	Remove(ctx context.Context, consumer model.ConsumerID) (bool, error)
}
