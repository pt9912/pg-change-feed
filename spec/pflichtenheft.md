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

- Offene Transaktionen sind nicht konsumierbar (LH-FA-CAP-006).
- Rollbacks erzeugen keine committed Changes (LH-FA-CAP-007).
- Jede persistierte Quelltransaktion besitzt eine interne ID und
  Commit-Position; Changes besitzen eine eindeutige Sequenz innerhalb der
  Transaktion (LH-FA-CAP-004, LH-FA-DAT-004).
- Große Transaktionen dürfen nicht unbegrenzt im RAM gehalten werden;
  der TransactionBufferPort (SPEC-005) ermöglicht Spooling.

### LH-FA-CFG-001.a — Aktivierung einer Tabelle

**Eingabe:** Tabellenname, CDC-Konfiguration. **Ausgabe:** aktivierte
Tabelle; aktualisierte Publication/Slot-Verwaltung.

1. Replica Identity prüfen. `REPLICA IDENTITY FULL` kann für vollständige
   alte Werte erforderlich sein (LH-FA-CAP-008), wird aber **nicht
   ungefragt** aktiviert.
2. Publication und Logical Replication Slot verwalten; als Output Plugin
   wird `pgoutput` verwendet.
3. Operationen erfassen: INSERT, UPDATE, DELETE. TRUNCATE wird erkannt und
   im MVP explizit als nicht unterstützt behandelt (Lastenheft fordert
   TRUNCATE nicht, §5 Out-of-Scope).

### LH-FA-REA-004.a — Deterministische Sortierung

**Eingabe:** Positionsbereich bzw. Startposition, Limit, optionaler
Tabellenfilter. **Ausgabe:** begrenzte, deterministisch sortierte
Ergebnismenge.

Sortierung nach (Commit-Position der Quelltransaktion, Transaktions-ID,
Sequenz innerhalb der Transaktion) — bei identischer Eingabe ist die
Reihenfolge damit stabil (LH-FA-REA-004, LH-FA-REA-003). Changes können
ab einer SourcePosition bzw. innerhalb eines Positionsbereichs gelesen
werden; Lesen verändert gespeicherte Positionen nicht.

### LH-FA-RET-004.a — Safe Watermark der Retention

**Eingabe:** Consumer-Positionen, Retention-Konfiguration. **Ausgabe:**
Bereinigungsmenge.

1. Die Safe Watermark basiert auf den relevanten Consumer-Positionen;
   Changes, die aktive Consumer benötigen, dürfen nicht automatisch
   gelöscht werden.
2. Zeitbasierte Mindestaufbewahrung ist zusätzlich konfigurierbar
   (LH-FA-RET-003).
3. Langsame Consumer, die die Safe Watermark blockieren, sind sichtbar
   (LH-FA-RET-005).

Retention ist eine Domain Policy; sie liegt nicht in einem Adapter.

### LH-FA-SCH-004.a — Schema-Änderungs-Klassifikation

**Eingabe:** Relation Metadata aus dem Replication Stream. **Ausgabe:**
TableSchema-/SchemaVersion-Modell (SPEC-004) oder sichtbarer Fehler.

Relation Metadata wird in technologieunabhängige
TableSchema-/SchemaVersion-Modelle übersetzt; jeder Change referenziert
eine Schema-Version (LH-FA-SCH-005). Nicht sicher interpretierbare
Schemaänderungen führen zu einem sichtbaren Fehler (Fehlerklasse `schema`,
§4) statt stiller Fehlinterpretation (LH-FA-SCH-004).

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
| `cdc.capture_state` | Betriebs-/Capture-Zustand |

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
Images; Verfügbarkeit je Operationstyp siehe LH-FA-CAP-008).

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
referenziert eine Schema-Version (LH-FA-SCH-005). Die Übersetzung
erfolgt in LH-FA-SCH-004.a.

### SPEC-005 — TransactionBuffer

Vertrag für das Spooling großer Quelltransaktionen: Zwischenstände einer
noch nicht abgeschlossenen Transaktion werden außerhalb des RAM
gehalten, bis die Transaktion committed und dauerhaft persistiert ist
(vorgesehene Implementierung: FileTransactionBufferAdapter).

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

