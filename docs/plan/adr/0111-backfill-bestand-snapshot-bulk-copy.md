# ADR-0111: Backfill des Bestands — Bulk-Copy in einem Slot-Snapshot, atomar committet, als `origin = 'backfill'` markiert

**Status:** Accepted

**Datum:** 2026-09-23

**Autor:** pt9912 (Architect-Rolle, Modul 8)

**Bezug:** [`LH-FA-CAP-009`](../../../spec/lastenheft.md) (Haupt-Bezug —
Initial-Snapshot/Backfill des Bestands; Out-of-Scope: keine
Instanz-zu-Instanz-Migration, Parallelisierung/Durchsatz ist Ausbaustufe),
[`LH-FA-CAP-004`](../../../spec/lastenheft.md) (logische Reihenfolge),
[`LH-FA-CAP-008`](../../../spec/lastenheft.md) (Row-Image-Abwesenheit),
[`LH-FA-DAT-006`](../../../spec/lastenheft.md) (Metadaten-Erweiterbarkeit),
[`LH-FA-REA-001`](../../../spec/lastenheft.md),
[`LH-FA-REA-003`](../../../spec/lastenheft.md),
[`LH-FA-REA-004`](../../../spec/lastenheft.md) (Lesezugriffsweg),
[`LH-FA-CON-004`](../../../spec/lastenheft.md),
[`LH-FA-CON-005`](../../../spec/lastenheft.md) (Consumer-Position),
[`LH-FA-ADM-001`](../../../spec/lastenheft.md),
[`LH-FA-SST-003`](../../../spec/lastenheft.md) (Administration und Diagnose),
[`LH-QA-SEC-004`](../../../spec/lastenheft.md) (Spaltenausschluss),
[ADR-0004](0004-cdc-domainmodell.md), [ADR-0005](0005-sourceposition-abstrahiert-lsn.md),
[ADR-0011](0011-persist-before-ack.md), [ADR-0012](0012-at-least-once.md),
[ADR-0014](0014-retention-domain-policy.md), [ADR-0016](0016-jsonb-row-images-mvp.md),
[ADR-0017](0017-generische-change-tabelle.md), [ADR-0027](0027-capture-application-service.md),
[ADR-0028](0028-inbound-use-cases.md), [ADR-0034](0034-ports-nach-faehigkeiten.md),
[ADR-0040](0040-clockport.md), [ADR-0043](0043-schemamigrationen-mit-d-migrate.md),
[ADR-0047](0047-rollenspezifische-dsn-verdrahtung.md),
[ADR-0050](0050-sql-administration-antragsqueue-und-live-reload.md),
[ADR-0053](0053-retention-loeschausfuehrung-cdc-admin-delete-grant.md),
[ADR-0058](0058-testansatz-fuenf-luecken.md), [ADR-0059](0059-spaltenauswahl-mechanismus.md),
[ADR-0065](0065-spaltenausschluss-dauerhafter-traeger.md),
[ADR-0081](0081-changes-lesen-ueber-die-http-api.md)

**Schärft:** [`LH-FA-CAP-009.a`](../../../spec/pflichtenheft.md) (die
offene, ADR-pflichtige Frage „Backfill-Mechanismus offen" — diese ADR
beantwortet sie), [`SPEC-002`](../../../spec/pflichtenheft.md) (Change-Form
um `origin` erweitert), [`SPEC-019`](../../../spec/pflichtenheft.md)
(Antragsart `backfill`), [`SPEC-022`](../../../spec/pflichtenheft.md)
(`GET /changes` trägt `origin`), [`ARC-005`](../../../spec/architecture.md) /
[`ARC-006`](../../../spec/architecture.md) (neue Adapter-Rollen)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

[`LH-FA-CAP-009`](../../../spec/lastenheft.md) verlangt, den Bestand einer
aktivierten Tabelle zusätzlich zu den künftigen Änderungen über
**denselben Lesezugriffsweg** lesbar zu machen, **als Backfill erkennbar**
(Happy Path), mit **definiertem Verhalten im Überlappungsfenster** zum
laufenden WAL-Pfad — weder stille Lücke noch unbegrenzte Dopplung
(Boundary) — und **ohne Verlust bei Unterbrechung** (Negative).
[`LH-FA-CAP-009.a`](../../../spec/pflichtenheft.md) lässt Mechanismus,
Markierung und Überlappungsverhalten offen und macht sie ADR-pflichtig.

### Befunde am Bestand (jeweils mit Anker)

- **Der Haupt-Slot entsteht einmal und exportiert keinen Snapshot.**
  `ensureSlot` (`internal/adapters/driving/replication/receive/receive.go`)
  legt den Slot nur an, wenn er fehlt, mit `SnapshotAction:
  "NOEXPORT_SNAPSHOT"`; danach startet `NewStream` sofort
  `START_REPLICATION` auf derselben Verbindung, die zugleich die ACK-Rolle
  trägt (`ADR-0007`, Option C).
- **Tabellen werden zur Laufzeit aktiviert.** `applyAdministrationRequest`
  (`internal/bootstrap/wiring.go`) führt `enable` über `EnableTableUseCase`
  aus und trägt die `Assembler`-Bindung synchron nach (`AddBinding`); der
  Assembler verwirft Änderungen an Tabellen ohne Bindung
  (`mapper.Assembler.change`, `lookupBinding`). Ein Snapshot zum Zeitpunkt
  der Slot-Anlage kann eine später aktivierte Tabelle deshalb nicht
  abdecken.
- **Change-Modell und Position.** `model.Change` trägt `ID`, `TransactionID`,
  `SourceTableID`, `Sequence`, `Operation` (`INSERT|UPDATE|DELETE`, CHECK
  `chk_change_operation` in `tools/schema/schema.yaml`), die Row Images und
  `SchemaVersion`; die Position liegt an der Transaktion
  (`model.ChangeTransaction.Commit`, `cdc.transaction.commit_position`).
  Die WAL-Kennung ist `<xid>-<sequence>` (`Assembler.change`). Die Lese-
  Ordnung ist `(commit_position, transaction_id, sequence)`
  (`queries.SelectChanges`, View `cdc.changes`); die Position ist
  Bereichsgröße (`[Start, End)`), das `Limit` schneidet Zeilen, nicht
  Positionen — ein Position mit mehr Changes als `Limit` ist über
  `GET /changes` nicht in Teilen fortsetzbar (Lesen ab derselben Position
  liefert wieder dieselben Zeilen; `ADR-0081` nennt die Cursor-Form als
  eigenen Re-Evaluierungs-Trigger).
