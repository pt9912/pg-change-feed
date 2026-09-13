# Benutzerhandbuch: PG Change Feed

Version: 1.7
Software-Version: 0.2.0-verdrahtung
Stand: 2026-09-13

## 1. Einleitung

### Zweck der Software

PG Change Feed stellt persistente Change Feeds für bestehende
PostgreSQL-Tabellen bereit: Änderungen (INSERT/UPDATE/DELETE) werden über
Logical Replication erfasst, dauerhaft gespeichert und über SQL lesbar
gemacht — ohne eine zweite Datenbank oder einen Message Broker.

### Zielgruppe

Dieses Handbuch richtet sich an **Betreiber und Integratoren**: Personen,
die den Container betreiben, die erfassten Änderungen per SQL auswerten
oder eine Integration darauf aufbauen. Es gibt keine grafische
Oberfläche — alle Aufgaben laufen über Umgebungsvariablen, `docker
compose`/`make` und SQL-Zugriffe (`psql` oder ein beliebiger
PostgreSQL-Client).

### Voraussetzungen

- Docker und Docker Compose.
- Eine PostgreSQL-Quelldatenbank mit `wal_level=logical`.
- Grundkenntnisse in SQL und im Umgang mit Docker-Umgebungsvariablen.

## 2. Installation und Zugriff

### Systemanforderungen

- PostgreSQL 17 oder 18 als Quelle und als CDC-Speicher.
- Docker-Runtime für den Feed-Container und den Schema-Rollout.

### Image bauen

Es gibt aktuell kein veröffentlichtes Container-Image — das Image wird
lokal aus dem Quellbaum gebaut:

```bash
make image
```

Das Ergebnis ist mit dem Tag `ghcr.io/pt9912/pg-change-feed:dev`
lokal verfügbar; der Image-Digest des letzten Baus steht in
`harness/image-hash.txt`.

### Schema anlegen

