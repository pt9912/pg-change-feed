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
   gegenüber PostgreSQL. Im Leerlauf des Streams — keine Quelltransaktion
   zwischen BEGIN und COMMIT — bestätigt die Application zusätzlich das
   WAL-Ende, das die Quelle in ihrer Keepalive-Nachricht nennt, sobald es hinter
   der zuletzt bestätigten Position liegt; diese Bestätigung braucht keine
   Persistenz, weil das WAL bis dahin keinen zu speichernden Change der
   Publication trägt. Inmitten einer Quelltransaktion und bei einem WAL-Ende
   nicht hinter der bestätigten Position bestätigt der Stream nichts.

Zentrale Invariante:

    ACK(position) => durable(all changes <= position)

Sie gilt auch für die Bestätigung im Leerlauf: jede Transaktion mit Changes der
Publication und Commit vor dem bestätigten WAL-Ende ist bereits gespeichert, und
eine Transaktion, die davor begann und danach committet, liefert die Quelle
nach der bestätigten Position vollständig.

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
4. Die Bereinigungsmenge wird seitenweise bestimmt (10.000 Kandidaten je
   Seite); eine Seite trägt je Change nur Kennung, Commit-Position und
   Commit-Zeitpunkt, keine Row Images. Der Arbeitsspeicher eines
   Bereinigungslaufs hängt an der Seitengröße und nicht mit nennenswertem
   Betrag an der Zahl der gespeicherten Changes.

Retention ist eine Domain Policy; sie liegt nicht in einem Adapter.

### LH-FA-SCH-004.a — Schema-Änderungs-Klassifikation

**Eingabe:** Relation Metadata aus dem Replication Stream. **Ausgabe:**
TableSchema-/SchemaVersion-Modell (SPEC-004) oder sichtbarer Fehler.

Relation Metadata wird in technologieunabhängige
TableSchema-/SchemaVersion-Modelle übersetzt; jeder Change referenziert
eine Schema-Version ([`LH-FA-SCH-005`](lastenheft.md)). Nicht sicher interpretierbare
Schemaänderungen führen zu einem sichtbaren Fehler (Fehlerklasse `schema`,
§4) statt stiller Fehlinterpretation ([`LH-FA-SCH-004`](lastenheft.md)).

### LH-FA-CAP-009.a — Backfill-Mechanismus

**Eingabe:** Bestand einer aktivierten Quelltabelle zum Startzeitpunkt eines
Backfill-Laufs („Run"). **Ausgabe:** Backfill-Changes über denselben
Lesezugriffsweg wie WAL-erfasste Changes.

Die Lastenheft-Fähigkeit ([`LH-FA-CAP-009`](lastenheft.md)) wird durch einen
ausdrücklich ausgelösten Run erfüllt. Die folgenden Sätze sind Zusagen an die
Umsetzung, keine Messergebnisse.

- **Auslösung.** Ausdrücklich, nie implizit von der Aktivierung: die
  SQL-Funktion `cdc.backfill_table(source_id, schema, table)` schreibt
  ausschließlich einen Antrag der Antragsart `backfill` (`SPEC-019`) und
  sendet `pg_notify`. Ein Run setzt eine aktivierte Tabelle mit laufender
  Bindung voraus; eine verletzte Vorbedingung endet den Antrag `failed` mit
  Text, ebenso ein zweiter Antrag für dieselbe Tabelle, solange ein Run dieser
  Tabelle `queued` oder `running` ist.
- **Annahme und Aufnahme.** Die Administrations-Verarbeitung nimmt den Antrag
  an, sie führt ihn nicht aus: in **einer** Transaktion entsteht die Run-Zeile
  `queued` (`SPEC-029`) und der Antrag wird `applied` vermerkt — `applied`
  heißt hier „angenommen". Ein einzelner Backfill-Worker (eine Goroutine, ein
  Run zugleich, keine Parallelisierung) nimmt `queued`-Runs seiner Quelle in
  Antragsreihenfolge beim Prozessstart und bei jedem Wecksignal auf; er prüft
  vor der Ausführung Bindung und Publication-Mitgliedschaft erneut (Abweichung
  endet den Run `failed`, Fehlerklasse `configuration`, ohne Slot und ohne
  Kopie).
- **Mechanismus.** Bulk-Copy des Tabellenbestands im Snapshot eines je Run
  angelegten temporären logischen Replication Slots (`pgoutput`, mit
  Snapshot-Export). Der `consistent_point` des Slots ist die Position `X` des
  Runs; der Bestand wird in einer `REPEATABLE READ`-Transaktion gelesen, die
  den exportierten Snapshot importiert, über einen Cursor in Blöcken
  begrenzter Größe, ohne Datei-, Dump- oder Zwischenspeicher. Der Import steht
  unter einer Lesesperre der Tabelle; wurde die Tabelle zwischen
  Snapshot-Export und Sperre umgeschrieben, endet der Run `failed`
  (Fehlerklasse `transient`) ohne Change, ein neuer Antrag beginnt neu. Der
  Slot besteht nur, bis der Snapshot importiert und die Replication-Verbindung beendet ist.
  Alle Blöcke werden in **einer** Store-Transaktion geschrieben, die einmal am
  Ende zusammen mit dem Run-Zustand `completed` committet. Eine leere Tabelle
  endet `completed` mit 0 Zeilen und schreibt keine Transaktion. Die Rolle des
  Capture-Prozesses braucht `SELECT` auf die Quelltabelle — eine
  Betriebs-Vorbedingung wie `REPLICATION`, keine Vergabe durch die Software;
  ihr Fehlen endet den Run `failed` (Fehlerklasse `permission`).
- **Markierung.** Jeder Backfill-Change trägt `origin = backfill`
  (`SPEC-002`), `operation = INSERT`, kein `old_data` und ein **vollständiges**
  `new_data` (auch Spalten, die ein `UPDATE`-Image als Abwesenheit tragen kann,
  [`LH-FA-CAP-008`](lastenheft.md)). Das Row Image ist inhaltsgleich
  (Schlüsselmenge und Werte) dem WAL-Image derselben Zeile: dieselbe
  Bild-Konstruktion, ausgeschlossene Spalten
  ([`LH-FA-CFG-005`](lastenheft.md)) und generierte Spalten fehlen, die
  Transformationsregeln der Tabelle ([`LH-FA-CFG-007`](lastenheft.md),
  `SPEC-030`) wirken, `NULL` entfällt, Werte im Text-Stand der Quelle. Die
  Schema-Version einer Backfill-Change ist die zum Run-Start aktuelle Version
  der Tabelle; sie unterscheidet, sie beschreibt die Bild-Spalten nicht.
- **Position und Ordnung.** Alle Blöcke eines Runs liegen auf der Position
  `X`. Jeder Block ist eine eigene synthetische Transaktion mit der Kennung
  `0bf-<run-id>-<Blocknummer>` (Blocknummer achtstellig, null-aufgefüllt) und
  der Sequenz `1…B`; die Kennung sortiert wegen des führenden `0` vor jeder
  WAL-Transaktions-Kennung derselben Position. Die Change-Kennung folgt der
  Regel des WAL-Pfads (`<Transaktions-ID>-<Sequenz>`). Die Lese-Ordnung
  bleibt (`commit_position`, `transaction_id`, `sequence`)
  ([`LH-FA-REA-004`](lastenheft.md)); innerhalb des Backfills ist sie die
  Cursor-Reihenfolge, einmal vergeben und stabil.
- **Überlappungs-Verhalten.** Jeder Commit mit Position ≤ `X` steckt im
  Bestand, jeder Commit mit Position > `X` kommt über den WAL-Pfad und sortiert
  hinter dem Backfill; die Bindung der Tabelle steht vor der Slot-Anlage, es
  entsteht **keine Lücke**. Eine Zeile, die im Fenster (Aktivierung, `X`]
  geändert wurde und deren WAL-Change persistiert ist, erscheint **doppelt**:
  der WAL-Change davor, der Backfill-Stand danach — begrenzt auf dieses
  Fenster, und **idempotent**: wendet ein Consumer das Log ab dem Log-Anfang in
  Lese-Ordnung an (`INSERT`/`UPDATE` als Upsert des Row Images, `DELETE` als
  Löschen), entspricht sein Stand je Schlüssel dem Quellstand.
- **Sichtbarkeits-Grenze.** Positionen wandern nur vorwärts
  ([`LH-FA-CON-004`](lastenheft.md)); der Bestand liegt auf `X`. Er ist ein
  Zustandsabzug für Consumer **vor** `X` (neu registrierte, zurückliegende und
  ein eigens für den Abzug registrierter Consumer). Ein Consumer, dessen
  bestätigte Position beim Commit des Runs bereits `X` erreicht hat, sieht den
  Bestand nicht über seinen Fortschritt; der Bestand bleibt über das
  Bereichslesen ([`LH-FA-REA-001`](lastenheft.md)) lesbar. Ein Bestandsabzug
  teilt **eine** Commit-Position; ein `limit` kann innerhalb einer Position
  nicht fortsetzen (`SPEC-022`), er wird deshalb ohne `limit` oder über den
  Schlüsselvergleich (`commit_position`, `transaction_id`, `sequence`) auf der
  View `cdc.changes` gelesen.
- **Unterbrechung und Neubeginn.** Ein abgebrochener Run (Prozessende,
  Verbindungsverlust) hinterlässt keine Change-Zeile. Beim Prozessstart wird
  jeder `running`-Run `interrupted`; `queued`-Runs laufen weiter; ein
  `interrupted`-Run wird **nicht** automatisch neu gestartet. Ein erneuter
  Antrag startet einen neuen Run mit neuem Snapshot und beginnt damit neu.
- **Fail-closed vor dem Commit.** Unmittelbar vor dem Commit prüft der Run,
  dass die Bindung der Tabelle besteht und der Ausschlussstand
  ([`LH-FA-CFG-005`](lastenheft.md)), der Regelstand der Transformationen
  ([`LH-FA-CFG-007`](lastenheft.md)) sowie der Routing-Regelstand
  ([`LH-FA-CFG-008`](lastenheft.md)) dem Stand entsprechen, mit dem die Blöcke
  gebaut wurden; jede Abweichung rollt den Run zurück (`failed`, Fehlerklasse
  `configuration`, Grund im Fehlertext). Der Routing-Regelstand wird einmal
  nach dem Öffnen des Snapshots gelesen und auf Anwendbarkeit geprüft (Klasse
  `schema`, siehe unten); jeder weitere Block und der Zustand unmittelbar vor
  dem Commit lesen ihn neu und vergleichen ihn als Menge mit diesem Stand. Ein
  nicht lesbarer Stand beendet den Run (`failed`, kein Commit) mit der Klasse
  der Ursache; die Klasse `configuration` gilt dem Wechsel des Standes. Eine Regel, die zwischen zwei
  Lesungen gesetzt und wieder entfernt wird, bleibt unsichtbar, weil der Stand
  keinen Verlauf trägt. Ein ausgeschlossener Wert wird nie serialisiert.
- **Sichtbarkeit und Fehler des Runs.** Die View `cdc.backfill_status` und die
  CLI-Diagnose zeigen den Run (`SPEC-029`). Run-Fehler tragen die Klasse aus
  `SPEC-008` im Fehlertext und sind **run-lokal**: sie setzen weder den
  Heartbeat-Fehlerzustand noch stoppen sie den Capture-Pfad. Ist eine Regel des
  Regelstands der Tabelle — eine Transformationsregel oder eine Routing-Regel
  — auf die Spalten des Snapshots nicht anwendbar (`SPEC-008`, Zeile `schema`),
  endet der Run `failed` mit der Klasse `schema`, bevor die erste Zeile gelesen
  wird — ohne Change und ohne den Erfassungspfad zu berühren.
- **Ziel der Backfill-Changes.** Backfill-Changes durchlaufen dieselbe
  Bestimmung des Zustellziels wie WAL-Changes (`LH-FA-CFG-008.a`) und tragen
  `route_target` des Regelstands zum Run (`SPEC-002`); „Regelstand zum Run" ist
  der Routing-Regelstand, der nach dem Öffnen des Snapshots gelesen wurde und
  für alle Blöcke des Runs gilt. Die Replay-Invariante bleibt unberührt, weil
  das Label nicht zum Zeilenzustand gehört.
- **Zustellung und Retention.** Backfill-Changes gehen in keinen Live-Weg
  (`SPEC-020`, `SPEC-021`, `SPEC-024`); nach dem Commit sendet der Run je
  Tabelle ein Wecksignal (`SPEC-017`, best effort). Die Retention behandelt sie
  wie jeden Change; `committed_at` der synthetischen Transaktionen ist der
  Snapshot-Zeitpunkt.
- **Große Tabellen.** Die Ein-Transaktions-Form hält für die Kopierdauer einen
  Snapshot an der Quelle und eine offene Schreibtransaktion im CDC-Speicher.
  Dafür gibt es zwei Warnungen — **keine Ablehnung, kein Abbruch, keine
  Statusänderung**: (1) beim Antrag, wenn die **geschätzte** Zeilenzahl über
  einer Richtgröße liegt; (2) zur Laufzeit, wenn die Kopierdauer eine Toleranz
  überschreitet. Das Ergebnis steht in den beiden Warn-Spalten (`SPEC-029`).
  Toleranz und Richtgröße sind keine Konstanten dieses Dokuments; `SPEC-029`
  trägt nur ihr Ergebnis.

### LH-FA-CFG-007.a — Transformationsform

**Eingabe:** Rohform einer erfassten Change. **Ausgabe:** transformierte
Change.

Die Lastenheft-Fähigkeit ([`LH-FA-CFG-007`](lastenheft.md)) wird durch einen
geschlossenen Satz deklarativer Transformationsregeln erfüllt — keine
Transformationssprache, kein Skripting-, kein Plugin-Modell. Die folgenden
Sätze sind Zusagen an die Umsetzung, keine Messergebnisse.

- **Konfigurationsmechanismus.** Zwei Antragsarten der Antrags-Queue
  (`SPEC-019`): `set_transformation` legt eine Regel für eine Tabelle an,
  `remove_transformation` nimmt sie heraus; die SQL-Funktionen
  `cdc.set_transformation(...)` und `cdc.remove_transformation(...)`
  schreiben ausschließlich den Antrag. Die Regeln gelten je Tabelle, wirken am
  laufenden Prozess ohne Neustart und überleben ihn: der Regelstand einer
  Tabelle wird aus den `applied`-Zeilen der beiden Antragsarten abgeleitet
  (`SPEC-019`). Es gibt keinen zweiten Konfigurationsweg — weder die
  Konfigurationsdatei noch eine Umgebungsvariable trägt Regeln.
- **Ausdrucksform.** Zwei Regeltypen, `rename_column` und `map_value`, jeder
  mit fester, vollständig prüfbarer Semantik (`SPEC-030`). Der Satz ist
  geschlossen: ein weiterer Regeltyp ist eine Änderung dieser Festlegung.
- **Auswertungsreihenfolge.** Erst der Spaltenausschluss, dann die
  Spaltenregeln; die Reihenfolge der Spaltenregeln und den beobachtbaren
  Inhalt des Ergebnisses führt `SPEC-030` (Reihenfolge und Inhalt).
- **Mehrdeutigkeit.** Sie wird statisch ausgeschlossen, nicht zur Laufzeit
  aufgelöst: vier Konfliktfreiheits-Invarianten je Tabelle (`SPEC-019`)
  verhindern, dass zwei Regeln denselben Schlüssel oder denselben Wert
  beanspruchen. Ein Antrag, der eine Invariante verletzt, endet `failed` mit
  Fehlertext, der Regelstand bleibt unverändert.
- **Verhältnis zum Spaltenausschluss** ([`LH-FA-CFG-005`](lastenheft.md)).
  Der Ausschluss gilt zuerst: eine ausgeschlossene Spalte wird von keiner
  Regel gelesen, ihr Schlüssel erscheint weder unter dem Quell- noch unter
  einem Zielnamen, und ihr Wert nirgends ([`LH-QA-SEC-004`](lastenheft.md)).
  `cdc.exclude_column` und `cdc.include_column` bleiben gegen eine Spalte mit
  Regel zulässig; der Ausschlussstand entscheidet immer vor dem Regelstand.
- **Wirkort.** Die Regeln werden bei der Konstruktion des Row Images
  ausgewertet — vor jeder Serialisierung und vor der Persistierung, an
  **einer** Stelle, die jeder Pfad nutzt, der Changes erzeugt (Replication
  Stream und Backfill). Speicher, SQL-Lesezugriff `cdc.changes`, `GET /changes`,
  gRPC-, SSE- und NATS-Vollinhalts-Stream tragen dadurch denselben Inhalt der
  Row Images (Schlüsselmenge und Werte), live wie beim erneuten Lesen
  ([`LH-FA-REA-005`](lastenheft.md)). Die Rohform wird nicht
  gespeichert: eine Regeländerung wirkt nur auf danach erfasste Changes,
  bereits gespeicherte behalten ihre Form, und die Rohform einer Change ist
  nicht rekonstruierbar. `map_value` ist nicht umkehrbar, sobald mehrere
  Quellwerte auf denselben Zielwert abgebildet werden. Regeln wirken
  ausschließlich auf Schlüssel und Werte der Row Images: `change_id`,
  `transaction_id`, `source_table_id`, `sequence`, `operation`,
  `schema_version` und die Tabellen-Identität (Schema und Tabelle, damit
  Subjekt und Tabellenfilter der Lesewege) bleiben Quell-Identität; die
  Nachrichtenschemata (`SPEC-020`, `SPEC-021`, `SPEC-022`, `SPEC-024`) bleiben
  unverändert, nur die Schlüsselmenge der Images folgt dem Regelstand.
- **Nicht anwendbare Regel.** Wann eine Regel auf eine Change nicht anwendbar
  ist, definiert `SPEC-030` (Anwendbarkeit). Der Erfassungspfad endet dann
  sichtbar mit der Fehlerklasse `schema` (`SPEC-008`): die Transaktion wird
  weder persistiert noch bestätigt, es geht keine Change verloren, die
  Erfassung setzt nach der Abhilfe ab der bestätigten Position fort. Im Run
  eines Backfills gilt dieselbe Ursache mit derselben Klasse, run-lokal
  (`LH-FA-CAP-009.a`, `SPEC-008`).
- **Abhilfe (Zusage).** Diese Stelle führt die Abhilfe im gescheiterten
  Prozess: `cdc.remove_transformation` beantragen — oder die Regel durch eine
  passende ersetzen, erst entfernen, dann neu setzen (`SPEC-019` K1); die
  Anträge bleiben `pending`, solange der Prozess steht —, den Prozess neu
  starten, offene Anträge werden beim Start innerhalb der Frist des Vorlaufs
  verarbeitet, **bevor** die erste Transaktion der Tabelle assembliert wird,
  und die zuvor nicht bestätigte Transaktion erscheint danach über
  `cdc.changes`, ohne dass der Prozess erneut an derselben Regel endet. Diese
  Abfolge muss die Umsetzung liefern; sie ist erst mit dem Beleg am laufenden
  System eine Tatsache.
- **Abgrenzung.** Die Regeln bestimmen die Form einer Change, nicht ihr
  Zustellziel; das Routing bleibt bei `LH-FA-CFG-008.a`.

### LH-FA-CFG-008.a — Routingform

**Eingabe:** erfasste Change. **Ausgabe:** Change mit dem durch eine
Routing-Regel bestimmten Zustellziel, je Ziel abrufbar über jeden Zustellweg.

Die Lastenheft-Fähigkeit ([`LH-FA-CFG-008`](lastenheft.md)) wird durch eine
geordnete Liste einfacher Regeln je Tabelle erfüllt — keine Routing-Sprache,
kein Skripting-, kein Plugin-Modell. Die folgenden Sätze sind Zusagen an die
Umsetzung, keine Messergebnisse: bis die Umsetzung sie am laufenden System
belegt, sind sie keine geprüften Tatsachen.

- **Zielmodell.** Ein Zustellziel ist ein benannter Kanal je Quelle. Sein Name
  folgt `[a-z0-9][a-z0-9_-]{0,62}`; das Alphabet enthält weder Punkt noch
  NATS-Platzhalter, damit der Name ein einzelnes Subjekt-Token sein kann
  (hergeleitet aus der Subjekt-Syntax, nicht gegen einen NATS-Server geprüft).
  „Zustellung" heißt Abruf oder Abonnement des Kanals über die vorhandenen
  Zustellwege; es gibt keine aktive Zustellung an externe Senken und keine feste
  Bindung eines Consumers an ein Ziel, Consumer-Positionen bleiben
  Log-Positionen. Eine Change trägt höchstens ein Ziel, und es gibt kein
  benanntes Standardziel.
- **Konfigurationsmechanismus.** Zwei Antragsarten der Antrags-Queue
  (`SPEC-019`): `set_route` legt eine Regel für eine Tabelle an, `remove_route`
  nimmt sie heraus; die SQL-Funktionen `cdc.set_route(...)` und
  `cdc.remove_route(...)` schreiben ausschließlich den Antrag. Die Regeln gelten
  je Tabelle, wirken am laufenden Prozess ohne Neustart und überleben ihn: der
  Regelstand wird aus den `applied`-Zeilen der beiden Antragsarten abgeleitet
  (`SPEC-019`). Es gibt keinen zweiten Konfigurationsweg.
- **Ausdrucksform.** Eine Regel trägt ein Ziel (`target`), eine Ordnungszahl
  (`order`) und optional eine Bedingung (`when`), die die Gleichheit des
  Textwerts einer Spalte der Quelle prüft (`SPEC-032`). Operatoren über die
  Gleichheit hinaus, ein Standardziel und mehrere Ziele je Change gehören nicht
  zu dieser Festlegung; jede Erweiterung ist eine Änderung dieser Festlegung.
- **Auflösung bei Mehrdeutigkeit.** Der erste Treffer in aufsteigender `order`
  bestimmt das Ziel. Fälle, in denen eine Regel nie wirken könnte, werden
  statisch ausgeschlossen: sechs Invarianten R1 bis R6 je Tabelle (`SPEC-019`).
  Ein Antrag, der eine Invariante verletzt, endet `failed` mit Fehlertext, der
  Regelstand bleibt unverändert.
