# Pflichtenheft — PG Change Feed

**Bezug zum Lastenheft:** Dieses Pflichtenheft präzisiert die in
`spec/lastenheft.md` formulierten Anforderungen (`LH-*`-IDs). Bei
Konflikt gewinnt das Lastenheft — präzisieren ja, erweitern nie.

**Rolle:** Technik-Stratum — fortschreibbar ohne Change Request; eine ADR darf
sie schärfen, das Lastenheft nicht. Regeln: Baseline-Regelwerk
`modul-03-spec.md` §Ziel-Form: Spezifikation.

---

## 1. Algorithmen und Datenflüsse

Regeln dieser Sektion: Tatsächlicher Code gehört in `src/`, nicht hierher.
ID-Schema `LH-<STRATUM>-<BEREICH>-<NN>.<Buchstabe>` für Verfeinerungen
einzelner Lastenheft-IDs (Baseline-Regelwerk
`grundlagen-source-precedence.md` §ID-Schema als Klammer). Was **keine**
einzelne Lastenheft-ID verfeinert, trägt eine `SPEC-<NNN>` — siehe §2 bis §6.

### LH-QA-REL-001.a — Persist-before-ACK

**Eingabe:** dekodierte `pgoutput`-Nachrichten eines Logical Replication
Streams. **Ausgabe:** dauerhaft gespeicherte Changes; bestätigte
Quellpositionen. **Schritte:**

1. Receive — der Replication Stream ist ein Driving Adapter; er baut die
   Replication-Verbindung auf und empfängt die Nachrichten.
2. Decode — `pgoutput`-Nachrichten werden dekodiert und in Aufrufe des
   CaptureInboundPort übersetzt.
3. Persist — committed Quelltransaktionen werden über den ChangeStorePort
   dauerhaft gespeichert.
4. COMMIT Store — die Persistenz ist abgeschlossen.
5. ACK Source — erst jetzt bestätigt der ReplicationAckPort die Position
   gegenüber PostgreSQL.

Zentrale Invariante:

    ACK(position) => durable(all changes <= position)

**Fehlermodi:** Persistenzfehler → kein Source-ACK; Dekodierfehler →
sichtbarer Fehler (Fehlerklasse `schema`, §4), kein stilles Überspringen.
Ein Crash zwischen Persistenz und ACK darf höchstens zu erneuter
Verarbeitung führen — wiederholte WAL-Daten sind deduplizierbar bzw.
idempotent persistierbar. Nach Crash wird anhand des Replication Slots und
persistierter Zustände fortgesetzt; Wiederholung wird gegenüber möglichem
Datenverlust bevorzugt.

### LH-FA-CAP-006.a — Transaktionsgrenzen und Sichtbarkeit

**Eingabe:** Replication-Stream-Ereignisse. **Ausgabe:** vollständige,
konsumierbare Quelltransaktionen am CaptureInboundPort.

BEGIN, COMMIT und relevante Streaming-Grenzen werden berücksichtigt:

- Offene Transaktionen sind nicht konsumierbar ([`LH-FA-CAP-006`](lastenheft.md)).
- Rollbacks erzeugen keine committed Changes ([`LH-FA-CAP-007`](lastenheft.md)).
- Jede persistierte Quelltransaktion besitzt eine interne ID und
  Commit-Position; Changes besitzen eine eindeutige Sequenz innerhalb der
  Transaktion ([`LH-FA-CAP-004`](lastenheft.md), [`LH-FA-DAT-004`](lastenheft.md)).
- Große Transaktionen dürfen nicht unbegrenzt im RAM gehalten werden;
  der TransactionBufferPort (SPEC-005) ermöglicht Spooling.

### LH-FA-CFG-001.a — Aktivierung einer Tabelle

**Eingabe:** Tabellenname, CDC-Konfiguration. **Ausgabe:** aktivierte
Tabelle; aktualisierte Publication/Slot-Verwaltung.

1. Replica Identity prüfen. `REPLICA IDENTITY FULL` kann für vollständige
   alte Werte erforderlich sein ([`LH-FA-CAP-008`](lastenheft.md)), wird aber **nicht
   ungefragt** aktiviert.
2. Publication und Logical Replication Slot verwalten; als Output Plugin
   wird `pgoutput` verwendet.
3. Operationen erfassen: INSERT, UPDATE, DELETE. TRUNCATE wird erkannt und
   im MVP explizit als nicht unterstützt behandelt (Lastenheft fordert
   TRUNCATE nicht, §5 Out-of-Scope).

