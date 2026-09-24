# Verifikations-Report: slice-backfill-snapshot-reader — 2026-09-24

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich +
Entscheidungs-Konformität + Plan-vs-Code-Diff + Gates. Review-Artefakte des
Reviewers: [`review-slice-backfill-snapshot-reader.md`](review-slice-backfill-snapshot-reader.md)
und [`review-slice-backfill-snapshot-reader-fixrunde.md`](review-slice-backfill-snapshot-reader-fixrunde.md);
Formvorbild dieses Reports:
[`verifikation-slice-backfill-change-origin.md`](verifikation-slice-backfill-change-origin.md).

**Gegenstand:** Slice-Plan `slice-backfill-snapshot-reader` (Welle
`welle-backfill-bestand`), Diff-Range `62d00abe..HEAD` (`7f75444f`), 19 Commits:
Lifecycle und Verantwortlich (`7673e774`, `1ae3fe3f`, `48aa388a`), Umsetzung
(`89b21079`, `d3b47b7e`, `9a2efe6a`, `641b2ba4`, `63351d4c`, `6cf9bafa`), Review-Report
`5e5d0b89`, Fixrunde 1 (`d308b3a5`, `eda2e41f`, `bf05067d`, `406810e4`),
Nachprüfungs-Review `56305a4d`, Fixrunde 2 (`d3420bd9`, `7f75444f`). Kontext,
nicht Prüfling: [`ADR-0115`](../plan/adr/0115-backfill-spaltenwerte-text-ergebnisformat.md)
(`04229b94`, `d7539e2c`). Dieser Lauf ändert weder Code noch Plan, Spec oder Doku; er
schreibt nur diesen Report. Das Arbeitsverzeichnis blieb während des Laufs sauber
(`git status --short` leer nach jedem Tier-Lauf, nachdem `tools/schema/plan.yaml` per
`git checkout` zurückgenommen war). Alle Mutationen liefen an Wegwerf-Kopien
(`git archive HEAD` in ein Scratchpad-Verzeichnis), nicht im Arbeitsbaum.

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — [`AGENTS.md`](../../AGENTS.md) §3.9)

Jeder Lauf schrieb in eine Log-Datei; der Exit-Code wurde im selben Aufruf gesondert
gesichert (`make … > <log> 2>&1; echo …=$?`), die Logs danach gelesen.

| Sensor | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `make gates` (vor dem Report-Commit) | **EXIT=0** | baseline-verify `v6.9.0 OK — 54 Dateien` · `db-package-lists-check: OK — … dieselben 4 Pakete` · `coverage-gate: OK — Coverage 83.10% erfüllt Schwelle 80%` · `d-check: 1001 Datei(en) geprüft, 0 Befund(e)` · `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD"` · generated-sync OK · a-check `gesamt: 0 Befund(e)` |
| `make test` (Race-Detector) | **EXIT=0** | `ok` für `postgressnapshot` und `postgressnapshot/snapshotlogic`, kein `FAIL` |
| `make a-check` | **EXIT=0** | `gesamt: 0 Befund(e)` |
| `make coverage-gate` | **EXIT=0** | Vorlauf `db-package-lists-check: OK — Dockerfile-Filter, Messlaeufe und DB_COVERAGE_PKGS nennen dieselben 4 Pakete`, `total: (statements) 83.1%`, `coverage-gate: OK — Coverage 83.10% erfüllt Schwelle 80%` |
| `make test-store` | **EXIT=0** | `DB-Adapter-Coverage: 80.07% (gedeckt 651 von 813 Statements; Profile gemergt: store,replication)`, `db-coverage: OK — … erfuellt Schwelle 70%` |
| `make test-replication` (PostgreSQL 18.6, `PG_TEST_IMAGE` = Repo-Default), 59 s | **EXIT=0** | `DB-Adapter-Coverage: 80.07% (gedeckt 651 von 813 Statements; Profile gemergt: store,replication)`; `--- PASS: TestSlotReserveExhaustedIsConfiguration (0.12s)` |
| `PG_TEST_IMAGE=postgres:17-alpine@sha256:7456ef82… make test-replication` (17er-Digest aus `e2e.yml` Z. 80; PostgreSQL 17.11) | **EXIT=0** | dieselbe Zeile `DB-Adapter-Coverage: 80.07% (gedeckt 651 von 813 …)`; `--- PASS: TestSlotReserveExhaustedIsConfiguration (0.16s)` |
| `bash tools/harness/run-replication-tests.sh tier` (Zeit mit `date`) | **EXIT=0 / EXIT=0** | PostgreSQL 18: 44,4 s; PostgreSQL 17: 44,5 s |
| `make commit-traceability RANGE=62d00abe..HEAD` | **EXIT=0** | `OK — 19 Commit(s) in "62d00abe..HEAD", Betreffs ohne Struktur-ID`; `git log --format=%s 62d00abe..HEAD \| grep -E 'SPEC-\|ARC-'` → 0 Treffer |
| `make doc-commits RANGE=62d00abe..HEAD` | **EXIT=0** | 1001 Dateien, 0 Befunde |
| `make doc-immutable RANGE=62d00abe..HEAD` | **EXIT=0** | 1001 Dateien, 0 Befunde; im Diff von `docs/plan/adr/` nur `0115-…` (neu) und die Indexzeile in `README.md` |
| `make doc-trace` (advisory, nur zur Kenntnis) | EXIT=0 | 80 Anforderung(en), 3 Waise(n): `LH-FA-CAP-009`, `LH-FA-CFG-007`, `LH-FA-CFG-008` — die Backfill-Waise bleibt bis zum Liefer-Slice erwartet |

Die Store-Zeile des Tiers `make test-store` merget mit dem zuletzt im `DB_COVERAGE_DIR`
liegenden Replication-Profil (Eigenschaft des Rollout-Pfads, nicht dieses Slice); die
Zahl 651/813 beider Tier-Läufe stimmt mit dem gemergten Profil überein (§3.4).

## 2. DoD — Verdikt je Zeile (§2 des Plans)