- **Kein Treffer, abwesender Wert.** Trifft keine Regel, ist das Ziel leer
  (`route_target` ist `NULL`, `SPEC-002`): die Change geht nicht verloren und
  nicht an ein Standardziel, sie bleibt im Log, in der SQL-Sicht mit
  `route_target IS NULL` und über jeden ungefilterten Weg sichtbar; für einen
  Leser, der nur ein Ziel abruft, ist sie unsichtbar. Ein abwesender Wert im
  Bild (`NULL`, unverändertes TOAST, bei `DELETE` ohne volle Replica-Identität
  jede Nicht-Schlüsselspalte) ist ein Nicht-Treffer der Bedingung, keine
  Nichtanwendbarkeit; die nächste Regel wird geprüft. Gemessen (PostgreSQL 17 und
  18, E2E-Messung): ein `DELETE` ohne volle Replica-Identität ist für eine
  Bedingung auf eine Nicht-Schlüsselspalte ein Nicht-Treffer (das Alt-Bild trägt
  nur den Schlüssel; ohne weitere Regel ist `route_target` `NULL`, mit einer
  Regel ohne `when` deren Ziel), unter voller Replica-Identität ein Treffer.
- **Verhältnis zu Transformationen** ([`LH-FA-CFG-007`](lastenheft.md)). Die
  Bedingung liest den Quellwert vor jeder Transformation; das Ziel ist ein
  Metadatum der Change und ändert das Row Image nicht, eine Transformation
  ändert das Ziel nicht. Beide Regelwerke haben getrennte Namensräume und
  getrennte Stände.
- **Verhältnis zum Spaltenausschluss** ([`LH-FA-CFG-005`](lastenheft.md),
  [`LH-QA-SEC-004`](lastenheft.md)). Eine ausgeschlossene Spalte darf keine
  Routing-Bedingung tragen, und eine Spalte mit Routing-Bedingung lässt sich
  nicht ausschließen (R3, `SPEC-019`): sonst verriete das Ziel eine Eigenschaft
  des ausgeschlossenen Werts.
- **Wirkort.** Das Ziel wird bei der Konstruktion der Change bestimmt — vor der
  Persistierung, an **einer** Stelle, die jeder Pfad nutzt, der Changes erzeugt
  (Replication Stream und Backfill) — und als Label `route_target` an der Change
  persistiert: ein Log, kein getrenntes Log je Ziel. Live und beim erneuten
  Lesen ([`LH-FA-REA-005`](lastenheft.md)) trägt eine Change dasselbe Ziel. Die
  Bestimmung verwirft nie eine Change. Das Label ist zum Erfassungszeitpunkt
  festgeschrieben: eine Regeländerung wirkt nur auf danach erfasste Changes, es
  gibt kein Umetikettieren. Backfill-Changes tragen das Label des Regelstands
  zum Run (`LH-FA-CAP-009.a`); der Altbestand erhält ein Ziel durch einen
  Backfill-Run mit dem aktuellen Regelstand.
- **Auswahl, keine Ausblendung.** Ein ungefilterter Leser sieht weiterhin alle
  Changes, geroutete eingeschlossen. Ein Leser wählt ein Ziel über den optionalen
  Parameter `target` (leer = kein Filter, Konjunktion mit `schema`/`table`) am
  Bereichslesen (`SPEC-022`, `SPEC-031` `ReadChanges`) und an den beiden
  Live-Streams (`SPEC-020`, `SPEC-021`), über die Spalte `route_target` der
  SQL-Sicht (`SPEC-002`) und über das Zusatz-Subjekt des NATS-Vollinhalts-Streams
  (`SPEC-024`). Das Wecksignal (`SPEC-017`) trägt kein Ziel.
- **Nicht anwendbare Regel.** Wann eine Regel auf eine Change nicht anwendbar
  ist, definiert `SPEC-032` (Anwendbarkeit). Der Erfassungspfad endet dann
  sichtbar mit der Fehlerklasse `schema` (`SPEC-008`): die Transaktion wird
  weder persistiert noch bestätigt, Heartbeat und Diagnose zeigen den Zustand,
  es geht keine Change verloren. Im Run eines Backfills gilt dieselbe Ursache mit
  derselben Klasse, run-lokal (`LH-FA-CAP-009.a`, `SPEC-008`).
- **Abhilfe (Zusage).** Diese Stelle führt die Abhilfe im gescheiterten
  Prozess: `cdc.remove_route` beantragen — die Anträge bleiben `pending`,
  solange der Prozess steht —, den Prozess neu starten, offene Anträge werden
  beim Start innerhalb der Frist des Vorlaufs verarbeitet, **bevor** die erste
  Transaktion der Tabelle assembliert wird, und die zuvor nicht bestätigte
  Transaktion erscheint danach über `cdc.changes`. Diese Abfolge muss die
  Umsetzung liefern. Die Zusage gilt für den Fall, in dem der Erfassungspfad an
  der Nichtanwendbarkeit der Regel endet. Eine entfernte Spalte mit bekannter
  Spaltenform endet dagegen zuerst im Pfad der inkompatiblen Schemaänderung
  (`LH-FA-SCH-003`, `LH-FA-SCH-004.a`), weil die Prüfung der Relation der Prüfung
  der Change vorausgeht und keine Regel liest (hergeleitet, nicht am laufenden
  System gemessen); dort genügt das Entfernen der Regel nicht, die Abhilfe ist die
  der inkompatiblen Schemaänderung, und die Regel auf die entfernte Spalte ist
  zusätzlich zu entfernen (hergeleitet). Gemessen (PostgreSQL 17 und 18,
  E2E-Messung): die Nichtanwendbarkeit der Regel entsteht am laufenden System bei
  Erstaktivierung ohne bekannte Spaltenform und bei einer Publication mit
  Spaltenliste, die die Bedingungsspalte nicht trägt, erzeugt auf einer Tabelle
  ohne bekannte Spaltenform; der Erfassungspfad endet dort mit der Fehlerklasse
  `schema`, und die Abfolge (Entfernen der Regel beantragen, Neustart) nimmt die
  zuvor nicht bestätigte Transaktion ohne Ziel auf. Hergeleitet bleibt der Fall
  einer Publication mit Spaltenliste bei bereits bekannter Spaltenform (dort ist
  der erste Halt die inkompatible Schemaänderung).
- **Abgrenzung.** Die Routing-Regeln bestimmen das Ziel einer Change, nicht ihre
  Form; die Form bleibt bei `LH-FA-CFG-007.a`.

### LH-FA-SST-009.a — Sprachmatrix und Vertriebsweg offen

**Eingabe:** bestehende Zustellwege (`SPEC-018`, `SPEC-020`, `SPEC-021`,
`SPEC-022`, `SPEC-024`). **Ausgabe:** offiziell gepflegtes, versioniertes Package je
bedienter Sprache.

Die Lastenheft-Fähigkeit ist gefordert ([`LH-FA-SST-009`](lastenheft.md)),
welche Sprache(n) zuerst bedient werden und über welchen
Paket-Vertriebsweg (z. B. NuGet, npm, Maven, Go-Modul-Registry, PyPI), ist
offen, ADR-pflichtig. Die bestehenden Beispiel-Client-Werkzeugketten
(`SPEC-023`) sind kein Vorgriff auf diese Anforderung — sie bleiben
unversioniertes Vorbild ohne Paketveröffentlichung.

Für C#/NuGet ist die Frage beantwortet: `PgChangeFeed.Client` (`SPEC-026`)
deckt HTTP-API, gRPC-Stream, SSE und NATS-Vollinhalt, real Docker-only
paketierbar (`make sdk-pack-csharp`) und real geprüft
(`PgChangeFeed.Client.0.2.1.nupkg`). Eine zweite Sprache oder ein zweiter
Vertriebsweg bleibt offen — diese Kennung bleibt ihre Adresse.

Für Python/PyPI ist die Frage ebenfalls beantwortet: `pgchangefeed`
(`SPEC-027`) deckt HTTP-API, gRPC-Stream, SSE und NATS-Vollinhalt, real
Docker-only paketierbar (`make sdk-pack-python`) und real geprüft
(`pgchangefeed-0.2.1-py3-none-any.whl`, `pgchangefeed-0.2.1.tar.gz`). Eine
dritte Sprache oder ein dritter Vertriebsweg bleibt offen — diese Kennung
bleibt ihre Adresse.

Für Kotlin ist die Frage ebenfalls beantwortet: `pgchangefeed-kotlin`
(`SPEC-028`) deckt HTTP-API, gRPC-Stream, SSE und NATS-Vollinhalt, real
Docker-only paketierbar (`make sdk-pack-kotlin`) und real geprüft
(`pgchangefeed-kotlin-0.2.2.jar`). Vertriebsweg sind zwei Ziele: Cloudsmith
(anonym lesbar, ohne Konto und Token beziehbar) und GitHub Packages (der
Bezug verlangt einen Token); beide Ziele erhalten dieselben Artefakte. Eine
vierte Sprache oder ein weiterer Vertriebsweg bleibt offen — diese Kennung
bleibt ihre Adresse.

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
| `cdc.change` | die einzelnen Changes (SPEC-002), mit dem Zustellziel `route_target` ([`LH-FA-CFG-008`](lastenheft.md)) |
| `cdc.consumer` | registrierte Consumer |
| `cdc.consumer_position` | bestätigte Position je Consumer |
| `cdc.schema_version` | Schema-Versionen (SPEC-004) |
| `cdc.process_heartbeat` | Betriebs-/Capture-Zustand: periodisches Lebenszeichen des Capture-Prozesses ([`LH-FA-ADM-002`](lastenheft.md), [`LH-QA-OPS-002`](lastenheft.md)) und der zuletzt beobachtete Fehlerzustand als Klasse (`error_class`, `text`, nullable) und Meldungscode (`error_code`, `text`, nullable, ein Code der Klasse, `SPEC-008`); beide `NULL` im Normalbetrieb ([`LH-FA-ADM-003`](lastenheft.md)). Die View `cdc.heartbeat` projiziert beide Spalten, `error_code` als letzte Spalte hinter `age_seconds`, und zeigt den Code nur neben einer Klasse |
| `cdc.backfill_run` | Run-Zustand eines Backfills des Tabellenbestands (`SPEC-029`, [`LH-FA-CAP-009`](lastenheft.md)) |

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
  "schema_version": "Referenz auf cdc.schema_version",
  "origin": "wal | backfill (fehlender Wert liest als wal; LH-FA-CAP-009)",
  "route_target": "Name des Zustellziels | null (null = kein Ziel; LH-FA-CFG-008)"
}
```

Im MVP werden `old_data` und `new_data` als `jsonb` gespeichert (Row
Images; Verfügbarkeit je Operationstyp siehe [`LH-FA-CAP-008`](lastenheft.md)).

Die Schlüsselmenge eines Row Images folgt dem Regelstand der Tabelle zum
Erfassungszeitpunkt (`SPEC-030`): ein umbenannter Schlüssel steht unter seinem
Zielnamen, nie zusätzlich unter dem Quellnamen; die Werte sind Zeichenketten
im Text-Stand der Quelle, ein abgebildeter Wert steht als der zugeordnete
Wert. Alle übrigen Felder und die Identität der Tabelle bleiben
Quell-Identität; `schema_version` referenziert weiterhin die Struktur der
Quelltabelle ([`LH-FA-SCH-005`](lastenheft.md)), nicht den Regelstand.

`origin` benennt die Herkunft des Changes: `wal` für einen über den
Replication Stream erfassten Change, `backfill` für einen Bestands-Change
eines Backfill-Runs (`LH-FA-CAP-009.a`). Die Menge ist geschlossen; ein
gespeicherter Change ohne das Feld liest sich als `wal`. Die Live-Wege
(`SPEC-020`, `SPEC-021`, `SPEC-024`) tragen es nicht.

`route_target` benennt das Zustellziel der Change (`LH-FA-CFG-008.a`,
`SPEC-032`): den Namen eines Ziels oder `NULL` für „kein Ziel". Die Spalte an
`cdc.change` ist nullable, ohne `DEFAULT` und ohne `CHECK`; ein gespeicherter
Change ohne das Feld und eine Change ohne Regeltreffer lesen sich gleich als
`NULL`. Das Label wird bei der Erfassung festgeschrieben (`LH-FA-CFG-008.a`,
Wirkort) und ändert sich danach nicht. Es gehört weder zum Row Image noch zur
Schlüsselmenge der Images, und die Live-Wege (`SPEC-020`, `SPEC-021`,
`SPEC-024`) und `GET /changes` (`SPEC-022`) tragen es nicht als Feld ihrer
Nachrichten; gelesen wird es über die View `cdc.changes`. Die View führt
`route_target` als **letzte** Spalte, `origin` steht davor.

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
Precedence; additiv heißt: eine gesetzte Umgebungsvariable wirkt auch
unter geladener Datei, und ein Feld **ohne** Datei-Gegenstück fällt
dadurch nicht weg. Felder, die Zugangsdaten tragen **können**, sind
**nicht zulässig** und brechen das Laden über die Fehlerklasse
`configuration` mit einer eigenen, den Grund nennenden Fehlerzeile ab
([`LH-QA-SEC-001`](lastenheft.md)/[`LH-QA-SEC-002`](lastenheft.md):
Least-Privilege/Secret-Trennung — eine Konfigurationsdatei ist für andere
Aufbewahrungs-/Verteilwege bestimmt als eine Umgebungsvariable). Die
Klasse umfasst die drei DSN-Schlüssel (`capture_dsn`/`admin_dsn`/
`reader_dsn`), die drei Token-Schlüssel (`api_token_reader`/
`api_token_admin`/`nats_stream_token`), `nats_url`, `otlp_endpoint`,
`otlp_headers` und die beiden Token-Listen (`api_tokens_reader`/
`api_tokens_admin`): die URL-Formen dieser Klasse können Benutzer/Passwort
einbetten, Header und Listen tragen Token, und die Abwesenheit von
Zugangsdaten in einem konkreten Wert ist keine Eigenschaft des Feldes.
`nats_stream_token` trägt den Verbindungs-Token des dritten,
vollinhaltstragenden NATS-Zustellwegs (`SPEC-024`) — derselbe
Zugangsdaten-Charakter wie die beiden API-Token-Schlüssel. Striktes
Decoding (jeder andere unbekannte Schlüssel → Fehlerklasse `configuration`).

| Schlüssel | Typ | Entspricht (Env-Var) | Pflicht in der Datei |
|---|---|---|---|
| `source_id` | string | `CDC_SOURCE_ID` | nein — Pflichtfeld nach Merge (Datei oder Env) |
| `publication` | string | `CDC_PUBLICATION` | nein — Pflichtfeld nach Merge (Datei oder Env) |
| `slot` | string | `CDC_SLOT` | nein — Pflichtfeld nach Merge (Datei oder Env) |
| `tables` | Mapping `<schema.tabelle>: {table_id, schema_version}` | `CDC_TABLES` (Zeichenkettenform, unverändert bestehen) | nein — Pflichtfeld nach Merge (Datei oder Env) |
| `log_level` | string (`debug`/`info`/`warn`/`error`) | `CDC_LOG_LEVEL` | nein, Default `info` |
| `http_addr` | string (`host:port`) | `CDC_HTTP_ADDR` | nein — ungesetzt bleibt die HTTP-/JSON-API deaktiviert |
| `grpc_addr` | string (`host:port`) | `CDC_GRPC_ADDR` | nein — ungesetzt bleibt der gRPC-Stream deaktiviert |
| `tls_cert_file` | string (Dateipfad) | `CDC_TLS_CERT_FILE` | nein — nur zusammen mit `tls_key_file` (`SPEC-034`) |
| `tls_key_file` | string (Dateipfad) | `CDC_TLS_KEY_FILE` | nein — nur zusammen mit `tls_cert_file` (`SPEC-034`) |
| `otlp_interval` | ganze Zahl (Sekunden, 5 bis 3600) | `CDC_OTLP_INTERVAL_SECONDS` | nein, Default 60 (`SPEC-033`) |
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
http_addr: ":8090"
```

`tables` als Feld wird **als Ganzes** ersetzt, nicht Zeile für Zeile
zusammengeführt: eine gesetzte `CDC_TABLES` schlägt die gesamte
Datei-`tables`-Mapping vollständig, ohne Vermischung einzelner Tabellen aus
beiden Quellen. `CDC_CONFIG_FILE` selbst trägt den Dateipfad; leer/unbenannt
bedeutet kein Dateizugriff, der bestehende Env-only-Pfad bleibt unverändert
Default.

Die env-exklusiven Variablen (`CDC_NATS_URL`, `CDC_NATS_STREAM_TOKEN`,
`CDC_API_TOKEN_READER`, `CDC_API_TOKEN_ADMIN`, `CDC_API_TOKENS_READER`,
`CDC_API_TOKENS_ADMIN`, `CDC_OTLP_ENDPOINT`, `CDC_OTLP_HEADERS`) werden auch
unter geladener Datei aus der Umgebung gelesen — sie haben kein Datei-Gegenstück, ihre
Herkunft ist die Umgebungsvariable auf beiden Pfaden.

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
Fähigkeiten des HTTP-Adapters, sowie den zehnten, lesenden Endpunkt
`GET /diagnose`. Das Changes-Lesen ([`LH-FA-REA-*`](lastenheft.md)) ist in
`SPEC-022` ausgestaltet; Health (`Healthcheck`, das binäre Docker-Verdikt)
bleibt außerhalb.

**Authn-Header** (für alle Endpunkte gleich): `Authorization: Bearer
<token>` — fehlend oder einer nicht konfigurierten Klasse entsprechend →
`401`; ein bekanntes `reader`-Token gegen einen administrativen Endpunkt →
`403`; ein bekanntes `admin`-Token erreicht jeden Endpunkt (`admin` deckt
implizit die lesende Klasse ab). Token-Umgebungsvariablen:
`CDC_API_TOKEN_READER` (lesende Rechtsklasse), `CDC_API_TOKEN_ADMIN`
(schreibend/administrativ); mehrere Token je Klasse: `SPEC-035`, TLS: `SPEC-034` — orthogonal zum DB-Rollenmodell der
Verdrahtung: die API-Token-Prüfung entscheidet an der HTTP-Schicht, welcher
Use Case erreichbar ist; welche DSN der Adapter darunter benutzt, bleibt
die bei der Verdrahtung fixierte. Aktivierung optional über
`CDC_HTTP_ADDR`; ungesetzt bedeutet deaktiviertes Feature, kein
HTTP-Server, unverändertes Bestandsverhalten.

**Fehler-Antwortform** (für alle Endpunkte gleich):
`{"error": "<Klartext>", "code": "<Meldungscode>"}` bei
`400`/`401`/`403`/`404`/`500`/`503` — ungültiger JSON-Body oder eine verletzte
Domänen-Invariante (`400`), fehlender/unbekannter Bearer-Token (`401`),
bekanntes Token mit unzureichender Rechtsklasse (`403`), eine an der Quelle
physisch fehlende Tabelle (`404`, nur `EnableTable`/`DisableTable`/`GetStatus`),
unerwarteter interner Fehler (`500`), ein Stream ohne verdrahteten Broadcaster
(`503`, nur `GET /changes/stream`). `code` ist additiv und der Meldungscode der
Ursache (`SPEC-008`): eine Ablehnung einer Aufrufer-Eingabe (`400`, `404`) trägt
einen Code des Bereichs `E8` (`PCF-E8050` Body kein gültiges JSON, `PCF-E8051`
Pflichtfeld oder Pflichtparameter fehlt, `PCF-E8052` unbekannter
Query-Parameter, `PCF-E8053` nicht lesbarer oder außerhalb des Bereichs liegender
Wert, `PCF-E8054` Position unter 1, `PCF-E8055` Position rückt nicht vor,
`PCF-E8056` Endposition vor Startposition, `PCF-E8057` Position einer anderen
Quelle, `PCF-E8025` Tabelle fehlt an der Quelle); ein `500` trägt den Code der
klassifizierten Ursache in der Fehlerkette, sonst den Rückfall `PCF-E7000`; das
`503` trägt `PCF-E2001`. `401` und `403` tragen kein Feld `code`: die Tabelle der
Codes führt für sie keinen Eintrag. Ein Leser, der nur `error` liest, bleibt
unverändert lauffähig; die Statuscodes ändern sich nicht.

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
| `Diagnose` ([`LH-FA-SST-003`](lastenheft.md), deckt [`LH-FA-ADM-002`](lastenheft.md)…[`005`](lastenheft.md), [`LH-FA-RET-005`](lastenheft.md), [`LH-FA-RET-006`](lastenheft.md), [`LH-FA-CAP-009`](lastenheft.md)) | `GET /diagnose?source=<string>` | `reader` oder `admin` | Query-Parameter `source` Pflicht, **kein** weiterer Parameter zulässig (`400` sonst, strenger als die neun Bestandsendpunkte oben, dieselbe Begründung wie `GET /changes` in `SPEC-022`) | `200`: `{"heartbeat_age_seconds": <float64 \| null>, "error_class": "<string> \| null", "error_code": "<string> \| null", "capture_lag": <float64>, "consumer_lags": [{"consumer_id": "<string>", "lag": <float64 \| null>}, …], "retention_blocker": {"consumer_id": "<string>", "name": "<string>", "acknowledged_position": <int64>, "backlog": <int64 \| null>} \| null, "storage_bytes": <float64>, "backfill": [{"schema": "<string>", "table": "<string>", "status": "<string>", "rows_copied": <int64>, "estimated_rows": <int64 \| null>, "warn_estimated_size": <bool>, "warn_duration": <bool>, "error_message": "<string>"}, …]}` — `heartbeat_age_seconds == null` bedeutet „kein Lebenszeichen"; in diesem Fall tragen `error_class` und `error_code` ebenfalls `null`; `error_code` ist der Meldungscode des Fehlerzustands (`SPEC-008`) und wie `error_class` im Normalbetrieb `null`; `retention_blocker == null` bedeutet „kein Blocker"; `consumer_lags`/`backfill` sind leere, gesetzte Listen ohne Treffer, nie `null`; ein Lesefehler an einer der Diagnose-Views → `500` |

### SPEC-019 — `cdc.administration_request` (Antrags-Datensatz)

Feldform des Antrags-Datensatzes der schreibenden SQL-Administration
([`LH-FA-ADM-001`](lastenheft.md), [`LH-FA-CFG-005`](lastenheft.md),
[`LH-FA-CFG-007`](lastenheft.md), [`LH-FA-CFG-008`](lastenheft.md),
[`LH-FA-CAP-009`](lastenheft.md)):
`cdc.enable_table`/`cdc.disable_table`/`cdc.exclude_column`/
`cdc.include_column`/`cdc.backfill_table`/`cdc.set_transformation`/
`cdc.remove_transformation`/`cdc.set_route`/`cdc.remove_route` schreiben
ausschließlich eine Zeile hierher und
senden `pg_notify` auf dem Kanal `cdc_administration`; der laufende
Capture-Prozess liest die offenen Anträge und vermerkt das Ergebnis in
derselben Zeile.

