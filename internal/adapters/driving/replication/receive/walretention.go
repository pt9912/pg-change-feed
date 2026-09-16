package receive

import (
	"context"
	"fmt"
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
	dsn string
	// session ist dieselbe Naht wie die des Stream-Adapters (`ADR-0080`):
	// die eigene Verbindung hinter der Treiber-Hülle, damit die
	// Rückstands-Messung netzlos fahrbar ist.
	session driverSession
	slot    string
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
	return newWALRetentionCheckerOnSession(connSession{conn: conn}, dsn, slot), nil
}

// newWALRetentionCheckerOnSession verdrahtet den Checker auf eine Naht —
// der paket-interne Einstieg der netzlosen Tests, die einen Fake anstelle
// des Treibers fahren; `NewWALRetentionChecker` reicht die Treiber-Hülle
// durch. `dsn` bleibt der Neuaufbau-Träger von `reconnectAfterError`.
func newWALRetentionCheckerOnSession(session driverSession, dsn, slot string) *WALRetentionChecker {
	return &WALRetentionChecker{dsn: dsn, session: session, slot: slot}
}

// Close schließt die eigene Verbindung.
func (c *WALRetentionChecker) Close(ctx context.Context) error {
	return c.session.Close(ctx)
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
	current, err := c.session.IdentifySystem(ctx)
	if err != nil {
		c.reconnectAfterError(ctx)
		return 0, fmt.Errorf("%w: IDENTIFY_SYSTEM: %v", ErrReplication, err)
	}
	values, exists, err := querySingle(ctx, c.session, slotLSNQuery(c.slot))
	if err != nil {
		c.reconnectAfterError(ctx)
		return 0, err
	}
	if !exists || values[0] == "" {
		return 0, fmt.Errorf("%w: Slot %q trägt keine confirmed_flush_lsn", ErrReplication, c.slot)
	}
	confirmed, err := parseLSN("confirmed_flush_lsn", values[0])
	if err != nil {
		return 0, err
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
	_ = c.session.Close(ctx)
	if conn, err := connectReplication(ctx, c.dsn); err == nil {
		c.session = connSession{conn: conn}
	}
}
