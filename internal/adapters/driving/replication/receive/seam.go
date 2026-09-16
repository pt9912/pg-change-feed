package receive

import (
	"context"

	"github.com/jackc/pglogrepl"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgproto3"
)

// driverSession trägt die Fläche, die der Stream-Adapter vom Treiber
// braucht: die sieben Operationen, die er wirklich ausführt — die vier
// `pglogrepl`-Paketfunktionen, den Katalog-Abruf, den Meldungs-Empfang
// und das Schließen. Die Signaturen sind die der realen Aufrufe: die
// vier Paketfunktionen verlangen den konkreten `*pgconn.PgConn` in der
// Signatur und sind damit selbst nicht von einem Fake erfüllbar, und
// `pgconn.PgConn.Exec` liefert den nicht fake-baren
// `*pgconn.MultiResultReader` — die Hülle erfüllt die Fläche mit dem Typ
// in der Hand, ein Fake mit den exportierten Werttypen von
// `pgconn`/`pglogrepl`. Was oberhalb der Naht liegt, kennt den konkreten
// Typ nicht mehr (`ADR-0080`).
type driverSession interface {
	IdentifySystem(ctx context.Context) (pglogrepl.IdentifySystemResult, error)
	CreateReplicationSlot(ctx context.Context, slot, outputPlugin string, options pglogrepl.CreateReplicationSlotOptions) (pglogrepl.CreateReplicationSlotResult, error)
	StartReplication(ctx context.Context, slot string, startLSN pglogrepl.LSN, options pglogrepl.StartReplicationOptions) error
	SendStandbyStatusUpdate(ctx context.Context, ssu pglogrepl.StandbyStatusUpdate) error
	Exec(ctx context.Context, sql string) ([]*pgconn.Result, error)
	ReceiveMessage(ctx context.Context) (pgproto3.BackendMessage, error)
	Close(ctx context.Context) error
}

// connSession ist die Treiber-Hülle: sie hält die konkrete
// Replication-Verbindung und delegiert 1:1 an `pgconn`/`pglogrepl` —
// dieselbe Prüfung der Bibliothek (etwa die LSN-Form) läuft damit
// weiter.
type connSession struct {
	conn *pgconn.PgConn
}

func (c connSession) IdentifySystem(ctx context.Context) (pglogrepl.IdentifySystemResult, error) {
	return pglogrepl.IdentifySystem(ctx, c.conn)
}

func (c connSession) CreateReplicationSlot(ctx context.Context, slot, outputPlugin string, options pglogrepl.CreateReplicationSlotOptions) (pglogrepl.CreateReplicationSlotResult, error) {
	return pglogrepl.CreateReplicationSlot(ctx, c.conn, slot, outputPlugin, options)
}

func (c connSession) StartReplication(ctx context.Context, slot string, startLSN pglogrepl.LSN, options pglogrepl.StartReplicationOptions) error {
	return pglogrepl.StartReplication(ctx, c.conn, slot, startLSN, options)
}

func (c connSession) SendStandbyStatusUpdate(ctx context.Context, ssu pglogrepl.StandbyStatusUpdate) error {
	return pglogrepl.SendStandbyStatusUpdate(ctx, c.conn, ssu)
}

func (c connSession) Exec(ctx context.Context, sql string) ([]*pgconn.Result, error) {
	return c.conn.Exec(ctx, sql).ReadAll()
}

func (c connSession) ReceiveMessage(ctx context.Context) (pgproto3.BackendMessage, error) {
	return c.conn.ReceiveMessage(ctx)
}

func (c connSession) Close(ctx context.Context) error {
	return c.conn.Close(ctx)
}

// Die Hülle erfüllt die Naht. Die Zusicherung ist der Kompilier-Beleg,
// dass die delegierenden Aufrufe die Signaturen der Bibliothek treffen:
// ändert `pglogrepl` eine von ihnen, bricht dieser Bau — nicht erst der
// Aufruf einer Adapter-Methode.
var _ driverSession = connSession{}
