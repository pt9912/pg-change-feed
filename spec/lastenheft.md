# Lastenheft – PG Change Feed

**Projektname:** PG Change Feed  
**Dokumenttyp:** Lastenheft / Zielbild  
**Vorgehensmodell:** V-Modell-ähnlich  
**Status:** Entwurf  
**Version:** 0.2

## 1. Dokumentzweck

### LH-ZWE-001 – Zweck des Dokuments
Dieses Lastenheft beschreibt die funktionalen und nichtfunktionalen Anforderungen an ein Open-Source-System zur Erfassung, Speicherung und Bereitstellung von Änderungen aus PostgreSQL-Datenbanken. Das System wird im Folgenden als **PG Change Feed** bezeichnet.

Das Lastenheft beschreibt **was** das System leisten soll, nicht **wie** die technische Umsetzung erfolgt.

### LH-ZWE-002 – Zielgruppe
Zielgruppen sind insbesondere:
- Softwareentwickler
- Datenbankadministratoren
- Software- und Systemarchitekten
- DevOps-/Plattform-Teams
- PostgreSQL-Betreiber
- Entwickler von ETL-, Integrations- und Synchronisationslösungen
- Open-Source-Maintainer

## 2. Ausgangssituation

### LH-AUS-001 – PostgreSQL
PostgreSQL stellt mit WAL, Logical Decoding, Logical Replication, Publications und Replication Slots technische Mechanismen zur Erfassung und Übertragung von Datenänderungen bereit.

Diese Mechanismen bilden jedoch primär technische Replikations- und Streaming-Primitiven.

### LH-AUS-002 – Vergleich zu SQL Server CDC
SQL Server CDC bietet eine höherwertige Abstraktion, über die Änderungen aktiviert, gespeichert und positions- bzw. bereichsbezogen abgefragt werden können, ohne dass Consumer unmittelbar das Transaktionslog verarbeiten müssen.

### LH-AUS-003 – Problemstellung
Für PostgreSQL soll eine einfach nutzbare, allgemeine CDC-Abstraktion geschaffen werden, die insbesondere folgende Fähigkeiten kombiniert:
- Aktivieren und Deaktivieren einzelner Tabellen
- persistente Speicherung von Änderungen
- nachvollziehbare Reihenfolge
- positionsbasiertes Lesen
- mehrere unabhängige Consumer
- Retention
- SQL-basierte Administration
- einfache Integration
- möglichst geringe Änderungen an Quellanwendungen

## 3. Zielbild

### LH-ZIE-001 – Übergeordnetes Ziel
PG Change Feed soll eine allgemeine Change-Data-Capture-Abstraktion für PostgreSQL bereitstellen.

### LH-ZIE-002 – PostgreSQL-native Orientierung
Das System soll PostgreSQL-native Fähigkeiten nutzen, ohne eine vollständige technische Kopie von SQL Server CDC zu verlangen.

### LH-ZIE-003 – CDC-Nutzungsmodell
Das angestrebte Nutzungsmodell lautet:

    Tabelle aktivieren
          |
          v
    Änderungen entstehen
          |
          v
    Änderungen werden persistent verfügbar
          |
          v
    Consumer liest ab bekannter Position
          |
          v
    Consumer bestätigt Verarbeitung
          |
          v
    Daten können gemäß Retention bereinigt werden

### LH-ZIE-004 – Open Source
PG Change Feed soll als Open-Source-Projekt nutzbar sein und für den Grundbetrieb keine proprietären Komponenten voraussetzen.

## 4. Systemkontext

### LH-KON-001 – Quelle
Quelle sind eine oder mehrere PostgreSQL-Datenbanken.

### LH-KON-002 – Consumer
Consumer können unter anderem sein:
- ETL-/ELT-Prozesse
- Data-Warehouse-Loader
- Suchindizes
- Cache-Synchronisation
- Replikations- und Integrationsdienste
- Microservices
- kundenspezifische Anwendungen

### LH-KON-003 – Systemgrenze
PG Change Feed umfasst:
- Erfassung relevanter Änderungen
- persistente Bereitstellung
- Konfiguration
- Consumer-Fortschritt
- Retention
- Betriebs- und Diagnoseinformationen

