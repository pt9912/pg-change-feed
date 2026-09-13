// Package natsnotify trägt den `NatsChangeNotificationAdapter` als
// Driven-Implementierung des `ChangeNotificationPort` (`ARC-013`,
// `ADR-0055`): ein dünner Publish-Aufruf gegen Core NATS, ohne eigene
// Zustellgarantie und ohne eigenen Zustand — die Nachvollziehbarkeit
// bleibt beim bestehenden `ChangeStorePort`-Lesezugriffsweg (`SPEC-017`).
// Der Verbindungsaufbau (`nats.Connect`) bleibt Composition-Root-Detail
// (`ADR-0026`), wie schon der `PostgresReplicationAckAdapter` seine
// Replication-Verbindung von außen erhält.
package natsnotify

import (
	"context"
	"fmt"
	"strings"
	"unicode"

	"github.com/nats-io/nats.go"

	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
)

// subjectPrefix trägt das Subjekt-Schema
// `cdc.changes.<source_id>.<schema>.<table>` (`SPEC-017`, `ADR-0056`): ein
// Subjekt je Tabelle einer Quelle, konsistent mit der übrigen
// Quelle-Skopierung (Publication-/Slot-Namen).
const subjectPrefix = "cdc.changes."

// reservedSubjectChars trägt die NATS-Subjekt-Sonderzeichen, die ein
// Schema- oder Tabellenname als eigenständiges Token nicht enthalten darf
// (`ADR-0056` Folgepflicht): `.` trennt Subjekt-Ebenen, `*`/`>` sind
// Wildcards. Ein Token mit einem dieser Zeichen würde das vier-Ebenen-
// Subjekt unbeabsichtigt aufspalten oder eine Wildcard-Bedeutung annehmen.
const reservedSubjectChars = ".*>"

// Option konfiguriert den Adapter bei der Konstruktion (`New`); aktuell
// trägt sie nur den optionalen `LogPort` (`LH-QA-OPS-004`, `ADR-0024`) —
// dasselbe Muster wie `postgresack.Option`.
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
// `postgresack.WithLog`, eigenständig definiert, weil beide Pakete keine
// gemeinsame Options-Infrastruktur teilen (kein Adapter-→-Adapter-Import).
func WithLog(log outbound.LogPort) Option {
	return func(o *options) { o.log = log }
}

// NatsChangeNotificationAdapter sendet das Wecksignal über eine bestehende
// Core-NATS-Verbindung (`ADR-0055`, Option C): kein eigener Zustand, kein
// Payload-Inhalt. `log` trägt die strukturierte Protokollierung über den
// injizierten `LogPort` (`LH-QA-OPS-004`, `ADR-0024`) — Default
// `outbound.NoopLog`.
type NatsChangeNotificationAdapter struct {
	conn *nats.Conn
	log  outbound.LogPort
}

// New legt den Notify-Adapter auf eine bestehende NATS-Verbindung; die
// Composition Root verdrahtet den Verbindungsaufbau über `CDC_NATS_URL`
// (`ADR-0055`).
func New(conn *nats.Conn, opts ...Option) (*NatsChangeNotificationAdapter, error) {
	if conn == nil {
		return nil, fmt.Errorf("%w: keine NATS-Verbindung", outbound.ErrNotify)
	}
	o := newOptions(opts)
	// Kein `context.Context`-Parameter an dieser Konstruktionsstelle
	// (Signatur-Vertrag von `New`, analog `postgresack.New`) — der
	// Log-Aufruf trägt deshalb `context.Background()`.
	o.log.Info(context.Background(), "natsnotify: verdrahtet")
	return &NatsChangeNotificationAdapter{conn: conn, log: o.log}, nil
}

var _ outbound.ChangeNotificationPort = (*NatsChangeNotificationAdapter)(nil)

// notifyFailure trägt die Übersetzungsverantwortung des Adapters
// (`ADR-0023`, `SPEC-008`): Treiber-Fehler gehen an dieser Grenze in die
// Klasse `transient` (`outbound.ErrNotify`); die technische Ursache bleibt
// über die zweite Wrappung lesbar. Derselbe Aufruf trägt den
// strukturierten Fehler-Log über den injizierten `LogPort`
// (`LH-QA-OPS-004`, `ADR-0024`).
func notifyFailure(ctx context.Context, log outbound.LogPort, cause error) error {
	log.Error(ctx, "natsnotify: Signalfehler", "error", cause)
	return fmt.Errorf("%w: %v", outbound.ErrNotify, cause)
}

// containsReservedSubjectToken meldet, ob `token` ein NATS-Subjekt-
// Sonderzeichen oder Whitespace trägt (`ADR-0056` Folgepflicht) — die
// defensive Validierung vor dem ersten Publish-Versuch.
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

// Notify sendet das Wecksignal auf
// `cdc.changes.<source_id>.<schema>.<table>` mit leerem Payload
// (`SPEC-017`, `ADR-0056`): kein Change-Inhalt, keine Positionsangabe —
// jede Nachricht bedeutet ausschließlich „lies erneut über den
// bestehenden Zugriffsweg". Core NATS trägt keine Zustellgarantie
// (Fire-and-Forget); die Rückkehr ohne Fehler meldet den abgeschickten
// Publish-Versuch. Schema/Tabelle mit NATS-reservierten Zeichen
// (`.`, `*`, `>`) oder Whitespace werden vor dem Publish-Versuch
// abgelehnt (`ADR-0056` Folgepflicht) — sonst spaltet ein Punkt im Token
// das Subjekt unbeabsichtigt in weitere Ebenen auf.
func (a *NatsChangeNotificationAdapter) Notify(ctx context.Context, sourceID, schema, table string) error {
	if sourceID == "" {
		return fmt.Errorf("%w: Notify ohne Quelle", outbound.ErrNotify)
	}
	if schema == "" || table == "" {
		return fmt.Errorf("%w: Notify ohne Schema oder Tabelle", outbound.ErrNotify)
	}
	if containsReservedSubjectToken(schema) || containsReservedSubjectToken(table) {
		return fmt.Errorf("%w: Schema %q oder Tabelle %q trägt ein NATS-reserviertes Zeichen (.,*,>) oder Whitespace", outbound.ErrNotify, schema, table)
	}
	subject := subjectPrefix + sourceID + "." + schema + "." + table
	if err := a.conn.Publish(subject, nil); err != nil {
		return notifyFailure(ctx, a.log, err)
	}
	a.log.Debug(ctx, "natsnotify: Signal gesendet", "source_id", sourceID, "schema", schema, "table", table, "subject", subject)
	return nil
}
