# Pflichtenheft – PG Change Feed

**Projekt:** PG Change Feed
**Version:** 0.2
**Architektur:** Hexagonale Architektur / Ports & Adapters
**Status:** Entwurf

## PH-ZWE-001 – Zweck
PG Change Feed stellt PostgreSQL eine persistente Change-Data-Capture-Abstraktion bereit. Änderungen werden transaktionsgetreu erfasst, dauerhaft gespeichert und unabhängigen Consumern bereitgestellt.

## PH-ZIE-001 – Zielarchitektur
PostgreSQL Logical Replication mit `pgoutput` bildet die Änderungsquelle. Der Replication Stream treibt den Capture-Use-Case und ist ein **Driving Adapter**. Die Bestätigung einer WAL-Position wird vom Application Core über einen **Outbound Port** und einen **Driven Adapter** ausgelöst.

    PostgreSQL WAL
          |
          v
    Replication Stream
    Driving Adapter
          |
          v
    CaptureInboundPort
          |
          v
    Application / Domain
          |
          +------------------+
          |                  |
          v                  v
    ChangeStorePort   ReplicationAckPort
          |                  |
          v                  v
    PostgreSQL Store  PostgreSQL ACK
      Driven Adapter   Driven Adapter

## PH-ARC-001 – Terminologie
Verbindlich:
- Ports: Inbound / Outbound
- Adapters: Driving / Driven

Driving Adapters rufen Inbound Ports auf. Application Services verwenden Outbound Ports. Driven Adapters implementieren Outbound Ports.

## PH-ARC-002 – Abhängigkeiten
Abhängigkeiten zeigen nach innen:

    Adapters -> Application -> Domain

Der Domain Core kennt keine PostgreSQL-, HTTP-, CLI-, Metrics- oder Dateisystembibliotheken.

## PH-DRV-001 – Replication Stream
Der PostgreSQL Logical Replication Stream ist Driving Adapter. Er baut die Replication-Verbindung auf, empfängt und dekodiert `pgoutput` und übersetzt technische Nachrichten in Aufrufe des CaptureInboundPort.

## PH-DRV-002 – Weitere Driving Adapters
CLI, SQL API sowie spätere HTTP-/gRPC-Schnittstellen sind Driving Adapters.

## PH-IN-001 – CaptureInboundPort
Der CaptureInboundPort nimmt dekodierte Replication-Ereignisse bzw. vollständige Quelltransaktionen entgegen.

## PH-IN-002 – Administrative Inbound Ports
Vorgesehen sind Use Cases für Tabellenaktivierung/-deaktivierung, Consumer-Registrierung, Lesen, Consumer-ACK, Reset, Retention und Status.

## PH-OUT-001 – ChangeStorePort
Persistiert committed CDC-Transaktionen und stellt persistierte Changes bereit.

## PH-OUT-002 – ReplicationAckPort
Bestätigt gegenüber PostgreSQL ausschließlich Positionen, deren abhängige Changes bereits dauerhaft gespeichert wurden.

## PH-OUT-003 – Weitere Outbound Ports
Vorgesehen sind ConsumerStatePort, SchemaStorePort, TransactionBufferPort und MetricsPort.

## PH-CAP-001 – Logical Replication
Als Output Plugin wird `pgoutput` verwendet. PG Change Feed verwaltet eine Publication und einen konfigurierbaren Logical Replication Slot.

## PH-CAP-002 – Operationen
INSERT, UPDATE und DELETE müssen erfasst werden. TRUNCATE muss erkannt werden und darf im MVP explizit als nicht unterstützt behandelt werden.

## PH-TRX-001 – Transaktionsgrenzen
BEGIN, COMMIT und relevante Streaming-Grenzen werden berücksichtigt. Offene Transaktionen sind nicht konsumierbar. Rollbacks erzeugen keine committed Changes.

## PH-TRX-002 – Reihenfolge
Jede persistierte Quelltransaktion besitzt eine interne ID und Commit-Position. Changes besitzen eine eindeutige Sequenz innerhalb der Transaktion.

