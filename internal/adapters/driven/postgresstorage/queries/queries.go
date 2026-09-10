// Package queries trägt die SQL-Texte der postgresstorage-Adapter
// (Paketstruktur je `ADR-0042`, der die Struktur-Regeln des abgelösten
// `ADR-0039` als Rest fortgilt): die Tabellennamen nach `SPEC-001` stehen
// hier und nirgends sonst im Adapter; die Zeilen-Übersetzung trägt der
// Mapper.
package queries

// InsertTransaction persistiert eine Quelltransaktion. Die
// Deduplizierungsbasis der Idempotenz (`ADR-0011`) ist der
// Primärschlüssel `transaction_id`: die erneut persistierte Transaktion
// konfligiert und bleibt ohne Wirkung — die Rückkehr meldet keinen Fehler.
// committed_at trägt die Instanzzeit (DEFAULT), der Store schreibt sie
// nicht.
const InsertTransaction = `
INSERT INTO cdc.transaction (transaction_id, source_id, commit_position)
VALUES ($1, $2, $3)
ON CONFLICT (transaction_id) DO NOTHING`

// InsertChange persistiert einen Change; die Deduplizierungsbasis ist der
// Primärschlüssel `change_id` (`SPEC-002`). Die Row Images gehen als Text
// in die `jsonb`-Spalten; ein fehlendes Bild geht als NULL
// (Abwesenheit, `LH-FA-CAP-008` Boundary).
const InsertChange = `
INSERT INTO cdc.change
    (change_id, transaction_id, source_table_id, sequence, operation, old_data, new_data, schema_version)
VALUES ($1, $2, $3, $4, $5, $6::jsonb, $7::jsonb, $8)
ON CONFLICT (change_id) DO NOTHING`

// SelectChanges liest deterministisch sortiert (`LH-FA-REA-004.a`):
// Sortierung nach (Commit-Position der Quelltransaktion, Transaktions-ID,
// Sequenz innerhalb der Transaktion). Der Start ist inklusive, das Ende
// exklusiv (`LH-FA-REA-001`); NULL-Grenzen grenzen nicht ein, LIMIT NULL
// liest unbegrenzt. Lesen trägt nur SELECT — gespeicherte Positionen
// bleiben unverändert (`LH-FA-REA-002`).
const SelectChanges = `
SELECT
    t.source_id,
    t.commit_position,
    c.change_id,
    c.transaction_id,
    c.source_table_id,
    c.sequence,
    c.operation,
    c.old_data,
    c.new_data,
    c.schema_version
FROM cdc.change AS c
JOIN cdc.transaction AS t
    ON c.transaction_id = t.transaction_id
WHERE t.source_id = $1
  AND ($2::bigint IS NULL OR t.commit_position >= $2)
  AND ($3::bigint IS NULL OR t.commit_position < $3)
  AND ($4::text IS NULL OR c.source_table_id = $4)
ORDER BY t.commit_position, c.transaction_id, c.sequence
LIMIT $5`

// CountTransactions zählt die Transaktions-Zeilen einer Quelle; die Tests
// tragen den Deduplizierungs-Stand darüber, nicht der Adapter.
const CountTransactions = `
SELECT count(*) FROM cdc.transaction WHERE source_id = $1`

// CountChanges zählt die Change-Zeilen einer Quelle über die
// Transaktions-Kennung; dieselbe Verwendung wie CountTransactions.
const CountChanges = `
SELECT count(*) FROM cdc.change WHERE transaction_id = $1`

// InsertSourceTable trägt die Bindungs-Zeile einer Aktivierung
// (`LH-FA-CFG-001`); die Deduplizierung der Idempotenz läuft über
// Primärschlüssel und UNIQUE-Kante (source_id, schema_name, table_name) —
// die erneut aktivierte Tabelle bleibt ohne Wirkung.
const InsertSourceTable = `
INSERT INTO cdc.source_table (source_table_id, source_id, schema_name, table_name)
VALUES ($1, $2, $3, $4)
ON CONFLICT DO NOTHING`

// InsertSchemaVersion trägt die Schema-Version-Zeile einer Aktivierung
// (`LH-FA-CFG-001`, `SPEC-004`); dieselbe Idempotenz über den
// Primärschlüssel.
const InsertSchemaVersion = `
INSERT INTO cdc.schema_version (schema_version_id, source_table_id, version)
VALUES ($1, $2, $3)
ON CONFLICT DO NOTHING`

// SelectSourceTable liest die Bindungs-Zeile einer Tabelle — der Zustand
// „aktiviert" (`LH-FA-CFG-003`).
const SelectSourceTable = `
SELECT source_table_id FROM cdc.source_table
WHERE source_id = $1 AND schema_name = $2 AND table_name = $3`