Offene Festlegungen — Delegationen aus dem Lastenheft, bewusst noch nicht
gesetzt, jede mit ihrem Schlusspunkt:

| Gegenstand | Referenz | Schlusspunkt |
|---|---|---|
| PostgreSQL-Major-Versionen | LH-QA-POR-001 (SPEC-010) | Festlegung bis zur MVP-Abnahme; Vorschlagsregel: die zwei neuesten aktiven Major-Versionen |
| Latenz-Schwellen-Initialwerte | LH-QA-PER-004 | Festlegung mit dem Benchmark-Design vor der MVP-Abnahme |
| Skalierungs-Volumina | LH-QA-PER-002 | Festlegung mit dem Benchmark-Design vor der MVP-Abnahme |

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
| `SPEC-008` | `replication` | Replication-Stream/Slot-Störung | Überwachung über Schwellen (§5, WAL-Rückstand); kontrollierte Fortsetzung |
| `SPEC-008` | `internal` | unerwarteter interner Fehler | Sichtbarer Fehler; Restart-Strategie nach LH-QA-REL-002 |

Keine Credentials in Logs.

---

## 5. Metriken und Tracing-Felder

Regeln dieser Sektion: verbindliche Felder der Betriebsschnittstelle
(vorgesehen: Prometheus/OpenTelemetry-Adapter). Inhalt deckt
LH-QA-OPS-003 ab.

| ID | Metrik | Inhalt | Quelle |
|---|---|---|---|
| `SPEC-009` | `cdc_changes_processed` | verarbeitete Changes | ChangeStore |
| `SPEC-009` | `cdc_changes_pending` | ausstehende/unbestätigte Changes | ChangeStore |
| `SPEC-009` | `cdc_transactions_total` | persistierte Transaktionen | ChangeStore |
| `SPEC-009` | `cdc_errors_total{class}` | Fehler je Klasse (§4) | Application |
| `SPEC-009` | `cdc_capture_lag` | Abstand Quelländerung → CDC-Verfügbarkeit (LH-FA-ADM-004) | Capture |
| `SPEC-009` | `cdc_consumer_lag` | Rückstand je Consumer (LH-FA-ADM-005) | ConsumerState |
| `SPEC-009` | `cdc_consumer_position` | bestätigte Position je Consumer, roh (LH-QA-OPS-003) | ConsumerState |
| `SPEC-009` | `cdc_wal_retention_bytes` | WAL-Rückstand/Replication-Slot-Zustand | Replication Stream |
| `SPEC-009` | `cdc_storage_bytes` | Speicherverbrauch der CDC-Daten | ChangeStore |
| `SPEC-009` | `cdc_oldest_change_age` | Alter des ältesten aufbewahrten Changes | ChangeStore |

Warn- und Fehlerschwellen für WAL-Rückstand und Capture-Lag sind
konfigurierbar; Replication-Slot-Zustand, WAL-Rückstand und Capture-Lag
werden überwacht.

---

## 6. Externe Verträge

| ID | System | Version | Vertrag-Datei |
|---|---|---|---|
| `SPEC-010` | PostgreSQL Logical Replication (`pgoutput`) | mehrere aktive Major-Versionen (LH-QA-POR-001); konkrete Versionen: offene Festlegung, §3 | — (Vertrag steht in diesem Dokument, §1 LH-FA-CFG-001.a; Zeiger auf `spec/architecture.md` entfällt, bis diese gefüllt ist) |
| `SPEC-011` | OCI-Container-Runtime | OCI-Image-Spec | — (Deployment; keine privilegierten Rechte nötig) |

---

## 7. Historie

Regeln dieser Sektion: **kein ADR- und kein Slice-Verweis.** Die Decken-Regel
gilt für alle drei Spec-Straten, auch hier — welche ADR eine Festlegung
schärft, deklariert die ADR aufwärts in ihrem `Schärft:`-Feld.

| Datum | Änderung |
|---|---|
| 2026-09-09 | Initial — Technik-Inhalt überführt aus dem zurückgezogenen Pflichtenheft-Entwurf (Git: angelegt in c70ee1c, zurückgezogen in 5a6f8ea); Algorithmen, CDC-Schema, Fehlerklassen, Metriken, externe Verträge |