## PH-TRX-003 – Große Transaktionen
Große Transaktionen dürfen nicht unbegrenzt im RAM gehalten werden. Ein TransactionBufferPort ermöglicht Spooling.

## PH-LSN-001 – SourcePosition
Der Core verwendet `SourcePosition`. Der PostgreSQL-Adapter mappt PostgreSQL LSN darauf.

## PH-LSN-002 – Persist-before-ACK
Zentrale Invariante:

    ACK(position) => durable(all changes <= position)

Ablauf:

    Receive -> Decode -> Persist -> COMMIT Store -> ACK Source

Ein Crash zwischen Persistenz und ACK darf höchstens zu erneuter Verarbeitung führen.

## PH-LSN-003 – Idempotenz
Wiederholte WAL-Daten müssen deduplizierbar bzw. idempotent persistierbar sein.

## PH-DAT-001 – CDC-Schema
Standardmäßig wird das Schema `cdc` verwendet.

Vorgesehene Tabellen:
- cdc.source
- cdc.source_table
- cdc.transaction
- cdc.change
- cdc.consumer
- cdc.consumer_position
- cdc.schema_version
- cdc.capture_state

## PH-DAT-002 – Change
`cdc.change` enthält mindestens Change-ID, Transaction-ID, Source-Table-ID, Sequenz, Operation, old_data, new_data und Schema-Version.

## PH-DAT-003 – Row Images
Im MVP werden `old_data` und `new_data` als `jsonb` gespeichert.

## PH-RID-001 – Replica Identity
Beim Aktivieren einer Tabelle wird Replica Identity geprüft. `REPLICA IDENTITY FULL` kann für vollständige alte Werte erforderlich sein, wird aber nicht ungefragt aktiviert.

## PH-CON-001 – Consumer
Mehrere benannte Consumer besitzen unabhängige persistierte Positionen.

## PH-CON-002 – Consumer ACK
Lesen verändert die Position nicht. ACK erfolgt explizit und bewegt die Position regulär nur vorwärts. Reset ist eine administrative Sonderoperation.

## PH-CON-003 – Delivery
Die Delivery-Semantik ist At-Least-Once. Consumer sollen idempotent arbeiten.

## PH-REA-001 – Lesen
Changes können ab einer SourcePosition bzw. innerhalb eines Positionsbereichs gelesen werden. Ergebnismengen sind begrenzbar und deterministisch sortiert.

## PH-RET-001 – Retention
Retention ist eine Domain Policy. Die Safe Watermark basiert auf den relevanten Consumer-Positionen. Von aktiven Consumern benötigte Changes dürfen nicht automatisch gelöscht werden.

## PH-RET-002 – Mindestaufbewahrung
Eine zeitbasierte Mindestaufbewahrung ist konfigurierbar. Langsame Consumer müssen sichtbar sein.

## PH-WAL-001 – WAL-Rückstau
Replication-Slot-Zustand, WAL-Rückstand und Capture-Lag werden überwacht. Warn- und Fehlerschwellen sind konfigurierbar.

## PH-SCH-001 – Schema Evolution
Relation Metadata wird in technologieunabhängige TableSchema-/SchemaVersion-Modelle übersetzt. Jeder Change referenziert eine Schema-Version.

## PH-SCH-002 – Inkompatible Änderungen
Nicht sicher interpretierbare Schemaänderungen führen zu einem sichtbaren Fehler statt stiller Fehlinterpretation.

## PH-ERR-001 – Fehlerklassen
Fehler werden mindestens klassifiziert als transient, configuration, permission, schema, storage, replication und internal.

## PH-ERR-002 – Fehlerstrategie
Transiente Fehler werden mit begrenztem Backoff erneut versucht. Bei Persistenzfehlern erfolgt kein Source-ACK. Dekodierungsfehler werden nicht still übersprungen.

## PH-REC-001 – Recovery
Nach Crash wird anhand des Replication Slots und persistierter Zustände fortgesetzt. Wiederholung wird gegenüber möglichem Datenverlust bevorzugt.