Die fachliche Verarbeitung der Änderungen durch Consumer liegt außerhalb der Systemgrenze.

## 5. Funktionale Anforderungen

### 5.1 CDC-Konfiguration

#### LH-FUN-CFG-001
CDC muss für einzelne Tabellen aktiviert werden können.

#### LH-FUN-CFG-002
CDC muss für einzelne Tabellen deaktiviert werden können.

#### LH-FUN-CFG-003
Der CDC-Status einer Tabelle muss abfragbar sein.

#### LH-FUN-CFG-004
Alle aktivierten Tabellen müssen aufgelistet werden können.

#### LH-FUN-CFG-005
Perspektivisch sollen einzelne Spalten von der Erfassung ausgeschlossen bzw. gezielt ausgewählt werden können.

#### LH-FUN-CFG-006
Das Aktivieren von CDC darf keine Änderungen am Anwendungscode der Quellanwendung erfordern.

### 5.2 Erfassung von Änderungen

#### LH-FUN-CAP-001
INSERT-Operationen müssen erfasst werden.

#### LH-FUN-CAP-002
UPDATE-Operationen müssen erfasst werden.

#### LH-FUN-CAP-003
DELETE-Operationen müssen erfasst werden.

#### LH-FUN-CAP-004
Es müssen genügend Informationen vorliegen, um die logische Reihenfolge der Änderungen eindeutig zu bestimmen.

#### LH-FUN-CAP-005
Änderungen derselben Quelltransaktion müssen als zusammengehörig identifizierbar sein.

#### LH-FUN-CAP-006
Reguläre Consumer dürfen Änderungen erst dann als dauerhaft erfolgreich ansehen, wenn die zugehörige Quelltransaktion erfolgreich abgeschlossen wurde.

#### LH-FUN-CAP-007
Zurückgerollte Änderungen dürfen nicht als erfolgreich committed Changes ausgeliefert werden.

#### LH-FUN-CAP-008
Soweit die Quelle dies zuverlässig ermöglicht, sollen vorherige und neue Werte bereitgestellt werden.

### 5.3 Änderungsdatensatz

#### LH-FUN-DAT-001
Jede Änderung muss eindeutig identifizierbar sein.

#### LH-FUN-DAT-002
Die Quelltabelle muss identifizierbar sein.

#### LH-FUN-DAT-003
Der Operationstyp muss mindestens INSERT, UPDATE und DELETE unterscheiden.

#### LH-FUN-DAT-004
Jede Änderung bzw. Transaktion muss einer eindeutig sortierbaren Quellposition oder daraus abgeleiteten Position zugeordnet werden können.

#### LH-FUN-DAT-005
Die für die konfigurierte CDC-Erfassung relevanten Datenwerte müssen bereitgestellt werden.

#### LH-FUN-DAT-006
Das Datenmodell muss um zusätzliche Metadaten erweiterbar sein.

### 5.4 Lesen von Änderungen

#### LH-FUN-REA-001
Änderungen müssen zwischen zwei Positionen abgefragt werden können.

#### LH-FUN-REA-002
Änderungen müssen ab einer bekannten Position gelesen werden können.

#### LH-FUN-REA-003
Die Anzahl zurückgegebener Änderungen muss begrenzt werden können.

#### LH-FUN-REA-004
Die Sortierung muss bei identischer Eingabe deterministisch sein.

#### LH-FUN-REA-005
Bereits gelesene Änderungen müssen innerhalb der Aufbewahrungszeit erneut gelesen werden können.

#### LH-FUN-REA-006
Änderungen müssen nach Quelltabelle filterbar sein.

### 5.5 Consumer-Verwaltung

#### LH-FUN-CON-001
Benannte Consumer müssen registriert werden können.

#### LH-FUN-CON-002
Mehrere Consumer müssen dieselben Änderungen unabhängig voneinander verarbeiten können.

#### LH-FUN-CON-003
Die zuletzt bestätigte Verarbeitungsposition muss pro Consumer gespeichert werden können.

#### LH-FUN-CON-004
Ein Consumer muss eine erfolgreich verarbeitete Position bestätigen können.

#### LH-FUN-CON-005
Nach einem Neustart muss ein Consumer ab seiner zuletzt bestätigten Position fortsetzen können.