### LH-FA-CON-001.a — Registrierung eines Consumers: Zugriffsweg offen

**Eingabe:** Consumer-Name. **Ausgabe:** registrierter Consumer.

Die Registrierungslogik (Ablehnung eines bereits vergebenen Namens,
Erzeugung der Consumer-Kennung) ist als eigenständig testbare Einheit
umgesetzt, aber **ohne einen von außen erreichbaren Zugriffsweg**: Eine
externe Anwendung kann sich heute über keinen unterstützten Kanal (CLI,
Netzwerkschnittstelle) registrieren. Ein direktes Schreiben der
CDC-Speichertabelle (`SPEC-001`) am Anwendungsdienst vorbei ist möglich,
umgeht aber jede künftige Prüfung dieser Logik. Welcher Zugriffsweg den
Aufruf trägt (CLI-Unterbefehl, Netzwerkschnittstelle, ein anderer
Mechanismus), ist eine offene technische Frage.

### LH-FA-CON-004.a — Bestätigung einer Position: Zugriffsweg offen, Vorwärts-Invariante nicht durchgesetzt

**Eingabe:** Consumer-Kennung, zu bestätigende Position. **Ausgabe:**
aktualisierte Consumer-Position.

Derselbe Befund wie LH-FA-CON-001.a: die Bestätigungslogik trägt die
Vorwärts-Invariante (eine Bestätigung bewegt die Position nur vorwärts,
nie zurück) als eigenständig getestete Prüfung, aber ohne von außen
erreichbaren Zugriffsweg. Ein direktes Schreiben der
Consumer-Positions-Tabelle (`SPEC-001`) umgeht diese Invariante
vollständig — sie ist an keiner anderen Stelle (auch nicht im
Datenbankschema selbst) durchgesetzt. Zugriffsweg wie LH-FA-CON-001.a
offen.

### LH-FA-REA-004.a — Deterministische Sortierung

**Eingabe:** Positionsbereich bzw. Startposition, Limit, optionaler
Tabellenfilter. **Ausgabe:** begrenzte, deterministisch sortierte
Ergebnismenge.

Sortierung nach (Commit-Position der Quelltransaktion, Transaktions-ID,
Sequenz innerhalb der Transaktion) — bei identischer Eingabe ist die
Reihenfolge damit stabil ([`LH-FA-REA-004`](lastenheft.md), [`LH-FA-REA-003`](lastenheft.md)). Changes können
ab einer SourcePosition bzw. innerhalb eines Positionsbereichs gelesen
werden; Lesen verändert gespeicherte Positionen nicht.

### LH-FA-RET-004.a — Safe Watermark der Retention

**Eingabe:** Consumer-Positionen, Retention-Konfiguration. **Ausgabe:**
Bereinigungsmenge.

1. Die Safe Watermark basiert auf den relevanten Consumer-Positionen;
   Changes, die aktive Consumer benötigen, dürfen nicht automatisch
   gelöscht werden.
2. Zeitbasierte Mindestaufbewahrung ist zusätzlich konfigurierbar
   ([`LH-FA-RET-003`](lastenheft.md)).
3. Langsame Consumer, die die Safe Watermark blockieren, sind sichtbar
   ([`LH-FA-RET-005`](lastenheft.md)).

Retention ist eine Domain Policy; sie liegt nicht in einem Adapter.

### LH-FA-SCH-004.a — Schema-Änderungs-Klassifikation

**Eingabe:** Relation Metadata aus dem Replication Stream. **Ausgabe:**
TableSchema-/SchemaVersion-Modell (SPEC-004) oder sichtbarer Fehler.

Relation Metadata wird in technologieunabhängige
TableSchema-/SchemaVersion-Modelle übersetzt; jeder Change referenziert
eine Schema-Version ([`LH-FA-SCH-005`](lastenheft.md)). Nicht sicher interpretierbare
Schemaänderungen führen zu einem sichtbaren Fehler (Fehlerklasse `schema`,
§4) statt stiller Fehlinterpretation ([`LH-FA-SCH-004`](lastenheft.md)).

---

## 2. Datenstrukturen und Schemas

Regeln dieser Sektion: Jede Struktur trägt eine `SPEC-<NNN>` — eine Adresse,
keine Anforderung. Gezählt wird fortlaufend je Datei, nicht je Sektion.

### SPEC-001 — CDC-Schema (Tabellen)

Standardmäßig wird das Schema `cdc` verwendet (Default: §3, `SPEC-006`).
Vorgesehene Tabellen:

