// Package natsstream trägt den `Publisher` als dritten `Broadcaster`-
// Abonnenten (`ADR-0100` Teilfrage 1): ein Driven-Adapter, der jeden vom
// `Broadcaster` verteilten Change als vollständiges JSON-Event über eine
// bestehende Core-NATS-Verbindung veröffentlicht — auf einem eigenen
// Subjekt-Namensraum (`cdc.stream.<source_id>.<schema>.<table>`,
// Teilfrage 2), fire-and-forget ohne eigene Zustellgarantie (Teilfrage 3).
// Die Nachvollziehbarkeit bleibt beim bestehenden Lesezugriffsweg
// — dieselbe Isolation wie die beiden bestehenden
// Zustellwege (gRPC, SSE): ein Publish-Fehlschlag bleibt lokal, er erreicht
// weder den `CaptureService` noch den kritischen Erfassungspfad.
// Das Wecksignal (`internal/adapters/driven/natsnotify`)
// bleibt davon byte-identisch unberührt — eigener Subjekt-Namensraum,
// eigenes Paket, kein Adapter-→-Adapter-Import (`ADR-0100` Teilfrage 1).
package natsstream

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/nats-io/nats.go"

	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// subjectPrefix trägt das Subjekt-Schema `cdc.stream.<source_id>.<schema>.<table>`
// (`ADR-0100` Teilfrage 2) — ein eigener Wurzel-Token, niemals `cdc.changes`
// (der Wecksignal-Namensraum bleibt unberührt).
const subjectPrefix = "cdc.stream."

// routeSubjectPrefix trägt das Zusatz-Subjekt `cdc.route.<source_id>.<ziel>`
// (`ADR-0137`) — ein eigener Wurzel-Token neben `cdc.stream` und
// `cdc.changes`.
const routeSubjectPrefix = "cdc.route."

// subjectPublisher trägt die eine Fähigkeit der NATS-Verbindung, die der
// Publisher braucht; `*nats.Conn` erfüllt sie.
type subjectPublisher interface {
	Publish(subject string, data []byte) error
}

// reservedSubjectChars trägt dieselben NATS-Subjekt-Sonderzeichen wie
// `natsnotify` — hier eigenständig geführt, kein
// Adapter-Paket importiert ein anderes (`ADR-0100` Teilfrage 1).
const reservedSubjectChars = ".*>"

// ErrPublish trägt die Konstruktions- und Validierungsgrenze dieses
// Pakets. Kein Fehlerklassen-Sentinel eines Outbound Ports: der `Publisher`
// implementiert keinen Port, er ist ein reiner `Broadcaster`-Abonnent
// (`ADR-0100` Teilfrage 1).
var ErrPublish = errors.New("natsstream: Publisher ohne gültige Verdrahtung")

// changeSubscriber trägt die eine Fähigkeit, die dieser Adapter vom
// `Broadcaster` braucht (`ADR-0100` Teilfrage 1) — dieselbe lokal
// deklarierte Schnittstelle wie beim gRPC-Server und dem SSE-Handler
// (`internal/adapters/driving/grpc/server.go`, `internal/adapters/driving/http/sse.go`).
// `*grpcstream.Broadcaster` erfüllt sie strukturell; die Composition Root
// verdrahtet ihn.
type changeSubscriber interface {
	Subscribe() (<-chan *model.Change, func())
}

// Option konfiguriert den Publisher bei der Konstruktion (`New`) — aktuell
// nur der optionale `LogPort` (`ADR-0024`), dasselbe Muster
// wie `natsnotify.Option`.
type Option func(*options)

type options struct {
	log outbound.LogPort
}

