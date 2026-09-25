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
	"github.com/pt9912/pg-change-feed/internal/domain/model"
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

// CaptureService implementiert `inbound.CaptureInboundPort` und
// `inbound.IdleConfirmationInboundPort` und hält die
// Persist-before-ACK-Ordnung an einer Stelle (`ADR-0027`):
//
//	Receive → Decode → Persist → COMMIT Store → ACK Source →
//	Notify (best effort) → Stream-Publish (best effort)
//
// (`LH-QA-REL-001.a`, `ADR-0055` für das optionale Wecksignal,
// `ADR-0056` für dessen tabellen-granulare Deduplizierung, `ADR-0060`
// Teilfrage 2 für den optionalen Stream-Publish nach dem Wecksignal)
type CaptureService struct {
	store  outbound.ChangeStorePort
	ack    outbound.ReplicationAckPort
	notify outbound.ChangeNotificationPort
	stream outbound.ChangeStreamPort
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

// WithChangeStream injiziert den optionalen `ChangeStreamPort`
// (`ADR-0060` Teilfrage 2): ungesetzt bleibt das Live-Streaming deaktiviert,
// kein Publish-Versuch — das bestehende Verhalten bleibt für jeden Aufrufer
// ohne diese Option unverändert. Der Aufruf reiht sich nach `ACK Source` und
// nach dem Notify-Schritt ein und trägt die Zustellsemantik des Ports: er
// blockiert nicht auf einen Abonnenten, sein Fehler geht nicht in den
// Rückgabewert von `Capture` ein (`ADR-0066`).
func WithChangeStream(stream outbound.ChangeStreamPort) Option {
	return func(s *CaptureService) { s.stream = stream }
}

// WithLog injiziert den `LogPort` (`ADR-0024`) für den Notify-Fehlerpfad —
// dasselbe Muster wie `postgresack.WithLog`. Ungesetzt bleibt die
// Protokollierung beim No-Op (`outbound.NoopLog`).
func WithLog(log outbound.LogPort) Option {
	return func(s *CaptureService) { s.log = log }
}

// NewCaptureService verdrahtet den Capture Use Case mit seinen beiden
// obligatorischen Outbound-Ports; `ChangeNotificationPort` und
// `ChangeStreamPort` sind optional (`WithChangeNotification`, `ADR-0055`;
// `WithChangeStream`, `ADR-0060` Teilfrage 2).
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
	// vollständig, bestehendes Verhalten bleibt bit-identisch. Ein Notify
	// je distinkter `(schema, table)`-Paarung der Transaktion, dedupliziert
	// (`ADR-0056`): mehrere Changes derselben Tabelle lösen genau ein
	// Signal aus.
	if s.notify != nil {
		for _, table := range distinctTables(tx) {
			if err := s.notify.Notify(ctx, string(position.SourceID), table.schema, table.table); err != nil {
				s.log.Warn(ctx, "capture: Wecksignal fehlgeschlagen", "error", err, "source_id", position.SourceID, "schema", table.schema, "table", table.table)
			}
		}
	}
	// Der Stream-Publish läuft als letzter Schritt der Best-Effort-Kette
	// (`ADR-0060` Teilfrage 2): genau ein Aufruf je Change der Transaktion in
	// der Reihenfolge von `tx.Changes()`, ohne Deduplizierung nach Tabelle —
	// `LH-FA-SST-008` verlangt den vollständigen Inhalt je Zeilen-Change. Der
	// Aufruf blockiert nicht auf einen Abonnenten; die Entkopplung trägt der
	// Port (`ADR-0066`), deshalb braucht diese Aufrufstelle keine eigene
	// Zeit-Isolation. Sein Fehler geht nie in den Rückgabewert dieses Aufrufs
	// ein, die bereits erfolgte Persistierung und Bestätigung bleiben
	// unberührt. Ohne konfigurierten Port (`s.stream == nil`) unterbleibt der
	// Versuch vollständig.
	if s.stream != nil {
		changes := changesOfCommittedTransaction(tx)
		for i := range changes {
			if err := s.stream.Publish(ctx, &changes[i]); err != nil {
				s.log.Warn(ctx, "capture: Stream-Publish fehlgeschlagen", "error", err, "change_id", changes[i].ID)
			}
		}
	}
	return CaptureResult{Acknowledged: position}, nil
}

// ErrMissingIdlePosition meldet, dass die Leerlauf-Bestätigung keine
// Position trägt.
var ErrMissingIdlePosition = stderrors.New("Leerlauf-Bestätigung ohne Position")

// ConfirmIdle bestätigt die gemeldete Leerlauf-Position über den
// `ReplicationAckPort` — ohne Persistenz, weil im Leerlauf kein Change der
// Publication zu speichern ist (`ADR-0120` Festlegung 1): der Aufruf erreicht
// nie den `ChangeStorePort`, benachrichtigt nicht und veröffentlicht nichts.
// Ein Fehler des Ports geht unverändert durch (Klasse `replication`,
// `ADR-0120` Festlegung 3, `SPEC-008`); die Bestätigung gilt erst mit der
// Rückkehr ohne Fehler.
func (s *CaptureService) ConfirmIdle(ctx context.Context, command inbound.IdleConfirmationCommand) (inbound.IdleConfirmationResult, error) {
	if command.Position.IsZero() {
		return inbound.IdleConfirmationResult{}, ErrMissingIdlePosition
	}
	if err := s.ack.Acknowledge(ctx, command.Position); err != nil {
		return inbound.IdleConfirmationResult{}, err
	}
	return inbound.IdleConfirmationResult{Acknowledged: command.Position}, nil
}

var _ inbound.IdleConfirmationInboundPort = (*CaptureService)(nil)

// changesOfCommittedTransaction liefert die Changes einer Transaktion, deren
// Commit-Status der Aufrufer bereits über `CommitPosition` geprüft hat. Der
// Fehler von `tx.Changes()` (offene Transaktion) kann deshalb nicht auftreten
// und wird hier verworfen.
func changesOfCommittedTransaction(tx *model.ChangeTransaction) []model.Change {
	changes, _ := tx.Changes()
	return changes
}

// schemaTable trägt ein distinktes Schema-/Tabellenpaar einer Transaktion
// für die Notify-Deduplizierung (`ADR-0056`).
type schemaTable struct {
	schema string
	table  string
}

// distinctTables sammelt die distinkten `(schema, table)`-Paare der
// bereits committed Transaktion in erster Auftrittsreihenfolge
// (`ADR-0056`): mehrere Changes derselben Tabelle liefern genau einen
// Eintrag.
func distinctTables(tx *model.ChangeTransaction) []schemaTable {
	changes := changesOfCommittedTransaction(tx)
	seen := make(map[schemaTable]struct{}, len(changes))
	tables := make([]schemaTable, 0, len(changes))
	for _, change := range changes {
		key := schemaTable{schema: change.Schema, table: change.Table}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		tables = append(tables, key)
	}
	return tables
}
