package postgresstorage

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage/queries"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// PostgresSchemaStoreAdapter ist die Referenzimplementierung des
// `SchemaStorePort` (`ARC-004`, `ADR-0015` Folgepflicht): er persistiert
// `TableSchema`-/`SchemaVersion`-Modelle je Change (`SPEC-004`) gegen die
// Tabelle `cdc.table_schema`, ausgerollt über d-migrate
// (`tools/schema/schema.yaml`, `ADR-0043`) — die handgeschriebene DDL des
// Store-Adapters (`schema.sql`) trägt sie nicht, dieselbe Abgrenzung wie
// bei den Consumer-State-Tabellen (`consumerstate.go`). Dieser Adapter
// trägt ausschließlich die Persistenz-Fähigkeit dieses Slice — die
// dynamische Re-Versionierung im laufenden Erfassungspfad und die
// Typ-Kompatibilitätsprüfung sind Folgepflichten anderer Slices
// (`slice-032`, `slice-033`). `log` trägt die strukturierte
// Protokollierung über den injizierten `LogPort` (`LH-QA-OPS-004`,
// `ADR-0024`, `WithLog`) — Default `outbound.NoopLog`.
type PostgresSchemaStoreAdapter struct {
	pool *pgxpool.Pool
	log  outbound.LogPort
}

// NewSchemaStore baut den Verbindungspool gegen die Instanz und meldet
// eine nicht erreichbare Instanz als Fehler der Klasse `storage`
// (`outbound.ErrSchemaStoreStorage`, `SPEC-008`).
func NewSchemaStore(ctx context.Context, dsn string, opts ...Option) (*PostgresSchemaStoreAdapter, error) {
	o := newOptions(opts)
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, schemaStoreFailure(ctx, o.log, err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, schemaStoreFailure(ctx, o.log, err)
	}
	o.log.Info(ctx, "schemastore: verbunden")
	return &PostgresSchemaStoreAdapter{pool: pool, log: o.log}, nil
}

// schemaStoreFailure trägt die Übersetzungsverantwortung dieses Adapters
// (`ADR-0023`, `SPEC-008`) — derselbe Aufbau wie `stateStorageFailure`
// (`consumerstate.go`), eigener Sentinel (`outbound.ErrSchemaStoreStorage`):
// die Klasse-Aktion des ChangeStore-Sentinels (kein Source-ACK,
// `LH-QA-REL-001.a`) trägt dieser Port nicht.
func schemaStoreFailure(ctx context.Context, log outbound.LogPort, cause error) error {
	log.Error(ctx, "schemastore: Datenbankfehler", "error", cause)
	return fmt.Errorf("%w: %w", outbound.ErrSchemaStoreStorage, cause)
}

// Close schließt den Verbindungspool.
func (a *PostgresSchemaStoreAdapter) Close() {
	a.pool.Close()
}

var _ outbound.SchemaStorePort = (*PostgresSchemaStoreAdapter)(nil)

// CurrentVersion liest die höchste registrierte Schema-Version einer
// Tabelle; die Abwesenheit meldet, dass die Tabelle noch keine über
// diesen Port registrierte Version trägt (`ADR-0015` Folgepflicht — die
// statische Erstaktivierung schreibt ihre Schema-Version-Zeile weiterhin
// über `TableActivationPort.Register`, nicht über diesen Port; ihre
// Version 1 wird über `CurrentVersion` erst nach der ersten
// `RegisterVersion` sichtbar).
func (a *PostgresSchemaStoreAdapter) CurrentVersion(ctx context.Context, table model.SourceTableID) (model.SchemaVersion, bool, error) {
	if table == "" {
		return model.SchemaVersion{}, false, domainerrors.ErrEmptyIdentifier
	}
	var id string
	var version int64
	err := a.pool.QueryRow(ctx, queries.SelectCurrentSchemaVersion, string(table)).Scan(&id, &version)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.SchemaVersion{}, false, nil
		}
		return model.SchemaVersion{}, false, schemaStoreFailure(ctx, a.log, err)
	}
	current, err := model.NewSchemaVersion(model.SchemaVersionID(id), table, version)
	if err != nil {
		return model.SchemaVersion{}, false, err
	}
	return current, true, nil
}

