# Benutzerhandbuch: PG Change Feed

Version: 1.62
Software-Version: siehe `docs/user/version.md`
Stand: 2026-09-25

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
| `cdc_capture` | Erfassungspfad des Feed-Containers (Store-Adapter, Replication-Stream) und Ausführung eines Backfills (Run-Zustand fortschreiben, Bestand im Snapshot lesen und schreiben, siehe [Bestand als Backfill überführen](#bestand-als-backfill-überführen)) | `CDC_CAPTURE_DSN` |
| `cdc_admin` | Verwaltungszugriff (Registrierung von Quellen und Tabellen, Heartbeat, `register-consumer`/`acknowledge-consumer`, Retention-Löschausführung, Verarbeitung der Antrags-Queue `cdc.administration_request` — offene Anträge lesen und ihren Ausgang vermerken, ohne Anträge selbst anzulegen oder zu löschen —, Annahme eines Backfill-Antrags, die die Run-Zeile in `cdc.backfill_run` anlegt; `cdc.backfill_table` und die übrigen Antragsfunktionen ruft nur diese Rolle auf) | `CDC_ADMIN_DSN` |
| `cdc_reader` | Nur-Lese-Zugriff auf die Diagnose- und Lese-Views (`cdc.active_tables`, `cdc.consumer_status`, `cdc.changes`, `cdc.metrics`, `cdc.heartbeat`, `cdc.retention_blockers`, `cdc.backfill_status`) — trägt auch `--healthcheck` und `diagnose` (siehe [Diagnose ausführen](#diagnose-ausführen)) | `CDC_READER_DSN` |

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

**Betriebs-Hinweis (Backfill):** Die Software vergibt kein Recht auf Ihre
Quelltabellen. Soll ein Backfill eine Tabelle überführen, braucht die
Login-Identität hinter `CDC_CAPTURE_DSN` zusätzlich `SELECT` auf sie —
`GRANT SELECT ON <schema>.<tabelle> TO feed_capture_login;`. Fehlt es, endet
der Run `failed` mit der Fehlerklasse `permission`.

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

### Spalte vom Ausschluss konfigurieren

**Voraussetzung:** eine Login-Identität mit `cdc_admin`-Mitgliedschaft,
verbunden über `CDC_ADMIN_DSN` (siehe [Zugriff und Rollen](#zugriff-und-rollen));
die physische Tabelle existiert an der Quelle.

**Vorgehen:** Derselbe Antrags-Weg wie bei der Live-Aktivierung, adressiert
eine einzelne Spalte einer Tabelle:

```sql
SELECT cdc.exclude_column('<source_id>', '<schema>', '<tabelle>', '<spalte>');
```

**Ergebnis:** Der Aufruf ist asynchron: Er schreibt einen Antrag nach
`cdc.administration_request` (Status `pending`) und gibt dessen Kennung
zurück — die Spalte ist damit noch nicht ausgeschlossen. Die
Administrations-Goroutine des laufenden Feed-Containers verarbeitet offene
Anträge (`LISTEN`/`NOTIFY`-Weckung mit periodischem Fallback-Poll) und wirkt
bei Erfolg auf den Ausschlusszustand der laufenden Erfassung. Den Fortschritt
prüfen Sie über den Antrags-Status:

```sql
SELECT status, error_message
FROM cdc.administration_request
WHERE administration_request_id = '<zurückgegebene-id>';
```

`status` wechselt von `pending` zu `applied` (Erfolg) oder `failed`
(Fehlertext in `error_message`). Nach `applied` trägt jede danach erfasste
Änderung dieser Tabelle den ausgeschlossenen Spaltenschlüssel nicht mehr im
Row Image und den Wert nirgends; die nicht ausgeschlossenen Spalten bleiben
enthalten. Ein Ausschluss gegen eine an der Quelle nicht existierende Spalte
endet `failed` mit einem Fehlertext, der die Adresse `schema.tabelle.spalte`
nennt.

**Den Ausschluss wieder aufheben:** spiegelbildlich über
`cdc.include_column('<source_id>', '<schema>', '<tabelle>', '<spalte>')` —
derselbe Antrags-Weg und Status-Poll; nach `applied` wird die Spalte wieder
erfasst.

**Dauerhaftigkeit:** Der Ausschlussstand ist **dauerhaft** und hängt nicht an
der Prozesslebensdauer. Die `applied`-Zeilen der beiden Spalten-Antragsarten
sind die einzige Herkunft des Standes einer Tabelle (`SPEC-019`); jeder Pfad,
der eine Erfassungs-Bindung anlegt — der Prozessstart **und** die laufende
Aktivierung — trägt den abgeleiteten Stand mit. Ein Neustart des
Feed-Containers verliert den Ausschluss deshalb **nicht**, und eine
Deaktivierung mit anschließender Aktivierung stellt ihn ebenso wieder her.
Die Antrags-Zeilen dieser beiden Arten sind damit tragend: Eine Bereinigung
von `cdc.administration_request` verlöre den Stand.

**Betriebs-Hinweis:** Ein Ausschluss gegen eine existierende Spalte einer
Tabelle **ohne laufende Erfassung** endet trotzdem `applied`, nicht `failed`
— er wirkt, sobald die Tabelle erfasst wird. Anders als die beiden
Tabellen-Antragsarten (`enable`/`disable`), die Bindungs- und
Publication-Zeilen tragen, wirkt diese Antrags-Art allein auf den
Filterzustand der laufenden Erfassung.

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
       operation, old_data, new_data, committed_at, origin
FROM cdc.changes
WHERE source_id = '<quelle-id>' AND commit_position > <letzte-gelesene-position>
ORDER BY commit_position, sequence
LIMIT 500;
```

**Ergebnis:** Jede Zeile ist eine committed Änderung in Anhang-Reihenfolge.
`old_data`/`new_data` sind `jsonb`; bei `INSERT` ist `old_data` NULL, bei
`DELETE` ist `new_data` NULL. `origin` nennt die Herkunft der Änderung:
`wal` für eine über den Replication Stream erfasste Änderung, `backfill`
für eine Bestands-Änderung eines Backfills (`LH-FA-CAP-009`); eine
Änderung, die ohne dieses Feld gespeichert wurde, liest als `wal`
(`LH-FA-DAT-006`). `origin` ist die letzte Spalte der View; die drei
Live-Zustellwege (gRPC, SSE, NATS-Vollinhalt) tragen das Feld nicht.

Dieselben Änderungen sind ohne SQL-Direktzugriff über die API lesbar:
`GET /changes` — derselbe Lesezugriff mit denselben Filtern und derselben
Reihenfolge (siehe
[Zugriff über die HTTP-/JSON-API](#zugriff-über-die-http-json-api)).

**Fortsetzen und `LIMIT`:** Ein `LIMIT` schneidet Zeilen, nicht Positionen.
Trug eine Commit-Position mehr Änderungen als das `LIMIT`, überspringt
`commit_position > <letzte-gelesene-position>` den Rest dieser Position
(`SPEC-022`, Zeile „Position und `limit`"). Das trifft den Bestandsabzug eines
Backfills, der alle seine Änderungen auf **eine** Position legt (siehe
[Bestand als Backfill überführen](#bestand-als-backfill-überführen)): lesen Sie
ihn ohne `LIMIT` oder setzen Sie über den Schlüsselvergleich fort:

```sql
SELECT change_id, commit_position, transaction_id, sequence, operation, new_data, origin
FROM cdc.changes
WHERE source_id = '<quelle-id>'
  AND (commit_position, transaction_id, sequence) > (<letzte-position>, '<letzte-transaction-id>', <letzte-sequence>)
ORDER BY commit_position, transaction_id, sequence
LIMIT 500;
```

### Bestand als Backfill überführen

Die Erfassung trägt nur Änderungen ab der Aktivierung. Ein **Backfill**
überführt zusätzlich die Zeilen, die die Tabelle zum Startzeitpunkt des Runs
bereits enthält, als `INSERT`-Änderungen mit `origin = 'backfill'` in den Feed
(`LH-FA-CAP-009`, `LH-FA-ADM-001`). Er wird **ausdrücklich** ausgelöst — die
Aktivierung einer Tabelle startet keinen.

**Voraussetzung:** Die Tabelle ist aktiviert (laufende Bindung, Mitglied der
Publication, siehe [Tabelle live aktivieren](#tabelle-live-aktivieren)); eine
Login-Identität mit `cdc_admin`-Mitgliedschaft für den Antrag und dessen
Vermerk; eine Login-Identität mit `cdc_reader`-Mitgliedschaft (die hinter
`CDC_READER_DSN`) für das Lesen des Runs; der Feed-Container läuft für die
Quelle. Der Snapshot entsteht mit dem Start des Runs, nicht mit dem Antrag: ein
Run, der hinter einem anderen wartet, liest den Bestand zu seinem späteren
Start. Für den Lauf selbst gelten vier Betriebs-Vorbedingungen an der Quelle:

- **`SELECT`-Recht:** die Login-Identität hinter `CDC_CAPTURE_DSN` liest die
  Tabelle (siehe [Zugriff und Rollen](#zugriff-und-rollen)).
- **Reserve für die Replication-Verbindung:** ein Run legt für die Dauer der
  Slot-Anlage einen temporären Replication Slot an und öffnet dafür eine
  Replication-Verbindung — beides zählt gegen `max_replication_slots` und
  `max_wal_senders` der Quell-Instanz, zusätzlich zu den beiden Verbindungen
  des Feed-Containers (siehe [Grenzwerte](#grenzwerte)).
- **Snapshot-Haltedauer:** die Kopie liest den Bestand in **einem** Snapshot der
  Quelle und schreibt ihn in **einer** Transaktion des CDC-Speichers, die einmal
  am Ende committet. Für die Dauer der Kopie hält die Quelle den Snapshot (er
  hindert die Bereinigung von Zeilenversionen, die er noch sieht), und der
  CDC-Speicher trägt eine offene Schreibtransaktion. Je größer die Tabelle,
  desto länger. Eine Ablehnung großer Tabellen gibt es nicht; ab welcher Größe
  ein Run warnt, steht unter [Grenzwerte](#grenzwerte).
- **WAL-Rückstand des Capture-Slots:** die Kopie schreibt den Bestand in den
  CDC-Speicher derselben Datenbank. Dieses WAL trägt keine Änderung einer
  veröffentlichten Tabelle; der Feed bestätigt es im Leerlauf seines Streams
  (siehe [WAL-Rückstand prüfen](#wal-rückstand-prüfen)), es hält den
  Rückstand des Capture-Slots also nicht, und der Run beendet den Feed-Container
  nicht über die Fehlerschwelle — dasselbe gilt für jeden Schreiber auf eine
  nicht aktivierte Tabelle. Was bleibt, ist die offene
  Schreibtransaktion des Runs: sie hält das WAL für den Slot auf der Platte der
  Quelle, und der Walsender der Quelle dekodiert sie in seinen Arbeitsspeicher
  und lagert sie ab einer Größe auf die Platte aus. Beides wächst mit der Größe
  der Tabelle und ist eine Eigenschaft der Ein-Transaktions-Form des Runs; die
  Bestätigung im Leerlauf ändert es nicht. Gemessene Werte stehen unter
  [Grenzwerte](#grenzwerte). Wächst der Rückstand dennoch über die
  Fehlerschwelle (etwa weil der Feed nicht antwortet), beendet sich der
  Feed-Container mit der Klasse `replication` (Ausgang 1) und ein laufender Run
  endet `interrupted`; das Datei-Feld `wal_retention_error_bytes` (Bytes, kein
  Umgebungsvariablen-Gegenstück, `SPEC-016`, siehe
  [Optionale YAML-Konfigurationsdatei](#optionale-yaml-konfigurationsdatei-cdc_config_file))
  hebt die Fehlerschwelle über den erwarteten Rückstand — es ändert die Ursache
  nicht und wirkt für jede Ursache eines Rückstands.

**Speicher des Feed-Containers:** Ein Backfill lässt die Zahl der Changes in
`cdc.change` um die Zeilenzahl der Tabelle wachsen. Der Bereinigungslauf des Feeds
liest sie in jedem Takt seitenweise (10.000 Changes je Seite, ohne die Row Images,
siehe [Aufbewahrung (Retention)](#aufbewahrung-retention)); der Speicherbedarf des
Feed-Containers hängt deshalb nicht an dieser Zahl. Gemessen liegt die Spitze des
Feed-Containers nach Backfills über je 1.000.000 Zeilen bei 3.000.000 Changes in
`cdc.change` bei 17,3 bis 17,6 MiB (zwei Läufe; Messwerte und Herkunft unter
[Grenzwerte](#grenzwerte)).

**Sperre der Tabelle:** Vom Beginn der Lese-Transaktion bis zu ihrem Ende hält
der Run eine Lesesperre (`ACCESS SHARE`) auf die Tabelle. Lesen und Schreiben
der Tabelle laufen weiter, solange keine DDL auf die Tabelle wartet. Eine DDL,
die `ACCESS EXCLUSIVE` verlangt (`ALTER TABLE` mit Umschreiben der Tabelle,
`TRUNCATE`, `VACUUM FULL`, `CLUSTER`, `DROP TABLE`), wartet bis zum Ende des
Runs, und **jeder weitere Zugriff auf die Tabelle stellt sich hinter sie**:
Schreiber und Leser der Tabelle ebenso wie die Abfrage der
Publication-Mitgliedschaft (`pg_publication_tables`), die ein Antrag und der
Start eines Runs ausführen. Gemessen (PostgreSQL 18, Ursprung:
[Review-Report](../reviews/review-slice-backfill-e2e.md) F-1): mit einer
wartenden `ALTER TABLE … ALTER COLUMN … TYPE` liefen ein `INSERT` in die
Tabelle, ein `SELECT count(*)` auf die Tabelle und die Abfrage von
`pg_publication_tables` je in ein Zeitlimit von 4 s. Führen Sie eine solche DDL
an einer Tabelle nicht aus, solange ein Run für sie `running` ist (Status in
`cdc.backfill_status`); die Dauer des Runs wächst mit der Größe der Tabelle.

In der Gegenrichtung wartet der Run, solange eine fremde Transaktion
`ACCESS EXCLUSIVE` auf der Tabelle hält: er steht `running` mit `rows_copied` 0
und hält für diese Zeit seinen Snapshot. Der Run trägt dafür keine eigene
Zeitgrenze; das Warten endet mit dem Ende der fremden Transaktion oder mit dem
Abbruch des Runs (Ablauf seines Kontexts; der Adapter beendet ihn mit der Klasse
`transient` und hinterlässt keine Sitzung, gemessen im Store-Tier von
`make test-replication`, PostgreSQL 18). Beenden Sie eine offene DDL-Transaktion
auf der Tabelle, statt den Run warten zu lassen.

**Vorgehen:**

```sql
SELECT cdc.backfill_table('<source_id>', '<schema>', '<tabelle>');
```

**Ergebnis:** Der Aufruf ist asynchron: Er schreibt einen Antrag nach
`cdc.administration_request` (Art `backfill`, Status `pending`) und gibt dessen
Kennung zurück — die Kopie hat damit noch nicht begonnen. Der Feed-Container
nimmt den Antrag an: die Zeile des Runs entsteht mit dem Status `queued`, und
der Antrag wird `applied` vermerkt. **`applied` heißt hier „angenommen", nicht
„Bestand kopiert"** — den Verlauf des Runs zeigt allein `cdc.backfill_status`.
Ein Antrag endet `failed` (Fehlertext in `error_message`, keine Run-Zeile), wenn
die Tabelle nicht aktiviert ist, keine laufende Bindung hat oder ein Run
derselben Tabelle `queued` oder `running` ist. Die Abfrage läuft unter der
`cdc_admin`-Identität des Antrags:

```sql
SELECT status, error_message
FROM cdc.administration_request
WHERE administration_request_id = '<zurückgegebene-id>';
```

Den Run lesen Sie über die View `cdc.backfill_status` — je Tabelle der zuletzt
beantragte Run. Die Abfrage läuft unter einer `cdc_reader`-Identität; die Rolle
`cdc_admin` trägt kein `SELECT` auf die View:

```sql
SELECT status, rows_copied, estimated_rows, started_at, finished_at,
       error_message, warn_estimated_size, warn_duration
FROM cdc.backfill_status
WHERE source_id = '<source_id>' AND schema_name = '<schema>' AND table_name = '<tabelle>';
```

| `status` | Bedeutung |
|---|---|
| `queued` | angenommen, wartet; ein Run läuft zugleich, weitere warten in der Reihenfolge ihrer Anträge |
| `running` | die Kopie läuft; `rows_copied` schreitet je Block fort |
| `completed` | alle Zeilen sind in **einer** Transaktion geschrieben; `rows_copied` ist die Zahl der Backfill-Änderungen (eine leere Tabelle endet `completed` mit 0) |
| `failed` | der Run endete mit einem Fehler; `error_message` trägt die Fehlerklasse (siehe [Fehlerklassen](#fehlerklassen)) vor der Ursache; der Run hinterlässt keine Änderung |
| `interrupted` | der Prozess endete während der Kopie; der Run hinterlässt keine Änderung, `rows_copied` nennt den zuletzt festgehaltenen Fortschritt der Kopie und zählt keine sichtbaren Änderungen |

- **`estimated_rows` ist eine Schätzung** der Quelle (aus dem Katalog), keine
  Zählung und keine Grenze. NULL heißt **unbekannt** — der Katalog führt keine
  Schätzung —, nie 0. Die Ausgabe von `diagnose` zeigt es als „unbekannt".
- **`warn_estimated_size` und `warn_duration`** sind eine Kennzeichnung, die die
  Zeilenzahl und die Kopierdauer betrifft; sie ändern weder `status` noch den
  Ablauf, und kein Antrag wird wegen der Größe abgelehnt (siehe
  [Grenzwerte](#grenzwerte)).
  - `warn_estimated_size` ist `true`, wenn die **geschätzte** Zeilenzahl beim
    Antrag über der Richtgröße liegt. Eine unbekannte Schätzung (NULL) setzt sie
    nie: die Warnung schweigt gerade bei einer frisch befüllten, noch nicht
    analysierten Tabelle, für die der Katalog keine Schätzung führt (Messung
    siehe [Grenzwerte](#grenzwerte)).
  - `warn_duration` ist `true`, wenn die Kopierdauer (`finished_at −
    started_at`; die Wartezeit in `queued` zählt nicht) die Toleranz
    überschritten hat. Geprüft wird bei jedem Fortschritts-Update (je Block) und
    beim Ende des Runs, auch bei `failed` und `interrupted`; ein einzelner
    Block, der länger als die Toleranz braucht, warnt erst danach — die Warnung
    ist eine Orientierung, kein Alarm. Eine gesetzte Warnung bleibt gesetzt.
- **Umschreiben der Tabelle im Fenster:** Zwischen dem Snapshot-Export und der
  Sperre kann eine fremde DDL die Tabelle umschreiben (`ALTER TABLE … ALTER
  COLUMN … TYPE`, das die Datei neu schreibt, oder `TRUNCATE`); der ältere
  Snapshot sieht die neue Datei leer. Der Run erkennt das an der Datei der
  Tabelle und endet `failed` mit `error_message` `transient: …`, Ursache
  „nach dem Snapshot-Export umgeschrieben“ und ohne Änderung; ein **neuer
  Antrag** beginnt neu und liest den Bestand des neuen Zustands. Dieselbe
  Erkennung schlägt auch bei `VACUUM FULL` und `CLUSTER` an, obwohl der Snapshot
  die Zeilen noch sieht — ein Fehlalarm mit derselben Abhilfe. Das Fenster ist
  die Zeit zwischen Export und Sperre; seine Dauer ist nicht gemessen. Ein
  `DROP COLUMN` im Fenster endet ebenfalls `failed`, mit der Klasse `storage`
  (Phase DDL-Fenster von `make test-integration`); ein `RENAME COLUMN` endet
  gleich (gemessen im [Review-Report](../reviews/review-slice-backfill-e2e.md),
  PostgreSQL 18, Klasse `storage`; kein Beleg im E2E-Runner).
- Ein Run-Fehler ist **run-lokal**: er setzt weder den Fehlerzustand des
  Lebenszeichens noch stoppt er die Erfassung.

**Neustart und Wiederholung:** Beim Start des Feed-Containers wird jeder
`running`-Run der Quelle `interrupted`; ein `queued`-Run überlebt den Start und
wird ausgeführt. Ein `interrupted`-Run startet **nicht** von selbst neu: ein
erneuter Antrag (`cdc.backfill_table`) beginnt einen neuen Run mit neuem
Snapshot und damit von vorn. Ein abgebrochener Run hinterlässt keine
Änderungszeile. Betreiben Sie je Quelle **eine** Instanz: eine zweite Instanz
gegen dieselbe Quelle würde die `running`-Runs der ersten beim Start als
`interrupted` vermerken. Eine Instanz nimmt nur Backfill-Anträge ihrer eigenen
Quelle an; ein Antrag für eine andere Quelle bleibt `pending`, bis die Instanz
dieser Quelle ihn annimmt.

**Was der Bestand im Feed bedeutet:**

- **Position:** alle Änderungen eines Runs liegen auf **einer** Commit-Position
  `X` (`snapshot_position` des Runs). Positionen wandern nur vorwärts: ein
  Consumer, dessen bestätigte Position beim Commit des Runs `X` bereits erreicht
  hat, sieht den Bestand nicht über seinen Fortschritt; über das Bereichslesen
  (siehe [Änderungen lesen](#änderungen-lesen)) bleibt er lesbar.
- **Startposition eines neuen Consumers:** Ein frisch registrierter Consumer
  (`register-consumer`, `POST /consumers`) trägt keine bestätigte Position.
  `GET /consumers/position` liest ihn mit `offset` 0 und `acknowledged`
  `false`, `cdc.consumer_status` mit leerer bestätigter Position. Diese
  Anfangsposition liegt vor der Snapshot-Position `X` jedes Runs: liest der
  Consumer ab ihr (`commit_position > 0`, über `GET /changes` ohne `from`),
  gehört der Bestand zu seinem Fortschritt, gefolgt von den Änderungen der
  Erfassung. Bestätigt er eine Position hinter `X`, liegt der Bestand vor seiner
  Position und erscheint nicht in seinem Fortschritt. *Ursprung:* gemessen im
  Lauf von `make test-integration`, Phase „Backfill-Startposition“ (Zeile in
  [`e2e-abdeckung.md`](e2e-abdeckung.md)): `offset` 0, `acknowledged` `false`,
  5 Backfill-Änderungen mit `commit_position > 0`, 0 mit `commit_position`
  hinter der bestätigten Position; die Lauf-Zeile lautet „… startet an Position
  0 (acknowledged=false …), die Snapshot-Position des Bestands ist …: ab der
  Anfangsposition sind 5 Backfill-Changes lesbar; nach der Bestätigung von …
  liegen 0 hinter ihr“.
- **Lesen:** Ein Bestandsabzug teilt **eine** Commit-Position; ein `LIMIT`
  kann innerhalb einer Position nicht fortsetzen. Lesen Sie ihn ohne `LIMIT`
  oder über den Schlüsselvergleich (siehe [Änderungen lesen](#änderungen-lesen)).
- **Überlappung:** Eine Zeile, die im Fenster zwischen Aktivierung und `X`
  geändert wurde, kann doppelt erscheinen — die Änderung aus der Erfassung und
  der Bestand danach. Wendet ein Consumer das Log ab dem Anfang in Lese-Ordnung
  an (`INSERT` und `UPDATE` als Upsert des Row Images, `DELETE` als Löschen),
  entspricht sein Stand je Schlüssel dem Quellstand.
- **`schema_version`:** Die Schema-Version einer Backfill-Änderung ist die zum
  Run-Start aktuelle Version der Tabelle. Sie unterscheidet Backfill-Änderungen
  von Änderungen einer später registrierten Version; sie beschreibt die Spalten
  des Bildes nicht — die Spalten stehen im Bild selbst (`new_data`). Das Bild
  kann eine Spalte tragen, die die referenzierte Version nicht führt: bei einer
  Erweiterung während des Runs, bei einer kompatiblen Erweiterung vor dem Run,
  wenn seit der letzten Relation-Nachricht keine Änderung der Tabelle erfasst
  wurde, und bei einer Tabelle, die seit ihrer Aktivierung keine Änderung
  erfasst hat (die Version hat dann noch keine Spaltenform). Der Run ändert die
  Version nicht und bricht deshalb nicht ab.
- **Ausgeschlossene Spalten** (siehe [Spalte vom Ausschluss
  konfigurieren](#spalte-vom-ausschluss-konfigurieren)) tragen die Backfill-
  Änderungen nicht; ein Ausschluss, der während des Runs geändert wird, endet den
  Run `failed` (Klasse `configuration`), bevor etwas sichtbar wird.
- **Zustellung:** Backfill-Änderungen gehen in keinen der Live-Zustellwege (gRPC,
  SSE, NATS-Vollinhalt); nach dem Commit sendet der Run je Tabelle ein
  Wecksignal (NATS-Wecksignal, wenn aktiviert). Sie unterliegen der Aufbewahrung
  wie jede Änderung.

**Diagnose:** `diagnose` gibt den Run je Tabelle mit Status, Fortschritt,
**geschätzter** Zeilenzahl und den beiden Kennzeichnungen aus (siehe [Diagnose
ausführen](#diagnose-ausführen)); ein `failed`- oder `interrupted`-Run ist
Berichtsinhalt, kein Befehlsfehler.

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
existiert dafür nicht. Jeder Durchlauf liest die Changes der Quelle seitenweise
(10.000 je Seite) und dabei je Change nur Kennung, Commit-Position und
Commit-Zeitpunkt, keine Row Images; er entscheidet je Change über die Löschung
und löscht die freigegebenen Changes einer Seite, bevor er die nächste liest.
Der Speicher des Feed-Containers hängt deshalb an der Seitengröße, nicht an der
Zahl der Changes in `cdc.change` (Messwerte unter [Grenzwerte](#grenzwerte)); die
Arbeit der Datenbank je Durchlauf wächst dagegen mit dieser Zahl. Ein Durchlauf
ist nicht atomar: bricht er ab, bleiben die Löschungen der bereits bearbeiteten
Seiten bestehen, und der nächste Takt setzt fort. Die bestätigten
Consumer-Positionen liest ein Durchlauf einmal zu Beginn.

**Betriebs-Hinweis (Consumer-Bindung):** Ein registrierter, aber gegen
eine Quelle noch nie bestätigender Consumer blockiert die Bereinigung
dieser Quelle **nicht** — er zählt erst ab seiner ersten Bestätigung
(`acknowledge-consumer` bzw. [Position bestätigen](#position-bestätigen))
als schützenswert. Ein neu angebundener Consumer, der vor seinem ersten
Lesezugriff vor Bereinigung geschützt sein soll, sollte deshalb einmal
bestätigen (auch die dokumentierte Anfangsposition genügt), bevor er mit
dem Lesen beginnt. Bestätigt ein Consumer erstmals, während ein Durchlauf
läuft, gilt seine Position erst im nächsten Durchlauf.

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
`cdc_consumer_lag` je registriertem Consumer, `cdc_changes_pending`
(`LH-QA-OPS-003`, Label = Consumer-ID, Wert = Anzahl noch nicht
bestätigter Changes dieses Consumers), `cdc_errors_total` (`LH-QA-OPS-003`,
Label = Fehlerklasse aus `cdc.process_heartbeat.error_class`, Wert =
Anzahl der Quellen aktuell in dieser Fehlerklasse), sowie
`cdc_storage_bytes` (`LH-FA-RET-006`, physische Speichergröße von
`cdc.change` über `pg_relation_size` — der mit dem Erfassungsvolumen
wachsenden Tabelle).

`cdc_wal_retention_bytes` (WAL-Rückstand des Capture-Slots, `SPEC-009`)
steht **nicht** in `cdc.metrics`: Die Erhebung braucht Systemkatalog-Zugriffe
außerhalb des `cdc`-Schemas (`pg_replication_slots`, `IDENTIFY_SYSTEM`), die
die Least-Privilege-Fläche von `cdc_reader` unnötig erweitern würden — siehe
[WAL-Rückstand prüfen](#wal-rückstand-prüfen).

### Diagnose ausführen

Statt der SQL-Abfragen oben einzeln zu stellen, liest der
`diagnose`-Sondermodus dieselben Views (`cdc.heartbeat`, `cdc.metrics`,
`cdc.retention_blockers`, `cdc.backfill_status`) über `CDC_READER_DSN` und gibt eine
menschenlesbare Zusammenfassung aus (`LH-FA-SST-003`, deckt
`LH-FA-ADM-002`…`005`, `LH-FA-RET-005`/`006`, `LH-FA-CAP-009`) — derselbe Image-Tag wie der
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
pg-change-feed diagnose: Quelle "src-e2e"
  Betriebsstatus (LH-FA-ADM-002): Lebenszeichen vor 1.203s
  Fehlerzustand (LH-FA-ADM-003): keiner (Normalbetrieb)
  CDC-Abstand cdc_capture_lag (LH-FA-ADM-004): 0.087s
  Verarbeitungsrückstand cdc_consumer_lag je Consumer (LH-FA-ADM-005, nur Consumer mit mindestens einer bestätigten Position):
    cli-e2e-consumer: 3
  Blockierender Consumer (LH-FA-RET-005): cli-e2e-consumer, bestätigte Position 42, Rückstand 3
  Speicherverbrauch cdc_storage_bytes (LH-FA-RET-006): 65536 Bytes
  Backfill je Tabelle (LH-FA-CAP-009, letzter Run; die Zeilenzahl ist geschätzt):
    public.orders: completed, 1200 Zeilen kopiert, geschätzt 1150, Warnung Größe false, Warnung Dauer false
    public.audit: failed, 0 Zeilen kopiert, geschätzt unbekannt, Warnung Größe false, Warnung Dauer false
      Fehler: permission: …
    public.events: completed, 4800000 Zeilen kopiert, geschätzt 4700000, Warnung Größe true, Warnung Dauer true
```

Der Abschnitt „Backfill je Tabelle" liest die View `cdc.backfill_status` (siehe
[Bestand als Backfill überführen](#bestand-als-backfill-überführen)): je
Tabelle der zuletzt beantragte Run mit Status, Fortschritt, der **geschätzten**
Zeilenzahl (eine unbekannte Schätzung erscheint als „unbekannt“, nie als 0) und
den beiden Kennzeichnungen (`Warnung Größe`: die geschätzte Zeilenzahl lag beim
Antrag über der Richtgröße; `Warnung Dauer`: die Kopierdauer hat die Toleranz
überschritten, siehe [Grenzwerte](#grenzwerte)); bei einem `failed`-Run folgt
der Fehlertext. Ist
noch nie ein Backfill beantragt worden, steht dort „(keiner — kein Backfill
beantragt)".

Trägt keine Quelle in `cdc.retention_blockers` gar keine Zeile (noch kein
Consumer hat je gegen sie bestätigt), zeigt die Zeile stattdessen:

```text
  Blockierender Consumer (LH-FA-RET-005): kein Blocker (kein Consumer hat je gegen diese Quelle bestätigt)
```

**Ergebnis:** Wie bei den Rohwerten der Views trifft der Befehl keine
Schwellenwert-Entscheidung (`SPEC-007` bleibt Sache des lesenden Systems) und
der Prozess-Ausgang trägt nur den Lese-Erfolg — ein gemeldeter Fehlerzustand
oder Rückstand ist Berichtsinhalt, kein Befehlsfehler (Ausgang bleibt 0). Ein
Consumer ohne je bestätigte Position erscheint nicht in der Rückstands-Liste
(dieselbe Grenze wie bei `cdc_consumer_lag` in [Metriken
lesen](#metriken-lesen)); ein Consumer mit bestätigter Position, dessen
Quelle noch nie eine Transaktion trug, erscheint mit dem Text „unbekannt"
statt einem irreführenden Rückstand von 0. Der blockierende Consumer und
`cdc_storage_bytes` folgen derselben Lese-Disziplin wie [Blockierende
Consumer erkennen](#blockierende-consumer-erkennen) und [Metriken
lesen](#metriken-lesen) — keine neue Berechnung, nur dieselben Sichten über
die CLI ausgegeben.

### WAL-Rückstand prüfen

Der Feed-Container misst den WAL-Rückstand des Capture-Slots periodisch
(gebunden an denselben Takt wie das Lebenszeichen, siehe
[Grenzwerte](#grenzwerte)) über die Rolle `cdc_capture` und protokolliert ihn
strukturiert:

```json
{"msg": "replication: WAL-Rückstand gemessen", "metric": "cdc_wal_retention_bytes", "bytes": 12345}
```

**Ergebnis:** `bytes` ist das WAL, das die Quelle dem Feed geliefert, der Feed
aber noch nicht bestätigt hat (aktuelles WAL-Ende der Instanz minus
`confirmed_flush_lsn` des Slots). WAL ohne Inhalt für die Publication — ein
Schreiber auf eine nicht aktivierte Tabelle, ein Backfill-Run — bestätigt der
Feed im Leerlauf seines Streams selbst und lässt den Wert damit nicht wachsen. `bytes` wächst, solange der Feed nicht
bestätigt: der Capture-Slot ist inaktiv und die Quelle schreibt weiter (z. B.
während eines Verbindungsabbruchs), oder der Feed antwortet nicht — ein
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

**Reihenfolge beim Upgrade.** Rollen Sie das Schema **vor** dem Tausch des
Feed-Containers aus, dann ersetzen Sie den Container durch die neue Version.
Der neue Container erwartet das neue Schema (er schreibt zum Beispiel die
Spalte `cdc.change.origin`); umgekehrt ist der Rollout vor dem Tausch für
den weiterlaufenden Container erwartungsgemäß unkritisch, weil neue Spalten
nullable und neue Tabellen oder Views additiv sind (abgeleitet, nicht mit
einem laufenden Alt-Container gemessen).

**Additive Änderungen** — neue Tabelle, neue nullable Spalte, neue View —
rollen ohne Zwischenschritt aus. Ein zweiter Lauf gegen ein bereits
ausgerolltes Ziel endet ebenfalls mit Exit 0.

**Rechte der drei Rollen.** Der Rollout setzt die Rechte der Gruppenrollen
bei jedem Lauf (idempotent), auch gegen ein bereits ausgerolltes Ziel. Die
Rechte von `cdc_admin` auf `cdc.administration_request` (`SELECT`, `UPDATE`)
und auf `cdc.backfill_run` (`SELECT`, `INSERT`), von `cdc_capture` auf
`cdc.backfill_run` (`SELECT`, `UPDATE`) und von `cdc_reader` auf die View
`cdc.backfill_status` (`SELECT`) kommen aus diesem Lauf: rollen
Sie das Schema nach einem Upgrade **vor** dem Tausch des Feed-Containers aus.
Ohne das Recht auf `cdc.administration_request` bleiben Anträge der
SQL-Administration `pending`, und der Feed-Container protokolliert „Anträge
lesen fehlgeschlagen" — das trifft eine Login-Identität, die nur `IN ROLE
cdc_admin` ist, nicht einen Superuser-Login. Die Antragsfunktion
`cdc.backfill_table` gehört wie die übrigen Antragsfunktionen `cdc_admin`
allein: eine Login-Identität ohne diese Mitgliedschaft scheitert mit
„permission denied for function".

**Änderung an der Spaltenliste einer View.** Ändert ein Release die Signatur
einer View des Schemas (Spalte anhängen, umordnen, umbenennen, Typ ändern —
so trägt `cdc.changes` seit der Einführung des Feldes `origin` eine
zusätzliche letzte Spalte), kann d-migrate die View nicht in-place ersetzen.
Der Lauf entfernt sie dann selbst und legt sie neu an; er meldet das:

```text
schema-rollout: Vorlauf (ADR-0114) - View-Signatur-Aenderung, DROP VIEW cdc.changes
```

Dabei gilt:

- Es gehen keine Daten verloren: die View trägt keine Daten, die erfassten
  Changes liegen in `cdc.change` und sind danach über `cdc.changes` unverändert
  lesbar. Die Rechte der Rolle `cdc_reader` setzt derselbe Lauf wieder.
- **Eigene Rechte:** Das Entfernen der View verwirft ihre gesamte
  Rechteliste. Nur `cdc_reader` bekommt das `SELECT`-Recht im selben Lauf
  zurück; ein von Ihnen an eine andere Rolle vergebenes `GRANT SELECT ON
  cdc.<view>` fehlt danach. Setzen Sie solche Rechte nach einem Rollout mit
  Vorlauf-Meldung erneut.
- **Vorbedingung:** Der Vorlauf adressiert das Schema `cdc` fest
  (`DROP VIEW cdc.<name>`); wie beim gesamten Rollout gilt der
  `search_path` `cdc` des Ziels.
- **Lesefenster:** Für SQL-Leser über `cdc_reader` fehlt die View — oder
  sie ist noch ohne Recht — für die Dauer des Rollouts. Richtwert rund 7
  Sekunden, aus einer einzelnen Architect-Messung auf einer Testinstanz mit
  einer Zeile (`ADR-0114`); eine Messung, nicht garantiert — auf einem
  größeren Ziel kann es länger dauern. Der Feed-Container liest den Store über
  die Tabellen, nicht über diese View (aus dem Quelltext abgeleitet, nicht
  während eines Rollouts gemessen).
- Der Schritt läuft nur in einem Lauf, der eine Signaturänderung ausliefert;
  ein Lauf gegen ein Ziel mit aktueller View-Signatur meldet keinen Vorlauf und
  hat kein Fenster.
- Hängt ein eigenes Objekt (etwa eine selbst angelegte View) an der
  betroffenen View, bricht der Lauf mit dem PostgreSQL-Fehler ab. Es wird nichts
  mitgelöscht: nehmen Sie das eigene Objekt vor dem Rollout weg und legen Sie es
  danach neu an.
- Bricht der Lauf nach dem Vorlauf ab, fehlt die View bis zum nächsten
  Lauf; ein Wiederholungslauf legt sie wieder an.

Ein Rollout, der am Precheck mit Exit 8 endet (ein nicht bekannter Blocker),
hat das Ziel nicht verändert — gemessen für eine ausstehende
View-Signaturänderung ohne Vorlauf: die neue Spalte war danach nicht angelegt.

### Zugriff über die HTTP-/JSON-API

Die API stellt die Verwaltungsfähigkeiten des Feed-Containers zusätzlich als
Netzwerkzugriffsweg bereit. Wo dieselbe Fähigkeit auch über CLI oder SQL
erreichbar ist, führt sie dieselbe Domänenlogik aus und ist fachlich
gleichwertig — kein Zweitpfad. Eine Ausnahme ist das Auslösen der
Aufbewahrung: Dafür gibt es keinen CLI-/SQL-Zugriffsweg — die API ist der
einzige manuelle Auslöser, und der automatische Hintergrundzug läuft
unabhängig davon (siehe [Aufbewahrung (Retention)](#aufbewahrung-retention)).

**Erreichbarkeit:** aktiv, sobald `CDC_HTTP_ADDR` gesetzt ist (`host:port`,
siehe [Konfiguration](#5-konfiguration)); ungesetzt bleibt sie vollständig
deaktiviert — kein Server, kein Port.

**Authentifizierung:** Jeder Aufruf trägt den Header
`Authorization: Bearer <token>`. Es gibt zwei Token-Klassen:
`CDC_API_TOKEN_READER` (lesend) und `CDC_API_TOKEN_ADMIN`
(administrativ/schreibend), wobei das Admin-Token die lesende Klasse implizit
mit abdeckt. Ein fehlender oder keinem konfigurierten Token entsprechender
Wert endet `401`, ein bekanntes Token mit unzureichender Klasse `403`. Die
Token-Klassen entscheiden an der API-Schicht, welche Fähigkeit erreichbar
ist; die Datenbankverbindung darunter trägt weiterhin die bei der Verdrahtung
fixierte Rolle. Fehlerantworten tragen die Form `{"error": "<Klartext>"}`.

| Fähigkeit | Methode und Pfad | Rechtsklasse |
|---|---|---|
| Consumer registrieren | `POST /consumers` | `admin` |
| Position bestätigen | `POST /consumers/acknowledge` | `admin` |
| Consumer-Position lesen | `GET /consumers/position` | `reader` |
| Consumer entfernen | `POST /consumers/remove` | `admin` |
| Tabelle aktivieren | `POST /tables/enable` | `admin` |
| Tabelle deaktivieren | `POST /tables/disable` | `admin` |
| Tabellen-Status | `GET /tables/status` | `reader` |
| Tabellen auflisten | `GET /tables` | `reader` |
| Änderungen lesen | `GET /changes` | `reader` |
| Aufbewahrung auslösen | `POST /retention/run` | `admin` |

Die lesenden Endpunkte sind mit dem Admin-Token ebenso erreichbar; mit dem
Reader-Token sind die administrativen Endpunkte nicht erreichbar (`403`).

**Changes lesen:** `GET /changes?source=<quelle-id>` liefert persistierte
Änderungen einer Quelle — optional gefiltert über `schema` und `table`
(je einzeln oder zusammen), eingegrenzt über `from` (inklusive) und `to`
(exklusive, jeweils ein `commit_position`-Wert) und begrenzt über `limit`
(≥ 1). Ohne `limit` liest der Aufruf unbegrenzt. Die Antwort trägt je
Änderung `commit_position`, `change_id`, `transaction_id`,
`source_table_id`, `schema`, `table`, `sequence`, `operation`,
`old_image`, `new_image`, `schema_version`, `committed_at` (RFC 3339,
UTC) und `origin` (`wal` oder `backfill`, als letztes Feld; eine ohne
dieses Feld gespeicherte Änderung liest als `wal`) — dieselbe Sicht wie
der SQL-Zugriff auf `cdc.changes`. Die
Reihenfolge ist deterministisch; die Fortsetzung liest ab
`from = <letzte gelieferte commit_position> + 1`, wenn das Lesen die letzte
Position vollständig erfasst hat. Ein `limit` schneidet Zeilen, nicht Positionen:
Enthält eine Commit-Position mehr Änderungen als `limit`, liefert dieses `from`
den Rest der Position nicht, und ein Lesen ab derselben Position liefert wieder
dieselben Zeilen (`SPEC-022`). Den Bestandsabzug eines Backfills, der alle seine
Änderungen auf **eine** Position legt (siehe [Bestand als Backfill
überführen](#bestand-als-backfill-überführen)), lesen Sie deshalb ohne `limit`
oder über den Schlüsselvergleich im SQL-Zugriff (siehe [Änderungen
lesen](#änderungen-lesen)). Ein Aufruf ohne Treffer
endet `200` mit leerer Liste (`{"changes": []}`), nie `404`; ein Parameter
außerhalb der genannten Liste endet `400`, ebenso ein fehlendes `source`,
eine nicht lesbare Zahl, `from`/`to` unter 1, `limit` unter 1 und
`from > to`.

**Zustellsemantik:** Diese Fähigkeiten sind synchrone Anfrage/Antwort — die
Antwort trägt das Ergebnis des Aufrufs, es gibt keine Warteschlange
dazwischen. Die Ausnahme ist der Live-Stream auf `GET /changes/stream` (siehe
unten), der die Verbindung offen hält — `GET /changes` ist demgegenüber die
nicht streamende Form desselben Gegenstands.

**Beispiele:** Jede Sprache ruft denselben `reader`-Endpunkt `GET /tables`
auf und gibt die Antwort aus; Adresse und Token liest jedes Beispiel aus
`CDC_HTTP_ADDR` und `CDC_API_TOKEN_READER` und lässt sich per Flag
übersteuern.

- **Go:** `examples/http-client` — Container-Aufruf über
  `make example-run-go SURFACE=http ARGS="-source <quelle> -publication <publication>"`
  (baut bei Bedarf `pg-change-feed-examples:go` aus `examples/Dockerfile`)
- **C#:** `examples/csharp/http-client` — Container-Aufruf über
  `make example-run-csharp SURFACE=http ARGS="--source <quelle> --publication <publication>"`
  (startet das mit `make examples-csharp` gebaute Image)
- **Kotlin:** `examples/kotlin/http-client` — Container-Aufruf über
  `make example-run-kotlin SURFACE=http ARGS="--source <quelle> --publication <publication>"`
  (startet das mit `make examples-kotlin` gebaute Image)

**SDK:** .NET-Anwendungen können statt der Beispiele das offizielle
NuGet-Package `PgChangeFeed.Client` einbinden (`LH-FA-SST-009`, `ADR-0106`,
`dotnet add package PgChangeFeed.Client`) — `PgChangeFeedHttpClient` deckt
alle zehn Fähigkeiten dieser Zugriffs-Oberfläche ab (die neun in der
Tabelle oben plus `GET /changes`) mit typisierten Requests/Responses und
einer typisierten Fehlerklasse für `400`/`401`/`403`/`404`/`500`, statt den
Draht-Vertrag selbst zu implementieren; siehe `sdks/csharp/README.md`.

Python-Anwendungen können statt der Beispiele das offizielle PyPI-Package
`pgchangefeed` einbinden (`LH-FA-SST-009`, `ADR-0107`,
`pip install pgchangefeed`) — `PgChangeFeedHttpClient` deckt dieselben zehn
Fähigkeiten dieser Zugriffs-Oberfläche ab (die neun in der Tabelle oben
plus `GET /changes`) mit typisierten Requests/Responses (`dataclasses`) und
einer typisierten Fehlerklasse für `400`/`401`/`403`/`404`/`500`; derselbe
Package trägt außerdem den gRPC-Change-Stream (siehe unten, „Zugriff über
den gRPC-Change-Stream"), den SSE-Stream (siehe „Zugriff über
Server-Sent-Events") und den NATS-Vollinhalts-Stream (siehe „Zugriff über
den NATS-Vollinhalts-Stream", `ADR-0110`). Siehe `sdks/python/README.md`.

Kotlin/JVM-Anwendungen können statt der Beispiele das offizielle
Gradle-/Maven-Package `pgchangefeed-kotlin` einbinden (`LH-FA-SST-009`,
`ADR-0109`, Koordinate `io.github.pt9912:pgchangefeed-kotlin`) —
`PgChangeFeedHttpClient` deckt dieselben zehn Fähigkeiten dieser
Zugriffs-Oberfläche ab (die neun in der Tabelle oben plus `GET /changes`)
mit typisierten Requests/Responses und einer versiegelten
(`sealed class`) Fehlerklasse für `400`/`401`/`403`/`404`/`500`; derselbe
Package trägt außerdem den gRPC-Change-Stream (siehe unten, „Zugriff über
den gRPC-Change-Stream"), den SSE-Stream (siehe „Zugriff über
Server-Sent-Events") und den NATS-Vollinhalts-Stream (siehe „Zugriff über
den NATS-Vollinhalts-Stream"). Anders als NuGet/PyPI wird dieses Package über
**GitHub Packages** vertrieben (`https://maven.pkg.github.com/pt9912/pg-change-feed`,
`ADR-0109` Festlegung 2) — GitHub Packages verlangt **immer** eine
Authentifizierung zum Lesen, auch für ein öffentliches Package: ein
GitHub-Konto und ein klassischer Personal Access Token (PAT) mit dem Scope
`read:packages` sind Voraussetzung für den Bezug, unabhängig davon, ob das
SDK öffentlich und quelloffen ist. Siehe `sdks/kotlin/pgchangefeed-kotlin/README.md`
für den vollständigen Installationsweg samt Gradle-Zugangsdaten-Konfiguration.

In allen drei Packages trägt jede über `GET /changes` gelesene Änderung das
Feld `origin` (`wal` oder `backfill`); eine Antwort ohne dieses Feld oder mit
JSON-`null` liest das Package als `wal`, jeden anderen Wert des Servers (auch
einen leeren oder unbekannten) gibt es unverändert weiter. Die Live-Wege
(gRPC, SSE, NATS-Vollinhalt) tragen kein `origin`.

### Zugriff über den gRPC-Change-Stream

**Erreichbarkeit:** aktiv, sobald `CDC_GRPC_ADDR` gesetzt ist (`host:port`);
ungesetzt bleibt der Streaming-Server vollständig deaktiviert — kein
Listener.

**Authentifizierung:** derselbe Token-Wert wie bei der HTTP-API, übertragen
als gRPC-Metadata-Eintrag `authorization` in der Form `Bearer <token>`.
Fehlt er oder entspricht er keiner konfigurierten Klasse, endet der Aufruf
mit dem gRPC-Status `Unauthenticated` — nicht mit einem stillen, leeren
Stream.

**Der Stream:** Der Server-Streaming-RPC `ChangeStream/StreamChanges`
(gRPC über HTTP/2 mit Protobuf) überträgt jedem verbundenen Consumer jeden
committed Change mit vollständigem Inhalt — eine Nachricht je Zeilen-Änderung.
Jede Nachricht trägt zehn Felder:

| Feld | Typ | Bedeutung |
|---|---|---|
| `change_id` | string | Kennung des Changes |
| `transaction_id` | string | Kennung der Quelltransaktion |
| `source_table_id` | string | Kennung der aktivierten Tabelle |
| `sequence` | int64 | Reihenfolge der Zeilen-Änderung in der Transaktion |
| `operation` | string | `INSERT`, `UPDATE` oder `DELETE` |
| `old_image` | bytes | Row Image vor der Änderung; bei `INSERT` leer |
| `new_image` | bytes | Row Image nach der Änderung; bei `DELETE` leer |
| `schema_version` | string | Schema-Version des Changes |
| `schema` | string | Schema-Name der Tabelle |
| `table` | string | Tabellenname |

Eine Filterung nach Tabelle ist nicht Teil dieser Version.

**Zustellsemantik:** Es gibt **keine** Zustellgarantie (Fire-and-Forget,
verlustbehaftet). Je Abonnent trägt der Server eine begrenzte
Empfangswarteschlange; ist sie voll, weil der Client langsamer liest als
Changes eintreffen, werden die betroffenen Nachrichten **verworfen**, statt
gepuffert zu werden. Ein nicht verbundener oder zu langsam lesender Consumer
verpasst sie damit ersatzlos — ein Replay innerhalb des Streams gibt es
nicht. Der Erzeuger hält nie auf einen Empfänger an: Ein langsamer oder
hängender Abonnent stoppt den Erfassungsbetrieb nicht.

**Nachvollziehbarkeit:** Verpasste Changes bleiben über den Lesezugriffsweg
[Änderungen lesen](#änderungen-lesen) und über die bestätigte
Consumer-Position ([Position bestätigen](#position-bestätigen)) nachholbar —
der Stream ersetzt diesen Zugriffsweg nicht.

**Beispiele:** Jede Sprache öffnet denselben Server-Streaming-RPC
`ChangeStream/StreamChanges` und gibt jede empfangene Nachricht aus; Adresse
und Token liest jedes Beispiel aus `CDC_GRPC_ADDR` und `CDC_API_TOKEN_READER`
und lässt sich per Flag übersteuern.

- **Go:** `examples/grpc-client` — Container-Aufruf über
  `make example-run-go SURFACE=grpc` (baut bei Bedarf
  `pg-change-feed-examples:go-grpc` aus `examples/Dockerfile`)
- **C#:** `examples/csharp/grpc-client` — Container-Aufruf über
  `make example-run-csharp SURFACE=grpc` (startet das mit
  `make examples-csharp` gebaute Image; der Stub entsteht im Bau aus der
  `.proto`, über einen zusätzlichen, benannten Bau-Kontext gelesen —
  `ADR-0090`)
- **Kotlin:** `examples/kotlin/grpc-client` — Container-Aufruf über
  `make example-run-kotlin SURFACE=grpc` (startet das mit
  `make examples-kotlin` gebaute Image; derselbe Stub-im-Bau-Mechanismus wie
  beim C#-Client, übertragen auf die Kotlin-Werkzeugkette — `ADR-0090`)

**SDK:** .NET-Anwendungen können statt der Beispiele das offizielle
NuGet-Package `PgChangeFeed.Client` einbinden (`LH-FA-SST-009`, `ADR-0106`,
`dotnet add package PgChangeFeed.Client`) — `PgChangeFeedGrpcClient.StreamChangesAsync`
öffnet den `ChangeStream/StreamChanges`-RPC und liefert ein
`IAsyncEnumerable<Change>` mit allen zehn Feldern der Tabelle oben; das
Bearer-Token landet im `authorization`-Metadata-Eintrag, ein fehlendes oder
ungültiges Token endet den Aufruf mit gRPC-Status `Unauthenticated`, statt
den Draht-Vertrag selbst zu implementieren; siehe `sdks/csharp/README.md`.

Kotlin/JVM-Anwendungen können statt des Beispiels dasselbe offizielle
Gradle-/Maven-Package `pgchangefeed-kotlin` einbinden (`LH-FA-SST-009`,
`ADR-0109`, Koordinate `io.github.pt9912:pgchangefeed-kotlin`) —
`PgChangeFeedGrpcClient.streamChanges()` öffnet denselben
`ChangeStream/StreamChanges`-RPC und liefert ein
`kotlinx.coroutines.flow.Flow<Change>` mit allen zehn Feldern der Tabelle
oben; das Bearer-Token landet im `authorization`-Metadata-Eintrag, ein
fehlendes oder ungültiges Token endet den Aufruf mit gRPC-Status
`Unauthenticated`, statt den Draht-Vertrag selbst zu implementieren. Derselbe
Package trägt außerdem den SSE-Stream (siehe „Zugriff über
Server-Sent-Events") und den NATS-Vollinhalts-Stream (siehe „Zugriff über
den NATS-Vollinhalts-Stream"). Wie beim HTTP-API-Zugriff oben verlangt der
Bezug über **GitHub Packages**
(`https://maven.pkg.github.com/pt9912/pg-change-feed`) immer eine
Authentifizierung, auch für dieses öffentliche Package: ein GitHub-Konto
und ein klassischer Personal Access Token (PAT) mit dem Scope
`read:packages` sind Voraussetzung für den Bezug (`ADR-0109`
Festlegung 2). Siehe `sdks/kotlin/pgchangefeed-kotlin/README.md`.

Python-Anwendungen können statt des Beispiels das offizielle
PyPI-Package `pgchangefeed` einbinden (`LH-FA-SST-009`, `ADR-0110`,
`pip install pgchangefeed`) — `PgChangeFeedGrpcClient.stream_changes()`
öffnet denselben `ChangeStream/StreamChanges`-RPC und liefert einen
Iterator über die generierten `Change`-Nachrichten mit allen zehn Feldern
der Tabelle oben; das Bearer-Token landet im `authorization`-Metadata-Eintrag,
ein fehlendes oder ungültiges Token endet den Aufruf mit gRPC-Status
`Unauthenticated`, statt den Draht-Vertrag selbst zu implementieren
(`timeout` ist der Gesamtfriestempel des Aufrufs in Sekunden, `None` =
unbegrenzt). Siehe `sdks/python/README.md`.

### Zugriff über Server-Sent-Events

**Erreichbarkeit:** derselbe HTTP-Server wie oben — aktiv, sobald
`CDC_HTTP_ADDR` gesetzt ist; der Endpunkt ist `GET /changes/stream`, seine
Rechtsklasse `reader` oder `admin`.

**Authentifizierung:** wie die übrigen Endpunkte der HTTP-API über
`Authorization: Bearer <token>`; ein fehlender oder unbekannter Token endet
`401`, bevor das erste Event läuft.

**Der Stream:** Die Antwort trägt `Content-Type: text/event-stream`; je
Change ein Event `event: change`, dessen `data:` ein JSON-Objekt mit
denselben zehn Feldern wie der gRPC-Stream trägt (siehe
[Zugriff über den gRPC-Change-Stream](#zugriff-über-den-grpc-change-stream));
ein fehlendes Row Image ist `null`. Jedes Event wird sofort ausgeliefert. Ist
die Adresse gesetzt, aber kein Live-Stream-Träger verdrahtet, antwortet der
Endpunkt mit `503`.

**Zustellsemantik:** keine Zustellgarantie (Fire-and-Forget): Ein nicht
verbundener oder langsamer lesender Client verpasst die betroffenen
Änderungen ersatzlos. Der `Last-Event-ID`-Header wird weder gesendet noch
ausgewertet — es gibt kein Stream-internes Replay. Der Erzeuger hält nie auf
einen Empfänger an.

**Nachvollziehbarkeit:** wie beim gRPC-Stream — verpasste Changes bleiben
über [Änderungen lesen](#änderungen-lesen) und die bestätigte
Consumer-Position ([Position bestätigen](#position-bestätigen)) nachholbar.

**Beispiele:** Jede Sprache öffnet `GET /changes/stream` und gibt jedes Event
aus; Adresse und Token liest jedes Beispiel aus `CDC_HTTP_ADDR` und
`CDC_API_TOKEN_READER` und lässt sich per Flag übersteuern.

- **Go:** `examples/sse-client` — Container-Aufruf über
  `make example-run-go SURFACE=sse` (baut bei Bedarf
  `pg-change-feed-examples:go-sse` aus `examples/Dockerfile`)
- **C#:** `examples/csharp/sse-client` — Container-Aufruf über
  `make example-run-csharp SURFACE=sse` (startet das mit
  `make examples-csharp` gebaute Image)
- **Kotlin:** `examples/kotlin/sse-client` — Container-Aufruf über
  `make example-run-kotlin SURFACE=sse` (startet das mit
  `make examples-kotlin` gebaute Image)

**SDK:** .NET-Anwendungen können statt des Beispiels dasselbe offizielle
NuGet-Package `PgChangeFeed.Client` einbinden (`LH-FA-SST-009`, `ADR-0106`,
`dotnet add package PgChangeFeed.Client`) — `PgChangeFeedSseClient.StreamChangesAsync`
öffnet `GET /changes/stream` und liefert ein `IAsyncEnumerable<Change>` mit
allen zehn Feldern der Tabelle oben; das Bearer-Token landet im
`Authorization`-Header, ein fehlender oder unbekannter Token endet den
Aufruf mit `PgChangeFeedUnauthorizedException` (HTTP-Status `401`), statt
den Draht-Vertrag selbst zu implementieren; siehe `sdks/csharp/README.md`.

Kotlin/JVM-Anwendungen können statt des Beispiels dasselbe offizielle
Gradle-/Maven-Package `pgchangefeed-kotlin` einbinden (`LH-FA-SST-009`,
`ADR-0109`, Koordinate `io.github.pt9912:pgchangefeed-kotlin`) —
`PgChangeFeedSseClient.streamChanges()` öffnet `GET /changes/stream` und
liefert eine `Sequence<Change>` mit allen zehn Feldern der Tabelle oben; das
Bearer-Token landet im `Authorization`-Header, ein fehlender oder
unbekannter Token endet den Aufruf mit `PgChangeFeedUnauthorizedException`
(HTTP-Status `401`), statt den Draht-Vertrag selbst zu implementieren. Wie
beim HTTP-API-Zugriff oben verlangt der Bezug über **GitHub Packages**
immer eine Authentifizierung, auch für dieses öffentliche Package
(`ADR-0109` Festlegung 2). Siehe `sdks/kotlin/pgchangefeed-kotlin/README.md`.

Python-Anwendungen können statt des Beispiels das offizielle
PyPI-Package `pgchangefeed` einbinden (`LH-FA-SST-009`, `ADR-0110`,
`pip install pgchangefeed`) — `PgChangeFeedSseClient.stream_changes()`
öffnet denselben Endpunkt `GET /changes/stream` und liefert einen Iterator
über die getypten `StreamChange`-Events mit allen zehn Feldern der Tabelle
oben; das Bearer-Token landet im `Authorization`-Header, ein fehlender oder
unbekannter Token endet den Aufruf mit `PgChangeFeedUnauthorizedError`
(HTTP-Status `401`), statt den Draht-Vertrag selbst zu implementieren.
Siehe `sdks/python/README.md`.

### Zugriff über das NATS-Wecksignal

Anders als die drei Abschnitte zuvor ist dieser Zugriffsweg **signal-tragend**,
nicht **daten-tragend**: man *bekommt* hier keine Änderung, man *erfährt*, dass
man nachsehen muss. Wer die Nachricht wie einen Datenstrom liest und Inhalt in
ihr erwartet, benutzt sie falsch.

**Erreichbarkeit:** optional über `CDC_NATS_URL`; ungesetzt bleibt das Feature
vollständig deaktiviert — keine NATS-Verbindung, unverändertes
Bestandsverhalten. Anders als `CDC_HTTP_ADDR` ist eine **gesetzte** URL eine
**Start-Vorbedingung** des Feed-Containers: schlägt die Verbindung fehl, startet
der Prozess nicht (Fehlerklasse `configuration`, siehe
[Konfiguration](#5-konfiguration)). Eine gesetzte HTTP-Adresse ist das nicht.

**Das Subjekt:** Ein Wecksignal wird je Transaktion und distinkter berührter
Tabelle auf dem tabellen-granularen Subjekt
`cdc.changes.<source_id>.<schema>.<table>` publiziert — `<source_id>` ist die
konfigurierte `CDC_SOURCE_ID`, `<schema>`/`<table>` sind die
Klartext-Bezeichner der Tabelle. Ein Consumer, der alle Tabellen einer Quelle
verfolgt, abonniert die Wildcard `cdc.changes.<source_id>.>`; ein Consumer, der
mehrere Quellen verfolgt, `cdc.changes.>`.

**Der leere Payload:** Das Signal trägt **keine** Daten — keinen Change-Inhalt,
keine Positionsangabe. Jede Nachricht bedeutet ausschließlich „lies erneut über
den bestehenden Zugriffsweg"; Fehlen oder Verdopplung einer Nachricht trägt
keine eigene Bedeutung.

**Zustellsemantik und Nachvollziehbarkeit:** Es gibt **keine** Zustellgarantie
(Core NATS, Fire-and-Forget) und kein Replay — ein nicht verbundener oder
gerade getrennter Consumer verpasst das Signal ersatzlos. Verpasste Änderungen
bleiben wie bei den Stream-Abschnitten über
[Änderungen lesen](#änderungen-lesen) und die bestätigte Consumer-Position
([Position bestätigen](#position-bestätigen)) nachholbar; das Wecksignal
ersetzt diesen Zugriffsweg nicht.

**Der zweiseitige Ablauf:** Auf das Subjekt lauschen → beim Weckruf die
Änderung **selbst** über die [HTTP-/JSON-API](#zugriff-über-die-http-json-api)
holen und ausgeben — die Abfrage ist der `reader`-Endpunkt `GET /changes` mit
den drei Filtern aus dem Subjekt (`source`, `schema`, `table`).

**Beispiele:** Jede Sprache abonniert dasselbe Subjekt und holt die Änderung
über denselben `GET /changes`-Aufruf; NATS-URL, Adresse und Token liest jedes
Beispiel aus `CDC_NATS_URL`, `CDC_HTTP_ADDR` und `CDC_API_TOKEN_READER` und
lässt sich per Flag übersteuern.

- **Go:** `examples/nats-client` — Container-Aufruf über
  `make example-run-go SURFACE=nats ARGS="-source <quelle> -schema <schema> -table <tabelle>"`
  (baut bei Bedarf `pg-change-feed-examples:go-nats` aus `examples/Dockerfile`)
- **C#:** `examples/csharp/nats-client` — Container-Aufruf über
  `make example-run-csharp SURFACE=nats ARGS="--source <quelle> --schema <schema> --table <tabelle>"`
  (startet das mit `make examples-csharp` gebaute Image)
- **Kotlin:** `examples/kotlin/nats-client` — Container-Aufruf über
  `make example-run-kotlin SURFACE=nats ARGS="--source <quelle> --schema <schema> --table <tabelle>"`
  (startet das mit `make examples-kotlin` gebaute Image)

Die Beispiele sind zum Lesen und Nachbauen gedacht; die E2E-Testclients des
Harness liegen unter `tools/harness/` und sind kein Vorbild.

### Zugriff über den NATS-Vollinhalts-Stream

Anders als das Wecksignal im vorigen Abschnitt ist dieser Zugriffsweg
**daten-tragend**: jede Nachricht trägt den vollständigen Change-Inhalt, kein
reines Trigger-Signal. Er ist ein dritter, unabhängig nutzbarer Zustellweg
neben gRPC und SSE — dasselbe Nachrichtenschema wie beim
[SSE-Stream](#zugriff-über-server-sent-events), hier über NATS statt HTTP
zugestellt.

**Erreichbarkeit:** Zwei-Bedingungen-Aktivierung — **sowohl** `CDC_NATS_URL`
**als auch** `CDC_NATS_STREAM_TOKEN` müssen gesetzt sein. Ist nur
`CDC_NATS_URL` gesetzt (das bestehende Wecksignal, siehe oben), bleibt dieser
dritte Weg deaktiviert; ist nur `CDC_NATS_STREAM_TOKEN` gesetzt, startet der
Feed-Container nicht (Fehlerklasse `configuration`).

**Das Subjekt:** Ein vollständiges Change-Ereignis wird je Zeilen-Änderung auf
dem tabellen-granularen Subjekt `cdc.stream.<source_id>.<schema>.<table>`
publiziert — derselbe vierstufige Aufbau wie beim Wecksignal, aber ein
eigener Wurzel-Token (`cdc.stream` statt `cdc.changes`): beide Fähigkeiten
bleiben unabhängig voneinander abonnierbar.

**Die Nachrichtenform:** Dasselbe JSON-Schema wie beim SSE-Stream (zehn
Felder: `change_id`, `transaction_id`, `source_table_id`, `sequence`,
`operation`, `old_image`, `new_image`, `schema_version`, `schema`, `table`)
als Bytes des JSON-Dokuments — kein drittes Nachrichtenschema für dieselben
Daten.

**Zustellsemantik:** Core NATS, Fire-and-Forget, kein Replay — dieselbe
Zusicherung wie gRPC und SSE: ein nicht verbundener oder gerade getrennter
Consumer verpasst die Änderung ersatzlos und holt sie über
[Änderungen lesen](#änderungen-lesen) nach.

**Die Auth-Nebenwirkung — wichtig beim Einschalten:** Ein gesetzter
`CDC_NATS_STREAM_TOKEN` verlangt vom NATS-Server denselben Token **serverweit**
— auch für die bislang anonyme Wecksignal-Verbindung. Wer diesen dritten Weg
aktiviert, betreibt den NATS-Server danach nicht mehr ohne Verbindungs-Token;
ein Verbindungsversuch ohne oder mit falschem Token wird vom Server selbst
abgelehnt (Verbindungsebene, nicht Anwendungsebene) — dieselbe Aussage wie
gRPCs `Unauthenticated` oder SSEs `401`, nur eine Ebene tiefer verortet.

**Beispiele:** Jede Sprache abonniert den Vollinhalts-Namensraum
`cdc.stream.>` und gibt jede empfangene Change aus; NATS-URL und Token liest
jedes Beispiel aus `CDC_NATS_URL` und `CDC_NATS_STREAM_TOKEN` und lässt sich
per Flag übersteuern.

- **Go:** `examples/nats-stream-client` — Container-Aufruf über
  `make example-run-go SURFACE=nats-stream` (baut bei Bedarf
  `pg-change-feed-examples:go-nats-stream` aus `examples/Dockerfile`)
- **C#:** `examples/csharp/nats-stream-client` — Container-Aufruf über
  `make example-run-csharp SURFACE=nats-stream` (startet das mit
  `make examples-csharp` gebaute Image)
- **Kotlin:** `examples/kotlin/nats-stream-client` — Container-Aufruf über
  `make example-run-kotlin SURFACE=nats-stream` (startet das mit
  `make examples-kotlin` gebaute Image)

Die Beispiele sind zum Lesen und Nachbauen gedacht; die E2E-Testclients des
Harness liegen unter `tools/harness/` (`natsstreamsub`) und sind kein
Vorbild.

**SDK:** .NET-Anwendungen können statt des Beispiels dasselbe offizielle
NuGet-Package `PgChangeFeed.Client` einbinden (`LH-FA-SST-009`, `ADR-0106`,
`dotnet add package PgChangeFeed.Client`) — `PgChangeFeedNatsStreamClient.StreamChangesAsync`
abonniert den Vollinhalts-Namensraum (Default `cdc.stream.>`, oder ein über
`BuildSubject`/`BuildSourceSubject` eingeschränktes Subjekt) und liefert ein
`IAsyncEnumerable<Change>` mit allen zehn Feldern der Tabelle oben; die
Authentifizierung ist verbindungsseitig (derselbe `CDC_NATS_STREAM_TOKEN`
wie oben), ein abgelehnter Verbindungsversuch endet die Aufzählung mit einer
NATS-eigenen Ausnahme, statt den Draht-Vertrag selbst zu implementieren;
siehe `sdks/csharp/README.md`.

Kotlin/JVM-Anwendungen können statt des Beispiels dasselbe offizielle
Gradle-/Maven-Package `pgchangefeed-kotlin` einbinden (`LH-FA-SST-009`,
`ADR-0109`, Koordinate `io.github.pt9912:pgchangefeed-kotlin`) —
`PgChangeFeedNatsStreamClient.streamChanges()` abonniert den
Vollinhalts-Namensraum (Default `cdc.stream.>`, oder ein über `buildSubject`/
`buildSourceSubject` eingeschränktes Subjekt) und liefert eine
`Sequence<Change>` mit allen zehn Feldern der Tabelle oben; die
Authentifizierung ist verbindungsseitig (derselbe `CDC_NATS_STREAM_TOKEN`
wie oben), ein abgelehnter Verbindungsversuch endet die Sequenz mit der
zugrunde liegenden `io.nats.client`-Ausnahme unverändert, statt den
Draht-Vertrag selbst zu implementieren oder eine zweite Fehlerklassen-Hierarchie
zu erfinden. Wie beim HTTP-API-Zugriff oben verlangt der Bezug über
**GitHub Packages** immer eine Authentifizierung, auch für dieses
öffentliche Package (`ADR-0109` Festlegung 2). Siehe
`sdks/kotlin/pgchangefeed-kotlin/README.md`.

Python-Anwendungen können statt des Beispiels das offizielle
PyPI-Package `pgchangefeed` einbinden (`LH-FA-SST-009`, `ADR-0110`,
`pip install pgchangefeed`) — `PgChangeFeedNatsStreamClient.stream_changes()`
abonniert den Vollinhalts-Namensraum `cdc.stream.<source_id>.>` (alle
Tabellen einer Quelle; `source_id` ist Konstruktor-Argument) und liefert
einen Iterator über die getypten `StreamChange`-Events mit allen zehn
Feldern der Tabelle oben; die Authentifizierung ist verbindungsseitig
(derselbe `CDC_NATS_STREAM_TOKEN` wie oben), ein abgelehnter
Verbindungsversuch endet am Connect-Fehler der NATS-Bibliothek unverändert,
statt den Draht-Vertrag selbst zu implementieren. Mit dieser Fläche deckt
das Package dieselbe Vier-Wege-Matrix wie die C#-/Kotlin-Pendants. Siehe
`sdks/python/README.md`.

## 5. Konfiguration

### Umgebungsvariablen des Feed-Containers

| Variable | Pflicht | Bedeutung |
|---|---|---|
| `CDC_CAPTURE_DSN` | ja | Verbindung über die Rolle `cdc_capture` (Store-Adapter, Replication-Stream, Ausführung eines Backfills) |
| `CDC_ADMIN_DSN` | ja | Verbindung über die Rolle `cdc_admin` (Tabellen-Aktivierung, Heartbeat, Verarbeitung der Antrags-Queue, Annahme eines Backfill-Antrags, `register-consumer`/`acknowledge-consumer`) |
| `CDC_READER_DSN` | ja | Verbindung über die Rolle `cdc_reader` (`--healthcheck`, `diagnose`) |
| `CDC_SOURCE_ID` | ja | Kennung der Quelle (muss in `cdc.source` registriert sein) |
| `CDC_PUBLICATION` | ja | Name der PostgreSQL-Publication |
| `CDC_SLOT` | ja | Name des Logical-Replication-Slots |
| `CDC_TABLES` | ja, falls keine Konfigurationsdatei dieselbe Aktivierung trägt | Aktivierte Tabellen, Format `schema.tabelle=tabelle-id:schema-version-id`, kommagetrennt |
| `CDC_LOG_LEVEL` | nein | Log-Level des strukturierten JSON-Loggers (Default `info`) |
| `CDC_NATS_URL` | nein | NATS-Server-URL für das Change-Notification-Wecksignal (`cdc.changes.<source_id>.<schema>.<table>`, tabellen-granular, leerer Payload, `ADR-0056`); ungesetzt bleibt das Feature vollständig deaktiviert, gesetzt ist eine erfolgreiche Verbindung Vorbedingung des Starts (Fehlerklasse `configuration`). Zusammen mit `CDC_NATS_STREAM_TOKEN` aktiviert dieselbe Variable zusätzlich den dritten, vollinhaltstragenden NATS-Zustellweg (`ADR-0100`, siehe [Zugriff über den NATS-Vollinhalts-Stream](#zugriff-über-den-nats-vollinhalts-stream)) |
| `CDC_NATS_STREAM_TOKEN` | nein | Verbindungs-Token des dritten, vollinhaltstragenden NATS-Zustellwegs (`ADR-0100`); wirkt nur zusammen mit gesetztem `CDC_NATS_URL` — ist nur `CDC_NATS_STREAM_TOKEN` gesetzt, aber `CDC_NATS_URL` leer, startet der Container nicht (Fehlerklasse `configuration`). Ein gesetzter Wert verlangt vom NATS-Server denselben Token **serverweit**, auch für die Wecksignal-Verbindung (siehe dortiger Abschnitt) |
| `CDC_HTTP_ADDR` | nein | Horch-Adresse der HTTP-/JSON-API (`host:port`); ungesetzt bleibt die API vollständig deaktiviert — kein Server, keine zusätzliche Verbindung. Anders als `CDC_NATS_URL` ist die Adresse **keine** Start-Vorbedingung: Ist sie gesetzt, öffnet der Prozess den Server in eigener Goroutine und läuft unverändert weiter; scheitert das Binden der Adresse (z. B. belegter Port), meldet er das im Log und der Erfassungsbetrieb bleibt davon unberührt |
| `CDC_API_TOKEN_READER` | nein | Bearer-Token der lesenden Rechtsklasse der HTTP- und gRPC-API; ungesetzt (leer) ist die Klasse nicht konfiguriert — ein Aufruf mit einem Token, das keiner konfigurierten Klasse entspricht, endet `401` |
| `CDC_API_TOKEN_ADMIN` | nein | Bearer-Token der administrativen Rechtsklasse der HTTP- und gRPC-API (deckt die lesende Klasse implizit mit ab); leer bedeutet dieselbe Deaktivierung wie bei `CDC_API_TOKEN_READER` |
| `CDC_GRPC_ADDR` | nein | Horch-Adresse des gRPC-Streaming-Servers (`host:port`); ungesetzt bleibt der Streaming-Server vollständig deaktiviert — kein Listener. Wie `CDC_HTTP_ADDR` keine Start-Vorbedingung |
| `CDC_CONFIG_FILE` | nein | Pfad zu einer optionalen YAML-Konfigurationsdatei (siehe unten) |

Fehlt eine Pflichtvariable und liefert auch keine Konfigurationsdatei
einen Wert für dasselbe Feld, startet der Container nicht (Fehlerklasse
`configuration`).

### Optionale YAML-Konfigurationsdatei (`CDC_CONFIG_FILE`)

Additiv zu den Umgebungsvariablen (`ADR-0052`, `ADR-0088`, `SPEC-016`): Ist
`CDC_CONFIG_FILE` gesetzt, liest der Container zusätzlich eine
YAML-Datei unter diesem Pfad (read-only in den Container gemountet). Jede
gesetzte Umgebungsvariable überschreibt das gleichnamige Feld der Datei
einzeln — Env-Var schlägt Datei, Feld für Feld. Ist `CDC_CONFIG_FILE`
nicht gesetzt, ändert sich am Env-only-Betrieb oben nichts.

Die zwei Oberflächen-Adressen sind Datei-Felder: `http_addr` und
`grpc_addr` tragen je eine Horch-Adresse in der Form `host:port` und
entsprechen `CDC_HTTP_ADDR` bzw. `CDC_GRPC_ADDR` — mit demselben
Feld-für-Feld-Vorrang. Trägt keine der beiden Quellen eine Adresse, bleibt
die jeweilige Oberfläche deaktiviert (dieselbe No-Op-Semantik wie bei der
entsprechenden Umgebungsvariable oben).

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
http_addr: ":8090"
grpc_addr: ":9090"
```

**Wichtig — Zugangsdaten bleiben env-var-exklusiv:** Die Schlüssel
`capture_dsn`, `admin_dsn`, `reader_dsn`, `api_token_reader`,
`api_token_admin`, `nats_url` und `nats_stream_token` dürfen in dieser
Datei **nicht** vorkommen. Ein Treffer bricht das Laden mit einer eigenen,
den Grund nennenden Zeile ab (Fehlerklasse `configuration`) — eine
Konfigurationsdatei landet typischerweise in Kanälen (Repository,
ConfigMap, Backup), die für Zugangsdaten nicht vorgesehen sind. Die Grenze
ist die **Form** des Feldes, nicht sein Wert: `http_addr`/`grpc_addr` sind
`host:port` und können keine Zugangsdaten tragen, `nats_url` ist eine URL
und kann Benutzer sowie Passwort einbetten (`nats://benutzer:passwort@host:4222`),
`nats_stream_token` trägt denselben Zugangsdaten-Charakter wie die beiden
API-Token-Schlüssel (`ADR-0100`). Ein unbekannter Schlüssel bricht das
Laden ebenfalls ab (striktes Decoding).

Die env-exklusiven Variablen `CDC_NATS_URL`, `CDC_NATS_STREAM_TOKEN`,
`CDC_API_TOKEN_READER` und `CDC_API_TOKEN_ADMIN` werden **auch unter
geladener Datei** aus der Umgebung gelesen — sie haben kein
Datei-Gegenstück, ihre Herkunft ist die Umgebungsvariable auf beiden Wegen;
die Datei kann sie weder setzen noch überschreiben.

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
| `transient` | vorübergehend nicht verfügbare Quelle/Speicher | Erneuter Versuch mit begrenztem Backoff; im Erfassungspfad deklariert, aktuell von keinem Adapter konstruiert — ein Backfill-Run trägt sie als Klasse seines Fehlertexts |
| `configuration` | ungültige oder fehlende Umgebungsvariable | Kein Start, sichtbarer Fehler |
| `permission` | fehlende Berechtigung | Sichtbarer Fehler, kein stiller Retry; im Erfassungspfad deklariert, aktuell von keinem Adapter konstruiert — ein Backfill-Run trägt sie als Klasse seines Fehlertexts (z. B. fehlendes `SELECT` auf die Quelltabelle) |
| `schema` | eine Replikationsnachricht ist nicht sicher interpretierbar (z. B. TRUNCATE, unbekannter Nachrichtentyp) | Sichtbarer Fehler, kein stilles Überspringen |
| `storage` | Persistenzfehler | Kein Source-ACK, damit keine Änderung verloren geht |
| `replication` | zwei Unterarten (`ADR-0049`): **Stream-Ordnungs-Verletzung** (z. B. Commit ohne offene Transaktion) oder **Transport-/Verbindungsstörung** (Verbindungsaufbau, Slot, Keepalive, Quell-Bestätigung) | Stream-Ordnungs-Verletzung: sofortiger, sichtbarer Abbruch, unabhängig vom WAL-Rückstand. Transport-/Verbindungsstörung: Schwellen-Überwachung über den WAL-Rückstand (siehe [WAL-Rückstand prüfen](#wal-rückstand-prüfen)) — kontrollierte Fortsetzung unterhalb 1 GiB, sichtbarer Abbruch darüber |
| `internal` | unerwarteter interner Fehler, der keiner anderen Klasse zuzuordnen ist | Sichtbarer Fehler; realer Fallback für jeden nicht erkannten Fehler |

`transient` und `permission` gehören zur deklarierten Menge der sieben
Klassen, werden im Erfassungspfad aber von keinem Adapter aktuell konstruiert
— beobachtbar sind sie heute nur als Klasse im Fehlertext eines
fehlgeschlagenen Backfill-Runs (`cdc.backfill_status.error_message`).
`internal` ist dagegen der real erreichbare
Fallback-Zweig: Jeder Fehler, der keiner der übrigen sechs Klassen
zugeordnet werden kann, fällt auf `internal` zurück.

Ein Fehler des Erfassungspfads jeder Klasse beendet den Container-Prozess mit Ausgang 1 (ein Fehler eines Backfill-Runs ist run-lokal und beendet ihn nicht); der
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
| Backfill | Die einmalige Überführung des Tabellenbestands (Zeilen, die zum Startzeitpunkt des Runs bestehen) als `INSERT`-Änderungen mit `origin = 'backfill'`; ausdrücklich über `cdc.backfill_table` ausgelöst |
| Run (Backfill) | Ein Backfill-Durchlauf für eine Tabelle; sein Zustand steht in `cdc.backfill_run`, gelesen über `cdc.backfill_status` |
| Angenommen (`applied` bei `backfill`) | Der Antrag ist angenommen und die Run-Zeile `queued` angelegt; sagt nichts über die Ausführung des Runs |
| Geschätzte Zeilenzahl | Eine Schätzung des Katalogs der Quelle für die Zeilenzahl der Tabelle (`estimated_rows`); keine Zählung und keine Grenze, NULL heißt unbekannt |

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
  zählen gegen `max_wal_senders` der Quell-Instanz. Während der Slot-Anlage
  eines Backfills kommt eine weitere hinzu (die Replication-Verbindung des
  temporären Slots, sie endet nach dem Import des Snapshots); sie zählt ebenfalls
  gegen `max_wal_senders`, der temporäre Slot zusätzlich gegen
  `max_replication_slots`.
- **Backfill, Toleranz der Kopierdauer:** 10 Minuten. Das ist ein **Startwert,
  Setzung ohne Messung**; er steuert allein die Warnung `warn_duration` (siehe
  [Bestand als Backfill überführen](#bestand-als-backfill-überführen)), keine
  Anforderung und kein Gate liest ihn.
- **Backfill, Richtgröße:** 4.000.000 **geschätzte** Zeilen. Das ist eine
  **Orientierung, keine Grenze**: ein Antrag wird nie abgelehnt, ein Run nie
  abgebrochen; liegt die geschätzte Zeilenzahl beim Antrag darüber, setzt der
  Antrag `warn_estimated_size`. *Ursprung (abgeleitet):* die Kopierrate der
  Stufe mit 200.000 Zeilen, Median von 3 Runs 7.693 Zeilen/s (Bereich 6.997 bis
  8.686; Lauf `20260925T000439Z`, übernommen aus dem Lauf-Bericht, im
  Repository nicht auflösbar), mal die Toleranz von 600 s ergibt 4.615.800
  Zeilen, auf eine Stelle abgerundet: 4.000.000. Keine gemessene Stufe
  kopierte die Toleranzdauer; die Rate ist hochgerechnet. Die Rate streut
  zwischen Läufen; dieselbe Rechnung an zwei weiteren Läufen desselben Hosts
  (Stufe 200.000, Median von 3 Runs): 8.933 Zeilen/s im Lauf
  `20260925T011036Z` (gemessen im
  [Review-Report](../reviews/review-slice-backfill-bench-richtgroesse.md)),
  8.559 Zeilen/s im Lauf `20260925T012459Z` (übernommen aus dem Lauf-Bericht
  des Implementers, im Repository nicht auflösbar), also je 5.000.000
  Zeilen nach Rundung; der Lauf `20260925T015600Z` (gemessen, gedruckt im
  [Verifikations-Report](../reviews/verifikation-slice-backfill-bench-richtgroesse.md)
  §3) ergibt 8.654 Zeilen/s und ebenfalls 5.000.000 Zeilen. Der Lauf
  `20260925T032925Z` (übernommen aus dem Lauf-Bericht des Implementers, im
  Repository nicht auflösbar) ergibt bei 5.532 Zeilen/s (Bereich 4.778 bis
  6.631) 3.319.200 Zeilen, abgerundet 3.000.000. Der Lauf `20260925T043056Z`
  (gemessen, gedruckt im
  [Review-Report](../reviews/review-slice-backfill-slot-leerlauf-bestaetigung.md);
  der Host trug dabei Last fremder Container) ergibt bei 4.504 Zeilen/s
  2.702.400 Zeilen, abgerundet 2.000.000. Die Spanne dieser sechs Läufe
  reicht von 2.000.000 bis 5.000.000 Zeilen nach Rundung; der Wert im Code
  (4.000.000) liegt innerhalb dieser Spanne und über den zwei niedrigsten
  Werten (3.000.000 und 2.000.000); die Konstante ist ein Startwert, den eine
  weitere Messung nachschärfen kann. Sechs weitere Läufe auf demselben Host lagen in den Stufen
  ab 100.000 Zeilen zwischen 4.088 und 9.425 Zeilen/s je Run (übernommen aus
  den Lauf-Berichten, nicht im Repository). Die Zahl gilt für den Host und die
  Bedingungen der Messung (siehe unten). Breite Zeilen (`jsonb`, `bytea`),
  andere Hardware und eine andere Einfügeform sind für die Kopierrate
  ungemessen. Die Richtgröße folgt der Kopierdauer; der Speicher des
  Feed-Containers geht nicht ein, weil er nicht an der Zahl der Changes hängt
  (siehe *Backfill, Speicher des Feed-Containers*).
- **Backfill, gemessene Werte** (Lauf `20260925T000439Z` von
  `tools/bench-backfill.sh`, Vertrag in
  [`harness/targets/bench-backfill.md`](../../harness/targets/bench-backfill.md);
  aus dem Lauf-Bericht übernommen, im Repository nicht auflösbar;
  Host: Linux 6.8.0-139-generic, Docker 29.8.1, 20 CPU, 31 GiB RAM,
  PostgreSQL 18 (`postgres:18-alpine`) mit Standard-Einstellungen; Tabellen mit
  fünf schmalen Spalten, mittlere Zeilenbreite etwa 74 Bytes, Blockgröße 1.000
  Zeilen, zeilenweise Einfügung in **eine** Transaktion; Median und Bereich von
  je 3 Runs):

  | Zeilen | Kopierdauer | Durchsatz |
  |---|---|---|
  | 10.000 | 1,19 s (1,14 bis 1,20 s) | 8.375 Zeilen/s |
  | 50.000 | 5,89 s (5,58 bis 5,96 s) | 8.489 Zeilen/s |
  | 200.000 | 26,0 s (23,0 bis 28,6 s) | 7.693 Zeilen/s |

  Eine Nachmessung am eigenen Lauf `20260925T012459Z` (gleicher Host, gleiche
  Stufen, Median von 3 Runs; übernommen aus dem Lauf-Bericht des
  Implementers, im Repository nicht auflösbar) lieferte 8.382, 8.975 und
  8.559 Zeilen/s. Zwei Runs über je 1.000.000 Zeilen (Lauf `20260924T233628Z`,
  `--full`, gleicher Host; übernommen, im Repository nicht auflösbar) dauerten
  125,2 s und 122,8 s (7.986 und 8.141 Zeilen/s). Die
  mittlere Blockdauer liegt bei etwa 0,12 bis 0,13 s (abgeleitet: Dauer durch
  Blockzahl); die längste Dauer eines einzelnen Blocks und die Dauer des Commits
  sind nicht gemessen.
- **Backfill, Speicher des Feed-Containers.** Der Speicher des Feed-Containers
  hängt nicht an der Zahl der Changes, die für die Quelle in `cdc.change` stehen:
  der periodische Bereinigungslauf (alle 10 s, siehe
  [Aufbewahrung (Retention)](#aufbewahrung-retention)) liest sie seitenweise
  (10.000 Changes je Seite) und ohne Row Images. Gemessen wurde mit
  `tools/bench-backfill-memory.sh` (Vertrag in
  [`harness/targets/bench-backfill.md`](../../harness/targets/bench-backfill.md);
  Host: Linux 6.8.0-139-generic, Docker 29.8.1, 20 CPU, 31 GiB RAM, PostgreSQL 18
  mit Standard-Einstellungen; schmale Zeilen von 74 Bytes; drei Runs über je
  1.000.000 Zeilen hintereinander, ein frischer Feed-Container je Run,
  `--memory 6g`; zwei Läufe, `20260925T183239Z` und `20260925T184503Z`; gedruckte
  Zeilen im
  [Messbericht](../reviews/messbericht-slice-retention-lauf-speicher-begrenzung.md)).
  Die Spitze des Feed-Containers (`memory.peak`, gedruckt 60 s nach dem Run) gegen
  die Zahl der Changes, die nach dem Run in `cdc.change` stehen:

  | Changes in `cdc.change` | Spitze in MiB (gemessen, zwei Läufe) |
  |---|---|
  | 1.000.000 | 14,9 und 15,1 |
  | 2.000.000 | 16,5 und 16,9 |
  | 3.000.000 | 17,3 und 17,6 |

  Zwischen 1.000.000 und 3.000.000 Changes wächst die Spitze um 2,4 und 2,5 MiB,
  etwa 1,3 Bytes je Change (abgeleitet); die Ursache dieses Wachstums ist nicht
  untersucht. Die Bereinigungs-Takte laufen in allen sechs Runs bis zum Ende weiter
  (gedruckt: 20 bis 23 Zeilen „Bereinigung gelaufen“ seit dem Start des
  Feed-Containers, keine fehlgeschlagene). Bei leerem `cdc.change` liegt die Spitze
  des Prozesses im Run bei 8,9 bis 10,5 MiB, von 10.000 bis 1.000.000 Zeilen
  (Blockgröße 1.000; gemessen, [Messbericht der
  Untersuchung](../reviews/messbericht-slice-backfill-speicher-untersuchung.md),
  Abschnitt 3.2). Ungemessen bleiben breite Zeilen, ein Speicherlimit des
  Containers unter 6 GiB und die laufende Erfassung als Quelle der Changes; die
  Seite trägt keine Row Images, ihr Bedarf hängt deshalb nicht an der Zeilenbreite
  (aus dem Aufbau der Seite hergeleitet, nicht gemessen).

  **Vorversionen.** Die Server-Versionen `v0.1.0` bis `v0.1.2` lesen in jedem Takt
  alle Changes der Quelle samt Row Images in den Speicher (Quelltext der drei
  Versionen verglichen, nicht am Image der Version gemessen). Ihr Speicher wächst
  mit der Zahl der Changes: am Stand vor der seitenweisen Lesung gemessen 1,03 bis
  1,59 KiB je Change bei schmalen Zeilen (abgeleitet, ab 100.000 Changes) und
  1.082,7 bis 1.273,5 MiB bei 1.000.000 Changes ([Messbericht der
  Untersuchung](../reviews/messbericht-slice-backfill-speicher-untersuchung.md),
  Abschnitt 3.3).
- **Backfill, WAL-Rückstand des Capture-Slots** (Größe von
  `cdc_wal_retention_bytes`, siehe [WAL-Rückstand prüfen](#wal-rückstand-prüfen)).
  Der Feed bestätigt WAL ohne Inhalt für die Publication im Leerlauf seines
  Streams. Gemessen (Lauf `20260925T032925Z` von `tools/bench-backfill.sh`,
  übernommen aus dem Lauf-Bericht des Implementers, im Repository nicht
  auflösbar; nachgemessen im Lauf `20260925T043056Z`, gedruckt im
  [Review-Report](../reviews/review-slice-backfill-slot-leerlauf-bestaetigung.md),
  mit gleichem Ergebnis; Vertrag in
  [`harness/targets/bench-backfill.md`](../../harness/targets/bench-backfill.md);
  Host und Tabellen wie oben, PostgreSQL 18; Rückstand im Abstand von 1 bis 2 s
  gelesen und auf ganze MiB gerundet): in allen neun Runs der Stufen mit 10.000,
  50.000 und 200.000 Zeilen (je 3 Runs) lag die Spitze des Rückstands im Run bei
  0 MiB, unmittelbar nach dem Run ebenfalls bei 0 MiB; das Skript schrieb
  zwischen den Runs nichts. Ohne die Bestätigung — gemessen an einem
  Stand, der WAL ohne Inhalt für die Publication nicht bestätigte (Lauf
  `20260925T000439Z`, Stufe 200.000 Zeilen, übernommen aus dem Lauf-Bericht, im
  Repository nicht auflösbar; ebenso die Läufe `20260925T012459Z`,
  `20260924T233628Z`) — lag die Spitze im Run im Median bei 140 MiB (etwa 735
  Bytes je Zeile), bei 719 und 834 MiB in den zwei Runs über je 1.000.000 Zeilen,
  und der Rückstand blieb bis zu einem Commit auf einer aktivierten Tabelle
  bestehen (776 MiB unmittelbar nach dem ersten Run über 1.000.000 Zeilen); ein
  dritter Run über 1.000.000 Zeilen im selben CDC-Speicher überschritt vor
  seinem Ende die Fehlerschwelle von 1 GiB (Messwert 1.098.218.616 Bytes,
  Feed-Container mit Ausgang 1, Run `interrupted`). Ein Schreiber auf eine nicht
  aktivierte Tabelle erzeugte dort ohne Run 33,5 MiB WAL je 200.000 Zeilen (175
  Bytes je Zeile; Lauf `wal-verdikt[20260925T004731Z]`, gedruckt in
  [`ADR-0120`](../plan/adr/0120-capture-slot-leerlauf-bestaetigung.md)
  §Gemessen; Architekt-Verdikt
  [`architect-verdict-backfill-wal-rueckstand-und-bench-rot`](../reviews/architect-verdict-backfill-wal-rueckstand-und-bench-rot.md)).
  **Grenze der Ein-Transaktions-Form:** das vom Slot auf der Platte der Quelle
  **gehaltene** WAL (`restart_lsn`) und der Spill des Walsenders bleiben. Das
  gehaltene WAL erreichte in der Stufe mit 200.000 Zeilen im Median 140 MiB
  (Spitze im Run, Lauf `20260925T032925Z`, übernommen, im Repository nicht
  auflösbar; 31 MiB bei 50.000 und 7 MiB bei 10.000 Zeilen) und 141 MiB (35
  MiB bei 50.000, 7 MiB bei 10.000 Zeilen; Lauf `20260925T043056Z`, gedruckt im
  [Review-Report](../reviews/review-slice-backfill-slot-leerlauf-bestaetigung.md)),
  in den zwei Runs über je 1.000.000 Zeilen 782 und 1.613 MiB
  (Lauf `20260924T233628Z`, übernommen, im Repository nicht auflösbar). Der
  Walsender lagerte bei einem Run über 200.000 Zeilen 79 MB der offenen
  Transaktion aus (`spill_bytes` des Slots, `logical_decoding_work_mem` 64 MB;
  Lauf `spill-verdikt[20260925T005729Z]`, übernommen aus
  [`ADR-0120`](../plan/adr/0120-capture-slot-leerlauf-bestaetigung.md)
  §Gemessen). Beides wächst mit der Größe der Tabelle; die Bestätigung im
  Leerlauf entlastet den Capture-Pfad, nicht die Platte der Quelle. Bemessen
  Sie den Plattenplatz der Quelle danach (abgeleitet: 1.613 MiB gehaltenes WAL
  bei 1.000.000 Zeilen von etwa 74 Bytes, gut das Zwanzigfache der
  Zeilenbytes).
- **Backfill, Wirkung auf die Live-Erfassung** (je ein Run über 200.000 Zeilen
  bei 100 Live-Änderungen/s in eine andere aktivierte Tabelle, dazu eine
  Referenz gleicher Dauer ohne Run; **jeder Vergleich ist ein Einzellauf ohne
  Wiederholung**): `cdc_capture_lag` im Run gegenüber ohne Run — Lauf
  `20260925T000439Z` (übernommen, im Repository nicht auflösbar): 0,049 bis
  1,005 s (Median 0,48 s, 24 Proben) gegenüber 0,049 bis 1,022 s (Median 0,50
  s, 29 Proben); Lauf `20260925T011036Z`
  ([Review-Report](../reviews/review-slice-backfill-bench-richtgroesse.md)):
  0,123 bis 1,074 s (Median 0,59 s, 21 Proben) gegenüber 0,041 bis 1,011 s
  (Median 0,40 s, 25 Proben); Lauf `20260925T012459Z` (übernommen, im Repository nicht auflösbar): 0,090 bis
  0,961 s (Median 0,456 s, 22 Proben) gegenüber 0,050 bis 1,008 s (Median
  0,413 s, 27 Proben). Die Mediane liegen in beiden Richtungen auseinander, die
  Maxima bei etwa 1 s; bei dieser Last und Größe ist aus diesen drei
  Einzelläufen kein Unterschied ableitbar.
- **Backfill, Schätzung der Zeilenzahl** (`pg_class.reltuples`, PostgreSQL 18,
  Tabelle mit 100.000 Zeilen, Lauf `20260925T000439Z`, übernommen und im
  Repository nicht auflösbar; dieselben Werte in den Läufen
  `20260925T011036Z` ([Review-Report](../reviews/review-slice-backfill-bench-richtgroesse.md))
  und `20260925T012459Z`, übernommen, im Repository nicht auflösbar): eine frisch befüllte
  Tabelle trägt **keine** Schätzung (NULL, „unbekannt“) — mit und ohne
  Autovacuum; mit Autovacuum lag die Schätzung nach 35 s vor (Abfrage im
  Abstand von 5 s), ohne Autovacuum blieb sie unbekannt; nach `ANALYZE` stimmte
  sie (100.000); nach 20.000 weiteren Zeilen ohne erneutes `ANALYZE` lag sie
  16,7 % unter der tatsächlichen Zahl (120.000). Die Schätzung ist so alt wie
  das letzte `ANALYZE`; `warn_estimated_size` schweigt bei einer unbekannten
  Schätzung.

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
| 1.8 | 2026-09-13 | Sichtbarkeit blockierender Consumer ergänzt (`LH-FA-RET-005`, slice-045): §4 neuer Abschnitt „Blockierende Consumer erkennen" (`cdc.retention_blockers`) |
| 1.9 | 2026-09-13 | `cdc_storage_bytes`-Metrik ergänzt (`LH-FA-RET-006`, slice-046): §4 „Metriken lesen" nennt die neue `cdc.metrics`-Zeile |
| 1.10 | 2026-09-13 | `diagnose`-Ausgabe um Retention-Sichtbarkeit erweitert (`LH-FA-SST-003`, deckt `LH-FA-RET-005`/`006`, slice-047): §4 „Diagnose ausführen" trägt jetzt den aktuell blockierenden Consumer je Quelle (inkl. „kein Blocker"-Fall) und `cdc_storage_bytes` |
| 1.11 | 2026-09-13 | `CDC_NATS_URL`-Zeile (§5) auf das tabellen-granulare Subjekt-Schema `cdc.changes.<source_id>.<schema>.<table>` korrigiert (`ADR-0056`, slice-058) |
| 1.12 | 2026-09-14 | Diagnose-Beispielausgabe (§4) auf den umbenannten E2E-Quellnamen `src-e2e` aktualisiert (reines Namensrelikt aus der ursprünglichen MVP-Testumgebung, slice-057) |
| 1.13 | 2026-09-15 | Betreiber-Oberfläche nachgezogen: §5 um die HTTP-Gruppe (`CDC_HTTP_ADDR`, `CDC_API_TOKEN_READER`, `CDC_API_TOKEN_ADMIN`) und `CDC_GRPC_ADDR` erweitert (je mit Aktivierungs-/No-Op-Semantik); §4 um „Spalte vom Ausschluss konfigurieren" (`cdc.exclude_column`/`cdc.include_column`, `LH-FA-CFG-005`, dauerhafter Ausschlussstand) und die drei Netzwerk-Zugriffswege (HTTP-/JSON-API `LH-FA-SST-006`, gRPC-Change-Stream und Server-Sent-Events `LH-FA-SST-008`) |
| 1.14 | 2026-09-15 | Review-Nachzug: Rahmen-Aussage der HTTP-§4 auf die tatsächlich gelistete Fähigkeitsmenge gezogen (die Retention-Auslösung ist nicht CLI-/SQL-gleichwertig, sondern API-exklusiv); der gRPC-Abschnitt nennt die zehn Nachrichtenfelder, und der SSE-Abschnitt verweist darauf statt auf die Spaltenliste von `cdc.changes` |
| 1.15 | 2026-09-15 | Changes-Lesen über die API ergänzt (`LH-FA-SST-006`, `LH-FA-REA-001`…`006`, `ADR-0081`, slice-086): §4 Fähigkeits-Tabelle um `GET /changes` erweitert, Parameter-/Antwort-Beschreibung samt Fehlerfällen, „Änderungen lesen" verweist auf den Endpunkt, und die Zustellsemantik nennt die nicht streamende Form neben dem Live-Stream |
| 1.16 | 2026-09-15 | Vierter Zugriffs-Abschnitt ergänzt: §4 „Zugriff über das NATS-Wecksignal" (`LH-FA-SST-007`, `ADR-0055`/`ADR-0056`/`ADR-0079`, slice-083) — Subjekt-Schema, leerer Payload, Zustellsemantik und der zweiseitige Ablauf (lauschen, dann über `GET /changes` holen) samt Beispiel `examples/nats-client` |
| 1.17 | 2026-09-17 | Beispiel-Programme der HTTP-Familie in den Zugriffs-Abschnitten ergänzt (`ADR-0076`, slice-095): §4 „Zugriff über die HTTP-/JSON-API" nennt `examples/http-client` samt Startbefehl, „Zugriff über Server-Sent-Events" nennt `examples/sse-client`; beide lesen Adresse und Token aus `CDC_HTTP_ADDR` und `CDC_API_TOKEN_READER` |
| 1.18 | 2026-09-17 | Konfigurationsdatei nachgezogen (`ADR-0088`, `SPEC-016`, slice-096): §5.2 führt die zwei neuen Datei-Felder `http_addr`/`grpc_addr` samt Precedence, die Zugangsdaten-Klasse auf sechs Schlüssel gezogen (`capture_dsn`/`admin_dsn`/`reader_dsn`/`api_token_reader`/`api_token_admin`/`nats_url`, Grenze ist die Feld-Form) und festgehalten, dass `CDC_NATS_URL` und die zwei Token-Klassen auch unter geladener Datei aus der Umgebung wirken |
| 1.19 | 2026-09-17 | Beispiel-Programm der gRPC-Familie ergänzt (`ADR-0076`, `ADR-0060`, slice-095): „Zugriff über den gRPC-Change-Stream" nennt `examples/grpc-client` samt Startbefehl; es liest Adresse und Token aus `CDC_GRPC_ADDR` und `CDC_API_TOKEN_READER` |
| 1.20 | 2026-09-17 | Erster C#-Beispiel-Client ergänzt (`ADR-0087`, `ADR-0090`, slice-098): §4 „Zugriff über die HTTP-/JSON-API" trägt jetzt einen `**Beispiele:**`-Block mit einer Zeile je Sprache (Go, C#) statt eines einzelnen `**Beispiel:**`-Absatzes — die Ziel-Form für die volle Matrix; `examples/csharp/http-client` ruft denselben `reader`-Endpunkt `GET /tables` über einen Container-Aufruf (`make examples-csharp`) |
| 1.21 | 2026-09-17 | Erster Kotlin-Beispiel-Client ergänzt (`ADR-0087`, `ADR-0090`, slice-099): §4 „Zugriff über die HTTP-/JSON-API" — dritte Zeile im `**Beispiele:**`-Block; `examples/kotlin/http-client` ruft denselben `reader`-Endpunkt `GET /tables` über einen Container-Aufruf (`make examples-kotlin`) |
| 1.22 | 2026-09-17 | C#- und Kotlin-SSE-Client ergänzt (`ADR-0090`, slice-100): §4 „Zugriff über Server-Sent-Events" — `**Beispiel:**`-Absatz (nur Go) wird zu einem `**Beispiele:**`-Block mit einer Zeile je Sprache (Go, C#, Kotlin); `examples/csharp/sse-client` und `examples/kotlin/sse-client` öffnen denselben Endpunkt `GET /changes/stream` über einen Container-Aufruf (`make examples-csharp`/`make examples-kotlin`, Image-Tags `pg-change-feed-examples:csharp-sse`/`:kotlin-sse`) |
| 1.23 | 2026-09-17 | C#- und Kotlin-NATS-Client ergänzt (`ADR-0090`, `ADR-0055`/`ADR-0056`/`ADR-0079`, slice-101): §4 „Zugriff über das NATS-Wecksignal" — Fließtext-Absatz wird zu einem `**Beispiele:**`-Block mit einer Zeile je Sprache (Go, C#, Kotlin); `examples/csharp/nats-client` und `examples/kotlin/nats-client` lauschen auf dasselbe tabellen-granulare Subjekt und holen die Änderung über denselben `GET /changes`-Aufruf, über einen Container-Aufruf (`make examples-csharp`/`make examples-kotlin`, Image-Tags `pg-change-feed-examples:csharp-nats`/`:kotlin-nats`) |
| 1.24 | 2026-09-17 | Erster C#-gRPC-Client ergänzt (`ADR-0090`, `ADR-0060`, slice-102): §4 „Zugriff über den gRPC-Change-Stream" — `**Beispiel:**`-Absatz (nur Go) wird zu einem `**Beispiele:**`-Block mit einer Zeile je Sprache (Go, C#); `examples/csharp/grpc-client` öffnet denselben Server-Streaming-RPC `ChangeStream/StreamChanges` über einen Container-Aufruf (`make examples-csharp`, Image-Tag `pg-change-feed-examples:csharp-grpc`) — der C#-Stub entsteht dabei im Bau aus der `.proto`, gelesen über einen zusätzlichen, benannten Bau-Kontext (erste reale Bauprobe dieser Form, bislang nur isoliert gemessen) |
| 1.25 | 2026-09-17 | Kotlin-gRPC-Client ergänzt (`ADR-0090`, `ADR-0060`, slice-103): §4 „Zugriff über den gRPC-Change-Stream" — dritte und letzte Zeile im `**Beispiele:**`-Block; `examples/kotlin/grpc-client` öffnet denselben Server-Streaming-RPC `ChangeStream/StreamChanges` über einen Container-Aufruf (`make examples-kotlin`, Image-Tag `pg-change-feed-examples:kotlin-grpc`) — der Kotlin-Stub entsteht dabei im Bau aus der `.proto`, gelesen über denselben zusätzlichen, benannten Bau-Kontext, jetzt auf die Kotlin-Werkzeugkette übertragen (`protoc-gen-grpc-java`/`protoc-gen-grpc-kotlin`). Mit dieser Zeile ist die volle Matrix (vier Zugriffs-Oberflächen × drei Sprachen, zwölf Programme) im Handbuch vollständig |
| 1.26 | 2026-09-18 | Go-Startform auf `make`+Dockerfile umgestellt (`ADR-0098`, Supersedes `ADR-0076` Startform-Bullet, slice-beispiele-go-dockerfile-start): alle vier Go-Zeilen der `**Beispiele:**`-Blöcke zitieren jetzt `make example-run-go SURFACE=<oberfläche>` (baut bei Bedarf `pg-change-feed-examples:go[-<surface>]` aus dem neuen `examples/Dockerfile`) statt `go run ./examples/<name>`; `go run` bleibt technisch funktionsfähig, ist aber nicht mehr die zitierte Startform |
| 1.27 | 2026-09-18 | C#-/Kotlin-Startform auf echte Start-Make-Ziele umgestellt (`ADR-0098` Festlegung 2, slice-beispiele-csharp-kotlin-start-target): alle acht C#-/Kotlin-Zeilen der vier `**Beispiele:**`-Blöcke zitieren jetzt `make example-run-csharp SURFACE=<oberfläche>`/`make example-run-kotlin SURFACE=<oberfläche>` statt des rohen `docker run --rm -e … <Image> …`-Aufrufs; beide Ziele bauen nichts, sie starten den bereits von `make examples-csharp`/`make examples-kotlin` gebauten Image-Tag |
| 1.28 | 2026-09-18 | Dritter, vollinhaltstragender NATS-Zustellweg ergänzt (`ADR-0100`, `LH-FA-SST-008`, slice-nats-drittstream-core): neuer §4-Abschnitt „Zugriff über den NATS-Vollinhalts-Stream" (Subjekt-Namensraum `cdc.stream.<...>`, dasselbe Nachrichtenschema wie SSE, Zwei-Bedingungen-Aktivierung, serverweite Auth-Nebenwirkung auf das Wecksignal); §5 trägt die neue Variable `CDC_NATS_STREAM_TOKEN` und die auf sieben Schlüssel (drei DSN, drei Token, `nats_url`) gewachsene Zugangsdaten-Klasse der Konfigurationsdatei (`nats_stream_token` neu) |
| 1.29 | 2026-09-18 | Erster Go-Client für den NATS-Vollinhalts-Stream ergänzt (`ADR-0100`, `LH-FA-SST-008`, slice-nats-drittstream-example-go): §4 „Zugriff über den NATS-Vollinhalts-Stream" — der Platzhalter-Absatz wird zu einem `**Beispiele:**`-Block (zunächst nur Go); `examples/nats-stream-client` abonniert `cdc.stream.>` und gibt jede empfangene Change aus, über einen Container-Aufruf (`make example-run-go SURFACE=nats-stream`, Image-Tag `pg-change-feed-examples:go-nats-stream`) |
| 1.30 | 2026-09-18 | C#- und Kotlin-Client für den NATS-Vollinhalts-Stream ergänzt (`ADR-0100`, `LH-FA-SST-008`, slice-nats-drittstream-example-csharp-kotlin): §4 „Zugriff über den NATS-Vollinhalts-Stream" — `**Beispiele:**`-Block komplettiert; `examples/csharp/nats-stream-client` und `examples/kotlin/nats-stream-client` abonnieren `cdc.stream.>` und geben jede empfangene Change aus, über einen Container-Aufruf (`make example-run-csharp`/`make example-run-kotlin SURFACE=nats-stream`, Image-Tags `pg-change-feed-examples:csharp-nats-stream`/`:kotlin-nats-stream`). Mit dieser Zeile ist die volle Matrix (fünf Zugriffs-Oberflächen × drei Sprachen, fünfzehn Programme) im Handbuch vollständig |
| 1.31 | 2026-09-19 | `cdc_changes_pending`/`cdc_errors_total`-Metriken nachgetragen (`LH-QA-OPS-003`, slice-e2e-drei-rtm-luecken): §4 „Metriken lesen" — beide Kennzahlen existierten in `cdc.metrics` bereits seit diesem Slice, waren aber nicht im Handbuch-Text genannt |
| 1.32 | 2026-09-19 | Kopf-Feld `Software-Version` korrigiert: trug seit Ersteinführung unverändert `0.2.0-verdrahtung`, nie mit dem später eingeführten `docs/user/version.md` (`ADR-0051`, welle-release-pipeline-adr-0051) synchronisiert und nirgends sonst referenziert — auf einen Verweis auf die tatsächliche Versionsquelle umgestellt |
| 1.33 | 2026-09-19 | C#-SDK-Hinweis für die HTTP-Oberfläche ergänzt (`LH-FA-SST-009`, `ADR-0106`, slice-sdk-csharp-http-client-flaeche): §4 „Zugriff über die HTTP-/JSON-API" trägt jetzt einen `**SDK:**`-Absatz nach dem `**Beispiele:**`-Block — das NuGet-Package `PgChangeFeed.Client` deckt alle zehn Fähigkeiten dieser Oberfläche (neun aus `SPEC-018` plus `GET /changes`, `SPEC-022`) mit typisierten Requests/Responses und einer typisierten Fehlerklasse ab |
| 1.34 | 2026-09-19 | C#-SDK-Hinweis für den gRPC-Change-Stream ergänzt (`LH-FA-SST-009`, `ADR-0106`, slice-sdk-csharp-grpc-client-flaeche): §4 „Zugriff über den gRPC-Change-Stream" trägt jetzt einen `**SDK:**`-Absatz nach dem `**Beispiele:**`-Block — `PgChangeFeedGrpcClient.StreamChangesAsync` öffnet `ChangeStream/StreamChanges` und liefert ein `IAsyncEnumerable<Change>` mit allen zehn Feldern der Tabelle oben |
| 1.35 | 2026-09-19 | Python-SDK-Hinweis für die HTTP-Oberfläche ergänzt (`LH-FA-SST-009`, `ADR-0107`, slice-sdk-python-http-client-flaeche): §4 „Zugriff über die HTTP-/JSON-API" trägt im `**SDK:**`-Absatz jetzt zusätzlich das PyPI-Package `pgchangefeed` — `PgChangeFeedHttpClient` deckt dieselben zehn Fähigkeiten dieser Oberfläche (neun aus `SPEC-018` plus `GET /changes`, `SPEC-022`) mit typisierten Requests/Responses und einer typisierten Fehlerklasse ab |
| 1.36 | 2026-09-20 | Kotlin-SDK-Hinweis für die HTTP-Oberfläche ergänzt (`LH-FA-SST-009`, `ADR-0109`, slice-sdk-kotlin-http-client-flaeche): §4 „Zugriff über die HTTP-/JSON-API" trägt im `**SDK:**`-Absatz jetzt zusätzlich das GitHub-Packages-Gradle-/Maven-Package `pgchangefeed-kotlin` — `PgChangeFeedHttpClient` deckt dieselben zehn Fähigkeiten dieser Oberfläche (neun aus `SPEC-018` plus `GET /changes`, `SPEC-022`) mit typisierten Requests/Responses und einer versiegelten Fehlerklassen-Hierarchie ab, samt explizitem Hinweis auf die PAT-Pflicht (`read:packages`) beim Bezug über GitHub Packages (`ADR-0109` Festlegung 2) |
| 1.37 | 2026-09-20 | Kotlin-SDK-Hinweis für den gRPC-Change-Stream ergänzt (`LH-FA-SST-009`, `ADR-0109`, slice-sdk-kotlin-grpc-client-flaeche): §4 „Zugriff über den gRPC-Change-Stream" trägt jetzt einen zweiten Absatz im `**SDK:**`-Block — `PgChangeFeedGrpcClient.streamChanges()` öffnet `ChangeStream/StreamChanges` und liefert ein `kotlinx.coroutines.flow.Flow<Change>` mit allen zehn Feldern der Tabelle oben, samt erneutem Hinweis auf die PAT-Pflicht (`read:packages`) beim Bezug über GitHub Packages (`ADR-0109` Festlegung 2); der vorangehende HTTP-Absatz oben wurde korrigiert — er behauptete fälschlich, der gRPC-Change-Stream „folge in einem Folge-Release" |
| 1.38 | 2026-09-22 | C#-SDK-Hinweis für SSE und den NATS-Vollinhalts-Stream ergänzt (`LH-FA-SST-009`, `ADR-0106`, `welle-sdk-csharp-vollabdeckung`, slice-sdk-csharp-sse-client-flaeche + slice-sdk-csharp-nats-stream-client-flaeche): §4 „Zugriff über Server-Sent-Events" und §4 „Zugriff über den NATS-Vollinhalts-Stream" tragen jetzt je einen `**SDK:**`-Absatz — `PgChangeFeedSseClient.StreamChangesAsync` bzw. `PgChangeFeedNatsStreamClient.StreamChangesAsync` liefern ein `IAsyncEnumerable<Change>` mit allen zehn Feldern der jeweiligen Tabelle; das NuGet-Package `PgChangeFeed.Client` ist dafür auf `0.2.0` gehoben |
| 1.39 | 2026-09-22 | Kotlin-SDK-Hinweis für SSE und den NATS-Vollinhalts-Stream ergänzt (`LH-FA-SST-009`, `ADR-0109`, `welle-sdk-kotlin-vollabdeckung`, slice-sdk-kotlin-sse-client-flaeche + slice-sdk-kotlin-nats-stream-client-flaeche): §4 „Zugriff über Server-Sent-Events" und §4 „Zugriff über den NATS-Vollinhalts-Stream" tragen jetzt je einen `**SDK:**`-Absatz — `PgChangeFeedSseClient.streamChanges()` bzw. `PgChangeFeedNatsStreamClient.streamChanges()` liefern eine `Sequence<Change>` mit allen zehn Feldern der jeweiligen Tabelle, samt erneutem Hinweis auf die PAT-Pflicht (`read:packages`) beim Bezug über GitHub Packages; die vorangehenden HTTP-/gRPC-Absätze oben wurden korrigiert — sie behaupteten fälschlich, SSE und der NATS-Vollinhalts-Stream blieben für dieses Package „vorerst außerhalb"; das GitHub-Packages-Gradle-/Maven-Package `pgchangefeed-kotlin` ist dafür auf `0.2.0` gehoben |
| 1.40 | 2026-09-23 | Python-SDK-Hinweis für den gRPC-Change-Stream ergänzt (`LH-FA-SST-009`, `ADR-0110`, `welle-sdk-python-vollabdeckung`, slice-sdk-python-grpc-client-flaeche): der §4-Absatz zum PyPI-Package `pgchangefeed` trägt jetzt die gRPC-Stream-Fläche `PgChangeFeedGrpcClient` (siehe „Zugriff über den gRPC-Change-Stream") statt „gRPC, SSE und der NATS-Vollinhalts-Stream bleiben vorerst außerhalb"; SSE und der NATS-Vollinhalts-Stream folgen im selben Folge-Release |
| 1.41 | 2026-09-23 | Python-SDK-Absatz im gRPC-Handbuch-Abschnitt ergänzt (`LH-FA-SST-009`, `ADR-0110`, slice-sdk-python-grpc-client-flaeche Fixrunde): „Zugriff über den gRPC-Change-Stream" trägt jetzt den `**SDK:**`-Absatz des PyPI-Packages — `PgChangeFeedGrpcClient.stream_changes()` liefert einen Iterator über die generierten `Change`-Nachrichten mit allen zehn Feldern (dritte Sprache neben C#/Kotlin im selben Abschnitt), samt `timeout`-Form; `Stand:`-Datum auf diesen Zug gezogen |
| 1.42 | 2026-09-23 | Python-SDK-Absatz im SSE-Handbuch-Abschnitt ergänzt (`LH-FA-SST-009`, `ADR-0110`, `welle-sdk-python-vollabdeckung`, slice-sdk-python-sse-client-flaeche): „Zugriff über Server-Sent-Events" trägt jetzt den `**SDK:**`-Absatz des PyPI-Packages — `PgChangeFeedSseClient.stream_changes()` liefert einen Iterator über die getypten `StreamChange`-Events mit allen zehn Feldern (dritte Sprache neben C#/Kotlin im selben Abschnitt); der NATS-Vollinhalts-Stream folgt im selben Folge-Release |
| 1.43 | 2026-09-23 | Python-SDK-Absatz für den NATS-Vollinhalts-Stream ergänzt (`LH-FA-SST-009`, `ADR-0110`, `welle-sdk-python-vollabdeckung`, slice-sdk-python-nats-stream-client-flaeche): §4 „Zugriff über den NATS-Vollinhalts-Stream" trägt jetzt den dritten Sprach-`**SDK:**`-Absatz — `PgChangeFeedNatsStreamClient.stream_changes()` abonniert `cdc.stream.<source_id>.>` und liefert einen Iterator über die getypten `StreamChange`-Events mit allen zehn Feldern; das PyPI-Package `pgchangefeed` ist dafür auf `0.2.0` gehoben — die volle Vier-Wege-Matrix ist damit für alle drei SDK-Sprachen im Handbuch vollständig |
| 1.44 | 2026-09-24 | Feld `origin` in den Lesewegen ergänzt (`LH-FA-CAP-009`, `LH-FA-DAT-006`, `ADR-0111`, slice-backfill-change-origin): §4 „Änderungen lesen" trägt `origin` als letzte Spalte des SQL-Beispiels über `cdc.changes` samt Bedeutung (`wal` \| `backfill`, ein fehlender Wert liest als `wal`), §4 „Zugriff über die HTTP-/JSON-API" nennt `origin` als letztes Feld der `GET /changes`-Antwort; die drei Live-Zustellwege tragen das Feld nicht |
| 1.45 | 2026-09-24 | Schema-Upgrade über eine View-Signaturänderung dokumentiert (`LH-QA-OPS-005`, `ADR-0114`, slice-backfill-change-origin Fixrunde): §4 „Schema aktualisieren" nennt die Reihenfolge (Schema-Rollout vor dem Container-Tausch), den automatischen Vorlauf `DROP VIEW cdc.<name>` samt Meldung, das Lesefenster für SQL-Leser (Richtwert aus einer einzelnen Messung, nicht garantiert) und das Verhalten bei einem abhängigen Objekt oder einem Abbruch nach dem Vorlauf |
| 1.46 | 2026-09-24 | Hinweis zu Rechten und Vorbedingung des View-Signatur-Vorlaufs ergänzt (`LH-QA-OPS-005`, `ADR-0114`, slice-backfill-change-origin Fixrunde): §4 „Schema aktualisieren" benennt, dass `DROP VIEW` die Rechteliste der View verwirft und nur `cdc_reader` im selben Lauf sein `SELECT`-Recht zurückbekommt (eigene Grants an andere Rollen setzt der Betreiber erneut), sowie die feste Adressierung des Schemas `cdc` |
| 1.47 | 2026-09-24 | Rechteschnitt von `cdc_admin` ergänzt (`LH-QA-SEC-001`, `LH-QA-SEC-002`, `ADR-0047`, `ADR-0050`, slice-backfill-run-store Fixrunde): §2 „Zugriff und Rollen" nennt die Verarbeitung der Antrags-Queue `cdc.administration_request` (lesen, Ausgang vermerken) als Zweck der Rolle, §5 „Umgebungsvariablen des Feed-Containers" die Zeile `CDC_ADMIN_DSN`; §4 „Schema aktualisieren" trägt den Absatz „Rechte der drei Rollen" (der Rollout setzt die Rechte bei jedem Lauf, Schema-Rollout vor dem Container-Tausch, Anträge bleiben ohne das Recht `pending`) |
| 1.48 | 2026-09-24 | SQL-Auslösung des Backfills dokumentiert (`LH-FA-CAP-009`, `LH-FA-ADM-001`, `LH-FA-SST-003`, `ADR-0111`, `ADR-0113`, `ADR-0116`, slice-backfill-sql-administration): §4 neuer Abschnitt „Bestand als Backfill überführen" (`cdc.backfill_table`, `applied` heißt „angenommen", View `cdc.backfill_status`, geschätzte Zeilenzahl, Neustart-Verhalten, Sichtbarkeits-Grenze, Bedeutung der Schema-Version einer Backfill-Änderung, Betriebs-Vorbedingungen); §4 „Änderungen lesen" trägt die Regel „Position und `limit`" mit dem Schlüsselvergleich, „Diagnose ausführen" den Abschnitt „Backfill je Tabelle", „Schema aktualisieren" die Rechte von `cdc_capture`/`cdc_reader`; §2 Rollen und Betriebs-Hinweis zum `SELECT`-Recht, §5 die beiden DSN-Zeilen, §6 Fehlerklassen, §8 Glossar, §9 Grenzwerte |
| 1.49 | 2026-09-24 | Backfill-Abschnitt an Rollen und Stichtag angeglichen (`LH-FA-CAP-009`, `ADR-0111`, `ADR-0113`, slice-backfill-sql-administration Fixrunde): §4 „Bestand als Backfill überführen" nennt den Antrag und dessen Vermerk unter `cdc_admin`, das Lesen von `cdc.backfill_status` unter einer `cdc_reader`-Identität (die Rolle `cdc_admin` trägt kein `SELECT` auf die View), und den Snapshot als Start des Runs statt des Antrags (auch §8 Glossar); ein Backfill-Antrag einer anderen Quelle bleibt für die Instanz dieser Quelle `pending` |
| 1.50 | 2026-09-24 | Gemessene Startposition eines frisch registrierten Consumers dokumentiert (`LH-FA-CAP-009`, `LH-FA-CON-005`, `ADR-0111`, slice-backfill-e2e): §4 „Bestand als Backfill überführen" trägt den Punkt „Startposition eines neuen Consumers" (`offset` 0, `acknowledged` `false`, vor der Snapshot-Position jedes Runs; Ursprung: der Lauf von `make test-integration`) und nennt in der Zustandstabelle, dass `rows_copied` eines `interrupted`-Runs den zuletzt festgehaltenen Fortschritt trägt |
| 1.51 | 2026-09-24 | Sperre des Backfill-Runs und Ausgang bei umgeschriebener Tabelle dokumentiert (`LH-FA-CAP-009`, `ADR-0111`, `ADR-0118`, slice-backfill-e2e Fixrunde): §4 „Bestand als Backfill überführen“ trägt den Absatz „Sperre der Tabelle“ (Lesesperre bis zum Ende des Runs, Wirkung auf DDL mit `ACCESS EXCLUSIVE`) und den Punkt „Umschreiben der Tabelle im Fenster“ (`failed`/`transient` ohne Änderung, neuer Antrag als Abhilfe, Fehlalarme `VACUUM FULL`/`CLUSTER`) |
| 1.52 | 2026-09-24 | Wirkung der Tabellensperre des Backfill-Runs vollständig und mit Ursprung dokumentiert (`LH-FA-CAP-009`, `ADR-0111`, `ADR-0118`, slice-backfill-e2e Fixrunde): §4 „Bestand als Backfill überführen“, Absatz „Sperre der Tabelle“, nennt neben Lesern auch Schreiber und die Publication-Abfrage der Administration als hinter einer wartenden DDL gestaut (gemessen, PostgreSQL 18) und die Gegenrichtung (der Run wartet ohne eigene Zeitgrenze auf eine offene `ACCESS EXCLUSIVE`-Transaktion); `RENAME COLUMN` im Fenster ist als im Review gemessen, nicht im E2E-Runner belegt gekennzeichnet |
| 1.53 | 2026-09-25 | Warnungen des Backfill-Runs und gemessene Richtgrößen dokumentiert (`LH-FA-CAP-009`, `ADR-0111`, `ADR-0113`, slice-backfill-bench-richtgroesse): §4 „Bestand als Backfill überführen“ nennt, wann `warn_estimated_size` und `warn_duration` gesetzt werden, und den WAL-Rückstand des Capture-Slots als vierte Betriebs-Vorbedingung; „Diagnose ausführen“ deutet die beiden Kennzeichnungen; §9 „Grenzwerte“ trägt die Toleranz der Kopierdauer (Startwert, Setzung ohne Messung), die Richtgröße (abgeleitet, Orientierung, keine Grenze) und die gemessenen Werte (Kopierdauer, Speicher, WAL-Rückstand, Wirkung auf die Live-Erfassung, Schätzung der Zeilenzahl) mit Host und Lauf |
| 1.54 | 2026-09-25 | Ursache und Abhilfen des WAL-Rückstands sowie Ursprung der Messwerte nachgezogen (`LH-FA-CAP-009`, `ADR-0111`, `ADR-0113`, `ADR-0120`, slice-backfill-bench-richtgroesse Fixrunde): §4 „Bestand als Backfill überführen“ nennt den WAL-Rückstand als nicht an den Backfill gebunden, den Folge-Slice `slice-backfill-slot-leerlauf-bestaetigung` und die Betriebs-Abhilfen (Commit auf einer aktivierten Tabelle, Datei-Feld `wal_retention_error_bytes`); §9 „Grenzwerte“ nennt zur Richtgröße die Werte dreier Läufe (7.693, 8.933, 8.559 Zeilen/s; Konstante = kleinster Wert), zum Speicher des Feed-Containers die Spitze im Run **und** die Probe 20 s nach dem Run (641,7 MiB als höchster gemessener Wert), zur Live-Wirkung drei Einzelläufe statt einer Aussage „kein Unterschied“, und kennzeichnet die im Repository nicht auflösbaren Läufe als übernommen |
| 1.55 | 2026-09-25 | Herkunft der Zahlen des Laufs `20260925T012459Z` in §9 „Grenzwerte“ als übernommen aus dem Lauf-Bericht des Implementers gekennzeichnet (im Repository nicht auflösbar), Richtgröße um den gedruckten Lauf `20260925T015600Z` des Verifikations-Reports ergänzt (`LH-FA-CAP-009`, `ADR-0111`, `ADR-0113`, slice-backfill-bench-richtgroesse Closure) |
| 1.56 | 2026-09-25 | Bestätigung von WAL ohne Inhalt für die Publication im Leerlauf des Streams dokumentiert (`LH-FA-CAP-009`, `LH-QA-REL-001`, `ADR-0120`, slice-backfill-slot-leerlauf-bestaetigung): §4 „WAL-Rückstand prüfen“ nennt die Bedeutung von `cdc_wal_retention_bytes` (vom Feed noch nicht bestätigtes WAL) und dass WAL ohne Inhalt für die Publication den Wert nicht wachsen lässt; §4 „Bestand als Backfill überführen“ trägt den WAL-Rückstand als Punkt ohne Abbruch über die Fehlerschwelle und die offene Schreibtransaktion des Runs als verbleibende Last; §9 „Grenzwerte“ führt den Rückstand mit Bestätigung (Lauf `20260925T032925Z`), die Messwerte ohne Bestätigung mit ihrem Lauf, das gehaltene WAL und den Spill als Grenze der Ein-Transaktions-Form und die Richtgröße um den Lauf `20260925T032925Z` ergänzt; der Satz „Diese Schwelle kann bei weniger Zeilen greifen als die Richtgröße“ entfällt |
| 1.57 | 2026-09-25 | Herkunft der Zahlen des Laufs `20260925T032925Z` in §9 „Grenzwerte“ als übernommen aus dem Lauf-Bericht des Implementers gekennzeichnet (im Repository nicht auflösbar); Nachmessung des Laufs `20260925T043056Z` aus dem Review-Report ergänzt (Rückstand 0 MiB in neun Runs, gehaltenes WAL 141 MiB bei 200.000 Zeilen, Richtgröße 2.000.000 bei 4.504 Zeilen/s); Spanne der Richtgröße über sechs Läufe 2.000.000 bis 5.000.000 (`LH-FA-CAP-009`, `ADR-0111`, `ADR-0113`, `ADR-0120`, slice-backfill-slot-leerlauf-bestaetigung Fixrunde) |
| 1.58 | 2026-09-25 | Feld `origin` der über `GET /changes` gelesenen Änderungen in den SDK-Absätzen von §4 „Zugriff über die HTTP-/JSON-API“ ergänzt: die drei Packages tragen es, eine Antwort ohne das Feld liest als `wal`, die Live-Wege tragen es nicht (`LH-FA-SST-009`, `LH-FA-SST-006`, `ADR-0111`, slice-backfill-sdk-origin) |
| 1.59 | 2026-09-25 | Regel für `origin` in den SDK-Absätzen von §4 „Zugriff über die HTTP-/JSON-API“ präzisiert: fehlendes Feld oder JSON-`null` liest als `wal`, jeder andere Server-Wert (auch leer oder unbekannt) kommt unverändert an; der Python-Absatz führt SSE und den NATS-Vollinhalts-Stream als vom Package getragen statt als folgend (`LH-FA-SST-009`, `LH-FA-SST-006`, `ADR-0111`, slice-backfill-sdk-origin Fixrunde) |
| 1.60 | 2026-09-25 | Speicher des Feed-Containers im Backfill gemessen und die Ursache benannt (`LH-FA-CAP-009`, `ADR-0111`, `ADR-0113`, slice-backfill-speicher-untersuchung): §9 „Grenzwerte“ ersetzt die übernommenen Zahlen durch gemessene — die Spitze hängt an der Zahl der Changes in `cdc.change` (Bereinigungslauf liest je Takt alle Changes), nicht an der Tabellengröße; 1,03 bis 1,59 KiB je Change bei schmalen, 2,65 bis 4,19 KiB bei breiten Zeilen; Bemessung des Speicherlimits; die Richtgröße um den Speicher ergänzt; §4 „Bestand als Backfill überführen“ nennt den Speicher des Feed-Containers, „Aufbewahrung (Retention)“ die Lesung aller Changes je Durchlauf |
| 1.61 | 2026-09-25 | Zahlen und Herkunftsangaben der Speicher-Messung in §9 „Grenzwerte“ und §4 „Bestand als Backfill überführen“ nachgezogen (`LH-FA-CAP-009`, `ADR-0111`, `ADR-0113`, slice-backfill-speicher-untersuchung Fixrunde): Höchstwert je Change 1,59 statt 1,57 KiB, `GOGC=25` senkt die Spitze um 11 bis 21 % statt 18 bis 21 %, die Runs mit 2.000.000 Changes vor dem Run stehen in einer eigenen Zeile; der Speicher im Run bei nicht leerem `cdc.change` ist benannt; Herkunft der Aussage zur laufenden Erfassung; die Server-Versionen `v0.1.0` bis `v0.1.2` tragen denselben Bereinigungslauf; Reihe N als Gegenprobe zum Ausbleiben der Bereinigungs-Takte |
| 1.62 | 2026-09-25 | Bereinigungslauf liest Kandidaten seitenweise ohne Row Images (`LH-FA-RET-004`, `LH-FA-CAP-009`, `ADR-0124`, slice-retention-lauf-speicher-begrenzung): §4 „Aufbewahrung (Retention)“ beschreibt die Seiten (10.000 Changes, nicht atomar, Consumer-Positionen einmal je Durchlauf), §4 „Bestand als Backfill überführen“ und §9 „Grenzwerte“ ersetzen die Bemessung des Speicherlimits je Change durch die Nachmessung (Spitze 14,9 bis 17,6 MiB bei 1.000.000 bis 3.000.000 Changes, zwei Läufe) und nennen den Stand der Vorversionen `v0.1.0` bis `v0.1.2`; die Richtgröße bleibt und folgt der Kopierdauer |
