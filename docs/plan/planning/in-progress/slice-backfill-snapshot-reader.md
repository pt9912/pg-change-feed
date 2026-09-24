# Slice backfill-snapshot-reader: Snapshot-Leser — Port `TableSnapshotPort` und Driven Adapter: temporärer Slot mit Export-Snapshot, Import, Cursor-Blöcke, GUC-Parität

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [welle-backfill-bestand](../welle-backfill-bestand.md).

**Bezug:** [`LH-FA-CAP-009`](../../../../spec/lastenheft.md) (Bestand über denselben Lesezugriffsweg),
[`LH-FA-CAP-008`](../../../../spec/lastenheft.md) (Row-Image-Abwesenheit), [`LH-FA-CAP-006.a`](../../../../spec/pflichtenheft.md) (keine
unbegrenzte RAM-Haltung), [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 1 (Mechanismus: Slot-Snapshot),
[`ADR-0034`](../../adr/0034-ports-nach-faehigkeiten.md) (Ports nach Fähigkeiten), [`ADR-0080`](../../adr/0080-nahtform-pgconn-adapter-treiberhuelle.md) (Nahtform der
Treiber-Hülle), [`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md) Punkt 1 und 3 (Messgegenstand des Coverage-Gates und der DB-Adapter-Coverage;
Re-Evaluierungs-Trigger (a): die namentliche Liste nachziehen),
[`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 3, Punkt 5 (`reltuples` = `−1` für nie analysierte Tabellen),
[`ADR-0115`](../../adr/0115-backfill-spaltenwerte-text-ergebnisformat.md) (Spaltenwerte im Text-Ergebnisformat der Ausgabefunktion, ersetzt die Cast-Stellen von [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md)),
[`ADR-0030`](../../adr/0030-testpyramide.md) (Testpyramide).

**Berührte Spec-Stellen:** [`SPEC-012`](../../../../spec/pflichtenheft.md) (PostgreSQL 17 und 18), [`ARC-004`](../../../../spec/architecture.md)
(Outbound Port), [`ARC-006`](../../../../spec/architecture.md) (Driven Adapter), [`ARC-008`](../../../../spec/architecture.md) (PostgreSQL Logical
Replication) — gelesen, nicht geändert.

**Verantwortlich:** Implementer-Agent, 2026-09-24.

**Autor:** Planner-Agent, Welle-Eröffnung [welle-backfill-bestand](../welle-backfill-bestand.md). **Datum:** 2026-09-23.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein Outbound Port `TableSnapshotPort` und sein Driven Adapter
liefern den Bestand einer Tabelle als konsistenten Snapshot: die
Replication-Verbindung legt `CREATE_REPLICATION_SLOT cdc_bf_<run> TEMPORARY
LOGICAL pgoutput EXPORT_SNAPSHOT` an und liefert `consistent_point` **X** und
`snapshot_name`; eine reguläre Verbindung öffnet `BEGIN ISOLATION LEVEL
REPEATABLE READ READ ONLY; SET TRANSACTION SNAPSHOT '…'`, danach endet die
Replication-Verbindung; im Snapshot liest der Adapter die Spaltenliste
(`pg_attribute`, `attnum > 0`, nicht gelöscht, **nicht generiert**) und die
Zeilen über einen `NO SCROLL`-Cursor in Blöcken von `B` Zeilen; jede Spalte steht
nur beim gequoteten Namen, ohne Cast und ohne Funktion, und kommt im
Text-Ergebnisformat als Ergebnis der Ausgabefunktion des Spaltentyps
([`ADR-0115`](../../adr/0115-backfill-spaltenwerte-text-ergebnisformat.md)
Festlegung 1/2; Startwert `B` = 1.000, der Slice legt ihn fest). Die Anlage trägt ein
Zeitlimit; sein Ablauf endet als Fehlerklasse `transient`. Alle Verbindungen
laufen über `CDC_CAPTURE_DSN`. Als Katalog-Fähigkeit derselben Quelle liest der
Adapter außerdem die **geschätzte Zeilenzahl** einer Tabelle
(`pg_class.reltuples`; ein negativer Wert heißt „unbekannt" — erwartet, in
diesem Slice an PostgreSQL 17 und 18 zu messen) für den Antrag des Runs. Der
Adapter liegt in einem **neuen eigenen Paket** (`postgressnapshot`), nicht in
`postgresstorage`: er liest die Quelltabelle über eine Replikationsverbindung, der
Run-Store bleibt in `postgresstorage` — zwei Verantwortungen.

Der Slice trägt die real gemessenen Befunde M1–M5 aus [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) als Tests und
die **Bild-Parität**: dasselbe Zeilenbild über WAL-Pfad und Backfill-Pfad
byte-gleich — mit der gemeinsamen Funktion aus `row-image-gemeinsam`.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Der Run** (Vorbedingungen, Fail-closed, Schreiben, Wecksignal) — `run-usecase`;
  der Snapshot-Leser kennt weder Run-Zustand noch Store.
- **Persistenz und Position im Store** — der Adapter liefert `X` als
  Rohwert der LSN; die Abbildung auf die Position ([`ADR-0005`](../../adr/0005-sourceposition-abstrahiert-lsn.md)) und ihre
  Ablage gehören dem Run und dem Store.
- **Die Replay-Invariante und der Abbruch mitten im Lauf** — sie brauchen den
  komponierten Run und gehören dem `e2e`-Slice.
- **Die Betriebs-Vorbedingungen** (`SELECT` auf die Quelltabelle,
  `max_replication_slots`/`max_wal_senders`-Reserve) — Doku im
  `sql-administration`-Slice; dieser Slice belegt nur, dass ihr Fehlen als
  Fehlerklasse (`permission`, `configuration`) sichtbar wird.

## 2. Definition of Done

- [x] Snapshot-Träger real belegt: der Slot-Snapshot sieht die vor `X`
      committeten Zeilen und keine danach committete (M1); eine Rolle nur mit
      `LOGIN REPLICATION` ohne `SELECT` scheitert mit der Klasse `permission`
      (M2); das Zeitlimit der Slot-Anlage gegen eine offene Schreibtransaktion
      endet als `transient` (M3); nach Import und Ende der Replication-Verbindung
      ist der temporäre Slot weg und der Cursor liest weiter den Snapshot-Stand
      (M4). *Zu belegen durch:* `make test-replication`; die Export-Syntax
      zusätzlich gegen PostgreSQL 17 (`PG_TEST_IMAGE` auf den 17er-Digest der
      CI-Matrix in `.github/workflows/e2e.yml`), damit beide Major-Versionen
      belegt sind ([`SPEC-012`](../../../../spec/pflichtenheft.md)).
- [x] Blöcke und Spalten: die Zeilen kommen in Blöcken von höchstens `B` (Test
      mit mehr als einem Block und exakt an der Blockgrenze), die Spaltenliste
      lässt generierte und gelöschte Spalten aus, eine leere Tabelle liefert
      keinen Block, die geschätzte Zeilenzahl einer frisch befüllten Tabelle ohne
      `ANALYZE` ist „unbekannt" (`reltuples = −1`, erwartet — der Bericht nennt den
      gemessenen Wert je Version mit seinem Lauf) und nach `ANALYZE` dicht an der
      tatsächlichen Zeilenzahl. *Zu belegen durch:* `make test-replication` (auf
      PostgreSQL 17 und 18).
- [x] Bild-Parität ([`ADR-0115`](../../adr/0115-backfill-spaltenwerte-text-ergebnisformat.md)
      Festlegung 4): über den Typ-Satz der Tabelle in `parity_types_test.go` — 85
      Typ-Spalten, darunter `bool`, `char(n)`, `inet`, `cidr`, `macaddr`, `money`,
      `bit`/`varbit`, `uuid`, `xml`, `oid`/`regclass`, Enum, vier Domains (über
      `bool`, `int`, `text`, `int[]`), zwei Composites, Range/Multirange,
      `tsvector`/`tsquery`, sieben Geometrie-Typen, Zeit-Typen, die Extension-Typen
      `hstore`, `citext`, `ltree`, 19 Array-Typen (darunter `hstore[]`), dazu
      eine `NULL`-Zeile, eine Rand-Zeile, eine gelöschte und eine generierte Spalte
      — ist der Roh-Text des Backfill-Pfads je Spalte byte-gleich dem des
      WAL-Pfads und das über `model.BuildRowImage` gebaute Bild ebenso, unter drei
      Rollen-GUC-Lagen (Standard; `timezone`/`datestyle`/`intervalstyle`/
      `bytea_output`/`extra_float_digits`; `postgres_verbose`/`SQL, MDY`/
      `America/New_York`) auf PostgreSQL 17 und 18 (M5). Der Test legt die drei
      Extensions in einem eigenen Schema an und entfernt sie im Cleanup; ein Test-Guard
      (`requiredTypes`) meldet einen aus der Tabelle gestrichenen Pflicht-Typ. *Zu
      belegen durch:* `make test-replication` (PostgreSQL 18 und
      PostgreSQL 17 über `PG_TEST_IMAGE`; der WAL-Zweig des Vergleichs läuft im Tier
      über Publication und Slot) und, netzlos, `make test` (die Cursor-Anweisung
      trägt an einfachen Namen kein `::` und keine Klammer).
- [x] Gate-Zuordnung (a): das neue Paket steht in den namentlichen Stellen des
      Messgegenstands — `Dockerfile` (Ausschluss-Muster der Stufe `coverage`),
      `DB_COVERAGE_PKGS` in `tools/harness/db-coverage.sh` und die Paketliste der
      Messphase in `tools/harness/run-replication-tests.sh` — sowie in den
      Beschreibungen `harness/sensors/coverage-gate.md` und
      `harness/sensors/db-adapter-coverage.md` und, falls sie eine Paketliste
      nennt, in `harness/README.md` ([`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md) Punkt 1: die Regel greift ohne
      Textänderung, Trigger (a) verlangt die Liste nachzuziehen). *Zu belegen
      durch:* der Suchlauf in §3 (beide Stände) und ein grüner `make coverage-gate`
      am Diff — dessen Vorlauf `tools/harness/db-package-lists-check.sh` hält die
      vier namentlichen Stellen (Dockerfile-Filter, `DB_COVERAGE_PKGS`, die
      Paketlisten beider Messläufe) gleich.
- [x] Gate-Zuordnung (b): die Tests des neuen Pakets überspringen ohne Datenbank
      **tatsächlich** — das trägt die Ausnahme. *Zu belegen durch:* ein `go test -v`
      des Pakets ohne `CDC_REPLICATION_TEST_DSN` (Docker-only), dessen Ausgabe für
      **jeden** Test `SKIP` zeigt; der Bericht nennt Befehl und gedruckte Ausgabe.
- [x] Gate-Zuordnung (c): netzlos prüfbare Logik (Konstruktionsvalidierung,
      Slot-Name, Cursor- und Fetch-Anweisung, Bezeichner-Quoting, Wertübernahme,
      Schätzungs-Abbildung, Fehlerklassifikation) liegt **nicht** im ausgenommenen
      Paket, sondern im Unterpaket `snapshotlogic` und damit im Gegenstand des
      blockierenden Gates; die Bild-Konstruktion liegt in der Domäne
      (`model.BuildRowImage`), die Blockbildung ist `FETCH FORWARD B`. Im
      ausgenommenen Paket bleiben die Schritte mit Verbindung. *Zu belegen durch:*
      Review des Diffs und der Statement-Zahl des Gate-Laufs vor und nach dem Diff
      (Zähler und Nenner mit Lauf, [`AGENTS.md`](../../../../AGENTS.md) §3.12).
      *Beleg:* Nenner der `coverage`-Stufe 2040 → 2082, gedeckt 1691 → 1731, das
      Unterpaket trägt 42 von 42 (Läufe in §3 Messbefunde); der END-verankerte
      Filter nimmt es nicht aus.
- [x] Gate-Zuordnung (d): der Adapter läuft in der Tier `make test-replication`
      (Replikationsverbindung); die DB-Adapter-Coverage wird mit dem Slice
      nachgemessen und der Bericht nennt ihre Zahl mit Zähler, Nenner und Lauf
      gegen die geltende Schwelle (`DB_COVERAGE_THRESHOLD`, [`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md) Punkt 3).
      *Zu belegen durch:* `make test-replication` (Exit ungefiltert gesichert).
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein Self-Review (Modul 8).
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [x] Doku-Update: `harness/README.md` §Sensors (`make test-replication` und die Coverage-Zeilen) und die Sensor-Dateien `harness/sensors/coverage-gate.md`/`harness/sensors/db-adapter-coverage.md` tragen das neue Paket, soweit sie Pakete nennen (Suchlauf in §3).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel ·
      neuer Sensor · benannte Spec-Lücke).