// RegisterVersion trägt die Schema-Version- und TableSchema-Zeilen in
// EINEM Store-Commit ein (`ADR-0015` Folgepflicht): die
// Schema-Version-Zeile trägt dieselbe Anweisung wie die statische
// Erstaktivierung (`queries.InsertSchemaVersion`,
// `TableActivationAdapter.Register`) — ihre Idempotenz trägt derselbe
// Primärschlüssel. Die Spaltenform geht als eine Zeile je Spalte
// (`queries.InsertTableSchemaColumn`), nur wenn die Version neu ist — eine
// erneut registrierte Version schreibt keine zweite Spaltenform. Eine
// widersprüchliche Kennung zwischen Version und TableSchema endet vor dem
// ersten SQL-Aufruf über `outbound.ErrSchemaVersionMismatch`.
func (a *PostgresSchemaStoreAdapter) RegisterVersion(ctx context.Context, version model.SchemaVersion, schema model.TableSchema) (bool, error) {
	if version.ID != schema.VersionID {
		return false, outbound.ErrSchemaVersionMismatch
	}
	tx, err := a.pool.Begin(ctx)
	if err != nil {
		return false, schemaStoreFailure(ctx, a.log, err)
	}
	defer tx.Rollback(ctx)

	tag, err := tx.Exec(ctx, queries.InsertSchemaVersion,
		string(version.ID), string(version.SourceTableID), version.Version)
	if err != nil {
		return false, schemaStoreFailure(ctx, a.log, err)
	}
	registered := tag.RowsAffected() == 1
	if registered {
		for i, column := range schema.Columns {
			if _, err := tx.Exec(ctx, queries.InsertTableSchemaColumn,
				string(schema.VersionID), int64(i+1), column.Name, int64(column.OID),
			); err != nil {
				return false, schemaStoreFailure(ctx, a.log, err)
			}
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return false, schemaStoreFailure(ctx, a.log, err)
	}
	if registered {
		a.log.Info(ctx, "schemastore: Schema-Version registriert",
			"schema_version_id", version.ID, "source_table_id", version.SourceTableID, "version", version.Version)
	}
	return registered, nil
}

// TableSchema liest die Spaltenform einer Schema-Version in
// Spalten-Reihenfolge; die Abwesenheit endet über
// `outbound.ErrSchemaVersionUnknown` sichtbar (`ADR-0015` Folgepflicht —
// keine stille Fehlinterpretation, `LH-FA-SCH-004`).
func (a *PostgresSchemaStoreAdapter) TableSchema(ctx context.Context, versionID model.SchemaVersionID) (model.TableSchema, error) {
	if versionID == "" {
		return model.TableSchema{}, domainerrors.ErrEmptyIdentifier
	}
	rows, err := a.pool.Query(ctx, queries.SelectTableSchemaColumns, string(versionID))
	if err != nil {
		return model.TableSchema{}, schemaStoreFailure(ctx, a.log, err)
	}
	defer rows.Close()

	var columns []model.Column
	for rows.Next() {
		var name string
		var oid int64
		if err := rows.Scan(&name, &oid); err != nil {
			return model.TableSchema{}, schemaStoreFailure(ctx, a.log, err)
		}
		columns = append(columns, model.Column{Name: name, OID: model.ColumnOID(oid)})
	}
	if err := rows.Err(); err != nil {
		return model.TableSchema{}, schemaStoreFailure(ctx, a.log, err)
	}
	if len(columns) == 0 {
		return model.TableSchema{}, outbound.ErrSchemaVersionUnknown
	}
	return model.NewTableSchema(versionID, columns)
}
