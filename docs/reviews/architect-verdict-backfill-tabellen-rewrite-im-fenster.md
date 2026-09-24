# Architect-Verdikt: Backfill — Tabellen-Rewrite im Fenster zwischen Snapshot-Export und Cursor

**Rolle:** Architect (Modul 8)

**Anlass:** Eine Übergabe aus dem Ansatz-Ergebnis von `slice-backfill-e2e`
(Implementer-Messung im Wegwerf-Lauf, nicht als Test committet): bei einem
`ALTER TABLE … ALTER COLUMN … TYPE varchar(64)` im Fenster zwischen Snapshot-Export
und `DECLARE` endet der Run `completed` mit `rows_copied` 0 bei drei Zeilen der
Tabelle. Der Implementer hat den Befund ausdrücklich als **Entscheidung** an den
Architect gereicht (Folge-ADR oder Schärfung), keine Fixrunde und keinen
Testanpassungs-Pfad gewählt. Der Rollenwechsel Planner → Architect → Planner folgt
Modul 8 §Konflikt-Pfad.

**Rolleninhaber:** pt9912 (Architect-Zug, anderer Kontext als der Implementer-Lauf
des Slice)

**Datum:** 2026-09-24

**Bezug:** [`LH-FA-CAP-009`](../../spec/lastenheft.md),
[`LH-FA-DAT-003`](../../spec/lastenheft.md),
[`LH-FA-CAP-009.a`](../../spec/pflichtenheft.md),
[`SPEC-008`](../../spec/pflichtenheft.md),
[`SPEC-029`](../../spec/pflichtenheft.md),
[`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md),
[`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md),
[`ADR-0116`](../plan/adr/0116-backfill-schema-version-referenz-reichweite.md),
[`ADR-0117`](../plan/adr/0117-backfill-run-fehlerklasse-schema.md),
[`ADR-0118`](../plan/adr/0118-backfill-umschreiben-im-snapshot-fenster.md); der Slice
ist als Kennung genannt, nicht als Pfad-Link (ein Slice wechselt die
Lifecycle-Ablage, ein Pfad-Link bräche mit)

**Erzeugte Artefakte dieses Zugs:**

- [`ADR-0118`](../plan/adr/0118-backfill-umschreiben-im-snapshot-fenster.md)
  (Accepted, `Supersedes ADR-0111` — nur der Ablauf-Schritt 3 der Teilfrage 1, um
  die Schritte Sperre und Vergleich ergänzt); Index-Zeile in `docs/plan/adr/README.md`
- dieses Dokument

Die Pläne, die Spec, das Handbuch und der Code bleiben in diesem Zug unberührt.

---

## Verdikt

**1 — Defekt des Vertrags, vor der Welle-Closure zu beheben (Verdikt 2/3: die
Entscheidung wird per Folge-Entscheidung geschärft).** Der Befund gilt, ist aber
enger, als der Ansatz-Ergebnis-Satz es fasst:

- `LH-FA-CAP-009` verlangt im Happy Path, dass jede zum Startzeitpunkt vorhandene
  Zeile lesbar ist, und im Negative, dass kein Bestand „still" verloren geht.
  `completed` mit `rows_copied` 0 bei vorhandenen Zeilen ist genau dieser stille
  Verlust; ein Fehler wird nicht gemeldet, der Betreiber kann ihn nicht erkennen.
  Eine „benannte Grenze" wäre eine dokumentierte Lücke gegen den Vertrag, für die
  eine Korrektur von zwei Anweisungen zur Verfügung steht (Punkt 2); sie ist hier
  keine akzeptierbare Größe.
- `ADR-0111` hat die Aussage „Bestand vollständig" über die Lückenfreiheit
  gegenüber dem **WAL-Pfad** hergeleitet (Teilfrage 3). Das Fenster zwischen
  Snapshot-Export und Lesesperre war dort nicht benannt. Die ADR wird nicht
  überschrieben; `ADR-0118` ergänzt den Ablauf.
- **Korrektur an der Übergabe:** „ebenso `VACUUM FULL`/`CLUSTER`" trifft für den
  **Verlust** nicht zu. Gemessen (PostgreSQL 18.6): beide lassen 3 Zeilen im
  Snapshot; die Datei ändert sich trotzdem (§Befundlage). Verlust tragen das
  Umschreiben durch `ALTER … TYPE` und `TRUNCATE` (bei `TRUNCATE` ist der
  Zustand nach dem Snapshot leer, siehe `ADR-0118` Festlegung 4).

