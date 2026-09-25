# Architect-Verdikt: Retention-Lauf — Speicher an die Seite binden statt an die Zahl der Changes

**Rolle:** Architect (Modul 8)

**Anlass:** Der Messbericht `messbericht-slice-backfill-speicher-untersuchung` weist
nach, dass der Speicher des Feed-Containers nach einem Backfill an der **Zahl der
Changes in `cdc.change`** hängt und im periodischen Retention-Lauf entsteht: er liest
je Takt alle Changes der Quelle samt beider Row Images in eine Liste. Der Planner hat
daraus den Slice `slice-retention-lauf-speicher-begrenzung` angelegt und die
Entscheidung, ob der Vertrag von `ChangeStorePort` eine ADR verlangt, dem Architect
gelassen. Der Nutzer hat entschieden: der Slice wird umgesetzt, danach folgt der
Server-Release v0.2.0. Der Zug arbeitet auf frischem Kontext (Plan, Spec,
Entscheidungs-Bestand, Code, eigene Messung), ohne Rückfragen.

**Rolleninhaber:** pt9912 (Architect-Zug, anderer Kontext als der Implementer-Lauf
der Untersuchung)

**Datum:** 2026-09-25

**Bezug:** [`LH-FA-RET-002`](../../spec/lastenheft.md)…[`004`](../../spec/lastenheft.md),
[`LH-FA-CAP-009`](../../spec/lastenheft.md),
[`LH-FA-RET-004.a`](../../spec/pflichtenheft.md),
[`SPEC-022`](../../spec/pflichtenheft.md),
[`ADR-0014`](../plan/adr/0014-retention-domain-policy.md),
[`ADR-0009`](../plan/adr/0009-change-store-outbound-port.md),
[`ADR-0034`](../plan/adr/0034-ports-nach-faehigkeiten.md),
[`ADR-0042`](../plan/adr/0042-transport-typen-am-port.md),
[`ADR-0053`](../plan/adr/0053-retention-loeschausfuehrung-cdc-admin-delete-grant.md),
[`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md),
[`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md),
[`ADR-0124`](../plan/adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md); der
Slice ist als Kennung genannt, nicht als Pfad-Link (ein Slice wechselt die
Lifecycle-Ablage, ein Pfad-Link bräche mit)

**Erzeugte Artefakte dieses Zugs:**

- [`ADR-0124`](../plan/adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md)
  (Accepted, ergänzt `ADR-0111` — löst dessen dritten Re-Evaluierungs-Trigger ein);
  Index-Zeile in `docs/plan/adr/README.md`
- dieses Dokument

Die Pläne, die Spec, das Handbuch und der Code bleiben in diesem Zug unberührt.

---

## Verdikt

**1 — Was die Retention von einem Change braucht: drei Größen, keine davon im Row
Image.** [`ADR-0014`](../plan/adr/0014-retention-domain-policy.md) legt die
Entscheidung in den Domain Core; `RetentionPolicy.AllowsDeletion(age,
changePosition, consumerPositions)` liest das Alter (aus `CommittedAt` gegen die
Wanduhr) und die Commit-Position, `DeleteChanges` die Kennung. `CommittedAt` und
Position sind Eigenschaften der **Transaktion**. Row Images, Operation, Tabelle,
Sequenz, Schema-Version und `origin` sind keine Eingabe der Freigabe (Beleg: die
Schleife in `internal/application/usecase/retention/service.go` liest nur
`record.CommittedAt`, `record.Position`, `record.Change.ID`). Die Projektion mit
beiden Bildern ist damit ein Fehler für sich, unabhängig von der fehlenden
Seitenbegrenzung. Die Regel selbst ist Zeit **und** Consumer-Position (kein Größen-
oder Anzahl-Kriterium); was bleiben muss und wie es gehalten wird:

