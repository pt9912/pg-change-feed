# ADR-0118: Backfill — Umschreiben der Tabelle zwischen Snapshot-Export und Lesesperre wird erkannt und beendet den Run (Supersedes ADR-0111, Teilfrage 1 Ablauf teilweise)

**Status:** Accepted — Supersedes [`ADR-0111`](0111-backfill-bestand-snapshot-bulk-copy.md)
in **Teilfrage 1, Absatz „Ablauf eines Runs", Schritt 3, teilweise** (genau eine
Ergänzung: die Schritte, die dem Lesen der Spaltenliste vorangehen, siehe
§Entscheidung); alles Übrige von `ADR-0111` bleibt in Kraft, insbesondere
Teilfrage 3 (Lückenfreiheit gegenüber dem WAL-Pfad) und Teilfrage 5
(Fehlerklassen, run-lokal). [`ADR-0117`](0117-backfill-run-fehlerklasse-schema.md)
bleibt unberührt.

**Datum:** 2026-09-24

**Autor:** pt9912 (Architect-Rolle, Modul 8; ausgelöst durch den Befund „Tabellen-Rewrite
im Fenster" aus dem Ansatz-Ergebnis von `slice-backfill-e2e`, den diese ADR an
Wegwerf-PostgreSQL-Instanzen mit den Snapshot-Schritten des Adapters nachgestellt
und in seiner Reichweite gemessen hat)

**Bezug:** [`LH-FA-CAP-009`](../../../spec/lastenheft.md) (Haupt-Bezug — Happy
Path: jede zum Startzeitpunkt vorhandene Zeile ist lesbar; Negative: kein stiller
Verlust), [`LH-FA-DAT-003`](../../../spec/lastenheft.md) (nur zur Abgrenzung:
TRUNCATE), [`ADR-0111`](0111-backfill-bestand-snapshot-bulk-copy.md) (teilweise
superseded — Haupt-Bezug), [`ADR-0117`](0117-backfill-run-fehlerklasse-schema.md)
(Klassenmenge des Runs, nur zur Abgrenzung), [`ADR-0116`](0116-backfill-schema-version-referenz-reichweite.md)
(kompatible Spalten-Erweiterung, nur zur Abgrenzung),
[ADR-0023](0023-fehlerklassifikation.md) (Fehlerklassen)

