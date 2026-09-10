package postgresstorage

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage/queries"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// ErrActivationConfiguration trägt die Fehlerklasse `configuration` des
// Aktivierungs-Adapters (`SPEC-008`, `ADR-0023`): Bezeichner, die nicht
// ins Bezeichner-Alphabet der Quelle passen, enden vor dem ersten
// SQL-Aufruf — kein Start im falschen Stand.
var ErrActivationConfiguration = fmt.Errorf("Fehlerklasse configuration: Aktivierung ohne gültige Bezeichner")

// identifierShape begrenzt die Bezeichner der Aktivierung auf das
// Alphabet der Quelle; Publication, Schema und Tabellenname gehen als
// Bezeichner-Literale in die Publication-DDL. Die Quelle der Regel trägt
// der Stream-Adapter (`receive.identifierShape`); dieser Ausdruck hält
// denselben Alphabet-Vertrag für denselben Aufrufgegenstand — die
// Adapter-Schicht importiert keine Adapter-Kante (Kopplung).
var identifierShape = regexp.MustCompile(`^[a-z0-9_]{1,63}$`)

// TableActivationAdapter implementiert den `TableActivationPort`
// (`outbound`, `ARC-004`) gegen dieselbe Instanz: im MVP trägt eine
// Instanz die Quelle und den CDC-Speicher gleichermaßen (Abschnitt 1
// Lastenheft) — die Bindungs-Zeilen der CDC-Referenztabellen (`SPEC-001`)
// und die Publication der Quelle laufen über denselben Verbindungspool
// (`LH-FA-CFG-001.a`). `log` trägt die strukturierte Protokollierung über
// den injizierten `LogPort` (`LH-QA-OPS-004`, `ADR-0024`, `WithLog`) —
// Default `outbound.NoopLog`.
type TableActivationAdapter struct {
	pool *pgxpool.Pool
	log  outbound.LogPort
}

// NewTableActivation baut den Verbindungspool gegen die Instanz und
// meldet eine nicht erreichbare Instanz als Fehler der Klasse `storage`
// (`outbound.ErrStorage`, `SPEC-008`).
func NewTableActivation(ctx context.Context, dsn string, opts ...Option) (*TableActivationAdapter, error) {
	o := newOptions(opts)
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, storageFailure(ctx, o.log, err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, storageFailure(ctx, o.log, err)
	}
	o.log.Info(ctx, "tableactivation: verbunden")
	return &TableActivationAdapter{pool: pool, log: o.log}, nil
}

// Close schließt den Verbindungspool.
func (a *TableActivationAdapter) Close() {
	a.pool.Close()
}

var _ outbound.TableActivationPort = (*TableActivationAdapter)(nil)

// TableExists prüft die physische Tabelle über den Katalog; die
// Negative-Pfade der Aktivierung, Deaktivierung und Status-Abfrage enden
// über die Abwesenheit sichtbar (`LH-FA-CFG-001`/`002`/`003`).
func (a *TableActivationAdapter) TableExists(ctx context.Context, schema, table string) (bool, error) {
	if err := validateIdentifier(schema); err != nil {
		return false, err
	}
	if err := validateIdentifier(table); err != nil {
		return false, err
	}
	var count int
	if err := a.pool.QueryRow(ctx,
		"SELECT count(*) FROM information_schema.tables WHERE table_schema = $1 AND table_name = $2",
		schema, table,
	).Scan(&count); err != nil {
		return false, storageFailure(ctx, a.log, err)
	}
	return count > 0, nil
}

// Registered liest die Bindungs-Zeile der Tabelle; die Abwesenheit ist
// der Zustand „nicht aktiviert" (`LH-FA-CFG-003` Boundary).
func (a *TableActivationAdapter) Registered(ctx context.Context, source model.SourceID, schema, table string) (model.SourceTable, bool, error) {
	var id string
	if err := a.pool.QueryRow(ctx, queries.SelectSourceTable, string(source), schema, table).Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.SourceTable{}, false, nil
		}
		return model.SourceTable{}, false, storageFailure(ctx, a.log, err)
	}
	registered, err := model.NewSourceTable(model.SourceTableID(id), source, schema, table)
	if err != nil {
		return model.SourceTable{}, false, err
	}
	return registered, true, nil
}

