package sqlexec

import (
	"context"
	"time"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage/mapper"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Die Funktionen dieser Datei sind die Zeilen-Übersetzung der Adapter: sie
// setzen eine Anweisung über die Naht ab, laufen die Ergebnis-Zeilen ab und
// übersetzen sie über den Mapper (`ADR-0039`) in Domänen-Werte. Ein
// Treiber-Fehler (Abfrage, Scan, Iteration) läuft durch `Statement.fail` —
// die Klasse des Aufrufers; ein Domänen-Fehler der Übersetzung kommt
// unverändert zurück (`ADR-0029`). Sie sind damit netzlos prüfbar: der
// Träger ist ein `Executor`, kein `*pgxpool.Pool`.

// ReadChanges setzt die Change-Abfrage ab (`SPEC-001`) und trägt die
// Ergebnis-Zeilen in ChangeRecords: die Commit-Position kommt von der
// Transaktion, der Change aus seiner Zeile samt den Klartext-Bezeichnern
// der Tabelle aus dem Join. Die Reihenfolge der Rückgabe trägt die
// SQL-Sortierung (`LH-FA-REA-004.a`).
func ReadChanges(ctx context.Context, exec Executor, statement Statement) ([]outbound.ChangeRecord, error) {
	rows, err := exec.Query(ctx, statement.SQL, statement.Args...)
	if err != nil {
		return nil, statement.fail(err)
	}
	defer rows.Close()

	records := make([]outbound.ChangeRecord, 0)
	for rows.Next() {
		var row mapper.ChangeRow
		var source string
		var commitPosition int64
		var committedAt time.Time
		if err := rows.Scan(
			&source,
			&commitPosition,
			&row.ChangeID,
			&row.TransactionID,
			&row.SourceTableID,
			&row.Schema,
			&row.Table,
			&row.Sequence,
			&row.Operation,
			&row.OldData,
			&row.NewData,
			&row.SchemaVersion,
			&committedAt,
		); err != nil {
			return nil, statement.fail(err)
		}
		position, err := mapper.ToPosition(source, commitPosition)
		if err != nil {
			return nil, err
		}
		change, err := mapper.ToChange(row)
		if err != nil {
			return nil, err
		}
		records = append(records, outbound.ChangeRecord{
			Position:    position,
			Change:      change,
			CommittedAt: model.NewTimePoint(committedAt.UnixNano()),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, statement.fail(err)
	}
	return records, nil
}

// ReadConsumerPosition liest die bestätigte Position eines Consumers
// (`LH-FA-CON-003` Happy Path); die Abwesenheit der Zeile liest sich als
// Nullwert — die definierte Anfangsposition eines Consumers ohne
// Bestätigung (`LH-FA-CON-005` Boundary). Die Zeile läuft zurück durch den
// Domänen-Konstruktor.
func ReadConsumerPosition(ctx context.Context, exec Executor, consumer model.ConsumerID, statement Statement) (model.ConsumerPosition, error) {
	var source string
	var offset int64
	if err := exec.QueryRow(ctx, statement.SQL, statement.Args...).Scan(&source, &offset); err != nil {
		if IsAbsent(err) {
			return model.ConsumerPosition{}, nil
		}
		return model.ConsumerPosition{}, statement.fail(err)
	}
	position, err := mapper.ToPosition(source, offset)
	if err != nil {
		return model.ConsumerPosition{}, err
	}
	return model.ConsumerPosition{ConsumerID: consumer, Position: position}, nil
}

// ReadConsumerPositions liest die bestätigten Positionen aller Consumer einer
// Quelle (`LH-FA-RET-004`): je Zeile in `cdc.consumer_position` ein Eintrag —
// ein Consumer ohne Zeile trägt keine Bestätigung gegen diese Quelle und
// erscheint nicht in der Rückgabe (dieselbe Abwesenheits-Lesart wie
// `ReadConsumerPosition`).
func ReadConsumerPositions(ctx context.Context, exec Executor, source model.SourceID, statement Statement) ([]model.ConsumerPosition, error) {
	rows, err := exec.Query(ctx, statement.SQL, statement.Args...)
	if err != nil {
		return nil, statement.fail(err)
	}
	defer rows.Close()

	positions := make([]model.ConsumerPosition, 0)
	for rows.Next() {
		var consumerID string
		var offset int64
		if err := rows.Scan(&consumerID, &offset); err != nil {
			return nil, statement.fail(err)
		}
		acknowledged, err := mapper.ToPosition(string(source), offset)
		if err != nil {
			return nil, err
		}
		positions = append(positions, model.ConsumerPosition{
			ConsumerID: model.ConsumerID(consumerID),
			Position:   acknowledged,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, statement.fail(err)
	}
	return positions, nil
}

// ReadSourceTables liest die Bindungs-Zeilen einer Quelle in Schema- und
// Tabellen-Ordnung (`LH-FA-CFG-004`); jede Zeile läuft durch den
// Domänen-Konstruktor (`ADR-0029`).
func ReadSourceTables(ctx context.Context, exec Executor, statement Statement) ([]model.SourceTable, error) {
	rows, err := exec.Query(ctx, statement.SQL, statement.Args...)
	if err != nil {
		return nil, statement.fail(err)
	}
	defer rows.Close()

	tables := make([]model.SourceTable, 0)
	for rows.Next() {
		var id, source, schema, table string
		if err := rows.Scan(&id, &source, &schema, &table); err != nil {
			return nil, statement.fail(err)
		}
		entry, err := model.NewSourceTable(model.SourceTableID(id), model.SourceID(source), schema, table)
		if err != nil {
			return nil, err
		}
		tables = append(tables, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, statement.fail(err)
	}
	return tables, nil
}

// ReadExcludedColumns liest den dauerhaften Ausschlussstand je Tabelle einer
// Quelle (`LH-FA-CFG-005`, `ADR-0065`): die `applied`-Zeilen der beiden
// Spalten-Antragsarten, in der Ordnung der Abfrage (Antrags-Zeitpunkt mit
// der Antrags-ID als Zweitschlüssel). `exclude_column` trägt den
// Spaltennamen ein, `include_column` nimmt ihn wieder heraus — derselbe
// Schreibpfad wie der Live-Reload (`ADR-0059` Teilfrage 5), nur über die
// dauerhafte Herkunft statt über den Prozessspeicher. Eine nicht vermerkte
// (`pending`/`failed`) Zeile und eine Zeile einer der beiden
// Tabellen-Antragsarten tragen keinen Stand; eine Quelle ohne
// Spalten-Anträge liefert eine leere Map.
func ReadExcludedColumns(ctx context.Context, exec Executor, statement Statement) (map[string][]string, error) {
	rows, err := exec.Query(ctx, statement.SQL, statement.Args...)
	if err != nil {
		return nil, statement.fail(err)
	}
	defer rows.Close()

	excluded := make(map[string][]string)
	for rows.Next() {
		var schema, table, kind, column string
		if err := rows.Scan(&schema, &table, &kind, &column); err != nil {
			return nil, statement.fail(err)
		}
		qualified := schema + "." + table
		switch model.AdministrationRequestKind(kind) {
		case model.AdministrationRequestExcludeColumn:
			if !containsColumn(excluded[qualified], column) {
				excluded[qualified] = append(excluded[qualified], column)
			}
		case model.AdministrationRequestIncludeColumn:
			// Ein Einschluss, der den letzten geführten Namen nimmt, lässt
			// keinen Eintrag stehen: ein Stand ohne Ausschluss ist ein
			// fehlender Eintrag, keine leere Liste.
			if remaining := removeColumn(excluded[qualified], column); len(remaining) == 0 {
				delete(excluded, qualified)
			} else {
				excluded[qualified] = remaining
			}
		}
	}
	if err := rows.Err(); err != nil {
		return nil, statement.fail(err)
	}
	return excluded, nil
}

// containsColumn meldet, ob ein Spaltenname in einem Ausschlussstand steht;
// der Stand einer Tabelle bleibt klein, die lineare Suche damit ohne
// eigenen Index.
func containsColumn(columns []string, name string) bool {
	for _, column := range columns {
		if column == name {
			return true
		}
	}
	return false
}

// removeColumn liefert den Ausschlussstand ohne den übergebenen
// Spaltennamen; die Rückgabe ist neu aufgebaut, damit kein Leser auf dem
// Speicher der übergebenen Liste liegt.
func removeColumn(columns []string, name string) []string {
	remaining := make([]string, 0, len(columns))
	for _, column := range columns {
		if column != name {
			remaining = append(remaining, column)
		}
	}
	return remaining
}

// ReadTableSchema liest die Spaltenform einer Schema-Version in
// Spalten-Reihenfolge; eine Version ohne Spaltenform endet über
// `outbound.ErrSchemaVersionUnknown` sichtbar (`ADR-0015` Folgepflicht —
// keine stille Fehlinterpretation, `LH-FA-SCH-004`).
func ReadTableSchema(ctx context.Context, exec Executor, versionID model.SchemaVersionID, statement Statement) (model.TableSchema, error) {
	rows, err := exec.Query(ctx, statement.SQL, statement.Args...)
	if err != nil {
		return model.TableSchema{}, statement.fail(err)
	}
	defer rows.Close()

	var columns []model.Column
	for rows.Next() {
		var name string
		var oid int64
		if err := rows.Scan(&name, &oid); err != nil {
			return model.TableSchema{}, statement.fail(err)
		}
		columns = append(columns, model.Column{Name: name, OID: model.ColumnOID(oid)})
	}
	if err := rows.Err(); err != nil {
		return model.TableSchema{}, statement.fail(err)
	}
	if len(columns) == 0 {
		return model.TableSchema{}, outbound.ErrSchemaVersionUnknown
	}
	return model.NewTableSchema(versionID, columns)
}

// ReadPendingRequests liest die offenen Anträge in Anlage-Reihenfolge; jede
// Zeile läuft durch den Domänen-Konstruktor (`ADR-0029`) — eine Zeile
// außerhalb der Antrags-Invarianten endet sichtbar, nicht als still
// gefälschter Antrag.
func ReadPendingRequests(ctx context.Context, exec Executor, statement Statement) ([]model.AdministrationRequest, error) {
	rows, err := exec.Query(ctx, statement.SQL, statement.Args...)
	if err != nil {
		return nil, statement.fail(err)
	}
	defer rows.Close()

	requests := make([]model.AdministrationRequest, 0)
	for rows.Next() {
		var id, source, schema, table, column, kind string
		if err := rows.Scan(&id, &source, &schema, &table, &column, &kind); err != nil {
			return nil, statement.fail(err)
		}
		request, err := model.NewAdministrationRequest(
			model.AdministrationRequestID(id), model.SourceID(source), schema, table, column, model.AdministrationRequestKind(kind),
		)
		if err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}
	if err := rows.Err(); err != nil {
		return nil, statement.fail(err)
	}
	return requests, nil
}

// RegisterConsumer trägt die Consumer-Zeile ein und meldet den Ausgang:
// genau eine betroffene Zeile heißt „registriert", keine heißt „war schon
// da" — die Idempotenz trägt der Primärschlüssel über ON CONFLICT DO
// NOTHING (`LH-FA-CON-001` Boundary), kein Fehler.
func RegisterConsumer(ctx context.Context, exec Executor, statement Statement) (bool, error) {
	tag, err := exec.Exec(ctx, statement.SQL, statement.Args...)
	if err != nil {
		return false, statement.fail(err)
	}
	return tag.RowsAffected() == 1, nil
}
