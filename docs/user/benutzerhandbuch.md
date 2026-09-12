# Benutzerhandbuch: PG Change Feed

Version: 1.1
Software-Version: 0.2.0-verdrahtung
Stand: 2026-09-12

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
| `cdc_admin` | Verwaltungszugriff (Registrierung von Quellen und Tabellen, Heartbeat, `register-consumer`/`acknowledge-consumer`) | `CDC_ADMIN_DSN` |
| `cdc_reader` | Nur-Lese-Zugriff auf die Diagnose- und Lese-Views (`cdc.active_tables`, `cdc.consumer_status`, `cdc.changes`, `cdc.metrics`, `cdc.heartbeat`) | `CDC_READER_DSN` |

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
| `CDC_READER_DSN` | ja | Verbindung über die Rolle `cdc_reader` (`--healthcheck`) |
| `CDC_SOURCE_ID` | ja | Kennung der Quelle (muss in `cdc.source` registriert sein) |
| `CDC_PUBLICATION` | ja | Name der PostgreSQL-Publication |
| `CDC_SLOT` | ja | Name des Logical-Replication-Slots |
| `CDC_TABLES` | ja | Aktivierte Tabellen, Format `schema.tabelle=tabelle-id:schema-version-id`, kommagetrennt |
| `CDC_LOG_LEVEL` | nein | Log-Level des strukturierten JSON-Loggers (Default `info`) |

Fehlt eine Pflichtvariable oder ist `CDC_TABLES` leer, startet der
Container nicht (Fehlerklasse `configuration`).

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
| `replication` | eine Störung der Stream-Ordnung (z. B. Commit ohne offene Transaktion) | Sichtbarer Fehler |
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
- Ein Container-Lauf bindet genau eine Quelle.

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