| Tabelle | Zweck |
|---|---|
| `cdc.source` | erfasste Quellen |
| `cdc.source_table` | aktivierte Tabellen je Quelle |
| `cdc.transaction` | persistierte Quelltransaktionen (interne ID, Commit-Position) |
| `cdc.change` | die einzelnen Changes (SPEC-002) |
| `cdc.consumer` | registrierte Consumer |
| `cdc.consumer_position` | bestätigte Position je Consumer |
| `cdc.schema_version` | Schema-Versionen (SPEC-004) |
| `cdc.process_heartbeat` | Betriebs-/Capture-Zustand: periodisches Lebenszeichen des Capture-Prozesses ([`LH-FA-ADM-002`](lastenheft.md), [`LH-QA-OPS-002`](lastenheft.md)) |

### SPEC-002 — `cdc.change`

```json
{
  "change_id": "eindeutige Change-Kennung (LH-FA-DAT-001)",
  "transaction_id": "interne Transaktions-ID (LH-FA-CAP-005)",
  "source_table_id": "Quelltabelle (LH-FA-DAT-002)",
  "sequence": "eindeutige Sequenz innerhalb der Transaktion",
  "operation": "INSERT | UPDATE | DELETE",
  "old_data": "jsonb | null",
  "new_data": "jsonb | null",
  "schema_version": "Referenz auf cdc.schema_version"
}
```

Im MVP werden `old_data` und `new_data` als `jsonb` gespeichert (Row
Images; Verfügbarkeit je Operationstyp siehe [`LH-FA-CAP-008`](lastenheft.md)).

### SPEC-003 — SourcePosition

Der Core verwendet ein technologieunabhängiges Positionsmodell; der
PostgreSQL-Adapter mappt die PostgreSQL-LSN darauf:

```json
{
  "source_id": "Referenz auf cdc.source",
  "position": "technologieunabhängige sortierbare Position (LH-FA-DAT-004)"
}
```

### SPEC-004 — TableSchema / SchemaVersion

Technologieunabhängige Modelle für Relation Metadata; jeder Change
referenziert eine Schema-Version ([`LH-FA-SCH-005`](lastenheft.md)). Die Übersetzung
erfolgt in LH-FA-SCH-004.a.

### SPEC-005 — TransactionBuffer

Vertrag für das Spooling großer Quelltransaktionen: Zwischenstände einer
noch nicht abgeschlossenen Transaktion werden außerhalb des RAM
gehalten, bis die Transaktion committed und dauerhaft persistiert ist
(vorgesehene Implementierung: FileTransactionBufferAdapter).

### SPEC-016 — Konfigurationsdatei (`CDC_CONFIG_FILE`)

Feldform der optionalen YAML-Konfigurationsdatei: additiv zu den
Umgebungsvariablen, mit Umgebungsvariable-schlägt-Datei-Feld-für-Feld-
Precedence; die drei DSN-Schlüssel (`capture_dsn`/`admin_dsn`/
`reader_dsn`) sind **nicht zulässig** und brechen das Laden über die
Fehlerklasse `configuration` ab
([`LH-QA-SEC-001`](lastenheft.md)/[`LH-QA-SEC-002`](lastenheft.md):
Least-Privilege/Secret-Trennung — eine Konfigurationsdatei ist für andere
Aufbewahrungs-/Verteilwege bestimmt als eine Umgebungsvariable). Striktes
Decoding (unbekannter Schlüssel → Fehlerklasse `configuration`).

| Schlüssel | Typ | Entspricht (Env-Var) | Pflicht in der Datei |
|---|---|---|---|
| `source_id` | string | `CDC_SOURCE_ID` | nein — Pflichtfeld nach Merge (Datei oder Env) |
| `publication` | string | `CDC_PUBLICATION` | nein — Pflichtfeld nach Merge (Datei oder Env) |
| `slot` | string | `CDC_SLOT` | nein — Pflichtfeld nach Merge (Datei oder Env) |
| `tables` | Mapping `<schema.tabelle>: {table_id, schema_version}` | `CDC_TABLES` (Zeichenkettenform, unverändert bestehen) | nein — Pflichtfeld nach Merge (Datei oder Env) |
| `log_level` | string (`debug`/`info`/`warn`/`error`) | `CDC_LOG_LEVEL` | nein, Default `info` |
| `wal_retention_warn_bytes` | int64 | — (kein Env-Gegenstück) | nein, Default SPEC-013 |
| `wal_retention_error_bytes` | int64 | — (kein Env-Gegenstück) | nein, Default SPEC-013 |

