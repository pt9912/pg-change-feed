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
// Logik liegt außerhalb des Pakets (`model.BuildRowImage`), damit der
// Gegenstand der DB-Adapter-Coverage nur real Gedecktes zählt
// (`ADR-0071` Punkt 1).
package postgressnapshot

import (
	"context"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pglogrepl"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
)

const (
	// DefaultBlockSize ist die Blockgröße `B` (Zeilen je `NextBlock`):
	// der Speicherbedarf des Lesens ist durch sie begrenzt.
	DefaultBlockSize = 1000

	// DefaultSlotTimeout begrenzt Verbindungsaufbau und Slot-Anlage. Die
	// Anlage wartet auf die beim Aufruf laufenden Schreibtransaktionen der
	// Quelle; ihr Ablauf endet als Fehlerklasse `transient`. Der Wert ist
	// ein Startwert (Setzung ohne Messung).
	DefaultSlotTimeout = 30 * time.Second

	outputPlugin  = "pgoutput"
	slotPrefix    = "cdc_bf_"
	cursorName    = "cdc_bf_cursor"
	closeTimeout  = 5 * time.Second
	maxIdentifier = 63
)

// identifierShape ist das Bezeichner-Alphabet der Quelle für Slot-Namen
// (`receive.identifierShape` hält denselben Ausdruck: die Adapter-Schicht
// importiert keine Adapter-Kante).
var identifierShape = regexp.MustCompile(`^[a-z0-9_]+$`)

// snapshotNameShape begrenzt den vom Server gelieferten Snapshot-Namen auf
// das Alphabet, in dem er als Literal in `SET TRANSACTION SNAPSHOT` steht.
var snapshotNameShape = regexp.MustCompile(`^[0-9A-Fa-f-]+$`)

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
	dsn         string
	blockSize   int
	slotTimeout time.Duration
}

var _ outbound.TableSnapshotPort = (*PostgresTableSnapshotAdapter)(nil)