// SelectSourceTableChanges liest den Change-Bestand der Tabelle; der
// Bindungs-Zeilen-Entzug liest darüber die Herkunft der persistierten
// Changes (`LH-FA-CFG-002` Out-of-Scope) und entzieht bei Bestand nichts.
const SelectSourceTableChanges = `
SELECT EXISTS (SELECT 1 FROM cdc.change WHERE source_table_id = $1)`

// DeleteSchemaVersions entzieht die Schema-Version-Zeilen der Tabelle; der
// Entzug läuft nach der Change-Bestands-Prüfung (SelectSourceTableChanges)
// und vor dem Bindungs-Zeilen-Entzug.
const DeleteSchemaVersions = `
DELETE FROM cdc.schema_version WHERE source_table_id = $1`

// DeleteSourceTable entzieht die Bindungs-Zeile; die
// Fremdschlüssel-Kante der Changes schützt den Bestand — der Entzug läuft
// nach der Change-Bestands-Prüfung (SelectSourceTableChanges).
const DeleteSourceTable = `
DELETE FROM cdc.source_table WHERE source_table_id = $1`

// SelectSourceTables liest die Bindungs-Zeilen einer Quelle — die Liste
// der aktivierten Tabellen (`LH-FA-CFG-004`).
const SelectSourceTables = `
SELECT source_table_id, source_id, schema_name, table_name
FROM cdc.source_table
WHERE source_id = $1
ORDER BY schema_name, table_name`

// SelectPublication liest den Publication-Bestand an der Quelle.
const SelectPublication = `
SELECT 1 FROM pg_publication WHERE pubname = $1`

// SelectPublicationMember liest die Mitgliedschaft einer Tabelle in der
// Publication.
const SelectPublicationMember = `
SELECT 1 FROM pg_publication_tables
WHERE pubname = $1 AND schemaname = $2 AND tablename = $3`

// InsertConsumer trägt die Consumer-Zeile einer Registrierung
// (`LH-FA-CON-001`); die Deduplizierung der Idempotenz läuft über den
// Primärschlüssel — die erneut registrierte Kennung bleibt ohne Wirkung
// und die betroffene Zeilen-Zahl liest den Ausgang.
const InsertConsumer = `
INSERT INTO cdc.consumer (consumer_id, name)
VALUES ($1, $2)
ON CONFLICT (consumer_id) DO NOTHING`

// SelectConsumer liest die Consumer-Zeile; die Bestätigung prüft die
// Registrierung über sie vor dem Schreiben — die Abwesenheit endet über
// den benannten Sentinel am Port, nicht über den Fremdschlüssel (`SPEC-001`).
const SelectConsumer = `
SELECT 1 FROM cdc.consumer WHERE consumer_id = $1`

// SelectConsumerPosition liest die bestätigte Position eines Consumers
// (`LH-FA-CON-003`); das Lesen trägt keine Schreibwirkung — der
// gesperrte Lese der Bestätigung läuft als eigene Abfrage
// (SelectConsumerPositionLocked).
const SelectConsumerPosition = `
SELECT source_id, acknowledged_position FROM cdc.consumer_position
WHERE consumer_id = $1`

// SelectConsumerPositionLocked liest und sperrt die bestätigte Position
// innerhalb des Bestätigungs-Commits; die Zeilen-Sperre hält die
// konkurrierenden Bestätigungen desselben Consumers in Ordnung — der
// Monotonie-Vergleich läuft gegen den gesperrten Stand (`LH-FA-CON-002`
// trennt die Consumer, ihr Fortschritt kollidiert nicht).
const SelectConsumerPositionLocked = `
SELECT source_id, acknowledged_position FROM cdc.consumer_position
WHERE consumer_id = $1
FOR UPDATE`

// UpsertConsumerPosition trägt die bestätigte Position fort
// (`LH-FA-CON-003` Boundary); die Monotonie trägt der Domänen-Vergleich
// vor dem Schreiben (`ADR-0029`, Regel 2), nicht der SQL-Ausdruck — die
// gesperrte Zeile liest der Adapter vor diesem Upsert.
const UpsertConsumerPosition = `
INSERT INTO cdc.consumer_position (consumer_id, source_id, acknowledged_position)
VALUES ($1, $2, $3)
ON CONFLICT (consumer_id) DO UPDATE
SET source_id = EXCLUDED.source_id, acknowledged_position = EXCLUDED.acknowledged_position`

// DeleteConsumerPosition entzieht die bestätigte Position; der Entzug
// läuft vor dem Consumer-Zeilen-Entzug — der Fremdschlüssel der
// Positions-Zeile setzt die Consumer-Zeile voraus (`LH-FA-CON-006`
// Boundary: die Position geht mit der Entfernung).
const DeleteConsumerPosition = `
DELETE FROM cdc.consumer_position WHERE consumer_id = $1`

// DeleteConsumer entzieht die Consumer-Zeile; der Entzug läuft nach dem
// Positions-Zeilen-Entzug (DeleteConsumerPosition).
const DeleteConsumer = `
DELETE FROM cdc.consumer WHERE consumer_id = $1`