**Schärft:** [`LH-FA-CAP-009.a`](../../../spec/pflichtenheft.md) (Absatz
„Mechanismus": ein Satz, siehe Folgepflicht 2), [`SPEC-008`](../../../spec/pflichtenheft.md)
(Klasse `transient`, unverändert wahr)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

[`ADR-0111`](0111-backfill-bestand-snapshot-bulk-copy.md) ist `Accepted` und
unberührbar (`AGENTS.md` §3.5). Er begründet die Vollständigkeit des Bestands über
den Slot-Snapshot (Teilfrage 1, Teilfrage 3) und führt im Schritt 3 des Ablaufs die
Spaltenliste und den Cursor „im Snapshot". Die Ergänzung ändert den Referenten
(den Ablauf) und ist keine Zitat-Korrektur.

### Befunde am Bestand (Code-Anker als Symbolnamen, Stand `26ea0c2e`)

- **Der Snapshot entsteht vor jeder Sperre der Tabelle.** `exportSnapshot`
  (`internal/adapters/driven/postgressnapshot/snapshot.go`) legt den temporären
  Slot an und liefert `snapshot_name`; `importSnapshot` öffnet danach eine
  `REPEATABLE READ`-Transaktion, importiert den Snapshot, liest die Spaltenliste
  (`readColumns`) und führt `DECLARE` des Cursors aus. Die `ACCESS SHARE`-Sperre
  der Tabelle entsteht erst mit dem `DECLARE`; zwischen Export und `DECLARE`
  liegt ein Fenster, in dem ein Fremdschreiber die Tabelle umschreiben kann.
- **Der Katalog des Cursors ist der aktuelle, der Bestand der des Snapshots.**
  `DECLARE` plant und öffnet die Relation über den aktuellen Katalog (die
  Relation trägt ihre neue Datei), die Zeilen filtert der importierte Snapshot.
- **Die Fehler des Fensters kennt der Run nur teilweise.** `DROP COLUMN` im
  Fenster endet am `DECLARE` mit SQLSTATE `42703`, Klasse `storage` (E2E-belegt in
  `slice-backfill-e2e`, Phase DDL-Fenster); ein Umschreiben der Tabelle endet
  ohne Fehler.

### Gemessen (Wegwerf-Läufe, nicht committet)

**Probe** (Shell-Skript im Scratchpad, 2026-09-24): Wegwerf-Container mit
`wal_level=logical`, je Läuf **eine** Tabelle `t(id int primary key, name text)`
mit 3 Zeilen; PostgreSQL **18.6** (`postgres:18-alpine`, Digest von
`PG_TEST_IMAGE` im `Makefile`) und **17.11** (`postgres:17-alpine`, Digest der
E2E-Matrix in `.github/workflows/e2e.yml`). Ablauf je Variante: Sitzung 1
(`replication=database`) legt `CREATE_REPLICATION_SLOT … TEMPORARY LOGICAL
pgoutput (SNAPSHOT 'export')` an und bleibt offen; Sitzung 2 führt die DDL-Variante
aus; Sitzung 3 (`BEGIN ISOLATION LEVEL REPEATABLE READ READ ONLY; SET TRANSACTION
SNAPSHOT '<name>'; LOCK TABLE public.t IN ACCESS SHARE MODE;`) zählt die Zeilen
und liest `pg_class.relfilenode` (SQL-Lesung im Snapshot) neben
`pg_relation_filenode(oid)` (Lesung über den aktuellen Katalog); die Ausdrücke
„Zeilen" und „Filenode verschieden" sind die gedruckten Zeilen
`rows_in_snapshot=…` und `differs=…`. Container danach mit `docker rm -fv`
entfernt.

| DDL im Fenster | Zeilen im Snapshot (gedruckt) | Filenode verschieden (gedruckt) | Versionen |
|---|---|---|---|
| keine (Kontrolle) | 3 | `false` | 17.11, 18.6 |
| `ALTER COLUMN name TYPE varchar(64)` (Umschreiben) | **0** | `true` | 17.11, 18.6 |
| `VACUUM FULL t` | 3 | `true` | 18.6 |
| `CLUSTER t USING t_pkey` | 3 | `true` | 18.6 |
| `TRUNCATE t` | **0** | `true` | 18.6 |
| `ALTER TABLE t DROP COLUMN name` | 3 | `false` | 18.6 |
| `DROP TABLE t; CREATE TABLE t …` (neue Tabelle) | 0 | Zeile in `pg_class` im Snapshot **fehlt** (`row_found_in_pg_class=0`) | 18.6 |
| partitionierte Tabelle, keine DDL | 3 | `false` (Ausdruck mit `coalesce(pg_relation_filenode(oid), 0)`; ohne `coalesce` stünde `0` gegen `NULL`) | 18.6 |
| `ALTER … TYPE` in einer offenen Transaktion, die 6 s hält, während `LOCK` läuft | **0** | `true`, die Sperre wartete (`lock_wait_s=5`) | 17.11, 18.6 |

Die Zeile mit `LOCK TABLE` vor dem Zählen belegt zugleich: **die Sperre allein
schließt die Lücke nicht** — bei `ALTER … TYPE` bleibt die Zahl `0`, auch mit
gehaltener Sperre. Die Werte von `pg_class.reltuples` taugen nicht als
Plausibilität: vor dem Umschreiben `-1` (Tabelle nie analysiert), im Snapshot nach
dem Umschreiben weiter `-1`, im aktuellen Katalog danach `3` (18.6 und 17.11,
Zeile `reltuples=…`). Eine Rolle mit ausschließlich `SELECT` auf die Tabelle führt
`LOCK TABLE … IN ACCESS SHARE MODE` in einer `READ ONLY`-Transaktion aus; ohne
`SELECT` endet die Anweisung mit `permission denied for table` (17.11).

**Lesart.** Ein Umschreiben, das nach dem Snapshot committet, lässt den älteren
Snapshot die neu geschriebene Datei leer sehen: `ALTER … TYPE` (Umschreibung) und
`TRUNCATE` verlieren dem Snapshot die Zeilen (nach der PostgreSQL-Dokumentation
nicht MVCC-sicher); `VACUUM FULL` und `CLUSTER` behalten sie (gemessen), ändern
die Datei aber ebenso. Die Lesung von `pg_class` im Snapshot zeigt den Eintrag
**zum Snapshot**, `pg_relation_filenode` den **aktuellen** — der Unterschied ist
das Merkmal des Umschreibens.

**Nicht gemessen:** die Dauer des Fensters (Export bis
Sperre); die Wirkung der Sperr-Warteschlange auf andere Leser der Quelltabelle,
solange eine DDL auf die Sperre des Runs wartet; die Datei-Änderungen durch
`ALTER TABLE … SET TABLESPACE`/`SET LOGGED` (Umschreibung nach Bauart, hergeleitet).

### Konstraints

- `LH-FA-CAP-009`: `completed` bedeutet, dass der Bestand zum Startzeitpunkt
  gelesen ist; ein `completed` mit fehlenden Zeilen ist ein stiller Verlust, den
  weder Happy Path noch Negative zulassen.
- Ein Run hinterlässt bei jedem Fehler nichts (`ADR-0111` Teilfrage 4, eine
  Schreibtransaktion); ein früher Abbruch verliert keinen Bestand.
- Die Menge der Fehlerklassen ist geschlossen (`ADR-0023`, `SPEC-008`); die
  Klasse `schema` ist im Run der Nichtanwendbarkeit einer Regel vorbehalten
  (`ADR-0117` Festlegung 5).
- Der Port `TableSnapshotPort` und die Fehler-Sentinels bleiben unverändert
  (kein Port-Diff): die Erkennung liegt im Adapter.

## Entscheidung

Wir wählen **eine Lesesperre und einen Filenode-Vergleich als die beiden ersten
Schritte der Lese-Transaktion, mit Abbruch als `transient`**. Fünf Festlegungen.

**Was diese ADR von `ADR-0111` ersetzt** — genau eine Ergänzung im Ablauf, Schritt 3:
zwischen dem Import (Schritt 2) und der Spaltenliste steht jetzt Festlegung 1 und 2.
Alles Übrige des Schritts 3 (Spaltenliste, Cursor, Blockgröße, Textform) bleibt
wörtlich in Kraft.

### Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun, das Fenster als bekannte Grenze führen | keine Änderung | ein stilles `completed` mit fehlenden Zeilen bricht `LH-FA-CAP-009` (Happy Path, Negative); die Ursache ist kein Fehler des Betreibers, der Betreiber kann sie nicht erkennen; gemessen an einem Umschreiben (Zeilen `0` bei 3 vorhandenen) |
| B — nur `LOCK TABLE … IN ACCESS SHARE MODE` als erste Anweisung | eine Anweisung | schließt nur das Fenster ab der Sperre, nicht das zwischen Export und Sperre; gemessen: mit gehaltener Sperre weiter `rows_in_snapshot=0` |
| C — nur der Filenode-Vergleich, ohne Sperre | eine Anweisung | zwischen Vergleich und `DECLARE` bleibt ein Fenster (ein Umschreiben dazwischen ist unerkannt); der Vergleich ist ohne gehaltene Sperre nicht stabil |
| D — Zeilenzahl-Plausibilität gegen `reltuples` (`ADR-0113` Warnkriterium) | kein Katalogvergleich | `reltuples` ist eine Schätzung: `-1` für nie analysierte Tabellen, im Snapshot nach dem Umschreiben unverändert (gemessen); keine Prüfung, ein Warnkriterium |
| E — Umschreiben über `relfrozenxid`/`xmin` des `pg_class`-Eintrags erkennen | keine Funktion neben der Tabelle | die SQL-Lesung im Snapshot zeigt den Eintrag **vor** dem Umschreiben (gemessen an `reltuples`, `relfilenode`); der neue Eintrag ist für den Snapshot unsichtbar (aus der gemessenen Snapshot-Sicht hergeleitet) |
| F — Sperre in einer eigenen Sitzung **vor** der Slot-Anlage halten, bis der Snapshot importiert ist (Verhinderung statt Erkennung) | schließt das Fenster vollständig, kein Abbruch | zwei Sitzungen je Run; die Sperre steht während der Slot-Anlage, die auf laufende Schreibtransaktionen wartet (`ADR-0111` M3, bis zum Zeitlimit): eine DDL, die in einer Transaktion mit Kennung auf sie wartet, und die Slot-Anlage, die auf jene Transaktion wartet, verklemmen bis zum Zeitlimit (hergeleitet, nicht gemessen); die Sperr-Warteschlange staut Leser der Quelltabelle über die Dauer der Slot-Anlage, nicht nur über die des Runs |
| G — den Slot-Export der Replication-Verbindung für den Katalogvergleich nutzen | atomar zum Snapshot | der exportierte Snapshot gilt nur bis zum nächsten Befehl der Replication-Verbindung (`ADR-0111` Teilfrage 1, Option A, Contra), sie kann keinen Vergleich mehr ausführen |
| **H — Sperre, dann Filenode-Vergleich, danach Spaltenliste und Cursor (gewählt)** | schließt die Lücke ab der Sperre und erkennt jedes Umschreiben davor (gemessen); zwei Anweisungen; kein Port-, Modell- oder Schema-Diff; die Sperre gilt ohnehin ab dem Cursor bis zum Run-Ende | Fehlalarm bei `VACUUM FULL`, `CLUSTER`, `TRUNCATE` (Festlegung 4); der Run wartet auf eine DDL, die die Sperre im Fenster hält |

### Festlegung 1 — Lesesperre zuerst

Nach dem Import des Snapshots und dem Ende der Replication-Verbindung führt die
Lese-Transaktion `LOCK TABLE <schema>.<table> IN ACCESS SHARE MODE` aus, mit
gequoteten Bezeichnern (`snapshotlogic.QuoteIdent`). Die Sperre gilt bis zum Ende
der Transaktion, also über die Spaltenliste, den Cursor und alle Blöcke. Sie
verlangt `SELECT` auf die Tabelle — dieselbe Betriebs-Vorbedingung wie das Lesen
(`ADR-0111` Teilfrage 1); ihr Fehlen endet als `permission` über die bestehende
Klassifikation (`42501`). Hält eine fremde Transaktion `ACCESS EXCLUSIVE`,
wartet die Anweisung bis zu deren Ende (gemessen: 5 s Wartezeit an einer 6-s-Sperre);
die Wartezeit trägt der Kontext des Aufrufs, sein Ablauf endet als `transient`
(bestehendes `snapshotlogic.Classify`).

### Festlegung 2 — Filenode-Vergleich

Als nächste Anweisung liest die Lese-Transaktion, in einer einzigen Abfrage, den
`relfilenode` des Eintrags der Tabelle in `pg_class` **im Snapshot** und vergleicht
ihn mit `pg_relation_filenode(oid)`:
`SELECT c.relfilenode <> coalesce(pg_relation_filenode(c.oid), 0) FROM pg_class c
WHERE c.oid = to_regclass(quote_ident($1) || '.' || quote_ident($2))`. **Ein
Unterschied** und **keine Zeile** (der Eintrag existiert im Snapshot nicht: neue
oder neu angelegte Tabelle) enden den Run; `false` setzt ihn mit der Spaltenliste
fort. Das `coalesce` hält den Ausdruck für Relationen ohne eigene Datei
(partitionierte Tabellen: `0` gegen `NULL`) frei von Fehlalarmen (gemessen).
Nach der Sperre (Festlegung 1) kann sich die Datei nicht mehr ändern: jede Form,
die sie ändert, verlangt `ACCESS EXCLUSIVE`.

### Festlegung 3 — Fehlerklasse und Ausgang

Der Abbruch trägt die Klasse **`transient`**: der Fehler wird als
`outbound.ErrSnapshotTransient` gewickelt und nennt Tabelle, die Ursache („nach dem
Snapshot-Export umgeschrieben") und die Abhilfe („neuer Antrag"). Keine neue
Klasse, kein neuer Sentinel, keine Änderung an `classifyError`. Die Verbindung wird
geschlossen (`snapshot.abort`), es entsteht keine Change (der Run schreibt erst
danach). Der Run endet `failed`; er ist run-lokal (`ADR-0111` Teilfrage 5).

Begründung der Klasse: eine Wiederholung trägt (der zweite Snapshot entsteht nach
dem Umschreiben), kein Eingriff in Konfiguration oder Regelstand ist nötig, und der
Erfassungspfad bleibt unberührt. `schema` wäre für `VACUUM FULL`, `CLUSTER` und
`TRUNCATE` (Fehlalarme, Festlegung 4) falsch und stünde gegen
[`ADR-0117`](0117-backfill-run-fehlerklasse-schema.md) Festlegung 5, die `schema`
der Nichtanwendbarkeit einer Regel vorbehält; `storage` verwiese den Betreiber auf
den Speicher.

### Festlegung 4 — Bekannte Fehlalarme, akzeptiert

`VACUUM FULL`, `CLUSTER` und `TRUNCATE` ändern die Datei, ohne dass im Fenster
Zeilen verloren gingen, die der Run braucht (`VACUUM FULL`/`CLUSTER`: 3 Zeilen im
Snapshot, gemessen; `TRUNCATE`: 0 Zeilen, aber der Zustand nach dem Snapshot ist
leer — Verhalten gegenüber dem WAL-Pfad: `LH-FA-DAT-003` Boundary). Der Run endet
in diesen Fällen `failed`/`transient`, obwohl die Zeilen vorhanden wären. Das ist
akzeptiert: das Fenster ist der Zeitraum zwischen Export und Sperre (Dauer nicht
gemessen), die Abhilfe ist ein neuer Antrag, und ein Unterscheiden von
`VACUUM FULL`/`CLUSTER` gegenüber `ALTER … TYPE` hätte keinen Katalog-Beleg, der
ohne Lesen der unsichtbaren Zeilen auskommt (Option E).

### Festlegung 5 — Was diese ADR nicht ändert

`DROP COLUMN` und `RENAME COLUMN` im Fenster enden weiter am `DECLARE` mit `42703`,
Klasse `storage` — sichtbar, kein Verlust (E2E-belegt). Eine kompatible
Spalten-Erweiterung im Fenster bleibt der Fall von [`ADR-0116`](0116-backfill-schema-version-referenz-reichweite.md).
Der Ablauf aus Teilfrage 1 (Slot, Export, Import, Cursor, Blöcke), die Position `X`
und die Lückenfreiheit gegenüber dem WAL-Pfad (Teilfrage 3) bleiben.

## Konsequenzen

- Positiv: `completed` heißt für die gemessenen Fälle wieder „Bestand zum Startzeitpunkt
  gelesen" — ein Umschreiben im Fenster endet sichtbar `failed`, statt still Zeilen zu
  verlieren.
- Positiv: zwei Anweisungen, kein Port-, Modell-, Schema-, Rollen- oder
  Klassen-Diff; `ADR-0117` und die Fehlertext-Konvention bleiben unberührt.
- Negativ: Fehlalarme bei `VACUUM FULL`, `CLUSTER`, `TRUNCATE` (Festlegung 4).
- Negativ: der Run hält seine `ACCESS SHARE`-Sperre bereits ab dem Beginn der
  Lese-Transaktion statt ab dem `DECLARE` (einige Anweisungen früher; gleiche Sperre,
  gleiches Ende); eine DDL, die `ACCESS EXCLUSIVE` will, wartet bis zum Run-Ende, und
  Leser der Quelltabelle, die nach ihr anfragen, stauen sich hinter ihr (Sperr-
  Warteschlange, hergeleitet, nicht gemessen). Die Sperre besteht bereits seit dem
  `DECLARE` (`slice-backfill-e2e` §6, hergeleitet); die Benutzerdoku nennt sie an
  sichtbarer Stelle (Folgepflicht 3).
- **Akzeptiertes Negativ** (kurz begründet, keine Folgepflicht): *ein Umschreiben nach
  der Sperre* ist unmöglich, solange die Transaktion läuft; *ein Umschreiben nach dem
  Run-Ende* ist kein Fenster dieses Runs.

### Folgepflichten

Jede Pflicht hat einen Träger im Architect-Verdikt zu dieser ADR (Verzeichnis
`docs/reviews/`, Verdikt „backfill-tabellen-rewrite-im-fenster"); **die genannten
Träger ändert diese ADR nicht**.

1. **Umsetzung** (Slice-Zuschnitt: das Verdikt): `importSnapshot` führt Festlegung 1
   und 2 zwischen Import und `readColumns` aus; die beiden Anweisungen bilden netzlos
   prüfbare Funktionen in `snapshotlogic`; die Klassifikation bleibt.
2. **Spec:** `LH-FA-CAP-009.a`, Absatz „Mechanismus", ein Satz — Wortlaut-Vorschlag:
   „Der Import steht unter einer Lesesperre der Tabelle; wurde die Tabelle zwischen
   Snapshot-Export und Sperre umgeschrieben, endet der Run `failed` (Klasse
   `transient`) ohne Change, ein neuer Antrag beginnt neu."
3. **Handbuch:** Betriebshinweis im Abschnitt „Bestand als Backfill überführen":
   die Sperre des Runs, ihre Wirkung auf DDL und der Ausgang `failed`/`transient`
   bei einem Umschreiben im Fenster.
4. **Träger nachziehen** (`AGENTS.md` §3.13): der Plan von `slice-backfill-e2e`
   (§3 Ansatz-Ergebnis, §6 Risiko „Tabellen-Rewrite im Fenster"), die Zeile
   `make test-integration` der `harness/README.md` (Phase DDL-Fenster) und die
   E2E-Abdeckungstabelle, soweit die Umsetzung sie bewegt.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Test, reale PostgreSQL (Tier `test-replication`, interner Test: `exportSnapshot`, DDL, `importSnapshot`) | ein Umschreiben (`ALTER … TYPE`) zwischen Export und Import endet mit `ErrSnapshotTransient`, ohne geöffneten Cursor; die Kontrolle ohne DDL liest alle Zeilen; eine DDL, die `ACCESS EXCLUSIVE` über die Sperranweisung hält, lässt den Import warten und danach abbrechen; eine Mutation (Vergleich entfernt) färbt den ersten Test rot (Zeilen `0`, Fehler `nil`) | `make test-replication` |
| Go-Test | die Anweisungen der Sperre und des Vergleichs sind gequotet und enthalten `coalesce(pg_relation_filenode(c.oid), 0)` | `make test` |
| E2E | Phase DDL-Fenster: `ALTER … TYPE` (Umschreiben) im Fenster nach dem Export → Run `failed`, Klasse `transient`, keine Change; ein neuer Antrag danach endet `completed` mit allen Zeilen | `make test-integration` |
| Review-Prüfpflicht | die Reihenfolge Import → Sperre → Vergleich → Spaltenliste → Cursor; kein `DECLARE` vor dem Vergleich | — (kein Gate) |

## Re-Evaluierungs-Trigger

- **`CREATE_REPLICATION_SLOT … EXPORT_SNAPSHOT` ändert sich** in einer neuen
  PostgreSQL-Hauptversion, oder PostgreSQL macht ein Umschreiben MVCC-sicher:
  Erkennung und Fehlalarme neu bewerten (`SPEC-012`: 17/18 gemessen).
- **Fehlalarme treten im Betrieb auf** (`VACUUM FULL`/`CLUSTER` im Fenster,
  beobachtet als Run `failed` mit dieser Ursache): Folge-ADR zur Unterscheidung
  oder zur Verhinderung (Option F).
- **Fortsetzen mit Checkpoint oder Parallelisierung** (`ADR-0111`-Ausbaustufe):
  Ort und Zeitpunkt von Sperre und Vergleich neu bewerten.
- Sonst permanent.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-24 | Proposed — Architect-Vorschlag zum Befund „Tabellen-Rewrite im Fenster" des Ansatz-Ergebnisses von `slice-backfill-e2e`; Reichweite an Wegwerf-PostgreSQL 17.11/18.6 gemessen | [`LH-FA-CAP-009`](../../../spec/lastenheft.md) |
| 2026-09-24 | Accepted — Annahme durch den Auftraggeber (Vollmacht zum Architect-Zug) samt Sperre + Filenode-Vergleich, Klasse `transient` | [`ADR-0111`](0111-backfill-bestand-snapshot-bulk-copy.md) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0118` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
