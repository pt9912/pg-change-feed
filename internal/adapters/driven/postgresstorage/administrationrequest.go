package postgresstorage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage/queries"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// administrationChannel trägt den `NOTIFY`-Kanal-Namen der Antrags-Queue
// (`ADR-0050`) — dieselbe Implementer-Entscheidung wie
// `tools/schema/nacharbeit-administration.sql` (`cdc.enable_table`/
// `cdc.disable_table`), real über `LISTEN`/pgx `WaitForNotification`
// getestet (`administrationrequest_test.go`).
const administrationChannel = "cdc_administration"

// administrationStorageFailure trägt die Übersetzungsverantwortung dieses
// Adapters (`ADR-0023`, `SPEC-008`): Treiber-Fehler gehen an dieser Grenze
// in die Klasse `storage` über `outbound.ErrAdministrationStorage`.
func administrationStorageFailure(ctx context.Context, log outbound.LogPort, cause error) error {
	log.Error(ctx, "administrationrequest: Datenbankfehler", "error", cause)
	return fmt.Errorf("%w: %w", outbound.ErrAdministrationStorage, cause)
}

// AdministrationRequestAdapter implementiert den
// `AdministrationRequestPort` (`outbound`, `ARC-004`) gegen dieselbe
// Instanz: Lesen offener Anträge und Ergebnis-Vermerk laufen über den
// Verbindungspool der Administrations-Goroutine (`cdc_admin`-Rolle,
// `ADR-0047`), getrennt vom dedizierten `LISTEN`-Verbindungsträger
// (`AdministrationListener` unten) — Pool-Verbindungen sind für
// `WaitForNotification` ungeeignet, weil jede Anfrage eine beliebige
// Pool-Verbindung ziehen kann.
type AdministrationRequestAdapter struct {
	pool *pgxpool.Pool
	log  outbound.LogPort
}

// NewAdministrationRequest baut den Verbindungspool gegen die Instanz und
// meldet eine nicht erreichbare Instanz als Fehler der Klasse `storage`
// (`outbound.ErrAdministrationStorage`, `SPEC-008`).
func NewAdministrationRequest(ctx context.Context, dsn string, opts ...Option) (*AdministrationRequestAdapter, error) {
	o := newOptions(opts)
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, administrationStorageFailure(ctx, o.log, err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, administrationStorageFailure(ctx, o.log, err)
	}
	o.log.Info(ctx, "administrationrequest: verbunden")
	return &AdministrationRequestAdapter{pool: pool, log: o.log}, nil
}

// Close schließt den Verbindungspool.
func (a *AdministrationRequestAdapter) Close() {
	a.pool.Close()
}

var _ outbound.AdministrationRequestPort = (*AdministrationRequestAdapter)(nil)