| # | DoD-Zeile | Verdikt | Realer Beleg |
|---|---|---|---|
| 1 | Snapshot-Träger M1–M4 real belegt (PostgreSQL 17 und 18) | **erfüllt** | `TestSnapshotPairedWithPoint` (M1: `beforePoint <= X < afterPoint`, Zeile 3 zwischen Slot-Anlage und Import committet, nicht im Snapshot), `TestPermissionClassWithoutSelect` (M2: `LOGIN REPLICATION` ohne `SELECT` → `permission`, danach mit Grant lesbar), `TestSlotCreationTimeout` (M3: 1-s-Limit gegen offene Schreibtransaktion → `transient`, kein Slot, danach lesbar), `TestTemporarySlotEndsWithConnection` (M4: Slot im Katalog zwischen Anlage und Import = 1, danach weg; Cursor liest den Snapshot-Stand trotz `UPDATE`/`INSERT` danach). Eigener `-v`-Lauf gegen 17.11 und 18.6: alle `--- PASS`, `ok`, `GOTEST_EXIT=0` |
| 2 | Blöcke und Spalten; `reltuples` je Version | **erfüllt** | `TestBlocks` (7 Zeilen/`B`=3: `[3 3 1 0]`; exakt an der Grenze `[3 3 0]`; `[2 0]`; leere Tabelle `[0]`), `TestColumnsSkipDroppedAndGenerated`, `TestEstimatedRows`: gedruckt `reltuples der nie analysierten Tabelle = -1 auf PostgreSQL 17.11` und `… = -1 auf PostgreSQL 18.6`, nach `ANALYZE` 40 (Zeilenzahl); `TestEstimatedRowsAnalyzedEmptyTableIsKnownZero` (`reltuples = 0` → bekannt 0) grün auf beiden |
| 3 | Bild-Parität über den Typ-Satz, drei GUC-Lagen, 17 und 18 | **erfüllt** | §3.1 und §3.2: 85 Typ-Spalten, 19 Array-Spalten, drei Lagen, WAL-Bild gegen Backfill-Bild, beide Versionen grün; Mutationen rot (§5) |
| 4 | Gate-Zuordnung (a): Paket in den namentlichen Stellen | **erfüllt** | `Dockerfile:101` (Filter mit `postgressnapshot`), `db-coverage.sh` (`DB_COVERAGE_PKGS`), `run-replication-tests.sh` (Messphase), beide Sensor-Dokumente, `harness/README.md`; Suchlauf beider Stände nachgemessen (§3.6); `make coverage-gate` EXIT=0 mit dem Listen-Vorlauf |
| 5 | Gate-Zuordnung (b): jeder Test überspringt ohne Datenbank | **erfüllt** | `go test -count=1 -v ./internal/adapters/driven/postgressnapshot` im Toolchain-Image mit `--network none`, ohne `CDC_REPLICATION_TEST_DSN`: **EXIT=0**, gedruckt **21** `--- SKIP`, **0** `--- PASS`, 0 `--- FAIL`; im Quelltext 21 `func Test…` (gezählt) |
| 6 | Gate-Zuordnung (c): netzlose Logik im Gegenstand (`snapshotlogic`) | **erfüllt** | Nenner der Stufe 2040 → 2082, das Unterpaket trägt 42 von 42 (§3.3); Filter END-verankert |
| 7 | Gate-Zuordnung (d): DB-Adapter-Coverage nachgemessen | **erfüllt** | 80.07 % (651 von 813) gegen `DB_COVERAGE_THRESHOLD` 70, PostgreSQL 17 und 18 (§1) |
| 8 | `make gates` grün, Exit gesondert | **erfüllt** | eigener Lauf EXIT=0 (§1) |
| 9 | Review durchgeführt, Report unter `docs/reviews/` | **erfüllt, mit V-3 (INFO)** | zwei Reports liegen vor (2 HIGH/4 MEDIUM/6 LOW/6 INFO und 1 HIGH/1 MEDIUM/3 LOW/5 INFO); die zweite Fixrunde nach dem Nachprüfungs-Review hat keinen weiteren Report — ich habe die Behebung von F-1…F-5 und F-8…F-10 des Nachprüfungs-Reviews selbst am Code und an den Läufen nachgemessen (§3.7); keine der Fixrunden-2-Änderungen ist unbelegt |
| 10 | §3.13-Suchlauf: Feld trägt Gefundenes und Nichtgefundenes, beide Stände | **erfüllt, mit V-1/V-2 (INFO)** | die Felder der Erstlieferung, der Fixrunde und der zweiten Fixrunde sind gegen beide Stände nachgemessen (§3.6): die Zahlen stimmen bis auf stand-relative Zellen |
| 11 | Doku-Update (`harness/README.md`, Sensor-Dateien) | **erfüllt** | `git diff 62d00abe..HEAD -- harness/README.md harness/sensors`: Paketlisten, Nenner 813 mit Lauf, Rücknahme-Messung, Grenzen Nr. 7/8 |
| 12 | Closure-Notiz mit Lerneintrag | **korrekt offen** | Plan §7 trägt „*(zu tragen bei Closure)*" |
| 13 | Reconciliation-Register — entfällt | **korrekt offen / entfällt** (V-5, INFO) | Greenfield; die Zeile trägt „entfällt", steht aber auf `[ ]` |
| 14 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | Closure-Pflicht; der Plan führt den Flake-Befund ausdrücklich als „Sache der Closure" |
| 15 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | zehn Risiko-Zeilen tragen „**Ausgang:** *(bei Closure)*" (gezählt); der Runner-Beleg von `e2e.yml` ist eine davon (§6) |
| 16 | Drei Paarungen | **korrekt offen** | hängen an der Closure von `welle-backfill-bestand` |

Kein `[x]` ohne Beleg, kein `[ ]`, das über die Rollen-Sequenz hinaus belegt wäre.

## 3. Kernaussagen, selbst gemessen

### 3.1 [`ADR-0115`](../plan/adr/0115-backfill-spaltenwerte-text-ergebnisformat.md) im Lesepfad

- **Lese-Anweisung ohne Cast und Funktion:** `snapshotlogic.CursorStatement`
  (`logic.go:90-97`) quotet jeden Namen über `QuoteIdent` und fügt sie mit `, ` zu
  `DECLARE cdc_bf_cursor NO SCROLL CURSOR FOR SELECT "a", "b" FROM "s"."t"` zusammen; der
  netzlose Test `TestCursorStatementSelectsOnlyQuotedNames` bindet die Form und verbietet
  `::`, `(`, `)`, `CAST`, `CASE`; der zweite Test bindet das Quoting an Namen mit
  Anführungszeichen, Leerzeichen und Großbuchstaben.
- **Text-Ergebnisformat, Rohbytes, NULL:** `snapshot.go` `NextBlock` liest über
  `conn.Exec(…).ReadAll()` (einfacher Query-Pfad, Text ist dort Protokoll-Eigenschaft) und
  reicht `snapshotlogic.Values` weiter (`nil` → `nil`, sonst `string(field)` als `*string`,
  leeres Feld = leerer Text; `TestValues`). Der Adapter setzt keine Sitzungs-GUC
  (`connConfig` ändert nur `replication`; das einzige `SET` ist `SET TRANSACTION SNAPSHOT`
  mit über `ValidSnapshotName` geprüftem Namen).
- **`model.BuildRowImage` unverändert:** `git diff --stat 62d00abe..HEAD -- internal/domain
  internal/bootstrap cmd` ist leer.