**2 — Minimale, erzwingbare Abhilfe: Lesesperre, dann Filenode-Vergleich (gemessen).**
Zwei Anweisungen nach dem Import, vor der Spaltenliste:
`LOCK TABLE … IN ACCESS SHARE MODE`, dann ein Vergleich von `pg_class.relfilenode`
(im Snapshot gelesen) mit `pg_relation_filenode(oid)` (aktueller Katalog). Nur die
Kombination schließt die Lücke; jede der beiden Anweisungen allein tut es nicht
(Belege in der nächsten Sektion). Träger: `ADR-0118` (Accepted).

**3 — Zuschnitt: Fixrunde im laufenden `slice-backfill-e2e`; Klasse `transient`; eine
neue ADR, keine Ablösung von `ADR-0117`.**

- **Zuschnitt (Empfehlung):** kein eigener Slice, sondern eine Fixrunde im
  laufenden `slice-backfill-e2e`, **vor** dessen Closure und damit vor der
  Welle-Closure. Gründe: (a) der Slice hat den Befund gemessen und trägt die
  Haltepunkt-Mechanik (`docker pause` öffnet das Fenster nach dem Export), die die
  Belege der Fixrunde braucht; (b) die Änderung ist klein (eine Funktion im
  Adapter, zwei Anweisungs-Funktionen in `snapshotlogic`, je ein Test in den drei
  Tiers); (c) ein eigener Slice verlangte wegen des WIP-Limits 1 die Closure von
  `slice-backfill-e2e` **davor** und ließe die Phase DDL-Fenster zweimal anfassen.
  Der Planner passt den Plan an: „Berührte Dateien" bekommt
  `internal/adapters/driven/postgressnapshot/snapshot.go` und
  `internal/adapters/driven/postgressnapshot/snapshotlogic/logic.go`, der Slice ist
  damit nicht mehr „nur Test"; der Reviewer liest den kleinen Produktions-Diff mit.
  **Fallback:** ist der Slice bereits im Review mit eingefrorenem Umfang, entsteht
  ein eigener kleiner Slice in `welle-backfill-bestand` mit Start nach dem `done/` von
  `slice-backfill-e2e`, und die Welle schließt erst danach.
- **Fehlerklasse: `transient`** (`ADR-0118` Festlegung 3): eine Wiederholung trägt, der
  Betreiber ändert keine Konfiguration, kein neuer Sentinel, keine Änderung an
  `classifyError`. `schema` verböte `ADR-0117` Festlegung 5 (der Nichtanwendbarkeit
  einer Regel vorbehalten) und wäre für die Fehlalarme falsch.
- **Neue ADR: ja, `ADR-0118`.** `ADR-0111` (Accepted) schreibt den Ablauf; die
  Ergänzung ändert den Referenten und ist keine Zitat-Korrektur (`AGENTS.md` §3.5).
  Umfang wie bei `ADR-0115`/`ADR-0116`/`ADR-0117`: `Supersedes ADR-0111` teilweise.

**4 — Auswirkung auf Spec und Handbuch:** je ein Satz, siehe §Folge-Arbeit.
`LH-FA-CAP-009` selbst ändert sich nicht: die Korrektur stellt die Zusage her, sie
lockert sie nicht.

---

## Befundlage

**Messung** (2026-09-24, Wegwerf-Container mit `wal_level=logical`, je Läuf eine
Tabelle `t(id int primary key, name text)` mit 3 Zeilen; PostgreSQL **18.6**, Digest
von `PG_TEST_IMAGE` im `Makefile`, und **17.11**, Digest der E2E-Matrix in
`.github/workflows/e2e.yml`; ein Container je Version, danach `docker rm -fv`).
Ablauf je Variante, wie ihn der Adapter fährt: Sitzung 1 (`replication=database`)
`CREATE_REPLICATION_SLOT … TEMPORARY LOGICAL pgoutput (SNAPSHOT 'export')`, offen
gehalten; Sitzung 2 die DDL-Variante; Sitzung 3 `BEGIN ISOLATION LEVEL REPEATABLE
READ READ ONLY; SET TRANSACTION SNAPSHOT '<name>'; LOCK TABLE public.t IN ACCESS
SHARE MODE;`, danach `count(*)` und ein Vergleich von `pg_class.relfilenode` mit
`pg_relation_filenode(oid)`. Gedruckte Zeilen (Kurzform, die Tabelle mit allen Zeilen
steht in `ADR-0118` §Gemessen):

| DDL im Fenster | `rows_in_snapshot=` | `differs=` | Versionen |
|---|---|---|---|
| keine | 3 | `false` | 17.11, 18.6 |
| `ALTER COLUMN name TYPE varchar(64)` | **0** | `true` | 17.11, 18.6 |
| `VACUUM FULL` | 3 | `true` | 18.6 |
| `CLUSTER` | 3 | `true` | 18.6 |
| `TRUNCATE` | **0** | `true` | 18.6 |
| `DROP COLUMN` | 3 | `false` | 18.6 |

