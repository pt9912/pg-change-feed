package postgresstorage

import (
	"context"
	"errors"
	"math"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage/mapper"
	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage/queries"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// PostgresConsumerStateAdapter ist die Referenzimplementierung des
// `ConsumerStatePort` (`ARC-004`, `ADR-0013`): er trägt die
// Registrierung, die Bestätigung, den Fortsetzungs-Lese und die
// administrative Entfernung gegen die Consumer-State-Tabellen
// (`SPEC-001`, `cdc.consumer`/`cdc.consumer_position`). Der Treiber
// (pgx/v5) und die PostgreSQL-Typen bleiben Adapterdetail (`ADR-0032`,
// rein Go, CGO-frei); die Tabellenform trägt der d-migrate-Rollout
// (`ADR-0043`) — die DDL des Store-Adapters (schema.sql) trägt die
// Consumer-State-Tabellen nicht.
type PostgresConsumerStateAdapter struct {
	pool *pgxpool.Pool
}

// NewConsumerState baut den Verbindungspool gegen die Instanz, die Quelle
// und den CDC-Speicher gleichermaßen trägt (Abschnitt 1 Lastenheft), und
// meldet eine nicht erreichbare Instanz als Fehler der Klasse `storage`
// (`outbound.ErrStorage`, `SPEC-008`).
func NewConsumerState(ctx context.Context, dsn string) (*PostgresConsumerStateAdapter, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, storageFailure(err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, storageFailure(err)
	}
	return &PostgresConsumerStateAdapter{pool: pool}, nil
}

// Close schließt den Verbindungspool.
func (a *PostgresConsumerStateAdapter) Close() {
	a.pool.Close()
}

var _ outbound.ConsumerStatePort = (*PostgresConsumerStateAdapter)(nil)

// Register trägt die Consumer-Zeile ein; die Idempotenz
// (`LH-FA-CON-001` Boundary) trägt der Primärschlüssel über
// ON CONFLICT DO NOTHING — die erneut registrierte Kennung bleibt ohne
// Wirkung und die Rückkehr meldet den Ausgang. Treiber-Fehler gehen in
// die Klasse `storage` (storageFailure).
func (a *PostgresConsumerStateAdapter) Register(ctx context.Context, consumer model.Consumer) (bool, error) {
	tag, err := a.pool.Exec(ctx, queries.InsertConsumer, string(consumer.ID), consumer.Name)
	if err != nil {
		return false, storageFailure(err)
	}
	return tag.RowsAffected() == 1, nil
}

// Position liest die bestätigte Position des Consumers (`LH-FA-CON-003`
// Happy Path); die Abwesenheit der Zeile liest sich als Nullwert — die
// definierte Anfangsposition eines Consumers ohne Bestätigung
// (`LH-FA-CON-005` Boundary, auch über einen Neustart hinweg). Die Zeile
// läuft zurück durch den Domänen-Konstruktor; die Spalten-Kante der DDL
// hält die Position positiv (`SPEC-001`).
func (a *PostgresConsumerStateAdapter) Position(ctx context.Context, consumer model.ConsumerID) (model.ConsumerPosition, error) {
	if consumer == "" {
		return model.ConsumerPosition{}, domainerrors.ErrEmptyIdentifier
	}
	var source string
	var offset int64
	// Der Lese trägt die Positions-Sperre nicht: er läuft als eigener
	// Aufruf ohne Transaktion, die Zeilen-Sperre der Bestätigung
	// (SelectConsumerPositionLocked) hält ihn nicht auf.
	err := a.pool.QueryRow(ctx, queries.SelectConsumerPosition, string(consumer)).Scan(&source, &offset)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.ConsumerPosition{}, nil
		}
		return model.ConsumerPosition{}, storageFailure(err)
	}
	position, err := mapper.ToPosition(source, offset)
	if err != nil {
		return model.ConsumerPosition{}, err
	}
	return model.ConsumerPosition{ConsumerID: consumer, Position: position}, nil
}