- **Paritätstest** (`TestImageParityWalAndBackfill`, eigene `-v`-Läufe, ungekürzt gegen
  Wegwerf-PostgreSQL mit der Tier-Konfiguration `wal_level=logical`,
  `max_wal_senders=10`, `max_replication_slots=10`, `wal_sender_timeout=2000`):

  | Version | Ergebniszeilen (gedruckt) |
  |---|---|
  | PostgreSQL 17.11 | `PostgreSQL 17.11, 85 Typ-Spalten`; drei `Bild Zeile 1 in der Lage Standard/Tokyo/Verbose`; `--- PASS: TestImageParityWalAndBackfill (0.71s)`; `ok … 4.748s`, `GOTEST_EXIT=0` |
  | PostgreSQL 18.6 | `PostgreSQL 18.6, 85 Typ-Spalten`; dieselben drei Lagen; `--- PASS: TestImageParityWalAndBackfill (0.79s)`; `ok … 5.000s`, `GOTEST_EXIT=0` |

  Anzahl der Vergleiche, **abgeleitet** aus dem Test (`checkParity`: je Lage drei Zeilen ×
  `id` plus 85 Typ-Spalten, dazu je Zeile ein Bild-Vergleich): 3 Lagen × 3 Zeilen × 86
  Spalten = 774 Roh-Text-Vergleiche und 9 Bild-Vergleiche je Version. Die Rollen-GUC
  prägen das WAL-Bild real (gedruckt: Lage Tokyo `24.09.2026 17:11:12.5 JST`,
  `+1-2 +3 +4:05:06.7`, `bytea` als `\336\255\276\357`; Lage Verbose `EDT`,
  `@ 1 year 2 mons 3 days 4 hours 5 mins 6.7 secs`); der Test bindet das über `inImage`.
  Die NULL-Zeile trägt `{"id":"3"}`, gelöschte (`dropme`) und generierte Spalten (`gen`,
  unter 18 auch `genv`) fehlen.
- **Typ-Satz gegen [`ADR-0115`](../plan/adr/0115-backfill-spaltenwerte-text-ergebnisformat.md) §Kontext und §Folgepflichten, Zeile für Zeile:** jede Klasse
  der ADR-Typliste hat eine Spalte in `parity_types_test.go` — Boolean; `int2/4/8`, `real`,
  `float8`, `numeric`; `money`; `char(n)`, `varchar(n)`, `text`, `"char"`, `name`; `bytea`;
  `date`, `time`, `timetz`, `timestamp`, `timestamptz`, `interval`; `bit`, `varbit`;
  `uuid`; `inet`, `cidr`, `macaddr`, `macaddr8`; `json`, `jsonb`, `jsonpath`, `xml`;
  `oid`, `regclass`, `regtype`, `regproc`, `xid`, `tid`, `pg_lsn` (zusätzlich
  `pg_snapshot`); `tsvector`, `tsquery`; die sieben Geometrie-Typen; vier Range-/
  Multirange-Typen und ein eigener Range-Typ; Enum; vier Domains (`int`, `text`, `int[]`,
  `bool`); zwei Composites (`c2` mit `bool`/`inet`/`char(4)`); `hstore`, `citext`, `ltree`
  in eigenem Schema; alle 18 Array-Formen der ADR (auch zweidimensional und mit
  Untergrenze, `hstore[]`) plus das Domain-über-`bool`-Array = 19 Array-Spalten
  (gezählt). Die Zahl der `add(`-Zeilen ist 85 (gezählt), gleich der gedruckten Spaltenzahl.
  Die „mindestens"-Liste der Folgepflicht 2 ist ebenfalls vollständig gedeckt. **Es fehlt
  nichts** — bestätigt. Die Lage `lc_monetary='de_DE.utf8'` ist auf dem Alpine-Image nicht
  herstellbar: gemessen auf `postgres:18-alpine` (18.6) und `postgres:17-alpine` (17.11)
  liefert `SET lc_monetary='de_DE.utf8'; SELECT 1234567.89::money::text` beide Male
  `$1,234,567.89` — die Ausnahme des Implementers ist bestätigt; die Fitness Function
  der ADR nennt drei Lagen ohne sie.
- **Guard der Pflicht-Typen:** `requiredTypes`/`missingRequired` melden einen gestrichenen
  Pflicht-Typ vor dem ersten SQL (`t.Fatalf("der Typ-Satz trägt die Typen … nicht")`).

### 3.2 Snapshot-Träger, Klassen

- **M1–M4:** siehe DoD 1. **Slot verschwindet:** `TestTemporarySlotEndsWithConnection`
  und `TestOpenSnapshotLeavesNoSlot` (kein Slot nach `OpenSnapshot`); der Slot im Katalog
  ist zwischen Anlage und Import da (1 Zeile), danach weg.
- **Klassen:** `permission` (`TestPermissionClassWithoutSelect`,
  `TestWrongPasswordIsPermission`, `TestCatalogQueryFailuresKeepTheirClass`),
  `configuration` (`TestConfigurationClass`: fehlende Tabelle, ungültige Run-Kennungen,
  ungültige Konstruktion; `TestSlotReserveExhaustedIsConfiguration`: `SQLSTATE 53400`),
  `storage` (`TestImportOfUnknownSnapshotIsStorage`, `TestNextBlockAfterClose`),
  `transient` (`TestSlotCreationTimeout`, `TestUnreachableSourceIsTransient`,
  `TestReadWithCancelledContextIsTransient`, `TestReadAfterBackendTermination`),
  `replication` (`TestSlotNameCollisionIsReplication`) — alle `--- PASS` auf 17.11 und
  18.6 im eigenen Lauf.