```yaml
source_id: quelle-1
publication: pub_quelle_1
slot: slot_quelle_1
tables:
  public.orders:
    table_id: tbl-orders
    schema_version: sv-orders-1
  public.customers:
    table_id: tbl-customers
    schema_version: sv-customers-1
log_level: info
```

`tables` als Feld wird **als Ganzes** ersetzt, nicht Zeile für Zeile
zusammengeführt: eine gesetzte `CDC_TABLES` schlägt die gesamte
Datei-`tables`-Mapping vollständig, ohne Vermischung einzelner Tabellen aus
beiden Quellen. `CDC_CONFIG_FILE` selbst trägt den Dateipfad; leer/unbenannt
bedeutet kein Dateizugriff, der bestehende Env-only-Pfad bleibt unverändert
Default.

### SPEC-017 — NATS-Wecksignal (Subjekt- und Nachrichtenform)

Technische Ausgestaltung von [`LH-FA-SST-007`](lastenheft.md): Core NATS
(kein JetStream) als reines, verlustbehaftetes Wecksignal — keine eigene
Zustellgarantie, keine eigene Nachvollziehbarkeit. Nachvollziehbarkeit bleibt
ausschließlich beim bestehenden Lesezugriffsweg
([`LH-FA-REA-001`](lastenheft.md) ff., `cdc.changes`).

| Merkmal | Festlegung |
|---|---|
| Subjekt-Schema | `cdc.changes.<source_id>.<schema>.<table>` — ein Subjekt je Tabelle einer Quelle, `<source_id>` identisch zur konfigurierten `CDC_SOURCE_ID`, `<schema>`/`<table>` die Klartext-Bezeichner der betroffenen Tabelle. Ein Consumer mit Tabelleninteresse abonniert das vollständige vierstufige Subjekt; ein Consumer, der alle Tabellen einer Quelle verfolgt, abonniert `cdc.changes.<source_id>.>` (NATS-Wildcard); ein Consumer, der mehrere Quellen verfolgt, abonniert `cdc.changes.>`. Ein Notify je Transaktion und distinkter berührter Tabelle (dedupliziert) |
| Nachrichteninhalt | leerer Payload (Trigger ohne Daten) — kein Change-Inhalt, keine Positionsangabe. Jede Nachricht bedeutet ausschließlich „lies erneut über den bestehenden Zugriffsweg"; Fehlen oder Verdopplung einer Nachricht trägt keine eigene Bedeutung |
| Zustellgarantie | keine (Core NATS, Fire-and-Forget); ein nicht verbundener oder gerade getrennter Consumer verpasst das Signal ersatzlos — zulässig nach `LH-FA-SST-007` Boundary/Negative |
| Reconnect-Verhalten | die Client-Bibliothek (`github.com/nats-io/nats.go`) trägt automatisches Reconnect mit eingebautem Backoff auf Verbindungsebene; auf Nachrichtenebene gibt es keinen gesonderten Replay — der Consumer holt entfallene Änderungen ausschließlich über den bestehenden Lesezugriffsweg nach |
| Fehlerklasse bei Notify-Fehlschlag | `transient` (`SPEC-008`) — der Fehler wird an der Aufrufstelle (`CaptureService`) abgefangen und propagiert **nicht** in den Rückgabewert des Capture-Aufrufs; er darf die bereits erfolgte Persistierung oder das bereits erfolgte Source-ACK (`LH-QA-REL-001.a`) nicht beeinflussen |
| Aktivierung | optional über `CDC_NATS_URL`; ungesetzt bedeutet deaktiviertes Feature, keine NATS-Verbindung, unverändertes Bestandsverhalten |

### SPEC-018 — HTTP-API: Endpunkte und Token-Header-Form

Technische Ausgestaltung von [`LH-FA-SST-006`](lastenheft.md): Endpunkt,
Methode, JSON-Schema und Token-Header-Form für alle neun Port-gedeckten
Fähigkeiten des HTTP-Adapters. Changes-Lesen
([`LH-FA-REA-*`](lastenheft.md)) und Diagnose/Health bleiben außerhalb.