| Spalte | Typ | Pflicht | Bedeutung |
|---|---|---|---|
| `administration_request_id` | text (PK) | ja | von der SQL-Funktion vergeben; zugleich der `pg_notify`-Payload |
| `source_id` | text (FK `cdc.source`) | ja | Quelle des Antrags |
| `schema_name` / `table_name` | text | ja | adressierte Tabelle; ein leerer Wert endet den Antrag `failed` (Tabelle „Zeilen, die kein Antrag sind“ unten) |
| `column_name` | text | nein | Ziel-Spalte der beiden Spalten-Antragsarten, dort nichtleer; die sieben übrigen Antragsarten (`enable`, `disable`, `backfill`, `set_transformation`, `remove_transformation`, `set_route`, `remove_route`) tragen hier NULL |
| `rule_name` | text | nein | Regelname der Transformations- und der Routing-Antragsarten, je Tabelle und Regelwerk eindeutig, Alphabet in `SPEC-030` (Bezeichner); die fünf übrigen Antragsarten tragen hier NULL |
| `rule_spec` | jsonb | nein | Regelform der Antragsarten `set_transformation` (`SPEC-030`) und `set_route` (`SPEC-032`); die sieben übrigen Antragsarten tragen hier NULL |
| `request_kind` | text | ja | geschlossene Menge `enable` \| `disable` \| `exclude_column` \| `include_column` \| `backfill` \| `set_transformation` \| `remove_transformation` \| `set_route` \| `remove_route` |
| `requested_at` | timestamptz | ja, Default `current_timestamp` | Zeitpunkt des Funktionsaufrufs (die Funktionen setzen ihn ausdrücklich auf `clock_timestamp()`); die Verarbeitungs-Ordnung |
| `status` | text | ja, Default `pending` | geschlossene Menge `pending` \| `applied` \| `failed` |
| `error_message` | text | nein | Fehlertext eines `failed`-Antrags: `abgelehnt [<Code>]: <Klartext>` bei der Ablehnung einer Aufrufer-Eingabe, der Kopf `Fehlerklasse <klasse> [<Code>]: <Ursache>` bei einem klassifizierten Fehler (`SPEC-008`); `applied` trägt NULL |

**Ordnung der Verarbeitung.** Die neun SQL-Funktionen setzen `requested_at` auf den
Zeitpunkt ihres Aufrufs (`clock_timestamp()`), nicht auf den Beginn der Transaktion.
Aufrufe **einer** Transaktion tragen dadurch verschiedene Zeitstempel in der Reihenfolge
des Aufrufs: `remove_transformation` vor `set_transformation` derselben Regel,
`remove_route` vor `set_route` derselben Regel,
`exclude_column` vor `include_column` derselben Spalte und `disable_table` vor
`enable_table` werden in dieser Folge verarbeitet und abgeleitet. Die offenen Anträge
werden in der Ordnung `requested_at`, bei gleichem Zeitstempel nach
`administration_request_id` verarbeitet — dieselbe Ordnung, in der die dauerhaften Stände
abgeleitet werden; die Kennung ordnet nur deterministisch, nicht zeitlich. Grenzen: Die
Ordnung ist der Zeitpunkt des Aufrufs, nicht der des `COMMIT`. Schreibt die Transaktion
eines früher aufgerufenen Antrags später fest als die eines später aufgerufenen Antrags
auf dieselbe Regel oder Spalte, verarbeitet die Queue den früheren nach dem späteren, die
Ableitung ordnet ihn davor; der laufende Stand weicht dann bis zum nächsten Prozessstart
vom abgeleiteten Stand ab. Anträge auf dieselbe Regel oder Spalte werden deshalb nicht aus
überlappenden Transaktionen abgesetzt. Ein Rückwärtssprung der Serveruhr zwischen zwei
Aufrufen kehrt deren Ordnung um.

`column_name` ist für die beiden Spalten-Antragsarten Pflicht (ein Antrag
ohne Spalte adressiert kein Ziel und endet `failed`, Tabelle „Zeilen, die
kein Antrag sind“ unten); ein Spaltenausschluss gegen eine an der Quelle nicht
existierende Spalte endet als `failed` mit dem Fehlertext der
`ErrSourceColumnMissing`-Ausprägung (Klartext „Spalte existiert nicht an der
Quelle", gefolgt von der Adresse `schema.table.column`, getrennt durch
Punkte). Die Antragsarten `enable`/`disable` tragen unverändert
Bindungs-Zeilen und Publication nach.

**Zeilen, die kein Antrag sind.** Die Lesung der Queue lehnt keine Zeile ab.
Die SQL-Funktionen prüfen Quelle, Schema, Tabelle und Spalte nicht; eine
`pending`-Zeile, die der Antrags-Konstruktor verwirft, endet in der Verarbeitung
`failed`, und die Zeilen dahinter werden in der Ordnung der Queue
weiterverarbeitet. `error_message` beginnt mit `abgelehnt [<Code>]: `
(`SPEC-008`), es folgt der Klartext der Tabelle, ein
Doppelpunkt, ein Leerzeichen und die Adresse — hier die
Antrags-Kennung (`administration_request_id`), weil Schema, Tabelle oder Spalte
selbst leer sind. Die erste verletzte Prüfung bestimmt den Text, in der
Reihenfolge der Tabelle von oben nach unten:

| Verletzung | Klartext | Adresse |
|---|---|---|
| `source_id` ist leer | `Quelle ist leer` | Antrags-Kennung |
| `schema_name` ist leer | `Schemaname ist leer` | Antrags-Kennung |
| `table_name` ist leer | `Tabellenname ist leer` | Antrags-Kennung |
| `request_kind` liegt außerhalb der geschlossenen Menge | `Antragsart ist unbekannt` | Antrags-Kennung |
| `column_name` ist bei `exclude_column`/`include_column` leer oder NULL | `Spaltenname ist leer` | Antrags-Kennung |
| jeder andere Grund, aus dem der Antrags-Konstruktor die Zeile verwirft | `Antrag ist ungültig` | Antrags-Kennung |

Grenze: eine Zeile ohne Kennung (`administration_request_id` leer) lässt sich
nicht vermerken. Die Verarbeitung überspringt sie mit einer Warnung im Log; sie
bleibt `pending` und hält die Zeilen dahinter nicht an. Ein Fehler beim Vermerk
einer verworfenen Zeile wird protokolliert und hält den Durchlauf nicht an; die
Zeile bleibt `pending` und wird im nächsten Durchlauf erneut vermerkt.

Für die Antragsart `backfill` heißt `applied` **angenommen**: die Run-Zeile
(`cdc.backfill_run`, `SPEC-029`) entsteht mit dem Status `queued` in
**derselben Transaktion**, die den Antrag von `pending` auf `applied` setzt —
der Vermerk trifft genau **eine** `pending`-Zeile, jede Abweichung rollt beide
Schreibvorgänge zurück. Die Ausführung des Runs steht im Run-Zustand, nicht im
Antragsstatus. Ein Antrag, dessen Vorbedingung verletzt ist (Tabelle nicht
aktiviert oder ohne laufende Bindung, ein aktiver Run — `queued`/`running` —
derselben Tabelle), endet `failed` mit Fehlertext und hinterlässt keine
Run-Zeile.

Für die beiden Spalten-Antragsarten trägt `applied` darüber hinaus eine
zweite Bedeutung: der Stand ist **dauerhaft vermerkt**. Ihre
`applied`-Zeilen sind die einzige Herkunft des Ausschlussstandes einer
Tabelle — ausgewertet in der Reihenfolge `requested_at`, bei gleichem
Zeitstempel deterministisch nach `administration_request_id`;
`exclude_column` trägt den Spaltennamen ein, `include_column` nimmt ihn
wieder heraus. Jeder Pfad, der eine Erfassungs-Bindung anlegt, trägt den so
abgeleiteten Stand der adressierten Tabelle mit — der Prozessstart **und**
der Aktivierungs-Zweig der Antrags-Verarbeitung; eine Deaktivierung mit
anschließender Aktivierung stellt ihn damit ebenso her wie ein Neustart. Ein
Antrag gegen eine existierende Spalte einer Tabelle ohne laufende Bindung
endet deshalb `applied`, nicht `failed`: er wirkt, sobald die Tabelle erfasst
wird. Die Antrags-Zeilen dieser beiden Arten sind dadurch tragend — eine
Bereinigung der Tabelle verlöre den Stand.

**Transformations-Antragsarten.** `rule_name` ist für `set_transformation`
und `remove_transformation` Pflicht, `rule_spec` nur für
`set_transformation`. Der Aufruf prüft beides nicht, und die Lesung reicht
einen Antrag mit fehlendem oder ungültigem Wert wie jede Zeile durch: die
Prüfung liegt in der Verarbeitung, der Antrag endet `failed` mit dem
Fehlertext der Tabelle unten.
`cdc.set_transformation` nimmt die Regelform als `json`-Parameter an und
schreibt sie als `jsonb` in die Spalte `rule_spec`: ein Literal und ein
`::json`-Wert werden angenommen, ein `jsonb`-typisierter Wert (`::jsonb`, das
Ergebnis von `jsonb_build_object`) braucht den Cast `::json`. Angenommen wird
`NULL` (SQL-`NULL`) und jeder Text, den PostgreSQL als `json` liest und nach
`jsonb` umwandeln kann — auch JSON-`null`, ein Wert ohne Objekt und ein doppelter
Schlüssel (die Spalte hält den letzten Wert). Der Aufruf scheitert ohne
Antrags-Zeile, wenn der Wert für PostgreSQL kein JSON ist (Syntaxfehler, leerer
Text, `NaN`, einzelnes Surrogat-Escape) oder die Umwandlung nach `jsonb`
scheitert: das Zeichen U+0000 als Escape `\u0000` in einem Schlüssel oder Wert,
eine Zahl außerhalb des Zahlbereichs von `numeric`, eine Schachtelung jenseits der
Stapeltiefe des Servers. Die Aussage ist an PostgreSQL 17 und 18 gemessen. Die
Prüfung der Regelform liegt vollständig im Capture-Prozess, nicht in SQL; die
Annahme des Aufrufs ist keine Zusage, dass der Antrag verarbeitet wird. Ein Antrag, der eine der
folgenden Bedingungen verletzt, endet `failed` mit dem Fehlertext der
Tabelle unten und lässt den Regelstand unverändert — die **Konfliktfreiheit**
je Tabelle, die Mehrdeutigkeit statt sie aufzulösen ausschließt:

- **K1** — `rule_name` ist je Tabelle eindeutig; ein `set_transformation`
  gegen einen vergebenen Namen endet `failed` (erst entfernen, dann neu
  setzen; beide Aufrufe dürfen in einer Transaktion stehen, siehe „Ordnung der
  Verarbeitung“).
- **K2** — jede Quellspalte trägt höchstens eine Spaltenregel
  (`rename_column` oder `map_value`); Umbenennung und Wertabbildung derselben
  Spalte sind nicht kombinierbar.
- **K3** — die Zielnamen (`to` von `rename_column`) sind untereinander
  verschieden und verschieden von jedem Spaltennamen der Quelltabelle,
  ausgeschlossene Spalten und die Quellspalte selbst eingeschlossen;
  „verschieden" heißt zeichengenau verschieden (`SPEC-030`, Bezeichner).
- **K4** — die Spalte der Regel (`column`) existiert an der Quelle;
  `remove_transformation` gegen einen nicht geführten Regelnamen endet
  `failed`.

Der Fehlertext beginnt mit `abgelehnt [<Code>]: ` (`SPEC-008`); es folgt der
Klartext der Zeile, ein Doppelpunkt, ein Leerzeichen und die Adresse; die
Adresse eines Regelnamens ist
`schema.table.rule_name`, die einer Spalte `schema.table.column`, die eines
Zielnamens `schema.table.zielname`; der Name steht zeichengenau, wie beantragt.
Die erste verletzte Prüfung bestimmt den Text, in der Reihenfolge der Tabelle
von oben nach unten: erst die fünf Formzeilen (Regelname, Form der `rule_spec`,
Regeltyp, Schlüssel, Pflichtschlüssel und Werte), dann K1, K2, K3, K4.
`remove_transformation` durchläuft nur die Zeile zum Regelnamen und K4.

| Verletzung | Klartext | Adresse |
|---|---|---|
| `rule_name` ist leer oder NULL oder liegt außerhalb des Alphabets (`SPEC-030`, Bezeichner) | `Regelname ist ungültig` | Regelname, wie beantragt |
| `set_transformation` mit `rule_spec` NULL (SQL oder JSON) oder ohne JSON-Objekt, oder `kind` fehlt oder ist keine Zeichenkette | `rule_spec ist ungültig` | Regelname |
| `kind` ist kein Regeltyp aus `SPEC-030` | `unbekannter Regeltyp` | der `kind`-Wert |
| `rule_spec` trägt einen Schlüssel, den der Regeltyp nicht kennt | `unbekannter Schlüssel in rule_spec` | der Schlüsselname |
| ein Pflichtschlüssel des Regeltyps fehlt oder hat den falschen Typ, `column`/`to` verletzt die Bezeichner-Form, `values` ist leer oder trägt einen Wert, der keine Zeichenkette ist (`SPEC-030`) | `rule_spec ist ungültig` | Regelname |
| K1 | `Regelname bereits vergeben` | Regelname |
| K2 | `Spalte trägt bereits eine Regel` | Spalte |
| K3, Zielname gleicht dem Zielnamen einer anderen Regel | `Zielname kollidiert mit einer anderen Regel` | Zielname |
| K3, Zielname gleicht einem Spaltennamen der Quelltabelle | `Zielname kollidiert mit einer Spalte der Tabelle` | Zielname |
| K4, Spalte fehlt an der Quelle | `Spalte existiert nicht an der Quelle` | Spalte |
| K4, `remove_transformation` gegen einen nicht geführten Regelnamen | `Regelname nicht geführt` | Regelname |

Für die beiden Transformations-Antragsarten trägt `applied` — wie bei den
Spalten-Antragsarten — eine zweite Bedeutung: der Stand ist **dauerhaft
vermerkt**. Ihre `applied`-Zeilen sind die einzige Herkunft des Regelstands
einer Tabelle, ausgewertet in der Reihenfolge `requested_at`, bei gleichem
Zeitstempel deterministisch nach `administration_request_id`;
`set_transformation` trägt die Regel ein, `remove_transformation` nimmt sie
heraus. Jeder Pfad, der eine Erfassungs-Bindung anlegt, trägt den abgeleiteten
Regelstand mit — der Prozessstart **und** der Aktivierungs-Zweig der
Antrags-Verarbeitung. Ein Antrag gegen eine Tabelle ohne laufende Bindung endet
`applied`, nicht `failed`, sofern keine der Bedingungen oben verletzt ist: er
wirkt, sobald die Tabelle erfasst wird. Die Antrags-Zeilen dieser beiden Arten
sind dadurch tragend — eine Bereinigung der Tabelle verlöre den Regelstand.
Der aktive Regelstand ist über diese Tabelle lesbar; eine eigene Sicht gibt es
nicht.

**Routing-Antragsarten.** `rule_name` ist für `set_route` und `remove_route`
Pflicht, `rule_spec` nur für `set_route`; der Namensraum der Routing-Regeln ist
getrennt von dem der Transformationsregeln (derselbe Name darf in beiden
vorkommen). Die Prüfung liegt wie bei den Transformations-Antragsarten
vollständig in der Verarbeitung (`cdc.set_route` nimmt die Regelform wie
`cdc.set_transformation` als `json`-Parameter an; Annahmemenge, Cast und
Grenzen wie dort), der Antrag endet bei einer Verletzung `failed` mit dem
Fehlertext der Tabelle unten und lässt den Regelstand unverändert. Die
**Konfliktfreiheit** je Tabelle schließt Mehrdeutigkeit statt sie aufzulösen
aus:

- **R1** — `rule_name` ist je Tabelle unter den Routing-Regeln eindeutig
  (erst entfernen, dann neu setzen).
- **R2** — `order` ist je Tabelle eindeutig.
- **R3** — `when.column` existiert an der Quelle und ist keine
  ausgeschlossene Spalte; umgekehrt endet `exclude_column` gegen eine Spalte,
  die eine Routing-Bedingung des Regelstands trägt, `failed`. `include_column`
  bleibt zulässig.
- **R4** — höchstens eine Regel ohne `when` je Tabelle, und sie trägt die
  höchste `order` der Tabelle.
- **R5** — das Paar (`when.column`, `when.equals`) kommt je Tabelle höchstens
  einmal vor; eine Regel ohne `when` trägt kein Paar.
- **R6** — `remove_route` gegen einen nicht geführten Regelnamen endet `failed`.

Der Fehlertext beginnt mit `abgelehnt [<Code>]: ` (`SPEC-008`); es folgt der
Klartext der Zeile, ein Doppelpunkt, ein Leerzeichen und die Adresse; die
Adresse eines Regelnamens ist
`schema.table.rule_name`, die einer Spalte `schema.table.column`, die eines
Zielnamens `schema.table.zielname`, die einer `order` `schema.table.<Wert>` in
Dezimalschreibweise, die eines Schlüssels der Schlüsselname; ein Name steht
zeichengenau, wie beantragt. Die erste verletzte
Prüfung bestimmt den Text, in der Reihenfolge der Tabelle von oben nach unten:
erst die Formzeilen, dann R1 bis R5; `remove_route` durchläuft nur die Zeile zum
Regelnamen und R6. Die letzte Zeile gilt für `exclude_column` zusätzlich zu
den Prüfungen der Spalten-Antragsarten.

| Verletzung | Klartext | Adresse |
|---|---|---|
| `rule_name` ist leer oder NULL oder liegt außerhalb des Alphabets (`SPEC-030`, Bezeichner) | `Regelname ist ungültig` | Regelname |
| `set_route` mit `rule_spec` NULL (SQL oder JSON) oder ohne JSON-Objekt | `rule_spec ist ungültig` | Regelname |
| `rule_spec` oder `when` trägt einen Schlüssel, den die Regelform nicht kennt (`SPEC-032`) | `unbekannter Schlüssel in rule_spec` | der Schlüsselname |
| `target` fehlt oder ist keine Zeichenkette, `order` fehlt oder ist keine positive Ganzzahl, `when` ist kein Objekt oder trägt nicht beide Schlüssel, `column` oder `equals` hat die falsche Form (`SPEC-032`) | `rule_spec ist ungültig` | Regelname |
| `target` liegt außerhalb des Alphabets des Zielnamens (`SPEC-032`) | `Zielname ist ungültig` | Zielname |
| R1 | `Regelname bereits vergeben` | Regelname |
| R2 | `order bereits vergeben` | order |
| R3, `when.column` fehlt an der Quelle | `Spalte existiert nicht an der Quelle` | Spalte |
| R3, `when.column` ist ausgeschlossen | `Spalte ist ausgeschlossen` | Spalte |
| R4, eine Regel ohne `when` ist bereits geführt | `Regel ohne when bereits vorhanden` | Regelname der geführten Regel |
| R4, die Regel ohne `when` trägt nicht die höchste `order` | `Regel ohne when trägt nicht die höchste order` | Regelname |
| R4, die Regel mit `when` trägt eine höhere `order` als die geführte Regel ohne `when` | `order liegt hinter der Regel ohne when` | Regelname der geführten Regel ohne `when` |
| R5 | `Bedingung bereits vergeben` | Spalte |
| R6, `remove_route` gegen einen nicht geführten Regelnamen | `Regelname nicht geführt` | Regelname |
| R3, `exclude_column` gegen eine Spalte mit Routing-Bedingung | `Spalte trägt eine Routing-Bedingung` | Spalte |

Für die beiden Routing-Antragsarten trägt `applied` wie bei den
Transformations-Antragsarten eine zweite Bedeutung: der Stand ist **dauerhaft
vermerkt**. Ihre `applied`-Zeilen sind die einzige Herkunft des Routing-
Regelstands einer Tabelle, ausgewertet in der Reihenfolge `requested_at`, bei
gleichem Zeitstempel deterministisch nach `administration_request_id`;
`set_route` trägt die Regel ein, `remove_route` nimmt sie heraus. Jeder Pfad, der
eine Erfassungs-Bindung anlegt, trägt den abgeleiteten Stand mit — der
Prozessstart **und** der Aktivierungs-Zweig der Antrags-Verarbeitung. Ein Antrag
gegen eine Tabelle ohne laufende Bindung endet `applied`, nicht `failed`, sofern
keine der Bedingungen oben verletzt ist: er wirkt, sobald die Tabelle erfasst
wird. Die Antrags-Zeilen dieser beiden Arten sind dadurch tragend — eine
Bereinigung der Tabelle verlöre den Regelstand. Der aktive Regelstand ist über
diese Tabelle lesbar; eine eigene Sicht gibt es nicht. Die Wortlaute und
Reihenfolgen dieses Absatzes sind Zusagen an die Umsetzung (`SPEC-032`); die
Gleichheit der Annahmemenge mit `cdc.set_transformation` ist hergeleitet, die
Messung an PostgreSQL 17 und 18 gilt dort.

**Grants** (Rollen nach der Zuordnung der DSN-Verdrahtung):

| Rolle | Recht auf `cdc.administration_request` | Träger |
|---|---|---|
| `cdc_admin` | `SELECT`, `UPDATE` | die Administrations-Verarbeitung: offene Anträge lesen, den Ausgang (`applied`/`failed`) vermerken, die dauerhaften Stände aus den `applied`-Zeilen ableiten; die Annahme eines Backfills vermerkt in derselben Transaktion |
| `cdc_capture` | **keines** | — |
| `cdc_reader` | **keines** | — |

Niemand trägt `INSERT` oder `DELETE`: Anträge legen ausschließlich die
SQL-Funktionen an (`SECURITY DEFINER`, unter den Rechten ihres Eigentümers);
das Recht, eine Funktion aufzurufen, trägt allein `cdc_admin`.

### SPEC-020 — gRPC-Live-Change-Stream (Nachrichtenschema, RPC-Name, Stream-Semantik)

Technische Ausgestaltung von [`LH-FA-SST-008`](lastenheft.md): ein
gRPC-Server-Streaming-RPC, der jedem verbundenen Consumer jeden einzelnen
committed Change mit vollständigem Inhalt überträgt. Die
Nachvollziehbarkeit bleibt beim bestehenden Lesezugriffsweg
([`LH-FA-REA-001`](lastenheft.md) ff.) und der bestätigten Consumer-Position
([`LH-FA-CON-003`](lastenheft.md)/[`LH-FA-CON-005`](lastenheft.md)); der Stream selbst
trägt kein Replay ([`LH-FA-SST-008`](lastenheft.md) Boundary).

