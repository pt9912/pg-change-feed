package receive

import (
	"context"
	"fmt"

	"github.com/jackc/pglogrepl"
	"github.com/jackc/pgx/v5/pgconn"
)

// WALRetentionChecker misst den WAL-Rückstand eines Logical-Replication-
// Slots (`SPEC-009` `cdc_wal_retention_bytes`, `ADR-0049`) über eine
// eigene, von der Stream-Verbindung getrennte Verbindung derselben Rolle
// (`cdc_capture`): Die Stream-Verbindung (`Stream.conn`) befindet sich
// während `Stream.Run` durchgehend im COPY-Modus des Replication-
// Protokolls und nimmt dort keine weiteren Abfragen entgegen — die
// periodische Messung braucht deshalb eine zweite, unabhängige
// Verbindung mit demselben `replication=database`-Verbindungsparameter.
// Anders als die Stream-Verbindung steht diese zwischen zwei Messungen
// untätig (kein Keepalive-Verkehr) und kann dadurch von der Quelle oder
// einem dazwischenliegenden Netzwerkelement getrennt werden (Idle-Timeout)
// — `dsn` bleibt deshalb erhalten, damit `Measure` eine so gestörte
// Verbindung selbst ersetzen kann.
type WALRetentionChecker struct {
	dsn  string
	conn *pgconn.PgConn
	slot string
}

// NewWALRetentionChecker baut die eigene Verbindung auf (dieselbe
// Rollen-DSN wie der Stream-Adapter, `connectReplication`); der
// Slot-Name trägt dasselbe Bezeichner-Alphabet wie die
// Stream-Konfiguration (`identifierShape`). Der Slot selbst muss bereits
// bestehen — `NewStream` legt ihn vor dem ersten Aufruf von `Measure` an.
func NewWALRetentionChecker(ctx context.Context, dsn, slot string) (*WALRetentionChecker, error) {
	if !identifierShape.MatchString(slot) {
		return nil, fmt.Errorf("%w: Slot-Name %q trägt nicht das Bezeichner-Alphabet", ErrConfiguration, slot)
	}
	conn, err := connectReplication(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return &WALRetentionChecker{dsn: dsn, conn: conn, slot: slot}, nil
}

// Close schließt die eigene Verbindung.
func (c *WALRetentionChecker) Close(ctx context.Context) error {
	return c.conn.Close(ctx)
}

// Measure liefert den WAL-Rückstand des Slots in Bytes: die Differenz
// zwischen der aktuellen WAL-Schreibposition der Quelle (`IDENTIFY_SYSTEM`)
// und `confirmed_flush_lsn` des Slots (`pg_replication_slots`). Beide
// LSN-Werte tragen denselben durchgehenden 64-Bit-Byte-Offset
// (`pglogrepl.LSN`) — ihre Differenz ist dieselbe Bytezahl, die
// PostgreSQLs `pg_wal_lsn_diff` berechnet, hier ohne den sonst nötigen
// Aufruf einer auf die Rolle `pg_monitor` beschränkten Systemfunktion
// (`pg_current_wal_lsn`): `IDENTIFY_SYSTEM` ist Teil des
// Replication-Protokolls und braucht nur das `REPLICATION`-Attribut, das
// `cdc_capture` bereits trägt.
func (c *WALRetentionChecker) Measure(ctx context.Context) (int64, error) {
	current, err := pglogrepl.IdentifySystem(ctx, c.conn)
	if err != nil {
		c.reconnectAfterError(ctx)
		return 0, fmt.Errorf("%w: IDENTIFY_SYSTEM: %v", ErrReplication, err)
	}
	values, exists, err := querySingle(ctx, c.conn,
		"SELECT confirmed_flush_lsn::text FROM pg_replication_slots WHERE slot_name = '"+c.slot+"' AND slot_type = 'logical'")
	if err != nil {
		c.reconnectAfterError(ctx)
		return 0, err
	}
	if !exists || values[0] == "" {
		return 0, fmt.Errorf("%w: Slot %q trägt keine confirmed_flush_lsn", ErrReplication, c.slot)
	}
	confirmed, err := pglogrepl.ParseLSN(values[0])
	if err != nil {
		return 0, fmt.Errorf("%w: confirmed_flush_lsn %q: %v", ErrReplication, values[0], err)
	}
	return int64(current.XLogPos - confirmed), nil
}

// reconnectAfterError ersetzt die eigene Verbindung, nachdem ein
// Protokoll- oder Katalogaufruf auf ihr fehlgeschlagen ist: Eine zwischen
// zwei Ticks dauerhaft unterbrochene Verbindung (Idle-Timeout, s. o.)
// bleibt so nicht für den Rest des Prozesslaufs bestehen — der nächste
// `Measure`-Aufruf versucht erneut, statt dass die Metrik endgültig
// verstummt. Ein Fehler beim Schließen der alten Verbindung wird
// verworfen — sie gilt bereits als unbrauchbar. Schlägt der Neuaufbau
// selbst fehl, bleibt die alte (bereits geschlossene) Verbindung stehen;
// der nächste `Measure`-Aufruf scheitert dann erneut auf demselben Pfad
// und versucht wieder — kein Aufgeben, aber auch kein
// Endlos-Reconnect-Loop innerhalb eines einzigen Aufrufs.
func (c *WALRetentionChecker) reconnectAfterError(ctx context.Context) {
	_ = c.conn.Close(ctx)
	if conn, err := connectReplication(ctx, c.dsn); err == nil {
		c.conn = conn
	}
}