#### LH-FUN-CON-006
Consumer müssen administrativ entfernt werden können.

### 5.6 Speicherung und Retention

#### LH-FUN-RET-001
Erfasste Changes müssen persistent gespeichert werden. Ein Neustart darf bereits dauerhaft gespeicherte Änderungen nicht verlieren.

#### LH-FUN-RET-002
Die Aufbewahrung muss konfigurierbar sein.

#### LH-FUN-RET-003
Zeitbasierte Retention muss möglich sein.

#### LH-FUN-RET-004
Consumer-basierte Retention soll ermöglichen, Änderungen erst zu entfernen, wenn alle relevanten Consumer die erforderliche Position bestätigt haben.

#### LH-FUN-RET-005
Langsame Consumer, die Retention blockieren, müssen erkennbar sein.

#### LH-FUN-RET-006
Unbegrenztes Wachstum der CDC-Daten muss erkannt und betrieblich kontrollierbar sein.

### 5.7 Schemaänderungen

#### LH-FUN-SCH-001
Für CDC relevante Schemaänderungen müssen erkannt werden.

#### LH-FUN-SCH-002
Das Verhalten bei hinzugefügten Spalten muss definiert und für Consumer nachvollziehbar sein.

#### LH-FUN-SCH-003
Das Verhalten bei entfernten Spalten muss definiert sein.

#### LH-FUN-SCH-004
Inkompatible Datentypänderungen müssen erkannt werden; eine stille Fehlinterpretation ist nicht zulässig.

#### LH-FUN-SCH-005
Änderungen sollen einer Schema-Version zugeordnet werden können.

### 5.8 Administration

#### LH-FUN-ADM-001
Zentrale administrative Funktionen sollen über SQL verfügbar sein.

Konzeptionelle Beispiele:

    SELECT cdc.enable_table(...);
    SELECT cdc.disable_table(...);
    SELECT * FROM cdc.tables;
    SELECT * FROM cdc.consumers;

Die exakten technischen Namen werden nicht durch das Lastenheft vorgeschrieben.

#### LH-FUN-ADM-002
Der Betriebsstatus muss abfragbar sein.

#### LH-FUN-ADM-003
Fehlerzustände müssen sichtbar sein.

#### LH-FUN-ADM-004
Der zeitliche bzw. positionsbezogene Abstand zwischen Quelländerung und CDC-Verfügbarkeit muss messbar sein.

#### LH-FUN-ADM-005
Ein Verarbeitungsrückstand muss sichtbar sein.

## 6. Nichtfunktionale Anforderungen

### 6.1 Zuverlässigkeit

#### LH-NFA-REL-001
Erkannte und bestätigte Quelländerungen dürfen nicht still verloren gehen.

#### LH-NFA-REL-002
Das System muss kontrolliert neu gestartet und fortgesetzt werden können.

#### LH-NFA-REL-003
Kann das System Änderungen nicht zuverlässig verarbeiten, muss dies deutlich erkennbar sein.

#### LH-NFA-REL-004
Modell und Schnittstellen sollen idempotente Consumer-Verarbeitung ermöglichen. Generisches Exactly-Once über externe Systeme hinweg wird nicht vorausgesetzt.

### 6.2 Performance

#### LH-NFA-PER-001
Die Auswirkung auf schreibende Quelltransaktionen soll minimiert werden.

#### LH-NFA-PER-002
Das System soll von kleinen Datenbanken bis zu kontinuierlichen Änderungsvolumina skalieren können.

#### LH-NFA-PER-003
Effiziente Batch-Verarbeitung muss möglich sein.

#### LH-NFA-PER-004
Die Latenz zwischen Commit und CDC-Verfügbarkeit soll gering sein und über Benchmarks bewertet werden.

### 6.3 Sicherheit

#### LH-NFA-SEC-001
PostgreSQL-Berechtigungen sollen nach Least-Privilege-Prinzip vergeben werden können.

#### LH-NFA-SEC-002
Administrative CDC-Funktionen und reguläre Lesezugriffe sollen getrennt berechtigbar sein.

#### LH-NFA-SEC-003
CDC-Datenzugriffe müssen über PostgreSQL-Berechtigungen oder einen gleichwertigen Mechanismus beschränkbar sein.

