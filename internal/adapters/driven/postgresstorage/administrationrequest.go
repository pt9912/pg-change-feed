package postgresstorage

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage/queries"
	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage/sqlexec"
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
	return sqlexec.Classify(outbound.ErrAdministrationStorage, cause)
}

// AdministrationRequestAdapter implementiert den
// `AdministrationRequestPort` (`outbound`, `ARC-004`) gegen dieselbe
// Instanz: Lesen offener Anträge und Ergebnis-Vermerk laufen über den
// Verbindungspool der Administrations-Goroutine (`cdc_admin`-Rolle,
// `ADR-0047`), getrennt vom dedizierten `LISTEN`-Verbindungsträger
// (`AdministrationListener` unten) — Pool-Verbindungen sind für
// `WaitForNotification` ungeeignet, weil jede Anfrage eine beliebige
// Pool-Verbindung ziehen kann. `db` trägt die Ausführung über die schmale
// Naht (`sqlexec`, `ADR-0071` Punkt 5).
type AdministrationRequestAdapter struct {
	db  sqlexec.DB
	log outbound.LogPort
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
	return &AdministrationRequestAdapter{db: pool, log: o.log}, nil
}

// Close schließt den Verbindungspool.
func (a *AdministrationRequestAdapter) Close() {
	a.db.Close()
}

var _ outbound.AdministrationRequestPort = (*AdministrationRequestAdapter)(nil)

// ListPending liest die offenen Anträge in Aufruf-Reihenfolge; eine vom
// Antrags-Konstruktor verworfene Zeile steht mit Kennung und Fehlertext an
// ihrer Stelle (`sqlexec.ReadPendingRequests`).
func (a *AdministrationRequestAdapter) ListPending(ctx context.Context) ([]outbound.PendingAdministrationRequest, error) {
	return sqlexec.ReadPendingRequests(ctx, a.db, sqlexec.Statement{
		SQL:  queries.SelectPendingAdministrationRequests,
		Fail: func(cause error) error { return administrationStorageFailure(ctx, a.log, cause) },
	})
}

// MarkApplied vermerkt einen erfolgreich verarbeiteten Antrag; die
// WHERE-Klausel der Query trägt die Idempotenz — ein bereits vermerkter
// Antrag (etwa ein `backfill`-Antrag, den die Annahme vermerkt hat) bleibt
// unverändert, kein Fehler und keine Erfolgsmeldung.
func (a *AdministrationRequestAdapter) MarkApplied(ctx context.Context, id model.AdministrationRequestID) error {
	if id == "" {
		return domainerrors.ErrEmptyIdentifier
	}
	tag, err := a.db.Exec(ctx, queries.UpdateAdministrationRequestApplied, string(id))
	if err != nil {
		return administrationStorageFailure(ctx, a.log, err)
	}
	if tag.RowsAffected() > 0 {
		a.log.Info(ctx, "administrationrequest: Antrag erledigt", "request_id", id)
	}
	return nil
}

// MarkFailed vermerkt einen gescheiterten Antrag samt Fehlertext.
func (a *AdministrationRequestAdapter) MarkFailed(ctx context.Context, id model.AdministrationRequestID, message string) error {
	if id == "" {
		return domainerrors.ErrEmptyIdentifier
	}
	if _, err := a.db.Exec(ctx, queries.UpdateAdministrationRequestFailed, string(id), message); err != nil {
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
	// reconnectBackoff trägt die Wartezeit vor dem nächsten
	// `WaitForNotification`-Rückgabewert nach einem gescheiterten
	// Wiederverbindungsversuch: 0 vor dem ersten Fehlschlag und nach jedem
	// erfolgreichen Wiederaufbau — `nextAdministrationReconnectBackoff`
	// verdoppelt sie ab da, gedeckelt bei
	// `administrationReconnectMaxBackoff`. Ohne sie liefe ein dauerhaft
	// unerreichbares `AdminDSN` in eine ungedrosselte Wiederholschleife
	// (jeder `WaitForNotification`-Aufruf des Aufrufers träfe sofort auf
	// denselben Fehler).
	reconnectBackoff time.Duration
}

// administrationReconnectInitialBackoff und administrationReconnectMaxBackoff
// tragen die Backoff-Grenzen des `LISTEN`-Wiederverbindungspfads: die erste
// Wartezeit nach einem Fehlschlag ist kurz (schnelle Erholung von einem
// kurzen Netzwerk-Blip), die Obergrenze verhindert, dass ein dauerhaft
// unerreichbares `AdminDSN` den Fallback-Poll-Takt der
// Administrations-Goroutine (`administrationPollInterval`, `wiring.go`)
// deutlich überschreitet.
const (
	administrationReconnectInitialBackoff = 200 * time.Millisecond
	administrationReconnectMaxBackoff     = 30 * time.Second
)

// nextAdministrationReconnectBackoff verdoppelt die aktuelle Wartezeit
// (0 startet bei der Initial-Backoff), gedeckelt bei der Obergrenze —
// reine Funktion, netzlos testbar.
func nextAdministrationReconnectBackoff(current time.Duration) time.Duration {
	if current <= 0 {
		return administrationReconnectInitialBackoff
	}
	doubled := current * 2
	if doubled > administrationReconnectMaxBackoff {
		return administrationReconnectMaxBackoff
	}
	return doubled
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
// Anfragen-Verarbeitung. Ein gescheiterter Wiederverbindungsversuch wartet
// zusätzlich den Backoff aus `reconnectBackoff` ab (oder bis `ctx` endet),
// bevor der Aufruf zurückkehrt — ein dauerhaft unerreichbares `AdminDSN`
// läuft damit nicht in eine ungedrosselte Wiederholschleife. Ein
// erfolgreiches Wecksignal oder ein erfolgreicher Wiederaufbau setzt den
// Backoff auf 0 zurück.
func (l *AdministrationListener) WaitForNotification(ctx context.Context) error {
	_, err := l.conn.WaitForNotification(ctx)
	if err == nil {
		l.reconnectBackoff = 0
		return nil
	}
	if ctx.Err() != nil {
		return err
	}
	_ = l.conn.Close(context.Background())
	if conn, reconnectErr := connectAdministrationListener(context.Background(), l.dsn); reconnectErr == nil {
		l.conn = conn
		l.reconnectBackoff = 0
		return err
	}
	l.reconnectBackoff = nextAdministrationReconnectBackoff(l.reconnectBackoff)
	select {
	case <-time.After(l.reconnectBackoff):
	case <-ctx.Done():
	}
	return err
}