| Merkmal | Festlegung |
|---|---|
| Protokoll / Dienst | gRPC über HTTP/2 mit Protobuf (`proto3`); Paket `cdc.stream.v1`, Dienst `ChangeStream`, Quelldatei `proto/cdc/stream/v1/changestream.proto` |
| RPC | `StreamChanges(StreamChangesRequest) returns (stream Change)` — ein Server-Streaming-Aufruf: ein Öffnungsversuch, viele Antwortnachrichten über die Zeit |
| Request | `StreamChangesRequest` trägt drei optionale, unabhängig setzbare Felder `schema`/`table`/`target` (string, Feldnummern 1/2/3); ein leeres Feld trägt keinen Filter auf dieser Dimension, ein leerer Request (der proto3-Zero-Value, auch einer noch nicht neu generierten `StreamChangesRequest{}` alter Clients) liefert wie zuvor jeden Change ungefiltert. Ein gesetztes `schema` ohne `table` filtert auf alle Tabellen dieses Schemas, ein gesetztes `table` ohne `schema` auf jede Tabelle dieses Namens unabhängig vom Schema, beide gesetzt filtert exakt eine Tabelle. Ein gesetztes `target` liefert nur Changes, deren Zustellziel (`route_target`, `SPEC-002`) dieser Name ist, als Konjunktion mit `schema`/`table`; ein Name, den keine Change trägt (auch einer außerhalb des Alphabets des Zielnamens, `SPEC-032`), liefert keine Nachricht und keinen Fehler; das gilt auch für einen Wert mit dem Zeichen U+0000, der nie ein vergebenes Ziel trifft (der Vergleich läuft im Speicher des Prozesses; der SQL-Zugriff `cdc.changes` ist ausgenommen, der Aufrufer bestimmt sein Prädikat selbst). Die Filterprüfung läuft im Driving-Handler nach dem Empfang aus dem `Broadcaster`, dessen Fan-out unverändert ungefiltert bleibt |
| Nachricht `Change` | dieselben Felder wie der Domain-Typ `model.Change` (`OldImage`/`NewImage`, `internal/domain/model/change.go`): `change_id` (string), `transaction_id` (string), `source_table_id` (string), `sequence` (int64), `operation` (string, eine der drei Operationen `INSERT`, `UPDATE`, `DELETE`), `old_image` (bytes), `new_image` (bytes), `schema_version` (string), `schema` (string), `table` (string); die Felder `origin` und `route_target` (`SPEC-002`) gehören nicht zur Nachricht |
| Granularität | eine Nachricht je Zeilen-Change der committed Transaktion, in deren Reihenfolge — keine Deduplizierung nach Tabelle |
| Zustellgarantie | keine (Fire-and-Forget, verlustbehaftet): ein Consumer, der nicht verbunden ist **oder langsamer liest als Changes eintreffen**, verpasst die betroffenen Nachrichten ersatzlos — je Abonnent trägt der `Broadcaster` eine begrenzte Empfangs-Warteschlange, deren Überlauf verworfen wird. Ein Replay innerhalb des Streams gibt es nicht; verpasste Changes bleiben über den bestehenden Lesezugriffsweg ([`LH-FA-REA-001`](lastenheft.md) ff.) und die bestätigte Consumer-Position ([`LH-FA-CON-003`](lastenheft.md)/[`LH-FA-CON-005`](lastenheft.md)) nachholbar ([`LH-FA-SST-008`](lastenheft.md) Boundary) |
| Erzeuger-Blockade | keine: `Publish` blockiert nie auf einen Abonnenten — auch ein verbundener, gerade nicht lesender Consumer hält den Capture-Pfad nicht an. Der Aufruf liefert nur bei ungültigem Aufruf (fehlender Change oder bereits beendeter Aufruf-Kontext) einen Fehler; ein Fehlschlag geht nicht in den Rückgabewert des Capture-Aufrufs ein (Zeile *Fehler bei Publish-Fehlschlag*) |
| Fehler bei Publish-Fehlschlag | der Aufrufer (`CaptureService`) fängt ihn ab; er geht nicht in den Rückgabewert des Capture-Aufrufs ein und beeinflusst weder die bereits erfolgte Persistierung noch das Source-ACK ([`LH-QA-REL-001.a`](lastenheft.md)) |
| Authentifizierung | gRPC-Metadata-Eintrag `authorization` in der Wertform `Bearer <token>` — dieselben zwei Token-Klassen wie die HTTP-API (`SPEC-018`, `CDC_API_TOKEN_READER`/`CDC_API_TOKEN_ADMIN`); ein fehlender oder keiner Klasse entsprechender Wert endet mit gRPC-Status `Unauthenticated`, nicht mit einem stillen leeren Stream |
| Aktivierung | optional über `CDC_GRPC_ADDR`; ungesetzt bedeutet deaktiviertes Feature, kein Listener, unverändertes Bestandsverhalten |

### SPEC-021 — HTTP-Server-Sent-Events für den Live-Change-Stream

Technische Ausgestaltung von [`LH-FA-SST-008`](lastenheft.md), zweiter Zustellweg neben dem gRPC-Stream (`SPEC-020`): ein lang laufender Server-Stream im bestehenden HTTP-Adapter, auf demselben `Broadcaster`. Eigener Eintrag statt einer Erweiterung von `SPEC-018`, weil jener Abschnitt ausdrücklich die neun Port-gedeckten Fähigkeiten von [`LH-FA-SST-006`](lastenheft.md) ausgestaltet und der Stream keine davon ist.


| Merkmal | Festlegung |
|---|---|
| Endpunkt / Methode | `GET /changes/stream` |
| Rechtsklasse | `reader` oder `admin` — Streaming ist rein lesend |
| Query-Parameter | `schema`, `table`, `target` (je optional und unabhängig, dieselben Feldnamen und dieselbe Kombinatorik wie `SPEC-020`s gRPC-Request und `GET /changes`, `SPEC-022`; `target` wählt das Zustellziel, `SPEC-020`; ein `target`, das keine Change trägt, auch eines außerhalb des Alphabets des Zielnamens oder mit dem Zeichen U+0000, liefert keinen Event und keinen Fehler); ohne Parameter liefert der Endpunkt wie zuvor jeden Change ungefiltert. Ein Parameter außerhalb dieser Menge endet mit `400`, vor jedem SSE-Event |
| Response-Form | `Content-Type: text/event-stream`; je Change ein Event, sofort über `http.Flusher` ausgeliefert |
| Event-Typ | `event: change` |
| Event-Daten | `data:` trägt ein JSON-Objekt mit denselben zehn Feldern wie der Domain-Typ `model.Change`, ohne die Felder `origin` und `route_target` (`SPEC-002`): `change_id` (string), `transaction_id` (string), `source_table_id` (string), `sequence` (int64), `operation` (string, `INSERT`/`UPDATE`/`DELETE`), `old_image`, `new_image`, `schema_version` (string), `schema` (string), `table` (string). Die Row Images stehen als eingebettete JSON-Werte; ein fehlendes Bild ist `null` |
| Zustellgarantie | keine (Fire-and-Forget): ein nicht verbundener **oder langsamer lesender** Consumer verpasst die betroffenen Nachrichten ersatzlos; ein Erzeuger hält nie auf einen Empfänger an. Verpasste Changes bleiben über den bestehenden Lesezugriffsweg ([`LH-FA-REA-001`](lastenheft.md) ff.) und die bestätigte Consumer-Position ([`LH-FA-CON-003`](lastenheft.md)/[`LH-FA-CON-005`](lastenheft.md)) nachholbar |
| Replay | kein Stream-internes Replay; der `Last-Event-ID`-Header wird weder gesendet noch ausgewertet |
| Aktivierung | wie die übrigen Endpunkte über `CDC_HTTP_ADDR`; ungesetzt bedeutet deaktiviertes Feature, kein HTTP-Server. Ist die Adresse gesetzt, aber kein `Broadcaster` verdrahtet, antwortet der Endpunkt mit `503` |

### SPEC-022 — HTTP-API: Changes lesen (`GET /changes`)

Technische Ausgestaltung von [`LH-FA-SST-006`](lastenheft.md) für das
**Changes-Lesen** ([`LH-FA-REA-001`](lastenheft.md)…[`006`](lastenheft.md)):
der lesende Netzwerkzugriffsweg über denselben `ChangeStorePort`, den der
SQL-View-Direktzugriff auf `cdc.changes` trägt — **kein** zweiter
Lesepfad, dieselbe Bereichs-, Limit- und Filtersemantik, dieselbe
deterministische Sortierung (`LH-FA-SST-006` Boundary). Eigener Eintrag
statt einer Erweiterung von `SPEC-018`, weil jener Abschnitt die neun
Port-gedeckten Fähigkeiten ausgestaltet.

| Fähigkeit | Endpunkt | Rechtsklasse | Request | Response |
|---|---|---|---|---|
| `ReadChanges` ([`LH-FA-SST-006`](lastenheft.md), [`LH-FA-REA-001`](lastenheft.md)…[`006`](lastenheft.md)) | `GET /changes` | `reader` oder `admin` | Query-Parameter `source` (**Pflicht**), `schema`, `table`, `target` (je optional und **unabhängig**), `from`, `to` (optional, `commit_position`-Werte ≥ 1, `from` **inklusiv** / `to` **exklusiv**), `limit` (optional, ≥ 1) — **kein** Default-Limit | `200`: `{"changes": [{"commit_position": <int64>, "change_id": "<string>", "transaction_id": "<string>", "source_table_id": "<string>", "schema": "<string>", "table": "<string>", "sequence": <int64>, "operation": "<INSERT\|UPDATE\|DELETE>", "old_image": <eingebettetes JSON \| null>, "new_image": <eingebettetes JSON \| null>, "schema_version": "<string>", "committed_at": "<RFC 3339 in UTC, Bruchteil-Sekunden mit bis zu neun Stellen (abschließende Nullen entfallen)>", "origin": "<wal\|backfill>"}, …]}` — kein Treffer → leere Liste, **nie `404`** |

Die Feldnamen folgen dem API-Vokabular, nicht dem Spaltenvokabular der
View: `schema`/`table` statt `schema_name`/`table_name` (wie
`listTablesResponse`), `old_image`/`new_image` statt `old_data`/`new_data`
(wie der gRPC- und der SSE-Stream, `SPEC-020`/`SPEC-021`). `source_table_id`
steht zusätzlich daneben — die API adressiert Tabellen an anderer Stelle
über `table_id` (`POST /tables/enable`).

| Merkmal | Festlegung |
|---|---|
| Reihenfolge | deterministisch nach (`commit_position`, `transaction_id`, `sequence`) — [`LH-FA-REA-004`](lastenheft.md); die Fortsetzung ist `from = <letzte gelieferte commit_position> + 1`, wenn das Lesen die letzte Position vollständig erfasst hat; enthält sie mehr Changes als `limit`, gilt die Zeile „Position und `limit`" |
| Leere Menge | `{"changes": []}`, nie `null` — ohne Treffer (unbekannte Quelle, unbekanntes Schema, unbekannte Tabelle, leerer Bereich) endet der Aufruf `200`; ein leerer Bestand ist kein Fehler ([`LH-FA-REA-006`](lastenheft.md) Boundary) |
| Fehler-Antwortform | `{"error": "<Klartext>", "code": "<Meldungscode>"}` wie in `SPEC-018` (Codes des Bereichs `E8` für die Ablehnungen); `400` für ein fehlendes `source`, einen **Parameter außerhalb der Liste** (strenger als die neun Bestandsendpunkte — ein unbekannter *Filter* änderte den Ergebnisstand sonst still), eine nicht als Ganzzahl lesbare Zahl, `from`/`to` `< 1`, `limit` `< 1` ([`LH-FA-REA-003`](lastenheft.md) Negative) oder `from > to` ([`LH-FA-REA-001`](lastenheft.md) Negative); fehlender/unbekannter Bearer-Token `401`; Store-Fehler der Klasse `storage` `500`. Kein `404`-Pfad: das Lesen prüft nichts an der Quelle, es liest einen Bestand |
| Herkunft | `origin` trägt `wal` für einen über den Replication Stream erfassten Change und `backfill` für einen Bestands-Change (`SPEC-002`, `LH-FA-CAP-009.a`); ein gespeicherter Change ohne das Feld liest als `wal`. Der Abschnitt bietet keinen Filter auf `origin` |
| Zustellziel | `target` ist ein optionaler Filter auf das Zustellziel der Change (`route_target`, `SPEC-002`): leer oder ohne den Parameter kein Filter, gesetzt nur Changes mit genau diesem Ziel, als Konjunktion mit `schema`/`table`. Ein Name, den keine Change trägt (auch einer außerhalb des Alphabets des Zielnamens, `SPEC-032`), ist bei sonst gültiger Anfrage kein `400`, sondern ein Filter ohne Treffer: `{"changes": []}`; das gilt auch für einen Wert mit dem Zeichen U+0000, weil der gemeinsame Use Case das Alphabet vor dem Speicherzugriff prüft und bei Verletzung die leere Liste liefert, ohne den Speicher anzufragen. Ein Fehler des Lese-Kontrakts (`limit`, Bereich, Quelle) hat Vorrang: er endet mit demselben Fehler wie ohne `target`. Der SQL-Zugriff `cdc.changes` ist davon ausgenommen: der Aufrufer bestimmt sein Prädikat auf `route_target` selbst. Das Ziel steht nicht in der Antwort; die View `cdc.changes` trägt es als Spalte |
| Position und `limit` | Ein `limit` schneidet Zeilen, nicht Positionen: enthält eine Commit-Position mehr Changes als `limit`, liefert `from = <letzte gelieferte commit_position> + 1` den Rest dieser Position nicht, und ein Lesen ab derselben Position liefert wieder dieselben Zeilen. Ein Bestandsabzug legt alle seine Changes auf **eine** Position; er wird ohne `limit` gelesen (oder über den Schlüsselvergleich auf der View `cdc.changes`) |
| Noch nicht begrenzt | Kein Default-Limit und **keine** harte Obergrenze: ohne `limit` liest der Aufruf unbegrenzt, wie der View-Direktzugriff. Eine eingebaute Grenze wäre eine eigene Festlegung dieses Abschnitts |
| Aktivierung | wie die übrigen Endpunkte über `CDC_HTTP_ADDR`; ungesetzt bedeutet deaktiviertes Feature, kein HTTP-Server |

### SPEC-023 — Beispiel-Clients

Technische Ausgestaltung der **öffentlichen Zugriffswege für Integratoren**:
die Klasse der Beispiel-Clients, die einen dokumentierten Draht real
ansprechen. Eigener Eintrag statt einer Erweiterung der Draht-Festlegungen
(`SPEC-018`, `SPEC-020`, `SPEC-021`, `SPEC-022`), weil die Beispiele keine
davon ausgestalten: sie **benutzen** sie, laufen nicht im Feed und tragen
keine Anforderung des Lastenhefts. Der Eintrag trägt die Form für **alle**
Sprachen — die Clients, die es gibt, sind damit verankert, nicht nur die,
die hinzukommen.

| Merkmal | Festlegung |
|---|---|
| Klasse | Ein Beispiel-Client ist ein eigenständiges Programm, das **eine** dokumentierte Zugriffs-Oberfläche real anspricht und ihre Antwort ausgibt. Er ist **Vorbild** — lesbar, kopierbar, startbar —, nicht Belegträger: die E2E-Belege des Repos tragen die Wegwerf-Clients des Harness. Er trägt **keine** Zustandsmaschine (Reconnect, Deduplizierung, Rückstand) |
| Verhältnis zum Draht | Das Beispiel **benutzt** die Festlegungen dieses Dokuments (`SPEC-018` HTTP-Endpunkte, `SPEC-020` gRPC-Stream, `SPEC-021` SSE, `SPEC-022` Changes lesen) und fügt ihnen nichts hinzu. Verlangt ein Beispiel eine Vertragsänderung, ist das eine Spec-Änderung, kein Beispiel-Umbau |
| Ort und Form | `examples/` auf der Repo-Wurzel; **je Client ein Verzeichnis mit einem Programm und einem Einstiegspunkt**; die Namen tragen `-client`. Verlangt die Werkzeugkette einer Sprache ein Projektverzeichnis, liegt der Client unter einem **Sprach-Wurzelverzeichnis** (`examples/<sprache>/<client>/`); Go verlangt das nicht und bleibt flach (`examples/<client>/`) |
| Import-Grenze | Ausschließlich die Standardbibliothek/Runtime der Sprache und **öffentliche** Fremdmodule; **kein** Import eines privaten Baums dieses Repositories. Die Quelle bleibt außerhalb dieses Repositories kopierbar, übersetzbar bzw. nachbaubar — ein Beispiel, das nur hier baubar ist, verfehlt seinen einzigen Leser |
| Laufzeit und Startform | **Docker-only**: kein Host-Compiler. Jede Sprache wird aus einem **digest-gepinnten** Basis-Image gebaut, Abhängigkeiten sind auf feste Versionen gepinnt, und die Startform ist ein Container-Aufruf. Adresse, Token und die fachlichen Parameter kommen aus denselben Umgebungsvariablen, die §5 führt, und lassen sich per Flag übersteuern |
| Bau- und Prüfweg | **Je Sprache ein Bau-/Testziel**, das die Beispiele der Sprache übersetzt und die **netzlos** prüfbaren Teile testet (Aufbau der Anfrage, Zerlegen eines Stream-Frames). Das Ziel ist **Werkzeug, kein Gate**: die Gate-Kette des Repos bleibt netzlos, ein Bauziel hängt nicht an ihr |
| Verhältnis zur Konfigurationsdatei | Beispiele lesen `CDC_CONFIG_FILE` **nicht**: sie beziehen ihre Eingaben aus Umgebungsvariablen und Flags. Die Konfigurationsdatei ist der Deployment-Eingang des Feeds, nicht der eines Integrator-Programms (`SPEC-016`) |
| Handbuch | Jedes Beispiel wird im Zugriffs-Abschnitt **seiner** Oberfläche namentlich mit Programm-Pfad und Startform genannt; die Änderungshistorie des Handbuchs trägt die Zeile. Beispiel und Handbuch-Zeile gehören in **denselben** Zug |
| Sprachen und Umfang | Das Repo führt Beispiel-Clients in **Go**, **C#** und **Kotlin** über die **volle Matrix**: alle fünf Zugriffs-Oberflächen (`SPEC-018`, `SPEC-020`, `SPEC-021`, `SPEC-022`, `SPEC-024`) in jeder der drei Sprachen. Ein fremdsprachiger gRPC-Client erreicht die `.proto`-Quelle über **einen** zusätzlichen, benannten Bau-Kontext (das Sprach-Wurzelverzeichnis bleibt der Bau-Kontext); die daraus erzeugten Fremdsprachen-Stubs entstehen **im Bau** und werden **nicht** committet — `gen/**` bleibt die Go-Bindung. Ein NATS-Client bekommt je Sprache eine gepinnte öffentliche Client-Bibliothek — dieselbe Bibliothek trägt beide NATS-basierten Zugriffswege (Wecksignal, Vollinhalts-Stream) je Sprache; ein fremdsprachiger Vollinhalts-Stream-Client kann zusätzlich eine gepinnte JSON-Bibliothek brauchen, wenn die Sprache selbst keine öffentliche JSON-Dekodierung trägt |
| Kein Belegträger | Beispiele erscheinen **nicht** in der E2E-Abdeckungstabelle; sie sind Doku mit Bau-Bindung. Was sie nicht leisten — „übersetzt" ist nicht „holt am laufenden Feed eine Änderung" — bleibt eine benannte Grenze, deren Wächter das Review ist |

### SPEC-024 — NATS-Vollinhalts-Stream für Live-Streaming (Subjekt- und Nachrichtenform)

Technische Ausgestaltung von [`LH-FA-SST-008`](lastenheft.md), dritter
Zustellweg neben dem gRPC-Stream (`SPEC-020`) und SSE (`SPEC-021`): derselbe
`Broadcaster` bekommt einen dritten Abonnenten, der jeden Change zusätzlich
über NATS Core veröffentlicht. Eigenständig von [`SPEC-017`](pflichtenheft.md)
(NATS-Wecksignal für [`LH-FA-SST-007`](lastenheft.md)) — beide Fähigkeiten
laufen über denselben Server und dieselbe Verbindung, aber mit
unterschiedlichem Subjekt-Namensraum, unterschiedlichem Nachrichteninhalt
und unterschiedlicher Kardinalität; `SPEC-017` bleibt durch diesen Eintrag
unverändert.

| Merkmal | Festlegung |
|---|---|
| Subjekt-Schema | `cdc.stream.<source_id>.<schema>.<table>` — ein Subjekt je Tabelle einer Quelle, strukturell wie `SPEC-017`s Wecksignal-Subjekt, aber mit dem Wurzel-Token `cdc.stream` statt `cdc.changes`, um `SPEC-017`s leeren Payload nicht zu berühren. Ein Consumer mit Tabelleninteresse abonniert das vollständige vierstufige Subjekt; `cdc.stream.<source_id>.>` deckt alle Tabellen einer Quelle, `cdc.stream.>` alle Quellen |
| Zusatz-Subjekt (Zustellziel) | `cdc.route.<source_id>.<ziel>` — jede Change mit einem Zustellziel (`route_target`, `SPEC-002`) wird zusätzlich auf diesem Subjekt veröffentlicht, mit demselben Payload wie auf dem `cdc.stream…`-Subjekt; eine Change ohne Ziel erzeugt keine zweite Nachricht. Das Subjekt `cdc.stream.<source_id>.<schema>.<table>` bleibt unverändert, das Alphabet des Zielnamens (`SPEC-032`) enthält weder Punkt noch NATS-Platzhalter. Ein Consumer abonniert `cdc.route.<source_id>.<ziel>` für ein Ziel, `cdc.route.<source_id>.>` für alle Ziele einer Quelle. Die zweite Veröffentlichung ist eine Zusage an die Umsetzung; die Zeit für `publish` samt Flush ist am Testcontainer-NATS ohne Abonnent und ohne Schwelle gemessen (Median von fünf Läufen, die Zahl streut zwischen den Läufen); eine Last-Zusage folgt daraus nicht. Backfill-Changes gehen in keinen Live-Weg (`LH-FA-CAP-009.a`), also auch nicht auf dieses Subjekt |
| Nachrichteninhalt | JSON-Objekt mit denselben zehn Feldern wie `SPEC-021`s SSE-Event: `change_id`, `transaction_id`, `source_table_id`, `sequence`, `operation` (`INSERT`/`UPDATE`/`DELETE`), `old_image`, `new_image`, `schema_version`, `schema`, `table` — dasselbe Nachrichtenschema, kein drittes |
| Granularität | eine Nachricht je Zeilen-Change der committed Transaktion, wie `SPEC-020`/`SPEC-021` — **nicht** `SPEC-017`s Tabellen-Dedup-Kardinalität |
| Zustellgarantie | keine (Core NATS, Fire-and-Forget) — ein nicht verbundener oder langsamer Consumer verpasst Nachrichten ersatzlos; Nachholen ausschließlich über den bestehenden Lesezugriffsweg ([`LH-FA-REA-001`](lastenheft.md) ff.) und die bestätigte Consumer-Position ([`LH-FA-CON-003`](lastenheft.md)/[`LH-FA-CON-005`](lastenheft.md)) |
| Erzeuger-Blockade | keine: Der Publisher liest ausschließlich aus dem bereits vom `Broadcaster` isolierten Kanal, dieselbe Fehlerisolation wie `SPEC-020`/`SPEC-021` |
| Authentifizierung | Verbindungsebene: Der NATS-Server verlangt einen Token, sobald `CDC_NATS_STREAM_TOKEN` konfiguriert ist — ein Verbindungsversuch ohne oder mit falschem Token wird vom Server abgelehnt. Dieser Token gilt serverweit (auch für die bislang anonyme `SPEC-017`-Verbindung), sobald er konfiguriert ist |
| Aktivierung | optional, **beide** Bedingungen: `CDC_NATS_URL` **und** `CDC_NATS_STREAM_TOKEN` gesetzt. Nur `CDC_NATS_URL` gesetzt bedeutet unverändertes `SPEC-017`-Bestandsverhalten, kein dritter Weg. `CDC_NATS_STREAM_TOKEN` gesetzt ohne `CDC_NATS_URL` ist ein Konfigurationsfehler beim Start (`configuration`, `SPEC-008`) |