#### LH-NFA-SEC-004
Sensible Spalten sollen von der CDC-Erfassung ausgeschlossen werden können.

### 6.4 Betrieb

#### LH-NFA-OPS-001
Containerisierter Betrieb muss unterstützt werden.

#### LH-NFA-OPS-002
Automatisierte Health Checks müssen möglich sein.

#### LH-NFA-OPS-003
Maschinenlesbare Metriken sollen mindestens CDC-Lag, verarbeitete und ausstehende Changes, Speicherverbrauch, Alter des ältesten Changes, Consumer-Positionen und Fehler abbilden.

#### LH-NFA-OPS-004
Strukturiertes Logging soll unterstützt werden.

#### LH-NFA-OPS-005
Upgrades dürfen persistierte CDC-Daten nicht verlieren.

### 6.5 Portabilität

#### LH-NFA-POR-001
Mehrere aktiv unterstützte PostgreSQL-Major-Versionen sollen unterstützt werden.

#### LH-NFA-POR-002
Linux ist primäre Zielplattform.

#### LH-NFA-POR-003
Ein offizielles bzw. reproduzierbares Docker-/OCI-Deployment soll bereitgestellt werden.

## 7. Schnittstellen

### LH-SST-DB-001
Es muss eine definierte Schnittstelle zur PostgreSQL-Quelle bestehen.

### LH-SST-SQL-001
Konfiguration, Status und Kernlesezugriffe sollen über SQL möglich sein.

### LH-SST-CLI-001
Eine CLI soll Installation, Diagnose und Administration unterstützen können.

### LH-SST-MET-001
Eine Monitoring-kompatible Betriebsschnittstelle muss möglich sein.

### LH-SST-API-001
Die Architektur soll eine spätere HTTP-/gRPC-API ermöglichen, ohne das interne CDC-Modell grundlegend zu verändern.

## 8. Qualitätsziele

### LH-QUA-001
Datenintegrität: sehr hoch.

### LH-QUA-002
Betriebssicherheit: sehr hoch.

### LH-QUA-003
Einfachheit: hoch. PostgreSQL-Entwickler sollen PG Change Feed ohne tiefgehende WAL-Kenntnisse einsetzen können.

### LH-QUA-004
Performance: hoch.

### LH-QUA-005
Transparenz und Diagnosefähigkeit: hoch.

### LH-QUA-006
Erweiterbarkeit: mittel bis hoch.

## 9. Abgrenzung

### LH-ABG-001
PG Change Feed ist kein revisionssicheres Audit-System.

### LH-ABG-002
PG Change Feed ist kein Message Broker wie Kafka oder RabbitMQ.

### LH-ABG-003
Technische Datenänderungen sind nicht automatisch semantische Business Events.

### LH-ABG-004
PG Change Feed ersetzt nicht die native physische oder logische PostgreSQL-Replikation.

### LH-ABG-005
Eine vollständige Kompatibilität mit SQL-Server-CDC-API und -Datenmodell wird nicht garantiert.

## 10. MVP

### LH-MVP-001
Der MVP soll nachweisen, dass eine robuste persistente CDC-Abstraktion mit vertretbarem Betriebsaufwand realisierbar ist.

### LH-MVP-002
Einzelne Tabellen müssen aktiviert und deaktiviert werden können.

### LH-MVP-003
INSERT, UPDATE und DELETE müssen unterstützt werden.

### LH-MVP-004
Changes müssen persistent gespeichert werden.

### LH-MVP-005
Jede Änderung bzw. Transaktion muss eine eindeutig sortierbare Position besitzen.

### LH-MVP-006
Changes müssen ab einer definierten Position gelesen werden können.

### LH-MVP-007
Ein Neustart darf keine dauerhaft erfassten Changes verlieren.

### LH-MVP-008
Eine reproduzierbare lokale Docker-Umgebung muss existieren.

### LH-MVP-009
Ein automatisierter Integrationstest muss mindestens folgenden Ablauf prüfen:

    PostgreSQL starten
          |
    CDC aktivieren
          |
    INSERT
          |
    UPDATE
          |
    DELETE
          |
    Changes lesen
          |
    Reihenfolge und Inhalt prüfen