## PH-DRN-001 – Driven Adapters
Vorgesehen sind:
- PostgresChangeStoreAdapter
- PostgresReplicationAckAdapter
- PostgresMetadataAdapter
- FileTransactionBufferAdapter
- Prometheus/OpenTelemetry Adapter

## PH-SEC-001 – Sicherheit
Dedizierter PostgreSQL-Benutzer, Least Privilege, getrennte Capture-/Reader-/Admin-Rechte, keine Credentials in Logs und TLS-Unterstützung.

## PH-OBS-001 – Observability
Strukturierte Logs sowie Metriken für Changes, Transaktionen, Fehler, Capture-Lag, Consumer-Lag, WAL-Rückstand und Storage-Größe.

## PH-OBS-002 – Health
Health-Zustände: healthy, degraded, unhealthy. Readiness zeigt, ob die Instanz ihre Betriebsrolle erfüllen kann.

## PH-DEP-001 – Deployment
PG Change Feed wird als eigenständiges Linux-Executable und OCI-Container ausgeliefert. Der Container benötigt keine privilegierten Rechte.

## PH-DEP-002 – Entwicklung
Eine Docker-Compose-Umgebung stellt PostgreSQL, PG Change Feed und Test-Consumer reproduzierbar bereit.

## PH-STR-001 – Paketstruktur

    src/
    ├── domain/
    ├── application/
    │   ├── ports/
    │   │   ├── inbound/
    │   │   └── outbound/
    │   └── services/
    ├── adapters/
    │   ├── driving/
    │   │   ├── postgres_replication_stream/
    │   │   ├── cli/
    │   │   ├── sql/
    │   │   └── http/
    │   └── driven/
    │       ├── postgres_storage/
    │       ├── postgres_replication_ack/
    │       ├── postgres_metadata/
    │       ├── filesystem_spool/
    │       └── telemetry/
    └── bootstrap/

## PH-TST-001 – Teststrategie
Viele Domain-Tests, Application-Tests mit Fake Ports, gezielte Adaptertests, reale PostgreSQL-Integrationstests und wenige kritische E2E-Tests.

## PH-TST-002 – Kritische Tests
Mindestens:
- INSERT/UPDATE/DELETE
- offene Transaktion nicht sichtbar
- Rollback ohne Changes
- Persistenzfehler verhindert ACK
- Crash nach Persistenz vor ACK erzeugt keine Lücke
- unabhängige Consumer-Positionen
- sichere Retention
- große Transaktion ohne unbegrenzten RAM-Verbrauch

## PH-MVP-001 – MVP
Enthalten: eine Source, pgoutput, Publication, Logical Slot, INSERT/UPDATE/DELETE, Transaktionen, persistente Changes, SourcePosition, SQL-Lesen, Restart, Docker und Integrationstests.

Nicht enthalten: Kafka, NATS, RabbitMQ, Kubernetes Operator, automatische HA, vollständige DDL-Replikation und revisionssicheres Auditing.

## PH-INV-001 bis PH-INV-008 – Invarianten
1. Source-ACK nur nach dauerhafter Persistenz.
2. Consumer-ACK regulär nur vorwärts.
3. Offene Transaktionen sind nicht konsumierbar.
4. Rollbacks erzeugen keine regulären Changes.
5. Retention löscht keine benötigten Changes.
6. Jeder Change gehört exakt einer Quelltransaktion.
7. Reihenfolge innerhalb einer Transaktion ist eindeutig.
8. Jeder Change referenziert eine Schema-Version.

## PH-DOD-001 – Definition of Done 0.1
Version 0.1 ist abgeschlossen, wenn INSERT/UPDATE/DELETE transaktionsgetreu erfasst, persistent gespeichert, deterministisch gelesen und nach Crash ohne stille Datenlücke fortgesetzt werden können; WAL erst nach Persistenz bestätigt wird; gefährlicher WAL-Rückstau sichtbar ist; und kritische Anforderungen automatisiert gegen eine reale PostgreSQL-Testinstanz verifiziert werden.
