// Package postgresack trägt den `PostgresReplicationAckAdapter` als
// Driven-Implementierung des `ReplicationAckPort` (`ARC-006`,
// `ADR-0007`): er bestätigt die vom Capture Use Case gemeldete Position
// gegenüber der Quelle — erst nach bestätigter Persistenz
// (`LH-QA-REL-001.a`); die Ordnung trägt die Application (`ADR-0027`).
// Der Treiber (pglogrepl über den pgx-Verbindungstyp) bleibt
// Adapterdetail (`ADR-0032`).
package postgresack

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pglogrepl"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// PostgresReplicationAckAdapter bestätigt Positionen gegenüber dem
// Quell-PostgreSQL über die Replication-Verbindung des Streams: Stream
// und ACK sind getrennte Rollen an einer technischen Verbindung
// (`ADR-0007`, Option C).
type PostgresReplicationAckAdapter struct {
	conn *pgconn.PgConn
}

// New legt den ACK-Adapter auf die Replication-Verbindung des Streams;
// die Composition Root verdrahtet beide Adapter über
// `receive.Stream.Conn` (`ADR-0026`).
func New(conn *pgconn.PgConn) (*PostgresReplicationAckAdapter, error) {
	if conn == nil {
		return nil, fmt.Errorf("%w: keine Replication-Verbindung", outbound.ErrReplication)
	}
	slog.Info("replicationack: verdrahtet")
	return &PostgresReplicationAckAdapter{conn: conn}, nil
}

var _ outbound.ReplicationAckPort = (*PostgresReplicationAckAdapter)(nil)

// replicationFailure trägt die Übersetzungsverantwortung des Adapters
// (`ADR-0023`, `SPEC-008`): Treiber-Fehler gehen an dieser Grenze in die
// Klasse `replication` (`outbound.ErrReplication`); die technische
// Ursache bleibt über die zweite Wrappung lesbar. Derselbe Aufruf trägt
// den strukturierten Fehler-Log (`LH-QA-OPS-004`); `New` oben trägt kein
// `context.Context`, deshalb `slog.Info` statt `InfoContext`.
func replicationFailure(cause error) error {
	slog.Error("replicationack: Bestätigungsfehler", "error", cause)
	return fmt.Errorf("%w: %v", outbound.ErrReplication, cause)
}

// Acknowledge bestätigt die Position gegenüber der Quelle: der
// Standby-Status-Update trägt die Position als Write-, Flush- und
// Apply-Stand, damit der Slot sie als confirmed_flush_lsn trägt. Die
// Bestätigung gilt mit der Rückkehr ohne Fehler (`LH-QA-REL-001.a`,
// Schritt ACK Source); ein Fehler endet ohne Bestätigung.
func (a *PostgresReplicationAckAdapter) Acknowledge(ctx context.Context, position model.SourcePosition) error {
	if position.IsZero() {
		return fmt.Errorf("%w: Bestätigung ohne Position", outbound.ErrReplication)
	}
	lsn := pglogrepl.LSN(position.Offset)
	if err := pglogrepl.SendStandbyStatusUpdate(ctx, a.conn, pglogrepl.StandbyStatusUpdate{
		WALWritePosition: lsn,
		WALFlushPosition: lsn,
		WALApplyPosition: lsn,
	}); err != nil {
		return replicationFailure(err)
	}
	slog.DebugContext(ctx, "replicationack: Position bestätigt", "offset", position.Offset)
	return nil
}
