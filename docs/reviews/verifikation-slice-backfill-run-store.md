# Verifikations-Report: slice-backfill-run-store — 2026-09-24

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich + Entscheidungs-Konformität +
Plan-vs-Code-Diff + Gates. Review-Artefakt des Reviewers:
[`review-slice-backfill-run-store.md`](review-slice-backfill-run-store.md);
Formvorbild dieses Reports:
[`verifikation-slice-backfill-snapshot-reader.md`](verifikation-slice-backfill-snapshot-reader.md).

**Gegenstand:** Slice-Plan `slice-backfill-run-store` (Welle
`welle-backfill-bestand`), Diff-Range `afe2392f..HEAD` (`552cb437`), 11 Commits:
Lifecycle und Verantwortlich (`86347a9c`, `38d4031d`, `03cd3804`), Umsetzung
(`6c7dbbb4` Tabelle, Grants, Rollen-Test der Rollout-Datei; `92603002` Adapter;
`d095e5be` CHECK-Test, Ordnungs-Soll, Kommentare), Plan-Nachzüge (`9bfd8246`,
`54fc81d3`), Review-Report `51c69243`, Fixrunde (`bd470a0b` Code, Tests, Handbuch,
Harness-Doku; `552cb437` Pläne). Dieser Lauf ändert weder Code noch Plan, Spec oder
Doku; er schreibt nur diesen Report. Der Arbeitsbaum blieb während des Laufs sauber
(`git status --short` leer nach jedem Tier-Lauf, nachdem `tools/schema/plan.yaml` und
`tools/schema/down.sql` per `git checkout` zurückgenommen waren). Alle Mutationen
liefen an einem Wegwerf-Klon (`git clone` des Repos in ein Scratchpad-Verzeichnis,
`git status` des Repos danach leer), nicht im Arbeitsbaum.

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — [`AGENTS.md`](../../AGENTS.md) §3.9)

Jeder Lauf schrieb in eine Log-Datei; der Exit-Code wurde im selben Aufruf gesondert
gesichert (`make … > <log> 2>&1; echo EXIT=$?`), die Logs danach gelesen.

| Sensor | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `make gates` | **EXIT=0** | `baseline-verify: v6.9.0 OK — 54 Dateien` · `db-package-lists-check: OK — … dieselben 4 Pakete` · `coverage-gate: OK — Coverage 84.70% erfüllt Schwelle 80%` · `d-check: 1033 Datei(en) geprüft, 0 Befund(e)` · `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD"` · `generated-sync: OK` · a-check `gesamt: 0 Befund(e)` |
| `make test` (Race-Detector) | **EXIT=0** | 42 Pakete `ok`, 0 `FAIL` (`postgresstorage`, `mapper`, `sqlexec`, `usecase/backfill`, `internal/bootstrap`, `domain/model` ok) |
| `make a-check` | **EXIT=0** | `gesamt: 0 Befund(e)` |
| `make coverage-gate` (eigener Lauf) | **EXIT=0** | `total: (statements) 84.8%`, `coverage-gate: OK — Coverage 84.80% erfüllt Schwelle 80%` — der Lauf im `make gates` druckte `84.70%` (lauf-gebundene Streuung, V-6); `mapper/backfillrun.go` (`ToBackfillRun`, `TimeArgument`, `EstimateArgument`, `PositionArgument`) und `sqlexec/translate.go` `ReadBackfillRuns` stehen im Profil mit je 100,0 % |
| `make test-store` (PostgreSQL 18.6, Repo-Default) | **EXIT=0** | `ok … postgresstorage 6.907s`; `DB-Adapter-Coverage: 81.69% (gedeckt 830 von 1016 Statements; Profile gemergt: store,replication)` |
| `make test-replication` (PostgreSQL 18.6) | **EXIT=0** | dieselbe Zeile `DB-Adapter-Coverage: 81.69% (gedeckt 830 von 1016 Statements; …)`, `db-coverage: OK — … erfuellt Schwelle 70%` |
| `PG_TEST_IMAGE=postgres:17-alpine@sha256:7456ef82… make test-store` (17er-Digest aus `.github/workflows/e2e.yml` Z. 80) | **EXIT=0** | `ok … postgresstorage 6.936s`, `DB-Adapter-Coverage: 81.69% (gedeckt 830 von 1016 …)` |
| `PG_TEST_IMAGE=postgres:17-alpine@sha256:7456ef82… make test-replication` | **EXIT=0** | dieselbe Zeile `81.69% (gedeckt 830 von 1016 …)` |
| `bash tools/harness/run-schema-rollout-guard-test.sh` (ungekürzt) | **EXIT=0** | `Lauf 5 OK — Tag v0.1.2: Exit 0 (Rollout des Tags), Exit 0 (Arbeitsbaum, mit Vorlauf), Exit 0 (Arbeitsbaum, zweiter Lauf); Zeile alttag-ch über cdc.changes lesbar; cdc_admin-Rechte auf administration_request und backfill_run gesetzt`; Schlusszeile `OK — alle Belege real erbracht (…)` |
| `bash tools/harness/db-package-lists-check.sh` | **EXIT=0** | `OK — Dockerfile-Filter, Messlaeufe und DB_COVERAGE_PKGS nennen dieselben 4 Pakete`; `git diff afe2392f..HEAD` über `Dockerfile`, `harness/mk`, `db-coverage.sh`, `run-replication-tests.sh`, `db-package-lists-check.sh` ist leer (die drei namentlichen DB-Listen unverändert) |
| `make commit-traceability RANGE=afe2392f..HEAD` | **EXIT=0** | `OK — 11 Commit(s) in "afe2392f..HEAD", Betreffs ohne Struktur-ID`; `git log --format=%s afe2392f..HEAD` enthält in keinem Betreff `SPEC-`/`ARC-` und in jedem `LH-` oder `ADR-` |
| `make doc-commits RANGE=afe2392f..HEAD` | **EXIT=0** | 1033 Dateien, 0 Befunde |
| `make doc-immutable RANGE=afe2392f..HEAD` | **EXIT=0** | 1033 Dateien, 0 Befunde |
| `make doc-trace` (advisory, nur zur Kenntnis) | EXIT=0 | 80 Anforderung(en), 3 Waise(n): [`LH-FA-CAP-009`](../../spec/lastenheft.md), `LH-FA-CFG-007`, `LH-FA-CFG-008` — die Backfill-Waise bleibt bis zum Liefer-Slice erwartet |

Die Store-Zeile von `make test-store` merget mit dem zuletzt im `DB_COVERAGE_DIR`
liegenden Replication-Profil und druckt deshalb schon vor einem eigenen
Replication-Lauf dieselbe Zahl (V-6); die Zahl 830/1016 ist der Beleg der
Läufe `make test-replication` auf PostgreSQL 18 und 17 dieses Verifier-Laufs.

