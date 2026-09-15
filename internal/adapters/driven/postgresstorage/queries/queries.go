// Package queries trägt die SQL-Texte der postgresstorage-Adapter
// (Paketstruktur je `ADR-0042`, der die Struktur-Regeln des abgelösten
// `ADR-0039` als Rest fortgilt): die Tabellennamen nach `SPEC-001` stehen
// hier und nirgends sonst im Adapter; die Zeilen-Übersetzung trägt der
// Mapper.
package queries

// InsertTransaction persistiert eine Quelltransaktion. Die
// Deduplizierungsbasis der Idempotenz (`ADR-0011`) ist der
// Primärschlüssel `transaction_id`: die erneut persistierte Transaktion
// konfligiert und bleibt ohne Wirkung — die Rückkehr meldet keinen Fehler,
// und der zuerst geschriebene committed_at-Wert bleibt bestehen.
// committed_at trägt den realen Quell-Commit-Zeitpunkt
// (`LH-FA-ADM-004`) — der Store übergibt ihn explizit; die Spalten-DEFAULT
// (`current_timestamp`) greift nur außerhalb dieses Anwendungspfads.
const InsertTransaction = `
INSERT INTO cdc.transaction (transaction_id, source_id, commit_position, committed_at)
VALUES ($1, $2, $3, $4)
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
// liest unbegrenzt. Der Tabellenfilter läuft über die Klartext-Bezeichner
// der Bindungs-Zeile (`st.schema_name`/`st.table_name`), je optional und
// unabhängig; der Join auf `cdc.source_table` trägt dieselben Bezeichner
// in die Projektion, damit die Rückgabe die Tabellen-Identität in Klartext
// führt — dieselbe Projektion wie die View `cdc.changes`. Lesen trägt nur
// SELECT — gespeicherte Positionen bleiben unverändert (`LH-FA-REA-002`).
// committed_at trägt den realen Quell-Commit-Zeitpunkt der Transaktion
// (`LH-FA-ADM-004`) — die zeitbasierte Retention (`LH-FA-RET-003`) liest
// ihr Alter dagegen.
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
    c.schema_version,
    st.schema_name,
    st.table_name,
    t.committed_at
FROM cdc.change AS c
JOIN cdc.transaction AS t
    ON c.transaction_id = t.transaction_id
JOIN cdc.source_table AS st
    ON c.source_table_id = st.source_table_id
WHERE t.source_id = $1
  AND ($2::bigint IS NULL OR t.commit_position >= $2)
  AND ($3::bigint IS NULL OR t.commit_position < $3)
  AND ($4::text IS NULL OR st.schema_name = $4)
  AND ($5::text IS NULL OR st.table_name = $5)
ORDER BY t.commit_position, c.transaction_id, c.sequence
LIMIT $6`

// DeleteChanges entfernt genau die übergebenen Change-Zeilen
// (`LH-FA-RET-002`…`004`); die Freigabe je Change trägt der aufrufende Use
// Case über `RetentionPolicy.AllowsDeletion` — diese Abfrage führt nur die
// bereits freigegebene Menge aus. Eine Kennung ohne Zeile bleibt ohne
// Wirkung (Idempotenz). RETURNING liefert die betroffenen
// Transaktions-Kennungen an den Aufrufer zurück — Grundlage für
// DeleteOrphanedTransactions: die Fremdschlüssel-Kante
// (`change.transaction_id → transaction.transaction_id`) kaskadiert nur in
// Richtung Parent-Löschung auf die Children, nie umgekehrt; eine durch
// diese Löschung verwaiste Elternzeile in `cdc.transaction` bräuchte sonst
// keinen weiteren Träger und bliebe stehen.
const DeleteChanges = `
DELETE FROM cdc.change WHERE change_id = ANY($1)
RETURNING transaction_id`

