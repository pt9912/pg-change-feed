-- CDC-Schema der Referenzimplementierung: die Tabellen aus
-- `spec/pflichtenheft.md` §2 (SPEC-001). Angewendet wird die DDL
-- über `postgresstorage.ApplySchema` (schema.go); sie ist idempotent
-- (IF NOT EXISTS) und läuft gegen eine frische Instanz wie gegen eine mit
-- Bestand.
--
-- Umfang: die Store-Seite — `cdc.source`, `cdc.source_table`,
-- `cdc.schema_version` als referenzierte Tabellen (Fremdschlüssel-Ziele der
-- Changes), `cdc.transaction` und `cdc.change` als Träger der Transaktionen
-- und Changes. `cdc.consumer`, `cdc.consumer_position` (Consumer-State)
-- und `cdc.process_heartbeat` (Betriebs-Lebenszeichen) trägt die DDL nicht —
-- ihre Ports (`ConsumerStatePort`, `HeartbeatPort`) liegen außerhalb
-- dieses Store-Adapters.
--
-- Die Fremdschlüssel setzen voraus, dass Quelle, Tabelle und Schema-Version
-- vor der ersten Persistenz registriert sind; das ist Aufgabe der
-- Metadaten-Registrierung, nicht des Store-Schreibpfads.

CREATE SCHEMA IF NOT EXISTS cdc;

-- Erfasste Quelle (`SPEC-001`, Tabelle `cdc.source`).
CREATE TABLE IF NOT EXISTS cdc.source (
    source_id  text PRIMARY KEY,
    name       text NOT NULL
);

-- Aktivierte Tabelle je Quelle (`SPEC-001`, Tabelle `cdc.source_table`);
-- Schema- und Tabellenname machen gleichnamige Tabellen unterscheidbar.
CREATE TABLE IF NOT EXISTS cdc.source_table (
    source_table_id text PRIMARY KEY,
    source_id       text NOT NULL REFERENCES cdc.source (source_id),
    schema_name     text NOT NULL,
    table_name      text NOT NULL,
    UNIQUE (source_id, schema_name, table_name)
);

-- Schema-Version (`SPEC-004`); jeder Change referenziert eine Version.
CREATE TABLE IF NOT EXISTS cdc.schema_version (
    schema_version_id text PRIMARY KEY,
    source_table_id   text NOT NULL REFERENCES cdc.source_table (source_table_id),
    version           bigint NOT NULL CHECK (version >= 1)
);

-- Persistierte Quelltransaktion (`SPEC-001`, Tabelle
-- `cdc.transaction`): der Primärschlüssel trägt die interne
-- Transaktions-ID und ist damit die Deduplizierungsbasis der
-- Idempotenz — dieselbe Transaktion erneut persistiert ändert keinen
-- Stand. Die Commit-Position trägt die sortierbare
-- Quellposition (der Adapter mappt die
-- PostgreSQL-LSN auf den Offset).
CREATE TABLE IF NOT EXISTS cdc.transaction (
    transaction_id  text PRIMARY KEY,
    source_id       text NOT NULL REFERENCES cdc.source (source_id),
    commit_position bigint NOT NULL CHECK (commit_position > 0),
    committed_at    timestamptz NOT NULL DEFAULT now()
);

-- Index über der Lese-Ordnung (`LH-FA-REA-004.a`): Bereich, Limit und
-- Filter lesen über Quelle und Commit-Position.
CREATE INDEX IF NOT EXISTS cdc_transaction_source_position_idx
    ON cdc.transaction (source_id, commit_position);

-- Einzelner Change (`SPEC-002`, Tabelle `cdc.change`): Row Images als
-- `jsonb`, ein fehlendes Bild ist NULL (Abwesenheit, kein Fehler).
-- Die Sequenz ist innerhalb der Transaktion
-- eindeutig und mindestens 1; die UNIQUE-Kante erzwingt die
-- Domänen-Invariante auch gegen Bestand außerhalb
-- der Domänen-Konstruktoren. `origin` trägt die Herkunft (`SPEC-002`):
-- nullable, ohne DEFAULT und ohne CHECK — NULL liest
-- als `wal`, die geschlossene Menge erzwingt die Domäne
-- (`model.ChangeOrigin`); dieselbe Spaltenform wie `tools/schema/schema.yaml`.
-- `route_target` trägt das Zustellziel: nullable, ohne DEFAULT und ohne
-- CHECK — NULL ist „kein Ziel", den Zielnamen erzwingt die Domäne
-- (`model.RouteTarget`).
CREATE TABLE IF NOT EXISTS cdc.change (
    change_id       text PRIMARY KEY,
    transaction_id  text NOT NULL REFERENCES cdc.transaction (transaction_id),
    source_table_id text NOT NULL REFERENCES cdc.source_table (source_table_id),
    sequence        bigint NOT NULL CHECK (sequence >= 1),
    operation       text NOT NULL CHECK (operation IN ('INSERT', 'UPDATE', 'DELETE')),
    old_data        jsonb,
    new_data        jsonb,
    schema_version  text NOT NULL REFERENCES cdc.schema_version (schema_version_id),
    origin          text,
    route_target    text,
    UNIQUE (transaction_id, sequence)
);

-- Index über der Lese-Ordnung (`LH-FA-REA-004.a`): die Sortierung läuft
-- über (Commit-Position, Transaktions-ID, Sequenz) — Transaktion und
-- Sequenz tragen die UNIQUE-Kante oben, die Commit-Position kommt über
-- `cdc.transaction_source_position_idx` in den Join.