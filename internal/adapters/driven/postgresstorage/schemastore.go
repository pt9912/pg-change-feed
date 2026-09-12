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
// trägt ausschließlich die Persistenz-Fähigkeit: weder die dynamische
// Re-Versionierung im laufenden Erfassungspfad noch die
// Typ-Kompatibilitätsprüfung gehören zu ihm. `log` trägt die strukturierte
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
// statische Erstaktivierung schreibt ihre Schema-Version-Zeile über
// `TableActivationPort.Register`, nicht über diesen Port; sie beschreibt
// aber dieselbe Tabelle `cdc.schema_version`, und macht ihre Version 1
// dadurch bereits ohne einen `RegisterVersion`-Aufruf über
// `CurrentVersion` sichtbar).
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
// Primärschlüssel, unabhängig davon, ob die Zeile hier oder dort zuerst
// entsteht. Die Spaltenform geht als eine Zeile je Spalte
// (`queries.InsertTableSchemaColumn`), aber nur, wenn diese
// `SchemaVersionID` noch keine Spaltenform trägt (`CountTableSchemaColumns`)
// — das trägt sowohl die Erstregistrierung als auch das Nachtragen der
// Spaltenform zu einer `SchemaVersionID`, deren Schema-Version-Zeile
// bereits über einen anderen Schreibpfad besteht (Backfill); eine
// `SchemaVersionID` mit bestehender Spaltenform schreibt keine zweite.
// Eine widersprüchliche Kennung zwischen Version und TableSchema endet
// vor dem ersten SQL-Aufruf über `outbound.ErrSchemaVersionMismatch`.
func (a *PostgresSchemaStoreAdapter) RegisterVersion(ctx context.Context, version model.SchemaVersion, schema model.TableSchema) (bool, error) {
	if version.ID != schema.VersionID {
		return false, outbound.ErrSchemaVersionMismatch
	}
	tx, err := a.pool.Begin(ctx)
	if err != nil {
		return false, schemaStoreFailure(ctx, a.log, err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, queries.InsertSchemaVersion,
		string(version.ID), string(version.SourceTableID), version.Version); err != nil {
		return false, schemaStoreFailure(ctx, a.log, err)
	}

	var existingColumns int
	if err := tx.QueryRow(ctx, queries.CountTableSchemaColumns, string(schema.VersionID)).Scan(&existingColumns); err != nil {
		return false, schemaStoreFailure(ctx, a.log, err)
	}
	written := existingColumns == 0
	if written {
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
	if written {
		a.log.Info(ctx, "schemastore: Schema-Version registriert",
			"schema_version_id", version.ID, "source_table_id", version.SourceTableID, "version", version.Version)
	}
	return written, nil
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
