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
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein Self-Review (Modul 8).
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [x] Doku-Update: die Sensor-Zeile `make test-store` in `harness/README.md` und die Schema-Kommentare, soweit sie die Tabellenliste oder die Rollenverteilung aufzählen; das Benutzerhandbuch bleibt bis `sql-administration` unberührt.
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
| `tools/schema/schema.yaml` | update | Tabelle `backfill_run`; CHECK auf `status` bei Erstanlage (spätere Erweiterung über Nacharbeit-SQL, nach dem Muster von `request_kind`); `estimated_rows` nullable; die zwei Warn-Spalten. |
| `tools/schema/nacharbeit-roles.sql` | update | Grants: `cdc_admin` `SELECT`, `INSERT`; `cdc_capture` `SELECT`, `UPDATE`. |
| `internal/bootstrap/roles_rollout_file_internal_test.go` | update | hält den Rollenschnitt gegen die Rollout-Datei; die neue Tabelle darf für `cdc_reader` **nicht** als Basistabellen-Grant erscheinen. |
| `internal/adapters/driven/postgresstorage/` (Arbeitsname `backfilladmission.go`, `backfillrun.go`, `backfillwriter.go`, + Tests) | neu | die drei Adapter über die schmale Ausführungs-Naht des Pakets (`sqlexec`); explizite Spaltenlisten; der Annahme-Adapter nutzt `Begin`. |
| `internal/adapters/driven/postgresstorage/backfillrun.go`, `backfillwriter.go` | Pflicht aus dem Port-Vertrag | Zeitgrenze für `Finish`, `InterruptRunning`, `Rollback` und `Finish` als wirkungsloser Erfolg auf einen beendeten Run (DoD-Punkt „Adapter-Pflichten"). |
| `internal/adapters/driven/postgresstorage/queries/queries.go` | update | die neuen Anweisungen. |
| `tools/schema/plan.yaml`, `tools/schema/down.sql` | regeneriert | Ergebnis von `make schema-rollout`, committet. |
| `tools/harness/run-schema-rollout-guard-test.sh` | ausführen (nicht ändern) | der Alt-Tag-Lauf belegt die Tabelle gegen einen Alt-Bestand ([`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md) Entscheidung 7). |
| `internal/adapters/driven/postgresstorage/schema.sql` | geprüft, nicht geändert | die eingebettete DDL trägt nur die Store-Seite (`source`, `source_table`, `schema_version`, `transaction`, `change`); `cdc.backfill_run` liegt allein im neutralen Modell, die Backfill-Store-Tests laufen vor den Tests, die das Schema per `DROP SCHEMA cdc CASCADE` samt hand-DDL neu aufbauen (Dateinamen `backfill*_test.go` vor `consumerstate_test.go`). |
| `internal/application/port/outbound/backfilladmission.go`, `backfillwriter.go` | update (Plan-Nachzug) | drei Sentinel-Fehler an den Ports, damit die Adapter-Ablehnungen über `errors.Is` unterscheidbar sind: `ErrBackfillRequestNotPending` (die Annahme trifft keine `pending`-Zeile), `ErrBackfillRunInvalid` (Run passt nicht zu `Admit`/`Begin`), `ErrBackfillBlockInvalid` (Block oder Commit verletzt den Vertrag des Schreibers); die Port-Kommentare tragen den geschärften Vertrag (Blocknummern 1, 2, 3, …, Kennungen, Herkunft `backfill`, Zähler = angehängte Changes). Keine Änderung an Methoden-Signaturen. |
| `internal/adapters/driven/postgresstorage/backfill.go` | neu (Plan-Nachzug) | die zwei Frist-Startwerte und die Fehlerübersetzung (`ErrBackfillStorage`) der drei Adapter an einer Stelle. |
| `internal/adapters/driven/postgresstorage/mapper/backfillrun.go`, `sqlexec/translate.go` (+ Tests) | neu/update (Plan-Nachzug) | die Übersetzung der Run-Zeile (`BackfillRunRow`, `ToBackfillRun`, `ReadBackfillRuns`, NULL-Schätzung als „unbekannt“) liegt in den beiden Paketen, die im Coverage-Gate stehen und netzlos geprüft werden — die Logik im DB-Paket bliebe außerhalb des Gates ([`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)); keine neue Paketgrenze, die drei namentlichen DB-Listen bleiben unverändert. |
| `internal/adapters/driven/postgresstorage/backfill_internal_test.go`, `backfillhelpers_test.go`, `backfilladmission_test.go`, `backfillrun_test.go`, `backfillwriter_test.go`, `roles_test.go` | neu/update (Plan-Nachzug) | Store-Tier-Tests der DoD (Annahme, Zustand, Atomarität, Ordnung, Rollen) und der netzlose Frist-/Vertragstest an der Ausführungs-Naht (`backfill_internal_test.go`, läuft in `make test`). |
| `internal/bootstrap/roles_rollout_file_internal_test.go` | update | Prüfung (4a): Grants und Verbote auf `cdc.backfill_run` gegen den Rollout-Text. |
| `harness/README.md` | update | die Sensor-Zeile `make test-store` nennt die Backfill-Tests im Store-Tier. |

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
  **Ausgang:** *(bei Closure)*
- **Der CHECK bei Erstanlage konvergiert** (die Zusage von [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage
  4): die Tabelle `change` trägt CHECKs bei Erstanlage bereits (`chk_change_operation`
  in `tools/schema/schema.yaml`); für **diese** Tabelle ungemessen. *Erwartet, zu
  belegen durch:* Erst- und Folgelauf des Rollouts und der Alt-Tag-Lauf (die
  Tabelle ist additiv gegenüber dem Schema des jüngsten `v*`-Tags; additive
  Änderungen rollen über einen Alt-Bestand, gemessen im Architect-Verdikt
  `architect-verdict-schema-rollout-view-signatur`, Szenario 5, für eine neue
  Tabelle mit Fremdschlüssel und Default — der CHECK und die Warn-Spalten
  gehören nicht zu dieser Messung). **Ausgang:** *(bei Closure)*
- **Der Fortschritt außerhalb der Daten-Transaktion** braucht eine zweite
  Verbindung und kollidiert nicht mit dem Commit derselben Run-Zeile. *Erwartet,
  zu belegen durch:* der Zustands-Test mit gleichzeitigem Fortschritts-Update
  und Commit. **Ausgang:** *(bei Closure)*
- **Die Rollenlage** ist mit [`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 1 gesetzt (`cdc_admin` legt an, `cdc_capture`
  ändert, `cdc_reader` liest nur die View); ein `INSERT` für `cdc_capture`, ein
  `UPDATE` oder `DELETE` für `cdc_admin` oder ein `SELECT`-Grant für `cdc_reader`
  auf die Basistabelle wäre ein Rollen-Verstoß ([`LH-QA-SEC-001`](../../../../spec/lastenheft.md)…[`LH-QA-SEC-003`](../../../../spec/lastenheft.md)).
  *Erwartet, zu belegen durch:* der Rollen-Test im Store-Tier und der netzlose
  Rollen-Test der Rollout-Datei. **Ausgang:** *(bei Closure)*
- **Die Annahme-Transaktion ist nicht atomar oder greift zu weit.** Der
  Antragsvermerk muss genau eine `pending`-Zeile treffen, sonst Rollback; der
  Verzicht auf eine Sperre zwischen mehreren Annehmenden beruht auf der Annahme,
  dass `processAdministrationRequests` Anträge sequenziell in **einer** Goroutine
  verarbeitet und je Quelle eine Instanz Anträge annimmt ([`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 1,
  Punkt 2). *Erwartet, zu belegen durch:* der Rollback-Test des zweiten
  Liefer-Punkts. **Ausgang:** *(bei Closure)*
- **Geteilter Zustand in den Store-Tests** (`BEO-PGC/test-isolation-geteilter-zustand`,
  offen, 1×): unskopierte Bereinigung in den bestehenden Tests würde die Zeilen
  eines parallelen Tests treffen. *Erwartet, zu belegen durch:* die neuen Tests
  bereinigen nur ihre Run-Kennung. **Ausgang:** *(bei Closure)*
- **DB-Adapter-Coverage** ([`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md) Punkt 3): neuer, DB-gestützt geprüfter
  Code in `postgresstorage` bewegt Zähler und Nenner. *Erwartet, zu belegen
  durch:* der Bericht nennt die Zahl mit ihrem Lauf ([`AGENTS.md`](../../../../AGENTS.md) §3.12 Instanz A).
  **Ausgang:** *(bei Closure)*
- **Gemeldet (Implementer) — `cdc_admin` trägt kein Recht auf `cdc.administration_request`.**
  Gemessen am ausgerollten Schema (`has_table_privilege('cdc_admin', 'cdc.administration_request', 'SELECT')` und
  `'UPDATE'`: beide `false`; `SET ROLE cdc_admin; SELECT count(*) FROM cdc.administration_request` endet mit
  „permission denied for table administration_request"). `Admit` vermerkt den Antrag als `cdc_admin`
  ([`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 1) und braucht dafür
  `SELECT` und `UPDATE`; `ListPending`/`MarkApplied` laufen laut `internal/bootstrap/wiring.go` (`NewAdministrationRequest` über `cfg.AdminDSN`) über denselben Pool — gelesen, nicht an einem Login gemessen. Die Store-Tests dieses
  Slice laufen mit dem Superuser des Testcontainers, der Rollen-Test deckt nur `cdc.backfill_run` — die Lücke ist damit
  nicht durch einen Test rot, sondern gemessen. Vorschlag (Architect-Entscheidung, Tabelle in
  [`ADR-0047`](../../adr/0047-rollenspezifische-dsn-verdrahtung.md)):
  `GRANT SELECT, UPDATE ON cdc.administration_request TO cdc_admin` samt Rollen-Test; Adresse:
  `slice-backfill-sql-administration` (dort läuft die Verdrahtung) oder ein eigener Slice. **Ausgang:** *(bei Closure)*
- **Gemeldet (Implementer) — der Annahme-Adapter liest die Antragsart nicht.** Die geschlossene `request_kind`-Menge
  trägt `backfill` noch nicht (`tools/schema/nacharbeit-administration.sql`); die Store-Tests nehmen deshalb einen
  offenen Antrag der Art `enable` an. Der Art-Vergleich (`request_kind = 'backfill'` in der `WHERE`-Klausel des
  Vermerks) gehört an `slice-backfill-sql-administration`, das die Antragsart einführt. **Ausgang:** *(bei Closure)*
- **Gemeldet (Implementer) — Startwerte ohne Messung.** Die Fristen der Adapter (30 s je Zustands-Operation, 5 min je
  Block und Commit, `internal/adapters/driven/postgresstorage/backfill.go`) und die zeilenweise Einfügung eines
  Blocks (je Change eine Anweisung) sind Setzungen; ihre Messung gehört zu `slice-backfill-bench-richtgroesse`.
  **Ausgang:** *(bei Closure)*

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