func newOptions(opts []Option) options {
	o := options{log: outbound.NoopLog}
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

// WithLog injiziert den `LogPort` (`ADR-0024`) — dieselbe Rolle wie
// `natsnotify.WithLog`, eigenständig definiert (kein Adapter-→-Adapter-Import).
func WithLog(log outbound.LogPort) Option {
	return func(o *options) { o.log = log }
}

// Publisher veröffentlicht jeden vom `Broadcaster` verteilten Change als
// vollständiges JSON-Event auf `cdc.stream.<source_id>.<schema>.<table>`
// (`ADR-0100`) und, bei gesetztem Zustellziel, mit demselben Payload zusätzlich
// auf `cdc.route.<source_id>.<ziel>`. Er trägt keinen eigenen Zustand über den
// Veröffentlichungszeitpunkt hinaus und keine Zustellgarantie — dieselbe
// Fire-and-Forget-Haltung wie der `Broadcaster` selbst (Teilfrage 3).
type Publisher struct {
	conn       subjectPublisher
	subscriber changeSubscriber
	sourceID   string
	log        outbound.LogPort
}

// New legt den Publisher auf eine bestehende NATS-Verbindung, den
// abonnierten `Broadcaster` und die konfigurierte Quelle fest; die
// Composition Root verdrahtet Verbindungsaufbau und Aktivierung
// (`ADR-0100` Teilfrage 5).
func New(conn *nats.Conn, subscriber changeSubscriber, sourceID string, opts ...Option) (*Publisher, error) {
	if conn == nil {
		return nil, fmt.Errorf("%w: keine NATS-Verbindung", ErrPublish)
	}
	if subscriber == nil {
		return nil, fmt.Errorf("%w: kein Broadcaster zum Abonnieren", ErrPublish)
	}
	if sourceID == "" {
		return nil, fmt.Errorf("%w: keine Quelle", ErrPublish)
	}
	o := newOptions(opts)
	// Kein `context.Context`-Parameter an dieser Konstruktionsstelle
	// (Signatur-Vertrag von `New`, analog `natsnotify.New`) — der
	// Log-Aufruf trägt deshalb `context.Background()`.
	o.log.Info(context.Background(), "natsstream: verdrahtet")
	return &Publisher{conn: conn, subscriber: subscriber, sourceID: sourceID, log: o.log}, nil
}

// Run abonniert den `Broadcaster` und veröffentlicht jeden eintreffenden
// Change, bis `ctx` endet — dieselbe Struktur wie die übrigen
// Hintergrund-Züge der Composition Root (`internal/bootstrap`), hier ohne
// Ticker: der Auslöser ist der Broadcaster-Kanal selbst. Die Abmeldung
// läuft beim Verlassen (`defer cancel()`).
func (p *Publisher) Run(ctx context.Context) {
	changes, cancel := p.subscriber.Subscribe()
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return
		case change := <-changes:
			if change == nil {
				continue
			}
			p.publish(ctx, change)
		}
	}
}

// streamMessage trägt dasselbe Nachrichtenschema wie die SSE-Ereignisse
// (`ADR-0100` Teilfrage 2) — dieselben zehn Felder wie
// `model.Change` ohne `Origin`, hier eigenständig
// geführt (kein Adapter-→-Adapter-Import, `ADR-0100` Teilfrage 1). Die
// Row Images stehen als eingebettete JSON-Werte; ein fehlendes Bild
// (Boundary) wird zu `null`.
type streamMessage struct {
	ChangeID      string          `json:"change_id"`
	TransactionID string          `json:"transaction_id"`
	SourceTableID string          `json:"source_table_id"`
	Sequence      int64           `json:"sequence"`
	Operation     string          `json:"operation"`
	OldImage      json.RawMessage `json:"old_image"`
	NewImage      json.RawMessage `json:"new_image"`
	SchemaVersion string          `json:"schema_version"`
	Schema        string          `json:"schema"`
	Table         string          `json:"table"`
}

// rowImage liefert ein Row Image als eingebetteten JSON-Wert; ein leeres
// Bild wird zu `null`, damit die Nachricht auch ohne Bild gültiges JSON
// bleibt (`LH-FA-CAP-008` Boundary) — dieselbe Übersetzung wie im
// SSE-Adapter.
func rowImage(image []byte) json.RawMessage {
	if len(image) == 0 {
		return json.RawMessage("null")
	}
	return json.RawMessage(image)
}

// toStreamMessage übersetzt einen Domain-Change in seine NATS-Nachrichtenform
// (`SPEC-021`): dieselben Felder wie das SSE-Event, Row Images
// unverändert übernommen.
func toStreamMessage(change *model.Change) streamMessage {
	return streamMessage{
		ChangeID:      string(change.ID),
		TransactionID: string(change.TransactionID),
		SourceTableID: string(change.SourceTableID),
		Sequence:      change.Sequence,
		Operation:     string(change.Operation),
		OldImage:      rowImage(change.OldImage),
		NewImage:      rowImage(change.NewImage),
		SchemaVersion: string(change.SchemaVersion),
		Schema:        change.Schema,
		Table:         change.Table,
	}
}

// containsReservedSubjectToken meldet, ob `token` ein NATS-Subjekt-
// Sonderzeichen oder Whitespace trägt (`ADR-0056` Folgepflicht) — dieselbe
// defensive Validierung wie `natsnotify`, hier eigenständig geführt.
func containsReservedSubjectToken(token string) bool {
	if strings.ContainsAny(token, reservedSubjectChars) {
		return true
	}
	for _, r := range token {
		if unicode.IsSpace(r) {
			return true
		}
	}
	return false
}

// subjectFor baut das Subjekt für Schema/Tabelle unter der konfigurierten
// Quelle (`ADR-0100` Teilfrage 2) — eine reine Funktion, unabhängig von der
// NATS-Verbindung: sie trägt die Regressionsgrenze „Wurzel-Token niemals
// cdc.changes" testbar ohne jede Verbindung.
func subjectFor(sourceID, schema, table string) string {
	return subjectPrefix + sourceID + "." + schema + "." + table
}