// DeleteOrphanedTransactions entfernt aus der übergebenen Menge genau die
// Transaktions-Zeilen, die keine Change-Zeile mehr referenziert. `cdc.metrics`
// (`tools/schema/nacharbeit-observability.sql`) berechnet
// `cdc_oldest_change_age_seconds`/`cdc_transactions_total` direkt über
// `cdc.transaction`, ungefiltert nach verbliebenen Changes — eine verwaiste
// Zeile läse dort als bestehende Aktivität. Scope bewusst eng auf die von
// DeleteChanges betroffene Menge — eine Transaktion, die von Geburt an nie
// eine Change-Zeile trug, ist ein anderer Fall (keine Rolle der Retention)
// und bleibt unberührt.
const DeleteOrphanedTransactions = `
DELETE FROM cdc.transaction
WHERE transaction_id = ANY($1)
  AND NOT EXISTS (
      SELECT 1 FROM cdc.change c WHERE c.transaction_id = cdc.transaction.transaction_id
  )`

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

// SelectTableColumnExists liest die physische Spalte einer Tabelle über
// den Katalog; der Negative-Pfad des Spaltenausschlusses/-einschlusses
// endet über die Abwesenheit (`LH-FA-CFG-005`) — dieselbe Katalog-Quelle
// wie `TableExists` (information_schema.tables), hier eine Ebene tiefer.
const SelectTableColumnExists = `
SELECT count(*) FROM information_schema.columns WHERE table_schema = $1 AND table_name = $2 AND column_name = $3`

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

// SelectConsumerPositionsBySource liest die bestätigten Positionen aller
// Consumer einer Quelle (`LH-FA-RET-004`): je Zeile ein Consumer, der
// bereits gegen diese Quelle bestätigt hat — ein Consumer ohne Zeile trägt
// keine Bestätigung gegen diese Quelle und blockiert ihre Retention nicht
// (dieselbe Abwesenheits-Lesart wie bei SelectConsumerPosition).
const SelectConsumerPositionsBySource = `
SELECT consumer_id, acknowledged_position FROM cdc.consumer_position
WHERE source_id = $1`

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

// UpsertHeartbeat trägt das Lebenszeichen der Quelle fort
// (`cdc.process_heartbeat`, `LH-FA-ADM-002`): die Instanzzeit der
// Speicherseite trägt den Zeitstempel (`current_timestamp`) — der
// Aufrufer übergibt keine Uhr. Eine bestehende Zeile aktualisiert ihren
// Zeitstempel; je Quelle bleibt genau eine Zeile. Ein erfolgreicher Beat
// löscht einen zuvor gemeldeten Fehlerzustand (`error_class`) wieder —
// der Fehlerzustand endet dadurch selbst erkennbar (`LH-FA-ADM-003`
// Boundary).
const UpsertHeartbeat = `
INSERT INTO cdc.process_heartbeat (source_id, heartbeat_at, error_class)
VALUES ($1, current_timestamp, NULL)
ON CONFLICT (source_id) DO UPDATE SET heartbeat_at = current_timestamp, error_class = NULL`

// InsertTableSchemaColumn persistiert eine Spalten-Zeile einer
// TableSchema-Version (`ADR-0015` Folgepflicht, `SPEC-004`) — eine Zeile
// je Spalte, in Anlage-Reihenfolge über `ordinal_position` sortiert
// (`SelectTableSchemaColumns`); dieselbe Deduplizierungsbasis wie
// `InsertChange`.
const InsertTableSchemaColumn = `
INSERT INTO cdc.table_schema (schema_version_id, ordinal_position, column_name, column_oid)
VALUES ($1, $2, $3, $4)
ON CONFLICT DO NOTHING`

// SelectCurrentSchemaVersion liest die höchste registrierte Schema-Version
// einer Tabelle — die „aktuelle" Version (`ADR-0015` Folgepflicht).
const SelectCurrentSchemaVersion = `
SELECT schema_version_id, version FROM cdc.schema_version
WHERE source_table_id = $1
ORDER BY version DESC
LIMIT 1`

// SelectTableSchemaColumns liest die Spaltenform einer Schema-Version in
// Spalten-Reihenfolge (`ordinal_position`).
const SelectTableSchemaColumns = `
SELECT column_name, column_oid FROM cdc.table_schema
WHERE schema_version_id = $1
ORDER BY ordinal_position`