- **Row-Image-Form.** `mapper.rowImage` baut ein JSON-Objekt aus den
  gesendeten Spalten in Relation-Reihenfolge, Werte als JSON-Strings im
  Text-Stand der Quelle, `NULL` entfällt, ausgeschlossene Spalten entfallen
  (`ADR-0016`, `ADR-0059` Teilfrage 3).
- **Persistenz.** `PersistTransaction` schreibt eine im Speicher gehaltene
  `ChangeTransaction` in **einem** Store-Commit
  (`postgresstorage.PostgresChangeStoreAdapter`); `cdc_capture` trägt
  `INSERT` auf `cdc.transaction`/`cdc.change` und das Attribut
  `REPLICATION`, **kein** `SELECT` auf Quelltabellen
  (`tools/schema/nacharbeit-roles.sql`).
- **Retention.** `RunRetentionService.Run` liest **alle** Changes der Quelle
  (`ChangeQuery{Source}` ohne Limit) und gibt je Change frei, wenn Alter ≥
  `MinAge` und jede bestätigte Consumer-Position an oder hinter der
  Change-Position liegt (`RetentionPolicy.AllowsDeletion`); das produktive
  `MinAge` ist 24 h (`retentionMinAge` in `internal/bootstrap/wiring.go`),
  das Alter liest `committed_at` der Transaktion (`ADR-0040`).
- **Administration.** Schreibende SQL-Funktionen legen ausschließlich einen
  Antrag in `cdc.administration_request` ab und senden `pg_notify`; die
  Administrations-Goroutine verarbeitet Anträge **sequenziell** und
  synchron (`processAdministrationRequests`); die geschlossene
  `request_kind`-Menge trägt ein CHECK aus `nacharbeit-administration.sql`,
  weil d-migrate 1.3.1 eine CHECK-Änderung an einer bestehenden Tabelle
  nicht konvergiert (Kopfkommentar dort und in `schema.yaml`). Der
  Idempotenz-Guard des Schema-Rollouts kennt die sechs Fremdobjekte
  namentlich (`tools/schema/rolloutguard/guard.go`,
  `knownForeignObjects`).
- **Wire-Schemata.** Die Live-Wege (gRPC `SPEC-020`, SSE `SPEC-021`,
  NATS-Vollinhalt `SPEC-024`) speisen sich ausschließlich aus dem
  In-Prozess-`Broadcaster` (`ChangeStreamPort`) nach dem Source-ACK
  (`CaptureService.Capture`); `GET /changes` (`SPEC-022`) und die View
  `cdc.changes` lesen den Store. Die SDK-Decoder ignorieren unbekannte
  JSON-Felder bzw. lesen Felder explizit (C# `System.Text.Json`-Default,
  Kotlin Gson, Python `json`+`from_json`; gelesen an
  `sdks/*/…/{Http,sse,nats}`), ein zusätzliches Feld bricht sie nicht.

### Real geprüft (Wegwerf-Läufe, nicht committet)

Lauf vom 2026-09-23 in einem lokalen Wegwerf-Container
(`postgres:18-alpine`, PostgreSQL 18.6, `wal_level=logical`; die
Export-Syntax zusätzlich gegen PostgreSQL 17.11); der Slice, der diese ADR
umsetzt, wiederholt jede Zeile als Test.

| # | Befund | Herkunft |
|---|---|---|
| M1 | `CREATE_REPLICATION_SLOT … TEMPORARY LOGICAL pgoutput EXPORT_SNAPSHOT` (die Syntax, die `pglogrepl` erzeugt) liefert `consistent_point` und `snapshot_name`. Eine zweite Verbindung mit `BEGIN ISOLATION LEVEL REPEATABLE READ READ ONLY; SET TRANSACTION SNAPSHOT '…'` liest den Stand zum Slot-Punkt: gedruckt `im Snapshot sichtbar: 1` gegen `aktuell: 1,2` (Zeile 2 danach committet). Der Commit-Record der Folgezeile liegt bei `0/17D7FF8`, der `consistent_point` bei `0/17D7F28` (`pg_walinspect`). | gemessen |
| M2 | Eine Rolle nur mit `LOGIN REPLICATION` (kein Superuser) legt den temporären Slot an; ihr `SELECT` auf die Tabelle scheitert ohne Grant mit `permission denied for table t`. | gemessen |
| M3 | Die Slot-Anlage blockiert, bis eine beim Aufruf laufende Schreibtransaktion endet: 5,27 s bei rund 5 s Restlaufzeit. | gemessen |
| M4 | Nach dem Import darf die exportierende Verbindung enden (`count(*)` aus `pg_replication_slots` danach `0`); die importierende Transaktion liest mit einem `NO SCROLL`-Cursor weiter den Snapshot-Stand (5 Zeilen, weder die danach eingefügte noch die danach geänderte). | gemessen |
| M5 | Der `pgoutput`-Text der Werte folgt den GUC der Sitzung der Rolle: mit `ALTER ROLE … SET timezone='Asia/Tokyo'`, `datestyle='German, DMY'` liefert der Walsender (`pg_recvlogical`) `23.09.2026 17:11:12.5 JST`, `SHOW` über `replication=database` meldet dieselben Werte; mit Standard-GUC `2026-09-23 08:11:12.5+00`. | gemessen |
| M6 | `pglogrepl` (Pin in `go.mod`) trägt `CreateReplicationSlotOptions{Temporary, SnapshotAction}` und im Ergebnis `SnapshotName`; die Repo-Naht `driverSession.CreateReplicationSlot` reicht sie unverändert durch. | gelesen |

Unbelegt und deshalb als **erwartet** geführt: dass eine reguläre pgx-
Verbindung mit demselben DSN dieselben GUC-Werte wie der Walsender trägt
(M5 belegt nur, dass der Walsender die Rollen-Defaults trägt); der Slice
belegt es mit dem Paritätstest aus §Fitness Function.

### Konstraints

- Der Capture-kritische Pfad (Persist → ACK, `ADR-0011`/`ADR-0027`) bleibt
  unberührt; der Backfill läuft daneben (eigene Goroutine, eigene
  Verbindungen — das Muster von `runHeartbeat`, `runRetentionCleanup`,
  `runAdministration`).
- Eine `Accepted` ADR wird nicht überschrieben (`AGENTS.md` §3.5): `ADR-0050`
  (Antrags-Queue), `ADR-0059`/`ADR-0065` (Spaltenausschluss), `ADR-0081`
  (Lesen über HTTP) bleiben in Kraft; diese ADR ergänzt sie.
- Die Row-Image-Form und die Ausschluss-Zusage sind Vertrag
  (`LH-FA-CAP-008`, `LH-QA-SEC-004`): ein Backfill-Image darf sich von einem
  WAL-Image derselben Zeile nicht unterscheiden.