**Authn-Header** (für alle Endpunkte gleich): `Authorization: Bearer
<token>` — fehlend oder einer nicht konfigurierten Klasse entsprechend →
`401`; ein bekanntes `reader`-Token gegen einen administrativen Endpunkt →
`403`; ein bekanntes `admin`-Token erreicht jeden Endpunkt (`admin` deckt
implizit die lesende Klasse ab). Token-Umgebungsvariablen:
`CDC_API_TOKEN_READER` (lesende Rechtsklasse), `CDC_API_TOKEN_ADMIN`
(schreibend/administrativ) — orthogonal zum DB-Rollenmodell der
Verdrahtung: die API-Token-Prüfung entscheidet an der HTTP-Schicht, welcher
Use Case erreichbar ist; welche DSN der Adapter darunter benutzt, bleibt
die bei der Verdrahtung fixierte. Aktivierung optional über
`CDC_HTTP_ADDR`; ungesetzt bedeutet deaktiviertes Feature, kein
HTTP-Server, unverändertes Bestandsverhalten.

**Fehler-Antwortform** (für alle Endpunkte gleich):
`{"error": "<Klartext>"}` bei `400`/`401`/`403`/`404`/`500` — ungültiger
JSON-Body oder eine verletzte Domänen-Invariante (`400`),
fehlender/unbekannter Bearer-Token (`401`), bekanntes Token mit
unzureichender Rechtsklasse (`403`), eine an der Quelle physisch fehlende
Tabelle (`404`, nur `EnableTable`/`DisableTable`/`GetStatus`), unerwarteter
interner Fehler (`500`).

| Fähigkeit | Endpunkt | Rechtsklasse | Request | Response |
|---|---|---|---|---|
| `RegisterConsumer` ([`LH-FA-CON-001`](lastenheft.md)) | `POST /consumers` | `admin` | `{"consumer_id": "<string>", "name": "<string>"}` — beide Pflicht | `201`: `{"consumer_id": "<string>", "name": "<string>", "already_registered": <bool>}` — `already_registered` trägt die Idempotenz fort, kein gesonderter Statuscode |
| `AcknowledgeConsumer` ([`LH-FA-CON-004`](lastenheft.md)) | `POST /consumers/acknowledge` | `admin` | `{"consumer_id": "<string>", "source_id": "<string>", "offset": <uint64>}` — alle Pflicht | `200`: `{"consumer_id": "<string>", "source_id": "<string>", "offset": <uint64>}` — die fortgeführte Position; eine frühere Position oder eine Position einer anderen Quelle → `400` |
| `GetConsumerPosition` ([`LH-FA-CON-003`](lastenheft.md), [`LH-FA-CON-005`](lastenheft.md)) | `GET /consumers/position?consumer_id=<string>` | `reader` oder `admin` | Query-Parameter `consumer_id` Pflicht | `200`: `{"consumer_id": "<string>", "source_id": "<string>", "offset": <uint64>, "acknowledged": <bool>}` — `acknowledged=false` liest die definierte Anfangsposition ohne Bestätigung |
| `RemoveConsumer` ([`LH-FA-CON-006`](lastenheft.md)) | `POST /consumers/remove` | `admin` | `{"consumer_id": "<string>"}` — Pflicht | `200`: `{"consumer_id": "<string>", "removed": <bool>}` — ein nie registrierter Consumer meldet `removed=false`, kein `404` (Idempotenz statt Fehler gegen eine unbekannte Ressource, fachlich gleichwertig zu CLI/SQL) |
| `EnableTable` ([`LH-FA-CFG-001`](lastenheft.md)) | `POST /tables/enable` | `admin` | `{"source": "<string>", "schema": "<string>", "table": "<string>", "table_id": "<string>", "schema_version_id": "<string>", "version": <int64>, "publication": "<string>"}` — alle Pflicht, `version` ≥ 1 | `201`: `{"table_id": "<string>", "source": "<string>", "schema": "<string>", "table": "<string>", "already_enabled": <bool>}`; physisch fehlende Tabelle an der Quelle → `404` |
| `DisableTable` ([`LH-FA-CFG-002`](lastenheft.md)) | `POST /tables/disable` | `admin` | `{"source": "<string>", "schema": "<string>", "table": "<string>", "publication": "<string>"}` — alle Pflicht | `200`: `{"removed": <bool>, "retained": <bool>}`; physisch fehlende Tabelle an der Quelle → `404` |
| `GetStatus` ([`LH-FA-CFG-003`](lastenheft.md)) | `GET /tables/status?source=<string>&schema=<string>&table=<string>&publication=<string>` | `reader` oder `admin` | alle vier Query-Parameter Pflicht | `200`: `{"enabled": <bool>, "retained": <bool>}` — beide `false` liest eine nie aktivierte Tabelle; physisch fehlende Tabelle an der Quelle → `404` |
| `ListTables` ([`LH-FA-CFG-004`](lastenheft.md)) | `GET /tables?source=<string>&publication=<string>` | `reader` oder `admin` | beide Query-Parameter Pflicht | `200`: `{"tables": [{"table_id": "<string>", "source": "<string>", "schema": "<string>", "table": "<string>"}, …], "retained": [...]}` — ohne Aktivierung beide Listen leer |
| `RunRetention` ([`LH-FA-RET-002`](lastenheft.md)…[`004`](lastenheft.md)) | `POST /retention/run` | `admin` | `{"source": "<string>", "min_age_nanos": <int64>}` — `source` Pflicht, `min_age_nanos` ≥ 0 | `200`: `{"deleted": <int>}` — Anzahl real gelöschter Changes |

