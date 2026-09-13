// Package capture trägt den Capture Use Case (`ARC-002`): die
// Orchestrierung von Persistenz und Source-ACK liegt im Application Layer
// (`ADR-0027`) — der Stream-Adapter entscheidet nicht selbst, wann eine
// Position dauerhaft verarbeitet ist.
package capture

import (
	"context"
	stderrors "errors"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

// CaptureCommand und CaptureResult sind die Transport-Typen des Capture
// Use Cases (`ADR-0042`). Ihre Definition liegt am Inbound-Port — der Port
// trägt seinen Vertrag einschließlich der Transport-Typen, Ports
// referenzieren nur die Domain —; die Use-Case-Adressen sind Aliase
// desselben Typs.
type (
	CaptureCommand = inbound.CaptureCommand
	CaptureResult  = inbound.CaptureResult
)

// ErrMissingTransaction meldet, dass der Capture-Aufruf keine
// Quelltransaktion trägt; der Port-Kontrakt verlangt eine vollständige,
// committed Quelltransaktion.
var ErrMissingTransaction = stderrors.New("Capture ohne Quelltransaktion")

// CaptureService implementiert `inbound.CaptureInboundPort` und hält die
// Persist-before-ACK-Ordnung an einer Stelle (`ADR-0027`):
//
//	Receive → Decode → Persist → COMMIT Store → ACK Source → Notify (best effort)
//
// (`LH-QA-REL-001.a`, `ADR-0055` für den optionalen letzten Schritt)
type CaptureService struct {
	store  outbound.ChangeStorePort
	ack    outbound.ReplicationAckPort
	notify outbound.ChangeNotificationPort
	log    outbound.LogPort
}

// Option konfiguriert `CaptureService` bei der Konstruktion (`NewCaptureService`);
// variadisch, damit bestehende Aufrufstellen unverändert kompilieren
// (dasselbe Muster wie `postgresack.Option`).
type Option func(*CaptureService)

// WithChangeNotification injiziert den optionalen `ChangeNotificationPort`
// (`ADR-0055`): ungesetzt bleibt das Wecksignal-Feature deaktiviert, kein
// Notify-Versuch — das bestehende Verhalten bleibt für jeden Aufrufer ohne
// diese Option unverändert.
func WithChangeNotification(notify outbound.ChangeNotificationPort) Option {
	return func(s *CaptureService) { s.notify = notify }
}

// WithLog injiziert den `LogPort` (`ADR-0024`) für den Notify-Fehlerpfad —
// dasselbe Muster wie `postgresack.WithLog`. Ungesetzt bleibt die
// Protokollierung beim No-Op (`outbound.NoopLog`).
func WithLog(log outbound.LogPort) Option {
	return func(s *CaptureService) { s.log = log }
}

// NewCaptureService verdrahtet den Capture Use Case mit seinen beiden
// obligatorischen Outbound-Ports; der `ChangeNotificationPort` ist optional
// (`WithChangeNotification`, `ADR-0055`).
func NewCaptureService(store outbound.ChangeStorePort, ack outbound.ReplicationAckPort, opts ...Option) *CaptureService {
	s := &CaptureService{store: store, ack: ack, log: outbound.NoopLog}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

var _ inbound.CaptureInboundPort = (*CaptureService)(nil)

// Capture persistiert die committed Quelltransaktion über den
// `ChangeStorePort` und bestätigt ihre Commit-Position erst nach dem
// Store-Commit über den `ReplicationAckPort`. Ein Persistenzfehler endet
// ohne Source-ACK (`SPEC-008`, Klasse `storage`); ein Fehler beim ACK
// lässt die Persistenz bestehen — der Crash zwischen Persistenz und ACK
// erzeugt höchstens erneute Verarbeitung, Wiederholung wird gegenüber
// möglichem Datenverlust bevorzugt (`ADR-0012`).
func (s *CaptureService) Capture(ctx context.Context, command CaptureCommand) (CaptureResult, error) {
	tx := command.Transaction
	if tx == nil {
		return CaptureResult{}, ErrMissingTransaction
	}
	// Offene Transaktionen sind nicht konsumierbar (`ADR-0029`, Regel 3);
	// zurückgerollte Transaktionen erreichen den Commit nicht
	// (`LH-FA-CAP-007`). Die Position liest der Pfad nur am committed
	// Transaktionsträger.
	position, committed := tx.CommitPosition()
	if !committed {
		return CaptureResult{}, domainerrors.ErrTransactionNotCommitted
	}
	if err := s.store.PersistTransaction(ctx, tx); err != nil {
		// Persistenzfehler → kein Source-ACK (`LH-QA-REL-001.a`).
		return CaptureResult{}, err
	}
	if err := s.ack.Acknowledge(ctx, position); err != nil {
		return CaptureResult{}, err
	}
	// Notify läuft NACH der Bestätigung und best-effort (`ADR-0055`): sein
	// Fehler geht nie in den Rückgabewert dieses Aufrufs ein — die bereits
	// erfolgte Persistierung und Bestätigung bleiben davon unberührt. Ohne
	// konfigurierten Port (`s.notify == nil`) unterbleibt der Versuch
	// vollständig, bestehendes Verhalten bleibt bit-identisch.
	if s.notify != nil {
		if err := s.notify.Notify(ctx, string(position.SourceID)); err != nil {
			s.log.Warn(ctx, "capture: Wecksignal fehlgeschlagen", "error", err, "source_id", position.SourceID)
		}
	}
	return CaptureResult{Acknowledged: position}, nil
}