Das CDC-Schema (`cdc`) wird über [`d-migrate`](https://github.com/pt9912/d-migrate)
ausgerollt, nicht von Hand:

```bash
make schema-rollout SCHEMA_TARGET="db:<ihre-postgres-dsn>"
```

Das legt alle `cdc.*`-Tabellen, die drei Least-Privilege-Rollen
(`cdc_capture`, `cdc_admin`, `cdc_reader`) und die Diagnose-Views an.

### Zugriff und Rollen

Es gibt keinen eigenen Login für PG Change Feed — der Zugriff läuft über
die PostgreSQL-Verbindung selbst. Der Schema-Rollout legt drei
Gruppenrollen an; alle drei bleiben `NOLOGIN` (Gruppenrollen ohne eigenes
Passwort). Der Feed-Container verbindet sich rollen-spezifisch — jede der
drei Verbindungs-Umgebungsvariablen ([Konfiguration](#5-konfiguration))
trägt eine eigene, anmeldefähige Login-Identität, die Sie der passenden
Gruppenrolle zuweisen:

| Rolle | Zweck | Umgebungsvariable |
|---|---|---|
| `cdc_capture` | Erfassungspfad des Feed-Containers (Store-Adapter, Replication-Stream) | `CDC_CAPTURE_DSN` |
| `cdc_admin` | Verwaltungszugriff (Registrierung von Quellen und Tabellen, Heartbeat, `register-consumer`/`acknowledge-consumer`, Retention-Löschausführung) | `CDC_ADMIN_DSN` |
| `cdc_reader` | Nur-Lese-Zugriff auf die Diagnose- und Lese-Views (`cdc.active_tables`, `cdc.consumer_status`, `cdc.changes`, `cdc.metrics`, `cdc.heartbeat`, `cdc.retention_blockers`) — trägt auch `--healthcheck` und `diagnose` (siehe [Diagnose ausführen](#diagnose-ausführen)) | `CDC_READER_DSN` |

```sql
CREATE ROLE feed_capture_login LOGIN PASSWORD '<geheim>' IN ROLE cdc_capture;
CREATE ROLE feed_admin_login LOGIN PASSWORD '<geheim>' IN ROLE cdc_admin;
CREATE ROLE feed_reader_login LOGIN PASSWORD '<geheim>' IN ROLE cdc_reader;
```

Für einen bereits bestehenden Login genügt die Mitgliedschaft:

```sql
GRANT cdc_reader TO ihr_login;
```

**Betriebs-Hinweis (`REPLICATION`-Attribut):** PostgreSQL vererbt
Rollen-*Attribute* (`LOGIN`, `REPLICATION`, …) nicht über Mitgliedschaft —
nur Objekt-Privilegien (`GRANT SELECT`/`INSERT`/…) tun das. Das
`REPLICATION`-Attribut liegt auf der Gruppenrolle `cdc_capture` selbst;
eine Login-Identität, die nur `IN ROLE cdc_capture` erhält, kann dadurch
**keine** Replication-Verbindung aufbauen. Setzen Sie das Attribut
zusätzlich direkt auf die Login-Identität hinter `CDC_CAPTURE_DSN`:

```sql
ALTER ROLE feed_capture_login REPLICATION;
```

## 3. Erste Schritte

### Schneller Einstieg

**Voraussetzung:** Schema ausgerollt (siehe oben), eine physische
Quelltabelle existiert bereits.

**Vorgehen:**

1. Quelle registrieren (einmalig je Quelle):

   ```sql
   INSERT INTO cdc.source (source_id, name) VALUES ('meine-quelle', 'Produktionsdatenbank');
   ```

2. Feed-Container mit den Aktivierungs-Umgebungsvariablen starten (siehe
   [Konfiguration](#5-konfiguration)):

   ```bash
   docker run --rm \
     -e CDC_CAPTURE_DSN="postgres://feed_capture_login:pass@host:5432/db?sslmode=disable" \
     -e CDC_ADMIN_DSN="postgres://feed_admin_login:pass@host:5432/db?sslmode=disable" \
     -e CDC_READER_DSN="postgres://feed_reader_login:pass@host:5432/db?sslmode=disable" \
     -e CDC_SOURCE_ID="meine-quelle" \
     -e CDC_PUBLICATION="pub_meine_quelle" \
     -e CDC_SLOT="slot_meine_quelle" \
     -e CDC_TABLES="public.orders=tbl-orders:sv-orders-1" \
     ghcr.io/pt9912/pg-change-feed:dev
   ```

3. Änderungen an `public.orders` lesen:

   ```sql
   SELECT * FROM cdc.changes WHERE source_id = 'meine-quelle' ORDER BY commit_position, sequence;
   ```

**Ergebnis:** Der Container aktiviert `public.orders` beim Start
automatisch (Publication, Bindungs- und Schema-Version-Zeile), nimmt den
Replication-Stream auf und schreibt jede committed Änderung nach
`cdc.change`.

## 4. Aufgaben

### Quelle registrieren

**Voraussetzung:** Zugriff auf die CDC-Speicherdatenbank mit
Schreibrecht (`cdc_admin` oder höher).

**Vorgehen:**

```sql
INSERT INTO cdc.source (source_id, name) VALUES ('<quelle-id>', '<lesbarer-name>');
```

**Ergebnis:** Die Quelle ist bekannt; Tabellen dieser Quelle können jetzt
aktiviert werden. Eine bereits registrierte Quelle bleibt unverändert
(`source_id` ist Primärschlüssel).

### Tabelle aktivieren

**Voraussetzung:** Die Quelle ist registriert; die physische Tabelle
existiert; `REPLICA IDENTITY` ist wie benötigt gesetzt (Standard genügt,
außer Sie brauchen alte Zeilenwerte bei UPDATE/DELETE — dann `REPLICA
IDENTITY FULL`).

**Vorgehen:** Die Aktivierung trägt der Feed-Container selbst beim
Start, gesteuert über die Umgebungsvariable `CDC_TABLES` (Format:
`schema.tabelle=tabelle-id:schema-version-id`, mehrere Einträge durch
Komma getrennt):

```bash
CDC_TABLES="public.orders=tbl-orders:sv-orders-1,public.customers=tbl-customers:sv-customers-1"
```

**Ergebnis:** Beim Start legt der Container die Publication-Mitgliedschaft
sowie die Bindungs- und Schema-Version-Zeile an. Eine bereits aktivierte
Tabelle bleibt unverändert (idempotent).

### Tabelle live aktivieren

**Voraussetzung:** Die Quelle ist registriert (siehe oben); die physische
Tabelle existiert; `REPLICA IDENTITY` ist wie benötigt gesetzt (siehe
[Tabelle aktivieren](#tabelle-aktivieren)); eine Login-Identität mit
`cdc_admin`-Mitgliedschaft, verbunden über `CDC_ADMIN_DSN` (siehe [Zugriff
und Rollen](#zugriff-und-rollen)); der Feed-Container läuft bereits für
die betroffene Quelle.

**Vorgehen:** Anders als die `CDC_TABLES`-gesteuerte Aktivierung beim
Containerstart (siehe oben) lässt sich eine Tabelle zusätzlich **während
des laufenden Betriebs** aktivieren, ohne den Container neu zu starten:

```sql
SELECT cdc.enable_table('<source_id>', '<schema>', '<tabelle>');
```

**Ergebnis:** Der Aufruf ist asynchron: Er schreibt einen Antrag nach
`cdc.administration_request` (Status `pending`) und gibt dessen Kennung
zurück — die Tabelle ist damit noch nicht aktiv. Die Administrations-
Goroutine des laufenden Feed-Containers verarbeitet offene Anträge
(`LISTEN`/`NOTIFY`-Weckung mit periodischem Fallback-Poll) und trägt bei
Erfolg die Publication-Mitgliedschaft sowie die laufende Erfassungs-Bindung
nach — ohne Neustart. Den Fortschritt prüfen Sie über den Antrags-Status:

```sql
SELECT status, error_message
FROM cdc.administration_request
WHERE administration_request_id = '<zurückgegebene-id>';
```

`status` wechselt von `pending` zu `applied` (Erfolg) oder `failed`
(Fehlertext in `error_message`). Erst nach `applied` erfasst der laufende
Prozess Änderungen an der Tabelle.

### Tabelle deaktivieren

**Voraussetzung:** wie bei der Live-Aktivierung — die physische Tabelle
existiert, `REPLICA IDENTITY` ist wie benötigt gesetzt (siehe [Tabelle
aktivieren](#tabelle-aktivieren)), `cdc_admin`-Mitgliedschaft über
`CDC_ADMIN_DSN`, der Feed-Container läuft bereits.

**Vorgehen:** Derselbe Antrags-Weg, spiegelbildlich:

```sql
SELECT cdc.disable_table('<source_id>', '<schema>', '<tabelle>');
```

**Ergebnis:** Der Aufruf ist asynchron: Er schreibt einen Antrag nach
`cdc.administration_request` (Status `pending`) und gibt dessen Kennung
zurück — die Tabelle wird damit noch nicht deaktiviert. Die
Administrations-Goroutine des laufenden Feed-Containers verarbeitet offene
Anträge (`LISTEN`/`NOTIFY`-Weckung mit periodischem Fallback-Poll) und
entzieht bei Erfolg die Publication-Mitgliedschaft — ohne Neustart. Den
Fortschritt prüfen Sie über den Antrags-Status:

```sql
SELECT status, error_message
FROM cdc.administration_request
WHERE administration_request_id = '<zurückgegebene-id>';
```

`status` wechselt von `pending` zu `applied` (Erfolg) oder `failed`
(Fehlertext in `error_message`). Nach `applied` endet die Erfassung dieser
einen Tabelle; der Feed-Container läuft für alle übrigen aktivierten
Tabellen unverändert weiter — kein Neustart, keine Unterbrechung des
laufenden Prozesses.

### Aktivierte Tabellen auflisten

```sql
SELECT * FROM cdc.active_tables;
```

**Ergebnis:** Eine Zeile je aktivierter Tabelle mit der aktuell
gebundenen Schema-Version.

### Consumer registrieren

**Voraussetzung:** eine Login-Identität für `CDC_ADMIN_DSN` (Mitglied von
`cdc_admin`, siehe [Zugriff und Rollen](#zugriff-und-rollen)). Der
Sondermodus liest dieselben Umgebungs-Vorbedingungen wie der reguläre
Lauf, verbindet sich für den Zugriffsweg selbst aber ausschließlich über
`CDC_ADMIN_DSN`.

**Vorgehen:** Der Feed-Container trägt einen Sondermodus, der den Aufruf
über `RegisterConsumerUseCase` führt, statt die CDC-Speichertabellen
direkt zu schreiben — derselbe Image-Tag wie der Daemon, als einmaliger,
kurzlebiger Lauf statt als Dauerdienst:

```bash
docker run --rm -e CDC_CAPTURE_DSN -e CDC_ADMIN_DSN -e CDC_READER_DSN \
  -e CDC_SOURCE_ID -e CDC_PUBLICATION -e CDC_SLOT -e CDC_TABLES \
  ghcr.io/pt9912/pg-change-feed:dev register-consumer <name>
```

In der Compose-Umgebung: `docker compose run --rm pg-change-feed
register-consumer <name>`. `<name>` trägt zugleich Kennung und Namen des
Consumers.

**Ergebnis:** Der Consumer ist registriert und kann fortan lesen und
Positionen bestätigen (`LH-FA-CON-001`). Ein bereits registrierter Name
bleibt unverändert; der Lauf meldet das über eine eigene Ausgabe-Zeile,
Exit-Code bleibt 0 (Idempotenz, `LH-FA-CON-001` Boundary).

### Position bestätigen

**Voraussetzung:** Der Consumer ist registriert (siehe oben); eine
Login-Identität für `CDC_ADMIN_DSN` wie bei der Registrierung — derselbe
Zugriffsweg verbindet sich ausschließlich über `CDC_ADMIN_DSN`.

**Vorgehen:** Derselbe Sondermodus-Mechanismus wie bei der Registrierung
führt den Aufruf über `AcknowledgeConsumerUseCase`, statt die
Consumer-Positions-Tabelle direkt zu schreiben — derselbe Image-Tag wie
der Daemon, als einmaliger, kurzlebiger Lauf statt als Dauerdienst:

```bash
docker run --rm -e CDC_CAPTURE_DSN -e CDC_ADMIN_DSN -e CDC_READER_DSN \
  -e CDC_SOURCE_ID -e CDC_PUBLICATION -e CDC_SLOT -e CDC_TABLES \
  ghcr.io/pt9912/pg-change-feed:dev acknowledge-consumer <consumer-id> <position>
```

In der Compose-Umgebung: `docker compose run --rm pg-change-feed
acknowledge-consumer <consumer-id> <position>`. `<position>` trägt den
Offset innerhalb der konfigurierten Quelle (`CDC_SOURCE_ID`) — derselbe
Wert, den `cdc.changes.commit_position` für die zuletzt verarbeitete
Änderung trägt.

**Ergebnis:** Die Position ist bestätigt und fortgeschrieben
(`LH-FA-CON-004`). Die Wiederholung derselben Position bleibt ohne Wirkung
(Idempotenz, `LH-FA-CON-004` Boundary); ein echter Rückschritt (eine
Position vor dem bereits bestätigten Stand) wird abgelehnt, der
gespeicherte Fortschritt bleibt unverändert — die Vorwärts-Invariante gilt
für jeden Aufruf über diesen Zugriffsweg (`LH-FA-CON-004.a`).

### Änderungen lesen

```sql
SELECT source_id, commit_position, change_id, schema_name, table_name,
       operation, old_data, new_data, committed_at
FROM cdc.changes
WHERE source_id = '<quelle-id>' AND commit_position > <letzte-gelesene-position>
ORDER BY commit_position, sequence
LIMIT 500;
```

**Ergebnis:** Jede Zeile ist eine committed Änderung in Anhang-Reihenfolge.
`old_data`/`new_data` sind `jsonb`; bei `INSERT` ist `old_data` NULL, bei
`DELETE` ist `new_data` NULL.

### Aufbewahrung (Retention)

Der Feed-Container bereinigt `cdc.change`-Zeilen automatisch über einen
Hintergrundzug — kein CLI-Befehl und keine manuelle Auslösung nötig. Der
Zug läuft alle 10 Sekunden und entfernt je Durchlauf genau die Zeilen, die
zwei Bedingungen zugleich erfüllen: ihr Alter (gemessen am realen
Quell-Commit-Zeitpunkt) erreicht mindestens 24 Stunden, **und** jeder
Consumer, der für die Quelle bereits einmal bestätigt hat, hat eine
Position an oder hinter der jeweiligen Zeile bestätigt. Beide Werte
(Takt, Mindestalter) sind fest im Feed-Container hinterlegt; eine
Laufzeit-Konfiguration über Umgebungsvariablen oder die YAML-Datei
existiert dafür nicht.

**Betriebs-Hinweis (Consumer-Bindung):** Ein registrierter, aber gegen
eine Quelle noch nie bestätigender Consumer blockiert die Bereinigung
dieser Quelle **nicht** — er zählt erst ab seiner ersten Bestätigung
(`acknowledge-consumer` bzw. [Position bestätigen](#position-bestätigen))
als schützenswert. Ein neu angebundener Consumer, der vor seinem ersten
Lesezugriff vor Bereinigung geschützt sein soll, sollte deshalb einmal
bestätigen (auch die dokumentierte Anfangsposition genügt), bevor er mit
dem Lesen beginnt.

Direkter SQL-Zugriff auf die betroffenen Tabellen bleibt zulässig
(`cdc_admin`-Mitgliedschaft, verbunden über `CDC_ADMIN_DSN`):

```sql
SELECT change_id, committed_at FROM cdc.changes WHERE source_id = '<quelle-id>' ORDER BY commit_position;
```

Eine Zeile, die dort nicht mehr erscheint, wurde bereits bereinigt.

### Blockierende Consumer erkennen

`cdc.retention_blockers` zeigt je Quelle den Consumer, dessen bestätigte
Position aktuell die Löschgrenze der Bereinigung trägt (`LH-FA-RET-005`) —
also genau den Consumer, den `RunRetentionUseCase` als
weitesten-zurückliegend behandelt, bevor er weitere Zeilen freigibt:

```sql
SELECT consumer_id, name, acknowledged_position, backlog
FROM cdc.retention_blockers
WHERE source_id = '<quelle-id>';
```

**Ergebnis:** Höchstens eine Zeile je Quelle — `backlog` trägt den Abstand
zwischen der bestätigten Position dieses Consumers und der letzten
Commit-Position der Quelle. Eine Quelle ohne Zeile hat aktuell keinen
Consumer-Blocker (entweder hat noch nie ein Consumer gegen sie bestätigt,
oder alle bestätigenden Consumer sind bereits auf Höhe der letzten
Commit-Position). Ein registrierter, aber gegen diese Quelle noch nie
bestätigender Consumer erscheint hier nicht — dieselbe Abwesenheits-Lesart
wie beim [Betriebs-Hinweis](#aufbewahrung-retention) oben: Schutz vor
Bereinigung entsteht erst mit seiner ersten Bestätigung. Die Sicht macht
nur sichtbar, was `RunRetentionUseCase` bereits real entscheidet — sie
berechnet die Freigabe nicht neu.

### Betriebsstatus prüfen

Der Container meldet seinen Zustand über den Docker-Healthcheck
(`docker inspect --format '{{.State.Health.Status}}' <container>`), der
intern `--healthcheck` ausführt und das Lebenszeichen der Quelle über
`CDC_READER_DSN` prüft (`SELECT`-Grant der Rolle `cdc_reader` auf
`cdc.heartbeat`). Direkter SQL-Zugriff:

```sql
SELECT source_id, heartbeat_at, age_seconds, error_class FROM cdc.heartbeat;
```

**Ergebnis:** `error_class` ist `NULL` im Normalbetrieb; ein gesetzter
Wert nennt die zuletzt beobachtete Fehlerklasse (siehe
[Fehlerbehebung](#6-fehlerbehebung)). Ein Lebenszeichen gilt als
veraltet, wenn `age_seconds` mehr als das Dreifache des Schreibtakts
(15 Sekunden bei einem Takt von 5 Sekunden) überschreitet.

### Metriken lesen

```sql
SELECT metric_name, label, value FROM cdc.metrics;
```

Verfügbare Kennzahlen: `cdc_transactions_total`, `cdc_changes_processed`,
`cdc_oldest_change_age_seconds`, `cdc_capture_lag` (Abstand zwischen der
letzten Quelländerung und der CDC-Verfügbarkeit, gemessen über den
Commit-Zeitstempel aus dem WAL), `cdc_consumer_position` und
`cdc_consumer_lag` je registriertem Consumer.

`cdc_wal_retention_bytes` (WAL-Rückstand des Capture-Slots, `SPEC-009`)
steht **nicht** in `cdc.metrics`: Die Erhebung braucht Systemkatalog-Zugriffe
außerhalb des `cdc`-Schemas (`pg_replication_slots`, `IDENTIFY_SYSTEM`), die
die Least-Privilege-Fläche von `cdc_reader` unnötig erweitern würden — siehe
[WAL-Rückstand prüfen](#wal-rückstand-prüfen).

### Diagnose ausführen

Statt der beiden SQL-Abfragen oben einzeln zu stellen, liest der
`diagnose`-Sondermodus dieselben Views (`cdc.heartbeat`, `cdc.metrics`) über
`CDC_READER_DSN` und gibt eine menschenlesbare Zusammenfassung aus
(`LH-FA-SST-003`, deckt `LH-FA-ADM-002`…`005`) — derselbe Image-Tag wie der
Daemon, als einmaliger, kurzlebiger Lauf statt als Dauerdienst:

```bash
docker run --rm -e CDC_CAPTURE_DSN -e CDC_ADMIN_DSN -e CDC_READER_DSN \
  -e CDC_SOURCE_ID -e CDC_PUBLICATION -e CDC_SLOT -e CDC_TABLES \
  ghcr.io/pt9912/pg-change-feed:dev diagnose
```

In der Compose-Umgebung: `docker compose run --rm pg-change-feed diagnose`,
oder gegen einen bereits laufenden Feed-Container: `docker exec <container>
/pg-change-feed diagnose`.

**Ausgabe (Beispiel):**

```text
pg-change-feed diagnose: Quelle "src-mvp"
  Betriebsstatus (LH-FA-ADM-002): Lebenszeichen vor 1.203s
  Fehlerzustand (LH-FA-ADM-003): keiner (Normalbetrieb)
  CDC-Abstand cdc_capture_lag (LH-FA-ADM-004): 0.087s
  Verarbeitungsrückstand cdc_consumer_lag je Consumer (LH-FA-ADM-005, nur Consumer mit mindestens einer bestätigten Position):
    cli-e2e-consumer: 0
```

**Ergebnis:** Wie bei den Rohwerten der beiden Views trifft der Befehl keine
Schwellenwert-Entscheidung (`SPEC-007` bleibt Sache des lesenden Systems) und
der Prozess-Ausgang trägt nur den Lese-Erfolg — ein gemeldeter Fehlerzustand
oder Rückstand ist Berichtsinhalt, kein Befehlsfehler (Ausgang bleibt 0). Ein
Consumer ohne je bestätigte Position erscheint nicht in der Rückstands-Liste
(dieselbe Grenze wie bei `cdc_consumer_lag` in [Metriken
lesen](#metriken-lesen)); ein Consumer mit bestätigter Position, dessen
Quelle noch nie eine Transaktion trug, erscheint mit dem Text „unbekannt"
statt einem irreführenden Rückstand von 0.

### WAL-Rückstand prüfen

Der Feed-Container misst den WAL-Rückstand des Capture-Slots periodisch
(gebunden an denselben Takt wie das Lebenszeichen, siehe
[Grenzwerte](#grenzwerte)) über die Rolle `cdc_capture` und protokolliert ihn
strukturiert:

```json
{"msg": "replication: WAL-Rückstand gemessen", "metric": "cdc_wal_retention_bytes", "bytes": 12345}
```

**Ergebnis:** `bytes` wächst, solange der Capture-Slot inaktiv ist und die
Quelle weiterschreibt (z. B. während eines Verbindungsabbruchs) — ein
dauerhaft wachsender Wert ist ein Warnsignal für WAL-Erschöpfung auf der
Quelle. Der Feed-Container vergleicht den gemessenen Wert bei jedem Takt
gegen zwei Schwellen (`SPEC-013`, `ADR-0049`):

- **Unterhalb 100 MiB:** unauffällig — derselbe Log-Eintrag wie oben.
- **Zwischen 100 MiB und 1 GiB:** kontrollierte Fortsetzung — der
  Capture-Betrieb läuft unverändert weiter, sichtbar über eine
  Log-Warnung:
  ```json
  {"level": "WARN", "msg": "replication: WAL-Rückstand über Warnschwelle — kontrollierte Fortsetzung", "metric": "cdc_wal_retention_bytes", "bytes": 150000000, "threshold_bytes": 104857600}
  ```
- **Oberhalb 1 GiB:** kontrollierter Abbruch — der Container protokolliert
  den Fehler, klassifiziert den Lauf als `replication` (Transport-/
  Verbindungsstörung, siehe [Fehlerklassen](#fehlerklassen)) und beendet
  sich (Ausgang 1). Ein Neustart setzt den Stream über den bestehenden
  Slot-Stand fort.

Diese Schwellen betreffen ausschließlich Transport-/Verbindungsstörungen
(`receive.ErrReplication`/`outbound.ErrReplication`); eine
Stream-Ordnungs-Verletzung bricht unabhängig vom WAL-Rückstand weiterhin
sofort ab (siehe [Fehlerklassen](#fehlerklassen)).

### Schema aktualisieren

Änderungen am neutralen Schema (`tools/schema/schema.yaml`) rollen Sie
mit demselben Befehl wie bei der Ersteinrichtung erneut aus:

```bash
make schema-rollout SCHEMA_TARGET="db:<ihre-postgres-dsn>"
```

Der Lauf erzeugt einen Pflicht-Report (`tools/schema/plan.yaml`) und ein
Rollback-Artefakt (`tools/schema/down.sql`).

## 5. Konfiguration

### Umgebungsvariablen des Feed-Containers

| Variable | Pflicht | Bedeutung |
|---|---|---|
| `CDC_CAPTURE_DSN` | ja | Verbindung über die Rolle `cdc_capture` (Store-Adapter, Replication-Stream) |
| `CDC_ADMIN_DSN` | ja | Verbindung über die Rolle `cdc_admin` (Tabellen-Aktivierung, Heartbeat, `register-consumer`/`acknowledge-consumer`) |
| `CDC_READER_DSN` | ja | Verbindung über die Rolle `cdc_reader` (`--healthcheck`, `diagnose`) |
| `CDC_SOURCE_ID` | ja | Kennung der Quelle (muss in `cdc.source` registriert sein) |
| `CDC_PUBLICATION` | ja | Name der PostgreSQL-Publication |
| `CDC_SLOT` | ja | Name des Logical-Replication-Slots |
| `CDC_TABLES` | ja, falls keine Konfigurationsdatei dieselbe Aktivierung trägt | Aktivierte Tabellen, Format `schema.tabelle=tabelle-id:schema-version-id`, kommagetrennt |
| `CDC_LOG_LEVEL` | nein | Log-Level des strukturierten JSON-Loggers (Default `info`) |
| `CDC_CONFIG_FILE` | nein | Pfad zu einer optionalen YAML-Konfigurationsdatei (siehe unten) |

Fehlt eine Pflichtvariable und liefert auch keine Konfigurationsdatei
einen Wert für dasselbe Feld, startet der Container nicht (Fehlerklasse
`configuration`).

### Optionale YAML-Konfigurationsdatei (`CDC_CONFIG_FILE`)

Additiv zu den Umgebungsvariablen (`ADR-0052`, `SPEC-016`): Ist
`CDC_CONFIG_FILE` gesetzt, liest der Container zusätzlich eine
YAML-Datei unter diesem Pfad (read-only in den Container gemountet). Jede
gesetzte Umgebungsvariable überschreibt das gleichnamige Feld der Datei
einzeln — Env-Var schlägt Datei, Feld für Feld. Ist `CDC_CONFIG_FILE`
nicht gesetzt, ändert sich am Env-only-Betrieb oben nichts.

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

**Wichtig — Secrets bleiben env-var-exklusiv:** `CDC_CAPTURE_DSN`,
`CDC_ADMIN_DSN` und `CDC_READER_DSN` dürfen in dieser Datei **nicht**
vorkommen (Schlüssel `capture_dsn`/`admin_dsn`/`reader_dsn`). Ein Treffer
bricht das Laden ab (Fehlerklasse `configuration`) — eine
Konfigurationsdatei landet typischerweise in Kanälen (Repository,
ConfigMap, Backup), die für Zugangsdaten nicht vorgesehen sind. Ein
unbekannter Schlüssel bricht das Laden ebenfalls ab (striktes Decoding).

Ist sowohl `CDC_TABLES` als auch `tables` in der Datei gesetzt, schlägt
`CDC_TABLES` die gesamte Datei-Tabellenliste vollständig — es findet keine
Vermischung einzelner Tabellen aus beiden Quellen statt (`tables` gilt als
ein Feld, nicht als Menge einzeln überschreibbarer Einträge).

## 6. Fehlerbehebung

### Fehlerklassen

Jeder Fehler des Feed-Containers gehört zu einer von sieben stabilen
Klassen (`ADR-0023`, `SPEC-008`):

| Klasse | Bedeutung | Verhalten |
|---|---|---|
| `transient` | vorübergehend nicht verfügbare Quelle/Speicher | Erneuter Versuch mit begrenztem Backoff; deklariert, aktuell von keinem Adapter konstruiert |
| `configuration` | ungültige oder fehlende Umgebungsvariable | Kein Start, sichtbarer Fehler |
| `permission` | fehlende Berechtigung | Sichtbarer Fehler, kein stiller Retry; deklariert, aktuell von keinem Adapter konstruiert |
| `schema` | eine Replikationsnachricht ist nicht sicher interpretierbar (z. B. TRUNCATE, unbekannter Nachrichtentyp) | Sichtbarer Fehler, kein stilles Überspringen |
| `storage` | Persistenzfehler | Kein Source-ACK, damit keine Änderung verloren geht |
| `replication` | zwei Unterarten (`ADR-0049`): **Stream-Ordnungs-Verletzung** (z. B. Commit ohne offene Transaktion) oder **Transport-/Verbindungsstörung** (Verbindungsaufbau, Slot, Keepalive, Quell-Bestätigung) | Stream-Ordnungs-Verletzung: sofortiger, sichtbarer Abbruch, unabhängig vom WAL-Rückstand. Transport-/Verbindungsstörung: Schwellen-Überwachung über den WAL-Rückstand (siehe [WAL-Rückstand prüfen](#wal-rückstand-prüfen)) — kontrollierte Fortsetzung unterhalb 1 GiB, sichtbarer Abbruch darüber |
| `internal` | unerwarteter interner Fehler, der keiner anderen Klasse zuzuordnen ist | Sichtbarer Fehler; realer Fallback für jeden nicht erkannten Fehler |

`transient` und `permission` gehören zur deklarierten Menge der sieben
Klassen, werden aber von keinem Adapter aktuell konstruiert — sie sind
heute nicht beobachtbar. `internal` ist dagegen der real erreichbare
Fallback-Zweig: Jeder Fehler, der keiner der übrigen sechs Klassen
zugeordnet werden kann, fällt auf `internal` zurück.

Ein Fehler jeder Klasse beendet den Container-Prozess mit Ausgang 1; der
zuletzt beobachtete Fehlerzustand wird zusätzlich in
`cdc.heartbeat.error_class` festgehalten und bei einem erfolgreichen
Neustart automatisch wieder gelöscht.

### Container startet nicht

**Ursache:** eine Pflicht-Umgebungsvariable fehlt oder ist falsch
formatiert.

**Lösung:**

1. Prüfen Sie die Container-Logs (`docker logs <container>`) — die
   Fehlermeldung nennt die betroffene Variable.
2. Prüfen Sie das Format von `CDC_TABLES`
   (`schema.tabelle=tabelle-id:schema-version-id`, kommagetrennt).
3. Starten Sie den Container erneut.

### Container startet, erfasst aber keine Änderungen

**Ursache:** Die Quelle ist nicht in `cdc.source` registriert — die
Tabellen-Aktivierung schlägt dann an der Fremdschlüsselbedingung fehl.

**Lösung:**

1. Prüfen Sie, ob die Quelle registriert ist:
   `SELECT * FROM cdc.source WHERE source_id = '<quelle-id>';`
2. Falls nicht: registrieren (siehe [Quelle registrieren](#quelle-registrieren))
   und den Container neu starten.

### Neustart nach einem Fehler

Der Container startet nicht automatisch neu (`restart: "no"` im
mitgelieferten `compose.yaml`) — der Neustart liegt beim Aufrufer
(Orchestrator, Supervisor). Ein Neustart setzt am zuletzt bestätigten
Slot-Stand fort; keine bereits gespeicherte Änderung geht dabei verloren.

## 7. FAQ

**Kann ich mehrere Quellen mit einem Container betreiben?**
Nein — ein Container-Lauf bindet genau eine Quelle (`CDC_SOURCE_ID`).
Mehrere Quellen brauchen mehrere Container-Instanzen.

**Was passiert, wenn ich eine bereits aktivierte Tabelle erneut in
`CDC_TABLES` nenne?**
Nichts — die Aktivierung ist idempotent.

**Werden gelöschte Zeilen mit alten Werten erfasst?**
Nur, wenn die Quelltabelle `REPLICA IDENTITY FULL` trägt; sonst ist
`old_data` bei DELETE nur mit den Schlüsselspalten gefüllt.

## 8. Glossar

| Begriff | Bedeutung |
|---|---|
| Quelle (Source) | Eine überwachte PostgreSQL-Datenbank |
| Aktivierte Tabelle | Eine Tabelle, deren Änderungen erfasst werden |
| Commit-Position | Die WAL-Position, an der eine Quelltransaktion committed wurde |
| Change | Eine einzelne erfasste INSERT-/UPDATE-/DELETE-Zeile |
| Schema-Version | Kennung der Spaltenform einer aktivierten Tabelle zum Zeitpunkt eines Changes |
| Consumer | Ein benannter, unabhängiger Leser des Change Feeds |
| Slot | Der PostgreSQL Logical-Replication-Slot, der den Fortsetzungspunkt trägt |
| Lebenszeichen (Heartbeat) | Periodischer Nachweis, dass der Capture-Prozess aktiv ist |

## 9. Anhang

### Grenzwerte

- Lebenszeichen-Takt: 5 Sekunden; als veraltet gilt ein Lebenszeichen
  nach mehr als 15 Sekunden (Faktor 3).
- WAL-Rückstand-Messtakt: derselbe Takt wie das Lebenszeichen (5 Sekunden).
- WAL-Rückstand-Schwellen (`SPEC-013`, `ADR-0049`): Warnschwelle 100 MiB,
  Fehlerschwelle 1 GiB — betrifft ausschließlich die Fehlerklasse
  `replication`, Unterart Transport-/Verbindungsstörung (siehe
  [Fehlerklassen](#fehlerklassen)).
- Ein Container-Lauf bindet genau eine Quelle.
- Ein Container-Lauf hält zwei gleichzeitige Replication-Protokoll-
  Verbindungen zur Quelle (Stream-Adapter, WAL-Rückstand-Messung) — beide
  zählen gegen `max_wal_senders` der Quell-Instanz.

### Support und Kontakt

Fragen und Fehler bitte über das Projekt-Repository melden.

### Lizenz

MIT — siehe `LICENSE`.

### Änderungshistorie

| Version | Datum | Änderung |
|---|---|---|
| 1.0 | 2026-09-12 | Erste Fassung |
| 1.1 | 2026-09-12 | Rollen-spezifische DSN-Verdrahtung (`ADR-0047`): `CDC_SOURCE_DSN` ersatzlos ersetzt durch `CDC_CAPTURE_DSN`/`CDC_ADMIN_DSN`/`CDC_READER_DSN`, Betriebs-Hinweis zum `REPLICATION`-Attribut ergänzt |
| 1.2 | 2026-09-12 | Fehlerklassen-Tabelle (§6) auf alle sieben Klassen aus `ADR-0023`/`SPEC-008` vervollständigt (`transient`, `permission`, `internal` ergänzt) |
| 1.3 | 2026-09-12 | WAL-Rückstand-Metrik `cdc_wal_retention_bytes` (`SPEC-009`) ergänzt: periodische Messung, strukturierte Log-Ausgabe, Abgrenzung gegen `cdc.metrics` |
| 1.4 | 2026-09-12 | Fehlerklasse `replication` auf zwei Unterarten präzisiert (`ADR-0049`): Stream-Ordnungs-Verletzung bleibt sofortiger Abbruch, Transport-/Verbindungsstörung trägt jetzt die Schwellen-Überwachung über den WAL-Rückstand (Warn 100 MiB, Fehler 1 GiB) mit kontrollierter Fortsetzung/Abbruch |
| 1.5 | 2026-09-13 | Neuer `diagnose`-Sondermodus ergänzt (`LH-FA-SST-003`, deckt `LH-FA-ADM-002`…`005`, slice-038): §4 „Diagnose ausführen", `cdc_reader`-Zeile und `CDC_READER_DSN`-Zeile aktualisiert |
| 1.6 | 2026-09-13 | Optionale YAML-Konfigurationsdatei (`CDC_CONFIG_FILE`, `ADR-0052`, `SPEC-016`, slice-041) ergänzt: §5 neue Unterüberschrift, Env-Var-Tabelle um `CDC_CONFIG_FILE` erweitert, `CDC_TABLES`-Pflichtangabe präzisiert |
| 1.7 | 2026-09-13 | SQL-Administration nachdokumentiert (`LH-FA-ADM-001`, `LH-FA-CFG-002`, `ADR-0050`, slice-036, slice-037, slice-042): §4 zwei neue Abschnitte „Tabelle live aktivieren" und „Tabelle deaktivieren" (`cdc.enable_table`/`cdc.disable_table`, asynchrone Antrags-Queue, Status-Polling) |