### SPEC-019 — `cdc.administration_request` (Antrags-Datensatz)

Feldform des Antrags-Datensatzes der schreibenden SQL-Administration
([`LH-FA-ADM-001`](lastenheft.md), [`LH-FA-CFG-005`](lastenheft.md)):
`cdc.enable_table`/`cdc.disable_table`/`cdc.exclude_column`/
`cdc.include_column` schreiben ausschließlich eine Zeile hierher und senden
`pg_notify` auf dem Kanal `cdc_administration`; der laufende Capture-Prozess
liest die offenen Anträge und vermerkt das Ergebnis in derselben Zeile.

| Spalte | Typ | Pflicht | Bedeutung |
|---|---|---|---|
| `administration_request_id` | text (PK) | ja | von der SQL-Funktion vergeben; zugleich der `pg_notify`-Payload |
| `source_id` | text (FK `cdc.source`) | ja | Quelle des Antrags |
| `schema_name` / `table_name` | text | ja | adressierte Tabelle |
| `column_name` | text | nein | Ziel-Spalte der beiden Spalten-Antragsarten; die beiden Tabellen-Antragsarten tragen hier NULL |
| `request_kind` | text | ja | geschlossene Menge `enable` \| `disable` \| `exclude_column` \| `include_column` |
| `requested_at` | timestamptz | ja, Default `current_timestamp` | Anlage-Zeitpunkt; die Verarbeitungs-Ordnung |
| `status` | text | ja, Default `pending` | geschlossene Menge `pending` \| `applied` \| `failed` |
| `error_message` | text | nein | Fehlertext eines `failed`-Antrags; `applied` trägt NULL |

`column_name` ist für die beiden Spalten-Antragsarten Pflicht (ein Antrag
ohne Spalte adressiert kein Ziel, Domänen-Invariante des
Antrags-Konstruktors); ein Spaltenausschluss gegen eine an der Quelle nicht
existierende Spalte endet als `failed` mit dem Fehlertext der
`ErrSourceColumnMissing`-Ausprägung (Klartext „Spalte existiert nicht an der
Quelle", gefolgt von der Adresse `schema.table.column`, getrennt durch
Punkte). Die beiden Tabellen-Antragsarten tragen unverändert Bindungs-Zeilen
und Publication nach.

---

## 3. Defaults und Konstanten

Regeln dieser Sektion: Die ADR, die einen Wert festlegt, deklariert das
aufwärts in ihrem `Schärft:`-Feld — kein ADR-Rückzeiger hier. Die
`SPEC-<NNN>` ist das, was ihr `Schärft:`-Feld **benennt**; ohne sie kann
eine ADR nur den ganzen Abschnitt nennen.

| ID | Name | Wert | Begründung |
|---|---|---|---|
| `SPEC-006` | `CDC_SCHEMA` | `cdc` | Standard-Schema-Name für alle CDC-Objekte (§2, SPEC-001) |
| `SPEC-007` | `HEALTH_STATES` | `healthy`, `degraded`, `unhealthy` | Health-Zustände; Readiness zeigt, ob die Instanz ihre Betriebsrolle erfüllen kann |
| `SPEC-012` | `PG_MAJOR_VERSIONS` | 17, 18 | Vorschlagsregel: die zwei neuesten aktiven Major-Versionen (Stand 2026-09-09; PostgreSQL 19 unmittelbar vor Release — Aufnahme als spätere Ausweitung) |
| `SPEC-013` | `CDC_THRESHOLDS` | Capture-Lag p95 ≤ 1 s · Warn > 5 s · Fehler > 60 s; WAL-Rückstand Warn > 100 MiB · Fehler > 1 GiB | Initialwerte Commit→CDC-Verfügbarkeit ([`LH-QA-PER-004`](lastenheft.md), `cdc_capture_lag`) und WAL-Wachstum inaktiver Slots ([`LH-QA-REL-003`](lastenheft.md), `cdc_wal_retention_bytes`); über ADR schärfbar |
| `SPEC-014` | `LOAD_TIERS` | klein: ≤ 10 Changes/s · mittel: 100 Changes/s über 30 min · groß: 1.000 Changes/s über 60 min | Benchmark-Stufen für die Skalierbarkeits-Prüfung ([`LH-QA-PER-002`](lastenheft.md)): von kleinen Datenbanken bis zu kontinuierlichen Änderungsvolumina; über ADR schärfbar |

