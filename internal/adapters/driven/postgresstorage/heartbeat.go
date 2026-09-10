package postgresstorage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage/queries"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// PostgresHeartbeatAdapter ist die Referenzimplementierung des
// `HeartbeatPort` (`ARC-004`, `ADR-0024`): er trägt das periodische
// Lebenszeichen des Capture-Prozesses gegen `cdc.process_heartbeat`
// fort. Der Treiber (pgx/v5) bleibt Adapterdetail (`ADR-0032`, rein Go,
// CGO-frei); die Tabellenform trägt der d-migrate-Rollout (`ADR-0043`) —
// die DDL des Store-Adapters (schema.sql) trägt sie nicht (dieselbe
// Abgrenzung wie bei den Consumer-State-Tabellen). `log` trägt die
// strukturierte Protokollierung über den injizierten `LogPort`
// (`LH-QA-OPS-004`, `ADR-0024`, `WithLog`) — Default `outbound.NoopLog`.
type PostgresHeartbeatAdapter struct {
	pool *pgxpool.Pool
	log  outbound.LogPort
}

// NewHeartbeat baut den Verbindungspool gegen die Instanz, die Quelle und
// CDC-Speicher gleichermaßen trägt (Abschnitt 1 Lastenheft), und meldet
// eine nicht erreichbare Instanz als Fehler der Klasse `storage`
// (`outbound.ErrHeartbeatStorage`, `SPEC-008`). Der Pool bleibt vom
// Store- und Aktivierungs-Pool getrennt (`internal/bootstrap/wiring.go`)
// — der periodische Schreib-Zug teilt keine Verbindung mit der
// Capture-Persist-ACK-Schleife (`LH-QA-REL-001.a`, slice-012 §6-Risiko).
func NewHeartbeat(ctx context.Context, dsn string, opts ...Option) (*PostgresHeartbeatAdapter, error) {
	o := newOptions(opts)
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, heartbeatStorageFailure(ctx, o.log, err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, heartbeatStorageFailure(ctx, o.log, err)
	}
	o.log.Info(ctx, "heartbeat: verbunden")
	return &PostgresHeartbeatAdapter{pool: pool, log: o.log}, nil
}

// heartbeatStorageFailure trägt die Übersetzungsverantwortung dieses
// Adapters (`ADR-0023`, `SPEC-008`): Treiber-Fehler gehen an dieser
// Grenze in die Klasse `storage` über den eigenen Sentinel des Ports
// (`outbound.ErrHeartbeatStorage`) — die technische Ursache bleibt über
// die zweite Wrappung lesbar. Derselbe Aufruf trägt den strukturierten
// Fehler-Log über den injizierten `LogPort` (`LH-QA-OPS-004`,
// `ADR-0024`), aus demselben Grund mit explizitem `ctx`/`log`-Parameter
// wie `storageFailure` (`store.go`).
func heartbeatStorageFailure(ctx context.Context, log outbound.LogPort, cause error) error {
	log.Error(ctx, "heartbeat: Datenbankfehler", "error", cause)
	return fmt.Errorf("%w: %w", outbound.ErrHeartbeatStorage, cause)
}

// Close schließt den Verbindungspool.
func (a *PostgresHeartbeatAdapter) Close() {
	a.pool.Close()
}

var _ outbound.HeartbeatPort = (*PostgresHeartbeatAdapter)(nil)

// Beat trägt die Lebenszeichen-Zeile der Quelle fort (UPSERT): der
// Zeitstempel trägt die Instanzzeit der Speicherseite; eine erneute
// Quelle überschreibt den vorigen Zeitstempel, keine zweite Zeile. Ein
// zuvor gemeldeter Fehlerzustand (Fault) geht dabei verloren — bewusst:
// der nächste erfolgreiche Beat ist das Zeichen, dass die Erfassung
// wieder normal läuft (`LH-FA-ADM-003` Boundary). Treiber-Fehler gehen in
// die Klasse `storage` (heartbeatStorageFailure).
func (a *PostgresHeartbeatAdapter) Beat(ctx context.Context, source model.SourceID) error {
	if source == "" {
		return domainerrors.ErrEmptyIdentifier
	}
	if _, err := a.pool.Exec(ctx, queries.UpsertHeartbeat, string(source)); err != nil {
		return heartbeatStorageFailure(ctx, a.log, err)
	}
	a.log.Debug(ctx, "heartbeat: Lebenszeichen geschrieben", "source", source)
	return nil
}

// Fault trägt den zuletzt beobachteten Fehlerzustand der Quelle fort
// (`slice-013`, `LH-FA-ADM-003`, `LH-QA-REL-003`) — dieselbe Zeile und
// derselbe Zeitstempel-Mechanismus wie Beat, nur mit Klasse. Eine leere
// Quelle oder eine Klasse außerhalb der sieben `ADR-0023`-Kategorien
// erreicht keinen SQL-Aufruf (Port-Grenze, wie bei Beat).
func (a *PostgresHeartbeatAdapter) Fault(ctx context.Context, source model.SourceID, class model.ErrorClass) error {
	if source == "" {
		return domainerrors.ErrEmptyIdentifier
	}
	validated, err := model.NewErrorClass(string(class))
	if err != nil {
		return err
	}
	if _, err := a.pool.Exec(ctx, queries.UpsertHeartbeatFault, string(source), string(validated)); err != nil {
		return heartbeatStorageFailure(ctx, a.log, err)
	}
	a.log.Warn(ctx, "heartbeat: Fehlerzustand gemeldet", "source", source, "class", validated)
	return nil
}