// Register trägt die Bindungs- und Schema-Version-Zeile in EINEM
// Store-Commit; die Idempotenz (`LH-FA-CFG-001` Boundary) tragen die
// UNIQUE-Kanten über ON CONFLICT DO NOTHING — die erneut aktivierte
// Tabelle bleibt ohne Wirkung und die Rückkehr meldet den Ausgang. Die
// Zeile `cdc.source` trägt der Aufrufer vor der Aktivierung; ein
// fremdschlüssel-verletzender Aufruf endet über die Klasse `storage`
// sichtbar.
func (a *TableActivationAdapter) Register(ctx context.Context, table model.SourceTable, version model.SchemaVersion) (bool, error) {
	_, registered, err := a.Registered(ctx, table.SourceID, table.Schema, table.Table)
	if err != nil {
		return false, err
	}
	if registered {
		return false, nil
	}
	tx, err := a.pool.Begin(ctx)
	if err != nil {
		return false, storageFailure(ctx, a.log, err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, queries.InsertSourceTable,
		string(table.ID),
		string(table.SourceID),
		table.Schema,
		table.Table,
	); err != nil {
		return false, storageFailure(ctx, a.log, err)
	}
	if _, err := tx.Exec(ctx, queries.InsertSchemaVersion,
		string(version.ID),
		string(version.SourceTableID),
		version.Version,
	); err != nil {
		return false, storageFailure(ctx, a.log, err)
	}
	if err := tx.Commit(ctx); err != nil {
		return false, storageFailure(ctx, a.log, err)
	}
	a.log.Info(ctx, "tableactivation: Tabelle registriert",
		"table_id", table.ID, "schema", table.Schema, "table", table.Table)
	return true, nil
}

// Unregister entzieht die Bindungs-Zeile und die Schema-Version-Zeilen in
// EINEM Store-Commit; der Change-Bestand läuft vorab: bei persistierten
// Changes bleibt die Zeile als Herkunft bestehen (ActivationRetained),
// ohne Bestand trägt der Entzug beide Zeilenmengen — kein Change und
// keine Schema-Version werden gelöscht.
func (a *TableActivationAdapter) Unregister(ctx context.Context, table model.SourceTable) (outbound.ActivationRemoval, error) {
	var hasChanges bool
	if err := a.pool.QueryRow(ctx, queries.SelectSourceTableChanges, string(table.ID)).Scan(&hasChanges); err != nil {
		return "", storageFailure(ctx, a.log, err)
	}
	if hasChanges {
		return outbound.ActivationRetained, nil
	}
	tx, err := a.pool.Begin(ctx)
	if err != nil {
		return "", storageFailure(ctx, a.log, err)
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, queries.DeleteSchemaVersions, string(table.ID)); err != nil {
		return "", storageFailure(ctx, a.log, err)
	}
	tag, err := tx.Exec(ctx, queries.DeleteSourceTable, string(table.ID))
	if err != nil {
		return "", storageFailure(ctx, a.log, err)
	}
	if err := tx.Commit(ctx); err != nil {
		return "", storageFailure(ctx, a.log, err)
	}
	if tag.RowsAffected() == 0 {
		return outbound.ActivationAbsent, nil
	}
	return outbound.ActivationRemoved, nil
}

// List liest die Bindungs-Zeilen der Quelle in Schema- und
// Tabellen-Ordnung (`LH-FA-CFG-004`).
func (a *TableActivationAdapter) List(ctx context.Context, source model.SourceID) ([]model.SourceTable, error) {
	rows, err := a.pool.Query(ctx, queries.SelectSourceTables, string(source))
	if err != nil {
		return nil, storageFailure(ctx, a.log, err)
	}
	defer rows.Close()

	tables := make([]model.SourceTable, 0)
	for rows.Next() {
		var id, tableSource, schema, table string
		if err := rows.Scan(&id, &tableSource, &schema, &table); err != nil {
			return nil, storageFailure(ctx, a.log, err)
		}
		entry, err := model.NewSourceTable(model.SourceTableID(id), model.SourceID(tableSource), schema, table)
		if err != nil {
			return nil, err
		}
		tables = append(tables, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, storageFailure(ctx, a.log, err)
	}
	return tables, nil
}

// Publish trägt die Publication (`LH-FA-CFG-001.a`): eine fehlende
// Publication legt der Aufruf mit der Tabelle und den erfassten
// Operationen an (INSERT, UPDATE, DELETE — TRUNCATE trägt die Publication
// nicht, `LH-FA-CFG-001.a` Schritt 3), eine bestehende trägt die Tabelle
// nach, wenn sie kein Mitglied ist. Die Bezeichner laufen vor der DDL über
// das Bezeichner-Alphabet und gehen als Literale ein.
func (a *TableActivationAdapter) Publish(ctx context.Context, publication, schema, table string) error {
	publicationName, qualified, err := publicationIdentifiers(publication, schema, table)
	if err != nil {
		return err
	}
	var exists int
	if err := a.pool.QueryRow(ctx, queries.SelectPublication, publication).Scan(&exists); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			if _, err := a.pool.Exec(ctx,
				fmt.Sprintf("CREATE PUBLICATION %s FOR TABLE %s WITH (publish = 'insert, update, delete')",
					publicationName, qualified),
			); err != nil {
				return storageFailure(ctx, a.log, err)
			}
			a.log.Info(ctx, "tableactivation: Publication angelegt", "publication", publication, "table", qualified)
			return nil
		}
		return storageFailure(ctx, a.log, err)
	}
	var member int
	if err := a.pool.QueryRow(ctx, queries.SelectPublicationMember, publication, schema, table).Scan(&member); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			if _, err := a.pool.Exec(ctx,
				fmt.Sprintf("ALTER PUBLICATION %s ADD TABLE %s", publicationName, qualified),
			); err != nil {
				return storageFailure(ctx, a.log, err)
			}
			a.log.Info(ctx, "tableactivation: Tabelle zur Publication hinzugefügt", "publication", publication, "table", qualified)
			return nil
		}
		return storageFailure(ctx, a.log, err)
	}
	return nil
}