---

## 4. Fehler-Codes und Logging-Felder

Regeln dieser Sektion: Der Fehler-Code ist ein Laufzeit-Symbol, die
`SPEC-<NNN>` benennt die *Festlegung* darüber — beide stehen nebeneinander.

Fehler werden mindestens in die folgenden Klassen klassifiziert
(`SPEC-008`):

| ID | Klasse | Bedingung | Aktion |
|---|---|---|---|
| `SPEC-008` | `transient` | vorübergehend nicht verfügbare Quelle/Speicher | Erneut versuchen mit begrenztem Backoff |
| `SPEC-008` | `configuration` | ungültige/falsch gesetzte Konfiguration | Sichtbarer Fehler; kein Start und keine Fortsetzung im falschen Stand |
| `SPEC-008` | `permission` | fehlende Berechtigung | Sichtbarer Fehler; kein stiller Retry |
| `SPEC-008` | `schema` | nicht sicher interpretierbare Schemaänderung/Dekodierfehler | Sichtbarer Fehler (LH-FA-SCH-004.a); kein stilles Überspringen |
| `SPEC-008` | `storage` | Persistenzfehler im ChangeStore | **Kein Source-ACK** (LH-QA-REL-001.a) |
| `SPEC-008` | `replication` | Replication-Stream/Slot-Störung: **Stream-Ordnungsverletzung** (BEGIN/COMMIT/Change außerhalb der erwarteten Reihenfolge) oder **Transport-/Verbindungsstörung** (Verbindungsaufbau, Start, Keepalive, Quell-Bestätigung) | Stream-Ordnungsverletzung: sichtbarer Fehler, harter Abbruch, keine Fortsetzung im widersprüchlichen Stand; Transport-/Verbindungsstörung: Überwachung über Schwellen (§5, WAL-Rückstand); kontrollierte Fortsetzung |
| `SPEC-008` | `internal` | unerwarteter interner Fehler | Sichtbarer Fehler; Restart-Strategie nach [`LH-QA-REL-002`](lastenheft.md) |

Keine Credentials in Logs.

---

## 5. Metriken und Tracing-Felder

Regeln dieser Sektion: verbindliche Felder der Betriebsschnittstelle
(vorgesehen: Prometheus/OpenTelemetry-Adapter). Inhalt deckt
[`LH-QA-OPS-003`](lastenheft.md) ab.

| ID | Metrik | Inhalt | Quelle |
|---|---|---|---|
| `SPEC-009` | `cdc_changes_processed` | verarbeitete Changes | ChangeStore |
| `SPEC-009` | `cdc_changes_pending` | ausstehende/unbestätigte Changes | ChangeStore |
| `SPEC-009` | `cdc_transactions_total` | persistierte Transaktionen | ChangeStore |
| `SPEC-009` | `cdc_errors_total{class}` | Fehler je Klasse (§4) | Application |
| `SPEC-009` | `cdc_capture_lag` | Abstand Quelländerung → CDC-Verfügbarkeit ([`LH-FA-ADM-004`](lastenheft.md)) | Capture |
| `SPEC-009` | `cdc_consumer_lag` | Rückstand je Consumer ([`LH-FA-ADM-005`](lastenheft.md)) | ConsumerState |
| `SPEC-009` | `cdc_consumer_position` | bestätigte Position je Consumer, roh ([`LH-QA-OPS-003`](lastenheft.md)) | ConsumerState |
| `SPEC-009` | `cdc_wal_retention_bytes` | WAL-Rückstand/Replication-Slot-Zustand | Replication Stream |
| `SPEC-009` | `cdc_storage_bytes` | Speicherverbrauch der CDC-Daten | ChangeStore |
| `SPEC-009` | `cdc_oldest_change_age` | Alter des ältesten aufbewahrten Changes | ChangeStore |

