# Review-Report: slice-backfill-run-store — 2026-09-24

**Review-Art:** Code — der Diff führt die Tabelle `cdc.backfill_run` (Schema, Grants,
Rollen-Test der Rollout-Datei) und die drei Postgres-Adapter des Backfills (Annahme,
Run-Zustand, atomarer Schreiber) samt Store-Tier-Tests ein; geprüft gegen Plan, ADRs,
Pflichtenheft und `AGENTS.md` Hard Rules (Modul 10 §Drei Review-Arten). Kein
DoD-Abgleich — das ist Verifier-Aufgabe (Modul 11).

**Gegenstand:** Slice `slice-backfill-run-store`, Diff-Range `afe2392f..HEAD`
(`54fc81d3`). Slice-Commits `86347a9c`, `38d4031d`, `03cd3804` (Lifecycle,
Verantwortlich), `6c7dbbb4` (Tabelle, Grants, Rollen-Test der Rollout-Datei),
`92603002` (Adapter), `d095e5be` (CHECK-Test, Ordnungs-Soll), `9bfd8246`, `54fc81d3`
(Plan).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09" (seither um
weitere HIGH-Klassen ergänzt, u. a. Zahl-im-Träger, Beleg-Satz,
Zusage-ohne-Eingabeseite, Kommentar-Chronik, Handbuch-Zug).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-24.

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-NNN` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** (`v6.9.0` ·
> `regelwerk/<datei>.md` §<Abschnitt>). Ein `pfad`-Feld auf den **geprüften
> Gegenstand** zitiert den Stand des Laufs und darf ihn festhalten.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-backfill-run-store` (§1 Ziel, §2 DoD als Prüfmaßstab für
  Plan-Zusagen, §3 Plan samt Suchlauf-Feld, §6 Risiken und gemeldete Lücken) und
  Welle `welle-backfill-bestand`; der Nachbar-Plan `slice-backfill-sql-administration`
  (Empfänger der Übergaben)