## 2. DoD — Verdikt je Zeile (§2 des Plans)

| # | DoD-Zeile | Verdikt | Realer Beleg |
|---|---|---|---|
| 1 | Schema und Grants (`schema.yaml`, `nacharbeit-roles.sql`, Rollout zweimal, Alt-Tag-Lauf, `plan.yaml`/`down.sql`, Rollen-Test (4a), Store-Rollen-Test) | **erfüllt** | §3.1 bis §3.3: `backfill_run` mit den 14 Spalten von [`SPEC-029`](../../spec/pflichtenheft.md), CHECK bei Erstanlage, `estimated_rows` nullable, Warn-Spalten `boolean NOT NULL DEFAULT false`; Grants gemessen (§3.2); Guard-Lauf 5 zweimal Exit 0 (§1); `plan.yaml`/`down.sql` gegen frisches Rollout gleich bis auf die `target`-Zeile (§3.3); Rollen-Tests `TestCdcCaptureRoleUpdatesBackfillRunButNeitherInsertsNorDeletes`, `TestCdcAdminRoleInsertsBackfillRunButNeitherUpdatesNorDeletes`, `TestCdcReaderRoleCannotReadTheBackfillRunBaseTable` grün auf 17 und 18 |
| 2 | Annahme-Adapter (queued + applied zugleich, kein Vermerk ohne `pending`, zweiter Antrag, `NULL` ≠ 0) | **erfüllt** | `TestBackfillAdmitLeavesQueuedRunAndAppliedRequestTogether`, `…RollsBackWhenTheRequestIsNotPending` (kein Antrag, `applied`, `failed`), `…RefusesASecondActiveRunOfTheSameTable` (`queued` und `running`), `…StoresUnknownEstimateAsNullAndNotAsZero` (`estimated_rows IS NULL` und `= 0` per SQL); `make test-store` EXIT=0 auf 17 und 18; Mutationen M2, M7, M8 rot (§4) |
| 3 | Atomarität (zweiter Leser, Rollback, `origin`/`operation`/`old_data`) | **erfüllt** | `TestBackfillWriterIsInvisibleBeforeTheCommitAndCompleteAfterIt` (vor dem Commit 0/0 Zeilen, danach 2/5; `origin IS DISTINCT FROM 'backfill' OR operation <> 'INSERT' OR old_data IS NOT NULL` zählt 0), `TestBackfillWriterRollbackLeavesNoRow`, `TestBackfillWriterDoesNotCommitAnEndedRun`; M3 rot |
| 4 | Ordnung und Zustand (WAL-Commit auf derselben Position, Übergänge, Fortschritt außerhalb, Abgleich; Kollation) | **erfüllt** | `TestBackfillBlocksSortBeforeTheWALTransactionOnTheSamePosition` (Lesezugriff des Store-Adapters **und** `cdc.changes` mit `ORDER BY commit_position, transaction_id, sequence`), gedruckt `datcollate der Test-Datenbank: en_US.utf8` auf 17.11 und 18.6; `TestBackfillProgressDoesNotWaitForTheOpenWriteTransaction`, `TestBackfillRunFollowsTheRunThroughItsTransitions`, `TestBackfillRunInterruptRunningTouchesOnlyRunningRunsOfTheSource`; Ordnung selbst an der DB nachgemessen (§3.4) |
| 5 | Adapter-Pflichten (Zeitbegrenzung, Endzustand idempotent) | **erfüllt** | netzlos `TestBackfillRunOperationsEndAtTheirOwnTimeoutOnDetachedContext` und `TestBackfillAdmissionAndWriterEndAtTheirOwnTimeoutOnDetachedContext` (blockierende Naht, Frist 40 ms, abgelöster Kontext); `TestBackfillRunFinishOnAnEndedRunIsANoOpSuccess` (`Finish` auf `completed`/`failed`/`interrupted` → nil, Zeile unverändert); M1 und M4 rot |
| 6 | Rollen unter echtem Login (Review F-1, F-2) | **erfüllt** | §3.2: Rechte gemessen; `TestAdministrationPathRunsUnderLeastPrivilegeLogins` (`internal/bootstrap`), `TestBackfillAdaptersRunUnderTheirRoles`, `TestBackfillAdaptersFailUnderTheOtherRole` grün auf 17 und 18 (eigener `-v`-Lauf); die Logins sind `IN ROLE`, gemessen `rolsuper = f`, 0 Objekte im Eigentum (V-4); (4b) und Guard-Lauf 5 rot unter der Grant-Mutation (M11) |
| 7 | `make gates` grün, Exit gesondert | **erfüllt** | eigener Lauf EXIT=0 (§1) |
| 8 | Review durchgeführt, Report unter `docs/reviews/` | **erfüllt, mit V-3 (INFO)** | Report liegt vor (1 HIGH, 1 MEDIUM, 2 LOW, 5 INFO); das neue Test-/Doku-Material der Fixrunde habe ich selbst geprüft (§3.5), einen zweiten Review-Report dazu gibt es nicht |
| 9 | §3.13-Suchlauf (Feld in §3, beide Stände) | **erfüllt, mit V-1 (LOW) und V-2 (INFO)** | die Zahlen beider Felder nachgemessen (§3.6): alle stimmen bis auf „Diff sieben Dateien" (gemessen sechs); der bewegte Träger „DB-Adapter-Nenner" steht nicht im Feld (V-1) |
| 10 | Doku-Update (`harness/README.md`, Schema-Kommentare, Handbuch §2/§4/§5) | **erfüllt** | §3.7: Sensor-Zeile `make test-store` trägt Backfill-Tests, Rollen, Adapter unter Login; Handbuch `Version: 1.47` und Zeile 1.47 in der Änderungshistorie; `schema.yaml`-Kopf trägt `backfill_run` |
| 11 | Closure-Notiz mit Lerneintrag | **korrekt offen** | §7 des Plans trägt „*(zu tragen bei Closure)*" |
| 12 | Reconciliation-Register — entfällt | **korrekt offen / entfällt** | Greenfield; Zeile steht auf `[ ]` wie in den Vorgänger-Slices (unverändert) |
| 13 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | Closure-Pflicht |
| 14 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | jede der zehn Risiko-Zeilen trägt „**Ausgang:** *(bei Closure)*" (gezählt) |
| 15 | Die drei Paarungen | **korrekt offen** | hängen an der Closure von `welle-backfill-bestand` |

Kein `[x]` ohne Beleg, kein `[ ]`, das über die Rollen-Sequenz hinaus belegt wäre.
`[x]` sind zehn Zeilen (Nr. 1–10), `[ ]` fünf (Nr. 11–15), zusammen 15 — gezählt am Plan.

## 3. Kernaussagen, selbst gemessen