// Acknowledge trägt die bestätigte Position fort (`LH-FA-CON-004`): der
// gespeicherte Fortschritt liest innerhalb des Store-Commits gegen die
// Zeilen-Sperre (SelectConsumerPositionLocked) und führt über den
// Domänen-Vergleich (`model.ConsumerPosition.Advance`, `ADR-0029` Regel 2)
// — eine frühere Position und eine Position einer anderen Quelle enden
// über die Invarianten-Sentinels, die Wiederholung derselben Position ist
// idempotent. Ein Consumer ohne Zeile trägt seine erste Bestätigung als
// neue Zeile. Treiber-Fehler gehen in die Klasse `storage`
// (storageFailure); der Aufruf ohne Registrierung endet über den
// Fremdschlüssel sichtbar (`SPEC-001`).
func (a *PostgresConsumerStateAdapter) Acknowledge(ctx context.Context, position model.ConsumerPosition) (model.ConsumerPosition, error) {
	if position.ConsumerID == "" {
		return model.ConsumerPosition{}, domainerrors.ErrEmptyIdentifier
	}
	if position.Position.IsZero() {
		return model.ConsumerPosition{}, domainerrors.ErrInvalidPosition
	}
	if position.Position.Offset > math.MaxInt64 {
		return model.ConsumerPosition{}, mapper.ErrPositionOutOfRange
	}

	tx, err := a.pool.Begin(ctx)
	if err != nil {
		return model.ConsumerPosition{}, storageFailure(err)
	}
	defer tx.Rollback(ctx)

	var source string
	var offset int64
	err = tx.QueryRow(ctx, queries.SelectConsumerPositionLocked, string(position.ConsumerID)).Scan(&source, &offset)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return model.ConsumerPosition{}, storageFailure(err)
	}
	stored, err := storedPosition(position.ConsumerID, source, offset, err)
	if err != nil {
		return model.ConsumerPosition{}, err
	}
	carried, err := stored.Advance(position.Position)
	if err != nil {
		return model.ConsumerPosition{}, err
	}
	if _, err := tx.Exec(ctx, queries.UpsertConsumerPosition,
		string(carried.ConsumerID),
		string(carried.Position.SourceID),
		int64(carried.Position.Offset),
	); err != nil {
		return model.ConsumerPosition{}, storageFailure(err)
	}
	if err := tx.Commit(ctx); err != nil {
		return model.ConsumerPosition{}, storageFailure(err)
	}
	return carried, nil
}

// storedPosition trägt den gesperrten Fortschritt aus dem Lese: eine
// Zeile liest sich als Quellposition, ihre Abwesenheit als Nullwert des
// Consumers — der Träger der ersten Bestätigung.
func storedPosition(consumer model.ConsumerID, source string, offset int64, scanErr error) (model.ConsumerPosition, error) {
	if errors.Is(scanErr, pgx.ErrNoRows) {
		return model.NewConsumerPosition(consumer)
	}
	if scanErr != nil {
		return model.ConsumerPosition{}, storageFailure(scanErr)
	}
	stored, err := mapper.ToPosition(source, offset)
	if err != nil {
		return model.ConsumerPosition{}, err
	}
	return model.ConsumerPosition{ConsumerID: consumer, Position: stored}, nil
}

// Remove entzieht die Consumer-Zeile und ihre Position-Zeile in EINEM
// Store-Commit; der dokumentierte Ausgang der bestätigten Position ist
// der Entzug mit der Zeile (`LH-FA-CON-006` Boundary) — kein Change und
// keine Quell-Zeile werden gelöscht. Die Abwesenheit der Zeile bleibt
// ohne Wirkung (Idempotenz); die Rückkehr meldet den Ausgang.
func (a *PostgresConsumerStateAdapter) Remove(ctx context.Context, consumer model.ConsumerID) (bool, error) {
	if consumer == "" {
		return false, domainerrors.ErrEmptyIdentifier
	}
	tx, err := a.pool.Begin(ctx)
	if err != nil {
		return false, storageFailure(err)
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, queries.DeleteConsumerPosition, string(consumer)); err != nil {
		return false, storageFailure(err)
	}
	tag, err := tx.Exec(ctx, queries.DeleteConsumer, string(consumer))
	if err != nil {
		return false, storageFailure(err)
	}
	if err := tx.Commit(ctx); err != nil {
		return false, storageFailure(err)
	}
	return tag.RowsAffected() == 1, nil
}