## 11. Zukünftige Erweiterungen

### LH-FUT-001
Mehrere unabhängige Consumer.

### LH-FUT-002
Erweiterte Retention.

### LH-FUT-003
Schema Evolution.

### LH-FUT-004
Produktionsreife Observability.

### LH-FUT-005
High Availability.

### LH-FUT-006
Exportadapter, beispielsweise Kafka, NATS, RabbitMQ, HTTP/Webhooks oder Object Storage. Diese sind nicht Bestandteil des MVP.

## 12. Lösungneutrale Systemzerlegung

    +-----------------------------------------+
    |              PostgreSQL                 |
    |                                         |
    |  Anwendungstabellen                     |
    |          |                              |
    |          v                              |
    |   Änderungsquelle                       |
    +----------+------------------------------+
               |
               v
    +-----------------------------------------+
    |               PG Change Feed                    |
    |                                         |
    |  Capture                                |
    |     |                                   |
    |     v                                   |
    |  Persistenz                             |
    |     |                                   |
    |     +--------------+                    |
    |     v              v                    |
    | Consumer State   Retention              |
    +----------+------------------------------+
               |
          +----+-----+
          v          v
      Consumer A  Consumer B

Die konkrete technische Änderungsquelle wird im Pflichtenheft bzw. in ADRs festgelegt.

## 13. Traceability

Verbindliche Präfixe:
- LH-ZWE – Dokumentzweck
- LH-AUS – Ausgangssituation
- LH-ZIE – Ziel
- LH-KON – Kontext
- LH-FUN-CFG – Konfiguration
- LH-FUN-CAP – Capture
- LH-FUN-DAT – Datenmodell
- LH-FUN-REA – Lesen
- LH-FUN-CON – Consumer
- LH-FUN-RET – Retention/Persistenz
- LH-FUN-SCH – Schemaänderungen
- LH-FUN-ADM – Administration
- LH-NFA-REL – Zuverlässigkeit
- LH-NFA-PER – Performance
- LH-NFA-SEC – Sicherheit
- LH-NFA-OPS – Betrieb
- LH-NFA-POR – Portabilität
- LH-SST – Schnittstellen
- LH-QUA – Qualität
- LH-ABG – Abgrenzung
- LH-MVP – MVP
- LH-FUT – Zukunft
- LH-ABN – Abnahme

Kennungen werden nach Veröffentlichung nicht wiederverwendet.

Traceability-Beispiel:

    LH-FUN-CAP-002
          |
          +-- PH-FUN-CAP-002
          +-- ADR/ARC-CAP-002
          +-- TST-CAP-002

## 14. MVP-Abnahmekriterien

### LH-ABN-001
INSERT, UPDATE und DELETE werden Ende-zu-Ende erfasst und gelesen.

### LH-ABN-002
Die logische Reihenfolge ist korrekt.

### LH-ABN-003
Rollback-Änderungen werden nicht ausgeliefert.

### LH-ABN-004
Ein Neustart verliert keine dauerhaft erfassten CDC-Daten.

### LH-ABN-005
Aufbewahrte Changes können erneut gelesen werden.

### LH-ABN-006
Die Quellanwendung benötigt keine CDC-spezifischen SQL-Anpassungen.

### LH-ABN-007
Die vollständige Testumgebung ist automatisiert und reproduzierbar.

## 15. Erfolgskriterien

PG Change Feed gilt technisch als erfolgreich, wenn:

1. PostgreSQL-Änderungen zuverlässig erfasst werden.
2. Consumer die zugrunde liegende Änderungsquelle nicht kennen müssen.
3. Consumer Änderungen über stabile Positionen inkrementell verarbeiten können.
4. Neustarts ohne stille Datenverluste möglich sind.
5. die Auswirkung auf die Quelle akzeptabel bleibt.
6. Installation und Betrieb für PostgreSQL-Entwickler nachvollziehbar sind.
7. keine proprietäre Infrastruktur erforderlich ist.

Angestrebte Developer Experience:

    CDC aktivieren
          |
          v
    Änderungen entstehen
          |
          v
    Änderungen lesen
          |
          v
    Position bestätigen
          |
          v
    Verarbeitung fortsetzen