- [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 4 (Atomarität, Run-Zustand), 6 (Ordnung, `change_id`), 7
  (Retention), [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 1 (Rollenschnitt, Annahme in einer Transaktion,
  Grants), [`ADR-0114`](../plan/adr/0114-schema-rollout-vorlauf-view-signatur.md) (additive Tabelle rollt ohne Vorlauf, Alt-Tag-Lauf), [`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md),
  [`ADR-0047`](../plan/adr/0047-rollenspezifische-dsn-verdrahtung.md) (rollenspezifische DSN-Verdrahtung), [`ADR-0050`](../plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md) (Antrags-Queue), [`ADR-0034`](../plan/adr/0034-ports-nach-faehigkeiten.md),
  [`ADR-0023`](../plan/adr/0023-fehlerklassifikation.md), [`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
- [`SPEC-029`](../../spec/pflichtenheft.md) (Feldform, Grants), [`SPEC-019`](../../spec/pflichtenheft.md), [`SPEC-001`](../../spec/pflichtenheft.md), [`SPEC-008`](../../spec/pflichtenheft.md); [`LH-FA-CAP-009`](../../spec/lastenheft.md),
  [`LH-FA-CAP-004`](../../spec/lastenheft.md), [`LH-FA-REA-004`](../../spec/lastenheft.md), [`LH-QA-SEC-001`](../../spec/lastenheft.md), [`LH-QA-SEC-002`](../../spec/lastenheft.md), [`LH-QA-SEC-003`](../../spec/lastenheft.md)
- `AGENTS.md` (Hard Rules §3.1, §3.3, §3.6, §3.7, §3.11, §3.12, §3.13),
  `harness/conventions.md` (MR-000/MR-001)
- Report-Gerüst: `docs/reviews/review-report.template.md`, Formvorbild
  `docs/reviews/review-slice-backfill-snapshot-reader.md`

**Eigenständig durchgeführte Prüfungen** (gemessen, nicht aus dem
Implementer-Bericht übernommen; Exit-Codes ungepiped in Log-Dateien gesichert):

- **Gates am Stand `54fc81d3`:** `make test` Exit 0; `make a-check` Exit 0;
  `make coverage-gate` Exit 0, gedruckt „coverage-gate: OK — Coverage 84.80% erfüllt
  Schwelle 80%"; `make test-store` (PostgreSQL 18) Exit 0, gedruckt
  „DB-Adapter-Coverage: 81.59% (gedeckt 829 von 1016 Statements; Profile gemergt:
  store,replication)"; `make gates` Exit 0 (vor dem Report-Commit), darin gedruckt
  „coverage-gate: OK — Coverage 84.90% erfüllt Schwelle 80%" — dieselbe Messung
  streut zwischen zwei Läufen um 0,1 Punkte (siehe F-8);
  `db-package-lists-check` „nennen dieselben 4 Pakete" (die drei namentlichen
  DB-Listen unverändert).
- **Rollout-Belege:** `bash tools/harness/run-schema-rollout-guard-test.sh`
  ungekürzt, `EXIT=0` (sechs Läufe; Lauf 1 validiert „Tables: 11 found"; Lauf 5
  gedruckt „Tag v0.1.2: Exit 0 (Rollout des Tags), Exit 0 (Arbeitsbaum, mit Vorlauf),
  Exit 0 (Arbeitsbaum, zweiter Lauf); Zeile alttag-ch über cdc.changes lesbar" — der
  „Vorlauf" dort ist der View-Signatur-Vorlauf aus [`ADR-0114`](../plan/adr/0114-schema-rollout-vorlauf-view-signatur.md), die neue Tabelle braucht
  keinen). Container, Netz und Temp-Verzeichnis danach aufgeräumt, `git status`
  sauber, `tools/schema/plan.yaml`/`down.sql` nicht im Arbeitsbaum-Diff. Das
  `CREATE TABLE "backfill_run"` in `plan.yaml` trägt Spalten, PK, FK und
  `chk_backfill_run_status` wie [`SPEC-029`](../../spec/pflichtenheft.md); die `target`-Zeile ist der committete
  Compose-Wert (`cdc-test-postgres`).
- **Recht von `cdc_admin` auf `cdc.administration_request` (F-1), Wegwerf-PostgreSQL 18,
  Rollout des Arbeitsbaums per `git archive` + `make -C … schema-rollout`, Exit 0:**
  `has_table_privilege` für `cdc_admin` auf `cdc.administration_request` —
  `SELECT`/`INSERT`/`UPDATE`/`DELETE` alle `f`; `relacl` der Tabelle leer (nur
  Eigentümer `postgres`); `cdc_capture` und `cdc_reader` ebenfalls `f`. Auf
  `cdc.backfill_run` `cdc_admin` `SELECT`/`INSERT` `t`, `cdc_capture`
  `SELECT`/`UPDATE` `t`, alle übrigen `f`, `relacl` =
  `{postgres=arwdDxtm/postgres,cdc_admin=ar/postgres,cdc_capture=rw/postgres}`.
  Ein Login `IN ROLE cdc_admin` (so legt es `docs/user/benutzerhandbuch.md` §2 an):
  die SQL von `ListPending` (`queries.go:309-311`) und `UpdateAdministrationRequestApplied`
  (`queries.go:338`) enden je mit „ERROR: permission denied for table
  administration_request"; `SELECT cdc.enable_table(...)` läuft dagegen bis zum
  `INSERT` (SECURITY DEFINER, EXECUTE-Grant). Nach `GRANT SELECT, UPDATE ON
  cdc.administration_request TO cdc_admin` laufen beide Anweisungen (0 Zeilen,
  kein Fehler).
- **Kollation (F-6):** `datcollate = en_US.utf8`, `datlocprovider = c`; dieselbe
  Textsortierung von `1, 748, 9999999999, 0bf-run-00000001, 0bf-a-00000002, a`
  unter dem Datenbank-Default und unter `COLLATE "C"` ergibt `0bf-a-…, 0bf-run-…, 1,
  748, 9999999999, a`; unter `COLLATE "und-x-icu"` (ICU) `0, 0bf-a-…, 0bf-run-…, 1,
  748, 9999999999, a`.
- **Suchlauf-Feld (Plan §3) an beiden Ständen nachgefahren** (`03cd3804`,
  `d095e5be`, Befehle wörtlich aus dem Feld): Tabellenlisten — beide Stände elf
  Dateien mit den genannten Zahlen (9/2/1/2/2/5/6/2/7/6/4); Rollenverteilung — acht
  Dateien, einzige bewegte Zahl `nacharbeit-roles.sql` 23 → 29; `backfill_run` in
  `internal/bootstrap` — Parent kein Treffer, Diff 5; Test-Bereinigung — Parent sieben,
  Diff elf Dateien (`backfillhelpers_test.go` 3, `backfillrun_test.go` 5,
  `backfillwriter_test.go` 4, `roles_test.go` 4), der zweite Befehl druckt am Diff-Stand
  null Treffer. Alle Zahlen stimmen. Offene Pläne mit `backfill_run`:
  `slice-backfill-sql-administration` und `welle-backfill-bestand` (siehe F-4).
- **Kommentare (§3.7):** hinzugefügte Zeilen aller geänderten `*.go`/`*.sql`/
  `*.yaml`/`*.md` (ohne Plan-Dateien, `plan.yaml`, `down.sql`) per Textsuche auf
  „seit slice"/`slice-`/`welle-`/„jetzt"/„vorher"/„nur noch"/„früher"/„bisher"/
  „wäre"/„würde"/„künftig"/„später"/„noch nicht": vier Treffer, alle inhaltlich
  zulässig („ein früherer Antrag", „ein späterer Fortschritt", „erst mit der letzten
  Anweisung", „noch nicht angelegt"). Konjunktive stehen nur in Test-Godoc
  („Rot färbende Mutation …", Subjekt der Test). Kein `//nolint`; keine
  host-lokalen absoluten Pfade im Diff.
- **Commit-Struktur:** `git show --stat -M` über `86347a9c` und `03cd3804` — reine
  Renames (`open → next`, `next → in-progress`, 0 Zeilen); `38d4031d` setzt
  `Verantwortlich` in `next/` **vor** dem zweiten Move — §3.3 erfüllt. Alle Betreffs
  tragen [`LH-FA-CAP-009`](../../spec/lastenheft.md)/`ADR-*`, keine `SPEC-`/`ARC-`-Kennung
  (`make commit-traceability` „5 Commit(s) … Betreffs ohne Struktur-ID"). `docs/user/`
  ohne Diff; der Diff führt keine Betreiber-Oberfläche ein (keine `CDC_*`-Variable, keine
  `cdc.*`-Funktion, kein Endpunkt).
- **Port-/Fehlerklassen:** `classifyError` (`usecase/backfill/service.go:313`) ordnet
  die drei neuen Sentinels `ErrBackfillRequestNotPending`, `ErrBackfillRunInvalid`,
  `ErrBackfillBlockInvalid` und `ErrInvalidBackfillTransition` der Klasse `internal`
  zu (Vertragsverletzung, kein Persistenzfehler); Treiberfehler laufen über
  `ErrBackfillStorage` (`storage`) — [`SPEC-008`](../../spec/pflichtenheft.md) gelesen, konsistent.

### Mutationen (Eingabeseite, selbst ausgeführt)

Datei nach jeder Mutation per `git checkout` zurückgenommen; Endstand `git diff` leer
(`tools/schema/plan.yaml` als bekannter Tier-Nebeneffekt jeweils zurückgenommen).
Netzlos: `go test -race` im gepinnten Race-Image mit `--network none`; Store-Tier:
`make test-store`; Rechte-Mutationen (R1–R3): ausgerolltes Wegwerf-Schema, Grant per
`psql` gesetzt, dann die Store-Tier-Rollen-Tests.

| # | Mutation | Ort | Ergebnis |
|---|---|---|---|
| M1 | Herkunftsprüfung `change.Origin != backfill` entfernt | `backfillwriter.go` (`blockRows`) | rot: `TestBackfillWriterRejectsBlocksOutsideTheContract` |
| M2 | Commit-Zähler-Prüfung `RowsCopied != written` entfernt | `backfillwriter.go` (`Commit`) | rot: `TestBackfillWriterCommitWritesTheRunRowLastAndChecksTheCounter` |
| M3 | Blocknummer-Prüfung `block.ID != expected` entfernt | `backfillwriter.go` (`blockRows`) | rot: `…RejectsBlocksOutsideTheContract`, `…RequiresConsecutiveBlockNumbers` |
| M4 | `boundedContext` aus `Finish` entfernt | `backfillrun.go` | rot: `TestBackfillRunOperationsEndAtTheirOwnTimeoutOnDetachedContext` |
| M7 | dasselbe in `InterruptRunning` | `backfillrun.go` | rot (dieselbe Test) |
| M8 | dasselbe in `Rollback` | `backfillwriter.go` | rot: `TestBackfillAdmissionAndWriterEndAtTheirOwnTimeoutOnDetachedContext` |
| M9 | dasselbe in `Commit` | `backfillwriter.go` | rot (dieselbe Test) |
| M10 | dasselbe in `Admit` | `backfilladmission.go` | rot (dieselbe Test) |
| M5 | `GRANT INSERT ON cdc.backfill_run TO cdc_capture` angehängt | `nacharbeit-roles.sql` | rot (netzlos): `TestRolloutDateiTraegtDieRechteDerVerdrahtung` |
| M6 | `UPDATE`-Grant von `cdc_capture` gestrichen | `nacharbeit-roles.sql` | rot (netzlos): dieselbe Test |
| S1 | `RowsAffected`-Prüfung des Antragsvermerks entfernt | `backfilladmission.go` | rot: `TestBackfillAdmitRollsBackWhenTheRequestIsNotPending` („Antrag besteht nicht: Fehler = <nil>") |
| S2 | `Finish`-`WHERE` trifft auch `completed` | `queries.go` | rot: `TestBackfillRunFinishOnAnEndedRunIsANoOpSuccess` |
| S3 | `Commit` des Schreibers rollt zurück statt zu committen | `backfillwriter.go` | rot: vier Tests (u. a. „nach dem Commit: 0 Transaktionen, 0 Changes") |
| S4 | `InterruptRunning` trifft auch `queued` | `queries.go` | rot: `…InterruptRunningTouchesOnlyRunningRunsOfTheSource` („meldet 3 Runs, erwartet 2") |
| S5 | `Queued` ohne Quellen-Filter (Eingabe `source` unbenutzt) | `queries.go` | rot: `TestBackfillRunQueuedIsOrderedByRequestedAtThenRunID` |
| R1 | `UPDATE` an `cdc_admin` auf `backfill_run` | Rechte im Schema | rot: `TestCdcAdminRoleInsertsBackfillRunButNeitherUpdatesNorDeletes` |
| R2 | `INSERT` an `cdc_capture` | Rechte im Schema | rot: `TestCdcCaptureRoleUpdatesBackfillRunButNeitherInsertsNorDeletes` |
| R3 | `SELECT` an `cdc_reader` | Rechte im Schema | rot: `TestCdcReaderRoleCannotReadTheBackfillRunBaseTable` |

Alle 18 Mutationen färben mindestens einen Test rot; die Zusagen des Diffs sind an ihrer
Eingabeseite gebunden — **für das, was die Tests prüfen**. Was sie **nicht** prüfen, ist
die Rolle, unter der die Adapter laufen (F-1, F-2).

---

## Findings

<!-- Kein Fließtext, kein Lösungsvorschlag im Befund. -->

### F-1 — `cdc_admin` trägt kein Recht auf `cdc.administration_request`: die Annahme ([`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 1) und der bestehende Administrationsweg scheitern unter einem Least-Privilege-Login

- `kategorie`: HIGH
- `quelle`: [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 1 (die Annahme läuft als `cdc_admin` und vermerkt
  den Antrag `applied`), [`ADR-0047`](../plan/adr/0047-rollenspezifische-dsn-verdrahtung.md) (Rollen je DSN), [`ADR-0050`](../plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md), [`LH-QA-SEC-001`](../../spec/lastenheft.md)…003;
  Skill-HIGH-Klasse „Korrektheitsfehler im kritischen Pfad" (Administrationsweg),
  Hard Rule §3.13 (Träger der bewegten Eigenschaft)
- `pfad`: `tools/schema/nacharbeit-roles.sql:143-144` (Grants des Diffs — die
  Tabelle `administration_request` kommt in der Datei an keiner Stelle vor),
  `internal/adapters/driven/postgresstorage/backfilladmission.go:109`
  (`UpdateAdministrationRequestApplied` in der Annahme-Transaktion),
  `internal/bootstrap/wiring.go:722` (`NewAdministrationRequest(ctx, cfg.AdminDSN, …)`),
  `compose.yaml:89-91` (alle drei DSNs als `postgres`),
  `docs/user/benutzerhandbuch.md:81` (`CREATE ROLE feed_admin_login LOGIN … IN ROLE
  cdc_admin`), Plan §6 (gemeldeter Befund)
- `befund`: Selbst gemessen (Wegwerf-PostgreSQL 18, Rollout des Arbeitsbaums): `cdc_admin`
  trägt weder `SELECT` noch `UPDATE` auf `cdc.administration_request`; ein Login `IN ROLE
  cdc_admin` scheitert mit „permission denied for table administration_request" an
  `ListPending` und `MarkApplied` — `Admit` (dieser Diff) braucht denselben `UPDATE`. Das ist
  **kein Defekt des Diffs, sondern des Bestands**: die Tabelle trägt seit der Antrags-Queue
  ([`ADR-0050`](../plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md)) keinen Grant an `cdc_admin`, `git log -S` findet in `nacharbeit-roles.sql`
  und in [`ADR-0047`](../plan/adr/0047-rollenspezifische-dsn-verdrahtung.md) nie einen; unentdeckt, weil `compose.yaml` alle drei DSNs als Superuser
  `postgres` fährt und jeder Store-Test mit dem Superuser läuft. Unter dem im Handbuch
  beschriebenen Betrieb bleibt jeder `cdc.enable_table`-/`exclude_column`-Antrag
  `pending` (Warnlog „Anträge lesen fehlgeschlagen"). Der Slice setzt die Annahme als
  `cdc_admin` **voraus** ([`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 1) und belegt sie nur als Superuser: der
  Rollen-Test deckt `cdc.backfill_run`, nicht die zweite Hälfte der Annahme-Transaktion;
  [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) (Accepted) führt den Befund nicht (er prüft die Rechte auf `backfill_run`,
  nicht auf `administration_request`) und nennt die bestehende Schleife „unverändert
  lauffähig". Der Implementer hat die Lücke gefunden und in Plan §6 gemeldet — sie ist
  offen, ohne Ausgang und ohne Träger im Empfänger-Plan (`slice-backfill-sql-administration`
  führt sie nicht).
- `verifizierbar`: ja — `SELECT has_table_privilege('cdc_admin','cdc.administration_request',
  'SELECT')` am ausgerollten Schema; ein Rollen-Test im Store-Tier, der die Anweisungen
  von `ListPending`/`MarkApplied`/`MarkFailed` und den Vermerk der Annahme als
  `cdc_admin`-Login ausführt (kein Test tut das heute)
- `klasse`: Rollen-Grant fehlt am Nachbar-Objekt eines Schreib-/Lesepfads — grün unter
  Superuser (`BEO-PGC/architect-verdikt-rollen-scope-luecke`-Familie; `slice-044`
  Retention-`DELETE`, Heartbeat, `table_schema` sind die früheren Fälle derselben Art)

**Prüfungsteile zur Frage (a)–(c) der Aufgabe:**

- (a) Heute berechtigt sind `ListPending`/`MarkApplied` **gar nicht** — weder über Grants
  noch über Eigentümer-Mitgliedschaft (`relowner = postgres`, `relacl` leer,
  `cdc_admin` ist `NOLOGIN` und nicht Mitglied einer Eigentümer-Rolle); der reale Betrieb
  trägt nur mit Superuser-/Eigentümer-Login. Der Compose-Aufbau
  (`compose.yaml:89-91`, `tools/schema/compose-init/01-cdc-schema.sql`) loggt sich für
  `CDC_ADMIN_DSN` als `postgres` ein, der Lauf `make test-integration` prüft die Rollen
  also nur an den Rollen-DSN-Belegen ([`LH-QA-SEC-001`](../../spec/lastenheft.md)…`003`), nicht am Administrationsweg.
- (b) Ein Defekt des Bestands, den dieser Slice sichtbar macht; der Slice **erbt** ihn
  zusätzlich in seiner Annahme (Annahme des Slice, ohne sie geprüft zu haben).
- (c) Schwere HIGH: der Administrationspfad ist der einzige Weg der SQL-Administration
  ([`LH-FA-ADM-001`](../../spec/lastenheft.md)) und der Auslösung des Backfills; Maßnahme und Träger benennt der
  Architect — Reviewer-Sicht auf den Umfang: ein Grant in `nacharbeit-roles.sql`
  (nach dem Muster der Lückenschließungen für `process_heartbeat` und `table_schema`,
  beide ohne eigene ADR gezogen), die Prüfung (4b) in
  `roles_rollout_file_internal_test.go`, ein Store-Tier-Rollen-Test und ein
  Login-Test nach dem Muster von `roles_wiring_test.go` für den Administrationsweg, kein
  Schema-Objekt (additiv, Alt-Tag-Lauf unberührt). Ob [`ADR-0047`](../plan/adr/0047-rollenspezifische-dsn-verdrahtung.md) oder [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) eine
  Ergänzung braucht — beide sind Accepted und tragen die Aussage nicht — ist die
  Architect-Entscheidung; die Rückführung zu diesem Slice (Fixrunde) oder zu einem
  eigenen Slice vor `slice-backfill-sql-administration` (der `Admit` verdrahtet) ist
  ebenfalls dort zu treffen. **Was der Slice hätte prüfen müssen:** den Rollen-Schnitt der
  **Annahme-Transaktion als Ganzes** (beide Tabellen) unter dem `cdc_admin`-Login, nicht
  nur die Rechte auf die neue Tabelle; und im Suchlauf-Feld (Rollenverteilung) nach den
  **Rechten** statt nach den Rollennamen fragen (`has_table_privilege` an den Tabellen, die
  die neuen Adapter berühren).

### F-2 — Kein Adapter-Test läuft unter der Rolle des Adapters; die Rollen-Zusage der drei Adapter ist an ihrer Eingabeseite (Rolle) nicht gebunden

- `kategorie`: MEDIUM
- `quelle`: [`ADR-0047`](../plan/adr/0047-rollenspezifische-dsn-verdrahtung.md) (Rollen-Verdrahtung, Fitness Function), [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 1,
  Skill-Klasse „Zusage ohne Bindung an ihre Eingabeseite", Plan §2 DoD
  (Annahme-Adapter, Adapter-Pflichten)
- `pfad`: `internal/adapters/driven/postgresstorage/backfillhelpers_test.go:44-63`
  (`newBackfillFixture`: `CDC_STORE_TEST_DSN` = Superuser), `backfilladmission_test.go:15-22`,
  `backfillrun_test.go`, `backfillwriter_test.go` (alle Adapter über `f.dsn`),
  `roles_test.go:419-525` (`SET ROLE` nur auf Statements an `cdc.backfill_run`)
- `befund`: `NewBackfillAdmission`, `NewBackfillRun` und `NewBackfillWriter` werden in
  allen Store-Tests mit der Superuser-DSN gebaut; die Rollen-Tests des Diffs führen rohe
  SQL-Statements per `SET ROLE` aus, nicht die Adapter. Ein Recht, das ein Adapter-Statement
  braucht und die Rolle nicht trägt (F-1 ist der reale Fall), färbt keinen Test — jede
  Mutation der Rechte (R1–R3) färbt nur die drei Statement-Tests. Für `cdc_capture`
  (Schreiber, Zustand) habe ich die Rechte auf `cdc.backfill_run` gemessen (`SELECT`,
  `UPDATE` `t`) und die Rechte auf `cdc.transaction`/`cdc.change` aus der Rollout-Datei
  gelesen (Bestand, `TestCdcCaptureRoleWritesChangesNotSchema`), aber nicht den
  Schreiber als Login ausgeführt.
- `verifizierbar`: ja — ein Store-Tier-Test, der `NewBackfillAdmission` mit der DSN eines
  `cdc_admin`-Logins und `NewBackfillRun`/`NewBackfillWriter` mit der eines
  `cdc_capture`-Logins baut (Muster `roles_wiring_test.go`, `newTestLoginRole`)
- `klasse`: Zusage ohne Bindung an ihre Eingabeseite (Rolle der Verbindung) — grün ohne Aussage

### F-3 — Kopplungs-Kommentar nennt `consumerstate_test.go` als Schema-Neuaufbau-Test; er ruft keinen auf, die Rollen-Tests des Diffs hängen an derselben Dateiordnung ohne Nennung

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.7 (Kopplung beschreibt, was da ist), Skill-Klasse „Beleg trägt
  seinen Satz nicht"
- `pfad`: `internal/adapters/driven/postgresstorage/backfillhelpers_test.go:23-27`
  („diese Dateien müssen vor `consumerstate_test.go`, `store_test.go` und
  `tableactivation_test.go` laufen — deren Tests räumen das Schema per `DROP SCHEMA cdc
  CASCADE`"), Plan §3 Zeile `schema.sql` (nennt `consumerstate_test.go` ebenfalls)
- `befund`: Gemessen (`grep -n "DROP SCHEMA" *_test.go`, `grep -ln "newTestStore(\|newTestActivation("`):
  die Ausführenden sind `store_test.go:50` und `tableactivation_test.go:39`;
  `consumerstate_test.go` enthält weder `DROP SCHEMA` noch einen Aufruf der beiden
  Helfer. Die neuen Rollen-Tests in `roles_test.go` laufen nur, weil `r` vor `s`/`t`
  sortiert — die Kopplung steht im Kopf von `backfillhelpers_test.go` und nicht an
  `roles_test.go`. Die Reihenfolge trägt heute (`make test-store` grün).
- `verifizierbar`: ja — `grep` wie oben
- `klasse`: Zitat nennt die falsche Stelle (Kopplungs-Kommentar)

### F-4 — Übergabe „Antragsart lesen" ist im Empfänger-Plan nicht getragen; `Admit` bindet Antrag und Run nur über die Kennung

- `kategorie`: LOW
- `quelle`: Modul 8 „Kein Pfeil ohne benennbares Artefakt", Plan §6 (gemeldet, Adresse
  `slice-backfill-sql-administration`), `AGENTS.md` §3.13
- `pfad`: `internal/adapters/driven/postgresstorage/queries/queries.go:338-341`
  (`UpdateAdministrationRequestApplied`, `WHERE administration_request_id = $1 AND status
  = 'pending'`), `backfilladmission.go:109`;
  `docs/plan/planning/open/slice-backfill-sql-administration.md:185-186` (Plan-Zeilen)
- `befund`: `Admit` nimmt jeden `pending`-Antrag an, unabhängig von `request_kind`, und
  prüft nicht, dass Quelle/Schema/Tabelle des Antrags die des Runs sind (die Kennung
  genügt). Der Implementer nennt den Art-Vergleich als Sache von
  `slice-backfill-sql-administration`; deren Plan (§3-Tabelle, §2 DoD) führt weder eine
  Zeile zu `queries.go`/`backfilladmission.go` noch ein Kriterium dafür — der Zeiger ist
  ein Absatz in einem anderen Plan, kein Eintrag im Empfänger. Bis dahin hat `Admit` keinen
  Aufrufer außerhalb des Use Cases (`grep "Admit("`: `service.go:111`), die Abgrenzung
  ist deshalb vertretbar.
- `verifizierbar`: nein — Plan-Lesen (Empfänger-Plan am Start des Slice)
- `klasse`: Übergabe ohne Träger im Empfänger-Plan

### F-5 — Annahme ohne Sperre und ohne Unique-Kante: die Ein-Instanz-Annahme ist benannt, aber nirgends erzwungen und ungetestet

- `kategorie`: INFO
- `quelle`: [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 1 Punkt 2 („Annahme, zu belegen im Store-Tier-Test der
  Antragsverarbeitung"), [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 4
- `pfad`: `backfilladmission.go:61-66` (Godoc), `tools/schema/schema.yaml:269-294`
  (`backfill_run`: PK auf `run_id`, kein Index auf `(source_id, schema_name, table_name)`
  für aktive Runs), Plan §6 (Risiko „Die Annahme-Transaktion …")
- `befund`: Zwei gleichzeitig annehmende Verbindungen sehen unter `READ COMMITTED` beide
  „kein aktiver Run" und legen zwei `queued`-Zeilen derselben Tabelle an; der Godoc, [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md)
  und Plan §6 nennen die Annahme (eine Goroutine, eine Instanz je Quelle) ausdrücklich, kein
  Test und keine Schema-Kante beleuchtet den Fall. Die Vereinbarkeit mit [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) ist
  gegeben (die ADR nimmt genau das an); der Risiko-Ausgang steht bei der Closure aus.
- `verifizierbar`: nein — Annahme, keine Prüfung am Diff
- `klasse`: Angenommene Ein-Instanz-Invariante ohne Erzwingung (benannt)

### F-6 — Ordnungs-Zusage `0bf-…` vor WAL-Kennungen: der Beleg deckt nur die Kollation der Test-Datenbank; nachgemessen trägt sie unter `C`, ICU und `en_US.utf8`

- `kategorie`: INFO
- `quelle`: [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 6, Plan §6 (erstes Risiko), [`LH-FA-CAP-004`](../../spec/lastenheft.md), [`LH-FA-REA-004`](../../spec/lastenheft.md)
- `pfad`: `backfillwriter_test.go:351-440`
  (`TestBackfillBlocksSortBeforeTheWALTransactionOnTheSamePosition`)
- `befund`: Der Test hält `datcollate` im Log fest und prüft `ReadChanges` sowie die View
  (mit dem `ORDER BY` **des Tests**, nicht mit einer Ordnung der View). Selbst gemessen ist
  die Ordnung unter dem Datenbank-Default (`en_US.utf8`, libc), `"C"` und `"und-x-icu"`
  identisch; die Zusage hängt an der Textsortierung, weil `0` unter jeder gängigen
  Kollation vor `1…9` steht — der Betreiber-Hinweis, den §6 bei Abweichung vorsieht, ist
  für diese drei nicht nötig. Das Risiko kann mit dieser Messung einen Ausgang tragen.
- `verifizierbar`: ja — `ORDER BY v COLLATE …` über die genannten Werte
- `klasse`: Beleg deckt eine Kollation von mehreren — Messung schließt die Lücke

### F-7 — Zeilenweise Einfügung eines Blocks: Leistungsgrenze steht im Plan, nicht am Code

- `kategorie`: INFO
- `quelle`: [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) Festlegung 3 (Richtgröße aus Messung), Plan §6 („Startwerte ohne
  Messung"), `AGENTS.md` §3.7 (Grenze als Kommentar-Klasse)
- `pfad`: `backfillwriter.go:96-125` (`AppendBlock`: eine `Exec`-Anweisung je Change, Zeile 116),
  `backfill.go:16-24` (Fristen 30 s/5 min)
- `befund`: Jeder Change ist ein eigener Round-Trip in der einen Schreibtransaktion; die
  Fristen sind als „Startwerte ohne Messung" gekennzeichnet (`backfill.go`), die
  Zeilen-für-Zeilen-Form ist nur in Plan §6 als Setzung geführt. Ein Leser von
  `AppendBlock` sieht die Grenze nicht.
- `verifizierbar`: nein — Messung gehört zu `slice-backfill-bench-richtgroesse`
- `klasse`: Setzung ohne Messung, Grenze nicht am Code

### F-8 — Coverage-Zahl streut zwischen zwei Läufen; der Implementer-Bericht nennt 84,90 %

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.12 Instanz A (Zahl gehört zu ihrem Lauf)
- `pfad`: keiner im Diff — die Zahl steht in keinem committeten Träger
- `befund`: Nachgemessen: `make coverage-gate` allein „Coverage 84.80%", innerhalb von
  `make gates` „84.90%", beide Exit 0, beide gegen dieselbe Schwelle 80 %; DB-Adapter-Coverage
  „81.59% (gedeckt 829 von 1016 …)" stimmt mit dem Bericht überein. Die Streuung ist die
  gemessene Eigenschaft der `coverage`-Stufe, kein Defekt des Diffs; wer die Zahl in
  einen Träger schreibt, nennt den Lauf.
- `verifizierbar`: ja — `make coverage-gate` wiederholen
- `klasse`: Zahl streut zwischen Läufen

### F-9 — Vorbestehend, nicht Gegenstand: `harness/targets/schema-rollout.md` „drei deklarierte Views" (Zeilen 13 und 71), vier im Leser; nicht wiederholbare Rollen-Tests

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.12, Plan §3 (Suchlauf Rollenverteilung), Plan §6
- `pfad`: `harness/targets/schema-rollout.md:13,71` (`views:` in `tools/schema/schema.yaml`
  trägt `active_tables`, `consumer_status`, `changes`, `retention_blockers`);
  `internal/adapters/driven/postgresstorage/roles_test.go:72-160`
  (`TestCdcCaptureRoleWritesChangesNotSchema`/`TestCdcAdminRoleManagesSchemaNotChanges`,
  feste Kennungen ohne Bereinigung)
- `befund`: Der Implementer nennt die Zeile 71; die Aussage steht auch in Zeile 13
  (`drei` mit den ersten drei Views). Beide Zeilen sind vom Diff unberührt
  (`git diff afe2392f..HEAD -- harness/targets` leer). Die zwei Rollen-Tests des Bestands schreiben unter festen
  Kennungen (`tx-capture-role`) und räumen nicht ab — ein zweiter Lauf gegen dieselbe
  Instanz träfe den Primärschlüssel; jeder `make test-store`-Lauf startet einen frischen
  Container. Beides ist Bestand.
- `verifizierbar`: ja — `grep -n "drei" harness/targets/schema-rollout.md`
- `klasse`: Zahl im Träger driftet gegen die Messung (Bestand) · Test-Isolation

---

## Negativbefunde

- geprüft, ohne Befund: `tools/schema/schema.yaml` (`backfill_run`) gegen [`SPEC-029`](../../spec/pflichtenheft.md) —
  Spalten, Typen, PK, FK auf `cdc.source`, `estimated_rows` nullable, beide Warn-Spalten
  `boolean NOT NULL DEFAULT false`, `rows_copied` `NOT NULL DEFAULT 0`,
  `chk_backfill_run_status` bei Erstanlage; ausgerollt gemessen (`plan.yaml`,
  `TestBackfillRunTableRejectsAStatusOutsideTheClosedSet` grün).
- geprüft, ohne Befund: `tools/schema/nacharbeit-roles.sql` (Grants) und
  `internal/bootstrap/roles_rollout_file_internal_test.go` (4a) — `cdc_admin`
  `SELECT`/`INSERT`, `cdc_capture` `SELECT`/`UPDATE`, `cdc_reader` kein Recht, niemand
  `DELETE` (am Schema gemessen: `relacl`); die Prüfung (4a) färbt bei M5/M6 rot.
- geprüft, ohne Befund: `tools/schema/plan.yaml`/`down.sql` — Ergebnis eines Rollouts
  (elf Tabellen, `backfill_run` im `CREATE TABLE` und im `DROP TABLE`), `target` ist der
  committete Compose-Wert; die Rollout-Läufe 1–6 der Wache Exit 0, additive Tabelle ohne
  Vorlauf, Alt-Tag `v0.1.2` Exit 0 zweimal.
- geprüft, ohne Befund: `internal/adapters/driven/postgresstorage/schema.sql` (zweiter
  Schema-Träger) — trägt laut Kopf nur die Store-Seite und führt `consumer`,
  `process_heartbeat`, `table_schema`, `administration_request` ebenfalls nicht; die neuen
  Store-Tests laufen gegen das ausgerollte Schema (`newBackfillFixture` bricht mit einer
  Fehlermeldung ab, wenn `cdc.backfill_run` fehlt), nicht gegen `schema.sql`.
  Die Dateiordnung als Träger der Reihenfolge ist F-3.
- geprüft, ohne Befund: `backfilladmission.go` — eine Transaktion (Lesen vor Einfügen,
  `INSERT`, Vermerk, Commit), `RowsAffected() != 1` bindet (S1), aktiver Run als
  `ErrBackfillRunActive`, Rollback über einen vom Abbruch gelösten Kontext, Frist, parametrisierte SQL,
  keine Bezeichner-Konkatenation.
- geprüft, ohne Befund: `backfillrun.go` — jede schreibende Anweisung trägt den
  erlaubten Ausgangszustand im `WHERE`; `Finish` auf einen beendeten Run ist ein
  wirkungsloser Erfolg (S2); `InterruptRunning` trifft nur `running` der Quelle (S4);
  `Queued` in der Ordnung `(requested_at, run_id)` und quellengefiltert (S5); Fristen (M4,
  M7); NULL-Schätzung liest als „unbekannt", schreibt als NULL (`TestBackfillAdmitStoresUnknownEstimateAsNullAndNotAsZero`).
- geprüft, ohne Befund: `backfillwriter.go` — Blockvertrag (Blocknummern lückenlos ab 1,
  Position `X`, Kennungen, `INSERT`, `origin = 'backfill'`, kein `old_data`, Zähler = Zahl
  der Changes; M1–M3), Commit mit `WHERE status = 'running'` (ein zwischenzeitlich
  `interrupted` gesetzter Run schreibt nichts, `TestBackfillWriterDoesNotCommitAnEndedRun`),
  Run-Zeile erst mit der letzten Anweisung, Rollback nach Fehler und Commit ohne Wirkung,
  Frist (M8, M9), Verbindungsfreigabe über `pgx` bei fehlgeschlagenem Commit; Vertrag des
  Ports und Kommentare in `backfilladmission.go`/`backfillwriter.go` (outbound) stimmen mit
  dem Verhalten.
- geprüft, ohne Befund: `queries/queries.go`, `mapper/backfillrun.go`,
  `sqlexec/translate.go` (`ReadBackfillRuns`) — explizite Spaltenlisten, `ON CONFLICT`-frei
  wo Doppelanlage ein Fehler sein soll, Domänen-Konstruktoren im Lesepfad, `mapper` und
  `sqlexec` liegen im Coverage-Gate (Unterpakete), drei DB-Listen unverändert.
- geprüft, ohne Befund: Fehlerklassen ([`SPEC-008`](../../spec/pflichtenheft.md)) — Treiberfehler `storage`,
  Vertragsverletzungen `internal` (siehe Prüfungen oben).
- geprüft, ohne Befund: `harness/README.md` (Sensor-Zeile `make test-store`) — der
  ergänzte Satz nennt keine Zahl; Handbuch (`docs/user/`) unberührt, keine
  Betreiber-Oberfläche im Diff (Aufschub mit Adresse `slice-backfill-sql-administration`
  laut Plan §2 und Suchlauf-Feld).
- geprüft, ohne Befund: Kommentare (§3.7), Suppression (§3.2: kein `//nolint`),
  host-lokale Pfade (§3.11), Commit-Struktur (§3.3, `commit-traceability`), Docker-only
  (§3.1: alle Läufe über `make`/gepinnte Images).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 1 |
| LOW | 2 |
| INFO | 5 |

**Finding-Klassen dieses Laufs:** Rollen-Grant fehlt am Nachbar-Objekt eines Pfads (grün
unter Superuser) · Zusage ohne Bindung an ihre Eingabeseite (Rolle der Verbindung) ·
Zitat nennt die falsche Stelle (Kopplungs-Kommentar) · Übergabe ohne Träger im
Empfänger-Plan · Angenommene Ein-Instanz-Invariante ohne Erzwingung · Beleg deckt eine
Kollation von mehreren · Setzung ohne Messung, Grenze nicht am Code · Zahl streut
zwischen Läufen · Zahl im Träger driftet (Bestand)

## Verdikt

**Merge-blockierend:** ja — F-1 (HIGH). Die Lücke liegt im Bestand und wird von diesem
Diff nicht verschlimmert, aber der Slice setzt mit `Admit` einen Pfad voraus, der unter
der Rolle, für die [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) ihn vorsieht, nicht läuft, und belegt ihn nur als Superuser;
die Rückführung (Fixrunde am Slice oder eigener Slice vor
`slice-backfill-sql-administration`) und die Frage einer ADR-Ergänzung entscheidet der
Architect (der Implementer hat sie selbst als Architect-Entscheidung gemeldet — kein
Rollen-Widerspruch). F-2 (MEDIUM) gehört zur selben Fixrunde. Alles Übrige — Schema,
Grants der neuen Tabelle, Blockvertrag, Zustandsübergänge, Fristen, Ordnung, Rollout — ist
am Diff gemessen und gebunden.

Die DoD-Zeile „Review durchgeführt" im Slice-Plan bleibt **offen**, weil eine Fixrunde
(oder ein Architect-Zug mit Rückführung) aussteht.

**Übergabe:** F-1 und F-2 gehen an den Architect (Träger und Umfang) und von dort an den
Implementer; F-3 bis F-9 gehen an den Implementer bzw. den Planner (F-4: Plan von
`slice-backfill-sql-administration` beim Start). Die **Finding-Klassen** gehen zusätzlich in
die Slice-Closure §7 und von dort in den Zähler. Dieser Report selbst ist ein
**Lauf-Beleg** (Audit: dieser Diff, dieser Skill, dieses Modell, dieses Verdikt) — er wird
über Läufe hinweg nicht wieder gelesen, und muss es nicht. Der Report ersetzt keine
Verifikation — DoD-/Spec-Konformität prüft der Verifier separat (Modul 11; anderes
Prüf-Artefakt, anderer Eingabe-Kontext).
