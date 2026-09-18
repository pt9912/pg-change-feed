// Package natsstream trägt den `Publisher` als dritten `Broadcaster`-
// Abonnenten (`ADR-0100` Teilfrage 1): ein Driven-Adapter, der jeden vom
// `Broadcaster` verteilten Change als vollständiges JSON-Event über eine
// bestehende Core-NATS-Verbindung veröffentlicht — auf einem eigenen
// Subjekt-Namensraum (`cdc.stream.<source_id>.<schema>.<table>`,
// Teilfrage 2), fire-and-forget ohne eigene Zustellgarantie (Teilfrage 3).
// Die Nachvollziehbarkeit bleibt beim bestehenden Lesezugriffsweg
// (`LH-FA-REA-001` ff.) — dieselbe Isolation wie die beiden bestehenden
// Zustellwege (gRPC, SSE): ein Publish-Fehlschlag bleibt lokal, er erreicht
// weder den `CaptureService` noch den kritischen Erfassungspfad.
// `ADR-0055`/`ADR-0056`s Wecksignal (`internal/adapters/driven/natsnotify`)
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
// (`ADR-0055`/`ADR-0056`s Wecksignal-Namensraum bleibt unberührt).
const subjectPrefix = "cdc.stream."

// reservedSubjectChars trägt dieselben NATS-Subjekt-Sonderzeichen wie
// `natsnotify` (`ADR-0056` Folgepflicht) — hier eigenständig geführt, kein
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
// nur der optionale `LogPort` (`LH-QA-OPS-004`, `ADR-0024`), dasselbe Muster
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
// (`ADR-0100`). Er trägt keinen eigenen Zustand über den
// Veröffentlichungszeitpunkt hinaus und keine Zustellgarantie — dieselbe
// Fire-and-Forget-Haltung wie der `Broadcaster` selbst (Teilfrage 3).
type Publisher struct {
	conn       *nats.Conn
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
// (`SPEC-021`, `ADR-0100` Teilfrage 2) — dieselben zehn Felder wie
// `model.Change`, hier eigenständig geführt (kein Adapter-→-Adapter-Import,
// `ADR-0100` Teilfrage 1). Die Row Images stehen als eingebettete
// JSON-Werte; ein fehlendes Bild (`LH-FA-CAP-008` Boundary) wird zu `null`.
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
// (`SPEC-021`/`SPEC-024`): dieselben Felder wie das SSE-Event, Row Images
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

// publish veröffentlicht einen einzelnen Change als vollständiges
// JSON-Event (`ADR-0100` Teilfrage 2/3): ein leerer Schema-/Tabellenname,
// ein Name mit NATS-reserviertem Zeichen oder Whitespace, ein Kodierfehler
// oder ein Publish-Fehlschlag bleiben lokal — kein Rückgabewert, kein
// propagierter Fehler: derselbe Fire-and-Forget-Vertrag wie der
// `Broadcaster` selbst (`ADR-0060` Teilfrage 3). Die Leerwert- und
// Zeichen-Grenze entspricht `natsnotify.Notify` (`ADR-0056` Folgepflicht):
// ohne die Leerwert-Prüfung entstünde aus einem leeren Relationsnamen ein
// verkürztes Subjekt (`cdc.stream.<source>.<schema>.`), das still
// publiziert würde. Ein Publish-Fehlschlag erreicht weder den
// `CaptureService` noch den kritischen Erfassungspfad — der Aufrufer liest
// ausschließlich aus dem bereits isolierten Broadcaster-Kanal.
func (p *Publisher) publish(ctx context.Context, change *model.Change) {
	if change.Schema == "" || change.Table == "" {
		p.log.Warn(ctx, "natsstream: Schema oder Tabelle leer — Publish übersprungen",
			"schema", change.Schema, "table", change.Table)
		return
	}
	if containsReservedSubjectToken(change.Schema) || containsReservedSubjectToken(change.Table) {
		p.log.Warn(ctx, "natsstream: Schema/Tabelle trägt ein NATS-reserviertes Zeichen (.,*,>) oder Whitespace — Publish übersprungen",
			"schema", change.Schema, "table", change.Table)
		return
	}
	subject := subjectFor(p.sourceID, change.Schema, change.Table)
	payload, err := json.Marshal(toStreamMessage(change))
	if err != nil {
		p.log.Warn(ctx, "natsstream: Change nicht kodierbar — Publish übersprungen", "error", err)
		return
	}
	if err := p.conn.Publish(subject, payload); err != nil {
		p.log.Warn(ctx, "natsstream: Publish fehlgeschlagen", "error", err, "subject", subject)
		return
	}
	p.log.Debug(ctx, "natsstream: Change publiziert", "subject", subject, "change_id", change.ID)
}
