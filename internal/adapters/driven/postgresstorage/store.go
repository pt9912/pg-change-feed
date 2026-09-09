package postgresstorage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage/mapper"
	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage/queries"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// PostgresChangeStoreAdapter ist die Referenzimplementierung des
// `ChangeStorePort` (`ADR-0009`, `ADR-0010`): er persistiert committed
// CDC-Transaktionen und liest persistierte Changes gegen die CDC-Tabellen
// (`SPEC-001`). Der Treiber (pgx/v5) und die PostgreSQL-Typen bleiben
// Adapterdetail (`ADR-0032`, rein Go, CGO-frei).
type PostgresChangeStoreAdapter struct {
	pool *pgxpool.Pool
}

// New baut den Verbindungspool gegen die CDC-Instanz und meldet eine
// nicht erreichbare Instanz als Fehler der Klasse `storage`
// (`outbound.ErrStorage`, `SPEC-008`) — ein Speicher, der nicht erreichbar
// ist, trägt keine Persistenz; der Träger ist storageFailure.
func New(ctx context.Context, dsn string) (*PostgresChangeStoreAdapter, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, storageFailure(err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, storageFailure(err)
	}
	return &PostgresChangeStoreAdapter{pool: pool}, nil
}

// storageFailure trägt die Übersetzungsverantwortung des Adapters
// (`ADR-0023`, `SPEC-008`): Treiber-Fehler gehen an dieser Grenze in die
// Klasse `storage` — Application und Betrieb klassifizieren über
// `errors.Is(err, outbound.ErrStorage)` und kennen keinen Treibertyp; die
// technische Ursache bleibt über die zweite Wrappung lesbar.
func storageFailure(cause error) error {
	return fmt.Errorf("%w: %w", outbound.ErrStorage, cause)
}

// Close schließt den Verbindungspool.
func (a *PostgresChangeStoreAdapter) Close() {
	a.pool.Close()
}

var _ outbound.ChangeStorePort = (*PostgresChangeStoreAdapter)(nil)

// PersistTransaction persistiert die committed Quelltransaktion in EINEM
// Store-Commit (`LH-QA-REL-001.a`, Schritte Persist und COMMIT Store) —
// Changes und Transaktions-Zeile gehen gemeinsam, die Rückkehr ohne Fehler
// meldet den abgeschlossenen Store-Commit. Die Idempotenz (`ADR-0011`)
// trägt die Primärschlüssel über den internen IDs: die erneut persistierte
// Transaktion dedupliziert über ON CONFLICT DO NOTHING, ändert keinen
// Stand und meldet keinen Fehler. Treiber-Fehler gehen in die Klasse
// `storage` (storageFailure); ein Persistenzfehler endet ohne
// Source-ACK (`LH-QA-REL-001.a`). Eine offene Transaktion verwirft der
// Domänen-Träger selbst (`model.ChangeTransaction.Changes`,
// `ADR-0029` Regel 3).
func (a *PostgresChangeStoreAdapter) PersistTransaction(ctx context.Context, transaction *model.ChangeTransaction) error {
	position, committed := transaction.CommitPosition()
	if !committed {
		// Offene Transaktionen sind nicht konsumierbar (`ADR-0029`,
		// Regel 3); ihre Changes liest der Träger nicht.
		return domainerrors.ErrTransactionNotCommitted
	}
	changes, err := transaction.Changes()
	if err != nil {
		return err
	}
	transactionRow, err := mapper.NewTransactionRow(transaction, position)
	if err != nil {
		return err
	}
	changeRows, err := mapper.NewChangeRows(changes)
	if err != nil {
		return err
	}

	tx, err := a.pool.Begin(ctx)
	if err != nil {
		return storageFailure(err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, queries.InsertTransaction,
		transactionRow.TransactionID,
		transactionRow.SourceID,
		transactionRow.CommitPosition,
	); err != nil {
		return storageFailure(err)
	}
	for _, row := range changeRows {
		if _, err := tx.Exec(ctx, queries.InsertChange,
			row.ChangeID,
			row.TransactionID,
			row.SourceTableID,
			row.Sequence,
			row.Operation,
			mapper.JSONImage(row.OldData),
			mapper.JSONImage(row.NewData),
			row.SchemaVersion,
		); err != nil {
			return storageFailure(err)
		}
	}
	return tx.Commit(ctx)
}

// ReadChanges liest persistierte Changes über die `SPEC-001`-Tabellen;
// Bereich, Limit und Filter liegen in der Abfrage, die Ordnung trägt die
// SQL-Sortierung (`LH-FA-REA-004.a`). Treiber-Fehler gehen in die Klasse
// `storage` (storageFailure); die Zeilen laufen zurück durch die
// Domänen-Konstruktoren — eine Zeile, die die Change-Invarianten verletzt,
// endet als sichtbarer Fehler, nicht als still gefälschter Change
// (`ADR-0029`).
func (a *PostgresChangeStoreAdapter) ReadChanges(ctx context.Context, query outbound.ChangeQuery) ([]outbound.ChangeRecord, error) {
	if err := query.Validate(); err != nil {
		return nil, err
	}

	rows, err := a.pool.Query(ctx, queries.SelectChanges,
		string(query.Source),
		positionArgument(query.Start),
		positionArgument(query.End),
		tableArgument(query.Table),
		query.Limit,
	)
	if err != nil {
		return nil, storageFailure(err)
	}
	defer rows.Close()

	records, err := collectRecords(rows)
	if err != nil {
		return nil, err
	}
	return records, nil
}

// positionArgument trägt eine Positions-Grenze als bigint-Argument; nil
// grenzt nicht ein (LIMIT- und NULL-Semantik der Abfrage).
func positionArgument(position *model.SourcePosition) any {
	if position == nil {
		return nil
	}
	return int64(position.Offset)
}

// tableArgument trägt den Tabellenfilter; nil filtert nicht.
func tableArgument(table *model.SourceTableID) any {
	if table == nil {
		return nil
	}
	return string(*table)
}

// collectRecords trägt die Ergebnis-Zeilen in ChangeRecords; die
// Commit-Position kommt von der Transaktion, der Change aus seiner Zeile.
// Lese- und Scan-Fehler des Treibers tragen die Klasse `storage`; die
// Domänen-Konstruktoren tragen ihre Invarianten-Sentinels selbst.
func collectRecords(rows pgx.Rows) ([]outbound.ChangeRecord, error) {
	records := make([]outbound.ChangeRecord, 0)
	for rows.Next() {
		var row mapper.ChangeRow
		var source string
		var commitPosition int64
		if err := rows.Scan(
			&source,
			&commitPosition,
			&row.ChangeID,
			&row.TransactionID,
			&row.SourceTableID,
			&row.Sequence,
			&row.Operation,
			&row.OldData,
			&row.NewData,
			&row.SchemaVersion,
		); err != nil {
			return nil, storageFailure(err)
		}
		position, err := mapper.ToPosition(source, commitPosition)
		if err != nil {
			return nil, err
		}
		change, err := mapper.ToChange(row)
		if err != nil {
			return nil, err
		}
		records = append(records, outbound.ChangeRecord{Position: position, Change: change})
	}
	if err := rows.Err(); err != nil {
		return nil, storageFailure(err)
	}
	return records, nil
}
