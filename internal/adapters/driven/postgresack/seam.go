package postgresack

import (
	"context"

	"github.com/jackc/pglogrepl"
	"github.com/jackc/pgx/v5/pgconn"
)

// standbySender trägt die Fläche, die der ACK-Adapter vom Treiber
// braucht: den Standby-Status-Update absetzen. Die Signatur ist die des
// realen Aufrufs — `pglogrepl.SendStandbyStatusUpdate` verlangt den
// konkreten `*pgconn.PgConn` in der Signatur und ist damit selbst nicht
// von einem Fake erfüllbar; die Hülle erfüllt sie mit dem Typ in der
// Hand. Was oberhalb der Naht liegt, kennt den konkreten Typ nicht mehr
// (`ADR-0080`).
type standbySender interface {
	SendStandbyStatusUpdate(ctx context.Context, ssu pglogrepl.StandbyStatusUpdate) error
}

// connSender ist die Treiber-Hülle: sie hält die konkrete
// Replication-Verbindung und delegiert 1:1 an `pglogrepl` — dieselbe
// Prüfung der Bibliothek (etwa die LSN-Form) läuft damit weiter.
type connSender struct {
	conn *pgconn.PgConn
}

func (c connSender) SendStandbyStatusUpdate(ctx context.Context, ssu pglogrepl.StandbyStatusUpdate) error {
	return pglogrepl.SendStandbyStatusUpdate(ctx, c.conn, ssu)
}

// Die Hülle erfüllt die Naht. Die Zusicherung ist der Kompilier-Beleg,
// dass der delegierende Aufruf die Signatur der Bibliothek trifft: ändert
// `pglogrepl` sie, bricht dieser Bau — nicht erst der Aufruf einer
// Adapter-Methode.
var _ standbySender = connSender{}