// Unpublish entzieht die Tabelle der Publication; eine fehlende
// Publication und eine fehlende Mitgliedschaft bleiben ohne Wirkung
// (Idempotenz, `LH-FA-CFG-002` Boundary).
func (a *TableActivationAdapter) Unpublish(ctx context.Context, publication, schema, table string) error {
	publicationName, qualified, err := publicationIdentifiers(publication, schema, table)
	if err != nil {
		return err
	}
	var exists int
	if err := a.pool.QueryRow(ctx, queries.SelectPublication, publication).Scan(&exists); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return storageFailure(ctx, a.log, err)
	}
	var member int
	if err := a.pool.QueryRow(ctx, queries.SelectPublicationMember, publication, schema, table).Scan(&member); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return storageFailure(ctx, a.log, err)
	}
	if _, err := a.pool.Exec(ctx,
		fmt.Sprintf("ALTER PUBLICATION %s DROP TABLE %s", publicationName, qualified),
	); err != nil {
		return storageFailure(ctx, a.log, err)
	}
	return nil
}

// Published liest die Mitgliedschaft der Tabelle in der Publication; eine
// fehlende Publication liest als Abwesenheit der Mitgliedschaft — die
// Status- und Listen-Abfragen trennen darüber den Erfassungs-Zustand von
// der Herkunft (`LH-FA-CFG-003`, `LH-FA-CFG-004`).
func (a *TableActivationAdapter) Published(ctx context.Context, publication, schema, table string) (bool, error) {
	if err := validateIdentifier(publication); err != nil {
		return false, err
	}
	if err := validateIdentifier(schema); err != nil {
		return false, err
	}
	if err := validateIdentifier(table); err != nil {
		return false, err
	}
	var member int
	if err := a.pool.QueryRow(ctx, queries.SelectPublicationMember, publication, schema, table).Scan(&member); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, storageFailure(ctx, a.log, err)
	}
	return true, nil
}

// publicationIdentifiers trägt die Bezeichner der Publication-DDL: die
// Prüfung läuft über das Bezeichner-Alphabet, die Interpolation trägt die
// geprüften Literale — die DDL trägt keine Parameter.
func publicationIdentifiers(publication, schema, table string) (string, string, error) {
	for _, name := range []string{publication, schema, table} {
		if err := validateIdentifier(name); err != nil {
			return "", "", err
		}
	}
	return publication, schema + "." + table, nil
}

// validateIdentifier meldet Bezeichner außerhalb des Alphabets über die
// Klasse `configuration` (ErrActivationConfiguration).
func validateIdentifier(name string) error {
	if !identifierShape.MatchString(name) {
		return fmt.Errorf("%w: Bezeichner %q trägt nicht das Bezeichner-Alphabet", ErrActivationConfiguration, name)
	}
	return nil
}