// ListPending liest die offenen Anträge in Anlage-Reihenfolge.
func (a *AdministrationRequestAdapter) ListPending(ctx context.Context) ([]model.AdministrationRequest, error) {
	rows, err := a.pool.Query(ctx, queries.SelectPendingAdministrationRequests)
	if err != nil {
		return nil, administrationStorageFailure(ctx, a.log, err)
	}
	defer rows.Close()

	requests := make([]model.AdministrationRequest, 0)
	for rows.Next() {
		var id, source, schema, table, kind string
		if err := rows.Scan(&id, &source, &schema, &table, &kind); err != nil {
			return nil, administrationStorageFailure(ctx, a.log, err)
		}
		requests = append(requests, model.AdministrationRequest{
			ID:     model.AdministrationRequestID(id),
			Source: model.SourceID(source),
			Schema: schema,
			Table:  table,
			Kind:   model.AdministrationRequestKind(kind),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, administrationStorageFailure(ctx, a.log, err)
	}
	return requests, nil
}

// MarkApplied vermerkt einen erfolgreich verarbeiteten Antrag; die
// WHERE-Klausel der Query trägt die Idempotenz — ein bereits vermerkter
// Antrag bleibt unverändert, kein Fehler.
func (a *AdministrationRequestAdapter) MarkApplied(ctx context.Context, id model.AdministrationRequestID) error {
	if id == "" {
		return domainerrors.ErrEmptyIdentifier
	}
	if _, err := a.pool.Exec(ctx, queries.UpdateAdministrationRequestApplied, string(id)); err != nil {
		return administrationStorageFailure(ctx, a.log, err)
	}
	a.log.Info(ctx, "administrationrequest: Antrag erledigt", "request_id", id)
	return nil
}

// MarkFailed vermerkt einen gescheiterten Antrag samt Fehlertext.
func (a *AdministrationRequestAdapter) MarkFailed(ctx context.Context, id model.AdministrationRequestID, message string) error {
	if id == "" {
		return domainerrors.ErrEmptyIdentifier
	}
	if _, err := a.pool.Exec(ctx, queries.UpdateAdministrationRequestFailed, string(id), message); err != nil {
		return administrationStorageFailure(ctx, a.log, err)
	}
	a.log.Warn(ctx, "administrationrequest: Antrag gescheitert", "request_id", id, "error", message)
	return nil
}

// AdministrationListener trägt die `LISTEN`-Wecksignal-Fähigkeit der
// Administrations-Goroutine über eine dedizierte, vom Pool getrennte
// Verbindung (`ADR-0050`): zwischen zwei Wecksignalen steht sie untätig
// und kann von der Quelle oder einem dazwischenliegenden Netzwerkelement
// getrennt werden (Idle-Timeout) — dasselbe Muster wie
// `receive.WALRetentionChecker`. `dsn` bleibt erhalten, damit
// `WaitForNotification` eine so gestörte Verbindung selbst ersetzen kann,
// statt für den Rest des Prozesslaufs stumm zu bleiben; der periodische
// Fallback-Poll der Administrations-Goroutine deckt die Lücke bis zum
// nächsten erfolgreichen `WaitForNotification`-Aufruf ab.
type AdministrationListener struct {
	dsn  string
	conn *pgx.Conn
}

// NewAdministrationListener baut die eigene Verbindung auf und registriert
// `LISTEN` auf dem Administrations-Kanal.
func NewAdministrationListener(ctx context.Context, dsn string) (*AdministrationListener, error) {
	conn, err := connectAdministrationListener(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return &AdministrationListener{dsn: dsn, conn: conn}, nil
}

// connectAdministrationListener baut eine Verbindung auf und registriert
// `LISTEN` — geteilt zwischen `NewAdministrationListener` und dem
// Wiederverbindungspfad in `WaitForNotification`.
func connectAdministrationListener(ctx context.Context, dsn string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", outbound.ErrAdministrationStorage, err)
	}
	if _, err := conn.Exec(ctx, "LISTEN "+administrationChannel); err != nil {
		_ = conn.Close(ctx)
		return nil, fmt.Errorf("%w: %w", outbound.ErrAdministrationStorage, err)
	}
	return conn, nil
}

// Close schließt die eigene Verbindung.
func (l *AdministrationListener) Close(ctx context.Context) error {
	return l.conn.Close(ctx)
}

// WaitForNotification blockiert, bis ein Wecksignal eintrifft oder `ctx`
// endet (Timeout des Aufrufers trägt den Fallback-Poll-Takt, `ADR-0050`).
// Eine gestörte Verbindung (jeder Fehler außerhalb eines abgelaufenen
// `ctx`) ersetzt der Aufruf selbst — der nächste Aufruf `LISTEN`t erneut,
// statt dass das Wecksignal für den Rest des Prozesslaufs ausbleibt; bis
// zum erfolgreichen Wiederaufbau trägt der Fallback-Poll die
// Anfragen-Verarbeitung.
func (l *AdministrationListener) WaitForNotification(ctx context.Context) error {
	_, err := l.conn.WaitForNotification(ctx)
	if err != nil && ctx.Err() == nil {
		_ = l.conn.Close(context.Background())
		if conn, reconnectErr := connectAdministrationListener(context.Background(), l.dsn); reconnectErr == nil {
			l.conn = conn
		}
	}
	return err
}