// CountTableSchemaColumns zählt die Spalten-Zeilen einer Schema-Version —
// die Backfill-Prüfung von `RegisterVersion`: eine Version ohne
// bestehende Spaltenform bekommt sie geschrieben, unabhängig davon, ob
// ihre `schema_version`-Zeile neu ist oder bereits besteht.
const CountTableSchemaColumns = `
SELECT count(*) FROM cdc.table_schema WHERE schema_version_id = $1`

// UpsertHeartbeatFault trägt den zuletzt beobachteten Fehlerzustand der
// Quelle fort (`cdc.process_heartbeat.error_class`, `LH-FA-ADM-003`,
// `LH-QA-REL-003`): dieselbe Zeile wie UpsertHeartbeat,
// derselbe fortlaufende Zeitstempel — ein Fehlerzustand ist ein
// Lebenszeichen mit Klasse, keine zweite Tabelle.
const UpsertHeartbeatFault = `
INSERT INTO cdc.process_heartbeat (source_id, heartbeat_at, error_class)
VALUES ($1, current_timestamp, $2)
ON CONFLICT (source_id) DO UPDATE SET heartbeat_at = current_timestamp, error_class = EXCLUDED.error_class`

// SelectPendingAdministrationRequests liest die offenen Anträge der
// Antrags-Queue (`cdc.administration_request`, `LH-FA-ADM-001`) in
// Anlage-Reihenfolge (`requested_at`) — die Administrations-Goroutine
// verarbeitet sie in dieser Ordnung, sowohl nach `NOTIFY` als auch
// periodisch als Fallback-Poll. Die vier Antragsarten teilen sich eine
// Tabelle; die beiden Tabellen-Antragsarten tragen keine Spalte
// (`column_name` NULL) — `COALESCE` normalisiert das auf den leeren Wert.
const SelectPendingAdministrationRequests = `
SELECT administration_request_id, source_id, schema_name, table_name, COALESCE(column_name, ''), request_kind
FROM cdc.administration_request
WHERE status = 'pending'
ORDER BY requested_at`

// SelectAppliedColumnRequests liest die `applied`-Zeilen der beiden
// Spalten-Antragsarten einer Quelle (`LH-FA-CFG-005`, `ADR-0065`): der
// Adapter wertet sie zur Reihenfolge aus und trägt damit den dauerhaften
// Ausschlussstand. `requested_at` trägt den Transaktionszeitstempel
// (`current_timestamp` der schreibenden Funktion) und ist zwischen zwei
// Anträgen derselben Transaktion nicht unterscheidend — deshalb der
// deterministische Zweitschlüssel `administration_request_id`: dieselbe
// Antrags-Menge trägt damit unabhängig von der Ausführungsreihenfolge
// genau eine Reihenfolge. Die beiden Tabellen-Antragsarten bleiben außen
// vor; `COALESCE` normalisiert das für sie NULL-bare `column_name` wie in
// SelectPendingAdministrationRequests.
const SelectAppliedColumnRequests = `
SELECT schema_name, table_name, request_kind, COALESCE(column_name, '')
FROM cdc.administration_request
WHERE source_id = $1
  AND status = 'applied'
  AND request_kind IN ('exclude_column', 'include_column')
ORDER BY requested_at, administration_request_id`

// UpdateAdministrationRequestApplied vermerkt einen erfolgreich
// verarbeiteten Antrag; die WHERE-Klausel trägt die Idempotenz — ein
// bereits vermerkter Antrag (nicht mehr `pending`) bleibt unverändert und
// die Rückkehr trägt `RowsAffected() == 0`, kein Fehler.
const UpdateAdministrationRequestApplied = `
UPDATE cdc.administration_request
SET status = 'applied', error_message = NULL
WHERE administration_request_id = $1 AND status = 'pending'`

// UpdateAdministrationRequestFailed vermerkt einen gescheiterten Antrag
// samt Fehlertext; dieselbe Idempotenz-Klausel wie
// UpdateAdministrationRequestApplied.
const UpdateAdministrationRequestFailed = `
UPDATE cdc.administration_request
SET status = 'failed', error_message = $2
WHERE administration_request_id = $1 AND status = 'pending'`
