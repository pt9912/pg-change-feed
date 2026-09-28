package postgresstorage

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage/queries"
	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage/sqlexec"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// PostgresDiagnosticsAdapter ist die Referenzimplementierung des
// `DiagnosticsPort` (`ADR-0132`): er bündelt die sechs Diagnose-
// Lesezugriffe, die bislang als freie Funktion in
// `internal/bootstrap/wiring.go` standen, hinter einem Port — dieselben
// SQL-Texte (`queries.SelectDiagnostics*`), mechanisch verschoben, kein
// neuer SQL-Text. `db` trägt die Ausführung über die schmale Naht
// (`sqlexec`).
type PostgresDiagnosticsAdapter struct {
	db  sqlexec.DB
	log outbound.LogPort
}

// NewDiagnostics baut den Verbindungspool gegen die Instanz (`ADR-0132`):
// dieselbe Konstruktions-Form wie jeder andere `postgresstorage`-Adapter
// (`NewHeartbeat`, …). Der Aufrufer übergibt `cfg.ReaderDSN` — alle sechs
// Views tragen ein `SELECT`-Grant an `cdc_reader`.
func NewDiagnostics(ctx context.Context, dsn string, opts ...Option) (*PostgresDiagnosticsAdapter, error) {
	o := newOptions(opts)
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, diagnosticsStorageFailure(ctx, o.log, err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, diagnosticsStorageFailure(ctx, o.log, err)
	}
	o.log.Info(ctx, "diagnostics: verbunden")
	return &PostgresDiagnosticsAdapter{db: pool, log: o.log}, nil
}

// diagnosticsStorageFailure trägt die Übersetzungsverantwortung dieses
// Adapters (`ADR-0132`): Treiber-Fehler gehen an dieser Grenze in die
// Klasse `storage` über den Port-Sentinel (`outbound.ErrDiagnosticsStorage`)
// — die technische Ursache bleibt über die zweite Wrappung lesbar.
func diagnosticsStorageFailure(ctx context.Context, log outbound.LogPort, cause error) error {
	log.Error(ctx, "diagnostics: Datenbankfehler", "error", cause)
	return sqlexec.Classify(outbound.ErrDiagnosticsStorage, cause)
}

// Close schließt den Verbindungspool.
func (a *PostgresDiagnosticsAdapter) Close() {
	a.db.Close()
}

var _ outbound.DiagnosticsPort = (*PostgresDiagnosticsAdapter)(nil)

// Read liest die sechs Diagnose-Signalgruppen der Quelle in einem Aufruf
// (`ADR-0132`): dieselbe Reihenfolge und dieselben Abwesenheits-Lesarten wie
// der bisherige CLI-Text (Heartbeat, CDC-Abstand, Verarbeitungsrückstand je
// Consumer, blockierender Consumer, Speicherverbrauch, Backfill je Tabelle).
func (a *PostgresDiagnosticsAdapter) Read(ctx context.Context, source model.SourceID) (outbound.DiagnosticsSnapshot, error) {
	var snapshot outbound.DiagnosticsSnapshot

	var ageSeconds float64
	var errorClass *string
	switch err := a.db.QueryRow(ctx, queries.SelectDiagnosticsHeartbeat, string(source)).Scan(&ageSeconds, &errorClass); {
	case sqlexec.IsAbsent(err):
		// kein Lebenszeichen — Instanz hat noch nie geschlagen.
	case err != nil:
		return outbound.DiagnosticsSnapshot{}, diagnosticsStorageFailure(ctx, a.log, err)
	default:
		snapshot.HeartbeatAgeSeconds = &ageSeconds
		snapshot.ErrorClass = errorClass
	}

	if err := a.db.QueryRow(ctx, queries.SelectDiagnosticsCaptureLag).Scan(&snapshot.CaptureLag); err != nil {
		return outbound.DiagnosticsSnapshot{}, diagnosticsStorageFailure(ctx, a.log, err)
	}

	consumerLagRows, err := a.db.Query(ctx, queries.SelectDiagnosticsConsumerLags)
	if err != nil {
		return outbound.DiagnosticsSnapshot{}, diagnosticsStorageFailure(ctx, a.log, err)
	}
	snapshot.ConsumerLags, err = scanConsumerLags(consumerLagRows)
	if err != nil {
		return outbound.DiagnosticsSnapshot{}, diagnosticsStorageFailure(ctx, a.log, err)
	}

	var blockerConsumer, blockerName string
	var blockerAckPos int64
	var blockerBacklog *int64
	switch err := a.db.QueryRow(ctx, queries.SelectDiagnosticsRetentionBlocker, string(source)).Scan(
		&blockerConsumer, &blockerName, &blockerAckPos, &blockerBacklog,
	); {
	case sqlexec.IsAbsent(err):
		// kein Blocker — kein Consumer hat je gegen diese Quelle bestätigt.
	case err != nil:
		return outbound.DiagnosticsSnapshot{}, diagnosticsStorageFailure(ctx, a.log, err)
	default:
		snapshot.RetentionBlocker = &outbound.RetentionBlockerSnapshot{
			ConsumerID: blockerConsumer, Name: blockerName,
			AcknowledgedPosition: blockerAckPos, Backlog: blockerBacklog,
		}
	}

	if err := a.db.QueryRow(ctx, queries.SelectDiagnosticsStorageBytes).Scan(&snapshot.StorageBytes); err != nil {
		return outbound.DiagnosticsSnapshot{}, diagnosticsStorageFailure(ctx, a.log, err)
	}

	backfillRows, err := a.db.Query(ctx, queries.SelectDiagnosticsBackfillStatus, string(source))
	if err != nil {
		return outbound.DiagnosticsSnapshot{}, diagnosticsStorageFailure(ctx, a.log, err)
	}
	snapshot.Backfill, err = scanBackfillStatus(backfillRows)
	if err != nil {
		return outbound.DiagnosticsSnapshot{}, diagnosticsStorageFailure(ctx, a.log, err)
	}

	return snapshot, nil
}

// scanConsumerLags liest die Zeilen von `SelectDiagnosticsConsumerLags` in
// eine leere, gesetzte Liste — nie `nil` (`ADR-0132`).
func scanConsumerLags(rows pgx.Rows) ([]outbound.ConsumerLagSnapshot, error) {
	defer rows.Close()
	lags := make([]outbound.ConsumerLagSnapshot, 0)
	for rows.Next() {
		var consumer string
		var lag *float64
		if err := rows.Scan(&consumer, &lag); err != nil {
			return nil, err
		}
		lags = append(lags, outbound.ConsumerLagSnapshot{ConsumerID: consumer, Lag: lag})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return lags, nil
}

// scanBackfillStatus liest die Zeilen von `SelectDiagnosticsBackfillStatus`
// in eine leere, gesetzte Liste — nie `nil` (`ADR-0132`).
func scanBackfillStatus(rows pgx.Rows) ([]outbound.BackfillTableSnapshot, error) {
	defer rows.Close()
	tables := make([]outbound.BackfillTableSnapshot, 0)
	for rows.Next() {
		var t outbound.BackfillTableSnapshot
		if err := rows.Scan(&t.Schema, &t.Table, &t.Status, &t.RowsCopied, &t.EstimatedRows, &t.WarnEstimatedSize, &t.WarnDuration, &t.ErrorMessage); err != nil {
			return nil, err
		}
		tables = append(tables, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tables, nil
}