### SPEC-029 — `cdc.backfill_run` und `cdc.backfill_status` (Run-Zustand eines Backfills)

Feldform des Run-Zustands zu `LH-FA-CAP-009.a` und der lesenden View
([`LH-FA-CAP-009`](lastenheft.md), [`LH-FA-SST-003`](lastenheft.md)): eine
Zeile je Run, angelegt bei der Annahme eines Antrags der Antragsart
`backfill` (`SPEC-019`), fortgeschrieben vom Backfill-Worker.

| Spalte | Typ | Pflicht | Bedeutung |
|---|---|---|---|
| `run_id` | text (PK) | ja | die `administration_request_id` des annehmenden Antrags |
| `source_id` | text (FK `cdc.source`) | ja | Quelle des Runs |
| `schema_name` / `table_name` | text | ja | kopierte Tabelle |
| `status` | text | ja | geschlossene Menge `queued` \| `running` \| `completed` \| `failed` \| `interrupted`; die Menge steht bei Erstanlage der Tabelle als Prüfbedingung |
| `requested_at` | timestamptz | ja | Anlage-Zeitpunkt; die Aufnahme-Ordnung ist (`requested_at`, `run_id`) |
| `started_at` | timestamptz | nein | Übergang nach `running`, NULL solange `queued`; der Beginn der Kopierdauer — die Wartezeit in `queued` zählt nicht |
| `finished_at` | timestamptz | nein | Übergang nach `completed`, `failed` oder `interrupted` |
| `snapshot_position` | bigint | nein | die Position `X` des Runs (`LH-FA-CAP-009.a`, `SPEC-003`); NULL bis zur Anlage des Snapshots |
| `rows_copied` | bigint | ja, Default 0 | Fortschritt je Block, außerhalb der Daten-Transaktion fortgeschrieben und damit für Leser sichtbar; bei `completed` die Zahl der geschriebenen Backfill-Changes |
| `estimated_rows` | bigint | nein | die beim Antrag **geschätzte** Zeilenzahl der Tabelle (eine Schätzung der Quelle, keine Zählung); NULL heißt „unbekannt", **nie** `0` |
| `warn_estimated_size` | boolean | ja, Default `false` | Warnung (1): die **geschätzte** Zeilenzahl liegt über der Richtgröße; bleibt `false`, solange `estimated_rows` unbekannt ist oder keine Richtgröße festgelegt ist. Schreibt die Annahme |
| `warn_duration` | boolean | ja, Default `false` | Warnung (2): die Kopierdauer hat die Toleranz überschritten; gesetzt beim Fortschritts-Update (je Block) oder beim Abschluss und danach unverändert. Schreibt der Worker |
| `error_message` | text | nein | Fehlertext eines `failed`-Runs in der Form `<Klasse> [<Code>]: <Ursache>` (Klasse und Meldungscode aus `SPEC-008`); sonst NULL |

Die beiden Warn-Spalten sind eine Kennzeichnung, kein Wert: sie nennen weder
Toleranz noch Richtgröße, und eine gesetzte Warnung ändert weder `status` noch
den Ablauf des Runs — keine Ablehnung, kein Abbruch. Es sind zwei Spalten mit
je einem Schreiber: die Annahme fügt nur ein und der Worker schreibt nur
fort (Grants unten), jede Spalte hat deshalb genau eine schreibende Rolle.

**Grants** (Rollen nach der Zuordnung der DSN-Verdrahtung):

| Rolle | Recht auf `cdc.backfill_run` | Träger |
|---|---|---|
| `cdc_admin` | `SELECT`, `INSERT` | die Annahme: Prüfung „kein aktiver Run derselben Tabelle" und Anlage der Zeile `queued` |
| `cdc_capture` | `SELECT`, `UPDATE` | der Worker: `queued`-Zeilen lesen, Statuswechsel, `started_at`/`finished_at`, `snapshot_position`, `rows_copied`, `error_message`, `warn_duration`; der Abgleich `running` → `interrupted` beim Prozessstart |
| `cdc_reader` | **keines** | liest ausschließlich über die View |

Niemand trägt `DELETE`; `cdc_admin` trägt kein `UPDATE`, `cdc_capture` kein
`INSERT`. `cdc_reader` trägt `SELECT` auf die View `cdc.backfill_status`.

**View `cdc.backfill_status`:** je Tabelle (`source_id`, `schema_name`,
`table_name`) die Zeile des zuletzt beantragten Runs (größtes `requested_at`,
Zweitschlüssel `run_id`) mit denselben Spalten wie `cdc.backfill_run`. Jede
Ausgabe der **geschätzten** Zeilenzahl (View-Leser, CLI-Diagnose) zeigt NULL
als „unbekannt" und nie als `0`.

### SPEC-030 — Transformationsregel (`rule_spec`) und ihre Wirkung auf Row Images

Feldform und Semantik der Regel zu `LH-FA-CFG-007.a`
([`LH-FA-CFG-007`](lastenheft.md)): das Feld `rule_spec` der Antragsart
`set_transformation` (`SPEC-019`) und die Wirkung der Regel auf `old_data` und
`new_data` einer Change (`SPEC-002`).

`rule_spec` ist ein JSON-Objekt mit dem Pflichtschlüssel `kind`. Jeder
Regeltyp kennt genau die Schlüssel seiner Zeile; ein anderer Schlüssel und ein
unbekannter `kind` enden den Antrag `failed` (Fehlertexte in `SPEC-019`).

| `kind` | Schlüssel (alle Pflicht) | Form | Wirkung auf das Row Image (`old_data` und `new_data`) |
|---|---|---|---|
| `rename_column` | `column`, `to` | zwei Zeichenketten, Form siehe Bezeichner | der Schlüssel `column` heißt im Image `to`; der Wert bleibt unverändert |
| `map_value` | `column`, `values` | `column`: Zeichenkette, Form siehe Bezeichner; `values`: nicht leeres Objekt `alt → neu`, dessen Werte Zeichenketten sind | ist der Wert von `column` als Zeichenkette ein Schlüssel von `values`, steht der zugeordnete Wert im Image; jeder andere Wert bleibt unverändert |

- **Abwesenheit bleibt Abwesenheit.** Fehlt der Wert der Spalte im Image
  (`NULL`, unverändertes TOAST, ausgeschlossene oder generierte Spalte), bleibt
  er für beide Regeltypen abwesend — dieselbe Vereinbarung wie ohne Regel
  ([`LH-FA-DAT-005`](lastenheft.md)), kein Platzhalter, kein Zielschlüssel
  ohne Wert. Ein fehlendes Bild einer Operation bleibt fehlend
  ([`LH-FA-CAP-008`](lastenheft.md)).
- **Bezeichner.** `column`, `to` und der Regelname (`rule_name`, `SPEC-019`)
  werden **zeichengenau** verglichen: Groß-/Kleinschreibung wird nicht gefaltet
  (`Name` und `name` sind verschiedene Namen), Anführungszeichen und Punkte
  tragen keine Bedeutung, ein Name wird nicht als quotierter Bezeichner
  gelesen. Das ist der Vergleich der Spaltenauswahl: `cdc.exclude_column` und
  `cdc.include_column` prüfen den Spaltennamen zeichengenau gegen den Katalog
  der Quelle (`SPEC-019`), der Ausschlussstand vergleicht Namen zeichengenau.
  - `column` ist die Quellspalte: nicht leer, ohne das Zeichen U+0000; sie muss
    zeichengenau als Spaltenname der Quelltabelle im Katalog stehen (K4). Eine
    Längen- oder Zeichenprüfung über den Katalog hinaus gibt es nicht.
  - `to` ist der Zielname und damit der Schlüssel im Image: nicht leer, ohne
    das Zeichen U+0000, höchstens 63 Byte in UTF-8 (die Bezeichner-Länge der
    Quelle); jedes weitere Zeichen ist zulässig, der Name wird nicht als
    Bezeichner der Quelle gedeutet. K3 vergleicht `to` zeichengenau mit den
    Spaltennamen der Quelltabelle (Katalog, ausgeschlossene Spalten
    eingeschlossen) und mit dem `to` jeder anderen Regel der Tabelle.
  - `rule_name` hat 1 bis 63 Zeichen aus `a`–`z`, `0`–`9` und `_` — das
    Alphabet, das die Aktivierung für Schema- und Tabellennamen verlangt. Ein
    Name mit Großbuchstaben liegt außerhalb des Alphabets und endet den Antrag
    `failed`, er wird nicht gefaltet. K1 vergleicht zeichengenau.
  - Die Adresse `schema.table.name` (`SPEC-019`) ist eindeutig zerlegbar:
    Schema und Tabelle tragen das Alphabet oben und damit keinen Punkt; alles
    nach dem zweiten Punkt ist der Name, auch wenn er Punkte trägt.
- **Vergleich.** `map_value` vergleicht den Wert im Text-Stand der Quelle mit
  den Schlüsseln von `values` zeichengenau: Groß-/Kleinschreibung zählt, es wird
  nichts abgeschnitten oder normalisiert.
- **Reihenfolge und Inhalt.** Die Spaltenregeln werden in der Spaltenreihenfolge
  der Relation der Change ausgewertet, nach dem Spaltenausschluss. K2 und K3
  machen das Ergebnis von dieser Reihenfolge unabhängig: gleiche Regelmenge und
  gleiche Relation ergeben dieselbe Schlüsselmenge mit denselben Werten. Die
  Reihenfolge der Schlüssel im Image ist nicht zugesagt, `jsonb` bewahrt sie
  nicht (`SPEC-002`); zugesagt ist der Inhalt als JSON-Wert — Schlüsselmenge
  und Werte.
- **Randfälle.** Ein leeres `values`-Objekt ist eine formale Verletzung (eine
  Regel ohne Zuordnung wirkt nie) und endet den Antrag `failed`. Ein
  Zielname gleich dem Quellnamen verletzt K3 (`SPEC-019`: die Quellspalte
  selbst ist ein Spaltenname der Quelltabelle) und endet den Antrag `failed`.
  Die Abbildung eines Werts auf sich selbst (`"a": "a"`) ist zulässig und
  ändert diesen Wert nicht. Mehrere Quellwerte dürfen auf denselben Zielwert
  abgebildet werden; die Abbildung ist dann nicht umkehrbar.
- **Anwendbarkeit** (die führende Stelle; `LH-FA-CFG-007.a` und `SPEC-008`
  verweisen hierher). Eine Regel ist auf eine Change anwendbar, wenn ihre
  Spalte in der Relation der Change vorkommt und, für `rename_column`, ihr
  Zielname mit keiner Spalte der Relation kollidiert (zeichengenau, siehe
  Bezeichner). Im Run eines Backfills gilt dieselbe Definition gegen die
  Spalten des Snapshots. Anwendbarkeit hängt an Regelmenge und Spaltenmenge,
  nie am Wert einer Zeile: kein Wert macht eine Regel unanwendbar. Eine nicht
  anwendbare Regel endet mit der Fehlerklasse `schema` (`SPEC-008`).
- **Wirkung nur auf das Bild.** Die Regeln ändern weder ein anderes Feld der
  Change noch die Tabellen-Identität (`SPEC-002`); die Row Images bleiben
  JSON-Objekte mit Zeichenketten-Werten.

**Beispiele** (Tabelle `public.orders` mit den Spalten `id`, `name`, `status`;
die Schlüsselreihenfolge der Images ist nicht zugesagt):

- `rename_column`: Regelname `kundenname`, `rule_spec`
  `{"kind": "rename_column", "column": "name", "to": "customer_name"}`. Eine
  Change mit dem Image `{"id": "7", "name": "Ada", "status": "o"}` trägt danach
  das Image `{"id": "7", "customer_name": "Ada", "status": "o"}`: der Schlüssel
  `name` fehlt, der Wert `Ada` steht unverändert unter dem Zielnamen.
- `map_value`: Regelname `status_lesbar`, `rule_spec`
  `{"kind": "map_value", "column": "status", "values": {"o": "open", "c": "closed"}}`.
  Der Wert `o` steht danach als `open`, der Wert `x` bleibt `x`; ist `status`
  `NULL`, fehlt der Schlüssel wie ohne Regel.
- Abgelehnter Antrag: `rule_spec`
  `{"kind": "rename_column", "column": "name", "to": "id"}` gegen dieselbe
  Tabelle endet `failed` mit dem Fehlertext
  `Zielname kollidiert mit einer Spalte der Tabelle: public.orders.id`
  (`SPEC-019`, K3); der Regelstand bleibt unverändert.
- Nicht anwendbare Regel: Wird der Tabelle nach der Regel `kundenname` die
  Spalte `customer_name` hinzugefügt, kollidiert der Zielname mit einer Spalte
  der Relation; die nächste Change der Tabelle endet im Erfassungspfad mit der
  Fehlerklasse `schema`, bis die Regel entfernt ist (`LH-FA-CFG-007.a`, Abhilfe).

### SPEC-031 — gRPC-Verwaltungs-API: Dienst, RPCs, Nachrichtenschema, Fehlercodes

Technische Ausgestaltung von [`LH-FA-CON-001`](lastenheft.md),
[`LH-FA-CON-003`](lastenheft.md)…[`LH-FA-CON-006`](lastenheft.md),
[`LH-FA-CFG-001`](lastenheft.md)…[`LH-FA-CFG-004`](lastenheft.md),
[`LH-FA-RET-002`](lastenheft.md)…[`LH-FA-RET-004`](lastenheft.md) und
[`LH-FA-REA-001`](lastenheft.md)…[`006`](lastenheft.md) und
[`LH-FA-SST-003`](lastenheft.md) (deckt [`LH-FA-ADM-002`](lastenheft.md)…[`005`](lastenheft.md),
[`LH-FA-RET-005`](lastenheft.md), [`LH-FA-RET-006`](lastenheft.md),
[`LH-FA-CAP-009`](lastenheft.md)) über einen zweiten Zugriffsweg (gRPC):
dieselben elf Fähigkeiten wie `SPEC-018`s/`SPEC-022`s HTTP-Kontrakt,
dieselben Inbound Use Cases — kein zweiter Domänenpfad, fachlich
gleichwertig (`LH-FA-SST-006` Boundary).

**Dienst:** Paket `cdc.administration.v1`, Dienst `Administration`,
Quelldatei `proto/cdc/administration/v1/administration.proto` — eigenständig
von Paket `cdc.stream.v1`/Dienst `ChangeStream` (`SPEC-020`): Letzterer
bleibt durch diesen Eintrag byte-identisch unverändert. Beide Dienste laufen
auf demselben `grpc.Server`, derselben Adresse `CDC_GRPC_ADDR`.

**Authentifizierung** (für alle elf RPCs gleich): gRPC-Metadata-Eintrag
`authorization` in der Wertform `Bearer <token>` — dieselben zwei
Token-Klassen wie die HTTP-API (`SPEC-018`, `CDC_API_TOKEN_READER`/
`CDC_API_TOKEN_ADMIN`); ein fehlender oder keiner Klasse entsprechender Wert
endet mit gRPC-Status `Unauthenticated`, ein bekanntes `reader`-Token gegen
eine `admin`-RPC mit `PermissionDenied`. Ein `admin`-Token erreicht jede RPC.

**RPCs** (alle unär, Request/Response — kein Streaming-Vorteil bei
Anfrage/Antwort-Aufrufen, dieselbe Begründung wie `SPEC-018`):

| RPC | Rechtsklasse | Request-Felder | Response-Felder |
|---|---|---|---|
| `RegisterConsumer` ([`LH-FA-CON-001`](lastenheft.md)) | `admin` | `consumer_id`, `name` — beide Pflicht | `consumer_id`, `name`, `already_registered` |
| `AcknowledgeConsumer` ([`LH-FA-CON-004`](lastenheft.md)) | `admin` | `consumer_id`, `source_id`, `offset` (`uint64`) — alle Pflicht | `consumer_id`, `source_id`, `offset` (`uint64`) |
| `GetConsumerPosition` ([`LH-FA-CON-003`](lastenheft.md), [`LH-FA-CON-005`](lastenheft.md)) | `reader` | `consumer_id` Pflicht | `consumer_id`, `source_id`, `offset` (`uint64`), `acknowledged` (`bool`) |
| `RemoveConsumer` ([`LH-FA-CON-006`](lastenheft.md)) | `admin` | `consumer_id` Pflicht | `consumer_id`, `removed` (`bool`) |
| `EnableTable` ([`LH-FA-CFG-001`](lastenheft.md)) | `admin` | `source`, `schema`, `table`, `table_id`, `schema_version_id`, `version` (`int64`), `publication` — alle Pflicht | `table_id`, `source`, `schema`, `table`, `already_enabled` (`bool`) |
| `DisableTable` ([`LH-FA-CFG-002`](lastenheft.md)) | `admin` | `source`, `schema`, `table`, `publication` — alle Pflicht | `removed` (`bool`), `retained` (`bool`) |
| `GetTableStatus` ([`LH-FA-CFG-003`](lastenheft.md)) | `reader` | `source`, `schema`, `table`, `publication` — alle Pflicht | `enabled` (`bool`), `retained` (`bool`) |
| `ListTables` ([`LH-FA-CFG-004`](lastenheft.md)) | `reader` | `source`, `publication` — beide Pflicht | `tables` (`repeated SourceTable`), `retained` (`repeated SourceTable`) — ohne Aktivierung beide leer |
| `RunRetention` ([`LH-FA-RET-002`](lastenheft.md)…[`004`](lastenheft.md)) | `admin` | `source` Pflicht, `min_age_nanos` (`int64`) ≥ 0 | `deleted` (`int64`) — Anzahl real gelöschter Changes |
| `ReadChanges` ([`LH-FA-REA-001`](lastenheft.md)…[`006`](lastenheft.md)) | `reader` | `source` Pflicht, `schema`, `table` (je optional, unabhängig), `from`, `to` (`uint64`, `0` = nicht gesetzt, Start inklusiv/Ende exklusiv), `limit` (`int64`, `0` = unbegrenzt), `target` (`string`, Feldnummer 7, optional: leer = kein Filter, gesetzt = nur Changes mit diesem Zustellziel, Konjunktion mit `schema`/`table`) | `changes` (`repeated ChangeRecord`) — kein Treffer liefert eine leere, gesetzte Liste (`LH-FA-REA-006` Boundary), kein `NotFound`; ein `target`, das keine Change trägt (auch eines außerhalb des Alphabets des Zielnamens, `SPEC-032`), liefert bei sonst gültiger Anfrage eine leere Liste, keinen Fehler; das gilt auch für einen Wert mit dem Zeichen U+0000, weil der gemeinsame Use Case das Alphabet nach der Prüfung des Lese-Kontrakts und vor dem Speicherzugriff prüft und bei Verletzung die leere Liste liefert, ohne den Speicher anzufragen; ein Fehler des Lese-Kontrakts (`limit`, Bereich, Quelle) hat Vorrang und endet mit demselben Code wie ohne `target` (der SQL-Zugriff `cdc.changes` ist ausgenommen: der Aufrufer bestimmt sein Prädikat selbst) |
| `Diagnose` ([`LH-FA-SST-003`](lastenheft.md), deckt [`LH-FA-ADM-002`](lastenheft.md)…[`005`](lastenheft.md), [`LH-FA-RET-005`](lastenheft.md), [`LH-FA-RET-006`](lastenheft.md), [`LH-FA-CAP-009`](lastenheft.md)) | `reader` | `source` Pflicht | `heartbeat` (`HeartbeatStatus`), `capture_lag` (`double`), `consumer_lags` (`repeated ConsumerLag`), `retention_blocker` (`RetentionBlocker`), `storage_bytes` (`double`), `backfill` (`repeated BackfillTableStatus`) — dieselben sechs Signalgruppen wie `GET /diagnose` (`SPEC-018`) |

Hilfsnachricht `SourceTable` (für `ListTables`): `table_id`, `source`,
`schema`, `table` — dieselben vier Felder wie `sourceTableResponse` im
HTTP-Adapter. Hilfsnachricht `ChangeRecord` (für `ReadChanges`):
`commit_position` (`int64`), `change_id`, `transaction_id`,
`source_table_id`, `schema`, `table` (`string`), `sequence` (`int64`),
`operation`, `old_image`, `new_image` (`bytes`, rohe JSON-Bytes, dieselbe
Form wie `streamv1.Change`), `schema_version`, `committed_at` (`string`,
RFC 3339 mit Nanosekunden, UTC), `origin` — dieselben dreizehn Felder wie
`readChangeResponse` im HTTP-Adapter (`SPEC-022`), in derselben Reihenfolge;
das Zustellziel (`route_target`, `SPEC-002`) ist kein Feld von `ChangeRecord`.
Die Wire-Kompatibilität des zusätzlichen Request-Felds ist aus den
Proto3-Regeln hergeleitet, nicht gegen einen alten Server gefahren: ein alter
Client sendet das Feld nicht (leer, kein Filter), ein alter Server ignoriert
ein unbekanntes Feld und liefert ungefiltert.
Alle String-Felder tragen proto3-`string`; `snake_case` im `.proto` wird zu
`camelCase` im generierten Go-/C#-/Kotlin-Code — Formsache der Zielsprache,
keine inhaltliche Abweichung.

