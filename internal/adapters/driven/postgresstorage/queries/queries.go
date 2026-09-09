// Package queries trägt die SQL-Texte des PostgresChangeStoreAdapters
// (`ADR-0039`): die Tabellennamen nach `SPEC-001` stehen hier und nirgends
// sonst im Adapter; die Zeilen-Übersetzung trägt der Mapper.
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