### 3.1 Schema der Tabelle

`git show HEAD:tools/schema/plan.yaml` trägt die `CREATE TABLE "backfill_run"`-Anweisung
mit `run_id`, `source_id` (FK `source`), `schema_name`, `table_name`, `status`,
`requested_at` (`NOT NULL DEFAULT CURRENT_TIMESTAMP`), `started_at`, `finished_at`,
`snapshot_position`, `rows_copied` (`NOT NULL DEFAULT 0`), `estimated_rows` (**nullable**),
`warn_estimated_size` und `warn_duration` (je `BOOLEAN NOT NULL DEFAULT FALSE`),
`error_message`, `PRIMARY KEY (run_id)` und `CONSTRAINT chk_backfill_run_status CHECK
(status IN ('queued', 'running', 'completed', 'failed', 'interrupted'))` — die 14 Spalten
und die Prüfbedingung von [`SPEC-029`](../../spec/pflichtenheft.md), die Menge **bei Erstanlage**. `down.sql` trägt
`DROP TABLE "backfill_run";`.

### 3.2 Rechte, selbst gemessen an Wegwerf-PostgreSQL 18.6 und 17.11

Aufbau: Wegwerf-Container, Schema-Rollout des Klons von `HEAD` über
`tools/schema/apply-rollout.sh` (Exit 0 auf beiden Versionen), Abfrage
`has_table_privilege(<rolle>, <tabelle>, <recht>)`; gedruckt identisch auf 18.6 und 17.11:

| Tabelle | Rolle | `SELECT` | `INSERT` | `UPDATE` | `DELETE` |
|---|---|---|---|---|---|
| `cdc.administration_request` | `cdc_admin` | **t** | f | **t** | f |
| `cdc.administration_request` | `cdc_capture` | f | f | f | f |
| `cdc.administration_request` | `cdc_reader` | f | f | f | f |
| `cdc.backfill_run` | `cdc_admin` | **t** | **t** | f | f |
| `cdc.backfill_run` | `cdc_capture` | **t** | f | **t** | f |
| `cdc.backfill_run` | `cdc_reader` | f | f | f | f |

`relacl` gedruckt: `administration_request` `{cdc=arwdDxtm/cdc,cdc_admin=rw/cdc}`,
`backfill_run` `{cdc=arwdDxtm/cdc,cdc_admin=ar/cdc,cdc_capture=rw/cdc}` — kein weiterer
Eintrag. `tools/schema/nacharbeit-roles.sql` trägt `GRANT SELECT, UPDATE ON
cdc.administration_request TO cdc_admin;` (Zeile 158, Kommentar davor), kein `INSERT`/`DELETE`.

**Login `IN ROLE cdc_admin`** (Rolle `vfy_adm`, `rolsuper = f`, 0 Objekte im Eigentum), direkt per SQL, auf 18.6
und 17.11: `SELECT administration_request_id … WHERE status='pending'` liefert `vr1`;
`UPDATE … SET status='applied' … AND status='pending'` → `UPDATE 1`; `INSERT` →
`ERROR: permission denied for table administration_request`; `DELETE` → `ERROR:
permission denied for table administration_request`. Der Bestandsdefekt von Review F-1
(kein Recht auf die Queue) ist behoben und gegen die Mutation „Grant streichen" gebunden
(§4, M11: `Admit` unter `cdc_admin` scheitert ohne den Grant mit `permission denied for
table administration_request (SQLSTATE 42501)`).

**Alt-Tag-Beleg** ([`ADR-0114`](../plan/adr/0114-schema-rollout-vorlauf-view-signatur.md) Entscheidung 7): Guard-Lauf 5 rollt das Schema von
`v0.1.2` (jüngstes `v*`-Tag im Repo, `git tag`) per `git archive` aus, prüft
vorab, dass `cdc_admin` dort kein `UPDATE` auf `administration_request` trägt, rollt danach den
Arbeitsbaum zweimal aus (Exit 0, Exit 0, die Zeile `alttag-ch` über `cdc.changes` unverändert
lesbar) und prüft danach die sieben Rechte (`administration_request` `SELECT` t / `UPDATE` t /
`INSERT` f / `DELETE` f, `backfill_run` `SELECT` t / `INSERT` t / `UPDATE` f) — gedruckt `Lauf 5 OK`.

### 3.3 `plan.yaml` und `down.sql` als Ergebnis eines Rollouts gegen eine leere DB

`make test-store` (Wegwerf-PostgreSQL 18.6, frische DB) erzeugte `plan.yaml` und `down.sql`
neu. Gegen die committeten Stände: `down.sql` **byte-gleich**; `plan.yaml` gleich bis auf die
`target`-Zeile (committet der Compose-Wert `postgres://postgres:***@cdc-test-postgres:5432/cdc?sslmode=disable`,
gedruckt im Tier `postgres://cdc:***@cdc-store-test-pg:5432/cdc_test?sslmode=disable`;
`git diff --stat` nannte genau diese eine Zeile). Der Compose-Wert steht auch im Parent
`afe2392f` und im Commit `5f126971`. Zweiter Rollout gegen ein migriertes Ziel: Guard-Läufe 2
und 5 (zweimal Exit 0), ohne neuen Eintrag in `knownForeignObjects` (`git diff` über
`tools/schema/rolloutguard` leer).

### 3.4 Ordnung `0bf-…` an der realen DB, selbst gemessen

`string_agg(v, ', ' ORDER BY v COLLATE …)` über `1, 748, 9999999999, 0bf-run-00000001,
0bf-a-00000002, a`, gedruckt:

| Datenbank | `datcollate` | Default | `"C"` | `"und-x-icu"` | `"en_US.utf8"` |
|---|---|---|---|---|---|
| `postgres:18-alpine` (18.6) | `en_US.utf8` | `0bf-a-…, 0bf-run-…, 1, 748, 9999999999, a` | dieselbe | dieselbe | Collation existiert nicht (`ERROR: collation "en_US.utf8" … does not exist`) |
| `postgres:17-alpine` (17.11) | `en_US.utf8` | dieselbe | dieselbe | dieselbe | Collation existiert nicht |
| `postgres:18` (Debian, glibc, 18.6) | `en_US.utf8` | dieselbe | dieselbe | dieselbe | dieselbe |

Die Reihenfolge `0bf-…` vor `1`, `748`, `9999999999` hält in allen gemessenen Lagen, auch unter
glibc `en_US.utf8`. Auf dem Alpine-Testcontainer ist `en_US.utf8` nur der gedruckte
`datcollate`-Name, keine wählbare Collation (musl; V-5). Der Test druckt `datcollate der
Test-Datenbank: en_US.utf8`.

### 3.5 Neuer Test-/Doku-Code der Fixrunde, selbst geprüft

