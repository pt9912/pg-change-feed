# Benutzerhandbuch: PG Change Feed

Version: 1.99
Software-Version: siehe `docs/user/version.md`
Stand: 2026-10-04

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
| `cdc_admin` | Verwaltungszugriff (Registrierung von Quellen und Tabellen, Heartbeat, `register-consumer`/`acknowledge-consumer`, Retention-Löschausführung, Verarbeitung der Antrags-Queue `cdc.administration_request` — offene Anträge lesen und ihren Ausgang vermerken, ohne Anträge selbst anzulegen oder zu löschen —, Annahme eines Backfill-Antrags, die die Run-Zeile in `cdc.backfill_run` anlegt; `cdc.backfill_table`, `cdc.set_transformation`/`cdc.remove_transformation` (siehe [Transformationsregel konfigurieren](#transformationsregel-konfigurieren)), `cdc.set_route`/`cdc.remove_route` (siehe [Routing-Regel konfigurieren](#routing-regel-konfigurieren)) und die übrigen Antragsfunktionen ruft nur diese Rolle auf) | `CDC_ADMIN_DSN` |
| `cdc_reader` | Nur-Lese-Zugriff auf die Diagnose- und Lese-Views (`cdc.active_tables`, `cdc.consumer_status`, `cdc.changes`, `cdc.metrics`, `cdc.heartbeat`, `cdc.retention_blockers`, `cdc.backfill_status`) — trägt auch `--healthcheck`, `diagnose` (siehe [Diagnose ausführen](#diagnose-ausführen)) und das Lesen für den OTLP-Export (siehe [Metriken per OTLP übertragen](#metriken-per-otlp-übertragen)) | `CDC_READER_DSN` |

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
sind die einzige Herkunft des Standes einer Tabelle; jeder Pfad,
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

### Transformationsregel konfigurieren

Eine Transformationsregel benennt eine Spalte im Row Image um oder bildet
ihren Wert nach einer festen Zuordnung ab, bevor die Change gespeichert
wird.

**Voraussetzung:** eine Login-Identität mit `cdc_admin`-Mitgliedschaft,
verbunden über `CDC_ADMIN_DSN` (siehe [Zugriff und Rollen](#zugriff-und-rollen));
die physische Tabelle existiert an der Quelle.

**Vorgehen:**

1. Wählen Sie einen Regeltyp:

   | `kind` | Pflichtschlüssel | Wirkung |
   |---|---|---|
   | `rename_column` | `column`, `to` | benennt den Schlüssel `column` im Image auf `to` um; der Wert bleibt unverändert. |
   | `map_value` | `column`, `values` (nicht leeres Objekt) | ist der Wert von `column` ein Schlüssel von `values`, steht der zugeordnete Wert im Image; jeder andere Wert bleibt unverändert. |

2. Rufen Sie die Regel auf — `rule_spec` als `json`, **nicht** als `jsonb`:

   ```sql
   SELECT cdc.set_transformation('<source_id>', '<schema>', '<tabelle>', '<regelname>',
     '{"kind": "rename_column", "column": "name", "to": "customer_name"}');
   ```

   Beispiel `map_value`: `{"kind": "map_value", "column": "status", "values": {"o": "open", "c": "closed"}}`
   bildet den Wert `o` auf `open` ab; ein anderer Wert (z. B. `x`) bleibt
   unverändert. Übergeben Sie `rule_spec` als Literal oder `::json`-Wert —
   ein `::jsonb`-Wert wird abgelehnt (`function cdc.set_transformation(unknown,
   unknown, unknown, unknown, jsonb) does not exist`); casten Sie in diesem
   Fall mit `::json`.

3. Prüfen Sie den Fortschritt über die zurückgegebene Antrags-Kennung:

   ```sql
   SELECT status, error_message
   FROM cdc.administration_request
   WHERE administration_request_id = '<zurückgegebene-id>';
   ```

**Ergebnis:** `status` wechselt von `pending` zu `applied` (Erfolg) oder
`failed` (Fehlertext in `error_message`). Ab `applied` trägt jede danach
erfasste Änderung dieser Tabelle die transformierte Form; eine zuvor
gespeicherte Change behält ihre bisherige Form (die Regel wirkt **nicht
rückwirkend**). Ein Wert, der im Image fehlt (`NULL`, unverändertes TOAST,
ausgeschlossene oder generierte Spalte), bleibt für beide Regeltypen
abwesend.

**Regel wieder entfernen:** `SELECT cdc.remove_transformation('<source_id>',
'<schema>', '<tabelle>', '<regelname>')` — derselbe Antrags-Weg und
Status-Poll wie beim Setzen. Um eine Regel zu ersetzen: erst entfernen,
dann neu setzen; beide Aufrufe dürfen in derselben Transaktion stehen.

**Fehler: Antrag endet `failed` mit einem Konfliktfehler**

**Ursache:** einer der folgenden Konflikte — Meldungscode und Klartext stehen in
`error_message` (siehe [Meldungscodes](#meldungscodes)):

| Ursache | Fehlertext |
|---|---|
| Regelname bereits vergeben (je Tabelle eindeutig) | `abgelehnt [PCF-E8020]: Regelname bereits vergeben` |
| Spalte trägt bereits eine andere Regel | `abgelehnt [PCF-E8021]: Spalte trägt bereits eine Regel` |
| Zielname (`to`) kollidiert mit dem `to` einer anderen Regel | `abgelehnt [PCF-E8022]: Zielname kollidiert mit einer anderen Regel` |
| Zielname kollidiert mit einer Spalte der Tabelle (auch mit ausgeschlossenen) | `abgelehnt [PCF-E8022]: Zielname kollidiert mit einer Spalte der Tabelle` |
| Spalte existiert nicht an der Quelle | `abgelehnt [PCF-E8024]: Spalte existiert nicht an der Quelle` |
| `remove_transformation` gegen einen unbekannten Regelnamen | `abgelehnt [PCF-E8023]: Regelname nicht geführt` |

Ein Tippfehler beim Aufruf selbst (leeres Schema, leerer Tabellenname, leere
Spalte) endet ebenso `failed`, ohne nachfolgende Anträge zu blockieren.

**Lösung:** den betroffenen Namen ändern, oder — bei „Spalte trägt bereits
eine Regel"/„Regelname bereits vergeben" — die vorhandene Regel zuerst
entfernen (siehe „Regel wieder entfernen" oben).

**Fehler: Erfassung endet mit der Fehlerklasse `schema`**

**Ursache:** Eine Regel ist auf eine Change nicht anwendbar — ihr Zielname
kollidiert mit einer Spalte (etwa nach einer Tabellen-Erweiterung), oder ihre
`column` fehlt in der Relation der Change. Die Erfassung der **gesamten Quelle** endet dann sichtbar mit
dieser Fehlerklasse (siehe [Fehlerklassen](#fehlerklassen)) — kein
Datenverlust. Wurde die Spalte `column` an der Quelle entfernt, während die
Spaltenform der Tabelle bekannt ist, endet die Erfassung schon an der
inkompatiblen Schemaänderung, nicht an der Regel: das
Entfernen der Regel genügt dort nicht, und die Regel auf die entfernte Spalte
ist zusätzlich zu entfernen (*abgeleitet*, nicht am laufenden System
ausprobiert; dieselbe
Lage wie bei Routing-Regeln, siehe
[Routing-Regel konfigurieren](#routing-regel-konfigurieren)).

**Lösung:** Bei einer nicht anwendbaren Regel: Regel entfernen oder ersetzen,
Prozess neu starten; die zuvor nicht bestätigte Transaktion erscheint danach
über `cdc.changes` in Rohform. Im Backfill gilt dieselbe Prüfung — der
betroffene Run endet `failed`; hier ist die Lösung ein **neuer**
`cdc.backfill_table`-Antrag, ohne Prozessneustart.

**Hinweise:**

- **`map_value` ist nicht umkehrbar:** Bilden mehrere Quellwerte auf
  denselben Zielwert ab, meldet das System nichts — kein Fehler, keine
  Warnung. Der Vergleich ist zeichengenau (Groß-/Kleinschreibung zählt).
  Bei einer Regel mit sehr vielen Paaren (ab etwa 10.000) beobachten Sie
  `cdc_capture_lag` (siehe [Metriken lesen](#metriken-lesen)).
- **Verhältnis zum Spaltenausschluss:** Der Ausschluss gilt zuerst — eine
  ausgeschlossene Spalte wird nie gelesen, auch nicht von einer Regel.
- **Reihenfolge bei mehreren Sitzungen:** Setzen Sie Anträge auf dieselbe
  Regel oder Spalte nicht aus zeitlich überlappenden Transaktionen
  mehrerer Sitzungen ab — sonst kann der laufende Stand bis zum nächsten
  Prozessstart vom zuletzt beantragten Stand abweichen.
- **Dauerhaft über Neustart:** Der Regelstand übersteht einen Neustart des
  Feed-Containers sowie eine Deaktivierung mit anschließender Aktivierung.
  Rollen Sie den Feed-Container niemals auf einen älteren Stand zurück,
  der einen bereits vermerkten Regeltyp nicht kennt (etwa `map_value`
  unter einer Version vor dessen Einführung) — Prozessstart und jeder
  weitere Regel-Antrag der Quelle bleiben sonst angehalten.
- **Start des Feed-Containers:** Ein Antrag, der beim Prozessstart noch in
  Bearbeitung ist (etwa durch eine Tabellensperre eines Betreibers),
  verzögert den Beginn der Erfassung höchstens 30 Sekunden; der
  Healthcheck bleibt währenddessen gesund, und der Antrag wird danach
  regulär weiterverarbeitet.
- Alle Zustellwege (`cdc.changes`, `GET /changes`, gRPC-Stream,
  SSE-Stream, NATS-Vollinhalts-Stream) tragen dieselbe transformierte
  Form.

Ein Beispiel gegen eine laufende Demo-Umgebung liefert
`make example-transformation-demo` (siehe [`examples/README.md`](../../examples/README.md)).

### Routing-Regel konfigurieren

Eine Routing-Regel weist einer erfassten Change ein **Zustellziel** zu: einen
benannten Kanal je Quelle, den ein Leser über den Parameter `target` oder das
NATS-Zusatz-Subjekt auswählt (siehe „Ziel lesen" unten). Die Regel ändert weder
das Row Image noch ein anderes Feld der Change als ihr Ziel-Label
`route_target`; eine Change trägt höchstens ein Ziel, und es gibt kein
Standardziel.

**Voraussetzung:** eine Login-Identität mit `cdc_admin`-Mitgliedschaft,
verbunden über `CDC_ADMIN_DSN` (siehe [Zugriff und Rollen](#zugriff-und-rollen));
die physische Tabelle existiert an der Quelle; der Feed-Container läuft für die
Quelle, denn er verarbeitet den Antrag. Das Lesen der Ergebnisse braucht eine
Login-Identität mit `cdc_reader`-Mitgliedschaft.

**Vorgehen:**

1. Legen Sie die Regel fest. Die `rule_spec` ist ein JSON-Objekt mit diesen
   Schlüsseln; jeder andere Schlüssel, auch innerhalb von `when`, endet den
   Antrag `failed`:

   | Schlüssel | Pflicht | Bedeutung |
   |---|---|---|
   | `target` | ja | Name des Zustellziels, 1 bis 63 Zeichen: das erste aus `a`–`z` und `0`–`9`, die übrigen aus `a`–`z`, `0`–`9`, `_` und `-`. Der Name wird zeichengenau verglichen, nicht gefaltet oder gekürzt. |
   | `order` | ja | Auswertungsreihenfolge: eine positive ganze Zahl bis 2147483647, je Tabelle eindeutig; die kleinere `order` wird zuerst geprüft. |
   | `when` | nein | Objekt mit den beiden Pflichtschlüsseln `column` und `equals` (beide Zeichenketten): die Regel trifft, wenn der Wert der Quellspalte `column` im Bild zeichengenau `equals` ist. Ohne `when` trifft die Regel jede Change der Tabelle (**Abschlussregel**). |

   Der erste Treffer in aufsteigender `order` bestimmt das Ziel, keine weitere
   Regel wird geprüft; trifft keine Regel, ist `route_target` `NULL`. Die
   Bedingung liest den Quellwert vor jeder Transformation (siehe
   [Transformationsregel konfigurieren](#transformationsregel-konfigurieren)).
   Die Obergrenze von `order` ist eine Setzung der Umsetzung, die
   Spezifikation nennt keine. Als `order` wird jede JSON-Zahl mit dem Wert einer
   positiven ganzen Zahl angenommen (`10`, `10.0`, `1e1`); ein Bruchteil, 0, ein
   negativer Wert und ein Wert über der Obergrenze enden den Antrag `failed` mit
   `rule_spec ist ungültig`.

2. Rufen Sie die Regel auf — `rule_spec` als `json`, **nicht** als `jsonb`:

   ```sql
   SELECT cdc.set_route('<source_id>', '<schema>', '<tabelle>', '<regelname>',
     '{"target": "eu", "order": 10, "when": {"column": "region", "equals": "eu"}}');
   ```

   Der Regelname ist je Tabelle unter den Routing-Regeln eindeutig; der
   Namensraum ist von dem der Transformationsregeln getrennt. Ein
   `::jsonb`-Wert wird abgelehnt (`function cdc.set_route(unknown, unknown,
   unknown, unknown, jsonb) does not exist`); casten Sie in diesem Fall mit
   `::json`.

3. Prüfen Sie den Fortschritt über die zurückgegebene Antrags-Kennung:

   ```sql
   SELECT status, error_message
   FROM cdc.administration_request
   WHERE administration_request_id = '<zurückgegebene-id>';
   ```

**Ergebnis:** `status` wechselt von `pending` zu `applied` (Erfolg) oder `failed`
(Fehlertext in `error_message`). Ab `applied` trägt jede danach erfasste Change
dieser Tabelle das Ziel der ersten treffenden Regel, am laufenden Prozess ohne
Neustart. Der Regelstand übersteht einen Neustart des Feed-Containers und gilt
auch für eine Tabelle, die über `EnableTable` des gRPC- oder HTTP-Zugriffswegs
aktiviert wird (*Ursprung:* übernommen aus dem E2E-Lauf von `make
test-integration`: zwei reale `docker restart`, je Aktivierungsweg eine Tabelle
mit Regel mit Ziel `api` und eine ohne Regel mit `NULL`).

**Beispiel** (Quelle `meine-quelle`, Tabelle `public.orders` mit den Spalten
`id`, `name`, `region`, aktiviert; die Zeilen 1 bis 3 sind vor den Regeln
eingefügt). Die Aufrufe von `cdc.set_route` laufen unter einer
`cdc_admin`-Identität:

```sql
SELECT cdc.set_route('meine-quelle', 'public', 'orders', 'eu_orders',
  '{"target": "eu", "order": 10, "when": {"column": "region", "equals": "eu"}}');
SELECT cdc.set_route('meine-quelle', 'public', 'orders', 'rest',
  '{"target": "sonstige", "order": 100}');
```

Beide Anträge enden `applied`. Danach fügt die Quelle drei Zeilen ein
(`INSERT INTO public.orders VALUES (4, 'Di', 'eu'), (5, 'Ed', 'us'), (6, 'Fy', NULL);`);
die Abfrage läuft unter einer `cdc_reader`-Identität:

```sql
SELECT new_data->>'id' AS id, new_data->>'region' AS region, route_target
FROM cdc.changes
WHERE source_id = 'meine-quelle'
ORDER BY commit_position, sequence;
```

```text
 id | region | route_target
----+--------+--------------
 1  | eu     |
 2  | us     |
 3  |        |
 4  | eu     | eu
 5  | us     | sonstige
 6  |        | sonstige
```

Die Zeilen 1 bis 3 sind vor den Regeln erfasst und tragen kein Ziel. Zeile 4
trifft `eu_orders`; Zeile 5 trifft keine Bedingung, die Abschlussregel `rest`
bestimmt `sonstige`; bei Zeile 6 fehlt der Wert von `region` im Bild, die
Bedingung trifft nicht, und `rest` bestimmt ebenfalls `sonstige`. Ohne die
Regel `rest` bliebe `route_target` bei den Zeilen 5 und 6 `NULL`. Ein Ziel wählt
`WHERE route_target = 'eu'` (liefert die Zeile 4), die Changes ohne Ziel
`WHERE route_target IS NULL`. *Ursprung:* gemessen in einer Compose-Umgebung
(PostgreSQL 18.6, Image `:dev`); die Aufrufe liefen unter Login-Identitäten mit
`cdc_admin`- bzw. `cdc_reader`-Mitgliedschaft, die gedruckten Zeilen entsprechen
der Ausgabe oben.

**Regel wieder entfernen:** `SELECT cdc.remove_route('<source_id>', '<schema>',
'<tabelle>', '<regelname>')` — derselbe Antrags-Weg und Status-Poll wie beim
Setzen. Um eine Regel zu ersetzen: erst entfernen, dann neu setzen.

**Fehler: Antrag endet `failed`**

**Ursache:** eine Prüfung der Regel verletzt — `error_message` beginnt mit
`abgelehnt [<Code>]: ` (siehe [Meldungscodes](#meldungscodes)); es folgt der
Klartext, ein Doppelpunkt, ein Leerzeichen und die
Adresse (`schema.tabelle.regelname`, `schema.tabelle.spalte`,
`schema.tabelle.zielname`, bei `order` `schema.tabelle.<Wert>`; bei einem
unbekannten Schlüssel der bloße Schlüsselname ohne Tabelle). Die führende
Stelle für Wortlaut und Reihenfolge ist die Fehlertext-Tabelle unten; die Konfliktfreiheit
schließt Mehrdeutigkeit aus, statt sie aufzulösen:

| Prüfung | Fehlertext (Beispiel mit Adresse) |
|---|---|
| R1 — Regelname je Tabelle eindeutig | `abgelehnt [PCF-E8020]: Regelname bereits vergeben: public.orders.eu_orders` |
| R2 — `order` je Tabelle eindeutig | `abgelehnt [PCF-E8030]: order bereits vergeben: public.orders.10` |
| R3 — `when.column` existiert an der Quelle | `abgelehnt [PCF-E8024]: Spalte existiert nicht an der Quelle: public.orders.regio` |
| R3 — `when.column` ist keine ausgeschlossene Spalte | `abgelehnt [PCF-E8033]: Spalte ist ausgeschlossen: public.orders.name` |
| R3 — umgekehrt: `cdc.exclude_column` gegen eine Spalte mit Routing-Bedingung | `abgelehnt [PCF-E8034]: Spalte trägt eine Routing-Bedingung: public.orders.region` |
| R4 — höchstens eine Abschlussregel (ohne `when`) | `abgelehnt [PCF-E8032]: Regel ohne when bereits vorhanden: public.orders.rest` |
| R4 — die Abschlussregel trägt die höchste `order` | `abgelehnt [PCF-E8032]: Regel ohne when trägt nicht die höchste order: public.orders.rest` |
| R4 — eine Regel mit `when` trägt keine höhere `order` als die Abschlussregel | `abgelehnt [PCF-E8032]: order liegt hinter der Regel ohne when: public.orders.rest` |
| R5 — das Paar (`when.column`, `when.equals`) kommt je Tabelle einmal vor | `abgelehnt [PCF-E8031]: Bedingung bereits vergeben: public.orders.region` |
| R6 — `cdc.remove_route` gegen einen unbekannten Regelnamen | `abgelehnt [PCF-E8023]: Regelname nicht geführt: public.orders.nope` |

Dazu kommen die Formprüfungen vor R1: `abgelehnt [PCF-E8010]: Regelname ist
ungültig`, `abgelehnt [PCF-E8011]: rule_spec ist ungültig`, `abgelehnt
[PCF-E8011]: unbekannter Schlüssel in rule_spec` (Beispiel: `unbekannter
Schlüssel in rule_spec: foo`, auch für einen Schlüssel innerhalb von `when`) und
`abgelehnt [PCF-E8012]: Zielname ist ungültig`.
Der Regelstand bleibt bei jedem `failed` unverändert. *Ursprung:* übernommen — Code
und Text der Tabelle sind die Ausgabe von `error_message`, die der E2E-Test von
`make test-integration` gegen eine Compose-Umgebung (`cdc_admin`-Identität)
erwartet; `PCF-E8010` (Regelname) ist im Store-Test von `make test-store` unter
einem `cdc_admin`-Login belegt, nicht im E2E-Lauf. Beide Läufe wurden für diesen
Absatz nicht neu gefahren; die Adressen nennen die Namen des Beispiels.

**Lösung:** den betroffenen Namen oder die `order` ändern; bei R4 mit
Abschlussregel: wer eine Regel mit höherer `order` als die Abschlussregel
ergänzen will, entfernt die Abschlussregel und setzt sie nach der neuen Regel
mit der dann höchsten `order` neu.

**Fehler: Erfassung endet mit der Fehlerklasse `schema`**

Eine Routing-Regel kann auf zwei Wegen zu dieser Klasse führen; die Abhilfe
unterscheidet sich.

**Ursache 1 — die Regel ist auf die Change nicht anwendbar.** Die
Bedingungsspalte `when.column` fehlt in der Relation der Change, ohne dass die
bekannte Spaltenform der Tabelle eine Entfernung zeigt. Das entsteht bei der
Erstaktivierung einer Tabelle ohne bekannte Spaltenform (die Version der Tabelle
trägt noch keine Spaltenzeilen), etwa mit einer Publication mit Spaltenliste, die
die Bedingungsspalte nicht enthält (*Ursprung:* übernommen aus dem E2E-Lauf von
`make test-integration`; gemessen an PostgreSQL 17.11, an PostgreSQL 18.6
übernommen). Eine
Publication mit Spaltenliste an einer Tabelle **mit** bekannter Spaltenform
gehört nicht hierher, sondern zu Ursache 2. Die Erfassung der **gesamten Quelle** endet sichtbar
mit der Klasse `schema` (siehe [Fehlerklassen](#fehlerklassen)): der Log nennt
die Regel und die Spalte, `cdc.heartbeat.error_class` zeigt `schema`, keine
Change der Tabelle ist persistiert, es geht keine Change verloren.

```text
Fehlerklasse schema [PCF-E4005]: Routing-Regel auf die Änderung nicht anwendbar: Regel "rt_eu" an public.rt2: Spalte der Routing-Regel fehlt in den Spalten der Änderung: region
```

**Lösung:** `cdc.remove_route` beantragen — der Antrag bleibt `pending`, solange
der Prozess steht —, den Prozess neu starten; offene Anträge werden beim Start
innerhalb der Frist des Vorlaufs verarbeitet, bevor die erste Transaktion der
Tabelle verarbeitet wird, und die zuvor nicht bestätigte Transaktion erscheint
danach über `cdc.changes` ohne Ziel (`route_target IS NULL`), ohne zweiten
`schema`-Fehler. *Ursprung:* übernommen aus demselben E2E-Lauf (Antrag
`pending`, nach dem Neustart `applied`) und gemessen in einer Compose-Umgebung
(PostgreSQL 18.6: nach dem Neustart `applied`, die Zeile ohne Ziel,
Feed-Container `healthy`). Der Kopf mit Klasse und Meldungscode im Log-Text oben
ist *gemessen* im Container-Log des E2E-Laufs von `make test-integration`, nicht
in dieser Compose-Umgebung.

**Ursache 2 — die Spalte wurde an der Quelle entfernt oder fehlt in der
Relation, die Spaltenform der Tabelle ist bekannt.** Zwei Anlässe: die Spalte
wird an der Quelle entfernt, oder eine Publication mit Spaltenliste, die eine
bekannte Spalte nicht enthält, wird für die Tabelle gesetzt. Die Erfassung endet
nicht an der Regel, sondern als inkompatible Schemaänderung: der Log nennt
„Relation-Änderung nicht sicher als Obermenge
interpretierbar", und die Change nach der Entfernung hat keine Zeile in
`cdc.changes` (*Ursprung:* für die entfernte Spalte übernommen aus dem E2E-Lauf
von `make test-integration`: Spalte entfernt, nachdem eine Change der Tabelle
die Spaltenform angelegt hat).
Für die Spaltenliste gemessen in einer Compose-Umgebung (PostgreSQL 18.6, Feed
als Superuser, Tabelle `public.orders` mit bekannter Spaltenform und einer Regel
auf `region`): nach `ALTER PUBLICATION … SET TABLE public.orders (id, name)` und
einer eingefügten Zeile endet der Prozess mit Exit 1, gedruckt „Relation-Änderung
nicht sicher als Obermenge interpretierbar: public.orders“,
`cdc.heartbeat.error_class` zeigt `schema`. Der Kopf mit Klasse und Meldungscode
(`Fehlerklasse schema [PCF-E4003]: …`) ist für diese Fehlerart im Container-Log des
E2E-Laufs von `make test-integration` *gemessen*, nicht in dieser
Compose-Umgebung. **Das Entfernen der
Regel genügt dort nicht:** nach dem Entfernen beider Regeln der Tabelle und einem
Neustart endete die Erfassung am System mit demselben Text, ebenso nach dem
Zurücksetzen der Publication auf die volle Spaltenliste (*Ursprung:* übernommen,
am System gefahren). Die Abhilfe dieser
Ursache ist die der inkompatiblen Schemaänderung; zusätzlich ist die Regel auf
die entfernte oder nicht mehr
gelieferte Spalte zu entfernen, sonst endet die Erfassung danach an der fehlenden
Bedingungsspalte (*abgeleitet*). Eine Anleitung zur Abhilfe der inkompatiblen
Schemaänderung ist in diesem Handbuch nicht beschrieben.

**Im Backfill:** Ein Run (siehe [Bestand als Backfill
überführen](#bestand-als-backfill-überführen)) prüft dieselbe Anwendbarkeit gegen
die Spalten des Snapshots. Eine Regel, deren `when.column` der Snapshot nicht
trägt, endet den Run `failed` mit der Klasse `schema` vor der ersten Zeile; der
Text nennt Regelname und Spalte, der Run ist run-lokal und stoppt die Erfassung
nicht. Die Lösung ist `cdc.remove_route` und ein **neuer** `cdc.backfill_table`-Antrag,
ohne Prozessneustart. Wird ein `set_route` oder `remove_route` während des Runs
`applied`, endet der Run `failed` mit der Klasse `configuration` und dem Text
„Routing-Regelstand während des Backfills geändert", bevor etwas sichtbar wird;
die Lösung ist ein neuer Antrag. Ein Lesefehler des Regelstands endet den Run mit
der Klasse seiner Ursache (z. B. `storage`), nicht mit `configuration`.
*Ursprung:* durch Unit-Tests belegt, am System nicht gefahren.

**Ziel lesen:** Ein Leser wählt ein Ziel; ohne Auswahl sieht er weiterhin alle
Changes, geroutete eingeschlossen.

| Weg | Auswahl |
|---|---|
| SQL (`cdc.changes`) | Spalte `route_target`: `WHERE route_target = '<ziel>'` bzw. `IS NULL` für „ohne Ziel" |
| `GET /changes` | Query-Parameter `target` (siehe [Zugriff über die HTTP-/JSON-API](#zugriff-über-die-http-json-api)) |
| gRPC-RPC `ReadChanges` | Feld `target` (siehe [Zugriff über die gRPC-Verwaltungs-API](#zugriff-über-die-grpc-verwaltungs-api)) |
| gRPC-Stream | Feld `target` des `StreamChangesRequest` (siehe [Zugriff über den gRPC-Change-Stream](#zugriff-über-den-grpc-change-stream)) |
| SSE (`GET /changes/stream`) | Query-Parameter `target` (siehe [Zugriff über Server-Sent-Events](#zugriff-über-server-sent-events)) |
| NATS-Vollinhalt | Subjekt `cdc.route.<source_id>.<ziel>` (siehe [Zugriff über den NATS-Vollinhalts-Stream](#zugriff-über-den-nats-vollinhalts-stream)) |

Das NATS-Wecksignal trägt kein Ziel. Das Ziel steht in keiner Antwort und in
keiner Stream-Nachricht der Zugriffswege, nur in der Spalte `route_target`.

**Hinweise:**

- **Das Ziel ist zum Erfassungszeitpunkt fest.** Eine Regeländerung wirkt nur
  auf danach erfasste Changes; es gibt kein Umetikettieren. Eine vor der Regel
  erfasste Change bleibt ohne Ziel, und dieselbe `change_id` liefert nach einer
  Regeländerung über `cdc.changes` und `GET /changes?target=` dasselbe Ziel
  (*Ursprung:* übernommen aus dem E2E-Lauf von `make test-integration`,
  dort gedruckte Zeile
  `ROUTING-REPLAY WAL-Ziele ["" "eu" "" "europa"], Backfill-Ziele zeilenweise 1="europa" 2="europa" 3="" 4="europa"`).
- **Altbestand über einen neuen Backfill-Run.** Ein Backfill-Run
  (`cdc.backfill_table`, siehe [Bestand als Backfill
  überführen](#bestand-als-backfill-überführen)) trägt das Ziel des Regelstands
  zum Run in seinen Changes (`origin = 'backfill'`); die WAL-Changes derselben
  Zeilen behalten ihr Label. *Ursprung:* gemessen in einer Compose-Umgebung
  (PostgreSQL 18.6): mit den Regeln des Beispiels und nach dem Löschen der
  Zeile 5 trugen die fünf Backfill-Changes `eu` (Region `eu`) bzw. `sonstige`
  (übrige Werte, auch fehlende); die sechs WAL-Changes derselben Zeilen trugen
  unverändert `NULL` (Zeilen 1 bis 3), `eu` und `sonstige` wie bei ihrer
  Erfassung.
- **Eine Change ohne Treffer bleibt sichtbar.** Trifft keine Regel, ist
  `route_target` `NULL`: die Change geht nicht verloren und nicht an ein
  Standardziel, sie bleibt im Log, über jeden ungefilterten Weg und über
  `route_target IS NULL` sichtbar. Für einen Leser, der nur ein Ziel abruft, ist
  sie unsichtbar; eine Abschlussregel (ohne `when`, höchste `order`) vergibt ein
  Ziel an alles, was keine Bedingung trifft.
- **Ein abwesender Wert trifft nicht.** Fehlt der Wert der Bedingungsspalte im
  gelesenen Bild (`NULL`, unverändertes TOAST, generierte Spalte), trifft die
  Bedingung nicht; die nächste Regel wird geprüft — das ist keine
  Nichtanwendbarkeit und kein Fehler. Gemessen ist `NULL` (Zeile 6 des
  Beispiels, `DELETE` unten); für unverändertes TOAST und generierte Spalten gilt
  die Zusage der Spezifikation, am System nicht gefahren. Der
  Vergleich ist zeichengenau (Groß-/Kleinschreibung zählt).
- **`DELETE` ohne volle Replica-Identität.** Das Alt-Bild trägt dort nur den
  Schlüssel; eine Inhaltsregel auf eine Nicht-Schlüsselspalte **trifft nicht**.
  Mit einer Abschlussregel bestimmt diese das Ziel, ohne sie bleibt
  `route_target` `NULL`; unter `REPLICA IDENTITY FULL` trifft die Inhaltsregel.
  `INSERT` und `UPDATE` tragen das Ziel der Inhaltsregel. *Ursprung:* übernommen
  aus dem E2E-Lauf von `make test-integration` (Zeilen `ROUTING-DELETE-MESSUNG`)
  an PostgreSQL 17.11 und 18.6, beide mit denselben Zielen: `INSERT` und `UPDATE`
  `eu`; `DELETE` mit Abschlussregel `sonstige` (Alt-Bild `{"id": "1"}`), ohne sie
  nicht gesetzt, unter `REPLICA IDENTITY FULL` `eu`. Eine Bedingung auf die
  Schlüsselspalte trifft dort, weil das Alt-Bild den Schlüssel trägt
  (*abgeleitet*, nicht gefahren). In einer Compose-Umgebung (PostgreSQL 18.6)
  ergab ein `DELETE` ohne volle Replica-Identität mit den Regeln des Beispiels
  das Alt-Bild `{"id": "5"}` (die Zeile 5 hatte die Region `us`) und das Ziel
  `sonstige`.
- **Auswahl, kein Zugriffsschutz.** Jeder Leser mit `reader`-Token kann jedes
  Ziel wählen oder gar nicht filtern. Das Ziel schützt keinen Inhalt; wer einen
  Wert vor Lesern verbergen will, schließt die Spalte aus (siehe [Spalte vom
  Ausschluss konfigurieren](#spalte-vom-ausschluss-konfigurieren)). Eine Spalte mit
  Routing-Bedingung lässt sich nicht ausschließen, und eine ausgeschlossene Spalte
  kann keine Bedingung tragen (R3): sonst verriete das Ziel eine Eigenschaft des
  ausgeschlossenen Werts.
- **Verhältnis zu Transformationen.** Die Bedingung liest den Quellwert vor jeder
  Transformation; eine `rename_column`- oder `map_value`-Regel ändert das
  Ergebnis nicht, und das Ziel ändert das Row Image nicht. Beide Regelwerke haben
  getrennte Namensräume und getrennte Stände.

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
Positionen bestätigen. Ein bereits registrierter Name
bleibt unverändert; der Lauf meldet das über eine eigene Ausgabe-Zeile,
Exit-Code bleibt 0 (idempotent).

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

**Ergebnis:** Die Position ist bestätigt und fortgeschrieben.
Die Wiederholung derselben Position bleibt ohne Wirkung
(idempotent); ein echter Rückschritt (eine
Position vor dem bereits bestätigten Stand) wird abgelehnt, der
gespeicherte Fortschritt bleibt unverändert — die Vorwärts-Invariante gilt
für jeden Aufruf über diesen Zugriffsweg.

### Änderungen lesen

```sql
SELECT source_id, commit_position, change_id, schema_name, table_name,
       operation, old_data, new_data, committed_at, origin, route_target
FROM cdc.changes
WHERE source_id = '<quelle-id>' AND commit_position > <letzte-gelesene-position>
ORDER BY commit_position, sequence
LIMIT 500;
```

**Ergebnis:** Jede Zeile ist eine committed Änderung in Anhang-Reihenfolge.
`old_data`/`new_data` sind `jsonb`; bei `INSERT` ist `old_data` NULL, bei
`DELETE` ist `new_data` NULL. `origin` nennt die Herkunft der Änderung:
`wal` für eine über den Replication Stream erfasste Änderung, `backfill`
für eine Bestands-Änderung eines Backfills; eine
Änderung, die ohne dieses Feld gespeichert wurde, liest als `wal`.
`route_target` nennt das Zustellziel der Change, das eine
Routing-Regel bei der Erfassung vergeben hat (siehe [Routing-Regel
konfigurieren](#routing-regel-konfigurieren)); `NULL` heißt „nicht geroutet" —
keine Regel hat getroffen, oder die Change entstand vor der Regel —, ein
Standardziel gibt es nicht (anders als bei `origin` liest ein fehlender Wert
nicht als Vorgabe). `route_target` ist die letzte Spalte der View,
`origin` steht davor; die drei Live-Zustellwege (gRPC, SSE, NATS-Vollinhalt)
tragen keines der beiden Felder.

Dieselben Änderungen sind ohne SQL-Direktzugriff über die API lesbar:
`GET /changes` — derselbe Lesezugriff mit denselben Filtern und derselben
Reihenfolge (siehe
[Zugriff über die HTTP-/JSON-API](#zugriff-über-die-http-json-api)).

**Fortsetzen und `LIMIT`:** Ein `LIMIT` schneidet Zeilen, nicht Positionen.
Trug eine Commit-Position mehr Änderungen als das `LIMIT`, überspringt
`commit_position > <letzte-gelesene-position>` den Rest dieser Position.
Das trifft den Bestandsabzug eines
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
bereits enthält, als `INSERT`-Änderungen mit `origin = 'backfill'` in den Feed.
Er wird **ausdrücklich** ausgelöst — die
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
  nicht über die Fehlerschwelle — dasselbe gilt für einen Schreiber auf eine
  nicht aktivierte Tabelle, dessen WAL der Feed zwischen den Schreibvorgängen
  bestätigt. Ein einzelner Stoß, der die Fehlerschwelle in einem Zug
  überschreitet, steht bis zur nächsten Bestätigung im Rückstand (Grenze unter
  [WAL-Rückstand prüfen](#wal-rückstand-prüfen)). Was bleibt, ist die offene
  Schreibtransaktion des Runs: sie hält das WAL für den Slot auf der Platte der
  Quelle, und der Walsender der Quelle dekodiert sie in seinen Arbeitsspeicher
  und lagert sie ab einer Größe auf die Platte aus. Beides wächst mit der Größe
  der Tabelle und ist eine Eigenschaft der Ein-Transaktions-Form des Runs; die
  Bestätigung im Leerlauf ändert es nicht. Gemessene Werte stehen unter
  [Grenzwerte](#grenzwerte). Wächst der Rückstand dennoch über die
  Fehlerschwelle (etwa weil der Feed nicht antwortet), beendet sich der
  Feed-Container mit der Klasse `replication` (Ausgang 1) und ein laufender Run
  endet `interrupted`; das Datei-Feld `wal_retention_error_bytes` (Bytes, kein
  Umgebungsvariablen-Gegenstück, siehe
  [Optionale YAML-Konfigurationsdatei](#optionale-yaml-konfigurationsdatei-cdc_config_file))
  hebt die Fehlerschwelle über den erwarteten Rückstand — es ändert die Ursache
  nicht und wirkt für jede Ursache eines Rückstands.

**Speicher des Feed-Containers:** Ein Backfill lässt die Zahl der Changes in
`cdc.change` um die Zeilenzahl der Tabelle wachsen. Der Bereinigungslauf des Feeds
liest sie in jedem Takt seitenweise (10.000 Changes je Seite, ohne die Row Images,
siehe [Aufbewahrung (Retention)](#aufbewahrung-retention)); der Speicherbedarf des
Feed-Containers hängt deshalb nicht mit nennenswertem Betrag an dieser Zahl.
Gemessen liegt die Spitze des Feed-Containers nach Backfills über je 1.000.000
Zeilen bei 1.000.000 bis 3.000.000 Changes in `cdc.change` bei 14,9 bis 17,6 MiB
(zwei Läufe; über 3.000.000 Changes ist nichts gemessen; Messwerte und Herkunft
unter [Grenzwerte](#grenzwerte)).

**Sperre der Tabelle:** Vom Beginn der Lese-Transaktion bis zu ihrem Ende hält
der Run eine Lesesperre (`ACCESS SHARE`) auf die Tabelle. Lesen und Schreiben
der Tabelle laufen weiter, solange keine DDL auf die Tabelle wartet. Eine DDL,
die `ACCESS EXCLUSIVE` verlangt (`ALTER TABLE` mit Umschreiben der Tabelle,
`TRUNCATE`, `VACUUM FULL`, `CLUSTER`, `DROP TABLE`), wartet bis zum Ende des
Runs, und **jeder weitere Zugriff auf die Tabelle stellt sich hinter sie**:
Schreiber und Leser der Tabelle ebenso wie die Abfrage der
Publication-Mitgliedschaft (`pg_publication_tables`), die ein Antrag und der
Start eines Runs ausführen: mit einer wartenden `ALTER TABLE … ALTER COLUMN …
TYPE` liefen ein `INSERT` in die Tabelle, ein `SELECT count(*)` auf die
Tabelle und die Abfrage von `pg_publication_tables` je in ein Zeitlimit von
4 s. Führen Sie eine solche DDL an einer Tabelle nicht aus, solange ein Run
für sie `running` ist (Status in `cdc.backfill_status`); die Dauer des Runs
wächst mit der Größe der Tabelle.

In der Gegenrichtung wartet der Run, solange eine fremde Transaktion
`ACCESS EXCLUSIVE` auf der Tabelle hält: er steht `running` mit `rows_copied` 0
und hält für diese Zeit seinen Snapshot. Der Run trägt dafür keine eigene
Zeitgrenze; das Warten endet mit dem Ende der fremden Transaktion oder mit
dem Abbruch des Runs (Ablauf seines Kontexts; der Adapter beendet ihn mit
der Klasse `transient` und hinterlässt keine Sitzung). Beenden Sie eine
offene DDL-Transaktion auf der Tabelle, statt den Run warten zu lassen.

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
| `failed` | der Run endete mit einem Fehler; `error_message` trägt Fehlerklasse und Meldungscode (siehe [Meldungscodes](#meldungscodes)) vor der Ursache; der Run hinterlässt keine Änderung |
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
  Tabelle und endet `failed` mit `error_message` `transient [PCF-E1003]: …`, Ursache
  „nach dem Snapshot-Export umgeschrieben“ und ohne Änderung; ein **neuer
  Antrag** beginnt neu und liest den Bestand des neuen Zustands. Dieselbe
  Erkennung schlägt auch bei `VACUUM FULL` und `CLUSTER` an, obwohl der Snapshot
  die Zeilen noch sieht — ein Fehlalarm mit derselben Abhilfe. Das Fenster ist
  die Zeit zwischen Export und Sperre; seine Dauer ist nicht gemessen. Ein
  `DROP COLUMN` im Fenster endet ebenfalls `failed`, mit der Klasse `storage`;
  ein `RENAME COLUMN` endet gleich.
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
  Lauf von `make test-integration` (siehe
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
- **Routing-Regeln** (siehe [Routing-Regel
  konfigurieren](#routing-regel-konfigurieren)) vergeben den Backfill-Änderungen
  das Ziel des Regelstands zum Run (Spalte `route_target`); der Altbestand erhält
  ein Ziel durch einen neuen Run mit dem aktuellen Regelstand. Ein
  `cdc.set_route` oder `cdc.remove_route`, das während des Runs `applied` wird,
  endet den Run `failed` (Klasse `configuration`), bevor etwas sichtbar wird; eine
  Regel, deren Bedingungsspalte der Snapshot nicht trägt, endet ihn `failed`
  (Klasse `schema`) vor der ersten Zeile, die Abhilfe steht unter „Fehler:
  Erfassung endet mit der Fehlerklasse `schema`" im Abschnitt Routing.
- **Zustellung:** Backfill-Änderungen gehen in keinen der Live-Zustellwege (gRPC,
  SSE, NATS-Vollinhalt, auch nicht auf das Zusatz-Subjekt `cdc.route.…`); nach dem Commit sendet der Run je Tabelle ein
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
Der Speicher des Feed-Containers hängt deshalb an der Seitengröße und nicht mit
nennenswertem Betrag an der Zahl der Changes in `cdc.change` (gemessen 14,9 bis
17,6 MiB bei 1.000.000 bis 3.000.000 Changes; Messwerte unter
[Grenzwerte](#grenzwerte)); die Arbeit der Datenbank je Durchlauf wächst dagegen
mit dieser Zahl. Ein Durchlauf
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
läuft, gilt seine Position erst im nächsten Durchlauf: Changes, die dieser
Durchlauf nach dem Lesen der Positionen löscht, sind für diesen Consumer
verloren, sein Schutz beginnt mit dem nächsten Durchlauf. Ein Durchlauf kann
eine Transaktion in mehreren Schritten löschen; ein Consumer mit bestätigter
Position sieht davon nichts, ein Consumer ohne Bestätigung kann eine alte
Transaktion währenddessen unvollständig lesen.

Direkter SQL-Zugriff auf die betroffenen Tabellen bleibt zulässig
(`cdc_admin`-Mitgliedschaft, verbunden über `CDC_ADMIN_DSN`):

```sql
SELECT change_id, committed_at FROM cdc.changes WHERE source_id = '<quelle-id>' ORDER BY commit_position;
```

Eine Zeile, die dort nicht mehr erscheint, wurde bereits bereinigt.

### Blockierende Consumer erkennen

`cdc.retention_blockers` zeigt je Quelle den Consumer, dessen bestätigte
Position aktuell die Löschgrenze der Bereinigung trägt —
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
SELECT source_id, heartbeat_at, age_seconds, error_class, error_code FROM cdc.heartbeat;
```

**Ergebnis:** `error_class` und `error_code` sind `NULL` im Normalbetrieb; ein
gesetzter Wert nennt die zuletzt beobachtete Fehlerklasse und den Meldungscode
des Fehlerzustands (siehe [Fehlerbehebung](#6-fehlerbehebung) und
[Meldungscodes](#meldungscodes)). `error_code` steht in der Sicht hinter
`age_seconds`; der Code erscheint nur neben einer Klasse. Ein Lebenszeichen gilt als
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
`cdc_consumer_lag` je registriertem Consumer (WAL-Strecke in Bytes, nicht
Datenvolumen), `cdc_changes_pending`
(Label = Consumer-ID, Wert = Anzahl noch nicht
bestätigter Changes dieses Consumers), `cdc_errors_total`
(Label = Fehlerklasse aus `cdc.process_heartbeat.error_class`, Wert =
Anzahl der Quellen aktuell in dieser Fehlerklasse), sowie
`cdc_storage_bytes` (physische Speichergröße von
`cdc.change` über `pg_relation_size` — der mit dem Erfassungsvolumen
wachsenden Tabelle).

`cdc_wal_retention_bytes` (WAL-Rückstand des Capture-Slots)
steht **nicht** in `cdc.metrics`: Die Erhebung braucht Systemkatalog-Zugriffe
außerhalb des `cdc`-Schemas (`pg_replication_slots`, `IDENTIFY_SYSTEM`), die
die Least-Privilege-Fläche von `cdc_reader` unnötig erweitern würden — siehe
[WAL-Rückstand prüfen](#wal-rückstand-prüfen). Der OTLP-Export (nächster
Abschnitt) überträgt ihn zusätzlich zu den Zeilen der Sicht.

### Metriken per OTLP übertragen

Statt die Sicht abzufragen, kann der Feed-Container seine Kennzahlen aktiv an
einen OpenTelemetry-Empfänger übertragen (Push), etwa an einen
OpenTelemetry-Collector. Der Export ist **aus**, solange `CDC_OTLP_ENDPOINT`
nicht gesetzt ist: dann baut der Container keine Verbindung zu einem Empfänger
auf, und die SQL-Sicht `cdc.metrics` sowie der Betrieb bleiben unverändert. Die
Sicht bleibt auch mit Export bestehen; der Export ersetzt sie nicht.

| Variable | Bedeutung | Gültig |
|---|---|---|
| `CDC_OTLP_ENDPOINT` | Basis-URL des Empfängers; gesetzt schaltet den Export ein | Schema `http` oder `https` und ein Host |
| `CDC_OTLP_HEADERS` | zusätzliche Header jeder Anfrage, Form `Schlüssel=Wert`, mehrere durch Komma getrennt | jedes Element mit nicht leerem Schlüssel und einem `=` |
| `CDC_OTLP_INTERVAL_SECONDS` | Takt der Übertragung in Sekunden, Default 60 | ganze Zahl von 5 bis 3600 |

```yaml
services:
  pg-change-feed:
    environment:
      CDC_OTLP_ENDPOINT: "http://otel-collector:4318"
      CDC_OTLP_HEADERS: "Authorization=Bearer <token>"
      CDC_OTLP_INTERVAL_SECONDS: "30"
```

Der Container sendet OTLP über HTTP mit Protobuf-Körper
(`POST <CDC_OTLP_ENDPOINT>/v1/metrics`, Inhaltstyp `application/x-protobuf`);
der Pfad `/v1/metrics` kommt zur Basis-URL hinzu, mit oder ohne abschließenden
Schrägstrich. Ein Pfadanteil der Basis-URL bleibt erhalten. Ein Status `2xx`
gilt als angenommen; jeder andere Status, eine abgewiesene Verbindung und eine
überschrittene Frist gelten als Fehlschlag. Eine Weiterleitung (`3xx`) wird nicht
verfolgt und ist ein Fehlschlag.

**`https`-Empfänger.** Das Container-Image enthält ein CA-Zertifikatsbündel
(`/etc/ssl/certs/ca-certificates.crt`). Ein `https`-Empfänger mit einem
selbstsignierten Zertifikat wird erreicht, wenn die Umgebungsvariable
`SSL_CERT_FILE` des Containers auf die Zertifikatsdatei des Empfängers zeigt
(die Datei wird in den Container eingebunden); ohne sie scheitert der
TLS-Aufbau zum Empfänger, weil dessen Zertifikat nicht vertrauenswürdig ist; der
Container warnt mit `PCF-W6001`, läuft weiter, und es kommt kein Export an.
Gemessen mit OpenTelemetry-Collector 0.162.0 und einem selbstsignierten
Zertifikat; über Zertifizierungsstellen öffentlicher Empfänger macht diese
Messung keine Aussage. Ein Client-Zertifikat, eine eigene Zertifizierungsstelle
als Einstellung des Exports und Proxy-Einstellungen sind nicht konfigurierbar.

**Zugangsdaten in der Endpunkt-URL.** Ein Benutzerteil der URL
(`http://benutzer:passwort@host:4318`) wird nicht abgelehnt: der Container sendet
ihn als `Authorization: Basic …`. Setzt `CDC_OTLP_HEADERS` selbst einen
`Authorization`-Header, gilt dieser. Weder ein Log-Eintrag noch ein Fehlertext
des Exports nennt die URL, den Benutzerteil oder einen Header-Wert. Endpunkt und
Header gehören zu den Angaben mit Zugangsdaten-Charakter und stehen nur in der
Umgebung des Containers, nicht in der Konfigurationsdatei.

**Header.** Das Komma trennt die Header, das **erste** `=` trennt Schlüssel und
Wert; ein Wert darf weitere `=` tragen (etwa die Auffüllung einer
Base64-Zeichenkette) und kein Komma. Der Schlüssel ist ein gültiger
Header-Name ohne Leerraum; Leerraum im Wert bleibt erhalten
(`Authorization=Bearer abc`). Ein Element ohne `=`, mit leerem Schlüssel, mit
Leerraum oder einem unzulässigen Zeichen im Schlüssel oder mit einem
Steuerzeichen im Wert verhindert den Start. Header ohne Endpunkt verhindern den
Start ebenso: ein Wert, der ohne Endpunkt nie wirkt, bleibt nicht stillschweigend
stehen. Die Header können Zugangsdaten tragen und stehen in keinem Log-Eintrag.
Der Takt ohne Endpunkt ist zulässig (die Konfigurationsdatei darf ihn tragen).

**Kennzahlen.** Übertragen werden die Zeilen der Sicht `cdc.metrics` — der
Container liest sie über `CDC_READER_DSN` (Rolle `cdc_reader`, ohne zusätzliche
Rechte) — und der zuletzt vom Container gemessene `cdc_wal_retention_bytes`
(siehe [WAL-Rückstand prüfen](#wal-rückstand-prüfen)). Ist der WAL-Rückstand noch
nicht gemessen, entfällt genau dieser Datenpunkt.

| Kennzahl | Einheit | Attribut des Datenpunkts |
|---|---|---|
| `cdc_transactions_total` | `1` | — |
| `cdc_changes_processed` | `1` | — |
| `cdc_oldest_change_age_seconds` | `s` | — |
| `cdc_capture_lag` | `s` | — |
| `cdc_consumer_position` | `1` | `consumer` |
| `cdc_consumer_lag` | `By` | `consumer` |
| `cdc_changes_pending` | `1` | `consumer` |
| `cdc_errors_total` | `1` | `class` |
| `cdc_storage_bytes` | `By` | — |
| `cdc_wal_retention_bytes` | `By` | — |

Die Resource trägt die Attribute `service.name` mit dem Wert `pg-change-feed`
und `cdc.source_id` mit der Kennung der Quelle des Containers; jeder Datenpunkt
trägt den Zeitpunkt seiner Messung. `cdc_consumer_lag` ist die
Strecke im WAL in Bytes zwischen der letzten Commit-Position und der bestätigten
Position des Consumers — eine WAL-Strecke in Bytes, nicht das Datenvolumen der noch
nicht bestätigten Changes. Zähler, Positionen und Byte-Zahlen gehen als
ganze Zahlen hinaus und behalten ihre volle 64-Bit-Genauigkeit; die Kennzahlen mit
Einheit `s` gehen als Fließkommazahlen hinaus.

**Alle Kennzahlen sind Momentstände (Gauge), auch die mit dem Namensteil
`_total`.** Ein Wert ist kein monoton wachsender Zähler: die Zahl der Changes fällt
bei einer Löschung nach der Aufbewahrung, `cdc_errors_total` zählt die Quellen, die
aktuell in einer Fehlerklasse stehen. Ein Werkzeug, das den Namensteil `_total` als
Zähler liest, darf darauf keine Zählerfunktionen wie `rate()` oder `increase()`
anwenden; ein Gauge wird als Momentwert gelesen.

**Wert der Sicht.** Ein übertragener Wert entspricht dem Wert, den die SQL-Sicht
zum Zeitpunkt seiner Messung liefert. `cdc_oldest_change_age_seconds` und
`cdc_capture_lag` hängen an der Uhr der Abfrage und weichen um den Abstand beider
Messzeitpunkte ab; `cdc_wal_retention_bytes` steht nicht in der Sicht, übertragen
wird der zuletzt gemessene Wert.

**Takt, Frist und Ausfall.** Der erste Versuch folgt nach einem vollen Takt, nicht
beim Start. Jeder Versuch hat eine Frist von der Länge des Taktes, höchstens
10 Sekunden. Es gibt keine Warteschlange: der nächste Takt ist die Wiederholung,
ein ausgefallener Takt wird nicht nachgeholt, und der nächste Takt überträgt den
dann aktuellen Stand. Der Export läuft unabhängig von der Erfassung in einem
eigenen Ablauf des Prozesses und hält weder die Erfassung noch die Persistierung
noch die Bestätigung auf; ein Empfänger, der
nicht erreichbar ist oder nicht antwortet, ändert weder den Health-Zustand noch
`error_class` des Lebenszeichens. Die Übertragung setzt ein, sobald der Empfänger
wieder antwortet.

Ein Fehlschlag steht im Log als Warnung: `PCF-W6001` für die Übertragung,
`PCF-W6002` für das Lesen der Sicht. Die erste Warnung seit dem Start oder seit der
letzten Wiederaufnahme steht sofort, weitere höchstens alle 5 Minuten; die
Wiederaufnahme erzeugt eine Info-Zeile. Die Zeile nennt Code und Stufe sowie —
bei der Übertragung — den Status oder die Art des Fehlers, nie eine URL oder einen
Header-Wert.

```json
{"level": "WARN", "msg": "metrics-export: Zyklus fehlgeschlagen", "code": "PCF-W6001", "stage": "Übertragung", "detail": "Status 503"}
```

**Ungültige Konfiguration.** Ein ungültiger Endpunkt, Header oder Takt beendet
den Start mit der Fehlerklasse `configuration`; die Zeile nennt die Variable und
bei Header die Position des Elements, nie den Wert von Endpunkt oder Header (beide
können Zugangsdaten tragen); beim Takt nennt sie den eingegebenen Wert der
Intervall-Einstellung: `PCF-E2011` (Endpunkt), `PCF-E2012` (Header),
`PCF-E2013` (Takt), `PCF-E2014` (Header ohne Endpunkt).

### Diagnose ausführen

Statt der SQL-Abfragen oben einzeln zu stellen, liest der
`diagnose`-Sondermodus dieselben Views (`cdc.heartbeat`, `cdc.metrics`,
`cdc.retention_blockers`, `cdc.backfill_status`) über `CDC_READER_DSN` und gibt eine
menschenlesbare Zusammenfassung aus — derselbe Image-Tag wie der
Daemon, als einmaliger, kurzlebiger Lauf statt als Dauerdienst:

```bash
docker run --rm -e CDC_CAPTURE_DSN -e CDC_ADMIN_DSN -e CDC_READER_DSN \
  -e CDC_SOURCE_ID -e CDC_PUBLICATION -e CDC_SLOT -e CDC_TABLES \
  ghcr.io/pt9912/pg-change-feed:dev diagnose
```

In der Compose-Umgebung: `docker compose run --rm pg-change-feed diagnose`,
oder gegen einen bereits laufenden Feed-Container: `docker exec <container>
/pg-change-feed diagnose`.

Derselbe Bericht ist zusätzlich über zwei Netzwerkzugriffswege lesbar —
ohne `docker exec` und ohne SQL-Direktzugriff, fachlich gleichwertig, kein
zweiter Domänenpfad: `GET /diagnose?source=<quelle-id>` (siehe [Zugriff über
die HTTP-/JSON-API](#zugriff-über-die-http-json-api)) und der gRPC-RPC
`Diagnose` (siehe [Zugriff über die gRPC-Verwaltungs-API](#zugriff-über-die-grpc-verwaltungs-api)),
beide mit der Rechtsklasse `reader`. Der CLI-Sondermodus bleibt der einzige
Zugriffsweg für ein Deployment ohne `CDC_HTTP_ADDR`/`CDC_GRPC_ADDR`.

**Ausgabe (Beispiel):**

```text
pg-change-feed diagnose: Quelle "src-e2e"
  Betriebsstatus: Lebenszeichen vor 1.203s
  Fehlerzustand: keiner (Normalbetrieb)
  CDC-Abstand cdc_capture_lag: 0.087s
  Verarbeitungsrückstand cdc_consumer_lag je Consumer (nur Consumer mit mindestens einer bestätigten Position):
    cli-e2e-consumer: 3
  Blockierender Consumer: Orders-Reader (cli-e2e-consumer), bestätigte Position 42, Rückstand 3
  Speicherverbrauch cdc_storage_bytes: 65536 Bytes
  Backfill je Tabelle (letzter Run; die Zeilenzahl ist geschätzt):
    public.orders: completed, 1200 Zeilen kopiert, geschätzt 1150, Warnung Größe false, Warnung Dauer false
    public.audit: failed, 0 Zeilen kopiert, geschätzt unbekannt, Warnung Größe false, Warnung Dauer false
      Fehler: permission: …
    public.events: completed, 4800000 Zeilen kopiert, geschätzt 4700000, Warnung Größe true, Warnung Dauer true
```

Bei einem Fehlerzustand nennt die Zeile „Fehlerzustand“ die Fehlerklasse und
den Meldungscode in eckigen Klammern (siehe [Meldungscodes](#meldungscodes));
alle übrigen Zeilen des Berichts tragen keinen Code. *Ursprung:* gemessen mit
dem Sondermodus `diagnose` gegen eine mit `make schema-rollout` ausgerollte
Instanz, in der `error_class` und `error_code` direkt in
`cdc.process_heartbeat` gesetzt waren:

```text
  Fehlerzustand: schema [PCF-E4003]
```

Trägt der Heartbeat nur eine Klasse (geschrieben von einer Version ohne
Meldungscode), steht dort `Fehlerzustand: schema` (dieselbe Messung).

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
  Blockierender Consumer: kein Blocker (kein Consumer hat je gegen diese Quelle bestätigt)
```

**Ergebnis:** Wie bei den Rohwerten der Views trifft der Befehl keine
Schwellenwert-Entscheidung (die bleibt Sache des lesenden Systems) und
der Prozess-Ausgang trägt nur den Lese-Erfolg — ein gemeldeter Fehlerzustand
oder Rückstand ist Berichtsinhalt, kein Befehlsfehler (Ausgang bleibt 0). Ein
Consumer ohne je bestätigte Position erscheint nicht in der Rückstands-Liste
(dieselbe Grenze wie bei `cdc_consumer_lag` in [Metriken
lesen](#metriken-lesen)); ein Consumer mit bestätigter Position, dessen
Quelle noch nie eine Transaktion trug, erscheint mit dem Text „unbekannt"
statt einem irreführenden Rückstand von 0. Der Rückstand je Consumer ist die
WAL-Strecke in Bytes, nicht das Datenvolumen der noch nicht bestätigten Changes.
Der blockierende Consumer und
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
Feed im Leerlauf seines Streams selbst, sobald die Quelle es geliefert hat; der
Wert steigt dabei um das WAL, das die Quelle seit der letzten Bestätigung
geschrieben hat, und fällt mit der nächsten zurück. Der Integrationstest
(`make test-integration`) belegt das mit einer Last in Stücken, die einzeln
unter der Warnschwelle liegen und zusammen die Fehlerschwelle übersteigen, mit
Bestätigung des Slots zwischen den Stücken. `bytes` wächst, solange der Feed nicht
bestätigt: der Capture-Slot ist inaktiv und die Quelle schreibt weiter (z. B.
während eines Verbindungsabbruchs), oder der Feed antwortet nicht — ein
dauerhaft wachsender Wert ist ein Warnsignal für WAL-Erschöpfung auf der
Quelle. Der Feed-Container vergleicht den gemessenen Wert bei jedem Takt
gegen zwei Schwellen:

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

**Grenze der Bestätigung im Leerlauf.** Die Bestätigung folgt der Quelle: WAL,
das die Quelle in einem Zug schreibt und dem Feed noch nicht geliefert hat,
steht bis zur nächsten Bestätigung im Rückstand. Schreibt Ihre Quelle in einer
Bestätigungs-Runde mehr WAL, als `wal_retention_error_bytes` zulässt, kann der
Container an diesem Stoß enden (Klasse `replication`), obwohl der Feed gesund
ist. Heben Sie die Fehlerschwelle deutlich über das WAL, das Ihre Quelle in
wenigen Sekunden schreibt
([Optionale YAML-Konfigurationsdatei](#optionale-yaml-konfigurationsdatei-cdc_config_file));
bei einer gesenkten Fehlerschwelle liegt diese Grenze entsprechend niedriger.

### Schema aktualisieren

Änderungen am neutralen Schema (`tools/schema/schema.yaml`) rollen Sie
mit demselben Befehl wie bei der Ersteinrichtung erneut aus:

```bash
make schema-rollout SCHEMA_TARGET="db:<ihre-postgres-dsn>"
```

Der Lauf erzeugt einen Pflicht-Report (`plan.yaml`), ein Rollback-Artefakt
(`down.sql`) und den Precheck-Report (`rollout-precheck.yaml`). Sie liegen in
dem Verzeichnis, das die Variable `SCHEMA_ARTEFACT_DIR` nennt (Default
`.tmp/schema-rollout`, durch `.gitignore` ausgenommen); der Lauf druckt den
Pfad am Ende. Das Verzeichnis wird bei jedem Lauf überschrieben — wer die
Berichte je Rollout als Beleg hält, setzt `SCHEMA_ARTEFACT_DIR` auf ein Ziel je
Rollout (etwa ein Verzeichnis mit Datum) oder kopiert die Dateien nach dem
Lauf. Der Lauf mountet den Arbeitsbaum nicht in einen Container und lässt ihn
unverändert.

**Reihenfolge beim Upgrade.** Rollen Sie das Schema **vor** dem Tausch des
Feed-Containers aus, dann ersetzen Sie den Container durch die neue Version.
Der neue Container erwartet das neue Schema (er schreibt zum Beispiel die
Spalten `cdc.change.origin`, `cdc.change.route_target` und
`cdc.process_heartbeat.error_code`; ohne die letzte schlägt jeder Schreib-Zug
des Lebenszeichens fehl, *abgeleitet* aus dem SQL-Text, nicht gegen ein Alt-Schema
gefahren); umgekehrt ist der Rollout vor dem Tausch für
den weiterlaufenden Container erwartungsgemäß unkritisch, weil neue Spalten
nullable und neue Tabellen oder Views additiv sind (abgeleitet, nicht mit
einem laufenden Alt-Container gemessen).

**Rückweg auf eine ältere Version.** Ein älterer Server kennt `error_code`
nicht: bei einem Fehlerzustand schreibt er nur `error_class` und lässt einen
früher von der neuen Version gesetzten `error_code` stehen. Die Sicht
`cdc.heartbeat` blendet den Code nur aus, solange `error_class` leer ist; nach
einem Fehlerzustand anderer Klasse kann sie Klasse und Code verschiedener
Klassen zeigen. Maßgeblich ist dann `error_class`; ein Code, dessen erste
Ziffer nicht zur Klasse passt (siehe [Meldungscodes](#meldungscodes)), ist ein
Rest und zu ignorieren. Der nächste Fehlerzustand der neuen Version überschreibt
ihn (abgeleitet aus den Schreib-Zügen, nicht gegen einen Alt-Server gefahren).

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
so trägt `cdc.changes` mit den Feldern `origin` und `route_target` zwei
zusätzliche letzte Spalten), kann d-migrate die View nicht in-place ersetzen.
Der Lauf entfernt sie dann selbst und legt sie neu an; er meldet das:

```text
schema-rollout: Vorlauf - View-Signatur-Aenderung, DROP VIEW cdc.changes
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
  Sekunden, gemessen einmalig auf einer Testinstanz mit
  einer Zeile; eine Messung, nicht garantiert — auf einem
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
`Authorization: Bearer <token>`. Es gibt zwei Token-Klassen: lesend
(`CDC_API_TOKEN_READER`, `CDC_API_TOKENS_READER`) und administrativ/schreibend
(`CDC_API_TOKEN_ADMIN`, `CDC_API_TOKENS_ADMIN`), wobei ein Admin-Token die
lesende Klasse implizit mit abdeckt. Je Klasse können mehrere Token gleichzeitig
gültig sein (siehe [API-Token in zwei Neustarts wechseln](#api-token-in-zwei-neustarts-wechseln)).
Ein fehlender oder keinem konfigurierten Token entsprechender
Wert endet `401`, ein bekanntes Token mit unzureichender Klasse `403`. Die
Token-Klassen entscheiden an der API-Schicht, welche Fähigkeit erreichbar
ist; die Datenbankverbindung darunter trägt weiterhin die bei der Verdrahtung
fixierte Rolle.

**Fehlerantworten:** Eine Fehlerantwort trägt die Form
`{"error": "<Klartext>", "code": "<Meldungscode>"}`. `error` ist der Klartext,
`code` der Meldungscode der Ursache (Katalog unter [Meldungscodes](#meldungscodes));
`401` und `403` tragen kein Feld `code`. Ein Aufruf, den der Server als Eingabe
ablehnt (`400`, `404`), trägt einen Code mit der Ziffer `8`
(`PCF-E8050` bis `PCF-E8057`, `PCF-E8025` für eine an der Quelle fehlende
Tabelle); ein `500` trägt den Code der Ursache, den Rückfall `PCF-E7000`, wenn
keine nähere Ursache bekannt ist, und im Log steht dazu die Warnung `PCF-W4008`.
`GET /changes/stream` antwortet mit `503` und `PCF-E2001`, wenn der Feed-Container
ohne verdrahteten Broadcaster für den Stream läuft (eine unvollständige
Verdrahtung des Prozesses); das ist ein Fehler der Einrichtung, keiner der
Eingabe.
Die Statuscodes ändern sich nicht; ein Client, der nur `error` liest, arbeitet
unverändert weiter. Die SDK-Packages liefern ab Package-Version 0.6.0 `code` als
Eigenschaft des Fehlertyps (siehe die Fehlerform der gRPC-Verwaltungs-API weiter unten).
Gemessen am laufenden Feed-Container (`make test-integration`,
gedruckte Zeile `run-integration-tests: Fehlerzustand-Code-Beleg (HTTP und gRPC) belegt`):

```text
GET /changes                 (reader-Token, ohne source)
400 {"error":"source ist Pflichtfeld","code":"PCF-E8051"}

GET /changes                 (ohne Token)
401 {"error":"fehlender oder unbekannter Bearer-Token"}
```

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
| Diagnose lesen | `GET /diagnose` | `reader` |

Die lesenden Endpunkte sind mit dem Admin-Token ebenso erreichbar; mit dem
Reader-Token sind die administrativen Endpunkte nicht erreichbar (`403`).

**Tabelle aktivieren/deaktivieren wirkt sofort:** `POST /tables/enable`/`POST
/tables/disable` aktualisieren den laufenden Erfassungsprozess unmittelbar
mit der Antwort — ohne Neustart und ohne die Antrags-Queue des SQL-Zugriffswegs
(siehe [Tabelle live aktivieren](#tabelle-live-aktivieren)) zu durchlaufen:
derselbe Live-Reload-Vertrag, ein anderer Netzwerkzugriffsweg.

**Changes lesen:** `GET /changes?source=<quelle-id>` liefert persistierte
Änderungen einer Quelle — optional gefiltert über `schema`, `table` und
`target` (je einzeln oder in beliebiger Kombination, gemeinsam als
Konjunktion), eingegrenzt über `from` (inklusive) und `to`
(exklusive, jeweils ein `commit_position`-Wert) und begrenzt über `limit`
(≥ 1). Ohne `limit` liest der Aufruf unbegrenzt. Die Antwort trägt je
Änderung `commit_position`, `change_id`, `transaction_id`,
`source_table_id`, `schema`, `table`, `sequence`, `operation`,
`old_image`, `new_image`, `schema_version`, `committed_at` (RFC 3339,
UTC) und `origin` (`wal` oder `backfill`, als letztes Feld; eine ohne
dieses Feld gespeicherte Änderung liest als `wal`) — dieselbe Sicht wie
der SQL-Zugriff auf `cdc.changes`, ohne die Spalte `route_target`: das
Zustellziel steht nicht in der Antwort, es wählt nur den Filter `target`. Die
Reihenfolge ist deterministisch; die Fortsetzung liest ab
`from = <letzte gelieferte commit_position> + 1`, wenn das Lesen die letzte
Position vollständig erfasst hat. Ein `limit` schneidet Zeilen, nicht Positionen:
Enthält eine Commit-Position mehr Änderungen als `limit`, liefert dieses `from`
den Rest der Position nicht, und ein Lesen ab derselben Position liefert wieder
dieselben Zeilen. Den Bestandsabzug eines Backfills, der alle seine
Änderungen auf **eine** Position legt (siehe [Bestand als Backfill
überführen](#bestand-als-backfill-überführen)), lesen Sie deshalb ohne `limit`
oder über den Schlüsselvergleich im SQL-Zugriff (siehe [Änderungen
lesen](#änderungen-lesen)). Ein Aufruf ohne Treffer
endet `200` mit leerer Liste (`{"changes": []}`), nie `404`; ein Parameter
außerhalb der genannten Liste endet `400`, ebenso ein fehlendes `source`,
eine nicht lesbare Zahl, `from`/`to` unter 1, `limit` unter 1 und
`from > to`.

**Zustellziel wählen:** `target` liefert nur Changes, deren Zustellziel
(`route_target`, siehe [Routing-Regel konfigurieren](#routing-regel-konfigurieren))
genau dieser Name ist; ohne den Parameter oder mit leerem Wert gilt kein Filter.
Ein Name, den keine Change trägt — auch einer außerhalb des Alphabets des
Zielnamens —, ist bei sonst gültiger Anfrage kein `400`, sondern ein Filter ohne
Treffer: `200` mit `{"changes": []}`. Ein Fehler des Lese-Kontrakts (fehlendes
`source`, `limit` unter 1, `from > to`, Start- oder Endposition einer anderen
Quelle) hat Vorrang und endet mit demselben Fehler wie ohne `target`. Der
SQL-Zugriff `cdc.changes` ist davon ausgenommen: dort bestimmt der Aufrufer sein
Prädikat auf `route_target` selbst. Das Ziel ist Auswahl, kein Zugriffsschutz:
jedes Token der Klasse `reader` kann jedes Ziel wählen. *Ursprung:*
gemessen in einer Compose-Umgebung (PostgreSQL 18.6, Token
der Klasse `reader`): `?target=eu` lieferte die eine `eu`-Change, `?target=Gross`
`200` mit `{"changes":[]}`, `?target=Gross&limit=0` `400`, `?target=eu&schema=other`
`200` mit leerer Liste (Konjunktion). Ein Server ohne diese Funktion lehnt den
unbekannten Parameter an `GET /changes` und `GET /changes/stream` mit `400` ab
(*abgeleitet* aus der strengen Parameter-Menge, an einem
Server ohne diese Funktion nicht gefahren). Die Beispiel-Clients (Go, C#, Kotlin) nehmen
`target` über das Flag `-target` bzw. `--target` des Verbs `changes` entgegen
(`ARGS="-verb=changes -source=<quelle> -target=<ziel>"`), und die drei
SDK-Packages tragen den Parameter `target` an `ReadChangesAsync`
(`PgChangeFeed.Client`, als letzter Parameter, per Name übergeben),
`readChanges` (`pgchangefeed-kotlin`) und `read_changes` (`pgchangefeed`);
ohne Angabe bleibt der Aufruf unverändert.

**Diagnose lesen:** `GET /diagnose?source=<quelle-id>` liefert denselben
Bericht wie [Diagnose ausführen](#diagnose-ausführen) — Betriebsstatus und
Fehlerzustand, CDC-Abstand, Verarbeitungsrückstand je Consumer, blockierender
Consumer, Speicherverbrauch und der zuletzt beantragte Backfill-Run je
Tabelle — strukturiert statt als Text. `source` ist Pflicht, ein Parameter
außerhalb dieser einen Angabe endet `400`, ebenso wie bei `GET /changes`.
`error_code` trägt den Meldungscode des Fehlerzustands (siehe
[Meldungscodes](#meldungscodes)) und ist wie `error_class` im Normalbetrieb `null`.
Fehlende Werte stehen als `null` (`heartbeat_age_seconds`/`error_class`/`error_code`
gemeinsam bei fehlendem Lebenszeichen, `retention_blocker` ohne Blocker,
`lag`/`estimated_rows` bei unbekanntem Rückstand bzw. unbekannter Schätzung);
`consumer_lags`/`backfill` sind ohne Treffer leere, gesetzte Listen. Ein
Lesefehler an einer der Diagnose-Views endet `500`.

**Zustellsemantik:** Diese Fähigkeiten sind synchrone Anfrage/Antwort — die
Antwort trägt das Ergebnis des Aufrufs, es gibt keine Warteschlange
dazwischen. Die Ausnahme ist der Live-Stream auf `GET /changes/stream` (siehe
unten), der die Verbindung offen hält — `GET /changes` ist demgegenüber die
nicht streamende Form desselben Gegenstands.

**Beispiele:** Jede Sprache deckt alle zehn Fähigkeiten der Tabelle oben über
ein Verb-Flag ab (Default `tables`); Adresse und Token liest jedes
Beispiel aus `CDC_HTTP_ADDR`/`CDC_API_TOKEN_READER`/`CDC_API_TOKEN_ADMIN` und
lässt sich per Flag übersteuern. `Diagnose lesen` (`GET /diagnose`) ist davon
nicht abgedeckt — anders als beim elften gRPC-RPC `Diagnose` (siehe
[Zugriff über die gRPC-Verwaltungs-API](#zugriff-über-die-grpc-verwaltungs-api)),
den das Go-gRPC-Beispiel abdeckt; die HTTP-Beispiel-Clients aller drei
Sprachen, die C#-/Kotlin-gRPC-Beispiele und die drei SDK-Packages decken
`Diagnose` nicht ab.

- **Go:** `examples/http-client` — Container-Aufruf über
  `make example-run-go SURFACE=http ARGS="-verb=<verb> ..."`
  (baut bei Bedarf `pg-change-feed-examples:go` aus `examples/Dockerfile`)
- **C#:** `examples/csharp/http-client` — Container-Aufruf über
  `make example-run-csharp SURFACE=http ARGS="--verb=<verb> ..."`
  (startet das mit `make examples-csharp` gebaute Image)
- **Kotlin:** `examples/kotlin/http-client` — Container-Aufruf über
  `make example-run-kotlin SURFACE=http ARGS="--verb=<verb> ..."`
  (startet das mit `make examples-kotlin` gebaute Image)

`<verb>` ist eine der zehn Fähigkeiten der Tabelle oben (z. B. `tables`,
`changes`, `register-consumer`, `enable-table`); die je Verb nötigen
zusätzlichen Flags (z. B. `-source`/`--source`, `-consumer-id`/
`--consumer-id`) entsprechen den Request-Feldern der jeweiligen Zeile. Ein
Aufruf mit `-verb=changes`/`--verb=changes` zeigt auch die Wirkung einer
zuvor über [Transformationsregel konfigurieren](#transformationsregel-konfigurieren)
angewendeten Regel — etwa nach einem Lauf von `make
example-transformation-demo`.

**SDK:** .NET-Anwendungen können statt der Beispiele das offizielle
NuGet-Package `PgChangeFeed.Client` einbinden
(`dotnet add package PgChangeFeed.Client`) — `PgChangeFeedHttpClient` deckt
alle zehn Fähigkeiten dieser Zugriffs-Oberfläche ab (die neun in der
Tabelle oben plus `GET /changes`) mit typisierten Requests/Responses und
einer typisierten Fehlerklasse für `400`/`401`/`403`/`404`/`500`, statt den
Draht-Vertrag selbst zu implementieren; siehe `sdks/csharp/README.md`.

Python-Anwendungen können statt der Beispiele das offizielle PyPI-Package
`pgchangefeed` einbinden
(`pip install pgchangefeed`) — `PgChangeFeedHttpClient` deckt dieselben zehn
Fähigkeiten dieser Zugriffs-Oberfläche ab (die neun in der Tabelle oben
plus `GET /changes`) mit typisierten Requests/Responses (`dataclasses`) und
einer typisierten Fehlerklasse für `400`/`401`/`403`/`404`/`500`; derselbe
Package trägt außerdem den gRPC-Change-Stream (siehe unten, „Zugriff über
den gRPC-Change-Stream"), den SSE-Stream (siehe „Zugriff über
Server-Sent-Events") und den NATS-Vollinhalts-Stream (siehe „Zugriff über
den NATS-Vollinhalts-Stream"). Siehe `sdks/python/README.md`.

Kotlin/JVM-Anwendungen können statt der Beispiele das offizielle
Gradle-/Maven-Package `pgchangefeed-kotlin` einbinden
(Koordinate `io.github.pt9912:pgchangefeed-kotlin`) —
`PgChangeFeedHttpClient` deckt dieselben zehn Fähigkeiten dieser
Zugriffs-Oberfläche ab (die neun in der Tabelle oben plus `GET /changes`)
mit typisierten Requests/Responses und einer versiegelten
(`sealed class`) Fehlerklasse für `400`/`401`/`403`/`404`/`500`; derselbe
Package trägt außerdem den gRPC-Change-Stream (siehe unten, „Zugriff über
den gRPC-Change-Stream"), den SSE-Stream (siehe „Zugriff über
Server-Sent-Events") und den NATS-Vollinhalts-Stream (siehe „Zugriff über
den NATS-Vollinhalts-Stream"). Anders als NuGet/PyPI wird dieses Package
nicht über die Registry der Sprache vertrieben, sondern über zwei Ziele: den
**Cloudsmith-Repository-Pfad** `https://dl.cloudsmith.io/public/pt9912/pg-change-feed/maven/`
(anonym lesbar — **ohne Konto und ohne Token**) und **GitHub
Packages** (`https://maven.pkg.github.com/pt9912/pg-change-feed`). GitHub Packages verlangt **immer** eine
Authentifizierung zum Lesen, auch für ein öffentliches Package: ein
GitHub-Konto und ein klassischer Personal Access Token (PAT) mit dem Scope
`read:packages` sind dort Voraussetzung, unabhängig davon, ob das SDK
öffentlich und quelloffen ist. Der einfachere Weg ist Cloudsmith. Siehe
`sdks/kotlin/pgchangefeed-kotlin/README.md` für den vollständigen
Installationsweg samt Gradle-Repository-Konfiguration.

In allen drei Packages trägt jede über `GET /changes` gelesene Änderung das
Feld `origin` (`wal` oder `backfill`); eine Antwort ohne dieses Feld oder mit
JSON-`null` liest das Package als `wal`, jeden anderen Wert des Servers (auch
einen leeren oder unbekannten) gibt es unverändert weiter. Die Live-Wege
(gRPC, SSE, NATS-Vollinhalt) tragen kein `origin`.

### API-Token in zwei Neustarts wechseln

Die Token der HTTP- und der gRPC-Schnittstelle gelten je Klasse (lesend,
administrativ) als Menge: gültig ist jedes Token der Menge, auf beiden
Schnittstellen gleich. Die Menge einer Klasse besteht aus dem einen Token der
Variable `CDC_API_TOKEN_READER` bzw. `CDC_API_TOKEN_ADMIN` und den Token der
Liste in `CDC_API_TOKENS_READER` bzw. `CDC_API_TOKENS_ADMIN` (kommagetrennt,
siehe [Umgebungsvariablen des Feed-Containers](#umgebungsvariablen-des-feed-containers)).
So lässt sich ein Token ersetzen, ohne dass ein Client in der Zwischenzeit
abgewiesen wird. Der Server liest die Token beim Start; ein Neuladen zur
Laufzeit gibt es nicht, der Wechsel besteht aus zwei Neustarts:

1. **Neues Token ergänzen.** Tragen Sie das neue Token in die Liste der Klasse
   ein (`CDC_API_TOKENS_READER=neues-token`) und lassen Sie das alte Token
   unverändert. Starten Sie den Container neu: beide Token werden bedient.
2. **Clients umstellen.** Stellen Sie jeden Client auf das neue Token um.
3. **Altes Token entfernen.** Nehmen Sie das alte Token aus der Konfiguration
   (die Variable leeren oder löschen) und starten Sie den Container erneut.
   Ein Aufruf mit dem alten Token endet danach auf der HTTP-API mit `401` und
   auf der gRPC-API mit `Unauthenticated`; das neue Token wird weiter bedient.

Regeln für die Werte:

- Die Variablen `CDC_API_TOKEN_READER` und `CDC_API_TOKEN_ADMIN` tragen genau
  ein Token; ein Komma gehört dort zum Token. Die Liste trennt am Komma, ein
  Token mit Komma lässt sich in der Liste nicht ausdrücken.
- Ein leeres Element der Liste (`a,,b`, ein Komma am Anfang oder am Ende) und
  ein Element mit Leerraum (Leerzeichen, Tabulator, Zeilenumbruch) verhindern
  den Start: der Container endet mit der Fehlerklasse `configuration` und dem
  Meldungscode `PCF-E2008`; die Zeile nennt die Variable und die Nummer des
  Elements, nie den Wert.
- Eine leere Variable gilt als nicht gesetzt.
- Ein leeres Token ist in keiner Klasse gültig. Steht derselbe Wert in der
  lesenden und der administrativen Klasse, gilt er als administrativ.
- Alle vier Variablen gehören zu den Zugangsdaten: sie stehen ausschließlich in
  der Umgebung, nicht in der Konfigurationsdatei.

Das Token des NATS-Vollinhalts-Zustellwegs (`CDC_NATS_STREAM_TOKEN`) gehört
nicht dazu: der NATS-Server vergibt und wechselt es.

### Schnittstellen mit TLS verschlüsseln

Die HTTP-/JSON-API (einschließlich des SSE-Endpunkts `GET /changes/stream`)
und die gRPC-Schnittstelle (Change-Stream und Verwaltungs-API) bedienen ohne
Konfiguration unverschlüsselt. Mit einem TLS-Paar bedienen beide Schnittstellen
verschlüsselt. Dasselbe Paar gilt für beide.

**Einrichten:** Legen Sie Zertifikat und privaten Schlüssel als PEM-Dateien an
(das Zertifikat mit dem Namen, unter dem Ihre Clients den Server erreichen),
binden Sie beide lesend in den Container ein und setzen Sie beide Variablen —
oder die beiden Datei-Felder `tls_cert_file` und `tls_key_file`:

```yaml
    environment:
      CDC_TLS_CERT_FILE: /etc/cdc/tls/server.pem
      CDC_TLS_KEY_FILE: /etc/cdc/tls/server-key.pem
    volumes:
      - ./tls:/etc/cdc/tls:ro
```

Der Prozess im Container läuft als unprivilegierter Benutzer; beide Dateien
müssen für ihn lesbar sein.

**Was gilt:**

- Beide Schnittstellen sprechen mindestens TLS 1.2. Ein Client, der höchstens
  TLS 1.1 spricht, kommt nicht zustande (gemessen: der Server beendet den
  Verbindungsaufbau mit „protocol version not supported“).
- Auf derselben Adresse gibt es keinen Klartext. Ein Klartext-Aufruf wird nicht
  bedient, auch nicht mit gültigem Token: ein HTTP-Aufruf über `http://`
  bekommt den Status `400` mit dem Text „Client sent an HTTP request to an
  HTTPS server“ (gemessen), ein gRPC-Aufruf im Klartext endet mit dem Status
  `Unavailable` (gemessen).
- Der Server verlangt kein Client-Zertifikat. Die Authentifizierung bleibt das
  Token (siehe [Zugriff über die HTTP-/JSON-API](#zugriff-über-die-http-json-api)).
- Der Server liest Zertifikat und Schlüssel beim Start. Ein Zertifikatswechsel
  ist ein Neustart; ein Neuladen zur Laufzeit gibt es nicht.
- Der Server prüft beim Start weder das Ablaufdatum noch den Namen des
  Zertifikats: ein abgelaufenes oder für einen anderen Namen ausgestelltes
  Zertifikat startet den Server, und der Client lehnt die Verbindung ab. Prüfen
  Sie Gültigkeit und Namen vor dem Einsatz.
- Die Verbindung zu NATS und zu PostgreSQL bleibt unberührt: TLS für NATS stellen
  Sie in der NATS-Konfiguration ein, TLS für PostgreSQL über die Parameter der
  Verbindungs-URL.
- `--healthcheck` (der Healthcheck des Containers) liest den Herzschlag in der
  Datenbank und erreicht die beiden Schnittstellen nicht; er läuft mit und ohne
  TLS gleich (gemessen: Ausgang 0 und Docker-Zustand `healthy` mit gesetztem
  Paar).

**Start-Fehler:** Beide Fälle beenden den Prozess vor dem Aufbau der
Datenbank-Verbindungen, ohne angelegten Replication-Slot, mit der
Fehlerklasse `configuration`:

| Fall | Meldungscode | Prozessausgang |
|---|---|---|
| nur eine der beiden Angaben ist gesetzt (Umgebung und Datei zusammen gelesen) | `PCF-E2009` | 2 |
| eine Datei ist nicht lesbar, enthält keine gültigen PEM-Daten, oder Zertifikat und Schlüssel gehören nicht zusammen | `PCF-E2010` | 1 |

Die Fehlerzeile nennt die Variablen bzw. die Pfade und die Ursache, nie den
Inhalt des Schlüssels.

**Clients:** Die gRPC-Beispielprogramme in Go, C# und Kotlin verbinden über TLS,
wenn ihnen ein Vertrauensanker genannt wird: das Flag `-ca-file <Pfad>` (Go) bzw.
`--ca-file <Pfad>` (C#, Kotlin) oder die Umgebungsvariable `CDC_TLS_CA_FILE`
nennt das Zertifikat des Servers als PEM-Datei. Das Beispiel prüft dann Kette
und Namen des Serverzertifikats gegen diese Datei und überspringt die Prüfung
nie; ohne Angabe verbindet es im Klartext wie bisher. Eine nicht lesbare Datei
und eine Datei ohne PEM-Zertifikat enden mit Ausgang 2 und einer Fehlermeldung,
bevor eine Verbindung entsteht. Die Datei muss im Container des Beispiels
liegen, zum Beispiel für das Go-Beispiel mit dem Stream als Verb:

```text
docker run --rm --network cdc-examples --env-file examples/.env \
  -v <Pfad zum Zertifikat>:/tls/ca.pem:ro -e CDC_TLS_CA_FILE=/tls/ca.pem \
  pg-change-feed-examples:go-grpc
```

Das Image `pg-change-feed-examples:go-grpc` entsteht mit `make example-run-go
SURFACE=grpc`, die Images `pg-change-feed-examples:csharp-grpc` und
`pg-change-feed-examples:kotlin-grpc` mit `make examples-csharp` bzw.
`make examples-kotlin`. Gemessen gegen einen Feed-Container mit TLS-Paar: alle
drei Beispiele empfangen im Verb `stream` mit `CDC_TLS_CA_FILE` eine danach
committete Änderung (dieselbe Kennung wie in `cdc.changes`); ohne Angabe endet
jedes Beispiel mit Ausgang 1 und einer Fehlermeldung zum nicht aufgebauten
Kanal, ohne eine Antwort des Servers. Die übrigen Verben der Beispiele sind
nicht gegen einen TLS-Server gefahren. Die Optionen der drei Client-Pakete (C#,
Kotlin, Python) bieten keine TLS-Einstellung; für die Client-Pakete und für
HTTP- und SSE-Clients gegen einen solchen Server ist hier nichts zugesagt.

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

**Filterung:** Der Request `StreamChangesRequest` trägt drei optionale,
unabhängig setzbare Felder `schema`/`table`/`target` (Feldnummern 1/2/3). Ein
gesetztes `schema` ohne `table` liefert alle Tabellen dieses Schemas, ein
gesetztes `table` ohne `schema` jede Tabelle dieses Namens unabhängig vom
Schema, beide gesetzt liefert exakt eine Tabelle. Ein gesetztes `target`
liefert nur Changes, deren Zustellziel (`route_target`, siehe
[Routing-Regel konfigurieren](#routing-regel-konfigurieren)) dieser Name ist, als
Konjunktion mit `schema`/`table`; ein Name, den keine Change trägt — auch einer
außerhalb des Alphabets des Zielnamens —, liefert keine Nachricht und keinen
Fehler. Bleiben alle drei Felder leer, liefert
der Stream jeden Change aller aktivierten Tabellen. Die Prüfung läuft
serverseitig, bevor eine nicht passende Change über das Netz geht. Ein Server
ohne das Feld `target` ignoriert es und liefert ungefilterte Changes (proto3
verwirft unbekannte Felder: übernommen aus den Tests der Lesewege,
an einem Server ohne das Feld *abgeleitet*). Die Nachricht
`Change` trägt das Ziel nicht. *Ursprung:* übernommen aus dem E2E-Lauf von `make
test-integration`: ein gRPC-Client
mit Ziel empfing von einer festen Menge gemischter Changes genau die Changes
seines Ziels und im Ruhefenster von 15 s keine Change eines anderen Ziels oder
ohne Ziel. Go, C# und Kotlin nehmen den Filter über `-schema`/`-table`/`-target`
bzw. `--schema`/`--table`/`--target` entgegen; das NuGet-Package
`PgChangeFeed.Client` nimmt ihn ebenfalls entgegen (`StreamChangesAsync(schema,
table, cancellationToken, target)`, alle drei Parameter optional, `target` als
letzter Parameter per Name übergeben), ebenso das PyPI-Package `pgchangefeed`
(`stream_changes(timeout, schema, table, target)`, alle drei Parameter
optional) und das Gradle-/Maven-Package `pgchangefeed-kotlin`
(`streamChanges(schema, table, target)`, alle drei Parameter optional, Default
`null`). Der RPC `ReadChanges` der Verwaltungs-API trägt `target` als Feld des
Requests (siehe [Zugriff über die gRPC-Verwaltungs-API](#zugriff-über-die-grpc-verwaltungs-api)).

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
und lässt sich per Flag übersteuern. Gegen einen Server mit TLS-Paar nimmt jedes
Beispiel das Zertifikat über `-ca-file`/`--ca-file` oder `CDC_TLS_CA_FILE`
entgegen (siehe [Schnittstellen mit TLS
verschlüsseln](#schnittstellen-mit-tls-verschlüsseln)).

- **Go:** `examples/grpc-client` — Container-Aufruf über
  `make example-run-go SURFACE=grpc` (baut bei Bedarf
  `pg-change-feed-examples:go-grpc` aus `examples/Dockerfile`); das
  Default-Verb `stream` nimmt den optionalen Filter über
  `ARGS="-schema=<schema> -table=<tabelle> -target=<ziel>"` entgegen (alle drei
  leer liefert jeden Change aller aktivierten Tabellen); dasselbe
  Programm deckt
  über `-verb` zusätzlich die elf RPCs der [gRPC-Verwaltungs-API](#zugriff-über-die-grpc-verwaltungs-api)
  ab
- **C#:** `examples/csharp/grpc-client` — Container-Aufruf über
  `make example-run-csharp SURFACE=grpc` (startet das mit
  `make examples-csharp` gebaute Image; die beiden Stubs entstehen im Bau aus
  den `.proto`-Dateien, über einen zusätzlichen, benannten Bau-Kontext gelesen);
  das Default-Verb `stream` nimmt den optionalen
  `--schema`/`--table`/`--target`-Filter entgegen; dasselbe Programm deckt
  über `--verb` zusätzlich die elf RPCs der
  [gRPC-Verwaltungs-API](#zugriff-über-die-grpc-verwaltungs-api) ab
- **Kotlin:** `examples/kotlin/grpc-client` — Container-Aufruf über
  `make example-run-kotlin SURFACE=grpc` (startet das mit
  `make examples-kotlin` gebaute Image; derselbe Stub-im-Bau-Mechanismus wie
  beim C#-Client, übertragen auf die Kotlin-Werkzeugkette); das
  Default-Verb `stream` nimmt den optionalen
  `--schema`/`--table`/`--target`-Filter entgegen; dasselbe
  Programm deckt über `--verb` zusätzlich die
  elf RPCs der [gRPC-Verwaltungs-API](#zugriff-über-die-grpc-verwaltungs-api)
  ab

**SDK:** .NET-Anwendungen können statt der Beispiele das offizielle
NuGet-Package `PgChangeFeed.Client` einbinden
(`dotnet add package PgChangeFeed.Client`) — `PgChangeFeedGrpcClient.StreamChangesAsync`
öffnet den `ChangeStream/StreamChanges`-RPC und liefert ein
`IAsyncEnumerable<Change>` mit allen zehn Feldern der Tabelle oben; drei
optionale Parameter `schema`/`table`/`target` tragen den Filter über Schema,
Tabelle und Zustellziel wie oben beschrieben, alle `null` (der
Default) liefert jeden Change (`target` ist der letzte Parameter und
wird per Name übergeben). Das Bearer-Token landet im `authorization`-Metadata-Eintrag, ein
fehlendes oder ungültiges Token endet den Aufruf mit gRPC-Status
`Unauthenticated`, statt den Draht-Vertrag selbst zu implementieren; siehe
`sdks/csharp/README.md`.

Kotlin/JVM-Anwendungen können statt des Beispiels dasselbe offizielle
Gradle-/Maven-Package `pgchangefeed-kotlin` einbinden
(Koordinate `io.github.pt9912:pgchangefeed-kotlin`) —
`PgChangeFeedGrpcClient.streamChanges()` öffnet denselben
`ChangeStream/StreamChanges`-RPC und liefert ein
`kotlinx.coroutines.flow.Flow<Change>` mit allen zehn Feldern der Tabelle
oben; drei optionale Parameter `schema`/`table`/`target` tragen den Filter über
Schema, Tabelle und Zustellziel wie oben beschrieben, alle `null`
(der Default) liefert jeden Change. Das Bearer-Token landet im `authorization`-Metadata-Eintrag, ein
fehlendes oder ungültiges Token endet den Aufruf mit gRPC-Status
`Unauthenticated`, statt den Draht-Vertrag selbst zu implementieren. Derselbe
Package trägt außerdem den SSE-Stream (siehe „Zugriff über
Server-Sent-Events") und den NATS-Vollinhalts-Stream (siehe „Zugriff über
den NATS-Vollinhalts-Stream"). Wie beim HTTP-API-Zugriff oben ist der Bezug
über Cloudsmith
(`https://dl.cloudsmith.io/public/pt9912/pg-change-feed/maven/`) ohne Konto
und ohne Token möglich; der Bezug über **GitHub Packages**
(`https://maven.pkg.github.com/pt9912/pg-change-feed`) verlangt dagegen
immer eine Authentifizierung: ein GitHub-Konto und ein klassischer Personal
Access Token (PAT) mit dem Scope `read:packages`. Siehe `sdks/kotlin/pgchangefeed-kotlin/README.md`.

Python-Anwendungen können statt des Beispiels das offizielle
PyPI-Package `pgchangefeed` einbinden
(`pip install pgchangefeed`) — `PgChangeFeedGrpcClient.stream_changes()`
öffnet denselben `ChangeStream/StreamChanges`-RPC und liefert einen
Iterator über die generierten `Change`-Nachrichten mit allen zehn Feldern
der Tabelle oben; das Bearer-Token landet im `authorization`-Metadata-Eintrag,
ein fehlendes oder ungültiges Token endet den Aufruf mit gRPC-Status
`Unauthenticated`, statt den Draht-Vertrag selbst zu implementieren
(`timeout` ist der Gesamtfriestempel des Aufrufs in Sekunden, `None` =
unbegrenzt); drei optionale Parameter `schema`/`table`/`target` tragen den
Filter über Schema, Tabelle und Zustellziel wie oben beschrieben,
alle `None` (der Default) liefert jeden Change. Siehe
`sdks/python/README.md`.

### Zugriff über die gRPC-Verwaltungs-API

Dieselben Verwaltungsfähigkeiten wie die HTTP-/JSON-API (siehe [Zugriff über
die HTTP-/JSON-API](#zugriff-über-die-http-json-api)) sind zusätzlich über
gRPC erreichbar — jede Fähigkeit ruft dieselbe Domänenlogik über denselben
Inbound Use Case auf wie ihr HTTP-Äquivalent und ist fachlich gleichwertig,
kein Zweitpfad.

**Erreichbarkeit:** aktiv, sobald `CDC_GRPC_ADDR` gesetzt ist — dieselbe
Horch-Adresse wie der [gRPC-Change-Stream](#zugriff-über-den-grpc-change-stream)
oben; ein gesetzter Wert aktiviert beide Services gemeinsam auf demselben
Server/Port. Ungesetzt bleibt auch diese Fläche vollständig deaktiviert.

**Authentifizierung:** derselbe Mechanismus wie beim gRPC-Change-Stream —
Bearer-Token im gRPC-Metadata-Eintrag `authorization` in der Form
`Bearer <token>`, dieselben zwei Token-Klassen (`CDC_API_TOKEN_READER`,
`CDC_API_TOKEN_ADMIN` mit den Listen `CDC_API_TOKENS_READER`,
`CDC_API_TOKENS_ADMIN`) wie bei der HTTP-API. Ein fehlender oder keinem
konfigurierten Token entsprechender Wert endet mit dem gRPC-Status
`Unauthenticated`, ein bekanntes Token mit unzureichender Klasse mit
`PermissionDenied`.

Der Dienst `Administration` (Paket `cdc.administration.v1`) trägt elf
unäre RPCs:

| Fähigkeit | RPC | Rechtsklasse |
|---|---|---|
| Consumer registrieren | `RegisterConsumer` | `admin` |
| Position bestätigen | `AcknowledgeConsumer` | `admin` |
| Consumer-Position lesen | `GetConsumerPosition` | `reader` |
| Consumer entfernen | `RemoveConsumer` | `admin` |
| Tabelle aktivieren | `EnableTable` | `admin` |
| Tabelle deaktivieren | `DisableTable` | `admin` |
| Tabellen-Status | `GetTableStatus` | `reader` |
| Tabellen auflisten | `ListTables` | `reader` |
| Aufbewahrung auslösen | `RunRetention` | `admin` |
| Changes lesen (begrenzter Bereich) | `ReadChanges` | `reader` |
| Diagnose lesen | `Diagnose` | `reader` |

Die lesenden RPCs sind mit dem Admin-Token ebenso erreichbar; mit dem
Reader-Token sind die administrativen RPCs nicht erreichbar
(`PermissionDenied`). `EnableTable`/`DisableTable` wirken wie ihr
HTTP-Äquivalent sofort auf den laufenden Erfassungsprozess, ohne Neustart
(siehe [Zugriff über die HTTP-/JSON-API](#zugriff-über-die-http-json-api)).
Nachrichtenfelder entsprechen 1:1 den Request-/
Response-Feldern der gleichnamigen HTTP-Fähigkeit (Tabelle unter [Zugriff
über die HTTP-/JSON-API](#zugriff-über-die-http-json-api)); `ListTables`
liefert `tables` und `retained` als je eine Liste einer `SourceTable`
genannten Nachricht (`table_id`, `source`, `schema`, `table`). `ReadChanges`
liest einen begrenzten Bereich persistierter Änderungen — dieselbe Filter-
und Bereichs-Semantik wie `GET /changes` (Quelle Pflicht, `schema`/`table`/`target`
optional und unabhängig, `from`/`to` als Positions-Bereich `[from, to)`,
`limit`; `target` ist `string`, Feldnummer 7: leer heißt kein Filter, gesetzt
liefert nur Changes mit diesem Zustellziel, als Konjunktion mit `schema`/`table`),
die Antwort trägt `changes` als Liste einer `ChangeRecord`
genannten Nachricht mit denselben dreizehn Feldern wie die HTTP-Antwort
(`commit_position`, `change_id`, `transaction_id`, `source_table_id`,
`schema`, `table`, `sequence`, `operation`, `old_image`, `new_image`,
`schema_version`, `committed_at`, `origin`; das Zustellziel gehört nicht dazu);
kein Treffer liefert eine leere, gesetzte Liste statt eines Fehlers. Ein
`target`, das keine Change trägt — auch eines außerhalb des Alphabets des
Zielnamens —, liefert bei sonst gültiger Anfrage eine leere Liste, keinen Fehler;
ein Fehler des Lese-Kontrakts (`limit`, Bereich, Quelle) hat Vorrang und endet
mit demselben Code wie ohne `target`. *Ursprung:* die Auswahl nach Ziel über
`ReadChanges` ist übernommen aus dem E2E-Lauf von `make test-integration`
(die Kennungen entsprechen der SQL-Auswahl); den leeren Treffer für einen Namen außerhalb des Alphabets und den
Vorrang des Lese-Kontrakts belegen Unit-Tests des gemeinsamen Use Cases, am
gRPC-Weg des Systems nicht gefahren (den HTTP-Weg desselben Use Cases siehe
oben). Das Ziel ist
Auswahl, kein Zugriffsschutz. Ein Server ohne das Feld `target` ignoriert es und liefert
ungefilterte Changes (proto3 verwirft unbekannte Felder; übernommen aus den Tests
der Lesewege, an einem Server ohne das Feld *abgeleitet*). Der [gRPC-Change-Stream](#zugriff-über-den-grpc-change-stream)
oben bleibt der Zugriffsweg für neue, laufend eintreffende Changes.
`Diagnose` liefert denselben Bericht wie [Diagnose ausführen](#diagnose-ausführen)
und `GET /diagnose` (siehe [Zugriff über die HTTP-/JSON-API](#zugriff-über-die-http-json-api))
— dieselben sechs Signalgruppen, strukturiert über die Nachrichten
`HeartbeatStatus`, `ConsumerLag`, `RetentionBlocker` und
`BackfillTableStatus`; ein `known`/`present`/`*_known`-Feld auf `false`
trägt dieselben Abwesenheits-Fälle wie der CLI-Text (kein Lebenszeichen,
kein Blocker, unbekannte Schätzung/unbekannter Rückstand). `HeartbeatStatus`
trägt neben `error_class` das Feld `error_code` mit dem Meldungscode des
Fehlerzustands (siehe [Meldungscodes](#meldungscodes)); beide Felder sind im
Normalbetrieb leer.

**Fehlerform:** Ein gRPC-Status ersetzt den jeweiligen HTTP-Statuscode:

| HTTP-Status (Ursache) | gRPC-Code |
|---|---|
| `400` — Domänen-Invariante verletzt | `InvalidArgument` |
| `401` — fehlender/unbekannter Bearer-Token | `Unauthenticated` |
| `403` — Rechtsklasse unzureichend | `PermissionDenied` |
| `404` — Tabelle nicht gefunden | `NotFound` |
| `500` — unerwarteter interner Fehler | `Internal` |

Ein Status mit `InvalidArgument`, `NotFound` oder `Internal` trägt zusätzlich das
Statusdetail `google.rpc.ErrorInfo`: `reason` ist der Meldungscode der Ursache
(Katalog unter [Meldungscodes](#meldungscodes)), `domain` ist `pg-change-feed`.
`Unauthenticated` und `PermissionDenied` tragen kein Detail. Statuscode und
Statustext bleiben unverändert; ein Client, der nur den Statuscode liest,
arbeitet unverändert weiter. Ein Go-Client liest das Detail über
`status.Convert(err).Details()`. Die drei SDK-Packages (C#, Kotlin, Python)
tragen ab Package-Version 0.6.0 den Code als Eigenschaft ihrer Fehlertypen (`MessageCode`, `messageCode`,
`message_code`; leer, wenn der Server keinen Code sendet): beim HTTP- und beim
SSE-Client aus dem Feld `code`, beim gRPC-Verwaltungs-Client aus dem Statusdetail;
die gRPC-Stream-Clients lassen den Fehler der gRPC-Bibliothek unverändert durch,
das Statusdetail bleibt dort über deren eigene Schnittstelle lesbar. Die
Beispiele dieses Projekts werten es nicht aus.

**Zustellsemantik:** wie die HTTP-API — synchrone Anfrage/Antwort je RPC,
keine Warteschlange dazwischen.

**Beispiele:** Das Go-Beispiel `examples/grpc-client`, das C#-Beispiel
`examples/csharp/grpc-client` und das Kotlin-Beispiel
`examples/kotlin/grpc-client` (siehe [gRPC-Change-Stream](#zugriff-über-den-grpc-change-stream)
oben) decken über dasselbe `-verb`-/`--verb`-Flag alle elf Fähigkeiten der
Tabelle oben ab (Default `stream`); Adresse und Token lesen alle drei
aus `CDC_GRPC_ADDR`/`CDC_API_TOKEN_READER`/`CDC_API_TOKEN_ADMIN` und lassen
sich per Flag übersteuern, z. B.
`make example-run-go SURFACE=grpc ARGS="-verb=list-tables -source=<quelle> -publication=<publication>"`,
`make example-run-csharp SURFACE=grpc ARGS="--verb=list-tables --source=<quelle> --publication=<publication>"`
bzw.
`make example-run-kotlin SURFACE=grpc ARGS="--verb=list-tables --source=<quelle> --publication=<publication>"`.
Das Verb `read-changes` nimmt zusätzlich `-target`/`--target` entgegen (leer
ist kein Filter).

**SDK:** .NET-Anwendungen können statt der Beispiele das offizielle
NuGet-Package `PgChangeFeed.Client` einbinden
(`dotnet add package PgChangeFeed.Client`) — `PgChangeFeedAdministrationClient`
trägt alle elf RPCs der Tabelle oben als eigene async-Methode
(`RegisterConsumerAsync`, `AcknowledgeConsumerAsync`,
`GetConsumerPositionAsync`, `RemoveConsumerAsync`, `EnableTableAsync`,
`DisableTableAsync`, `GetTableStatusAsync`, `ListTablesAsync`,
`RunRetentionAsync`, `ReadChangesAsync`, `DiagnoseAsync`); Requests und
Responses sind die generierten Protobuf-Nachrichten unverändert, kein
eigener DTO-Layer. Ein nicht-`OK`-Status wird zu einer typisierten
`PgChangeFeedGrpcException`-Unterklasse je gRPC-Code der Fehlerform-Tabelle
oben, statt den Draht-Vertrag selbst zu implementieren; siehe
`sdks/csharp/README.md`.

Python-Anwendungen können statt der Beispiele das offizielle PyPI-Package
`pgchangefeed` einbinden (`pip install pgchangefeed`) —
`PgChangeFeedAdministrationClient` trägt alle elf RPCs der Tabelle oben als
eigene Methode (`register_consumer`, `acknowledge_consumer`,
`get_consumer_position`, `remove_consumer`, `enable_table`, `disable_table`,
`get_table_status`, `list_tables`, `run_retention`, `read_changes`,
`diagnose`); Requests und Responses sind die generierten Protobuf-Nachrichten
unverändert, kein eigener DTO-Layer. Ein nicht-`OK`-Status wird zu einer
typisierten `PgChangeFeedGrpcError`-Unterklasse je gRPC-Code der
Fehlerform-Tabelle oben, statt den Draht-Vertrag selbst zu implementieren;
siehe `sdks/python/README.md`.

Kotlin/JVM-Anwendungen können statt der Beispiele dasselbe offizielle
Gradle-/Maven-Package `pgchangefeed-kotlin` einbinden
(Koordinate `io.github.pt9912:pgchangefeed-kotlin`) —
`PgChangeFeedAdministrationClient` trägt alle elf RPCs der Tabelle oben als
eigene `suspend fun`-Methode (`registerConsumer`, `acknowledgeConsumer`,
`getConsumerPosition`, `removeConsumer`, `enableTable`, `disableTable`,
`getTableStatus`, `listTables`, `runRetention`, `readChanges`, `diagnose`);
Requests und Responses sind die generierten Protobuf-Nachrichten unverändert,
kein eigener DTO-Layer. Ein nicht-`OK`-Status wird zu einer versiegelten
(`sealed class`) `PgChangeFeedGrpcException`-Unterklasse je gRPC-Code der
Fehlerform-Tabelle oben, statt den Draht-Vertrag selbst zu implementieren;
siehe `sdks/kotlin/pgchangefeed-kotlin/README.md`. Die gRPC-Verwaltungs-API ist
damit in allen drei Sprach-SDKs (C#, Python, Kotlin) abgedeckt.

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

**Filterung:** Drei optionale, unabhängig setzbare Query-Parameter
`schema`/`table`/`target` — dieselben Feldnamen und dieselbe Kombinatorik wie
beim gRPC-Stream oben und bei `GET /changes` (siehe
[Änderungen lesen](#änderungen-lesen)); `target` wählt das Zustellziel (siehe
[Routing-Regel konfigurieren](#routing-regel-konfigurieren)) und liefert nur
Changes dieses Ziels, als Konjunktion mit `schema`/`table`. Ein Name, den keine
Change trägt — auch einer außerhalb des Alphabets des Zielnamens —, liefert
keinen Event und keinen Fehler. Ohne Parameter liefert der Endpunkt
jeden Change; ein Parameter außerhalb dieser drei Namen endet mit `400`, bevor
das erste Event läuft; ein Server ohne den Parameter `target` lehnt ihn mit `400`
ab (*abgeleitet* aus der strengen Parameter-Menge, an einem
Server ohne diese Funktion nicht gefahren). Das Event trägt das Ziel nicht. *Ursprung:*
übernommen aus dem E2E-Lauf von `make test-integration`: ein SSE-Client mit Ziel empfing von einer festen Menge gemischter
Changes genau die Changes seines Ziels und im Ruhefenster von 15 s keine Change
eines anderen Ziels oder ohne Ziel. Die Beispiel-Clients (Go/C#/Kotlin) und
die drei SDK-Packages nehmen alle drei Query-Parameter als eigene
Aufrufparameter entgegen (`-schema`/`-table`/`-target` bzw.
`--schema`/`--table`/`--target`, `StreamChangesAsync(target: …, schema: …,
table: …)`, `streamChanges(target, schema, table)`,
`stream_changes(target, schema, table)`); ein nicht gesetzter oder leerer Wert
des Beispiels erscheint nicht auf dem Draht, am Package erscheint ein gesetzter
Parameter — auch ein leerer — auf dem Draht, und der Server liest einen leeren
Wert als „kein Filter“.

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
  `pg-change-feed-examples:go-sse` aus `examples/Dockerfile`); das optionale
  Flag `ARGS="-target=<ziel>"` wählt das Zustellziel, `-schema=<schema>` und
  `-table=<tabelle>` wählen Schema und Tabelle
- **C#:** `examples/csharp/sse-client` — Container-Aufruf über
  `make example-run-csharp SURFACE=sse` (startet das mit
  `make examples-csharp` gebaute Image); das optionale Flag
  `ARGS="--target=<ziel>"` wählt das Zustellziel, `--schema=<schema>` und
  `--table=<tabelle>` wählen Schema und Tabelle
- **Kotlin:** `examples/kotlin/sse-client` — Container-Aufruf über
  `make example-run-kotlin SURFACE=sse` (startet das mit
  `make examples-kotlin` gebaute Image); das optionale Flag
  `ARGS="--target=<ziel>"` wählt das Zustellziel, `--schema=<schema>` und
  `--table=<tabelle>` wählen Schema und Tabelle

**SDK:** .NET-Anwendungen können statt des Beispiels dasselbe offizielle
NuGet-Package `PgChangeFeed.Client` einbinden
(`dotnet add package PgChangeFeed.Client`) — `PgChangeFeedSseClient.StreamChangesAsync`
öffnet `GET /changes/stream` und liefert ein `IAsyncEnumerable<Change>` mit
allen zehn Feldern der Tabelle oben; das Bearer-Token landet im
`Authorization`-Header, ein fehlender oder unbekannter Token endet den
Aufruf mit `PgChangeFeedUnauthorizedException` (HTTP-Status `401`), statt
den Draht-Vertrag selbst zu implementieren; die optionalen Parameter `target`,
`schema` und `table` (per Name übergeben) wählen Zustellziel, Schema und
Tabelle; siehe `sdks/csharp/README.md`.

Kotlin/JVM-Anwendungen können statt des Beispiels dasselbe offizielle
Gradle-/Maven-Package `pgchangefeed-kotlin` einbinden
(Koordinate `io.github.pt9912:pgchangefeed-kotlin`) —
`PgChangeFeedSseClient.streamChanges()` öffnet `GET /changes/stream` und
liefert eine `Sequence<Change>` mit allen zehn Feldern der Tabelle oben; das
Bearer-Token landet im `Authorization`-Header, ein fehlender oder
unbekannter Token endet den Aufruf mit `PgChangeFeedUnauthorizedException`
(HTTP-Status `401`), statt den Draht-Vertrag selbst zu implementieren; die
optionalen Parameter `target`, `schema` und `table` wählen Zustellziel,
Schema und Tabelle. Wie
beim HTTP-API-Zugriff oben ist der Bezug über Cloudsmith ohne Konto und
ohne Token möglich; der Bezug über **GitHub Packages** verlangt
dagegen immer eine Authentifizierung. Siehe
`sdks/kotlin/pgchangefeed-kotlin/README.md`.

Python-Anwendungen können statt des Beispiels das offizielle
PyPI-Package `pgchangefeed` einbinden
(`pip install pgchangefeed`) — `PgChangeFeedSseClient.stream_changes()`
öffnet denselben Endpunkt `GET /changes/stream` und liefert einen Iterator
über die getypten `StreamChange`-Events mit allen zehn Feldern der Tabelle
oben; das Bearer-Token landet im `Authorization`-Header, ein fehlender oder
unbekannter Token endet den Aufruf mit `PgChangeFeedUnauthorizedError`
(HTTP-Status `401`), statt den Draht-Vertrag selbst zu implementieren; die
optionalen Parameter `target`, `schema` und `table` wählen Zustellziel,
Schema und Tabelle. Siehe `sdks/python/README.md`.

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
mehrere Quellen verfolgt, `cdc.changes.>`. Das Wecksignal trägt kein
Zustellziel; die Wurzel `cdc.changes` bekommt von der Zusatz-Veröffentlichung des
[NATS-Vollinhalts-Streams](#zugriff-über-den-nats-vollinhalts-stream) nichts
(*Ursprung:* Unit-Test des Vollinhalts-Adapters, nur negativ belegt: ohne
Wecksignal-Sender null Nachrichten).

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

Die Beispiele sind zum Lesen und Nachbauen gedacht.

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
Daten. Das Zustellziel gehört nicht zur Nachricht.

**Das Zusatz-Subjekt (Zustellziel):** Jede Change mit einem Zustellziel
(`route_target`, siehe [Routing-Regel konfigurieren](#routing-regel-konfigurieren))
wird zusätzlich auf `cdc.route.<source_id>.<ziel>` veröffentlicht, mit demselben
Payload wie auf dem `cdc.stream…`-Subjekt; eine Change ohne Ziel erzeugt keine
zweite Nachricht, das Subjekt `cdc.stream.<source_id>.<schema>.<table>` bleibt
unverändert. Ein Consumer abonniert `cdc.route.<source_id>.<ziel>` für ein Ziel,
`cdc.route.<source_id>.*` oder `cdc.route.<source_id>.>` für alle Ziele einer
Quelle. Das Alphabet des Zielnamens enthält weder Punkt noch NATS-Platzhalter:
ein Zielname mit `-` und `_` (`eu-west_1`) ist an einem realen NATS-Server ein
einzelnes Subjekt-Token, und `.*` wie `.>` empfangen jedes Ziel (*Ursprung:*
übernommen aus dem Lauf von `make test-notify` an einem realen NATS-Server).
Die zweite
Veröffentlichung ist von der ersten in der Prüfung unabhängig (Tabellen-Subjekt
zuerst, danach das Ziel): ein Fehlschlag oder Überspringen der einen verändert
die andere nicht. Backfill-Changes gehen auf keines der beiden Subjekte (siehe
[Bestand als Backfill überführen](#bestand-als-backfill-überführen)).

**Kosten der zweiten Veröffentlichung:** Die Zeit für `publish` samt Flush zum
Server (je 10.000 Changes, Median von fünf Durchgängen, ohne Abonnent,
Testcontainer-NATS, `make test-notify`) lag mit Ziel-Veröffentlichung beim
0,92- bis 1,30-fachen der Zeit ohne (Spanne mehrerer Läufe am Testcontainer,
ohne Schwelle; das Verhältnis ist aus den gedruckten Zeiten *abgeleitet*).
*Ursprung:* übernommen aus mehreren Läufen am Testcontainer, darunter ein Lauf
mit ohne Ziel 20,86 ms (479.477 Changes/s) und mit Ziel 19,27 ms (518.898
Changes/s), Verhältnis 0,92; weitere Verhältnisse waren 1,05, 1,16 und 1,30. Ein Aufschlag der zweiten Veröffentlichung ist an diesem
Messaufbau nicht auflösbar. Die Messung hat keine Schwelle und deckt die
Verteilung an Abonnenten nicht ab; eine Last-Zusage folgt daraus nicht.

**Zustellsemantik:** Core NATS, Fire-and-Forget, kein Replay — dieselbe
Zusicherung wie gRPC und SSE, auch für das Zusatz-Subjekt: ein nicht
verbundener oder gerade getrennter Consumer verpasst die Änderung ersatzlos und
holt sie über [Änderungen lesen](#änderungen-lesen) nach — die Lücke schließt
`WHERE route_target = '<ziel>'` auf `cdc.changes` oder `GET /changes?target=<ziel>`.

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
  `pg-change-feed-examples:go-nats-stream` aus `examples/Dockerfile`); mit
  `ARGS="-source=<quelle> -target=<ziel>"` abonniert es das Zusatz-Subjekt
  `cdc.route.<source_id>.<ziel>` statt `cdc.stream.>`
- **C#:** `examples/csharp/nats-stream-client` — Container-Aufruf über
  `make example-run-csharp SURFACE=nats-stream` (startet das mit
  `make examples-csharp` gebaute Image); mit
  `ARGS="--source=<quelle> --target=<ziel>"` das Zusatz-Subjekt
- **Kotlin:** `examples/kotlin/nats-stream-client` — Container-Aufruf über
  `make example-run-kotlin SURFACE=nats-stream` (startet das mit
  `make examples-kotlin` gebaute Image); mit
  `ARGS="--source=<quelle> --target=<ziel>"` das Zusatz-Subjekt

`-source` und `-target` gelten nur zusammen; ein Wert mit Punkt, `*`, `>` oder
Leerraum wird mit Exit 2 abgelehnt, bevor ein Abonnement entsteht.

Die Beispiele sind zum Lesen und Nachbauen gedacht.

**SDK:** .NET-Anwendungen können statt des Beispiels dasselbe offizielle
NuGet-Package `PgChangeFeed.Client` einbinden
(`dotnet add package PgChangeFeed.Client`) — `PgChangeFeedNatsStreamClient.StreamChangesAsync`
abonniert den Vollinhalts-Namensraum (Default `cdc.stream.>`, oder ein über
`BuildSubject`/`BuildSourceSubject` eingeschränktes Subjekt; `BuildTargetSubject`
und `BuildSourceTargetsSubject` bauen das Zusatz-Subjekt eines Zustellziels bzw.
aller Ziele einer Quelle) und liefert ein
`IAsyncEnumerable<Change>` mit allen zehn Feldern der Tabelle oben; die
Authentifizierung ist verbindungsseitig (derselbe `CDC_NATS_STREAM_TOKEN`
wie oben), ein abgelehnter Verbindungsversuch endet die Aufzählung mit einer
NATS-eigenen Ausnahme, statt den Draht-Vertrag selbst zu implementieren;
siehe `sdks/csharp/README.md`.

Kotlin/JVM-Anwendungen können statt des Beispiels dasselbe offizielle
Gradle-/Maven-Package `pgchangefeed-kotlin` einbinden
(Koordinate `io.github.pt9912:pgchangefeed-kotlin`) —
`PgChangeFeedNatsStreamClient.streamChanges()` abonniert den
Vollinhalts-Namensraum (Default `cdc.stream.>`, oder ein über `buildSubject`/
`buildSourceSubject` eingeschränktes Subjekt; `buildTargetSubject` und
`buildSourceTargetsSubject` bauen das Zusatz-Subjekt eines Zustellziels bzw.
aller Ziele einer Quelle) und liefert eine
`Sequence<Change>` mit allen zehn Feldern der Tabelle oben; die
Authentifizierung ist verbindungsseitig (derselbe `CDC_NATS_STREAM_TOKEN`
wie oben), ein abgelehnter Verbindungsversuch endet die Sequenz mit der
zugrunde liegenden `io.nats.client`-Ausnahme unverändert, statt den
Draht-Vertrag selbst zu implementieren oder eine zweite Fehlerklassen-Hierarchie
zu erfinden. Wie beim HTTP-API-Zugriff oben ist der Bezug über Cloudsmith
ohne Konto und ohne Token möglich; der Bezug über **GitHub
Packages** verlangt dagegen immer eine Authentifizierung. Siehe `sdks/kotlin/pgchangefeed-kotlin/README.md`.

Python-Anwendungen können statt des Beispiels das offizielle
PyPI-Package `pgchangefeed` einbinden
(`pip install pgchangefeed`) — `PgChangeFeedNatsStreamClient.stream_changes()`
abonniert den Vollinhalts-Namensraum `cdc.stream.<source_id>.>` (alle
Tabellen einer Quelle; `source_id` ist Konstruktor-Argument; der optionale
Parameter `target` abonniert stattdessen das Zusatz-Subjekt
`cdc.route.<source_id>.<ziel>`) und liefert
einen Iterator über die getypten `StreamChange`-Events mit allen zehn
Feldern der Tabelle oben; die Authentifizierung ist verbindungsseitig
(derselbe `CDC_NATS_STREAM_TOKEN` wie oben), ein abgelehnter
Verbindungsversuch endet am Connect-Fehler der NATS-Bibliothek unverändert,
statt den Draht-Vertrag selbst zu implementieren. Das Package deckt damit
dieselben vier Zugriffswege wie die C#-/Kotlin-Pendants ab. Siehe
`sdks/python/README.md`.

## 5. Konfiguration

### Umgebungsvariablen des Feed-Containers

| Variable | Pflicht | Bedeutung |
|---|---|---|
| `CDC_CAPTURE_DSN` | ja | Verbindung über die Rolle `cdc_capture` (Store-Adapter, Replication-Stream, Ausführung eines Backfills) |
| `CDC_ADMIN_DSN` | ja | Verbindung über die Rolle `cdc_admin` (Tabellen-Aktivierung, Heartbeat, Verarbeitung der Antrags-Queue, Annahme eines Backfill-Antrags, `register-consumer`/`acknowledge-consumer`) |
| `CDC_READER_DSN` | ja | Verbindung über die Rolle `cdc_reader` (`--healthcheck`, `diagnose`, Lesen der Sicht für den OTLP-Export) |
| `CDC_SOURCE_ID` | ja | Kennung der Quelle (muss in `cdc.source` registriert sein) |
| `CDC_PUBLICATION` | ja | Name der PostgreSQL-Publication |
| `CDC_SLOT` | ja | Name des Logical-Replication-Slots |
| `CDC_TABLES` | ja, falls keine Konfigurationsdatei dieselbe Aktivierung trägt | Aktivierte Tabellen, Format `schema.tabelle=tabelle-id:schema-version-id`, kommagetrennt |
| `CDC_LOG_LEVEL` | nein | Log-Level des strukturierten JSON-Loggers (Default `info`) |
| `CDC_NATS_URL` | nein | NATS-Server-URL für das Change-Notification-Wecksignal (`cdc.changes.<source_id>.<schema>.<table>`, tabellen-granular, leerer Payload); ungesetzt bleibt das Feature vollständig deaktiviert, gesetzt ist eine erfolgreiche Verbindung Vorbedingung des Starts (Fehlerklasse `configuration`). Zusammen mit `CDC_NATS_STREAM_TOKEN` aktiviert dieselbe Variable zusätzlich den dritten, vollinhaltstragenden NATS-Zustellweg (siehe [Zugriff über den NATS-Vollinhalts-Stream](#zugriff-über-den-nats-vollinhalts-stream)) |
| `CDC_NATS_STREAM_TOKEN` | nein | Verbindungs-Token des dritten, vollinhaltstragenden NATS-Zustellwegs; wirkt nur zusammen mit gesetztem `CDC_NATS_URL` — ist nur `CDC_NATS_STREAM_TOKEN` gesetzt, aber `CDC_NATS_URL` leer, startet der Container nicht (Fehlerklasse `configuration`). Ein gesetzter Wert verlangt vom NATS-Server denselben Token **serverweit**, auch für die Wecksignal-Verbindung (siehe dortiger Abschnitt) |
| `CDC_HTTP_ADDR` | nein | Horch-Adresse der HTTP-/JSON-API (`host:port`); ungesetzt bleibt die API vollständig deaktiviert — kein Server, keine zusätzliche Verbindung. Anders als `CDC_NATS_URL` ist die Adresse **keine** Start-Vorbedingung: Ist sie gesetzt, öffnet der Prozess den Server in eigener Goroutine und läuft unverändert weiter; scheitert das Binden der Adresse (z. B. belegter Port), meldet er das im Log und der Erfassungsbetrieb bleibt davon unberührt |
| `CDC_API_TOKEN_READER` | nein | Genau ein Bearer-Token der lesenden Rechtsklasse der HTTP- und gRPC-API (ein Komma gehört zum Token); ungesetzt (leer) trägt das Singular kein Token — ein Aufruf mit einem Token, das keiner konfigurierten Klasse entspricht, endet `401` |
| `CDC_API_TOKEN_ADMIN` | nein | Genau ein Bearer-Token der administrativen Rechtsklasse der HTTP- und gRPC-API (deckt die lesende Klasse implizit mit ab); leer bedeutet dasselbe wie bei `CDC_API_TOKEN_READER` |
| `CDC_API_TOKENS_READER` | nein | Kommagetrennte Liste von Bearer-Token der lesenden Klasse; gültig ist die Vereinigung mit `CDC_API_TOKEN_READER`. Ein leeres Element oder ein Element mit Leerraum verhindert den Start (Fehlerklasse `configuration`, `PCF-E2008`); eine leere Variable gilt als nicht gesetzt. Nur in der Umgebung, nicht in der Konfigurationsdatei (siehe [API-Token in zwei Neustarts wechseln](#api-token-in-zwei-neustarts-wechseln)) |
| `CDC_API_TOKENS_ADMIN` | nein | Kommagetrennte Liste von Bearer-Token der administrativen Klasse; gültig ist die Vereinigung mit `CDC_API_TOKEN_ADMIN`, dieselben Regeln wie bei `CDC_API_TOKENS_READER`. Steht ein Wert in beiden Klassen, gilt er als administrativ |
| `CDC_GRPC_ADDR` | nein | Horch-Adresse des gRPC-Servers (`host:port`); aktiviert gemeinsam den Change-Stream (`ChangeStream`) und die Verwaltungs-API (`Administration`) auf demselben Port. Ungesetzt bleibt der Server vollständig deaktiviert — kein Listener. Wie `CDC_HTTP_ADDR` keine Start-Vorbedingung |
| `CDC_TLS_CERT_FILE` | nein | Pfad zur Zertifikatsdatei (PEM, das Zertifikat und dahinter seine Kette) der HTTP- und der gRPC-Schnittstelle. Zertifikat und Schlüssel gehören zusammen: ist nur eine der beiden Variablen gesetzt, startet der Container nicht (Fehlerklasse `configuration`, `PCF-E2009`). Beide gesetzt verschlüsseln HTTP (einschließlich `GET /changes/stream`) und gRPC (siehe [Schnittstellen mit TLS verschlüsseln](#schnittstellen-mit-tls-verschlüsseln)); beide ungesetzt lassen beide Schnittstellen unverschlüsselt. Die Variable trägt einen Pfad, keine Zugangsdaten, und hat das Datei-Feld `tls_cert_file` |
| `CDC_TLS_KEY_FILE` | nein | Pfad zur Datei mit dem privaten Schlüssel (PEM) zum Zertifikat; dieselben Regeln wie bei `CDC_TLS_CERT_FILE`. Der Schlüssel selbst steht weder in einer Umgebungsvariable noch in der Konfigurationsdatei; das Datei-Feld `tls_key_file` trägt nur den Pfad |
| `CDC_OTLP_ENDPOINT` | nein | Basis-URL des OpenTelemetry-Empfängers (Schema `http` oder `https`, ein Host); gesetzt schaltet den periodischen Export der Kennzahlen ein, ungesetzt bleibt er aus — keine Verbindung zu einem Empfänger, Betrieb und SQL-Sicht unverändert. Ein Wert mit anderem Schema oder ohne Host verhindert den Start (Fehlerklasse `configuration`, `PCF-E2011`). Die URL kann Zugangsdaten tragen und steht nur in der Umgebung, nicht in der Konfigurationsdatei (siehe [Metriken per OTLP übertragen](#metriken-per-otlp-übertragen)) |
| `CDC_OTLP_HEADERS` | nein | Zusätzliche Header jeder Export-Anfrage, Form `Schlüssel=Wert`, mehrere durch Komma getrennt (geteilt am ersten `=`). Ein Element ohne `=`, mit leerem Schlüssel oder Leerraum im Schlüssel verhindert den Start (Fehlerklasse `configuration`, `PCF-E2012`); ohne `CDC_OTLP_ENDPOINT` gesetzt ebenso (`PCF-E2014`). Trägt Zugangsdaten und steht nur in der Umgebung |
| `CDC_OTLP_INTERVAL_SECONDS` | nein | Takt der Übertragung in Sekunden, ganze Zahl von 5 bis 3600, Default 60; ein Wert außerhalb oder keine ganze Zahl verhindert den Start (Fehlerklasse `configuration`, `PCF-E2013`). Hat das Datei-Feld `otlp_interval`; die Umgebungsvariable schlägt es |
| `CDC_CONFIG_FILE` | nein | Pfad zu einer optionalen YAML-Konfigurationsdatei (siehe unten) |

Fehlt eine Pflichtvariable und liefert auch keine Konfigurationsdatei
einen Wert für dasselbe Feld, startet der Container nicht (Fehlerklasse
`configuration`).

### Optionale YAML-Konfigurationsdatei (`CDC_CONFIG_FILE`)

Additiv zu den Umgebungsvariablen: Ist
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

Die zwei Pfade des TLS-Paars sind ebenfalls Datei-Felder: `tls_cert_file` und
`tls_key_file` entsprechen `CDC_TLS_CERT_FILE` bzw. `CDC_TLS_KEY_FILE`, mit
demselben Feld-für-Feld-Vorrang der Umgebungsvariable. Es sind Pfade, keine
Zugangsdaten; der private Schlüssel selbst steht in keiner Konfigurationsdatei.
Beide Felder gehören zusammen (siehe [Schnittstellen mit TLS
verschlüsseln](#schnittstellen-mit-tls-verschlüsseln)).

Der Takt des Metrik-Exports ist ein Datei-Feld: `otlp_interval` trägt die
Sekunden (5 bis 3600) und entspricht `CDC_OTLP_INTERVAL_SECONDS`, mit demselben
Feld-für-Feld-Vorrang der Umgebungsvariable und dem Default 60. Eine Zahl trägt
keine Zugangsdaten. Endpunkt und Header des Exports (`CDC_OTLP_ENDPOINT`,
`CDC_OTLP_HEADERS`) haben kein Datei-Feld (siehe [Metriken per OTLP
übertragen](#metriken-per-otlp-übertragen)).

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
tls_cert_file: /etc/cdc/tls/server.pem
tls_key_file: /etc/cdc/tls/server-key.pem
otlp_interval: 60
```

**Wichtig — Zugangsdaten bleiben env-var-exklusiv:** Die Schlüssel
`capture_dsn`, `admin_dsn`, `reader_dsn`, `api_token_reader`,
`api_token_admin`, `api_tokens_reader`, `api_tokens_admin`, `nats_url`,
`nats_stream_token`, `otlp_endpoint` und `otlp_headers` — elf Schlüssel — dürfen
in dieser Datei **nicht** vorkommen. Ein Treffer bricht das Laden mit einer eigenen,
den Grund nennenden Zeile ab (Fehlerklasse `configuration`) — eine
Konfigurationsdatei landet typischerweise in Kanälen (Repository,
ConfigMap, Backup), die für Zugangsdaten nicht vorgesehen sind. Die Grenze
ist die **Form** des Feldes, nicht sein Wert: `http_addr`/`grpc_addr` sind
`host:port` und können keine Zugangsdaten tragen, `nats_url` ist eine URL
und kann Benutzer sowie Passwort einbetten (`nats://benutzer:passwort@host:4222`),
`nats_stream_token` und die vier API-Token-Schlüssel tragen denselben
Zugangsdaten-Charakter, `otlp_endpoint` ist wie `nats_url` eine URL, die
Benutzer und Passwort einbetten kann, und `otlp_headers` trägt Zugangsdaten als
Header-Werte (`otlp_interval` ist eine Zahl und kein Mitglied dieser Klasse). Ein unbekannter Schlüssel bricht das
Laden ebenfalls ab (striktes Decoding).

Die env-exklusiven Variablen `CDC_NATS_URL`, `CDC_NATS_STREAM_TOKEN`,
`CDC_API_TOKEN_READER`, `CDC_API_TOKEN_ADMIN`, `CDC_API_TOKENS_READER`,
`CDC_API_TOKENS_ADMIN`, `CDC_OTLP_ENDPOINT` und `CDC_OTLP_HEADERS` werden **auch unter
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
Klassen:

| Klasse | Bedeutung | Verhalten |
|---|---|---|
| `transient` | vorübergehend nicht verfügbare Quelle/Speicher | Erneuter Versuch mit begrenztem Backoff (siehe [Neustart nach einem Fehler](#neustart-nach-einem-fehler)); der Erfassungspfad trägt die Klasse, wenn das Wiederholungsfenster erschöpft ist — ein Backfill-Run trägt sie als Klasse seines Fehlertexts |
| `configuration` | ungültige oder fehlende Umgebungsvariable | Kein Start, sichtbarer Fehler |
| `permission` | fehlende Berechtigung | Sichtbarer Fehler, kein stiller Retry; der Erfassungspfad trägt sie bei SQLSTATE 42501 und bei Fehlern der Klasse 28 (Authentifizierung/Autorisierung) am Quellzugriff, ein Backfill-Run als Klasse seines Fehlertexts (z. B. fehlendes `SELECT` auf die Quelltabelle) |
| `schema` | eine Replikationsnachricht ist nicht sicher interpretierbar (z. B. TRUNCATE, unbekannter Nachrichtentyp) oder eine Transformations- oder Routing-Regel ist auf die Relation einer Change nicht anwendbar (siehe [Transformationsregel konfigurieren](#transformationsregel-konfigurieren) und [Routing-Regel konfigurieren](#routing-regel-konfigurieren)) | Sichtbarer Fehler, kein stilles Überspringen; ein Backfill-Run trägt dieselbe Ursache als Klasse seines Fehlertexts, run-lokal — bei einer Routing-Regel, deren Bedingungsspalte der Snapshot nicht trägt, vor der ersten Zeile |
| `storage` | Persistenzfehler | Kein Source-ACK, damit keine Änderung verloren geht |
| `replication` | zwei Unterarten: **Stream-Ordnungs-Verletzung** (z. B. Commit ohne offene Transaktion) oder **Transport-/Verbindungsstörung** (Verbindungsaufbau, Slot, Keepalive, Quell-Bestätigung) | Stream-Ordnungs-Verletzung: sofortiger, sichtbarer Abbruch, unabhängig vom WAL-Rückstand. Transport-/Verbindungsstörung: Schwellen-Überwachung über den WAL-Rückstand (siehe [WAL-Rückstand prüfen](#wal-rückstand-prüfen)) — kontrollierte Fortsetzung unterhalb 1 GiB, sichtbarer Abbruch darüber |
| `internal` | unerwarteter interner Fehler, der keiner anderen Klasse zuzuordnen ist | Sichtbarer Fehler; realer Fallback für jeden nicht erkannten Fehler |

`transient` und `permission` sind im Erfassungspfad beobachtbar (Heartbeat-Fehlerzustand
nach erschöpftem Wiederholungsfenster bzw. bei einer Berechtigungs-Abweisung des Servers)
und als Klasse im Fehlertext eines fehlgeschlagenen Backfill-Runs
(`cdc.backfill_status.error_message`).
`internal` ist der real erreichbare
Fallback-Zweig: Jeder Fehler, der keiner der übrigen sechs Klassen
zugeordnet werden kann, fällt auf `internal` zurück.

Ein Fehler des Erfassungspfads jeder Klasse beendet den Container-Prozess mit Ausgang 1 (ein Fehler eines Backfill-Runs ist run-lokal und beendet ihn nicht); der
zuletzt beobachtete Fehlerzustand wird zusätzlich in
`cdc.heartbeat.error_class` und sein Meldungscode in
`cdc.heartbeat.error_code` festgehalten (siehe [Meldungscodes](#meldungscodes))
und bei einem erfolgreichen Neustart automatisch wieder gelöscht.

### Meldungscodes

Jeder klassifizierte Fehler, jede abgelehnte Eingabe und jede Warnung mit einer
Maßnahme für den Betreiber trägt einen Meldungscode, der die Ursache feiner
benennt als die Fehlerklasse. Ein Code hat die Form `PCF-E4003`: der Buchstabe
`E` kennzeichnet einen Fehler, die erste Ziffer die Fehlerklasse (1 `transient`,
2 `configuration`, 3 `permission`, 4 `schema`, 5 `storage`, 6 `replication`,
7 `internal`), die übrigen drei Ziffern die Ursache. Die Endung `000` ist der
Rückfall einer Klasse: ein Fehler, dem keine einzelne Ursache zugeordnet ist,
trägt den Rückfall seiner Klasse und nie keinen Code. Die erste Ziffer `8`
kennzeichnet die Ablehnung einer Eingabe, etwa eines Antrags oder eines
API-Aufrufs; sie trägt keine Fehlerklasse. Eine Warnung trägt den Buchstaben `W` (`PCF-W3002`); ihre erste
Ziffer ist der Bereich (1 Erfassung und Replikation, 2 Backfill, 3 Retention und
Speicher, 4 Verwaltung einschließlich der Aufrufe der HTTP- und gRPC-API, 6 Beobachtbarkeit und Transport; ein
Bereich 5 für Konfiguration und Start ist vorgesehen und trägt keine Warnung), eine
Warnung trägt keine Fehlerklasse. Zeilen des
Berichts von `diagnose` mit Ausnahme der Zeile „Fehlerzustand“ tragen keinen Code.

Der Code steht an diesen Stellen:

| Stelle | Form |
|---|---|
| Log-Zeile einer Warnung | eigenes Attribut `code=PCF-W3002`; der Meldungstext trägt den Code nicht. Die Warn-Zeile `heartbeat: Fehlerzustand gemeldet` trägt unter `code` den Code des Fehlerzustands (ein Code mit `E`), ein Filter auf Warn-Codes erfasst sie nicht |
| `cdc.heartbeat.error_code` | Meldungscode des Fehlerzustands neben `error_class`, `NULL` im Normalbetrieb |
| Zeile „Fehlerzustand“ von `diagnose`, `GET /diagnose` (`error_code`) und RPC `Diagnose` (`error_code`) | `schema [PCF-E4003]` bzw. das Feld `error_code` |
| Log-Zeile und Zeile beim Prozessende | `Fehlerklasse schema [PCF-E4003]: <Ursache>` |
| `cdc.backfill_status.error_message` eines fehlgeschlagenen Runs | `schema [PCF-E4004]: <Ursache>` |
| `cdc.administration_request.error_message` eines abgelehnten Antrags | `abgelehnt [PCF-E8021]: <Klartext>` |
| `cdc.administration_request.error_message` eines Antrags, der an einem klassifizierten Fehler scheiterte | `Fehlerklasse internal [PCF-E7002]: <Ursache>` |
| Meldungen von `make schema-rollout` | `FEHLER [PCF-E2007]: <Text>` |
| Fehlerantwort der HTTP-API (`400`, `404`, `500`, `503`) | Feld `code` neben `error`: `{"error": "<Klartext>", "code": "PCF-E8051"}`; bei `401` und `403` fehlt das Feld |
| Fehlerstatus der gRPC-API (`InvalidArgument`, `NotFound`, `Internal`) | Statusdetail `google.rpc.ErrorInfo`: `reason` ist der Code (`PCF-E8051`), `domain` ist `pg-change-feed`; `Unauthenticated` und `PermissionDenied` tragen kein Detail |

**Der Text ist nicht Vertrag, der Code ist es.** Stabil sind der Code, die Klasse,
das Wort `Fehlerklasse` am Anfang einer Fehlerzeile und der Ausgang des Prozesses; der Text
nach dem Code nennt Laufzeitdetails (Regelname, Spalte, Adresse) und kann
sich ändern. Wer Meldungen maschinell auswertet, wertet den Code oder die Klasse
aus, nicht den Text. Ein Code wird nie neu belegt und seine Klasse ändert sich nie.
Entfällt eine Ursache, bleibt ihr Code in dieser Tabelle und trägt den Vermerk
„zurückgezogen“.

| Code | Klasse | Bedeutung | Maßnahme |
|---|---|---|---|
| `PCF-E1000` | `transient` | vorübergehend nicht verfügbare Quelle oder nicht näher zugeordnete vorübergehende Störung (Rückfall) | Erreichbarkeit der Quelle und des Netzes prüfen, Container neu starten |
| `PCF-E1002` | `transient` | das Wiederholungsfenster (5 Minuten) bei einer vorübergehenden Störung am Quellzugriff ist erschöpft; der Prozess endet mit Ausgang 1 | Erreichbarkeit der Quelle prüfen, Container neu starten (siehe [Neustart nach einem Fehler](#neustart-nach-einem-fehler)) |
| `PCF-E1003` | `transient` | die Quelle war für den Snapshot eines Backfills vorübergehend nicht verfügbar (Zeitlimit der Slot-Anlage, Verbindungsabbruch, Umschreiben der Tabelle zwischen Export und Import) | Backfill neu beantragen |
| `PCF-E2000` | `configuration` | ungültige oder falsch gesetzte Konfiguration ohne nähere Zuordnung (Rückfall) | Konfiguration prüfen |
| `PCF-E2001` | `configuration` | eine Pflicht-Umgebungsvariable oder ein Wert der Konfigurationsdatei fehlt oder ist falsch gesetzt; der Container startet nicht — oder der laufende Prozess ist für den Stream unvollständig eingerichtet und `GET /changes/stream` bzw. der gRPC-Stream antwortet mit `503` bzw. `Internal` | beim Start nennt die Meldung die Variable; korrigieren und neu starten (siehe [Container startet nicht](#container-startet-nicht)); bei der Antwort des Streams die Einrichtung des Prozesses prüfen und neu starten |
| `PCF-E2002` | `configuration` | der Replikationszugriff ist falsch konfiguriert: Slot- oder Publication-Name außerhalb des Bezeichner-Alphabets, Publication fehlt an der Quelle, DSN fehlt | Namen und Publication prüfen |
| `PCF-E2003` | `configuration` | ein Bezeichner einer Aktivierung (Schema, Tabelle, Publication) liegt außerhalb des Bezeichner-Alphabets | Namen korrigieren |
| `PCF-E2004` | `configuration` | der Snapshot eines Backfills steht im falschen Stand der Konfiguration: Tabelle nicht vorhanden, ungültige Kennung oder keine freie Reserve bei `max_replication_slots`/`max_wal_senders` | Tabelle und Einstellungen der Instanz prüfen, Backfill neu beantragen |
| `PCF-E2005` | `configuration` | der Backfill-Run fand bei der Ausführung keine Aktivierung der Tabelle (Bindung oder Mitgliedschaft in der Publication fehlt) | Tabelle aktivieren, Backfill neu beantragen |
| `PCF-E2006` | `configuration` | während des Backfill-Runs hat sich der Ausschluss- oder Regelstand der Tabelle geändert; der Run endet ohne Change | Änderung des Regelstands abschließen, Backfill neu beantragen |
| `PCF-E2007` | `configuration` | die Schema-Quelldatei von `make schema-rollout` fehlt | Variable `SCHEMA_SOURCE` und Arbeitsbaum prüfen |
| `PCF-E2008` | `configuration` | eine Token-Liste (`CDC_API_TOKENS_READER`, `CDC_API_TOKENS_ADMIN`) trägt ein leeres Element oder ein Element mit Leerraum; der Container startet nicht, die Zeile nennt Variable und Nummer des Elements | Liste prüfen (siehe [API-Token in zwei Neustarts wechseln](#api-token-in-zwei-neustarts-wechseln)) |
| `PCF-E2009` | `configuration` | nur eine der beiden TLS-Angaben (`CDC_TLS_CERT_FILE`, `CDC_TLS_KEY_FILE` bzw. `tls_cert_file`, `tls_key_file`) ist gesetzt; der Container startet nicht (Prozessausgang 2), die Zeile nennt die gesetzte und die fehlende Angabe | beide Angaben setzen oder beide entfernen (siehe [Schnittstellen mit TLS verschlüsseln](#schnittstellen-mit-tls-verschlüsseln)) |
| `PCF-E2010` | `configuration` | das TLS-Paar ist nicht ladbar: eine Datei ist nicht lesbar, trägt keine gültigen PEM-Daten, oder Zertifikat und Schlüssel gehören nicht zusammen; der Container startet nicht (Prozessausgang 1), die Zeile nennt die Pfade und die Ursache | Pfade, Dateirechte für den Benutzer des Containers und die Zusammengehörigkeit von Zertifikat und Schlüssel prüfen (siehe [Schnittstellen mit TLS verschlüsseln](#schnittstellen-mit-tls-verschlüsseln)) |
| `PCF-E2011` | `configuration` | `CDC_OTLP_ENDPOINT` ist keine URL mit Schema `http` oder `https` und einem Host; der Container startet nicht (Prozessausgang 2), die Zeile nennt die Variable und den Grund, nie die URL | Basis-URL des Empfängers prüfen (siehe [Metriken per OTLP übertragen](#metriken-per-otlp-übertragen)) |
| `PCF-E2012` | `configuration` | `CDC_OTLP_HEADERS` trägt ein Element ohne `=`, mit leerem Schlüssel, mit Leerraum oder unzulässigem Zeichen im Schlüssel oder mit einem Steuerzeichen im Wert; der Container startet nicht, die Zeile nennt Variable und Nummer des Elements, nie dessen Inhalt | Form `Schlüssel=Wert,Schlüssel2=Wert2` prüfen |
| `PCF-E2013` | `configuration` | der Takt des Exports (`CDC_OTLP_INTERVAL_SECONDS` oder `otlp_interval`) ist keine ganze Zahl von 5 bis 3600; der Container startet nicht, die Zeile nennt Quelle und Wert | Takt auf eine ganze Zahl von 5 bis 3600 setzen |
| `PCF-E2014` | `configuration` | `CDC_OTLP_HEADERS` ist gesetzt, aber `CDC_OTLP_ENDPOINT` fehlt; der Container startet nicht | Endpunkt setzen oder die Header entfernen |
| `PCF-E3000` | `permission` | fehlende Berechtigung ohne nähere Zuordnung (Rückfall) | Rechte der Rollen prüfen (siehe [Zugriff und Rollen](#zugriff-und-rollen)) |
| `PCF-E3001` | `permission` | der Server weist den Replikationszugriff ab (SQLSTATE 42501 oder Klasse 28) | Rechte und `REPLICATION`-Attribut der Capture-Rolle prüfen |
| `PCF-E3002` | `permission` | die Capture-Rolle darf die Quelltabelle eines Backfills nicht lesen oder den temporären Slot nicht anlegen | `SELECT` auf die Tabelle und das `REPLICATION`-Attribut prüfen, Backfill neu beantragen |
| `PCF-E4000` | `schema` | nicht sicher interpretierbare Schemaänderung ohne nähere Zuordnung (Rückfall) | Meldung im Log lesen |
| `PCF-E4001` | `schema` | eine Replikationsnachricht ist nicht sicher interpretierbar; die Erfassung endet sichtbar, es geht keine Änderung verloren | Ursache im Log prüfen, Container neu starten |
| `PCF-E4002` | `schema` | `TRUNCATE` an einer erfassten Tabelle wird nicht unterstützt; die Erfassung endet sichtbar | `TRUNCATE` an erfassten Tabellen vermeiden, stattdessen `DELETE` verwenden |
| `PCF-E4003` | `schema` | Relation-Änderung nicht sicher als Obermenge interpretierbar (Spalte entfernt, Typ geändert, Spalte umbenannt); die Erfassung endet sichtbar | die Schemaänderung an der Quelle prüfen; das Entfernen einer Regel genügt nicht (siehe [Routing-Regel konfigurieren](#routing-regel-konfigurieren)) |
| `PCF-E4004` | `schema` | eine Transformationsregel ist auf die Änderung nicht anwendbar, im Backfill auf die Spalten des Snapshots | Regelstand ändern (siehe [Transformationsregel konfigurieren](#transformationsregel-konfigurieren)) |
| `PCF-E4005` | `schema` | eine Routing-Regel ist auf die Änderung nicht anwendbar, im Backfill auf die Spalten des Snapshots | Regelstand ändern (siehe [Routing-Regel konfigurieren](#routing-regel-konfigurieren)) |
| `PCF-E5000` | `storage` | Persistenzfehler ohne nähere Zuordnung (Rückfall) | Erreichbarkeit, Speicherplatz und Rechte der CDC-Datenbank prüfen |
| `PCF-E5001` | `storage` | Persistenzfehler im Change-Speicher; es wird keine Quellposition bestätigt | Erreichbarkeit, Speicherplatz und Rechte der CDC-Datenbank prüfen |
| `PCF-E5003` | `storage` | Persistenzfehler im Backfill-Speicher | Erreichbarkeit und Rechte der CDC-Datenbank prüfen, Backfill neu beantragen |
| `PCF-E5004` | `storage` | Persistenzfehler im Speicher der Consumer-Stände | Erreichbarkeit und Rechte der CDC-Datenbank prüfen |
| `PCF-E5006` | `storage` | Persistenzfehler im Heartbeat-Speicher | Erreichbarkeit und Rechte der CDC-Datenbank prüfen |
| `PCF-E5007` | `storage` | Persistenzfehler im Schema-Speicher | Erreichbarkeit und Rechte der CDC-Datenbank prüfen |
| `PCF-E5008` | `storage` | Lesefehler im Snapshot eines Backfills (etwa eine Spalte, die zwischen Export und Import entfernt wurde) | Tabelle prüfen, Backfill neu beantragen |
| `PCF-E6000` | `replication` | Störung des Replication-Streams oder Slots ohne nähere Zuordnung (Rückfall) | Quelle, Slot und WAL-Rückstand prüfen (siehe [WAL-Rückstand prüfen](#wal-rückstand-prüfen)) |
| `PCF-E6001` | `replication` | Störung an Replikationsverbindung oder Slot (Verbindungsaufbau, Start, Keepalive) | Quelle, Slot und WAL-Rückstand prüfen (siehe [WAL-Rückstand prüfen](#wal-rückstand-prüfen)) |
| `PCF-E6002` | `replication` | die Bestätigung der Position an der Quelle ist fehlgeschlagen, oder der WAL-Rückstand liegt über der Fehlerschwelle | WAL-Rückstand und Slot prüfen (siehe [WAL-Rückstand prüfen](#wal-rückstand-prüfen)) |
| `PCF-E6003` | `replication` | Verletzung der Stream-Ordnung (Änderung oder Commit ohne offene Transaktion, BEGIN bei offener Transaktion); sofortiger Abbruch | Log sichern, Container neu starten |
| `PCF-E6004` | `replication` | Störung an Slot oder Replikationsverbindung beim Snapshot eines Backfills | Reserve von `max_replication_slots` prüfen, Backfill neu beantragen |
| `PCF-E7000` | `internal` | unerwarteter interner Fehler, der keiner anderen Klasse zuzuordnen ist; auch der Rückfallcode eines API-Aufrufs mit `500`, dessen Ursache keine Fehlerklasse trägt | in der Log-Zeile der Warnung `PCF-W4008` steht die Ursache im Attribut `error`; Log sichern, Container neu starten, bei Wiederholung das Log dem Support übergeben |
| `PCF-E7001` | `internal` | das Wecksignal über NATS konnte nicht veröffentlicht werden; die Änderung selbst ist erfasst und über SQL lesbar | NATS-Server und `CDC_NATS_URL` prüfen |
| `PCF-E7002` | `internal` | Persistenzfehler an der Antrags-Queue | Erreichbarkeit und Rechte der CDC-Datenbank prüfen, Antrag erneut stellen |
| `PCF-E7003` | `internal` | Lesefehler an den Diagnose-Views | Erreichbarkeit und Rechte der Rolle `cdc_reader` prüfen |
| `PCF-E7004` | `internal` | Persistenzfehler im Backfill-Speicher in einer Meldung außerhalb des Fehlertexts eines Runs (dort steht `PCF-E5003`) | Erreichbarkeit und Rechte der CDC-Datenbank prüfen, Backfill neu beantragen |
| `PCF-E7005` | `internal` | Persistenzfehler im Schema-Speicher in einer Meldung außerhalb des Fehlertexts eines Runs (dort steht `PCF-E5007`) | Erreichbarkeit und Rechte der CDC-Datenbank prüfen |
| `PCF-E7006` | `internal` | fehlende Berechtigung für den Snapshot eines Backfills in einer Meldung außerhalb des Fehlertexts eines Runs (dort steht `PCF-E3002`) | `SELECT` auf die Tabelle und das `REPLICATION`-Attribut prüfen, Backfill neu beantragen |
| `PCF-E7007` | `internal` | der Snapshot eines Backfills steht im falschen Stand der Konfiguration, in einer Meldung außerhalb des Fehlertexts eines Runs (dort steht `PCF-E2004`) | Tabelle und Einstellungen der Instanz prüfen, Backfill neu beantragen |
| `PCF-E7008` | `internal` | die Quelle war für den Snapshot eines Backfills vorübergehend nicht verfügbar, in einer Meldung außerhalb des Fehlertexts eines Runs (dort steht `PCF-E1003`) | Backfill neu beantragen |
| `PCF-E7009` | `internal` | Störung an Slot oder Replikationsverbindung beim Snapshot eines Backfills, in einer Meldung außerhalb des Fehlertexts eines Runs (dort steht `PCF-E6004`) | Reserve von `max_replication_slots` prüfen, Backfill neu beantragen |
| `PCF-E7010` | `internal` | Lesefehler im Snapshot eines Backfills, in einer Meldung außerhalb des Fehlertexts eines Runs (dort steht `PCF-E5008`) | Tabelle prüfen, Backfill neu beantragen |
| `PCF-E8000` | Ablehnung | ein Antrag ist ungültig, ohne dass eine nähere Ursache benannt ist (Rückfall der Ablehnungen) | Antrag prüfen und erneut stellen |
| `PCF-E8001` | Ablehnung | die Antrags-Zeile trägt keine Kennung; sie bleibt `pending` und erscheint als Warnung im Log | Zeile in `cdc.administration_request` prüfen |
| `PCF-E8002` | Ablehnung | die Quelle des Antrags ist leer | Quelle angeben |
| `PCF-E8003` | Ablehnung | der Schemaname des Antrags ist leer | Schemaname angeben |
| `PCF-E8004` | Ablehnung | der Tabellenname des Antrags ist leer | Tabellenname angeben |
| `PCF-E8005` | Ablehnung | die Antragsart liegt außerhalb der bekannten Menge | eine der dokumentierten Funktionen `cdc.*` verwenden |
| `PCF-E8006` | Ablehnung | der Spaltenname eines Spaltenausschluss-Antrags ist leer | Spalte angeben |
| `PCF-E8010` | Ablehnung | der Regelname ist leer oder liegt außerhalb des Alphabets `a`–`z`, `0`–`9`, `_` mit 1 bis 63 Zeichen | Regelnamen korrigieren |
| `PCF-E8011` | Ablehnung | die Regelform `rule_spec` ist ungültig: kein JSON-Objekt, unbekannter Regeltyp oder Schlüssel, fehlender oder falsch typisierter Pflichtschlüssel | Regelform korrigieren |
| `PCF-E8012` | Ablehnung | die Routing-Regel oder ihr Zielname ist ungültig | Zielname und `order` korrigieren |
| `PCF-E8020` | Ablehnung | der Regelname ist je Tabelle bereits vergeben | andere Regel zuerst entfernen oder anderen Namen wählen |
| `PCF-E8021` | Ablehnung | die Quellspalte trägt bereits eine Regel | vorhandene Regel zuerst entfernen |
| `PCF-E8022` | Ablehnung | der Zielname kollidiert mit dem Zielnamen einer anderen Regel oder mit einer Spalte der Tabelle | anderen Zielnamen wählen |
| `PCF-E8023` | Ablehnung | die Tabelle führt keine Regel mit diesem Namen | Regelnamen prüfen |
| `PCF-E8024` | Ablehnung | die Spalte existiert nicht an der Quelle | Spaltennamen prüfen |
| `PCF-E8025` | Ablehnung | die Tabelle existiert nicht an der Quelle | Tabellennamen prüfen |
| `PCF-E8030` | Ablehnung | die `order` einer Routing-Regel ist je Tabelle bereits vergeben | andere `order` wählen |
| `PCF-E8031` | Ablehnung | die Bedingung (Spalte und Wert) kommt je Tabelle bereits vor | Bedingung ändern oder die vorhandene Regel entfernen |
| `PCF-E8032` | Ablehnung | die Lage der Regel ohne `when` verletzt die Ordnung: sie ist bereits vorhanden, trägt nicht die höchste `order`, oder eine Regel mit `when` liegt hinter ihr | `order` anpassen |
| `PCF-E8033` | Ablehnung | die Bedingungsspalte der Routing-Regel ist ausgeschlossen | andere Spalte wählen oder den Ausschluss aufheben |
| `PCF-E8034` | Ablehnung | die Spalte trägt eine Routing-Bedingung und lässt sich nicht ausschließen | Routing-Regel zuerst entfernen |
| `PCF-E8040` | Ablehnung | die Tabelle ist nicht aktiviert oder nicht Mitglied der Publication (Vorbedingung eines Backfills) | Tabelle aktivieren, Backfill erneut beantragen |
| `PCF-E8041` | Ablehnung | für die Tabelle besteht bereits ein aktiver Backfill-Run | Ende des Runs abwarten |
| `PCF-E8050` | Ablehnung | der Request-Body eines Aufrufs der HTTP-API ist kein gültiges JSON | Body prüfen |
| `PCF-E8051` | Ablehnung | ein Pflichtfeld oder Pflichtparameter eines API-Aufrufs fehlt oder ist leer (Kennung des Consumers, Quelle, Schema, Tabelle, Publication) | die fehlende Angabe ergänzen |
| `PCF-E8052` | Ablehnung | ein Query-Parameter liegt außerhalb der Menge, die der Endpunkt kennt | Parameternamen prüfen |
| `PCF-E8053` | Ablehnung | ein Wert eines API-Aufrufs ist nicht lesbar oder liegt außerhalb des zulässigen Bereichs (keine Ganzzahl, negative Dauer, Version oder `limit` kleiner 1) | Wert korrigieren |
| `PCF-E8054` | Ablehnung | eine Position ist kleiner als 1; eine Position 0 gibt es nicht | eine Position ab 1 angeben |
| `PCF-E8055` | Ablehnung | die zu bestätigende Position liegt vor der bereits bestätigten Position des Consumers | eine Position hinter der bestätigten wählen; die bestätigte Position nennt `GET /consumers/position` bzw. der RPC `GetConsumerPosition` |
| `PCF-E8056` | Ablehnung | die Endposition eines Lesebereichs liegt vor der Startposition | `from` und `to` prüfen |
| `PCF-E8057` | Ablehnung | die Position gehört zu einer anderen Quelle als der Stand des Consumers | Quelle der Position prüfen |

Warnungen stehen im Log des Feed-Containers mit dem Attribut `code`; eine Warnung
beendet den Prozess nicht und trägt keine Fehlerklasse:

| Code | Bereich | Bedeutung | Maßnahme |
|---|---|---|---|
| `PCF-W1001` | Erfassung und Replikation | das Wecksignal über NATS konnte nicht veröffentlicht werden (bei der Erfassung oder bei einem Backfill-Run); die Änderungen sind erfasst und über SQL und HTTP lesbar | NATS-Server und `CDC_NATS_URL` prüfen |
| `PCF-W1002` | Erfassung und Replikation | die Veröffentlichung einer Change im NATS-Vollinhalts-Stream ist fehlgeschlagen (Zeile des NATS-Adapters, mit dem Subjekt); die Change ist erfasst und über SQL lesbar, sie fehlt nur auf diesem Zustellweg | NATS-Server und Verbindung prüfen, die fehlende Change über `GET /changes` oder SQL lesen |
| `PCF-W1003` | Erfassung und Replikation | eine Change wurde im NATS-Vollinhalts-Stream übersprungen, weil Schema- oder Tabellenname leer ist oder ein für NATS-Subjekte reserviertes Zeichen (`.`, `*`, `>`) oder Leerraum enthält | Namen von Schema und Tabelle prüfen, die Change über SQL oder HTTP lesen |
| `PCF-W1004` | Erfassung und Replikation | der Zielname einer Change trägt ein reserviertes Zeichen oder Leerraum; die Change fehlt im Ziel-Subjekt, das Tabellen-Subjekt ist unberührt; Reserve: das erlaubte Alphabet der Zielnamen schließt diese Zeichen aus, über eine Routing-Regel ist die Warnung nicht erreichbar | Zielname der Routing-Regel korrigieren (siehe [Routing-Regel konfigurieren](#routing-regel-konfigurieren)) |
| `PCF-W1005` | Erfassung und Replikation | eine Change ist nicht als JSON kodierbar und wird nicht zugestellt: im NATS-Vollinhalts-Stream übersprungen, im Server-Sent-Events-Stream endet die Verbindung des Clients | Log sichern, die Change über SQL lesen, den Support kontaktieren (siehe [Support und Kontakt](#support-und-kontakt)) |
| `PCF-W1006` | Erfassung und Replikation | ein Zyklus des Replication-Streams ist an einer vorübergehenden Störung gescheitert und wird mit Backoff wiederholt | Erreichbarkeit der Quelle prüfen; dauert die Störung länger als das Wiederholungsfenster, endet der Prozess mit `PCF-E1002` (siehe [Neustart nach einem Fehler](#neustart-nach-einem-fehler)) |
| `PCF-W2001` | Backfill | das Aufräumen eines Backfill-Runs (Snapshot schließen, Schreibtransaktion zurückrollen) ist fehlgeschlagen; der Run behält sein Ergebnis | beim Schließen des Snapshots die Erreichbarkeit der Quelle prüfen, beim Rollback der Schreibtransaktion Erreichbarkeit und Rechte der CDC-Datenbank; die Zeile im Log nennt, welcher der beiden Schritte gescheitert ist |
| `PCF-W2002` | Backfill | die wartenden Backfill-Runs konnten nicht gelesen werden; der Durchgang endet und wird wiederholt | Erreichbarkeit und Rechte der CDC-Datenbank prüfen |
| `PCF-W2003` | Backfill | ein Backfill-Run ist nach seinem Versuch weiter `queued`; der Durchgang endet und wird wiederholt | Zeile des Runs in `cdc.backfill_status` prüfen, Erreichbarkeit der CDC-Datenbank prüfen |
| `PCF-W2004` | Backfill | der Zustand eines Backfill-Runs konnte nicht festgehalten werden; der Durchgang endet und wird wiederholt | Erreichbarkeit und Rechte der CDC-Datenbank prüfen, `cdc.backfill_status` lesen |
| `PCF-W2005` | Backfill | ein Backfill-Run ist fehlgeschlagen; `error_message` in `cdc.backfill_status` nennt Klasse und Meldungscode der Ursache | Fehlertext lesen, die Ursache nach ihrem Code beheben, Backfill neu beantragen |
| `PCF-W2006` | Backfill | ein Backfill-Run ist unterbrochen: beim Prozessstart waren Runs im Zustand `running`, oder ein Run endete ohne Abschluss | Backfill neu beantragen |
| `PCF-W3001` | Retention und Speicher | die periodische Messung des WAL-Rückstands ist fehlgeschlagen; die Erfassung läuft weiter, die nächste Messung folgt im Takt | Erreichbarkeit und Rechte der Quelle prüfen |
| `PCF-W3002` | Retention und Speicher | der WAL-Rückstand des Capture-Slots liegt über der Warnschwelle; die Erfassung läuft kontrolliert weiter | Rückstand und Ursache prüfen (siehe [WAL-Rückstand prüfen](#wal-rückstand-prüfen)) |
| `PCF-W3003` | Retention und Speicher | die periodische Bereinigung nach der Retention ist fehlgeschlagen; die Erfassung ist nicht betroffen, der nächste Takt versucht es erneut | Erreichbarkeit und Rechte der CDC-Datenbank prüfen |
| `PCF-W4001` | Verwaltung | das Wecksignal der Antrags-Queue ist gestört; die Anträge werden weiter im Takt gelesen | Erreichbarkeit der CDC-Datenbank und `CDC_ADMIN_DSN` prüfen |
| `PCF-W4002` | Verwaltung | die offenen Anträge konnten nicht gelesen werden; der nächste Durchlauf versucht es erneut | Erreichbarkeit der CDC-Datenbank und Rechte der Rolle `cdc_admin` prüfen |
| `PCF-W4003` | Verwaltung | die Frist des Vorlaufs beim Prozessstart ist abgelaufen; der Antrag bleibt `pending` | die Sperre auf der Tabelle lösen; der Antrag wird im nächsten Durchlauf erneut verarbeitet |
| `PCF-W4004` | Verwaltung | ein Antrag ist gescheitert und `failed` vermerkt; `error_message` in `cdc.administration_request` nennt Klasse und Meldungscode der Ursache | Fehlertext lesen, die Ursache nach ihrem Code beheben, den Antrag neu stellen |
| `PCF-W4005` | Verwaltung | ein Antrag ist abgelehnt worden; `error_message` nennt den Code der Ablehnung | den Antrag nach dem Code der Ablehnung korrigieren und neu stellen |
| `PCF-W4006` | Verwaltung | der Ausgang eines Antrags (`applied` oder `failed`) konnte nicht vermerkt werden; der Antrag bleibt `pending` und wird im nächsten Durchlauf erneut verarbeitet | Erreichbarkeit und Rechte der Rolle `cdc_admin` prüfen |
| `PCF-W4007` | Verwaltung | eine Zeile der Antrags-Queue trägt keine Kennung und wird übersprungen; sie bleibt `pending` | Zeile in `cdc.administration_request` prüfen |
| `PCF-W4008` | Verwaltung | ein Aufruf der HTTP- oder gRPC-API ist an einem unerwarteten Fehler gescheitert; der Aufrufer erhielt `500` bzw. `Internal` mit dem Code der Ursache, das Attribut `error` der Zeile nennt die Ursache | die Ursache nach dem Code der Antwort beheben, Aufruf wiederholen |
| `PCF-W6001` | Beobachtbarkeit und Transport | die Übertragung der Kennzahlen an den OTLP-Empfänger ist fehlgeschlagen (Status außerhalb `2xx`, Verbindung abgewiesen, Frist überschritten); Erfassung, Persistierung und Bestätigung laufen unverändert, der nächste Takt versucht es erneut; die Warnung steht beim ersten Fehlschlag und danach höchstens alle 5 Minuten | Erreichbarkeit, Adresse und Header des Empfängers prüfen (siehe [Metriken per OTLP übertragen](#metriken-per-otlp-übertragen)) |
| `PCF-W6002` | Beobachtbarkeit und Transport | die Sicht `cdc.metrics` konnte für den Export nicht gelesen werden; der Zyklus entfällt, der nächste Takt versucht es erneut | Erreichbarkeit der CDC-Datenbank und `CDC_READER_DSN` prüfen |

### Container startet nicht

**Ursache:** eine Pflicht-Umgebungsvariable fehlt oder ist falsch
formatiert.

**Lösung:**

1. Prüfen Sie die Container-Logs (`docker logs <container>`) — die
   Fehlermeldung nennt die betroffene Variable.
2. Prüfen Sie das Format von `CDC_TABLES`
   (`schema.tabelle=tabelle-id:schema-version-id`, kommagetrennt).
3. Prüfen Sie die Token-Listen `CDC_API_TOKENS_READER` und
   `CDC_API_TOKENS_ADMIN`: kein leeres Element, kein Leerraum in einem
   Element (Meldungscode `PCF-E2008`).
4. Prüfen Sie bei gesetztem TLS-Paar beide Angaben (`CDC_TLS_CERT_FILE` und
   `CDC_TLS_KEY_FILE`, Meldungscode `PCF-E2009`) und die beiden Dateien:
   lesbar für den Benutzer des Containers, gültiges PEM, ein zusammengehöriges
   Paar (Meldungscode `PCF-E2010`).
5. Prüfen Sie bei gesetztem Export die drei Angaben `CDC_OTLP_ENDPOINT`
   (Schema `http` oder `https`, ein Host; `PCF-E2011`), `CDC_OTLP_HEADERS`
   (`Schlüssel=Wert`, durch Komma getrennt; `PCF-E2012`, ohne Endpunkt
   `PCF-E2014`) und den Takt (ganze Zahl von 5 bis 3600; `PCF-E2013`).
6. Starten Sie den Container erneut.

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
(Orchestrator, Supervisor). Eine vorübergehend nicht verfügbare Quelle
(Fehlerklasse `transient`) wiederholt der Capture-Pfad vorher selbst mit
begrenztem Backoff: Anfangsverzögerung 2 s, Verdopplung bis 30 s je
Warteschritt, Gesamtfenster 5 Minuten. „Gestreamt“ heißt: der Server hat
`START_REPLICATION` bestätigt; ein Zyklus, der ab dieser Bestätigung
mindestens 30 s bis zu seinem Fehler gestreamt hat, setzt die Episode zurück.
Verbindungsaufbau, Katalog- und Slot-Abfragen zählen nie dazu, und ein Zyklus,
der den Streaming-Zustand nicht erreicht, setzt die Episode nicht zurück. Der
Aufbau eines Zyklus (Verbindung, Katalog- und Slot-Abfragen) hat eine Frist
von 30 s; ihr Ablauf ist ein wiederholbarer Fehler. Jede Wiederholung steht
als WARN mit Versuchszähler im Log, die Fortsetzung als INFO, sobald der
Zyklus den Streaming-Zustand erreicht. Wiederholt werden Fehler ohne SQLSTATE
(Verbindungsabbruch, Fristablauf), die SQLSTATE-Klassen 08, 40, 53, 55, 57 und
58 sowie der Code 25006 (schreibgeschützter Knoten während eines Failovers).
Nicht wiederholt werden Berechtigungsfehler (SQLSTATE 42501, Klasse 28) und
jede andere Server-Abweisung; der Prozess endet dann sofort mit Ausgang 1.
Erst nach erschöpftem Fenster
endet der Prozess mit Ausgang 1 und dem sichtbaren Fehlerzustand im
Heartbeat. Ein Neustart
setzt am zuletzt bestätigten Slot-Stand fort; keine bereits gespeicherte
Änderung geht dabei verloren.

**Nicht anwendbare Routing-Regel.** Endet der Prozess mit der Fehlerklasse
`schema`, weil eine Routing-Regel auf eine Change nicht anwendbar ist (der Log
nennt „Routing-Regel auf die Änderung nicht anwendbar", Regel und Spalte), hilft
der bloße Neustart nicht: dieselbe Change trifft dieselbe Regel erneut. Die
Abhilfe ist `cdc.remove_route` zu beantragen — der Antrag bleibt `pending`, solange
der Prozess steht — und den Prozess dann neu zu starten; offene Anträge werden im
Vorlauf des Starts verarbeitet, bevor die erste Transaktion der Tabelle
verarbeitet wird (Frist des Vorlaufs: 30 Sekunden), und die zuvor nicht
bestätigte Transaktion erscheint danach über `cdc.changes` ohne Ziel (Einzelheiten
und Ursprung: [Routing-Regel konfigurieren](#routing-regel-konfigurieren)). Das
gilt nur für diese Ursache: endete der Prozess, weil eine Spalte an der Quelle
entfernt wurde (inkompatible Schemaänderung, Log „Relation-Änderung nicht sicher
als Obermenge interpretierbar"), genügt das Entfernen der Regel nicht. Im
Backfill-Run ist die Regel dieselbe Ursache, run-lokal: der Run endet `failed`,
der Prozess läuft weiter, und die Abhilfe ist `cdc.remove_route` samt einem
neuen `cdc.backfill_table`-Antrag ohne Neustart.

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
| Transformationsregel | Eine deklarative Regel (`rename_column` oder `map_value`), die im Row Image eine Spalte umbenennt oder ihren Wert abbildet, ausgewertet vor der Persistierung; konfiguriert über `cdc.set_transformation`/`cdc.remove_transformation` |
| Routing-Regel | Eine Regel je Tabelle (`target`, `order`, optional `when`), die einer erfassten Change bei der Erfassung ein Zustellziel zuweist; der erste Treffer in aufsteigender `order` bestimmt das Ziel; konfiguriert über `cdc.set_route`/`cdc.remove_route` |
| Zustellziel (`route_target`) | Der benannte Kanal einer Change, vergeben von einer Routing-Regel und mit der Change gespeichert; `NULL` heißt „nicht geroutet"; ein Leser wählt es über `target` oder das Subjekt `cdc.route.<source_id>.<ziel>` |

## 9. Anhang

### Grenzwerte

- Lebenszeichen-Takt: 5 Sekunden; als veraltet gilt ein Lebenszeichen
  nach mehr als 15 Sekunden (Faktor 3).
- WAL-Rückstand-Messtakt: derselbe Takt wie das Lebenszeichen (5 Sekunden).
- WAL-Rückstand-Schwellen: Warnschwelle 100 MiB,
  Fehlerschwelle 1 GiB — betrifft ausschließlich die Fehlerklasse
  `replication`, Unterart Transport-/Verbindungsstörung (siehe
  [Fehlerklassen](#fehlerklassen)).
- OTLP-Export (festgelegt, nicht gemessen): Takt Default 60 Sekunden, erlaubt
  5 bis 3600; Frist je Versuch die Länge des Taktes, höchstens 10 Sekunden; die
  Warnung bei einem andauernden Fehlschlag höchstens alle 5 Minuten (siehe
  [Metriken per OTLP übertragen](#metriken-per-otlp-übertragen)).
- Ein Container-Lauf bindet genau eine Quelle.
- Zielname einer Routing-Regel (`target`): 1 bis 63 Zeichen, das
  erste aus `a`–`z` und `0`–`9`, die übrigen aus `a`–`z`, `0`–`9`, `_` und `-`
  (`[a-z0-9][a-z0-9_-]{0,62}`); ein Name außerhalb des Alphabets endet den Antrag
  `failed` (`Zielname ist ungültig`). Das Alphabet ist aus der Subjekt-Syntax
  abgeleitet, nicht gegen einen NATS-Server geprüft; ein Name mit `-` und `_` ist
  an einem realen NATS-Server ein einzelnes Token (siehe [Zugriff über den
  NATS-Vollinhalts-Stream](#zugriff-über-den-nats-vollinhalts-stream)).
- `order` einer Routing-Regel: eine positive ganze Zahl bis 2147483647
  (**Setzung** der Umsetzung, die Spezifikation nennt keine Obergrenze;
  gemessen in einer Compose-Umgebung (PostgreSQL 18.6): `2147483647` wird
  angenommen, `2147483648` endet `failed`). Die
  Spezifikation führt keine Obergrenze für die Zahl der Regeln je Tabelle.
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
  Stufe mit 200.000 Zeilen (Median von 3 Runs je Lauf) mal die Toleranz von
  600 s, auf eine Stelle abgerundet. In sechs Läufen auf demselben Host
  (gemessen bzw. übernommen; ein Lauf bei Last fremder
  Container) lag die Rate zwischen 4.504 und 8.933 Zeilen/s, die Rechnung ergab
  nach Rundung zwischen 2.000.000 und 5.000.000 Zeilen; der Wert 4.000.000 liegt
  innerhalb dieser Spanne und ist ein Startwert, den eine weitere Messung
  nachschärfen kann. Keine gemessene Stufe kopierte die Toleranzdauer; die Rate
  ist hochgerechnet. Sechs weitere Läufe auf demselben Host lagen in den Stufen
  ab 100.000 Zeilen zwischen 4.088 und 9.425 Zeilen/s je Run (übernommen). Die
  Zahl gilt für den Host und die
  Bedingungen der Messung (siehe unten). Breite Zeilen (`jsonb`, `bytea`),
  andere Hardware und eine andere Einfügeform sind für die Kopierrate
  ungemessen. Die Richtgröße folgt der Kopierdauer; der Speicher des
  Feed-Containers geht nicht ein, weil er nicht mit nennenswertem Betrag an der
  Zahl der Changes hängt (siehe *Backfill, Speicher des Feed-Containers*).
- **Backfill, gemessene Werte** (Messung mit `tools/bench-backfill.sh`;
  übernommen;
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

  Eine Nachmessung (gleicher Host, gleiche Stufen, Median von 3 Runs;
  übernommen) lieferte 8.382, 8.975 und
  8.559 Zeilen/s. Zwei Runs über je 1.000.000 Zeilen (`--full`, gleicher Host;
  übernommen) dauerten
  125,2 s und 122,8 s (7.986 und 8.141 Zeilen/s). Die
  mittlere Blockdauer liegt bei etwa 0,12 bis 0,13 s (abgeleitet: Dauer durch
  Blockzahl); die längste Dauer eines einzelnen Blocks und die Dauer des Commits
  sind nicht gemessen.
- **Backfill, Speicher des Feed-Containers.** Der Speicher des Feed-Containers
  hängt nicht mit nennenswertem Betrag an der Zahl der Changes, die für die
  Quelle in `cdc.change` stehen: gemessen 14,9 bis 17,6 MiB bei 1.000.000 bis
  3.000.000 Changes; der periodische Bereinigungslauf (alle 10 s, siehe
  [Aufbewahrung (Retention)](#aufbewahrung-retention)) liest sie seitenweise
  (10.000 Changes je Seite) und ohne Row Images. Gemessen wurde mit
  `tools/bench-backfill-memory.sh` (Host: Linux 6.8.0-139-generic, Docker
  29.8.1, 20 CPU, 31 GiB RAM, PostgreSQL 18
  mit Standard-Einstellungen; schmale Zeilen von 74 Bytes; drei Runs über je
  1.000.000 Zeilen hintereinander, ein frischer Feed-Container je Run,
  `--memory 6g`; zwei Läufe).
  Die Spitze des Feed-Containers (`memory.peak`, gedruckt 60 s nach dem Run) gegen
  die Zahl der Changes, die nach dem Run in `cdc.change` stehen:

  | Changes in `cdc.change` | Spitze in MiB (gemessen, zwei Läufe) |
  |---|---|
  | 1.000.000 | 14,9 und 15,1 |
  | 2.000.000 | 16,5 und 16,9 |
  | 3.000.000 | 17,3 und 17,6 |

  Zwischen 1.000.000 und 3.000.000 Changes liegt die Spitze um 2,4 und 2,5 MiB
  höher; das sind etwa 1,3 Bytes je Change, abgeleitet aus der Differenz der
  beiden Endpunkte (keine gemessene Steigung). Die Werte der drei Stufen
  überlappen einander nicht, der Sprung von 1.000.000 auf 2.000.000 Changes fällt
  mit dem Ausgangszustand des Containers zusammen (leeres `cdc.change` im ersten
  Run, gefülltes in den folgenden); die Ursache des Anstiegs ist nicht
  untersucht. Ein dritter Lauf lag bei 2.000.000 und
  3.000.000 Changes bei 16,6 und 16,8 MiB; bei 1.000.000 Changes stand die Spitze
  bei 31,8 MiB, davon 16,5 MiB Seiten-Cache des Containers (`memory.peak`
  schließt ihn ein; der Speicher des Prozesses lag dort bei 9,1 MiB, gemessen).
  Die Bereinigungs-Takte laufen in allen sechs Runs bis zum Ende weiter
  (gedruckt: 20 bis 23 Zeilen „Bereinigung gelaufen“ seit dem Start des
  Feed-Containers, keine fehlgeschlagene). Bei leerem `cdc.change` liegt die Spitze
  des Prozesses im Run bei 8,9 bis 10,5 MiB, von 10.000 bis 1.000.000 Zeilen
  (Blockgröße 1.000; gemessen). Ungemessen bleiben breite Zeilen, ein Speicherlimit des
  Containers unter 6 GiB und die laufende Erfassung als Quelle der Changes; die
  Seite trägt keine Row Images, ihr Bedarf hängt deshalb nicht an der Zeilenbreite
  (aus dem Aufbau der Seite abgeleitet, nicht gemessen).

  **Ältere Server-Versionen.** Die Server-Versionen `v0.1.0` bis `v0.1.2` lesen in
  jedem Takt alle Changes der Quelle samt Row Images in den Speicher (aus dem
  Quelltext der drei Versionen abgeleitet, nicht am Image der Version gemessen).
  Ihr Speicher wächst mit der Zahl der Changes: 1,03 bis 1,59 KiB je Change bei
  schmalen Zeilen (abgeleitet, ab 100.000 Changes; an einem Build mit dieser
  Lesung gemessen, nicht am Image der Version) und 1.082,7 bis 1.273,5 MiB bei
  1.000.000 Changes (gemessen).
- **Backfill, WAL-Rückstand des Capture-Slots** (Größe von
  `cdc_wal_retention_bytes`, siehe [WAL-Rückstand prüfen](#wal-rückstand-prüfen)).
  Der Feed bestätigt WAL ohne Inhalt für die Publication im Leerlauf seines
  Streams. Gemessen (mit `tools/bench-backfill.sh`, übernommen aus dem
  übernommen; in einem zweiten Lauf mit gleichem Ergebnis nachgemessen;
  Host und Tabellen wie oben, PostgreSQL 18; Rückstand im Abstand von 1 bis 2 s
  gelesen und auf ganze MiB gerundet): in allen neun Runs der Stufen mit 10.000,
  50.000 und 200.000 Zeilen (je 3 Runs) lag die Spitze des Rückstands im Run bei
  0 MiB, unmittelbar nach dem Run ebenfalls bei 0 MiB; das Skript schrieb
  zwischen den Runs nichts. Ohne die Bestätigung — gemessen an einem
  Stand, der WAL ohne Inhalt für die Publication nicht bestätigte (Stufe
  200.000 Zeilen, übernommen) — lag die Spitze im Run im Median bei 140 MiB (etwa 735
  Bytes je Zeile), bei 719 und 834 MiB in den zwei Runs über je 1.000.000 Zeilen,
  und der Rückstand blieb bis zu einem Commit auf einer aktivierten Tabelle
  bestehen (776 MiB unmittelbar nach dem ersten Run über 1.000.000 Zeilen); ein
  dritter Run über 1.000.000 Zeilen im selben CDC-Speicher überschritt vor
  seinem Ende die Fehlerschwelle von 1 GiB (Messwert 1.098.218.616 Bytes,
  Feed-Container mit Ausgang 1, Run `interrupted`). Ein Schreiber auf eine nicht
  aktivierte Tabelle erzeugte dort ohne Run 33,5 MiB WAL je 200.000 Zeilen (175
  Bytes je Zeile; übernommen).
  **Grenze der Ein-Transaktions-Form:** das vom Slot auf der Platte der Quelle
  **gehaltene** WAL (`restart_lsn`) und der Spill des Walsenders bleiben. Das
  gehaltene WAL erreichte in der Stufe mit 200.000 Zeilen im Median 140 MiB
  (Spitze im Run, übernommen; 31 MiB bei 50.000 und 7 MiB bei 10.000 Zeilen)
  und 141 MiB (35 MiB bei 50.000, 7 MiB bei 10.000 Zeilen; gemessen im zweiten
  Lauf), in den zwei Runs über je 1.000.000 Zeilen 782 und 1.613 MiB
  (übernommen). Der
  Walsender lagerte bei einem Run über 200.000 Zeilen 79 MB der offenen
  Transaktion aus (`spill_bytes` des Slots, `logical_decoding_work_mem` 64 MB;
  übernommen). Beides wächst mit der Größe der Tabelle; die Bestätigung im
  Leerlauf entlastet den Capture-Pfad, nicht die Platte der Quelle. Bemessen
  Sie den Plattenplatz der Quelle danach (abgeleitet: 1.613 MiB gehaltenes WAL
  bei 1.000.000 Zeilen von etwa 74 Bytes, gut das Zwanzigfache der
  Zeilenbytes).
- **Backfill, Wirkung auf die Live-Erfassung** (je ein Run über 200.000 Zeilen
  bei 100 Live-Änderungen/s in eine andere aktivierte Tabelle, dazu eine
  Referenz gleicher Dauer ohne Run; **jeder Vergleich ist ein Einzellauf ohne
  Wiederholung**): `cdc_capture_lag` im Run gegenüber ohne Run — erster Lauf
  (übernommen): 0,049 bis
  1,005 s (Median 0,48 s, 24 Proben) gegenüber 0,049 bis 1,022 s (Median 0,50
  s, 29 Proben); zweiter Lauf (gemessen):
  0,123 bis 1,074 s (Median 0,59 s, 21 Proben) gegenüber 0,041 bis 1,011 s
  (Median 0,40 s, 25 Proben); dritter Lauf (übernommen): 0,090 bis
  0,961 s (Median 0,456 s, 22 Proben) gegenüber 0,050 bis 1,008 s (Median
  0,413 s, 27 Proben). Die Mediane liegen in beiden Richtungen auseinander, die
  Maxima bei etwa 1 s; bei dieser Last und Größe ist aus diesen drei
  Einzelläufen kein Unterschied ableitbar.
- **Backfill, Schätzung der Zeilenzahl** (`pg_class.reltuples`, PostgreSQL 18,
  Tabelle mit 100.000 Zeilen, übernommen; dieselben Werte in zwei weiteren
  Läufen): eine frisch befüllte
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
| 1.1 | 2026-09-12 | Rollenspezifische Verbindungen: `CDC_SOURCE_DSN` entfällt, stattdessen `CDC_CAPTURE_DSN`, `CDC_ADMIN_DSN` und `CDC_READER_DSN`; Hinweis zum `REPLICATION`-Attribut ergänzt |
| 1.2 | 2026-09-12 | Fehlerklassen-Tabelle um die Klassen `transient`, `permission` und `internal` vervollständigt (alle sieben Klassen) |
| 1.3 | 2026-09-12 | Metrik `cdc_wal_retention_bytes` (WAL-Rückstand) ergänzt: periodische Messung, Log-Ausgabe, Abgrenzung gegen `cdc.metrics` |
| 1.4 | 2026-09-12 | Fehlerklasse `replication` in zwei Unterarten gefasst: Verletzung der Stream-Ordnung beendet sofort, Verbindungsstörungen laufen über die WAL-Schwellen (Warnung 100 MiB, Fehler 1 GiB) mit kontrollierter Fortsetzung bzw. Abbruch |
| 1.5 | 2026-09-13 | Neuer Sondermodus `diagnose` ergänzt (Abschnitt „Diagnose ausführen“); Rollen- und Variablentabellen angepasst |
| 1.6 | 2026-09-13 | Optionale YAML-Konfigurationsdatei (`CDC_CONFIG_FILE`) beschrieben; Pflichtangabe von `CDC_TABLES` präzisiert |
| 1.7 | 2026-09-13 | Neue Abschnitte „Tabelle live aktivieren“ und „Tabelle deaktivieren“ (`cdc.enable_table`/`cdc.disable_table`, asynchrone Antrags-Queue, Status abfragen) |
| 1.8 | 2026-09-13 | Neuer Abschnitt „Blockierende Consumer erkennen“ (`cdc.retention_blockers`) |
| 1.9 | 2026-09-13 | Metrik `cdc_storage_bytes` in „Metriken lesen“ ergänzt |
| 1.10 | 2026-09-13 | Die `diagnose`-Ausgabe nennt jetzt je Quelle den aktuell blockierenden Consumer (oder „kein Blocker“) und `cdc_storage_bytes` |
| 1.11 | 2026-09-13 | `CDC_NATS_URL`: Subjekt-Schema auf `cdc.changes.<source_id>.<schema>.<table>` (je Tabelle) korrigiert |
| 1.12 | 2026-09-14 | Beispielausgabe der Diagnose auf den Quellnamen `src-e2e` umgestellt |
| 1.13 | 2026-09-15 | Netzwerk-Zugriff ergänzt: Variablen `CDC_HTTP_ADDR`, `CDC_API_TOKEN_READER`, `CDC_API_TOKEN_ADMIN` und `CDC_GRPC_ADDR`; Abschnitt „Spalte vom Ausschluss konfigurieren“; Zugriffswege HTTP-/JSON-API, gRPC-Change-Stream und Server-Sent-Events |
| 1.14 | 2026-09-15 | Redaktionelle Korrekturen: Fähigkeitsliste der HTTP-API (Retention-Auslösung nur über die API) sowie Feldangaben bei gRPC (zehn Felder) und SSE präzisiert |
| 1.15 | 2026-09-15 | Änderungen lassen sich über `GET /changes` lesen (Parameter, Antwort, Fehlerfälle); die Zustellsemantik nennt die nicht streamende Form neben dem Live-Stream |
| 1.16 | 2026-09-15 | Neuer Abschnitt „Zugriff über das NATS-Wecksignal“: Subjekt-Schema, leerer Payload, Zustellsemantik, Ablauf (lauschen, dann über `GET /changes` holen) mit Beispielprogramm |
| 1.17 | 2026-09-17 | Go-Beispielprogramme für die HTTP-API und für SSE mit Startbefehl ergänzt; sie lesen Adresse und Token aus `CDC_HTTP_ADDR` und `CDC_API_TOKEN_READER` |
| 1.18 | 2026-09-17 | Konfigurationsdatei: neue Felder `http_addr` und `grpc_addr` mit Vorrangregel; Zugangsdaten-Felder auf sechs Schlüssel erweitert; `CDC_NATS_URL` und die beiden Token wirken auch bei geladener Datei aus der Umgebung |
| 1.19 | 2026-09-17 | Go-Beispielprogramm für den gRPC-Change-Stream ergänzt; es liest Adresse und Token aus `CDC_GRPC_ADDR` und `CDC_API_TOKEN_READER` |
| 1.20 | 2026-09-17 | Erstes C#-Beispielprogramm (HTTP-API) ergänzt; die Beispiele stehen jetzt je Sprache in einer Liste |
| 1.21 | 2026-09-17 | Erstes Kotlin-Beispielprogramm (HTTP-API) ergänzt |
| 1.22 | 2026-09-17 | C#- und Kotlin-Beispielprogramme für Server-Sent-Events ergänzt |
| 1.23 | 2026-09-17 | C#- und Kotlin-Beispielprogramme für das NATS-Wecksignal ergänzt |
| 1.24 | 2026-09-17 | C#-Beispielprogramm für den gRPC-Change-Stream ergänzt (`ChangeStream/StreamChanges`, Start über `make examples-csharp`); der gRPC-Code entsteht beim Bau aus der `.proto`-Datei, die dem Bau als zusätzlicher benannter Kontext übergeben wird |
| 1.25 | 2026-09-17 | Kotlin-Beispielprogramm für den gRPC-Change-Stream ergänzt (Start über `make examples-kotlin`, gRPC-Code beim Bau aus der `.proto`-Datei erzeugt); die Beispiele für den gRPC-Change-Stream sind damit in allen drei Sprachen vorhanden, die Beispielmatrix aus vier Zugriffsflächen und drei Sprachen ist vollständig |
| 1.26 | 2026-09-18 | Go-Beispiele starten jetzt über `make example-run-go SURFACE=<oberfläche>` statt über `go run` |
| 1.27 | 2026-09-18 | C#- und Kotlin-Beispiele starten jetzt über `make example-run-csharp SURFACE=<oberfläche>` bzw. `make example-run-kotlin SURFACE=<oberfläche>`; die Ziele starten die bereits gebauten Images |
| 1.28 | 2026-09-18 | Neuer Zustellweg „Zugriff über den NATS-Vollinhalts-Stream“ (Subjekte `cdc.stream.<...>`, gleiches Nachrichtenschema wie SSE, Aktivierungsbedingungen, Auswirkung der Authentifizierung auf das Wecksignal); neue Variable `CDC_NATS_STREAM_TOKEN`; Zugangsdaten-Felder der Konfigurationsdatei auf sieben erweitert |
| 1.29 | 2026-09-18 | Go-Beispielprogramm für den NATS-Vollinhalts-Stream ergänzt |
| 1.30 | 2026-09-18 | C#- und Kotlin-Beispielprogramme für den NATS-Vollinhalts-Stream ergänzt (abonnieren `cdc.stream.>` und geben jede Change aus; Start über `make example-run-csharp`/`make example-run-kotlin` mit `SURFACE=nats-stream`); alle Zustellwege haben damit Beispiele in drei Sprachen |
| 1.31 | 2026-09-19 | Metriken `cdc_changes_pending` und `cdc_errors_total` in „Metriken lesen“ genannt |
| 1.32 | 2026-09-19 | Kopffeld `Software-Version` verweist jetzt auf die maßgebliche Versionsquelle statt auf einen veralteten Wert |
| 1.33 | 2026-09-19 | C#-SDK (NuGet-Package `PgChangeFeed.Client`) für die HTTP-API beschrieben: alle zehn Fähigkeiten mit typisierten Anfragen, Antworten und Fehlerklasse |
| 1.34 | 2026-09-19 | C#-SDK für den gRPC-Change-Stream beschrieben: `PgChangeFeedGrpcClient.StreamChangesAsync` liefert die Changes mit allen zehn Feldern |
| 1.35 | 2026-09-19 | Python-SDK (PyPI-Package `pgchangefeed`) für die HTTP-API beschrieben: dieselben zehn Fähigkeiten, typisiert |
| 1.36 | 2026-09-20 | Kotlin-SDK `pgchangefeed-kotlin` für die HTTP-API beschrieben, samt Hinweis, dass der Bezug über GitHub Packages einen Token mit `read:packages` verlangt |
| 1.37 | 2026-09-20 | Kotlin-SDK für den gRPC-Change-Stream beschrieben (`streamChanges()` als `Flow`); Hinweis zum Token beim Bezug wiederholt |
| 1.38 | 2026-09-22 | C#-SDK für Server-Sent-Events und den NATS-Vollinhalts-Stream beschrieben (`PgChangeFeedSseClient` und `PgChangeFeedNatsStreamClient`, je `StreamChangesAsync` als `IAsyncEnumerable<Change>`); das Package `PgChangeFeed.Client` erscheint als Version 0.2.0 und deckt damit alle vier Zugriffswege ab |
| 1.39 | 2026-09-22 | Kotlin-SDK für Server-Sent-Events und den NATS-Vollinhalts-Stream beschrieben (`PgChangeFeedSseClient` und `PgChangeFeedNatsStreamClient`, je `streamChanges()` als `Sequence<Change>`); `pgchangefeed-kotlin` erscheint als Version 0.2.0; der Bezug über GitHub Packages verlangt weiterhin einen Token mit `read:packages` |
| 1.40 | 2026-09-23 | Python-SDK deckt zusätzlich den gRPC-Change-Stream ab (`PgChangeFeedGrpcClient`) |
| 1.41 | 2026-09-23 | Python-SDK-Absatz im Abschnitt zum gRPC-Change-Stream ergänzt (`stream_changes()` mit `timeout`-Parameter) |
| 1.42 | 2026-09-23 | Python-SDK für Server-Sent-Events beschrieben (`PgChangeFeedSseClient.stream_changes()`) |
| 1.43 | 2026-09-23 | Python-SDK für den NATS-Vollinhalts-Stream beschrieben (`PgChangeFeedNatsStreamClient`); das Package `pgchangefeed` erscheint als Version 0.2.0, alle vier Zugriffswege sind in allen drei SDK-Sprachen abgedeckt |
| 1.44 | 2026-09-24 | Feld `origin` ergänzt: `cdc.changes` und `GET /changes` unterscheiden `wal` und `backfill` (fehlender Wert liest als `wal`); die Live-Zustellwege tragen das Feld nicht |
| 1.45 | 2026-09-24 | „Schema aktualisieren“ beschreibt den Ablauf bei geänderter View-Signatur: Schema-Rollout vor dem Container-Tausch, automatisches Löschen der View, Lesefenster für SQL-Leser, Verhalten bei abhängigen Objekten oder Abbruch |
| 1.46 | 2026-09-24 | Hinweis zu Rechten beim Neuanlegen der View: nur `cdc_reader` erhält sein `SELECT`-Recht zurück, eigene Rechte für andere Rollen setzt der Betreiber erneut; Schema `cdc` wird fest adressiert |
| 1.47 | 2026-09-24 | Rechte der drei Rollen beschrieben: `cdc_admin` verarbeitet die Antrags-Queue (`CDC_ADMIN_DSN` in den Umgebungsvariablen aufgeführt); der Rollout setzt die Rechte bei jedem Lauf; ohne das Recht bleiben Anträge `pending` |
| 1.48 | 2026-09-24 | Neuer Abschnitt „Bestand als Backfill überführen“ (`cdc.backfill_table`, View `cdc.backfill_status`, geschätzte Zeilenzahl, Verhalten beim Neustart, Sichtbarkeit, Schema-Version, Vorbedingungen); „Änderungen lesen“ erklärt Position und `limit`; „Diagnose ausführen“ zeigt Backfill-Runs; Rechte von `cdc_capture` und `cdc_reader` in „Schema aktualisieren“, Rollen, Umgebungsvariablen (zwei DSN-Zeilen), Fehlerklassen, Glossar und Grenzwerte ergänzt |
| 1.49 | 2026-09-24 | Backfill-Abschnitt präzisiert: Antrag unter `cdc_admin`, `cdc.backfill_status` unter `cdc_reader` lesen; der Run beginnt mit dem Snapshot; Anträge anderer Quellen bleiben für diese Instanz `pending` |
| 1.50 | 2026-09-24 | Backfill: Startposition eines frisch registrierten Consumers beschrieben (`offset` 0, nicht bestätigt); `rows_copied` eines unterbrochenen Runs nennt den zuletzt gesicherten Fortschritt |
| 1.51 | 2026-09-24 | Backfill: Lesesperre der Tabelle während des Runs und ihre Wirkung auf DDL beschrieben; Verhalten bei umgeschriebener Tabelle (Run `failed`, Klasse `transient`, neuer Antrag als Abhilfe) |
| 1.52 | 2026-09-24 | Backfill: Wirkung der Tabellensperre vollständig beschrieben (wartende DDL staut Leser, Schreiber und Administration; der Run wartet ohne Zeitgrenze auf offene `ACCESS EXCLUSIVE`-Transaktionen) |
| 1.53 | 2026-09-25 | Backfill: Wann `warn_estimated_size` und `warn_duration` gesetzt werden, WAL-Rückstand als vierte Vorbedingung; „Grenzwerte“ nennt Toleranz der Kopierdauer, Richtgröße und gemessene Werte |
| 1.54 | 2026-09-25 | Backfill: Ursache und Abhilfen des WAL-Rückstands beschrieben (Commit auf einer aktivierten Tabelle, Datei-Feld `wal_retention_error_bytes`); Messwerte in „Grenzwerte“ ergänzt |
| 1.55 | 2026-09-25 | „Grenzwerte“: Herkunft der Messzahlen gekennzeichnet, Richtgröße um eine weitere Messung ergänzt |
| 1.56 | 2026-09-25 | WAL ohne Inhalt für die Publication wird im Leerlauf des Streams bestätigt und lässt `cdc_wal_retention_bytes` (vom Feed noch nicht bestätigtes WAL) nicht wachsen; der Backfill-Abschnitt nennt den WAL-Rückstand als Punkt ohne Abbruch über die Fehlerschwelle und die offene Schreibtransaktion des Runs als verbleibende Last; „WAL-Rückstand prüfen“ und „Grenzwerte“ entsprechend angepasst |
| 1.57 | 2026-09-25 | „Grenzwerte“: Herkunft weiterer Messzahlen gekennzeichnet, Nachmessung ergänzt (Rückstand 0 MiB in neun Runs), Richtgröße als Spanne von 2.000.000 bis 5.000.000 Zeilen angegeben |
| 1.58 | 2026-09-25 | Die drei SDK-Packages tragen das Feld `origin` bei `GET /changes`; eine Antwort ohne das Feld liest als `wal` |
| 1.59 | 2026-09-25 | Regel für `origin` in den SDK-Packages präzisiert: fehlendes Feld oder `null` liest als `wal`, jeder andere Wert kommt unverändert an; das Python-Package trägt SSE und den NATS-Vollinhalts-Stream |
| 1.60 | 2026-09-25 | Speicherbedarf des Feed-Containers im Backfill gemessen: die Spitze hängt an der Zahl der Changes in `cdc.change`, nicht an der Tabellengröße; Hinweise zur Bemessung des Speicherlimits in „Grenzwerte“ und im Backfill-Abschnitt |
| 1.61 | 2026-09-25 | Messzahlen und Herkunftsangaben zum Speicherbedarf berichtigt (Höchstwert je Change, Wirkung von `GOGC=25`, Messung mit sehr vielen Changes); die Server-Versionen `v0.1.0` bis `v0.1.2` tragen das beschriebene Verhalten |
| 1.62 | 2026-09-25 | Der Bereinigungslauf liest die Kandidaten seitenweise (10.000 Changes, nicht atomar, Consumer-Positionen einmal je Durchlauf) ohne Row Images; „Aufbewahrung (Retention)“ beschreibt die Seiten; Speicherspitze 14,9 bis 17,6 MiB bei 1.000.000 bis 3.000.000 Changes, die Bemessung des Speicherlimits im Backfill-Abschnitt und in „Grenzwerte“ folgt dieser Messung; die Server-Versionen `v0.1.0` bis `v0.1.2` lesen noch alle Changes je Durchlauf |
| 1.63 | 2026-09-25 | Aussagen zum Speicher angeglichen: der Bedarf hängt an der Seitengröße, nicht nennenswert an der Zahl der Changes (über 3.000.000 Changes nicht gemessen) |
| 1.64 | 2026-09-25 | „Aufbewahrung (Retention)“: Changes, die ein Durchlauf löscht, sind für einen erstmals bestätigenden Consumer verloren; ein Durchlauf kann eine Transaktion in mehreren Schritten löschen |
| 1.65 | 2026-09-25 | Kotlin-SDK ohne Token beziehbar: Cloudsmith als anonym lesbarer Bezugsweg neben GitHub Packages; `pgchangefeed-kotlin` erscheint als Version 0.2.2 |
| 1.66 | 2026-09-27 | Zusage der Bestätigung im Leerlauf präzisiert: bestätigt wird WAL ohne Inhalt für die Publication zwischen den Schreibvorgängen; neuer Absatz „Grenze der Bestätigung im Leerlauf“ (ein Stoß über `wal_retention_error_bytes` bleibt bis zur nächsten Bestätigung im Rückstand, Abhilfe: Fehlerschwelle heben) |
| 1.67 | 2026-09-27 | Neuer Abschnitt „Transformationsregel konfigurieren“ (`cdc.set_transformation`/`cdc.remove_transformation`, Regeltypen, Konfliktfreiheit, Wirkung, Dauerhaftigkeit, Nichtanwendbarkeit samt Abhilfe, Backfill-Bezug, Form auf allen Zustellwegen); Rollen, Fehlerklassen und Glossar ergänzt |
| 1.68 | 2026-09-27 | Fehlertext eines `::jsonb`-Aufrufs von `cdc.set_transformation` auf den tatsächlichen Wortlaut korrigiert; Absatz „Zeilen, die kein Antrag sind“ um ein Beispiel ergänzt |
| 1.69 | 2026-09-28 | „Transformationsregel konfigurieren“ als Anleitung überarbeitet (Voraussetzung, nummeriertes Vorgehen, Ergebnis, zwei Fehlerblöcke mit Ursache und Lösung); interne Verweise und Messrohwerte aus dem Fließtext entfernt |
| 1.70 | 2026-09-28 | Beispiele der HTTP-API in allen drei Sprachen decken jetzt alle zehn Fähigkeiten ab: ein Verb-Flag (Standard `tables`); Hinweis auf `make example-transformation-demo` |
| 1.71 | 2026-09-28 | Neuer Abschnitt „Zugriff über die gRPC-Verwaltungs-API“: die Fähigkeiten der HTTP-API als gRPC-Dienst `Administration`, Rechteklassen, Fehlercodes; `CDC_GRPC_ADDR` gilt für beide gRPC-Dienste |
| 1.72 | 2026-09-28 | gRPC-Verwaltungs-API um den RPC `ReadChanges` ergänzt (zehn RPCs; Filter, Bereiche und Nachrichtenschema wie `GET /changes`; leerer Treffer ist kein `NotFound`) |
| 1.73 | 2026-09-28 | Diagnose auch über `GET /diagnose` und den gRPC-RPC `Diagnose` (elfter RPC) abrufbar, gleichwertig zum CLI-Aufruf; Antwortschema beschrieben |
| 1.74 | 2026-09-28 | gRPC-Change-Stream und SSE lassen sich nach Tabelle filtern: optionales `schema`/`table`-Paar mit derselben Kombinatorik wie `GET /changes` |
| 1.75 | 2026-09-28 | Go-Beispiel für gRPC deckt mit dem Flag `-verb` zusätzlich alle elf Verwaltungs-RPCs ab; der Stream-Modus nimmt den Filter `-schema`/`-table` entgegen |
| 1.76 | 2026-09-28 | C#-Beispiel für gRPC deckt wie das Go-Beispiel alle elf Verwaltungs-RPCs ab (`--verb`, Filter `--schema`/`--table`) |
| 1.77 | 2026-09-28 | Kotlin-Beispiel für gRPC deckt ebenfalls alle elf Verwaltungs-RPCs ab; die gRPC-Verwaltungs-API hat damit Beispiele in allen drei Sprachen |
| 1.78 | 2026-09-28 | `EnableTable`/`DisableTable` über HTTP und gRPC wirken unmittelbar mit der Antwort auf den laufenden Erfassungsprozess, ohne Neustart |
| 1.79 | 2026-09-28 | C#-SDK deckt die volle gRPC-Fläche ab: `PgChangeFeedAdministrationClient` mit allen elf RPCs, Stream-Filter `schema`/`table` |
| 1.80 | 2026-09-28 | Python-SDK deckt die volle gRPC-Fläche ab: `PgChangeFeedAdministrationClient` mit allen elf RPCs, Stream-Filter `schema`/`table` |
| 1.81 | 2026-09-28 | Kotlin-SDK deckt die volle gRPC-Fläche ab; die gRPC-Verwaltungs-API ist damit in allen drei SDK-Sprachen verfügbar |
| 1.82 | 2026-09-30 | „Neustart nach einem Fehler“: begrenzte Wiederholung bei Fehlerklasse `transient` im Erfassungspfad (Rücksetzung nach einem Zyklus von mindestens 30 s, Log-Meldungen mit Versuchszähler); keine Wiederholung bei Berechtigungsfehlern und Server-Abweisungen |
| 1.83 | 2026-09-30 | Wiederholung präzisiert: ein Zyklus zählt als gestreamt ab Bestätigung von `START_REPLICATION` (mindestens 30 s), der Aufbau hat eine Frist von 30 s; wiederholt werden nur bestimmte SQLSTATE-Klassen, jede andere Server-Abweisung beendet sofort |
| 1.84 | 2026-10-01 | Neuer Abschnitt „Routing-Regel konfigurieren“ (Voraussetzung `cdc_admin`, `cdc.set_route`/`cdc.remove_route`, Regelform, Konfliktfälle, Fehlerklasse `schema`, Ziel lesen, Hinweise zu festem Label, nicht treffenden Changes und `DELETE`); der Filter `target` an allen Zugriffswegen, das Zusatz-Subjekt `cdc.route.<source_id>.<ziel>` im NATS-Vollinhalts-Stream, die Spalte `route_target` in `cdc.changes`; Fehlerklassen („Neustart nach einem Fehler“), Rollen, Glossar, „Grenzwerte“ (Kosten der zweiten Veröffentlichung) und „Schema aktualisieren“ nachgezogen; die Transformationsregeln trennen die entfernte Spalte von der nicht anwendbaren Regel |
| 1.85 | 2026-10-01 | Parameter `target` in den drei SDK-Packages und in den Go-, C#- und Kotlin-Beispielen (HTTP-Lesen, gRPC-Stream, SSE, `ReadChanges`, NATS-Vollinhalts-Stream) |
| 1.86 | 2026-10-01 | Redaktionelle Korrektur: `target` ist am SSE-Client der einzige Filter-Parameter; der Verweis auf einen offenen Folge-Schritt entfällt |
| 1.87 | 2026-10-02 | Der SSE-Client der SDK-Packages und die SSE-Beispiele in Go, C# und Kotlin nehmen zusätzlich `schema` und `table` entgegen (Flags `-schema`/`-table` bzw. `--schema`/`--table`) |
| 1.88 | 2026-10-02 | Die Erzeugnisse von `make schema-rollout` (Pflicht-Report, Rollback-Artefakt, Precheck-Report) liegen in `SCHEMA_ARTEFACT_DIR` (Standard `.tmp/schema-rollout`); die Aufbewahrung je Rollout liegt beim Betreiber; der Lauf mountet den Arbeitsbaum nicht |
| 1.89 | 2026-10-03 | Die Ausgabe von `diagnose`, die Meldungen von `make schema-rollout` und der Fehlertext der Konfigurationsdatei tragen keine Anforderungs- oder Entscheidungskennung mehr; die Beispiele zeigen den tatsächlichen Ausgabetext |
| 1.90 | 2026-10-03 | Neuer Abschnitt „Meldungscodes“ (Fehlerbehebung) mit dem Katalog aller Codes: Fehlerzeilen tragen jetzt `Fehlerklasse <Klasse> [<Code>]: …`, `error_message` von Backfill-Runs `<Klasse> [<Code>]: …` und abgelehnte Anträge `abgelehnt [<Code>]: …`, die Meldungen von `make schema-rollout` `FEHLER [<Code>]: …`; der Code ist stabil, der Text nicht; ein `grep` auf den bisherigen Wortlaut `Fehlerklasse schema:` greift nicht mehr, ein `grep` auf die Klasse oder den Code schon |
| 1.91 | 2026-10-03 | Warnungen im Log tragen das Attribut `code` mit einem Meldungscode wie `PCF-W3002` (Katalog der Warnungen im Abschnitt „Meldungscodes“); `cdc.process_heartbeat` und `cdc.heartbeat` tragen die neue Spalte `error_code` mit dem Meldungscode des Fehlerzustands (in der Sicht hinter `age_seconds`, `NULL` im Normalbetrieb); die Zeile „Fehlerzustand“ von `diagnose` nennt Klasse und Code (`schema [PCF-E4003]`), `GET /diagnose` und der RPC `Diagnose` tragen das Feld `error_code`; die neue Spalte kommt mit dem Schema-Rollout vor dem Container-Tausch |
| 1.92 | 2026-10-03 | Fehlerantworten der HTTP-API tragen das Feld `code` mit einem Meldungscode neben `error` (`401` und `403` ohne Feld); Fehlerstatus der gRPC-API (`InvalidArgument`, `NotFound`, `Internal`) tragen das Statusdetail `ErrorInfo` mit dem Code als `reason` und `pg-change-feed` als `domain`; der Katalog führt die neuen Codes `PCF-E8050` bis `PCF-E8057` für abgelehnte API-Aufrufe und `PCF-W4008` für die Warnung eines gescheiterten API-Aufrufs; Statuscodes und Fehlertexte bleiben unverändert |
| 1.93 | 2026-10-03 | Die Fehlertypen der drei SDK-Packages (C#, Kotlin, Python) tragen den Meldungscode des Servers als Eigenschaft (`MessageCode`, `messageCode`, `message_code`): beim HTTP- und beim SSE-Client aus dem Feld `code`, beim gRPC-Verwaltungs-Client aus dem Statusdetail; ohne Code auf dem Draht ist die Eigenschaft leer; das Feld `error_code` der Diagnose kommt unverändert an |
| 1.94 | 2026-10-03 | Die SDK-Eigenschaft für den Meldungscode gilt ab Package-Version 0.6.0 (Abschnitte zu Fehlerantworten und zur Fehlerform der gRPC-Verwaltungs-API); trägt ein Fehlerkörper `code` als String, aber `error` nicht als String, liefern alle drei SDKs den Code und den rohen Körper als Fehlertext |
| 1.95 | 2026-10-04 | Mehrere API-Token je Klasse: die neuen Variablen `CDC_API_TOKENS_READER` und `CDC_API_TOKENS_ADMIN` (kommagetrennte Listen) gelten neben den bisherigen Variablen, auf HTTP und gRPC gleich; neuer Abschnitt „API-Token in zwei Neustarts wechseln“; eine Liste mit leerem Element oder Leerraum verhindert den Start (neuer Meldungscode `PCF-E2008`); die Zugangsdaten-Schlüssel der Konfigurationsdatei umfassen jetzt neun Schlüssel |
| 1.96 | 2026-10-04 | Die HTTP- und die gRPC-Schnittstelle lassen sich mit einem gemeinsamen TLS-Paar verschlüsseln: die neuen Variablen `CDC_TLS_CERT_FILE` und `CDC_TLS_KEY_FILE` (Datei-Felder `tls_cert_file` und `tls_key_file`, Pfade ohne Zugangsdaten-Charakter), neuer Abschnitt „Schnittstellen mit TLS verschlüsseln“ (TLS ab Version 1.2, kein Klartext auf derselben Adresse, kein Client-Zertifikat, Zertifikatswechsel nur mit Neustart); ein unvollständiges oder nicht ladbares Paar verhindert den Start (neue Meldungscodes `PCF-E2009` und `PCF-E2010`); die gRPC-Beispielprogramme verbinden im Klartext und die Optionen der Client-Pakete bieten keine TLS-Einstellung |
| 1.97 | 2026-10-04 | Der Feed-Container kann seine Kennzahlen periodisch per OpenTelemetry-Protokoll (OTLP/HTTP, Protobuf) an einen Empfänger übertragen: die neuen Variablen `CDC_OTLP_ENDPOINT`, `CDC_OTLP_HEADERS` und `CDC_OTLP_INTERVAL_SECONDS` (Datei-Feld `otlp_interval` für den Takt), neuer Abschnitt „Metriken per OTLP übertragen“ mit Kennzahlen, Einheiten und dem Hinweis, dass alle Kennzahlen Momentstände (Gauge) sind, auch die mit dem Namensteil `_total`; ohne Endpunkt bleibt der Export aus; ein Ausfall des Empfängers beeinträchtigt Erfassung und Health nicht und erzeugt die Warnungen `PCF-W6001` und `PCF-W6002` (neuer Warn-Bereich 6 „Beobachtbarkeit und Transport“); eine ungültige Export-Konfiguration verhindert den Start (neue Meldungscodes `PCF-E2011` bis `PCF-E2014`; die Meldung nennt nie den Wert von Endpunkt oder Header, beim Takt den eingegebenen Wert); die zugangsdaten-tragenden Schlüssel der Konfigurationsdatei, die dort nicht stehen dürfen, wachsen von neun auf elf (`otlp_endpoint`, `otlp_headers`) |
| 1.98 | 2026-10-04 | Metrik-Export gegen einen echten OpenTelemetry-Collector gemessen: die Einheit von `cdc_consumer_lag` ist `By` (WAL-Strecke in Bytes, nicht Datenvolumen; Abschnitte „Metriken lesen“, „Metriken per OTLP übertragen“ und die Diagnose-Beschreibung); der Abschnitt zum Export nennt den `https`-Empfänger (CA-Bündel im Image, selbstsigniertes Zertifikat über `SSL_CERT_FILE`, ohne sie Warnung `PCF-W6001`) und den Benutzerteil der Endpunkt-URL (wird als Basic-Authorization gesendet, ein `Authorization`-Header aus `CDC_OTLP_HEADERS` gewinnt, keine Zugangsdaten in Log oder Fehlertext) |
| 1.99 | 2026-10-05 | Die gRPC-Beispielprogramme in Go, C# und Kotlin verbinden über TLS, wenn ihnen ein Vertrauensanker genannt wird (Flag `-ca-file` bzw. `--ca-file`, Umgebungsvariable `CDC_TLS_CA_FILE`; ohne Angabe Klartext wie bisher); der Abschnitt „Schnittstellen mit TLS verschlüsseln“ nennt die Aufrufform und die Messung gegen einen Feed-Container mit TLS-Paar |