// New legt den Adapter auf den DSN der Capture-Rolle. Ein leerer oder
// nicht parsbarer DSN, eine Blockgröße unter 1 und ein Zeitlimit ≤ 0
// enden als Fehlerklasse `configuration`, ohne Verbindung.
func New(dsn string, opts ...Option) (*PostgresTableSnapshotAdapter, error) {
	a := &PostgresTableSnapshotAdapter{dsn: dsn, blockSize: DefaultBlockSize, slotTimeout: DefaultSlotTimeout}
	for _, opt := range opts {
		opt(a)
	}
	if dsn == "" {
		return nil, fmt.Errorf("%w: DSN fehlt", outbound.ErrSnapshotConfiguration)
	}
	if _, err := pgconn.ParseConfig(dsn); err != nil {
		return nil, fmt.Errorf("%w: DSN: %v", outbound.ErrSnapshotConfiguration, err)
	}
	if a.blockSize < 1 {
		return nil, fmt.Errorf("%w: Blockgröße %d unter 1", outbound.ErrSnapshotConfiguration, a.blockSize)
	}
	if a.slotTimeout <= 0 {
		return nil, fmt.Errorf("%w: Zeitlimit der Slot-Anlage nicht positiv", outbound.ErrSnapshotConfiguration)
	}
	return a, nil
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
	slot, err := slotName(runID)
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

// slotName bildet den Slot-Namen `cdc_bf_<run>`: die Run-Kennung geht ohne
// Bindestriche ein und trägt danach das Bezeichner-Alphabet.
func slotName(runID string) (string, error) {
	name := slotPrefix + strings.ReplaceAll(runID, "-", "")
	if len(name) > maxIdentifier || !identifierShape.MatchString(name) || name == slotPrefix {
		return "", fmt.Errorf("%w: Run-Kennung %q ergibt keinen gültigen Slot-Namen", outbound.ErrSnapshotConfiguration, runID)
	}
	return name, nil
}

// exportSnapshot verbindet über die Replication-Verbindung und legt den
// temporären Slot mit Snapshot-Export an. Verbindungsaufbau und Anlage
// tragen das Zeitlimit; die Anlage wartet auf laufende
// Schreibtransaktionen der Quelle.
func (a *PostgresTableSnapshotAdapter) exportSnapshot(ctx context.Context, slot string) (*slotExport, error) {
	slotCtx, cancel := context.WithTimeout(ctx, a.slotTimeout)
	defer cancel()

	config, err := pgconn.ParseConfig(a.dsn)
	if err != nil {
		return nil, fmt.Errorf("%w: DSN: %v", outbound.ErrSnapshotConfiguration, err)
	}
	config.RuntimeParams["replication"] = "database"
	conn, err := pgconn.ConnectConfig(slotCtx, config)
	if err != nil {
		return nil, classify(slotCtx, err, outbound.ErrSnapshotReplication, "Replication-Verbindung")
	}
	result, err := pglogrepl.CreateReplicationSlot(slotCtx, conn, slot, outputPlugin, pglogrepl.CreateReplicationSlotOptions{
		Temporary:      true,
		SnapshotAction: "EXPORT_SNAPSHOT",
		Mode:           pglogrepl.LogicalReplication,
	})
	if err != nil {
		closeConn(conn)
		return nil, classify(slotCtx, err, outbound.ErrSnapshotReplication, "Slot-Anlage")
	}
	lsn, err := pglogrepl.ParseLSN(result.ConsistentPoint)
	if err != nil || lsn == 0 {
		closeConn(conn)
		return nil, fmt.Errorf("%w: consistent_point %q nicht lesbar", outbound.ErrSnapshotReplication, result.ConsistentPoint)
	}
	if !snapshotNameShape.MatchString(result.SnapshotName) {
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

	config, err := pgconn.ParseConfig(a.dsn)
	if err != nil {
		return nil, fmt.Errorf("%w: DSN: %v", outbound.ErrSnapshotConfiguration, err)
	}
	delete(config.RuntimeParams, "replication")
	conn, err := pgconn.ConnectConfig(ctx, config)
	if err != nil {
		return nil, classify(ctx, err, outbound.ErrSnapshotStorage, "Verbindung zum Lesen")
	}
	snap := &snapshot{conn: conn, offset: export.offset, blockSize: a.blockSize}
	for _, statement := range []string{
		"BEGIN ISOLATION LEVEL REPEATABLE READ READ ONLY",
		"SET TRANSACTION SNAPSHOT '" + export.snapshotName + "'",
	} {
		if err := snap.exec(ctx, statement); err != nil {
			snap.abort()
			return nil, classify(ctx, err, outbound.ErrSnapshotStorage, "Snapshot-Import")
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
	if err := snap.exec(ctx, cursorStatement(columns, schema, table)); err != nil {
		snap.abort()
		return nil, classify(ctx, err, outbound.ErrSnapshotStorage, "Cursor")
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
		return nil, classify(ctx, result.Err, outbound.ErrSnapshotStorage, "Spaltenliste")
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

// cursorStatement bildet die Cursor-Deklaration: `col::text` je Spalte
// liefert den Text-Stand der Quelle, `NO SCROLL` hält keinen Rückwärtspuffer.
func cursorStatement(columns []string, schema, table string) string {
	quoted := make([]string, len(columns))
	for i, column := range columns {
		quoted[i] = quoteIdent(column) + "::text"
	}
	return "DECLARE " + cursorName + " NO SCROLL CURSOR FOR SELECT " +
		strings.Join(quoted, ", ") + " FROM " + quoteIdent(schema) + "." + quoteIdent(table)
}

// quoteIdent setzt einen Bezeichner in Anführungszeichen.
func quoteIdent(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

// EstimatedRows liest `pg_class.reltuples` der Tabelle über eine kurze
// eigene Verbindung; ein negativer Wert meldet „unbekannt".
func (a *PostgresTableSnapshotAdapter) EstimatedRows(ctx context.Context, schema, table string) (int64, bool, error) {
	if schema == "" || table == "" {
		return 0, false, fmt.Errorf("%w: Tabelle ohne Namen", outbound.ErrSnapshotConfiguration)
	}
	config, err := pgconn.ParseConfig(a.dsn)
	if err != nil {
		return 0, false, fmt.Errorf("%w: DSN: %v", outbound.ErrSnapshotConfiguration, err)
	}
	delete(config.RuntimeParams, "replication")
	conn, err := pgconn.ConnectConfig(ctx, config)
	if err != nil {
		return 0, false, classify(ctx, err, outbound.ErrSnapshotStorage, "Verbindung zum Lesen")
	}
	defer closeConn(conn)

	result := conn.ExecParams(ctx,
		`SELECT c.reltuples::float8::text
		   FROM pg_class c
		  WHERE c.oid = to_regclass(quote_ident($1) || '.' || quote_ident($2))`,
		[][]byte{[]byte(schema), []byte(table)}, nil, nil, nil).Read()
	if result.Err != nil {
		return 0, false, classify(ctx, result.Err, outbound.ErrSnapshotStorage, "Zeilenschätzung")
	}
	if len(result.Rows) == 0 {
		return 0, false, fmt.Errorf("%w: Tabelle %s.%s nicht vorhanden", outbound.ErrSnapshotConfiguration, schema, table)
	}
	value, err := strconv.ParseFloat(string(result.Rows[0][0]), 64)
	if err != nil {
		return 0, false, fmt.Errorf("%w: Zeilenschätzung %q nicht lesbar", outbound.ErrSnapshotStorage, result.Rows[0][0])
	}
	if value < 0 {
		return 0, false, nil
	}
	return int64(math.Round(value)), true, nil
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

// NextBlock holt die nächsten höchstens `blockSize` Zeilen des Cursors.
func (s *snapshot) NextBlock(ctx context.Context) ([][]*string, error) {
	if s.closed {
		return nil, fmt.Errorf("%w: Snapshot geschlossen", outbound.ErrSnapshotStorage)
	}
	results, err := s.conn.Exec(ctx, "FETCH FORWARD "+strconv.Itoa(s.blockSize)+" FROM "+cursorName).ReadAll()
	if err != nil {
		return nil, classify(ctx, err, outbound.ErrSnapshotStorage, "Cursor lesen")
	}
	block := make([][]*string, 0, s.blockSize)
	for _, result := range results {
		for _, row := range result.Rows {
			values := make([]*string, len(row))
			for i, field := range row {
				if field != nil {
					text := string(field)
					values[i] = &text
				}
			}
			block = append(block, values)
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

// classify übersetzt einen Treiber-Fehler an der Adapter-Grenze in die
// Fehlerklasse (`ADR-0023`): ein Kontextende und ein Verbindungsfehler sind
// `transient`, `42501` und die Anmelde-Klasse `28…` sind `permission`, die
// Konfigurationsgrenzen der Quelle (`53400`, `53300`) sind `configuration`;
// der Rest trägt die Klasse der Phase (`fallback`).
func classify(ctx context.Context, cause, fallback error, phase string) error {
	var pgErr *pgconn.PgError
	var connectErr *pgconn.ConnectError
	switch {
	case ctx.Err() != nil || errors.Is(cause, context.DeadlineExceeded) || errors.Is(cause, context.Canceled) || pgconn.Timeout(cause):
		return fmt.Errorf("%w: %s: %v", outbound.ErrSnapshotTransient, phase, cause)
	case errors.As(cause, &pgErr) && (pgErr.Code == "42501" || strings.HasPrefix(pgErr.Code, "28")):
		return fmt.Errorf("%w: %s: %v", outbound.ErrSnapshotPermission, phase, cause)
	case errors.As(cause, &pgErr) && (pgErr.Code == "53400" || pgErr.Code == "53300"):
		return fmt.Errorf("%w: %s: %v", outbound.ErrSnapshotConfiguration, phase, cause)
	case errors.As(cause, &connectErr):
		return fmt.Errorf("%w: %s: %v", outbound.ErrSnapshotTransient, phase, cause)
	default:
		return fmt.Errorf("%w: %s: %v", fallback, phase, cause)
	}
}