Weitere gedruckte Zeilen: `DROP TABLE`/`CREATE TABLE` in der Fenster-Zeit:
`rows_in_snapshot=0`, `row_found_in_pg_class=0` (der Eintrag der neuen Tabelle
existiert im Snapshot nicht — der Vergleich behandelt „keine Zeile" als Abbruch);
partitionierte Tabelle: `snap_filenode=0 cur=NULL differs_coalesce=false` (ohne
`coalesce` ergäbe `0 <> NULL` den Wert NULL statt `false`); `ALTER … TYPE` in einer offenen
Transaktion, die 6 s hält: `lock_wait_s=5`, danach `rows_in_snapshot=0` und
`differs_coalesce=true` (17.11 und 18.6) — die Sperre wartet, und der Vergleich
danach sieht das Umschreiben. Eine Rolle mit nur `SELECT` führt die Sperranweisung in
einer `READ ONLY`-Transaktion aus (`lock_ok_select_only=false`, 17.11); ohne `SELECT`:
`permission denied for table t2`.

Die Lücke ist damit **an den Anweisungen gemessen**, nicht vermutet:

| Variante | schließt die Lücke? | Beleg |
|---|---|---|
| (a) nur `LOCK` als erste Anweisung | **nein** | `ALTER … TYPE` vor der Sperre: `rows_in_snapshot=0` trotz gehaltener Sperre |
| (b) Filenode-Vergleich (Snapshot-Lesung gegen aktuellen Katalog) | erkennt jedes Umschreiben **vor** der Sperre (`differs=true`); ohne Sperre bliebe ein Fenster bis `DECLARE` | `differs=true` für `ALTER … TYPE`, `VACUUM FULL`, `CLUSTER`, `TRUNCATE` |
| (c) Zeilenzahl gegen `reltuples` | **nein** | `reltuples` vor dem Umschreiben `-1`, im Snapshot danach weiter `-1`, im aktuellen Katalog `3` (Zeile `reltuples=…`, beide Versionen) |
| (d) `relfrozenxid`/`xmin` des `pg_class`-Eintrags | **nein** | die SQL-Lesung im Snapshot zeigt den Eintrag **vor** dem Umschreiben (dieselbe Snapshot-Sicht wie `reltuples_snapshot=-1`); für `xmin` hergeleitet |
| **(a)+(b): Sperre, dann Vergleich** | **ja** (ab Sperre bis Transaktionsende lückenlos: jede Datei-Änderung verlangt `ACCESS EXCLUSIVE`) | alle Zeilen oben |

**Katalog-Lesungen und der Transaktions-Snapshot (die Frage der Übergabe):** gemessen
verhalten sich beide anders. `SELECT … FROM pg_class` als SQL-Abfrage in der
`REPEATABLE READ`-Transaktion sieht den Eintrag **zum Snapshot** (`class_snapshot_filenode`
alt); `pg_relation_filenode(oid)` und `to_regclass` lesen den **aktuellen** Katalog
(`current_filenode` neu). Gerade die Differenz beider ist das Merkmal des Umschreibens.

**Falsch-Positive** (`VACUUM FULL`, `CLUSTER`, `TRUNCATE`; `SET TABLESPACE`/`SET LOGGED`
hergeleitet): der Vergleich meldet ein Umschreiben, obwohl der Run die Zeilen hätte
lesen können (Zeilen 3 im Snapshot bei `VACUUM FULL`/`CLUSTER`). Akzeptiert
(`ADR-0118` Festlegung 4): das Fenster ist der Zeitraum zwischen Export und Sperre
(Dauer nicht gemessen), die Abhilfe ist ein neuer Antrag; ein Unterscheiden bräuchte
ein Lesen der für den Snapshot unsichtbaren Zeilen.

**Nicht Teil der Abhilfe:** `DROP COLUMN`/`RENAME COLUMN` im Fenster enden am
`DECLARE` mit `42703`, Klasse `storage` — sichtbar, E2E-belegt; ein Verlust ohne
Meldung entsteht dort nicht (`rows_in_snapshot=3`, `differs=false`).

**Speicher:** `free -m`, Zeile „Speicher": benutzt 11426 MB (Beginn) → 11767 MB
(nach dem letzten `docker rm -fv`), frei 1881 → 1361 MB (Schwankung durch den
Seitencache der beiden Image-Pulls); dangling Volumes (`docker volume ls -qf
dangling=true | wc -l`): 34 vor, 34 nach den Läufen; kein `prune`. Die beiden Images
`postgres:17-alpine`/`postgres:18-alpine` (gepinnte Digests) liegen als Pull im lokalen
Cache.

---

## Folge-Arbeit

Jede Pflicht hat einen Träger; die Pläne ändert dieses Dokument nicht — der Planner
zieht sie nach (`AGENTS.md` §3.13).

| Träger | Nachzug | Herkunft |
|---|---|---|
| `slice-backfill-e2e` (Fixrunde) | **Umsetzung** (`ADR-0118` Folgepflicht 1): `importSnapshot` führt nach dem Import und dem Ende der Replication-Verbindung die Sperre und den Vergleich aus, **vor** `readColumns`; die beiden Anweisungen als Funktionen in `snapshotlogic` (netzlos testbar); Abbruch als `ErrSnapshotTransient`-Wicklung mit Tabelle, Ursache und Abhilfe („neuer Antrag"), Verbindung über `snapshot.abort` schließen. Der Kommentar von `importSnapshot` (Abfolge) und der Paketkopf tragen die neue Reihenfolge | Punkt 2, `ADR-0118` |
| `slice-backfill-e2e` (Fixrunde), Tests | **Store-Tier** (`make test-replication`, interner Test, `exportSnapshot` → DDL → `importSnapshot`): `ALTER … TYPE` → `ErrSnapshotTransient`, kein Cursor; Kontrolle ohne DDL → alle Zeilen; offene DDL-Transaktion → der Import wartet und bricht danach ab; die Mutation (Vergleich entfernt) färbt rot (Zeilen 0, Fehler `nil`). **E2E** (`make test-integration`, Phase DDL-Fenster): zweiter Lauf im Fenster nach dem Export mit dem Umschreiben (`docker pause`-Fenster) → Run `failed`, Klasse `transient`, keine Change; ein neuer Antrag danach → `completed` mit allen Zeilen. Der bestehende Lauf mit `DROP COLUMN` bleibt (Klasse `storage`). Der Test bindet damit das **korrigierte** Verhalten, nicht die Fehlform. Der bestehende `permission`-Test (Rolle ohne `SELECT`) ändert seine Phase (`LOCK` statt `DECLARE`), nicht die Klasse — gegen den Text der Fehlermeldung prüfen | Punkt 3 |
| `slice-backfill-e2e` (Plan) | §3 Ansatz-Ergebnis („Fenster zwischen Export und Cursor") und §6 Risiko „Tabellen-Rewrite im Fenster" auf den Ausgang „behoben durch die Fixrunde, `ADR-0118`" ziehen; DoD-Punkt ergänzen; „Berührte Dateien" um die zwei Adapter-Dateien und die drei Testdateien erweitern; §3.13-Suchlauf-Feld um die Träger dieser Änderung (Zeile `make test-integration`, `make test-replication` der `harness/README.md`, `docs/user/e2e-abdeckung.md`, Handbuch-Abschnitt, `LH-FA-CAP-009.a`) | `AGENTS.md` §3.13 |
| `spec/pflichtenheft.md` (Planner-Spec-Zug) | `LH-FA-CAP-009.a`, Absatz „Mechanismus": ein Satz — Wortlaut-Vorschlag in `ADR-0118` Folgepflicht 2 | Punkt 4 |
| `docs/user/benutzerhandbuch.md` (Planner) | Abschnitt „Bestand als Backfill überführen": (1) der Run hält eine Lesesperre der Tabelle bis zu seinem Ende, eine DDL mit `ACCESS EXCLUSIVE` wartet, Leser hinter ihr stauen sich; (2) `failed`/`transient` mit Ursache „umgeschrieben" bei einem Umschreiben im Fenster: neuer Antrag; (3) die Fehlalarme (`VACUUM FULL`, `CLUSTER`, `TRUNCATE`). Die Tabelle der Run-Zustände (`failed`) bleibt | Punkt 4 |

**Reihenfolge:** Fixrunde im laufenden Slice vor dessen Closure → `slice-backfill-e2e`
nach `done/` → die offenen Slices der Welle (`slice-backfill-bench-richtgroesse`,
`slice-backfill-sdk-origin`) sind von der Fixrunde nicht abhängig; die Welle-Closure
setzt den Ausgang „behoben" des Risikos voraus.

---

## Auftraggeber-Frage

**Keine offen.** `ADR-0118` steht `Accepted` (Vollmacht zum Architect-Zug, Aussage an
der Messung belegt, `AGENTS.md` §3.12). Hält der Auftraggeber statt der Fixrunde den
Fallback (eigener Slice) für richtig, ist das eine Planner-Entscheidung ohne Wirkung
auf die ADR.