- **`57P01/02/03` und beendete Verbindung → `transient`:** `snapshotlogic.Classify` ordnet
  die drei SQLSTATE zu, `NextBlock` setzt bei `conn.IsClosed()` den Fallback
  `ErrSnapshotTransient`; die Port-Doku führt „die Verbindung ist abgebrochen" unter
  `transient` (`tablesnapshot.go:23-26`); `grep -rln ErrSnapshot` über alle `.go`-Dateien findet die Sentinels nur im Port, im
  Paket und seinem Unterpaket samt Tests (kein weiterer Verbraucher).
  Gedeckt durch Port-Doku und Plan §3 („Fehlerklasse einer beendeten Sitzung").
- **Keine Warn-Schwelle im Code:** `grep -niE 'warn|schwelle|threshold' snapshot.go logic.go
  tablesnapshot.go` findet nur den Kommentar zu `Estimate`; `EstimatedRows` liefert
  `(rows, known, err)` und nichts sonst ([`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md)
  Festlegung 3).

### 3.3 Gate-Nenner und Unterpaket `snapshotlogic`

Zählverfahren wie in der Sensor-Doku beschrieben, von mir nachgestellt: die Stufe im
gepinnten Toolchain-Image (`go test -coverpkg=<Liste der Stufe> -covermode=atomic`, ohne
Netz), dann Deduplizierung über die Block-Position (`awk`).

| Stand | Filter | Nenner | gedeckt (dedupliziert) | gedruckt |
|---|---|---|---|---|
| `62d00abe` (Parent des Slice) | 3 Pakete | **2040** | 1692 | `82.9%` |
| `d7539e2c` (Vorstand der Fixrunde) | 4 Pakete | **2040** | 1689 | `82.8%` |
| `7f75444f` (Diff) | 4 Pakete | **2082** | 1733 | `83.2%` (Nachstellung); `make coverage-gate` druckte in diesem Lauf `83.10%` |

Nenner 2040 → 2082 (+42) bestätigt; `snapshotlogic` trägt **42 von 42**; kein Statement
aus `postgressnapshot` selbst steht im Profil des Diffs (`grep` über das Profil: 0 Blöcke
außerhalb des Unterpakets). Der Filter des `Dockerfile` (Zeile 101) ist END-verankert
(`(^|/)(postgresstorage|postgresack|postgressnapshot|replication/receive)$`) und nimmt
das Unterpaket nicht aus. Die gedeckte Zahl ist lauf-gebunden: meine Läufe liegen zwischen
1689 und 1692 am Vorstand und bei 1733 am Diff, die Plan-/Sensor-Zahlen (1691, 1731) liegen
in dieser Spanne — das Dokument führt sie ausdrücklich als lauf-gebunden (V-1, INFO).
`postgressnapshot` (DB-Paket) ist ausgenommen; die Rücknahme wird gemessen (§5, M7).

### 3.4 DB-Adapter-Coverage

`make test-replication` auf PostgreSQL 18 und 17 druckt je `80.07% (gedeckt 651 von 813
Statements; Profile gemergt: store,replication)`. Aus dem Replication-Profil nachgerechnet
(`awk`, Block-Position dedupliziert): 245 Positionen × 3 = 735 Zeilen (Vielfachheit 3 für
alle 245); `postgresack` 32/32, `postgressnapshot` **118/122**, `replication/receive`
**153/187**; Replication-Profil gesamt 303/341 = 88,86 %, mit „nur erstes Vorkommen"
32/341 = 9,38 %; `postgressnapshot` und `receive` mit „nur erstes Vorkommen" je 0 — alle
Zahlen des Sensor-Dokuments (`db-adapter-coverage.md` §Zählbasis) stimmen. Store-Anteil
651 − 303 = 348 (abgeleitet). Die 4 ungedeckten Statements von `postgressnapshot` sind die
beiden Fehlformat-Zweige in `exportSnapshot` (`consistent_point`/Snapshot-Name nicht
lesbar), im Plan benannt.

### 3.5 Listen-Prüfskript und Skip-Eigenschaft

`tools/harness/db-package-lists-check.sh` (Aufruf in `harness/mk/coverage.mk` als erste
Rezeptzeile vor dem `docker build`): vergleicht Filter, `DB_COVERAGE_PKGS` und die beiden
Messlauf-Listen; es lockert nichts, es kommt eine Prüfung hinzu (§3.6 der Repo-Regeln;
`THRESHOLD`/`DB_COVERAGE_THRESHOLD` unverändert, `git diff` zeigt keine Schwelle).
Mutationen in §5 (M3–M6): jede endet mit Exit 1 und nennt das abweichende Paket.

### 3.6 Zahlen der Suchlauf-Felder (§3.12), beide Stände selbst gemessen

| Plan-Zahl | Meine Messung | Ergebnis |
|---|---|---|
| Pakete-Suchlauf, Parent `d7539e2c`: `Dockerfile` 1 · `antragsweg-usecase` 1 · `kern-rename` 1 · `welle-backfill-bestand` 1 · `harness/README.md` 1 · `coverage-gate.md` 16 · `db-adapter-coverage.md` 11 · `db-coverage.sh` 7 · `run-replication-tests.sh` 7 | identisch (`git grep -c -E 'postgresack\|replication/receive\|postgressnapshot' d7539e2c -- . <Ausschlüsse>`) | bestätigt |
| Diff-Stand `bf05067d`: `kern-rename` 2, `coverage-gate.md` 20, `db-adapter-coverage.md` 14, `db-coverage.sh` 4, `run-replication-tests.sh` 8 | identisch an `bf05067d`; an `HEAD` `db-adapter-coverage.md` **15** (die zweite Fixrunde ergänzt eine Zeile) | bestätigt für den benannten Stand, an `HEAD` verschoben (V-2, INFO) |
| Kein Treffer in `.github/workflows/*.yml`, `Makefile`, `harness/mk/*.mk` für die vier Pakete | `git grep` über diese Pfade an beiden Ständen: leer | bestätigt |
| `::tex[t]`-Suchlauf: 9 Zeilen in 4 Dateien am Vorgänger `5e5d0b89`; am Diff-Stand `bf05067d` 21 Zeilen in 4 Dateien (`ADR-0111` 2, `ADR-0115` 13, `done/…row-image-gemeinsam` 1, Review-Report 5; dieser Plan 0) | `5e5d0b89`: 9/4 · `bf05067d`: 21/4, gleiche Aufteilung; `HEAD`: 23 Zeilen (Plan 1, Nachprüfungs-Report 1) | bestätigt; an `HEAD` verschoben (V-2, INFO) |
| Zweite Fixrunde, Phasen-/Beschreibungstexte: Parent `56305a4d` 17 Dateien, Diff `d3420bd9` dieselben 17 mit `e2e.yml` 6 → 7, `db-adapter-coverage.md` 15 → 19, `run-replication-tests.sh` 10 → 11 | identisch (`git grep -c -i -w -E 'tier\|Messphase\|Slot-Reserve\|max_replication_slots' <Stand> -- …`) | bestätigt |
| Typ-Satz-Suchlauf: „zwei Treffer zu `nicht mitbringt`" in `BEO-PGC/roter-test-ohne-leser` bleiben | 2 Treffer dort, kein „81 Typ", „nicht anlegbar" oder „bringt sie nicht mit" mehr im Baum (außer dem Suchbefehl selbst) | bestätigt |
| Nenner 813 (472 · 32 · 122 · 187), gedeckt 651 (348 · 32 · 118 · 153) | §3.4 | bestätigt |
| Rücknahme: Dockerfile-Filter ohne `postgressnapshot` → `Coverage 78.60% unter Schwelle 80%` | §5 M7: gedruckt `coverage-gate: FAIL — Coverage 78.60% unter Schwelle 80%`, `docker build` Exit 1 | bestätigt (dieselbe Zahl) |
| 21 von 21 `SKIP` | §2 Nr. 5 | bestätigt |
| Slot-Reserve-Reihenfolge: ein verschobener Testname und `max_replication_slots=10` färben die Phase rot | §5 M8/M9 | bestätigt |
| Laufzeit der Phase `tier`: 46,3 s (18), 44,4 s (17); Reserve-Lauf 7,8 s/7,5 s | ganze Phase 44,4 s (18) und 44,5 s (17), gemessen mit `date` um den Skript-Aufruf; den Reserve-Anteil habe ich nicht getrennt zeitgemessen (`--- PASS … (0,12 s)`/`(0,16 s)` gedruckt) | Größenordnung bestätigt (rund 45 s); Zeitangaben sind Lauf-Belege |
| Aufräumen im Testcontainer: nach zwei Läufen 0 Rollen, 0 Datenbanken, 0 Slots, 0 Extensions, 0 Schemata | zwei `go test -count=1`-Läufe gegen Wegwerf-PostgreSQL 18, danach `psql`: Rollen `snap_role%` **0**, Datenbanken `snap_%` **0**, Slots **0**, Extensions `hstore`/`citext`/`ltree` **0**, Schemata `snap_ext%` **0**, Tabellen `snap_%` **0**, benutzerdefinierte Typen `snap_%` **0**, Publications **0** | bestätigt |

### 3.7 Nachprüfungs-Review (Fixrunde 1) — F-1…F-10 am Stand `HEAD` nachgemessen

- **F-1 (HIGH), Typ-Satz ohne `hstore`/`citext`/`ltree`:** behoben — die Tabelle trägt `hs`,
  `ci2`, `lt`, `ah` (85 Spalten, gedruckt); `CREATE EXTENSION … SCHEMA snap_ext_<sfx>` gelingt
  auf 17.11 und 18.6, `set.drop` räumt (Aufräum-Messung §3.6); Mutation M2 färbt genau `hs`/`ah`.
- **F-2 (MEDIUM), Flake 55006:** behoben im Test, Produktionscode unverändert:
  `git diff --stat 62d00abe..HEAD -- internal/adapters/driving/replication/receive` nennt
  **nur** `stream_test.go` (24 Zeilen hinzugefügt: `awaitSlotInactive` und ein Aufruf).
  Zählung §5.2.
- **F-3 (LOW), Phasen-Beschreibung:** die Kopfzeile, der Kommentar an der Phase `tier` in
  `run-replication-tests.sh` und Kommentar samt `name:` in `e2e.yml` nennen den
  Slot-Reserve-Lauf (Diff gelesen).
- **F-4 (LOW), Test-Hygiene:** 0 zurückbleibende Rollen (§3.6); die Rolle wird vor der
  eigenen Datenbank angelegt, die Replication-Verbindung eines Exports per `t.Cleanup`
  geschlossen.
- **F-5 (LOW), Doku-Prosa:** `coverage-gate.md` §Zählbasis nennt die Zusammensetzung des
  Nenners (2040 + 42) mit Lauf; `db-adapter-coverage.md` §Grenze Nr. 8 nennt den Test der
  Slot-Reserve mit seiner eigenen Vorbedingungs-Prüfung (der Test beginnt mit
  `os.Getenv(exclusiveDSNEnv)`, `snapshot_test.go:834-837`).
- **F-8 (INFO):** Port-Doku und `DefaultBlockSize`-Kommentar benennen `B` als Zeilenzahl;
  „nicht gemessen" steht dabei. **F-9 (INFO):** siehe §6. **F-10 (INFO):** die beiden
  Zeilen in den offenen Plänen sind minimal (je ein Aufzählungspunkt um
  `postgressnapshot` ergänzt) und im Plan benannt.

## 4. Entscheidungs-Konformität

| Entscheidung | Festlegung | Beleg | Verdikt |
|---|---|---|---|
| [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 1 | temporärer Slot `cdc_bf_<run>` mit `EXPORT_SNAPSHOT`, Import in `REPEATABLE READ READ ONLY`, Replication-Verbindung endet, Spaltenliste (`attnum > 0`, nicht gelöscht, nicht generiert), `NO SCROLL`-Cursor in Blöcken `B`, Zeitlimit → `transient`, alle Verbindungen über denselben DSN | `snapshot.go` `exportSnapshot`/`importSnapshot`/`readColumns`, `DefaultBlockSize` 1000, `DefaultSlotTimeout` 30 s (beide Startwerte, im Plan als Setzung ohne Messung geführt); M1–M4 real | konform |
| [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 3 | Position `X` = `consistent_point`, jeder Commit ≤ X im Bestand, jeder spätere nicht | `TestSnapshotPairedWithPoint`; der Adapter liefert `X` als Rohwert (`Offset()`), die Abbildung auf `SourcePosition` bleibt beim Run ([`ADR-0005`](../plan/adr/0005-sourceposition-abstrahiert-lsn.md)) | konform (Reader-Anteil) |
| [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 4 | Atomarität liegt beim Run (eine Store-Transaktion) | der Reader hält keinen Run-Zustand; ein abgebrochener Lesezustand endet als `transient`/`storage`, keine Teilablage | konform (nicht Gegenstand des Slice) |
| [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 3 | `reltuples = −1` = unbekannt; keine Ablehnung, keine Richtgröße im Code | §3.2, §2 Nr. 2; kein Schwellenwert im Quelltext | konform |
| [`ADR-0115`](../plan/adr/0115-backfill-spaltenwerte-text-ergebnisformat.md) Festlegung 1 | nur gequotete Namen, Rohbytes, `nil` ↔ NULL | §3.1 | konform |
| [`ADR-0115`](../plan/adr/0115-backfill-spaltenwerte-text-ergebnisformat.md) Festlegung 2 | Text-Format Teil der Zusage, keine Sitzungs-GUC durch den Adapter | einfacher Query-Pfad; kein `SET`/`options`; der Paritätstest fängt einen Formatwechsel rot | konform |
| [`ADR-0115`](../plan/adr/0115-backfill-spaltenwerte-text-ergebnisformat.md) Festlegung 3 | `BuildRowImage` unverändert | `internal/domain` im Diff leer; der Test ruft die Domänenfunktion für beide Bilder | konform |
| [`ADR-0115`](../plan/adr/0115-backfill-spaltenwerte-text-ergebnisformat.md) Festlegung 4 + Fitness Function | Typ-Satz als Daten, Roh-Text und Bild, drei Lagen, 17 und 18, netzloser Strukturtest | §3.1 | konform |
| [`ADR-0115`](../plan/adr/0115-backfill-spaltenwerte-text-ergebnisformat.md) Festlegung 5 (Folgepflichten 1–5) | Adapter/Port-Doku, Fixture, Strukturtest, Plan-Träger, Spec unberührt | Diff: `git diff --stat 62d00abe..HEAD -- spec` leer; der Plan-Träger (Ziel-Absatz, DoD, Messbefunde, Risiken) ist nachgezogen | konform |
| [`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md) Punkt 1 und Trigger (a) | ausgenommenes Paket = Eigenschaft „Testlauf setzt einen externen Dienst voraus"; die namentliche Liste wird nachgezogen | alle vier Stellen und beide Sensor-Dokumente; Skip-Beleg (21/21); Rücknahme rot (§5 M7) | konform |
| [`ADR-0034`](../plan/adr/0034-ports-nach-faehigkeiten.md) | Ports nach Fähigkeit | `TableSnapshotPort`: Snapshot lesen plus Katalog-Schätzung derselben Quelle, eine Konsistenzgrenze, kein Entity-Repository | konform |
| [`ADR-0080`](../plan/adr/0080-nahtform-pgconn-adapter-treiberhuelle.md) | Hülle für `postgresack` und `receive` (Festlegung 1/5); Festlegung 4 lässt ausgenommene Pakete zu | die Nahtform ist für die zwei benannten Pakete entschieden, nicht für neue Pakete; der Plan wählt den anderen Weg (Logik in `snapshotlogic`, Schritte mit Verbindung als reale Tier-Läufe); §Kontext (5) der ADR führt Unterpakete ausgenommener Pakete ausdrücklich im Unit-Gegenstand | **kein Verstoß** (Review F-13 bestätigt) |

## 5. Mutationen (Eingabeseite, selbst gesehen)

Je Mutation: Edit per `sed` an einer Wegwerf-Kopie des `HEAD`-Stands (Scratchpad), Sensor
gefahren, Exit-Code und `FAIL`-Zeilen gelesen. Der Arbeitsbaum des Repos blieb unberührt
(`git status --short` leer).

| # | Zusage | Mutation | Sensor | Ergebnis |
|---|---|---|---|---|
| M1 | Lese-Anweisung ohne Cast | `quoted[i] = QuoteIdent(column) + "::text"` in `CursorStatement` | Tier-Test `TestImageParityWalAndBackfill` gegen 18.6 und 17.11; netzlos `snapshotlogic` | **rot**, `GOTEST_EXIT=1` auf beiden Versionen, in allen drei Lagen; die Meldungen nennen genau die Spalten `bo`, `ch1`, `ch5`, `db`, `ia` (je 6×) und `xm` (3×); netzlos rot: `TestCursorStatementSelectsOnlyQuotedNames`, `TestCursorStatementQuotesEveryIdentifier` |
| M2 | Wertseite trägt Extension-Text unverändert | in `Values` `strings.ReplaceAll(string(field), "=>", "->")` | Tier-Test gegen 18.6 und 17.11 | **rot** auf beiden Versionen, genau `hs` und `ah` (je 3×) |
| M3 | Listen-Gleichheit: Paket aus `DB_COVERAGE_PKGS` | `,./internal/adapters/driven/postgressnapshot` gestrichen | `db-package-lists-check.sh` | **Exit 1**, „Filter-Eintrag 'postgressnapshot' trifft nicht genau ein Paket von DB_COVERAGE_PKGS" |
| M4 | Paket aus dem Dockerfile-Filter | `postgressnapshot\|` gestrichen | dasselbe Skript | **Exit 1**, „Dockerfile-Filter weicht von DB_COVERAGE_PKGS ab", Diff nennt `postgressnapshot` |
| M5 | Paket aus der Messphase von `run-replication-tests.sh` | Zeile mit `./internal/adapters/driven/postgressnapshot \` gestrichen | dasselbe Skript | **Exit 1**, „Messlaeufe (store + replication) weicht ab", Diff nennt `postgressnapshot` |
| M6 | zweites Paket aus dem Filter | `\|replication/receive` gestrichen | dasselbe Skript | **Exit 1**, Diff nennt `replication/receive` |
| M7 | Coverage-Gate hält die Rücknahme | Dockerfile-Filter ohne `postgressnapshot` | `docker build --target coverage` (dieselbe Anweisung wie `make coverage-gate`, an der Kopie) | **BUILD_EXIT=1**, gedruckt `coverage-gate: FAIL — Coverage 78.60% unter Schwelle 80%` |
| M8 | Slot-Reserve-Test ist Lauf-Bedingung | Testname im `-run` von `run-replication-tests.sh tier` verschoben | Phase `tier` an der Kopie | **EXIT=1**, „testing: warning: no tests to run" und „TestSlotReserveExhaustedIsConfiguration ist nicht als PASS gelaufen" |
| M9 | Test läuft nur gegen `max_replication_slots=1` | `max_replication_slots=10` statt `1` im eigenen Container | Phase `tier` an der Kopie | **EXIT=1**, `--- FAIL: TestSlotReserveExhaustedIsConfiguration` |

### 5.2 Flake `TestStreamRestartsOnExistingSlot` (Paket `receive`, SQLSTATE 55006)

Aufbau (eigenes Skript im Scratchpad): Wegwerf-PostgreSQL mit der Tier-Konfiguration
(`wal_level=logical`, `max_wal_senders=10`, `max_replication_slots=10`,
`wal_sender_timeout=2000`), beide Testbinaries vorab gebaut; Restart-Test
`go test -count=100 -v -run '^TestStreamRestartsOnExistingSlot$'
./internal/adapters/driving/replication/receive` parallel zum Paket
`go test -count=3 -v ./internal/adapters/driven/postgressnapshot` (Lastfenster 24 bis 26 s,
gedruckt). „Ohne Poll" heißt: eine Kopie, in der genau die Zeile
`awaitSlotInactive(t, pool, env.slot, 10*time.Second)` gestrichen ist (Diff gegen `HEAD`:
eine Zeile).

| Version | ohne Poll | mit Poll (`HEAD`) | Paket `postgressnapshot` |
|---|---|---|---|
| PostgreSQL 18.6 | **3 von 100 rot**, Meldung `START_REPLICATION: ERROR: replication slot "slot_pgc_test_restart" is active for PID … (SQLSTATE 55006)` | **0 von 100** | 0 `--- FAIL` |
| PostgreSQL 17.11 | **3 von 100 rot**, dieselbe Meldung | **0 von 100** | 0 `--- FAIL` |

Die Größenordnung stimmt mit den Plan-Zahlen (9 bis 11 von 300 ohne, 0 von 300 mit Poll)
überein. Nur die Testdatei ändert sich (`git diff --stat`, §3.7). Container und Netze sind
danach entfernt (`docker ps -a`, `docker network ls`: nichts von mir übrig).

## 6. `.github/workflows/e2e.yml` — Diff-Urteil und offener Runner-Beleg

`git diff 62d00abe..HEAD -- .github/workflows/e2e.yml` (8 Zeilen): (1) der Kommentar vor
den Replication-Schritten wird umformuliert und nennt den Slot-Reserve-Lauf; (2) der
`name:` des letzten Schritts wird von `Replication-Tier (go test ./...)` zu
`Replication-Tier (go test ./... und Slot-Reserve)`. **Kein** `run:`, **keine** Matrix,
**keine** `uses:`-Zeile und damit keine Action-Pin ([`AGENTS.md`](../../AGENTS.md) §3.8
unberührt), keine Job-Abhängigkeit, kein neuer Schritt.

**Ist die Namensänderung eine „strukturelle Änderung" im Sinn von [`AGENTS.md`](../../AGENTS.md)
§3.10? Nein.** §3.10 zählt eine neue Matrix-Achse, einen neuen Schritt oder eine neue
Job-Abhängigkeit; ein reiner `name:` und ein Kommentar ändern weder das Verhalten noch die
Graph-Form. **Die Wirkung ist trotzdem real:** der bestehende Schritt `tier` ruft
`run-replication-tests.sh tier`, und dieses Skript fährt seit diesem Slice je Matrix-Leg
(PostgreSQL 17 und 18) einen zweiten Container (`max_replication_slots=1`, Bereitschaft
über `pg_isready` im Container) und den Slot-Reserve-Test. Ob das auf dem GitHub-Runner
ebenso grün läuft wie lokal (rund 45 s für die Phase), ist **weiter offen** bis zum ersten
realen Post-Push-Lauf von `e2e.yml` (`gh run list`/`gh run view`). **Der Plan führt es so:**
§6 letzte Risiko-Zeile („Der Tier-Schritt in `e2e.yml` fährt den Slot-Reserve-Lauf je
Matrix-Leg … bis zum ersten realen Lauf nicht geprüft … Ausgang: *(bei Closure)*"), und
§3 (`.github/workflows/e2e.yml`-Zeile) nennt die Wirkung ohne Workflow-Änderung ausdrücklich.
Der Workflow ist nicht blockierend (eigener Workflow, kein Required-Status-Check). Ich habe
den Runner-Lauf nicht gefahren und nicht behauptet; die Slice-Closure darf das Risiko erst
nach dem realen Lauf auflösen.

## 7. Harte Regeln

- **[`AGENTS.md`](../../AGENTS.md) §3.3** — die zwei Lifecycle-Moves sind reine Renames:
  `git show --stat -M 7673e774` (`{open => next}/… | 0`) und `48aa388a`
  (`{next => in-progress}/… | 0`); die Inhaltsänderung „Verantwortlich" (`1ae3fe3f`) liegt in
  einem eigenen Commit dazwischen.
- **§3.5** — keine bestehende `Accepted`-ADR verändert (`make doc-immutable` grün); die
  Reports unter `docs/reviews/` sind je in einem Commit angelegt und danach unverändert
  (`git log -- docs/reviews`: je ein Commit).
- **§3.6** — der Diff **verschärft**: `coverage.mk` ruft zusätzlich das Listen-Skript,
  `THRESHOLD` und `DB_COVERAGE_THRESHOLD` bleiben unverändert (`git diff … | grep -i
  THRESHOLD` zeigt nur Kontextzeilen und den Kopfkommentar), der Nenner der Stufe wächst um
  das Unterpaket, das Paket `postgressnapshot` wird mit belegter Skip-Eigenschaft (21/21)
  und gemessener Rücknahme (M7) ausgenommen.
- **§3.2** — kein `//nolint` im Go-Diff (`git diff -- '*.go' | grep -c nolint` = 0).
- **§3.7** — diff-skopierter Kandidatenlauf über die `+`-Zeilen von `*.go`, `*.sh`,
  `*.mk`, `Dockerfile` und `.github` (Chronik-/Vorher-Vokabular: `früher`, `zuvor`,
  `nicht mehr`, `jetzt`, `Fixrunde`, `Review`, `Befund`, `F-<n>`, `slice-`, `welle-`,
  `wäre`/`würde`/`hätte`, `ersetzt`, `wurde`, `statt`): neun Treffer, alle Indikativ oder
  Fehlermeldungen („roleDSN ersetzt Nutzer und Passwort", „waitSlotGone wartet, bis …",
  „Testcontainer … wurde nicht bereit"); der Herkunftsanker „gemessen an der Version dieses
  Servers" ist eine Ursprungsangabe. Längste Kommentarblöcke der neuen Dateien: 17 Zeilen
  (Paket-Doku `snapshot.go`), 13 (`db-package-lists-check.sh`), 11 (Port-Doku); die
  Kopfkommentare von `db-coverage.sh` (43 Zeilen) und `run-replication-tests.sh` (26) sind
  Bestand, der Diff verkürzt den ersten und ersetzt Zahlen durch einen Verweis. Der lange
  Vertrag liegt in den Sensor-Dokumenten und im Plan. Sensor-Doku-Prosa: die Chronik-Formen
  aus dem Nachprüfungs-Befund F-5 sind durch Zusammensetzungs-Aussagen mit Lauf-Ursprung
  ersetzt; die Zeile „an deren Stand `d7539e2c` im selben Verfahren 1691 gedeckt" nennt einen
  Stand als Ursprung einer Zahl (Instanz A), keine Vorher-Nachher-Erzählung. Kein Befund.
- **§3.11** — `docs-check` 0 Befunde (`hostpaths`); der Report trägt keinen host-lokalen
  Pfad, auch nicht in Grep-Mustern (Mutationen und Aufbauten sind mit Repo-relativen Pfaden
  beschrieben).
- **§3.12** — Zahlen im Plan und in den Sensor-Dokumenten gegen die Messung geprüft (§3.3,
  §3.4, §3.6): Nenner exakt (2040, 2082, 813, 122, 42, 245 × 3), gedeckte Unit-Zahl
  lauf-gebunden mit benannter Spanne (V-1).
- **§3.13** — die Suchlauf-Felder (Erstlieferung, Fixrunde, zweite Fixrunde) tragen
  Gefundenes und Nichtgefundenes je Träger, beide Stände; fremde Träger sind im Plan benannt,
  die zwei offenen Pläne (`slice-transformationen-kern-rename`,
  `slice-transformationen-antragsweg-usecase`) sind minimal gezogen (Diff: 4 Zeilen, je ein
  Aufzählungspunkt um `postgressnapshot`); im Baum bleiben keine weiteren offenen Pläne, die
  die Ausschlussliste nennen (`git grep -n -E 'postgresack|replication/receive' HEAD --
  docs/plan/planning` ohne `done/` und `observations/`: genau diese zwei Pläne).
- **Handbuch unberührt** — `git diff --name-only 62d00abe -- docs/user/benutzerhandbuch.md`
  ist leer; der Slice führt keine Betreiber-Oberfläche ein (keine `CDC_*`-Variable, keine
  `cdc.*`-Funktion, kein Endpunkt); die Testvariable `CDC_SNAPSHOT_TEST_EXCLUSIVE_DSN` ist
  keine Betreiber-Variable.
- **Commit-Traceability** — 19 Commits, jede Message nennt `LH-*` oder `ADR-*`, keine
  `SPEC-`/`ARC-`-Kennung im Betreff (§1).
- **[`AGENTS.md`](../../AGENTS.md) §3.1** — der Diff trägt keine Host-Toolchain; die Tier-
  Läufe und der Skip-Lauf laufen im gepinnten Toolchain-Image.

## 8. Befunde

Kein V-Befund blockiert. V-1…V-8 sind INFO ohne erwartete Aktion im Slice.

- **V-1 (INFO) — gedeckte Unit-Zahl lauf-gebunden, Nenner exakt.** Plan §3 (Messbefunde)
  und `coverage-gate.md` §Zählbasis nennen gedeckt 1691 (Stand `d7539e2c`) und 1731 (Diff);
  ich messe 1692 (`62d00abe`), 1689 (`d7539e2c`) und 1733 (`HEAD`), `make coverage-gate`
  druckt `83.10%`. Die Nenner 2040/2082 sind exakt; beide Dokumente führen die gedeckte
  Zahl ausdrücklich als lauf-gebunden (Spanne der Ausdrucke 83,1 bis 83,3 %). Kein Nachzug.
- **V-2 (INFO) — stand-relative Zellen im Suchlauf-Feld.** Die Zahlen gelten für den
  genannten Diff-Stand: `::tex[t]` „dieser Plan 0" (am Diff-Stand `bf05067d`; an `HEAD`
  trägt der Plan eine Zeile, die Beschreibung der Mutation im Messbefund zur Parität, kein
  Träger des Lese-Wortlauts), `db-adapter-coverage.md` 14 (an `HEAD` 15). Kein Widerspruch;
  Klasse „Stand-relative Zelle" wie im Report zu `slice-backfill-change-origin`.
- **V-3 (INFO) — zweite Fixrunde ohne eigenen Review-Report.** Die DoD-Zeile „Review
  durchgeführt" ist durch die zwei Reports erfüllt; die Behebung der Nachprüfungs-Befunde
  habe ich selbst am Code und an den Läufen nachgemessen (§3.7). Wer eine dritte
  Reviewer-Durchsicht der Fixrunde-2-Änderungen (Testdaten, Poll, Beschreibungen) verlangt,
  verlangt sie als eigenen Schritt; die Zeile hängt nicht daran.
- **V-4 (INFO) — Runner-Beleg von `e2e.yml` offen.** §6: bis zum ersten realen
  Post-Push-Lauf weiter offen, im Plan als Risiko geführt; Auflösung Sache der Closure.
- **V-5 (INFO) — DoD-Zeile „Reconciliation-Register — entfällt" auf `[ ]`.** Beim
  Closure-Nachzug einheitlich behandeln (im Vorgänger-Slice `row-image-gemeinsam` auf `[x]`).
- **V-6 (INFO) — Nebenwirkung der DB-Tiers.** `make test-store` merget mit dem im
  `DB_COVERAGE_DIR` liegenden Replication-Profil und druckt deshalb schon vor einem eigenen
  Replication-Lauf 651/813; `make test-replication` überschreibt `tools/schema/plan.yaml`
  (Zielzeile), ich habe die Datei nach jedem Lauf per `git checkout` zurückgenommen.
  Eigenschaft der Tier-Pfade, nicht dieses Slice (wie V-6 im Report zu
  `slice-backfill-change-origin`).
- **V-7 (INFO) — Extension-Typen im gemeinsamen Testcontainer.** Die drei Extensions
  liegen datenbankweit im gemeinsamen Testcontainer, in eigenem Schema und mit
  `DROP EXTENSION … CASCADE` im Cleanup; kein anderer Test des Trees nutzt sie (Aufräum-Messung
  §3.6: 0 Rückstände). Zur Kenntnis für künftige Tests, die dieselben Extensions anlegen.
- **V-8 (INFO) — Zählung des Flakes kleiner als beim Implementer.** 3 von 100 (18 und 17)
  gegen 9 bis 11 von 300: gleiche Größenordnung, andere Lastdauer (`-count=3` gegen `-count=6`
  des Pakets `postgressnapshot`); mit Poll je 0.

## 9. Ergebnis

| Prüfpunkt | Ergebnis |
|---|---|
| DoD §2 — `[x]`-Zeilen (11) | **11 von 11 erfüllt**, je mit eigenem Beleg (Zeile 9 mit V-3, Zeile 10 mit V-1/V-2) |
| DoD §2 — „Review durchgeführt" | **erfüllt** (zwei Reports; die Fixrunde-2-Behebungen selbst nachgemessen) |
| DoD §2 — übrige `[ ]`-Zeilen (5) | **5 von 5 korrekt offen** (Closure-/Rollen-Sequenz; Reconciliation „entfällt", V-5) |
| Plan-vs-Code-Diff | **deckungsgleich**: jede Zeile der Plan-Tabelle §3 ist im Diff vertreten, `seam.go` ist als „nicht realisiert" begründet und nicht im Diff; ungeplant nichts außerhalb der benannten Dateien (24 Dateien im Diff; außer den Reports, `ADR-0115` samt Indexzeile und dem Plan selbst steht jede in der Plan-Tabelle) |
| [`ADR-0115`](../plan/adr/0115-backfill-spaltenwerte-text-ergebnisformat.md) (Festlegungen 1–5, Fitness Function) | **bestätigt**: Paritätstest PostgreSQL 17.11 und 18.6 grün (85 Typ-Spalten, drei Lagen, 774 + 9 Vergleiche je Version); Mutationen M1 und M2 rot |
| Typ-Satz gegen die ADR | **vollständig**; `lc_monetary`-Lage auf Alpine nicht herstellbar (gemessen) |
| Snapshot-Träger, Klassen, `reltuples` | **bestätigt** (17.11 und 18.6) |
| Gate-Nenner, `snapshotlogic`, DB-Adapter-Coverage | **bestätigt**: 2040 → 2082, 42/42; 80.07 % (651/813) auf 17 und 18; Rücknahme rot (78.60 %) |
| Listen-Prüfskript | **bestätigt**: vier Mutationen Exit 1 |
| Tier (Slot-Reserve, Laufzeit) | **bestätigt**: EXIT=0 auf 17 und 18, Phase rund 45 s, zwei Mutationen rot |
| Flake `TestStreamRestartsOnExistingSlot` | Fix im Test (nur `stream_test.go`); 3/100 ohne Poll, 0/100 mit Poll auf 17 und 18 |
| `e2e.yml` | nur Kommentar und `name:`; §3.10-Runner-Beleg **weiter offen**, im Plan geführt |
| Harte Regeln (§3.2, §3.3, §3.5, §3.6, §3.7, §3.9, §3.11, §3.12, §3.13, Handbuch) | **erfüllt** |
| Commit-Traceability, MR-Immutabilität | **grün** (`RANGE=62d00abe..HEAD`) |
| Gates | **`make gates` EXIT=0**, `make test` EXIT=0, `make a-check` EXIT=0, `make coverage-gate` EXIT=0 (`83.10%`), `make test-store` EXIT=0, `make test-replication` EXIT=0 (PostgreSQL 18 und 17) |

## Verdikt

**Bestätigt.** Die DoD-Behauptung des Implementers ist durch eigene Messung belegt: der
Snapshot-Träger (M1–M4), die Blöcke und Spalten, die Schätzung (`reltuples = −1` auf 17.11
und 18.6), die Klassen und die Bild-Parität über den 85-Spalten-Typ-Satz in drei GUC-Lagen
gegen das reale WAL-Bild halten auf beiden Major-Versionen; die Lese-Anweisung trägt weder
Cast noch Funktion und `BuildRowImage` ist unverändert. Die Mutationen `::text` an den
Spaltennamen (rot für `bo`/`ch1`/`ch5`/`db`/`ia`/`xm`) und die Veränderung der Extension-
Wertseite (rot für `hs`/`ah`) färben den Tier auf 17 und 18 rot, ebenso Streichungen in den
namentlichen Paketlisten, die Rücknahme aus dem Filter (78.60 %) und die zwei Tier-
Mutationen des Slot-Reserve-Laufs. Das Gate verschärft (Nenner 2040 → 2082 mit dem
Unterpaket, Listen-Prüfskript vor dem Bau), keine Schwelle sinkt. Kein Befund blockiert;
V-1…V-8 sind INFO. Der Runner-Lauf von `e2e.yml` (§6, V-4) bleibt bis zum ersten realen
Post-Push-Lauf offen und steht so im Plan.

**Übergabe:** Verifier → Planner. Offen bleibt die Closure: Closure-Notiz mit Lerneintrag
(die Finding-Klassen beider Reviews, darunter „Tatsachenbehauptung im Träger widerspricht
der Messung" und „Flake im geteilten Testcontainer ohne Träger", in das
Beobachtungs-Register), Ausgänge der zehn §6-Risiken (das `e2e.yml`-Risiko erst nach dem
realen Lauf), Reconciliation-Zeile (V-5), Welle-Paarungen; der sachliche Befund
„`receive` wiederholt `START_REPLICATION` bei 55006 nicht" gehört ins Beobachtungs-Register
(Sache der Closure, im Plan so benannt). Danach der reine `git mv` nach `done/`
([`AGENTS.md`](../../AGENTS.md) §3.3, Fall 2).