// routeSubjectFor baut das Zusatz-Subjekt eines Zustellziels unter der
// konfigurierten Quelle — eine reine Funktion wie `subjectFor`.
func routeSubjectFor(sourceID string, target model.RouteTarget) string {
	return routeSubjectPrefix + sourceID + "." + string(target)
}

// publish veröffentlicht einen einzelnen Change als vollständiges
// JSON-Event (`ADR-0100`): zuerst auf dem Tabellen-Subjekt, danach — nur bei
// gesetztem Zustellziel — mit demselben Payload auf dem Ziel-Subjekt. Die
// beiden Veröffentlichungen sind in ihrer Prüfung unabhängig: ein übersprungenes
// Tabellen-Subjekt (leerer oder reservierter Name) lässt das Ziel-Subjekt
// bestehen, ein übersprungenes oder gescheitertes Ziel-Subjekt lässt das
// Tabellen-Subjekt bestehen. Ein leerer Schema-/Tabellenname, ein Name mit
// NATS-reserviertem Zeichen oder Whitespace, ein Kodierfehler oder ein
// Publish-Fehlschlag bleiben lokal — kein Rückgabewert, kein propagierter
// Fehler: derselbe Fire-and-Forget-Vertrag wie der `Broadcaster` selbst. Die
// Prüfung in `tableSubject` lässt einen leeren Relationsnamen nicht als
// verkürztes Subjekt (`cdc.stream.<source>.<schema>.`) durch. Ein
// Fehlschlag erreicht weder den `CaptureService` noch den kritischen
// Erfassungspfad — der Aufrufer liest ausschließlich aus dem bereits
// isolierten Broadcaster-Kanal.
func (p *Publisher) publish(ctx context.Context, change *model.Change) {
	payload, err := json.Marshal(toStreamMessage(change))
	if err != nil {
		p.log.Warn(ctx, "natsstream: Change nicht kodierbar — Publish übersprungen", "error", err)
		return
	}
	if subject, ok := p.tableSubject(ctx, change); ok {
		p.send(ctx, subject, payload, change)
	}
	if change.RouteTarget == "" {
		return
	}
	if subject, ok := p.routeSubject(ctx, change); ok {
		p.send(ctx, subject, payload, change)
	}
}

// tableSubject liefert das Tabellen-Subjekt des Changes; ok ist false, wenn
// Schema oder Tabelle leer sind oder ein reserviertes Zeichen tragen.
func (p *Publisher) tableSubject(ctx context.Context, change *model.Change) (string, bool) {
	if change.Schema == "" || change.Table == "" {
		p.log.Warn(ctx, "natsstream: Schema oder Tabelle leer — Publish übersprungen",
			"schema", change.Schema, "table", change.Table)
		return "", false
	}
	if containsReservedSubjectToken(change.Schema) || containsReservedSubjectToken(change.Table) {
		p.log.Warn(ctx, "natsstream: Schema/Tabelle trägt ein NATS-reserviertes Zeichen (.,*,>) oder Whitespace — Publish übersprungen",
			"schema", change.Schema, "table", change.Table)
		return "", false
	}
	return subjectFor(p.sourceID, change.Schema, change.Table), true
}

// routeSubject liefert das Ziel-Subjekt des Changes (`ADR-0137`); ok ist
// false, wenn der Zielname ein reserviertes Zeichen oder Whitespace trägt.
// Das Alphabet des Zielnamens schließt beides aus; die Prüfung bindet die
// Subjekt-Bildung unabhängig vom Alphabet des Aufrufers.
func (p *Publisher) routeSubject(ctx context.Context, change *model.Change) (string, bool) {
	if containsReservedSubjectToken(string(change.RouteTarget)) {
		p.log.Warn(ctx, "natsstream: Zielname trägt ein NATS-reserviertes Zeichen (.,*,>) oder Whitespace — Publish übersprungen",
			"target", string(change.RouteTarget))
		return "", false
	}
	return routeSubjectFor(p.sourceID, change.RouteTarget), true
}

// send veröffentlicht den Payload auf einem Subjekt; ein Fehlschlag bleibt
// lokal (Warnung im Log).
func (p *Publisher) send(ctx context.Context, subject string, payload []byte, change *model.Change) {
	if err := p.conn.Publish(subject, payload); err != nil {
		p.log.Warn(ctx, "natsstream: Publish fehlgeschlagen", "error", err, "subject", subject)
		return
	}
	p.log.Debug(ctx, "natsstream: Change publiziert", "subject", subject, "change_id", change.ID)
}
