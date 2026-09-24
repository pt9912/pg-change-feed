# Slice backfill-run-store: Run-Store — Tabelle `cdc.backfill_run`, Grants, Postgres-Adapter für Run-Zustand und den einen atomaren Schreiber

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [welle-backfill-bestand](../welle-backfill-bestand.md).

**Bezug:** [`LH-FA-CAP-009`](../../../../spec/lastenheft.md) (Negative: Neubeginn ohne Verlust — ein
abgebrochener Run hinterlässt nichts), [`LH-FA-CAP-004`](../../../../spec/lastenheft.md) (Ordnung),
[`LH-FA-REA-004`](../../../../spec/lastenheft.md) (deterministische Sortierung), [`LH-FA-DAT-006`](../../../../spec/lastenheft.md)
(Metadaten-Erweiterbarkeit), [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 4/6/7 (Atomarität,
Run-Zustand, Ordnung, Retention), [`ADR-0043`](../../adr/0043-schemamigrationen-mit-d-migrate.md) (Schemamigrationen mit
d-migrate), [`ADR-0047`](../../adr/0047-rollenspezifische-dsn-verdrahtung.md) (rollenspezifische DSN-Verdrahtung), [`ADR-0017`](../../adr/0017-generische-change-tabelle.md)
(generische Change-Tabelle), [`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 1/3 (Rollenschnitt, Grants, Annahme in einer
Transaktion, Warn-Spalten, `estimated_rows` = `NULL` als „unbekannt").

**Berührte Spec-Stellen:** [`SPEC-029`](../../../../spec/pflichtenheft.md) (Feldform `cdc.backfill_run`, durch
`spec-nachzug`), [`SPEC-001`](../../../../spec/pflichtenheft.md), [`SPEC-002`](../../../../spec/pflichtenheft.md) — gelesen, nicht geändert;
[`ARC-006`](../../../../spec/architecture.md) (Driven Adapter).

**Verantwortlich:** Implementer-Agent, 2026-09-24.

**Autor:** Planner-Agent, Welle-Eröffnung [welle-backfill-bestand](../welle-backfill-bestand.md). **Datum:** 2026-09-23.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die Postgres-Seite des Runs: die Tabelle `cdc.backfill_run`
(`run_id` = die `administration_request_id` des Antrags, `source_id`,
`schema_name`, `table_name`, `status`, `requested_at`, `started_at`,
`finished_at`, `snapshot_position`, `rows_copied`, `estimated_rows`,
`error_message`; `estimated_rows` ist **nullable** — `NULL` heißt „unbekannt",
nie `0` —; dazu die **zwei Warn-Spalten** `warn_estimated_size` und `warn_duration`
für die beiden Warnungen (Typ `boolean NOT NULL DEFAULT false` wie in [`SPEC-029`](../../../../spec/pflichtenheft.md)
durch `spec-nachzug` festgelegt; `false`, solange keine Auswertung sie setzt); die geschlossene `status`-Menge als CHECK **bei
Erstanlage** der Tabelle), die Grants, und die Adapter für die drei Ports aus
`run-usecase`:

- der **Annahme-Adapter** (`Admit`, Paket `postgresstorage`, Pool der
  Administrations-Goroutine aus `CDC_ADMIN_DSN`, Naht `sqlexec.DB` mit `Begin`):
  **eine** Transaktion — Prüfung „kein aktiver Run (`queued`/`running`) derselben
  Tabelle" (Lesen vor Einfügen), `INSERT` der Run-Zeile `queued` (`run_id` = die
  `administration_request_id`, geschätzte Zeilenzahl, ggf. Warnung 1),
  `UPDATE … applied` des Antrags, der genau **eine** `pending`-Zeile treffen muss,
  Commit; jede Abweichung ist Rollback und Fehler, ein aktiver Run ein
  Sentinel-Fehler;
- der **Run-Zustands-Adapter** (Pool aus `CDC_CAPTURE_DSN`, **keine**
  Anlage-Operation): Übergänge, Fortschritt (`rows_copied` je Block
  **außerhalb** der Daten-Transaktion, damit er für Leser sichtbar ist; trägt
  auch die Warnung 2), der Abgleich `running` → `interrupted`, die `queued`-Zeilen
  der eigenen Quelle in Antragsreihenfolge (`ORDER BY requested_at, run_id`);
- der **atomare Schreiber**: **eine** Store-Transaktion über alle Blöcke — je
  Block eine `cdc.transaction`-Zeile (Kennung `0bf-<run-id>-<Block>`, Position
  `X`, `committed_at`) und `cdc.change`-Zeilen (`origin = 'backfill'`,
  `operation = 'INSERT'`, `new_data` = Bild) — und **ein** Commit am Ende
  zusammen mit der Run-Zeile `completed`; Rollback hinterlässt keine Zeile.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die SQL-Funktion, die Antragsart, der Worker, die View `cdc.backfill_status`
  und der Idempotenz-Guard** — `sql-administration`; diese Tabelle ist eine
  d-migrate-fähige Tabelle des neutralen Modells und braucht **keinen**
  Guard-Eintrag (nur Objekte außerhalb des Modells stehen in
  `knownForeignObjects`). *Zu belegen durch:* zweiter `make schema-rollout`.
- **Das Lesen des Bestands** — `cdc.changes` und `GET /changes` liefern ihn
  bereits über die Spalte `origin` aus `change-origin`.
- **Block-Transaktionen und Checkpoint** — Welle §6; die Ein-Transaktions-Form ist
  die Festlegung von [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 4.
- **Die Retention-Regel für Backfill-Changes** — sie gilt unverändert
  ([`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 7, kein Sonderpfad); der `e2e`-Slice trägt keinen
  Retention-Beleg, weil die Regel nichts Neues festlegt.

## 2. Definition of Done

- [x] Schema und Grants: `cdc.backfill_run` liegt in `tools/schema/schema.yaml`
      (CHECK auf `status` bei Erstanlage, `estimated_rows` nullable, die
      zwei Warn-Spalten nach [`SPEC-029`](../../../../spec/pflichtenheft.md)), die Grants stehen in
      `tools/schema/nacharbeit-roles.sql` — `cdc_admin` mit `SELECT`, `INSERT`,
      `cdc_capture` mit `SELECT`, `UPDATE`, niemand mit `DELETE`, `cdc_reader` ohne
      Recht auf die Basistabelle ([`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 1, Punkt 3). *Zu belegen
      durch:* `make schema-rollout` zweimal
      hintereinander mit Exit 0, der **Alt-Tag-Lauf** von
      `tools/harness/run-schema-rollout-guard-test.sh`
      ([`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md)
      Entscheidung 7: das Schema des jüngsten `v*`-Tags per `git archive`
      ausrollen, danach den Arbeitsbaum — Exit 0 zweimal, der zuvor eingefügte
      Datenstand über `cdc.changes` lesbar; der Bericht nennt den Tag und die
      gedruckten Exit-Codes; der Lauf entsteht in der Fixrunde von
      `slice-backfill-change-origin`), `tools/schema/plan.yaml` und
      `tools/schema/down.sql` neu erzeugt und mitcommittet, der Rollen-Test
      `internal/bootstrap/roles_rollout_file_internal_test.go` (netzlos, Teil der
      Gates) angepasst, und ein Rollen-Test im Store-Tier (`roles_test.go`-Muster,
      je Rolle ein Login): `cdc_capture` — `INSERT` scheitert mit SQLSTATE `42501`,
      `UPDATE` gelingt; `cdc_admin` — `INSERT` gelingt, `UPDATE` und `DELETE`
      scheitern; `cdc_reader` — `SELECT` auf die Basistabelle scheitert.
- [x] Annahme-Adapter: ein Store-Test gegen reale PostgreSQL zeigt, dass eine
      Annahme Run-Zeile `queued` und Antragsvermerk `applied` **zugleich**
      hinterlässt, dass ein Antragsvermerk ohne `pending`-Zeile (Rollback-Fall)
      **weder** Run-Zeile **noch** Vermerk hinterlässt, dass ein zweiter Antrag bei
      aktivem Run ohne zweite Zeile endet, und dass `estimated_rows` als `NULL` und
      nicht als `0` gespeichert und gelesen wird. *Zu belegen durch:*
      `make test-store`.
- [x] Atomarität: ein Store-Test gegen reale PostgreSQL zeigt, dass ein zweiter
      Leser vor dem Commit **keine** Zeile des Runs sieht und danach alle
      Blöcke zugleich, dass ein Rollback **keine** `cdc.transaction`- und keine
      `cdc.change`-Zeile hinterlässt, und dass jeder Change `origin = 'backfill'`,
      `operation = 'INSERT'` und `old_data IS NULL` trägt. *Zu belegen durch:*
      `make test-store`.
- [x] Ordnung und Zustand: mit einem WAL-Commit **auf derselben Position** `X`
      liest `cdc.changes` die Backfill-Blöcke **vor** dem WAL-Commit
      (`(commit_position, transaction_id, sequence)`); die Übergänge, der
      Fortschritt außerhalb der Daten-Transaktion und der Abgleich `running` →
      `interrupted` verhalten sich wie im Use-Case-Vertrag. *Zu belegen durch:*
      `make test-store` in der Test-Datenbank; der Bericht nennt die
      Kollation der Datenbank (`datcollate`), weil die Ordnung `0bf-…` vor
      Ziffern-Kennungen eine Textsortierung ist.
- [x] Adapter-Pflichten aus den Port-Verträgen (`BackfillRunPort`,
      `BackfillTransaction`): (1) **Zeitbegrenzung** — `Finish`,
      `InterruptRunning` und `Rollback` sind adapterseitig zeitbegrenzt und
      beenden sich bei einem vom Abbruch gelösten Kontext (ohne Frist des
      Aufrufers) nicht vorzeitig; Vorbild ist `closeTimeout` des
      Snapshot-Adapters (`internal/adapters/driven/postgressnapshot`);
      (2) **Endzustand idempotent** — `Finish` auf einen bereits beendeten Run
      ist ein wirkungsloser Erfolg: der Endzustand bleibt unverändert, der
      Aufruf meldet keinen Fehler (die Anweisung wirkt nur auf `queued`/`running`).
      *Zu belegen durch:* `make test-store` (`Finish(failed)` auf einen
      `completed`-Run liefert nil, die Zeile bleibt `completed`) und ein
      netzloser Test an der Ausführungs-Naht (eine blockierende Naht: jede der
      drei Operationen endet nach der eigenen Frist mit Fehler, auch bei
      abgelöstem Kontext); die Gegenseite am Use Case trägt
      `TestExecuteContextEndedInterrupts` in `usecase/backfill`.
- [x] Rollen unter echtem Login (Review F-1, F-2): der Administrationsweg
      (`processAdministrationRequests` mit allen vier Antragsarten) läuft mit
      einer `cdc_admin`- und einer `cdc_capture`-Login-Identität ohne Superuser
      bis zum Vermerk `applied`/`failed`; die drei Backfill-Adapter laufen unter
      dem Login ihrer Rolle (Annahme `cdc_admin`, Run-Zustand und Schreiber
      `cdc_capture`) und scheitern unter der Rolle des anderen mit SQLSTATE
      `42501`; `cdc_admin` trägt `SELECT`, `UPDATE` auf
      `cdc.administration_request`, sonst nichts auf ihr; der Rollout des
      Arbeitsbaums über den Stand des jüngsten `v*`-Tags setzt die Rechte. *Zu
      belegen durch:* `make test-store` (`internal/bootstrap` und
      `postgresstorage`), der netzlose Rollen-Test der Rollout-Datei (4b) und
      der Alt-Tag-Lauf von `tools/harness/run-schema-rollout-guard-test.sh`.
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein Self-Review (Modul 8).
      Nach der Fixrunde bleibt kein offenes HIGH oder MEDIUM.
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [x] Doku-Update: die Sensor-Zeile `make test-store` in `harness/README.md` und die Schema-Kommentare, soweit sie die Tabellenliste oder die Rollenverteilung aufzählen; das Benutzerhandbuch trägt den Rechteschnitt von `cdc_admin` (§2, §4, §5), die Betreiber-Oberfläche des Backfills bleibt bis `sql-administration` unberührt.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel ·
      neuer Sensor · benannte Spec-Lücke).
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in
      diesem Repo (Greenfield).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls
      eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen /
      weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      von der Closure der Welle [welle-backfill-bestand](../welle-backfill-bestand.md) (die Roadmap führt sie unter
      *Offene Wellen*, das Ereignis kann eintreten).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/schema/schema.yaml` | update | Tabelle `backfill_run`; CHECK auf `status` bei Erstanlage (spätere Erweiterung über Nacharbeit-SQL, nach dem Muster von `request_kind`); `estimated_rows` nullable; die zwei Warn-Spalten. |
| `tools/schema/nacharbeit-roles.sql` | update | Grants: `cdc_admin` `SELECT`, `INSERT`; `cdc_capture` `SELECT`, `UPDATE` auf `cdc.backfill_run`. Fixrunde (Review F-1): `cdc_admin` zusätzlich `SELECT`, `UPDATE` auf `cdc.administration_request` (kein `INSERT`, kein `DELETE`); der Kopf-Kommentar nennt den Ist-Rollenschnitt samt Antrags-Queue. Weitere Grants ergab die Messung des ganzen Administrationswegs nicht (§6). |
| `internal/bootstrap/roles_rollout_file_internal_test.go` | update | hält den Rollenschnitt gegen die Rollout-Datei; die neue Tabelle darf für `cdc_reader` **nicht** als Basistabellen-Grant erscheinen. |
| `internal/adapters/driven/postgresstorage/` (Arbeitsname `backfilladmission.go`, `backfillrun.go`, `backfillwriter.go`, + Tests) | neu | die drei Adapter über die schmale Ausführungs-Naht des Pakets (`sqlexec`); explizite Spaltenlisten; der Annahme-Adapter nutzt `Begin`. |
| `internal/adapters/driven/postgresstorage/backfillrun.go`, `backfillwriter.go` | Pflicht aus dem Port-Vertrag | Zeitgrenze für `Finish`, `InterruptRunning`, `Rollback` und `Finish` als wirkungsloser Erfolg auf einen beendeten Run (DoD-Punkt „Adapter-Pflichten"). |
| `internal/adapters/driven/postgresstorage/queries/queries.go` | update | die neuen Anweisungen. |
| `tools/schema/plan.yaml`, `tools/schema/down.sql` | regeneriert | Ergebnis von `make schema-rollout`, committet. |
| `tools/harness/run-schema-rollout-guard-test.sh` | ausführen; Fixrunde: update | der Alt-Tag-Lauf belegt die Tabelle gegen einen Alt-Bestand ([`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md) Entscheidung 7). Fixrunde (Review F-1): Lauf 5 fragt nach dem Upgrade des Schemas vom jüngsten `v*`-Tag die Rechte von `cdc_admin` ab (`has_table_privilege`: `administration_request` `SELECT`/`UPDATE` gesetzt, `INSERT`/`DELETE` nicht; `backfill_run` `SELECT`/`INSERT` gesetzt, `UPDATE` nicht) und prüft vor dem Upgrade, dass `UPDATE` auf `administration_request` fehlt. Entscheidung: im Guard-Skript statt an anderer Stelle, weil der Alt-Tag-Lauf dort schon das Ausrollen des Alt-Standes und des Arbeitsbaums trägt — der Zusatz sind sieben Abfragen ohne eigenen Aufbau. |
| `internal/adapters/driven/postgresstorage/schema.sql` | geprüft, nicht geändert | die eingebettete DDL trägt nur die Store-Seite (`source`, `source_table`, `schema_version`, `transaction`, `change`); `cdc.backfill_run` liegt allein im neutralen Modell, die Backfill-Store-Tests laufen vor den Tests, die das Schema per `DROP SCHEMA cdc CASCADE` samt hand-DDL neu aufbauen (`store_test.go`, `tableactivation_test.go`; die Dateinamen `backfill*_test.go` und `roles_test.go` sortieren vor ihnen). |
| `internal/application/port/outbound/backfilladmission.go`, `backfillwriter.go` | update (Plan-Nachzug) | drei Sentinel-Fehler an den Ports, damit die Adapter-Ablehnungen über `errors.Is` unterscheidbar sind: `ErrBackfillRequestNotPending` (die Annahme trifft keine `pending`-Zeile), `ErrBackfillRunInvalid` (Run passt nicht zu `Admit`/`Begin`), `ErrBackfillBlockInvalid` (Block oder Commit verletzt den Vertrag des Schreibers); die Port-Kommentare tragen den geschärften Vertrag (Blocknummern 1, 2, 3, …, Kennungen, Herkunft `backfill`, Zähler = angehängte Changes). Keine Änderung an Methoden-Signaturen. |
| `internal/adapters/driven/postgresstorage/backfill.go` | neu (Plan-Nachzug) | die zwei Frist-Startwerte und die Fehlerübersetzung (`ErrBackfillStorage`) der drei Adapter an einer Stelle. |
| `internal/adapters/driven/postgresstorage/mapper/backfillrun.go`, `sqlexec/translate.go` (+ Tests) | neu/update (Plan-Nachzug) | die Übersetzung der Run-Zeile (`BackfillRunRow`, `ToBackfillRun`, `ReadBackfillRuns`, NULL-Schätzung als „unbekannt“) liegt in den beiden Paketen, die im Coverage-Gate stehen und netzlos geprüft werden — die Logik im DB-Paket bliebe außerhalb des Gates ([`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)); keine neue Paketgrenze, die drei namentlichen DB-Listen bleiben unverändert. |
| `internal/adapters/driven/postgresstorage/backfill_internal_test.go`, `backfillhelpers_test.go`, `backfilladmission_test.go`, `backfillrun_test.go`, `backfillwriter_test.go`, `roles_test.go` | neu/update (Plan-Nachzug) | Store-Tier-Tests der DoD (Annahme, Zustand, Atomarität, Ordnung, Rollen) und der netzlose Frist-/Vertragstest an der Ausführungs-Naht (`backfill_internal_test.go`, läuft in `make test`). |
| `internal/bootstrap/roles_rollout_file_internal_test.go` | update | Prüfung (4a): Grants und Verbote auf `cdc.backfill_run` gegen den Rollout-Text; Fixrunde: Prüfung (4b) für `cdc.administration_request` (`cdc_admin` `SELECT`, `UPDATE`; `INSERT`/`DELETE` und jedes Recht der beiden anderen Rollen verboten). |
| `internal/bootstrap/administration_roles_internal_test.go` | neu (Fixrunde, Review F-1) | der ganze Antragsverarbeitungs-Pfad (`processAdministrationRequests`: `ListPending`, alle vier Antragsarten samt Use Case, `Registered`, `ExcludedColumns`, `MarkApplied`/`MarkFailed`, Schema-Speicher) unter einer `cdc_admin`- und einer `cdc_capture`-Login-Identität (`IN ROLE`, kein Superuser, kein Eigentum an `cdc`-Objekten; Quelltabelle gehört `cdc_admin`) — Test im Paket `internal/bootstrap`, weil es die Composition-Root-Funktion `processAdministrationRequests` aufruft und `tools/harness/run-store-tests.sh` das Paket vor den `postgresstorage`-Tests ausführt. |
| `internal/adapters/driven/postgresstorage/backfillroles_test.go`, `backfillwriter_test.go` (`newWriterRunAs`), `roles_test.go` | neu/update (Fixrunde, Review F-2, F-1) | die drei Backfill-Adapter unter dem Login ihrer Rolle (Annahme `cdc_admin`; Run-Zustand und Schreiber `cdc_capture`) und je unter der Rolle des anderen mit SQLSTATE `42501`; Rollen-Tests für `cdc.administration_request` je Rolle (`SET ROLE`-Muster wie bei `backfill_run`). |
| `internal/adapters/driven/postgresstorage/backfillhelpers_test.go` | update (Fixrunde, Review F-3) | der Kopplungs-Kommentar nennt die Dateien, die das Schema neu aufbauen (`store_test.go`, `tableactivation_test.go`), und die Dateinamen, die davor sortieren. |
| `harness/README.md` | update | die Sensor-Zeile `make test-store` nennt die Backfill-Tests im Store-Tier; Fixrunde: die Rollen-Tests, die Adapter unter Login und den Antragsverarbeitungs-Pfad. |
| `docs/user/benutzerhandbuch.md` | update (Fixrunde, Review F-1) | der Rechteschnitt von `cdc_admin` ist eine Betreiber-Oberfläche: §2 „Zugriff und Rollen" (Rollen-Tabelle), §4 „Schema aktualisieren" (Absatz „Rechte der drei Rollen"), §5 (`CDC_ADMIN_DSN`), `Version:` 1.47 und Zeile in der Änderungshistorie. |
| `harness/targets/schema-rollout.md` | update (Fixrunde, Review F-9) | die zwei Stellen „drei deklarierte Views" nennen vier (`active_tables`, `consumer_status`, `changes`, `retention_blockers`, am Schema gezählt: die vier Einträge unter `views:` in `tools/schema/schema.yaml`); die Beleg-Zeile zu Lauf 5 nennt die Rechte. |
| `docs/plan/planning/open/slice-backfill-sql-administration.md` | update (Fixrunde, Review F-4) | die Übergabe „Antragsart lesen / Bezug von Antrag und Run" ist dort als DoD-Punkt und §3-Zeile getragen; sonst unverändert. |

**Übergaben aus `slice-backfill-run-usecase`** (gemeldet, kein zusätzlicher Umfang; die
Ports liegen in `internal/application/port/outbound/backfill*.go`, der Run-Wert in
`internal/domain/model/backfillrun.go`):

- **Ganzer Run je Übergang.** `MarkRunning`, `RecordProgress` und `Finish` tragen den
  ganzen Run (`model.BackfillRun`), `Commit` den Run im Endzustand `completed`: der
  Run-Zustands-Adapter schreibt die Spalten je Übergang — `snapshot_position`,
  `rows_copied`, beide Warn-Spalten, `started_at`, `finished_at`, `error_message`.
- **Leere Zeitspalten und Fehlertext.** `started_at` ist leer, bis der Run `running` ist,
  auch bei `queued` → `failed` (die erneute Prüfung der Vorbedingungen endet einen Run,
  der nie `running` war); `error_message` trägt bei `failed` die Klasse vor dem Text
  (`<Klasse>: <Ursache>`), sonst ist er leer, `interrupted` trägt keinen Text.
- **Zulässige Übergänge.** `failed` folgt auf `queued` und `running`, `completed` und
  `interrupted` auf `running`; jeder Endzustand ist endgültig (Adapter-Pflicht (2) im
  DoD).
- **Kennungen.** Transaktions-Kennung `0bf-<run-id>-<Blocknummer, 8 Stellen>` aus
  `model.BackfillTransactionID` (Blocknummern zählen ab 1, ein Überlauf der achten
  Stelle ist ein Fehler) und `change_id` aus `model.ChangeIDFor`; beide stehen im
  übergebenen `model.ChangeTransaction`, der Schreiber bildet sie nicht neu.
- **Fehlerklassen.** Persistenzfehler der Adapter tragen `outbound.ErrBackfillStorage`;
  ein aktiver Run bei `Admit` ist `domainerrors.ErrBackfillRunActive`.

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „die Tabellenliste des `cdc`-Schemas und die Rollenverteilung der Grants"; beide Stände gemessen: Parent `03cd3804`, Diff `d095e5be` = letzter Code-Commit, danach nur Plan-Commits; die Zahlen sind die gedruckten Zeilen der Läufe):**

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| Tabellenlisten in Doku, Spec und Schema-Kommentaren | `git grep -c -e 'cdc.consumer_position' -e 'cdc.process_heartbeat' -e 'cdc.table_schema' -e 'cdc.administration_request' <Stand> -- docs/user harness spec tools/schema README.md AGENTS.md internal/adapters/driven/postgresstorage/schema.sql ':!tools/schema/plan.yaml' ':!tools/schema/down.sql'` | Parent und Diff drucken dieselben elf Dateien mit denselben Zahlen (`docs/user/benutzerhandbuch.md` 9, `docs/user/e2e-abdeckung.md` 2, `harness/README.md` 1, `harness/sensors/db-adapter-coverage.md` 2, `internal/adapters/driven/postgresstorage/schema.sql` 2, `spec/pflichtenheft.md` 5, `tools/schema/nacharbeit-administration.sql` 6, `nacharbeit-heartbeat.sql` 2, `nacharbeit-observability.sql` 7, `nacharbeit-roles.sql` 6, `tools/schema/schema.yaml` 4). Gelesen: die eine vollständige Tabellenliste ([`SPEC-001`](../../../../spec/pflichtenheft.md)) trägt `cdc.backfill_run` seit `spec-nachzug`; `db-adapter-coverage.md` Nr. 6 nennt die Tabellen, die `bootstrap.Run` liest — `cdc.backfill_run` gehört (noch) nicht dazu; der Kopf von `schema.sql` zählt Tabellen auf, die die DDL **nicht** trägt, und war schon ohne `table_schema` und `administration_request` unvollständig, eine Aufzählung der Modellliste ist er nicht; der Kopf von `schema.yaml` trägt eine Aufzählung der Tabellen des Modells | `schema.yaml`-Kopf um einen Absatz zu `backfill_run` ergänzt; die übrigen Träger unverändert, Begründung je Träger in der Spalte Befund |
| Rollenverteilung in Kommentaren und Doku | `git grep -c -e 'cdc_capture' -e 'cdc_admin' <Stand> -- docs/user harness spec README.md AGENTS.md tools/schema/nacharbeit-roles.sql` | Parent und Diff drucken dieselben acht Dateien; einzige Zahl, die sich bewegt: `tools/schema/nacharbeit-roles.sql` 23 → 29 (die eigenen Zeilen). Gelesen: der Kopf-Absatz „Rollenschnitt" von `nacharbeit-roles.sql` beschreibt `cdc_capture` und `cdc_admin` (nachgezogen); die Rollen-Tabelle in `docs/user/benutzerhandbuch.md` §Zugriff und Rollen nennt die Zwecke der Rollen — die Betreiber-Oberfläche des Backfills entsteht erst mit `slice-backfill-sql-administration`; `harness/targets/schema-rollout.md` Zeile „Schritt 1" nennt „die drei deklarierten Views", der Leser trägt vier (`retention_blockers`), der Stand ist von diesem Slice nicht bewegt | Rollen-Datei-Kopf nachgezogen; Handbuch unverändert (Aufschub mit Adresse `slice-backfill-sql-administration`, laut §2 dieses Plans); die „drei Views" in `schema-rollout.md` gemeldet, nicht geändert (nicht von diesem Slice bewegt) |
| Rollen-Test der Rollout-Datei | `git grep -c -e 'backfill_run' <Stand> -- internal/bootstrap` | Parent: kein Treffer; Diff: `internal/bootstrap/roles_rollout_file_internal_test.go` 5. Gelesen: Regel (7) („der Leser teilt kein Objekt mit einer schreibenden Rolle") trägt die neue Tabelle unverändert, weil `cdc_reader` keinen Grant auf sie hat; Regel (3) prüft nur `transaction`/`change` | Prüfung (4a) für `cdc.backfill_run` ergänzt (Grants und Verbote je Rolle) |
| Test-Bereinigung anderer Store-Tests | `git grep -c -e 'DELETE FROM cdc' -e 'DROP SCHEMA' <Stand> -- 'internal/adapters/driven/postgresstorage/*_test.go'` und `git grep -n -e 'DELETE FROM cdc.transaction"' -e 'DELETE FROM cdc.change"' -e 'DELETE FROM cdc.backfill_run"' -e 'DELETE FROM cdc.source"' d095e5be -- 'internal/adapters/driven/postgresstorage/*_test.go'` | Parent sieben Dateien mit Treffern, Diff elf (neu: `backfillhelpers_test.go` 3, `backfillrun_test.go` 5, `backfillwriter_test.go` 4, `roles_test.go` 4 — die vier Zahlen sind die eigenen Bereinigungen); der zweite Befehl druckt am Diff-Stand null Treffer: keine unskopierte Löschung auf `transaction`, `change`, `backfill_run` oder `source` | die neuen Tests bereinigen nur Zeilen unter ihren eigenen Kennungen (Präfix je Test); die Dateinamen `backfill*_test.go` sortieren vor den Tests mit `DROP SCHEMA cdc CASCADE` |

**§3.13-Suchlauf der Fixrunde (committetes Feld — bewegte Eigenschaft: „die Rechte, die `cdc_admin` auf die Antrags-Queue `cdc.administration_request` trägt"; gesucht nach Rechten und Tabellen, nicht nach Rollennamen; beide Stände gemessen: Parent `51c69243`, Diff `bd470a0b` = Code-Commit der Fixrunde; die Zahlen sind die gedruckten Zeilen der Läufe):**

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| Rechte-Vokabular in Doku, Spec, Schema, Tests | `git grep -c -e 'has_table_privilege' -e 'GRANT SELECT' -e 'GRANT INSERT' -e 'Least-Privilege' -e 'insufficient_privilege' -e 'SQLSTATE 42501' <Stand> -- docs/user README.md harness spec AGENTS.md compose.yaml examples tools/schema internal/bootstrap ':!tools/schema/plan.yaml' ':!tools/schema/down.sql'` | Parent dreizehn Dateien, Diff vierzehn; bewegt: `tools/schema/nacharbeit-roles.sql` 13 → 14, `internal/bootstrap/roles_rollout_file_internal_test.go` 5 → 6, neu `internal/bootstrap/administration_roles_internal_test.go` 1; unverändert u. a. `docs/user/benutzerhandbuch.md` 4, `compose.yaml` 1, `spec/pflichtenheft.md` 1, `spec/lastenheft.md` 2. Gelesen: der Kommentar in `compose.yaml` (die Umgebung verbindet als Superuser, die Rollentrennung belegen die Testcontainer-Tests) bleibt wahr; `spec/pflichtenheft.md` trägt die Grants-Tabelle von [`SPEC-029`](../../../../spec/pflichtenheft.md) nur für `cdc.backfill_run` | keine Änderung an `compose.yaml` und Spec; die Grants-Tabelle von `SPEC-029` an Architect gemeldet (§6) |
| Beschreibungen der Rolle `cdc_admin` | `git grep -c -e 'cdc_admin' <Stand> -- docs/user README.md harness spec compose.yaml examples` | Parent fünf Dateien (`README.md` 1, `compose.yaml` 1, `docs/user/benutzerhandbuch.md` 10, `harness/targets/schema-rollout.md` 2, `spec/pflichtenheft.md` 3), Diff sechs (dieselben fünf und `harness/README.md`): `docs/user/benutzerhandbuch.md` 10 → 13, `harness/README.md` 0 → 1, `harness/targets/schema-rollout.md` 2 → 3 (die eigenen Zeilen). Gelesen: die Rollen-Tabelle in §2 und die Zeile `CDC_ADMIN_DSN` in §5 des Handbuchs nannten den Zweck ohne Antrags-Queue (nachgezogen); die Zeilen mit „`cdc_admin`-Mitgliedschaft" als Voraussetzung in §4 (Aktivieren, Deaktivieren, Ausschließen, Consumer) sind mit dem Grant wahr; `README.md` nennt Least-Privilege ohne Rechte | Handbuch §2, §4 („Schema aktualisieren"), §5 nachgezogen, `harness/README.md` (Sensor-Zeile `make test-store`) nachgezogen |
| Rollenschnitt-Tabelle von [`ADR-0047`](../../adr/0047-rollenspezifische-dsn-verdrahtung.md) | `git grep -c -e 'cdc_admin' <Stand> -- docs/plan/adr/0047-rollenspezifische-dsn-verdrahtung.md` | Parent und Diff drucken 7; die Tabelle der Rollen und Aufrufer nennt `cdc.administration_request` an keiner Stelle | `Accepted`, nicht geändert; an Architect gemeldet (§6) |
| Träger der Antrags-Queue | `git grep -c -e 'administration_request' <Stand> -- docs/user harness spec docs/plan/adr README.md compose.yaml tools/schema/compose-init ':!tools/schema/plan.yaml' ':!tools/schema/down.sql'` | Parent neun Dateien, Diff zehn: `docs/user/benutzerhandbuch.md` 10 → 14, `harness/README.md` 0 → 1, `harness/targets/schema-rollout.md` 1 → 2; unverändert: `spec/pflichtenheft.md` 6, `docs/plan/adr/0050-…` 1, `0059-…` 4, `0065-…` 5, `0111-…` 3, `0112-…` 5, `harness/sensors/db-adapter-coverage.md` 1. Gelesen: `db-adapter-coverage.md` nennt `cdc.administration_request` als Lesezugriff des `bootstrap.Run`-Fixtures — betrifft den Schema-Stand des Tier-Laufs, nicht die Rolle | keine Änderung an ADRs, Spec und Sensor-Dokument |
| „drei Views" (F-9) | `git grep -n -e 'drei deklarierten Views' -e 'drei Views' <Stand> -- harness docs/user spec tools/schema` | Parent drei Treffer (`harness/targets/schema-rollout.md` Zeile 12 und 71, `tools/schema/schema.yaml` Zeile 39), Diff einer (`schema.yaml` Zeile 39). Gelesen: `schema.yaml` Zeile 39 spricht von den ersten drei Views und nennt die vierte, `retention_blockers`, im folgenden Absatz (Zeile 50) | die zwei Treffer in `schema-rollout.md` auf „vier" gezogen (gezählt am Schema: vier Einträge unter `views:`); `schema.yaml` unverändert |

**§3.13-Suchlauf der Closure (committetes Feld — bewegte Eigenschaft: „der Nenner der beiden Coverage-Messungen, den die neuen Pakete und Dateien des Slice verschieben"; nachgetragen aus Verifikation V-1; beide Stände gemessen: Parent `03cd3804`, Diff `bd470a0b` = letzter Code-Commit; die Zahlen sind die gedruckten Zeilen der Läufe):**

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| Nenner und gedeckte Zahl in den Sensor-Dokumenten | `git grep -c -w -e 2082 -e 813 -e 1731 -e 651 -e 2040 -e 472 <Stand> -- harness docs/user spec README.md AGENTS.md ':!tools/schema/plan.yaml'` | Parent und Diff drucken dieselben zwei Dateien mit denselben Zahlen: `harness/sensors/coverage-gate.md` 6, `harness/sensors/db-adapter-coverage.md` 5; keine Zahl in `harness/README.md`, `docs/user`, `spec`, `README.md`, `AGENTS.md`. Gelesen: `db-adapter-coverage.md` führte 813 (`postgresstorage` 472) und den Lauf des Vorgänger-Slice, `coverage-gate.md` 2082 und `83.2%` — je als Beleg des benannten Laufs wahr, als „Zustandsgröße" des Ist-Stands nicht mehr; `coverage-gate.md` Zeilen zu `slice-097` (1936) und `slice-085` (1903, `31/472`) sind datierte Messungen anderer Läufe | Zahlen im Lauf der Closure neu gemessen (Stand `c7045f81`, `make test-store` gefolgt von `make test-replication`, `make coverage-gate`, die Stufe nachgestellt) und mit Lauf und Commit-Kennung nachgezogen: DB-Adapter 830 von 1016 (`81.69%`; `postgresstorage` 675, gedeckt 527), Unit-Stufe 2033 von 2398 (gedruckt `84.8%`, Gate-Zeile `84.90%`); die datierten Absätze unverändert; am Stand `510f6ffe` (nach dem Nachzug) druckt derselbe Befehl nur noch `harness/sensors/coverage-gate.md` 2 (`31/472`, `(1903 + 472)`: datierte Messung des Laufs `slice-085`) |

## 4. Trigger

**Start** (`next` → `in-progress`): wenn `run-usecase` und `change-origin` in
`done/` liegen und kein anderer Slice in `in-progress/` liegt (WIP-Limit 1).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls Schema/Grants
  und die beiden Adapter nicht in einem Review tragen — der abtrennbare Teil
  ist Schema samt Grants (dann als eigener Schritt vor den Adaptern).
- `in-progress` → `open` (blockiert — Carveout?): falls d-migrate die Tabelle
  oder den CHECK bei Erstanlage nicht konvergiert (Exit 5) — dann Nacharbeit-SQL
  nach [`ADR-0043`](../../adr/0043-schemamigrationen-mit-d-migrate.md) §Re-Evaluierungs-Trigger und ein Eintrag in
  `knownForeignObjects` im selben Commit; Architect-Frage.

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + `make test-store` und
`make schema-rollout` (zweimal) real grün + der Alt-Tag-Lauf von
`tools/harness/run-schema-rollout-guard-test.sh` real grün + Closure-Notiz mit
Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Die Textsortierung trägt die Ordnung `0bf-…` vor Ziffern-Kennungen**
  ([`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 6 leitet sie her: WAL-Kennungen beginnen nie mit
  `0`). Sie hängt an der Kollation der Datenbank (`ORDER BY t.commit_position,
  c.transaction_id, c.sequence` in der View). *Erwartet, zu belegen durch:* der
  Ordnungs-Test in der Test-Datenbank samt genannter Kollation; weicht eine
  gängige Kollation ab, wird der Betreiber-Hinweis Teil des Handbuchs.
  *Messung (Implementer, PostgreSQL 18-alpine, `datcollate = en_US.utf8`):* `string_agg(v, ', ' ORDER BY v COLLATE …)` über
  `1, 748, 9999999999, 0bf-run-00000001, 0bf-a-00000002, a` liefert unter dem Datenbank-Default, `"C"` und `"und-x-icu"`
  dieselbe Reihenfolge (`0bf-a-00000002, 0bf-run-00000001, 1, 748, 9999999999, a`); der Reviewer maß dasselbe (Review F-6,
  Messung des Reviewers, hier nachgemessen).
  **Ausgang:** *entfallen* — die Messung trägt: die Reihenfolge `0bf-…` vor `1`, `748`,
  `9999999999` hält unter dem Datenbank-Default, `"C"` und `"und-x-icu"` an PostgreSQL 17.11 und
  18.6 und unter glibc `"en_US.utf8"` am Image `postgres:18` (Debian) (Verifikation
  `verifikation-slice-backfill-run-store` §3.4); auf den Alpine-Testcontainern ist `en_US.utf8` nur
  der gedruckte `datcollate`-Name, keine wählbare Kollation, der Beleg dort trägt Default, `"C"` und
  ICU (Verifikation V-5). Eine gemessene Kollation weicht nicht ab, ein Betreiber-Hinweis im Handbuch
  entfällt.
- **Der CHECK bei Erstanlage konvergiert** (die Zusage von [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage
  4): die Tabelle `change` trägt CHECKs bei Erstanlage bereits (`chk_change_operation`
  in `tools/schema/schema.yaml`); für **diese** Tabelle ungemessen. *Erwartet, zu
  belegen durch:* Erst- und Folgelauf des Rollouts und der Alt-Tag-Lauf (die
  Tabelle ist additiv gegenüber dem Schema des jüngsten `v*`-Tags; additive
  Änderungen rollen über einen Alt-Bestand, gemessen im Architect-Verdikt
  `architect-verdict-schema-rollout-view-signatur`, Szenario 5, für eine neue
  Tabelle mit Fremdschlüssel und Default — der CHECK und die Warn-Spalten
  gehören nicht zu dieser Messung). **Ausgang:** *entfallen* — der Rollout konvergiert: `plan.yaml`
  trägt `chk_backfill_run_status` und die zwei Warn-Spalten (`BOOLEAN NOT NULL DEFAULT FALSE`) im
  `CREATE TABLE`; Lauf 5 des Guard-Tests rollt das Schema von `v0.1.2` und danach den Arbeitsbaum
  aus, gedruckt `Exit 0 (Rollout des Tags), Exit 0 (Arbeitsbaum, mit Vorlauf), Exit 0
  (Arbeitsbaum, zweiter Lauf)` mit lesbarem Altbestand (Verifikation §1); kein Eintrag in
  `knownForeignObjects` (`git diff` über `tools/schema/rolloutguard` leer, Verifikation §3.3);
  `TestBackfillRunTableRejectsAStatusOutsideTheClosedSet`.
- **Der Fortschritt außerhalb der Daten-Transaktion** braucht eine zweite
  Verbindung und kollidiert nicht mit dem Commit derselben Run-Zeile. *Erwartet,
  zu belegen durch:* der Zustands-Test mit gleichzeitigem Fortschritts-Update
  und Commit. **Ausgang:** *entfallen* —
  `TestBackfillProgressDoesNotWaitForTheOpenWriteTransaction` schreibt den Fortschritt auf der
  zweiten Verbindung, während die Schreibtransaktion offen ist (Wartezeit unter 2 s, sichtbarer Zähler und
  Warn-Spalte, keine Daten vor dem Commit), und committet danach dieselbe Run-Zeile als
  `completed`; die Mutation „Commit rollt zurück" färbt ihn rot (Verifikation §2 Zeile 4, §4 M3).
- **Die Rollenlage** ist mit [`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 1 gesetzt (`cdc_admin` legt an, `cdc_capture`
  ändert, `cdc_reader` liest nur die View); ein `INSERT` für `cdc_capture`, ein
  `UPDATE` oder `DELETE` für `cdc_admin` oder ein `SELECT`-Grant für `cdc_reader`
  auf die Basistabelle wäre ein Rollen-Verstoß ([`LH-QA-SEC-001`](../../../../spec/lastenheft.md)…[`LH-QA-SEC-003`](../../../../spec/lastenheft.md)).
  *Erwartet, zu belegen durch:* der Rollen-Test im Store-Tier und der netzlose
  Rollen-Test der Rollout-Datei. **Ausgang:** *entfallen* — die Rechte auf `cdc.backfill_run` sind
  an PostgreSQL 17.11 und 18.6 gemessen (`cdc_admin` `SELECT`/`INSERT`, `cdc_capture`
  `SELECT`/`UPDATE`, `cdc_reader` nichts, niemand `DELETE`; Verifikation §3.2); die Rollen-Tests je
  Rolle im Store-Tier und die Prüfung (4a) der Rollout-Datei färben unter den Mutationen R1 bis R3,
  M5 und M6 rot (Review „Mutationen"). Der Rollenschnitt auf die Antrags-Queue ist das Risiko
  darunter.
- **Die Annahme-Transaktion ist nicht atomar oder greift zu weit.** Der
  Antragsvermerk muss genau eine `pending`-Zeile treffen, sonst Rollback; der
  Verzicht auf eine Sperre zwischen mehreren Annehmenden beruht auf der Annahme,
  dass `processAdministrationRequests` Anträge sequenziell in **einer** Goroutine
  verarbeitet und je Quelle eine Instanz Anträge annimmt ([`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 1,
  Punkt 2). Zwei gleichzeitig annehmende Verbindungen sähen unter `READ COMMITTED` beide „kein aktiver Run"; das ist mit
  [`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) vereinbar (die ADR nimmt genau das an),
  benannt, aber weder erzwungen (keine Sperre, keine Unique-Kante auf aktive Runs) noch getestet (Review F-5).
  *Erwartet, zu belegen durch:* der Rollback-Test des zweiten
  Liefer-Punkts. **Ausgang:** *weiter offen* — Register `BEO-PGC/ein-instanz-annahme-ohne-erzwingung`
  (1×). Die Atomarität trägt: `TestBackfillAdmitLeavesQueuedRunAndAppliedRequestTogether` und
  `TestBackfillAdmitRollsBackWhenTheRequestIsNotPending` (kein Antrag, `applied`, `failed`; Mutationen
  M2 und S1 rot, Verifikation §2 Zeile 2, §4). Offen bleibt die Annahme ohne Sperre: benannt, weder
  erzwungen noch getestet (Review F-5); Adresse ist der Re-Evaluierungs-Trigger von
  [`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) („Mehr als eine
  Instanz je Quelle nimmt Anträge an"), den sequenziellen Fall trägt der DoD „Verarbeitung" in
  `slice-backfill-sql-administration`.
- **Geteilter Zustand in den Store-Tests** (`BEO-PGC/test-isolation-geteilter-zustand`,
  offen, 1×): unskopierte Bereinigung in den bestehenden Tests würde die Zeilen
  eines parallelen Tests treffen. *Erwartet, zu belegen durch:* die neuen Tests
  bereinigen nur ihre Run-Kennung. **Ausgang:** *entfallen* — der Suchlauf (§3, vierte Zeile) druckt
  am Diff-Stand null unskopierte Löschungen auf `transaction`, `change`, `backfill_run` und
  `source`, die neuen Tests bereinigen unter dem Präfix ihrer eigenen Kennungen (Verifikation §3.6).
  Die Dateiordnungs-Kopplung der Tests mit `DROP SCHEMA cdc CASCADE` (Review F-3) ist ein anderer
  Zustand desselben Musters und im Register geführt (`BEO-PGC/test-isolation-geteilter-zustand`,
  2×).
- **DB-Adapter-Coverage** ([`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md) Punkt 3): neuer, DB-gestützt geprüfter
  Code in `postgresstorage` bewegt Zähler und Nenner. *Erwartet, zu belegen
  durch:* der Bericht nennt die Zahl mit ihrem Lauf ([`AGENTS.md`](../../../../AGENTS.md) §3.12 Instanz A).
  **Ausgang:** *eingetreten*, im Slice gelöst (kein Carveout, kein Folge-Slice): Zähler und Nenner
  der beiden Messungen bewegen sich mit dem Slice (Verifikation V-1); die Sensor-Dokumente führen die
  Zahlen mit Lauf und Commit-Kennung (Stand `c7045f81`, `make test-store` gefolgt von
  `make test-replication`): DB-Adapter-Coverage `81.69%`, gedeckt 830 von 1016 Statements
  (`postgresstorage` 675, gedeckt 527); Unit-Stufe gedeckt 2033 von 2398, gedruckt `84.8%`, Gate-Zeile
  `84.90%` (Läufe desselben Stands druckten `84.70%` bis `84.90%`). Beide Schwellen (70 % und 80 %) sind
  erfüllt, das Suchlauf-Feld der Closure (§3) trägt den Träger.
- **Eingetreten und in der Fixrunde behoben (Review F-1) — `cdc_admin` trug kein Recht auf `cdc.administration_request`.**
  *Sachverhalt:* die Tabelle trägt keinen Grant an `cdc_admin`; `ListPending`, `ExcludedColumns` und der Vermerk
  `applied`/`failed` (dieselbe Anweisung wie im Vermerk von `Admit`,
  [`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 1) scheitern unter einer
  Login-Identität `IN ROLE cdc_admin` mit „permission denied for table administration_request". *Messung:* Implementer am
  ausgerollten Schema (`has_table_privilege('cdc_admin', 'cdc.administration_request', 'SELECT')` und `'UPDATE'`: beide
  `false`); Reviewer an PostgreSQL 18 (`relacl` leer, Login-Lauf der Anweisungen, `docs/reviews/review-slice-backfill-run-store.md`
  F-1); Fixrunde: `TestAdministrationPathRunsUnderLeastPrivilegeLogins` (der ganze Weg von `processAdministrationRequests`
  unter einem `cdc_admin`- und einem `cdc_capture`-Login) ist ohne den Grant rot gesehen (Status bleibt `pending`, Log „Anträge
  lesen fehlgeschlagen") und mit `SELECT, UPDATE` grün — die Messung des ganzen Weges (`ListPending`, `enable`/`disable`/
  `exclude_column`/`include_column` samt Publication-DDL an einer `cdc_admin`-eigenen Quelltabelle, `Registered`,
  `ExcludedColumns`, `CurrentVersion` unter `cdc_capture`, Vermerk `applied` und `failed`) fand **kein** weiteres fehlendes
  Recht. *Herkunft:* Lücke im Bestand seit der Antrags-Queue
  ([`ADR-0050`](../../adr/0050-sql-administration-antragsqueue-und-live-reload.md)) — `nacharbeit-roles.sql` trug nie einen
  Grant auf die Tabelle; unentdeckt, weil `compose.yaml` alle drei DSNs als Superuser `postgres` fährt und jeder Store-Test
  mit dem Superuser lief. *Wirkung:* Anträge der SQL-Administration bleiben mit einem Least-Privilege-Login `pending`.
  *Maßnahme:* `GRANT SELECT, UPDATE` in `nacharbeit-roles.sql`, Prüfung (4b) im netzlosen Rollen-Test, Rollen-Tests im
  Store-Tier, der Login-Test über den ganzen Weg, die Adapter unter Login (F-2), Rechte-Abfragen im Alt-Tag-Lauf, Handbuch.
  **An Architect gemeldet:** [`ADR-0047`](../../adr/0047-rollenspezifische-dsn-verdrahtung.md) nennt den Rollenschnitt von
  `cdc_admin` ohne `administration_request`; ob [`ADR-0047`](../../adr/0047-rollenspezifische-dsn-verdrahtung.md) oder
  [`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) (beide `Accepted`) eine Ergänzung brauchen,
  und ob die Grants-Tabelle von [`SPEC-029`](../../../../spec/pflichtenheft.md) (nur `cdc.backfill_run`) den Vermerk auf
  `administration_request` nennt, entscheidet der Architect — kein ADR- oder Spec-Auftrag dieses Slice. **Ausgang:** *eingetreten*, im Slice
  behoben (Fixrunde `bd470a0b`; kein Carveout, kein Folge-Slice): `cdc_admin` trägt `SELECT`, `UPDATE` auf
  `cdc.administration_request`, gemessen an PostgreSQL 17.11 und 18.6 (Verifikation §3.2), und die Träger
  färben unter der Mutation „Grant streichen" rot (Verifikation §4, M11 und M11b). Die Frage an den Architect
  steht als Beleg im Register (`BEO-PGC/adr-aussage-breiter-als-ihre-messung`, 3×; Adresse Architect), der
  Befund als weiterer Beleg der Klasse `BEO-PGC/rollen-test-abdeckungsluecken` (3×).
- **Gemeldet (Implementer) — der Annahme-Adapter liest die Antragsart nicht.** Die geschlossene `request_kind`-Menge
  trägt `backfill` noch nicht (`tools/schema/nacharbeit-administration.sql`); die Store-Tests nehmen deshalb einen
  offenen Antrag der Art `enable` an. Der Art-Vergleich (`request_kind = 'backfill'` in der `WHERE`-Klausel des
  Vermerks) und der Bezug von Antrag und Run (Quelle, Schema, Tabelle) gehören an `slice-backfill-sql-administration`, das
  die Antragsart einführt; dort sind sie als DoD-Punkt „Bezug von Antrag und Run" und als §3-Zeile getragen (Review F-4).
  **Ausgang:** *eingetreten* → Folge-Slice `slice-backfill-sql-administration` (DoD „Bezug von Antrag und Run" und
  §3-Zeile, im Plan des Empfängers gelesen; Verifikation §8 Punkt 1). `Admit` nimmt am Stand `c7045f81` jeden
  `pending`-Antrag über seine Kennung an, unabhängig von Art, Quelle, Schema und Tabelle; bis zum Folge-Slice ist keiner der drei Adapter verdrahtet
  (`git grep` der drei Konstruktoren außerhalb der Tests: nur ihre Definitionen, keine Aufrufstelle; `Admit` ruft allein der Use Case).
- **Gemeldet (Implementer) — Startwerte ohne Messung.** Die Fristen der Adapter (30 s je Zustands-Operation, 5 min je
  Block und Commit, `internal/adapters/driven/postgresstorage/backfill.go`) und die zeilenweise Einfügung eines
  Blocks (je Change eine Anweisung) sind Setzungen; ihre Messung gehört zu `slice-backfill-bench-richtgroesse`.
  **Ausgang:** *weiter offen* — Register `BEO-PGC/backfill-adapter-startwerte-ohne-messung` (1×); Träger der
  Messung ist `slice-backfill-bench-richtgroesse` (§3 „Übergabe aus `slice-backfill-run-store`", in der Closure
  gezogen).

## 7. Closure-Notiz

- **Was hat funktioniert:** die Rollen-Kette lief in getrennten Kontexten und fand, was kein Gate
  las. Der Review (1 HIGH · 1 MEDIUM · 2 LOW · 5 INFO, Summary des Reports) fand den Bestandsdefekt F-1
  mit einer Rechte-Messung an PostgreSQL 18 — seine 18 Mutationen färbten alle rot, prüften aber nicht die
  Rolle der Verbindung (F-2); die Fixrunde (`bd470a0b`,
  `552cb437`) behob F-1 bis F-5 und F-9; die Verifikation bestätigte mit eigenen Läufen (Verifikation §1:
  `make gates`, `make test`, `make a-check`, `make coverage-gate` je EXIT=0; `make test-store` und
  `make test-replication` je EXIT=0 an PostgreSQL 17 und 18), mit der Rechte-Matrix an 17.11 und 18.6
  (§3.2), mit elf Eingabeseiten-Mutationen, alle rot (§4, M1 bis M11b), und mit dem Alt-Tag-Lauf von
  `v0.1.2` (Lauf 5: Exit 0/0/0). Die Adapter-Pflichten aus dem Port-Vertrag des Run-Use-Cases (Zeitbegrenzung
  bei abgelöstem Kontext, Endzustand idempotent) sind am Adapter belegt und rot färbbar (M1, M4); die
  Ordnung `0bf-…` trägt unter vier gemessenen Kollationen (§3.4).
- **Was ging anders als geplant:** der Plan wuchs um drei Sentinel-Fehler an den Ports, `backfill.go`
  (Fristen, Fehlerübersetzung) und die Übersetzung der Run-Zeile in `mapper` und `sqlexec` (beide im
  Coverage-Gate, [`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)); die Fixrunde trug zusätzlich Grant, Prüfung (4b), Login-Tests, Alt-Tag-Rechte,
  Handbuch 1.47 und den Nachzug in `harness/targets/schema-rollout.md` (F-9). Die Fixrunde lief ohne
  Nachprüfungs-Review; der Verifier las, fuhr und mutierte das neue Material selbst (Verifikation V-3,
  §3.5, §4). Die Nachmessung der Closure am Sensor-Träger (Verifikation V-1) misst die beiden Nenner am
  Stand `c7045f81`: DB-Adapter 1016, Unit-Stufe 2398; beide stehen mit Lauf und Commit-Kennung in den
  Sensor-Dokumenten (`510f6ffe`), das committete Suchlauf-Feld trägt den Träger (§3, Suchlauf der Closure).
- **Verifier-Beobachtungen (V-1 bis V-7):** *V-1* (LOW) gezogen: Sensor-Dokumente und Suchlauf-Feld der
  Closure (`510f6ffe`), Risiko 7 mit diesem Lauf. *V-2* gezogen: „Diff sieben" in §3 ist „sechs" (am Stand
  `bd470a0b` gemessen). *V-3* zulässig, kein Nachzug. *V-4* benannt, kein Nachzug: „kein Superuser" der
  Login-Tests ist Konstruktion (`CREATE ROLE … LOGIN … IN ROLE`), keine Assertion; der Verifier maß
  `rolsuper = f`, ein Superuser-Login würde die Tests nicht rot färben. *V-5* benannt: auf den
  Alpine-Testcontainern ist `en_US.utf8` nur ein gedruckter Name, der Beleg trägt Default, `"C"` und ICU
  (Risiko 1). *V-6* im Register (`BEO-PGC/test-schreibt-in-committete-datei`, 3×). *V-7* im Register
  (`BEO-PGC/adr-aussage-breiter-als-ihre-messung`, 3×).
- **Steering-Loop-Eintrag (Lerneintrag):** *Geschärfte Regel (Ausgangs-Kandidat, nicht entschieden):* eine
  Rollen-Zusage — ein Grant, ein Verbot, die Zuordnung eines Adapters zu einer Rolle — braucht einen Test
  unter dem **Login der Rolle** (`CREATE ROLE … LOGIN … IN ROLE <rolle>`, kein Superuser) und eine Mutation
  an der Rolle der Verbindung („Grant streichen"), nicht einen Test unter dem Superuser. Der Slice fand so
  einen Defekt des Bestands: `cdc_admin` trug auf `cdc.administration_request` kein einziges Recht
  (`has_table_privilege` `SELECT` und `UPDATE` beide `f`, vom Reviewer an PostgreSQL 18 gemessen, vom Verifier
  nach der Fixrunde an 17.11 und 18.6 als gesetzt), ein Login `IN ROLE cdc_admin` (Handbuch §2) scheiterte an
  `ListPending` und am Vermerk `applied`; im Compose-Aufbau blieb es unentdeckt, weil `CDC_ADMIN_DSN` dort
  als Superuser `postgres` läuft und jeder Store-Test ebenfalls. Die Klasse ist nicht neu:
  `BEO-PGC/rollen-test-abdeckungsluecken` (Punkt 2, Adapter-Ebene) erreicht mit diesem Slice 3×, der
  Ausgangs-Kandidat steht in dessen `state.md` (Träger-Kandidaten: Schritt 19 von `implement-slice` und
  der HIGH-Punkt „Zusage ohne Bindung an ihre Eingabeseite" im Reviewer-Skill; Architect-Entscheidung im
  Lese-Schritt). Kein Beleg von `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (8×): der Mechanismus
  ist derselbe, die Ausprägung „Rolle der Verbindung" trägt der Rollen-Eintrag. *Neuer Sensor-Umfang:* die
  Rechte sind an drei Ebenen gebunden — der netzlose Rollen-Test der Rollout-Datei (4b, Teil der Gates), die
  Store-Tier-Tests und der Login-Test über den ganzen Administrationsweg (`internal/bootstrap`, mit den
  drei Backfill-Adaptern unter ihren echten Rollen) und die Rechte-Abfragen in Lauf 5 des Guard-Tests nach
  dem Upgrade vom jüngsten `v*`-Tag; die Mutation „Grant streichen" färbt alle rot (Verifikation §4, M11
  und M11b). *Benannte Lücke (Architect):* [`ADR-0047`](../../adr/0047-rollenspezifische-dsn-verdrahtung.md)
  führt den Rollenschnitt von `cdc_admin` ohne die Antrags-Queue,
  [`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) nennt die bestehende
  Schleife „unverändert lauffähig", und die Grants-Tabelle von [`SPEC-029`](../../../../spec/pflichtenheft.md)
  deckt nur `cdc.backfill_run`; beide ADRs sind `Accepted` und unberührbar, ob eine Ergänzung nötig ist (neue
  ADR mit `Supersedes`, Spec-Nachzug), entscheidet der Architect
  (`BEO-PGC/adr-aussage-breiter-als-ihre-messung`, 3×; verwandt, nicht doppelt gezählt:
  `BEO-PGC/architect-verdikt-rollen-scope-luecke`, 1×).
- **Beobachtungs-Register (`../observations/`):** je Vorkommen eine `evidence/`-Datei
  `slice-backfill-run-store.md`, Zähler = Zahl der Dateien (real ausgezählt).
  *Bestehende Klassen:* `BEO-PGC/rollen-test-abdeckungsluecken` (F-1, F-2) **3×**,
  `BEO-PGC/adr-aussage-breiter-als-ihre-messung` (F-1, V-7) **3×** und
  `BEO-PGC/test-schreibt-in-committete-datei` (V-6) **3×** erreichen **mit diesem Slice die Schwelle
  3×**; ihr Ausgang gehört dem Lese-Schritt der Closure von
  [welle-backfill-bestand](../welle-backfill-bestand.md); `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`
  (V-1, V-2) **17×** (verkörpert); `BEO-PGC/test-isolation-geteilter-zustand` (F-3, Dateiordnungs-Kopplung)
  **2×**; `BEO-PGC/zitat-nennt-die-falsche-stelle` (F-3) **8×** (verkörpert). *Neue Klassen (je 1×, offen):*
  `BEO-PGC/ein-instanz-annahme-ohne-erzwingung` (F-5) und `BEO-PGC/backfill-adapter-startwerte-ohne-messung`
  (F-7), für die zwei Risiken mit Ausgang „weiter offen" (§6). *Fortgeführt:*
  `BEO-PGC/adapter-pflicht-eines-ports-ohne-traeger-im-adapter-slice` (1×): sein `state.md` führt den Beleg
  im Adapter (M1, M4). *Benannt, nicht gezählt:* F-4 (Übergabe ohne Träger im Empfänger-Plan, in der Fixrunde
  gezogen; der Gegenstand — eine Aufrufer-Pflicht — liegt neben dem der Klasse
  `adapter-pflicht-eines-ports-ohne-traeger-im-adapter-slice`, eine Adapter-Pflicht), F-6 (der Beleg deckte eine
  Kollation von mehreren, die Messung schließt die Lücke: Risiko 1), F-8 (Coverage-Zahl streut zwischen Läufen,
  die Läufe stehen bei den Zahlen), V-4 (Konstruktion statt Assertion) und die nicht wiederholbaren
  Rollen-Tests des Bestands mit festen Kennungen (F-9, im Beleg von `test-isolation-geteilter-zustand`
  benannt).
- **Träger-Übergaben (§3.13):** die vom Slice gemeldeten Eigenschaften stehen in den Plänen der Empfänger,
  geprüft an deren Text (die Pläne bleiben offen, ihr Planner prüft sie beim Start):
  `slice-backfill-sql-administration` (§3 „Übergaben aus `slice-backfill-run-store`": Rechte der drei Rollen
  auf die zwei Basistabellen, Pools je Adapter, drei Sentinel-Fehler und ihre Klasse, Test-Aufbau mit Logins,
  Handbuch-Stand `Version: 1.47` mit Zeile 1.48 als Zug; §3-Zeile für den Login-Test),
  `slice-transformationen-antragsweg-usecase` (§3: der Regelstand-Adapter läuft über `CDC_ADMIN_DSN`,
  `cdc_capture` trägt auf der Queue kein Recht; §3-Zeile für den Login-Test mit den zwei neuen Arten),
  `slice-transformationen-antragsweg-schema` (§3: die Tabellen-Grants decken die zwei neuen Spalten),
  `slice-backfill-bench-richtgroesse` (§3: Fristen und Einfügeform als Setzungen ohne Messung). Ohne
  Änderung geprüft: `slice-backfill-e2e` und `slice-backfill-sdk-origin` (keine Aussage zu Rollen, Grants
  oder Handbuch-Version; der Compose-Aufbau fährt die DSNs als Superuser und trägt keine Rollen-Aussage),
  `welle-backfill-bestand` und die Roadmap (keine Zustandsaussage berührt; die Register-Zähler in den
  §8-Sichtungslisten der offenen Pläne sind Stände der jeweiligen Planungszeit, der Zähler steht im
  Register).
- **Validator (Modul 8):** entfällt ausdrücklich — der Slice liefert Adapter ohne End-Nutzer-Verhalten:
  der Backfill-Lauf, den ein Nutzer auslösen und lesen könnte, existiert erst mit den Folge-Slices
  (`slice-backfill-sql-administration`, `slice-backfill-e2e`); die drei Adapter sind bis dahin nirgends
  verdrahtet. Die einzige Wirkung am Betrieb ist der Grant `SELECT, UPDATE` von `cdc_admin` auf die
  Antrags-Queue; er ist gegen die Handbuch-Anleitung des Least-Privilege-Logins real belegt
  (`TestAdministrationPathRunsUnderLeastPrivilegeLogins`, Verifikation §3.5). Der Bedarf aus
  [`LH-FA-CAP-009`](../../../../spec/lastenheft.md) wird mit dem Wellen-Beleg (`slice-backfill-e2e`)
  validierbar. Kein stilles Überspringen.
- **Closure-Notiz-Review (`.harness/skills/closure-note-reviewer.md`):** eine getrennte Rolle im frischen
  Kontext, kein Schritt der Planner-Closure; der Skill prüft Slices in `done/` und greift daher erst nach
  dem `git mv` — hier nicht ausgeführt.
- **Folge-Slices:** keine neuen — die Folge-Slices der Welle
  [welle-backfill-bestand](../welle-backfill-bestand.md) liegen als Dateien in `open/`; die Übergaben oben
  tragen ihre Adressen. Die Frage an den Architect zu `ADR-0047`, `ADR-0113` und `SPEC-029` ist kein Slice,
  sie ist im Register adressiert.
- **Risiken aus §6:** je ein Ausgang am Ort — Risiko 1 (Textsortierung) **entfallen**, gemessen; Risiko 2
  (CHECK bei Erstanlage) **entfallen**, Rollout und Alt-Tag-Lauf; Risiko 3 (Fortschritt außerhalb der
  Daten-Transaktion) **entfallen**, Test und Mutation; Risiko 4 (Rollenlage) **entfallen**, Rechte gemessen;
  Risiko 5 (Annahme-Transaktion) **weiter offen** im Register
  (`BEO-PGC/ein-instanz-annahme-ohne-erzwingung`), Atomarität belegt; Risiko 6 (geteilter Zustand)
  **entfallen** für die neuen Tests; Risiko 7 (DB-Adapter-Coverage) **eingetreten**, im Slice gelöst;
  Risiko 8 (Grant auf die Antrags-Queue) **eingetreten**, im Slice behoben, Frage an den Architect im
  Register; Risiko 9 (Antragsart) **eingetreten** → Folge-Slice `slice-backfill-sql-administration`;
  Risiko 10 (Startwerte) **weiter offen** im Register
  (`BEO-PGC/backfill-adapter-startwerte-ohne-messung`).
- **Drei Paarungen:** dieser Slice gehört zu [welle-backfill-bestand](../welle-backfill-bestand.md) (offen) — die
  Prüfung läuft regelkonform bei deren Closure.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area
`*`/`PGC` (Greenfield); Schema, Grants und der Store-Adapter sind keine
eigenen Sub-Areas — kein Anlass zur Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/d-migrate-nacharbeit` (verkörpert, 6×, einschlägig — Rückführung §4),
`BEO-PGC/schema-rollout-fremdobjekte` (verkörpert, 3×, Tabelle bleibt im
neutralen Modell — kein Guard-Eintrag nötig, zu belegen durch den zweiten Lauf),
`BEO-PGC/rollen-verdrahtung` (eingetreten) und
`BEO-PGC/rollen-test-abdeckungsluecken` (offen, 2×, einschlägig — Rollen-Test
im DoD; ein weiterer Beleg erreicht 3×),
`BEO-PGC/test-isolation-geteilter-zustand` (offen, 1×, Risiko §6),
`BEO-PGC/laufzeitzustand-ohne-dauerhaften-traeger` (verkörpert — der
Run-Zustand ist eine Tabelle, kein Prozessspeicher),
`BEO-PGC/db-gegenstand-enthaelt-netzlos-geprueften-code` (offen, 2×, gesichtet
— hier wächst der DB-Gegenstand mit DB-gestütztem Code),
`BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert, 26×, Suchlauf §3).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
