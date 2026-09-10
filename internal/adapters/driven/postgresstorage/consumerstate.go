package postgresstorage

import (
	"context"
	"errors"
	"fmt"
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
// Consumer-State-Tabellen nicht. `log` trägt die strukturierte
// Protokollierung über den injizierten `LogPort` (`LH-QA-OPS-004`,
// `ADR-0024`, `WithLog`) — Default `outbound.NoopLog`.
type PostgresConsumerStateAdapter struct {
	pool *pgxpool.Pool
	log  outbound.LogPort
}

// NewConsumerState baut den Verbindungspool gegen die Instanz, die Quelle
// und den CDC-Speicher gleichermaßen trägt (Abschnitt 1 Lastenheft), und
// meldet eine nicht erreichbare Instanz als Fehler der Klasse `storage`
// (`outbound.ErrConsumerStateStorage`, `SPEC-008`).
func NewConsumerState(ctx context.Context, dsn string, opts ...Option) (*PostgresConsumerStateAdapter, error) {
	o := newOptions(opts)
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, stateStorageFailure(ctx, o.log, err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, stateStorageFailure(ctx, o.log, err)
	}
	o.log.Info(ctx, "consumerstate: verbunden")
	return &PostgresConsumerStateAdapter{pool: pool, log: o.log}, nil
}

// stateStorageFailure trägt die Übersetzungsverantwortung dieses Adapters
// (`ADR-0023`, `SPEC-008`): Treiber-Fehler gehen an dieser Grenze in die
// Klasse `storage` über den eigenen Sentinel des Ports
// (`outbound.ErrConsumerStateStorage`) — die Klasse-Aktion des
// ChangeStore-Sentinels (kein Source-ACK, `LH-QA-REL-001.a`) trägt dieser
// Adapter nicht; die technische Ursache bleibt über die zweite Wrappung
// lesbar. Derselbe Aufruf trägt den strukturierten Fehler-Log über den
// injizierten `LogPort` (`LH-QA-OPS-004`, `ADR-0024`), aus demselben
// Grund mit explizitem `ctx`/`log`-Parameter wie `storageFailure`
// (`store.go`).
func stateStorageFailure(ctx context.Context, log outbound.LogPort, cause error) error {
	log.Error(ctx, "consumerstate: Datenbankfehler", "error", cause)
	return fmt.Errorf("%w: %w", outbound.ErrConsumerStateStorage, cause)
}

// Close schließt den Verbindungspool.
func (a *PostgresConsumerStateAdapter) Close() {
	a.pool.Close()
}

var _ outbound.ConsumerStatePort = (*PostgresConsumerStateAdapter)(nil)

// Register trägt die Consumer-Zeile ein; die Kennungs- und Namens-Grenze
// läuft vor dem ersten SQL-Aufruf über den Domänen-Konstruktor — dieselbe
// Grenze wie bei Position, Acknowledge und Remove (`ADR-0029`); die
// Idempotenz (`LH-FA-CON-001` Boundary) trägt der Primärschlüssel über
// ON CONFLICT DO NOTHING — die erneut registrierte Kennung bleibt ohne
// Wirkung und die Rückkehr meldet den Ausgang. Treiber-Fehler gehen in
// die Klasse `storage` (stateStorageFailure).
func (a *PostgresConsumerStateAdapter) Register(ctx context.Context, consumer model.Consumer) (bool, error) {
	valid, err := model.NewConsumer(consumer.ID, consumer.Name)
	if err != nil {
		return false, err
	}
	tag, err := a.pool.Exec(ctx, queries.InsertConsumer, string(valid.ID), valid.Name)
	if err != nil {
		return false, stateStorageFailure(ctx, a.log, err)
	}
	registered := tag.RowsAffected() == 1
	if registered {
		a.log.Info(ctx, "consumerstate: Consumer registriert", "consumer_id", valid.ID)
	}
	return registered, nil
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
		return model.ConsumerPosition{}, stateStorageFailure(ctx, a.log, err)
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
// (stateStorageFailure); der Aufruf ohne Registrierung endet über
// `outbound.ErrConsumerUnregistered` sichtbar — die Rest-Grenze zwischen
// Prüfung und Schreiben trägt der Fremdschlüssel der DDL über dieselbe
// Klasse (`SPEC-001`).
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
		return model.ConsumerPosition{}, stateStorageFailure(ctx, a.log, err)
	}
	defer tx.Rollback(ctx)

	var registered int
	if err := tx.QueryRow(ctx, queries.SelectConsumer, string(position.ConsumerID)).Scan(&registered); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.ConsumerPosition{}, outbound.ErrConsumerUnregistered
		}
		return model.ConsumerPosition{}, stateStorageFailure(ctx, a.log, err)
	}

	var source string
	var offset int64
	err = tx.QueryRow(ctx, queries.SelectConsumerPositionLocked, string(position.ConsumerID)).Scan(&source, &offset)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return model.ConsumerPosition{}, stateStorageFailure(ctx, a.log, err)
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
		return model.ConsumerPosition{}, stateStorageFailure(ctx, a.log, err)
	}
	if err := tx.Commit(ctx); err != nil {
		return model.ConsumerPosition{}, stateStorageFailure(ctx, a.log, err)
	}
	a.log.Debug(ctx, "consumerstate: Position bestätigt",
		"consumer_id", carried.ConsumerID, "offset", carried.Position.Offset)
	return carried, nil
}

// storedPosition trägt den gesperrten Fortschritt aus dem Lese: eine
// Zeile liest sich als Quellposition, ihre Abwesenheit als Nullwert des
// Consumers — der Träger der ersten Bestätigung. `stateStorageFailure`
// hier läuft ohne Port-Log (kein `log`/`ctx` in dieser Funktionssignatur):
// der einzige Fehlerpfad ist ein bereits von `Acknowledge` gelesener
// Scan-Fehler, dessen Log-Aufruf am Lese-Aufruf selbst nicht dupliziert
// werden soll — siehe Aufrufstelle.
func storedPosition(consumer model.ConsumerID, source string, offset int64, scanErr error) (model.ConsumerPosition, error) {
	if errors.Is(scanErr, pgx.ErrNoRows) {
		return model.NewConsumerPosition(consumer)
	}
	if scanErr != nil {
		return model.ConsumerPosition{}, fmt.Errorf("%w: %w", outbound.ErrConsumerStateStorage, scanErr)
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
		return false, stateStorageFailure(ctx, a.log, err)
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, queries.DeleteConsumerPosition, string(consumer)); err != nil {
		return false, stateStorageFailure(ctx, a.log, err)
	}
	tag, err := tx.Exec(ctx, queries.DeleteConsumer, string(consumer))
	if err != nil {
		return false, stateStorageFailure(ctx, a.log, err)
	}
	if err := tx.Commit(ctx); err != nil {
		return false, stateStorageFailure(ctx, a.log, err)
	}
	return tag.RowsAffected() == 1, nil
}