- `internal/bootstrap/administration_roles_internal_test.go`: legt zwei Logins `IN ROLE cdc_admin` bzw.
  `cdc_capture` an, überträgt die Quelltabelle an `cdc_admin` (dokumentierte Betriebs-Vorbedingung,
  `nacharbeit-roles.sql` Zeilen 89–98), zieht je einen Antrag der vier Arten durch
  `processAdministrationRequests`, prüft **Status `applied`** je Antrag, die **Publication-Mitgliedschaft**
  (1 nach `enable_table`, 0 nach `disable_table`, aus `pg_publication_tables`) und einen fehlgeschlagenen
  Antrag (`failed` samt nicht leerem Fehltext); unter der Grant-Mutation meldet er
  `Status "pending" (), erwartet applied — Log: [WARN: administration: Anträge lesen fehlgeschlagen]`.
  Der Test prüft Erwartetes, nicht bloß „kein Panic". Die Nicht-Superuser-Eigenschaft ist per Konstruktion
  (`CREATE ROLE … LOGIN … IN ROLE`) gegeben, nicht per Assertion (V-4; von mir gemessen `rolsuper = f`).
- `postgresstorage/backfillroles_test.go`: `TestBackfillAdaptersRunUnderTheirRoles` führt Annahme
  (`cdc_admin`-Login), Run-Zustand und Schreiber (`cdc_capture`-Login) bis zum Commit (`completed`, 1
  Transaktion, 2 Changes), `Queued`, `Finish(failed)` und `InterruptRunning`;
  `TestBackfillAdaptersFailUnderTheOtherRole` prüft je Gegenrolle `SQLSTATE 42501`
  (`permissionDenied`) **und** dass nichts hinterlassen wird (keine Run-Zeile, Antrag bleibt `pending`, Run bleibt
  `running`, 0 Transaktionen).
- `roles_test.go` (Erweiterung) und `roles_rollout_file_internal_test.go` (4b): Rechte je Rolle real mit
  `SET ROLE`, erlaubt = Erfolg, verboten = 42501; (4b) prüft `SELECT`, `UPDATE` erlaubt und `INSERT`/`DELETE`
  sowie jedes Recht der beiden anderen Rollen verboten gegen den Rollout-Text; unter M11 rot mit der
  Meldung `cdc_admin fehlt SELECT auf cdc.administration_request im Rollout-Text …`.
- Guard-Skript Lauf 5: die sieben Rechte-Abfragen und die Vorbedingung; unter M11 rot mit `FEHLER — Lauf 5:
  cdc_admin trägt nach dem Upgrade über v0.1.2 auf cdc.administration_request das Recht SELECT nicht als t`.