Warn- und Fehlerschwellen für WAL-Rückstand und Capture-Lag sind
konfigurierbar; Initialwerte: SPEC-013. Replication-Slot-Zustand,
WAL-Rückstand und Capture-Lag werden überwacht.

---

## 6. Externe Verträge

| ID | System | Version | Vertrag-Datei |
|---|---|---|---|
| `SPEC-010` | PostgreSQL Logical Replication (`pgoutput`) | PostgreSQL 17 und 18 (SPEC-012, [`LH-QA-POR-001`](lastenheft.md)) | — (Vertrag steht in diesem Dokument, §1 LH-FA-CFG-001.a; Zeiger auf `spec/architecture.md` entfällt, bis diese gefüllt ist) |
| `SPEC-011` | OCI-Container-Runtime | OCI-Image-Spec | — (Deployment; keine privilegierten Rechte nötig) |
| `SPEC-015` | Eigenständiges Executable (Deployment-Form neben SPEC-011) | Cross-Compile: Linux amd64/arm64 (primär, [`LH-QA-POR-002`](lastenheft.md)); darwin/amd64, darwin/arm64, windows/amd64 perspektivisch; `CGO_ENABLED=0` | — (Deployment-Artefakt; Cross-Compile in CI/CD) |
| `SPEC-017` | NATS Core (Wecksignal, kein JetStream) | NATS-Server 2.x, Go-Client `github.com/nats-io/nats.go` | — (Vertrag steht in diesem Dokument, §2 SPEC-017) |

---

## 7. Historie

Regeln dieser Sektion: **kein ADR- und kein Slice-Verweis.** Die Decken-Regel
gilt für alle drei Spec-Straten, auch hier — welche ADR eine Festlegung
schärft, deklariert die ADR aufwärts in ihrem `Schärft:`-Feld.

| Datum | Änderung |
|---|---|
| 2026-09-09 | Initial — Technik-Inhalt überführt aus dem zurückgezogenen Pflichtenheft-Entwurf (Git: angelegt in c70ee1c, zurückgezogen in 5a6f8ea); Algorithmen, CDC-Schema, Fehlerklassen, Metriken, externe Verträge |
| 2026-09-09 | SPEC-015 ergänzt: eigenständiges Executable — Cross-Compile Linux amd64/arm64 primär, darwin/windows perspektivisch, `CGO_ENABLED=0`; Deployment-Form neben SPEC-011 (OCI) |
| 2026-09-12 | LH-FA-CON-001.a und LH-FA-CON-004.a ergänzt: Registrierungs- und Bestätigungslogik sind eigenständig getestet, aber ohne von außen erreichbaren Zugriffsweg — die Wahl des Zugriffswegs bleibt eine offene technische Frage |
| 2026-09-13 | SPEC-016 ergänzt: Feldform der optionalen YAML-Konfigurationsdatei (`CDC_CONFIG_FILE`) — Schlüsselnamen, Precedence-Verweis, DSN-Ausschluss |
| 2026-09-13 | SPEC-017 ergänzt: NATS-Wecksignal — Subjekt-Schema (`cdc.changes.<source_id>`), leerer Payload, Zustellgarantie, Reconnect-Verhalten, Fehlerklasse `transient`, Aktivierung über `CDC_NATS_URL`; externe-Verträge-Zeile in §6 |
| 2026-09-13 | SPEC-017 Subjekt-Schema korrigiert: tabellen-granulares Subjekt `cdc.changes.<source_id>.<schema>.<table>` statt quellen-weit; übrige SPEC-017-Festlegungen unverändert |
| 2026-09-14 | SPEC-018 ergänzt: HTTP-API `RegisterConsumer` — Endpunkt/Methode, JSON-Request-/Response-Schema, Fehler-Antwortform `400`/`401`/`403`/`500`, Token-Header-Form, Aktivierung über `CDC_HTTP_ADDR` |
| 2026-09-14 | SPEC-018 erweitert: acht weitere Endpunkte (Acknowledge-/Position-/Remove-Consumer, Enable-/Disable-/Status-/List-Table, Retention-Lauf) — Fehler-Antwortform um `404` (physisch fehlende Tabelle an der Quelle) ergänzt |
| 2026-09-14 | SPEC-019 ergänzt: Feldform von `cdc.administration_request` — Spalte `column_name`, erweiterte `request_kind`-Menge `enable`/`disable`/`exclude_column`/`include_column`, `failed`-Fehlertext der fehlenden Spalte |