- [ ] Reconciliation-Register — entfällt: keine Reconciliation-Datei in
      diesem Repo (Greenfield).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls
      eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen /
      weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      von der Closure der Welle [welle-backfill-bestand](../welle-backfill-bestand.md) (die Roadmap führt sie unter
      *Offene Wellen*, das Ereignis kann eintreten).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/application/port/outbound/tablesnapshot.go` | neu | `TableSnapshotPort` (`OpenSnapshot`, `EstimatedRows`) und `TableSnapshot` (`Offset`, `Columns`, `NextBlock`, `Close`) samt den fünf Fehlerklassen-Sentinels (`permission`, `configuration`, `transient`, `replication`, `storage`); Fähigkeits-Port ([`ADR-0034`](../../adr/0034-ports-nach-faehigkeiten.md)). Die Schätzung liegt im selben Port (§1: Katalog-Fähigkeit derselben Quelle). |
| `internal/adapters/driven/postgressnapshot/snapshot.go` | neu | neues eigenes Paket: Slot-Anlage, Import, Cursor, Katalog-Lesen von Spaltenliste und Schätzung — die Schritte mit Verbindung; netzlos prüfbare Logik liegt außerhalb (`snapshotlogic`, Domäne bzw. `model.BuildRowImage`). Blockgröße `B` = 1.000 (`DefaultBlockSize`) und Zeitlimit der Slot-Anlage 30 s (`DefaultSlotTimeout`) sind Startwerte (Setzung ohne Messung), beide über Optionen übersteuerbar. Der Adapter parst den DSN einmal in `New` und kopiert die Konfiguration je Verbindung; ein beendeter Lese-Zustand endet als `transient`. |
| `internal/adapters/driven/postgressnapshot/snapshotlogic/logic.go`, `logic_test.go` | neu (Fixrunde) | Unterpaket mit der netzlos prüfbaren Logik: `Validate`, `SlotName`, `ValidSnapshotName`, `QuoteIdent`, `CursorStatement` (nur gequotete Namen, kein Cast, kein Funktionsaufruf — [`ADR-0115`](../../adr/0115-backfill-spaltenwerte-text-ergebnisformat.md) Festlegung 1), `FetchStatement`, `Values`, `Estimate`, `Classify` (samt `57P01`/`57P02`/`57P03` als `transient`); netzlose Tests laufen in `make test`, das Paket liegt im Unit-Gegenstand (DoD (c), Review-Befund F-5/F-6). |
| `internal/adapters/driven/postgressnapshot/seam.go` (Treiber-Hülle nach [`ADR-0080`](../../adr/0080-nahtform-pgconn-adapter-treiberhuelle.md)) | **nicht realisiert** | Wahl, keine Pflicht: [`ADR-0080`](../../adr/0080-nahtform-pgconn-adapter-treiberhuelle.md) Festlegung 1/5 gilt für `postgresack` und `receive`, Festlegung 4 lässt netzlose Fake-Tests in ausgenommenen Paketen zu — der Gegenstand wächst dann um die Hüllen-Statements. Das Paket wählt den anderen Weg: die Schritte mit Verbindung hängen direkt an `*pgconn.PgConn` und laufen als reale Tier-Läufe, die netzlos prüfbare Logik liegt in `snapshotlogic`. |
| `internal/adapters/driven/postgressnapshot/snapshot_test.go`, `parity_types_test.go` | neu, Fixrunde und zweite Fixrunde erweitert | Snapshot-Träger M1–M4, Blöcke, Spalten, Schätzung (auch die Nullgrenze `reltuples = 0`), Klassen, Bild-Parität (M5) über die Typ-Tabelle aus `parity_types_test.go` in drei GUC-Lagen (Roh-Text je Spalte und Bild; die Fehlermeldung nennt Spalte und Typ); Quoting von Schema und Tabelle (`"Sch ""ema"`, `"Mixed Case ""x"`); Fehlerpfade: belegter Slot-Name (`replication`), unbekannter Snapshot beim Import (`storage`, beide Verbindungen danach geschlossen), Quelle ohne Listener (`transient`), Katalog-Abfragen ohne `SELECT` auf `pg_attribute`/`pg_class` in einer eigenen Datenbank (`permission`), beendete Sitzung mitten im Lesen (`pg_terminate_backend`, `transient`), beendeter Kontext (`transient`); **jeder** Test überspringt ohne `CDC_REPLICATION_TEST_DSN`. Der Test liegt im Paket selbst (`package postgressnapshot`), weil M1 die zwei Schritte `exportSnapshot`/`importSnapshot` einzeln fährt (eine Zeile wird zwischen Slot-Anlage und Import committet). Zusätzlich `TestSlotReserveExhaustedIsConfiguration` (`SQLSTATE 53400` im Fehlertext), das ohne `CDC_SNAPSHOT_TEST_EXCLUSIVE_DSN` überspringt; die Variable setzt die Phase `tier` von `tools/harness/run-replication-tests.sh` (eigener PostgreSQL mit `max_replication_slots=1`). Zweite Fixrunde (Nachprüfungs-Review F-1, F-4): der Typ-Satz trägt 85 Typ-Spalten einschließlich `hstore`, `citext`, `ltree` und `hstore[]` — die drei Extensions liegen in einem eigenen Schema `snap_ext_<sfx>` (`CREATE EXTENSION … SCHEMA`, im Testcontainer als Superuser anlegbar) und `set.drop` entfernt sie samt Schema; `requiredTypes` und `missingRequired` melden einen aus der Tabelle gestrichenen Pflicht-Typ; die Rolle von `TestCatalogQueryFailuresKeepTheirClass` wird vor der eigenen Datenbank angelegt, damit der Cleanup die Datenbank vor der Rolle entfernt; die Replication-Verbindung eines Exports wird per `t.Cleanup` geschlossen. |
| `internal/adapters/driving/replication/receive/stream_test.go` | update (zweite Fixrunde, Nachprüfungs-Review F-2) | `TestStreamRestartsOnExistingSlot` wartet nach dem Ende des ersten Laufs auf `pg_replication_slots.active = false` (`awaitSlotInactive`, Zeitgrenze 10 s), bevor der Restart den Slot wieder auflegt; nur der Test ändert sich, der Produktionscode des Pakets `receive` bleibt. Sachlicher Befund (Beobachtungs-Register: Sache der Closure): der Adapter wiederholt `START_REPLICATION` bei SQLSTATE 55006 („slot … is active for PID") nicht — ein Restart, der den Slot vor der serverseitigen Freigabe des vorigen Walsenders auflegt, endet als Fehler; im Betrieb endet der Prozess auf jeden Adapter-Fehler mit Ausgang 1, gemeldet als Fehlerklasse `replication` (`internal/bootstrap/wiring.go`, Kommentar an `Run`), und die Fortsetzung trägt der Prozess-Neustart; im Test ist es ein Fehlschlag. |
| `tools/harness/run-replication-tests.sh` | update, Fixrunde erweitert | Paketliste der Messphase um das neue Paket; Kommentare der Messphase. Die Phase `tier` fährt nach `go test ./...` den Slot-Reserve-Test gegen einen eigenen PostgreSQL mit `max_replication_slots=1` (Review-Befund F-10): ein Lauf ohne `--- PASS` des Tests ist rot. Der zweite Lauf gegen PostgreSQL 17 ist `PG_TEST_IMAGE=<17er-Digest> make test-replication`. Zweite Fixrunde (Nachprüfungs-Review F-3): der Kommentar an der Phase `tier` nennt beide Läufe (Tier-weiter `go test ./...` und Slot-Reserve-Lauf). |
| `.github/workflows/e2e.yml` | update (nur Beschreibungstext, zweite Fixrunde, Nachprüfungs-Review F-3/F-9) | der Kommentar vor den Replication-Schritten und der `name:` des Schritts `tier` nennen den Slot-Reserve-Lauf; **kein** Schritt, keine Matrix, keine `uses:`-Zeile, kein `run:` ändert sich ([`AGENTS.md`](../../../../AGENTS.md) §3.8/§3.10 unberührt). **Wirkung ohne Workflow-Änderung:** der Schritt `tier` ruft `run-replication-tests.sh tier` und fährt damit je Matrix-Leg (PostgreSQL 17 und 18) einen zweiten Container (`max_replication_slots=1`) und den Slot-Reserve-Test; der Workflow ist nicht blockierend. Die Zeitannahmen der Phase mit Ursprung stehen in den Messbefunden. Ob der Lauf auf dem GitHub-Runner ebenso grün läuft, ist bis zum ersten realen Lauf ungeprüft (§6). |
| `tools/harness/db-package-lists-check.sh`, `harness/mk/coverage.mk` | neu, update (Fixrunde) | netzloses Prüfskript, das Dockerfile-Filter, `DB_COVERAGE_PKGS` und die Paketlisten der beiden Messläufe vergleicht; `coverage-gate` ruft es vor dem Bau — eine Verschärfung des bestehenden Gate-Ziels, keine Lockerung (`AGENTS.md` §3.6). Entscheidung zu Review-Befund F-3: das kleinste Mittel ist ein Vergleich, kein Umbau der Listen auf eine Quelle. Grenze: die Gleichheit ist gewächtert, die Eigenschaft und die Skip-Eigenschaft (F-11) bleiben Disziplin, benannt in `db-adapter-coverage.md` §Grenze Nr. 8. |
| `Dockerfile` (Stufe `coverage`), `tools/harness/db-coverage.sh` (`DB_COVERAGE_PKGS`) | update | das Paket steht im Ausschluss-Muster bzw. in der Paketliste ([`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md) Punkt 1): es ist ohne DB-Verbindung ungeprüft und füllt den netzlosen Nenner nicht; Kopfkommentare von `db-coverage.sh` und Dockerfile-Stufe zählen vier statt drei Pakete. Der Kopfkommentar von `db-coverage.sh` nennt die Positionsvielfachheit ohne Zahl und verweist für Zahlen auf das Sensor-Dokument (Review-Befund F-2/F-8). |
| `harness/sensors/coverage-gate.md`, `harness/sensors/db-adapter-coverage.md`, `harness/README.md` | update | Paketlisten, Nenner **813** mit Lauf, Positionsvielfachheit 3 im Replication-Profil, Rücknahme-Messung (`postgressnapshot` färbt die Stufe rot), Grenze Nr. 7 (geteilter Testcontainer, Reserve-Test in der Phase `tier`) und Nr. 8 (Listengleichheit gewächtert, Eigenschaft und Skip-Eigenschaft Disziplin); `harness/README.md` §Sensors verweist in der Zeile `make test-replication` für die Paketliste auf `db-adapter-coverage.md` §Gegenstand (Review-Befund F-12) und nennt in der Zeile `make coverage-gate` das Prüfskript. Zweite Fixrunde (Nachprüfungs-Review F-5): `coverage-gate.md` §Zählbasis nennt die Zusammensetzung des Nenners statt eines Vorgänger-Stands; die Skip-Aussage in `db-adapter-coverage.md` §Grenze Nr. 8 nennt den Slot-Reserve-Test mit seiner eigenen Vorbedingungs-Prüfung; Grenze Nr. 7 nennt die Wirkung der Last auf den Restart-Test und den Poll, der sie im Test abfängt. |
| `internal/application/port/outbound/tablesnapshot.go`, `internal/adapters/driven/postgressnapshot/snapshot.go` | update (zweite Fixrunde, Nachprüfungs-Review F-8) | Port-Doku und Kommentar an `DefaultBlockSize` benennen die Grenze `B` als Zeilenzahl, nicht Bytezahl; der Block liegt bis zur Rückgabe als Treiberbytes und als Zeichenketten vor (abgeleitet aus dem Treiber-Quelltext, nicht gemessen). Kein Code ändert sich. |
| `tools/schema/nacharbeit-roles.sql`, `docs/user/benutzerhandbuch.md` | **nicht berührt** | `SELECT` auf die Quelltabelle ist Betriebs-Vorbedingung ([`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 1, keine Vergabe durch die Software), `REPLICATION` trägt `cdc_capture` bereits: kein Grant, kein Schema-Objekt, damit kein Alt-Tag-Lauf ([`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md)). Keine Betreiber-Oberfläche (keine `CDC_*`-Variable, keine `cdc.*`-Funktion, kein Endpunkt): kein Handbuch-Zug, das Handbuch der Auslösung trägt `slice-backfill-sql-administration`. |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „die Menge der Pakete, deren Testlauf einen externen PostgreSQL voraussetzt"; beide Stände gemessen; der Befund ist der **Planner-Vorbefund am Parent-Stand `18cee124`**, der Implementer bestätigt ihn am Diff und trägt Nichtgefundenes nach):**

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| Ausschluss-Muster der Coverage-Stufe | `grep -n 'postgresstorage\|postgresack\|replication/receive' Dockerfile` | `Dockerfile:101`: `grep -vE '(^|/)(postgresstorage|postgresack|replication/receive)$'` | `postgressnapshot` ergänzen |
| Paketliste der DB-Adapter-Coverage | `grep -n 'DB_COVERAGE_PKGS' tools/harness/db-coverage.sh` | `tools/harness/db-coverage.sh:48` (drei Pakete); Kopfkommentar Zeilen 7–33 nennt die Pakete und Statement-Zahlen des Stands `slice-085` | Liste ergänzen; Kommentar nur ergänzen, wo er aufzählt, Zahlen mit Ursprung neu messen |
| Paketliste der Messphase im Tier-Skript | `grep -n 'postgresack' tools/harness/run-replication-tests.sh` | `tools/harness/run-replication-tests.sh:98-103`: der `go test`-Aufruf der Messphase nennt `postgresack` und `replication/receive` namentlich (dritte namentliche Stelle neben `Dockerfile` und `db-coverage.sh`); die Kommentare in den Zeilen 87–88 und 116–117 nennen die Pakete ebenfalls | Paket ergänzen, sonst läuft der Adapter in der Messphase nicht mit |
| Sensor-Beschreibungen | `grep -rn 'postgresack' harness docs` | `harness/sensors/db-adapter-coverage.md` (Zeilen 28, 45–68: Paketliste, Nenner „691 Statements"), `harness/sensors/coverage-gate.md` (Zeilen 20, 131, 194–217: Paket-Zahlen des Stands `slice-085`); kein Treffer in `harness/README.md` (nennt keine Paketliste) | aktive Träger nachziehen; `Accepted` ADRs ([`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md), [`ADR-0080`](../../adr/0080-nahtform-pgconn-adapter-treiberhuelle.md)) nicht ändern |
| Träger der DB-Adapter-Coverage im Workflow und in Make | `grep -n 'postgresack\|replication/receive' .github/workflows/*.yml Makefile harness/mk/*.mk` | kein Treffer | keine Workflow-Änderung ([`AGENTS.md`](../../../../AGENTS.md) §3.10 greift nicht); der Implementer bestätigt am Diff |
| Build-Kontext | Lesen von `.dockerignore` | `.dockerignore:28`: `!internal/`, an Parent (`48aa388a`) und Diff (`641b2ba4`) identisch; das neue Paket liegt unter `internal/`, kein neuer Pfad | `internal/` ist freigegeben; ein Pfad außerhalb bräuchte eine Freigabe (`BEO-PGC/dockerignore-default-deny-blockiert-neuen-pfad`) |

**Bestätigung des Implementers, Erstlieferung (beide Stände gemessen; Parent-Stand `48aa388a`, Diff-Stand `641b2ba4`):**
`git grep -c -E 'postgresack|replication/receive|postgressnapshot' <Stand> -- . ':!docs/plan/adr' ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!internal' ':!docs/plan/planning/in-progress/slice-backfill-snapshot-reader.md'`
druckt am **Parent** (Treffer je Datei): `Dockerfile` 1 · `slice-transformationen-antragsweg-usecase.md` 1 · `slice-transformationen-kern-rename.md` 1 · `harness/sensors/coverage-gate.md` 11 · `harness/sensors/db-adapter-coverage.md` 8 · `tools/harness/db-coverage.sh` 8 · `tools/harness/run-replication-tests.sh` 6; am **Diff**: dieselben Dateien (`coverage-gate.md` 16, `db-adapter-coverage.md` 11, `db-coverage.sh` 7, `run-replication-tests.sh` 7) plus `harness/README.md` 1 und `welle-backfill-bestand.md` 1 (dessen Erwähnung von `postgressnapshot` ist Bestand der Welle, keine Nachzug-Stelle).

- **Gefunden und nachgezogen:** die drei namentlichen Stellen (`Dockerfile`, `db-coverage.sh`, `run-replication-tests.sh`) und beide Sensor-Beschreibungen. Zahlenträger: der Nenner 691 → **848** mit Lauf, die Positionsvielfachheit im Replication-Profil zwei → **drei**, die Beispielzahlen des Replication-Profils durch die Zahlen des Diff-Laufs ersetzt; die datierte Rückrechnung von `coverage-gate.md` Grenze 4 (Stand `slice-085`, drei Pakete) bleibt unverändert gültig und trägt den vierten Fall als gemessen.
- **Nicht gefunden:** kein Treffer in `.github/workflows/*.yml`, `Makefile`, `harness/mk/*.mk` (beide Stände: `git grep -c` über diese Pfade nennt keine Datei) — keine Workflow-Änderung, [`AGENTS.md`](../../../../AGENTS.md) §3.10 greift nicht; kein Treffer in `docs/user/`, `spec/`, `README.md`, `AGENTS.md`; `harness/README.md` trug an beiden Ständen **keine** Paketliste (die ergänzte Nennung in der Zeile `make test-replication` ist ein Zusatz, kein Nachzug).
- **Gefunden, in fremden Dateien:** `docs/plan/planning/open/slice-transformationen-kern-rename.md` (Zeile 119) und `docs/plan/planning/open/slice-transformationen-antragsweg-usecase.md` (Zeile 102) behaupteten, die `coverage`-Stufe des Dockerfile schließe „nur die Pakete `postgresstorage`, `postgresack` und `replication/receive`" aus — am Diff-Stand sind es vier. Die Fixrunde zieht beide Sätze minimal auf vier Pakete nach (Review-Befund F-16; die Pläne bleiben offen, ihr Planner prüft sie beim Start).
- **Vorbestehende Drift:** `harness/sensors/coverage-gate.md` §Zählbasis nannte „Der aktuelle Nenner ist 1936" (Lauf `slice-097`); der Nenner der `coverage`-Stufe misst am Parent wie an der Erstlieferung **2040** (abgeleitet aus dem Profil der Stufe, dedupliziert über die Block-Position). Die Fixrunde führt den geltenden Nenner mit Lauf (siehe die Fixrunden-Zeilen unten) und datiert die 1936 auf `slice-097`.
- **Build-Kontext:** kein neuer Pfad, siehe die Zeile oben.

**Suchlauf der Fixrunde** (bewegte Eigenschaften: der Wortlaut des Lesepfads — Cast statt Ausgabefunktion —, die Menge der ausgenommenen Pakete samt Unterpaket `snapshotlogic`, der Träger der Slot-Reserve-Zusage; Parent-Stand `d7539e2c`, Diff-Stand `bf05067d`):

- **Wortlaut des Lesepfads:** `grep -rn --include=*.md -e '::tex[t]' spec docs harness AGENTS.md README.md` (die Klammer hält diese Zeile aus dem Treffer-Satz). Am Parent des Vorgängers `5e5d0b89` (`git grep -n` mit denselben Pfaden) **9** Zeilen in **4** Dateien: `ADR-0111` 2, `done/slice-backfill-row-image-gemeinsam.md` 1, dieser Plan 1, der Review-Report 5. Am Diff-Stand `bf05067d` **21** Zeilen in **4** Dateien: `ADR-0111` 2, `ADR-0115` 13, `done/slice-backfill-row-image-gemeinsam.md` 1, Review-Report 5; dieser Plan 0. **Gefunden und nachgezogen:** dieser Plan (Ziel-Absatz, DoD Bild-Parität, Messbefunde, Risiken). **Nicht angefasst, mit Grund:** `ADR-0111` (`Accepted`, ersetzt an drei Stellen durch `ADR-0115`), der `done/`-Record und der Review-Report (Records, `AGENTS.md` §3.5), `ADR-0115` (trägt den Wortlaut als Gegenstand seiner Entscheidung). Kein Treffer in `spec/`, `harness/`, `AGENTS.md`, `README.md`.
- **Menge der ausgenommenen Pakete (Träger mit Paketnamen):** `git grep -c -E 'postgresack|replication/receive|postgressnapshot' <Stand> -- . ':!docs/plan/adr' ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!internal' ':!docs/plan/planning/in-progress/slice-backfill-snapshot-reader.md'` (Diff-Stand mit `--untracked` gemessen). Am Parent `d7539e2c`: `Dockerfile` 1 · `slice-transformationen-antragsweg-usecase.md` 1 · `slice-transformationen-kern-rename.md` 1 · `welle-backfill-bestand.md` 1 · `harness/README.md` 1 · `coverage-gate.md` 16 · `db-adapter-coverage.md` 11 · `db-coverage.sh` 7 · `run-replication-tests.sh` 7. Am Diff: dieselben Dateien mit `kern-rename` 2, `coverage-gate.md` 20, `db-adapter-coverage.md` 14, `db-coverage.sh` 4, `run-replication-tests.sh` 8. **Nachgezogen:** beide offenen Pläne (vier Pakete), `harness/README.md` Zeile `make test-replication` (Verweis statt vierter Paketliste), beide Sensor-Beschreibungen. **Nicht gefunden:** kein Treffer für die vier Pakete oder `snapshotlogic` in `.github/workflows/*.yml`, `Makefile`, `harness/mk/*.mk` (beide Stände) — keine Workflow-Änderung, [`AGENTS.md`](../../../../AGENTS.md) §3.10 greift nicht.
- **Träger der Slot-Reserve-Zusage:** `git grep -n --untracked -E 'SNAPSHOT_TEST_EXCLUSIVE_DSN|TestSlotReserve'`: Test, `run-replication-tests.sh` (setzt die Variable, Phase `tier`), `db-adapter-coverage.md` §Grenze Nr. 7 und dieser Plan; die vorige Behauptung „nur manuell belegt" steht in keinem Träger mehr.
- **Träger der Skip-Eigenschaft und der Listengleichheit:** `db-adapter-coverage.md` §Grenze Nr. 8, `coverage-gate.md` §Grenze Nr. 8, `harness/README.md` Zeile `make coverage-gate`.
- **Suchlauf der zweiten Fixrunde** (bewegte Eigenschaften: der Inhalt der Phase `tier` — zwei Läufe statt eines — und der Typ-Satz mit seiner Größe und Ausschluss-Begründung; Parent-Stand `56305a4d`, Diff-Stand `d3420bd9`):
  - **Phasen- und Beschreibungstexte:** `git grep -c -i -w -E 'tier|Messphase|Slot-Reserve|max_replication_slots' <Stand> -- .github tools/harness harness docs/user docs/plan/planning/open docs/plan/planning/in-progress ':!docs/plan/planning/in-progress/slice-backfill-snapshot-reader.md' README.md AGENTS.md Makefile compose.yaml`. Am Parent 17 Dateien, am Diff dieselben 17 mit drei geänderten Zählern: `e2e.yml` 6 → 7, `db-adapter-coverage.md` 15 → 19, `run-replication-tests.sh` 10 → 11; unverändert `compose.yaml` 1 · `harness/README.md` 2 · `coverage-gate.md` 1 · `generated-sync.md` 1 · `db-package-lists-check.sh` 1 · `run-integration-tests.sh` 2 · `run-store-tests.sh` 1 · die sieben offenen Pläne (`slice-backfill-e2e` 3, `-run-store` 2, `-sql-administration` 3, `slice-transformationen-antragsweg-schema` 2, `-antragsweg-usecase` 1, `-backfill-pfad` 3, `-e2e-wirkung` 1). **Gefunden und nachgezogen:** der Kommentar an der Phase `tier` in `tools/harness/run-replication-tests.sh` (Zeile 111 am Parent: „der Tier-weite `go test ./...`; der Exit dieses Aufrufs ist das Verdikt dieser Phase") und der Kommentar sowie der Schrittname der Phase in `.github/workflows/e2e.yml` (Zeilen 109 und 122 am Parent); `db-adapter-coverage.md` §Grenze Nr. 7 (Wirkung der Last auf den Restart-Test, Zeit der Phase). **Gelesen, ohne Nachzug:** `compose.yaml:30` (`max_replication_slots=10` der Betriebsumgebung, keine Beschreibung der Phase), die offenen Pläne (nennen `max_replication_slots` als Betriebs-Vorbedingung oder „Tier" als Testebene, keine Beschreibung der Phase `tier`), `harness/README.md` (Zeile `make test-replication` nennt den Slot-Reserve-Test bereits; Zeile `make test-integration` betrifft den Compose-Tier). **Nicht gefunden:** kein Treffer in `docs/user/`, `README.md`, `AGENTS.md`, `Makefile`.
  - **Typ-Satz:** `git grep -n -E '81 Typ|18 Array|nicht anlegbar|Nicht anlegbar|bringt sie nicht mit|nicht mitbringt|Extension-Typen' <Stand> -- . ':!docs/reviews'`. Am Parent: dieser Plan (DoD Bild-Parität, Messbefund zur Parität), `parity_types_test.go:28`, `ADR-0115` §Kontext (nennt die Typen als Bestandteil der Probe, `Accepted`); am Diff: alle drei Plan-/Test-Stellen nachgezogen, `ADR-0115` unverändert. Zwei Treffer zu `nicht mitbringt` in `docs/plan/planning/observations/BEO-PGC/roter-test-ohne-leser/` betreffen ein anderes Fixture (Schema-Tabellen) und bleiben.

**Messbefunde des Implementers** (Läufe an diesem Diff; die Nachweise stehen mit Befehl im Bericht der Übergabe):

- **Unit-Zahl, Erstlieferung** (`make coverage-gate`, gedruckt): Parent `82.90%`, Diff `82.90%` (zwei weitere `make gates`-Läufe am Diff: `82.90%`, `82.80%`); gedeckt/Nenner, abgeleitet aus dem Profil der Stufe (dedupliziert): Parent **1692 von 2040**, Diff **1691 von 2040**.
- **Unit-Zahl, Fixrunde** (die Stufe nachgestellt: `go test -coverpkg=<Paketliste der Stufe> -covermode=atomic`, dedupliziert über die Block-Position): am Stand `d7539e2c` **1691 von 2040** (gedruckt `total: (statements) 82.9%`), am Diff **1731 von 2082** (gedruckt `83.1%`); `make coverage-gate` am Diff druckt `Coverage 83.20%`. Der Nenner wächst um **42** Statements — das Unterpaket `snapshotlogic` liegt im Gegenstand, gedeckt **42 von 42** (aus demselben Profil abgeleitet). Die gedeckte Zahl ist lauf-gebunden (zwei Läufe am Diff druckten `83.1%` und `83.2%`).
- **DB-Adapter-Coverage, Erstlieferung** (frischer `make test-store` gefolgt von `make test-replication`, PostgreSQL 18, gedruckt): Parent **77.13% (gedeckt 533 von 691)**, Diff **79.25% (gedeckt 672 von 848)** gegen `DB_COVERAGE_THRESHOLD` 70; `postgressnapshot` trug davon **139 von 157** (abgeleitet aus dem gemergten Profil).
- **DB-Adapter-Coverage, Fixrunde** (frischer `make test-store` gefolgt von `make test-replication`, PostgreSQL 18, gedruckt): **80.07% (gedeckt 651 von 813 Statements; Profile gemergt: store,replication)**; derselbe Wert druckt `PG_TEST_IMAGE=<17er-Digest> make test-replication` gegen PostgreSQL 17. Der Nenner sinkt um 35 (`postgressnapshot` 157 → 122, abgeleitet): die Logik des Unterpakets verlässt den Gegenstand, die Konfigurations-Kopie je Verbindung und der Zweig für beendete Verbindungen kommen hinzu; `postgressnapshot` trägt **118 von 122** (abgeleitet aus dem gemergten Profil). **Ungedeckt bleiben 2 Blöcke (4 Statements)**: die Zweige „`consistent_point` nicht lesbar" und „Snapshot-Name nicht lesbar" in `exportSnapshot` — der Server liefert beide Werte im geprüften Format, ein Fehlformat ist ohne Fake nicht auslösbar. Der Abbruch des komponierten Runs mitten im Lauf gehört dem `e2e`-Slice (§1). Die Zeitlage des Fensters zwischen Export und Cursor (gleichzeitiges `ALTER TABLE`, Review-Befund F-17) ist nicht gemessen und bleibt Sache des `e2e`-Slice.
- **`reltuples`:** `−1` für die frisch befüllte, nie analysierte Tabelle auf PostgreSQL **17.11** und **18.6** (der Test druckt den Wert und die Version), nach `ANALYZE` 40 (= tatsächliche Zeilenzahl) auf beiden — die Erwartung aus [`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 3, Punkt 5 ist belegt, die Abbildung „unbekannt" im Port bleibt.
- **GUC-Parität und Typ-Parität:** hält auf beiden Versionen — eine reguläre Verbindung mit demselben DSN trägt die Rollen-GUC wie der Walsender; Roh-Text je Spalte und Bild sind über den Typ-Satz (85 Typ-Spalten einschließlich `hstore`, `citext`, `ltree` und `hstore[]`; PostgreSQL 17.11 und 18.6, je drei GUC-Lagen; der Test druckt die Spaltenzahl) byte-gleich; der Fallback der expliziten `SET`-Angleichung entfällt. `CREATE EXTENSION hstore`/`citext`/`ltree` gelingt im Testcontainer an beiden Versionen (Testnutzer ist Superuser, Lauf `make test-replication` 17.11 und 18.6). Der Cast des Lesepfads verletzt die Parität für `bool`, `char(1)`, `char(n)`, `inet`, `xml` und die Domain über `bool` (gemessen im rot gesehenen Lauf, [`ADR-0115`](../../adr/0115-backfill-spaltenwerte-text-ergebnisformat.md)); die drei Extension-Typen bleiben unter dem Cast grün (die Mutation `::text` an jedem Spaltennamen färbt genau diese sechs Spalten rot), der Lesepfad ohne Cast trägt alle. **Abgleich mit dem Typ-Satz der ADR** (Zeile für Zeile, die Typliste aus §Kontext gegen `parity_types_test.go`): jeder dort genannte Typ hat eine Spalte, die 19 Array-Spalten tragen die 18 genannten Array-Formen samt Domain-über-`bool`-Array; **es fehlt nichts**. Nicht Teil der drei Lagen ist die Lage `lc_monetary='de_DE.utf8'`, die die Probe der ADR auf einem Debian-Image fuhr: das Alpine-Image des Tiers nimmt `SET lc_monetary='de_DE.utf8'` an, die `money`-Ausgabe bleibt `$1,234,567.89` (gemessen an `postgres:18-alpine`) — die Lage ist dort nicht herstellbar, und die Fitness Function der ADR nennt drei Lagen ohne sie.
- **Generierte Spalten unter PostgreSQL 18:** der Walsender sendet die generierte Spalte für die Publication dieses Repos nicht (Test auf 17.11 und 18.6: Spalte fehlt im WAL-Bild), das Backfill-Bild lässt sie ebenfalls aus.
- **Slot-Reserve:** `53400` endet als `configuration`; belegt durch die Phase `tier` von `tools/harness/run-replication-tests.sh` (eigener PostgreSQL mit `max_replication_slots=1`, `CDC_SNAPSHOT_TEST_EXCLUSIVE_DSN`; der Test prüft `SQLSTATE 53400` im Fehlertext, sein `--- PASS` ist Lauf-Bedingung).
- **Fehlerklasse einer beendeten Sitzung:** `pg_terminate_backend` auf die Lese-Sitzung endet als `57P01`; der Port führt „Verbindung abgebrochen" unter `transient`, die Erstlieferung klassifizierte es als `storage`. Die Fixrunde ordnet `57P01`/`57P02`/`57P03` und jeden Fehler auf beendeter Verbindung `transient` zu (`snapshotlogic.Classify`, `NextBlock`).
- **Verbindungslimit der Rolle:** eine Rolle mit `CONNECTION LIMIT 1` und `REPLICATION` öffnet den Snapshot trotzdem — die Replication-Verbindung zählt nicht gegen das Limit; `53300` an der Lese-Verbindung ist damit im Tier nicht auslösbar (gemessen; der Zweig steht als Unit-Fall in `TestClassify`), der Fehler der Lese-Verbindung ist über eine Quelle ohne Listener belegt.
- **Zweite Fixrunde, Läufe** (Ursprung: `make test-replication`, ungefiltert, Exit je Lauf gesichert, `tools/schema/plan.yaml` danach zurückgenommen): PostgreSQL 18.6 drei Läufe Exit 0, PostgreSQL 17.11 (`PG_TEST_IMAGE` auf den 17er-Digest aus `e2e.yml`) zwei Läufe Exit 0; jeder Lauf druckt `DB-Adapter-Coverage: 80.07% (gedeckt 651 von 813 Statements; Profile gemergt: store,replication)` (der Nenner ändert sich durch die Test-Änderungen nicht: sie liegen in Test-Dateien); der Test druckt `PostgreSQL 17.11, 85 Typ-Spalten` bzw. `PostgreSQL 18.6, 85 Typ-Spalten` (`-v`-Läufe gegen Wegwerf-Container mit der Tier-Konfiguration). Unit-Seite: `make coverage-gate` am Diff druckt `Coverage 83.20% erfüllt Schwelle 80%`, sein Vorlauf `db-package-lists-check: OK — … dieselben 4 Pakete`; `make test` und `make a-check` (`gesamt: 0 Befund(e)`) Exit 0; das Paket `postgressnapshot` überspringt ohne Datenbank 21 von 21 Tests (`go test -count=1 -v` mit `--network none`, gedruckt 21 `--- SKIP`, 0 `--- PASS`).
- **Flake `TestStreamRestartsOnExistingSlot` (SQLSTATE 55006):** Aufbau: Wegwerf-PostgreSQL mit der Tier-Konfiguration (`wal_level=logical`, `max_replication_slots=10`, `max_wal_senders=10`, `wal_sender_timeout=2000`), dieser Test mit `-count=300` parallel zum Paket `postgressnapshot` mit `-count=6` (Lastfenster des Pakets 46 bis 47 s, gedruckt). **Ohne** den Poll auf `active = false`: PostgreSQL 17 **11** und **9** von 300 Läufen rot, PostgreSQL 18 **10** von 300, jede Meldung `--- FAIL` mit 55006; **mit** dem Poll: PostgreSQL 17 **0** und **0**, PostgreSQL 18 **0** von 300. Das Paket `postgressnapshot` lief in allen Fenstern grün. Der Poll ändert nur den Test; das Produktionsverhalten (kein Wiederholungsversuch bei 55006) bleibt unverändert.
- **Laufzeit der Phase `tier`** (Ursprung: `bash -x` mit Zeitstempel je Kommando, ein Lauf je Version, Exit 0): die ganze Phase dauert 46,3 s (PostgreSQL 18) und 44,4 s (PostgreSQL 17); der Slot-Reserve-Lauf (Container-Start bis `--- PASS`, ohne den Modul-Cache-Schritt) 7,8 s und 7,5 s; der gedruckte Test selbst 0,08 s bis 0,21 s.
- **Aufräumen im Testcontainer** (Ursprung: zwei aufeinanderfolgende `go test -count=1`-Läufe des Pakets `postgressnapshot` gegen einen Wegwerf-PostgreSQL 18, danach `pg_roles`, `pg_database`, `pg_replication_slots`, `pg_extension`, `pg_namespace`): 0 Rollen `snap_role%`, 0 Datenbanken `snap_%`, 0 Slots, 0 Extensions `hstore`/`citext`/`ltree`, 0 Schemata `snap_ext%`. Mit der Rolle **nach** der Datenbank angelegt (die vorige Reihenfolge) bleiben nach denselben zwei Läufen 2 Rollen zurück.
- **Speicher eines Blocks:** die Grenze `B` zählt Zeilen, nicht Bytes; der Block liegt bis zur Rückgabe kurzzeitig doppelt vor (Treiberbytes und Zeichenketten der Werte) — aus dem Code und dem Treiber-Quelltext abgeleitet, nicht gemessen; benannt in der Port-Doku, Änderung Sache der Ausbaustufe für Durchsatz (`ADR-0111` Re-Evaluierung).

## 4. Trigger

**Start** (`next` → `in-progress`): wenn `run-usecase` **nicht** vorausgesetzt
ist (dieser Slice geht ihm voraus), `row-image-gemeinsam` in `done/` liegt und kein
anderer Slice in `in-progress/` liegt. Die Zuordnung des neuen Pakets zum
Messgegenstand ist entschieden und kein Start-Trigger: ein eigenes Paket
`postgressnapshot`, dessen Testlauf einen externen Dienst voraussetzt, ist nach
[`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md) Punkt 1 nicht Gegenstand des Coverage-Gates; die namentliche Liste
zieht dieser Slice nach (DoD, Gate-Zuordnung (a)).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls Snapshot-Träger
  und Bild-Parität nicht in einem Review tragen — der abtrennbare Teil ist die
  Bild-Parität (dann eigener Slice `backfill-bild-paritaet` mit Start nach
  diesem).
- `in-progress` → `open` (blockiert): falls die Bild-Parität real **nicht**
  besteht und die GUC-Angleichung (`SET` der Sitzungs-GUC auf die Werte des
  Walsenders) eine Entscheidung über den Betriebsvertrag braucht — Architect-Frage.

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + `make test-replication` real grün
(beide Phasen) + Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Der Test-Ausschluss trägt die Ausnahme nicht.** Das neue Paket ist nur dann
  aus dem Nenner des Coverage-Gates herausgehalten ([`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md) Punkt 1), wenn seine
  Tests ohne Datenbank tatsächlich überspringen; ein netzlos laufender Test im
  Paket wäre ein Widerspruch zum Ausschluss, ein vergessener Eintrag in einer der
  namentlichen Stellen (ihre Gleichheit hält `db-package-lists-check.sh`) füllte
  den Nenner mit ungedeckten Zeilen. *Erwartet,
  zu belegen durch:* DoD Gate-Zuordnung (a)–(c) und der grüne `make coverage-gate`
  am Diff. **Ausgang:** *(bei Closure)*
- **Die GUC-Parität ist ungemessen.** [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) belegt nur, dass der
  Walsender die Rollen-Defaults trägt (M5); dass eine reguläre Verbindung mit
  demselben DSN dieselben Werte trägt, führt die ADR als „erwartet".
  *Erwartet, zu belegen durch:* der Paritätstest über den Typ-Satz unter drei
  Rollen-GUC-Lagen (`ALTER ROLE … SET`, [`ADR-0115`](../../adr/0115-backfill-spaltenwerte-text-ergebnisformat.md) Festlegung 4).
  Fällt er rot aus, ist der Fallback eine explizite `SET`-Angleichung der
  Sitzung an die Walsender-Werte — dann Rückführung §4. **Ausgang:** *(bei
  Closure)*
- **Generierte Spalten in PostgreSQL 18.** Der WAL-Pfad sendet generierte
  Spalten nicht, solange die Publication sie nicht ausdrücklich veröffentlicht;
  ob das für die Publication dieses Repos unter PostgreSQL 18 gilt, ist ungemessen.
  *Erwartet, zu belegen durch:* der Paritätstest mit generierter Spalte auf 17
  und 18. **Ausgang:** *(bei Closure)*
- **Slot-Anlage wartet auf laufende Schreibtransaktionen** (M3: 5,27 s bei rund
  5 s Restlaufzeit, gemessen in [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md)) und braucht Reserve in
  `max_replication_slots`/`max_wal_senders`. *Erwartet, zu belegen durch:* der
  Zeitlimit-Test; die Tier-Umgebung setzt beides auf 10
  (`tools/harness/run-replication-tests.sh`). **Ausgang:** *(bei Closure)*
- **Netzlos geprüfter Code im DB-Gegenstand** (`BEO-PGC/db-gegenstand-enthaelt-netzlos-geprueften-code`,
  offen, 2×): netzlos prüfbare Logik im Adapter-Paket höbe die DB-Zahl ohne
  DB-Beleg und widerspräche dem Test-Ausschluss; sie liegt deshalb im Unterpaket
  `snapshotlogic` (DoD Gate-Zuordnung (c)). *Erwartet, zu belegen durch:* der Bericht nennt Zähler
  und Nenner je Messung; ein weiterer Beleg der Klasse erreicht die Schwelle 3×.
  **Ausgang:** *(bei Closure)*
- **`reltuples` = `−1` für nie analysierte Tabellen** ist Wissen aus der
  PostgreSQL-Dokumentation und in diesem Repo nicht gemessen ([`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md)
  Festlegung 3, Punkt 5). *Erwartet, zu belegen durch:* die Messung an
  PostgreSQL 17 und 18 im Test (DoD Blöcke und Spalten). Weicht ein Wert ab, ist
  die Abbildung „unbekannt" im Port zu ändern und der Befund geht als Frage an
  den Architect. **Ausgang:** *(bei Closure)*
- **`pglogrepl` trägt die Optionen** (`CreateReplicationSlotOptions{Temporary,
  SnapshotAction}`, Ergebnis `SnapshotName`) — in [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) als „gelesen"
  geführt, nicht als gefahren. *Erwartet, zu belegen durch:* M1 im Tier.
  **Ausgang:** *(bei Closure)*
- **Restart auf einem noch aktiven Slot.** Der Adapter im Paket `receive`
  wiederholt `START_REPLICATION` bei SQLSTATE 55006 („slot … is active for PID")
  nicht; ein Neustart des Streams, der den Slot vor der serverseitigen Freigabe
  des vorigen Walsenders auflegt, endet als Fehlerklasse `replication`. Im
  Test zeigt sich das als Flake unter Last (Flake-Messung in §3); der Test wartet
  auf die Freigabe; das Produktionsverhalten ist unverändert und
  gehört nicht zu diesem Slice. *Erwartet, zu belegen durch:* die Flake-Messung
  (0 von 300 in drei Läufen mit dem Poll). **Ausgang:** *(bei Closure)*
- **Blockgröße zählt Zeilen, nicht Bytes.** Bei Zeilen im MB-Bereich (`jsonb`,
  `bytea`) ist der Speicherbedarf eines `NextBlock` `B` mal die Zeilenbreite,
  und der Block liegt kurzzeitig doppelt vor (abgeleitet, nicht gemessen; Port-
  Doku benennt die Grenze). *Erwartet, zu belegen durch:* eine Messung mit
  breiten Zeilen in der Ausbaustufe für Durchsatz. **Ausgang:** *(bei Closure)*
- **Der Tier-Schritt in `e2e.yml` fährt den Slot-Reserve-Lauf je Matrix-Leg.**
  Ob der zweite Container (`max_replication_slots=1`, Bereitschaft über
  `pg_isready` im Container) auf dem GitHub-Runner ebenso grün läuft wie lokal
  (rund 7,5 bis 8 s je Lauf), ist bis zum ersten realen Lauf nicht geprüft
  ([`AGENTS.md`](../../../../AGENTS.md) §3.10: der Workflow ändert nur
  Beschreibungstext, das Verhalten des bestehenden Schritts umfasst den
  zweiten Container). *Erwartet, zu belegen durch:* der erste Post-Push-Lauf
  von `e2e.yml` (`gh run list`). **Ausgang:** *(bei Closure)*

## 7. Closure-Notiz

- **Was hat funktioniert:** *(zu tragen bei Closure)*
- **Was ging anders als geplant:** *(zu tragen bei Closure)*
- **Steering-Loop-Eintrag (Lerneintrag):** *(zu tragen bei Closure —
  geschärfte Regel · neuer Sensor · benannte Spec-Lücke; ohne ihn kein
  `done/`-Übergang)*
- **Beobachtungs-Register (`../observations/`):** *(je Anfall Beleg oder
  „keine Beobachtung angefallen" als notierte Antwort)*
- **Folge-Slices:** *(zu tragen bei Closure)*
- **Risiken aus §6:** *(je ein Ausgang)*
- **Drei Paarungen:** dieser Slice gehört zu [welle-backfill-bestand](../welle-backfill-bestand.md) (offen) — die
  Prüfung läuft regelkonform bei deren Closure.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area
`*`/`PGC` (Greenfield); der Snapshot-Adapter ist ein Driven Adapter wie die
bestehenden DB-Pakete, keine eigene Sub-Area — kein Anlass zur
Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/vorab-bedingung-nach-umsetzung-geprueft` (offen, 2×, gesichtet — die
Zuordnung des Pakets ist vor dem Start entschieden), `BEO-PGC/db-gegenstand-enthaelt-netzlos-geprueften-code`
(offen, 2×, einschlägig — Risiko §6),
`BEO-PGC/endstufe-unter-eigenem-messgegenstand-unerreichbar` (verkörpert),
`BEO-PGC/coverage-stage-dockerignore-blockiert-tooling` (offen, 2×) und
`BEO-PGC/dockerignore-default-deny-blockiert-neuen-pfad` (offen, 1×) —
Suchlauf-Zeile Build-Kontext, `BEO-PGC/arbeit-ueberholt-stehenden-traeger`
(verkörpert, 26×, Suchlauf §3), `BEO-PGC/gate-scope-erweiterung-ohne-adr-traeger`
(offen, 2×, einschlägig — der Träger der Erweiterung ist [`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md) Punkt 1 samt
Trigger (a), keine neue ADR), `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`
(verkörpert, 13×, Zahlen der Coverage tragen den Lauf).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