- Handbuch 1.47: §2 (Rollen-Tabelle), §4 (Absatz „Rechte der drei Rollen"), §5 (`CDC_ADMIN_DSN`) und Zeile 1.47.
  Die Aussage „Ohne das Recht … bleiben Anträge `pending`, der Feed-Container protokolliert ‚Anträge lesen
  fehlgeschlagen'" habe ich am Test selbst gesehen (M11-Meldung oben); „setzt die Rechte bei jedem Lauf
  (idempotent)" trägt der zweite Guard-Lauf (Exit 0) und Lauf 2.

### 3.6 Zahlen der Suchlauf-Felder (§3.12), beide Stände selbst gemessen

| Plan-Zahl | Meine Messung | Ergebnis |
|---|---|---|
| Tabellenlisten: Parent `03cd3804` und Diff `d095e5be` dieselben elf Dateien (`benutzerhandbuch.md` 9, `e2e-abdeckung.md` 2, `harness/README.md` 1, `db-adapter-coverage.md` 2, `schema.sql` 2, `pflichtenheft.md` 5, `nacharbeit-administration.sql` 6, `-heartbeat.sql` 2, `-observability.sql` 7, `-roles.sql` 6, `schema.yaml` 4) | identisch an beiden Ständen (`git grep -c`, gleiche Pfad-Ausschlüsse) | bestätigt |
| Rollenverteilung: `nacharbeit-roles.sql` 23 → 29, sonst gleich | Parent 23, Diff 29; übrige Dateien gleich (`README.md` 1, `bench-abdeckung.md` 1, `benutzerhandbuch.md` 18, `e2e-abdeckung.md` 2, `harness/README.md` 2, `schema-rollout.md` 2, `pflichtenheft.md` 6) | bestätigt |
| Rollen-Test der Rollout-Datei: Parent kein Treffer, Diff `roles_rollout_file_internal_test.go` 5 | identisch | bestätigt |
| Test-Bereinigung: Parent sieben Dateien, Diff elf (`backfillhelpers_test.go` 3, `backfillrun_test.go` 5, `backfillwriter_test.go` 4, `roles_test.go` 4); zweiter Befehl am Diff-Stand null Treffer | identisch, zweiter Befehl Exit 1 (kein Treffer) | bestätigt |
| Fixrunde, Rechte-Vokabular: Parent 13 Dateien, Diff 14 (`nacharbeit-roles.sql` 13 → 14, `roles_rollout_file_internal_test.go` 5 → 6, neu `administration_roles_internal_test.go` 1) | Parent 13, Diff 14, gleiche Zellen | bestätigt |
| Fixrunde, `cdc_admin`: „Parent fünf Dateien, Diff **sieben**" (`benutzerhandbuch.md` 10 → 13, `harness/README.md` 0 → 1, `schema-rollout.md` 2 → 3) | Parent fünf, Diff **sechs** (`README.md` 1, `compose.yaml` 1, `benutzerhandbuch.md` 13, `harness/README.md` 1, `schema-rollout.md` 3, `pflichtenheft.md` 3); die drei genannten Bewegungen stimmen | Zahl „sieben" weicht ab (V-2) |
| Fixrunde, `administration_request`: Parent neun, Diff zehn (`benutzerhandbuch.md` 10 → 14, `harness/README.md` 0 → 1, `schema-rollout.md` 1 → 2) | identisch | bestätigt |
| Fixrunde, [`ADR-0047`](../plan/adr/0047-rollenspezifische-dsn-verdrahtung.md): Parent und Diff 7 | identisch | bestätigt |
| Fixrunde, „drei Views": Parent drei Treffer, Diff einer (`schema.yaml` Zeile 39) | identisch | bestätigt |
| Coverage-/DB-Adapter-Zahlen im Plan | der Plan nennt keine Prozentzahl; §6 verlangt „der Bericht nennt die Zahl mit ihrem Lauf" (§1 dieses Reports) | — |

## 4. Mutationen (Eingabeseite, selbst gesehen)

Aufbau: Wegwerf-Klon von `HEAD`, alle Mutationen zugleich angewendet (die getesteten Pfade
sind disjunkt; zugeordnet über die roten Testnamen), `bash tools/harness/run-store-tests.sh`
im Klon. Im Klon habe ich zusätzlich die zwei `go test`-Zeilen der Pakete `internal/bootstrap` und
„übrige Pakete" mit `|| true` versehen, damit der Lauf nach einem roten Paket weiterläuft; der erste Versuch
brach an der Mutation M7 im netzlosen `mapper`-Paket ab (`set -e`), der zweite lief durch. Der Arbeitsbaum
des Repos wurde nicht verändert. Ergebnis des zweiten Laufs: `EXIT=1`, `--- FAIL` für:

| # | Zusage | Mutation | Rot gesehen (Testname, gedruckte Meldung) |
|---|---|---|---|
| M1 | `Finish` wirkt nur auf `queued`/`running` | `UpdateBackfillRunFinish`: `status IN ('running', 'completed')` | `TestBackfillRunFinishOnAnEndedRunIsANoOpSuccess`: `completed: die Zeile änderte sich: failed … rows 99` |
| M2 | Annahme prüft genau eine `pending`-Zeile | `Admit`: `RowsAffected`-Prüfung entfernt | `TestBackfillAdmitRollsBackWhenTheRequestIsNotPending`: `Antrag besteht nicht: Fehler = <nil>, erwartet ErrBackfillRequestNotPending` |
| M3 | ein Commit am Ende | `Commit`: `t.tx.Rollback(ctx)` statt `t.tx.Commit(ctx)` | `TestBackfillWriterIsInvisibleBeforeTheCommitAndCompleteAfterIt`, `TestBackfillProgressDoesNotWaitForTheOpenWriteTransaction`, `TestBackfillBlocksSortBeforeTheWALTransactionOnTheSamePosition` |
| M4 | jede Operation zeitbegrenzt | `Finish`: `boundedContext` entfernt | `TestBackfillRunOperationsEndAtTheirOwnTimeoutOnDetachedContext` (3,00 s = Obergrenze der Naht): `Finish: Fehler = … fake: Operation ohne Frist blockiert, erwartet die eigene Frist` |
| M5 | Blockvertrag: Herkunft `backfill` | `blockRows`: `case change.Origin != …` → `case false:` | `TestBackfillWriterRejectsBlocksOutsideTheContract` (Fall „Herkunft wal") |
| M6 | Blockvertrag: Blocknummern 1, 2, 3, … | `blockRows`: `if false && block.ID != expected` | `TestBackfillWriterRequiresConsecutiveBlockNumbers`, `TestBackfillWriterRejectsAContractViolationAndStaysUsable`, `TestBackfillWriterRejectsBlocksOutsideTheContract` (M5 und M6 zugleich angewendet, je Testname zuordenbar; einzeln habe ich sie nicht gefahren) |
| M7 | `NULL` = unbekannt, nie 0 | `EstimateArgument`: `int64(0)` für „unbekannt" | `TestBackfillAdmitStoresUnknownEstimateAsNullAndNotAsZero` und netzlos `TestBackfillArgumentsCarryNullAndValues`: `unbekannte Schätzung = 0, wollen nil` |
| M8 | kein zweiter aktiver Run | `SelectActiveBackfillRun`: nur `status IN ('queued')` | `TestBackfillAdmitRefusesASecondActiveRunOfTheSameTable` |
| M9 | `MarkRunning` nur aus `queued` | `UpdateBackfillRunRunning`: `WHERE run_id = $1` | `TestBackfillRunMarkRunningAndProgressRequireTheirSourceState` |
| M10 | Aufnahme-Ordnung `(requested_at, run_id)` | `SelectQueuedBackfillRuns`: `ORDER BY requested_at` | `TestBackfillRunQueuedIsOrderedByRequestedAtThenRunID` |
| M11 | `cdc_admin` `SELECT, UPDATE` auf `administration_request` | Grant-Zeile aus `nacharbeit-roles.sql` gestrichen | `TestAdministrationPathRunsUnderLeastPrivilegeLogins` (`Status "pending" …`), `TestCdcAdminRoleReadsAndMarksAdministrationRequestsButNeitherInsertsNorDeletes` (`SELECT … permission denied … SQLSTATE 42501`), `TestBackfillAdaptersRunUnderTheirRoles` und `TestBackfillAdaptersFailUnderTheOtherRole` (`Admit … permission denied for table administration_request (SQLSTATE 42501)`), netzlos `TestRolloutDateiTraegtDieRechteDerVerdrahtung` (4b) |
| M11b | Alt-Tag-Beleg bindet den Grant | dieselbe Grant-Streichung im Klon, `bash tools/harness/run-schema-rollout-guard-test.sh` | **EXIT=1**, gedruckt `FEHLER — Lauf 5: cdc_admin trägt nach dem Upgrade über v0.1.2 auf cdc.administration_request das Recht SELECT nicht als t` (Läufe 1–4 vorher grün) |

Zurücknahme: Die Mutationen liegen nur im Klon; `git status --short` und `git diff` des
Repos sind leer, keine Container/Netze übrig (`docker ps -a`, `docker network ls`: nichts von mir).

## 5. Entscheidungs-Konformität

| Entscheidung | Festlegung | Beleg | Verdikt |
|---|---|---|---|
| [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 4 | Ein-Transaktions-Form, Run-Zeile `completed` als letzte Anweisung vor dem Commit, Fortschritt außerhalb der Daten-Transaktion | `backfillwriter.go` (`Commit` schreibt `UpdateBackfillRunCompleted` vor `tx.Commit`; `AppendBlock` berührt `backfill_run` nicht); Tests `…WaitForTheOpenWriteTransaction`, M3 | konform |
| [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 6 | Ordnung `0bf-…` vor WAL-Kennungen derselben Position (Textsortierung) | §3.4 und `TestBackfillBlocksSort…` (Lesezugriff und View) | konform, Grenze der Kollation gemessen |
| [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 7 | Retention ohne Sonderpfad | der Diff berührt keine Retention (`git diff --stat` ohne Retention-Dateien) | konform |
| [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 1 | `cdc_admin` legt `queued` an, in **derselben Transaktion** wie `pending → applied`, Prüfung „kein aktiver Run" davor; Grants `SELECT, INSERT` / `SELECT, UPDATE`, niemand `DELETE`, `cdc_reader` nichts | `backfilladmission.go` (`Begin`, `SelectActiveBackfillRun`, `InsertBackfillRun`, `UpdateAdministrationRequestApplied` mit `RowsAffected`-Prüfung, `Commit`); Rechte gemessen (§3.2); M2, M8 | konform |
| [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 1 Punkt 2 | keine Sperre zwischen Annehmenden (eine Goroutine je Quelle) | im Code keine Sperre und keine Unique-Kante; Plan §6 führt es als benannte, nicht erzwungene und nicht getestete Annahme (Review F-5) | konform (Annahme der ADR, benannt); ich habe zwei gleichzeitig annehmende Verbindungen **nicht** gemessen |
| [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 3 | `estimated_rows` = `NULL` heißt unbekannt; zwei Warn-Spalten mit je einem Schreiber | Schema §3.1; M7; `warn_estimated_size` schreibt nur `InsertBackfillRun`, `warn_duration` nur die Worker-Anweisungen (`warn_duration OR $n`) | konform |
| [`ADR-0114`](../plan/adr/0114-schema-rollout-vorlauf-view-signatur.md) | additive Tabelle ohne Vorlauf; Alt-Tag-Lauf | Guard-Lauf 5 (`v0.1.2`, „mit Vorlauf" für die View-Signatur, Tabelle additiv) Exit 0/0/0 | konform |
| [`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md) | Tabelle im neutralen Modell, Grants in der Nacharbeit-Datei, kein Guard-Eintrag | §3.1, §3.3; `tools/schema/rolloutguard` unverändert | konform |
| [`ADR-0047`](../plan/adr/0047-rollenspezifische-dsn-verdrahtung.md) | rollenspezifische DSN-Verdrahtung | Annahme unter `cdc_admin`-, Run-Zustand und Schreiber unter `cdc_capture`-Login belegt (§3.5); die Rollen-/Aufrufer-Tabelle der ADR nennt `administration_request` nicht: Plan §6 (Zeile „Eingetreten und in der Fixrunde behoben") führt es als **an Architect gemeldet** („entscheidet der Architect — kein ADR- oder Spec-Auftrag dieses Slice", Ausgang bei Closure), nicht als geklärt | konform; Lücke offen gemeldet |
| [`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md) | `mapper/backfillrun.go` und `sqlexec.ReadBackfillRuns` im Gate; DB-Listen unverändert | §1 (100,0 % je Funktion, Listen-Diff leer, `db-package-lists-check` EXIT=0) | konform |
| [`ADR-0034`](../plan/adr/0034-ports-nach-faehigkeiten.md) | Ports nach Fähigkeit | drei Sentinel-Fehler an den Ports, keine Methoden-Signatur geändert (`git diff` der Ports gelesen) | konform |
| [`SPEC-029`](../../spec/pflichtenheft.md) | Feldform; Grants-Tabelle nennt nur `backfill_run` | Schema und Rechte gemessen; die Lücke der Grants-Tabelle (Antrags-Queue) steht im Plan §6 als an Architect gemeldet | konform; Lücke gemeldet |

## 6. Plan-vs-Code-Diff

Jede Zeile der Plan-Tabelle §3 ist im Diff vertreten: `schema.yaml` (Tabelle, Kopf-Absatz),
`nacharbeit-roles.sql`, Rollen-Test (4a) und (4b), die drei Adapter samt `backfill.go`,
`queries.go`, `mapper/backfillrun.go`, `sqlexec/translate.go` (jeweils mit Tests),
`plan.yaml`/`down.sql`, Guard-Skript (Lauf 5), die zwei Port-Dateien (drei Sentinel-Fehler),
`administration_roles_internal_test.go`, `backfillroles_test.go`, `roles_test.go`,
`backfillhelpers_test.go`, `harness/README.md`, `benutzerhandbuch.md`,
`harness/targets/schema-rollout.md`, der offene Plan `slice-backfill-sql-administration`. Die als
„geprüft, nicht geändert" geführte `postgresstorage/schema.sql` steht nicht im Diff.
Ungeplant im Diff (32 Dateien): nichts außer dem Plan selbst (`open/` → `in-progress/`) und dem
Review-Report; kein Fremdverzeichnis im Baum (`ls` der Wurzel, `git ls-files --others` leer).
Die Lifecycle-Moves sind reine Renames (`git show --stat -M`: `{open => next}/… | 0`,
`{next => in-progress}/… | 0`); `Verantwortlich` wurde in `next/` gesetzt (`38d4031d`,
1 Zeile) **vor** dem Move nach `in-progress/` (`03cd3804`).

## 7. Harte Regeln

- **[`AGENTS.md`](../../AGENTS.md) §3.3** — die zwei Moves sind reine Renames, die Inhaltsänderung „Verantwortlich"
  liegt in einem eigenen Commit dazwischen.
- **§3.5** — keine bestehende `Accepted`-ADR im Diff (`make doc-immutable` grün; der Diff berührt
  unter `docs/plan/adr/` nichts).
- **§3.6** — keine Schwelle gesenkt: `THRESHOLD`/`DB_COVERAGE_THRESHOLD`, `Dockerfile`, `harness/mk`
  im Diff leer; das Gate wird nicht gelockert.
- **§3.2** — `git diff afe2392f..HEAD -- '*.go' | grep -c nolint` = 0.
- **§3.7** — diff-skopierter Kandidatenlauf über die `+`-Zeilen aller geänderten `.go`, `.sql`, `.sh`
  und `schema.yaml` (Chronik-/Vorher-Vokabular: `früher`, `zuvor`, `nicht mehr`, `jetzt`, `Fixrunde`,
  `Review`, `Befund`, `F-<n>`, `slice-`, `welle-`, `wäre`/`würde`/`hätte`, `ersetzt`, `wurde`, `statt`,
  `bisher`, `vorher`, `nachher`): elf Treffer, alle Indikativ, Testdaten oder Fehlermeldungen
  („2 statt 1", `'zuvor'` als Fehltext-Literal, „wurde verändert" als Testmeldung, „ein früherer
  Antrag"); kein Verweis auf Review-Findings oder Slices in Kommentaren. Längste zusammenhängende
  Kommentarblöcke (gezählt): 65 Zeilen im Guard-Skript-Kopf und 44 in `nacharbeit-roles.sql`
  (beide Bestandsköpfe, der Diff ergänzt sie um wenige Zeilen), 50 in
  `roles_rollout_file_internal_test.go` (Bestand, um die Mutationsliste ergänzt), 20 in
  `administration_roles_internal_test.go`, 19 in `backfillhelpers_test.go`; der lange Vertrag steht in
  Plan, ADR und Sensor-Doku. Kein Befund.
- **§3.11** — `docs-check` 0 Befunde (`hostpaths`); der Report trägt keinen host-lokalen Pfad.
- **§3.12** — Zahlen im Plan gegen die Messung geprüft (§3.6): eine Abweichung („sieben"), sonst
  exakt; der Plan trägt keine Coverage-Zahl. Die Zahlen dieses Reports sind **gemessen** (Lauf im §1
  genannt) oder als **abgeleitet**/**übernommen** gekennzeichnet (V-1: 813 aus der Sensor-Doku).
- **§3.13** — Rechte statt Rollennamen gesucht (§3.6): Handbuch §2, §4, §5, `harness/README.md`,
  `harness/targets/schema-rollout.md` (vier Views), [`ADR-0047`](../plan/adr/0047-rollenspezifische-dsn-verdrahtung.md)-/[`SPEC-029`](../../spec/pflichtenheft.md)-Träger gemeldet, nicht still
  geändert. Der bewegte Träger „DB-Adapter-Nenner" fehlt im Feld (V-1).
- **Handbuch-Pflicht** — `git diff afe2392f..HEAD -- docs/user/benutzerhandbuch.md`: `Version: 1.46` →
  `1.47` **und** neue Zeile 1.47 in der Änderungshistorie; drei Stellen (§2, §4, §5) nachgezogen.
- **Commit-Traceability** — 11 Commits, jede Message nennt `LH-*` oder `ADR-*`, keine `SPEC-`/`ARC-`-Kennung im
  Betreff (§1).
- **[`AGENTS.md`](../../AGENTS.md) §3.1** — der Diff trägt keine Host-Toolchain; alle Läufe im gepinnten Image.
- **§3.10** — der Diff berührt keinen GitHub-Actions-Workflow.

## 8. Bewertung der Implementer-Neben-Aussagen

1. **`Admit` prüft die Antragsart nicht** — bestätigt: `UpdateAdministrationRequestApplied` trägt
   `WHERE administration_request_id = $1 AND status = 'pending'`; die Store-Tests nehmen einen Antrag der Art
   `enable` an. Der Träger steht im Plan von `slice-backfill-sql-administration` als DoD-Zeile „Bezug von
   Antrag und Run" und als §3-Zeile (`git diff` gelesen). Getragen.
2. **Annahme ohne Sperre und ohne Unique-Kante** — bestätigt (kein Lock, keine Unique-Kante im Schema);
   mit [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 1 Punkt 2 vereinbar, im Plan §6 benannt; nicht getestet, nicht von mir gemessen.
3. **Grenze des Login-Tests (Quelltabelle gehört `cdc_admin`)** — bestätigt: `nacharbeit-roles.sql`
   Zeilen 89–98 nennen die Eigentümer-Vorbedingung je aktivierter Tabelle; der Test setzt sie
   (`ALTER TABLE … OWNER TO cdc_admin`) und kommentiert sie. Dokumentierte Betriebs-Vorbedingung, keine
   verdeckte Test-Vereinfachung.
4. **Bestehende, nicht wiederholbare Rollen-Tests mit festen Kennungen** — bestätigt als Bestand:
   `internal/bootstrap/roles_wiring_test.go` steht nicht im Diff und trägt feste Login-Namen mit
   `DROP ROLE IF EXISTS`; die neuen Tests räumen ihre Logins und Zeilen ab (`t.Cleanup`) und laufen auf
   17 und 18 wiederholt (je zweimal in diesem Lauf, ohne Rückstand).

## 9. Befunde

Kein V-Befund blockiert. V-1 ist LOW, V-2…V-7 sind INFO ohne erwartete Aktion im Slice.

- **V-1 (LOW) — Sensor-Doku führt einen DB-Adapter-Nenner, den dieser Slice bewegt, ohne dass das
  [`AGENTS.md`](../../AGENTS.md) §3.13-Feld ihn suchte.** `harness/sensors/db-adapter-coverage.md` Zeilen 87–98 nennt „Der
  gemergte Nenner ist **813 Statements** (`postgresstorage` 472 · … ) — die **Zustandsgröße** dieses
  Gegenstands … hängt am **Code-Stand**", gebunden an den Lauf `slice-backfill-snapshot-reader`. Am
  Stand `HEAD` druckt der Lauf `830 von 1016` (je PostgreSQL 17 und 18; abgeleitet gegen die
  übernommenen 813: **+203** Statements im Nenner, **+179** gedeckt gegenüber 651). Der Satz bleibt
  als Beleg des benannten Laufs wahr, beschreibt aber nicht mehr den Ist-Stand, den die Überschrift
  „Zustandsgröße" behauptet. Ebenso `harness/sensors/coverage-gate.md` Zeilen 65–74 (Nenner 2082, gedruckt
  `83.2%`; hier `84.70%`/`84.80%`, den Unit-Nenner selbst habe ich nicht neu gezählt). Der Plan §6
  führt das Risiko „DB-Adapter-Coverage bewegt Zähler und Nenner" mit **offenem Ausgang**; das
  Suchlauf-Feld in §3 enthält den Träger nicht. Aktion: bei Closure den Träger nachziehen oder als
  lauf-gebunden lesbar machen und den Ausgang des §6-Risikos mit diesem Lauf tragen. Kein Gate betroffen.
- **V-2 (INFO) — Zahl im Suchlauf-Feld.** Plan §3, Zeile „Beschreibungen der Rolle `cdc_admin`": „Diff
  sieben"; gemessen **sechs** Dateien (`README.md`, `compose.yaml`, `benutzerhandbuch.md`,
  `harness/README.md`, `schema-rollout.md`, `pflichtenheft.md`). Die drei genannten Bewegungen stimmen; die
  Behandlung ist unberührt.
- **V-3 (INFO) — Fixrunde ohne eigenen Review-Report.** DoD „Review durchgeführt" ist durch den Report
  erfüllt; das neue Test-/Doku-Material der Fixrunde (Login-Test, `backfillroles_test.go`, Erweiterungen in
  `roles_test.go`, (4b), Guard-Lauf 5, Handbuch 1.47) habe ich selbst gelesen, gefahren und mutiert (§3.5,
  §4). Wer eine Reviewer-Durchsicht der Fixrunde verlangt, verlangt sie als eigenen Schritt.
- **V-4 (INFO) — „kein Superuser" der Login-Tests ist Konstruktion, keine Assertion.** Weder der
  Login-Test noch `backfillroles_test.go` prüft `rolsuper`; gemessen `rolsuper = f`, 0 Objekte im Eigentum
  (§3.2). Ein Superuser-Login würde die Tests nicht rot färben, aber die Konstruktion
  (`CREATE ROLE … LOGIN … IN ROLE`) erzeugt keinen.
- **V-5 (INFO) — `en_US.utf8` auf dem Alpine-Testcontainer ist nur ein Name.** Auf `postgres:17-alpine`
  und `postgres:18-alpine` existiert die Collation nicht (musl), `datcollate` druckt trotzdem
  `en_US.utf8`. Die Ordnung `0bf-…` vor Ziffern hält unter Default, `"C"`, `"und-x-icu"` und — auf dem
  Debian-Image `postgres:18` (glibc) gemessen — unter `"en_US.utf8"`. Zur Kenntnis für einen künftigen
  Betreiber-Hinweis (Plan §6 nennt ihn nur für den Abweichungsfall).
- **V-6 (INFO) — Lauf-Streuung und Nebenwirkung der Tier-Pfade.** Coverage-Gate 84.70 % (`make gates`)
  gegen 84.80 % (eigener `make coverage-gate`), gleicher Stand: lauf-gebunden. `make test-store`
  merget mit dem im `DB_COVERAGE_DIR` liegenden Replication-Profil und druckt deshalb vorab dieselbe
  Zahl; die Tier-Läufe überschreiben `tools/schema/plan.yaml` (Zielzeile), ich habe die Datei nach jedem
  Lauf per `git checkout` zurückgenommen. Eigenschaft der Tier-Pfade, nicht dieses Slice.
- **V-7 (INFO) — Der Slice-Umfang hängt am Bestand.** Der Defekt von Review F-1 (kein Grant auf
  `cdc.administration_request`) stammt aus dem Bestand seit [`ADR-0050`](../plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md); die Compose-Umgebung fährt alle
  DSNs als Superuser (`compose.yaml`), deshalb blieb er unentdeckt. Die Fixrunde schließt ihn; ob
  [`ADR-0047`](../plan/adr/0047-rollenspezifische-dsn-verdrahtung.md)/[`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md)/[`SPEC-029`](../../spec/pflichtenheft.md) eine Ergänzung brauchen, entscheidet der Architect (im Plan §6 als
  gemeldet geführt).

## 10. Ergebnis

| Prüfpunkt | Ergebnis |
|---|---|
| DoD §2 — `[x]`-Zeilen (10) | **10 von 10 erfüllt**, je mit eigenem Beleg (Zeile 8 mit V-3, Zeile 9 mit V-1/V-2) |
| DoD §2 — übrige `[ ]`-Zeilen (5) | **5 von 5 korrekt offen** (Closure-/Rollen-Sequenz; Reconciliation „entfällt") |
| Plan-vs-Code-Diff | **deckungsgleich**: jede Zeile der Plan-Tabelle §3 im Diff, ungeplant nichts außer Plan und Review-Report |
| Rechte (Grant-Lücke, Review F-1) | **behoben und gemessen**: `cdc_admin` `SELECT`/`UPDATE` t, `INSERT`/`DELETE` f auf `administration_request`; `backfill_run` `cdc_admin` `SELECT`/`INSERT`, `cdc_capture` `SELECT`/`UPDATE`, `cdc_reader` nichts; Login-Beleg auf 17.11 und 18.6 |
| Alt-Tag-/Guard-Lauf | **EXIT=0**, Lauf 5 `v0.1.2` Exit 0/0/0; unter der Grant-Mutation EXIT=1 an Lauf 5 |
| Schema/Tabelle | **bestätigt** ([`SPEC-029`](../../spec/pflichtenheft.md)-Spalten, CHECK bei Erstanlage, `estimated_rows` nullable, `plan.yaml`/`down.sql` gleich dem frischen Rollout bis auf die `target`-Zeile) |
| Adapter-Verträge | **bestätigt**: Annahme, `Finish` idempotent, `MarkRunning`, `Queued`-Ordnung, Writer-Blockvertrag, Commit auf beendeten Run, Zeitbegrenzung; elf Mutationen rot |
| Ordnung `0bf-…` | **bestätigt** unter Default, `"C"`, `"und-x-icu"` (17.11 und 18.6) und glibc `"en_US.utf8"` |
| PostgreSQL 17 und 18 | `make test-store` **EXIT=0** (17 und 18), `make test-replication` **EXIT=0** (17 und 18); `DB-Adapter-Coverage: 81.69% (gedeckt 830 von 1016 Statements; Profile gemergt: store,replication)` je Version |
| [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md), [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md), [`ADR-0114`](../plan/adr/0114-schema-rollout-vorlauf-view-signatur.md), [`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md), [`ADR-0047`](../plan/adr/0047-rollenspezifische-dsn-verdrahtung.md), [`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md), [`SPEC-029`](../../spec/pflichtenheft.md) | **konform**; die Lücken von [`ADR-0047`](../plan/adr/0047-rollenspezifische-dsn-verdrahtung.md) und der Grants-Tabelle von [`SPEC-029`](../../spec/pflichtenheft.md) sind im Plan als gemeldet geführt |
| Harte Regeln (§3.2, §3.3, §3.5, §3.6, §3.7, §3.9, §3.11, §3.12, §3.13, Handbuch) | **erfüllt**, §3.13 mit V-1 (LOW) |
| Commit-Traceability, MR-Immutabilität | **grün** (`RANGE=afe2392f..HEAD`) |
| Gates | **`make gates` EXIT=0**, `make test` EXIT=0, `make a-check` EXIT=0, `make coverage-gate` EXIT=0 (`84.80%`; im `make gates` `84.70%`) |

## Verdikt

**Bestätigt.** Die DoD-Behauptung des Implementers ist durch eigene Messung belegt: die Tabelle
`cdc.backfill_run` trägt die Feldform von [`SPEC-029`](../../spec/pflichtenheft.md) mit der Status-Menge als CHECK bei Erstanlage,
die Grants sind an Wegwerf-PostgreSQL 17.11 und 18.6 als Rechte-Matrix gemessen (auch der behobene
Bestandsdefekt: `cdc_admin` liest und vermerkt Anträge, legt keine an und löscht keine), die drei
Adapter laufen unter dem Login ihrer Rolle und scheitern unter der anderen mit SQLSTATE 42501, der
Alt-Tag-Lauf von `v0.1.2` über den Arbeitsbaum endet dreimal mit Exit 0 und bindet die Rechte, und
`make test-store` sowie `make test-replication` sind auf beiden Versionen grün. Elf Eingabeseiten-Mutationen
(Endzustand, Annahme, Commit, Zeitbegrenzung, Blockvertrag, `NULL`-Schätzung, aktiver Run, Übergang,
Ordnung, Grant) färben die Store-Tests rot, die Grant-Mutation zusätzlich den netzlosen Rollen-Test (4b),
den Login-Test und Lauf 5 des Guard-Skripts. Kein Befund blockiert; V-1 (LOW) verlangt bei Closure den
Nachzug der Nenner-Träger in den Sensor-Dokumenten, V-2…V-7 sind INFO.

**Übergabe:** Verifier → Planner. Offen bleibt die Closure: Closure-Notiz mit Lerneintrag, Ausgänge der
zehn §6-Risiken (darunter der DB-Adapter-Nenner mit diesem Lauf, V-1, und die an den Architect gemeldeten
Lücken von [`ADR-0047`](../plan/adr/0047-rollenspezifische-dsn-verdrahtung.md) und [`SPEC-029`](../../spec/pflichtenheft.md)), Beobachtungs-Register, Reconciliation-Zeile, Welle-Paarungen. Danach
der reine `git mv` nach `done/` ([`AGENTS.md`](../../AGENTS.md) §3.3, Fall 2).