Hilfsnachrichten von `Diagnose`: `HeartbeatStatus` (`known` (`bool`),
`age_seconds` (`double`, nur gültig wenn `known`), `error_class` (leer
bedeutet Normalbetrieb, nur gültig wenn `known`) und `error_code`
(Feldnummer 4: der Meldungscode des Fehlerzustands (`SPEC-008`), leer im
Normalbetrieb, nur gültig wenn `known`); ein Leser ohne das Feld ignoriert es, ein Server ohne das
Feld lässt es leer, *abgeleitet* aus den Proto3-Regeln, nicht gegen einen alten
Leser gefahren); `ConsumerLag`
(`consumer_id`, `known` (`bool`), `lag` (`double`, nur gültig wenn
`known`)); `RetentionBlocker` (`present` (`bool`), `consumer_id`, `name`,
`acknowledged_position` (`int64`), `backlog_known` (`bool`), `backlog`
(`int64`, nur gültig wenn `backlog_known`)); `BackfillTableStatus`
(`schema`, `table`, `status`, `rows_copied` (`int64`),
`estimated_rows_known` (`bool`), `estimated_rows` (`int64`, nur gültig wenn
`estimated_rows_known`), `warn_estimated_size` (`bool`), `warn_duration`
(`bool`), `error_message`) — explizite Präsenz-Flags statt eines
0-als-„nicht gesetzt"-Sentinels, weil `0.0` für `age_seconds`/`lag` ein real
möglicher, gültiger Wert ist.

**Fehlercodes** (ersetzt den jeweiligen HTTP-Statuscode aus `SPEC-018`; ein
JSON-Decode-Fehler entfällt strukturell, das Nachrichtenschema ist bereits
typisiert):

| Ursache | gRPC-Code |
|---|---|
| Domänen-Invariante verletzt (Pflichtfeld leer, ungültiger Wertebereich) | `codes.InvalidArgument` |
| fehlender/unbekannter Bearer-Token | `codes.Unauthenticated` |
| bekanntes Token mit unzureichender Rechtsklasse | `codes.PermissionDenied` |
| an der Quelle physisch fehlende Tabelle (nur `EnableTable`/`DisableTable`/`GetTableStatus`) | `codes.NotFound` |
| `ReadChanges` ohne Treffer (unbekannte Quelle/Schema/Tabelle, leerer Bereich) | kein Fehler — `200`-Äquivalent mit leerer, gesetzter `changes`-Liste |
| unerwarteter interner Fehler | `codes.Internal` |

**Meldungscode im Status** (`SPEC-008`): ein Status mit `codes.InvalidArgument`,
`codes.NotFound` oder `codes.Internal` trägt das Statusdetail
`google.rpc.ErrorInfo` mit `reason` = Meldungscode der Ursache und `domain` =
`pg-change-feed`. Die Zuordnung ist die der HTTP-API (`SPEC-018`): eine
Ablehnung einer Aufrufer-Eingabe trägt einen Code des Bereichs `E8`
(`PCF-E8025` Tabelle fehlt an der Quelle, `PCF-E8051` leere Kennung,
`PCF-E8053` Wert außerhalb des Bereichs, `PCF-E8054` Position unter 1,
`PCF-E8055` Position rückt nicht vor, `PCF-E8056` Endposition vor
Startposition, `PCF-E8057` Position einer anderen Quelle), ein
`codes.Internal` den Code der klassifizierten Ursache in der Fehlerkette, sonst
den Rückfall `PCF-E7000`; ein `ChangeStream` ohne verdrahteten Broadcaster
endet mit `codes.Internal` und `PCF-E2001`. `codes.Unauthenticated` und
`codes.PermissionDenied` tragen kein Detail: die Tabelle der Codes führt für sie
keinen Eintrag. Statuscode und Statustext bleiben unverändert; ein Client, der
nur den Statuscode liest, bleibt lauffähig. Der generierte Code ändert sich
nicht: `ErrorInfo` ist ein Detail des Status, kein Feld der Nachrichten.

**Aktivierung:** optional über `CDC_GRPC_ADDR`, dieselbe Variable wie
`ChangeStream` (`SPEC-020`); ungesetzt bedeutet deaktiviertes Feature, kein
Listener, unverändertes Bestandsverhalten.

### SPEC-032 — Routing-Regel (`rule_spec`) und ihre Wirkung auf das Zustellziel

Feldform und Semantik der Regel zu `LH-FA-CFG-008.a`
([`LH-FA-CFG-008`](lastenheft.md)): das Feld `rule_spec` der Antragsart
`set_route` (`SPEC-019`) und die Bestimmung des Feldes `route_target` einer
Change (`SPEC-002`). Die Sätze dieses Eintrags sind Zusagen an die Umsetzung,
keine Messergebnisse; wo die Spec eine Grenze festlegt, die über den Wortlaut
des Lastenhefts hinausgeht, ist das eine Festlegung dieses Eintrags, keine
Messung.

`rule_spec` ist ein JSON-Objekt. Jeder Schlüssel steht in der Tabelle; ein
anderer Schlüssel, auch innerhalb von `when`, endet den Antrag `failed`
(Fehlertexte in `SPEC-019`).

| Schlüssel | Pflicht | Form | Bedeutung |
|---|---|---|---|
| `target` | ja | Zeichenkette, Alphabet siehe Zielname | Name des Zustellziels, das die Regel einer treffenden Change zuweist |
| `order` | ja | JSON-Zahl, positive ganze Zahl | Auswertungsreihenfolge: eine kleinere `order` wird zuerst geprüft |
| `when` | nein | Objekt mit den beiden Pflichtschlüsseln `column` und `equals` (beide Zeichenketten) | Bedingung auf den Textwert einer Spalte der Quelle; ohne `when` trifft die Regel jede Change der Tabelle |

- **Herkunft und Inhalt.** Die Tabelle der Antragszeile ist die Herkunft: eine
  Regel gilt für genau diese Tabelle. Der Inhalt ist die Bedingung `when`;
  weitere Operatoren (Präfix, Menge, Vergleich) kennt die Regelform nicht.
- **Zielname.** `target` hat 1 bis 63 Zeichen: das erste aus `a`–`z` und
  `0`–`9`, die übrigen aus `a`–`z`, `0`–`9`, `_` und `-`. Ein Name außerhalb des
  Alphabets endet den Antrag `failed`, er wird nicht gefaltet oder gekürzt. Der
  Zielname wird zeichengenau verglichen. Das Alphabet enthält weder Punkt noch
  NATS-Platzhalter und ist hergeleitet aus der Subjekt-Syntax, nicht gegen einen
  NATS-Server geprüft.
- **Bedingung.** `when.column` ist die Quellspalte: nicht leer, ohne das Zeichen
  U+0000; sie muss zeichengenau als Spaltenname der Quelltabelle im Katalog
  stehen (R3, `SPEC-019`), Groß-/Kleinschreibung wird nicht gefaltet.
  `when.equals` ist eine Zeichenkette, begrenzt nur durch die Grenzen von
  `jsonb`, ohne eigene Längenvorgabe; die leere Zeichenkette ist zulässig. Die
  Bedingung trifft, wenn der Wert der Spalte im Bild als Zeichenkette
  zeichengenau gleich `when.equals` ist (Text-Stand der Quelle); es wird nichts
  abgeschnitten oder normalisiert.
- **Bildbasis.** Die Bedingung liest bei `INSERT` und `UPDATE` das Neu-Bild, bei
  `DELETE` das Alt-Bild, und zwar den Quellwert vor jeder Transformation
  (`SPEC-030`): eine `rename_column`- oder `map_value`-Regel der Tabelle ändert
  das Ergebnis der Bedingung nicht. Eine Backfill-Change (`INSERT`) wird
  wie ein `INSERT` des WAL-Pfads behandelt.
- **Abwesenheit.** Fehlt der Wert der Spalte im gelesenen Bild (`NULL`,
  unverändertes TOAST, bei `DELETE` ohne volle Replica-Identität jede
  Nicht-Schlüsselspalte, generierte Spalte), trifft die Bedingung nicht; die
  nächste Regel wird geprüft. Das ist keine Nichtanwendbarkeit. Gemessen
  (PostgreSQL 17 und 18, E2E-Messung): ein `DELETE` ohne volle Replica-Identität
  trifft für eine Nicht-Schlüsselspalte nicht, unter voller Replica-Identität
  trifft er (Alt-Bild nach [`LH-FA-DAT-005`](lastenheft.md)).
- **Auswertung.** Die Regeln der Tabelle werden in aufsteigender `order`
  geprüft; der erste Treffer bestimmt `route_target`, keine weitere Regel wird
  geprüft. Trifft keine Regel, ist `route_target` `NULL`. R2 macht die
  Reihenfolge eindeutig, R4 und R5 schließen Regeln aus, die nie wirken könnten
  (`SPEC-019`); gleicher Regelstand und gleiche Change ergeben dasselbe Ziel.
- **Anwendbarkeit** (die führende Stelle; `LH-FA-CFG-008.a` und `SPEC-008`
  verweisen hierher). Eine Regel ohne `when` ist auf jede Change anwendbar. Eine
  Regel mit `when` ist auf eine Change anwendbar, wenn `when.column` in der
  Relation der Change vorkommt. Im Run eines Backfills gilt dieselbe Definition
  gegen die Spalten des Snapshots. Anwendbarkeit hängt an Regelstand und
  Spaltenmenge, nie am Wert einer Zeile: kein Wert macht eine Regel
  unanwendbar. Eine nicht anwendbare Regel endet mit der Fehlerklasse `schema`
  (`SPEC-008`). Die Prüfung der Relation auf eine entfernte oder geänderte Spalte
  geht der Prüfung der Change voraus und ist regelunabhängig: eine Relation ohne
  eine bekannte Spalte endet als inkompatible Schemaänderung (`LH-FA-SCH-004.a`),
  mit oder ohne Regel auf diese Spalte.
- **Wirkung nur auf das Ziel.** Die Regeln ändern weder das Row Image noch ein
  anderes Feld der Change als `route_target`; die Auswertung einer Transformation
  ändert das Ziel nicht.

**Beispiele** (Tabelle `public.orders` mit den Spalten `id`, `name`, `region`;
Regel `eu_orders` mit `rule_spec`
`{"target": "eu", "order": 10, "when": {"column": "region", "equals": "eu"}}`
und Regel `rest` mit `rule_spec` `{"target": "sonstige", "order": 100}`):

- Eine Change mit dem Bild `{"id": "7", "name": "Ada", "region": "eu"}` trägt
  `route_target = eu`; mit `"region": "us"` trägt sie `sonstige`; ist `region`
  `NULL`, fehlt der Schlüssel im Bild, die Bedingung trifft nicht, und die Regel
  `rest` bestimmt `sonstige`. Ohne die Regel `rest` bliebe `route_target`
  `NULL`: die Change ist im Log und über jeden ungefilterten Weg sichtbar, aber
  unter keinem Ziel.
- Abgelehnter Antrag: eine weitere Regel mit `"order": 10` endet `failed` mit dem
  Fehlertext `abgelehnt [<Code>]: order bereits vergeben: public.orders.10`
  (`SPEC-019`, R2); der
  Regelstand bleibt unverändert.
- Nicht anwendbare Regel: Fehlt `region` in der Relation einer Change der
  Tabelle, ist `eu_orders` nicht anwendbar; die Change endet im Erfassungspfad
  mit der Fehlerklasse `schema`, bis die Regel entfernt ist
  (`LH-FA-CFG-008.a`, Abhilfe). Dieser Fall entsteht am laufenden System bei
  Erstaktivierung ohne bekannte Spaltenform und bei einer Publication mit
  Spaltenliste ohne `region` (gemessen, PostgreSQL 17 und 18, E2E-Messung);
  wird `region` aus der Tabelle entfernt, endet der Pfad zuerst als
  inkompatible Schemaänderung, und es entsteht keine Zeile.

### SPEC-033 — OTLP-Metrik-Export (Drahtform, Konfiguration, Verhalten)

Technische Ausgestaltung von [`LH-FA-SST-010`](lastenheft.md). Die SQL-Sicht
`cdc.metrics` ([`LH-FA-SST-004`](lastenheft.md), [`LH-QA-OPS-003`](lastenheft.md))
bleibt unverändert; der Export ist ein zweiter Zugriffsweg auf dieselben
Kennzahlen, ergänzt um `cdc_wal_retention_bytes`.

**Drahtform (bindend).** OTLP/HTTP mit Protobuf-Körper: `POST
<endpoint>/v1/metrics`, `Content-Type: application/x-protobuf`, Nachricht
`ExportMetricsServiceRequest`. `<endpoint>` ist die Basis-URL aus
`CDC_OTLP_ENDPOINT` ohne den Pfad `/v1/metrics`. Ein Status `2xx` gilt als
angenommen, jeder andere Status und jeder Verbindungs- oder Zeitfehler als
Fehlschlag. Welche Bibliothek die Nachricht erzeugt, ist nicht Teil des
Vertrags; die Drahtform ist es.

**Quelle der Werte.** Die Zeilen der Sicht `cdc.metrics`, gelesen über
`CDC_READER_DSN` (Rolle `cdc_reader`, keine zusätzlichen Rechte), und der
zuletzt im Prozess gemessene `cdc_wal_retention_bytes` (`SPEC-009`); ist noch
keiner gemessen, entfällt dieser Datenpunkt. Gleichheit mit der Sicht gilt je
Messzeitpunkt: `cdc_oldest_change_age_seconds` und `cdc_capture_lag` hängen an
`now()` der Abfrage und weichen um den Abstand der Messzeitpunkte ab;
`cdc_wal_retention_bytes` steht nicht in der Sicht.

**Kennzahlen.** Namen sind die der Sicht (`SPEC-009`). Alle sind OTLP-`Gauge`,
auch die mit dem Namensteil `_total`: ein Wert der Sicht ist ein
Momentstand und kein monoton wachsender Zähler (hergeleitet aus dem SQL-Text
der Sicht: die Zahl der Changes fällt bei einer Retention-Löschung, die Zahl der
Fehler zählt die Quellen im Fehlerzustand).

| Metrik | Einheit | Datenpunkt-Attribut |
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

`cdc_consumer_lag` ist der Abstand zwischen der letzten gespeicherten
Commit-Position der Quelle und der bestätigten Position des Consumers als
WAL-Strecke in Bytes (beide Positionen sind PostgreSQL-LSNs, `SPEC-003`), nicht
das Datenvolumen der Changes: das WAL trägt auch Änderungen nicht erfasster
Tabellen. `cdc_consumer_position` ist eine Position und trägt `1`.

Zahlenform: eine Kennzahl geht als ganze Zahl hinaus, wenn ihr Wert in der Sicht
eine ist und ihre Einheit nicht `s` ist; die zwei Kennzahlen mit Einheit `s`
(`cdc_oldest_change_age_seconds`, `cdc_capture_lag`) gehen immer als
Fließkommazahl. Eine Zeile der Sicht mit einem Namen außerhalb der Tabelle wird
nicht übertragen.

Resource-Attribute: `service.name` = `pg-change-feed`, `cdc.source_id` = die
Quellkennung des Prozesses.

**Konfiguration.**

| Variable | Bedeutung | Gültig |
|---|---|---|
| `CDC_OTLP_ENDPOINT` | Basis-URL des Empfängers; gesetzt aktiviert den Export, ungesetzt ist er aus | Schema `http` oder `https`, Host vorhanden |
| `CDC_OTLP_HEADERS` | zusätzliche Header je Anfrage, Form `k=v,k2=v2` | jedes Element mit einem gültigen Header-Namen als Schlüssel und einem `=` (Regeln unten) |
| `CDC_OTLP_INTERVAL_SECONDS` | Takt der Übertragung in Sekunden, Default 60 | ganze Zahl von 5 bis 3600 |

Die Konfigurationsdatei (`SPEC-016`) trägt nur den Takt (`otlp_interval`);
Endpunkt und Header sind env-exklusiv, weil sie Zugangsdaten tragen können. Ein
ungültiger Wert endet den Start mit der Fehlerklasse `configuration`.

- **Header.** Das Komma trennt die Elemente, das erste `=` Schlüssel und Wert
  (der Wert darf weitere `=` tragen). Der Schlüssel ist ein gültiger
  HTTP-Header-Name: Ziffern, ASCII-Buchstaben und die Zeichen
  ``!#$%&'*+-.^_`|~``; Leerraum im Schlüssel ist ein Fehler. Der Wert trägt kein
  Steuerzeichen (`U+0000` bis `U+001F` ohne Tabulator, `U+007F`); Leerraum im
  Wert bleibt erhalten. Eine leere Zeichenkette liefert keine Header. Der
  Fehlertext nennt die Position des Elements, nie dessen Inhalt.
- **Header ohne Endpunkt.** `CDC_OTLP_HEADERS` gesetzt bei nicht gesetztem
  `CDC_OTLP_ENDPOINT` endet den Start mit der Fehlerklasse `configuration`; ein
  Takt ohne Endpunkt ist zulässig. Eine leer gesetzte Variable gilt als nicht
  gesetzt.
- **Benutzerteil der Endpunkt-URL.** `http://benutzer:passwort@host` wird nicht
  abgelehnt und geht als `Authorization: Basic` hinaus; setzt
  `CDC_OTLP_HEADERS` selbst einen `Authorization`-Header, gilt dieser. Weder ein
  Log-Eintrag noch ein Fehlertext des Exports nennt die URL, den Benutzerteil
  oder einen Header-Wert.
- **Datei-Feld.** `otlp_interval` unterliegt derselben Prüfung (5 bis 3600) wie
  die Variable; die Variable geht vor.

**Verhalten.** Der Export läuft in einer eigenen Goroutine und blockiert weder
Erfassung, Persistierung noch Bestätigung. Der erste Versuch folgt einen vollen
Takt nach dem Start, nicht beim Start. Jeder Versuch hat eine Frist von
`min(Takt, 10 s)`, die für das Lesen der Sicht und die Übertragung zusammen
gilt. Es gibt keine Warteschlange: der nächste Takt ist die
Wiederholung, ein ausgefallener Takt wird nicht nachgeholt. Eine Weiterleitung
(`3xx`) wird nicht verfolgt und ist ein Fehlschlag, damit kein Header einen
anderen Host erreicht. Der erste
Fehlschlag seit dem Start oder seit der letzten Wiederaufnahme erzeugt eine Warnung, weitere
Fehlschläge höchstens alle 5 Minuten eine erneute; Lesen und Übertragen teilen
diesen Zustand, und die Wiederaufnahme ist ein Versuch, in dem beide gelungen
sind; sie erzeugt
eine Info-Zeile. Der Export beeinflusst weder den Health-Zustand noch
`error_class` des Heartbeats. Ein Fehlschlag der Übertragung und ein
Fehlschlag beim Lesen der Sicht tragen je einen eigenen Warn-Meldungscode im
Bereich 6 (`SPEC-008`); die Zuweisung der Codes liegt in der Code-Tabelle im
Quelltext.

### SPEC-034 — TLS der HTTP- und gRPC-Schnittstellen

Technische Ausgestaltung von [`LH-FA-SST-011`](lastenheft.md).

| Variable | Datei-Schlüssel (`SPEC-016`) | Bedeutung |
|---|---|---|
| `CDC_TLS_CERT_FILE` | `tls_cert_file` | Pfad der PEM-Datei mit Zertifikat und Kette |
| `CDC_TLS_KEY_FILE` | `tls_key_file` | Pfad der PEM-Datei mit dem privaten Schlüssel |

Beide Pfade sind Dateinamen, keine Zugangsdaten; der private Schlüssel selbst
steht nie in einer Umgebungsvariable oder Konfigurationsdatei. Ein Paar gilt
für HTTP (einschließlich `GET /changes/stream`, `SPEC-021`) und gRPC
(`SPEC-020`, `SPEC-031`) gemeinsam.

- Keines der beiden gesetzt: beide Schnittstellen laufen unverschlüsselt wie
  bisher.
- Nur eines gesetzt, ein Pfad nicht lesbar oder die beiden Dateien bilden kein
  Paar: Fehlerklasse `configuration`, kein Start.
- „Ladbar“ ist die ganze Prüfung beim Start: beide Dateien sind lesbares PEM
  und der Schlüssel gehört zum Zertifikat. Der Start prüft nicht den
  Gültigkeitszeitraum, nicht die Namen (Hostname, SAN) und nicht die Kette zu
  einer Wurzel; ein abgelaufenes oder für den Namen unpassendes Zertifikat
  startet den Server, und der Client lehnt es beim Verbinden ab.
- Beide gesetzt und ladbar: jede der Schnittstellen, die über ihre Adresse
  (`CDC_HTTP_ADDR`, `CDC_GRPC_ADDR`) aktiviert ist, wird nur über TLS bedient;
  die Mindestversion ist TLS 1.2, explizit gesetzt. Ein Klartext-Zugriff auf
  dieselbe Adresse wird nicht bedient.
- Kein Neuladen: Zertifikat und Schlüssel werden beim Start gelesen.
- Kein Client-Zertifikat wird verlangt oder geprüft.
- Die Verbindung zu NATS bleibt unberührt (Sache der NATS-Konfiguration); ein
  Fehler beim Binden der Adresse bleibt wie bisher ein Log-Eintrag.

### SPEC-035 — Mehrere API-Token je Klasse

Technische Ausgestaltung von [`LH-FA-SST-012`](lastenheft.md) für die Token
der HTTP- und gRPC-Schnittstellen (`SPEC-018`, `SPEC-031`).

| Variable | Klasse | Form |
|---|---|---|
| `CDC_API_TOKEN_READER`, `CDC_API_TOKEN_ADMIN` | lesend, administrativ | genau ein Token, unverändert (kein Trennzeichen, ein Komma gehört zum Token) |
| `CDC_API_TOKENS_READER`, `CDC_API_TOKENS_ADMIN` | lesend, administrativ | kommagetrennte Liste |

Die wirksame Menge einer Klasse ist die Vereinigung des Singulars mit der
Liste. Alle vier Variablen sind env-exklusiv (`SPEC-016`).

- Ein leeres Element der Liste oder ein Element mit Leerraum (Whitespace) ist
  ein Fehler der Klasse `configuration`, kein Start.
- Eine leer gesetzte Variable gilt als ungesetzt, in der Liste wie im Singular.
- Ein Komma im Token ist in der Liste nicht ausdrückbar (das Komma trennt); ein
  Token mit Komma steht im Singular.
- Die Prüfung auf Leerraum gilt für jedes Element der zwei Listen-Variablen;
  der Singular wird nicht auf Leerraum geprüft.
- Ein leeres Token ist nie gültig, in keiner Klasse.
- Steht derselbe Wert in beiden Klassen, gilt `admin`.
- Die Prüfung vergleicht den vorgelegten Wert zeitkonstant gegen alle Token
  beider Klassen, ohne beim ersten Treffer abzubrechen (Erwartung an die
  Umsetzung).
- HTTP und gRPC bilden die Klasse eines Tokens aus derselben Menge und
  gleich.
