// Package postgressnapshot trägt den `PostgresTableSnapshotAdapter` als
// Driven-Implementierung des `TableSnapshotPort` (`ARC-006`,
// `LH-FA-CAP-009`, `ADR-0111` Teilfrage 1): er liest den Bestand einer
// Tabelle in dem `REPEATABLE READ`-Snapshot, den ein je Run angelegter
// temporärer logischer Slot exportiert. Snapshot und Position `X` (der
// `consistent_point` des Slots) sind dadurch exakt gepaart.
//
// Alle Verbindungen laufen über denselben DSN (`CDC_CAPTURE_DSN`): die
// Replication-Verbindung legt den Slot an und endet nach dem Import, eine
// reguläre Verbindung liest im importierten Snapshot. Der Adapter liegt in
// einem eigenen Paket, weil er die Quelltabelle über eine
// Replication-Verbindung liest und der Run-Store (`postgresstorage`) eine
// andere Verantwortung trägt. Sein Testlauf setzt einen PostgreSQL mit
// `wal_level=logical` voraus (`make test-replication`); netzlos prüfbare
// Logik liegt außerhalb des Pakets (Unterpaket `snapshotlogic`,
// `model.BuildRowImage`), damit der Gegenstand der DB-Adapter-Coverage nur
// real Gedecktes zählt (`ADR-0071` Punkt 1).
package postgressnapshot

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pglogrepl"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgressnapshot/snapshotlogic"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
)

const (
	// DefaultBlockSize ist die Blockgröße `B` (Zeilen je `NextBlock`):
	// die Zeilenzahl des Lesens ist durch sie begrenzt, nicht die Bytezahl.
	// Der Wert ist ein Startwert (Setzung ohne Messung, `ADR-0111`
	// Teilfrage 1).
	DefaultBlockSize = 1000

	// DefaultSlotTimeout begrenzt Verbindungsaufbau und Slot-Anlage. Die
	// Anlage wartet auf die beim Aufruf laufenden Schreibtransaktionen der
	// Quelle; ihr Ablauf endet als Fehlerklasse `transient`. Der Wert ist
	// ein Startwert (Setzung ohne Messung).
	DefaultSlotTimeout = 30 * time.Second

	outputPlugin = "pgoutput"
	closeTimeout = 5 * time.Second
)

// Option konfiguriert den Adapter bei der Konstruktion.
type Option func(*PostgresTableSnapshotAdapter)

// WithBlockSize setzt die Blockgröße `B` (mindestens 1).
func WithBlockSize(rows int) Option {
	return func(a *PostgresTableSnapshotAdapter) { a.blockSize = rows }
}

// WithSlotTimeout setzt das Zeitlimit von Verbindungsaufbau und
// Slot-Anlage (größer als 0).
func WithSlotTimeout(timeout time.Duration) Option {
	return func(a *PostgresTableSnapshotAdapter) { a.slotTimeout = timeout }
}

// PostgresTableSnapshotAdapter liest Tabellenbestände im Slot-Snapshot.
type PostgresTableSnapshotAdapter struct {
	config      *pgconn.Config
	blockSize   int
	slotTimeout time.Duration
}

var _ outbound.TableSnapshotPort = (*PostgresTableSnapshotAdapter)(nil)

// New legt den Adapter auf den DSN der Capture-Rolle. Ein leerer oder
// nicht parsbarer DSN, eine Blockgröße unter 1 und ein Zeitlimit ≤ 0
// enden als Fehlerklasse `configuration`, ohne Verbindung.
func New(dsn string, opts ...Option) (*PostgresTableSnapshotAdapter, error) {
	a := &PostgresTableSnapshotAdapter{blockSize: DefaultBlockSize, slotTimeout: DefaultSlotTimeout}
	for _, opt := range opts {
		opt(a)
	}
	config, err := snapshotlogic.Validate(dsn, a.blockSize, a.slotTimeout)
	if err != nil {
		return nil, err
	}
	a.config = config
	return a, nil
}

// connConfig liefert eine Kopie der Verbindungskonfiguration: mit dem
// Replication-Modus der Datenbank für die Slot-Anlage, sonst als reguläre
// Verbindung zum Lesen.
func (a *PostgresTableSnapshotAdapter) connConfig(replication bool) *pgconn.Config {
	config := a.config.Copy()
	if replication {
		config.RuntimeParams["replication"] = "database"
	} else {
		delete(config.RuntimeParams, "replication")
	}
	return config
}

// slotExport ist der Zwischenstand nach der Slot-Anlage: die
// Replication-Verbindung hält den exportierten Snapshot, bis er importiert
// ist.
type slotExport struct {
	conn         *pgconn.PgConn
	offset       uint64
	snapshotName string
}