## Entscheidung

Wir wählen **einen Backfill als eigenen, explizit ausgelösten Lauf
(„Run"): eine Bulk-Copy des Tabellenbestands in einem
`REPEATABLE READ`-Snapshot, den ein je Run angelegter temporärer logischer
Slot exportiert; die Kopie wird in einer **einzigen** Store-Transaktion
committet, die Changes tragen `operation = INSERT` und ein neues Feld
`origin = 'backfill'`.** Acht Teilfragen.

### Teilfrage 1 — Mechanismus

| Option | Pro | Contra |
|---|---|---|
| A — Export-Snapshot des **Haupt**-Slots bei dessen Anlage | ein Slot, keine Zusatzverbindung | trägt nur die Erstanlage des Slots (`ensureSlot` legt ihn einmal an); eine später aktivierte Tabelle (`ADR-0050`) und jeder Wiederholungs-Backfill bleiben unbedient; `NOEXPORT_SNAPSHOT` ist heute bewusst gesetzt, der exportierte Snapshot lebt nur bis zum nächsten Befehl auf der Verbindung, die sofort in `START_REPLICATION` wechselt |
| B — Bulk-Copy in `REPEATABLE READ`, Position aus `pg_current_wal_lsn()` neben dem Snapshot gelesen | kein Slot, keine Replication-Rolle nötig | Snapshot und LSN sind nicht atomar gepaart: eine Transaktion kann ihren Commit-Record schon geschrieben, aber noch nicht für Snapshots sichtbar gemacht haben — ihr Commit liegt dann vor der gelesenen LSN und fehlt im Snapshot; wird der Backfill an diese LSN gesetzt, sortiert der Bestand hinter der neueren WAL-Änderung und überschreibt sie mit veraltetem Stand (hergeleitet aus dem Commit-Ablauf, nicht gemessen) |
| C — Bulk-Copy im **Slot-Snapshot**, temporärer Slot je Run **(gewählt)** | Snapshot und Punkt sind exakt gepaart (M1); wirkt für jede Tabelle zu jedem Zeitpunkt; der Slot verschwindet, sobald der Snapshot importiert und die Verbindung beendet ist (M4), kein dauerhaftes WAL-Pinning; nur `cdc_capture` (M2) | ein zusätzlicher Slot und Walsender für die Dauer der Anlage (`max_replication_slots`/`max_wal_senders` brauchen Reserve); die Anlage wartet auf laufende Schreibtransaktionen (M3); die Quelle hält einen Snapshot für die Kopierdauer (Vacuum-Horizont) |
| D — Zeilen an der Quelle „anstoßen" (`UPDATE … SET c = c`), damit der WAL-Pfad sie liefert | kein neuer Lesepfad | schreibt in die Quelle (Trigger, `updated_at`-Spalten, Schreiblast, Zeilenversionen — gegen `LH-QA-PER-001`), erzeugt `UPDATE`-Changes ohne Backfill-Herkunft, verletzt „Bestand zum Startzeitpunkt" bei nebenläufigen Schreibern |
| E — nichts tun | kein Aufwand | `LH-FA-CAP-009` bleibt unerfüllt |

**Ablauf eines Runs** (Verbindungen alle über `CDC_CAPTURE_DSN`):

1. Replication-Verbindung (`replication=database`):
   `CREATE_REPLICATION_SLOT cdc_bf_<run-id> TEMPORARY LOGICAL pgoutput
   EXPORT_SNAPSHOT` → `consistent_point` **X** (LSN, über die Positions-
   Abbildung von `ADR-0005` ein Offset) und `snapshot_name`. Der Slot-Name
   trägt das Bezeichner-Alphabet (`identifierShape`); die Run-Kennung geht
   ohne Bindestriche ein. Die Anlage trägt ein Zeitlimit; sein Ablauf
   endet den Run als `failed` (Fehlerklasse `transient`, `SPEC-008`).
2. Reguläre Verbindung: `BEGIN ISOLATION LEVEL REPEATABLE READ READ ONLY;
   SET TRANSACTION SNAPSHOT '<snapshot_name>'`. Danach endet die
   Replication-Verbindung (Slot fällt weg, M4).
3. Im Snapshot: Spaltenliste der Tabelle (`pg_attribute`, `attnum > 0`,
   nicht gelöscht, **nicht generiert** — der WAL-Pfad sendet generierte
   Spalten nicht), Zeilen über einen `NO SCROLL`-Cursor in **Blöcken** von
   `B` Zeilen (Startwert 1.000, Wert legt der Slice fest) mit
   `col::text` je Spalte — der Speicherbedarf ist durch `B` begrenzt
   (`LH-FA-CAP-006.a`: keine unbegrenzte RAM-Haltung).
4. Blöcke gehen in die Store-Transaktion nach Teilfrage 4.

Die Tabelle wird über **keinen** Datei-, Dump- oder Zwischenspeicher
geführt; ein Backfill ohne Zeilen (leere Tabelle) endet als `completed`
mit 0 Zeilen und schreibt keine Transaktion.

Ein Backfill setzt eine **aktivierte Tabelle mit laufender Bindung**
voraus (Bindungs-Zeile über `TableActivationPort.Registered`, Publication-
Mitgliedschaft über `Published`). Diese Vorbedingung ist der Grund, warum
keine Lücke entsteht (Teilfrage 3). Die Rolle `cdc_capture` braucht
zusätzlich `SELECT` auf die Quelltabelle — eine **Betriebs-Vorbedingung**
wie `REPLICATION` (`ADR-0047`), keine Vergabe durch diese Software; ihr
Fehlen endet den Run als `failed` (Klasse `permission`, M2).

### Teilfrage 2 — Markierung: Feld `origin`, Operation `INSERT`

| Option | Pro | Contra |
|---|---|---|
| A — neuer Operationswert (`SNAPSHOT`/`READ`) | im bestehenden Feld sichtbar | erweitert `chk_change_operation`: d-migrate konvergiert eine CHECK-Änderung an einer bestehenden Tabelle nicht (Kopfkommentar `schema.yaml`), also Nacharbeit-SQL; ändert die geschlossene Operationsmenge aus `SPEC-002`/`LH-FA-DAT-003` und bricht jeden Consumer, der über `INSERT\|UPDATE\|DELETE` verzweigt |
| B — Kennungs-Konvention (`0bf-…`-Präfix als einzige Markierung) | keine Schema-Änderung | Kopplung an ein Zeichenketten-Format; kein sauberer SQL-Filter |
| C — additive Spalte `origin` auf `cdc.change` **(gewählt)** | genau der Weg, den `LH-FA-DAT-006` und `ADR-0058` Entscheidung 2 real belegen (additive Spalte, alte Zeilen bleiben lesbar); filterbar; Operation bleibt in der Menge | Änderung an View und Lesepfaden |
| D — eigene Tabelle `cdc.backfill_change`, in `cdc.changes` per `UNION ALL` gemischt | keine Änderung an `cdc.change` | zweiter Persistenz-, Lese-, Retention- und Lösch-Pfad; `ORDER BY` über die Vereinigung |

**Festlegung.** `cdc.change.origin` ist `text`, **nullable**, ohne
Default und ohne CHECK (`NULL` ≙ `wal`); der Store schreibt `wal` bzw.
`backfill`; die geschlossene Menge `wal|backfill` erzwingt die Domäne
(`model.ChangeOrigin`, Konstruktor-Default `wal`), nicht die Datenbank —
dieselbe Begründung wie bei `request_kind`. Die View `cdc.changes` führt
`COALESCE(c.origin, 'wal') AS origin` als **letzte** Spalte (`LH-FA-DAT-006`
Boundary: gespeicherte Changes ohne das Feld lesen sich unverändert). Ein
Backfill-Change ist `operation = INSERT`, `old_data` fehlt, `new_data` trägt
das Row Image: Consumer, die `INSERT` als Upsert behandeln, brauchen keine
Änderung; Consumer, die `INSERT` strikt als „neu" lesen, filtern auf
`origin`. **Row-Image-Parität:** die Bild-Konstruktion ist **eine** Funktion,
die WAL- und Backfill-Pfad gemeinsam aufrufen (`mapper.rowImage` im
Driving-Adapter ist der Ausgangspunkt; die Funktion liegt an einer Stelle,
die beide Adapter-Schichten importieren dürfen) — `col::text` in derselben Sitzungs-GUC-Lage wie der
Walsender, generierte Spalten ausgenommen, `NULL` entfällt, Ausschluss
angewandt. Eine zweite JSON-Erzeugung wäre ein zweiter Träger derselben
Aussage.

### Teilfrage 3 — Verhältnis zum WAL-Pfad, Überlappungsfenster

**Position.** Alle Blöcke eines Runs liegen auf **`p = X`**, dem
`consistent_point` des Snapshots. Aus M1 folgt die Trennung: jeder Commit
mit Position ≤ X steckt im Bestand, jeder Commit mit Position > X ist
nicht darin und kommt über den Stream, sortiert hinter dem Backfill.

| Fall | Ergebnis im Log |
|---|---|
| Zeile seit Aktivierung nicht geändert | nur Backfill-Change |
| Zeile mit Commit ≤ X geändert, WAL-Change persistiert (Position ≤ X) | WAL-Change(s) davor, Backfill-Stand danach — **bekannte, begrenzte Dopplung** (das Fenster `(Aktivierung, X]`), Backfill-Stand ist neuer oder gleich |
| Zeile mit Commit > X geändert | Backfill-Stand, danach der WAL-Change — **keine** Dopplung |
| Zeile mit Commit ≤ X gelöscht | nicht im Backfill; der WAL-`DELETE` (falls persistiert) steht davor |
| Commit einer aktivierten Tabelle, den der Assembler vor `AddBinding` verwarf | Commit-Position ≤ X (die Bindung steht vor der Slot-Anlage) → im Bestand enthalten, **keine Lücke** |

**Begründung der Lückenfreiheit** (hergeleitet): der Backfill startet erst
nach Abschluss der Aktivierung (`enable` und `backfill` werden von
derselben Goroutine nacheinander verarbeitet, `AddBinding` liegt vor dem
Slot-Anlegen); `X` liegt nicht vor der Stream-Position zum Zeitpunkt von
`AddBinding`; jeder Commit > X wird mit gesetzter Bindung dekodiert.

**Replay-Invariante** (hergeleitet, im Slice als Eigenschaftstest zu
belegen): wendet ein Consumer das Log ab einem Punkt vor `p` in
Lese-Ordnung an — `INSERT`/`UPDATE` als Upsert des Row Images, `DELETE` als
Löschen —, entspricht sein Stand je Schlüssel dem Quellstand; die Dopplung
ist damit **idempotent**, nicht nur begrenzt (`ADR-0012`: Consumer arbeiten
ohnehin idempotent). Ausnahme ist die bestehende Eigenschaft, dass ein
`UPDATE`-Image unveränderte TOAST-Spalten als Abwesenheit trägt
(`LH-FA-CAP-008` Boundary) — der Backfill-Change trägt sie **vollständig**.

**Sichtbarkeit — die Grenze dieser Entscheidung.** Positionen wandern nur
vorwärts (`ADR-0029` Regel 2, `LH-FA-CON-004`); ein Backfill liegt auf `X`.
Ein Consumer, dessen bestätigte Position bei Commit des Runs bereits
`≥ X` ist, sieht ihn beim Lesen ab seiner Position **nicht** — der Backfill
ist ein Zustandsabzug für Consumer **vor** `X`: neu registrierte Consumer
(Anfangsposition = Log-Anfang, `LH-FA-CON-005` Boundary), zurückliegende
Consumer und ein eigens für den Abzug registrierter Consumer. Ein live
folgender Consumer holt den Bestand über den expliziten Bereichslesezugriff
(`LH-FA-REA-001`, `from`/`to`, Filter Tabelle) oder über einen zweiten
Consumer. Diese Grenze ist Folge des Positionsmodells, kein Fehler des
Mechanismus.

### Teilfrage 4 — Atomarität, Unterbrechung, Wiederaufnahme

**Festlegung.** Der Run schreibt alle Blöcke in **eine** Store-Transaktion
und committet sie **einmal am Ende**, zusammen mit der Zeile
`cdc.backfill_run` (`status = completed`, `rows_copied`). Bis dahin ist
nichts für Leser sichtbar; ein Prozessende, ein Verbindungsverlust oder ein
Abbruch rollt alles zurück.

| Option | Pro | Contra |
|---|---|---|
| A — Block-weise committen | kürzere Schreibtransaktion | ein Consumer liest den Teilbestand, bestätigt `X` und verpasst die später committeten Blöcke still — genau die „stille Lücke" der Boundary/Negative-Kriterien |
| B — Fortsetzen mit Checkpoint (`last_key`) | keine Neuarbeit | braucht eine Schlüsselordnung (Tabellen ohne Primärschlüssel gehen nicht), einen neuen Snapshot bei Wiederaufnahme und Sichtbarkeits-Gating je Block; Ausbaustufe (Durchsatz sehr großer Tabellen ist laut Lastenheft nicht Voraussetzung) |
| C — eine Transaktion, Neubeginn nach Abbruch **(gewählt)** | Negative-Kriterium trivial („beginnt neu, ohne Bestand zu verlieren": ein abgebrochener Run hinterlässt nichts); keine Teilzustände, keine Dopplung durch Abbrüche | die Schreibtransaktion dauert so lange wie die Kopie; sehr große Tabellen belasten Vacuum-Horizont (Quelle: Snapshot; CDC-Speicher: offene Transaktion) |

**Run-Zustand.** `cdc.backfill_run` (neue Tabelle; Schlüssel = die
`administration_request_id` des Antrags, die die SQL-Funktion vergibt):
`run_id`, `source_id`, `schema_name`, `table_name`, `status`
(`queued|running|completed|failed|interrupted`, CHECK bei Erstanlage der
Tabelle — konvergiert; eine spätere Erweiterung läuft über Nacharbeit-SQL),
`requested_at`, `started_at`, `finished_at`, `snapshot_position`,
`rows_copied` (Fortschritt je Block, außerhalb der Daten-Transaktion,
damit er sichtbar ist), `error_message`. Beim Prozessstart setzt die
Composition Root jeden `running`-Run auf `interrupted` (eine Instanz je
Quelle trägt den Slot, `ADR-0050`-Muster für `pending`); `queued`-Runs
laufen weiter. **Kein automatischer Neustart** eines `interrupted`-Runs:
ein erneuter Antrag (`cdc.backfill_table`) startet einen neuen Run mit neuem
Snapshot — „wenn er erneut gestartet wird, beginnt er neu"
(`LH-FA-CAP-009` Negative). Ein zweiter Antrag für dieselbe Tabelle, solange
ein Run `queued|running` ist, endet als `failed` mit Grund (kein zweiter
paralleler Run je Tabelle).

**Fail-closed vor dem Commit** (Sicherheit, `LH-QA-SEC-004`): unmittelbar
vor dem Commit prüft der Run, dass (a) die Bindung der Tabelle noch
besteht und (b) der Ausschlussstand (`ColumnExclusionPort.ExcludedColumns`)
demselben Stand entspricht, mit dem die Blöcke gebaut wurden; jede
Abweichung rollt den Run zurück (`failed`, Grund im `error_message`). Jeder
Block liest den Ausschlussstand neu und baut das Bild in Go mit dem
aktuellen Stand (der Cursor liefert alle Spalten; ausgeschlossene Werte
werden nie serialisiert und nie an die Persistenz übergeben, dieselbe
Zusage wie im WAL-Pfad). Ohne die Prüfung würde ein Run, der vor einem
`exclude_column` Blöcke gebaut hat, nach dem Ausschluss einen Change mit
dem ausgeschlossenen Wert sichtbar machen.

### Teilfrage 5 — Auslösung, Administration, Sichtbarkeit

Der Backfill wird **explizit** ausgelöst, nicht implizit von `enable`
(`LH-FA-CAP-009` Happy Path: „aktiviert und ein Backfill ausgelöst").

- **Antragsart** `backfill` in `cdc.administration_request` und
  `model.AdministrationRequestKind`; SQL-Funktion
  `cdc.backfill_table(p_source_id, p_schema_name, p_table_name) RETURNS text`
  (SECURITY DEFINER, `REVOKE … FROM PUBLIC`, `GRANT EXECUTE` an
  `cdc_admin`) im Muster von `cdc.enable_table`: **schreibt nur den Antrag**
  und sendet `pg_notify` (`ADR-0050` Fitness Function bleibt erfüllt).
- **Ausführung** über einen Inbound Port `BackfillTableUseCase`
  (`ADR-0028`); die Administrations-Goroutine blockiert **nicht** auf der
  Kopie: sie prüft die Vorbedingungen (Tabelle aktiviert, kein aktiver Run
  derselben Tabelle), legt die Run-Zeile `queued` an, vermerkt den Antrag
  `applied` („angenommen") und übergibt an einen **einzelnen
  Backfill-Worker** (eine Goroutine, ein Run zugleich — keine
  Parallelisierung, Lastenheft Out-of-Scope; Folge-Anträge warten `queued`
  in Antragsreihenfolge). Vorbedingungs-Fehler enden den Antrag als
  `failed` mit Text. Die Ausführung des Runs steht in `cdc.backfill_run`,
  nicht im Antragsstatus.
- **Verbindungen und Rollen.** Der Worker hält einen eigenen Pool aus
  `CDC_CAPTURE_DSN` (Schreiben `cdc.transaction`/`cdc.change`,
  `SELECT/INSERT/UPDATE` auf `cdc.backfill_run`, `SELECT` auf der
  Quelltabelle, `REPLICATION` für den temporären Slot) — Muster der
  Retention-Goroutine mit eigenem Pool, keine Berührung des
  Capture-Pfads. `cdc_reader` bekommt `SELECT` auf eine neue View
  `cdc.backfill_status` (letzter Run je Tabelle: Status, Zeilen, Zeiten,
  Fehlertext, geschätzte Zeilenzahl).
- **Diagnose.** Der Sondermodus `diagnose` (`bootstrap.Diagnose`,
  `LH-FA-SST-003`) liest `cdc.backfill_status` zusätzlich und gibt je
  Tabelle Status und Fortschritt aus; ein `failed`/`interrupted`-Run ist
  Berichtsinhalt, kein Befehlsfehler.
- **Fehlerklassen und Heartbeat.** Run-Fehler tragen die `SPEC-008`-Klasse
  im Fehlertext (`permission`, `configuration`, `storage`, `transient`,
  `replication`) und sind **run-lokal**: sie setzen weder den
  Heartbeat-Fehlerzustand noch stoppen sie den Capture-Pfad.
- **Wecksignal.** Nach dem Commit sendet der Run je Tabelle **ein**
  Wecksignal über den bestehenden `ChangeNotificationPort` (best effort,
  `ADR-0055`). Backfill-Changes gehen **nicht** in den
  `ChangeStreamPort`/`Broadcaster`: der Live-Stream ist ein
  nichtblockierender Live-Weg mit Verwerfen langsamer Abonnenten
  (`ADR-0066`) und trägt keinen Zustandsabzug in Millionenzahl.

### Teilfrage 6 — Ordnung, `change_id`, Position (`LH-FA-CAP-004`)

- **Position:** alle Blöcke eines Runs `p = X` (Teilfrage 3);
  `commit_position` ist nicht eindeutig je Transaktion (Index
  `cdc_transaction_source_position_idx`, kein UNIQUE), die Lese-Ordnung
  bleibt `(commit_position, transaction_id, sequence)`.
- **Block = Transaktion.** Jeder Block ist eine eigene synthetische
  `ChangeTransaction` mit Kennung `0bf-<run-id>-<Blocknummer, 8 Stellen,
  null-aufgefüllt>` und Sequenz `1…B`; die Kennung sortiert wegen des
  Präfix `0` **vor** jeder WAL-Kennung (WAL-Kennungen sind Ziffernfolgen
  ohne führende Null) — bei gleicher Position, sollte je ein WAL-Commit
  genau auf `X` liegen, liegt der Backfill davor (die sichere Richtung,
  Teilfrage 3). `change_id` = `<Transaktions-ID>-<Sequenz>`, dieselbe
  Bildungsregel wie im WAL-Pfad; eindeutig, weil Run-Kennung und
  Blocknummer eingehen.
- **Reihenfolge innerhalb des Backfills** ist die Cursor-Reihenfolge,
  einmal vergeben und in `sequence`/Blocknummer festgeschrieben, damit
  stabil (`LH-FA-REA-004` Boundary); eine fachliche Ausführungsreihenfolge
  gibt es für einen Zustandsabzug nicht (`LH-FA-CAP-004` verlangt sie für
  Änderungen, nicht für den Bestand) — der Abzug ist gegenüber den
  WAL-Changes eindeutig geordnet (Position `X`).
- **Paging-Grenze.** Alle Blöcke teilen `X`; `GET /changes` mit `Limit`
  schneidet innerhalb einer Position und kann sie nicht fortsetzen
  (Befund oben; dasselbe gilt heute für jede Quelltransaktion mit mehr
  Changes als `Limit`). Bestandsabzüge lesen deshalb **ohne `Limit`** oder
  über die View mit Schlüsselvergleich
  (`WHERE (commit_position, transaction_id, sequence) > (…)`). Eine
  Cursor-Form am HTTP-Lesepfad ist der Re-Evaluierungs-Trigger aus
  `ADR-0081`; sie ist **nicht** Teil dieser Entscheidung und **keine**
  Voraussetzung für die Akzeptanzkriterien, verschärft aber ihre
  Dringlichkeit (Folgepflicht 6).

### Teilfrage 7 — Retention (`ADR-0014`, `ADR-0053`, `ADR-0040`)

Backfill-Changes folgen **unverändert** der bestehenden Retention-Regel,
kein Sonderpfad: `committed_at` der synthetischen Transaktion ist der
Snapshot-Zeitpunkt (Uhr über `ClockPort`, `ADR-0040`), das Alter zählt
also ab dem Backfill; das produktive `MinAge` von 24 h hält den Bestand
mindestens so lange lesbar; danach gibt `AllowsDeletion` ihn frei, sobald
jede bestätigte Consumer-Position an oder hinter `X` liegt. Ein
zurückliegender Consumer blockiert die Löschung wie bei jedem Change
(`cdc.retention_blockers`). `RunRetentionService.Run` liest alle Changes
der Quelle in den Speicher und übergibt die freigegebene Menge als
Kennungsliste an `DeleteChanges`; der Bestand einer großen Tabelle
vergrößert beide — die Grenze besteht schon, wird aber sichtbar (siehe
§Konsequenzen und Re-Evaluierungs-Trigger).

### Teilfrage 8 — Reichweite in Lesewegen und SDKs

- **Betroffen:** View `cdc.changes` und `queries.SelectChanges` (Spalte
  `origin` als letzte), `postgresstorage/mapper` (Zeile ↔ Change),
  `GET /changes` (`SPEC-022`: Feld `origin`), Doku.
- **Unverändert:** die drei Live-Wege (gRPC `SPEC-020`, SSE `SPEC-021`,
  NATS-Vollinhalt `SPEC-024`) — sie tragen ausschließlich WAL-Changes, ein
  Feld `origin` wäre dort konstant `wal`; die Nachrichtenschemata und die
  generierten Proto-Artefakte bleiben, `make generated-sync` ist nicht
  berührt. Reihen sich Backfill-Changes später in den Live-Weg ein, ist das
  eine Folge-ADR (additives Feld, `string origin = 11`).
- **SDKs:** die HTTP-Lesemodelle (C# `Http`, Kotlin `http.model.Changes`,
  Python `http_client`) und der Go-HTTP-Beispiel-Client bekommen `origin`
  als optionales Feld (fehlt es, gilt `wal`); ohne die Anpassung
  funktionieren sie unverändert weiter (Decoder ignorieren das Feld).
  Die `LH-FA-SST-009`-Package-Version hebt der SDK-Slice nach dem
  bestehenden Muster.

## Konsequenzen

- Positiv: `LH-FA-CAP-009` ist mit einem Mechanismus erfüllbar, der für
  jede aktivierte Tabelle zu jedem Zeitpunkt wirkt (auch für
  Wiederholungen), die Quelle nicht beschreibt (kein Trigger, kein
  `UPDATE`) und den WAL-Pfad nicht berührt.
- Positiv: Boundary und Negative sind strukturell erfüllt statt
  aufgefangen — der Bestand ist ab `X` lückenlos an den Stream angeschlossen
  (Teilfrage 3), ein abgebrochener Run hinterlässt keine Teilzustände
  (Teilfrage 4).
- Positiv: eine gemeinsame Row-Image-Funktion beseitigt die Möglichkeit
  zweier auseinanderlaufender Bildformen.
- Positiv: `origin` ist additiv (`LH-FA-DAT-006`); kein Live-Schema, kein
  Proto, keine SDK-Zusage bricht.
- Negativ: die Sichtbarkeits-Grenze (Teilfrage 3): der Backfill ist für
  Consumer vor `X` da, nicht für live folgende. Ein Betreiber, der das
  nicht liest, erwartet „mein laufender Consumer bekommt den Bestand"; die
  Benutzerdoku muss es an sichtbarer Stelle sagen.
- Negativ: eine lange Kopie hält einen Snapshot an der Quelle
  (Vacuum-Horizont) und eine offene Schreibtransaktion im CDC-Speicher;
  große Tabellen verschieben das Problem in die Ausbaustufe (Checkpoint,
  Blöcke mit eigener Transaktion, Parallelisierung).
- Negativ: die Slot-Anlage wartet auf laufende Schreibtransaktionen der
  Quelle (M3) und braucht Reserve in `max_replication_slots`/
  `max_wal_senders`; beides sind Betriebs-Vorbedingungen, die der Run als
  sichtbare Fehler meldet.
- Negativ: `SELECT` auf die Quelltabelle für `cdc_capture` und eine erweiterte
  Schema-/Rollen-Fläche (Spalte, Tabelle, zwei Views, eine Funktion,
  Grants, Idempotenz-Guard).
- **Akzeptierte Negative** (kurz begründet, keine Folgepflicht):
  - *Schema-Version-Verweis:* Backfill-Changes referenzieren die zum
    Run-Start aktuelle Version; eine kompatible Spalten-Erweiterung
    während des Runs erscheint im Image ohne Versions-Bump — harmlos, das
    Image ist ein selbstbeschreibendes JSON-Objekt.
  - *Historien-Consumer und Fenster ≤ X:* ein Consumer, der den
    Backfill bestätigt, verpasst WAL-Events mit Position ≤ X, die der
    Capture-Pfad erst **nach** dem Run-Commit persistiert (nur bei
    Capture-Lag größer als die Kopierdauer, `SPEC-013`: p95 ≤ 1 s); ihr
    Zustand ist im Bestand enthalten, Zustands-Consumer verlieren nichts.
  - *Backfill-Changes ohne Live-Zustellung:* Consumer, die nur einen
    Live-Weg abonnieren, erhalten den Bestand nicht — er ist ein
    Lese-Weg-Ereignis (Wecksignal genügt zum Nachholen).

### Folgepflichten

1. **Spec-Nachzug (Planner/Architect, ohne ADR-Bezüge in der Sicht):**
   `LH-FA-CAP-009.a` (Mechanismus, Markierung, Überlappung an die
   Entscheidung binden), `SPEC-002` (`origin`), `SPEC-019` (Antragsart
   `backfill`), `SPEC-022` (`origin`), eine **neue** `SPEC-*`-Kennung für
   `cdc.backfill_run`/`cdc.backfill_status` (Kennung vergibt der Nachzug);
   `spec/architecture.md` (Backfill-Sequenz im Sicht-Stratum, ohne
   ADR-/Slice-Bezug); `docs/user/benutzerhandbuch.md` (Auslösung,
   Sichtbarkeits-Grenze, `SELECT`-Grant, Lesen ohne `Limit`).
2. **Umsetzung — Slice-Schnitt-Empfehlung (Planner entscheidet final),
   eine Welle:**
   - **S1 `backfill-change-origin`** — `model.ChangeOrigin`/`Change.Origin`,
     Store (`InsertChange`, `SelectChanges`, Mapper), `schema.yaml`
     (Spalte, View `changes` +`origin`), `GET /changes` (`SPEC-022`),
     gemeinsame Row-Image-Funktion (Verschiebung ohne Verhaltensänderung).
     Eigenständig lieferbar; belegt `LH-FA-DAT-006` Boundary (alte Zeilen
     `NULL` ≙ `wal`).
   - **S2 `backfill-snapshot-reader`** — Outbound Port
     (`TableSnapshotPort`: Snapshot öffnen → Position `X` + Blöcke lesen →
     schließen) und Driven Adapter (temporärer Slot, Import, Cursor,
     Spaltenliste, GUC-Parität) mit Tests gegen reale PostgreSQL in der
     Tier von `make test-replication`. Trägt M1–M5 als Tests.
   - **S3 `backfill-run-usecase`** — Domäne (`BackfillRun`), Port für
     Run-Zustand und atomaren Schreiber (`ADR-0034`: Fähigkeits-Ports),
     `BackfillTableUseCase` samt Fail-closed-Prüfung, `cdc.backfill_run`
     (Schema, Grants, geschätzte Zeilenzahl), Wecksignal. Unit-Tests mit Fakes plus Store-Tier.
   - **S4 `backfill-sql-administration`** — Antragsart und
     `cdc.backfill_table` (`nacharbeit-administration.sql`, CHECK-Menge),
     Erweiterung von `knownForeignObjects` im Idempotenz-Guard, Übergabe
     Administrations-Goroutine → Worker, Start-Abgleich
     (`running → interrupted`), View `cdc.backfill_status`, Grant,
     `diagnose`-Ausgabe.
   - **S5 `backfill-e2e`** — `make test-integration`: Happy Path (Bestand
     über `cdc.changes` und `GET /changes`, `origin = 'backfill'`),
     Boundary (nebenläufige Schreiber während des Runs, Replay-Invariante),
     Negative (`docker kill` mitten im Run, `interrupted`, erneuter Antrag);
     E2E-Abdeckung für `LH-FA-CAP-009`; Startposition eines frisch
     registrierten Consumers gemessen und dokumentiert; Bench-Beleg der
     Kopierdauer je Tabellengröße, aus dem die Warn-Richtgröße folgt.
   - **S6 `backfill-sdk-origin`** — `origin` in den drei SDK-HTTP-Modellen
     und dem Go-Beispiel-Client (unabhängig von S2–S5 nach S1 lieferbar).
3. **Nicht Teil dieser Entscheidung (Ausbaustufe/Folge-ADR):** Fortsetzen
   mit Checkpoint und Blöcke mit eigener Transaktion, Parallelisierung,
   `POST /backfill`-Endpunkt und CLI-Auslösung (additiv über denselben Use
   Case), Live-Zustellung von Backfill-Changes, Filter `origin` an
   `GET /changes`, Auslösung als Option von `enable`.
4. **Kommentar-Nachzug:** der Kommentar zu `model.Change` und zur
   Operationsmenge in `internal/domain/model/change.go` trägt nach S1 das
   Feld `origin` (`AGENTS.md` §3.7, §3.13: Träger, die eine bewegte
   Eigenschaft beschreiben, mitziehen).
5. **Idempotenz-Guard:** jede neue Funktion/View außerhalb des neutralen
   Modells gehört in `knownForeignObjects` — sonst bricht ein zweiter
   Rollout mit Exit 8 (S4).
6. **Cursor-Form (Priorität nach S5):** Re-Evaluierung von `ADR-0081`
   für eine Fortsetzung innerhalb einer Position; bis dahin trägt die Doku
   die Lese-Regel „Bestandsabzug ohne `Limit`". Die Beobachtung
   `BEO-PGC/limit-fortsetzung-innerhalb-einer-position` führt den Befund.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Test, reale PostgreSQL (Tier `test-replication`) | **Snapshot-Träger:** Slot-Snapshot sieht die vor `X` committeten Zeilen und keine danach committete (M1); Rolle ohne `SELECT` → Klasse `permission` (M2); Slot fällt nach Import und Verbindungsende weg (M4) | `make test-replication` |
| Go-Test, reale PostgreSQL | **Bild-Parität:** dasselbe Zeilenbild über WAL-Pfad und Backfill-Pfad byte-gleich — Typen `timestamptz`, `float8`, `bytea`, `interval`, `numeric`, `jsonb`, Array, `date`, generierte Spalte — unter Rollen-GUC ≠ Standard (`timezone`, `datestyle`; M5) | `make test-replication` |
| Go-Test, reale PostgreSQL (Eigenschaftstest) | **Replay-Invariante:** Backfill mit nebenläufigen `INSERT`/`UPDATE`/`DELETE` an der Tabelle; das Log ab Log-Anfang, angewandt als Upsert/Delete, ergibt den Quellstand | `make test-replication` |
| Go-Test (Fakes) | Run bricht ab und hinterlässt keine Zeilen, wenn Bindung oder Ausschlussstand vor dem Commit abweichen (`LH-QA-SEC-004`); Fehlerklasse je Ursache; `running → interrupted` beim Start | `make test` |
| Go-Test | `model.ChangeOrigin`-Menge geschlossen, Konstruktor-Default `wal`; `NULL`-Zeile liest als `wal` | `make test` |
| Go-Test, reale PostgreSQL | `cdc.backfill_table` schreibt ausschließlich einen Antrag; `PUBLIC` ohne `EXECUTE`; Rollout zweimal hintereinander idempotent (`knownForeignObjects`) | `make test-store`, `make schema-rollout` |
| E2E | Happy/Boundary/Negative gegen den laufenden Feed-Container (S5) | `make test-integration` |
| `.a-check` | neue Pakete liegen in den bestehenden Globs `app`/`ports`/`adapters`; der Snapshot-Adapter importiert keinen anderen Adapter | `make a-check` |
| Review-Prüfpflicht | Row-Image-Konstruktion an genau einer Stelle; keine zweite JSON-Erzeugung für Row Images | — (kein Gate) |

## Re-Evaluierungs-Trigger

- **Tabellen, deren Kopierdauer den Snapshot-/Transaktions-Horizont
  untragbar macht** (beobachtbar an Vacuum-Rückstand der Quelle oder
  Laufzeit eines Runs über der Betriebs-Toleranz): Folge-ADR für
  Block-Transaktionen mit Checkpoint und Sichtbarkeits-Gating (Option B aus
  Teilfrage 4) und Parallelisierung — die Lastenheft-Ausbaustufe.
- **Ein Live-Weg soll Backfill-Changes zustellen** oder ein Consumer soll
  Backfill-Changes nach `origin` serverseitig filtern: Folge-ADR
  (additives Feld in `SPEC-020`/`021`/`024`, Filter-Grammatik nach
  `ADR-0081`).
- **`RunRetentionService` liest den Bestand nicht mehr tragfähig** (Arbeits-
  speicher oder Laufzeit je Lauf; die Grenze besteht unabhängig vom
  Backfill): Folge-ADR zu Retention-Paging bzw. Lösch-Prädikat statt
  Kennungsliste.
- **`CREATE_REPLICATION_SLOT … EXPORT_SNAPSHOT` entfällt oder ändert sich**
  in einer neuen PostgreSQL-Hauptversion (`SPEC-012`: 17/18 gemessen):
  Mechanismus neu prüfen (`SNAPSHOT 'export'`-Syntax ab PostgreSQL 15
  vorhanden).
- **d-migrate konvergiert CHECK-Änderungen**: `origin`-CHECK und
  `request_kind`-Menge können deklarativ werden (`ADR-0043`-Trigger).
- **Ein Betreiber muss einen live folgenden Consumer ohne Neuregistrierung
  „nachladen"**: Folge-ADR für einen administrativen Consumer-Reset
  (`ADR-0013`: „explizit administrativ", bewusste, protokollierte
  Sonderoperation).
- **Ein Betreiber muss den Backfill ohne SQL-Zugang auslösen**: additiver
  HTTP-Endpunkt oder CLI-Aufruf über denselben `BackfillTableUseCase`
  (eigener Slice, Admin-Berechtigung).
- Sonst permanent — Mechanismus, Markierung und Positionsregel gelten
  unabhängig davon, wie viele Tabellen einen Backfill nutzen.

### Festlegungen des Auftraggebers

1. **Sichtbarkeits-Grenze (Teilfrage 3):** Das Positionsmodell bleibt
   unverändert; ein Consumer-Reset ist nicht Teil des Erstumfangs. Ein
   Consumer, dessen Position beim Run-Commit bereits `X` erreicht hat,
   sieht den Bestand nicht über seinen Fortschritt; der Bestand bleibt über
   das Bereichslesen (`LH-FA-REA-001`, unabhängig von jeder Consumer-Position)
   und über einen frisch registrierten Consumer erreichbar. Von welcher
   Position ein frisch registrierter Consumer startet, ist **erwartet, nicht
   geprüft**; `S5` misst es und die Doku trägt das Ergebnis als Bezugsmuster.
   `ADR-0013` führt den Reset als „explizit administrativ"; er ist eine eigene
   Fähigkeit mit eigener Entscheidung (Re-Evaluierungs-Trigger).
2. **Auslösungs-Weg (Teilfrage 5):** Im Erstumfang genügt die SQL-Funktion
   `cdc.backfill_table`; `cdc.backfill_status` und der `diagnose`-Sondermodus
   sind der Lesezugriff. Ein HTTP-Endpunkt oder CLI-Aufruf ist Folgezug über
   denselben Use Case (Re-Evaluierungs-Trigger).
3. **Große Tabellen (Teilfrage 4):** Die Ein-Transaktions-Form gilt für den
   Erstumfang. `cdc.backfill_run` trägt die beim Antrag geschätzte Zeilenzahl
   (`pg_class.reltuples`, eine Schätzung) und die Laufzeit;
   `cdc.backfill_status` und `diagnose` zeigen beides. Es gibt eine
   Warnung, keine Ablehnung. Die Richtgröße, ab der gewarnt wird, ist **zu
   messen** (`S5`, `make bench`-Infrastruktur, `ADR-0054`) und wird danach
   in der Doku genannt.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-23 | Proposed — Architect-Vorschlag zur offenen, ADR-pflichtigen Frage `LH-FA-CAP-009.a`; Machbarkeit an einem Wegwerf-Container real geprüft (M1–M5), Annahme beim Auftraggeber | [`LH-FA-CAP-009.a`](../../../spec/pflichtenheft.md) |
| 2026-09-23 | Accepted — Annahme durch den Auftraggeber samt den Festlegungen zu Sichtbarkeit, Auslösung und großen Tabellen | [`LH-FA-CAP-009.a`](../../../spec/pflichtenheft.md) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0111` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