- Ein Wechsel besteht aus zwei Neustarts: das neue Token ergänzen, die Clients
  umstellen, das alte Token entfernen. Ein Neuladen zur Laufzeit gibt es nicht.

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
| `SPEC-025` | `CDC_BENCH_THRESHOLDS` | Zusätzliche Commit-Latenz je Quelltransaktion ≤ 0,10 ms in `SPEC-036` (Vorschlag) · Batch-Vorteil ≥ 10× | Pass/Fail für [`LH-QA-PER-001`](lastenheft.md)/[`LH-QA-PER-003`](lastenheft.md); [`LH-QA-PER-002`](lastenheft.md) nutzt `SPEC-013`; Einzelheiten unter der Tabelle; über ADR schärfbar |
| `SPEC-036` | `BENCH_SOURCE_ENV` | PostgreSQL 18 · Feed co-located · N = 5000, Median von 5 · `fdatasync` ≤ 0,5 ms | Messumgebung der Obergrenze `SPEC-025`; Einzelheiten unter der Tabelle; über ADR schärfbar |

**Zu `SPEC-025`.** Die Latenz ist die Differenz der Gesamtdauern von N einzeln
committeten Quelltransaktionen mit und ohne Feed, geteilt durch N. Der Wert
0,10 ms ist ein Vorschlag zur Bestätigung durch die Messung in der Umgebung
`SPEC-036` (abgeleitet aus der Obergrenze früherer Messungen, nicht erprobt).
Außerhalb der Umgebung druckt die Messung ihren Wert samt Verhältnis zur
gemessenen Latenz des Festschreibens (`fdatasync`) ohne Verdikt.

**Zu `SPEC-036`.** PostgreSQL 18 im Digest-Pin des Bench-Skripts; der
Feed-Container läuft auf demselben Host wie die Quelle; N = 5000 einzeln
committete `INSERT`s, Median von 5 Läufen je Phase; die Latenz einer
`fdatasync`-Operation auf dem Datenträger des Datenverzeichnisses liegt bei
höchstens 0,5 ms (gemessen mit `pg_test_fsync`). Die Latenz je Commit hängt am
Festschreiben auf dem Datenträger, deshalb gilt die absolute Zahl nur in dieser
Umgebung.

---

## 4. Fehler-Codes und Logging-Felder

Regeln dieser Sektion: Der Fehler-Code ist ein Laufzeit-Symbol, die
`SPEC-<NNN>` benennt die *Festlegung* darüber — beide stehen nebeneinander.

Fehler werden mindestens in die folgenden Klassen klassifiziert
(`SPEC-008`):

| ID | Klasse | Bedingung | Aktion |
|---|---|---|---|
| `SPEC-008` | `transient` | vorübergehend nicht verfügbare Quelle/Speicher | Erneut versuchen mit begrenztem Backoff — Bedingung: wiederholt wird die Transport-/Verbindungsstörung am Quellzugriff (`receive.ErrReplication`/`outbound.ErrReplication`), am Server ausgewählt nach SQLSTATE: Fehler ohne SQLSTATE, die Klassen 08, 40, 53, 55, 57, 58 und der Code 25006; nicht wiederholt enden die Klassen `permission`/`configuration`/`schema`/`storage`, Stream-Ordnungs-Verletzungen und Server-Abweisungen außerhalb dieser Auswahl (`receive.ErrRejected`) mit sichtbarem Fehler |
| `SPEC-008` | `configuration` | ungültige/falsch gesetzte Konfiguration | Sichtbarer Fehler; kein Start und keine Fortsetzung im falschen Stand |
| `SPEC-008` | `permission` | fehlende Berechtigung | Sichtbarer Fehler; kein stiller Retry |
| `SPEC-008` | `schema` | nicht sicher interpretierbare Schemaänderung/Dekodierfehler; eine Transformations- oder Routing-Regel, die auf die Spalten der Change bzw. des Snapshots nicht anwendbar ist (`SPEC-030`, `SPEC-032`) — im Erfassungspfad und im Run | Sichtbarer Fehler (LH-FA-SCH-004.a); kein stilles Überspringen. Erfassungspfad: keine Persistierung und keine Bestätigung der Transaktion, Erfassung endet sichtbar (Heartbeat, Diagnose); Run: `failed`, run-lokal. Abhilfe nur bei nicht anwendbarer Regel: Regelstand ändern (Absatz unten) |
| `SPEC-008` | `storage` | Persistenzfehler im ChangeStore | **Kein Source-ACK** (LH-QA-REL-001.a) |
| `SPEC-008` | `replication` | Replication-Stream/Slot-Störung: **Stream-Ordnungsverletzung** (BEGIN/COMMIT/Change außerhalb der erwarteten Reihenfolge) oder **Transport-/Verbindungsstörung** (Verbindungsaufbau, Start, Keepalive, Quell-Bestätigung) | Stream-Ordnungsverletzung: sichtbarer Fehler, harter Abbruch, keine Fortsetzung im widersprüchlichen Stand; Transport-/Verbindungsstörung: Überwachung über Schwellen (§5, WAL-Rückstand); kontrollierte Fortsetzung |
| `SPEC-008` | `internal` | unerwarteter interner Fehler | Sichtbarer Fehler; Restart-Strategie nach [`LH-QA-REL-002`](lastenheft.md) |

**Meldungscode.** Jede klassifizierte Fehlerursache, jede Ablehnung einer
Aufrufer-Eingabe und jede Warnung mit einer Maßnahme für den Betreiber trägt
einen Meldungscode der Form `PCF-<S><NNNN>` (ERE
`PCF-[EWI][0-9]{4}`); `S` ist die Schwere (`E` Fehler, `W` Warnung, `I` ist
reserviert). Bei Fehlern ist die erste Ziffer die Klasse der Tabelle oben
(1 `transient`, 2 `configuration`, 3 `permission`, 4 `schema`, 5 `storage`,
6 `replication`, 7 `internal`); die Ziffer 8 kennzeichnet die Ablehnung einer
Aufrufer-Eingabe, etwa eines Antrags, und trägt keine Klasse. Die Endung `000`
ist der Rückfall der Klasse: ein klassifizierter Fehler ohne Einzelursache
trägt den Rückfall seiner Klasse, nie keinen Code. Die Klasse folgt aus dem
Code; trägt die Kette eines Fehlers Codes mehrerer Klassen, gilt die
Vorrangfolge `transient`, `configuration`, `schema`, `permission`,
`replication`, `storage`, `internal` und der erste Code der gewählten Klasse in
der Kette; ein Fehler ohne Code ist `internal` mit dem Rückfall `PCF-E7000`.
`cdc.process_heartbeat.error_class` und das Metrik-Label `class` bleiben die
Klasse; `cdc.process_heartbeat.error_code` trägt den Code des Fehlerzustands
(`NULL` im Normalbetrieb), `GET /diagnose` und der RPC `Diagnose` das Feld
`error_code` (`SPEC-018`, `SPEC-031`), die Zeile „Fehlerzustand“ der
CLI-Diagnose nennt Klasse und Code (`schema [<code>]`); die übrigen Zeilen des
Berichts tragen keinen Code. Bei einer Warnung ist die erste Ziffer der
Bereich (1 Erfassung und Replikation, 2 Backfill, 3 Retention und Speicher,
4 Verwaltung, 5 Konfiguration und Start, 6 Beobachtbarkeit und Transport,
9 reserviert); sie trägt keine Klasse
und steht als eigenes Log-Attribut `code`, der Meldungstext trägt keinen Code.
Die Warn-Zeile des Heartbeat-Adapters, die einen Fehlerzustand meldet, trägt
unter `code` den `E`-Code des Fehlerzustands; ein Filter auf `code=PCF-W…`
erfasst sie nicht. Die Fehlerantwort der HTTP-API trägt den Code im Feld `code`
des Körpers (`SPEC-018`), ein Fehlerstatus der gRPC-API im Statusdetail
`ErrorInfo` (`SPEC-031`); die Warn-Zeile eines gescheiterten API-Aufrufs trägt
`PCF-W4008`.
Der Fehlertext (Fehlerwert, Zeile beim Prozessende, Attribut
`error` einer Log-Zeile) beginnt mit dem Kopf
`Fehlerklasse <klasse> [<code>]: <Ursache>`; der Prozessausgang bleibt unverändert,
kein Code ändert ihn.
`error_message` eines `failed`-Runs beginnt mit `<klasse> [<code>]: `, die eines
abgelehnten Antrags mit `abgelehnt [<code>]: `, eines an einem klassifizierten
Fehler gescheiterten Antrags mit dem Kopf des Fehlertexts. Ein Code wird nie
neu belegt; ein entfallener Code wird zurückgezogen und bleibt vergeben; die
Klasse eines Codes ändert sich nie. Stabil sind Code, Klasse, das Wort
`Fehlerklasse` am Anfang des Kopfes und der Prozessausgang; der Text nach dem Kopf
ist nicht Vertrag. Die Code-Tabelle liegt im Quelltext, der Katalog (Code,
Klasse, Bedeutung, Maßnahme) im Benutzerhandbuch.

**Nicht anwendbare Regel (Klasse `schema`).** Die Ursache definiert `SPEC-030`
(Anwendbarkeit) für eine Transformationsregel und `SPEC-032` (Anwendbarkeit)
für eine Routing-Regel; sie trägt in beiden Pfaden dieselbe Klasse. Das
Verhalten des Erfassungspfads und die Abhilfe führt `LH-FA-CFG-007.a` bzw.
`LH-FA-CFG-008.a` (Nicht anwendbare Regel, Abhilfe). Im Run einer Routing-Regel
prüft dieselbe Funktion wie im Erfassungspfad einmal je Run, vor der
Schreibtransaktion und der ersten Zeile; der Fehlertext beginnt mit Klasse und
Meldungscode (`schema [<code>]: `) und nennt Regelname und Spalte, keine Change
entsteht. Das ist eine
Zusage an die Umsetzung. Im Run steht nur der Run — `cdc.backfill_status` und die
CLI-Diagnose zeigen ihn `failed` mit der Klasse im Fehlertext, der
Erfassungspfad läuft weiter. Die Abhilfe ist die Regelstand-Änderung der
Abhilfe-Zusage; ein Prozessneustart für den neuen Run ist nicht Teil der
Zusage. Nach der Abhilfe im Run beginnt ein **neuer** Antrag
`cdc.backfill_table` einen neuen Run, ein `failed`-Run wird nicht fortgesetzt.
Die Abhilfe im Erfassungspfad gilt nur, wenn er an der Nichtanwendbarkeit der
Regel endete; endete er an einer inkompatiblen Schemaänderung (eine entfernte
Spalte), ist die Abhilfe die dieser Ursache (`LH-FA-SCH-003`, `LH-FA-SCH-004.a`),
und das Entfernen der Regel genügt dort nicht.

Keine Credentials in Logs.

---

## 5. Metriken und Tracing-Felder

Regeln dieser Sektion: verbindliche Felder der Betriebsschnittstelle
(Zugriffswege: die SQL-Sicht `cdc.metrics` und der OTLP-Export, `SPEC-033`).
Inhalt deckt [`LH-QA-OPS-003`](lastenheft.md) ab.