// OpenSnapshot legt den temporären Slot an, importiert seinen Snapshot in
// einer regulären `REPEATABLE READ`-Transaktion, beendet die
// Replication-Verbindung und liest Spalten und Cursor im Snapshot.
func (a *PostgresTableSnapshotAdapter) OpenSnapshot(ctx context.Context, runID, schema, table string) (outbound.TableSnapshot, error) {
	slot, err := snapshotlogic.SlotName(runID)
	if err != nil {
		return nil, err
	}
	if schema == "" || table == "" {
		return nil, fmt.Errorf("%w: Tabelle ohne Namen", outbound.ErrSnapshotConfiguration)
	}
	export, err := a.exportSnapshot(ctx, slot)
	if err != nil {
		return nil, err
	}
	return a.importSnapshot(ctx, export, schema, table)
}

// exportSnapshot verbindet über die Replication-Verbindung und legt den
// temporären Slot mit Snapshot-Export an. Verbindungsaufbau und Anlage
// tragen das Zeitlimit; die Anlage wartet auf laufende
// Schreibtransaktionen der Quelle.
func (a *PostgresTableSnapshotAdapter) exportSnapshot(ctx context.Context, slot string) (*slotExport, error) {
	slotCtx, cancel := context.WithTimeout(ctx, a.slotTimeout)
	defer cancel()

	conn, err := pgconn.ConnectConfig(slotCtx, a.connConfig(true))
	if err != nil {
		return nil, snapshotlogic.Classify(slotCtx, err, outbound.ErrSnapshotReplication, "Replication-Verbindung")
	}
	result, err := pglogrepl.CreateReplicationSlot(slotCtx, conn, slot, outputPlugin, pglogrepl.CreateReplicationSlotOptions{
		Temporary:      true,
		SnapshotAction: "EXPORT_SNAPSHOT",
		Mode:           pglogrepl.LogicalReplication,
	})
	if err != nil {
		closeConn(conn)
		return nil, snapshotlogic.Classify(slotCtx, err, outbound.ErrSnapshotReplication, "Slot-Anlage")
	}
	lsn, err := pglogrepl.ParseLSN(result.ConsistentPoint)
	if err != nil || lsn == 0 {
		closeConn(conn)
		return nil, fmt.Errorf("%w: consistent_point %q nicht lesbar", outbound.ErrSnapshotReplication, result.ConsistentPoint)
	}
	if !snapshotlogic.ValidSnapshotName(result.SnapshotName) {
		closeConn(conn)
		return nil, fmt.Errorf("%w: Snapshot-Name %q nicht lesbar", outbound.ErrSnapshotReplication, result.SnapshotName)
	}
	return &slotExport{conn: conn, offset: uint64(lsn), snapshotName: result.SnapshotName}, nil
}

// importSnapshot importiert den exportierten Snapshot in einer regulären
// Verbindung, beendet danach die Replication-Verbindung (der temporäre Slot
// fällt weg, der Snapshot bleibt lesbar) und bereitet Spalten und Cursor
// vor. Nach der Rückkehr — mit oder ohne Fehler — ist die
// Replication-Verbindung geschlossen.
func (a *PostgresTableSnapshotAdapter) importSnapshot(ctx context.Context, export *slotExport, schema, table string) (outbound.TableSnapshot, error) {
	defer closeConn(export.conn)

	conn, err := pgconn.ConnectConfig(ctx, a.connConfig(false))
	if err != nil {
		return nil, snapshotlogic.Classify(ctx, err, outbound.ErrSnapshotStorage, "Verbindung zum Lesen")
	}
	snap := &snapshot{conn: conn, offset: export.offset, blockSize: a.blockSize}
	for _, statement := range []string{
		"BEGIN ISOLATION LEVEL REPEATABLE READ READ ONLY",
		"SET TRANSACTION SNAPSHOT '" + export.snapshotName + "'",
	} {
		if err := snap.exec(ctx, statement); err != nil {
			snap.abort()
			return nil, snapshotlogic.Classify(ctx, err, outbound.ErrSnapshotStorage, "Snapshot-Import")
		}
	}
	// Der Snapshot ist importiert: die Replication-Verbindung endet, der
	// Slot fällt weg, die Lese-Transaktion behält den Stand.
	closeConn(export.conn)

	columns, err := readColumns(ctx, conn, schema, table)
	if err != nil {
		snap.abort()
		return nil, err
	}
	snap.columns = columns
	if err := snap.exec(ctx, snapshotlogic.CursorStatement(columns, schema, table)); err != nil {
		snap.abort()
		return nil, snapshotlogic.Classify(ctx, err, outbound.ErrSnapshotStorage, "Cursor")
	}
	return snap, nil
}