| Eigenschaft | Träger heute und nach dem Slice |
|---|---|
| Nie löschen vor Bestätigung ([`LH-FA-RET-004`](../../spec/lastenheft.md)) | `AllowsDeletion` unverändert; jeder Kandidat läuft durch sie; die veralteten Consumer-Positionen eines längeren Laufs verzögern nur (Positionen wandern vorwärts) |
| Mindestalter ([`LH-FA-RET-003`](../../spec/lastenheft.md)) | Alter je Kandidat aus `CommittedAt`, eine Uhr-Lesung je Lauf |
| Ordnung | die Retention braucht sie nicht (Menge, nicht Folge); die fachliche Ordnung der Lesewege bleibt unberührt |
| Backfill-Changes (`origin` `backfill`, Kennung `0bf-…`) | keine Sonderrolle ([`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 7): Alter zählt ab dem Snapshot-Zeitpunkt der Block-Transaktion; die Kennung ist ein undurchsichtiger Text, der Vergleich der Seitengrenze läuft in der Datenbank unter derselben Sortierung wie der Schlüssel, nicht in Go |
| Kein Sprung über eine Zeile | eine hinter dem Cursor sichtbar werdende Zeile wird in diesem Lauf nicht gesehen, im nächsten ja — die Löschung verschiebt sich, sie kommt nie früher |

**2 — Empfehlung: Lösungsform (d), „Projektion **und** Seiten“, die Entscheidung
bleibt im Domain Core.** Kandidaten (`ChangeID`, `Position`, `CommittedAt`) in Seiten
zu 10.000, Schlüssel `change_id` (Index des Primärschlüssels, keine Sortierung), je
Seite `AllowsDeletion` und ein `DeleteChanges` über die Freigaben **dieser** Seite.
Bewertung der Formen (Tabelle mit Zahlen: `ADR-0124` §Verglichene Alternativen):

- **(a) nur Projektion:** die Zeit fällt auf 374 bis 829 ms je 1.000.000, der Speicher
  bleibt linear (hergeleitet etwa 150 MiB je 1.000.000) und die einzige
  Lösch-Anweisung hält 1,5 s ihre Sperren. Behebt den Faktor, nicht die Abhängigkeit.
- **(b) Seiten ohne Projektion:** die Bilder bleiben im Transport; kein Beitrag zur
  Ursache.
- **(c) alles in SQL:** verletzt `ADR-0014` und den Satz „liegt nicht in einem
  Adapter“ von [`LH-FA-RET-004.a`](../../spec/pflichtenheft.md); die Regel stünde
  zweimal, die Sicherheit von [`LH-FA-RET-004`](../../spec/lastenheft.md) wäre nicht mehr
  domänentestbar. Abgelehnt; nur Löschen und Zählen dürfen in der Datenbank
  laufen (so ist es heute).
- **(d) Kombination:** gewählt.
- **Nicht gewählt, aber gemessen:** Kandidaten auf Transaktions-Ebene
  (`ADR-0124` Alternative E): 1.000 statt 1.000.000 Kandidaten bei einem Backfill,
  175 statt 2.003 ms in der WAL-Form. Sie ändern den Lösch-Vertrag und das Verhalten
  für Transaktionen ohne Change; das braucht der Speicher-Defekt nicht. Sie ist die
  Antwort auf den Trigger (a) der ADR, nicht auf dieses Ziel.

Wirkung im Einzelnen:

| Frage | Antwort | Beleg |
|---|---|---|
| Port | eine Methode mehr an `ChangeStorePort` (`ReadRetentionCandidates`), ein Transport-Typ `RetentionCandidate`; `ReadChanges`, `DeleteChanges` und ihre Aufrufer bleiben. Adapter und vier Test-Fakes tragen die Methode mit | `ADR-0009` benennt die Kosten; Fakes: `grep` über `DeleteChanges(` in `*_test.go` |
| Transaktionsgrenzen | je Seite eine Transaktion (Löschung samt Waisen-Bereinigung, unverändert); Lauf nicht atomar, ein Abbruch hinterlässt ein Präfix, der nächste Takt setzt fort | `ADR-0124` Festlegung 6 |
| Sperrdauer | eine Seite: 13,9 bis 242,2 ms statt einer Anweisung über die ganze Menge: 1,5 bis 1,6 s je 1.000.000 | gedruckte Zeilen unten, Q6 |
| `cdc_capture_lag` und Live-Last | **nicht gemessen**. Hergeleitet: die Kandidaten-Abfrage liest ohne Zeilensperre (MVCC), die Löschung sperrt nur die gelöschten Zeilen und die Waisen-Transaktionen; die Erfassung schreibt neue Transaktionen (Ausnahme: die Wiederholung einer Transaktion nach einem Absturz). Die Retention läuft auf einer eigenen `cdc_admin`-Verbindung. Zusage: `make bench` (`tools/bench-scaling.sh`) einmal nach der Änderung mit der bestehenden 60-s-Grenze von `cdc_capture_lag` | Code: `internal/bootstrap/wiring.go` (eigene Verbindung `retentionStore`) |
| Rechenzeit der Datenbank je Takt | linear in der Zahl der Changes: 0,38 s (Backfill-Form) bis 2,0 s (WAL-Form) je 1.000.000; nicht mehr als der Ist-Stand (1,9 bis 2,0 s für die Abfrage) | gedruckte Zeilen unten |

**3 — Messung.** Siehe Abschnitt „Messung und gedruckte Zeilen“. Kurz: die Seiten-
Abfrage kostet je Seite höchstens 80 ms und sortiert nicht; die Gesamtzeit je Lauf
entspricht der der bisherigen Abfrage; der Gewinn ist der Speicher (Seite statt
Menge) und die Sperrdauer, nicht die Rechenzeit.

**4 — ADR: ja, [`ADR-0124`](../plan/adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md).**
Gründe: (a) der Port-Vertrag wächst um eine Methode mit eigener Semantik (Seite,
Cursor, Ordnung des Schlüssels statt der fachlichen Ordnung); (b) `ADR-0111` nennt
genau diesen Fall als Re-Evaluierungs-Trigger („Folge-ADR zu Retention-Paging bzw.
Lösch-Prädikat statt Kennungsliste“) — die Antwort steht in einer Entscheidung, nicht
im Slice; (c) die Festlegungen „Lauf nicht atomar“ und „Regel bleibt im Domain Core,
kein SQL-Prädikat“ sind Entscheidungen, deren Verletzung der Reviewer gegen einen
Text prüfen kann. Form: **ergänzt** `ADR-0111` und `ADR-0014` (kein `Supersedes`):
`ADR-0014`s Entscheidung gilt unverändert, und der Satz in `ADR-0111` Teilfrage 7 zur
Lesung aller Changes beschreibt den Stand vor dieser ADR (`AGENTS.md` §3.5:
Accepted-ADRs bleiben stehen; ein `Supersedes` verlangte die Ablösung einer
Entscheidung, hier wird ein dort benannter Trigger eingelöst). Vollmacht: Nutzer.

**5 — Zuschnitt und Tests.** Umfang M bleibt. Der Slice-Plan des Planners ist die
Vorgabe; die Abweichungen stehen in „Abweichungen vom Plan“.

**6 — Andere Lesepfade und Reichweite.** Die Ursache betrifft **keinen** weiteren
produktiven Aufrufer von `ReadChanges` mit Wirkung im Feed-Container:

- `ReadChangesService` (`GET /changes`) reicht `limit` durch; ohne `limit` liest der
  Aufruf unbegrenzt, das ist eine bewusste Festlegung
  ([`SPEC-022`](../../spec/pflichtenheft.md), „Noch nicht begrenzt“). Ein Aufruf über
  1.000.000 Changes hätte denselben Bedarf je Change (aus dem Code gelesen, nicht
  gemessen); die Last liegt beim Aufrufer. Akzeptiertes Negativ, kein neuer Vorgang:
  der Satz steht im Pflichtenheft.
- Der View-Direktzugriff `cdc.changes` liest in der Datenbank, nicht im
  Feed-Container.
- Die drei Live-Wege (gRPC, SSE, NATS-Vollinhalt) lesen nicht aus dem Change Store,
  sondern aus dem In-Prozess-Broadcaster (`internal/bootstrap/wiring.go`,
  `changeBroadcaster`); `internal/adapters/driving/grpc` und
  `internal/adapters/driven/natsstream` tragen keinen Aufruf von `ReadChanges` (Beleg:
  `grep` über beide Pakete ohne Treffer).
- `POST /retention/run` ruft denselben Use Case und erbt die Seiten; die Antwortzeit
  ist die Summe über die Seiten.

**Aussage für die Release-Beschreibung (v0.1.2 und Vorgänger).** Der Defekt sitzt in
allen drei veröffentlichten Server-Versionen. Belege: `git diff v0.1.2 HEAD --
internal/application/usecase/retention internal/application/port/outbound/changestore.go`
ist leer; `git show v0.1.0:internal/application/usecase/retention/service.go` und
`v0.1.1:…` tragen dieselbe Zeile `ReadChanges(ctx, outbound.ChangeQuery{Source:
command.Source})`, und `retentionInterval = 10 * time.Second` steht in beiden
`wiring.go`; der Baum von `v0.1.2` enthält keinen Backfill. Der Feed-Container von
v0.1.x liest deshalb je Takt (10 s) alle Changes der Quelle mit beiden Row Images;
sein Speicher hängt an der Zahl der Changes in `cdc.change`. Ohne blockierenden
Consumer sind das die Changes der letzten 24 Stunden (Mindestalter 24 h, `retentionMinAge`),
mit einem Blocker mehr. Hergeleitet aus dem Messbericht (1,03 bis 1,57 KiB je Change bei
schmalen Zeilen mit **einem** Bild; Zeilen mit zwei Bildern, also UPDATE und DELETE, sind
größer und ungemessen): 10 Changes/s halten 864.000 Changes (Rate mal 86.400 s), das
sind etwa 0,85 bis 1,3 GiB (864.000 × 1,03 bis 1,57 KiB) an Spitze je Takt. Nicht am
Bild von `v0.1.2` gemessen: der Code ist gleich, die Messungen liefen am Baum vor dieser
Änderung. Vorschlag für die Beschreibung von v0.2.0: „Behoben: der Bereinigungslauf
las je Takt alle Changes der Quelle in den Speicher; er liest jetzt seitenweise
(10.000 Kandidaten, ohne Row Images). In v0.1.x wächst der Speicher des Feed-Containers
mit der Zahl der gespeicherten Changes.“

---

## Befundlage

**Warum die Untersuchung nicht bei „Paging“ endet.** Der Messbericht hat den Schalter
„Bereinigung aus“ gemessen (flach) und die Zahl je Change (1,03 bis 1,57 KiB). Er hat
nicht gefragt, **was** der Lauf braucht. Die Antwort (Verdikt 1) macht aus einem
Paging-Slice einen Slice mit zwei Schnitten: die Projektion (Faktor sieben bis zehn,
hergeleitet) und die Seite (Abhängigkeit von der Zahl). Beide zusammen sind der kleinste
Schnitt, der das Ziel erreicht.

**Die ausbleibenden Takte ab 2.000.000 Changes** (Messbericht Abschnitt 3.3, Ursache dort
nicht belegt): die Rechenzeit der Abfrage in der Datenbank erklärt sie nicht — 1,9 bis
2,0 s je 1.000.000 Changes (gemessen, Client-Lauf der bisherigen Abfrage) ergeben linear
etwa 4 s je 2.000.000 (hergeleitet), gegen ein Fenster von mehr als 60 s ohne Takt. Die
plausible Erklärung liegt im Feed-Prozess (Dekodierung von 2.000.000 Zeilen mit Bildern und
Garbage Collection bei einem Heap von 3 GiB); **Hypothese, nicht gemessen**. Der Slice
klärt sie ohne Zusatzaufwand: `BENCH_MEM_STAGES=1000000` mit drei Runs lässt den Feed bei
1.000.000, 2.000.000 und 3.000.000 Changes laufen; bleiben die Takte dort aus, ist die
Hypothese falsch und der Slice trägt den Befund in seinen Bericht.

### Messung und gedruckte Zeilen

**Umgebung.** Host: Linux 6.8.0-139-generic, 20 CPU, `free -m` „gesamt“ 31.817 MB; ein
Container `postgres:18-alpine` (Digest aus `tools/bench-lib.sh`, PostgreSQL 18.6, Standard-
Einstellungen), Netz und Container über `bench::start_postgres` mit eigenen Namen
(`pgc-arch-pg`), Ende `docker rm -fv` und `docker network rm` (Falle). Tabellen nach
`tools/schema/schema.yaml` per DDL nachgebaut (nicht per d-migrate-Rollout; die
Spalten, `PRIMARY KEY (change_id)`, `PRIMARY KEY (transaction_id)`,
`UNIQUE (transaction_id, sequence)`, die Fremdschlüssel und der Index
`cdc_transaction_source_position_idx (source_id, commit_position)` sind dieselben):

```sql
CREATE TABLE cdc.transaction (transaction_id text PRIMARY KEY, source_id text NOT NULL REFERENCES cdc.source(source_id), commit_position bigint NOT NULL CHECK (commit_position > 0), committed_at timestamptz NOT NULL DEFAULT current_timestamp);
CREATE INDEX cdc_transaction_source_position_idx ON cdc.transaction (source_id, commit_position);
CREATE TABLE cdc.change (change_id text PRIMARY KEY, transaction_id text NOT NULL REFERENCES cdc.transaction(transaction_id), source_table_id text NOT NULL REFERENCES cdc.source_table(source_table_id), sequence bigint NOT NULL CHECK (sequence >= 1), operation text NOT NULL, old_data jsonb, new_data jsonb, schema_version text NOT NULL REFERENCES cdc.schema_version(schema_version_id), origin text, CONSTRAINT uq_change_transaction_sequence UNIQUE (transaction_id, sequence));
-- Füllung: 1.000.000 Changes über generate_series; new_data = jsonb_build_object('id', g, 'name', 'name-'||g, 'qty', g%97, 'price', (g%1000)/10.0, 'flag', (g%2=0)), old_data NULL;
-- Backfill-Form: 1.000 Transaktionen zu je 1.000 Changes, commit_position 1000, origin 'backfill';
-- WAL-Form: 1.000.000 Transaktionen zu je 1 Change, commit_position g+1, origin NULL; committed_at = now() - 1 h.
```

**Speicher und Bestand.** `free -m`, Zeile „Speicher“, Spalte „verfügbar“: 13.949 MB vor
dem ersten Lauf, 13.991 MB vor dem zweiten, 14.550 MB nach beiden (Schwankung durch den
Seitencache); dangling Volumes (`docker volume ls -qf dangling=true | wc -l`): 34 vor,
34 nach den Läufen; kein `prune`; ein Container je Lauf, danach `docker rm -fv`. Ein
zweiter Satz Container (`pgc-bench-*`, ein Review-Lauf) lief parallel auf demselben Host:
die Zeiten sind Einzelläufe (`n` = 1) unter fremder Last, bei warmem Cache, Ausgabe der
Abfragen an `/dev/null` (`\o /dev/null`), Server-Zeit aus `\timing` und `EXPLAIN (ANALYZE)`.

**Gedruckte Zeilen — Backfill-Form** (Größen, dann je Abfrage; PostgreSQL 18.6):

```text
 changes | transactions | change_size | tx_size | avg_new_data_bytes
 1000000 |         1000 | 297 MB      | 184 kB  |                107
Q0 alt (SelectChanges, beide Bilder, Ordnung), Ausgabe verworfen:  Time: 1974.831 ms (00:01.975)
   Plan: Sort Method: external merge  Disk: 94952kB / Worker 0: 77784kB / Worker 1: 47736kB
         Execution Time: 1007.593 ms
Q1 Projektion ohne Bilder, ohne Ordnung, ganze Menge:               Time: 373.905 ms
Q2 Plan einer Seite (change_id > '', LIMIT 10000): Index Scan using change_pkey, Nested Loop, Execution Time: 5.254 ms
Q3 Änderungs-Ebene Seite 1000 : 1000000 Zeilen, 1000 Seiten, gesamt 630 ms, erste Seite 1.33 ms, längste Seite 1.33 ms
Q3 Änderungs-Ebene Seite 10000 : 1000000 Zeilen, 100 Seiten, gesamt 379 ms, erste Seite 3.86 ms, längste Seite 4.60 ms
Q3 Änderungs-Ebene Seite 100000 : 1000000 Zeilen, 10 Seiten, gesamt 352 ms, erste Seite 35.81 ms, längste Seite 35.84 ms
Q4 Transaktions-Ebene ganze Menge:                                  Time: 0.395 ms
Q5 Transaktions-Ebene Seite 1000 : 1000 Zeilen, 1 Seiten, gesamt 0 ms
Q6 (Transaktion, ROLLBACK) DELETE change ANY(ids) RETURNING / Waisen-DELETE, Seite 1000:   1.527 ms / 12.421 ms
                                                                     Seite 10000:  11.561 ms /  9.879 ms
   DELETE change, 100000 Kennungen:                                  119.963 ms
   DELETE change, alle 1000000 Kennungen:                            1532.862 ms (00:01.533)
```

**Gedruckte Zeilen — WAL-Form:**

```text
 changes | transactions | change_size | tx_size | avg_new_data_bytes
 1000000 |      1000000 | 282 MB      | 134 MB  |                107
Q0 alt (SelectChanges, beide Bilder, Ordnung), Ausgabe verworfen:  Time: 1903.475 ms (00:01.903)
   Plan: Sort Method: external merge  Disk: 65928kB / 65896kB / 66120kB; Execution Time: 958.495 ms
Q1 Projektion ohne Bilder, ohne Ordnung, ganze Menge:               Time: 828.586 ms
Q2 Plan einer Seite (change_id > '', LIMIT 10000): Index Scan, Nested Loop, Execution Time: 35.118 ms
Q3 Änderungs-Ebene Seite 1000 : 1000000 Zeilen, 1000 Seiten, gesamt 2202 ms, erste Seite 2.87 ms, längste Seite 4.23 ms
Q3 Änderungs-Ebene Seite 10000 : 1000000 Zeilen, 100 Seiten, gesamt 2003 ms, erste Seite 19.07 ms, längste Seite 79.78 ms
Q3 Änderungs-Ebene Seite 100000 : 1000000 Zeilen, 10 Seiten, gesamt 2029 ms, erste Seite 191.30 ms, längste Seite 295.25 ms
Q4 Transaktions-Ebene ganze Menge:                                  Time: 313.061 ms
Q5 Transaktions-Ebene Seite 1000 : 1000000 Zeilen, 1000 Seiten, gesamt 243 ms
Q5 Transaktions-Ebene Seite 10000 : 1000000 Zeilen, 100 Seiten, gesamt 175 ms
Q6 (Transaktion, ROLLBACK) DELETE change / Waisen-DELETE, Seite 1000:   1.500 ms / 22.140 ms
                                                        Seite 10000:  12.305 ms / 229.907 ms
   DELETE change, 100000 Kennungen:                                  163.290 ms
   DELETE change, alle 1000000 Kennungen:                            1598.090 ms (00:01.598)
```

**Lesart.** (1) Die Seiten-Abfrage sortiert nicht und ist von der Tiefe unabhängig
(längste Seite 4,6 ms in der Backfill-Form, 79,8 ms in der WAL-Form); ihre Summe ist der
bisherigen Abfrage gleich (WAL-Form 2.003 gegen 1.903 ms), in der Backfill-Form
kleiner (379 gegen 1.975 ms). (2) Die bisherige Abfrage sortiert 198 bis 220 MB auf die
Platte (`external merge`). (3) Die Seite zu 10.000 ist der Punkt, an dem die Zeit je Lauf
nicht mehr fällt (100.000: 2.029 ms) und die Sperrdauer je Löschung klein bleibt.
(4) Die Waisen-Bereinigung ist in der WAL-Form der teurere Teil der Löschung (229,9 gegen
12,3 ms): sie ist bereits heute Teil von `DeleteChanges` und wächst nicht mit dem Schnitt
dieses Slice. **Grenzen:** `n` = 1; warmer Cache (kalter Cache nicht gemessen); die
Löschzeilen laufen bis zum `ROLLBACK`, die Kosten des `COMMIT` sind nicht enthalten; nur
1.000.000 Changes, nicht 3.000.000 (der Slice misst dort, siehe „Tests“); die Größe der
Kandidaten im Go-Prozess ist nicht gemessen; kein Lauf mit gleichzeitiger Erfassung.

---

## Folge-Arbeit

Jede Pflicht hat einen Träger; die Pläne ändert dieses Dokument nicht — der Planner
zieht sie nach (`AGENTS.md` §3.13).

### Zuschnitt des Slice `slice-retention-lauf-speicher-begrenzung`

| Datei / Komponente | Änderung |
|---|---|
| `internal/application/port/outbound/changestore.go` | `RetentionCandidate` (`ChangeID`, `Position`, `CommittedAt`) und `ReadRetentionCandidates(ctx, source, after, limit)` am `ChangeStorePort`; der Kommentar nennt den Seitenvertrag (`ADR-0124` Festlegung 2), dass die Ordnung die des Schlüssels ist und dass eine leere Seite das Ende ist |
| `internal/application/usecase/retention/service.go` | Schleife über Seiten (`PageSize` = 10.000 als Konstante im Paket), `AllowsDeletion` je Kandidat, ein `DeleteChanges` je gelesener Seite (auch mit leerer Menge: der Adapter kehrt ohne Datenbank-Aufruf zurück, die bestehenden Tests bleiben gültig), Positionen und Uhr einmal je Lauf; Paketkopf und `Run`-Kommentar tragen den Ablauf |
| `internal/adapters/driven/postgresstorage/queries/queries.go`, `sqlexec/translate.go`, `store.go` | die Abfrage aus `ADR-0124` §Entscheidung als Konstante `SelectRetentionCandidates`; eine Übersetzungsfunktion neben `ReadChanges` (Scan in drei Werte, `mapper.ToPosition(source, commitPosition)`); Validierung (`limit` ≥ 1, Quelle nicht leer) im Adapter wie bei `ReadChanges` |
| Test-Fakes: `internal/bootstrap/heartbeat_internal_test.go`, `internal/application/usecase/readchanges/service_test.go`, `internal/application/usecase/retention/service_test.go`, `internal/application/usecase/capture/service_test.go` | die Methode implementieren (die ersten drei tragen nur die Schnittstelle, der der Retention ist der prüfende Fake) |
| `internal/bootstrap/wiring.go` | der Kommentar an `retentionInterval` („Jeder Takt liest alle Changes der Quelle …“) nennt die Seiten |
| Benutzerhandbuch, Pflichtenheft | siehe Träger-Tabelle unten |

**Tests (Pyramide, [`ADR-0030`](../plan/adr/0030-testpyramide.md)):**

- **Unit (`make test`, `retention/service_test.go`).** Fake mit Begrenzung je Aufruf und
  Aufruf-Zähler: (1) dieselbe freigegebene Menge bei Begrenzung 1, 2, 3, N und größer als
  N, verglichen mit einer **handgeschriebenen** Erwartung (nicht mit dem Ergebnis des
  Use Case); (2) jeder Lese-Aufruf trägt `limit` = `retention.PageSize`, `after` ist beim
  ersten Aufruf leer und danach die letzte Kennung der vorigen nichtleeren Seite;
  (3) jeder `DeleteChanges`-Aufruf trägt höchstens eine Seite und nur Kennungen dieser
  Seite; (4) **„ab Aufruf n“:** der Lese-Fehler beim Aufruf `n` hinterlässt genau die
  Löschungen der Seiten davor und liefert den Fehler, der Lösch-Fehler bei Seite `n`
  ebenso; (5) ein Fake, der kürzere Seiten liefert als `limit`, beendet den Lauf nicht
  vorzeitig (nur die leere Seite endet ihn); (6) die bestehenden Tests zu
  `LH-FA-RET-002`…`004` laufen mit angepasstem Fake ohne geänderte Erwartung — eine
  Erwartung, die sich ändern müsste, ist ein Befund im Bericht, keine stille Anpassung.
- **Store-Tier (`make test-store`, `postgresstorage`).** Reale PostgreSQL: (1) die Seiten
  einer Quelle decken genau ihre Changes ab (Vereinigung, keine Doppelten, keine Kennung
  einer anderen Quelle), bei Begrenzung 1, 2 und größer als die Menge; (2) `Position` und
  `CommittedAt` je Kandidat sind gleich denen des Datensatzes von `ReadChanges` für
  dieselbe Kennung, mit einer Kennung `0bf-…` (Backfill) und einer WAL-Kennung derselben
  Position; (3) ein Commit **hinter** dem Cursor zwischen zwei Seiten fehlt in diesem
  Durchlauf und steht im nächsten, ein Commit **davor** steht im laufenden; (4) ein
  Durchlauf über mehr als eine Seite in `PageSize` (etwa 25.000 Changes per
  `generate_series`, eine Anweisung) löscht über den echten Use Case dieselbe Menge, die
  eine unabhängige SQL-Zählung nach denselben Regeln nennt; (5) das Lesen gelingt unter
  einer `cdc_admin`-Login-Identität (Rechte aus `ADR-0053`); (6) **Mutationen der
  Eingabeseite**, je mit gesehenem Rot: der Cursor `>=` statt `>` (die letzte Kennung
  käme doppelt), der Quellfilter entfernt (eine Fremdquelle erscheint), `LIMIT`
  entfernt (eine Seite liefert alles), `ORDER BY` entfernt (der Cursor überspringt
  Kennungen).
- **Messung (Werkzeug, kein Gate).** `tools/bench-backfill-memory.sh` mit
  `BENCH_MEM_STAGES=1000000` (Default-Image, drei Runs: Bestand vor dem Run 0, 1.000.000,
  2.000.000 Changes), gedruckte Zeilen im Bericht des Slice. **Erwartung (Zusage):** die
  Spitze nach dem Run liegt in jedem der drei Runs im Bereich der Spitze im Run
  (8,9 bis 10,5 MiB, Messbericht Abschnitt 3.2) plus dem Bedarf einer Seite (hergeleitet
  etwa 1,5 MiB), und die Bereinigungs-Takte laufen bis zum Ende weiter (die Zeile
  „Bereinigung gelaufen“ zählt in 60 s etwa fünf bis sechs, hergeleitet aus dem Takt von
  10 s). Trifft das nicht zu, steht der Befund im Bericht und der Slice geht nicht nach
  `done/` (so schon die DoD). Ergänzend `make bench` (Skalierung) einmal, damit
  `cdc_capture_lag` gegen die bestehende 60-s-Grenze läuft (Regressions-Beleg der
  Live-Last, kein neues Gate).
- **E2E.** Kein neuer Test: `make test-integration` trägt den Retention-Rundlauf über
  reale Lösch-Takte (`tools/harness/run-integration-tests.sh`, Wartezeit über mehr als zwei
  Takte) und bleibt die Regression der Verdrahtung; die Seite ist an der Eingabe
  (Store-Tier) und am Speicher (Messung) belegt.

### Abweichungen vom Plan

| Stelle im Plan | Abweichung |
|---|---|
| DoD 1: „Seiten oder eine Projektion … ob eine ADR nötig ist“ | **beides**, und die ADR ist nötig: [`ADR-0124`](../plan/adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md) (Accepted). Die Rückführung `in-progress` → `next` wegen fehlender ADR entfällt |
| DoD 1: „ein Test, der bei Seitengröße 1 dieselbe Menge liefert wie bei unbegrenzter Größe“ | Die Seitengröße ist eine **Konstante**, keine Einstellung; die Begrenzung auf 1 kommt vom Fake (kürzere Seiten), der Test bleibt inhaltlich, seine Mechanik ändert sich. Dazu ein Store-Tier-Test über mehr als eine echte Seite |
| §3, Tabelle: „Seite, Schlüssel der Ordnung“ | Der Schlüssel ist **`change_id`**, nicht die fachliche Ordnung; die Datei-Liste wächst um `queries.go`, `sqlexec/translate.go`, vier Test-Fakes und den Kommentar in `wiring.go` |
| §6, Risiko 2: „Sortierung über die ganze Menge, 3.000.000 Zeilen ausbleibende Takte“ | Die Sortierung entfällt (gemessen: 198 bis 220 MB auf die Platte bei der bisherigen Abfrage); die Rechenzeit der Datenbank erklärt das Ausbleiben nicht (hergeleitet: etwa 4 s je 2.000.000); die Hypothese liegt im Feed-Prozess und wird von der Messung bei 3.000.000 beantwortet |
| §6, Risiko 3: „die Seitengrenze darf keine Zeile überspringen“ | Sie **darf** eine Zeile in diesem Lauf überspringen (Commit hinter dem Cursor), nie **dauerhaft**, und sie löscht nie früher: der Test bindet „nicht in diesem Lauf, im nächsten ja“, nicht „nie übersprungen“ |
| §1 „NICHT in diesem Slice“: „Andere Aufrufer von `ReadChanges` ohne Limit“ | bleibt richtig; Befund in Verdikt 6 (keiner mit Wirkung im Feed-Container) |
| §1: der Wert der Richtgröße | bleibt ein Liefer-Punkt nach der Nachmessung; **Erwartung** des Architects (keine Entscheidung, `ADR-0113` verlangt die Messung): die Kopierdauer trägt die Richtgröße unverändert, das Speicher-Argument für eine niedrigere Zahl entfällt; die Bewertung nach der Nachmessung bleibt beim Slice |
| DoD: `make bench` | ergänzt um einen Lauf der Skalierung als Regressions-Beleg der Live-Last (kein Gate) |

### Träger-Nachzug (§3.13, Suchlauf durch den Architect gefahren)

Suchlauf: `git grep -n -i -E 'alle Changes der Quelle|liest alle|Retention-Lauf|Bereinigungslauf|Bereinigungs-Takt|retentionInterval' -- docs spec harness internal tools README.md`
(ohne `docs/reviews` und `done/`). **Gefunden und zu ziehen:** die Träger der Tabelle.
**Gefunden und nicht zu ziehen:** `ADR-0111` (Kontext und Teilfrage 7, Accepted, der Trigger
ist der Anker dieser ADR), `ADR-0081` („produktiver Aufrufer ist allein der Retention-Lauf“,
Accepted, beschreibt den Stand der ADR), `ADR-0057`/`ADR-0118`, Beobachtungs-Belege,
`tools/harness/run-integration-tests.sh` (Kommentare zum Takt, nicht zur Lesung),
`harness/targets/bench-backfill.md` (Zeile „Bereinigungs-Takte“ beschreibt die Zählung,
nicht die Lesung). **Grenze:** der Suchlauf trifft Wortlaut, keine Zahlen; die Zahlen
(Bemessungsregel je Change, Richtgröße mit 7,7 GiB) liest der Reviewer am Diff.

| Träger | Nachzug | Herkunft |
|---|---|---|
| `spec/pflichtenheft.md` (Planner-Spec-Zug) | [`LH-FA-RET-004.a`](../../spec/pflichtenheft.md): Punkt 4 ergänzen — Wortlaut-Vorschlag ohne ADR- und Slice-Bezug: „4. Die Bereinigungsmenge wird seitenweise bestimmt: der Arbeitsspeicher eines Bereinigungslaufs hängt an der Seitengröße (10.000 Kandidaten), nicht an der Zahl der gespeicherten Changes; eine Seite trägt je Change nur Kennung, Commit-Position und Commit-Zeitpunkt, keine Row Images.“ Eine Zeile in der Änderungshistorie des Pflichtenhefts | `ADR-0124` Schärft |
| `docs/user/benutzerhandbuch.md`, „Aufbewahrung (Retention)“ | der Satz „Jeder Durchlauf liest dazu alle Changes der Quelle in den Speicher …“ wird ersetzt: seitenweise, ohne Row Images, der Speicher hängt nicht an der Zahl der Changes; die Datenbank-Arbeit je Durchlauf wächst mit ihr (Zahl aus der Messung dieses Zugs oder der Nachmessung, mit Ursprung) | Folgepflicht 2 |
| `docs/user/benutzerhandbuch.md`, „Bestand als Backfill überführen“, Absatz „Speicher des Feed-Containers“ | der Absatz („der Bereinigungslauf des Feeds liest in jedem Takt alle Changes der Quelle“) entfällt oder nennt den Stand nach der Änderung; der Verweis auf „Grenzwerte“ bleibt | Folgepflicht 2 |
| `docs/user/benutzerhandbuch.md`, „Grenzwerte“ | der Punkt „Backfill, Speicher des Feed-Containers“ und die Bemessungsregel je Change (samt Hinweis zu `docker --memory`) werden durch die **Messung nach der Änderung** ersetzt; die Sätze der Richtgröße („bezieht den Speicher … nicht ein … 7,7 GiB“) folgen der Bewertung nach `ADR-0113` (Trigger „Eine Messung liegt vor“). Version und Änderungshistorie tragen eine Zeile | Folgepflicht 2 und 3, Slice-DoD |
| `internal/bootstrap/wiring.go` | Kommentar an `retentionInterval` | Folgepflicht 2 |
| `docs/plan/planning/open/slice-retention-lauf-speicher-begrenzung.md` (Planner) | Bezug um `ADR-0124` ergänzen; DoD 1, §3-Tabelle, §6 (drei Risiken) und den §3.13-Suchlauf-Block auf dieses Verdikt ziehen (siehe „Abweichungen vom Plan“) | `AGENTS.md` §3.13 |

**Reihenfolge:** Planner zieht den Slice-Plan nach → Slice läuft (Implementer, Reviewer,
Verifier) → `done/` mit der Nachmessung im Bericht → Handbuch und Pflichtenheft im selben Zug →
Release v0.2.0 mit der Aussage zu v0.1.x (Verdikt 6).

---

## Auftraggeber-Frage

**Keine offen.** `ADR-0124` steht `Accepted` (Vollmacht zum Architect-Zug, Zahlen an den
Messungen dieses Dokuments und des Messberichts belegt, `AGENTS.md` §3.12; die Größe eines
Kandidaten im Go-Prozess ist als hergeleitet und als Zusage geführt). Hält der Auftraggeber
die Kandidaten auf Transaktions-Ebene schon jetzt für richtig, ist das eine Folge-ADR
(Trigger (a) der ADR) mit einer Festlegung für Transaktionen ohne Change, kein Eingriff in
diesen Zug.