| ID | Metrik | Inhalt | Quelle |
|---|---|---|---|
| `SPEC-009` | `cdc_changes_processed` | verarbeitete Changes | ChangeStore |
| `SPEC-009` | `cdc_changes_pending` | ausstehende/unbestätigte Changes | ChangeStore |
| `SPEC-009` | `cdc_transactions_total` | persistierte Transaktionen | ChangeStore |
| `SPEC-009` | `cdc_errors_total{class}` | Fehler je Klasse (§4) | Application |
| `SPEC-009` | `cdc_capture_lag` | Abstand Quelländerung → CDC-Verfügbarkeit ([`LH-FA-ADM-004`](lastenheft.md)) | Capture |
| `SPEC-009` | `cdc_consumer_lag` | Rückstand je Consumer als WAL-Strecke in Bytes zwischen letzter Commit-Position und bestätigter Position ([`LH-FA-ADM-005`](lastenheft.md)) | ConsumerState |
| `SPEC-009` | `cdc_consumer_position` | bestätigte Position je Consumer, roh ([`LH-QA-OPS-003`](lastenheft.md)) | ConsumerState |
| `SPEC-009` | `cdc_wal_retention_bytes` | WAL-Rückstand/Replication-Slot-Zustand: das vom Feed noch nicht bestätigte WAL (aktuelles WAL-Ende der Instanz minus `confirmed_flush_lsn` des Slots); WAL ohne Inhalt für die Publication bestätigt der Feed im Leerlauf, es hält den Rückstand nicht | Replication Stream |
| `SPEC-009` | `cdc_storage_bytes` | Speicherverbrauch der CDC-Daten | ChangeStore |
| `SPEC-009` | `cdc_oldest_change_age_seconds` | Alter des ältesten aufbewahrten Changes | ChangeStore |

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
| `SPEC-020` | gRPC Server-Streaming (HTTP/2 mit Protobuf) | gRPC-Go `google.golang.org/grpc`, Protobuf-Runtime `google.golang.org/protobuf` | — (Vertrag steht in diesem Dokument, §2 SPEC-020) |
| `SPEC-023` | Beispiel-Client-Werkzeugketten (Go, C#/.NET, Kotlin/JVM) | digest-gepinnte Basis-Images, auf feste Versionen gepinnte Abhängigkeiten (Pin-Hebung = bewusster Commit) | — (Vertrag steht in diesem Dokument, §2 SPEC-023; die Werkzeugketten-Dateien liegen im jeweiligen Sprach-Wurzelverzeichnis) |
| `SPEC-024` | NATS Core (Vollinhalts-Stream, kein JetStream) | NATS-Server 2.x, Go-Client `github.com/nats-io/nats.go` (bereits im Baum, `SPEC-017`) | — (Vertrag steht in diesem Dokument, §2 SPEC-024) |
| `SPEC-026` | `PgChangeFeed.Client` NuGet-Package (C#/.NET, erstes SDK-Package für `LH-FA-SST-009`) | SemVer 2.0, `0.x.y` (aktuell `0.2.1`) | `sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.csproj` als Metadaten-Quelle (`<PackageId>`/`<Version>`) — kein eigener §2-Eintrag, das Package deckt bereits dokumentierte Drahtverträge (`SPEC-018`, `SPEC-020`, `SPEC-021`, `SPEC-022`, `SPEC-024`) |
| `SPEC-027` | `pgchangefeed` PyPI-Package (Python, zweites SDK-Package für `LH-FA-SST-009`) | PEP 440, `0.x.y` (aktuell `0.2.1`) | `sdks/python/pgchangefeed/pyproject.toml` als Metadaten-Quelle (`[project] name`/`version`) — kein eigener §2-Eintrag, das Package deckt bereits dokumentierte Drahtverträge (`SPEC-018`, `SPEC-020`, `SPEC-021`, `SPEC-022`, `SPEC-024`) |
| `SPEC-028` | `pgchangefeed-kotlin` Gradle-/Maven-Package auf Cloudsmith und GitHub Packages (Kotlin, drittes SDK-Package für `LH-FA-SST-009`) | SemVer 2.0, `0.x.y` (aktuell `0.2.2`) | `sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts` als Metadaten-Quelle (`group`/`version`, ein Repository je Vertriebsziel) — kein eigener §2-Eintrag, das Package deckt bereits dokumentierte Drahtverträge (`SPEC-018`, `SPEC-020`, `SPEC-021`, `SPEC-022`, `SPEC-024`) |
| `SPEC-033` | OTLP/HTTP (Protobuf), Metrik-Signal | OpenTelemetry-Protokoll, `POST /v1/metrics`; Typen aus `go.opentelemetry.io/proto/otlp` bevorzugt (nach Kenntnisstand, ungeprüft) | — (Vertrag steht in diesem Dokument, §2 SPEC-033) |

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
| 2026-09-14 | SPEC-020 ergänzt: gRPC-Live-Change-Stream — Dienst `ChangeStream` mit Server-Streaming-RPC `StreamChanges`, Protobuf-Nachrichtenschema der Change-Nachricht, Fire-and-Forget-Zustellsemantik ohne Stream-internes Replay, Authentifizierung über den Metadata-Eintrag `authorization` (`Bearer`-Form, dieselben Token-Klassen wie SPEC-018), Aktivierung über `CDC_GRPC_ADDR`; externe-Verträge-Zeile in §6 |
| 2026-09-15 | `SPEC-021` ergänzt: HTTP-Server-Sent-Events für den Live-Change-Stream — Endpunkt `GET /changes/stream` (`text/event-stream`), Event-Typ `change`, JSON-Nachrichtenschema mit denselben zehn Change-Feldern, kein Stream-internes Replay (der `Last-Event-ID`-Header bleibt ungenutzt), Fire-and-Forget-Zustellsemantik, Aktivierung über `CDC_HTTP_ADDR` samt `503`-Pfad ohne verdrahteten `Broadcaster` (aus `SPEC-018` herausgelöst — jener Abschnitt gilt den neun Port-gedeckten Fähigkeiten aus `LH-FA-SST-006`) |
| 2026-09-15 | SPEC-019 Fließtext ergänzt: `applied` heißt für `exclude_column`/`include_column` dauerhaft vermerkt — die `applied`-Zeilen sind die einzige Herkunft des Ausschlussstandes einer Tabelle, abgeleitet in `requested_at`-Ordnung mit `administration_request_id` als Zweitschlüssel, mitgeführt bei jedem Anlegen einer Erfassungs-Bindung; der Antrag gegen eine Tabelle ohne laufende Bindung endet `applied` statt `failed` |
| 2026-09-15 | `SPEC-022` ergänzt: HTTP-API `GET /changes` — Query-Parameter `source` (Pflicht), `schema`/`table` (optional, unabhängig), `from`/`to` (`commit_position` ≥ 1, Start inklusiv/Ende exklusiv), `limit` (optional, kein Default-Limit), JSON-Antwortform der Changes samt Klartext-Identität der Tabelle, deterministische Reihenfolge, leere Liste statt `404`, `400` für Parameter außerhalb der Liste; `SPEC-018`s Abgrenzungssatz trägt das Changes-Lesen nicht mehr als außerhalb |
| 2026-09-15 | `SPEC-022` Antwort-Zelle präzisiert: `committed_at` trägt RFC 3339 in UTC mit Bruchteil-Sekunden bis zu neun Stellen — abschließende Nullen im Bruchteil entfallen |
| 2026-09-17 | `SPEC-016` nachgezogen: Feldmenge um `http_addr`/`grpc_addr` erweitert (additiv, Env schlägt feldweise), Ausschlussklausel von der Drei-Schlüssel-Liste auf die **Klasse der zugangsdaten-tragenden Felder** gezogen (drei DSN-Schlüssel, zwei Token-Schlüssel, `nats_url`) samt eigener, den Grund nennender Fehlerzeile, und klargestellt, dass die env-exklusiven Variablen unter geladener Datei aus der Umgebung wirken |
| 2026-09-17 | `SPEC-023` ergänzt: Beispiel-Clients — Klasse (Vorbild, kein Belegträger, keine Zustandsmaschine), Verhältnis zum Draht, Ort und Form samt Sprach-Wurzel, Import-Grenze, Docker-only-Startform, Bau-/Testziel je Sprache als Werkzeug, kein Lesen der Konfigurationsdatei, Handbuch-Bindung, Sprachen und Umfang (Go: vier Oberflächen; C#/Kotlin: HTTP-Familie), kein Eintrag in der E2E-Abdeckung; externe-Verträge-Zeile in §6 |
| 2026-09-17 | `SPEC-023` Zeile *Sprachen und Umfang* auf die volle Matrix gezogen (vier Zugriffs-Oberflächen in Go, C# und Kotlin statt der HTTP-Familie in C#/Kotlin), samt dem benannten Zusatzkontext für einen fremdsprachigen gRPC-Bau und dem Ort der erzeugten Stubs (im Bau, nicht committet) |
| 2026-09-18 | `SPEC-024` ergänzt: NATS-Vollinhalts-Stream — Subjekt-Schema (`cdc.stream.<source_id>.<schema>.<table>`, eigener Namensraum neben `SPEC-017`s Wecksignal-Subjekt), JSON-Nachrichtenschema identisch zu `SPEC-021`, Granularität je Zeilen-Change, Fire-and-Forget ohne Replay, Authentifizierung über einen serverweiten NATS-Verbindungs-Token (`CDC_NATS_STREAM_TOKEN`), Aktivierung nur bei gesetztem `CDC_NATS_URL` **und** `CDC_NATS_STREAM_TOKEN`; externe-Verträge-Zeile in §6 |
| 2026-09-18 | `SPEC-016` nachgezogen: `nats_stream_token` in die Klasse der zugangsdaten-tragenden Felder aufgenommen (jetzt drei Token-Schlüssel statt zwei) und in die Liste der env-exklusiven, auch unter geladener Datei aus der Umgebung wirkenden Variablen ergänzt |
| 2026-09-18 | `SPEC-023` Zeile *Sprachen und Umfang* auf fünf Zugriffs-Oberflächen gezogen (`SPEC-024` ergänzt), samt der Klarstellung, dass eine Sprache ohne öffentliche JSON-Dekodierung für den Vollinhalts-Stream-Client zusätzlich eine gepinnte JSON-Bibliothek braucht |
| 2026-09-19 | `LH-FA-CAP-009.a` ergänzt: Backfill-Mechanismus (Export-Snapshot vs. Bulk-Copy, Markierung, Verhältnis zum WAL-Erfassungspfad) als offene, ADR-pflichtige technische Frage zur neuen Lastenheft-Anforderung `LH-FA-CAP-009` |
| 2026-09-19 | `LH-FA-CFG-007.a` ergänzt: Konfigurationsmechanismus, Ausdrucksform und Regel-Auswertungsreihenfolge von Transformationen/Routing als offene, ADR-pflichtige technische Frage zur neuen Lastenheft-Anforderung `LH-FA-CFG-007` |
| 2026-09-23 | `LH-FA-CFG-007.a` auf Transformationen begrenzt; `LH-FA-CFG-008.a` ergänzt: Routingform (Zielmodell, Konfigurationsmechanismus, Auflösung) als offene, ADR-pflichtige technische Frage zur aus `LH-FA-CFG-007` herausgelösten Lastenheft-Anforderung `LH-FA-CFG-008` |
| 2026-09-19 | `LH-FA-SST-009.a` ergänzt: Sprachmatrix und Paket-Vertriebsweg für offizielle Client-Bibliotheken als offene, ADR-pflichtige technische Frage zur neuen Lastenheft-Anforderung `LH-FA-SST-009` — die bestehenden Beispiel-Client-Werkzeugketten (`SPEC-023`) erfüllen sie nicht |
| 2026-09-19 | `SPEC-026` ergänzt: `PgChangeFeed.Client` NuGet-Package (C#/.NET, erstes SDK-Package für `LH-FA-SST-009`) — System, SemVer 2.0 `0.x.y`, Vertrag-Datei-Verweis auf die `.csproj` als Metadaten-Quelle; externe-Verträge-Zeile in §6 |
| 2026-09-19 | `LH-FA-SST-009.a` nachgezogen: Für C#/NuGet ist die Sprachmatrix-/Vertriebsweg-Frage beantwortet und `PgChangeFeed.Client` real paketierbar (`make sdk-pack-csharp`, `PgChangeFeed.Client.0.1.0.nupkg`) — die Kennung bleibt bestehen, eine zweite Sprache oder ein zweiter Vertriebsweg bleibt offen |
| 2026-09-19 | `SPEC-027` ergänzt: `pgchangefeed` PyPI-Package (Python, zweites SDK-Package für `LH-FA-SST-009`) — System, PEP 440 `0.x.y`, Vertrag-Datei-Verweis auf die `pyproject.toml` als Metadaten-Quelle; externe-Verträge-Zeile in §6 |
| 2026-09-19 | `LH-FA-SST-009.a` nachgezogen: Für Python/PyPI ist die Sprachmatrix-/Vertriebsweg-Frage beantwortet und `pgchangefeed` real paketierbar (`make sdk-pack-python`, `pgchangefeed-0.1.0-py3-none-any.whl`, `pgchangefeed-0.1.0.tar.gz`) — die Kennung bleibt bestehen, eine dritte Sprache oder ein dritter Vertriebsweg bleibt offen |
| 2026-09-20 | `SPEC-028` ergänzt: `pgchangefeed-kotlin` GitHub-Packages-Gradle-/Maven-Package (Kotlin, drittes SDK-Package für `LH-FA-SST-009`) — System, SemVer 2.0 `0.x.y`, Vertrag-Datei-Verweis auf die `build.gradle.kts` als Metadaten-Quelle; externe-Verträge-Zeile in §6 |
| 2026-09-20 | `LH-FA-SST-009.a` nachgezogen: Für Kotlin/GitHub Packages ist die Sprachmatrix-/Vertriebsweg-Frage beantwortet und `pgchangefeed-kotlin` real paketierbar (`make sdk-pack-kotlin`, `pgchangefeed-kotlin-0.1.0.jar`) — die Kennung bleibt bestehen, eine vierte Sprache oder ein vierter Vertriebsweg bleibt offen |
| 2026-09-22 | `LH-FA-SST-009.a`/`SPEC-026` nachgezogen (`welle-sdk-csharp-vollabdeckung`, slice-sdk-csharp-sse-client-flaeche + slice-sdk-csharp-nats-stream-client-flaeche): `PgChangeFeed.Client` deckt jetzt HTTP-API, gRPC-Stream, SSE und NATS-Vollinhalt statt nur HTTP-API und gRPC-Stream — Version auf `0.2.0` gehoben, real paketiert (`PgChangeFeed.Client.0.2.0.nupkg`); die Kennung bleibt bestehen, eine zweite Sprache oder ein zweiter Vertriebsweg bleibt offen |
| 2026-09-22 | `LH-FA-SST-009.a`/`SPEC-028` nachgezogen (`welle-sdk-kotlin-vollabdeckung`, slice-sdk-kotlin-sse-client-flaeche + slice-sdk-kotlin-nats-stream-client-flaeche): `pgchangefeed-kotlin` deckt jetzt HTTP-API, gRPC-Stream, SSE und NATS-Vollinhalt statt nur HTTP-API und gRPC-Stream — Version auf `0.2.0` gehoben, real paketiert (`pgchangefeed-kotlin-0.2.0.jar`); die Kennung bleibt bestehen, eine vierte Sprache oder ein vierter Vertriebsweg bleibt offen |
| 2026-09-23 | `LH-FA-SST-009.a`/`SPEC-027` nachgezogen (`welle-sdk-python-vollabdeckung`, slice-sdk-python-grpc-client-flaeche + slice-sdk-python-sse-client-flaeche + slice-sdk-python-nats-stream-client-flaeche): `pgchangefeed` deckt jetzt HTTP-API, gRPC-Stream, SSE und NATS-Vollinhalt statt nur HTTP-API — Version auf `0.2.0` gehoben, real paketiert (`pgchangefeed-0.2.0-py3-none-any.whl`, `pgchangefeed-0.2.0.tar.gz`); die Kennung bleibt bestehen, eine dritte Sprache oder ein dritter Vertriebsweg bleibt offen |
| 2026-09-24 | `LH-FA-CAP-009.a` beantwortet: Backfill-Mechanismus als Zusagen an die Umsetzung — ausdrückliche Auslösung über die Antragsart `backfill`, Annahme in einer Transaktion und Aufnahme durch einen einzelnen Worker, Bulk-Copy im Snapshot eines je Run angelegten temporären Slots in einer Store-Transaktion, Markierung über `origin`, Position `X` je Run, Überlappungs-Verhalten (keine Lücke, begrenzte idempotente Dopplung), Sichtbarkeits-Grenze, Neubeginn nach Abbruch, Fail-closed-Prüfung, zwei Warnungen für große Tabellen; die Überschrift verliert „offen" |
| 2026-09-24 | `SPEC-001` um `cdc.backfill_run` erweitert; `SPEC-002` um das Feld `origin` (`wal` \| `backfill`, fehlender Wert liest als `wal`, letzte Spalte der View `cdc.changes`); `SPEC-019` um die Antragsart `backfill` (fünf Werte) und die Bedeutung von `applied` bei `backfill` („angenommen"); `SPEC-022` um das Antwort-Feld `origin` und die Position-und-`limit`-Anmerkung; `SPEC-020`/`SPEC-021` grenzen das Feld `origin` aus der Nachricht aus |
| 2026-09-24 | `SPEC-029` ergänzt: Feldform von `cdc.backfill_run` und `cdc.backfill_status` — Spalten, zwei Warn-Spalten (`warn_estimated_size`, `warn_duration`), `estimated_rows` NULL als „unbekannt", Grants je Rolle |
| 2026-09-24 | `SPEC-019` um den Absatz „Grants" der Antrags-Queue erweitert (`cdc_admin` `SELECT`, `UPDATE`; `cdc_capture` und `cdc_reader` kein Recht; kein `INSERT`, kein `DELETE`); `LH-FA-CAP-009.a` Absatz „Markierung" um die Bedeutung der Schema-Version einer Backfill-Change ergänzt (Kennung zum Run-Start, keine Beschreibung der Bild-Spalten) |
| 2026-09-24 | `LH-FA-CAP-009.a` Absatz „Mechanismus" um die Lesesperre des Imports und das Ende des Runs bei umgeschriebener Tabelle ergänzt (`failed`, Fehlerklasse `transient`, ohne Change) |
| 2026-09-25 | `LH-QA-REL-001.a` Schritt 5 und Invariante um die Bestätigung im Leerlauf ergänzt (WAL-Ende der Keepalive-Nachricht, ohne Persistenz, nie inmitten einer Quelltransaktion, nie zurück); `SPEC-009` Zeile `cdc_wal_retention_bytes`: Bedeutung „vom Feed noch nicht bestätigtes WAL" |
| 2026-09-25 | `SPEC-026`/`SPEC-027`/`SPEC-028` nachgezogen: die HTTP-Lesemodelle der drei Packages tragen das `GET /changes`-Antwortfeld `origin` (`wal` \| `backfill`, ein fehlendes Feld oder JSON-`null` liest als `wal`, jeder andere Wert kommt unverändert an), die Live-Flächen (gRPC, SSE, NATS-Vollinhalt) tragen es nicht; die Zeilen nennen `SPEC-022` unter den gedeckten Drahtverträgen |
| 2026-09-25 | `LH-FA-SST-009.a` Eingabe um `SPEC-022` ergänzt (der Zustellweg `GET /changes` ist Teil der HTTP-Fläche der Packages) |
| 2026-09-25 | `SPEC-026`/`SPEC-027`/`SPEC-028` nachgezogen: die Version der drei Packages ist `0.2.1`; die Paketbeschreibung (README und Metadaten-Felder) ist Anwender-Dokumentation ohne interne Kennungen |
| 2026-09-25 | `LH-FA-RET-004.a` um Punkt 4 ergänzt: die Bereinigungsmenge wird seitenweise bestimmt (10.000 Kandidaten je Seite, ohne Row Images); der Arbeitsspeicher eines Bereinigungslaufs hängt an der Seitengröße und nicht mit nennenswertem Betrag an der Zahl der gespeicherten Changes |
| 2026-09-25 | `LH-FA-SST-009.a`/`SPEC-028` nachgezogen: `pgchangefeed-kotlin` hat zwei Vertriebsziele — Cloudsmith (anonym lesbar, ohne Konto und Token beziehbar) und GitHub Packages (Bezug mit Token); beide Ziele erhalten dieselben Artefakte, die Version ist `0.2.2`. Die Kennung bleibt bestehen, eine vierte Sprache oder ein weiterer Vertriebsweg bleibt offen |
| 2026-09-26 | `LH-FA-CFG-007.a` beantwortet: Transformationen als Zusagen an die Umsetzung — Konfiguration über zwei Antragsarten der Antrags-Queue, geschlossener Regelsatz (`rename_column`, `map_value`), Auswertung nach dem Spaltenausschluss in Relation-Spaltenreihenfolge, Mehrdeutigkeit statisch ausgeschlossen, Wirkort vor der Persistierung, nicht anwendbare Regel endet sichtbar mit Klasse `schema`, Abhilfe im gescheiterten Prozess; die Überschrift verliert „offen"; `LH-FA-CFG-008.a` bleibt offen |
| 2026-09-26 | `SPEC-019` um die Antragsarten `set_transformation` und `remove_transformation` (sieben Werte in `request_kind`), die Spalten `rule_name`/`rule_spec`, die zweite Bedeutung von `applied` für beide Arten, die Konfliktfreiheit K1–K4 und die `failed`-Fehlertexte erweitert |
| 2026-09-26 | `SPEC-030` ergänzt: Regelform (`rule_spec`) und Wirkung auf Row Images — Regeltypen, Abwesenheit, Vergleich, Position, Randfälle (leeres `values`, Zielname gleich Quellname, Abbildung auf sich selbst), Anwendbarkeit |
| 2026-09-26 | `SPEC-030` um den Bezeichner-Vergleich (`column`, `to`, Regelname), die Schlüsselreihenfolge als nicht zugesagt statt „Position", Beispiele je Regeltyp und die führende Stelle der Anwendbarkeit erweitert; `SPEC-019` Fehlertext-Tabelle um Regelname, fehlende `rule_spec` und die Prüfreihenfolge der Formzeilen erweitert; `SPEC-008` Zeile `schema` ordnet die Abhilfe der nicht anwendbaren Regel zu |
| 2026-09-26 | `SPEC-002` um die Aussage ergänzt, dass die Schlüsselmenge der Row Images dem Regelstand zum Erfassungszeitpunkt folgt; `SPEC-008` Zeile `schema` und der Absatz darunter um die Nichtanwendbarkeit einer Transformationsregel im Erfassungspfad und im Run samt Abhilfe erweitert; `LH-FA-CAP-009.a` um die Regelwirkung im Bild, den Regelstand in der Fail-closed-Prüfung und den Run-Fehler der Klasse `schema` ergänzt |
| 2026-09-26 | `LH-FA-CAP-009.a` Absatz „Markierung": das Backfill-Row-Image ist inhaltsgleich (Schlüsselmenge und Werte) dem WAL-Image derselben Zeile; `SPEC-008` Absatz „Nicht anwendbare Regel": ein Prozessneustart für den neuen Run ist nicht Teil der Zusage |
| 2026-09-26 | `SPEC-019` Absatz „Transformations-Antragsarten": `cdc.set_transformation` nimmt die Regelform als `json`-Parameter an (Spalte `rule_spec` bleibt `jsonb`); Aufrufform: Literal oder `::json`, ein `jsonb`-typisierter Wert braucht den Cast; syntaktisch ungültiges JSON scheitert beim Aufruf, gültiges JSON wird als Antrag angenommen und in Go geprüft |
| 2026-09-26 | `SPEC-019` Absatz „Transformations-Antragsarten": die Annahmemenge von `rule_spec` benannt statt „gültiges JSON" — angenommen `NULL` und jeder Text, den PostgreSQL als `json` liest und nach `jsonb` umwandelt; abgelehnt ohne Antrags-Zeile Syntaxfehler, `\u0000`-Escape, Zahl außerhalb des `numeric`-Bereichs, Schachtelung jenseits der Stapeltiefe |
| 2026-09-26 | `SPEC-019` Spalte `requested_at` und neuer Absatz „Ordnung der Verarbeitung“: `requested_at` ist der Zeitpunkt des Funktionsaufrufs, Aufrufe einer Transaktion werden in Aufrufreihenfolge verarbeitet, Verarbeitungs- und Ableitungsordnung sind dieselbe (Kennung als deterministischer Zweitschlüssel); Grenzen: Aufruf- statt Festschreibungs-Zeitpunkt, Serveruhr; K1 nennt, dass Entfernen und Neusetzen in einer Transaktion stehen dürfen |
| 2026-09-26 | `SPEC-019` Absatz „Transformations-Antragsarten“: die Prüfung von Regelname und Regelform liegt in der Verarbeitung, nicht im Antrags-Konstruktor — ein fehlender oder ungültiger Wert endet `failed` mit dem Fehlertext der Tabelle, die Queue lehnt den Antrag nicht beim Lesen ab |
| 2026-09-26 | `SPEC-019` neuer Absatz „Zeilen, die kein Antrag sind“: die Lesung der Queue lehnt keine Zeile ab, eine vom Antrags-Konstruktor verworfene Zeile (leere Quelle, leeres Schema, leerer Tabellenname, unbekannte Antragsart, leere Spalte der Spalten-Antragsarten) endet in der Verarbeitung `failed` mit Klartext und Antrags-Kennung als Adresse, die Zeilen dahinter laufen weiter; Grenze: eine Zeile ohne Kennung wird mit Warnung übersprungen; Zeilen `schema_name`/`table_name`/`column_name` und Absatz „Transformations-Antragsarten“ nachgezogen |
| 2026-09-28 | `SPEC-020` Zeile „Request“: `StreamChangesRequest` trägt zwei optionale, unabhängig setzbare Felder `schema`/`table` statt keiner Felder — Kombinatorik und Filterort ergänzt; `SPEC-021` neue Zeile „Query-Parameter“: `GET /changes/stream` nimmt dieselben zwei Felder als Query-Parameter, ein Parameter außerhalb der Menge endet `400`; `LH-FA-CFG-008.a` umformuliert — „gRPC und SSE liefern ungefiltert" trifft nicht mehr zu, die offene Frage bleibt auf das serverseitig konfigurierbare Routing-Zielmodell beschränkt |
| 2026-10-01 | `LH-FA-CFG-008.a` beantwortet: Routing als Zusagen an die Umsetzung — Zustellziel als benannter Kanal je Quelle, Konfiguration über zwei Antragsarten der Antrags-Queue, Regeln mit Ziel, Ordnungszahl und optionaler Gleichheitsbedingung, erster Treffer in aufsteigender Ordnungszahl, Mehrdeutigkeit statisch ausgeschlossen, Wirkort vor der Persistierung mit dem Ziel als Label an der Change, kein Treffer ergibt kein Ziel, nicht anwendbare Regel endet sichtbar mit Klasse `schema`, Abhilfe als Zusage (Erreichbarkeit offen); die Überschrift verliert „offen" |
| 2026-10-01 | `SPEC-019` um die Antragsarten `set_route` und `remove_route` (neun Werte in `request_kind`), die Bedeutung von `rule_name`/`rule_spec` für sie, die sechs Konfliktfreiheits-Invarianten R1 bis R6 mit Prüfreihenfolge und Fehlertexten, die zweite Bedeutung von `applied` für die Routing-Regelstand-Ableitung und die Sperre von `exclude_column` gegen eine Spalte mit Routing-Bedingung ergänzt |
| 2026-10-01 | `SPEC-032` ergänzt: Routing-Regel (`rule_spec`) und ihre Wirkung auf das Zustellziel — Schlüssel, Zielname, Bedingung, Bildbasis, Abwesenheit, Auswertung, Anwendbarkeit, Beispiele |
| 2026-10-01 | `SPEC-001`/`SPEC-002`: Spalte `route_target` an `cdc.change`, als letzte Spalte der View `cdc.changes`; `SPEC-020`, `SPEC-021`, `SPEC-022`, `SPEC-031` (Zeile `ReadChanges`, Feld `target = 7`): optionaler Filter `target` als Konjunktion mit `schema`/`table`; `SPEC-022`, `SPEC-031`: Fehler des Lese-Kontrakts haben bei ungültigem `target` Vorrang vor der leeren Antwort; `SPEC-024`: Zusatz-Subjekt `cdc.route.<source_id>.<ziel>`; `SPEC-008` Zeile `schema` und Absatz „Nicht anwendbare Regel" sowie `LH-FA-CAP-009.a` (Fehler des Runs, Ziel der Backfill-Changes) um die Routing-Regel ergänzt |
| 2026-10-01 | `LH-FA-CAP-009.a`: Fail-closed-Prüfung vor dem Commit um den Routing-Regelstand ergänzt (Stand des Runs, Klasse `configuration`), „Regelstand zum Run" definiert; `SPEC-020`, `SPEC-021`, `SPEC-022`, `SPEC-031`: ein `target` außerhalb des Alphabets, auch mit U+0000, liefert auf den Lesewegen eine leere Antwort, der SQL-Zugriff ist ausgenommen; `SPEC-032`: Wertebereich von `order` auf „positive ganze Zahl" zurückgenommen |
| 2026-10-01 | `LH-FA-CFG-008.a` Absatz „Abhilfe (Zusage)": Geltung auf den Fall der nicht anwendbaren Regel geschärft, die entfernte Spalte endet zuerst als inkompatible Schemaänderung (hergeleitet), Erreichbarkeit am System weiter offen; `SPEC-032` (Anwendbarkeit): Relation-Prüfung geht der Change-Prüfung voraus und ist regelunabhängig, Beispiel „Nicht anwendbare Regel" nachgezogen; `SPEC-008` Absatz „Nicht anwendbare Regel": Grenze der Abhilfe |
| 2026-10-01 | `SPEC-024` (Zusatz-Subjekt): Aussage zur Last der zweiten Veröffentlichung von „nicht gemessen" auf den gemessenen Umfang (Testcontainer, ohne Abonnent, ohne Schwelle) gezogen, ohne Zahl und ohne Last-Zusage |
| 2026-10-01 | `LH-FA-CFG-008.a`, `SPEC-032`: die Aussagen zu `DELETE` ohne volle Replica-Identität und zur Erreichbarkeit der Nichtanwendbarkeit von „nicht gemessen" auf den gemessenen Stand gezogen (PostgreSQL 17 und 18, E2E-Messung); `SPEC-019`: Adressspalte der Zeilen `Regelname ist ungültig` und `Zielname ist ungültig` der Routing-Tabelle auf die Adressform der Prosa angeglichen |
| 2026-10-03 | `SPEC-008`: Absatz „Meldungscode“ (Form `PCF-[EWI][0-9]{4}`, Klasse in der ersten Ziffer, Ziffer 8 für Ablehnungen, Rückfall `…000`, Kopf `Fehlerklasse <klasse> [<code>]: …`, Stabilität); der Satz zum Fehlertext-Beginn im Run (`schema: `) nennt Klasse und Code; `SPEC-029` (`error_message` des Runs) und `SPEC-019` (`error_message` des Antrags, Texte der Ablehnungen mit dem Kopf `abgelehnt [<Code>]: `) |
| 2026-10-03 | `SPEC-008`: Warncodes (Schwere `W`, erste Ziffer der Bereich, Log-Attribut `code`), Vorrangfolge der Klassen bei mehreren Codes in einer Kette und der Code des Fehlerzustands in Heartbeat und Diagnose; Tabelle von `cdc.process_heartbeat` (Spalte `error_code`, View `cdc.heartbeat`); `SPEC-018` (`GET /diagnose`, Feld `error_code`) und `SPEC-031` (`HeartbeatStatus.error_code`, Feldnummer 4) |
| 2026-10-03 | `SPEC-008`: Warn-Zeile des Fehlerzustands trägt den `E`-Code unter `code`; `SPEC-031`: `error_code` als eigener Aufzählungspunkt von `HeartbeatStatus` |
| 2026-10-03 | `SPEC-018` (Fehler-Antwortform, `SPEC-022` Zeile „Fehler-Antwortform“): additives Feld `code` im Fehlerkörper, Codes der Ablehnungen `PCF-E8050` bis `PCF-E8057`, `401`/`403` ohne Code; `SPEC-031`: Statusdetail `ErrorInfo` (`reason` Code, `domain` `pg-change-feed`), `Unauthenticated`/`PermissionDenied` ohne Detail; `SPEC-008`: Verweis auf die beiden Wege und `PCF-W4008` |
| 2026-10-04 | `SPEC-033` (OTLP-Metrik-Export), `SPEC-034` (TLS der HTTP- und gRPC-Schnittstellen), `SPEC-035` (mehrere API-Token je Klasse) neu; `SPEC-036` (Messumgebung `BENCH_SOURCE_ENV`) neu in §3; `SPEC-025` neu gefasst (absolute Obergrenze der zusätzlichen Commit-Latenz je Quelltransaktion als Vorschlag, Verhalten außerhalb der Messumgebung); `SPEC-016`: Schlüssel `tls_cert_file`, `tls_key_file`, `otlp_interval`, ausgeschlossene Klasse um `otlp_endpoint`, `otlp_headers` und die zwei Token-Listen erweitert, env-exklusive Liste ergänzt; `SPEC-008`: Warn-Bereich 6 „Beobachtbarkeit und Transport“; `SPEC-009`: `cdc_oldest_change_age` heißt `cdc_oldest_change_age_seconds` (Name der Sicht); §5 Kopfsatz auf die beiden Zugriffswege gezogen; §6: `SPEC-033` als externer Vertrag; `SPEC-018`: Verweise auf `SPEC-034`/`SPEC-035` |
| 2026-10-04 | `SPEC-033`: Einheit von `cdc_consumer_lag` ist `By` (WAL-Strecke), Zahlenform, Header-Regeln, Header ohne Endpunkt, Benutzerteil der Endpunkt-URL, erster Versuch nach einem vollen Takt, Frist für Lesen und Übertragen gemeinsam, gemeinsamer Fehlerzustand, keine Weiterleitung; `SPEC-034`: Prüfung beim Start auf das Laden des Paares begrenzt (Ablauf, Namen, Kette nicht geprüft); `SPEC-035`: leere Variable, Komma im Token, Leerraum nur in den Listen; `SPEC-009`: Bedeutung von `cdc_consumer_lag` — schreibt den Ist-Zustand fest, kein neues Verhalten außer der Einheit |