// readColumns liest die Spaltenliste im Snapshot: Spalten der Tabelle in
// Tabellenordnung, ohne gelöschte und ohne generierte Spalten. Eine nicht
// vorhandene Tabelle und eine Tabelle ohne lesbare Spalte enden als
// Fehlerklasse `configuration`.
func readColumns(ctx context.Context, conn *pgconn.PgConn, schema, table string) ([]string, error) {
	result := conn.ExecParams(ctx,
		`SELECT a.attname
		   FROM pg_attribute a
		  WHERE a.attrelid = to_regclass(quote_ident($1) || '.' || quote_ident($2))
		    AND a.attnum > 0
		    AND NOT a.attisdropped
		    AND a.attgenerated = ''
		  ORDER BY a.attnum`,
		[][]byte{[]byte(schema), []byte(table)}, nil, nil, nil).Read()
	if result.Err != nil {
		return nil, snapshotlogic.Classify(ctx, result.Err, outbound.ErrSnapshotStorage, "Spaltenliste")
	}
	if len(result.Rows) == 0 {
		return nil, fmt.Errorf("%w: Tabelle %s.%s nicht vorhanden oder ohne lesbare Spalte", outbound.ErrSnapshotConfiguration, schema, table)
	}
	columns := make([]string, len(result.Rows))
	for i, row := range result.Rows {
		columns[i] = string(row[0])
	}
	return columns, nil
}

// EstimatedRows liest `pg_class.reltuples` der Tabelle über eine kurze
// eigene Verbindung; ein negativer Wert meldet „unbekannt".
func (a *PostgresTableSnapshotAdapter) EstimatedRows(ctx context.Context, schema, table string) (int64, bool, error) {
	if schema == "" || table == "" {
		return 0, false, fmt.Errorf("%w: Tabelle ohne Namen", outbound.ErrSnapshotConfiguration)
	}
	conn, err := pgconn.ConnectConfig(ctx, a.connConfig(false))
	if err != nil {
		return 0, false, snapshotlogic.Classify(ctx, err, outbound.ErrSnapshotStorage, "Verbindung zum Lesen")
	}
	defer closeConn(conn)

	result := conn.ExecParams(ctx,
		`SELECT c.reltuples::float8::text
		   FROM pg_class c
		  WHERE c.oid = to_regclass(quote_ident($1) || '.' || quote_ident($2))`,
		[][]byte{[]byte(schema), []byte(table)}, nil, nil, nil).Read()
	if result.Err != nil {
		return 0, false, snapshotlogic.Classify(ctx, result.Err, outbound.ErrSnapshotStorage, "Zeilenschätzung")
	}
	if len(result.Rows) == 0 {
		return 0, false, fmt.Errorf("%w: Tabelle %s.%s nicht vorhanden", outbound.ErrSnapshotConfiguration, schema, table)
	}
	return snapshotlogic.Estimate(string(result.Rows[0][0]))
}

// snapshot ist die Lese-Transaktion im importierten Snapshot.
type snapshot struct {
	conn      *pgconn.PgConn
	offset    uint64
	columns   []string
	blockSize int
	closed    bool
}

var _ outbound.TableSnapshot = (*snapshot)(nil)

func (s *snapshot) Offset() uint64 { return s.offset }

func (s *snapshot) Columns() []string { return s.columns }

// exec führt ein Statement ohne Ergebniszeilen aus.
func (s *snapshot) exec(ctx context.Context, statement string) error {
	_, err := s.conn.Exec(ctx, statement).ReadAll()
	return err
}

// abort beendet die Verbindung eines nicht fertig geöffneten Snapshots.
func (s *snapshot) abort() {
	closeConn(s.conn)
	s.closed = true
}

// NextBlock holt die nächsten höchstens `blockSize` Zeilen des Cursors. Der
// einfache Query-Pfad von `pgconn` liefert jedes Feld im Text-Ergebnisformat
// — Protokoll-Eigenschaft, keine Einstellung (`ADR-0115` Festlegung 2); die
// Rohbytes gehen als Wert weiter.
func (s *snapshot) NextBlock(ctx context.Context) ([][]*string, error) {
	if s.closed {
		return nil, fmt.Errorf("%w: Snapshot geschlossen", outbound.ErrSnapshotStorage)
	}
	results, err := s.conn.Exec(ctx, snapshotlogic.FetchStatement(s.blockSize)).ReadAll()
	if err != nil {
		// Eine beendete Verbindung ist vorübergehend, ein Fehler auf
		// lebender Verbindung ein Lesefehler.
		fallback := outbound.ErrSnapshotStorage
		if s.conn.IsClosed() {
			fallback = outbound.ErrSnapshotTransient
		}
		return nil, snapshotlogic.Classify(ctx, err, fallback, "Cursor lesen")
	}
	block := make([][]*string, 0, s.blockSize)
	for _, result := range results {
		for _, row := range result.Rows {
			block = append(block, snapshotlogic.Values(row))
		}
	}
	return block, nil
}

// Close beendet die Lese-Transaktion und die Verbindung.
func (s *snapshot) Close(context.Context) error {
	if s.closed {
		return nil
	}
	s.abort()
	return nil
}

// closeConn beendet eine Verbindung mit eigenem Zeitlimit, unabhängig vom
// Kontext des Aufrufers: ein abgebrochener Aufrufer lässt keine Sitzung
// und keinen Slot zurück.
func closeConn(conn *pgconn.PgConn) {
	ctx, cancel := context.WithTimeout(context.Background(), closeTimeout)
	defer cancel()
	_ = conn.Close(ctx)
}
