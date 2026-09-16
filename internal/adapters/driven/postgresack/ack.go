// Package postgresack trägt den `PostgresReplicationAckAdapter` als
// Driven-Implementierung des `ReplicationAckPort` (`ARC-006`,
// `ADR-0007`): er bestätigt die vom Capture Use Case gemeldete Position
// gegenüber der Quelle — erst nach bestätigter Persistenz
// (`LH-QA-REL-001.a`); die Ordnung trägt die Application (`ADR-0027`).
// Der Treiber (pglogrepl über den pgx-Verbindungstyp) bleibt
// Adapterdetail (`ADR-0032`) und liegt hinter der Naht (`seam.go`,
// `ADR-0080`).
package postgresack

import (
	"context"
	"fmt"

	"github.com/jackc/pglogrepl"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Option konfiguriert den Adapter bei der Konstruktion (`New`); aktuell
// trägt sie nur den optionalen `LogPort` (`LH-QA-OPS-004`, `ADR-0024`) —
// variadisch, damit bestehende Aufrufstellen (Tests) unverändert
// kompilieren. Ungesetzt bleibt die Protokollierung beim No-Op
// (`outbound.NoopLog`).
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
// `postgresstorage.WithLog`, eigenständig definiert, weil beide Pakete
// keine gemeinsame Options-Infrastruktur teilen (kein Adapter-→-Adapter-
// Import).
func WithLog(log outbound.LogPort) Option {
	return func(o *options) { o.log = log }
}

// PostgresReplicationAckAdapter bestätigt Positionen gegenüber dem
// Quell-PostgreSQL über die Replication-Verbindung des Streams: Stream
// und ACK sind getrennte Rollen an einer technischen Verbindung
// (`ADR-0007`, Option C). `sender` ist die Naht zum Treiber (`ADR-0080`),
// `log` trägt die strukturierte Protokollierung über den injizierten
// `LogPort` (`LH-QA-OPS-004`, `ADR-0024`) — Default `outbound.NoopLog`.
type PostgresReplicationAckAdapter struct {
	sender standbySender
	log    outbound.LogPort
}

// errNoConnection trägt die Konstruktions-Grenze des Adapters: ohne
// Träger gibt es keinen ACK-Stand (`SPEC-008`, Klasse `replication`).
var errNoConnection = fmt.Errorf("%w: keine Replication-Verbindung", outbound.ErrReplication)

// New legt den ACK-Adapter auf die Replication-Verbindung des Streams;
// die Composition Root verdrahtet beide Adapter über
// `receive.Stream.Conn` (`ADR-0026`). Die Verbindung geht in die
// Treiber-Hülle (`seam.go`), damit die Logik an der Naht hängt
// (`ADR-0080`).
func New(conn *pgconn.PgConn, opts ...Option) (*PostgresReplicationAckAdapter, error) {
	if conn == nil {
		return nil, errNoConnection
	}
	return newOnSender(connSender{conn: conn}, opts...)
}

// newOnSender verdrahtet den Adapter auf eine Naht — der paket-interne
// Einstieg der netzlosen Tests, die einen Fake anstelle des Treibers
// fahren; `New` reicht die Treiber-Hülle durch.
func newOnSender(sender standbySender, opts ...Option) (*PostgresReplicationAckAdapter, error) {
	if sender == nil {
		return nil, errNoConnection
	}
	o := newOptions(opts)
	// Kein `context.Context`-Parameter an dieser Konstruktionsstelle
	// (Signatur-Vertrag von `New`, unverändert ggü. dem Bestand) — der
	// Log-Aufruf trägt deshalb `context.Background()`, wie an jeder
	// anderen ctx-losen Konstruktionsstelle dieses Adapters
	// (`postgresstorage.New`/`NewTableActivation`/`NewHeartbeat`/
	// `NewConsumerState` haben dagegen `ctx` als Parameter und nutzen
	// ihn).
	o.log.Info(context.Background(), "replicationack: verdrahtet")
	return &PostgresReplicationAckAdapter{sender: sender, log: o.log}, nil
}

var _ outbound.ReplicationAckPort = (*PostgresReplicationAckAdapter)(nil)

// ackLSN trägt die Null-Positions-Grenze und die LSN-Form (`ADR-0005`):
// ohne Position gibt es nichts zu bestätigen (`SPEC-008`, Klasse
// `replication`); sonst wandert der Offset als LSN weiter.
func ackLSN(position model.SourcePosition) (pglogrepl.LSN, error) {
	if position.IsZero() {
		return 0, fmt.Errorf("%w: Bestätigung ohne Position", outbound.ErrReplication)
	}
	return pglogrepl.LSN(position.Offset), nil
}

// standbyStatus bildet die Standby-Status-Meldung aus der LSN: die drei
// Positionen tragen denselben Stand, damit der Slot ihn als
// confirmed_flush_lsn trägt.
func standbyStatus(lsn pglogrepl.LSN) pglogrepl.StandbyStatusUpdate {
	return pglogrepl.StandbyStatusUpdate{
		WALWritePosition: lsn,
		WALFlushPosition: lsn,
		WALApplyPosition: lsn,
	}
}

// replicationClass trägt die Übersetzungsverantwortung des Adapters
// (`ADR-0023`, `SPEC-008`): Treiber-Fehler gehen an dieser Grenze in die
// Klasse `replication` (`outbound.ErrReplication`); die technische
// Ursache bleibt über die zweite Wrappung lesbar.
func replicationClass(cause error) error {
	return fmt.Errorf("%w: %v", outbound.ErrReplication, cause)
}

// replicationFailure trägt den strukturierten Fehler-Log über den
// injizierten `LogPort` (`LH-QA-OPS-004`, `ADR-0024`) und reicht die
// Fehlerklasse durch (`replicationClass`).
func replicationFailure(ctx context.Context, log outbound.LogPort, cause error) error {
	log.Error(ctx, "replicationack: Bestätigungsfehler", "error", cause)
	return replicationClass(cause)
}

// Acknowledge bestätigt die Position gegenüber der Quelle: der
// Standby-Status-Update trägt die Position als Write-, Flush- und
// Apply-Stand, damit der Slot sie als confirmed_flush_lsn trägt. Die
// Bestätigung gilt mit der Rückkehr ohne Fehler (`LH-QA-REL-001.a`,
// Schritt ACK Source); ein Fehler endet ohne Bestätigung.
func (a *PostgresReplicationAckAdapter) Acknowledge(ctx context.Context, position model.SourcePosition) error {
	lsn, err := ackLSN(position)
	if err != nil {
		return err
	}
	if err := a.sender.SendStandbyStatusUpdate(ctx, standbyStatus(lsn)); err != nil {
		return replicationFailure(ctx, a.log, err)
	}
	a.log.Debug(ctx, "replicationack: Position bestätigt", "offset", position.Offset)
	return nil
}
